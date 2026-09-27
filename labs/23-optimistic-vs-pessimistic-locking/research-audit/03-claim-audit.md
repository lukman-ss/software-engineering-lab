# Claim Audit: Optimistic vs Pessimistic Locking

## Claim 1

Claim: A lost update occurs when two transactions read the same data, both modify it based on the old value, and the second commit overwrites the first.

Location: `research/03-evidence.md:3-11`, `research/05-report.md:9`

Evidence Provided: Wikipedia "Write-write conflict", PostgreSQL docs 13.2, Oracle 19c Table 10-2.

Source: Wikipedia, PostgreSQL docs, Oracle docs.

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Standard definition in database literature.

---

## Claim 2

Claim: Default READ COMMITTED isolation level in PostgreSQL and Oracle does not prevent lost updates.

Location: `research/03-evidence.md:57-73`, `research/05-report.md:20-30`

Evidence Provided: PostgreSQL docs explain that subsequent statements re-evaluate against latest committed row versions; Oracle Table 10-2 shows Session 2 overwriting Session 1's committed change.

Source: PostgreSQL Docs 13.2, Oracle 19c Data Concurrency and Consistency.

Source Actually Supports Claim: YES

Classification: FACT

Severity: HIGH (Critical insight that counters naive assumption)

Notes: Verified against primary sources.

---

## Claim 3

Claim: Pessimistic locking via `SELECT ... FOR UPDATE` acquires an exclusive row lock that blocks concurrent writers until transaction completion.

Location: `research/03-evidence.md:13-29`, `research/05-report.md:32-48`

Evidence Provided: PostgreSQL 13.3 explicitly describes row-level `FOR UPDATE` blocking `UPDATE`, `DELETE`, and locking queries; Oracle TX locks behave identically.

Source: PostgreSQL Docs 13.3, Oracle 19c Concepts.

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Verified directly.

---

## Claim 4

Claim: Explicit pessimistic locking increases the risk of deadlocks, which databases detect automatically and resolve by aborting one transaction.

Location: `research/03-evidence.md:93-100`, `research/05-report.md:41-44`

Evidence Provided: PostgreSQL 13.3.4 explicitly documents lock inversion deadlocks and automated deadlock resolution via transaction abort.

Source: PostgreSQL Docs 13.3.4.

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Verified directly.

---

## Claim 5

Claim: Optimistic concurrency control (OCC) operates across Begin, Modify, Validate, Commit phases, detecting conflicts before committing without holding locks during data processing.

Location: `research/03-evidence.md:39-46`, `research/05-report.md:50-58`

Evidence Provided: Kung & Robinson (1981) cited via Wikipedia OCC.

Source: Wikipedia Optimistic Concurrency Control.

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Theoretical basis accurately reflected.

---

## Claim 6

Claim: In application systems, optimistic locking is implemented using `UPDATE ... SET ..., version = version + 1 WHERE id = ? AND version = ?`, detecting conflicts when affected rows equal 0.

Location: `research/03-evidence.md:48-56`, `research/05-report.md:51-57`

Evidence Provided: Microsoft EF Core documentation details concurrency token check and `DbUpdateConcurrencyException`.

Source: Microsoft Learn EF Core Concurrency.

Source Actually Supports Claim: YES

Classification: IMPLEMENTATION-SPECIFIC

Severity: LOW

Notes: Accurate description of application-level OCC pattern.

---

## Claim 7

Claim: Under PostgreSQL REPEATABLE READ, concurrent updates to rows modified by another transaction raise `ERROR: could not serialize access due to concurrent update`.

Location: `research/03-evidence.md:57-65`, `research/05-report.md:71-76`

Evidence Provided: PostgreSQL docs 13.2.2.

Source: PostgreSQL Docs 13.2.2.

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Verbatim match with PostgreSQL documentation.

---

## Claim 8

Claim: Under Oracle SERIALIZABLE, modifying rows changed by another transaction after transaction start raises `ORA-08177: Cannot serialize access for this transaction`.

Location: `research/03-evidence.md:66-74`, `research/05-report.md:76`

Evidence Provided: Oracle 19c Concepts Chapter 10 Table 10-3.

Source: Oracle 19c Concepts.

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Verbatim match with Oracle documentation.

---

## Claim 9

Claim: Atomic database operations (`UPDATE ... SET stock = stock - N WHERE id = ? AND stock >= N`) eliminate the race window for counter/decrement operations without explicit locking.

Location: `research/03-evidence.md:120-127`, `research/05-report.md:83-94`

Evidence Provided: PostgreSQL documentation and standard SQL single-statement atomicity semantics.

Source: PostgreSQL Docs 13.2; SQL Standard.

Source Actually Supports Claim: PARTIAL

Classification: INTERPRETATION

Severity: MEDIUM

Notes: While single-row conditional updates are atomic and eliminate the read-modify-write anomaly on that row, multi-row business constraints or complex predicates require explicit locking or higher isolation levels. The research correctly acknowledges this boundary in `research/06-open-questions.md:67`.

---

## Claim 10

Claim: MySQL/InnoDB implements `SELECT ... FOR UPDATE` and gap locking in REPEATABLE READ.

Location: `research/03-evidence.md:31-38`, `research/04-contradictions.md:20-27`

Evidence Provided: Secondary references; primary MySQL documentation was inaccessible (HTTP 403).

Source: Secondary sources / SQL Standard.

Source Actually Supports Claim: PARTIAL

Classification: HYPOTHESIS / UNVERIFIED_PRIMARY

Severity: MEDIUM

Notes: The research author transparently disclosed that MySQL docs returned 403 and marked confidence as MEDIUM. This is honest research practice, but the primary source remains unverified.
