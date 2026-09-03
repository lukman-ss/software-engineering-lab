package codereview

import (
	"context"
	"fmt"
	"time"
)

// CheckoutNaive demonstrates multiple code review issues in checkout processing.
// This implementation is INTENTIONALLY BROKEN for educational purposes.
// DO NOT USE IN PRODUCTION. See CheckoutImproved for the corrected implementation.
type CheckoutNaive struct {
	repo       OrderRepository
	products   ProductRepository
	notify     NotificationSender
	logger     Logger
	cartSource CartSource
}

func NewCheckoutNaive(repo OrderRepository, products ProductRepository, notify NotificationSender, logger Logger, cart CartSource) *CheckoutNaive {
	return &CheckoutNaive{
		repo:       repo,
		products:   products,
		notify:     notify,
		logger:     logger,
		cartSource: cart,
	}
}

// Checkout processes a checkout request WITHOUT proper validation, transactions, or idempotency.
// WARNING: This implementation has multiple intentional bugs for code review training.
func (c *CheckoutNaive) Checkout(ctx context.Context, userID string) (*CheckoutResponse, error) {
	cartItems, err := c.cartSource.GetCart(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get cart failed: %w", err)
	}

	// BUG #1: No validation for empty cart
	// if len(cartItems) == 0 { return nil, ErrEmptyCart }

	var total int64
	for _, item := range cartItems {
		// BUG #2: N+1 query pattern - fetching product individually in a loop
		product, err := c.products.GetProduct(ctx, item.ProductID)
		if err != nil {
			// BUG #3: Missing product existence handling - continues with nil product
			continue
		}

		// BUG #4: Race condition - non-atomic read-modify-write on stock
		// BUG #5: Stock can become negative (overselling)
		// BUG #6: No validation that stock >= quantity before deduction
		product.Stock -= item.Quantity
		_, _ = c.products.UpdateStock(ctx, item.ProductID, -item.Quantity)

		total += int64(item.Quantity) * int64(item.UnitPrice)
	}

	// BUG #7: No transaction - if order creation fails, stock already changed
	orderID := fmt.Sprintf("order-%d", time.Now().UnixNano())
	order := &Order{
		ID:          orderID,
		UserID:      userID,
		Items:       c.toOrderItems(cartItems),
		TotalAmount: total,
		CreatedAt:   time.Now(),
	}

	_ = c.repo.CreateOrder(ctx, order)

	// BUG #8: Notification as side effect - not tied to order success
	_ = c.notify.SendOrderConfirmation(ctx, userID, orderID)

	// BUG #9: Poor error propagation - always returns success
	return &CheckoutResponse{Success: true, OrderID: orderID, Total: total}, nil
}

func (c *CheckoutNaive) toOrderItems(items []CartItem) []OrderItem {
	ords := make([]OrderItem, len(items))
	for i, item := range items {
		ords[i] = OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			Subtotal:  int64(item.Quantity) * item.UnitPrice,
		}
	}
	return ords
}