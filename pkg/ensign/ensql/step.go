package ensql

type step uint16

const (
	stepInit step = iota
	stepTerm
	stepSelect
	stepSelectField
	stepSelectFieldAlias
	stepSelectFrom
	stepSelectFromSchema
	stepSelectFromVersion
	stepWhere
	stepWhereField
	stepWhereOperator
	stepWhereValue
	stepWhereLogical
	stepWhereCloseParens
	stepOffset
	stepOffsetValue
	stepLimit
	stepLimitValue
)
