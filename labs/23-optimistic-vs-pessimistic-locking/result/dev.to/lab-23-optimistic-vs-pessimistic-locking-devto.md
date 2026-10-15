---
title: "Optimistic vs Pessimistic Locking & Atomic Updates: Lost Update Anomaly Mitigation, Version Guards, and Statement Atomicity in Go"
published: false
description: "Explore Optimistic vs Pessimistic Locking & Atomic Updates in Go: mechanics, architecture, caveats, and testable invariants."
tags: go, testing, performance, softwareengineering
---

# Optimistic vs Pessimistic Locking & Atomic Updates: Lost Update Anomaly Mitigation, Version Guards, and Statement Atomicity in Go

A production system must remain resilient under load. Lab 23 makes that behavior visible with a focused Go simulation and deterministic test assertions.

> Independent read-modify-write cycles over shared state silently obliterate intermediate committed modifications.

## Anomali Lost Update: Korupsi Data Tanpa Sinyal Error

When two concurrent requests read an inventory stock of 100 and decrement by 1, both independently compute and persist 99. The initial decrement is silently overwritten without surfacing any database error. Encapsulating operations in basic SQL transactions does not automatically guarantee safety against this corruption.

## Mengapa Isolasi Default Database Gagal Mencegah Lost Update

Manual code reviews often falsely assume ACID transactions inherently prevent write collisions. The default READ COMMITTED isolation level does not automatically prevent lost updates across decoupled read-modify-write steps, as each statement re-evaluates snapshots without locking competing writers.

## Anatomi Tabrakan Waktu: Dua Transaksi Berebut Stok

This chronological trace illustrates the fatal window opened when dual threads read identical base snapshots and compute off-database. The secondary transaction overwrites the prior commit without detecting its stale view. Only along this path does a double decrement degenerate into a single lossy write.

## Tiga Taksonomi Strategi Pengendalian Konkurensi

Modern database architecture presents three distinct paradigms: serialized blocking via pessimistic locks, commit-time detection via optimistic version guards, and collapsed single-statement atomic operations. Our lab validates these control structures, but has not yet proven multi-million row physical lock tables on production RDBMS engines.

## Pola Naif Read-Modify-Write yang Membuka Celah Tabrakan

NaiveDeduct snapshots row stock outside an exclusive lock, calculates the reduction, and writes back later. Injecting a microsecond sleep does not artificially break valid code; it accurately mirrors realistic application compute delays and I/O latency across distributed networks.

```go
func (s *Store) NaiveDeduct(id int, qty int) error {
    p, err := s.Get(id) // 1. Read snapshot tanpa lock
    if err != nil { return err }
    if p.Stock < qty { return ErrInsufficientStock }
    time.Sleep(100 * time.Microsecond) // perlebar celah balapan
    s.mu.Lock()
    curr := s.products[id]
    curr.Stock = p.Stock - qty // 2. Tulis memakai snapshot usang
    s.mu.Unlock()
    atomic.AddInt64(&s.NaivelyDrawn, int64(qty))
    return nil
}
```

## Kronologi Kegagalan: 50 Goroutine dan Hilangnya 49 Mutasi

When 50 goroutines invoke NaiveDeduct simultaneously against an initial stock of 100, all workers capture identical base state and overwrite stock as 99. Only along this path does inventory discrepancy manifest: logs register 50 sold units while remaining stock sits erroneously at 99.

## Matriks Strategi: Evaluasi Tiga Pendekatan Mutasi Konkuren

This matrix outlines technical trade-offs across naive execution, pessimistic locks, optimistic guards, and atomic statements. Pessimistic locking excels for high-value contested records, whereas optimistic locking delivers superior scalability in read-heavy systems with sparse write collisions.

Dimensi Evaluasi | Naive Read-Modify-Write | Pessimistic Locking | Optimistic Locking | Atomic Update
|---|---|---|---|---|
Mekanisme Kunci | Tanpa isolasi mutasi | Row lock eksklusif | Version check saat commit | Single-statement lockless
Performa Beban Baca | Sangat cepat (tanpa lock) | Cepat (MVCC read aman) | Sangat cepat (tanpa lock) | Sangat cepat
Tingkat Konflik Tinggi | Korupsi data massal | Terblokir antrean (antre) | Lonjakan retry / error 409 | Optimal (sekuensial cepat)
Risiko Sistem | Lost update tanpa error | Deadlock / antrean panjang | Starvation pada retry loop | Terbatas pada aritmatika simpel
Skenario Ideal | Hanya sistem single-thread | Saldo bank / tiket terbatas | Edit profil / dokumen CMS | Counter hit / decrement stok

## Alur Pengambilan Keputusan Strategi Penguncian Data

Architectural selection hinges on mutation complexity. For single-row arithmetic decrements, atomic updates provide the cleanest solution. Opting for optimistic locking does not eliminate error handling responsibilities; calling services must implement robust backoff retry routines.

