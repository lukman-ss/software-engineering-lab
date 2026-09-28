package store

import (
	"sync"
	"sync/atomic"

	"bloomfilters/internal/bloom"
)

// Segment represents an immutable sorted string table (SST) component.
type Segment struct {
	data   map[string]string
	filter *bloom.Filter // nil if no filter configured
}

// NewSegment builds an immutable segment. If useFilter is true, a Bloom filter is created.
func NewSegment(kvs map[string]string, useFilter bool, epsilon float64) *Segment {
	seg := &Segment{
		data: make(map[string]string, len(kvs)),
	}
	for k, v := range kvs {
		seg.data[k] = v
	}
	if useFilter && len(kvs) > 0 {
		bf := bloom.New(uint(len(kvs)), epsilon)
		for k := range kvs {
			bf.Add([]byte(k))
		}
		seg.filter = bf
	}
	return seg
}

// LSMStore simulates an LSM-tree store with multiple segments.
type LSMStore struct {
	mu        sync.RWMutex
	segments  []*Segment
	diskReads uint64 // total segment data lookups (simulating disk I/O)
}

// NewLSMStore creates an empty LSMStore.
func NewLSMStore() *LSMStore {
	return &LSMStore{}
}

// AddSegment adds a new segment to the store.
func (s *LSMStore) AddSegment(seg *Segment) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.segments = append(s.segments, seg)
}

// Get looks up a key across all segments.
// If a segment has a Bloom filter and Check() returns false, the segment is skipped
// entirely without incrementing the diskRead counter.
func (s *LSMStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keyBytes := []byte(key)
	for i := len(s.segments) - 1; i >= 0; i-- {
		seg := s.segments[i]
		if seg.filter != nil && !seg.filter.Check(keyBytes) {
			// Bloom filter: definitely not in this segment; skip disk read!
			continue
		}
		// Candidate: simulate disk read to load data block
		atomic.AddUint64(&s.diskReads, 1)
		if v, ok := seg.data[key]; ok {
			return v, true
		}
	}
	return "", false
}

// DiskReads returns the total number of simulated disk reads performed.
func (s *LSMStore) DiskReads() uint64 {
	return atomic.LoadUint64(&s.diskReads)
}

// ResetCounters clears the disk read counter.
func (s *LSMStore) ResetCounters() {
	atomic.StoreUint64(&s.diskReads, 0)
}
