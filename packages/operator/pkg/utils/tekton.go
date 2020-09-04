/*
 * Copyright 2019 EPAM Systems
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import (
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	corev1 "k8s.io/api/core/v1"
	k8s_resource "k8s.io/apimachinery/pkg/api/resource"
)

const (
	tektonContainerPrefix = "step-"
)

var (
	DefaultHelperLimits = corev1.ResourceList{
		corev1.ResourceCPU:    k8s_resource.MustParse("128m"),
		corev1.ResourceMemory: k8s_resource.MustParse("128Mi"),
	}
	EmptyHelperContainerRequestRes = corev1.ResourceList{
		corev1.ResourceMemory: *k8s_resource.NewQuantity(0, k8s_resource.DecimalSI),
		corev1.ResourceCPU:    *k8s_resource.NewQuantity(0, k8s_resource.DecimalSI),
	}
)

// Generate tekton container name base on the step name
func TektonContainerName(stepName string) string {
	return tektonContainerPrefix + stepName
}

// Resources of helper containers are copy of main trainer/packager resources, but
// limit doesn't contain GPU part and all requests res are zeroes.
// If core limit resources is nill then defaultHelperLimits will be used.
func CalculateHelperContainerResources(
	res corev1.ResourceRequirements, gpuResourceName string,
) corev1.ResourceRequirements {
	clippedResources := res.DeepCopy()
	delete(clippedResources.Limits, corev1.ResourceName(gpuResourceName))

	if clippedResources.Limits == nil {
		clippedResources.Limits = DefaultHelperLimits.DeepCopy()
	} else {
		if _, ok := clippedResources.Limits[corev1.ResourceMemory]; !ok {
			clippedResources.Limits[corev1.ResourceMemory] = DefaultHelperLimits[corev1.ResourceMemory].DeepCopy()
		}
		if _, ok := clippedResources.Limits[corev1.ResourceCPU]; !ok {
			clippedResources.Limits[corev1.ResourceCPU] = DefaultHelperLimits[corev1.ResourceCPU].DeepCopy()
		}
	}

	clippedResources.Requests = EmptyHelperContainerRequestRes.DeepCopy()

	return *clippedResources
}

// Build affinity that matches all nodes from nodePools list
func BuildNodeAffinity(nodePools []config.NodePool) *corev1.Affinity {
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

	return &corev1.Affinity{
		NodeAffinity: &corev1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{NodeSelectorTerms: nodeSelectorTerms},
		},
	}
}
