# Lab 05 — Race Condition: Bug yang Tidak Muncul Saat Development, Tapi Menghancurkan Sistem di Production

> **Mental Model**: Kode yang benar ketika dijalankan satu request belum tentu benar ketika banyak request menjalankannya secara bersamaan.

---

## Sebelumnya

Di Lab 01–04, pembaca sudah mempelajari:

- **Lab 01 — Idempotency**: Request dapat di-retry tanpa efek samping berulang
- **Lab 02 — Database Index**: Optimasi query, bukan concurrency protection
- **Lab 03 — Database Transaction**: ACID guarantee, tapi isolation level terbatas
- **Lab 04 — Caching**: Consistent cache invalidation, stale read problem

## Masuk ke Race Condition

Sistem terlihat normal saat dites satu user. Problem baru muncul ketika ratusan atau ribuan request berjalan bersamaan.

Race condition seringkali:

- **Intermittent** — tidak selalu terjadi
- **Sulit direproduksi** — butuh timing spesifik
- **Tidak terlihat di development** — biasanya single-user / low-load
- **Baru terlihat di production** — high concurrency, network latency
- **Merusak business state** — meski database value terlihat valid

Jika concurrent execution dapat membroken **business invariant**, itu adalah **application-level race condition** — bukan hanya Go memory data race.

---

## Studi Kasus Utama: Oli Mesin — Lost Update

**Setup**: Toko memiliki **1 unit** Oli Mesin. Stok = 1.

### Inventory Invariant

```
initial_stock = successful_sales + final_stock
```

Ini adalah **conservation law** — total unit tidak boleh bertambah atau berkurang.

### Sequential (benar)

```
Request A
  READ stock = 1
  CHECK stock > 0? YES
  CALCULATE new_stock = 0
  WRITE stock = 0
  Result: 1 penjualan berhasil, final_stock = 0
  Invariant: 1 == 1 + 0 ✅
```

### Concurrent (runtuh)

```
Kasir A                  Kasir B
   │                        │
T0 │ READ stock = 1 ◄──────┼── READ stock = 1
   │                        │
T1 │ CHECK > 0? YES        │ CHECK > 0? YES
   │                        │
T2 │ CALCULATE new = 0     │ CALCULATE new = 0
   │                        │
T3 │ WRITE stock = 0 ◄─────┼── WRITE stock = 0
   │                        │
```

**Hasil**:

```
final_stock = 0
successful_sales = 2

1 != 2 + 0  →  Invariant broken ❌
```

**Final stock 0 terlihat valid.** Database tidak error. Tapi state sudah corrupt — sistem mengira 2 unit terjual padahal stok hanya 1.

### Timeline Detail Lost Update

```
T0  Request A  READ stock = 1
T1  Request B  READ stock = 1        ← stale read (A belum WRITE)
T2  Request A  CALCULATE newStock = 0
T3  Request B  CALCULATE newStock = 0 ← hitung dari nilai stale yang sama
T4  Request A  WRITE stock = 0
T5  Request B  WRITE stock = 0       ← overwrite hasil A, tidak ada delta
```

### Kenapa Database Tidak Melindungi?

Pada isolation level default (`READ COMMITTED`), setiap statement mendapat snapshot terbaru. Tapi sequence READ → CHECK → CALCULATE → WRITE adalah **beberapa statement terpisah**. Database tidak tahu ini adalah satu logical operation.

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

## Read-Check-Write Pattern (Unsafe)

Setiap kali business logic melibatkan shared mutable state, pola berikut hampir selalu muncul:

```
READ
 ↓
CHECK (apakah operation ini valid?)
 ↓
WRITE
```

### Ilustrasi PHP

```php
$product = Product::find($id);

if ($product->stock > 0) {
    $product->stock--;
    $product->save();
}
```

Ini sama persis dengan trySell di Go — hanya berbeda bahasa.

### Pertanyaan Kunci

**Apa yang terjadi jika request lain mengubah data setelah READ tapi sebelum WRITE?**

Jika jawabannya dapat menyebabkan **business state salah**, berarti ada **potential race condition**.

#### Contoh: Stock Oli Mesin

```
Request A: READ stock = 1
Request B: READ stock = 1     ← stale, A belum WRITE
Request A: CHECK > 0 → YES
Request A: WRITE stock = 0
Request B: CHECK > 0 → YES    ← lewat karena baca stale
Request B: WRITE stock = 0    ← overwrite, unit "melayang"
```

Result: 2 sales, 1 unit berhak. **Invariant rusak.**

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

### Timeline Detail

