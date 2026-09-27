# Audit Plan: Optimistic vs Pessimistic Locking (Research Stage)

## Target Lab
`labs/23-optimistic-vs-pessimistic-locking`

## Files Reviewed
- `labs/23-optimistic-vs-pessimistic-locking/research/01-plan.md`
- `labs/23-optimistic-vs-pessimistic-locking/research/02-sources.md`
- `labs/23-optimistic-vs-pessimistic-locking/research/03-evidence.md`
- `labs/23-optimistic-vs-pessimistic-locking/research/04-contradictions.md`
- `labs/23-optimistic-vs-pessimistic-locking/research/05-report.md`
- `labs/23-optimistic-vs-pessimistic-locking/research/06-open-questions.md`

## Claims To Verify
1. Lost update definition and non-serializability under concurrent execution.
2. Read Committed isolation level in PostgreSQL and Oracle does NOT prevent lost updates.
3. Pessimistic locking semantics via `SELECT ... FOR UPDATE` (row-level exclusive lock, blocking writers).
4. Deadlock risks and mitigation via lock ordering in pessimistic locking.
5. Optimistic locking mechanisms (version checking at commit time, Kung & Robinson 1981, 0-rows-affected conflict detection).
6. Application-level vs database-level optimistic concurrency patterns (Martin Fowler, EF Core).
7. Database isolation level variations: PostgreSQL REPEATABLE READ (Snapshot Isolation) vs Oracle SERIALIZABLE (ORA-08177).
8. Atomic `UPDATE ... WHERE condition` as an alternative to explicit locking.
9. MVCC as the underlying architecture enabling non-blocking reads.

## Code To Execute
*Pipeline override*: Research audit only. No code or implementation execution required in this stage.

## Primary Risks
1. **Unverified Primary Source**: MySQL 8.0 documentation was inaccessible (HTTP 403) during research; claims regarding MySQL locking reads and gap locks rely on secondary attribution.
2. **Overgeneralization of "Atomic UPDATE is always safe"**: Atomic updates avoid lost updates for single-row updates, but predicate/range integrity or multi-entity invariants still require explicit locking or higher isolation.
3. **Integer Counter Overflow**: Practical limitations of integer version columns in optimistic locking are flagged as open questions but need clear qualification.
4. **Distinction between Database OCC and Application OCC**: Clarifying differences between engine-level OCC (e.g., Kung-Robinson validate phase, MVCC snapshot serializability) and application-level OCC (version column with `UPDATE ... WHERE version = ?`).

## Audit Strategy
1. Independently fetch and verify all 9 cited URLs from primary and secondary sources.
2. Cross-examine claims in `03-evidence.md` and `05-report.md` against extracted source text.
3. Assess factual accuracy, severity of gaps, and validity of conclusions.
4. Generate the full suite of research audit markdown reports (`01-audit-plan.md` to `07-verdict.md`).
