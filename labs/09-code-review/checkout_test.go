package codereview_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	codereview "github.com/lukman-ss/software-engineering-lab/labs/09-code-review"
)

type testCartSource struct {
	mu    sync.Mutex
	carts map[string][]codereview.CartItem
}

func newTestCartSource() *testCartSource {
	return &testCartSource{carts: make(map[string][]codereview.CartItem)}
}

func (c *testCartSource) GetCart(ctx context.Context, userID string) ([]codereview.CartItem, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	items := c.carts[userID]
	if items == nil {
		return []codereview.CartItem{}, nil
	}
	result := make([]codereview.CartItem, len(items))
	copy(result, items)
	return result, nil
}

func (c *testCartSource) SetCart(userID string, items []codereview.CartItem) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.carts[userID] = items
}

// TestNaiveCheckout_RaceCondition demonstrates race condition bug
// This test intentionally shows the bug in the naive implementation
func TestNaiveCheckout_RaceCondition(t *testing.T) {
	repo := codereview.NewMockOrderRepository()
	products := codereview.NewMockProductRepository(map[string]int{"p1": 1})
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()

	cart.SetCart("user1", []codereview.CartItem{{ProductID: "p1", Quantity: 1, UnitPrice: 1000}})
	cart.SetCart("user2", []codereview.CartItem{{ProductID: "p1", Quantity: 1, UnitPrice: 1000}})

	naive := codereview.NewCheckoutNaive(repo, products, notify, logger, cart)

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err := naive.Checkout(context.Background(), "user1")
		errChan <- err
	}()
	go func() {
		defer wg.Done()
		_, err := naive.Checkout(context.Background(), "user2")
		errChan <- err
	}()

	wg.Wait()
	close(errChan)

	successCount := 0
	for err := range errChan {
		if err == nil {
			successCount++
		}
	}

	stock, _ := products.GetStock(context.Background(), "p1")

	t.Logf("Naive: successes=%d final_stock=%d", successCount, stock)

	if stock < 0 || successCount > 1 {
		t.Skip("Naive implementation has race condition bug (expected)")
	}
}

func TestImprovedCheckout_NoRaceCondition(t *testing.T) {
	repo := codereview.NewMockOrderRepository()
	products := codereview.NewMockProductRepository(map[string]int{"p1": 1})
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()
	idem := codereview.NewMockIdempotencyRepository()

	cart.SetCart("user1", []codereview.CartItem{{ProductID: "p1", Quantity: 1, UnitPrice: 1000}})
	cart.SetCart("user2", []codereview.CartItem{{ProductID: "p1", Quantity: 1, UnitPrice: 1000}})

	improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem)

	var wg sync.WaitGroup
	errChan := make(chan error, 2)
	successCount := atomic.Int32{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err := improved.Checkout(context.Background(), "user1", "key-user1")
		if err == nil {
			successCount.Add(1)
		}
		errChan <- err
	}()
	go func() {
		defer wg.Done()
		_, err := improved.Checkout(context.Background(), "user2", "key-user2")
		if err == nil {
			successCount.Add(1)
		}
		errChan <- err
	}()

	wg.Wait()
	close(errChan)

	stock, _ := products.GetStock(context.Background(), "p1")
	orders := len(repo.GetOrders(context.Background()))

	t.Logf("Improved: successes=%d stock=%d orders=%d", successCount.Load(), stock, orders)

	if stock < 0 {
		t.Errorf("INVALID: Stock cannot be negative: %d", stock)
	}

	if orders != 1 {
		t.Errorf("Expected exactly 1 order, got %d", orders)
	}
}

func TestImprovedCheckout_InsufficientStock(t *testing.T) {
	repo := codereview.NewMockOrderRepository()
	products := codereview.NewMockProductRepository(map[string]int{"p1": 5})
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()
	idem := codereview.NewMockIdempotencyRepository()

	cart.SetCart("user1", []codereview.CartItem{{ProductID: "p1", Quantity: 10, UnitPrice: 1000}})

	improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem)

	_, err := improved.Checkout(context.Background(), "user1", "test-idem-key")

	if err == nil {
		t.Error("Expected error for insufficient stock")
	}

	stock, _ := products.GetStock(context.Background(), "p1")
	if stock != 5 {
		t.Errorf("Stock should remain 5, got %d", stock)
	}

	orders := len(repo.GetOrders(context.Background()))
	if orders != 0 {
		t.Errorf("No orders should be created, got %d", orders)
	}
}

