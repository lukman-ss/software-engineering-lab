# 01 - Audit Plan

## Target Lab
`labs/13-backward-compatibility`

## Files Reviewed
- `README.md`
- `research/01-research-plan.md`
- `research/02-sources.md`
- `research/03-core-concepts.md`
- `research/04-database-migration.md`
- `research/05-api-compatibility.md`
- `research/06-expand-migrate-contract.md`
- `research/07-deployment-and-rollback.md`
- `research/08-failure-modes.md`
- `research/09-case-studies.md`
- `research/10-open-questions.md`
- `research/11-final-research.md`

## Claims To Verify
1. Expand-Migrate-Contract (Parallel Change) guarantees zero-downtime evolution.
2. Dropping columns requires a 3-release migration strategy (Ignore -> Drop -> Remove Ignore) across releases M, M+1, M+2.
3. Additive changes are always backward-compatible (Google AIP-180 & Fowler).
4. Breaking changes classification: Wire, Semantic, and Source incompatibility.
5. Backfill must be batched and idempotent (`ON CONFLICT DO NOTHING`) to avoid lock contention.
6. Rolling deployment requires N and N+1 compatibility with a shared schema.
7. Dual write introduces partial inconsistency risks unless transactional/atomic.

## Code To Execute
- None. No source code or tests currently implemented in `labs/13-backward-compatibility/`.

## Primary Risks
- Generalizing framework-specific mechanics (e.g. Rails ActiveRecord's `ignore_column` and schema cache) as universal database truths.
- Stating additive changes are always safe without qualifying enum/default value edge cases.
- Missing validation on distributed transactions vs single-database dual writes.

## Audit Strategy
1. Inspect claimed URLs and verify accessibility, publisher, content match, and scope.
2. Check claim fidelity against primary documentation (Martin Fowler, Stripe, GitLab, Google AIP-180).
3. Evaluate whether recommendations distinguish platform-specific behavior (PostgreSQL/Rails) from universal patterns.
4. Record internal contradictions and research gaps.
5. Issue verdict.