```
T0  Request A membaca stock = 1
T1  Request B membaca stock = 1 (stale)
T2  Request A menghitung newStock = 0
T3  Request B menghitung newStock = 0
T4  Request A menulis stock = 0
T5  Request B menulis stock = 0
```

**Penjelasan teknis**:

1. Request B tidak mengetahui bahwa state sedang/sudah diubah oleh Request A.
2. Request B membaca nilai stale (1).
3. Request A dan B menghasilkan state yang sama (0).
4. Hasil `WRITE` dari Request A dioverwrite oleh Request B tanpa memperhitungkan mutasi A.
5. **Salah satu state transition secara efektif hilang (lost update).**

**Formula**:

```
stale read + read-modify-write + concurrent execution = lost update
```

### Hubungan Sebab-Akibat

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
- Request A dan B wajib BLOCK sampai keduanya selesai READ (`aReadDone`, `bReadDone`)
- Baru lanjut ke CALCULATE → WRITE (`aCalcDone`, `bCalcDone`)
- WRITE juga disinkronisasi (`aWriteDone`)

```go
aReadDone := make(chan struct{})
bReadDone := make(chan struct{})
// A reads → close(aReadDone) → B reads → close(bReadDone)
// A calculates → close(aCalcDone) → B calculates → close(bCalcDone)
// A writes → close(aWriteDone) → B writes (overwrite)
```

**Hasil deterministic**: `successful_sales = 2, final_stock = 0` — **every time, every platform, every run**.

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

## Solution: Row Locking (Pessimistic)

Row locking memperoleh exclusive lock sebelum membaca, sehingga tidak ada goroutine lain yang dapat membaca state tersebut hingga transaksi selesai.

### SQL Pattern

```sql
BEGIN;

SELECT * FROM products
WHERE id = $1
FOR UPDATE;

-- Aplikasi melakukan CHECK di sini
-- UPDATE stock = new_value

COMMIT;
```

### Transaction Flow

```
Transaction A                     Transaction B
    │                                │
    ▼ BEGIN                          ▼ BEGIN
    ▼ SELECT ... FOR UPDATE         │  (blocked — A holds lock)
    ▼ CHECK stock > 0 → YES         │
    ▼ UPDATE stock = 0              │
    ▼ COMMIT                        ▼ SELECT ... FOR UPDATE (now unblocks)
                                   ▼ CHECK stock > 0 → NO
                                   ▼ ROLLBACK / return error
```

### Expected Behavior

- **Transaction A** memperoleh lock row pertama.
- **Transaction B** menunggu sampai A COMMIT/ROLLBACK.
- Setelah A commit, B membaca nilai terbaru (stock = 0) → CHECK gagal → REJECT.

### PostgreSQL Implementation

```go
func (r *PostgresRowLockRepository) TrySell(ctx context.Context, productID string) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Lock row — blocks other FOR UPDATE until commit
    var stock int
    err = tx.QueryRowContext(ctx,
        "SELECT stock FROM products WHERE id = $1 FOR UPDATE",
        productID).Scan(&stock)
    if err != nil {
        return err
    }

    if stock <= 0 {
        return ErrOutOfStock
    }

    _, err = tx.ExecContext(ctx,
        "UPDATE products SET stock = $1 WHERE id = $2",
        stock-1, productID)
    if err != nil {
        return err
    }

    return tx.Commit()
}
```

> Implementasi penuh di [Lab 11 — Pessimistic Locking](../11-pessimistic-locking/)

---

## Solution: Atomic Conditional Update

Atomic conditional UPDATE melakukan CHECK + WRITE sekaligus dalam satu statement SQL.

### SQL Pattern

```sql
UPDATE products
SET stock = stock - 1
WHERE id = $1
  AND stock > 0
RETURNING stock;
```

### How It Works

- **CHECK + UPDATE dilakukan oleh database** sebagai satu statement yang atomic.
- **1 row affected** → decrement sukses.
- **0 rows affected** → stock habis (condition `stock > 0` gagal) → REJECT.

### Flow

```
Request A: UPDATE ... WHERE stock > 0 → 1 row affected → success
Request B: UPDATE ... WHERE stock > 0 → 0 rows affected → rejected
```

### Keuntungan

- Tidak ada SELECT terlebih dahulu → tidak ada stale read window.
- Read-modify-write atomik di database.
- Concurrent execution di-handle oleh database engine.
- Lebih ringan daripada row lock (tidak perlu hold lock selama aplikasi memproses).

`atomic_inventory.go` implements pola ini.

---

## Solution: Optimistic Locking

Optimistic locking memakai **version column** untuk mendeteksi conflict, bukan mencegahnya.

### Schema

