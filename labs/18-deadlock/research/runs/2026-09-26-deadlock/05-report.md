# Research Report — Deadlock: Kenapa Dua Transaksi yang Benar Bisa Saling Mengunci Selamanya

## Research Question

Bagaimana transaksi yang berjalan bersamaan dapat saling mengunci (deadlock), bagaimana database modern (PostgreSQL/MySQL) menanganinya, dan bagaimana aplikasi (studi kasus Sistem PPOB) mengurangi serta menangani dampaknya melalui lock ordering, transaction duration, contention control, retry, dan observability?

## Executive Summary

Deadlock adalah kondisi siklus tunggu antar transaksi yang masing-masing memegang resource yang dibutuhkan transaksi lain. Empat Coffman conditions harus bersamaan untuk deadlock terjadi. PostgreSQL dan MySQL tidak membiarkan siklus berlangsung selamanya: sistem mendeteksi deadlock dan membatalkan satu transaksi korban (victim) dengan error `40P01 deadlock_detected` (PostgreSQL). Pencegahan paling efektif adalah konsistensi urutan lock di seluruh code path, transaksi sesingkat mungkin, dan menghindari operasi eksternal di dalam transaksi. Untuk operasi idempoten, retry dengan backoff setelah `40P01`/`40001` adalah strategi standar. Observability mengandalkan `log_lock_waits`, `pg_locks`, `pg_stat_activity` (`wait_event_type = Lock`), dan metrik hitungan error 40P01.

## Findings

### Finding 1 — Definisi dan Empat Syarat Deadlock (Coffman)

Claim: Deadlock hanya terjadi bila empat kondisi bersamaan: mutual exclusion, hold and wait, no preemption, dan circular wait.

Evidence: PostgreSQL docs: "two (or more) transactions each hold locks that the other wants... Transaction 1 holds A waits B, Transaction 2 holds B waits A, neither can proceed." Wikipedia merumuskan "only if all of the following conditions occur simultaneously" lalu daftar keempat Coffman conditions (1971).

Sources:
- https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-DEADLOCKS
- https://en.wikipedia.org/wiki/Deadlock_(computer_science)#Conditions

Confidence: HIGH

### Finding 2 — Database Membatalkan Satu Transaksi (Deadlock Victim)

Claim: PostgreSQL otomatis mendeteksi deadlock dan abort salah satu transaksi agar yang lain lanjut. Transaksi mana yang dibatalkan tidak dapat diprediksi dan tidak boleh diandalkan.

Evidence: "PostgreSQL automatically detects deadlock situations and resolves them by aborting one of the transactions involved, allowing the other(s) to complete. Exactly which transaction will be aborted is difficult to predict and should not be relied upon."

Sources:
- https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-DEADLOCKS

Confidence: HIGH

### Finding 3 — Error Code 40P01 dan 40001 (Retry Signal)

Claim: PostgreSQL mengembalikan `SQLSTATE 40P01 deadlock_detected` untuk deadlock dan `40001 serialization_failure` untuk konflik serialisasi di Repeatable Read/Serializable. Keduanya menandakan aplikasi harus retry seluruh transaksi.

Evidence: Table A.1 Appendix A: `40P01 deadlock_detected` (Class 40 Transaction Rollback), `40001 serialization_failure`. Docs 13.2.2: "Applications using this level must be prepared to retry transactions due to serialization failures... retry the whole transaction from the beginning." Docs 13.3.4: "deadlocks can be handled on-the-fly by retrying transactions that abort due to deadlocks."

Sources:
- https://www.postgresql.org/docs/current/errcodes-appendix.html
- https://www.postgresql.org/docs/current/transaction-iso.html#XACT-REPEATABLE-READ
- https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-DEADLOCKS

Confidence: HIGH

### Finding 4 — deadlock_timeout Bukan Lock Timeout

Claim: `deadlock_timeout` (default 1s) adalah interval sebelum pengecekan deadlock dijalankan; `lock_timeout` (default 0 = disabled) adalah batas tunggu per lock attempt lalu statement dibatalkan. Keduanya berbeda dari deadlock.

Evidence: `deadlock_timeout`: "amount of time to wait on a lock before checking to see if there is a deadlock condition... default is one second (1s)". `lock_timeout`: "Abort any statement that waits longer than the specified amount of time while attempting to acquire a lock... A value of zero (the default) disables the timeout."

Sources:
- https://www.postgresql.org/docs/current/runtime-config-locks.html
- https://www.postgresql.org/docs/current/runtime-config-client.html

