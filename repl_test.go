package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " Hello WORLD ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "HELLO",
			expected: []string{"hello"},
		},
	}

	// loop over cases
	// for each case check actual function call with input
	// compare each word in the actual output with expected
	// if there are differences, raise error
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) == 0 {
			t.Errorf("resulting slice is empty")
		}
		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("actual value does not match expected value")
			}
		}
	}
}
