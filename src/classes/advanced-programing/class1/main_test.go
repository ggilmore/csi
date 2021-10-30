package main

import (
	"math"
	"os"
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
	if os.Getenv("CI") == "true" {
		t.Skip("this is unsafe pointer math")
	}

	var input []int
	f := fuzz.New().NumElements(0, 10)

	for i := 0; i < 1000; i++ {
		f.Fuzz(&input)

		expected := correctSumSlice(input)
		actual := SumSlice(input)

		if actual != expected {
			t.Errorf("for %+v got: %d, wanted: %d", input, actual, expected)
		}
	}
}

func TestMapMax(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("this is unsafe pointer math")
	}

	var input map[int]int
	f := fuzz.New().NumElements(0, 10)

	for i := 0; i < 1; i++ {
		f.Fuzz(&input)

		expected := correctMapMax(input)
		actual := MapMax(input)

		if actual != expected {
			t.Errorf("for %+v got: %d, wanted: %d", input, actual, expected)
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

func correctMapMax(xs map[int]int) int {
	max := math.MinInt

	for _, v := range xs {
		if v > max {
			max = v
		}
	}

	return max
}
