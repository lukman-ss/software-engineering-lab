package service

import "testing"

// TestCalculateDiscount_Weak achieves 100% line coverage but uses weak assertions.
// It fails to assert exact values, edge boundaries, or boolean branches.
func TestCalculateDiscount_Weak(t *testing.T) {
	// Case 1: VIP customer
	res1 := CalculateDiscount(Order{
		TotalAmount: 1200.0,
		ItemCount:   3,
		Tier:        TierVIP,
		HasCoupon:   true,
	})
	if res1.FinalAmount <= 0 {
		t.Errorf("Expected positive final amount, got %v", res1.FinalAmount)
	}

	// Case 2: Premium customer with high amount
	res2 := CalculateDiscount(Order{
		TotalAmount: 600.0,
		ItemCount:   1,
		Tier:        TierPremium,
		HasCoupon:   false,
	})
	if res2.DiscountTotal < 0 {
		t.Errorf("Discount total should not be negative")
	}

	// Case 3: Standard customer with medium amount
	res3 := CalculateDiscount(Order{
		TotalAmount: 150.0,
		ItemCount:   1,
		Tier:        TierStandard,
		HasCoupon:   false,
	})
	if res3.DiscountRate < 0 {
		t.Errorf("Rate should not be negative")
	}

	// Case 4: Standard customer with low amount and coupon
	res4 := CalculateDiscount(Order{
		TotalAmount: 50.0,
		ItemCount:   4,
		Tier:        TierStandard,
		HasCoupon:   true,
	})
	if res4.OriginalTotal != 50.0 {
		t.Errorf("Original total mismatch")
	}
}
