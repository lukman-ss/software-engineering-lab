# 03-evidence.md

## Evidence 1

Claim: Deadlock terjadi ketika dua atau lebih transaksi masing-masing memegang lock yang diinginkan oleh yang lain, membentuk siklus tunggu (circular wait).
Evidence: "PostgreSQL automatically detects deadlock situations and resolves them by aborting one of the transactions involved... The best defense against deadlocks is generally to avoid them by being certain that all applications using a database acquire locks on multiple objects in a consistent order."
Source: PostgreSQL 18 Docs — 13.3.4 Deadlocks
URL: https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-DEADLOCKS
Confidence: HIGH
Corroborated By: Wikipedia Deadlock (computer science) — circular wait condition; Wikipedia Coffman conditions.

## Evidence 2

Claim: Empat Coffman conditions wajib bersamaan untuk deadlock: (1) mutual exclusion, (2) hold and wait, (3) no preemption, (4) circular wait.
Evidence: "A deadlock situation on a resource can arise only if all of the following conditions occur simultaneously in a system: 1. Mutual exclusion... 2. Hold and wait... 3. No preemption... 4. Circular wait... These four conditions are known as the Coffman conditions from their first description in a 1971 article by Edward G. Coffman, Jr."
Source: Wikipedia Deadlock (computer science) — Conditions
URL: https://en.wikipedia.org/wiki/Deadlock_(computer_science)#Conditions
Confidence: HIGH
Corroborated By: Silberschatz OS Principles (2006) p.239 via Wikipedia ref 6.

## Evidence 3

Claim: PostgreSQL error code untuk deadlock adalah SQLSTATE 40P01 (deadlock_detected).
Evidence: Table A.1 Class 40 — Transaction Rollback: 40P01 deadlock_detected.
Source: PostgreSQL 18 Docs — Appendix A Error Codes
URL: https://www.postgresql.org/docs/current/errcodes-appendix.html
Confidence: HIGH
Corroborated By: -

## Evidence 4

Claim: PostgreSQL mendeteksi deadlock secara periodik (bukan tiap wait) dan membatalkan salah satu transaksi (victim).
Evidence: "PostgreSQL automatically detects deadlock situations and resolves them by aborting one of the transactions involved... Exactly which transaction will be aborted is difficult to predict and should not be relied upon."
Source: PostgreSQL 18 Docs — 13.3.4 Deadlocks
URL: https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-DEADLOCKS
Confidence: HIGH
Corroborated By: -

## Evidence 5

Claim: deadlock_timeout (default 1s) adalah interval waktu menunggu sebelum server memeriksa adanya deadlock; bukan lock timeout.
Evidence: "This is the amount of time to wait on a lock before checking to see if there is a deadlock condition... If this value is specified without units, it is taken as milliseconds. The default is one second (1s)... On a heavily loaded server you might want to raise it."
Source: PostgreSQL 18 Docs — 19.12 Lock Management
URL: https://www.postgresql.org/docs/current/runtime-config-locks.html
Confidence: HIGH
Corroborated By: -

## Evidence 6

Claim: lock_timeout (default 0 = disabled) membatalkan statement jika menunggu lock melebihi durasi tertentu — berbeda dari deadlock.
Evidence: "Abort any statement that waits longer than the specified amount of time while attempting to acquire a lock... The time limit applies separately to each lock acquisition attempt... A value of zero (the default) disables the timeout."
Source: PostgreSQL 18 Docs — 19.11 Client Connection Defaults
URL: https://www.postgresql.org/docs/current/runtime-config-client.html
Confidence: HIGH
Corroborated By: -

## Evidence 7

Claim: Mitigasi paling efektif deadlock: konsistensi urutan lock (lock ordering) di seluruh aplikasi.
Evidence: "The best defense against deadlocks is generally to avoid them by being certain that all applications using a database acquire locks on multiple objects in a consistent order."
Source: PostgreSQL 18 Docs — 13.3.4 Deadlocks
URL: https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-DEADLOCKS
Confidence: HIGH
Corroborated By: Wikipedia Deadlock Prevention — "Approaches that avoid circular waits include ... using a hierarchy to determine a partial ordering of resources."

## Evidence 8

Claim: Transaksi panjang (memegang lock lama) memperbesar peluang deadlock.
Evidence: "So long as no deadlock situation is detected, a transaction seeking either a table-level or row-level lock will wait indefinitely for conflicting locks to be released. This means it is a bad idea for applications to hold transactions open for long periods of time (e.g., while waiting for user input)."
Source: PostgreSQL 18 Docs — 13.3.4 Deadlocks
URL: https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-DEADLOCKS
Confidence: HIGH
Corroborated By: -

## Evidence 9

Claim: Operasi eksternal (API call, PDF, S3) di dalam transaksi memperpanjang lock hold time dan meningkatkan risiko deadlock.
Evidence: Contoh anti-pattern dokumen: BEGIN → UPDATE → Call API WhatsApp → Generate PDF → Upload S3 → COMMIT — semua langkah memegang lock.
Source: Lab topic spec case study (derived)
URL: (internal)
Confidence: HIGH (inferred from PG docs Evidence 8 + general principle)
Corroborated By: -