```sql
CREATE TABLE products(
    id      VARCHAR(50) PRIMARY KEY,
    stock   INT NOT NULL,
    version INT NOT NULL DEFAULT 1
);
```

### SQL Pattern

```sql
-- App reads version = 10
UPDATE products
SET stock = $newStock,
    version = version + 1
WHERE id = $id
  AND version = $oldVersion;   -- 10
```

### Flow

```
read version = 10
calculate new_stock
conditional update WHERE version = 10
  ↓
  1 row affected → commit (version increment to 11)
  0 rows affected → conflict detected (someone else changed it)
```

### Conflict Handling

Jika `0 rows affected`, artinya state berubah sejak dibaca. Aplikasi harus:

1. **Reload** data terbaru
2. **Retry** operasi
3. **Atau reject** dengan pesan conflict

Cocok untuk high concurrency dengan **low contention**. Untuk inventory hot-stock (stock = 1, banyak yang compete), atomic conditional UPDATE biasanya lebih baik.

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

## Invoice Workflow Race Case

Salah satu skenario nyata:

```
Request
    ↓
Generate Invoice Number
    ↓
Update Stock
    ↓
Generate Komisi
    ↓
Save to Database
```

### Race Example: Invoice Number

```
Request A → SELECT MAX(inv_no) = 100
Request B → SELECT MAX(inv_no) = 100
Request A → App calculates INV-000101
Request B → App calculates INV-000101  ← DUPLICATE!
```

**Expected:**
- INV-000100
- INV-000101

**Database harus menjaga uniqueness** melalui UNIQUE constraint, bukan hanya application check.

### Solution Landscape

1. **Database Sequence** (paling direkomendasikan)
   ```sql
   CREATE SEQUENCE invoice_no_seq START 100;
   SELECT nextval('invoice_no_seq') FROM invoices;
   ```

2. **Unique Constraint + Retry**
   ```sql
   ALTER TABLE invoices ADD CONSTRAINT uq_invoice_no UNIQUE (invoice_no);
   -- Jika conflict, retry
   ```

---

## Race Condition Bukan Hanya Database

Race condition juga terjadi di luar database — file system, cache, dan sistem eksternal.

### Contoh: File Upload

```
POST /upload → file = "profile.jpg"
POST /upload → file = "profile.jpg"  ← teroverride!
```

Request A: upload "profile.jpg"
Request B: upload "profile.jpg" (status 200, tapi file A tertimpa)

**Solusi**: Generate unique filename di sisi server (UUID, timestamp).

### Contoh: Cron Singleton

```
Server A: generate_report()
Server B: generate_report()  ← duplicate!
```

Dua server sekaligus menjalankan cron yang sama → duplicate processing, notifikasi ganda.

**Solusi**: Distributed lock (Redis, PostgreSQL advisory lock) untuk memastikan satu job hanya dapat dijalankan satu instance.

### Contoh Lain-lain

- **Quota** — Rate limiting per user harus konsisten di semua instance
- **Seat Reservation** — Online ticket booking: satu kursi hanya untuk satu pembeli
- **Balance** — Wallet topup/cashout: race dapat bikin duplikat transfer
- **Commission** — Hitung komisi harus akurat, bukan double-count
- **Sequence Number** — Invoice no, order no harus unik

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

## Solution: Distributed Lock

Jika koordinasi melibatkan *system eksternal* atau *multiple services* di luar satu database, pertimbangkan distributed lock.

### Contoh Kasus (Conceptual Redis)

`lock:invoice:10291`

Use case yang umum:
- Payment / Top Up / Cash Out (menghindari double charge di API eksternal)
- Pembayaran Vendor
- Cron singleton (memastikan hanya 1 server yang menjalankan task)
- Report generation (menghindari duplikasi heavy processing)

### Peringatan Teknis

> **Distributed lock jauh lebih kompleks daripada database lock.** Jangan berasumsi Redis lock selalu menjadi solusi terbaik.

Minimal harus memikirkan:
1. **Ownership**: Hanya pemilik lock yang boleh unlock.
2. **TTL (Time-To-Live)**: Bagaimana jika proses crash sebelum memanggil unlock? Lock harus memiliki auto-expiration.
3. **Expiration & Stale Lock**: Bagaimana jika proses melambat dan lock expired, lalu proses lain mengambil lock, dan proses pertama bangun lalu menulis state?
4. **Safe Unlock**: Memerlukan Lua script untuk mengecek kepemilikan sebelum `DEL`.
5. **Network Partition**: Split-brain bisa menyebabkan >1 pemegang lock.

Jika invariant bisa ditegakkan via database (atomic/constraint), gunakan mekanisme database.

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

## Cara Mengenali Potential Race Condition

