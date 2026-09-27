# Research Audit Plan: Database Constraints

## Target Lab
`labs/27-database-constraints`

## PIPELINE OVERRIDE
- Research audit only.
- Do NOT audit implementation/code/tests/demo.
- Output location: `labs/27-database-constraints/research-audit/`

## Files Reviewed
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md`
- `research/05-report.md`
- `research/06-open-questions.md`

## Major Claims To Verify
1. Application-only validation (`exists() + insert`) is vulnerable to race conditions under concurrency.
2. `UNIQUE` constraints atomically prevent duplicate key inserts and raise SQLSTATE `23505` (`unique_violation`).
3. `NOT NULL` constraint is functionally equivalent to `CHECK (col IS NOT NULL)` but implemented more efficiently in PostgreSQL.
4. `FOREIGN KEY` constraints do not automatically create indexes on referencing columns in PostgreSQL.
5. `CHECK` constraints are row-scoped, assume immutability, and cannot reference other rows/tables or contain subqueries.
6. Partial unique indexes (`CREATE UNIQUE INDEX ... WHERE deleted_at IS NULL`) solve soft-delete uniqueness patterns.
7. Partial index planner usage requires strict predicate implication; parameterized query clauses do not match partial index predicates.
8. Database constraint violations return SQLSTATE class 23 codes and structured error fields, which applications should map to domain/HTTP errors (e.g. HTTP 409).
9. Idempotency using `UNIQUE(reference_number)` + transaction prevents duplicate processing in payment/webhook integration patterns.
10. Production constraint addition should use `ALTER TABLE ... ADD CONSTRAINT ... NOT VALID` followed by `VALIDATE CONSTRAINT`.
11. Partitioned tables require all partition key columns to be included in any unique or primary key constraint.

## External URLs To Verify
- `https://www.postgresql.org/docs/current/ddl-constraints.html`
- `https://www.postgresql.org/docs/current/indexes-unique.html`
- `https://www.postgresql.org/docs/current/indexes-partial.html`
- `https://www.postgresql.org/docs/current/errcodes-appendix.html`
- `https://www.postgresql.org/docs/current/applevel-consistency.html`
- `https://www.postgresql.org/docs/current/explicit-locking.html`
- `https://www.postgresql.org/docs/current/ddl-partitioning.html`
- `https://www.postgresql.org/docs/current/sql-altertable.html`
- `https://www.sqlite.org/lang_createtable.html`
- `https://docs.stripe.com/api/idempotent_requests`

## Primary Risks
- Inaccurate attribution of documentation quotes to specific PostgreSQL 18 doc sections.
- Over-generalization of PostgreSQL-specific behavior to all SQL databases.
- Over-reliance on secondary/community assumptions without official documentation backing.

## Audit Strategy
1. Cross-reference source citations in `02-sources.md` and `03-evidence.md` against official PostgreSQL/SQLite/Stripe docs.
2. Evaluate claim accuracy and severity in `03-claim-audit.md`.
3. Audit recorded contradictions, gaps, and open questions in `04-contradictions.md` and `06-gaps.md`.
4. Issue final verdict in `07-verdict.md`.
