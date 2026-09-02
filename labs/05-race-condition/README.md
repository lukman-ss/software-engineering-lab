# Lab 05 — Race Condition: Bug yang Tidak Muncul Saat Development, Tapi Menghancurkan Sistem di Production

> **Mental Model**: Kode yang benar ketika dijalankan satu request belum tentu benar ketika banyak request menjalankannya secara bersamaan.

Sistem sering terlihat normal saat dites satu user, tetapi masalah baru muncul ketika ratusan request berjalan bersamaan. Race condition sangat berbahaya karena: intermittent, sulit direproduksi, tidak terlihat saat development, dan dapat merusak *business state* walaupun nilai di database terlihat valid.

---

## 1. Studi Kasus: Stok Sparepart (Oli Mesin)

Bayangkan sebuah toko bengkel memiliki **1 unit** Oli Mesin tersisa.

### Inventory Invariant

```text
initial_stock == successful_sales + final_stock
```

Ini adalah *conservation law* — total barang tidak boleh tiba-tiba bertambah atau hilang.

### Sequential Execution (Benar)

```text
Request A
  READ stock = 1
  CHECK stock > 0? YES
  CALCULATE new_stock = 0
  WRITE stock = 0
  Result: 1 penjualan berhasil, final_stock = 0
  Invariant: 1 == 1 + 0 ✅
```

---

## 2. Timeline Dua Kasir

Apa yang terjadi jika Kasir A dan Kasir B menekan tombol "Bayar" secara hampir bersamaan?

```text
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

**Hasil Akhir**:
- `final_stock = 0`
- `successful_sales = 2`

**Invariant broken**: `1 != 2 + 0` ❌

> **Catatan:** `final_stock = 0` tidak otomatis benar hanya karena tidak ada negative stock. Ini adalah bukti bahwa dua transaksi berhasil, padahal stok hanya satu. Untuk sistem yang **tidak mengizinkan overselling**, invariant yang benar adalah `stock >= 0` dan hanya **satu penjualan** yang boleh berhasil. Hasil `successful_sales = 2` menunjukkan state sudah corrupt — sistem mengira 2 unit terjual padahal stok fisik hanya 1 unit. Toko harus meretur uang satu pembeli.

---

## 3. Kenapa Bisa Terjadi? (Read-Check-Write Pattern)

Setiap kali business logic melibatkan shared mutable state, pola ini sering muncul:

```text
READ
 ↓
CHECK (apakah operation ini valid?)
 ↓
WRITE
```

> **TOCTOU (Time-Of-Check To Time-Of-Use)**: Antara `CHECK` dan `WRITE`, state dapat berubah karena request lain. Jika tidak ada atomicity, hasil `CHECK` tidak lagi berlaku saat `WRITE` dieksekusi.

Contoh: check slot kosong → beberapa milidetik berlalu (state berubah) → insert booking.

### Pertanyaan Kunci

**Apa yang terjadi jika request lain mengubah data setelah READ tapi sebelum WRITE?**

Jika jawabannya "state bisnis akan salah", berarti ada **potential race condition**.

Contoh PHP (hanya ilustrasi):
```php
$product = Product::find($id); // READ

