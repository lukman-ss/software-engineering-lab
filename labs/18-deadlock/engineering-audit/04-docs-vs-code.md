# Docs vs Code Audit

Target Lab: labs/18-deadlock

## Mismatch Analysis

### README vs Code
- Description of `cmd/demo/main.go`: matches.
- Description of `internal/bank/account.go`: matches.
- Description of `internal/transfer/transfer.go`: matches.
- Commands provided (`go test`, `go test -race`, `go run`): exact match, working as documented.

### Research Claims vs Code
- **Finding 1 (Circular Wait)**: Accurately modeled via concurrent naive transfers blocking each other.
- **Finding 2 (Deadlock Victim)**: Accurately modeled via context timeout aborting one transaction and returning `ErrDeadlock`.
- **Finding 3 (Lock Ordering)**: Accurately modeled via lexicographical sorting in `TransferOrdered`.
- **Finding 4 (Transaction Duration)**: Accurately modeled and empirically proven in `TestTransactionDurationImpact` using simulated delays.
- **Finding 5 (Retry Recovery)**: Accurately modeled in `TransferWithRetry`.

### Engineering Notes vs Code
- The design document `01-design.md` states "Tests verify transaction delay increases deadlock frequency." This is explicitly present in `TestTransactionDurationImpact`.
- The implementation notes state "Retry logic uses simple loop with basic sleep instead of full exponential backoff with jitter". This perfectly aligns with `TransferWithRetry` implementation.
- Phrasing "boolean ID string comparison" in `02-implementation-notes.md` is slightly awkward but functionally corresponds to the string inequality check (`acc1.ID > acc2.ID`).

## Conclusion
PASS. No DOC_CODE_MISMATCH. Implementation faithfully represents all approved research claims without overstating capabilities.
