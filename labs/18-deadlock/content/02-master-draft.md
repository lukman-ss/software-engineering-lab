# Deadlock: Mekanisme Terjadinya, Penanganan Database, dan Strategi Pencegahan di Tingkat Aplikasi

## Problem
Pada sistem yang memproses transaksi secara konkuren, beberapa proses atau goroutine sering kali mengakses dan memodifikasi sumber daya bersama (shared resources) seperti baris data (rows) atau tabel dalam basis data. Ketika dua transaksi berusaha mengunci (lock) sumber daya yang sama namun dengan urutan yang berbeda, situasi dapat berujung pada kondisi saling menunggu tanpa akhir.

Kondisi saling menunggu ini dikenal sebagai **deadlock**. Kedua transaksi tidak dapat melanjutkan operasinya secara mandiri karena masing-masing pihak memegang sumber daya yang dibutuhkan oleh pihak lain. Tanpa campur tangan eksternal dari database engine atau runtime sistem, transaksi-transaksi tersebut akan terhenti selamanya (*deadly embrace*), menahan resource sistem dan menyebabkan penurunan performa secara drastis hingga sistem berhenti merespons.

## Why This Matters
Deadlock bukanlah sekadar isu teoretis, melainkan peristiwa riil yang kerap muncul seiring bertambahnya volume transaksi per detik (TPS) pada sistem produksi:
1. **Peningkatan Latensi**: Transaksi yang terlibat deadlock tertahan hingga terdeteksi oleh sistem, memperpanjang waktu respons.
2. **Kegagalan Transaksi (Aborted Transactions)**: Sistem basis data (seperti PostgreSQL atau SQL Server) akan membatalkan salah satu transaksi untuk memutus siklus tunggu. Jika aplikasi tidak siap menangani pembatalan ini, pengguna akan menerima kegagalan (error).
3. **Pemborosan Sumber Daya Komputasi**: Transaksi yang dibatalkan telah mengonsumsi waktu CPU, memori, dan I/O sebelum akhirnya harus di-rollback.

Deadlock sering kali disalahartikan sebagai bug pada database engine, padahal mayoritas penyebabnya terletak pada pola pengaksesan data yang tidak teratur di tingkat aplikasi.

## Mental Model
Model mental deadlock dapat diilustrasikan melalui siklus dependensi (circular wait):

```text
Transaksi 1 memegang Lock A ---> Transaksi 1 membutuhkan Lock B
             ^                                      |
             |                                      v
Transaksi 2 membutuhkan Lock A <--- Transaksi 2 memegang Lock B
```

- **Transaksi 1** mengunci Akun A, lalu mencoba mengunci Akun B.
- **Transaksi 2** mengunci Akun B, lalu mencoba mengunci Akun A secara bersamaan.
- **Transaksi 1** menunggu Akun B dibebaskan oleh Transaksi 2.
- **Transaksi 2** menunggu Akun A dibebaskan oleh Transaksi 1.

Siklus ini merupakan *deadly embrace*. Keduanya tidak akan pernah selesai tanpa adanya entitas luar yang secara sepihak memutus siklus tersebut.

## Core Concept
Ada lima konsep utama yang mendasari dinamika deadlock:

1. **Circular Wait (Siklus Tunggu)**: Terjadi ketika himpunan proses saling menunggu sumber daya yang ditahan oleh proses lain dalam himpunan tersebut secara berputar.
2. **Deadlock Monitor & Deadlock Victim**: Database engine modern memiliki thread monitor independen yang mendeteksi siklus tunggu secara periodik. Jika siklus ditemukan, database memilih satu transaksi sebagai *deadlock victim* dan membatalkannya (*abort/rollback*) untuk membebaskan kuncinya.
3. **Deadlock Timeout (`deadlock_timeout`)**: Pengecekan siklus secara berkala membutuhkan beban komputasi yang tidak murah. Database umumnya menunggu jeda waktu tertentu (misalnya default 1 detik di PostgreSQL atau interval 5 detik di SQL Server) sebelum menjalankan algoritma deteksi siklus.
4. **Lock Ordering**: Pola penguncian sumber daya yang seragam dan konsisten berdasarkan urutan tertentu (misalnya urutan alfabetis atau numerik ID sumber daya). Teknik ini memutus kemungkinan terbentuknya *circular wait*.
5. **Application-Level Retry**: Deadlock tidak selalu dapat dihindari 100% pada sistem dengan konkurensi ekstrem. Oleh karena itu, aplikasi harus memiliki kemampuan untuk mendeteksi error deadlock dan melakukan eksekusi ulang (*retry*) secara otomatis.

