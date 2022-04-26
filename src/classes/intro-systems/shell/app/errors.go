package app

import (
	"fmt"
)

type ErrorList []error

func (e *ErrorList) Add(line int, message string) {
	*e = append(*e, &parseError{line, message})
}

func (e ErrorList) Error() string {
	if len(e) == 0 {
		return "no errors"
	}

	out := fmt.Sprintf("There were %d error(s)\n", len(e))

	for _, err := range e {
		out += fmt.Sprintf("- %s\n", err.Error())
	}

	return out
}

func (e ErrorList) ErrorOrNil() error {
	if len(e) == 0 {
		return nil
	}

	return e
}

type RecoverableError interface {
	IsRecoverableError()
	Error() string
}

type parseError struct {
	Line    int
	Message string
}

func (e *parseError) Error() string {
	return fmt.Sprintf("[line %d] Error: %s", e.Line, e.Message)
}

type runtimeError struct {
	Message string
}

func (e *runtimeError) Error() string {
	return e.Message
}

func (e *parseError) IsRecoverableError()   {}
func (e *runtimeError) IsRecoverableError() {}

var (
	_ RecoverableError = &parseError{}
	_ RecoverableError = &runtimeError{}
)
