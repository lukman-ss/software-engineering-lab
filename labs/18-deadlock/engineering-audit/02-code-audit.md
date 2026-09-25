# Code Audit

Target Lab: labs/18-deadlock

## Finding 1

Location: internal/bank/account.go:21-28
Claimed Behavior: Deadlock victim aborted via lock timeout without unrecoverable crash.
Observed Implementation: Account lock uses buffered channel (cap 1) with `select` against `ctx.Done()`, returning `ErrDeadlock`.
Assessment: PASS
Severity: LOW
Notes: Accurately simulates RDBMS deadlock abort mechanism cleanly in Go.

## Finding 2

Location: internal/transfer/transfer.go:10-26
Claimed Behavior: Unordered transfers trigger circular wait under concurrent bidirectional execution.
Observed Implementation: `TransferNaive` locks `from` before `to` with simulated work delay.
Assessment: PASS
Severity: LOW
Notes: Creates circular wait condition reliably when invoked concurrently with swapped accounts.

## Finding 3

Location: internal/transfer/transfer.go:29-49
Claimed Behavior: Lock ordering prevents deadlock completely.
Observed Implementation: `TransferOrdered` sorts lock acquisition by lexicographical comparison of account IDs.
Assessment: PASS
Severity: LOW
Notes: Completely breaks circular wait condition.

## Finding 4

Location: internal/transfer/transfer.go:52-69
Claimed Behavior: Application retry loop recovers aborted transactions.
Observed Implementation: `TransferWithRetry` runs bounded loop with 10ms child timeout and 2ms backoff upon encountering `ErrDeadlock`.
Assessment: PASS
Severity: LOW
Notes: Correctly handles retry semantics as specified in research finding 5.

## Finding 5

Location: internal/transfer/transfer.go:10-49
Claimed Behavior: Accounts transfer money safely.
Observed Implementation: Self-transfer (`from == to`) not validated; results in self-deadlock on second lock acquisition.
Assessment: WARNING
Severity: LOW
Notes: Minor missing validation on identical accounts. Out of scope for lab core claims.
