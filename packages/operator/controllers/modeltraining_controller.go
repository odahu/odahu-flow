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

package controllers

import (
	"context"
	"fmt"
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	train_api_client "github.com/odahu/odahu-flow/packages/operator/pkg/apiclient/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	kube_client "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/trainingclient"
	"github.com/odahu/odahu-flow/packages/operator/pkg/odahuflow"
	tektonv1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

const (
	controllerName   = "modeltraining_controller"
	evictedPodReason = "Evicted"
)

const (
	trainingIDLabel = "odahu.org/trainingID"
)

var log = logf.Log.WithName(controllerName)

// ModelTrainingReconciler reconciles a ModelTraining object
type ModelTrainingReconciler struct {
	client.Client
	scheme          *runtime.Scheme
	k8sConfig       *rest.Config
	trainKubeClient kube_client.Client
	trainAPIClient  train_api_client.Client
	trainingConfig  config.ModelTrainingConfig
	operatorConfig  config.OperatorConfig
	gpuResourceName string
}

// newReconciler returns a new reconcile.Reconciler
func NewModelTrainingReconciler(
	mgr manager.Manager, cfg config.Config, trainAPIClient train_api_client.Client,
) *ModelTrainingReconciler {

	k8sClient := mgr.GetClient()

	return &ModelTrainingReconciler{
		Client:    k8sClient,
		k8sConfig: mgr.GetConfig(),
		scheme:    mgr.GetScheme(),
		trainKubeClient: kube_client.NewClient(
			cfg.Training.Namespace,
			cfg.Training.ToolchainIntegrationNamespace,
			k8sClient,
			mgr.GetConfig(),
		),
		trainAPIClient:  trainAPIClient,
		trainingConfig:  cfg.Training,
		operatorConfig:  cfg.Operator,
		gpuResourceName: cfg.Common.ResourceGPUName,
	}
}

func (r *ModelTrainingReconciler) SetupBuilder(mgr ctrl.Manager) *ctrl.Builder {
	return ctrl.NewControllerManagedBy(mgr).
		For(&odahuflowv1alpha1.ModelTraining{}).
		Owns(&corev1.Pod{}).
		Owns(&tektonv1beta1.TaskRun{})
}

func (r *ModelTrainingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return r.SetupBuilder(mgr).Complete(r)
}

const (
	mtConfig = "mt.json"
)

// Determine crd state by child pod.
// If pod has RUNNING state then we determine crd state by state of trainer container in the pod
func (r *ModelTrainingReconciler) syncCrdState(
	taskRun *tektonv1beta1.TaskRun,
	trainingCR *odahuflowv1alpha1.ModelTraining,
) error {
	if len(taskRun.Status.Conditions) > 0 {
		if err := r.calculateStateByTaskRun(taskRun, trainingCR); err != nil {
			return err
		}
	} else {
		trainingCR.Status.State = odahuflowv1alpha1.ModelTrainingScheduling
	}

	log.Info("Setup training state", "mt_id", trainingCR.Name, "state", trainingCR.Status.State)

	trainingCR.Status.PodName = taskRun.Status.PodName

	return r.Update(context.TODO(), trainingCR)
}

func (r *ModelTrainingReconciler) calculateStateByTaskRun(
	taskRun *tektonv1beta1.TaskRun,
	trainingCR *odahuflowv1alpha1.ModelTraining,
) error {
	lastCondition := taskRun.Status.Conditions[len(taskRun.Status.Conditions)-1]

	switch lastCondition.Status {
	case corev1.ConditionUnknown:
		if len(taskRun.Status.PodName) != 0 {
			if err := r.calculateStateByPod(taskRun.Status.PodName, trainingCR); err != nil {
				return err
			}
		} else {
			trainingCR.Status.State = odahuflowv1alpha1.ModelTrainingScheduling
		}
	case corev1.ConditionTrue:
		trainingCR.Status.State = odahuflowv1alpha1.ModelTrainingSucceeded
		trainingCR.Status.Message = &lastCondition.Message
		trainingCR.Status.Reason = &lastCondition.Reason

		result, err := r.trainKubeClient.GetModelTrainingResult(trainingCR.Name)
		if err != nil {
			return err
		}

		trainingCR.Status.Artifacts = append(trainingCR.Status.Artifacts, *result)
	case corev1.ConditionFalse:
		trainingCR.Status.State = odahuflowv1alpha1.ModelTrainingFailed
		trainingCR.Status.Message = &lastCondition.Message
		trainingCR.Status.Reason = &lastCondition.Reason
	default:
		trainingCR.Status.State = odahuflowv1alpha1.ModelTrainingScheduling
	}
	return nil
}

// When tekton task run has the unknown state, we calculate CRD state by pod
func (r *ModelTrainingReconciler) calculateStateByPod(
	trainerPodName string, trainingCR *odahuflowv1alpha1.ModelTraining) error {
	trainerPod := &corev1.Pod{}
	if err := r.Get(
		context.TODO(),
		types.NamespacedName{
			Name:      trainerPodName,
			Namespace: trainingCR.Namespace,
		},
		trainerPod,
	); err != nil {
		return err
	}

	if trainerPod.Status.Reason == evictedPodReason {
		trainingCR.Status.State = odahuflowv1alpha1.ModelTrainingFailed
		trainingCR.Status.Message = &trainerPod.Status.Message

		return nil
	}

	switch trainerPod.Status.Phase {
	case corev1.PodPending:
		trainingCR.Status.State = odahuflowv1alpha1.ModelTrainingScheduling
	case corev1.PodUnknown:
		trainingCR.Status.State = odahuflowv1alpha1.ModelTrainingScheduling
	case corev1.PodRunning:
		trainingCR.Status.State = odahuflowv1alpha1.ModelTrainingRunning
	}

	return nil
}

func (r *ModelTrainingReconciler) getToolchainIntegration(trainingCR *odahuflowv1alpha1.ModelTraining) (
	*training.ToolchainIntegration, error,
) {
	var ti *training.ToolchainIntegration
	ti, err := r.trainAPIClient.GetToolchainIntegration(trainingCR.Spec.Toolchain)
	if err != nil {
		return nil, err
	}
	return &training.ToolchainIntegration{Spec: ti.Spec}, nil
}

func (r *ModelTrainingReconciler) getNodeSelectorAndAffinity(trainingCR *odahuflowv1alpha1.ModelTraining) (
	map[string]string, *corev1.Affinity) {

	// Specific node pool was chosen
	if len(trainingCR.Spec.NodeSelector) > 0 {
		return trainingCR.Spec.NodeSelector, nil
	}

	// Node pool wasn't specified, define affinity to any of training pools
	var nodePools []config.NodePool
	if trainingCR.Spec.IsGPUResourceSet() {
		nodePools = r.trainingConfig.GPUNodePools
	} else {
		nodePools = r.trainingConfig.NodePools
	}

	nodeSelectorTerms := make([]corev1.NodeSelectorTerm, 0, len(nodePools))
	for _, nodePool := range nodePools {
		selector := nodePool.NodeSelector
		matchExpressions := make([]corev1.NodeSelectorRequirement, 0, len(selector))

		for label, value := range selector {
			matchExpressions = append(matchExpressions, corev1.NodeSelectorRequirement{
				Key:      label,
				Operator: corev1.NodeSelectorOpIn,
				Values:   []string{value},
			})
		}

		nodeSelectorTerms = append(nodeSelectorTerms, corev1.NodeSelectorTerm{
			MatchExpressions: matchExpressions,
		})
	}

	return nil, &corev1.Affinity{
		NodeAffinity: &corev1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{NodeSelectorTerms: nodeSelectorTerms},
		},
	}
}

