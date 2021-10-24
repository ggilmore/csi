package main

import (
	"testing"

	fuzz "github.com/google/gofuzz"
)

func TestStringLen(t *testing.T) {
	var input string
	f := fuzz.New()

	for i := 0; i < 1000; i++ {
		f.Fuzz(&input)

		expected := len(input)
		actual := StringLength(input)

		if actual != expected {
			t.Errorf("for %q got: %d, wanted: %d", input, actual, expected)
		}
	}
}
