package tests

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"compat/internal/compat"
)

func TestConcurrency(t *testing.T) {
	store := compat.NewMemoryStore()
	flags := compat.NewFeatureFlags()
	obs := compat.NewObservability()
	svc := compat.NewService(store, flags, obs)

	// Pre-populate legacy users
	for i := 0; i < 50; i++ {
		_, _ = svc.CreateUser(fmt.Sprintf("User%d", i), fmt.Sprintf("+628%d", i))
	}

	flags.SetWriteMode(compat.WriteDual)
	flags.SetReadMode(compat.ReadFallback)

	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 1. Concurrent Legacy Writers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(writerID int) {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				_, _ = svc.CreateUser(fmt.Sprintf("Writer%d_User%d", writerID, j), "+628X")
			}
		}(i)
	}

	// 2. Concurrent Legacy and Modern Readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 1; j <= 50; j++ {
				_, _ = svc.GetLegacyUser(j)
				_, _ = svc.GetModernUser(j)
			}
		}()
	}

	// 3. Concurrent Backfill Worker
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, _ = svc.BackfillWorker().RunAll(ctx)
	}()

	// 4. Concurrent Drift Reconciliation
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			_, _ = svc.ReconcileData()
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Wait for all goroutines
	wg.Wait()

	// Verify no panics and data integrity is maintained
	if obs.DualWriteErrors.Load() > 0 {
		t.Errorf("Experienced %d dual write errors under concurrency", obs.DualWriteErrors.Load())
	}
}
