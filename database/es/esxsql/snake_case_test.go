package esxsql

import (
	"testing"
)

func TestSnakeCase(t *testing.T) {
	tcs := []struct{ expected, in string }{
		{
			expected: "checking_the_snake_case",
			in:       "Checking The Snake Case",
		},
		{
			expected: "snake_case_is_nice",
			in:       "SnakeCaseIsNice",
		},
	}
	for _, tc := range tcs {
		res := ToSnakeCase(tc.in)
		if res != tc.expected {
			t.Errorf("to snake case failed - expected: '%s' is: '%s'", tc.expected, res)
		}
	}
}
