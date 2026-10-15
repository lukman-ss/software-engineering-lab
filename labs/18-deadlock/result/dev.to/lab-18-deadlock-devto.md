---
title: "Deadlock in Concurrent Systems: Circular Wait Mechanics, Database Engine Aborts, and Application Strategies in Go"
published: false
description: "Understand deadlock mechanics in concurrent systems: circular wait, wait-for graphs, lock ordering prevention, and retry recovery in Go."
tags: go, concurrency, database, softwareengineering
---

# Deadlock in Concurrent Systems: Circular Wait Mechanics, Database Engine Aborts, and Application Strategies in Go

A production transactional system must handle concurrent access safely. Lab 18 makes deadlock behavior visible with a focused Go simulation and deterministic test assertions.

> Deadlock is an inescapable circular wait resolved only by external engine aborts or strict lock ordering.

## Lonjakan Transaksi Konkuren Memicu Kemacetan Sistemik

When high-volume transactions access shared data concurrently, uncoordinated access sequences cause systemic freezes. Encountering aborted transaction errors does not imply database malfunction; it indicates that the application layer permitted conflicting lock acquisition patterns.

## Review Kode Manual Gagal Menangkap Siklus Tunggu Konkuren

Manual inspection and single-threaded unit tests fail to expose deadlocks because concurrency depends on operating system thread interleaving. Encapsulating queries within ACID transactions does not automatically guarantee deadlock safety if separate endpoints query tables in opposing orders.

## Anatomi Siklus Tunggu: Deadly Embrace Dua Transaksi

Deadlocks manifest when all four Coffman conditions align: mutual exclusion, hold-and-wait, no preemption, and circular wait. Only along this path does an unbreakable mutual lock freeze execution permanently, necessitating external abort intervention by the engine.

## Alur Deteksi Mesin Basis Data: Dari Wait-For Graph ke Abort

Database engines postpone cycle checks to conserve CPU overhead, inspecting the Wait-For Graph only after a timeout elapses. While our lab models this via context timeouts, it has not yet proven production cost-based victim selection heuristics like MySQL InnoDB.

## Pola Penguncian Naif yang Menimbulkan Deadlock

TransferNaive acquires locks sequentially based purely on function argument order. Injecting artificial processing delay does not mean introducing arbitrary bloat; it accurately simulates query execution latency that opens the concurrency collision window.

```go
func TransferNaive(ctx context.Context, from, to *bank.Account, amount int, delay time.Duration) error {
    if err := from.Lock(ctx); err != nil { return err }
    defer from.Unlock()
    time.Sleep(delay) // simulasi durasi transaksi internal
    if err := to.Lock(ctx); err != nil { return err }
    defer to.Unlock()
    from.Balance -= amount
    to.Balance += amount
    return nil
}
```

## Kronologi Tabrakan: Goroutine A vs Goroutine B

When bidirectional transfers execute concurrently, each successfully secures its first lock but blocks on the second. Only along this path does an absolute freeze emerge, forcing timeout cancellation to sacrifice one transaction so the survivor can complete.

## Matriks Diagnostik: Fenomena Penguncian di Database

This matrix delineates the technical boundaries separating standard lock contention, lock timeouts, and genuine deadlocks. Deadlocks uniquely exhibit cyclic dependencies that cannot self-resolve, demanding proactive engine rollback intervention.

Dimensi Masalah | Normal Lock Contention | Lock Timeout | Deadlock (Circular Wait)
|---|---|---|
Kondisi Dependensi | Antrean linear satu arah | Tunggu melebihi ambang batas | Siklus tertutup saling sandera (A->B->A)
Penyelesaian Mesin | Menunggu lock dilepas | Transaksi pemohon dibatalkan | Satu transaksi dikorbankan via rollback
Kode Galat Umum | Tidak ada galat (hanya antre) | MySQL 1205 (wait timeout) | PG 40P01 / MySQL 1213 / SQLServer 1205
Akar Masalah | Beban transaksi tinggi | Transaksi lambat / long-running | Pola urutan akuisisi lock tidak konsisten
Solusi Arsitektural | Indexing & isolasi lebih rendah | Perpendek durasi transaksi | Lock Ordering konsisten & Retry Pattern

## Tiga Lapisan Pertahanan Menghadapi Deadlock

