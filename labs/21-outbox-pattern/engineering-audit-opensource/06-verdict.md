# Engineering Audit Verdict

Target Lab: labs/21-outbox-pattern
Audit Date: 2026-09-27

## Summary

Code Files Reviewed:
- internal/outbox/{model,db,broker,service,relay,consumer}.go
- cmd/demo/main.go
- tests/outbox_test.go
Tests Reviewed: 6 (all in tests/outbox_test.go)
Commands Executed: go test -count=1 ./... ; go test -count=1 -race ./... ; go run ./cmd/demo
Failures: 0
Warnings: 4 (see Non-Blocking Issues)

## Quality Gates

Compilation: PASS (go build ./... exit 0)
Tests: PASS (ok github.com/software-engineering-lab/labs/21-outbox-pattern/tests 0.212s)
Race Detector: PASS (ok ... 1.227s, zero races)
Demo: PASS (output matches demo expectations; atomic write, dispatch, idempotent duplicate, dual-write inconsistency shown)
Research Alignment: WARNING (design.md describes SQLite tables; actual impl is in-memory map DB — intentional per impl notes, but design.md stale)
Documentation Accuracy: WARNING (README aligned to code; engineering-design.md out of sync re: SQLite, Broker->Consumer wire, retry test coverage claim)

## Blocking Issues
1. None

## Non-Blocking Issues
1. Relay Publish-then-Mark is non-atomic; crash between Publish success and Mark processed leaves message PENDING -> DUPLICATE publish on retry (acceptable at-least-once semantics, documented in design as expected).
2. Relay.Start() / Relay.Stop() not lifecycle-safe: double-close of stopChan panics; multiple Start() leaks goroutines (MEDIUM severity latent bug).
3. design.md claims SQLite table schema while code uses in-memory maps; design.md architecture shows Broker->Consumer delivery leg that does not exist in code (consumer pulls published list manually).
4. No test exercises relay retry after transient broker failure (design.md test strategy claims retry mechanism coverage).

## Required Revisions
1. (Before full release) Harden Relay.Stop() with sync.Once and guard Start() against double-invoke — or document single-lifecycle usage.
2. (Future) Align engineering/01-design.md with chosen in-memory implementation: remove SQLite tables or implement them.
3. (Future) Add test for relay recovery after broker failure: inject SetFailNext on first poll, then allow success, assert message reaches consumer after retry.
4. (Future) Clarify relay->consumer delivery responsibility in README/architecture (currently consumer reads broker.GetPublished() manually).

## Final Status

APPROVED_WITH_WARNINGS

Code compiles, tests pass, race detector clean, demo output genuine, core claims (atomicity, idempotency, dual-write flaw, rollback, purge) proven by tests. Warnings are documentation drift and latent relay lifecycle safety; no HIGH/CRITICAL blocking issues. Technical Writer may proceed but should be aware design.md requires reconciliation with implementation and that Broker->Consumer wiring is conceptual, not literal.
