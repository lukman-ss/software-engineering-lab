---
title: "Deadlock in Concurrent Systems: Mechanics, Lock Ordering, and Application Retries"
published: false
description: "Simulating circular wait deadlocks in Go, preventing deadly embrace via lock ordering, and recovering aborted transactions with application retries."
tags: "go, architecture, concurrency, database"
---

When concurrent processes mutate multiple shared resources without coordination, they risk stepping into a mutual circular wait. Neither transaction can proceed because each holds a lock that the other requires. Left unresolved, this condition—known as a **deadlock** or *deadly embrace*—freezes transactional throughput and exhausts connection pools.

In production databases, engine-level deadlock monitors resolve this dilemma forcefully by terminating one of the transactions as a *deadlock victim*. However, treating deadlocks as purely database errors overlooks application-level design flaws.

In **Software Engineering Lab 18**, we simulate deadlock mechanics in Go using buffered channels and context deadlines. We examine how uncoordinated acquisition causes circular wait, how deterministic **Lock Ordering** prevents cycles altogether, and how **Application Retries** recover aborted transactions.

```
Transaction 1 holds Lock A ---> Transaction 1 waits for Lock B
                ^                                      |
                |                                      v
Transaction 2 waits for Lock A <--- Transaction 2 holds Lock B
```

---

## The Root Cause: Circular Wait and Deadlock Victims

Four conditions (the Coffman conditions) must hold simultaneously for a deadlock to exist: mutual exclusion, hold and wait, no preemption, and **circular wait**.

When two concurrent transfers execute in reverse directions:
- **Transfer 1** locks `Account A` and prepares to lock `Account B`.
- **Transfer 2** locks `Account B` and prepares to lock `Account A`.

Because neither transaction yields its held lock, neither can acquire its second lock. Relational database engines (such as PostgreSQL, MySQL InnoDB, or Microsoft SQL Server) maintain a background *deadlock monitor* that checks a **Wait-For Graph (WFG)** periodically. When a directed cycle appears, the database engine aborts one transaction (the victim), rolling back its state and returning an error (e.g., PostgreSQL `40P01` or SQL Server 1205).

---

## 1. Simulating Naive Locking (Deadlock Occurrence)

In `internal/bank/account.go`, we model each account using a 1-slot buffered channel acting as a mutex. The `Lock` method honors a `context.Context` deadline:

```go
package bank

import (
	"context"
	"errors"
)

var ErrDeadlock = errors.New("deadlock victim")

type Account struct {
	ID      string
	Balance int
	ch      chan struct{}
}

func NewAccount(id string, bal int) *Account {
	ch := make(chan struct{}, 1)
	ch <- struct{}{}
	return &Account{ID: id, Balance: bal, ch: ch}
}

func (a *Account) Lock(ctx context.Context) error {
	select {
	case <-a.ch:
		return nil
	case <-ctx.Done():
		return ErrDeadlock
	}
}

func (a *Account) Unlock() {
	a.ch <- struct{}{}
}
```

When two transfers execute naively, they lock resources in the sequence provided by the caller:

```go
// TransferNaive creates a potential deadlock by locking from then to.
func TransferNaive(ctx context.Context, from, to *bank.Account, amount int, delay time.Duration) error {
	if err := from.Lock(ctx); err != nil {
		return err
	}
	defer from.Unlock()

	time.Sleep(delay) // simulate transaction duration

	if err := to.Lock(ctx); err != nil {
		return err
	}
	defer to.Unlock()

	from.Balance -= amount
	to.Balance += amount
	return nil
}
```

Running two bidirectional transfers concurrently with `delay > 0` produces a circular wait. Once `ctx` expires, one or both transactions return `bank.ErrDeadlock`.

---

## 2. Prevention: Consistent Lock Ordering

Deadlocks cannot occur if transactions acquire shared locks in a globally consistent order. Sorting resources by a deterministic key (such as account ID) breaks the circular wait condition.

