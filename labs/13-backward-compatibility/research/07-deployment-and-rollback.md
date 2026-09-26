# Deployment Sequence dan Rollback Strategies

## 1. Rolling Deployment Realities

Selama rolling deployment, selalu ada jendela di mana:

- Instance A menjalankan **Version N** (Kode lama)
- Instance B menjalankan **Version N+1** (Kode baru)

Kedua instance berbagi database yang sama secara bersamaan.

### 1.1 Fundamental Rule of Migration Safety

**Evidence** (Martin Fowler, Blue Green Deployment):
"The trick is to separate the deployment of schema changes from application upgrades. So first apply a database refactoring to change the schema to support both the new and old version of the application, deploy that, check everything is working fine so you have a rollback point, then deploy the new version of the application. (And when the upgrade has bedded down remove the database support for the old version.)"

> **Database schema harus SELALU kompatibel dengan Version N dan Version N+1 secara bersamaan.**

Jika Version N+1 mengubah kolom dengan cara yang memutus Version N, Instance A akan error SQL saat rolling update.

## 2. Urutan Deployment Aman

Untuk menjamin kompatibilitas lintas versi, migrasi harus dipecah ke deployment terpisah:

### 2.1 Deployment 1: Expand Schema

**Aksi**: Terapkan migrasi database untuk membuat tabel/kolom baru (semua non-breaking, nullable, atau dengan default).

**Validasi**: Instance Version N tetap berjalan normal karena kolom lama tidak tersentuh.

### 2.2 Deployment 2: Dual Write & Support New Schema

**Aksi**: Deploy versi aplikasi yang mendukung dual write dan worker backfill. Instance lama terus membaca field lama tanpa crash.

### 2.3 Run Backfill

**Aksi**: Jalankan script backfill data di background untuk sinkronkan data legacy ke skema baru.

### 2.4 Deployment 3: Switch Read Path

**Aksi**: Deploy versi aplikasi yang membaca eksklusif dari skema baru (atau menggunakan feature flag untuk ramp up baca).

### 2.5 Deployment 4: Stop Dual Write

**Aksi**: Deploy versi aplikasi yang berhenti menulis ke skema legacy.

### 2.6 Deployment 5: Contract Schema

**Aksi**: Drop kolom atau tabel legacy di database.

## 3. Skenario Rollback

Apa yang terjadi jika Application N+1 di-deploy, menemui masalah, dan harus di-rollback ke Application N?

### 3.1 Jika Schema Diperluas (Non-Destructive) — **AMAN**

**Evidence** (Fowler, Blue Green):
"first apply a database refactoring to change the schema to support both the new and old version of the application, deploy that, check everything is working fine so you have a rollback point, then deploy the new version of the application."

Application N masih berfungsi karena kolom/tabel legacy tidak tersentuh.

### 3.2 Jika Rollback Terjadi Selama Dual Write — **DATA TETAP AMAN**

Jika kedua skema dipertahankan secara sinkron (dalam transaksi yang sama), data tidak hilang.

### 3.3 Jika Rollback Terjadi **SETELAH** Menghentikan Dual Write — **RISIKO DATA LOSS**

**Klasifikasi**: **Critical Risk**
- Application N mengharapkan data di skema lama
- Skema lama tidak diperbarui selama runtime N+1
- Data yang ditulis saat N+1 berjalan tidak ada di skema lama

**Kebijakan**: Jangan pernah menjalankan fase Contract sampai versi aplikasi baru menunjukkan operasi stabil di produksi selama periode observasi berkelanjutan.

### 3.4 Rollback vs Roll-Forward

**Evidence** (Google AIP-180):
Kompatibilitas forward memengaruhi kemampuan rollback. Jika migrasi forward-compatible (consumer lama bisa baca data baru), maka rollback aman. Jika tidak, rollback memerlukan data migration down (tidak disarankan).

## 4. Blue-Green Deployment dan Schema

**Evidence** (Fowler, Blue Green Deployment):
"Databases can often be a challenge with this technique, particularly when you need to change the schema to support a new version of the software. The trick is to separate the deployment of schema changes from application upgrades."

**Implementasi**:
1. Apply schema migration (expand) ke **dua** environment blue dan green
2. Deploy aplikasi baru ke green (dengan dual write)
3. Verify dan cutover
4. Setelah stabil, apply contract migration

---

## Catatan Kualitas Bukti

| Pernyataan | Sumber | Confidence | Keterangan |
|------------|--------|------------|------------|
| Schema harus support N dan N+1 | Fowler Blue Green | HIGH | Explisit di dokumentasi |
| 6-deployment sequence | Fowler + prinsip umum | HIGH | Validasi dari Fowler schema first |
| Rollback aman jika hanya expand | Fowler Blue Green | HIGH | Terbukti pattern |
| Rollback berisiko setelah stop dual write | Logika deduktif | HIGH | Data loss karena gap writes |