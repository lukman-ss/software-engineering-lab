package db

import (
	"errors"
	"strings"
	"sync"
)

var (
	ErrNotFound = errors.New("record not found")
)

type UserRecord struct {
	ID        string
	Name      string
	FirstName string
	LastName  string
}

type UserStore struct {
	mu      sync.RWMutex
	records map[string]UserRecord
}

func NewUserStore() *UserStore {
	return &UserStore{
		records: make(map[string]UserRecord),
	}
}

func (s *UserStore) InsertLegacy(id, name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records[id] = UserRecord{
		ID:   id,
		Name: name,
	}
}

func (s *UserStore) SaveExpand(id, firstName, lastName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	combined := strings.TrimSpace(firstName + " " + lastName)
	s.records[id] = UserRecord{
		ID:        id,
		Name:      combined,
		FirstName: firstName,
		LastName:  lastName,
	}
}

func (s *UserStore) GetUser(id string) (UserRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rec, ok := s.records[id]
	if !ok {
		return UserRecord{}, ErrNotFound
	}

	if rec.FirstName == "" && rec.LastName == "" && rec.Name != "" {
		parts := strings.SplitN(rec.Name, " ", 2)
		rec.FirstName = parts[0]
		if len(parts) > 1 {
			rec.LastName = parts[1]
		}
	} else if rec.Name == "" {
		rec.Name = strings.TrimSpace(rec.FirstName + " " + rec.LastName)
	}

	return rec, nil
}
