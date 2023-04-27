package ensql

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
)

var ReservedWords = []string{
	SELECT, FROM, WHERE, AS, OFFSET, LIMIT,
	EQ, NE, GT, LT, GTE, LTE, AND, OR, LIKE, ILIKE,
	ASTERISK, COMMA, DOT, LP, RP, SC,
}

type Token struct {
	Token string
	Type  TokenType
}

type TokenType uint8

const (
	UnknownTokenType TokenType = iota
	ReservedWord
	Identifier
	Asterisk
	OperatorToken
	QuotedString
	Numeric
)