if ($product->stock > 0) {     // CHECK
    $product->stock--;
    $product->save();          // WRITE
}
```

Implementasi yang sama di Go (Unsafe):
```go
func (repo *UnsafeInventory) TrySell(ctx context.Context, productID string) error {
    stock, err := repo.GetStock(ctx, productID) // READ
    if stock <= 0 {                             // CHECK
        return ErrOutOfStock
    }
    return repo.SetStock(ctx, productID, stock-1) // CALCULATE & WRITE
}
```

Kenapa tidak aman? Karena **sequence READ → CHECK → WRITE bukanlah proses yang atomic**.

---

## 4. Lost Update

*Lost update* adalah konsekuensi langsung dari check-then-act yang dijalankan bersamaan.

**Formula:**
`stale read` + `read-modify-write` + `concurrent execution` = `lost update`

1. Request B tidak mengetahui bahwa state sudah diubah (atau sedang diubah) oleh Request A.
2. Request B membaca nilai stale (1).
3. Hasil perhitungan A (0) dioverwrite oleh hasil perhitungan B (0).
4. **Salah satu state transition secara efektif hilang (lost).**

---

## 5. Kenapa Sulit Muncul Saat Development?

- **Timing-dependent**: Hanya terjadi pada *interleaving* spesifik (milidetik).
- **Single-user environment**: Saat dev, tidak ada concurrent user.
- **Network & DB Latency**: Latency di production memperbesar "window of opportunity" antara READ dan WRITE.

---

## 6. Business Invariant

Invariant adalah "kebenaran absolut" yang harus selalu dijaga sistem, tidak peduli seberapa banyak concurrent request yang masuk.

- **Inventory**: `initial == success + final`
- **Booking**: `count(booking_slot) <= 1`
- **Invoice**: `invoice_no is unique`

---

## 7. Solusi 1 — Row Lock (Pessimistic)

Gunakan ketika business decision membutuhkan read data yang stabil sebelum melakukan update.

### SQL Pattern
```sql
BEGIN;

SELECT * FROM products
WHERE id = $1
FOR UPDATE;

-- Aplikasi melakukan perhitungan kompleks di sini
-- UPDATE products SET stock = ...

COMMIT;
```

### Penjelasan Teknis (MVCC & Row Locking)

`SELECT ... FOR UPDATE` mengambil **row-level lock** pada row yang dipilih. Lock ini:

- Membuat transaksi lain yang mencoba memperoleh **conflicting row lock** (misalnya melalui `FOR UPDATE` lain pada row yang sama) **menunggu (blocked)** hingga lock dilepas.
- Membuat `UPDATE` atau `DELETE` yang bersaing pada row yang sama **menunggu**.
- **Tidak** memblokir *ordinary non-locking `SELECT`*: reader tanpa `FOR UPDATE` masih dapat membaca snapshot MVCC row tersebut.
- Lock bertahan sampai **transaction commit atau rollback**.

> **Bukan** `SELECT ... FOR UPDATE` yang "memblokir semua reader", tetapi hanya writer dan locker lain yang konflik. MVCC (Multi-Version Concurrency Control) memungkinkan reader melihat snapshot yang konsisten tanpa perlu lock.

### Flow

1. **Transaction A** mulai, memanggil `FOR UPDATE`, dan mendapat row-level lock.
2. **Transaction B** mulai, memanggil `FOR UPDATE` di row yang sama, lalu **menunggu (blocked)** karena konflik lock.
3. **Transaction A** memvalidasi stok, update ke 0, lalu `COMMIT` (lock dilepas).
4. **Transaction B** unblocked, membaca stok terbaru (0), memvalidasi stok → gagal, lalu `ROLLBACK`.

Lihat: `postgres_rowlock.go`

---

## 8. Solusi 2 — Atomic Conditional Update

Melakukan operasi read, check, dan write di dalam **satu single statement database**.

### SQL Pattern
```sql
UPDATE inventory_products
SET stock = stock - 1
WHERE id = $1
  AND stock > 0
RETURNING stock;
```

### Flow

- **CHECK + UPDATE** dilakukan oleh database engine.
- Jika stok > 0: `1 row affected` (sukses).
- Jika stok 0: `0 rows affected` (gagal, out of stock).

Sangat cepat dan efisien. Lihat: `atomic_inventory.go`

---

## 9. Solusi 3 — Optimistic Locking

Menggunakan kolom `version` untuk mendeteksi perubahan state, bukan me-lock di awal.

```sql
UPDATE products
SET stock = $1, version = version + 1
WHERE id = $2
  AND version = $old_version;
