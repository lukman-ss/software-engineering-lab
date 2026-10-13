---
title: "Database Isolation: Why Two Committed Transfers Can Create Money"
published: false
description: "A PostgreSQL wallet lab tracing lost updates, row locks, consistent snapshots, and bounded transaction retries in Go."
tags: go, postgres, database, backend
---

Two wallet transfers can both commit successfully while creating money. In Lab 08, the problem is a specific query pattern: read a balance, validate it in Go, calculate a replacement value, and write that value later. A transaction boundary does not make the original read remain current.

The useful question is whether overlapping transactions preserve the business invariant. This lab approaches that question through a three-account wallet, PostgreSQL isolation experiments, and concurrent tests.

## Start with what must remain true

Alice, Bob, and Charlie each begin with 1,000,000. A transfer should change the distribution of those balances while keeping their sum at 3,000,000. Every balance must also remain non-negative.

The schema implements the row-level invariant:

```sql
balance BIGINT NOT NULL CHECK (balance >= 0)
```

That constraint does not enforce conservation of money across accounts. A broken transfer can leave every individual balance positive while increasing their sum. Correctness therefore needs both checks.

`TransferNaive` uses `sql.LevelReadCommitted`. It begins a transaction, reads the sender, checks sufficient funds, reads the receiver, updates both accounts, inserts an audit record, and commits. These operations share an atomic boundary. Their decisions can still depend on stale application values.

## Follow the overlap instead of the happy path

Transfer A sends 800,000 from Alice to Bob. Transfer B sends 800,000 from Alice to Charlie. Both transactions read Alice at 1,000,000 before either writes. Each concludes that the transfer is allowed.

Each then calculates Alice's replacement balance as 200,000. The naive sender update uses that calculated value rather than deriving the new balance from the current row:

```go
// 4. Overwrite sender balance with calculated value (Lost Update vulnerability)
	_, err = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = $1 WHERE id = $2", fromBalance-amount, fromID)
	if err != nil {
		return fmt.Errorf("update from balance: %w", err)
	}
```

*Excerpt from `transfer.go`; surrounding validation, receiver update, audit, and commit are omitted.*

The second write can wait for the first writer and still overwrite Alice with the same 200,000. Bob and Charlie receive separate credits of 800,000. The expected final state in `TestNaiveTransfer_LostUpdate` is Alice 200,000, Bob 1,800,000, Charlie 1,800,000: total 3,800,000.

The test orchestrates this overlap with channels. Both transactions signal that their reads have completed; the coordinator then releases their writes. It reproduces the read-calculate-write pattern explicitly, rather than calling `TransferNaive` directly. That distinction matters when describing what the test covers.

These numbers are assertions in the repository, not measurements from a test run for this article. No lab tests were executed during preparation.

### Why the constraint does not reject it

Alice still has 200,000. There is no negative balance to reject. The missing debit is a lost update, and the violated invariant is the cross-account sum. A successful commit and a satisfied CHECK constraint are insufficient evidence that the transfer was correct.

This does not mean every READ COMMITTED update loses data. The vulnerable implementation reads, calculates in application memory, and later writes an absolute value. Its query pattern is part of the explanation.

## Separate the SQL standard from PostgreSQL behavior

The README describes the standard anomaly ladder: READ UNCOMMITTED permits dirty reads; READ COMMITTED prevents dirty reads but permits non-repeatable and phantom reads; standard REPEATABLE READ prevents dirty and non-repeatable reads but may permit phantoms; SERIALIZABLE prevents those anomalies.

PostgreSQL provides stronger behavior in two important places. READ UNCOMMITTED behaves as READ COMMITTED and does not permit dirty reads. PostgreSQL REPEATABLE READ uses a stable snapshot and prevents the classic phantom read demonstrated by the lab.

| PostgreSQL level | Visibility in this lab | Remaining concern |
| --- | --- | --- |
| READ UNCOMMITTED | Same behavior as READ COMMITTED | Does not grant a stable transaction snapshot |
| READ COMMITTED | New snapshot per statement | Repeated reads may change; naive writes may lose updates |
| REPEATABLE READ | Snapshot acquired at the first read | Concurrent same-row update can fail with 40001 |
| SERIALIZABLE | SSI with serial execution equivalence | Conflicts may abort a transaction; retry must be handled |

