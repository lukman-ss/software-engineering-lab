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

// GetDashboard mengembalikan statistik dashboard dengan caching (single-tenant).
func (s *DashboardCacheService) GetDashboard(ctx context.Context, branchID int64) (Dashboard, error) {
	return s.GetDashboardWithTenant(ctx, 1, branchID, time.Now().UTC())
}

// GetDashboardWithTenant mengembalikan statistik dashboard dengan tenant isolation.
// Primary entry point for multi-tenant deployments.
func (s *DashboardCacheService) GetDashboardWithTenant(ctx context.Context, tenantID, branchID int64, businessDate time.Time) (Dashboard, error) {
	key := NewDashboardKey(branchID).
		WithTenant(tenantID).
		WithDate(businessDate).
		Build()

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
	d, err := s.repo.GetDashboard(ctx, tenantID, branchID, businessDate)
	if err != nil {
		return Dashboard{}, err
	}

	data, err := json.Marshal(d)
	if err != nil {
		return Dashboard{}, fmt.Errorf("marshal dashboard: %w", err)
	}

	// Add jitter to prevent synchronized cache expiration across branches
	jitteredTTL := TTLWithJitter(30*time.Second, 10*time.Second)
	if err := s.cache.Set(ctx, key, string(data), jitteredTTL); err != nil {
		// Log cache failure but don't fail the request
		fmt.Printf("warn: cache set failed: %v\n", err)
	}

	return d, nil
}

// InvalidateDashboard invalidates cache for a branch and specific business date.
func (s *DashboardCacheService) InvalidateDashboard(ctx context.Context, tenantID, branchID int64, businessDate time.Time) error {
	key := NewDashboardKey(branchID).
		WithTenant(tenantID).
		WithDate(businessDate).
		Build()
	err := s.cache.Delete(ctx, key)
	if err != nil && !errors.Is(err, ErrCacheMiss) {
		return fmt.Errorf("invalidate cache: %w", err)
	}
	return nil
}

// InvalidateCurrentDashboard invalidates cache for today using the injected clock.
func (s *DashboardCacheService) InvalidateCurrentDashboard(ctx context.Context, tenantID, branchID int64) error {
	return s.InvalidateDashboard(ctx, tenantID, branchID, defaultClock.Now().UTC())
}

// --- REPRESENTATIVE MUTATION METHODS (Bagian 8) ---

// CreateInvoice mensimulasikan pembuatan invoice baru (mutasi data).
func (s *DashboardCacheService) CreateInvoice(ctx context.Context, branchID int64, amount float64) error {
	_ = amount
	return s.InvalidateCurrentDashboard(ctx, 1, branchID)
}

// PayInvoice mensimulasikan pembayaran invoice.
func (s *DashboardCacheService) PayInvoice(ctx context.Context, branchID int64, invoiceID int64) error {
	_ = invoiceID
	return s.InvalidateCurrentDashboard(ctx, 1, branchID)
}

// FinishService mensimulasikan penyelesaian servis oleh mekanik.
func (s *DashboardCacheService) FinishService(ctx context.Context, branchID int64, mechanicID int64) error {
	_ = mechanicID
	return s.InvalidateCurrentDashboard(ctx, 1, branchID)
}

// UseSparepart mensimulasikan penggunaan sparepart.
func (s *DashboardCacheService) UseSparepart(ctx context.Context, branchID int64, partID int64) error {
	_ = partID
	return s.InvalidateCurrentDashboard(ctx, 1, branchID)
}

// SaveCustomer mensimulasikan pembuatan/perubahan customer.
func (s *DashboardCacheService) SaveCustomer(ctx context.Context, branchID int64, customerName string) error {
	_ = customerName
	return s.InvalidateCurrentDashboard(ctx, 1, branchID)
}

// CreateVehicle mensimulasikan pembuatan kendaraan baru.
func (s *DashboardCacheService) CreateVehicle(ctx context.Context, branchID int64, plate string) error {
	_ = plate
	return s.InvalidateCurrentDashboard(ctx, 1, branchID)
}

func (s *DashboardCacheService) QueryCount() int64 {
	return s.queryCounter.Load()
}
