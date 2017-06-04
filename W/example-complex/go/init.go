package main

import (
	"github.com/kokizzu/gotro/D/Pg"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/W"
	"github.com/kokizzu/gotro/W/example-complex/model"
	"github.com/kokizzu/gotro/W/example-complex/model/mUsers"
)

func main() {
	model.PG_W.CreateBaseTable(mUsers.TABLE, mUsers.TABLE)
	model.PG_W.CreateBaseTable(`todos`, mUsers.TABLE) // 2nd table
	ajax := W.NewAjax()
	rm := W.NewRequestModel_ById_ByDbActor_ByAjax(`1`, `1`, ajax)
	model.PG_W.DoTransaction(func(tx *Pg.Tx) string {
		dm := Pg.NewRow(tx, mUsers.TABLE, rm)
		dm.SetVal(`full_name`, `CHANGEME`)
		dm.Set_UserEmails(`CHANGEME@gmail.com`)
		dm.UpsertRow()
		return ajax.LastError()
	})
	L.Print(ajax)
}
