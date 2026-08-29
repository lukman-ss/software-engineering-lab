//go:build integration
// +build integration

package race

import (
	"context"
	"database/sql"
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
			return 0, fmt.Errorf("stock habis")
		}
		return 0, fmt.Errorf("decrement stock: %w", err)
	}
	return newStock, nil
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

// TestPostgresUnsafe_LostUpdate menguji lost update dengan PostgreSQL.
//
// Timeline:
// T0: Transaction A membaca stock = 1
// T1: Transaction B membaca stock = 1 (stale read)
// T2: Transaction A menulis stock = 0
// T3: Transaction B menulis stock = 0
//
// Result:successful_sales = 2, final_stock = 0
// Broken invariant: 1 != 2 + 0
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
	initialStock := 1

	// Setup
	if err := setupTestInventory(ctx, db, productID, initialStock); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	defer cleanupTestInventory(ctx, db, productID)

	repo := NewPostgresUnsafeRepository(db)
	service := NewInventoryService(repo)

	var wg sync.WaitGroup
	var successCount int
	var mu sync.Mutex

	wg.Add(2)

	// Request A dan B
	for _, name := range []string{"A", "B"} {
		go func(reqName string) {
			defer wg.Done()

			// TrySell menggunakan check-then-act
			err := service.TrySell(ctx, productID)
			if err == nil {
				mu.Lock()
				successCount++
				mu.Unlock()
				t.Logf("Request %s: SUCCESS", reqName)
			} else {
				t.Logf("Request %s: FAILED (%v)", reqName, err)
			}
		}(name)
	}

	wg.Wait()

	// Get final stock
	finalStock, _ := repo.GetStock(ctx, productID)

	t.Logf("=== POSTGRES UNSAFE RESULTS ===")
	t.Logf("Initial stock: %d", initialStock)
	t.Logf("Successful sales: %d", successCount)
	t.Logf("Final stock: %d", finalStock)

	// Verify lost update occurred
	if successCount == 2 && finalStock == 0 {
		t.Logf("❌ LOST UPDATE DETECTED via PostgreSQL: 2 sales, 0 final, 1 initial")
	}
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

	var wg sync.WaitGroup
	var successCount int
	var rejectedCount int
	var mu sync.Mutex

	wg.Add(attempts)

	for i := 0; i < attempts; i++ {
		go func() {
			defer wg.Done()
			_, err := repo.DecrementStock(ctx, productID)
			if err != nil {
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

	finalStock, _ := repo.GetStock(ctx, productID)

	t.Logf("=== POSTGRES ATOMIC RESULTS ===")
	t.Logf("Initial stock: %d", initialStock)
	t.Logf("Successful: %d", successCount)
	t.Logf("Rejected: %d", rejectedCount)
	t.Logf("Final stock: %d", finalStock)

	// Assertions
	if finalStock < 0 {
		t.Errorf("BUG: final_stock < 0: %d", finalStock)
	}

	if successCount > initialStock {
		t.Errorf("BUG: successful_sales (%d) > initial_stock (%d)", successCount, initialStock)
	}

	// Main invariant check
	if successCount != initialStock {
		t.Errorf("expected success = %d, got %d", initialStock, successCount)
	}
	if rejectedCount != 400 {
		t.Errorf("expected rejected = 400, got %d", rejectedCount)
	}
	if finalStock != 0 {
		t.Errorf("expected final_stock = 0, got %d", finalStock)
	}

	t.Logf("✅ POSTGRES ATOMIC INTEGRATION TEST PASSED")
}
