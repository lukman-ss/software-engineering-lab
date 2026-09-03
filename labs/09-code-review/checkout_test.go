package codereview_test

import (
	"context"
	"errors"
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

func setupImproved(t *testing.T, initialStock map[string]int) (*codereview.CheckoutImproved, *codereview.MockProductRepository, *codereview.MockOrderRepository, *codereview.MockIdempotencyRepository, *codereview.MockNotificationSender, *codereview.MockTransactionManager, *testCartSource) {
	t.Helper()
	repo := codereview.NewMockOrderRepository()
	products := codereview.NewMockProductRepository(initialStock)
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()
	idem := codereview.NewMockIdempotencyRepository()
	txManager := codereview.NewMockTransactionManager(products, repo)
	improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem, txManager)
	return improved, products, repo, idem, notify, txManager, cart
}

// TestNaiveCheckout_DeterministicallyDemonstratesOverselling memaksa race condition pada
// implementasi naive dengan barrier. Bug overselling direproduksi secara deterministik.
func TestNaiveCheckout_DeterministicallyDemonstratesOverselling(t *testing.T) {
	repo := codereview.NewMockOrderRepository()
	products := codereview.NewMockProductRepository(map[string]int{"p1": 1})
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()

	cart.SetCart("user1", []codereview.CartItem{{ProductID: "p1", Quantity: 1, UnitPrice: 1000}})
	cart.SetCart("user2", []codereview.CartItem{{ProductID: "p1", Quantity: 1, UnitPrice: 1000}})

	var readBarrier, writeBarrier sync.WaitGroup
	readBarrier.Add(2)
	writeBarrier.Add(2)

	products.SetReadHook(func(productID string) {
		readBarrier.Done()
		readBarrier.Wait()
		writeBarrier.Done()
		writeBarrier.Wait()
	})

	naive := codereview.NewCheckoutNaive(repo, products, notify, logger, cart)

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(2)
	successCount := atomic.Int32{}
	go func() {
		defer wg.Done()
		_, err := naive.Checkout(context.Background(), "user1")
		if err == nil {
			successCount.Add(1)
		}
		errChan <- err
	}()
	go func() {
		defer wg.Done()
		_, err := naive.Checkout(context.Background(), "user2")
		if err == nil {
			successCount.Add(1)
		}
		errChan <- err
	}()

	wg.Wait()
	close(errChan)

	stock, _ := products.GetStock(context.Background(), "p1")

	t.Logf("Naive: successes=%d final_stock=%d", successCount.Load(), stock)

	// Naive sengaja broken: kedua request berhasil logis walhampos, menghasilkan oversell.
	if successCount.Load() != 2 {
		t.Errorf("Naive should demonstrate overselling with 2 logical successes, got %d", successCount.Load())
	}
	if stock >= 0 {
		t.Errorf("Naive should demonstrate negative stock (overselling), got %d", stock)
	}
}

// TestImprovedCheckout_DeterministicNoRace memastikan improved implementation aman
// terhadap race condition yang sama menggunakan barrier yang sama.
func TestImprovedCheckout_DeterministicNoRace(t *testing.T) {
	improved, products, repo, _, _, _, cart := setupImproved(t, map[string]int{"p1": 1})

	cart.SetCart("user1", []codereview.CartItem{{ProductID: "p1", Quantity: 1, UnitPrice: 1000}})
	cart.SetCart("user2", []codereview.CartItem{{ProductID: "p1", Quantity: 1, UnitPrice: 1000}})

	var readBarrier, writeBarrier sync.WaitGroup
	readBarrier.Add(2)
	writeBarrier.Add(2)

	products.SetReadHook(func(productID string) {
		readBarrier.Done()
		readBarrier.Wait()
		writeBarrier.Done()
		writeBarrier.Wait()
	})

	var wg sync.WaitGroup
	errChan := make(chan error, 2)
	successCount := atomic.Int32{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err := improved.Checkout(context.Background(), codereview.Principal{UserID: "user1"}, codereview.CheckoutCommand{
			CartOwnerID:    "user1",
			IdempotencyKey: "key-user1",
		})
		if err == nil {
			successCount.Add(1)
		}
		errChan <- err
	}()
	go func() {
		defer wg.Done()
		_, err := improved.Checkout(context.Background(), codereview.Principal{UserID: "user2"}, codereview.CheckoutCommand{
			CartOwnerID:    "user2",
			IdempotencyKey: "key-user2",
		})
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
		t.Errorf("Stock cannot be negative: %d", stock)
	}
	if orders != 1 {
		t.Errorf("Expected exactly 1 order, got %d", orders)
	}
	if successCount.Load() != 1 {
		t.Errorf("Expected exactly 1 success, got %d", successCount.Load())
	}
}

