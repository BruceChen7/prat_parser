package parser

import (
	"testing"
)

var addOpTests = []struct {
	expectValue int64
	input       string
}{
	{1, "1 + 0"},
	{5, "2 + 3"},
	{5, "+2 + 3"},
	{7, "+3 + +4"},
}

// TestAddOperation is testing add operation
func TestAddOperation(t *testing.T) {
	for _, tt := range addOpTests {
		p := NewParser(tt.input)
		val := p.Parse(0)

		if value, ok := val.(int64); ok {
			if value != tt.expectValue {
				t.Fatalf("expected %d, actual %d", tt.expectValue, value)
			}
		} else {
			t.Fatalf("could't not get invalid type, %s ", tt.input)
		}
	}

}
