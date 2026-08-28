package safe_test

import (
	"context"
	"testing"

	"github.com/lukman/software-engineer-lab/labs/01-idempotency/safe"
)

// TestUnsafeImplementationAllowsDuplicatePayment demonstrates that without idempotency key,
// sequential duplicate requests create duplicate side effects.
func TestUnsafeImplementationAllowsDuplicatePayment(t *testing.T) {
	gw := safe.MockGateway()
	svc := safe.NewService(gw)

	req := safe.PaymentRequest{
		OrderID: "order-123",
		Amount:  10000,
	}

	ctx := context.Background()

	// First request
	result1, status1, err := svc.ProcessPayment(ctx, "POST", "/payments", "key-1", req)
	if err != nil {
		t.Fatalf("first payment failed: %v", err)
	}
	if status1 != 200 {
		t.Fatalf("expected status 200, got %d", status1)
	}

	// Second request with DIFFERENT key = new operation
	result2, status2, err := svc.ProcessPayment(ctx, "POST", "/payments", "key-2", req)
	if err != nil {
		t.Fatalf("second payment failed: %v", err)
	}
	if status2 != 200 {
		t.Fatalf("expected status 200, got %d", status2)
	}

	// Different keys produce different payment IDs (duplicate side effect)
	if result1.PaymentID == result2.PaymentID {
		t.Fatal("expected different payment IDs for different keys")
	}

	// Gateway charged twice (duplicate side effect)
	if safe.GatewayChargeCount(gw) != 2 {
		t.Fatalf("expected gateway charged 2 times (duplicate), got %d", safe.GatewayChargeCount(gw))
	}

	t.Logf("UNSAFE: Duplicate keys create duplicate payments:")
	t.Logf("  Payment 1: %s", result1.PaymentID)
	t.Logf("  Payment 2: %s", result2.PaymentID)
}

// TestSameIdempotencyKeyDoesNotCreateSecondPayment verifies that using the same
// idempotency key with the same payload only charges the gateway once.
func TestSameIdempotencyKeyDoesNotCreateSecondPayment(t *testing.T) {
	gw := safe.MockGateway()
	svc := safe.NewService(gw)

	req := safe.PaymentRequest{
		OrderID: "order-456",
		Amount:  50000,
	}

	ctx := context.Background()

	// First request
	_, status1, err := svc.ProcessPayment(ctx, "POST", "/payments", "pay-key-789", req)
	if err != nil {
		t.Fatalf("first payment failed: %v", err)
	}
	if status1 != 200 {
		t.Fatalf("expected status 200, got %d", status1)
	}

	// Second request with SAME key and SAME payload = response replay
	_, status2, err := svc.ProcessPayment(ctx, "POST", "/payments", "pay-key-789", req)
	if err != nil {
		t.Fatalf("retry payment failed: %v", err)
	}
	if status2 != 200 {
		t.Fatalf("expected status 200 on replay, got %d", status2)
	}

	// Gateway still only charged once
	if safe.GatewayChargeCount(gw) != 1 {
		t.Fatalf("expected gateway charged 1 time, got %d", safe.GatewayChargeCount(gw))
	}
}

// TestRetryWithSameIdempotencyKeyReturnsCachedResponse verifies that retrying with
// the same idempotency key returns the original response without invoking the gateway again.
func TestRetryWithSameIdempotencyKeyReturnsCachedResponse(t *testing.T) {
	gw := safe.MockGateway()
	svc := safe.NewService(gw)

	req := safe.PaymentRequest{
		OrderID: "order-789",
		Amount:  75000,
	}

	ctx := context.Background()

	// First request
	result1, _, _ := svc.ProcessPayment(ctx, "POST", "/payments", "replay-key", req)

	// Retry with same key
	result2, _, _ := svc.ProcessPayment(ctx, "POST", "/payments", "replay-key", req)

	// Same response returned
	if result1.PaymentID != result2.PaymentID {
		t.Fatal("expected same payment ID from replay")
	}
	if result1.Amount != result2.Amount {
		t.Fatal("expected same amount from replay")
	}

	// No additional gateway charge
	if safe.GatewayChargeCount(gw) != 1 {
		t.Fatalf("expected 1 gateway call, got %d", safe.GatewayChargeCount(gw))
	}
}

// TestSameIdempotencyKeyWithDifferentPayloadReturns409Conflict verifies that using
// the same idempotency key with a different payload results in a conflict error.
func TestSameIdempotencyKeyWithDifferentPayloadReturns409Conflict(t *testing.T) {
	gw := safe.MockGateway()
	svc := safe.NewService(gw)
	ctx := context.Background()

	// First request: pay 100
	req1 := safe.PaymentRequest{OrderID: "order-conflict", Amount: 100}
	_, status1, _ := svc.ProcessPayment(ctx, "POST", "/payments", "conflict-key", req1)
	if status1 != 200 {
		t.Fatalf("expected status 200, got %d", status1)
	}

	// Second request: SAME key but different payload = 1000 (malicious or bug)
	req2 := safe.PaymentRequest{OrderID: "order-conflict", Amount: 1000}
	_, status2, err := svc.ProcessPayment(ctx, "POST", "/payments", "conflict-key", req2)

	// Should reject with 409 Conflict
	if err == nil {
		t.Fatal("expected conflict error for different payload")
	}
	if status2 != 409 {
		t.Fatalf("expected status 409 for conflict, got %d", status2)
	}

	// Gateway NOT charged again for the conflicting request
	if safe.GatewayChargeCount(gw) != 1 {
		t.Fatalf("expected 1 gateway call, got %d", safe.GatewayChargeCount(gw))
	}
}

// TestDifferentIdempotencyKeysTreatedAsDifferentOperations verifies that
// requests with different idempotency keys are treated as new operations.
func TestDifferentIdempotencyKeysTreatedAsDifferentOperations(t *testing.T) {
	gw := safe.MockGateway()
	svc := safe.NewService(gw)

	ctx := context.Background()

	req1 := safe.PaymentRequest{OrderID: "order-a", Amount: 1000}
	req2 := safe.PaymentRequest{OrderID: "order-b", Amount: 2000}

	// Different keys = different operations
	res1, _, _ := svc.ProcessPayment(ctx, "POST", "/payments", "key-op-1", req1)
	res2, _, _ := svc.ProcessPayment(ctx, "POST", "/payments", "key-op-2", req2)

	if res1.PaymentID == res2.PaymentID {
		t.Fatal("different keys should produce different payments")
	}

	if safe.GatewayChargeCount(gw) != 2 {
		t.Fatalf("expected 2 gateway calls, got %d", safe.GatewayChargeCount(gw))
	}
}