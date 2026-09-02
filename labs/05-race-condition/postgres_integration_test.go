//go:build integration
// +build integration

package race

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/lukman-ss/software-engineering-lab/pkg/database"
)

// PostgresUnsafeRepository implementasi unsafe check-then-act dengan PostgreSQL.
// Mendemonstrasikan lost update di level database.
type PostgresUnsafeRepository struct {
	db *sql.DB
}

func NewPostgresUnsafeRepository(db *sql.DB) *PostgresUnsafeRepository {
	return &PostgresUnsafeRepository{db: db}
}

// GetStock melakukan SELECT stock
func (r *PostgresUnsafeRepository) GetStock(ctx context.Context, productID string) (int, error) {
	var stock int
	err := r.db.QueryRowContext(ctx, "SELECT stock FROM inventory_products WHERE id = $1", productID).Scan(&stock)
	if err != nil {
		return 0, fmt.Errorf("get stock: %w", err)
	}
	return stock, nil
}

// SetStock melakukan UPDATE stock
func (r *PostgresUnsafeRepository) SetStock(ctx context.Context, productID string, stock int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE inventory_products SET stock = $1 WHERE id = $2", stock, productID)
	if err != nil {
		return fmt.Errorf("set stock: %w", err)
	}
	return nil
}

// PostgresAtomicRepository implementasi safe atomic conditional update.
type PostgresAtomicRepository struct {
	db *sql.DB
}

func NewPostgresAtomicRepository(db *sql.DB) *PostgresAtomicRepository {
	return &PostgresAtomicRepository{db: db}
}

// DecrementStock mengurangi stock secara atomik dengan kondisi.
// SQL: UPDATE products SET stock = stock - 1 WHERE id = $1 AND stock > 0 RETURNING stock
func (r *PostgresAtomicRepository) DecrementStock(ctx context.Context, productID string) (int, error) {
	var newStock int
	err := r.db.QueryRowContext(ctx,
		"UPDATE inventory_products SET stock = stock - 1 WHERE id = $1 AND stock > 0 RETURNING stock",
		productID).Scan(&newStock)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrOutOfStock
		}
		return 0, fmt.Errorf("decrement stock: %w", err)
	}
	return newStock, nil
}

// GetStock reads the current stock value (for verification purposes).
func (r *PostgresAtomicRepository) GetStock(ctx context.Context, productID string) (int, error) {
	var stock int
	err := r.db.QueryRowContext(ctx,
		"SELECT stock FROM inventory_products WHERE id = $1",
		productID).Scan(&stock)
	if err != nil {
		return 0, fmt.Errorf("get stock: %w", err)
	}
	return stock, nil
}

// setupTestInventory membuat data test di database.
func setupTestInventory(ctx context.Context, db *sql.DB, productID string, initialStock int) error {
	_, err := db.ExecContext(ctx,
		`INSERT INTO inventory_products (id, name, stock)
		 VALUES ($1, 'Test Product', $2)
		 ON CONFLICT (id) DO UPDATE SET stock = $2, updated_at = NOW()`,
		productID, initialStock)
	return err
}

// cleanupTestInventory menghapus data test.
func cleanupTestInventory(ctx context.Context, db *sql.DB, productID string) error {
	_, err := db.ExecContext(ctx, "DELETE FROM inventory_products WHERE id = $1", productID)
	return err
}

