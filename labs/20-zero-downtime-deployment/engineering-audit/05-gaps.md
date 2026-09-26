# Gap Analysis

Target Lab: labs/20-zero-downtime-deployment

---

## GAP-01

Type: MISSING_TEST
Location: internal/db/db.go — GetUser with ErrNotFound
Description: No test verifies that requesting a non-existent ID returns ErrNotFound.
Severity: LOW
Blocking: No — error path is simple and unambiguous from code inspection.

---

## GAP-02

Type: MISSING_TEST
Location: internal/db/db.go — single-name legacy user
Description: No test for InsertLegacy with a name containing no space (e.g., "Madonna"). GetUser would return FirstName="Madonna", LastName="" — correct behavior, but untested.
Severity: LOW
Blocking: No.

---

## GAP-03

Type: MISSING_TEST
Location: internal/db/db.go — SaveExpand with empty fields
Description: No test for SaveExpand("id", "", "Smith") or SaveExpand("id", "Jane", ""). TrimSpace logic handles these but is untested.
Severity: LOW
Blocking: No.

---

## GAP-04

Type: MISSING_TEST
Location: internal/worker/worker.go — Enqueue after Stop
Description: Enqueue after Stop panics (send on closed channel). No guard exists, no test covers this. In the lab's controlled usage pattern this never occurs.
Severity: MEDIUM
Blocking: No — lab usage is controlled; caller responsibility is documented implicitly by the pattern.

---

## GAP-05

Type: MISSING_TEST
Location: internal/worker/worker.go — concurrency > 1
Description: All worker tests use concurrency=1. Multi-worker concurrent job completion order and shared completed slice access under concurrent writes are untested at concurrency > 1. Race detector passes for the single-worker case.
Severity: LOW
Blocking: No — completedMu correctly protects the shared slice; race detector passes.

---

## GAP-06

Type: MISSING_TEST
Location: internal/server/server.go — /work with invalid duration query parameter
Description: No test verifies that an invalid `d` parameter falls back to 50ms default. The code path exists and is correct, but is untested.
Severity: LOW
Blocking: No.

---

## GAP-07

Type: DOC_CODE_MISMATCH
Location: engineering/01-design.md line 47
Description: Design describes Shutdown signature as `Shutdown(ctx, preStopDelay)`. Implementation uses `Shutdown(ctx context.Context)` with preStopDelay as a constructor argument. Minor documentation drift from design-to-implementation iteration.
Severity: LOW
Blocking: No — functional behavior matches intent.

---

## GAP-08

Type: MISSING_TEST
Location: internal/server/server.go — /healthz/ready state transition READY→UNREADY
Description: Tests cover NOT READY → READY via SetReady(true). The reverse transition (marking unready again, as happens during shutdown) is exercised indirectly inside Shutdown() but not as an explicit probe test.
Severity: LOW
Blocking: No — Shutdown test implicitly validates this by testing a server that starts ready and then shuts down.

---

## No CRITICAL or HIGH gaps found.

All primary claimed behaviors are implemented, proven by tests, and confirmed by execution.
