package app

import (
	"fmt"
)

type ErrorList []error

func (e *ErrorList) Add(line int, message string) {
	*e = append(*e, &ParseError{line, message})
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

type ParseError struct {
	Line    int
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("[line %d] Error: %s", e.Line, e.Message)
}
