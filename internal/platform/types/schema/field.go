package schema

type Field struct {
	column string // database column name
	field  string // json field name
}

func (f Field) Column() string {
	return f.column
}

func (f Field) JSONField() string {
	return f.field
}

func NewField(dbField, jsonField string) Field {
	return Field{
		column: dbField,
		field:  jsonField,
	}
}

var (
	CountField = NewField("COUNT(*) as count", "count")
)