## Failure Scenario
Skenario kegagalan paling umum adalah transfer dana secara konkuren antara dua akun:

- **Goroutine/Koneksi 1**: Melakukan transfer dari `Akun A` ke `Akun B`. Goroutine 1 mengakuisisi lock pada `Akun A`.
- **Goroutine/Koneksi 2**: Melakukan transfer dari `Akun B` ke `Akun A`. Goroutine 2 mengakuisisi lock pada `Akun B`.
- Goroutine 1 berupaya mengunci `Akun B`, namun terblokir karena `Akun B` sedang dikunci oleh Goroutine 2.
- Goroutine 2 berupaya mengunci `Akun A`, namun terblokir karena `Akun A` sedang dikunci oleh Goroutine 1.
- Keduanya menunggu satu sama lain hingga context timeout tercapai, memunculkan error `ErrDeadlock` (simulasi *deadlock victim*).

## How It Works
### Deteksi dan Penanganan oleh Database
Di dalam database engine relasional:
1. Setiap transaksi meminta lock atas data sebelum melakukan mutasi.
2. Jika lock tidak langsung tersedia, transaksi dimasukkan ke antrean tunggu (*wait queue*).
3. Database engine memantau relasi antar transaksi menggunakan *Wait-For Graph* (WFG). Dalam graf ini:
   - Node merepresentasikan transaksi.
   - Sisi berarah (*directed edge*) merepresentasikan relasi tunggu (Transaksi X menunggu Transaksi Y).
4. Komponen *deadlock monitor* berjalan di latar belakang secara periodik. Jika ditemukan siklus tertutup dalam graf (*cycle in WFG*), status deadlock dinyatakan terjadi.
5. Mesin basis data memilih satu transaksi sebagai korban (*victim*) berdasarkan metrik tertentu (misalnya transaksi yang paling sedikit mengubah log, atau berbasis prioritas) dan membatalkannya. Transaksi tersebut menerima error spesifik (misalnya error code 1205 pada SQL Server atau error `deadlock detected` pada PostgreSQL).

### Penanganan oleh Aplikasi
Aplikasi memutus siklus sebelum terbentuk dengan:
1. **Lock Ordering**: Mengurutkan entitas yang akan dimutasi berdasarkan identitas unik (misalnya `id_a < id_b`). Setiap transaksi selalu mengunci entitas dengan identifier lebih kecil terlebih dahulu, sehingga siklus mustahil terbentuk.
2. **Durasi Transaksi Minimum**: Menjaga blok transaksi tetap sesingkat mungkin dan menghindari panggilan eksternal (seperti HTTP call atau pemrosesan berat) di dalam blok transaksi aktif.
3. **Mekanisme Retry**: Menangkap error deadlock dan melakukan *backoff* singkat sebelum mencoba ulang transaksi secara transparan bagi pengguna.

## Architecture
Implementasi lab ini memodelkan sistem penguncian baris basis data menggunakan channel Go dengan timeout context untuk mensimulasikan deadlock detection:

```text
+------------------------------------------------------------+
|                        cmd/demo                            |
|       Menjalankan demo Naive, Ordered, dan Retry           |
+------------------------------------------------------------+
                             |
                             v
+------------------------------------------------------------+
|                   internal/transfer                        |
|  - TransferNaive: Meminta lock sesuai urutan parameter    |
|  - TransferOrdered: Mengurutkan lock berdasarkan ID Akun   |
|  - TransferWithRetry: Menjalankan TransferNaive + Retry    |
+------------------------------------------------------------+
                             |
                             v
+------------------------------------------------------------+
|                     internal/bank                          |
|  - Account: Entitas dengan ID, Balance, dan channel lock   |
|  - Lock(ctx): Mengakuisisi lock atau return ErrDeadlock     |
|  - Unlock(): Membebaskan channel lock                      |
+------------------------------------------------------------+
```

