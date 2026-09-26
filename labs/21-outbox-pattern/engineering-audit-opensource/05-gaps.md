# Engineering Gap Analysis

All findings mapped to allowed gap types. Severity assessed per rubric:
- LOW: minor issue.
- MEDIUM: incomplete coverage or docs mismatch.
- HIGH: core behavior unproven or incorrect.
- CRITICAL: fabricated result or fundamentally invalid implementation.

## Gap 1: Missing Cleanup Worker

Type: DOC_CODE_MISMATCH
Location: engineering/01-design.md (Expected Behavior #5, Success Criteria #5, Component "OutboxRelay ... with retry / cleanup support")
Observed: No cleanup/purge/delete API exists in db.go; relay only marks PENDING→PROCESSED; processed records accumulate indefinitely.
Severity: MEDIUM
Notes: Auxiliary to the core transactional guarantee but an explicit, unmet success criterion. README does not over-claim cleanup.

## Gap 2: Database Engine Mismatch

Type: DOC_CODE_MISMATCH
Location: engineering/01-design.md (Components: "SQLite Database"; Architecture diagram); engineering/02-implementation-notes.md (Implementation Decisions list modernc.org/sqlite or stdlib database/sql)
Observed: DB is sync.RWMutex-guarded maps; no SQLite/database/sql dependency. go.mod confirms no such library. README correctly describes an in-memory DB.
Severity: MEDIUM
Notes: Transactional semantics preserved via map-based Tx; simplification documented in implementation notes but contradicts design doc's stated engine. Not a correctness bug.

## Gap 3: Relay Stop Not Idempotent

Type: MISSING_EDGE_CASE
Location: internal/outbox/relay.go:39-41 (Relay.Stop)
Observed: close(r.stopChan) with no guard; calling Stop() twice panics "close of closed channel".
Severity: LOW
Notes: Not exercised by current usage (single defer Stop/test/demo). Fragile for reuse but not a correctness issue in the scoped deployment.

## Gap 4: Weak Concurrent Writes Test

Type: MISSING_TEST
Location: tests/outbox_test.go:TestTransactionalOutbox_ConcurrentWrites
Observed: 100 CreateOrderWithOutbox calls use identical key "o-concurrent-1"; test ends with zero assertions on persisted orders, outbox state, broker publishes, or consumer state.
Severity: MEDIUM
Notes: Exists for race detector (go test -race passes) but lacks correctness verification. Does not test concurrent distinct-key commits or validate throughput.

## Gap 5: Missing Broker Failure Recovery Test

Type: MISSING_TEST
Location: N/A (omitted test)
Observed: No test where the relay encounters a broker Publish failure (e.g., via simulated network error) and later recovers (e.g., after failure clears) to verify the message eventually publishes and is marked processed.
Severity: MEDIUM
Notes: Implicit re-poll retry provides at-least-once, but the claimed "Broker failure retry mechanism" (engineering/02-implementation-notes.md Test Strategy) lacks dedicated validation.

## Gap 6: Missing MarkOutboxProcessed Failure Path Test

Type: MISSING_TEST
Location: N/A (omitted test)
Observed: No test where broker.Publish succeeds but db.MarkOutboxProcessed fails (simulated error) and verify the message remains PENDING for redelivery on the next poll tick.
Severity: MEDIUM
Notes: Tests the at-least-once semantics when marking fails but publish succeeds (message stays PENDING → redelivered → duplicate handled by idempotent consumer).

## Gap 7: Missing Service-Layer Rollback Test

Type: MISSING_TEST
Location: N/A (omitted test)
Observed: No test where service.CreateOrderWithOutbox returns an error (e.g., JSON marshal failure) and verifies that neither the order nor outbox event is persisted (rollback behavior at the service layer).
Severity: MEDIUM
Notes: Tests rollback via direct Tx exist (TestTransactionalOutbox_Rollback), but service-layer error paths are not explicitly validated for rollback.

## Gap 8: Missing Exactly-Once Semantics Test

Type: MISSING_TEST
Location: N/A (omitted test)
Observed: No test verifying that N successful CreateOrderWithOutbox calls (with distinct keys) result in exactly N broker messages published and exactly N consumer receptions (e.g., consumer.GetReceivedCount() == N).
Severity: MEDIUM
Notes: Concurrent test uses identical keys, so cannot verify this property. Ideal-condition correctness is untested at scale.

## Gap 9: Design Overclaim on Relay Locking

Type: DOC_CODE_MISMATCH
Location: engineering/01-design.md (Architecture diagram line 39 showing "Poll (FOR UPDATE / Lock)")
Observed: Code uses GetPendingOutbox (snapshot read) and per-message Publish/MarkOutboxProcessed; no row-level locking or claiming step. Only a single relay instance is used; README does not promise multi-instance relay.
Severity: LOW
Notes: Out of scope for the lab's single-instance deployment; idempotency absorbs duplicates if they occurred.

## Summary of Gap Types
- DOC_CODE_MISMATCH: 4 instances (Gaps 1, 2, 3, 9)
- MISSING_TEST: 5 instances (Gaps 4, 5, 6, 7, 8)
- MISSING_EDGE_CASE: 1 instance (Gap 8? Wait no, Gap 3 is MISSING_EDGE_CASE)
  - Actually: Gap 3 is MISSING_EDGE_CASE
- BROKEN_IMPLEMENTATION: 0
- RACE_CONDITION: 0 (go test -race clean)
- UNHANDLED_ERROR: 0
- MISSING_EDGE_CASE: 1 (Gap 3)
- IMPLEMENTATION_OVERCLAIM: 0
- RESEARCH_MISMATCH: 0 (override: not auditing research)
- FAKE_DEMO: 0 (demo output verified real)
- FAKE_BENCHMARK: 0 (no benchmarks)
- UNVERIFIED_RESULT: 0 (03-execution-result.md matches observed)