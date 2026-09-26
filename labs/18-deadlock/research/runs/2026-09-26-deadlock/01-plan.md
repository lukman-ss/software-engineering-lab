# 01-plan.md

## Research Topic

Deadlock pada sistem PPOB (Pembayaran Primer Otomatis Berbasis) - bagaimana transaksi yang berjalan bersamaan dapat saling mengunci selamanya.

## Objective

Memahami bagaimana transaksi database yang berjalan bersamaan dapat menciptakan deadlock, bagaimana database mendeteksi dan mengatasi deadlock, serta bagaimana aplikasi dapat mengurangi risiko deadlock melalui perancangan transaksi yang tepat.

## Research Questions

1. Apa saja syarat (Coffman conditions) yang harus ada untuk terbentuknya deadlock?
2. Bagaimana PostgreSQL mendeteksi deadlock secara otomatis?
3. Bagaimana MySQL/InnoDB mendeteksi dan mengatasi deadlock?
4. Apa error code PostgreSQL yang ditunjukkan saat deadlock terjadi?
5. Bagaimana cara mencegah deadlock melalui lock ordering?
6. Strategi retry deadlock dengan exponential backoff di Go/C++?
7. Apa metrik dan logging yang harus dipantau untuk deteksi deadlock di production?
8. Perbedaan antara deadlock dan lock timeout.

## Search Strategy

1. Baca dokumentasi resmi PostgreSQL tentang Explicit Locking dan Deadlock Detection
2. Baca dokumentasi resmi MySQL/InnoDB tentang deadlock handling
3. Referensi academic: Coffman conditions (1971)
4. Studi kasus PPOB: analisis flow transfer antar tabel akun, transaksi, dan riwayat

## Expected Primary Sources

- PostgreSQL Documentation: 13.3 Explicit Locking, 13.3.4 Deadlocks
- PostgreSQL Documentation: 19.12 Lock Management
- PostgreSQL Documentation: Appendix A Error Codes (SQLSTATE 40P01)
- PostgreSQL Documentation: 53.13 pg_locks view
- MySQL Documentation: InnoDB Deadlock Detection (dinamisakan 403 - fallback ke alternative)
- Wikipedia: Deadlock (computer science) - Coffman conditions
- Perancangan aplikasi: retry pattern dengan exponential backoff

## Risks / Unknowns

- MySQL docs resmi tidak dapat diakses (403) - perlu alternatif MariaDB atau sumber lain
- Beberapa konfigurasi PostgreSQL mungkin berbeda antar versi