func TestImprovedCheckout_NoRaceCondition(t *testing.T) {
	repo := codereview.NewMockOrderRepository()
	products := codereview.NewMockProductRepository(map[string]int{"p1": 1})
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()
	idem := codereview.NewMockIdempotencyRepository()
	txManager := codereview.NewMockTransactionManager(products, repo)

	cart.SetCart("user1", []codereview.CartItem{{ProductID: "p1", Quantity: 1, UnitPrice: 1000}})
	cart.SetCart("user2", []codereview.CartItem{{ProductID: "p1", Quantity: 1, UnitPrice: 1000}})

	improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem, txManager)

	var wg sync.WaitGroup
	errChan := make(chan error, 2)
	successCount := atomic.Int32{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err := improved.Checkout(context.Background(), codereview.Principal{UserID: "user1"}, codereview.CheckoutCommand{
			CartOwnerID:    "user1",
			IdempotencyKey: "key-user1",
		})
		if err == nil {
			successCount.Add(1)
		}
		errChan <- err
	}()
	go func() {
		defer wg.Done()
		_, err := improved.Checkout(context.Background(), codereview.Principal{UserID: "user2"}, codereview.CheckoutCommand{
			CartOwnerID:    "user2",
			IdempotencyKey: "key-user2",
		})
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
	txManager := codereview.NewMockTransactionManager(products, repo)

	cart.SetCart("user1", []codereview.CartItem{{ProductID: "p1", Quantity: 10, UnitPrice: 1000}})

	improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem, txManager)

	_, err := improved.Checkout(context.Background(), codereview.Principal{UserID: "user1"}, codereview.CheckoutCommand{
		CartOwnerID:    "user1",
		IdempotencyKey: "test-idem-key",
	})

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

func TestImprovedCheckout_BatchLoadsProducts(t *testing.T) {
	repo := codereview.NewMockOrderRepository()
	products := codereview.NewMockProductRepository(map[string]int{"p1": 10, "p2": 10})
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()
	idem := codereview.NewMockIdempotencyRepository()
	txManager := codereview.NewMockTransactionManager(products, repo)

	cart.SetCart("user1", []codereview.CartItem{
		{ProductID: "p1", Quantity: 2, UnitPrice: 1000},
		{ProductID: "p2", Quantity: 1, UnitPrice: 2000},
	})

	improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem, txManager)

	_, err := improved.Checkout(context.Background(), codereview.Principal{UserID: "user1"}, codereview.CheckoutCommand{
		CartOwnerID:    "user1",
		IdempotencyKey: "batch-test-key",
	})
	if err != nil {
		t.Fatalf("checkout failed: %v", err)
	}

	if products.BatchGetCalls() != 1 {
		t.Errorf("Expected exactly 1 batch GetProducts call, got %d", products.BatchGetCalls())
	}
	if products.GetCalls() != 0 {
		t.Errorf("Expected no single GetProduct calls on improved path, got %d", products.GetCalls())
	}
}

func TestImprovedCheckout_DuplicateRequestReturnsSameResult(t *testing.T) {
	repo := codereview.NewMockOrderRepository()
	products := codereview.NewMockProductRepository(map[string]int{"p1": 10})
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()
	idem := codereview.NewMockIdempotencyRepository()
	txManager := codereview.NewMockTransactionManager(products, repo)

	cart.SetCart("user1", []codereview.CartItem{{ProductID: "p1", Quantity: 2, UnitPrice: 1000}})

	key := "duplicate-test-key"
	improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem, txManager)

	principal := codereview.Principal{UserID: "user1"}
	cmd := codereview.CheckoutCommand{CartOwnerID: "user1", IdempotencyKey: key}

	resp1, err1 := improved.Checkout(context.Background(), principal, cmd)
	if err1 != nil {
		t.Fatalf("First request failed: %v", err1)
	}

	resp2, err2 := improved.Checkout(context.Background(), principal, cmd)
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

