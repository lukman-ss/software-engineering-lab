package storage

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrStaleFencingToken = errors.New("stale fencing token: rejected to prevent split-brain write")
)

type Record struct {
	FencingToken int64
	Author       string
	Value        string
	Timestamp    time.Time
}

type FencedStorage struct {
	mu           sync.RWMutex
	lastSeenToken int64
	records      []Record
}

func New() *FencedStorage {
	return &FencedStorage{
		records: make([]Record, 0),
	}
}

func (s *FencedStorage) Write(author string, token int64, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if token <= s.lastSeenToken {
		return fmt.Errorf("%w: token %d <= last seen %d (author: %s)", ErrStaleFencingToken, token, s.lastSeenToken, author)
	}

	s.lastSeenToken = token
	s.records = append(s.records, Record{
		FencingToken: token,
		Author:       author,
		Value:        value,
		Timestamp:    time.Now(),
	})
	return nil
}

func (s *FencedStorage) LastSeenToken() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastSeenToken
}

func (s *FencedStorage) History() []Record {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Record, len(s.records))
	copy(out, s.records)
	return out
}
