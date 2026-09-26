# Gap Analysis

## GAP-01

Type: MISSING_TEST
Severity: MEDIUM
Location: internal/httputil/middleware_test.go
Claim affected: RFC 6585 middleware with per-tenant isolation.
Observation: Only one tenant key tested; Retry-After value not validated; anonymous fallback path untested.
Evidence: middleware_test.go:11-43 covers single key, asserts 200 then 429 + header present.

## GAP-02

Type: MISSING_TEST
Severity: MEDIUM
Location: internal/ratelimit/bucket_test.go:83-98
Claim affected: "TokenBucket accurately track and reject ... under concurrent access" (01-design.md success criteria).
Observation: Concurrency test asserts nothing beyond no-panic; no check that total allowed tokens ≤ capacity + refill.
Evidence: Test calls Allow() 500x, then wg.Wait() with no assertion.

## GAP-03

Type: DOC_CODE_MISMATCH
Severity: MEDIUM
Location: engineering/01-design.md:93 (Decision 3)
Observation: Design claims deterministic simulated-time tests; code uses real time.Sleep + math/rand.
Impact: Flaky tests possible; NoJitter deterministic but jitter outputs random.

## GAP-04

Type: MISSING_TEST
Severity: MEDIUM
Location: internal/retry/backoff_test.go
Claim affected: AWS jitter spec compliance (EqualJitter exact half formula, DecorrelatedJitter range).
Observation: Tests only check bounds [base,cap], not formula lower values e.g. EqualJitter ≥ temp/2, DecorrelatedJitter ≥ base. No determinism via seed.
Evidence: backoff_test.go:8-48 bounds checks only.

## GAP-05

Type: UNVERIFIED_RESULT
Severity: MEDIUM
Location: engineering/03-execution-result.md:79-83
Claim affected: Demo jitter values recorded as fixed (FullJitter 99/109/111/125ms).
Observation: Backoff uses math/rand.Float64 unseeded; values differ each run (observed FullJitter 8/122/376/141ms). Recorded values unreproducible.
Impact: Execution result implies determinism; should note randomness or seed.

## GAP-06

Type: MISSING_EDGE_CASE
Severity: MEDIUM
Location: internal/ratelimit/bucket.go:69-78 (RetryAfterSeconds with refillRate=0)
Observation: Division by refillRate without zero guard. Current `if secs <= 0 return 1` hides +Inf→int conversion edge; contract unclear.
Impact: Zero-rate bucket (never refills) behavior undefined.

## GAP-07

Type: IMPLEMENTATION_OVERCLAIM
Severity: LOW
Location: cmd/demo/main.go vs design claims
Observation: Demo exercises Token/Leaky/BoundedQueue/No/Full/Equal jitter but omits DecorrelatedJitter and HTTP middleware 429 path.
Impact: "Interactive CLI displaying ... retry jitter distributions" partially unproven in demo.

## GAP-08

Type: DOC_CODE_MISMATCH
Severity: LOW
Location: engineering/01-design.md:52, engineering/02-implementation-notes.md implications
Observation: Architecture shows "Queue Full? -> 503 Overloaded" but code's HTTP middleware returns only 429; BoundedQueue returns ErrQueueFull (no HTTP mapping). No 503 path implemented.
Impact: Design over-describes HTTP surface not present.

## No gaps found (explicitly checked, clean)

- Compilation: PASS (go build ./... EXIT 0, go vet clean)
- Race detector: PASS (go test -race ./... all ok)
- Token/Leaky burst + refill/leak: proven by tests + demo
- BoundedQueue fast rejection: proven
- No RACE_CONDITION found beyond weak assertions
- No BROKEN_IMPLEMENTATION of claimed core paths
- No FAKE_BENCHMARK (no benchmarks claimed)
