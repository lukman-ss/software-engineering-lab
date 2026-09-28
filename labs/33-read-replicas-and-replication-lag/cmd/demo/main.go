package main

import (
	"context"
	"fmt"
	"time"

	"labs/33-read-replicas-and-replication-lag/internal/cluster"
	"labs/33-read-replicas-and-replication-lag/internal/router"
)

func main() {
	fmt.Println("=== Lab 33: Read Replicas and Replication Lag Demo ===")

	c := cluster.NewCluster(2, 200*time.Millisecond, cluster.AsyncReplication)
	defer c.Close()

	cfg := router.DefaultConfig()
	cfg.StickyDuration = 500 * time.Millisecond
	cfg.WaitTimeout = 300 * time.Millisecond
	r := router.NewRouter(c, cfg)

	fmt.Println("\n--- 1. Demonstrating Async Replication Lag Anomaly ---")
	sessionID := "user-123"
	lsn, err := r.Write(sessionID, "user:profile:123", `{"name":"Alice","tier":"premium"}`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[Write Primary] Wrote key 'user:profile:123', Primary LSN: %d\n", lsn)

	val, nodeID, readLSN, err := r.ReadNaive("user:profile:123")
	fmt.Printf("[Naive Read] Node: %s, Applied LSN: %d, Val: '%s', Err: %v\n", nodeID, readLSN, val, err)

	fmt.Println("\n--- 2. Demonstrating Read-Your-Own-Writes with Sticky Session ---")
	stickyVal, stickyNode, stickyLSN, err := r.ReadWithStickySession(sessionID, "user:profile:123")
	fmt.Printf("[Sticky Read (within 500ms)] Node: %s, LSN: %d, Val: '%s', Err: %v\n", stickyNode, stickyLSN, stickyVal, err)

	fmt.Println("Sleeping 600ms for sticky duration to expire and replica to catch up...")
	time.Sleep(600 * time.Millisecond)

	expiredVal, expiredNode, expiredLSN, err := r.ReadWithStickySession(sessionID, "user:profile:123")
	fmt.Printf("[Sticky Read (after TTL)] Node: %s, LSN: %d, Val: '%s', Err: %v\n", expiredNode, expiredLSN, expiredVal, err)

	fmt.Println("\n--- 3. Demonstrating Causal Token / LSN Wait ---")
	lsn2, _ := r.Write(sessionID, "order:101", `{"status":"completed"}`)
	fmt.Printf("[Write Primary] Wrote key 'order:101', LSN: %d\n", lsn2)

	tokenVal, tokenNode, tokenLSN, err := r.ReadWithToken(context.Background(), lsn2, "order:101")
	fmt.Printf("[Token Read (MinLSN=%d)] Node: %s, LSN: %d, Val: '%s', Err: %v\n", lsn2, tokenNode, tokenLSN, tokenVal, err)

	fmt.Println("\n--- 4. Demonstrating Synchronous Replication (remote_apply) ---")
	cSync := cluster.NewCluster(2, 50*time.Millisecond, cluster.SyncReplication)
	defer cSync.Close()
	rSync := router.NewRouter(cSync, cfg)

	start := time.Now()
	syncLSN, _ := rSync.Write(sessionID, "user:profile:456", `{"name":"Bob"}`)
	duration := time.Since(start)
	fmt.Printf("[Sync Write] LSN: %d, Write Duration: %v\n", syncLSN, duration)

	syncVal, syncNode, syncReadLSN, err := rSync.ReadNaive("user:profile:456")
	fmt.Printf("[Sync Naive Read] Node: %s, LSN: %d, Val: '%s', Err: %v\n", syncNode, syncReadLSN, syncVal, err)

	fmt.Println("\n=== Demo Complete ===")
}
