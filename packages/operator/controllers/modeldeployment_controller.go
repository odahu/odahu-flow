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
	"github.com/go-logr/logr"
	conn_api_client "github.com/odahu/odahu-flow/packages/operator/pkg/apiclient/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/odahuflow"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/kubernetes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"knative.dev/serving/pkg/apis/serving"
	knservingv1 "knative.dev/serving/pkg/apis/serving/v1"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
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
	defaultLivenessFailureThreshold      = 15
	defaultLivenessPeriod                = 1
	defaultLivenessTimeout               = 1
	defaultReadinessFailureThreshold     = 15
	defaultReadinessPeriod               = 1
	defaultReadinessTimeout              = 1
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
	ModelDeploymentVersionKey            = "modelDeploymentVersion"

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
) *ModelDeploymentReconciler {
	return &ModelDeploymentReconciler{
		Client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
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
	}
}

// ModelDeploymentReconciler reconciles a ModelDeployment object
type ModelDeploymentReconciler struct {
	client.Client
	scheme           *runtime.Scheme
	connAPIClient    conn_api_client.Client
	deploymentConfig config.ModelDeploymentConfig
	operatorConfig   config.OperatorConfig
	gpuResourceName  string
}

func KnativeConfigurationName(md *odahuflowv1alpha1.ModelDeployment) string {
	return md.Name
}

func knativeDeploymentName(revisionName string) string {
	return fmt.Sprintf("%s-deployment", revisionName)
}

func (r *ModelDeploymentReconciler) ReconcileKnativeConfiguration(
	log logr.Logger,
	modelDeploymentCR *odahuflowv1alpha1.ModelDeployment,
	predictor odahuflow.Predictor,
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

	knativeConfiguration := &knservingv1.Configuration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      KnativeConfigurationName(modelDeploymentCR),
			Namespace: modelDeploymentCR.Namespace,
			Annotations: map[string]string{
				ModelDeploymentVersionKey: modelDeploymentCR.ResourceVersion,
			},
		},
		Spec: knservingv1.ConfigurationSpec{
			Template: knservingv1.RevisionTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						ModelNameAnnotationKey: modelDeploymentCR.Name,
						deploymentIDLabel:      modelDeploymentCR.Name,
						// We must provide it because we don't use knative built-in routing
						// And w/o this label revision is considered as "Unreachable" so then minScale is ignored and
						// replicas are scaled to 0 even if user set minScale > 0
						// TODO in future: move from creating knservingv1.Configuration -> knservingv1.Service
						// https://github.com/knative/serving/blob/a333742324081d769d1b234622f3fc4cfd181ca4/pkg/apis/autoscaling/v1alpha1/pa_lifecycle.go#L85
						serving.RouteLabelKey:   modelDeploymentCR.Name,
						OdahuAuthorizationLabel: "enabled",
						podPolicyLabel:          getCMPolicyName(modelDeploymentCR),
					},
					Annotations: map[string]string{
						KnativeAutoscalingClass:          DefaultKnativeAutoscalingClass,
						KnativeAutoscalingMetric:         DefaultKnativeAutoscalingMetric,
						KnativeMinReplicasKey:            strconv.Itoa(int(*modelDeploymentCR.Spec.MinReplicas)),
						KnativeMaxReplicasKey:            strconv.Itoa(int(*modelDeploymentCR.Spec.MaxReplicas)),
						KnativeAutoscalingTargetKey:      KnativeAutoscalingTargetDefaultValue,
						IstioRewriteHTTPProbesAnnotation: "true",
					},
				},
				Spec: knservingv1.RevisionSpec{
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
				},
			},
		},
	}

	if err := controllerutil.SetControllerReference(modelDeploymentCR, knativeConfiguration, r.scheme); err != nil {
		return err
	}

	found := &knservingv1.Configuration{}
	err = r.Get(context.TODO(), types.NamespacedName{
		Name: knativeConfiguration.Name, Namespace: knativeConfiguration.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info(fmt.Sprintf("Creating %s k8s Knative Configuration", knativeConfiguration.ObjectMeta.Name))
		err = r.Create(context.TODO(), knativeConfiguration)
		return err
	} else if err != nil {
		return err
	}

	modelDeploymentVersion := found.Annotations[ModelDeploymentVersionKey]
	if modelDeploymentVersion == modelDeploymentCR.ResourceVersion {
		log.Info("Model Deployment version is up to date, skipping")
		return nil
	}

	log.Info(fmt.Sprintf(
		"Knative Configuration bases on version %s, update to reflect version %s",
		modelDeploymentVersion, modelDeploymentCR.ResourceVersion,
	))

	found.Spec = knativeConfiguration.Spec
	for k, v := range knativeConfiguration.Labels {
		found.Labels[k] = v
	}
	for k, v := range knativeConfiguration.Annotations {
		found.Annotations[k] = v
	}

	log.Info(fmt.Sprintf("Updating %s Knative Configuration", knativeConfiguration.ObjectMeta.Name))
	err = r.Update(context.TODO(), found)
	if err != nil {
		return err
	}

	return nil
}

