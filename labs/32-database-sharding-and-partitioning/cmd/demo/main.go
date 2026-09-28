package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	"labs/32-database-sharding-and-partitioning/internal/idgen"
	"labs/32-database-sharding-and-partitioning/internal/partitioning"
	"labs/32-database-sharding-and-partitioning/internal/sharding"
)

func main() {
	fmt.Println("================================================================================")
	fmt.Println("LAB 32: DATABASE SHARDING AND PARTITIONING DEMONSTRATION")
	fmt.Println("================================================================================")

	demoPartitioning()
	demoShardingHotspots()
	demoRoutingScaleOutComparison()
	demoScatterGatherVsGSI()
	demoIDGeneration()

	fmt.Println("\n[DEMO COMPLETE] All database sharding & partitioning concepts successfully executed.")
}

func demoPartitioning() {
	fmt.Println("\n--- 1. Single-Node Logical Table Partitioning & Range Pruning ---")
	ranges := []partitioning.PartitionRange{
		{Name: "p2026_q1", Start: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)},
		{Name: "p2026_q2", Start: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)},
		{Name: "p2026_q3", Start: time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)},
		{Name: "p2026_q4", Start: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)},
	}

	tbl := partitioning.NewTable("sales_logs", ranges)

	_ = tbl.Insert(partitioning.Record{ID: "r1", CreatedAt: time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC), Payload: "Q1 sale"})
	_ = tbl.Insert(partitioning.Record{ID: "r2", CreatedAt: time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC), Payload: "Q2 sale"})
	_ = tbl.Insert(partitioning.Record{ID: "r3", CreatedAt: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC), Payload: "Q3 sale"})

	start := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)
	res := tbl.QueryRange(start, end)

	fmt.Printf("Query Range: %s to %s\n", start.Format("2006-01-02"), end.Format("2006-01-02"))
	fmt.Printf("Records Found: %d | Partitions Scanned: %d / %d (Pruned %d partitions)\n",
		len(res.Records), res.PartitionsScanned, res.TotalPartitions, res.TotalPartitions-res.PartitionsScanned)
}

func demoShardingHotspots() {
	fmt.Println("\n--- 2. Sharding Key Selection: Monotonic Key vs High-Cardinality Key ---")
	shardIDs := []string{"shard-0", "shard-1", "shard-2", "shard-3"}

	// Monotonic key routing (e.g., date-based shard key)
	fmt.Println("\n[Scenario A] Monotonic Key (date-string) Sharding:")
	dateRouter := sharding.NewModuloRouter(shardIDs)
	dateCluster := sharding.NewCluster(dateRouter)

	// Ingest sequential timestamp keys
	for i := 0; i < 1000; i++ {
		// All writes on same day get same key -> write hotspot!
		key := "2026-09-28"
		_ = dateCluster.Insert(sharding.Record{ID: fmt.Sprintf("id-%d", i), ShardKey: key, Data: "item"})
	}

	for id, count := range dateCluster.GetShardCounts() {
		fmt.Printf("  Node %s: %d records %s\n", id, count, barChart(count, 1000))
	}
	fmt.Println("Result: Severe Write Hotspot! 100% writes hit single shard.")

	// High cardinality key routing (tenant_id / user_id)
	fmt.Println("\n[Scenario B] High-Cardinality Key (user_id) Sharding:")
	userRouter := sharding.NewConsistentHashRouter(100, shardIDs)
	userCluster := sharding.NewCluster(userRouter)

	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("user_tenant_%d", i)
		_ = userCluster.Insert(sharding.Record{ID: fmt.Sprintf("id-%d", i), ShardKey: key, Data: "item"})
	}

	for id, count := range userCluster.GetShardCounts() {
		fmt.Printf("  Node %s: %d records %s\n", id, count, barChart(count, 1000))
	}
	fmt.Println("Result: Uniform distribution across physical shards.")
}

