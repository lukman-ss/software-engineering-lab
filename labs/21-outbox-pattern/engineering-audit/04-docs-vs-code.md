# Docs vs Code Audit

## Comparison Matrix

| Item | Documented Claim | Implemented Reality | Result |
| :--- | :--- | :--- | :--- |
| **Architecture** | In-memory transactional DB (`BeginTx`, `Commit`, `Rollback`), mock broker, service, relay, consumer. | Exactly matches `internal/outbox/` files. | PASS |
| **Dual-Write Vulnerability** | Direct write fails on broker error, leaving DB committed. | Implemented in `service.CreateOrderDualWriteNaive` and proven in `TestDualWriteProblem_Failure` & Scenario 1 of demo. | PASS |
| **Outbox Atomic Commit** | Entity and outbox event written in one transaction. | Implemented in `service.CreateOrderWithOutbox` via single `Tx`. | PASS |
| **Relay Polling** | Background polling worker queries pending messages and pushes to broker. | Implemented in `relay.Relay` with ticker-based loop. | PASS |
| **Consumer Idempotency** | Event ID deduplication cache drops duplicate messages. | Implemented in `consumer.Consumer` with `processedIDs` map. | PASS |
| **Demo Output** | README references running `go run ./cmd/demo`. Execution result file records scenario output. | Actual `go run ./cmd/demo` output matches `engineering/03-execution-result.md` verbatim. | PASS |

## Discrepancy Checks
- `DOC_CODE_MISMATCH`: None detected. README and engineering notes match code structure and behavior.
- `TEST_CLAIM_MISMATCH`: None detected. All tests directly assert the documented guarantees.
- `RESEARCH_IMPLEMENTATION_MISMATCH`: None detected. Implements Polling Publisher and Idempotent Consumer patterns matching research recommendations.
