package ensql_test

import (
	"testing"

	. "github.com/rotationalio/ensign/pkg/ensign/ensql"
	"github.com/stretchr/testify/require"
)

func TestValidatePredicate(t *testing.T) {
	// Define some tokens for testing predicates
	col := Token{"color", Identifier, 3}
	str := Token{"'red'", QuotedString, 5}
	blu := Token{"'blue'", QuotedString, 6}
	itg := Token{"42", Numeric, 2}
	flt := Token{"3.14", Numeric, 4}
	bol := Token{"True", Boolean, 4}
	pre := Token{"'prefix%'", QuotedString, 9}
	cer := Predicate{col, Eq, str}
	ceb := Predicate{col, Eq, blu}
	inv := Predicate{ceb, And, 23}

	testCases := []struct {
		input Predicate
		err   error
		msg   string
	}{
		{Predicate{col, Eq, str}, nil, "color = 'red'"},
		{Predicate{col, Gt, itg}, nil, "color > 42"},
		{Predicate{col, Gte, itg}, nil, "color >= 42"},
		{Predicate{col, Lt, flt}, nil, "color < 3.14"},
		{Predicate{col, Lte, flt}, nil, "color <= 3.14"},
		{Predicate{col, Ne, bol}, nil, "color != True"},
		{Predicate{42, Eq, flt}, ErrInvalidPredicate, "? = 3.14"},
		{Predicate{flt, Gte, itg}, ErrInvalidPredicate, "3.14 >= 42"},
		{Predicate{col, Eq, 42}, ErrInvalidPredicate, "red = ?"},
		{Predicate{col, Lte, col}, ErrInvalidPredicate, "red = red"},

		{Predicate{col, Like, pre}, nil, "color like 'prefix%'"},
		{Predicate{col, ILike, pre}, nil, "color ilike 'prefix%'"},
		{Predicate{flt, ILike, pre}, ErrInvalidPredicate, "3.14 ilike 'prefix%'"},
		{Predicate{col, Like, itg}, ErrInvalidPredicate, "color like 42"},
		{Predicate{cer, ILike, pre}, ErrInvalidPredicate, "color = 'red' ilike 'prefix%'"},
		{Predicate{col, Like, cer}, ErrInvalidPredicate, "color ilike color = 'red'"},

		{Predicate{ceb, And, cer}, nil, "color = 'blue' AND color = 'red'"},
		{Predicate{cer, Or, ceb}, nil, "color = 'red' OR color = 'blue'"},
		{Predicate{ceb, And, 23}, ErrInvalidPredicate, "color = 'blue' AND ?"},
		{Predicate{23, Or, ceb}, ErrInvalidPredicate, "? OR color = 'blue'"},
		{Predicate{inv, Or, cer}, ErrInvalidPredicate, "color = 'blue AND ? OR color = 'red'"},
		{Predicate{ceb, And, inv}, ErrInvalidPredicate, "color = 'blue AND color = 'blue' AND ?"},

		{Predicate{col, UnknownOperator, itg}, ErrPredicateType, "unknown operator"},
	}

	for i, tc := range testCases {
		err := tc.input.Validate()
		if tc.err == nil {
			require.NoError(t, err, "expected valid predicate for test case %d: %s", i, tc.msg)
		} else {
			require.ErrorIs(t, err, tc.err, "expected invalid predicate for test case %d: %s", i, tc.msg)
		}
	}
}

func TestNestedPredicate(t *testing.T) {
	// (color = 'red' OR color = 'blue' AND country = 'fr') AND (age > 18 OR AGE < 35)
	pred := Predicate{
		Predicate{
			Predicate{
				Token{"color", Identifier, 5},
				Eq,
				Token{"'red'", QuotedString, 5},
			},
			Or,
			Predicate{
				Predicate{
					Token{"color", Identifier, 5},
					Eq,
					Token{"'blue'", QuotedString, 6},
				},
				And,
				Predicate{
					Token{"country", Identifier, 7},
					Eq,
					Token{"'fr'", QuotedString, 4},
				},
			},
		},
		And,
		Predicate{
			Predicate{
				Token{"age", Identifier, 3},
				Gt,
				Token{"18", Numeric, 2},
			},
			Or,
			Predicate{
				Token{"age", Identifier, 3},
				Lt,
				Token{"35", Numeric, 2},
			},
		},
	}

	require.NoError(t, pred.Validate(), "expected nested predicate to be valid")
}
