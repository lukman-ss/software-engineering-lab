# Source Map

## Mental Model & Core Concept (Circular Wait & Deadlock Victim)

Research:
research/05-report.md (Finding 1, Finding 2, Finding 3)
research/03-evidence.md (Evidence 1, Evidence 2, Evidence 3, Evidence 4)

Implementation:
internal/bank/account.go (Timeout-based locking models deadlock victim abort)

Tests:
tests/transfer_test.go (TestDeadlockOccurrence)

## How It Works & Lock Ordering Prevention

Research:
research/05-report.md (Finding 5)
research/03-evidence.md (Evidence 5)

Implementation:
internal/transfer/transfer.go (TransferOrdered)

Tests:
tests/transfer_test.go (TestLockOrderingPreventsDeadlock)

## Recovery / Application-Level Retry

Research:
research/05-report.md (Finding 7)
research/03-evidence.md (Evidence 7, Evidence 8)

Implementation:
internal/transfer/transfer.go (TransferWithRetry)

Tests:
tests/transfer_test.go (TestRetryRecoversDeadlock)

## Production Considerations (Transaction Duration & Timeout)

Research:
research/05-report.md (Finding 4, Finding 6, Finding 8, Finding 9)
research/03-evidence.md (Evidence 6, Evidence 9, Evidence 10, Evidence 11, Evidence 12, Evidence 14)

Implementation:
internal/transfer/transfer.go (delay parameter, TransferNaive)

Tests:
tests/transfer_test.go (TestTransactionDurationImpact)

## Source Gaps / Caveats (Preserved from Audits)

Research Audit non-blocking issues:
- Coffman Conditions / Two-Phase Locking sourced via Wikipedia (Tier 3) — primary literature (Coffman 1971, Bernstein 1987) not directly accessed.
- Oracle Database omitted due to URL retrieval failure.
- "Sistem PPOB" context is illustrative only, unbacked by primary domain literature.
- MySQL documentation accessed via Wayback Machine snapshots (archived 2024-2025), not live URLs.

Engineering Audit non-blocking issue:
- TransferWithRetry uses fixed 2ms backoff (not exponential backoff with jitter) — documented deliberate simplification for readability.
- Lab is an in-memory simulation; does not implement actual Wait-For Graph (WFG) cycle detection; models victim abort via context timeout.