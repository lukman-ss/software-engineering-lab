package store

import (
	"fmt"
	"testing"
)

func TestLSMStoreWithAndWithoutFilter(t *testing.T) {
	const numSegments = 5
	const keysPerSegment = 1000

	// Build 5 segments of data
	var segmentsData []map[string]string
	for segIdx := 0; segIdx < numSegments; segIdx++ {
		data := make(map[string]string, keysPerSegment)
		for k := 0; k < keysPerSegment; k++ {
			key := fmt.Sprintf("seg%d-key%d", segIdx, k)
			data[key] = fmt.Sprintf("val-%d-%d", segIdx, k)
		}
		segmentsData = append(segmentsData, data)
	}

	// 1. LSM Store WITHOUT filters
	lsmNoFilter := NewLSMStore()
	for _, data := range segmentsData {
		lsmNoFilter.AddSegment(NewSegment(data, false, 0))
	}

	// 2. LSM Store WITH Bloom filters (eps = 0.01)
	lsmWithFilter := NewLSMStore()
	for _, data := range segmentsData {
		lsmWithFilter.AddSegment(NewSegment(data, true, 0.01))
	}

	// Query for present keys: both should find them
	for segIdx := 0; segIdx < numSegments; segIdx++ {
		key := fmt.Sprintf("seg%d-key500", segIdx)
		if _, ok := lsmNoFilter.Get(key); !ok {
			t.Fatalf("lsmNoFilter missed existing key %s", key)
		}
		if _, ok := lsmWithFilter.Get(key); !ok {
			t.Fatalf("lsmWithFilter missed existing key %s (false negative!)", key)
		}
	}

	lsmNoFilter.ResetCounters()
	lsmWithFilter.ResetCounters()

	// Query for 5000 ABSENT keys
	const absentQueries = 5000
	for i := 0; i < absentQueries; i++ {
		absentKey := fmt.Sprintf("absent-key-%d", i)
		lsmNoFilter.Get(absentKey)
		lsmWithFilter.Get(absentKey)
	}

	readsNoFilter := lsmNoFilter.DiskReads()
	readsWithFilter := lsmWithFilter.DiskReads()

	t.Logf("Absent queries: %d", absentQueries)
	t.Logf("Disk reads WITHOUT Bloom filter: %d (expected %d)", readsNoFilter, absentQueries*numSegments)
	t.Logf("Disk reads WITH Bloom filter:    %d", readsWithFilter)

	if readsNoFilter != uint64(absentQueries*numSegments) {
		t.Errorf("Expected %d disk reads without filter, got %d", absentQueries*numSegments, readsNoFilter)
	}

	// With Bloom filter (eps=0.01 across 5 segments), reads should be < 5% of unfiltered reads
	if float64(readsWithFilter) > float64(readsNoFilter)*0.05 {
		t.Errorf("Bloom filter disk reads (%d) exceeded 5%% of unfiltered reads (%d)", readsWithFilter, readsNoFilter)
	}
}

type mockBackend struct {
	data map[string]string
}

func (m *mockBackend) Fetch(key string) (string, bool) {
	v, ok := m.data[key]
	return v, ok
}

func TestCachePenetrationMitigation(t *testing.T) {
	backendData := map[string]string{
		"user:101": "Alice",
		"user:102": "Bob",
		"user:103": "Charlie",
	}
	backend := &mockBackend{data: backendData}

	// 1. Unprotected cache (no filter): every absent query reaches backend.
	cacheUnprotected := NewCache(0, 0, backend)
	for k, v := range backendData {
		cacheUnprotected.Add(k, v) // populate in-memory only
	}
	const attacks = 1000
	for i := 0; i < attacks; i++ {
		cacheUnprotected.Get(fmt.Sprintf("hacker-query-%d", i))
	}
	if cacheUnprotected.BackendCalls() != attacks {
		t.Errorf("Unprotected cache expected %d backend calls, got %d", attacks, cacheUnprotected.BackendCalls())
	}

	// 2. Protected cache: Bloom filter pre-populated with known keys.
	// Only keys already known to the filter pass through; unknown absent keys are blocked.
	cacheProtected := NewCache(1000, 0.01, backend)
	for k, v := range backendData {
		cacheProtected.Add(k, v) // populates both in-mem map and filter
	}
	// Sanity: fetch a known key – should come from in-mem map, no backend call.
	val, ok := cacheProtected.Get("user:101")
	if !ok || val != "Alice" {
		t.Fatalf("Expected Alice from cache, got %q ok=%v", val, ok)
	}
	if cacheProtected.BackendCalls() != 0 {
		t.Errorf("Expected 0 backend calls for cached key, got %d", cacheProtected.BackendCalls())
	}

	// Cache-penetration attack: 1000 requests for non-existent keys.
	for i := 0; i < attacks; i++ {
		cacheProtected.Get(fmt.Sprintf("hacker-query-%d", i))
	}

	penetrations := cacheProtected.BackendCalls()
	t.Logf("%d penetration attempts: unprotected=%d, protected=%d (FP rate %.2f%%)",
		attacks, cacheUnprotected.BackendCalls(), penetrations, float64(penetrations)/float64(attacks)*100)

	// Bloom filter (eps=0.01) should block ~99% → at most ~30 backend calls.
	if penetrations > 30 {
		t.Errorf("Too many cache penetrations (%d > 30 threshold)", penetrations)
	}
	if penetrations >= cacheUnprotected.BackendCalls() {
		t.Errorf("Protected cache (%d) had >= backend calls than unprotected (%d)", penetrations, cacheUnprotected.BackendCalls())
	}
}