// TestPostgresUnsafe_LostUpdate menguji lost update dengan PostgreSQL secara deterministik.
//
// Menggunakan barrier synchronization untuk memaksa interleaving tertentu:
// A READ → B READ → A WRITE → B WRITE
//
// Timeline yang direproduksi:
// T0: Transaction A  READ stock = 1
// T1: Transaction B  READ stock = 1   (stale read — A belum WRITE)
// T4: Transaction A  WRITE stock = 0
// T5: Transaction B  WRITE stock = 0 (overwrite)
//
// Expected outcome:
//
//	successful_sales = 2
//	final_stock = 0
//	1 != 2 + 0 (invariant rusak)
//
// Test PASS ketika lost update berhasil direproduksi (membuktikan bahwa implementation unsafe bermasalah).
func TestPostgresUnsafe_LostUpdate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to PostgreSQL
	db, err := database.Connect(ctx, database.FromEnv())
	if err != nil {
		t.Skipf("PostgreSQL not available: %v", err)
	}
	defer db.Close()

	productID := "unit-oli-mesin"
	const initialStock = 1

	// Setup
	if err := setupTestInventory(ctx, db, productID, initialStock); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	defer cleanupTestInventory(ctx, db, productID)

	repo := NewPostgresUnsafeRepository(db)

	// Synchronization barriers — pure signal channels, no payload needed.
	// Semantic: signal is sent AFTER the named phase completes.
	aReadDone := make(chan struct{})
	bReadDone := make(chan struct{})
	aWriteDone := make(chan struct{})

	// Result channel for deterministic error reporting from goroutines.
	// Main goroutine owns all t.Fatal / t.Errorf calls.
	type unsafeWorkerResult struct {
		name string
		err  error
	}
	resultsCh := make(chan unsafeWorkerResult, 2)

	var mu sync.Mutex
	var successCount int

	// Transaction A
	go func() {
		// A: READ first
		stockA, err := repo.GetStock(ctx, productID)
		if err != nil {
			resultsCh <- unsafeWorkerResult{name: "A", err: fmt.Errorf("READ: %w", err)}
			return
		}
		t.Logf("Request A: READ stock = %d", stockA)

		// Signal: A READ is complete — B may now read
		select {
		case aReadDone <- struct{}{}:
		case <-ctx.Done():
			resultsCh <- unsafeWorkerResult{name: "A", err: ctx.Err()}
			return
		}

		// A: wait for B to finish READ before writing (barrier)
		select {
		case <-bReadDone:
		case <-ctx.Done():
			resultsCh <- unsafeWorkerResult{name: "A", err: fmt.Errorf("timeout waiting for B READ: %w", ctx.Err())}
			return
		}

		// A: CHECK + CALCULATE + WRITE
		if stockA > 0 {
			if err := repo.SetStock(ctx, productID, stockA-1); err != nil {
				resultsCh <- unsafeWorkerResult{name: "A", err: fmt.Errorf("WRITE: %w", err)}
				return
			}
			t.Logf("Request A: WRITE stock = %d", stockA-1)
			mu.Lock()
			successCount++
			mu.Unlock()
		}

		// Signal: A WRITE is complete — B may now write
		select {
		case aWriteDone <- struct{}{}:
		case <-ctx.Done():
			resultsCh <- unsafeWorkerResult{name: "A", err: ctx.Err()}
			return
		}

		resultsCh <- unsafeWorkerResult{name: "A"}
	}()

	// Transaction B
	go func() {
		// B: wait for A READ to complete (barrier)
		select {
		case <-aReadDone:
		case <-ctx.Done():
			resultsCh <- unsafeWorkerResult{name: "B", err: ctx.Err()}
			return
		}

		// B: READ (A has read but not yet written — stale read guaranteed)
		stockB, err := repo.GetStock(ctx, productID)
		if err != nil {
			resultsCh <- unsafeWorkerResult{name: "B", err: fmt.Errorf("READ: %w", err)}
			return
		}
		t.Logf("Request B: READ stock = %d", stockB)

		// Signal: B READ is complete — A may now write
		select {
		case bReadDone <- struct{}{}:
		case <-ctx.Done():
			resultsCh <- unsafeWorkerResult{name: "B", err: ctx.Err()}
			return
		}

		// B: wait for A WRITE to complete before writing stale value (barrier)
		select {
		case <-aWriteDone:
		case <-ctx.Done():
			resultsCh <- unsafeWorkerResult{name: "B", err: fmt.Errorf("timeout waiting for A WRITE: %w", ctx.Err())}
			return
		}

		// B: CHECK + CALCULATE + WRITE (stale value — lost update)
		if stockB > 0 {
			if err := repo.SetStock(ctx, productID, stockB-1); err != nil {
				resultsCh <- unsafeWorkerResult{name: "B", err: fmt.Errorf("WRITE: %w", err)}
				return
			}
			t.Logf("Request B: WRITE stock = %d (stale — lost update)", stockB-1)
			mu.Lock()
			successCount++
			mu.Unlock()
		}

		resultsCh <- unsafeWorkerResult{name: "B"}
	}()

	// Drain results with bounded timeout
	for i := 0; i < 2; i++ {
		select {
		case res := <-resultsCh:
			if res.err != nil {
				t.Fatalf("worker %s failed: %v", res.name, res.err)
			}
		case <-ctx.Done():
			t.Fatal("test timeout: deadlock in synchronization barriers")
		}
	}

	// Get final stock
	finalStock, err := repo.GetStock(ctx, productID)
	if err != nil {
		t.Fatalf("get final stock: %v", err)
	}

	t.Logf("\n=== POSTGRES UNSAFE RESULTS ===")
	t.Logf("Initial stock:         %d", initialStock)
	t.Logf("Successful operations: %d", successCount)
	t.Logf("Final stock:           %d", finalStock)
	t.Logf("Invariant check:       %d == %d + %d → %s",
		initialStock, successCount, finalStock,
		func() string {
			if initialStock == successCount+finalStock {
				return "HOLDS (unexpected)"
			}
			return "BROKEN"
		}())

	// ASSERTION: Lost update harus terjadi (invariant rusak)
	if successCount != 2 {
		t.Fatalf("expected successful operations = 2 (both wrote), got %d", successCount)
	}
	if finalStock != 0 {
		t.Fatalf("expected final_stock = 0, got %d", finalStock)
	}

	// Verify invariant is BROKEN (ini yang menjadi bukti bahwa implementation unsafe)
	actualSum := successCount + finalStock
	if actualSum == initialStock {
		t.Fatalf("INVARIANT ERROR: expected invariant to be BROKEN but %d + %d == %d", successCount, finalStock, initialStock)
	}
	t.Logf("✅ LOST UPDATE CONFIRMED: %d != %d + %d (initial != success + final)",
		initialStock, successCount, finalStock)
}

