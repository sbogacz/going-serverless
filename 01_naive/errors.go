package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
)

type notFoundError interface {
	NotFound()
}

type notFoundErr struct {
	err error
}

func newNotFoundErr(err error, msg string) error {
	if msg == "" {
		return notFoundErr{
			err: err,
		}
	}
	return notFoundErr{
		err: errors.Wrap(err, msg),
	}
}

func (ne notFoundErr) Error() string {
	return ne.err.Error()
}

func (ne notFoundErr) NotFound() {}

type internalServerError interface {
	InternalServer()
}

type internalServerErr struct {
	err error
}

func newInternalServerErr(err error, msg string) error {
	if msg == "" {
		return internalServerErr{
			err: err,
		}
	}
	return internalServerErr{
		err: errors.Wrap(err, msg),
	}
}

func (ie internalServerErr) Error() string {
	return ie.err.Error()
}

func (ie internalServerErr) InternalServer() {}

type badRequestError interface {
	BadRequest()
}

type badRequestErr struct {
	err error
}

func newBadRequestErr(err error, msg string) error {
	if msg == "" {
		return badRequestErr{
			err: err,
		}
	}
	return badRequestErr{
		err: errors.Wrap(err, msg),
	}
}

func (be badRequestErr) Error() string {
	return be.err.Error()
}

func (be badRequestErr) BadRequest() {}

func errorResponse(err error) *events.APIGatewayProxyResponse {
	var code int
	switch err.(type) {
	case badRequestError:
		code = http.StatusBadRequest
	case internalServerError:
		code = http.StatusInternalServerError
	case notFoundError:
		code = http.StatusNotFound
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: code,
	}
}
