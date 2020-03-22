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

type JWKS struct {
	// JWKS URL
	URL     string `json:"url"`
	// Issuer claim value
	Issuer  string `json:"issuer"`
	// Model authorization enabled
	Enabled bool   `json:"enabled"`
}

type EdgeConfig struct {
	// External model host
	Host string `json:"host"`
}

type ModelDeploymentIstioConfig struct {
	// Istio ingress gateway service name
	ServiceName string `json:"serviceName"`
	// Istio ingress gateway namespace
	Namespace   string `json:"namespace"`
}

type ModelDeploymentSecurityConfig struct {
	JWKS     JWKS   `json:"jwks"`
	// Deprecated
	RoleName string `json:"roleName"`
}

type ModelDeploymentConfig struct {
	// Kubernetes namespace, where model deployments will be deployed
	Namespace string `json:"namespace"`
	// Enable deployment API/operator
	Enabled                   bool                          `json:"enabled"`
	Security                  ModelDeploymentSecurityConfig `json:"security"`
	// Default connection ID which will be used if a user doesn't specify it in a model deployment
	DefaultDockerPullConnName string                        `json:"defaultDockerPullConnName"`
	Edge                      EdgeConfig                    `json:"edge"`
	// Kubernetes node selector for model deployments
	NodeSelector              map[string]string             `json:"nodeSelector"`
	// Kubernetes tolerations for model deployments
	Toleration                map[string]string             `json:"toleration"`
	Istio                     ModelDeploymentIstioConfig    `json:"istio"`
}

func NewDefaultModelDeploymentConfig() ModelDeploymentConfig {
	return ModelDeploymentConfig{
		Namespace: "odahu-flow-deployment",
		Enabled:   true,
		Security: ModelDeploymentSecurityConfig{
			RoleName: "default-odahu-flow",
		},
		DefaultDockerPullConnName: "",
		Istio: ModelDeploymentIstioConfig{
			ServiceName: "istio-ingressgateway",
			Namespace:   "istio-system",
		},
	}
}
