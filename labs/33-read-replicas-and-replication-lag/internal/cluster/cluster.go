package cluster

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrNotFound      = errors.New("key not found")
	ErrClusterClosed = errors.New("cluster closed")
	ErrWriteToReplica = errors.New("cannot write to read-only replica")
)

type WALEntry struct {
	LSN       uint64
	Key       string
	Value     string
	Timestamp time.Time
}

type Node struct {
	id          string
	isPrimary   bool
	mu          sync.RWMutex
	data        map[string]string
	appliedLSN  uint64
	walChannel  chan WALEntry
	lagDuration time.Duration
	cond        *sync.Cond
}

func newNode(id string, isPrimary bool, lag time.Duration) *Node {
	n := &Node{
		id:          id,
		isPrimary:   isPrimary,
		data:        make(map[string]string),
		walChannel:  make(chan WALEntry, 1024),
		lagDuration: lag,
	}
	n.cond = sync.NewCond(&n.mu)
	return n
}

func (n *Node) ID() string {
	return n.id
}

func (n *Node) IsPrimary() bool {
	return n.isPrimary
}

func (n *Node) Read(key string) (string, uint64, error) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	val, ok := n.data[key]
	if !ok {
		return "", n.appliedLSN, ErrNotFound
	}
	return val, n.appliedLSN, nil
}

func (n *Node) AppliedLSN() uint64 {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.appliedLSN
}

func (n *Node) SetLag(d time.Duration) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.lagDuration = d
}

func (n *Node) WaitForLSN(ctx context.Context, targetLSN uint64) error {
	ch := make(chan struct{})
	go func() {
		n.mu.Lock()
		defer n.mu.Unlock()
		for n.appliedLSN < targetLSN {
			n.cond.Wait()
		}
		close(ch)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
		return nil
	}
}

type ReplicationMode int

const (
	AsyncReplication ReplicationMode = iota
	SyncReplication
)

type Cluster struct {
	mu           sync.RWMutex
	primary      *Node
	replicas     []*Node
	currentLSN   atomic.Uint64
	replMode     ReplicationMode
	closed       bool
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
}

func NewCluster(numReplicas int, defaultLag time.Duration, mode ReplicationMode) *Cluster {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Cluster{
		primary:  newNode("primary", true, 0),
		replMode: mode,
		ctx:      ctx,
		cancel:   cancel,
	}

	for i := 1; i <= numReplicas; i++ {
		replica := newNode(fmt.Sprintf("replica-%d", i), false, defaultLag)
		c.replicas = append(c.replicas, replica)

		c.wg.Add(1)
		go c.replicaWorker(replica)
	}

	return c
}

func (c *Cluster) replicaWorker(replica *Node) {
	defer c.wg.Done()

	for {
		select {
		case <-c.ctx.Done():
			return
		case entry, ok := <-replica.walChannel:
			if !ok {
				return
			}

			replica.mu.RLock()
			lag := replica.lagDuration
			replica.mu.RUnlock()

			if lag > 0 {
				select {
				case <-c.ctx.Done():
					return
				case <-time.After(lag):
				}
			}

			replica.mu.Lock()
			replica.data[entry.Key] = entry.Value
			replica.appliedLSN = entry.LSN
			replica.cond.Broadcast()
			replica.mu.Unlock()
		}
	}
}

func (c *Cluster) Primary() *Node {
	return c.primary
}

func (c *Cluster) Replicas() []*Node {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]*Node, len(c.replicas))
	copy(out, c.replicas)
	return out
}

func (c *Cluster) Write(key, value string) (uint64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return 0, ErrClusterClosed
	}

	lsn := c.currentLSN.Add(1)
	entry := WALEntry{
		LSN:       lsn,
		Key:       key,
		Value:     value,
		Timestamp: time.Now(),
	}

	c.primary.mu.Lock()
	c.primary.data[key] = value
	c.primary.appliedLSN = lsn
	c.primary.mu.Unlock()

	if c.replMode == SyncReplication {
		for _, replica := range c.replicas {
			replica.mu.RLock()
			lag := replica.lagDuration
			replica.mu.RUnlock()

			if lag > 0 {
				time.Sleep(lag)
			}

			replica.mu.Lock()
			replica.data[key] = value
			replica.appliedLSN = lsn
			replica.cond.Broadcast()
			replica.mu.Unlock()
		}
	} else {
		for _, replica := range c.replicas {
			select {
			case replica.walChannel <- entry:
			default:
			}
		}
	}

	return lsn, nil
}

func (c *Cluster) CurrentLSN() uint64 {
	return c.currentLSN.Load()
}

func (c *Cluster) SetReplicationMode(mode ReplicationMode) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.replMode = mode
}

func (c *Cluster) Close() {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return
	}
	c.closed = true
	c.cancel()
	c.mu.Unlock()

	c.wg.Wait()
}
