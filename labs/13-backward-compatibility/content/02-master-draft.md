# Panduan Kompatibilitas Mundur: Pola Expand → Migrate → Contract untuk Evolusi Skema dan API Tanpa Downtime

## Problem

Sistem produksi sering memerlukan perubahan struktural pada data — misalnya mengubah relasi
`1:1` (`users.phone`, satu nomor per pengguna) menjadi `1:N` (`user_phones`, banyak nomor per
pengguna), atau mengubah satu field string `phone` menjadi array `phones` pada kontrak JSON.
Dilakukan sekaligus, perubahan ini *breaking*:

- **Crash klien lama**: klien V1 yang mem-parse `phone` akan rusak jika field dihapus atau tipe diganti.
- **Downtime migrasi**: backfill tabel penuh dalam satu transaksi bisa menguncang blokir selama berjam-jam.
- **Data loss pada rollover**: jika kolom legacy sudah dihapus, rollback ke versi lama menyebabkan data
  yang ditulis *setelah* transisi menghilang bagi V1.

Tanpa strategi, setiap evolusi skema menjadi *release-coordinated* — memaksa semua klien dan semua
instance aplikasi berkoordinasi serentak, yang justru meningkatkan risiko kegagalan.

Lab ini membuktikan (8/8 tes lolos termasuk `-race`, demo berjalan) bahwa transisi
`1:1 → 1:N` dapat dilakukan dengan aman dengan tiga fase bertahap yang masing-masing
*non-breaking*.

## Why This Matters

Kompatibilitas mundur adalah fondasi SLO service yang dapat diandalkan. Tanpa pendekatan
inkremental:

- Tim terpaksa menjadwalkan *maintenance window* mahal untuk setiap migrasi skema.
- Klien eksternal (aplikasi mobile, mitra API) yang tidak pernah upgrade menciptakan
  *technical debt* permanen dan beban dukungan operasional.
- Setiap perubahan skema menjadi *high-sca" change* yang menghambat kecepatan pengiriman
  fitur di bawah prinsip *Continuous Delivery*.

Dengan pola *parallel change*, organisasi dapat *deploy-any-phase* — kode baru dapat
dideploy ke produksi bahkan hanya untuk satu fase, karena setiap fase dirancang untuk
tidak memutus klien lama.

## Mental Model

Satu breaking change dipecah menjadi tiga fase yang masing-masing *independently deployable*:

```text
EXPAND   ->   MIGRATE   ->   CONTRACT
(tambah       (pindah      (hapus lama)
 keduanya)    data/traffic)
```

- **Expand**: tambahkan representasi baru **di samping** representasi lama. Tidak ada
  yang berkurang — klien lama tidak menyadari adanya perubahan.
- **Migrate**: alirkan data historis dan traffic ke representasi baru secara bertahap.
  Selama ini, sistem menulis dan membaca dari **kedua** representasi.
- **Contract**: hapus representasi lama **hanya ketika** traffic legacy sudah nol.

Kunci psikologis: tidak ada fase yang mengharuskan koordinasi simultan antara penyedia
layanan dan semua konsumen. Sistem dapat berada di fase mana saja secara independen,
bahkan dengan versi aplikasi yang berbeda berjalan bersamaan (rolling deployment).

## Core Concept

### Expand — Tambahkan tanpa mengubah yang ada

Pada fase ini, sistem belajar "berbicara dua bahasa" sekaligus:

- **Schema**: buat tabel/kolom baru (`user_phones`) yang *nullable*, tanpa menyentuh
  kolom lama (`users.phone`). Pada produksi PostgreSQL ini berarti
  `ALTER TABLE ... ADD COLUMN ... DEFAULT` yang tidak menulis ulang baris.
- **API/Payload**: tambahkan field baru secara aditif (`phones` array) bersamaan
  dengan field lama (`phone`). Respons JSON menjadi *enriched*.

