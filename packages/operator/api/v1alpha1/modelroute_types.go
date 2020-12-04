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

type ModelDeploymentTarget struct {
	// Model Deployment name
	Name string `json:"mdName"`
	// The proportion of traffic to be forwarded to the Model Deployment.
	Weight *int32 `json:"weight,omitempty"`
}

type ModelRouteSpec struct {
	// URL prefix for model deployment. For example: /custom/test
	// Prefix must start with slash
	// "/feedback" and "/model" are reserved for internal usage
	URLPrefix string `json:"urlPrefix"`
	// Mirror HTTP traffic to a another Model deployment in addition to forwarding
	// the requests to the model deployments.
	Mirror *string `json:"mirror,omitempty"`
	// A http rule can forward traffic to Model Deployments.
	ModelDeploymentTargets []ModelDeploymentTarget `json:"modelDeployments"`
}

type ModelRouteState string

const (
	ModelRouteStateReady      = ModelRouteState("Ready")
	ModelRouteStateProcessing = ModelRouteState("Processing")
)

// ModelRouteStatus defines the observed state of ModelRoute
type ModelRouteStatus struct {
	// Full url with prefix to a model deployment service
	EdgeURL string `json:"edgeUrl"`
	// State of Model Route
	State ModelRouteState `json:"state"`
	// Info about create and update
	//CreatedAt *metav1.Time `json:"createdAt,omitempty"`
	//UpdatedAt *metav1.Time `json:"updatedAt,omitempty"`
	Modifiable `json:",inline"`
}

func (in ModelRouteSpec) Value() (driver.Value, error) {
	return json.Marshal(in)
}

func (in *ModelRouteSpec) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	res := json.Unmarshal(b, &in)
	return res
}

func (in ModelRouteStatus) Value() (driver.Value, error) {
	return json.Marshal(in)
}

func (in *ModelRouteStatus) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	res := json.Unmarshal(b, &in)
	return res
}


// +kubebuilder:object:root=true

// ModelRoute is the Schema for the modelroutes API
// +kubebuilder:printcolumn:name="Edge URL",type="string",JSONPath=".status.edgeUrl"
// +kubebuilder:printcolumn:name="Routes",type="string",JSONPath=".spec.modelDeployments"
// +kubebuilder:printcolumn:name="Mirror",type="string",JSONPath=".spec.mirror"
// +kubebuilder:resource:shortName=mr
type ModelRoute struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ModelRouteSpec   `json:"spec,omitempty"`
	Status ModelRouteStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ModelRouteList contains a list of ModelRoute
type ModelRouteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ModelRoute `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ModelRoute{}, &ModelRouteList{})
}
