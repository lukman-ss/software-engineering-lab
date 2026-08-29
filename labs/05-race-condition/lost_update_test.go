package race

import (
	"context"
	"sync"
	"testing"
)

// UnsafeInventory adalah implementasi check-then-act yang TIDAK atomic.
type UnsafeInventory struct {
	mu    sync.RWMutex
	stock int
}

func (r *UnsafeInventory) GetStock(_ context.Context) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.stock, nil
}

func (r *UnsafeInventory) SetStock(_ context.Context, stock int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stock = stock
	return nil
}

func (r *UnsafeInventory) TrySellUnsafe(ctx context.Context) error {
	stock, err := r.GetStock(ctx)
	if err != nil {
		return err
	}
	if stock <= 0 {
		return ErrOutOfStock
	}
	newStock := stock - 1
	return r.SetStock(ctx, newStock)
}

// TestLostUpdate_Deterministic mereproduksi lost update secara deterministic.
//
// Timeline yang direproduksi:
// T0: Request A membaca stock = 1
// T1: Request B membaca stock = 1 (stale read)
// T2: Request A menghitung newStock = 0
// T3: Request B menghitung newStock = 0
// T4: Request A menulis stock = 0
// T5: Request B menulis stock = 0
//
// Di akhir:
// - final_stock = 0
// - successful_sales = 2 (keduanya melewati CHECK)
// - 1 != 2 + 0 (invariant rusak)
func TestLostUpdate_Deterministic(t *testing.T) {
	const productID = "Oli Mesin"
	initialStock := 1

	repo := &UnsafeInventory{}
	ctx := context.Background()

	// Inisialisasi stock
	repo.SetStock(ctx, initialStock)

	var wg sync.WaitGroup
	var mu sync.Mutex
	successfulSales := 0

	// Channel untuk sinkronisasi langkah-langkah
	// channel[0]: A selesai READ → B boleh READ
	// channel[1]: B selesai READ → A boleh hitung
	// channel[2]: A selesai CALCULATE → B boleh hitung
	// channel[3]: A selesai WRITE → B boleh write
	aReadDone := make(chan struct{})
	bReadDone := make(chan struct{})
	aCalcDone := make(chan struct{})
	bCalcDone := make(chan struct{})
	aWriteDone := make(chan struct{})

	// Goroutine A
	wg.Add(1)
	go func() {
		defer wg.Done()

		// T0: A reads
		stock, _ := repo.GetStock(ctx)
		t.Logf("Request A: READ stock = %d", stock)
		close(aReadDone)

		// Tunggu sampai B selesai READ
		<-bReadDone

		// T2: A calculates
		newStock := stock - 1
		t.Logf("Request A: CALCULATE new_stock = %d", newStock)
		close(aCalcDone)

		// Tunggu sampai B selesai CALCULATE
		<-bCalcDone

		// T4: A writes
		repo.SetStock(ctx, newStock)
		t.Logf("Request A: WRITE stock = %d", newStock)
		close(aWriteDone)

		mu.Lock()
		successfulSales++
		mu.Unlock()
	}()

	// Goroutine B
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Tunggu A selesai READ
		<-aReadDone

		// T1: B reads (baca yang sama dengan A)
		stock, _ := repo.GetStock(ctx)
		t.Logf("Request B: READ stock = %d", stock)
		close(bReadDone)

		// Tunggu A selesai CALCULATE
		<-aCalcDone

		// T3: B calculates
		newStock := stock - 1
		t.Logf("Request B: CALCULATE new_stock = %d", newStock)
		close(bCalcDone)

		// Tunggu A selesai WRITE
		<-aWriteDone

		// T5: B writes (overwrite hasil A)
		repo.SetStock(ctx, newStock)
		t.Logf("Request B: WRITE stock = %d", newStock)

		mu.Lock()
		successfulSales++
		mu.Unlock()
	}()

	wg.Wait()

	// Verify final state
	finalStock, _ := repo.GetStock(ctx)

	t.Logf("\n=== FINAL STATE ===")
	t.Logf("Initial stock:    %d", initialStock)
	t.Logf("Successful sales: %d", successfulSales)
	t.Logf("Final stock:      %d", finalStock)
	t.Logf("Expected:         1 = 2 + 0 → BROKEN")

	// Assertions
	if successfulSales != 2 {
		t.Errorf("expected successful_sales = 2, got %d", successfulSales)
	}
	if finalStock != 0 {
		t.Errorf("expected final_stock = 0, got %d", finalStock)
	}

	// Verify invariant is broken
	var expectedStock int
	mu.Lock()
	// Re-read final state
	finalStock, _ = repo.GetStock(ctx)
	expectedStock = successfulSales + finalStock
	mu.Unlock()

	if initialStock != expectedStock {
		t.Logf("✅ LOST UPDATE CONFIRMED: invariant %d != %d", initialStock, expectedStock)
	} else {
		t.Logf("❌ Invariant holds unexpectedly")
	}
}
