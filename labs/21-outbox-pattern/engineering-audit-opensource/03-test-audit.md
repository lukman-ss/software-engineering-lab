# Test Audit

Target: labs/21-outbox-pattern
Test file: tests/outbox_test.go

Execution results:
```
go test ./...                       -> ok  tests  (cached)
go test -race -v ./...              -> ok  tests  (PASS, 5/5)
```

Tests present (tests/outbox_test.go):
1. `TestTransactionalOutbox_HappyPath` — happy path.
2. `TestTransactionalOutbox_Rollback` — rollback path.
3. `TestTransactionalOutbox_Idempotency_DuplicateDelivery` — duplicate delivery / idempotency.
4. `TestDualWriteProblem_Failure` — naive dual-write failure path.
5. `TestTransactionalOutbox_ConcurrentWrites` — concurrency / race check.

Coverage matrix:

| Coverage area        | Covered? | Where                          |
|----------------------|----------|--------------------------------|
| Happy path (atomic outbox -> relay -> broker -> consumer) | YES | TestTransactionalOutbox_HappyPath |
| Rollback / transaction abort              | YES | TestTransactionalOutbox_Rollback |
| Marshal failure rollback                  | NO  | (not exercised; `SaveOrder`/`SaveOutbox` cannot fail here) |
| Dual-write failure (broker down, order left committed) | YES | TestDualWriteProblem_Failure |
| Broker failure during relay dispatch + retry | NO  | demo only, no test |
| Idempotency / duplicate delivery          | YES | TestTransactionalOutbox_Idempotency_DuplicateDelivery |
| Concurrency (distinct orders)             | PARTIAL/NO | TestTransactionalOutbox_ConcurrentWrites (see below) |
| At-least-once retry (pending message re-published) | NO | no test |
| Consumer concurrency                      | NO   | no concurrent Handle test |
| Edge: empty DB / no pending               | NO   | no test |

## Finding T1

Location: tests/outbox_test.go:141-169 (TestTransactionalOutbox_ConcurrentWrites)
Claimed Behavior: Concurrent transactional writes are safe.
Observed Implementation: 10 workers each call `CreateOrderWithOutbox` 10 times, all with the **same** order ID `"o-concurrent-1"` (and therefore the same event ID `evt-o-concurrent-1`). After `wg.Wait` it sleeps 50ms and asserts **nothing**.
Assessment: WARNING — test is structurally present but functionally weak.
Severity: MEDIUM
Notes: Because every transaction writes the same map key, the test exercises contention on a single key rather than concurrent distinct transactions; it cannot detect a lost-update or duplicate-publish bug across distinct orders. No post-condition is asserted (not even `len(broker.GetPublished())`), so the test only proves "no crash / no data race" — it does not prove the outbox dispatches each message exactly once. The race detector still passes, so there is no data-race defect, but the concurrency behavior is not actually verified.

## Finding T2

Location: tests/outbox_test.go (absence)
Claimed Behavior: (none explicitly)
Observed Implementation: No test covers the relay's retry-on-failure behavior (broker `failNext` set during dispatch) or verifies that a publish failure leaves the message `PENDING` for retry and that it succeeds on the next poll.
Assessment: WARNING — missing failure-path test.
Severity: MEDIUM
Notes: The demo shows this conceptually, but no automated test asserts it. The at-least-once retry guarantee (Finding 4 in code audit) is therefore unproven by tests.

## Finding T3

Location: tests/outbox_test.go:61-90 (TestTransactionalOutbox_Rollback)
Claimed Behavior: Rollback discards staged mutations.
Observed Implementation: Stages order + outbox, calls `Rollback()`, asserts neither is retrievable and no message was published.
Assessment: PASS
Severity: LOW
Notes: Correctly covers the rollback path. (Does not cover the marshal-error rollback path inside `CreateOrderWithOutbox`, but that path cannot fail for the inputs used.)

## Finding T4

Location: tests/outbox_test.go:11-59 (TestTransactionalOutbox_HappyEnd-to-end)
Claimed Behavior: End-to-end atomic outbox flow.
Observed Implementation: Verifies order persisted, broker received 1 message with correct type, outbox marked `PROCESSED`, and consumer processed it.
Assessment: PASS
Severity: LOW
Notes: Relies on a `time.Sleep(50ms)` for relay timing, which is acceptable for a lab but is a mild flakiness source. No race or correctness defect observed.

## Summary

- 4/5 tests assert meaningful post-conditions (PASS).
- 3 of 5 coverage gaps are MISSING cases (broker-failure retry, at-least-once retry, consumer concurrency) — not defects in implementation, but gaps in proof.
- Race detector: PASS (no data races detected under `-race`).
- Overall test strength: the suite proves happy path, rollback, dual-write failure, and idempotency, but does **not** prove the relay's retry/recovery behavior under broker failure.
