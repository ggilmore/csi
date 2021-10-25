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

func TestStructField(t *testing.T) {
	var input point
	f := fuzz.New()

	for i := 0; i < 1000; i++ {
		f.Fuzz(&input)

		expected := input.Y
		actual := StructField(input)

		if actual != expected {
			t.Errorf("for (%+v).Y got: %d, wanted: %d", input, actual, expected)
		}
	}
}
func TestSumSlice(t *testing.T) {
	var input []int
	f := fuzz.New().NumElements(0, 10)

	for i := 0; i < 1000; i++ {
		f.Fuzz(&input)

		expected := correctSumSlice(input)
		actual := SumSlice(input)

		if actual != expected {
			t.Errorf("for (%+v).Y got: %d, wanted: %d", input, actual, expected)
		}
	}
}

func correctSumSlice(xs []int) int {
	sum := 0
	for _, x := range xs {
		sum += x
	}

	return sum
}
