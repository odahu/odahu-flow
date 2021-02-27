/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/


// TODO:
//	 1. If modelPath is .zip, .gz, .tar then use rclone cp but not sync
//		and then unzip model in validate-model step
//	 2. Pass on helm chart level env variable with tools image name
package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	tektonv1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
)

const (
	batchIDLabel = "odahu.org/batchID"
)

type ConnGetter interface {
	GetConnection(id string) (*connection.Connection, error)
}

type PodGetter interface {
	GetPod(ctx context.Context, name string, namespace string) (corev1.Pod, error)
}


// BatchInferenceJobReconciler reconciles a BatchInferenceJob object
type BatchInferenceJobReconciler struct {
	client.Client
	Log         logr.Logger
	Scheme      *runtime.Scheme
	connAPI     ConnGetter
	cfg         config.BatchConfig
	gpuResName  string
	podGetter   PodGetter
}

type defaultPodGetter struct {
	client.Client
}

func (p defaultPodGetter) GetPod(ctx context.Context, name string, namespace string) (corev1.Pod, error) {
	trainerPod := &corev1.Pod{}
	err := p.Get(
		ctx,
		types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
		trainerPod,
	)
	return *trainerPod, err
}


func setDefaultOptions(options *BatchInferenceJobReconcilerOptions) {
	if options.PodGetter == nil {
		options.PodGetter = defaultPodGetter{options.Client}
	}

}

type BatchInferenceJobReconcilerOptions struct {
	Client client.Client
	Schema *runtime.Scheme
	ConnGetter        ConnGetter
	PodGetter         PodGetter
	Cfg               config.BatchConfig
	ResourceGPUName	  string
}


func NewBatchInferenceJobReconciler(opts BatchInferenceJobReconcilerOptions) *BatchInferenceJobReconciler {

	setDefaultOptions(&opts)

	return &BatchInferenceJobReconciler{
		Client: opts.Client,
		Scheme: opts.Schema,
		podGetter: opts.PodGetter,
		connAPI: opts.ConnGetter,
		cfg: opts.Cfg,
		gpuResName: opts.ResourceGPUName,
		Log: logf.Log.WithName("batch-inference-controller"),
	}
}

// +kubebuilder:rbac:groups=odahuflow.odahu.org,resources=batchinferencejobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=odahuflow.odahu.org,resources=batchinferencejobs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=tekton.dev,resources=taskruns,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tekton.dev,resources=taskruns/status,verbs=get;update;patch

func (r *BatchInferenceJobReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log = r.Log.WithValues("BatchInferenceJobName", req.NamespacedName)
	reconcileLoopUID := uuid.New().String()
	log = log.WithValues("ReconcileLoopUID", reconcileLoopUID)

	batchJob := &odahuflowv1alpha1.BatchInferenceJob{}

	if err := r.Get(ctx, req.NamespacedName, batchJob); err != nil {
		if k8serrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		log.Error(err, "Unable to fetch BatchInferenceJob from Kube API")
		return reconcile.Result{}, err
	}

	// Run new TaskRun or get already created
	tr, err := r.reconcileTaskRun(batchJob,  log)
	if err != nil {
		return ctrl.Result{}, err
	}

	if len(tr.Status.Conditions) > 0 {
		// Fill BatchInferenceJob
		if err := r.syncStatusFromTaskRun(batchJob, tr, log); err != nil {
			log.Error(err,"Unable to sync status for BatchInferenceJob from  TaskRun")
			return ctrl.Result{}, err
		}
	} else {
		batchJob.Status.State = odahuflowv1alpha1.BatchScheduling
	}

	log.Info("Setup batchJob state", "state", batchJob.Status.State)

	batchJob.Status.PodName = tr.Status.PodName

	if err := r.Update(ctx, batchJob); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil

}

func (r *BatchInferenceJobReconciler) generateTaskSpec(
	batchJob *odahuflowv1alpha1.BatchInferenceJob,
	) (*tektonv1beta1.TaskSpec, error) {
	return BatchJobToTaskSpec(
		batchJob, r.connAPI, r.gpuResName, r.cfg.RCloneImage, r.cfg.ToolsSecret, r.cfg.ToolsImage)
}

