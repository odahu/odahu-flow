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

package config

import (
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"os"
	"time"
)

var (
	defaultTrainingMemoryLimit    = "256Mi"
	defaultTrainingCPULimit       = "250m"
	defaultTrainingMemoryRequests = "128Mi"
	defaultTrainingCPURequests    = "125m"
)

type ModelTrainingConfig struct {
	// Kubernetes namespace, where model trainings will be deployed
	Namespace                     string `json:"namespace"`
	ToolchainIntegrationNamespace string `json:"toolchainIntegrationNamespace"`
	// Enable deployment API/operator
	Enabled            bool   `json:"enabled"`
	ServiceAccount     string `json:"serviceAccount"`
	OutputConnectionID string `json:"outputConnectionID"`

	// Node pools to run training tasks on
	NodePools []NodePool `json:"nodePools"`
	// Kubernetes tolerations for model trainings pods
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// Node pools to run GPU training tasks on
	GPUNodePools []NodePool `json:"gpuNodePools"`
	// Kubernetes tolerations for GPU model trainings pods
	GPUTolerations []corev1.Toleration `json:"gpuTolerations,omitempty"`

	MetricURL         string `json:"metricUrl"`
	ModelTrainerImage string `json:"modelTrainerImage"`
	// Timeout for full training process
	Timeout time.Duration `json:"timeout"`
	// Default resources for training pods
	DefaultResources odahuflowv1alpha1.ResourceRequirements `json:"defaultResources"`

	// Storage backend for toolchain integrations. Available options:
	//   * kubernetes
	//   * postgres
	ToolchainIntegrationRepositoryType RepositoryType `json:"toolchainIntegrationRepositoryType"`
}

func NewDefaultModelTrainingConfig() ModelTrainingConfig {
	return ModelTrainingConfig{
		Namespace:                     "odahu-flow-training",
		ToolchainIntegrationNamespace: "odahu-flow",
		Enabled:                       true,
		Timeout:                       4 * time.Hour,
		ServiceAccount:                "odahu-flow-model-trainer",
		// workaround https://github.com/spf13/viper/issues/761
		ModelTrainerImage: os.Getenv("TRAINING_MODEL_TRAINER_IMAGE"),
		DefaultResources: odahuflowv1alpha1.ResourceRequirements{
			Requests: &odahuflowv1alpha1.ResourceList{
				CPU:    &defaultTrainingCPURequests,
				Memory: &defaultTrainingMemoryRequests,
			},
			Limits: &odahuflowv1alpha1.ResourceList{
				CPU:    &defaultTrainingCPULimit,
				Memory: &defaultTrainingMemoryLimit,
			},
		},
		ToolchainIntegrationRepositoryType: RepositoryPostgresType,
	}
}

type NodePool struct {
	NodeSelector map[string]string `json:"nodeSelector"`
	Tags         []string          `json:"tags"`
}
