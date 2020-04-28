//
//    Copyright 2019 EPAM Systems
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
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestPodMutator(t *testing.T) {
	intVal := 42
	int64Val := int64(intVal)
	nodeSelector := map[string]string{"label": "labelValue"}
	toleration := &corev1.Toleration{
		Key:               "key",
		Operator:          corev1.TolerationOpExists,
		Value:             "value",
		Effect:            corev1.TaintEffectNoSchedule,
		TolerationSeconds: &int64Val}
	defaultToleration := corev1.Toleration{
		Key:               "defaultKey",
		Operator:          corev1.TolerationOpEqual,
		Value:             "defaultValue",
		Effect:            corev1.TaintEffectPreferNoSchedule,
		TolerationSeconds: &int64Val}

	pm := podMutator{deploymentConfig: config.ModelDeploymentConfig{NodeSelector: nodeSelector, Toleration: toleration}}
	pod := &corev1.Pod{}
	pod.Spec.Tolerations = append(pod.Spec.Tolerations, defaultToleration)
	_ = pm.addNodeSelectors(pod)

	expectedPod := &corev1.Pod{}
	expectedPod.Spec.Tolerations = append(expectedPod.Spec.Tolerations, defaultToleration, *toleration)
	expectedPod.Spec.NodeSelector = nodeSelector

	g := NewGomegaWithT(t)
	g.Expect(pod).Should(BeEquivalentTo(expectedPod))
}
