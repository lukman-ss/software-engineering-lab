## Snippet 1 — Sticky Session Routing

Source File: `internal/router/router.go`
Purpose: Mengarahkan read ke primary selama jendela waktu `StickyDuration` sejak last-write oleh session yang sama.

```go
func (r *Router) ReadWithStickySession(sessionID, key string) (string, string, uint64, error) {
    if stateVal, ok := r.sessions.Load(sessionID); ok {
        state := stateVal.(SessionState)
        if time.Since(state.LastWriteTime) < r.config.StickyDuration {
            val, lsn, err := r.cluster.Primary().Read(key)
            return val, r.cluster.Primary().ID(), lsn, err
        }
    }
    return r.ReadLagAware(key)
}
```

Explanation: Jika session memiliki catatan tulis terakhir dan selisih waktu kurang dari `StickyDuration`, read dialihkan ke primary. Setelah TTL habis, flow jatuh ke `ReadLagAware`.

---

## Snippet 2 — Causal Token / LSN Wait

Source File: `internal/router/router.go`
Purpose: Melayani read dengan menunggu replica mencapai `minLSN` yang ditentukan, atau fallback ke primary jika timeout.

```go
func (r *Router) ReadWithToken(ctx context.Context, minLSN uint64, key string) (string, string, uint64, error) {
    replicas := r.cluster.Replicas()

    for _, replica := range replicas {
        if replica.AppliedLSN() >= minLSN {
            val, lsn, err := replica.Read(key)
            return val, replica.ID(), lsn, err
        }
    }

    if len(replicas) > 0 {
        targetReplica := replicas[0]
        waitCtx, cancel := context.WithTimeout(ctx, r.config.WaitTimeout)
        defer cancel()

        if err := targetReplica.WaitForLSN(waitCtx, minLSN); err == nil {
            val, lsn, err := targetReplica.Read(key)
            return val, targetReplica.ID(), lsn, err
        }
    }

    val, lsn, err := r.cluster.Primary().Read(key)
    return val, r.cluster.Primary().ID() + " (fallback)", lsn, err
}
```

Explanation: Iterasi replica mencari yang sudah `AppliedLSN >= minLSN`. Jika tidak ada, tunggu hingga `WaitTimeout`. Setelah itu fallback ke primary.

---

## Snippet 3 — Lag-Aware Routing dengan Fallback

Source File: `internal/router/router.go`
Purpose: Memfilter replica yang melebihi `MaxLSNDiff` dan mengalihkan ke primary jika tidak ada replica yang layak.

```go
func (r *Router) ReadLagAware(key string) (string, string, uint64, error) {
    replicas := r.cluster.Replicas()
    primaryLSN := r.cluster.CurrentLSN()

    var eligible []*cluster.Node
    for _, replica := range replicas {
        appliedLSN := replica.AppliedLSN()
        if primaryLSN >= appliedLSN {
            diff := primaryLSN - appliedLSN
            if diff <= r.config.MaxLSNDiff {
                eligible = append(eligible, replica)
            }
        }
    }

    if len(eligible) == 0 {
        val, lsn, err := r.cluster.Primary().Read(key)
        return val, r.cluster.Primary().ID() + " (fallback-lag)", lsn, err
    }

    idx := r.rrIndex.Add(1) % uint64(len(eligible))
    target := eligible[idx]
    val, lsn, err := target.Read(key)
    return val, target.ID(), lsn, err
}
```

Explanation: Hitung ΔLSN = `primaryLSN - replicaAppliedLSN`. Hanya replica dengan ΔLSN ≤ `MaxLSNDiff` yang masuk pool `eligible`. Jika kosong, read jatuh ke primary.

---

## Snippet 4 — Write dengan Mode Sync dan Async

Source File: `internal/cluster/cluster.go`
Purpose: Membedakan perilaku write pada mode async (WAL streaming) dan sync (`remote_apply`).

```go
func (c *Cluster) Write(key, value string) (uint64, error) {
    lsn := c.currentLSN.Add(1)
    entry := WALEntry{LSN: lsn, Key: key, Value: value, Timestamp: time.Now()}

    c.primary.mu.Lock()
    c.primary.data[key] = value
    c.primary.appliedLSN = lsn
    c.primary.mu.Unlock()

    if c.replMode == SyncReplication {
        for _, replica := range c.replicas {
            if lag > 0 { time.Sleep(lag) }
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
```

Explanation: Pada sync mode, write diblokir sampai replica menerapkan entry (simulasi `remote_apply`). Pada async mode, entry dikirim ke channel tanpa tunggu.

