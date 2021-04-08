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
	"encoding/json"
	"fmt"
	conn_api_client "github.com/odahu/odahu-flow/packages/operator/pkg/apiclient/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/inspectors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/odahuflow"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/http"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/kubernetes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	"go.uber.org/zap"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/hash"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"knative.dev/pkg/apis"
	knservingv1 "knative.dev/serving/pkg/apis/serving/v1"
	"odahu-commons/predictors"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"strconv"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
)

const (
	deploymentControllerName             = "modeldeployment_controller"
	DefaultModelPort                     = int32(5000)
	DefaultRequeueDelay                  = 10 * time.Second
	DefaultPortName                      = "http1"
	KnativeMinReplicasKey                = "autoscaling.knative.dev/minScale"
	KnativeMaxReplicasKey                = "autoscaling.knative.dev/maxScale"
	KnativeAutoscalingTargetKey          = "autoscaling.knative.dev/target"
	KnativeAutoscalingTargetDefaultValue = "10"
	KnativeAutoscalingClass              = "autoscaling.knative.dev/class"
	KnativeAutoscalingMetric             = "autoscaling.knative.dev/metric"
	DefaultKnativeAutoscalingMetric      = "concurrency"
	DefaultKnativeAutoscalingClass       = "kpa.autoscaling.knative.dev"
	ModelNameAnnotationKey               = "modelName"
	AppliedModelDeploymentSpecKey        = "odahu.org/applied-model-deployment-spec"
	AppliedPolicyHashKey                 = "odahu.org/applied-policy-hash"

	IstioRewriteHTTPProbesAnnotation = "sidecar.istio.io/rewriteAppHTTPProbers"
	OdahuAuthorizationLabel          = "odahu-flow-authorization"

	deploymentIDLabel = "odahu.org/deploymentID"
	cmPolicySuffix    = "opa-policy"
	podPolicyLabel    = "opa-policy-config-map-name"
)

var (
	DefaultTerminationPeriod = int64(600)
)

func NewModelDeploymentReconciler(
	mgr manager.Manager,
	cfg config.Config,
	rootLogger *zap.Logger,
) *ModelDeploymentReconciler {
	authCfg := cfg.Operator.Auth

	noPrefixHTTPClient := http.NewBaseAPIClient(
		cfg.ServiceCatalog.EdgeURL, authCfg.APIToken, authCfg.ClientID,
		authCfg.ClientSecret, authCfg.OAuthOIDCTokenEndpoint, "",
	)

	inspectorsMap, err := inspectors.NewInspectorsMap(cfg.ServiceCatalog.EdgeURL, &noPrefixHTTPClient)
	if err != nil {
		panic(err)
	}

	return &ModelDeploymentReconciler{
		Client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
		log:    rootLogger.Sugar().Named(deploymentControllerName),
		connAPIClient: conn_api_client.NewClient(
			cfg.Operator.Auth.APIURL,
			cfg.Operator.Auth.APIToken,
			cfg.Operator.Auth.ClientID,
			cfg.Operator.Auth.ClientSecret,
			cfg.Operator.Auth.OAuthOIDCTokenEndpoint,
		),
		deploymentConfig: cfg.Deployment,
		operatorConfig:   cfg.Operator,
		gpuResourceName:  cfg.Common.ResourceGPUName,
		inspectors:       inspectorsMap,
	}
}

// ModelDeploymentReconciler reconciles a ModelDeployment object
type ModelDeploymentReconciler struct {
	client.Client
	scheme           *runtime.Scheme
	log              *zap.SugaredLogger
	connAPIClient    conn_api_client.Client
	deploymentConfig config.ModelDeploymentConfig
	operatorConfig   config.OperatorConfig
	gpuResourceName  string
	inspectors       map[string]inspectors.ModelServerInspector
}

func KnativeServiceName(md *odahuflowv1alpha1.ModelDeployment) string {
	return md.Name
}

