# Gap Analysis

Allowed gap types:
- MISSING_TEST
- BROKEN_IMPLEMENTATION
- DOC_CODE_MISMATCH
- RACE_CONDITION
- UNHANDLED_ERROR
- MISSING_EDGE_CASE
- IMPLEMENTATION_OVERCLAIM
- RESEARCH_MISMATCH
- FAKE_DEMO
- FAKE_BENCHMARK
- UNVERIFIED_RESULT

All gaps found during code & test audit (out of scope: research).

1. MISSING_TEST
   Location: `internal/circuit/circuit_test.go`
   Detail: No test for half-open failure → Open transition; no test for half-open success with SuccessThreshold != 2; no concurrent trial test in half-open; no zero-config test.

2. MISSING_TEST
   Location: `internal/idempotency/idempotency_test.go`
   Detail: Concurrent test makes no assertions about correctness (lost updates, stale reads); no test for check-then-set race outcome; no test for expired entry eviction on Get.

3. MISSING_TEST
   Location: `internal/retry/retry_test.go`
   Detail: No test verifying jitter range (base*2^(n-1) ≤ sleep ≤ max); no test for attempt <= 0 (Config validation); no test for zero-base/max-backoff leading to defaults.

4. MISSING_EDGE_CASE
   Location: `internal/deadline/deadline.go`
   Detail: Function does not enforce worker context-cancellation; worker may ignore ctx and leak (though select will return error). Not a bug but an edge-case where deadline claim assumes cooperation.

5. MISSING_EDGE_CASE
   Location: `internal/idempotency/idempotency.go`
   Detail: Expired entries not removed on Get (lazy deletion). Accumulate under high-cardinality keys until overwrite or restart.

6. MISSING_EDGE_CASE
   Location: `internal/idempotency/idempotency.go`
   Detail: Get-then-Set race (check-then-set) allows multiple Sets for same key; idempotency guarantee weakened under concurrent retries for same logical request.

7. DOC_CODE_MISMATCH
   Location: `labs/28-timeouts-and-deadlines/engineering/01-design.md:27`
   Detail: Backoff formula misprinted as `base * 2^attempt + jitter` (code uses `base * 2^(attempt-1)`). Corrected in engineering/02-implementation-notes.md. Low.

No BROKEN_IMPLEMENTATION, RACE_CONDITION (race detector passed), UNHANDLED_ERROR, IMPLEMENTATION_OVERCLAIM (claims bounded), FAKE_DEMO, FAKE_BENCHMARK, UNVERIFIED_RESULT found.

Severity mapping:
- MISSING_TEST: MEDIUM (reduces confidence; core behavior still unit-tested)
- MISSING_EDGE_CASE: MEDIUM (leak, weakened idempotency) to LOW (context cooperation)
- DOC_CODE_MISMATCH: LOW