// TestPostgresAtomic_ConcurrentUpdate menguji atomic update dengan PostgreSQL.
//
// Setup:
// - initial_stock = 100
// - attempts = 500 goroutines
//
// Expected:
// - success = 100
// - rejected = 400
// - final_stock = 0
//
// Assert invariant:
// - final_stock >= 0
// - successful <= initial_stock
// - initial_stock == successful + final_stock
func TestPostgresAtomic_ConcurrentUpdate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, database.FromEnv())
	if err != nil {
		t.Skipf("PostgreSQL not available: %v", err)
	}
	defer db.Close()

	productID := "unit-oli-mesin"
	const initialStock = 100
	const attempts = 500

	// Setup
	if err := setupTestInventory(ctx, db, productID, initialStock); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	defer cleanupTestInventory(ctx, db, productID)

	repo := NewPostgresAtomicRepository(db)

	// errorCount tracks unexpected DB failures (distinct from expected out-of-stock).
	var mu sync.Mutex
	var wg sync.WaitGroup
	var successCount int
	var rejectedCount int
	var errorCount int

	ready, release := startGate(attempts)

	wg.Add(attempts)

	for i := 0; i < attempts; i++ {
		go func() {
			defer wg.Done()
			ready <- struct{}{} // signal ready
			<-release           // wait for synchronized start

			_, err := repo.DecrementStock(ctx, productID)
			if err != nil {
				// Distinguish expected out-of-stock from unexpected DB errors.
				if !errors.Is(err, ErrOutOfStock) {
					t.Errorf("unexpected DB error during atomic update: %v", err)
					mu.Lock()
					errorCount++
					mu.Unlock()
					return
				}
				mu.Lock()
				rejectedCount++
				mu.Unlock()
				return
			}
			mu.Lock()
			successCount++
			mu.Unlock()
		}()
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

	finalStock, err := repo.GetStock(ctx, productID)
	if err != nil {
		t.Fatalf("get final stock: %v", err)
	}

	t.Logf("=== POSTGRES ATOMIC RESULTS ===")
	t.Logf("Initial stock: %d", initialStock)
	t.Logf("Successful: %d", successCount)
	t.Logf("Rejected (out-of-stock): %d", rejectedCount)
	t.Logf("Error (unexpected DB failure): %d", errorCount)
	t.Logf("Final stock: %d", finalStock)

	// Assertions
	if finalStock < 0 {
		t.Errorf("BUG: final_stock < 0: %d", finalStock)
	}

	if successCount > initialStock {
		t.Errorf("BUG: successful_sales (%d) > initial_stock (%d)", successCount, initialStock)
	}

	// Main invariant check
	if successCount+finalStock != initialStock {
		t.Fatalf("INVARIANT ERROR: initial_stock (%d) != successful (%d) + final_stock (%d)", initialStock, successCount, finalStock)
	}

	if errorCount > 0 {
		t.Errorf("expected 0 unexpected DB errors, got %d", errorCount)
	}
	if rejectedCount != 400 {
		t.Errorf("expected rejected = 400, got %d", rejectedCount)
	}
	if finalStock != 0 {
		t.Errorf("expected final_stock = 0, got %d", finalStock)
	}

	t.Logf("✅ POSTGRES ATOMIC INTEGRATION TEST PASSED: invariant holds")
}
