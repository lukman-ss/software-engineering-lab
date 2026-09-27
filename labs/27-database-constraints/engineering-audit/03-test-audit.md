# Test Audit

Target Lab: labs/27-database-constraints

## Test Suite Execution

Commands:
```bash
go test -v ./...
go test -race ./...
```

Results:
- Total tests executed: 8
- Passed: 8
- Failed: 0
- Race detector warnings: 0

## Coverage Analysis

1. **NOT NULL Constraints (`TestNotNullConstraints`)**
   - Missing email rejected: PASS
   - Missing username rejected: PASS
   - Missing foreign key ID in child order rejected: PASS

2. **CHECK Constraints (`TestCheckConstraints`)**
   - Underage boundary violation (`age < 18`): PASS
   - Enum domain check violation (`status IN (...)`): PASS
   - Numerical boundary check violation (`total_cents <= 0`): PASS

3. **UNIQUE Constraints (`TestUniqueConstraint`)**
   - Sequential duplicate email rejected: PASS

4. **FOREIGN KEY Constraints (`TestForeignKeyConstraint`)**
   - Non-existent parent reference rejected: PASS
   - Valid parent reference accepted: PASS

5. **Partial Unique Index (`TestPartialUniqueIndex`)**
   - Initial active record created: PASS
   - Duplicate active record rejected: PASS
   - Soft-delete clears partial index: PASS
   - Re-registration of identical email succeeds after soft delete: PASS
   - Second active duplicate rejected: PASS
   - Direct insertion of soft-deleted duplicate allowed: PASS

6. **Concurrency Safety & Race Conditions**
   - Safe registration under 20 concurrent goroutines (`TestConcurrentRegistration_Safe_EnforcesUniqueness`): Exactly 1 succeeds, 19 rejected with unique constraint violation. PASS.
   - Unsafe registration under 20 concurrent goroutines (`TestConcurrentRegistration_Unsafe_SuffersRaceCondition`): Multiple records inserted due to read-then-write race. PASS.

7. **Error Classification (`TestErrorClassification`)**
   - SQLSTATE mapping verified for `23502`, `23503`, `23505`, and `23514`. PASS.
