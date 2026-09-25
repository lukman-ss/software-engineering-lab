# Source Map

## Mental Model & Core Concept (Circular Wait & Deadlock Victim)

Research:
- `research/runs/2026-09-25-deadlock/05-report.md` (Finding 1, Finding 2)
- `research/runs/2026-09-25-deadlock/03-evidence.md` (Evidence 1, Evidence 2)
- `research/runs/2026-09-25-deadlock/02-sources.md` (Source 1, Source 3)

Implementation:
- `internal/bank/account.go` (Timeout logic memodelkan Deadlock Victim)

Tests:
- `tests/transfer_test.go` (`TestDeadlockOccurrence`)

## How It Works & Lock Ordering Prevention

Research:
- `research/runs/2026-09-25-deadlock/05-report.md` (Finding 3)
- `research/runs/2026-09-25-deadlock/03-evidence.md` (Evidence 3)
- `research/runs/2026-09-25-deadlock/02-sources.md` (Source 1, Source 3)

Implementation:
- `internal/transfer/transfer.go` (`TransferOrdered`)

Tests:
- `tests/transfer_test.go` (`TestLockOrderingPreventsDeadlock`)

## Recovery / Application-Level Retry

Research:
- `research/runs/2026-09-25-deadlock/05-report.md` (Finding 5)
- `research/runs/2026-09-25-deadlock/03-evidence.md` (Evidence 4)
- `research/runs/2026-09-25-deadlock/02-sources.md` (Source 1, Source 3)

Implementation:
- `internal/transfer/transfer.go` (`TransferWithRetry`)

Tests:
- `tests/transfer_test.go` (`TestRetryRecoversDeadlock`)

## Production Considerations (Transaction Duration & Timeout)

Research:
- `research/runs/2026-09-25-deadlock/05-report.md` (Finding 4)
- `research/runs/2026-09-25-deadlock/03-evidence.md` (Evidence 3, Evidence 5)
- `research/runs/2026-09-25-deadlock/02-sources.md` (Source 2, Source 3)

Implementation:
- `internal/transfer/transfer.go` (`delay` parameter)

Tests:
- `tests/transfer_test.go` (`TestTransactionDurationImpact`)

## Common Mistakes (Self-Transfer Caveat)

Engineering Audit:
- `engineering-audit/05-gaps.md` (MISSING_EDGE_CASE: self-transfer missing check `from == to`)
