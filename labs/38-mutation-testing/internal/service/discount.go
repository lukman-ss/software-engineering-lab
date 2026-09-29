package service

type CustomerTier string

const (
	TierStandard CustomerTier = "STANDARD"
	TierPremium  CustomerTier = "PREMIUM"
	TierVIP      CustomerTier = "VIP"
)

type Order struct {
	TotalAmount float64
	ItemCount   int
	Tier        CustomerTier
	HasCoupon   bool
}

type DiscountResult struct {
	OriginalTotal float64
	DiscountRate  float64
	DiscountTotal float64
	FinalAmount   float64
	FreeShipping  bool
}

// CalculateDiscount computes discount and shipping based on customer tier, amount, and item count.
func CalculateDiscount(order Order) DiscountResult {
	rate := 0.0

	// Relational & boolean: VIP or high spender gets primary rate
	if order.Tier == TierVIP || order.TotalAmount >= 1000.0 {
		rate = 0.20
	} else if order.Tier == TierPremium && order.TotalAmount >= 500.0 {
		rate = 0.10
	} else if order.TotalAmount >= 100.0 {
		rate = 0.05
	}

	// Boundary & arithmetic: Coupon adds extra discount if minimum items present
	if order.HasCoupon && order.ItemCount > 2 {
		rate = rate + 0.05
	}

	discountAmount := order.TotalAmount * rate
	finalAmount := order.TotalAmount - discountAmount

	// Free shipping threshold: after discount, VIP or amount >= 200
	freeShipping := (order.Tier == TierVIP) || (finalAmount >= 200.0)

	return DiscountResult{
		OriginalTotal: order.TotalAmount,
		DiscountRate:  rate,
		DiscountTotal: discountAmount,
		FinalAmount:   finalAmount,
		FreeShipping:  freeShipping,
	}
}
