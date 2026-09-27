# Documentation vs Code Audit Report

## Target Lab
`labs/25-rate-limiting-and-backpressure`

## Comparison Matrix

| Component / Claim | README.md Claim | Code & Test Implementation | Status |
| :--- | :--- | :--- | :--- |
| **Token Bucket** | Burst to capacity $B$, continuous refill rate $R$ | `TokenBucket` in `internal/ratelimit/bucket.go`, tested in `bucket_test.go` | PASS (MATCH) |
| **Leaky Bucket** | Constant drain rate $R$, burst rejection when full | `LeakyBucket` in `internal/ratelimit/bucket.go`, tested in `bucket_test.go` | PASS (MATCH) |
| **Tenant Registry** | Per-tenant rate limiter key isolation (RFC 6598) | `Registry` in `internal/ratelimit/registry.go`, tested in `bucket_test.go` | PASS (MATCH) |
| **Bounded Backpressure** | Non-blocking `TrySubmit` with fast `ErrQueueFull` drop | `BoundedQueue` in `internal/backpressure/queue.go`, tested in `queue_test.go` | PASS (MATCH) |
| **HTTP Middleware** | Status 429 Too Many Requests with `Retry-After` header | `RateLimitMiddleware` in `internal/httputil/middleware.go`, tested in `middleware_test.go` | PASS (MATCH) |
| **AWS Jitter Retries** | Full, Equal, NoJitter, Decorrelated Jitter algorithms | `ComputeBackoff` in `internal/retry/backoff.go`, tested in `backoff_test.go` | PASS (MATCH) |
| **Demo Execution** | `go run ./cmd/demo` runnable output matching claims | `cmd/demo/main.go` executes all 4 components cleanly | PASS (MATCH) |

## Discrepancy Checks

- **DOC_CODE_MISMATCH**: None detected. README directory layout, APIs, commands, and behavior explanations match source files exactly.
- **TEST_CLAIM_MISMATCH**: None detected. All claimed edge cases (exhaustion, refill, drop on overflow, post-shutdown submission, retry headers) have corresponding passing tests.
- **RESEARCH_IMPLEMENTATION_MISMATCH**: None detected. Implements core research tenants (Token Bucket, Leaky Bucket, RFC 6585, RFC 6598, Marc Brooker AWS jitter backoff, bounded queue load shedding).
