# Docs vs Code Audit

Target Lab: `labs/25-rate-limiting-and-backpressure`

## Comparison Matrix

| Item | Documentation / Design | Implementation | Assessment |
|---|---|---|---|
| Token Bucket & Leaky Bucket | Burst accommodation up to $B$, constant refill $R$, leaky smoothing | `internal/ratelimit/bucket.go` | PASS |
| Bounded Queue Backpressure | Fast drop `ErrQueueFull` on queue saturation | `internal/backpressure/queue.go` | PASS |
| AWS Jitter Retries | Full Jitter, Equal Jitter, No Jitter, Decorrelated Jitter | `internal/retry/backoff.go` | PASS |
| HTTP 429 Middleware | RFC 6585 status 429 and `Retry-After` header | `internal/httputil/middleware.go` | PASS |
| Tenant Quota Isolation | RFC 6598 tenant key isolation over CGNAT IP limits | `internal/ratelimit/registry.go` | PASS |
| README Instructions | Commands `go test ./...`, `go test -race ./...`, `go run ./cmd/demo` | Fully working & matching | PASS |

## Discrepancy Findings

No blocking mismatches found between design, documentation, and implementation.
