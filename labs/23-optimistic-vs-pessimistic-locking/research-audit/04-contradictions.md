# Contradiction Audit: Optimistic vs Pessimistic Locking

## Summary
No material contradictions found across Tier 1 and Tier 2 sources.

The research maintains internal consistency across definitions, database behaviors, and trade-offs.

## Noted Implementation Nuances (Non-Contradictory)

### Nuance 1: Isolation Level Semantics Across Vendors

**Statement A (PostgreSQL REPEATABLE READ)**:
Implemented via Snapshot Isolation; prevents phantom reads; raises `could not serialize access due to concurrent update` on concurrent row modification (`research/03-evidence.md:59-62`).

**Statement B (Oracle Isolation Levels)**:
Does not offer an ANSI REPEATABLE READ level; implements READ COMMITTED, SERIALIZABLE (which raises `ORA-08177`), and READ ONLY (`research/03-evidence.md:68-71`).

**Type**: INTERNAL / SYSTEM_SPECIFIC_VARIATION

**Impact**: None. Accurately highlights that standard SQL isolation names have vendor-specific behavioral implementations.

**Assessment**: PASS

---

### Nuance 2: Concurrency Enforcement Location (Database Engine vs Application Code)

**Statement A (Database-level Concurrency)**:
PostgreSQL SSI / Repeatable Read and Oracle Serializable track read/write dependencies at the storage/transaction engine level and abort on serialization conflicts (`research/05-report.md:71-78`).

**Statement B (Application-level Concurrency)**:
Optimistic locking via version column (`UPDATE ... WHERE version = ?`) is executed at the application layer by checking affected rows (`research/05-report.md:51-64`).

**Type**: INTERNAL / SYSTEM_SPECIFIC_VARIATION

**Impact**: None. Both mechanisms represent distinct architectural layers for achieving concurrency control.

**Assessment**: PASS
