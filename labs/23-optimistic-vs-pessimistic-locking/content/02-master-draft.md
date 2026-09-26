# Optimistic vs Pessimistic Locking & Atomic Updates

## Problem

Ketika dua proses bersamaan membaca stok produk (misalnya 100 unit), mengurangi secara lokal (masing-masing kurangi 1), lalu menulis kembali, transaksi kedua dapat menimpa hasil transaksi pertama tanpa memberi tahu.hasil ini disebut **lost update** — anomali korupsi data diam-diam yang tidak menghasilkan error.

Contoh: Thread A baca stok=100, kurangi menjadi 99. Thread B baca stok=100 (belum commit A), kurangi menjadi 99. Kedua tulis 99. Stok akhir = 99, padahal harus = 98.

## Why This Matters

Lost update melanggar invariant bisnis (stok = awal - total ketika). Dalam sistem inventory, keuangan, atau reservasi, korupsi ini menyebabkan:

- Laporan keuangan tidak akurat
- Kesalahan persediaan
- Double-booking (seat, tiket)
- Kerugian uang nyata

Database dengan isolasi READ COMMITTED (default PostgreSQL, Oracle) tidak mencegahnya jika aplikasi menggunakan pola **read-modify-write** tanpa penguncian eksplisit.

## Mental Model

```text
┌────────────────┐     ┌────────────────┐
│   TRANSaksi 1  │     │   Transaksi 2  │
│  READ (stock)  │     │  READ (stock)  │
│   ↆ 100         │     │   ↆ 100        │
│  CALCULATE     │     │  CALCULATE     │
│  100-1=99      │     │  100-1=99      │
│   ↇ WRITE 99   │────▶│   ↇ WRITE 99   │
│                │     │                │
│   COMMIT       │     │   COMMIT       │
└────────────────┘     └────────────────┘
    Akhir: 99 ❌     Harus: 98
```

## Core Concept

Tiga strategi utama untuk mencegahnya:

### 1. Pessimistic Locking

`SELECT ... FOR UPDATE` memegang **row-level exclusive lock** hingga commit/rollback. Thread kedua menunggu sampai pertama selesai.

```sql
BEGIN;
SELECT stock FROM products WHERE id=1 FOR UPDATE;
-- stok ada di memori, tidak boleh diubah orang lain
UPDATE products SET stock=99 WHERE id=1;
COMMIT;
```

**Karakteristik:** Mencegah konflik dengan memblokir. Cocok untuk data yang sangat kritis (balance, stok terbatas, nomor antrian).

### 2. Optimistic Locking

Versi/timestamp disimpan di baris. UPDATE memeriksa versi lama di WHERE clause:

```sql
BEGIN;
SELECT stock, version FROM products WHERE id=1;  -- version=5
UPDATE products 
  SET stock=99, version=6 
  WHERE id=1 AND version=5;  -- zero rows affected = konflik
COMMIT;
```

Jika `affected_rows = 0`, ada konsekuensi stok yang berbeda. Aplikasi harus retry atau return 409 Conflict.

**Karakteristik:** Deteksi konflik, bukan pencegahan. Cocok untuk beban baca-banyak, konflik jarang (CMS, profil pengguna).

### 3. Atomic Single-Statement Operation

Hilangkan jeda read-modify-write lewat satu pernyataan yang atoms:

```sql
UPDATE products SET stock = stock - 3 WHERE id=1 AND stock >= 3;
-- cek affected_rows == 1 untuk konfirmasi
```

**Karakteristik:** Paling simpel untuk operasi sederhana (decrement kuantitas). Tidak memerlukan transaction panjang atau manual lock.

## Architecture

```text
┌─────────────────────────────────────────┐
│  cmd/demo/main.go  (CLI Demo Runner)    │
│                                         │
│  tests/locking_test.go (Concurrency)    │
│                                         │
│  ┌───────────────────────────────────┐  │
│  │  internal/inventory/              │  │
│  │  ├── model.go   (Product, errors) │  │
│  │  ├── store.go   (simulated DB)  │  │
│  │  └── service.go  (business ops)  │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
```

## Implementation

### Store — Simulated Engine (`internal/inventory/store.go`)