Setiap kali menemukan pola berikut, tanyakan: **"Apa yang terjadi jika request lain menjalankan operasi yang sama di tengah proses?"**

### Pattern 1: Read-Check-Write

```
READ
 ↓
CHECK (apakah operation valid?)
 ↓
WRITE
```

**Contoh**: Cek stok → kurangi stok.

### Pattern 2: Check-Create

```
CHECK (apakah slot kosong?)
 ↓
CREATE (booking baru)
```

**Contoh**: Cek slot booking → buat booking.

### Pattern 3: Read-Calculate-Save

```
READ current value
 ↓
CALCULATE next value
 ↓
SAVE
```

**Contoh**: Baca invoice number terakhir → hitung next → save.

### High-Risk Domain

Jika kode kamu berurusan dengan salah satu berikut, **langsung anggap sebagai race candidate**:

- **Stock** / **Inventory** — sparepart, oli, barang fisik
- **Balance** / **Wallet** — saldo uang, e-wallet
- **Payment** — pembayaran, top up, cash out
- **Booking** / **Reservation** — slot waktu, kursi, kamar
- **Sequence Number** — invoice no, order no
- **Commission** — perhitungan komisi
- **Quota** — rate limit, kuota pengguna
- **Voucher** — nomor voucher, diskon
- **Inventory Allocation** — alokasi stok ke order

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

## Unique Constraint Sebagai Last Defense

Jika suatu business value **harus unik**, **jangan hanya mengandalkan application-level check.**

### Contoh

- `invoice_no` → harus unik
- `booking slot` → hanya satu booking per slot
- `external reference` → tidak boleh duplikat yang sama

### Contoh Booking Table

```sql
CREATE TABLE service_bookings(
    branch_id    VARCHAR(50),
    service_date DATE,
    slot_time    TIME,
    customer_id  VARCHAR(50),
    UNIQUE(branch_id, service_date, slot_time)  -- Composite unique constraint
);
```

### Konsep: Application Check ≠ Database Constraint

| Layer | Purpose | Protection |
|-------|---------|------------|
| Application check | UX protection — beri tahu user cepat | Tidak kuat — bisa dipass lewat tab baru, refresh, atau concurrent request |
| Database constraint | Integrity protection — garanti tidak mungkin | Kuat — sudah dipaksakan oleh database engine |

**Never rely on application check alone for invariants.**

---

## Common Mistakes

Kesalahan-kesalahan umum yang harus dihindari:

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

# PostgreSQL Row Lock tests (requires docker-compose up postgres)
go test -v -tags=integration -run TestPostgresRowLock

# PostgreSQL Atomic Update tests (requires docker-compose up postgres)
go test -v -tags=integration -run TestPostgresAtomic

# PostgreSQL Unsafe Lost Update test (requires docker-compose up postgres)
go test -v -tags=integration -run TestPostgresUnsafe

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
- `lost_update_test.go`: Deterministic lost update reproduction (barrier synchronization)
- `atomic_update_test.go`: Concurrent stress tests (1 dan 500 workers)
- `booking_test.go`: Booking race condition + 500 concurrent test
- `postgres_rowlock.go`: Pessimistic row locking implementation
- `postgres_rowlock_test.go`: Row lock integration tests (2 req, 500 req scenarios)
- `postgres_integration_test.go`: PostgreSQL unsafe + atomic integration tests
- `datarace.go`: Go memory data race demo (secondary topic)
- `datarace_test.go`: Go race detector tests (secondary topic)
- `schema.sql`: Database schema for inventory_products, service_bookings, invoices

---

## Expected Learning Outcome

Setelah menyelesaikan lab ini, pembaca harus mampu:

1. ✅ Membedakan memory data race dari business race condition
2. ✅ Merekonstruksi race condition secara deterministic (barrier, bukan sleep)
3. ✅ Membuktikan invariant corruption via concurrent execution
4. ✅ Menjelaskan pola Read-Check-Write dan mengapa itu unsafe
5. ✅ Menerapkan atomic conditional UPDATE untuk integrity protection
6. ✅ Menerapkan row locking (FOR UPDATE) untuk complex read-modify-write
7. ✅ Menjelaskan mengapa BEGIN/COMMIT tidak otomatis melindungi race condition
8. ✅ Memilih tool yang benar berdasarkan invariant dan architecture
9. ✅ Menulis test yang assert global invariant setelah concurrent execution
10. ✅ Membedakan process-local mutex dari distributed coordination

---

## Navigasi

- **Previous**: [Lab 04 — Caching](../04-caching/)
- **Next**: [Lab 06 — Deadlock](../06-deadlock/)
- **All Labs**: [](../)