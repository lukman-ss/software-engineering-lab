package sharding

import (
	"errors"
	"fmt"
	"hash/fnv"
	"sort"
	"strconv"
	"sync"
)

var ErrShardNotFound = errors.New("shard not found")

type Record struct {
	ID        string
	ShardKey  string
	Email     string
	Data      string
}

// Router interface for determining shard location.
type Router interface {
	GetShard(key string) (string, error)
	GetShards() []string
	AddShard(shardID string)
	RemoveShard(shardID string)
}

// ModuloRouter routes keys using hash(key) % N.
type ModuloRouter struct {
	mu     sync.RWMutex
	shards []string
}

func NewModuloRouter(shardIDs []string) *ModuloRouter {
	shards := make([]string, len(shardIDs))
	copy(shards, shardIDs)
	return &ModuloRouter{shards: shards}
}

func (m *ModuloRouter) GetShard(key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.shards) == 0 {
		return "", ErrShardNotFound
	}
	h := hashKey(key)
	idx := int(h % uint64(len(m.shards)))
	return m.shards[idx], nil
}

func (m *ModuloRouter) GetShards() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]string, len(m.shards))
	copy(res, m.shards)
	return res
}

func (m *ModuloRouter) AddShard(shardID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shards = append(m.shards, shardID)
}

func (m *ModuloRouter) RemoveShard(shardID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, s := range m.shards {
		if s == shardID {
			m.shards = append(m.shards[:i], m.shards[i+1:]...)
			return
		}
	}
}

type vnode struct {
	hash    uint64
	shardID string
}

// ConsistentHashRouter routes keys using consistent hashing ring with virtual nodes.
type ConsistentHashRouter struct {
	mu          sync.RWMutex
	vnodeCount  int
	ring        []vnode
	shardSet    map[string]struct{}
}

func NewConsistentHashRouter(vnodeCount int, shardIDs []string) *ConsistentHashRouter {
	if vnodeCount <= 0 {
		vnodeCount = 100
	}
	ch := &ConsistentHashRouter{
		vnodeCount: vnodeCount,
		shardSet:   make(map[string]struct{}),
	}
	for _, id := range shardIDs {
		ch.AddShard(id)
	}
	return ch
}

func (ch *ConsistentHashRouter) AddShard(shardID string) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	if _, exists := ch.shardSet[shardID]; exists {
		return
	}
	ch.shardSet[shardID] = struct{}{}

	for i := 0; i < ch.vnodeCount; i++ {
		vkey := shardID + "#" + strconv.Itoa(i)
		h := hashKey(vkey)
		ch.ring = append(ch.ring, vnode{hash: h, shardID: shardID})
	}
	sort.Slice(ch.ring, func(i, j int) bool {
		return ch.ring[i].hash < ch.ring[j].hash
	})
}

func (ch *ConsistentHashRouter) RemoveShard(shardID string) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	if _, exists := ch.shardSet[shardID]; !exists {
		return
	}
	delete(ch.shardSet, shardID)

	newRing := make([]vnode, 0, len(ch.ring))
	for _, vn := range ch.ring {
		if vn.shardID != shardID {
			newRing = append(newRing, vn)
		}
	}
	ch.ring = newRing
}

func (ch *ConsistentHashRouter) GetShard(key string) (string, error) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	if len(ch.ring) == 0 {
		return "", ErrShardNotFound
	}
	h := hashKey(key)

	idx := sort.Search(len(ch.ring), func(i int) bool {
		return ch.ring[i].hash >= h
	})
	if idx == len(ch.ring) {
		idx = 0
	}
	return ch.ring[idx].shardID, nil
}

func (ch *ConsistentHashRouter) GetShards() []string {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	shards := make([]string, 0, len(ch.shardSet))
	for s := range ch.shardSet {
		shards = append(shards, s)
	}
	sort.Strings(shards)
	return shards
}

func hashKey(key string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(key))
	return h.Sum64()
}

// Shard simulates a single database node storing records.
type Shard struct {
	ID   string
	mu   sync.RWMutex
	data map[string]Record
}

func NewShard(id string) *Shard {
	return &Shard{
		ID:   id,
		data: make(map[string]Record),
	}
}

func (s *Shard) Put(rec Record) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[rec.ID] = rec
}

func (s *Shard) Get(id string) (Record, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rec, ok := s.data[id]
	return rec, ok
}

func (s *Shard) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

func (s *Shard) AllRecords() []Record {
	s.mu.RLock()
	defer s.mu.RUnlock()
	recs := make([]Record, 0, len(s.data))
	for _, r := range s.data {
		recs = append(recs, r)
	}
	return recs
}

func (s *Shard) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, id)
}

// Global Secondary Index mapping secondary key (e.g. Email) -> ShardKey
type GlobalSecondaryIndex struct {
	mu    sync.RWMutex
	index map[string]string // email -> shardKey
}

