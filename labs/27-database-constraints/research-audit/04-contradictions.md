# Contradictions & Conflicts Audit: Database Constraints

## Audit Summary
No material contradictions found across the research documents (`01-plan.md`, `02-sources.md`, `03-evidence.md`, `04-contradictions.md`, `05-report.md`, `06-open-questions.md`).

## Specific Checks

### 1. NULL Semantics in UNIQUE Constraints
- **Statement A (`02-sources.md`, `03-evidence.md`):** PostgreSQL and SQLite treat NULLs in UNIQUE constraints as distinct by default (multiple NULL values permitted).
- **Statement B (`04-contradictions.md`, `06-open-questions.md`):** Multi-column combinations with mixed NULLs (e.g., `(1, NULL)` vs `(1, NULL)`) require testing because exact multi-column NULL equality rules are subtle.
- **Type:** CLARIFICATION / EDGE CASE
- **Impact:** None on core claims.
- **Assessment:** Consistent. The research correctly distinguishes standard single-column NULL behavior from open multi-column edge cases.

### 2. Application Idempotency vs Database Idempotency
- **Statement A (`02-sources.md`, `05-report.md`):** Stripe idempotency caches response payloads per idempotency key at the application/API layer.
- **Statement B (`05-report.md`, `03-evidence.md`):** Database idempotency uses `UNIQUE(reference_number)` to atomically reject duplicate records at the storage layer.
- **Type:** ARCHITECTURAL DISTINCTION
- **Impact:** None.
- **Assessment:** Consistent. The research clearly articulates that the two mechanisms are complementary layers in defense-in-depth, not contradictory implementations.

### 3. PostgreSQL vs SQLite Behavior
- **Statement A:** PostgreSQL auto-indexes UNIQUE constraints with B-trees and supports `NULLS NOT DISTINCT`.
- **Statement B:** SQLite implements UNIQUE via internal unique indexes and treats NULLs as distinct.
- **Type:** CROSS-DATABASE VALIDATION
- **Impact:** None.
- **Assessment:** Consistent. Cross-check against SQLite corroborates the fundamental storage-layer enforcement thesis.

## Conclusion
No internal contradictions or source conflicts found.
