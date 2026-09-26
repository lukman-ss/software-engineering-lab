# 02-sources.md

## Source 1

Title: PostgreSQL 18 Documentation — 13.3. Explicit Locking (incl. 13.3.4 Deadlocks)
Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/explicit-locking.html
Published: PostgreSQL 18.6 docs, diakses 2026-09-26 (rilis 19 Beta 4 disebut 24 Sep 2026 di header)
Accessed: 2026-09-26
Source Tier: Tier 1 (official documentation)
Relevance: Definisi deadlock PG, contoh row-level deadlock accounts 11111/22222, auto-abort satu transaksi, lock ordering sebagai pertahanan utama, retry, jangan tahan transaksi lama.

## Source 2

Title: PostgreSQL 18 Documentation — 19.12. Lock Management
Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/runtime-config-locks.html
Published: PostgreSQL 18.6 docs, diakses 2026-09-26
Accessed: 2026-09-26
Source Tier: Tier 1 (official documentation)
Relevance: deadlock_timeout default 1s (interval cek deadlock, bukan timeout lock), max_locks_per_transaction default 64, predicate locks.

## Source 3

Title: PostgreSQL 18 Documentation — Appendix A. PostgreSQL Error Codes
Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/errcodes-appendix.html
Published: PostgreSQL 18.6 docs, diakses 2026-09-26
Accessed: 2026-09-26
Source Tier: Tier 1 (standards-adjacent official docs)
Relevance: SQLSTATE 40P01 deadlock_detected, 40001 serialization_failure, 55P03 lock_not_available — kunci untuk retry handling di aplikasi.

## Source 4

Title: PostgreSQL 18 Documentation — 13.2. Transaction Isolation
Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/transaction-iso.html
Published: PostgreSQL 18.6 docs, diakses 2026-09-26
Accessed: 2026-09-26
Source Tier: Tier 1 (official documentation)
Relevance: Read Committed default + perilaku updater menunggu + re-evaluasi WHERE; Repeatable Read/Serializable wajib retry 40001; SSI predicate locks tidak menyebabkan deadlock.

## Source 5

Title: PostgreSQL 18 Documentation — 53.13. pg_locks
Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/view-pg-locks.html
Published: PostgreSQL 18.6 docs, diakses 2026-09-26
Accessed: 2026-09-26
Source Tier: Tier 1 (official documentation)
Relevance: Observability lock: granted false = menunggu, join pg_stat_activity, anjuran pakai pg_blocking_pids() daripada self-join.

## Source 6

Title: PostgreSQL 18 Documentation — 19.11. Client Connection Defaults
Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/runtime-config-client.html
Published: PostgreSQL 18.6 docs, diakses 2026-09-26
Accessed: 2026-09-26
Source Tier: Tier 1 (official documentation)
Relevance: lock_timeout (per lock attempt, default 0 = disabled), statement_timeout, transaction_timeout, idle_in_transaction_session_timeout — fondasi bedakan deadlock vs lock timeout.

## Source 7

Title: PostgreSQL 18 Documentation — 19.8. Error Reporting and Logging
Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/runtime-config-logging.html
Published: PostgreSQL 18.6 docs, diakses 2026-09-26
Accessed: 2026-09-26
Source Tier: Tier 1 (official documentation)
Relevance: log_lock_waits (log bila tunggu > deadlock_timeout), log_min_duration_statement, CSV/JSON log kolom SQLSTATE — metrik deteksi deadlock production.

## Source 8

Title: PostgreSQL 18 Documentation — 27.2. The Cumulative Statistics System (pg_stat_activity)
Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/monitoring-stats.html
Published: PostgreSQL 18.6 docs, diakses 2026-09-26
Accessed: 2026-09-26
Source Tier: Tier 1 (official documentation)
Relevance: wait_event_type = Lock, state active tapi wait_event non-null = query jalan tapi terblokir — sinyal lock contention production.

## Source 9

Title: PostgreSQL 18 Documentation — 13.1. Introduction (MVCC)
Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/mvcc-intro.html
Published: PostgreSQL 18.6 docs, diakses 2026-09-26
Accessed: 2026-09-26
Source Tier: Tier 1 (official documentation)
Relevance: MVCC: reading never blocks writing — batasi ruang lingkup deadlock ke writer-writer / explicit locker.