func (r *ModelDeploymentReconciler) createModelContainer(
	modelDeploymentCR *odahuflowv1alpha1.ModelDeployment,
	predictor odahuflow.Predictor,
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

// Retrieve current configuration and return last revision name.
// If the latest revision name equals the latest created revision, then last deployment changes were applied.
func (r *ModelDeploymentReconciler) getLatestReadyRevision(
	log logr.Logger,
	modelDeploymentCR *odahuflowv1alpha1.ModelDeployment,
) (string, bool, error) {
	knativeConfiguration := &knservingv1.Configuration{}
	if err := r.Get(context.TODO(), types.NamespacedName{
		Name: KnativeConfigurationName(modelDeploymentCR), Namespace: modelDeploymentCR.Namespace,
	}, knativeConfiguration); errors.IsNotFound(err) {
		return "", false, nil
	} else if err != nil {
		log.Error(err, "Getting Knative Configuration")
		return "", false, err
	}

	latestReadyRevisionName := knativeConfiguration.Status.LatestReadyRevisionName
	configurationReady := len(latestReadyRevisionName) != 0 &&
		latestReadyRevisionName == knativeConfiguration.Status.LatestCreatedRevisionName

	return latestReadyRevisionName,
		configurationReady,
		nil
}

func (r *ModelDeploymentReconciler) reconcileStatus(
	log logr.Logger, modelDeploymentCR *odahuflowv1alpha1.ModelDeployment,
	state odahuflowv1alpha1.ModelDeploymentState, latestReadyRevision string) error {

	modelDeploymentCR.Status.State = state

	if len(latestReadyRevision) != 0 {
		modelDeploymentCR.Status.ServiceURL = fmt.Sprintf(
			"%s.%s.svc.cluster.local", modelDeploymentCR.Name, modelDeploymentCR.Namespace,
		)
		modelDeploymentCR.Status.LastRevisionName = latestReadyRevision
	}

	if err := r.Update(context.TODO(), modelDeploymentCR); err != nil {
		log.Error(err, fmt.Sprintf(
			"Update status of %s model deployment custom resource", modelDeploymentCR.Name,
		))
		return err
	}

	return nil
}

// Reconciles a separate kubernetes service for a model deployment with stable name
func (r *ModelDeploymentReconciler) reconcileService(
	log logr.Logger,
	modelDeploymentCR *odahuflowv1alpha1.ModelDeployment,
) error {

	// Ensures that Service has a required port
	fulfilServiceObject := func(service *corev1.Service) (updated bool) {
		requiredPort := corev1.ServicePort{
			Name:       "http",
			Protocol:   corev1.ProtocolTCP,
			Port:       80,
			TargetPort: intstr.FromInt(8012),
		}

		found := false
		for i, port := range service.Spec.Ports {
			if port.Name != requiredPort.Name {
				continue
			}
			found = true
			if port != requiredPort {
				service.Spec.Ports[i] = requiredPort
				updated = true
			}
		}

		if !found {
			service.Spec.Ports = append(service.Spec.Ports, requiredPort)
			updated = true
		}

		return updated
	}

	serviceNamespacedName := types.NamespacedName{
		Name: modelDeploymentCR.Name, Namespace: modelDeploymentCR.Namespace,
	}

	foundService := &corev1.Service{}
	err := r.Get(context.TODO(), serviceNamespacedName, foundService)
	if err != nil && errors.IsNotFound(err) {
		log.Info(fmt.Sprintf("Creating %s k8s service", serviceNamespacedName.Name))

		newService := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      serviceNamespacedName.Name,
				Namespace: serviceNamespacedName.Namespace,
			},
			Spec: corev1.ServiceSpec{
				Type: corev1.ServiceTypeClusterIP,
			},
		}
		fulfilServiceObject(newService)

		if err := controllerutil.SetControllerReference(modelDeploymentCR, newService, r.scheme); err != nil {
			return err
		}

		err = r.Create(context.TODO(), newService)
		return err
	} else if err != nil {
		return err
	}

	needToUpdate := fulfilServiceObject(foundService)

	if !needToUpdate {
		log.Info("K8s Service already has required port, skipping")
		return nil
	}

	log.Info(fmt.Sprintf("Update ports of %s K8s Service...", foundService.Name))

	err = r.Update(context.TODO(), foundService)
	if err != nil {
		return err
	}

	return nil
}

