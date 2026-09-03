# Lab 08 — Database Isolation Level

## Problem
Banyak engineer berasumsi bahwa membungkus sederet query SQL di dalam `BEGIN TRANSACTION` dan `COMMIT` secara otomatis membuat aplikasi bebas dari race condition dan data anomaly. Di production dengan ratusan request per detik, asumsi ini sering kali menyebabkan saldo bocor, overbooking, double spending, hingga laporan keuangan yang tidak konsisten.

`BEGIN TRANSACTION` hanya menjamin **Atomicity** (All-or-Nothing) dari scope transaksi tersebut, bukan **Isolation** dari query yang dieksekusi bersamaan oleh transaksi lain.

## Apa itu Isolation
Isolation adalah komponen "I" dalam ACID yang menentukan sejauh mana perubahan yang dilakukan oleh satu transaksi dapat dilihat atau diinterferensi oleh transaksi lain yang berjalan secara bersamaan.

## ACID
- **Atomicity**: Seluruh operasi dalam transaksi sukses bersama, atau jika gagal seluruhnya di-rollback.
- **Consistency**: Transaksi membawa database dari satu state yang valid ke state valid lainnya (menjaga seluruh database constraints).
- **Isolation**: Menentukan visibilitas state data antar transaksi concurrent.
- **Durability**: Sekali transaksi di-commit, perubahannya permanen dan bertahan meskipun terjadi crash/mati listrik.

## Dirty Read
Dirty Read terjadi saat Transaksi A membaca baris data yang sedang diubah oleh Transaksi B, padahal Transaksi B **belum di-commit**. Jika Transaksi B kemudian melakukan `ROLLBACK`, maka data yang dibaca oleh Transaksi A adalah data "fiktif" (dirty) yang sebenarnya tidak pernah sah ada di database.

## Non-Repeatable Read
Non-Repeatable Read (atau Fuzzy Read) terjadi saat Transaksi A membaca baris data yang sama dua kali dalam satu transaksi, namun mendapatkan nilai yang berbeda karena di antara dua pembacaan tersebut, Transaksi B melakukan `UPDATE` atau `DELETE` pada baris tersebut dan melakukan `COMMIT`.

## Phantom Read
Phantom Read terjadi pada query berbasis range (misalnya `SELECT COUNT(*) WHERE status = 'PAID'`). Transaksi A menjalankan query range pertama, kemudian Transaksi B melakukan `INSERT` baris baru yang cocok dengan kriteria tersebut dan me-`COMMIT`-nya. Ketika Transaksi A menjalankan query range yang sama untuk kedua kalinya, muncul baris data baru ("hantu" / phantom) yang sebelumnya tidak ada.

## Isolation Levels

Standar ANSI/ISO SQL-92 mendefinisikan 4 tingkat isolasi berdasarkan anomali yang dicegah:

| Isolation Level | Dirty Read | Non-Repeatable Read | Phantom Read |
|---|---|---|---|
| **READ UNCOMMITTED** | Diizinkan | Diizinkan | Diizinkan |
| **READ COMMITTED** | Dicegah | Diizinkan | Diizinkan |
| **REPEATABLE READ** | Dicegah | Dicegah | Diizinkan (Standard ANSI) |
| **SERIALIZABLE** | Dicegah | Dicegah | Dicegah |

## PostgreSQL dan MVCC

PostgreSQL mengimplementasikan transaksi menggunakan **Multi-Version Concurrency Control (MVCC)**. Dalam MVCC, setiap penulisan data menghasilkan versi tuple baru (`xmin`, `xmax`) tanpa menimpa data lama secara langsung, sehingga pembaca tidak pernah memblokir penulis dan penulis tidak pernah memblokir pembaca.

Perbedaan penting implementasi PostgreSQL dibanding database lain:
- **READ UNCOMMITTED**: PostgreSQL **tidak mendukung Dirty Read**. Jika Anda menjalankan `SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;`, PostgreSQL secara internal memperlakukannya sama persis seperti `READ COMMITTED`.
- **REPEATABLE READ**: PostgreSQL mengambil snapshot MVCC **sekali saja pada awal statement pertama dalam transaksi**. Karena snapshot ini statis sepanjang transaksi, PostgreSQL REPEATABLE READ **secara otomatis mencegah Phantom Read klasik**. Jangan berasumsi Repeatable Read selalu memunculkan Phantom Read di semua database engine!
- **SERIALIZABLE**: PostgreSQL menggunakan **SSI (Serializable Snapshot Isolation)**. SSI tidak memblokir query atau me-lock seluruh baris seperti database lawas, melainkan mencatat dependensi pembacaan (SIREAD locks) di memori dan membatalkan (abort) transaksi yang menimbulkan dependensi anomali (Write Skew).

