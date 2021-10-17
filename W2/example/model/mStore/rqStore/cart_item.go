package rqStore

import (
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/W2/example/conf"
)

func (c *CartItems) FindByOwnerIdInvoiceId() (res []*CartItems, total uint32) {
	query := `
SELECT ` + c.sqlSelectAllFields() + `
FROM ` + c.sqlTableName() + `
WHERE ` + c.sqlOwnerId() + ` = ` + I.UToS(c.OwnerId) +`
	AND ` + c.sqlInvoiceId() + ` = 0
ORDER BY ` + c.sqlProductId() + `
`
	if conf.DEBUG_MODE {
		L.Print(query)
	}
	
	c.Adapter.QuerySql(query, func(row []interface{}) {
		obj := &CartItems{}
		obj.FromArray(row)
		total += uint32(obj.Qty)
		res = append(res, obj)
	})
	return
}
