package deadline

import (
	"context"
	"errors"
	"time"
)

var ErrDeadlineExceeded = errors.New("deadline exceeded")

type WorkerFunc func(ctx context.Context) error

func ExecuteWithBudget(ctx context.Context, budget time.Duration, fn WorkerFunc) error {
	childCtx, cancel := context.WithTimeout(ctx, budget)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- fn(childCtx)
	}()

	select {
	case <-childCtx.Done():
		return childCtx.Err()
	case err := <-done:
		return err
	}
}
