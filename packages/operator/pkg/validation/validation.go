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

package validation

import (
	"fmt"
	connection "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection"
)

const (
	SpecSectionValidationFailedMessage = "\"Spec.%q\" validation errors: %s"
	EmptyValueStringError              = "%q parameter must be not empty"
)

func ValidateEmpty(parameterName, value string) error {
	if len(value) == 0 {
		return fmt.Errorf(EmptyValueStringError, parameterName)
	}
	return nil
}

func ValidateExistsInRepository(name string, repository connection.Repository) error {
	if len(name) > 0 {
		if _, odahuError := repository.GetConnection(name); odahuError != nil {
			return odahuError
		}
	}
	return nil
}