## Read Committed
Level default di PostgreSQL. Setiap query SQL dalam transaksi mengambil **snapshot baru** dari data yang telah di-commit pada saat query tersebut dieksekusi.
- Anomali yang mungkin muncul: *Non-Repeatable Read*, *Phantom Read*, dan *Lost Update* (pada pola check-then-act).

## Repeatable Read
Snapshot database diambil hanya satu kali saat transaksi pertama kali membaca data.
- Menjamin pembacaan berulang selalu konsisten.
- Jika transaksi mencoba meng-`UPDATE` baris yang telah di-update & di-commit oleh transaksi lain setelah snapshot diambil, PostgreSQL akan menggagalkan transaksi dengan error code `SQLSTATE 40001 (serialization_failure)`.

## Serializable
Tingkat isolasi tertinggi. **Serializable menjamin hasil concurrent transaction ekuivalen dengan suatu urutan eksekusi serial.** Ini bukan berarti database menjalankan transaksi satu per satu secara harfiah. Transaksi tetap berjalan paralel, namun engine SSI akan memvalidasi apakah ada siklus anomali (misalnya Write Skew). Jika terjadi konflik, engine me-return `SQLSTATE 40001`.

## SELECT FOR UPDATE
Mekanisme Pessimistic Locking eksplisit. Mengambil row-level exclusive lock pada baris terpilih:
- Transaksi lain yang mencoba mengambil lock `FOR UPDATE` atau meng-`UPDATE` baris yang sama akan **tertahan (blocked)** sampai transaksi pertama commit atau rollback.
- Pembaca biasa (`SELECT` tanpa lock) **tetap tidak terblokir** berkat snapshot MVCC.

**Deadlock Prevention**:
Jika Transaksi 1 mentransfer `A -> B` (lock A lalu B) dan Transaksi 2 mentransfer `B -> A` (lock B lalu A) secara bersamaan, database akan mengalami deadlock (`SQLSTATE 40P01`).
Solusi wajib: Gunakan **Deterministic Lock Ordering** — selalu urutkan lock dari ID akun terkecil ke terbesar.

```go
firstID, secondID := fromID, toID
if firstID > secondID {
    firstID, secondID = secondID, firstID
}
tx.QueryRowContext(ctx, "SELECT balance FROM accounts WHERE id = $1 FOR UPDATE", firstID)
tx.QueryRowContext(ctx, "SELECT balance FROM accounts WHERE id = $2 FOR UPDATE", secondID)
```

## Lost Update
Terjadi saat aplikasi melakukan pola **Check-Then-Act**:
1. Transaksi 1 membaca saldo Alice = 1.000.000.
2. Transaksi 2 membaca saldo Alice = 1.000.000 (belum ada commit).
3. Transaksi 1 mengurangi 800.000, menulis saldo Alice = 200.000, lalu commit.
4. Transaksi 2 mengurangi 800.000, menulis saldo Alice = 200.000, lalu commit.
Hasil: Alice mentransfer total 1.600.000 padahal saldonya hanya 1.000.000! Pengurangan pertama ditimpa dan hilang (lost).

## Serialization Failure
Ketika menggunakan `REPEATABLE READ` atau `SERIALIZABLE`, PostgreSQL menolak eksekusi concurrent yang berbenturan dengan mengembalikan error:
```text
ERROR: could not serialize access due to concurrent update (SQLSTATE 40001)
```
Aplikasi yang menggunakan isolation level tinggi **wajib menyiapkan mekanisme retry terarah** (bounded retries dengan exponential backoff dan jitter).

## Wallet Transfer Case
Studi kasus sistem dompet digital dengan 3 akun:
```text
Alice   : Rp 1.000.000
Bob     : Rp 1.000.000
Charlie : Rp 1.000.000
```
- **Invariant 1**: `balance >= 0` (Saldo tidak boleh negatif).
- **Invariant 2**: `total_money == 3.000.000` (Konservasi total uang dalam ekosistem).

## Experiments

Lab ini menyajikan eksperimen executable:
1. `read_committed_test.go`: Membuktikan Non-Repeatable Read di bawah default PostgreSQL.
2. `repeatable_read_test.go`: Membuktikan Snapshot Isolation dan pendeteksian concurrent update conflict (`40001`).
3. `phantom_read_test.go`: Membuktikan Phantom Read di Read Committed dan pencegahannya secara otomatis di Repeatable Read PostgreSQL.
4. `transfer_test.go`:
   - Reproduksi Lost Update / Uang gaib pada implementasi Naive.
   - Pembuktian invariant kokoh pada `SELECT ... FOR UPDATE` dengan Deterministic Lock Order.
   - Pembuktian Deadlock-free pada transfer bolak-balik (`A->B` dan `B->A`).
   - Stress test 100 concurrent transfer.
