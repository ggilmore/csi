package app

import "fmt"

type TokenType int

const (
	KindEOF TokenType = iota

	KindWord
)

func (t TokenType) String() string {
	switch t {
	case KindEOF:
		return "EOF"
	case KindWord:
		return "Word"
	}

	panic("unhandled token type")
}

type Token struct {
	Kind      TokenType
	Pos, Line int

	Lexeme  string
	Literal interface{}
}

func (t Token) String() string {
	return fmt.Sprintf("Token<%s>{%q, [%v]}", t.Kind.String(), t.Lexeme, t.Literal)
}
