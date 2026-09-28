# Source Map

Pemetaan tiap bagian article ke sumber riset, sumber eksternal, file implementasi, dan test yang mendukungnya.

## Problem

Research:
- `research/05-report.md` — Executive Summary (risk cascading failure tanpa protective mechanisms)

Research Sources:
- [10] Google SRE Workbook: Handling Overload (https://landing.google.com/sre/sre-book/chapters/handling-overload/)
- [7] Stripe Engineering Blog: Scaling your API with rate limiters (https://stripe.com/blog/rate-limiters)

## Why This Matters

Research:
- `research/05-report.md` — Finding 6 (Multi-Layer Rate Limiting Approach)
- `research/05-report.md` — Finding 8 (Google SRE Per-Customer Limits)

Research Sources:
- [7] Stripe Engineering Blog: Scaling your API with rate limiters
- [10] Google SRE Workbook: Handling Overload

## Mental Model

Research:
- `research/05-report.md` — Finding 1 (Token Bucket vs Leaky Bucket), Finding 2 (Backpressure vs Rate Limiting), Finding 3 (Exponential Backoff with Jitter)

Research Sources:
- [1] Token bucket — Wikipedia
- [12] Leaky bucket — Wikipedia
- [8] AWS Architecture Blog: Exponential Backoff And Jitter

Implementation:
- `internal/ratelimit/bucket.go` — TokenBucket, LeakyBucket
- `internal/backpressure/queue.go` — BoundedQueue
- `internal/retry/backoff.go` — ComputeBackoff

## Core Concept: Token Bucket vs. Leaky Bucket

Research:
- `research/05-report.md` — Finding 1
- `research/03-evidence.md` — Evidence 1 (Token Bucket Properties), Evidence 2 (Leaky Bucket Meter vs Queue)

Research Sources:
- [1] Token bucket — Wikipedia
- [12] Leaky bucket — Wikipedia
- [5] NGINX limit_req_module

Implementation:
- `internal/ratelimit/bucket.go` — `NewTokenBucket`, `Allow()`, `AllowN()`, `NewLeakyBucket`, `Allow()`

## Failure Scenario

Research:
- `research/05-report.md` — Finding 5 (Little's Law), Executive Summary

Research Sources:
- [3] Little's law — Wikipedia
- [6] RFC 6585: Additional HTTP Status Codes

Calculation:
- 5,000,000 items ÷ 2,000 items/sec = 2,500 seconds ≈ 41 minutes 40 seconds (verified against Little's Law)

## Architecture

Engineering:
- `engineering/01-design.md` — Architecture diagram, Components section

Implementation:
- `internal/httputil/middleware.go` — HTTP Middleware / Tenant Extractor
- `internal/ratelimit/bucket.go` + `internal/ratelimit/registry.go` — Token Bucket + Registry
- `internal/backpressure/queue.go` — Bounded Queue
- `internal/retry/backoff.go` — Jitter Backoff
- `cmd/demo/main.go` — CLI demo

## Implementation

Engineering:
- `engineering/01-design.md` — Implementation Decisions (pure stdlib, monotonic time, mutex vs CAS)
- `engineering/02-implementation-notes.md` — Core Design Decisions, Known Limitations

Implementation Files:
- `internal/ratelimit/bucket.go` — TokenBucket, LeakyBucket
- `internal/ratelimit/registry.go` — Registry (per-tenant)
- `internal/backpressure/queue.go` — BoundedQueue
- `internal/retry/backoff.go` — ComputeBackoff (FullJitter, EqualJitter, DecorrelatedJitter, NoJitter)
- `internal/httputil/middleware.go` — RateLimitMiddleware

## Code Walkthrough

Implementation Files:
- `cmd/demo/main.go` — Demo sections 1-4
- `internal/httputil/middleware.go:17-40` — Middleware 429 RFC 6585 + Retry-After

## What the Tests Prove

Engineering:
- `engineering/03-execution-result.md` — Test results (15/15 pass, race detector clean)

Test Files:
- `internal/ratelimit/bucket_test.go` — TestTokenBucket_BurstAndRefill, TestLeakyBucket_LeakRate, TestRegistry_TenantIsolation, TestTokenBucket_RetryAfterSeconds, TestTokenBucket_ConcurrencyRace, TestLeakyBucket_ConcurrencyRace, TestRegistry_ConcurrentSameKeyGet, TestTokenBucket_EdgeCases
- `internal/backpressure/queue_test.go` — TestBoundedQueue_RejectionUnderLoad, TestBoundedQueue_ConcurrencySafety, TestBoundedQueue_SubmitAfterStop, TestBoundedQueue_ConcurrentStopAndSubmit
- `internal/retry/backoff_test.go` — TestComputeBackoff_Bounds, TestDecorrelatedJitter_Bounds, TestComputeBackoff_UnknownStrategy
- `internal/httputil/middleware_test.go` — TestRateLimitMiddleware_RFC6585, TestRateLimitMiddleware_AnonymousFallbackAndBody

## Recovery / Rollback

Implementation:
- `internal/backpressure/queue.go:91-103` — `Stop()` graceful shutdown
- `internal/ratelimit/bucket.go:54-79` — `RetryAfterSeconds()`

## Production Considerations

Research:
- `research/05-report.md` — Finding 7 (Redis as Backend), Finding 9 (Client-Side Throttling and Retry Budget)

Research Sources:
- [11] Redis rate limiter documentation (https://redis.io/docs/latest/develop/use-cases/rate-limiter/)
- [9] AWS SDK Documentation: Retry behavior
- [10] Google SRE Workbook: Handling Overload

Research Revision:
- `research-revision/03-revision-result.md` — AWS SDK constants contextualized as implementation-specific defaults

Engineering Revision:
- `engineering-revision/03-revision-result.md` — Resolved all gaps (high/medium/low)

## Common Mistakes

Research:
- `research/05-report.md` — Finding 2, Finding 8, Finding 9

Research Sources:
- [6] RFC 6585
- [10] Google SRE Workbook: retry budgets
- Engineering revision: CGNAT RFC 6598 per-tenant key limiting

## Case Study: Stripe

Research:
- `research/05-report.md` — Finding 6

Research Sources:
- [7] Stripe Engineering Blog: Scaling your API with rate limiters

## Case Study: Google SRE

Research:
- `research/05-report.md` — Finding 8, Finding 9

Research Sources:
- [10] Google SRE Workbook: Handling Overload

## Audit Status

Research Audit Verdict: APPROVED (`research-audit/07-verdict.md`)
Engineering Audit Verdict: APPROVED (`engineering-audit/06-verdict.md`, `engineering-audit-opensource/06-verdict.md`)
Research Revision: READY_FOR_RESEARCH_REAUDIT → re-audit completed, APPROVED
Engineering Revision: READY_FOR_ENGINEERING_REAUDIT → re-audit completed, APPROVED
