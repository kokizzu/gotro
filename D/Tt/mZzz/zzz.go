package mZzz

import "github.com/kokizzu/gotro/D/Tt"

const (
	TableZzz    Tt.TableName = `zzz`
	Id                       = `id`
	CreatedAt                = `createdAt`
	Coords                   = `coords`
	Name                     = `name`
	HeightMeter              = `heightMeter`
)

var TarantoolTables = map[Tt.TableName]*Tt.TableProp{
	// can only adding fields on back, and must IsNullable: true
	// primary key must be first field and set to UniqueX: fieldName or AutoIncrementId: true
	TableZzz: {
		Fields: []Tt.Field{
			{Id, Tt.Unsigned},
			{CreatedAt, Tt.Integer},
			{Coords, Tt.Array},
			{Name, Tt.String},
			{HeightMeter, Tt.Double},
		},
		Engine:          Tt.Memtx,
		AutoIncrementId: true,
		Unique1:         Name,
		Spatial:         Coords, // spatial index only works for memtx
	},
}
