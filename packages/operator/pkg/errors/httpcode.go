package errors

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"net/http"
)

func CalculateHTTPStatusCode(err error) int {

	// TODO: Remove StatusError because after ensure that we wrap all kubernetes errors by odahu errors
	errStatus, ok := err.(*errors.StatusError)
	if ok {
		return int(errStatus.ErrStatus.Code)
	}

	if IsNotFoundError(err) {
		return http.StatusNotFound
	}

	if IsAlreadyExistError(err) {
		return http.StatusConflict
	}

	if IsForbiddenError(err) {
		return http.StatusForbidden
	}

	if _, ok = err.(InvalidEntityError); ok {
		return http.StatusBadRequest
	}

	if _, ok = err.(DeletingServiceHasJobs); ok {
		return http.StatusBadRequest
	}

	if _, ok = err.(CreatingJobServiceNotFound); ok {
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}

