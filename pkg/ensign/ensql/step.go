package ensql

type step uint16

const (
	stepInit step = iota
	stepTerm
	stepSelect
	stepSelectField
	stepSelectFieldAlias
	stepSelectFrom
	stepSelectComma
	stepSelectFromTopic
	stepSelectDot
	stepSelectFromSchema
	stepSelectFromVersion
	stepWhere
	stepWhereField
	stepWhereOperator
	stepWhereValue
	stepWhereAnd
	stepWhereOr
	stepWhereOpenParen
	stepWhereCloseParen
	stepOffset
	stepOffsetValue
	stepLimit
	stepLimitValue
)
