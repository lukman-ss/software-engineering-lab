package caching_test

import (
	"strings"
	"testing"
	"time"

	caching "github.com/lukman-ss/software-engineering-lab/labs/04-caching"
)

// TestMultiTenancyKeyIsolation memastikan tenant/branch berbeda menghasilkan key berbeda.
func TestMultiTenancyKeyIsolation(t *testing.T) {
	now := time.Now()

	keyTenant1Branch1 := caching.DashboardCacheKey(1, 10, now)
	keyTenant1Branch2 := caching.DashboardCacheKey(1, 20, now)
	keyTenant2Branch1 := caching.DashboardCacheKey(2, 10, now)

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
// Uses a fixed UTC instant where Jakarta and New York have different calendar dates:
// 2026-08-29 18:30 UTC = 2026-08-30 01:30 Jakarta, 2026-08-29 14:30 New York
func TestBusinessTimezone(t *testing.T) {
	locJakarta, _ := time.LoadLocation("Asia/Jakarta")
	locNY, _ := time.LoadLocation("America/New_York")

	// Fixed instant: 2026-08-29 18:30 UTC
	utcNow := time.Date(2026, 8, 29, 18, 30, 0, 0, time.UTC)

	// Jakarta: 2026-08-30 01:30 (next day)
	// New York: 2026-08-29 14:30 (same day)
	dateJakarta := caching.TodayInLocation(utcNow, locJakarta)
	dateNY := caching.TodayInLocation(utcNow, locNY)

	// Assert exact expected dates
	expectedJakarta := "2026-08-30"
	expectedNY := "2026-08-29"

	if dateJakarta != expectedJakarta {
		t.Errorf("Jakarta date: expected %s, got %s", expectedJakarta, dateJakarta)
	}
	if dateNY != expectedNY {
		t.Errorf("New York date: expected %s, got %s", expectedNY, dateNY)
	}

	// Keys must reflect business timezone
	bJ, _ := caching.NewTenantDashboardKey(1, 1, utcNow.In(locJakarta))
	keyJakarta := bJ.Build()
	bNY, _ := caching.NewTenantDashboardKey(1, 1, utcNow.In(locNY))
	keyNY := bNY.Build()

	// Assert exact expected key dates
	if !strings.Contains(keyJakarta, "date:2026-08-30") {
		t.Errorf("Jakarta key should contain date:2026-08-30, got %s", keyJakarta)
	}
	if !strings.Contains(keyNY, "date:2026-08-29") {
		t.Errorf("New York key should contain date:2026-08-29, got %s", keyNY)
	}

	t.Logf("Jakarta key: %s", keyJakarta)
	t.Logf("New York key: %s", keyNY)
	t.Log("✓ Business timezone: Jakarta = 2026-08-30, New York = 2026-08-29")
}

// TestTimezoneBoundary verifikasi bahwa time.Time yang diberikan ke key builder
// tidak di-convert lagi ke UTC, sehingga timezone bisnis dihormati.
func TestTimezoneBoundary(t *testing.T) {
	locJakarta, _ := time.LoadLocation("Asia/Jakarta")

	// 2026-08-30 00:30 Jakarta time = 2026-08-29 17:30 UTC
	// Jakarta business date = 2026-08-30
	ts := time.Date(2026, 8, 30, 0, 30, 0, 0, locJakarta)

	// Key must reflect Jakarta business date, NOT UTC date
	b, _ := caching.NewTenantDashboardKey(1, 1, ts)
	key := b.Build()

	// Expected: date:2026-08-30 (NOT 2026-08-29)
	if !strings.Contains(key, "date:2026-08-30") {
		t.Errorf("business date should be 2026-08-30, key: %s", key)
	}
	t.Logf("Key with Jakarta timezone: %s", key)
	t.Log("✓ Timezone boundary respected: Jakarta 00:30 = business date 2026-08-30")
}

// TestDashboardKeyBuilderVersatility demonstrates the unified key builder
func TestDashboardKeyBuilderVersatility(t *testing.T) {
	// Single-tenant dashboard key (tenant defaults to 1)
	singleKey := caching.NewDashboardKey(42).Build()
	t.Logf("Single-tenant key: %s", singleKey)

	// Multi-tenant dashboard key
	bM, _ := caching.NewTenantDashboardKey(5, 42, time.Now().UTC())
	multiKey := bM.Build()
	t.Logf("Multi-tenant key: %s", multiKey)

	// With explicit date (use a fixed past date to differentiate from today)
	specificDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
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

// TestNewTenantDashboardKeyValidation verifies zero tenant, zero branch, zero businessDate return error.
func TestNewTenantDashboardKeyValidation(t *testing.T) {
	validDate := time.Date(2026, 9, 3, 0, 0, 0, 0, time.UTC)

	// Zero tenant -> error
	_, err := caching.NewTenantDashboardKey(0, 1, validDate)
	if err == nil {
		t.Error("expected error for zero tenant ID, got nil")
	}

	// Negative tenant -> error
	_, err = caching.NewTenantDashboardKey(-1, 1, validDate)
	if err == nil {
		t.Error("expected error for negative tenant ID, got nil")
	}

	// Zero branch -> error
	_, err = caching.NewTenantDashboardKey(1, 0, validDate)
	if err == nil {
		t.Error("expected error for zero branch ID, got nil")
	}

	// Zero businessDate -> error
	_, err = caching.NewTenantDashboardKey(1, 1, time.Time{})
	if err == nil {
		t.Error("expected error for zero businessDate, got nil")
	}

	// Valid values -> correct key
	builder, err := caching.NewTenantDashboardKey(1, 42, validDate)
	if err != nil {
		t.Fatalf("unexpected error for valid values: %v", err)
	}

	expectedKey := "cmms:dashboard:v1:tenant:1:branch:42:date:2026-09-03"
	if key := builder.Build(); key != expectedKey {
		t.Errorf("expected key %s, got %s", expectedKey, key)
	}

	t.Log("✓ NewTenantDashboardKey validation validated")
}
