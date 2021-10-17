package rqStore

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

import (
	"github.com/kokizzu/gotro/W2/example/model/mStore"

	"github.com/tarantool/go-tarantool"

	"github.com/graphql-go/graphql"
	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/X"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file rqStore__ORM.GEN.go
//go:generate replacer 'Id" form' 'Id,string" form' type rqStore__ORM.GEN.go
//go:generate replacer 'json:"id"' 'json:"id,string"' type rqStore__ORM.GEN.go
//go:generate replacer 'By" form' 'By,string" form' type rqStore__ORM.GEN.go
// go:generate msgp -tests=false -file rqStore__ORM.GEN.go -o rqStore__MSG.GEN.go

type CartItems struct {
	Adapter    *Tt.Adapter `json:"-" msg:"-" query:"-" form:"-" long:"adapter"`
	Id         uint64      `json:"id,string" form:"id" query:"id" long:"id" msg:"id"`
	CreatedAt  int64       `json:"createdAt" form:"createdAt" query:"createdAt" long:"createdAt" msg:"createdAt"`
	CreatedBy  uint64      `json:"createdBy,string" form:"createdBy" query:"createdBy" long:"createdBy" msg:"createdBy"`
	UpdatedAt  int64       `json:"updatedAt" form:"updatedAt" query:"updatedAt" long:"updatedAt" msg:"updatedAt"`
	UpdatedBy  uint64      `json:"updatedBy,string" form:"updatedBy" query:"updatedBy" long:"updatedBy" msg:"updatedBy"`
	DeletedAt  int64       `json:"deletedAt" form:"deletedAt" query:"deletedAt" long:"deletedAt" msg:"deletedAt"`
	DeletedBy  uint64      `json:"deletedBy,string" form:"deletedBy" query:"deletedBy" long:"deletedBy" msg:"deletedBy"`
	IsDeleted  bool        `json:"isDeleted" form:"isDeleted" query:"isDeleted" long:"isDeleted" msg:"isDeleted"`
	RestoredAt int64       `json:"restoredAt" form:"restoredAt" query:"restoredAt" long:"restoredAt" msg:"restoredAt"`
	RestoredBy uint64      `json:"restoredBy,string" form:"restoredBy" query:"restoredBy" long:"restoredBy" msg:"restoredBy"`
	OwnerId    uint64      `json:"ownerId,string" form:"ownerId" query:"ownerId" long:"ownerId" msg:"ownerId"`
	InvoiceId  uint64      `json:"invoiceId,string" form:"invoiceId" query:"invoiceId" long:"invoiceId" msg:"invoiceId"`
	ProductId  uint64      `json:"productId,string" form:"productId" query:"productId" long:"productId" msg:"productId"`
	NameCopy   string      `json:"nameCopy" form:"nameCopy" query:"nameCopy" long:"nameCopy" msg:"nameCopy"`
	PriceCopy  int64       `json:"priceCopy" form:"priceCopy" query:"priceCopy" long:"priceCopy" msg:"priceCopy"`
	Qty        int64       `json:"qty" form:"qty" query:"qty" long:"qty" msg:"qty"`
	Discount   uint64      `json:"discount" form:"discount" query:"discount" long:"discount" msg:"discount"`
	SubTotal   int64       `json:"subTotal" form:"subTotal" query:"subTotal" long:"subTotal" msg:"subTotal"`
	Info       string      `json:"info" form:"info" query:"info" long:"info" msg:"info"`
}

func NewCartItems(adapter *Tt.Adapter) *CartItems {
	return &CartItems{Adapter: adapter}
}

func (c *CartItems) SpaceName() string { //nolint:dupl false positive
	return string(mStore.TableCartItems)
}

func (c *CartItems) sqlTableName() string { //nolint:dupl false positive
	return `"cartItems"`
}

func (c *CartItems) UniqueIndexId() string { //nolint:dupl false positive
	return `id`
}

func (c *CartItems) FindById() bool { //nolint:dupl false positive
	res, err := c.Adapter.Select(c.SpaceName(), c.UniqueIndexId(), 0, 1, tarantool.IterEq, A.X{c.Id})
	if L.IsError(err, `CartItems.FindById failed: `+c.SpaceName()) {
		return false
	}
	rows := res.Tuples()
	if len(rows) == 1 {
		c.FromArray(rows[0])
		return true
	}
	return false
}

var GraphqlFieldCartItemsById = &graphql.Field{
	Type:        GraphqlTypeCartItems,
	Description: `list of CartItems`,
	Args: graphql.FieldConfigArgument{
		`Id`: &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
	},
}

func (g *CartItems) GraphqlFieldCartItemsByIdWithResolver() *graphql.Field {
	field := *GraphqlFieldCartItemsById
	field.Resolve = func(p graphql.ResolveParams) (interface{}, error) {
		q := g
		v, ok := p.Args[`id`]
		if !ok {
			v, _ = p.Args[`Id`]
		}
		q.Id = X.ToU(v)
		if q.FindById() {
			return q, nil
		}
		return nil, nil
	}
	return &field
}

func (c *CartItems) UniqueIndexOwnerIdProductIdInvoiceId() string { //nolint:dupl false positive
	return `ownerId__productId__invoiceId`
}

func (c *CartItems) FindByOwnerIdProductIdInvoiceId() bool { //nolint:dupl false positive
	res, err := c.Adapter.Select(c.SpaceName(), c.UniqueIndexOwnerIdProductIdInvoiceId(), 0, 1, tarantool.IterEq, A.X{c.OwnerId, c.ProductId, c.InvoiceId})
	if L.IsError(err, `CartItems.FindByOwnerIdProductIdInvoiceId failed: `+c.SpaceName()) {
		return false
	}
	rows := res.Tuples()
	if len(rows) == 1 {
		c.FromArray(rows[0])
		return true
	}
	return false
}

func (c *CartItems) sqlSelectAllFields() string { //nolint:dupl false positive
	return ` "id"
	, "createdAt"
	, "createdBy"
	, "updatedAt"
	, "updatedBy"
	, "deletedAt"
	, "deletedBy"
	, "isDeleted"
	, "restoredAt"
	, "restoredBy"
	, "ownerId"
	, "invoiceId"
	, "productId"
	, "nameCopy"
	, "priceCopy"
	, "qty"
	, "discount"
	, "subTotal"
	, "info"
	`
}

func (c *CartItems) ToUpdateArray() A.X { //nolint:dupl false positive
	return A.X{
		A.X{`=`, 0, c.Id},
		A.X{`=`, 1, c.CreatedAt},
		A.X{`=`, 2, c.CreatedBy},
		A.X{`=`, 3, c.UpdatedAt},
		A.X{`=`, 4, c.UpdatedBy},
		A.X{`=`, 5, c.DeletedAt},
		A.X{`=`, 6, c.DeletedBy},
		A.X{`=`, 7, c.IsDeleted},
		A.X{`=`, 8, c.RestoredAt},
		A.X{`=`, 9, c.RestoredBy},
		A.X{`=`, 10, c.OwnerId},
		A.X{`=`, 11, c.InvoiceId},
		A.X{`=`, 12, c.ProductId},
		A.X{`=`, 13, c.NameCopy},
		A.X{`=`, 14, c.PriceCopy},
		A.X{`=`, 15, c.Qty},
		A.X{`=`, 16, c.Discount},
		A.X{`=`, 17, c.SubTotal},
		A.X{`=`, 18, c.Info},
	}
}

func (c *CartItems) IdxId() int { //nolint:dupl false positive
	return 0
}

func (c *CartItems) sqlId() string { //nolint:dupl false positive
	return `"id"`
}

func (c *CartItems) IdxCreatedAt() int { //nolint:dupl false positive
	return 1
}

func (c *CartItems) sqlCreatedAt() string { //nolint:dupl false positive
	return `"createdAt"`
}

func (c *CartItems) IdxCreatedBy() int { //nolint:dupl false positive
	return 2
}

func (c *CartItems) sqlCreatedBy() string { //nolint:dupl false positive
	return `"createdBy"`
}

func (c *CartItems) IdxUpdatedAt() int { //nolint:dupl false positive
	return 3
}

func (c *CartItems) sqlUpdatedAt() string { //nolint:dupl false positive
	return `"updatedAt"`
}

func (c *CartItems) IdxUpdatedBy() int { //nolint:dupl false positive
	return 4
}

func (c *CartItems) sqlUpdatedBy() string { //nolint:dupl false positive
	return `"updatedBy"`
}

func (c *CartItems) IdxDeletedAt() int { //nolint:dupl false positive
	return 5
}

func (c *CartItems) sqlDeletedAt() string { //nolint:dupl false positive
	return `"deletedAt"`
}

func (c *CartItems) IdxDeletedBy() int { //nolint:dupl false positive
	return 6
}

func (c *CartItems) sqlDeletedBy() string { //nolint:dupl false positive
	return `"deletedBy"`
}

func (c *CartItems) IdxIsDeleted() int { //nolint:dupl false positive
	return 7
}

func (c *CartItems) sqlIsDeleted() string { //nolint:dupl false positive
	return `"isDeleted"`
}

func (c *CartItems) IdxRestoredAt() int { //nolint:dupl false positive
	return 8
}

func (c *CartItems) sqlRestoredAt() string { //nolint:dupl false positive
	return `"restoredAt"`
}

func (c *CartItems) IdxRestoredBy() int { //nolint:dupl false positive
	return 9
}

func (c *CartItems) sqlRestoredBy() string { //nolint:dupl false positive
	return `"restoredBy"`
}

func (c *CartItems) IdxOwnerId() int { //nolint:dupl false positive
	return 10
}

func (c *CartItems) sqlOwnerId() string { //nolint:dupl false positive
	return `"ownerId"`
}

func (c *CartItems) IdxInvoiceId() int { //nolint:dupl false positive
	return 11
}

func (c *CartItems) sqlInvoiceId() string { //nolint:dupl false positive
	return `"invoiceId"`
}

func (c *CartItems) IdxProductId() int { //nolint:dupl false positive
	return 12
}

func (c *CartItems) sqlProductId() string { //nolint:dupl false positive
	return `"productId"`
}

func (c *CartItems) IdxNameCopy() int { //nolint:dupl false positive
	return 13
}

func (c *CartItems) sqlNameCopy() string { //nolint:dupl false positive
	return `"nameCopy"`
}

func (c *CartItems) IdxPriceCopy() int { //nolint:dupl false positive
	return 14
}

func (c *CartItems) sqlPriceCopy() string { //nolint:dupl false positive
	return `"priceCopy"`
}

func (c *CartItems) IdxQty() int { //nolint:dupl false positive
	return 15
}

func (c *CartItems) sqlQty() string { //nolint:dupl false positive
	return `"qty"`
}

func (c *CartItems) IdxDiscount() int { //nolint:dupl false positive
	return 16
}

func (c *CartItems) sqlDiscount() string { //nolint:dupl false positive
	return `"discount"`
}

func (c *CartItems) IdxSubTotal() int { //nolint:dupl false positive
	return 17
}

func (c *CartItems) sqlSubTotal() string { //nolint:dupl false positive
	return `"subTotal"`
}

func (c *CartItems) IdxInfo() int { //nolint:dupl false positive
	return 18
}

func (c *CartItems) sqlInfo() string { //nolint:dupl false positive
	return `"info"`
}

func (c *CartItems) ToArray() A.X { //nolint:dupl false positive
	return A.X{
		c.Id,         // 0
		c.CreatedAt,  // 1
		c.CreatedBy,  // 2
		c.UpdatedAt,  // 3
		c.UpdatedBy,  // 4
		c.DeletedAt,  // 5
		c.DeletedBy,  // 6
		c.IsDeleted,  // 7
		c.RestoredAt, // 8
		c.RestoredBy, // 9
		c.OwnerId,    // 10
		c.InvoiceId,  // 11
		c.ProductId,  // 12
		c.NameCopy,   // 13
		c.PriceCopy,  // 14
		c.Qty,        // 15
		c.Discount,   // 16
		c.SubTotal,   // 17
		c.Info,       // 18
	}
}

func (c *CartItems) FromArray(a A.X) *CartItems { //nolint:dupl false positive
	c.Id = X.ToU(a[0])
	c.CreatedAt = X.ToI(a[1])
	c.CreatedBy = X.ToU(a[2])
	c.UpdatedAt = X.ToI(a[3])
	c.UpdatedBy = X.ToU(a[4])
	c.DeletedAt = X.ToI(a[5])
	c.DeletedBy = X.ToU(a[6])
	c.IsDeleted = X.ToBool(a[7])
	c.RestoredAt = X.ToI(a[8])
	c.RestoredBy = X.ToU(a[9])
	c.OwnerId = X.ToU(a[10])
	c.InvoiceId = X.ToU(a[11])
	c.ProductId = X.ToU(a[12])
	c.NameCopy = X.ToS(a[13])
	c.PriceCopy = X.ToI(a[14])
	c.Qty = X.ToI(a[15])
	c.Discount = X.ToU(a[16])
	c.SubTotal = X.ToI(a[17])
	c.Info = X.ToS(a[18])
	return c
}

func (c *CartItems) Total() int64 { //nolint:dupl false positive
	rows := c.Adapter.CallBoxSpace(c.SpaceName()+`:count`, A.X{})
	if len(rows) > 0 && len(rows[0]) > 0 {
		return X.ToI(rows[0][0])
	}
	return 0
}

var GraphqlTypeCartItems = graphql.NewObject(
	graphql.ObjectConfig{
		Name: `cartItems`,
		Fields: graphql.Fields{
			`id`: &graphql.Field{
				Type: graphql.Int,
			},
			`createdAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`createdBy`: &graphql.Field{
				Type: graphql.ID,
			},
			`updatedAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`updatedBy`: &graphql.Field{
				Type: graphql.ID,
			},
			`deletedAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`deletedBy`: &graphql.Field{
				Type: graphql.ID,
			},
			`isDeleted`: &graphql.Field{
				Type: graphql.Boolean,
			},
			`restoredAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`restoredBy`: &graphql.Field{
				Type: graphql.ID,
			},
			`ownerId`: &graphql.Field{
				Type: graphql.ID,
			},
			`invoiceId`: &graphql.Field{
				Type: graphql.ID,
			},
			`productId`: &graphql.Field{
				Type: graphql.ID,
			},
			`nameCopy`: &graphql.Field{
				Type: graphql.String,
			},
			`priceCopy`: &graphql.Field{
				Type: graphql.Int,
			},
			`qty`: &graphql.Field{
				Type: graphql.Int,
			},
			`discount`: &graphql.Field{
				Type: graphql.Int,
			},
			`subTotal`: &graphql.Field{
				Type: graphql.Int,
			},
			`info`: &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

type Invoices struct {
	Adapter        *Tt.Adapter `json:"-" msg:"-" query:"-" form:"-" long:"adapter"`
	Id             uint64      `json:"id,string" form:"id" query:"id" long:"id" msg:"id"`
	CreatedAt      int64       `json:"createdAt" form:"createdAt" query:"createdAt" long:"createdAt" msg:"createdAt"`
	CreatedBy      uint64      `json:"createdBy,string" form:"createdBy" query:"createdBy" long:"createdBy" msg:"createdBy"`
	UpdatedAt      int64       `json:"updatedAt" form:"updatedAt" query:"updatedAt" long:"updatedAt" msg:"updatedAt"`
	UpdatedBy      uint64      `json:"updatedBy,string" form:"updatedBy" query:"updatedBy" long:"updatedBy" msg:"updatedBy"`
	DeletedAt      int64       `json:"deletedAt" form:"deletedAt" query:"deletedAt" long:"deletedAt" msg:"deletedAt"`
	DeletedBy      uint64      `json:"deletedBy,string" form:"deletedBy" query:"deletedBy" long:"deletedBy" msg:"deletedBy"`
	IsDeleted      bool        `json:"isDeleted" form:"isDeleted" query:"isDeleted" long:"isDeleted" msg:"isDeleted"`
	RestoredAt     int64       `json:"restoredAt" form:"restoredAt" query:"restoredAt" long:"restoredAt" msg:"restoredAt"`
	RestoredBy     uint64      `json:"restoredBy,string" form:"restoredBy" query:"restoredBy" long:"restoredBy" msg:"restoredBy"`
	OwnerId        uint64      `json:"ownerId,string" form:"ownerId" query:"ownerId" long:"ownerId" msg:"ownerId"`
	TotalWeight    uint64      `json:"totalWeight" form:"totalWeight" query:"totalWeight" long:"totalWeight" msg:"totalWeight"`
	TotalPrice     uint64      `json:"totalPrice" form:"totalPrice" query:"totalPrice" long:"totalPrice" msg:"totalPrice"`
	TotalDiscount  uint64      `json:"totalDiscount" form:"totalDiscount" query:"totalDiscount" long:"totalDiscount" msg:"totalDiscount"`
	DeliveryMethod uint64      `json:"deliveryMethod" form:"deliveryMethod" query:"deliveryMethod" long:"deliveryMethod" msg:"deliveryMethod"`
	DeliveryPrice  uint64      `json:"deliveryPrice" form:"deliveryPrice" query:"deliveryPrice" long:"deliveryPrice" msg:"deliveryPrice"`
	TotalPaid      uint64      `json:"totalPaid" form:"totalPaid" query:"totalPaid" long:"totalPaid" msg:"totalPaid"`
	PaidAt         uint64      `json:"paidAt" form:"paidAt" query:"paidAt" long:"paidAt" msg:"paidAt"`
	PaymentMethod  uint64      `json:"paymentMethod" form:"paymentMethod" query:"paymentMethod" long:"paymentMethod" msg:"paymentMethod"`
	DeadlineAt     uint64      `json:"deadlineAt" form:"deadlineAt" query:"deadlineAt" long:"deadlineAt" msg:"deadlineAt"`
	PromoRuleIds   string      `json:"promoRuleIds" form:"promoRuleIds" query:"promoRuleIds" long:"promoRuleIds" msg:"promoRuleIds"`
}

func NewInvoices(adapter *Tt.Adapter) *Invoices {
	return &Invoices{Adapter: adapter}
}

func (i *Invoices) SpaceName() string { //nolint:dupl false positive
	return string(mStore.TableInvoices)
}

func (i *Invoices) sqlTableName() string { //nolint:dupl false positive
	return `"invoices"`
}

func (i *Invoices) UniqueIndexId() string { //nolint:dupl false positive
	return `id`
}

func (i *Invoices) FindById() bool { //nolint:dupl false positive
	res, err := i.Adapter.Select(i.SpaceName(), i.UniqueIndexId(), 0, 1, tarantool.IterEq, A.X{i.Id})
	if L.IsError(err, `Invoices.FindById failed: `+i.SpaceName()) {
		return false
	}
	rows := res.Tuples()
	if len(rows) == 1 {
		i.FromArray(rows[0])
		return true
	}
	return false
}

var GraphqlFieldInvoicesById = &graphql.Field{
	Type:        GraphqlTypeInvoices,
	Description: `list of Invoices`,
	Args: graphql.FieldConfigArgument{
		`Id`: &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
	},
}

func (g *Invoices) GraphqlFieldInvoicesByIdWithResolver() *graphql.Field {
	field := *GraphqlFieldInvoicesById
	field.Resolve = func(p graphql.ResolveParams) (interface{}, error) {
		q := g
		v, ok := p.Args[`id`]
		if !ok {
			v, _ = p.Args[`Id`]
		}
		q.Id = X.ToU(v)
		if q.FindById() {
			return q, nil
		}
		return nil, nil
	}
	return &field
}

func (i *Invoices) sqlSelectAllFields() string { //nolint:dupl false positive
	return ` "id"
	, "createdAt"
	, "createdBy"
	, "updatedAt"
	, "updatedBy"
	, "deletedAt"
	, "deletedBy"
	, "isDeleted"
	, "restoredAt"
	, "restoredBy"
	, "ownerId"
	, "totalWeight"
	, "totalPrice"
	, "totalDiscount"
	, "deliveryMethod"
	, "deliveryPrice"
	, "totalPaid"
	, "paidAt"
	, "paymentMethod"
	, "deadlineAt"
	, "promoRuleIds"
	`
}

func (i *Invoices) ToUpdateArray() A.X { //nolint:dupl false positive
	return A.X{
		A.X{`=`, 0, i.Id},
		A.X{`=`, 1, i.CreatedAt},
		A.X{`=`, 2, i.CreatedBy},
		A.X{`=`, 3, i.UpdatedAt},
		A.X{`=`, 4, i.UpdatedBy},
		A.X{`=`, 5, i.DeletedAt},
		A.X{`=`, 6, i.DeletedBy},
		A.X{`=`, 7, i.IsDeleted},
		A.X{`=`, 8, i.RestoredAt},
		A.X{`=`, 9, i.RestoredBy},
		A.X{`=`, 10, i.OwnerId},
		A.X{`=`, 11, i.TotalWeight},
		A.X{`=`, 12, i.TotalPrice},
		A.X{`=`, 13, i.TotalDiscount},
		A.X{`=`, 14, i.DeliveryMethod},
		A.X{`=`, 15, i.DeliveryPrice},
		A.X{`=`, 16, i.TotalPaid},
		A.X{`=`, 17, i.PaidAt},
		A.X{`=`, 18, i.PaymentMethod},
		A.X{`=`, 19, i.DeadlineAt},
		A.X{`=`, 20, i.PromoRuleIds},
	}
}

func (i *Invoices) IdxId() int { //nolint:dupl false positive
	return 0
}

func (i *Invoices) sqlId() string { //nolint:dupl false positive
	return `"id"`
}

func (i *Invoices) IdxCreatedAt() int { //nolint:dupl false positive
	return 1
}

func (i *Invoices) sqlCreatedAt() string { //nolint:dupl false positive
	return `"createdAt"`
}

func (i *Invoices) IdxCreatedBy() int { //nolint:dupl false positive
	return 2
}

func (i *Invoices) sqlCreatedBy() string { //nolint:dupl false positive
	return `"createdBy"`
}

func (i *Invoices) IdxUpdatedAt() int { //nolint:dupl false positive
	return 3
}

func (i *Invoices) sqlUpdatedAt() string { //nolint:dupl false positive
	return `"updatedAt"`
}

func (i *Invoices) IdxUpdatedBy() int { //nolint:dupl false positive
	return 4
}

func (i *Invoices) sqlUpdatedBy() string { //nolint:dupl false positive
	return `"updatedBy"`
}

func (i *Invoices) IdxDeletedAt() int { //nolint:dupl false positive
	return 5
}

func (i *Invoices) sqlDeletedAt() string { //nolint:dupl false positive
	return `"deletedAt"`
}

func (i *Invoices) IdxDeletedBy() int { //nolint:dupl false positive
	return 6
}

func (i *Invoices) sqlDeletedBy() string { //nolint:dupl false positive
	return `"deletedBy"`
}

func (i *Invoices) IdxIsDeleted() int { //nolint:dupl false positive
	return 7
}

func (i *Invoices) sqlIsDeleted() string { //nolint:dupl false positive
	return `"isDeleted"`
}

func (i *Invoices) IdxRestoredAt() int { //nolint:dupl false positive
	return 8
}

func (i *Invoices) sqlRestoredAt() string { //nolint:dupl false positive
	return `"restoredAt"`
}

func (i *Invoices) IdxRestoredBy() int { //nolint:dupl false positive
	return 9
}

func (i *Invoices) sqlRestoredBy() string { //nolint:dupl false positive
	return `"restoredBy"`
}

func (i *Invoices) IdxOwnerId() int { //nolint:dupl false positive
	return 10
}

func (i *Invoices) sqlOwnerId() string { //nolint:dupl false positive
	return `"ownerId"`
}

func (i *Invoices) IdxTotalWeight() int { //nolint:dupl false positive
	return 11
}

func (i *Invoices) sqlTotalWeight() string { //nolint:dupl false positive
	return `"totalWeight"`
}

func (i *Invoices) IdxTotalPrice() int { //nolint:dupl false positive
	return 12
}

func (i *Invoices) sqlTotalPrice() string { //nolint:dupl false positive
	return `"totalPrice"`
}

func (i *Invoices) IdxTotalDiscount() int { //nolint:dupl false positive
	return 13
}

func (i *Invoices) sqlTotalDiscount() string { //nolint:dupl false positive
	return `"totalDiscount"`
}

func (i *Invoices) IdxDeliveryMethod() int { //nolint:dupl false positive
	return 14
}

func (i *Invoices) sqlDeliveryMethod() string { //nolint:dupl false positive
	return `"deliveryMethod"`
}

func (i *Invoices) IdxDeliveryPrice() int { //nolint:dupl false positive
	return 15
}

func (i *Invoices) sqlDeliveryPrice() string { //nolint:dupl false positive
	return `"deliveryPrice"`
}

func (i *Invoices) IdxTotalPaid() int { //nolint:dupl false positive
	return 16
}

func (i *Invoices) sqlTotalPaid() string { //nolint:dupl false positive
	return `"totalPaid"`
}

func (i *Invoices) IdxPaidAt() int { //nolint:dupl false positive
	return 17
}

func (i *Invoices) sqlPaidAt() string { //nolint:dupl false positive
	return `"paidAt"`
}

func (i *Invoices) IdxPaymentMethod() int { //nolint:dupl false positive
	return 18
}

func (i *Invoices) sqlPaymentMethod() string { //nolint:dupl false positive
	return `"paymentMethod"`
}

func (i *Invoices) IdxDeadlineAt() int { //nolint:dupl false positive
	return 19
}

func (i *Invoices) sqlDeadlineAt() string { //nolint:dupl false positive
	return `"deadlineAt"`
}

func (i *Invoices) IdxPromoRuleIds() int { //nolint:dupl false positive
	return 20
}

func (i *Invoices) sqlPromoRuleIds() string { //nolint:dupl false positive
	return `"promoRuleIds"`
}

func (i *Invoices) ToArray() A.X { //nolint:dupl false positive
	return A.X{
		i.Id,             // 0
		i.CreatedAt,      // 1
		i.CreatedBy,      // 2
		i.UpdatedAt,      // 3
		i.UpdatedBy,      // 4
		i.DeletedAt,      // 5
		i.DeletedBy,      // 6
		i.IsDeleted,      // 7
		i.RestoredAt,     // 8
		i.RestoredBy,     // 9
		i.OwnerId,        // 10
		i.TotalWeight,    // 11
		i.TotalPrice,     // 12
		i.TotalDiscount,  // 13
		i.DeliveryMethod, // 14
		i.DeliveryPrice,  // 15
		i.TotalPaid,      // 16
		i.PaidAt,         // 17
		i.PaymentMethod,  // 18
		i.DeadlineAt,     // 19
		i.PromoRuleIds,   // 20
	}
}

func (i *Invoices) FromArray(a A.X) *Invoices { //nolint:dupl false positive
	i.Id = X.ToU(a[0])
	i.CreatedAt = X.ToI(a[1])
	i.CreatedBy = X.ToU(a[2])
	i.UpdatedAt = X.ToI(a[3])
	i.UpdatedBy = X.ToU(a[4])
	i.DeletedAt = X.ToI(a[5])
	i.DeletedBy = X.ToU(a[6])
	i.IsDeleted = X.ToBool(a[7])
	i.RestoredAt = X.ToI(a[8])
	i.RestoredBy = X.ToU(a[9])
	i.OwnerId = X.ToU(a[10])
	i.TotalWeight = X.ToU(a[11])
	i.TotalPrice = X.ToU(a[12])
	i.TotalDiscount = X.ToU(a[13])
	i.DeliveryMethod = X.ToU(a[14])
	i.DeliveryPrice = X.ToU(a[15])
	i.TotalPaid = X.ToU(a[16])
	i.PaidAt = X.ToU(a[17])
	i.PaymentMethod = X.ToU(a[18])
	i.DeadlineAt = X.ToU(a[19])
	i.PromoRuleIds = X.ToS(a[20])
	return i
}

func (i *Invoices) Total() int64 { //nolint:dupl false positive
	rows := i.Adapter.CallBoxSpace(i.SpaceName()+`:count`, A.X{})
	if len(rows) > 0 && len(rows[0]) > 0 {
		return X.ToI(rows[0][0])
	}
	return 0
}

var GraphqlTypeInvoices = graphql.NewObject(
	graphql.ObjectConfig{
		Name: `invoices`,
		Fields: graphql.Fields{
			`id`: &graphql.Field{
				Type: graphql.Int,
			},
			`createdAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`createdBy`: &graphql.Field{
				Type: graphql.ID,
			},
			`updatedAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`updatedBy`: &graphql.Field{
				Type: graphql.ID,
			},
			`deletedAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`deletedBy`: &graphql.Field{
				Type: graphql.ID,
			},
			`isDeleted`: &graphql.Field{
				Type: graphql.Boolean,
			},
			`restoredAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`restoredBy`: &graphql.Field{
				Type: graphql.ID,
			},
			`ownerId`: &graphql.Field{
				Type: graphql.ID,
			},
			`totalWeight`: &graphql.Field{
				Type: graphql.Int,
			},
			`totalPrice`: &graphql.Field{
				Type: graphql.Int,
			},
			`totalDiscount`: &graphql.Field{
				Type: graphql.Int,
			},
			`deliveryMethod`: &graphql.Field{
				Type: graphql.Int,
			},
			`deliveryPrice`: &graphql.Field{
				Type: graphql.Int,
			},
			`totalPaid`: &graphql.Field{
				Type: graphql.Int,
			},
			`paidAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`paymentMethod`: &graphql.Field{
				Type: graphql.Int,
			},
			`deadlineAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`promoRuleIds`: &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

type Products struct {
	Adapter      *Tt.Adapter `json:"-" msg:"-" query:"-" form:"-" long:"adapter"`
	Id           uint64      `json:"id,string" form:"id" query:"id" long:"id" msg:"id"`
	CreatedAt    int64       `json:"createdAt" form:"createdAt" query:"createdAt" long:"createdAt" msg:"createdAt"`
	CreatedBy    uint64      `json:"createdBy,string" form:"createdBy" query:"createdBy" long:"createdBy" msg:"createdBy"`
	UpdatedAt    int64       `json:"updatedAt" form:"updatedAt" query:"updatedAt" long:"updatedAt" msg:"updatedAt"`
	UpdatedBy    uint64      `json:"updatedBy,string" form:"updatedBy" query:"updatedBy" long:"updatedBy" msg:"updatedBy"`
	DeletedAt    int64       `json:"deletedAt" form:"deletedAt" query:"deletedAt" long:"deletedAt" msg:"deletedAt"`
	DeletedBy    uint64      `json:"deletedBy,string" form:"deletedBy" query:"deletedBy" long:"deletedBy" msg:"deletedBy"`
	IsDeleted    bool        `json:"isDeleted" form:"isDeleted" query:"isDeleted" long:"isDeleted" msg:"isDeleted"`
	RestoredAt   int64       `json:"restoredAt" form:"restoredAt" query:"restoredAt" long:"restoredAt" msg:"restoredAt"`
	RestoredBy   uint64      `json:"restoredBy,string" form:"restoredBy" query:"restoredBy" long:"restoredBy" msg:"restoredBy"`
	Sku          string      `json:"sku" form:"sku" query:"sku" long:"sku" msg:"sku"`
	Name         string      `json:"name" form:"name" query:"name" long:"name" msg:"name"`
	Price        uint64      `json:"price" form:"price" query:"price" long:"price" msg:"price"`
	InventoryQty uint64      `json:"inventoryQty" form:"inventoryQty" query:"inventoryQty" long:"inventoryQty" msg:"inventoryQty"`
	WeightGram   uint64      `json:"weightGram" form:"weightGram" query:"weightGram" long:"weightGram" msg:"weightGram"`
}

func NewProducts(adapter *Tt.Adapter) *Products {
	return &Products{Adapter: adapter}
}

func (p *Products) SpaceName() string { //nolint:dupl false positive
	return string(mStore.TableProducts)
}

func (p *Products) sqlTableName() string { //nolint:dupl false positive
	return `"products"`
}

func (p *Products) UniqueIndexId() string { //nolint:dupl false positive
	return `id`
}

func (p *Products) FindById() bool { //nolint:dupl false positive
	res, err := p.Adapter.Select(p.SpaceName(), p.UniqueIndexId(), 0, 1, tarantool.IterEq, A.X{p.Id})
	if L.IsError(err, `Products.FindById failed: `+p.SpaceName()) {
		return false
	}
	rows := res.Tuples()
	if len(rows) == 1 {
		p.FromArray(rows[0])
		return true
	}
	return false
}

var GraphqlFieldProductsById = &graphql.Field{
	Type:        GraphqlTypeProducts,
	Description: `list of Products`,
	Args: graphql.FieldConfigArgument{
		`Id`: &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
	},
}

func (g *Products) GraphqlFieldProductsByIdWithResolver() *graphql.Field {
	field := *GraphqlFieldProductsById
	field.Resolve = func(p graphql.ResolveParams) (interface{}, error) {
		q := g
		v, ok := p.Args[`id`]
		if !ok {
			v, _ = p.Args[`Id`]
		}
		q.Id = X.ToU(v)
		if q.FindById() {
			return q, nil
		}
		return nil, nil
	}
	return &field
}

func (p *Products) UniqueIndexSku() string { //nolint:dupl false positive
	return `sku`
}

func (p *Products) FindBySku() bool { //nolint:dupl false positive
	res, err := p.Adapter.Select(p.SpaceName(), p.UniqueIndexSku(), 0, 1, tarantool.IterEq, A.X{p.Sku})
	if L.IsError(err, `Products.FindBySku failed: `+p.SpaceName()) {
		return false
	}
	rows := res.Tuples()
	if len(rows) == 1 {
		p.FromArray(rows[0])
		return true
	}
	return false
}

func (p *Products) sqlSelectAllFields() string { //nolint:dupl false positive
	return ` "id"
	, "createdAt"
	, "createdBy"
	, "updatedAt"
	, "updatedBy"
	, "deletedAt"
	, "deletedBy"
	, "isDeleted"
	, "restoredAt"
	, "restoredBy"
	, "sku"
	, "name"
	, "price"
	, "inventoryQty"
	, "weightGram"
	`
}

func (p *Products) ToUpdateArray() A.X { //nolint:dupl false positive
	return A.X{
		A.X{`=`, 0, p.Id},
		A.X{`=`, 1, p.CreatedAt},
		A.X{`=`, 2, p.CreatedBy},
		A.X{`=`, 3, p.UpdatedAt},
		A.X{`=`, 4, p.UpdatedBy},
		A.X{`=`, 5, p.DeletedAt},
		A.X{`=`, 6, p.DeletedBy},
		A.X{`=`, 7, p.IsDeleted},
		A.X{`=`, 8, p.RestoredAt},
		A.X{`=`, 9, p.RestoredBy},
		A.X{`=`, 10, p.Sku},
		A.X{`=`, 11, p.Name},
		A.X{`=`, 12, p.Price},
		A.X{`=`, 13, p.InventoryQty},
		A.X{`=`, 14, p.WeightGram},
	}
}

func (p *Products) IdxId() int { //nolint:dupl false positive
	return 0
}

func (p *Products) sqlId() string { //nolint:dupl false positive
	return `"id"`
}

func (p *Products) IdxCreatedAt() int { //nolint:dupl false positive
	return 1
}

func (p *Products) sqlCreatedAt() string { //nolint:dupl false positive
	return `"createdAt"`
}

func (p *Products) IdxCreatedBy() int { //nolint:dupl false positive
	return 2
}

func (p *Products) sqlCreatedBy() string { //nolint:dupl false positive
	return `"createdBy"`
}

func (p *Products) IdxUpdatedAt() int { //nolint:dupl false positive
	return 3
}

func (p *Products) sqlUpdatedAt() string { //nolint:dupl false positive
	return `"updatedAt"`
}

func (p *Products) IdxUpdatedBy() int { //nolint:dupl false positive
	return 4
}

func (p *Products) sqlUpdatedBy() string { //nolint:dupl false positive
	return `"updatedBy"`
}

func (p *Products) IdxDeletedAt() int { //nolint:dupl false positive
	return 5
}

func (p *Products) sqlDeletedAt() string { //nolint:dupl false positive
	return `"deletedAt"`
}

func (p *Products) IdxDeletedBy() int { //nolint:dupl false positive
	return 6
}

func (p *Products) sqlDeletedBy() string { //nolint:dupl false positive
	return `"deletedBy"`
}

func (p *Products) IdxIsDeleted() int { //nolint:dupl false positive
	return 7
}

func (p *Products) sqlIsDeleted() string { //nolint:dupl false positive
	return `"isDeleted"`
}

func (p *Products) IdxRestoredAt() int { //nolint:dupl false positive
	return 8
}

func (p *Products) sqlRestoredAt() string { //nolint:dupl false positive
	return `"restoredAt"`
}

func (p *Products) IdxRestoredBy() int { //nolint:dupl false positive
	return 9
}

func (p *Products) sqlRestoredBy() string { //nolint:dupl false positive
	return `"restoredBy"`
}

func (p *Products) IdxSku() int { //nolint:dupl false positive
	return 10
}

func (p *Products) sqlSku() string { //nolint:dupl false positive
	return `"sku"`
}

func (p *Products) IdxName() int { //nolint:dupl false positive
	return 11
}

func (p *Products) sqlName() string { //nolint:dupl false positive
	return `"name"`
}

func (p *Products) IdxPrice() int { //nolint:dupl false positive
	return 12
}

func (p *Products) sqlPrice() string { //nolint:dupl false positive
	return `"price"`
}

func (p *Products) IdxInventoryQty() int { //nolint:dupl false positive
	return 13
}

func (p *Products) sqlInventoryQty() string { //nolint:dupl false positive
	return `"inventoryQty"`
}

func (p *Products) IdxWeightGram() int { //nolint:dupl false positive
	return 14
}

func (p *Products) sqlWeightGram() string { //nolint:dupl false positive
	return `"weightGram"`
}

func (p *Products) ToArray() A.X { //nolint:dupl false positive
	return A.X{
		p.Id,           // 0
		p.CreatedAt,    // 1
		p.CreatedBy,    // 2
		p.UpdatedAt,    // 3
		p.UpdatedBy,    // 4
		p.DeletedAt,    // 5
		p.DeletedBy,    // 6
		p.IsDeleted,    // 7
		p.RestoredAt,   // 8
		p.RestoredBy,   // 9
		p.Sku,          // 10
		p.Name,         // 11
		p.Price,        // 12
		p.InventoryQty, // 13
		p.WeightGram,   // 14
	}
}

func (p *Products) FromArray(a A.X) *Products { //nolint:dupl false positive
	p.Id = X.ToU(a[0])
	p.CreatedAt = X.ToI(a[1])
	p.CreatedBy = X.ToU(a[2])
	p.UpdatedAt = X.ToI(a[3])
	p.UpdatedBy = X.ToU(a[4])
	p.DeletedAt = X.ToI(a[5])
	p.DeletedBy = X.ToU(a[6])
	p.IsDeleted = X.ToBool(a[7])
	p.RestoredAt = X.ToI(a[8])
	p.RestoredBy = X.ToU(a[9])
	p.Sku = X.ToS(a[10])
	p.Name = X.ToS(a[11])
	p.Price = X.ToU(a[12])
	p.InventoryQty = X.ToU(a[13])
	p.WeightGram = X.ToU(a[14])
	return p
}

func (p *Products) Total() int64 { //nolint:dupl false positive
	rows := p.Adapter.CallBoxSpace(p.SpaceName()+`:count`, A.X{})
	if len(rows) > 0 && len(rows[0]) > 0 {
		return X.ToI(rows[0][0])
	}
	return 0
}

var GraphqlTypeProducts = graphql.NewObject(
	graphql.ObjectConfig{
		Name: `products`,
		Fields: graphql.Fields{
			`id`: &graphql.Field{
				Type: graphql.Int,
			},
			`createdAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`createdBy`: &graphql.Field{
				Type: graphql.ID,
			},
			`updatedAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`updatedBy`: &graphql.Field{
				Type: graphql.ID,
			},
			`deletedAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`deletedBy`: &graphql.Field{
				Type: graphql.ID,
			},
			`isDeleted`: &graphql.Field{
				Type: graphql.Boolean,
			},
			`restoredAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`restoredBy`: &graphql.Field{
				Type: graphql.ID,
			},
			`sku`: &graphql.Field{
				Type: graphql.String,
			},
			`name`: &graphql.Field{
				Type: graphql.String,
			},
			`price`: &graphql.Field{
				Type: graphql.Int,
			},
			`inventoryQty`: &graphql.Field{
				Type: graphql.Int,
			},
			`weightGram`: &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

type Promos struct {
	Adapter         *Tt.Adapter `json:"-" msg:"-" query:"-" form:"-" long:"adapter"`
	Id              uint64      `json:"id,string" form:"id" query:"id" long:"id" msg:"id"`
	CreatedAt       int64       `json:"createdAt" form:"createdAt" query:"createdAt" long:"createdAt" msg:"createdAt"`
	CreatedBy       uint64      `json:"createdBy,string" form:"createdBy" query:"createdBy" long:"createdBy" msg:"createdBy"`
	UpdatedAt       int64       `json:"updatedAt" form:"updatedAt" query:"updatedAt" long:"updatedAt" msg:"updatedAt"`
	UpdatedBy       uint64      `json:"updatedBy,string" form:"updatedBy" query:"updatedBy" long:"updatedBy" msg:"updatedBy"`
	DeletedAt       int64       `json:"deletedAt" form:"deletedAt" query:"deletedAt" long:"deletedAt" msg:"deletedAt"`
	DeletedBy       uint64      `json:"deletedBy,string" form:"deletedBy" query:"deletedBy" long:"deletedBy" msg:"deletedBy"`
	IsDeleted       bool        `json:"isDeleted" form:"isDeleted" query:"isDeleted" long:"isDeleted" msg:"isDeleted"`
	RestoredAt      int64       `json:"restoredAt" form:"restoredAt" query:"restoredAt" long:"restoredAt" msg:"restoredAt"`
	RestoredBy      uint64      `json:"restoredBy,string" form:"restoredBy" query:"restoredBy" long:"restoredBy" msg:"restoredBy"`
	StartAt         int64       `json:"startAt" form:"startAt" query:"startAt" long:"startAt" msg:"startAt"`
	EndAt           int64       `json:"endAt" form:"endAt" query:"endAt" long:"endAt" msg:"endAt"`
	ProductId       uint64      `json:"productId,string" form:"productId" query:"productId" long:"productId" msg:"productId"`
	ProductCount    uint64      `json:"productCount" form:"productCount" query:"productCount" long:"productCount" msg:"productCount"`
	FreeProductId   uint64      `json:"freeProductId,string" form:"freeProductId" query:"freeProductId" long:"freeProductId" msg:"freeProductId"`
	DiscountCount   uint64      `json:"discountCount" form:"discountCount" query:"discountCount" long:"discountCount" msg:"discountCount"`
	DiscountPercent float64     `json:"discountPercent" form:"discountPercent" query:"discountPercent" long:"discountPercent" msg:"discountPercent"`
}

func NewPromos(adapter *Tt.Adapter) *Promos {
	return &Promos{Adapter: adapter}
}

func (p *Promos) SpaceName() string { //nolint:dupl false positive
	return string(mStore.TablePromos)
}

func (p *Promos) sqlTableName() string { //nolint:dupl false positive
	return `"promos"`
}

func (p *Promos) UniqueIndexId() string { //nolint:dupl false positive
	return `id`
}

func (p *Promos) FindById() bool { //nolint:dupl false positive
	res, err := p.Adapter.Select(p.SpaceName(), p.UniqueIndexId(), 0, 1, tarantool.IterEq, A.X{p.Id})
	if L.IsError(err, `Promos.FindById failed: `+p.SpaceName()) {
		return false
	}
	rows := res.Tuples()
	if len(rows) == 1 {
		p.FromArray(rows[0])
		return true
	}
	return false
}

var GraphqlFieldPromosById = &graphql.Field{
	Type:        GraphqlTypePromos,
	Description: `list of Promos`,
	Args: graphql.FieldConfigArgument{
		`Id`: &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
	},
}

func (g *Promos) GraphqlFieldPromosByIdWithResolver() *graphql.Field {
	field := *GraphqlFieldPromosById
	field.Resolve = func(p graphql.ResolveParams) (interface{}, error) {
		q := g
		v, ok := p.Args[`id`]
		if !ok {
			v, _ = p.Args[`Id`]
		}
		q.Id = X.ToU(v)
		if q.FindById() {
			return q, nil
		}
		return nil, nil
	}
	return &field
}

func (p *Promos) sqlSelectAllFields() string { //nolint:dupl false positive
	return ` "id"
	, "createdAt"
	, "createdBy"
	, "updatedAt"
	, "updatedBy"
	, "deletedAt"
	, "deletedBy"
	, "isDeleted"
	, "restoredAt"
	, "restoredBy"
	, "startAt"
	, "endAt"
	, "productId"
	, "productCount"
	, "freeProductId"
	, "discountCount"
	, "discountPercent"
	`
}

func (p *Promos) ToUpdateArray() A.X { //nolint:dupl false positive
	return A.X{
		A.X{`=`, 0, p.Id},
		A.X{`=`, 1, p.CreatedAt},
		A.X{`=`, 2, p.CreatedBy},
		A.X{`=`, 3, p.UpdatedAt},
		A.X{`=`, 4, p.UpdatedBy},
		A.X{`=`, 5, p.DeletedAt},
		A.X{`=`, 6, p.DeletedBy},
		A.X{`=`, 7, p.IsDeleted},
		A.X{`=`, 8, p.RestoredAt},
		A.X{`=`, 9, p.RestoredBy},
		A.X{`=`, 10, p.StartAt},
		A.X{`=`, 11, p.EndAt},
		A.X{`=`, 12, p.ProductId},
		A.X{`=`, 13, p.ProductCount},
		A.X{`=`, 14, p.FreeProductId},
		A.X{`=`, 15, p.DiscountCount},
		A.X{`=`, 16, p.DiscountPercent},
	}
}

func (p *Promos) IdxId() int { //nolint:dupl false positive
	return 0
}

func (p *Promos) sqlId() string { //nolint:dupl false positive
	return `"id"`
}

func (p *Promos) IdxCreatedAt() int { //nolint:dupl false positive
	return 1
}

func (p *Promos) sqlCreatedAt() string { //nolint:dupl false positive
	return `"createdAt"`
}

func (p *Promos) IdxCreatedBy() int { //nolint:dupl false positive
	return 2
}

func (p *Promos) sqlCreatedBy() string { //nolint:dupl false positive
	return `"createdBy"`
}

func (p *Promos) IdxUpdatedAt() int { //nolint:dupl false positive
	return 3
}

func (p *Promos) sqlUpdatedAt() string { //nolint:dupl false positive
	return `"updatedAt"`
}

func (p *Promos) IdxUpdatedBy() int { //nolint:dupl false positive
	return 4
}

func (p *Promos) sqlUpdatedBy() string { //nolint:dupl false positive
	return `"updatedBy"`
}

func (p *Promos) IdxDeletedAt() int { //nolint:dupl false positive
	return 5
}

func (p *Promos) sqlDeletedAt() string { //nolint:dupl false positive
	return `"deletedAt"`
}

func (p *Promos) IdxDeletedBy() int { //nolint:dupl false positive
	return 6
}

func (p *Promos) sqlDeletedBy() string { //nolint:dupl false positive
	return `"deletedBy"`
}

func (p *Promos) IdxIsDeleted() int { //nolint:dupl false positive
	return 7
}

func (p *Promos) sqlIsDeleted() string { //nolint:dupl false positive
	return `"isDeleted"`
}

func (p *Promos) IdxRestoredAt() int { //nolint:dupl false positive
	return 8
}

func (p *Promos) sqlRestoredAt() string { //nolint:dupl false positive
	return `"restoredAt"`
}

func (p *Promos) IdxRestoredBy() int { //nolint:dupl false positive
	return 9
}

func (p *Promos) sqlRestoredBy() string { //nolint:dupl false positive
	return `"restoredBy"`
}

func (p *Promos) IdxStartAt() int { //nolint:dupl false positive
	return 10
}

func (p *Promos) sqlStartAt() string { //nolint:dupl false positive
	return `"startAt"`
}

func (p *Promos) IdxEndAt() int { //nolint:dupl false positive
	return 11
}

func (p *Promos) sqlEndAt() string { //nolint:dupl false positive
	return `"endAt"`
}

func (p *Promos) IdxProductId() int { //nolint:dupl false positive
	return 12
}

func (p *Promos) sqlProductId() string { //nolint:dupl false positive
	return `"productId"`
}

func (p *Promos) IdxProductCount() int { //nolint:dupl false positive
	return 13
}

func (p *Promos) sqlProductCount() string { //nolint:dupl false positive
	return `"productCount"`
}

func (p *Promos) IdxFreeProductId() int { //nolint:dupl false positive
	return 14
}

func (p *Promos) sqlFreeProductId() string { //nolint:dupl false positive
	return `"freeProductId"`
}

func (p *Promos) IdxDiscountCount() int { //nolint:dupl false positive
	return 15
}

func (p *Promos) sqlDiscountCount() string { //nolint:dupl false positive
	return `"discountCount"`
}

func (p *Promos) IdxDiscountPercent() int { //nolint:dupl false positive
	return 16
}

func (p *Promos) sqlDiscountPercent() string { //nolint:dupl false positive
	return `"discountPercent"`
}

func (p *Promos) ToArray() A.X { //nolint:dupl false positive
	return A.X{
		p.Id,              // 0
		p.CreatedAt,       // 1
		p.CreatedBy,       // 2
		p.UpdatedAt,       // 3
		p.UpdatedBy,       // 4
		p.DeletedAt,       // 5
		p.DeletedBy,       // 6
		p.IsDeleted,       // 7
		p.RestoredAt,      // 8
		p.RestoredBy,      // 9
		p.StartAt,         // 10
		p.EndAt,           // 11
		p.ProductId,       // 12
		p.ProductCount,    // 13
		p.FreeProductId,   // 14
		p.DiscountCount,   // 15
		p.DiscountPercent, // 16
	}
}

func (p *Promos) FromArray(a A.X) *Promos { //nolint:dupl false positive
	p.Id = X.ToU(a[0])
	p.CreatedAt = X.ToI(a[1])
	p.CreatedBy = X.ToU(a[2])
	p.UpdatedAt = X.ToI(a[3])
	p.UpdatedBy = X.ToU(a[4])
	p.DeletedAt = X.ToI(a[5])
	p.DeletedBy = X.ToU(a[6])
	p.IsDeleted = X.ToBool(a[7])
	p.RestoredAt = X.ToI(a[8])
	p.RestoredBy = X.ToU(a[9])
	p.StartAt = X.ToI(a[10])
	p.EndAt = X.ToI(a[11])
	p.ProductId = X.ToU(a[12])
	p.ProductCount = X.ToU(a[13])
	p.FreeProductId = X.ToU(a[14])
	p.DiscountCount = X.ToU(a[15])
	p.DiscountPercent = X.ToF(a[16])
	return p
}

func (p *Promos) Total() int64 { //nolint:dupl false positive
	rows := p.Adapter.CallBoxSpace(p.SpaceName()+`:count`, A.X{})
	if len(rows) > 0 && len(rows[0]) > 0 {
		return X.ToI(rows[0][0])
	}
	return 0
}

var GraphqlTypePromos = graphql.NewObject(
	graphql.ObjectConfig{
		Name: `promos`,
		Fields: graphql.Fields{
			`id`: &graphql.Field{
				Type: graphql.Int,
			},
			`createdAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`createdBy`: &graphql.Field{
				Type: graphql.ID,
			},
			`updatedAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`updatedBy`: &graphql.Field{
				Type: graphql.ID,
			},
			`deletedAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`deletedBy`: &graphql.Field{
				Type: graphql.ID,
			},
			`isDeleted`: &graphql.Field{
				Type: graphql.Boolean,
			},
			`restoredAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`restoredBy`: &graphql.Field{
				Type: graphql.ID,
			},
			`startAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`endAt`: &graphql.Field{
				Type: graphql.Int,
			},
			`productId`: &graphql.Field{
				Type: graphql.ID,
			},
			`productCount`: &graphql.Field{
				Type: graphql.Int,
			},
			`freeProductId`: &graphql.Field{
				Type: graphql.ID,
			},
			`discountCount`: &graphql.Field{
				Type: graphql.Int,
			},
			`discountPercent`: &graphql.Field{
				Type: graphql.Float,
			},
		},
	},
)

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go
