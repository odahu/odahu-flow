package deploymenthook

import (
	"context"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/prometheus/common/log"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apitypes "k8s.io/apimachinery/pkg/types"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/builder"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

const (
	webhookName        = "modeldeployment-webhook"
	webhookServerName  = "modeldeployment-webhook-server"
	webhookServiceName = "modeldeployment-webhook-service"
	webhookSecretName  = "modeldeployment-webhook-secret"
	webhookconfigName  = "modeldeployment-webhook-config"
)

func Add(
	mgr manager.Manager,
	deploymentConfig config.ModelDeploymentConfig,
	_ config.OperatorConfig,
	_ string,
) error {
	log := logf.Log.WithName(webhookName)
	log.Info("Creating model deployment webhook for knative pods")

	wh, err := builder.NewWebhookBuilder().
		Name(webhookName).
		Mutating().
		Operations(admissionregistrationv1beta1.Create, admissionregistrationv1beta1.Update).
		WithManager(mgr).
		ForType(&corev1.Pod{}).
		NamespaceSelector(&metav1.LabelSelector{MatchLabels: map[string]string{"namespace": deploymentConfig.Namespace}}).
		Handlers(&podMutator{deploymentConfig: deploymentConfig}).
		Build()
	if err != nil {
		return err
	}

	log.Info("Setting up deployment webhook server")
	as, err := webhook.NewServer(webhookServerName, mgr, webhook.ServerOptions{
		Port: 6443,
		BootstrapOptions: &webhook.BootstrapOptions{
			MutatingWebhookConfigName: webhookconfigName,
			Secret: &apitypes.NamespacedName{
				Namespace: deploymentConfig.Namespace,
				Name:      webhookSecretName,
			},
			Service: &webhook.Service{
				Namespace: deploymentConfig.Namespace,
				Name:      webhookServiceName,
				Selectors: deploymentConfig.NodeSelector,
			},
		},
	})
	if err != nil {
		return err
	}

	log.Info("Registering deployment webhook to the server")
	err = as.Register(wh)
	if err != nil {
		return err
	}
	return nil
}

type podMutator struct {
	client           client.Client
	decoder          types.Decoder
	deploymentConfig config.ModelDeploymentConfig
}

// Implement admission.Handler so the controller can handle admission request.
var _ admission.Handler = &podMutator{}

// podMutator adds an annotation to every incoming pods.
func (pm *podMutator) Handle(_ context.Context, req types.Request) types.Response {
	pod := &corev1.Pod{}

	err := pm.decoder.Decode(req, pod)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	podCopy := pod.DeepCopy()
	err = pm.addNodeSelectors(podCopy)
	if err != nil {
		return admission.ErrorResponse(http.StatusInternalServerError, err)
	}
	return admission.PatchResponse(pod, podCopy)
}

//Adds node selectors from deployment config to knative pods
func (pm *podMutator) addNodeSelectors(pod *corev1.Pod) error {
	nodeSelector := pm.deploymentConfig.NodeSelector
	if len(nodeSelector) > 0 {
		pod.Spec.NodeSelector = nodeSelector
		log.Infof("Assigning node selector %v to a pod %v", nodeSelector, pod.Name)
	} else {
		log.Warnf("Got empty node selector from deployment config, doing nothing to a pod %v", pod.Name)
	}
	return nil
}

//Client and Decoder are auto injected
var _ inject.Client = &podMutator{}

func (pm *podMutator) InjectClient(c client.Client) error {
	pm.client = c
	return nil
}

var _ inject.Decoder = &podMutator{}

func (pm *podMutator) InjectDecoder(d types.Decoder) error {
	pm.decoder = d
	return nil
}
