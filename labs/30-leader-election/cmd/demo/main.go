package main

import (
	"context"
	"fmt"
	"time"

	"labs/30-leader-election/internal/candidate"
	"labs/30-leader-election/internal/coordinator"
	"labs/30-leader-election/internal/storage"
)

func main() {
	fmt.Println("=== Distributed Leader Election with Fencing Tokens Demo ===")

	coord := coordinator.New()
	store := storage.New()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ttl := 200 * time.Millisecond
	renewInterval := 50 * time.Millisecond

	nodeA := candidate.NewNode("Node-A", "app-primary", ttl, renewInterval, coord, store)
	nodeB := candidate.NewNode("Node-B", "app-primary", ttl, renewInterval, coord, store)

	fmt.Println("\n[1] Starting Node-A and Node-B election campaigns...")
	nodeA.Start(ctx)
	nodeB.Start(ctx)

	time.Sleep(100 * time.Millisecond)

	fmt.Printf("Status: Node-A=%s, Node-B=%s\n", nodeA.State(), nodeB.State())

	var leader, standby *candidate.Node
	if nodeA.State() == candidate.StateLeader {
		leader, standby = nodeA, nodeB
	} else {
		leader, standby = nodeB, nodeA
	}

	token1 := leader.CurrentToken()
	fmt.Printf("[2] Leader elected: %s (Fencing Token: %d)\n", leader.ID(), token1)

	fmt.Printf("[3] %s performing legitimate fenced write...\n", leader.ID())
	err := leader.PerformFencedWrite("record-from-initial-leader")
	if err != nil {
		fmt.Printf("Write error: %v\n", err)
	} else {
		fmt.Printf("Write SUCCESS by %s (Token %d)\n", leader.ID(), token1)
	}

	fmt.Printf("\n[4] Simulating Stop-the-World GC Pause / Network Partition on %s (350ms > TTL %s)...\n", leader.ID(), ttl)
	leader.SimulatePause(350 * time.Millisecond)

	// Wait for TTL expiration + standby promotion
	time.Sleep(250 * time.Millisecond)

	fmt.Printf("[5] Status after TTL expiration: %s=%s, %s=%s\n", leader.ID(), leader.State(), standby.ID(), standby.State())
	token2 := standby.CurrentToken()
	fmt.Printf("[6] New Leader promoted: %s (Fencing Token: %d)\n", standby.ID(), token2)

	fmt.Printf("[7] New Leader %s performing fenced write...\n", standby.ID())
	err = standby.PerformFencedWrite("record-from-failover-leader")
	if err != nil {
		fmt.Printf("Write error: %v\n", err)
	} else {
		fmt.Printf("Write SUCCESS by %s (Token %d)\n", standby.ID(), token2)
	}

	// Wait for old leader to wake up
	time.Sleep(150 * time.Millisecond)

	fmt.Printf("\n[8] Old Leader %s wakes up from GC pause (still holding stale token %d)...\n", leader.ID(), token1)
	fmt.Printf("[9] Old Leader %s attempts write to shared storage using stale token %d...\n", leader.ID(), token1)
	err = store.Write(leader.ID(), token1, "split-brain-stale-write")
	if err != nil {
		fmt.Printf("Write REJECTED by FencedStorage: %v\n", err)
		fmt.Println("Result: Split-brain prevented! Old leader's stale write was safely rejected.")
	} else {
		fmt.Println("CRITICAL ERROR: Stale write succeeded! Fencing check failed.")
	}

	fmt.Println("\n[10] Shared Storage History:")
	for idx, rec := range store.History() {
		fmt.Printf("  %d. Token: %d | Author: %s | Value: %s\n", idx+1, rec.FencingToken, rec.Author, rec.Value)
	}

	nodeA.Stop()
	nodeB.Stop()
	fmt.Println("\n=== Demo completed successfully ===")
}
