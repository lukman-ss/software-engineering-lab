package codereview

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"
)

// CheckoutImproved is a production-ready checkout implementation.
type CheckoutImproved struct {
	repo        OrderRepository
	products    ProductRepository
	notify      NotificationSender
	logger      Logger
	cartSource  CartSource
	idempotency IdempotencyRepository
	orderCounter atomic.Int64
}

func NewCheckoutImproved(repo OrderRepository, products ProductRepository, notify NotificationSender, logger Logger, cart CartSource, idem IdempotencyRepository) *CheckoutImproved {
	return &CheckoutImproved{
		repo:        repo,
		products:    products,
		notify:      notify,
		logger:      logger,
		cartSource:  cart,
		idempotency: idem,
	}
}

// Checkout processes a checkout request with proper validation, transactions, and idempotency.
func (c *CheckoutImproved) Checkout(ctx context.Context, userID string, idempotencyKey string) (*CheckoutResponse, error) {
	cartItems, err := c.cartSource.GetCart(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get cart failed: %w", err)
	}

	if len(cartItems) == 0 {
		return nil, ErrEmptyCart
	}

	hash := c.hashRequest(cartItems)

	existing, err := c.idempotency.Get(ctx, idempotencyKey)
	if err != nil {
		return nil, fmt.Errorf("check idempotency failed: %w", err)
	}

	if existing != nil {
		if existing.RequestHash != hash {
			return nil, ErrIdempotencyConflict
		}
		if existing.Processed && existing.Response != nil {
			return existing.Response, nil
		}
		return nil, ErrDuplicateRequest
	}

	if _, err := c.idempotency.TryInsert(ctx, idempotencyKey, hash); err != nil {
		return nil, fmt.Errorf("idempotency insert failed: %w", err)
	}

	// Validate and deduct stock atomically per product
	var total int64
	itemsForOrder := make([]OrderItem, 0, len(cartItems))

	for _, item := range cartItems {
		if item.Quantity <= 0 {
			return nil, ErrInvalidQuantity
		}

		product, err := c.products.GetProduct(ctx, item.ProductID)
		if err != nil {
			return nil, ErrProductNotFound
		}

		if product.Stock < item.Quantity {
			return nil, ErrInsufficientStock
		}

		newStock, err := c.products.UpdateStock(ctx, item.ProductID, -item.Quantity)
		if err != nil {
			return nil, fmt.Errorf("stock update failed: %w", err)
		}

		if newStock < 0 {
			return nil, ErrInsufficientStock
		}

		subtotal := int64(item.Quantity) * item.UnitPrice
		total += subtotal
		itemsForOrder = append(itemsForOrder, OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			Subtotal:  subtotal,
		})
	}

	orderID := fmt.Sprintf("order-%d", c.orderCounter.Add(1))
	order := &Order{
		ID:             orderID,
		UserID:         userID,
		IdempotencyKey: idempotencyKey,
		Items:          itemsForOrder,
		TotalAmount:    total,
		CreatedAt:      time.Now(),
	}

	if err := c.repo.CreateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("create order failed: %w", err)
	}

	// Notification after successful order (fire-and-forget)
	_ = c.notify.SendOrderConfirmation(ctx, userID, orderID)

	resp := &CheckoutResponse{Success: true, OrderID: orderID, Total: total}
	_ = c.idempotency.MarkProcessed(ctx, idempotencyKey, resp)

	c.logger.Info(ctx, "checkout completed", "userID", userID, "orderID", orderID, "total", total)

	return resp, nil
}

func (c *CheckoutImproved) hashRequest(items []CartItem) string {
	b, _ := json.Marshal(items)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}