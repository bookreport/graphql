package gqlerrors

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"runtime/debug"
	"strings"

	"github.com/bookreport/graphql/language/location"
)

type FormattedError struct {
	Message             string                    `json:"message"`
	Locations           []location.SourceLocation `json:"locations"`
	LocalizedStackTrace string
}

func (g FormattedError) Error() string {
	return g.Message
}

func MakeLocalizedStackTrace(msg string, stackTrace []byte, fieldNames ...string) string {
	var stack bytes.Buffer
	stack.WriteString("\n\n|---------------------------------------------------------------------------\n")
	stack.WriteString("|  graphql error")
	if len(fieldNames) > 0 {
		stack.WriteString(" on field ")
		stack.WriteString(strings.Join(fieldNames, " - "))
	}
	stack.WriteString("\n")
	stack.WriteString("|---------------------------------------------------------------------------\n")
	stack.WriteString("|\n|  ")
	stack.WriteString(msg)
	stack.WriteString("\n|  ...\n")
	scanner := bufio.NewScanner(bytes.NewReader(stackTrace))
	dir, _ := os.Getwd()
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "graphql-go") {
			continue
		}
		if strings.Contains(line, "runtime") {
			continue
		}
		stack.WriteString("|  ")
		stack.WriteString(strings.Replace(line, dir, "", -1))
		stack.WriteString("\n")
	}
	stack.WriteString("|  ...\n")
	stack.WriteString("|---------------------------------------------------------------------------\n")
	return stack.String()
}

func NewFormattedError(message string) FormattedError {
	err := errors.New(message)
	return FormatError(err)
}

func FormatError(err error, fieldNames ...string) FormattedError {
	stackTrace := debug.Stack()
	switch err := err.(type) {
	case FormattedError:
		return err
	case *Error:
		return FormattedError{
			Message:             err.Error(),
			Locations:           err.Locations,
			LocalizedStackTrace: MakeLocalizedStackTrace(err.Error(), stackTrace, fieldNames...),
		}
	case Error:
		return FormattedError{
			Message:             err.Error(),
			Locations:           err.Locations,
			LocalizedStackTrace: MakeLocalizedStackTrace(err.Error(), stackTrace, fieldNames...),
		}
	default:
		return FormattedError{
			Message:             err.Error(),
			Locations:           []location.SourceLocation{},
			LocalizedStackTrace: MakeLocalizedStackTrace(err.Error(), stackTrace, fieldNames...),
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
