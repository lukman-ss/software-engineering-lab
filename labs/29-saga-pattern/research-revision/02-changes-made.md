# Changes Made

## Revision 1

Audit Issue:
LOW — Dead link in Source Citations (Outdated Source)

Location:
`labs/29-saga-pattern/research/02-sources.md` (Source 5)
`labs/29-saga-pattern/research/runs/2026-09-28-saga-pattern/02-sources.md` (Source 5)

Problem:
`https://learn.microsoft.com/en-us/dotnet/architecture/cloud-native/saga-pattern` returns HTTP 404.

Required Revision:
Replace with active .NET Saga documentation or remove Source 5 since Sources 1 and 2 already fully support all core claims.

Action:
- Removed Source 5 from `02-sources.md` (main and run copy) because Sources 1 (Azure Architecture Center) and 2 (microservices.io) provide authoritative, comprehensive coverage of the Saga pattern, including all core claims verified in the claim audit.

Verification:
- Source 5 URL checked and confirmed 404 Not Found
- Sources 1-4 remain verified and reachable per source audit
- Core claims in claim audit (03-claim-audit.md) remain supported by Sources 1-4

Status:
RESOLVED

## Revision 2

Audit Issue:
MEDIUM — Recovery Mechanism for Failed Compensations (Missing Case)

Location:
`labs/29-saga-pattern/research/03-evidence.md` (Evidence 9), `06-open-questions.md` (Unanswered Question 4)
`labs/29-saga-pattern/research/runs/2026-09-28-saga-pattern/03-evidence.md` (Evidence 9), `06-open-questions.md` (Unanswered Question 4)

Problem:
The research notes that compensating transactions can fail, but does not provide architectural patterns for handling permanently failed compensations (e.g., dead-letter queue processing, human-in-the-loop manual reconciliation consoles, out-of-band balance adjustments).

Required Revision:
Detail operational remediation workflows when compensating actions fail in production.

Action:
- Added new Evidence 9a in `03-evidence.md` (main and run copy) documenting operational remediation patterns: dead-letter queues, alerting/monitoring, manual reconciliation via admin consoles, out-of-band adjustments, citing Microsoft Azure Architecture Center statements about compensating transaction limitations and need for monitoring/tracking.
- Updated `06-open-questions.md` Question 4 (main and run copy) to expand coverage from LOW to EXPANDED, adding concrete operational patterns: DLQ for failed compensations, alert escalation to on-call teams, manual reconciliation via admin console, out-of-band data adjustments, referencing Azure documentation.

Verification:
- Evidence 9a quotes verified against actual Microsoft Azure Saga pattern page
- Cited passages are verbatim from the source: "Limitations of compensating transactions: Compensating transactions might not always succeed, which can leave the system in an inconsistent state." and "Need for monitoring and tracking sagas: Monitoring and tracking the workflow of a saga are essential tasks to maintain operational oversight."
- Open question coverage updated to reflect added detail

Status:
RESOLVED