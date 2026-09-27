# Docs vs Code Audit

Target Lab: labs/25-rate-limiting-and-backpressure
Date: 2026-09-27
Scope: README.md, engineering/01-design.md, engineering/02-implementation-notes.md, engineering/03-execution-result.md vs implementation/tests/demo.

## Match

- README Feature list matches code: TokenBucket, LeakyBucket, Registry, BoundedQueue, jitter strategies, 429 middleware. PASS.
- README structure tree matches repo layout. PASS.
- Design claims (burst cap, leak rate, fast shed, Full Jitter formula, 429+Retry-After) implemented. PASS.
- README "no third-party deps" true (stdlib only). PASS.
- README command snippets correct and verified live. PASS.

## Mismatch 1 (DOC_CODE_MISMATCH)

Location: engineering/03-execution-result.md Section 3 vs cmd/demo + queue.go
Claimed Result: Jobs 4,5,6 rejected; Stats Accepted=3 Rejected=3 Processed=1.
Observed Run: Job 4 rejected; Job 5 accepted; Job 6 rejected; Stats Accepted=4 Rejected=2 Processed=1.
Assessment: WARNING
Severity: MEDIUM
Notes: Documented stats contradict reproducible execution. Root cause: 50ms simulated job vs microsecond-fast submit loop races; worker may or may not drain a slot before later submits. Output is run-dependent. Recorded snapshot not the deterministic "truth." README claims "CLI demo runs showing..."; claim valid but specific recorded numbers are misleading.

## Mismatch 2 (DOC_CODE_MISMATCH)

Location: engineering/03-execution-result.md test list vs actual
Claimed Result lists 9 tests: TestBoundedQueue_*, TestRateLimitMiddleware_RFC6585, TestTokenBucket_BurstAndRefill, TestLeakyBucket_LeakRate, TestRegistry_TenantIsolation, TestTokenBucket_ConcurrencyRace, TestComputeBackoff_Bounds, TestDecorrelatedJitter_Bounds.
Observed: 11 tests present — additionally TestTokenBucket_RetryAfterSeconds, TestBoundedQueue_SubmitAfterStop.
Assessment: PASS (superset; doc understated, harmless) but stale.
Severity: LOW

## Mismatch 3 (IMPLEMENTATION_OVERCLAIM)

Location: README "LeakyBucket: Enforces constant drain rate R, rejecting bursts when water reaches capacity."
Claimed Behavior: Bursts rejected at capacity.
Observed: LeakyBucket.Allow admits only +1 unit per call; a burst of N calls fills N, capped at capacity. A single call with capacity headroom is always admitted; "burst" is a sequence of immediate calls, not a single n>1 call. Semantics correct but wording implies per-call burst-size rejection.
Assessment: WARNING
Severity: LOW
Notes: Minor imprecision, no behavioral gap.

## Match (verified live)

- demo sections 1,2,4 behaviors align with recorded description (429, jitter ranges). Section 3 numbers only diverge. PASS for behavior.

## Doc Gaps (MISSING_TEST/MISSING_EDGE_CASE)

- Zero leakRate/refillRate not validated anywhere; math division hazards (bucket.go RetryAfterSeconds) not guarded nor claimed. README/design silent.
- Distributed/registry TTL not in README "Known Limitations" scope? Implementation Notes lists distributed sync limitation; README omits. Minor.
- No mention of Stop-vs-Submit race condition in design notes.
