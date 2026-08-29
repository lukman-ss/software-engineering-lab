package caching

import (
	"fmt"
	"time"
)

// DashboardKeyBuilder creates standardized cache keys for dashboard statistics.
// Standard format: cmms:dashboard:v1:tenant:{tenantID}:branch:{branchID}:date:{YYYY-MM-DD}
type DashboardKeyBuilder struct {
	tenantID    int64
	branchID    int64
	businessDate time.Time
	version     int
}

// NewDashboardKey creates a new dashboard key builder with default v1.
func NewDashboardKey(branchID int64) *DashboardKeyBuilder {
	return &DashboardKeyBuilder{
		tenantID:   1, // default tenant
		branchID:   branchID,
		businessDate: time.Now().UTC(),
		version:    1,
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
	dateStr := b.businessDate.UTC().Format("2006-01-02")
	return fmt.Sprintf("cmms:dashboard:v%d:tenant:%d:branch:%d:date:%s",
		b.version, b.tenantID, b.branchID, dateStr)
}

// String implements Stringer for convenience.
func (b *DashboardKeyBuilder) String() string {
	return b.Build()
}

// TenantDashboardKey creates cache key for multi-tenant dashboard.
func TenantDashboardKey(tenantID, branchID int64, date time.Time) string {
	return NewDashboardKey(branchID).
		WithTenant(tenantID).
		WithDate(date).
		Build()
}

// KeyForDashboard membuat key sederhana untuk demo (single-tenant).
// Deprecated: Use NewDashboardKey(branchID).Build() for consistency.
func KeyForDashboard(branchID int64, date string) string {
	t, _ := time.Parse("2006-01-02", date)
	return NewDashboardKey(branchID).WithDate(t).Build()
}

// NewTenantDashboardKey creates a tenant-scoped dashboard key (legacy wrapper).
func NewTenantDashboardKey(tenantID, branchID int64, loc *time.Location) string {
	return TenantDashboardKey(tenantID, branchID, time.Now().In(loc))
}

// MultiTenantKeyExample demonstrasi key per tenant/branch
func MultiTenantKeyExample() {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	// Tenant A, Branch 1
	keyA1 := TenantDashboardKey(100, 1, time.Now().In(loc))
	// Tenant A, Branch 2
	keyA2 := TenantDashboardKey(100, 2, time.Now().In(loc))
	// Tenant B, Branch 1
	keyB1 := TenantDashboardKey(200, 1, time.Now().In(loc))

	fmt.Printf("Tenant A Branch 1: %s\n", keyA1)
	fmt.Printf("Tenant A Branch 2: %s\n", keyA2)
	fmt.Printf("Tenant B Branch 1: %s\n", keyB1)
}