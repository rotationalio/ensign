package ensql

type Token struct {
	Token string
	Type  TokenType
}

type TokenType uint8

const (
	UnknownTokenType TokenType = iota
	ReservedWord
	QuotedString
	Identifier
)
