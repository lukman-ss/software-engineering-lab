package caching

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"
)

// DashboardCacheService menggunakan cache aside pattern.
type DashboardCacheService struct {
	db           *sql.DB
	cache        CacheInterface
	queryCounter atomic.Int64
}

func NewDashboardCacheService(db *sql.DB, cache CacheInterface) *DashboardCacheService {
	return &DashboardCacheService{
		db:    db,
		cache: cache,
	}
}

func dashboardCacheKey(branchID int64) string {
	return fmt.Sprintf("cmms:dashboard:v1:branch:%d:date:%s", branchID, Today())
}

// GetDashboard mengembalikan statistik dashboard dengan caching.
func (s *DashboardCacheService) GetDashboard(ctx context.Context, branchID int64) (Dashboard, error) {
	key := dashboardCacheKey(branchID)

	cached, err := s.cache.Get(ctx, key)
	if err == nil && cached != "" {
		var d Dashboard
		if err := json.Unmarshal([]byte(cached), &d); err == nil {
			return d, nil
		}
	}

	s.queryCounter.Add(1)
	d, err := s.computeDashboard(ctx, branchID)
	if err != nil {
		return Dashboard{}, err
	}

	data, _ := json.Marshal(d)
	_ = s.cache.Set(ctx, key, string(data), 30*time.Second)

	return d, nil
}

// InvalidateBranchDashboard menginvalidasi cache untuk branch tertentu.
// Prefered flow: COMMIT DB -> Invalidate Cache
func (s *DashboardCacheService) InvalidateBranchDashboard(ctx context.Context, branchID int64) error {
	key := dashboardCacheKey(branchID)
	// Set value kosong atau hapus key di cache
	return s.cache.Set(ctx, key, "", -1*time.Second)
}

// --- REPRESENTATIVE MUTATION METHODS (Bagian 8) ---

// CreateInvoice mensimulasikan pembuatan invoice baru (mutasi data).
// Flow yang benar: BEGIN -> DB update -> COMMIT -> Invalidate Cache
func (s *DashboardCacheService) CreateInvoice(ctx context.Context, branchID int64, amount float64) error {
	// 1. BEGIN TRANSACTION (Disimulasikan)
	// 2. Insert ke database
	_, err := s.db.ExecContext(ctx, "INSERT INTO invoices (branch_id, amount, status) VALUES ($1, $2, 'unpaid')", branchID, amount)
	if err != nil {
		return fmt.Errorf("db insert invoice: %w", err)
	}
	// 3. COMMIT (Disimulasikan berhasil)

	// 4. INVALIDATE CACHE SETELAH COMMIT
	// Jangan delete sebelum commit untuk menghindari race condition (data belum commit tapi cache udah kosong)
	if err := s.InvalidateBranchDashboard(ctx, branchID); err != nil {
		// Log error invalidation, tapi jangan gagalkan transaction commit.
		// TTL akan menjadi safety net jika invalidation gagal.
		_ = err
	}

	return nil
}

// PayInvoice mensimulasikan pembayaran invoice.
func (s *DashboardCacheService) PayInvoice(ctx context.Context, branchID int64, invoiceID int64) error {
	_, err := s.db.ExecContext(ctx, "UPDATE invoices SET status = 'paid' WHERE id = $1", invoiceID)
	if err != nil {
		return fmt.Errorf("db pay invoice: %w", err)
	}

	// Invalidate cache setelah commit
	_ = s.InvalidateBranchDashboard(ctx, branchID)
	return nil
}

// FinishService mensimulasikan penyelesaian servis oleh mekanik.
func (s *DashboardCacheService) FinishService(ctx context.Context, branchID int64, mechanicID int64) error {
	_, err := s.db.ExecContext(ctx, "UPDATE service_records SET status = 'completed' WHERE mechanic_id = $1", mechanicID)
	if err != nil {
		return fmt.Errorf("db finish service: %w", err)
	}

	_ = s.InvalidateBranchDashboard(ctx, branchID)
	return nil
}

// UseSparepart mensimulasikan penggunaan sparepart.
func (s *DashboardCacheService) UseSparepart(ctx context.Context, branchID int64, partID int64) error {
	_, err := s.db.ExecContext(ctx, "UPDATE parts SET stock = stock - 1 WHERE id = $1", partID)
	if err != nil {
		return fmt.Errorf("db use sparepart: %w", err)
	}

	_ = s.InvalidateBranchDashboard(ctx, branchID)
	return nil
}

// SaveCustomer mensimulasikan pembuatan/perubahan customer.
func (s *DashboardCacheService) SaveCustomer(ctx context.Context, branchID int64, customerName string) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO customers (name, last_purchase) VALUES ($1, datetime('now'))", customerName)
	if err != nil {
		return fmt.Errorf("db save customer: %w", err)
	}

	_ = s.InvalidateBranchDashboard(ctx, branchID)
	return nil
}

// CreateVehicle mensimulasikan pembuatan kendaraan baru.
func (s *DashboardCacheService) CreateVehicle(ctx context.Context, branchID int64, plate string) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO vehicles (branch_id, plate) VALUES ($1, $2)", branchID, plate)
	if err != nil {
		return fmt.Errorf("db create vehicle: %w", err)
	}

	_ = s.InvalidateBranchDashboard(ctx, branchID)
	return nil
}

func (s *DashboardCacheService) computeDashboard(ctx context.Context, branchID int64) (Dashboard, error) {
	d := Dashboard{
		BranchID: branchID,
		Date:     Today(),
	}

	row := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*), COALESCE(SUM(amount), 0) FROM invoices WHERE branch_id = $1
	`, branchID)
	_ = row.Scan(&d.InvoiceCountToday, &d.TotalRevenueToday)

	return d, nil
}

func (s *DashboardCacheService) QueryCount() int64 {
	return s.queryCounter.Load()
}