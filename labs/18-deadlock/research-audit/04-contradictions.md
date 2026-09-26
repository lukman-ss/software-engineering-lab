# Contradiction Audit

## Contradiction 1 (False Contradiction / Expected Variation)

Statement A:
Victim selection is unpredictable and should not be relied upon.
Location: PostgreSQL Docs §13.3.4 (via `05-report.md`)

Statement B:
Victim selection attempts to roll back the smallest transaction based on the number of rows inserted/updated/deleted.
Location: MySQL 8.4 Reference Manual §17.7.5.2 (via `05-report.md`)

Statement C:
Victim selection chooses the transaction that is least expensive to roll back, subject to configurable deadlock priority.
Location: Microsoft Learn Deadlocks Guide (via `05-report.md`)

Type:
CODE_DOC_MISMATCH / IMPLEMENTATION_VARIATION

Impact:
Applications cannot rely on a universal "victim selection" behavior if they intend to be database-agnostic.

Assessment:
Not a logical contradiction. Correctly classified by the research as an "Implementation Variation" between DB engines.

---

## Contradiction 2 (False Contradiction / Expected Variation)

Statement A:
Locking waits indefinitely unless a deadlock is detected. `deadlock_timeout` (default 1s) is only the delay before checking for deadlocks.
Location: PostgreSQL Docs §19.12 (via `05-report.md`)

Statement B:
Locking falls back to `innodb_lock_wait_timeout` (default 50s) if a lock cannot be acquired, irrespective of circular waits (or when `innodb_deadlock_detect` is disabled).
Location: MySQL 8.4 Reference Manual §17.7.5.2 (via `05-report.md`)

Type:
IMPLEMENTATION_VARIATION

Impact:
A PostgreSQL deadlock might cause applications to hang indefinitely if the conflict is not a circular deadlock but a pure contention timeout. MySQL bounds this automatically.

Assessment:
Not a logical contradiction. Correctly identified and clarified in `05-report.md` Finding 8.

---
## Summary

No material contradictions found.
The research accurately distinguishes universal database theory (e.g. Coffman conditions, wait-for cycles) from vendor-specific implementation details (e.g. error codes, timeout semantics, and victim selection logic).
