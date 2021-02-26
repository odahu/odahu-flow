/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BatchJobState defines current state
type BatchJobState string

const (
	BatchScheduling BatchJobState = "scheduling"
	BatchRunning    BatchJobState = "running"
	BatchSucceeded  BatchJobState = "succeeded"
	BatchFailed     BatchJobState = "failed"
	BatchUnknown    BatchJobState = "unknown"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BatchInferenceJobSpec defines the desired state of BatchInferenceJob
type BatchInferenceJobSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	BatchInferenceServiceID string `json:"batchInferenceService"`
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
	// InputPath is a source directory for BatchInferenceJob input data
	// relative to bucket root
	InputPath string `json:"inputPath"`
	// OutputPath is a destination directory for BatchInferenceJob results
	// relative to bucket root
	OutputPath string `json:"outputPath"`
	// Node selector for specifying a node pool
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// Resources for model container
	// The same format like k8s uses for pod resources.
	Resources *ResourceRequirements `json:"resources,omitempty"`
	// BatchRequestID is unique identifier for BatchInferenceJob that helps to correlate between
	// Model input, model output and feedback
	// Take into account that it is not the same as kubeflow InferenceRequest id
	// Each BatchInferenceJob can process more than one InferenceRequest (delivered in separate input file)
	// So each BatchRequestID has set of corresponding InferenceRequest and their IDs
	BatchRequestID string `json:"requestId"`
}

// BatchInferenceJobStatus defines the observed state of BatchInferenceJob
type BatchInferenceJobStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	State BatchJobState `json:"state"`
	Message string `json:"message"`
	Reason string `json:"reason"`
}

// +kubebuilder:object:root=true

// BatchInferenceJob is the Schema for the batchinferencejobs API
type BatchInferenceJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BatchInferenceJobSpec   `json:"spec,omitempty"`
	Status BatchInferenceJobStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BatchInferenceJobList contains a list of BatchInferenceJob
type BatchInferenceJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BatchInferenceJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BatchInferenceJob{}, &BatchInferenceJobList{})
}
