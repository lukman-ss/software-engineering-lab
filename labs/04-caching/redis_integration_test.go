//go:build integration
// +build integration

package caching_test

import (
	"context"
	"os"
	"testing"
	"time"

	caching "github.com/lukman-ss/software-engineering-lab/labs/04-caching"
)

// To run these tests:
// go test -tags=integration ./...
// Requires a running Redis instance at REDIS_ADDR (default: localhost:6379)

func setupRedisCache(t *testing.T) *caching.RedisCache {
	// Let test run with default locally or configured via env
	if os.Getenv("REDIS_ADDR") == "" {
		os.Setenv("REDIS_ADDR", "localhost:6379")
	}

	cache := caching.NewRedisCacheWithTTL(5 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Ping to check if Redis is up, skip test if not
	// Since CacheInterface doesn't have Ping, we'll try a dummy Set/Get
	err := cache.Set(ctx, "test_ping", "pong", 1*time.Second)
	if err != nil {
		t.Skipf("Skipping Redis integration test: Redis not available at %s: %v", os.Getenv("REDIS_ADDR"), err)
	}

	return cache
}

func TestRedisCache_SetAndGet(t *testing.T) {
	cache := setupRedisCache(t)
	defer cache.Close()
	ctx := context.Background()

	key := "integration:test:setget"
	value := "hello_redis"

	// Set
	err := cache.Set(ctx, key, value, 10*time.Second)
	if err != nil {
		t.Fatalf("Failed to set: %v", err)
	}

	// Get
	got, err := cache.Get(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get: %v", err)
	}

	if got != value {
		t.Errorf("Expected %s, got %s", value, got)
	}

	// Cleanup
	_ = cache.Delete(ctx, key)
}

func TestRedisCache_TTLExpiry(t *testing.T) {
	cache := setupRedisCache(t)
	defer cache.Close()
	ctx := context.Background()

	key := "integration:test:expiry"
	value := "temporary_data"

	// Set with 1 second TTL
	err := cache.Set(ctx, key, value, 1*time.Second)
	if err != nil {
		t.Fatalf("Failed to set: %v", err)
	}

	// Immediate Get should work
	got, err := cache.Get(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get immediately: %v", err)
	}
	if got != value {
		t.Errorf("Expected %s, got %s", value, got)
	}

	// Wait for expiry
	time.Sleep(1500 * time.Millisecond)

	// Get should fail with ErrCacheMiss
	_, err = cache.Get(ctx, key)
	if err == nil {
		t.Fatal("Expected error after expiry, got nil")
	}

	if err != caching.ErrCacheMiss {
		t.Errorf("Expected ErrCacheMiss, got %v", err)
	}
}

func TestRedisCache_DeleteAndInvalidation(t *testing.T) {
	cache := setupRedisCache(t)
	defer cache.Close()
	ctx := context.Background()

	key := "integration:test:delete"

	// Set
	cache.Set(ctx, key, "to_be_deleted", 10*time.Second)

	// Delete
	err := cache.Delete(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete: %v", err)
	}

	// Get should fail with ErrCacheMiss
	_, err = cache.Get(ctx, key)
	if err != caching.ErrCacheMiss {
		t.Errorf("Expected ErrCacheMiss after delete, got %v", err)
	}
}

func TestRedisCache_CacheMiss(t *testing.T) {
	cache := setupRedisCache(t)
	defer cache.Close()
	ctx := context.Background()

	// Get non-existent key
	_, err := cache.Get(ctx, "integration:test:missing_key")

	if err != caching.ErrCacheMiss {
		t.Errorf("Expected ErrCacheMiss, got %v", err)
	}
}

func TestRedisCache_Locking(t *testing.T) {
	cache := setupRedisCache(t)
	defer cache.Close()
	ctx := context.Background()

	key := "integration:test:lock"

	// First acquire should succeed
	acquired, err := cache.SetNX(ctx, key, "lock_value", 5*time.Second)
	if err != nil {
		t.Fatalf("Failed SetNX: %v", err)
	}
	if !acquired {
		t.Fatal("First lock acquisition failed")
	}

	// Second acquire should fail
	acquired, err = cache.SetNX(ctx, key, "lock_value2", 5*time.Second)
	if err != nil {
		t.Fatalf("Failed second SetNX: %v", err)
	}
	if acquired {
		t.Fatal("Second lock acquisition succeeded, expected failure")
	}

	// Delete lock
	deleted, err := cache.Del(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete lock: %v", err)
	}
	if !deleted {
		t.Fatal("Delete reported lock didn't exist")
	}
}
