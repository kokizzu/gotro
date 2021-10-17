package model

import (
	"github.com/kokizzu/gotro/D/Ch"
	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/W2/example/conf"
	"github.com/kokizzu/gotro/W2/example/model/mAuth"
	"github.com/kokizzu/gotro/W2/example/model/mAuth/wcAuth"
	"github.com/kokizzu/gotro/W2/example/model/mStore"
	"github.com/kokizzu/gotro/W2/example/model/mStore/wcStore"
	"github.com/kokizzu/id64"
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
	m.Taran.MigrateTables(mStore.TarantoolTables)
	//m.Click.MigrateTables(mStore.ClickhouseTables)

	rootUser := wcAuth.NewUsersMutator(m.Taran)
	rootUser.SetId(1)
	rootUser.FindById()
	rootUser.SetEmail(`root@localhost`)
	rootUser.SetEncryptPassword(`test123`)
	rootUser.DoReplace()

	product1 := wcStore.NewProductsMutator(m.Taran)
	product1.SetSku(`120P90`)
	if !product1.FindBySku() {
		product1.SetId(id64.UID())
	}
	product1.SetName(`Google Home`)
	product1.SetPrice(49_99)
	product1.SetInventoryQty(10)
	product1.DoReplace()

	product2 := wcStore.NewProductsMutator(m.Taran)
	product2.SetSku(`43N23P`)
	if !product2.FindBySku() {
		product2.SetId(id64.UID())
	}
	product2.SetName(`MacBook Pro`)
	product2.SetPrice(5_399_99)
	product2.SetInventoryQty(5)
	product2.DoReplace()

	product3 := wcStore.NewProductsMutator(m.Taran)
	product3.SetSku(`A304SD`)
	if !product3.FindBySku() {
		product3.SetId(id64.UID())
	}
	product3.SetName(`Alexa Speaker`)
	product3.SetPrice(109_50)
	product3.SetInventoryQty(10)
	product3.DoReplace()

	product4 := wcStore.NewProductsMutator(m.Taran)
	product4.SetSku(`234234`)
	if !product4.FindBySku() {
		product4.SetId(id64.UID())
	}
	product4.SetName(`Raspberry Pi B`)
	product4.SetPrice(30_00)
	product4.SetInventoryQty(2)
	product4.DoReplace()

	promo1 := wcStore.NewPromosMutator(m.Taran)
	promo1.SetId(11)
	promo1.SetProductId(product2.Id)
	promo1.SetProductCount(1)
	promo1.SetFreeProductId(product4.Id) // free product4
	promo1.FindById()
	promo1.DoReplace()

	promo2 := wcStore.NewPromosMutator(m.Taran)
	promo2.SetId(12)
	promo2.SetProductId(product1.Id)
	promo2.SetProductCount(3)
	promo2.SetDiscountCount(1) // free 1
	promo2.FindById()
	promo2.DoReplace()

	promo3 := wcStore.NewPromosMutator(m.Taran)
	promo3.SetId(13)
	promo3.SetProductId(product3.Id)
	promo3.SetProductCount(2)
	promo3.SetDiscountPercent(10) // 10% to all product3
	promo3.FindById()
	promo3.DoReplace()
}
