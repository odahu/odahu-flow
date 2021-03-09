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
	"database/sql/driver"
	"encoding/json"
	"errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DataBindingDir struct {
	// Connection name for data
	Connection string `json:"connName"`
	// Local training path
	LocalPath string `json:"localPath"`
	// Overwrite remote data path in connection
	RemotePath string `json:"remotePath,omitempty"`
}

type ModelIdentity struct {
	// Model name
	Name string `json:"name"`
	// Model version
	Version string `json:"version"`
	// Template of output artifact name
	ArtifactNameTemplate string `json:"artifactNameTemplate,omitempty"`
}

// ModelTrainingSpec defines the desired state of ModelTraining
type ModelTrainingSpec struct {
	// Model Identity
	Model ModelIdentity `json:"model"`
	// IntegrationName of toolchain
	Toolchain string `json:"toolchain"`
	// Custom environment variables that should be set before entrypoint invocation.
	CustomEnvs []EnvironmentVariable `json:"envs,omitempty"`
	// Model training hyperParameters in parameter:value format
	HyperParameters map[string]string `json:"hyperParameters,omitempty"`
	// Directory with model scripts/files in a git repository
	WorkDir string `json:"workDir,omitempty"`
	// Model training file. It can be python\bash script or jupiter notebook
	Entrypoint          string   `json:"entrypoint"`
	EntrypointArguments []string `json:"args,omitempty"`
	// AlgorithmSource for training
	AlgorithmSource AlgorithmSource `json:"algorithmSource"`
	// Name of Connection to storage where training output artifact will be stored.
	// Permitted connection types are defined by specific toolchain
	OutputConnection string `json:"outputConnection,omitempty"`
	// Train image
	Image string `json:"image,omitempty"`
	// Resources for model container
	// The same format like k8s uses for pod resources.
	Resources *ResourceRequirements `json:"resources,omitempty"`
	// Input data for a training
	Data []DataBindingDir `json:"data,omitempty"`
	// Node selector for specifying a node pool
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
}

// The function returns true if one of the GPU resources is set up.
func (spec *ModelTrainingSpec) IsGPUResourceSet() bool {
	isPresent := func(s *string) bool { return s != nil && *s != "" }
	return spec.Resources != nil && ((spec.Resources.Limits != nil && isPresent(spec.Resources.Limits.GPU)) ||
		(spec.Resources.Requests != nil && isPresent(spec.Resources.Requests.GPU)))
}

// ModelTrainingState defines current state
type ModelTrainingState string

// These are the valid statuses of pods.
const (
	ModelTrainingScheduling ModelTrainingState = "scheduling"
	ModelTrainingRunning    ModelTrainingState = "running"
	ModelTrainingSucceeded  ModelTrainingState = "succeeded"
	ModelTrainingFailed     ModelTrainingState = "failed"
	ModelTrainingUnknown    ModelTrainingState = "unknown"
)

type TrainingResult struct {
	// Mlflow run ID
	RunID string `json:"runId"`
	// Trained artifact name
	ArtifactName string `json:"artifactName"`
	// VCS commit
	CommitID string `json:"commitID"`
}

// ModelTrainingStatus defines the observed state of ModelTraining
type ModelTrainingStatus struct {
	// Pod package for name
	PodName string `json:"podName,omitempty"`
	// Model Packaging State
	State ModelTrainingState `json:"state,omitempty"`
	// Pod exit code
	ExitCode *int32 `json:"exitCode,omitempty"`
	// Pod reason
	Reason *string `json:"reason,omitempty"`
	// Pod last log
	Message *string `json:"message,omitempty"`
	// List of training results
	Artifacts []TrainingResult `json:"artifacts,omitempty"`
	// DEPRECATED Info about create and update
	//CreatedAt *metav1.Time `json:"createdAt,omitempty"`
	//UpdatedAt *metav1.Time `json:"updatedAt,omitempty"`
	Modifiable `json:",inline"`
}

func (spec ModelTrainingSpec) Value() (driver.Value, error) {
	return json.Marshal(spec)
}

func (spec *ModelTrainingSpec) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	res := json.Unmarshal(b, &spec)
	return res
}

func (in ModelTrainingStatus) Value() (driver.Value, error) {
	return json.Marshal(in)
}

func (in *ModelTrainingStatus) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	res := json.Unmarshal(b, &in)
	return res
}

// +kubebuilder:object:root=true

// ModelTraining is the Schema for the modeltrainings API
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.state"
// +kubebuilder:printcolumn:name="Toolchain",type="string",JSONPath=".spec.toolchain"
// +kubebuilder:printcolumn:name="VCS name",type="string",JSONPath=".spec.vcsName"
// +kubebuilder:printcolumn:name="Model name",type="string",JSONPath=".spec.model.name"
// +kubebuilder:printcolumn:name="Model version",type="string",JSONPath=".spec.model.version"
// +kubebuilder:printcolumn:name="Model image",type="string",JSONPath=".spec.image"
// +kubebuilder:resource:shortName=mt
type ModelTraining struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ModelTrainingSpec   `json:"spec,omitempty"`
	Status ModelTrainingStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ModelTrainingList contains a list of ModelTraining
type ModelTrainingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ModelTraining `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ModelTraining{}, &ModelTrainingList{})
}
