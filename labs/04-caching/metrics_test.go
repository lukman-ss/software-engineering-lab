package caching_test

import (
	"math"
	"testing"
	"time"

	caching "github.com/lukman-ss/software-engineering-lab/labs/04-caching"
)

func TestCacheMetricsReset(t *testing.T) {
	m := caching.NewCacheMetrics()

	// Record counters + latency
	m.IncHit()
	m.IncMiss()
	m.IncError()
	m.IncCacheOperationError()
	m.IncRebuildAttempt()
	m.IncRebuildSuccess()
	m.IncDBQuery()
	m.IncLockWait()
	m.IncDBFallback()
	m.IncEvictedKey()
	m.IncExpiredKey()
	m.IncCacheSetError()
	m.IncCacheInvalidationError()
	m.IncCacheGetOp()
	m.IncCacheSetOp()
	m.IncCacheInvalidateOp()

	m.RecordCacheGetLatency(10 * time.Millisecond)
	m.RecordCacheSetLatency(20 * time.Millisecond)
	m.RecordDBFallbackLatency(30 * time.Millisecond)
	m.RecordRebuildLatency(40 * time.Millisecond)

	// Assert non-zero
	if m.Hits() == 0 || m.Misses() == 0 || m.Errors() == 0 || m.CacheOperationErrors() == 0 ||
		m.RebuildAttempts() == 0 || m.RebuildSuccesses() == 0 ||
		m.DBQueries() == 0 || m.LockWaits() == 0 || m.DBFallbacks() == 0 ||
		m.EvictedKeys() == 0 || m.ExpiredKeys() == 0 ||
		m.CacheSetErrors() == 0 || m.CacheInvalidationErrors() == 0 ||
		m.CacheGetOps() == 0 || m.CacheSetOps() == 0 || m.CacheInvalidateOps() == 0 {
		t.Fatal("expected all counters to be non-zero before reset")
	}

	if m.AverageCacheGetLatency() == 0 || m.AverageCacheSetLatency() == 0 || m.AverageDBFallbackLatency() == 0 {
		t.Fatal("expected all average latencies to be non-zero before reset")
	}

	// Reset
	m.Reset()

	// Assert semua counters = 0
	if m.Hits() != 0 || m.Misses() != 0 || m.Errors() != 0 || m.CacheOperationErrors() != 0 ||
		m.RebuildAttempts() != 0 || m.RebuildSuccesses() != 0 ||
		m.DBQueries() != 0 || m.LockWaits() != 0 || m.DBFallbacks() != 0 ||
		m.EvictedKeys() != 0 || m.ExpiredKeys() != 0 ||
		m.CacheSetErrors() != 0 || m.CacheInvalidationErrors() != 0 ||
		m.CacheGetOps() != 0 || m.CacheSetOps() != 0 || m.CacheInvalidateOps() != 0 {
		t.Errorf("expected all counters to be 0 after reset")
	}

	if m.AverageCacheGetLatency() != 0 {
		t.Errorf("expected AverageCacheGetLatency = 0, got %v", m.AverageCacheGetLatency())
	}
	if m.AverageCacheSetLatency() != 0 {
		t.Errorf("expected AverageCacheSetLatency = 0, got %v", m.AverageCacheSetLatency())
	}
	if m.AverageDBFallbackLatency() != 0 {
		t.Errorf("expected AverageDBFallbackLatency = 0, got %v", m.AverageDBFallbackLatency())
	}
}

func TestCacheErrorRate(t *testing.T) {
	m := caching.NewCacheMetrics()

	// 1. Initial state (0 operations) -> 0.0
	if rate := m.CacheErrorRate(); rate != 0.0 {
		t.Errorf("expected 0.0 for 0 operations, got %v", rate)
	}

	// 2. 1 GET backend error out of 2 cache ops -> 50%
	m.IncCacheGetOp() // op 1: GET success
	m.IncHit()

	m.IncCacheGetOp() // op 2: GET backend error
	m.IncError()
	m.IncCacheOperationError()

	rate := m.CacheErrorRate()
	if math.Abs(rate-50.0) > 0.001 {
		t.Errorf("expected 50.0%% operation error rate, got %v", rate)
	}

	// 3. Corrupt JSON -> aggregate Errors increases, CacheOperationErrors does not increase
	errsBefore := m.Errors()
	opErrsBefore := m.CacheOperationErrors()
	m.IncError() // Corrupt JSON

	if m.Errors() != errsBefore+1 {
		t.Errorf("expected aggregate Errors to increase by 1")
	}
	if m.CacheOperationErrors() != opErrsBefore {
		t.Errorf("expected CacheOperationErrors to remain unchanged on corrupt JSON, got %d", m.CacheOperationErrors())
	}

	// 4. Marshal error -> aggregate Errors increases, operation error rate does not change
	rateBefore := m.CacheErrorRate()
	m.IncError() // Marshal error
	if rateAfter := m.CacheErrorRate(); rateAfter != rateBefore {
		t.Errorf("expected operation error rate to remain unchanged on marshal error, before %v, after %v", rateBefore, rateAfter)
	}
}
