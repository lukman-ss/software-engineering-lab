# Open Questions

## Technical Questions Requiring Testing/Research

### 1. NULL Semantics in Multi-Column UNIQUE Constraints
**Question:** With a multi-column UNIQUE (a, b), and the default NULL handling (NULLs are distinct), does the combination `(1, NULL)` violate the constraint if another row also has `(1, NULL)`?

**Analysis from docs:** "By default, two null values are not considered equal in this comparison" — suggesting NULL != NULL. But whether the index treats (1, NULL) and (1, NULL) as the same entry is the question.

**Suggested verification:** Live database test with INSERT of two rows with identical non-null + null combination.

---

### 2. Partial Index Planner Recognition Capability
**Question:** Beyond "x < 1 implies x < 2", what mathematical implication classes does PostgreSQL's planner recognize when matching query predicates to partial index predicates?

**Analysis from docs:** "PostgreSQL does not have a sophisticated theorem prover that can recognize mathematically equivalent expressions that are written in different forms."

**Suggested verification:** 
- Test `WHERE status IN ('a', 'b')` index vs `WHERE status = 'a' OR status = 'b'` query
- Test `WHERE status = 'active' AND tenant_id = 123` index vs filtered query
- Test expression-based predicates

---

### 3. Constraint Validation Order in Mixed Scenarios
**Question:** When a single INSERT violates both a NOT NULL constraint and a CHECK constraint, which error is reported?

**Analysis from docs:** CHECK constraints are "tested for each row in alphabetical order by name, after checking NOT NULL constraints." This suggests NOT NULL first. But no explicit ordering stated for FOREIGN KEY and UNIQUE.

**Suggested verification:** Create table with all constraint types, insert row violating multiple, observe which error surfaces.

---

### 4. Deferred UNIQUE Constraint + Partial Index Interaction
**Question:** If a UNIQUE constraint is declared DEFERRABLE and creates a partial unique index (via underlying index predicate), at what point is the index checked — row-by-row or at transaction end?

**Analysis from docs:** DEFERRABLE constraints can be checked at transaction end. Partial indexes are index-level structures. The interaction is unclear — the constraint is deferred but the index enforces immediately.

**Suggested verification:** Test DEFERRABLE UNIQUE with simultaneous constraint violation across two rows in same transaction.

---

### 5. Concurrent Partial Index Creation
**Question:** Can `CREATE UNIQUE INDEX ... WHERE ... CONCURRENTLY` detect and reject existing data that violates the partial uniqueness?

**Analysis from docs:** CONCURRENTLY prevents table lock during build. The index scans existing data during construction. If data violates uniqueness, the build must fail.

**Risk:** During CONCURRENTLY build, a row satisfying the predicate may be inserted by another transaction — does the index creation catch it?

**Suggested verification:** Read PostgreSQL source / commitfest notes on concurrent partial index behavior; test race scenario.

---

### 6. Generated Column Evaluation and Constraint Interaction
**Question:** If a generated column (STORED or VIRTUAL) has a UNIQUE constraint, when is the generation expression evaluated relative to the uniqueness check?

**Analysis from docs:** Generated columns: "Any functions and operators used must be immutable." UNIQUE constraint creates a B-tree index.

**Suggested verification:** Test with a volatile-looking but immutable function in a generated column under concurrent load.

---

### 7. Foreign Key Action + Trigger Ordering Under Concurrency
**Question:** For a foreign key with `ON DELETE CASCADE`, if a trigger on the child table also fires, what determines the order of cascade action vs trigger execution?

**Analysis from docs:** FK cascades "are executed as part of the data changing command." Triggers fire "before/after" row level.

**Suggested verification:** Document behavior with FK CASCADE + BEFORE/AFTER triggers on same table — which fires first?

---

### 8. EXCLUDE Constraint with Temporal (Time Range) Validity
**Question:** Can an EXCLUDE constraint enforce that "no two rows have the same key with overlapping time ranges" where the time range is itself dynamic (updated at runtime)?

**Analysis from docs:** EXCLUDE USING gist supports range-overlap operators. But "PostgreSQL assumes that CHECK constraints' conditions are immutable."

**Suggested verification:** EXCLUDE constraints don't use CHECK's immutability assumption — they use index-level validation. Test whether EXCLUDE properly handles range updates that create overlaps.

---

### 9. MySQL Behavior Parity (Unverified)
**Question:** What are the exact MySQL-specific differences in constraint behavior (NULL handling, DEFERRABLE support, error codes)?

**Analysis from docs:** Unable to access MySQL documentation (HTTP 403). Need Docker testbed or alternate mirror.

**Suggested verification:** 
- Run MySQL 8 in Docker, test NULL in UNIQUE constraints
- Test DEFERRABLE support (MySQL historically lacks this)
- Map MySQL error codes to PostgreSQL equivalents

---

### 10. Constraint Violation Error Detail Fields
**Question:** Beyond SQLSTATE, what structured fields does PostgreSQL return for constraint violations (table name, column name, constraint name)?

**Analysis from docs:** "such names are supplied in separate fields of the error report message" — mentions table/column/constraint.

**Suggested verification:** Use a PostgreSQL driver to capture all error detail fields from different constraint violation types.

---

## Research Environment Notes

- **PostgreSQL version tested against (documentation):** PostgreSQL 18 (current), supported back to 14
- **MySQL version to test against:** MySQL 8.0 (documentation blocked)
- **Docker available:** Yes (for local verification of open questions)

## Priority for Follow-Up

| Priority | Question | Effort | Impact |
|---|---|---|---|
| High | Q1: NULL in multi-column UNIQUE | Low | High - common bug source |
| High | Q4: Deferred + partial index | Medium | High - transaction correctness |
| Medium | Q2: Partial index planner | Medium | Medium - performance tuning |
| Medium | Q5: Concurrent partial index | High | High - deployment safety |
| Low | Q3: Constraint validation order | Low | Low - academic interest |
| Low | Q10: Error detail fields | Low | Medium - error handling