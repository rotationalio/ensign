package ensql

import (
	"regexp"
	"strings"
	"unicode"
)

// Parse an EnSQL statement to create a Query object for an Ensign SQL execution. An
// error is returned on syntax or validation errors that occur during parsing.
func Parse(sql string) (Query, error) {
	// Remove any space before or after the query string.
	sql = strings.TrimSpace(sql)

	// Create the parser object initialized and ready to parse.
	// NOTE: the parser is not reusable and must be allocated for each parse.
	parser := &parser{
		sql:  sql,
		idx:  0,
		step: stepInit,
		query: Query{
			Raw: sql,
		},
	}

	// Execute the parse and return any errors.
	return parser.parse()
}

// Parser implements a feed-forward SQL parsing mechanism.
type parser struct {
	sql   string
	idx   int
	step  step
	query Query
	err   error
}

// Parse executes the parse but ensures that the parse isn't executed a second time if
// it has already been executed by saving any parse errors locally. Exec advances the
// state of the of the parser so if the state isn't init then exec will not be called
// a second time.
func (p *parser) parse() (Query, error) {
	if p.step == stepInit {
		if p.err = p.exec(); p.err == nil {
			p.err = p.validate()
		}
	}
	return p.query, p.err
}

// Exec implements a feed-forward parser, advancing the index and checking the current
// step to determine how to parse the next section of the SQL string. Parsing stops
// when the end of the string is reached or the parser reaches a state where it cannot
// continue parsing using the FSM described by the SQL statement.
func (p *parser) exec() error {
	// Continue until we reach the end of the string.
	// NOTE: p.pop() must be called to advance the index and guarantee termination but
	// because of look-aheads, pop is called in each case rather than in the loop.
	for p.idx < len(p.sql) {
		switch p.step {
		case stepInit:
			// At the initial step we expect a query determiner such as SELECT or WITH
			// This means that the very first token should be a reserved word.
			token := p.peek()
			if token.Type != ReservedWord {
				return Error(p.idx, token.Token, "invalid query type")
			}

			switch token.Token {
			case SELECT:
				p.query.Type = SelectQuery
				p.step = stepSelectField
			default:
				return Error(p.idx, token.Token, "invalid query type")
			}

			// Advance the index, ready for the next step.
			p.pop()

		case stepSelectField:
			// After a SELECT statement we expect a comma separated list of fields or *
			field := p.peek()
			if field.Type != Identifier && field.Type != Asterisk {
				if field.Type == ReservedWord {
					return ErrNoFieldsSelected
				}
				return Error(p.idx, field.Token, "invalid field identifier")
			}

			// Add the field to the list of query fields
			p.query.Fields = append(p.query.Fields, field)

			// Advance the index so that we can peek ahead to determine the next state
			p.pop()

			next := p.peek()
			switch next.Token {
			case COMMA:
				// Pop the comma and advance to the next select field
				p.pop()
				p.step = stepSelectField
			case AS:
				// Advance to the aliasing step
				p.step = stepSelectFieldAlias
			case FROM:
				// Advance to the from table step
				p.step = stepSelectFrom
			default:
				return Error(p.idx, next.Token, "invalid select fields statement")
			}

		case stepSelectFieldAlias:
			// Pop the aliasing reserved word and ensure that the step is correct
			if rword := p.pop(); rword.Token != AS {
				panic(InvalidState(AS, rword.Token))
			}

			// After AS we expect an identifier that we can alias the field to
			alias := p.peek()
			if alias.Type != Identifier {
				return Error(p.idx, alias.Token, "invalid alias identifier")
			}

			// Add the alias to the list of aliases in the query, associating the alias
			// with the previous field added to the list of fields
			field := p.query.Fields[len(p.query.Fields)-1]
			if field.Type == Asterisk {
				return Error(p.idx, alias.Token, "cannot alias *")
			}

			if p.query.Aliases == nil {
				p.query.Aliases = make(map[string]string)
			}
			p.query.Aliases[field.Token] = alias.Token

			// Pop the alias and peek next to determine which state to go to
			p.pop()
			next := p.peek()

			switch next.Token {
			case COMMA:
				// Pop the comma and advance to the next select field
				p.pop()
				p.step = stepSelectField
			case FROM:
				// Advance to the from topic state
				p.step = stepSelectFrom
			default:
				return Error(p.idx, next.Token, "invalid select fields statement")
			}

		case stepSelectFrom:
			// Pop the FROM reserved word and ensure that the step is correct
			if rword := p.pop(); rword.Token != FROM {
				panic(InvalidState(FROM, rword.Token))
			}

			// After FROM we expect an identifier with the topic name
			topic := p.peek()
			if topic.Type != Identifier {
				return Error(p.idx, topic.Token, "invalid topic identifier")
			}

			// Add the topic to the query
			if p.query.Topic.Topic != "" {
				return Error(p.idx, topic.Token, "topic has already been identified")
			}
			p.query.Topic.Topic = topic.Token

			// Pop the topic and peek next to determine the next state
			p.pop()
			next := p.peek()

			switch next.Token {
			case DOT:
				p.pop()
				p.step = stepSelectFromSchema
			case SC, Empty.Token:
				p.step = stepTerm
			case WHERE:
				p.step = stepWhere
			case OFFSET:
				p.step = stepOffset
			case LIMIT:
				p.step = stepLimit
			default:
				return Error(p.idx, next.Token, "invalid from clause")
			}

		case stepSelectFromSchema:
			// The schema follows a dot after the topic name
			schema := p.peek()
			if schema.Type != Identifier && schema.Type != Asterisk {
				return Error(p.idx, schema.Token, "invalid schema identifier")
			}

			// Add the schema to the query topic
			if p.query.Topic.Schema != "" {
				return Error(p.idx, schema.Token, "topic schema has already been identified")
			}
			p.query.Topic.Schema = schema.Token

			// Pop the schema and peek next to determine the next state
			p.pop()
			next := p.peek()

			switch next.Token {
			case DOT:
				// A version can follow a schema if the schema isn't *
				if schema.Type == Asterisk {
					return Error(p.idx, next.Token, "cannot specify version for * schema")
				}
				// Otherwise pop the dot and parse the version
				p.pop()
				p.step = stepSelectFromVersion
			case SC, Empty.Token:
				p.step = stepTerm
			case WHERE:
				p.step = stepWhere
			case OFFSET:
				p.step = stepOffset
			case LIMIT:
				p.step = stepLimit
			default:
				return Error(p.idx, next.Token, "invalid from clause")
			}

		case stepSelectFromVersion:
			// The version follows a dot after the schema name
			version := p.peek()
			if version.Type != Numeric {
				return Error(p.idx, version.Token, "invalid version identifier")
			}

			// Add the version to the query topic
			if p.query.Topic.Version != 0 {
				return Error(p.idx, version.Token, "topic version has already been identified")
			}

			vint, err := version.ParseUint(0, 32)
			if err != nil {
				return Error(p.idx, version.Token, "could not parse schema version")
			}
			p.query.Topic.Version = uint32(vint)

			// Pop the version and peek next to determine the next state
			p.pop()
			next := p.peek()

			switch next.Token {
			case SC, Empty.Token:
				p.step = stepTerm
			case WHERE:
				p.step = stepWhere
			case OFFSET:
				p.step = stepOffset
			case LIMIT:
				p.step = stepLimit
			default:
				return Error(p.idx, next.Token, "invalid from clause")
			}
		case stepWhere:
			// Pop the WHERE reserved word and ensure that the step is correct
			if rword := p.pop(); rword.Token != WHERE {
				panic(InvalidState(WHERE, rword.Token))
			}

		case stepOffset:
			// Pop the OFFSET reserved word and ensure that the step is correct
			if rword := p.pop(); rword.Token != OFFSET {
				panic(InvalidState(OFFSET, rword.Token))
			}

			// After OFFSET we expect a numeric identifier
			offset := p.peek()
			if offset.Type != Numeric {
				return Error(p.idx, offset.Token, "invalid offset")
			}

			// If we've already had an offset defined then error
			if p.query.HasOffset {
				return Error(p.idx, offset.Token, "offset has already been set")
			}

			// Set the offset on the query
			var err error
			if p.query.Offset, err = offset.ParseUint(10, 64); err != nil {
				return Error(p.idx, offset.Token, "could not parse offset")
			}
			p.query.HasOffset = true

			// Pop the offset and peek next to determine next state
			p.pop()
			next := p.peek()

			switch next.Token {
			case LIMIT:
				p.step = stepLimit
			case SC, Empty.Token:
				p.step = stepTerm
			default:
				return Error(p.idx, next.Token, "invalid offset clause")
			}

		case stepLimit:
			// Pop the LIMIT reserved word and ensure that the next step is correct
			if rword := p.pop(); rword.Token != LIMIT {
				panic(InvalidState(LIMIT, rword.Token))
			}

			// After LIMIT we expect a numeric identifier
			limit := p.peek()
			if limit.Type != Numeric {
				return Error(p.idx, limit.Token, "invalid limit")
			}

			// If we've already had a limit defined then error
			if p.query.HasLimit {
				return Error(p.idx, limit.Token, "limit has already been set")
			}

			// Set the limit on the query
			var err error
			if p.query.Limit, err = limit.ParseUint(10, 64); err != nil {
				return Error(p.idx, limit.Token, "could not parse limit")
			}
			p.query.HasLimit = true

			// Pop the limit and peek next to determine next state
			p.pop()
			next := p.peek()

			switch next.Token {
			case OFFSET:
				p.step = stepOffset
			case SC, Empty.Token:
				p.step = stepTerm
			default:
				return Error(p.idx, next.Token, "invalid limit clause")
			}

		case stepTerm:
			// If we reach the termination step then the next token should be empty or ;
			if token := p.pop(); !(token.Token == SC || token.Type == EmptyToken) {
				panic(InvalidState(";", token.Token))
			}

		default:
			// This is a developer error; the parser should never reach a state that is
			// unhandled in this switch statement. Panic ensures tests catch it.
			panic(ErrUnhandledStep)
		}

	}

	// If we've reached the end of the sql query return any errors on the parser.
	return p.err
}

