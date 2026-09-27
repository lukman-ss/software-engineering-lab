# Open Questions

## Technical Questions Requiring Testing/Research

### 1. Multi-Column UNIQUE Constraint NULL Semantics

**Question:** With a multi-column UNIQUE constraint on `(user_id, email)`, and the default NULL handling (NULLs are distinct), does `(7, NULL)` violate uniqueness if another row is `(7, NULL)`?

**Analysis from docs:** PostgreSQL: "By default, two null values are not considered equal in this comparison"  
SQLite: "NULL values are considered distinct from all other values, including other NULLs"

**Missing evidence:** Documentation doesn't explicitly state whether two identical rows with non-NULL + NULL values violate the constraint.

**Expected:** `(7, NULL)` vs `(7, NULL)` should violate (NULL == NULL in comparison), but not verified.

**Suggested verification:** 
```sql
CREATE TABLE users (user_id INT, email TEXT, UNIQUE (user_id, email));
INSERT INTO users VALUES (7, NULL);
INSERT INTO users VALUES (7, NULL);  -- Does this violate?
```

**Priority:** High  
**Impact:** Common bug source for multi-column unique constraints  
**Effort:** Low (single test query)

---

### 2. Partial Index Planner Recognition Depth

**Question:** Beyond "x < 1 implies x < 2", what mathematical implication classes does PostgreSQL's planner recognize when matching query predicates to partial index predicates?

**Analysis from docs:** "PostgreSQL does not have a sophisticated theorem prover that can recognize mathematically equivalent expressions that are written in different forms"  
"The system can recognize simple inequality implications"

**Missing evidence:** 
- Exact algorithm for implication recognition
- Comprehensive list of recognized patterns
- Examples of non-recognized patterns

**Test scenarios needed:**
```sql
-- Can planner recognize?
-- WHERE status IN ('a', 'b') vs index WHERE status = 'a' OR status = 'b'
-- WHERE status = 'active' AND tenant_id = 123 vs index WHERE status = 'active'
-- WHERE created_at >= '2024-01-01' vs index WHERE created_at >= '2023-12-01'
```

**Suggested verification:** Test with EXPLAIN on partial index with various query predicates.

**Priority:** Medium  
**Impact:** Performance tuning - planner may not use partial index  
**Effort:** Medium (multiple test queries with EXPLAIN)

---

### 3. Constraint Validation Order in Mixed Violations

**Question:** When a single INSERT violates both a NOT NULL constraint and a CHECK constraint, which error is reported?

**Analysis from docs:** "CHECK constraints are tested for each row in alphabetical order by name, after checking NOT NULL constraints"

**Missing evidence:** No explicit ordering stated for FOREIGN KEY vs UNIQUE validation.

**Expected:** NOT NULL first (explicit implementation), then CHECK (alphabetical), then FK/UNIQUE (index lookups), but not verified.

**Suggested verification:** 
```sql
CREATE TABLE test (
  col1 INT NOT NULL,
  col2 INT CHECK (col2 > 0),
  col3 INT CHECK (col3 < 100),
  UNIQUE (col1)
);
INSERT INTO test VALUES (NULL, -1, 200);  -- Which error surfaces?
```

**Priority:** Low  
**Impact:** Minimal (deterministic behavior in practice)  
**Effort:** Low (single test query)

---

### 4. DEFERRABLE UNIQUE Constraint + Partial Index Interaction

**Question:** If a UNIQUE constraint is declared DEFERRABLE and creates a partial unique index, at what point is the index checked — row-by-row or at transaction end?

**Analysis from docs:** DEFERRABLE constraints can be checked at transaction end  
Partial indexes are index-level structures  
Documentation: "Only UNIQUE, PRIMARY KEY, EXCLUDE, and REFERENCES constraints accept this clause"

**Missing evidence:** Documentation doesn't explicitly state that partial index check respects DEFERRABLE timing.

**Risk:** Index may enforce immediately while constraint logic defers, creating inconsistency.

**Expected:** Constraint deferral applies to entire constraint including index enforcement.

**Suggested verification:**
```sql
CREATE TABLE test (
  user_id INT,
  email TEXT,
  UNIQUE (user_id, email) DEFERRABLE INITIALLY DEFERRED
);
CREATE UNIQUE INDEX test_email_idx ON test (email) WHERE deleted_at IS NULL;
-- Test: insert two rows violating partial uniqueness, then commit
```

**Priority:** High  
**Impact:** High - transaction correctness  
**Effort:** Medium (multiple test queries)

---

### 5. Concurrent Partial Index Creation Safety

**Question:** Can `CREATE UNIQUE INDEX ... WHERE ... CONCURRENTLY` detect and reject existing data that violates the partial uniqueness?

**Analysis from docs:** CONCURRENTLY prevents table lock during build  
Index scans existing data during construction  
Documentation: "Concurrently" builds avoid table locks

**Missing evidence:** Documentation doesn't explicitly address behavior when existing data violates the partial uniqueness.

**Risk:** If data violates uniqueness during CONCURRENTLY build, the build must fail. But if concurrent INSERT occurs during build, does the index creation catch it?

