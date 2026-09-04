# Lab 08 — Database Isolation Level

## Problem
Banyak engineer berasumsi bahwa membungkus sederet query SQL di dalam `BEGIN TRANSACTION` dan `COMMIT` secara otomatis membuat aplikasi bebas dari race condition dan data anomaly. Di production dengan ratusan request per detik, asumsi ini sering kali menyebabkan saldo bocor, overbooking, double spending, hingga laporan keuangan yang tidak konsisten.

`BEGIN TRANSACTION` dan `COMMIT` hanya menentukan **transaction boundary** (cakupan atomisitas). Kebenaran (correctness) aplikasi di bawah eksekusi *concurrent* tetap bergantung pada: **isolation level**, cara membaca/mengunci data (*locking strategy*), *business invariant*, dan penanganan *serialization conflict* (retry strategy).

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

### Default Behavior & Implementation Reality (PostgreSQL)

Standar SQL hanyalah spesifikasi minimum. Database engine memiliki implementasi internal yang dapat lebih kuat daripada spesifikasi ANSI. Di lab ini (PostgreSQL):
- **Default Isolation Level**: `READ COMMITTED`.
- **READ UNCOMMITTED**: PostgreSQL **tidak memiliki** implementasi Read Uncommitted yang mengizinkan Dirty Read. Jika aplikasi meminta `READ UNCOMMITTED`, PostgreSQL secara diam-diam mengeksekusinya setara dengan `READ COMMITTED`.
- **READ COMMITTED**: Menggunakan snapshot baru pada setiap statement SQL yang dieksekusi.
- **REPEATABLE READ**: Menggunakan snapshot MVCC statis di awal transaksi (transaction snapshot) sehingga PostgreSQL tidak mengizinkan Phantom Read klasik. Namun, level ini **masih dapat mengalami serialization anomaly** (seperti Write Skew) dan akan menggagalkan transaksi concurrent dengan error `SQLSTATE 40001 (serialization_failure)` jika terjadi benturan update pada baris yang sama.
- **SERIALIZABLE**: Bukan berarti semua transaksi benar-benar dijalankan satu per satu atau semua operasi dikunci secara global. PostgreSQL Serializable menggunakan **Serializable Snapshot Isolation (SSI)**: transaksi concurrent tetap dieksekusi paralel, dan database membatalkan salah satu transaksi dengan `serialization_failure` (40001) apabila execution tidak dapat diserialisasi.

*(Catatan: Karakteristik locking/anomali pada engine lain seperti MySQL InnoDB berbeda. InnoDB defaultnya REPEATABLE READ dan menggunakan Gap Locks/Next-Key Locks).*

## PostgreSQL dan MVCC

PostgreSQL mengimplementasikan transaksi menggunakan **Multi-Version Concurrency Control (MVCC)**. Dalam MVCC, setiap penulisan data menghasilkan versi tuple baru (`xmin`, `xmax`) tanpa menimpa data lama secara langsung. Pembacaan ordinary `SELECT` tidak memblokir penulisan `UPDATE`, dan sebaliknya.

## Read Committed
Level default di PostgreSQL. Setiap query SQL dalam transaksi mengambil **snapshot baru** dari data yang telah di-commit pada saat query tersebut dieksekusi.
- Anomali yang mungkin muncul: *Non-Repeatable Read*, *Phantom Read*, dan *Lost Update* (pada pola check-then-act).

## Repeatable Read
Snapshot database diambil hanya satu kali saat transaksi pertama kali membaca data.
- Menjamin pembacaan berulang selalu konsisten.
- Jika transaksi mencoba meng-`UPDATE` baris yang telah di-update & di-commit oleh transaksi lain setelah snapshot diambil, PostgreSQL akan menggagalkan transaksi dengan error code `SQLSTATE 40001 (serialization_failure)`.

## Serializable
Tingkat isolasi tertinggi. PostgreSQL Serializable menggunakan **Serializable Snapshot Isolation (SSI)**.

**PERINGATAN KONSEP**: Serializable **tidak** berarti database engine akan mengunci seluruh tabel dan mengeksekusi request murni satu-per-satu secara literal (Global Locking). SSI membiarkan transaksi concurrent berjalan paralel. Database hanya akan *mendeteksi dependency/conflict cycle*. Jika hasilnya **tidak ekuivalen** dengan suatu valid serial execution, database akan menggagalkan salah satu transaksi dengan melemparkan error:
```text
ERROR: could not serialize access due to concurrent update (SQLSTATE 40001)
```

Aplikasi **wajib** menangkap error ini dan melakukan *Bounded Retry*.

##### Analogi SERIALIZABLE (yang akurat)
Transaksi concurrent dapat berjalan bersamaan. Namun hasil commit harus ekuivalen dengan suatu serial execution; database dapat membatalkan transaksi dan perlu di-retry jika terjadi serialization conflict. Buktinya: *TestSerializable_ConcurrentUpdate_SerializationFailure*.

## Wallet Transfer: REPEATABLE READ BUKAN solusi otomatis untuk read-check-write

