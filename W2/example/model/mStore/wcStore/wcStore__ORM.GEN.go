package wcStore

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

import (
	"github.com/kokizzu/gotro/W2/example/model/mStore/rqStore"

	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/D/Tt"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/X"
)

//go:generate gomodifytags -all -add-tags json,form,query,long,msg -transform camelcase --skip-unexported -w -file wcStore__ORM.GEN.go
//go:generate replacer 'Id" form' 'Id,string" form' type wcStore__ORM.GEN.go
//go:generate replacer 'json:"id"' 'json:"id,string"' type wcStore__ORM.GEN.go
//go:generate replacer 'By" form' 'By,string" form' type wcStore__ORM.GEN.go
// go:generate msgp -tests=false -file wcStore__ORM.GEN.go -o wcStore__MSG.GEN.go

type CartItemsMutator struct {
	rqStore.CartItems
	mutations []A.X
}

func NewCartItemsMutator(adapter *Tt.Adapter) *CartItemsMutator {
	return &CartItemsMutator{CartItems: rqStore.CartItems{Adapter: adapter}}
}

func (c *CartItemsMutator) HaveMutation() bool { //nolint:dupl false positive
	return len(c.mutations) > 0
}

// Overwrite all columns, error if not exists
func (c *CartItemsMutator) DoOverwriteById() bool { //nolint:dupl false positive
	_, err := c.Adapter.Update(c.SpaceName(), c.UniqueIndexId(), A.X{c.Id}, c.ToUpdateArray())
	return !L.IsError(err, `CartItems.DoOverwriteById failed: `+c.SpaceName())
}

// Update only mutated, error if not exists, use Find* and Set* methods instead of direct assignment
func (c *CartItemsMutator) DoUpdateById() bool { //nolint:dupl false positive
	if !c.HaveMutation() {
		return true
	}
	_, err := c.Adapter.Update(c.SpaceName(), c.UniqueIndexId(), A.X{c.Id}, c.mutations)
	return !L.IsError(err, `CartItems.DoUpdateById failed: `+c.SpaceName())
}

func (c *CartItemsMutator) DoDeletePermanentById() bool { //nolint:dupl false positive
	_, err := c.Adapter.Delete(c.SpaceName(), c.UniqueIndexId(), A.X{c.Id})
	return !L.IsError(err, `CartItems.DoDeletePermanentById failed: `+c.SpaceName())
}

// func (c *CartItemsMutator) DoUpsert() bool { //nolint:dupl false positive
//	_, err := c.Adapter.Upsert(c.SpaceName(), c.ToArray(), A.X{
//		A.X{`=`, 0, c.Id},
//		A.X{`=`, 1, c.CreatedAt},
//		A.X{`=`, 2, c.CreatedBy},
//		A.X{`=`, 3, c.UpdatedAt},
//		A.X{`=`, 4, c.UpdatedBy},
//		A.X{`=`, 5, c.DeletedAt},
//		A.X{`=`, 6, c.DeletedBy},
//		A.X{`=`, 7, c.IsDeleted},
//		A.X{`=`, 8, c.RestoredAt},
//		A.X{`=`, 9, c.RestoredBy},
//		A.X{`=`, 10, c.OwnerId},
//		A.X{`=`, 11, c.InvoiceId},
//		A.X{`=`, 12, c.ProductId},
//		A.X{`=`, 13, c.NameCopy},
//		A.X{`=`, 14, c.PriceCopy},
//		A.X{`=`, 15, c.Qty},
//		A.X{`=`, 16, c.Discount},
//		A.X{`=`, 17, c.SubTotal},
//		A.X{`=`, 18, c.Info},
//	})
//	return !L.IsError(err, `CartItems.DoUpsert failed: `+c.SpaceName())
// }

// Overwrite all columns, error if not exists
func (c *CartItemsMutator) DoOverwriteByOwnerIdInvoiceIdProductId() bool { //nolint:dupl false positive
	_, err := c.Adapter.Update(c.SpaceName(), c.UniqueIndexOwnerIdInvoiceIdProductId(), A.X{c.OwnerId, c.InvoiceId, c.ProductId}, c.ToUpdateArray())
	return !L.IsError(err, `CartItems.DoOverwriteByOwnerIdInvoiceIdProductId failed: `+c.SpaceName())
}

// Update only mutated, error if not exists, use Find* and Set* methods instead of direct assignment
func (c *CartItemsMutator) DoUpdateByOwnerIdInvoiceIdProductId() bool { //nolint:dupl false positive
	if !c.HaveMutation() {
		return true
	}
	_, err := c.Adapter.Update(c.SpaceName(), c.UniqueIndexOwnerIdInvoiceIdProductId(), A.X{c.OwnerId, c.InvoiceId, c.ProductId}, c.mutations)
	return !L.IsError(err, `CartItems.DoUpdateByOwnerIdInvoiceIdProductId failed: `+c.SpaceName())
}

