// Package caching explores real-world caching patterns and their trade-offs.
package caching

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache implements CacheInterface using real Redis.
type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

const defaultRedisTTL = 5 * time.Minute

// NewRedisCache creates a Redis-backed cache from environment variables.
// Required env vars: REDIS_ADDR, REDIS_PASSWORD, REDIS_DB
func NewRedisCache() (*RedisCache, error) {
	return NewRedisCacheFromEnv(defaultRedisTTL)
}

// NewRedisCacheFromEnv creates a Redis-backed cache from environment variables with specific default TTL.
func NewRedisCacheFromEnv(ttl time.Duration) (*RedisCache, error) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	password := os.Getenv("REDIS_PASSWORD")
	db := 0
	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		parsedDB, err := strconv.Atoi(dbStr)
		if err != nil {
			return nil, fmt.Errorf("invalid REDIS_DB format %q: %w", dbStr, err)
		}
		db = parsedDB
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCache{client: client, ttl: ttl}, nil
}

// NewRedisCacheWithTTL creates a Redis-backed cache bypassing env validation errors (legacy wrapper).
// Will panic if REDIS_DB is invalid.
func NewRedisCacheWithTTL(ttl time.Duration) *RedisCache {
	cache, err := NewRedisCacheFromEnv(ttl)
	if err != nil {
		panic(err)
	}
	return cache
}

// NewRedisClient creates a RedisCache with explicit client (for testing) with default TTL.
func NewRedisClient(client *redis.Client) *RedisCache {
	return &RedisCache{client: client, ttl: defaultRedisTTL}
}

// Get retrieves a value from Redis.
func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", ErrCacheMiss
		}
		return "", fmt.Errorf("redis get: %w", err)
	}
	return val, nil
}

// Set stores a value in Redis with the given TTL.
func (c *RedisCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = c.ttl
	}
	err := c.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("redis set: %w", err)
	}
	return nil
}

// Delete removes a key from Redis.
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("redis del: %w", err)
	}
	return nil
}

// SetNX sets a key only if it does not exist (for distributed locking).
func (c *RedisCache) SetNX(ctx context.Context, key, value string, ttl time.Duration) (bool, error) {
	success, err := c.client.SetNX(ctx, key, value, ttl).Result()
	if err != nil {
		return false, fmt.Errorf("redis setnx: %w", err)
	}
	return success, nil
}

// CompareAndDel atomically deletes only if value matches via Lua script.
func (c *RedisCache) CompareAndDel(ctx context.Context, key, value string) (bool, error) {
	script := `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end`
	result, err := c.client.Eval(ctx, script, []string{key}, value).Result()
	if err != nil {
		return false, fmt.Errorf("redis eval: %w", err)
	}

	if val, ok := result.(int64); ok {
		return val == 1, nil
	}
	return false, nil
}

// GetWithExpiry retrieves a value and its remaining TTL.
func (c *RedisCache) GetWithExpiry(ctx context.Context, key string) (string, time.Time, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", time.Time{}, ErrCacheMiss
		}
		return "", time.Time{}, fmt.Errorf("redis get: %w", err)
	}

	ttl, err := c.client.TTL(ctx, key).Result()
	if err != nil {
		// TTL check failed, return value without expiry
		return val, time.Time{}, nil
	}

	// Handle Redis TTL edge cases
	if ttl == -2 {
		// Key does not exist (expired between GET and TTL)
		return "", time.Time{}, ErrCacheMiss
	}
	if ttl == -1 {
		// Key exists but has no associated expire
		return val, time.Time{}, nil // zero value time.Time signals no expiry
	}

	expiry := time.Now().Add(ttl)
	return val, expiry, nil
}

// Close closes the Redis connection.
func (c *RedisCache) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// Ensure interface is satisfied
var _ CacheInterface = (*RedisCache)(nil)
var _ LockInterface = (*RedisCache)(nil)
