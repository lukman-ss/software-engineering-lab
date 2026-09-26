package deadline

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestExecuteWithBudget_Success(t *testing.T) {
	ctx := context.Background()
	err := ExecuteWithBudget(ctx, 100*time.Millisecond, func(c context.Context) error {
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestExecuteWithBudget_Timeout(t *testing.T) {
	ctx := context.Background()
	err := ExecuteWithBudget(ctx, 20*time.Millisecond, func(c context.Context) error {
		select {
		case <-time.After(100 * time.Millisecond):
			return nil
		case <-c.Done():
			return c.Err()
		}
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context.DeadlineExceeded, got %v", err)
	}
}

func TestExecuteWithBudget_ParentTimeoutInherited(t *testing.T) {
	parentCtx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := ExecuteWithBudget(parentCtx, 500*time.Millisecond, func(c context.Context) error {
		select {
		case <-time.After(100 * time.Millisecond):
			return nil
		case <-c.Done():
			return c.Err()
		}
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected parent context deadline propagation, got %v", err)
	}
}
