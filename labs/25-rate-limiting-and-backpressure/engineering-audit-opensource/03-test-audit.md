# Test Audit

Target Lab: labs/25-rate-limiting-and-backpressure
Date: 2026-09-27
Commands executed (fresh, -count=1):
- go build ./... : SUCCESS
- go test -v -count=1 ./... : ALL PASS (11 tests: backpressure 3, httputil 1, ratelimit 5, retry 2)
- go test -race -count=1 ./... : ALL PASS, no data races
- go vet ./... : CLEAN
- go run ./cmd/demo : SUCCESS (output captured, see 02-code-audit Finding 10)

## Coverage vs Required Dimensions

- Happy path: PASS. Burst/refill, leak/admit, enqueue, jitter bounds, 200 then 429.
- Failure path: PASS. Exhaustion rejection, ErrQueueFull, ErrQueueStopped, 429+header.
- Edge cases: PARTIAL. Empty->full transitions, stop-idempotency, decorrelated chain, anonymous fallback implicit only. Missing: zero refill/leak rate, negative attempt, unknown strategy, concurrent Stop+Submit.
- Transitions: PASS. Refill-after-sleep, leak-after-sleep, retry-after 0 after refill.
- Recovery: PASS (time-based refill/leak proven by sleeps). No restart/crash recovery claimed.
- Rollback: NOT_APPLICABLE (no transactions/state rollback in design).
- Concurrency: PASS with WARNING. 50-goroutine token pounding + 30-goroutine queue submits, race detector clean. But TokenBucket_ConcurrencyRace asserts nothing (smoke only); Queue concurrency asserts only accepted+rejected==30. Stop+Submit race untested.
- Negative cases: PARTIAL. Wrong behavior asserted absent (4th token denied, 3rd water denied). Missing negative checks listed under edge cases.

## Test Quality Notes

- Wall-clock sleeps (200/250/600ms) used; deterministic enough at chosen margins, slightly slow (~1.1s ratelimit suite). Acceptable.
- Backpressure rejection test correctly blocks worker via blockCh/startedCh to force full buffer. Sound.
- Retry tests assert bounds over attempts 0..9 including cap saturation. No distribution/statistical test; bounds-only is adequate for stated claim.
- Middleware test verifies status + Retry-After presence, not body JSON or header numeric value.
- Suite stronger than appearance: includes RetryAfterSeconds round-trip and SubmitAfterStop idempotency, beyond the 9 tests listed in 03-execution-result.md (which is stale).
