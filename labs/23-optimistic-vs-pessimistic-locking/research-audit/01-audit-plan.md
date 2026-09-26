# Audit Plan — labs/23-optimistic-vs-pessimistic-locking

Target Lab: labs/23-optimistic-vs-pessimistic-locking
Audit Date: 2026-09-25
Scope: RESEARCH ONLY per pipeline override. No code audit, no implementation changes.

## Files Reviewed
- research/01-plan.md (15 questions, search strategy)
- research/02-sources.md (12 sources claimed)
- research/03-evidence.md (20 evidence items)
- research/05-report.md (16 findings + appendices)
- research/06-open-questions.md (gaps disclosed)
- research-revision/01-revision-plan.md, 02-changes-made.md, 03-revision-result.md

## Claims To Verify
1. Formal OCC/PCC definitions attributed to Bernstein (1987) and Silberschatz (2019) with page numbers.
2. MVCC/FOR UPDATE/deadlock/isolation quotes attributed to PostgreSQL docs sections 13.1–13.4.
3. `SELECT FOR UPDATE` vendor syntax matrix (PG / SQL Server / Oracle / MySQL).
4. Numeric claims: contention crossover ≤10% / >30% / 15–20%; retry 3–5x, 10–50ms start, 1–5s cap.
5. ORM claims: Hibernate `@Version` + "exponential backoff" quote; EF Core `[Timestamp]`; Laravel `lockForUpdate()/sharedLock()` + `$timestamps` as optimistic locking.
6. Distributed DB claims (CockroachDB HLC, Spanner TrueTime).
7. Version claims (PG 14 NOWAIT/SKIP LOCKED, PG 18/19).

## Code To Execute
N/A — research-only audit, no implementation in lab.

## Primary Risks
- Dead URLs cited as supporting evidence (Bernstein PDF, Oracle SELECT, EF Core page).
- Numeric/benchmark quotes attributed to vendor docs that contain no such numbers.
- Page-level academic quotes that cannot be verified from cited URLs.
- MySQL behavior inferred without a source (honestly disclosed — lower risk).
- Internal numeric inconsistency (≤10%/>30% vs 15–20%).

## Audit Strategy
Fetch every URL source directly. Compare claimed quotes against actual page content. Check each numeric claim for anchoring. Cross-check evidence vs report vs revision-result for contradictions. Record gaps with severity; no silent repair.
