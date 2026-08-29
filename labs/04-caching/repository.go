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
	// For multi-tenant deployments, use GetDashboardWithTenant which includes tenant scoping.
	GetDashboard(ctx context.Context, branchID int64, businessDate time.Time) (Dashboard, error)

	// GetDashboardWithTenant retrieves dashboard stats with explicit tenant scoping.
	// This is the recommended method for multi-tenant deployments.
	GetDashboardWithTenant(ctx context.Context, tenantID, branchID int64, businessDate time.Time) (Dashboard, error)
}

// FakeDashboardRepository is a test double that counts invocations.
type FakeDashboardRepository struct {
	callCount atomic.Int64
	mu        sync.RWMutex
	nextValue func() Dashboard
}

func NewFakeDashboardRepository() *FakeDashboardRepository {
	return &FakeDashboardRepository{}
}

// CallCount returns number of repository calls made.
func (r *FakeDashboardRepository) CallCount() int64 {
	return r.callCount.Load()
}

// SetNextValue configures the next value to return (thread-safe).
func (r *FakeDashboardRepository) SetNextValue(fn func() Dashboard) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextValue = fn
}

// Reset resets call count and next value.
func (r *FakeDashboardRepository) Reset() {
	r.callCount.Store(0)
	r.mu.Lock()
	r.nextValue = nil
	r.mu.Unlock()
}

func (r *FakeDashboardRepository) GetDashboard(ctx context.Context, branchID int64, businessDate time.Time) (Dashboard, error) {
	r.callCount.Add(1)

	d := Dashboard{
		BranchID:          branchID,
		Date:              businessDate.Format("2006-01-02"),
		InvoiceCountToday: 42, // default value
		TotalRevenueToday: 150000.0,
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.nextValue != nil {
		return r.nextValue(), nil
	}
	return d, nil
}

// GetDashboardWithTenant implements tenant-aware retrieval.
func (r *FakeDashboardRepository) GetDashboardWithTenant(ctx context.Context, tenantID, branchID int64, businessDate time.Time) (Dashboard, error) {
	r.callCount.Add(1)

	d := Dashboard{
		BranchID:          branchID,
		Date:              businessDate.Format("2006-01-02"),
		InvoiceCountToday: 42 + int(tenantID%10), // Different per tenant for demo
		TotalRevenueToday: 150000.0,
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.nextValue != nil {
		return r.nextValue(), nil
	}
	return d, nil
}
