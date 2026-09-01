// Package caching explores real-world caching patterns.
package caching

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// DashboardRepository defines the interface for dashboard data access.
// Supports tenant isolation for multi-tenant deployments.
type DashboardRepository interface {
	// GetDashboard retrieves dashboard stats for a branch on a specific business date.
	// This canonical method explicitly scopes by tenantID.
	GetDashboard(ctx context.Context, tenantID, branchID int64, businessDate time.Time) (Dashboard, error)
}

// FakeDashboardRepository is a test double that counts invocations.
type FakeDashboardRepository struct {
	callCount     atomic.Int64
	mu            sync.RWMutex
	tenantDataMap map[int64]Dashboard
	nextValue     func() Dashboard
}

func NewFakeDashboardRepository() *FakeDashboardRepository {
	return &FakeDashboardRepository{
		tenantDataMap: make(map[int64]Dashboard),
	}
}

// CallCount returns number of repository calls made.
func (r *FakeDashboardRepository) CallCount() int64 {
	return r.callCount.Load()
}

// SetTenantData sets explicit data for a tenant (useful for behavioral isolation testing).
func (r *FakeDashboardRepository) SetTenantData(tenantID int64, data Dashboard) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tenantDataMap[tenantID] = data
}

// SetNextValue configures the next value to return globally (thread-safe).
func (r *FakeDashboardRepository) SetNextValue(fn func() Dashboard) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextValue = fn
}

// Reset resets call count and mocked data.
func (r *FakeDashboardRepository) Reset() {
	r.callCount.Store(0)
	r.mu.Lock()
	r.tenantDataMap = make(map[int64]Dashboard)
	r.nextValue = nil
	r.mu.Unlock()
}

// GetDashboard implements tenant-aware retrieval.
func (r *FakeDashboardRepository) GetDashboard(ctx context.Context, tenantID, branchID int64, businessDate time.Time) (Dashboard, error) {
	r.callCount.Add(1)

	r.mu.RLock()
	defer r.mu.RUnlock()

	// Prioritize global override if set
	if r.nextValue != nil {
		return r.nextValue(), nil
	}

	// Then check tenant-specific configured data
	if d, ok := r.tenantDataMap[tenantID]; ok {
		// Override date to match requested date for correctness
		d.Date = businessDate.Format("2006-01-02")
		return d, nil
	}

	// Fallback to default demo data
	return Dashboard{
		BranchID:          branchID,
		Date:              businessDate.Format("2006-01-02"),
		InvoiceCountToday: 42 + int(tenantID%10), // Different per tenant for demo
		TotalRevenueToday: 150000.0,
	}, nil
}
