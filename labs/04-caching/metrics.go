package caching

import (
	"sync/atomic"
	"time"
)

// CacheMetrics adalah counter sederhana untuk monitoring cache behavior.
// Production: ganti dengan Prometheus/OpenTelemetry.
//
// Tiga kategori metrik:
//  1. COUNTERS: hit, miss, error, rebuild_attempt, rebuild_success, db_query, lock_wait
//  2. LATENCY: cache_get, cache_set, db_fallback, rebuild
//  3. CAPACITY: memory_usage_pct, evicted_keys, expired_keys, key_count
//
// Derived metrics: hit_ratio, cache_error_rate, rebuild_success_rate
//
// Catatan: Threshold alert (seperti hit_ratio < 50%) bersifat illustrative
// dan harus disesuaikan dengan baseline serta SLO service.
type CacheMetrics struct {
	// COUNTERS
	hits             atomic.Int64
	misses           atomic.Int64
	errors           atomic.Int64
	rebuildAttempts  atomic.Int64
	rebuildSuccesses atomic.Int64
	dbQueries        atomic.Int64
	lockWaits        atomic.Int64
	dbFallbacks      atomic.Int64

	// LATENCY (total nanoseconds untuk menghitung average)
	cacheGetLatency   atomic.Int64
	cacheSetLatency   atomic.Int64
	dbFallbackLatency atomic.Int64
	rebuildLatency    atomic.Int64

	// CAPACITY
	evictedKeys atomic.Int64
	expiredKeys atomic.Int64
}

func NewCacheMetrics() *CacheMetrics {
	return &CacheMetrics{}
}

// --- COUNTERS ---

func (m *CacheMetrics) IncHit()            { m.hits.Add(1) }
func (m *CacheMetrics) IncMiss()           { m.misses.Add(1) }
func (m *CacheMetrics) IncError()          { m.errors.Add(1) }
func (m *CacheMetrics) IncRebuildAttempt() { m.rebuildAttempts.Add(1) }
func (m *CacheMetrics) IncRebuildSuccess() { m.rebuildSuccesses.Add(1) }
func (m *CacheMetrics) IncDBQuery()        { m.dbQueries.Add(1) }
func (m *CacheMetrics) IncLockWait()       { m.lockWaits.Add(1) }
func (m *CacheMetrics) IncDBFallback()     { m.dbFallbacks.Add(1) }
func (m *CacheMetrics) IncEvictedKey()     { m.evictedKeys.Add(1) }
func (m *CacheMetrics) IncExpiredKey()     { m.expiredKeys.Add(1) }

// --- LATENCY ---

func (m *CacheMetrics) RecordCacheGetLatency(d time.Duration) { m.cacheGetLatency.Add(d.Nanoseconds()) }
func (m *CacheMetrics) RecordCacheSetLatency(d time.Duration) { m.cacheSetLatency.Add(d.Nanoseconds()) }
func (m *CacheMetrics) RecordDBFallbackLatency(d time.Duration) {
	m.dbFallbackLatency.Add(d.Nanoseconds())
}
func (m *CacheMetrics) RecordRebuildLatency(d time.Duration) { m.rebuildLatency.Add(d.Nanoseconds()) }

// --- GETTERS ---

func (m *CacheMetrics) Hits() int64             { return m.hits.Load() }
func (m *CacheMetrics) Misses() int64           { return m.misses.Load() }
func (m *CacheMetrics) Errors() int64           { return m.errors.Load() }
func (m *CacheMetrics) RebuildAttempts() int64  { return m.rebuildAttempts.Load() }
func (m *CacheMetrics) RebuildSuccesses() int64 { return m.rebuildSuccesses.Load() }
func (m *CacheMetrics) DBQueries() int64        { return m.dbQueries.Load() }
func (m *CacheMetrics) LockWaits() int64        { return m.lockWaits.Load() }
func (m *CacheMetrics) DBFallbacks() int64      { return m.dbFallbacks.Load() }
func (m *CacheMetrics) EvictedKeys() int64      { return m.evictedKeys.Load() }
func (m *CacheMetrics) ExpiredKeys() int64      { return m.expiredKeys.Load() }

// --- DERIVED METRICS ---

// HitRatio menghitung persentase hit dari total (hit + miss).
// Catatan: Hit ratio tinggi ≠ cache pasti bernilai.
// ROI cache harus dilihat bersama cost query yang dihindari, cache latency,
// memory cost, invalidation complexity, dan failure amplification.
func (m *CacheMetrics) HitRatio() float64 {
	total := m.hits.Load() + m.misses.Load()
	if total == 0 {
		return 0.0
	}
	return float64(m.hits.Load()) / float64(total) * 100.0
}

// CacheErrorRate menghitung persentase error dari total operasi cache.
func (m *CacheMetrics) CacheErrorRate() float64 {
	total := m.hits.Load() + m.misses.Load() + m.errors.Load()
	if total == 0 {
		return 0.0
	}
	return float64(m.errors.Load()) / float64(total) * 100.0
}

// RebuildSuccessRate menghitung persentase rebuild yang berhasil.
func (m *CacheMetrics) RebuildSuccessRate() float64 {
	attempts := m.rebuildAttempts.Load()
	if attempts == 0 {
		return 0.0
	}
	return float64(m.rebuildSuccesses.Load()) / float64(attempts) * 100.0
}

// AverageCacheGetLatency mengembalikan rata-rata latency cache GET.
func (m *CacheMetrics) AverageCacheGetLatency() time.Duration {
	ops := m.hits.Load() + m.misses.Load() + m.errors.Load()
	if ops == 0 {
		return 0
	}
	return time.Duration(m.cacheGetLatency.Load() / ops)
}

// AverageDBFallbackLatency mengembalikan rata-rata latency fallback ke DB.
func (m *CacheMetrics) AverageDBFallbackLatency() time.Duration {
	if m.dbFallbacks.Load() == 0 {
		return 0
	}
	return time.Duration(m.dbFallbackLatency.Load() / m.dbFallbacks.Load())
}

// --- BACKWARD COMPATIBILITY ---

func (m *CacheMetrics) IncRebuild()           { m.IncRebuildAttempt() }
func (m *CacheMetrics) Rebuilds() int64       { return m.RebuildAttempts() }
func (m *CacheMetrics) RebuildSuccess() int64 { return m.RebuildSuccesses() }
