package coordinator

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrLeaseHeld      = errors.New("lease is held by another candidate")
	ErrLeaseExpired   = errors.New("lease has expired or does not exist")
	ErrNotLeaseOwner  = errors.New("caller is not the lease owner")
)

type Lease struct {
	Key          string
	HolderID     string
	FencingToken int64
	ExpiresAt    time.Time
	TTL          time.Duration
}

type Coordinator struct {
	mu           sync.Mutex
	leases       map[string]*Lease
	revision     int64
	clockOffset  time.Duration
}

func New() *Coordinator {
	return &Coordinator{
		leases: make(map[string]*Lease),
	}
}

func (c *Coordinator) Acquire(key, holderID string, ttl time.Duration) (*Lease, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().Add(c.clockOffset)
	existing, found := c.leases[key]

	if found && now.Before(existing.ExpiresAt) {
		if existing.HolderID == holderID {
			existing.ExpiresAt = now.Add(ttl)
			existing.TTL = ttl
			leaseCopy := *existing
			return &leaseCopy, nil
		}
		return nil, ErrLeaseHeld
	}

	c.revision++
	lease := &Lease{
		Key:          key,
		HolderID:     holderID,
		FencingToken: c.revision,
		ExpiresAt:    now.Add(ttl),
		TTL:          ttl,
	}
	c.leases[key] = lease

	leaseCopy := *lease
	return &leaseCopy, nil
}

func (c *Coordinator) Renew(key, holderID string, token int64, ttl time.Duration) (*Lease, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().Add(c.clockOffset)
	existing, found := c.leases[key]
	if !found {
		return nil, ErrLeaseExpired
	}

	if now.After(existing.ExpiresAt) || now.Equal(existing.ExpiresAt) {
		delete(c.leases, key)
		return nil, ErrLeaseExpired
	}

	if existing.HolderID != holderID || existing.FencingToken != token {
		return nil, ErrNotLeaseOwner
	}

	existing.ExpiresAt = now.Add(ttl)
	existing.TTL = ttl
	leaseCopy := *existing
	return &leaseCopy, nil
}

func (c *Coordinator) Release(key, holderID string, token int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	existing, found := c.leases[key]
	if !found {
		return nil
	}

	if existing.HolderID == holderID && existing.FencingToken == token {
		delete(c.leases, key)
	}
	return nil
}

func (c *Coordinator) GetLease(key string) (*Lease, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	l, found := c.leases[key]
	if !found {
		return nil, false
	}
	now := time.Now().Add(c.clockOffset)
	if now.After(l.ExpiresAt) || now.Equal(l.ExpiresAt) {
		delete(c.leases, key)
		return nil, false
	}
	leaseCopy := *l
	return &leaseCopy, true
}