func (c *CartItemsMutator) DoDeletePermanentByOwnerIdInvoiceIdProductId() bool { //nolint:dupl false positive
	_, err := c.Adapter.Delete(c.SpaceName(), c.UniqueIndexOwnerIdInvoiceIdProductId(), A.X{c.OwnerId, c.InvoiceId, c.ProductId})
	return !L.IsError(err, `CartItems.DoDeletePermanentByOwnerIdInvoiceIdProductId failed: `+c.SpaceName())
}

// insert, error if exists
func (c *CartItemsMutator) DoInsert() bool { //nolint:dupl false positive
	row, err := c.Adapter.Insert(c.SpaceName(), c.ToArray())
	if err == nil {
		tup := row.Tuples()
		if len(tup) > 0 && len(tup[0]) > 0 && tup[0][0] != nil {
			c.Id = X.ToU(tup[0][0])
		}
	}
	return !L.IsError(err, `CartItems.DoInsert failed: `+c.SpaceName())
}

// replace = upsert, only error when there's unique secondary key
func (c *CartItemsMutator) DoReplace() bool { //nolint:dupl false positive
	_, err := c.Adapter.Replace(c.SpaceName(), c.ToArray())
	return !L.IsError(err, `CartItems.DoReplace failed: `+c.SpaceName())
}

func (c *CartItemsMutator) SetId(val uint64) bool { //nolint:dupl false positive
	if val != c.Id {
		c.mutations = append(c.mutations, A.X{`=`, 0, val})
		c.Id = val
		return true
	}
	return false
}

func (c *CartItemsMutator) SetCreatedAt(val int64) bool { //nolint:dupl false positive
	if val != c.CreatedAt {
		c.mutations = append(c.mutations, A.X{`=`, 1, val})
		c.CreatedAt = val
		return true
	}
	return false
}

func (c *CartItemsMutator) SetCreatedBy(val uint64) bool { //nolint:dupl false positive
	if val != c.CreatedBy {
		c.mutations = append(c.mutations, A.X{`=`, 2, val})
		c.CreatedBy = val
		return true
	}
	return false
}

func (c *CartItemsMutator) SetUpdatedAt(val int64) bool { //nolint:dupl false positive
	if val != c.UpdatedAt {
		c.mutations = append(c.mutations, A.X{`=`, 3, val})
		c.UpdatedAt = val
		return true
	}
	return false
}

func (c *CartItemsMutator) SetUpdatedBy(val uint64) bool { //nolint:dupl false positive
	if val != c.UpdatedBy {
		c.mutations = append(c.mutations, A.X{`=`, 4, val})
		c.UpdatedBy = val
		return true
	}
	return false
}

func (c *CartItemsMutator) SetDeletedAt(val int64) bool { //nolint:dupl false positive
	if val != c.DeletedAt {
		c.mutations = append(c.mutations, A.X{`=`, 5, val})
		c.DeletedAt = val
		return true
	}
	return false
}

func (c *CartItemsMutator) SetDeletedBy(val uint64) bool { //nolint:dupl false positive
	if val != c.DeletedBy {
		c.mutations = append(c.mutations, A.X{`=`, 6, val})
		c.DeletedBy = val
		return true
	}
	return false
}

func (c *CartItemsMutator) SetIsDeleted(val bool) bool { //nolint:dupl false positive
	if val != c.IsDeleted {
		c.mutations = append(c.mutations, A.X{`=`, 7, val})
		c.IsDeleted = val
		return true
	}
	return false
}

func (c *CartItemsMutator) SetRestoredAt(val int64) bool { //nolint:dupl false positive
	if val != c.RestoredAt {
		c.mutations = append(c.mutations, A.X{`=`, 8, val})
		c.RestoredAt = val
		return true
	}
	return false
}

func (c *CartItemsMutator) SetRestoredBy(val uint64) bool { //nolint:dupl false positive
	if val != c.RestoredBy {
		c.mutations = append(c.mutations, A.X{`=`, 9, val})
		c.RestoredBy = val
		return true
	}
	return false
}

func (c *CartItemsMutator) SetOwnerId(val uint64) bool { //nolint:dupl false positive
	if val != c.OwnerId {
		c.mutations = append(c.mutations, A.X{`=`, 10, val})
		c.OwnerId = val
		return true
	}
	return false
}

func (c *CartItemsMutator) SetInvoiceId(val uint64) bool { //nolint:dupl false positive
	if val != c.InvoiceId {
		c.mutations = append(c.mutations, A.X{`=`, 11, val})
		c.InvoiceId = val
		return true
	}
	return false
}

func (c *CartItemsMutator) SetProductId(val uint64) bool { //nolint:dupl false positive
	if val != c.ProductId {
		c.mutations = append(c.mutations, A.X{`=`, 12, val})
		c.ProductId = val
		return true
	}
	return false
}

