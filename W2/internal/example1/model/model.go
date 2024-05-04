package model

import (
	"example1/conf"
	"example1/model/mAuth"
	"example1/model/mAuth/wcAuth"
	"github.com/kokizzu/gotro/D/Ch"
	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/L"
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
