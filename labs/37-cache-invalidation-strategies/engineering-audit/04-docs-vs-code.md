# Docs vs Code Audit

## Comparison Matrix

| Component / Claim | Documented Claim (README & Engineering Notes) | Actual Implementation in Code | Mismatch Detected |
|---|---|---|---|
| Module Path | `github.com/lukman/labs/37-cache-invalidation-strategies` | `go.mod`: `github.com/lukman/labs/37-cache-invalidation-strategies` | None |
| Cache-Aside | Reads DB on miss, populates cache; invalidates on update | `internal/cache/patterns.go:12-48` | None |
| Write-Through | Updates DB and cache synchronously | `internal/cache/patterns.go:50-90` | None |
| Write-Behind | Updates cache immediately; flushes to DB asynchronously | `internal/cache/patterns.go:92-166` | None |
| SingleFlight | Coalesces concurrent cache misses to 1 DB query | `internal/cache/stampede.go:45-85` | None |
| XFetch Formula | `-Δ · β · ln(U) > TTL_remaining` | `internal/cache/stampede.go:122-136` | None |
| Stale-While-Revalidate | Serves stale data immediately; triggers async revalidation | `internal/cache/stampede.go:173-254` | None |
| TTL Jitter | Adds random offset `[0, maxJitter)` to base TTL | `internal/cache/store.go:79-85` | None |
| Demo Output | 7 sections matching claimed behaviors | `cmd/demo/main.go` | None |

---

## Detailed Checks

### 1. README vs Code
- **Structure**: All directory tree entries in `README.md` (lines 8-24) exist at exact locations.
- **Commands**: `go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo` all run cleanly without errors.
- **Features List**: All 7 features described in `README.md` (lines 44-50) are implemented and exercised in tests.

### 2. Engineering Notes vs Code
- In-memory store decision accurately reflects lack of Redis requirement.
- Singleflight in-process limitation accurately recorded.
- Write-behind drop behavior matches `select { default: }` in `patterns.go:158`.
- XFetch formula guard `u <= 0 || u >= 1` in `stampede.go:126` matches engineering note 28.

### 3. Research Claims vs Code
- Finding 1 (Cache-Aside read/write lifecycle): Implemented in `patterns.go`.
- Finding 2 (Write-Through synchronous write): Implemented in `patterns.go`.
- Finding 3 (Write-Behind async flush durability risk): Implemented in `patterns.go`.
- Finding 4 (Stampede naive N vs singleflight 1): Implemented in `stampede.go`.
- Finding 5 (XFetch formula with negative sign): Correct formula `-deltaSec * beta * math.Log(u)` in `stampede.go`.
- Finding 6 (XFetch beta parameter default 1.0): Used in `cmd/demo/main.go` and `tests/cache_test.go`.
- Finding 7 (SWR stale return + async revalidation): Implemented in `stampede.go`.
- Finding 8 (TTL Jitter anti-synchronization): Implemented in `store.go`.

---

## Discrepancies / Overclaims

No documentation mismatches or overclaims detected.
The code strictly delivers the scope promised by the documentation and research reports.