Klien lama tetap bekerja karena:
1. Field lama (`phone`) masih ada dan diisi.
2. Parser JSON standar meng-*ignore*-kan field baru (`phones`) yang tidak dikenal.

#### Bukti (Martin Fowler, *Parallel Change*):

> "In the *expand* phase you augment the interface to support both the old and the new
> versions. In our example, we introduce a new `Map<Coordinate, Cell>` data structure and
> the new methods that can receive `Coordinate` instances without changing the existing code."

### Migrate — Alihkan data dan traffic

Fase ini berlangsung paling lama karena bergantung pada adopsi klien eksternal:

- **Dual-write**: setiap tulisan baru ditulis ke kedua skema. Pada lab ini dilakukan
  *atomic* dalam satu *critical section* `sync.RWMutex`; pada produksi memerlukan transaksi
  database tunggal atau pola *transactional outbox*.
- **Backfill**: worker batch menyalin data historis dari skema lama ke baru, dengan
  *checkpoint* `last_processed_id` dan logika *idempotent* (cek keberadaan sebelum insert).
- **Fallback read**: baca dari skema baru dulu; jika kosong, baca dari skema lama dan
  *lazy backfill* ke skema baru.

Feature flag memisahkan *deployment* dari *aktivasi*: `WriteMode` dan `ReadMode` dapat
diubah runtime tanpa redeploy, memungkinkan *canary* 1% → 10% → 100%.

### Contract — Hapus lama hanya pada traffic nol

Setelah semua klien bermigrasi dan `LegacyReadHits` = 0:
1. Berhenti menulis ke skema lama (`WriteNewOnly`).
2. Hapus *code path* legacy dan *adapter transform*.
3. Drop kolom/tabel lama (`ALTER TABLE users DROP COLUMN phone`).

#### Bukti (Martin Fowler):

> "Once all usages have been migrated to the new version, you perform the *contract* phase
> to remove the old version."
> "If the contract phase is not executed you might end up in a worse state than you started."

## Failure Scenario

Jika fase dilewati atau urutannya salah, kegagalan sistematis terjadi:

| # | Failure Mode | Dampak | Penyebab |
|---|-------------|--------|----------|
| 1 | Destructive alter | klien V1 crash / 500 | `DROP COLUMN phone` sebelum semua klien migrasi |
| 2 | Dual-write drift | data antar skema berbeda | write tak atomic atau outbox tidak idle |
| 3 | Missing backfill | data starvation | baca ke skema baru sebelum semua legacy ter-backfill |
| 4 | Abandoned expand | teknical debt / kemusatan | lupa fase Contract |
| 5 | Type-change corruption | nil/err salah tipe | `phone TEXT` → `phone INT` (non-additive) |
| 6 | Premature rollback | data loss | rollback setelah `WriteNewOnly` berhenti |

Pada produksi, semua failure mode ini berpotensi menyebabkan *customer-affecting incident*.
Lab ini tidak melakukan failure mode #5 (type change) karena tidak bagian dari 1:1→1:N,
tapi mem-buktikan #3 (fallback read mencegahnya) dan #6 (Skenario B pada tes).

## How It Works

```text
┌──────────────────────────────────────────────────────────────┐
│                    HTTP / Service Layer                       │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │  Feature Flags (WriteMode, ReadMode, ContractApplied)   │  │
│  │  Observability (LegacyReadHits, NewReadHits, ...)       │  │
│  └─────────────────────────────────────────────────────────┘  │
└──────┬───────────────────────────────┬─────────────┬─────────┘
       │                               │             │
       ▼                               ▼             ▼
┌──────────────┐              ┌──────────────┐  ┌──────────┐
│  /api/v1/users  │              │  /api/v2/users │  │ Backfill │
│  (LegacyConsumerDTO)│              │  (ModernConsumerDTO)│  │ Worker   │
└──────────────┘              └──────────────┘  └────┬─────┘
       │                                    │       │
       │  Deprecation: true / Sunset date   │       │
       ▼                                    ▼       │
┌────────────────────────────────────────────────────│────────┐
│             MemoryStore (simulasi relasional)         │        │
│  users(id, name, phone*)        user_phones(id,    │        │
│                            │ user_id FK, │        │
│                            │ number, is_primary) │        │
└────────────────────────────────────────────────────┴────────┘
```

