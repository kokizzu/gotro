package domain

import "testing"

func TestCartFlow(t *testing.T) {
	t.Run(`add new item to cart`, func(t *testing.T) {

	})
	t.Run(`add other item to cart`, func(t *testing.T) {

	})
	t.Run(`add existing item to cart`, func(t *testing.T) {

	})
	t.Run(`remove existing item to cart`, func(t *testing.T) {

	})
	t.Run(`remove nonexistent item to cart`, func(t *testing.T) {

	})
	t.Run(`remove more than available to cart`, func(t *testing.T) {

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
