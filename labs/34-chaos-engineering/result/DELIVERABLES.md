# Publication Deliverables — Lab 34: Chaos Engineering & Fault Injection

## Daftar File Deliverable (24 slide, 2 cover, 2 dokumen teks, 1 PDF, 1 PPTX, 1 metadata, 1 caption)

| No | Tipe Deliverable | Lokasi File | Deskripsi |
|---|---|---|---|
| 0 | Python Content Definition | [content_def.py](content_def.py) | Definisi 24 slide, metadata, warna, dan mental models (reproducible builder source) |
| 1 | PowerPoint Carousel (editable) | [linkedin/lab-34-chaos-engineering-carousel.pptx](linkedin/lab-34-chaos-engineering-carousel.pptx) | Tepat 24 slide 16:9, speaker notes lengkap, dark navy cover, editorial light slides 2‑23, dark navy penutup, ukuran 120 KB (>100 KB) |
| 2 | LinkedIn Carousel PDF (selectable text) | [linkedin/lab-34-chaos-engineering-linkedin.pdf](linkedin/lab-34-chaos-engineering-linkedin.pdf) | Tepat 24 halaman 16:9, dark navy cover & penutup, editorial light interior, ukuran 749 KB (>200 KB) |
| 3 | Dev.to Markdown Article (English) | [dev.to/lab-34-chaos-engineering-devto.md](dev.to/lab-34-chaos-engineering-devto.md) | 2.498 kata, frontmatter YAML lengkap, syntax-highlighted Go code, link repo & commit hash, batas implementasi transparan |
| 4 | Medium HTML Article (Indonesian) | [medium/lab-34-chaos-engineering-medium.html](medium/lab-34-chaos-engineering-medium.html) | 2.430 kata (22.423 bytes), 22 section editorial lengkap, cover image, tabel komparasi, callout, tanpa JavaScript, styling Georgia/Arial |
| 5 | Cover Image (Indonesia) | [image-cover/lab-34-chaos-engineering-cover-id.png](image-cover/lab-34-chaos-engineering-cover-id.png) | 1920×1080 PNG retro CRT monitor, scanlines, subtitle Bahasa Indonesia, ukuran 733 KB (>50 KB) |
| 6 | Cover Image (English) | [image-cover/lab-34-chaos-engineering-cover-en.png](image-cover/lab-34-chaos-engineering-cover-en.png) | 1920×1080 PNG retro CRT monitor, scanlines, subtitle Bahasa Inggris, ukuran 729 KB (>50 KB) |
| 7 | Medium Metadata (Markdown) | [medium/lab-34-chaos-engineering-metadata.md](medium/lab-34-chaos-engineering-metadata.md) | Title, subtitle, SEO title (41 kar), SEO description (147 kar), slug, 5 tags, alt-text cover, social metadata |
| 8 | LinkedIn Caption (plain text) | [linkedin/lab-34-chaos-engineering-caption.md](linkedin/lab-34-chaos-engineering-caption.md) | 234 kata (target 180‑250), awalan persis *Lab 34 kali ini membahas …*, mental model Inggris, tanpa markdown bold |

## Catatan Verifikasi & Batas Implementasi
- **Commit SHA**: `b41fffba81521102d1ddfa66c51ba9b4cdae5758`
- **Status Test**: 5/5 test cases PASS (`go test -v -race ./...` durasi ~1.5s, clean tanpa data race)
- **Batas Implementasi**: Simulasi fault injection in-memory dalam runtime Go lokal; tidak memanipulasi layer jaringan kernel (iptables/eBPF); metrik kesehatan menggunakan counter kumulatif sederhana tanpa sliding time window ring-buffer.

## Caption LinkedIn (Siap Salin)
Lab 34 kali ini membahas Chaos Engineering dan Fault Injection, atau bagaimana sebuah service terdistribusi dapat membuktikan ketahanan secara empiris sebelum insiden produksi menimpa pelanggan.

Masalah teknis nyata terjadi ketika sebuah layanan bergantung pada payment gateway downstream yang mulai melambat atau mengembalikan error 500 beruntun. Koneksi HTTP client terjebak, connection pool habis dalam hitungan detik, dan kerusakan menyebar ke seluruh pohon call tree hulu. Pengujian integrasi statis dan mock unit saja tidak mampu mereproduksi fluktuasi jaringan yang tidak terduga ini.

Lab ini membangun eksperimen chaos deterministik di Go 1.22+ murni dengan standard library. Frameworknya menggabungkan tiga mekanisme kunci:

• Steady-state monitor atomik (sync/atomic) yang menilai error rate transaksional bukan CPU atau RAM.
• Circuit Breaker tiga fase (Closed, Open, Half-Open) yang memutus downstream yang sekarat dan memasukkan antrean fallback asinkron.
• Auto-abort runner berbasis ticker yang menetralkan fault injector dalam hitungan milidetik saat steady-state terlanggar.

Demo interaktif menunjukkan steady-state error rate tetap 0.00% di tengah simulasi pemadaman 100% di downstream, sementunya fase eksperimen tak terproteksi memicu pembatalan otomatis dan pemulihan penuh dalam 200ms cooldown.

Mental model: Build safety rails before injecting turbulence; measure customer-visible output, never component internals.

Batasan penting: Lab ini menggunakan injeksi in-memory dan counter kumulatif sederhana, bukan manipulasi lapisan jaringan (iptables/eBPF) atau sliding time window eksponensial seperti di lingkungan produksi sejati.

Pelajari kode sumber penuh dan lima test suite yang lolos race detector di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/34-chaos-engineering

#softwareengineering #golang #chaosengineering #distributedsystems #resilience #systemdesign
