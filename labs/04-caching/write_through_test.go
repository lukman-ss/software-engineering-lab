package caching_test

import (
	"context"
	"encoding/json"
	"testing"

	caching "github.com/lukman-ss/software-engineering-lab/labs/04-caching"
	"github.com/DATA-DOG/go-sqlmock"
)

// TestWriteThroughSuccess memverifikasi update DB diikuti update cache.
func TestWriteThroughSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	cache := caching.NewMockCache()
	svc := caching.NewWriteThroughService(db, cache)
	ctx := context.Background()

	p := caching.Product{ID: "100", Name: "Updated Product", Price: 15.0}

	// Expect DB update
	mock.ExpectExec("UPDATE products SET name = \\$1, price = \\$2 WHERE id = \\$3").
		WithArgs(p.Name, p.Price, p.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = svc.UpdateProduct(ctx, p)
	if err != nil {
		t.Fatalf("UpdateProduct failed: %v", err)
	}

	// Verify Cache was updated (Write Through)
	key := caching.CacheKey("product", p.ID, 1)
	cached, err := cache.Get(ctx, key)
	if err != nil {
		t.Fatalf("cache should be updated: %v", err)
	}

	var cachedProduct caching.Product
	if err := json.Unmarshal([]byte(cached), &cachedProduct); err != nil {
		t.Fatalf("failed to unmarshal cached data: %v", err)
	}

	if cachedProduct.Name != "Updated Product" {
		t.Errorf("expected cache to have 'Updated Product', got %s", cachedProduct.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	t.Log("✓ Write Through success flow validated")
}

// TestWriteThroughCacheFailure memverifikasi jika DB sukses tapi Cache gagal,
// aplikasi menangani dengan aman (fallback delete cache).
func TestWriteThroughCacheFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	// Use failing cache (Set/Delete will return error)
	cache := caching.NewFailingMockCache()
	svc := caching.NewWriteThroughService(db, cache)
	ctx := context.Background()

	p := caching.Product{ID: "101", Name: "Failure Test", Price: 10.0}

	// Expect DB update (succeeds)
	mock.ExpectExec("UPDATE products SET name = \\$1, price = \\$2 WHERE id = \\$3").
		WithArgs(p.Name, p.Price, p.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Operation still succeeds because DB (Source of Truth) succeeded
	err = svc.UpdateProduct(ctx, p)
	if err != nil {
		t.Fatalf("UpdateProduct should succeed even if cache fails: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	t.Log("✓ Write Through handles cache failure safely (partial failure handled)")
}