## Implementation
Komponen utama lab terdiri dari:
- `internal/bank/account.go`: Mendefinisikan model `Account` dan metode `Lock(ctx)` berbasis channel buffer 1 slot. Ketika channel terisi, lock tersedia. Saat channel kosong, pemanggil harus menunggu atau dibatalkan ketika context timeout terpicu (`ctx.Done()`), yang mengembalikan `ErrDeadlock`.
- `internal/transfer/transfer.go`:
  - `TransferNaive`: Mengakuisisi lock akun asal (`from`) kemudian akun tujuan (`to`). Jika dua transfer berlawanan arah terjadi bersamaan, terjadi deadlock.
  - `TransferOrdered`: Mengurutkan pointer akun berdasarkan nilai `ID`. Akun dengan nilai `ID` lebih kecil selalu dikunci terlebih dahulu, memutus circular wait.
  - `TransferWithRetry`: Membungkus `TransferNaive` dengan context timeout pendek (simulasi penundaan monitor) dan melakukan perulangan hingga batas `maxRetries`.

## Code Walkthrough
### 1. Model Akun dan Mekanisme Penguncian Berbatas Waktu
File: `internal/bank/account.go`
```go
package bank

import (
	"context"
	"errors"
)

var ErrDeadlock = errors.New("deadlock victim")

type Account struct {
	ID      string
	Balance int
	ch      chan struct{}
}

func NewAccount(id string, bal int) *Account {
	ch := make(chan struct{}, 1)
	ch <- struct{}{}
	return &Account{ID: id, Balance: bal, ch: ch}
}

func (a *Account) Lock(ctx context.Context) error {
	select {
	case <-a.ch:
		return nil
	case <-ctx.Done():
		return ErrDeadlock
	}
}

func (a *Account) Unlock() {
	a.ch <- struct{}{}
}
```
Ketika `Lock` dipanggil dan channel kosong, goroutine memblokir hingga channel kembali terisi atau context kedaluwarsa. Context timeout mewakili tindakan *deadlock monitor* yang mengorbankan transaksi.

### 2. Pendekatan Naif yang Memicu Deadlock
File: `internal/transfer/transfer.go`
```go
// TransferNaive creates a potential deadlock by locking from then to.
func TransferNaive(ctx context.Context, from, to *bank.Account, amount int, delay time.Duration) error {
	if err := from.Lock(ctx); err != nil {
		return err
	}
	defer from.Unlock()

	time.Sleep(delay) // simulate transaction duration

	if err := to.Lock(ctx); err != nil {
		return err
	}
	defer to.Unlock()

	from.Balance -= amount
	to.Balance += amount
	return nil
}
```
Urutan penguncian ditentukan sepenuhnya oleh urutan input parameter (`from` lalu `to`). Jika Transaksi 1 mentransfer A ke B dan Transaksi 2 mentransfer B ke A secara konkuren, siklus terbentuk. Parameter `delay` mensimulasikan durasi pemrosesan di dalam transaksi.

### 3. Pencegahan Menggunakan Lock Ordering
File: `internal/transfer/transfer.go`
```go
// TransferOrdered prevents deadlock by always locking accounts in alphabetical order.
func TransferOrdered(ctx context.Context, acc1, acc2 *bank.Account, amount int, delay time.Duration) error {
	first, second := acc1, acc2
	if acc1.ID > acc2.ID {
		first, second = acc2, acc1
	}

	if err := first.Lock(ctx); err != nil {
		return err
	}
	defer first.Unlock()

	time.Sleep(delay)

	if err := second.Lock(ctx); err != nil {
		return err
	}
	defer second.Unlock()

	acc1.Balance -= amount
	acc2.Balance += amount
	return nil
}
```
`TransferOrdered` memeriksa nilai leksikografis `acc1.ID` dan `acc2.ID`. Sumber daya dengan ID lebih kecil selalu dikunci lebih dulu. Kedua transaksi konkuren akan bersaing memperebutkan lock pertama yang sama, sehingga salah satu harus mengantre secara tertib alih-alih saling mengunci.

