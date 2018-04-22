package golog

import (
	"errors"

	"github.com/joaosoft/go-error/service"
)

type Addition struct {
	message string
}

// newAddition ...
func newAddition(message string) IAddition {
	addition := &Addition{
		message: message,
	}

	return addition
}

// ToError
func (addition *Addition) ToError(err *error) IAddition {
	*err = errors.New(addition.message)
	return addition
}

// ToErrorData
func (addition *Addition) ToErrorData(err *goerror.ErrorData) IAddition {
	newErr := errors.New(addition.message)
	err.AddError(newErr)

	return addition
}