func TestImprovedCheckout_IdempotencyFailureCanRetry(t *testing.T) {
	repo := codereview.NewMockOrderRepository()
	repo.FailNextCreate(errors.New("database unavailable"))
	products := codereview.NewMockProductRepository(map[string]int{"p1": 10})
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()
	idem := codereview.NewMockIdempotencyRepository()
	txManager := codereview.NewMockTransactionManager(products, repo)

	cart.SetCart("user1", []codereview.CartItem{{ProductID: "p1", Quantity: 2, UnitPrice: 1000}})
	improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem, txManager)

	key := "retryable-key"
	cmd := codereview.CheckoutCommand{CartOwnerID: "user1", IdempotencyKey: key}

	resp, err := improved.Checkout(context.Background(), codereview.Principal{UserID: "user1"}, cmd)
	if err == nil {
		t.Fatal("Expected first request to fail")
	}
	if resp != nil {
		t.Fatal("Expected response to be nil on failure")
	}

	stock, _ := products.GetStock(context.Background(), "p1")
	if stock != 10 {
		t.Errorf("Stock should remain 10 on failure, got %d", stock)
	}
	if len(repo.GetOrders(context.Background())) != 0 {
		t.Error("No order should be created")
	}
	if len(notify.GetSent()) != 0 {
		t.Error("No notification should be sent")
	}

	repo.FailNextCreate(nil)
	retryResp, err := improved.Checkout(context.Background(), codereview.Principal{UserID: "user1"}, cmd)
	if err != nil {
		t.Fatalf("Retry request should succeed, got error: %v", err)
	}
	if retryResp == nil || !retryResp.Success {
		t.Fatal("Expected successful retry response")
	}
	if len(repo.GetOrders(context.Background())) != 1 {
		t.Errorf("Expected 1 order created after retry, got %d", len(repo.GetOrders(context.Background())))
	}
	finalStock, _ := products.GetStock(context.Background(), "p1")
	if finalStock != 8 {
		t.Errorf("Stock should be deducted exactly once (8), got %d", finalStock)
	}
}

func TestImprovedCheckout_MultiProductFailureRollsBackEverything(t *testing.T) {
	repo := codereview.NewMockOrderRepository()
	repo.FailNextCreate(errors.New("database unavailable"))
	products := codereview.NewMockProductRepository(map[string]int{"p1": 10, "p2": 1})
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()
	idem := codereview.NewMockIdempotencyRepository()
	txManager := codereview.NewMockTransactionManager(products, repo)

	cart.SetCart("user1", []codereview.CartItem{
		{ProductID: "p1", Quantity: 2, UnitPrice: 1000},
		{ProductID: "p2", Quantity: 1, UnitPrice: 2000},
	})

	improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem, txManager)

	_, err := improved.Checkout(context.Background(), codereview.Principal{UserID: "user1"}, codereview.CheckoutCommand{
		CartOwnerID:    "user1",
		IdempotencyKey: "fail-multi-key",
	})
	if err == nil {
		t.Fatal("Expected checkout to fail")
	}

	p1Stock, _ := products.GetStock(context.Background(), "p1")
	if p1Stock != 10 {
		t.Errorf("p1 stock should be rolled back to 10, got %d", p1Stock)
	}

	p2Stock, _ := products.GetStock(context.Background(), "p2")
	if p2Stock != 1 {
		t.Errorf("p2 stock should be rolled back to 1, got %d", p2Stock)
	}

	if len(repo.GetOrders(context.Background())) != 0 {
		t.Error("No order should be created")
	}
}

