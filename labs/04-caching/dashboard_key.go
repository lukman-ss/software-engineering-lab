package caching

import (
	"fmt"
	"time"
)

// DashboardKeyBuilder creates standardized cache keys for dashboard statistics.
// Standard format: cmms:dashboard:v1:tenant:{tenantID}:branch:{branchID}:date:{YYYY-MM-DD}
//
// IMPORTANT: Business date MUST come from tenant/branch timezone, not system UTC.
// Cache key date semantics must match the business day being reported.
// Using time.Now().UTC() as default is for convenience only - NOT production-ready.
type DashboardKeyBuilder struct {
	tenantID     int64
	branchID     int64
	businessDate time.Time
	version      int
}

// NewDashboardKey creates a new dashboard key builder for single-tenant/demo usage.
// WARNING: Default tenant is 1 and businessDate is today in UTC - this is for convenience only.
// Production multi-tenant calls MUST provide explicit tenantID, branchID, and businessDate via NewTenantDashboardKey.
func NewDashboardKey(branchID int64) *DashboardKeyBuilder {
	return &DashboardKeyBuilder{
		tenantID:     1, // default tenant - DEMO ONLY
		branchID:     branchID,
		businessDate: time.Now().UTC(), // DEFAULT: convenience only
		version:      1,
	}
}

// NewTenantDashboardKey creates a dashboard key builder with validated tenant, branch, and business date.
// Production-safe: validates tenantID > 0, branchID > 0, and non-zero businessDate at construction.
// Multi-tenant isolation and validity are enforced by construction.
func NewTenantDashboardKey(tenantID, branchID int64, businessDate time.Time) (*DashboardKeyBuilder, error) {
	if tenantID <= 0 {
		return nil, fmt.Errorf("tenant ID must be positive: got %d", tenantID)
	}
	if branchID <= 0 {
		return nil, fmt.Errorf("branch ID must be positive: got %d", branchID)
	}
	if businessDate.IsZero() {
		return nil, fmt.Errorf("business date cannot be zero")
	}
	return &DashboardKeyBuilder{
		tenantID:     tenantID,
		branchID:     branchID,
		businessDate: businessDate,
		version:      1,
	}, nil
}

// WithTenant sets the tenant ID for multi-tenant isolation.
func (b *DashboardKeyBuilder) WithTenant(tenantID int64) *DashboardKeyBuilder {
	b.tenantID = tenantID
	return b
}

// WithBranch sets the branch ID.
func (b *DashboardKeyBuilder) WithBranch(branchID int64) *DashboardKeyBuilder {
	b.branchID = branchID
	return b
}

// WithDate sets the business date from tenant/branch timezone.
func (b *DashboardKeyBuilder) WithDate(date time.Time) *DashboardKeyBuilder {
	b.businessDate = date
	return b
}

// WithVersion sets the key version for migrations.
func (b *DashboardKeyBuilder) WithVersion(v int) *DashboardKeyBuilder {
	b.version = v
	return b
}

// Build constructs the final cache key string.
func (b *DashboardKeyBuilder) Build() string {
	return fmt.Sprintf("cmms:dashboard:v%d:tenant:%d:branch:%d:date:%s",
		b.version, b.tenantID, b.branchID, b.businessDate.Format("2006-01-02"))
}

// String implements Stringer for convenience.
func (b *DashboardKeyBuilder) String() string {
	return b.Build()
}
