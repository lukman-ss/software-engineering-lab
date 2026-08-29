// Package caching explores real-world caching patterns.
package caching

import (
	"context"
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
	CallCount int64
	NextValue func() Dashboard // If set, returns increasingly higher invoice counts
}

func NewFakeDashboardRepository() *FakeDashboardRepository {
	return &FakeDashboardRepository{}
}

func (r *FakeDashboardRepository) GetDashboard(ctx context.Context, branchID int64, businessDate time.Time) (Dashboard, error) {
	r.CallCount++

	d := Dashboard{
		BranchID:          branchID,
		Date:              businessDate.Format("2006-01-02"),
		InvoiceCountToday: 42, // default value
		TotalRevenueToday: 150000.0,
	}

	if r.NextValue != nil {
		return r.NextValue(), nil
	}
	return d, nil
}

// GetDashboardWithTenant implements tenant-aware retrieval.
func (r *FakeDashboardRepository) GetDashboardWithTenant(ctx context.Context, tenantID, branchID int64, businessDate time.Time) (Dashboard, error) {
	r.CallCount++

	d := Dashboard{
		BranchID:          branchID,
		Date:              businessDate.Format("2006-01-02"),
		InvoiceCountToday: 42 + int(tenantID%10), // Different per tenant for demo
		TotalRevenueToday: 150000.0,
	}

	if r.NextValue != nil {
		return r.NextValue(), nil
	}
	return d, nil
}

func (r *FakeDashboardRepository) Reset() {
	r.CallCount = 0
}
