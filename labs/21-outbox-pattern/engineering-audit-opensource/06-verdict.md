# Engineering Audit Verdict

Target Lab: labs/21-outbox-pattern
Audit Date: 2026-09-26

## Summary

Code Files Reviewed:
- internal/outbox/model.go
- internal/outbox/db.go
- internal/outbox/service.go
- internal/outbox/broker.go
- internal/outbox/relay.go
- internal/outbox/consumer.go
- cmd/demo/main.go

Tests Reviewed:
- tests/outbox_test.go

Commands Executed:
- go build ./... → PASS (exit code 0)
- go test ./... → PASS (ok, tests)
- go test -race ./... → PASS (ok, tests, no races)
- go run ./cmd/demo → PASS (output matches engineering/03-execution-result.md exactly)

Failures: None
Warnings: See Non-Blocking Issues below.

## Quality Gates

Compilation: PASS
Tests: PASS
Race Detector: PASS
Demo: PASS
Research Alignment: WARNING (design doc overclaims SQLite engine and cleanup worker relative to implementation; implementation notes document simplifications)
Documentation Accuracy: WARNING (README accurate; design doc contains unmet claims about cleanup and SQLite engine)

## Blocking Issues
1. None (no HIGH or CRITICAL severity issues)

## Non-Blocking Issues
1. MEDIUM: Missing cleanup worker (design claims it; code lacks implementation or tests) [DOC_CODE_MISMATCH]
2. MEDIUM: Database engine mismatch (design claims SQLite; code uses in-memory map DB) [DOC_CODE_MISMATCH]
3. MEDIUM: Weak concurrent writes test (no correctness assertions; identical keys) [MISSING_TEST]
4. MEDIUM: Missing broker failure recovery test for relay [MISSING_TEST]
5. MEDIUM: Missing MarkOutboxProcessed failure path test [MISSING_TEST]
6. MEDIUM: Missing service-layer rollback test [MISSING_TEST]
7. MEDIUM: Missing exactly-once semantics test [MISSING_TEST]
8. LOW: Relay Stop() not idempotent (double close panics) [MISSING_EDGE_CASE]
9. LOW: Design overclaim on relay locking (FOR UPDATE vs single-instance) [DOC_CODE_MISMATCH]

## Required Revisions
1. Either implement cleanup worker for processed outbox records or remove the claim from engineering/01-design.md (Expected Behavior #5, Success Criteria #5, Component description) and related success criteria in engineering/03-execution-result.md.
2. Align database documentation: either switch to SQLite/database/sql per design doc OR update engineering/01-design.md and engineering/02-implementation-notes.md to reflect the intentional in-memory map DB simplification (noting it preserves atomic Tx semantics).
3. Make Relay.Stop idempotent by adding a guard against double-close (e.g., atomic boolean or closed channel check) OR document that Stop must be called only once.
4. Strengthen TestTransactionalOutbox_ConcurrentWrites: vary order IDs to test distinct-key commits, add assertions on persisted orders, outbox messages, broker publishes, and consumer received count.
5. Add a test where the relay encounters a broker Publish failure (e.g., via temporary SetFailNext(true)) and later recovers to verify message eventual delivery and processing.
6. Add a test where broker.Publish succeeds but db.MarkOutboxProcessed fails (simulated error) and verify the message remains PENDING for redelivery and is not lost.
7. Add a test where service.CreateOrderWithOutbox returns an error (e.g., JSON marshal failure) and verify that neither order nor outbox event is persisted (service-layer rollback).
8. Add a test verifying that N successful CreateOrderWithOutbox calls (with distinct keys) result in exactly N broker messages and N consumer receptions (ideal-condition exactly-once).
9. Either implement row-level locking or claiming for multi-instance relay safety (to match design's "FOR UPDATE / Lock") OR update engineering/01-design.md to clarify the single-instance scope and remove the locking implication from the architecture diagram.

## Final Status

APPROVED_WITH_WARNINGS