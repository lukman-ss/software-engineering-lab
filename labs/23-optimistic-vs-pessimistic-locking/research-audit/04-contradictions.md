# Contradictions Audit: Optimistic vs Pessimistic Locking

## Contradiction 1

Statement A: PostgreSQL implements REPEATABLE READ as Snapshot Isolation (preventing both non-repeatable reads and lost updates, throwing serialization errors on conflict).  
Location: `research/04-contradictions.md:9-11`

Statement B: MySQL InnoDB implements REPEATABLE READ using gap/next-key locks with strict 2PL properties.  
Location: `research/04-contradictions.md:12-14`

Statement C: ANSI SQL standard defines REPEATABLE READ as preventing non-repeatable reads while allowing phantom reads.  
Location: `research/04-contradictions.md:15-16`

Type: INTERNAL / SOURCE_CONFLICT (Engine Semantic Divergence)

Impact: MEDIUM — Developers cannot assume "REPEATABLE READ" behaves identically across PostgreSQL, MySQL, and Oracle.

Assessment: Correctly identified and documented in research. Not a flaw in the research, but an important semantic divergence across implementations.

---

## Contradiction 2

Statement A: PostgreSQL treats `READ UNCOMMITTED` as identical to `READ COMMITTED` because MVCC does not support dirty reads.  
Location: `research/04-contradictions.md:24-25`

Statement B: SQL Server / traditional non-MVCC databases permit true dirty reads under `READ UNCOMMITTED`.  
Location: `research/04-contradictions.md:27-28`

Type: INTERNAL / ENGINE_SPECIFIC

Impact: LOW — Highlights engine-specific behavior.

Assessment: Correctly documented in research.

---

## Summary of Audit Findings

No material unhandled contradictions found in the research files. All semantic variations between engines (PostgreSQL, MySQL, Oracle) and specifications (ANSI SQL) have been explicitly identified, categorized, and resolved in `research/04-contradictions.md` and `research/05-report.md`.