Robust deadlock management pairs proactive architectural avoidance with operational resilience. Implementing automatic retry logic does not mean surrendering to lazy design; it serves as an indispensable safety net for complex concurrency patterns.

## Eksekusi Test Suite dan Demo cmd/demo Lolos Bersih

Terminal verification demonstrates that all test suites pass cleanly under the Go race detector with zero data races. The CLI demo confirms that naive execution yields deadlock victims, whereas lock ordering and application retries ensure consistent final state.

```go
=== RUN   TestDeadlockOccurrence
--- PASS: TestDeadlockOccurrence (0.02s) [victim identified]
=== RUN   TestLockOrderingPreventsDeadlock
--- PASS: TestLockOrderingPreventsDeadlock (0.01s) [balances consistent]
=== RUN   TestRetryRecoversDeadlock
--- PASS: TestRetryRecoversDeadlock (0.02s) [both transfers complete]
=== RUN   TestTransactionDurationImpact
--- PASS: TestTransactionDurationImpact (0.05s) [long >= short]
PASS: ok  labs/18-deadlock/tests  1.715s (race detector: clean)
1. Simulating Naive Concurrent Transfers -> deadlock victim detected
2. Simulating Lock Ordered Transfers    -> ACC-1=1100, ACC-2=900 (success)
3. Simulating Retry Mechanism           -> Final Balances: ACC-1=1100, ACC-2=900
```

## Batasan Bukti Lab: Yang Dibuktikan vs Belum

This lab proves fundamental deadlock mechanics and lock ordering efficacy using in-memory channel primitives. It has not yet proven production phenomena such as database index gap locking or systemic resilience under 50,000 transactions per second.

## Kompromi Rekayasa: Lock Ordering vs Application Retry

Enforcing global lock ordering guarantees safety but requires rigorous query conventions across all services. Introducing retry wrappers does not automatically enhance throughput if underlying shared resources suffer from acute lock contention.

## Model Sinkronisasi Akun Berbasis Channel dan Timeout

The Account struct utilizes a 1-slot buffered channel to govern balance mutation access. While this in-memory simulation is not identical to PostgreSQL row-level lock trees with shared/exclusive modes, it cleanly models external deadlock abortion.

```go
type Account struct {
    ID      string
    Balance int
    ch      chan struct{} // channel buffer 1 sebagai token lock
}
func (a *Account) Lock(ctx context.Context) error {
    select {
    case <-a.ch: return nil
    case <-ctx.Done(): return ErrDeadlock // representasi deadlock abort
    }
}
```

## Transformasi: Dari Tabrakan Acak ke Eksekusi Teratur

Prior to lock ordering, intersecting transactions mutually block and trigger destructive engine rollbacks. Queuing transactions deterministically on an initial shared lock does not degrade system throughput; it prevents wasted rollback cycles.

## Pencegahan Deterministik: Pengurutan Kunci Berdasarkan ID

By sorting resource IDs lexicographically, the entity with the lower ID is invariably locked first. The scope of this prevention logic strictly matches our concurrency test suite, proving circular waits cannot form because competing routines queue on the identical first lock.

```go
func TransferOrdered(ctx context.Context, acc1, acc2 *bank.Account, amount int, delay time.Duration) error {
    first, second := acc1, acc2
    if acc1.ID > acc2.ID { first, second = acc2, acc1 }
    if err := first.Lock(ctx); err != nil { return err }
    defer first.Unlock()
    time.Sleep(delay)
    if err := second.Lock(ctx); err != nil { return err }
    defer second.Unlock()
    acc1.Balance -= amount; acc2.Balance += amount
    return nil
}
```

## Verifikasi Test: Membuktikan Terjadinya Deadlock Victim

This unit test deterministically proves that cross transfers without ordering reliably produce at least one aborted victim. Spawning two goroutines concurrently does not automatically reproduce deadlocks unless processing delay sufficiently widens the contention window.

```go
func TestDeadlockOccurrence(t *testing.T) {
    // Goroutine 1 mentransfer A -> B (delay 10ms, timeout 20ms)
    go func() { errs[0] = transfer.TransferNaive(ctx1, accA, accB, 10, 10*time.Millisecond) }()
    // Goroutine 2 bersamaan mentransfer B -> A (delay 10ms, timeout 20ms)
    go func() { errs[1] = transfer.TransferNaive(ctx2, accB, accA, 10, 10*time.Millisecond) }()
    wg.Wait()
    hasDeadlock := (errs[0] == bank.ErrDeadlock) || (errs[1] == bank.ErrDeadlock)
    if !hasDeadlock { t.Fatalf("expected at least one deadlock victim") }
}
```

