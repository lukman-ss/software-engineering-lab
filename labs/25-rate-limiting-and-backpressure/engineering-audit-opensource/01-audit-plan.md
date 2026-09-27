# Engineering Audit Plan

Target Lab: labs/25-rate-limiting-and-backpressure

Implementation Files:
- internal/ratelimit/bucket.go (TokenBucket, LeakyBucket)
- internal/ratelimit/registry.go (Registry multi-tenant)
- internal/backpressure/queue.go (BoundedQueue)
- internal/retry/backoff.go (AWS jitter strategies)
- internal/httputil/middleware.go (HTTP 429 middleware)
- cmd/demo/main.go (CLI demo)

Tests:
- internal/ratelimit/bucket_test.go
- internal/backpressure/queue_test.go
- internal/httputil/middleware_test.go
- internal/retry/backoff_test.go

Executable/Demo: `go run ./cmd/demo`

Approved Research Inputs: (out of scope for this pipeline stage — implementation & tests only)

Main Claims To Verify:
1. Token/Leaky bucket enforce capacity + refill/leak rates under concurrency.
2. Bounded queue fast-rejects (non-blocking select-default) when capacity full.
3. Retry backoff adheres to AWS Full/Equal/No/Decorrelated jitter bounds.
4. HTTP middleware emits RFC 6585 429 + Retry-After header.
5. Race detector clean: `go test -race ./...`.
6. Demo output reflects real behavior (not fabricated).

Commands To Run:
- go build ./cmd/... ./internal/...
- go test ./cmd/... ./internal/...
- go test -race ./cmd/... ./internal/...
- go run ./cmd/demo

Primary Risks:
- `go test ./...` matched no packages (module root has no nested dirs at repo root — using explicit `./cmd/... ./internal/...`).
- Demo Stats output (Accepted=4, Rejected=2, Processed=1) differs from engineering notes (Accepted=3, Rejected=3, Processed=1) — race in demo stats read vs worker.
- LeakyBucket capacity check uses `+1.0 <= capacity`; boundary semantics subtle.
- Retry jitter is random — bounds must hold probabilistically, not exactly.
