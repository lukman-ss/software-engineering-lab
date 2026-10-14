---
title: "Database Constraints in Go: Atomic Integrity, Zero Race Conditions"
published: false
description: "Learn how moving validation to database constraints (NOT NULL, CHECK, UNIQUE, FK, Partial Index) in Go prevents race conditions atomically."
tags: go, databases, architecture, backend
---

# Database Constraints in Go: Atomic Integrity, Zero Race Conditions

## 1. Introduction: The Fragility of Application-Layer Validation

Most backend engineers write entity creation and user registration logic using the familiar check-then-act pattern:

```go
// Anti-pattern: Application-level existence check
func (s *UnsafeStore) RegisterUser(u User) error {
    existing := s.FindByEmail(u.Email) // SELECT
    if existing != nil {
        return ErrAlreadyExists
    }
    s.SimulateLatency(1 * time.Millisecond) // scheduling jitter or I/O delay
    return s.InsertUserUnsafe(u)            // INSERT
}
```

This pattern passes every unit test written with sequential invocations. The engineer verifies that invalid data is rejected, duplicates are flagged, and happy-path records are inserted. But in production, requests arrive concurrently across dozens of goroutines and multiple load-balanced application instances.

Under concurrent load, the check-then-act pattern collapses. The gap between the `SELECT` and the `INSERT` is an open race window. Two requests for the same email address arrive simultaneously, both read "record does not exist", both proceed past the condition, and both perform an `INSERT`. The database ends up with duplicate rows that violate business invariants.

Lab 27 of the Software Engineering Lab demonstrates this vulnerability empirically and proves its architectural solution. An in-memory Go relational engine models the semantics of PostgreSQL constraint enforcement. When tested with 20 concurrent goroutines attempting to register the exact same email address, the unprotected `UnsafeStore` systematically produces multiple duplicate records. Conversely, the constraint-backed `SafeStore` allows exactly one insert to succeed and reliably rejects the other 19 attempts via SQLSTATE `23505`.

## 2. Why Data Integrity Must Live in the Storage Engine

Data duplication is never merely a presentation defect. Once inconsistent data crosses the write boundary and commits to disk, remedial costs grow exponentially:

1. **Downstream Pipeline Failures**: Batch synchronization jobs, search indexing workers, and event stream processors crash when expecting single-entity lookups on unique fields.
2. **Expensive Data Reconciliation**: Reconciling duplicates in a production database containing millions of rows requires table locks, complex deduplication scripts, and risky manual customer support interventions.
3. **Broken Business Rules**: In financial workflows, duplicate accounts or orders can trigger double charges, improper credit allocations, or orphaned fulfillment requests.

PostgreSQL documentation on DDL Constraints (Section 5.5) and Explicit Locking (Section 13.3) explains why the database can guarantee what application code cannot:

> "Adding a unique constraint will automatically create a unique B-tree index on the column or group of columns used in the constraint."

When an `INSERT` or `UPDATE` operation executes, PostgreSQL acquires a `ROW EXCLUSIVE` lock and updates the underlying unique B-tree index atomically within the storage engine. There is zero scheduling gap between checking for key existence and writing the index leaf node. If a conflicting key exists, the write is aborted before the transaction commits.

Furthermore, PostgreSQL and SQLite enforce `NOT NULL` and `CHECK` constraints strictly at the moment of insertion or modification, rather than on query reads. This ensures that the storage layer serves as an immutable, non-bypassable barrier protecting data integrity.

## 3. The 2-Layer Mental Model: Defense in Depth

To design robust systems, engineers must treat application validation and database constraints as complementary layers with distinct responsibilities:

