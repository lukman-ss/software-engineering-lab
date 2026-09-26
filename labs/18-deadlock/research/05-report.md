# Research Report: Deadlock — Kenapa Dua Transaksi yang Benar Bisa Saling Mengunci Selamanya?

## Research Question

Bagaimana transaksi database yang masing-masing benar secara individual dapat saling mengunci selamanya (deadlock), bagaimana database (PostgreSQL, MySQL/InnoDB, SQL Server) mendeteksi dan menanganinya, serta bagaimana aplikasi mengurangi frekuensi dan menangani dampaknya (lock ordering, transaction duration, lock contention, retry, lock timeout)? Studi kasus: Sistem PPOB.

## Executive Summary

Deadlock bukan bug database, melainkan konsekuensi cara aplikasi mengakses data bersamaan. Dua transaksi yang masing-masing benar dapat deadlock ketika mengunci resource dalam urutan berbeda sehingga terjadi circular wait — masing-masing memegang lock yang dibutuhkan pihak lain. Semua DBMS modern (PostgreSQL, MySQL/InnoDB, SQL Server) mendeteksi kondisi ini otomatis dan mem-batalkan satu transaksi (deadlock victim) agar yang lain lanjut; korban menerima error (PG 40P01, MySQL 1213, SQL Server 1205). Penyebab paling umum: urutan lock tidak konsisten dan transaksi terlalu lama. Defense paling efektif menurut ketiga vendor: lock ordering konsisten + transaksi singkat + retry di application-level dengan backoff. Lock timeout berbeda dari deadlock: timeout = menunggu terlalu lama tanpa siklus; deadlock = siklus saling tunggu. Research date: 2026-09-26.

## Findings

### Finding 1 — Definisi: deadlock = circular wait (4 syarat Coffman)

**Claim:** Deadlock terjadi hanya jika empat syarat Coffman terpenuhi bersamaan: mutual exclusion, hold-and-wait, no preemption, circular wait.

**Evidence:** "A deadlock situation on a resource can arise only if all of the following conditions occur simultaneously: mutual exclusion, hold and wait, no preemption, circular wait." Contoh konkret dua transaksi UPDATE accounts (11111 vs 22222 dalam urutan terbalik) menghasilkan "transaction one is blocked on transaction two, and transaction two is blocked on transaction one: a deadlock condition."

**Sources:**
- Coffman et al. 1971, via https://en.wikipedia.org/wiki/Deadlock_(computer_science)
- PostgreSQL 18 Docs §13.3.4, https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-DEADLOCKS (2026-09-24)

**Confidence:** HIGH

### Finding 2 — Dua transaksi benar bisa deadlock karena urutan lock berbeda

**Claim:** Akar penyebab paling umum: dua transaksi menyentuh set data sama dalam urutan berbeda (A: lock 1→2, B: lock 2→1).

**Evidence:** PostgreSQL: jika kedua transaksi update baris dalam urutan sama, "no deadlock would have occurred." MySQL: Client A lock Animals lalu minta Birds; Client B lock Birds lalu minta Animals → "ERROR 1213 (40001): Deadlock found when trying to get lock; try restarting transaction", InnoDB rollback transaction (2). SQL Server: contoh T1 lock Supplier lalu Part vs T2 sebaliknya membentuk cycle.

**Sources:**
- PostgreSQL §13.3.4 (URL di atas)
- MySQL 8.4 §17.7.5.1, https://web.archive.org/web/20241214192732/https://dev.mysql.com/doc/refman/8.4/en/innodb-deadlock-example.html (©2024 Oracle)
- Microsoft Learn Deadlocks Guide, https://learn.microsoft.com/en-us/sql/relational-databases/sql-server-deadlocks-guide?view=sql-server-ver17 (2025-10-09, upd. 2026-09-21)

**Confidence:** HIGH

### Finding 3 — Database deteksi otomatis dan korbankan satu transaksi

**Claim:** PostgreSQL, InnoDB, SQL Server semuanya deteksi otomatis (wait-for graph / lock monitor) dan abort satu transaksi agar siklus pecah; transaksi mana yang dikorbankan tidak boleh diandalkan aplikasi.

**Evidence:** PG: "PostgreSQL automatically detects deadlock situations and resolves them by aborting one of the transactions involved... Exactly which transaction will be aborted is difficult to predict and should not be relied upon." InnoDB: "automatically detects transaction deadlocks and rolls back a transaction... tries to pick small transactions to roll back, where size = number of rows inserted/updated/deleted." Batas wait-for list 200 transaksi. SQL Server: deadlock monitor default tiap 5 detik (turun ke 100ms bila sering), pilih korban termurah di-rollback atau prioritas terendah, kembalikan error 1205.

