package compat

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrLegacyDropped   = errors.New("legacy column 'phone' dropped after contract phase")
	ErrMissingRequired = errors.New("missing required field")
)

// ponytail: in-memory mock storage; replace with database/sql for persistent store.
type MemoryStore struct {
	mu            sync.RWMutex
	users         map[int]*User
	userPhones    map[int][]PhoneEntry // keyed by user_id
	phoneSeq      int
	userSeq       int
	legacyDropped bool
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users:      make(map[int]*User),
		userPhones: make(map[int][]PhoneEntry),
	}
}

func (s *MemoryStore) CreateLegacy(name, phone string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.legacyDropped {
		return nil, ErrLegacyDropped
	}

	s.userSeq++
	id := s.userSeq
	phoneCopy := phone
	u := &User{
		ID:        id,
		Name:      name,
		Phone:     &phoneCopy,
		CreatedAt: time.Now(),
	}
	s.users[id] = u
	return u, nil
}

func (s *MemoryStore) CreateDual(name, phone string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.userSeq++
	id := s.userSeq

	var phonePtr *string
	if !s.legacyDropped {
		p := phone
		phonePtr = &p
	}

	u := &User{
		ID:        id,
		Name:      name,
		Phone:     phonePtr,
		CreatedAt: time.Now(),
	}
	s.users[id] = u

	// Write to modern user_phones table
	s.phoneSeq++
	entry := PhoneEntry{
		ID:        s.phoneSeq,
		UserID:    id,
		Number:    phone,
		IsPrimary: true,
	}
	s.userPhones[id] = append(s.userPhones[id], entry)

	return u, nil
}

func (s *MemoryStore) CreateModern(name string, phones []string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(phones) == 0 {
		return nil, ErrMissingRequired
	}

	s.userSeq++
	id := s.userSeq

	u := &User{
		ID:        id,
		Name:      name,
		Phone:     nil, // Legacy column not written
		CreatedAt: time.Now(),
	}
	s.users[id] = u

	for i, p := range phones {
		s.phoneSeq++
		entry := PhoneEntry{
			ID:        s.phoneSeq,
			UserID:    id,
			Number:    p,
			IsPrimary: i == 0,
		}
		s.userPhones[id] = append(s.userPhones[id], entry)
	}

	return u, nil
}

func (s *MemoryStore) GetUser(id int) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, exists := s.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}
	// Return a copy to avoid external mutation
	copyUser := *u
	return &copyUser, nil
}

func (s *MemoryStore) GetPhones(userID int) ([]PhoneEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exists := s.users[userID]; !exists {
		return nil, ErrUserNotFound
	}

	entries := s.userPhones[userID]
	res := make([]PhoneEntry, len(entries))
	copy(res, entries)
	return res, nil
}

func (s *MemoryStore) SavePhoneEntry(userID int, number string, isPrimary bool) (PhoneEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[userID]; !exists {
		return PhoneEntry{}, ErrUserNotFound
	}

	// Check if number already exists for user (idempotency)
	for _, p := range s.userPhones[userID] {
		if p.Number == number {
			return p, nil
		}
	}

	s.phoneSeq++
	entry := PhoneEntry{
		ID:        s.phoneSeq,
		UserID:    userID,
		Number:    number,
		IsPrimary: isPrimary,
	}
	s.userPhones[userID] = append(s.userPhones[userID], entry)
	return entry, nil
}

func (s *MemoryStore) GetUserIDs(afterID int, limit int) []int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var ids []int
	for id := range s.users {
		if id > afterID {
			ids = append(ids, id)
		}
	}
	// Sort IDs
	for i := 0; i < len(ids); i++ {
		for j := i + 1; j < len(ids); j++ {
			if ids[i] > ids[j] {
				ids[i], ids[j] = ids[j], ids[i]
			}
		}
	}

	if len(ids) > limit {
		return ids[:limit]
	}
	return ids
}

func (s *MemoryStore) TotalUsers() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.users)
}

func (s *MemoryStore) ApplyContractDropLegacyColumn() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.legacyDropped = true
	for _, u := range s.users {
		u.Phone = nil
	}
}

func (s *MemoryStore) IsLegacyDropped() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.legacyDropped
}
