---
title: "Database Constraints: The Last Line of Data Integrity"
published: false
description: "Prevent race conditions, enforce schema invariants, and reject duplicates atomically in Go — with SQLSTATE taxonomy, partial indexes, and 50-goroutine stress tests."
tags: go,database,integrity,concurrency,sqlstate
---

# Database Constraints: The Last Line of Data Integrity

Many backend architectures delegate data integrity entirely to application code. Under concurrent load, naive validation lets duplicate records and orphaned relations slip into persistent storage. Lab 27 empirically demonstrates how database constraints act as the unbreachable, atomic last line of defense — using a pure Go in-memory simulator under 50 concurrent goroutines.

> Database constraints provide declarative, ACID-guaranteed enforcement that eliminates check-then-act vulnerabilities.

## Real-World Problem: Application Validation Collapses Under Concurrency

The check-then-act pattern inspecting records via SELECT before issuing an INSERT invariably succeeds in sequential testing. Under concurrent traffic, however, two threads observe the same pre-write state before either commits. This failure does not stem from faulty conditional logic, but from lacking atomic isolation at the persistence tier.

Duplicate detection via application-level memory checks leaves a race window that CPU scheduling can exploit. Manual deduplication in production demands reconciliation, complex migrations, and downtime. Unit tests without concurrency are misleading — they create an illusion of safety.

## Why Manual Review and Naive Approaches Fail

Code reviews routinely miss integrity vulnerabilities because application-level checks appear impeccably structured on screen. Sequential unit tests pass reliably, fostering false confidence across engineering teams. Running validations purely in application services does not automatically ensure integrity without atomic database constraints.

Three specific weaknesses exist:

1. **CPU scheduling interleaving** — milliseconds between SELECT and INSERT are enough for another thread to interleave.
2. **Sequential unit tests deceive** — 100% pass rates with no goroutines create a false sense of correctness.
3. **No automatic protection** — running validation in a separate microservice still lacks atomic locking.

## Two-Layer Defense: Application and Database Collaboration

Modern software architectures divide integrity responsibilities into two layers: the application handles user experience, while the database guarantees absolute correctness.

| Layer | Role | Responsibility |
|---|---|---|
| Application Validation | User Experience (UX) | Friendly error messages, fast format validation, input sanitization before the request reaches the database |
| Database Constraints | Last Line of Truth | Declarative, atomic schema invariant enforcement with internal engine locks at write time |
| Error Mapping Bridge | SQLSTATE Taxonomy | Maps database error 23xxx into structured domain errors like HTTP 409 Conflict |

Application validation filters invalid formats swiftly. Only along the path where concurrent collisions occur does the database reject violations atomically with standard SQLSTATE codes.

## Constraint Evaluation Pipeline: Deterministic Internal Flow

The storage engine evaluates schema constraints in a deterministic sequence:

```
NOT NULL (23502) → CHECK (23514) → UNIQUE / PARTIAL (23505) → FOREIGN KEY (23503) → ATOMIC COMMIT
```

1. **NOT NULL** — checks column emptiness; fails immediately if a mandatory field is empty.
2. **CHECK** — evaluates boolean expressions on the row (e.g., `age >= 18`).
3. **UNIQUE / PARTIAL** — scans B-tree indexes under exclusive lock protection.
4. **FOREIGN KEY** — verifies parent record existence before commit.
5. **ATOMIC COMMIT** — writes row data and index entries simultaneously.

If any check fails, execution aborts immediately without mutating tables or indexes. This lab models this sequential evaluation atomically within a unified critical section.

## Atomic Constraint Evaluation Under Mutex Lock

The `InsertUser` method executes all constraint checks inside an exclusive mutex lock representing PostgreSQL's ROW EXCLUSIVE lock:

