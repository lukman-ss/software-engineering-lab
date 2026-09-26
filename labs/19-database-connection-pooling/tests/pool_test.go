package tests

import (
	"context"
	"errors"
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

	if errCount < 10 {
		t.Errorf("expected at least 10 connection attempts to fail (server max=10, pool size=20), got %d errors", errCount)
	}
}

func TestConnectionStarvationDueToLeak(t *testing.T) {
	mockDriver := pool.NewMockDriver(10, 0)
	db := pool.OpenDB(mockDriver)
	defer db.Close()

	// Strict pool size of 1
	db.SetMaxOpenConns(1)

	svc := pool.NewOrderService(db)

	acquired := make(chan struct{})

	// Start unsafe request that holds connection for 100ms
	go func() {
		_ = svc.ProcessOrderUnsafeLeak(context.Background(), 1, func() error {
			close(acquired)
			time.Sleep(100 * time.Millisecond)
			return nil
		})
	}()

	// Wait until goroutine has acquired the connection and is in externalCall
	select {
	case <-acquired:
	case <-time.After(2 * time.Second):
		t.Fatal("leak goroutine did not acquire connection within 2s")
	}

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

func TestPoolLockingDeadlock(t *testing.T) {
	// Pool of size 1: acquiring a second connection while holding the first results in deadlock/timeout
	mockDriver := pool.NewMockDriver(10, 0)
	db := pool.OpenDB(mockDriver)
	defer db.Close()
	db.SetMaxOpenConns(1)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	conn1, err := db.Conn(ctx)
	if err != nil {
		t.Fatalf("failed to acquire first connection: %v", err)
	}
	defer conn1.Close()

	_, err = db.Conn(ctx)
	if err == nil {
		t.Fatalf("expected context deadline exceeded when acquiring 2nd connection on pool of size 1, got nil")
	}
}

func TestMockConnDoubleClose(t *testing.T) {
	mockDriver := pool.NewMockDriver(10, 0)
	conn, err := mockDriver.Open("")
	if err != nil {
		t.Fatalf("failed to open mock connection: %v", err)
	}

	err = conn.Close()
	if err != nil {
		t.Fatalf("first close failed: %v", err)
	}

	// Should not panic, return error, or double-decrement active conns
	err = conn.Close()
	if err != nil {
		t.Fatalf("second close failed: %v", err)
	}

	if mockDriver.ActiveConnections() != 0 {
		t.Errorf("expected 0 active connections, got %d", mockDriver.ActiveConnections())
	}
}

func TestExternalCallErrorPropagation(t *testing.T) {
	callErr := errors.New("external service unavailable")

	t.Run("ProcessOrderSafe", func(t *testing.T) {
		mockDriver := pool.NewMockDriver(10, 0)
		db := pool.OpenDB(mockDriver)
		svc := pool.NewOrderService(db)

		err := svc.ProcessOrderSafe(context.Background(), 1, func() error {
			return callErr
		})
		if err != callErr {
			t.Errorf("expected %v, got %v", callErr, err)
		}
		db.Close()
		if mockDriver.ActiveConnections() != 0 {
			t.Errorf("expected 0 active connections after safe error, got %d", mockDriver.ActiveConnections())
		}
	})

	t.Run("ProcessOrderUnsafeLeak", func(t *testing.T) {
		mockDriver := pool.NewMockDriver(10, 0)
		db := pool.OpenDB(mockDriver)
		svc := pool.NewOrderService(db)

		err := svc.ProcessOrderUnsafeLeak(context.Background(), 2, func() error {
			return callErr
		})
		if err != callErr {
			t.Errorf("expected %v, got %v", callErr, err)
		}
		db.Close()
		if mockDriver.ActiveConnections() != 0 {
			t.Errorf("expected 0 active connections after unsafe error (defer close), got %d", mockDriver.ActiveConnections())
		}
	})
}

func TestPreCancelledContextProcessOrderSafe(t *testing.T) {
	mockDriver := pool.NewMockDriver(10, 0)
	db := pool.OpenDB(mockDriver)
	defer db.Close()

	svc := pool.NewOrderService(db)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := svc.ProcessOrderSafe(ctx, 1, nil)
	if err == nil {
		t.Fatal("expected error with pre-cancelled context, got nil")
	}
}
