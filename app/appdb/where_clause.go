package appdb

type WhereClause struct {
	Column    string
	Condition DBCondition
	Value     interface{}
	Logic     DBLogic
}

func NewWhere(column string, condition DBCondition, value interface{}, logic ...DBLogic) WhereClause {
	if len(logic) > 0 {
		return WhereClause{Column: column, Condition: condition, Value: value, Logic: logic[0]}
	}
	return WhereClause{Column: column, Condition: condition, Value: value}
}
