# Panduan Backward Compatibility: Pola Expand -> Migrate -> Contract pada Schema dan API

## Problem
Mengubah schema database atau contract API di sistem produksi sering kali memicu insiden downtime dan data loss jika dilakukan secara langsung (in-place modification). Ketika kolom diubah namanya atau dihapus, consumer versi lama yang masih berjalan langsung gagal melakukan deserialisasi atau query database.

## Why This Matters
Dalam arsitektur modern berbasis microservices dan continuous delivery:
1. Rolling deployment menyebabkan aplikasi versi N dan versi N+1 berjalan bersamaan.
2. Client eksternal (seperti mobile app atau third-party integration) tidak dapat dipaksa untuk meng-upgrade versi secara serentak.
3. Mengabaikan backward compatibility berisiko menimbulkan cascading failures, downtime layanan, dan inkonsistensi data.

## Mental Model
Pola dasar untuk evolusi sistem yang aman adalah **Expand -> Migrate -> Contract** (dikenal juga sebagai *Parallel Change*):
- **Expand**: Perluas schema atau contract baru tanpa mengubah atau menghapus yang lama. Keduanya berjalan berdampingan.
- **Migrate**: Sinkronkan data lama ke baru, alihkan pembacaan data, dan arahkan consumer ke interface baru secara bertahap.
- **Contract**: Hapus schema atau contract lama hanya setelah seluruh consumer bermigrasi dan traffic ke format lama bernilai nol.

## Core Concept
- **Backward Compatibility**: Kemampuan sistem versi baru untuk tetap memahami dan memproses format atau data dari sistem lama.
- **Forward Compatibility**: Kemampuan sistem versi lama untuk menerima input dari sistem baru secara aman tanpa crash (misalnya dengan mengabaikan field baru).
- **Breaking Change**: Perubahan yang melanggar kontrak yang disepakati, seperti menghapus field, mengubah tipe data, atau menambahkan validasi wajib baru.

## Failure Scenario
Beberapa skenario kegagalan umum dalam evolusi schema:
1. **Destructive Alteration**: Mengubah nama kolom atau menghapus kolom secara langsung menyebabkan consumer V1 mengalami crash.
2. **Premature Read Switch**: Mengalihkan pembacaan ke struktur baru sebelum proses backfill selesai mengakibatkan pembacaan data kosong (data starvation).
3. **Dual-Write Drift**: Menulis ke dua tempat secara tidak atomic memicu inkonsistensi saat salah satu operasi write gagal.
4. **Premature Contract**: Menghapus interface lama saat masih ada traffic aktif memicu kegagalan runtime bagi consumer yang belum sempat update.

## How It Works
Evolusi sistem dilakukan dalam 5 tahap:
1. **Tahap 1 (Expand)**: Tambahkan tabel atau field baru. Struktur lama tetap dipertahankan.
2. **Tahap 2 (Dual-Write)**: Aplikasi menulis ke struktur lama dan struktur baru secara bersamaan untuk menjaga konsistensi data baru.
3. **Tahap 3 (Backfill)**: Worker berjalan di background memindahkan data historis dari struktur lama ke struktur baru secara bertahap (batching) dan idempoten.
4. **Tahap 4 (Read Switch & Fallback)**: Aplikasi dialihkan untuk membaca dari struktur baru. Jika data belum ter-backfill, sistem menggunakan fallback read ke struktur lama.
5. **Tahap 5 (Contract)**: Setelah metrik traffic lama mencapai nol, matikan dual-write dan hapus struktur lama secara aman.

## Architecture
```text
[ Client V1 (Legacy) ]     [ Client V2 (Modern) ]
          │                         │
          ▼                         ▼
┌──────────────────────────────────────────────────┐
│              HTTP / Service Layer                │
│  - Transformation / Deprecation Pipeline         │
│  - Feature Flags (WriteMode, ReadMode)           │
│  - Observability Metrics (Legacy/New Traffic)    │
└─────────┬────────────────────────────────┬───────┘
          │                                │
          ▼                                ▼
┌──────────────────┐             ┌─────────────────┐
│ Legacy Storage   │◄──Backfill──┤ Modern Storage  │
│ (users.phone)    │   Worker    │ (user_phones)   │
└──────────────────┘             └─────────────────┘
```

## Implementation
Lab ini mengimplementasikan pola evolusi dari relasi 1:1 (`users.phone`) menuju 1:N (`user_phones` table) menggunakan Go:
- `MemoryStore`: Mengelola persistensi data dan mensimulasikan penghapusan kolom legasi.
- `FeatureFlagManager`: Mengatur mode baca (`LegacyOnly`, `FallbackRead`, `NewOnly`) dan tulis (`LegacyOnly`, `DualWrite`, `NewOnly`).
- `BackfillWorker`: Memindahkan data secara chunked/batched dengan pencatatan checkpoint.
- `MetricsCollector`: Memonitor hit legacy vs modern untuk memvalidasi kesiapan fase contract.

## What the Tests Prove
Pengujian otomatis (`tests/migration_test.go`, `tests/concurrency_test.go`, `internal/compat/service_test.go`) memvalidasi:
1. **Consumer Independence**: Consumer lama dan baru dapat membaca payload yang sesuai tanpa error selama transisi.
2. **Atomic Dual-Write**: Mutasi baru tersimpan di kedua representasi tanpa kehilangan data.
3. **Backfill Idempotency**: Menjalankan backfill berulang kali tidak menduplikasi data.
4. **Fallback Safety**: Pembacaan data yang belum di-backfill tidak menghasilkan data kosong.
5. **Rollback Reversibility**: Revert flag kembali ke V1 tetap aman dan tidak merusak data.
6. **Concurrency Safety**: Eksekusi paralel pembacaan, penulisan, dan backfill lolos uji race condition (`go test -race`).

## Production Considerations
- **Dual-Write Latency & Atomicity**: Menulis ke dua tempat menambah overhead latency. Di database SQL, gunakan transaksi database atau Transactional Outbox pattern.
- **Backfill Throttling**: Lakukan pembacaan data historis dalam batch kecil (misal 500-1000 record) dengan jeda waktu untuk menghindari database locking dan lonjakan CPU/IO.
- **Rolling Deployment**: Pastikan database schema selalu kompatibel dengan versi kode N dan N+1 sekaligus.
- **Monitoring Window**: Window monitoring (misalnya 30 hari tanpa traffic legacy) adalah panduan operasional (operational guideline) yang fleksibel tergantung SLA organisasi, bukan standar mutlak.

## Checklist
- [ ] Schema baru dibuat tanpa menyentuh atau mengubah schema lama (Expand).
- [ ] Dual-write diaktifkan via feature flag.
- [ ] Backfill berjalan secara batching, idempoten, dan dapat di-resume jika gagal.
- [ ] Fallback read aktif saat read switch diuji.
- [ ] Metrik legacy traffic dipantau hingga bernilai nol secara konsisten.
- [ ] Schema lama dihapus setelah diverifikasi aman (Contract).

## Key Takeaways
1. Jangan pernah melakukan perubahan destruktif langsung pada schema atau API yang aktif digunakan.
2. Pola Expand-Migrate-Contract memisahkan migrasi menjadi tahapan yang independen dan dapat dibatalkan (reversible).
3. Backfill harus dibuat idempoten dan berbasis checkpoint.
4. Observabilitas adalah prasyarat mutlak sebelum mengeksekusi fase Contract.
