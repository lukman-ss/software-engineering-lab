# Code Snippets

## Snippet 1 — Monotonic Fencing Token on Lease Grant

Source File: `internal/coordinator/coordinator.go`
Purpose: Menunjukkan counter revision yang selalu naik setiap lease baru diberikan.

```go
c.revision++
lease := &Lease{
	Key:          key,
	HolderID:     holderID,
	FencingToken: c.revision,
	ExpiresAt:    now.Add(ttl),
	TTL:          ttl,
}
c.leases[key] = lease
```

Explanation: Field `revision` diincrement di bawah `sync.Mutex` sebelum membuat `Lease`. Nilai ini menjadi `FencingToken` — token monotonik yang dibawa leader untuk setiap tulisan berikutnya.

## Snippet 2 — Lease Acquisition with Contention Guard

Source File: `internal/coordinator/coordinator.go`
Purpose: Menolak akuisisi bila lease masih valid milik holder lain; memperpanjang bila holder sama.

```go
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
```

Explanation: `time.Now().Add(c.clockOffset)` memakai monotonic clock Go. Pengecekan `now.Before(existing.ExpiresAt)` menjaga single-leader invariant selama lease belum kedaluwarsa.

## Snippet 3 — Lease Renewal with Expiry and Ownership Check

Source File: `internal/coordinator/coordinator.go`
Purpose: Memperbarui lease hanya bila belum kedaluwarsa dan pemanggil adalah pemegang dengan token yang cocok.

```go
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
```

Explanation: Dua gerbang keamanan: kedaluwarsa (`After` atau `Equal`) dan kepemilikan (holder + token harus cocok). Kegagalan memaksa kandidat kembali ke `FOLLOWER`.

## Snippet 4 — Fenced Storage Check-and-Set

Source File: `internal/storage/storage.go`
Purpose: Gerbang tulisan yang menolak stale token dan menerima token lebih besar secara atomik.

```go
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
```

Explanation: Perbandingan dan update terjadi atomik di bawah satu mutex. Aturan strict `token <= lastSeenToken` → ditolak. Penerimaan mengupdate `lastSeenToken` dan menambah history.

## Snippet 5 — Election Loop with GC Pause Simulation

Source File: `internal/candidate/candidate.go`
Purpose: Loop periodik yang menangani renew saat LEADER dan acquire saat FOLLOWER, dengan simulasi pause.

```go
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
				n.mu.Unlock()
				continue
			}

			if n.state == StateLeader {
				if n.currentLease == nil {
					n.state = StateFollower
					n.mu.Unlock()
					continue
				}
				renewed, err := n.coord.Renew(n.leaseKey, n.id, n.currentLease.FencingToken, n.ttl)
				if err != nil {
					n.state = StateFollower
					n.currentLease = nil
				} else {
					n.currentLease = renewed
				}
			} else {
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
```

Explanation: Setiap tick memeriksa `pauseUntil` — jika masih di-pause, tick dilewati (mensimulasikan stop-the-world GC atau partition). Leader memperbarui lease; follower mencoba memperoleh lease.

## Snippet 6 — Pause Simulation and Fenced Write

Source File: `internal/candidate/candidate.go`
Purpose: Meniru pause proses dan melakukan tulisan berfencing.

```go
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
```

Explanation: `SimulatePause` mengatur batas waktu monotonik di masa depan — elegan dan tanpa dependensi eksternal. `PerformFencedWrite` menolak bila bukan leader, lalu meneruskan token ke `FencedStorage`.

## Snippet 7 — Demo Lifecycle

Source File: `cmd/demo/main.go`
Purpose: Menjalankan siklus penuh: election → pause → failover → penolakan tulisan stale.

```go
ttl := 200 * time.Millisecond
renewInterval := 50 * time.Millisecond

nodeA := candidate.NewNode("Node-A", "app-primary", ttl, renewInterval, coord, store)
nodeB := candidate.NewNode("Node-B", "app-primary", ttl, renewInterval, coord, store)

nodeA.Start(ctx)
nodeB.Start(ctx)

time.Sleep(100 * time.Millisecond)

token1 := leader.CurrentToken()
leader.PerformFencedWrite("record-from-initial-leader")

leader.SimulatePause(350 * time.Millisecond)
time.Sleep(250 * time.Millisecond)

token2 := standby.CurrentToken()
standby.PerformFencedWrite("record-from-failover-leader")

time.Sleep(150 * time.Millisecond)
err = store.Write(leader.ID(), token1, "split-brain-stale-write")
// Diharapkan: ErrStaleFencingToken — token 1 <= last seen 2
```

Explanation: Nilai waktu (350 ms > TTL 200 ms) dipilih agar pause pasti melebihi lease. Tulisan terakhir memverifikasi pertahanan split-brain dengan token stale.

## Snippet 8 — Integration Test: Failover and Split-Brain Defense

Source File: `tests/election_test.go`
Purpose: Bukti terverifikasi dari failover otomatis dan penolakan tulisan stale.

```go
ttl := 150 * time.Millisecond
renewRate := 30 * time.Millisecond

nodeA := candidate.NewNode("node-A", "cluster-leader", ttl, renewRate, coord, store)
nodeB := candidate.NewNode("node-B", "cluster-leader", ttl, renewRate, coord, store)

nodeA.Start(ctx)
nodeB.Start(ctx)

time.Sleep(80 * time.Millisecond)

var leader1 *candidate.Node
if nodeA.State() == candidate.StateLeader {
	leader1 = nodeA
} else if nodeB.State() == candidate.StateLeader {
	leader1 = nodeB
}

err := leader1.PerformFencedWrite("leader-1-write")
token1 := leader1.CurrentToken()

leader1.SimulatePause(250 * time.Millisecond)
time.Sleep(200 * time.Millisecond)

if leader2.State() != candidate.StateLeader {
	t.Fatalf("expected candidate 2 to take over leadership after pause")
}

err = leader2.PerformFencedWrite("leader-2-write")
token2 := leader2.CurrentToken()
if token2 <= token1 {
	t.Fatalf("expected leader 2 token (%d) > leader 1 token (%d)", token2, token1)
}

err = store.Write(leader1.ID(), token1, "stale-write-after-gc-pause")
if !errors.Is(err, storage.ErrStaleFencingToken) {
	t.Fatalf("expected stale write from paused leader to be rejected by fencing storage, got %v", err)
}
```

Explanation: Test ini adalah bukti utama bahwa failover terjadi dalam TTL dan fencing menolak tulisan leader lama — dua klaim paling penting dari lab.
