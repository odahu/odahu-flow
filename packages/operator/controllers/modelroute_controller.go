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
	v1alpha3_istio_api "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	v1alpha3_istio "istio.io/api/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sort"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
)

const (
	knativeRevisionHeader    = "knative-serving-revision"
	knativeNamespaceHeader   = "knative-serving-namespace"
	defaultRetryAttempts     = 30
	defaultListOfRetryCauses = "5xx,connect-failure,refused-stream"
	routeForLabelPrefix      = "odahu-route-for-"
	ModelRouteVersionKey     = "modelRouteVersion"
)

var (
	defaultTimeoutPerTry = gogotypes.DurationProto(time.Second)
)

func NewModelRouteReconciler(
	mgr manager.Manager,
	cfg config.Config,
) *ModelRouteReconciler {
	return &ModelRouteReconciler{
		Client:           mgr.GetClient(),
		scheme:           mgr.GetScheme(),
		deploymentConfig: cfg.Deployment,
		operatorConfig:   cfg.Operator,
		gpuResourceName:  cfg.Common.ResourceGPUName,
	}
}

// ModelRouteReconciler reconciles a ModelRoute object
type ModelRouteReconciler struct {
	client.Client
	scheme           *runtime.Scheme
	deploymentConfig config.ModelDeploymentConfig
	operatorConfig   config.OperatorConfig
	gpuResourceName  string
}

func VirtualServiceName(mr *odahuflowv1alpha1.ModelRoute) string {
	return mr.Name
}