```go
func (e *Engine) InsertUser(u model.User, usePartial bool) (model.User, error) {
    e.mu.Lock()
    defer e.mu.Unlock()
    if u.Email == "" { return model.User{}, dberr.NewNotNullViolation("users", "email", "users_email_not_null") }
    if u.Age < 18 { return model.User{}, dberr.NewCheckViolation("users", "users_age_check", "age >= 18") }
    if exists := e.emailIndex[u.Email]; exists {
        return model.User{}, dberr.NewUniqueViolation("users", "users_email_key", "duplicate email")
    }
    u.ID = e.userSeq.Add(1); e.users[u.ID] = u; e.emailIndex[u.Email] = u.ID
    return u, nil
}
```

Locking the critical section models atomic guarantees cleanly — it does not implement low-level B-tree page latches, but ensures no competing goroutine can insert a conflicting email during evaluation.

## Race Condition Chronicle: Two Goroutines on Naive Validation

When two goroutines execute `UnsafeStore.RegisterUser` concurrently, both observe that the requested email does not yet exist. CPU scheduling delays expose a vulnerability window, prompting both threads to issue INSERT calls.

The chronology:

1. **Concurrent check** — Goroutine A and B both scan the table and read that the email is absent.
2. **CPU scheduling gap** — Both enter a scheduling delay (`time.Sleep`) exposing the race window.
3. **Dual write success** — Goroutine A writes the row; Goroutine B immediately follows with the same email.
4. **Integrity collapse** — The table now holds 2 rows with identical email, violating the uniqueness invariant.

Uncoordinated concurrent reads before unconstrained writes guarantee permanent data corruption.

## Taxonomy Matrix: Constraint Types vs Characteristics

This matrix charts the five constraint mechanisms verified in our lab alongside standard PostgreSQL Class 23 SQLSTATE codes:

| Constraint Type | SQLSTATE | Evaluation Scope | Failure Prevention | HTTP Domain Mapping |
|---|---|---|---|---|
| NOT NULL | 23502 | Single Column | Prevents identity-less empty records | 400 Bad Request |
| CHECK | 23514 | Column / Row | Rejects values outside business-logical bounds | 422 Unprocessable Entity |
| UNIQUE | 23505 | Cross-Row (Index) | Prevents duplicate registration race condition | 409 Conflict |
| FOREIGN KEY | 23503 | Cross-Table | Prevents orphan records (missing parent) | 400 Reference Error |
| PARTIAL UNIQUE | 23505 | Conditional Rows | Allows email reuse after soft-delete | 409 Conflict |

Each mechanism addresses a distinct integrity dimension, ranging from single-column presence to cross-table referential integrity.

## Three Key Decisions for Schema Integrity Design

Robust schema design begins by defining mandatory fields with NOT NULL and value boundaries with CHECK. When managing unique entities that support soft deletion, use a partial unique index so recycled identifiers remain valid. Adding foreign keys does not mean the job is complete — referenced foreign key columns require manual indexing to sustain query performance.

1. **Mandatory & Row Predicates** — Apply NOT NULL and CHECK on every column with an explicit domain boundary.
2. **Uniqueness & Deletion Pattern** — Choose standard UNIQUE for permanent entities, or Partial Index for soft-delete systems.
3. **Relational Integrity & Indexing** — Install FOREIGN KEY to prevent orphan rows; manual indexing on referencing columns is still required.

## Execution Evidence: Test Suite and 50-Goroutine Stress Test

Execution output confirms that all test suites pass with zero warnings from the Go race detector:

```
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
=== RUN   TestErrorClassification
--- PASS: TestErrorClassification (0.00s)
PASS: ok  labs/27-database-constraints/internal/store  1.330s (race: clean)
[DEMO] CONCURRENCY STRESS TEST: 50 GOROUTINES FOR SAME EMAIL
Results: Total=50 | Successful=1 | Rejected 23505=49 | Integrity Intact=true
```

