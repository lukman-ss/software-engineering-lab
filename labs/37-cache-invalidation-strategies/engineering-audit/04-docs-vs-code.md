# Documentation vs Code Audit

Target Lab: labs/37-cache-invalidation-strategies

## README vs Code

| README Claim | Code Reality | Status |
|---|---|---|
| Correct command `go test -v ./...` | Verified working | MATCH |
| Correct command `go test -race ./...` | Verified working | MATCH |
| Correct command `go run ./cmd/demo` | Verified working and output matches | MATCH |
| "Cache-Aside: Checks cache, loads from DB on miss, populates cache. Invalidates cache on update." | Implemented in `patterns.go:22-47`. `Update` calls DB write then `cache.Delete`. | MATCH |
| "Write-Through: Updates DB and cache synchronously on write; subsequent reads hit cache." | Implemented in `patterns.go:78-89`. Confirmed no extra DB query after write. | MATCH |
| "Write-Behind: Updates cache immediately; flushes to DB asynchronously via background worker queue." | Implemented in `patterns.go:98-166` with goroutine worker and buffered channel. | MATCH |
| "SingleFlight: Coalesces concurrent cache miss requests on a hot key down to a single DB query using `golang.org/x/sync/singleflight`." | Implemented in `stampede.go:45-84`. | MATCH |
| "XFetch: Early refresh using formula `-Δ · β · ln(U) > TTL_remaining`." | Formula `expiryCompute := -deltaSec * beta * math.Log(u)` in `stampede.go:134`. Correct sign. | MATCH |
| "SWR: Returns stale cached data immediately while asynchronously triggering background DB revalidation." | Three-tier logic in `stampede.go:197-224`, background goroutine in `stampede.go:226-253`. | MATCH |
| "TTL Jitter: Adds randomized offset `[0, maxJitter)` to base TTL." | `rand.Int63n(int64(maxJitter))` in `store.go:84` — correct `[0, maxJitter)` range. | MATCH |
| Code structure: `internal/cache/store.go`, `repo.go`, `patterns.go`, `stampede.go` | All exist and match described purpose | MATCH |
| Code structure: no `jitter.go` listed | README does not list `jitter.go`. No discrepancy in README. | MATCH |

## Engineering Notes vs Code

| Design Doc Claim | Code Reality | Status |
|---|---|---|
| `engineering/01-design.md:46` — Architecture lists `jitter.go: TTL jitter calculation` | `TTLWithJitter` is implemented in `store.go`. No `jitter.go` file exists. | DOC_CODE_MISMATCH (MINOR) |
| Double-check inside singleflight `Do` function | Present in `stampede.go:64-67` | MATCH |
| `GetRaw` needed for SWR and XFetch | `GetRaw` implemented in `store.go:46-51`; called by both XFetch and SWR services | MATCH |
| `ReadDelta` stored per item for XFetch | `ReadDelta` field in `Item` struct, stored in `Set`, used in `ShouldRecompute` | MATCH |
| Write-Behind `Close()` drains remaining queue | `flushWorker` drains `writeQueue` on quit channel signal | MATCH |
| `SetRandFunc` injectable rand for deterministic tests | Implemented in `stampede.go:109-111`; used in `TestXFetchService_Get` | MATCH |

## Research vs Implementation

| Research Requirement | Implementation | Status |
|---|---|---|
| XFetch formula correct sign: `-Δ · β · ln(U) > TTL_remaining` (research revision finding) | `expiryCompute := -deltaSec * beta * math.Log(u)` with `expiryCompute > ttlRemainingSec` | MATCH |
| Singleflight reduces N concurrent misses to 1 DB query | Verified by test and demo output (20 goroutines → 1 DB query) | MATCH |
| SWR: stale-within-window returns cached value, async revalidation | Implemented and verified by test | MATCH |
| Write-Behind: data loss risk on un-flushed queue | Acknowledged in implementation notes; overflow drops writes silently | MATCH |
| TTL Jitter: desynchronize key expiration | Implemented and tested with 100-iteration bounds check | MATCH |

## Demo Output vs Code

| Demo Output Claim | Code Behavior | Status |
|---|---|---|
| Naive 20 goroutines → 20 DB queries | Demo shows 20 queries; confirmed by `NaiveStampedeService` design | MATCH (non-deterministic but reliably >1 in practice) |
| SingleFlight 20 goroutines → 1 DB query | Demo shows 1 query; `singleflight.Group.Do` coalesces requests | MATCH |
| Write-Behind writes=0 immediately, writes=1 after 50ms | Confirmed by background worker channel timing | MATCH |
| SWR serves stale v1, then fresh v2 | Confirmed by 3-tier SWR logic with 30ms sleep + 50ms async wait | MATCH |
| TTL Jitter sampled TTLs all > base (5m0s) and < base+maxJitter (5m30s) | Verified by `rand.Int63n(int64(maxJitter))` range | MATCH |

## Issues

DOC_CODE_MISMATCH (LOW): `engineering/01-design.md` lists `jitter.go` as a separate architecture component, but implementation consolidates TTL jitter into `store.go`. No `jitter.go` exists. README correctly omits this; only the design doc is inconsistent.
