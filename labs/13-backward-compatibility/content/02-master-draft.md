# Panduan Kompatibilitas Mundur: Implementasi Pola Expand → Migrate → Contract untuk Evolusi Skema dan API

## Masalah
Sistem produksi sering menghadapi kebutuhan untuk mengubah kontrak data—misalnya mengubah relasi 1:1 (satu nomor telepon per pengguna) menjadi 1:N (banyak nomor telepon per pengguna). Perubahan langsung seperti `DROP COLUMN phone` atau mengubah struktur JSON API dapat menyebabkan:
- Crash aplikasi bagi klien lama yang masih mengandalkan struktur lama
- Downtime selama migrasi skema berat
- Risiko kehilangan data jika rollback dilakukan setelah kontrak lama dihapus

Lab ini membuktikan bahwa perubahan struktural dapat dilakukan dengan aman tanpa downtime atau merusak kompatibilitas dengan menggunakan pola **Expand → Migrate → Contract** (Parallel Change Pattern).

## Mengapa Ini Penting
Kompatibilitas mundur adalah fondasian layanan yang dapat diandalkan. Tanpa strategi ini:
- Tim engineering terpaksa menjadwalkan downtime yang mahal
- Klien eksternal (mobile app, mitra API) mungkin tidak pernah mengupdate, menciptakan beban dukungan permanen
- Setiap perubahan skema menjadi kegiatan berisiko tinggi yang menghambat kecepatan pengiriman fitur

Dengan parallel change, organisasi dapat:
- Melatih deploy-anytime kultur melalui continuous delivery
- Menjaga pengalaman pengguna selama transisi infrastruktur
- Mengurangi beban operasional dari pemeliharaan versi ganda

## Model Mental
Pola parallel change membagi perubahan yang traditionally breaking menjadi tiga fase yang masing-masing non-breaking:
1. **Expand**: Tambahkan representasi baru bersama-sama dengan lama tanpa mengubah yang sudah ada
2. **Migrate**: Secara gradual salurkan traffic dan data dari representasi lama ke baru
3. **Contract**: Hapus representasi lama hanya setelah penggunaan baru mencapai 100% dan penggunaan lama mencapai 0%

Kunci psikologis: tidak ada fase yang mengharuskan koordinasi simultan antara penyedia layanan dan semua konsumen. Sistem dapat beroperasi dalam salah satu fase tiganya secara independen.

## Konsep Inti

### Expand (Perluasan)
Fase pertama menambah kemampuan baru tanpa menyentuh yang sudah ada:
- Database: Buat tabel/kolom baru (misal: `user_phones` bersama-sama dengan kolom lama `users.phone`)
- API: Tambahkan field baru secara aditif (misal: field `phones` array bersama-sama dengan field lama `phone`)
- Payload: Respons JSON menjadi enriched—berisi baik field lama maupun field baru

Klien lama tetap bekerja karena:
- Mereka tidak melihat perubahan apa-apa (kolom/field lama tetap ada)
- Mereka mengabaikan field baru yang tidak dikenal (perilaku standar JSON parser)

### Migrate (Migrasi)
Fase tengah menangkap aliran data dan traffic:
- **Dual-write**: Setiap tulisan baru dilakukan ke kedua skema (legacy dan modern) secara atomik
- **Backfill**: Worker latar belakang meng-copy data historis dari skema lama ke baru dalam batch kecil dengan checkpoint resumable
- **Fallback read (dual-read)**: Baca dari skema baru terlebih dahulu; jika kosong/null, baca dari skema lama dan secara opsional lakukan lazy backfill
- **Migrasi konsumen**: Tim klien (mobile, backend internal) diupdate untuk menggunakan struktur baru

Fase ini adalah yang terpanjang karena tergantung pada adopsi eksternal.

## Struktur Implementasi
Lab ini mengimplementasikan pola ini dengan arsitektur berikut:

```
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

### Komponen Kunci
1. **Storage**: Implementasi thread-safe in-memory yang mensimulasikan tabel relasional dengan dukungan untuk dual-write dan operasi drop kolom
2. **FeatureFlagManager**: Mengontrol tiga mode tulis (`WriteLegacyOnly`, `WriteDual`, `WriteNewOnly`) dan tiga mode baca (`ReadLegacyOnly`, `ReadFallback`, `ReadNewOnly`)
3. **MetricsCollector**: Melacak counter untuk traffic lama/baru, kesalahan dual-write, progres backfill, dan deteksi drift data
4. **BackfillWorker**: Worker yang dapat dihentikan dan dilanjutkan lagi dengan idempotent melalui checkpoint `last_processed_id`
5. **CompatService**: Facade domain yang mengorkestrasi operasi baca/tulis, rekonsiliasi data, dan siklus deprecasi
6. **HTTPHandler**: Mengembalikan kontrak JSON dengan header Deprecation dan Sunset sesuai RFC 8594

## Arsitekural Detail

### Dual-Write dalam Transaksi
Pada mode `WriteDual`, setiap operasi CREATE:
1. Menulis ke tabel `users` (kolom `phone` legacy)
2. Menulis ke tabel `user_phones` (record dengan `is_primary = true`)
3. Kedua operasi dilakukan dalam kritial section yang sama menggunakan `sync.RWMutex`

Ini menjamin konsistensi data selama migrasi—jika sistem crash, baik keduanya diperbarui atau keduanya tidak berubah.

### Mekanisme Backfill
BackfillWorker bekerja dengan pola:
1. Mengambil chunk ID pengguna (misal: 100 ID per batch) mulai dari `last_processed_id`
2. Untuk setiap ID, memeriksa apakah record legacy memiliki `phone` dan apakah `user_phones` kosong untuk pengguna tersebut
3. Jika ya, melakukan insert idempotent ke `user_phones` (menggunakan UPSERT logic via pengecekan duplikasi)
4. Menerapkan batasan waktu antar batch untuk mencegah kontensi database
5. Melanjutkan dari checkpoint terakhir jika proses dihentikan

### Fallback Read
Pada mode `ReadFallback`, operasi GET:
1. Mencoba membaca dari `user_phones` terlebih dahulu
2. Jika kosong dan mode adalah `ReadFallback`, maka baca dari `users.phone`
3. Jika berhasil, secara opsional menulis kembali ke `user_phones` (lazy backfill) untuk mengurangi beban bacaan di kemudian hari

## Kode yang Diverifikasi

### Model Data dan DTO
File: `internal/compat/model.go`

```go
type User struct {
    ID        int
    Name      string
    Phone     *string // Legacy field; nil once contracted
    CreatedAt time.Time
}

