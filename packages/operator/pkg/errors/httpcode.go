package errors

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"net/http"
)

// We should develop a custom exception on the repository layer.
// But we rely on kubernetes exceptions for now.
// TODO: implement Odahuflow exceptions
func CalculateHTTPStatusCode(err error) int {
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

	return http.StatusInternalServerError
}

