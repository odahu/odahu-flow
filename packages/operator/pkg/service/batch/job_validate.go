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
	"github.com/google/uuid"
	api_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/batch"
	odahu_errors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
)

const (
	EmptyServiceJobField = `field "%s"" must be set for BatchInferenceJob 
or as a default for a corresponding BatchInferenceService`
	ConnectionNotFound = `connection from "%s" with name "%s" is not found`
)

func DefaultJob(job *api_types.InferenceJob, service api_types.InferenceService) {

	job.ID = "job-" + service.ID + "-" + uuid.New().String()[:5]

	if len(job.Spec.BatchRequestID) == 0 {
		job.Spec.BatchRequestID = uuid.New().String()
	}

	if job.Spec.DataSource == nil {
		job.Spec.DataSource = service.Spec.DataSource
	}
	if job.Spec.OutputDestination == nil {
		job.Spec.OutputDestination = service.Spec.OutputDestination
	}
	if job.Spec.NodeSelector == nil {
		job.Spec.NodeSelector = service.Spec.NodeSelector
	}
	if job.Spec.Resources == nil {
		job.Spec.Resources = service.Spec.Resources
	}
}


// ValidateJobInput validates a job before it was defaulted by BatchInferenceService values
func ValidateJobInput(job api_types.InferenceJob) (errs []error) {

	if len(job.Spec.InferenceServiceID) == 0 {
		errs = append(errs, fmt.Errorf(EmptySpecFieldErrorMessage, "service"))
	}

	return errs
}


// ValidateJob validates a job after it was defaulted by BatchInferenceService values
func ValidateJob(
	job api_types.InferenceJob,
	connGetter ConnectionGetter,
	service api_types.InferenceService) (valErrs []error, err error) {

	if service.Spec.Triggers.Webhook == nil || !service.Spec.Triggers.Webhook.Enabled {
		valErrs = append(valErrs, fmt.Errorf("InferenceService: %s webhook trigger is disabled", service.ID))
	}

	// Validate empty values
	if len(job.Spec.BatchRequestID) == 0 {
		valErrs = append(valErrs, fmt.Errorf(EmptySpecFieldErrorMessage, "requestId"))
	}
	if job.Spec.DataSource == nil {
		valErrs = append(valErrs, fmt.Errorf(EmptyServiceJobField, "dataSource"))
	}
	if job.Spec.OutputDestination == nil {
		valErrs = append(valErrs, fmt.Errorf(EmptyServiceJobField, "outputDestination"))
	}

	// Validate that connections exist
	if job.Spec.DataSource != nil {
		name := job.Spec.DataSource.Connection
		_, err := connGetter.GetConnection(name, true)
		if !odahu_errors.IsNotFoundError(err) {
			return valErrs, err
		}
		valErrs = append(valErrs, fmt.Errorf(ConnectionNotFound, "dataSource", name))
	}
	if job.Spec.OutputDestination != nil {
		name := job.Spec.OutputDestination.Connection
		_, err := connGetter.GetConnection(name, true)
		if !odahu_errors.IsNotFoundError(err) {
			return valErrs, err
		}
		valErrs = append(valErrs, fmt.Errorf(ConnectionNotFound, "outputDestination", name))
	}

	return valErrs, err

}
