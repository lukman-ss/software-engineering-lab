# Claim Audit — labs/23-optimistic-vs-pessimistic-locking

Audit Date: 2026-09-25
Method: Extract every claim from research/03-evidence.md, research/05-report.md, and 02-sources.md. Match claim to source URL, fetch if needed, verify if quote is verbatim.

## Claim 1
Claim: Optimistic concurrency control validates transactions at commit time; pessimistic concurrency control prevents conflicts via upfront locking.
Location: evidence.md:4-5, report.md:12-13
Evidence Provided: Bernstein et al. (1987) p.47; Silberschatz et al. (2019) p.695
Source: "Concurrency Control and Recovery" (Source 8, URL 404); "Database System Concepts" (Source 9)
Source Actually Supports Claim: PARTIAL — Bernstein PDF URL 404 (FAIL); Silberschatz book page citable but not independently verifiable (no ISBN/page fetch possible). Quote is standard textbook knowledge but citation unverifiable per anti-hallucination rules.
Classification: FACT
Severity: HIGH (source URL dead)
Notes: Claim itself is textbook-correct; citation method for Bernstein PDF is broken.

## Claim 2
Claim: MVCC provides the mechanism enabling optimistic locking by allowing readers to access consistent snapshots without blocking writers.
Location: evidence.md:12-18, report.md:19-23
Evidence Provided: PostgreSQL docs Section 13.1
Source: PostgreSQL MVCC Introduction (Source 4)
Source Actually Supports Claim: YES — fetched page contains: "locks acquired for querying (reading) data do not conflict with locks acquired for writing data, and so reading never blocks writing and writing never blocks reading."
Classification: FACT
Severity: LOW
Notes: Source PASS; quote verified verbatim.

## Claim 3
Claim: Pessimistic locking is implemented via SELECT FOR UPDATE (PostgreSQL, Oracle), SELECT ... WITH (UPDLOCK) (SQL Server), acquiring exclusive row locks that block concurrent writers.
Location: evidence.md:21-27, report.md:26-30
Evidence Provided: PostgreSQL docs Section 13.3.2
Source: PostgreSQL Explicit Locking (Source 1)
Source Actually Supports Claim: YES — fetched page contains exact quote about FOR UPDATE locking rows. Oracle and SQL Server behavior inferred from PostgreSQL.
Classification: FACT
Severity: MEDIUM — scope-limited (PostgreSQL-focused); vendor-specific details for SQL Server/Oracle rely on inference rather than independent source verification.
Notes: SQL Server URL (Source 5) over-cited for this claim; Oracle URL (Source 6) dead.

## Claim 4
Claim: Optimistic locking is implemented via version columns checked in UPDATE WHERE clauses, ORM-managed @Version/[Timestamp] annotations, or atomic conditional UPDATE statements.
Location: evidence.md:29-36, report.md:32-37
Evidence Provided: SQL Server Row Versioning Guide (Source 5), PostgreSQL atomic UPDATE (Source 3)
Source: SQL Server docs (Source 5); PostgreSQL applevel-consistency (Source 3)
Source Actually Supports Claim: PARTIAL — PostgreSQL Section 13.4.2 does not contain the verbatim UPDATE accounts example quoted in evidence. SQL Server docs support row-versioning concept but not specific ORM behavior.
Classification: FACT (concept true)
Severity: MEDIUM — specific SQL example not verifiable at cited location.
Notes: Evidence 12's UPDATE accounts SET balance... example not found in fetched applevel-consistency page.

## Claim 5
Claim: Pessimistic locking behavior under concurrent updates varies by isolation level.
Location: evidence.md:38-45, report.md:39-44
Evidence Provided: PostgreSQL docs Section 13.3.2
Source: PostgreSQL Explicit Locking (Source 1)
Source Actually Supports Claim: YES — fetched page contains: "In Repeatable Read or Serializable transactions, however, an error will be thrown if a row to be locked has changed since the transaction started."
Classification: FACT
Severity: LOW
Notes: Scope correctly noted as PostgreSQL; other vendors implied consistent but not verified.

## Claim 6
Claim: Both strategies prevent lost update anomaly when correctly implemented.
Location: evidence.md:47-53, report.md:81-86
Evidence Provided: Silberschatz et al. (2019, p. 784)
Source: Database System Concepts (Source 9)
Source Actually Supports Claim: PARTIAL — page reference not independently verifiable via fetch (no URL for book). Academic source needed external verification not possible.
Classification: FACT
Severity: MEDIUM
Notes: Concept is fundamental, but cited resource unverifiable per audit rules.

## Claim 7
Claim: Pessimistic locking can cause deadlocks; optimistic locking avoids deadlocks but may cause validation failures.
Location: evidence.md:55-62, report.md:95-100
Evidence Provided: PostgreSQL docs Section 13.3.4
Source: PostgreSQL Explicit Locking (Source 1)
Source Actually Supports Claim: YES — fetched page contains: "use of explicit locking can increase the likelihood of deadlocks...best defense...acquire locks on multiple objects in a consistent order."
Classification: FACT
Severity: LOW
Notes: Quote verified.

