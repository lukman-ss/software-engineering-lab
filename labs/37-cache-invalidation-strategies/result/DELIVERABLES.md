# Deliverables Rekapitulasi: Lab 37 — Cache Invalidation Strategies

Paket deliverable publikasi lengkap untuk **Software Engineering Lab 37: Cache Invalidation Strategies — Pola Write, Mitigasi Stampede, dan Jitter TTL dalam Go** oleh **Lukman / lukman-ss**.

---

## Tabel Unduhan & Berkas Publikasi

| No | Tipe Deliverable | File Output | Status / Ukuran | Deskripsi |
| :---: | :--- | :--- | :--- | :--- |
| **0** | **Python Definition** | [`result/content_def.py`](content_def.py) | 59.4 KB (848 baris) | Definisi 24 slide, konstanta warna, metadata, dan mental model (DILARANG DIHAPUS) |
| **1** | **Dev.to Article** | [`result/dev.to/lab-37-cache-invalidation-strategies-devto.md`](dev.to/lab-37-cache-invalidation-strategies-devto.md) | ~15.4 KB (2.115 kata) | Artikel teknis Bahasa Inggris lengkap dengan frontmatter dan 22 subbagian |
| **2** | **Medium Article** | [`result/medium/lab-37-cache-invalidation-strategies-medium.html`](medium/lab-37-cache-invalidation-strategies-medium.html) | ~23.4 KB (2.543 kata) | Standalone HTML Bahasa Indonesia dengan 22 section dan styling Georgia/Code |
| **3** | **Medium Metadata** | [`result/medium/lab-37-cache-invalidation-strategies-metadata.md`](medium/lab-37-cache-invalidation-strategies-metadata.md) | Markdown | Metadata SEO, slug, social title, dan 5 tags Medium |
| **4** | **Cover Image (ID)**| [`result/image-cover/lab-37-cache-invalidation-strategies-cover-id.png`](image-cover/lab-37-cache-invalidation-strategies-cover-id.png) | 751 KB (1920x1080) | Cover monitor CRT retro pixel-art dengan subtitle Bahasa Indonesia |
| **5** | **Cover Image (EN)**| [`result/image-cover/lab-37-cache-invalidation-strategies-cover-en.png`](image-cover/lab-37-cache-invalidation-strategies-cover-en.png) | 754 KB (1920x1080) | Cover monitor CRT retro pixel-art dengan subtitle Bahasa Inggris |
| **6** | **LinkedIn PDF** | [`result/linkedin/lab-37-cache-invalidation-strategies-linkedin.pdf`](linkedin/lab-37-cache-invalidation-strategies-linkedin.pdf) | 1.171 KB (24 slide) | Carousel dokumen PDF 16:9 24 halaman (Cover Dark Navy, 22 Editorial Cream, Closing Dark Navy) |
| **7** | **LinkedIn PPTX** | [`result/linkedin/lab-37-cache-invalidation-strategies-carousel.pptx`](linkedin/lab-37-cache-invalidation-strategies-carousel.pptx) | 119 KB (24 slide) | File presentasi PowerPoint 24 slide dengan layout 16:9 dan Speaker Notes |
| **8** | **LinkedIn Caption**| [`result/linkedin/lab-37-cache-invalidation-strategies-caption.md`](linkedin/lab-37-cache-invalidation-strategies-caption.md) | 240 kata | Teks caption LinkedIn tanpa bold markdown siap salin |

---

## Verifikasi Kualitas & Catatan Teknis

1. **Commit Snapshot**: `eac3d80a`
2. **Status Pengujian**:
   - `go test -v ./...` -> 16 passing tests (race-clean verified)
   - `go test -race ./...` -> PASS (Race-Clean Verified Go 1.22, 0 data races)
   - `go run ./cmd/demo` -> Simulasi seluruh pola (Cache-Aside, Write-Through, Write-Behind, Singleflight, XFetch, SWR, TTL Jitter) berjalan sukses
3. **Batas Implementasi Utama**:
   - Seluruh modul menggunakan simulasi in-memory cache dan MockDB Go murni, bukan cluster Redis atau database PostgreSQL nyata.
   - Buffer antrean Write-Behind menggunakan channel Go in-memory; crash pada proses menyebabkan data yang belum di-flush hilang (sistem produksi memerlukan durable WAL / Kafka).
   - Singleflight beroperasi pada level single process/node; sistem multi-pod membutuhkan distributed locking (misal Redis Redlock / SETNX).

---

## Salinan Teks Caption LinkedIn

```text
Lab 37 kali ini membahas strategi cache invalidation dalam sistem terdistribusi, dengan fokus pada pola write, mitigasi stampede, dan teknik TTL jitter dalam implementasi Go.

Cache invalidation adalah salah satu dilema terbesar dalam arsitektur sistem terdistribusi skala tinggi. Tanpa mekanisme yang disiplin dan terukur, sistem akan menghadapi tiga masalah konkret yang fatal:
• Cache stampede: Ribuan request konkuren membanjiri database secara serentak saat kunci populer kedaluwarsa mendadak
• Stale data: Goroutine pembaca menerima data usang karena inkonsistensi jeda waktu antara penulisan database dan pembaruan cache
• Risiko kehilangan data: Pilihan throughput ekstrem pada Write-Behind mengorbankan durabilitas saat antrean memory overflow atau proses crash

Eksperimen Go standard library pada lab ini membuktikan implementasi native:
• Cache-Aside meminimalkan penggunaan memori namun rentan terhadap race condition pembacaan usang
• Write-Through menjamin konsistensi read-after-write sinkron dengan kompromi kenaikan latensi penulisan
• Write-Behind mengalirkan mutasi asinkron ke worker latar belakang dengan throughput tinggi
• Singleflight memadatkan 20 request konkuren menjadi tepat 1 query database tunggal
• XFetch menggeser titik rebuild proaktif secara probabilistik sebelum batas waktu kedaluwarsa
• Stale-While-Revalidate menyajikan data usang seketika sambil merevalidasi kunci di latar belakang
• TTL Jitter mengacak rentang waktu kedaluwarsa untuk mengeliminasi lonjakan serentak

Mental model: Cache invalidation is the hardest problem because consistency is eventual, never instant.

Batasan implementasi utama: Seluruh simulasi menggunakan memory cache dan MockDB Go murni tanpa Redis atau PostgreSQL persisten eksternal.

Pelajari kode lengkap dan pengujian di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/37-cache-invalidation-strategies

#softwareengineering #golang #systemdesign #distributed #performance #caching
```
