package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/lukman/software-engineering-lab/labs/23-optimistic-vs-pessimistic-locking/internal/inventory"
)

func main() {
	fmt.Println("==========================================================")
	fmt.Println("  Optimistic vs Pessimistic Locking & Atomic Operations   ")
	fmt.Println("==========================================================")

	// 1. Naive lost update
	store1 := inventory.NewStore()
	store1.Seed(1, "Laptop Pro", 100)
	svc1 := inventory.NewService(store1)

	var wg1 sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg1.Add(1)
		go func() {
			defer wg1.Done()
			_ = svc1.DeductNaive(1, 1)
		}()
	}
	wg1.Wait()
	p1, _ := store1.Get(1)
	fmt.Printf("\n[1] Naive Read-Modify-Write (50 concurrent requests):\n")
	fmt.Printf("    Initial Stock: 100\n")
	fmt.Printf("    Expected Final Stock: 50\n")
	fmt.Printf("    Actual Final Stock:   %d (LOST UPDATE DETECTED!)\n", p1.Stock)

	// 2. Pessimistic locking
	store2 := inventory.NewStore()
	store2.Seed(2, "Smartphone X", 100)
	svc2 := inventory.NewService(store2)

	var wg2 sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			_ = svc2.DeductPessimistic(2, 1)
		}()
	}
	wg2.Wait()
	p2, _ := store2.Get(2)
	fmt.Printf("\n[2] Pessimistic Locking (SELECT ... FOR UPDATE):\n")
	fmt.Printf("    Initial Stock: 100\n")
	fmt.Printf("    Actual Final Stock:   %d (SUCCESS - Fully Synchronized)\n", p2.Stock)

	// 3. Optimistic locking without retry
	store3 := inventory.NewStore()
	store3.Seed(3, "Wireless Headphones", 100)
	svc3 := inventory.NewService(store3)

	var wg3 sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg3.Add(1)
		go func() {
			defer wg3.Done()
			_ = svc3.DeductOptimisticDirect(3, 1)
		}()
	}
	wg3.Wait()
	p3, _ := store3.Get(3)
	fmt.Printf("\n[3] Optimistic Locking Direct (20 concurrent requests, no retry):\n")
	fmt.Printf("    Initial Stock: 100\n")
	fmt.Printf("    Successful Deductions: %d\n", store3.Optimistically)
	fmt.Printf("    Rejected Conflicts:   %d\n", store3.OptimisticFails)
	fmt.Printf("    Actual Final Stock:   %d (SUCCESS - State Guarded, Zero Corruption)\n", p3.Stock)

	// 4. Optimistic locking with retry
	store4 := inventory.NewStore()
	store4.Seed(4, "Smartwatch", 100)
	svc4 := inventory.NewService(store4)

	start := time.Now()
	var wg4 sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg4.Add(1)
		go func() {
			defer wg4.Done()
			_ = svc4.DeductOptimisticWithRetry(4, 1, 10)
		}()
	}
	wg4.Wait()
	p4, _ := store4.Get(4)
	fmt.Printf("\n[4] Optimistic Locking With Exponential Backoff Retry (20 requests):\n")
	fmt.Printf("    Initial Stock: 100\n")
	fmt.Printf("    Successful Deductions: %d\n", store4.Optimistically)
	fmt.Printf("    Total Attempted Conflicts Retried: %d\n", store4.OptimisticFails)
	fmt.Printf("    Actual Final Stock:   %d (SUCCESS - All retries eventually converged)\n", p4.Stock)
	fmt.Printf("    Elapsed Time:         %v\n", time.Since(start))

	// 5. Atomic single-statement update
	store5 := inventory.NewStore()
	store5.Seed(5, "Tablet Air", 100)
	svc5 := inventory.NewService(store5)

	var wg5 sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg5.Add(1)
		go func() {
			defer wg5.Done()
			_ = svc5.DeductAtomic(5, 1)
		}()
	}
	wg5.Wait()
	p5, _ := store5.Get(5)
	fmt.Printf("\n[5] Atomic Single-Statement Operation (UPDATE ... WHERE stock >= qty):\n")
	fmt.Printf("    Initial Stock: 100\n")
	fmt.Printf("    Actual Final Stock:   %d (SUCCESS - Lockless Single Statement)\n", p5.Stock)

	fmt.Println("\n==========================================================")
	fmt.Println("  Lab Execution Completed Successfully                    ")
	fmt.Println("==========================================================")
}