Pola `read balance → validate → write balance` di bawah `REPEATABLE READ`:
- *Snapshot consistency* aman untuk pembacaan berulang.
- Namun, dua transaksi masih bisa membaca saldo yang sama, dan yang pertama commit berhasil. Yang kedua akan gagal dengan `40001` **hanya jika benturan pada UPDATE yang sama**.
- Jika aplikasi hanya `read` untuk validasi lalu lakukan logika lokal (bukan *check-modify-write* pada baris yang sama), tidak ada transaction conflict.

Untuk pola read-check-write kritis, pertimbangkan:
- `SELECT ... FOR UPDATE` (pessimistic locking), atau
- Atomic conditional update (`WHERE balance >= check_value`), atau  
- `SERIALIZABLE` + retry (optimistic locking).

## SELECT FOR UPDATE
Mekanisme Pessimistic Locking eksplisit. Mengambil row-level exclusive lock pada baris terpilih:
- Transaksi lain yang mencoba mengambil lock `FOR UPDATE` atau meng-`UPDATE` baris yang sama akan **tertahan (blocked)** sampai transaksi pertama commit atau rollback.
- Pembaca biasa (`SELECT` tanpa lock) **tetap tidak terblokir** berkat snapshot MVCC.

**Deadlock Prevention dengan Deterministic Lock Ordering**:
Jika Transaksi 1 mentransfer `A -> B` (lock A lalu B) dan Transaksi 2 mentransfer `B -> A` (lock B lalu A) secara bersamaan tanpa urutan, database dapat mengalami deadlock (`SQLSTATE 40P01`) karena circular wait.
Solusi wajib: Gunakan **Deterministic Lock Ordering** — selalu urutkan lock dari ID akun terkecil ke terbesar. Dengan konsisten me-lock `MIN(id)` lalu `MAX(id)`, circular wait terhindari.

```go
firstID, secondID := fromID, toID
if firstID > secondID {
    firstID, secondID = secondID, firstID
}
tx.QueryRowContext(ctx, "SELECT balance FROM accounts WHERE id = $1 FOR UPDATE", firstID)
tx.QueryRowContext(ctx, "SELECT balance FROM accounts WHERE id = $2 FOR UPDATE", secondID)
```

## Lost Update
Terjadi saat aplikasi melakukan pola **Check-Then-Act** atau `READ-CHECK-CALCULATE-WRITE` menggunakan query bacaan biasa di isolation level rendah:
1. Transaksi 1 membaca saldo Alice = 1.000.000.
2. Transaksi 2 membaca saldo Alice = 1.000.000 (belum ada commit).
3. Transaksi 1 memvalidasi saldo cukup, lalu mengurangi 800.000, menulis saldo Alice = 200.000, lalu commit.
4. Transaksi 2 memvalidasi saldo cukup (berdasarkan bacaan lama), mengurangi 800.000, menulis saldo Alice = 200.000, lalu commit.
Hasil: Alice mentransfer total 1.600.000 padahal saldonya hanya 1.000.000! Pengurangan pertama ditimpa dan hilang (lost). Hal ini dapat diamati pada `TestNaiveTransfer_LostUpdate`.

## Serialization Failure
Ketika menggunakan `REPEATABLE READ` atau `SERIALIZABLE`, PostgreSQL menolak eksekusi concurrent yang berbenturan dengan mengembalikan error:
```text
ERROR: could not serialize access due to concurrent update (SQLSTATE 40001)
```
Aplikasi yang menggunakan isolation level tinggi **wajib menyiapkan mekanisme retry terarah** (bounded retries dengan exponential backoff dan jitter).

Contoh implementasi Retry Wrapper yang aman ada pada `isolation.RetryTransaction()`:
```go
func (r *PostgresWalletRepo) TransferSerializableWithRetry(ctx context.Context, fromID, toID int, amount int64, maxAttempts int) error {
	return RetryTransaction(ctx, maxAttempts, policy, func(ctx context.Context) error {
		return r.TransferSerializable(ctx, fromID, toID, amount)
	})
}
```

## Database Engine vs Standard SQL
Penjelasan di README ini membedakan secara tegas antara **ANSI/ISO SQL-92 Standard** dengan perilaku konkret **PostgreSQL engine** (MVCC & SSI). Karakteristik performa dan penguncian pada database lain (misalnya MySQL InnoDB yang memiliki Gap Locks / Next-Key Locks) dapat berbeda dan memiliki perilaku anomali/locking tersendiri.

## Wallet Transfer Case
Studi kasus sistem dompet digital dengan 3 akun:
```text
Alice   : Rp 1.000.000
Bob     : Rp 1.000.000
Charlie : Rp 1.000.000
```
- **Invariant 1**: `balance >= 0` (Saldo tidak boleh negatif).
- **Invariant 2**: `total_money == 3.000.000` (Konservasi total uang dalam ekosistem).

### Rekomendasi Strategi & Trade-offs
Jangan menentukan isolation level hanya berdasarkan data berkaitan dengan uang. Business invariant harus menjadi dasar keputusan:

