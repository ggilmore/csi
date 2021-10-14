package app

import "fmt"

// command: word (" " word)+
// word: any non-space, non-null character

type Parser struct {
	current int
	tokens  []Token
}

func NewParser(input []Token) *Parser {
	return &Parser{
		tokens:  input,
		current: 0,
	}
}

func (p *Parser) Parse() ([]*CommandExpression, error) {
	var commands []*CommandExpression

	for !p.isAtEnd() {
		command, err := p.command()
		if err != nil {
			return nil, err
		}

		commands = append(commands, command)
	}

	return commands, nil
}

func (p *Parser) command() (*CommandExpression, error) {
	command := &CommandExpression{}

	if !p.match(KindWord) {
		return nil, fmt.Errorf("expected program name")
	}
	command.Program = p.previous()

	for p.match(KindWord) {
		command.Args = append(command.Args, p.previous())
	}

	return command, nil
}

// match consume tokens while any of the provided
// kinds match the current token
func (p *Parser) match(kinds ...TokenType) bool {
	for _, k := range kinds {
		if p.check(k) {
			p.advance()
			return true
		}
	}

	return false
}

// advance consumes the current token and
// advances the token stream
func (p *Parser) advance() Token {
	out := p.peek()
	p.current++

	return out
}

// previous returns to
func (p *Parser) previous() Token {
	out := p.tokens[p.current-1]
	return out
}

func (p *Parser) check(kind TokenType) bool {
	return p.peek().Kind == kind
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Kind == KindEOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}
