package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/lukman/software-engineering-lab/labs/27-database-constraints/internal/dberr"
	"github.com/lukman/software-engineering-lab/labs/27-database-constraints/internal/engine"
	"github.com/lukman/software-engineering-lab/labs/27-database-constraints/internal/model"
	"github.com/lukman/software-engineering-lab/labs/27-database-constraints/internal/store"
)

func main() {
	fmt.Println("=================================================================")
	fmt.Println("LAB 27: DATABASE CONSTRAINTS & DATA INTEGRITY DEMONSTRATION")
	fmt.Println("=================================================================")

	ctx := context.Background()
	eng := engine.NewEngine()
	safeStore := store.NewSafeStore(eng)

	// 1. NOT NULL Constraint Violation Demo
	fmt.Println("\n[1] DEMONSTRATING NOT NULL CONSTRAINT (SQLSTATE 23502)")
	_, err := safeStore.RegisterUser(ctx, model.User{Username: "missing_email", Age: 25, Status: "active"})
	fmt.Printf("Attempt insert with missing email -> Error: %v\n", err)

	// 2. CHECK Constraint Violation Demo
	fmt.Println("\n[2] DEMONSTRATING CHECK CONSTRAINT (SQLSTATE 23514)")
	_, err = safeStore.RegisterUser(ctx, model.User{Email: "minor@test.com", Username: "minor", Age: 15, Status: "active"})
	fmt.Printf("Attempt insert with age=15 (CHECK age >= 18) -> Error: %v\n", err)

	_, err = safeStore.RegisterUser(ctx, model.User{Email: "badstatus@test.com", Username: "badstatus", Age: 20, Status: "invalid_status"})
	fmt.Printf("Attempt insert with invalid status -> Error: %v\n", err)

	// 3. FOREIGN KEY Constraint Violation Demo
	fmt.Println("\n[3] DEMONSTRATING FOREIGN KEY CONSTRAINT (SQLSTATE 23503)")
	_, err = safeStore.CreateOrder(ctx, model.Order{UserID: 9999, TotalCents: 2500})
	fmt.Printf("Attempt insert order for non-existent UserID=9999 -> Error: %v\n", err)

	// 4. PARTIAL UNIQUE INDEX (Soft Delete Pattern)
	fmt.Println("\n[4] DEMONSTRATING PARTIAL UNIQUE INDEX (WHERE deleted_at IS NULL)")
	email := "alice@company.com"
	u1, err := safeStore.RegisterUserPartial(ctx, model.User{Email: email, Username: "alice_active1", Age: 30, Status: "active"})
	fmt.Printf("Created active user ID=%d (%s)\n", u1.ID, email)

	_, err = safeStore.RegisterUserPartial(ctx, model.User{Email: email, Username: "alice_active2", Age: 30, Status: "active"})
	fmt.Printf("Duplicate active user insert rejected -> Error: %v\n", err)

	now := time.Now()
	_ = eng.SoftDeleteUser(u1.ID, model.User{DeletedAt: &now})
	fmt.Printf("Soft-deleted user ID=%d (deleted_at set)\n", u1.ID)

	u2, err := safeStore.RegisterUserPartial(ctx, model.User{Email: email, Username: "alice_active3", Age: 30, Status: "active"})
	fmt.Printf("New active user re-using email after soft delete -> Created user ID=%d (%s)\n", u2.ID, email)

	// 5. CONCURRENT RACE CONDITION vs DATABASE UNIQUE CONSTRAINT (SQLSTATE 23505)
	fmt.Println("\n[5] CONCURRENCY STRESS TEST: 50 CONCURRENT REGISTRATIONS FOR SAME EMAIL")
	concurrentEmail := "race-target@domain.com"
	workers := 50
	var wg sync.WaitGroup
	wg.Add(workers)

	var successCount, errorCount int
	var mu sync.Mutex

	for i := 0; i < workers; i++ {
		go func(id int) {
			defer wg.Done()
			_, regErr := safeStore.RegisterUser(ctx, model.User{
				Email:    concurrentEmail,
				Username: fmt.Sprintf("user_%d", id),
				Age:      25,
				Status:   "active",
			})

			mu.Lock()
			defer mu.Unlock()
			if regErr == nil {
				successCount++
			} else {
				errorCount++
			}
		}(i)
	}

	wg.Wait()

	fmt.Printf("Results:\n")
	fmt.Printf("  - Total Goroutines: %d\n", workers)
	fmt.Printf("  - Successful Registrations: %d\n", successCount)
	fmt.Printf("  - Rejected with UNIQUE VIOLATION (23505): %d\n", errorCount)
	fmt.Printf("  - Database Integrity Intact: %v\n", successCount == 1 && errorCount == workers-1)

	// Verify SQLSTATE taxonomy
	testErr := dberr.NewUniqueViolation("users", "users_email_key", "duplicate key")
	fmt.Printf("\nSQLSTATE Taxonomy Verification: code=%s isUniqueViolation=%v\n", testErr.Code, dberr.IsConstraintViolation(testErr, dberr.SQLStateUniqueViolation))
}
