// Package race explores application-level race conditions in concurrent systems.
// This lab demonstrates how concurrent execution can corrupt business invariants,
// not just Go memory data races.
package race

import "context"

// Product mewakili produk inventory sederhana.
type Product struct {
	ID   string
	Name string
}

// InventoryItem mewakili item inventory dengan stock.
type InventoryItem struct {
	ProductID string
	Stock     int
}

// ErrOutOfStock dikembalikan ketika stock tidak mencukupi.
var ErrOutOfStock = &Error{Msg: "stock tidak mencukupi"}

// Error merepresentasikan error bisnis.
type Error struct {
	Msg string
}

func (e *Error) Error() string { return e.Msg }

// InventoryRepository mendefinisikan contract untuk operasi inventory.
type InventoryRepository interface {
	GetStock(ctx context.Context, productID string) (int, error)
	SetStock(ctx context.Context, productID string, stock int) error
}

// AtomicInventoryRepository mendefinisikan contract untuk operasi atomic conditional update.
// Menggunakan pola: UPDATE ... WHERE condition RETURNING rows_affected
type AtomicInventoryRepository interface {
	// DecrementStock mengurangi stock secara atomik dengan kondisi.
	// Mengembalikan (newStock, true) jika berhasil, (0, false) jika stock habis atau condition gagal.
	DecrementStock(ctx context.Context, productID string) (int, error)
}

// StockHistory merekam transaksi perubahan stock.
type StockHistory struct {
	ProductID   string
	OldStock    int
	NewStock    int
	Timestamp   int64
	Description string
}

// InventoryService mengelola logika business inventory.
type InventoryService struct {
	repo    InventoryRepository
	history []StockHistory
}

// NewInventoryService membuat service inventory baru.
func NewInventoryService(repo InventoryRepository) *InventoryService {
	return &InventoryService{repo: repo}
}

// TrySell mencoba melakukan penjualan dengan check-then-act pattern.
// Tujuan demo: menunjukkan race condition pada business state.
func (s *InventoryService) TrySell(ctx context.Context, productID string) error {
	// READ: Baca stock saat ini
	stock, err := s.repo.GetStock(ctx, productID)
	if err != nil {
		return err
	}

	// CHECK: Verifikasi cukup stok
	if stock <= 0 {
		return ErrOutOfStock
	}

	// CALCULATE: Hitung stok baru
	newStock := stock - 1

	// WRITE: Simpan stok baru
	return s.repo.SetStock(ctx, productID, newStock)
}
