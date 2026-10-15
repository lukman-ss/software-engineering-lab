# Rekapitulasi Deliverables — Lab 25: Rate Limiting & Backpressure

Seluruh berkas keluaran publikasi untuk **Software Engineering Lab 25: Rate Limiting & Backpressure** telah selesai dibuat dan diverifikasi sesuai standar mutu publikasi Lab 13–15 dan panduan modular pipeline.

---

## 1. Tabel Berkas Publikasi

| No | Deliverable | Jalur Berkas | Format / Ukuran | Deskripsi |
|---|---|---|---|---|
| 0 | Slide & Metadata Definition | [`content_def.py`](content_def.py) | Python (21 KB) | Definisi 24 slide editorial, visual metadata, limiting phrases, dan source mapping |
| 1 | Cover Visual Master | [`image-cover/cover_art.jpg`](image-cover/cover_art.jpg) | JPEG 1920×1080 (325 KB) | Hero visual retro monitor CRT terminal amber/cyan dengan scanlines fotorealistik |
| 2 | Cover Poster Bahasa Indonesia | [`image-cover/lab-25-rate-limiting-and-backpressure-cover-id.png`](image-cover/lab-25-rate-limiting-and-backpressure-cover-id.png) | PNG 1920×1080 (1.7 MB) | Poster sinematik Bahasa Indonesia, satu hero visual, ZERO micro-text |
| 3 | Cover Poster Bahasa Inggris | [`image-cover/lab-25-rate-limiting-and-backpressure-cover-en.png`](image-cover/lab-25-rate-limiting-and-backpressure-cover-en.png) | PNG 1920×1080 (1.7 MB) | Poster sinematik Bahasa Inggris, satu hero visual, ZERO micro-text |
| 4 | Slide Deck LinkedIn PDF | [`linkedin/lab-25-rate-limiting-and-backpressure-linkedin.pdf`](linkedin/lab-25-rate-limiting-and-backpressure-linkedin.pdf) | PDF 16:9 (2.4 MB) | Tepat 24 slide landscape resolusi tinggi dengan layout editorial |
| 5 | Slide Carousel PPTX | [`linkedin/lab-25-rate-limiting-and-backpressure-carousel.pptx`](linkedin/lab-25-rate-limiting-and-backpressure-carousel.pptx) | PPTX 16:9 (112 KB) | Tepat 24 slide presentasi lengkap dengan Speaker Notes di tiap slide |
| 6 | Caption LinkedIn | [`linkedin/lab-25-rate-limiting-and-backpressure-caption.md`](linkedin/lab-25-rate-limiting-and-backpressure-caption.md) | Markdown (1.6 KB) | Naskah posting LinkedIn siap salin tanpa bold markdown formatting |
| 7 | Artikel Dev.to | [`dev.to/lab-25-rate-limiting-and-backpressure-devto.md`](dev.to/lab-25-rate-limiting-and-backpressure-devto.md) | Markdown (14 KB) | Artikel teknis mendalam Bahasa Inggris lengkap dengan frontmatter |
| 8 | Artikel Medium | [`medium/lab-25-rate-limiting-and-backpressure-medium.html`](medium/lab-25-rate-limiting-and-backpressure-medium.html) | HTML Standalone (15 KB) | Artikel komprehensif Bahasa Indonesia mencakup tabel komparasi, kode, dan arsitektur |
| 9 | Metadata Medium | [`medium/lab-25-rate-limiting-and-backpressure-metadata.md`](medium/lab-25-rate-limiting-and-backpressure-metadata.md) | Markdown (1.5 KB) | Konfigurasi SEO, judul alternatif, slug, dan tepat 5 topik tag |

---

## 2. Catatan Teknis & Verifikasi Mutu

- **Commit SHA**: `45cc429`
- **Hasil Pengujian**: 15/15 unit test berstatus **PASS** (`go test -v -race ./...`).
- **Race Detector**: Bersih tanpa data race terdeteksi (`PASS` dalam 0.28 detik).
- **Demo Eksekusi**: Program `cmd/demo/main.go` membuktikan burst token bucket, smoothing leaky bucket, fast rejection bounded queue, dan sebaran AWS Full Jitter.
- **Kepatuhan Mutu Publikasi**:
  1. *Series Identity Lock*: Wajib CRT raster cover art (`cover_art.jpg`) tanpa elemen geometris vektor generik.
  2. *Slide Rules*: Rata-rata kata per slide terkontrol (<= 45 kata narasi, <= 8 baris kode dominan), layout beragam (cards, schema, table, code, checklist, closing).
  3. *Quality Bar*: Memuat >= 12 slide dengan limiting phrases eksplisit ("bukan berarti", "tidak otomatis", "hanya pada jalur ini", "belum membuktikan", "simulasi ini tidak sama dengan", "production masih perlu").
  4. *Batas Implementasi*:
     - *In-Memory Architecture*: Menggunakan mutex Go lokal dalam single process; cluster multi-pod produksi membutuhkan Redis cluster (Lua script / atomic hash).
     - *Konstanta Delay AWS*: Delay base 100ms dan cap 1000ms adalah referensi implementasi yang perlu dikalibrasi terhadap SLA layanan nyata.

---

## 3. Teks Siap Salin: Caption LinkedIn

```text
Lab 25 kali ini membahas strategi pengendalian beban terdistribusi melalui implementasi Token Bucket, Leaky Bucket, Bounded Queue backpressure, dan AWS Full Jitter retry di Go.

Masalah mendasar pada sistem tanpa pembatas adalah ketimpangan laju: ketika trafik masuk melompat melampaui kemampuan downstream, antrean bertumpuk tanpa batas di RAM. Berdasarkan Little's Law (L = lambda * W), 5.000.000 request dengan processing rate 2.000/detik membutuhkan waktu 41 menit 40 detik untuk menguras buffer, memicu lonjakan latensi masif dan crash akibat Out of Memory (OOM).

Lab ini mengimplementasikan empat pilar pertahanan konkret murni menggunakan Go standard library:
• Token Bucket & Leaky Bucket: Pengisian waktu berjalan (lazy monotonic refill) mengizinkan burst terukur atau memaksakan aliran output konstan tanpa background ticker yang boros CPU.
• Isolasi Kunci Per-Tenant: Registri bucket berbasis API key (X-API-Key) mencegah fenomena collateral throttling akibat ribuan pengguna berbagi IP publik di balik Carrier-Grade NAT (RFC 6598).
• Bounded Queue Backpressure: Buffer berkapasitas tetap dengan pola select-default menolak kelebihan pekerjaan seketika (ErrQueueFull) tanpa alokasi memori tambahan.
• AWS Full Jitter Retry: Randomisasi jeda tidur sleep = random(0, min(cap, base * 2^attempt)) membubarkan sinkronisasi thundering herd saat downstream pulih.
• HTTP 429 RFC 6585 Middleware: Menyertakan header Retry-After presisi untuk memandu penjadwalan ulang request klien.

Mental model inti: Rate limiting mengontrol pintu gerbang masuk; backpressure mempropagasi sinyal penolakan ke belakang; Full Jitter mendesinkronisasi gelombang percobaan ulang klien.

Batasan penting implementasi: Lab ini menggunakan state in-memory single process. Pada skala multi-pod produksi terdistribusi, koordinasi token bucket membutuhkan shared state eksternal seperti Redis cluster dengan skrip atomic Lua.

Pelajari arsitektur lengkap, kode implementasi, dan pengujian race-clean di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/25-rate-limiting-and-backpressure

#softwareengineering #golang #distributedsystems #systemdesign #resilience #performance #backend
```
