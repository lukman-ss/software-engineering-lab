# Gap Analysis

Target Lab: labs/21-outbox-pattern

---

## Gap 1

Type: MISSING_TEST
Severity: MEDIUM
Location: tests/outbox_test.go:TestTransactionalOutbox_ConcurrentWrites
Description: Concurrent test uses same ID for all 100 writers; no outcome assertions.
Impact: Test proves absence of data races (−race passes) but not concurrent write correctness. Last-writer-wins overwrite hides real concurrent behavior. Unique IDs per worker would prove N distinct outbox messages eventually published.
Action: Use fmt.Sprintf with worker + j to generate unique order IDs and assert final published count + PROCESSED status for all.

## Gap 2

Type: MISSING_TEST
Severity: MEDIUM
Location: tests/outbox_test.go (missing)
Description: No concurrent consumer test under race detector. Consumer.Handle / IsProcessed / GetReceivedCount share mutex-protected state but are never exercised concurrently.
Impact: Consumer race safety unverified by test (though lock discipline looks sound and global -race run is clean for current paths).
Action: Add test with concurrent Handle() calls of same + distinct IDs under -race.

## Gap 3

Type: MISSING_TEST
Severity: LOW
Location: internal/outbox/db.go:MarkOutboxProcessed (not-found error path)
Description: Error return when message ID not found is untested.
Impact: Minor — relay's error branch (log + no increment) unverified.
Action: Add negative test calling MarkOutboxProcessed with unknown ID, expect error.

## Gap 4

Type: MISSING_TEST
Severity: LOW
Location: internal/outbox/db.go:SaveOrder/SaveOutbox/Commit after close
Description: ErrTxClosed paths untested (double Commit, Save after Commit, Save after Rollback).
Impact: Minor — Tx lifecycle edge cases unverified.
Action: Add test for double-commit / post-close save errors.

## Gap 5

Type: MISSING_EDGE_CASE
Severity: MEDIUM
Location: internal/outbox/model.go (OutboxMessage.ID)
Description: No UUID enforcement; deterministic `evt-<orderID>` IDs risk collision in concurrent/production use. Design doc claims "unique UUID event identifiers" but implementation uses caller-provided string IDs.
Impact: In concurrent test all IDs collide (Finding 9). In production this would cause lost events via map overwrite. Acceptable for demo scope.
Action: Document ID uniqueness requirement; optionally generate UUIDs in service (stdlib, no new dep) — out of audit-fixing scope.

## Gap 6

Type: DOC_CODE_MISMATCH
Severity: MEDIUM
Location: engineering/01-design.md (SQLite schema, relay retry/cleanup, consumer_log)
Description: Design doc describes SQLite DB, outbox_events/consumer_log tables, relay with retry/cleanup support. Implementation is in-memory maps, no automated retry policy/backoff/DLQ, no cleanup worker, no consumer_log table.
Impact: Design overstates implementation. Implementation-notes + master-draft correctly scope to in-memory polling demo. A reader following design doc alone would expect more.
Action: Align design doc with as-built implementation (in-memory simulation, polling only, manual purge) or implement missing pieces.

## Gap 7

Type: IMPLEMENTATION_OVERCLAIM
Severity: MEDIUM
Location: engineering/01-design.md Success Criteria #5 ("Cleanup worker successfully purges processed events")
Description: No cleanup worker exists. PurgeProcessedOutbox() is manual, untested in relay/demo flow.
Impact: Success criterion unproven by code.
Action: Either implement periodic purge in relay/demo or reword criterion to "manual purge function verified by unit test".

## Gap 8

Type: UNVERIFIED_RESULT
Severity: LOW
Location: engineering/03-execution-result.md
Description: Recorded results claim PASS for build/tests/race/demo.
Impact: Verification performed in this audit — all four commands re-run and PASS. Recorded results confirmed real.
Action: None. Results reproduced 2026-09-27 (go 1.26.7): build exit 0, 7/7 tests PASS, -race PASS, demo output matches recorded output verbatim.

## Gap 9

Type: UNHANDLED_ERROR
Severity: LOW
Location: internal/outbox/relay.go:PollAndDispatch
Description: MarkOutboxProcessed failure is logged only, message stays PENDING → will be re-published next cycle (duplicate). Publish uses fire-and-forget with log on failure.
Impact: By design (at-least-once + idempotent consumer covers it). No DLQ for poison messages — a message that always fails to publish is retried forever each pollInterval.
Action: Document no-backoff/no-DLQ as known limitation (already done in master-draft). Out of scope for demo.

## Gap 10

Type: MISSING_TEST
Severity: LOW
Location: tests/outbox_test.go (missing)
Description: No test for concurrent relay-poll + purge interaction (MarkOutboxProcessed racing PurgeProcessedOutbox delete).
Impact: Purge deletes PROCESSED rows; relay marks PENDING→PROCESSED. A purge between publish and mark could cause MarkOutboxProcessed "not found" → logged, message counted as undispatched but already published. Edge case unverified.
Action: Add test: publish → purge → mark, assert error handled gracefully. Low priority for demo scope.

---

## Severity Summary

CRITICAL: 0
HIGH: 0
MEDIUM: 4 (Gap 1, 2, 5, 6, 7 — weak concurrency test; consumer concurrency untested; ID uniqueness; design-doc overclaim incl. cleanup criterion)
LOW: 5 (Gap 3, 4, 8-resolved, 9, 10)

No HIGH/CRITICAL gaps. Blocking: none. All gaps are non-blocking warnings consistent with an in-memory demo lab.