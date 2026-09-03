package codereview

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"
)

// CheckoutImproved demonstrates a safer checkout implementation
// after addressing the review findings in this lab.
// Test-only ID generator: atomic.Int64 produces unique IDs within a single process.
// Production systems should use UUID/ULID, database sequences, or Snowflake-style IDs.
type CheckoutImproved struct {
	repo         OrderRepository
	products     ProductRepository
	notify       NotificationSender
	logger       Logger
	cartSource   CartSource
	idempotency  IdempotencyRepository
	txManager    TransactionManager
	orderCounter atomic.Int64
}

func NewCheckoutImproved(
	repo OrderRepository,
	products ProductRepository,
	notify NotificationSender,
	logger Logger,
	cart CartSource,
	idem IdempotencyRepository,
	txManager TransactionManager,
) *CheckoutImproved {
	return &CheckoutImproved{
		repo:        repo,
		products:    products,
		notify:      notify,
		logger:      logger,
		cartSource:  cart,
		idempotency: idem,
		txManager:   txManager,
	}
}

// Checkout processes a checkout request with proper validation, transactions, and idempotency.
func (c *CheckoutImproved) Checkout(ctx context.Context, principal Principal, cmd CheckoutCommand) (*CheckoutResponse, error) {
	if principal.UserID != cmd.CartOwnerID {
		c.logger.Error(ctx, "authorization failed", "userID", principal.UserID, "cartOwnerID", cmd.CartOwnerID)
		return nil, ErrForbidden
	}

	if strings.TrimSpace(cmd.IdempotencyKey) == "" {
		return nil, ErrInvalidIdempotencyKey
	}

	scopedKey := fmt.Sprintf("checkout:%s:%s", principal.UserID, cmd.IdempotencyKey)

	// Check if already completed first before relying on current mutable cart state.
	// In production, an idempotency record stores the committed result. If completed, return cached response directly.
	if existing, err := c.idempotency.Get(ctx, scopedKey); err == nil && existing != nil {
		if existing.Status == IdempotencyStatusCompleted && existing.Response != nil {
			return existing.Response, nil
		}
	}

	cartItems, err := c.cartSource.GetCart(ctx, cmd.CartOwnerID)
	if err != nil {
		c.logger.Error(ctx, "load cart failed", "userID", principal.UserID, "error", err.Error())
		return nil, fmt.Errorf("load cart: %w", err)
	}

	if len(cartItems) == 0 {
		return nil, ErrEmptyCart
	}

	seen := make(map[string]struct{})
	productIDs := make([]string, 0, len(cartItems))
	for _, item := range cartItems {
		if item.Quantity <= 0 {
			return nil, ErrInvalidQuantity
		}
		if _, ok := seen[item.ProductID]; !ok {
			seen[item.ProductID] = struct{}{}
			productIDs = append(productIDs, item.ProductID)
		}
	}

	// Hash should be based ONLY on the immutable command payload, not mutable server state like cart items
	hash := c.hashRequest(cmd)

	status, cachedResp, err := c.idempotency.Claim(ctx, scopedKey, hash)
	if err != nil {
		c.logger.Error(ctx, "idempotency claim failed", "userID", principal.UserID, "idempotencyKey", scopedKey, "error", err.Error())
		return nil, err
	}
	if status == IdempotencyStatusCompleted && cachedResp != nil {
		return cachedResp, nil
	}

	productsMap, err := c.products.GetProducts(ctx, productIDs)
	if err != nil {
		if relErr := c.idempotency.Release(ctx, scopedKey); relErr != nil {
			c.logger.Error(ctx, "failed to release idempotency after load products failure", "userID", principal.UserID, "idempotencyKey", scopedKey, "error", relErr.Error())
			err = errors.Join(err, relErr)
		}
		c.logger.Error(ctx, "load products failed", "userID", principal.UserID, "error", err.Error())
		if errors.Is(err, ErrProductNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, fmt.Errorf("load products: %w", err)
	}

	var total int64
	itemsForOrder := make([]OrderItem, 0, len(cartItems))
	for _, item := range cartItems {
		product, exists := productsMap[item.ProductID]
		if !exists {
			if relErr := c.idempotency.Release(ctx, scopedKey); relErr != nil {
				c.logger.Error(ctx, "failed to release idempotency after missing product in batch", "userID", principal.UserID, "idempotencyKey", scopedKey, "error", relErr.Error())
			}
			c.logger.Error(ctx, "product missing in batch", "userID", principal.UserID, "productID", item.ProductID)
			return nil, ErrProductNotFound
		}

		subtotal := int64(item.Quantity) * product.UnitPrice
		total += subtotal
		itemsForOrder = append(itemsForOrder, OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: product.UnitPrice,
			Subtotal:  subtotal,
		})
	}

	// Demo-only ID generation.
	// Production systems should use a database sequence, UUID/ULID, Snowflake-style ID,
	// or another distributed-safe strategy (atomic is only safe within a single process).
	orderID := fmt.Sprintf("order-%d", c.orderCounter.Add(1))
	order := &Order{
		ID:             orderID,
		UserID:         principal.UserID,
		IdempotencyKey: scopedKey,
		Items:          itemsForOrder,
		TotalAmount:    total,
		CreatedAt:      time.Now(),
	}

	txErr := c.txManager.WithinTransaction(ctx, func(tx CheckoutTx) error {
		for _, item := range cartItems {
			if reserveErr := tx.ReserveStock(ctx, item.ProductID, item.Quantity); reserveErr != nil {
				return reserveErr
			}
		}
		if createErr := tx.CreateOrder(ctx, order); createErr != nil {
			return createErr
		}
		return nil
	})

	if txErr != nil {
		if relErr := c.idempotency.Release(ctx, scopedKey); relErr != nil {
			c.logger.Error(ctx, "failed to release idempotency after transaction failure", "userID", principal.UserID, "idempotencyKey", scopedKey, "error", relErr.Error(), "originalError", txErr.Error())
			txErr = errors.Join(txErr, relErr)
		}
		c.logger.Error(ctx, "checkout transaction failed", "userID", principal.UserID, "idempotencyKey", scopedKey, "error", txErr.Error())
		if errors.Is(txErr, ErrInsufficientStock) || errors.Is(txErr, ErrProductNotFound) || errors.Is(txErr, ErrInvalidQuantity) {
			return nil, txErr
		}
		return nil, fmt.Errorf("checkout transaction: %w", txErr)
	}

	resp := &CheckoutResponse{Success: true, OrderID: orderID, Total: total}
	if markErr := c.idempotency.MarkCompleted(ctx, scopedKey, resp); markErr != nil {
		c.logger.Error(ctx, "failed to finalize idempotency record; business transaction was committed", "userID", principal.UserID, "idempotencyKey", scopedKey, "error", markErr.Error())
		return nil, fmt.Errorf("%w: %v", ErrIdempotencyFinalize, markErr)
	}

	// Best-effort post-commit notification (Synchronous).
	// In production, consider the Outbox Pattern for reliable delivery.
	if notifyErr := c.notify.SendOrderConfirmation(ctx, principal.UserID, orderID); notifyErr != nil {
		c.logger.Error(ctx, "send order confirmation failed", "userID", principal.UserID, "orderID", orderID, "error", notifyErr.Error())
	}

	c.logger.Info(ctx, "checkout completed", "userID", principal.UserID, "orderID", orderID, "total", total)

	return resp, nil
}

func (c *CheckoutImproved) hashRequest(cmd CheckoutCommand) string {
	b, _ := json.Marshal(cmd)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
