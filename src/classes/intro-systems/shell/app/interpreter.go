package app

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type Interpreter struct{}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Interpret(program *Program) error {
	for _, c := range program.Commands {
		err := i.execute(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) execute(c Command) error {
	name := c.Program.Lexeme
	path, err := exec.LookPath(name)

	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return i.runtimeErrorf("%q: command not found", name)
		}

		return i.runtimeErrorf("looking up %q in $PATH: %s", name, err)
	}

	arguments := []string{path}
	for _, a := range c.Args {
		arguments = append(arguments, a.Lexeme)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return i.runtimeErrorf("failed to get current working directory: %s", err)
	}

	attr := &syscall.ProcAttr{
		Dir: cwd,
		Env: os.Environ(),
		Files: []uintptr{
			os.Stdin.Fd(),
			os.Stdout.Fd(),
			os.Stderr.Fd(),
		},
	}

	pid, err := syscall.ForkExec(path, arguments, attr)
	if err != nil {
		return i.runtimeErrorf("running %q: %s", path, err)
	}

	var ws syscall.WaitStatus
	_, err = syscall.Wait4(pid, &ws, 0, nil)
	if err != nil {
		return i.runtimeErrorf("waiting on %q (pid %d): %s", path, pid, err)
	}

	return nil
}

func (i *Interpreter) runtimeErrorf(format string, a ...interface{}) error {
	return &runtimeError{
		Message: fmt.Sprintf(format, a...),
	}
}
