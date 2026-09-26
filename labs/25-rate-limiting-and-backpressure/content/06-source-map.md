# Source Map

## Rate Limiting: Token Bucket

Research:
- research/05-report.md (Finding 1)
- research/02-sources.md (Source 1: RFC 6585, Source 3: Token Bucket Wikipedia)

Implementation:
- internal/ratelimit/bucket.go (TokenBucket struct, AllowN, RetryAfterSeconds)

Tests:
- internal/ratelimit/bucket_test.go (TestTokenBucket_BurstAndRefill, TestTokenBucket_ConcurrencyRace, TestTokenBucket_RetryAfterSeconds)

Demo:
- cmd/demo/main.go (Phase 1: Token Bucket Burst & Rate Limiting)

---

## Rate Limiting: Leaky Bucket

Research:
- research/05-report.md (Finding 1, Notes)
- research/02-sources.md (Source 6: Leaky Bucket Wikipedia)

Implementation:
- internal/ratelimit/bucket.go (LeakyBucket struct, Allow)

Tests:
- internal/ratelimit/bucket_test.go (TestLeakyBucket_LeakRate)

Demo:
- cmd/demo/main.go (Phase 2: Leaky Bucket Traffic Smoothing)

---

## HTTP 429 Middleware (RFC 6585)

Research:
- research/05-report.md (Finding 2)
- research/02-sources.md (Source 1: RFC 6585)

Implementation:
- internal/httputil/middleware.go (RateLimitMiddleware, RateLimitResponse)

Tests:
- internal/httputil/middleware_test.go (TestRateLimitMiddleware_RFC6585)

---

## Multi-tenant Registry (Tenant Isolation)

Research:
- research/05-report.md (Finding 5)
- research/02-sources.md (Source 5: Datacenter Traffic Control paper)

Implementation:
- internal/ratelimit/registry.go (Registry struct, Get method)

Tests:
- internal/ratelimit/bucket_test.go (TestRegistry_TenantIsolation)

---

## Bounded Queue Backpressure

Research:
- research/05-report.md (Finding 3, Finding 7)
- research/02-sources.md (Source 9: Little's Law Wikipedia)

Implementation:
- internal/backpressure/queue.go (BoundedQueue, TrySubmit, ErrQueueFull)

Tests:
- internal/backpressure/queue_test.go (TestBoundedQueue_RejectionUnderLoad, TestBoundedQueue_ConcurrencySafety)

Demo:
- cmd/demo/main.go (Phase 3: Bounded Queue Backpressure)

---

## Exponential Backoff with Jitter (AWS)

Research:
- research/05-report.md (Finding 4)
- research/02-sources.md (Source 10: AWS Architecture Blog)

Implementation:
- internal/retry/backoff.go (ComputeBackoff, 4 strategy variants)

Tests:
- internal/retry/backoff_test.go (TestComputeBackoff_Bounds, TestDecorrelatedJitter_Bounds)

Demo:
- cmd/demo/main.go (Phase 4: AWS Retry Backoff Strategies)

---

## Little's Law & Capacity Planning

Research:
- research/05-report.md (Finding 3, Finding 7)
- research/02-sources.md (Source 9: Little's Law Wikipedia)

Notes:
- Contoh perhitungan backlog 5M / 2000/sec = 2500 detik di dalam research/05-report.md
- Demo CLI tidak menjalankan contoh ini (hanya visualisasi queue bounded)

---

## CGNAT / IP Limitation (RFC 6598)

Research:
- research/05-report.md (Finding 2 Notes)
- research/02-sources.md (tidak terdaftar secara eksplisit, dikutip dari lab document)

Implementation:
- internal/httputil/middleware.go (ekstrak X-API-Key, fallback "anonymous")

---

## Reactive Streams / Backpressure Standard

Research:
- research/05-report.md (Finding 6)
- research/02-sources.md (Source 8: Reactive Streams Wikipedia)

Notes:
- Konsep referensi; implementasi lab menggunakan bounded queue Go channel, bukan Reactive Streams API

---

## Research Gaps / Caveats (Preserved from Audits)

Research-Audit:
- research-audit/06-gaps.md (Gap 1: 5/10 sources Wikipedia; Gap 2: Missing empirical benchmarks)
- research-audit/07-verdict.md (APPROVED dengan 2 non-blocking issues)

Engineering-Audit:
- engineering-audit/06-verdict.md (APPROVED, 0 failures, 0 warnings, race detector PASS)

Engineering Notes:
- engineering/02-implementation-notes.md (Known Limitations: in-memory state only, no distributed sync)
- engineering/03-execution-result.md (Full test + race detector + demo output)

---

## File Inventory (Verified)

Source Code:
- internal/ratelimit/bucket.go (121 lines)
- internal/ratelimit/registry.go (37 lines)
- internal/ratelimit/bucket_test.go (98 lines)
- internal/backpressure/queue.go (80 lines)
- internal/backpressure/queue_test.go (68 lines)
- internal/retry/backoff.go (63 lines)
- internal/retry/backoff_test.go (49 lines)
- internal/httputil/middleware.go (40 lines)
- internal/httputil/middleware_test.go (43 lines)
- cmd/demo/main.go (65 lines)
- go.mod, README.md

Audit Artifacts:
- research-audit/01-audit-plan.md through 07-verdict.md
- engineering-audit/01-audit-plan.md through 06-verdict.md
- engineering/01-design.md through 03-execution-result.md