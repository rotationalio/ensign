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
		{
			Name:     "invalid state transition from",
			SQL:      "SELECT * FROM topic 1234",
			Expected: nil,
			Err:      Error(20, "1234", "invalid from clause").Error(),
		},
		{
			Name:     "invalid state transition from select",
			SQL:      "SELECT * FROM topic SELECT",
			Expected: nil,
			Err:      Error(20, "SELECT", "invalid from clause").Error(),
		},
		{
			Name:     "with offset",
			SQL:      "SELECT * FROM topic OFFSET 42",
			Expected: &Query{Type: SelectQuery, Fields: []Token{{"*", Asterisk, 1}}, Topic: Topic{Topic: "topic"}, Offset: 42, HasOffset: true},
			Err:      "",
		},
		{
			Name:     "with offset terminated",
			SQL:      "SELECT * FROM topic OFFSET 42;",
			Expected: &Query{Type: SelectQuery, Fields: []Token{{"*", Asterisk, 1}}, Topic: Topic{Topic: "topic"}, Offset: 42, HasOffset: true},
			Err:      "",
		},
		{
			Name:     "with limit",
			SQL:      "SELECT * FROM topic LIMIT 42",
			Expected: &Query{Type: SelectQuery, Fields: []Token{{"*", Asterisk, 1}}, Topic: Topic{Topic: "topic"}, Limit: 42, HasLimit: true},
			Err:      "",
		},
		{
			Name:     "with limit terminated",
			SQL:      "SELECT * FROM topic LIMIT 42;",
			Expected: &Query{Type: SelectQuery, Fields: []Token{{"*", Asterisk, 1}}, Topic: Topic{Topic: "topic"}, Limit: 42, HasLimit: true},
			Err:      "",
		},
		{
			Name:     "with offset and limit",
			SQL:      "SELECT * FROM topic OFFSET 23 LIMIT 42",
			Expected: &Query{Type: SelectQuery, Fields: []Token{{"*", Asterisk, 1}}, Topic: Topic{Topic: "topic"}, Limit: 42, HasLimit: true, Offset: 23, HasOffset: true},
			Err:      "",
		},
		{
			Name:     "can specify 0 offset and limit",
			SQL:      "SELECT * FROM topic OFFSET 0 LIMIT 0",
			Expected: &Query{Type: SelectQuery, Fields: []Token{{"*", Asterisk, 1}}, Topic: Topic{Topic: "topic"}, Limit: 0, HasLimit: true, Offset: 0, HasOffset: true},
			Err:      "",
		},
		{
			Name:     "with offset and limit terminated",
			SQL:      "SELECT * FROM topic OFFSET 23 LIMIT 42;",
			Expected: &Query{Type: SelectQuery, Fields: []Token{{"*", Asterisk, 1}}, Topic: Topic{Topic: "topic"}, Limit: 42, HasLimit: true, Offset: 23, HasOffset: true},
			Err:      "",
		},
		{
			Name:     "with offset and limit reversed",
			SQL:      "SELECT * FROM topic LIMIT 42 OFFSET 23",
			Expected: &Query{Type: SelectQuery, Fields: []Token{{"*", Asterisk, 1}}, Topic: Topic{Topic: "topic"}, Limit: 42, HasLimit: true, Offset: 23, HasOffset: true},
			Err:      "",
		},
		{
			Name:     "with offset and limit reversed terminated",
			SQL:      "SELECT * FROM topic LIMIT 42 OFFSET 23;",
			Expected: &Query{Type: SelectQuery, Fields: []Token{{"*", Asterisk, 1}}, Topic: Topic{Topic: "topic"}, Limit: 42, HasLimit: true, Offset: 23, HasOffset: true},
			Err:      "",
		},
		{
			Name:     "invalid state transition offset",
			SQL:      "SELECT * FROM topic OFFSET 32 WHERE",
			Expected: nil,
			Err:      Error(30, "WHERE", "invalid offset clause").Error(),
		},
		{
			Name:     "invalid state transition limit",
			SQL:      "SELECT * FROM topic LIMIT 100 WHERE",
			Expected: nil,
			Err:      Error(30, "WHERE", "invalid limit clause").Error(),
		},
		{
			Name:     "invalid offset",
			SQL:      "SELECT * FROM topic OFFSET abcd",
			Expected: nil,
			Err:      Error(27, "abcd", "invalid offset").Error(),
		},
		{
			Name:     "invalid limit",
			SQL:      "SELECT * FROM topic LIMIT abcd",
			Expected: nil,
			Err:      Error(26, "abcd", "invalid limit").Error(),
		},
		{
			Name:     "invalid offset quoted",
			SQL:      "SELECT * FROM topic OFFSET '32'",
			Expected: nil,
			Err:      Error(27, "32", "invalid offset").Error(),
		},
		{
			Name:     "invalid limit quoted",
			SQL:      "SELECT * FROM topic LIMIT '100'",
			Expected: nil,
			Err:      Error(26, "100", "invalid limit").Error(),
		},
		{
			Name:     "invalid offset negative",
			SQL:      "SELECT * FROM topic OFFSET -42",
			Expected: nil,
			Err:      Error(27, "-42", "could not parse offset").Error(),
		},
		{
			Name:     "invalid limit negative",
			SQL:      "SELECT * FROM topic LIMIT -1000",
			Expected: nil,
			Err:      Error(26, "-1000", "could not parse limit").Error(),
		},
		{
			Name:     "invalid offset float",
			SQL:      "SELECT * FROM topic OFFSET 3.24",
			Expected: nil,
			Err:      Error(27, "3.24", "could not parse offset").Error(),
		},
		{
			Name:     "invalid limit float",
			SQL:      "SELECT * FROM topic LIMIT 7.77",
			Expected: nil,
			Err:      Error(26, "7.77", "could not parse limit").Error(),
		},
		{
			Name:     "canot duplicate offset",
			SQL:      "SELECT * FROM topic OFFSET 5 LIMIT 3 OFFSET 6",
			Expected: nil,
			Err:      Error(44, "6", "offset has already been set").Error(),
		},
		{
			Name:     "canot duplicate limit",
			SQL:      "SELECT * FROM topic LIMIT 5 OFFSET 3 LIMIT 6",
			Expected: nil,
			Err:      Error(43, "6", "limit has already been set").Error(),
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