// When the parsing concludes we need to ensure we've reached a valid query state, this
// method checks all of the ways a query can be invalid or partially processed before
// returning a "valid" query struct back to the user. It is assumed that this method is
// called after exec() when parsing has been completed and the index and step have been
// advanced as far as possible.
func (p *parser) validate() error {
	if p.sql == "" {
		return ErrEmptyQuery
	}

	if p.query.Topic.Topic == "" {
		return ErrMissingTopic
	}

	// Validate based on the query type if necessary
	switch p.query.Type {
	case SelectQuery:
		if err := p.validateSelectQuery(); err != nil {
			return err
		}
	}

	return nil
}

func (p *parser) validateSelectQuery() error {
	// The query must have a projection of fields or *
	if len(p.query.Fields) == 0 {
		return ErrNoFieldsSelected
	}

	// If * is in the fields, then it should be the only field
	for _, field := range p.query.Fields {
		if field.Type == Asterisk && len(p.query.Fields) != 1 {
			return ErrInvalidSelectAllFields
		}
	}

	return nil
}

// Pop returns the next token and advances the index of the parser to the end of the
// next token and removes any whitespace that follows it (including new lines).
func (p *parser) pop() Token {
	peeked := p.peek()
	p.idx += peeked.Length
	p.strip()
	return peeked
}

