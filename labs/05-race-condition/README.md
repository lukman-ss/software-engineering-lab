# Lab 05 — Race Condition: Bug yang Tidak Muncul Saat Development, Tapi Menghancurkan Sistem di Production

> **Mental Model**: Concurrent execution dapat merusak correctness. Kode yang benar secara sequential belum tentu benar ketika dieksekusi concurrent terhadap shared state.

---

## Mental Model

**Kunci utama**: Business invariant = kontrak yang harus selalu benar, bahkan di tengah concurrent chaos.

Jika concurrent execution bisa membroken invariant, itu adalah **race condition pada level aplikasi** — bukan hanya memory data race di Go.

---

## Studi Kasus Inventory: Oli Mesin

**Setup**: Toko memiliki 1 unit Oli Mesin.

**Inventory invariant**:
```
initial_stock == successful_sales + final_stock
final_stock >= 0
successful_sales <= initial_stock
```

---

## Sequential vs Concurrent Execution

### Sequential (benar)

```
Request A
  READ stock = 1
  CHECK > 0? YES
  CALCULATE new = 0
  WRITE stock = 0
  Result: 1 successful sale, final_stock = 0 ✅
```
Invariant: `1 == 1 + 0` ✅

### Concurrent (runtuh)

```
Request A              Request B
   │                        │
T0 │ READ stock = 1 ◄──────┼── READ stock = 1
   │                        │
T1 │ CHECK > 0? YES   ┌────▼────┐
   │                   │  (B blocked) │
T2 │ CALCULATE 0 ◄────┼─WRITE 0─┼─► CALCULATE 0
   │                   │         │
T3 │ WRITE 0 ─────────┼─────────┼─► WRITE 0
   │                        │
Result:
  successful_sales = 2  ❌
  final_stock = 0
  1 ≠ 2 + 0  →  Invariant broken ❌
```

---

## Definisi Race Condition

**Application/Business Race Condition**: Ketika outcome bergantung pada timing/interleaving eksekusi concurrent terhadap shared state.

Ini **bukan** memory data race. Semua variable lokal Go dapat thread-safe, tapi business state tetap corrupt.

---

## Shared State

Race condition terjadi ketika ada **shared mutable state** yang diakses concurrent.

Untuk inventory: `stock` di database adalah shared mutable state.

Untuk booking: `slot_available` adalah shared mutable state.

---

## Check-Then-Act Pattern (Unsafe)

Pattern ini menyebabkan race condition pada business state:

```go
func TrySell(ctx context.Context, productID string) error {
    stock, err := repo.GetStock(ctx, productID) // READ
    if err != nil {
        return err
    }
    if stock <= 0 {                              // CHECK
        return ErrOutOfStock
    }
    newStock := stock - 1                        // CALCULATE
    return repo.SetStock(ctx, productID, newStock) // WRITE
}
```

**Mengapa ini tidak safe?**

- `GetStock` dan `SetStock` masing-masing thread-safe (Go memory model)
- Tapi **sequence READ → CHECK → CALCULATE → WRITE bukan atomic**
- Antara READ dan WRITE, goroutine lain dapat mengubah state
- Kedua goroutine membaca stale value, menulis nilai yang sama

---

## Lost Update

**Lost update** adalah consequence langsung dari check-then-act race + concurrent stale reads.

### Timeline

```
T0  Request A membaca stock = 1
T1  Request B membaca stock = 1 (stale)
T2  Request A menghitung newStock = 0
T3  Request B menghitung newStock = 0
T4  Request A menulis stock = 0
T5  Request B menulis stock = 0
```

### Hubungan sebab-akibat

```
check-then-act race
    ↓
concurrent stale reads
    ↓
lost update
    ↓
broken business invariant
```

---

## Broken Business Invariant

Setelah race condition:

```
initial_stock = 1
successful_sales = 2
final_stock = 0

1 ≠ 2 + 0
```

**Masalahnya bukan stock = -1.** Masalahnya adalah **konservasi jumlah**: sistem mengira 2 unit terjual padahal stock hanya 1.

---

## Kenapa Sulit Direproduksi?

1. **Timing-dependent**: Race hanya terjadi pada interleaving spesifik
2. **Environment-dependent**: Development biasanya single-user
3. **Latency masking**: Slow database di production mengeksponen race condition
4. **Flaky reproduction**: `time.Sleep` tidak dapat mereproduksi race secara deterministik

