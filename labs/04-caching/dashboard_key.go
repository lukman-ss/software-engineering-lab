package caching

import (
	"fmt"
	"time"
)

// DashboardKeyBuilder creates standardized cache keys for dashboard statistics.
// Standard format: cmms:dashboard:v1:tenant:{tenantID}:branch:{branchID}:date:{YYYY-MM-DD}
type DashboardKeyBuilder struct {
	tenantID     int64
	branchID     int64
	businessDate time.Time
	version      int
}

// NewDashboardKey creates a new dashboard key builder with default v1.
func NewDashboardKey(branchID int64) *DashboardKeyBuilder {
	return &DashboardKeyBuilder{
		tenantID:     1, // default tenant
		branchID:     branchID,
		businessDate: time.Now().UTC(),
		version:      1,
	}
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

// WithDate sets the business date (timezone-aware).
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
