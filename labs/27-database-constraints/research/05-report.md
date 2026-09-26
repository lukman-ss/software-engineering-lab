# Research Report: Database Constraints & Data Integrity

## Executive Summary

Database constraints provide declarative, ACID-guaranteed data integrity enforcement. By embedding validation rules directly into the storage engine, they offer stronger guarantees than application-level checks alone — including automatic race condition prevention during concurrent writes. The key tradeoff is that constraints are limited to row-scoped logic and cannot reference mutable external data.

---

## 1. Constraint Types & Integrity Enforcement

| Constraint | Role | Auto-Index? | Deferrable? | NULL Behavior |
|---|---|---|---|---|
| NOT NULL | Prevents NULL in column | No (uses heap tuple check) | No (always immediate) | N/A - null not allowed |
| CHECK | Boolean expression on row values | No | No | NULL result → constraint passes |
| UNIQUE | Enforces distinct values across rows | Yes (B-tree) | Yes | NULLs distinct by default; `NULLS NOT DISTINCT` opts in |
| PRIMARY KEY | Identity + uniqueness guarantee | Yes (B-tree) | Yes | NOT NULL enforced automatically |
| FOREIGN KEY | Referential integrity between tables | Index on referenced table, not referencing | Yes | NULLs bypass constraint (MATCH SIMPLE) |
| EXCLUDE | Custom operator-based non-overlap guarantee | Yes (GiST) | Yes | Same as UNIQUE with operator |

### Key Details

- **CHECK constraints** cannot contain subqueries or reference other rows. PostgreSQL explicitly warns: "it cannot guarantee that the database will not reach a state in which the constraint condition is false." Cross-row logic requires triggers.

- **UNIQUE/PRIMARY KEY** create B-tree indexes automatically — no need to create separate indexes.

- **`NULLS NOT DISTINCT`** option (PostgreSQL 15+) treats NULLs as equal in unique constraints — useful for "at most one NULL" semantics.

- **Foreign keys** create index on referenced (parent) columns but NOT on referencing columns. Explicit index on referencing columns recommended for performance with CASCADE operations.

---

## 2. Race Conditions & Constraint-Based Prevention

### How Constraints Prevent Races

Constraints eliminate classic **read-then-write** race conditions:

```
-- Without constraint (vulnerable):
tx1: SELECT count(*) FROM items WHERE code = 'x' → 0
tx2: SELECT count(*) FROM items WHERE code = 'x' → 0
tx1: INSERT INTO items (code) VALUES ('x') → INSERT succeeds
tx2: INSERT INTO items (code) VALUES ('x') → INSERT succeeds (BUG)

-- With UNIQUE constraint (safe):
tx1: INSERT INTO items (code) VALUES ('x') → INSERT succeeds
tx2: INSERT INTO items (code) VALUES ('x') → ERROR 23505 unique_violation
```

Constraints work because the database engine uses **row-level locks acquired during INSERT/UPDATE** that prevent conflicting writes concurrently. The B-tree index for a UNIQUE constraint is updated atomically with the row insert, under a row-exclusive lock.

### Lock Mechanism
- INSERT acquires `ROW EXCLUSIVE` lock on the table
- The B-tree index insertion uses internal page/row locks
- Concurrent inserts with identical key values will block, and the engine detects the conflict, raising `unique_violation` (SQLSTATE 23505)

### What Constraints Cannot Solve

- **Multi-statement logical invariants** (e.g., "sum of balances must be zero") require either:
  - Serializable isolation (SQLSTATE 40001 on serialization failure)
  - Explicit locking (`SELECT FOR UPDATE`)
  - Triggers with explicit locks

