# Test Audit

## Test Run Records (observed, uncached race run)
- `go build ./...` → ok (BUILD_OK)
- `go vet ./...` → clean (exit 0)
- `go test -v ./...` → 5 tests PASS (cache not cached: 0.05s+0.03s+0.00s+0.09s+0.00s)
- `go test -race -count=1 ./...` → ok (1.614s) fresh (not cached)
- `go run ./cmd/demo` → output matches engineering/03-execution-result.md claim (jitter
  sample values differ only by randomness). Demo verified real.

## Coverage Matrix

| Claim / Behavior | Test | Result | Strong? |
|---|---|---|---|
| Cache-Aside miss→DB→populate | TestCachePatterns/Cache-Aside | PASS | YES |
| Cache-Aside hit→no query | TestCachePatterns/Cache-Aside | PASS | YES |
| Cache-Aside update→invalidate+refetch | TestCachePatterns/Cache-Aside | PASS | YES |
| Write-Through sync write+cache | TestCachePatterns/Write-Through | PASS | YES |
| Write-Through next read hits cache | TestCachePatterns/Write-Through | PASS | YES |
| Write-Behind immediate cache | TestCachePatterns/Write-Behind | PASS | YES |
| Write-Behind async DB flush (count) | TestCachePatterns/Write-Behind | PASS (50ms Sleep) | YES (timing-coupled) |
| Naive stampede N→N queries | TestStampedeMitigation/Naive | PASS (>1) | YES |
| SingleFlight N→1 queries + correct | TestStampedeMitigation/SingleFlight | PASS (exactly 1) | YES |
| SingleFlight concurrency race-free | go test -race | PASS | YES (20 goroutines) |
| XFetch formula sign (math check) | TestXFetchLogic | PASS | YES (pure fn only) |
| XFetch end-to-end early refresh | _NONE (demo only) | — | NO |
| SWR stale-immediate + async refresh | TestStaleWhileRevalidate | PASS (value check) | PARTIAL |
| TTL jitter range | TestJitter | PASS (100 iterations) | YES |

## Coverage Gaps

### Happy path
All primary happy paths covered. PASS.

### Failure / error paths
- NO test for any DB error in service layer: CacheAside.Get with ErrNotFound, WriteThrough.Update
  with db.Write failure, SingleFlight with a failing flight func, XFetch recompute-fallback-on-error
  path (stampede.go:159-164 returns stale on error), SWR sync-fetch error.
- MockDB.SetData is non-failing; Write never tested to return error (delay+context only).
- These paths exist in code but are unreachable-under-test → confidence gap.

## Finding 21

Location: tests/cache_test.go:84-103 (WriteBehind test)
Claimed Behavior: Verifies async flush.
Observed Implementation: `svc.Update(...); svc.Get; time.Sleep(50ms); assert WriteCount==1`. No
  assertion on queue-overflow drop path; no Close-drain assertion (drain is implicit via defer
  after sleep).
Assessment: WARNING
Severity: MEDIUM
Notes: 50ms Sleep coupled. Overflow drop (patterns.go:156-160 default branch) never asserted.
Gap type: MISSING_TEST

## Finding 22

Location: tests/cache_test.go:161-184 (TestXFetchLogic)
Claimed Behavior: Verifies XFetch formula correctness.
Observed Implementation: Tests `ShouldRecompute` pure function with deterministic u values;
  asserts correct sign vs erroneous formula. Does NOT test `XFetchService.Get` early-refresh
  behavior end-to-end.
Assessment: WARNING
Severity: MEDIUM
Notes: Math proven; service Get() (stampede.go:138-169) only exercised via demo output.
Gap type: MISSING_TEST

## Finding 23

Location: tests/cache_test.go:186-221 (TestStaleWhileRevalidate)
Claimed Behavior: SWR serves stale immediately and revalidates.
Observed Implementation: Asserts stale value returned, then fresh value returned after 50ms
  sleep. Does NOT assert `RevalidateCount()>`; does NOT assert exactly one in-flight
  revalidation (single-revalidation-in-flight guard at stampede.go:228-231 not verified);
  does NOT assert the stale value persisted during async refresh.
Assessment: WARNING
Severity: MEDIUM
Notes: Test passes but only value-checked, not behaviorally strict; timing-coupled (Sleep).
Gap type: MISSING_TEST / MISSING_EDGE_CASE

## Finding 24

Location: tests/cache_test.go (whole file)
Claimed Behavior: Concurrency + stampede coverage.
Observed Implementation: Only SingleFlight/Naive use concurrency (20 goroutines). XFetch/SWR
  Get never tested concurrently; `XFetchService.SetRandFunc` writes randFunc with NO lock while
  `getRand` reads it — latent data race if SetRandFunc called concurrently with Get (unexercised).
Assessment: WARNING
Severity: MEDIUM
Notes: -race passes (not invoked), but cross-conc of SetRandFunc/Get is unsafe. Singleflight Get is
  the only concurrent service path tested. Race-free in tested scope, but scope is narrow.
Gap type: MISSING_TEST

## Finding 25

Location: tests/cache_test.go (whole file)
Claimed Behavior: Error / context-cancellation coverage.
Observed Implementation: ctx is always context.Background() with no deadline; no cancelled ctx
  passed to any service.Get/Update. MockDB honours ctx (repo.go:37-42) but never challenged.
Assessment: WARNING
Severity: MEDIUM
Notes: No failure-path or cancellation test anywhere.
Gap type: MISSING_TEST

## Summary
- Tests PASS and prove: 3 write patterns happy path, naive vs singleflight query accounting,
  XFetch pure-formula sign correctness, SWR stale-then-fresh value transition, jitter range.
- -race PASS on tested concurrency (20-goroutine singleflight). go vet clean.
- Weaknesses: no error/failure-path tests; XFetch service Get e2e unverified; SWR behavioural
  assertions weak (value-only, timing-coupled); SetRandFunc/goroutine data race latent; WriteBehind
  overflow + drain not asserted; no ctx-cancellation tests.
- No fake results: all counters sourced from code (atomic) + observed demo run; test counts
  match engineering/03-execution-result.md exactly.
