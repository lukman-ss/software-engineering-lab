# Rekapitulasi Deliverables — Lab 27: Database Constraints & Data Integrity

Seluruh berkas keluaran untuk **Software Engineering Lab 27: Database Constraints & Data Integrity** telah selesai dibuat dan diverifikasi sesuai standar mutu publikasi Lab 14 & 15.

---

## 1. Tabel Berkas Publikasi

| No | Deliverable | Jalur Berkas | Format / Ukuran | Deskripsi |
|---|---|---|---|---|
| 0 | Slide & Metadata Definition | [`content_def.py`](content_def.py) | Python (7.9 KB) | Sumber kode reproducible untuk 24 slide, palet warna, dan metadata |
| 1 | Artikel Dev.to | [`dev.to/lab-27-database-constraints-devto.md`](dev.to/lab-27-database-constraints-devto.md) | Markdown (2,457 kata) | Artikel teknis mendalam Bahasa Inggris lengkap dengan frontmatter |
| 2 | Artikel Medium | [`medium/lab-27-database-constraints-medium.html`](medium/lab-27-database-constraints-medium.html) | HTML Standalone (2,807 kata) | Artikel komprehensif Bahasa Indonesia mencakup 9 section wajib |
| 3 | Metadata Medium | [`medium/lab-27-database-constraints-metadata.md`](medium/lab-27-database-constraints-metadata.md) | Markdown (1.3 KB) | Konfigurasi SEO, judul alternatif, slug, dan tepat 5 topik tag |
| 4 | Cover Visual ID | [`image-cover/lab-27-database-constraints-cover-id.png`](image-cover/lab-27-database-constraints-cover-id.png) | PNG 1920×1080 (728 KB) | Visual retro monitor CRT bezel dengan tipografi Bahasa Indonesia |
| 5 | Cover Visual EN | [`image-cover/lab-27-database-constraints-cover-en.png`](image-cover/lab-27-database-constraints-cover-en.png) | PNG 1920×1080 (731 KB) | Visual retro monitor CRT bezel dengan tipografi Bahasa Inggris |
| 6 | Slide Deck LinkedIn PDF | [`linkedin/lab-27-database-constraints-linkedin.pdf`](linkedin/lab-27-database-constraints-linkedin.pdf) | PDF 16:9 (473 KB) | Tepat 24 halaman landscape resolusi tinggi, selectable text |
| 7 | Slide Carousel PPTX | [`linkedin/lab-27-database-constraints-carousel.pptx`](linkedin/lab-27-database-constraints-carousel.pptx) | PPTX 16:9 (108 KB) | Tepat 24 slide presentasi lengkap dengan Speaker Notes di tiap slide |
| 8 | Caption LinkedIn | [`linkedin/lab-27-database-constraints-caption.md`](linkedin/lab-27-database-constraints-caption.md) | Markdown (217 kata) | Naskah posting LinkedIn siap salin tanpa bold markdown |

---

## 2. Catatan Teknis & Verifikasi Mutu

- **Commit SHA**: `bdb285d35a71c9a309e83a0e4497b9dbd842ceb8`
- **Hasil Pengujian**: 8/8 unit dan integration test berstatus **PASS** (`go test -v -race ./...`).
- **Race Detector**: Bersih tanpa data race terdeteksi (`PASS`).
- **Demo Eksekusi**: Program `cmd/demo/main.go` berhasil menjalankan comparasi antara `UnsafeStore` dan `SafeStore` pada 50 goroutine secara deterministik.
- **Batas Implementasi Utama**:
  1. *In-Memory Engine Simulator*: Engine menggunakan Go map dan `sync.RWMutex` untuk memodelkan semantik PostgreSQL tanpa Docker. Bukan PostgreSQL live — tidak ada B-tree page latch atau MVCC.
  2. *Coarse-Grained Locking*: Penguncian onggok pada level tabel, bukan baris (`row-level lock`). Cocok untuk pembelajaran konsep, tidak untuk benchmarking produksi.
  3. *Tanpa Query Counter*: Simulator tidak memiliki query execution plan atau cost estimator PostgreSQL.

---

## 3. Teks Siap Salin: Caption LinkedIn

```text
Lab 27 kali ini membahas pencegahan race condition dan kerusakan integritas data melalui pemindahan validasi ke level database constraints di Go.

Masalah teknis yang sering terjadi adalah pola check-then-act di application layer (SELECT lalu INSERT). Di bawah konkurensi, dua goroutine dapat membaca bahwa email belum ada secara bersamaan, sehingga keduanya mengeksekusi INSERT dan menghasilkan record duplikat yang merusak database.

Lab ini membuktikan lima pilar integritas data konkret:
• NOT NULL Enforcement: Kolom wajib seperti email, username, dan user_id ditolak seketika pada level engine dengan SQLSTATE 23502.
• CHECK Boolean Predicate: Validasi batas nilai (age >= 18, enum status, total_cents > 0) dievaluasi atomik dengan SQLSTATE 23514.
• UNIQUE B-Tree Protection: Penulisan atomik di bawah lock menolak duplikasi konkuren, menghasilkan tepat 1 pemenang dan N-1 error SQLSTATE 23505.
• Referential Integrity: FOREIGN KEY mencegah orphan rows pada tabel orders dengan SQLSTATE 23503.
• Partial Unique Index: Predikat WHERE deleted_at IS NULL mengizinkan registrasi ulang email setelah soft delete tanpa mengorbankan keunikan aktif.

Mental model inti: Application validation is for user experience; database constraints are for absolute correctness. Never rely on application checks alone.

Satu batasan penting: Lab ini menggunakan simulator in-memory engine Go dengan sync.RWMutex untuk membuktikan konsep relasional tanpa dependensi Docker live.

Pelajari implementasi lengkap dan seluruh pembuktian test di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/27-database-constraints

#softwareengineering #golang #databases #systemdesign #dataintegrity #postgresql
```