package tests

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/lukman/software-engineering-lab/labs/19-database-connection-pooling/internal/pool"
)

func TestDirectConnectionOverhead(t *testing.T) {
	// Connect delay 5ms
	mockDriver := pool.NewMockDriver(100, 5*time.Millisecond)

	// DB without connection reuse
	unpooledDB := pool.OpenDB(mockDriver)
	defer unpooledDB.Close()
	unpooledDB.SetMaxIdleConns(0) // Forces a new connection every time

	start := time.Now()
	for i := 0; i < 5; i++ {
		_, err := unpooledDB.Exec("SELECT 1")
		if err != nil {
			t.Fatalf("unpooled failed: %v", err)
		}
	}
	unpooledDuration := time.Since(start)

	// DB with connection pooling
	pooledDB := pool.OpenDB(mockDriver)
	defer pooledDB.Close()
	pooledDB.SetMaxIdleConns(5)
	pooledDB.SetMaxOpenConns(5)

	start = time.Now()
	for i := 0; i < 5; i++ {
		_, err := pooledDB.Exec("SELECT 1")
		if err != nil {
			t.Fatalf("pooled failed: %v", err)
		}
	}
	pooledDuration := time.Since(start)

	if pooledDuration >= unpooledDuration {
		t.Errorf("expected pooled (%v) to be faster than unpooled (%v)", pooledDuration, unpooledDuration)
	}
}

func TestOversizedPoolExhaustsServerConnections(t *testing.T) {
	// Server allows max 10 connections
	mockDriver := pool.NewMockDriver(10, 0)

	// Client opens oversized pool of 20 connections
	db := pool.OpenDB(mockDriver)
	defer db.Close()
	db.SetMaxOpenConns(20)

	var wg sync.WaitGroup
	var errCount int32
	var mu sync.Mutex

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := db.Conn(context.Background())
			if err != nil {
				mu.Lock()
				errCount++
				mu.Unlock()
				return
			}
			// Hold connection briefly
			time.Sleep(20 * time.Millisecond)
			conn.Close()
		}()
	}

	wg.Wait()

	if errCount == 0 {
		t.Errorf("expected some connection attempts to fail due to server limits, but got 0 errors")
	}
}

func TestConnectionStarvationDueToLeak(t *testing.T) {
	mockDriver := pool.NewMockDriver(10, 0)
	db := pool.OpenDB(mockDriver)
	defer db.Close()

	// Strict pool size of 1
	db.SetMaxOpenConns(1)

	svc := pool.NewOrderService(db)

	// Start unsafe request that holds connection for 100ms
	go func() {
		_ = svc.ProcessOrderUnsafeLeak(context.Background(), 1, func() error {
			time.Sleep(100 * time.Millisecond)
			return nil
		})
	}()

	// Give goroutine time to acquire connection
	time.Sleep(10 * time.Millisecond)

	// Subsequent request with 20ms timeout should fail due to pool starvation
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := svc.ProcessOrderSafe(ctx, 2, nil)
	if err == nil {
		t.Fatalf("expected context deadline exceeded due to pool starvation, got nil")
	}
}

func TestSafeProcessingConcurrently(t *testing.T) {
	mockDriver := pool.NewMockDriver(5, 0)
	db := pool.OpenDB(mockDriver)
	defer db.Close()
	db.SetMaxOpenConns(5)

	svc := pool.NewOrderService(db)

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			err := svc.ProcessOrderSafe(ctx, id, func() error {
				time.Sleep(5 * time.Millisecond)
				return nil
			})
			if err != nil {
				t.Errorf("order %d failed: %v", id, err)
			}
		}(i)
	}

	wg.Wait()
}