```text
Incoming Request
       │
       ▼
┌────────────────────────────────────────┐
│  Layer 1: Application Validation       │
│  - Form formatting, email syntax check │
│  - UX feedback & early HTTP 400s       │
│  - Saves DB network roundtrips         │
└───────────────────┬────────────────────┘
                    │ (Valid Format)
                    ▼
┌────────────────────────────────────────┐
│  Layer 2: Database Constraints         │
│  - NOT NULL (23502), CHECK (23514)     │
│  - UNIQUE (23505), FOREIGN KEY (23503) │
│  - Partial Index (23505 on predicate)  │
│  - Atomic, thread-safe, single truth   │
└────────────────────────────────────────┘
```

The mental model comprises two essential principles:

- **Application Validation is for User Experience**: It provides immediate, contextual error messages on form inputs, validates formatting (such as RFC 5322 email syntax), and prevents pointless network trips to the database. However, it can be bypassed by direct database queries, background jobs, migration scripts, or concurrency interleaving.
- **Database Constraints are for Correctness Guarantee**: They represent the final, unbreakable authority on schema invariants. Even if application code has bugs, even if two microservices race against each other, the database rejects invalid state.

Application validation is the first line of defense; database constraints are the last line of defense. Never discard the safety net just because the first line looks clean.

## 4. Architecture and Lab Components

The codebase in `labs/27-database-constraints` provides a fully self-contained Go implementation with zero external dependencies:

```text
labs/27-database-constraints/
├── cmd/demo/main.go            # Interactive demonstration of safe vs unsafe concurrency
├── internal/
│   ├── dberr/errors.go         # SQLSTATE Class 23 taxonomy and domain error mapping
│   ├── engine/engine.go        # Thread-safe in-memory table store and constraint evaluator
│   ├── model/model.go          # Domain entities: User and Order
│   └── store/
│       ├── store.go            # UnsafeStore vs SafeStore implementations
│       └── store_test.go       # Concurrency test suite running with -race
└── go.mod
```

### Component Roles

- **`internal/model`**: Defines the data models `User` (fields: `ID`, `Email`, `Username`, `Age`, `Status`, `DeletedAt`) and `Order` (fields: `ID`, `UserID`, `TotalCents`, `Status`).
- **`internal/dberr`**: Implements structured error types holding the ANSI SQLSTATE code, constraint rule name, table name, and diagnostic details. It includes a domain mapper (`MapToDomainError`) to convert database error codes into appropriate HTTP and domain responses.
- **`internal/engine`**: Implements the simulated database engine using Go maps and a `sync.RWMutex`. It enforces NOT NULL, CHECK, UNIQUE, FOREIGN KEY, and partial index rules atomically inside critical sections.
- **`internal/store`**: Exposes two store implementations:
  - `UnsafeStore`: Mimics naive application-level checks, explicitly demonstrating race conditions.
  - `SafeStore`: Directly delegates operations to the constraint engine, handling SQLSTATE responses gracefully.
- **`cmd/demo`**: Runs an interactive demonstration executing 50 concurrent goroutines against both store implementations to visualize the difference in integrity outcomes.

## 5. Implementation of the Five Constraint Types

The engine enforces five standard constraint types, modeled directly after PostgreSQL semantics.

### 1. NOT NULL Constraint (`SQLSTATE 23502`)

Guarantees that mandatory columns cannot contain zero or empty values. In `internal/engine/engine.go`:

```go
func (e *Engine) validateUserSchema(u *model.User) error {
    if u.Email == "" {
        return dberr.NewNotNullViolation("users", "email", "users_email_not_null")
    }
    if u.Username == "" {
        return dberr.NewNotNullViolation("users", "username", "users_username_not_null")
    }
    // ...
}
```

For orders, foreign keys must also not be zero:

```go
func (e *Engine) validateOrderSchema(o *model.Order) error {
    if o.UserID == 0 {
        return dberr.NewNotNullViolation("orders", "user_id", "orders_user_id_not_null")
    }
    // ...
}
```

### 2. CHECK Constraint (`SQLSTATE 23514`)

Enforces boolean predicates on the record's column values:

```go
// Age boundary check
if u.Age < 18 {
    return dberr.NewCheckViolation("users", "users_age_check",
        "email %q violates check: age >= 18", u.Email)
}

// Enumeration membership check
validStatuses := map[string]bool{"active": true, "suspended": true, "pending": true}
if !validStatuses[u.Status] {
    return dberr.NewCheckViolation("users", "users_status_check",
        "status %q not in (active, suspended, pending)", u.Status)
}

// Positive monetary value check
if o.TotalCents <= 0 {
    return dberr.NewCheckViolation("orders", "orders_total_cents_check",
        "total_cents must be positive, got %d", o.TotalCents)
}
```

CHECK constraints are row-scoped: they evaluate the current row in isolation without performing subqueries or cross-table lookups.

### 3. UNIQUE Constraint (`SQLSTATE 23505`)

Enforces uniqueness across the entire dataset. In `internal/engine/engine.go`, the engine maintains a dedicated hash index `emailIndex`:

```go
if existingID, exists := e.emailIndex[u.Email]; exists && existingID != u.ID {
    return dberr.NewUniqueViolation("users", "users_email_key",
        "duplicate key value violates unique constraint %q", "users_email_key")
}
```

Because `InsertUser` holds `e.mu.Lock()` throughout the entire duration of index checking and map insertion, no concurrent worker can interleave between verification and record creation.

### 4. FOREIGN KEY Constraint (`SQLSTATE 23503`)

Guarantees referential integrity between tables, preventing orphan child records:

```go
// internal/engine/engine.go - InsertOrder foreign key check
e.mu.Lock()
defer e.mu.Unlock()

user, exists := e.users[o.UserID]
if !exists || user == nil {
    return dberr.NewForeignKeyViolation("orders", "fk_orders_user",
        "insert or update on table %q violates foreign key constraint", "orders")
}
```

Attempting to insert an order with `UserID = 99999` fails immediately with code `23503`.

### 5. PARTIAL UNIQUE INDEX (`SQLSTATE 23505` with Predicate)

A standard unique constraint on `email` prevents users from re-registering with an email address if an old account was soft-deleted (`deleted_at IS NOT NULL`). A partial unique index solves this by scoping uniqueness only to active rows:

```sql
CREATE UNIQUE INDEX users_active_email_idx ON users (email) WHERE deleted_at IS NULL;
```

In the lab's engine, this is implemented using a specialized `activeEmails` map:

```go
if partialUnique {
    if u.DeletedAt == nil {
        if existingID, exists := e.activeEmails[u.Email]; exists && existingID != u.ID {
            return dberr.NewUniqueViolation("users", "users_active_email_idx",
                "email %q already exists in active records", u.Email)
        }
        e.activeEmails[u.Email] = u.ID
    }
}
```

When a user is soft-deleted via `SoftDeleteUser`, its entry is removed from `activeEmails`, freeing the email for future active registrations while preserving the historical deleted record.

## 6. Concurrency Testing and Empirical Verification

The lab includes a test suite in `internal/store/store_test.go` verifying the behavior under parallel execution.

### The Race Condition in UnsafeStore

```go
func TestConcurrentRegistration_Unsafe_SuffersRaceCondition(t *testing.T) {
    eng := engine.NewEngine()
    store := store.NewUnsafeStore(eng)

    var wg sync.WaitGroup
    email := "race-unsafe@example.com"

    for i := 0; i < 20; i++ {
        wg.Add(1)
        go func(idx int) {
            defer wg.Done()
            _ = store.RegisterUser(model.User{
                Email:    email,
                Username: fmt.Sprintf("user-%d", idx),
                Age:      20,
                Status:   "active",
            })
        }(i)
    }
    wg.Wait()

    // Assert that race condition corrupted data by inserting duplicates
    require.Greater(t, eng.GetUsersCount(), 1)
}
```

When run, `eng.GetUsersCount()` consistently evaluates to greater than 1 (often 5 to 15 duplicate rows), proving that application-level checking fails under concurrency.

