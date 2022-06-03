package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/chzyer/readline"
	"github.com/peterbourgon/ff/v3/ffcli"
)

var setRegex = regexp.MustCompile(`^\s*set\s+(\w+)=(\w*)\s*$`)
var setUsage = strings.TrimSpace(`
Set: Storage store a key/value pair in the store, overwriting the value if <key> already exists

Usage:
	set <key>=[value]
`)

var getRegex = regexp.MustCompile(`^\s*get\s+(\w+)\s*$`)
var getUsage = strings.TrimSpace(`
Get: Retrieve the value of a key from the store

Usage:
	get <key>
`)

func clientCommand() *ffcli.Command {
	var port int
	var hostname string

	fs := flag.NewFlagSet("client", flag.ExitOnError)

	fs.IntVar(&port, "port", defaultPort, "the port that the server is listening on")
	fs.StringVar(&hostname, "hostname", defaultHostname, "the hostname of the server")

	cmd := &ffcli.Command{
		Name:       "client",
		ShortUsage: fmt.Sprintf("%s client [-port N] [-hostname foo.com]", os.Args[0]),
		ShortHelp:  "initiate a REPL session with a key/value server",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			baseAddress := fmt.Sprintf("%s:%d", hostname, port)

			endpoint, err := url.Parse(fmt.Sprintf("http://%s/v1/store", baseAddress))
			if err != nil {
				return fmt.Errorf("failed to parse %q: %w", baseAddress, err)
			}

			c := client{endpoint: endpoint}
			return c.runREPL()
		},
	}

	return cmd
}

type client struct {
	endpoint *url.URL
}

func (c *client) runREPL() error {
	rl, err := readline.New("> ")
	if err != nil {
		return fmt.Errorf("setting up REPL: %s", err)
	}

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if matches := setRegex.FindAllStringSubmatch(line, -1); len(matches) > 0 {
			m := matches[0]
			key, value := m[1], m[2]

			err := c.sendSet(key, value)
			if err != nil {
				fmt.Printf("error: failed to store %q=%q: %s\n", key, value, err)
			}

		} else if matches := getRegex.FindAllStringSubmatch(line, -1); len(matches) > 0 {
			m := matches[0]
			key := m[1]

			value, err := c.sendGet(key)
			if err != nil {
				fmt.Printf("error: failed to retrieve value for %q: %s\n", key, err)
			}

			fmt.Println(value)
		} else {
			fmt.Printf("error: invalid command %q\n", line)
			fmt.Printf("Known commands:\n%s\n\n%s\n", setUsage, getUsage)
		}
	}

	return nil
}

func (c *client) sendSet(key, value string) error {
	// deep copy of URL - this is safe according to: https://github.com/golang/go/issues/38351
	u := *c.endpoint

	query := u.Query()
	query.Set("key", key)
	query.Set("value", value)
	u.RawQuery = query.Encode()

	r, err := http.Post(u.String(), "text/plain", nil)
	if err != nil {
		return fmt.Errorf("POST %s: %w", c.endpoint, err)
	}

	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		bs, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("decoding response body from failed request: %w", err)
		}

		return errors.New(string(bs))
	}

	return nil
}

func (c *client) sendGet(key string) (value string, err error) {
	// deep copy of URL - this is safe according to: https://github.com/golang/go/issues/38351
	u := *c.endpoint

	query := u.Query()
	query.Set("key", key)
	u.RawQuery = query.Encode()

	r, err := http.Get(u.String())
	if err != nil {
		return "", fmt.Errorf("GET %s: %w", c.endpoint, err)
	}

	defer r.Body.Close()
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("reading response body: %w", err)
	}

	text := string(raw)
	if r.StatusCode != http.StatusOK {
		return "", errors.New(text)
	}

	return text, nil
}
