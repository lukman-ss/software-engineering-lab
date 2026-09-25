# Final Research Synthesis: Backward Compatibility

## Executive Summary
Backward compatibility is an essential capability for mission-critical systems. The core challenge of modern software evolution is not adding new features, but doing so without causing production downtime or breaking existing clients, consumers, and services. The fundamental solution is decoupling changes into independent, reversible phases via the **Expand -> Migrate -> Contract (Parallel Change)** pattern.

---

## Answers to Core Research Questions

### 1. Apa itu backward compatibility?
Backward compatibility adalah karakteristik sistem, API, atau database di mana komponen baru masih dapat menerima input, request, atau struktur data yang dibuat berdasarkan spesifikasi lama tanpa mengalami crash, kegagalan logic, atau merusak integritas data.

### 2. Apa bedanya dengan forward compatibility?
- **Backward compatibility**: Sistem baru memahami dan memproses format dari sistem lama.
- **Forward compatibility**: Sistem lama mampu menerima input dari sistem baru secara aman (misalnya dengan mengabaikan field-field baru yang belum dikenalnya) tanpa mengalami error.

### 3. Apa yang membuat sebuah perubahan menjadi breaking change?
Perubahan menjadi breaking change jika ada ekspektasi contract yang dilanggar:
- Menghapus atau me-rename field/kolom/endpoint.
- Mengubah tipe data (misal: string menjadi list/array).
- Menambahkan parameter wajib (required) baru.
- Memperketat validasi (contoh: batas panjang string diperkecil).
- Mengubah arti semantik dari suatu nilai tanpa mengubah struktur.

### 4. Bagaimana Expand → Migrate → Contract bekerja?
Pola ini membagi satu perubahan menjadi 3 tahap terpisah:
1. **Expand**: Menambahkan struktur baru tanpa menyentuh struktur lama. Keduanya eksis berdampingan.
2. **Migrate**: Memindahkan data lama ke struktur baru, memperbarui aplikasi agar menulis dan membaca ke struktur baru, dan memigrasikan consumer secara bertahap.
3. **Contract**: Menghapus struktur lama setelah dipastikan tidak ada lagi sistem atau client yang menggunakannya.

### 5. Bagaimana melakukan schema migration tanpa downtime?
- Hindari destructive changes (`DROP COLUMN`, `RENAME COLUMN`).
- Buat kolom baru sebagai nullable atau sediakan default value.
- Buat index secara non-blocking (`CONCURRENTLY` di Postgres).
- Pastikan perubahan schema selalu kompatibel dengan versi aplikasi saat ini dan versi yang akan di-deploy (rolling update compatibility).

### 6. Kapan dual read dibutuhkan?
Dual read (fallback read) dibutuhkan ketika proses migrasi data (backfill) masih berlangsung, sehingga aplikasi membaca dari struktur baru; jika data belum ada, aplikasi mengambilnya dari struktur lama sebagai fallback.

### 7. Apa risiko dual write?
- **Peningkatan latency**: Aplikasi harus menunggu dua operasi write selesai.
- **Data drift / Inconsistency**: Jika salah satu write gagal atau terjadi network failure, kedua struktur data menjadi tidak sinkron.
- **Kompleksitas transaksi**: Membutuhkan transaksi database atau outbox pattern agar atomic.

### 8. Bagaimana melakukan backfill secara aman?
- Jalankan secara bertahap dalam ukuran batch kecil (misal 500-1000 record per batch).
- Gunakan throttling/sleep antar batch agar tidak membebani CPU/IO database.
- Buat prosesnya **idempotent** dan **resumable** (menggunakan checkpoint/last processed ID).

### 9. Bagaimana menjaga compatibility saat rolling deployment?
Karena versi N dan N+1 berjalan bersamaan selama proses rolling update, **database schema harus selalu kompatibel dengan N dan N+1 sekaligus**. Jangan pernah menjalankan schema migration yang merusak kebutuhan kode versi N sebelum seluruh instance berhasil di-update ke N+1.

### 10. Bagaimana mengetahui consumer lama masih aktif?
Gunakan instrumentasi observabilitas:
- Metric hit pada legacy endpoints atau legacy field access.
- Structured logging yang mencatat `client_version` atau field access.
- Sunset warning headers (e.g., `Deprecation: true`, `Sunset: <date>`) pada response HTTP.

### 11. Kapan field, column, atau endpoint lama boleh dihapus?
Field, kolom, atau endpoint lama hanya boleh dihapus saat metric penggunaan legacy menunjukkan angka **nol secara konsisten** selama periode waktu tertentu (misal: 30 hari berturut-turut).

### 12. Bagaimana observability membantu migration?
Observability menyediakan data real-time mengenai:
- Jumlah read/write ke schema lama vs schema baru (`legacy_read_count`, `new_read_count`).
- Error rate selama migrasi (`migration_error_count`).
- Data mismatch count untuk mendeteksi data drift antara old dan new tables.

### 13. Bagaimana feature flag membantu rollout dan rollback?
Feature flag memisahkan deployment kode dari aktivasi fungsionalitas. Flag memungkinkan perubahan jalur baca/tulis diaktifkan secara bertahap (canary: 1% -> 10% -> 100%) dan dapat di-rollback secara instan (flip toggle ke `false`) tanpa perlu me-redeploy aplikasi jika ditemukan bug.

### 14. Bagaimana rollback memengaruhi desain migration?
Desain migration harus selalu **reversible**:
- Jangan terburu-buru menghapus kolom lama.
- Pertahankan dual write selama fase transisi, sehingga jika aplikasi di-rollback ke versi lama, data terbaru masih tersedia di schema lama dan tidak terjadi data loss.

### 15. Apa failure mode paling umum dalam backward-compatible migration?
- Menghapus kolom lama sebelum seluruh instance dan consumer bermigrasi.
- Data drift akibat dual write parsial yang tidak atomic.
- Backfill query besar yang me-lock tabel produksi dan menyebabkan cascading outage.
- Menelantarkan fase Contract sehingga menciptakan timbunan technical debt permanen.

### 16. Bagaimana membuat migration resumable dan idempotent?
- **Idempotent**: Eksekusi berulang terhadap data yang sama menghasilkan state akhir yang identik (misalnya menggunakan `UPSERT` atau query `WHERE new_field IS NULL`).
- **Resumable**: Menyimpan checkpoint (misalnya `last_processed_id`) ke disk atau tabel metadata sehingga jika backfill terhenti, ia dapat melanjutkan dari record terakhir.

### 17. Apa yang harus diuji sebelum contract lama dihapus?
- Verifikasi observabilitas bahwa traffic ke legacy interface adalah nol.
- Menjalankan contract testing (misalnya Pact) dan end-to-end integration tests untuk memastikan tidak ada consumer internal yang masih bergantung pada format lama.
- Uji rollback plan sebelum mengeksekusi penghapusan permanen.

### 18. Apa perbedaan database compatibility dan API compatibility?
- **Database Compatibility**: Mengatur kompatibilitas antara kode aplikasi dan schema storage. Dipengaruhi oleh lifecycle deployment aplikasi dan locking engine database.
- **API Compatibility**: Mengatur interaksi antara server dan client luar (web, mobile app, partner integration). Mobile app seringkali tidak dapat dipaksa untuk update seketika, sehingga backward compatibility API seringkali harus dipertahankan dalam jangka waktu yang jauh lebih panjang (berbulan-bulan hingga bertahun-tahun).
