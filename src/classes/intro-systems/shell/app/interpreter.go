package app

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type Interpreter struct {
}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Interpret(input []*CommandExpression) error {
	for _, s := range input {
		err := i.execute(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) execute(c *CommandExpression) error {
	name := c.Program.Lexeme
	path, err := exec.LookPath(name)

	if err != nil {
		var pathErr *os.PathError
		if !errors.As(err, &pathErr) {
			return fmt.Errorf("while looking up %q in PATH: %w", name, err)
		}

		return err
	}

	arguments := []string{path}
	for _, arg := range c.Args {
		arguments = append(arguments, arg.Lexeme)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
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
		return fmt.Errorf("executing %q: %w", path, err)
	}

	var ws syscall.WaitStatus
	_, err = syscall.Wait4(pid, &ws, 0, nil)
	if err != nil {
		return err
	}

	return nil
}
