# Claim Audit

## Claim 1

Claim:
Deadlock occurs if and only if all four Coffman conditions (mutual exclusion, hold-and-wait, no preemption, circular wait) are simultaneously satisfied.

Location:
`05-report.md` §Finding 1, `03-evidence.md` §Evidence 1

Evidence Provided:
Quotation of 4 Coffman conditions and demonstration via two-transaction UPDATE scenario.

Source:
Coffman et al. 1971 / Silberschatz via Wikipedia (`https://en.wikipedia.org/wiki/Deadlock_(computer_science)`)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurate theoretical foundation. Minor nuance: Coffman conditions are necessary conditions; circular wait in a resource allocation graph with multi-unit resources is necessary but not always sufficient unless resources are single-unit. In standard database lock managers with exclusive locks, they are necessary and sufficient.

---

## Claim 2

Claim:
Two individually correct, valid transactions will deadlock if they acquire the same set of locks in opposing order.

Location:
`05-report.md` §Finding 2, `03-evidence.md` §Evidence 4

Evidence Provided:
Concrete examples from PostgreSQL docs (accounts 11111 vs 22222) and MySQL docs (Animals vs Birds).

Source:
PostgreSQL Docs §13.3.4, MySQL 8.4 Manual §17.7.5.1, Microsoft Learn Deadlocks Guide

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Extremely well-supported across all three major RDBMS vendors.

---

## Claim 3

Claim:
PostgreSQL, MySQL/InnoDB, and SQL Server detect deadlocks automatically and abort/rollback one transaction (deadlock victim), but victim selection logic varies and should not be relied upon by application code.

Location:
`05-report.md` §Finding 3, `03-evidence.md` §Evidence 2, 3

Evidence Provided:
PostgreSQL docs state victim selection is unpredictable. MySQL documentation describes picking the smallest transaction (by rows modified). SQL Server evaluates rollback cost and priority.

Source:
PostgreSQL Docs §13.3.4, MySQL Manual §17.7.5.2, Microsoft Learn Deadlocks Guide

Source Actually Supports Claim:
YES

Classification:
FACT / IMPLEMENTATION-SPECIFIC

Severity:
LOW

Notes:
Distinction between generic circular wait breaking and vendor-specific victim choice is cleanly documented.

---

## Claim 4

Claim:
Longer transaction durations increase deadlock probability because locks are held longer, expanding the window of contention. Non-database operations (APIs, file writes, user waits) inside transactions aggravate this risk.

Location:
`05-report.md` §Finding 4, `03-evidence.md` §Evidence 6

Evidence Provided:
Direct citations from PostgreSQL, MySQL, and Microsoft Learn cautioning against long transactions and external waits.

Source:
PostgreSQL §13.3.4, MySQL §17.7.5.3, Microsoft Learn Deadlocks Guide, Percona Blog

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Universal consensus across all vendors and operational engineering guides.

---

## Claim 5

Claim:
Consistent lock ordering across all concurrent application paths is the primary and most effective architectural defense against deadlocks.

Location:
`05-report.md` §Finding 5, `03-evidence.md` §Evidence 5

Evidence Provided:
PostgreSQL notes identical ordering eliminates the cycle; MySQL notes transactions form well-defined queues; SQL Server and Percona recommend serializing access order.

Source:
PostgreSQL §13.3.4, Microsoft Learn Deadlocks Guide, MySQL §17.7.5.3

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Mathematically breaks the circular wait condition. Proven industry best practice.

---

## Claim 6

Claim:
Lock contention can be reduced by appropriate indexing, selecting lower isolation levels (READ COMMITTED), and utilizing row versioning (MVCC / SNAPSHOT), thereby minimizing gap locks and lock footprints.

Location:
`05-report.md` §Finding 6, `03-evidence.md` §Evidence 10

Evidence Provided:
MySQL/Percona guidance on eliminating gap locking via READ COMMITTED; SQL Server documentation on READ_COMMITTED_SNAPSHOT; PostgreSQL SSI predicate locking behavior.

Source:
Microsoft Learn Deadlocks Guide, MySQL §17.7.5.3, Percona Blog, PostgreSQL §13.2.3

Source Actually Supports Claim:
YES

Classification:
FACT / IMPLEMENTATION-SPECIFIC

Severity:
MEDIUM

Notes:
Trade-offs of lower isolation levels (e.g., non-repeatable reads, phantom reads) must be understood by implementers, but the contention-reduction mechanism is accurately stated.

---

## Claim 7

Claim:
Deadlock errors must be caught by the application and retried with exponential backoff and randomized jitter, but retries are safe only if the transaction/operation is strictly idempotent.

Location:
`05-report.md` §Finding 7, `03-evidence.md` §Evidence 7, 8

Evidence Provided:
MySQL and Microsoft documentation specifying automatic query resubmission with randomized delays; AWS prescriptive guidance on exponential backoff and idempotency requirements.

Source:
MySQL §17.7.5.3, Microsoft Learn Deadlocks Guide, Percona Blog, AWS Architecture Guidance

Source Actually Supports Claim:
YES

Classification:
INTERPRETATION / BEST_PRACTICE

Severity:
MEDIUM

Notes:
Idempotency requirement is critical in financial/PPOB systems to prevent double processing when retrying after a partial failure.

---

## Claim 8

Claim:
Deadlock is fundamentally distinct from lock timeout: deadlock is an active cycle of mutually blocked transactions requiring termination, whereas lock timeout is waiting beyond a configured duration without necessarily having a circular dependency.

Location:
`05-report.md` §Finding 8, `03-evidence.md` §Evidence 12

Evidence Provided:
Comparative breakdown of PostgreSQL `deadlock_timeout` vs MySQL `innodb_lock_wait_timeout` vs SQL Server `LOCK_TIMEOUT`.

Source:
Microsoft Learn Deadlocks Guide, PostgreSQL §13.3.4 & §19.12, MySQL §17.7.5.2

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately clarifies a frequent point of confusion among engineers.

---

## Claim 9

Claim:
PostgreSQL `deadlock_timeout` defaults to 1 second; the server waits this duration before executing the computationally expensive cycle detection algorithm.

Location:
`03-evidence.md` §Evidence 9, `05-report.md` §Finding 8

Evidence Provided:
Direct quote from PostgreSQL 18 Documentation §19.12 runtime configuration.

Source:
PostgreSQL 18 Documentation §19.12

Source Actually Supports Claim:
YES

Classification:
IMPLEMENTATION-SPECIFIC

Severity:
LOW

Notes:
Verified against live PostgreSQL 18 docs.

---

## Claim 10

Claim:
Production deadlock monitoring relies on database error logs (`innodb_print_all_deadlocks`, `log_lock_waits`, SQL Server XEvent) and status inspection (`SHOW ENGINE INNODB STATUS`).

Location:
`05-report.md` §Finding 9, `03-evidence.md` §Evidence 14

Evidence Provided:
MySQL and Percona operational references; PostgreSQL monitoring documentation.

Source:
MySQL §17.7.5.x, PostgreSQL Chapter 27, Percona Blog

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Correctly identifies standard production telemetry and diagnostic tools.
