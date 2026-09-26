# Audit Plan: Optimistic vs Pessimistic Locking Research

## Target Lab
`labs/23-optimistic-vs-pessimistic-locking`

## Pipeline Mode
- Research-only audit (per pipeline override instructions).
- Implementation / runnable code audit deferred to implementation phase.

## Files Reviewed
- `labs/23-optimistic-vs-pessimistic-locking/research/01-plan.md`
- `labs/23-optimistic-vs-pessimistic-locking/research/02-sources.md`
- `labs/23-optimistic-vs-pessimistic-locking/research/03-evidence.md`
- `labs/23-optimistic-vs-pessimistic-locking/research/04-contradictions.md`
- `labs/23-optimistic-vs-pessimistic-locking/research/05-report.md`
- `labs/23-optimistic-vs-pessimistic-locking/research/06-open-questions.md`

## Claims To Verify
1. Definition and reproducibility of "lost update" under READ COMMITTED isolation.
2. Pessimistic locking semantics via `SELECT ... FOR UPDATE` (PostgreSQL, MySQL, Oracle).
3. Deadlock risks, lock wait duration, and holding locks across slow/external calls as an anti-pattern.
4. Optimistic locking mechanics via version guard in `WHERE` clause and checking affected rows.
5. Selection criteria (pessimistic prevents vs optimistic detects; high contention vs low contention).
6. Atomic single-statement updates (`UPDATE ... SET stock = stock - N WHERE stock >= N`) eliminating read-modify-write windows.
7. Isolation level vs transaction misconception (transaction boundary alone does not prevent lost update).
8. Default isolation levels and multi-database behavioral divergence.

## Primary Risks
- Overgeneralizing database-specific semantics across vendors (e.g., PostgreSQL Snapshot Isolation vs MySQL 2PL next-key locking).
- Relying on secondary/unverified sources for edge cases (e.g. Hibernate versionless locking, Redis distributed lock boundaries).
- Arbitrary numeric claims or unmeasured benchmark generalizations.

## Audit Strategy
1. Validate all 14 cited sources for reachability, tier categorization, and substantive topical support.
2. Audit all 9 major findings and 16 evidence items against source citations and claim scope.
3. Review documented contradictions to verify they reflect actual engine behaviors without ungrounded assertions.
4. Check gap analysis for completeness in capturing unverified areas (e.g., benchmark numbers, atomic decrement vendor recipes).
5. Produce evidence-based verdict on research soundness.
