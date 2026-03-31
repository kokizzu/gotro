package mZzz

import "github.com/kokizzu/gotro/D/Ch"

const (
	TableZzz    Ch.TableName = `zzz`
	Id                       = `id`
	CreatedAt                = `createdAt`
	Name                     = `name`
	HeightMeter              = `heightMeter`
)

var ClickhouseTables = map[Ch.TableName]*Ch.TableProp{
	TableZzz: {
		Fields: []Ch.Field{
			{Id, Ch.UInt64},
			{CreatedAt, Ch.DateTime},
			{Name, Ch.String},
			{HeightMeter, Ch.Float64},
		},
		Engine: Ch.ReplacingMergeTree,
		Orders: []string{Id},
	},
}

func GenerateORM() {
	Ch.GenerateOrm(ClickhouseTables)
}
