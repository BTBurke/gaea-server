package routes

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestSubTotal(t *testing.T) {
	total := calcSubTotal(5, decimal.NewFromFloat(10.13))
	if !total.Equals(decimal.NewFromFloat(50.65)) {
		t.Fatalf("Subtotal incorrect: Expected 50.65 Got: %v", total)
	}
}

func TestCasePricing(t *testing.T) {
	total := calcSubTotalCasePricing(5, decimal.NewFromFloat(10.00), 4, 20)
	expect := decimal.NewFromFloat(float64(5*10 + 1*2))
	if !total.Equals(expect) {
		t.Fatalf("Subtotal incorrect: Expected: %v Got: %v", expect, total)
	}
}