```

Jika `0 rows affected`, artinya record sudah diubah transaksi lain sejak terakhir dibaca. Aplikasi bisa melakukan *retry* atau return error. Cocok untuk environment dengan read yang banyak tapi tingkat konflik/update (contention) rendah.

> **Catatan:** Optimistic locking adalah salah satu strategi concurrency control. Lab ini tidak mengimplementasikannya secara penuh — lihat [`labs/12-optimistic-locking`](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/12-optimistic-locking) untuk production-ready implementation.

---

## 10. Solusi 4 — Distributed Lock

Distributed lock **bukan solusi default** untuk semua race condition. Gunakan hanya ketika coordination benar-benar melintasi:

- **Process** (misalnya multiple replicas di kubernetes)
- **Instance** (bukan hanya thread yang sama)
- **Node** (bukan hanya satu mesin)
- **Resource eksternal** (misalnya API bank, file storage, layanan lain)

Jika invariant sebenarnya berada di **satu database**, maka database lock / atomic update / unique constraint biasanya **lebih sederhana dan lebih andal** dibanding distributed lock.

Contoh: `lock:payment:INV-123`

Skenario:
- Cash Out / Top Up (mencegah double hit ke API Bank)
- Cron Singleton (menjamin hanya 1 worker memproses report)

### Peringatan Teknis

Distributed lock **jauh lebih kompleks** dari database lock. Anda harus memikirkan:
- **Ownership**: Hanya pembuat lock yang boleh unlock.
- **TTL & Expiration**: Lock harus otomatis hilang jika worker mati (crash). Jika TTL terlalu pendek, lock bisa hangus sebelum work selesai. Jika terlalu panjang, sistem tidak pulih dari crash.
- **Safe Unlock**: Evaluasi token kepemilikan (via Lua script) sebelum menghapus lock, agar satu proses tidak meng-unlock lock milik proses lain.
- **Network Partition (Split-Brain)**: Bagaimana jika ada dua master Redis? Lock bisa menjadi tidak konsisten.
- **Lock Acquisition Failure**: Jangan anggap lock selalu berhasil. Handle failure gracefully (retry dengan backoff, atau abort).

> **Rule:** Jika invariant bisa diselesaikan oleh database (Atomic / Row Lock / Unique), jangan pakai Distributed Lock. Kalau Anda memakai Redis untuk mengganti `UNIQUE` constraint database, berarti Anda sedang menambah kompleksitas tanpa mendapatkan integrity yang sama.

---

## 10.1. Distributed Lock vs Database Constraint

| Use Case | Rekomendasi |
|----------|-------------|
| Stock / Booking / Invoice Number (single DB) | **Database Atomic Update / Row Lock / UNIQUE** |
| Koordinasi multi-replica / eksternal API | **Distributed Lock** |
| Mencegah double-click pembayaran | **Idempotency key** (lebih baik dari pada lock) |

---

---

## 11. Kasus Bengkel / CMMS: Invoice Number Race

Alur pembuatan tagihan servis bengkel:
```text
Request
  ↓
Generate Invoice Number
  ↓
Update Stock
  ↓