**Lifecycle 6-fase** (demo `cmd/demo/main.go` + tes `tests/migration_test.go`):

1. **Baseline**: `WriteLegacyOnly` + `ReadLegacyOnly`. V1 baca/tulis `users.phone`.
2. **Expand**: deploy kode yang mendukung `user_phones`. `WriteDual` menulis ke keduanya.
3. **Migrate — Dual-write**: semua write baru ke kedua skema. Payload *enriched*.
4. **Migrate — Backfill**: worker menyalin data historis ke `user_phones`.
5. **Read switch**: `ReadNewOnly` — V2 baca dari skema baru. V1 masih baca dari `phone`.
6. **Contract**: `ApplyContract` setelah `LegacyReadHits == 0`. Drop kolom legacy. V1 dapat `410 Gone`.

## Architecture

Empat komponen inti di `internal/compat/`:

- **`store.go`** — `MemoryStore` dengan `sync.RWMutex`, tabel `users` (legacy) dan
  `user_phones` (modern), dual-write atomik via `CreateDual`, drop kolom via
  `ApplyContractDropLegacyColumn`.
- **`flags.go`** — `FeatureFlags` dengan `atomic.Value`/`atomic.Bool` untuk
  `WriteMode` (`WriteLegacyOnly`/`WriteDual`/`WriteNewOnly`),
  `ReadMode` (`ReadLegacyOnly`/`ReadFallback`/`ReadNewOnly`),
  dan `ContractApplied`.
- **`metrics.go`** — `Observability` counter atomik: `LegacyReadHits`,
  `NewReadHits`, `DualWriteCount`, `DualWriteErrors`, `BackfillProcessed`,
  `DriftDetected`.
- **`backfill.go`** — `BackfillWorker` resumable + idempotent via
  `BackfillCheckpoint.LastProcessedID`.
- **`service.go`** — `CompatService` facade: `CreateUser`, `GetUser`,
  `GetLegacyUser`, `GetModernUser`, `ReconcileData`, `ApplyContract`.
- **`handler.go`** — `APIHandler` HTTP `/api/v1/users` + `/api/v2/users`.

> **Boundary yang disadari**: lab ini in-memory (bukan PostgreSQL). `CREATE INDEX CONCURRENTLY`,
> `NOT VALID`, lock-timeout tidak dieksekusi — hanya didokumentasikan di `schema.sql` dan
> `research/04-database-migration.md`. Dual-write memakai mutex, bukan transaksi DB. Ini
> ilustrasi yang disederhanakan untuk edukasi; *swapping* ke `database/sql` adalah jalur upgrade
> yang tersedia.

## Implementation

### Model & DTO — payload *enriched* bersatu

File: `internal/compat/model.go`

```go
type User struct {
    ID        int
    Name      string
    Phone     *string // Legacy field; nil once contracted
    CreatedAt time.Time
}

// UserResponse is an enriched additive response supporting both legacy and modern clients
type UserResponse struct {
    ID     int          `json:"id"`
    Name   string       `json:"name"`
    Phone  string       `json:"phone,omitempty"` // Legacy field maintained for v1 consumers
    Phones []PhoneEntry `json:"phones"`          // New field for v2 consumers
}
```

`Phone` memakai `*string` agar dapat `nil` setelah kontrak diterapkan (simulasi kolom
yang *tadi dipaksa*). `UserResponse` mengekspor field lama maupun baru — ini inti
kompatibilitas mundur/aditif: klien lama pakai `phone`, klien baru pakai `phones`.

### Dual-write atomic

File: `internal/compat/store.go` (`CreateDual`)

