package httperrs

import (
	"net/http"

	"github.com/pkg/errors"
)

// NotFoundError lets us know that an error constitutes
// http.StatusNotFound
type NotFoundError interface {
	NotFound()
}

type notFoundErr struct {
	err error
}

// NotFound wraps a go error with our interface
func NotFound(err error, msg string) error {
	if msg == "" {
		return notFoundErr{
			err: err,
		}
	}
	return notFoundErr{
		err: errors.Wrap(err, msg),
	}
}

// Error makes our errors satisfy the stdlib
// error interface
func (ne notFoundErr) Error() string {
	return ne.err.Error()
}

// NotFound shows this error type should reflect
// a not found error
func (ne notFoundErr) NotFound() {}

// InternalServerError lets us know that an error constitutes
// http.StatusInternalServerError
type InternalServerError interface {
	InternalServer()
}

type internalServerErr struct {
	err error
}

// InternalServer wraps a go error with our interface
func InternalServer(err error, msg string) error {
	if msg == "" {
		return internalServerErr{
			err: err,
		}
	}
	return internalServerErr{
		err: errors.Wrap(err, msg),
	}
}

// Error makes our errors satisfy the stdlib
// error interface
func (ie internalServerErr) Error() string {
	return ie.err.Error()
}

// InternalServer shows this error type should reflect
// an internal server error
func (ie internalServerErr) InternalServer() {}

// BadRequestError lets us know that an error constitutes
// http.StatusBadRequest
type BadRequestError interface {
	BadRequest()
}

type badRequestErr struct {
	err error
}

// BadRequest wraps a go error with our interface
func BadRequest(err error, msg string) error {
	if msg == "" {
		return badRequestErr{
			err: err,
		}
	}
	return badRequestErr{
		err: errors.Wrap(err, msg),
	}
}

// Error makes our errors satisfy the stdlib
// error interface
func (be badRequestErr) Error() string {
	return be.err.Error()
}

// BadRequest shows this error type should reflect
// a bad request error
func (be badRequestErr) BadRequest() {}

// StatusCode takes an error and tries to determine what
// http Status Code it corresponds to
func StatusCode(err error) int {
	switch err.(type) {
	case BadRequestError:
		return http.StatusBadRequest
	case InternalServerError:
		return http.StatusInternalServerError
	case NotFoundError:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