func knativeDeploymentName(revisionName string) string {
	return fmt.Sprintf("%s-deployment", revisionName)
}

func (r *ModelDeploymentReconciler) ReconcileKnativeService(
	log *zap.SugaredLogger,
	modelDeploymentCR *odahuflowv1alpha1.ModelDeployment,
	predictor predictors.Predictor,
) error {
	container, err := r.createModelContainer(modelDeploymentCR, predictor)
	if err != nil {
		return err
	}

	serviceAccountName := ""
	if modelDeploymentCR.Spec.ImagePullConnectionID != nil &&
		*modelDeploymentCR.Spec.ImagePullConnectionID != "" {
		serviceAccountName = odahuflow.GenerateDeploymentConnectionSecretName(modelDeploymentCR.Name)
	}

	var affinity *corev1.Affinity
	if len(modelDeploymentCR.Spec.NodeSelector) == 0 {
		affinity = utils.BuildNodeAffinity(r.deploymentConfig.NodePools)
	}

	fulfilKnativeService := func(kService *knservingv1.Service) error {
		templateLabelsToAdd := map[string]string{
			ModelNameAnnotationKey:  modelDeploymentCR.Name,
			deploymentIDLabel:       modelDeploymentCR.Name,
			OdahuAuthorizationLabel: "enabled",
		}
		if modelDeploymentCR.Spec.RoleName != nil {  // Otherwise default ConfigMap will be mounted to container
			templateLabelsToAdd[podPolicyLabel] = GetCMPolicyName(modelDeploymentCR)
		}
		templateAnnotationsToAdd := map[string]string{
			KnativeAutoscalingClass:          DefaultKnativeAutoscalingClass,
			KnativeAutoscalingMetric:         DefaultKnativeAutoscalingMetric,
			KnativeMinReplicasKey:            strconv.Itoa(int(*modelDeploymentCR.Spec.MinReplicas)),
			KnativeMaxReplicasKey:            strconv.Itoa(int(*modelDeploymentCR.Spec.MaxReplicas)),
			KnativeAutoscalingTargetKey:      KnativeAutoscalingTargetDefaultValue,
			IstioRewriteHTTPProbesAnnotation: "true",
			// Annotation to trigger pod restart if policy is changed
			AppliedPolicyHashKey: modelDeploymentCR.Annotations[AppliedPolicyHashKey],
		}
		revisionSpec := knservingv1.RevisionSpec{
			TimeoutSeconds: &DefaultTerminationPeriod,
			PodSpec: corev1.PodSpec{
				ServiceAccountName: serviceAccountName,
				Containers: []corev1.Container{
					*container,
				},
				NodeSelector: modelDeploymentCR.Spec.NodeSelector,
				Tolerations:  r.deploymentConfig.Tolerations,
				Affinity:     affinity,
			},
		}

		appliedMDSpecJSON, err := json.Marshal(modelDeploymentCR.Spec)
		if err != nil {
			return err
		}
		serviceAnnotationsToAdd := map[string]string{
			// Annotation to detect changes in MD
			AppliedModelDeploymentSpecKey: string(appliedMDSpecJSON),
		}

		revisionTemplate := kService.Spec.ConfigurationSpec.Template

		if revisionTemplate.Labels == nil {
			revisionTemplate.Labels = make(map[string]string, len(templateLabelsToAdd))
		}
		for k, v := range templateLabelsToAdd {
			revisionTemplate.Labels[k] = v
		}

		if revisionTemplate.Annotations == nil {
			revisionTemplate.Annotations = make(map[string]string, len(templateAnnotationsToAdd))
		}
		for k, v := range templateAnnotationsToAdd {
			revisionTemplate.Annotations[k] = v
		}
		revisionTemplate.Spec = revisionSpec
		kService.Spec.ConfigurationSpec.Template = revisionTemplate

		if kService.Annotations == nil {
			kService.Annotations = make(map[string]string, len(serviceAnnotationsToAdd))
		}
		for k, v := range serviceAnnotationsToAdd {
			kService.Annotations[k] = v
		}

		return nil
	}

	knativeServiceName := KnativeServiceName(modelDeploymentCR)

	found := &knservingv1.Service{}
	err = r.Get(context.TODO(), types.NamespacedName{
		Name: knativeServiceName, Namespace: modelDeploymentCR.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {
		newKService := &knservingv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      knativeServiceName,
				Namespace: modelDeploymentCR.Namespace,
			},
		}
		err = fulfilKnativeService(newKService)
		if err != nil {
			return err
		}

		if err := controllerutil.SetControllerReference(modelDeploymentCR, newKService, r.scheme); err != nil {
			return err
		}

		log.Info(fmt.Sprintf("Creating %s k8s Knative Service, applied MD spec: %s",
			newKService.ObjectMeta.Name, newKService.Annotations[AppliedModelDeploymentSpecKey]))
		err = r.Create(context.TODO(), newKService)
		return err
	} else if err != nil {
		return err
	}

	lastAppliedMDSpec := &odahuflowv1alpha1.ModelDeploymentSpec{}
	err = json.Unmarshal([]byte(found.Annotations[AppliedModelDeploymentSpecKey]), lastAppliedMDSpec)
	if err != nil {
		return err
	}

	depSpecChanged := !reflect.DeepEqual(*lastAppliedMDSpec, modelDeploymentCR.Spec)

	oldPolicyHash := found.Spec.ConfigurationSpec.Template.Annotations[AppliedPolicyHashKey]
	newPolicyHash := modelDeploymentCR.Annotations[AppliedPolicyHashKey]
	policyIsChanged := newPolicyHash != oldPolicyHash


	if depSpecChanged || policyIsChanged {
		if policyIsChanged {
			log.Info("Policy hash was changed",
				"old", oldPolicyHash,
				"new", newPolicyHash,
			)
		}
		if depSpecChanged {
			log.Info("ModelDeployment spec was changed","old", lastAppliedMDSpec,
				"new", modelDeploymentCR.Spec)
		}

		err = fulfilKnativeService(found)
		if err != nil {
			return err
		}
		log.Info(fmt.Sprintf("Updating '%s' Knative Service, new MD generation: %d",
			knativeServiceName, modelDeploymentCR.Generation))
		err = r.Update(context.TODO(), found)
		if err != nil {
			return err
		}

		return nil
	}
	log.Info("ModelDeployment Spec is not changed. Policy is not changed. Skipping.")
	return nil

}