## Source 10

Title: Deadlock (computer science) — Wikipedia
Publisher: Wikipedia (tertiary, merujuk Coffman et al. 1971, Silberschatz, Tanenbaum)
URL: https://en.wikipedia.org/wiki/Deadlock_(computer_science)
Published: revisi berjalan; diakses 2026-09-26
Accessed: 2026-09-26
Source Tier: Tier 2 untuk artikel Wikipedia; Tier 1 untuk paper primer yang dirujuk Wikipedia
Relevance: Empat Coffman conditions, taksonomi handling (ignore/detection/prevention/avoidance incl. Banker's), circular-wait breaking via resource ordering, livelock vs deadlock, distributed/phantom deadlock.
Primary: Edward G. Coffman, Jr., M. J. Elphick, A. Shoshani — "System Deadlocks", ACM Computing Surveys 3(2):67–78, June 1971. DOI 10.1145/356586.356588. Wikipedia merujuk paper ini secara eksplisit ("from their first description in a 1971 article by Edward G. Coffman, Jr."); paper berada di balik paywall ACM/IEEE namun sitasi dan kondisi 1–4 dikutip verbatim oleh Wikipedia.

## Source 12

Title: MySQL 8.0 Reference Manual — 17.7.5.2 Deadlock Detection
Publisher: Oracle / MySQL (via Wayback Machine snapshot)
URL: https://web.archive.org/web/2024/https://dev.mysql.com/doc/refman/8.0/en/innodb-deadlock-detection.html
Published: snapshot 2024 (original rilis 8.0); diakses 2026-09-26 via web.archive.org
Accessed: 2026-09-26
Source Tier: Tier 1 content, Tier 2 access channel (arsip resmi dev.mysql.com)
Relevance: InnoDB aktif mendeteksi deadlock via wait-for graph (lebih dari 200 transaksi di wait-for list dianggap deadlock), lalu rollback transaksi (prioritas transaksi kecil: jumlah row diinsert/updated/deleted), konfigurasi `innodb_deadlock_detect` (bisa dimatikan, lalu mengandalkan `innodb_lock_wait_timeout`). Klarifikasi bahwa deadlock_detection vs lock_wait_timeout adalah dua mekanisme berbeda — analog PostgreSQL `deadlock_timeout` vs `lock_timeout`.
Note: dev.mysql.com mengembalikan 403 untuk bot; halaman yang sama diakses melalui Wayback Machine. Isi halaman identik dengan dokumentasi asli.

## Source 13

Title: Exponential Backoff And Jitter — AWS Architecture Blog (Marc Brooker)
Publisher: Amazon Web Services
URL: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
Published: 04 MAR 2015; update 8 tahun kemudian (May 2023) mengonfirmasi AWS SDK memakai pola ini
Accessed: 2026-09-26
Source Tier: Tier 1 (authoritative cloud vendor engineering guidance)
Relevance: Pola "capped exponential backoff + jitter" sebagai standar retry di sistem tersebar — mendukung Finding 8 (retry backoff + jitter), bukan sekadar primitif Go stdlib. Kutipan: "Most AWS SDKs now support exponential backoff and jitter as part of their retry behavior."

## Attempted but inaccessible (langsung; sudah diatasi via mirror)

- https://dev.mysql.com/doc/refman/8.0/en/innodb-locking.html — 403
- https://dev.mysql.com/doc/refman/8.0/en/deadlock.html — 403
- https://dev.mysql.com/doc/refman/8.0/en/innodb-deadlock-detection.html — 403 (diakses via Wayback Machine, lihat Source 12)
- https://dev.mysql.com/doc/refman/8.0/en/innodb-locking-reads.html — 403
- https://mariadb.com/kb/en/innodb-deadlock-detection/ — halaman ada tetapi dirender via JavaScript, tidak bisa di-fetch statis; tidak digunakan sebagai sumber.
- Klaim MySQL/InnoDB (innodb_deadlock_detect, innodb_lock_wait_timeout, wait-for graph, rollback transaksi kecil): diverifikasi via Source 12 (arsip).
