# Docs vs Code Audit

Target Lab: `labs/25-rate-limiting-and-backpressure`

## Comparison Matrix

| Claim / Specification | Documentation Location | Code / Demo Location | Verdict | Notes |
|---|---|---|---|---|
| Token Bucket & Leaky Bucket Algorithms | `README.md:40-43`, `01-design.md:40-43` | `internal/ratelimit/bucket.go:8-121` | MATCH | Standard implementation matching design math. |
| Bounded Queue Backpressure | `README.md:44-45`, `01-design.md:44-45` | `internal/backpressure/queue.go:14-70` | MATCH | `TrySubmit` rejects non-blocking when capacity is reached. |
| AWS Retry Jitter Backoff | `README.md:46-47`, `01-design.md:46-47` | `internal/retry/backoff.go:24-62` | MATCH | Formulas match Marc Brooker's AWS specifications. |
| RFC 6585 HTTP 429 & `Retry-After` Header | `README.md:48-49`, `01-design.md:48-49` | `internal/httputil/middleware.go:17-40` | MATCH | Headers set correctly and status 429 returned. |
| Per-Tenant Registry (RFC 6598 CGNAT Avoidance) | `README.md:43`, `01-design.md:63` | `internal/ratelimit/registry.go:6-37` | MATCH | Key-based registry isolates tenant states. |
| Project Directory Structure | `README.md:11-36` | File system tree | MATCH | Structure matches README exactly. |

## Mismatch Analysis
- `DOC_CODE_MISMATCH`: None observed.
- `TEST_CLAIM_MISMATCH`: None observed.
- `RESEARCH_IMPLEMENTATION_MISMATCH`: None observed.