The stress test features 50 concurrent goroutines submitting an identical email. The storage engine admitted exactly 1 registration while rejecting 49 attempts with error 23505. This confirms 100% data integrity under sustained parallel load.

## Evidence Boundary: What Is and Is Not Proven

This laboratory has definitive limitations that must be understood to prevent overclaiming.

**Proven:** SQLSTATE 23xxx semantics, race condition prevention via atomic lock, and the lifecycle of partial indexes — all using a pure Go in-memory simulator.

**Not Proven:** Distributed multi-node constraints, SERIALIZABLE isolation, GiST exclusion operators, real B-tree latches, and `ALTER TABLE ADD CONSTRAINT NOT VALID` migration techniques on massive production tables.

The test suite has not been verified against actual PostgreSQL or any distributed storage engine. Maintaining this empirical boundary ensures engineers do not conflate simulator testing with production database benchmarking.

## Engineering Trade-offs: Constraints vs Schema Flexibility

Enforcing database constraints introduces minor write overhead during INSERTs as the engine inspects B-tree indexes and parent records. Furthermore, attaching constraints to massive production tables demands phased migration workflows to prevent prolonged table locks.

However, this overhead does not degrade overall throughput when weighed against the catastrophic cost of data repair. The cost model is:

- **Write overhead** — every INSERT/UPDATE must evaluate B-tree indexes and FK relationships, adding marginal latency.
- **Migration friction** — adding constraints on billion-row tables requires exclusive table locks unless phased validation patterns are used.
- **Not inherently slow** — constraint index additions do not automatically degrade system performance; write overhead is repaid by guaranteed integrity.

## Concurrent Pipeline: Goroutines vs Mutex Lock

This diagram visualizes how the storage engine handles 50 concurrent goroutines racing to register an identical email:

1. **50 Goroutines Spawn** — 50 workers run simultaneously attempting to register `race-target@domain.com`.
2. **Mutex Lock Serialization** — The first goroutine acquires the exclusive lock; 49 workers queue orderly.
3. **First Insert Success** — Worker 1 passes checks, commits the email to the table and index map, then releases the lock.
4. **49 Abort SQLSTATE 23505** — Workers 2–50 enter the lock sequentially, detect the email already exists, and receive error 23505.

Atomic serialization at the write boundary transforms concurrent race conditions into deterministic rejections. The mutex ensures only one writer succeeds; all competitors receive clean constraint violations.

## Transformation: UnsafeStore to SafeStore

Prior to implementing database constraints, UnsafeStore relied on in-memory loops that collapsed under 20 concurrent goroutines, allowing duplicate rows. Migrating to SafeStore delegates integrity to the storage engine, guaranteeing exactly one winner among 50 parallel requests. Retaining application validation upfront does not represent redundancy — it provides an essential bridge for friendly client feedback.

| Before: UnsafeStore (App Check) | After: SafeStore (DB Constraint) |
|---|---|
| Relies on in-memory loop scanning | Delegates to storage engine |
| 20 goroutines produce >1 duplicate user | 50 goroutines produce exactly 1 user + 49 SQLSTATE 23505 errors |
| No atomicity guarantee | Atomic constraint enforcement |

Retaining application validation serves as the initial filter for user-facing error messages; the database constraint serves as the final, unyielding gatekeeper.

## Deep Analysis: UnsafeStore Race Condition Test

This unit test empirically demonstrates that application-level validation collapses under concurrent workloads:

```go
func TestConcurrentRegistration_Unsafe_SuffersRaceCondition(t *testing.T) {
    eng := engine.NewEngine(); s := store.NewUnsafeStore(eng)
    goroutines := 20; var wg sync.WaitGroup; wg.Add(goroutines)
    for i := 0; i < goroutines; i++ {
        go func(idx int) {
            defer wg.Done()
            _, _ = s.RegisterUser(ctx, model.User{Email: "unsafe@example.com", Username: fmt.Sprintf("u_%d", idx)})
        }(i)
    }
    wg.Wait()
    if eng.GetUsersCount() <= 1 { t.Errorf("expected race condition causing >1 user entries, got %d", eng.GetUsersCount()) }
}
```

