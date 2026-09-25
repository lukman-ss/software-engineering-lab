package transfer

import (
	"context"
	"time"
	"labs/18-deadlock/internal/bank"
)

// TransferNaive creates a potential deadlock by locking from then to.
func TransferNaive(ctx context.Context, from, to *bank.Account, amount int, delay time.Duration) error {
	if err := from.Lock(ctx); err != nil {
		return err
	}
	defer from.Unlock()

	time.Sleep(delay) // simulate transaction duration

	if err := to.Lock(ctx); err != nil {
		return err
	}
	defer to.Unlock()

	from.Balance -= amount
	to.Balance += amount
	return nil
}

// TransferOrdered prevents deadlock by always locking accounts in alphabetical order.
func TransferOrdered(ctx context.Context, acc1, acc2 *bank.Account, amount int, delay time.Duration) error {
	first, second := acc1, acc2
	if acc1.ID > acc2.ID {
		first, second = acc2, acc1
	}

	if err := first.Lock(ctx); err != nil {
		return err
	}
	defer first.Unlock()

	time.Sleep(delay)

	if err := second.Lock(ctx); err != nil {
		return err
	}
	defer second.Unlock()

	acc1.Balance -= amount
	acc2.Balance += amount
	return nil
}

// TransferWithRetry handles deadlocks by retrying the operation.
func TransferWithRetry(ctx context.Context, from, to *bank.Account, amount int, delay time.Duration, maxRetries int) error {
	for i := 0; i < maxRetries; i++ {
		// Use a short timeout for the attempt to simulate deadlock monitor aborting fast
		attemptCtx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
		err := TransferNaive(attemptCtx, from, to, amount, delay)
		cancel()

		if err == nil {
			return nil // Success
		}
		if err != bank.ErrDeadlock && err != context.DeadlineExceeded {
			return err // Some other error
		}
		
		// Backoff before retry
		time.Sleep(2 * time.Millisecond)
	}
	return bank.ErrDeadlock
}
