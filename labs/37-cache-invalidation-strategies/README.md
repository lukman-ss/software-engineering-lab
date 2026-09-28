# Lab 37: Cache Invalidation Strategies

Demonstrates cache invalidation patterns (Cache-Aside, Write-Through, Write-Behind) and stampede mitigations (Singleflight, XFetch, Stale-While-Revalidate, TTL Jitter) in Go.

## Code Structure

```text
37-cache-invalidation-strategies/
├── cmd/
│   └── demo/main.go           # Runnable demonstration
├── internal/
│   └── cache/
│       ├── store.go           # Memory cache with TTL & jitter
│       ├── repo.go            # Mock DB with atomic query counters
│       ├── patterns.go        # Cache-Aside, Write-Through, Write-Behind
│       └── stampede.go        # Singleflight, XFetch, SWR implementations
├── tests/
│   └── cache_test.go          # Unit, integration, & concurrency tests
├── engineering/
│   ├── 01-design.md           # Architecture & test strategy
│   ├── 02-implementation-notes.md  # Core decisions & limitations
│   └── 03-execution-result.md # Build & test outputs
└── go.mod
```

## Running the Lab

### Run Tests
```bash
go test -v ./...
```

### Run Race Detector
```bash
go test -race ./...
```

### Run Demo
```bash
go run ./cmd/demo
```

## Implemented Features
- **Cache-Aside**: Checks cache, loads from DB on miss, populates cache. Invalidates cache on update.
- **Write-Through**: Updates DB and cache synchronously on write; subsequent reads hit cache.
- **Write-Behind**: Updates cache immediately; flushes to DB asynchronously via background worker queue.
- **SingleFlight**: Coalesces concurrent cache miss requests on a hot key down to a single DB query using `golang.org/x/sync/singleflight`.
- **XFetch (Probabilistic Expiration)**: Early refresh using formula `-Δ · β · ln(U) > TTL_remaining` where `U ~ Uniform(0,1)`.
- **Stale-While-Revalidate (SWR)**: Returns stale cached data immediately while asynchronously triggering background DB revalidation.
- **TTL Jitter**: Adds randomized offset `[0, maxJitter)` to base TTL to prevent synchronized key expiry.
