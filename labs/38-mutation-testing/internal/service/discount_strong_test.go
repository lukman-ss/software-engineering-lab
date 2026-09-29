package service

import "testing"

// TestCalculateDiscount_Strong tests all boundary conditions, rates, free shipping flags, and coupon additions.
func TestCalculateDiscount_Strong(t *testing.T) {
	tests := []struct {
		name                 string
		order                Order
		expectedRate         float64
		expectedDiscount     float64
		expectedFinal        float64
		expectedFreeShipping bool
	}{
		{
			name: "VIP customer gets 20% discount and free shipping regardless of amount",
			order: Order{
				TotalAmount: 100.0,
				ItemCount:   1,
				Tier:        TierVIP,
				HasCoupon:   false,
			},
			expectedRate:         0.20,
			expectedDiscount:     20.0,
			expectedFinal:        80.0,
			expectedFreeShipping: true,
		},
		{
			name: "High spender non-VIP gets 20% discount",
			order: Order{
				TotalAmount: 1000.0,
				ItemCount:   1,
				Tier:        TierStandard,
				HasCoupon:   false,
			},
			expectedRate:         0.20,
			expectedDiscount:     200.0,
			expectedFinal:        800.0,
			expectedFreeShipping: true,
		},
		{
			name: "Premium customer with amount 500 gets 10% discount",
			order: Order{
				TotalAmount: 500.0,
				ItemCount:   2,
				Tier:        TierPremium,
				HasCoupon:   false,
			},
			expectedRate:         0.10,
			expectedDiscount:     50.0,
			expectedFinal:        450.0,
			expectedFreeShipping: true,
		},
		{
			name: "Standard customer with amount 100 gets 5% discount",
			order: Order{
				TotalAmount: 100.0,
				ItemCount:   2,
				Tier:        TierStandard,
				HasCoupon:   false,
			},
			expectedRate:         0.05,
			expectedDiscount:     5.0,
			expectedFinal:        95.0,
			expectedFreeShipping: false,
		},
		{
			name: "Coupon with >2 items adds 5% rate boost",
			order: Order{
				TotalAmount: 100.0,
				ItemCount:   3,
				Tier:        TierStandard,
				HasCoupon:   true,
			},
			expectedRate:         0.10,
			expectedDiscount:     10.0,
			expectedFinal:        90.0,
			expectedFreeShipping: false,
		},
		{
			name: "Coupon with <=2 items does NOT add extra rate",
			order: Order{
				TotalAmount: 100.0,
				ItemCount:   2,
				Tier:        TierStandard,
				HasCoupon:   true,
			},
			expectedRate:         0.05,
			expectedDiscount:     5.0,
			expectedFinal:        95.0,
			expectedFreeShipping: false,
		},
		{
			name: "Free shipping threshold at exact 200 final amount",
			order: Order{
				TotalAmount: 200.0,
				ItemCount:   1,
				Tier:        TierStandard,
				HasCoupon:   false,
			},
			expectedRate:         0.05,
			expectedDiscount:     10.0,
			expectedFinal:        190.0,
			expectedFreeShipping: false,
		},
		{
			name: "Degenerate zero amount and zero items",
			order: Order{
				TotalAmount: 0.0,
				ItemCount:   0,
				Tier:        TierStandard,
				HasCoupon:   false,
			},
			expectedRate:         0.0,
			expectedDiscount:     0.0,
			expectedFinal:        0.0,
			expectedFreeShipping: false,
		},
		{
			name: "Below-boundary amount 499.99 for Premium tier",
			order: Order{
				TotalAmount: 499.99,
				ItemCount:   2,
				Tier:        TierPremium,
				HasCoupon:   false,
			},
			expectedRate:         0.05,
			expectedDiscount:     24.9995,
			expectedFinal:        474.9905,
			expectedFreeShipping: true,
		},
		{
			name: "Below-boundary amount 99.99 for Standard tier",
			order: Order{
				TotalAmount: 99.99,
				ItemCount:   1,
				Tier:        TierStandard,
				HasCoupon:   false,
			},
			expectedRate:         0.0,
			expectedDiscount:     0.0,
			expectedFinal:        99.99,
			expectedFreeShipping: false,
		},
		{
			name: "Coupon true with zero items",
			order: Order{
				TotalAmount: 150.0,
				ItemCount:   0,
				Tier:        TierStandard,
				HasCoupon:   true,
			},
			expectedRate:         0.05,
			expectedDiscount:     7.5,
			expectedFinal:        142.5,
			expectedFreeShipping: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := CalculateDiscount(tt.order)
			if res.DiscountRate != tt.expectedRate {
				t.Errorf("DiscountRate: got %.2f, want %.2f", res.DiscountRate, tt.expectedRate)
			}
			if res.DiscountTotal != tt.expectedDiscount {
				t.Errorf("DiscountTotal: got %.2f, want %.2f", res.DiscountTotal, tt.expectedDiscount)
			}
			if res.FinalAmount != tt.expectedFinal {
				t.Errorf("FinalAmount: got %.2f, want %.2f", res.FinalAmount, tt.expectedFinal)
			}
			if res.FreeShipping != tt.expectedFreeShipping {
				t.Errorf("FreeShipping: got %v, want %v", res.FreeShipping, tt.expectedFreeShipping)
			}
		})
	}
}
