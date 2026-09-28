package candidate

import (
	"context"
	"errors"
	"sync"
	"time"

	"labs/30-leader-election/internal/coordinator"
	"labs/30-leader-election/internal/storage"
)

type State string

const (
	StateFollower State = "FOLLOWER"
	StateLeader   State = "LEADER"
	StateStopped  State = "STOPPED"
)

type Node struct {
	mu           sync.RWMutex
	id           string
	leaseKey     string
	ttl          time.Duration
	renewRate    time.Duration
	coord        *coordinator.Coordinator
	store        *storage.FencedStorage
	state        State
	currentLease *coordinator.Lease
	cancelLoop   context.CancelFunc
	pauseUntil   time.Time
}

func NewNode(id, leaseKey string, ttl, renewRate time.Duration, coord *coordinator.Coordinator, store *storage.FencedStorage) *Node {
	return &Node{
		id:        id,
		leaseKey:  leaseKey,
		ttl:       ttl,
		renewRate: renewRate,
		coord:     coord,
		store:     store,
		state:     StateFollower,
	}
}

func (n *Node) ID() string {
	return n.id
}

func (n *Node) State() State {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.state
}

func (n *Node) CurrentToken() int64 {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if n.currentLease != nil {
		return n.currentLease.FencingToken
	}
	return 0
}

func (n *Node) Start(ctx context.Context) {
	n.mu.Lock()
	if n.cancelLoop != nil {
		n.mu.Unlock()
		return
	}
	loopCtx, cancel := context.WithCancel(ctx)
	n.cancelLoop = cancel
	n.state = StateFollower
	n.mu.Unlock()

	go n.runElectionLoop(loopCtx)
}

func (n *Node) Stop() {
	n.mu.Lock()
	if n.cancelLoop != nil {
		n.cancelLoop()
		n.cancelLoop = nil
	}
	if n.state == StateLeader && n.currentLease != nil {
		_ = n.coord.Release(n.leaseKey, n.id, n.currentLease.FencingToken)
	}
	n.state = StateStopped
	n.currentLease = nil
	n.mu.Unlock()
}

func (n *Node) SimulatePause(duration time.Duration) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.pauseUntil = time.Now().Add(duration)
}

func (n *Node) PerformFencedWrite(value string) error {
	n.mu.RLock()
	isLeader := n.state == StateLeader
	lease := n.currentLease
	n.mu.RUnlock()

	if !isLeader || lease == nil {
		return errors.New("node is not active leader")
	}

	return n.store.Write(n.id, lease.FencingToken, value)
}

func (n *Node) runElectionLoop(ctx context.Context) {
	ticker := time.NewTicker(n.renewRate)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			n.mu.Lock()
			if now.Before(n.pauseUntil) {
				// Node is paused (simulating GC or unresponsiveness)
				n.mu.Unlock()
				continue
			}

			if n.state == StateLeader {
				// Renew lease
				if n.currentLease == nil {
					n.state = StateFollower
					n.mu.Unlock()
					continue
				}
				renewed, err := n.coord.Renew(n.leaseKey, n.id, n.currentLease.FencingToken, n.ttl)
				if err != nil {
					// Lost leadership
					n.state = StateFollower
					n.currentLease = nil
				} else {
					n.currentLease = renewed
				}
			} else {
				// Try to acquire
				acquired, err := n.coord.Acquire(n.leaseKey, n.id, n.ttl)
				if err == nil {
					n.state = StateLeader
					n.currentLease = acquired
				}
			}
			n.mu.Unlock()
		}
	}
}

func (n *Node) ForceStepDown() {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.state == StateLeader && n.currentLease != nil {
		_ = n.coord.Release(n.leaseKey, n.id, n.currentLease.FencingToken)
		n.state = StateFollower
		n.currentLease = nil
	}
}
