# 04-contradictions.md

## Contradiction 1

**Perbedaan definisi deadlock vs lock timeout (MySQL docs tidak tersedia untuk konfirmasi)**

### SOURCE A
PostgreSQL: deadlock adalah kondisi siklus tunggu (A menunggu B, B menunggu A). Database mendeteksi siklus dan membatalk salah satu transaksi. Jika tidak ada siklus, hanya ada tunggu.

### SOURCE B
MySQL/InnoDB: (klaim tidak diverifikasi - docs 403) MySQL juga mendeteksi deadlock via wait-for graph dan memilih victim berdasarkan transaction weight. Bisa saja error 1213 (ER_LOCK_DEADLOCK) untuk deadlock, atau 1205 (ER_LOCK_WAIT_TIMEOUT) untuk timeout.

### ASSESSMENT
Tidak ada bertentangan secara konsisten — keduanya mendeteksi deadlock via mekanisme berbeda:
- PostgreSQL: cek periode deadlock_timeout
- MySQL: wait-for graph traversal aktif

Perbedaan nyata: PostgreSQL menggunakan timeout periodik, MySQL menggunakan pendeteksian aktif. Kedua sistem kemudian memilih victim dan rollback.

**Status**: Tidak ada kontradiksi konstanta — perbedaan arsitektural yang valid.

---

## Contradiction 2

**Apakah deadlock dapat terjadi di PostgreSQL dengan Repeatable Read?**

### SOURCE A
PostgreSQL Docs 13.3.4: Deadlock dapat terjadi bahkan tanpa explicit locking karena row-level locks otomatis pada UPDATE/DELETE.

### SOURCE B
PostgreSQL Docs 13.2.2: Repeatable Read menggunakan MVCC/Timestamp ordering; aplikasi harus retry saat terima 40001 serialization_failure.

### ASSESSMENT
Tidak ada kontradiksi. Deadlock dan serialization failure adalah mekanisme yang berbeda:
- Deadlock (40P01): dua transaksi saling menunggu lock resource yang ditahan
- Serialization failure (40001): konflik snapshot/serialisasi yang tidak memungkinkan urutan serial konsisten

Kedua memerlukan retry transaction.

**Status**:Tidak ada kontradiksi — dua mekanisme concurrency control yang berbeda dalam MVCC PostgreSQL.

---

## Contradiction 3

**Prioritas sumber MySQL tidak tersedia**

MySQL docs resmi (dev.mysql.com) mengembalikan 403 Forbidden. Sumber alternatif MariaDB atau situs pihak ketiga diperlukan untuk konfirmasi definisi error 1213/1205 dan mekanisme deadlock detection InnoDB.

**Status [2026-09-26 Revisi]**: Diverifikasi via Wayback Machine (snapshot 2024, konten identik dev.mysql.com §17.7.5.2). InnoDB: (1) deteksi aktif via wait-for thread saat `innodb_deadlock_detect=ON` (default); (2) victim = transaksi kecil (jumlah row); (3) jika dimatikan, fallback ke `innodb_lock_wait_timeout`. Bukan kontradiksi, melainkan perbedaan arsitektural PG (cek periodik `deadlock_timeout`) vs InnoDB (wait-for graph). Gap akses teratasi.

---

No material contradictions discovered in verified sources.