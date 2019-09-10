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

package deployment

import (
	"github.com/legion-platform/legion/legion/operator/pkg/apis/legion/v1alpha1"
)

type ModelDeployment struct {
	// Model deployment id
	Id string `json:"id"`
	// Model deployment specification
	Spec v1alpha1.ModelDeploymentSpec `json:"spec,omitempty"`
	// Model deployment status
	Status *v1alpha1.ModelDeploymentStatus `json:"status,omitempty"`
}
