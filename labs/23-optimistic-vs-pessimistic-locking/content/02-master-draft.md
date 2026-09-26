# Optimistic vs Pessimistic Locking — Mencegah Lost Update pada Aplikasi Concurrent

## Problem

Dua transaksi konkurensi membaca nilai yang sama, memodifikasinya secara independen, dan menulis kembali secara berurutan menyebabkan penimpaan data *diam-diam* (lost update) tanpa error. Ini bukan bug—ini adalah konsekuensi logis dari isolasi data tanpa mekanisme concurrency control yang tepat.

Lost update bukanlah dirty read, non-repeatable read, atau phantom read. Ini adalah penimpaan *committed* write oleh write kedua tanpa deteksi. Karena tidak ada error, sistem tetap berjalan normal, tetapi data sudah corrupt.

```text
Initial stock: 100

Transaction A              Transaction B
---------                  ---------
READ stock = 100
                           READ stock = 100
WRITE stock = 99
                           WRITE stock = 99

Final stock: 99 (expected 98)
```

Dua *deduct* berhasil dieksekusi, tetapi stok hanya berkurang satu unit. Tidak ada error. Ini *silent data corruption*.

## Why This Matters

Lost update adalah masalah nyata yang terjadi di database produksi dengan isolation level default:

- PostgreSQL default: **READ COMMITTED**
- Oracle default: **READ COMMITTED**
- MySQL/InnoDB default: **REPEATABLE READ** (mencegah lost update hanya jika menggunakan SERIALIZABLE atau explicit locking)

Di READ COMMITTED, bahkan pembungkusan dalam transaksi tidak mencegah lost update jika menggunakan pola `read-modify-write` naive tanpa `SELECT ... FOR UPDATE` atau version guard.

Akibatnya:
- Stok produk tidak akurat
- Saldo rekening salah
- Reservasi seat double-booked
- Queue number duplikat

Ketika ini terjadi di production, biasanya terdeteksi *by accident* melalui discrepancy report, bukan error log—karena tidak ada yang *error*.

## Mental Model

Bayangkan sebuah *shared whiteboard* di ruang rapat. Beberapa orang ingin mengupdate data pada saat yang sama.

### Pessimistic Locking (SELECT ... FOR UPDATE)
Anda mengambil spidol, menulis "DILARANG MENGHAPUS / MENGUPDATE" di pojok kiri whiteboard, lalu mengerjakan update. Orang lain yang ingin update harus tunggu sampai Anda selesai dan hapus catatan tersebut. Ini *blocking*—anda pasti punya eksklusivitas, tetapi orang lain menunggu.

### Optimistic Locking (version guard + affected_rows check)
Anda merekam versi whiteboard saat Anda mulai ("Versi 5"). Anda mengerjakan update. Saat mau menulis, Anda cek: "Masih Versi 5?" Jika ya, tulis dan naikkan ke Versi 6. Jika tidak (ada yang update ke Versi 7), Anda batalkan—data sudah lama berubah. Ini *non-blocking* untuk readers, tapi Anda mungkin harus ulang (retry).

### Atomic Single-Statement
Anda menulis satu perintah: "Jika stok >= 1, kurangi 1." Database mengerjakan ini sebagai satu *atomic operation*. Tidak ada read-modify-write window—ini dianggap satu statement yang tidak bisa terpotong.

## Core Concept

### Lost Update Anomaly

Lost update terjadi ketika dua transaksi:
1. Membaca nilai yang sama
2. Memodifikasi secara independen
3. Menulis kembali secara berurutan

Transaksi pertama yang commit ditimpa oleh transaksi kedua. Tidak ada error—data corrupt secara diam-diam.

Ini mungkin terjadi di READ COMMITTED (PostgreSQL, Oracle) dan REPEATABLE READ (MySQL) *tanpa* explicit locking atau version guard.

**Key insight:** Transaction alone tidak cukup. Query design menentukan correctness.

### Pessimistic Locking

`SELECT ... FOR UPDATE` mengambil row-level exclusive lock pada baris yang dibaca. Lock ini:
- Bertahan sampai commit/rollback
- Memblokir write lain (UPDATE/DELETE/SELECT FOR UPDATE) pada baris yang sama
- Tidak memblokir plain SELECT (non-locking read)

```sql
BEGIN;

SELECT stock FROM products WHERE id = 10 FOR UPDATE;

-- hitung new_stock = stock - 1 di application

UPDATE products SET stock = new_stock WHERE id = 10;

COMMIT;
```

Karena lock diambil di awal, transaksi kedua yang ingin update baris yang sama harus *block* (tunggu) sampai transaksi pertama commit/rollback.