func TestImprovedCheckout_ConcurrentStockReservationPreservesInvariant(t *testing.T) {
	repo := codereview.NewMockOrderRepository()
	products := codereview.NewMockProductRepository(map[string]int{"p1": 1})
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()
	idem := codereview.NewMockIdempotencyRepository()
	txManager := codereview.NewMockTransactionManager(products, repo)

	cart.SetCart("user1", []codereview.CartItem{{ProductID: "p1", Quantity: 1, UnitPrice: 1000}})
	cart.SetCart("user2", []codereview.CartItem{{ProductID: "p1", Quantity: 1, UnitPrice: 1000}})

	improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem, txManager)

	var wg sync.WaitGroup
	errChan := make(chan error, 2)
	successCount := atomic.Int32{}
	failCount := atomic.Int32{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err := improved.Checkout(context.Background(), codereview.Principal{UserID: "user1"}, codereview.CheckoutCommand{
			CartOwnerID:    "user1",
			IdempotencyKey: "key-user1",
		})
		if err == nil {
			successCount.Add(1)
		} else {
			failCount.Add(1)
		}
		errChan <- err
	}()
	go func() {
		defer wg.Done()
		_, err := improved.Checkout(context.Background(), codereview.Principal{UserID: "user2"}, codereview.CheckoutCommand{
			CartOwnerID:    "user2",
			IdempotencyKey: "key-user2",
		})
		if err == nil {
			successCount.Add(1)
		} else {
			failCount.Add(1)
		}
		errChan <- err
	}()

	wg.Wait()
	close(errChan)

	stock, _ := products.GetStock(context.Background(), "p1")
	orders := len(repo.GetOrders(context.Background()))

	if successCount.Load() != 1 {
		t.Errorf("Expected 1 successful checkout, got %d", successCount.Load())
	}
	if failCount.Load() != 1 {
		t.Errorf("Expected 1 failed checkout, got %d", failCount.Load())
	}
	if stock != 0 {
		t.Errorf("Final stock must be 0, got %d", stock)
	}
	if orders != 1 {
		t.Errorf("Expected exactly 1 order, got %d", orders)
	}
	if len(notify.GetSent()) > 1 {
		t.Errorf("Expected at most 1 notification sent, got %d", len(notify.GetSent()))
	}
}

func TestImprovedCheckout_ConcurrentSameIdempotencyKeyRunsOnce(t *testing.T) {
	repo := codereview.NewMockOrderRepository()
	products := codereview.NewMockProductRepository(map[string]int{"p1": 10})
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()
	idem := codereview.NewMockIdempotencyRepository()
	txManager := codereview.NewMockTransactionManager(products, repo)

	cart.SetCart("user1", []codereview.CartItem{{ProductID: "p1", Quantity: 2, UnitPrice: 1000}})
	improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem, txManager)

	var startGate sync.WaitGroup
	startGate.Add(1)

	var wg sync.WaitGroup
	errChan := make(chan error, 2)
	successCount := atomic.Int32{}

	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()
			startGate.Wait()
			_, err := improved.Checkout(context.Background(), codereview.Principal{UserID: "user1"}, codereview.CheckoutCommand{
				CartOwnerID:    "user1",
				IdempotencyKey: "concurrent-same-key",
			})
			if err == nil {
				successCount.Add(1)
			}
			errChan <- err
		}()
	}

	startGate.Done()
	wg.Wait()
	close(errChan)

	orders := len(repo.GetOrders(context.Background()))
	if orders != 1 {
		t.Errorf("Expected exactly 1 order, got %d", orders)
	}

	stock, _ := products.GetStock(context.Background(), "p1")
	if stock != 8 {
		t.Errorf("Stock should be deducted only once (8), got %d", stock)
	}

	if len(notify.GetSent()) != 1 {
		t.Errorf("Notification should be sent at most once, got %d", len(notify.GetSent()))
	}
}

func TestImprovedCheckout_IdempotencyClaimFailure(t *testing.T) {
	repo := codereview.NewMockOrderRepository()
	products := codereview.NewMockProductRepository(map[string]int{"p1": 10})
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()
	idem := codereview.NewMockIdempotencyRepository()
	txManager := codereview.NewMockTransactionManager(products, repo)

	claimErr := errors.New("redis timeout during claim")
	idem.FailNextClaim(claimErr)

	cart.SetCart("user1", []codereview.CartItem{{ProductID: "p1", Quantity: 2, UnitPrice: 1000}})
	improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem, txManager)

	_, err := improved.Checkout(context.Background(), codereview.Principal{UserID: "user1"}, codereview.CheckoutCommand{
		CartOwnerID:    "user1",
		IdempotencyKey: "claim-fail-key",
	})

	if err == nil {
		t.Fatal("Expected error when claim fails")
	}
	if !errors.Is(err, claimErr) {
		t.Errorf("Expected claim error %v, got %v", claimErr, err)
	}

	stock, _ := products.GetStock(context.Background(), "p1")
	if stock != 10 {
		t.Errorf("Stock should remain 10, got %d", stock)
	}
	if len(repo.GetOrders(context.Background())) != 0 {
		t.Error("Order should not be created")
	}
}

