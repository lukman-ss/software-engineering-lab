# Claim Audit: Optimistic vs Pessimistic Locking

## Claim 1

Claim:
The "lost update" anomaly occurs under default database isolation (READ COMMITTED in PostgreSQL/Oracle) when concurrent transactions execute uncoordinated read-modify-write sequences.

Location:
`research/05-report.md`, Finding 1; `research/03-evidence.md`, Evidence 01, 02

Evidence Provided:
Wikipedia lost update definition citing Bernstein et al. 1987; Oracle 19c Concepts Table 10-2 Banda salary scenario; PostgreSQL 13.2 Read Committed re-evaluation semantics.

Source:
Sources 1, 2, 3, 10

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Solidly supported across multiple database engines and standard textbook definitions.

---

## Claim 2

Claim:
Pessimistic locking via `SELECT ... FOR UPDATE` acquires row-level exclusive locks held until commit/rollback, preventing concurrent UPDATE/DELETE/locking reads on the targeted rows while plain SELECT remains unblocked.

Location:
`research/05-report.md`, Finding 2; `research/03-evidence.md`, Evidence 03, 04

Evidence Provided:
PostgreSQL 13.3.2 row-level lock documentation; MySQL 8.0 locking reads documentation; Oracle TX row lock semantics.

Source:
Sources 1, 7, 10, 16

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Fully verified against PostgreSQL, MySQL/InnoDB, and Oracle documentation.

---

## Claim 3

Claim:
Pessimistic locking increases lock wait time, reduces concurrency, and introduces deadlock risk; holding locks across external HTTP/gateway calls is an anti-pattern.

Location:
`research/05-report.md`, Finding 3, Finding 9; `research/03-evidence.md`, Evidence 05, 06, 13

Evidence Provided:
PostgreSQL 13.3.4 deadlocks section ("applications should not hold transactions open for long periods waiting for user input / external events"); Wikipedia 2PL deadlock discussion.

Source:
Sources 1, 3, 5

Source Actually Supports Claim:
YES

Classification:
FACT / BEST_PRACTICE

Severity:
LOW

Notes:
Consistently warned against in all primary database documentation.

---

## Claim 4

Claim:
Optimistic locking detects conflicts at commit/update time by including the read version/timestamp in the `WHERE` clause (`UPDATE ... WHERE id = ? AND version = ?`), where zero affected rows signals a conflict requiring application-level handling (retry or 409).

Location:
`research/05-report.md`, Finding 4; `research/03-evidence.md`, Evidence 07, 16

Evidence Provided:
Martin Fowler (Optimistic Offline Lock); Oracle 19c WHERE-clause original value recommendation; Hibernate/JPA OptimisticLockException behavior.

Source:
Sources 4, 10, 12, 14

Source Actually Supports Claim:
YES

Classification:
FACT / DESIGN_PATTERN

Severity:
LOW

Notes:
Standard implementation pattern supported by enterprise architecture literature and vendor guides.

---

## Claim 5

Claim:
Atomic single-statement updates (`UPDATE products SET stock = stock - N WHERE id = ? AND stock >= N` with `affected_rows == 1` check) eliminate read-modify-write race windows without requiring explicit multi-statement locks.

Location:
`research/05-report.md`, Finding 6; `research/03-evidence.md`, Evidence 10

Evidence Provided:
PostgreSQL 13.4.2 statement-level consistency ("must actually update the row"); SQL statement-level ACID atomicity from Oracle/PostgreSQL/MySQL.

Source:
Sources 11, 15, 16

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Statement atomicity and predicate evaluation in SQL engines guarantee atomic execution of single DML statements.

---

## Claim 6

Claim:
Default transaction isolation levels vary across database engines (PostgreSQL and Oracle default to READ COMMITTED; MySQL/InnoDB defaults to REPEATABLE READ) and do not uniformly eliminate lost updates.

Location:
`research/05-report.md`, Finding 8; `research/03-evidence.md`, Evidence 11, 12

Evidence Provided:
PostgreSQL 13.2 isolation levels; MySQL 15.7.2.1 isolation levels; Oracle 19c Concepts 9.

Source:
Sources 2, 9, 10

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately reflects differences in vendor defaults and MVCC/2PL implementations.

---

## Claim 7

Claim:
Using distributed locks (e.g. Redis) when the resource lives entirely in a single relational database is an anti-pattern.

Location:
`research/05-report.md`, Finding 9; `research/06-open-questions.md`, OQ-2

Evidence Provided:
Inference from PostgreSQL advisory locks and database-native capabilities.

Source:
Source 1, Source 13

Source Actually Supports Claim:
PARTIAL

Classification:
INTERPRETATION / BEST_PRACTICE

Severity:
MEDIUM

Notes:
While single-DB native locks are universally recommended before introducing Redis/Redlock overhead, the boundary conditions (e.g. high-throughput rate-limiting vs transactional state) require careful architectural nuance. Properly flagged in Open Questions (OQ-2).