func (r *ModelRouteReconciler) reconcileVirtualService(modelRouteCR *odahuflowv1alpha1.ModelRoute) (bool, error) {
	httpTargets := []*v1alpha3_istio.HTTPRouteDestination{}
	reconileAgain := false

	for _, md := range modelRouteCR.Spec.ModelDeploymentTargets {
		modelDeployment := &odahuflowv1alpha1.ModelDeployment{}
		if err := r.Get(context.TODO(), types.NamespacedName{
			Name: md.Name, Namespace: modelRouteCR.Namespace,
		}, modelDeployment); errors.IsNotFound(err) {
			log.Error(
				err, "Model Deployment is not found",
				"Model Deployment Name", md.Name,
				"Model Route Name", modelRouteCR.Name,
			)

			reconileAgain = true
			continue
		} else if err != nil {
			log.Error(
				err, "Getting of the Model Deployment",
				"Model Deployment Name", md.Name,
				"Model Route Name", modelRouteCR.Name,
			)

			return reconileAgain, err
		}

		if modelDeployment.Status.State != odahuflowv1alpha1.ModelDeploymentStateReady {
			log.Info("Model deployment is not ready", "Model Deployment Name", md.Name, "Model Route Name", modelRouteCR.Name)
			reconileAgain = true

			continue
		}

		requestHeaders := &v1alpha3_istio.Headers_HeaderOperations{
			Add: map[string]string{
				knativeRevisionHeader:  modelDeployment.Status.LastRevisionName,
				knativeNamespaceHeader: r.deploymentConfig.Namespace,
			},
		}

		httpTargets = append(httpTargets,
			&v1alpha3_istio.HTTPRouteDestination{
				Destination: &v1alpha3_istio.Destination{
					Host: modelDeployment.Status.ServiceURL,
					Port: &v1alpha3_istio.PortSelector{
						Number: uint32(80),
					},
				},
				Weight: *md.Weight,
				Headers: &v1alpha3_istio.Headers{
					Request: requestHeaders,
				},
			})
	}

	if len(httpTargets) == 0 {
		log.Info("Number of http targets is zero", "Model Route Name", modelRouteCR.Name)
		return reconileAgain, nil
	}

	var mirror *v1alpha3_istio.Destination
	if modelRouteCR.Spec.Mirror != nil && len(*modelRouteCR.Spec.Mirror) != 0 {
		modelDeployment := &odahuflowv1alpha1.ModelDeployment{}
		if err := r.Get(context.TODO(), types.NamespacedName{
			Name: *modelRouteCR.Spec.Mirror, Namespace: modelRouteCR.Namespace,
		}, modelDeployment); errors.IsNotFound(err) {

		} else if err != nil {
			log.Error(err, fmt.Sprintf("Getting of %s Model Deployment mirror", *modelRouteCR.Spec.Mirror))

			return reconileAgain, err
		}

		if modelDeployment.Status.State != odahuflowv1alpha1.ModelDeploymentStateReady {
			log.Info(
				"Model deployment is not ready",
				"Model Deployment Name", modelRouteCR.Spec.Mirror,
				"Model Route Name", modelRouteCR.Name,
			)

			reconileAgain = true
		} else {
			mirror = &v1alpha3_istio.Destination{
				Host: modelDeployment.Status.ServiceURL,
			}
		}
	}

	vservice := &v1alpha3_istio_api.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      VirtualServiceName(modelRouteCR),
			Namespace: modelRouteCR.Namespace,
			Annotations: map[string]string{
				ModelRouteVersionKey: modelRouteCR.ResourceVersion,
			},
		},
		Spec: v1alpha3_istio_api.VirtualServiceSpec{
			VirtualService: v1alpha3_istio.VirtualService{
				Hosts:    []string{"*"},
				Gateways: []string{"edge"},
				Http: []*v1alpha3_istio.HTTPRoute{
					{
						Retries: &v1alpha3_istio.HTTPRetry{
							Attempts:      defaultRetryAttempts,
							PerTryTimeout: defaultTimeoutPerTry,
							RetryOn:       defaultListOfRetryCauses,
						},
						Match: []*v1alpha3_istio.HTTPMatchRequest{
							{
								Uri: &v1alpha3_istio.StringMatch{
									MatchType: &v1alpha3_istio.StringMatch_Prefix{
										Prefix: modelRouteCR.Spec.URLPrefix + "/",
									},
								},
							},
						},
						Rewrite: &v1alpha3_istio.HTTPRewrite{
							Uri: "/",
						},
						Route:  httpTargets,
						Mirror: mirror,
					},
				},
			},
		},
	}

	if err := controllerutil.SetControllerReference(modelRouteCR, vservice, r.scheme); err != nil {
		return reconileAgain, err
	}

	found := &v1alpha3_istio_api.VirtualService{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name: vservice.Name, Namespace: vservice.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info(fmt.Sprintf("Creating %s k8s Istio Virtual Service", vservice.ObjectMeta.Name))
		err = r.Create(context.TODO(), vservice)
		return reconileAgain, err
	} else if err != nil {
		return reconileAgain, err
	}

	modelDeploymentVersion := found.Annotations[ModelRouteVersionKey]
	if modelDeploymentVersion == modelRouteCR.ResourceVersion {
		log.Info("Istio VService is associated with up-to-date Model Route version, skipping")
		return reconileAgain, err
	}

	log.Info(fmt.Sprintf("Istio Virtual Service hashes don't equal. Update the %s Model route", vservice.Name))

	found.Spec = vservice.Spec
	for k, v := range vservice.Labels {
		found.Labels[k] = v
	}
	for k, v := range vservice.Annotations {
		found.Annotations[k] = v
	}

	log.Info(fmt.Sprintf("Updating %s k8s Istio Virtual Service", vservice.ObjectMeta.Name))
	err = r.Update(context.TODO(), found)
	if err != nil {
		return reconileAgain, err
	}

	return reconileAgain, err
}

func (r *ModelRouteReconciler) reconcileStatus(modelRouteCR *odahuflowv1alpha1.ModelRoute,
	state odahuflowv1alpha1.ModelRouteState) error {
	modelRouteCR.Status.EdgeURL = fmt.Sprintf(
		"%s%s", r.deploymentConfig.Edge.Host, modelRouteCR.Spec.URLPrefix,
	)
	modelRouteCR.Status.State = state

	if err := r.Update(context.TODO(), modelRouteCR); err != nil {
		log.Error(err, "Update status of model deployment custom resource", "Model Deployment Name", modelRouteCR.Name)
		return err
	}

	return nil
}

// +kubebuilder:rbac:groups=odahuflow.odahu.org,resources=modelroutes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=odahuflow.odahu.org,resources=modelroutes/status,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.istio.io,resources=virtualservices,verbs=get;list;watch;create;update;patch;delete

