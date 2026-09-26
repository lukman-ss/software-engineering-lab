# Code Audit

Target Lab: labs/18-deadlock

## Finding 1

Location: internal/bank/account.go:22-29 (Lock method)
Claimed Behavior: Lock returns ErrDeadlock when context expires, simulating deadlock victim.
Observed Implementation: Lock uses select on a.ch (token channel) and ctx.Done(). If ctx.Done(), returns ErrDeadlock.
Assessment: PASS
Severity: -
Notes: Correctly models timeout-based deadlock detection.

## Finding 2

Location: internal/bank/account.go:31-33 (Unlock method)
Claimed Behavior: Unlock returns the token to allow other goroutines to acquire the lock.
Observed Implementation: Unlock sends a struct{} into the buffered channel ch.
Assessment: PASS
Severity: -
Notes: Matches lock/unlock semantics. Channel size 1 ensures mutual exclusion.

## Finding 3

Location: internal/transfer/transfer.go:10-26 (TransferNaive)
Claimed Behavior: Naive transfer locks in caller-specified order (from then to), risking deadlock.
Observed Implementation: Locks from.Account then to.Account with defer unlocks. Includes delay between locks.
Assessment: PASS
Severity: -
Notes: Precisely implements naive locking. Delay increases deadlock window.

## Finding 4

Location: internal/transfer/transfer.go:28-50 (TransferOrdered)
Claimed Behavior: Ordered transfer locks accounts by ID order to prevent circular wait.
Observed Implementation: Sorts accounts by ID, locks first then second. Updates original account balances.
Assessment: PASS
Severity: -
Notes: Correct ordering prevents deadlock. Balance updates on locked accounts are safe.

## Finding 5

Location: internal/transfer/transfer.go:52-71 (TransferWithRetry)
Claimed Behavior: Retry wrapper handles deadlocks by retrying naive transfer with backoff.
Observed Implementation: Loop with short timeout attempts, checks for ErrDeadlock or timeout, sleeps between retries.
Assessment: PASS
Severity: -
Notes: Retry logic recovers from deadlock victims. Potential improvement: context expiration could cause misleading ErrDeadlock return after maxRetries.

## Finding 6

Location: internal/transfer/transfer_test.go
Claimed Behavior: Tests cover deadlock occurrence, prevention, recovery, and duration impact.
Observed Implementation: Four concurrency tests validate all claims using synchronized goroutines and context timeouts.
Assessment: PASS
Severity: -
Notes: Tests are deterministic and pass consistently. TestTransactionDurationImpact flaky? Ran 20 times, all passed.
