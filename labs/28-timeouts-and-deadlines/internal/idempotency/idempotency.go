package idempotency

import (
	"sync"
	"time"
)

type Record struct {
	Response  string
	CreatedAt time.Time
}

type Store struct {
	mu      sync.RWMutex
	records map[string]Record
	ttl     time.Duration
}

func NewStore(ttl time.Duration) *Store {
	if ttl <= 0 {
		ttl = 1 * time.Minute
	}
	return &Store{
		records: make(map[string]Record),
		ttl:     ttl,
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rec, ok := s.records[key]
	if !ok {
		return "", false
	}
	if time.Since(rec.CreatedAt) > s.ttl {
		return "", false
	}
	return rec.Response, true
}

func (s *Store) Set(key string, response string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.records[key] = Record{
		Response:  response,
		CreatedAt: time.Now(),
	}
}
