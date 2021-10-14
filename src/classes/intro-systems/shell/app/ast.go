package app

import (
	"fmt"
	"strings"
)

type Expression interface {
	String() string
	isExpression()
}

type Program struct {
	Commands []CommandExpression
}

type CommandExpression struct {
	Program Token
	Args    []Token
}

func (c *CommandExpression) isExpression() {}
func (c *CommandExpression) String() string {
	var argumentStrings []string
	for _, a := range c.Args {
		argumentStrings = append(argumentStrings, a.String())
	}

	return fmt.Sprintf("CommandExpression{%q -> %s}", c.Program, strings.Join(argumentStrings, ","))
}

var (
	_ Expression = &CommandExpression{}
)
