package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    " mayonnaise is good ",
			expected: []string{"mayonnaise", "is", "good"},
		},
		{
			input:    " bUt KeTcHuP is my FAVOURITE ",
			expected: []string{"but", "ketchup", "is", "my", "favourite"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(c.expected) != len(actual) {
			t.Errorf("Expected length %d, got %d", len(c.expected), len(actual))
		}
		for i := range actual {
			if word := actual[i]; word != c.expected[i] {
				t.Errorf("Expected word %q, got %q", c.expected[i], word)
			}
		}
	}
}
