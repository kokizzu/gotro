package domain

import (
	"testing"

	"github.com/hexops/autogold"

	"github.com/kokizzu/gotro/W2/example/model/mStore/rqStore"
	"github.com/kokizzu/id64"
	"github.com/kokizzu/lexid"
	"github.com/stretchr/testify/assert"
)

func prepareTestUser(d *Domain, t *testing.T) RequestCommon {
	name := id64.ID().String()
	pass := lexid.ID()
	email := name + testDomain
	{
		in := &UserRegister_In{
			Email:    email,
			Password: pass,
			UserName: name,
		}
		out := d.UserRegister(in)
		assert.Empty(t, out.Error)
	}
	sessionToken := ``
	{
		in := &UserLogin_In{
			Email:    email,
			Password: pass,
		}
		out := d.UserLogin(in)
		assert.Empty(t, out.Error)
		sessionToken = out.SessionToken
	}
	return NewRC(sessionToken)
}

func TestCartFlow(t *testing.T) {
	d := NewDomain()
	rc := prepareTestUser(d, t)

	t.Run(`add new item to cart`, func(t *testing.T) {
		in := &StoreCartItemsAdd_In{
			RequestCommon: rc,
			ProductId:     1,
			DeltaQty:      4,
		}
		out := d.StoreCartItemsAdd(in)
		assert.Empty(t, out.Error)
		assert.Equal(t, len(out.CartItems), 1)
		assert.Equal(t, out.CartItems[0].Qty, in.DeltaQty)
	})
	t.Run(`add other item to cart`, func(t *testing.T) {
		in := &StoreCartItemsAdd_In{
			RequestCommon: rc,
			ProductId:     2,
			DeltaQty:      5,
		}
		out := d.StoreCartItemsAdd(in)
		assert.Empty(t, out.Error)
		assert.Equal(t, len(out.CartItems), 2)
		assert.Equal(t, out.CartItems[1].Qty, in.DeltaQty)
	})
	t.Run(`add existing item to cart`, func(t *testing.T) {
		in := &StoreCartItemsAdd_In{
			RequestCommon: rc,
			ProductId:     1,
			DeltaQty:      2,
		}
		out := d.StoreCartItemsAdd(in)
		assert.Empty(t, out.Error)
		assert.Equal(t, len(out.CartItems), 2)
		assert.Equal(t, out.CartItems[0].Qty, 4+in.DeltaQty)
	})
	t.Run(`remove existing item to cart`, func(t *testing.T) {
		in := &StoreCartItemsAdd_In{
			RequestCommon: rc,
			ProductId:     1,
			DeltaQty:      -3,
		}
		out := d.StoreCartItemsAdd(in)
		assert.Empty(t, out.Error)
		assert.Equal(t, len(out.CartItems), 2)
		assert.Equal(t, out.CartItems[0].Qty, 4+2+in.DeltaQty)
	})
	t.Run(`remove nonexistent item to cart`, func(t *testing.T) {
		in := &StoreCartItemsAdd_In{
			RequestCommon: rc,
			ProductId:     3,
			DeltaQty:      -4,
		}
		out := d.StoreCartItemsAdd(in)
		assert.NotEmpty(t, out.Error)
	})
	t.Run(`remove more than available to cart`, func(t *testing.T) {
		in := &StoreCartItemsAdd_In{
			RequestCommon: rc,
			ProductId:     2,
			DeltaQty:      -9,
		}
		out := d.StoreCartItemsAdd(in)
		assert.Empty(t, out.Error)
		assert.Equal(t, len(out.CartItems), 2)
		assert.Equal(t, out.CartItems[0].Qty, int64(4+2-3))
		assert.Equal(t, out.CartItems[1].Qty, int64(0))
	})
	t.Run(`remove when already zero`, func(t *testing.T) {
		in := &StoreCartItemsAdd_In{
			RequestCommon: rc,
			ProductId:     2,
			DeltaQty:      -9,
		}
		out := d.StoreCartItemsAdd(in)
		assert.NotEmpty(t, out.Error)
	})
	t.Run(`add more than max when not yet max`, func(t *testing.T) {
		in := &StoreCartItemsAdd_In{
			RequestCommon: rc,
			ProductId:     3,
			DeltaQty:      11,
		}
		out := d.StoreCartItemsAdd(in)
		assert.Empty(t, out.Error)
		assert.Equal(t, out.CartItems[0].Qty, int64(4+2-3))
		assert.Equal(t, out.CartItems[1].Qty, int64(0))
		assert.Equal(t, out.CartItems[2].Qty, int64(10)) // max stock of this item
	})
	t.Run(`add more when already max`, func(t *testing.T) {
		in := &StoreCartItemsAdd_In{
			RequestCommon: rc,
			ProductId:     3,
			DeltaQty:      11,
		}
		out := d.StoreCartItemsAdd(in)
		assert.NotEmpty(t, out.Error)
	})
}

