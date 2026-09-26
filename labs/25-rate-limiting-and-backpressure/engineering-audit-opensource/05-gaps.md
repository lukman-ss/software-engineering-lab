# Gap Analysis

Target Lab: labs/25-rate-limiting-and-backpressure

Allowed gap types used:
BROKEN_IMPLEMENTATION, MISSING_TEST, DOC_CODE_MISMATCH, RACE_CONDITION, UNHANDLED_ERROR, MISSING_EDGE_CASE, IMPLEMENTATION_OVERCLAIM, UNVERIFIED_RESULT

## Gaps

1. BROKEN_IMPLEMENTATION — `internal/backpressure/queue.go:66-81` `BoundedQueue.TrySubmit`.
   - `Stop()` closes `bq.queue`; `TrySubmit`'s select success branch `case bq.queue <- job:` can fire on a closed channel.
   - Confirmed: `panic: send on closed channel` at `queue.go:71`, reproduced by a temporary audit probe (removed, no repo change), under concurrent submit + Stop.
   - Core concurrency-safety claim broken for the shutdown path.
   - Severity: HIGH.

2. RACE_CONDITION — Same as #1, framed as a Go concurrency hazard. The success case of the select is selected over `ctx.Done()` non-deterministically when both are ready after `close`+`cancel`. No mutex/lock-free guard prevents the send.
   - Severity: HIGH.

3. MISSING_TEST — `internal/backpressure/queue_test.go` has no concurrent-submit-and-Stop test. This is the exact path that ships a panic; absence of the test is why `go test -race` is green while the bug lives.
   - Severity: HIGH.

4. MISSING_TEST — `internal/httputil/middleware_test.go` single test covers only one tenant's 200→429. Missing: anonymous fallback, JSON body shape/content, Retry-After numeric correctness, multi-tenant independence, 200 passthrough body.
   - Severity: MEDIUM.

5. MISSING_TEST — `internal/retry/backoff_test.go` checks bounds only. No exact-value assertions (e.g., NoJitter == min(cap, base*2^attempt)); no DecorrelatedJitter prev-clamp edge; no distribution sanity for Equal/Full jitter over many samples.
   - Severity: MEDIUM.

6. MISSING_EDGE_CASE — `internal/ratelimit/bucket.go` `RetryAfterSeconds`: `refillRate == 0` yields `needed/0` → +Inf → `int(+Inf)` (undefined). Not reachable via current constructors/tests but an unhandled-error surface.
   - Severity: LOW.

7. MISSING_EDGE_CASE — `internal/backpressure/queue.go` `Stats()` returns `len(bq.queue)` concurrently with sends; safe (no panic) but racy-read-ish approximation; documented? no. Minor, not a correctness bug.
   - Severity: LOW.

8. DOC_CODE_MISMATCH — `engineering/01-design.md` Decision 3 + `engineering/02-implementation-notes.md` claim "deterministic simulated time helpers in tests". Actual tests use real `time.Sleep`. No fake clock exists.
   - Severity: MEDIUM.

9. DOC_CODE_MISMATCH — `engineering/01-design.md` §Test Strategy references test files `ratelimit_test.go`, `backpressure_test.go`, `retry_test.go`, `http_test.go`. Actual files: `bucket_test.go`, `queue_test.go`, `backoff_test.go`, `middleware_test.go`.
   - Severity: LOW.

10. DOC_CODE_MISMATCH — `engineering/03-execution-result.md` records queue stats `Accepted=3, Rejected=3, Processed=1`. Auditor replay across 3 runs: split varies (3/4 accepted, 2/3 rejected). The recorded value is one TOCTOU snapshot, not reproducible deterministically.
    - Severity: MEDIUM.

11. MISSING_EDGE_CASE / resource leak — `internal/ratelimit/registry.go` never evicts tenant buckets. High-cardinality tenants → unbounded map growth. Not listed in "Known Limitations" (which only mentions distributed state).
    - Severity: MEDIUM.

12. UNVERIFIED_RESULT — `engineering/03-execution-result.md` Retry section prints specific FullJitter/EqualJitter millisecond values that are random per run. Expected for jitter but recorded as if static. No claim of determinism in prose; low concern.
    - Severity: LOW.

## Gap Summary by Severity

| Severity | Count | IDs |
|---|---|---|
| CRITICAL | 0 | — |
| HIGH | 3 | 1, 2, 3 |
| MEDIUM | 6 | 4, 5, 8, 10, 11 |
| LOW | 3 | 6, 7, 9, 12 |

## Notes

- No fabricated benchmarks/results: all tests genuinely pass; demo genuinely runs. The defect is a real concurrency bug, not a falsified result.
- The HIGH gaps (1–3) block APPROVED status (see verdict) because the design explicitly targets concurrency safety and the only thing masking it is a missing test.
- Gaps #1 and #2 are the same root cause described in two ways (implementation defect + race-condition framing).