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

package training

import (
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/spf13/viper"
	"time"
)

const (
	Namespace = "training.namespace"
	// Enable training API/operator
	Enabled                       = "training.enabled"
	ToolchainIntegrationNamespace = "training.ti_namespace"
	TrainingServiceAccount        = "training.service_account"
	OutputConnectionName          = "training.output_connection"
	ModelBuilderImage             = "training.model_trainer.image"
	NodeSelector                  = "training.node_selector"
	Toleration                    = "training.toleration"
	GPUNodeSelector               = "training.gpu_node_selector"
	GPUToleration                 = "training.gpu_toleration"
	// Kubernetes can consume the GPU resource in the <vendor>.com/gpu format.
	// For example, amd.com/gpu or nvidia.com/gpu.
	ResourceGPUName = "training.gpu_resource_name"
	MetricURL       = "training.metric_url"
	// Timeout for full training process
	Timeout = "training.timeout"
)

const (
	TolerationKey      = "key"
	TolerationOperator = "operator"
	TolerationValue    = "value"
	TolerationEffect   = "effect"
)

func init() {
	viper.SetDefault(Enabled, true)

	viper.SetDefault(Namespace, "odahu-flow-training")
	config.PanicIfError(viper.BindEnv(Namespace))

	viper.SetDefault(TrainingServiceAccount, "odahu-flow-model-trainer")
	viper.SetDefault(ToolchainIntegrationNamespace, "odahu-flow")

	viper.SetDefault(MetricURL, "")

	viper.Set(Timeout, 4*time.Hour)

	viper.Set(ResourceGPUName, "nvidia.com/gpu")
}
