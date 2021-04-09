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

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.


type RemoteModelSource struct {
	// ModelConnection is name of connection to object storage bucket where ML model files are expected
	ModelConnection string `json:"modelConnection"`
	// ModelPath is a directory inside ModelConnection where ML model files are located
	ModelPath string `json:"modelPath"`
}

type ModelMeta struct {
	Name string `json:"name"`
	Version string `json:"version"`
}

type LocalModelSource struct {
	ModelMeta ModelMeta `json:"meta"`
	// ModelPath is a directory inside container where ML model files are located
	ModelPath string `json:"modelPath"`
}

type ModelSource struct {
	// Remote fetch model from remote model registry using ODAHU connections mechanism
	Remote *RemoteModelSource `json:"remote,omitempty"`
	// Local does not fetch model and assume that model is embedded into container
	Local *LocalModelSource `json:"local,omitempty"`
}

// BatchInferenceJobSpec defines the desired state of BatchInferenceJob
type BatchInferenceJobSpec struct {
	// Docker image
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
	// InputConnection is name of connection to object storage bucket where input data are expected
	InputConnection string `json:"inputConnection"`
	// InputPath is a source directory for BatchInferenceJob input data
	// relative to bucket root of InputConnection
	InputPath string `json:"inputPath"`
	// OutputConnection is name of connection to object storage bucket where results should be
	// saved
	OutputConnection string `json:"outputConnection"`
	// OutputPath is a destination directory for BatchInferenceJob results
	// relative to bucket root of OutputConnection
	OutputPath string `json:"outputPath"`
	ModelSource ModelSource `json:"modelSource"`
	// Node selector for specifying a node pool
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// Resources for model container
	// The same format like k8s uses for pod resources.
	Resources *ResourceRequirements `json:"resources,omitempty"`
	// requestId is unique identifier for InferenceJob that helps to correlate between
	// Model input, model output and feedback.
	// Take into account that it is not the same as kubeflow InferenceRequest id.
	// Each BatchInferenceJob can process more than one InferenceRequest (delivered in separate input file).
	// So each RequestID has set of corresponding InferenceRequest and their IDs.
	BatchRequestID string `json:"requestId"`
}

// BatchInferenceJobStatus defines the observed state of BatchInferenceJob
type BatchInferenceJobStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	State BatchJobState `json:"state"`
	Message string `json:"message"`
	Reason string `json:"reason"`
	PodName string `json:"podName"`
}

// +kubebuilder:object:root=true

// BatchInferenceJob is the Schema for the batchinferencejobs API
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.state"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.message"
// +kubebuilder:printcolumn:name="Reason",type="string",JSONPath=".status.reason"
// +kubebuilder:resource:shortName=bij
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
