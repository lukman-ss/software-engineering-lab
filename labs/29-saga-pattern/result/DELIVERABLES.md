# Publication Deliverables — Lab 29: Saga Pattern

## Daftar File Deliverable (8 Deliverables + Python Definition)

| No | Tipe Deliverable | Lokasi File | Deskripsi | Ukuran |
|---|---|---|---|---|
| 0 | Python Slide Definition (Master) | [content_def.py](content_def.py) | Master script Python berisi definisi 24 slide, konstanta warna, metadata, dan mental models (DILARANG DIHAPUS) | 16 KB |
| 1 | PowerPoint Carousel (editable) | [linkedin/lab-29-saga-pattern-carousel.pptx](linkedin/lab-29-saga-pattern-carousel.pptx) | Tepat 24 slide 16:9 landscape, editable vector shapes/text, speaker notes lengkap di setiap slide | 116 KB |
| 2 | LinkedIn Carousel PDF (selectable) | [linkedin/lab-29-saga-pattern-linkedin.pdf](linkedin/lab-29-saga-pattern-linkedin.pdf) | Tepat 24 halaman 16:9, selectable text, dark navy cover & closing, editorial light interior, rendered via Chrome headless | 750 KB |
| 3 | Dev.to Markdown Article (English) | [dev.to/lab-29-saga-pattern-devto.md](dev.to/lab-29-saga-pattern-devto.md) | Artikel teknis bahasa Inggris, valid YAML frontmatter, 2.541 kata (>1.500 kata), syntax-highlighted Go code, link repo & commit hash | 21 KB |
| 4 | Medium HTML Article (Indonesian) | [medium/lab-29-saga-pattern-medium.html](medium/lab-29-saga-pattern-medium.html) | Artikel teknis bahasa Indonesia, standalone HTML Georgia/Arial styling, 4.250 kata (>2.000 kata), lengkap 9 section wajib, tanpa JS | 33 KB |
| 5 | Cover Image (Indonesia) | [image-cover/lab-29-saga-pattern-cover-id.png](image-cover/lab-29-saga-pattern-cover-id.png) | 1920×1080 PNG retro CRT terminal monitor, subtitle Bahasa Indonesia, rendered via Chrome screenshot | 760 KB |
| 6 | Cover Image (English) | [image-cover/lab-29-saga-pattern-cover-en.png](image-cover/lab-29-saga-pattern-cover-en.png) | 1920×1080 PNG retro CRT terminal monitor, subtitle Bahasa Inggris, rendered via Chrome screenshot | 768 KB |
| 7 | Medium Metadata (Markdown) | [medium/lab-29-saga-pattern-metadata.md](medium/lab-29-saga-pattern-metadata.md) | Title, subtitle, SEO title (45 kar), SEO description (154 kar), slug, 5 tags, alt-text, social metadata | 1.4 KB |
| 8 | LinkedIn Caption (plain text) | [linkedin/lab-29-saga-pattern-caption.md](linkedin/lab-29-saga-pattern-caption.md) | Tepat 233 kata (180–250 kata), awalan persis *Lab 29 kali ini membahas ...*, English mental model, tanpa markdown bold | 1.8 KB |

---

## Catatan Verifikasi & Batas Implementasi

- **Commit SHA**: `9c525944dc2de00e16513830a25ffa1c4834d742`
- **Status Test**: 9/9 test suites PASS (`go test -v -race ./...` in 0.428s, 0 data races)
  * `TestOrchestrator_HappyPath` (PASS)
  * `TestOrchestrator_FailureCompensatesLIFO` (PASS)
  * `TestPayment_Idempotency` (PASS)
  * `TestSemanticLock` (PASS)
  * `TestOrchestrator_Concurrency` (PASS)
  * `TestChoreography_Flow` (PASS)
  * `TestChoreography_FailureCompensates` (PASS)
  * `TestOrchestrator_CompensationErrorPropagated` (PASS)
  * `TestOrchestrator_ContextCancellation` (PASS)
- **Batas Implementasi Utama**:
  1. Lab ini menggunakan in-memory state dan data structure (`map` dan `sync.Mutex`) di Go, bukan database SQL/NoSQL fisik.
  2. Komunikasi antar service bersifat in-process function call, bukan gRPC atau HTTP lintas jaringan.
  3. EventBus menggunakan slice handler in-memory, bukan message broker terdistribusi seperti Apache Kafka atau RabbitMQ.
  4. Penanganan status `COMPENSATE_FAILED` dalam implementasi nyata memerlukan Dead-Letter Queue (DLQ) dan intervensi rekonsiliasi manual.

---

## Caption LinkedIn (Siap Salin)

Lab 29 kali ini membahas implementasi Saga Pattern untuk mengelola transaksi terdistribusi tanpa Two-Phase Commit pada arsitektur microservices berbasis Go.

Tantangan utama pada pola Database-per-Service adalah hilangnya jaminan transaksi ACID lintas layanan. Memaksakan protokol Two-Phase Commit memicu penguncian baris data secara sinkron lintas jaringan, meningkatkan latensi drastis, serta memunculkan titik kegagalan tunggal yang melumpuhkan ketersediaan sistem. Ketika proses pemesanan gagal di tengah jalan, database lokal yang telah melakukan commit tidak dapat dibatalkan melalui mekanisme rollback storage engine standar.

Lab ini membuktikan lima prinsip arsitektur terdistribusi:
• Orchestration Engine: Koordinator sentral mengeksekusi rangkaian transaksi lokal secara sekuensial dan memantau status setiap tahapan bisnis.
• LIFO Rollback: Kegagalan reservasi stok memicu kompensasi mundur dari pemotongan dana kartu kredit hingga pembatalan pesanan.
• Event-Driven Choreography: Komunikasi asinkron via EventBus memungkinkan layanan bereaksi mandiri terhadap domain event.
• Idempotency Protection: Pengecekan payment ID mencegah duplikasi pendebetan saldo saat terjadi retry request akibat network glitch.
• Semantic Locking: Status pesanan berstatus pending mengunci modifikasi konkuren sehingga mencegah terjadinya anomali data dirty reads.

Mental model inti: A saga is a sequence of local transactions where each step commits immediately. If failure occurs, compensating transactions execute in reverse order.

Satu batasan penting: Lab ini menguji logika koordinasi menggunakan in-memory state dan pemanggilan fungsi internal Go, bukan database relasional fisik atau message broker eksternal seperti Apache Kafka.

Pelajari implementasi dan seluruh test suite di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/29-saga-pattern

#softwareengineering #golang #microservices #distributedsystems #systemdesign #backend
