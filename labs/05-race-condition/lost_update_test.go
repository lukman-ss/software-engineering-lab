package race

import (
	"context"
	"sync"
	"testing"
)

// TestLostUpdate_Deterministic mereproduksi lost update secara deterministic.
//
// Timeline yang direproduksi:
// T0: Request A  READ stock = 1
// T1: Request B  READ stock = 1   (stale read — A belum WRITE)
// T2: Request A  CHECK > 0? YES, CALCULATE newStock = 0
// T3: Request B  CHECK > 0? YES, CALCULATE newStock = 0
// T4: Request A  WRITE stock = 0
// T5: Request B  WRITE stock = 0   (overwrite, tidak ada delta)
//
// Di akhir:
//   - final_stock = 0
//   - successful_sales = 2 (keduanya lewat CHECK)
//   - 1 != 2 + 0  →  invariant rusak
//
// Test ini tidak bergantung pada time.Sleep.
// Menggunakan channel barrier untuk sinkronisasi fase READ → CALCULATE → WRITE.
func TestLostUpdate_Deterministic(t *testing.T) {
	const productID = "Oli Mesin"
	const initialStock = 1

	repo := NewUnsafeInventoryRepository(map[string]int{productID: initialStock})
	ctx := context.Background()

	// Set stock via SetStock untuk inisialisasi yang eksplisit
	if err := repo.SetStock(ctx, productID, initialStock); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// svc demonstrates the service layer — not used in barrier test,
	// but documents TrySell logic. Test controls timing manually for determinism.
	svc := NewInventoryService(repo)
	_ = svc // service layer available for non-deterministic stress tests

	var wg sync.WaitGroup
	var mu sync.Mutex
	successfulSales := 0

	// Channel barrier untuk kontrol timing
	// Fase 1: kedua goroutine selesai READ
	// Fase 2: kedua goroutine selesai CALCULATE
	// Fase 3: kedua goroutine selesai WRITE
	aReadDone := make(chan struct{})
	bReadDone := make(chan struct{})
	aCalcDone := make(chan struct{})
	bCalcDone := make(chan struct{})
	aWriteDone := make(chan struct{})

	// Goroutine A
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Override TrySell dengan step yang sama tapi barrier
		stock, err := repo.GetStock(ctx, productID)
		if err != nil {
			t.Errorf("A GetStock: %v", err)
			return
		}
		t.Logf("Request A: READ stock = %d", stock)
		close(aReadDone)

		// Tunggu B selesai READ
		<-bReadDone

		// CHECK + CALCULATE
		if stock <= 0 {
			t.Error("A: stock should be > 0")
			return
		}
		newStock := stock - 1
		t.Logf("Request A: CHECK>0 YES, CALCULATE new_stock = %d", newStock)
		close(aCalcDone)

		// Tunggu B selesai CALCULATE
		<-bCalcDone

		// WRITE
		if err := repo.SetStock(ctx, productID, newStock); err != nil {
			t.Errorf("A SetStock: %v", err)
			return
		}
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

		stock, err := repo.GetStock(ctx, productID)
		if err != nil {
			t.Errorf("B GetStock: %v", err)
			return
		}
		t.Logf("Request B: READ stock = %d", stock)
		close(bReadDone)

		// Tunggu A selesai CALCULATE
		<-aCalcDone

		// CHECK + CALCULATE (baca stale value yang sama)
		if stock <= 0 {
			t.Error("B: stock should be > 0")
			return
		}
		newStock := stock - 1
		t.Logf("Request B: CHECK>0 YES, CALCULATE new_stock = %d", newStock)
		close(bCalcDone)

		// Tunggu A selesai WRITE
		<-aWriteDone

		// WRITE (overwrite hasil A)
		if err := repo.SetStock(ctx, productID, newStock); err != nil {
			t.Errorf("B SetStock: %v", err)
			return
		}
		t.Logf("Request B: WRITE stock = %d", newStock)

		mu.Lock()
		successfulSales++
		mu.Unlock()
	}()

	wg.Wait()

	// Verify final state
	finalStock, err := repo.GetStock(ctx, productID)
	if err != nil {
		t.Fatalf("final GetStock: %v", err)
	}

	t.Logf("\n=== FINAL STATE ===")
	t.Logf("Initial stock:    %d", initialStock)
	t.Logf("Successful sales: %d", successfulSales)
	t.Logf("Final stock:      %d", finalStock)
	t.Logf("Invariant check:  %d == %d + %d → BROKEN", initialStock, successfulSales, finalStock)

	// Assertions
	if successfulSales != 2 {
		t.Errorf("expected successful_sales = 2, got %d", successfulSales)
	}
	if finalStock != 0 {
		t.Errorf("expected final_stock = 0, got %d", finalStock)
	}

	// Verify invariant is broken — ini adalah expected outcome dari unsafe pattern
	var expectedStock int
	mu.Lock()
	finalStock, _ = repo.GetStock(ctx, productID)
	expectedStock = successfulSales + finalStock
	mu.Unlock()

	if initialStock != expectedStock {
		t.Logf("✅ LOST UPDATE CONFIRMED: invariant %d != %d", initialStock, expectedStock)
	} else {
		t.Logf("❌ Invariant holds unexpectedly — timing might need adjustment")
	}
}
