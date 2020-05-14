//
//    Copyright 2020 EPAM Systems
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package deploymenthook

import (
	"context"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" //nolint
	"net/http"
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
)

var log = logf.Log.WithName(webhookName)

func Add(
	mgr manager.Manager,
	deploymentConfig config.ModelDeploymentConfig,
	operatorConfig config.OperatorConfig,
	_ string,
) error {
	log.Info("Creating model deployment webhook for knative pods")

	wh, err := builder.NewWebhookBuilder().
		Mutating().
		Operations(admissionregistrationv1beta1.Create, admissionregistrationv1beta1.Update).
		WithManager(mgr).
		ForType(&corev1.Pod{}).
		NamespaceSelector(&metav1.LabelSelector{
			MatchLabels: map[string]string{namespaceSelectorKey: namespaceSelectorValue}}).
		Handlers(&podMutator{deploymentConfig: deploymentConfig}).
		Build()
	if err != nil {
		return err
	}

	log.Info("Setting up deployment webhook server", "service namespace", operatorConfig.Namespace)
	as, err := webhook.NewServer(webhookServerName, mgr, webhook.ServerOptions{
		Port: 6443,
		BootstrapOptions: &webhook.BootstrapOptions{
			MutatingWebhookConfigName: webhookconfigName,
			Service: &webhook.Service{
				Namespace: operatorConfig.Namespace,
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
	decoder          types.Decoder
	deploymentConfig config.ModelDeploymentConfig
}

var _ admission.Handler = &podMutator{}

func (pm *podMutator) Handle(_ context.Context, req types.Request) types.Response {
	pod := &corev1.Pod{}

	err := pm.decoder.Decode(req, pod)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	podCopy := pod.DeepCopy()
	pm.addNodeSelectors(podCopy)
	return admission.PatchResponse(pod, podCopy)
}

//Adds node selectors and tolerations from the deployment config to knative pods
func (pm *podMutator) addNodeSelectors(pod *corev1.Pod)  {
	nodeSelector := pm.deploymentConfig.NodeSelector
	if len(nodeSelector) > 0 {
		pod.Spec.NodeSelector = nodeSelector
		log.Info("Assigning node selector to a pod", "nodeSelector", nodeSelector, "pod name", pod.Name)
	} else {
		log.Info("Got empty node selector from deployment config, skipping", "pod name", pod.Name)
	}

	toleration := pm.deploymentConfig.Toleration
	if toleration != nil {
		pod.Spec.Tolerations = append(pod.Spec.Tolerations, *toleration)
		log.Info("Assigning toleration to a pod", "toleration", toleration, "pod name", pod.Name)
	} else {
		log.Info("Got empty toleration from deployment config, skipping", "pod name", pod.Name)
	}
}

var _ inject.Decoder = &podMutator{}

func (pm *podMutator) InjectDecoder(d types.Decoder) error {
	pm.decoder = d
	return nil
}
