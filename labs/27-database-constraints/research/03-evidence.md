# Evidence Gathered

## 1. Constraint Fundamentals

### UNIQUE Constraints
- "By default, two null values are not considered equal in this comparison." [DDL Constraints: 5.5.3]
- "NULLS NOT DISTINCT option modifies this and causes the index to treat nulls as equal" [Unique Indexes]
- "Adding a unique constraint will automatically create a unique btree index on the column or group of columns used in the constraint." [DDL Constraints: 5.5.3]
- For multi-column: "This specifies that the combination of values in the indicated columns is unique across the whole table" [DDL Constraints: 5.5.3]

### NOT NULL Constraints
- "A not-null constraint is functionally equivalent to creating a check constraint CHECK (column_name IS NOT NULL), but in PostgreSQL creating an explicit not-null constraint is more efficient." [DDL Constraints: 5.5.2]
- "In most database designs the majority of columns should be marked not null." [Tip in 5.5.2]

### FOREIGN KEY Constraints
- "This controls whether the constraint can be deferred. A constraint that is not deferrable will be checked immediately after every command." [DDL Constraints: FK section]
- Matching options:
  - MATCH FULL: "will not allow one column of a multicolumn foreign key to be null unless all foreign key columns are null"
  - MATCH SIMPLE: "allows any of the foreign key columns to be null; if any of them are null, the row is not required to have a match in the referenced table"
  - MATCH PARTIAL: "is not yet implemented" [DDL Constraints: FK section]
- Actions: NO ACTION (default), RESTRICT, CASCADE, SET NULL, SET DEFAULT

### CHECK Constraints
- "A check constraint can also refer to several columns." [DDL Constraints: 5.5.1]
- "CHECK expressions cannot contain subqueries nor refer to variables other than columns of the current row" [DDL Constraints: 5.5.1 Note]
- "PostgreSQL assumes that CHECK constraints' conditions are immutable" [DDL Constraints: 5.5.1 Note]
- "If what you desire is a one-time check against other rows at row insertion, rather than a continuously-maintained consistency guarantee, a custom trigger can be used to implement that" [DDL Constraints: 5.5.1 Note]

### Primary Key
- "Primary keys are useful both for documentation purposes and for client applications" [DDL Constraints: 5.5.4]
- "Adding a primary key will automatically create a unique btree index on the column or group of columns listed in the primary key" [DDL Constraints: 5.5.4]

## 2. Race Condition Prevention

### Constraint Atomicity
- Constraints are checked "when rows are inserted or updated" [DDL Constraints: 5.5.1 Note]
- The assumption of immutability for CHECK constraints "justifies examining CHECK constraints only when rows are inserted or updated" [DDL Constraints: 5.5.1 Note]
- For exclusion constraints: "Each exclude_element defines a column of the index" [DDL Constraints: Exclusion section]
- UNIQUE constraints use B-tree indexes which provide "readers/writers lock semantics" at the index level

### Row-Level Locks
- "FOR UPDATE causes the rows retrieved by the SELECT statement to be locked as though for update. This prevents them from being locked, modified or deleted by other transactions until the current transaction ends." [Explicit Locking: 13.3.2]
- "Within a REPEATABLE READ or SERIALIZABLE transaction, however, an error will be thrown if a row to be locked has changed since the transaction started." [Explicit Locking: 13.3.2]

### Serializable Transactions
- "Serializable transactions are just Repeatable Read transactions which add nonblocking monitoring for dangerous patterns of read/write conflicts." [App-Level Consistency: 13.4.1]
- "When a pattern is detected which could cause a cycle in the apparent order of execution, one of the transactions involved is rolled back to break the cycle." [App-Level Consistency: 13.4.1]

## 3. Partial Unique Indexes

### Syntax & Behavior
- "A partial index is an index built over a subset of a table; the subset is defined by a conditional expression (called the predicate of the partial index)." [Partial Indexes: 11.8]
- "To create a partial index that suits our example, use a command such as this: CREATE INDEX access_log_client_ip_ix ON access_log (client_ip) WHERE NOT (client_ip > inet '192.168.100.0' AND client_ip < inet '192.168.100.255');" [Partial Indexes: Example 11.1]
- "The predicate must match the conditions used in the queries that are supposed to benefit from the index" [Partial Indexes: Note]
- "PostgreSQL does not have a sophisticated theorem prover that can recognize mathematically equivalent expressions that are written in different forms" [Partial Indexes: Note]

