# Lab 08 — Database Isolation Level: Kenapa BEGIN TRANSACTION Saja Tidak Menjamin Data Selalu Benar?

> **Mental Model**: Membungkus beberapa query dengan `BEGIN` dan `COMMIT` **hanya menjamin Atomicity (All-or-Nothing)** dari scope transaksi tersebut, **bukan Isolation dari interleaving query transaksi lain**. Tanpa pemahaman isolation level dan locking, race condition dan data anomaly tetap terjadi.

---

## 1. Empat Fenomena Anomali (SQL Standard vs Realita)

Standar SQL-92 mendefinisikan fenomena anomali data concurrency:

| Anomali | Deskripsi | Minimal Level Pencegah (SQL-92) |
|---|---|---|
| **Dirty Read** | Membaca data yang sedang diubah oleh transaksi lain yang **belum di-commit** (jika rollback, data tadi dianggap fiktif). | `READ COMMITTED` |
| **Non-Repeatable Read (Fuzzy Read)** | Membaca row yang sama dua kali dalam satu transaksi, tetapi menghasilkan nilai berbeda karena transaksi lain meng-`UPDATE` / `DELETE` dan `COMMIT` di antaranya. | `REPEATABLE READ` |
| **Phantom Read** | Menjalankan query range (misal `SELECT COUNT(*) WHERE balance > 1000`), lalu query kedua menemukan baris baru ("hantu") karena transaksi lain meng-`INSERT` dan `COMMIT`. | `SERIALIZABLE` |
| **Lost Update** | Dua transaksi membaca nilai yang sama, menghitung nilai baru secara terpisah di aplikasi, lalu menulis balik sehingga update transaksi pertama tertimpa dan hilang. | `REPEATABLE READ` / Locking |

---

## 2. Perbedaan ANSI SQL vs PostgreSQL vs MySQL/InnoDB

Penting untuk dipahami bahwa setiap database engine mengimplementasikan isolasi dengan cara berbeda via **MVCC (Multi-Version Concurrency Control)**:

| Isolation Level | ANSI SQL-92 Standard | PostgreSQL Implementation | MySQL / InnoDB Implementation |
|---|---|---|---|
| **READ UNCOMMITTED** | Mengizinkan Dirty Read | **Diperlakukan sama persis seperti READ COMMITTED** (Dirty Read tidak mungkin terjadi di Postgres karena MVCC tidak membaca uncommitted tuple). | Mengizinkan Dirty Read (tanpa lock/undo log snapshot). |
| **READ COMMITTED** (Default) | Mencegah Dirty Read. Setiap statement melihat data committed terbaru. Non-Repeatable Read & Phantom Read mungkin terjadi. | Setiap statement SQL dalam transaksi mengambil **snapshot baru** dari data committed saat statement itu mulai berjalan. | Menggunakan read view baru untuk setiap query `SELECT`. |
| **REPEATABLE READ** | Mencegah Dirty Read & Non-Repeatable Read. Standar SQL masih mengizinkan Phantom Read. | Snapshot diambil **sekali saja pada awal transaksi (first query)**. Mencegah Dirty Read, Non-Repeatable Read, **DAN Phantom Read**. Bila terjadi concurrent update pada row yang sama, transaksi menghasilkan error `40001 (serialization failure)`. | Mencegah Phantom Read pada plain `SELECT` via snapshot MVCC, dan menggunakan **Gap Locking / Next-Key Locks** pada locking reads (`FOR UPDATE`). |
| **SERIALIZABLE** | Menjamin eksekusi transaksi ekuivalen dengan eksekusi sekuensial murni. | Menggunakan **SSI (Serializable Snapshot Isolation)**. Melacak dependency antar transaksi via SIREAD locks di memory tanpa blocking, lalu abort transaksi (`40001`) bila terdeteksi cycle/anomaly (Write Skew). | Mengubah semua plain `SELECT` menjadi `SELECT ... FOR SHARE` (locking reader). |

---

## 3. Skenario Utama: Wallet Transfer