- **Cross-row validations** that read other rows (CHECK constraints can't do this)

- **External service checks** (e.g., validating an ID against an external API)

### Serializable Transactions
For complex multi-row invariants, PostgreSQL's SERIALIZABLE isolation level:
- "Serializable transactions are just Repeatable Read transactions which add nonblocking monitoring for dangerous patterns of read/write conflicts."
- On detection of cycle: "one of the transactions involved is rolled back to break the cycle" (SQLSTATE 40001 — `serialization_failure`)

---

## 3. Partial Unique Indexes

### Syntax
```sql
CREATE UNIQUE INDEX users_active_email_idx ON users (email) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX tests_success_constraint ON tests (subject, target) WHERE success;
```

### When to Use
- **Soft delete pattern**: `WHERE deleted_at IS NULL` — only one active record per unique key
- **Conditional uniqueness**: "only one record can have `is_default = true` per user"
- **Performance optimization**: Exclude rows that never need uniqueness enforcement

### Critical Limitations
- **Predicate matching**: "a partial index can be used in a query only if the system can recognize that the WHERE condition of the query mathematically implies the predicate of the index"
- **No parameterized implication**: Prepared statements with `?` parameters cannot imply a constant predicate
- **Planner sophistication**: Only recognizes "simple inequality implications" (e.g., `x < 1` implies `x < 2`)

### Constraint vs Index
- Unique constraints create indexes but also appear in `information_schema.table_constraints`
- Partial unique indexes created via `CREATE UNIQUE INDEX` are not "constraints" — they don't appear in the constraint catalog
- Constraint names appear in error messages: "constraint names like col must be positive can be used to communicate helpful constraint information"

---

## 4. Constraint Validation & Error Handling

### SQLSTATE Error Codes (Class 23 — Integrity Constraint Violation)

| Code | Name | Meaning |
|---|---|---|
| `23502` | `not_null_violation` | Row violates NOT NULL constraint |
| `23503` | `foreign_key_violation` | Foreign key violation |
| `23505` | `unique_violation` | Unique/PK constraint violated |
| `23514` | `check_violation` | CHECK constraint violated |
| `23P01` | `exclusion_violation` | EXCLUDE constraint violated |
| `23001` | `restrict_violation` | Foreign key RESTRICT action |

### Best Practices for Error Handling
1. **Always check SQLSTATE, not error text** — text is localized; codes are stable
2. **Name constraints explicitly** — names appear in error details:
   ```sql
   CONSTRAINT email_unique UNIQUE (email)
   ```
   The constraint name is reported in error objects, allowing application to map to user-facing messages
3. **Map 23xxx codes to user-friendly errors** — e.g., 23505 → "A user with this email already exists"
4. **Retry on 40001** (`serialization_failure`) — wrap transactional work in retry logic

### Constraint Validation Timing
- Constraints checked at INSERT/UPDATE time (not at query time)
- NOT NULL is checked first and is fastest (no index lookup)
- CHECK constraints tested "in alphabetical order by name" (for multi-check scenarios)
- UNIQUE/FK require index lookups — slower than CHECK

---

## 5. When NOT to Use Database Constraints

### Anti-Patterns Where Constraints Fail

| Scenario | Why Constraints Don't Work | Alternative |
|---|---|---|
| Cross-row aggregates (sum, count) | CHECK cannot reference other rows | Serializable transactions + explicit checks |
| External service validation | Function may return different values; violates immutability | Application-level validation before DB write |
| Mutable function in CHECK | "PostgreSQL assumes that CHECK constraints' conditions are immutable" — violating this can lead to corrupt state | TRIGGERS with careful design or app-level |
| Complex multi-step business rules | Single-row constraint can't express workflow state | State machine in application layer |
| Non-deterministic time windows | `CURRENT_TIMESTAMP` breaks constraint exclusion/pruning | Explicit locking; design with fixed boundaries |

### Performance Considerations
- **Index creation cost**: UNIQUE/PRIMARY KEY create B-tree indexes — writes slower
- **Validation during INSERT/UPDATE**: All constraints must pass before commit
- **Foreign key overhead**: Referenced table rows are locked; CASCADE operations require index scans
- **Deferred constraints**: Checking at transaction commit can cause unexpected rollback if earlier operations relied on "eventual validity"

### Partitioning Limitations
- **Unique constraints on partitioned tables**: Must include ALL partition key columns
- **No cross-partition uniqueness** enforcement without full partition key in constraint
  > "the individual indexes making up the constraint can only directly enforce uniqueness within their own partitions; therefore, the partition structure itself must guarantee that there are not duplicates in different partitions"
- EXCLUDE constraints have the same limitation.

### When to Choose Triggers Instead
- Need cross-row validation
- Need to log violations before rejecting
- Need conditional behavior (soft-reject vs hard-reject)
- Constraints documented as "appear to work in simple tests" but fail under concurrent load — triggers with explicit locking are the correct solution for cross-row checks

---

## Recommendations (Senior Engineer Level)

1. **Use constraints aggressively** for single-row and referential integrity — they are your last line of defense against data corruption
2. **Design schemas with constraints first**, then add application validation as user experience enhancement
3. **Always name constraints** for predictable error messages
4. **Use partial unique indexes** for soft-delete patterns and conditional uniqueness — far more efficient than triggers
5. **For multi-row invariants**: prefer SERIALIZABLE isolation with retry logic; fall back to explicit locking only when serialization is too costly
6. **Never put mutable functions or external service calls** in CHECK constraints
7. **Validate partitioning keys** against unique constraints — if partition key isn't in the constraint, uniqueness is not guaranteed across partitions
8. **Map error codes, not error text** — error messages are localized and change between versions