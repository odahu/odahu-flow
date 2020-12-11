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

type ModelRoute struct {
	// Model route id
	ID string `json:"id"`
	// Default routes cannot be deleted by user. They are managed by system
	// One ModelDeployment has exactly one default Route that gives 100% traffic to the model
	Default bool `json:"default,omitempty"`
	// Deletion mark
	DeletionMark bool `json:"deletionMark,omitempty" swaggerignore:"true"`
	// CreatedAt
	CreatedAt time.Time `json:"createdAt,omitempty"`
	// UpdatedAt
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	// Model route specification
	Spec v1alpha1.ModelRouteSpec `json:"spec,omitempty"`
	// Model route status
	Status v1alpha1.ModelRouteStatus `json:"status,omitempty"`
}

func (in ModelRoute) Value() (driver.Value, error) {
	return json.Marshal(in)
}

func (in *ModelRoute) Scan(value interface{}) error {
	switch b := value.(type) {
	case nil:
		return nil
	case []byte:
		return json.Unmarshal(b, &in)
	default:
		return errors.New("type assertion to []byte or nil is failed")
	}
}
