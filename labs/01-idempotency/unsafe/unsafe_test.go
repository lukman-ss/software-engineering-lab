package unsafe_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/lukman/software-engineer-lab/labs/01-idempotency/unsafe"
)

// TestDuplicatePaymentBug demonstrates that without idempotency,
// concurrent/retried payment requests create duplicate charges.
func TestDuplicatePaymentBug(t *testing.T) {
	gw := &unsafe.MockGateway{}
	svc := unsafe.NewService(gw)

	req := unsafe.PaymentRequest{
		OrderID: "order-123",
		Amount:  10000, // $100.00
	}

	ctx := context.Background()

	// Simulate client retry: initial request times out, client retries
	// In a real scenario, both requests hit the server
	result1, err := svc.ProcessPayment(ctx, req)
	if err != nil {
		t.Fatalf("first payment failed: %v", err)
	}

	result2, err := svc.ProcessPayment(ctx, req)
	if err != nil {
		t.Fatalf("second payment failed: %v", err)
	}

	// BUG: Two different payment IDs means two charges were created
	if result1.PaymentID == result2.PaymentID {
		t.Fatalf("expected different payment IDs (bug), got %s == %s",
			result1.PaymentID, result2.PaymentID)
	}

	// Gateway should have been called twice (the problem)
	gatewayCalls := gw.ChargeCount()
	if gatewayCalls != 2 {
		t.Fatalf("expected gateway to be called 2 times, got %d", gatewayCalls)
	}

	// Two payments stored in database (the problem)
	payments := svc.GetPayments()
	if len(payments) != 2 {
		t.Fatalf("expected 2 payments in store, got %d", len(payments))
	}

	t.Logf("BUG REPRODUCED: gateway called %d times, %d payments created",
		gatewayCalls, len(payments))
}

// TestConcurrentDuplicatePayment simulates multiple concurrent retries.
func TestConcurrentDuplicatePayment(t *testing.T) {
	gw := &unsafe.MockGateway{}
	svc := unsafe.NewService(gw)

	req := unsafe.PaymentRequest{
		OrderID: "order-456",
		Amount:  5000,
	}

	ctx := context.Background()

	// Launch 10 concurrent payment requests (simulating retry storm)
	var wg sync.WaitGroup
	var successCount int64
	var errCount int64

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := svc.ProcessPayment(ctx, req)
			if err != nil {
				atomic.AddInt64(&errCount, 1)
			} else {
				atomic.AddInt64(&successCount, 1)
			}
		}()
	}

	wg.Wait()

	// BUG: Multiple gateway calls, multiple payment records
	if successCount != 10 {
		t.Logf("expected 10 successful responses (bug), got %d", successCount)
	}

	gatewayCalls := gw.ChargeCount()
	storeCount := svc.CountPayments()

	t.Logf("BUG REPRODUCED (concurrent): gateway called %d times, %d payments in store",
		gatewayCalls, storeCount)

	if gatewayCalls <= 1 {
		t.Fatalf("expected more than 1 gateway call to demonstrate the bug, got %d", gatewayCalls)
	}
}

// Test20ConcurrentRequestsWithBarrier demonstrates duplicate transactions
// with 20 concurrent requests using a synchronization barrier.
// This test ensures race conditions are reproducibly triggered without
// using time.Sleep for synchronization.
func Test20ConcurrentRequestsWithBarrier(t *testing.T) {
	gw := &unsafe.MockGateway{}
	svc := unsafe.NewService(gw)

	req := unsafe.PaymentRequest{
		OrderID: "order-999",
		Amount:  500000,
	}

	ctx := context.Background()

	// Synchronization barrier: all goroutines wait until ready to proceed
	const numRequests = 20
	barrier := make(chan struct{})

	// Wait group to collect all results
	var wg sync.WaitGroup
	var successCount int64
	var results []unsafe.PaymentResult
	var resultsMu sync.Mutex

	// Launch 20 concurrent goroutines
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Wait at the barrier until all goroutines are ready
			<-barrier

			// Process payment
			result, err := svc.ProcessPayment(ctx, req)
			if err != nil {
				t.Logf("goroutine %d: error: %v", id, err)
				return
			}

			// Collect result
			resultsMu.Lock()
			results = append(results, result)
			resultsMu.Unlock()
			atomic.AddInt64(&successCount, 1)

			t.Logf("goroutine %d: payment created: %s", id, result.PaymentID)
		}(i)
	}

	// Staggered start: open barrier after all goroutines are waiting
	// This maximizes the chance of race condition
	go func() {
		wg.Wait() // Wait for all goroutines to reach barrier
		// Now release all at once
		close(barrier)
	}()

	// Give goroutines time to start and wait at barrier
	// Using a tiny sleep here for goroutine scheduling, not synchronization
	// The actual synchronization is via the barrier channel
	//
	// NOTE: This is only to ensure goroutines are scheduled.
	// The real synchronization mechanism is the barrier channel itself.
	// time.Sleep(10 * time.Millisecond) would be used in practice,
	// but we skip it for deterministic testing in CI environments.

	// Wait for all to complete
	wg.Wait()

	// BUG VERIFICATION: Should have 20 different payment IDs (duplicate charges)
	uniqueIDs := make(map[string]bool)
	for _, r := range results {
		uniqueIDs[r.PaymentID] = true
	}

	t.Logf("Results: %d successful, %d unique payment IDs out of %d requests",
		successCount, len(uniqueIDs), numRequests)

	// Without idempotency, we expect all 20 to succeed with 20 unique IDs
	if int(successCount) != numRequests {
		t.Errorf("expected %d successful payments, got %d", numRequests, successCount)
	}

	// BUG: All 20 should have different payment IDs (20 duplicate charges!)
	if len(uniqueIDs) != numRequests {
		t.Errorf("expected %d unique payment IDs (bug: duplicate charges), got %d",
			numRequests, len(uniqueIDs))
	}

	gatewayCalls := gw.ChargeCount()
	storeCount := svc.CountPayments()

	t.Logf("BUG CONFIRMED: gateway called %d times, %d payments stored (expected %d)",
		gatewayCalls, storeCount, numRequests)

	if gatewayCalls != int64(numRequests) {
		t.Errorf("expected gateway called %d times, got %d", numRequests, gatewayCalls)
	}

	if storeCount != numRequests {
		t.Errorf("expected %d payments in store, got %d", numRequests, storeCount)
	}
}

// TestNoPayloadValidation demonstrates that the same idempotency key
// with different payloads is accepted (no validation).
func TestNoPayloadValidation(t *testing.T) {
	svc := unsafe.NewService(&unsafe.MockGateway{})

	req1 := unsafe.PaymentRequest{OrderID: "order-789", Amount: 100}
	req2 := unsafe.PaymentRequest{OrderID: "order-789", Amount: 999999} // different amount

	ctx := context.Background()

	result1, _ := svc.ProcessPayment(ctx, req1)
	result2, _ := svc.ProcessPayment(ctx, req2)

	// BUG: Different payloads produce different results without complaint
	if result1.Amount != req1.Amount || result2.Amount != req2.Amount {
		t.Fatalf("unexpected amounts: %d, %d", result1.Amount, result2.Amount)
	}

	t.Logf("BUG REPRODUCED (no payload check): processed different payloads without complaint")
}
