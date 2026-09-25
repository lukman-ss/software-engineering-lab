package ratelimit

import (
	"sync"
	"time"
)

// PerClientLimiter implements rate limiting per logical client.
type PerClientLimiter struct {
	mu             sync.RWMutex
	buckets        map[string]*TokenBucket
	globalLimit    float64 // Global rate limit across all clients
	globalCapacity float64
	globalMu       sync.Mutex
	globalCount    float64 // Tokens used recently
}

// PerClientConfig holds configuration for per-client rate limiting.
type PerClientConfig struct {
	PerClientRate     float64 // Tokens/sec per client
	PerClientCapacity float64 // Burst capacity per client
	GlobalRate        float64 // Global rate limit (0 = unlimited)
	GlobalCapacity    float64 // Global burst capacity (0 = unlimited)
}

// NewPerClientLimiter creates a new per-client rate limiter.
func NewPerClientLimiter(cfg PerClientConfig) *PerClientLimiter {
	return &PerClientLimiter{
		buckets:        make(map[string]*TokenBucket),
		globalLimit:    cfg.GlobalRate,
		globalCapacity: cfg.GlobalCapacity,
		globalCount:    0,
	}
}

// Limiter returns a new limiter for the specified client identifier.
func (pc *PerClientLimiter) Bucket(clientID string, cfg PerClientConfig) *TokenBucket {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	bucket, exists := pc.buckets[clientID]
	if !exists {
		bucket = NewTokenBucket(TokenBucketConfig{
			Capacity: cfg.PerClientCapacity,
			Rate:     cfg.PerClientRate,
		})
		pc.buckets[clientID] = bucket
	}
	return bucket
}

// Use attempts to consume tokens from both per-client and global limits.
func (pc *PerClientLimiter) Use(clientID string, n float64, cfg PerClientConfig) error {
	// Check per-client limit
	bucket := pc.Bucket(clientID, cfg)

	// Refill bucket based on its own rate
	allowed, err := bucket.Use(n)
	if !allowed {
		return ErrRateLimitExceeded
	}

	// Check global limit
	if pc.globalLimit > 0 {
		pc.globalMu.Lock()
		// Simplified: just check global count
		if pc.globalCount >= pc.globalCapacity {
			pc.globalMu.Unlock()
			return ErrRateLimitExceeded
		}
		pc.globalCount += n
		pc.globalMu.Unlock()
	}

	_ = err // bucket.Use doesn't error on success
	return nil
}

// TokenBucket returns the bucket for a client.
func (pc *PerClientLimiter) TokenBucket(clientID string, cfg PerClientConfig) *TokenBucket {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	bucket, exists := pc.buckets[clientID]
	if !exists {
		bucket = NewTokenBucket(TokenBucketConfig{
			Capacity: cfg.PerClientCapacity,
			Rate:     cfg.PerClientRate,
		})
		pc.buckets[clientID] = bucket
	}
	return bucket
}

// ClientCount returns the number of tracked clients.
func (pc *PerClientLimiter) ClientCount() int {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	return len(pc.buckets)
}

// Cleanup removes stale client entries.
// Call periodically to prevent unbounded map growth.
func (pc *PerClientLimiter) Cleanup() {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	// In production, check lastAccess time and remove old entries
	// This is a placeholder for the concept
	for clientID := range pc.buckets {
		delete(pc.buckets, clientID) // This would be selective in real impl
	}
}

/*
RISIKO DAN KOMPONEN KEAMANAN (Prompt 053)
========================================

1. Unbounded Map (Memory Growth)
   --------------------------------
   Masalah: Setiap clientID baru membuat entry baru di map.
   Bisa dipergunakan untuk memory exhaustion attack.

   Solusi:
   - Implement LRU eviction atau TTL-based cleanup
   - Set max_clients limit
   - Rate limit pada level gateway sebelum masuk rate limiter

   Contoh implementasi:
*/

type SafePerClientLimiter struct {
	mu              sync.RWMutex
	buckets         map[string]*TokenBucket
	maxClients      int
	cleanupInterval time.Duration
	lastCleanup     time.Time
}

/*
   2. IP Spoofing
   --------------
   Masalah: Client bisa mengubah IP/address header untuk bypass limit.

   Solusi:
   - Gunakan endpoint IP sebagai fallback (remote_addr)
   - Validasi dengan IP reputation service
   - Implement ASYNC rate limiting (bukan blocking) bila suspicion

   3. Proxy Headers & Load Balancer
   -------------------------------
   Masalah: X-Forwarded-For dan X-Real-IP bisa dimanipulasi.
   Beberapa client melalui proxy NAT memiliki IP yang sama.

   Solusi:
   - Konfigurasi trusted proxies di load balancer
   - Gunakan combined key: (proxy_user OR session_id) + IP
   - Rate limit di API gateway level, bukan upstream

   4. NAT Sharing
   --------------
   Masalah: Banyak klien dibalun NAT memiliki IP yang sama.
   Satu user menerima seluruh batas rate limit bersama.

   Solusi:
   - Gunakan session/cookie/token sebagai client identifier
   - Rate limit per endpoint untuk critical operations
   - Adaptive rate limiting based on user tier

   Contoh implementasi hybrid:
*/

type HybridClientIdentifier struct {
	// Priority:
	// 1. Session/API token (most accurate)
	// 2. User ID (from auth token)
	// 3. Client IP (from remote_addr)
	// 4. Fingerprint (User-Agent + Accept-Language)

	SessionID string
	UserID    string
	RemoteIP  string
	UserAgent string
}

// SafePerClientLimiterWithSts implements safer client identification.
type SafePerClientLimiterWithSts struct {
	mu         sync.RWMutex
	buckets    map[string]*TokenBucket
	maxClients int
	sts        *time.Ticker
	stopCh     chan struct{}
}

// NewSafePerClientLimiterWithSts creates a limiter with safety protections.
func NewSafePerClientLimiterWithSts(maxClients int) *SafePerClientLimiterWithSts {
	return &SafePerClientLimiterWithSts{
		buckets:    make(map[string]*TokenBucket),
		maxClients: maxClients,
		sts:        time.NewTicker(time.Minute),
		stopCh:     make(chan struct{}),
	}
}

// GetOrCreate returns or creates bucket for client, with eviction.
func (l *SafePerClientLimiterWithSts) GetOrCreate(clientID string, cfg TokenBucketConfig) (*TokenBucket, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.buckets[clientID]; !exists {
		if len(l.buckets) >= l.maxClients {
			// Eviction strategy: remove oldest (simple FIFO, could be LRU)
			for key := range l.buckets {
				delete(l.buckets, key)
				break
			}
		}
		l.buckets[clientID] = NewTokenBucket(cfg)
		return l.buckets[clientID], true
	}
	return l.buckets[clientID], false
}
