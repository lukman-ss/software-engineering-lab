package cache

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var ErrNotFound = errors.New("record not found")

type MockDB struct {
	mu         sync.RWMutex
	data       map[string]string
	queryDelay time.Duration
	queryCount atomic.Int64
	writeCount atomic.Int64
}

func NewMockDB(delay time.Duration) *MockDB {
	return &MockDB{
		data:       make(map[string]string),
		queryDelay: delay,
	}
}

func (db *MockDB) SetData(key, value string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = value
}

func (db *MockDB) Query(ctx context.Context, key string) (string, error) {
	db.queryCount.Add(1)
	if db.queryDelay > 0 {
		select {
		case <-time.After(db.queryDelay):
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	val, ok := db.data[key]
	if !ok {
		return "", ErrNotFound
	}
	return val, nil
}

func (db *MockDB) Write(ctx context.Context, key, value string) error {
	db.writeCount.Add(1)
	if db.queryDelay > 0 {
		select {
		case <-time.After(db.queryDelay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = value
	return nil
}

func (db *MockDB) QueryCount() int64 {
	return db.queryCount.Load()
}

func (db *MockDB) WriteCount() int64 {
	return db.writeCount.Load()
}

func (db *MockDB) ResetCounts() {
	db.queryCount.Store(0)
	db.writeCount.Store(0)
}