func TestFreeItemPromo(t *testing.T) {
	// buy N product X, got free M product Y
	// buy 1 macbook, got 1 free rapsberry pi (stock 2)
	d := NewDomain()
	rc := prepareTestUser(d, t)

	t.Run(`equal rule`, func(t *testing.T) {
		{
			in := &StoreCartItemsAdd_In{
				RequestCommon: rc,
				ProductId:     2,
				DeltaQty:      1,
			}
			_ = d.StoreCartItemsAdd(in)
		}
		in := &StoreInvoice_In{
			RequestCommon: rc,
			Recalculate:   true,
		}
		out := d.StoreInvoice(in)
		assert.Empty(t, out.Error)
		for _, ci := range out.CartItems {
			ci.Id = 0
			ci.OwnerId = 0
		}
		want := autogold.Want("A1", StoreInvoice_Out{
			CartItems: []*rqStore.CartItems{
				{
					ProductId: 2,
					NameCopy:  "MacBook Pro",
					PriceCopy: 539999,
					Qty:       1,
					SubTotal:  539999,
				},
				{
					ProductId: 4,
					NameCopy:  "Raspberry Pi B",
					PriceCopy: 3000,
					Qty:       1,
					Discount:  3000,
					Info: `got 1 free (total: 1) every purchase of 1 MacBook Pro
`,
				},
			},
			Invoice: rqStore.Invoices{
				TotalDiscount: 3000,
				TotalPaid:     539999,
			},
		})
		want.Equal(t, out)
	})
	t.Run(`twice above rule`, func(t *testing.T) {
		{
			in := &StoreCartItemsAdd_In{
				RequestCommon: rc,
				ProductId:     2,
				DeltaQty:      1,
			}
			_ = d.StoreCartItemsAdd(in)
		}
		in := &StoreInvoice_In{
			RequestCommon: rc,
			Recalculate:   true,
		}
		out := d.StoreInvoice(in)
		assert.Empty(t, out.Error)
		for _, ci := range out.CartItems {
			ci.Id = 0
			ci.OwnerId = 0
		}
		want := autogold.Want("A2", StoreInvoice_Out{
			CartItems: []*rqStore.CartItems{
				{
					ProductId: 2,
					NameCopy:  "MacBook Pro",
					PriceCopy: 539999,
					Qty:       2,
					SubTotal:  1079998,
				},
				{
					ProductId: 4,
					NameCopy:  "Raspberry Pi B",
					PriceCopy: 3000,
					Qty:       2,
					Discount:  6000,
					Info: `got 1 free (total: 2) every purchase of 1 MacBook Pro
`,
				},
			},
			Invoice: rqStore.Invoices{
				TotalDiscount: 6000,
				TotalPaid:     1079998,
			},
		})
		want.Equal(t, out)
	})
	t.Run(`twice above rule but already purchase 1`, func(t *testing.T) {
		{
			in := &StoreCartItemsAdd_In{
				RequestCommon: rc,
				ProductId:     4,
				DeltaQty:      1,
			}
			_ = d.StoreCartItemsAdd(in)
		}
		in := &StoreInvoice_In{
			RequestCommon: rc,
			Recalculate:   true,
		}
		out := d.StoreInvoice(in)
		assert.Empty(t, out.Error)
		for _, ci := range out.CartItems {
			ci.Id = 0
			ci.OwnerId = 0
		}
		want := autogold.Want("A3", StoreInvoice_Out{
			CartItems: []*rqStore.CartItems{
				&rqStore.CartItems{
					ProductId: 2,
					NameCopy:  "MacBook Pro",
					PriceCopy: 539999,
					Qty:       2,
					SubTotal:  1079998,
				},
				&rqStore.CartItems{
					ProductId: 4,
					NameCopy:  "Raspberry Pi B",
					PriceCopy: 3000,
					Qty:       2,
					Discount:  3000,
					SubTotal:  3000,
					Info: `got 1 free (total: 2) every purchase of 1 MacBook Pro
but we don't have enough free item in inventory (missing: 1)
`,
				},
			},
			Invoice: rqStore.Invoices{
				TotalDiscount: 3000,
				TotalPaid:     1082998,
			},
		})
		want.Equal(t, out)
	})
}