Confidence: HIGH

### Finding 5 — Lock Ordering Mencegah Circular Wait

Claim: Pertahanan terbaik adalah semua transaksi mengunci resource dalam urutan konsisten (partial ordering). Jika A dan B keduanya mengurutkan `accounts` by id atau `User -> Invoice`, siklus tidak terbentuk.

Evidence: "The best defense against deadlocks is generally to avoid them by being certain that all applications using a database acquire locks on multiple objects in a consistent order. In the example above, if both transactions had updated the rows in the same order, no deadlock would have occurred."

Sources:
- https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-DEADLOCKS
- https://en.wikipedia.org/wiki/Deadlock_(computer_science)#Prevention (resource hierarchy / partial ordering)

Confidence: HIGH

### Finding 6 — Transaction Duration dan Contention

Claim: Menahan transaksi lama meningkatkan peluang deadlock dan waste cek deadlock. Operasi eksternal (API WhatsApp, generate PDF, upload S3) serta UPDATE tanpa filter selektif yang mengunci jutaan row memperparah contention.

Evidence: "So long as no deadlock situation is detected, a transaction seeking either a table-level or row-level lock will wait indefinitely... it is a bad idea for applications to hold transactions open for long periods of time (e.g., while waiting for user input)." Docs log: "Ideally the setting [deadlock_timeout] should exceed your typical transaction time". MVCC intro: limit lock contention via MVCC, tapi row-level locks tetap menahan writer.

Sources:
- https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-DEADLOCKS
- https://www.postgresql.org/docs/current/runtime-config-locks.html
- https://www.postgresql.org/docs/current/mvcc-intro.html

Confidence: HIGH

### Finding 7 — Pola PPOB: Di Dalam vs Setelah COMMIT

Claim: Flow PPOB ideal: di dalam transaksi hanya (1) kurangi saldo agen (row lock terurut by id) + (2) buat transaksi + (4) tambah riwayat. Kirim request ke provider eksternal dilakukan setelah COMMIT (outbox/polling). Jika provider perlu id transaksi, gunakan staging row + status PENDING.

Evidence: Prinsip Finding 6: operasi eksternal tidak memerlukan atomicity DB dan memperpanjang hold time. Praktik industri konsisten dengan "BEGIN -> DB -> COMMIT -> Email/WhatsApp" (topic spec).

Sources:
- https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-DEADLOCKS (indirect)
- Topic spec PPOB (application pattern, bukan DB docs)

Confidence: MEDIUM (DB docs mendukung prinsip; mapping spesifik PPOB adalah interpretasi engineering yang wajar)

### Finding 8 — Retry Idempoten dengan Backoff (Go)

Claim: Untuk transaksi idempoten, retry setelah `40P01`/`40001` dengan batas maksimum dan jeda singkat (exponential backoff + jitter) adalah praktik standar. Go stdlib `time.Sleep` / `time.After` + `context` menyediakan primitif; backoff+jitter didukung oleh AWS Architecture Blog dan AWS SDK (bukan sekadar pola umum). Hindari retry buta untuk transaksi non-idempoten.

Evidence: PG docs sarankan retry on-the-fly. Go `time.After(d)` / `Sleep(d)` untuk jeda; `Since`/`Until` monotonic untuk ukur durasi backoff tanpa terpengaruh wall-clock skew. AWS SDK: "Most AWS SDKs now support exponential backoff and jitter as part of their retry behavior" (Brooker, 2015).

Sources:
- https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-DEADLOCKS
- https://pkg.go.dev/time
- https://web.archive.org/web/2024/https://dev.mysql.com/doc/refman/8.0/en/innodb-deadlock-detection.html (MySQL retry context via backoff)
- https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/ (AWS backoff+jitter pattern)

Confidence: HIGH (retry disarankan PG; backoff+jitter didukung oleh AWS Architecture Blog sebagai best practice industri — Evidence 16)

### Finding 9 — Observability Production

Claim: Tanda deadlock: error `40P01` acak saat trafik tinggi, tidak reproduksibel di laptop, hilang saat retry. Pantau via `log_lock_waits = on` (log bila tunggu > deadlock_timeout), `log_min_duration_statement`, `pg_locks` (granted=false), `pg_stat_activity` (wait_event_type=Lock), dan `pg_blocking_pids()`.

