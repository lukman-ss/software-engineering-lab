# Docs vs Code Audit

Target Lab: labs/18-deadlock

## Comparison Matrix

| Component / Claim | Documented Behavior | Code / Test Implementation | Alignment |
|---|---|---|---|
| Naive Transfer Deadlock | Described in README & Design | `transfer.TransferNaive` creates circular wait; `TestDeadlockOccurrence` validates failure | PASS |
| Lock Ordering Prevention | Described in README & Design | `transfer.TransferOrdered` sorts IDs before locking; `TestLockOrderingPreventsDeadlock` passes | PASS |
| Application Retry Recovery | Described in README & Design | `transfer.TransferWithRetry` retries upon `ErrDeadlock`; `TestRetryRecoversDeadlock` passes | PASS |
| Transaction Duration Impact | Claimed in Research & Design | `TestTransactionDurationImpact` compares 5ms delay vs 0 delay across 10 iterations | PASS |
| CLI Commands | Documented in `README.md` (`go test ./...`, `go test -race ./...`, `go run ./cmd/demo`) | All commands execute cleanly and match output documented in `03-execution-result.md` | PASS |

## Identified Mismatches
None detected.

- DOC_CODE_MISMATCH: None.
- TEST_CLAIM_MISMATCH: None.
- RESEARCH_IMPLEMENTATION_MISMATCH: None.
