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

package batch

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"time"
)

// ConnectionReference refers to specific Connection. Connection path can be overridden using Path
type ConnectionReference struct {
	// Next connection types are supported
	Connection string `json:"connection"`
	// User can override path otherwise Connection path will be used
	Path string `json:"path"`
}

type PredictorWebhookTrigger struct {
	// Enabled. If True then it possible to run InferenceJob by creating it using REST API
	Enabled bool `json:"enabled"`
}

type InferenceServiceTriggers struct {
	// Webhook provides a REST API to execute InferenceJob that correspond to this service
	Webhook *PredictorWebhookTrigger `json:"webhook,omitempty"`
}

type InferenceServiceSpec struct {
	// Image is OCI image that contains user defined prediction code
	Image string `json:"image"`
	// Entrypoint array. Not executed within a shell.
	// The docker image's ENTRYPOINT is used if this is not provided.
	// Variable references $(VAR_NAME) are expanded using the container's environment. If a variable
	// cannot be resolved, the reference in the input string will be unchanged. The $(VAR_NAME) syntax
	// can be escaped with a double $$, ie: $$(VAR_NAME). Escaped references will never be expanded,
	// regardless of whether the variable exists or not.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell
	Command []string `json:"command"`
	// Arguments to the entrypoint.
	// The docker image's CMD is used if this is not provided.
	// Variable references $(VAR_NAME) are expanded using the container's environment. If a variable
	// cannot be resolved, the reference in the input string will be unchanged. The $(VAR_NAME) syntax
	// can be escaped with a double $$, ie: $$(VAR_NAME). Escaped references will never be expanded,
	// regardless of whether the variable exists or not.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell
	Args []string `json:"args"`
	// ModelSource defines location of ML model files
	ModelSource ConnectionReference `json:"modelSource"`
	// DataSource defines location input data files
	// Can be overridden in BatchInferenceJob definition
	DataSource *ConnectionReference `json:"dataSource,omitempty"`
	// OutputDestination defines location of directory with output files
	// Can be overridden in BatchInferenceJob definition
	OutputDestination *ConnectionReference `json:"outputDestination,omitempty"`
	// Triggers are describe how to run InferenceService
	Triggers InferenceServiceTriggers `json:"triggers"`
	// Node selector for specifying a node pool
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// Resources for model container
	// The same format like k8s uses for pod resources.
	Resources *v1alpha1.ResourceRequirements `json:"resources,omitempty"`
}

type InferenceServiceStatus struct {}

type InferenceService struct {
	ID string `json:"id"`
	// Deletion mark. Managed by system. Cannot be overridden by User
	DeletionMark bool `json:"deletionMark,omitempty" swaggerignore:"true"`
	// When resource was created. Managed by system. Cannot be overridden by User
	CreatedAt time.Time `json:"createdAt"`
	// When resource was updated. Managed by system. Cannot be overridden by User
	UpdatedAt time.Time              `json:"updatedAt"`
	Spec      InferenceServiceSpec   `json:"spec"`
	Status    InferenceServiceStatus `json:"status"`
}


func (spec InferenceServiceSpec) Value() (driver.Value, error) {
	return json.Marshal(spec)
}

func (spec *InferenceServiceSpec) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	res := json.Unmarshal(b, &spec)
	return res
}

func (in InferenceServiceStatus) Value() (driver.Value, error) {
	return json.Marshal(in)
}

func (in *InferenceServiceStatus) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	res := json.Unmarshal(b, &in)
	return res
}


const TagKey = "name"
type InferenceServiceFilter struct {

}