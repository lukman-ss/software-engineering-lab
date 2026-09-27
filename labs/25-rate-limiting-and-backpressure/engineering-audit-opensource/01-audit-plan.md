# Engineering Audit Plan

Target Lab: `labs/25-rate-limiting-and-backpressure` — Rate Limiting & Backpressure

Implementation Files:
- internal/ratelimit/bucket.go — TokenBucket, LeakyBucket
- internal/ratelimit/registry.go — per-tenant Registry
- internal/backpressure/queue.go — BoundedQueue (TrySubmit)
- internal/retry/backoff.go — AWS Jitter backoff strategies
- internal/httputil/middleware.go — HTTP 429 + Retry-After middleware
- cmd/demo/main.go — CLI demo driver

Tests:
- internal/ratelimit/bucket_test.go
- internal/backpressure/queue_test.go
- internal/retry/backoff_test.go
- internal/httputil/middleware_test.go

Executable/Demo:
- `go run ./cmd/demo` (cmd/demo/main.go)

Approved Research Inputs:
- research/01-plan.md, 02-sources.md, 03-evidence.md, 04-contradictions.md, 05-report.md
(content/01-content-brief.md lists Verified Behaviors — used as claims baseline)

Main Claims To Verify:
1. Token bucket allows burst up to capacity B, then rejects until refill.
2. Token bucket computes accurate `RetryAfterSeconds` when empty.
3. Leaky bucket rejects burst when capacity full, allows after leak.
4. Registry gives isolated per-tenant quotas (tenant A ≠ tenant B).
5. BoundedQueue.TrySubmit drops immediately with ErrQueueFull under load.
6. Token bucket + bounded queue race-free under 50 goroutines (race detector clean).
7. Full Jitter backoff ∈ [0, min(cap, base·2^attempt)].
8. HTTP middleware returns 429 + Retry-After under tenant overload.
9. Demo CLI displays all behaviors live.

Commands To Run:
- `go vet ./...`
- `go test -count=1 -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Time-based tests (sleep/retry thresholds) can be flaky but not fabricated.
- `RetryAfterSeconds` rounding returns min 1 for any positive fractional need — verify claim vs code.
- LeakyBucket.Allow condition `water+1.0 <= capacity` permits equal fill; check boundary.
- Full Jitter bound uses `<= cap` inclusive; verify strict vs inclusive.
