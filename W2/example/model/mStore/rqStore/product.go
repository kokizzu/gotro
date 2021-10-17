package rqStore

import (
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/W2/example/conf"
	"github.com/kokizzu/gotro/X"
)

func (p *Products) FindOffsetLimit(offset, limit uint32) (res []*Products) {
	query := `
SELECT ` + p.sqlSelectAllFields() + `
FROM ` + p.sqlTableName() + `
ORDER BY ` + p.sqlId() + `
LIMIT ` + X.ToS(limit) + `
OFFSET ` + X.ToS(offset)
	if conf.DEBUG_MODE {
		L.Print(query)
	}
	p.Adapter.QuerySql(query, func(row []interface{}) {
		obj := &Products{}
		obj.FromArray(row)
		res = append(res, obj)
	})
	return
}
