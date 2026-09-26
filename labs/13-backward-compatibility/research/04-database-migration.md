# Database Schema Evolution dan Zero-Downtime Migration

## 1. Prinsip Migrasi Tanpa Downtime

Migrasi skema database harus dapat dijalankan selama database aktif menangani transaksi secara bersamaan. Prinsip utama:

### 1.1 Additive Schema Change (Penambahan Non-Destructive)

**Evidence**:
- PostgreSQL ALTER TABLE dokumentasi: "Adding a column with a constant default value does not require each row of the table to be updated when the ALTER TABLE statement is executed. Instead, the default value will be returned the next time the row is accessed, and applied when the table is rewritten."
- PostgreSQL: DROP COLUMN tidak menghapus data secara fisik, hanya membuat kolom tidak terlihat — data tetap ada hingga UPDATE row atau VACUUM FULL.

**Praktik**:
```sql
-- Expand: Tambah kolom nullable dengan default di luar tabel
ALTER TABLE customers ADD COLUMN phone_new VARCHAR(20) DEFAULT NULL;
```

### 1.2 Index Creation Concurrently

**Evidence**:
- PostgreSQL: "CREATE UNIQUE INDEX CONCURRENTLY" memungkinkan pembuatan index tanpa lock tabel untuk UPDATE/INSERT.

**Praktik**:
```sql
-- Expand: Buat index untuk tabel baru secara parallel
CREATE INDEX CONCURRENTLY idx_customer_phone_customer_id ON customer_phones(customer_id);
```

## 2. Table Locking dan Transaksi Panjang

### 2.1 LOCKING LEVEL pada ALTER TABLE

**Evidence** (PostgreSQL):
- `ADD COLUMN` dengan constant default: tidak perlu table rewrite (ACCESS EXCLUSIVE lock tidak diperlukan secara signifikan)
- `ADD CONSTRAINT` umumnya memerlukan scan tabel; `NOT VALID` memungkinkan skip scan
- `SET DATA TYPE` biasanya merewrite seluruh tabel

### 2.2 Avoiding Lock pada Large Tables

**Evidence** (PostgreSQL Notes):
"Table and/or index rebuilds may take a significant amount of time for a large table, and will temporarily require as much as double the disk space."

**Praktik**:
1. Jika memungkinkan, gunakan `USING` clause untuk konversi tipe data dengan logika idempotent.
2. Jalankan migrasi di maintenance window JIKA required rewrite tak terhindar.
3. Gunakan `NOT VALID` untuk constraint, validasi nanti secara async.

## 3. Data Backfill Strategies

### 3.1 Batch Processing

**Evidence**:
- PostgreSQL: Menambah constraint wajib (CHECK, NOT NULL, FOREIGN KEY) memerlukan scan seluruh tabel untuk memverifikasi semua rows.
- Konvensi: batch size 1000-5000 adalah praktik umum untuk hindari memory exhaustion dan lock timeout.

**Praktik**:
```sql
-- Contoh backfill batch
UPDATE customer_phones 
SET phone_number = (SELECT phone FROM customers c WHERE c.id = customer_phones.customer_id)
WHERE customer_phones.phone_number IS NULL
AND customer_phones.customer_id IN (
    SELECT id FROM customers 
    WHERE phone IS NOT NULL 
    LIMIT 1000
);
```

### 3.2 Throttling dan Sleep

**Evidence**:
- Umum: batch tanpa throttling dapat menghasilkan spike I/O yang menghancurkan performa produksi.

**Praktik**:
```sql
-- Gunakan pseudocode atau external scheduler
PERFORM pg_sleep(0.1); -- 100ms delay antar batch
```

### 3.3 Idempotency dan Resumable Migration

**Evidence** (PostgreSQL):
- `NOT VALID` constraint memungkinkan penambahan constraint tanpa memverifikasi data lama, sehingga idempotent.
- UPSERT (`INSERT ... ON CONFLICT ... DO UPDATE`) adalah pola idempotent standar.

**Praktik**:
```sql
-- Idempotent: hanya insert jika belum ada
INSERT INTO customer_phones (customer_id, phone_number)
SELECT id, phone FROM customers
WHERE phone IS NOT NULL
ON CONFLICT (customer_id, phone_number) DO NOTHING;
```

## 4. Dual Write dan Dual Read

### 4.1 Dual Write Pattern

**Evidence**:
- Martin Fowler Parallel Change: "During the *migrate* phase you update all clients using the old version to the new version. This can be done incrementally and, in the case of external clients, this will be the longest phase."

**Risiko** (dokumentasi tidak ada sumber utama):
- Data drift: Jika write ke dua tabel tidak atomic
- Latency peningkatan: Aplikasi menunggu dua operasi selesai

### 4.2 Fallback Read (Dual Read)

**Evidence**:
- Prinsip umum: hingga semua data terbackfill, aplikasi harus dapat membaca dari kedua sumber.

**Implementasi**:
```
read_customer(id):
  result = query new_table where customer_id = id
  if result.empty():
      result = query legacy_table where id = id
      if result.found():
          async_write_to_new_table(result)  // lazy migration
  return result
```

## 5. Kontrak pada Migrasi Tabel Terhubung

### 5.1 Foreign Key Constraint Timing

**Evidence** (PostgreSQL):
"ADD FOREIGN KEY requires only a SHARE ROW EXCLUSIVE lock on the referenced table, and a foreign lock on the table on which the constraint is declared."

### 5.2 Adding nullable column first, then making required

**Evidence** (Stripe API versioning):
Mengubah field dari tidak required ke required memerlukan:
1. Expand: Tambah field baru sebagai nullable
2. Populate: Backfill semua data
3. Contract: Tambahkan NOT NULL constraint

---

## Catatan Kualitas Bukti

| Pernyataan | Sumber | Confidence | Keterangan |
|------------|--------|------------|------------|
| ALTER TABLE constant default tidak require rewrite | PostgreSQL Docs | HIGH | Ditunjukkan di dokumentasi resmi |
| CREATE INDEX CONCURRENTLY | PostgreSQL Docs | HIGH | Fitur resmi PostgreSQL |
| NOT VALID untuk constraint | PostgreSQL Docs | HIGH | Dokumentasi resmi |
| Batch processing untuk backfill | Umum prinsip | MEDIUM | Inferensi dari lock behavior PostgreSQL |
| Dual write/fallback read pattern | Martin Fowler | HIGH | Dijelaskan secara eksplisit dalam Parallel Change |