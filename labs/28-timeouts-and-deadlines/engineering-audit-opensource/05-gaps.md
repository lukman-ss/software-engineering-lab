# Gap Analysis — labs/28-timeouts-and-deadlines

## Gap 1

Type: UNHANDLED_ERROR
Location: internal/deadline/deadline.go:13-27
Severity: MEDIUM
Description: On timeout path the worker goroutine result is discarded; fn error/panic never surfaces. Background worker leak until fn returns; no timeout on worker itself.
Status: accepted limitation; document only.

## Gap 2

Type: MISSING_TEST
Location: internal/deadline (no test)
Severity: LOW
Description: No test for select-race branch (fn completes exactly as ctx expires) or for fn ignoring ctx (leak path). No goroutine-leak assertion.
Status: open.

## Gap 3

Type: MISSING_TEST
Location: internal/circuit (no test)
Severity: LOW
Description: No concurrent HALF_OPEN probe test; unlimited concurrent probes allowed by design. No test that single success below SuccessThreshold stays HALF_OPEN.
Status: open.

## Gap 4

Type: MISSING_TEST
Location: internal/retry (no test)
Severity: LOW
Description: No statistical test proving jitter decorrelates retries (e.g. variance > 0 over N samples); bounds-only check cannot prove anti-thundering-herd claim. No ctx-cancel-during-backoff-after-success case.
Status: open.

## Gap 5

Type: MISSING_EDGE_CASE
Location: internal/idempotency/idempotency.go:29-42
Severity: LOW
Description: Expired-but-never-read keys never evicted (no sweeper); unbounded map growth for write-once keys. Set overwrites + refreshes TTL without documenting refresh semantics.
Status: accepted for lab scope; implementation notes already flag non-persistence.

## Gap 6

Type: UNVERIFIED_RESULT
Location: engineering/01-design.md Success Criteria 1 ("without leakage")
Severity: MEDIUM
Description: "Without leakage" claimed in design; test suite has no leak check (e.g. goleak or goroutine count). Race detector passes, leak-freedom unproven for deadline path (see Gap 1).
Status: open; downgrade claim to "returns promptly; worker bounded by fn's ctx cooperation" or add leak test.

No BROKEN_IMPLEMENTATION. No RACE_CONDITION. No FAKE_DEMO. No FAKE_BENCHMARK. No IMPLEMENTATION_OVERCLAIM beyond Gap 6 wording. No DOC_CODE_MISMATCH.