Save Database
```

### Anti-Pattern: MAX() + 1
```sql
SELECT MAX(invoice_no) FROM invoices;
-- App menghitung: next = max + 1
```

**Race**:
- Request A baca MAX = 100, hitung next = 101.
- Request B baca MAX = 100, hitung next = 101.
- Hasilnya: Dua invoice dengan nomor `INV-000101`!

> **Catatan:** `SELECT MAX(invoice_no) + 1` adalah pattern berbahaya dalam concurrent system. Tidak ada lock, tidak ada atomicity, dan tidak ada jaminan unik.

### Solusi Invoice Number

1. **Database Sequence (PostgreSQL)**
   ```sql
   CREATE SEQUENCE invoice_no_seq;
   SELECT nextval('invoice_no_seq');
   ```
   Fungsi `nextval()` atomik dan pasti unik. (Catatan: sequence *tidak* menjamin gapless/urut tanpa bolong).

2. **Dedicated Counter + Row Lock**
   Jika secara regulasi pajak nomor invoice *harus gapless*, gunakan tabel counter khusus yang di-lock secara pesimistik setiap kali generate nomor baru.

### Unique Constraint Sebagai Safety Net

Meskipun Anda sudah pakai sequence atau counter, **tetap pasang `UNIQUE(invoice_no)`** sebagai final safety net. Database constraint adalah lapisan terakhir yang menolak duplikat — bahkan jika ada bug di aplikasi atau race condition yang tidak terdeteksi.

```sql
CREATE TABLE invoices (
    id SERIAL PRIMARY KEY,
    invoice_no INT NOT NULL,
    customer_id VARCHAR(50) NOT NULL,
    CONSTRAINT uq_invoice_no UNIQUE (invoice_no)
);
```

> **Prinsip:** Sequence/counter menghasilkan angka unik. `UNIQUE` constraint memastikan tidak ada duplikat yang lolos ke database. Keduanya saling melengkapi.

---

## 12. Unique Constraint Sebagai Last Defense

Jika nilai harus unik (seperti `invoice_no` atau `slot_booking`), **jangan hanya mengandalkan pengecekan aplikasi.**

### Invariant Booking

Dalam satu kombinasi `branch_id`, `service_date`, dan `slot_time`, maksimal boleh ada **satu booking** untuk slot eksklusif yang sama.

```text
count(booking WHERE branch_id=X AND service_date=Y AND slot_time=Z) <= 1
```

### TOCTOU Race pada Check-then-Insert

Kode aplikasi yang "cek dulu, baru insert" memiliki *Time-Of-Check to Time-Of-Use (TOCTOU)* race condition:

```text
Customer A → CHECK slot kosong   ✅ (kosong)
Customer B → CHECK slot kosong   ✅ (kosong, belum ada commit A)
Customer A → INSERT booking      ✅
Customer B → INSERT booking      ✅ → 2 booking untuk slot yang sama! ❌
```

Jika hanya mengandalkan `SELECT ... WHERE branch_id=? AND date=? AND slot=?` tanpa database constraint, race ini **tidak dapat dideteksi oleh `go test -race`** dan hanya muncul di production.

### Contoh Booking Service

```sql
CREATE TABLE service_bookings (
    branch_id    VARCHAR(50),
    service_date DATE,
    slot_time    TIME,
    customer_id  VARCHAR(50),
    UNIQUE(branch_id, service_date, slot_time)
);
```

| Layer | Tujuan | Kekuatan |
|-------|--------|----------|
| **Application Check** (`SELECT ...`) | UX — memberitahu user lebih awal | Lemah (mudah dijebol race condition) |
| **Database Constraint** (`UNIQUE`) | Integrity — menjaga Invariant | Absolut (engine database menolak duplikat) |

### Rekomendasi Senior (Prompt 9)

Untuk kasus slot eksklusif sederhana seperti booking, **jangan berpikir Anda harus memilih tepat satu** dari row locking, optimistic locking, atau unique constraint. Rekomendasi senior adalah:

1. **Database `UNIQUE` constraint** — **baseline yang sangat kuat**. Ini adalah final correctness guarantee. Let database engine yang menolak duplikat.
2. **Handle conflict secara proper** — baca error kode `23505` (unique_violation) dan ubah jadi user-friendly response.
3. **Row locking (`SELECT ... FOR UPDATE`)** — gunakan **hanya jika** flow membutuhkan *read-modify-write yang kompleks* (misalnya update stok, generate invoice nomor, atau validasi business rule tambahan sebelum commit).

> **Pokoknya:** Jika invariant bisa diekspresikan langsung sebagai database constraint, gunakan itu. Jangan paksa lock hanya untuk kasus sederhana.

**Konsep response conflict (tanpa fokus HTTP):**
- Result: `conflict` (bukan `error`)
- Pesan: `Slot sudah diambil customer lain. Silakan pilih slot lain.`
- Kode error aplikasi: `ErrDuplicateKey` / `ErrAlreadyBooked`

---

## 13. Race Condition Bukan Hanya Database

Race condition terjadi setiap kali ada **shared resource**:

- **File System**: Dua request mengupload file dengan nama yang sama `profile.jpg` bersamaan. Salah satu file tertimpa.
- **Sistem Eksternal**: Memanggil payment gateway API dua kali karena tombol diklik double.
- **Memory Cache**: Dua goroutine meng-update value map di local memory.

---

## 14. Frontend Validation Tidak Cukup

- "Tombol Submit sudah saya disable setelah diklik."
- Disable button adalah UX. User masih bisa buka tab baru, mematikan JavaScript, curl langsung ke API, atau request terhambat oleh proxy dan dikirim ulang (retry).

---

## 15. Transaction Saja Tidak Cukup

**Transaction ≠ Race-Free.** Membungkus kode dalam `BEGIN` dan `COMMIT` **tidak otomatis** menghilangkan race condition.

```sql
BEGIN;
SELECT stock FROM products WHERE id = 1; -- (membaca 1)
-- (aplikasi cek stock > 0)
UPDATE products SET stock = 0 WHERE id = 1;
COMMIT;
```

Pada default isolation level (`READ COMMITTED`):

- Transaction A membaca `stock = 1`.
- Transaction B membaca `stock = 1` (belum ada commit A).
- A update `stock = 0`, commit.
- B update `stock = 0`, commit → **lost update!**

Correctness bergantung pada tiga hal:

1. **Query** — apakah Anda menggunakan `SELECT` biasa, `SELECT ... FOR UPDATE`, atau `UPDATE ... WHERE stock > 0`?
2. **Lock yang digunakan** — apakah ada row-level lock yang menahan writer lain?
3. **Isolation semantics** — apakah `READ COMMITTED` cukup, atau butuh `REPEATABLE READ` / `SERIALIZABLE`?

> **Intinya:** Transaction menyediakan atomicity dan konsistensi *dari scope transaksi itu sendiri*, bukan dari interleaving dengan transaksi lain. Jika Anda melakukan *check-then-act* (`SELECT` lalu `UPDATE`), maka isolation level default **tidak melindungi** Anda. Gunakan lock atau atomic statement.

---

## 16. Perbedaan: Memory Data Race vs Business Race Condition

**Memory Data Race** vs **Business Race Condition** adalah dua konsep yang berbeda:

### Memory Data Race
- Dua thread mengakses alamat memory yang sama tanpa sinkronisasi, dan setidaknya salah satunya mengubah nilai.
- Dapat dideteksi dengan `go test -race` (Go race detector).
- Contoh: `counter++` dari banyak goroutine tanpa mutex/atomic.

### Business Race Condition
- Terjadi pada shared mutable state di database akibat interleaving timing antara READ → CHECK → WRITE.
- **Tidak dapat dideteksi oleh `go test -race`** karena semua variable di memory sudah thread-safe.
- Contoh: dua transaction membaca `stock = 1`, kemudian keduanya menjual barang yang sama secara bersamaan.

```go
// Kode ini LULUS go test -race (tiap variable thread-safe),
// tapi MASIH punya business race condition (lost update).
func (s *Service) TrySell(ctx context.Context, id string) error {
    stock := s.repo.GetStock(ctx, id)   // thread-safe, tapi bisa stale
    if stock <= 0 {
        return ErrOutOfStock
    }
    return s.repo.SetStock(ctx, id, stock-1) // tidak atomic — lost update!
}
```

> **Penting:** `go test -race` tidak membuktikan kode Anda bebas dari business race condition. Race detector hanya mendeteksi memory-level concurrent access yang tidak sinkron. Jika business logic Anda melibatkan shared mutable state di database, Anda tetap perlu lock, atomic statement, atau constraint — bahkan jika race detector memberi "PASS".

---

## 17. Cara Mengenali Potential Race Condition

Setiap kali Anda menemukan:

1. `READ` → `CHECK` → `WRITE`
2. `CHECK` (slot kosong?) → `CREATE` (booking baru)
3. `READ` (current max_id) → `CALCULATE` (max + 1) → `SAVE`

Langsung tanyakan: **"Apa yang terjadi jika request lain memotong masuk di antara langkah-langkah ini?"**

Domain dengan risiko tinggi: *Stock, Balance, Wallet, Payment, Booking, Reservation, Sequence Number, Invoice Number, Commission, Quota, Seat, Voucher.*

---

## 18. Testing Concurrent System

### Pendekatan yang Salah
⛔ Mengandalkan `time.Sleep(...)` untuk memaksa timing (flaky & tidak deterministik).
⛔ Hanya menjalankan 2 request biasa dan berharap tabrakan.
⛔ Menganggap pass dari `go test -race` sudah cukup aman.

### Pendekatan yang Benar (Seperti di Lab Ini)
✅ **Deterministic Barrier**: Menggunakan *Channel Synchronization* untuk memblokir goroutine tepat sebelum fase WRITE, sehingga eksekusi serentak dipaksakan.
✅ **Stress Test**: Menembak fungsi dengan 500+ goroutine sekaligus (High Contention).
✅ **Invariant Assertion**: Test harus memvalidasi invariant bisnis di akhir (contoh: pastikan `initial_stock == success_count + final_stock`).

> **Catatan Penting pada Test di Lab Ini**:
> Test implementasi "Unsafe" (seperti `TestPostgresUnsafe_LostUpdate`) sengaja dibuat agar **PASS** ketika race condition berhasil direproduksi dan invariant hancur. Ini membuktikan bahwa tanpa perlindungan, sistem tersebut corrupt. Sebaliknya, test implementasi "Safe" dianggap **PASS** jika invariant tetap kokoh berdiri meski diserang 500 concurrent request.

---

## 19. Kesalahan Umum

1. ⛔ **Hanya validasi frontend** — Bypass mudah via API.
2. ⛔ **Tidak memakai Unique Constraint** — Mengandalkan query "SELECT sebelum INSERT".
3. ⛔ **Menganggap BEGIN/COMMIT cukup** — Lupa mengatur isolation atau row lock.
4. ⛔ **Mutex untuk Microservices** — Mutex hanya mengunci satu instance/proses, bukan semua server.

---

## 20. Mindset Senior Hari Ini

"Kode saya berjalan mulus di laptop dan staging (saat dites satu orang)."

Mental model engineer:
> "Apakah sistem tetap stabil, data tetap akurat, dan uang tidak bocor ketika 1.000 user memanggil fungsi ini tepat di milidetik yang sama?"

---

## Running the Lab

### 1. Normal Unit Tests

```bash
# Semua unit tests (in-memory/mock repository)
go test ./...