func TestImprovedCheckout_IdempotencyFinalizeFailureAfterCommit(t *testing.T) {
	repo := codereview.NewMockOrderRepository()
	products := codereview.NewMockProductRepository(map[string]int{"p1": 10})
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()
	idem := codereview.NewMockIdempotencyRepository()
	txManager := codereview.NewMockTransactionManager(products, repo)

	idem.FailNextMarkCompleted(errors.New("redis unavailable"))
	cart.SetCart("user1", []codereview.CartItem{{ProductID: "p1", Quantity: 2, UnitPrice: 1000}})
	improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem, txManager)

	_, err := improved.Checkout(context.Background(), codereview.Principal{UserID: "user1"}, codereview.CheckoutCommand{
		CartOwnerID:    "user1",
		IdempotencyKey: "finalize-fail-key",
	})

	if err == nil {
		t.Fatal("Expected error due to idempotency finalize failure")
	}
	if !errors.Is(err, codereview.ErrIdempotencyFinalize) {
		t.Errorf("Expected ErrIdempotencyFinalize, got: %v", err)
	}

	stock, _ := products.GetStock(context.Background(), "p1")
	if stock != 8 {
		t.Errorf("Stock should be deducted because business tx committed, got %d", stock)
	}
	if len(repo.GetOrders(context.Background())) != 1 {
		t.Errorf("Order should be created because business tx committed, got %d", len(repo.GetOrders(context.Background())))
	}
	if len(notify.GetSent()) != 0 {
		t.Error("Notification should not be sent on finalize failure")
	}

	record, _ := idem.Get(context.Background(), "checkout:user1:finalize-fail-key")
	if record != nil && record.Status == codereview.IdempotencyStatusCompleted {
		t.Error("Idempotency record should not be marked as COMPLETED")
	}
}

func TestImprovedCheckout_ReleaseFailureIsNotIgnored(t *testing.T) {
	repo := codereview.NewMockOrderRepository()
	products := codereview.NewMockProductRepository(map[string]int{"p1": 10})
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()
	idem := codereview.NewMockIdempotencyRepository()
	txManager := codereview.NewMockTransactionManager(products, repo)

	// Make transaction fail to trigger release
	repo.FailNextCreate(errors.New("db error"))
	// Make release fail
	releaseErr := errors.New("redis error during release")
	idem.FailNextRelease(releaseErr)

	cart.SetCart("user1", []codereview.CartItem{{ProductID: "p1", Quantity: 2, UnitPrice: 1000}})
	improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem, txManager)

	_, err := improved.Checkout(context.Background(), codereview.Principal{UserID: "user1"}, codereview.CheckoutCommand{
		CartOwnerID:    "user1",
		IdempotencyKey: "release-fail-key",
	})

	if err == nil {
		t.Fatal("Expected error")
	}
	// Error should contain both transaction error and release error
	if !errors.Is(err, releaseErr) {
		t.Errorf("Error should wrap release error, got: %v", err)
	}

	stock, _ := products.GetStock(context.Background(), "p1")
	if stock != 10 {
		t.Errorf("Stock should remain 10, got %d", stock)
	}
}

func TestImprovedCheckout_EmptyIdempotencyKeyFails(t *testing.T) {
	repo := codereview.NewMockOrderRepository()
	products := codereview.NewMockProductRepository(map[string]int{"p1": 10})
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()
	idem := codereview.NewMockIdempotencyRepository()
	txManager := codereview.NewMockTransactionManager(products, repo)

	improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem, txManager)

	_, err := improved.Checkout(context.Background(), codereview.Principal{UserID: "user1"}, codereview.CheckoutCommand{
		CartOwnerID:    "user1",
		IdempotencyKey: "",
	})

	if err == nil {
		t.Fatal("Expected error for empty idempotency key")
	}
	if !errors.Is(err, codereview.ErrInvalidIdempotencyKey) {
		t.Errorf("Expected ErrInvalidIdempotencyKey, got: %v", err)
	}
}

