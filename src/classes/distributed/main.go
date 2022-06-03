package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"
)

const defaultPort int = 2222
const defaultHostname string = "localhost"

func main() {
	serverCmd, err := serverCommand()
	if err != nil {
		die("failed to initialize server command: %s", err.Error())
	}

	clientCmd := clientCommand()
	root := &ffcli.Command{
		ShortUsage:  fmt.Sprintf("%s <subcommand>", os.Args[0]),
		Subcommands: []*ffcli.Command{serverCmd, clientCmd},
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
	}

	if err := root.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}

		os.Exit(1)
	}
}

func die(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintf(os.Stderr, "\n")

	os.Exit(1)
}