func TestDiscountDeductItemPromo(t *testing.T) {
	// buy N only pay N-M
	// buy 3 google home, only need to pay 2
	d := NewDomain()
	rc := prepareTestUser(d, t)

	t.Run(`below rule`, func(t *testing.T) {
		{
			in := &StoreCartItemsAdd_In{
				RequestCommon: rc,
				ProductId:     1,
				DeltaQty:      2,
			}
			_ = d.StoreCartItemsAdd(in)
		}
		in := &StoreInvoice_In{
			RequestCommon: rc,
			Recalculate:   true,
		}
		out := d.StoreInvoice(in)
		assert.Empty(t, out.Error)
		for _, ci := range out.CartItems {
			ci.Id = 0
			ci.OwnerId = 0
		}
		want := autogold.Want("B1", StoreInvoice_Out{
			CartItems: []*rqStore.CartItems{
				{
					ProductId: 1,
					NameCopy:  "Google Home",
					PriceCopy: 4999,
					Qty:       2,
					SubTotal:  9998,
				},
			},
			Invoice: rqStore.Invoices{TotalPaid: 9998},
		})
		want.Equal(t, out)
	})
	t.Run(`equal rule`, func(t *testing.T) {
		{
			in := &StoreCartItemsAdd_In{
				RequestCommon: rc,
				ProductId:     1,
				DeltaQty:      1,
			}
			_ = d.StoreCartItemsAdd(in)
		}
		in := &StoreInvoice_In{
			RequestCommon: rc,
			Recalculate:   true,
		}
		out := d.StoreInvoice(in)
		assert.Empty(t, out.Error)
		for _, ci := range out.CartItems {
			ci.Id = 0
			ci.OwnerId = 0
		}
		want := autogold.Want("B2", StoreInvoice_Out{
			CartItems: []*rqStore.CartItems{
				&rqStore.CartItems{
					ProductId: 1,
					NameCopy:  "Google Home",
					PriceCopy: 4999,
					Qty:       3,
					Discount:  4999,
					SubTotal:  9998,
					Info:      "discount 1 for every 3 purchase\n",
				},
			},
			Invoice: rqStore.Invoices{
				TotalDiscount: 4999,
				TotalPaid:     9998,
			},
		})
		want.Equal(t, out)
	})
	t.Run(`above rule`, func(t *testing.T) {
		{
			in := &StoreCartItemsAdd_In{
				RequestCommon: rc,
				ProductId:     1,
				DeltaQty:      1,
			}
			_ = d.StoreCartItemsAdd(in)
		}
		in := &StoreInvoice_In{
			RequestCommon: rc,
			Recalculate:   true,
		}
		out := d.StoreInvoice(in)
		assert.Empty(t, out.Error)
		for _, ci := range out.CartItems {
			ci.Id = 0
			ci.OwnerId = 0
		}
		want := autogold.Want("B3", StoreInvoice_Out{
			CartItems: []*rqStore.CartItems{
				&rqStore.CartItems{
					ProductId: 1,
					NameCopy:  "Google Home",
					PriceCopy: 4999,
					Qty:       4,
					Discount:  4999,
					SubTotal:  14997,
					Info:      "discount 1 for every 3 purchase\n",
				},
			},
			Invoice: rqStore.Invoices{
				TotalDiscount: 4999,
				TotalPaid:     14997,
			},
		})
		want.Equal(t, out)
	})
	t.Run(`twice above rule`, func(t *testing.T) {
		{
			in := &StoreCartItemsAdd_In{
				RequestCommon: rc,
				ProductId:     1,
				DeltaQty:      2,
			}
			_ = d.StoreCartItemsAdd(in)
		}
		in := &StoreInvoice_In{
			RequestCommon: rc,
			Recalculate:   true,
		}
		out := d.StoreInvoice(in)
		assert.Empty(t, out.Error)
		for _, ci := range out.CartItems {
			ci.Id = 0
			ci.OwnerId = 0
		}
		want := autogold.Want("B4", StoreInvoice_Out{
			CartItems: []*rqStore.CartItems{
				&rqStore.CartItems{
					ProductId: 1,
					NameCopy:  "Google Home",
					PriceCopy: 4999,
					Qty:       6,
					Discount:  9998,
					SubTotal:  19996,
					Info:      "discount 1 for every 3 purchase\n",
				},
			},
			Invoice: rqStore.Invoices{
				TotalDiscount: 9998,
				TotalPaid:     19996,
			},
		})
		want.Equal(t, out)
	})
}
func TestDiscountPercentPromo(t *testing.T) {
	// buy >= N, discount P %
	d := NewDomain()
	rc := prepareTestUser(d, t)

	t.Run(`below rule`, func(t *testing.T) {
		{
			in := &StoreCartItemsAdd_In{
				RequestCommon: rc,
				ProductId:     3,
				DeltaQty:      2,
			}
			_ = d.StoreCartItemsAdd(in)
		}
		in := &StoreInvoice_In{
			RequestCommon: rc,
			Recalculate:   true,
		}
		out := d.StoreInvoice(in)
		assert.Empty(t, out.Error)
		for _, ci := range out.CartItems {
			ci.Id = 0
			ci.OwnerId = 0
		}
		want := autogold.Want("C1", StoreInvoice_Out{
			CartItems: []*rqStore.CartItems{
				&rqStore.CartItems{
					ProductId: 3,
					NameCopy:  "Alexa Speaker",
					PriceCopy: 10950,
					Qty:       2,
					SubTotal:  21900,
				},
			},
			Invoice: rqStore.Invoices{TotalPaid: 21900},
		})
		want.Equal(t, out)
	})
	t.Run(`equal rule`, func(t *testing.T) {
		{
			in := &StoreCartItemsAdd_In{
				RequestCommon: rc,
				ProductId:     3,
				DeltaQty:      1,
			}
			_ = d.StoreCartItemsAdd(in)
		}
		in := &StoreInvoice_In{
			RequestCommon: rc,
			Recalculate:   true,
		}
		out := d.StoreInvoice(in)
		assert.Empty(t, out.Error)
		for _, ci := range out.CartItems {
			ci.Id = 0
			ci.OwnerId = 0
		}
		want := autogold.Want("C2", StoreInvoice_Out{
			CartItems: []*rqStore.CartItems{
				&rqStore.CartItems{
					ProductId: 3,
					NameCopy:  "Alexa Speaker",
					PriceCopy: 10950,
					Qty:       3,
					Discount:  3285,
					SubTotal:  29565,
					Info:      "discount 10% for 3 purchase\n",
				},
			},
			Invoice: rqStore.Invoices{
				TotalDiscount: 3285,
				TotalPaid:     29565,
			},
		})
		want.Equal(t, out)

	})
	t.Run(`above rule`, func(t *testing.T) {
		{
			in := &StoreCartItemsAdd_In{
				RequestCommon: rc,
				ProductId:     3,
				DeltaQty:      1,
			}
			_ = d.StoreCartItemsAdd(in)
		}
		in := &StoreInvoice_In{
			RequestCommon: rc,
			Recalculate:   true,
		}
		out := d.StoreInvoice(in)
		assert.Empty(t, out.Error)
		for _, ci := range out.CartItems {
			ci.Id = 0
			ci.OwnerId = 0
		}
		want := autogold.Want("C3", StoreInvoice_Out{
			CartItems: []*rqStore.CartItems{
				&rqStore.CartItems{
					ProductId: 3,
					NameCopy:  "Alexa Speaker",
					PriceCopy: 10950,
					Qty:       4,
					Discount:  4380,
					SubTotal:  39420,
					Info:      "discount 10% for 3 purchase\n",
				},
			},
			Invoice: rqStore.Invoices{
				TotalDiscount: 4380,
				TotalPaid:     39420,
			},
		})
		want.Equal(t, out)

	})
	t.Run(`twice above rule`, func(t *testing.T) {

		{
			in := &StoreCartItemsAdd_In{
				RequestCommon: rc,
				ProductId:     3,
				DeltaQty:      1,
			}
			_ = d.StoreCartItemsAdd(in)
		}
		in := &StoreInvoice_In{
			RequestCommon: rc,
			Recalculate:   true,
		}
		out := d.StoreInvoice(in)
		assert.Empty(t, out.Error)
		for _, ci := range out.CartItems {
			ci.Id = 0
			ci.OwnerId = 0
		}
		want := autogold.Want("C4", StoreInvoice_Out{
			CartItems: []*rqStore.CartItems{
				&rqStore.CartItems{
					ProductId: 3,
					NameCopy:  "Alexa Speaker",
					PriceCopy: 10950,
					Qty:       5,
					Discount:  5475,
					SubTotal:  49275,
					Info:      "discount 10% for 3 purchase\n",
				},
			},
			Invoice: rqStore.Invoices{
				TotalDiscount: 5475,
				TotalPaid:     49275,
			},
		})
		want.Equal(t, out)
	})
}