The table describes the engine used by the repository. It should not be treated as a portable implementation table for other databases.

## Read Committed: two statements, two views

`TestReadCommitted_NonRepeatableRead` coordinates a reader and a writer. The reader first sees Alice at 1,000,000. The writer adds 500,000 and commits. The reader's second SELECT, still inside the same transaction, is expected to see 1,500,000.

The change is committed, so this is not a dirty read. It is a non-repeatable read caused by statement snapshots.

The phantom experiment uses a predicate instead of a single account. Two invoices initially have status PAID. A writer inserts a third PAID invoice and commits between the reader's queries:

```sql
SELECT COUNT(*) FROM isolation_invoices WHERE status = 'PAID'
```

Under READ COMMITTED, the test expects counts of two and three. Under PostgreSQL REPEATABLE READ, it expects two and two. The predicate is unchanged; the visible set of matching rows differs only in the first strategy.

## A stable snapshot is not exclusive access

PostgreSQL MVCC lets ordinary readers observe versions without taking a FOR UPDATE row lock. `TestRepeatableRead_SnapshotIsolation` expects both reads to return the original 1,000,000 even after another transaction commits an increase.

The concurrent-update experiment is different. Two REPEATABLE READ transactions read the same account, then try to deduct 800,000. Its assertions require exactly one success and one serialization failure, SQLSTATE 40001.

The application must therefore handle an aborted transaction. Choosing REPEATABLE READ does not make the second transfer silently succeed correctly. It changes the failure behavior.

The README also discusses write skew as a possible serialization anomaly under snapshot isolation. The supplied experiments demonstrate same-row update conflicts; they do not reproduce a separate write-skew scenario. Keep the conceptual warning separate from the experiment's actual evidence.

## Row locking moves validation behind coordination

`TransferWithLock` stays at READ COMMITTED but obtains both account locks before validating the sender. Here is the locking excerpt:

```go
firstID, secondID := DeterministicLockOrder(fromID, toID)

	// Lock both accounts deterministically
	var b1, b2 int64
	err = tx.QueryRowContext(ctx, "SELECT balance FROM isolation_accounts WHERE id = $1 FOR UPDATE", firstID).Scan(&b1)
	if err != nil {
		return fmt.Errorf("lock account %d: %w", firstID, err)
	}
	err = tx.QueryRowContext(ctx, "SELECT balance FROM isolation_accounts WHERE id = $1 FOR UPDATE", secondID).Scan(&b2)
	if err != nil {
		return fmt.Errorf("lock account %d: %w", secondID, err)
	}
```

*Excerpt from `transfer.go`; the rest of the transfer is omitted.*

A competing FOR UPDATE or UPDATE on those rows waits until the holder ends its transaction. An ordinary SELECT continues to use MVCC. The lock is coordination for the write path, not a promise that every reader is blocked.

After locking, the implementation identifies which locked balance belongs to the sender and checks sufficient funds. Debit and credit use arithmetic SQL updates, and the audit record is written before commit. The second competing transfer validates after obtaining access to the row, rather than trusting an earlier unlocked balance.

### Acquire locks in a shared order

If A-to-B locks A first while B-to-A locks B first, the two transactions can create circular wait. The lab uses a helper independent of transfer direction:

```go
func DeterministicLockOrder(id1, id2 int) (firstID, secondID int) {
	if id1 < id2 {
		return id1, id2
	}
	return id2, id1
}
```

Both directions acquire the smaller account ID, then the larger. The bidirectional test runs two workers with 50 transfers each. It checks errors and total balance. This order addresses the two-account lock pattern used here; it does not guarantee that arbitrary additional locking elsewhere can never deadlock.

Lock waits are a real cost. The README's phrase about deterministic latency should not be read as a guarantee from this implementation. Contention and transaction duration still affect waiting, and the repository provides no latency benchmark.

## Serializable means serial equivalence, not a global queue

PostgreSQL SERIALIZABLE uses Serializable Snapshot Isolation. Concurrent transactions can run together; committed results must be equivalent to an allowed serial execution. Conflicts can require an abort.