// Syncs subsets endpoints of a model deployment service with subsets endpoints of the latest model knative revision
func (r *ModelDeploymentReconciler) reconcileEndpoints(
	log logr.Logger,
	modelDeploymentCR *odahuflowv1alpha1.ModelDeployment,
) error {
	lastRevisionName := modelDeploymentCR.Status.LastRevisionName
	if len(lastRevisionName) == 0 {
		log.Info("Last revision name is empty")
		return nil
	}

	knativeEndpoints := &corev1.Endpoints{}
	if err := r.Get(context.TODO(), types.NamespacedName{
		Namespace: modelDeploymentCR.Namespace,
		Name:      lastRevisionName,
	}, knativeEndpoints); err != nil {
		log.Error(err, "Cannot get the knative endpoints endpoints",
			"last revision name", lastRevisionName)
		return err
	}

	endpointsNamespacedName := types.NamespacedName{
		Name:      modelDeploymentCR.Name,
		Namespace: modelDeploymentCR.Namespace,
	}

	found := &corev1.Endpoints{}
	err := r.Get(context.TODO(), endpointsNamespacedName, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info(fmt.Sprintf("Creating %s k8s endpoints", endpointsNamespacedName.Name))

		endpoints := &corev1.Endpoints{
			ObjectMeta: metav1.ObjectMeta{
				Name:      endpointsNamespacedName.Name,
				Namespace: endpointsNamespacedName.Namespace,
			},
			Subsets: knativeEndpoints.Subsets,
		}

		if err := controllerutil.SetControllerReference(modelDeploymentCR, endpoints, r.scheme); err != nil {
			return err
		}

		err = r.Create(context.TODO(), endpoints)
		return err
	} else if err != nil {
		return err
	}

	if reflect.DeepEqual(found.Subsets, knativeEndpoints.Subsets) {
		log.Info("Endpoints are associated with up-to-date Model Deployment version, skipping")
		return nil
	}

	log.Info(fmt.Sprintf("Endpoints are not up-to-date with latest Knative revision."+
		"Update the %s endpoints", found.Name))

	found.Subsets = knativeEndpoints.Subsets

	err = r.Update(context.TODO(), found)
	if err != nil {
		return err
	}

	return nil
}

