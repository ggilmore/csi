package app

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParser(t *testing.T) {
	programName := Token{Kind: KindWord, Literal: "hi", Lexeme: "hi"}
	arg1 := Token{Kind: KindWord, Literal: "arg1", Lexeme: "arg1"}
	arg2 := Token{Kind: KindWord, Literal: "arg2", Lexeme: "arg2"}

	for _, test := range []struct {
		name     string
		input    []Token
		expected []Command
	}{
		{
			name:     "(empty)",
			input:    []Token{},
			expected: nil,
		},
		{
			name: "(simpleCommand)",
			input: []Token{
				programName,
				arg1,
				arg2,
			},
			expected: []Command{{Program: programName, Args: []Token{arg1, arg2}}},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.input = append(test.input, Token{Kind: KindEOF})

			p := NewParser(test.input)

			actual, err := p.Parse()
			if err != nil {
				t.Errorf("received unexpected error while parsing: %s", err)
			}

			if diff := cmp.Diff(test.expected, actual); diff != "" {
				t.Errorf("unexpected diff (-expected +actual)\n:%s", diff)
			}
		})
	}
}
