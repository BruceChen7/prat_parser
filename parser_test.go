package parser

import (
	"testing"
)

type testCase struct {
	expectValue int64
	input       string
}


func testHelper(tests []testCase, t *testing.T) {
	for _, tt := range tests {
		p := NewParser(tt.input)
		val := p.Expr()

		if value, ok := val.(int64); ok {
			if value != tt.expectValue {
				t.Fatalf("%s expected %d, actual %d", tt.input, tt.expectValue, value)
			}
		} else {
			t.Fatalf("could't not get invalid type, %s ", tt.input)
		}
	}
}

var addOpTests = []testCase{
	{1, "1 + 0"},
	{5, "2 + 3"},
	{7, "2 + 3 + 2"},
	{-1, "-2 + 3 + -2"},
}

// TestAddOperation is testing add operation
func TestAddOperation(t *testing.T) {
    testHelper(addOpTests, t)
}

var minusOpTests = []testCase {
	{-1, "0 - 1"},
	{-2, "0 - 1 - 1"},
}

func TestMinusOperation(t *testing.T) {
    testHelper(addOpTests, t)
}
