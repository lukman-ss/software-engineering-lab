# 06-open-questions.md

## Unanswered Questions

1. ~~**MySQL/InnoDB deadlock detection mechanism**~~ [RESOLVED 2026-09-26] Diverifikasi via Wayback Machine snapshot: `innodb_deadlock_detect` (default on), wait-for graph aktif, rollback transaksi kecil (berdasarkan jumlah row), `innodb_lock_wait_timeout` sebagai fallback bila deteksi dimatikan. Error code 1213/1205 belum diekstraksi eksplisit dari arsip — catatan: di luar cakupan utama (PostgreSQL).

2. **Backoff + jitter parameter optimal untuk PPOB**: berapa `initial_backoff`, `max_attempts`, `max_backoff` yang tepat untuk latency provider eksternal PPOB (WhatsApp/API)? butuh data empiris / SLO dari production.

3. **Idempotency key vs transaction retry trade-off**: apakah retry `TransferNaive` setelah deadlock selalu aman jika tiap UPDATE idempotent (balance += delta)? Perlu analisis double-spend edge case ketika transaction commit namun response hilang (network partition).

4. **deadlock_timeout tuning heuristik**: rekomendasi PG adalah `>` typical transaction time, tapi berapa nilai optimal untuk PPOB dengan transaksi sub-second? butuh workload benchmark.

5. **SHOW ENGINE INNODB STATUS output format**: detail field `latest detected deadlock` (trx id, wait in sec, table record lock) — butuh sumber MySQL primer untuk reproducibility.

6. **Distributed deadlock di sistem multi-shard/multi-region PPOB**: wait-for graph global, phantom deadlock. butuh literatur tambahan (Chandy-Misra-Haas) — belum di-fetch.

## Weak Evidence

- ~~Claim "sistem PPOB dengan 500 -> 5000 user" (topic spec)~~ — tetap ilustrasi, tidak ada dataset produksi yang dikutip.
- ~~Claim MySQL `innodb_deadlock_detect` dan algoritma weight-based victim selection~~ — **RESOLVED** via Wayback Machine (Evidence 15 di 03-evidence.md).

## Claims Needing Deeper Research

- ~~Apakah `lock_timeout` di MySQL semantik sama dengan PostgreSQL?~~ — MySQL menggunakan `innodb_lock_wait_timeout` (detik, default 50s). Perbedaan disinggung: PG `lock_timeout` per-lock-attempt; MySQL `innodb_lock_wait_timeout` per-transaksi. Tidak kontradiksi, melainkan perbedaan konfigurasi.
- ~~Implementasi retry di Go: paket error resmi `database/sql` tidak wrap SQLSTATE `40P01`~~ — ditambahkan sumber AWS (Evidence 16); pola backoff+jitter sekarang HIGH confidence.

## Possible Next Research Directions

A. Cross-check MySQL deadlock behavior via Wayback Machine / dokumen versi lama sebagai Tier-2.
B. Fetch Coffman et al. "System Deadlocks" (1971) PDF via archive atau unpaywall untuk memastiksi klaim Coffman conditions (naikkan dari Tier-2 -> Tier-1).
C. Fetch ISO/IEC 9075-2 SQL standard bagian deadlock semantics bila perlu cross-reference.
D. Survey production post-mortem (Stripe, Shopify) untuk pola retry deadlock di payments system — naikkan confidence Finding 8 dari MEDIUM ke HIGH.
