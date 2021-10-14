package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ggilmore/csi/src/classes/intro-systems/shell/app"
)

var command string

func main() {
	flag.StringVar(&command, "c", "", "the command to execute")
	flag.Parse()

	if len(command) > 0 {
		runner := newRunner()
		err := runner.Run(strings.NewReader(command))
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}

		return
	}

	runPrompt(os.Stdin)
}

func runPrompt(r io.Reader) {
	runner := newRunner()
	s := bufio.NewScanner(r)

	prompt := "> "
	fmt.Print(prompt)

	for s.Scan() {
		line := s.Text()

		err := runner.Run(strings.NewReader(line))
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		fmt.Printf("\n%s", prompt)
	}

	if err := s.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "while processing input: %s\n", err)
		os.Exit(1)
	}
}

type runner struct {
	interpreter *app.Interpreter
}

func newRunner() *runner {
	return &runner{
		interpreter: app.NewInterpreter(),
	}
}

func (r *runner) Run(input io.Reader) error {
	s, err := app.NewScanner(input)
	if err != nil {
		return fmt.Errorf("intializing scanner: %w", err)
	}

	tokens, err := s.Scan()
	if err != nil {
		return fmt.Errorf("scanning for tokens: %w", err)
	}

	statements, err := app.NewParser(tokens).Parse()
	if err != nil {
		return fmt.Errorf("while parsing: %w", err)
	}

	err = r.interpreter.Interpret(statements)
	if err != nil {
		return fmt.Errorf("while interpreting: %w", err)
	}

	return nil
}