func (r *ModelDeploymentReconciler) createModelContainer(
	modelDeploymentCR *odahuflowv1alpha1.ModelDeployment,
	predictor predictors.Predictor,
) (*corev1.Container, error) {

	depResources, err := kubernetes.ConvertOdahuflowResourcesToK8s(modelDeploymentCR.Spec.Resources, r.gpuResourceName)
	if err != nil {
		return nil, err
	}

	// Merge Probes from predictor and MD Spec
	livenessProbe := predictor.LivenessProbe
	livenessProbe.InitialDelaySeconds = *modelDeploymentCR.Spec.LivenessProbeInitialDelay
	readinessProbe := predictor.ReadinessProbe
	readinessProbe.InitialDelaySeconds = *modelDeploymentCR.Spec.ReadinessProbeInitialDelay

	return &corev1.Container{
		Image:          modelDeploymentCR.Spec.Image,
		Resources:      depResources,
		Ports:          predictor.Ports,
		LivenessProbe:  &livenessProbe,
		ReadinessProbe: &readinessProbe,
	}, nil
}

// Shortcut to get corresponding Knative Service
func (r *ModelDeploymentReconciler) getKnativeService(
	modelDeploymentCR *odahuflowv1alpha1.ModelDeployment,
) (*knservingv1.Service, error) {
	knativeService := &knservingv1.Service{}

	err := r.Get(context.TODO(), types.NamespacedName{
		Name: KnativeServiceName(modelDeploymentCR), Namespace: modelDeploymentCR.Namespace,
	}, knativeService)

	return knativeService, err
}

