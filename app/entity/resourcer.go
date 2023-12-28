package entity

type Resourcer interface {
	GetEntityName() string
	GetFieldStructTag(fieldName string, structTag string) string
}
