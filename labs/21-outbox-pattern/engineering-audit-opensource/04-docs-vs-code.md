# Docs vs Code Audit

Scope compared:
- README.md (user-facing)
- engineering/01-design.md (design spec)
- engineering/02-implementation-notes.md (implementation notes)
- engineering/03-execution-result.md (recorded execution results)
- internal/outbox/*.go (implementation)
- tests/outbox_test.go (tests)
- cmd/demo/main.go (demo)
- Verified runtime: go build, go test, go test -race, go run ./cmd/demo

## 1. README vs Code — Alignment

README claim: "In-memory transactional database simulating BeginTx, Commit, and Rollback across orders and outbox records."
Code reality: db.go implements a `DB` with `orders`/`outbox` maps and a `Tx` struct that stages and atomically commits/rolls back. Matches.
Result: PASS

README component descriptions (db.go broker.go service.go relay.go consumer.go) all match actual file responsibilities. Result: PASS

README top-line claim: "proving atomic persistence between domain entities and event logs, decoupled polling relay dispatch to a message broker, and downstream consumer idempotency."
All three are implemented and proven by tests + demo. Result: PASS

## 2. README execution commands

README specifies: `go test ./...`, `go test -race ./...`, `go run ./cmd/demo`.
All run as documented (see 06-verdict). Result: PASS

## Finding 1

Location: README.md (line 3, 9) vs engineering/01-design.md (lines 36-37, 50)
Claimed Behavior: engineering/01-design.md components describe a "SQLite Database" / "SQLite DB: Embedded SQL engine storing business data (orders) and outbox queue (outbox_events)."
Observed Implementation: No SQLite; DB is sync.RWMutex-guarded maps (internal/outbox/db.go). go.mod has no SQLite dependency. README accurately describes an in-memory DB.
Type: DOC_CODE_MISMATCH (design doc vs code)
Severity: MEDIUM
Notes: README is correct; only the internal design doc diverges. Transactional semantics are faithfully preserved via the map-based Tx. The simplification is acknowledged in engineering/02-implementation-notes.md ("In-Memory Transactional DB"), but the design doc's SQLite claim remains unmet.

## Finding 2

Location: engineering/01-design.md (lines 10, 21-26, 52) vs internal/outbox/relay.go
Claimed Behavior: Expected Behavior #5 "Outbox cleanup job removes or archives processed outbox records after retention interval"; Success Criteria #5 "Cleanup worker successfully purges processed events"; Component "OutboxRelay ... with retry / cleanup support."
Observed Implementation: relay.PollAndDispatch only transitions PENDING→PROCESSED via db.MarkOutboxProcessed. No delete/purge/cleanup API exists in db.go (grep-confirmed). README makes no cleanup claim.
Type: DOC_CODE_MISMATCH (design doc vs code)
Severity: MEDIUM
Notes: An aspirational success criterion with no implementation or test. README does not over-claim, so user-facing docs stay accurate.

## Finding 3

Location: engineering/01-design.md (line 39) vs internal/outbox/relay.go
Claimed Behavior: Architecture diagram shows "Poll (FOR UPDATE / Lock)" implying row locking to prevent duplicate dispatch across concurrent relay instances.
Observed Implementation: relay.PollAndDispatch reads a snapshot via db.GetPendingOutbox (RLock) then publishes each message; no per-message lock or claim step. Only a single relay instance is used in demo/tests; README does not promise multi-instance relay.
Type: DOC_CODE_MISMATCH (design doc vs code)
Severity: LOW
Notes: Out of the lab's single-instance scope; idempotency absorbs any duplicates.

## Finding 4

Location: engineering/01-design.md Test Strategy (lines 61-62, 64-66) vs tests/outbox_test.go
Claimed Behavior: Test strategy lists "Broker failure retry mechanism" under Relay & Idempotency Tests.
Observed Implementation: No test exercises the relay encountering a broker Publish failure and recovering (no SetFailNext-based relay failure test). Covered test areas (happy path, rollback, idempotency, dual-write, concurrency) match engineering/02-implementation-notes.md's list.
Type: TEST_CLAIM_MISMATCH
Severity: MEDIUM
Notes: Stated test strategy partially unmet; see 03-test-audit Finding 6.

## Finding 5

Location: engineering/03-execution-result.md vs actual runtime execution
Claimed Behavior: Build passes (exit 0); go test ./... ok; go test -race ./... ok; demo output as recorded.
Observed Implementation: Auditor independently ran all four commands and observed passing build, passing tests (0.222s), passing race (1.230s), and demo output identical to the recorded 03-execution-result.md line-for-line (same messages, same broker counts, same consumer accepted=1/accepted=false, total processed=1). Only timing differs by run.
Type: RESEARCH/RESULT VERIFICATION
Assessment: PASS (no fabricated result)
Notes: The recorded execution result is an honest, reproducible snapshot. Minor timing variance (0.497s vs 0.222s; 1.487s vs 1.230s) is normal machine-to-machine difference, not fabrication.

## Summary

| Artifact | vs Code | Status |
|---|---|---|
| README user-facing claims | code | PASS — accurate |
| engineering/01-design.md SQLite | code (maps) | MEDIUM mismatch |
| engineering/01-design.md cleanup worker | code (absent) | MEDIUM mismatch |
| engineering/01-design.md FOR UPDATE locking | code (no lock) | LOW mismatch |
| engineering/01-design.md retry test | tests (absent) | MEDIUM mismatch |
| engineering/03-execution-result.md | runtime | PASS — reproducible |