```go
func (s *MemoryStore) CreateDual(name, phone string) (*User, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    // ... write both users.Phone and user_phones in the same critical section ...
}
```

Kedua penulisan terjadi dalam satu *critical section* yang sama — jika crash di tengah,
atau keduanya masuk atau keduanya tidak. Isolasi `sync.RWMutex` menggantikan transaksi DB
pada lab ini. Audit mencatat: dual-write atomik *within the store lock* (engineering-audit-opensource/02-code-audit.md, Finding 1).

### Backfill resumable + idempotent

File: `internal/compat/backfill.go` (`RunBatch`)

```go
ids := b.store.GetUserIDs(b.checkpoint.LastProcessedID, b.batchSize)
// ... idempotent: only SavePhoneEntry if GetPhones is empty ...
b.checkpoint.LastProcessedID = id
```

Idempotensi pada dua lapisan:
1. **Storage** (`SavePhoneEntry`): cek nomor sudah ada → return existing entry, tidak insert duplikat.
2. **Checkpoint** (`LastProcessedID`): lanjut dari ID terakhir yang selesai diproses.

Tes `TestBackfillIdempotentAndResumable` membuktikan: 10 record, batch=3 → 3+3+4,
run kedua = 0 migrasi.

### Fallback read + lazy backfill

File: `internal/compat/service.go` (`GetUser`, mode `ReadFallback`)

```go
case ReadFallback:
    if len(phones) > 0 { /* baca modern, pilih primary */ }
    else if u.Phone != nil {
        phoneVal = *u.Phone
        _, _ = s.store.SavePhoneEntry(id, phoneVal, true) // lazy backfill
    }
```

Mencegah *data starvation* pada record yang belum terbackfill. Tes
`TestFallbackRead` memverifikasi: sebelum baca, `user_phones` kosong; sesudah
fallback GET, modern store terisi dengan satu entry.

### Contract dengan traffic guard

File: `internal/compat/service.go` (`ApplyContract`)

```go
if !force && s.obs.LegacyReadHits.Load() > 0 {
    return fmt.Errorf("%w: recorded %d legacy reads", ErrContractViolation, ...)
}
s.flags.SetWriteMode(WriteNewOnly)
s.flags.SetReadMode(ReadNewOnly)
s.flags.SetContractApplied(true)
s.store.ApplyContractDropLegacyColumn()
```

Guard memakai counter kumulatif `LegacyReadHits`. Lab dokumentasi pelariangan ini
melalui parameter `force=true` ("simulating after 30-day zero traffic window"). Pada
produksi, counter ini perlu *sliding window* reset — audit mencatat ini sebagai
gap LOW, bukan bug kebenaran (engineering-audit-opensource/05-gaps.md, Gap 3).

### Header Deprecasi

File: `internal/compat/handler.go`

```go
w.Header().Set("Deprecation", "true")
w.Header().Set("Sunset", "Mon, 31 Dec 2026 23:59:59 GMT")
```

`Deprecation` memindahkan ke RFC 9224; `Sunset` ke RFC 8594. README menyatukan keduanya
sebagai RFC 8594 — catatan akurasi LOW (engineering-audit-opensource/05-gaps.md, Gap 2).
Setelah kontrak, V1 mendapatkan `410 Gone`.

## Code Walkthrough

Demo `cmd/demo/main.go` mengeksekusi seluruh siklus dalam 6 langkah:

1. **Baseline**: `CreateUser("Alice","+62811111111")`, `CreateUser("Bob","+62822222222")`
   — klien V1 baca via `GetLegacyUser`.
2. **Expand**: `SetWriteMode(WriteDual)`, `CreateUser("Charlie","+62833333333","+62833334444")`
   — response *enriched*: V1 parset `phone`, V2 parse `phones` (2 entry, `is_primary=true`).
3. **Migrate**: `RunAll` backfill 2 — output: "2 legacy records migrated".
   `ReconcileData` → 0 drift.
