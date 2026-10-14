# Publication Deliverables — Lab 39: Bloom Filters

## Daftar File Deliverable (24 slide, 2 cover, 2 dokumen teks, 1 PDF, 1 PPTX, 1 metadata)

| No | Tipe Deliverable | Lokasi File | Deskripsi |
|---|---|---|---|
| 1 | PowerPoint Carousel (editable) | [linkedin/lab-39-bloom-filters-carousel.pptx](linkedin/lab-39-bloom-filters-carousel.pptx) | 24 slide 16:9, speaker notes lengkap, dark navy cover, editorial light slides 2‑23, dark navy penutup |
| 2 | LinkedIn Carousel PDF (selectable text) | [linkedin/lab-39-bloom-filters-linkedin.pdf](linkedin/lab-39-bloom-filters-linkedin.pdf) | 24 halaman 16:9, dark navy cover & penutup, editorial light interior, ukuran 954 KB (>200 KB) |
| 3 | Dev.to Markdown Article (English) | [dev.to/lab-39-bloom-filters-devto.md](dev.to/lab-39-bloom-filters-devto.md) | 2.408 kata (>1.500 kata), frontmatter lengkap, syntax‑highlighted Go code, link repo & commit, batas implementasi transparan |
| 4 | Medium HTML Article (Indonesian) | [medium/lab-39-bloom-filters-medium.html](medium/lab-39-bloom-filters-medium.html) | 2.101 kata (>2.000 kata), 22 section editorial komprehensif, cover image, tabel perbandingan, callout, tanpa JavaScript |
| 5 | Cover Image (Indonesia) | [image-cover/lab-39-bloom-filters-cover-id.png](image-cover/lab-39-bloom-filters-cover-id.png) | 1920×1080 PNG retro CRT terminal, subtitle Bahasa Indonesia, ukuran 742 KB (>50 KB) |
| 6 | Cover Image (English) | [image-cover/lab-39-bloom-filters-cover-en.png](image-cover/lab-39-bloom-filters-cover-en.png) | 1920×1080 PNG retro CRT terminal, subtitle Bahasa Inggris, ukuran 741 KB (>50 KB) |
| 7 | Medium Metadata (Markdown) | [medium/lab-39-bloom-filters-metadata.md](medium/lab-39-bloom-filters-metadata.md) | Title, subtitle, SEO title (54 kar), SEO description (153 kar), slug, 5 tags, alt‑text cover, social title/description |
| 8 | LinkedIn Caption (plain text) | [linkedin/lab-39-bloom-filters-caption.md](linkedin/lab-39-bloom-filters-caption.md) | 230 kata (rentang 180‑250 kata), awalan *Lab 39 kali ini membahas …*, mental model Bahasa Inggris, URL repo, 6 hashtag, tanpa markdown bold |

## Catatan Verifikasi & Batas Implementasi
- **Commit SHA**: `c977525a4a31ba811006df2bbc04fd031ee11765`
- **Status Test**: 7/7 test cases PASS (`go test -v -race ./...`)
- **Batas Implementasi**: Menggunakan FNV-1a dengan rotasi double hashing Kirsch-Mitzenmacher; tidak mendukung operasi deletion secara fundamental (immutable segment architecture); estimasi kapasitas n harus ditentukan di muka untuk mencegah saturasi kerapatan bit array.

## Caption LinkedIn (Siap Salin)
Lab 39 kali ini membahas struktur data probabilistik Bloom filter, mekanisme eliminasi operasi I/O disk yang berlebihan pada sistem penyimpanan berbasis LSM-Tree, serta proteksi backend terhadap serangan cache penetration.

Masalah teknis yang kerap melumpuhkan arsitektur data adalah biaya pencarian kunci yang sebenarnya tidak ada. Pada basis data LSM-Tree, kueri data non-existent memaksa sistem membaca seluruh file SSTable segment di disk secara berurutan. Di lapisan caching, banjir permintaan untuk identifier yang tidak valid melompati memori cache dan langsung menghantam basis data hilir, memicu kelelahan connection pool dan lonjakan latensi p99.

Lab ini membuktikan empat hal konkret:
• Zero False Negatives: Pengujian 10.000 kunci membuktikan seluruh elemen yang terdaftar selalu dikembalikan true tanpa kegagalan.
• Presisi False Positive: Pengujian 50.000 kueri kunci absen menghasilkan tingkat error empiris 0.996 persen, konsisten terhadap target teoretis 1.00 persen.
• Reduksi 99 Persen I/O Disk: Simulasi 10.000 kueri pada 10 segment memangkas pembacaan disk dari 100.000 kali menjadi hanya 988 kali.
• Perlindungan Cache Mutlak: Pemasangan admission gate menahan 100 persen kueri penetrasi berbahaya sehingga 0 panggilan tembus ke database.

Mental model inti: Zero false negatives is a feature, not a bug. Trade certainty for compression and prune before you peek.

Satu batasan penting: Bloom filter standar tidak mendukung operasi penghapusan dan membutuhkan estimasi kapasitas di muka untuk mencegah saturasi bit array.

Pelajari implementasi lengkap dan unit test di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/39-bloom-filters

#softwareengineering #golang #database #performance #systemdesign #backend