**Trade-offs:**
- Pros: Tidak ada conflict—data selalu konsisten
- Cons: Concurrency loss, wait time, deadlock risk

PostgreSQL dokumentasi: "A transaction seeking a lock will wait indefinitely ... it is a bad idea to hold transactions open for long periods (e.g., while waiting for user input)."

Jangan pegang lock saat network call (payment gateway, PDF generation). Ini anti-pattern.

### Optimistic Locking

Optimistic locking *mengasumsikan* konflik jarang terjadi. Itu tidak blokir reader/writer lain. Konflik dideteksi saat commit dengan memasukkan version/timestamp di WHERE clause:

```sql
UPDATE products
SET stock = stock - 1,
    version = version + 1
WHERE id = 10 AND version = 5;
```

Jika `affected_rows == 1`, update berhasil. Jika `affected_rows == 0`, berarti ada transaksi lain yang update baris ini (version != 5)—konflik terdeteksi.

Application harus handle ini: reload data, recalculate, retry (dengan backoff), atau return error 409 Conflict.

**Trade-offs:**
- Pros: High concurrency, no blocking
- Cons: Retry overhead, must handle conflict explicitly

### Atomic Single-Statement

Untuk simple counter/quota/stock decrement, satu conditional UPDATE bisa dianggap atomic di level statement:

```sql
UPDATE products
SET stock = stock - 3
WHERE id = 10 AND stock >= 3;
```

Jika `affected_rows == 1`, decrement berhasil. Jika `affected_rows == 0`, stock tidak cukup.

Statement atomic di semua major database: jika statement gagal di tengah, efeknya rollback. Tidak ada intermediate state yang terlihat oleh transaction lain.

**Trade-offs:**
- Pros: Lockless, single statement, no version overhead
- Cons: Hanya cocok untuk simple arithmetic; tidak bisa complex business logic

## Failure Scenario

Bayangkan dua customer ingin beli product terakhir (stok = 1).

```
Initial stock: 1

Customer A                     Customer B
------------                   ------------
READ stock = 1
                               READ stock = 1
WRITE stock = 0
                               WRITE stock = 0

Final stock: 0 (expected -1 atau error)
```

Stok mencapai -1 atau 0, padahal dua customer merasa sudah berhasil. Ini *silent inventory corruption*.

Atau bayangkan 50 goroutines yang concurrent mengurangi stok 100:

```
Naive:
- 50 concurrent deduct calls
- Expected final stock: 50
- Actual final stock: 99 (LOST UPDATE DETECTED!)

Pessimistic:
- 50 concurrent deduct calls
- Actual final stock: 50 (SUCCESS)

Optimistic (no retry):
- 20 concurrent deduct calls
- Successful: 1, Conflicts rejected: 19
- Actual final stock: 99 (STATE GUARDED, ZERO CORRUPTION)

Atomic:
- 50 concurrent deduct calls
- Actual final stock: 50 (SUCCESS)
```

## How It Works

### Naive Read-Modify-Write (Failure Path)

```go
func (s *Store) NaiveDeduct(id int, qty int) error {
    p, err := s.Get(id) // READ
    if err != nil {
        return err
    }
    if p.Stock < qty {
        return ErrInsufficientStock
    }
    time.Sleep(100 * time.Microsecond) // artificial delay simulating calc window
    
    s.mu.Lock()
    curr := s.products[id]
    curr.Stock = p.Stock - qty // STALE WRITE: p.Stock masih old value
    s.mu.Unlock()
    
    atomic.AddInt64(&s.NaivelyDrawn, int64(qty))
    return nil
}
```

`p.Stock` dibaca di awal, tetapi saat ditulis kembali, nilai tersebut sudah *stale*. Jika transaksi lain update `products[id]` di antara `Get()` dan write, perubahan mereka ditimpa.

### Pessimistic Locking (Row-Level Exclusive)

```go
func (s *Store) PessimisticDeduct(id int, qty int) error {
    rowLock := s.GetRowLock(id)
    rowLock.Lock()
    defer rowLock.Unlock()
    
    s.mu.Lock()
    p, exists := s.products[id]
    if !exists {
        s.mu.Unlock()
        return ErrNotFound
    }
    if p.Stock < qty {
        s.mu.Unlock()
        return ErrInsufficientStock
    }
    p.Stock -= qty
    s.mu.Unlock()
    
    atomic.AddInt64(&s.Pessimistically, int64(qty))
    return nil
}
```

Simulasi row-level lock: per-row `sync.Mutex` diambil sebelum read dan held sampai write selesai. Ini memastikan hanya satu goroutine yang bisa execute critical section untuk row tertentu pada satu waktu.

