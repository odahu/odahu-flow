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

package config

type Vault struct {
	// Vault URL
	URL string `json:"url"`
	// Vault secret engine path where connection will be stored
	SecretEnginePath string `json:"secretEnginePath"`
	// Vault role for access to the secret engine path
	Role string `json:"role"`
	// Optionally. Token for access to the vault server
	// If it is empty then client will use the k8s auth
	Token string `json:"token"`
}

type ConnectionConfig struct {
	// Enable connection API/operator
	Namespace string `json:"namespace"`
	// Connection API server and operator are enabled
	Enabled bool `json:"enabled"`
	// Storage backend for connections. Available options:
	//   * kubernetes
	//   * vault
	RepositoryType RepositoryType `json:"repositoryType"`
	// Connection Vault configuration
	Vault Vault `json:"vault"`
}

func NewDefaultConnectionConfig() ConnectionConfig {
	return ConnectionConfig{
		Namespace:      "odahu-flow",
		Enabled:        true,
		RepositoryType: RepositoryKubernetesType,
		Vault: Vault{
			URL:              "http://127.0.0.1:8200",
			SecretEnginePath: "odahu-flow/connections",
			Role:             "odahu-flow",
			Token:            "",
		},
	}
}
