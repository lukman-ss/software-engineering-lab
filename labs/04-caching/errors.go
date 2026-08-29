package caching

import "errors"

// Sentinel errors for cache operations.
// Digunakan untuk membedakan antara cache miss, cache error, dan database error.
var (
	// ErrCacheMiss indicates the cache does not contain the requested key.
	// This is NOT an error condition - it's expected behavior.
	ErrCacheMiss = errors.New("cache miss")

	// ErrCacheDown indicates the cache service is unavailable.
	// Trigger fallback to database.
	ErrCacheDown = errors.New("cache service unavailable")

	// ErrCacheCorrupt indicates the cached value is malformed.
	// Delete and rebuild from source of truth.
	ErrCacheCorrupt = errors.New("cache value corrupt")

	// ErrDatabaseError indicates a database failure.
	// This is a real error - not a cache issue.
	ErrDatabaseError = errors.New("database error")
)

// IsCacheMiss returns true if err is a cache miss (not a failure).
func IsCacheMiss(err error) bool {
	return err != nil && err == ErrCacheMiss
}

// IsCacheDown returns true if cache is unavailable, allowing fallback.
func IsCacheDown(err error) bool {
	return err != nil && err == ErrCacheDown
}

// IsCacheCorrupt returns true if cached data is malformed.
func IsCacheCorrupt(err error) bool {
	return err != nil && err == ErrCacheCorrupt
}
