package tests

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"labs/32-database-sharding-and-partitioning/internal/idgen"
	"labs/32-database-sharding-and-partitioning/internal/partitioning"
	"labs/32-database-sharding-and-partitioning/internal/sharding"
)

func TestPartitionPruning(t *testing.T) {
	ranges := []partitioning.PartitionRange{
		{Name: "p2026_01", Start: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)},
		{Name: "p2026_02", Start: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)},
		{Name: "p2026_03", Start: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)},
	}

	tbl := partitioning.NewTable("orders", ranges)

	// Insert into jan and feb
	rec1 := partitioning.Record{ID: "1", CreatedAt: time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC), Payload: "Jan order"}
	rec2 := partitioning.Record{ID: "2", CreatedAt: time.Date(2026, 2, 10, 10, 0, 0, 0, time.UTC), Payload: "Feb order"}
	rec3 := partitioning.Record{ID: "3", CreatedAt: time.Date(2026, 2, 20, 10, 0, 0, 0, time.UTC), Payload: "Feb order 2"}

	if err := tbl.Insert(rec1); err != nil {
		t.Fatalf("insert failed: %v", err)
	}
	if err := tbl.Insert(rec2); err != nil {
		t.Fatalf("insert failed: %v", err)
	}
	if err := tbl.Insert(rec3); err != nil {
		t.Fatalf("insert failed: %v", err)
	}

	// Query only Feb
	start := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC)

	res := tbl.QueryRange(start, end)
	if len(res.Records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(res.Records))
	}
	if res.PartitionsScanned != 1 {
		t.Fatalf("expected 1 partition scanned due to pruning, got %d", res.PartitionsScanned)
	}
	if res.TotalPartitions != 3 {
		t.Fatalf("expected 3 total partitions, got %d", res.TotalPartitions)
	}

	// Fast drop partition
	if !tbl.DropPartition("p2026_01") {
		t.Fatalf("failed to drop partition")
	}
	if tbl.PartitionCount() != 2 {
		t.Fatalf("expected 2 partitions after drop, got %d", tbl.PartitionCount())
	}
}

func TestRoutingAndConsistentHashRelocation(t *testing.T) {
	numKeys := 5000
	keys := make([]string, numKeys)
	for i := 0; i < numKeys; i++ {
		keys[i] = fmt.Sprintf("tenant_%d", i)
	}

	// 1. Modulo router test: 3 nodes -> 4 nodes
	initialShards := []string{"shard-1", "shard-2", "shard-3"}
	modRouter := sharding.NewModuloRouter(initialShards)
	modAssignments := make(map[string]string)
	for _, k := range keys {
		s, _ := modRouter.GetShard(k)
		modAssignments[k] = s
	}

	modRouter.AddShard("shard-4")
	modMoved := 0
	for _, k := range keys {
		s, _ := modRouter.GetShard(k)
		if s != modAssignments[k] {
			modMoved++
		}
	}
	modMoveRatio := float64(modMoved) / float64(numKeys)

	// 2. Consistent Hash router test: 3 nodes -> 4 nodes
	chRouter := sharding.NewConsistentHashRouter(150, initialShards)
	chAssignments := make(map[string]string)
	for _, k := range keys {
		s, _ := chRouter.GetShard(k)
		chAssignments[k] = s
	}

	chRouter.AddShard("shard-4")
	chMoved := 0
	for _, k := range keys {
		s, _ := chRouter.GetShard(k)
		if s != chAssignments[k] {
			chMoved++
		}
	}
	chMoveRatio := float64(chMoved) / float64(numKeys)

	t.Logf("Hash Modulo moved %d / %d keys (%.2f%%)", modMoved, numKeys, modMoveRatio*100)
	t.Logf("Consistent Hash moved %d / %d keys (%.2f%%)", chMoved, numKeys, chMoveRatio*100)

	// Modulo should move roughly 3/4 = 75% keys
	if modMoveRatio < 0.65 {
		t.Errorf("expected modulo move ratio >= 65%%, got %.2f%%", modMoveRatio*100)
	}

	// Consistent hashing should move roughly 1/4 = 25% keys (+/- variance) and strictly > 0%
	if chMoveRatio <= 0.05 || chMoveRatio > 0.40 {
		t.Errorf("expected consistent hash move ratio between 5%% and 40%%, got %.2f%%", chMoveRatio*100)
	}
}

