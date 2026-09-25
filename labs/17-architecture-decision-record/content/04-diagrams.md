# Diagrams

Diagram-diagram berikut merepresentasikan alur kerja, siklus hidup, dan arsitektur validasi yang diimplementasikan pada lab `labs/17-architecture-decision-record`.

---

## 1. Siklus Hidup Status ADR (ADR Lifecycle)

Diagram ini menggambarkan transisi status keputusan yang didukung oleh model sistem (`internal/adr/models.go`):

```text
               ┌──────────────┐
               │   Proposed   │
               └──────┬───────┘
                      │
         ┌────────────┴────────────┐
         ▼                         ▼
  ┌──────────────┐          ┌──────────────┐
  │   Accepted   │          │   Rejected   │
  └──────┬───────┘          └──────────────┘
         │
    ┌────┴────────────────────────┐
    ▼                             ▼
┌──────────────┐           ┌──────────────┐
│  Superseded  │           │  Deprecated  │
└──────────────┘           └──────────────┘
```

Keterangan:
- **Proposed:** Tahap usulan dan peninjauan awal.
- **Accepted:** Usulan disetujui; menjadi pedoman arsitektur aktif (tidak dapat diubah secara retroaktif).
- **Rejected:** Usulan ditolak setelah evaluasi; disimpan sebagai catatan riwayat untuk mencegah perdebatan berulang.
- **Superseded:** Keputusan lama digantikan oleh ADR baru (memiliki tautan timbal balik ke ADR pengganti).
- **Deprecated:** Keputusan dihentikan tanpa pengganti langsung.

---

## 2. Hubungan Silsilah Penggantian Timbal-Balik (Bidirectional Supersession Lineage)

Diagram ini mengilustrasikan hubungan graf keputusan terarah antara ADR awal dan ADR pengganti sebagaimana diverifikasi dalam percontohan (`cmd/demo/main.go`):

```text
┌──────────────────────────────────────┐
│ ADR 0001                             │
│ Title: Use Modular Monolith...       │
│ Status: Superseded by 2 ─────────────┼────────┐
└──────────────────────────────────────┘        │
                   ▲                            │
                   │ (Tautan Verifikasi Timbal-Balik)
                   │                            │
┌──────────────────┴───────────────────┐        │
│ ADR 0002                             │        │
│ Title: Extract Notification Service..│        │
│ Status: Accepted                     │        │
│ Supersedes: 1 ◄──────────────────────┼────────┘
└──────────────────────────────────────┘
```

Keterangan:
- Linter memverifikasi bahwa penanda `Superseded by 2` pada ADR 1 merujuk ke ADR 2 yang benar-benar ada.
- Linter juga memverifikasi bahwa ADR 2 mendeklarasikan secara eksplisit klausa `Supersedes: 1`. Jika salah satu sisi hilang, linter menolak rangkaian keputusan tersebut.

---

## 3. Alur Kerja Pemeriksaan Integritas Linter (Validation Pipeline)

Diagram ini menunjukkan alur pemrosesan dari dokumen Markdown hingga status integritas akhir:

```text
[ Markdown Documents ]
        │
        ▼
[ adr.Parse() ] ──────────────► Ekstraksi ID, Title, Status, SupersededBy, Supersedes
        │
        ▼
[ []*Record ]
        │
        ▼
[ adr.Linter.Validate() ]
        │
        ├─► [ 1. Pemeriksaan Monotonik ]
        │     - Urutkan ID: ids[i] == i + 1
        │     - Deteksi duplikasi atau nomor lompat
        │
        └─► [ 2. Validasi Graf Konkuren (Goroutines) ]
              - Cek rec.Status == Superseded -> target exists & supersedes rec.ID
              - Cek rec.Supersedes != 0 -> target exists & target is superseded by rec.ID
              - Mutex-guarded aggregation of errors
        │
        ▼
[ Status Hasil: Lulus / Kumpulan Galat Integritas ]
```
