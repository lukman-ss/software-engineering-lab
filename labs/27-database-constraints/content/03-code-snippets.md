## Snippet 1 — SQLSTATE Error Codes and Taxonomy

Source File: `internal/dberr/errors.go:10-18`
Purpose: Defines ANSI SQLSTATE class 23 error codes for constraint violations; enables programmatic error inspection instead of parsing error text.

```go
const (
	SQLStateRestrictViolation   SQLState = "23001"
	SQLStateNotNullViolation    SQLState = "23502"
	SQLStateForeignKeyViolation SQLState = "23503"
	SQLStateUniqueViolation     SQLState = "23505"
	SQLStateCheckViolation      SQLState = "23514"
	SQLStateExclusionViolation  SQLState = "23P01"
	SQLStateSerializationFail   SQLState = "40001"
)
```

Explanation:
Codes follow PostgreSQL Class 23 (Integrity Constraint Violation). Applications test `IsConstraintViolation(err, SQLStateUniqueViolation)` rather than checking error text. This is locale-safe and version-stable.

---

## Snippet 2 — ConstraintError Type with Structured Fields

Source File: `internal/dberr/errors.go:20-37`
Purpose: Error type carrying constraint name, table name, and detail for structured error reporting and mapping to domain errors.

```go
type ConstraintError struct {
	Code           SQLState
	ConstraintName string
	TableName      string
	Detail         string
	Err            error
}

func (e *ConstraintError) Error() string {
	if e.ConstraintName != "" {
		return fmt.Sprintf("db constraint error [%s] violated on table '%s' (SQLSTATE %s): %s", e.ConstraintName, e.TableName, e.Code, e.Detail)
	}
	return fmt.Sprintf("db error on table '%s' (SQLSTATE %s): %s", e.TableName, e.Code, e.Detail)
}
```

Explanation:
Separate fields allow domain mappers to extract constraint name for user-friendly messages like "Email already used". Detail field contains technical information for debugging.

---

## Snippet 3 — NOT NULL Violation Constructors

Source File: `internal/dberr/errors.go:39-46, 66-73`
Purpose: Factory functions creating specific constraint error instances for NOT NULL and FOREIGN KEY violations.

```go
func NewNotNullViolation(table, column, constraintName string) *ConstraintError {
	return &ConstraintError{
		Code:           SQLStateNotNullViolation,
		ConstraintName: constraintName,
		TableName:      table,
		Detail:         fmt.Sprintf("null value in column %q violates not-null constraint", column),
	}
}

func NewForeignKeyViolation(table, constraintName, detail string) *ConstraintError {
	return &ConstraintError{
		Code:           SQLStateForeignKeyViolation,
		ConstraintName: constraintName,
		TableName:      table,
		Detail:         detail,
	}
}
```

Explanation:
Each constraint type has its own constructor. Error messages follow PostgreSQL format. `column` and `constraintName` are parameterized for flexibility.

---

## Snippet 4 — MapToDomainError: Bridging Storage to User

Source File: `internal/dberr/errors.go:83-100`
Purpose: Maps constraint error codes to domain-friendly error messages with constraint names for troubleshooting.

```go
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

Explanation:
This function is the adapter between storage layer errors and application responses. It converts `23505` to "conflict" with constraint name for debugging. Each message type corresponds to appropriate HTTP response: `409 Conflict` for unique violation, `400 Bad Request` for not null, `422 Unprocessable Entity` for check violation.

---

## Snippet 5 — Engine Structure and Thread-Safety

Source File: `internal/engine/engine.go:12-30`
Purpose: Core engine struct providing thread-safe in-memory relational store with constraint indexes.

```go
type Engine struct {
	mu           sync.RWMutex
	userSeq      atomic.Int64
	orderSeq     atomic.Int64
	users        map[int64]model.User
	orders       map[int64]model.Order
	emailIndex   map[string]int64 // UNIQUE (email)
	activeEmails map[string]int64 // PARTIAL UNIQUE (email) WHERE deleted_at IS NULL
}

