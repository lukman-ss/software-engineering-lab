package caching

import (
	"fmt"
	"time"
)

// TenantKey adalah cache key builder yang mempertimbangkan multi-tenancy.
// Format: cmms:tenant:{tenantID}:branch:{branchID}:entity:v1:{date}
type TenantKey struct {
	Namespace string
	TenantID  int64
	BranchID  int64
	Entity    string
	Version   int
	Date      string
}

// Build membuat cache key dengan tenant isolation.
func (k TenantKey) Build() string {
	return fmt.Sprintf("%s:tenant:%d:branch:%d:%s:v%d:%s",
		k.Namespace,
		k.TenantID,
		k.BranchID,
		k.Entity,
		k.Version,
		k.Date,
	)
}

// NewTenantDashboardKey membuat key untuk dashboard dengan tenant isolation
func NewTenantDashboardKey(tenantID, branchID int64, loc *time.Location) TenantKey {
	return TenantKey{
		Namespace: "cmms",
		TenantID:  tenantID,
		BranchID:  branchID,
		Entity:    "dashboard",
		Version:   1,
		Date:      TodayInLocation(defaultClock.Now(), loc),
	}
}

// KeyForDashboard membuat key sederhana untuk demo (single-tenant).
func KeyForDashboard(branchID int64, date string) string {
	return fmt.Sprintf("cmms:dashboard:v1:branch:%d:date:%s", branchID, date)
}

// MultiTenantKeyExample demonstrasi key per tenant/branch
func MultiTenantKeyExample() {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	// Tenant A, Branch 1
	keyA1 := NewTenantDashboardKey(100, 1, loc).Build()
	// Tenant A, Branch 2
	keyA2 := NewTenantDashboardKey(100, 2, loc).Build()
	// Tenant B, Branch 1
	keyB1 := NewTenantDashboardKey(200, 1, loc).Build()

	fmt.Printf("Tenant A Branch 1: %s\n", keyA1)
	fmt.Printf("Tenant A Branch 2: %s\n", keyA2)
	fmt.Printf("Tenant B Branch 1: %s\n", keyB1)
}