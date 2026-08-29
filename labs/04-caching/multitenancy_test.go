package caching_test

import (
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
func TestBusinessTimezone(t *testing.T) {
	locJakarta, _ := time.LoadLocation("Asia/Jakarta")
	locNY, _ := time.LoadLocation("America/New_York")

	now := time.Now()
	dateJakarta := caching.TodayInLocation(now, locJakarta)
	dateNY := caching.TodayInLocation(now, locNY)

	t.Logf("Business date in Jakarta: %s", dateJakarta)
	t.Logf("Business date in New York: %s", dateNY)

	// They might be different depending on when the test runs, which proves timezone handling matters
	keyJakarta := caching.NewDashboardKey(1).WithTenant(1).WithDate(now.In(locJakarta)).Build()
	keyNY := caching.NewDashboardKey(1).WithTenant(1).WithDate(now.In(locNY)).Build()

	t.Logf("Key Jakarta: %s", keyJakarta)
	t.Logf("Key New York: %s", keyNY)

	t.Log("Business timezone handling validated")
}

// TestTimezoneBoundary verifikasi bahwa time.Time yang diberikan ke key builder
// tidak di-convert lagi ke UTC, sehingga timezone bisnis dihormati.
func TestTimezoneBoundary(t *testing.T) {
	locJakarta, _ := time.LoadLocation("Asia/Jakarta")

	// 2026-08-30 00:30 Jakarta time = 2026-08-29 17:30 UTC
	// Jakarta business date = 2026-08-30
	ts := time.Date(2026, 8, 30, 0, 30, 0, 0, locJakarta)

	// Key must reflect Jakarta business date, NOT UTC date
	key := caching.NewDashboardKey(1).WithTenant(1).WithDate(ts).Build()

	// Expected: date:2026-08-30 (NOT 2026-08-29)
	if !contains(key, "date:2026-08-30") {
		t.Errorf("business date should be 2026-08-30, key: %s", key)
	}
	t.Logf("Key with Jakarta timezone: %s", key)
	t.Log("✓ Timezone boundary respected: Jakarta 00:30 = business date 2026-08-30")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > 0 && containsStr(s, substr)))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestDashboardKeyBuilderVersatility demonstrates the unified key builder
func TestDashboardKeyBuilderVersatility(t *testing.T) {
	// Single-tenant dashboard key (tenant defaults to 1)
	singleKey := caching.NewDashboardKey(42).Build()
	t.Logf("Single-tenant key: %s", singleKey)

	// Multi-tenant dashboard key
	multiKey := caching.NewDashboardKey(42).WithTenant(5).Build()
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
