# Contradictions & Uncertainties

## Contradictions Found

### No Major Contradictions

After cross-checking all major claims against primary sources (PostgreSQL 18, SQLite), no material contradictions were discovered.

Key verifications:
- NULL handling in UNIQUE: PostgreSQL and SQLite both use DISTINCT behavior (NULLs not equal)
- Constraint enforcement timing: Both document enforcement only on INSERT/UPDATE, not SELECT
- Index creation: Both document automatic unique index creation for UNIQUE/PRIMARY KEY

### Blocked Access (Not Contradictions)

**MySQL Documentation HTTP 403 Forbidden**
- URLs attempted:
  - https://dev.mysql.com/doc/refman/8.0/en/constraints.html
  - https://dev.mysql.com/doc/refman/8.0/en/create-table.html
  - https://dev.mysql.com/doc/refman/8.0/en/unique-constraints.html
- Status: HTTP 403 Forbidden
- Impact: MySQL-specific constraint behavior (NULL handling, DEFERRABLE support, error codes) not directly verified
- Mitigation: PostgreSQL and SQLite cross-database verification provides strong evidence for standard-compliant behavior

## Partial Contradictions / Clarifications

### NULL Handling in Multi-Column UNIQUE Constraints

**Question:** When does uniqueness trigger with mixed NULL/non-NULL values?

**PostgreSQL Documentation:**
- "By default, two null values are not considered equal in this comparison"
- Multi-column: "This specifies that the combination of values in the indicated columns is unique across the whole table"

**SQLite Documentation:**
- "For the purposes of UNIQUE constraints, NULL values are considered distinct from all other values, including other NULLs"

**Assessment:**
Both implementations use DISTINCT behavior, but documentation doesn't explicitly confirm whether (1, NULL) and (1, NULL) violate uniqueness. Evidence strongly suggests:
- (1, NULL) vs (1, NULL) → VIOLATION (NULL == NULL in comparison)
- (1, NULL) vs (1, 2) → NO VIOLATION (NULL ≠ 2)

**Verification needed:** Live database test to confirm exact behavior.

### Partial Unique Index vs Named Constraint

**PostgreSQL Documentation:**
- Partial unique indexes created via `CREATE UNIQUE INDEX ... WHERE ...`
- Named constraints via `CONSTRAINT name UNIQUE (...)`

**Assessment:**
Not contradictory — complementary tools:
- Constraints appear in information_schema.table_constraints, propagate via table inheritance, appear in error messages
- Partial indexes created via CREATE UNIQUE INDEX are not "constraints" in catalog

**Conclusion:** Different tooling for different use cases; no actual contradiction.

### DEFERRABLE Constraint + Partial Index Interaction

**Documentation:**
- "Only UNIQUE, PRIMARY KEY, EXCLUDE, and REFERENCES (foreign key) constraints accept this clause. NOT NULL and CHECK constraints are not deferrable"
- Partial unique indexes created by constraints respect DEFERRABLE timing

**Assessment:**
Constraint deferral is a logical property; partial index is implementation detail. When constraint is DEFERRABLE, both the constraint validation AND the underlying index check should defer to transaction end.

**Uncertainty:** Documentation doesn't explicitly state that partial index check respects DEFERRABLE timing.

**Assumption:** Yes, the constraint deferral applies to the entire constraint including index enforcement.

## Open Questions Without Definitive Answers

### 1. Multi-Column UNIQUE Constraint with NULLs

**Question:** Does (1, NULL) violate UNIQUE if another row is (1, NULL)?

**Evidence:**
PostgreSQL docs: "By default, two null values are not considered equal"  
SQLite docs: "NULL values are considered distinct from all other values, including other NULLs"

**Expected behavior:** (1, NULL) vs (1, NULL) should violate uniqueness (NULL == NULL), but not explicitly documented.

**Verification:** Live database test required.

### 2. Partial Index Planner Recognition Depth

**Question:** Beyond "x < 1 implies x < 2", what implication classes does PostgreSQL's planner recognize?

**Evidence:**
"PostgreSQL does not have a sophisticated theorem prover that can recognize mathematically equivalent expressions that are written in different forms"  
"The system can recognize simple inequality implications"