## Claim 8
Claim: Optimistic locking requires application-level retry logic when validation fails due to concurrent updates.
Location: evidence.md:64-71, report.md:102-107
Evidence Provided: Hibernate docs Section 11.1
Source: Hibernate ORM User Guide (Source 7)
Source Actually Supports Claim: PARTIAL — fetched page lists locking sections but exact quote "must be prepared to handle OptimisticLockException and typically implement retry logic with exponential backoff" not found in fetched content (truncated). Concept supported by TOC presence of optimistic locking section.
Classification: FACT
Severity: MEDIUM
Notes: Quote attribution questionable; full page not searchable in truncated fetch.

## Claim 9
Claim: Under high contention (>30% conflicting updates), pessimistic locking outperforms optimistic locking.
Location: evidence.md:73-79, report.md:53-58
Evidence Provided: Microsoft SQL Server docs Performance section
Source: SQL Server docs (Source 5)
Source Actually Supports Claim: PARTIAL — fetched page does NOT contain the quoted "≤10% / >30%" benchmark numbers. Bernstein (1987) p.136 cited for theoretical analysis.
Classification: INTERPRETATION
Severity: HIGH
Notes: CRITICAL: Numeric benchmark paragraph attributed to SQL Server docs is NOT PRESENT in the fetched doc. Claim is INTERPRETATION with theoretical foundation (Bernstein) but the specific vendor benchmark quote is fabricated.

## Claim 10
Claim: Laravel lockForUpdate()/sharedLock() method generates native FOR UPDATE/SHARE SQL; Eloquent $timestamps with updated_at enables optimistic locking.
Location: evidence.md:81-88, report.md:110-114, evidence.md:170-177 (Laravel Evidence 20)
Evidence Provided: Laravel docs (Source 11)
Source: Laravel Database: Locking docs (Source 11)
Source Actually Supports Claim: PARTIAL — Laravel doc shows lockForUpdate()/sharedLock() correctly. The phrase "Eloquent models can use $timestamps = true with updated_at for optimistic locking" (Evidence 10) is NOT in the Laravel locking page; $timestamps is automatic timestamp handling, NOT an optimistic locking mechanism per Laravel docs.
Classification: FACT (pessimistic part) / MISLEADING (optimistic $timestamps)
Severity: HIGH
Notes: The claim conflates Laravel's automatic timestamp maintenance with optimistic locking; Laravel uses $timestamps for auto-updating created_at/updated_at, not version checking.

## Claim 11
Claim: Optimistic locking preferred for read-heavy, low contention, long user-think-time, distributed systems; pessimistic for write-heavy, high contention, short transactions.
Location: evidence.md:90-96, report.md:60-65
Evidence Provided: Bernstein et al. (1987) p.135-138
Source: Concurrency Control and Recovery (Source 8)
Source Actually Supports Claim: PARTIAL — Source 8 URL 404; page numbers not verifiable. Concept aligns with theory but cited resource not accessible.
Classification: FACT (general consensus)
Severity: MEDIUM — source inaccessible.

## Claim 12
Claim: Atomic UPDATE with WHERE clause comparing original values provides lock-free optimistic concurrency.
Location: evidence.md:98-105, report.md:116-121
Evidence Provided: PostgreSQL docs Section 13.4.2
Source: PostgreSQL App-Level Consistency (Source 3)
Source Actually Supports Claim: PARTIAL — claim conceptually true for all SQL databases. Specific `UPDATE accounts SET balance = balance - 100.00 WHERE acctnum = 12345 AND balance = 500.00` example NOT found in fetched applevel-consistency page (which covers locking, not this example).
Classification: FACT
Severity: MEDIUM
Notes: Example not verifiable from cited section despite concept being solid.

## Claim 13
Claim: Distributed databases (CockroachDB, Spanner, TiDB) use hybrid logical clocks/timestamp ordering for optimistic concurrency.
Location: evidence.md:107-113, report.md:74-79
Evidence Provided: CockroachDB Architecture docs (URL not provided, noted "inferred from academic sources")
Source: Research notes cite HLC/Timestamp ordering theory
Source Actually Supports Claim: NOT VERIFIED — URL "inaccessible due to environment" (claimed); no direct vendor URL available. Concept correct but source unverifiable.
Classification: FACT
Severity: MEDIUM — source URL not independently reachable.