### Optimistic Locking (Version Guard)

```go
func (s *Store) OptimisticDeduct(id int, qty int) error {
    p, err := s.Get(id) // READ tanpa lock
    if err != nil {
        return err
    }
    if p.Stock < qty {
        return ErrInsufficientStock
    }
    time.Sleep(50 * time.Microsecond) // artificial delay
    
    s.mu.Lock()
    defer s.mu.Unlock()
    curr, exists := s.products[id]
    if !exists {
        return ErrNotFound
    }
    
    if curr.Version != p.Version { // VERSION CHECK
        atomic.AddInt64(&s.OptimisticFails, 1)
        return ErrOptimisticLock // CONFLICT DETECTED
    }
    
    curr.Stock -= qty
    curr.Version++ // VERSION INCREMENT
    atomic.AddInt64(&s.Optimistically, int64(qty))
    return nil
}
```

`p.Version` dibaca di awal. Saat write, cek apakah `curr.Version == p.Version`. Jika tidak (ada perubahan di antara read dan write), return `ErrOptimisticLock`.

### Optimistic Locking With Retry (Convergence)

```go
func (svc *Service) DeductOptimisticWithRetry(id int, qty int, maxRetries int) error {
    for attempt := 0; attempt <= maxRetries; attempt++ {
        err := svc.store.OptimisticDeduct(id, qty)
        if err == nil {
            return nil
        }
        if err != ErrOptimisticLock {
            return err
        }
        if attempt == maxRetries {
            return ErrOptimisticLock
        }
        // Jittered exponential backoff
        sleepDuration := time.Duration(1<<attempt)*time.Millisecond + time.Duration(rand.Intn(5))*time.Millisecond
        time.Sleep(sleepDuration)
    }
    return ErrOptimisticLock
}
```

Pada konflik (`ErrOptimisticLock`), retry dengan jittered exponential backoff (`1ms, 2ms, 4ms, ... + jitter 0..5ms`). Ini mengurangi thundering herd saat high contention.

### Atomic Single-Statement (Lockless)

```go
func (s *Store) AtomicDeduct(id int, qty int) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    curr, exists := s.products[id]
    if !exists {
        return ErrNotFound
    }
    if curr.Stock < qty {
        return ErrInsufficientStock
    }
    
    curr.Stock -= qty // single statement: check + decrement dalam satu lock hold
    atomic.AddInt64(&s.Atomically, int64(qty))
    return nil
}
```

Simulasi statement-level atomicity: read-check-write dalam satu hold lock, tanpa external sleep atau read-modify-write window.

## Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                            cmd/demo/main.go                          │
│  (CLI runner: 5 scenarios side-by-side)                             │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        internal/inventory/                           │
│  ┌──────────────────┐         ┌──────────────────┐                  │
│  │     model.go     │         │     service.go   │                  │
│  │ - Product struct │◄───────►│ - DeductNaive    │                  │
│  │ - Errors         │         │ - DeductPess.    │                  │
│  └──────────────────┘         │ - DeductOpt.     │                  │
│                               │ - DeductOpt.retry│                  │
│                               │ - DeductAtomic   │                  │
│                               └──────────────────┘                  │
│                                          │                          │
│                                          ▼                          │
│  ┌──────────────────┐         ┌──────────────────┐                  │
│  │     store.go     │◄────────┤     tests/       │                  │
│  │ - NewStore()     │         │ - TestNaive...   │                  │
│  │ - Seed()         │         │ - TestPess...    │                  │
│  │ - Get()          │         │ - TestOpt...     │                  │
│  │ - NaiveDeduct()  │         │ - TestOpt.retry  │                  │
│  │ - Pessimistic... │         │ - TestAtomic...  │                  │
│  │ - Optimistic...  │         └──────────────────┘                  │
│  │ - AtomicDeduct() │                                               │
│  │ - Row locks      │                                               │
│  │ - Version counter│                                               │
│  └──────────────────┘                                               │
└─────────────────────────────────────────────────────────────────────┘
```

- **model.go**: Domain model (`Product` struct) dan error definitions (`ErrNotFound`, `ErrInsufficientStock`, `ErrOptimisticLock`, `ErrInvalidQuantity`)
- **store.go**: Simulated storage engine dengan in-memory `map[int]*Product`, row-level mutexes (`map[int]*sync.Mutex`), atomic counters (`NaivelyDrawn`, `Pessimistically`, `Optimistically`, `OptimisticFails`, `Atomically`)
- **service.go**: Business layer yang wrap `store.go` methods, menambahkan retry logic di optimistic path
- **locking_test.go**: Automated concurrency tests dengan goroutines 20-50, invariant assertions
- **cmd/demo/main.go**: CLI runner yang menampilkan 5 skenario side-by-side

## Implementation

### Data Model (`model.go`)

```go
type Product struct {
    ID      int
    Name    string
    Stock   int
    Version int
}
```

- `ID`: primary key
- `Stock`: quantity available (int, bukan pointer untuk avoid nil ambiguity)
- `Version`: optimistic lock guard (increment setiap successful update)

Errors:
- `ErrNotFound`: product tidak ditemukan
- `ErrInsufficientStock`: stock < requested qty
- `ErrOptimisticLock`: version mismatch (konflik terdeteksi)
- `ErrInvalidQuantity`: qty <= 0

### Store (`store.go`)

```go
type Store struct {
    mu       sync.Mutex
    rowLocks map[int]*sync.Mutex
    products map[int]*Product
    
    NaivelyDrawn    int64
    Pessimistically int64
    Optimistically  int64
    OptimisticFails int64
    Atomically      int64
}
```

- `mu`: global lock untuk `products` map access (tidak untuk per-row—itu `rowLocks`)
- `rowLocks`: per-row mutex map untuk pessimistic locking simulasi
- `products`: in-memory store
- Atomic counters untuk tracking statistics

### Service (`service.go`)

```go
func (svc *Service) DeductOptimisticWithRetry(id int, qty int, maxRetries int) error {
    for attempt := 0; attempt <= maxRetries; attempt++ {
        err := svc.store.OptimisticDeduct(id, qty)
        if err == nil {
            return nil
        }
        if err != ErrOptimisticLock {
            return err
        }
        if attempt == maxRetries {
            return ErrOptimisticLock
        }
        sleepDuration := time.Duration(1<<attempt)*time.Millisecond + time.Duration(rand.Intn(5))*time.Millisecond
        time.Sleep(sleepDuration)
    }
    return ErrOptimisticLock
}
```

Jittered exponential backoff:
- `1ms` on attempt 0
- `2ms + jitter` on attempt 1
- `4ms + jitter` on attempt 2
- `8ms + jitter` on attempt 3, dst

Jitter (`rand.Intn(5)`) mengurangi synchronized retry pattern yang bisa memicu thundering herd.

## Code Walkthrough

### Test: Lost Update Demonstrated

```go
func TestNaiveLostUpdate(t *testing.T) {
    store := inventory.NewStore()
    store.Seed(1, "Item A", 100)
    svc := inventory.NewService(store)
    
    const goroutines = 50
    var wg sync.WaitGroup
    wg.Add(goroutines)
    
    for i := 0; i < goroutines; i++ {
        go func() {
            defer wg.Done()
            _ = svc.DeductNaive(1, 1)
        }()
    }
    wg.Wait()
    
    p, err := store.Get(1)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if p.Stock == 50 {
        t.Errorf("expected lost update anomaly to manifest (stock != 50), got stock = %d", p.Stock)
    }
    t.Logf("Lost Update Demonstrated: 50 deduct calls occurred, but final stock is %d (expected 50 under proper locking)", p.Stock)
}
```

- 50 goroutines concurrent memanggil `DeductNaive(1, 1)`
- Expected invariant: `InitialStock - Deductions = FinalStock` → `100 - 50 = 50`
- Actual: stock 99, bukan 50 (lost update)
- Test *pass* jika stock != 50 (lost update terbukti)

### Test: Pessimistic Locking Exact Invariant

```go
func TestPessimisticLocking(t *testing.T) {
    store := inventory.NewStore()
    store.Seed(2, "Item B", 100)
    svc := inventory.NewService(store)
    
    const goroutines = 50
    var wg sync.WaitGroup
    wg.Add(goroutines)
    
    for i := 0; i < goroutines; i++ {
        go func() {
            defer wg.Done()
            err := svc.DeductPessimistic(2, 1)
            if err != nil {
                t.Errorf("deduct failed: %v", err)
            }
        }()
    }
    wg.Wait()
    
    p, err := store.Get(2)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if p.Stock != 50 {
        t.Fatalf("expected stock 50, got %d", p.Stock)
    }
}
```

- 50 goroutines concurrent memanggil `DeductPessimistic(2, 1)`
- Expected: stock = 50 (exact invariant)
- Actual: stock = 50 (success—fully synchronized)

### Test: Optimistic Locking Conflict Detection

```go
func TestOptimisticLockingConflict(t *testing.T) {
    store := inventory.NewStore()
    store.Seed(3, "Item C", 100)
    svc := inventory.NewService(store)
    
    const goroutines = 20
    var wg sync.WaitGroup
    wg.Add(goroutines)
    
    var conflictCount int64
    var successCount int64
    var mu sync.Mutex
    
    for i := 0; i < goroutines; i++ {
        go func() {
            defer wg.Done()
            err := svc.DeductOptimisticDirect(3, 1)
            mu.Lock()
            if err == inventory.ErrOptimisticLock {
                conflictCount++
            } else if err == nil {
                successCount++
            }
            mu.Unlock()
        }()
    }
    wg.Wait()
    
    p, _ := store.Get(3)
    if successCount+int64(p.Stock) != 100 {
        t.Fatalf("stock invariant violated: success %d + stock %d != 100", successCount, p.Stock)
    }
    if conflictCount == 0 {
        t.Fatalf("expected at least one optimistic lock conflict, got 0")
    }
    t.Logf("Optimistic Locking Conflict Demonstrated: %d succeeded, %d conflicts rejected", successCount, conflictCount)
}
```

- Invariant: `successCount + stock == initialStock` (semua successful deduction + remaining stock = 100)
- Conflicts > 0 harus terjadi (setidaknya 1)
- Actual: 1 success, 19 conflicts (1 + 99 = 100)

### Test: Optimistic Locking With Retry Convergence

```go
func TestOptimisticLockingWithRetry(t *testing.T) {
    store := inventory.NewStore()
    store.Seed(4, "Item D", 100)
    svc := inventory.NewService(store)
    
    const goroutines = 20
    var wg sync.WaitGroup
    wg.Add(goroutines)
    
    for i := 0; i < goroutines; i++ {
        go func() {
            defer wg.Done()
            _ = svc.DeductOptimisticWithRetry(4, 1, 10)
        }()
    }
    wg.Wait()
    
    p, _ := store.Get(4)
    if p.Stock != 100-int(store.Optimistically) {
        t.Fatalf("inconsistent stock: expected %d, got %d", 100-int(store.Optimistically), p.Stock)
    }
    t.Logf("Optimistic Retry Successful: %d total updates applied, retried conflicts resolved", store.Optimistically)
}
```

- `store.Optimistically` counter mencatat totalqty yang successful via optimistic path
- Expected: `stock == initial - applied` → `100 - 20 = 80`
- Actual: stock = 80 (semua 20 goroutines eventually succeed via retry)

### Test: Atomic Conditional Update

```go
func TestAtomicConditionalUpdate(t *testing.T) {
    store := inventory.NewStore()
    store.Seed(5, "Item E", 100)
    svc := inventory.NewService(store)
    
    const goroutines = 50
    var wg sync.WaitGroup
    wg.Add(goroutines)
    
    for i := 0; i < goroutines; i++ {
        go func() {
            defer wg.Done()
            err := svc.DeductAtomic(5, 1)
            if err != nil {
                t.Errorf("unexpected atomic deduct failure: %v", err)
            }
        }()
    }
    wg.Wait()
    
    p, _ := store.Get(5)
    if p.Stock != 50 {
        t.Fatalf("expected stock 50, got %d", p.Stock)
    }
}
```

- 50 goroutines concurrent memanggil `DeductAtomic(5, 1)`
- Expected: stock = 50 (exact invariant)
- Actual: stock = 50 (success—single statement atomicity)

## What the Tests Prove

| Test Name | Focus Path | Invariant Checked | Result |
|---|---|---|---|
| `TestNaiveLostUpdate` | Failure / Anomaly | Final stock > 50 despite 50 deducts | PASS (lost update demonstrated) |
| `TestPessimisticLocking` | Happy Path / Concurrency | Exact stock (50) under 50 goroutines | PASS (fully synchronized) |
| `TestPessimisticLockingInsufficientStock` | Negative / Edge Case | Stock boundary error when stock < qty | PASS (error returned correctly) |
| `TestOptimisticLockingConflict` | Conflict Detection | Successes + Stock == 100; Conflicts > 0 | PASS (state guarded, zero corruption) |
| `TestOptimisticLockingWithRetry` | Recovery & Convergence | Stock == 100 - Optimistically applied | PASS (all retries converged) |
| `TestAtomicConditionalUpdate` | Happy Path / Concurrency | Stock exactly matches 100 - 50 = 50 | PASS (lockless single statement) |

**Race detector:** `go test -race ./...` → PASS (zero race conditions detected)

**Demo output:**
```
[1] Naive Read-Modify-Write (50 concurrent requests):
    Initial Stock: 100
    Actual Final Stock:   99 (LOST UPDATE DETECTED!)

