package gqlerrors

import (
	"errors"
	"runtime/debug"

	"github.com/graphql-go/graphql/language/location"
)

type FormattedError struct {
	Message   string                    `json:"message"`
	Locations []location.SourceLocation `json:"locations"`
	Stack     []byte
}

func (g FormattedError) Error() string {
	return g.Message
}

func NewFormattedError(message string) FormattedError {
	err := errors.New(message)
	return FormatError(err)
}

func FormatError(err error) FormattedError {
	switch err := err.(type) {
	case FormattedError:
		return err
	case *Error:
		return FormattedError{
			Message:   err.Error(),
			Locations: err.Locations,
			Stack:     debug.Stack(),
		}
	case Error:
		return FormattedError{
			Message:   err.Error(),
			Locations: err.Locations,
			Stack:     debug.Stack(),
		}
	default:
		return FormattedError{
			Message:   err.Error(),
			Locations: []location.SourceLocation{},
			Stack:     debug.Stack(),
		}
	}
}

func FormatErrors(errs ...error) []FormattedError {
	formattedErrors := []FormattedError{}
	for _, err := range errs {
		formattedErrors = append(formattedErrors, FormatError(err))
	}
	return formattedErrors
}
