package store_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/lukman/software-engineering-lab/labs/27-database-constraints/internal/dberr"
	"github.com/lukman/software-engineering-lab/labs/27-database-constraints/internal/engine"
	"github.com/lukman/software-engineering-lab/labs/27-database-constraints/internal/model"
	"github.com/lukman/software-engineering-lab/labs/27-database-constraints/internal/store"
)

func TestNotNullConstraints(t *testing.T) {
	eng := engine.NewEngine()
	s := store.NewSafeStore(eng)
	ctx := context.Background()

	// Missing Email
	_, err := s.RegisterUser(ctx, model.User{Username: "alice", Age: 25, Status: "active"})
	if err == nil {
		t.Fatal("expected NOT NULL error for missing email, got nil")
	}

	// Missing Username
	_, err = s.RegisterUser(ctx, model.User{Email: "alice@example.com", Age: 25, Status: "active"})
	if err == nil {
		t.Fatal("expected NOT NULL error for missing username, got nil")
	}

	// Missing UserID in order
	_, err = s.CreateOrder(ctx, model.Order{TotalCents: 100})
	if err == nil {
		t.Fatal("expected NOT NULL error for missing user_id in order, got nil")
	}
}

func TestCheckConstraints(t *testing.T) {
	eng := engine.NewEngine()
	s := store.NewSafeStore(eng)
	ctx := context.Background()

	// Underage (age < 18)
	_, err := s.RegisterUser(ctx, model.User{Email: "underage@example.com", Username: "bob", Age: 16, Status: "active"})
	if err == nil {
		t.Fatal("expected CHECK violation for age < 18, got nil")
	}

	// Invalid status
	_, err = s.RegisterUser(ctx, model.User{Email: "invalid@example.com", Username: "charlie", Age: 22, Status: "banned"})
	if err == nil {
		t.Fatal("expected CHECK violation for invalid status, got nil")
	}

	// Order total <= 0
	u, err := s.RegisterUser(ctx, model.User{Email: "dave@example.com", Username: "dave", Age: 30, Status: "active"})
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	_, err = s.CreateOrder(ctx, model.Order{UserID: u.ID, TotalCents: 0})
	if err == nil {
		t.Fatal("expected CHECK violation for total_cents <= 0, got nil")
	}
}

func TestUniqueConstraint(t *testing.T) {
	eng := engine.NewEngine()
	s := store.NewSafeStore(eng)
	ctx := context.Background()

	u1 := model.User{Email: "unique@example.com", Username: "u1", Age: 20, Status: "active"}
	_, err := s.RegisterUser(ctx, u1)
	if err != nil {
		t.Fatalf("first registration failed: %v", err)
	}

	// Duplicate insertion
	u2 := model.User{Email: "unique@example.com", Username: "u2", Age: 22, Status: "active"}
	_, err = s.RegisterUser(ctx, u2)
	if err == nil {
		t.Fatal("expected unique constraint violation on duplicate email, got nil")
	}
}

func TestForeignKeyConstraint(t *testing.T) {
	eng := engine.NewEngine()
	s := store.NewSafeStore(eng)
	ctx := context.Background()

	// Non-existent user reference
	_, err := s.CreateOrder(ctx, model.Order{UserID: 99999, TotalCents: 5000})
	if err == nil {
		t.Fatal("expected foreign key violation for non-existent user_id, got nil")
	}

	// Valid user reference
	u, err := s.RegisterUser(ctx, model.User{Email: "fk@example.com", Username: "fkuser", Age: 25, Status: "active"})
	if err != nil {
		t.Fatalf("user registration failed: %v", err)
	}

	o, err := s.CreateOrder(ctx, model.Order{UserID: u.ID, TotalCents: 5000})
	if err != nil {
		t.Fatalf("order creation failed: %v", err)
	}
	if o.ID == 0 {
		t.Fatal("expected valid order ID")
	}
}

