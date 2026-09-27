package backpressure

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrQueueFull    = errors.New("backpressure: queue capacity exceeded")
	ErrQueueStopped = errors.New("backpressure: queue stopped")
)

type Job func(ctx context.Context) error

type BoundedQueue struct {
	capacity  int
	queue     chan Job
	workers   int
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	stopMu    sync.RWMutex
	stopped   bool
	accepted  atomic.Int64
	rejected  atomic.Int64
	processed atomic.Int64
}

func NewBoundedQueue(capacity int, workers int) *BoundedQueue {
	ctx, cancel := context.WithCancel(context.Background())
	bq := &BoundedQueue{
		capacity: capacity,
		queue:    make(chan Job, capacity),
		workers:  workers,
		ctx:      ctx,
		cancel:   cancel,
	}

	for i := 0; i < workers; i++ {
		bq.wg.Add(1)
		go bq.workerLoop()
	}

	return bq
}

func (bq *BoundedQueue) workerLoop() {
	defer bq.wg.Done()
	for {
		select {
		case <-bq.ctx.Done():
			return
		case job, ok := <-bq.queue:
			if !ok {
				return
			}
			_ = job(bq.ctx)
			bq.processed.Add(1)
		}
	}
}

// TrySubmit enqueues the job if capacity allows, otherwise drops immediately with ErrQueueFull.
// Returns ErrQueueStopped if the queue has been stopped.
func (bq *BoundedQueue) TrySubmit(job Job) error {
	bq.stopMu.RLock()
	defer bq.stopMu.RUnlock()

	if bq.stopped {
		return ErrQueueStopped
	}

	select {
	case <-bq.ctx.Done():
		return ErrQueueStopped
	case bq.queue <- job:
		bq.accepted.Add(1)
		return nil
	default:
		bq.rejected.Add(1)
		return ErrQueueFull
	}
}

func (bq *BoundedQueue) Stats() (accepted, rejected, processed int64, queueLen int) {
	return bq.accepted.Load(), bq.rejected.Load(), bq.processed.Load(), len(bq.queue)
}

func (bq *BoundedQueue) Stop() {
	bq.stopMu.Lock()
	if bq.stopped {
		bq.stopMu.Unlock()
		return
	}
	bq.stopped = true
	bq.cancel()
	close(bq.queue)
	bq.stopMu.Unlock()

	bq.wg.Wait()
}
