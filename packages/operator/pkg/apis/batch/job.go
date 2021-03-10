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

type JobState string

const (
	Scheduling JobState = "scheduling"
	Running    JobState = "running"
	Succeeded  JobState = "succeeded"
	Failed     JobState = "failed"
	Unknown    JobState = "unknown"
)

type InferenceJobSpec struct {
	// InferenceServiceID refers to BatchInferenceService
	InferenceServiceID string `json:"inferenceServiceId"`
	// DataSource defines location input data files
	// If nil then will be filled from BatchInferenceService
	DataSource *ConnectionReference `json:"dataSource"`
	// OutputDestination defines location of directory with output files
	// If nil then will be filled from BatchInferenceService
	OutputDestination *ConnectionReference `json:"outputDestination"`
	// Node selector for specifying a node pool
	NodeSelector map[string]string `json:"nodeSelector"`
	// Resources for model container
	// The same format like k8s uses for pod resources.
	Resources *v1alpha1.ResourceRequirements `json:"resources"`
	// BatchRequestID is unique identifier for InferenceJob that helps to correlate between
	// Model input, model output and feedback
	// Take into account that it is not the same as kubeflow InferenceRequest id
	// Each InferenceJob can process more than one InferenceRequest (delivered in separate input file)
	// So each BatchRequestID has set of corresponding InferenceRequest and their IDs
	BatchRequestID string `json:"requestId"`
}

type InferenceJobStatus struct {
	// State describes current state of InferenceJob
	State JobState `json:"state"`
	// Message is any message from runtime service about status of InferenceJob
	Message string `json:"message"`
	// Reason is a reason of some InferenceJob state that was retrieved from runtime service.
	// for example reason of failure
	Reason string `json:"reason"`
	// PodName is a name of Pod in Kubernetes that is running under the hood of InferenceJob
	PodName string `json:"podName"`
}

type InferenceJob struct {
	// Resource ID
	ID string
	// Deletion mark
	DeletionMark bool `json:"deletionMark,omitempty" swaggerignore:"true"`
	// CreatedAt describes when InferenceJob was launched
	CreatedAt time.Time `json:"createdAt,omitempty"`
	// CreatedAt describes when InferenceJob was updated (status was changed)
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	// Spec describes parameters of InferenceJob
	Spec InferenceJobSpec `json:"spec,omitempty"`
	// Spec describes execution status of InferenceJob
	Status InferenceJobStatus `json:"status,omitempty"`
}


func (spec InferenceJobSpec) Value() (driver.Value, error) {
	return json.Marshal(spec)
}

func (spec *InferenceJobSpec) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	res := json.Unmarshal(b, &spec)
	return res
}

func (in InferenceJobStatus) Value() (driver.Value, error) {
	return json.Marshal(in)
}

func (in *InferenceJobStatus) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	res := json.Unmarshal(b, &in)
	return res
}


const JobTagKey = "name"
type InferenceJobFilter struct {

}