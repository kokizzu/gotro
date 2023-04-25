package mZzz

import "github.com/kokizzu/gotro/D/Tt"

const (
	TableZzz  Tt.TableName = `zzz`
	Id                     = `id`
	CreatedAt              = `created_at`
	Coords                 = `coords`
)

var TarantoolTables = map[Tt.TableName]*Tt.TableProp{
	// can only adding fields on back, and must IsNullable: true
	// primary key must be first field and set to UniqueX: fieldName or AutoIncrementId: true
	TableZzz: {
		Fields: []Tt.Field{
			{Id, Tt.Unsigned},
			{CreatedAt, Tt.Integer},
			{Coords, Tt.Array},
		},
		Engine:          Tt.Memtx,
		AutoIncrementId: true,
		Spatial:         Coords, // spatial index only works for memtx
	},
}
