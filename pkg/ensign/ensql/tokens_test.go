package ensql_test

import (
	"testing"

	. "github.com/rotationalio/ensign/pkg/ensign/ensql"
	"github.com/stretchr/testify/require"
)

func TestReservedWordsTyped(t *testing.T) {
	for _, word := range ReservedWords {
		_, ok := ReservedWordType[word]
		require.True(t, ok, "the reserved word %q is not typed", word)
	}
}

func TestReservedWordTokenization(t *testing.T) {
	// Test that different combinations of reserved words are tokenized correctly and
	// that the parsing of reserved words is case and whitespace insensitive.
	expected := []Token{
		{SELECT, ReservedWord, len(SELECT)},
		{FROM, ReservedWord, len(FROM)},
		{WHERE, ReservedWord, len(WHERE)},
		{AS, ReservedWord, len(AS)},
		{OFFSET, ReservedWord, len(OFFSET)},
		{LIMIT, ReservedWord, len(LIMIT)},
		{EQ, OperatorToken, len(EQ)},
		{NE, OperatorToken, len(NE)},
		{GT, OperatorToken, len(GT)},
		{LT, OperatorToken, len(LT)},
		{GTE, OperatorToken, len(GTE)},
		{LTE, OperatorToken, len(LTE)},
		{AND, OperatorToken, len(AND)},
		{OR, OperatorToken, len(OR)},
		{LIKE, OperatorToken, len(LIKE)},
		{ILIKE, OperatorToken, len(ILIKE)},
		{ASTERISK, Asterisk, 1},
		{COMMA, Punctuation, 1},
		{DOT, Punctuation, 1},
		{LP, Punctuation, 1},
		{RP, Punctuation, 1},
		{SC, Punctuation, 1},
	}

	testCases := []struct {
		sql string
		msg string
	}{
		{
			"SELECT FROM WHERE AS OFFSET LIMIT = != > < >= <= AND OR LIKE ILIKE * , . ( ) ;",
			"simple tokenization with spaces",
		},
		{
			"SELECTFROMWHEREASOFFSETLIMIT=!=><>=<=ANDORLIKEILIKE*,.();",
			"no whitespace at all",
		},
		{
			"select from where as offset limit = != > < >= <= and or like ilike * , . ( ) ;",
			"all lowercase reserved words",
		},
		{
			"Select From Where As Offset Limit = != > < >= <= And Or Like ILike * , . ( ) ;",
			"title casing reserved words",
		},
		{
			"SELECT  FROM      WHERE\t AS \tOFFSET\n\n LIMIT\r\n =  !=\t\t\t\t > \t  < \n\t >=\t \n <= AND \r\n  OR\r LIKE  \t ILIKE     * , . ( )\t\t   ;\n\n",
			"crazy whitespace",
		},
	}

	for _, tc := range testCases {
		actual := Tokenize(tc.sql)
		require.Equal(t, expected, actual, tc.msg)
	}
}

func TestQuotedStringTokenization(t *testing.T) {
	testCases := []struct {
		sql      string
		expected Token
		msg      string
	}{
		{"'foo'", Token{"foo", QuotedString, 5}, "regular quoted string"},
		{`'foo\'s'`, Token{`foo\'s`, QuotedString, 8}, "escaped quoted string"},
		{"'foo", Token{"", EmptyToken, 4}, "unclosed quote"},
		{`'foo\'s`, Token{"", EmptyToken, 7}, "unclosed, escaped quote"},
	}

	for _, tc := range testCases {
		tokens := Tokenize(tc.sql)
		require.Len(t, tokens, 1, tc.msg)
		require.Equal(t, tc.expected, tokens[0], tc.msg)
	}
}