// Peek returns the next token without modifying the underlying state of the parser.
func (p *parser) peek() Token {
	if p.idx >= len(p.sql) {
		return Empty
	}

	// Check to see if the next token is any of our reserved words.
	for _, rWord := range ReservedWords {
		token := strings.ToUpper(p.sql[p.idx:min(len(p.sql), p.idx+len(rWord))])
		if token == rWord {
			return Token{token, ReservedWordType[token], len(token)}
		}
	}

	// If the next char is a single quote attempt to get the quoted value
	if p.sql[p.idx] == SQUOTE {
		return p.peekQuotedString()
	}

	// If the next char is a digit or a - (for negative numbers) get the numeric value
	if p.sql[p.idx] == MINUS || unicode.IsDigit(rune(p.sql[p.idx])) {
		return p.peekNumeric()
	}

	// Finally, attempt to peek an identifier (e.g. a value that is not reserved)
	identifier := p.peekIdentifier()

	// If the identifier is boolean (t, T, TRUE, true, True, f, F, FALSE, false, False)
	// then return the boolean token. Note that 0 and 1 will be returned as numeric and
	// must be converted to a bool from numeric if the boolean type is required.
	if boolre.MatchString(identifier.Token) {
		identifier.Type = Boolean
	}

	return identifier
}

