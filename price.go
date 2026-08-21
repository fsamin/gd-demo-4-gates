package main

// Item is one line of an order.
type Item struct {
	Label string
	Cents int
	Qty   int
}

// discountThresholdCents is the order value from which the volume discount
// applies. It is a business rule, not a magic number — changing it changes what
// every customer pays, which is why the tests pin it and its boundary.
const discountThresholdCents = 10_000

// discountPercent is the volume discount granted from the threshold on.
const discountPercent = 10

// Subtotal is the order value before any discount.
func Subtotal(items []Item) int {
	total := 0
	for _, item := range items {
		total += item.Cents * item.Qty
	}
	return total
}

// Total applies the volume discount to the subtotal. The discount is granted on
// the whole order, not only on the part above the threshold, and an order has to
// pass the threshold to earn it.
func Total(items []Item) int {
	subtotal := Subtotal(items)
	if subtotal <= discountThresholdCents {
		return subtotal
	}
	return subtotal - subtotal*discountPercent/100
}
