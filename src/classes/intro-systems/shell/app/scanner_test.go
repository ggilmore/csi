package app

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func eofToken(line int) Token {
	return Token{
		Kind: KindEOF,

		Line: line,
	}
}

func TestScanner(t *testing.T) {
	for _, test := range []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "(empty)",
			input: "",
			expected: []Token{
				{Kind: KindEOF},
			},
		},
		{
			name:  "simple command",
			input: "hi",
			expected: []Token{
				{Kind: KindWord, Lexeme: "hi", Literal: "hi"},
				{Kind: KindEOF},
			},
		},
		{
			name:  "command with args",
			input: "hi             arg          ",
			expected: []Token{
				{Kind: KindWord, Lexeme: "hi", Literal: "hi"},
				{Kind: KindWord, Lexeme: "arg", Literal: "arg"},
				{Kind: KindEOF},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			s, err := NewScanner(strings.NewReader(test.input))
			if err != nil {
				t.Fatalf("failed to initialize scanner with input %q: %s", test.input, err)
			}

			actual, err := s.Scan()
			if err != nil {
				t.Errorf("received unexpected error while scanning: %s", err)
			}

			if diff := cmp.Diff(test.expected, actual, cmpopts.IgnoreFields(Token{}, "Literal")); diff != "" {
				t.Errorf("unexpected diff (-expected +actual)\n:%s", diff)
			}
		})
	}
}
