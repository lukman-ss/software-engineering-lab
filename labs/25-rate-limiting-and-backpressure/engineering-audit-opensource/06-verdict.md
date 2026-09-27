# Engineering Audit Verdict

Target Lab: labs/25-rate-limiting-and-backpressure
Audit Date: 2026-09-27

## Summary

Code Files Reviewed: internal/ratelimit/bucket.go, internal/ratelimit/registry.go, internal/backpressure/queue.go, internal/httputil/middleware.go, internal/retry/backoff.go, cmd/demo/main.go (6 .go files)

Tests Reviewed: 11 unit tests (bucket_test.go:5, queue_test.go:3, middleware_test.go:1, backoff_test.go:2)

Commands Executed: go build ./..., go test -count=1 -v ./..., go test -race -count=1 ./..., go run ./cmd/demo, go vet ./..., gofmt -l ., plus a throwaway concurrent Stop+Submit repro (removed post-verification).

Failures: 1 (race-condition panic under concurrent shutdown: throwaway test confirmed send-on-closed-channel panic)

Warnings: 2 (03-execution-result.md stale test transcript; non-reproducible demo backpressure stats)

## Quality Gates

Compilation: PASS (go build ./... -> SUCCESS)
Tests: PASS (11/11 PASS, fresh count=1, exit 0)
Race Detector: PASS (no data races, -race count=1 -> all ok)
Demo: PASS (go run ./cmd/demo -> exit 0, real output)
Research Alignment: NOT_APPLICABLE (pipeline override: research/content not audited)
Documentation Accuracy: WARNING (README structure omits revision/ audit artifacts; 03-execution-result.md test transcript stale and presents nondeterministic demo output as canonical)

## Blocking Issues

1. Race condition in BoundedQueue: concurrent Stop + TrySubmit can panic with send-on-closed-channel. Not covered by existing tests; latent crash hazard under documented usage. (Severity: HIGH)
2. Design over-claim: concurrency-safety guarantee extends to a code path that panics (inherits from #1). (Severity: HIGH)

## Non-Blocking Issues

1. Test record stale: engineering/03-execution-result.md omits TestBoundedQueue_SubmitAfterStop and TestTokenBucket_RetryAfterSeconds from its transcript. (Severity: MEDIUM)
2. Demo non-reproducibility: backpressure section of 03-execution-result.md presents one timing-dependent outcome as canonical; jitter values differ by design. (Severity: MEDIUM)
3. Missing test: Registry concurrent same-key Get (double-checked locking) unproven. (Severity: LOW)
4. Missing test: LeakyBucket concurrency unproven. (Severity: LOW)
5. Missing test: Middleware assertions shallow (Retry-After value, body JSON, anonymous fallback). (Severity: LOW)
6. Missing edge cases: TokenBucket.AllowN>1, zero capacity/rate, negative input, unknown BackoffStrategy default. (Severity: LOW)
7. Documentation drift: README omits engineering-revision/ and engineering-audit-*/ directories from structure listing. (Severity: LOW)
8. Code formatting: 3 files flagged by gofmt (whitespace/alignment). (Severity: LOW)

## Required Revisions

1. Fix BoundedQueue.TrySubmit race: guard the channel send against a closed channel or re-check ctx.Done after the stopped check. Example patterns: use a drained flag, or move stopped check into the first select case, or use a mutex to serialize Stop vs submit.
2. Update engineering/03-execution-result.md: either replace with a fresh run labeled as example (not canon), add timing caveats, or narrow claims to functional behavior only (burst, leak, jitter spread, backpressure concept).
3. Augment test suite: add concurrent Stop+Submit test (expecting no panic), improve Middleware assertions, add Registry same-key Get concurrency test, etc.
4. Address minor items (missing edge cases, docs, fmt) at discretion.

## Final Status

NEEDS_REVISION

Reason: Unresolved HIGH severity issues exist (concurrent shutdown panic; design over-claim). APPROVED requires no unresolved HIGH/CRITICAL; APPROVED_WITH_WARNINGS is not defined to permit HIGH. The core implementation is correct and tested, but the identified defect is a real crash hazard under plausible usage and must be fixed before the engineering record can be trusted as defect-free.