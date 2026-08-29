package caching

import (
	"sync/atomic"
)

// CacheMetrics adalah counter sederhana untuk monitoring cache behavior.
// Production: ganti dengan Prometheus/OpenTelemetry.
type CacheMetrics struct {
	hits      atomic.Int64
	misses    atomic.Int64
	errors    atomic.Int64
	rebuilds  atomic.Int64
	dbQueries atomic.Int64
	lockWaits atomic.Int64
}

func NewCacheMetrics() *CacheMetrics {
	return &CacheMetrics{}
}

func (m *CacheMetrics) IncHit()      { m.hits.Add(1) }
func (m *CacheMetrics) IncMiss()     { m.misses.Add(1) }
func (m *CacheMetrics) IncError()    { m.errors.Add(1) }
func (m *CacheMetrics) IncRebuild()  { m.rebuilds.Add(1) }
func (m *CacheMetrics) IncDBQuery()  { m.dbQueries.Add(1) }
func (m *CacheMetrics) IncLockWait() { m.lockWaits.Add(1) }

func (m *CacheMetrics) Hits() int64      { return m.hits.Load() }
func (m *CacheMetrics) Misses() int64    { return m.misses.Load() }
func (m *CacheMetrics) Errors() int64    { return m.errors.Load() }
func (m *CacheMetrics) Rebuilds() int64  { return m.rebuilds.Load() }
func (m *CacheMetrics) DBQueries() int64 { return m.dbQueries.Load() }
func (m *CacheMetrics) LockWaits() int64 { return m.lockWaits.Load() }

// HitRatio menghitung persentase hit dari total (hit + miss)
func (m *CacheMetrics) HitRatio() float64 {
	total := m.hits.Load() + m.misses.Load()
	if total == 0 {
		return 0.0
	}
	return float64(m.hits.Load()) / float64(total) * 100.0
}
