package ensql_test

import (
	"testing"

	"github.com/rotationalio/ensign/pkg/ensign/ensql"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	Name     string       // description of the test case
	SQL      string       // input SQL to be parsed
	Expected *ensql.Query // expected resulting AST that represents the query
	Err      string       // expected error result from parsing
}

func TestParse(t *testing.T) {
	ts := []testCase{
		{
			Name:     "empty query is invalid",
			SQL:      "",
			Expected: nil,
			Err:      "empty query is invalid",
		},
		{
			Name:     "SELECT without FROM is invalid",
			SQL:      "SELECT",
			Expected: nil,
			Err:      "topic name cannot be empty",
		},
		{
			Name:     "Unclosed Quote",
			SQL:      "SELECT 'unclosed",
			Expected: nil,
			Err:      "syntax error at position 7 near \"'uncl\": quoted string missing closing quote",
		},
	}

	for _, tc := range ts {
		t.Run(tc.Name, func(t *testing.T) {
			actual, err := ensql.Parse(tc.SQL)

			// Expect an error if the test case has an error.
			if tc.Err != "" {
				require.Error(t, err, "expected an error to have occurred")
				require.EqualError(t, err, tc.Err, "unexpected error occurred")
			}

			// Expect a query tree if the expected value is not nil.
			if tc.Expected != nil {
				require.NotNil(t, actual, "expected a non-nil query tree")
				require.Equal(t, tc.Expected, &actual, "actual query tree did not match test case expectation")
			}
		})
	}

}
