package mStore

import (
	"github.com/kokizzu/gotro/D/Tt"
)

// table products, promotions, cart, invoice
const (
	TableProducts Tt.TableName = `products`
	Id                         = `id`
	CreatedBy                  = `createdBy`
	CreatedAt                  = `createdAt`
	UpdatedBy                  = `updatedBy`
	UpdatedAt                  = `updatedAt`
	DeletedBy                  = `deletedBy`
	DeletedAt                  = `deletedAt`
	IsDeleted                  = `isDeleted`
	RestoredBy                 = `restoredBy`
	RestoredAt                 = `restoredAt`
	Sku                        = `sku`
	Name                       = `name`
	Price                      = `price`
	InventoryQty               = `inventoryQty`
	WeightGram                 = `weightGram`
	Info                       = `info` // promo or other information (eg. product no longer exists)
)

const (
	TablePromos     Tt.TableName = `promos`
	StartAt                      = `startAt`
	EndAt                        = `endAt`
	ProductId                    = `productId` // product being purchased
	ProductCount                 = `productCount`
	FreeProductId                = `freeProductId`   // buy x get y
	DiscountCount                = `discountCount`   // buy 3 get 2
	DiscountPercent              = `discountPercent` // buy x discount y % for all x
)

const (
	TableInvoices Tt.TableName = `invoices`
	// StoreId                     = `storeId`
	OwnerId        = `ownerId`
	TotalWeight    = `totalWeight`
	TotalPrice     = `totalPrice`
	TotalDiscount  = `totalDiscount`
	DeliveryMethod = `deliveryMethod`
	DeliveryPrice  = `deliveryPrice`
	TotalPaid      = `totalPaid` // TotalPrice - TotalDiscount + DeliveryPrice
	PaymentMethod  = `paymentMethod`
	DeadlineAt     = `deadlineAt` // payment deadline
	PaidAt         = `paidAt`
	PromoRuleIds   = `promoRuleIds` // applied rules in separated by space
)

const (
	TableCartItems Tt.TableName = `cartItems`
	InvoiceId                   = `invoiceId`
	NameCopy                    = `nameCopy`
	PriceCopy                   = `priceCopy`
	Qty                         = `qty`
	Discount                    = `discount`
	SubTotal                    = `subTotal` // = PriceCopy x Qty - Discount
)

var TarantoolTables = map[Tt.TableName]*Tt.TableProp{
	// can only adding fields on back, and must IsNullable: true
	// primary key must be first field and set to Unique: fieldName
	TableProducts: {
		Fields: []Tt.Field{
			{Id, Tt.Unsigned},
			{CreatedAt, Tt.Integer},
			{CreatedBy, Tt.Unsigned},
			{UpdatedAt, Tt.Integer},
			{UpdatedBy, Tt.Unsigned},
			{DeletedAt, Tt.Integer},
			{DeletedBy, Tt.Unsigned},
			{IsDeleted, Tt.Boolean},
			{RestoredAt, Tt.Integer},
			{RestoredBy, Tt.Unsigned},
			{Sku, Tt.String},
			{Name, Tt.String},
			{Price, Tt.Unsigned},
			{InventoryQty, Tt.Unsigned},
			{WeightGram, Tt.Unsigned},
		},
		AutoIncrementId: true,
		Unique1:         Sku,
	},
	TablePromos: {
		Fields: []Tt.Field{
			{Id, Tt.Unsigned},
			{CreatedAt, Tt.Integer},
			{CreatedBy, Tt.Unsigned},
			{UpdatedAt, Tt.Integer},
			{UpdatedBy, Tt.Unsigned},
			{DeletedAt, Tt.Integer},
			{DeletedBy, Tt.Unsigned},
			{IsDeleted, Tt.Boolean},
			{RestoredAt, Tt.Integer},
			{RestoredBy, Tt.Unsigned},
			{StartAt, Tt.Integer},
			{EndAt, Tt.Integer},
			{ProductId, Tt.Unsigned},
			{ProductCount, Tt.Unsigned},
			{FreeProductId, Tt.Unsigned},
			{DiscountCount, Tt.Unsigned},
			{DiscountPercent, Tt.Number},
		},
		AutoIncrementId: true,
	},
	TableCartItems: {
		Fields: []Tt.Field{
			{Id, Tt.Unsigned},
			{CreatedAt, Tt.Integer},
			{CreatedBy, Tt.Unsigned},
			{UpdatedAt, Tt.Integer},
			{UpdatedBy, Tt.Unsigned},
			{DeletedAt, Tt.Integer},
			{DeletedBy, Tt.Unsigned},
			{IsDeleted, Tt.Boolean},
			{RestoredAt, Tt.Integer},
			{RestoredBy, Tt.Unsigned},
			{OwnerId, Tt.Unsigned},
			{InvoiceId, Tt.Unsigned},
			{ProductId, Tt.Unsigned},
			{NameCopy, Tt.String},
			{PriceCopy, Tt.Integer},
			{Qty, Tt.Integer}, // negative = refund
			{Discount, Tt.Unsigned},
			{SubTotal, Tt.Integer}, // negative = refund
			{Info, Tt.String},
		},
		AutoIncrementId: true,
		Uniques:         []string{OwnerId, InvoiceId, ProductId},
	},
	TableInvoices: {
		Fields: []Tt.Field{
			{Id, Tt.Unsigned},
			{CreatedAt, Tt.Integer},
			{CreatedBy, Tt.Unsigned},
			{UpdatedAt, Tt.Integer},
			{UpdatedBy, Tt.Unsigned},
			{DeletedAt, Tt.Integer},
			{DeletedBy, Tt.Unsigned},
			{IsDeleted, Tt.Boolean},
			{RestoredAt, Tt.Integer},
			{RestoredBy, Tt.Unsigned},
			{OwnerId, Tt.Unsigned},
			{TotalWeight, Tt.Unsigned},
			{TotalPrice, Tt.Unsigned},
			{TotalDiscount, Tt.Unsigned},
			{DeliveryMethod, Tt.Unsigned},
			{DeliveryPrice, Tt.Unsigned},
			{TotalPaid, Tt.Unsigned},
			{PaidAt, Tt.Unsigned},
			{PaymentMethod, Tt.Unsigned},
			{DeadlineAt, Tt.Unsigned},
			{PromoRuleIds, Tt.String},
		},
		AutoIncrementId: true,
		Indexes:         []string{OwnerId, PaidAt, DeadlineAt},
	},
}

func GenerateORM() {
	Tt.GenerateOrm(TarantoolTables)
	//Ch.GenerateOrm(ClickhouseTables) // find d.InitClickhouseBuffer to create chBuffer on NewDomain
}
