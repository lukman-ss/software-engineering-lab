# Execution Result

## Build
Command:
```bash
go build ./...
```
Result:
```text
Exit Code: 0 (Success)
```

## Tests
Command:
```bash
go test -v ./...
```
Result:
```text
?   	bloomfilters/cmd/demo	[no test files]
=== RUN   TestNoFalseNegatives
--- PASS: TestNoFalseNegatives (0.00s)
=== RUN   TestFalsePositiveRate
    bloom_test.go:49: Expected FP rate: 0.0100, Measured FP rate: 0.0098 (489 / 50000)
--- PASS: TestFalsePositiveRate (0.01s)
=== RUN   TestOptimalSizing
--- PASS: TestOptimalSizing (0.00s)
=== RUN   TestEmptyFilter
--- PASS: TestEmptyFilter (0.00s)
=== RUN   TestSyncFilterConcurrency
--- PASS: TestSyncFilterConcurrency (0.00s)
PASS
ok  	bloomfilters/internal/bloom	0.238s
=== RUN   TestLSMStoreWithAndWithoutFilter
    lsm_test.go:60: Absent queries: 5000
    lsm_test.go:61: Disk reads WITHOUT Bloom filter: 25000 (expected 25000)
    lsm_test.go:62: Disk reads WITH Bloom filter:    276
--- PASS: TestLSMStoreWithAndWithoutFilter (0.01s)
=== RUN   TestCachePenetrationMitigation
    lsm_test.go:125: 1000 penetration attempts: unprotected=1000, protected=0 (FP rate 0.00%)
--- PASS: TestCachePenetrationMitigation (0.00s)
PASS
ok  	bloomfilters/internal/store	0.349s
```

## Race Detector
Command:
```bash
go test -race ./...
```
Result:
```text
?   	bloomfilters/cmd/demo	[no test files]
ok  	bloomfilters/internal/bloom	1.487s
ok  	bloomfilters/internal/store	1.407s
```

## Demo
Command:
```bash
go run ./cmd/demo
```
Result:
```text
==================================================
LAB 39: BLOOM FILTERS IN STORAGE & CACHE SYSTEMS
==================================================

--- Scenario 1: Math Verification (n=10,000, eps=0.01) ---
Allocated Bits (m): 95851 (9.59 bits/element)
Hash Functions (k): 7
False Negatives across 10000 present keys: 0 (guarantee: 0)
False Positives across 50000 absent keys: 498 (0.996% vs target 1.00%)

--- Scenario 2: LSM-Tree Disk Read Reduction ---
Absent queries executed:     10000
Segment reads WITHOUT filter: 100000
Segment reads WITH filter:    988
Disk I/O Reduction:          99.01%

--- Scenario 3: Cache Penetration Prevention ---
Malicious / absent queries:      5000
Backend DB hits (unprotected):  5000 (100% penetration)
Backend DB hits (with Bloom):   0 (0.00% penetration)

==================================================
RESULT: All Bloom filter guarantees demonstrated.
==================================================
```

## Final Engineering Status
READY_FOR_ENGINEERING_AUDIT
