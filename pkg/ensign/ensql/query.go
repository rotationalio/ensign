package ensql

import (
	"fmt"
	"strings"
)

// Query is a parsed representation of an EnSQL statement that can be used to process
// EnSQL requests. All elements from the SQL statement must be represented in the query
// (including things like comments if supported by EnSQL). The raw sql is also available
// on the query for debugging purposes.
type Query struct {
	Type       QueryType
	Topic      Topic
	Conditions *ConditionGroup
	Fields     []Token
	Aliases    map[string]string
	Offset     uint64
	HasOffset  bool
	Limit      uint64
	HasLimit   bool
	Raw        string
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

// Represents the topic and event types that are processed by the query.
type Topic struct {
	Topic   string
	Schema  string
	Version uint32
}

// Condition represents a basic expression in a where clause.
type Condition struct {
	Left     Token
	Operator Token
	Right    Token
}

// Returns true if the condition has not been completed yet.
func (c Condition) IsPartial() bool {
	if c.Left.Type == UnknownTokenType || c.Operator.Type == UnknownTokenType || c.Right.Type == UnknownTokenType {
		return true
	}
	return c.Left == Empty || c.Operator == Empty || c.Right == Empty
}

func (c Condition) String() string {
	return fmt.Sprintf("%s %s %s", c.Left.Token, c.Operator.Token, c.Right.Token)
}

// ConditionGroup represents a list of subgroups, logical operators, and conditions.
// It is not an operational type, e.g. it cannot be used for evaluating an expression,
// but it is used as a substep to building a Predicate object from a series of nested
// conditions that are evaluated from a tokenizer.
type ConditionGroup struct {
	children []fmt.Stringer
	current  *ConditionGroup
	parent   *ConditionGroup
	isParens bool
}

// Initialize a ConditionGroup for parsing a predicate tree with precedence.
func NewConditionGroup() *ConditionGroup {
	cg := &ConditionGroup{
		children: make([]fmt.Stringer, 0),
		parent:   nil,
		isParens: false,
	}
	cg.current = cg
	return cg
}

// Create a new condition with the left-token, appending it to the current condition
// group. Will return an error if the previous object is not a logical operator.
func (g *ConditionGroup) ConditionLeft(token Token) error {
	if g.curlen() != 0 {
		if g.prevType() != ctLogicalOperator {
			return ErrAppendCondition
		}
	}

	// Append a new condition group to the current children and mark it as now current.
	condition := &Condition{Left: token}
	g.current.children = append(g.current.children, condition)
	return nil
}

// Update the previous condition with an operator. Will return an error if the previous
// object is not a partial condition.
func (g *ConditionGroup) ConditionOperator(token Token) error {
	if g.curlen() == 0 || g.prevType() != ctPartialCondition {
		return ErrAppendCondition
	}

	prev := g.prev()
	condition, ok := prev.(*Condition)
	if !ok {
		return ErrAppendCondition
	}

	condition.Operator = token
	return nil
}

// Update the previous condition with the right token. Will return an error if the
// previous object is not a partial condition. Expects that the condition will no longer
// be partial after setting the right side value.
func (g *ConditionGroup) ConditionRight(token Token) error {
	if g.curlen() == 0 || g.prevType() != ctPartialCondition {
		return ErrAppendCondition
	}

	prev := g.prev()
	condition, ok := prev.(*Condition)
	if !ok {
		return ErrAppendCondition
	}

	condition.Right = token
	if condition.IsPartial() {
		return ErrAppendCondition
	}
	return nil
}

// Append a logical operator to the current condition group. Will return an error if the
// previous object is not a closed parentheses group or a condition.
func (g *ConditionGroup) LogicalOperator(op Operator) error {
	// Append the operator to the current group list so long as the conditions are met.
	// The operator cannot be the first thing in the group unless it is NOT
	// TODO: refactor logic to handle NOT
	if g.curlen() == 0 {
		return ErrAppendOperator
	}

	// The prev type must be a condition or condition group
	prevType := g.prevType()
	if !(prevType == ctCondition || prevType == ctConditionGroup) {
		return ErrCloseParens
	}

	g.current.children = append(g.current.children, op)
	return nil
}

// OpenParens creates a new condition group that is a subgroup of the current group.
// It will return an error if the parens is not the first thing in the group or if it
// is not proceeded by a logical operator token.
func (g *ConditionGroup) OpenParens() error {
	// When opening parens, we're effectively creating a new subgroup and appending it
	// to the current node's children, then marking that node as the current node to
	// append any following items to until the parens are closed.

	// If this is not the first child in the current group, then the prev node must be
	// a logical operator otherwise we cannot open the parentheses.
	if g.curlen() != 0 {
		if g.prevType() != ctLogicalOperator {
			return ErrOpenParens
		}
	}

	// Append a new condition group to the current children and mark it as now current.
	parens := &ConditionGroup{children: make([]fmt.Stringer, 0), parent: g.current, isParens: true}
	g.current.children = append(g.current.children, parens)
	g.current = parens

	return nil
}

// CloseParens closes the current subgroup and returns control to the containing group.
// It returns an error if the last item in the parentheses is not another parentheses or
// a condition (e.g. parens cannot be closed after a logical operator nor can we parse
// empty parentheses in a sql query).
func (g *ConditionGroup) CloseParens() error {
	// When closing parens, we're effectively closing a subgroup and returning control
	// to that subgroups parent node to continue appending to.

	// Cannot close a group that is not parentheses or does not have a parent.
	if !g.current.isParens || g.current.parent == nil {
		return ErrCloseParens
	}

	// Cannot have empty parentheses
	if g.curlen() == 0 {
		return ErrCloseParens
	}

	// The prev type must be a condition or condition group
	prevType := g.prevType()
	if !(prevType == ctCondition || prevType == ctConditionGroup) {
		return ErrCloseParens
	}

	g.current = g.current.parent
	return nil
}

func (g *ConditionGroup) String() string {
	var sb strings.Builder
	for _, child := range g.children {
		sb.WriteString(child.String())
		sb.WriteString(" ")
	}

	if g.isParens {
		return "(" + strings.TrimSpace(sb.String()) + ")"
	}
	return strings.TrimSpace(sb.String())
}

func (g *ConditionGroup) curlen() int {
	return len(g.current.children)
}

func (g *ConditionGroup) prev() any {
	return g.current.children[len(g.current.children)-1]
}

func (g *ConditionGroup) prevType() conditionType {
	prev := g.prev()
	switch c := prev.(type) {
	case *Condition:
		if c.IsPartial() {
			return ctPartialCondition
		}
		return ctCondition
	case Operator:
		return ctLogicalOperator
	case *ConditionGroup:
		return ctConditionGroup
	default:
		return ctUnknown
	}
}

type conditionType uint8

const (
	ctUnknown conditionType = iota
	ctCondition
	ctPartialCondition
	ctLogicalOperator
	ctConditionGroup
)
