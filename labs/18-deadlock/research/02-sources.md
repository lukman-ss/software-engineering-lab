# Sources

## Source 1
**Title:** PostgreSQL 18 Documentation — Chapter 13: Concurrency Control, §13.3.4 Deadlocks
**Publisher:** The PostgreSQL Global Development Group
**URL:** https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-DEADLOCKS
**Published:** 2026-09-24 (PostgreSQL 18.6 release)
**Accessed:** 2026-09-26
**Source Tier:** Tier 1 — official documentation
**Relevance:** Primary reference untuk definisi, contoh konkret, deteksi otomatis, dan saran defense PostgreSQL.

## Source 2
**Title:** PostgreSQL 18 Documentation — Chapter 19: Server Configuration, §19.12 Lock Management (deadlock_timeout)
**Publisher:** The PostgreSQL Global Development Group
**URL:** https://www.postgresql.org/docs/current/runtime-config-locks.html#GUC-DEADLOCK-TIMEOUT
**Published:** 2026-09-24
**Accessed:** 2026-09-26
**Source Tier:** Tier 1
**Relevance:** Konfigurasi deteksi deadlock, default 1s, implikasi untuk transaction duration.

## Source 3
**Title:** MySQL 8.4 Reference Manual — §17.7.5.1 An InnoDB Deadlock Example
**Publisher:** Oracle (via Wayback Machine)
**URL:** https://web.archive.org/web/20241214192732/https://dev.mysql.com/doc/refman/8.4/en/innodb-deadlock-example.html
**Archived:** 2024-12-14 (original © 2024 Oracle)
**Accessed:** 2026-09-26
**Source Tier:** Tier 1 — official documentation (archived)
**Relevance:** Contoh dua transaksi klien A dan B mengunci dalam urutan terbalik; output SHOW ENGINE INNODB STATUS; InnoDB me-rollback transaksi (2).

## Source 4
**Title:** MySQL 8.4 Reference Manual — §17.7.5.2 Deadlock Detection
**Publisher:** Oracle (via Wayback Machine)
**URL:** https://web.archive.org/web/20250126061946/https://dev.mysql.com/doc/refman/8.4/en/innodb-deadlock-detection.html
**Archived:** 2025-01-26 (original © 2025 Oracle)
**Accessed:** 2026-09-26
**Source Tier:** Tier 1 — official documentation (archived)
**Relevance:** Wait-for graph, pembatasan 200 transaksi, rollback transaksi terkecil, innodb_deadlock_detect toggle, lock_wait_timeout fallback.

## Source 5
**Title:** MySQL 8.4 Reference Manual — §17.7.5.3 How to Minimize and Handle Deadlocks
**Publisher:** Oracle (via Wayback Machine)
**URL:** https://web.archive.org/web/20241201211118/https://dev.mysql.com/doc/refman/8.4/en/innodb-deadlocks-handling.html
**Archived:** 2024-12-01 (original © 2024 Oracle)
**Accessed:** 2026-09-26
**Source Tier:** Tier 1 — official documentation (archived)
**Relevance:** Teknik minimalisasi: urutan konsisten, transaksi singkat, commit segera, index, isolasi READ COMMITTED, LOCK TABLES serialisasi, retry application-level.

## Source 6
**Title:** Microsoft Learn — Deadlocks Guide (SQL Server / Azure SQL Database)
**Publisher:** Microsoft
**URL:** https://learn.microsoft.com/en-us/sql/relational-databases/sql-server-deadlocks-guide?view=sql-server-ver17
**Published:** 2025-10-09; updated 2026-09-21
**Accessed:** 2026-09-26
**Source Tier:** Tier 1 — official vendor documentation
**Relevance:** Definisi deadlock vs blocking, resource types (locks, worker threads, memory, parallel query, MARS), deadlock monitor interval (5s→100ms), victim selection (cost, priority), error 1205, TRY...CATCH handling, minimize techniques (same order, short tx, row versioning), contoh deadlock.

## Source 7
**Title:** How to deal with MySQL deadlocks
**Publisher:** Percona Blog
**Author:** Peiran Song
**URL:** https://www.percona.com/blog/how-to-deal-with-mysql-deadlocks/
**Published:** 2014-10-28
**Accessed:** 2026-09-26
**Source Tier:** Tier 2 — reputable technical publication (database expert vendor)
**Relevance:** Diagnosa praktis SHOW ENGINE INNODB STATUS, log binlog/slow log, kasus AUTO-INC + FK deadlock, kasus opposite-order UPDATE, saran split transaction, serialize access, index, isolation level.

## Source 8
**Title:** Coffman, E.G.; Elphick, M.J.; Shoshani, A. — "System Deadlocks"
**Publisher:** ACM Computing Surveys
**URL:** https://doi.org/10.1145/356586.356588
**Published:** 1971-06
**Accessed:** 2026-09-26 (via Wikipedia citation)
**Source Tier:** Tier 1 — academic paper (original source)
**Relevance:** 4 Coffman conditions (mutual exclusion, hold-and-wait, no preemption, circular wait) sebagai syarat perlu dan cukup deadlock.

## Source 9
**Title:** Wikipedia — Deadlock (computer science)
**Publisher:** Wikimedia Foundation
**URL:** https://en.wikipedia.org/wiki/Deadlock_(computer_science)
**Published:** Continuously updated
**Accessed:** 2026-09-26
**Source Tier:** Tier 3 — tertiary summary
**Relevance:** Ringkasan Coffman conditions dan referensi ke Silberschatz OS Concepts; digunakan untuk konteks sejarah, bukan bukti primer.

## Source 10
**Title:** Wikipedia — Two-phase locking
**Publisher:** Wikimedia Foundation
**URL:** https://en.wikipedia.org/wiki/Two-phase_locking
**Published:** Continuously updated
**Accessed:** 2026-09-26
**Source Tier:** Tier 3 — tertiary summary
**Relevance:** Ringkasan 2PL, expanding/shrinking phase, keterkaitan dengan deadlock, Conservative 2PL; rujukan primer: Bernstein et al. 1987, Weikum & Vossen 2001.

## Source 11
**Title:** AWS Prescriptive Guidance — Retry and Backoff Pattern
**Publisher:** Amazon Web Services
**URL:** https://docs.aws.amazon.com/prescriptive-guidance/latest/cloud-design-patterns/retry-backoff.html
**Published:** Undated (accessed 2026-09-26)
**Accessed:** 2026-09-26
**Source Tier:** Tier 2 — cloud provider best practice
**Relevance:** Exponential backoff, jitter, idempotency requirement untuk retry aman.

## Source 12
**Title:** Oracle Database Concepts 19c — Data Concurrency and Consistency (deadlock handling)
**Publisher:** Oracle
**URL:** (NOT VERIFIED — URL retrieval failed / mixed with MySQL docs)
**Accessed:** 2026-09-26
**Source Tier:** (Tidak terverifikasi)
**Relevance:** Dihapus dari evidence karena tidak berhasil dibuka sumber aslinya.