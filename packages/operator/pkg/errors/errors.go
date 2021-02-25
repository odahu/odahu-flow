/*
 * Copyright 2019 EPAM Systems
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package errors

import (
	"fmt"
	"net/http"
	"strings"
)

type NotFoundError struct {
	Entity string
}

func (nfe NotFoundError) Error() string {
	return fmt.Sprintf("entity %q is not found", nfe.Entity)
}

func IsNotFoundError(err error) bool {
	_, ok := err.(NotFoundError)
	return ok
}

type SerializationError struct{}

func (se SerializationError) Error() string {
	return "serialization is failed"
}

func IsSerializationError(err error) bool {
	_, ok := err.(SerializationError)
	return ok
}

type AlreadyExistError struct {
	Entity string
}

func (aee AlreadyExistError) Error() string {
	return fmt.Sprintf("entity %q already exists", aee.Entity)
}

func IsAlreadyExistError(err error) bool {
	_, ok := err.(AlreadyExistError)
	return ok
}

type ForbiddenError struct{}

func (aee ForbiddenError) Error() string {
	return "access forbidden"
}

type ExtendedForbiddenError struct{
	Message string
}

func (aee ExtendedForbiddenError) Error() string {
	return fmt.Sprintf("access forbidden: %v", aee.Message)
}

func IsForbiddenError(err error) bool {
	_, ok := err.(ForbiddenError)
	_, eok := err.(ExtendedForbiddenError)
	return ok || eok
}

type InvalidEntityError struct {
	Entity           string
	ValidationErrors []error
}

func (iee InvalidEntityError) Error() string {
	errorStrings := make([]string, 0, len(iee.ValidationErrors))
	for _, err := range iee.ValidationErrors {
		errorStrings = append(errorStrings, err.Error())
	}

	return fmt.Sprintf("entity %q is invalid; errors: %s", iee.Entity, strings.Join(errorStrings, ", "))
}


// Error means that spec of entity was touched (any field was changed).
// Could be raised in operations that require such behaviour
type SpecWasTouched struct {
	Entity string
}

func (swc SpecWasTouched) Error() string {
	return fmt.Sprintf("entity %q spec was changed", swc.Entity)
}

func IsSpecWasTouchedError(err error) bool {
	_, ok := err.(SpecWasTouched)
	return ok
}


type DeletingServiceHasJobs struct {
	// ID of BatchInferenceService
	Entity string
}

func (e DeletingServiceHasJobs) Error() string {
	return fmt.Sprintf(`Unable to delete service: "%s". Cause: there are child jobs`, e.Entity)
}

func (e DeletingServiceHasJobs) HTTPCode() int {
	return http.StatusBadRequest
}

type CreatingJobServiceNotFound struct {
	Entity string
	Service string
}

func (e CreatingJobServiceNotFound) Error() string {
	return fmt.Sprintf(`Unable to create job: "%s". There is no service with ID: %s`, e.Entity, e.Service)
}

func (e CreatingJobServiceNotFound) HTTPCode() int {
	return http.StatusNotFound
}

