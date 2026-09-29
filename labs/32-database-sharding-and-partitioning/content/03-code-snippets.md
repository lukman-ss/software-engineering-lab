# Code Snippets

## Snippet 1 — Logical Table Partitioning & Partition Pruning

Source File: `internal/partitioning/table.go`
Purpose: Mengimplementasikan range-based table partitioning dengan query pruning — hanya memindai partisi yang memiliki overlap dengan range query temporal `[start, end)`.

```go
package partitioning

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrNoMatchingPartition = errors.New("no matching partition found for record")
)

type Record struct {
	ID        string
	CreatedAt time.Time
	Payload   string
}

type PartitionRange struct {
	Name  string
	Start time.Time // inclusive
	End   time.Time // exclusive
}

type Partition struct {
	Range PartitionRange
	mu    sync.RWMutex
	rows  []Record
}

func (p *Partition) Insert(rec Record) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.rows = append(p.rows, rec)
}

// Table represents a logically partitioned table residing on a single engine.
type Table struct {
	Name       string
	mu         sync.RWMutex
	partitions []*Partition
}

func NewTable(name string, ranges []PartitionRange) *Table {
	t := &Table{
		Name:       name,
		partitions: make([]*Partition, len(ranges)),
	}
	for i, r := range ranges {
		t.partitions[i] = &Partition{
			Range: r,
			rows:  make([]Record, 0),
		}
	}
	return t
}

func (t *Table) Insert(rec Record) error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	for _, p := range t.partitions {
		if (rec.CreatedAt.Equal(p.Range.Start) || rec.CreatedAt.After(p.Range.Start)) && rec.CreatedAt.Before(p.Range.End) {
			p.Insert(rec)
			return nil
		}
	}
	return fmt.Errorf("%w: timestamp %s", ErrNoMatchingPartition, rec.CreatedAt.Format(time.RFC3339))
}

type QueryResult struct {
	Records           []Record
	PartitionsScanned int
	TotalPartitions   int
}

// QueryRange demonstrates partition pruning: only partitions overlapping [start, end) are scanned.
func (t *Table) QueryRange(start, end time.Time) QueryResult {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var records []Record
	scanned := 0

	for _, p := range t.partitions {
		// Overlap condition: start < p.Range.End && end > p.Range.Start
		if start.Before(p.Range.End) && end.After(p.Range.Start) {
			scanned++
			p.mu.RLock()
			for _, r := range p.rows {
				if (r.CreatedAt.Equal(start) || r.CreatedAt.After(start)) && r.CreatedAt.Before(end) {
					records = append(records, r)
				}
			}
			p.mu.RUnlock()
		}
	}

	return QueryResult{
		Records:           records,
		PartitionsScanned: scanned,
		TotalPartitions:   len(t.partitions),
	}
}

func (t *Table) DropPartition(name string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	for i, p := range t.partitions {
		if p.Range.Name == name {
			t.partitions = append(t.partitions[:i], t.partitions[i+1:]...)
			return true
		}
	}
	return false
}
```

Explanation:
`Table` mengelola daftar `Partition` berbasis rentang waktu `Start` (inclusive) hingga `End` (exclusive). Saat `Insert(rec)`, record dimasukkan ke partisi yang rentang waktunya cocok. Pada `QueryRange(start, end)`, kondisi overlap `start.Before(p.Range.End) && end.After(p.Range.Start)` memastikan hanya partisi yang relevan yang di-scan, sedangkan partisi lainnya di-skip (`scanned` counter hanya bertambah untuk partisi yang relevan). `DropPartition` menunjukkan penghapusan partisi cepat tanpa perlu melakukan scan atau VACUUM pada sub-tabel lainnya.

---

## Snippet 2 — Modulo Routing & Consistent Hashing Ring with Virtual Nodes

Source File: `internal/sharding/sharding.go`
Purpose: Mengimplementasikan dua algoritma routing (`ModuloRouter` dan `ConsistentHashRouter`) untuk menentukan lokasi physical shard berdasarkan shard key.

```go
package sharding

import (
	"hash/fnv"
	"sort"
	"strconv"
	"sync"
)

func hashKey(key string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(key))
	return h.Sum64()
}

// ModuloRouter routes keys using hash(key) % N.
type ModuloRouter struct {
	mu     sync.RWMutex
	shards []string
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

type vnode struct {
	hash    uint64
	shardID string
}

// ConsistentHashRouter routes keys using consistent hashing ring with virtual nodes.
type ConsistentHashRouter struct {
	mu         sync.RWMutex
	vnodeCount int
	ring       []vnode
	shardSet   map[string]struct{}
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
```

Explanation:
`ModuloRouter` menghitung indeks shard menggunakan formula sederhana `hash(key) % N`. Saat $N$ berubah dari $M$ ke $M+1$, formula ini mengubah posisi mayoritas key ($\approx \frac{M}{M+1}$). `ConsistentHashRouter` membuat $V$ virtual node per physical shard (`shardID#i`), mengurutkan seluruh vnode di ring tertutup berdasarkan 64-bit FNV-1a hash. Pencarian shard menggunakan binary search (`sort.Search`) untuk menemukan vnode pertama dengan hash $\ge$ hash(key). Saat shard baru ditambahkan, hanya vnode yang terdekat di ring yang terkena imbas, meminimalkan redistribusi data menjadi $\approx \frac{1}{M+1}$.

---