func (c *CartItemsMutator) SetNameCopy(val string) bool { //nolint:dupl false positive
	if val != c.NameCopy {
		c.mutations = append(c.mutations, A.X{`=`, 13, val})
		c.NameCopy = val
		return true
	}
	return false
}

func (c *CartItemsMutator) SetPriceCopy(val int64) bool { //nolint:dupl false positive
	if val != c.PriceCopy {
		c.mutations = append(c.mutations, A.X{`=`, 14, val})
		c.PriceCopy = val
		return true
	}
	return false
}

func (c *CartItemsMutator) SetQty(val int64) bool { //nolint:dupl false positive
	if val != c.Qty {
		c.mutations = append(c.mutations, A.X{`=`, 15, val})
		c.Qty = val
		return true
	}
	return false
}

func (c *CartItemsMutator) SetDiscount(val uint64) bool { //nolint:dupl false positive
	if val != c.Discount {
		c.mutations = append(c.mutations, A.X{`=`, 16, val})
		c.Discount = val
		return true
	}
	return false
}

func (c *CartItemsMutator) SetSubTotal(val int64) bool { //nolint:dupl false positive
	if val != c.SubTotal {
		c.mutations = append(c.mutations, A.X{`=`, 17, val})
		c.SubTotal = val
		return true
	}
	return false
}