func NewEngine() *Engine {
	return &Engine{
		users:        make(map[int64]model.User),
		orders:       make(map[int64]model.Order),
		emailIndex:   make(map[string]int64),
		activeEmails: make(map[string]int64),
	}
}
```

Explanation:
Single `sync.RWMutex` provides coarse-grained table locking. `atomic.Int64` yields collision-free IDs without lock contention. Two email maps distinguish full UNIQUE (`emailIndex`) from partial UNIQUE with soft-delete (`activeEmails`).

---

## Snippet 6 — InsertUser: Atomic Constraint Evaluation

Source File: `internal/engine/engine.go:45-99`
Purpose: Inserts user atomically while evaluating NOT NULL, CHECK, UNIQUE, and PARTIAL UNIQUE constraints in lock-protected section.

```go
func (e *Engine) InsertUser(u model.User, usePartialUniqueIndex bool) (model.User, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// 1. NOT NULL checks
	if u.Email == "" {
		return model.User{}, dberr.NewNotNullViolation("users", "email", "users_email_not_null")
	}
	if u.Username == "" {
		return model.User{}, dberr.NewNotNullViolation("users", "username", "users_username_not_null")
	}

	// 2. CHECK constraints
	// CHECK (age >= 18)
	if u.Age < 18 {
		return model.User{}, dberr.NewCheckViolation("users", "users_age_check", fmt.Sprintf("value %d for column 'age' violates check constraint (age >= 18)", u.Age))
	}
	// CHECK (status IN ('active', 'suspended', 'pending'))
	switch u.Status {
	case "active", "suspended", "pending":
	default:
		return model.User{}, dberr.NewCheckViolation("users", "users_status_check", fmt.Sprintf("invalid status %q violates check constraint", u.Status))
	}

	// 3. UNIQUE / PARTIAL UNIQUE constraint
	if usePartialUniqueIndex {
		// Partial Unique Index: CREATE UNIQUE INDEX users_active_email_idx ON users (email) WHERE deleted_at IS NULL;
		if u.DeletedAt == nil {
			if existingID, exists := e.activeEmails[u.Email]; exists {
				return model.User{}, dberr.NewUniqueViolation("users", "users_active_email_idx", fmt.Sprintf("duplicate key value violates unique constraint 'users_active_email_idx' (email=%s, existing_id=%d)", u.Email, existingID))
			}
		}
	} else {
		// Standard full UNIQUE: CONSTRAINT users_email_key UNIQUE (email)
		if existingID, exists := e.emailIndex[u.Email]; exists {
			return model.User{}, dberr.NewUniqueViolation("users", "users_email_key", fmt.Sprintf("duplicate key value violates unique constraint 'users_email_key' (email=%s, existing_id=%d)", u.Email, existingID))
		}
	}

	// Primary Key generation
	if u.ID == 0 {
		u.ID = e.userSeq.Add(1)
	}

	// Commit row and indexes
	e.users[u.ID] = u
	if !usePartialUniqueIndex {
		e.emailIndex[u.Email] = u.ID
	} else if u.DeletedAt == nil {
		e.activeEmails[u.Email] = u.ID
	}

	return u, nil
}
```

Explanation:
Entire operation within single lock acquisition. Constraint checks happen before row commit. For partial index mode, check skipped if `DeletedAt != nil`. Index updates co-located with row insert for atomicity. This mirrors PostgreSQL's `ROW EXCLUSIVE` lock + B-tree index update pattern.

---

## Snippet 7 — InsertOrder: NOT NULL, CHECK, and FOREIGN KEY

Source File: `internal/engine/engine.go:122-148`
Purpose: Inserts order with NOT NULL, CHECK, and FOREIGN KEY constraint enforcement.

```go
func (e *Engine) InsertOrder(o model.Order) (model.Order, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// 1. NOT NULL constraint on user_id
	if o.UserID == 0 {
		return model.Order{}, dberr.NewNotNullViolation("orders", "user_id", "orders_user_id_not_null")
	}

	// 2. CHECK constraint (total_cents > 0)
	if o.TotalCents <= 0 {
		return model.Order{}, dberr.NewCheckViolation("orders", "orders_total_cents_check", fmt.Sprintf("total_cents %d must be greater than zero", o.TotalCents))
	}

	// 3. FOREIGN KEY constraint: CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES users(id)
	if _, exists := e.users[o.UserID]; !exists {
		return model.Order{}, dberr.NewForeignKeyViolation("orders", "fk_orders_user", fmt.Sprintf("key (user_id)=(%d) is not present in table \"users\"", o.UserID))
	}

	if o.ID == 0 {
		o.ID = e.orderSeq.Add(1)
	}

	e.orders[o.ID] = o
	return o, nil
}
```

Explanation:
Sequential evaluation: NOT NULL → CHECK → FK → commit. FOREIGN KEY check scans `e.users` map for existence. No separate index on `UserID` in orders table (lab simplicity) — production would index FK referencing columns.

---

## Snippet 8 — UnsafeStore: Demonstrating Race Condition

Source File: `internal/store/store.go:23-47`
Purpose: Shows vulnerable read-then-write pattern lacking database constraint protection under concurrent load.

```go
func (s *UnsafeStore) RegisterUser(ctx context.Context, u model.User) (model.User, error) {
	// Vulnerable app-level check: no database constraint protecting email uniqueness!
	// If two goroutines reach this concurrently, both see count == 0 and both proceed to insert.
	var duplicateFound bool
	count := s.eng.GetUsersCount()
	for i := 1; i <= count; i++ {
		existing, ok := s.eng.GetUser(int64(i))
		if ok && existing.Email == u.Email {
			duplicateFound = true
			break
		}
	}

	// Artificial yield to exacerbate read-then-write race window under CPU scheduling
	time.Sleep(1 * time.Millisecond)

	if duplicateFound {
		return model.User{}, fmt.Errorf("app check failed: email %s already taken", u.Email)
	}

	// Engine insert bypassing unique constraints (simulating unconstrained table)
	res, err := s.eng.InsertUserUnsafe(u)
	return res, err
}
```

Explanation:
Vulnerable pattern: read (scan) → `Sleep(1ms)` → write (`InsertUserUnsafe`). Sleep widens race window so concurrent goroutines see same count. `InsertUserUnsafe` bypasses index checks. Test proves duplicates occur.

---

## Snippet 9 — SafeStore: Delegating to Constraints

Source File: `internal/store/store.go:58-66`
Purpose: Safe store directly delegates to storage engine constraints, eliminating race window.

```go
func (s *SafeStore) RegisterUser(ctx context.Context, u model.User) (model.User, error) {
	// Directly insert; database constraint handles race condition and returns 23505
	res, err := s.eng.InsertUser(u, false)
	if err != nil {
		return model.User{}, dberr.MapToDomainError(err)
	}
	return res, nil
}
```

Explanation:
No read-then-write gap. `InsertUser` performs all checks atomically under lock. Second concurrent attempt fails with `23505` before it can see "succeeded" status of first. This is the atomicity guarantee of database constraints.

---

## Snippet 10 — Soft-Drop User (Partial Unique Index Update)

Source File: `internal/engine/engine.go:101-120`
Purpose: Removes user from partial unique index on soft delete, enabling email reuse.

```go
func (e *Engine) SoftDeleteUser(id int64, deletedAt model.User) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	u, exists := e.users[id]
	if !exists {
		return fmt.Errorf("user %d not found", id)
	}

	if deletedAt.DeletedAt == nil {
		return fmt.Errorf("deleted_at cannot be nil for soft delete")
	}

	// Remove from partial active index
	delete(e.activeEmails, u.Email)
	u.DeletedAt = deletedAt.DeletedAt
	e.users[id] = u
	return nil
}
```

Explanation:
Deletes email from `activeEmails` map (partial unique index). When user is re-registered with same email, constraint checks only `activeEmails`. This models `CREATE UNIQUE INDEX ... WHERE deleted_at IS NULL`. Note: duplicate soft-deleted rows allowed; they are not indexed by partial index.

---

## Snippet 11 — Concurrent Safety Test

Source File: `internal/store/store_test.go:162-208`
Purpose: Verifies exactly-once insertion under high concurrency with database constraints.

```go
func TestConcurrentRegistration_Safe_EnforcesUniqueness(t *testing.T) {
	eng := engine.NewEngine()
	s := store.NewSafeStore(eng)
	ctx := context.Background()

	goroutines := 20
	var wg sync.WaitGroup
	wg.Add(goroutines)

	successCount := 0
	errorCount := 0
	var mu sync.Mutex

	targetEmail := "concurrent@example.com"

	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			_, err := s.RegisterUser(ctx, model.User{
				Email:    targetEmail,
				Username: fmt.Sprintf("user_%d", idx),
				Age:      25,
				Status:   "active",
			})

			mu.Lock()
			defer mu.Unlock()
			if err == nil {
				successCount++
			} else {
				errorCount++
			}
		}(i)
	}

	wg.Wait()

	if successCount != 1 {
		t.Errorf("expected exactly 1 successful registration under concurrency, got %d", successCount)
	}
	if errorCount != goroutines-1 {
		t.Errorf("expected %d constraint errors, got %d", goroutines-1, errorCount)
	}
	if eng.GetUsersCount() != 1 {
		t.Errorf("expected total users count 1, got %d", eng.GetUsersCount())
	}
}
```

Explanation:
20 goroutines compete to register same email. Only one succeeds, 19 fail with constraint error. `GetUsersCount() == 1` proves no duplicates. Same test in demo uses 50 goroutines and prints results for visual verification.