func TestPartialUniqueIndex(t *testing.T) {
	eng := engine.NewEngine()
	s := store.NewSafeStore(eng)
	ctx := context.Background()

	email := "partial@example.com"

	// 1. Insert active user
	u1, err := s.RegisterUserPartial(ctx, model.User{Email: email, Username: "p1", Age: 25, Status: "active"})
	if err != nil {
		t.Fatalf("failed to insert active user: %v", err)
	}

	// 2. Attempt duplicate active user -> must fail
	_, err = s.RegisterUserPartial(ctx, model.User{Email: email, Username: "p2", Age: 26, Status: "active"})
	if err == nil {
		t.Fatal("expected partial index violation for active duplicate email, got nil")
	}

	// 3. Soft-delete active user
	now := time.Now()
	err = eng.SoftDeleteUser(u1.ID, model.User{DeletedAt: &now})
	if err != nil {
		t.Fatalf("failed to soft delete: %v", err)
	}

	// 4. Now inserting active user with same email should SUCCEED because u1 has deleted_at != NULL
	u3, err := s.RegisterUserPartial(ctx, model.User{Email: email, Username: "p3", Age: 27, Status: "active"})
	if err != nil {
		t.Fatalf("expected insert after soft-delete to succeed, got: %v", err)
	}

	// 5. Attempt second active user with same email -> must fail again
	_, err = s.RegisterUserPartial(ctx, model.User{Email: email, Username: "p4", Age: 28, Status: "active"})
	if err == nil {
		t.Fatal("expected partial index violation for second active user, got nil")
	}

	// 6. Insert soft-deleted user directly with same email -> allowed because WHERE deleted_at IS NULL doesn't index soft deleted
	delTime := time.Now()
	_, err = s.RegisterUserPartial(ctx, model.User{Email: email, Username: "p5", Age: 29, Status: "active", DeletedAt: &delTime})
	if err != nil {
		t.Fatalf("expected soft-deleted insert to bypass partial index, got: %v", err)
	}

	_ = u3
}

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

func TestConcurrentRegistration_Unsafe_SuffersRaceCondition(t *testing.T) {
	eng := engine.NewEngine()
	s := store.NewUnsafeStore(eng)
	ctx := context.Background()

	goroutines := 20
	var wg sync.WaitGroup
	wg.Add(goroutines)

	targetEmail := "unsafe_concurrent@example.com"

	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			_, _ = s.RegisterUser(ctx, model.User{
				Email:    targetEmail,
				Username: fmt.Sprintf("unsafe_user_%d", idx),
				Age:      25,
				Status:   "active",
			})
		}(i)
	}

	wg.Wait()

	// Without database constraints, read-then-write check permits duplicate inserts under concurrency
	if eng.GetUsersCount() <= 1 {
		t.Errorf("expected race condition causing >1 user entries, got %d", eng.GetUsersCount())
	}
}

func TestErrorClassification(t *testing.T) {
	nnErr := dberr.NewNotNullViolation("users", "email", "nn_idx")
	if !dberr.IsConstraintViolation(nnErr, dberr.SQLStateNotNullViolation) {
		t.Errorf("expected IsConstraintViolation true for SQLStateNotNullViolation")
	}

	uqErr := dberr.NewUniqueViolation("users", "uq_idx", "dup")
	if !dberr.IsConstraintViolation(uqErr, dberr.SQLStateUniqueViolation) {
		t.Errorf("expected IsConstraintViolation true for SQLStateUniqueViolation")
	}

	chkErr := dberr.NewCheckViolation("users", "chk_idx", "check fail")
	if !dberr.IsConstraintViolation(chkErr, dberr.SQLStateCheckViolation) {
		t.Errorf("expected IsConstraintViolation true for SQLStateCheckViolation")
	}

	fkErr := dberr.NewForeignKeyViolation("orders", "fk_idx", "fk fail")
	if !dberr.IsConstraintViolation(fkErr, dberr.SQLStateForeignKeyViolation) {
		t.Errorf("expected IsConstraintViolation true for SQLStateForeignKeyViolation")
	}
}
