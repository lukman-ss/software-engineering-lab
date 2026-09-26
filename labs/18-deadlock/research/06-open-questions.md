# Open Questions

## Unanswered Questions

1. **Berapa frekuensi deadlock "normal" vs "bermasalah" di production?**
   - Semua sumber menyebut deadlock sesekali wajar, sering = masalah — tetapi tidak ada
     yang memberi ambang kuantitatif (mis. deadlock/jam atau % transaksi). Perlu data
     operasional nyata atau benchmark.

2. **Bagaimana perilaku deadlock pada database terdistribusi / multi-node?**
   - Penelitian ini hanya mencakup single-node (PG, InnoDB, SQL Server). Deteksi deadlock
     lintas shard/node (2PC, distributed lock manager) belum diteliti.

3. **Apakah Conservative 2PL atau deadlock-prevention (wait-die / wound-wait) dipakai di
   DBMS modern?**
   - Disebut di literatur 2PL tetapi tidak ditemukan bukti pemakaian di PG/MySQL/SQL Server.
     Perlu verifikasi ke paper implementasi atau source code.

## Weak Evidence

4. **Klaim "query luas meningkatkan konflik" (Evidence 11) — MEDIUM.**
   - Berasal dari analisis kasus Percona, bukan pernyataan eksplisit vendor. Perlu contoh
     terukur (mis. UPDATE tanpa index vs dengan index, jumlah row locked).

5. **Coffman conditions dan 2PL via Wikipedia (Tier 3) — perlu verifikasi primer.**
   - Perlu buka langsung: Coffman et al. 1971 (DOI 10.1145/356586.356588), Silberschatz
     Operating System Concepts, Bernstein et al. 1987.

6. **MySQL docs via Wayback Machine (arsip, bukan live).**
   - Isi sesuai © Oracle tetapi tanggal arsip (Des 2024 / Jan 2025). Jika dev.mysql.com
     dapat diakses lagi, verifikasi ulang URL live:
     https://dev.mysql.com/doc/refman/8.4/en/innodb-deadlocks-handling.html

## Claims Needing Deeper Research

7. **Oracle ORA-00060 — NOT VERIFIED.** URL retrieval gagal; klaim "Oracle rollback satu
   statement" dikeluarkan. Perlu buka Oracle Database Concepts 19c langsung.

8. **Kasus AUTO-INC + FK deadlock (Percona) — perlu reproduksi.**
   - Root cause S-lock dari foreign key constraint menarik untuk lab PPOB (relasi
     agent↔transaksi↔riwayat). Perlu eksperimen InnoDB nyata.

9. **Optimal retry policy (max retries, base delay, jitter, circuit breaker interplay).**
   - Sumber hanya memberi pola umum (1–3 detik acak, exponential backoff). Nilai optimal
     untuk workload PPOB (latensi provider eksternal) belum ada.

10. **Metrik observability spesifik deadlock per DBMS.**
    - PG: pg_stat_database deadlocks count; MySQL: performance_schema / SHOW STATUS
      Innodb_deadlocks; SQL Server: system_health XEvent. Masing-masing disebut sepintas,
      belum dikumpulkan sistematis untuk dashboard lab.

## Possible Next Research Directions

- Eksperimen reproduksi deadlock dua transaksi (PG + MySQL) dengan skenario transfer
  A↔B seperti studi kasus lab; ukur efek lock ordering dan durasi transaksi.
- Teliti deadlock pada pola PPOB nyata: kurangi saldo agen → buat transaksi → request
  provider → tambah riwayat; petakan mana di dalam vs luar transaksi.
- Bandingkan strategi isolasi (READ COMMITTED vs REPEATABLE READ vs SERIALIZABLE/SSI)
  terhadap frekuensi deadlock pada workload contention tinggi.
- Kaji interaksi deadlock-retry dengan idempotency key dan exactly-once semantics pada
  pembayaran (risiko double-charge saat retry).
