# Key Takeaways

## 1. Database constraints prevent race conditions atomically.

Read-then-write patterns (`SELECT → INSERT`) fail under concurrency. Database UNIQUE constraints check during INSERT under lock, ensuring exactly one winner and N-1 failures with SQLSTATE 23505.

## 2. Application validation and database constraints serve complementary roles.

Application validation provides user-friendly error messages (first line, UX). Database constraints provide absolute correctness protection (last line, safety net).

## 3. Partial unique indexes solve the soft-delete uniqueness problem.

`CREATE UNIQUE INDEX ... WHERE deleted_at IS NULL` enforces uniqueness only for active records, allowing email reuse after soft delete while preventing duplicate active records.

## 4. Constraints are limited to row-scoped logic.

CHECK cannot reference other rows or external services. Cross-row invariants require SERIALIZABLE isolation or triggers — not demonstrated by this lab's in-memory engine.

## 5. Always inspect SQLSTATE codes, not error text.

Text can be localized or changed. SQLSTATE `23502/23503/23505/23514` are stable. Map codes to domain errors (HTTP `409 Conflict`, `400 Bad Request`, `422 Unprocessable Entity`).

## 6. Foreign key referencing columns are not auto-indexed.

Referential integrity is enforced, but DELETE/UPDATE on referenced table may scan referencing table without index. Manually index FK referencing columns for performance.

## 7. In-memory engine is a simulator, not production PostgreSQL.

This lab models SQLSTATE behavior and concurrency patterns without requiring Docker. It is suitable for learning and testing but does not benchmark production throughput or B-tree internals.

## 8. Test concurrency explicitly.

Without concurrent tests, race conditions are invisible. `go test -race ./...` and stress tests with goroutines reveal data integrity issues that sequential tests miss.

## 9. Constraint violations return structured fields.

Error messages include constraint name, table name, and column. This enables debugging and precise error mapping without parsing text.

## 10. Production migrations require special care.

NOT VALID + VALIDATE CONSTRAINT pattern avoids table lock and allows incremental validation. This pattern is documented in PostgreSQL docs but not implemented in this lab's simulator.