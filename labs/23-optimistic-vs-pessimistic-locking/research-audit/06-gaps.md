# Research Gaps — labs/23-optimistic-vs-pessimistic-locking

Audit Date: 2026-09-25

## Gap 1
Type: MISSING_SOURCE
Severity: HIGH
Location: research/02-sources.md Source 12; research/03-evidence.md:19
Problem: EF Core Concurrency source URL https://learn.microsoft.com/en-us/ef/core/performance/efficient-query-patterns/concurrent-updates returns 404. Claim that EF Core [Timestamp]/rowversion is covered by this URL is unsupported.
Required Revision: Replace with correct EF Core docs URL (e.g., https://learn.microsoft.com/en-us/ef/core/saving/concurrency or /ef/core/modeling/concurrency) and re-verify quote. Do not cite dead URL as evidence.
Can Be Approved Without Fix: NO — blocks Evidence 19 / Finding 15 (optimistic locking ORM support).

## Gap 2
Type: MISSING_SOURCE
Severity: HIGH
Location: research/02-sources.md Source 6; research/05-report.md:28, Appendix
Problem: Oracle SELECT FOR UPDATE docs URL https://docs.oracle.com/en/database/oracle/oracle-database/23/sql/SELECT.html returns 404. ORA_ROWSCN/ FOR UPDATE WAIT/NOWAIT/SKIP LOCKED syntax thus uncited.
Required Revision: Replace with valid Oracle 23c SQL Language Reference URL (e.g., https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html) and re-verify FOR UPDATE clause, WAIT/NOWAIT/SKIP LOCKED, and ORA_ROWSCN pseudo-column existence.
Can Be Approved Without Fix: NO — vendor matrix in appendix relies on Oracle syntax.

## Gap 3
Type: MISSING_SOURCE / WEAK_SOURCE
Severity: HIGH
Location: research/02-sources.md Source 8; research/03-evidence.md:5,92
Problem: Bernstein et al. (1987) PDF URL https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/ccontrol.pdf returns 404. Page-level quotes (p.47 validation phase; p.136 crossover) NOT VERIFIED. Book exists but citation broken.
Required Revision: Replace with stable identifier (Microsoft Research archive, ACM copy, or worldcat/ISBN) and flag page numbers as "NOT VERIFIED" unless inspected. Remove claim that PDF was fetched/verified.
Can Be Approved Without Fix: NO — core theoretical foundation (Q1, Q4, Q9, Q11) built on this source.

## Gap 4
Type: WEAK_SOURCE
Severity: HIGH
Location: research/03-evidence.md:74-79,99-105; research/02-sources.md Source 5
Problem: Performance benchmark numbers (≤10% low contention / >30% high contention / 15-20% crossover) and atomic UPDATE example verbatim-attributed to PostgreSQL/SQL Server docs but not present in fetched pages. Quotes appear fabricated or conflated from secondary commentary.
Required Revision: Either produce correct source with verbatim audit trail, or re-label these as INTERPRETATION/HYPOTHESIS with explicit "NOT VERIFIED — no vendor benchmark found" and remove fake verbatim quotes. Keep Bernstein theoretical claim separate from vendor benchmark claim.
Can Be Approved Without Fix: NO — numeric guidance most likely to be acted on; must be accurately anchored (see Claim Audit Claims 9, 12).

## Gap 5
Type: WEAK_SOURCE
Severity: MEDIUM
Location: research/03-evidence.md:66-71,124-131; research/02-sources.md Sources 7,14
Problem: Hibernate exponential-backoff quote and Azure retry "up to 5 times starting at 10ms" quote not found in fetched pages. Sources reachable but specific claims not evidenced.
Required Revision: Replace verbatim quotes with paraphrase accurately reflecting fetched content (Hibernate: optimistic lock exception + retry is recommended; Azure: exponential backoff with jitter recommended, no single 10ms/5x prescription). Mark as EXAMPLE with scope qualifier.
Can Be Approved Without Fix: YES — if downgraded to paraphrase/EXAMPLE; not blocking core conclusions.

## Gap 6
Type: MISSING_CASE
Severity: HIGH
Location: research/05-report.md:145-153 Appendix; research/03-evidence.md:115-122
Problem: PostgreSQL 14 "introduced SKIP LOCKED/NOWAIT" claim factually wrong (NOWAIT PG9.3, SKIP LOCKED PG9.5). Version-specific implementation guidance thus incorrect.
Required Revision: Correct to "PostgreSQL NOWAIT since 9.3, SKIP LOCKED since 9.5; PG 14 behavior unchanged for these options" citing https://www.postgresql.org/docs/9.5/sql-select.html and current docs. Fix Evidence 14 and Finding 3 scope note.
Can Be Approved Without Fix: NO — version guidance is actionable and currently wrong.

## Gap 7
Type: OVERGENERALIZATION
Severity: MEDIUM
Location: research/03-evidence.md:82-88; research/02-sources.md Source 11
Problem: Laravel `$timestamps = true` with `updated_at` presented as optimistic locking mechanism. Laravel docs describe this as automatic timestamp maintenance, not version-based conflict detection.
Required Revision: Remove `$timestamps`/updated_at as optimistic locking example. Replace with correct Laravel optimistic pattern (manual version column check in UPDATE WHERE, or note that Laravel has no built-in optimistic locking — as documented). Keep lockForUpdate()/sharedLock() pessimistic part.
Can Be Approved Without Fix: NO for publication — misleading framework guidance; YES for internal draft with warning.

## Gap 8
Type: MISSING_SOURCE
Severity: MEDIUM
Location: research/03-evidence.md:109-113; research/05-report.md:74-79
Problem: Distributed DB claims (CockroachDB HLC, Spanner TrueTime, TiDB) inferred without direct vendor doc URL; marked as MEDIUM confidence but no URL to verify.
Required Revision: Add direct vendor URLs (e.g., cockroachlabs docs on transaction layer, Spanner whitepaper) or mark as HYPOTHESIS and move detail to open questions. Do not present as FACT with vendor attribution.
Can Be Approved Without Fix: YES — if reclassified to HYPOTHESIS/INTERPRETATION and scoped.

## Gap 9
Type: MISSING_SOURCE
Severity: MEDIUM
Location: research/03-evidence.md:47-52; research/05-report.md:82-85
Problem: Silberschatz page numbers (p.695, p.784) and chapter numbers (15 vs 17) inconsistent within evidence; not verifiable without physical book inspection. Quoted sentence not fetch-verifiable.
Required Revision: Keep book citation (ISBN/edition) but mark page-level quotes as NOT VERIFIED; reconcile chapter (15 = Concurrency Control in 7th ed) and remove approximate page numbers unless inspected.
Can Be Approved Without Fix: YES — book is real and concept is standard; page precision not material.

## Gap 10
Type: MISSING_SOURCE (honestly disclosed)
Severity: LOW
Location: research/02-sources.md Source 10; research/06-open-questions.md:4-6
Problem: MySQL 8.0 docs 403 — behavior correctly marked "NOT VERIFIED / inferred from cross-vendor patterns" and listed in open questions.
Required Revision: Re-verify when MySQL docs accessible; until then keep "(NOT VERIFIED)" in appendix and evidence notes. No incorrect MySQL claim made — good handling.
Can Be Approved Without Fix: YES — properly disclosed; not blocking.

## Gap 11
Type: WEAK_SOURCE
Severity: LOW
Location: research-revision/02-changes-made.md:20-28; research-revision/03-revision-result.md:42-45
Problem: Revision claims "17 verifiable sources, all URLs verified reachable, excerpts match" — overstated. Actual distinct sources = 12; 4 URLs dead on audit fetch (6,8,12,10). Overstates remediation.
Required Revision: Correct revision-result counts to 12 sources, 6 reachable/4 unreachable (1 correctly disclosed), and downgrade "RESOLVED" to "PARTIALLY RESOLVED" for source integrity.
Can Be Approved Without Fix: YES — process hygiene, not research correctness.