### One-Active-Per-Key Pattern
- "Suppose that we have a table describing test outcomes. We wish to ensure that there is only one 'successful' entry for a given subject and target combination, but there might be any number of 'unsuccessful' entries." [Partial Indexes: Example 11.3]
- "CREATE UNIQUE INDEX tests_success_constraint ON tests (subject, target) WHERE success;" [Partial Indexes: Example 11.3]
- "This is a particularly efficient approach when there are few successful tests and many unsuccessful ones." [Partial Indexes: Example 11.3]

### Limitations
- "Matchings takes place at query planning time, not at run time. As a result, parameterized query clauses do not work with a partial index." [Partial Indexes: Note]
- "For example a prepared query with a parameter might specify 'x < ?' which will never imply 'x < 2' for all possible values of the parameter." [Partial Indexes: Note]
- "Setting up a partial index indicates that you know at least as much as the query planner knows" [Partial Indexes: Note]

## 4. Constraint Validation & Error Handling

### Error Codes
- "23502: not_null_violation" [Error Codes: Class 23]
- "23503: foreign_key_violation" [Error Codes: Class 23]
- "23505: unique_violation" [Error Codes: Class 23]
- "23514: check_violation" [Error Codes: Class 23]
- "23P01: exclusion_violation" [Error Codes: Class 23]

### Constraint Names
- "If the constraint is violated, the constraint name is present in error messages, so constraint names like col must be positive can be used to communicate helpful constraint information to client applications" [CREATE TABLE: CONSTRAINT section]
- "If you don't specify a constraint name in this way, the system chooses a name for you." [DDL Constraints: Check Constraints]

### Validation Timing
- "Currently, CHECK expressions cannot contain subqueries nor refer to variables other than columns of the current row" [DDL Constraints: 5.5.1]
- "The warning above about not referencing other table data is really a special case of this restriction" [DDL Constraints: 5.5.1]

## 5. When NOT to Use Constraints

### Cross-Row/Mutable Data
- "PostgreSQL does not support CHECK constraints that reference table data other than the new or updated row being checked. While a CHECK constraint that violates this rule may appear to work in simple tests, it cannot guarantee that the database will not reach a state in which the constraint condition is false" [DDL Constraints: 5.5.1 Note]
- "If possible, use UNIQUE, EXCLUDE, or FOREIGN KEY constraints to express cross-row and cross-table restrictions" [DDL Constraints: 5.5.1 Note]

### Performance & Locking
- "Adding a unique constraint will automatically create a unique btree index" implies index maintenance overhead
- Constraint checking requires lock acquisition similar to index updates
- "Exclusion constraints are implemented using an index that has the same name as the constraint" [DDL Constraints: Exclusion section]

### Business Logic Complexity
- "CHECK expressions cannot contain subqueries nor refer to variables other than columns of the current row" limits complex validations
- "References to other tables are not allowed" for generated columns (similar restriction) [DDL Constraints: Generated Columns]
- Alternatives: triggers, application validation, serializable transactions

## 6. Partitioning & Constraints

### Unique Constraints on Partitioned Tables
- "To create a unique or primary key constraint on a partitioned table, the partition keys must not include any expressions or function calls and the constraint's columns must include all of the partition key columns." [Partitioning: 5.12.2.3 Limitations]
- "This limitation exists because the individual indexes making up the constraint can only directly enforce uniqueness within their own partitions; therefore, the partition structure itself must guarantee that there are not duplicates in different partitions." [Partitioning: 5.12.2.3 Limitations]

### Exclusion Constraints on Partitioned Tables
- "Similarly an exclusion constraint must include all the partition key columns. Furthermore the constraint must compare those columns for equality (not e.g. &&)" [Partitioning: 5.12.2.3 Limitations]