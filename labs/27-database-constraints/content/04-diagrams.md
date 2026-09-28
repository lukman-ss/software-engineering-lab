## Diagram 1 — Pipeline Position

```text
Research
  ↓
Research Audit → APPROVED
  ↓
Engineering
  ↓
Engineering Audit → APPROVED
  ↓
Technical Writer (current)
  ↓
READY_FOR_CONTENT_ADAPTER
```

---

## Diagram 2 — Race Condition: Check-Then-Act Failure

Source: `research/05-report.md`, reproduced in lab via `UnsafeStore`

```text
Goroutine A                    Goroutine B
  │                               │
  ├─ scan → 0 rows found          │
  │                               ├─ scan → 0 rows found
  │                               │
  ├─ sleep(1ms)                   │
  │                               ├─ sleep(1ms)
  │                               │
  ├─ INSERT (email=X) → OK        │
  │                               ├─ INSERT (email=X) → OK (BUG: duplicate)
  │                               │
  ├─ count() → 2                  │ ← integrity violated
```

`UnsafeStore.RegisterUser` (`internal/store/store.go:24-47`) does scan → sleep → insert without atomic uniqueness check. Both goroutines see `count == 0` before either inserts.

---

## Diagram 3 — Safe Insert: Atomic Constraint Enforcement

Source: `internal/engine/engine.go:45-99`

```text
Goroutine A                    Goroutine B
  │                               │
  ├─ Lock acquired                │
  ├─ NOT NULL check               │
  ├─ CHECK validation             │
  ├─ UNIQUE check: no conflict    │
  ├─ COMMIT row + index           │
  ├─ Unlock                       │
  │                               ├─ Lock acquired
  │                               ├─ NOT NULL check
  │                               ├─ CHECK validation
  │                               ├─ UNIQUE check: CONFLICT (23505)
  │                               ├─ Return error immediately
  │                               ├─ Unlock
  │                               │
  ├─ count() → 1                  │ ← integrity intact
```

Single mutex. Second goroutine sees committed index entry before it can proceed past UNIQUE check. This mirrors PostgreSQL's `ROW EXCLUSIVE` lock + atomic B-tree index update during INSERT.

---

## Diagram 4 — Constraint Evaluation Order

Source: `internal/engine/engine.go:45-99, 122-148`

```text
INSERT operation
  │
  ├──[1] NOT NULL (23502)
  │     Column null or empty → reject
  │
  ├──[2] CHECK (23514)
  │     Boolean predicate false → reject
  │
  ├──[3] UNIQUE (23505) or PARTIAL UNIQUE
  │     Duplicate key detected → reject
  │
  ├──[4] FOREIGN KEY (23503)
  │     Referenced row missing → reject
  │
  └──[5] COMMIT
        Row + index written
```

Checks execute sequentially. First failure returns error immediately. No row or index entry is written if any check fails. This ordering matches PostgreSQL behavior: NOT NULL first, then CHECK (alphabetical by name), then FK/UNIQUE (index-based).

---

## Diagram 5 — Partial Unique Index: Soft-Delete Flow

Source: `internal/engine/engine.go:45-99, 101-120`

```text
activeEmails map (WHERE deleted_at IS NULL):
  ┌──────────────────────────────────────────────┐
  │ alice@company.com → ID=1 (active)            │
  └──────────────────────────────────────────────┘

Step 1: Insert active user with email alice@company.com
  ├── deleted_at IS NULL → check activeEmails
  ├── No conflict → INSERT, add to activeEmails
  └── activeEmails: {alice@company.com: ID=1}

Step 2: Insert duplicate active user (same email)
  ├── deleted_at IS NULL → check activeEmails
  ├── Conflict found → ERROR 23505

Step 3: Soft delete user ID=1
  ├── delete(activeEmails, alice@company.com)
  └── activeEmails: {}

Step 4: Insert active user with same email (ID=2)
  ├── deleted_at IS NULL → check activeEmails
  ├── No conflict (alice@company.com removed) → INSERT
  └── activeEmails: {alice@company.com: ID=2}
```

Partial index `WHERE deleted_at IS NULL` only indexes active rows. After soft delete, email is removed from index, allowing reuse. Multiple soft-deleted rows with same email coexist because they are not indexed by partial index.

---

## Diagram 6 — Architecture

Source: `internal/store/store.go`, `internal/engine/engine.go`, `internal/dberr/errors.go`, `internal/model/model.go`, `cmd/demo/main.go`

```text
cmd/demo/main.go
  │
  └── Store ─────────────────────────────────────────────────────┐
        │                                                        │
        ├── UnsafeStore                                          │
        │     ├── scan users (no lock)                           │
        │     ├── sleep(1ms)                                     │
        │     └── InsertUserUnsafe (no index check)              │
        │                                                        │
        └── SafeStore                                            │
              ├── RegisterUser → InsertUser(false)               │
              ├── RegisterUserPartial → InsertUser(true)         │
              └── CreateOrder → InsertOrder                      │
                     │                                           │
                     ├── Engine (sync.RWMutex)                   │
                     │     ├── users map (table)                 │
                     │     ├── orders map (table)                │
                     │     ├── emailIndex map (UNIQUE)           │
                     │     └── activeEmails map (PARTIAL UNIQUE) │
                     │                                           │
                     └── dberr (SQLSTATE taxonomy)                │
                           ├── MapToDomainError                  │
                           ├── IsConstraintViolation             │
                           └── ConstraintError                   │
```

---

## Diagram 7 — Error Code Mapping

Source: `internal/dberr/errors.go:83-100`

```text
SQLSTATE    Constraint Name              Domain Error (HTTP)
────────    ─────────────────            ──────────────────────
23502       users_email_not_null         invalid input: mandatory field is missing
23505       users_email_key              conflict: resource with this unique attribute already exists
23514       users_age_check              validation failed: value outside permissible boundary
23503       fk_orders_user               reference error: referenced entity does not exist
```

All messages include constraint/rule name for debugging. `IsConstraintViolation` inspects `ConstraintError.Code` against target SQLSTATE for programmatic error routing.