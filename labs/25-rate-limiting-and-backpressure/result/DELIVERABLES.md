# Publication Deliverables — Software Engineering Lab 25

## Tabel File Deliverable

| No | Tipe Deliverable | Lokasi File | Deskripsi |
|---|---|---|---|
| 1 | PowerPoint Carousel | [linkedin/lab-25-rate-limiting-and-backpressure-carousel.pptx](linkedin/lab-25-rate-limiting-and-backpressure-carousel.pptx) | File .pptx 24 slide 16:9 editable dengan Speaker Notes (via python-pptx) |
| 2 | LinkedIn Carousel PDF | [linkedin/lab-25-rate-limiting-and-backpressure-linkedin.pdf](linkedin/lab-25-rate-limiting-and-backpressure-linkedin.pdf) | File PDF 24 slide rasio 16:9 (Cover CRT + 22 slide krem editorial + Penutup navy) |
| 3 | Dev.to Article | [dev.to/lab-25-rate-limiting-and-backpressure-devto.md](dev.to/lab-25-rate-limiting-and-backpressure-devto.md) | Artikel Markdown Bahasa Inggris lengkap dengan frontmatter |
| 4 | Medium Article | [medium/lab-25-rate-limiting-and-backpressure-medium.html](medium/lab-25-rate-limiting-and-backpressure-medium.html) | Artikel HTML lengkap Bahasa Indonesia (~1.280 kata) |
| 5 | Cover Image (ID) | [image-cover/lab-25-rate-limiting-and-backpressure-cover-id.png](image-cover/lab-25-rate-limiting-and-backpressure-cover-id.png) | Cover PNG 1920×1080 retro CRT terminal (Bahasa Indonesia) |
| 6 | Cover Image (EN) | [image-cover/lab-25-rate-limiting-and-backpressure-cover-en.png](image-cover/lab-25-rate-limiting-and-backpressure-cover-en.png) | Cover PNG 1920×1080 retro CRT terminal (Bahasa Inggris) |
| 7 | Medium Metadata | [medium/lab-25-rate-limiting-and-backpressure-metadata.md](medium/lab-25-rate-limiting-and-backpressure-metadata.md) | Metadata SEO, slug, tags, social preview, & hitungan karakter |
| 8 | Caption LinkedIn | [linkedin/lab-25-rate-limiting-and-backpressure-caption.md](linkedin/lab-25-rate-limiting-and-backpressure-caption.md) | Caption teks 232 kata siap salin (tanpa markdown bold) |

## Catatan Verifikasi & Keterbatasan
- **Commit SHA**: `a39256cbf9ed820160fd28a0d842110780debda0`
- **Status Test**: 16/16 test cases PASS (`go test -race -v ./...` durasi ~6.5s, 0 data races).
- **Batas Implementasi**: Lab ini mengimplementasikan rate limiting dan queue backpressure in-memory Go pada satu host lokal. Wadah bucket disimpan dalam memori RAM instans dan buffer channel beroperasi di tingkat proses tunggal. Lab tidak menggunakan kluster Redis terdistribusi untuk sinkronisasi state lintas instance dan tidak menggunakan message broker eksternal (Kafka/RabbitMQ) untuk backpressure armada multi-node. Di lingkungan enterprise produksi, integrasi Redis Lua script untuk rate limiting dan monitoring consumer lag broker pesan direkomendasikan.

## Caption LinkedIn (Siap Salin)

```text
Lab 25 kali ini membahas strategi pertahanan beban bertingkat di Go: algoritma Token Bucket dan Leaky Bucket, bounded queue backpressure, penegakan HTTP 429 berbasis RFC 6585, dan exponential backoff dengan AWS Full Jitter.

Masalah teknis yang sering terjadi adalah lonjakan request yang melebihi kapasitas pemrosesan downstream worker. Tanpa mekanisme kendali beban, antrean memori tumbuh tanpa batas (unbounded queue), memicu degradasi latensi berbasis Little's Law, kegagalan memory exhaustion (OOM), dan badai retry tersinkronisasi (thundering herd) yang melipatgandakan beban saat sistem mencoba pulih.

Lab ini membuktikan empat pilar ketahanan sistem terdistribusi:
• Kendali Batas Token Bucket: Mengakomodasi lonjakan burst sesaat hingga kapasitas batas sambil mempertahankan laju rata-rata melalui penambahan token proporsional waktu.
• Isolasi Tenant RFC 6598: Middleware mengekstrak identitas X-API-Key untuk mencegah tabrakan kuota pada arsitektur Carrier-Grade NAT.
• Bounded Queue Backpressure: Buffer channel berukuran tetap menolak request baru secara non-blocking dengan ErrQueueFull saat antrean penuh.
• AWS Full Jitter Retry: Algoritma jitter mendistribusikan interval pengulangan klien secara seragam pada rentang eksponensial, melenyapkan sinkronisasi thundering herd.

Mental model inti: Rate limiting protects system boundaries from excessive ingress; bounded queue backpressure protects internal execution workers; full jitter retry prevents clients from synchronizing into destructive thundering herds.

Satu batasan penting: Lab ini menggunakan penyimpanan in-memory Go yang terikat pada satu host, sehingga deployment multi-node di lingkungan produksi memerlukan koordinasi terdistribusi seperti Redis dan message broker partisi Kafka.

Pelajari implementasi lengkap dan verifikasi pengujian di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/25-rate-limiting-and-backpressure

#softwareengineering #golang #systemdesign #distributedsystems #architecture #backend
```