---

## Unsafe Implementation

`unsafe_inventory.go` implements check-then-act pattern menggunakan local variable yang Go-memory-safe, tapi business state tidak safe.

```go
func (s *InventoryService) TrySell(ctx context.Context, productID string) error {
    stock, err := s.repo.GetStock(ctx, productID)
    if err != nil {
        return err
    }
    if stock <= 0 {
        return ErrOutOfStock
    }
    newStock := stock - 1
    return s.repo.SetStock(ctx, productID, newStock)
}
```

---

## Deterministic Reproduction

`lost_update_test.go` mereproduksi lost update tanpa `time.Sleep`.

Menggunakan **channel barrier synchronization** untuk mengontrol timing:
- Request A dan B **must diblokir sampai keduanya selesai READ**
- Kemudian baru melanjutkan ke CALCULATE → WRITE

```go
// aReadDone signal ketika A selesai READ
// bReadDone signal ketika B selesai READ
// Kedua goroutine diblokir sampai "read phase" selesai
// Baru lanjut ke calculate/write phase
```

**Hasil deterministic**: `successful_sales = 2, final_stock = 0` — **every time**.

---

## Safe Solution: Atomic Conditional Update

```sql
-- PostgreSQL style
UPDATE products SET stock = stock - 1
WHERE id = $1 AND stock > 0
RETURNING stock;
```

### RowsAffected / RETURNING

- **1 row affected** → decrement sukses
- **0 row affected** → stock habis, condition gagal

### Keuntungan

- Tidak ada SELECT terlebih dahulu → tidak ada stale read window
- Read-modify-write atomik di database
- Concurrent execution di-handle oleh database engine

`atomic_inventory.go` implements pola ini.

---

## Transaction Is Not Magic

**BEGIN/COMMIT tidak otomatis memperbaiki race condition.**

Isolation level `READ COMMITTED` (default PostgreSQL) tidak menjamin serialisasi.

```
Transaction A                    Transaction B
    │                                 │
    ▼ SELECT stock = 1                ▼ SELECT stock = 1 (stale read)
    ▼ UPDATE stock = 0                ▼ UPDATE stock = 0 (overwrite!)
    ▼ COMMIT                          ▼ COMMIT
```

### Isolation Awareness

| Property | Traditional Code | Transaction | Atomic Conditional |
|----------|-----------------|-------------|-------------------|
| **Atomicity** | Non-atomic | Terpisah | Atomik |
| **Isolation** | N/A | Seperah 1 isolation level | UPDATE with WHERE |
| **Locking** | Manual | Row lock | Database internal |
| **Invariant** | Bisa rusak | Bisa rusak | Terjamin |

---

## Secondary Scenario: Booking Race

### Service Booking Online

```
Customer A → SELECT slot 09:00 kosong
Customer B → SELECT slot 09:00 kosong
Customer A → INSERT booking A
Customer B → INSERT booking B
Result: 2 bookings pada slot yang sama
```

**Broken Invariant**: `COUNT(bookings per slot) <= 1`

### Solusi: Unique Constraint

```sql
CREATE TABLE service_bookings (
    customer_id   VARCHAR(50),
    service_date  DATE,
    slot_time     TIME,
    UNIQUE(service_date, slot_time)
);
```

Frontend validation → **UX protection**
Database constraint → **invariant protection**

Lihat `booking_test.go` untuk `Test500_ConcurrentBooking`.

---

## Invoice Number Race Example

```sql
SELECT MAX(invoice_no) FROM invoices;
-- App calculates: newNo = max + 1
-- Request A: MAX = 100, app calculates 101
-- Request B: MAX = 100, app calculates 101
-- Duplicate invoice number!
```

**MAX() + 1 rawan race.** Gunakan database sequence untuk internal sequencing.

---

## Pessimistic Locking Overview

```sql
SELECT * FROM products WHERE id = $1 FOR UPDATE;
```

### Timeline

```
Transaction A: lock → read → update → commit
Transaction B: wait (blocked on row lock)
```

### Trade-off

- **Blocking**: Request B menunggu
- **Contention**: Antrian panjang pada hot rows
- **Latency**: Lock wait menambah response time
- **Connection occupancy**: Koneksi terpakai selama lock
- **Deadlock risk**: Circular wait harus dihindari

