# Gap Analysis

Target: labs/25-rate-limiting-and-backpressure

## Gaps

### GAP-1: Stop() lifecycle not idempotent / submit-vs-close race

Type: RACE_CONDITION
Severity: MEDIUM
Location: internal/backpressure/queue.go:76-80
Detail: Stop() closes channel; concurrent TrySubmit panics (send on closed channel). No guard/atomic closed flag.
Status: Untested. Safe in tests/demo by sequencing, but real concurrent use unsafe.

### GAP-2: Zero/negative constructor inputs unvalidated

Type: UNHANDLED_ERROR
Severity: LOW
Location: ratelimit NewTokenBucket/NewLeakyBucket, backpressure NewBoundedQueue
Detail: refillRate==0 -> RetryAfterSeconds div-by-zero (float->int). capacity<=0/workers<=0 -> unbuffered/undrained queue. No constructor guards, no tests.
Status: Untested edge.

### GAP-3: Registry unbounded growth

Type: MISSING_EDGE_CASE
Severity: LOW
Location: internal/ratelimit/registry.go
Detail: No eviction/TTL; unbounded tenant cardinality grows map indefinitely. Acknowledged in-memory scope.
Status: Accepted scope gap.

### GAP-4: BoundedQueue job errors swallowed, queued jobs dropped on Stop

Type: MISSING_EDGE_CASE
Severity: LOW
Location: internal/backpressure/queue.go:54, 76-80
Detail: Job return value ignored (`_ = job(...)`); Stop discards buffered jobs instead of draining. Tests assert acceptance only.
Status: Deliberate minimal demo behavior.

### GAP-5: Retry tests assert bounds, not exact formulas

Type: MISSING_TEST
Severity: LOW
Location: internal/retry/backoff_test.go
Detail: EqualJitter lower bound (temp/2), FullJitter lower=0, NoJitter formula exactness, negative attempt handling not asserted.
Status: Coverage gap, behavior verified manually in code audit.

### GAP-6: Middleware test asserts header presence, not value/schema

Type: MISSING_TEST
Severity: LOW
Location: internal/httputil/middleware_test.go
Detail: Retry-After integer value correctness, JSON error body, anonymous fallback, post-refill recovery untested.
Status: Coverage gap.

### GAP-7: Design wording (simulated time helpers)

Type: DOC_CODE_MISMATCH
Severity: LOW
Location: engineering/01-design.md (Implementation Decisions)
Detail: Claims deterministic simulated-time helpers; tests use real time.Sleep. Behavior still proven.
Status: Wording inaccuracy only.

### GAP-8: Tokens() returns stale stored value (no pending refill applied)

Type: DOC_CODE_MISMATCH
Severity: LOW
Location: internal/ratelimit/bucket.go:48-52
Detail: Demo displays Remaining Tokens from stale field; limiting decision itself refills correctly.
Status: Cosmetic only.

## Gap Type Checklist

- MISSING_TEST: GAP-5, GAP-6 (bounds/presence only, not value-level)
- RACE_CONDITION: GAP-1
- UNHANDLED_ERROR: GAP-2
- MISSING_EDGE_CASE: GAP-3, GAP-4
- DOC_CODE_MISMATCH: GAP-7, GAP-8
- BROKEN_IMPLEMENTATION: none
- IMPLEMENTATION_OVERCLAIM: none
- RESEARCH_MISMATCH: out of scope (pipeline override)
- FAKE_DEMO: none (re-ran, real)
- FAKE_BENCHMARK: none (no benchmarks claimed)
- UNVERIFIED_RESULT: none (all green verified)