# Dengan race detector untuk memory data race
go test -race ./...
```

### 2. PostgreSQL Integration Tests

```bash
# Set up PostgreSQL (dari root repository)
docker-compose up -d postgres

# Jalankan integration tests (requires PostgreSQL)
go test -tags=integration ./...
```

#### 2.1 Detailed Integration Test Commands

```bash
# Lost Update detection (Unsafe pattern)
go test -v -tags=integration -run TestPostgresUnsafe_LostUpdate

# Atomic SQL Statement (Safe pattern)
go test -v -tags=integration -run TestPostgresAtomic_ConcurrentUpdate

# Pessimistic Row Lock (Safe pattern)
go test -v -tags=integration -run TestPostgresRowLock_ConcurrentStock
go test -v -tags=integration -run TestPostgresRowLock_HighContention

# PostgreSQL UNIQUE constraint (Booking)
go test -v -tags=integration -run TestPostgres_ConcurrentBooking
go test -v -tags=integration -run TestPostgres_Booking_SameBranchDifferentBranch
go test -v -tags=integration -run TestPostgres_Booking_MultipleBranches
```

#### 2.2 Booking Race Condition Tests

##### In-Memory Booking Concurrency Tests (Mock Repository)
```bash
go test -v -run Test500_ConcurrentBooking
go test -v -run TestBooking_SameBranchDifferentBranch
go test -v -run TestBooking_ErrorHandling
```
> Menggunakan `MockBookingRepository` dengan `sync.RWMutex` - application-level unique constraint demo.

##### PostgreSQL Booking Unique Constraint Tests
```bash
# 500 concurrent: 1 success, 499 SQLSTATE 23505 conflict
go test -v -tags=integration -run TestPostgres_ConcurrentBooking
go test -v -tags=integration -run TestPostgres_Booking_SameBranchDifferentBranch
go test -v -tags=integration -run TestPostgres_Booking_MultipleBranches
```

### 3. Intentional Go Memory Race Demo (Educational)

```bash
# Run intentional race demo (EXPECTED TO FAIL under race detector)
go test -race -v -tags=racedemo -run TestUnsafeCounter_Race

