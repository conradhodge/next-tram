package matcher_test

import (
	"fmt"
	"testing"

	"github.com/conradhodge/next-tram/src/lib/matcher"
)

func TestMatches(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected bool
	}{
		{
			input:    "abc",
			expected: false,
		},
		{
			input:    "123",
			expected: false,
		},
		{
			input:    "abc123",
			expected: false,
		},
		{
			input:    "0187582-820a-41fa-9f45-9bdcefce0791",
			expected: false,
		},
		{
			input:    "0c187582-820-41fa-9f45-9bdcefce0791",
			expected: false,
		},
		{
			input:    "0c187582-820a-41f-9f45-9bdcefce0791",
			expected: false,
		},
		{
			input:    "0c187582-820a-41fa-9f4-9bdcefce0791",
			expected: false,
		},
		{
			input:    "0c187582-820a-41fa-9f45-9bdcefce079",
			expected: false,
		},
		{
			input:    "0c187582-820a-41fa-9f45-9bdcefce0791",
			expected: true,
		},
	}

	for _, test := range tests {
		name := fmt.Sprintf("%v", test.input)
		t.Run(name, func(t *testing.T) {
			guidMatcher := matcher.IsGUID()

			if guidMatcher.Matches(test.input) != test.expected {
				t.Fatalf("Expected %v", test.expected)
			}
		})
	}
}

func TestString(t *testing.T) {
	guidMatcher := matcher.IsGUID()

	if guidMatcher.String() != "is a GUID" {
		t.Fatal("Expected string: is a GUID")
	}
}