func (r *ModelDeploymentReconciler) reconcileStatus(
	log *zap.SugaredLogger, modelDeploymentCR *odahuflowv1alpha1.ModelDeployment,
	state odahuflowv1alpha1.ModelDeploymentState) error {

	modelDeploymentCR.Status.State = state

	//if len(latestReadyRevision) != 0 {
	//	modelDeploymentCR.Status.ServiceURL = fmt.Sprintf(
	//		"%s.%s.svc.cluster.local", modelDeploymentCR.Name, modelDeploymentCR.Namespace,
	//	)
	//	modelDeploymentCR.Status.LastRevisionName = latestReadyRevision
	//}

	if err := r.Update(context.TODO(), modelDeploymentCR); err != nil {
		log.Error(err, fmt.Sprintf(
			"Update status of %s model deployment custom resource", modelDeploymentCR.Name,
		))
		return err
	}

	return nil
}

func GetCMPolicyName(modelDeploymentCR *odahuflowv1alpha1.ModelDeployment) string {
	return modelDeploymentCR.Name + "-" + cmPolicySuffix
}

func (r *ModelDeploymentReconciler) reconcilePolicyCM(log *zap.SugaredLogger,
	modelDeploymentCR *odahuflowv1alpha1.ModelDeployment, predictor predictors.Predictor) error {

	rn := modelDeploymentCR.Spec.RoleName

	roleNameIsNotSet := rn == nil

	if roleNameIsNotSet {
		log.Info(".Spec.RoleName is nil")
		// We should delete custom policy if roleName is not set
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:     GetCMPolicyName(modelDeploymentCR),
				Namespace: r.deploymentConfig.Namespace,
			},
		}
		err := r.Delete(context.TODO(), cm)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}

		log.Info("ConfigMap with polices was maybe deleted")

		// Set special hash value to trigger Pod recreation
		noCustomPolicyHash, err := hash.Hash("Role name is not set")
		if err != nil {
			log.Error(err, "Unable to produce configmap hash using policy configmap data")
			return err
		}
		log.Info("Setting no custom policy hash to ModelDeployment Annotation")
		modelDeploymentCR.Annotations[AppliedPolicyHashKey] = strconv.FormatUint(noCustomPolicyHash, 10)

		return nil
	}

	// Handle case when roleName is set

	policies, err := deployment.ReadDefaultPoliciesAndRender(*rn, predictor.OpaPolicyFilename)
	if err != nil {
		return err
	}

	cm := deployment.BuildDefaultPolicyConfigMap(
		GetCMPolicyName(modelDeploymentCR), r.deploymentConfig.Namespace, policies,
	)

	policyHash, err := hash.Hash(cm.Data)
	if err != nil {
		log.Error(err, "Unable to produce configmap hash using policy configmap data")
		return err
	}
	log.Info("Setting policy hash to ModelDeployment Annotation")
	modelDeploymentCR.Annotations[AppliedPolicyHashKey] = strconv.FormatUint(policyHash, 10)

	if err := controllerutil.SetControllerReference(modelDeploymentCR, cm, r.scheme); err != nil {
		return err
	}

	foundCM := &corev1.ConfigMap{}
	err = r.Get(context.TODO(), types.NamespacedName{
		Name: cm.Name, Namespace: cm.Namespace,
	}, foundCM)

	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating config map", "ID", cm.Name)
		return r.Create(context.TODO(), cm)
	}

	if !reflect.DeepEqual(foundCM.Data, cm.Data) {
		log.Info("Policy was changed. Updating config map")
		return r.Update(context.TODO(), cm)
	}

	if err != nil {
		return err
	}
	return nil
}

