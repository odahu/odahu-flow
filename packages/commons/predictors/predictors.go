/*
 * Copyright 2021 EPAM Systems
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

package predictors

import corev1 "k8s.io/api/core/v1"

type Predictor struct {
	// Predictor ID
	ID string
	// List of ports to expose
	Ports []corev1.ContainerPort
	// Endpoint to check Liveness
	LivenessProbe corev1.Probe
	// Endpoint to check Readiness
	ReadinessProbe corev1.Probe
	// OPA policy filename
	OpaPolicyFilename string
	// Inference endpoint regex
	InferenceEndpointRegex string
}

var (
	OdahuMLServer = Predictor{
		ID:                "odahu",
		OpaPolicyFilename: "odahu_ml_server.rego",
		Ports: []corev1.ContainerPort{{
			Name:          "http1",
			ContainerPort: int32(5000),
			Protocol:      corev1.ProtocolTCP,
		}},
		LivenessProbe: corev1.Probe{
			Handler: corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/healthcheck",
				},
			},
			FailureThreshold: 15,
			PeriodSeconds:    1,
			TimeoutSeconds:   1,
		},
		ReadinessProbe: corev1.Probe{
			Handler: corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/healthcheck",
				},
			},
			FailureThreshold: 15,
			PeriodSeconds:    1,
			TimeoutSeconds:   1,
		},
		InferenceEndpointRegex: ".*/api/model/invoke.*",
	}

	Triton = Predictor{
		ID:                "triton",
		OpaPolicyFilename: "triton.rego",
		Ports: []corev1.ContainerPort{
			{
				Name:          "http1",
				ContainerPort: int32(8000),
				Protocol:      corev1.ProtocolTCP,
			},
			// Currently disabled because of Knative limitations
			// https://github.com/knative/serving/issues/7140
			//{
			//	Name:          "grpc-inference",
			//	ContainerPort: int32(8001),
			//	Protocol:      corev1.ProtocolTCP,
			//},
			//{
			//	Name:          "http-metrics",
			//	ContainerPort: int32(8002),
			//	Protocol:      corev1.ProtocolTCP,
			//},
		},
		LivenessProbe: corev1.Probe{
			Handler: corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/v2/health/live",
				},
			},
			FailureThreshold: 15,
			PeriodSeconds:    1,
			TimeoutSeconds:   1,
		},
		ReadinessProbe: corev1.Probe{
			Handler: corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/v2/health/ready",
				},
			},
			FailureThreshold: 15,
			PeriodSeconds:    1,
			TimeoutSeconds:   1,
		},
		InferenceEndpointRegex: `.*/v2/models/.*/infer/?`,
	}

	Predictors = map[string]Predictor{
		OdahuMLServer.ID: OdahuMLServer,
		Triton.ID:        Triton,
	}
)
