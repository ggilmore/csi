package app

import (
	"fmt"
	"os"
	"strconv"
)

type builtin interface {
	Run(args []string) error
}

type exit struct{}

func (e exit) Run(args []string) error {
	if len(args) == 0 {
		return &runtimeError{"exit: expected an argument"}
	}

	if len(args) > 1 {
		return &runtimeError{"exit: too many arguments"}
	}

	rawCode := args[0]
	code, err := strconv.Atoi(rawCode)
	if err != nil {
		return &runtimeError{
			fmt.Sprintf("exit: can't convert %q to integer: %s", rawCode, err),
		}
	}

	os.Exit(code)
	return nil
}

type cd struct{}

func (c cd) Run(args []string) error {
	if len(args) > 1 {
		return &runtimeError{"cd: too many arguments"}
	}

	var target string

	if len(args) == 1 {
		target = args[0]
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return &runtimeError{
				fmt.Sprintf("cd: unable to get home directory: %s", err),
			}
		}

		target = home
	}

	err := os.Chdir(target)
	if err != nil {
		return &runtimeError{
			fmt.Sprintf("cd: can't change directories to %q: %s", target, err),
		}
	}

	return nil
}

var CDBuiltin = &cd{}
var ExitBuiltin = &exit{}
