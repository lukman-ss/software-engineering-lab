# Docs vs Code Audit

## Finding 1

Location: README.md:40-43 vs code (bucket.go, queue.go, retry/backoff.go, middleware.go)
Documented Features: 1) Token & Leaky bucket; 2) BoundedQueue Backpressure; 3) AWS Jitter (Full/Equal/No/Decorrelated); 4) HTTP 429 RFC 6585.
Code Coverage: Token/Leaky ✓, BoundedQueue ✓, HTTP 429 ✓. AWS Jitter: Full/Equal/No ✓, Decorrelated ✓ in code.
Doc gap: README "Matching AWS Architecture specifications (Marc Brooker)" — code matches; fine.
Assessment: PASS
Notes: README claims match implementation.

## Finding 2

Location: engineering/01-design.md:14 (Decision 3)
Claimed: "Deterministic simulated time helpers in tests to avoid flaky wall-clock sleeps."
Code Reality: bucket_test.go:25,45,77 use `time.Sleep` with real wall clock. Backoff tests rely on real rand.
Mismatch Type: DOC_CODE_MISMATCH
Assessment: FAIL
Severity: MEDIUM
Notes: Design docs assert simulated-time tests; actual tests sleep on wall clock → flakiness present in RetryAfterSeconds (0.60s) and leaky tests.

## Finding 3

Location: engineering/01-design.md:46 & README structure diagram vs cmd/demo
Claimed: "Interactive CLI displaying ... retry jitter distributions" + architecture includes Worker Pool → Client with Full Jitter Retry.
Code Reality: demo prints backoff but skips DecorrelatedJitter demonstration; no client retry loop shown; demo is non-interactive print-out.
Mismatch Type: DOC_CODE_MISMATCH
Assessment: WARNING
Severity: MEDIUM
Notes: "Interactive CLI" not interactive; DecorrelatedJitter omitted from demo output.

## Finding 4

Location: engineering/02-implementation-notes.md:22
Claimed: "AWS Jitter Specification Compliance: ... Full Jitter ..."
Code Matches: ✓
Assessment: PASS

## Finding 5

Location: engineering/03-execution-result.md demo output vs actual `go run ./cmd/demo` just now.
Claimed demo jitter lines:
  Attempt 0 -> NoJitter 100ms  FullJitter 99ms   EqualJitter 92ms
  Attempt 1 -> NoJitter 200ms  FullJitter 109ms  EqualJitter 183ms
  Attempt 2 -> NoJitter 400ms  FullJitter 111ms  EqualJitter 202ms
  Attempt 3 -> NoJitter 800ms  FullJitter 125ms  EqualJitter 709ms
Actual (just now):
  Attempt 0 -> NoJitter 100ms  FullJitter 8ms    EqualJitter 82ms
  Attempt 1 -> NoJitter 200ms  FullJitter 122ms  EqualJitter 198ms
  Attempt 2 -> NoJitter 400ms  FullJitter 376ms  EqualJitter 229ms
  Attempt 3 -> NoJitter 800ms  FullJitter 141ms  EqualJitter 782ms
Observation: FullJitter/EqualJitter are randomized; exact numbers differ every run. NoJitter deterministic matches (100/200/400/800). Documented specific jitter values cannot be reproduced — expected for random values, but recording exact values in 03-execution-result.md implies determinism that the code lacks.
Mismatch Type: FAKE_DEMO (partial)
Assessment: WARNING
Severity: MEDIUM
Notes: Not fabricated logic — values are genuinely random from rand.Float64(). But pinning specific random outputs as "results" is misleading. FullJitter 8ms vs documented 99ms is plausible (random in [0,100]).

## Finding 6

Location: README.md:51-65 vs engineering/03-execution-result.md commands.
Claimed: go build SUCCESS (03-report doesn't show), go test -v (pass), go test -race (pass), go run ./cmd/demo (listed).
Real execution just performed:
  - go build ./... → EXIT:0 ✓
  - go vet ./... → EXIT:0 ✓
  - go test -v -count=1 ./... → all pass ✓
  - go test -race -count=1 ./... → all pass ✓
  - go run ./cmd/demo → produces output ✓
Assessment: PASS
Notes: README commands all work in real environment.

## Finding 7

Location: README.md Features section vs actual file list.
Claimed files include "registry.go", "queue.go", "middleware.go", "backoff.go", "bucket.go" — all present.
Claimed "LeakyBucket" — present in bucket.go.
Claimed per-tenant "key isolation guarding against CGNAT" — registry.go.
Assessment: PASS

## Finding 8

Location: README.md:44 - Bounded Queue claims "fast ErrQueueFull when buffer capacity is reached" and "preventing queue age degradation and memory exhaustion."
Code: queue.go TrySubmit select-default ✓ fast drop ✓.
Assessment: PASS