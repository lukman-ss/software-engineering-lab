package tests

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"labs/33-read-replicas-and-replication-lag/internal/cluster"
	"labs/33-read-replicas-and-replication-lag/internal/router"
)

func TestNaiveReplicationLag_StaleRead(t *testing.T) {
	c := cluster.NewCluster(1, 500*time.Millisecond, cluster.AsyncReplication)
	defer c.Close()

	r := router.NewRouter(c, router.DefaultConfig())

	lsn, err := r.Write("session-1", "account:balance", "1000")
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if lsn != 1 {
		t.Fatalf("Expected LSN 1, got %d", lsn)
	}

	val, nodeID, _, err := r.ReadNaive("account:balance")
	if err == nil {
		t.Fatalf("Expected stale read / ErrNotFound on lagging replica, got val: %s from node %s", val, nodeID)
	}
	if err != cluster.ErrNotFound {
		t.Fatalf("Expected ErrNotFound, got %v", err)
	}

	time.Sleep(600 * time.Millisecond)

	val, nodeID, _, err = r.ReadNaive("account:balance")
	if err != nil {
		t.Fatalf("Expected successful read after lag delay, got error: %v", err)
	}
	if val != "1000" {
		t.Fatalf("Expected value '1000', got '%s'", val)
	}
}

func TestStickySessionRouting(t *testing.T) {
	c := cluster.NewCluster(1, 500*time.Millisecond, cluster.AsyncReplication)
	defer c.Close()

	cfg := router.DefaultConfig()
	cfg.StickyDuration = 300 * time.Millisecond
	r := router.NewRouter(c, cfg)

	sessionID := "user-session-abc"

	_, err := r.Write(sessionID, "cart:items", "item1,item2")
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	val, nodeID, _, err := r.ReadWithStickySession(sessionID, "cart:items")
	if err != nil {
		t.Fatalf("Expected successful sticky read, got %v", err)
	}
	if nodeID != "primary" {
		t.Fatalf("Expected read routed to 'primary', got '%s'", nodeID)
	}
	if val != "item1,item2" {
		t.Fatalf("Expected 'item1,item2', got '%s'", val)
	}

	otherVal, otherNode, _, otherErr := r.ReadWithStickySession("other-session", "cart:items")
	if otherErr == nil {
		t.Fatalf("Expected stale/not-found read for other-session on lagging replica, got val: %s from node %s", otherVal, otherNode)
	}

	time.Sleep(550 * time.Millisecond)

	valAfter, nodeIDAfter, _, errAfter := r.ReadWithStickySession(sessionID, "cart:items")
	if errAfter != nil {
		t.Fatalf("Expected successful read after sticky duration & replica catchup, got %v", errAfter)
	}
	if nodeIDAfter == "primary" {
		t.Fatalf("Expected read routed to replica after sticky TTL expiry, got '%s'", nodeIDAfter)
	}
	if valAfter != "item1,item2" {
		t.Fatalf("Expected 'item1,item2', got '%s'", valAfter)
	}
}

func TestReadWithToken_LSN(t *testing.T) {
	c := cluster.NewCluster(2, 200*time.Millisecond, cluster.AsyncReplication)
	defer c.Close()

	cfg := router.DefaultConfig()
	cfg.WaitTimeout = 300 * time.Millisecond
	r := router.NewRouter(c, cfg)

	lsn, err := r.Write("s1", "config:feature_flag", "enabled")
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	start := time.Now()
	val, nodeID, readLSN, err := r.ReadWithToken(context.Background(), lsn, "config:feature_flag")
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("ReadWithToken failed: %v", err)
	}
	if val != "enabled" {
		t.Fatalf("Expected 'enabled', got '%s'", val)
	}
	if readLSN < lsn {
		t.Fatalf("Expected read LSN >= %d, got %d", lsn, readLSN)
	}
	if elapsed < 150*time.Millisecond {
		t.Fatalf("Expected token read to wait for replica WAL catchup (~200ms), finished in %v (node: %s)", elapsed, nodeID)
	}
}

func TestReplicaLagThreshold_Fallback(t *testing.T) {
	c := cluster.NewCluster(1, 0, cluster.AsyncReplication)
	defer c.Close()

	cfg := router.DefaultConfig()
	cfg.MaxLSNDiff = 2
	r := router.NewRouter(c, cfg)

	replica := c.Replicas()[0]
	replica.SetLag(10 * time.Second)

	for i := 0; i < 5; i++ {
		_, _ = r.Write("s1", fmt.Sprintf("k%d", i), fmt.Sprintf("v%d", i))
	}

	val, nodeID, _, err := r.ReadLagAware("k0")
	if err != nil {
		t.Fatalf("ReadLagAware failed: %v", err)
	}
	if nodeID != "primary (fallback-lag)" {
		t.Fatalf("Expected fallback to primary due to excessive lag (>2 LSNs), got node '%s'", nodeID)
	}
	if val != "v0" {
		t.Fatalf("Expected 'v0', got '%s'", val)
	}
}

func TestSynchronousReplication_Freshness(t *testing.T) {
	c := cluster.NewCluster(2, 50*time.Millisecond, cluster.SyncReplication)
	defer c.Close()

	r := router.NewRouter(c, router.DefaultConfig())

	lsn, err := r.Write("s1", "user:balance", "500")
	if err != nil {
		t.Fatalf("Sync write failed: %v", err)
	}

	val, nodeID, readLSN, err := r.ReadNaive("user:balance")
	if err != nil {
		t.Fatalf("Sync naive read failed: %v", err)
	}
	if val != "500" {
		t.Fatalf("Expected '500', got '%s'", val)
	}
	if readLSN < lsn {
		t.Fatalf("Expected replica LSN >= %d, got %d (node: %s)", lsn, readLSN, nodeID)
	}
}

func TestConcurrentAccess_RaceFree(t *testing.T) {
	c := cluster.NewCluster(3, 10*time.Millisecond, cluster.AsyncReplication)
	defer c.Close()

	r := router.NewRouter(c, router.DefaultConfig())

	var wg sync.WaitGroup
	numWriters := 5
	numReaders := 10
	opsPerGoroutine := 20

	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(writerID int) {
			defer wg.Done()
			sessionID := fmt.Sprintf("writer-%d", writerID)
			for j := 0; j < opsPerGoroutine; j++ {
				key := fmt.Sprintf("key-%d-%d", writerID, j)
				val := fmt.Sprintf("val-%d-%d", writerID, j)
				token, err := r.Write(sessionID, key, val)
				if err != nil {
					t.Errorf("Concurrent write failed: %v", err)
					return
				}
				_, _, _, _ = r.ReadWithStickySession(sessionID, key)
				_, _, _, _ = r.ReadWithToken(context.Background(), token, key)
			}
		}(i)
	}

	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func(readerID int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				key := fmt.Sprintf("key-%d-%d", readerID%numWriters, j)
				_, _, _, _ = r.ReadNaive(key)
				_, _, _, _ = r.ReadLagAware(key)
			}
		}(i)
	}

	wg.Wait()
}

func TestWaitForLSN_ContextTimeout(t *testing.T) {
	c := cluster.NewCluster(1, 10*time.Second, cluster.AsyncReplication)
	defer c.Close()

	replica := c.Replicas()[0]

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := replica.WaitForLSN(ctx, 999)
	if err != context.DeadlineExceeded {
		t.Fatalf("Expected context.DeadlineExceeded, got %v", err)
	}
}
