package deploymenthook

import (
	"context"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	webhookconfigName  = "modeldeployment-webhook-config"

	//label for a watched  namespace where pods are created
	namespaceSelectorKey   = "modeldeployment-webhook"
	namespaceSelectorValue = "enabled"

	//selectors for pods for webhook to run on
	webhookServiceSelectorKey   = "app"
	webhookServiceSelectorValue = "operator"

	//namespace where webhook service will be created
	webhookServiceNamespace = "odahu-flow"
)

var log = logf.Log.WithName(webhookName)

func Add(
	mgr manager.Manager,
	deploymentConfig config.ModelDeploymentConfig,
	_ config.OperatorConfig,
	_ string,
) error {
	log.Info("Creating model deployment webhook for knative pods")

	wh, err := builder.NewWebhookBuilder().
		Mutating().
		Operations(admissionregistrationv1beta1.Create, admissionregistrationv1beta1.Update).
		WithManager(mgr).
		ForType(&corev1.Pod{}).
		NamespaceSelector(&metav1.LabelSelector{MatchLabels: map[string]string{namespaceSelectorKey: namespaceSelectorValue}}).
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
			Service: &webhook.Service{
				Namespace: webhookServiceNamespace,
				Name:      webhookServiceName,
				Selectors: map[string]string{webhookServiceSelectorKey: webhookServiceSelectorValue},
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

var _ admission.Handler = &podMutator{}

// +kubebuilder:rbac:groups=admissionregistration.k8s.io,resources=mutatingwebhookconfigurations,verbs=get;list;watch;create;update;patch;delete
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

//Adds node selectors and tolerations from the deployment config to knative pods
func (pm *podMutator) addNodeSelectors(pod *corev1.Pod) error {
	nodeSelector := pm.deploymentConfig.NodeSelector
	if len(nodeSelector) > 0 {
		pod.Spec.NodeSelector = nodeSelector
		log.Info("Assigning node selector to a pod", "nodeSelector", nodeSelector, "pod name", pod.Name)
	} else {
		log.Info("Got empty node selector from deployment config, skipping", "pod name", pod.Name)
	}

	toleration := pm.deploymentConfig.Toleration
	if toleration != nil {
		parsedToleration := corev1.Toleration{Key: toleration.Key,
			Operator:          corev1.TolerationOperator(toleration.Operator),
			Value:             toleration.Value,
			Effect:            corev1.TaintEffect(toleration.Effect),
			TolerationSeconds: toleration.TolerationSeconds}
		pod.Spec.Tolerations = append(pod.Spec.Tolerations, parsedToleration)
		log.Info("Assigning toleration to a pod", "toleration", toleration, "pod name", pod.Name)
	} else {
		log.Info("Got empty toleration from deployment config, skipping", "pod name", pod.Name)
	}

	return nil
}

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