## Claim 14
Claim: PostgreSQL 14 introduced SKIP LOCKED and NOWAIT options for SELECT FOR UPDATE.
Location: evidence.md:115-122, report.md:115-119
Evidence Provided: PostgreSQL 14 Release Notes (https://www.postgresql.org/docs/current/release-14.html)
Source: PostgreSQL 14 Release Notes (Source 13, not numbered in sources.md)
Source Actually Supports Claim: NO — fetched PG14 release notes do NOT mention SKIP LOCKED or NOWAIT context; both features introduced earlier (NOWAIT PG9.3, SKIP LOCKED PG9.5). Claim is historically incorrect.
Classification: IMPLEMENTATION-SPECIFIC (incorrect)
Severity: HIGH
Notes: CRITICAL: PostgreSQL 19 Beta 4 notes "Now available" for some features, but SKIP LOCKED/NOWAIT exist in earlier versions. This is a factual error, not just scoping.

## Claim 15
Claim: Retry limits: 3-5 attempts with exponential backoff starting at 10-50ms.
Location: evidence.md:124-131, report.md:103-107
Evidence Provided: Microsoft Azure Architecture Guide: Retry Pattern
Source: Azure transient-faults guidelines (Source 14, accessed in audit)
Source Actually Supports Claim: NO — fetched Azure page discusses retry strategies but does NOT contain "retry up to 5 times with exponential backoff starting at 10ms". That specific wording is not present.
Classification: EXAMPLE
Severity: HIGH
Notes: Quote fabricated or misattributed; retry patterns discussed but exact parameters not evidenced.

## Claim 16
Claim: Silberschatz 7th ed 2019 Chapter 15 covers Concurrency Control.
Location: evidence.md:133-138
Evidence Provided: Edition verified, ISBN cited
Source: Physical book / Library of Congress (not URL)
Source Actually Supports Claim: NOT VERIFIED — cannot fetch book page to confirm chapter numbers.
Classification: FACT
Severity: MEDIUM
Notes: Citation exists but not independently verifiable via web fetch.

## Claim 17
Claim: Concrete SQL syntax for locking varies by vendor but follows standard FOR UPDATE patterns.
Location: evidence.md:141-150, report.md:117-121, Appendix in report.md:145-154
Evidence Provided: PostgreSQL, SQL Server, Oracle vendor docs
Source: Cross-referenced vendor docs (Sources 1,5,6)
Source Actually Supports Claim: PARTIAL — PostgreSQL and SQL Server syntax accurately documented. Oracle SQL 404, Laravel syntax accurate. MySQL marked NOT VERIFIED.
Classification: EXAMPLE
Severity: MEDIUM — MySQL missing, Oracle URL dead.

## Claim 18
Claim: Hibernate @Version annotation enables automatic optimistic locking via WHERE clause inclusion.
Location: evidence.md:152-159
Evidence Provided: Hibernate docs Section 11.1.1
Source: Hibernate ORM User Guide (Source 7)
Source Actually Supports Claim: PARTIAL — concept correct per Hibernate features. Specific quote "automatically includes it in UPDATE statements to check for concurrent modifications" not verifiable due to truncated fetch.
Classification: FACT
Severity: MEDIUM
Notes: TOC suggests section exists; full content not searchable.

## Claim 19
Claim: EF Core [Timestamp] attribute maps to rowversion column in UPDATE WHERE.
Location: evidence.md:161-168
Evidence Provided: EF Core docs Concurrency page
Source: Entity Framework Core Concurrency (Source 12) — URL 404
Source Actually Supports Claim: NOT VERIFIED — cited URL invalid; no alternative EF Core URL verified.
Classification: FACT (concept correct)
Severity: HIGH
Notes: EF Core concurrency docs URL dead; specific guidance not evidenced.

## Claim 20
Claim: Laravel lockForUpdate()/sharedLock() methods translate to native FOR UPDATE/SHARE SQL per database driver.
Location: evidence.md:170-175, evidence.md:172 example
Evidence Provided: Laravel docs (Source 11)
Source: Laravel Database: Locking (Source 11)
Source Actually Supports Claim: PARTIAL — lockForUpdate() syntax verified; the specific PHP example `$users = DB::table('users')->where('votes', '>', 100)->lockForUpdate()->get(); generates SQL with FOR UPDATE` not directly shown (doc shows result with FOR UPDATE but implementation example not quoted as shown).
Classification: FACT
Severity: LOW
Notes: Concept verified; specific code example attribution uncertain.

## Summary

| Claim | Classification | Severity | Notes |
|-------|----------------|----------|-------|
| 1 | FACT | HIGH | Bernstein URL 404 |
| 2 | FACT | LOW | PASS |
| 3 | FACT | MEDIUM | Oracle/SQL Server inference |
| 4 | FACT | MEDIUM | SQL example unverified |
| 5 | FACT | LOW | PASS |
| 6 | FACT | MEDIUM | Page not verifiable |
| 7 | FACT | LOW | PASS |
| 8 | FACT | MEDIUM | Quote not searchable |
| 9 | INTERPRETATION | HIGH | Benchmark quote NOT in SQL Server doc |
| 10 | FACT/MISLEADING | HIGH | $timestamps ≠ optimistic locking |
| 11 | FACT | MEDIUM | Source URL 404 |
| 12 | FACT | MEDIUM | Example not in cited section |
| 13 | FACT | MEDIUM | Vendor URL inaccessible |
| 14 | IMPLEMENTATION-SPECIFIC | HIGH | PG14 did NOT "introduce" NOWAIT/SKIP LOCKED |
| 15 | EXAMPLE | HIGH | Retry pattern quote NOT in Azure doc |
| 16 | FACT | MEDIUM | Cannot fetch book |
| 17 | EXAMPLE | MEDIUM | MySQL/Oracle sources fail |
| 18 | FACT | MEDIUM | Truncated fetch |
| 19 | FACT | HIGH | URL 404 |
| 20 | FACT | LOW | PASS (partial) |