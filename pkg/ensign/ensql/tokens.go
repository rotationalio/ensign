package ensql

import (
	"strconv"
	"strings"
)

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
	NEALT    = "<>"
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
	EQ, NE, NEALT, GTE, LTE, GT, LT, AND, OR, LIKE, ILIKE,
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
	NEALT:    OperatorToken,
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

// Operator fields for where clauses and conditional queries
type Operator uint8

const (
	UnknownOperator Operator = iota
	Eq                       // =
	Ne                       // !=
	Gt                       // >
	Lt                       // <
	Gte                      // >=
	Lte                      // <=
	Like                     // like
	ILike                    // ilike
	And                      // AND
	Or                       // OR
)

func (o Operator) String() string {
	switch o {
	case Eq:
		return EQ
	case Ne:
		return NE
	case Gt:
		return GT
	case Lt:
		return LT
	case Gte:
		return GTE
	case Lte:
		return LTE
	case Like:
		return LIKE
	case ILike:
		return ILIKE
	case And:
		return AND
	case Or:
		return OR
	default:
		return "UnknownOperator"
	}
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
	Boolean
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

// Parse a numeric token as a signed integer using strconv.ParseInt. Generally, the base
// should be 10 and the bitSize should be 64 unless otherwise defined by the schema.
func (t Token) ParseInt(base, bitSize int) (int64, error) {
	if t.Type != Numeric {
		return 0, ErrNonNumeric
	}
	return strconv.ParseInt(t.Token, base, bitSize)
}

// Parse a numeric token as an unsigned integer using strconv.ParseUint. Generally, the
// base should be 10 and the bitSize should be 64 unless otherwise defined by the schema.
func (t Token) ParseUint(base, bitSize int) (uint64, error) {
	if t.Type != Numeric {
		return 0, ErrNonNumeric
	}
	return strconv.ParseUint(t.Token, base, bitSize)
}

// Parse a numeric token as a float using strconv.ParseFloat. Generally, the bitSize
// should be 64 (e.g. double) unless otherwise defined by the schema.
func (t Token) ParseFloat(bitSize int) (float64, error) {
	if t.Type != Numeric {
		return 0, ErrNonNumeric
	}
	return strconv.ParseFloat(t.Token, bitSize)
}

// Parse a boolean token as a bool using strconv.ParseBool.
func (t Token) ParseBool() (bool, error) {
	if t.Type != Boolean {
		return false, ErrNonBoolean
	}
	return strconv.ParseBool(t.Token)
}