### 4. Pemulihan Menggunakan Application-Level Retry
File: `internal/transfer/transfer.go`
```go
// TransferWithRetry handles deadlocks by retrying the operation.
func TransferWithRetry(ctx context.Context, from, to *bank.Account, amount int, delay time.Duration, maxRetries int) error {
	for i := 0; i < maxRetries; i++ {
		// Use a short timeout for the attempt to simulate deadlock monitor aborting fast
		attemptCtx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
		err := TransferNaive(attemptCtx, from, to, amount, delay)
		cancel()

		if err == nil {
			return nil // Success
		}
		if err != bank.ErrDeadlock && err != context.DeadlineExceeded {
			return err // Some other error
		}
		
		// Backoff before retry
		time.Sleep(2 * time.Millisecond)
	}
	return bank.ErrDeadlock
}
```
Jika operasi naif menghasilkan `ErrDeadlock` atau `context.DeadlineExceeded`, perulangan melakukan jeda waktu singkat (backoff) lalu mengulangi transaksi dari awal hingga batas `maxRetries`.

## What the Tests Prove
Pengujian otomatis pada file `tests/transfer_test.go` memverifikasi empat perilaku esensial:

1. **Deadlock Terjadi pada Pola Naif (`TestDeadlockOccurrence`)**:
   Dua goroutine yang menjalankan `TransferNaive` bolak-balik (A ke B, dan B ke A) secara konkuren terbukti memicu kegagalan dengan kembalian `bank.ErrDeadlock` pada salah satu transaksi.
2. **Lock Ordering Mencegah Deadlock (`TestLockOrderingPreventsDeadlock`)**:
   Dua transfer konkuren bolak-balik berhasil dieksekusi secara berurutan tanpa error (`err == nil`) ketika menggunakan `TransferOrdered`, dan saldo akhir kedua akun tetap konsisten.
3. **Mekanisme Retry Memulihkan Transaksi yang Gugur (`TestRetryRecoversDeadlock`)**:
   Operasi transfer naif konkuren yang dibungkus `TransferWithRetry` terbukti berhasil menyelesaikan kedua transaksi secara penuh setelah melakukan retry pada transaksi yang sempat dibatalkan.
4. **Durasi Transaksi Memperbesar Peluang Deadlock (`TestTransactionDurationImpact`)**:
   Pengujian membandingkan iterasi transfer dengan jeda waktu transaksi (`5ms`) dan tanpa jeda (`0ms`). Hasil pengujian membuktikan bahwa durasi transaksi yang lebih panjang menghasilkan frekuensi deadlock yang lebih tinggi atau sama dengan durasi singkat (`longDeadlocks >= shortDeadlocks`).

Semua pengujian lolos tanpa race condition saat dieksekusi menggunakan perintah:
```bash
go test -v ./...
go test -race ./...
```

## Recovery / Rollback
Ketika basis data mendeteksi siklus deadlock:
1. **Rollback Otomatis di Tingkat Basis Data**: Database engine secara otomatis membatalkan seluruh perubahan yang dibuat oleh transaksi korban (deadlock victim), melepaskan seluruh lock yang sedang dipegang oleh transaksi tersebut, dan mengembalikan error ke koneksi klien.
2. **Penanganan di Tingkat Aplikasi**:
   - Tangkap exception atau error code deadlock (misal: Postgres SQLSTATE `40P01` atau SQL Server Error 1205).
   - Jangan teruskan error sebagai kegagalan permanen ke pengguna akhir secara langsung.
   - Lakukan rollback pada transaction context di aplikasi jika library ORM/database driver belum melakukannya secara otomatis.
   - Berikan jeda waktu acak (jittered backoff) sebelum mengeksekusi ulang seluruh blok transaksi untuk menghindari resonansi konkurensi (thundering herd).

