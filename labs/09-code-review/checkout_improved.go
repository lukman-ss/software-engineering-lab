package codereview

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

// CheckoutImproved demonstrates a safer checkout implementation
// after addressing the review findings in this lab.
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

	cartItems, err := c.cartSource.GetCart(ctx, cmd.CartOwnerID)
	if err != nil {
		c.logger.Error(ctx, "load cart failed", "userID", principal.UserID, "error", err.Error())
		return nil, fmt.Errorf("load cart: %w", err)
	}

	if len(cartItems) == 0 {
		return nil, ErrEmptyCart
	}

	productIDs := make([]string, 0, len(cartItems))
	for _, item := range cartItems {
		if item.Quantity <= 0 {
			return nil, ErrInvalidQuantity
		}
		productIDs = append(productIDs, item.ProductID)
	}

	hash := c.hashRequest(principal.UserID, cartItems)
	scopedKey := fmt.Sprintf("checkout:%s:%s", principal.UserID, cmd.IdempotencyKey)

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
			if err := tx.ReserveStock(ctx, item.ProductID, item.Quantity); err != nil {
				return err
			}
		}
		if err := tx.CreateOrder(ctx, order); err != nil {
			return err
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
	if err := c.idempotency.MarkCompleted(ctx, scopedKey, resp); err != nil {
		c.logger.Error(ctx, "failed to finalize idempotency record; business transaction was committed", "userID", principal.UserID, "idempotencyKey", scopedKey, "error", err.Error())
		return nil, fmt.Errorf("%w: %v", ErrIdempotencyFinalize, err)
	}

	// Notification after successful transaction (fire-and-forget)
	if err := c.notify.SendOrderConfirmation(ctx, principal.UserID, orderID); err != nil {
		c.logger.Error(ctx, "send order confirmation failed", "userID", principal.UserID, "orderID", orderID, "error", err.Error())
	}

	c.logger.Info(ctx, "checkout completed", "userID", principal.UserID, "orderID", orderID, "total", total)

	return resp, nil
}

func (c *CheckoutImproved) hashRequest(userID string, items []CartItem) string {
	type payload struct {
		UserID string
		Items  []CartItem
	}
	b, _ := json.Marshal(payload{UserID: userID, Items: items})
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