> Implementasi penuh di [Lab 11 — Pessimistic Locking](../11-pessimistic-locking/)

---

## Optimistic Locking Overview

```sql
UPDATE products
SET stock = $newStock, version = version + 1
WHERE id = $id AND version = $oldVersion;
```

### Flow

```
read version → calculate → conditional write → rows_affected = 0? → conflict
```

Cocok untuk high concurrency dengan low contention.

> Implementasi penuh di [Lab 12 — Optimistic Locking](../12-optimistic-locking/)

---

## Process-local Synchronization

`sync.Mutex`, `sync/atomic`, channel hanya mengkoordinasikan **goroutine dalam satu process**.

```
Instance A → mutex A
Instance B → mutex B  (mutex berbeda, tidak saling koordinasikan)
Instance C → mutex C
```

`sync.Mutex` **bukan distributed lock**. Ia tidak bekerja untuk web app multi-instance.

### Kapan mutex/atomic/channel tepat pakai?

- Shared memory dalam satu process
- Single instance deployment
- Actor model (single goroutine owns state)

---

## Distributed Coordination

Untuk koordinasi lintas instance, butuh mekanisme terpisah.

### Contoh penggunaan

- Singleton cron: satu scheduler aktif
- Report generation: hindari duplicate generation
- Cross-instance critical section
- External process coordination

### Distributed lock konsep

- **Ownership**: Hanya pemilik yang boleh membuka kunci
- **TTL/Lease**: Otomatis expired jika pemilik crash
- **Expiration**: Bukan sekadar lock, tapi time-bounded
- **Safe release**: Hanya pemilik bisa melepas kunci
- **Fencing token** (advanced): Token unik setiap lease untuk mencegah stale lock

### Rule

> **Jika invariant bisa ditegakkan via database (atomic statement atau constraint), gunakan mekanisme database.** Distributed lock adalah pilihan terakhir, bukan default.

---

## Memory Data Race vs Business Race Condition

### Memory Data Race

Dua goroutine mengakses memory address yang sama tanpa synchronization, dengan setidaknya satu write.

```go
var counter int
go func() { counter++ }()
go func() { counter++ }()
```

**Deteksi**: `go test -race`

### Business/Application Race Condition

Correctness bergantung pada interleaving concurrent operations terhadap shared state.

```
Request A → SELECT stock = 1
Request B → SELECT stock = 1
```

Semua Go local variables bisa race-free, tapi **business state corrupt**.

### Key Insight

✅ **Bebas memory data race** (`go test -race` passes)  
❌ **Bisa masih ada business race condition**

> `go test -race` **tidak membuktikan** business correctness.

---

## Non-Database Race Conditions

Race condition juga bisa terjadi pada local shared state:

```
Service A (local cache):
  cache["user_123"] = {balance: 100}
  
Request A: cache["user_123"].balance = 80 (after purchase)
Request B: cache["user_123"].balance = 90 (after refund)
```

Jika cache tidak konsisten dengan database, race condition masih terjadi.

---

## Testing Concurrent Systems

### Pendekatan yang benar

1. **Deterministic reproduction** — Gunakan channel, WaitGroup, barrier
2. **Stress testing** — Jalankan 500+ concurrent workers
3. **Invariant checking** — Assert global invariant setelah semua worker selesai
4. **Repeatable** — Jalankan dengan `-count=20` untuk kepastian

### Pendekatan yang salah

1. ⛔ `time.Sleep` sebagai concurrency mechanism
2. ⛔ Mengandalkan execution order
3. ⛔ Hanya menguji single operation, bukan global state
4. ⛔ `go test -race` dianggap membuktikan business correctness

---

## Business Invariants

Tests harus assert **global invariant**, bukan hanya individual response.

### Inventory

```
final_stock >= 0
successful_sales <= initial_stock
initial_stock == successful_sales + final_stock
```

### Booking

```
count(exclusive_slot) <= 1
```

### Invoice

```
invoice_no unique
```

---

## Common Mistakes