func TestImprovedCheckout_RetryAfterCartMutationReturnsOriginalResponse(t *testing.T) {
	repo := codereview.NewMockOrderRepository()
	products := codereview.NewMockProductRepository(map[string]int{"p1": 10})
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()
	idem := codereview.NewMockIdempotencyRepository()
	txManager := codereview.NewMockTransactionManager(products, repo)

	cart.SetCart("user1", []codereview.CartItem{{ProductID: "p1", Quantity: 2, UnitPrice: 1000}})
	improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem, txManager)

	key := "retry-mutation-key"
	cmd := codereview.CheckoutCommand{CartOwnerID: "user1", IdempotencyKey: key}

	// 1. Initial success
	resp1, err := improved.Checkout(context.Background(), codereview.Principal{UserID: "user1"}, cmd)
	if err != nil {
		t.Fatalf("First request failed: %v", err)
	}

	// 2. Client empties/changes the cart
	cart.SetCart("user1", []codereview.CartItem{})

	// 3. Retry with same key
	resp2, err := improved.Checkout(context.Background(), codereview.Principal{UserID: "user1"}, cmd)
	if err != nil {
		t.Fatalf("Retry failed: %v", err)
	}

	// Should return identical response to resp1 (same order ID, success state)
	if resp1.OrderID != resp2.OrderID {
		t.Errorf("Retry should return same OrderID. Expected %s, got %s", resp1.OrderID, resp2.OrderID)
	}

	stock, _ := products.GetStock(context.Background(), "p1")
	if stock != 8 {
		t.Errorf("Stock should only be deducted once. Expected 8, got %d", stock)
	}
}

func TestImprovedCheckout_CannotCheckoutAnotherUsersCart(t *testing.T) {
	repo := codereview.NewMockOrderRepository()
	products := codereview.NewMockProductRepository(map[string]int{"p1": 10})
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()
	idem := codereview.NewMockIdempotencyRepository()
	txManager := codereview.NewMockTransactionManager(products, repo)

	// Cart milik user-B
	cart.SetCart("userB", []codereview.CartItem{{ProductID: "p1", Quantity: 2, UnitPrice: 1000}})

	improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem, txManager)

	// Authenticated principal = user-A, tapi cart milik user-B
	_, err := improved.Checkout(context.Background(), codereview.Principal{UserID: "userA"}, codereview.CheckoutCommand{
		CartOwnerID:    "userB",
		IdempotencyKey: "authz-test-key",
	})

	if !errors.Is(err, codereview.ErrForbidden) {
		t.Errorf("Expected ErrForbidden, got: %v", err)
	}

	stock, _ := products.GetStock(context.Background(), "p1")
	if stock != 10 {
		t.Errorf("Stock should remain unchanged at 10, got %d", stock)
	}

	orders := len(repo.GetOrders(context.Background()))
	if orders != 0 {
		t.Errorf("No orders should be created, got %d", orders)
	}
}

func TestImprovedCheckout_IdempotencyKeyIsScopedPerUser(t *testing.T) {
	repo := codereview.NewMockOrderRepository()
	products := codereview.NewMockProductRepository(map[string]int{"p1": 10})
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()
	idem := codereview.NewMockIdempotencyRepository()
	txManager := codereview.NewMockTransactionManager(products, repo)

	// Same cart content for two different users
	cart.SetCart("userA", []codereview.CartItem{{ProductID: "p1", Quantity: 2, UnitPrice: 1000}})
	cart.SetCart("userB", []codereview.CartItem{{ProductID: "p1", Quantity: 2, UnitPrice: 1000}})

	improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem, txManager)

	sameKey := "shared-key"
	cmdA := codereview.CheckoutCommand{CartOwnerID: "userA", IdempotencyKey: sameKey}
	cmdB := codereview.CheckoutCommand{CartOwnerID: "userB", IdempotencyKey: sameKey}

	_, err := improved.Checkout(context.Background(), codereview.Principal{UserID: "userA"}, cmdA)
	if err != nil {
		t.Fatalf("userA checkout failed: %v", err)
	}

	// userB with same idempotency key string harus diperlakukan sebagai request yang berbeda
	_, err = improved.Checkout(context.Background(), codereview.Principal{UserID: "userB"}, cmdB)
	if err != nil {
		t.Errorf("userB with same idempotency key should be independent request: %v", err)
	}

	orders := len(repo.GetOrders(context.Background()))
	if orders != 2 {
		t.Errorf("Expected 2 orders (one per user), got %d", orders)
	}

	stock, _ := products.GetStock(context.Background(), "p1")
	if stock != 6 {
		t.Errorf("Expected stock 6 (10 - 4), got %d", stock)
	}
}

