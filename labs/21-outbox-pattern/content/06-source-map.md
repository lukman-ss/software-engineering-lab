# Source Map

## Problem
Research: research/05-report.md:11-16, research/03-evidence.md:13-22
Implementation: internal/outbox/service.go (dual-write naive function), internal/outbox/broker.go (fail simulation)
Tests: tests/outbox_test.go:118-140 (TestDualWriteProblem_Failure)

## Why This Matters
Research: research/05-report.md:6-8 (Executive Summary)
Implementation: README.md:3-4 (Lab description)
Tests: Tests show inconsistency manifests in dual-write failure

## Mental Model
Research: research/05-report.md:23-28 (Finding 2), research/03-evidence.md:23-38 (Evidence 3)
Implementation: internal/outbox/service.go:17-53 (CreateOrderWithOutbox)
Tests: tests/outbox_test.go:12-60 (TestTransactionalOutbox_HappyPath)

## Core Concept
Research: research/03-evidence.md:39-52 (Evidence 4: Outbox Table Structure)
Implementation: internal/outbox/model.go:26-32 (OutboxMessage struct), internal/outbox/db.go:14-16 (orders, outbox maps)

## Failure Scenario
Research: research/05-report.md:13-16 (Finding 1)
Implementation: internal/outbox/service.go:55-90 (CreateOrderDualWriteNaive), internal/outbox/broker.go:28-36 (Publish with failNext)
Tests: tests/outbox_test.go:118-140

## How It Works
Implementation: internal/outbox/db.go:25-31 (BeginTx), 92-129 (Tx methods: SaveOrder, SaveOutbox, Commit, Rollback)
Implementation: internal/outbox/relay.go:50-67 (PollAndDispatch)

## Architecture
Research: research/05-report.md:60-72 (Finding 5: Outbox Table Structure and Payload Design), research/03-evidence.md:39-52 (Evidence 4: columns id, aggregatetype, aggregateid, type, payload)
Implementation: internal/outbox/model.go (id, aggregateid, aggregatetype, type, payload fields)
Implementation: internal/outbox/service.go:28-40 (payload marshal and message creation)

## Implementation
All implementation files under internal/outbox/:
- model.go (domain)
- db.go (transactional in-memory DB)
- service.go (business logic: atomic vs dual-write)
- broker.go (mock broker with failure injection)
- relay.go (polling publisher)
- consumer.go (idempotent consumer)
- cmd/demo/main.go (end-to-end demo)
All tests under tests/outbox_test.go

## Code Walkthrough
Atomic Write (Service): internal/outbox/service.go:17-53
Dual-Write Naive (Service): internal/outbox/service.go:55-90
Tx Staging (DB): internal/outbox/db.go:84-139
Relay Loop: internal/outbox/relay.go:27-41 (Start goroutine), 50-67 (PollAndDispatch)
Consumer Idempotency: internal/outbox/consumer.go:19-31

## What the Tests Prove
Test file: tests/outbox_test.go
- Happy path & atomicity: lines 12-60
- Rollback: lines 62-90
- Idempotency: lines 93-115
- Dual-write flaw: lines 118-140
- Concurrent writes: lines 142-194
- Purge: lines 196-216
- Relay retry: lines 219-264
- Concurrent consumers: lines 266-293

## Recovery / Rollback
Implementation: internal/outbox/db.go:131-139 (Tx.Rollback)
Implementation: internal/outbox/service.go:29-32, 42-45, 47-50 (rollback on error paths)
Tests: TestTransactionalOutbox_Rollback lines 62-90

## Production Considerations
Research: research/05-report.md:74-86 (Finding 6: Cleanup and Monitoring)
Implementation: internal/outbox/db.go:71-82 (PurgeProcessedOutbox)
Implementation: internal/outbox/relay.go:50-67 (retry on publish failure)
Research/evidence: research/03-evidence.md:100-118 (Evidence 9-10: Cleanup and Monitoring)

## Common Mistakes
Research: research/03-evidence.md:90-99 (Evidence 8: Payload Size Warning)
Research: research/05-report.md:69-73 (Finding 5: thin vs fat events trade-offs)
Implementation note: Avoid storing large objects in outbox.Payload

## Case Study
Demo walkthrough: cmd/demo/main.go (full file)
Tests: tests/outbox_test.go (all test functions as case studies)

## Checklist
Derived from success criteria in engineering/01-design.md:21-27:
- [x] 100% atomicity between business state and outbox state
- [x] Zero lost events under broker network disconnect / relay retries
- [x] Zero duplicate processing by idempotent consumers despite at-least-once relay delivery
- [x] Cleanup worker successfully purges processed events
- [x] All unit and concurrency tests pass with zero data races (`go test -race ./...`)

## Key Takeaways
See content/05-key-takeaways.md for final distilled list.

## Sources
Research sources list: research/02-sources.md (6 sources)
Evidence mappings: research/03-evidence.md (each evidence block cites source and URL)