1. ⛔ **Hanya validasi frontend** — Tidak menjaga cross-device/cross-instance
2. ⛔ **Disable button dianggap synchronization** — Bisa di refresh, buka tab baru
3. ⛔ **SELECT → IF → UPDATE** — Stale read + non-atomic write
4. ⛔ **Transaction dianggap otomatis cukup** — Isolation level tidak magic
5. ⛔ **Uniqueness hanya dicek aplikasi** — Race condition masih ada
6. ⛔ **Mutex dianggap distributed** — sync.Mutex hanya dalam-process
7. ⛔ **Redis lock untuk semua masalah** — Overhead + kompleksitas
8. ⛔ **go test -race dianggap coverage** — Race detector ≠ business correctness
9. ⛔ **Arbitrary sleep untuk testing** — Flaky dan tidak deterministic
10. ⛔ **500 concurrent tapi 500 connection** — Connection pool exhaustion

---

## Decision Table

| Problem | Preferred First Tool |
|---------|---------------------|
| Simple decrement stock | Atomic conditional UPDATE |
| Booking uniqueness | UNIQUE constraint |
| Complex read-modify-write + contention | Pessimistic locking |
| Conflict jarang | Optimistic locking |
| Shared memory satu process | Mutex / atomic / channel |
| Coordination lintas instance | Distributed coordination (bila perlu) |

> Tool dipilih berdasarkan **invariant** dan **architecture**, bukan kebiasaan.

---

## Exercise: Senior Software Engineer Daily #5

### Booking Servis Online

Sebuah sistem booking servis di mana customer dapat memilih slot waktu tertentu.

**Flow:**
```
Customer memilih slot 09:00
    ↓
System cek slot kosong
    ↓
System membuat booking
```

**Pertanyaan:**

1. Apa yang terjadi jika dua customer memilih slot yang sama hampir bersamaan?
2. Apakah validasi frontend (cek slot via API) cukup?
3. Apa business invariant-nya?
4. Apakah database constraint diperlukan?
5. Kapan atomic operation (UPDATE WHERE) cukup?
6. Kapan row locking (FOR UPDATE) masuk akal?
7. Kapan optimistic locking lebih sesuai?
8. Bagaimana caranya menguji dengan 500 concurrent request?
9. Apa yang harus di-assert setelah seluruh request selesai?
10. Apakah `go test -race` cukup membuktikan implementation aman?

---

## Commands

```bash
# Deterministic lost update test
go test -v -run TestLostUpdate_Deterministic

# Atomic conditional update tests
go test -v -run TestAtomicUpdate
go test -v -run TestAtomicUpdate_HighContention

# Booking race condition test
go test -v -run Test500_ConcurrentBooking

# PostgreSQL integration tests (requires docker-compose up postgres)
go test -v -tags=integration -run TestPostgres

# Go memory data race (secondary topic)
go test -race -v -run TestUnsafeCounter_Race

# All tests
go test -v -count=1 .

# Stress test (20 iterations)
go test -count=20 .
```

---

## Files

- `inventory.go`: Domain model (Product, InventoryItem, InventoryService)
- `unsafe_inventory.go`: Unsafe check-then-act repository
- `atomic_inventory.go`: Safe atomic conditional update repository
- `lost_update_test.go`: Deterministic lost update reproduction
- `atomic_update_test.go`: Concurrent stress tests (1 dan 500 workers)
- `booking_test.go`: Booking race condition + 500 concurrent test
- `postgres_integration_test.go`: PostgreSQL integration tests
- `datarace.go`: Go memory data race demo (secondary topic)
- `datarace_test.go`: Go race detector tests (secondary topic)

---

## Expected Learning Outcome

Setelah menyelesaikan lab ini, pembaca harus mampu:

1. ✅ Membedakan memory data race dari business race condition
2. ✅ Merekonstruksi race condition secara deterministic
3. ✅ Membuktikan invariant corruption via concurrent execution
4. ✅ Menerapkan atomic conditional update untuk integrity protection
5. ✅ Menjelaskan mengapa BEGIN/COMMIT tidak otomatis melindungi race condition
6. ✅ Memilih tool yang benar berdasarkan invariant dan architecture
7. ✅ Menulis test yang assert global invariant setelah concurrent execution
8. ✅ Membedakan process-local mutex dari distributed coordination

---

## Navigasi

- **Previous**: [Lab 04 — Caching](../04-caching/)
- **Next**: [Lab 06 — Deadlock](../06-deadlock/)
- **All Labs**: [](../)