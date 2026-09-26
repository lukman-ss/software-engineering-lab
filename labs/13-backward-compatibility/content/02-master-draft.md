# Pandangan Teknis: Evolusi Skema Database yang Kompatibel dengan Mundar-Melit

## Masalah
Kemampuan untuk mengubah skema produksi tanpa gangguan adalah kebutuhan kritis bagi sistem yang tidak pernah berhenti. Perubahan skema destruktif langsung (`DROP COLUMN`, `RENAME COLUMN`) atau penambahan kolom wajib dapat menyebabkan kesalahan deserialisasi pada klien lama dan downtime layanan selama pengulangan deployment. Tantangan inti adalah: bagaimana mengalihkan dari skema 1:1 (satu nilai per entitas) ke skema 1:N (nilai berganda per entitas) sambil mempertahankan kompatibilitas mundar-melit untuk klien yang belum diupgrade?

## Mengapa Ini Penting
Sistem kritis yang berjalan 24/7 tidak mampu membeli biaya downtime untuk migrasi skema. Tanpa strategi kompatibel mundar-melit, setiap perubahan skema berisiko:
- Gangguan layanan saat klien lama gagal memproses respons baru
- Kerusakan data akibat penulisan tidak sinkron selama transisi
- Keputusan rollback yang kompleks jika masalah ditemukan setelah deploy
- Akumulasi utang teknis saat fase kontrak selalu ditangguhkan

## Model Mental
Kompatibilitas mundar-melit dicapai dengan memecah perubahan skema struktural menjadi tiga fase operasional yang terpisah dan dapat dibalik: **Expand → Migrate → Contract** (Parallel Change). Ini memastikan bahwa pada setiap saat, sistem mendukung baik representasi lama maupun baru, memungkinkan pengulangan deployment dan rollback yang aman.

## Konsep Inti

### Apa Itu Kompatibilitas Mundar-Melit?
Kompatibilitas mundar-melit adalah karakteristik sistem, API, atau basis data di mana komponen baru masih dapat menerima input, request, atau struktur data yang dibuat berdasarkan spesifikasi lama tanpa mengalami crash, kegagalan logika, atau merusak integritas data.

### Mengapa Sebuah Perubahan Menjadi Perubahan yang Mengganggu (Breaking Change)?
Perubahan menjadi perubahan yang mengganggu jika ada ekspektasi kontrakt yang dilanggar:
- Menghapus atau mengubah nama field/kolom/endpoint.
- Mengubah tipe data (misal: string menjadi list/array).
- Menambahkan parameter wajib (required) baru.
- Memperketat validasi (contoh: batas panjang string diperkecil).
- Mengubah arti semantik dari suatu nilai tanpa mengubah struktur.

### Bagaimana Pola Expand → Migrate → Contract Bekerja?
Pola ini membagi satu perubahan skema menjadi tiga fase terpisah:
1. **Expand**: Menambahkan struktur baru tanpa menyentuh struktur lama. Keduanya eksis berdampingan.
2. **Migrate**: Memindahkan data lama ke struktur baru, memperbarui aplikasi agar menulis dan membaca ke struktur baru, dan memigrasikan consumer secara bertahap.
3. **Contract**: Menghapus struktur lama setelah dipastikan tidak ada lagi sistem atau klien yang menggunakannya.

### Bagaimana Melakukan Migrasi Skema Tanpa Downtime?
- Hindari perubahan destruktif (`DROP COLUMN`, `RENAME COLUMN`).
- Buat kolom baru sebagai nullable atau sediakan nilai default.
- Buat index secara non-blocking (`CONCURRENTLY` di Postgres).
- Pastikan perubahan skema selalu kompatibel dengan versi aplikasi saat ini dan versi yang akan di-deploy (kompatibilitas pengulangan deployment).

### Kapan Dual Read Diperlukan?
Dual read (pembacaan cadangan) diperlukan ketika proses migrasi data (backfill) masih berlangsung, sehingga aplikasi membaca dari struktur baru; jika data belum ada, aplikasi mengambilnya dari struktur lama sebagai fallback.

### Apa Risiko Dual Write?
- **Peningkatan latensi**: Aplikasi harus menunggu dua operasi write selesai.
- **Data drift / Ketidaksesuaian**: Jika salah satu write gagal atau terjadi kegagalan jaringan, kedua struktur data menjadi tidak sinkron.
- **Kompleksitas transaksi**: Membutuhkan transaksi basis data atau outbox pattern untuk operasi atomik.