Evidence: `log_lock_waits`: "Controls whether a log message is produced when a session waits longer than deadlock_timeout to acquire a lock." `pg_locks`: kolom `granted bool True if lock is held, false if lock is awaited`. `pg_stat_activity`: `wait_event_type Lock` + `wait_event` mengidentifikasi tipe lock (relation, tuple, transactionid).

Sources:
- https://www.postgresql.org/docs/current/runtime-config-logging.html
- https://www.postgresql.org/docs/current/view-pg-locks.html
- https://www.postgresql.org/docs/current/monitoring-stats.html

Confidence: HIGH

### Finding 10 — Isolation Level dan Ruang Lingkup Deadlock

Claim: Default PostgreSQL `Read Committed` sudah mencegah dirty read; updater yang menemukan row telah di-lock akan menunggu lalu re-evaluasi WHERE terhadap versi terbaru. `Repeatable Read`/`Serializable` menambah jaminan tetapi wajib handle `40001`. SSI predicate locks (`SIReadLock` di `pg_locks`) tidak memblokir, hanya mendeteksi anomali serialisasi.

Evidence: 13.2.1 Read Committed: "the would-be updater will wait... search condition is re-evaluated". 13.2.2: retry 40001. 13.2.3 SSI: "predicate locks... do not cause any blocking and therefore can not play any part in causing a deadlock."

Sources:
- https://www.postgresql.org/docs/current/transaction-iso.html

Confidence: HIGH

## Areas of Agreement

- Semua sumber primer PostgreSQL dan Wikipedia sepakat: deadlock = siklus tunggu yang memerlukan keempat Coffman conditions bersamaan; memutus circular wait via lock ordering adalah pencegahan paling efektif.
- Semua sumber sepakat: deadlock bukan resource exhaustion (bukan "tambah RAM") melainkan ordering + duration problem.
- Semua sumber sepakat: `40P01` dan `40001` wajib di-retry di level aplikasi untuk transaksi yang aman diulang.

## Areas of Disagreement

- Tidak ada kontradiksi material di antara sumber terverifikasi PostgreSQL. Perbedaan apparent antara `deadlock_timeout` vs `lock_timeout` sering tertukar di komunitas, tetapi dokumentasi resmi membedakannya secara tegas (interval deteksi vs batas tunggu).
- MySQL/InnoDB: Diverifikasi via Wayback Machine (snapshot 2024, halaman 17.7.5.2). InnoDB pakai **wait-for graph** aktif (bukan interval periodik seperti PG) dan memilih victim **transaksi kecil** (jumlah row). Deteksi otomatis dikontrol `innodb_deadlock_detect` (default ON); jika dimatikan, fallback ke `innodb_lock_wait_timeout`. Kontradiksi sebelumnya tentang MySQL tidak tersedia — sekarang tersedia via arsip.

## Limitations

- MySQL primary sources diakses via Wayback Machine (arsip resmi) — verifikasi deadlock detection dan victim selection sekarang可信 (Evidence 15).
- Paper asli Coffman 1971 (DOI 10.1145/356586.356588) tidak bisa diakses langsung (paywall ACM), tetapi Wikipedia merujuknya secara eksplisit dan menyatakan conditions berasal dari paper tersebut; Wikipedia adalah secondary source dengan confidence tinggi untuk fakta dasar.
- Metrik spesifik PPOB (throughput 500->5000 user) di topic spec adalah ilustrasi, bukan benchmark terukur — tidak ada dataset produksi yang dikutip.
- Backoff/jitter/idempotency key sekarang didukung oleh AWS Architecture Blog (Tier 1) dan primitif Go `time` (Evidence 16).
- Observability via CSV/JSON log (`log_destination`) dan `pg_stat_database` conflict stats belum dieksplor mendalam.

## Conclusion

Deadlock pada sistem PPOB bukan anomali database melainkan konsekuensi alami konkurensi: transaksi yang benar secara individual dapat membentuk siklus bila urutan lock tidak konsisten dan transaksi terlalu lama. Database mendeteksi siklus secara periodik (`deadlock_timeout` 1s) dan memilih korban (`40P01`). Aplikasi bertugas mengurangi frekuensi (lock ordering, transaksi singkat, kirim ke provider setelah COMMIT, hindari mengunci terlalu banyak row) dan menangani sisanya dengan benar (retry idempoten berbatas + backoff, bedakan dari `lock_timeout`/`55P03`, pantau `log_lock_waits`/`pg_locks`/`pg_stat_activity`). Pada isolasi lebih ketat, `40001 serialization_failure` menambah kebutuhan retry yang sama.
