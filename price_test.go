package main

import "testing"

// A small order pays its subtotal: the discount rewards volume, and granting it
// to everyone would just be a price cut.
func TestTotalBelowThresholdIsSubtotal(t *testing.T) {
	items := []Item{{Label: "widget", Cents: 2_500, Qty: 2}}
	if got, want := Total(items), 5_000; got != want {
		t.Errorf("Total = %d, want %d (no discount below %d)", got, want, discountThresholdCents)
	}
}

// The threshold itself earns the discount. A customer landing exactly on it
// must not be told to spend more to get what they were promised — this is the
// boundary the rule is most often broken on.
func TestTotalAtThresholdIsDiscounted(t *testing.T) {
	items := []Item{{Label: "widget", Cents: discountThresholdCents, Qty: 1}}
	if got, want := Total(items), 9_000; got != want {
		t.Errorf("Total = %d, want %d (the discount applies from %d on)", got, want, discountThresholdCents)
	}
}

// The discount is granted on the whole order, not only on the part above the
// threshold: 200.00 discounted by 10% is 180.00, not 190.00.
func TestTotalDiscountsTheWholeOrder(t *testing.T) {
	items := []Item{{Label: "widget", Cents: 20_000, Qty: 1}}
	if got, want := Total(items), 18_000; got != want {
		t.Errorf("Total = %d, want %d (the discount covers the whole order)", got, want)
	}
}

// Quantities count. A subtotal that ignored them would make every multi-unit
// order free money.
func TestSubtotalCountsQuantities(t *testing.T) {
	items := []Item{
		{Label: "widget", Cents: 1_000, Qty: 3},
		{Label: "gadget", Cents: 500, Qty: 2},
	}
	if got, want := Subtotal(items), 4_000; got != want {
		t.Errorf("Subtotal = %d, want %d", got, want)
	}
}
