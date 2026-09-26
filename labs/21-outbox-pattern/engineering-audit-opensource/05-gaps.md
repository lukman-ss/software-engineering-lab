# Gaps Analysis

Target: labs/21-outbox-pattern
Pipeline scope: implementation and tests only (research content not audited).

## GAP-1

Type: MISSING_TEST
Location: tests/outbox_test.go (absence); relay retry path in internal/outbox/relay.go:43-59
Description: No automated test verifies the relay's at-least-once retry behavior. The design doc (engineering/01-design.md:20) and success criteria (line 27) state "thread-safe implementations with zero race conditions," and the failure scenario (line 20) describes "Relay Crash / Duplicate Delivery" — relay publishes then crashes before updating outbox status, then republishes. The demo (cmd/demo/main.go:54-59) simulates a duplicate delivery to the *consumer* by re-handling `published[0]`, but no test asserts that the *relay itself* re-publishes a `PENDING` message after a broker failure and leaves it pending for retry (relay.go:56 logs and leaves message pending on publish error).
Severity: MEDIUM
Impact: The at-least-once recovery guarantee (code audit Finding 4) is documented and implemented but is unproven by the test suite.
Proof needed: Test that sets `broker.SetFailNext(true)`, runs one `PollAndDispatch` (publish fails, message stays PENDING), then clears failure and runs another `PollAndDispatch` (message re-published and marked PROCESSED).

## GAP-2

Type: MISSING_EDGE_CASE
Location: tests/outbox_test.go:141-169 (TestTransactionalOutbox_ConcurrentWrites)
Description: All 100 concurrent transactions use the same order ID `"o-concurrent-1"` (and thus the same event ID), so the test contends on a single key rather than verifying concurrent distinct transactions. After `wg.Wait` it sleeps and asserts nothing — no check of persisted count, published count, or deduplication.
Severity: MEDIUM
Impact: The concurrency success criterion (design doc:01-design.md:50, 01-execution-result.md) is only proven for "no data race / no crash," not for correctness (each distinct order published exactly once). A lost-update or duplicate-publish bug across distinct orders would not be caught.
Proof needed: Concurrently create distinct order IDs and assert `len(broker.GetPublished())` equals the number of distinct orders processed.

## GAP-3

Type: UNHANDLED_ERROR (latent)
Location: internal/outbox/relay.go:39-41 (Relay.Stop)
Description: `Stop()` calls `close(r.stopChan)` with no guard against a second close. Calling `Stop()` twice panics with "close of closed channel." All current callers invoke `Stop()` exactly once (single `defer`), so this is not triggered, but it is a latent concurrency hazard if the relay lifecycle is reused.
Severity: LOW
Impact: No test failure observed. Latent crash under double-shutdown usage.
Upgrade path: Use `sync.Once` or a `stopped bool` check under a guard.

## GAP-4

Type: MISSING_EDGE_CASE
Location: tests/outbox_test.go (absence); internal/outbox/consumer.go:18-31
Description: `Consumer.Handle` is documented as thread-safe (mutex-guarded) and the success criteria require "thread-safe implementations," but no test exercises concurrent `Handle` calls. The single idempotency test is sequential.
Severity: LOW
Impact: Consumer concurrency under load is unproven, though the mutex makes it very likely correct.
Proof needed: Concurrently call `Handle` for distinct IDs plus one shared duplicate ID and assert exactly one acceptance per distinct ID.

## GAP-5

Type: MISSING_EDGE_CASE
Location: internal/outbox/relay.go:43-59; tests/outbox_test.go (absence)
Description: The relay has no poison-message handling or backoff (noted as a limitation in engineering/02-implementation-notes.md:29). A message that always fails to publish is retried every poll interval indefinitely. No test asserts retry count bounds or logging behavior, but this is an accepted limitation, not a defect.
Severity: LOW (informational; explicitly out of scope per design)
Notes: Listed under Known Limitations (implementation-notes.md:29); not promoted to a defect because polling-over-CDC was chosen by design.

## Summary

| Gap | Type                | Severity | Triggered? |
|-----|---------------------|----------|------------|
| GAP-1 | MISSING_TEST       | MEDIUM   | No (feature implemented) |
| GAP-2 | MISSING_EDGE_CASE  | MEDIUM   | No (feature implemented) |
| GAP-3 | UNHANDLED_ERROR    | LOW      | No (latent) |
| GAP-4 | MISSING_EDGE_CASE  | LOW      | No |
| GAP-5 | MISSING_EDGE_CASE  | LOW      | No (accepted limitation) |

No BROKEN_IMPLEMENTATION, DOC_CODE_MISMATCH, RACE_CONDITION (race detector clean), FAKE_DEMO, FAKE_BENCHMARK, or UNVERIFIED_RESULT issues found.
