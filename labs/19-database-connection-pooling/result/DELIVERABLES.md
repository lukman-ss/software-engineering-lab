# Publication Deliverables — Software Engineering Lab 19

## Tabel File Deliverable

| No | Tipe Deliverable | Lokasi File | Deskripsi |
|---|---|---|---|
| 1 | PowerPoint Carousel | [linkedin/lab-19-database-connection-pooling-carousel.pptx](linkedin/lab-19-database-connection-pooling-carousel.pptx) | File .pptx 24 slide 16:9 editable dengan Speaker Notes (via python-pptx) |
| 2 | LinkedIn Carousel PDF | [linkedin/lab-19-database-connection-pooling-linkedin.pdf](linkedin/lab-19-database-connection-pooling-linkedin.pdf) | File PDF 24 slide rasio 16:9 (Cover CRT + 22 slide krem editorial + Penutup navy) |
| 3 | Dev.to Article | [dev.to/lab-19-database-connection-pooling-devto.md](dev.to/lab-19-database-connection-pooling-devto.md) | Artikel Markdown Bahasa Inggris lengkap dengan frontmatter |
| 4 | Medium Article | [medium/lab-19-database-connection-pooling-medium.html](medium/lab-19-database-connection-pooling-medium.html) | Artikel HTML lengkap Bahasa Indonesia (~1.397 kata) |
| 5 | Cover Image (ID) | [image-cover/lab-19-database-connection-pooling-cover-id.png](image-cover/lab-19-database-connection-pooling-cover-id.png) | Cover PNG 1920×1080 retro CRT terminal (Bahasa Indonesia) |
| 6 | Cover Image (EN) | [image-cover/lab-19-database-connection-pooling-cover-en.png](image-cover/lab-19-database-connection-pooling-cover-en.png) | Cover PNG 1920×1080 retro CRT terminal (Bahasa Inggris) |
| 7 | Medium Metadata | [medium/lab-19-database-connection-pooling-metadata.md](medium/lab-19-database-connection-pooling-metadata.md) | Metadata SEO, slug, tags, social preview, & hitungan karakter |
| 8 | Caption LinkedIn | [linkedin/lab-19-database-connection-pooling-caption.md](linkedin/lab-19-database-connection-pooling-caption.md) | Caption teks 218 kata siap salin (tanpa markdown bold) |

## Catatan Verifikasi & Keterbatasan
- **Commit SHA**: `ed2fe192cf5794f4430348df73e174611246be0e`
- **Status Test**: 10/10 test cases PASS (`go test -race -v ./tests/...` durasi 0.51s, 0 data races).
- **Batas Implementasi**: Lab ini memodelkan simulasi tingkat aplikasi menggunakan `database/sql/driver` mock dalam memori untuk menguji batas konkurensi (`maxConnections`) dan latensi inisialisasi jaringan (`connectDelay`) secara deterministik. Lab tidak mensimulasikan terminasi OOM proses backend database, Wait-For Graph (WFG) internal storage engine, maupun kompatibilitas penuh PgBouncer session mode. Rumus sizing hardware `((cores * 2) + spindles)` berfungsi sebagai baseline titik awal yang wajib divalidasi lewat load testing empiris.

## Caption LinkedIn (Siap Salin)

```text
Lab 19 kali ini membahas analisis mendalam database connection pooling, perumusan kapasitas berbasis hardware, serta mitigasi kegagalan fatal akibat kebocoran koneksi.

Masalah teknis yang sering terjadi adalah lonjakan error FATAL: sorry, too many clients already yang keliru direspons dengan memperbesar connection pool atau menaikkan max_connections server, sehingga latensi memburuk drastis dan server mengalami kehabisan slot.

Lab ini membuktikan empat hal konkret:
• Overhead Koneksi: Mode unpooled membutuhkan waktu berkali-kali lipat lebih lama akibat TCP handshake berulang, sedangkan pooled memakai ulang koneksi aktif seketika.
• Server Exhaustion: Mengizinkan pool aplikasi melebihi kapasitas max_connections server secara deterministik memicu penolakan koneksi ErrServerOverloaded.
• Connection Starvation: Menahan koneksi database saat menjalankan I/O eksternal lambat menyebabkan starvation dan timeout context deadline exceeded pada request lain yang sehat.
• Anti-Deadlock Sizing: Mengakuisisi koneksi kedua secara simultan pada pool berukuran 1 membuktikan risiko pool-locking deadlock yang memerlukan rumus kapasitas minimum.

Mental model inti: Connection capacity is dictated by database hardware constraints, not by client concurrency. Keep pools lean, reuse connections, and never hold sockets during external I/O.

Satu batasan penting: Lab ini menggunakan simulasi in-memory berbasis driver Go database/sql untuk menguji batas konkurensi dan latensi handshake, bukan proses database PostgreSQL nyata, serta belum memodelkan OOM process backend atau PgBouncer session mode.

Pelajari implementasi lengkap dan test assertion di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/19-database-connection-pooling

#softwareengineering #golang #database #backend #systemdesign #performance
```