func TestNaiveCheckout_NPlusOneProductLookups(t *testing.T) {
	repo := codereview.NewMockOrderRepository()
	products := codereview.NewMockProductRepository(map[string]int{"p1": 10, "p2": 10, "p3": 10})
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()

	cart.SetCart("user1", []codereview.CartItem{
		{ProductID: "p1", Quantity: 1, UnitPrice: 1000},
		{ProductID: "p2", Quantity: 1, UnitPrice: 1000},
		{ProductID: "p3", Quantity: 1, UnitPrice: 1000},
	})

	naive := codereview.NewCheckoutNaive(repo, products, notify, logger, cart)
	resp, err := naive.Checkout(context.Background(), "user1")
	if err != nil || resp == nil {
		t.Fatalf("Naive checkout failed: %v", err)
	}

	if products.GetCalls() != 3 {
		t.Errorf("Naive expected 3 N+1 GetProduct calls, got %d", products.GetCalls())
	}
	if products.BatchGetCalls() != 0 {
		t.Errorf("Naive should not use batch GetProducts")
	}
}

func TestImprovedCheckout_BatchLoadUsesUniqueProductIDs(t *testing.T) {
	repo := codereview.NewMockOrderRepository()
	products := codereview.NewMockProductRepository(map[string]int{"p1": 10, "p2": 10})
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()
	idem := codereview.NewMockIdempotencyRepository()
	txManager := codereview.NewMockTransactionManager(products, repo)

	cart.SetCart("user1", []codereview.CartItem{
		{ProductID: "p1", Quantity: 1, UnitPrice: 1000},
		{ProductID: "p1", Quantity: 2, UnitPrice: 1000},
		{ProductID: "p2", Quantity: 1, UnitPrice: 2000},
	})

	improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem, txManager)
	_, err := improved.Checkout(context.Background(), codereview.Principal{UserID: "user1"}, codereview.CheckoutCommand{
		CartOwnerID:    "user1",
		IdempotencyKey: "unique-batch-key",
	})
	if err != nil {
		t.Fatalf("Checkout failed: %v", err)
	}

	lastIDs := products.GetLastBatchIDs()
	if len(lastIDs) != 2 {
		t.Errorf("Expected exactly 2 unique product IDs requested in batch, got %d: %v", len(lastIDs), lastIDs)
	}
}

func TestImprovedCheckout_DuplicateProductInCartAccumulatesReservation(t *testing.T) {
	t.Run("success when total accumulated stock is sufficient", func(t *testing.T) {
		repo := codereview.NewMockOrderRepository()
		products := codereview.NewMockProductRepository(map[string]int{"p1": 5})
		notify := codereview.NewMockNotificationSender()
		logger := codereview.NewMockLogger()
		cart := newTestCartSource()
		idem := codereview.NewMockIdempotencyRepository()
		txManager := codereview.NewMockTransactionManager(products, repo)

		cart.SetCart("user1", []codereview.CartItem{
			{ProductID: "p1", Quantity: 2, UnitPrice: 1000},
			{ProductID: "p1", Quantity: 3, UnitPrice: 1000},
		})

		improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem, txManager)
		_, err := improved.Checkout(context.Background(), codereview.Principal{UserID: "user1"}, codereview.CheckoutCommand{
			CartOwnerID:    "user1",
			IdempotencyKey: "dup-cart-success",
		})
		if err != nil {
			t.Fatalf("Checkout should succeed: %v", err)
		}

		stock, _ := products.GetStock(context.Background(), "p1")
		if stock != 0 {
			t.Errorf("Final stock should be 0 (5 - 2 - 3), got %d", stock)
		}
		if len(repo.GetOrders(context.Background())) != 1 {
			t.Errorf("Expected 1 order created, got %d", len(repo.GetOrders(context.Background())))
		}
	})

	t.Run("fails and rolls back when accumulated reservation exceeds stock", func(t *testing.T) {
		repo := codereview.NewMockOrderRepository()
		products := codereview.NewMockProductRepository(map[string]int{"p1": 4})
		notify := codereview.NewMockNotificationSender()
		logger := codereview.NewMockLogger()
		cart := newTestCartSource()
		idem := codereview.NewMockIdempotencyRepository()
		txManager := codereview.NewMockTransactionManager(products, repo)

		cart.SetCart("user1", []codereview.CartItem{
			{ProductID: "p1", Quantity: 2, UnitPrice: 1000},
			{ProductID: "p1", Quantity: 3, UnitPrice: 1000},
		})

		improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem, txManager)
		_, err := improved.Checkout(context.Background(), codereview.Principal{UserID: "user1"}, codereview.CheckoutCommand{
			CartOwnerID:    "user1",
			IdempotencyKey: "dup-cart-fail",
		})
		if err == nil {
			t.Fatal("Expected checkout to fail due to insufficient accumulated stock")
		}

		stock, _ := products.GetStock(context.Background(), "p1")
		if stock != 4 {
			t.Errorf("Final stock should remain 4, got %d", stock)
		}
		if len(repo.GetOrders(context.Background())) != 0 {
			t.Errorf("Expected 0 orders created, got %d", len(repo.GetOrders(context.Background())))
		}
	})
}