`TestSerializable_ConcurrentUpdate_SerializationFailure` coordinates two same-row updates and expects one success plus one 40001. It is a conflict experiment, not a throughput comparison or proof of every possible invariant.

The recovery path needs a new transaction with fresh reads and validation. Retrying only the failed statement would not repeat the decision that caused it.

## Retry the whole operation, with a limit

The wrapper passes a complete transfer as the operation:

```go
func (r *PostgresWalletRepo) TransferSerializableWithRetryPolicy(ctx context.Context, fromID, toID int, amount int64, maxAttempts int, policy RetryPolicy) error {
	repo := r
	return RetryTransaction(ctx, maxAttempts, policy, func(ctx context.Context) error {
		return repo.TransferSerializable(ctx, fromID, toID, amount)
	})
}
```

*Complete wrapper from `transfer.go`.*

Each attempt calls `TransferSerializable` again, including BeginTx, read, validation, debit, credit, audit, and commit. Failed attempts defer rollback. The retry helper accepts only serialization failure 40001 or deadlock 40P01, through the PostgreSQL error or corresponding sentinel.

```go
func IsRetryableTxError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrSerializationFailure) || errors.Is(err, ErrDeadlockDetected) {
		return true
	}
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "40001" || pqErr.Code == "40P01"
	}
	return false
}
```

*Complete filter from `transfer.go`.*

`WrapTxError` preserves both sentinel matching with `errors.Is` and access to `*pq.Error` with `errors.As`. Constraint violations such as 23505 and 23503, generic errors, and insufficient funds are not automatically retried.

`maxAttempts` includes the initial attempt. Non-positive values are rejected. The default base delay is 10 milliseconds and the maximum is one second. Before another attempt, the helper computes an exponential upper bound, caps it, and multiplies it by a random value for full jitter.

Context cancellation is checked before attempts and can interrupt the default sleep. When retryable errors exhaust the limit, `ErrMaxRetryExceeded` wraps the last error. `retry_test.go` covers filtering, error preservation, success after a conflict, attempt exhaustion, cancellation, invalid limits, and delay behavior.

## Test the invariant under contention

The stress fixture differs from the introductory wallet: three accounts begin at 50,000 each, giving total 150,000. One hundred goroutines attempt transfers of 1,000 using `TransferWithLock` after a shared start gate is released.

The test checks transfer errors, non-negative final balances, and the conserved sum of 150,000. A start gate creates competing work; it does not force every goroutine to reach every query at an identical moment. The naive test uses a more specific barrier after reads to reproduce its anomaly.

To run the supplied experiments, start the repository's PostgreSQL infrastructure and execute:

```bash
cd labs/08-database-isolation-level
go test -v ./...
```

The database helper skips integration tests when PostgreSQL is unavailable unless `REQUIRE_POSTGRES=1` makes absence a failure. A test command that reports success alongside skipped experiments is not evidence that database behavior was verified.

## Choose from the invariant outward

For this wallet's known rows, READ COMMITTED with FOR UPDATE and a shared lock order coordinates the critical read-check-write path. REPEATABLE READ serves a stable multi-query view but requires conflict handling for updates. SERIALIZABLE addresses serial equivalence at the cost of aborts, retry work, and dependency tracking.

The README mentions an atomic conditional update as an alternative, but the lab does not implement that strategy. It also offers no measured throughput or latency ranking. The practical conclusion is to specify the invariant, identify the relevant reads and writes, choose coordination or conflict detection, and verify the resulting state under overlap.

## Repository and source map

- [Lab 08](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/08-database-isolation-level)
- [README.md](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/08-database-isolation-level/README.md)
- [schema.sql](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/08-database-isolation-level/schema.sql)
- [isolation.go](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/08-database-isolation-level/isolation.go)
- [transfer.go](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/08-database-isolation-level/transfer.go)
- [read_committed_test.go](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/08-database-isolation-level/read_committed_test.go)
- [repeatable_read_test.go](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/08-database-isolation-level/repeatable_read_test.go)
- [phantom_read_test.go](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/08-database-isolation-level/phantom_read_test.go)
- [transfer_test.go](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/08-database-isolation-level/transfer_test.go)
- [serializable_test.go](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/08-database-isolation-level/serializable_test.go)
- [retry_test.go](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/08-database-isolation-level/retry_test.go)
