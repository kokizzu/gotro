package domain

import "github.com/kokizzu/gotro/W2/example/model/mStore/rqStore"

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file template.go
//go:generate replacer 'Id" form' 'Id,string" form' type template.go
//go:generate replacer 'json:"id"' 'json:id,string" form' type template.go
//go:generate replacer 'By" form' 'By,string" form' type template.go
// go:generate msgp -tests=false -file template.go -o template__MSG.GEN.go

type (
	StoreProducts_In struct {
		RequestCommon
	}
	StoreProducts_Out struct {
		ResponseCommon
		Products []rqStore.Products
	}
)

const StoreProducts_Url = `/StoreProducts`

func (d *Domain) StoreProducts(in *StoreProducts_In) (out StoreProducts_Out) {
	// TODO: continue this
	return
}

type (
	StoreCartItemsAdd_In struct {
		RequestCommon
	}
	StoreCartItemsAdd_Out struct {
		ResponseCommon
	}
)

const StoreCartItemsAdd_Url = `/StoreCartItemsAdd`

func (d *Domain) StoreCartItemsAdd(in *StoreCartItemsAdd_In) (out StoreCartItemsAdd_Out) {
	// TODO: continue this
	return
}

type (
	StoreInvoice_In struct {
		RequestCommon
	}
	StoreInvoice_Out struct {
		ResponseCommon
	}
)

const StoreInvoice_Url = `/StoreInvoice`

func (d *Domain) StoreInvoice(in *StoreInvoice_In) (out StoreInvoice_Out) {
	// TODO: continue this
	return
}
