# Docs vs Code Audit

## Overview

This audit compares the claims made in `README.md`, `engineering/01-design.md`, and `engineering/02-implementation-notes.md` against the actual Go codebase and executable output.

## Comparison Table

| Feature / Claim | Documented Claim | Code Implementation | Status |
| --- | --- | --- | --- |
| Token Bucket | Bursts up to $B$, continuous refill rate $R$ | `TokenBucket` in `internal/ratelimit/bucket.go` | PASS |
| Leaky Bucket | Constant drain rate $R$, rejects bursts when full | `LeakyBucket` in `internal/ratelimit/bucket.go` | PASS |
| Multi-tenant Isolation | Per-tenant key registry to avoid CGNAT IP collision (RFC 6598) | `Registry` in `internal/ratelimit/registry.go` | PASS |
| Bounded Queue Backpressure | Non-blocking `TrySubmit` returning fast `ErrQueueFull` | `BoundedQueue.TrySubmit` in `internal/backpressure/queue.go` | PASS |
| AWS Jitter Retries | Full, Equal, No, and Decorrelated Jitter (Marc Brooker) | `ComputeBackoff` in `internal/retry/backoff.go` | PASS |
| HTTP 429 Middleware | Standard RFC 6585 response with `Retry-After` header | `RateLimitMiddleware` in `internal/httputil/middleware.go` | PASS |
| Clean Race Execution | Zero race condition warnings | `go test -race ./...` passes cleanly | PASS |
| Demo Output | Shows rate limits, backpressure drops, and retry jitter | `cmd/demo/main.go` runs without errors | PASS |

## Discrepancies Found

None. README and design document accurately describe all implemented structures, package locations, execution commands, and behavior.