### Absolute Integrity in SafeStore

```go
func TestConcurrentRegistration_Safe_EnforcesUniqueness(t *testing.T) {
    eng := engine.NewEngine()
    store := store.NewSafeStore(eng)

    var wg sync.WaitGroup
    email := "race-safe@example.com"
    successCount := atomic.Int64{}
    violationCount := atomic.Int64{}

    for i := 0; i < 20; i++ {
        wg.Add(1)
        go func(idx int) {
            defer wg.Done()
            err := store.RegisterUser(model.User{
                Email:    email,
                Username: fmt.Sprintf("user-%d", idx),
                Age:      20,
                Status:   "active",
            })
            if err == nil {
                successCount.Add(1)
            } else if dberr.IsConstraintViolation(err, dberr.SQLStateUniqueViolation) {
                violationCount.Add(1)
            }
        }(i)
    }
    wg.Wait()

    require.Equal(t, int64(1), successCount.Load())
    require.Equal(t, int64(19), violationCount.Load())
    require.Equal(t, 1, eng.GetUsersCount())
}
```

**Outcome**: Exactly 1 goroutine succeeds, 19 goroutines are rejected with `SQLSTATE 23505`, and the total user count in storage is strictly 1.

Running the full suite with the Go race detector:

```bash
go test -v -race ./...
```

Output:
```text
=== RUN   TestNotNullConstraints
--- PASS: TestNotNullConstraints (0.00s)
=== RUN   TestCheckConstraints
--- PASS: TestCheckConstraints (0.00s)
=== RUN   TestUniqueConstraint
--- PASS: TestUniqueConstraint (0.00s)
=== RUN   TestForeignKeyConstraint
--- PASS: TestForeignKeyConstraint (0.00s)
=== RUN   TestPartialUniqueIndex
--- PASS: TestPartialUniqueIndex (0.00s)
=== RUN   TestConcurrentRegistration_Safe_EnforcesUniqueness
--- PASS: TestConcurrentRegistration_Safe_EnforcesUniqueness (0.00s)
=== RUN   TestConcurrentRegistration_Unsafe_SuffersRaceCondition
--- PASS: TestConcurrentRegistration_Unsafe_SuffersRaceCondition (0.00s)
=== RUN   TestErrorClassification
--- PASS: TestErrorClassification (0.00s)
PASS
ok      github.com/lukman/software-engineering-lab/labs/27-database-constraints/internal/store  1.024s
```

## 7. Error Classification and Domain Mapping

Applications should never inspect raw error strings like `"duplicate key"` using regular expressions. Error strings vary between database versions and change when database locale settings switch (e.g., from English to German or Indonesian).

Instead, applications must inspect ANSI SQLSTATE codes. In `internal/dberr/errors.go`:

```go
func IsConstraintViolation(err error, code SQLState) bool {
    var cErr *ConstraintError
    if errors.As(err, &cErr) {
        return cErr.Code == code
    }
    return false
}

func MapToDomainError(err error) error {
    var cErr *ConstraintError
    if errors.As(err, &cErr) {
        switch cErr.Code {
        case SQLStateUniqueViolation:
            return fmt.Errorf("conflict: resource with this unique attribute already exists (rule: %s)", cErr.ConstraintName)
        case SQLStateNotNullViolation:
            return fmt.Errorf("invalid input: mandatory field is missing (rule: %s)", cErr.ConstraintName)
        case SQLStateCheckViolation:
            return fmt.Errorf("validation failed: value outside permissible boundary (rule: %s)", cErr.ConstraintName)
        case SQLStateForeignKeyViolation:
            return fmt.Errorf("reference error: referenced entity does not exist (rule: %s)", cErr.ConstraintName)
        default:
            return fmt.Errorf("database integrity violation: %s", cErr.Detail)
        }
    }
    return err
}
```