1. **READ COMMITTED + SELECT ... FOR UPDATE + Deterministic Lock Ordering**
   - **Cocok ketika**: Invariant berpusat pada beberapa row tertentu, contention relatif tinggi, dan kita ingin serialize akses terhadap row tersebut secara eksplisit.
   - **Karakteristik**: Solusi ini lebih predictable dibanding Serializable + retry pada wallet dengan row contention tinggi.
2. **REPEATABLE READ**
   - **Cocok ketika**: Membutuhkan consistent snapshot dan concurrent update conflict dapat ditangani dengan retry.
3. **SERIALIZABLE**
   - **Cocok ketika**: Terdapat complex cross-row/business invariant, write skew harus dicegah, dan contention memungkinkan retry strategy.

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

Jalankan tiap skenario untuk melihat visualisasi dan perilaku engine:

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

# Serializable Conflict (Bukti 40001) & Retry
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

### Panduan Pemilihan Isolation Strategy
Tidak ada isolation level yang otomatis benar hanya karena domain-nya financial. Tentukan business invariant terlebih dahulu sebelum memilih.

| Scenario | Strategy |
|---|---|
| CRUD umum | `Read Committed` |
| Dashboard / agregasi statistik sederhana | `Read Committed` |
| Consistent report snapshot (multi-table cross-section) | `Repeatable Read` / explicit reporting snapshot |
| Wallet transfer dengan high contention | `Read Committed` + `SELECT ... FOR UPDATE` (Pessimistic) |
| Concurrent update dengan snapshot statis | `Repeatable Read` + retry (Optimistic) |
| Complex multi-row invariant / Pencegahan Write Skew | `Serializable` + retry (SSI) |

## FAQs: Lab Hari #8

1. **Isolation level apa yang paling masuk akal untuk wallet transfer?**
   `READ COMMITTED` dipadukan dengan `SELECT ... FOR UPDATE` (Pessimistic Locking) dan Deterministic Lock Ordering. Skenario wallet umumnya memiliki contention tinggi pada row user tertentu.
2. **Apa risiko jika menggunakan READ COMMITTED tanpa mekanisme tambahan?**
   *Lost Update*. Dua transaksi yang melakukan "read-calculate-write" secara bersamaan akan menimpa satu sama lain, menyebabkan saldo tercipta atau menguap dari kehampaan.
3. **Kapan menggunakan SELECT ... FOR UPDATE?**
   Ketika aplikasi melakukan *Check-Then-Act* business logic yang kritis terhadap suatu state (misal saldo >= 0), dan keputusan tidak boleh didasarkan pada data *stale*.
4. **Bagaimana menguji dua transfer concurrent?**
   Jangan mengandalkan `time.Sleep()`. Gunakan *barrier synchronization* (seperti `sync.WaitGroup` atau blocking channel) agar semua goroutine tertahan di titik query, lalu lepaskan serentak untuk memaksimalkan tabrakan.
5. **Mengapa BEGIN TRANSACTION saja belum menjamin saldo benar?**
   Transaksionalitas (`A`tomicity) hanya menjamin semua query sukses bersama atau gagal bersama. Namun cara transaksi memandang transaksi *lain* bergantung pada `I`solation level.
6. **Apa perbedaan snapshot consistency dan row locking?**
   Snapshot (`REPEATABLE READ`) memberi jaminan data *terbaca* tidak berubah (consistency statis optimis). Row locking (`FOR UPDATE`) mencegah transaksi lain *memodifikasi* row sampai lock dilepas (konsistensi pesimis).
7. **Mengapa SERIALIZABLE bukan pilihan default untuk seluruh aplikasi?**
   Serializable memakan *cost* tinggi: abort rate akibat conflict `40001` akan meroket di skenario high-contention, memerlukan *retry logic* mutlak di aplikasi, dan overhead dependency tracking.
8. **Mengapa application mungkin perlu retry ketika menggunakan SERIALIZABLE?**
   Karena *Serializable* pada SSI PostgreSQL beroperasi optimis: membiarkan concurrency terjadi dan baru akan membatalkan salah satu dengan `SQLSTATE 40001 (serialization_failure)` jika ditemukan conflict cycle.

## Senior Engineer Checklist
- [ ] Apakah business invariant telah didefinisikan secara eksplisit sebelum memilih solusi?
- [ ] Apakah transaksi yang melibatkan multi-row lock sudah menggunakan urutan deterministik (misal `MIN(id)` lalu `MAX(id)`)?
- [ ] Apakah aplikasi telah menangani error code `40001` (serialization_failure) jika menggunakan REPEATABLE READ atau SERIALIZABLE?
- [ ] Apakah retry logic dibatasi dengan max attempts, exponential backoff, dan full jitter?
- [ ] Apakah test concurrency menggunakan deterministic barrier dan memverifikasi invariant sistem, bukan sekadar memeriksa ketiadaan error?

---

## Navigation

- **Previous**: [Lab 07 — Observability](../07-observability/)
- **Next**: [Lab 09 — Code Review](../09-code-review/)