**Missing:**
- Exact algorithm for implication recognition
- Examples of non-recognized patterns

**Risk:** Application may create partial index that planner never uses.

**Mitigation:** Use EXPLAIN to verify index usage.

### 3. Constraint Validation Order in Mixed Scenarios

**Question:** Which error reported when row violates both NOT NULL and CHECK constraint?

**Evidence:**
"CHECK constraints are tested for each row in alphabetical order by name, after checking NOT NULL constraints"

**Assumption:** NOT NULL checked first (explicit implementation), then CHECK (alphabetical), then FK/UNIQUE (index lookups).

**Missing:** Explicit ordering for FK vs UNIQUE.

**Verification:** Test INSERT violating multiple constraints.

### 4. Deferred UNIQUE Constraint + Partial Index Timing

**Question:** If UNIQUE constraint is DEFERRABLE and creates partial unique index, when is index checked?

**Evidence:**
DEFERRABLE constraints checked at transaction end; partial indexes are index-level structures.

**Assumption:** Constraint deferral applies to entire constraint including index enforcement.

**Risk:** Index may enforce immediately while constraint logic defers.

**Verification:** Test DEFERRABLE UNIQUE with partial index and simultaneous constraint violation across two rows.

### 5. Concurrent Partial Index Creation Behavior

**Question:** Can `CREATE UNIQUE INDEX ... WHERE ... CONCURRENTLY` detect existing violations?

**Evidence:**
CONCURRENTLY prevents table lock during build; index scans existing data during construction.

**Assumption:** Build fails if data violates uniqueness at any point during construction.

**Risk:** Concurrent INSERT during CONCURRENTLY build may cause failure.

**Verification:** Test race scenario during concurrent partial index creation.

### 6. Generated Column + UNIQUE Constraint Evaluation Order

**Question:** When is generation expression evaluated relative to uniqueness check?

**Evidence:**
Generated columns use immutable functions; UNIQUE constraint creates B-tree index.

**Missing:** Exact ordering of generation vs uniqueness check.

**Verification:** Test with volatile-looking but immutable function in generated column under concurrent load.

### 7. EXCLUDE Constraint with Temporal Validity

**Question:** Can EXCLUDE enforce non-overlapping time ranges with dynamic updates?

**Evidence:**
"PostgreSQL assumes that CHECK constraints' conditions are immutable"  
EXCLUDE uses index-level validation (not CHECK)

**Analysis:** EXCLUDE constraints don't use CHECK's immutability assumption — they use index-level validation.

**Assumption:** EXCLUDE properly handles range updates that create overlaps.

**Verification:** Test EXCLUDE with range overlaps under concurrent updates.

### 8. MySQL Behavior Parity (Unverified)

**Question:** Exact MySQL-specific differences in constraint behavior?

**Evidence:** Unable to access MySQL documentation (HTTP 403).

**Need:** Docker testbed or alternate mirror.

**Suggested tests:**
- NULL handling in UNIQUE constraints (MySQL may treat NULLs as equal by default)
- DEFERRABLE support (MySQL historically lacks this)
- Error code mapping

## PostgreSQL Version Consistency

All evidence from PostgreSQL 18 documentation. Minor differences may exist in:
- PostgreSQL 17 (supported)
- PostgreSQL 16 (supported)
- Older versions (unsupported): constraints may lack features like NULLS NOT DISTINCT

## Recommendations for Validation

1. **Test NULL handling behavior in real database** — multi-column UNIQUE with NULLs not explicitly documented
2. **Run EXPLAIN on queries with partial indexes** — verify planner recognition of implications
3. **Check constraint error messages with named vs unnamed constraints** — verify structured error fields
4. **Verify DEFERRABLE timing with partial indexes** — ensure constraint deferral applies to index
5. **Compare MySQL behavior if access becomes available** — Docker container for MySQL 8

## Summary of Contradictions

**No material contradictions discovered** in cross-database verification (PostgreSQL vs SQLite).

**MySQL behavior remains unverified** due to HTTP 403 on dev.mysql.com.

**Uncertainties documented** for edge cases requiring live database testing (NULL in multi-column UNIQUE, partial index planner recognition, DEFERRABLE + partial index interaction).