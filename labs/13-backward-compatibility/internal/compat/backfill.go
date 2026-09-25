package compat

import (
	"context"
	"fmt"
	"sync"
)

type BackfillCheckpoint struct {
	mu              sync.Mutex
	LastProcessedID int
	TotalMigrated   int
	IsComplete      bool
}

type BackfillWorker struct {
	store      *MemoryStore
	obs        *Observability
	batchSize  int
	checkpoint *BackfillCheckpoint
}

func NewBackfillWorker(store *MemoryStore, obs *Observability, batchSize int) *BackfillWorker {
	if batchSize <= 0 {
		batchSize = 100
	}
	return &BackfillWorker{
		store:      store,
		obs:        obs,
		batchSize:  batchSize,
		checkpoint: &BackfillCheckpoint{},
	}
}

// RunBatch runs a single batch migration from the checkpoint
func (b *BackfillWorker) RunBatch(ctx context.Context) (int, bool, error) {
	b.checkpoint.mu.Lock()
	defer b.checkpoint.mu.Unlock()

	if b.checkpoint.IsComplete {
		return 0, true, nil
	}

	ids := b.store.GetUserIDs(b.checkpoint.LastProcessedID, b.batchSize)
	if len(ids) == 0 {
		b.checkpoint.IsComplete = true
		return 0, true, nil
	}

	migratedInBatch := 0
	for _, id := range ids {
		select {
		case <-ctx.Done():
			return migratedInBatch, false, ctx.Err()
		default:
		}

		user, err := b.store.GetUser(id)
		if err != nil {
			continue
		}

		// Only backfill if legacy phone exists and user_phones is empty
		if user.Phone != nil && *user.Phone != "" {
			phones, _ := b.store.GetPhones(id)
			if len(phones) == 0 {
				_, err := b.store.SavePhoneEntry(id, *user.Phone, true)
				if err != nil {
					return migratedInBatch, false, fmt.Errorf("failed backfilling user %d: %w", id, err)
				}
				migratedInBatch++
				b.obs.BackfillProcessed.Add(1)
			}
		}
		b.checkpoint.LastProcessedID = id
	}

	b.checkpoint.TotalMigrated += migratedInBatch

	if len(ids) < b.batchSize {
		b.checkpoint.IsComplete = true
	}

	return migratedInBatch, b.checkpoint.IsComplete, nil
}

// RunAll executes batches until completion
func (b *BackfillWorker) RunAll(ctx context.Context) (int, error) {
	total := 0
	for {
		migrated, done, err := b.RunBatch(ctx)
		if err != nil {
			return total, err
		}
		total += migrated
		if done {
			break
		}
	}
	return total, nil
}

func (b *BackfillWorker) GetProgress() (lastID int, total int, done bool) {
	b.checkpoint.mu.Lock()
	defer b.checkpoint.mu.Unlock()
	return b.checkpoint.LastProcessedID, b.checkpoint.TotalMigrated, b.checkpoint.IsComplete
}
