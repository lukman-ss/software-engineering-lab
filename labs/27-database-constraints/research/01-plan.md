# Research Plan: Database Constraints & Data Integrity Enforcement

## Context
Senior Software Engineer lab exploring how database constraints enforce data integrity and prevent race conditions.

## Research Questions

### 1. Constraint Types & Integrity Roles
- Q1.1: How do UNIQUE, NOT NULL, FOREIGN KEY, and CHECK constraints each enforce integrity?
- Q1.2: What are the exact semantics (NULL handling, deferrability, evaluation timing)?
- Q1.3: Which constraints automatically create indexes?

### 2. Race Conditions & Constraint Prevention
- Q2.1: How do constraints prevent write-write race conditions without explicit locking?
- Q2.2: What is the relationship between constraint enforcement, row locks, and transaction isolation?
- Q2.3: What races remain that constraints alone cannot solve?

### 3. Partial Unique Indexes
- Q3.1: How do partial unique indexes (WHERE predicate) differ from UNIQUE constraints?
- Q3.2: Use cases: one-active-per-key, one-null-row patterns?
- Q3.3: Limitations (planner implication rules, parameterized queries)?

### 4. Validation & Error Handling
- Q4.1: What SQLSTATE error codes map to each constraint type?
- Q4.2: How should applications map constraint violations to user-facing errors?
- Q4.3: Named constraints' role in error communication?

### 5. When NOT to Use Constraints
- Q5.1: What business logic is unsuitable for CHECK constraints (cross-row, mutable functions)?
- Q5.2: Performance costs (locks, validation scans)?
- Q5.3: Alternatives: triggers, application validation, serializable transactions?

## Source Priority
1. PostgreSQL documentation (primary)
2. MySQL documentation (blocked: 403 on dev.mysql.com — see 04-contradictions.md)
3. Academic sources (referenced via PG biblio: stonebraker89, olson93, seshadri95)
4. Industry articles (attempted, several 404)

## Method
- Fetch primary doc pages, extract claims with quotes
- Note contradictions/gaps honestly
- Synthesize into report
