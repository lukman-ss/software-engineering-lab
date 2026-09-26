# Claim Audit: Optimistic vs Pessimistic Locking

## Claim 1

Claim: The "lost update" problem occurs when a second transaction overwrites a first transaction's update without seeing it, causing the first value to be lost.

Location: `research/05-report.md:21-25`, `research/03-evidence.md:7-22`

Evidence Provided: Academic definition from Bernstein et al. 1987 / Weikum & Vossen 2001 via Wikipedia, and concrete salary update example from Oracle 19c Concepts Table 10-2.

Source: Source 3 (Wikipedia Concurrency Control), Source 10 (Oracle 19c Concepts), Source 2 (PostgreSQL Transaction Isolation)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW (Fully Supported)

Notes: Verified against classical database literature and official RDBMS documentation.

---

## Claim 2

Claim: Pessimistic locking via `SELECT ... FOR UPDATE` causes retrieved rows to be locked exclusively until transaction end; other transactions attempting write or locking read on those rows are blocked.

Location: `research/05-report.md:35-46`, `research/03-evidence.md:43-58`

Evidence Provided: PostgreSQL 18 Documentation 13.3.2, MySQL 8.0 InnoDB Locking Reads, Oracle 19c Concepts.

Source: Source 1 (PostgreSQL Explicit Locking), Source 7 (MySQL Locking Reads), Source 10 (Oracle Data Concurrency)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW (Fully Supported)

Notes: Verified via PostgreSQL docs Section 13.3.2 and Table 13.3.

---

## Claim 3

Claim: Pessimistic locking reduces concurrency, increases wait times, and introduces deadlock risks. Applications should avoid holding locks during long operations like network calls or user input.

Location: `research/05-report.md:48-58`, `research/03-evidence.md:80-112`

Evidence Provided: PostgreSQL Documentation 13.3.4 (Deadlocks) and Section 13.3 warning against holding transactions open for long periods.

Source: Source 1 (PostgreSQL 18 Docs 13.3)

Source Actually Supports Claim: YES

Classification: FACT / IMPLEMENTATION-SPECIFIC

Severity: LOW (Fully Supported)

Notes: Matches official PostgreSQL and database engineering principles.

---

## Claim 4

Claim: Optimistic locking validates at commit time that data has not changed since it was read (via version/timestamp check in WHERE clause). A result of 0 affected rows signals a conflict that must be handled by the application (reload/retry/409).

Location: `research/05-report.md:60-73`, `research/03-evidence.md:116-130`, `research/03-evidence.md:278-290`

Evidence Provided: Martin Fowler EAA (Optimistic Offline Lock), Oracle 19c Concepts WHERE-guard recommendation, Hibernate / Baeldung JPA docs.

Source: Source 4 (Fowler Optimistic Offline Lock), Source 10 (Oracle Data Concurrency), Source 12 (Baeldung JPA), Source 14 (Hibernate Guide)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW (Fully Supported)

Notes: Verified with Fowler EAA pattern and standard SQL WHERE-guard semantics.

---

## Claim 5

Claim: Atomic single-statement updates (`UPDATE ... SET stock = stock - N WHERE stock >= N` + checking `affected_rows == 1`) eliminate the read-modify-write race window without explicit locking.

Location: `research/05-report.md:89-100`, `research/03-evidence.md:170-184`

Evidence Provided: Oracle ACID statement atomicity, PostgreSQL single-statement MVCC lock release, MySQL row-locking behavior during update.

Source: Source 11 (Oracle Transactions), Source 1 (PostgreSQL Explicit Locking), Source 9 (MySQL Isolation)

Source Actually Supports Claim: YES

Classification: INTERPRETATION / IMPLEMENTATION-SPECIFIC

Severity: LOW (Supported with proper confidence caveat)

Notes: The research accurately classified this finding as MEDIUM confidence because while single-statement atomicity and WHERE predicates are standard SQL guarantees, the exact counter-decrement recipe is an industry idiom rather than a formal vendor spec quote.

---

## Claim 6

Claim: Wrapping code in a database transaction alone does not automatically prevent lost updates; prevention depends on query structure and isolation level.

Location: `research/05-report.md:102-114`, `research/03-evidence.md:206-220`

Evidence Provided: PostgreSQL Read Committed re-evaluation allows overwrite; Oracle documented example occurs within transactions; MySQL semi-consistent reads.

Source: Source 1, Source 2, Source 9, Source 10

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW (Fully Supported)

Notes: Essential distinction for software engineers.

---

## Claim 7

Claim: Default transaction isolation levels differ across database vendors: PostgreSQL defaults to READ COMMITTED (READ UNCOMMITTED is a no-op identical to READ COMMITTED), Oracle defaults to READ COMMITTED, and MySQL InnoDB defaults to REPEATABLE READ.

Location: `research/05-report.md:116-128`, `research/03-evidence.md:188-202`

Evidence Provided: PostgreSQL Docs 13.2, MySQL Docs 15.7.2.1, Oracle Docs Chapter 9.

Source: Source 1, Source 2, Source 9, Source 10

Source Actually Supports Claim: YES

Classification: FACT / IMPLEMENTATION-SPECIFIC

Severity: LOW (Fully Supported)

Notes: Cross-engine isolation level semantics are correctly distinguished.
