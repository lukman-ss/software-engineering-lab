# Rekapitulasi Deliverables — Lab 28: Timeouts & Deadlines

Seluruh berkas keluaran untuk **Software Engineering Lab 28: Timeouts & Deadlines** telah selesai dibuat dan diverifikasi sesuai standar mutu publikasi Lab 14 & 15.

---

## 1. Tabel Berkas Publikasi

| No | Deliverable | Jalur Berkas | Format / Ukuran | Deskripsi |
|---|---|---|---|---|
| 0 | Slide & Metadata Definition | [`content_def.py`](content_def.py) | Python (7.8 KB) | Sumber kode reproducible untuk 24 slide, palet warna, dan metadata |
| 1 | Artikel Dev.to | [`dev.to/lab-28-timeouts-and-deadlines-devto.md`](dev.to/lab-28-timeouts-and-deadlines-devto.md) | Markdown (2,873 kata) | Artikel teknis mendalam Bahasa Inggris lengkap dengan frontmatter |
| 2 | Artikel Medium | [`medium/lab-28-timeouts-and-deadlines-medium.html`](medium/lab-28-timeouts-and-deadlines-medium.html) | HTML Standalone (3,829 kata) | Artikel komprehensif Bahasa Indonesia mencakup 9 section wajib |
| 3 | Metadata Medium | [`medium/lab-28-timeouts-and-deadlines-metadata.md`](medium/lab-28-timeouts-and-deadlines-metadata.md) | Markdown (1.2 KB) | Konfigurasi SEO, judul alternatif, slug, dan tepat 5 topik tag |
| 4 | Cover Visual ID | [`image-cover/lab-28-timeouts-and-deadlines-cover-id.png`](image-cover/lab-28-timeouts-and-deadlines-cover-id.png) | PNG 1920×1080 (759 KB) | Visual retro monitor CRT bezel dengan tipografi Bahasa Indonesia |
| 5 | Cover Visual EN | [`image-cover/lab-28-timeouts-and-deadlines-cover-en.png`](image-cover/lab-28-timeouts-and-deadlines-cover-en.png) | PNG 1920×1080 (752 KB) | Visual retro monitor CRT bezel dengan tipografi Bahasa Inggris |
| 6 | Slide Deck LinkedIn PDF | [`linkedin/lab-28-timeouts-and-deadlines-linkedin.pdf`](linkedin/lab-28-timeouts-and-deadlines-linkedin.pdf) | PDF 16:9 (553 KB) | Tepat 24 halaman landscape resolusi tinggi, selectable text |
| 7 | Slide Carousel PPTX | [`linkedin/lab-28-timeouts-and-deadlines-carousel.pptx`](linkedin/lab-28-timeouts-and-deadlines-carousel.pptx) | PPTX 16:9 (115 KB) | Tepat 24 slide presentasi lengkap dengan Speaker Notes di tiap slide |
| 8 | Caption LinkedIn | [`linkedin/lab-28-timeouts-and-deadlines-caption.md`](linkedin/lab-28-timeouts-and-deadlines-caption.md) | Markdown (216 kata) | Naskah posting LinkedIn siap salin tanpa bold markdown |

---

## 2. Catatan Teknis & Verifikasi Mutu

- **Commit SHA**: `bdb285d35a71c9a309e83a0e4497b9dbd842ceb8`
- **Hasil Pengujian**: 13/13 unit dan integration test berstatus **PASS** (`go test -v -race ./...`).
- **Race Detector**: Bersih tanpa data race terdeteksi (`PASS`).
- **Demo Eksekusi**: Program `cmd/demo/main.go` berhasil menjalankan 4 skenario pembuktian secara deterministik.
- **Batas Implementasi Utama**:
  1. *In-Memory Deduplication*: Idempotency key disimpan dalam memori RAM lokal dengan `sync.Mutex` dan lazy TTL eviction, belum menggunakan Redis cluster terdistribusi.
  2. *Single-Process Context*: Pewarisan context deadline diuji dalam proses tunggal Go runtime; integrasi wire-protocol lintas jaringan (seperti header `grpc-timeout` atau `Request-Timeout`) dimodelkan secara konseptual.
  3. *Static Time Budget*: Budget durasi bersifat statis (tidak memakai adaptive moving-window latency percentile).

---

## 3. Teks Siap Salin: Caption LinkedIn

```text
Lab 28 kali ini membahas mitigasi cascading failure akibat dependency lambat melalui penerapan context deadline, exponential backoff dengan full jitter, circuit breaker, dan idempotency key di Go.

Masalah teknis yang sering terjadi adalah satu dependency eksternal melambat dari 300ms menjadi 60 detik. Berdasarkan Little's Law (L = lambda * W), jumlah request in-flight meledak dan menguras seluruh worker thread pool, sehingga endpoint independen ikut lumpuh total.

Lab ini membuktikan empat pilar ketahanan konkret:
• Context Deadline Propagation: Child context secara hierarkis mewarisi sisa deadline parent dan melepaskan goroutine melalui buffered channel tanpa memory leak.
• Full Jitter Exponential Backoff: Distribusi acak seragam sleep = random(0, min(max, base * 2^(attempt-1))) membubarkan thundering herd retry storm saat dependency pulih.
• Three-State Circuit Breaker: Transisi CLOSED, OPEN, dan HALF_OPEN memotong rantai retry seketika saat downstream terbukti down, memberi ruang isolasi pemulihan.
• Idempotency Deduplication: Penyimpanan hash response transaksi melindungi operasi finansial dari risiko double charge ketika status timeout bersifat indeterminate.

Mental model inti: Fast failure frees resources immediately; slow dependencies hoard threads until total starvation. Timeouts are resource guardrails, not arbitrary numbers.

Satu batasan penting: Lab ini menggunakan in-memory storage untuk idempotency key dalam satu proses Go, bukan Redis terdistribusi atau cluster database transaksional.

Pelajari implementasi lengkap dan seluruh pembuktian test di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/28-timeouts-and-deadlines

#softwareengineering #golang #distributedsystems #systemdesign #resilience #microservices
```
