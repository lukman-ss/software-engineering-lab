package caching_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	caching "github.com/lukman-ss/software-engineering-lab/labs/04-caching"
)

// TestWriteThroughSuccess memverifikasi update DB dengan RETURNING diikuti update cache.
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

	// Expect DB update with RETURNING
	mock.ExpectQuery("UPDATE products SET name = \\$1, price = \\$2 WHERE id = \\$3 RETURNING id, name, price").
		WithArgs(p.Name, p.Price, p.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).
			AddRow(p.ID, p.Name, p.Price))

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

	t.Log("✓ Write Through success flow validated (RETURNING guarantees authoritative value)")
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

	// Expect DB update with RETURNING (succeeds)
	mock.ExpectQuery("UPDATE products SET name = \\$1, price = \\$2 WHERE id = \\$3 RETURNING id, name, price").
		WithArgs(p.Name, p.Price, p.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).
			AddRow(p.ID, p.Name, p.Price))

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
	mock.ExpectQuery("UPDATE products SET name = \\$1, price = \\$2 WHERE id = \\$3 RETURNING id, name, price").
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

	// Expect DB update (no rows = sql.ErrNoRows)
	mock.ExpectQuery("UPDATE products SET name = \\$1, price = \\$2 WHERE id = \\$3 RETURNING id, name, price").
		WithArgs(p.Name, p.Price, p.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"})) // empty rows → ErrNoRows

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

// TestWriteThroughEmptyIDValidation memverifikasi validasi ID kosong sebelum DB
func TestWriteThroughEmptyIDValidation(t *testing.T) {
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

	// Expect DB update with RETURNING (succeeds)
	mock.ExpectQuery("UPDATE products SET name = \\$1, price = \\$2 WHERE id = \\$3 RETURNING id, name, price").
		WithArgs(p.Name, p.Price, p.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).
			AddRow(p.ID, p.Name, p.Price))

	// Operation still succeeds because DB (Source of Truth) succeeded
	// even though both cache SET and DELETE operations failed
	err = svc.UpdateProduct(ctx, p)
	if err != nil {
		t.Fatalf("UpdateProduct should succeed even if both cache SET and DELETE fail: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}

	t.Log("✓ Business operation succeeds even when cache SET+DELETE both fail")
	t.Log("✓ Stale cache persists until TTL (safety net)")
}

// TestWriteThroughGetProduct memverifikasi GetProduct menangani: ErrCacheMiss, nil error + empty value, backend error, valid hit, corrupt JSON.
func TestWriteThroughGetProduct(t *testing.T) {
	ctx := context.Background()

	t.Run("ErrCacheMiss", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		metrics := caching.NewCacheMetrics()
		cache := caching.NewMockCache()
		svc := caching.NewWriteThroughServiceWithMetrics(db, cache, metrics)

		mock.ExpectQuery("SELECT id, name, price FROM products WHERE id =").
			WithArgs("p1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).
				AddRow("p1", "Item 1", 50.0))

		p, err := svc.GetProduct(ctx, "p1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.ID != "p1" || p.Name != "Item 1" {
			t.Errorf("unexpected product: %+v", p)
		}
		if metrics.Misses() != 1 {
			t.Errorf("expected 1 miss, got %d", metrics.Misses())
		}
		if metrics.Errors() != 0 {
			t.Errorf("expected 0 errors, got %d", metrics.Errors())
		}
		if metrics.DBFallbacks() != 0 {
			t.Errorf("expected 0 db fallbacks for miss, got %d", metrics.DBFallbacks())
		}
	})

	t.Run("NilErrorAndEmptyValue", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		metrics := caching.NewCacheMetrics()
		cache := caching.NewMockCache()
		cache.Set(ctx, caching.CacheKey("product", "p2", 1), "", time.Minute)
		metrics.Reset()

		svc := caching.NewWriteThroughServiceWithMetrics(db, cache, metrics)

		mock.ExpectQuery("SELECT id, name, price FROM products WHERE id =").
			WithArgs("p2").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).
				AddRow("p2", "Item 2", 75.0))

		p, err := svc.GetProduct(ctx, "p2")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.ID != "p2" || p.Name != "Item 2" {
			t.Errorf("unexpected product: %+v", p)
		}
		if metrics.Misses() != 1 {
			t.Errorf("expected 1 miss for empty cached value, got %d", metrics.Misses())
		}
		if metrics.Errors() != 0 {
			t.Errorf("expected 0 errors for empty cached value, got %d", metrics.Errors())
		}
		if metrics.DBFallbacks() != 0 {
			t.Errorf("expected 0 db fallbacks for empty cached value, got %d", metrics.DBFallbacks())
		}
	})

	t.Run("BackendError", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		metrics := caching.NewCacheMetrics()
		cache := caching.NewFailingMockCache()
		svc := caching.NewWriteThroughServiceWithMetrics(db, cache, metrics)

		mock.ExpectQuery("SELECT id, name, price FROM products WHERE id =").
			WithArgs("p3").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).
				AddRow("p3", "Item 3", 100.0))

		p, err := svc.GetProduct(ctx, "p3")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.ID != "p3" {
			t.Errorf("unexpected product ID: %s", p.ID)
		}
		if metrics.Errors() == 0 {
			t.Errorf("expected >0 errors for backend failure")
		}
		if metrics.DBFallbacks() != 1 {
			t.Errorf("expected 1 DB fallback for backend error, got %d", metrics.DBFallbacks())
		}
	})

	t.Run("ValidHit", func(t *testing.T) {
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		metrics := caching.NewCacheMetrics()
		cache := caching.NewMockCache()
		svc := caching.NewWriteThroughServiceWithMetrics(db, cache, metrics)

		product := caching.Product{ID: "p4", Name: "Item 4", Price: 120.0}
		data, _ := json.Marshal(product)
		cache.Set(ctx, caching.CacheKey("product", "p4", 1), string(data), time.Minute)
		metrics.Reset()

		p, err := svc.GetProduct(ctx, "p4")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.Name != "Item 4" {
			t.Errorf("expected Item 4, got %s", p.Name)
		}
		if metrics.Hits() != 1 {
			t.Errorf("expected 1 hit, got %d", metrics.Hits())
		}
	})

	t.Run("CorruptJSON", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock db: %v", err)
		}
		defer db.Close()

		metrics := caching.NewCacheMetrics()
		cache := caching.NewMockCache()
		cache.Set(ctx, caching.CacheKey("product", "p5", 1), "{corrupt_json}", time.Minute)
		metrics.Reset()

		svc := caching.NewWriteThroughServiceWithMetrics(db, cache, metrics)

		mock.ExpectQuery("SELECT id, name, price FROM products WHERE id =").
			WithArgs("p5").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price"}).
				AddRow("p5", "Item 5", 150.0))

		p, err := svc.GetProduct(ctx, "p5")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.ID != "p5" {
			t.Errorf("unexpected product: %+v", p)
		}
		if metrics.Errors() != 1 {
			t.Errorf("expected 1 error for corrupt json, got %d", metrics.Errors())
		}
	})
}
