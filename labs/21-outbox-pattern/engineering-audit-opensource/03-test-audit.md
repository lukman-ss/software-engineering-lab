# Engineering Test Audit

Audited File:
- tests/outbox_test.go

Test Suite Overview:
5 unit tests covering happy path, rollback, idempotency, dual-write flaw, and concurrent writes. The race detector passes cleanly for the concurrent test.

## Finding 1: Happy Path Test

Location: TestTransactionalOutbox_HappyPath
Coverage: 
- Atomic persistence: service.CreateOrderWithOutbox stages order+outbox and commits
- Pre-relay state: verifies order persisted in DB before relay runs
- Relay dispatch: after time.Sleep, verifies broker received exactly 1 message
- Outbox status update: verifies outbox message marked PROCESSED
- Idempotent processing: consumer.Handle returns true for first delivery
- Consumer state: verifies received count = 1
Assessment: PASS
Strengths: 
- Verifies each step of the pipeline in order.
- Uses time.Sleep to allow relay to act (adequate for unit test).
Notes: None.

## Finding 2: Rollback Test

Location: TestTransactionalOutbox_Rollback
Coverage: 
- Transaction rollback: Tx.Rollback() called after staging order+outbox
- DB state: verifies no order persisted in db.orders
- Outbox state: verifies no outbox persisted in db.outbox
- Broker state: verifies zero messages published (time.Sleep to let relay try)
Assessment: PASS
Strengths: 
- Tests rollback at the Tx level (same code path service uses on error).
- Confirms nothing leaked to broker despite relay running.
Notes: Uses direct Tx calls rather than service-layer error paths; service-layer rollback is covered indirectly via marshal/save error paths in other tests.

## Finding 3: Idempotency Test

Location: TestTransactionalOutbox_Idempotency_DuplicateDelivery
Coverage: 
- First delivery: consumer.Handle returns true (processed)
- Duplicate delivery: same msg.ID returns false (skipped)
- Consumer state: verifies received count remains 1
Assessment: PASS
Strengths: 
- Pure unit test of consumer idempotency with zero external dependencies.
- Matches the behavior needed to handle relay redelivery.
Notes: Does not test the relay's actual redelivery scenario (e.g., relay crashes after publish but before marking processed); simulates it via duplicate Handle call.

## Finding 4: Dual-Write Flaw Test

Location: TestDualWriteProblem_Failure
Coverage: 
- Setup: broker.SetFailNext(true) to simulate broker down
- Service call: CreateOrderDualWriteNaive expects error from broker failure
- DB state: verifies order persisted in db.orders (broker failure after DB commit)
- Broker state: verifies zero messages published
- Inconsistency note: comments that system is now inconsistent (order exists, event never sent)
Assessment: PASS
Strengths: 
- Directly demonstrates the dual-write vulnerability the outbox pattern solves.
- Matches demo Scenario 1 and engineering notes description.
Notes: Error message matches broker's "broker unavailable" exactly.

## Finding 5: Concurrent Writes Test

Location: TestTransactionalOutbox_ConcurrentWrites
Coverage: 
- Concurrency: 10 goroutines × 10 orders each = 100 CreateOrderWithOutbox calls
- Race detector: test exists to allow go test -race to run with concurrent load
- Relay: started with short poll interval (5ms)
- Wait: time.Sleep after WaitGroup to let relay catch up
Assessment: WEAK
Weaknesses:
- **No correctness assertions**: test ends with no checks on persisted orders, outbox messages, broker publishes, or consumer state.
- **Identical keys**: all 100 calls use the same order ID "o-concurrent-1". Due to map overwrites, only the last commit for that key survives (effectively testing 1 order, not 100).
- **Does not verify relay throughput**: no check that all (or any) staged messages were eventually published.
Missed Opportunities:
- Assert that number of unique orders persisted equals expected distinct IDs (if varied).
- Assert that broker received at least 1 message (if IDs differed, would be up to 100).
- Use varied order IDs to test concurrent distinct-key commits.
- Assert outbox status transitions for published messages.
Strength:
- The test does successfully exercise the race detector (go test -race passes).

## Finding 6: Missing Test Coverage

Location: N/A (omitted tests)
Claimed Behavior: engineering/02-implementation-notes.md Test Strategy lists:
  - Atomic Tx: Order creation + Outbox event insertion. ✓ (HappyPath, Rollback)
  - Rollback handling on invalid order. ✓ (Rollback)
  - Relay & Idempotency Tests: Outbox polling dispatch and marking processed; duplicate message delivery handling; broker failure retry mechanism.
  - Concurrency & Race Detector: Multiple concurrent order creation requests and relay polling cycles; Go race detector verification.
Observed Gaps:
- **Broker failure retry mechanism**: no test where the relay encounters a broker failure during PollAndDispatch and later succeeds after recovery.
  - e.g., SetFailNext(true) during relay operation, then SetFailNext(false) and verify message eventually publishes and is marked processed.
- **MarkOutboxProcessed failure path**: no test where broker.Publish succeeds but db.MarkOutboxProcessed fails (simulate error) and verify message remains PENDING and is retried on next tick.
- **Cleanup verification**: no test for processed record purging (feature not implemented).
- **Transactional service-layer rollback**: Rollback test uses Tx directly; no test where service.CreateOrderWithOutbox returns error (e.g., marshal failure) and verifies nothing persisted.
- **Exactly-once semantics under ideal conditions**: no test verifying that N successful CreateOrderWithOutbox calls result in exactly N broker messages and N consumer receptions (concurrent test lacks assertions).
Assessment: WEAK
Severity: MEDIUM
Notes: Passing test suite does not guarantee comprehensive coverage. Core behaviors are tested (happy path, rollback, idempotency, dual-write flaw, race freedom) but several claimed test strategies are missing or weak.