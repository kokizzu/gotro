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
		Limit  uint32 `json:"limit" form:"limit" query:"limit" long:"limit" msg:"limit"`
		Offset uint32 `json:"offset" form:"offset" query:"offset" long:"offset" msg:"offset"`
	}
	StoreProducts_Out struct {
		ResponseCommon
		Limit    uint32 `json:"limit" form:"limit" query:"limit" long:"limit" msg:"limit"`
		Offset   uint32 `json:"offset" form:"offset" query:"offset" long:"offset" msg:"offset"`
		Total    uint32 `json:"total" form:"total" query:"total" long:"total" msg:"total"`
		Products []*rqStore.Products
	}
)

const StoreProducts_Url = `/StoreProducts`

func (d *Domain) StoreProducts(in *StoreProducts_In) (out StoreProducts_Out) {
	products := rqStore.NewProducts(d.Taran)
	out.Limit = in.Limit
	out.Offset = in.Offset
	out.Total = uint32(products.Total())
	out.Products = products.FindOffsetLimit(in.Offset, in.Limit)
	return
}

type (
	StoreCartItemsAdd_In struct {
		RequestCommon
		ProductId  uint64
		ProductQty int64
	}
	StoreCartItemsAdd_Out struct {
		ResponseCommon
		CartItems []*rqStore.CartItems
	}
)

const StoreCartItemsAdd_Url = `/StoreCartItemsAdd`

func (d *Domain) StoreCartItemsAdd(in *StoreCartItemsAdd_In) (out StoreCartItemsAdd_Out) {
	// TODO: uncomment check user login
	//if d.mustLogin(in.SessionToken, in.UserAgent, &in.RequestCommon) {
	//	return out.SetError(403,)
	//}
	// TODO: continue this

	return
}

type (
	StoreInvoice_In struct {
		RequestCommon
	}
	StoreInvoice_Out struct {
		ResponseCommon
		CartItems []rqStore.CartItems
		Invoice   rqStore.Invoices
	}
)

const StoreInvoice_Url = `/StoreInvoice`

func (d *Domain) StoreInvoice(in *StoreInvoice_In) (out StoreInvoice_Out) {
	// TODO: continue this
	return
}
