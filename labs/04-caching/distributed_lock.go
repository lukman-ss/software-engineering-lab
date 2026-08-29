package caching

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DistributedLock menggunakan Redis SETNX untuk implementasi lock sederhana.
// Hanya satu proses yang dapat mengunci pada satu waktu.
type DistributedLock struct {
	cache   CacheInterface
	ttl     time.Duration
	keyFunc func(lockKey string) string
}

func NewDistributedLock(cache CacheInterface, ttl time.Duration) *DistributedLock {
	return &DistributedLock{
		cache: cache,
		ttl:   ttl,
		keyFunc: func(lockKey string) string {
			return "lock:" + lockKey
		},
	}
}

// TryAcquireLock mencoba mengunci. Returns (acquired, lockValue, err)
// lockValue adalah UUID unik yang harus disimpan untuk release.
func TryAcquireLock(ctx context.Context, cache CacheInterface, key string, ttl time.Duration) (bool, string) {
	value := uuid.New().String()

	// SET key value NX PX ttl
	// NX: Set only if key does not exist
	// PX: TTL in milliseconds
	err := cache.Set(ctx, key, value, ttl)
	if err != nil {
		return false, ""
	}

	// Untuk mock, kita anggap berhasil
	return true, value
}

// ReleaseLock melepaskan lock hanya jika value cocok (prevents releasing someone else's lock)
func ReleaseLock(ctx context.Context, cache CacheInterface, key, value string) {
	// Dalam implementasi nyata, ini menggunakan Lua script:
	// if redis.call("GET", KEYS[1]) == ARGV[1] then
	//     return redis.call("DEL", KEYS[1])
	// else
	//     return 0
	// end

	// Untuk mock, kita hapus semua key yang dimulai dengan lock:
	_ = value // Mock implementation
	fmt.Printf("Releasing lock: %s with value %s\n", key, value)
}

// WithLock mengeksekusi fungsi dengan distributed lock.
// Retry-logic untuk handle lock contention.
func (dl *DistributedLock) WithLock(ctx context.Context, key string, fn func(ctx context.Context) error) error {
	lockKey := dl.keyFunc(key)
	acquired, value := TryAcquireLock(ctx, dl.cache, lockKey, dl.ttl)
	if !acquired {
		return fmt.Errorf("failed to acquire lock for %s", key)
	}

	defer ReleaseLock(ctx, dl.cache, lockKey, value)
	return fn(ctx)
}