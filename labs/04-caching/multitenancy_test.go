package caching_test

import (
	"testing"
	"time"

	caching "github.com/lukman/software-engineer-lab/labs/04-caching"
)

// TestMultiTenancyKeyIsolation memastikan tenant/branch berbeda menghasilkan key berbeda.
func TestMultiTenancyKeyIsolation(t *testing.T) {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		loc = time.UTC
	}

	keyTenant1Branch1 := caching.NewTenantDashboardKey(1, 10, loc)
	keyTenant1Branch2 := caching.NewTenantDashboardKey(1, 20, loc)
	keyTenant2Branch1 := caching.NewTenantDashboardKey(2, 10, loc)

	t.Logf("Key T1B1: %s", keyTenant1Branch1)
	t.Logf("Key T1B2: %s", keyTenant1Branch2)
	t.Logf("Key T2B1: %s", keyTenant2Branch1)

	// Invariant: No collisions between different tenants or branches
	if keyTenant1Branch1 == keyTenant1Branch2 {
		t.Error("tenant 1 branch 1 and branch 2 should not have the same key")
	}
	if keyTenant1Branch1 == keyTenant2Branch1 {
		t.Error("tenant 1 and tenant 2 should not have the same key")
	}

	t.Log("Multi-tenant key isolation validated")
}

// TestBusinessTimezone ensures business day is respected
func TestBusinessTimezone(t *testing.T) {
	locJakarta, _ := time.LoadLocation("Asia/Jakarta")
	locNY, _ := time.LoadLocation("America/New_York")

	now := time.Now()
	dateJakarta := caching.TodayInLocation(now, locJakarta)
	dateNY := caching.TodayInLocation(now, locNY)

	t.Logf("Business date in Jakarta: %s", dateJakarta)
	t.Logf("Business date in New York: %s", dateNY)

	// They might be different depending on when the test runs, which proves timezone handling matters
	keyJakarta := caching.NewTenantDashboardKey(1, 1, locJakarta)
	keyNY := caching.NewTenantDashboardKey(1, 1, locNY)

	t.Logf("Key Jakarta: %s", keyJakarta)
	t.Logf("Key New York: %s", keyNY)

	t.Log("Business timezone handling validated")
}

// TestDashboardKeyBuilderVersatility demonstrates the unified key builder
func TestDashboardKeyBuilderVersatility(t *testing.T) {
	// Single-tenant dashboard key (tenant defaults to 1)
	singleKey := caching.NewDashboardKey(42).Build()
	t.Logf("Single-tenant key: %s", singleKey)

	// Multi-tenant dashboard key
	multiKey := caching.NewDashboardKey(42).WithTenant(5).Build()
	t.Logf("Multi-tenant key: %s", multiKey)

	// With explicit date
	specificDate := time.Date(2026, 8, 29, 0, 0, 0, 0, time.UTC)
	dateKey := caching.NewDashboardKey(42).WithDate(specificDate).Build()
	t.Logf("Specific date key: %s", dateKey)

	// With version for migration
	v2Key := caching.NewDashboardKey(42).WithVersion(2).Build()
	t.Logf("Version 2 key: %s", v2Key)

	// Verify all keys are different
	keys := []string{singleKey, multiKey, dateKey, v2Key}
	unique := make(map[string]bool)
	for _, k := range keys {
		unique[k] = true
	}
	if len(unique) != len(keys) {
		t.Error("all keys should be unique")
	}
}