- `rowLocks map[int]*sync.Mutex` — tiap row punya mutex terpisah, meniru `SELECT ... FOR UPDATE`
- `products map[int]*Product` — data terseksi dengan field `Stock` dan `Version`
- `NaiveDeduct` — read, sleep, write; tidak aman
- `PessimisticDeduct` — acuisisi row lock, validasi, update; aman
- `OptimisticDeduct` — baca + version check di dalam lock; return `ErrOptimisticLock` bila konflik
- `AtomicDeduct` — satu operasi lock-unlock dengan update terpisah; aman

### Service Layer (`internal/inventory/service.go`)

- `DeductNaive` — wrapper langsung
- `DeductPessimistic` — wrapper yang aman
- `DeductOptimisticDirect` — tanpa retry
- `DeductOptimisticWithRetry(id, qty, maxRetries)` — loop retry dengan jittered exponential backoff `(1<<attempt)ms + rand(0-5ms)`
- `DeductAtomic` — wrapper untuk operasi tunggal

## Code Walkthrough

### Naive (Bug)

```go
func (s *Store) NaiveDeduct(id int, qty int) error {
    p, err := s.Get(id)           // 1. READ
    if err != nil { return err }
    if p.Stock < qty { return ErrInsufficientStock }
    time.Sleep(100 * time.Microsecond)  // jendela konkurensi
    s.mu.Lock()
    curr := s.products[id]
    curr.Stock = p.Stock - qty      // 2. WRITE pakai nilai lama
    s.mu.Unlock()
    return nil
}
```

### Pessimistic (Correct)

```go
func (s *Store) PessimisticDeduct(id int, qty int) error {
    rowLock := s.GetRowLock(id)
    rowLock.Lock()                  // BLOCK sampai selesai
    defer rowLock.Unlock()
    s.mu.Lock()
    p, exists := s.products[id]
    if !exists { s.mu.Unlock(); return ErrNotFound }
    if p.Stock < qty { s.mu.Unlock(); return ErrInsufficientStock }
    p.Stock -= qty                // update aman
    s.mu.Unlock()
    return nil
}
```

### Optimistic (Detection)

```go
func (s *Store) OptimisticDeduct(id int, qty int) error {
    p, err := s.Get(id)
    if err != nil { return err }
    if p.Stock < qty { return ErrInsufficientStock }
    time.Sleep(50 * time.Microsecond)    // jendela konkurensi
    s.mu.Lock()
    defer s.mu.Unlock()
    curr := s.products[id]
    if curr.Version != p.Version {        // VERSION MISMATCH
        s.OptimisticFails++
        return ErrOptimisticLock
    }
    curr.Stock -= qty
    curr.Version++
    return nil
}
```

### Atomic (Simplest)

```go
func (s *Store) AtomicDeduct(id int, qty int) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    curr, exists := s.products[id]
    if !exists { return ErrNotFound }
    if curr.Stock < qty { return ErrInsufficientStock }
    curr.Stock -= qty                  // UPDATE tunggal, aman
    return nil
}
```

## What the Tests Prove

| Test | Apa yang Dibuktikan | Hasil |
|------|---------------------|-------|
| `TestNaiveLostUpdate` | 50 deduct, stok ≠ 50 | PASS (stok=99) |
| `TestPessimisticLocking` | 50 deduct menggunakan lock | PASS (stok=50) |
| `TestPessimisticLockingInsufficientStock` | Refuse understock | PASS |
| `TestOptimisticLockingConflict` | Konflik terdeteksi | PASS (1 sukses, 19 konflik) |
| `TestOptimisticLockingWithRetry` | Retry konvergen | PASS (semua 20 berhasil) |
| `TestAtomicConditionalUpdate` | Atomic update aman | PASS (stok=50) |

Semua tes + race detector (`go test -race ./...`) lulus tanpa warning.

## Recovery / Rollback

- **Pessimistic:** Deadlock dapat terjadi. PostgreSQL auto-detect dan abort satu transaksi secara acak. Defensa: urutan lock konsisten, transaksi pendek.
- **Optimistic:** Konflik dikembalikan sebagai `ErrOptimisticLock`. Aplikasi harus implementasi retry atau kembalikan 409 Conflict ke klien.
- **Atomic:** Jika kondisi tidak terpenuhi (`stock < qty`), return error. Bukan korupsi, melainkan penolakan bisnis normal.

## Production Considerations