Twenty goroutines submitting an identical email produce multiple duplicate records in the underlying table. The scope of this failure aligns strictly with lab tests, proving that systems lacking database constraints suffer inevitable data corruption under traffic spikes.

## Deep Analysis: SafeStore Uniqueness Enforcement Test

This test verifies the resilience of SafeStore backed by engine constraints:

```go
func TestConcurrentRegistration_Safe_EnforcesUniqueness(t *testing.T) {
    eng := engine.NewEngine(); s := store.NewSafeStore(eng)
    goroutines := 20; var wg sync.WaitGroup; wg.Add(goroutines)
    successCount, errorCount := 0, 0; var mu sync.Mutex
    for i := 0; i < goroutines; i++ {
        go func(idx int) {
            defer wg.Done()
            _, err := s.RegisterUser(ctx, model.User{Email: "safe@example.com", Username: fmt.Sprintf("u_%d", idx), Age: 25, Status: "active"})
            mu.Lock(); defer mu.Unlock()
            if err == nil { successCount++ } else { errorCount++ }
        }(i)
    }
    wg.Wait()
    if successCount != 1 || errorCount != 19 || eng.GetUsersCount() != 1 { t.Fatalf("integrity violated!") }
}
```

Among 20 concurrently spawned goroutines, exactly one registration succeeds in persisting to storage, while the remaining 19 receive immediate rejections. Table integrity is strictly preserved at exactly 1 row, confirming the power of database-level protection.

## Deep Analysis: Partial Unique Index and Soft Delete

Applying a standard UNIQUE constraint to tables with soft deletion creates severe friction: deactivated emails can never be re-registered. By leveraging a conditional Partial Unique Index with `WHERE deleted_at IS NULL`, the database enforces uniqueness exclusively on active rows.

```go
// PARTIAL UNIQUE: CREATE UNIQUE INDEX users_active_email_idx ON users (email) WHERE deleted_at IS NULL;
if usePartialUniqueIndex {
    if u.DeletedAt == nil {
        if existingID, exists := e.activeEmails[u.Email]; exists {
            return model.User{}, dberr.NewUniqueViolation("users", "users_active_email_idx", "duplicate active email")
        }
    }
}
// On soft delete: delete(e.activeEmails, u.Email) -> email can be re-registered immediately
// Rows with deleted_at != nil are not indexed, so multi-duplicate soft-deleted rows remain valid.
```

Deactivated records are evicted from the index, allowing identifier recycling without conflict. Only active rows participate in the uniqueness predicate.

## Deep Analysis: Foreign Key Integrity — Preventing Orphan Records

The `InsertOrder` function sequentially evaluates three constraints:

```go
func (e *Engine) InsertOrder(o model.Order) (model.Order, error) {
    e.mu.Lock(); defer e.mu.Unlock()
    if o.UserID == 0 { return model.Order{}, dberr.NewNotNullViolation("orders", "user_id", "orders_user_id_not_null") }
    if o.TotalCents <= 0 { return model.Order{}, dberr.NewCheckViolation("orders", "orders_total_cents_check", "total > 0") }
    // FOREIGN KEY: CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES users(id)
    if _, exists := e.users[o.UserID]; !exists {
        return model.Order{}, dberr.NewForeignKeyViolation("orders", "fk_orders_user", fmt.Sprintf("key (user_id)=(%d) not present", o.UserID))
    }
    o.ID = e.orderSeq.Add(1); e.orders[o.ID] = o; return o, nil
}
```

If a UserID of 99999 is submitted without existing in the users table, the engine returns SQLSTATE 23503. This mechanism prevents orphaned order records that could compromise financial reconciliation.

