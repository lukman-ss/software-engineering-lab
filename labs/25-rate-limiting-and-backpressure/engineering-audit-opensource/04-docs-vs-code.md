# Docs vs Code Audit

Target Lab: labs/25-rate-limiting-and-backpressure

## Comparison Matrix

| Asset | Reference | Actual Code | Verdict |
|---|---|---|---|
| README Features | Token Bucket + Leaky Bucket, per-tenant registry, BoundedQueue non-blocking TrySubmit w/ ErrQueueFull, AWS 4 jitter strategies, RFC 6585 429 + Retry-After | bucket.go, registry.go, queue.go, backoff.go, middleware.go all implemented as described | PASS — claims match |
| README Structure | `internal/ratelimit`, `internal/backpressure`, `internal/retry`, `internal/httputil`, `cmd/demo` | exactly present on disk | PASS |
| README Commands | `go test ./...`, `go test -race ./...`, `go run ./cmd/demo` | all executed by auditor, PASS | PASS |
| Design §Expected Behavior | token burst/refill, leaky smoothing, bounded reject, AWS full-jitter [0, min(cap,...)], 429 Retry-After | implemented identically | PASS |
| Design §Components | lists ratelimit/backpressure/retry/http packages | package dirs match (httputil named http in prose) | PASS |
| Design §Test Strategy filenames | ratelimit_test.go, backpressure_test.go, retry_test.go, http_test.go | actual files: bucket_test.go, queue_test.go, backoff_test.go, middleware_test.go | **DOC_CODE_MISMATCH** — 4 filenames wrong (LOW) |
| Design Decision 3 / Impl Notes | "Deterministic simulated time helpers in tests to avoid flaky wall-clock sleeps" | tests use real `time.Sleep` (200ms/250ms/600ms). No fake clock anywhere. | **DOC_CODE_MISMATCH** — false claim |
| engineering/03-execution-result.md | demo queue output: Accepted=3 Rejected=3 Processed=1; specific jitter millisecond values | demo is run by auditor 3x — queue split varies (3-4 accepted, 2-3 rejected); jitter is random by design | **DOC_CODE_MISMATCH** — execution-result captures a non-deterministic snapshot presented as static |
| README "What Is Not Demonstrated" (not claimed) | distributed rate limiting, autoscaling | not implemented | PASS (consistent) |

## Research Alignment

Excluded from this audit per pipeline override (implementation + tests only). No research file contents cross-checked.

## Key Mismatches (Detail)

### MISMATCH-1 — Simulated time (FALSE CLAIM)
- Claimed: `engineering/01-design.md` Decision 3 and `engineering/02-implementation-notes.md` §What Is Demonstrated (implied) state tests use deterministic simulated-time helpers.
- Actual: `bucket_test.go` uses `time.Sleep(...)`. There is no fake clock / sim clock in any test file. Tests pass because sleeps are long, not because of simulated time.
- Impact: misrepresents test determinism. Tests are in fact real-time but stable. MEDIUM.

### MISMATCH-2 — Execution result non-determinism
- Recorded (`03-execution-result.md` §3):
  ```
  Job #4: REJECTED, Job #5: REJECTED, Job #6: REJECTED
  Stats: Accepted=3, Rejected=3, Processed=1
  ```
- Auditor replay (3 runs):
  - Run 1: Accepted=4, Rejected=2, Processed=2
  - Run 2: Accepted=3, Rejected=3, Processed=1
  - Run 3: Accepted=3, Rejected=3, Processed=1
- Cause: `Job #1` sleeps 50ms while submissions are instantaneous — worker may drain slot before Job #5 is enqueued. The exact accepted/rejected split is a TOCTOU, not deterministic.
- Impact: `03-execution-result.md` is a real (not fabricated) run, but it is presented as the reproducible "result." Jitter columns are likewise non-reproducible values (random), which is expected and flagged in prose, but the queue stat is not. MEDIUM.

### MISMATCH-3 — Test filenames in design doc
- Design §Test Strategy names four `_test.go` files that do not exist. The actual test files use different names. Aspirational planning text; LOW.

## Summary

| Dimension | Result |
|---|---|
| README vs implementation | PASS |
| Design claims vs code | PASS (algorithms, middleware, jitter, buckets) |
| Test strategy filenames | WARNING (3 wrong filenames) |
| Simulated-time test claim | FAIL/ WARNING (false — real sleeps) |
| Execution result determinism | WARNING (queue counts vary run to run) |
| Research alignment | NOT_AUDITED (excluded) |

Overall documentation accuracy: WARNING. Three doc/code mismatches; none indicate fabrication of results, but one (simulated time) is a direct false statement about the test methodology, and one (execution-result) over-specifies timing-dependent values.