# Normal validation (should pass - unsafe test excluded)
go test -race ./...
```

> **Catatan:** `TestUnsafeCounter_Race` ada di `datarace_unsafe_test.go` dengan build tag `//go:build racedemo`. Ini **diusir** dari normal test suite sehingga `go test -race ./...` selalu PASS. Demo hanya untuk edukasi - jalankan secara eksplisit dengan `-tags=racedemo`.

---

## Testing Best Practices (Lab Notes)

### Concurrency Test Methodology (Prompt 15)

Tests yang benar **tidak boleh** hanya meluncurkan 500 goroutine sekaligus tanpa sinkronisasi:

```go
// ❌ Flaky — scheduler tidak menjamin overlapping yang tinggi
for i := 0; i < 500; i++ {
    go doSomething()
}
```

Pakai **start gate** agar semua worker ready dulu, baru release bersamaan:

```go
// ✅ Deterministic high-contention
ready := make(chan struct{}, n)
release := make(chan struct{})
for i := 0; i < n; i++ {
    go func() {
        ready <- struct{}{} // signal ready
        <-release          // wait for gate open
        doSomething()
    }()
}
// fill gate
for i := 0; i < n; i++ { <-ready }
close(release)
```

### Timeout & Deadlock Prevention (Prompt 16)