func (r *ModelDeploymentReconciler) reconcileModelMeta(log *zap.SugaredLogger,
	modelDeploymentCR *odahuflowv1alpha1.ModelDeployment) error {

	inspector := r.inspectors[modelDeploymentCR.Spec.Predictor]
	servedModel, err := inspector.Inspect("", modelDeploymentCR.Status.HostHeader, log)
	if err != nil {
		return err
	}

	modelDeploymentCR.Status.ModelName = servedModel.Metadata.ModelName
	modelDeploymentCR.Status.ModelVersion = servedModel.Metadata.ModelVersion
	return nil
}

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services/status,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=endpoints,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=endpoints/status,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=authentication.istio.io,resources=policies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=authentication.istio.io,resources=policies/status,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=odahuflow.odahu.org,resources=modeldeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=odahuflow.odahu.org,resources=modeldeployments/status,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=serving.knative.dev,resources=configurations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=serving.knative.dev,resources=revisions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=serving.knative.dev,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=serving.knative.dev,resources=routes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.internal.knative.dev,resources=certificates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.internal.knative.dev,resources=serverlessservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.internal.knative.dev,resources=clusteringresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=caching.internal.knative.dev,resources=images,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=autoscaling.internal.knative.dev,resources=podautoscalers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.istio.io,resources=envoyfilters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=secrets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=serviceaccounts/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=odahuflow.odahu.org,resources=connections,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=odahuflow.odahu.org,resources=connections/status,verbs=get;update;patch

