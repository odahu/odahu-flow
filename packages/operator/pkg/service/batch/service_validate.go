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
	"fmt"
	api_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/batch"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/kubernetes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/validation"
	"go.uber.org/multierr"
)

const (
	EmptySpecFieldErrorMessage = "%s must be non-empty"
)


func validateRequiredFields(bis api_types.InferenceService) (err error) {

	idErr := validation.ValidateID(bis.ID)
	if idErr != nil {
		err = multierr.Append(err, idErr)
	}

	if len(bis.Spec.Image) == 0 {
		err = multierr.Append(err, fmt.Errorf(EmptySpecFieldErrorMessage, "image"))
	}
	if len(bis.Spec.Command) == 0 {
		err = multierr.Append(err, fmt.Errorf(EmptySpecFieldErrorMessage, "command"))
	}
	registry := bis.Spec.ModelRegistry
	switch  {
	case registry.Remote != nil:
		if len(registry.Remote.ModelConnection) == 0 {
			err = multierr.Append(err, fmt.Errorf(EmptySpecFieldErrorMessage, "modelRegistry.remote.modelConnection"))
		}
	case registry.Local != nil:
		if len(registry.Local.ModelMeta.Name) == 0 {
			err = multierr.Append(err, fmt.Errorf(EmptySpecFieldErrorMessage, "modelRegistry.local.meta.name"))
		}
		if len(registry.Local.ModelMeta.Version) == 0 {
			err = multierr.Append(err, fmt.Errorf(EmptySpecFieldErrorMessage, "modelRegistry.local.meta.version"))
		}
	default:
		err = multierr.Append(err,
			fmt.Errorf("whether modelRegistry.local.meta.name " +
				"or modelRegistry.local.meta.version must be defined for embedded models"))
	}
	if bis.Spec.ModelRegistry.Remote == nil || bis.Spec.ModelRegistry.Local == nil {
		err = multierr.Append(err, fmt.Errorf(EmptySpecFieldErrorMessage, "modelRegistry.connection"))
	}


	if bis.Spec.Resources != nil {
		_, resValidationErr := kubernetes.ConvertOdahuflowResourcesToK8s(bis.Spec.Resources, "nvidia")
		if resValidationErr != nil {
			err = multierr.Append(err, resValidationErr)
		}
	}


	return err
}


func ValidateCreateUpdate(bis api_types.InferenceService) (errs []error) {

	var err error

	err = multierr.Append(err, validateRequiredFields(bis))

	if err != nil {
		return multierr.Errors(err)
	}
	return nil
}

func DefaultCreate(bis *api_types.InferenceService) {

	// By default webhook trigger is enabled
	if bis.Spec.Triggers.Webhook == nil {
		bis.Spec.Triggers.Webhook = &api_types.PredictorWebhookTrigger{Enabled: true}
	}
}
