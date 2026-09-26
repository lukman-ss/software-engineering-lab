# Implementation Notes

## Files Added
- `go.mod`: Module definition for `github.com/software-engineering-lab/labs/21-outbox-pattern`.
- `internal/outbox/model.go`: Domain models (`Order`, `OutboxMessage`, statuses).
- `internal/outbox/db.go`: Thread-safe transactional mock database with staged mutations on `Tx`.
- `internal/outbox/broker.go`: Mock message broker with configurable failure simulation.
- `internal/outbox/service.go`: Order creation service implementing both naive dual-write and transactional outbox approaches.
- `internal/outbox/relay.go`: Polling message relay worker querying pending outbox records and updating status upon publish.
- `internal/outbox/consumer.go`: Idempotent consumer tracking processed message IDs.
- `cmd/demo/main.go`: End-to-end runnable demo showcasing the dual-write problem, outbox resolution, and duplicate suppression.
- `tests/outbox_test.go`: Automated tests for happy path, rollback safety, duplicate delivery deduplication, dual-write failure, and concurrent execution.
- `engineering/01-design.md`: Engineering design specification.
- `engineering/02-implementation-notes.md`: Implementation choices, trade-offs, and limitations.
- `engineering/03-execution-result.md`: Recorded execution outputs.
- `README.md`: Lab overview and execution guide.

## Core Design Decisions
- Atomicity: Used an in-memory transactional wrapper where staged mutations to both `orders` and `outbox` commit together or discard on rollback.
- Asynchronous Polling Relay: Employed a background ticker that fetches pending messages, publishes them to the broker, and updates outbox status to `PROCESSED`.
- Idempotent Consumption: Maintained an in-memory set of processed event IDs in the consumer to filter duplicates.

## Implementation-Specific Choices
- Implementation Decision: In-memory simulation rather than external SQL / CGO-based SQLite driver to keep dependencies minimal (pure Go standard library).
- Implementation Decision: Polling publisher pattern chosen instead of Transaction Log Tailing (CDC), matching the primary simplest relay option identified in the research.

## Known Limitations
- In-memory persistence does not survive process restarts.
- Polling frequency is fixed and does not implement exponential backoff on broker failures.
- No outbox cleanup/retention worker implemented for old `PROCESSED` events.

## Trade-offs
- Simplicity vs Production Scale: In-memory maps provide clear concurrency-safe semantics for proving the pattern without introducing external operational complexity.
- Polling Overhead: Polling introduces slight latency and database query overhead compared to log-tailing CDC solutions (e.g. Debezium), but is simpler to understand and test.

## What Is Demonstrated
- The dual-write flaw when broker write fails after database commit.
- Atomic commit of business entities and outbox events in a single transaction.
- Discarding outbox events on transaction rollback.
- Message relay polling and successfully publishing outbox events.
- Consumer deduplication ensuring at-least-once delivery does not cause duplicate processing.
- Thread-safe concurrent execution under the Go race detector.

## What Is Not Demonstrated
- Change Data Capture (CDC) via database transaction logs.
- Distributed broker partitions, consumer groups, or dead-letter queues.
- Outbox table compaction and long-term purging strategies.
