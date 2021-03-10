/*
 *
 *     Copyright 2021 EPAM Systems
 *
 *     Licensed under the Apache License, Version 2.0 (the "License");
 *     you may not use this file except in compliance with the License.
 *     You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *     Unless required by applicable law or agreed to in writing, software
 *     distributed under the License is distributed on an "AS IS" BASIS,
 *     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *     See the License for the specific language governing permissions and
 *     limitations under the License.
 */

package batch_test

import (
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	api_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/batch"
	"github.com/odahu/odahu-flow/packages/operator/pkg/service/batch"
	"github.com/stretchr/testify/require"
	"testing"
)

var resString = "resource"
var resDefInJobString = "resource-defined-in-job"

var testDefaultJobCases = []struct {
	testName string
	job      api_types.InferenceJob
	service  api_types.InferenceService
	expected api_types.InferenceJob
}{
	{
		testName: "empty job, filled service",
		job: api_types.InferenceJob{
			Spec: api_types.InferenceJobSpec{
				InferenceServiceID: "service",
				BatchRequestID:     "unique-req-id",
			},
		},
		service: api_types.InferenceService{
			Spec: api_types.InferenceServiceSpec{
				Image:   "image",
				Command: []string{"python"},
				Args:    []string{"--verbose"},
				ModelSource: api_types.ConnectionReference{
					Connection: "conn",
					Path:       "path",
				},
				DataSource: &api_types.ConnectionReference{
					Connection: "conn",
					Path:       "path",
				},
				OutputDestination: &api_types.ConnectionReference{
					Connection: "conn",
					Path:       "path",
				},
				Triggers:     api_types.InferenceServiceTriggers{},
				NodeSelector: map[string]string{"node": "gpu"},
				Resources: &v1alpha1.ResourceRequirements{
					Limits: &v1alpha1.ResourceList{
						GPU:    &resString,
						CPU:    &resString,
						Memory: &resString,
					},
					Requests: &v1alpha1.ResourceList{
						GPU:    &resString,
						CPU:    &resString,
						Memory: &resString,
					},
				},
			},
		},
		expected: api_types.InferenceJob{
			Spec: api_types.InferenceJobSpec{
				InferenceServiceID: "service",
				BatchRequestID:     "unique-req-id",
				DataSource: &api_types.ConnectionReference{
					Connection: "conn",
					Path:       "path",
				},
				OutputDestination: &api_types.ConnectionReference{
					Connection: "conn",
					Path:       "path",
				},
				NodeSelector: map[string]string{"node": "gpu"},
				Resources: &v1alpha1.ResourceRequirements{
					Limits: &v1alpha1.ResourceList{
						GPU:    &resString,
						CPU:    &resString,
						Memory: &resString,
					},
					Requests: &v1alpha1.ResourceList{
						GPU:    &resString,
						CPU:    &resString,
						Memory: &resString,
					},
				},
			},
		},
	},
	{
		testName: "filled job values have a priority",
		job: api_types.InferenceJob{Spec: api_types.InferenceJobSpec{
			InferenceServiceID: "service",
			BatchRequestID:     "unique-req-id",
			DataSource: &api_types.ConnectionReference{
				Connection: "conn-defined-in-job",
				Path:       "path-defined-in-job",
			},
			OutputDestination: &api_types.ConnectionReference{
				Connection: "conn-defined-in-job",
				Path:       "path-defined-in-job",
			},
			NodeSelector: map[string]string{"node": "gpu-defined-in-job"},
			Resources: &v1alpha1.ResourceRequirements{
				Limits: &v1alpha1.ResourceList{
					GPU:    &resDefInJobString,
					CPU:    &resDefInJobString,
					Memory: &resDefInJobString,
				},
				Requests: &v1alpha1.ResourceList{
					GPU:    &resDefInJobString,
					CPU:    &resDefInJobString,
					Memory: &resDefInJobString,
				},
			},
		}},
		service: api_types.InferenceService{
			Spec: api_types.InferenceServiceSpec{
				Image:   "image",
				Command: []string{"python"},
				Args:    []string{"--verbose"},
				ModelSource: api_types.ConnectionReference{
					Connection: "conn",
					Path:       "path",
				},
				DataSource: &api_types.ConnectionReference{
					Connection: "conn",
					Path:       "path",
				},
				OutputDestination: &api_types.ConnectionReference{
					Connection: "conn",
					Path:       "path",
				},
				Triggers:     api_types.InferenceServiceTriggers{},
				NodeSelector: map[string]string{"node": "gpu"},
				Resources: &v1alpha1.ResourceRequirements{
					Limits: &v1alpha1.ResourceList{
						GPU:    &resString,
						CPU:    &resString,
						Memory: &resString,
					},
					Requests: &v1alpha1.ResourceList{
						GPU:    &resString,
						CPU:    &resString,
						Memory: &resString,
					},
				},
			},
		},
		expected: api_types.InferenceJob{Spec: api_types.InferenceJobSpec{
			DataSource: &api_types.ConnectionReference{
				Connection: "conn-defined-in-job",
				Path:       "path-defined-in-job",
			},
			OutputDestination: &api_types.ConnectionReference{
				Connection: "conn-defined-in-job",
				Path:       "path-defined-in-job",
			},
			NodeSelector: map[string]string{"node": "gpu-defined-in-job"},
			Resources: &v1alpha1.ResourceRequirements{
				Limits: &v1alpha1.ResourceList{
					GPU:    &resDefInJobString,
					CPU:    &resDefInJobString,
					Memory: &resDefInJobString,
				},
				Requests: &v1alpha1.ResourceList{
					GPU:    &resDefInJobString,
					CPU:    &resDefInJobString,
					Memory: &resDefInJobString,
				},
			},
			InferenceServiceID: "service",
			BatchRequestID:     "unique-req-id",
		}},
	},
}

func TestDefaultJob(t *testing.T) {
	for _, test := range testDefaultJobCases {
		t.Run(test.testName, func(t *testing.T) {
			req := require.New(t)
			batch.DefaultJob(&test.job, test.service)
			req.Equal(test.job.Spec, test.expected.Spec)
		})
	}
}
