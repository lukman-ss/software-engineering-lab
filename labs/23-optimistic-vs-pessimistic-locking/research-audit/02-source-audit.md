# Source Audit — labs/23-optimistic-vs-pessimistic-locking

Audit Date: 2026-09-25
Method: Direct fetch via WebFetch for every URL; publisher/title matched against fetched page `<title>` and metadata.

## Source 1
Claimed Title: PostgreSQL Documentation: Explicit Locking (Section 13.3)
Claimed Publisher: The PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/explicit-locking.html
Reachable: YES (HTTP 200, PostgreSQL 18, fetched 2026-09-25)
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES — documents FOR UPDATE/FOR SHARE/FOR KEY SHARE, row-level conflict matrix (Table 13.3), deadlock handling (13.3.4).
Problems: Title slightly off — page title is "13.3. Explicit Locking" not "Explicit Locking (Section 13.3)" — trivial. Content matches claimed scope.
Assessment: PASS

## Source 2
Claimed Title: PostgreSQL Documentation: Transaction Isolation (Section 13.2)
Claimed Publisher: The PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/transaction-iso.html
Reachable: YES (HTTP 200, fetched 2026-09-25, title "13.2. Transaction Isolation")
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES — defines READ COMMITTED / REPEATABLE READ / SERIALIZABLE, snapshot behavior, serialization anomaly.
Problems: None.
Assessment: PASS

## Source 3
Claimed Title: PostgreSQL Documentation: Data Consistency Checks at the Application Level (Section 13.4)
Claimed Publisher: The PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/applevel-consistency.html
Reachable: YES (HTTP 200, fetched 2026-09-25, title "13.4. Data Consistency Checks at the Application Level")
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES — documents SELECT FOR UPDATE vs serializable for consistency checks, "temporarily blocks" nuance.
Problems: None for source itself. Note: Evidence 12's verbatim quote attributed to 13.4.2 is not found in fetched page (see Claim Audit).
Assessment: PASS

## Source 4
Claimed Title: PostgreSQL Documentation: MVCC Introduction (Section 13.1)
Claimed Publisher: The PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/mvcc-intro.html
Reachable: YES (HTTP 200, fetched 2026-09-25; page title is "13.1. Introduction" — redirects from /mvcc-intro.html)
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES — MVCC sentence quoted in Evidence 2 is verbatim in fetched page.
Problems: Minor title mismatch ("MVCC Introduction" vs "Introduction"). Section number shifted in PG 18 (13.1 is Introduction, not MVCC Introduction) but same content.
Assessment: PASS

## Source 5
Claimed Title: Microsoft SQL Server Documentation: Transaction Locking and Row Versioning Guide
Claimed Publisher: Microsoft
URL: https://learn.microsoft.com/en-us/sql/relational-databases/sql-server-transaction-locking-and-row-versioning-guide
Reachable: YES (HTTP 200, ms.date 2026-04-22, document_id verified, fetched 2026-09-25)
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: PARTIAL — documents locking vs row versioning, UPDLOCK/ROWLOCK, isolation levels, optimistic concurrency theory. Does NOT contain performance benchmark numbers or retry guidance claimed in Evidence 4/9/15.
Problems: Two verbatim quotes attributed to this source (Evidence 4 row-versioning paragraph, Evidence 9 "≤10% / >30%" benchmark paragraph) not found in fetched page via search. Source supports general concepts, not the specific numeric/benchmark quotes.
Assessment: WARNING — source is valid but over-cited for numeric claims it does not contain.

## Source 6
Claimed Title: Oracle Database SQL Language Reference: SELECT Statement
Claimed Publisher: Oracle Corporation
URL: https://docs.oracle.com/en/database/oracle/oracle-database/23/sql/SELECT.html
Reachable: NO (HTTP 404 fetched 2026-09-25)
Source Type: PRIMARY (claimed)
Relevant: PARTIAL — Oracle does document FOR UPDATE [WAIT/NOWAIT/SKIP LOCKED] but not at this URL.
Supports Claimed Topic: NOT VERIFIED — page not reachable; ORA_ROWSCN claim not verifiable from this URL.
Problems: URL dead. Oracle docs moved; current 23c path differs. No access date correction. Vendor fact (FOR UPDATE syntax) plausible but not evidenced by this URL.
Assessment: FAIL — URL invalid; citation does not support claim as cited.

## Source 7
Claimed Title: Hibernate ORM User Guide: Locking
Claimed Publisher: Red Hat Hibernate Team
URL: https://docs.jboss.org/hibernate/orm/6.6/userguide/html_single/Hibernate_User_Guide.html#locking
Reachable: YES (HTTP 200, version 6.6.58.Final chapter 11 Locking, fetched 2026-09-25; section anchors include locking, locking-optimistic)
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES — @Version, optimistic/pessimistic locking modes exist. Verbatim quote in Evidence 8 ("must be prepared to handle OptimisticLockException and typically implement retry logic with exponential backoff") not found verbatim in fetched TOC/page; appears paraphrased.
Problems: Over-precise quote attribution; URL anchor #locking may not match exact subsection IDs (#locking-optimistic etc.). Content supports concepts but not the verbatim sentence.
Assessment: WARNING — source valid, specific quote not verified verbatim.

