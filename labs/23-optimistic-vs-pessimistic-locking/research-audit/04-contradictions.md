# Contradictions — labs/23-optimistic-vs-pessimistic-locking

Audit Date: 2026-09-25

## Contradiction 1

Statement A: "under low contention (≤10% conflicting updates), optimistic concurrency reduces blocking; under high conflict (>30% conflicting updates), pessimistic locking... performs better" (Evidence 9, source-attributed to SQL Server docs).

Location: research/03-evidence.md:74-76

Statement B: "Under low contention (<15% update conflict rate), optimistic... yields higher throughput; under high contention (>30% conflict rate)... Bernstein theoretical analysis" (Report Finding 7). Areas of Disagreement notes "No consensus on exact contention crossover point (15-20% from SQL Server benchmarks, but highly workload-dependent)" (report.md:131).

Location: research/05-report.md:54-58, 131

Type: INTERNAL

Impact: Numeric thresholds inconsistent within same research (≤10% vs <15% vs 15-20%). Reader cannot determine which threshold the research endorses. Undermines credibility of performance guidance.

Assessment: Material — numeric thresholds presented as evidence-based but internally inconsistent and, per Source Audit, not actually present in cited vendor doc.

## Contradiction 2

Statement A: PostgreSQL 14 introduced SKIP LOCKED and NOWAIT for SELECT FOR UPDATE; earlier versions only supported WAIT. (Evidence 14: "PostgreSQL 14 introduced... earlier versions only supported WAIT")

Location: research/03-evidence.md:117, research/02-sources.md:50

Statement B: PostgreSQL explicit locking docs (Section 13.3) document NOWAIT and SKIP LOCKED as available in current (PG 18) without version qualifier implying PG 14 novelty; fetched PG14 release notes contain no such introduction note. Historical record: NOWAIT since PG 9.3, SKIP LOCKED since PG 9.5.

Location: research/03-evidence.md:119 (URL), fetched https://www.postgresql.org/docs/current/release-14.html (no SKIP LOCKED/NOWAIT mention)

Type: SOURCE_CONFLICT

Impact: Factually incorrect version claim creates implementation risk — users on PG 12/13 would incorrectly believe NOWAIT/SKIP LOCKED unavailable.

Assessment: Material — incorrect vendor history.

## Contradiction 3

Statement A: Revision Result claims "All sources: URLs verified reachable, titles/publishers confirmed... PostgreSQL docs accessed and excerpts verified... SQL Server fetched... Oracle docs fetched... Hibernate fetched... All implementation claims scoped... No unqualified universal claims" (02-changes-made.md:20-28, 03-revision-result.md:39-45).

Location: research-revision/02-changes-made.md:19-29, research-revision/03-revision-result.md:42-45

Statement B: Direct audit fetch 2026-09-25: Oracle SELECT URL 404, Bernstein PDF 404, EF Core 404, MySQL 403, SQL Server performance quote not in doc, Oracle ORA_ROWSCN not in doc, retry pattern quote not in Azure doc. At least 4 of 10 URLs dead.

Location: research-audit/02-source-audit.md (Sources 6, 8, 12, 10)

Type: INTERNAL (revision claim vs reality)

Impact: Revision verification overstates source integrity; masks remaining gaps documented elsewhere (06-open-questions.md correctly notes MySQL inaccessibility, but revision claims "RESOLVED" for source integrity).

Assessment: Material — undermines trust in revision process.

## Contradiction 4

Statement A: "Source 10 MySQL ... NOT VERIFIED (Technical Difficulties)... MySQL locking behavior in this research inferred from PostgreSQL/SQL Server/Oracle patterns" (02-sources.md:99-102) and "MySQL 8.0 official documentation inaccessible... behavior inferred" (06-open-questions.md:4-6).

Location: research/02-sources.md:99-102, research/06-open-questions.md:4-6

Statement B: Report Appendix quick reference includes MySQL column: "`SELECT ... FOR UPDATE [NOWAIT/SKIP LOCKED]`" and "`App-managed version in UPDATE WHERE`" (report.md:149-151) and Finding 3 scope says "PostgreSQL 14+ supports NOWAIT/SKIP LOCKED" while implying MySQL parity.

Location: research/05-report.md:145-153

Type: INTERNAL

Impact: Research correctly discloses MySQL not verified, but appendix presents MySQL syntax as if verified (header says "MySQL (NOT VERIFIED)" only in column header — subtle disclaimer). Mixed signal: is MySQL syntax trusted or inferred?

Assessment: Low — disclosed correctly, but presentation invites misread.

## Contradiction 5

Statement A: Evidence 10 / Finding 15 claim Laravel Eloquent `lockForUpdate()` and pessimistic locking, and "`$timestamps = true;` with `updated_at` for optimistic locking" (evidence.md:83, report.md:111).

Location: research/03-evidence.md:83, research/05-report.md:153

Statement B: Evidence 20 / report.md:111 correctly notes `lockForUpdate()` translates to FOR UPDATE per fetched Laravel docs (pessimistic). But `$timestamps` is NOT an optimistic locking mechanism in Laravel docs — it's auto-maintained timestamps.

Location: research/02-sources.md:111, fetched https://laravel.com/docs/11.x/queries#locking-rows (section is "Pessimistic Locking" only; no optimistic via $timestamps documented there)

Type: INTERNAL

Impact: ORM claim overstates Laravel capability; reader would implement `$timestamps` expecting version-based conflict detection and get none.

Assessment: Material — misleading framework guidance.

## Contradiction 6

Statement A: "Source 9: Chapter 15 covers concurrency control" (revision 02-changes-made.md:41, 03-evidence.md:6). Report Finding 1 cites "Silberschatz et al. (2019, p. 695)" for PCC definition.

Location: research/03-evidence.md:6, research-revision/02-changes-made.md:41

Statement B: Evidence 6 cites "Silberschatz et al. (2019, p. 784)" for same book "Chapter 17" lost update claim; report Finding 11 repeats p.784.

Location: research/03-evidence.md:49-50

Type: INTERNAL

Impact: Same textbook cited with inconsistent chapter (15 vs 17) and page numbers (695 vs 784) within same evidence file — suggests page numbers approximate, not verified.

Assessment: Low — academic citation hygiene issue.

## Contradiction 7

Statement A: Report Limitations: "Performance benchmarks lack standardized methodology across vendors" (report.md:138) and Finding 7 / Evidence 9 classified as INTERPRETATION with "actual crossover depends on workload... No universal numeric threshold" (report.md:58, evidence.md:79).

Location: research/05-report.md:58, 138

Statement B: Same documents present specific thresholds (≤10%/>30%, 15-20%, 15% crossover) as if actionable guidance, with Evidence 15 giving "3-5 retries, 10-50ms start, 1-5s cap" as Microsoft-recommended.

Location: research/03-evidence.md:125-131, research/05-report.md:103-107

Type: INTERNAL

Impact: Tone contradiction — correctly hedges that thresholds are workload-dependent while simultaneously presenting precise numbers as anchored guidance. Both can't be "precise recommendation" and "no universal threshold" without clearer framing.

Assessment: Low — defensible hedging, but could mislead skim readers.

---
No contradictions found between research/01-plan.md risks and later mitigations — risks properly addressed except where noted above.
