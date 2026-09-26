# Docs vs Code Audit

## README Claim 1: db.go simulates BeginTx, Commit, Rollback across orders and outbox records
Code reality: db.go implements Tx with BeginTx, Commit, Rollback. Staged orders and outbox messages are committed atomically under db.mu. Rollback discards staged data.
Status: DOC_CODE_MATCH

## README Claim 2: broker.go is a thread-safe mock message broker simulating publish failures and event reception
Code reality: MockBroker uses sync.RWMutex protecting published slice and failNext flag. Publish simulates failure. GetPublished returns copy.
Status: DOC_CODE_MATCH

## README Claim 3: service.go compares naive dual-write vs atomic outbox writes
Code reality: OrderService has CreateOrderWithOutbox (atomic, uses Tx) and CreateOrderDualWriteNaive (non-atomic, DB commit then broker publish). Both implemented.
Status: DOC_CODE_MATCH

## README Claim 4: relay.go is an asynchronous polling worker querying pending outbox records and dispatching to broker
Code reality: Relay.Start spawns goroutine, polls via ticker, calls PollAndDispatch which queries GetPendingOutbox and publishes. Stop via stopChan.
Status: DOC_CODE_MATCH

## README Claim 5: consumer.go enforces idempotency through event ID tracking
Code reality: Consumer tracks processedIDs map, rejects duplicates. Handle returns false for duplicates.
Status: DOC_CODE_MATCH

## README Claim: Running Tests - go test ./... and go test -race ./...
Code reality: Both commands executed successfully. All tests pass. Race detector finds no issues.
Status: DOC_CODE_MATCH

## README Claim: Running Demo - go run ./cmd/demo
Code reality: Demo executes successfully, demonstrating all three scenarios:
- Scenario 1: Dual-write inconsistency (order in DB, 0 broker messages)
- Scenario 2: Transactional outbox (order + outbox atomically saved, broker receives 1 message, consumer accepts)
- Scenario 3: At-least-once delivery & idempotent consumer (duplicate rejected)
Status: DOC_CODE_MATCH

## Research Claim: Atomic persistence, decoupled polling relay, downstream consumer idempotency
Code reality: Atomic persistence via shared Tx (Finding 3 PASS). Decoupled polling relay (Finding 5 PASS). Idempotent consumer (Finding 7 PASS).
Status: DOC_CODE_MATCH

Overall documentation accuracy: PASS - README accurately describes the implementation, commands work as documented, demo output matches described scenarios.