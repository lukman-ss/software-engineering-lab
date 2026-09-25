package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"compat/internal/compat"
)

func main() {
	fmt.Println("=================================================================")
	fmt.Println("DEMO: BACKWARD COMPATIBILITY — EXPAND -> MIGRATE -> CONTRACT")
	fmt.Println("=================================================================")

	store := compat.NewMemoryStore()
	flags := compat.NewFeatureFlags()
	obs := compat.NewObservability()
	svc := compat.NewService(store, flags, obs)

	// -------------------------------------------------------------
	// STEP 1: Baseline (V1 Production State)
	// -------------------------------------------------------------
	fmt.Println("\n--- [STEP 1] Baseline Production State (Version N) ---")
	fmt.Println("Schema: users(id, name, phone)")
	u1, _ := svc.CreateUser("Alice", "+62811111111")
	u2, _ := svc.CreateUser("Bob", "+62822222222")
	fmt.Printf("Created historical users in legacy store: ID %d (%s), ID %d (%s)\n", u1.ID, u1.Name, u2.ID, u2.Name)

	legacyUser, err := svc.GetLegacyUser(u1.ID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("V1 Legacy Client reads user 1: {id: %d, name: %q, phone: %q}\n", legacyUser.ID, legacyUser.Name, legacyUser.Phone)

	// -------------------------------------------------------------
	// STEP 2: Expand Phase (Additive Schema & Payload)
	// -------------------------------------------------------------
	fmt.Println("\n--- [STEP 2] Expand Phase (Version N+1 Deployed) ---")
	fmt.Println("Action: user_phones table created. WriteMode set to WriteDual.")
	flags.SetWriteMode(compat.WriteDual)

	u3, _ := svc.CreateUser("Charlie", "+62833333333", "+62833334444")
	fmt.Printf("New user created via Dual-Write: ID %d (%s)\n", u3.ID, u3.Name)

	resp, _ := svc.GetUser(u3.ID)
	payload, _ := compat.SerializeResponse(resp)
	fmt.Printf("API Output (Enriched Additive Payload):\n%s\n", string(payload))

	var v1Client compat.LegacyConsumerDTO
	_ = json.Unmarshal(payload, &v1Client)
	fmt.Printf("V1 Legacy Client consumes payload: phone=%q (Compatibility Kept)\n", v1Client.Phone)

	var v2Client compat.ModernConsumerDTO
	_ = json.Unmarshal(payload, &v2Client)
	fmt.Printf("V2 Modern Client consumes payload: %d phones parsed (is_primary=%v)\n", len(v2Client.Phones), v2Client.Phones[0].IsPrimary)

	// -------------------------------------------------------------
	// STEP 3: Migrate Phase (Backfill & Fallback Read)
	// -------------------------------------------------------------
	fmt.Println("\n--- [STEP 3] Migrate Phase (Backfill & Data Reconciliation) ---")
	fmt.Println("Action: Running resumable batch backfill for historical data...")
	worker := compat.NewBackfillWorker(store, obs, 2)
	totalBackfilled, err := worker.RunAll(context.Background())
	if err != nil {
		fmt.Printf("Backfill error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Backfill worker complete: %d legacy records migrated to user_phones table.\n", totalBackfilled)

	driftCount, _ := svc.ReconcileData()
	fmt.Printf("Data Reconciliation check: detected %d drifting records.\n", driftCount)

	// -------------------------------------------------------------
	// STEP 4: Switch Read Path
	// -------------------------------------------------------------
	fmt.Println("\n--- [STEP 4] Switch Read Path (ReadMode: ReadNewOnly) ---")
	flags.SetReadMode(compat.ReadNewOnly)
	modernU1, _ := svc.GetModernUser(u1.ID)
	fmt.Printf("V2 Client reading historical User 1 from new schema: %+v\n", modernU1.Phones)

	// -------------------------------------------------------------
	// STEP 5: Rollback Demonstration
	// -------------------------------------------------------------
	fmt.Println("\n--- [STEP 5] Safe Rollback Demonstration ---")
	fmt.Println("Simulating rollback from N+1 back to Version N while in Dual-Write...")
	flags.SetWriteMode(compat.WriteLegacyOnly)
	flags.SetReadMode(compat.ReadLegacyOnly)
	rolledBackRead, err := svc.GetLegacyUser(u3.ID)
	if err != nil || rolledBackRead.Phone == "" {
		fmt.Println("Rollback failed! Data was lost.")
		os.Exit(1)
	}
	fmt.Printf("Rollback SUCCESS: Legacy instance read user 3's phone: %q (No data loss!)\n", rolledBackRead.Phone)

	// Restore dual write and finish migration
	flags.SetWriteMode(compat.WriteDual)
	flags.SetReadMode(compat.ReadNewOnly)

	// -------------------------------------------------------------
	// STEP 6: Contract Phase (Drop Legacy Schema & Fields)
	// -------------------------------------------------------------
	fmt.Println("\n--- [STEP 6] Contract Phase ---")
	fmt.Println("Simulating sunset of legacy interface (traffic to legacy interface = 0)...")
	err = svc.ApplyContract(true)
	if err != nil {
		fmt.Printf("Contract error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Contract applied successfully: legacy column dropped from users table.")

	_, err = svc.GetLegacyUser(u1.ID)
	fmt.Printf("Legacy read attempt post-contract: %v (Legacy safely retired)\n", err)

	finalModernRead, _ := svc.GetModernUser(u1.ID)
	fmt.Printf("Modern client read post-contract: %s with %d phones.\n", finalModernRead.Name, len(finalModernRead.Phones))

	// -------------------------------------------------------------
	// Observability Summary
	// -------------------------------------------------------------
	fmt.Println("\n--- [METRICS & OBSERVABILITY SNAPSHOT] ---")
	metrics := obs.Snapshot()
	for k, v := range metrics {
		fmt.Printf("- %s: %d\n", k, v)
	}
	fmt.Println("\nDEMO COMPLETED SUCCESSFULLY.")
}
