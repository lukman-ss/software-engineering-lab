# Content Accuracy Audit

## Scope
`content/01-content-brief.md` through `content/07-revision-record.md` vs `research/05-report.md`, `research-audit/07-verdict.md`, `engineering-audit/06-verdict.md`, and source under `internal/outbox/`, `tests/outbox_test.go`, `cmd/demo/main.go`.

## Verified Accurate
- Dual-write problem definition (Tx only controls DB, cannot ROLLBACK broker) — matches research Finding 1, engineering `CreateOrderDualWriteNaive`.
- Atomic outbox (order + outbox in single Tx, commit/rollback together) — matches `service.go:17-53`, `db.go:112-139`, `TestTransactionalOutbox_HappyPath`, `TestTransactionalOutbox_Rollback`.
- Polling Publisher relay (fetch PENDING, publish, mark PROCESSED, retry on fail) — matches `relay.go:50-67`, `TestTransactionalOutbox_RelayRetryAfterBrokerFailure`.
- At-least-once + idempotent consumer via `processedIDs` map — matches `consumer.go:19-31`, research Finding 4, both idempotency tests.
- Concurrent safety 10x10=100, 0 pending, race clean — matches `TestTransactionalOutbox_ConcurrentWrites` + engineering audit race PASS.
- PurgeProcessedOutbox deletes PROCESSED, retains PENDING — matches `db.go:71-82`, `TestTransactionalOutbox_PurgeProcessed`.
- Demo 3 scenarios (dual-write fail, outbox happy, duplicate) — matches `cmd/demo/main.go:10-63`.
- Code snippet line ranges in `03-code-snippets.md` verified against current source — all correct within 1 line.
- SLA thresholds (2s normal / 47m abnormal) correctly marked illustrative, not universal — consistent with research Finding 6 LIMITATIONS.

## Inaccuracies / Drifts
1. **Table schema vs implementation gap** — `02-master-draft.md:22-27` and `01-content-brief.md:20` list columns `aggregatetype`, `aggregateid` as table design. `internal/outbox/model.go:26-32` `OutboxMessage` has only `ID`, `EventType`, `Payload`, `Status`, `CreatedAt` — no aggregate fields. Research Finding 5 describes those columns for production/Debezium, lab intentionally simplified. Master draft presents them without explicitly noting they are absent from lab struct — reader may expect fields that do not compile.
2. **Revision record stale** — `07-revision-record.md:79` claims "All 6 tests pass" and throughout tables counts 5→6. Actual suite is 8 test functions; engineering-audit/06-verdict.md:9 and content brief/master-draft/source-map correctly state 8. Revision record also references line numbers (326, 558, 658 etc.) exceeding current `02-master-draft.md` length (227 lines) — indicates revision record not regenerated after content resize.
3. **Source map checklist unchecked** — `06-source-map.md:87-91` shows unchecked boxes against `engineering/01-design.md` criteria despite engineering audit APPROVED. Not inaccurate but incomplete closure.
