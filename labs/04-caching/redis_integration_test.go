//go:build integration
// +build integration

package caching_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	caching "github.com/lukman-ss/software-engineering-lab/labs/04-caching"
)

// To run these tests:
// make lab-04-integration
// Requires a running Redis instance at REDIS_ADDR (default: localhost:6379)

func setupRedisCache(t *testing.T) *caching.RedisCache {
	// Let test run with default locally or configured via env
	if os.Getenv("REDIS_ADDR") == "" {
		t.Setenv("REDIS_ADDR", "localhost:6379")
	}

	cache := caching.NewRedisCacheWithTTL(5 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Wait, since integration tests must not silently pass CI if Redis is down,
	// we shouldn't use t.Skip. We should fail.
	err := cache.Set(ctx, "integration:test_ping", "pong", 1*time.Second)
	if err != nil {
		t.Fatalf("Redis integration test failed: Redis not available at %s: %v", os.Getenv("REDIS_ADDR"), err)
	}

	return cache
}

func TestRedisCache_SetAndGet(t *testing.T) {
	cache := setupRedisCache(t)
	defer cache.Close()
	ctx := context.Background()

	key := "integration:test:setget:" + uuid.New().String()
	t.Cleanup(func() { _ = cache.Delete(context.Background(), key) })

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

	// Cleanup is handled by t.Cleanup
}

func TestRedisCache_TTLExpiry(t *testing.T) {
	cache := setupRedisCache(t)
	defer cache.Close()
	ctx := context.Background()

	key := "integration:test:expiry:" + uuid.New().String()
	t.Cleanup(func() { _ = cache.Delete(context.Background(), key) })

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

	if !errors.Is(err, caching.ErrCacheMiss) {
		t.Errorf("Expected ErrCacheMiss, got %v", err)
	}
}

func TestRedisCache_DeleteAndInvalidation(t *testing.T) {
	cache := setupRedisCache(t)
	defer cache.Close()
	ctx := context.Background()

	key := "integration:test:delete:" + uuid.New().String()
	t.Cleanup(func() { _ = cache.Delete(context.Background(), key) })

	// Set
	cache.Set(ctx, key, "to_be_deleted", 10*time.Second)

	// Delete
	err := cache.Delete(ctx, key)
	if err != nil {
		t.Fatalf("Failed to delete: %v", err)
	}

	// Get should fail with ErrCacheMiss
	_, err = cache.Get(ctx, key)
	if !errors.Is(err, caching.ErrCacheMiss) {
		t.Errorf("Expected ErrCacheMiss after delete, got %v", err)
	}
}

func TestRedisCache_CacheMiss(t *testing.T) {
	cache := setupRedisCache(t)
	defer cache.Close()
	ctx := context.Background()

	// Get non-existent key
	_, err := cache.Get(ctx, "integration:test:missing_key")

	if !errors.Is(err, caching.ErrCacheMiss) {
		t.Errorf("Expected ErrCacheMiss, got %v", err)
	}
}

func TestRedisCache_Locking(t *testing.T) {
	cache := setupRedisCache(t)
	defer cache.Close()
	ctx := context.Background()

	key := "integration:test:lock:" + uuid.New().String()
	t.Cleanup(func() { _ = cache.Delete(context.Background(), key) })

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

	// Try release with wrong token should fail
	deleted, err := cache.CompareAndDel(ctx, key, "wrong_token")
	if err != nil {
		t.Fatalf("CompareAndDel error: %v", err)
	}
	if deleted {
		t.Fatal("CompareAndDel with wrong token deleted the lock")
	}

	// Release with correct token should succeed
	deleted, err = cache.CompareAndDel(ctx, key, "lock_value")
	if err != nil {
		t.Fatalf("CompareAndDel error: %v", err)
	}
	if !deleted {
		t.Fatal("CompareAndDel with correct token failed to delete")
	}

	// After release, should acquire again
	acquired, _ = cache.SetNX(ctx, key, "lock_value3", 5*time.Second)
	if !acquired {
		t.Fatal("Lock acquisition after safe release failed")
	}
}

func TestRedisCache_LockTTL(t *testing.T) {
	cache := setupRedisCache(t)
	defer cache.Close()
	ctx := context.Background()

	key := "integration:test:lock_ttl:" + uuid.New().String()
	t.Cleanup(func() { _ = cache.Delete(context.Background(), key) })

	// Acquire lock with short TTL
	acquired, err := cache.SetNX(ctx, key, "ttl_token", 500*time.Millisecond)
	if err != nil || !acquired {
		t.Fatalf("Failed to acquire lock: %v", err)
	}

	// Should not be able to acquire immediately
	acquired, _ = cache.SetNX(ctx, key, "other_token", 500*time.Millisecond)
	if acquired {
		t.Fatal("Should not acquire lock before TTL expires")
	}

	// Wait for TTL to expire safely
	time.Sleep(600 * time.Millisecond)

	// Should be able to acquire now
	acquired, err = cache.SetNX(ctx, key, "new_token", 500*time.Millisecond)
	if err != nil || !acquired {
		t.Fatalf("Failed to acquire lock after TTL expiration: %v", err)
	}
}
