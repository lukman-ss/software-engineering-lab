# Docs vs Code

## Comparison Matrix

| Artifact | Claim | Actual | Status |
|---|---|---|---|
| README.md architecture | `internal/outbox/db.go` transactional DB simulating BeginTx/Commit/Rollback across orders + outbox | in-memory Tx with staged maps, Commit applies atomically | PASS |
| README.md architecture | `internal/outbox/broker.go` thread-safe mock broker with publish failures + reception | MockBroker with SetFailNext fault injection | PASS |
| README.md architecture | `internal/outbox/service.go` naive dual-write vs atomic outbox | CreateOrderWithOutbox + CreateOrderDualWriteNaive | PASS |
| README.md architecture | `internal/outbox/relay.go` async polling worker dispatching to broker | Relay with ticker goroutine PollAndDispatch | PASS |
| README.md architecture | `internal/outbox/consumer.go` idempotency via event ID tracking | Consumer with processedIDs map | PASS |
| README.md commands | `go test ./...` | PASS (6 tests, 3 packages) | PASS |
| README.md commands | `go test -race ./...` | PASS (no races) | PASS |
| README.md commands | `go run ./cmd/demo` | PASS (exact output matches engineering/03-execution-result.md) | PASS |
| engineering/01-design.md | Atomic insert order + outbox in single transaction | PASS |
| engineering/01-design.md | Rollback leaves neither persisted | PASS |
| engineering/01-design.md | Relay dispatches and marks PROCESSED | PASS |
| engineering/01-design.md | Consumer tracks event IDs for idempotency | PASS |
| engineering/01-design.md | Outbox cleanup job purges processed records | PASS |
| engineering/01-design.md | SQLite DB implementation | FAIL - design references SQLite (`modernc.org/sqlite`, `outbox_events` table, `FOR UPDATE`, CDC/Debezium) but actual implementation uses in-memory Go map-based Tx with no SQL engine | MISMATCH |
| engineering/01-design.md | Relay polls with FOR UPDATE lock | FAIL - no row-level locking; uses RLock on map | MISMATCH |
| engineering/01-design.md | Unique UUID event identifiers | FAIL - IDs are `evt-<orderID>` deterministic strings, not UUIDs | MISMATCH |
| research-audit (not audited) | - | Per pipeline override, research not audited in this stage | N/A |

## Mismatches Found:

1. DOC_CODE_MISMATCH: Design doc (`engineering/01-design.md`) describes SQLite/database implementation but code uses in-memory map-based simulation. README accurately describes actual implementation.
2. DOC_CODE_MISMATCH: Design doc references `FOR UPDATE` row locking; code uses `sync.RWMutex` on map.
3. DOC_CODE_MISMATCH: Design doc claims UUID event identifiers; code uses deterministic `evt-<orderID>` IDs.

README vs implementation: PASS (README accurately describes in-memory implementation).

Demo output vs engineering/03-execution-result.md: PASS (identical output).

## Conclusion

README matches code. Engineering design notes over-claim SQLite features that were not implemented. Implementation notes accurately document what was built.