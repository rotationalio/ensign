package ensql_test

import (
	"fmt"
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

func TestNumericTokenization(t *testing.T) {
	testCases := []struct {
		sql      string
		expected Token
		msg      string
	}{
		{"42", Token{"42", Numeric, 2}, "integer numeric"},
		{"192.321", Token{"192.321", Numeric, 7}, "floating point numeric"},
		{"-7", Token{"-7", Numeric, 2}, "negative integer"},
		{"-0.83", Token{"-0.83", Numeric, 5}, "negative floating point"},
	}

	for _, tc := range testCases {
		tokens := Tokenize(tc.sql)
		require.Len(t, tokens, 1, tc.msg)
		require.Equal(t, tc.expected, tokens[0], tc.msg)
	}
}

func TestIdentifierTokenization(t *testing.T) {
	testCases := []struct {
		sql      string
		expected Token
		msg      string
	}{
		{"foo", Token{"foo", Identifier, 3}, "regular identifier"},
		{"*", Token{"*", Asterisk, 1}, "asterisk identifier"},
		{"snake_case", Token{"snake_case", Identifier, 10}, "identifier with underscores"},
		{"CamelCase", Token{"CamelCase", Identifier, 9}, "identifier with uppercase"},
		{"blue42blue42", Token{"blue42blue42", Identifier, 12}, "identifier with digits"},
		{"_private", Token{"_private", Identifier, 8}, "identifier started by underscore"},
		{"true_light", Token{"true_light", Identifier, 10}, "boolean prefix"},
	}

	for _, tc := range testCases {
		tokens := Tokenize(tc.sql)
		require.Len(t, tokens, 1, tc.msg)
		require.Equal(t, tc.expected, tokens[0], tc.msg)
	}
}

func TestBooleanTokenization(t *testing.T) {
	testCases := []struct {
		sql      string
		expected Token
		msg      string
	}{
		{"f", Token{"f", Boolean, 1}, "lowercase f"},
		{"F", Token{"F", Boolean, 1}, "uppercase F"},
		{"False", Token{"False", Boolean, 5}, "title case False"},
		{"FALSE", Token{"FALSE", Boolean, 5}, "uppercase FALSE"},
		{"false", Token{"false", Boolean, 5}, "lowercase false"},
		{"FaLsE", Token{"FaLsE", Identifier, 5}, "unparseable FaLsE"},
		{"t", Token{"t", Boolean, 1}, "lowercase t"},
		{"T", Token{"T", Boolean, 1}, "uppercase T"},
		{"True", Token{"True", Boolean, 4}, "title case True"},
		{"TRUE", Token{"TRUE", Boolean, 4}, "uppercase TRUE"},
		{"true", Token{"true", Boolean, 4}, "lowercase true"},
		{"TrUe", Token{"TrUe", Identifier, 4}, "unparseable TrUe"},
	}

	for _, tc := range testCases {
		tokens := Tokenize(tc.sql)
		require.Len(t, tokens, 1, tc.msg)
		require.Equal(t, tc.expected, tokens[0], tc.msg)
	}
}

func TestTokenize(t *testing.T) {
	sql := `SELECT identifier FROM table_identifier WHERE 'quoted string with spaces' = -32.31 AND * ilike 41; '-31.31' 'foo\'s' TRUE f`

	expected := []Token{
		{SELECT, ReservedWord, 6},
		{"identifier", Identifier, len("identifier")},
		{FROM, ReservedWord, 4},
		{"table_identifier", Identifier, len("table_identifier")},
		{WHERE, ReservedWord, 5},
		{"quoted string with spaces", QuotedString, len("quoted string with spaces") + 2},
		{EQ, OperatorToken, 1},
		{"-32.31", Numeric, 6},
		{AND, OperatorToken, 3},
		{"*", Asterisk, 1},
		{ILIKE, OperatorToken, 5},
		{"41", Numeric, 2},
		{SC, Punctuation, 1},
		{"-31.31", QuotedString, 8},
		{"foo\\'s", QuotedString, 8},
		{"TRUE", Boolean, 4},
		{"f", Boolean, 1},
	}

	for i, actual := range Tokenize(sql) {
		require.Equal(t, expected[i], actual)
	}
}

func TestTokenizeSQL(t *testing.T) {
	sql := `
SELECT name, age, favorite_color AS color, title AS profession, salary
	FROM hiring.employee.8
	WHERE company = 'rotational' AND salary < 250000 OR specialist = true
OFFSET 2300
LIMIT 100;`

	expected := []Token{
		{SELECT, ReservedWord, 6},
		{"name", Identifier, 4},
		{COMMA, Punctuation, 1},
		{"age", Identifier, 3},
		{COMMA, Punctuation, 1},
		{"favorite_color", Identifier, 14},
		{AS, ReservedWord, 2},
		{"color", Identifier, 5},
		{COMMA, Punctuation, 1},
		{"title", Identifier, 5},
		{AS, ReservedWord, 2},
		{"profession", Identifier, 10},
		{COMMA, Punctuation, 1},
		{"salary", Identifier, 6},
		{FROM, ReservedWord, 4},
		{"hiring", Identifier, 6},
		{DOT, Punctuation, 1},
		{"employee", Identifier, 8},
		{DOT, Punctuation, 1},
		{"8", Numeric, 1},
		{WHERE, ReservedWord, 5},
		{"company", Identifier, 7},
		{EQ, OperatorToken, 1},
		{"rotational", QuotedString, 12},
		{AND, OperatorToken, 3},
		{"salary", Identifier, 6},
		{LT, OperatorToken, 1},
		{"250000", Numeric, 6},
		{OR, OperatorToken, 2},
		{"specialist", Identifier, 10},
		{EQ, OperatorToken, 1},
		{"true", Boolean, 4},
		{OFFSET, ReservedWord, 6},
		{"2300", Numeric, 4},
		{LIMIT, ReservedWord, 5},
		{"100", Numeric, 3},
		{SC, Punctuation, 1},
	}

	for i, actual := range Tokenize(sql) {
		require.Equal(t, expected[i], actual)
	}
}

func TestParseInt(t *testing.T) {
	testCases := []struct {
		token    Token
		expected int64
		err      error
	}{
		{Token{"1234", Numeric, 4}, 1234, nil},
		{Token{"0", Numeric, 1}, 0, nil},
		{Token{"-42", Numeric, 3}, -42, nil},
		{Token{"abc", Numeric, 3}, 0, fmt.Errorf("strconv.ParseInt: parsing %q: invalid syntax", "abc")},
		{Token{"abc", QuotedString, 5}, 0, ErrNonNumeric},
	}

	for _, tc := range testCases {
		actual, err := tc.token.ParseInt(10, 64)
		require.Equal(t, tc.expected, actual)

		if tc.err != nil {
			require.Error(t, err)
			require.EqualError(t, err, tc.err.Error())
		} else {
			require.NoError(t, err)
		}
	}
}

func TestParseUint(t *testing.T) {
	testCases := []struct {
		token    Token
		expected uint64
		err      error
	}{
		{Token{"1234", Numeric, 4}, 1234, nil},
		{Token{"0", Numeric, 1}, 0, nil},
		{Token{"-42", Numeric, 3}, 0, fmt.Errorf("strconv.ParseUint: parsing %q: invalid syntax", "-42")},
		{Token{"abc", Numeric, 3}, 0, fmt.Errorf("strconv.ParseUint: parsing %q: invalid syntax", "abc")},
		{Token{"abc", QuotedString, 5}, 0, ErrNonNumeric},
	}

	for _, tc := range testCases {
		actual, err := tc.token.ParseUint(10, 64)
		require.Equal(t, tc.expected, actual)

		if tc.err != nil {
			require.Error(t, err)
			require.EqualError(t, err, tc.err.Error())
		} else {
			require.NoError(t, err)
		}
	}
}

func TestParseFloat(t *testing.T) {
	testCases := []struct {
		token    Token
		expected float64
		err      error
	}{
		{Token{"12.34", Numeric, 4}, 12.34, nil},
		{Token{"0", Numeric, 1}, 0.0, nil},
		{Token{"-4.23333212", Numeric, 3}, -4.23333212, nil},
		{Token{"abc", Numeric, 3}, 0, fmt.Errorf("strconv.ParseFloat: parsing %q: invalid syntax", "abc")},
		{Token{"abc", QuotedString, 5}, 0, ErrNonNumeric},
	}

	for _, tc := range testCases {
		actual, err := tc.token.ParseFloat(64)
		require.Equal(t, tc.expected, actual)

		if tc.err != nil {
			require.Error(t, err)
			require.EqualError(t, err, tc.err.Error())
		} else {
			require.NoError(t, err)
		}
	}
}

func TestParseBool(t *testing.T) {
	testCases := []struct {
		token    Token
		expected bool
		err      error
	}{
		{Token{"t", Boolean, 1}, true, nil},
		{Token{"T", Boolean, 1}, true, nil},
		{Token{"1", Boolean, 1}, true, nil},
		{Token{"true", Boolean, 4}, true, nil},
		{Token{"TRUE", Boolean, 4}, true, nil},
		{Token{"True", Boolean, 4}, true, nil},
		{Token{"0", Boolean, 1}, false, nil},
		{Token{"f", Boolean, 1}, false, nil},
		{Token{"F", Boolean, 1}, false, nil},
		{Token{"false", Boolean, 5}, false, nil},
		{Token{"False", Boolean, 5}, false, nil},
		{Token{"FALSE", Boolean, 5}, false, nil},
		{Token{"abc", Boolean, 3}, false, fmt.Errorf("strconv.ParseBool: parsing %q: invalid syntax", "abc")},
		{Token{"abc", QuotedString, 5}, false, ErrNonBoolean},
	}

	for _, tc := range testCases {
		actual, err := tc.token.ParseBool()
		require.Equal(t, tc.expected, actual)

		if tc.err != nil {
			require.Error(t, err)
			require.EqualError(t, err, tc.err.Error())
		} else {
			require.NoError(t, err)
		}
	}
}