func TestImprovedCheckout_DuplicateRequest(t *testing.T) {
	repo := codereview.NewMockOrderRepository()
	products := codereview.NewMockProductRepository(map[string]int{"p1": 10})
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()
	idem := codereview.NewMockIdempotencyRepository()

	cart.SetCart("user1", []codereview.CartItem{{ProductID: "p1", Quantity: 2, UnitPrice: 1000}})

	key := "duplicate-test-key"
	improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem)

	resp1, err1 := improved.Checkout(context.Background(), "user1", key)
	if err1 != nil {
		t.Fatalf("First request failed: %v", err1)
	}

	resp2, err2 := improved.Checkout(context.Background(), "user1", key)
	if err2 != nil {
		t.Errorf("Duplicate request should succeed with cached response: %v", err2)
	}

	if resp1.OrderID != resp2.OrderID {
		t.Error("Duplicate request should return same order ID")
	}

	orders := len(repo.GetOrders(context.Background()))
	if orders != 1 {
		t.Errorf("Expected 1 order from 2 requests, got %d", orders)
	}
}

func TestImprovedCheckout_EmptyCart(t *testing.T) {
	repo := codereview.NewMockOrderRepository()
	products := codereview.NewMockProductRepository(map[string]int{"p1": 10})
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()
	idem := codereview.NewMockIdempotencyRepository()

	improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem)

	_, err := improved.Checkout(context.Background(), "user1", "test-key")

	if err == nil {
		t.Error("Expected error for empty cart")
	}
}

func TestImprovedCheckout_ProductNotFound(t *testing.T) {
	repo := codereview.NewMockOrderRepository()
	products := codereview.NewMockProductRepository(map[string]int{})
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()
	idem := codereview.NewMockIdempotencyRepository()

	cart.SetCart("user1", []codereview.CartItem{{ProductID: "nonexistent", Quantity: 1, UnitPrice: 1000}})

	improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem)

	_, err := improved.Checkout(context.Background(), "user1", "test-key")

	if err == nil {
		t.Error("Expected error for product not found")
	}
}

func TestImprovedCheckout_InvalidQuantity(t *testing.T) {
	repo := codereview.NewMockOrderRepository()
	products := codereview.NewMockProductRepository(map[string]int{"p1": 10})
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()
	idem := codereview.NewMockIdempotencyRepository()

	cart.SetCart("user1", []codereview.CartItem{{ProductID: "p1", Quantity: 0, UnitPrice: 1000}})

	improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem)

	_, err := improved.Checkout(context.Background(), "user1", "test-key")

	if err == nil {
		t.Error("Expected error for invalid quantity")
	}
}

func TestConcurrentCheckout_InvariantPreserved(t *testing.T) {
	for i := 0; i < 3; i++ {
		t.Run("iteration-"+string(rune('A'+i)), func(t *testing.T) {
			repo := codereview.NewMockOrderRepository()
			products := codereview.NewMockProductRepository(map[string]int{"p1": 10})
			notify := codereview.NewMockNotificationSender()
			logger := codereview.NewMockLogger()
			cart := newTestCartSource()
			idem := codereview.NewMockIdempotencyRepository()

			for u := 1; u <= 4; u++ {
				userID := string(rune('A' + u))
				cart.SetCart(userID, []codereview.CartItem{{ProductID: "p1", Quantity: 2, UnitPrice: 1000}})
			}

			improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem)

			var wg sync.WaitGroup
			for u := 1; u <= 4; u++ {
				wg.Add(1)
				go func(userID string) {
					defer wg.Done()
					_, _ = improved.Checkout(context.Background(), userID, "key-"+userID)
				}(string(rune('A' + u)))
			}
			wg.Wait()

			stock, _ := products.GetStock(context.Background(), "p1")
			if stock < 0 {
				t.Errorf("Stock invariant violated: %d", stock)
			}
		})
	}
}