// Cleanup old Knative revisions
// Workaround for https://knative.dev/serving/issues/2720
// TODO: need to upgrade knative
func (r *ModelDeploymentReconciler) cleanupOldRevisions(
	log logr.Logger,
	modelDeploymentCR *odahuflowv1alpha1.ModelDeployment,
) error {
	lastRevisionName := modelDeploymentCR.Status.LastRevisionName
	if len(lastRevisionName) == 0 {
		log.Info("Last revision name is empty")

		return nil
	}

	lastKnativeRevision := &knservingv1.Revision{}
	if err := r.Get(context.TODO(), types.NamespacedName{
		Name: lastRevisionName, Namespace: modelDeploymentCR.Namespace,
	}, lastKnativeRevision); err != nil {
		log.Error(err, "Getting Knative Revision")

		return err
	}

	lastKnativeRevisionGenerationStr, ok := lastKnativeRevision.Labels[serving.ConfigurationGenerationLabelKey]
	if !ok {
		return fmt.Errorf(
			"cannot get the latest knative revision generation: %s",
			lastKnativeRevisionGenerationStr,
		)
	}

	lastKnativeRevisionGeneration, err := strconv.Atoi(lastKnativeRevisionGenerationStr)
	if err != nil {
		return err
	}

	knativeRevisions := &knservingv1.RevisionList{}

	labelSelectorReq, err := labels.NewRequirement(
		ModelNameAnnotationKey,
		selection.DoubleEquals,
		[]string{modelDeploymentCR.Name},
	)
	if err != nil {
		log.Error(
			err,
			"Creation of the label selector requirement",
		)
		return err
	}

	labelSelector := labels.NewSelector()
	labelSelector.Add(*labelSelectorReq)

	if err := r.List(context.TODO(), knativeRevisions, &client.ListOptions{
		LabelSelector: labelSelector,
		Namespace:     modelDeploymentCR.Namespace,
	}); err != nil {
		log.Error(err, "Get the list of knative revisions")

		return err
	}

	for _, kr := range knativeRevisions.Items {
		// pin variable
		kr := kr

		modelDeploymentName, ok := kr.Labels[ModelNameAnnotationKey]
		if !ok || modelDeploymentName != modelDeploymentCR.Name {
			continue
		}

		krGenerationStr, ok := kr.Labels[serving.ConfigurationGenerationLabelKey]
		if !ok {
			return fmt.Errorf("cannot get the latest knative revision generation: %s", kr.Name)
		}

		krGeneration, err := strconv.Atoi(krGenerationStr)
		if err != nil {
			return err
		}

		if krGeneration < lastKnativeRevisionGeneration {
			if err := r.Delete(context.TODO(), &kr); err != nil {
				log.Error(err, "Delete old knative revision",
					"knative revision name", kr.Name,
					"knative revision generation", krGeneration)
				return err
			}

			log.Info("Delete old knative revision",
				"model deployment id", modelDeploymentCR.Name,
				"knative revision name", kr.Name,
				"knative revision generation", krGeneration)
		}
	}

	return nil
}

func getCMPolicyName(modelDeploymentCR *odahuflowv1alpha1.ModelDeployment) string {
	return modelDeploymentCR.Name + "-" + cmPolicySuffix
}

