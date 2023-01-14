package model

import (
	"github.com/kokizzu/gotro/D/Ch"
	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/W2/example/conf"
	"github.com/kokizzu/gotro/W2/example/model/mAuth"
	"github.com/kokizzu/gotro/W2/example/model/mAuth/wcAuth"
)

type Migrator struct {
	Taran *Tt.Adapter
	Click *Ch.Adapter
}

func RunMigration() {
	L.Print(`run migration..`)
	m := Migrator{}
	m.Taran = &Tt.Adapter{Connection: conf.ConnectTarantool(), Reconnect: conf.ConnectTarantool}
	m.Click = &Ch.Adapter{DB: conf.ConnectClickhouse(), Reconnect: conf.ConnectClickhouse}
	m.Taran.MigrateTables(mAuth.TarantoolTables)
	m.Click.MigrateTables(mAuth.ClickhouseTables)
	//m.Click.MigrateTables(mStore.ClickhouseTables)

	rootUser := wcAuth.NewUsersMutator(m.Taran)
	rootUser.SetId(1)
	rootUser.FindById()
	rootUser.SetEmail(`root@localhost`)
	rootUser.SetEncryptPassword(`test123`)
	rootUser.DoUpsert()
}