4. **Read switch**: `SetReadMode(ReadNewOnly)` — V2 baca historical Alice dari
   `user_phones`: `[{Number:+62811111111 IsPrimary:true}]`.
5. **Rollback demo**: kembali ke `WriteLegacyOnly`+`ReadLegacyOnly`, V1 baca Charlie =
   "+62833333333" — tidak ada data loss selama dual-write.
6. **Contract**: `ApplyContract(true)` — legacy read error (410), V2 modern tetap bekerja.

Snapshot akhir: `legacy_reads:3`, `new_reads:2`, `dual_writes:1`, `dual_write_errors:0`,
`backfilled:2`, `drift_detected:0`.

## What the Tests Prove

| Perilaku | Tes | Asersi |
|---|---|---|
| Payload enriched kompatibel V1 + V2 | `TestSerializationBackwardCompatibility` | V1 dapat `phone` exact; V2 dapat 2 `phones` dengan `IsPrimary=true` pada pertama |
| Backfill resumable + idempotent | `TestBackfillIdempotentAndResumable` | batch 3+3+4; run ke-2 = 0 |
| Fallback read + lazy backfill | `TestFallbackRead` | sebelum baca kosong; sesudah 1 entry |
| Drift detect pre/post backfill | `TestDataReconciliationAndDrift` | drift=1 → drift=0 |
| Header deprecasi + contract guard | `TestDeprecationHeadersAndContractEnforcement` | header ada; `ApplyContract(false)` error; `ApplyContract(true)` aman; 410; V2 200 |
| Siklus penuh | `TestFullExpandMigrateContractLifecycle` | Expand dual-write ✓; backfill 2 ✓; drift 0 ✓; read switch ✓; contract legacy error ✓; V2 tetap ✓ |
| Rollback safety | `TestRollbackScenarios` | Skenario A: data aman saat dual-write; Skenario B: `phone=""` setelah rollback pasca-NewOnly (data loss terbukti, tidak diasumsikan) |
| Konkuren under `-race` | `TestConcurrency` | 10 writer × 20, 10 reader × 50, backfill, drift — `DualWriteErrors == 0`, tidak ada race |

Semua asersi memeriksa nilai eksak (count, field, status code), bukan sekadar "no error".
`go test -race ./...` → 8/8 PASS, tidak ada data race.

## Recovery / Rollback

Rollback hanya aman **selama fase dual-write**. Mekanisme dan batasannya:

- **Rollback ke V1 saat `WriteDual` aktif** → data aman. V1 menulis ke `users.phone`,
  V2 menulis ke `user_phones`. Setiap record baru ada di kedua skema. Demo langkah 5 +
  `TestRollbackScenarios` Skenario A membuktinya.
- **Rollback setelah `WriteNewOnly`** → **data loss**. Record yang dibuat di fase ini
  tidak pernah ditulis ke `users.phone`. `TestRollbackScenarios` Skenario B memverifikasi:
  `legacyU2.Phone == ""` (dokumen string kosong, bukan error).
- **Recovery backfill**: worker bersifat *resumable* — `LastProcessedID` checkpoint
  memungkinkan melanjutkan dari titik hentikan setelah crash.
- **Recovery drift**: `ReconcileData` mendeteksi ketidak-konsistenan; tidak ada *self-healing*
  otomatis (operasional harus diintervensi manual).

> **Batasan demonstrasi**: rollback simulasi hanya pada flag runtime, bukan pada versi
> biner yang sebenarnya. Lab tidak mensimulasikan *multi-instance rolling rollback*
> yang melibatkan multiple binary versions concurrently — hanya *mode flag* yang berubah.

## Production Considerations

### Apa yang perlu dipertimbangkan saat produksi

- **Database DDL nyata**: gunakan `ALTER TABLE ADD COLUMN ... DEFAULT <const>` (PostgreSQL
  tidak mere-write baris), `CREATE INDEX CONCURRENTLY` (hindari lock lama), dan
  `ADD CONSTRAINT ... NOT VALID` + `VALIDATE CONSTRAINT` (validasi async).
  Lihat `research/04-database-migration.md`. Lab hanya mensimulasikan secara in-memory.