func (r *ModelDeploymentReconciler) reconcilePolicyCM(log logr.Logger,
	modelDeploymentCR *odahuflowv1alpha1.ModelDeployment, predictor odahuflow.Predictor) error {

	rn := modelDeploymentCR.Spec.RoleName
	if rn == nil {
		log.Info(".Spec.RoleName is nil. Skip creating custom policies for model")
		return nil
	}

	policies, err := deployment.ReadDefaultPoliciesAndRender(*rn, predictor.OpaPolicyFilename)
	if err != nil {
		return err
	}

	cm := deployment.BuildDefaultPolicyConfigMap(
		getCMPolicyName(modelDeploymentCR), r.deploymentConfig.Namespace, policies,
	)

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

	if err != nil {
		return err
	}
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
	log := logf.Log.WithName(deploymentControllerName).WithValues(odahuflow.ModelDeploymentIDLogPrefix, request.Name)

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

			if err := r.reconcileStatus(log, modelDeploymentCR, odahuflowv1alpha1.ModelDeploymentStateDeleting, ""); err != nil {
				log.Error(err, "Set deleting deployment state")
				return reconcile.Result{}, err
			}

			return reconcile.Result{}, nil
		}
	}

	log.Info("Run reconciling of model deployment")

	predictor, ok := odahuflow.Predictors[modelDeploymentCR.Spec.Predictor]
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

	if err := r.ReconcileKnativeConfiguration(log, modelDeploymentCR, predictor); err != nil {
		log.Error(err, "Reconcile Knative Configuration")
		return reconcile.Result{}, err
	}

	latestReadyRevision, configurationReady, err := r.getLatestReadyRevision(log, modelDeploymentCR)
	if err != nil {
		log.Error(err, "Getting latest revision")
		return reconcile.Result{}, err
	}

	if !configurationReady {
		log.Info("Configuration was not updated yet. Update Status and Put the Model Deployment back in the queue")

		_ = r.reconcileStatus(log, modelDeploymentCR, odahuflowv1alpha1.ModelDeploymentStateProcessing, "")

		return reconcile.Result{RequeueAfter: DefaultRequeueDelay}, nil
	}

	log.Info("Reconcile K8s Service")

	if err := r.reconcileService(log, modelDeploymentCR); err != nil {
		log.Error(err, "Reconcile the k8s service")
		return reconcile.Result{}, err
	}

	if err := r.reconcileEndpoints(log, modelDeploymentCR); err != nil {
		log.Error(err, "Can not reconcile endpoints")
		return reconcile.Result{}, err
	}

	modelDeployment := &appsv1.Deployment{}
	modelDeploymentKey := types.NamespacedName{
		Name:      knativeDeploymentName(latestReadyRevision),
		Namespace: r.deploymentConfig.Namespace,
	}

	if err := r.Client.Get(context.TODO(), modelDeploymentKey, modelDeployment); errors.IsNotFound(err) {
		log.Info("Knative Revision Deployment not found. Re-queueing...")
		_ = r.reconcileStatus(log, modelDeploymentCR, odahuflowv1alpha1.ModelDeploymentStateProcessing, latestReadyRevision)

		return reconcile.Result{RequeueAfter: DefaultRequeueDelay}, nil
	} else if err != nil {
		log.Error(err, "Getting of model deployment", "k8s_deployment_name", modelDeploymentKey.Name)

		return reconcile.Result{}, err
	}

	modelDeploymentCR.Status.Replicas = modelDeployment.Status.Replicas
	modelDeploymentCR.Status.AvailableReplicas = modelDeployment.Status.AvailableReplicas
	modelDeploymentCR.Status.Deployment = modelDeployment.Name

	if modelDeploymentCR.Status.Replicas != modelDeploymentCR.Status.AvailableReplicas {
		log.Info(fmt.Sprintf("Not enough replicas running. Requeue after %s", DefaultRequeueDelay))
		_ = r.reconcileStatus(log, modelDeploymentCR, odahuflowv1alpha1.ModelDeploymentStateProcessing, latestReadyRevision)

		return reconcile.Result{RequeueAfter: DefaultRequeueDelay}, nil
	}

	if err := r.reconcileStatus(
		log,
		modelDeploymentCR,
		odahuflowv1alpha1.ModelDeploymentStateReady,
		latestReadyRevision,
	); err != nil {
		log.Error(err, "Reconcile Model Deployment Status")

		return reconcile.Result{}, err
	}

	if err := r.cleanupOldRevisions(log, modelDeploymentCR); err != nil {
		log.Error(err, "Cleanup old revisions")
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
		Owns(&knservingv1.Configuration{}).
		Owns(&odahuflowv1alpha1.ModelRoute{}).
		Owns(&corev1.ConfigMap{}).
		Watches(&source.Kind{Type: &appsv1.Deployment{}}, &EnqueueRequestForImplicitOwner{}).
		Watches(&source.Kind{Type: &knservingv1.Revision{}}, &EnqueueRequestForImplicitOwner{}).
		Watches(&source.Kind{Type: &corev1.Endpoints{}}, &EnqueueRequestForImplicitOwner{}).
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