func TestClusterScatterGatherAndGSI(t *testing.T) {
	shards := []string{"shard-0", "shard-1", "shard-2"}
	router := sharding.NewConsistentHashRouter(100, shards)
	cluster := sharding.NewCluster(router)

	records := []sharding.Record{
		{ID: "rec-1", ShardKey: "tenant-A", Email: "alice@example.com", Data: "Alice Data"},
		{ID: "rec-2", ShardKey: "tenant-B", Email: "bob@example.com", Data: "Bob Data"},
		{ID: "rec-3", ShardKey: "tenant-C", Email: "carol@example.com", Data: "Carol Data"},
	}

	for _, r := range records {
		if err := cluster.Insert(r); err != nil {
			t.Fatalf("failed to insert: %v", err)
		}
	}

	// Direct lookup with ShardKey
	rec, err := cluster.GetByShardKey("tenant-A", "rec-1")
	if err != nil || rec.Data != "Alice Data" {
		t.Fatalf("failed to get by shard key: %v", err)
	}

	// Lookup using Global Secondary Index (GSI)
	recGSI, err := cluster.GetByEmailUsingGSI("bob@example.com", "rec-2")
	if err != nil || recGSI.Data != "Bob Data" {
		t.Fatalf("failed to get by GSI: %v", err)
	}

	// Scatter Gather broadcast lookup
	sgRes := cluster.ScatterGatherBroadcast(func(r sharding.Record) bool {
		return r.Email == "carol@example.com"
	})
	if len(sgRes.Records) != 1 || sgRes.Records[0].ID != "rec-3" {
		t.Fatalf("scatter gather failed to find Carol")
	}
	if sgRes.ShardsQueried != 3 {
		t.Fatalf("expected 3 shards queried, got %d", sgRes.ShardsQueried)
	}

	// Scatter Gather with context timeout/cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // pre-canceled context
	canceledRes := cluster.ScatterGatherBroadcastWithContext(ctx, func(r sharding.Record) bool {
		return r.Email == "carol@example.com"
	})
	if canceledRes.ShardResponded != 0 {
		t.Fatalf("expected 0 shard responses on canceled context, got %d", canceledRes.ShardResponded)
	}
}

func TestIDGenerators(t *testing.T) {
	// UUIDv7 test
	u1, err := idgen.NewUUIDv7()
	if err != nil || len(u1) != 36 {
		t.Fatalf("uuidv7 generation failed: %v", err)
	}
	time.Sleep(2 * time.Millisecond)
	u2, err := idgen.NewUUIDv7()
	if err != nil {
		t.Fatalf("uuidv7 generation failed: %v", err)
	}
	if u1 >= u2 {
		t.Fatalf("expected u1 < u2 lexicographically for time-ordered UUIDv7")
	}

	// Central sequence allocator
	central := &idgen.MemoryCentralSequence{}
	allocator := idgen.NewSequenceBlockAllocator(10, central.AllocateBlock)

	ids := make([]int64, 25)
	for i := 0; i < 25; i++ {
		id, err := allocator.NextID()
		if err != nil {
			t.Fatalf("sequence error: %v", err)
		}
		ids[i] = id
	}

	for i := 0; i < 24; i++ {
		if ids[i+1] != ids[i]+1 {
			t.Fatalf("expected sequential IDs, got %d then %d", ids[i], ids[i+1])
		}
	}
}

func TestConcurrentClusterAccess(t *testing.T) {
	router := sharding.NewConsistentHashRouter(50, []string{"s1", "s2", "s3", "s4"})
	cluster := sharding.NewCluster(router)

	var wg sync.WaitGroup
	numOps := 200

	for i := 0; i < numOps; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sk := fmt.Sprintf("tenant-%d", idx%10)
			recID := fmt.Sprintf("id-%d", idx)
			email := fmt.Sprintf("user-%d@test.com", idx)

			_ = cluster.Insert(sharding.Record{
				ID:       recID,
				ShardKey: sk,
				Email:    email,
				Data:     "concurrent-test",
			})

			_, _ = cluster.GetByShardKey(sk, recID)
			_, _ = cluster.GetByEmailUsingGSI(email, recID)
		}(i)
	}

	wg.Wait()

	counts := cluster.GetShardCounts()
	total := 0
	for _, c := range counts {
		total += c
	}
	if total != numOps {
		t.Fatalf("expected %d total records, got %d", numOps, total)
	}
}
