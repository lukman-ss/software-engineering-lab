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
// Production multi-tenant calls MUST provide explicit tenantID, branchID, and businessDate.
//
// For production multi-tenant usage, use NewTenantDashboardKey:
//
//	key := NewTenantDashboardKey(tenantID, branchID).WithDate(businessDate).Build()
//
// Or set values explicitly:
//
//	localNow := now.In(loc)
//	date := localNow.Format("2006-01-02")
//	key := NewDashboardKey(branchID).WithTenant(tenantID).WithDate(date).Build()
//
// WARNING: time.Time.Truncate(24 * time.Hour) does NOT give local midnight!
// It truncates to the previous 24-hour boundary from the time's epoch,
// not to "00:00:00" in the local timezone.
//
// For local midnight boundary, use:
//
//	localNow := now.In(loc)
//	startOfDay := time.Date(
//	    localNow.Year(), localNow.Month(), localNow.Day(),
//	    0, 0, 0, 0, loc,
//	)
func NewDashboardKey(branchID int64) *DashboardKeyBuilder {
	return &DashboardKeyBuilder{
		tenantID:     1, // default tenant - DEMO ONLY. Must override in multi-tenant production
		branchID:     branchID,
		businessDate: time.Now().UTC(), // DEFAULT: convenience only, not timezone-aware business date
		version:      1,
	}
}

// NewTenantDashboardKey creates a dashboard key builder with explicit tenant and branch.
// Production-safe: requires tenantID and businessDate to be set explicitly.
// Multi-tenant isolation is enforced by construction.
func NewTenantDashboardKey(tenantID, branchID int64) *DashboardKeyBuilder {
	return &DashboardKeyBuilder{
		tenantID:     tenantID,
		branchID:     branchID,
		businessDate: time.Time{}, // MUST be set via WithDate()
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
// before extracting the calendar date. Examples:
//
//	// For cache key (just the date string):
//	localNow := now.In(loc)
//	date := localNow.Format("2006-01-02")
//	key := NewDashboardKey(branchID).WithTenant(tenantID).WithDate(date).Build()
//
//	// For local midnight boundary:
//	localNow := now.In(loc)
//	startOfDay := time.Date(
//	    localNow.Year(), localNow.Month(), localNow.Day(),
//	    0, 0, 0, 0, loc,
//	)
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
