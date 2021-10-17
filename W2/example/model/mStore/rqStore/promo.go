package rqStore

import (
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/W2/example/conf"
	"github.com/kpango/fastime"
)

func (p *Promos) FindActive() (res []*Promos) {
	query := `
SELECT ` + p.sqlSelectAllFields() + `
FROM ` + p.sqlTableName() + `
WHERE ` + p.sqlStartAt() + ` <= ` + I.ToS(fastime.UnixNow()) + `
	AND ` + p.sqlEndAt() + ` >= ` + I.ToS(fastime.UnixNow()) + `
ORDER BY ` + p.sqlId() + `
`
	if conf.DEBUG_MODE {
		L.Print(query)
	}

	p.Adapter.QuerySql(query, func(row []interface{}) {
		obj := &Promos{}
		obj.FromArray(row)
		res = append(res, obj)
	})
	return
}
