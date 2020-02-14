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

package packaging

import (
	"time"

	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/spf13/viper"
)

const (
	Namespace = "packaging.namespace"
	// Enable packaging API/operator
	Enabled                       = "packaging.enabled"
	PackagingIntegrationNamespace = "packaging.packager_integration_namespace"
	ServiceAccount                = "packaging.service_account"
	OutputConnectionName          = "packaging.output_connection"
	ModelPackagerImage            = "packaging.model_packager.image"
	NodeSelector                  = "packaging.node_selector"
	Toleration                    = "packaging.toleration"
	// Timeout for full packaging process
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

	viper.SetDefault(Namespace, "odahu-flow-packaging")
	viper.SetDefault(PackagingIntegrationNamespace, "odahu-flow")
	config.PanicIfError(viper.BindEnv(Namespace))

	viper.SetDefault(ServiceAccount, "odahu-flow-model-packager")

	viper.SetDefault(Timeout, 4 * time.Hour)
}
