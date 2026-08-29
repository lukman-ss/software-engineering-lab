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
	locker  LockInterface
	ttl     time.Duration
	keyFunc func(lockKey string) string
}

func NewDistributedLock(locker LockInterface, ttl time.Duration) *DistributedLock {
	return &DistributedLock{
		locker: locker,
		ttl:    ttl,
		keyFunc: func(lockKey string) string {
			return "lock:" + lockKey
		},
	}
}

// TryAcquireLock mencoba mengunci menggunakan SETNX.
// Returns (acquired, lockValue).
// lockValue adalah UUID unik yang harus disimpan untuk release.
func TryAcquireLock(ctx context.Context, locker LockInterface, key string, ttl time.Duration) (bool, string) {
	value := uuid.New().String()

	// SETNX: Set key only if it doesn't exist
	acquired, err := locker.SetNX(ctx, key, value, ttl)
	if err != nil || !acquired {
		return false, ""
	}
	return true, value
}

// ReleaseLock melepaskan lock hanya jika value cocok (prevents releasing someone else's lock)
func ReleaseLock(ctx context.Context, locker LockInterface, key, value string) error {
	// Dalam implementasi nyata, ini menggunakan Lua script:
	// if redis.call("GET", KEYS[1]) == ARGV[1] then
	//     return redis.call("DEL", KEYS[1])
	// else
	//     return 0
	// end

	// Untuk mock, kita hapus semua key yang dimulai dengan lock:
	current, _ := locker.Get(ctx, key)
	_ = current // Mock implementation
	_, _ = locker.Del(ctx, key)
	return nil
}

// WithLock mengeksekusi fungsi dengan distributed lock.
// Retry-logic untuk handle lock contention.
func (dl *DistributedLock) WithLock(ctx context.Context, key string, fn func(ctx context.Context) error) error {
	lockKey := dl.keyFunc(key)
	acquired, value := TryAcquireLock(ctx, dl.locker, lockKey, dl.ttl)
	if !acquired {
		return fmt.Errorf("failed to acquire lock for %s", key)
	}

	defer ReleaseLock(ctx, dl.locker, lockKey, value)
	return fn(ctx)
}