### Bagaimana Melakukan Backfill Secara Aman?
- Jalankan secara bertahap dalam ukuran batch kecil (misal 500-1000 record per batch).
- Gunakan throttling/sleep antar batch agar tidak membebani CPU/IO basis data.
- Buat prosesnya idempotent dan resumable (menggunakan checkpoint/last processed ID).

### Bagaimana Menjaga Kompatibilitas Selama Pengulangan Deployment?
Karena versi N dan N+1 berjalan bersamaan selama proses pengulangan deployment, **basis data skema harus selalu kompatibel dengan N dan N+1 sekaligus**. Jangan pernah menjalankan migrasi skema yang merusak kebutuhan kode versi N sebelum seluruh instance berhasil di-update ke N+1.

### Bagaimana Mengetahui Klien Lama Masih Aktif?
Gunakan instrumen osservabilitas:
- Metric hit pada endpoint lama atau akses field lama.
- Pencatatan terstruktur yang mencatat `client_version` atau akses field.
- Header peringatan sunset (misalnya, `Deprecation: true`, `Sunset: <date>`) pada respons HTTP.

### Kapan Field, Kolom, atau Endpoint Lama Boleh Dihapus?
Field, kolom, atau endpoint lama hanya boleh dihapus saat metrik penggunaan legacy menunjukkan angka **nol secara konsisten** selama periode waktu tertentu.
> [!NOTE] Indeks industri — `NOT VERIFIED`: Periode 30 hari belum didukung data Tingkat 1 untuk semua kasus; terapkan berdasarkan metrik penggunaan legacy = 0 secara nyata.

### Bagaimana Observabilitas Membantu Migrasi?
Observabilitas menyediakan data real-time mengenai:
- Jumlah baca/tulis ke skema lama vs skema baru (`legacy_read_count`, `new_read_count`).
- Tingkat kesalahan selama migrasi (`migration_error_count`).
- Hitung ketidaksesuaian data untuk mendeteksi drift data antara tabel lama dan baru.

### Bagaimana Fitur Flag Membantu Rollout dan Rollback?
Fitur flag memisahkan deployment kode dari aktivasi fungsionalitas. Flag memungkinkan perubahan jalur baca/tulis diaktifkan secara bertahap (canary: 1% -> 10% -> 100%) dan dapat di-rollback secara instan (alihkan toggle ke `false`) tanpa perlu mendeploy ulang aplikasi jika ditemukan bug.

### Bagaimana Rollback Memengaruhi Desain Migrasi?
Desain migrasi harus selalu **reversibel**:
- Jangan terburu-buru menghapus kolom lama.
- Pertahankan dual write selama fase transisi, sehingga jika aplikasi di-rollback ke versi lama, data terbaru masih tersedia di skema lama dan tidak terjadi kehilangan data.

### Apa Failure Mode Paling Umum dalam Migrasi yang Kompatibel Mundar-Melit?
- Menghapus kolom lama sebelum seluruh instance dan klien bermigrasi.
- Drift data akibat dual write parsial yang tidak atomik.
- Query backfill besar yang me-lock tabel produksi dan menyebabkan outage berantai.
- Menelantarkan fase Contract sehingga menciptakan timbunan utang teknis permanen.

### Bagaimana Membuat Migrasi yang Resumable dan Idempotent?
- **Idempotent**: Eksekusi berulang terhadap data yang sama menghasilkan state akhir yang identik (misalnya menggunakan `UPSERT` atau query `WHERE new_field IS NULL`).
- **Resumable**: Menyimpan checkpoint (misalnya `last_processed_id`) ke disk atau tabel metadata sehingga jika backfill terhenti, ia dapat melanjutkan dari record terakhir.

### Apa yang Harus Diuji Sebelum Kontrak Lama Dihapus?
- Verifikasi osservabilitas bahwa traffic ke interface legacy adalah nol.
- Jalankan contract testing (misalnya Pact) dan uji integrasi end-to-end untuk memastikan tidak ada klien internal yang masih bergantung pada format lama.
- Uji rencana rollback sebelum mengeksekusi penghapusan permanen.

### Apa Perbedaan Kompatibilitas Basis Data dan Kompatibilitas API?
- **Kompatibilitas Basis Data**: Mengatur kompatibilitas antara kode aplikasi dan basis data penyimpanan. Dipengaruhi oleh siklus hidup deployment aplikasi dan mesin kunci basis data.
- **Kompatibilitas API**: Mengatur interaksi antara server dan klien luar (seluler, aplikasi web, integrasi mitra). Aplikasi seluler seringkali tidak dapat dipaksa untuk update secara langsung, sehingga kompatibilitas API mundar-melit seringkali harus dipertahankan dalam jangka waktu yang jauh lebih lama (berbulan-bulan hingga bertahun-tahun).

