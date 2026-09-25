package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"labs/18-deadlock/internal/bank"
	"labs/18-deadlock/internal/transfer"
)

func main() {
	fmt.Println("=== Deadlock Simulation Demo ===")

	// 1. Naive Transfer (Deadlock Scenario)
	fmt.Println("\n1. Simulating Naive Concurrent Transfers (Circular Wait)...")
	acc1 := bank.NewAccount("ACC-1", 1000)
	acc2 := bank.NewAccount("ACC-2", 1000)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		err := transfer.TransferNaive(ctx, acc1, acc2, 100, 20*time.Millisecond)
		fmt.Printf("Transfer ACC-1 -> ACC-2: %v\n", err)
	}()

	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		err := transfer.TransferNaive(ctx, acc2, acc1, 200, 20*time.Millisecond)
		fmt.Printf("Transfer ACC-2 -> ACC-1: %v\n", err)
	}()

	wg.Wait()

	// 2. Lock Ordering (Deadlock Prevention)
	fmt.Println("\n2. Simulating Lock Ordered Transfers (Deadlock Prevention)...")
	acc1 = bank.NewAccount("ACC-1", 1000)
	acc2 = bank.NewAccount("ACC-2", 1000)

	wg.Add(2)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		err := transfer.TransferOrdered(ctx, acc1, acc2, 100, 20*time.Millisecond)
		fmt.Printf("TransferOrdered ACC-1 -> ACC-2: %v\n", err)
	}()

	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		err := transfer.TransferOrdered(ctx, acc2, acc1, 200, 20*time.Millisecond)
		fmt.Printf("TransferOrdered ACC-2 -> ACC-1: %v\n", err)
	}()

	wg.Wait()
	fmt.Printf("Balances: ACC-1=%d, ACC-2=%d\n", acc1.Balance, acc2.Balance)

	// 3. Retry Strategy (Deadlock Recovery)
	fmt.Println("\n3. Simulating Retry Mechanism (Deadlock Recovery)...")
	acc1 = bank.NewAccount("ACC-1", 1000)
	acc2 = bank.NewAccount("ACC-2", 1000)

	wg.Add(2)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		err := transfer.TransferWithRetry(ctx, acc1, acc2, 100, 10*time.Millisecond, 5)
		fmt.Printf("TransferWithRetry ACC-1 -> ACC-2: %v\n", err)
	}()

	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		err := transfer.TransferWithRetry(ctx, acc2, acc1, 200, 10*time.Millisecond, 5)
		fmt.Printf("TransferWithRetry ACC-2 -> ACC-1: %v\n", err)
	}()

	wg.Wait()
	fmt.Printf("Final Balances: ACC-1=%d, ACC-2=%d\n", acc1.Balance, acc2.Balance)
}
