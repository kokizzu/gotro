package domain

import (
	"github.com/kokizzu/gotro/F"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/W2/example/model/mStore/rqStore"
	"github.com/kokizzu/gotro/W2/example/model/mStore/wcStore"
	"github.com/kokizzu/id64"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file store.go
//go:generate replacer 'Id" form' 'Id,string" form' type store.go
//go:generate replacer 'json:"id"' 'json:"id,string"' type store.go
//go:generate replacer 'By" form' 'By,string" form' type store.go
// go:generate msgp -tests=false -file store.go -o store__MSG.GEN.go
//go:generate farify doublequote --file store.go

type (
	StoreProducts_In struct {
		RequestCommon
		Limit  uint32 `json:"limit" form:"limit" query:"limit" long:"limit" msg:"limit"`
		Offset uint32 `json:"offset" form:"offset" query:"offset" long:"offset" msg:"offset"`
	}
	StoreProducts_Out struct {
		ResponseCommon
		Limit    uint32              `json:"limit" form:"limit" query:"limit" long:"limit" msg:"limit"`
		Offset   uint32              `json:"offset" form:"offset" query:"offset" long:"offset" msg:"offset"`
		Total    uint32              `json:"total" form:"total" query:"total" long:"total" msg:"total"`
		Products []*rqStore.Products `json:"products" form:"products" query:"products" long:"products" msg:"products"`
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
		ProductId uint64 `json:"productId,string" form:"productId" query:"productId" long:"productId" msg:"productId"`
		DeltaQty  int64  `json:"deltaQty" form:"deltaQty" query:"deltaQty" long:"deltaQty" msg:"deltaQty"`
		// -n remove from cart, +n add to cart
	}
	StoreCartItemsAdd_Out struct {
		ResponseCommon
		CartItems  []*rqStore.CartItems `json:"cartItems" form:"cartItems" query:"cartItems" long:"cartItems" msg:"cartItems"`
		Total      uint32               `json:"total" form:"total" query:"total" long:"total" msg:"total"`
		IsOverflow bool                 `json:"isOverflow" form:"isOverflow" query:"isOverflow" long:"isOverflow" msg:"isOverflow"`
	}
)

const StoreCartItemsAdd_Url = `/StoreCartItemsAdd`

func (d *Domain) StoreCartItemsAdd(in *StoreCartItemsAdd_In) (out StoreCartItemsAdd_Out) {
	sess := d.mustLogin(in.SessionToken, in.UserAgent, &out.ResponseCommon)
	if sess == nil {
		out.SetError(403, `must login`)
		return
	}
	cartItem := wcStore.NewCartItemsMutator(d.Taran)
	cartItem.SetProductId(in.ProductId)
	cartItem.SetOwnerId(sess.UserId)
	// InvoiceId = 0 not yet purchased

	product := rqStore.NewProducts(d.Taran)
	product.Id = in.ProductId
	if !product.FindById() {
		if in.DeltaQty < 0 { // error unless removing from cart
			out.SetError(404, `product not found`)
			return
		}
	}

	inv := int64(product.InventoryQty)
	if !cartItem.FindByOwnerIdInvoiceIdProductId() {
		if in.DeltaQty < 0 {
			out.SetError(404, `cart item not found`)
			return
		}

		newQty := I.Max(0, in.DeltaQty)
		newQty = I.Min(newQty, inv)
		cartItem.SetQty(newQty)
		cartItem.Id = id64.UID()
		if !cartItem.DoInsert() {
			out.SetError(500, `failed insert to cart`)
			return
		}
	} else {
		if cartItem.Qty >= inv && in.DeltaQty > 0 {
			out.SetError(400, `cannot add more`)
			return
		}

		if cartItem.Qty <= 0 && in.DeltaQty < 0 {
			out.SetError(400, `cannot remove more`)
			return
		}
		newQty := I.Max(0, cartItem.Qty+in.DeltaQty)
		cartItem.SetQty(newQty)
		if !cartItem.DoUpdateById() {
			out.SetError(500, `failed add/remove item on cart`)
			return
		}
	}

	out.IsOverflow = cartItem.Qty >= inv
	out.CartItems, out.Total = cartItem.FindByOwnerIdInvoiceId()

	return
}

type (
	StoreInvoice_In struct {
		RequestCommon
		InvoiceId   uint64 `json:"invoiceId,string" form:"invoiceId" query:"invoiceId" long:"invoiceId" msg:"invoiceId"`
		Recalculate bool   `json:"recalculate" form:"recalculate" query:"recalculate" long:"recalculate" msg:"recalculate"`
		DoPurchase  bool   `json:"doPurchase" form:"doPurchase" query:"doPurchase" long:"doPurchase" msg:"doPurchase"`
	}
	StoreInvoice_Out struct {
		ResponseCommon
		CartItems []*rqStore.CartItems `json:"cartItems" form:"cartItems" query:"cartItems" long:"cartItems" msg:"cartItems"`
		Invoice   rqStore.Invoices     `json:"invoice" form:"invoice" query:"invoice" long:"invoice" msg:"invoice"`
	}
)

const StoreInvoice_Url = `/StoreInvoice`

func (d *Domain) StoreInvoice(in *StoreInvoice_In) (out StoreInvoice_Out) {
	sess := d.mustLogin(in.SessionToken, in.UserAgent, &out.ResponseCommon)
	if sess == nil {
		out.SetError(403, `must login`)
		return
	}

	invoice := wcStore.NewInvoicesMutator(d.Taran)
	invoice.Id = in.InvoiceId
	if in.InvoiceId == 0 && !invoice.FindById() {
		if in.DoPurchase {
			in.InvoiceId = id64.UID()
			// TODO: add rollback just in case saving interrupted
		}
	}

	cartItem := rqStore.NewCartItems(d.Taran)
	cartItem.OwnerId = sess.UserId
	cartItem.InvoiceId = in.InvoiceId // invoiceId = 0 not yet purchased

	promo := rqStore.NewPromos(d.Taran)
	promos := promo.FindActive()
	promoByProductId := map[uint64]*rqStore.Promos{}
	// ^ assuming no 2 promo at the same time for same product
	for _, promo := range promos {
		promoByProductId[promo.ProductId] = promo
	}

	// free product map
	type FreeInfo struct {
		Count int64  `json:"count" form:"count" query:"count" long:"count" msg:"count"`
		Label string `json:"label" form:"label" query:"label" long:"label" msg:"label"`
	}
	freeProductsMap := map[uint64]*FreeInfo{}
	addFreeProduct := func(productId uint64, count int64, label string) {
		free, ok := freeProductsMap[productId]
		if ok {
			free.Count += count
			free.Label += label + "\n"
		} else {
			free = &FreeInfo{count, label + "\n"}
		}
		freeProductsMap[productId] = free
	}

	// fetch product, TODO: change to only fetch by id that are in cart
	product := rqStore.NewProducts(d.Taran)
	products := product.FindOffsetLimit(0, 100)
	productIdMap := map[uint64]*rqStore.Products{}
	for _, product := range products {
		productIdMap[product.Id] = product
	}

	out.CartItems, _ = cartItem.FindByOwnerIdInvoiceId()
	cartItemIdMap := map[uint64]*rqStore.CartItems{}
	for _, ci := range out.CartItems {
		cartItemIdMap[ci.ProductId] = ci
	}

	if in.Recalculate || in.DoPurchase {
		for _, ci := range out.CartItems {
			ci.InvoiceId = in.InvoiceId
			ci.Info = ``
			product := productIdMap[ci.ProductId]
			if product == nil {
				// if product gone, set qty to 0
				ci.Qty = 0
				ci.Info = "product no longer exists\n"
			} else {
				ci.PriceCopy = int64(product.Price)
				ci.NameCopy = product.Name
				inv := int64(product.InventoryQty)
				if ci.Qty > inv {
					ci.Info = "qty in cart more than available stock\n"
					// make sure next purchase doesn't overflow
					ci.Qty = I.Min(ci.Qty, inv)
				}
			}

			ci.SubTotal = ci.PriceCopy * ci.Qty // negative = refund

			// apply promo
			promo, ok := promoByProductId[ci.ProductId]
			if ok {
				minCount := int64(promo.ProductCount)
				if ci.Qty >= minCount {
					multiplier := ci.Qty / minCount
					if promo.FreeProductId > 0 { // got other product for free
						addFreeProduct(promo.FreeProductId, multiplier, `got 1 free (total: `+I.ToS(multiplier)+`) every purchase of `+I.ToS(minCount)+` `+product.Name)
					} else if promo.DiscountPercent > 0 {
						orig := ci.SubTotal
						ci.SubTotal = int64(float64(ci.SubTotal) * (100 - promo.DiscountPercent) / 100)
						ci.Discount = uint64(orig - ci.SubTotal)
						ci.Info += `discount ` + F.ToS(promo.DiscountPercent) + `% for ` + I.ToS(minCount) + " purchase\n"
					} else if promo.DiscountCount > 0 {
						// eg. buy 3, 3rd one is discount, then buy 6, 3rd and 6th is discount
						orig := ci.SubTotal
						ci.SubTotal = ci.PriceCopy * (ci.Qty - int64(promo.DiscountCount)*multiplier)
						ci.Discount = uint64(orig - ci.SubTotal)
						ci.Info += `discount ` + I.UToS(promo.DiscountCount) + ` for every ` + I.ToS(minCount) + " purchase\n"
					}
				}
			}

			if in.DoPurchase {
				cartItem := wcStore.NewCartItemsMutator(d.Taran)
				ci.Adapter = d.Taran
				cartItem.CartItems = *ci
				cartItem.DoUpdateById()
				ci.Adapter = nil
			}
		}

		// add free item
		for productId, freeInfo := range freeProductsMap {
			if product == nil { // ignore if product doesn't exists
				continue
			}
			ci, ok := cartItemIdMap[productId]
			product := productIdMap[productId]
			if ok {
				ci.Qty += freeInfo.Count
				ci.Info += freeInfo.Label
				ci.Discount += product.Price * uint64(freeInfo.Count)
			} else {
				ci = &rqStore.CartItems{
					Id:        id64.UID(),
					ProductId: productId,
					OwnerId:   sess.UserId,
					Qty:       freeInfo.Count,
					Info:      freeInfo.Label,
					InvoiceId: out.Invoice.Id,
					PriceCopy: int64(product.Price),
					NameCopy:  product.Name,
					Discount:  product.Price * uint64(freeInfo.Count),
				}
				out.CartItems = append(out.CartItems, ci)
				if in.DoPurchase {
					cartItem := wcStore.NewCartItemsMutator(d.Taran)
					ci.Adapter = d.Taran
					cartItem.CartItems = *ci
					cartItem.DoUpdateById()
					ci.Adapter = nil
				}
			}
			inv := int64(product.InventoryQty)
			if ci.Qty > inv {
				missing := ci.Qty - inv
				ci.Info += `but we don't have enough free item in inventory (missing: ` + I.ToS(missing) + ")\n"
				ci.Qty = I.Min(ci.Qty, inv)
				ci.Discount -= product.Price * uint64(missing)
			}
		}
		out.Invoice = rqStore.Invoices{}

		total := int64(0)
		for _, ci := range out.CartItems {
			total += ci.SubTotal
			out.Invoice.TotalDiscount += ci.Discount
		}
		out.Invoice.TotalPaid = uint64(total)
	}

	return
}

// TODO: decrease stock after accepted by seller