## Production Considerations
1. **Urutan Penguncian Harus Universal**: Jika tabel A diakses sebelum tabel B pada satu use case, seluruh endpoint atau background job lain dalam sistem harus mengakses tabel A sebelum tabel B. Inkonsistensi urutan penguncian pada satu fungsi saja cukup untuk menciptakan deadlock.
2. **Hindari Interaksi Jaringan Eksternal dalam Transaksi**: Melakukan panggilan API pihak ketiga (misalnya HTTP request ke payment gateway) di dalam blok transaksi database akan memperpanjang waktu penahanan lock secara signifikan, meningkatkan risiko deadlock secara eksponensial.
3. **Penyetelan `deadlock_timeout`**: Pada PostgreSQL, nilai default `deadlock_timeout` adalah 1 detik. Menurunkan nilai ini membuat deteksi deadlock lebih cepat, namun menambah beban CPU server.
4. **Batasi Jumlah Retry**: Tetapkan batasan `maxRetries` (misalnya 3 hingga 5 kali) untuk mencegah loop tanpa akhir apabila terjadi kegagalan persisten atau degradasi sistem yang parah.

## Common Mistakes
1. **Membiarkan Koneksi Melakukan I/O Lambat dalam Transaksi**: Mengunggah file atau menunggu input pengguna saat transaksi database masih terbuka.
2. **Mengabaikan Pengurutan ID pada Operasi Batch**: Melakukan `UPDATE` atau `SELECT FOR UPDATE` pada sekumpulan row tanpa mengurutkan ID kumpulan row tersebut (`ORDER BY id ASC`).
3. **Mengasumsikan Retry Dapat Mengabaikan Idempotensi**: Menjalankan ulang transaksi yang memiliki efek samping non-idempotent (misalnya mengirim email atau menembak webhook eksternal) tanpa pemeriksaan status.
4. **Self-Transfer Tanpa Validasi**: Tidak memvalidasi kesamaan akun asal dan tujuan (`from == to`). Upaya mengunci akun yang sama dua kali dalam satu alur eksekusi akan menyebabkan self-deadlock.

## Checklist
Sebelum merilis kode yang memodifikasi multiple shared resources:
- [ ] Apakah seluruh entitas yang dimutasi dalam transaksi diakses dengan urutan identik (Lock Ordering) di seluruh codebase?
- [ ] Apakah operasi batch/bulk melakukan pengurutan ID sebelum penguncian?
- [ ] Apakah transaksi dibuat sesingkat mungkin tanpa operasi I/O eksternal di dalamnya?
- [ ] Apakah aplikasi menangkap error deadlock dari database driver?
- [ ] Apakah retry logic telah diimplementasikan dengan backoff dan batasan perulangan maksimum?
- [ ] Apakah validasi akun/sumber daya identik (`from == to`) telah diuji untuk mencegah self-deadlock?

## Key Takeaways
1. Deadlock adalah kondisi sistemik di mana dua transaksi atau lebih saling menunggu lock yang ditahan pihak lain (*circular wait*).
2. Database engine menyelesaikan deadlock secara paksa dengan membatalkan (*abort/rollback*) salah satu transaksi sebagai *deadlock victim*.
3. Pencegahan paling efektif di tingkat aplikasi adalah menerapkan **Lock Ordering** yang konsisten di semua jalur eksekusi.
4. Memperpendek durasi transaksi secara signifikan menurunkan probabilitas terjadinya deadlock.
5. Mekanisme **Application Retry** adalah standar operasional wajib pada sistem konkuren tinggi untuk memulihkan transaksi yang menjadi korban deadlock secara transparan.

## Sources
- **Source 1**: PostgreSQL 18 Documentation: 13.3. Explicit Locking (The PostgreSQL Global Development Group, https://www.postgresql.org/docs/current/explicit-locking.html)
- **Source 2**: PostgreSQL 18 Documentation: 19.12. Lock Management (The PostgreSQL Global Development Group, https://www.postgresql.org/docs/current/runtime-config-locks.html)
- **Source 3**: Deadlocks Guide - SQL Server (Microsoft Learn, https://learn.microsoft.com/en-us/sql/relational-databases/sql-server-deadlocks-guide)
