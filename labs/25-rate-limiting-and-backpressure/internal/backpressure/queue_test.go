package backpressure

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestBoundedQueue_RejectionUnderLoad(t *testing.T) {
	capacity := 5
	workers := 1
	bq := NewBoundedQueue(capacity, workers)
	defer bq.Stop()

	blockCh := make(chan struct{})
	startedCh := make(chan struct{})
	// First job blocks worker
	err := bq.TrySubmit(func(ctx context.Context) error {
		close(startedCh)
		<-blockCh
		return nil
	})
	if err != nil {
		t.Fatalf("first job should submit: %v", err)
	}

	<-startedCh // Ensure worker has popped first job from channel

	// Fill queue buffer (capacity items)
	for i := 0; i < capacity; i++ {
		err := bq.TrySubmit(func(ctx context.Context) error { return nil })
		if err != nil {
			t.Fatalf("expected enqueue for item %d, got %v", i, err)
		}
	}

	// Next submit must immediately fail with ErrQueueFull
	err = bq.TrySubmit(func(ctx context.Context) error { return nil })
	if err != ErrQueueFull {
		t.Fatalf("expected ErrQueueFull, got %v", err)
	}

	close(blockCh)
}

func TestBoundedQueue_ConcurrencySafety(t *testing.T) {
	bq := NewBoundedQueue(20, 4)
	defer bq.Stop()

	var wg sync.WaitGroup
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = bq.TrySubmit(func(ctx context.Context) error {
				time.Sleep(2 * time.Millisecond)
				return nil
			})
		}()
	}

	wg.Wait()
	accepted, rejected, _, _ := bq.Stats()
	if accepted+rejected != 30 {
		t.Fatalf("expected total 30 requests recorded, got accepted=%d, rejected=%d", accepted, rejected)
	}
}

func TestBoundedQueue_SubmitAfterStop(t *testing.T) {
	bq := NewBoundedQueue(5, 1)
	bq.Stop()

	// Repeated Stop calls must be safe (idempotent)
	bq.Stop()

	err := bq.TrySubmit(func(ctx context.Context) error { return nil })
	if err != ErrQueueStopped {
		t.Fatalf("expected ErrQueueStopped, got %v", err)
	}
}
