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
	"errors"
	"fmt"
	odahuv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	connection "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/kubernetes"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/kubernetes/pkg/apis/core/v1/validation"
	"regexp"
)

const (
	SpecSectionValidationFailedMessage = "\"Spec.%q\" validation errors: %s"
	EmptyValueStringError              = "empty %q"

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

var idRegex = regexp.MustCompile("^[a-z]([-a-z0-9]{0,61}[a-z0-9])?$")
var ErrIDValidation = errors.New("ID is not valid")

// Id restrictions:
//  * contain at most 63 characters
//  * contain only lowercase alphanumeric characters or ‘-’
//  * start with an alpha character
//  * end with an alphanumeric character
func ValidateID(id string) error {
	if idRegex.MatchString(id) {
		return nil
	}
	err := ErrIDValidation
	if len(id) == 0 {
	    err = multierr.Append(err, ValidateEmpty("ID", id))
	}
	return err
}

// K8s label must start/end with alphanumeric character, can consist of
var k8sLabelRegex = regexp.MustCompile("^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$")
var LabelValueValidationErrorTemplate = "%s must be valid Kubernetes label, i.e. match this pattern: %s"

func ValidateK8sLabel(label string) error {
	if k8sLabelRegex.MatchString(label) {
		return nil
	}
	return fmt.Errorf(LabelValueValidationErrorTemplate, label, k8sLabelRegex)
}


func ValidateResources(resources *odahuv1alpha1.ResourceRequirements, gpuResName string) (err error) {

	coreV1Resources, err := kubernetes.ConvertOdahuflowResourcesToK8s(resources, gpuResName)

	fErrs := validation.ValidateResourceRequirements(&coreV1Resources, field.NewPath("resources"))

	for _, fErr := range fErrs {
		err = multierr.Append(err, fErr)
	}

	return
}