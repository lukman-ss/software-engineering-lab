# Engineering Test Audit

## Execution Record (actual, fresh, 2026-09-27)

- `go build ./...` -> SUCCESS, exit 0, Go 1.26.7 darwin/arm64.
- `go test -count=1 -v ./...` -> all 11 tests PASS, exit 0. Counts: backpressure 3/3, httputil 1/1, ratelimit 5/5, retry 2/2.
- `go test -race -count=1 ./...` -> all 4 packages ok, no data races, exit 0. Fresh (uncached); earlier run was cached and was re-executed uncached before approving this statement.
- `go run ./cmd/demo` -> exit 0. Real output recorded in 04-docs-vs-code.md Finding 2.
- `go vet ./...` -> clean, exit 0.
- `gofmt -l .` -> flags 3 files: internal/backpressure/queue.go, internal/ratelimit/bucket.go, internal/ratelimit/registry.go (whitespace/alignment). Code compiles and tests pass; cosmetic only, not blocking.

## Coverage Matrix

| Area | Happy path | Failure path | Edge cases | Transitions | Recovery | Rollback | Concurrency | Negative |
|---|---|---|---|---|---|---|---|---|
| TokenBucket | PASS (burst 3) | PASS (4th rejected) | PARTIAL (no capacity=0, refill=0, AllowN>1, negative n) | PASS (refill after sleep) | N/A | N/A | PASS (50 goroutines, race clean) | WEAK (no invalid-input test) |
| LeakyBucket | PASS (2 allows) | PASS (burst rejected) | PARTIAL (no capacity=0, leak=0) | PASS (leak after sleep) | N/A | N/A | WEAK (no dedicated leaky concurrency test) | WEAK |
| Registry | PASS (isolation a vs b) | N/A | MISSING (concurrent Get same key; many tenants; unbounded map growth) | N/A | N/A | N/A | MISSING (no concurrent same-key Get test; double-checked locking unproven) | MISSING |
| BoundedQueue | PASS (enqueue) | PASS (ErrQueueFull) | PASS (Stop idempotent; submit-after-stop -> ErrQueueStopped) | PASS (accepted+rejected accounting) | WEAK (job error ignored `_ = job(ctx)`; no propagation) | N/A | PARTIAL (submit concurrency tested; Stop-during-submit NOT tested -> Finding 1) | PASS (ErrQueueStopped) |
| Retry backoff | PASS (bounds all strategies) | N/A | WEAK (attempt overflow large; base=0; cap=0; unknown strategy default) | N/A | N/A | N/A | N/A | MISSING (no distribution test; bounds only) |
| Middleware | PASS (200 then 429) | PASS (429 path) | MISSING (no anonymous fallback test; no per-tenant isolation at HTTP level; no Retry-After value correctness; no body shape check) | PASS (exhaustion transition) | N/A | N/A | MISSING (no concurrent HTTP test) | WEAK |

## Test-by-test Assessment

1. TestTokenBucket_BurstAndRefill: PASS, proves burst + exhaustion + refill. Timing-based (200ms sleep) but margin (~2 tokens) is safe. Adequate.
2. TestLeakyBucket_LeakRate: PASS, proves burst reject + drain. Adequate for happy/failure; no concurrency coverage.
3. TestRegistry_TenantIsolation: PASS, proves cross-tenant isolation. Does NOT prove same-key concurrent creation safety (double-checked locking pattern in registry.go:21-37). Registry has no mutex-guarded test under race.
4. TestTokenBucket_RetryAfterSeconds: PASS (not listed in 03-execution-result.md; record is stale). Proves >0 when empty and 0 after refill. Good.
5. TestTokenBucket_ConcurrencyRace: PASS under -race. 50 goroutines x 10 Allow. Proves no data race and no deadlock; does NOT assert token conservation (accepted <= capacity+refill), so it is a smoke test, not a correctness proof. Acceptable but weak.
6. TestBoundedQueue_RejectionUnderLoad: PASS. Deterministic (blocking first job + startedCh gate). Proves fast ErrQueueFull. Strong.
7. TestBoundedQueue_ConcurrencySafety: PASS. Asserts accepted+rejected == submitted. Good accounting proof; does not assert processed accounting or Stop interplay.
8. TestBoundedQueue_SubmitAfterStop: PASS. Proves idempotent Stop and ErrQueueStopped. Sequential only; concurrent Stop+Submit untested (see code-audit Finding 1).
9. TestRateLimitMiddleware_RFC6585: PASS. Asserts 200 then 429 and Retry-After presence. Does NOT assert Retry-After is a positive integer, correct value, Content-Type, or body JSON shape; no anonymous-tenant or multi-tenant case.
10. TestComputeBackoff_Bounds: PASS for attempts 0..9. Bounds-only; NoJitter lower bound asserted as >= Base which is wrong for attempt<0 (not tested) but fine. Note: for attempt >= 5 with Base=100ms Cap=2s, temp is capped so bounds hold; correct.
11. TestDecorrelatedJitter_Bounds: PASS. Chains prev=sleep, asserts [Base,Cap]. Good.

## Verdict on Suite Strength

Passing suite is real (fresh run, race clean) but has gaps: the most important missing test is concurrent Stop+Submit for BoundedQueue, where a confirmed panic exists. Middleware assertions are shallow. Registry concurrency is unproven. Jitter tests prove bounds, not distribution shape (acceptable for bounds claims; distribution claims in demo/README are illustrative, not asserted by tests).
