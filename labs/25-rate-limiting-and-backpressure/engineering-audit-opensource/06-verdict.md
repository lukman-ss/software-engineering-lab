# Engineering Audit Verdict

Target Lab: labs/25-rate-limiting-and-backpressure
Audit Date: 2026-09-26

## Summary

Audit scope: implementation + tests only (research excluded per pipeline override). No code was modified during the audit (one temporary probe test was created to confirm a defect, executed, then deleted; the working tree is restored).

Code Files Reviewed: 6
- internal/ratelimit/bucket.go, internal/ratelimit/registry.go
- internal/backpressure/queue.go
- internal/httputil/middleware.go
- internal/retry/backoff.go
- cmd/demo/main.go

Tests Reviewed: 4 files (11 test functions total)

Commands Executed (results recorded in 02–03):
- `go build ./...`                 -> exit 0 (SUCCESS)
- `go vet ./...`                   -> exit 0 (SUCCESS)
- `go test -v -count=1 ./...`      -> all packages PASS
- `go test -race ./...`            -> all packages PASS (see caveat)
- `go run ./cmd/demo` (×3)         -> exit 0 (SUCCESS); queue accept/reject split varies 3–4 vs 2–3 (timing-dependent)

Failures: 0 tests fail.
Warnings: documented in 04/05.

## Quality Gates

| Gate | Result | Notes |
|---|---|---|
| Compilation | PASS | `go build ./...` exit 0 |
| Tests | PASS | all 11 tests green |
| Race Detector | PASS * | green on shipped tests; the concurrent submit+Stop path is untested and panics (see Blocking) |
| Demo | PASS | runs cleanly; non-deterministic queue counts are timing, not fabrication |
| Research Alignment | NOT_AUDITED | excluded per pipeline override |
| Documentation Accuracy | WARNING | false "simulated time" claim; wrong test filenames; execution-result over-specifies timing-dependent numbers |

\* The race detector is green only because no test exercises concurrent `TrySubmit`+`Stop`; the defect is a logic panic, not a data race, so `-race` cannot surface it.

## Blocking Issues

1. `BoundedQueue.TrySubmit` (internal/backpressure/queue.go:66-81) panics with `send on closed channel` when invoked concurrently with `Stop()` (internal/backpressure/queue.go:87-94), which closes `bq.queue`. Reproduced by a temporary audit probe:
   ```
   panic: send on closed channel
   .../queue.go:71
   ```
   The design (§Concurrency safety / "thread-safe" success criterion, §BoundedQueue "rejecting jobs when full") and README ("Non-blocking submission ... preventing ... memory exhaustion") imply safe concurrent use; the implementation violates this on the shutdown path. Severity: HIGH. The shipped 3-test backpressure suite never tests this path, so CI is green while production shutdown would crash.

## Non-Blocking Issues

1. DOC_CODE_MISMATCH (MEDIUM): design/implementation-notes claim "deterministic simulated time helpers in tests"; actual tests use real `time.Sleep`.
2. DOC_CODE_MISMATCH (LOW): design §Test Strategy lists four test filenames that don't exist.
3. DOC_CODE_MISMATCH (MEDIUM): `03-execution-result.md` records fixed queue stats (Accepted=3/Rejected=3) that vary run-to-run (3–4 accepted).
4. MISSING_TEST (MEDIUM): middleware test covers only one tenant's 200→429; no anonymous fallback, body, Retry-After value, or multi-tenant coverage.
5. MISSING_TEST (MEDIUM): backoff tests check bounds only; no exact-formula assertions.
6. MISSING_EDGE_CASE (MEDIUM): `Registry` never evicts per-tenant buckets → unbounded memory growth under high tenant cardinality (not in known limitations).
7. MISSING_EDGE_CASE (LOW): `RetryAfterSeconds` divides by `refillRate`; no guard for zero refill rate.
8. UNVERIFIED_RESULT (LOW): recorded jitter millisecond values are random per run.

## Required Revisions

1. **Fix the shutdown panic (HIGH):** prevent `TrySubmit` from sending on a closed channel. Recommended approaches: (a) add an atomic "stopping/closing" flag set before `close(queue)` and re-checked inside the `select` success branch (double-check + recover or re-loop), or (b) have `Stop` drain-and-close with a guard, or (c) return `ErrQueueStopped` from the send branch if closed. Any accepted fix must keep the non-blocking fast-fail contract.
2. **Add a regression test** that runs many concurrent `TrySubmit` goroutines while another goroutine calls `Stop`, asserting no panic and all results are one of {accepted, rejected, ErrQueueFull, ErrQueueStopped}.
3. **Correct documentation (MEDIUM):** rewrite the false "simulated time helpers" claim; update test-strategy filenames to the actual `bucket_test.go`/`queue_test.go`/`backoff_test.go`/`middleware_test.go`; annotate the demo queue stats as timing-dependent rather than fixed.

## Final Status

NEEDS_REVISION

Rationale: compilation, tests, and the race detector all pass, the demo runs, and there is no fabricated benchmark/result. However, a confirmed HIGH-severity defect — a panic in `BoundedQueue` under concurrent submit+Stop, a path the design explicitly covers ("concurrency safety") but tests omit — is unresolved, and the documentation contains a direct false claim about test methodology. Per the audit gate rules (APPROVED requires no unresolved HIGH/CRITICAL issues and README/code alignment), the lab cannot be approved as-is.

Recommendation: implement revision #1 + #2 (fix + regression test), correct docs per #3, then re-run `go test -race ./...` and the demo before re-submission. Re-audit after revisions.