## Deep Analysis: SQLSTATE Taxonomy and Domain Error Mapping

Many engineers parse raw error strings like "duplicate key," which break whenever database locales change. Class 23 SQLSTATE codes provide a standardized, permanent machine contract:

```go
func MapToDomainError(err error) error {
    var cErr *ConstraintError
    if errors.As(err, &cErr) {
        switch cErr.Code {
        case SQLStateUniqueViolation: // 23505
            return fmt.Errorf("conflict: resource with this unique attribute already exists (rule: %s)", cErr.ConstraintName)
        case SQLStateNotNullViolation: // 23502
            return fmt.Errorf("invalid input: mandatory field is missing (rule: %s)", cErr.ConstraintName)
        case SQLStateCheckViolation: // 23514
            return fmt.Errorf("validation failed: value outside boundary (rule: %s)", cErr.ConstraintName)
        case SQLStateForeignKeyViolation: // 23503
            return fmt.Errorf("reference error: referenced entity does not exist (rule: %s)", cErr.ConstraintName)
        }
    }
    return err
}
```

Inspecting stable SQLSTATE codes rather than brittle localized text strings enables robust application error mapping without leaking internal database schemas to clients.

## Production Implementation Checklist

This checklist summarizes mandatory engineering steps before releasing schemas to production:

- Mark all mandatory columns with NOT NULL; do not rely solely on form validation.
- Apply UNIQUE constraints on unique attributes instead of SELECT-before-INSERT patterns.
- Use Partial Unique Index (`WHERE deleted_at IS NULL`) for entities supporting soft deletion.
- Define CHECK constraints to limit numeric ranges and enum status columns.
- Install explicit indexes on FOREIGN KEY referencing columns to prevent full table scans on parent DELETE.
- Inspect Class 23 SQLSTATE codes (23502, 23505, 23514, 23503) in application error handlers.

## Common Anti-Patterns in Data Integrity Maintenance

The most perilous anti-pattern is assuming safety from application-level SELECT-then-INSERT flows. Other prevalent traps include parsing localized error strings rather than SQLSTATE codes and assuming foreign keys automatically index referencing columns.

| Anti-Pattern | Risk |
|---|---|
| Check-Then-Act in Application | SELECT before INSERT always leaks under concurrent load |
| Guessing Error Strings | Regex matching on DB error text breaks across versions and locales |
| Assuming Auto-Indexing | FOREIGN KEY creation does not automatically index the referencing column |

Overlooking these realities creates ticking performance and data corruption time-bombs.

## Self-Directed Reproduction Commands

The entire lab simulation is engineered strictly with the Go standard library, requiring no external packages or Docker containers:

```sh
cd labs/27-database-constraints
go test -v ./...     # runs 7 unit tests for NOT NULL, CHECK, UNIQUE, FK, and Partial Index
go test -race ./...  # validates concurrent execution is free of data races (PASS in 1.33s)
go run ./cmd/demo    # runs interactive CLI demo: SQLSTATE demo & 50-worker stress test
```

While local tests prove constraint semantics deterministically, production environments still require phased migration strategies such as `ADD CONSTRAINT NOT VALID`.

## Long-Term Architectural Principles

Application layers will inevitably be refactored and microservices will proliferate, but persisted data outlives software lifecycles. Entrusting data integrity solely to application code is a high-risk architectural gamble.

| Pillar | Description |
|---|---|
| Declarative at Write Point | Place data validity rules as close to the storage medium as possible for atomic protection |
| Clear Role Separation | Application serves user interaction comfort; database constraints guarantee absolute truth |
| Not Automatically Immune | Installing constraints does not automatically resolve cross-table rules; understand row-scoped vs trigger boundaries |

Database constraints ensure that regardless of application-layer regressions, persistent data remains uncorrupted.

## Source

[Software Engineering Lab 27 — Database Constraints](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/27-database-constraints), snapshot commit `c317c1d`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
