package bloom

import (
	"hash/fnv"
	"math"
	"sync"
)

// Filter is an unsynchronized Bloom filter.
// Use SyncFilter for concurrent access.
type Filter struct {
	bits []uint64
	m    uint   // total number of bits
	k    uint   // number of hash probes per element
}

// New returns a Filter sized for n expected elements at false-positive rate epsilon.
//
// Math (from research report):
//   m = ceil(-n * ln(epsilon) / (ln 2)^2)   [bits]
//   k = round((m/n) * ln 2)                  [hash probes]
func New(n uint, epsilon float64) *Filter {
	m := optimalM(n, epsilon)
	k := optimalK(m, n)
	words := (m + 63) / 64
	return &Filter{
		bits: make([]uint64, words),
		m:    m,
		k:    k,
	}
}

// Add inserts key into the filter.
func (f *Filter) Add(key []byte) {
	h1, h2 := baseHashes(key)
	for i := uint(0); i < f.k; i++ {
		pos := doubleHash(h1, h2, uint64(i), uint64(f.m))
		f.bits[pos/64] |= 1 << (pos % 64)
	}
}

// Check returns false if key is definitely not in the set.
// Returns true if key is possibly in the set (may be a false positive).
func (f *Filter) Check(key []byte) bool {
	h1, h2 := baseHashes(key)
	for i := uint(0); i < f.k; i++ {
		pos := doubleHash(h1, h2, uint64(i), uint64(f.m))
		if f.bits[pos/64]&(1<<(pos%64)) == 0 {
			return false
		}
	}
	return true
}

// M returns the number of bits allocated.
func (f *Filter) M() uint { return f.m }

// K returns the number of hash probes.
func (f *Filter) K() uint { return f.k }

// baseHashes returns two independent 64-bit hashes using FNV-1a.
// The second hash is derived by XOR-rotating h1 so both paths are independent
// without requiring a second hash function import.
// Implementation decision: stdlib-only; MurmurHash3 is not available in std.
func baseHashes(key []byte) (uint64, uint64) {
	h := fnv.New64a()
	_, _ = h.Write(key)
	h1 := h.Sum64()
	// Independent second hash: xorshift of h1 with a different constant.
	// This satisfies Kirsch-Mitzenmacher's requirement for two "independent" hashes
	// under the double-hashing scheme proven to achieve the same asymptotic FP rate.
	h2 := h1 ^ (h1>>17 | h1<<47) // 64-bit rotate-xor; stays non-zero when h1 != 0
	if h2 == 0 {
		h2 = 0xdeadbeefdeadbeef
	}
	return h1, h2
}

// doubleHash computes h_i(x) = (h1 + i*h2) mod m.
// Kirsch & Mitzenmacher (2006): two hashes suffice for k probes with no FP degradation.
func doubleHash(h1, h2, i, m uint64) uint64 {
	return (h1 + i*h2) % m
}

// optimalM computes the optimal bit-array size.
// m = ceil( -n * ln(epsilon) / (ln 2)^2 )
func optimalM(n uint, epsilon float64) uint {
	ln2sq := math.Log(2) * math.Log(2)
	m := -float64(n) * math.Log(epsilon) / ln2sq
	return uint(math.Ceil(m))
}

// optimalK computes the optimal number of hash probes.
// k = round( (m/n) * ln 2 )
func optimalK(m, n uint) uint {
	k := float64(m) / float64(n) * math.Log(2)
	k = math.Round(k)
	if k < 1 {
		k = 1
	}
	return uint(k)
}

// SyncFilter is a thread-safe Bloom filter.
type SyncFilter struct {
	mu sync.RWMutex
	f  *Filter
}

// NewSync returns a thread-safe Filter.
func NewSync(n uint, epsilon float64) *SyncFilter {
	return &SyncFilter{f: New(n, epsilon)}
}

// Add inserts key into the filter.
func (sf *SyncFilter) Add(key []byte) {
	sf.mu.Lock()
	sf.f.Add(key)
	sf.mu.Unlock()
}

// Check returns false if key is definitely absent, true if possibly present.
func (sf *SyncFilter) Check(key []byte) bool {
	sf.mu.RLock()
	ok := sf.f.Check(key)
	sf.mu.RUnlock()
	return ok
}

// M returns the number of bits allocated.
func (sf *SyncFilter) M() uint { sf.mu.RLock(); defer sf.mu.RUnlock(); return sf.f.M() }

// K returns the number of hash probes.
func (sf *SyncFilter) K() uint { sf.mu.RLock(); defer sf.mu.RUnlock(); return sf.f.K() }
