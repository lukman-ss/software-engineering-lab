package caching

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

// DashboardCacheService menggunakan cache aside pattern.
type DashboardCacheService struct {
	repo         DashboardRepository
	cache        CacheInterface
	queryCounter atomic.Int64
}

func NewDashboardCacheService(repo DashboardRepository, cache CacheInterface) *DashboardCacheService {
	return &DashboardCacheService{
		repo:  repo,
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
		// Cache corrupt - delete and fetch fresh
		_ = s.cache.Delete(ctx, key)
	}

	s.queryCounter.Add(1)
	d, err := s.repo.GetDashboard(ctx, branchID, time.Now())
	if err != nil {
		return Dashboard{}, err
	}

	data, err := json.Marshal(d)
	if err != nil {
		return Dashboard{}, fmt.Errorf("marshal dashboard: %w", err)
	}
	if err := s.cache.Set(ctx, key, string(data), 30*time.Second); err != nil {
		// Log cache failure but don't fail the request
		fmt.Printf("warn: cache set failed: %v\n", err)
	}

	return d, nil
}

// InvalidateBranchDashboard menginvalidasi cache untuk branch tertentu.
// Preferred flow: COMMIT DB -> Invalidate Cache
func (s *DashboardCacheService) InvalidateBranchDashboard(ctx context.Context, branchID int64) error {
	key := dashboardCacheKey(branchID)
	err := s.cache.Delete(ctx, key)
	if err != nil && !errors.Is(err, ErrCacheMiss) {
		return fmt.Errorf("invalidate cache: %w", err)
	}
	return nil
}

// --- REPRESENTATIVE MUTATION METHODS (Bagian 8) ---

// CreateInvoice mensimulasikan pembuatan invoice baru (mutasi data).
// Flow yang benar: BEGIN -> DB update -> COMMIT -> Invalidate Cache
func (s *DashboardCacheService) CreateInvoice(ctx context.Context, branchID int64, amount float64) error {
	// In real implementation, this would use a transaction repository.
	// For demo, we just invalidate cache to simulate mutation occurred.
	_ = branchID
	_ = amount
	// In production: BEGIN -> INSERT invoice -> COMMIT -> Invalidate
	return s.InvalidateBranchDashboard(ctx, branchID)
}

// PayInvoice mensimulasikan pembayaran invoice.
func (s *DashboardCacheService) PayInvoice(ctx context.Context, branchID int64, invoiceID int64) error {
	// In real implementation, this would update the invoice
	_ = invoiceID
	// Invalidate cache setelah commit
	return s.InvalidateBranchDashboard(ctx, branchID)
}

// FinishService mensimulasikan penyelesaian servis oleh mekanik.
func (s *DashboardCacheService) FinishService(ctx context.Context, branchID int64, mechanicID int64) error {
	// In real implementation, this would update service_records
	_ = mechanicID
	return s.InvalidateBranchDashboard(ctx, branchID)
}

// UseSparepart mensimulasikan penggunaan sparepart.
func (s *DashboardCacheService) UseSparepart(ctx context.Context, branchID int64, partID int64) error {
	// In real implementation, this would update parts stock
	_ = partID
	return s.InvalidateBranchDashboard(ctx, branchID)
}

// SaveCustomer mensimulasikan pembuatan/perubahan customer.
func (s *DashboardCacheService) SaveCustomer(ctx context.Context, branchID int64, customerName string) error {
	// In real implementation, this would insert/update customer
	_ = customerName
	return s.InvalidateBranchDashboard(ctx, branchID)
}

// CreateVehicle mensimulasikan pembuatan kendaraan baru.
func (s *DashboardCacheService) CreateVehicle(ctx context.Context, branchID int64, plate string) error {
	// In real implementation, this would insert vehicle
	_ = plate
	return s.InvalidateBranchDashboard(ctx, branchID)
}

func (s *DashboardCacheService) QueryCount() int64 {
	return s.queryCounter.Load()
}