Setiap integration test concurrency memakai `context.WithTimeout` atau `select` dengan `time.After`. Jika timeout, test gagal dengan pesan diagnosis spesifik (bukan hang).

### Database Connection Pool (Prompt 17)

- **Application concurrency ≠ Database connection count.**
- 500 goroutine tidak berarti 500 koneksi database.
- `pkg/database` mengkonfigurasi pool: `MaxOpenConns=25`, `MaxIdleConns=5`.
- Workers yang melebihi pool size akan **menunggu** koneksi — ini dihitung sebagai bagian dari concurrency correctness, bukan bug.
- Tujuan lab: menjaga invariant, bukan membebani PostgreSQL dengan connection explosion.

### Test Isolation & Cleanup (Prompt 18)

- Setiap integration test membuat fixture sendiri (productID unik per test).
- `defer` cleanup menghapus data agar test bisa dijalankan berulang tanpa state leak.
- ID produk/booking dirancang secara unik untuk menghindari collision antar test.
- `Test500_ConcurrentBooking` memakai slot yang sama secara disengaja untuk trigger race — ini **bukan** bug test isolation, ini adalah fokus dari test itu sendiri.

### Expected vs Unexpected Errors (Prompt 19)

Untuk unique constraint tests:
- `pq.Error` code `23505` → **expected conflict** (`ErrDuplicateKey`).
- `connection refused`, `timeout`, `syntax error` → **unexpected error**, test harus gagal.

Assertion: `unexpectedErrorCount == 0`.

---

## Running the Lab

### Unit Tests (In-Memory)

```bash
# Semua unit tests (in-memory/mock repository)
go test ./...

# Dengan race detector untuk memory data race
go test -race ./...
```

### PostgreSQL Integration Tests

Prasyarat: PostgreSQL sudah berjalan di `localhost:5432` dengan database `se_lab`.

#### Fresh Database Setup

Schema otomatis di-init via Docker entrypoint pada first database creation:

