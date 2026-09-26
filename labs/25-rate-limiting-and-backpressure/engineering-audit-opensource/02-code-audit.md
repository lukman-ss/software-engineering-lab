# Engineering Code Audit

Target Lab: labs/25-rate-limiting-and-backpressure
Scope: implementation files only (source + tests executed). Research excluded per pipeline override.

Verification commands executed (recorded separately in 03-test-audit.md):
- `go build ./...`            -> SUCCESS (exit 0)
- `go vet ./...`              -> SUCCESS (exit 0)
- `go test -v -count=1 ./...` -> all packages PASS
- `go test -race ./...`       -> all packages PASS (race)

## Finding 1

Location: internal/ratelimit/bucket.go:8-79 (TokenBucket)
Claimed Behavior: Token bucket allows bursts up to capacity B, refills continuously at rate R tokens/sec, thread-safe.
Observed Implementation: `AllowN` acquires `sync.Mutex`, computes `elapsed = now - lastRefill`, `tokens = min(capacity, tokens + elapsed*rate)`, consumes n if `tokens >= n`, updates `lastRefill`. `Tokens()`/`RetryAfterSeconds` share the lock.
Assessment: PASS
Severity: N/A
Notes: Correct token-bucket semantics. `time.Now().Sub(...)` uses Go's monotonic clock component, so it is immune to wall-clock adjustments — the design's "monotonic clock" claim is accurate.

## Finding 2

Location: internal/ratelimit/bucket.go:54-79 (RetryAfterSeconds)
Claimed Behavior: Returns seconds to wait until n tokens are available; 0 when sufficient tokens exist.
Observed Implementation: Refills into a local `tokens` (not mutating state), returns 0 if `tokens >= n`, else `ceil((n-tokens)/rate)`, minimum 1.
Assessment: PASS (with edge caveat)
Severity: LOW
Notes: Ceiling rounding correct. If `refillRate == 0` the division `needed/tb.refillRate` yields +Inf -> `int(+Inf)` is implementation-defined (large/garbage). Not exercised by tests; capacity/refill are always > 0 in tests.

## Finding 3

Location: internal/ratelimit/bucket.go:81-121 (LeakyBucket)
Claimed Behavior: Constant drain at rate R, rejects burst when water reaches capacity.
Observed Implementation: `water = max(0, water - elapsed*leakRate)`; allow iff `water+1 <= capacity`. Mutex guards all reads/writes.
Assessment: PASS
Severity: N/A
Notes: Matches demo output (3/5 allowed at capacity 3). `Water()` returns last-drained level without refilling (cosmetic staleness).

## Finding 4

Location: internal/ratelimit/registry.go:6-37 (Registry)
Claimed Behavior: Per-tenant bucket isolation (CGNAT/RFC 6598 avoidance) via key-based buckets; thread-safe.
Observed Implementation: Double-checked locking with `sync.RWMutex`. `Get` returns existing bucket or creates+stores a new one.
Assessment: PASS (with resource gap)
Severity: MEDIUM
Notes: Concurrency pattern correct. Buckets are never evicted — map grows unbounded per unique tenant/API-key. Not enumerated in "Known Limitations" (which only lists distributed state). Resource-leak risk for high-cardinality, long-running tenants.

## Finding 5  [BLOCKING]

Location: internal/backpressure/queue.go:66-81 (TrySubmit) vs queue.go:87-94 (Stop)
Claimed Behavior: `BoundedQueue` is thread-safe; `Stop` is idempotent and safe; `TrySubmit` fast-fails with `ErrQueueFull`/`ErrQueueStopped` non-blockingly under concurrent callers.
Observed Implementation: `Stop()` does `CAS(stopped)=true` -> `cancel()` -> `close(bq.queue)` -> `wg.Wait()`. `TrySubmit` first checks `stopped` then enters a `select` whose success branch is `case bq.queue <- job:`. There is no flag preventing the send case from being selected after `close(bq.queue)`.
Assessment: FAIL
Severity: HIGH
Notes: CONFIRMED at runtime via an audit probe (removed after run):
```
panic: send on closed channel
...queue.go:71 -> case bq.queue <- job
```
When `Stop()` runs concurrently with `TrySubmit`, the select may choose the send case (`bq.queue <- job`) after the channel is closed, panicking the process. The `<-ctx.Done()` case is also ready at that moment, so Go selects randomly between it and the send — roughly 50% panic when buffer has space. The shipped tests never exercise concurrent submit+stop, so `go test -race` is green despite the defect. This violates the design's explicit "concurrency safety" success criterion for the shutdown path.

## Finding 6

Location: internal/backpressure/queue.go:87-94 (Stop / workerLoop)
Claimed Behavior: Graceful shutdown drains queue and waits for in-flight jobs; repeated Stop is safe.
Observed Implementation: `CompareAndSwap` makes Stop idempotent; `close(queue)` unblocks workers (`ok==false` return); `wg.Wait()` waits for in-flight jobs. `Stats` reads `len(queue)` safely (never panics).
Assessment: PASS
Severity: N/A
Notes: Correct for the orderly (non-concurrent-submit) path.

## Finding 7

Location: internal/httputil/middleware.go:17-40
Claimed Behavior: RFC 6585 `429 Too Many Requests` with `Retry-After` header + JSON error body; tenant key from `X-API-Key`, fallback `anonymous`.
Observed Implementation: Reads header (fallback `anonymous`), `registry.Get`, `!Allow()` -> sets Content-Type, Retry-After, writes 429, encodes JSON `{error:"rate_limit_exceeded", retry_after:N}`. `Allow()` path delegates to `next`.
Assessment: PASS
Severity: N/A
Notes: Status/headers/body all produced. Thread-safe via bucket locking. `Retry-After` value derived from `RetryAfterSeconds`.

## Finding 8

Location: internal/retry/backoff.go:24-62
Claimed Behavior: NoJitter/FullJitter/EqualJitter/DecorrelatedJitter match AWS (Marc Brooker) formulas.
Observed Implementation:
- NoJitter: `min(cap, base*2^attempt)` — matches.
- FullJitter: `rand()*min(cap, base*2^attempt)` in [0,temp] — matches.
- EqualJitter: `temp/2 + rand()*temp/2` in [temp/2, temp] — matches.
- DecorrelatedJitter: `min(cap, rand(base, prev*3))` with prev clamped to >=base — matches.
Assessment: PASS
Severity: LOW
Notes: Uses global `math/rand` (deprecated since Go 1.20, functionally fine; `rand.Float64` is goroutine-safe). `temp <= 0` guard in FullJitter prevents NaN. EqualJitter/DecorrelatedJitter do not guard `temp<=0` but cfg.Base>0 in all uses, so safe in practice.

## Finding 9

Location: cmd/demo/main.go
Claimed Behavior: Interactive CLI demonstrating burst allowance, leaky smoothing, backpressure shedding, retry jitter.
Observed Implementation: Exercises all four components; prints token counts, water levels, accept/reject, backoff table. No arguments/flag handling.
Assessment: PASS
Severity: N/A
Notes: Demo runs cleanly (exit 0). Queue accept/reject counts are timing-dependent (3-4 accepted across runs) — see 04-docs-vs-code.md.

## Summary

| Area | Status |
|---|---|
| Token/Leaky bucket | PASS |
| Registry isolation | PASS (MEDIUM resource-growth caveat) |
| BoundedQueue shutdown safety | FAIL / HIGH — panic on concurrent submit+Stop |
| Middleware 429 | PASS |
| Retry jitter (AWS spec) | PASS |
| Build / vet / tests / race / demo | PASS |