[2] Pessimistic Locking (SELECT ... FOR UPDATE):
    Actual Final Stock:   50 (SUCCESS - Fully Synchronized)

[3] Optimistic Locking Direct (20 concurrent requests, no retry):
    Successful Deductions: 1
    Rejected Conflicts:   19
    Actual Final Stock:   99 (SUCCESS - State Guarded, Zero Corruption)

[4] Optimistic Locking With Exponential Backoff Retry (20 requests):
    Successful Deductions: 20
    Total Attempted Conflicts Retried: ~61
    Actual Final Stock:   80 (SUCCESS - All retries eventually converged)

[5] Atomic Single-Statement Operation:
    Actual Final Stock:   50 (SUCCESS - Lockless Single Statement)
```

## Recovery / Rollback

### Optimistic Locking Recovery

Pada konflik (`ErrOptimisticLock`), recovery path:

1. **Reload** data dari database (fetch latest version)
2. **Recalculate** business logic dengan latest values
3. **Retry** update dengan version baru
4. **Backoff** (jittered exponential) sebelum retry

```go
for attempt := 0; attempt <= maxRetries; attempt++ {
    err := svc.store.OptimisticDeduct(id, qty)
    if err == nil {
        return nil // success
    }
    if err != ErrOptimisticLock {
        return err // non-optimistic error (not found, insufficient stock)
    }
    if attempt == maxRetries {
        return ErrOptimisticLock // max retries reached
    }
    // Jittered exponential backoff
    sleepDuration := time.Duration(1<<attempt)*time.Millisecond + time.Duration(rand.Intn(5))*time.Millisecond
    time.Sleep(sleepDuration)
}
```

Alternative: return `409 Conflict` ke client dan biarkan client retry (dengan backoff atau manual retry button).

### Pessimistic Locking Recovery

Pessimistic locking tidak perlu recovery di application layer karena conflict tidak terjadi—write diblokir sampai lock available.

**Butuh attention:** Deadlock detection (database auto-abort satu transaction) dan lock hold time (jangan pegang lock saat network call).

### Atomic Update Recovery

Atomic update tidak bisa "concurrent" dalam arti conflict—statement itu atomic. Satu statement bisa fail karena:
- Product tidak ditemukan (`ErrNotFound`)
- Stock tidak cukup (`ErrInsufficientStock`)

Recovery: error handling di application level.

## Production Considerations

### Pessimistic Locking

**When to use:**
- High contention (banyak concurrent writes ke resource yang sama)
- Correctness-critical (balances, seat reservations, queue numbers)
- Write frequency mendekati read frequency

**Production tips:**
- Keep transactions short (jangan pegang lock saat HTTP call)
- Consistent lock ordering (mencegah deadlock)
- Monitor lock wait time (high wait time = contention problem)
- Consider `NOWAIT` / `SKIP LOCKED` untuk non-blocking behavior
- Lock escalation (row → page → table) terjadi di beberapa database—monitor

**Anti-patterns:**
- Holding lock di goroutine yang tunggu HTTP response
- Lock acquired tapi deferUnlock dipanggil di deferred function yang dieksekusi setelah network call
- Lock ordering tidak konsisten (goroutine A lock X lalu Y, goroutine B lock Y lalu X → deadlock)

### Optimistic Locking

**When to use:**
- Low contention (banyak reads, sedikit writes)
- Read-heavy workloads (profile edit, CMS, CRM, master data)
- Business transactions yang span multiple requests (multi-step wizard, shopping cart)
- Scalability critical (avoid blocking)

**Production tips:**
- Retry logic wajib (exponential backoff dengan jitter)
- Version column wajib (integer atau timestamp)
- Handle 0-rows-affected sebagai conflict, bukan success
- Monitor retry rate (high retry rate = need pessimistic or caching)

**Anti-patterns:**
- Ignoring 0-rows-affected (silently lose update)
- Retry tanpa backoff (thundering herd)
- Retry tanpa reload (retry dengan stale data)

### Atomic Single-Statement

**When to use:**
- Simple counter/quota/stock decrement
- No business logic antara read dan write
- Lockless desired (avoid lock contention entirely)

**Production tips:**
- Gunakan single conditional UPDATE
- Check `affected_rows` untuk success/failure
- Jangan gunakan untuk complex business logic (perlu read-modify-write dengan external API call)

**Anti-patterns:**
- Using for complex logic (misal: deduct stock, create order, send email—ini bukan single statement)
- Ignoring insufficient stock condition (stock < qty)

### General

**Isolation level alone tidak cukup:**
- READ COMMITTED (PostgreSQL, Oracle) = lost update mungkin tanpa explicit locking
- REPEATABLE READ (MySQL) = lost update mungkin tanpa WHERE guard atau explicit locking
- SERIALIZABLE = lost update must throw error (tapi overhead tinggi)

**Selection heuristic:**
1. **Atomic update** jika simple counter/quota
2. **Optimistic locking** jika low contention, read-heavy, multi-request workflow
3. **Pessimistic locking** jika high contention, correctness-critical, single-request workflow
4. **Distributed locks** (Redis) hanya jika resource span multiple databases atau microservices—prefer database-native mechanism dulu

## Common Mistakes

### 1. Transaction Alone ≠ Lost Update Protection

**Mistake:** "Kita pakai `BEGIN ... COMMIT`—data kita aman."

**Reality:** Transaction tanpa explicit locking atau version guard atau atomic statement *tidak mencegah* lost update di default isolation.

```sql
-- MISTAKE: lost update tetap mungkin
BEGIN;
SELECT stock FROM products WHERE id = 10;
-- ... compute new_stock ...
UPDATE products SET stock = new_stock WHERE id = 10;
COMMIT;
```

**Fix:** Pakai `SELECT ... FOR UPDATE` atau version guard atau atomic update.

### 2. Distributed Locks When Database Can Solve It

**Mistake:** "Kita pakai Redis lock untuk stock update."

**Reality:** Database-native mechanism (row locks, version guard) lebih simple, lebih fast, dan lebih reliable kecuali resource span multiple databases.

PostgreSQL advisory locks doc: "A common use of advisory locks is to emulate pessimistic locking strategies... While a flag stored in a table could be used, advisory locks are faster"—implies database-native mechanism preferred.

**Fix:** Pakai `SELECT ... FOR UPDATE` dulu, baru pertimbangkan distributed locks untuk cross-database scenario.

### 3. Holding Locks Across Network Calls

**Mistake:** "Kita lock product, lalu panggil payment gateway, lalu update stock."

**Reality:** Holding lock selama network call = bad idea (PostgreSQL: "do not hold transactions open waiting for user input"). Ini:
- Block semua writer lain
- Increase deadlock risk
- Network failure = lock leak (jika tidak rollback)

**Fix:**
- Lock → read → unlock → network call → lock (jika perlu) → write
- Atau gunakan optimistic locking (tidak perlu lock saat network call)

### 4. Ignoring 0-Rows-Affected pada Optimistic Update

**Mistake:** "UPDATE succeeds (no error), tapi 0 rows affected—kita anggap success."

**Reality:** `affected_rows == 0` = version mismatch = konflik terdeteksi. Jangan ignore.

```go
// MISTAKE
_, err := db.Exec("UPDATE products SET stock = stock - 1, version = version + 1 WHERE id = ? AND version = ?", id, version)
if err != nil {
    return err // error handling
}
// no check for affected_rows == 0!
```

**Fix:** Check `affected_rows`, return `ErrOptimisticLock` jika 0.

### 5. Naive Read-Modify-Write Tanpa Concurrency Guard

**Mistake:** `$product->stock -= 1; $product->save();` tanpa version check atau locking.

**Reality:** Ini adalah source code untuk lost update anomaly.

**Fix:** Gunakan salah satu:
- Atomic update: `UPDATE products SET stock = stock - 1 WHERE id = ? AND stock >= 1`
- Optimistic locking: version check di WHERE
- Pessimistic locking: `SELECT ... FOR UPDATE` sebelum UPDATE

## Case Study

### Scenario: E-commerce Inventory System

**Problem:** Product "Laptop Pro" stock = 100. 50 concurrent customer requests masing-masing beli 1 unit.

**Naive approach (lost update):**
- 50 concurrent deduct calls
- Final stock: 99 (expected 50)
- Lost update: 1 unit (data corrupt, tidak ada error)

**Pessimistic locking solution:**
- 50 concurrent deduct calls
- Final stock: 50 (exact)
- Trade-off: Writers block, but no data corruption

**Optimistic locking solution:**
- 20 concurrent deduct calls (untuk test)
- Successful: 1, Conflicts: 19
- With retry: 20 successful (all goroutines converge)
- Trade-off: Retry overhead, but high concurrency

**Atomic single-statement solution:**
- 50 concurrent deduct calls
- Final stock: 50 (exact)
- Trade-off: Lockless, but only untuk simple decrement

**Selection decision:**
- High contention? Ya (50 concurrent requests)
- Correctness-critical? Ya (stock accuracy)
- Choice: **Pessimistic locking** atau **Atomic update** (jika simple decrement, no business logic)

**Production metric:**
- Monitor retry rate (optimistic) / lock wait time (pessimistic)
- High retry rate = need more pessimistic
- High wait time = need caching atau sharding

## Checklist

Before deployment:
- [ ] Lost update scenario reproduced (test with concurrent goroutines)
- [ ] Pessimistic locking test passes (exact invariant under concurrency)
- [ ] Optimistic locking conflict detected (affected_rows == 0 handled)
- [ ] Optimistic locking retry converges (backoff + jitter)
- [ ] Atomic update test passes (lockless single statement)
- [ ] Race detector passes (`go test -race ./...`)
- [ ] Demo runs successfully (`go run ./cmd/demo`)
- [ ] Error messages clear (`ErrOptimisticLock`, `ErrInsufficientStock`)
- [ ] Selection criteria documented (atomic/optimistic/pessimistic)
- [ ] Anti-patterns avoided (transaction alone, lock across network call, ignore 0-rows-affected)
- [ ] Monitoring in place (retry rate, lock wait time, stock discrepancy reports)

## Key Takeaways

1. **Lost update adalah silent data corruption**—dua transaksi konkurensi menimpa satu sama lain tanpa error.
2. **Transaction alone tidak cukup**—query design (locking, version guard, atau atomic statement) menentukan correctness.
3. **Pessimistic locking** (`SELECT ... FOR UPDATE`) = prevent conflict (blokir writer lain).
4. **Optimistic locking** (version guard + affected_rows check) = detect conflict (reload/recalculate/retry).
5. **Atomic single-statement** (`UPDATE ... SET stock = stock - N WHERE stock >= N`) = eliminate race window entirely.
6. **Selection criteria:** Atomic first (simple counter), optimistic (low contention/read-heavy), pessimistic (high contention/correctness-critical).
7. **Anti-patterns:** Transaction alone, holding lock across network call, ignoring 0-rows-affected, distributed locks when database can solve it.
8. **Isolation level:** READ COMMITTED (PostgreSQL, Oracle) = lost update mungkin; REPEATABLE READ (MySQL) = lost update mungkin tanpa guard; SERIALIZABLE = must throw error.
9. **Production tip:** Monitor retry rate (optimistic) / lock wait time (pessimistic).
10. **Demo test:** Naive → lost update (stock 99), Pessimistic → exact (stock 50), Optimistic → guarded (stock 99, 1 success), Optimistic+retry → converged (stock 80), Atomic → lockless (stock 50).

## Sources

### Research

- **Research Plan:** `research/01-plan.md`
- **Evidence:** `research/03-evidence.md`
- **Sources:** `research/02-sources.md`
- **Contradictions:** `research/04-contradictions.md`
- **Open Questions:** `research/06-open-questions.md`
- **Report:** `research/05-report.md`

### Engineering

- **Design:** `engineering/01-design.md`
- **Implementation Notes:** `engineering/02-implementation-notes.md`
- **Execution Result:** `engineering/03-execution-result.md`

### Audit

- **Research Audit Plan:** `research-audit/01-audit-plan.md`
- **Source Audit:** `research-audit/02-source-audit.md`
- **Claim Audit:** `research-audit/03-claim-audit.md`
- **Contradictions Audit:** `research-audit/04-contradictions.md`
- **Code Audit:** `research-audit/05-code-audit.md`
- **Gaps:** `research-audit/06-gaps.md`
- **Verdict:** `research-audit/07-verdict.md`

- **Engineering Audit Plan:** `engineering-audit/01-audit-plan.md`
- **Code Audit:** `engineering-audit/02-code-audit.md`
- **Test Audit:** `engineering-audit/03-test-audit.md`
- **Docs vs Code:** `engineering-audit/04-docs-vs-code.md`
- **Gaps:** `engineering-audit/05-gaps.md`
- **Verdict:** `engineering-audit/06-verdict.md`

### Revision

- **Research Revision Plan:** `research-revision/01-revision-plan.md`
- **Changes Made:** `research-revision/02-changes-made.md`
- **Revision Result:** `research-revision/03-revision-result.md`

- **Engineering Revision Plan:** `engineering-revision/01-revision-plan.md`
- **Changes Made:** `engineering-revision/02-changes-made.md`
- **Revision Result:** `engineering-revision/03-revision-result.md`

### Source Code

- **Model:** `internal/inventory/model.go`
- **Store:** `internal/inventory/store.go`
- **Service:** `internal/inventory/service.go`
- **Tests:** `tests/locking_test.go`
- **Demo:** `cmd/demo/main.go`