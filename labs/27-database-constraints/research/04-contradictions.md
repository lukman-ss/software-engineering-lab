# Contradictions & Uncertainties

## Contradictions Found

### MySQL Documentation Access
- **Status:** Multiple MySQL documentation URLs returned HTTP 403 Forbidden
- **Attempted URLs:**
  - https://dev.mysql.com/doc/refman/8.0/en/constraints.html
  - https://dev.mysql.com/doc/refman/8.0/en/create-table.html
  - https://dev.mysql.com/doc/refman/8.0/en/unique-constraints.html
- **Hypothesis:** Site implements access control or rate limiting; no clear contradiction in content, only access failure
- **Impact:** MySQL-specific constraint behavior not directly verified

### Missing Industry Articles
- **Status:** Attempts to retrieve industry best practice articles resulted in 404 Not Found
- **Attempted URLs:**
  - https://www.citusdata.com/blog/2018/03/28/five-ways-to-race-condition/ (404)
  - https://martinfowler.com/articles/dbc-constraints.html (404)
- **Impact:** External validation of best practices not directly accessible in current session
- **Mitigation:** PostgreSQL docs themselves cite industry practices and academic sources

## Partial Contradictions / Clarifications

### NULL Handling Behavior
- **PostgreSQL:** NULLs are distinct by default for UNIQUE constraints
- **SQL Standard:** Implementation-defined behavior (PostgreSQL documents this explicitly)
- **Other Databases:** Cannot verify due to MySQL access failure
- **Recommendation:** Explicitly document NULLS NOT DISTINCT or NULLS DISTINCT in schema

### Partial Unique Index vs Unique Constraint
- **Documentation:** States partial unique indexes are created via `CREATE UNIQUE INDEX ... WHERE`
- **Implementation:** Both use B-tree indexes (via CREATE INDEX) vs constraint machinery
- **Practical difference:** Constraints auto-name, constraints propagate via table inheritance, constraints appear in information_schema
- **No contradiction:** These are complementary tools, not conflicting approaches

## Open Questions Without Definitive Answers

### 1. Multi-Column Unique Constraints with NULLs
- **Question:** When does "the combination of values" trigger a violation if some columns are NULL?
- **Evidence:** "Two null values are not considered equal" but unclear if "one null and one non-null" is allowed
- **Test needed:** Verify behavior with sample INSERTs in live database
- **Expected:** NULL != non-NULL, so (1, NULL) and (1, NULL) should violate but (1, NULL) and (1, 2) should not

### 2. Partial Index Planner Implication Depth
- **Question:** How sophisticated is the "simple inequality implication" that the planner recognizes?
- **Evidence:** "The system can recognize simple inequality implications, for example 'x < 1' implies 'x < 2'" [Partial Indexes]
- **Missing:** Exact algorithm or examples of non-recognized implications
- **Risk:** Application may create partial index that planner never uses
- **Mitigation:** Use EXPLAIN to verify index usage

### 3. Constraint Validation Order
- **Question:** What is the exact evaluation order when multiple constraints apply to one row?
- **Evidence:** "When a table has multiple CHECK constraints, they will be tested for each row in alphabetical order by name" [DDL Constraints: 5.5.1]
- **Missing:** Order for NOT NULL vs CHECK vs FOREIGN KEY vs UNIQUE
- **Assumption:** NOT NULL is fastest (explicit implementation) followed by CHECK, then FK/UNIQUE requiring index lookups

### 4. Deferred Constraint Validation with Partial Indexes
- **Question:** Can partial unique indexes work with DEFERRABLE constraints?
- **Evidence:** "Only UNIQUE, PRIMARY KEY, EXCLUDE, and REFERENCES (foreign key) constraints accept this clause. NOT NULL and CHECK constraints are not deferrable." [CREATE TABLE]
- **Missing:** Partial unique index created by constraint respects DEFERRABLE timing
- **Assumption:** Yes (the constraint itself is deferrable, index is implementation detail)

### 5. Concurrent Partial Index Building
- **Question:** Can partial unique indexes be built CONCURRENTLY?
- **Evidence:** CREATE INDEX CONCURRENTLY supported for regular indexes
- **Missing:** Documentation does not explicitly mention partial unique indexes in CONCURRENTLY context
- **Risk:** Build may not be atomic if constraint violated during construction
- **Assumption:** CONCURRENTLY works if WHERE clause matches existing data; otherwise fails mid-build

## MySQL-Specific Uncertainties (Unverified)
Due to HTTP 403 on dev.mysql.com, cannot verify:

1. MySQL NULL handling in UNIQUE constraints (MySQL may treat NULLs as equal by default)
2. MySQL's DEFERRABLE support (may not exist)
3. MySQL's partial index support (MySQL doesn't support partial indexes natively)
4. MySQL EXCLUDE constraint (likely not supported)

## PostgreSQL Version Consistency
All evidence from PostgreSQL 18 documentation. Minor differences may exist in:
- PostgreSQL 17 (supported)
- PostgreSQL 16 (supported)
- Older versions (unsupported): constraints may lack features like NULLS NOT DISTINCT

## Recommendations for Validation
1. Test NULL handling behavior in real database
2. Run EXPLAIN on queries with partial indexes to verify planner usage
3. Check constraint error messages with named vs unnamed constraints
4. Verify DEFERRABLE timing with partial indexes
5. Compare MySQL behavior if access becomes available (or use Docker container)