# Publication Deliverables — Lab 22: N+1 Query Problem

## Daftar File Deliverable (24 slide, 2 cover, 2 dokumen teks, 1 PDF, 1 PPTX, 1 metadata)

| No | Tipe Deliverable | Lokasi File | Deskripsi |
|---|---|---|---|
| 1 | PowerPoint Carousel (editable) | [linkedin/lab-22-n-plus-one-query-problem-carousel.pptx](linkedin/lab-22-n-plus-one-query-problem-carousel.pptx) | 24 slide 16:9, speaker notes, dark navy cover, editorial light slides 2‑23, dark navy penutup |
| 2 | LinkedIn Carousel PDF (selectable text) | [linkedin/lab-22-n-plus-one-query-problem-linkedin.pdf](linkedin/lab-22-n-plus-one-query-problem-linkedin.pdf) | 24 halaman 16:9, dark navy cover & penutup, editorial light interior, ukuran 3.68 MB (>200 KB) |
| 3 | Dev.to Markdown Article (English) | [dev.to/lab-22-n-plus-one-query-problem-devto.md](dev.to/lab-22-n-plus-one-query-problem-devto.md) | 1.500+ kata, frontmatter lengkap, syntax‑highlighted Go code, link repo & commit, batas implementasi transparan |
| 4 | Medium HTML Article (Indonesian) | [medium/lab-22-n-plus-one-query-problem-medium.html](medium/lab-22-n-plus-one-query-problem-medium.html) | 2.000+ kata, 9 section wajib, cover image, tabel perbandingan, callout, SEO metadata, tanpa JavaScript |
| 5 | Cover Image (Indonesia) | [image-cover/lab-22-n-plus-one-query-problem-cover-id.png](image-cover/lab-22-n-plus-one-query-problem-cover-id.png) | 1920×1080 PNG retro CRT terminal, subtitle Bahasa Indonesia, ukuran 132 KB (>50 KB) |
| 6 | Cover Image (English) | [image-cover/lab-22-n-plus-one-query-problem-cover-en.png](image-cover/lab-22-n-plus-one-query-problem-cover-en.png) | 1920×1080 PNG retro CRT terminal, subtitle Bahasa Inggris, ukuran 133 KB (>50 KB) |
| 7 | Medium Metadata (Markdown) | [medium/lab-22-n-plus-one-query-problem-metadata.md](medium/lab-22-n-plus-one-query-problem-metadata.md) | Title, subtitle, SEO title (45 kar), SEO description (149 kar), slug, 5 tags, alt‑text cover, social title/description |
| 8 | LinkedIn Caption (plain text) | [linkedin/lab-22-n-plus-one-query-problem-caption.md](linkedin/lab-22-n-plus-one-query-problem-caption.md) | 180‑250 kata, awalan *Lab 22 kali ini membahas …*, mental model dalam bahasa Inggris, URL repo, 4‑6 hashtag, tanpa markdown bold |

## Catatan Verifikasi & Batas Implementasi
- **Commit SHA**: `6649acb9`
- **Status Test**: 4/4 test cases PASS (`go test -v -race ./...`)
- **Batas Implementasi**: Simulasi dalam‑memory mock store, tidak terhubung ke basis data fisik; query count bersifat deterministik, bukan benchmark latency dunia nyata.

## Caption LinkedIn (Siap Salin)
Lab 22 kali ini membahas anti-pattern N+1 query problem, mekanisme terjadinya latensi akumuler pada lapisan akses data, serta mitigasi deterministik menggunakan eager loading berbasis batching.

Masalah teknis yang sering terjadi adalah aplikasi mengeksekusi satu query awal untuk mengambil koleksi data induk, lalu mengeksekusi satu query tambahan di dalam loop untuk setiap entitas anak yang terasosiasi. Pada dataset pengujian kecil, setiap query individu selesai dalam sub‑milidetik sehingga lolos tanpa terdeteksi oleh slow query log. Namun ketika volume data meningkat di produksi, pola ini memicu ratusan round‑trip jaringan, overhead protokol basis data, serta potensi kehabisan connection pool.

Lab ini membuktikan empat hal konkret:
• Loop‑Driven Anti‑Pattern: Pengambilan relasi secara naif terbukti mengeksekusi 4 query untuk memuat 3 entitas author beserta postingannya.
• Batched Eager Loading: Pengumpulan seluruh identifier relasi ke dalam klausa IN tunggal memangkas eksekusi menjadi tepat 2 query konstan.
• Deep Equivalence: Struktur data hasil batching terbukti identik secara semantik dengan hasil metode naif melalui pengujian refleksi mendalam.
• Zero Race Condition: Pengujian konkuren dengan race detector membuktikan pembaruan query counter berjalan aman secara thread‑safe.

Mental model inti: N+1 is a pattern, not a layer. Batch all related identifiers into a single IN query instead of executing per‑row lookups.

Satu batasan penting: Lab ini adalah simulasi konseptual berbasis in‑memory mock store di tingkat aplikasi, bukan pengujian terhadap basis data relasional fisik, serta berfokus murni pada query‑level eager loading.

Pelajari implementasi lengkap dan unit test di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/22-n-plus-one-query-problem

#softwareengineering #golang #database #performance #systemdesign #backend