type UserResponse struct {
    ID     int          `json:"id"`
    Name   string       `json:"name"`
    Phone  string       `json:"phone,omitempty"` // Legacy field maintained for v1 consumers
    Phones []PhoneEntry `json:"phones"`          // New field for v2 consumers
}
```

Komentar: Field `Phone` menggunakan pointer `*string` sehingga dapat menjadi `nil` setelah kontrak diterapkan, menunjukkan kolom yang telah di-drop. Field `Phones` selalu berisi array, memenuhi kontrak konsumen V2.

### Handler HTTP dengan Header Deprecasi
File: `internal/compat/handler.go`

```go
func (h *APIHandler) GetUserV1(w http.ResponseWriter, r *http.Request) {
    // ... logic untuk mendapatkan data legacy ...
    
    // Inject Deprecation Headers
    w.Header().Set("Deprecation", "true")
    w.Header().Set("Sunset", "Mon, 31 Dec 2026 23:59:59 GMT") // Sunset date
    w.Header().Set("Content-Type", "application/json")
    
    json.NewEncoder(w).Encode(dto)
}
```

Komentar: Header `Deprecation: true` dan `Sunset: <tanggal>` sesuai RFC 8594 memberi tahu klien bahwa endpoint ini akan dihapus dan kapan hal tersebut akan terjadi.

### Mekanisme Kontrakt dengan Penghitung Lalu Lintas
File: `internal/compat/service.go`

```go
func (s *Service) ApplyContract(force bool) error {
    // Guard: Pastikan tidak ada trafic legacy kecuali jika dipaksa
    if !force && s.obs.LegacyReadHits.Load() > 0 {
        return fmt.Errorf("%w: recorded %d legacy reads", ErrContractViolation, s.obs.LegacyReadHits.Load())
    }
    
    // 1. Beralih mode tulis ke NewOnly
    s.flags.SetWriteMode(WriteNewOnly)
    // 2. Beralih mode baca ke NewOnly  
    s.flags.SetReadMode(ReadNewOnly)
    // 3. Tandai kontrak sebagai diterapkan
    s.flags.SetContractApplied(true)
    // 4. Drop kolom legacy dari storage
    s.store.ApplyContractDropLegacyColumn()
    
    return nil
}
```

Komentar: Kontrakt hanya dapat diterapkan ketika meter `LegacyReadHits` menunjukkan nol, memastikan tidak ada klien lama yang masih aktif.

### Backfill yang Idempotent dan Resumable
File: `internal/compat/backfill.go`

```go
func (b *BackfillWorker) RunBatch(ctx context.Context) (int, bool, error) {
    b.checkpoint.mu.Lock()
    defer b.checkpoint.mu.Unlock()

    if b.checkpoint.IsComplete {
        return 0, true, nil
    }

    ids := b.store.GetUserIDs(b.checkpoint.LastProcessedID, b.batchSize)
    if len(ids) == 0 {
        b.checkpoint.IsComplete = true
        return 0, true, nil
    }

    migratedInBatch := 0
    for _, id := range ids {
        // ... proses backfill untuk setiap ID ...
        
        // Idempotency check: hanya proses jika belum ada di user_phones
        if user.Phone != nil && *user.Phone != "" {
            phones, _ := b.store.GetPhones(id)
            if len(phones) == 0 {
                // Insert hanya jika belum ada entri
                _, err := b.store.SavePhoneEntry(id, *user.Phone, true)
                // ... error handling ...
                migratedInBatch++
            }
        }
        b.checkpoint.LastProcessedID = id
    }

    // ... update checkpoint dan cek selesai ...
}
```

Komentar: Idempotency dicapai dengan memeriksa keberadaan entri sebelum insert. Resumabilitas dicapai dengan menyimpan `LastProcessedID` dalam struktur yang dilindungi mutex.

## Apa yang Dibuktikan oleh Tes

### Serialisasi Kompatibel Mundur
File: `internal/compat/service_test.go` (TestSerializationBackwardCompatibility)
- Membuat user dengan nomor telepon utama dan tambahan
- Meng-serialize ke format enriched JSON
- Memverifikasi bahwa klien V1 (menggunakan LegacyConsumerDTO) dapat deserialize tanpa crash dan mendapatkan nomor telepon utama
- Memverifikasi bahwa klien V2 (menggunakan ModernConsumerDTO) dapat deserialize dan mendapatkan array nomor telepon lengkap

### Backfill Idempotent dan Resumable
File: `internal/compat/service_test.go` (TestBackfillIdempotentAndResumable)
- Membuat 10 pengguna legacy
- Mengjalankan backfill dengan batch size 3
- Memverifikasi bahwa batch pertama memproses 3 record, batch kedua memproses 3 record, dan sisa memproses 4 record
- Memverifikasi bahwa menjalankan backfill lagi menghasilkan 0 migrasi (idempotent)

### Mekanisme Fallback Read
File: `internal/compat/service_test.go` (TestFallbackRead)
- Membuat satu pengguna legacy
- Mengatur mode baca ke `ReadFallback`
- Memverifikasi bahwa sebelum baca, tabel modern kosong
- Memverifikasi bahwa operasi GET mengembalikan nomor telepon dari field legacy
- Memverifikasi bahwa setelah baca, tabel modern telah terisi secara lazy (nomor telepon disalin ke user_phones)

### Deteksi Drift Data dan Rekonsiliasi
File: `internal/compat/service_test.go` (TestDataReconciliationAndDrift)
- Membuat pengguna legacy tanpa menjalankan backfill
- Memverifikasi bahwa rekonsiliasi mendeteksi 1 drift (data ada di legacy tetapi tidak di modern)
- Menjalankan backfill hingga selesai
- Memverifikasi bahwa rekonsiliasi setelah backfill menunjukkan 0 drift

### Penegakan Kontrakt dan Header Deprecasi
File: `internal/compat/service_test.go` (TestDeprecationHeadersAndContractEnforcement)
- Memverifikasi bahwa endpoint legacy (`/v1/users`) mengembalikan header `Deprecation: true` dan `Sunset`
- Memverifikasi bahwa kontrak tidak dapat diterapkan (`ApplyContract(false)`) ketika masih ada trafic legacy
- Memverifikasi bahwa kontrak paksa (`ApplyContract(true)`) berhasil
- Memverifikasi bahwa setelah kontrakt, endpoint legacy mengembalikan HTTP 410 Gone
- Memverifikasi bahwa endpoint modern (`/v2/users`) tetap berfungsi normal setelah kontrakt

### Keseluruhan Siklus Hidup
File: `tests/migration_test.go` (TestFullExpandMigrateContractLifecycle)
- Memulai dengan data historis legacy
- Fase Expand: mengaktifkan dual-write dan membuat pengguna baru
- Memverifikasi bahwa data baru ditulis ke kedua skema
- Fase Migrate: menjalankan backfill hingga selesai
- Memverifikasi bahwa nol drift terdeteksi sebelum pergantian jalur baca
- Pergantian jalur baca ke `ReadNewOnly`
- Memverifikasi bahwa klien modern dapat membaca data historis dari skema baru
- Fase Contract: menerapkan kontrak dengan paksa
- Memverifikasi bahwa baca legacy setelah kontrak menghasilkan error
- Memverifikasi bahwa baca modern setelah kontrak terus berfungsi dengan data lengkap

### Keamanan Rollback
File: `tests/migration_test.go` (TestRollbackScenarios)
- **Skenario A (Rollback Aman selama Dual-Write)**:
  - Membuat pengguna selama mode dual-write
  - Mengganti kembali ke mode legacy-only (menyimulasikan rollback dari N+1 ke N)
  - Memverifikasi bahwa klien legacy masih dapat membaca nomor telepon yang dibuat selama dual-write (tidak ada kehilangan data)
  
- **Skenario B (Rollback Berhenti Dual-Write Terlalu Awal)**:
  - Membuat pengguna selama mode new-only (setelah dual-write dihentikan)
  - Mengganti kembali ke mode legacy-only
  - Memverifikasi bahwa klien legacy menerima string kosong untuk nomor telepon (mendemonstrasikan kehilangan data jika rollback dilakukan setelah dual-write dihentikan)

### Keseluruhan Concurrency
File: `tests/concurrency_test.go` (TestConcurrency)
- Menjalankan 10 goroutine penulis legacy bersamaan
- Menjalankan 10 goroutine pembaca legacy dan modern bersamaan  
- Menjalankan 1 goroutine worker backfill
- Menjalankan 1 goroutine perekoncilian drift
- Semua berjalan dengan timeout 2 deteksi
- Memverifikasi bahwa tidak ada kesalahan dual-write yang dilaporkan di bawah beban bersamaan

## Demonstrasi End-to-End
File: `cmd/demo/main.go`

Demo menunjukkan seluruh siklus hidup:
1. **Baseline**: Membuat pengguna historis Alice dan Bob; klien V1 membaca nomor telepon mereka
2. **Expand**: Mengaktifkan dual-write; membuat pengguna Charlie dengan dua nomor telepon; klien V1 masih bisa membaca nomor pertama
3. **Migrate**: Menjalankan backfill untuk data historis; rekonsiliasi menunjukkan nol drift
4. **Switch Read Path**: Mengalihkan mode baca ke NewOnly; klien V2 membaca data historis Alice dari skema baru
5. **Rollback Demo**: Mensimulasikan rollback ke Version N selama dual-write; klien V1 masih bisa membaca nomor telepon Charlie (bukti keamanan rollback)
6. **Contract**: Menerapkan kontrak; klien V1 menerima error 410 Gone; klien V2 terus bekerja dengan data lengkap

## Pertimbangan Produksi

### Keamanan Rollback
Keamanan rollback hanya terjamin selama fase dual-write. Jika dual-write dihentikan (beralih ke `WriteNewOnly`) kemudian rollback dilakukan ke Version N, data yang ditulis selama fase new-only akan hilang karena tidak pernah ditulis ke skema legacy.

### Deteksi Konsumen Lama
Sistem mengandalkan metrik untuk mengetahui ketika kontrak aman untuk diterapkan:
- Meter `LegacyReadHits` harus mencapai nol dan tetap pada nol selama jendela observasi
- Dalam produksi, ini biasanya dicapai dengan memonitor traffic selama 30 hari setelah traffic legacy menunjukkan penurunan signifikan
- Header Deprecasi dan Sunset memberikan waktu untuk klien eksternal melakukan migrasi

### Kompleksitas dan Overhead
- **Latensi Tulis**: Dual-write meningkatkan latensi tulis karena setiap operasi harus menyelesaikan dua penulisan
- **Kompleksitas Kode**: Perlu menjaga dua jalur baca dan logika sinkronisasi data
- **Beban Storage**: Diperlukan penyimpanan ganda selama fase migrasi (legacy + modern)
- **Beban Operational**: Perlu monitoring terus-meterus untuk metrik lalu lintas dan drift data

### Keterbatasan Implementasi Demonstrasi
Lab ini menggunakan beberapa penyederhanaannya untuk tujuan edukasi:
- **Penyimpanan Dalam-Memori**: Menggunakan `sync.RWMutex` dengan peta alih daripada PostgreSQL aktual (ditandai dengan komentar `ponytail: in-memory mock storage; replace with database/sql for persistent store.`)
- **Emulasi DDL**: Perintah PostgreSQL spesifik seperti `SET lock_timeout` dan `CREATE INDEX CONCURRENTLY` hanya didokumentasikan dalam `schema.sql` dan penelitian, bukan dieksekusi runtime karena penggunaan penyimpanan dalam-memori
- **Skala Transaksi**: Dual-write diimplementasikan dalam kritial section tunggal alih-alih transaksi database yang sebenarnya

## Kesalahan Umum yang Harus Dihindari

### 1. Mengabaikan Fase Contract
Menghgalkan kontrakt menyebabkan:
- Teknical debt yang terus-menerus menumpuk
- Pemborongan sumber dari memejaga dua skema untuk selamanya
- Kebingungan tentang mana skema yang "sebenarnya" digunakan

Mitasi: Buat tiket kontrak sebagai bagian dari setiap epic expand, dan gunakan metrik usage untuk memicu peringatan ketika penggunaan legacy = 0 selama periode observasi.

### 2. Beralih Jalan Baca Terlalu Awal
Beralih ke `ReadNewOnly` sebelum backfill selesai menyebabkan:
- Kelahiran data (404/respons kosong) untuk rekaman yang belum terbackfill
- Pengalaman pengguna yang buruk untuk klien yang beralih terlebih dahulu

Mitasi: Gunakan mode `ReadFallback` selama fase backfill, atau Pastikan backfill 100% selesai sebelum beralih ke `ReadNewOnly`.

### 3. Mengasumsikan Idempotensi Tanpa Pemeriksaan
Worker backfill yang tidak memeriksa duplikasi dapat:
- Membuat entri ganda jika diulang
- Menyebabkan pelanggaran constraint unik jika ada
- Memboroskan sumber I/O dengan operasi yang tidak perlu

Mitasi: Selalu terapkan logika idempotent (misal: periksa keberadaan sebelum insert, atau gunakan sintaks UPSERT/ON CONFLICT DO NOTHING).

### 4. Mengabaikan Overhead Dual-Write dalamencapaian Skala
Dual-write yang tidak di-throttle pada sistem dengan volume tinggi dapat:
- Menyebabkan kontensi database yang signifikan
- Meningkatkan latensi tulis secara berlebihan
- Mengkonsumsi lebih banyak koneksi pool database

Mitasi: Implementasikan dual-write melalui pola outbox (tulis ke tabel tunggal, lalu proses async ke tujuan kedua) untuk sistem dengan skala tinggi.

## Ringkasan
Pola Expand → Migrate → Contract memberikan rangka yang teruji untuk mengubah kontrak data dan skema secara bertahap tanpa downtime atau memutus klien lama. Lab ini membuktikan bahwa dengan kombinasi dual-write, backfill yang dapat dihentikan dan dilanjutkan lagi, fallback baca, dan pengontrol lintas lalu lintas berbasis fitur flag, organisasi dapat:
- Melakukan migrasi skema besar selama operasi normal
- Menjaga kompatibilitas dengan klien lama yang tidak bisa langsung diupdate  
- Memastikan keamanan rollback selama fase transisi
- Mengukur dan memverifikasi ketika aman untuk menghapus representasi lama

Kunci keberhasilan terletak pada memperlakukan perubahan infrastruktur sebagai proses pengalihan traffic gradual alih-alih kejadian yang bersifat putus-putus, dengan metrik dan observasi sebagai panduan untuk mengambil keputusan pada setiap fase transisi.