func (r *ModelRouteReconciler) Reconcile(request ctrl.Request) (ctrl.Result, error) {
	modelRouteCR := &odahuflowv1alpha1.ModelRoute{}
	err := r.Get(context.TODO(), request.NamespacedName, modelRouteCR)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	if err = r.reconcileRouteForLabels(modelRouteCR); err != nil {
		return reconcile.Result{}, err
	}

	if reconcileAgain, err := r.reconcileVirtualService(modelRouteCR); err != nil {
		log.Error(err, "Reconcile Istio Virtual Service")
		return reconcile.Result{}, err
	} else if reconcileAgain {
		_ = r.reconcileStatus(modelRouteCR, odahuflowv1alpha1.ModelRouteStateProcessing)

		log.Info("Put the Model Route back in the queue", "Model Route Name", modelRouteCR.Name)
		return reconcile.Result{RequeueAfter: DefaultRequeueDelay}, nil
	}

	if err := r.reconcileStatus(modelRouteCR, odahuflowv1alpha1.ModelRouteStateReady); err != nil {
		log.Info("Reconcile Status of Model Route", "error", err)
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func routeForLabelKey(mdName string) string {
	return routeForLabelPrefix + mdName
}

func (r *ModelRouteReconciler) reconcileRouteForLabels(route *odahuflowv1alpha1.ModelRoute) error {

	var oldRouteForLabels []string
	for key := range route.GetLabels() {
		if strings.HasPrefix(key, routeForLabelPrefix) {
			oldRouteForLabels = append(oldRouteForLabels, key)
		}
	}

	var newRouteForLabels []string //nolint
	for _, md := range route.Spec.ModelDeploymentTargets {
		newRouteForLabels = append(newRouteForLabels, routeForLabelKey(md.Name))
	}
	sort.Strings(oldRouteForLabels)
	sort.Strings(newRouteForLabels)
	if reflect.DeepEqual(oldRouteForLabels, newRouteForLabels) {
		log.Info("ModelRoute routeFor labels are already reconciled properly")
		return nil
	}

	// Otherwise let's update routeFor labels

	newLabels := map[string]string{}
	for key, el := range route.GetLabels() {
		if !strings.HasPrefix(key, routeForLabelPrefix) {
			newLabels[key] = el
		}
	}
	for _, key := range newRouteForLabels {
		newLabels[key] = "enabled"
	}

	route.SetLabels(newLabels)
	if err := r.Update(context.TODO(), route); err != nil {
		return err
	}
	log.Info("ModelRoute routeFor labels were changed")
	return nil
}

func (r *ModelRouteReconciler) SetupBuilder(mgr ctrl.Manager) *ctrl.Builder {
	return ctrl.NewControllerManagedBy(mgr).
		For(&odahuflowv1alpha1.ModelRoute{}).
		Owns(&v1alpha3_istio_api.VirtualService{}).
		// Route should reacts on changes in any ModelDeployment from its targets
		Watches(&source.Kind{Type: &odahuflowv1alpha1.ModelDeployment{}}, &handler.EnqueueRequestsFromMapFunc{
			ToRequests: handler.ToRequestsFunc(func(object handler.MapObject) []reconcile.Request {

				mdName := object.Meta.GetName()

				var result []reconcile.Request
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()

				selector := labels.NewSelector()
				req, err := labels.NewRequirement(routeForLabelKey(mdName), selection.Exists, nil)
				if err != nil {
					log.Error(err, "Unable to create label requirement")
					return result
				}
				selector.Add(*req)

				modelRoutes := &odahuflowv1alpha1.ModelRouteList{}
				if err := r.List(ctx, modelRoutes, &client.ListOptions{
					LabelSelector: selector,
				}); err != nil {
					log.Error(err, "Unable to fetch ModelRouteList for ModelDeployment")
					return result
				}

				for _, mr := range modelRoutes.Items {
					result = append(result, reconcile.Request{
						NamespacedName: types.NamespacedName{
							Namespace: mr.Namespace,
							Name:      mr.Name,
						},
					})
				}

				return result
			}),
		})
}

func (r *ModelRouteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return r.SetupBuilder(mgr).Complete(r)
}