## Evidence 10

Claim: Retry otomatis transaksi idempoten saat dapat SQLSTATE 40P01 adalah pola standar.
Evidence: "If it is not feasible to verify this in advance, then deadlocks can be handled on-the-fly by retrying transactions that abort due to deadlocks."
Source: PostgreSQL 18 Docs — 13.3.4 Deadlocks
URL: https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-DEADLOCKS
Confidence: HIGH
Corroborated By: -

## Evidence 11

Claim: Serialization_failure (SQLSTATE 40001) di Repeatable Read/Serializable memerlukan retry penuh transaksi.
Evidence: "Applications using this level must be prepared to retry transactions due to serialization failures... When an application receives this error message, it should abort the current transaction and retry the whole transaction from the beginning."
Source: PostgreSQL 18 Docs — 13.2.2 Repeatable Read Isolation Level
URL: https://www.postgresql.org/docs/current/transaction-iso.html
Confidence: HIGH
Corroborated By: -

## Evidence 12

Claim: Monitoring deadlock production: log_lock_waits = on, log_min_duration_statement, pg_stat_activity wait_event_type=Lock.
Evidence: "Controls whether a log message is produced when a session waits longer than deadlock_timeout to acquire a lock... wait_event_type Lock... wait_event will identify the type of lock awaited."
Source: PostgreSQL 18 Docs — 19.8.3 What to Log (log_lock_waits); 27.2.3 pg_stat_activity wait_event_type Lock
URL: https://www.postgresql.org/docs/current/runtime-config-logging.html; https://www.postgresql.org/docs/current/monitoring-stats.html
Confidence: HIGH
Corroborated By: -

## Evidence 13

Claim: MVCC PostgreSQL membuat read tidak memblokir write; deadlock hanya melibatkan writer/explicit locker.
Evidence: "The main advantage of using the MVCC model... locks acquired for querying (reading) data do not conflict with locks acquired for writing data, and so reading never blocks writing and writing never blocks reading."
Source: PostgreSQL 18 Docs — 13.1 Introduction (MVCC)
URL: https://www.postgresql.org/docs/current/mvcc-intro.html
Confidence: HIGH
Corroborated By: -

## Evidence 14

Claim: Go time package menyediakan After, Sleep, Ticker untuk backoff; Since/Until untuk durasi monotonic.
Evidence: "After waits for the duration to elapse... Sleep pauses the current goroutine for at least the duration... time.Since(start) robust against wall clock resets."
Source: Go pkg time documentation
URL: https://pkg.go.dev/time
Confidence: HIGH
Corroborated By: -

## Evidence 15

Claim: MySQL/InnoDB mendeteksi deadlock aktif via wait-for graph, kemudian me-rollback satu atau lebih transaksi untuk memutus siklus (memilih transaksi kecil berdasarkan jumlah row); `innodb_deadlock_detect` (default enabled) bisa dimatikan dan mengandalkan `innodb_lock_wait_timeout`.
Evidence: "When deadlock detection is enabled (the default), InnoDB automatically detects transaction deadlocks and rolls back a transaction or transactions to break the deadlock. InnoDB tries to pick small transactions to roll back, where the size of a transaction is determined by the number of rows inserted, updated, or deleted."
Source: MySQL 8.0 Reference Manual — 17.7.5.2 Deadlock Detection (Wayback Machine snapshot 2024, konten identik dengan dev.mysql.com)
URL: https://web.archive.org/web/2024/https://dev.mysql.com/doc/refman/8.0/en/innodb-deadlock-detection.html
Confidence: HIGH (kembali diverifikasi 2026-09-26)
Corroborated By: Struktur halaman juga mencatat "TOO DEEP OR LONG SEARCH IN THE LOCK TABLE WAITS-FOR GRAPH" saat wait-for list >200 transaksi → dipakai sebagai batas praktis deteksi siklus.

## Evidence 16

Claim: Exponential backoff + jitter adalah pola retry yang distandarkan vendor cloud untuk sistem tahan-gagal, bukan sekadar improvisasi.
Evidence: "Most AWS SDKs now support exponential backoff and jitter as part of their retry behavior when using standard or adaptive modes." Pola capped exponential backoff diperkenalkan untuk mengurangi sinkronisasi retry (herding) pada kontensi tinggi.
Source: AWS Architecture Blog — "Exponential Backoff And Jitter" (Marc Brooker, 04 Mar 2015; update Mei 2023)
URL: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
Confidence: HIGH
Corroborated By: AWS Builders' Library (timeouts, retries, and backoff with jitter) ditautkan di artikel yang sama. Untuk pola spesifik transaksi database, PG 13.3.4 merekomendasikan retry on-the-fly (Evidence 10).