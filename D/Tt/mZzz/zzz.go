package mZzz

import "github.com/kokizzu/gotro/D/Tt"

const (
	TableZzz  Tt.TableName = `zzz`
	Id                     = `id`
	CreatedAt              = `created_at`
)

var TarantoolTables = map[Tt.TableName]*Tt.TableProp{
	// can only adding fields on back, and must IsNullable: true
	// primary key must be first field and set to UniqueX: fieldName or AutoIncrementId: true
	TableZzz: {
		Fields: []Tt.Field{
			{Id, Tt.Unsigned},
			{CreatedAt, Tt.Integer},
		},
		Engine:          Tt.Memtx,
		AutoIncrementId: true,
	},
}