In an HTTP API context, this allows seamless, deterministic status code translation:
- `23505 Unique Violation` → `HTTP 409 Conflict`
- `23502 Not Null Violation` → `HTTP 400 Bad Request`
- `23514 Check Violation` → `HTTP 422 Unprocessable Entity`
- `23503 Foreign Key Violation` → `HTTP 409 Conflict`

## 8. Before vs After Comparison

| Evaluation Metric | Application Validation Only (`UnsafeStore`) | Database Constraints (`SafeStore`) |
|---|---|---|
| **Concurrency Protection** | Vulnerable to check-then-act race conditions | 100% atomic enforcement under storage engine lock |
| **Outcome on 50 Concurrent Inserts** | Multiple corrupt duplicate rows created | Exactly 1 success, 49 rejected with SQLSTATE 23505 |
| **Bypass Vulnerability** | Direct SQL, migrations, and other services bypass checks | Universal enforcement regardless of client origin |
| **Soft Delete Support** | Complex manual queries with potential race windows | Partial Unique Index (`WHERE deleted_at IS NULL`) |
| **Error Handling Method** | Brittle regex string matching on error text | Programmatic classification via ANSI SQLSTATE Class 23 |
| **Architectural Role** | First line of defense (User Experience) | Last line of defense (Absolute Data Correctness) |

## 9. Production Considerations and Limitations

When translating the findings of this lab into high-throughput production environments, consider the following technical boundaries:

1. **Foreign Key Indexing**: In PostgreSQL, creating a `FOREIGN KEY` does not automatically create an index on the referencing table's column (e.g., `orders.user_id`). Without a dedicated index, executing a `DELETE` or `UPDATE` on the parent table (`users`) triggers a full sequential scan on the referencing table (`orders`), severely impacting database throughput.
2. **Row-Scoped Limitation**: `CHECK` constraints can only inspect columns within the single row being modified. Cross-row invariants (such as "total user balance across accounts must equal zero" or "no more than 3 active trials per user") cannot be enforced via standard CHECK constraints. They require PostgreSQL triggers or `SERIALIZABLE` transaction isolation.
3. **Zero-Downtime Migrations**: Adding a constraint to a large production table can lock the table for hours. In PostgreSQL, use the two-step migration pattern:
   ```sql
   -- Step 1: Add constraint without validating historical rows (fast, brief lock)
   ALTER TABLE orders ADD CONSTRAINT fk_orders_user 
     FOREIGN KEY (user_id) REFERENCES users (id) NOT VALID;

   -- Step 2: Validate existing rows in background (no write lock)
   ALTER TABLE orders VALIDATE CONSTRAINT fk_orders_user;
   ```
4. **Lab Simulator Boundary**: The lab engine is an educational in-memory simulator written in Go. It uses coarse-grained `sync.RWMutex` locking to model PostgreSQL atomic semantics cleanly without requiring a live Docker container. It is not intended for benchmarking disk I/O or buffer pool latches.

## 10. Summary and Takeaways

- **Eliminate Check-Then-Act**: Never depend on `SELECT` followed by `INSERT` to prevent duplicates. Let the database constraint decide atomically.
- **Adopt Defense in Depth**: Use application validation to craft user-friendly error messages and save round-trips; rely on database constraints for absolute correctness.
- **Use Partial Unique Indexes**: Solve the soft-delete reuse dilemma cleanly with `WHERE deleted_at IS NULL`.
- **Inspect SQLSTATE Codes**: Inspect codes (`23505`, `23502`, `23514`, `23503`) rather than parsing localized error strings.
- **Always Test Concurrency**: Sequential unit tests hide race conditions. Test with goroutines and the Go race detector (`go test -race ./...`).

## Repository and Code Reference

- **Lab Source Code**: [https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/27-database-constraints](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/27-database-constraints)
- **Commit SHA**: `bdb285d35a71c9a309e83a0e4497b9dbd842ceb8`
- **Verification**: 8/8 tests passing with `go test -v -race ./...`