- **Dual-write cross-storage / cross-service**: tanpa 2PC, gunakan *transactional outbox*
  (satu tabel event + relay idempotent) — lihat `research/03-core-concepts.md` Evidence 7
  (Microservices.io). Lab memakai mutex, bukan outbox.
- **Siklus hapus kolom**: 3-release rule (M, M+1, M+2) adalah artefak ActiveRecord/GitLab,
  bukan aturan universul (research-audit/04-contradictions.md). GitHub memakai 24 bulan
  untuk public API. Jadewan observasi "30 hari" pada lab tidak divalidasi — ilustratif.
- **Monitoring cutover**: counter `LegacyReadHits` perlu *sliding window reset* pada
  produksi; laboratorium memakai counter kumulatif + `force` (gap LOW, documented escape).
- **Scale dual-write**: pada volume tinggi, throttle via rate limiter / outbox; mutex
  dalam memori tidak scale ke banyak instance.

### Apa yang TIDAK ditunjukkan

- Koordinasi transaksi terdistribusi (2PC) atau event streaming (Kafka/CDC).
- DDL runtime PostgreSQL (`SET lock_timeout`, `CREATE INDEX CONCURRENTLY`).
- Multi-release database refactorings berlapis (GitLab 3-release rule).
- Cross-microservice dual-write tanpa outbox — didokumentasikan sebagai *open research*
  (research/04-database-migration.md, research/10-open-questions.md).

## Common Mistakes

1. **Melewatkan fase Contract** → technical debt permanen, dua skema hidup untuk selalu.
   *Mitigasi*: buat *ticket* kontrak per epic; *alert* otomatis pada `LegacyReadHits == 0`
   selama periode observasi.

2. **Beralih baca terlalu awal** → data starvation (404/kosong) untuk record belum backfill.
   *Mitigasi*: pakai `ReadFallback` selama migrasi, atau pastikan backfill 100% selesai
   sebelum `ReadNewOnly`.

3. **Asumsikan idempotent tanpa cek eksplisit** → duplikat pada rerun backfill.
   *Mitigasi*: periksa eksistensi sebelum insert, atau pakai `UPSERT`/`ON CONFLICT DO NOTHING`.

4. **Abaiakan overhead dual-write pada skala** → kontensi DB, latency tulis naik.
   *Mitigasi*: pola outbox (satu tulisan lokal, proses async ke tujuan kedua).

## Case Study

Lab ini mementr. **Case Study A: Customer Phone 1:1 → 1:N** (research/09-case-studies.md):

- **Baseline**: `users.phone TEXT` satu nomor per pengguna (Alice `+62811111111`).
- **Target**: banyak nomor per pengguna, dengan bendera `is_primary`.

Transisi:
1. `schema.sql` V2: buat `user_phones(user_id, number, is_primary)` — kolom `phone` tidak disentuh.
2. Mode `WriteDual`: tiap `CreateUser` tulis ke kedua tabel.
3. Backfill historis: Alice, Bob — migrasi ke `user_phones`.
4. `ReadNewOnly`: V2 baca `user_phones`; `ReadFallback` melindungi sebelum selesai.
5. `ApplyContract`: drop `phone` setelah traffic V1 = 0.

> research/09-case-studies.md berisi tambahan **Case Study B** (CMMS Invoice 1:1→N:M)
> dan **Case Study C** (Multi-Currency schema split). Keduanya tidak diimplementasikan
> dalam kode lab — mereferensikan contoh konseptual, bukan perilaku lab yang diverifikasikan.

## Checklist

