package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"labs/30-leader-election/internal/candidate"
	"labs/30-leader-election/internal/coordinator"
	"labs/30-leader-election/internal/storage"
)

func TestCoordinatorLeaseAcquisitionAndRenewal(t *testing.T) {
	coord := coordinator.New()
	ttl := 200 * time.Millisecond

	lease1, err := coord.Acquire("leader-lock", "node-1", ttl)
	if err != nil {
		t.Fatalf("expected acquisition to succeed, got %v", err)
	}

	if lease1.FencingToken != 1 {
		t.Errorf("expected fencing token 1, got %d", lease1.FencingToken)
	}

	_, err = coord.Acquire("leader-lock", "node-2", ttl)
	if !errors.Is(err, coordinator.ErrLeaseHeld) {
		t.Errorf("expected ErrLeaseHeld, got %v", err)
	}

	renewed, err := coord.Renew("leader-lock", "node-1", lease1.FencingToken, ttl)
	if err != nil {
		t.Fatalf("expected renewal to succeed, got %v", err)
	}

	if renewed.FencingToken != 1 {
		t.Errorf("expected fencing token 1 after renewal, got %d", renewed.FencingToken)
	}

	time.Sleep(250 * time.Millisecond)

	_, err = coord.Renew("leader-lock", "node-1", lease1.FencingToken, ttl)
	if !errors.Is(err, coordinator.ErrLeaseExpired) {
		t.Errorf("expected ErrLeaseExpired after TTL lapse, got %v", err)
	}

	lease2, err := coord.Acquire("leader-lock", "node-2", ttl)
	if err != nil {
		t.Fatalf("expected node-2 to acquire after expiration, got %v", err)
	}

	if lease2.FencingToken != 2 {
		t.Errorf("expected fencing token 2 for second leader, got %d", lease2.FencingToken)
	}
}

func TestFencedStorageRejectsStaleTokens(t *testing.T) {
	store := storage.New()

	err := store.Write("node-1", 10, "state-1")
	if err != nil {
		t.Fatalf("expected first write to succeed, got %v", err)
	}

	err = store.Write("node-2", 11, "state-2")
	if err != nil {
		t.Fatalf("expected higher token write to succeed, got %v", err)
	}

	err = store.Write("node-1", 10, "stale-state")
	if !errors.Is(err, storage.ErrStaleFencingToken) {
		t.Fatalf("expected ErrStaleFencingToken for token 10 after token 11 was processed, got %v", err)
	}

	err = store.Write("node-2", 11, "duplicate-token")
	if !errors.Is(err, storage.ErrStaleFencingToken) {
		t.Fatalf("expected ErrStaleFencingToken for equal token 11, got %v", err)
	}

	history := store.History()
	if len(history) != 2 {
		t.Fatalf("expected exactly 2 records in store, got %d", len(history))
	}
}

func TestLeaderElectionFailoverAndSplitBrainDefense(t *testing.T) {
	coord := coordinator.New()
	store := storage.New()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ttl := 150 * time.Millisecond
	renewRate := 30 * time.Millisecond

	nodeA := candidate.NewNode("node-A", "cluster-leader", ttl, renewRate, coord, store)
	nodeB := candidate.NewNode("node-B", "cluster-leader", ttl, renewRate, coord, store)

	nodeA.Start(ctx)
	nodeB.Start(ctx)

	time.Sleep(80 * time.Millisecond)

	var leader1 *candidate.Node
	if nodeA.State() == candidate.StateLeader {
		leader1 = nodeA
	} else if nodeB.State() == candidate.StateLeader {
		leader1 = nodeB
	} else {
		t.Fatalf("expected one node to be leader")
	}

	err := leader1.PerformFencedWrite("leader-1-write")
	if err != nil {
		t.Fatalf("expected leader 1 write to succeed, got %v", err)
	}

	token1 := leader1.CurrentToken()

	// Simulate stop-the-world GC pause on Leader 1 longer than TTL
	leader1.SimulatePause(250 * time.Millisecond)

	time.Sleep(200 * time.Millisecond)

	var leader2 *candidate.Node
	if leader1 == nodeA {
		leader2 = nodeB
	} else {
		leader2 = nodeA
	}

	if leader2.State() != candidate.StateLeader {
		t.Fatalf("expected candidate 2 to take over leadership after pause")
	}

	err = leader2.PerformFencedWrite("leader-2-write")
	if err != nil {
		t.Fatalf("expected leader 2 write to succeed, got %v", err)
	}

	token2 := leader2.CurrentToken()
	if token2 <= token1 {
		t.Fatalf("expected leader 2 token (%d) > leader 1 token (%d)", token2, token1)
	}

	// Leader 1 wakes up and attempts write with old token
	err = store.Write(leader1.ID(), token1, "stale-write-after-gc-pause")
	if !errors.Is(err, storage.ErrStaleFencingToken) {
		t.Fatalf("expected stale write from paused leader to be rejected by fencing storage, got %v", err)
	}

	nodeA.Stop()
	nodeB.Stop()
}

func TestConcurrentElectionRace(t *testing.T) {
	coord := coordinator.New()
	store := storage.New()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nodes := make([]*candidate.Node, 5)
	for i := 0; i < 5; i++ {
		nodes[i] = candidate.NewNode(
			string(rune('A'+i)),
			"cluster-leader",
			100*time.Millisecond,
			20*time.Millisecond,
			coord,
			store,
		)
		nodes[i].Start(ctx)
	}

	time.Sleep(100 * time.Millisecond)

	leaderCount := 0
	for _, n := range nodes {
		if n.State() == candidate.StateLeader {
			leaderCount++
		}
	}

	if leaderCount != 1 {
		t.Fatalf("expected exactly 1 leader among concurrent nodes, found %d", leaderCount)
	}

	for _, n := range nodes {
		n.Stop()
	}
}
