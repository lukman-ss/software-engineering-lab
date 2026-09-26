# Failure Modes dan Mitigasi dalam Backward-Compatible Migration

## 1. Failure Mode: Immediate Destructive Schema Alteration

### Incident
Engineer menjalankan `ALTER TABLE customers DROP COLUMN phone;` atau mengubah tipe kolom secara langsung.

### Impact
Semua instance aplikasi yang berjalan langsung crash dengan error SQL syntax/column missing.

### Evidence
MongoDB menyebutkan: "Dropping a collection or a field is a breaking change" — hal yang serupa untuk tabel relasional.

### Mitigation
1. **CI/CD Checks**: Konfigurasi linting migrasi untuk membatalkan `DROP COLUMN`, `RENAME COLUMN`, atau penambahan kolom non-nullable tanpa default.
2. **Migration Review Process**: Setiap migrasi harus melalui review secara manual sebelum production.

## 2. Failure Mode: Dual-Write Inconsistency (Split-Brain)

### Incident
Aplikasi menulis ke Table A (lama) dan Table B (baru). Jika write kedua gagal atau network timeout, data di kedua tabel tidak sinkron.

### Impact
- Data drift antara schema lama dan baru
- Kesalahan laporan
- Inconsistensi read dari consumer yang bergantung pada schema mana

### Evidence
Microservices.io (Chris Richardson), Transactional Outbox: "The command must atomically update the database and send messages in order to avoid data inconsistencies and bugs. However, it is not viable to use a traditional distributed transaction (2PC) that spans the database and the message broker."
Microservices.io (Chris Richardson): "If the database transaction commits then the messages must be sent. Conversely, if the database rolls back, the messages must not be sent."

### Mitigation
1. **Transaksi Database**: Enclose dual writes dalam transaksi yang sama.
2. **Data Reconciliation**: Jalankan skrip audit secara periodik untuk mendeteksi perbedaan antara tabel lama dan baru.
3. **Outbox Pattern**: Gunakan message queue untuk memastikan konsumsi yang konsisten.

**Evidence** (Prisma/arkatechno/zero-downtime-postgresql):
"Using the outbox pattern avoids both the dual write problem and the risk of data loss in a failure scenario."

## 3. Failure Mode: Missing Backfill Records

### Incident
Read dialihkan ke tabel baru sebelum backfill selesai, mengembalikan 404/Empty untuk baris legacy.

### Impact
- Read error pada API
- Data tidak ditemukan untuk koneksi lama
- Pelanggan mengalami downtime

### Evidence
Martin Fowler, Parallel Change: "fallback read mencegah data starvation" secara implisit dijelaskan dalam fase migrate.

### Mitigation
1. **Fallback Reading**: Jika record tidak ditemukan di tabel baru, ambil dari tabel lama secara async dan lazy-write ke tabel baru.
2. **Checkpoint Validation**: Pastikan jumlah baris historis di old table sama dengan new table sebelum menghentikan fallback read.

## 4. Failure Mode: Permanent Expand State (Contract Phase Abandoned)

### Incident
Engineering team memigrasi read/write, tetapi lupa untuk membersihkan tabel/kolom legacy dan logika dual-write.

### Impact
- Technical debt yang menumpuk
- Penggunaan disk dan CPU yang tidak perlu
- Kehbingungan code apakah field masih diperlukan

### Evidence
Martin Fowler, Feature Toggle: "Release flags are a useful technique and lots of teams use them. However they should be your last choice when you're dealing with putting features into production." — menekankan pentingnya pembersihan.

### Mitigation
1. **Explicit Deprecation Deadlines**: Set deadline eksplicit untuk penghapusan kode.
2. **Legacy Usage Metrics**: Track metrik akses kode/field legacy; trigger alert otomatis ketika penggunaan legacy = 0 selama jendela observasi (30 hari sebagai heuristic).
3. **Code Ownership**: Setiap expand phase harus memiliki ticket follow-up untuk contract phase.

**Evidence** (NOT VERIFIED): Heuristic 30 hari tidak didukung oleh data Tier 1 untuk semua kasus; harus diterapkan berdasarkan metrik aktual.

## 5. Failure Mode: Data Corruption pada Column Type Change

### Incident
Mengubah tipe kolom `phone VARCHAR(20)` ke `phone VARCHAR(50)` atau sebaliknya dengan data yang tidak dapat dikonversi.

### Impact
- Query return error
- Aplikasi crash saat mencoba membaca nilai yang tidak kompatibel
- Data loss jika konversi gagal

### Evidence
PostgreSQL ALTER TABLE: "ALTER TABLE will attempt to convert the column's default value (if any) to the new type, as well as any constraints that involve the column. But these conversions might fail, or might produce surprising results."

### Mitigation
1. Gunakan `USING` clause untuk konversi eksplisit
2. Backup data sebelum konversi
3. Test konversi pada environment non-production dulu

## 6. Failure Mode: Transaction Deadlock pada Dual Write

### Incident
Aplikasi menulis ke dua tabel dalam urutan yang sama untuk semua requests, menyebabkan deadlock atau high contention.

### Impact
- Waktu respons meningkat
- Timeout pada request
- Deadlocks yang menghancurkan transaksi

### Mitigation
1. Pastikan urutan write konsisten (misal: selalu write ke table lama dulu, baru table baru)
2. Gunakan retry logic dengan exponential backoff
3. Pertimbangkan asynchronous dual write via queue

---

## Ringkasan Risk Priority

| Failure Mode | Risk Level | Detection | Mitigation Primary |
|--------------|------------|-----------|-------------------|
| Destructive Schema | HIGH | Deployment | CI/CD lint + review |
| Dual-Write Inconsistency | HIGH | Data audit | Transaksi + reconciliation |
| Missing Backfill | MEDIUM | Metrics | Fallback read |
| Contract Abandoned | MEDIUM | Metrics | Deprecation tracking |
| Type Change Corruption | HIGH | Testing | USING clause + backup |
| Transaction Deadlock | MEDIUM | Monitoring | Consistent ordering |