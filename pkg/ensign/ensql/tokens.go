package ensql

import "strings"

// Reserved Words constants
const (
	SELECT   = "SELECT"
	FROM     = "FROM"
	WHERE    = "WHERE"
	AS       = "AS"
	OFFSET   = "OFFSET"
	LIMIT    = "LIMIT"
	EQ       = "="
	NE       = "!="
	GT       = ">"
	LT       = "<"
	GTE      = ">="
	LTE      = "<="
	AND      = "AND"
	OR       = "OR"
	LIKE     = "LIKE"
	ILIKE    = "ILIKE"
	ASTERISK = "*"
	COMMA    = ","
	DOT      = "."
	LP       = "("
	RP       = ")"
	SC       = ";"
	SQUOTE   = '\''
	MINUS    = '-'
	ESCAPE   = '\\'
)

var (
	Empty = Token{"", EmptyToken, 0}
)

// NOTE: GT and LT must follow GTE and LTE in this list (or in general, any word that
// is a prefix of another word must follow that word to ensure parsing is correct).
var ReservedWords = []string{
	SELECT, FROM, WHERE, AS, OFFSET, LIMIT,
	EQ, NE, GTE, LTE, GT, LT, AND, OR, LIKE, ILIKE,
	ASTERISK, COMMA, DOT, LP, RP, SC,
}

var ReservedWordType = map[string]TokenType{
	SELECT:   ReservedWord,
	FROM:     ReservedWord,
	WHERE:    ReservedWord,
	AS:       ReservedWord,
	OFFSET:   ReservedWord,
	LIMIT:    ReservedWord,
	EQ:       OperatorToken,
	NE:       OperatorToken,
	GT:       OperatorToken,
	LT:       OperatorToken,
	GTE:      OperatorToken,
	LTE:      OperatorToken,
	AND:      OperatorToken,
	OR:       OperatorToken,
	LIKE:     OperatorToken,
	ILIKE:    OperatorToken,
	ASTERISK: Asterisk,
	COMMA:    Punctuation,
	DOT:      Punctuation,
	LP:       Punctuation,
	RP:       Punctuation,
	SC:       Punctuation,
}

// A token represents a parsed element from the SQL and is returned from peek. The
// token string may not match the original string in the query (for example it might be
// uppercased or have quotations or whitespace stripped). When evaluating tokens, both
// the token type and the token itself should be used in concert to ensure correct
// normalization has occurred. The length is used to advance the parser index and may
// not match the length of the parsed token string.
type Token struct {
	Token  string
	Type   TokenType
	Length int
}

type TokenType uint8

const (
	UnknownTokenType TokenType = iota
	EmptyToken
	ReservedWord
	OperatorToken
	Asterisk
	Punctuation
	Identifier
	QuotedString
	Numeric
)

// Tokenize returns the tokens parsed from the input string with no validation or FSM.
// This function is primarily used by tests but can also be used by debugging tools to
// determine how a SQL query is being parsed.
func Tokenize(sql string) []Token {
	parser := &parser{sql: strings.TrimSpace(sql), idx: 0, step: stepInit}
	tokens := make([]Token, 0)

	for parser.idx < len(parser.sql) {
		token := parser.pop()
		tokens = append(tokens, token)
	}

	return tokens
}
