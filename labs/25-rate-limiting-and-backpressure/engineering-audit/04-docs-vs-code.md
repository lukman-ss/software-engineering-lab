# Docs vs Code Analysis

## Verification Matrix

| Claim in README / Design Doc | Implementation Location | Test Verification | Verdict |
|---|---|---|---|
| Token Bucket with burst allowance and continuous refill | `internal/ratelimit/bucket.go:8-79` | `internal/ratelimit/bucket_test.go:9` | MATCH |
| Leaky Bucket traffic smoothing with constant leak rate | `internal/ratelimit/bucket.go:81-121` | `internal/ratelimit/bucket_test.go:32` | MATCH |
| Tenant key isolation guarding against CGNAT IP collisions | `internal/ratelimit/registry.go:6-37` | `internal/ratelimit/bucket_test.go:51` | MATCH |
| Bounded queue non-blocking submit with immediate rejection | `internal/backpressure/queue.go:60-70` | `internal/backpressure/queue_test.go:10` | MATCH |
| AWS exponential backoff with Full / Equal / Decorrelated Jitter | `internal/retry/backoff.go:24-63` | `internal/retry/backoff_test.go:8` | MATCH |
| HTTP RFC 6585 status 429 response with Retry-After header | `internal/httputil/middleware.go:17-40` | `internal/httputil/middleware_test.go:11` | MATCH |
| Interactive CLI Demo matching actual execution output | `cmd/demo/main.go` | Run verified via `go run ./cmd/demo` | MATCH |

## Findings

- `DOC_CODE_MISMATCH`: None detected.
- `TEST_CLAIM_MISMATCH`: None detected.
- `RESEARCH_IMPLEMENTATION_MISMATCH`: None detected.

Documentation accurately reflects the code and runtime output.
