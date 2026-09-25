package tests

import (
	"context"
	"testing"

	"compat/internal/compat"
)

func TestFullExpandMigrateContractLifecycle(t *testing.T) {
	// Baseline: Legacy store with existing historical records
	store := compat.NewMemoryStore()
	flags := compat.NewFeatureFlags()
	obs := compat.NewObservability()
	svc := compat.NewService(store, flags, obs)

	// Step 0: Baseline - Pre-migration historical users
	u1, err := svc.CreateUser("Historical User 1", "+628111111")
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	u2, err := svc.CreateUser("Historical User 2", "+628222222")
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// Verify legacy consumer can read
	leg1, err := svc.GetLegacyUser(u1.ID)
	if err != nil || leg1.Phone != "+628111111" {
		t.Fatalf("Legacy read failed: %v", err)
	}

	// Step 1: Expand Phase - Modern table exists (in store), Feature flag sets DualWrite
	flags.SetWriteMode(compat.WriteDual)

	// New user created during DualWrite
	u3, err := svc.CreateUser("DualWrite User 3", "+628333333", "+628333444")
	if err != nil {
		t.Fatalf("DualWrite creation failed: %v", err)
	}

	// Verify dual write wrote to both legacy and modern
	legacyRecord, err := svc.GetLegacyUser(u3.ID)
	if err != nil || legacyRecord.Phone != "+628333333" {
		t.Fatalf("Legacy record missing in dual write: %v", err)
	}
	modernRecord, err := svc.GetModernUser(u3.ID)
	if err != nil || len(modernRecord.Phones) < 2 {
		t.Fatalf("Modern record incomplete in dual write: %+v", modernRecord)
	}

	// Step 2: Migrate Phase - Historical Backfill
	backfiller := compat.NewBackfillWorker(store, obs, 10)
	migrated, err := backfiller.RunAll(context.Background())
	if err != nil {
		t.Fatalf("Backfill failed: %v", err)
	}
	if migrated != 2 { // u1 and u2 backfilled; u3 was already written via DualWrite
		t.Fatalf("Expected 2 backfilled users, got %d", migrated)
	}

	// Step 3: Verify Data Consistency before read switch
	drift, err := svc.ReconcileData()
	if err != nil || drift != 0 {
		t.Fatalf("Expected 0 drift before read switch, got %d (err: %v)", drift, err)
	}

	// Step 4: Migrate Phase - Switch Read Path
	flags.SetReadMode(compat.ReadNewOnly)

	// Modern consumer reads historical user u1
	modernU1, err := svc.GetModernUser(u1.ID)
	if err != nil || len(modernU1.Phones) != 1 || modernU1.Phones[0].Number != "+628111111" {
		t.Fatalf("Modern read of historical user failed: %+v", modernU1)
	}

	// Step 5: Contract Phase - Stop legacy write and drop legacy column
	err = svc.ApplyContract(true)
	if err != nil {
		t.Fatalf("ApplyContract failed: %v", err)
	}

	// Verify legacy column is dropped: legacy reads now fail
	_, err = svc.GetLegacyUser(u1.ID)
	if err == nil {
		t.Fatal("Expected error on legacy read after contract, got nil")
	}

	// Modern reads continue functioning flawlessly
	modernU2, err := svc.GetModernUser(u2.ID)
	if err != nil || len(modernU2.Phones) != 1 || modernU2.Phones[0].Number != "+628222222" {
		t.Fatalf("Modern read failed post-contract: %+v", modernU2)
	}
}

func TestRollbackScenarios(t *testing.T) {
	// Scenario A: Safe Rollback during Dual-Write
	store := compat.NewMemoryStore()
	flags := compat.NewFeatureFlags()
	obs := compat.NewObservability()
	svc := compat.NewService(store, flags, obs)

	// Step 1: Deploy DualWrite (N+1)
	flags.SetWriteMode(compat.WriteDual)
	u1, err := svc.CreateUser("RollbackUser", "+628777777")
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// Rollback triggered! Application rolled back to Version N (WriteLegacyOnly, ReadLegacyOnly)
	flags.SetWriteMode(compat.WriteLegacyOnly)
	flags.SetReadMode(compat.ReadLegacyOnly)

	// Legacy client reads record created during N+1: Data is safe!
	legacyU1, err := svc.GetLegacyUser(u1.ID)
	if err != nil || legacyU1.Phone != "+628777777" {
		t.Fatalf("Rollback safety violated: legacy data lost: %v", err)
	}

	// Scenario B: Unsafe Rollback after stopping Dual-Write prematurely
	flags.SetWriteMode(compat.WriteNewOnly)
	u2, err := svc.CreateUser("NewOnlyUser", "+628888888")
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// If rolled back to Version N now:
	flags.SetWriteMode(compat.WriteLegacyOnly)
	flags.SetReadMode(compat.ReadLegacyOnly)

	// Version N cannot see phone because legacy column was not populated in NewOnly mode!
	legacyU2, err := svc.GetLegacyUser(u2.ID)
	if err != nil {
		t.Fatalf("Unexpected read error: %v", err)
	}
	if legacyU2.Phone != "" {
		t.Fatalf("Expected empty legacy phone demonstrating data loss on premature rollback, got %q", legacyU2.Phone)
	}
}