**Expected:** Build fails if data violates uniqueness; concurrent INSERT during build may cause failure.

**Suggested verification:** Create table with existing duplicate partial uniqueness violations, attempt CONCURRENTLY build.

**Priority:** High  
**Impact:** High - deployment safety  
**Effort:** High (requires controlled test environment)

---

### 6. Generated Column + UNIQUE Constraint Interaction

**Question:** If a generated column (STORED) has a UNIQUE constraint, when is the generation expression evaluated relative to the uniqueness check?

**Analysis from docs:** Generated columns use immutable functions  
UNIQUE constraint creates B-tree index  
Documentation: "Any functions and operators used must be immutable"

**Missing evidence:** Exact ordering of generation vs uniqueness check.

**Expected:** Generation evaluated during INSERT, uniqueness checked during INSERT, atomic operation.

**Suggested verification:**
```sql
CREATE TABLE test (
  id INT,
  col1 TEXT,
  col2 TEXT GENERATED ALWAYS AS (col1 || '_suffix') STORED,
  UNIQUE (col2)
);
INSERT INTO test (id, col1) VALUES (1, 'a'), (2, 'a');  -- Does second fail?
```

**Priority:** Medium  
**Impact:** Correctness for generated columns  
**Effort:** Medium (test query)

---

### 7. EXCLUDE Constraint with Temporal Range Overlap

**Question:** Can an EXCLUDE constraint enforce that "no two rows have the same key with overlapping time ranges" where the time range is updated at runtime?

**Analysis from docs:** EXCLUDE uses GiST index with range operators  
"PostgreSQL assumes that CHECK constraints' conditions are immutable"

**Missing evidence:** Documentation doesn't explicitly test EXCLUDE with runtime-updated ranges.

**Expected:** EXCLUDE properly handles range updates that create overlaps because it uses index-level validation, not CHECK's immutability assumption.

**Suggested verification:**
```sql
CREATE TABLE bookings (
  room_id INT,
  starts_at TIMESTAMP,
  ends_at TIMESTAMP,
  EXCLUDE USING gist (room_id WITH =, tsrange(starts_at, ends_at) WITH &&)
);
-- Update one booking to overlap another
```

**Priority:** Medium  
**Impact:** High for booking/event systems  
**Effort:** Medium (test with EXCLUDE)

---

### 8. MySQL Constraint Behavior Parity (Unverified)

**Question:** What are the exact MySQL-specific differences in constraint behavior (NULL handling, DEFERRABLE support, error codes)?

**Analysis from docs:** Unable to access MySQL documentation (HTTP 403).

**Missing evidence:** Cannot verify MySQL-specific behavior.

**Suggested verification:**
- Run MySQL 8 in Docker
- Test NULL handling in UNIQUE constraints (MySQL may treat NULLs as equal)
- Test DEFERRABLE support (MySQL historically lacks this)
- Map MySQL error codes to PostgreSQL equivalents

**Priority:** Medium  
**Impact:** Multi-database compatibility  
**Effort:** Medium (Docker setup + tests)

---

### 9. Constraint Violation Error Detail Fields

**Question:** Beyond SQLSTATE, what structured fields does PostgreSQL return for constraint violations?

**Analysis from docs:** "such names are supplied in separate fields of the error report message"

**Missing evidence:** Documentation mentions fields exist but doesn't list all available fields.

**Expected:** constraint_name, table_name, column_name, maybe others.

**Suggested verification:** Use PostgreSQL driver to capture all error detail fields from different constraint violation types.

**Priority:** Low  
**Impact:** Medium - error handling in applications  
**Effort:** Low (test with driver)

---

## Research Environment Notes

- **PostgreSQL version tested against (documentation):** PostgreSQL 18 (current), supported back to 14
- **MySQL version to test against:** MySQL 8.0 (documentation blocked)
- **SQLite version used for cross-verification:** SQLite 3.45+ (current)
- **Docker available:** Yes (for local verification of open questions)

## Priority for Follow-Up

| Priority | Question | Effort | Impact |
|---|---|---|---|
| High | Q1: NULL in multi-column UNIQUE | Low | High - common bug source |
| High | Q4: Deferred + partial index | Medium | High - transaction correctness |
| High | Q5: Concurrent partial index | High | High - deployment safety |
| Medium | Q2: Partial index planner | Medium | Medium - performance tuning |
| Medium | Q6: Generated column + UNIQUE | Medium | Medium - correctness |
| Medium | Q7: EXCLUDE with time ranges | Medium | High - booking systems |
| Medium | Q8: MySQL behavior | Medium | Medium - multi-DB compatibility |
| Low | Q3: Constraint validation order | Low | Low - academic interest |
| Low | Q9: Error detail fields | Low | Medium - error handling |

## Next Steps for Verification

1. **Setup Docker test environment** for PostgreSQL and MySQL
2. **Run Q1 test** — multi-column UNIQUE with NULLs
3. **Run Q4 test** — DEFERRABLE + partial index interaction
4. **Run Q5 test** — concurrent partial index creation safety
5. **Document Docker setup and test commands** for future verification
