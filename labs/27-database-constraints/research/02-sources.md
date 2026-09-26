# Sources Accessed

## Successfully Accessed

### 1. PostgreSQL Documentation: CREATE TABLE Constraints
**URL:** https://www.postgresql.org/docs/current/sql-createtable.html#SQL-CREATETABLE-CONSTRAINT-UNIQUE
**Section:** Column and table constraint syntax (implicit reference)
**Key coverage:**
- Constraint syntax variants (column vs table constraints)
- DEFERRABLE options
- NOT NULL, CHECK, UNIQUE, PRIMARY KEY, FOREIGN KEY, EXCLUDE
- NULLS NOT DISTINCT option for unique constraints

### 2. PostgreSQL Documentation: Chapter 5.5 - Constraints
**URL:** https://www.postgresql.org/docs/current/ddl-constraints.html
**Section:** 5.5. Constraints
**Key coverage:**
- Check constraints: expression requirements, null evaluation
- Unique constraints: NULL handling, multi-column, index creation
- Primary keys: NOT NULL + UNIQUE, index auto-creation
- Foreign keys: MATCH options, ON DELETE/UPDATE actions, referential actions
- Exclusion constraints: GiST index requirement, comparison operators

### 3. PostgreSQL Documentation: Chapter 11.6 - Unique Indexes
**URL:** https://www.postgresql.org/docs/current/indexes-unique.html
**Section:** 11.6. Unique Indexes
**Key coverage:**
- UNIQUE INDEX syntax
- NULL handling (distinct by default, NULLS NOT DISTINCT)
- Auto-index creation for constraints

### 4. PostgreSQL Documentation: Chapter 11.8 - Partial Indexes
**URL:** https://www.postgresql.org/docs/current/indexes-partial.html
**Section:** 11.8. Partial Indexes
**Key coverage:**
- CREATE INDEX ... WHERE predicate
- Example: exclude common values, exclude uninteresting values
- Partial unique indexes for one-active-per-key pattern
- Planner implication requirements

### 5. PostgreSQL Documentation: Chapter 5.12 - Table Partitioning
**URL:** https://www.postgresql.org/docs/current/ddl-partitioning.html
**Section:** 5.12. Table Partitioning
**Key coverage:**
- Partition constraint inheritance
- Best practices for partition keys and count
- Limitations: unique constraints require partition key columns

### 6. PostgreSQL Documentation: Chapter 13.3 - Explicit Locking
**URL:** https://www.postgresql.org/docs/current/explicit-locking.html
**Section:** 13.3. Explicit Locking
**Key coverage:**
- Row-level lock modes (FOR UPDATE, FOR SHARE, etc.)
- Lock conflict matrix
- Deadlock detection
- Row lock acquisition during UPDATE/DELETE

### 7. PostgreSQL Documentation: Chapter 13.4 - Application-Level Consistency
**URL:** https://www.postgresql.org/docs/current/applevel-consistency.html
**Section:** 13.4. Data Consistency Checks at the Application Level
**Key coverage:**
- Serializable transactions for consistency
- SELECT FOR UPDATE/SHARE usage
- Read/write conflict scenarios

### 8. PostgreSQL Documentation: Appendix A - Error Codes
**URL:** https://www.postgresql.org/docs/current/errcodes-appendix.html
**Section:** Appendix A
**Key coverage:**
- SQLSTATE class 23: Integrity constraint violations
- 23502: not_null_violation
- 23503: foreign_key_violation
- 23505: unique_violation
- 23514: check_violation
- 23P01: exclusion_violation

## Blocked / Unreachable

### MySQL Documentation
**URLs attempted:**
- https://dev.mysql.com/doc/refman/8.0/en/constraints.html (403)
- https://dev.mysql.com/doc/refman/8.0/en/create-table.html (403)
- https://dev.mysql.com/doc/refman/8.0/en/unique-constraints.html (403)

**Status:** HTTP 403 Forbidden - site may have rate limiting or require different user agent

### Industry Articles
**URLs attempted:**
- https://www.citusdata.com/blog/2018/03/28/five-ways-to-race-condition/ (404)
- https://martinfowler.com/articles/dbc-constraints.html (404)

**Status:** Not found vs attempted archives

## Academic References (via PostgreSQL bibliography)
- stonebraker89: Architecture of a Database System
- olson93: Partial indexing in POSTGRES: research project
- seshadri95: Partial indexing research