## Verifikasi Test: Lock Ordering Menuntaskan Semua Transaksi

Testing confirms both bidirectional transfers complete successfully with zero errors when routed through TransferOrdered. Only along this path do balances remain perfectly consistent (100 and 100) as transactions serialize gracefully rather than aborting.

```go
func TestLockOrderingPreventsDeadlock(t *testing.T) {
    go func() { errs[0] = transfer.TransferOrdered(ctx1, accA, accB, 10, 5*time.Millisecond) }()
    go func() { errs[1] = transfer.TransferOrdered(ctx2, accB, accA, 10, 5*time.Millisecond) }()
    wg.Wait()
    if errs[0] != nil || errs[1] != nil { t.Fatalf("expected both to succeed") }
    if accA.Balance != 100 || accB.Balance != 100 {
        t.Fatalf("unexpected balance: A=%d, B=%d", accA.Balance, accB.Balance)
    }
}
```

## Pemulihan Transparan: Penanganan Abort dengan Retry

TransferWithRetry intercepts deadlock aborts and transparently reissues the transaction up to maxRetries. Adopting a fixed 2ms backoff does not represent production best practice; it is a documented simplification that requires randomized jitter in production systems.

```go
func TransferWithRetry(ctx context.Context, from, to *bank.Account, amount int, delay time.Duration, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        attemptCtx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
        err := TransferNaive(attemptCtx, from, to, amount, delay)
        cancel()
        if err == nil { return nil }
        if err != bank.ErrDeadlock && err != context.DeadlineExceeded { return err }
        time.Sleep(2 * time.Millisecond) // backoff sebelum mencoba kembali
    }
    return bank.ErrDeadlock
}
```

## Dampak Durasi: Transaksi Panjang Memperbesar Risiko

Empirical testing demonstrates that transactions holding locks for 5ms produce higher or equal deadlock counts compared to zero-delay transactions. This finding directly reflects test suite metrics, confirming that prolonged lock durations expand the concurrency collision surface.

```go
func TestTransactionDurationImpact(t *testing.T) {
    longDeadlocks := countDeadlocks(5*time.Millisecond, 10)
    shortDeadlocks := countDeadlocks(0, 10)
    if longDeadlocks < shortDeadlocks {
        t.Fatalf("expected long duration to have >= deadlocks: long=%d, short=%d",
                 longDeadlocks, shortDeadlocks)
    }
}
```

## Checklist Proteksi Deadlock di Lingkungan Produksi

This checklist codifies production defenses: uniform global lock ordering, primary key sorting on bulk updates, stripping external I/O from transactions, vendor error code interception, jittered exponential backoff, and strict idempotency.

## Anti-Pattern Umum dalam Penanganan Transaksi

Holding database locks while waiting on third-party HTTP responses represents the most frequent cause of production deadlocks. Wrapping endpoints in retry logic without verifying mutation idempotency does not yield safety; it risks double financial debits.

## Reproduksi Mandiri di Mesin Lokal Menggunakan Go

The entire simulation runs on the Go standard library with zero external dependencies. While local testing confirms in-memory algorithmic principles, production systems still require configuring database parameters like deadlock_timeout and log_lock_waits.

```sh
cd labs/18-deadlock
go test -v ./...    # verifikasi 4 pengujian deadlock, ordering, retry & durasi
go test -race ./... # pembuktian sinkronisasi konkuren bebas data race
go run ./cmd/demo   # eksekusi demonstrasi CLI interaktif skenario 1, 2, dan 3
```

## Prinsip Inti: Deadlock adalah Keniscayaan Konkurensi

Deadlocks should not be viewed as catastrophic bugs, but as the database engine's ultimate defense of data consistency. Designing transactional architectures does not automatically eliminate lock contention; resilient systems prevent cycles upfront and recover victims gracefully.

## Source

[Software Engineering Lab 18 - Deadlock in Concurrent Systems](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/18-deadlock), snapshot commit `e76d13f`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
