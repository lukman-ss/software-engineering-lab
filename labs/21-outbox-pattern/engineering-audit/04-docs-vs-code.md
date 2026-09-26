# Docs vs Code Audit

Target Lab: `labs/21-outbox-pattern`

## Documents Reviewed
1. `README.md`
2. `engineering/01-design.md`
3. `engineering/02-implementation-notes.md`
4. `engineering/03-execution-result.md`

## Comparison Matrix

| Component / Claim | README & Design Claim | Actual Implementation | Status |
|---|---|---|---|
| Architecture & File Structure | `db.go`, `broker.go`, `service.go`, `relay.go`, `consumer.go` | Exactly matches internal package structure | MATCH |
| Database Engine | `01-design.md` mentions SQLite DB or mock engine | Implemented as zero-dependency in-memory transactional DB (`internal/outbox/db.go`), documented in `02-implementation-notes.md` | MATCH |
| Dual-Write Failure Demonstration | Demonstrates state inconsistency on broker failure | Implemented in `service.CreateOrderDualWriteNaive` and demo Scenario 1 | MATCH |
| Atomic Outbox Persistence | Persists order and outbox record atomically in single transaction | Implemented in `service.CreateOrderWithOutbox` and `Tx.Commit` | MATCH |
| Polling Relay | Background polling worker with ticker interval | Implemented in `relay.Start` / `PollAndDispatch` | MATCH |
| Consumer Deduplication | Consumer tracking processed event IDs | Implemented in `consumer.Handle` with internal map | MATCH |
| Outbox Purge | Purges processed outbox records | Implemented in `db.PurgeProcessedOutbox` | MATCH |
| Demo Execution | `go run ./cmd/demo` executes scenarios 1, 2, and 3 | Verified real terminal execution matches `03-execution-result.md` | MATCH |

## Identified Discrepancies
None. Design notes accurately capture the in-memory transactional model choice, limitations (no CDC / log tailing), and trade-offs.
