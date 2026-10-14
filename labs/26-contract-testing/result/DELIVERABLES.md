# Publication Deliverables — Lab 26: Contract Testing

## Daftar File Deliverable (24 slide, 2 cover, 2 dokumen teks, 1 PDF, 1 PPTX, 1 metadata)

| No | Tipe Deliverable | Lokasi File | Deskripsi |
|---|---|---|---|
| 1 | PowerPoint Carousel (editable) | [linkedin/lab-26-contract-testing-carousel.pptx](linkedin/lab-26-contract-testing-carousel.pptx) | 24 slide 16:9, speaker notes, dark navy cover, editorial light slides 2‑23, dark navy penutup (118 KB) |
| 2 | LinkedIn Carousel PDF (selectable text) | [linkedin/lab-26-contract-testing-linkedin.pdf](linkedin/lab-26-contract-testing-linkedin.pdf) | 24 halaman 16:9, dark navy cover & penutup, editorial light interior, ukuran 796 KB (>200 KB) |
| 3 | Dev.to Markdown Article (English) | [dev.to/lab-26-contract-testing-devto.md](dev.to/lab-26-contract-testing-devto.md) | 2.241 kata (>1.500 kata), frontmatter lengkap, syntax‑highlighted Go code, link repo & commit, batas implementasi transparan |
| 4 | Medium HTML Article (Indonesian) | [medium/lab-26-contract-testing-medium.html](medium/lab-26-contract-testing-medium.html) | 2.717 kata (>2.000 kata), 9 section wajib, cover image, tabel perbandingan, callout, SEO metadata, tanpa JavaScript |
| 5 | Cover Image (Indonesia) | [image-cover/lab-26-contract-testing-cover-id.png](image-cover/lab-26-contract-testing-cover-id.png) | 1920×1080 PNG retro CRT terminal, subtitle Bahasa Indonesia, ukuran 743 KB (>50 KB) |
| 6 | Cover Image (English) | [image-cover/lab-26-contract-testing-cover-en.png](image-cover/lab-26-contract-testing-cover-en.png) | 1920×1080 PNG retro CRT terminal, subtitle Bahasa Inggris, ukuran 740 KB (>50 KB) |
| 7 | Medium Metadata (Markdown) | [medium/lab-26-contract-testing-metadata.md](medium/lab-26-contract-testing-metadata.md) | Title, subtitle, SEO title (43 kar), SEO description (151 kar), slug, 5 tags, alt‑text cover, social title/description |
| 8 | LinkedIn Caption (plain text) | [linkedin/lab-26-contract-testing-caption.md](linkedin/lab-26-contract-testing-caption.md) | 241 kata (180‑250 kata), awalan *Lab 26 kali ini membahas …*, mental model dalam bahasa Inggris, URL repo, 6 hashtag, tanpa markdown bold |

## Catatan Verifikasi & Batas Implementasi
- **Commit SHA**: `e963675c67e7d518e19a07ad930dbdb5be73fff1`
- **Status Test**: 7/7 test cases PASS (`go test -v -race ./...`) dalam 0.390s (test biasa) dan 1.425s (race detector).
- **Batas Implementasi**: Engine verifikasi CDC murni Go standard library tanpa binary daemon Pact; kontrak dipertukarkan in-process tanpa broker remote; state provider disimulasikan melalui deterministic routing path.

## Caption LinkedIn (Siap Salin)
Lab 26 kali ini membahas Consumer-Driven Contract testing, deteksi breaking schema changes otomatis, serta implementasi CI/CD deployment gate di Go.

Masalah teknis yang sering terjadi pada microservices adalah perubahan skema API tanpa koordinasi: penggantian casing enum status, penggantian nama field kunci JSON, atau mutasi tipe data numerik menjadi teks. Pada pengujian unit lokal dengan mock statis, kesalahan ini tidak terdeteksi karena mock cepat usang. Sementara pengujian end-to-end di staging memakan waktu lama dan sering terganggu flakiness jaringan, sehingga breaking change kerap lolos hingga merusak klien seluler di lingkungan produksi.

Lab ini membuktikan empat hal konkret:
• Minimal Subset Contract: Klien seluler mendefinisikan ekspektasi payload minimal yang hanya memuat field yang benar-benar dikonsumsi.
• Build-Time Verification: Engine verifikasi mengeksekusi interaksi HTTP terhadap server ephemeral dan memvalidasi struktur secara rekursif dalam waktu kurang dari satu detik.
• Breaking Change Barrier: Tiga mutasi payload terbukti langsung memicu kegagalan verifikasi dengan pesan kesalahan field yang spesifik, memblokir deployment otomatis di tingkat CI.
• Safe API Evolution: Pola dual DTO paralel mempertahankan kepatuhan kontrak V1 sekaligus mengekspos skema V2 untuk migrasi bertahap tanpa downtime.

Mental model inti: Mocks drift into false confidence; contract tests turn consumer assumptions into executable provider guardrails.

Satu batasan penting: Lab ini menggunakan engine verifikasi Go murni berbasis pustaka standar tanpa binary eksternal Pact, serta menyederhanakan manajemen provider state melalui routing path deterministik.

Pelajari implementasi lengkap dan unit test di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/26-contract-testing

#softwareengineering #golang #microservices #testing #devops #systemdesign