func (r *ModelTrainingReconciler) getTolerations(trainingCR *odahuflowv1alpha1.ModelTraining) []corev1.Toleration {
	var tolerations []corev1.Toleration

	if trainingCR.Spec.IsGPUResourceSet() {
		tolerations = r.trainingConfig.GPUTolerations
	} else {
		tolerations = r.trainingConfig.Tolerations
	}

	return tolerations
}

func (r *ModelTrainingReconciler) reconcileTaskRun(
	trainingCR *odahuflowv1alpha1.ModelTraining,
) (*tektonv1beta1.TaskRun, error) {
	if trainingCR.Status.State != "" && trainingCR.Status.State != odahuflowv1alpha1.ModelTrainingUnknown {
		taskRun := &tektonv1beta1.TaskRun{}
		err := r.Get(context.TODO(), types.NamespacedName{
			Name: trainingCR.Name, Namespace: r.trainingConfig.Namespace,
		}, taskRun)
		if err != nil {
			return nil, err
		}

		log.Info("Training has no unknown state. Skip the task run reconcile",
			"mt id", trainingCR.Name, "state", trainingCR.Status.State)
		return taskRun, nil
	}

	toolchainIntegration, err := r.getToolchainIntegration(trainingCR)
	if err != nil {
		return nil, err
	}

	taskSpec, err := r.generateTrainerTaskSpec(trainingCR, toolchainIntegration)
	if err != nil {
		return nil, err
	}

	nodeSelector, affinity := r.getNodeSelectorAndAffinity(trainingCR)

	taskRun := &tektonv1beta1.TaskRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      trainingCR.Name,
			Namespace: trainingCR.Namespace,
			Labels: map[string]string{
				trainingIDLabel: trainingCR.Name,
			},
		},
		Spec: tektonv1beta1.TaskRunSpec{
			TaskSpec: taskSpec,
			Timeout:  &metav1.Duration{Duration: r.trainingConfig.Timeout},
			PodTemplate: &tektonv1beta1.PodTemplate{
				Tolerations:  r.getTolerations(trainingCR),
				NodeSelector: nodeSelector,
				Affinity:     affinity,
			},
		},
	}

	if err := controllerutil.SetControllerReference(trainingCR, taskRun, r.scheme); err != nil {
		return nil, err
	}

	if err := odahuflow.StoreHash(taskRun); err != nil {
		log.Error(err, "Cannot apply obj hash")
		return nil, err
	}

	found := &tektonv1beta1.TaskRun{}
	err = r.Get(context.TODO(), types.NamespacedName{
		Name: taskRun.Name, Namespace: r.trainingConfig.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {
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

func (r *ModelTrainingReconciler) createResultConfigMap(trainingCR *odahuflowv1alpha1.ModelTraining) error {
	resultCM := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      odahuflow.GenerateTrainingResultCMName(trainingCR.Name),
			Namespace: r.trainingConfig.Namespace,
		},
		Data: map[string]string{},
	}

	if err := controllerutil.SetControllerReference(trainingCR, resultCM, r.scheme); err != nil {
		return err
	}

	if err := odahuflow.StoreHash(resultCM); err != nil {
		log.Error(err, "Cannot apply obj hash")
		return err
	}

	found := &corev1.ConfigMap{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name: resultCM.Name, Namespace: r.trainingConfig.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info(fmt.Sprintf("Creating %s k8s result config map", resultCM.ObjectMeta.Name))
		err = r.Create(context.TODO(), resultCM)
		return err
	}

	return err
}

func isTrainingFinished(mt *odahuflowv1alpha1.ModelTraining) bool {
	state := mt.Status.State

	return state == odahuflowv1alpha1.ModelTrainingSucceeded || state == odahuflowv1alpha1.ModelTrainingFailed
}

// Reconcile reads that state of the cluster for a StorageEntity object and makes changes based on the state read
// and what is in the StorageEntity.Spec
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=pods/exec,verbs=create
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=odahuflow.odahu.org,resources=modeltrainings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=odahuflow.odahu.org,resources=modeltrainings/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=odahuflow.odahu.org,resources=toolchainintegrations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=odahuflow.odahu.org,resources=toolchainintegrations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

func (r *ModelTrainingReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	trainingCR := &odahuflowv1alpha1.ModelTraining{}

	if err := r.Get(context.TODO(), request.NamespacedName, trainingCR); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		log.Error(err, "Cannot fetch CR status")
		return reconcile.Result{}, err
	}

	if isTrainingFinished(trainingCR) {
		log.Info("Training has been finished. Skip reconcile function", "mt id", trainingCR.Name)

		return reconcile.Result{}, nil
	}

	// The configmap is used to save a training result.
	if err := r.createResultConfigMap(trainingCR); err != nil {
		log.Error(err, "Cannot create result config map")

		return reconcile.Result{}, err
	}

	if taskRun, err := r.reconcileTaskRun(trainingCR); err != nil {
		log.Error(err, "Cannot synchronize desired K8S instances state to cluster")

		return reconcile.Result{}, err
	} else if err := r.syncCrdState(taskRun, trainingCR); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
