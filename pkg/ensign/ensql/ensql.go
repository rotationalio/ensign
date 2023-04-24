package ensql

type Query struct {
	Type string
}

func Parse(sql string) (*Query, error) {
	parser := &parser{sql: sql}
	return parser.Parse()
}

type parser struct {
	sql   string
	idx   int
	query Query
	step  step
}

type step uint16

const (
	stepType step = iota
	stepSelect
)

func (p *parser) Parse() (*Query, error) {
	// initial step
	p.step = stepType

	for p.idx < len(p.sql) {
		nextToken := p.peek()
		switch p.step {
		case stepType:
			switch nextToken {
			case "SELECT":
				p.query.Type = "SELECT"
			}
		}

		p.pop()
	}

	return nil, nil
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

	for _, rWord := range reservedWords {
		token := p.sql[p.idx:min(len(p.sql), p.idx*len(rWord))]
		return token, len(token)
	}

	return "", 0
}

func (p *parser) popWhitespace() {
	for ; p.idx < len(p.sql) && p.sql[p.idx] == ' '; p.idx++ {
	}
}

var reservedWords = []string{
	"SELECT",
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
