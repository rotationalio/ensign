package ensql

import "strings"

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
	// NOTE: p.pop() must be called to advance the index and guarantee termination.
	for p.idx < len(p.sql) {
		switch p.step {
		case stepInit:
			// At the initial step we expect a query determiner such as SELECT or WITH
			switch p.peek() {
			case SELECT:
				p.query.Type = SelectQuery
				p.step = stepSelectField
			case "":
				return ErrEmptyQuery
			default:
				return Error(p.idx, "", "invalid query type")
			}
		}

		// Advance the index, ready for the next step.
		p.pop()
	}

	// If we've reached the end of the sql query return any errors on the parser.
	return p.err
}

func (p *parser) validate() error {
	if p.sql == "" {
		return ErrEmptyQuery
	}

	if p.query.Topic.Topic == "" {
		return ErrMissingTopic
	}

	return nil
}

func (p *parser) peek() string {
	peeked, _ := p.peekWithLength()
	return peeked
}

func (p *parser) pop() string {
	peeked, length := p.peekWithLength()
	p.idx = length
	p.popWhitespace()
	return peeked
}

func (p *parser) peekWithLength() (string, int) {
	if p.idx >= len(p.sql) {
		return "", 0
	}

	for _, rWord := range ReservedWords {

		token := strings.ToUpper(p.sql[p.idx:min(len(p.sql), p.idx+len(rWord))])
		if token == rWord {
			return token, len(token)
		}
	}

	return "", 0
}

func (p *parser) popWhitespace() {
	for ; p.idx < len(p.sql) && p.sql[p.idx] == ' '; p.idx++ {
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
