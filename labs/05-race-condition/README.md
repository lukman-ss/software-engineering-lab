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

Final stock = 0 terlihat *valid* di database. Tidak ada negative stock. Tapi **state sudah corrupt** — sistem mengira 2 unit terjual padahal stok fisik hanya 1 unit. Toko harus meretur uang satu pembeli.

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

### Penjelasan Teknis (MVCC)

`SELECT ... FOR UPDATE` mengambil **row-level lock** yang mencegah transaksi lain melakukan *conflicting lock/update* terhadap row yang sama sampai transaction selesai.

> *Catatan MVCC: Ordinary non-locking `SELECT` masih dapat membaca snapshot row tersebut; row lock hanya memblokir writer atau locker lain (seperti transaksi lain yang juga memanggil `FOR UPDATE`).*

### Flow

1. **Transaction A** mulai, memanggil `FOR UPDATE`, dan mendapat lock.
2. **Transaction B** mulai, memanggil `FOR UPDATE` di row yang sama, lalu **menunggu (blocked)**.
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

Menggunakan versi (`version`) untuk mendeteksi perubahan state, bukan me-lock di awal.

```sql
UPDATE products
SET stock = $1, version = version + 1
WHERE id = $2
  AND version = $old_version;
```

Jika `0 rows affected`, artinya record sudah diubah transaksi lain sejak terakhir dibaca. Aplikasi bisa melakukan *retry* atau return error. Cocok untuk environment dengan read yang banyak tapi tingkat konflik/update (contention) rendah.

---

## 10. Solusi 4 — Distributed Lock

Jika koordinasi terjadi **di luar scope satu database** (misal memanggil API eksternal), gunakan Distributed Lock (seperti Redis).

Contoh: `lock:payment:INV-123`

Skenario:
- Cash Out / Top Up (mencegah double hit ke API Bank)
- Cron Singleton (menjamin hanya 1 worker memproses report)

### Peringatan Teknis

Distributed lock **jauh lebih kompleks** dari database lock. Anda harus memikirkan:
- **Ownership**: Hanya pembuat lock yang boleh unlock.
- **TTL & Expiration**: Lock harus otomatis hilang jika worker mati (crash).
- **Safe Unlock**: Evaluasi token kepemilikan (via Lua script) sebelum menghapus lock.
- **Network Partition (Split-Brain)**: Bagaimana jika ada dua master?

*Rule: Jika invariant bisa diselesaikan oleh database (Atomic / Row Lock / Unique), jangan pakai Distributed Lock.*

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

### Solusi Invoice Number

1. **Database Sequence**
   ```sql
   SELECT nextval('invoice_no_seq');
   ```
   Fungsi `nextval()` atomik dan pasti unik. (Catatan: sequence *tidak* menjamin gapless/urut tanpa bolong).

2. **Dedicated Counter + Row Lock**
   Jika secara regulasi pajak nomor invoice *harus gapless*, gunakan tabel counter khusus yang di-lock secara pesimistik setiap kali generate nomor baru.

---

## 12. Unique Constraint Sebagai Last Defense

Jika nilai harus unik (seperti `invoice_no` atau `slot_booking`), **jangan hanya mengandalkan pengecekan aplikasi.**

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

## 15. Transaction Tidak Otomatis Menyelesaikan Semua Masalah

**Transaction ≠ Automatic Serialization.**

```sql
BEGIN;
SELECT stock FROM products WHERE id = 1; -- (membaca 1)
-- (aplikasi cek stock > 0)
UPDATE products SET stock = 0 WHERE id = 1;
COMMIT;
```

Pada default isolation level (`READ COMMITTED`), dua transaksi tetap bisa membaca nilai awal yang sama dan menghasilkan *Lost Update*. Membungkus kode dalam `BEGIN` dan `COMMIT` tidak membebaskan Anda dari keharusan merancang strategi lock/atomic yang tepat.

---

## 16. Memory Data Race vs Business Race Condition

- **Memory Data Race**: Dua thread mengakses alamat memory yang sama tanpa sinkronisasi, dan setidaknya salah satunya mengubah nilai. (Dapat dideteksi dengan `go test -race`).
- **Business Race Condition**: Alur logika berantakan akibat *interleaving* timing, padahal semua variable di memory sudah *thread-safe* (menggunakan channel atau local variable).

> **`go test -race` tidak membuktikan kode Anda bebas dari business race condition.**

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

## Commands

### Deterministic Lost Update Reproduction (Barrier Synchronization)
Membuktikan bahwa Unsafe/Check-Then-Act melanggar invariant.
```bash
go test -v -run TestLostUpdate_Deterministic
```

### Atomic Conditional Update Tests
Membuktikan bahwa atomic decrement melindungi invariant.
```bash
go test -v -run TestAtomicUpdate_StockOne
go test -v -run TestAtomicUpdate_HighContention
```

### Booking Race Condition Test (Unique Constraint)
```bash
go test -v -run Test500_ConcurrentBooking
```

### PostgreSQL Integration Tests (Requires `docker-compose up postgres`)

```bash
# Membuktikan Unsafe Pattern di DB menghasilkan Lost Update
go test -v -tags=integration -run TestPostgresUnsafe_LostUpdate

# Membuktikan Atomic SQL Statement menjaga Invariant
go test -v -tags=integration -run TestPostgresAtomic_ConcurrentUpdate

# Membuktikan Pessimistic Row Lock menjaga Invariant
go test -v -tags=integration -run TestPostgresRowLock_ConcurrentStock
go test -v -tags=integration -run TestPostgresRowLock_HighContention
```

### Secondary Topic: Go Memory Data Race
```bash
go test -race -v -run TestUnsafeCounter_Race
```

### All Tests & Stress Testing
```bash
go test -v -count=1 .
go test -count=20 .   # Menjalankan semua test 20 kali (stress test)
```

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
- `schema.sql`: Skema Tabel Database
- `datarace.go` & `datarace_test.go`: Secondary topic (Go memory model)
