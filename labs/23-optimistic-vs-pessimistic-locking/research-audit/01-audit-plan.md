# Audit Plan: Optimistic vs Pessimistic Locking Research

## Target Lab
`labs/23-optimistic-vs-pessimistic-locking`

## Scope
Pipeline Override: Research Only. Implementation/code audit is out of scope for this stage and handled under code audit placeholder rules.

## Files Reviewed
- `labs/23-optimistic-vs-pessimistic-locking/research/01-plan.md`
- `labs/23-optimistic-vs-pessimistic-locking/research/02-sources.md`
- `labs/23-optimistic-vs-pessimistic-locking/research/03-evidence.md`
- `labs/23-optimistic-vs-pessimistic-locking/research/04-contradictions.md`
- `labs/23-optimistic-vs-pessimistic-locking/research/05-report.md`
- `labs/23-optimistic-vs-pessimistic-locking/research/06-open-questions.md`

## Claims To Verify
1. Lost update occurs under default isolation levels (READ COMMITTED in PostgreSQL/Oracle) during concurrent read-modify-write cycles.
2. `SELECT ... FOR UPDATE` acquires row-level exclusive locks blocking concurrent writers until commit/rollback.
3. Pessimistic locking induces deadlocks and concurrency bottlenecks when transactions are long-running.
4. Optimistic locking validates version/timestamp guards at commit, detecting conflicts via zero affected rows without holding database locks.
5. Selection heuristic: Pessimistic prevents conflicts (high contention/critical data); Optimistic detects conflicts (low contention/read-heavy/long-lived transactions).
6. Single-statement atomic conditional updates (`UPDATE ... SET stock = stock - N WHERE stock >= N`) remove read-modify-write race windows.
7. Plain database transactions alone without explicit locking or conditional guards do not prevent lost updates.
8. Default isolation levels differ significantly across vendors (PostgreSQL/Oracle: READ COMMITTED; MySQL: REPEATABLE READ).
9. Common anti-patterns: premature distributed locks (Redis), holding pessimistic locks across network calls, ignoring 0-rows-affected.

## Code To Execute
None. Per pipeline override instructions:
`PIPELINE OVERRIDE: Audit research only. Do not audit implementation/code in this stage.`

## Primary Risks
1. Verification gaps for vendor-specific documentation (e.g. dev.mysql.com returning HTTP 403, requiring verification via Oracle CDN mirror).
2. Unverified performance quantification claims (e.g. concrete latency/throughput numbers vs qualitative assertions).
3. Over-generalization of distributed lock trade-offs (asserting Redis is an anti-pattern without explicit boundary conditions).
4. Version counter integer overflow risks in optimistic locking implementations.

## Audit Strategy
1. Cross-reference all 16 listed sources in `research/02-sources.md` against claims in `research/03-evidence.md` and `research/05-report.md`.
2. Verify reachability, relevance, and fidelity of cited citations against primary documentation.
3. Validate claim classification (FACT vs INTERPRETATION vs IMPLEMENTATION-SPECIFIC) and ensure unsupported claims are flagged.
4. Examine internal and cross-database contradictions identified in `research/04-contradictions.md`.
5. Map documented research gaps and unresolved open questions in `research/06-open-questions.md`.
6. Issue formal verdict in `research-audit/07-verdict.md`.