## Source 8
Claimed Title: "Concurrency Control and Recovery in Database Systems" by Bernstein, Hadzilacos, Goodman
Claimed Publisher: Addison-Wesley (1987), PDF at https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/ccontrol.pdf
URL: https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/ccontrol.pdf
Reachable: NO (HTTP 404 fetched 2026-09-25)
Source Type: PRIMARY (academic)
Relevant: YES (foundational OCC/2PL theory)
Supports Claimed Topic: NOT VERIFIED via URL — book exists (worldcat/ISBN) but cited PDF URL dead. No alternative stable URL verified. Page numbers (p.47, p.136) cited in Evidence 1/11 not verifiable from this URL.
Problems: Dead URL presented as accessible; page-level quotes unverified. Research-revision says "verified" but fetch fails. Publication still real but citation broken.
Assessment: FAIL — URL invalid; page-level quotes NOT VERIFIED.

## Source 9
Claimed Title: "Database System Concepts" by Silberschatz, Korth, Sudarshan (7th Ed 2019)
Claimed Publisher: McGraw-Hill Education
URL: (none provided; ISBN 978-0078022159 only)
Reachable: NOT APPLICABLE — physical book, no URL to fetch.
Source Type: PRIMARY (academic textbook)
Relevant: YES
Supports Claimed Topic: PARTIAL — chapter/page claims (Chapter 15 vs 17, p.695, p.784) inconsistent across evidence/report; not verifiable without book inspection. Chapter 15 = Concurrency Control in 7th ed is plausible, but Evidence 6 cites Chapter 17 incorrectly.
Problems: Edition confusion (Evidence 6 cites Chapter 17 p.784; Report Finding 1 cites p.695). Neither page verified. No URL does not equal failure, but page-level quotes remain NOT VERIFIED.
Assessment: WARNING — book is real and relevant, but page/chapter citations not verified and internally inconsistent.

## Source 10
Claimed Title: MySQL 8.0 Reference Manual: InnoDB Locking Reads
Claimed Publisher: Oracle Corporation (MySQL)
URL: https://dev.mysql.com/doc/refman/8.0/en/innodb-locking-reads.html
Reachable: NO (HTTP 403 "Technical Difficulties" fetched 2026-09-25) — matches research disclosure.
Source Type: PRIMARY (claimed)
Relevant: YES
Supports Claimed Topic: NOT VERIFIED — site returns 403 across 8.0 and 5.7 per research.
Problems: None in disclosure — research correctly marks "NOT VERIFIED" and marks MySQL behavior as inferred. This is proper handling.
Assessment: WARNING — source correctly flagged as unverified; not counted as supporting evidence (research does this correctly).

## Source 11
Claimed Title: Laravel Database Documentation: Locking
Claimed Publisher: Laravel
URL: https://laravel.com/docs/11.x/queries#locking-rows
Reachable: YES (HTTP 200, title "Database: Query Builder | Laravel 11.x", section "Pessimistic Locking" exists, fetched 2026-09-25)
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: PARTIAL — lockForUpdate()/sharedLock() pessimistic locking is correctly documented. Optimistic locking via "$timestamps = true with updated_at" (Evidence 10) is NOT documented there; Laravel $timestamps is automatic timestamp maintenance, not optimistic locking.
Problems: Anchor #locking-rows not a real anchor (actual is #pessimistic-locking). Optimistic claim via $timestamps is inaccurate.
Assessment: WARNING — pessimistic part supported, optimistic $timestamps claim unsupported by this URL.

## Source 12
Claimed Title: Entity Framework Core Documentation: Concurrency
Claimed Publisher: Microsoft
URL: https://learn.microsoft.com/en-us/ef/core/performance/efficient-query-patterns/concurrent-updates
Reachable: NO (HTTP 404 fetched 2026-09-25)
Source Type: PRIMARY (claimed)
Relevant: PARTIAL — EF Core concurrency docs exist but at different path (/ef/core/modeling/concurrency or /ef/core/saving/concurrency).
Supports Claimed Topic: NOT VERIFIED at cited URL — [Timestamp]/rowversion behavior plausible but not evidenced by this URL. Evidence 19's verbatim quote not found.
Problems: URL incorrect; claimed evidence cannot be traced to this page.
Assessment: FAIL — URL invalid.

## Additional Source Referenced in Evidence (not in 02-sources.md)
- PostgreSQL 14 Release Notes (https://www.postgresql.org/docs/current/release-14.html) — reachable YES (fetched), but Evidence 14's claim that PG14 introduced SKIP LOCKED/NOWAIT is factually wrong (NOWAIT since PG 9.3, SKIP LOCKED since PG 9.5).
- Microsoft Azure Architecture Guide: Retry Pattern (https://learn.microsoft.com/en-us/azure/architecture/best-practices/transient-faults) — reachable YES (fetched 2026-09-25, title "Transient Fault Handling"). Does NOT contain the quoted "retry up to 5 times with exponential backoff starting at 10ms" — that sentence not found in fetched page.
- PostgreSQL Transaction Isolation doc (Source 2 already covers) — reachable.

## Summary Counts
Total claimed unique URLs in 02-sources.md: 10 URLs (Sources 1-7,10-12)
Reachable (200): 6 (1,2,3,4,5,11) — plus 7 partial
Unreachable (404/403): 4 (6,8,10,12) — 10 correctly disclosed as not verified
PASS: 4, WARNING: 5, FAIL: 3
Revision claim "17 verifiable sources" (research-revision/02-changes-made.md) overstates count; actual distinct sources in 02-sources.md = 12.