- **Pessimistic:** gunakan hanya pada sumber daya kritis dengan tingkat kontensi tinggi. Jangan pegang lock selama call eksternal (HTTP, payment gateway) — hal ini adalah anti-pattern.
- **Optimistic:** cocok untuk beban baca-banyak, konflik jarang. Implementasi retry wajib ada. Pilih 64-bit atau timestamp untuk version counter bila frekuensi tinggi.
- **Atomic:** pilihan paling simpel untuk operasi kuantitas. Tidak cocok bila validasi kompleks diperlukan sebelum update.
- **Isolation Level:** REPEATABLE READ/SERIALIZABLE di PostgreSQL melempar error bila terdeteksi write-skew. Namun, READ COMMITTED (default) masih rentan lost update.

## Common Mistakes

1. **Menganggap transaksi tunggal sudah cukup** — tidak, READ COMMITTED tidak mencegah lost update.
2. **Holding lock saat HTTP call** — menunda commit, meningkatkan window race + deadlock.
3. **Mengabaikan `affected_rows = 0`** — konflik optimis tidak ditangan, stok tidak berubah padahal aplikasi bilang berhasil.
4. **Naive read-modify-write** — pola `load() + modify() + save()` tanpa guard.
5. **Menggunakan Redis lock sebelum coba database-native lock** — PostgreSQL advisory locks lebih ringan; gunakan external lock hanya bila resource memang lintas database.

## Case Study

Demo CLI menampilkan Lima skenario berdampingan (stok awal = 100):

```
[1] Naive Read-Modify-Write (50 concurrent requests):
    Expected Final Stock: 50
    Actual Final Stock:   99 (LOST UPDATE DETECTED!)

[2] Pessimistic Locking (SELECT ... FOR UPDATE):
    Actual Final Stock:   50 (SUCCESS)

[3] Optimistic Locking Direct (20 concurrent requests, no retry):
    Successful Deductions: 1
    Rejected Conflicts:   19
    Actual Final Stock:   99 (State Guarded, Zero Corruption)

[4] Optimistic Locking With Exponential Backoff Retry (20 requests):
    Successful Deductions: 20
    Actual Final Stock:   80 (All retries converged)

[5] Atomic Single-Statement Operation:
    Actual Final Stock:   50 (Lockless Single Statement)
```

## Checklist

- [ ] Tentukan tingkat kontensi (tinggi/biasa/ rendah)
- [ ] Pilih strategi:
  - Kontensi tinggi atau data kritis → Pessimistic
  - Kontensi rendah, beban baca-banyak → Optimistic dengan retry
  - Operasi sederhana (counter, stok) → Atomic
- [ ] Jangan pegang lock selama I/O jaringan
- [ ] Implementasikan retry exponential backoff untuk optimistic
- [ ] Selalu cek `affected_rows` atau hasil error untuk konflik
- [ ] Jalankan race detector: `go test -race ./...`

## Key Takeaways

1. Lost update adalah korupsi data diam-diam yang terjadi pada read-modify-write tanpa penguncian.
2. Pessimistic locking (SELECT FOR UPDATE) mencegah konflik dengan memblokir baris.
3. Optimistic locking mendeteksi konflik via version guard + affected_rows check.
4. Atomic single-statement UPDATE menghilangkan jendela konkurensi.
5. Isolation level tidak cukup; query design adalah kunci konsistensi.
6. Pilih strategi berdasarkan frekuensi konflik dan kritis data.
7. Deadlock dapat terjadi pada pessimistic; transaksi harus pendek.
8. Optimistic memerlukan retry atau 409 Conflict response.
9. Race detector harus lulus tanpa warning.
10. Nilai stock akhir harus selalu = awal - total successful deduction.

## Sources

- Wikipedia Concurrency Control — definisi lost update, kategori optimistic/pessimistic
- Martin Fowler Optimistic Offline Lock — asumsi konflik jarang, pre-commit validation
- Martin Fowler Pessimistic Offline Lock — mencegah konflik dengan membeli lock
- PostgreSQL 18 Docs — 13.2 Transaction Isolation, 13.3 Explicit Locking, 13.4 Application-Level Consistency
- MySQL 8.0 Docs — InnoDB locking reads, transaction isolation levels
- Oracle 19c Concepts — data concurrency, row locks TX, WHERE-guard pattern