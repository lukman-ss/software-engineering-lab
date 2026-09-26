package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lukman/software-engineering-lab/labs/19-database-connection-pooling/internal/pool"
)

func main() {
	fmt.Println("--- Database Connection Pooling Demo ---")
	
	fmt.Println("\n1. Direct Connection Overhead Penalty")
	demoDirectOverhead()

	fmt.Println("\n2. Oversized Pool Exhausting Server Connections")
	demoOversizedPool()

	fmt.Println("\n3. Connection Leak Starving Pool")
	demoConnectionLeak()
}

func demoDirectOverhead() {
	driver := pool.NewMockDriver(100, 10*time.Millisecond) // 10ms handshake penalty

	unpooledDB := pool.OpenDB(driver)
	unpooledDB.SetMaxIdleConns(0)

	start := time.Now()
	for i := 0; i < 5; i++ {
		unpooledDB.Exec("SELECT 1")
	}
	fmt.Printf("Unpooled (5 requests): %v\n", time.Since(start))

	pooledDB := pool.OpenDB(driver)
	pooledDB.SetMaxIdleConns(5)
	pooledDB.SetMaxOpenConns(5)
	
	// Warm up pool
	pooledDB.Exec("SELECT 1")

	start = time.Now()
	for i := 0; i < 5; i++ {
		pooledDB.Exec("SELECT 1")
	}
	fmt.Printf("Pooled (5 requests): %v\n", time.Since(start))
}

func demoOversizedPool() {
	driver := pool.NewMockDriver(15, 0) // Database server max_connections = 15
	
	db := pool.OpenDB(driver)
	db.SetMaxOpenConns(50) // Client oversized pool = 50

	var wg sync.WaitGroup
	var successCount int32
	var failCount int32

	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := db.Conn(context.Background())
			if err != nil {
				atomic.AddInt32(&failCount, 1)
				return
			}
			time.Sleep(10 * time.Millisecond)
			conn.Close()
			atomic.AddInt32(&successCount, 1)
		}()
	}
	wg.Wait()

	fmt.Printf("Client attempted: 30, Succeeded: %d, Server Rejected: %d\n", successCount, failCount)
}

func demoConnectionLeak() {
	driver := pool.NewMockDriver(10, 0)
	db := pool.OpenDB(driver)
	db.SetMaxOpenConns(2) // strict small pool
	
	svc := pool.NewOrderService(db)

	fmt.Println("Starting 2 unsafe orders (holding connection during slow external IO)")
	for i := 0; i < 2; i++ {
		go func(id int) {
			svc.ProcessOrderUnsafeLeak(context.Background(), id, func() error {
				time.Sleep(500 * time.Millisecond)
				return nil
			})
		}(i)
	}

	time.Sleep(50 * time.Millisecond) // Let goroutines claim the pool

	fmt.Println("Attempting 3rd order with safe flow and short timeout...")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := svc.ProcessOrderSafe(ctx, 3, nil)
	if err != nil {
		fmt.Printf("Order 3 Failed: %v\n", err)
	} else {
		fmt.Println("Order 3 Succeeded")
	}
}