## Arsitektur
```text
[ Klien V1 (Legacy) ]     [ Klien V2 (Modern) ]
          │                         │
          ▼                         ▼
┌──────────────────────────────────────────────────┐
│              Lapisan HTTP / Layanan              │
│  - Transformasi / Pipa Pengusiran                │
│  - Fitur Flag (WriteMode, ReadMode)              │
│  - Metrik Observabilitas (Traffic Legacy/Baru)   │
└─────────┬────────────────────────────────┬───────┘
          │                                │
          ▼                                ▼
┌──────────────────┐             ┌─────────────────┐
│ Penyimpanan Lama │◄──Backfill──┤ Penyimpanan Baru│
│ (users.phone)    │   Worker    │ (user_phones)   │
└──────────────────┘             └─────────────────┘
```

## Implementasi
Implementasi laboratorium ini menggunakan penyimpanan dalam memori thread-safe untuk mensimulasikan perilaku basis data relasional tanpa dependensi eksternal, dengan fokus pada pola paralel perubahan.

### Komponen Utama
1. **Storage**: Map berdasarkan thread-safe yang mensimulasikan tabel `users` (kolom legacy `phone`) dan `user_phones` (tabel anak dengan `user_id`, `number`, `is_primary`). Mendukung dual-write atomik dan penurunan kontrak kolom.
2. **Pengelola Fitur Flag**: Mengontrol tahap migrasi:
   - `WriteMode`: `LegacyOnly`, `DualWrite`, `NewOnly`.
   - `ReadMode`: `LegacyOnly`, `FallbackRead`, `NewOnly`.
   - `ContractApplied`: Flag boolean yang menunjukkan skema lama telah dihapus.
3. **Kolektor Metrik**: Counter thread-safe untuk `LegacyReadHits`, `NewReadHits`, `DualWriteCount`, `DualWriteErrors`, `DriftDetectedCount`.
4. **Pekerja Backfill**: Pekerja migrasi data yang resumable, berorientasi batch menggunakan checkpoint ID dan logika upsert idempoten.
5. **Layanan Kompatibilitas**: Fasade domain yang mengatur operasi baca/tulis, audit reconciliasi data, dan siklus hidup pengusiran.
6. **Penangan HTTP**: Penangan standar yang mengembalikan kontrak JSON, header pengusiran (`Deprecation: @epoch`, `Sunset: @epoch`), dan kode status.

### Alur Kontrol Lintasan Kode
- **CreateUser**: Menulis sesuai dengan WriteMode saat ini (LegacyOnly, DualWrite, atau NewOnly).
- **GetUser**: Mengembalikan respons enrichment yang ditambahkan dengan kedua field lama dan baru sesuai ReadMode.
- **GetLegacyUser/GetModernUser**: Mensimulasikan baca klien V1 dan V2 dengan pelacakan observabilitas.
- **ApplyContract**: Memverifikasi kondisi (nol hit legacy kecsaia dipaksa) lalu mengalihkan mode tulis/baca ke NewOnly dan menghapus kolom legacy.
- **BackfillWorker**: Memproses rekaman dalam batch, meng-checkpoint progres, dan memastikan idempoten melalui pengecekan sebelum insert.

### Apa yang Diperlihatkan
1. Migrasi skema struktural tanpa downtime (1:1 ke 1:N).
2. Kompatibilitas mundar-melit JSON serialization yang mengonsumsi field aditif di seluruh klien lama (V1) dan modern (V2).
3. Rollback yang aman tanpa kehilangan data selama fase dual-write.
4. Deteksi drift data melalui reconciliasi.
5. Penegakan fase kontrak tergantung pada verifikasi observabilitas nol traffic legacy.

### Apa yang Tidak Diperlihatkan
- Koordinasi transaksi terdistribusi antar layanan terpisah (misalnya pola Outbox / streaming berbasis Kafka).
- Waktu tunggu kunci DDL PostgreSQL (`SET lock_timeout = '2s'`) dan eksekusi runtime `CREATE INDEX CONCURRENTLY`.

## Apa yang Dibuktikan oleh Tes

Tes unit, integrasi, dan konkurensi memverifikasi perilaku berikut:

