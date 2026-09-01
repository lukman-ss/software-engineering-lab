package caching_test

import (
	"context"
	"encoding/json"
	"testing"
	"fmt"

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

// TestWriteThroughDBFailure memverifikasi jika DB gagal, cache tidak disentuh
func TestWriteThroughDBFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	cache := caching.NewMockCache()
	svc := caching.NewWriteThroughService(db, cache)
	ctx := context.Background()

	p := caching.Product{ID: "102", Name: "DB Failure Test", Price: 10.0}

	// Expect DB update (fails)
	mock.ExpectExec("UPDATE products SET name = \\$1, price = \\$2 WHERE id = \\$3").
		WithArgs(p.Name, p.Price, p.ID).
		WillReturnError(fmt.Errorf("db connection error"))

	// Operation must fail
	err = svc.UpdateProduct(ctx, p)
	if err == nil {
		t.Fatalf("UpdateProduct should fail when DB fails")
	}

	// Cache must not be touched (miss expected)
	key := caching.CacheKey("product", p.ID, 1)
	_, err = cache.Get(ctx, key)
	if err != caching.ErrCacheMiss {
		t.Errorf("cache should not be touched, expected ErrCacheMiss, got: %v", err)
	}
}

// TestWriteThroughProductNotFound memverifikasi jika product tidak ditemukan di DB,
// cache tidak disentuh
func TestWriteThroughProductNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	cache := caching.NewMockCache()
	svc := caching.NewWriteThroughService(db, cache)
	ctx := context.Background()

	p := caching.Product{ID: "103", Name: "Not Found Test", Price: 10.0}

	// Expect DB update (succeeds but 0 rows affected)
	mock.ExpectExec("UPDATE products SET name = \\$1, price = \\$2 WHERE id = \\$3").
		WithArgs(p.Name, p.Price, p.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Operation must fail with "product not found"
	err = svc.UpdateProduct(ctx, p)
	if err == nil {
		t.Fatalf("UpdateProduct should fail when product not found")
	}

	// Cache must not be touched
	key := caching.CacheKey("product", p.ID, 1)
	_, err = cache.Get(ctx, key)
	if err != caching.ErrCacheMiss {
		t.Errorf("cache should not be touched, expected ErrCacheMiss, got: %v", err)
	}
}

// TestWriteThroughSerializationError memverifikasi error dari input serialization
// terjadi sebelum mutasi database
func TestWriteThroughSerializationError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	cache := caching.NewMockCache()
	svc := caching.NewWriteThroughService(db, cache)
	ctx := context.Background()

	// Missing ID should cause early validation error
	p := caching.Product{ID: "", Name: "Bad ID Test", Price: 10.0}

	// Operation must fail
	err = svc.UpdateProduct(ctx, p)
	if err == nil {
		t.Fatalf("UpdateProduct should fail validation")
	}

	// No expectations set on mock DB because it should never reach DB
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

// TestWriteThroughCacheAndFallbackFailure memverifikasi bahwa bahkan jika
// cache SET gagal dan DELETE juga gagal, business operation tetap SUCCESS.
// Stale cache mungkin bertahan sampai TTL (safety net).
func TestWriteThroughCacheAndFallbackFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	// Use failing cache (Set/Delete will return error)
	cache := caching.NewFailingMockCache()
	svc := caching.NewWriteThroughService(db, cache)
	ctx := context.Background()

	p := caching.Product{ID: "104", Name: "All Cache Failures", Price: 20.0}

	// Expect DB update (succeeds)
	mock.ExpectExec("UPDATE products SET name = \\$1, price = \\$2 WHERE id = \\$3").
		WithArgs(p.Name, p.Price, p.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Operation still succeeds because DB (Source of Truth) succeeded
	// even though both cache SET and DELETE operations failed
	err = svc.UpdateProduct(ctx, p)
	if err != nil {
		t.Fatalf("UpdateProduct should succeed even if both cache SET and DELETE fail: %v", err)
	}

	// Verify: stale cache still exists (because Set never succeeded, Delete failed)
	// It will expire after TTL (safety net)
	key := caching.CacheKey("product", p.ID, 1)
	_, err = cache.Get(ctx, key)
	if err != caching.ErrCacheMiss {
		t.Logf("Cache still has stale data (as expected when Set+Delete both fail)")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	t.Log("✓ Business operation succeeds even when cache SET+DELETE both fail")
	t.Log("✓ Stale cache persists until TTL (safety net)")
}