## Bukti Eksekusi: Verifikasi Lima Skenario dan Race Detector

Test suite and demo execution conclusively validate our theoretical concurrency assertions. The Go race detector confirms zero data races at the memory level, while logic assertions prove that unsynchronized mutations obliterate 49 inventory updates.

```go
=== RUN   TestNaiveLostUpdate
    locking_test.go:37: Lost Update Demonstrated: final stock is 99 (expected 50)
--- PASS: TestNaiveLostUpdate (0.00s)
=== RUN   TestPessimisticLocking -> PASS: (stock == 50)
=== RUN   TestOptimisticLockingConflict -> PASS: (1 succeeded, 19 conflicts rejected)
=== RUN   TestOptimisticLockingWithRetry -> PASS: (20 total updates applied, stock == 80)
=== RUN   TestAtomicConditionalUpdate -> PASS: (stock == 50)
PASS: ok  labs/23-optimistic-vs-pessimistic-locking/tests  2.510s (race detector: clean)
[1] Naive Read-Modify-Write: Stock=99 (LOST UPDATE DETECTED!)
[2] Pessimistic Locking:    Stock=50 (SUCCESS - Fully Synchronized)
[3] Optimistic Direct:      Stock=99 (SUCCESS - 1 Succeeded, 19 Conflicts Rejected)
[4] Optimistic With Retry:  Stock=80 (SUCCESS - 20 Succeeded, 61 Retries Absorbed)
[5] Atomic Update:          Stock=50 (SUCCESS - Lockless Single Statement)
```

## Batasan Bukti Lab: Yang Dibuktikan vs Belum

This laboratory employs in-memory sync.Mutex primitives to model row-level exclusive locks and version validation. It has not yet proven physical RDBMS internals such as database engine lock tables, gap locks, or multi-statement deadlock graph resolution.

## Kompromi Rekayasa: Beban Tunggu vs Beban Percobaan Ulang

Each architectural pattern introduces distinct operational overheads. Pessimistic locking incurs transaction queue wait times, whereas optimistic locking consumes CPU cycles and network rounds executing retries. Blindly applying retries does not automatically preserve throughput under relentless contention.

## Simulasi Row Lock Eksklusif Berbasis Mutex Per Baris

Store isolates synchronization by segregating the global map mutex from individual per-row mutexes. While this simulation is not identical to PostgreSQL lock managers governing shared modes and lock timeouts, it accurately replicates single-row serialization semantics.

```go
type Store struct {
    mu       sync.Mutex          // melindungi struktur map internal
    rowLocks map[int]*sync.Mutex // lock eksklusif per ID produk
    products map[int]*Product    // data penyimpanan stok dan versi
}
func (s *Store) GetRowLock(id int) *sync.Mutex {
    s.mu.Lock()
    defer s.mu.Unlock()
    lock, exists := s.rowLocks[id]
    if !exists { lock = &sync.Mutex{}; s.rowLocks[id] = lock }
    return lock
}
```

## Transformasi: Dari Overwrite Diam-diam ke State Terproteksi

Prior to applying concurrency strategies, inventory levels suffered silent corruption from blind overwrites. Following stabilization, mathematical stock invariants remain intact under heavy concurrency. Enforcing data safety does not mandate architectural sprawl when mutations can be expressed via atomic SQL statements.

## Pessimistic Locking: Mengunci Sebelum Membaca Data

PessimisticDeduct acquires the rowLock prior to inspecting stock levels and committing updates. The scope of this preventive logic strictly matches our laboratory concurrency tests, proving that all goroutines serialize predictably to decrement stock to exactly 50.

```go
func (s *Store) PessimisticDeduct(id int, qty int) error {
    rowLock := s.GetRowLock(id)
    rowLock.Lock() // akuisisi lock eksklusif baris
    defer rowLock.Unlock()
    s.mu.Lock()
    p, exists := s.products[id]
    if !exists { s.mu.Unlock(); return ErrNotFound }
    if p.Stock < qty { s.mu.Unlock(); return ErrInsufficientStock }
    p.Stock -= qty
    s.mu.Unlock()
    atomic.AddInt64(&s.Pessimistically, int64(qty))
    return nil
}
```

## Optimistic Locking: Version Guard Tanpa Memblokir Pembaca

Optimistic locking extracts record state alongside its version identifier without locking rows. At commit time, versions are evaluated; if the target row changed, the update aborts with ErrOptimisticLock. This mechanism does not automatically deliver success unless calling code manages the resulting conflict.