**Sources:**
- PostgreSQL §13.3.4 (URL di atas)
- MySQL 8.4 §17.7.5.2, https://web.archive.org/web/20250126061946/https://dev.mysql.com/doc/refman/8.4/en/innodb-deadlock-detection.html (©2025 Oracle)
- Microsoft Learn Deadlocks Guide (URL di atas)

**Confidence:** HIGH

### Finding 4 — Transaction duration panjang meningkatkan risiko deadlock

**Claim:** Semakin lama transaksi menahan lock, semakin besar peluang deadlock; operasi eksternal (API, PDF, S3, input user) di dalam transaksi memperburuk.

**Evidence:** PG: "it is a bad idea for applications to hold transactions open for long periods of time (e.g., while waiting for user input)." MS: "A deadlock typically occurs when several long-running transactions execute concurrently... The longer the transaction, the longer the exclusive or update locks are held." MySQL: "Keep transactions small and short in duration... Commit transactions immediately... do not leave an interactive mysql session open for a long time with an uncommitted transaction." Percona: "splitting a long transaction into smaller ones, so locks are released sooner."

**Sources:**
- PostgreSQL §13.3.4; MySQL §17.7.5.3 (https://web.archive.org/web/20241201211118/https://dev.mysql.com/doc/refman/8.4/en/innodb-deadlocks-handling.html); Microsoft Learn; Percona https://www.percona.com/blog/how-to-deal-with-mysql-deadlocks/ (2014-10-28)

**Confidence:** HIGH

### Finding 5 — Lock ordering konsisten adalah defense utama

**Claim:** Semua aplikasi harus acquire lock multi-objek dalam urutan konsisten; ini mengubah potensi deadlock menjadi lock wait biasa (antrean terdefinisi).

**Evidence:** PG: "The best defense against deadlocks is generally to avoid them by being certain that all applications using a database acquire locks on multiple objects in a consistent order." MySQL: "do those operations in a consistent order each time. Then transactions form well-defined queues and do not deadlock." MS: "Access objects in the same order... deadlocks are less likely to occur." Percona: "change them to access data in the same order, in another word, serialize the access."

**Sources:** Ketiga vendor + Percona (URL di Finding 4)

**Confidence:** HIGH

### Finding 6 — Lock contention: kurangi via index, isolasi lebih rendah, row versioning

**Claim:** Query skop luas tanpa index, gap locking pada isolasi tinggi, dan locking reads berlebihan memperbesar contention; mitigasi: index tepat, READ COMMITTED / row versioning, hindari SERIALIZABLE bila tak perlu.

**Evidence:** MySQL: "Add well-chosen indexes so queries scan fewer index records and set fewer locks"; "try using a lower isolation level such as READ COMMITTED" untuk locking reads. MS: "Enable READ_COMMITTED_SNAPSHOT... Use snapshot isolation... Avoid higher isolation levels such as REPEATABLE READ and SERIALIZABLE." Percona: "adding indexes to minimize the rows scanned and locked"; "change isolation level to read committed to avoid gap locking." PG SSI: predicate locks "do not cause any blocking and therefore can not play any part in causing a deadlock."

**Sources:** MySQL §17.7.5.3; Microsoft Learn; Percona; PG §13.2.3 (https://www.postgresql.org/docs/current/transaction-iso.html#XACT-SERIALIZABLE)

**Confidence:** HIGH

### Finding 7 — Deadlock retry: selalu siap retry, dengan backoff + jitter, hanya untuk operasi idempotent

**Claim:** Aplikasi harus catch error deadlock dan retry transaksi; jeda singkat acak (backoff+jitter) mencegah deadlock terulang; retry hanya aman bila operasi idempotent.

**Evidence:** MySQL: "you must write your applications so that they are always prepared to re-issue a transaction if it gets rolled back because of a deadlock." MS: "Implementing an error handler that catches error 1205... automatically resubmitting... pause briefly before resubmitting... Randomizing the duration of the pause minimizes likelihood of deadlock reoccurring." Percona: "have the applications catch deadlock error (MySQL error no. 1213) and handles it by retrying." AWS: "Exponential backoff... Operations should be idempotent when you use the retry with backoff pattern."

**Sources:** MySQL §17.7.5.3; Microsoft Learn; Percona; AWS https://docs.aws.amazon.com/prescriptive-guidance/latest/cloud-design-patterns/retry-backoff.html

**Confidence:** HIGH

### Finding 8 — Deadlock vs lock timeout: mekanisme dan solusi berbeda

**Claim:** Deadlock = siklus saling tunggu → DB deteksi dan batalkan satu pihak. Lock timeout = satu pihak menunggu terlalu lama tanpa siklus → ekspirasi batas waktu.

**Evidence:** MS: "Deadlocking is often confused with normal blocking... The requesting transaction is blocked, not deadlocked... Eventually the owning transaction completes... Deadlocks are resolved almost immediately, whereas blocking can persist indefinitely." PG: "So long as no deadlock situation is detected, a transaction seeking a lock will wait indefinitely." MySQL: `innodb_lock_wait_timeout` (default 50s) sebagai fallback bila `innodb_deadlock_detect` dimatikan; PG `deadlock_timeout` default 1s hanya jeda sebelum cek deadlock, bukan timeout lock umum.

**Sources:** Microsoft Learn; PG §13.3.4 + §19.12 (https://www.postgresql.org/docs/current/runtime-config-locks.html#GUC-DEADLOCK-TIMEOUT); MySQL §17.7.5.2

**Confidence:** HIGH

### Finding 9 — Deteksi di production: log DB + pola error acak saat trafik tinggi

**Claim:** Tanda deadlock: error deadlock di log DB, transaksi gagal acak saat trafik tinggi, tak terreproduksi di laptop, hilang saat retry. Alat: SHOW ENGINE INNODB STATUS / innodb_print_all_deadlocks (MySQL), pg_locks + log_lock_waits (PG), xml_deadlock_report / system_health (SQL Server).

**Evidence:** MySQL: "At any time, issue SHOW ENGINE INNODB STATUS to determine the cause of the most recent deadlock"; "enabling innodb_print_all_deadlocks... Information about each deadlock... is recorded in the MySQL error log." Percona: timestamp deadlock di INNODB STATUS dicocokkan dengan application log; thread id + user/host mengidentifikasi area aplikasi; binlog/slow log/performance_schema untuk statement sebelumnya. PG: "pg_locks system view"; deadlock_timeout juga mengatur kapan pesan log_lock_waits keluar.

**Sources:** MySQL §17.7.5.3; Percona; PG §19.12 + monitoring docs

**Confidence:** HIGH

## Areas of Agreement

Semua sumber otoritatif sepakat: definisi circular wait; urutan lock berbeda sebagai penyebab utama; transaksi panjang memperbesar risiko; deteksi otomatis + korban satu transaksi; defense = urutan konsisten + transaksi singkat + retry aplikasi; perbedaan deadlock vs lock timeout. Tidak ada kontradiksi material (lihat 04-contradictions.md).

## Areas of Disagreement

Tidak ada ketidaksepakatan substansial. Hanya variasi implementasi antar vendor (bukan kontradiksi): pemilihan korban (PG: tak terprediksi; InnoDB: terkecil; SQL Server: termurah/prioritas), tuning deteksi (PG deadlock_timeout 1s; InnoDB innodb_deadlock_detect + limit 200; SQL Server monitor 5s→100ms), error code (PG 40P01; MySQL 1213; SQL Server 1205), default lock timeout (PG: wait indefinite; MySQL: 50s; SQL Server: LOCK_TIMEOUT).

## Limitations

- dev.mysql.com memblokir fetch langsung (403); bukti MySQL diambil dari arsip Wayback Machine (capture 2024-12-01, 2024-12-14, 2025-01-26) — isi sesuai dokumentasi resmi Oracle namun tanggal akses arsip, bukan live.
- Klaim Coffman conditions dan 2PL diringkas via Wikipedia (Tier 3) yang mengutip sumber primer (Coffman 1971, Bernstein 1987); makalah/buku asli tidak dibuka langsung.
- Tidak ditemukan sumber primer kasus PPOB/fintech Indonesia; pemetaan ke PPOB pada laporan ini adalah interpretasi, bukan fakta bersumber.
- Klaim Oracle ORA-00060 tidak terverifikasi (URL gagal dibuka) dan dikeluarkan dari evidence.
- Research date 2026-09-26; default konfigurasi dapat berubah di versi DB mendatang.

## Conclusion

Dua transaksi yang benar bisa deadlock karena mengunci dalam urutan berbeda sehingga membentuk circular wait; database modern mendeteksinya otomatis dan mengorbankan satu transaksi. Risiko naik seiring durasi transaksi dan luasnya lock. Cara paling efektif: urutan lock konsisten di semua kode, transaksi sesingkat mungkin (operasi eksternal setelah commit), kurangi contention (index, isolasi tepat), dan tangani sisa deadlock yang tak terhindarkan dengan retry ber-backoff untuk operasi idempotent. Deadlock tidak bisa dihilangkan 100% pada concurrency tinggi — tujuannya mengurangi frekuensi dan menangani dampak dengan benar, dipantau via log dan metrik deadlock.