**Kompatibilitas Serialisasi Mundar-Melit** (`internal/compat/service_test.go:11-45`): Payload tunggal yang mengandung field `phone` (string) dan `phones` (array) di-unmarshal dengan benar oleh klien V1 (`LegacyConsumerDTO`) dan V2 (`ModernConsumerDTO`) tanpa kesalahan.

**Idempoten & Resumable Backfill** (`internal/compat/service_test.go:47-91`): Worker memproses batch 3 record sekaligus, melanjutkan dari checkpoint, dan menjalankan ulang tidak menghasilkan duplikasi (migrasi kedua = 0).

**Fallback Read dengan Lazy Backfill** (`internal/compat/service_test.go:93-125`): Mode `ReadFallback` mengembalikan data dari `users.phone` saat `user_phones` kosong, dan menuliskannya ke tabel baru secara lazy saat dibaca.

**Reconciliasi & Deteksi Drift** (`internal/compat/service_test.go:127-154`): Sebelum backfill, `ReconcileData` mendeteksi 1 drift; setelah backfill, drift = 0. Counter `DriftDetected` bertambah tepat.

**Header Deprecation & Penegakan Kontrak** (`internal/compat/service_test.go:156-205`): Endpoint V1 mengembalikan `Deprecation: true` dan `Sunset` header. `ApplyContract(false)` gagal saat `LegacyReadHits > 0`; `ApplyContract(true)` berhasil dan endpoint V1 mengembalikan `410 Gone`.

**Siklus Penuh Expand-Migrate-Contract** (`tests/migration_test.go:10-94`): Menjalankan 6 tahap deployment (Baseline → Expand → DualWrite → Backfill → ReadSwitch → Contract) dan memverifikasi klien V2 membaca data historis, kontrak menghapus kolom legacy, baca V1 gagal, baca V2 berhasil.

**Skenario Rollback** (`tests/migration_test.go:96-139`): Rollback selama DualWrite mempertahankan data legacy. Rollback *setelah* `WriteNewOnly` menyebabkan kehilangan data pada skema lama (bukti: legacy phone kosong).

**Keamanan Thread (Race Detector)** (`tests/concurrency_test.go:13-78`): 10 writer, 10 reader, 1 backfill worker, 1 reconciler berjalan paralel 2 detik tanpa race condition (`-race pass`). `DualWriteErrors = 0`.

**Demo End-to-End** (`cmd/demo/main.go`): Menjalankan seluruh siklus hidup secara berurutan dan mengeluarkan snapshot metrik (`legacy_reads: 3, new_reads: 2, dual_writes: 1, backfilled: 2, drift_detected: 0`).

## Pemulihan / Rollback

Strategi rollout dan rollback didasarkan pada prinsip **reversibilitas migrasi**:

1. **Expand Non-Destructive**: Kolom/tabel baru ditambahkan tanpa menyentuh skema lama. Rollback ke versi N selalu aman.
2. **Dual-Write sebagai Jembatan**: Selama `WriteDual`, data ditulis ke kedua skema. Jika rollback terjadi di fase ini, data terbaru tersedia di `users.phone` → klien V1 tetap berfungsi (`tests/migration_test.go:96-118`).
3. **Titik Non-Return**: Setelah `WriteNewOnly` aktif, skema lama tidak lagi diperbarui. Rollback ke versi N *setelah* titik ini berisiko kehilangan data (`tests/migration_test.go:120-138`).
4. **Fitur Flag = Saklar Instan**: `WriteMode` dan `ReadMode` dikontrol via atomic flag. Rollback = `SetWriteMode(WriteLegacyOnly)` + `SetReadMode(ReadLegacyOnly)` tanpa redeploy (dilakukan dalam demo Step 5: `cmd/demo/main.go:87-100`).
5. **Periode Observasi Sebelum Contract**: `ApplyContract(false)` memblokir penghapusan kolom selama `LegacyReadHits > 0`. Hanya boleh dipaksa (`force=true`) setelah jendela observasi nol traffic terpenuhi (`research/07-deployment-and-rollback.md:25-33`).

## Pertimbangan Produksi

Ketika menerapkan pola ini ke produksi, pertimbangkan:

