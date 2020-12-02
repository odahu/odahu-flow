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
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"time"
)

type ModelDeployment struct {
	// Model deployment id
	ID string `json:"id"`
	// Deletion mark
	DeletionMark bool `json:"deletionMark,omitempty" swaggerignore:"true"`
	// CreatedAt
	CreatedAt time.Time `json:"createdAt,omitempty"`
	// UpdatedAt
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	// Model deployment specification
	Spec v1alpha1.ModelDeploymentSpec `json:"spec,omitempty"`
	// Model deployment status
	Status v1alpha1.ModelDeploymentStatus `json:"status,omitempty"`
}

func (in ModelDeployment) Value() (driver.Value, error) {
	return json.Marshal(in)
}


func (in *ModelDeployment) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	res := json.Unmarshal(b, &in)
	return res
}
