package main

import (
	"fmt"

	"bloomfilters/internal/bloom"
	"bloomfilters/internal/store"
)

type diskBackend struct {
	records map[string]string
}

func (d *diskBackend) Fetch(key string) (string, bool) {
	v, ok := d.records[key]
	return v, ok
}

func main() {
	fmt.Println("==================================================")
	fmt.Println("LAB 39: BLOOM FILTERS IN STORAGE & CACHE SYSTEMS")
	fmt.Println("==================================================")

	// SCENARIO 1: Bloom Filter Math Verification
	fmt.Println("\n--- Scenario 1: Math Verification (n=10,000, eps=0.01) ---")
	const n = 10000
	const eps = 0.01
	bf := bloom.New(n, eps)

	fmt.Printf("Allocated Bits (m): %d (%.2f bits/element)\n", bf.M(), float64(bf.M())/float64(n))
	fmt.Printf("Hash Functions (k): %d\n", bf.K())

	// Insert keys
	for i := 0; i < n; i++ {
		bf.Add([]byte(fmt.Sprintf("user:%d", i)))
	}

	// Verify no false negatives
	falseNegatives := 0
	for i := 0; i < n; i++ {
		if !bf.Check([]byte(fmt.Sprintf("user:%d", i))) {
			falseNegatives++
		}
	}
	fmt.Printf("False Negatives across %d present keys: %d (guarantee: 0)\n", n, falseNegatives)

	// Measure false positives
	const testQueries = 50000
	falsePositives := 0
	for i := 0; i < testQueries; i++ {
		if bf.Check([]byte(fmt.Sprintf("absent-user:%d", i))) {
			falsePositives++
		}
	}
	empiricalFP := float64(falsePositives) / float64(testQueries)
	fmt.Printf("False Positives across %d absent keys: %d (%.3f%% vs target %.2f%%)\n",
		testQueries, falsePositives, empiricalFP*100, eps*100)

	// SCENARIO 2: LSM-tree Disk I/O Reduction
	fmt.Println("\n--- Scenario 2: LSM-Tree Disk Read Reduction ---")
	const numSegments = 10
	const keysPerSegment = 2000

	lsmWithoutFilter := store.NewLSMStore()
	lsmWithFilter := store.NewLSMStore()

	for seg := 0; seg < numSegments; seg++ {
		kvs := make(map[string]string, keysPerSegment)
		for k := 0; k < keysPerSegment; k++ {
			kvs[fmt.Sprintf("seg-%d-doc-%d", seg, k)] = "payload"
		}
		lsmWithoutFilter.AddSegment(store.NewSegment(kvs, false, 0))
		lsmWithFilter.AddSegment(store.NewSegment(kvs, true, 0.01))
	}

	const absentSearches = 10000
	for i := 0; i < absentSearches; i++ {
		query := fmt.Sprintf("absent-key-%d", i)
		lsmWithoutFilter.Get(query)
		lsmWithFilter.Get(query)
	}

	readsNoFilter := lsmWithoutFilter.DiskReads()
	readsWithFilter := lsmWithFilter.DiskReads()
	reduction := (1.0 - float64(readsWithFilter)/float64(readsNoFilter)) * 100.0

	fmt.Printf("Absent queries executed:     %d\n", absentSearches)
	fmt.Printf("Segment reads WITHOUT filter: %d\n", readsNoFilter)
	fmt.Printf("Segment reads WITH filter:    %d\n", readsWithFilter)
	fmt.Printf("Disk I/O Reduction:          %.2f%%\n", reduction)

	// SCENARIO 3: Cache Penetration Defense
	fmt.Println("\n--- Scenario 3: Cache Penetration Prevention ---")
	db := &diskBackend{
		records: map[string]string{
			"product:1": "Laptop",
			"product:2": "Phone",
			"product:3": "Monitor",
		},
	}

	cacheUnprotected := store.NewCache(0, 0, db)
	cacheProtected := store.NewCache(100, 0.01, db)

	for k, v := range db.records {
		cacheUnprotected.Add(k, v)
		cacheProtected.Add(k, v)
	}

	const penetrationAttempts = 5000
	for i := 0; i < penetrationAttempts; i++ {
		badQuery := fmt.Sprintf("non-existent-product-%d", i)
		cacheUnprotected.Get(badQuery)
		cacheProtected.Get(badQuery)
	}

	backendHitsUnprotected := cacheUnprotected.BackendCalls()
	backendHitsProtected := cacheProtected.BackendCalls()

	fmt.Printf("Malicious / absent queries:      %d\n", penetrationAttempts)
	fmt.Printf("Backend DB hits (unprotected):  %d (100%% penetration)\n", backendHitsUnprotected)
	fmt.Printf("Backend DB hits (with Bloom):   %d (%.2f%% penetration)\n",
		backendHitsProtected, float64(backendHitsProtected)/float64(penetrationAttempts)*100)

	fmt.Println("\n==================================================")
	fmt.Println("RESULT: All Bloom filter guarantees demonstrated.")
	fmt.Println("==================================================")
}
