# Engineering Code Audit

Audited Files:
- internal/outbox/model.go
- internal/outbox/db.go
- internal/outbox/service.go
- internal/outbox/broker.go
- internal/outbox/relay.go
- internal/outbox/consumer.go
- cmd/demo/main.go

## Finding 1

Location: internal/outbox/db.go:99-116 (Tx.Commit)
Claimed Behavior: Order record and outbox event are persisted atomically within a single transaction; an observer sees either both or neither.
Observed Implementation: Commit sets tx.closed=true then locks db.mu and copies the entire stagedOrders and stagedOutbox maps into db.orders/db.outbox in one critical section. Rollback sets closed=true and discards the staged maps without touching db.
Assessment: PASS
Severity: N/A
Notes: Atomic apply via a single mutex critical section. Not a real SQL transaction, but semantically equivalent atomic persistence for the demo. go test -race clean.

## Finding 2

Location: internal/outbox/db.go:118-127 (Tx.Rollback); service.go:30,42-52
Claimed Behavior: On rollback (marshal error, save error, invalid input), neither the order nor the outbox event is persisted.
Observed Implementation: service.CreateOrderWithOutbox calls tx.Rollback() on marshal/save failures and returns before Commit(). Rollback discards staged maps; Commit after Rollback returns ErrTxClosed.
Assessment: PASS
Severity: N/A
Notes: TestTransactionalOutbox_Rollback proves both records remain absent after rollback.

## Finding 3

Location: internal/outbox/db.go Commit/Commit lock order
Claimed Behavior: Concurrent transactions are race-free (no data races, no deadlock).
Observed Implementation: Lock order is fixed: Tx methods acquire tx.mu then, in Commit only, db.mu. DB read/write helpers (GetOrder, GetOutbox, GetPendingOutbox, MarkOutboxProcessed) acquire db.mu only and never touch tx.mu. No code path acquires db.mu then tx.mu, so the order is deadlock-free. db.mu serializes all Commit writes.
Assessment: PASS
Severity: N/A
Notes: `go test -race ./...` reported no data races.

## Finding 4

Location: internal/outbox/relay.go:43-60 (PollAndDispatch); db.go:47-69; consumer.go:19-31
Claimed Behavior: Zero lost events under broker/network failure; duplicates tolerated by idempotent consumer (at-least-once delivery).
Observed Implementation: GetPendingOutbox returns only PENDING messages. On Publish success + MarkOutboxProcessed success the message becomes PROCESSED and is no longer fetched. On Publish failure or MarkOutboxProcessed failure, the message stays PENDING and is re-delivered on the next poll tick. Consumer dedups by msg.ID and returns false on duplicates.
Assessment: PASS
Severity: N/A
Notes: Correct at-least-once semantics. MarkOutboxProcessed failure leaves the message PENDING by design rather than silently dropping it.

## Finding 5

Location: internal/outbox/service.go:57-90 (CreateOrderDualWriteNaive)
Claimed Behavior: Naive dual write commits the DB transaction and then publishes to the broker outside the transaction; a broker failure leaves the DB committed but the event lost (state inconsistency).
Observed Implementation: The order is staged+committed to db first, then broker.Publish is called. On Publish failure the error is wrapped and returned; the DB commit is NOT rolled back, so the order persists while the broker never receives the event.
Assessment: PASS
Severity: N/A
Notes: Intentional flaw demonstration, not a latent bug. Proven by TestDualWriteProblem_Failure and demo Scenario 1.

## Finding 6

Location: internal/outbox/relay.go, internal/outbox/db.go
Claimed Behavior: engineering/01-design.md Expected Behavior #5 and Success Criteria #5 require a cleanup worker that removes/archives processed outbox records after a retention interval. engineering/03-execution-result.md lists the same as a success criterion.
Observed Implementation: Relay only transitions PENDING→PROCESSED. There is no DeleteOutbox/Purge/Cleanup API in db.go (confirmed by grep). Processed records accumulate indefinitely in db.outbox with no eviction path and no test.
Assessment: FAIL
Severity: MEDIUM
Notes: Auxiliary to the core transactional guarantee, but it is an explicit, unmet success criterion in the engineering docs with no code or test coverage. Not a correctness bug in the proven core behavior.

## Finding 7

Location: internal/outbox/db.go (struct DB uses sync.RWMutex + maps); go.mod
Claimed Behavior: engineering/01-design.md architecture diagram and Components describe a "SQLite Database / embedded SQLite engine." Implementation Decisions list modernc.org/sqlite or stdlib database/sql with in-memory SQLite.
Observed Implementation: DB is a struct holding two Go maps guarded by sync.RWMutex. BeginTx returns a Tx with staged map copies; Commit copies staged maps under db.mu. go.mod contains no SQLite/database/sql dependency; grep confirms no sqlite/sql.Open/modernc usage anywhere in the module.
Assessment: WARNING
Severity: MEDIUM
Notes: README accurately describes an "in-memory transactional database simulating BeginTx, Commit, and Rollback," so user-facing docs match the code. The deviation is between the design doc's stated engine and the actual map-based implementation. Transactional semantics are preserved; this is an honest simplification, not a fabrication.

## Finding 8

Location: internal/outbox/relay.go:39-41 (Relay.Stop)
Claimed Behavior: Relay can be safely stopped.
Observed Implementation: Stop() executes close(r.stopChan) with no guard or recovered check.
Assessment: WARNING
Severity: LOW
Notes: Calling Stop() twice panics with "close of closed channel." Not exercised by current tests/demo (each uses a single defer Stop). Fragile API for reuse.

## Finding 9

Location: internal/outbox/relay.go:43-60
Claimed Behavior: engineering/01-design.md lists relay "retry / cleanup support" and the Test Strategy lists a "Broker failure retry mechanism."
Observed Implementation: On Publish failure the message is logged and left PENDING; retry occurs only on a subsequent ticker tick. There is no retry counter, backoff, max-attempts, or circuit breaker, and no test exercises a failure-then-recovery sequence for the relay. Cleanup is absent (see Finding 6).
Assessment: WARNING
Severity: MEDIUM
Notes: Implicit re-poll provides at-least-once delivery, which satisfies the core guarantee. The broader "retry mechanism" claim (backoff/circuit-breaker) is only minimally realized and untested.

## Finding 10

Location: internal/outbox/relay.go:43-59
Claimed Behavior: engineering/01-design.md architecture diagram implies Poll (FOR UPDATE / Lock) to prevent duplicate dispatch across concurrent relay instances.
Observed Implementation: GetPendingOutbox returns a snapshot read under RLock; each message is then individually Publish'd and MarkOutboxProcessed'd. There is no per-message locking or claiming step. With a single relay instance (the demo/tests scope) this is correct; with multiple concurrent relay instances two workers could fetch and dispatch the same PENDING message before either marks it processed.
Assessment: WARNING
Severity: LOW
Notes: README does not promise multi-instance relay; the lab consistently deploys a single relay. Idempotency absorbs any duplicate, so data is not corrupted. Out of scope for the current demo but noted.

## Code Quality Notes

- No unnecessary abstractions or dead code observed in core paths.
- Error handling: service wraps relay/broker errors with fmt.Errorf/%w; Tx errors propagate to callers.
- demo/main.go: uses time.Sleep-based synchronization to wait for relay ticks (flaky-prone in CI but adequate for a local demo).
