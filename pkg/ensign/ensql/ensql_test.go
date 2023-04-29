package ensql_test

import (
	"testing"

	. "github.com/rotationalio/ensign/pkg/ensign/ensql"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	Name     string // description of the test case
	SQL      string // input SQL to be parsed
	Expected *Query // expected resulting AST that represents the query
	Err      string // expected error result from parsing
}

func TestParse(t *testing.T) {
	ts := []testCase{
		{
			Name:     "empty query is invalid",
			SQL:      "",
			Expected: nil,
			Err:      ErrEmptyQuery.Error(),
		},
		{
			Name:     "SELECT without FROM is invalid",
			SQL:      "SELECT",
			Expected: nil,
			Err:      ErrMissingTopic.Error(),
		},
		{
			Name:     "invalid query",
			SQL:      "WHERE",
			Expected: nil,
			Err:      Error(0, "WHERE", "invalid query type").Error(),
		},
		{
			Name:     "must start with reserved word",
			SQL:      "foo",
			Expected: nil,
			Err:      Error(0, "foo", "invalid query type").Error(),
		},
		{
			Name:     "invalid field identifier numeric",
			SQL:      "SELECT 1234 FROM topic",
			Expected: nil,
			Err:      Error(7, "1234", "invalid field identifier").Error(),
		},
		{
			Name:     "invalid field identifier quoted",
			SQL:      "SELECT 'name' FROM topic",
			Expected: nil,
			Err:      Error(7, "name", "invalid field identifier").Error(),
		},
		{
			Name:     "must specify select fields",
			SQL:      "SELECT FROM topic",
			Expected: nil,
			Err:      ErrNoFieldsSelected.Error(),
		},
		{
			Name:     "select all fields",
			SQL:      "SELECT * FROM topic",
			Expected: &Query{Type: SelectQuery, Fields: []Token{{"*", Asterisk, 1}}, Topic: Topic{Topic: "topic"}},
			Err:      "",
		},
		{
			Name:     "select all fields with termination",
			SQL:      "SELECT * FROM topic;",
			Expected: &Query{Type: SelectQuery, Fields: []Token{{"*", Asterisk, 1}}, Topic: Topic{Topic: "topic"}},
			Err:      "",
		},
		{
			Name:     "invalid asterisk with field",
			SQL:      "SELECT *, name FROM topic",
			Expected: nil,
			Err:      ErrInvalidSelectAllFields.Error(),
		},
		{
			Name:     "invalid double asterisk",
			SQL:      "SELECT *, * FROM topic",
			Expected: nil,
			Err:      ErrInvalidSelectAllFields.Error(),
		},
		{
			Name:     "select single field",
			SQL:      "SELECT name FROM topic",
			Expected: &Query{Type: SelectQuery, Fields: []Token{{"name", Identifier, 4}}, Topic: Topic{Topic: "topic"}},
			Err:      "",
		},
		{
			Name:     "select single field with alias",
			SQL:      "SELECT name AS first_name FROM topic",
			Expected: &Query{Type: SelectQuery, Fields: []Token{{"name", Identifier, 4}}, Topic: Topic{Topic: "topic"}, Aliases: map[string]string{"name": "first_name"}},
			Err:      "",
		},
		{
			Name:     "select multiple fields",
			SQL:      "SELECT name, age, color FROM topic",
			Expected: &Query{Type: SelectQuery, Fields: []Token{{"name", Identifier, 4}, {"age", Identifier, 3}, {"color", Identifier, 5}}, Topic: Topic{Topic: "topic"}},
			Err:      "",
		},
		{
			Name:     "select multiple fields with aliases",
			SQL:      "SELECT name AS first_name, age, color AS favorite_color FROM topic",
			Expected: &Query{Type: SelectQuery, Fields: []Token{{"name", Identifier, 4}, {"age", Identifier, 3}, {"color", Identifier, 5}}, Topic: Topic{Topic: "topic"}, Aliases: map[string]string{"name": "first_name", "color": "favorite_color"}},
			Err:      "",
		},
		{
			Name:     "missing comma fields",
			SQL:      "SELECT name age color FROM topic",
			Expected: nil,
			Err:      Error(12, "age", "invalid select fields statement").Error(),
		},
		{
			Name:     "invalid alias",
			SQL:      "SELECT name AS 1234 FROM topic",
			Expected: nil,
			Err:      Error(15, "1234", "invalid alias identifier").Error(),
		},
		{
			Name:     "cannot alias *",
			SQL:      "SELECT * AS foo FROM topic",
			Expected: nil,
			Err:      Error(12, "foo", "cannot alias *").Error(),
		},
		{
			Name:     "missing comma alias",
			SQL:      "SELECT name AS first_name age color AS favorite_color FROM topic",
			Expected: nil,
			Err:      Error(26, "age", "invalid select fields statement").Error(),
		},
		{
			Name:     "invalid topic identifier numeric",
			SQL:      "SELECT * FROM 1234",
			Expected: nil,
			Err:      Error(14, "1234", "invalid topic identifier").Error(),
		},
		{
			Name:     "invalid topic identifier quoted",
			SQL:      "SELECT * FROM 'topic'",
			Expected: nil,
			Err:      Error(14, "topic", "invalid topic identifier").Error(),
		},
		// {
		// 	Name:     "unclosed quote",
		// 	SQL:      "SELECT 'unclosed",
		// 	Expected: nil,
		// 	Err:      "syntax error at position 7 near \"'uncl\": quoted string missing closing quote",
		// },
	}

	for _, tc := range ts {
		t.Run(tc.Name, func(t *testing.T) {
			actual, err := Parse(tc.SQL)

			// Expect an error if the test case has an error.
			if tc.Err != "" {
				require.Error(t, err, "expected an error to have occurred")
				require.EqualError(t, err, tc.Err, "unexpected error occurred")
			}

			// Expect a query tree if the expected value is not nil.
			if tc.Expected != nil {
				// Add the raw field to make it easier to compose tests cases
				tc.Expected.Raw = tc.SQL

				require.NotNil(t, actual, "expected a non-nil query tree")
				require.Equal(t, tc.Expected, &actual, "actual query tree did not match test case expectation")
			}
		})
	}

}
