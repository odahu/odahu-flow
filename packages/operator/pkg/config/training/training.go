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
)

const (
	Namespace = "training.namespace"
	// Enable training API/operator
	Enabled                       = "training.enabled"
	ToolchainIntegrationNamespace = "training.ti_namespace"
	TrainingServiceAccount        = "training.service_account"
	OutputConnectionName          = "training.output_connection"
	ModelBuilderImage             = "training.model_trainer.image"
	ModelValidatorImage           = "training.model_validator.image"
	NodeSelector                  = "training.node_selector"
	Toleration                    = "training.toleration"
	MetricURL                     = "training.metric_url"
	Timeout                       = "packaging.timeout"
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
}