func demoRoutingScaleOutComparison() {
	fmt.Println("\n--- 3. Resharding / Scale-out Comparison: Hash Modulo vs Consistent Hashing ---")
	numKeys := 10000
	keys := make([]string, numKeys)
	for i := 0; i < numKeys; i++ {
		keys[i] = fmt.Sprintf("account_%d", i)
	}

	initialShards := []string{"shard-0", "shard-1", "shard-2", "shard-3"}

	// Hash Modulo
	modRouter := sharding.NewModuloRouter(initialShards)
	modBefore := make(map[string]string)
	for _, k := range keys {
		s, _ := modRouter.GetShard(k)
		modBefore[k] = s
	}

	modRouter.AddShard("shard-4") // Scale from 4 -> 5 nodes
	modMoved := 0
	for _, k := range keys {
		s, _ := modRouter.GetShard(k)
		if s != modBefore[k] {
			modMoved++
		}
	}

	// Consistent Hashing (100 vnodes)
	chRouter := sharding.NewConsistentHashRouter(100, initialShards)
	chBefore := make(map[string]string)
	for _, k := range keys {
		s, _ := chRouter.GetShard(k)
		chBefore[k] = s
	}

	chRouter.AddShard("shard-4") // Scale from 4 -> 5 nodes
	chMoved := 0
	for _, k := range keys {
		s, _ := chRouter.GetShard(k)
		if s != chBefore[k] {
			chMoved++
		}
	}

	fmt.Printf("Cluster Resize: 4 Shards -> 5 Shards (Total Keys: %d)\n", numKeys)
	fmt.Printf("  Hash Modulo (N %% M) Keys Remapped   : %d / %d (%.2f%% moved)\n",
		modMoved, numKeys, float64(modMoved)/float64(numKeys)*100)
	fmt.Printf("  Consistent Hashing Keys Remapped      : %d / %d (%.2f%% moved)\n",
		chMoved, numKeys, float64(chMoved)/float64(numKeys)*100)
	fmt.Printf("  Theoretical Minimal Relocation (1/N) : ~%.2f%%\n", 1.0/5.0*100)
}

func demoScatterGatherVsGSI() {
	fmt.Println("\n--- 4. Querying Non-Sharded Attributes: Scatter-Gather vs Global Secondary Index (GSI) ---")
	shards := []string{"shard-0", "shard-1", "shard-2", "shard-3"}
	router := sharding.NewConsistentHashRouter(100, shards)
	cluster := sharding.NewCluster(router)

	// Populate dataset
	for i := 0; i < 500; i++ {
		rec := sharding.Record{
			ID:       fmt.Sprintf("usr-%d", i),
			ShardKey: fmt.Sprintf("tenant-%d", i%50),
			Email:    fmt.Sprintf("user_%d@company.com", i),
			Data:     fmt.Sprintf("Profile data %d", i),
		}
		_ = cluster.Insert(rec)
	}

	targetEmail := "user_342@company.com"

	// 1. Scatter Gather Broadcast
	startSG := time.Now()
	sgRes := cluster.ScatterGatherBroadcast(func(r sharding.Record) bool {
		return r.Email == targetEmail
	})
	durSG := time.Since(startSG)

	fmt.Printf("Scatter-Gather Query (by Email without Shard Key):\n")
	fmt.Printf("  Nodes Broadcasted : %d / %d\n", sgRes.ShardsQueried, len(shards))
	fmt.Printf("  Records Matched   : %d\n", len(sgRes.Records))
	fmt.Printf("  Execution Time    : %v\n", durSG)

	// 2. Point lookup via GSI
	startGSI := time.Now()
	rec, err := cluster.GetByEmailUsingGSI(targetEmail, "usr-342")
	durGSI := time.Since(startGSI)

	fmt.Printf("Global Secondary Index (Lookup Vindex) Query:\n")
	fmt.Printf("  Nodes Broadcasted : 1 (Direct Point Lookup via Shard Key mapping)\n")
	fmt.Printf("  Record Found      : ID=%s, Email=%s, ShardKey=%s\n", rec.ID, rec.Email, rec.ShardKey)
	fmt.Printf("  Execution Time    : %v (GSI Avoided Broadcast Overhead!)\n", durGSI)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}
}

func demoIDGeneration() {
	fmt.Println("\n--- 5. Distributed Unique ID Generation: UUIDv7 vs Central Sequence Block Allocation ---")

	// UUIDv7
	uuid, _ := idgen.NewUUIDv7()
	fmt.Printf("Generated UUIDv7 (Time-Ordered 128-bit) : %s\n", uuid)

	// Block Allocator
	central := &idgen.MemoryCentralSequence{}
	allocator := idgen.NewSequenceBlockAllocator(5, central.AllocateBlock)

	fmt.Print("Sequence Block Allocator IDs (Block Size=5) : ")
	for i := 0; i < 8; i++ {
		id, _ := allocator.NextID()
		fmt.Printf("%d ", id)
	}
	fmt.Println()
}

func barChart(val, max int) string {
	width := 30
	ratio := float64(val) / float64(max)
	bars := int(math.Round(ratio * float64(width)))
	return "[" + strings.Repeat("#", bars) + strings.Repeat(" ", width-bars) + "]"
}
