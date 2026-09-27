# Content vs Implementation Audit

Target Lab: `labs/21-outbox-pattern`

## 1. Domain & Table Schema
- Master draft distinguishes between full production/Debezium schema (`id`, `aggregateid`, `aggregatetype`, `type`, `payload`) and lab implementation (`ID`, `EventType`, `Payload`, `Status`, `CreatedAt`).
- Verified against `internal/outbox/model.go`. Matches.

## 2. Order Service
- `CreateOrderWithOutbox` (`internal/outbox/service.go:17-53`): Atomically persists Order and Outbox record in single Tx. Rollback triggers on marshal error or save error.
- `CreateOrderDualWriteNaive` (`internal/outbox/service.go:55-90`): Commits DB Tx first, then publishes to broker.
- Snippet and master draft accurately quote and reference line ranges.

## 3. Database & Staging Transaction
- `internal/outbox/db.go:84-140` implements `Tx` with `stagedOrders` and `stagedOutbox`.
- `Commit` flushes staged maps into live DB maps under lock.
- `Rollback` marks Tx closed and discards staged maps.
- Verified in `content/03-code-snippets.md` and `content/02-master-draft.md`.

## 4. Polling Relay
- `internal/outbox/relay.go:50-67` implements `PollAndDispatch()`.
- Fetches pending, publishes, marks processed on success; retains pending on failure for retry.
- Matches `content/02-master-draft.md` and `content/04-diagrams.md`.

## 5. Idempotent Consumer
- `internal/outbox/consumer.go:19-31` tracks `processedIDs map[string]bool`.
- First delivery returns `true`, duplicate returns `false`.
- Matches master draft and test assertions.

## 6. Test Suite Alignment
- Documented test count: 8 tests.
- Verified test suite: `TestTransactionalOutbox_HappyPath`, `TestTransactionalOutbox_Rollback`, `TestTransactionalOutbox_Idempotent_DuplicateDelivery`, `TestDualWriteProblem_Failure`, `TestTransactionalOutbox_ConcurrentWrites`, `TestTransactionalOutbox_PurgeProcessed`, `TestTransactionalOutbox_RelayRetryAfterBrokerFailure`, `TestTransactionalOutbox_ConcurrentConsumers`.
- Exactly 8 tests; all pass under race detector.
