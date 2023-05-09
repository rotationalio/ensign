package ensql

// Query is a parsed representation of an EnSQL statement that can be used to process
// EnSQL requests. All elements from the SQL statement must be represented in the query
// (including things like comments if supported by EnSQL). The raw sql is also available
// on the query for debugging purposes.
type Query struct {
	Type       QueryType
	Topic      Topic
	Conditions []Condition
	Fields     []Token
	Aliases    map[string]string
	Offset     uint64
	HasOffset  bool
	Limit      uint64
	HasLimit   bool
	Raw        string
}

// Represents the topic and event types that are processed by the query.
type Topic struct {
	Topic   string
	Schema  string
	Version uint32
}

type Condition struct {
	Left     Token
	Operator Operator
	Right    Token
}

// The raw query is returned as the string representation of the query.
func (q Query) String() string {
	return q.Raw
}

// The type of the EnSQL query (e.g. SELECT)
type QueryType uint8

const (
	UnknownQueryType = iota
	SelectQuery
)

func (q QueryType) String() string {
	switch q {
	case SelectQuery:
		return SELECT
	default:
		return "UNKNOWN"
	}
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