// Returns the token that is inside a pair of single quotes e.g. 'token' ensuring that
// any escaped quotes are included, e.g. 'token\'s' should return token's. Note that the
// enclosing quotes are removed from the token but the length includes the quotes to
// ensure the parser is advanced correctly.
func (p *parser) peekQuotedString() Token {
	// Sanity check -- callers should ensure that the parser is valid before calling
	if p.idx > len(p.sql) || p.sql[p.idx] != SQUOTE {
		return Empty
	}

	// Scan over all of the chars after the quote looking for the closing quote.
	for i := p.idx + 1; i < len(p.sql); i++ {
		// If the next character is a single quote and it is not escaped (e.g. the
		// previous character is not an escape character) then we've found the end.
		// Ensure we return only the part inside the quotes but add 2 to the length to
		// ensure the index is advanced past the single quotes.
		if p.sql[i] == SQUOTE && p.sql[i-1] != ESCAPE {
			token := p.sql[p.idx+1 : i]
			return Token{token, QuotedString, len(token) + 2}
		}
	}

	// If the opening quote is not terminated by an unescaped closing quote then empty
	// is returned. An error is set on the parser for the executor to return.
	p.err = Error(p.idx, p.sql[p.idx:min(len(p.sql), p.idx+5)], "quoted string missing closing quote")
	return Token{"", EmptyToken, len(p.sql) - p.idx}
}

var numre = regexp.MustCompile(`[-\.0-9]`)

func (p *parser) peekNumeric() Token {
	// Numeric matches any positive or negative decimal number (base10) including both
	// integers and floating point numbers that have a . to represent the decmial.
	// Numeric does not currently match scientific notation (e.g. 1e10) or other base
	// systems such as base8 or base16.
	for i := p.idx; i < len(p.sql); i++ {
		if !numre.MatchString(string(p.sql[i])) {
			token := p.sql[p.idx:i]
			return Token{token, Numeric, len(token)}
		}
	}

	// If we get to the end of the string return the remainder as numeric
	token := p.sql[p.idx:]
	return Token{token, Numeric, len(token)}
}

var (
	identre = regexp.MustCompile(`[a-zA-Z0-9_]`)
	boolre  = regexp.MustCompile(`^(1|t|T|True|true|TRUE|0|f|F|False|false|FALSE)$`)
)

func (p *parser) peekIdentifier() Token {
	// An identifier is any word that contains letters, digits, or underscore and is
	// not surrounded by quotation marks. Identifiers cannot begin with a digit,
	// otherwise they will be parsed as numeric; they can start with underscore. No
	// punctuation, including asterisk is parsed by the identifier.
	for i := p.idx; i < len(p.sql); i++ {
		if !identre.MatchString(string(p.sql[i])) {
			token := p.sql[p.idx:i]
			return Token{token, Identifier, len(token)}
		}
	}

	// Return the entire remainder of the string if we get to the end.
	token := p.sql[p.idx:]
	return Token{token, Identifier, len(token)}
}

// Strip whitespace by advancing the index of the parser until it is not pointing to a
// whitespace character (defined by unicode).
func (p *parser) strip() {
	for {
		if p.idx < len(p.sql) && unicode.IsSpace(rune(p.sql[p.idx])) {
			p.idx++
		} else {
			return
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
