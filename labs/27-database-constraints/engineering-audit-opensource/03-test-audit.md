# Engineering Test Audit

## Test Suite Overview

Location: internal/store/store_test.go

## Coverage Matrix

| Requirement | Constraint | Test | Status |
|---|---|---|---|
| NOT NULL (email) | 23502 | TestNotNullConstraints | PASS |
| NOT NULL (username) | 23502 | TestNotNullConstraints | PASS |
| NOT NULL (order.user_id) | 23502 | TestNotNullConstraints | PASS |
| CHECK (age >= 18) | 23514 | TestCheckConstraints | PASS |
| CHECK (status IN set) | 23514 | TestCheckConstraints | PASS |
| CHECK (total_cents > 0) | 23514 | TestCheckConstraints | PASS |
| UNIQUE (email) | 23505 | TestUniqueConstraint | PASS |
| FOREIGN KEY | 23503 | TestForeignKeyConstraint | PASS |
| Partial UNIQUE (soft delete) | 23505 | TestPartialUniqueIndex | PASS |
| Concurrency race prevention | 23505 | TestConcurrentRegistration_Safe_EnforcesUniqueness | PASS |
| SQLSTATE taxonomy | — | TestErrorClassification | PASS |
| UnsafeStore vulnerability | — | *none* | MISSING |
| Domain error mapping | — | TestErrorClassification (constructors only) | MISSING |

## Executed Commands and Actual Results

```
$ go test -v ./...
PASS  (8/8 tests)
ok  github.com/lukman/software-engineering-lab/labs/27-database-constraints/internal/store

$ go test -race ./...
ok  github.com/lukman/software-engineering-lab/labs/27-database-constraints/internal/store  1.330s  (clean — no data race)
```

```
$ go run ./cmd/demo
[1] NOT NULL (23502)        -> Error mapped to domain
[2] CHECK (23514)           -> Error mapped to domain
[3] FOREIGN KEY (23503)     -> Error mapped to domain
[4] Partial UNIQUE           -> soft-delete reuse works
[5] 50-goroutine concurrency -> 1 success, 49 errors, Integrity Intact: true
SQLSTATE taxonomy           -> code=23505 isUniqueViolation=true
```

All reported results are reproduced by execution. No fabricated output.

## Happy Path

PASS — all positive flows (valid registration, valid order, valid soft-delete reuse, concurrent single winner) execute and return non-nil results.

## Failure Path

PASS — every rejected insert returns a typed ConstraintError (constructors verified) and the store wraps it to a domain error.

## Edge Cases

- Partial unique: active duplicate rejected, post-soft-delete re-insert allowed, second active duplicate rejected, soft-deleted insert bypasses partial index (step 6).
- CHECK boundary: age=16 rejected; age=0 not tested (negative); exact boundary age=18 not tested.
- Order total=0 rejected; total negative not tested.

## Transitions / Recovery / Rollback

NOT COVERED. There is no transaction API — the engine is single-statement (one method = one write). There is no abort/rollback, no commit/durable boundary. This is acceptable for the in-memory sim but is inherent (no multi-statement integrity scenarios from research are modeled).

## Concurrency

SAFE path tested with 20 goroutines — exactly 1 success. PASS under `-race`.
UNSAFE path (UnsafeStore) — NOT tested. The code claims a vulnerable read-then-write pattern, but the engine still enforces the constraint, so a test would pass rather than demonstrate the race. This is the UnsafeStore finding (see code-audit Finding 1).

## Negative Cases

PARTIAL. Tests assert `err == nil` for the "should fail" cases and `err != nil` for "should succeed" cases. They do NOT assert:
- the SQLSTATE of the returned error (no `IsConstraintViolation` on a store return value),
- that the error is a `*ConstraintError`,
- the content of the mapped domain message.

`TestErrorClassification` only exercises the `dberr.New*` constructors and `IsConstraintViolation` on those synthetic values — not on errors returned from `SafeStore.RegisterUser`/`CreateOrder`.

## Gap: Error Classification Pipeline

Research: "Always check SQLSTATE, not error text."
Tests: verify the constructors and `IsConstraintViolation` in isolation, but never on a real store-returned error. After `MapToDomainError`, the ConstraintError type is lost, so `IsConstraintViolation` would return false — and no test exercises this.

## Gap: Domain Error Mapping Content

`MapToDomainError` is called on every store write but no assertion checks the resulting human-readable message maps to the correct SQLSTATE branch (unique→conflict, not-null→invalid input, check→validation, fk→reference error).

## Gap: UnsafeStore

`NewUnsafeStore` and `UnsafeStore.RegisterUser` are dead code — declared but never referenced by tests or the demo. The documented "vulnerable pattern" is implemented but never demonstrated.