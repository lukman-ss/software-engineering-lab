# Rate Limiting & Backpressure Engineering Lab

Target Lab: `labs/25-rate-limiting-and-backpressure`

## Overview

This lab provides a Go 1.22+ implementation demonstrating rate limiting, bounded queue backpressure, HTTP 429 response handling (RFC 6585), and exponential backoff with AWS Full Jitter.

## Structure

```text
labs/25-rate-limiting-and-backpressure/
├── cmd/
│   └── demo/
│       └── main.go
├── internal/
│   ├── backpressure/
│   │   ├── queue.go
│   │   └── queue_test.go
│   ├── httputil/
│   │   ├── middleware.go
│   │   └── middleware_test.go
│   ├── ratelimit/
│   │   ├── bucket.go
│   │   ├── bucket_test.go
│   │   └── registry.go
│   └── retry/
│       ├── backoff.go
│       └── backoff_test.go
├── engineering/
│   ├── 01-design.md
│   ├── 02-implementation-notes.md
│   └── 03-execution-result.md
├── engineering-revision/
│   ├── 01-revision-plan.md
│   ├── 02-changes-made.md
│   └── 03-revision-result.md
├── go.mod
└── README.md
```

## Features Implemented

1. **Token Bucket & Leaky Bucket Algorithms**:
   - `TokenBucket`: Allows bursts up to capacity $B$, refilling continuously at rate $R$.
   - `LeakyBucket`: Enforces constant drain rate $R$, rejecting bursts when water reaches capacity.
   - Per-tenant registry key isolation guarding against CGNAT IP collisions (RFC 6598).
2. **Bounded Queue Backpressure**:
   - Non-blocking submission (`TrySubmit`) returning fast `ErrQueueFull` when buffer capacity is reached, preventing queue age degradation and memory exhaustion.
3. **AWS Jitter Retry Strategies**:
   - Full Jitter, Equal Jitter, No Jitter, and Decorrelated Jitter algorithms matching AWS Architecture specifications (Marc Brooker).
4. **HTTP 429 Middleware**:
   - Enforces rate limits returning standard RFC 6585 `429 Too Many Requests` response with `Retry-After` header.

## Execution & Testing

Run unit tests:
```bash
go test ./...
```

Run race detector:
```bash
go test -race ./...
```

Run demo:
```bash
go run ./cmd/demo
```
