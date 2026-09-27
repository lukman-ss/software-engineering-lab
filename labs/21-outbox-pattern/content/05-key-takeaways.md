# Key Takeaways

1. Dual-write (DB commit + broker publish in separate steps) can leave the system inconsistent when the broker fails after a successful DB commit — the lab proves this with `TestDualWriteProblem_Failure` and demo Scenario 1.

2. The Transactional Outbox pattern moves the message write into the same database transaction as the business entity, so either both persist or neither does — verified by `CreateOrderWithOutbox` and `TestTransactionalOutbox_HappyPath`.

3. A rollback discards both the business entity and the outbox record at the staged level — nothing reaches the database, and no event is ever dispatched to the broker. Verified by `TestTransactionalOutbox_Rollback`.

4. The message relay implements the Polling Publisher strategy: it polls `PENDING` records, publishes to the broker, and only marks a record as `PROCESSED` after a successful publish. Failures leave the record retryable. Verified by `TestTransactionalOutbox_RelayRetryAfterBrokerFailure`.

5. The pattern provides at-least-once delivery, not exactly-once. A relay crash after publish but before marking processed causes a duplicate. The consumer must therefore be idempotent. Verified by `TestTransactionalOutbox_Idempotent_DuplicateDelivery`.

6. Idempotent consumers deduplicate by event ID (`processedIDs` map). Duplicate delivery returns `false` and is ignored, leaving the received count unchanged. Verified by demo Scenario 3.

7. Cleanup is mandatory for correctness and performance — `PurgeProcessedOutbox` removes `PROCESSED` records while retaining `PENDING` ones. Outbox tables grow without bound otherwise. Verified by `TestTransactionalOutbox_PurgeProcessed`.

8. Concurrency is safe in this implementation — 10 concurrent worker goroutines producing 100 orders result in exactly 100 published messages with zero pending and no data races under `go test -race`.

9. Monitoring the oldest unprocessed outbox event age is the leading indicator of relay or broker degradation. The lab's illustrative thresholds (~2 seconds normal, 47+ minutes abnormal) are examples, not universal constants — tune to your own SLOs.

10. This lab is a simplified, in-memory model: it uses a transactional memory DB instead of PostgreSQL/MySQL, implements only Polling Publisher (not CDC/log-tailing), and has no dead-letter queue for non-retryable payloads.
