# Test Audit

Target Lab: labs/21-outbox-pattern
Test File: tests/outbox_test.go

Run Results:
- `go test -count=1 ./...` — PASS (7 tests, exit 0)
- `go test -race ./...` — PASS (1.334s, exit 0)

---

## Test 1: TestTransactionalOutbox_HappyPath

Coverage: happy path (atomic save → relay dispatch → broker receive → status PROCESSED → consumer idempotent process)
Assessment: PASS
Notes: Verifies complete happy path including order persistence, relay polling, broker publishing, outbox status transition, and consumer processing.

## Test 2: TestTransactionalOutbox_Rollback

Coverage: failure path / rollback (staged data discarded on rollback, no persistence, no dispatch)
Assessment: PASS
Notes: Verifies that Rollback() discards both order and outbox message, and relay has nothing to dispatch.

## Test 3: TestTransactionalOutbox_Idempotency_DuplicateDelivery

Coverage: idempotency (duplicate message rejected by consumer)
Assessment: PASS
Notes: First delivery accepted, second delivery rejected. Received count remains 1.

## Test 4: TestDualWriteProblem_Failure

Coverage: negative case / failure path (dual-write inconsistency when broker fails)
Assessment: PASS
Notes: Broker fails after DB commit. Verifies order persists in DB but message not published. Demonstrates the problem the outbox pattern solves.

## Test 5: TestTransactionalOutbox_ConcurrentWrites

Coverage: concurrency safety (10 workers × 10 writes under race detector)
Assessment: WARNING
Severity: MEDIUM
Notes: All 10 workers use SAME orderID "o-concurrent-1" and SAME customerID "c-multi". The generated message ID `evt-o-concurrent-1` collides for all 100 calls. This means:
  - Orders overwrite each other (last writer wins in db.orders map under lock)
  - Only one outbox message survives (same ID, overwritten in map)
  - Relay only publishes 1 message total
  - No assertions on outcome — test only checks for data races, not correctness
  The race detector passes (no data race), but the test does NOT meaningfully verify concurrent write correctness. Using unique IDs would be more representative.
  PASS under -race but weak coverage of concurrent write semantics.

## Test 6: TestTransactionalOutbox_PurgeProcessed

Coverage: cleanup / edge case (purge removes only PROCESSED records)
Assessment: PASS
Notes: Verifies purge count and that PENDING record remains while PROCESSED record is deleted.

## Test 7: TestTransactionalOutbox_RelayRetryAfterBrokerFailure

Coverage: recovery / retry (relay retries after broker failure, eventually succeeds)
Assessment: PASS
Notes: Simulates broker failure via SetFailNext(true). First poll fails, message stays PENDING. Second poll succeeds. Tests Start()/Stop() idempotency. Verifies end-to-end retry recovery.

## Coverage Analysis

Coverage Matrix:
| Scenario              | Covered? | Test Name                                    |
|-----------------------|----------|----------------------------------------------|
| Happy path            | YES      | TestTransactionalOutbox_HappyPath            |
| Rollback              | YES      | TestTransactionalOutbox_Rollback             |
| Idempotency           | YES      | TestTransactionalOutbox_Idempotency_DuplicateDelivery |
| Dual-write failure    | YES      | TestDualWriteProblem_Failure                 |
| Concurrency           | PARTIAL  | TestTransactionalOutbox_ConcurrentWrites     |
| Relay retry            | YES      | TestTransactionalOutbox_RelayRetryAfterBrokerFailure |
| Purge/cleanup          | YES      | TestTransactionalOutbox_PurgeProcessed       |

## Gaps Identified

1. MISSING_TEST: No test verifies concurrent reads/writes to Consumer under race detector — Consumer has mutex but concurrency not exercised under -race.
2. MISSING_TEST: Concurrent relay + purge interaction not tested — PurgeProcessedOutbox could run while relay marks messages.
3. MISSING_TEST: No test for MarkOutboxProcessed on non-existent message (error path returns error but untested).
4. MISSING_TEST: No test for Commit() called twice (ErrTxClosed error path untested).
5. WEAK_TEST: TestTransactionalOutbox_ConcurrentWrites uses same IDs and no assertions on outcome — only race detection, not correctness.