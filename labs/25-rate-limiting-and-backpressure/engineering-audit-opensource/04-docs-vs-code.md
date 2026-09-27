# Docs vs Code

Scope note: per pipeline override, research/content is NOT audited. Comparisons below cover README.md, engineering/01-design.md, 02-implementation-notes.md, 03-execution-result.md against code/tests/demo behavior only.

## Finding 1 — STALE_TEST_RECORD (03-execution-result.md)

- engineering/03-execution-result.md test transcript lists 9 tests: RejectionUnderLoad, ConcurrencySafety, RFC6585, BurstAndRefill, LeakRate, TenantIsolation, ConcurrencyRace, ComputeBackoff_Bounds, DecorrelatedJitter_Bounds.
- Actual suite has 11 tests: the record omits TestBoundedQueue_SubmitAfterStop and TestTokenBucket_RetryAfterSeconds, both of which exist and pass.
- Type: TEST_CLAIM_MISMATCH. Severity: MEDIUM (record understates coverage; not a fabrication, but the transcript is not a faithful log of the current suite).

## Finding 2 — NON_REPRODUCIBLE_DEMO_TRANSCRIPT (03-execution-result.md vs demo)

- 03-execution-result.md §3 records: "Job #1..3 ACCEPTED, #4..6 REJECTED; Stats: Accepted=3, Rejected=3, Processed=1".
- Actual run 2026-09-27: "Job #1..3 ACCEPTED, #4 REJECTED, #5 ACCEPTED, #6 REJECTED; Stats: Accepted=4, Rejected=2, Processed=1". The demo submits into a 3-cap buffer with a 1-worker pool draining at 50ms/job; which submits land before a worker pop is timing-dependent. Jitter values also differ run to run by design (record shows FullJitter 99/109/111/125ms; actual 96/95/25/701ms).
- Type: UNVERIFIED_RESULT (transcript presented as fixed output for a nondeterministic program). Severity: MEDIUM. The demo is REAL (runs, exit 0, behavior demonstrated); the defect is presenting one timing-dependent sample as canonical output, not a fake demo.

## Finding 3 — README STRUCTURE DRIFT (minor)

- README §Structure lists engineering/01-design.md, 02-implementation-notes.md, 03-execution-result.md. Directory additionally contains engineering-revision/ (01/02/03). README does not mention engineering-revision/ or engineering-audit outputs.
- Type: DOC_CODE_MISMATCH. Severity: LOW.

## Finding 4 — README FEATURE CLAIMS vs CODE (checked, no mismatch)

- README claims TokenBucket bursts to B with continuous refill; LeakyBucket constant drain R; per-tenant registry; TrySubmit fast ErrQueueFull; Full/Equal/No/Decorrelated jitter per AWS; 429 + Retry-After. All match code behavior observed. No DOC_CODE_MISMATCH on functional claims.
- README "Run demo: go run ./cmd/demo" verified working.

## Finding 5 — DESIGN SUCCESS-CRITERIA OVERCLAIM

- 01-design.md success criteria: "TokenBucket and LeakyBucket accurately track and reject ... under concurrent access" — only TokenBucket has a concurrency test; LeakyBucket has none. "Concurrency test passes with Go race detector clean" — true for tested paths but concurrent shutdown path panics (code-audit Finding 1). So the design's concurrency claim is broader than what tests prove.
- Type: IMPLEMENTATION_OVERCLAIM (scoped to concurrency/shutdown). Severity: HIGH (inherits code-audit Finding 1).

## No FAKE_DEMO / FAKE_BENCHMARK

- Demo executes for real, output above is genuine. No benchmarks exist; none claimed.