func TestIdempotencyRepository_StateMachineDirect(t *testing.T) {
	idem := codereview.NewMockIdempotencyRepository()
	ctx := context.Background()

	key := "test-state-key"
	hash1 := "hash-1"
	hash2 := "hash-2"

	// 1. First claim: transitions to PROCESSING
	status, resp, err := idem.Claim(ctx, key, hash1)
	if err != nil {
		t.Fatalf("First claim failed: %v", err)
	}
	if status != codereview.IdempotencyStatusProcessing || resp != nil {
		t.Errorf("Expected PROCESSING and nil resp, got status=%s resp=%v", status, resp)
	}

	// 2. Second claim with same hash while PROCESSING: ErrDuplicateRequest
	_, _, err = idem.Claim(ctx, key, hash1)
	if !errors.Is(err, codereview.ErrDuplicateRequest) {
		t.Errorf("Expected ErrDuplicateRequest while processing, got %v", err)
	}

	// 3. Second claim with different hash: ErrIdempotencyConflict
	_, _, err = idem.Claim(ctx, key, hash2)
	if !errors.Is(err, codereview.ErrIdempotencyConflict) {
		t.Errorf("Expected ErrIdempotencyConflict for payload mismatch, got %v", err)
	}

	// 4. Mark completed
	mockResp := &codereview.CheckoutResponse{Success: true, OrderID: "order-1", Total: 5000}
	if err := idem.MarkCompleted(ctx, key, mockResp); err != nil {
		t.Fatalf("MarkCompleted failed: %v", err)
	}

	// 5. Claim again after COMPLETED: returns cached response
	status, cachedResp, err := idem.Claim(ctx, key, hash1)
	if err != nil {
		t.Fatalf("Claim after completion failed: %v", err)
	}
	if status != codereview.IdempotencyStatusCompleted {
		t.Errorf("Expected COMPLETED status, got %s", status)
	}
	if cachedResp == nil || cachedResp.OrderID != "order-1" {
		t.Errorf("Expected cached response with orderID 'order-1', got %v", cachedResp)
	}
}

func TestImprovedCheckout_PreservesRepositoryErrors(t *testing.T) {
	repo := codereview.NewMockOrderRepository()
	products := codereview.NewMockProductRepository(map[string]int{"p1": 10})
	notify := codereview.NewMockNotificationSender()
	logger := codereview.NewMockLogger()
	cart := newTestCartSource()
	idem := codereview.NewMockIdempotencyRepository()
	txManager := codereview.NewMockTransactionManager(products, repo)

	infraErr := errors.New("connection refused")
	products.SetBatchGetError(infraErr)

	cart.SetCart("user1", []codereview.CartItem{{ProductID: "p1", Quantity: 1, UnitPrice: 1000}})
	improved := codereview.NewCheckoutImproved(repo, products, notify, logger, cart, idem, txManager)

	_, err := improved.Checkout(context.Background(), codereview.Principal{UserID: "user1"}, codereview.CheckoutCommand{
		CartOwnerID:    "user1",
		IdempotencyKey: "error-test-key",
	})

	if err == nil {
		t.Fatal("Expected error")
	}
	// Infra error harus tidak sama dengan ErrProductNotFound secara eksis
	if errors.Is(err, codereview.ErrProductNotFound) {
		t.Error("Infrastructure error should not be reported as ErrProductNotFound")
	}
}
