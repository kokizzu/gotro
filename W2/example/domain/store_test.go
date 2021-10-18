package domain

import (
	"testing"

	"github.com/kokizzu/id64"
	"github.com/kokizzu/lexid"
	"github.com/stretchr/testify/assert"
)

func TestCartFlow(t *testing.T) {
	d := NewDomain()
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
	rc := NewRC(sessionToken)

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
	t.Run(`below rule`, func(t *testing.T) {

	})
	t.Run(`equal rule`, func(t *testing.T) {

	})
	t.Run(`above rule`, func(t *testing.T) {

	})
	t.Run(`twice above rule`, func(t *testing.T) {

	})
}

func TestDiscountPercentPromo(t *testing.T) {
	// buy >= N, discount P %

	t.Run(`below rule`, func(t *testing.T) {

	})
	t.Run(`equal rule`, func(t *testing.T) {

	})
	t.Run(`above rule`, func(t *testing.T) {

	})
	t.Run(`twice above rule`, func(t *testing.T) {

	})
}

func TestDiscountDeductItemPromo(t *testing.T) {
	// buy N only pay N-M

	t.Run(`below rule`, func(t *testing.T) {

	})
	t.Run(`equal rule`, func(t *testing.T) {

	})
	t.Run(`above rule`, func(t *testing.T) {

	})
	t.Run(`twice above rule`, func(t *testing.T) {

	})
}
