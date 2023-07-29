package ensql_test

import (
	"testing"

	. "github.com/rotationalio/ensign/pkg/ensign/ensql"
	"github.com/stretchr/testify/require"
)

func TestConditionGroup(t *testing.T) {
	age := Token{"age", Identifier, 3}
	a21 := Token{"21", Numeric, 2}
	a65 := Token{"65", Numeric, 2}
	color := Token{"color", Identifier, 5}
	eq := Token{EQ, OperatorToken, len(EQ)}
	gt := Token{GT, OperatorToken, len(GT)}
	lte := Token{LTE, OperatorToken, len(LTE)}
	red := Token{"'red'", QuotedString, 5}
	blue := Token{"'blue'", QuotedString, 6}

	// From an empty group should only be able to call ConditionLeft and OpenParen
	group := NewConditionGroup()
	require.ErrorIs(t, group.ConditionOperator(eq), ErrAppendCondition)
	require.ErrorIs(t, group.ConditionRight(red), ErrAppendCondition)
	require.ErrorIs(t, group.LogicalOperator(And), ErrAppendOperator)
	require.ErrorIs(t, group.CloseParens(), ErrCloseParens)

	require.NoError(t, group.OpenParens())
	require.Equal(t, "()", group.String())

	// Reset condition group - should be able to add a condition
	group = NewConditionGroup()
	require.NoError(t, group.ConditionLeft(color))
	require.NoError(t, group.ConditionOperator(eq))
	require.NoError(t, group.ConditionRight(red))
	require.Equal(t, "color = 'red'", group.String())

	// Should be able to add a logical operator after the condition
	require.NoError(t, group.LogicalOperator(Or))
	require.NoError(t, group.ConditionLeft(color))
	require.NoError(t, group.ConditionOperator(eq))
	require.NoError(t, group.ConditionRight(blue))
	require.Equal(t, "color = 'red' OR color = 'blue'", group.String())

	// Should be able to add a paren
	require.NoError(t, group.LogicalOperator(Or))
	require.NoError(t, group.OpenParens())
	require.NoError(t, group.ConditionLeft(age))
	require.NoError(t, group.ConditionOperator(gt))
	require.NoError(t, group.ConditionRight(a21))
	require.NoError(t, group.LogicalOperator(And))
	require.NoError(t, group.ConditionLeft(age))
	require.NoError(t, group.ConditionOperator(lte))
	require.NoError(t, group.ConditionRight(a65))
	require.Equal(t, "color = 'red' OR color = 'blue' OR (age > 21 AND age <= 65)", group.String())
}

func TestConditionIsPartial(t *testing.T) {
	partials := []Condition{
		{Token{"color", Identifier, 5}, Empty, Empty},
		{Empty, Token{EQ, OperatorToken, len(EQ)}, Empty},
		{Empty, Empty, Token{"'red'", QuotedString, 5}},
		{Empty, Empty, Empty},
		{Token{"color", Identifier, 5}, Token{EQ, OperatorToken, len(EQ)}, Empty},
		{Token{"color", Identifier, 5}, Empty, Token{"'red'", QuotedString, 5}},
		{Empty, Token{EQ, OperatorToken, len(EQ)}, Token{"'red'", QuotedString, 5}},
		{},
		{Left: Token{"color", Identifier, 5}},
		{Left: Token{"color", Identifier, 5}, Operator: Token{EQ, OperatorToken, len(EQ)}},
	}

	for i, partial := range partials {
		require.True(t, partial.IsPartial(), "expected test case %d to be partial", i)
	}

	valid := []Condition{
		{Token{"color", Identifier, 5}, Token{EQ, OperatorToken, len(EQ)}, Token{"'red'", QuotedString, 5}},
		{Token{"age", Identifier, 3}, Token{GTE, OperatorToken, len(GTE)}, Token{"21", Numeric, 2}},
	}

	for i, condition := range valid {
		require.False(t, condition.IsPartial(), "expected test case %d to not be partial", i)
	}
}
