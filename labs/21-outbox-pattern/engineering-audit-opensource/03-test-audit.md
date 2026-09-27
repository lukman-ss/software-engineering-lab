# Test Audit

Tests cover:
- Happy-path transactional outbox: atomic write, relay dispatch, consumer processing, idempotency (TestTransactionalOutbox_HappyPath)
- Tx rollback purity (TestTransactionalOutbox_Rollback)
- Consumer duplicate rejection (TestTransactionalOutbox_Idempotency_DuplicateDelivery)
- Dual-write failure demonstration (TestDualWriteProblem_Failure)
- Concurrent writes race detector sanity (TestTransactionalOutbox_ConcurrentWrites)
- Processed outbox purge (TestTransactionalOutbox_PurgeProcessed)

What tests DO NOT cover:
- Relay double-start / double-close panic (Stop() on closed channel)
- Relay goroutine leak on multiple Start()
- PurgeProcessedOutbox concurrency safety (called while relay polling)
- Broker failure handling within relay retry loop beyond one poll (simulated but not asserted)
- Message ordering guarantees (none claimed)
- Commit-time panics (no DB layer error injection)
- Tx.SaveOrder/Outbox error rollback (already covered indirectly)
- Very long polling intervals / context cancellation
- Stress test of concurrent relay + service + consumer with assertions on final state

Assessment:
- Core claims proven: atomicity, idempotency, dual-write flaw, rollback, purge.
- Concurrency: race detector passes but test weak (shared key, no asserts).
- Edge cases: missing relay lifecycle safety and retry loop resilience.
- Failure paths: broker fail simulated; DB layer assumes no panics.
- Demo scenarios validated: dual-write inconsistency, atomic success, idempotent duplicate.

Overall: Passing test suite validates core behavior but leaves operational concerns untested.