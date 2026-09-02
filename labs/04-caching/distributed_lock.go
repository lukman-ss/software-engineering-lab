package caching

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DistributedLock menggunakan Redis SETNX untuk implementasi lock sederhana.
// Redis distributed lock digunakan agar hanya satu holder/client/application
// instance yang memegang lock untuk key tertentu pada saat yang sama.
//
// PENTING - Lease Expiry Problem:
// Jika eksekusi lebih lama dari TTL lock:
// 1. Holder A memegang lock (TTL 5s)
// 2. Job A berjalan lambat (> 5s)
// 3. Lock expires dari Redis
// 4. Holder B memperoleh lock baru
// 5. A dan B sekarang bekerja bersamaan (duplicate rebuild)
//
// Oleh karena itu, Distributed Lock pada lab ini digunakan murni untuk
// MENGURANGI duplicate cache rebuild (optimization), BUKAN sebagai
// strict concurrency control atau correctness guarantee untuk transaksi DB.
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
// Returns (acquired, lockValue, error).
// lockValue adalah UUID unik yang harus disimpan untuk release.
func TryAcquireLock(ctx context.Context, locker LockInterface, key string, ttl time.Duration) (bool, string, error) {
	if ttl <= 0 {
		return false, "", fmt.Errorf("distributed lock requires positive TTL")
	}

	value := uuid.New().String()

	// SETNX: Set key only if it doesn't exist
	acquired, err := locker.SetNX(ctx, key, value, ttl)
	if err != nil {
		// Technical error (Redis down, context timeout, network)
		return false, "", err
	}
	if !acquired {
		// Normal lock contention (someone else holds it)
		return false, "", nil
	}

	return true, value, nil
}

// ReleaseLock melepaskan lock hanya jika value cocok (prevents releasing someone else's lock)
func ReleaseLock(ctx context.Context, locker LockInterface, key, value string) error {
	// Use atomic compare-and-delete to prevent releasing someone else's lock
	deleted, err := locker.CompareAndDel(ctx, key, value)
	if err != nil {
		return err
	}
	if !deleted {
		// Lock was held by someone else (lease expired and acquired by B) or didn't exist
		return fmt.Errorf("lock not released: value mismatch or key missing")
	}
	return nil
}

// WithLock mengeksekusi fungsi dengan distributed lock.
// Mencoba acquire satu kali (try-once), tanpa retry/backoff.
//
// Jika gagal karena contention atau error teknis, return error.
// Caller harus meng-handle fallback (misal: serve stale cache, atau error).
func (dl *DistributedLock) WithLock(ctx context.Context, key string, fn func(ctx context.Context) error) (err error) {
	lockKey := dl.keyFunc(key)

	acquired, value, acquireErr := TryAcquireLock(ctx, dl.locker, lockKey, dl.ttl)
	if acquireErr != nil {
		return fmt.Errorf("technical error acquiring lock for %s: %w", key, acquireErr)
	}
	if !acquired {
		return fmt.Errorf("failed to acquire lock for %s: contention", key)
	}

	defer func() {
		// Gunakan context baru tanpa cancellation untuk safe release.
		// Jika ctx caller sudah di-cancel, Redis call dengan ctx lama
		// akan gagal saat mencoba me-release lock.
		releaseCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()

		releaseErr := ReleaseLock(releaseCtx, dl.locker, lockKey, value)
		if releaseErr != nil {
			err = errors.Join(err, releaseErr)
		}
	}()

	return fn(ctx)
}