5. `serializable_test.go`: Menguji penanganan `40001` via bounded retry + backoff & jitter.

## Running the Lab

Prasyarat: PostgreSQL berjalan di `localhost:5432` (`make infra-up` atau `docker compose up -d postgres`).

```bash
# Masuk ke direktori lab
cd labs/08-database-isolation-level

# Jalankan seluruh test
go test -v ./...
```

## Running Tests

```bash
# Non-Repeatable Read
go test -v -run TestReadCommitted_NonRepeatableRead ./...

# Snapshot Isolation
go test -v -run TestRepeatableRead_SnapshotIsolation ./...

# Phantom Read
go test -v -run TestPhantomRead_ReadCommitted ./...
go test -v -run TestPhantomRead_RepeatableRead ./...

# Lost Update Naive vs Safe Locking
go test -v -run TestNaiveTransfer_LostUpdate ./...
go test -v -run TestSafeTransferWithLock_DeterministicLocking ./...

# Deadlock-Free Bidirectional Transfer
go test -v -run TestBidirectionalTransfers_DeadlockFree ./...

# 100 Concurrent Transfers Stress Test
go test -v -run TestHighContention_100ConcurrentTransfers ./...

# Serializable Conflict & Retry
go test -v -run TestSerializable_ConcurrentUpdate_SerializationFailure ./...
```

## Expected Results
- `TestNaiveTransfer_LostUpdate`: Mengamati total uang bertambah melewati batas (misal menjadi 3.800.000) yang membuktikan reproduksi bug.
- `TestSafeTransferWithLock`: Menolak transfer kedua dan mempertahankan saldo total tetap 3.000.000.
- `TestBidirectionalTransfers_DeadlockFree`: 100 transfer bolak-balik selesai tanpa error deadlock `40P01`.
- `TestHighContention_100ConcurrentTransfers`: Seluruh invariant (`balance >= 0`, `total_money == 150000`) terbukti valid 100%.

## Production Lessons
1. **BEGIN/COMMIT Bukan Perisai Ajaib**: Tanpa locking eksplisit atau isolasi yang tepat, race condition tetap tembus.
2. **Ketahui Database Engine Anda**: Standar ANSI SQL berbeda dengan realita engine (PostgreSQL tidak punya Dirty Read dan mencegah Phantom Read di Repeatable Read).
3. **Pilih Strategi Berdasarkan Beban**:
   - High Contention: `SELECT ... FOR UPDATE` dengan deterministic ordering lebih stabil daripada membiarkan ratusan serializable transaction gagal dan retry berulang kali.
   - Low Contention / Complex Invariants: `SERIALIZABLE` dengan retry policy terarah.
4. **Wajib Retryable Filter**: Jangan pernah me-retry error non-transient seperti constraint violation (`23505`) atau syntax error. Hanya retry `40001` (serialization failure) dan transient deadlock `40P01`.

## Trade-offs

| Pendekatan | Kelebihan | Kelemahan / Biaya |
|---|---|---|
| **Read Committed (Naive)** | Throughput maksimal, tanpa lock wait | Rentan lost update dan inkonsistensi data |
| **Read Committed + Row Lock** | Menjamin invariant, latensi deterministik | Membutuhkan kehati-hatian urutan lock agar tidak deadlock |
| **Repeatable Read** | Snapshot konsisten tanpa lock | Memerlukan handling conflict `40001` pada concurrent update |
| **Serializable (SSI)** | Menjamin serial execution equivalent, bebas anomali write skew | Abort rate tinggi pada skenario high contention, wajib retry handling |

## Senior Engineer Checklist
- [ ] Apakah business invariant telah didefinisikan secara eksplisit sebelum memilih solusi?
- [ ] Apakah transaksi yang melibatkan multi-row lock sudah menggunakan urutan deterministik (misal `MIN(id)` lalu `MAX(id)`)?
- [ ] Apakah aplikasi telah menangani error code `40001` (serialization_failure) jika menggunakan REPEATABLE READ atau SERIALIZABLE?
- [ ] Apakah retry logic dibatasi dengan max attempts, exponential backoff, dan full jitter?
- [ ] Apakah test concurrency menggunakan deterministic barrier dan memverifikasi invariant sistem, bukan sekadar memeriksa ketiadaan error?

---

## Navigation

- **Previous**: [Lab 07 — Observability](../07-observability/)
- **Next**: [Lab 09 — Circuit Breaker](../09-circuit-breaker/)
