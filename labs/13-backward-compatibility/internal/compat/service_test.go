package compat

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSerializationBackwardCompatibility(t *testing.T) {
	svc := NewService(nil, nil, nil)
	svc.Flags().SetWriteMode(WriteDual)

	userResp, err := svc.CreateUser("Lukman", "+628123456", "+628765432")
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	payload, err := SerializeResponse(userResp)
	if err != nil {
		t.Fatalf("Serialize failed: %v", err)
	}

	// 1. Legacy consumer unmarshals without errors
	var legacy LegacyConsumerDTO
	if err := json.Unmarshal(payload, &legacy); err != nil {
		t.Fatalf("Legacy consumer crashed during unmarshal: %v", err)
	}
	if legacy.Phone != "+628123456" {
		t.Errorf("Expected phone %q, got %q", "+628123456", legacy.Phone)
	}

	// 2. Modern consumer unmarshals new phones array
	var modern ModernConsumerDTO
	if err := json.Unmarshal(payload, &modern); err != nil {
		t.Fatalf("Modern consumer crashed during unmarshal: %v", err)
	}
	if len(modern.Phones) != 2 {
		t.Fatalf("Expected 2 phones, got %d", len(modern.Phones))
	}
	if modern.Phones[0].Number != "+628123456" || !modern.Phones[0].IsPrimary {
		t.Errorf("Primary phone mismatch: %+v", modern.Phones[0])
	}
}

func TestBackfillIdempotentAndResumable(t *testing.T) {
	svc := NewService(nil, nil, nil)
	svc.Flags().SetWriteMode(WriteLegacyOnly)

	// Create 10 legacy users
	for i := 1; i <= 10; i++ {
		_, err := svc.CreateUser("LegacyUser", "+628000000")
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
	}

	// Worker with batchSize = 3
	worker := NewBackfillWorker(svc.Store(), svc.Obs(), 3)

	// Run batch 1
	n, done, err := worker.RunBatch(context.Background())
	if err != nil || done || n != 3 {
		t.Fatalf("Batch 1 unexpected: n=%d, done=%v, err=%v", n, done, err)
	}

	// Run batch 2
	n, done, err = worker.RunBatch(context.Background())
	if err != nil || done || n != 3 {
		t.Fatalf("Batch 2 unexpected: n=%d, done=%v, err=%v", n, done, err)
	}

	// Finish remaining batches
	total, err := worker.RunAll(context.Background())
	if err != nil {
		t.Fatalf("RunAll failed: %v", err)
	}
	if total != 4 { // remaining 4 records
		t.Errorf("Expected 4 remaining records migrated, got %d", total)
	}

	// Check idempotency: running again should migrate 0
	secondRun, err := worker.RunAll(context.Background())
	if err != nil {
		t.Fatalf("Second RunAll failed: %v", err)
	}
	if secondRun != 0 {
		t.Errorf("Expected 0 migrations on second run, got %d", secondRun)
	}
}

func TestFallbackRead(t *testing.T) {
	svc := NewService(nil, nil, nil)
	svc.Flags().SetWriteMode(WriteLegacyOnly)

	u, err := svc.CreateUser("Alice", "+628999999")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Read mode: Fallback
	svc.Flags().SetReadMode(ReadFallback)

	// Prior to backfill, user_phones is empty
	phonesBefore, _ := svc.Store().GetPhones(u.ID)
	if len(phonesBefore) != 0 {
		t.Fatalf("Expected 0 phones in modern store before read, got %d", len(phonesBefore))
	}

	// Read via GetUser triggers fallback read and lazy backfill
	resp, err := svc.GetUser(u.ID)
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if resp.Phone != "+628999999" {
		t.Errorf("Expected fallback phone '+628999999', got %q", resp.Phone)
	}

	// Now modern store has been lazily populated
	phonesAfter, _ := svc.Store().GetPhones(u.ID)
	if len(phonesAfter) != 1 || phonesAfter[0].Number != "+628999999" {
		t.Errorf("Lazy backfill did not populate modern store: %+v", phonesAfter)
	}
}

func TestDataReconciliationAndDrift(t *testing.T) {
	svc := NewService(nil, nil, nil)
	svc.Flags().SetWriteMode(WriteLegacyOnly)

	// Create legacy user (un-backfilled)
	_, _ = svc.CreateUser("Bob", "+628111111")

	// Drift expected because user_phones has no corresponding entry
	drifts, err := svc.ReconcileData()
	if err != nil {
		t.Fatalf("ReconcileData failed: %v", err)
	}
	if drifts != 1 {
		t.Errorf("Expected 1 drift, got %d", drifts)
	}

	// Run backfill
	_, _ = svc.BackfillWorker().RunAll(context.Background())

	// Drifts should now be 0
	drifts, err = svc.ReconcileData()
	if err != nil {
		t.Fatalf("ReconcileData failed: %v", err)
	}
	if drifts != 0 {
		t.Errorf("Expected 0 drifts after backfill, got %d", drifts)
	}
}

func TestDeprecationHeadersAndContractEnforcement(t *testing.T) {
	svc := NewService(nil, nil, nil)
	svc.Flags().SetWriteMode(WriteDual)
	u, _ := svc.CreateUser("Charlie", "+628333333")

	handler := NewAPIHandler(svc)

	// 1. Hit legacy endpoint: expect Deprecation & Sunset headers
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users?id=1", nil)
	w := httptest.NewRecorder()
	handler.GetUserV1(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}
	if w.Header().Get("Deprecation") != "true" {
		t.Errorf("Missing Deprecation header")
	}
	if w.Header().Get("Sunset") == "" {
		t.Errorf("Missing Sunset header")
	}

	// 2. Attempt contract while legacy hits > 0 without force: should fail
	err := svc.ApplyContract(false)
	if err == nil {
		t.Fatal("Expected error applying contract while legacy traffic active, got nil")
	}

	// 3. Force contract (simulating after 30-day zero legacy traffic window)
	err = svc.ApplyContract(true)
	if err != nil {
		t.Fatalf("ApplyContract(true) failed: %v", err)
	}

	// 4. Hit legacy endpoint again: expect 410 Gone
	w2 := httptest.NewRecorder()
	handler.GetUserV1(w2, req)
	if w2.Code != http.StatusGone {
		t.Errorf("Expected 410 Gone after contract, got %d", w2.Code)
	}

	// 5. Hit modern endpoint: works as expected
	reqV2 := httptest.NewRequest(http.MethodGet, "/api/v2/users?id=1", nil)
	w3 := httptest.NewRecorder()
	handler.GetUserV2(w3, reqV2)
	if w3.Code != http.StatusOK {
		t.Errorf("Expected 200 OK for v2 endpoint, got %d", w3.Code)
	}
	_ = u
}
