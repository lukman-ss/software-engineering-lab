# Execution Result

## Build
Command: `go mod tidy`
Result:
```text
go: downloading golang.org/x/sync v0.7.0
```

## Tests
Command: `go test -v ./...`
Result:
```text
?   	github.com/lukman/labs/37-cache-invalidation-strategies/cmd/demo	[no test files]
?   	github.com/lukman/labs/37-cache-invalidation-strategies/internal/cache	[no test files]
=== RUN   TestCachePatterns
=== RUN   TestCachePatterns/Cache-Aside_Read_&_Write
=== RUN   TestCachePatterns/Write-Through_Read_&_Write
=== RUN   TestCachePatterns/Write-Behind_Asynchronous_Flush
--- PASS: TestCachePatterns (0.05s)
    --- PASS: TestCachePatterns/Cache-Aside_Read_&_Write (0.00s)
    --- PASS: TestCachePatterns/Write-Through_Read_&_Write (0.00s)
    --- PASS: TestCachePatterns/Write-Behind_Asynchronous_Flush (0.05s)
=== RUN   TestStampedeMitigation
=== RUN   TestStampedeMitigation/Naive_Stampede_Queries_DB_Concurrently
=== RUN   TestStampedeMitigation/SingleFlight_Coalesces_To_Single_Query
--- PASS: TestStampedeMitigation (0.03s)
    --- PASS: TestStampedeMitigation/Naive_Stampede_Queries_DB_Concurrently (0.01s)
    --- PASS: TestStampedeMitigation/SingleFlight_Coalesces_To_Single_Query (0.02s)
=== RUN   TestXFetchLogic
--- PASS: TestXFetchLogic (0.00s)
=== RUN   TestStaleWhileRevalidate
--- PASS: TestStaleWhileRevalidate (0.09s)
=== RUN   TestJitter
--- PASS: TestJitter (0.00s)
PASS
ok  	github.com/lukman/labs/37-cache-invalidation-strategies/tests	0.693s
```

## Race Detector
Command: `go test -race ./...`
Result:
```text
?   	github.com/lukman/labs/37-cache-invalidation-strategies/cmd/demo	[no test files]
?   	github.com/lukman/labs/37-cache-invalidation-strategies/internal/cache	[no test files]
ok  	github.com/lukman/labs/37-cache-invalidation-strategies/tests	1.553s
```

## Demo
Command: `go run ./cmd/demo`
Result:
```text
========================================
Lab 37: Cache Invalidation Strategies
========================================

--- CACHE PATTERNS ---

[Cache-Aside]
  1st GET (cache miss)  -> DB queries: 1, val: Laptop
  2nd GET (cache hit)   -> DB queries: 1, val: Laptop
  After update & GET    -> DB queries: 2, val: Laptop Pro

[Write-Through]
  1st GET (cache miss)  -> DB queries: 1
  After update & GET    -> DB queries: 1, val: Phone Pro

[Write-Behind]
  After Update, GET from cache -> DB writes: 0, val: Tablet
  After 50ms flush delay       -> DB writes: 1

--- STAMPEDE DEMONSTRATION ---

[Naive] 20 concurrent goroutines -> 20 DB queries (stampede!)
[SF]    20 concurrent goroutines -> 1 DB queries (singleflight), 20 correct results

--- XFETCH (Probabilistic Early Expiration) ---
  Initial GET       -> DB queries: 1, val: hot-data
  Low rand draw     -> DB queries: 2, val: hot-data (early proactive refresh)
  High rand draw    -> DB queries: 1, val: fresh-data (served from cache, no early refresh)

--- STALE-WHILE-REVALIDATE ---
  Initial GET       -> val: v1
  GET (stale)       -> val: v1  (stale served immediately)
  GET (fresh)       -> val: v2  (async revalidation complete)

--- TTL JITTER ---
  Base TTL: 5m0s, Max Jitter: 30s
  Sampled TTLs: 5m13.339152199s, 5m26.959108616s, 5m18.809347525s, 5m28.542930723s, 5m12.311981236s

[Demo Complete]
```

## Final Engineering Status
ENGINEERING_STATUS: READY_FOR_ENGINEERING_AUDIT