func (r *ModelDeploymentReconciler) Reconcile(request ctrl.Request) (ctrl.Result, error) {
	log := r.log.With(odahuflow.ModelDeploymentIDLogPrefix, request.Name)

	// Fetch the ModelDeployment modelDeploymentCR
	modelDeploymentCR := &odahuflowv1alpha1.ModelDeployment{}
	err := r.Get(context.TODO(), request.NamespacedName, modelDeploymentCR)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	for _, finalizer := range modelDeploymentCR.ObjectMeta.Finalizers {
		if finalizer == metav1.FinalizerDeleteDependents {
			log.Info(fmt.Sprintf("Found %s finalizer. Skip reconciling", metav1.FinalizerDeleteDependents))

			if err := r.reconcileStatus(log, modelDeploymentCR, odahuflowv1alpha1.ModelDeploymentStateDeleting); err != nil {
				log.Error(err, "Set deleting deployment state")
				return reconcile.Result{}, err
			}

			return reconcile.Result{}, nil
		}
	}

	log.Info("Run reconciling of model deployment")
	if modelDeploymentCR.Annotations == nil {
		modelDeploymentCR.Annotations = make(map[string]string)
	}

	predictor, ok := predictors.Predictors[modelDeploymentCR.Spec.Predictor]
	if !ok {
		return reconcile.Result{}, fmt.Errorf("unknown predictor %s", modelDeploymentCR.Spec.Predictor)
	}

	if err := r.reconcilePolicyCM(log, modelDeploymentCR, predictor); err != nil {
		log.Error(err, "Reconcile policy config map")
		return reconcile.Result{}, err
	}

	if err := r.reconcileDeploymentPullConnection(log, modelDeploymentCR); err != nil {
		log.Error(err, "Reconcile deployment pull connection")

		return reconcile.Result{}, nil
	}

	if err := r.ReconcileKnativeService(log, modelDeploymentCR, predictor); err != nil {
		log.Error(err, "Reconcile Knative Service")
		return reconcile.Result{}, err
	}

	knService, err := r.getKnativeService(modelDeploymentCR)
	if err != nil {
		log.Error(err, "Getting latest revision")
		return reconcile.Result{}, err
	}

	condition := knService.Status.GetCondition(apis.ConditionReady)
	if condition == nil || condition.Status != corev1.ConditionTrue {
		log.Info("Knative Service is not ready, re-schedule...",
			"Model Deployment Name", modelDeploymentCR.Name)

		err := r.reconcileStatus(log, modelDeploymentCR, odahuflowv1alpha1.ModelDeploymentStateProcessing)
		if err != nil {
			log.Error(err, "failed to update MD status")
		}

		return reconcile.Result{RequeueAfter: DefaultRequeueDelay}, nil
	}

	latestReadyRevision := knService.Status.LatestReadyRevisionName
	modelDeployment := &appsv1.Deployment{}
	modelDeploymentKey := types.NamespacedName{
		Name:      knativeDeploymentName(latestReadyRevision),
		Namespace: r.deploymentConfig.Namespace,
	}

	if err := r.Client.Get(context.TODO(), modelDeploymentKey, modelDeployment); errors.IsNotFound(err) {
		log.Info("Knative Revision Deployment not found. Re-queueing...")
		_ = r.reconcileStatus(log, modelDeploymentCR, odahuflowv1alpha1.ModelDeploymentStateProcessing)

		return reconcile.Result{RequeueAfter: DefaultRequeueDelay}, nil
	} else if err != nil {
		log.Error(err, "Getting of model deployment", "k8s_deployment_name", modelDeploymentKey.Name)

		return reconcile.Result{}, err
	}

	modelDeploymentCR.Status.Replicas = modelDeployment.Status.Replicas
	modelDeploymentCR.Status.AvailableReplicas = modelDeployment.Status.AvailableReplicas
	modelDeploymentCR.Status.Deployment = modelDeployment.Name
	modelDeploymentCR.Status.HostHeader = knService.Status.URL.Host

	if modelDeploymentCR.Status.Replicas != modelDeploymentCR.Status.AvailableReplicas {
		log.Info(fmt.Sprintf("Not enough replicas running. Requeue after %s", DefaultRequeueDelay))
		_ = r.reconcileStatus(log, modelDeploymentCR, odahuflowv1alpha1.ModelDeploymentStateProcessing)

		return reconcile.Result{RequeueAfter: DefaultRequeueDelay}, nil
	}

	log.Info("Inspecting model metadata...")
	if modelDeploymentCR.Status.AvailableReplicas > 0 {
		err = r.reconcileModelMeta(log, modelDeploymentCR)
		if err != nil {
			log.Errorw("failed to inspect model", "error", err)
			_ = r.reconcileStatus(log, modelDeploymentCR, odahuflowv1alpha1.ModelDeploymentStateProcessing)
			return reconcile.Result{RequeueAfter: DefaultRequeueDelay}, nil
		}
	}

	if err := r.reconcileStatus(
		log,
		modelDeploymentCR,
		odahuflowv1alpha1.ModelDeploymentStateReady,
	); err != nil {
		log.Error(err, "Reconcile Model Deployment Status")

		return reconcile.Result{}, err
	}

	return reconcile.Result{
		Requeue:      true,
		RequeueAfter: PeriodVerifyingDockerConnectionToken,
	}, nil
}

func (r *ModelDeploymentReconciler) SetupBuilder(mgr ctrl.Manager) *ctrl.Builder {
	return ctrl.NewControllerManagedBy(mgr).
		For(&odahuflowv1alpha1.ModelDeployment{}).
		Owns(&knservingv1.Service{}).
		Owns(&odahuflowv1alpha1.ModelRoute{}).
		Owns(&corev1.ConfigMap{}).
		Watches(&source.Kind{Type: &appsv1.Deployment{}}, &EnqueueRequestForImplicitOwner{}).
		Watches(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &odahuflowv1alpha1.Connection{},
		}).
		Watches(&source.Kind{Type: &corev1.ServiceAccount{}}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &odahuflowv1alpha1.Connection{},
		})
}

func (r *ModelDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return r.SetupBuilder(mgr).Complete(r)
}
