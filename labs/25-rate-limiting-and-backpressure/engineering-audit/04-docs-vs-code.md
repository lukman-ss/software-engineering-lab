# Docs vs Code Audit

Target Lab: `labs/25-rate-limiting-and-backpressure`

## Comparison Matrix

| Claim / Topic | Documentation (`README.md`, `engineering/`) | Codebase (`internal/*`, `cmd/*`) | Audit Assessment |
|---|---|---|---|
| Token Bucket Burst & Refill | Documented in README.md & 01-design.md | Implemented in `internal/ratelimit/bucket.go:TokenBucket` | MATCH |
| Leaky Bucket Traffic Smoothing | Documented in README.md & 01-design.md | Implemented in `internal/ratelimit/bucket.go:LeakyBucket` | MATCH |
| Multi-tenant Key Registry | Documented in README.md & 01-design.md (RFC 6598 CGNAT avoidance) | Implemented in `internal/ratelimit/registry.go:Registry` | MATCH |
| Bounded Queue Backpressure | Documented in README.md & 01-design.md (`TrySubmit`, `ErrQueueFull`) | Implemented in `internal/backpressure/queue.go:BoundedQueue` | MATCH |
| AWS Retry Jitter Backoff | Full Jitter, Equal Jitter, No Jitter, Decorrelated Jitter | Implemented in `internal/retry/backoff.go:ComputeBackoff` | MATCH |
| RFC 6585 HTTP 429 Middleware | Status 429, `Retry-After` header, JSON payload | Implemented in `internal/httputil/middleware.go:RateLimitMiddleware` | MATCH |
| Demo Execution Output | Documented in `engineering/03-execution-result.md` | Verified via real run of `cmd/demo/main.go` | MATCH |

## Discrepancies Found

None. The documentation accurately reflects the code structure, algorithms, and execution output.
