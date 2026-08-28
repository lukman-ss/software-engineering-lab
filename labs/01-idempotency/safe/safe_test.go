package safe_test

import (
	"context"
	"errors"
	"testing"

	"github.com/lukman/software-engineer-lab/labs/01-idempotency/safe"
)

// 1. Same key + same payload -> side effect only once
func TestSameIdempotencyKeyExecutesOnlyOnce(t *testing.T) {
	gw := safe.MockGateway()
	svc := safe.NewService(gw)
	ctx := context.Background()

	req := safe.PaymentRequest{Amount: 10000}

	// First request
	_, status1, err := svc.ProcessPayment(ctx, "POST", "/payments", "key-same-1", req)
	if err != nil || status1 != 200 {
		t.Fatalf("first payment failed: %v (status %d)", err, status1)
	}

	// Second request with same key & payload
	_, status2, err := svc.ProcessPayment(ctx, "POST", "/payments", "key-same-1", req)
	if err != nil || status2 != 200 {
		t.Fatalf("second payment failed: %v (status %d)", err, status2)
	}

	if safe.GatewayChargeCount(gw) != 1 {
		t.Fatalf("expected gateway charged exactly 1 time, got %d", safe.GatewayChargeCount(gw))
	}
}

// 2. Completed retry -> previous response replayed
func TestCompletedRetryReplaysPreviousResponse(t *testing.T) {
	gw := safe.MockGateway()
	svc := safe.NewService(gw)
	ctx := context.Background()

	req := safe.PaymentRequest{Amount: 50000}

	res1, _, _ := svc.ProcessPayment(ctx, "POST", "/payments", "key-replay", req)
	res2, _, _ := svc.ProcessPayment(ctx, "POST", "/payments", "key-replay", req)

	if res1.PaymentID != res2.PaymentID {
		t.Fatalf("expected replayed response to have payment ID %s, got %s", res1.PaymentID, res2.PaymentID)
	}
	if res1.ExternalID != res2.ExternalID {
		t.Fatalf("expected replayed response to have external ID %s, got %s", res1.ExternalID, res2.ExternalID)
	}
}

// 3. Same key + different payload -> conflict
func TestSameIdempotencyKeyDifferentPayloadReturnsConflict(t *testing.T) {
	gw := safe.MockGateway()
	svc := safe.NewService(gw)
	ctx := context.Background()

	req1 := safe.PaymentRequest{Amount: 100}
	_, status1, _ := svc.ProcessPayment(ctx, "POST", "/payments", "key-conflict", req1)
	if status1 != 200 {
		t.Fatalf("first request failed: status %d", status1)
	}

	req2 := safe.PaymentRequest{Amount: 99999}
	_, status2, err := svc.ProcessPayment(ctx, "POST", "/payments", "key-conflict", req2)

	if !errors.Is(err, safe.ErrIdempotencyKeyConflict) {
		t.Fatalf("expected ErrIdempotencyKeyConflict, got %v", err)
	}
	if status2 != 409 {
		t.Fatalf("expected status 409, got %d", status2)
	}

	if safe.GatewayChargeCount(gw) != 1 {
		t.Fatalf("expected gateway called 1 time, got %d", safe.GatewayChargeCount(gw))
	}
}

// 4. Missing key -> validation error
func TestMissingIdempotencyKeyReturnsBadRequest(t *testing.T) {
	gw := safe.MockGateway()
	svc := safe.NewService(gw)
	ctx := context.Background()

	req := safe.PaymentRequest{Amount: 1000}
	_, status, err := svc.ProcessPayment(ctx, "POST", "/payments", "", req)

	if !errors.Is(err, safe.ErrMissingIdempotencyKey) {
		t.Fatalf("expected ErrMissingIdempotencyKey, got %v", err)
	}
	if status != 400 {
		t.Fatalf("expected status 400, got %d", status)
	}
}

// 5. Different key -> treated as different logical operation
func TestDifferentIdempotencyKeysTreatedAsNewOperations(t *testing.T) {
	gw := safe.MockGateway()
	svc := safe.NewService(gw)
	ctx := context.Background()

	req1 := safe.PaymentRequest{Amount: 1000}
	req2 := safe.PaymentRequest{Amount: 2000}

	_, _, _ = svc.ProcessPayment(ctx, "POST", "/payments", "key-op-1", req1)
	_, _, _ = svc.ProcessPayment(ctx, "POST", "/payments", "key-op-2", req2)

	if safe.GatewayChargeCount(gw) != 2 {
		t.Fatalf("expected gateway called 2 times, got %d", safe.GatewayChargeCount(gw))
	}
}