func (c *CartItemsMutator) SetInfo(val string) bool { //nolint:dupl false positive
	if val != c.Info {
		c.mutations = append(c.mutations, A.X{`=`, 18, val})
		c.Info = val
		return true
	}
	return false
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

type InvoicesMutator struct {
	rqStore.Invoices
	mutations []A.X
}

func NewInvoicesMutator(adapter *Tt.Adapter) *InvoicesMutator {
	return &InvoicesMutator{Invoices: rqStore.Invoices{Adapter: adapter}}
}

func (i *InvoicesMutator) HaveMutation() bool { //nolint:dupl false positive
	return len(i.mutations) > 0
}

// Overwrite all columns, error if not exists
func (i *InvoicesMutator) DoOverwriteById() bool { //nolint:dupl false positive
	_, err := i.Adapter.Update(i.SpaceName(), i.UniqueIndexId(), A.X{i.Id}, i.ToUpdateArray())
	return !L.IsError(err, `Invoices.DoOverwriteById failed: `+i.SpaceName())
}

// Update only mutated, error if not exists, use Find* and Set* methods instead of direct assignment
func (i *InvoicesMutator) DoUpdateById() bool { //nolint:dupl false positive
	if !i.HaveMutation() {
		return true
	}
	_, err := i.Adapter.Update(i.SpaceName(), i.UniqueIndexId(), A.X{i.Id}, i.mutations)
	return !L.IsError(err, `Invoices.DoUpdateById failed: `+i.SpaceName())
}

func (i *InvoicesMutator) DoDeletePermanentById() bool { //nolint:dupl false positive
	_, err := i.Adapter.Delete(i.SpaceName(), i.UniqueIndexId(), A.X{i.Id})
	return !L.IsError(err, `Invoices.DoDeletePermanentById failed: `+i.SpaceName())
}

// func (i *InvoicesMutator) DoUpsert() bool { //nolint:dupl false positive
//	_, err := i.Adapter.Upsert(i.SpaceName(), i.ToArray(), A.X{
//		A.X{`=`, 0, i.Id},
//		A.X{`=`, 1, i.CreatedAt},
//		A.X{`=`, 2, i.CreatedBy},
//		A.X{`=`, 3, i.UpdatedAt},
//		A.X{`=`, 4, i.UpdatedBy},
//		A.X{`=`, 5, i.DeletedAt},
//		A.X{`=`, 6, i.DeletedBy},
//		A.X{`=`, 7, i.IsDeleted},
//		A.X{`=`, 8, i.RestoredAt},
//		A.X{`=`, 9, i.RestoredBy},
//		A.X{`=`, 10, i.OwnerId},
//		A.X{`=`, 11, i.TotalWeight},
//		A.X{`=`, 12, i.TotalPrice},
//		A.X{`=`, 13, i.TotalDiscount},
//		A.X{`=`, 14, i.DeliveryMethod},
//		A.X{`=`, 15, i.DeliveryPrice},
//		A.X{`=`, 16, i.TotalPaid},
//		A.X{`=`, 17, i.PaidAt},
//		A.X{`=`, 18, i.PaymentMethod},
//		A.X{`=`, 19, i.DeadlineAt},
//		A.X{`=`, 20, i.PromoRuleIds},
//	})
//	return !L.IsError(err, `Invoices.DoUpsert failed: `+i.SpaceName())
// }

// insert, error if exists
func (i *InvoicesMutator) DoInsert() bool { //nolint:dupl false positive
	row, err := i.Adapter.Insert(i.SpaceName(), i.ToArray())
	if err == nil {
		tup := row.Tuples()
		if len(tup) > 0 && len(tup[0]) > 0 && tup[0][0] != nil {
			i.Id = X.ToU(tup[0][0])
		}
	}
	return !L.IsError(err, `Invoices.DoInsert failed: `+i.SpaceName())
}

// replace = upsert, only error when there's unique secondary key
func (i *InvoicesMutator) DoReplace() bool { //nolint:dupl false positive
	_, err := i.Adapter.Replace(i.SpaceName(), i.ToArray())
	return !L.IsError(err, `Invoices.DoReplace failed: `+i.SpaceName())
}

func (i *InvoicesMutator) SetId(val uint64) bool { //nolint:dupl false positive
	if val != i.Id {
		i.mutations = append(i.mutations, A.X{`=`, 0, val})
		i.Id = val
		return true
	}
	return false
}

func (i *InvoicesMutator) SetCreatedAt(val int64) bool { //nolint:dupl false positive
	if val != i.CreatedAt {
		i.mutations = append(i.mutations, A.X{`=`, 1, val})
		i.CreatedAt = val
		return true
	}
	return false
}

func (i *InvoicesMutator) SetCreatedBy(val uint64) bool { //nolint:dupl false positive
	if val != i.CreatedBy {
		i.mutations = append(i.mutations, A.X{`=`, 2, val})
		i.CreatedBy = val
		return true
	}
	return false
}

func (i *InvoicesMutator) SetUpdatedAt(val int64) bool { //nolint:dupl false positive
	if val != i.UpdatedAt {
		i.mutations = append(i.mutations, A.X{`=`, 3, val})
		i.UpdatedAt = val
		return true
	}
	return false
}

func (i *InvoicesMutator) SetUpdatedBy(val uint64) bool { //nolint:dupl false positive
	if val != i.UpdatedBy {
		i.mutations = append(i.mutations, A.X{`=`, 4, val})
		i.UpdatedBy = val
		return true
	}
	return false
}

func (i *InvoicesMutator) SetDeletedAt(val int64) bool { //nolint:dupl false positive
	if val != i.DeletedAt {
		i.mutations = append(i.mutations, A.X{`=`, 5, val})
		i.DeletedAt = val
		return true
	}
	return false
}

func (i *InvoicesMutator) SetDeletedBy(val uint64) bool { //nolint:dupl false positive
	if val != i.DeletedBy {
		i.mutations = append(i.mutations, A.X{`=`, 6, val})
		i.DeletedBy = val
		return true
	}
	return false
}

func (i *InvoicesMutator) SetIsDeleted(val bool) bool { //nolint:dupl false positive
	if val != i.IsDeleted {
		i.mutations = append(i.mutations, A.X{`=`, 7, val})
		i.IsDeleted = val
		return true
	}
	return false
}

func (i *InvoicesMutator) SetRestoredAt(val int64) bool { //nolint:dupl false positive
	if val != i.RestoredAt {
		i.mutations = append(i.mutations, A.X{`=`, 8, val})
		i.RestoredAt = val
		return true
	}
	return false
}

func (i *InvoicesMutator) SetRestoredBy(val uint64) bool { //nolint:dupl false positive
	if val != i.RestoredBy {
		i.mutations = append(i.mutations, A.X{`=`, 9, val})
		i.RestoredBy = val
		return true
	}
	return false
}

func (i *InvoicesMutator) SetOwnerId(val uint64) bool { //nolint:dupl false positive
	if val != i.OwnerId {
		i.mutations = append(i.mutations, A.X{`=`, 10, val})
		i.OwnerId = val
		return true
	}
	return false
}

func (i *InvoicesMutator) SetTotalWeight(val uint64) bool { //nolint:dupl false positive
	if val != i.TotalWeight {
		i.mutations = append(i.mutations, A.X{`=`, 11, val})
		i.TotalWeight = val
		return true
	}
	return false
}

func (i *InvoicesMutator) SetTotalPrice(val uint64) bool { //nolint:dupl false positive
	if val != i.TotalPrice {
		i.mutations = append(i.mutations, A.X{`=`, 12, val})
		i.TotalPrice = val
		return true
	}
	return false
}

func (i *InvoicesMutator) SetTotalDiscount(val uint64) bool { //nolint:dupl false positive
	if val != i.TotalDiscount {
		i.mutations = append(i.mutations, A.X{`=`, 13, val})
		i.TotalDiscount = val
		return true
	}
	return false
}

func (i *InvoicesMutator) SetDeliveryMethod(val uint64) bool { //nolint:dupl false positive
	if val != i.DeliveryMethod {
		i.mutations = append(i.mutations, A.X{`=`, 14, val})
		i.DeliveryMethod = val
		return true
	}
	return false
}

func (i *InvoicesMutator) SetDeliveryPrice(val uint64) bool { //nolint:dupl false positive
	if val != i.DeliveryPrice {
		i.mutations = append(i.mutations, A.X{`=`, 15, val})
		i.DeliveryPrice = val
		return true
	}
	return false
}

func (i *InvoicesMutator) SetTotalPaid(val uint64) bool { //nolint:dupl false positive
	if val != i.TotalPaid {
		i.mutations = append(i.mutations, A.X{`=`, 16, val})
		i.TotalPaid = val
		return true
	}
	return false
}

func (i *InvoicesMutator) SetPaidAt(val uint64) bool { //nolint:dupl false positive
	if val != i.PaidAt {
		i.mutations = append(i.mutations, A.X{`=`, 17, val})
		i.PaidAt = val
		return true
	}
	return false
}

func (i *InvoicesMutator) SetPaymentMethod(val uint64) bool { //nolint:dupl false positive
	if val != i.PaymentMethod {
		i.mutations = append(i.mutations, A.X{`=`, 18, val})
		i.PaymentMethod = val
		return true
	}
	return false
}

func (i *InvoicesMutator) SetDeadlineAt(val uint64) bool { //nolint:dupl false positive
	if val != i.DeadlineAt {
		i.mutations = append(i.mutations, A.X{`=`, 19, val})
		i.DeadlineAt = val
		return true
	}
	return false
}

func (i *InvoicesMutator) SetPromoRuleIds(val string) bool { //nolint:dupl false positive
	if val != i.PromoRuleIds {
		i.mutations = append(i.mutations, A.X{`=`, 20, val})
		i.PromoRuleIds = val
		return true
	}
	return false
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

type ProductsMutator struct {
	rqStore.Products
	mutations []A.X
}

func NewProductsMutator(adapter *Tt.Adapter) *ProductsMutator {
	return &ProductsMutator{Products: rqStore.Products{Adapter: adapter}}
}

func (p *ProductsMutator) HaveMutation() bool { //nolint:dupl false positive
	return len(p.mutations) > 0
}

// Overwrite all columns, error if not exists
func (p *ProductsMutator) DoOverwriteById() bool { //nolint:dupl false positive
	_, err := p.Adapter.Update(p.SpaceName(), p.UniqueIndexId(), A.X{p.Id}, p.ToUpdateArray())
	return !L.IsError(err, `Products.DoOverwriteById failed: `+p.SpaceName())
}

// Update only mutated, error if not exists, use Find* and Set* methods instead of direct assignment
func (p *ProductsMutator) DoUpdateById() bool { //nolint:dupl false positive
	if !p.HaveMutation() {
		return true
	}
	_, err := p.Adapter.Update(p.SpaceName(), p.UniqueIndexId(), A.X{p.Id}, p.mutations)
	return !L.IsError(err, `Products.DoUpdateById failed: `+p.SpaceName())
}

func (p *ProductsMutator) DoDeletePermanentById() bool { //nolint:dupl false positive
	_, err := p.Adapter.Delete(p.SpaceName(), p.UniqueIndexId(), A.X{p.Id})
	return !L.IsError(err, `Products.DoDeletePermanentById failed: `+p.SpaceName())
}

// func (p *ProductsMutator) DoUpsert() bool { //nolint:dupl false positive
//	_, err := p.Adapter.Upsert(p.SpaceName(), p.ToArray(), A.X{
//		A.X{`=`, 0, p.Id},
//		A.X{`=`, 1, p.CreatedAt},
//		A.X{`=`, 2, p.CreatedBy},
//		A.X{`=`, 3, p.UpdatedAt},
//		A.X{`=`, 4, p.UpdatedBy},
//		A.X{`=`, 5, p.DeletedAt},
//		A.X{`=`, 6, p.DeletedBy},
//		A.X{`=`, 7, p.IsDeleted},
//		A.X{`=`, 8, p.RestoredAt},
//		A.X{`=`, 9, p.RestoredBy},
//		A.X{`=`, 10, p.Sku},
//		A.X{`=`, 11, p.Name},
//		A.X{`=`, 12, p.Price},
//		A.X{`=`, 13, p.InventoryQty},
//		A.X{`=`, 14, p.WeightGram},
//	})
//	return !L.IsError(err, `Products.DoUpsert failed: `+p.SpaceName())
// }

// Overwrite all columns, error if not exists
func (p *ProductsMutator) DoOverwriteBySku() bool { //nolint:dupl false positive
	_, err := p.Adapter.Update(p.SpaceName(), p.UniqueIndexSku(), A.X{p.Sku}, p.ToUpdateArray())
	return !L.IsError(err, `Products.DoOverwriteBySku failed: `+p.SpaceName())
}

// Update only mutated, error if not exists, use Find* and Set* methods instead of direct assignment
func (p *ProductsMutator) DoUpdateBySku() bool { //nolint:dupl false positive
	if !p.HaveMutation() {
		return true
	}
	_, err := p.Adapter.Update(p.SpaceName(), p.UniqueIndexSku(), A.X{p.Sku}, p.mutations)
	return !L.IsError(err, `Products.DoUpdateBySku failed: `+p.SpaceName())
}

func (p *ProductsMutator) DoDeletePermanentBySku() bool { //nolint:dupl false positive
	_, err := p.Adapter.Delete(p.SpaceName(), p.UniqueIndexSku(), A.X{p.Sku})
	return !L.IsError(err, `Products.DoDeletePermanentBySku failed: `+p.SpaceName())
}

// Overwrite all columns, error if not exists
func (p *ProductsMutator) DoOverwriteBySku() bool { //nolint:dupl false positive
	_, err := p.Adapter.Update(p.SpaceName(), p.UniqueIndexSku(), A.X{p.Sku}, p.ToUpdateArray())
	return !L.IsError(err, `Products.DoOverwriteBySku failed: `+p.SpaceName())
}

// Update only mutated, error if not exists, use Find* and Set* methods instead of direct assignment
func (p *ProductsMutator) DoUpdateBySku() bool { //nolint:dupl false positive
	if !p.HaveMutation() {
		return true
	}
	_, err := p.Adapter.Update(p.SpaceName(), p.UniqueIndexSku(), A.X{p.Sku}, p.mutations)
	return !L.IsError(err, `Products.DoUpdateBySku failed: `+p.SpaceName())
}

func (p *ProductsMutator) DoDeletePermanentBySku() bool { //nolint:dupl false positive
	_, err := p.Adapter.Delete(p.SpaceName(), p.UniqueIndexSku(), A.X{p.Sku})
	return !L.IsError(err, `Products.DoDeletePermanentBySku failed: `+p.SpaceName())
}

// insert, error if exists
func (p *ProductsMutator) DoInsert() bool { //nolint:dupl false positive
	row, err := p.Adapter.Insert(p.SpaceName(), p.ToArray())
	if err == nil {
		tup := row.Tuples()
		if len(tup) > 0 && len(tup[0]) > 0 && tup[0][0] != nil {
			p.Id = X.ToU(tup[0][0])
		}
	}
	return !L.IsError(err, `Products.DoInsert failed: `+p.SpaceName())
}

// replace = upsert, only error when there's unique secondary key
func (p *ProductsMutator) DoReplace() bool { //nolint:dupl false positive
	_, err := p.Adapter.Replace(p.SpaceName(), p.ToArray())
	return !L.IsError(err, `Products.DoReplace failed: `+p.SpaceName())
}

func (p *ProductsMutator) SetId(val uint64) bool { //nolint:dupl false positive
	if val != p.Id {
		p.mutations = append(p.mutations, A.X{`=`, 0, val})
		p.Id = val
		return true
	}
	return false
}

func (p *ProductsMutator) SetCreatedAt(val int64) bool { //nolint:dupl false positive
	if val != p.CreatedAt {
		p.mutations = append(p.mutations, A.X{`=`, 1, val})
		p.CreatedAt = val
		return true
	}
	return false
}

func (p *ProductsMutator) SetCreatedBy(val uint64) bool { //nolint:dupl false positive
	if val != p.CreatedBy {
		p.mutations = append(p.mutations, A.X{`=`, 2, val})
		p.CreatedBy = val
		return true
	}
	return false
}

func (p *ProductsMutator) SetUpdatedAt(val int64) bool { //nolint:dupl false positive
	if val != p.UpdatedAt {
		p.mutations = append(p.mutations, A.X{`=`, 3, val})
		p.UpdatedAt = val
		return true
	}
	return false
}

func (p *ProductsMutator) SetUpdatedBy(val uint64) bool { //nolint:dupl false positive
	if val != p.UpdatedBy {
		p.mutations = append(p.mutations, A.X{`=`, 4, val})
		p.UpdatedBy = val
		return true
	}
	return false
}

func (p *ProductsMutator) SetDeletedAt(val int64) bool { //nolint:dupl false positive
	if val != p.DeletedAt {
		p.mutations = append(p.mutations, A.X{`=`, 5, val})
		p.DeletedAt = val
		return true
	}
	return false
}

func (p *ProductsMutator) SetDeletedBy(val uint64) bool { //nolint:dupl false positive
	if val != p.DeletedBy {
		p.mutations = append(p.mutations, A.X{`=`, 6, val})
		p.DeletedBy = val
		return true
	}
	return false
}

func (p *ProductsMutator) SetIsDeleted(val bool) bool { //nolint:dupl false positive
	if val != p.IsDeleted {
		p.mutations = append(p.mutations, A.X{`=`, 7, val})
		p.IsDeleted = val
		return true
	}
	return false
}

func (p *ProductsMutator) SetRestoredAt(val int64) bool { //nolint:dupl false positive
	if val != p.RestoredAt {
		p.mutations = append(p.mutations, A.X{`=`, 8, val})
		p.RestoredAt = val
		return true
	}
	return false
}

func (p *ProductsMutator) SetRestoredBy(val uint64) bool { //nolint:dupl false positive
	if val != p.RestoredBy {
		p.mutations = append(p.mutations, A.X{`=`, 9, val})
		p.RestoredBy = val
		return true
	}
	return false
}

func (p *ProductsMutator) SetSku(val string) bool { //nolint:dupl false positive
	if val != p.Sku {
		p.mutations = append(p.mutations, A.X{`=`, 10, val})
		p.Sku = val
		return true
	}
	return false
}

func (p *ProductsMutator) SetName(val string) bool { //nolint:dupl false positive
	if val != p.Name {
		p.mutations = append(p.mutations, A.X{`=`, 11, val})
		p.Name = val
		return true
	}
	return false
}

func (p *ProductsMutator) SetPrice(val uint64) bool { //nolint:dupl false positive
	if val != p.Price {
		p.mutations = append(p.mutations, A.X{`=`, 12, val})
		p.Price = val
		return true
	}
	return false
}

func (p *ProductsMutator) SetInventoryQty(val uint64) bool { //nolint:dupl false positive
	if val != p.InventoryQty {
		p.mutations = append(p.mutations, A.X{`=`, 13, val})
		p.InventoryQty = val
		return true
	}
	return false
}

func (p *ProductsMutator) SetWeightGram(val uint64) bool { //nolint:dupl false positive
	if val != p.WeightGram {
		p.mutations = append(p.mutations, A.X{`=`, 14, val})
		p.WeightGram = val
		return true
	}
	return false
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go

type PromosMutator struct {
	rqStore.Promos
	mutations []A.X
}

func NewPromosMutator(adapter *Tt.Adapter) *PromosMutator {
	return &PromosMutator{Promos: rqStore.Promos{Adapter: adapter}}
}

func (p *PromosMutator) HaveMutation() bool { //nolint:dupl false positive
	return len(p.mutations) > 0
}

// Overwrite all columns, error if not exists
func (p *PromosMutator) DoOverwriteById() bool { //nolint:dupl false positive
	_, err := p.Adapter.Update(p.SpaceName(), p.UniqueIndexId(), A.X{p.Id}, p.ToUpdateArray())
	return !L.IsError(err, `Promos.DoOverwriteById failed: `+p.SpaceName())
}

// Update only mutated, error if not exists, use Find* and Set* methods instead of direct assignment
func (p *PromosMutator) DoUpdateById() bool { //nolint:dupl false positive
	if !p.HaveMutation() {
		return true
	}
	_, err := p.Adapter.Update(p.SpaceName(), p.UniqueIndexId(), A.X{p.Id}, p.mutations)
	return !L.IsError(err, `Promos.DoUpdateById failed: `+p.SpaceName())
}

func (p *PromosMutator) DoDeletePermanentById() bool { //nolint:dupl false positive
	_, err := p.Adapter.Delete(p.SpaceName(), p.UniqueIndexId(), A.X{p.Id})
	return !L.IsError(err, `Promos.DoDeletePermanentById failed: `+p.SpaceName())
}

// func (p *PromosMutator) DoUpsert() bool { //nolint:dupl false positive
//	_, err := p.Adapter.Upsert(p.SpaceName(), p.ToArray(), A.X{
//		A.X{`=`, 0, p.Id},
//		A.X{`=`, 1, p.CreatedAt},
//		A.X{`=`, 2, p.CreatedBy},
//		A.X{`=`, 3, p.UpdatedAt},
//		A.X{`=`, 4, p.UpdatedBy},
//		A.X{`=`, 5, p.DeletedAt},
//		A.X{`=`, 6, p.DeletedBy},
//		A.X{`=`, 7, p.IsDeleted},
//		A.X{`=`, 8, p.RestoredAt},
//		A.X{`=`, 9, p.RestoredBy},
//		A.X{`=`, 10, p.StartAt},
//		A.X{`=`, 11, p.EndAt},
//		A.X{`=`, 12, p.ProductId},
//		A.X{`=`, 13, p.ProductCount},
//		A.X{`=`, 14, p.FreeProductId},
//		A.X{`=`, 15, p.DiscountCount},
//		A.X{`=`, 16, p.DiscountPercent},
//	})
//	return !L.IsError(err, `Promos.DoUpsert failed: `+p.SpaceName())
// }

// insert, error if exists
func (p *PromosMutator) DoInsert() bool { //nolint:dupl false positive
	row, err := p.Adapter.Insert(p.SpaceName(), p.ToArray())
	if err == nil {
		tup := row.Tuples()
		if len(tup) > 0 && len(tup[0]) > 0 && tup[0][0] != nil {
			p.Id = X.ToU(tup[0][0])
		}
	}
	return !L.IsError(err, `Promos.DoInsert failed: `+p.SpaceName())
}

// replace = upsert, only error when there's unique secondary key
func (p *PromosMutator) DoReplace() bool { //nolint:dupl false positive
	_, err := p.Adapter.Replace(p.SpaceName(), p.ToArray())
	return !L.IsError(err, `Promos.DoReplace failed: `+p.SpaceName())
}

func (p *PromosMutator) SetId(val uint64) bool { //nolint:dupl false positive
	if val != p.Id {
		p.mutations = append(p.mutations, A.X{`=`, 0, val})
		p.Id = val
		return true
	}
	return false
}

func (p *PromosMutator) SetCreatedAt(val int64) bool { //nolint:dupl false positive
	if val != p.CreatedAt {
		p.mutations = append(p.mutations, A.X{`=`, 1, val})
		p.CreatedAt = val
		return true
	}
	return false
}

func (p *PromosMutator) SetCreatedBy(val uint64) bool { //nolint:dupl false positive
	if val != p.CreatedBy {
		p.mutations = append(p.mutations, A.X{`=`, 2, val})
		p.CreatedBy = val
		return true
	}
	return false
}

func (p *PromosMutator) SetUpdatedAt(val int64) bool { //nolint:dupl false positive
	if val != p.UpdatedAt {
		p.mutations = append(p.mutations, A.X{`=`, 3, val})
		p.UpdatedAt = val
		return true
	}
	return false
}

func (p *PromosMutator) SetUpdatedBy(val uint64) bool { //nolint:dupl false positive
	if val != p.UpdatedBy {
		p.mutations = append(p.mutations, A.X{`=`, 4, val})
		p.UpdatedBy = val
		return true
	}
	return false
}

func (p *PromosMutator) SetDeletedAt(val int64) bool { //nolint:dupl false positive
	if val != p.DeletedAt {
		p.mutations = append(p.mutations, A.X{`=`, 5, val})
		p.DeletedAt = val
		return true
	}
	return false
}

func (p *PromosMutator) SetDeletedBy(val uint64) bool { //nolint:dupl false positive
	if val != p.DeletedBy {
		p.mutations = append(p.mutations, A.X{`=`, 6, val})
		p.DeletedBy = val
		return true
	}
	return false
}

func (p *PromosMutator) SetIsDeleted(val bool) bool { //nolint:dupl false positive
	if val != p.IsDeleted {
		p.mutations = append(p.mutations, A.X{`=`, 7, val})
		p.IsDeleted = val
		return true
	}
	return false
}

func (p *PromosMutator) SetRestoredAt(val int64) bool { //nolint:dupl false positive
	if val != p.RestoredAt {
		p.mutations = append(p.mutations, A.X{`=`, 8, val})
		p.RestoredAt = val
		return true
	}
	return false
}

func (p *PromosMutator) SetRestoredBy(val uint64) bool { //nolint:dupl false positive
	if val != p.RestoredBy {
		p.mutations = append(p.mutations, A.X{`=`, 9, val})
		p.RestoredBy = val
		return true
	}
	return false
}

func (p *PromosMutator) SetStartAt(val int64) bool { //nolint:dupl false positive
	if val != p.StartAt {
		p.mutations = append(p.mutations, A.X{`=`, 10, val})
		p.StartAt = val
		return true
	}
	return false
}

func (p *PromosMutator) SetEndAt(val int64) bool { //nolint:dupl false positive
	if val != p.EndAt {
		p.mutations = append(p.mutations, A.X{`=`, 11, val})
		p.EndAt = val
		return true
	}
	return false
}

func (p *PromosMutator) SetProductId(val uint64) bool { //nolint:dupl false positive
	if val != p.ProductId {
		p.mutations = append(p.mutations, A.X{`=`, 12, val})
		p.ProductId = val
		return true
	}
	return false
}

func (p *PromosMutator) SetProductCount(val uint64) bool { //nolint:dupl false positive
	if val != p.ProductCount {
		p.mutations = append(p.mutations, A.X{`=`, 13, val})
		p.ProductCount = val
		return true
	}
	return false
}

func (p *PromosMutator) SetFreeProductId(val uint64) bool { //nolint:dupl false positive
	if val != p.FreeProductId {
		p.mutations = append(p.mutations, A.X{`=`, 14, val})
		p.FreeProductId = val
		return true
	}
	return false
}

func (p *PromosMutator) SetDiscountCount(val uint64) bool { //nolint:dupl false positive
	if val != p.DiscountCount {
		p.mutations = append(p.mutations, A.X{`=`, 15, val})
		p.DiscountCount = val
		return true
	}
	return false
}

func (p *PromosMutator) SetDiscountPercent(val float64) bool { //nolint:dupl false positive
	if val != p.DiscountPercent {
		p.mutations = append(p.mutations, A.X{`=`, 16, val})
		p.DiscountPercent = val
		return true
	}
	return false
}

// DO NOT EDIT, will be overwritten by github.com/kokizzu/D/Tt/tarantool_orm_generator.go
