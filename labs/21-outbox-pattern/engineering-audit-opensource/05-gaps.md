# Gap Analysis

## Identified Gaps

| # | Gap Type | Location | Description |
|---|---|---|---|
| 1 | WARNING | internal/outbox/relay.go:50-54 | Relay increments dispatched++ even when MarkOutboxProcessed fails, breaking at-least-once guarantee consistency. |
| 2 | MISSING_TEST | tests/outbox_test.go | No test for relay retry after broker transient failure (broker recovers, message re-dispatched, consumer dedups). |
| 3 | MISSING_TEST | tests/outbox_test.go | No test for MarkOutboxProcessed on non-existent ID (error path). |
| 4 | MISSING_EDGE_CASE | tests/outbox_test.go:141-170 | TestTransactionalOutbox_ConcurrentWrites has no behavioral assertions; writes same orderID so only 1 order persists. |
| 5 | DOC_CODE_MISMATCH | engineering/01-design.md | Design doc references SQLite, `FOR UPDATE`, UUID IDs; actual impl is in-memory map simulation without these features. |
| 6 | RACE_CONDITION | internal/outbox/relay.go:49-54 | Re-poll after failed publish+mark leaves message PENDING. If broker was actually up but MarkOutboxProcessed failed (race on another concurrent MarkOutboxProcessed), message could be re-published causing duplicate processing despite consumer idempotency (acceptable but unhandled). |

## Severity Assessment

All gaps are LOW-to-MEDIUM. No CRITICAL or fabricated results. Implementation is honest simulation; discrepancies exist only in design-doc over-claiming (not deception in actual code/README).