```bash
# Jalankan Docker Compose (dari root repository)
docker-compose up -d postgres

# Schema auto-created pada fresh database dari:
# ./labs/05-race-condition/schema.sql → /docker-entrypoint-initdb.d/001-race-condition.sql
```

> **Penting:** Init script PostgreSQL hanya berjalan pada **first database creation**. Jika volume lama sudah pernah dibuat sebelum schema mount ditambahkan, lakukan reset volume:

```bash
# Reset fresh test database (HANYA untuk development lab)
docker compose down -v
docker compose up -d postgres

# PERINGATAN: -v menghapus local Docker volume test database.
# JANGAN lakukan ini ke production.
```

#### Run Tests

```bash
# Set environment (opsional, default: postgres:5432, user=postgres, pass=postgres, db=se_lab)
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=postgres
export POSTGRES_DB=se_lab

# Jalankan integration tests (build tag: integration)
go test -v -tags=integration ./...
```

Catatan:
- Test akan `SKIP` bila PostgreSQL tidak tersedia.
- Connection pool: `MaxOpenConns=25`, `MaxIdleConns=5` (dari `pkg/database`).
- `setupTestInventory` membuat test fixture (INSERT/UPDATE), `setupBookingTable` memastikan table ada + TRUNCATE untuk isolation.

---

## Senior Concurrency Mindset Checklist

Sebelum memilih concurrency control strategy, tanyakan diri sendiri:

1. **Invariant apa yang harus selalu benar?** (contoh: `stock >= 0`, `booking_count <= 1`)
2. **Resource apa yang diperebutkan?** (stock, booking slot, invoice number, balance)
3. **Apakah invariant bisa dijaga dengan database constraint?** (UNIQUE, CHECK)
4. **Apakah operasi bisa dibuat atomic?** (single UPDATE/INSERT statement)
5. **Apakah perlu row lock?** (bisa saja — jika ada read-modify-write yang rumit)
6. **Seberapa sering konflik terjadi?** (rendah → optimistic, tinggi → pessimistic)
7. **Apakah koordinasi melintasi database/process/server?** (gunakan distributed lock hanya jika benar-benar perlu)
8. **Bagaimana membuktikannya dengan concurrent test?** (gunakan deterministic barrier + high contention + invariant assertion)

> **Rule:** Mulai dari "invariant apa yang harus dijaga?" bukan "lock apa yang dipakai?". Database constraint + atomic operation sering kali lebih kuat daripada application-level coordination.

---

## Files

- `inventory.go`: Domain model dan Interfaces
- `unsafe_inventory.go`: Implementasi Unsafe check-then-act
- `atomic_inventory.go`: Implementasi Safe Atomic Conditional Update
- `postgres_rowlock.go`: Implementasi Safe Pessimistic Row Locking
- `lost_update_test.go`: Test Barrier Synchronization (Membuktikan Unsafe)
- `atomic_update_test.go`: Test Safe Inventory (1 & 500 concurrent workers)
- `booking_test.go`: Test Uniqueness Constraint (500 concurrent bookings)
- `postgres_rowlock_test.go`: Test Integrasi PostgreSQL Pessimistic Lock
- `postgres_integration_test.go`: Test Integrasi PostgreSQL Unsafe & Atomic
- `postgres_booking.go`: Test Integrasi PostgreSQL Booking (unique constraint)
- `postgres_booking_test.go`: Test Integrasi PostgreSQL Booking (unique constraint)
- `schema.sql`: Skema Tabel Database (multi-branch booking)
- `datarace.go`: UnsafeCounter (demo race), MutexCounter, AtomicCounter, ChannelCounter implementations
- `datarace_test.go`: Safe counter tests (Mutex, Atomic, Channel)
- `datarace_unsafe_test.go`: Build tag `racedemo` - intentional Go data race demo (EXPECTED TO FAIL under race detector)

---

## Navigasi

- **Previous**: [Lab 04 — Caching](../04-caching/)
- **Next**: [Lab 06 — API Versioning](../06-api-versioning/)
