# Gap Analysis

Target Lab: labs/25-rate-limiting-and-backpressure
Date: 2026-09-27
Allowed gap types as per specification.

## MISSING_TEST
- TokenBucket refill exactness (epsilon) near boundary times.
- LeakyBucket zero leakRate handling (division by zero in theory, but Allow uses subtraction only; no test).
- BoundedQueue Stop concurrent with TrySubmit (race condition possible).
- Retry unknown BackoffStrategy fallback behavior (implementation default to temp; not tested).
- HTTP middleware 503 path (never exercised, queue errors bubble only via caller).
- Anonymous tenant fallback exercised (no test sets no header or blank header for anonymous).

## BROKEN_IMPLEMENTATION
- None found; all claimed behaviors verified or warn only.

## DOC_CODE_MISMATCH
- engineering/03-execution-result.md Section 3 (Bounded Queue stats) mismatched with live deterministic observation (see 04-docs-vs-code).
- README "Known Limitations" omits distributed sync (covered in engineering/02-implementation-notes but not README).
- engineering/03-execution-result.md stale test count (missing 2 tests).

## RACE_CONDITION
- WARNING: BoundedQueue.Stop vs TrySubmit send-on-close race (see 02-code-audit Finding 6). Mitigated by typical defer Stop after workloads, but not proven safe under adversarial timing.

## UNHANDLED_ERROR
- None; all errors propagated or logged/discarded (job error discarded; explicit in code).

## MISSING_EDGE_CASE
- TokenBucket refillRate = 0 leads to infinite retryAfter (division by zero). No validation in constructor; not in tests.
- LeakyBucket leakRate = 0 -> never drains; bucket fills permanently; no test.
- BoundedQueue capacity = 0 or workers = 0 panics on send or zero worker goroutines (Test not present).
- Backoff attempt negative (uint conversion in math.Pow; attempt int negative leads to huge temp; no test).
- Backoff prevSleep negative in DecorrelatedJitter (float64 negative; leads to weird rangeMax; not tested).

## IMPLEMENTATION_OVERCLAIM
- None; wording in README "burst when water reaches capacity" is minor impreciseness (see 04-docs-vs-code).

## RESEARCH_MISMATCH
- Not auditing research per PIPELINE OVERRIDE.

## FAKE_DEMO
- Not found; demo output genuine, only recorded snapshot in engineering/03-execution-result.md outdated/wrong for section 3.

## FAKE_BENCHMARK
- No benchmarks present.

## UNVERIFIED_RESULT
- engineering/03-execution-result.md Section 3 result unverified (contradicted by live runs). Recorded numbers incorrect for given code/timing.