```go
// TransferOrdered prevents deadlock by always locking accounts in alphabetical order.
func TransferOrdered(ctx context.Context, acc1, acc2 *bank.Account, amount int, delay time.Duration) error {
	first, second := acc1, acc2
	if acc1.ID > acc2.ID {
		first, second = acc2, acc1
	}

	if err := first.Lock(ctx); err != nil {
		return err
	}
	defer first.Unlock()

	time.Sleep(delay)

	if err := second.Lock(ctx); err != nil {
		return err
	}
	defer second.Unlock()

	acc1.Balance -= amount
	acc2.Balance += amount
	return nil
}
```

Regardless of transfer direction, both transactions compete for the lower ID first (`acc1.ID < acc2.ID`). The winner proceeds while the second transaction waits in line before acquiring any lock. No circular dependency can form.

---

## 3. Recovery: Application-Level Retries

In complex systems with dynamic queries, global lock ordering is not always feasible. Applications must be prepared to catch deadlock errors and retry the failed transaction transparently:

```go
// TransferWithRetry handles deadlocks by retrying the operation.
func TransferWithRetry(ctx context.Context, from, to *bank.Account, amount int, delay time.Duration, maxRetries int) error {
	for i := 0; i < maxRetries; i++ {
		attemptCtx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
		err := TransferNaive(attemptCtx, from, to, amount, delay)
		cancel()

		if err == nil {
			return nil
		}
		if err != bank.ErrDeadlock && err != context.DeadlineExceeded {
			return err
		}

		time.Sleep(2 * time.Millisecond)
	}
	return bank.ErrDeadlock
}
```

When a transaction is aborted as a deadlock victim, a brief pause followed by a clean retry allows the remaining transaction to finish, after which the retried transaction succeeds.

---

## Verification and Test Proofs

Our test suite in `tests/transfer_test.go` verifies all four operational behaviors:

1. **Deadlock Occurrence (`TestDeadlockOccurrence`)**: Concurrently executing bidirectional `TransferNaive` calls consistently triggers `ErrDeadlock` on at least one transfer.
2. **Lock Ordering Prevention (`TestLockOrderingPreventsDeadlock`)**: Bidirectional concurrent transfers using `TransferOrdered` complete with zero errors and perfect account balance consistency.
3. **Application Retry Recovery (`TestRetryRecoversDeadlock`)**: Wrapping naive transfers in `TransferWithRetry` yields 100% eventual success across both concurrent transactions.
4. **Transaction Duration Impact (`TestTransactionDurationImpact`)**: Demonstrates that holding locks longer (5ms vs 0ms) increases deadlock frequency due to wider contention windows.

```bash
$ go test -v -race ./...
=== RUN   TestDeadlockOccurrence
--- PASS: TestDeadlockOccurrence (0.02s)
=== RUN   TestLockOrderingPreventsDeadlock
--- PASS: TestLockOrderingPreventsDeadlock (0.01s)
=== RUN   TestRetryRecoversDeadlock
--- PASS: TestRetryRecoversDeadlock (0.02s)
=== RUN   TestTransactionDurationImpact
--- PASS: TestTransactionDurationImpact (0.21s)
PASS
ok      labs/18-deadlock/tests  1.730s
```

---

## Architectural Guardrails for Production

- **Universal Lock Ordering**: Apply sorting across all transaction paths. Sorting IDs in one service while executing unsorted updates in another will still trigger deadlocks.
- **Shorten Transaction Lifespans**: Keep transactional scopes tight. Never perform external network calls (such as payment gateway webhooks or email dispatch) while holding database row locks.
- **Batch Update Sorting**: When performing bulk mutations (`UPDATE ... WHERE id IN (...)`), sort the primary keys ascending (`ORDER BY id ASC`) before locking.
- **Idempotent Retries**: Ensure that any transaction submitted for automatic retry does not produce non-idempotent side effects (e.g., duplicate notifications or token generation).

---

## Implementation Boundaries and Trade-Offs

- **Simulation Scope**: This lab uses an in-memory Go channel simulation with context timeouts. It models the operational effect of deadlock victims rather than executing a true graph-cycle traversal (Wait-For Graph) inside a storage engine.
- **Fixed Backoff Simplification**: `TransferWithRetry` employs a fixed 2ms backoff for fast, deterministic unit test verification. Production environments should implement exponential backoff with full jitter to avoid retry storms.

---

## Repository

Complete source code, benchmarks, and test assertions are available at:
[https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/18-deadlock](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/18-deadlock)