Dua rekening awal:
```text
Alice   : Rp 1.000.000 (id: 1)
Bob     : Rp 1.000.000 (id: 2)
Charlie : Rp 1.000.000 (id: 3)
```

Total uang dalam sistem: **Rp 3.000.000** (Conservation Invariant).

### Pola Naive (Anti-Pattern: Check-then-Act)

```sql
-- Transaksi A (Transfer 800k ke Bob)
BEGIN; -- READ COMMITTED
SELECT balance FROM isolation_accounts WHERE id = 1; -- (membaca 1.000.000)
-- Aplikasi cek: 1.000.000 >= 800.000 (Valid!)
UPDATE isolation_accounts SET balance = 200000 WHERE id = 1;
UPDATE isolation_accounts SET balance = 1800000 WHERE id = 2;
COMMIT;
```

Jika Transaksi B (Transfer 800k ke Charlie) berjalan bersamaan dan membaca sebelum Transaksi A commit:
1. Transaksi B juga membaca saldo Alice `1.000.000`.
2. Transaksi B mengira saldo cukup, lalu menulis saldo Alice `200.000` dan Charlie `1.800.000`.
3. **Hasil**: Alice mengirim Rp 1.600.000 padahal uangnya hanya Rp 1.000.000. Total uang sistem menjadi **Rp 3.800.000** (Uang gaib bertambah Rp 800.000!).

---

## 4. Solusi & Perbandingannya

### 1. Read Committed + Row-Level Locking (`SELECT ... FOR UPDATE`)
- Mengunci baris data pengirim & penerima.
- Transaksi kedua akan **tertahan (blocked)** sampai transaksi pertama `COMMIT`.
- Transaksi kedua lalu membaca saldo terbaru (`200.000`), validasi saldo gagal (`200.000 < 800.000`), dan me-return error `ErrInsufficientFunds`.
- **Urutan Lock**: Selalu urutkan ID (misal `MIN(from, to)` lalu `MAX(from, to)`) untuk menghindari Deadlock.

### 2. Repeatable Read
- Begitu dua transaksi mencoba meng-`UPDATE` baris yang sama, PostgreSQL mendeteksi conflict dan menggagalkan transaksi kedua dengan error code `40001`:
  ```text
  ERROR: could not serialize access due to concurrent update (SQLSTATE 40001)
  ```

### 3. Serializable + Application Retry
- Menjamin isolasi penuh dari anomaly write-skew.
- Aplikasi harus siap menangani serialization failure dengan **Exponential Backoff dan Jitter**:
  ```go
  for attempt := 0; attempt <= maxRetries; attempt++ {
      err := repo.TransferSerializable(ctx, fromID, toID, amount)
      if isSerializationError(err) {
          sleepWithJitter(attempt)
          continue
      }
      return err
  }
  ```

---

## 5. Cara Menjalankan Test

Prasyarat: PostgreSQL aktif (`docker compose up -d postgres` atau `make infra-up`).

```bash
# Jalankan seluruh test isolasi database
go test -v ./labs/08-database-isolation-level/...

# Jalankan skenario pembuktian Non-Repeatable Read
go test -v ./labs/08-database-isolation-level/ -run TestReadCommitted_NonRepeatableRead

# Jalankan skenario pembuktian Snapshot Isolation
go test -v ./labs/08-database-isolation-level/ -run TestRepeatableRead_SnapshotIsolation

# Jalankan pembuktian Lost Update pada Naive Transfer
go test -v ./labs/08-database-isolation-level/ -run TestNaiveTransfer_LostUpdate

# Jalankan solusi Safe Locking
go test -v ./labs/08-database-isolation-level/ -run TestSafeTransferWithLock

# Jalankan pembuktian Serialization Failure & Retry
go test -v ./labs/08-database-isolation-level/ -run TestSerializableWithRetry
```

---

## Navigasi

- **Previous**: [Lab 07 — Observability](../07-observability/)
- **Next**: [Lab 09 — Circuit Breaker](../09-circuit-breaker/)