```go
func (s *Store) OptimisticDeduct(id int, qty int) error {
    p, err := s.Get(id) // ambil snapshot beserta nilai version
    if err != nil { return err }
    if p.Stock < qty { return ErrInsufficientStock }
    time.Sleep(50 * time.Microsecond)
    s.mu.Lock()
    defer s.mu.Unlock()
    curr, exists := s.products[id]
    if curr.Version != p.Version { // cek version guard
        atomic.AddInt64(&s.OptimisticFails, 1)
        return ErrOptimisticLock // equivalen affected_rows == 0
    }
    curr.Stock -= qty; curr.Version++
    atomic.AddInt64(&s.Optimistically, int64(qty)); return nil
}
```

## Konvergensi Optimistik: Exponential Backoff dengan Jitter

Upon encountering a version collision, the Service layer retries using exponential delays augmented with randomized jitter. Jitter staggers concurrent retries, defusing thundering-herd stampedes. Only along this path do 20 contending goroutines cleanly converge to complete all mutations.

```go
func (svc *Service) DeductOptimisticWithRetry(id, qty, maxRetries int) error {
    for attempt := 0; attempt <= maxRetries; attempt++ {
        err := svc.store.OptimisticDeduct(id, qty)
        if err == nil { return nil }
        if err != ErrOptimisticLock { return err }
        if attempt == maxRetries { return ErrOptimisticLock }
        sleepDuration := time.Duration(1<<attempt)*time.Millisecond +
            time.Duration(rand.Intn(5))*time.Millisecond
        time.Sleep(sleepDuration)
    }
    return ErrOptimisticLock
}
```

## Atomic Update: Memangkas Celah Read-Modify-Write

AtomicDeduct eliminates the temporal gap between reading and writing by binding condition checks and updates within a single critical section. This pattern does not represent a silver bullet for all workflows, as single-statement atomicity cannot enforce complex multi-table invariants.

```go
// Ekuivalen SQL: UPDATE products SET stock = stock - ? WHERE id = ? AND stock >= ?
func (s *Store) AtomicDeduct(id int, qty int) error {
    if qty <= 0 { return ErrInvalidQuantity }
    s.mu.Lock()
    defer s.mu.Unlock()
    curr, exists := s.products[id]
    if !exists { return ErrNotFound }
    if curr.Stock < qty { return ErrInsufficientStock }
    curr.Stock -= qty // validasi dan kalkulasi dalam satu seksi kritis
    atomic.AddInt64(&s.Atomically, int64(qty))
    return nil
}
```

## Verifikasi Invarian: Konsistensi Matematis di Bawah Beban

The test suite strictly verifies that optimistic rejections preserve state integrity: successful transactions plus final stock must sum exactly to 100. This empirical verification strictly adheres to the lab's test scope, proving zero data corruption despite widespread request rejection.

```go
func TestOptimisticLockingConflict(t *testing.T) {
    // 20 goroutines serempak DeductOptimisticDirect pada stok 100
    wg.Wait()
    p, _ := store.Get(3)
    // Pembuktian invarian: transaksi sukses + sisa stok wajib tepat 100
    if successCount+int64(p.Stock) != 100 {
        t.Fatalf("invarian stok rusak: sukses %d + stok %d != 100",
                 successCount, p.Stock)
    }
    if conflictCount == 0 { t.Fatalf("harus ada minimal 1 konflik") }
}
```

## Checklist Penerapan Concurrency Control di Produksi

This production checklist codifies operational best practices for concurrent database systems. Categorizing mutation lifecycles, minimizing lock durations, validating affected rows, implementing jittered retries, and running race detectors constitute essential backend defenses.

## Anti-Pattern Umum dalam Penanganan Transaksi Konkuren

Common production blunders include ignoring affected-row counts after optimistic updates and maintaining database locks across external HTTP requests. Reaching for distributed Redis locks does not automatically imply superior engineering; relational engines already provide native, low-latency locking primitives.

## Reproduksi Mandiri di Mesin Lokal Menggunakan Go

The entire test harness runs purely on the Go standard library without third-party packages. While local execution validates locking semantics in memory, production systems still require configuring database lock timeouts and monitoring contention metrics on physical RDBMS clusters.

```sh
cd labs/23-optimistic-vs-pessimistic-locking
go test -v ./...    # verifikasi 6 test skenario naive, pessimistic, optimistic & atomic
go test -race ./... # pembuktian sinkronisasi konkuren bebas data race di Go memory
go run ./cmd/demo   # eksekusi demonstrasi CLI interaktif kelima skenario berdampingan
```

## Prinsip Inti: Tidak Ada Solusi Tunggal untuk Konkurensi

Mature transactional engineering matches synchronization primitives to domain contention profiles. Mission-critical rows require pessimistic certainty, read-heavy workloads scale via optimistic guards, and simple counters thrive on atomic updates. Architecting concurrent systems does not automatically require complex infrastructure when native atomic statements suffice.

## Source

[Software Engineering Lab 23 - Optimistic vs Pessimistic Locking & Atomic Updates](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/23-optimistic-vs-pessimistic-locking), snapshot commit `aad5d78`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
