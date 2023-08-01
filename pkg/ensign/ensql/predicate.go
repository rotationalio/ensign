package ensql

import "errors"

// Predicate implements a binary abstract syntax tree that can be used to evaluate
// complex predicate expressions that are created from where clauses. Predicates should
// be parsed into this abstract syntax tree with a precedence order defined.
type Predicate struct {
	Left     any
	Operator Operator
	Right    any
}

// Validate the predicate tree to ensure it is syntactically correct.
func (p Predicate) Validate() error {
	switch p.Type() {
	case ComparisonPredicate:
		// Comparison predicate must be a leaf node with left and right tokens.
		if token, ok := p.Left.(Token); !ok {
			return ErrInvalidPredicate
		} else if token.Type != Identifier {
			return ErrInvalidPredicate
		}

		if token, ok := p.Right.(Token); !ok {
			return ErrInvalidPredicate
		} else {
			if !(token.Type == QuotedString || token.Type == Numeric || token.Type == Boolean) {
				return ErrInvalidPredicate
			}
		}
		return nil
	case SearchPredicate:
		// SearchPredicate must be a leaf node with a left identifier token and a right
		// quoted string token otherwise the predicate cannot be evaluated.
		if token, ok := p.Left.(Token); !ok {
			return ErrInvalidPredicate
		} else if token.Type != Identifier {
			return ErrInvalidPredicate
		}

		if token, ok := p.Right.(Token); !ok {
			return ErrInvalidPredicate
		} else if token.Type != QuotedString {
			return ErrInvalidPredicate
		}

		return nil
	case LogicalPredicate:
		// Logical predicates must have two leaf nodes that are themselves predicates.
		if left, ok := p.Left.(Predicate); !ok {
			return ErrInvalidPredicate
		} else {
			if err := left.Validate(); err != nil {
				return err
			}
		}

		if right, ok := p.Right.(Predicate); !ok {
			return ErrInvalidPredicate
		} else {
			if err := right.Validate(); err != nil {
				return err
			}
		}

		return nil
	default:
		return ErrPredicateType
	}
}

// Evaluate the predicate tree for the specified variables.
// TODO: from the predicate tree, extract a comparable representation that doesn't have
// to parse numbers/bools/etc every time and can be applied more efficiently to a large
// number of events.
func (p Predicate) Evaluate(vars ...any) (bool, error) {
	switch p.Type() {
	case ComparisonPredicate:
		return p.compare(vars...)
	case SearchPredicate:
		return p.search(vars...)
	case LogicalPredicate:
		return p.boolean(vars...)
	default:
		return false, ErrPredicateType
	}
}

// Compare requires the left value to be an identifier token and the right value to be
// either a numeric, quoted string, or boolean value to ensure the comparison operation
// happens correctly with the specified variables.
func (p Predicate) compare(vars ...any) (bool, error) {
	// TODO: implement
	return false, errors.New("not implemented yet")
}

// Search implements the like and ilike operators. The left value should be an
// identifier and the right token should be a quoted string.
func (p Predicate) search(vars ...any) (bool, error) {
	// TODO: implement
	return false, errors.New("not implemented yet")
}

// Boolean implements logical operations. The left and right values should be predicates
// such that the left predicate is evaluated first, then the right predicate.
func (p Predicate) boolean(vars ...any) (_ bool, err error) {
	var (
		ok    bool
		left  Predicate
		right Predicate
		lres  bool
		rres  bool
	)

	if left, ok = p.Left.(Predicate); !ok {
		return false, ErrInvalidPredicate
	}

	if right, ok = p.Right.(Predicate); !ok {
		return false, ErrInvalidPredicate
	}

	if lres, err = left.Evaluate(vars...); err != nil {
		return false, err
	}

	if rres, err = right.Evaluate(vars...); err != nil {
		return false, err
	}

	switch p.Operator {
	case And:
		return lres && rres, nil
	case Or:
		return lres || rres, nil
	default:
		return false, ErrInvalidPredicate
	}
}

func (p Predicate) Type() PredicateType {
	switch p.Operator {
	case And, Or:
		// Logical predicates are intermediate nodes
		return LogicalPredicate
	case Eq, Ne, Gt, Lt, Gte, Lte:
		// Comparison predicates are leaf nodes
		return ComparisonPredicate
	case Like, ILike:
		// Search predicates are leaf nodes
		return SearchPredicate
	default:
		return UnknownPredicateType
	}
}

type PredicateType uint8

const (
	UnknownPredicateType PredicateType = iota
	LogicalPredicate
	ComparisonPredicate
	SearchPredicate
)
