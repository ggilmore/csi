package app

import (
	"fmt"
	"io"
)

type Scanner struct {
	start, current, line int

	errs ErrorList

	input  []rune
	tokens []Token
}

func NewScanner(r io.Reader) (*Scanner, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}

	text := []rune(string(b))

	return &Scanner{
		input: text,
	}, nil
}

func (s *Scanner) Scan() ([]Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{
		Kind: KindEOF,
		Line: s.line,
	})

	return s.tokens, s.errs.ErrorOrNil()
}

func (s *Scanner) scanToken() {
	c := s.advance()

	switch c {
	case ' ', '\r', '\t':
		break

	default:
		if s.isWordCharacter(c) {
			s.word()
			break
		}

		s.errs = append(s.errs, &parseError{
			s.line,
			fmt.Sprintf("unexpected character %q", c),
		})
	}
}

func (s *Scanner) word() {
	for s.isWordCharacter(s.peek()) {
		s.advance()
	}

	contents := s.input[s.start:s.current]
	s.addTokenLiteral(KindWord, contents)
}

const null = '\x00'

func (s *Scanner) isWordCharacter(c rune) bool {
	switch c {
	case ' ', null:
		return false

	default:
		return true
	}
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return null
	}

	return s.input[s.current]
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.input)
}

func (s *Scanner) advance() rune {
	c := s.input[s.current]
	s.current++

	return c
}

func (s *Scanner) addToken(kind TokenType) {
	s.addTokenLiteral(kind, nil)
}

func (s *Scanner) addTokenLiteral(kind TokenType, literal interface{}) {
	s.tokens = append(s.tokens, Token{
		Kind: kind,

		Lexeme:  string(s.input[s.start:s.current]),
		Literal: literal,

		Line: s.line,
	})
}
