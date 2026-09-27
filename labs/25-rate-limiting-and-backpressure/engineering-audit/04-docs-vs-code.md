# Documentation vs Code Audit

## Comparison Matrix

| Component / Claim | Documentation (README / Design) | Code Implementation | Status |
|---|---|---|---|
| Token Bucket | Capacity $B$, continuous refill rate $R$ | `internal/ratelimit/bucket.go:8-46` | MATCH |
| Leaky Bucket | Capacity $B$, constant drain rate $R$ | `internal/ratelimit/bucket.go:81-121` | MATCH |
| Tenant Isolation | Multi-tenant key registry avoiding CGNAT collision | `internal/ratelimit/registry.go:6-37` | MATCH |
| Backpressure Queue | Non-blocking `TrySubmit` dropping excess jobs with `ErrQueueFull` | `internal/backpressure/queue.go:18-85` | MATCH |
| AWS Jitter Retries | Full Jitter, Equal Jitter, No Jitter, Decorrelated Jitter | `internal/retry/backoff.go:9-63` | MATCH |
| HTTP 429 & Retry-After | RFC 6585 compliance with `Retry-After` header | `internal/httputil/middleware.go:17-40` | MATCH |
| CLI Demo | Demonstrates all 4 components sequentially | `cmd/demo/main.go:13-65` | MATCH |

## Discrepancies Found
- None. Code structure, exported identifiers, and documented behaviors match 1:1.
