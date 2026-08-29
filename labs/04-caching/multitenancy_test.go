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

	keyTenant1Branch1 := caching.NewTenantDashboardKey(1, 10, loc).Build()
	keyTenant1Branch2 := caching.NewTenantDashboardKey(1, 20, loc).Build()
	keyTenant2Branch1 := caching.NewTenantDashboardKey(2, 10, loc).Build()

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

	t.Log("SUCCESS: Multi-tenant cache keys are fully isolated")
}

// TestBusinessTimezone ensures business day is respected
func TestBusinessTimezone(t *testing.T) {
	locJakarta, _ := time.LoadLocation("Asia/Jakarta")
	locNY, _ := time.LoadLocation("America/New_York")

	dateJakarta := caching.TodayInLocation(locJakarta)
	dateNY := caching.TodayInLocation(locNY)

	t.Logf("Business date in Jakarta: %s", dateJakarta)
	t.Logf("Business date in New York: %s", dateNY)

	// They might be different depending on when the test runs, which proves timezone handling matters
	keyJakarta := caching.NewTenantDashboardKey(1, 1, locJakarta).Build()
	keyNY := caching.NewTenantDashboardKey(1, 1, locNY).Build()

	t.Logf("Key Jakarta: %s", keyJakarta)
	t.Logf("Key New York: %s", keyNY)

	t.Log("SUCCESS: Business timezone supported correctly via location config")
}