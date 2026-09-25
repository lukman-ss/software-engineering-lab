package tests

import (
	"context"
	"sync"
	"testing"
	"time"

	"labs/18-deadlock/internal/bank"
	"labs/18-deadlock/internal/transfer"
)

func TestDeadlockOccurrence(t *testing.T) {
	accA := bank.NewAccount("A", 100)
	accB := bank.NewAccount("B", 100)

	var wg sync.WaitGroup
	errs := make([]error, 2)

	wg.Add(2)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
		defer cancel()
		errs[0] = transfer.TransferNaive(ctx, accA, accB, 10, 10*time.Millisecond)
	}()

	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
		defer cancel()
		errs[1] = transfer.TransferNaive(ctx, accB, accA, 10, 10*time.Millisecond)
	}()

	wg.Wait()

	hasDeadlock := (errs[0] == bank.ErrDeadlock) || (errs[1] == bank.ErrDeadlock)
	if !hasDeadlock {
		t.Fatalf("expected at least one deadlock victim, got errs: %v, %v", errs[0], errs[1])
	}
}

func TestLockOrderingPreventsDeadlock(t *testing.T) {
	accA := bank.NewAccount("A", 100)
	accB := bank.NewAccount("B", 100)

	var wg sync.WaitGroup
	errs := make([]error, 2)

	wg.Add(2)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		errs[0] = transfer.TransferOrdered(ctx, accA, accB, 10, 5*time.Millisecond)
	}()

	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		errs[1] = transfer.TransferOrdered(ctx, accB, accA, 10, 5*time.Millisecond)
	}()

	wg.Wait()

	if errs[0] != nil || errs[1] != nil {
		t.Fatalf("expected both transfers to succeed, got: %v, %v", errs[0], errs[1])
	}

	if accA.Balance != 100 || accB.Balance != 100 {
		t.Fatalf("unexpected balance: A=%d, B=%d", accA.Balance, accB.Balance)
	}
}

func TestRetryRecoversDeadlock(t *testing.T) {
	accA := bank.NewAccount("A", 100)
	accB := bank.NewAccount("B", 100)

	var wg sync.WaitGroup
	errs := make([]error, 2)

	wg.Add(2)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		errs[0] = transfer.TransferWithRetry(ctx, accA, accB, 10, 5*time.Millisecond, 10)
	}()

	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		errs[1] = transfer.TransferWithRetry(ctx, accB, accA, 10, 5*time.Millisecond, 10)
	}()

	wg.Wait()

	if errs[0] != nil || errs[1] != nil {
		t.Fatalf("expected retry to recover both transfers, got: %v, %v", errs[0], errs[1])
	}
}

func TestTransactionDurationImpact(t *testing.T) {
	// Runs multiple naive transfers with long vs zero duration to observe deadlock rate
	countDeadlocks := func(delay time.Duration, iterations int) int {
		deadlocks := 0
		for i := 0; i < iterations; i++ {
			accA := bank.NewAccount("A", 100)
			accB := bank.NewAccount("B", 100)

			var wg sync.WaitGroup
			errs := make([]error, 2)

			wg.Add(2)
			go func() {
				defer wg.Done()
				ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
				defer cancel()
				errs[0] = transfer.TransferNaive(ctx, accA, accB, 1, delay)
			}()

			go func() {
				defer wg.Done()
				ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
				defer cancel()
				errs[1] = transfer.TransferNaive(ctx, accB, accA, 1, delay)
			}()

			wg.Wait()
			if errs[0] == bank.ErrDeadlock || errs[1] == bank.ErrDeadlock {
				deadlocks++
			}
		}
		return deadlocks
	}

	longDeadlocks := countDeadlocks(5*time.Millisecond, 10)
	shortDeadlocks := countDeadlocks(0, 10)

	if longDeadlocks < shortDeadlocks {
		t.Fatalf("expected long transaction duration to have >= deadlocks than short: long=%d, short=%d", longDeadlocks, shortDeadlocks)
	}
}