- **Transaksi Dual-Write**: Gunakan transaksi database (`BEGIN ... COMMIT`) atau pola outbox untuk atomic dual-write lintas tabel. Penyimpanan memori menggunakan mutex tunggal (`store.go:53-84`); production membutuhkan atomicity di tingkat DB.
- **Batch Backfill di DB Nyata**: Gunakan cursor ID range (`WHERE id > last_id LIMIT batch`) bukan `OFFSET`. Tambahkan `pg_sleep` atau throttling antar batch. Demo menggunakan `batchSize` konfigurable (`backfill.go:23-32`).
- **Index Non-Blocking**: `CREATE INDEX CONCURRENTLY` pada kolom foreign key `user_phones.user_id` sebelum read switch.
- **Lock Timeout**: Set `SET lock_timeout = '2s'` saat migrasi DDL untuk mencegah blocking DML production.
- **Jendela Observasi Klien Mobile**: API compatibility memerlukan periode legacy lebih lama (bulan-tahun). Pertimbangkan versioned endpoint (`/api/v1`, `/api/v2`) dengan transformer.
- **CDC Alternatif**: Untuk skala tinggi, pertimbangkan Debezium/Kafka logical replication alih-alih dual-write aplikasi (lihat `research-audit/06-gaps.md:Gap 3`).

## Kesalahan Umum

1. **Menghapus kolom legacy terlalu dini** → crash V1 instances (research `08-failure-modes.md:1`).
2. **Beralih read path sebelum backfill selesai tanpa fallback** → data starvation / 404 pada record historis (research `08-failure-modes.md:14-18`).
3. **Dual-write non-atomik tanpa transaksi/outbox** → drift data tersembunyi (research `08-failure-modes.md:8-13`).
4. **Mengabaikan fase Contract** → akumulasi utang teknis permanen (research `08-failure-modes.md:20-23`).
5. **Rollback setelah dual-write dihentikan** → kehilangan data pada skema lama (bukti: `tests/migration_test.go:132-138`).
6. **Backfill tanpa checkpoint/idempotency** → duplikasi atau kehilangan data saat restart.

## Studi Kasus

### Studi Kasus A — Telepon Pelanggan (1:1 ke 1:N)
**Awal**: `customers (id, name, phone)`
**Target**: `customer_phone_numbers (id, customer_id, phone_number)`
Strategi: Expand tabel baru → DualWrite ke keduanya → Backfill historis → Switch read → Contract drop `customers.phone` (`research/09-case-studies.md:Case Study A`).

### Studi Kasus B — Mekanik Invoice CMMS (1:1 ke N:M)
**Awal**: `invoice (id, mechanic_id)`
**Target**: `invoice_mechanics (invoice_id, mechanic_id)`
Strategi: Expand tabel junction → DualWrite (simpan mekanik pertama ke kolom legacy, semua ke tabel baru) → Backfill → Read switch dengan transform array ke item pertama untuk V1 → Contract (`research/09-case-studies.md:Case Study B`).

### Studi Kasus C — Multi-Mata Uang (Pemecahan Skema)
**Awal**: `products (id, name, price)`
**Target**: `product_prices (id, product_id, currency, amount)`
Strategi: Expand → Backfill asumsikan mata uang dasar (USD) → API compatibility: endpoint lama tetap mengembalikan `.price` via transformer pull dari `product_prices` → Contract kolom fisik opsional jika transformer tetap diperlukan (`research/09-case-studies.md:Case Study C`).

## Checklist Deployment

- [ ] Skema Expand: tabel/kolom baru nullable, default, atau `CREATE INDEX CONCURRENTLY`.
- [ ] Deploy N+1: WriteMode = `DualWrite`, ReadMode = `LegacyOnly`/`FallbackRead`.
- [ ] Backfill: jalankan batch kecil, checkpoint, verifikasi idempoten, throttling.
- [ ] Reconciliation: `ReconcileData() == 0` sebelum read switch.
- [ ] Read Switch: Feature flag `ReadMode = ReadNewOnly`, canary 1% → 100%.
- [ ] Observasi: `LegacyReadHits == 0` selama periode penuh siklus bisnis (mis. 30+ hari).
- [ ] Contract: `ApplyContract(false)` → verifikasi → drop kolom legacy.
- [ ] Rollback Plan: uji rollback di setiap fase sebelum Contract.

## Sumber

- Riset Disetujui: `research/11-final-research.md`, `research/03-core-concepts.md` sampai `research/09-case-studies.md`
- Audit Riset: `research-audit/07-verdict.md` (APPROVED)
- Desain Engineering: `engineering/01-design.md`, `engineering/02-implementation-notes.md`
- Audit Engineering: `engineering-audit/06-verdict.md` (APPROVED)
- Implementasi: `internal/compat/*.go`, `cmd/demo/main.go`
- Tes: `internal/compat/service_test.go`, `tests/migration_test.go`, `tests/concurrency_test.go`
- Demo: `engineering/03-execution-result.md`