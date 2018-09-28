package ironic

import (
	"fmt"
	"github.com/micro/go-micro/errors"
	"net/http"
	"sync"
)

const DefaultID string = "ironic.error.default"

type ErrorManager struct {
	managers sync.Map
	// errors.Error
	baseID string `json`
}

func NewErrorManager(id string) *ErrorManager {
	if id == "" {
		id = DefaultID
	}
	return &ErrorManager{
		baseID: id,
	}
}

func (e *ErrorManager) WithID(id string) *ErrorManager {
	if foundEm, found := e.managers.Load(id); found {
		return foundEm.(*ErrorManager)
	}

	em := &ErrorManager{
		baseID: id,
	}

	e.managers.Store(id, em)

	return em
}

func (e *ErrorManager) New(code int32, format string, a ...interface{}) error {
	id := e.baseID
	if id == "" {
		id = DefaultID
	}
	return errors.New(id, fmt.Sprintf(format, a...), code)
}

func (e *ErrorManager) TemplateBadRequest(temp string, a ...interface{}) error {
	id := e.baseID
	if id == "" {
		id = DefaultID
	}

	return &errors.Error{
		Code:   400,
		Status: temp,
		Detail: fmt.Sprintf(temp, a...),
		Id:     id,
	}
}

func (e *ErrorManager) ActionError(temp string, a ...interface{}) error {
	id := e.baseID
	if id == "" {
		id = DefaultID
	}

	return &errors.Error{
		Code:   400,
		Status: temp,
		Detail: fmt.Sprintf(temp, a...),
		Id:     id,
	}
}

// BadRequest generates a 400 error.
func (e *ErrorManager) BadRequest(format string, a ...interface{}) error {
	return &errors.Error{
		Id:     e.baseID,
		Code:   400,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(400),
	}
}

// Unauthorized generates a 401 error.
func (e *ErrorManager) Unauthorized(format string, a ...interface{}) error {
	return &errors.Error{
		Id:     e.baseID,
		Code:   401,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(401),
	}
}

// Forbidden generates a 403 error.
func (e *ErrorManager) Forbidden(format string, a ...interface{}) error {
	return &errors.Error{
		Id:     e.baseID,
		Code:   403,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(403),
	}
}

// NotFound generates a 404 error.
func (e *ErrorManager) NotFound(format string, a ...interface{}) error {
	return &errors.Error{
		Id:     e.baseID,
		Code:   404,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(404),
	}
}

// InternalServerError generates a 500 error.
func (e *ErrorManager) InternalServerError(format string, a ...interface{}) error {
	return &errors.Error{
		Id:     e.baseID,
		Code:   500,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(500),
	}
}

// Conflict generates a 409 error.
func (e *ErrorManager) Conflict(format string, a ...interface{}) error {
	return &errors.Error{
		Id:     e.baseID,
		Code:   409,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(409),
	}
}
