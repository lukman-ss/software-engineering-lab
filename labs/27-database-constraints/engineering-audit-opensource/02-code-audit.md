## Finding 1

Location: engine.go:23-30 (NewEngine)
Claimed Behavior: Engine indexes correctly initialized.
Observed Implementation: emailIndex (full unique) and activeEmails (partial unique) both map[string]int64.
Assessment: PASS
Severity: LOW
Notes: Both maps are properly initialized for uniqueness enforcement.

## Finding 2

Location: engine.go:45-98 (InsertUser)
Claimed Behavior: NOT NULL, CHECK, UNIQUE/PARTIAL UNIQUE enforced atomically under engine.mu.
Observed Implementation: All validations and index updates occur under a single RWMutex write lock. Indexes updated only after validation passes.
Assessment: PASS
Severity: LOW
Notes: Correctly prevents partial updates on failure.

## Finding 3

Location: engine.go:51-56, 54-56 (NOT NULL)
Claimed Behavior: Email and username NOT NULL enforced (23502).
Observed Implementation: Blank string check returns dberr.NewNotNullViolation with correct table/column/constraint.
Assessment: PASS
Severity: LOW
Notes: Matches claim.

## Finding 4

Location: engine.go:59-68 (CHECK age >= 18), engine.go:63-68 (CHECK status IN)
Claimed Behavior: CHECK constraints reject age<18 and invalid status.
Observed Implementation: Exactly as claimed; error detail includes value and expression.
Assessment: PASS
Severity: LOW
Notes: Matches claim.

## Finding 5

Location: engine.go:70-83 (UNIQUE/PARTIAL UNIQUE logic)
Claimed Behavior: Standard UNIQUE always enforced; partial unique only for rows with deleted_at IS NULL.
Observed Implementation: Branch on usePartialUniqueIndex; activeEmails updated only when DeletedAt==NULL.
Assessment: PASS
Severity: LOW
Notes: Correct conditional index semantics.

## Finding 6

Location: engine.go:80, 81, 75, 76 (Unique violation detail)
Claimed Behavior: Unique violation includes duplicate key value and existing_id.
Observed Implementation: Format strings match: duplicate key value violates unique constraint (email=%s, existing_id=%d).
Assessment: PASS
Severity: LOW
Notes: Detail string matches PostgreSQL style.

## Finding 7

Location: engine.go:122-148 (InsertOrder)
Claimed Behavior: NOT NULL on user_id, CHECK total_cents>0, FOREIGN KEY users(id).
Observed Implementation: All three checks in order; foreign key uses e.users map lookup.
Assessment: PASS
Severity: LOW
Notes: FK enforced via in-memory parent table presence check.

## Finding 8

Location: engine.go:101-120 (SoftDeleteUser)
Claimed Behavior: Soft delete removes email from partial unique active index.
Observed Implementation: Deletes from activeEmails before setting DeletedAt; does not touch full emailIndex.
Assessment: PASS
Severity: LOW
Notes: Correct partial-index maintenance.

## Finding 9

Location: store.go:15-21 (UnsafeStore)
Claimed Behavior: UnsafeStore does app-level duplicate check then inserts unsafe.
Observed Implementation: RegisterUser reads all users via GetUsersCount/GetUser loop, then sleeps, then InsertUserUnsafe.
Assessment: PASS
Severity: MEDIUM
Notes: The check is correct but extremely narrow window; relies on sleep to expose race. Demonstrates concept but not a strong data race.

## Finding 10

Location: store.go:24-47 (UnsafeStore.RegisterUser)
Claimed Behavior: App-level check vulnerable to read-then-write race.
Observed Implementation: Loop over integer IDs assumes dense packing; real holes would cause false negatives. Sleep(1ms) exacerbates window.
Assessment: WARNING
Severity: MEDIUM
Notes: Implementation approximates vulnerability but is not a true data race on shared state; the engine itself is locked during insert, so duplicates only possible if two goroutines both pass the check before either inserts. Acceptable for demo.

## Finding 11

Location: store.go:58-66 (SafeStore.RegisterUser)
Claimed Behavior: Delegates to engine.InsertUser (constraint-enforced).
Observed Implementation: Direct call, then MapToDomainError.
Assessment: PASS
Severity: LOW
Notes: Correct.

## Finding 12

Location: store.go:68-75 (SafeStore.RegisterUserPartial)
Claimed Behavior: Delegates to engine.InsertUser with usePartialUniqueIndex=true.
Observed Implementation: Matches.
Assessment: PASS
Severity: LOW
Notes: Correct.

## Finding 13

Location: store.go:77-83 (SafeStore.CreateOrder)
Claimed Behavior: Delegates to engine.InsertOrder.
Observed Implementation: Matches.
Assessment: PASS
Severity: LOW
Notes: Correct.

## Finding 14

Location: dberr/errors.go:39-45, 48-54, 57-63, 66-72 (New*Violation)
Claimed Behavior: ConstraintError factory functions set correct SQLSTate and detail.
Observed Implementation: All return *ConstraintError with matching Code, ConstraintName, TableName, Detail.
Assessment: PASS
Severity: LOW
Notes: Matches PostgreSQL SQLSTATE class 23.

## Finding 15

Location: dberr/errors.go:83-99 (MapToDomainError)
Claimed Behavior: Maps SQLSTATE 23502/23503/23505/23514 to domain-friendly messages preserving constraint name.
Observed Implementation: Switch returns fmt.Errorf with rule: %s using cErr.ConstraintName.
Assessment: PASS
Severity: LOW
Notes: Correctly preserves constraint name in user message.

## Finding 16

Location: model/model.go:5-12, 14-19 (User, Order structs)
Claimed Behavior: Domain models match engine expectations.
Observed Implementation: User.Email string, Username string, Age int, Status string, DeletedAt *time.Time; Order.UserID int64, TotalCents int64.
Assessment: PASS
Severity: LOW
Notes: Types align.

## Finding 17

Location: cmd/demo/main.go:15-13 (imports)
Claimed Behavior: Demo imports correct packages.
Observed Implementation: Imports internal/engine, model, store, dberr.
Assessment: PASS
Severity: LOW
Notes: Correct.