func NewGlobalSecondaryIndex() *GlobalSecondaryIndex {
	return &GlobalSecondaryIndex{
		index: make(map[string]string),
	}
}

func (gsi *GlobalSecondaryIndex) Index(secondaryKey, shardKey string) {
	gsi.mu.Lock()
	defer gsi.mu.Unlock()
	gsi.index[secondaryKey] = shardKey
}

func (gsi *GlobalSecondaryIndex) Lookup(secondaryKey string) (string, bool) {
	gsi.mu.RLock()
	defer gsi.mu.RUnlock()
	sk, ok := gsi.index[secondaryKey]
	return sk, ok
}

// Cluster manages multiple shards with pluggable routing and GSI.
type Cluster struct {
	mu     sync.RWMutex
	shards map[string]*Shard
	router Router
	gsi    *GlobalSecondaryIndex
}

func NewCluster(router Router) *Cluster {
	c := &Cluster{
		shards: make(map[string]*Shard),
		router: router,
		gsi:    NewGlobalSecondaryIndex(),
	}
	for _, id := range router.GetShards() {
		c.shards[id] = NewShard(id)
	}
	return c
}

func (c *Cluster) AddShardNode(shardID string) {
	c.mu.Lock()
	c.shards[shardID] = NewShard(shardID)
	c.router.AddShard(shardID)
	c.mu.Unlock()
}

func (c *Cluster) Insert(rec Record) error {
	shardID, err := c.router.GetShard(rec.ShardKey)
	if err != nil {
		return err
	}

	c.mu.RLock()
	shard, ok := c.shards[shardID]
	c.mu.RUnlock()

	if !ok {
		return fmt.Errorf("%w: %s", ErrShardNotFound, shardID)
	}

	shard.Put(rec)
	if rec.Email != "" {
		c.gsi.Index(rec.Email, rec.ShardKey)
	}
	return nil
}

func (c *Cluster) GetByShardKey(shardKey, recordID string) (Record, error) {
	shardID, err := c.router.GetShard(shardKey)
	if err != nil {
		return Record{}, err
	}

	c.mu.RLock()
	shard, ok := c.shards[shardID]
	c.mu.RUnlock()

	if !ok {
		return Record{}, fmt.Errorf("%w: %s", ErrShardNotFound, shardID)
	}

	rec, found := shard.Get(recordID)
	if !found {
		return Record{}, errors.New("record not found")
	}
	return rec, nil
}

// GetByEmailUsingGSI performs point-lookup via Global Secondary Index.
func (c *Cluster) GetByEmailUsingGSI(email, recordID string) (Record, error) {
	shardKey, ok := c.gsi.Lookup(email)
	if !ok {
		return Record{}, errors.New("email not found in GSI")
	}
	return c.GetByShardKey(shardKey, recordID)
}

type ScatterGatherResult struct {
	Records          []Record
	ShardsQueried    int
	ShardResponded   int
}

// ScatterGatherBroadcast queries all shards in parallel when sharding key is unknown.
func (c *Cluster) ScatterGatherBroadcast(predicate func(Record) bool) ScatterGatherResult {
	c.mu.RLock()
	shards := make([]*Shard, 0, len(c.shards))
	for _, s := range c.shards {
		shards = append(shards, s)
	}
	c.mu.RUnlock()

	type shardResult struct {
		records []Record
		err     error
	}

	ch := make(chan shardResult, len(shards))
	var wg sync.WaitGroup

	for _, s := range shards {
		wg.Add(1)
		go func(shard *Shard) {
			defer wg.Done()
			var matched []Record
			for _, rec := range shard.AllRecords() {
				if predicate(rec) {
					matched = append(matched, rec)
				}
			}
			ch <- shardResult{records: matched}
		}(s)
	}

	wg.Wait()
	close(ch)

	var combined []Record
	responded := 0
	for res := range ch {
		responded++
		combined = append(combined, res.records...)
	}

	return ScatterGatherResult{
		Records:        combined,
		ShardsQueried:  len(shards),
		ShardResponded: responded,
	}
}

func (c *Cluster) GetShardCounts() map[string]int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	res := make(map[string]int)
	for id, s := range c.shards {
		res[id] = s.Count()
	}
	return res
}

func (c *Cluster) RebalanceData() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Gather all records
	var allRecs []Record
	for _, s := range c.shards {
		allRecs = append(allRecs, s.AllRecords()...)
	}

	// Reset shards
	for id := range c.shards {
		c.shards[id] = NewShard(id)
	}

	migrated := 0
	for _, rec := range allRecs {
		targetShard, _ := c.router.GetShard(rec.ShardKey)
		c.shards[targetShard].Put(rec)
		migrated++
	}
	return migrated
}