- [ ] Fase Expand: tabel/field baru ditambahkan, lama **tidak disentuh**
- [ ] Payload additive: klien lama parset field lama; field baru di-ignore
- [ ] Dual-write aktif: tiap record baru ada di kedua skema
- [ ] Backfill selesai: `ReconcileData` → drift = 0
- [ ] Fallback read siap: record un-backfill tetap readable
- [ ] Feature flag siap: `WriteMode`/`ReadMode` dapat di-flip runtime
- [ ] Observability online: `LegacyReadHits` dilacak
- [ ] Traffic lama = 0 selama jendela observasi
- [ ] Contract: `ApplyContract` sukses; V1 dapat `410 Gone`; V2 tetap `200`
- [ ] Rollback simulasi: V1 masih baca data V2 — tidak ada data loss
- [ ] Audit lapis bawaan: `go test -race ./...` — 0 data race, 0 error

## Key Takeaways

1. Perubahan struktural tidak harus berhenti jasa — pecahkan menjadi Expand/Migrate/Contract.
2. Dual-write adalah jaminan atomicity sempat-pinjam; pada memori pakai mutex, pada prod pakai outbox/transaksi.
3. Backfill harus *idempotent + resumable* (checkpoint + dedup) — bukan batch besar satu kali.
4. `ReadFallback` (dual-read) adalah *parachute* sebelum backfill selesai — wajib ada.
5. Contract (hapus kolom/field) hanya saat traffic legacy = 0 — dikunci oleh metrik, bukan asumsi.
6. Rollback hanya aman di fase dual-write — setelah `WriteNewOnly`, rollback = data loss.
7. Feature flags memisahkan *deploy* dari *activate* — enabler canary dan rollback instan.
8. Header `Deprecation`/`Sunset` + metrik adalah antarmuka manusia ke putus-kepaksaan teknis.
9. Penyederhanaan lab (in-memory, mutex) tidak mengaburkan mekanika relasional — tapi
   jalur upgrade ke `database/sql`/outbox jelas didokumentasikan.
10. "30 hari" jendela observasi adalah heuristik terbuka — tidak berlaku sebagai aturan universal.

## Sources

- **Martin Fowler, *Parallel Change* (2014-05-13)**: konsep Expand/Migrate/Contract, Grid contoh, feature flag, blue-green. `research/06-expand-migrate-contract.md`, `research/03-core-concepts.md` (Evidence 2,4,9,12,14).
- **PostgreSQL Global Development Group (PG 18)**: additive `ALTER TABLE`, `CREATE INDEX CONCURRENTLY`, `NOT VALID`/`VALIDATE CONSTRAINT`. `research/04-database-migration.md`, `research/03-core-concepts.md` (Evidence 5,13).
- **Google AIP-180 (2019-07-23)**: definisi additive/non-breaking, larangan type change, source/wire/semantic compatibility. `research/03-core-concepts.md`, `research/05-api-compatibility.md`.
- **GitHub REST API (2026-03-10)**: 24-month support window, Deprecation/Sunset header, breaking/non-breaking daftar. `research/05-api-compatibility.md`, `research/03-core-concepts.md` (Evidence 10,11,15).
- **Stripe API Versioning (2017-08-05)**: request-time version transformation pipeline, pin-by-first-request. `research/05-api-compatibility.md`.
- **Microservices.io / Chris Richardson, *Transactional Outbox***: dual-write 2PC kebutuhatan, relay idempotency. `research/03-core-concepts.md` (Evidence 7).
- **Confluent Schema Registry**: definisi BACKWARD / FORWARD / compatibility. `research/03-core-concepts.md` (Evidence 2).

Audit status: Research APPROVED (research-audit/07-verdict.md); Engineering APPROVED + APPROVED_WITH_WARNINGS (engineering-audit/06-verdict.md, engineering-audit-opensource/06-verdict.md — 7 temuan LOW).

Catatan akurasi: 7 temuan LOW pada engineering-audit-opensource tidak menggangu kebenaran inti. Lihat `content/05-key-takeaways.md` poin 9 dan `06-source-map.md` Failure Modes.