---

## Snippet 5 — WaitForLSN dengan sync.Cond dan Context Timeout

Source File: `internal/cluster/cluster.go`
Purpose: Event-driven blocking wait sampai replica `appliedLSN` mencapai target, atau return `context.DeadlineExceeded`.

```go
func (n *Node) WaitForLSN(ctx context.Context, targetLSN uint64) error {
    ch := make(chan struct{})
    stop := make(chan struct{})

    go func() {
        n.mu.Lock()
        defer n.mu.Unlock()
        for n.appliedLSN < targetLSN {
            select {
            case <-stop:
                return
            default:
            }
            n.cond.Wait()
        }
        close(ch)
    }()

    select {
    case <-ctx.Done():
        close(stop)
        n.mu.Lock()
        n.cond.Broadcast()
        n.mu.Unlock()
        return ctx.Err()
    case <-ch:
        return nil
    }
}
```

Explanation: Goroutine tunggu via `sync.Cond` sampai `appliedLSN >= targetLSN`. Jika context expired, goroutine dihentikan dan `ctx.Err()` dikembalikan.

---

## Snippet 6 — Test Stale Read Anomaly

Source File: `tests/replication_test.go`
Purpose: Membuktikan bahwa read naive ke replica lagging menghasilkan `ErrNotFound`.

```go
func TestNaiveReplicationLag_StaleRead(t *testing.T) {
    c := cluster.NewCluster(1, 500*time.Millisecond, cluster.AsyncReplication)
    defer c.Close()

    r := router.NewRouter(c, router.DefaultConfig())

    lsn, err := r.Write("session-1", "account:balance", "1000")
    if lsn != 1 {
        t.Fatalf("Expected LSN 1, got %d", lsn)
    }

    val, nodeID, _, err := r.ReadNaive("account:balance")
    if err != cluster.ErrNotFound {
        t.Fatalf("Expected ErrNotFound, got %v", err)
    }

    time.Sleep(600 * time.Millisecond)

    val, nodeID, _, err = r.ReadNaive("account:balance")
    if val != "1000" {
        t.Fatalf("Expected '1000', got '%s'", val)
    }
}
```

Explanation: Write (LSN=1) → read langsung ke replica yang belum apply → `ErrNotFound`. Setelah 600ms (melebihi lag 500ms), read berhasil.

---

## Snippet 7 — Test Concurrent Access (Race Detector)

Source File: `tests/replication_test.go`
Purpose: Memverifikasi tidak ada data race saat banyak goroutine membaca dan menulis secara bersamaan.

```go
func TestConcurrentAccess_RaceFree(t *testing.T) {
    c := cluster.NewCluster(3, 10*time.Millisecond, cluster.AsyncReplication)
    defer c.Close()

    r := router.NewRouter(c, router.DefaultConfig())

    var wg sync.WaitGroup
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(writerID int) {
            defer wg.Done()
            for j := 0; j < 20; j++ {
                token, err := r.Write(sessionID, key, val)
                if err != nil { t.Errorf("Concurrent write failed: %v", err) }
                _, _, _, _ = r.ReadWithStickySession(sessionID, key)
                _, _, _, _ = r.ReadWithToken(context.Background(), token, key)
            }
        }(i)
    }
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(readerID int) {
            defer wg.Done()
            for j := 0; j < 20; j++ {
                _, _, _, _ = r.ReadNaive(key)
                _, _, _, _ = r.ReadLagAware(key)
            }
        }(i)
    }
    wg.Wait()
}
```

Explanation: 5 writer + 10 reader × 20 operasi = 300 operasi konkuren; lolos `go test -race`.

---

## Snippet 8 — Demo: Async Replication Lag Anomaly

Source File: `cmd/demo/main.go`
Purpose: Menampilkan langsung terjadinya stale read saat read dialihkan ke replica yang belum catch-up.

```go
lsn, err := r.Write(sessionID, "user:profile:123", `{"name":"Alice","tier":"premium"}`)
fmt.Printf("[Write Primary] Wrote key 'user:profile:123', Primary LSN: %d\n", lsn)

val, nodeID, readLSN, err := r.ReadNaive("user:profile:123")
fmt.Printf("[Naive Read] Node: %s, Applied LSN: %d, Val: '%s', Err: %v\n", nodeID, readLSN, val, err)
```

Penjelasan: Tulis di primary (LSN=1), langsung baca ke replica (appliedLSN=0) → `Err: key not found`.
