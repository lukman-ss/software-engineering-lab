# Docs vs Code Audit

## Item Comparison

| Documented Item | Claimed Location / Behavior | Actual Implementation in Code | Discrepancy Found |
| :--- | :--- | :--- | :--- |
| `internal/outbox/db.go` | In-memory transactional DB simulating BeginTx, Commit, and Rollback across orders and outbox | Exact match (`DB`, `Tx`, `BeginTx`, `Commit`, `Rollback`) | None |
| `internal/outbox/broker.go` | Thread-safe mock message broker simulating failures and event reception | Exact match (`MockBroker` with mutex, `SetFailNext`, `Publish`, `GetPublished`) | None |
| `internal/outbox/service.go` | Business logic comparing naive dual-write vs atomic outbox writes | Exact match (`CreateOrderWithOutbox` vs `CreateOrderDualWriteNaive`) | None |
| `internal/outbox/relay.go` | Asynchronous polling worker querying pending outbox records and dispatching them to the broker | Exact match (`Relay` polling loop, `PollAndDispatch`) | None |
| `internal/outbox/consumer.go` | Subscriber enforcing idempotency through event ID tracking | Exact match (`Consumer.Handle` checking `processedIDs`) | None |
| Test commands | `go test ./...` and `go test -race ./...` | Both run cleanly and pass | None |
| Demo command | `go run ./cmd/demo` | Runs cleanly, outputs 3 scenarios matching README claims | None |

## Finding Summary

No documentation-to-code mismatch identified.
README accurately describes file structure, design concepts, and verification commands.
