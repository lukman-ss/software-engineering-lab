//go:build integration
// +build integration

package race

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/lukman-ss/software-engineering-lab/pkg/database"
)

// resetBookingTable truncates the service_bookings table for test isolation.
// Schema creation is handled by schema.sql via Docker PostgreSQL init.
func resetBookingTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, "TRUNCATE TABLE service_bookings RESTART IDENTITY CASCADE")
	if err != nil {
		return fmt.Errorf("reset service_bookings failed; ensure schema.sql has been initialized via docker compose up -d postgres: %w", err)
	}
	return nil
}

// TestPostgres_ConcurrentBooking menguji 500 concurrent requests booking exclusive slot dengan DB asli.
//
// Setup:
// - attempts = 500 goroutines
//
// Expected:
// - created = 1
// - conflict = 499
// - final_count = 1
func TestPostgres_ConcurrentBooking(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, database.FromEnv())
	if err != nil {
		t.Skipf("PostgreSQL not available: %v", err)
	}
	defer db.Close()

	if err := resetBookingTable(ctx, db); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	defer func() {
		if err := resetBookingTable(ctx, db); err != nil {
			t.Logf("cleanup warning: %v", err)
		}
	}()

	repo := NewPostgresBookingRepository(db)

	const attempts = 500
	var createdCount int
	var conflictCount int
	var unexpectedCount int
	var mu sync.Mutex

	ready, release := startGate(attempts)

	var wg sync.WaitGroup
	wg.Add(attempts)

	// All 500 goroutines compete for 1 exclusive slot (same branch)
	for i := 0; i < attempts; i++ {
		go func(id int) {
			defer wg.Done()
			ready <- struct{}{}
			<-release

			err := repo.CreateBooking(ctx, "branch-pg-01", "2026-09-01 09:00")
			if err != nil {
				// Assert error is our mapped constraint error
				if err == ErrDuplicateKey {
					mu.Lock()
					conflictCount++
					mu.Unlock()
				} else {
					t.Errorf("Request %d unexpected error: %v", id, err)
					mu.Lock()
					unexpectedCount++
					mu.Unlock()
				}
				return
			}
			mu.Lock()
			createdCount++
			mu.Unlock()
		}(i)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		t.Fatal("test timeout: potential goroutine leak or DB connection exhaustion")
	}

	// Verify final count
	finalCount, err := repo.CountBookings(ctx, "branch-pg-01", "2026-09-01 09:00")
	if err != nil {
		t.Fatalf("get final count: %v", err)
	}

	t.Logf("=== DB CONCURRENT BOOKING RESULTS ===")
	t.Logf("Created: %d", createdCount)
	t.Logf("Conflict (Unique Violation): %d", conflictCount)
	t.Logf("Unexpected Errors: %d", unexpectedCount)
	t.Logf("Final DB rows: %d", finalCount)

	// Assertions
	if unexpectedCount > 0 {
		t.Errorf("expected 0 unexpected errors, got %d", unexpectedCount)
	}

	if createdCount != 1 {
		t.Errorf("expected created = 1, got %d", createdCount)
	}

	if conflictCount != 499 {
		t.Errorf("expected conflict = 499, got %d", conflictCount)
	}

	if finalCount != 1 {
		t.Errorf("expected final_count = 1, got %d", finalCount)
	}

	t.Logf("✅ DB INVARIANT HOLDS: COUNT(bookings) = %d", finalCount)
}

// TestPostgres_Booking_SameBranchDifferentBranch verifies multi-branch invariant:
// - same branch + same date + same slot → only 1 booking (second gets unique violation)
// - different branch + same date + same slot → both allowed
func TestPostgres_Booking_SameBranchDifferentBranch(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, database.FromEnv())
	if err != nil {
		t.Skipf("PostgreSQL not available: %v", err)
	}
	defer db.Close()

	if err := resetBookingTable(ctx, db); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	repo := NewPostgresBookingRepository(db)

	// Same branch
	err = repo.CreateBooking(ctx, "branch-A", "2026-09-01 10:00")
	if err != nil {
		t.Fatalf("first booking for branch-A should succeed: %v", err)
	}

	err = repo.CreateBooking(ctx, "branch-A", "2026-09-01 10:00")
	if !errors.Is(err, ErrDuplicateKey) {
		t.Errorf("second booking for same branch should fail with duplicate key, got: %v", err)
	}

	// Different branch, same time
	err = repo.CreateBooking(ctx, "branch-B", "2026-09-01 10:00")
	if err != nil {
		t.Errorf("booking for branch-B same time should succeed, got: %v", err)
	}

	// Counts
	countA, _ := repo.CountBookings(ctx, "branch-A", "2026-09-01 10:00")
	countB, _ := repo.CountBookings(ctx, "branch-B", "2026-09-01 10:00")

	if countA != 1 || countB != 1 {
		t.Errorf("counts invalid: branch-A=%d, branch-B=%d", countA, countB)
	}
}

// TestPostgres_Booking_MultipleBranches proves that different branches can book the same slot.
// This directly demonstrates the UNIQUE(branch_id, service_date, slot_time) constraint.
func TestPostgres_Booking_MultipleBranches(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, database.FromEnv())
	if err != nil {
		t.Skipf("PostgreSQL not available: %v", err)
	}
	defer db.Close()

	if err := setupBookingTable(ctx, db); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	defer setupBookingTable(ctx, db)

	repo := NewPostgresBookingRepository(db)

	// Scenario: Branch 1 and Branch 2, same date, same slot → both should succeed
	err = repo.CreateBooking(ctx, "branch-01", "2026-09-01 09:00")
	if err != nil {
		t.Fatalf("Branch 1 first booking should succeed: %v", err)
	}
	err = repo.CreateBooking(ctx, "branch-02", "2026-09-01 09:00")
	if err != nil {
		t.Fatalf("Branch 2 booking should succeed (different branch): %v", err)
	}

	// Verify counts
	count1, _ := repo.CountBookings(ctx, "branch-01", "2026-09-01 09:00")
	count2, _ := repo.CountBookings(ctx, "branch-02", "2026-09-01 09:00")

	t.Logf("Branch 1 bookings: %d", count1)
	t.Logf("Branch 2 bookings: %d", count2)

	if count1 != 1 {
		t.Errorf("expected 1 booking for branch-01, got %d", count1)
	}
	if count2 != 1 {
		t.Errorf("expected 1 booking for branch-02, got %d", count2)
	}

	t.Logf("✅ UNIQUE(branch_id, service_date, slot_time) works: different branches can book same slot")
}
