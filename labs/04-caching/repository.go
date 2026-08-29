// Package caching explores real-world caching patterns and their trade-offs.
package caching

import (
	"context"
	"time"
)

// DashboardRepository defines the interface for dashboard data access.
type DashboardRepository interface {
	GetDashboard(ctx context.Context, branchID int64, businessDate time.Time) (Dashboard, error)
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

func (r *FakeDashboardRepository) Reset() {
	r.CallCount = 0
}