func (r *BatchInferenceJobReconciler) calculateStateByPod(
	podName string, job *odahuflowv1alpha1.BatchInferenceJob) error {

	pod, err := r.podGetter.GetPod(context.TODO(), podName, job.Namespace)
	if err != nil {
		return err
	}
	job.Status.PodName = podName
	job.Status.Message = ""
	job.Status.Reason = ""

	if pod.Status.Reason == evictedPodReason {
		job.Status.State = odahuflowv1alpha1.BatchFailed
		job.Status.Message = pod.Status.Message
		job.Status.Reason = "Pod evicted"

		return nil
	}

	switch pod.Status.Phase {
	case corev1.PodPending:
		job.Status.State = odahuflowv1alpha1.BatchScheduling
	case corev1.PodUnknown:
		job.Status.State = odahuflowv1alpha1.BatchScheduling
	case corev1.PodRunning:
		job.Status.State = odahuflowv1alpha1.BatchRunning
	}

	return nil
}

func (r *BatchInferenceJobReconciler) syncStatusFromTaskRun(
	batchJob *odahuflowv1alpha1.BatchInferenceJob, taskRun *tektonv1beta1.TaskRun, log logr.Logger) error {
	lastCondition := taskRun.Status.Conditions[len(taskRun.Status.Conditions)-1]

	switch lastCondition.Status {
	case corev1.ConditionUnknown:
		if len(taskRun.Status.PodName) != 0 {
			if err := r.calculateStateByPod(taskRun.Status.PodName, batchJob); err != nil {
				return err
			}
		} else {
			batchJob.Status.State = odahuflowv1alpha1.BatchScheduling
			batchJob.Status.Message = ""
			batchJob.Status.Reason = ""
		}
	case corev1.ConditionTrue:
		batchJob.Status.State = odahuflowv1alpha1.BatchSucceeded
		batchJob.Status.Message = lastCondition.Message
		batchJob.Status.Reason = lastCondition.Reason

	case corev1.ConditionFalse:
		batchJob.Status.State = odahuflowv1alpha1.BatchFailed
		batchJob.Status.Message = lastCondition.Message
		batchJob.Status.Reason = lastCondition.Reason
	default:
		batchJob.Status.State = odahuflowv1alpha1.BatchScheduling
	}
	return nil
}

func (r *BatchInferenceJobReconciler) reconcileTaskRun(
	job *odahuflowv1alpha1.BatchInferenceJob, log logr.Logger,
) (*tektonv1beta1.TaskRun, error) {
	if job.Status.State != "" && job.Status.State != odahuflowv1alpha1.BatchUnknown {
		taskRun := &tektonv1beta1.TaskRun{}
		err := r.Get(context.TODO(), types.NamespacedName{
			Name: job.Name, Namespace: r.cfg.Namespace,
		}, taskRun)
		if err != nil {
			return nil, err
		}

		log.Info("BatchInferenceJob is already finished or running, skip creating TaskRun",
			"mt id", job.Name, "state", job.Status.State)
		return taskRun, nil
	}


	taskSpec, err := r.generateTaskSpec(job)
	if err != nil {
		return nil, err
	}

	var affinity *corev1.Affinity
	if len(job.Spec.NodeSelector) == 0 {
		affinity = utils.BuildNodeAffinity(r.cfg.NodePools)
	}

	taskRun := &tektonv1beta1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      job.Name,
			Namespace: job.Namespace,
			Labels: map[string]string{
				batchIDLabel: job.Name,
			},
		},
		Spec: tektonv1beta1.TaskRunSpec{
			TaskSpec: taskSpec,
			Timeout:  &metav1.Duration{Duration: r.cfg.Timeout},
			PodTemplate: &tektonv1beta1.PodTemplate{
				Tolerations:  r.cfg.Tolerations,
				NodeSelector: job.Spec.NodeSelector,
				Affinity:     affinity,
			},
		},
	}

	if err := controllerutil.SetControllerReference(job, taskRun, r.Scheme); err != nil {
		return nil, err
	}

	found := &tektonv1beta1.TaskRun{}
	err = r.Get(context.TODO(), types.NamespacedName{
		Name: taskRun.Name, Namespace: r.cfg.Namespace,
	}, found)
	if err != nil && k8serrors.IsNotFound(err) {
		log.Info(fmt.Sprintf("Creating %s k8s task run", taskRun.ObjectMeta.Name))
		return taskRun, r.Create(context.TODO(), taskRun)
	} else if err != nil {
		return nil, err
	}

	if err := r.Delete(context.TODO(), found); err != nil {
		return nil, err
	}

	return taskRun, r.Create(context.TODO(), taskRun)
}


func (r *BatchInferenceJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&odahuflowv1alpha1.BatchInferenceJob{}).
		Owns(&tektonv1beta1.TaskRun{}).
		Complete(r)
}
