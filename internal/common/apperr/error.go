package apperr

import (
	"errors"
	"net/http"
)

type HttpError interface {
	error
	HttpStatus() int
}

type implementation struct {
	error
	status int
}

func (err *implementation) Unwrap() error {
	return err.error
}

func (err *implementation) HttpStatus() int {
	return err.status
}

func WithHTTPStatus(err error, status int) error {
	return &implementation{
		error:  err,
		status: status,
	}
}

func HTTPStatus(err error) int {
	if err == nil {
		return 0
	}

	if appErr, ok := GetAppError(err); ok {
		return appErr.HttpStatus()
	}

	return http.StatusInternalServerError
}

func GetAppError(err error) (HttpError, bool) {
	var converted HttpError
	if errors.As(err, &converted) {
		return converted, true
	}

	return nil, false
}
