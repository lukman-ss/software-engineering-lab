# Research Plan

## Research Topic

Deadlock — Kenapa Dua Transaksi yang Benar Bisa Saling Mengunci Selamanya?
(Lab #18, Senior Software Engineer Daily, kategori Concurrency/Database, studi kasus Sistem PPOB)

## Objective

Mengumpulkan bukti dari sumber otoritatif tentang: definisi deadlock, penyebab
(utamanya urutan lock berbeda dan transaksi terlalu lama), cara database mendeteksi
dan menanganinya, serta strategi aplikasi mengurangi dan menangani deadlock
(lock ordering, lock contention, deadlock retry, lock timeout).

## Research Questions

- RQ1: Apa definisi deadlock dan syarat-syarat terjadinya (Coffman conditions)?
- RQ2: Bagaimana dua transaksi yang benar bisa saling mengunci (urutan lock berbeda)?
- RQ3: Apa yang dilakukan database (PostgreSQL, MySQL/InnoDB, SQL Server) saat deadlock terjadi — deteksi, pemilihan korban, error code?
- RQ4: Bagaimana transaction duration dan lock contention memengaruhi kemungkinan deadlock?
- RQ5: Bagaimana aplikasi mengurangi dan menangani deadlock (lock ordering, retry, backoff)?
- RQ6: Apa bedanya deadlock dengan lock timeout, dan bagaimana penanganannya berbeda?
- RQ7: Bagaimana deadlock dideteksi/dipantau di production (log dan metrik)?

## Search Strategy

1. Cari dokumentasi resmi database: PostgreSQL (explicit locking, lock management),
   MySQL/InnoDB (Deadlocks in InnoDB, deadlock detection, handling), Microsoft SQL Server
   (Deadlocks Guide), Oracle.
2. Buka halaman sumber langsung (bukan snippet pencarian) dan kutip bukti.
3. Sumber akademik/standar: Coffman et al. 1971, Silberschatz Operating System Concepts,
   Bernstein et al. 2PL.
4. Sumber industri: Percona blog (MySQL deadlock diagnosis), AWS retry-backoff pattern.
5. Cross-check setiap klaim penting ke minimal 2 sumber independen.

## Expected Primary Sources

- PostgreSQL Documentation: 13.3.4 Deadlocks, 19.12 Lock Management (deadlock_timeout)
- MySQL 8.4 Reference Manual: 17.7.5.x Deadlocks in InnoDB (example, detection, handling)
- Microsoft Learn: Deadlocks Guide — SQL Server
- Oracle Database documentation (deadlock ORA-00060) — jika dapat diakses
- Coffman, Elphick, Shoshani 1971 "System Deadlocks"; Silberschatz et al. OS Concepts
- Bernstein, Hadzilacos, Goodman 1987 (2PL)
- Percona Blog: "How to deal with MySQL deadlocks" (2014)
- AWS Prescriptive Guidance: Retry and backoff pattern

## Risks / Unknowns

- dev.mysql.com memblokir fetch langsung (HTTP 403) → gunakan Wayback Machine snapshot
  dan catat sebagai sumber arsip dengan tanggal capture.
- Tidak ditemukan sumber primer spesifik sistem PPOB/fintech Indonesia → klaim konteks
  PPOB hanya ilustrasi studi kasus, bukan fakta bersumber.
- Oracle docs URL deadlock tercampur dengan MySQL URL pada hasil pencarian → verifikasi
  hati-hati; jika tidak yakin, tandai NOT VERIFIED.
- Mata penelitian: 2026-09-26. Perilaku default database bisa berubah di versi mendatang.
