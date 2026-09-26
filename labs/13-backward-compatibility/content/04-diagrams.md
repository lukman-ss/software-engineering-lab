# Backward Compatibility: Expand → Migrate → Contract Pattern

## Arsitektur Komponen Utama

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

*Sumber: engineering/01-design.md*

Keterangan:
- Kedua jenis klien (V1 legacy dan V2 modern) mengakses HTTP/Service Layer yang sama.
- Service Layer berisi pipeline transformasi, feature flags untuk kontrol mode baca/tulis, dan pengumpul metrik observasi.
- Backfill Worker menyalin data dari Legacy Storage ke Modern Storage dalam batch.
- Storage dibagi menjadi dua: legacy (`users.phone` kolom tunggal) dan modern (`user_phones` tabel terpisah dengan kolom `user_id`, `number`, `is_primary`).

## Urutan Fase Migrasi

```text
Fase 0: Baseline
┌─────────────────┐
│ users(phone)    │ ◄── V1 Client baca/tulis
│ WriteLegacyOnly │
│ ReadLegacyOnly  │
└─────────────────┘

Fase 1: Expand
┌─────────────────┐      ┌──────────────────┐
│ users(phone)    │ ◄────┤ user_phones      │ ← tabel baru, kosong
│ WriteDual       │      │ (dibuat)         │
└─────────────────┘      └──────────────────┘
       │ V1 dan V2 keduanya bisa baca (enriched payload)

Fase 2: Migrate - Dual Write
┌─────────────────┐      ┌──────────────────┐
│ users(phone)    │ ◄──► │ user_phones      │ ◄── dual-write atomic
│ (tulis+ baca)   │      │ (tulis+ baca)    │
└─────────────────┘      └──────────────────┘

Fase 3: Migrate - Backfill
┌─────────────────┐      ┌──────────────────┐
│ users(phone)    │ ───► │ user_phones      │ ◄── backfill batch
│ (historis)      │ copy │ (terisi penuh)   │
└─────────────────┘      └──────────────────┘

Fase 4: Switch Read Path
┌─────────────────┐      ┌──────────────────┐
│ users(phone)    │      │ user_phones      │ ◄── V2 Client baca
│ (tidak dibaca)  │      │ ReadNewOnly      │
└─────────────────┘      └──────────────────┘

Fase 5: Contract
┌─────────────────┐      ┌──────────────────┐
│ users           │      │ user_phones      │ ◄── V2 Client baca/tulis
│ (kolom di-drop) │      │ WriteNewOnly     │
│                 │      │ ReadNewOnly      │
└─────────────────┘      └──────────────────┘
       ▲ V1 Client dapatkan 410 Gone
```

*Sumber: engineering/01-design.md, research/06-expand-migrate-contract.md*

## Skema Database: Evolusi Langkah-demi-Langkah

```text
V1 Baseline (1:1 relationship):
┌──────────────┐
│ users        │
├──────────────┤
│ id           │
│ name         │
│ phone (TEXT) │ ← kolom tunggal
└──────────────┘

V2 Expand (1:N ditambahkan tanpa menyentuh kolom lama):
┌──────────────┐      ┌─────────────────┐
│ users        │      │ user_phones     │
├──────────────┤      ├─────────────────┤
│ id           │◄─────┤ user_id (FK)    │
│ name         │      │ number          │
│ phone (TEXT) │      │ is_primary      │
└──────────────┘      └─────────────────┘
   (tetap)              (tabel baru)

V3 Contract (kolom lama dihapus SETELAH trafic legacy = 0):
┌──────────────┐      ┌─────────────────┐
│ users        │      │ user_phones     │
├──────────────┤      ├─────────────────┤
│ id           │◄─────┤ user_id (FK)    │
│ name         │      │ number          │
│ (kolom drop) │      │ is_primary      │
└──────────────┘      └─────────────────┘
```

*Sumber: schema.sql*

## Diagram Alir Fallback Read (ReadFallback)

```text
GET /user/{id} dengan ReadMode = ReadFallback
                    │
                    ▼
        ┌─────────────────────┐
        │ Baca user_phones    │
        │ untuk user_id       │
        └─────────┬───────────┘
                  │
        ┌─────────┴──────────┐
        │                    │
   Ditemukan?           Kosong?
        │                    │
        ▼                    ▼
┌──────────────┐   ┌──────────────────┐
│ Kembalikan   │   │ Baca users.phone │
│ dari modern  │   │ (legacy column)  │
└──────────────┘   └────────┬─────────┘
                            │
                   ┌────────┴────────┐
                   │                 │
              Ditemukan?        Kosong?
                   │                 │
                   ▼                 ▼
         ┌────────────────┐  ┌──────────────┐
         │ Kembalikan &   │  │ Kembalikan   │
         │ lazy backfill  │  │ kosong       │
         │ ke user_phones │  │              │
         └────────────────┘  └──────────────┘
```

*Sumber: internal/compat/service.go (fungsi GetUser, mode ReadFallback)*

## Diagram Alir Keputusan Deprecation/Contract

```text
Endpoint /v1/users menerima request
                │
                ▼
   ┌────────────────────────┐
   │ ContractApplied?       │
   └──────────┬─────────────┘
              │
     ┌────────┴────────┐
     │                 │
    Ya                Tidak
     │                 │
     ▼                 ▼
┌────────────┐  ┌──────────────────┐
│ 410 Gone   │  │ Kembalikan data  │
│ Legacy     │  │ + headers:       │
│ dihapus    │  │ Deprecation:true │
│            │  │ Sunset:<date>    │
└────────────┘  └──────────────────┘
```

*Sumber: internal/compat/handler.go (fungsi GetUserV1)*

## Diagram Status WriteMode dan ReadMode

```text
WriteMode:
  WriteLegacyOnly ──► WriteDual ──► WriteNewOnly
       │                  │               │
       │                  │               │
       ▼                  ▼               ▼
  Hanya ke          Ke kedua         Hanya ke
  users.phone       tabel            user_phones

ReadMode:
  ReadLegacyOnly ──► ReadFallback ──► ReadNewOnly
       │                  │               │
       │                  │               │
       ▼                  ▼               ▼
  Hanya dari        Modern dulu,    Hanya dari
  users.phone       fallback ke    user_phones
                    legacy
```

*Sumber: internal/compat/flags.go*

Catatan: Kombinasi yang digunakan dalam lab:
- Baseline: WriteLegacyOnly + ReadLegacyOnly
- Expand/Migrate: WriteDual + ReadFallback (atau ReadLegacyOnly di awal)
- Read Switch: WriteDual + ReadNewOnly
- Contract: WriteNewOnly + ReadNewOnly + ContractApplied=true