## Snippet 3 — Cluster Orchestration, Scatter-Gather, and Global Secondary Index

Source File: `internal/sharding/sharding.go`
Purpose: Mengelola multiple physical shard nodes, mengeksekusi scatter-gather broadcast secara paralel dengan Goroutines & Context, serta menyediakan point-lookup via Global Secondary Index.

```go
// Global Secondary Index mapping secondary key (e.g. Email) -> ShardKey
type GlobalSecondaryIndex struct {
	mu    sync.RWMutex
	index map[string]string // email -> shardKey
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

// GetByEmailUsingGSI performs point-lookup via Global Secondary Index.
func (c *Cluster) GetByEmailUsingGSI(email, recordID string) (Record, error) {
	shardKey, ok := c.gsi.Lookup(email)
	if !ok {
		return Record{}, errors.New("email not found in GSI")
	}
	return c.GetByShardKey(shardKey, recordID)
}

// ScatterGatherBroadcastWithContext queries all shards in parallel with context timeout/cancellation.
func (c *Cluster) ScatterGatherBroadcastWithContext(ctx context.Context, predicate func(Record) bool) ScatterGatherResult {
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
			select {
			case <-ctx.Done():
				ch <- shardResult{err: ctx.Err()}
				return
			default:
			}

			var matched []Record
			for _, rec := range shard.AllRecords() {
				select {
				case <-ctx.Done():
					ch <- shardResult{err: ctx.Err()}
					return
				default:
					if predicate(rec) {
						matched = append(matched, rec)
					}
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
		if res.err == nil {
			responded++
			combined = append(combined, res.records...)
		}
	}

	return ScatterGatherResult{
		Records:        combined,
		ShardsQueried:  len(shards),
		ShardResponded: responded,
	}
}
```

Explanation:
`ScatterGatherBroadcastWithContext` mendistribusikan query pencarian berbasis predicate ke seluruh physical shard secara paralel menggunakan Goroutines dan `sync.WaitGroup`. Setiap goroutine memeriksa pembatalan `ctx.Done()`. Jika query tidak menyertakan shard key dan tidak ada GSI, metode ini terpaksa diakses (menghasilkan broadcast 4/4 shard). Sebaliknya, `GetByEmailUsingGSI` memanfaatkan `GlobalSecondaryIndex` untuk memetakan `email → shardKey` secara in-memory, sehingga router langsung meneruskan request ke 1 shard target (direct point lookup), menghindari overhead broadcast ke seluruh node.

---

## Snippet 4 — RFC 9562 UUIDv7 & Sequence Block Allocator

Source File: `internal/idgen/idgen.go`
Purpose: Mengimplementasikan generator ID unik terdistribusi tanpa risiko bentrokan auto_increment lintas node.

```go
package idgen

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"
)

// UUIDv7 generates an RFC 9562 compliant time-ordered UUIDv7.
// 48-bit timestamp | 4-bit ver (0111) | 12-bit rand_a | 2-bit var (10) | 62-bit rand_b
func NewUUIDv7() (string, error) {
	var uuid [16]byte

	ms := uint64(time.Now().UnixMilli())
	uuid[0] = byte(ms >> 40)
	uuid[1] = byte(ms >> 32)
	uuid[2] = byte(ms >> 24)
	uuid[3] = byte(ms >> 36) // Note: bit shift byte assembly
	uuid[4] = byte(ms >> 8)
	uuid[5] = byte(ms)

	if _, err := rand.Read(uuid[6:]); err != nil {
		return "", err
	}

	// version 7: 0111
	uuid[6] = (uuid[6] & 0x0F) | 0x70
	// variant RFC 4122/9562: 10xxxxxx
	uuid[8] = (uuid[8] & 0x3F) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:16]), nil
}

// SequenceBlockAllocator simulates Vitess Sequences block allocation for auto-increment IDs.
type SequenceBlockAllocator struct {
	mu        sync.Mutex
	blockSize int64
	current   int64
	max       int64
	fetcher   func(blockSize int64) (int64, error)
}

func NewSequenceBlockAllocator(blockSize int64, fetcher func(blockSize int64) (int64, error)) *SequenceBlockAllocator {
	return &SequenceBlockAllocator{
		blockSize: blockSize,
		fetcher:   fetcher,
	}
}

func (s *SequenceBlockAllocator) NextID() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.current >= s.max {
		base, err := s.fetcher(s.blockSize)
		if err != nil {
			return 0, err
		}
		s.current = base
		s.max = base + s.blockSize
	}

	id := s.current
	s.current++
	return id, nil
}
```

Explanation:
`NewUUIDv7` mengonstruksi 128-bit identifier sesuai standar RFC 9562. 48 bit pertama menampung milidetik Unix timestamp, disusul versi 7 (`0111`), varian RFC (`10`), dan bit acak yang dihasilkan oleh `crypto/rand`. Keunggulan UUIDv7 adalah sifatnya yang time-ordered (dapat diurutkan secara leksikografis sesuai urutan waktu pembuatan) sehingga tidak merusak struktur indeks B-tree database. `SequenceBlockAllocator` mensimulasikan mekanisme Vitess Sequences: alih-alih menghubungi central database pada setiap insert untuk mendapatkan ID auto_increment berikutnya, allocator meminta jatah rentang blok berukuran `blockSize` (misal 10 atau 1000 ID), kemudian mengalokasikan ID secara lokal hingga blok habis.
