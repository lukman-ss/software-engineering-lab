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

// NewDashboardKey creates a new dashboard key builder with defaults.
// WARNING: Default businessDate is today in UTC - this is for convenience only.
// Production calls MUST set explicit date via WithDate() with proper timezone handling.
// Business date should be: now.In(tenantTimezone).Truncate(24h)
func NewDashboardKey(branchID int64) *DashboardKeyBuilder {
	return &DashboardKeyBuilder{
		tenantID:     1, // default tenant - MUST override in multi-tenant
		branchID:     branchID,
		businessDate: time.Now().UTC(), // DEFAULT: convenience only, not timezone-aware business date
		version:      1,
	}
}

// WithTenant sets the tenant ID for multi-tenant isolation.
// REQUIRED: Must call this in multi-tenant deployments.
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
// CRITICAL: Caller must convert current instant to the tenant's business timezone
// before extracting the calendar date. Example:
//
//	businessDate := time.Now().In(loc).Truncate(24 * time.Hour)
//	key := NewDashboardKey(branchID).WithTenant(tenantID).WithDate(businessDate).Build()
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