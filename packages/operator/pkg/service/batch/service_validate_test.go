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
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	api_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/batch"
	"github.com/odahu/odahu-flow/packages/operator/pkg/service/batch"
	"github.com/stretchr/testify/require"
	"testing"
)


var testValidateCreateCases = []struct {
	testName string
	spec api_types.InferenceServiceSpec
	expectedErrString string
}{
	{
		testName: "ok",
		spec: api_types.InferenceServiceSpec{
			Image:       "image",
			Command:     []string{"python"},
			ModelRegistry: v1alpha1.ModelSource{
				Remote: &v1alpha1.RemoteModelSource{
					ModelConnection: "gcs-conn",
					ModelPath:       "path",
				},
			},
		},
	},
	{
		testName: "empty image",
		spec: api_types.InferenceServiceSpec{
			Command:     []string{"python"},
			ModelRegistry: v1alpha1.ModelSource{
				Remote: &v1alpha1.RemoteModelSource{
					ModelConnection: "gcs-conn",
					ModelPath:       "path",
				},
			},
		},
		expectedErrString: fmt.Sprintf(batch.EmptySpecFieldErrorMessage, "image"),
	},
	{
		testName: "empty command",
		spec: api_types.InferenceServiceSpec{
			Image:       "image",
			ModelRegistry: v1alpha1.ModelSource{
				Remote: &v1alpha1.RemoteModelSource{
					ModelConnection: "gcs-conn",
					ModelPath:       "path",
				},
			},
		},
		expectedErrString: fmt.Sprintf(batch.EmptySpecFieldErrorMessage, "command"),
	},
	{
		testName: "empty modelSource.connection",
		spec: api_types.InferenceServiceSpec{
			Image:       "image",
			Command:     []string{"python"},
			ModelRegistry: v1alpha1.ModelSource{
				Remote: &v1alpha1.RemoteModelSource{
					ModelPath:       "path",
				},
			},
		},
		expectedErrString: fmt.Sprintf(batch.EmptySpecFieldErrorMessage, "modelRegistry.remote.modelConnection"),

	},
}
func TestValidateCreateUpdate(t *testing.T) {
	for _, test := range testValidateCreateCases {
		t.Run(test.testName, func(t *testing.T) {
			req := require.New(t)
			bis := api_types.InferenceService{
				ID:   "bis",
				Spec: test.spec,
			}
			errs := batch.ValidateCreateUpdate(bis)
			if len(test.expectedErrString) > 0 {
				for _, err := range errs {
					if err.Error() == test.expectedErrString {
						return
					}
				}
				req.FailNow(fmt.Sprintf("Expected error not found: %s", test.expectedErrString))
			} else {
				req.Empty(errs)
			}

		})

	}
}



func TestDefaultCreate(t *testing.T) {
	req := require.New(t)
	bis := api_types.InferenceService{
		ID:           "bis",
		Spec:         api_types.InferenceServiceSpec{},
	}

	batch.DefaultCreate(&bis)
	req.True(bis.Spec.Triggers.Webhook.Enabled)
}