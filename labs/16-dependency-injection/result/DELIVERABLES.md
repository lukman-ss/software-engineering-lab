# Publication Deliverables — Software Engineering Lab 16

## Tabel File Deliverable

| No | Tipe Deliverable | Lokasi File | Deskripsi |
|---|---|---|---|
| 1 | PowerPoint Carousel | [linkedin/lab-16-dependency-injection-carousel.pptx](linkedin/lab-16-dependency-injection-carousel.pptx) | File .pptx 24 slide 16:9 editable dengan Speaker Notes |
| 2 | LinkedIn Carousel PDF | [linkedin/lab-16-dependency-injection-linkedin.pdf](linkedin/lab-16-dependency-injection-linkedin.pdf) | File PDF 24 slide (Cover CRT + 22 slide krem + Penutup) |
| 3 | Dev.to Article | [dev.to/lab-16-dependency-injection-devto.md](dev.to/lab-16-dependency-injection-devto.md) | Artikel Markdown Bahasa Inggris dengan frontmatter |
| 4 | Medium Article | [medium/lab-16-dependency-injection-medium.html](medium/lab-16-dependency-injection-medium.html) | Artikel HTML lengkap Bahasa Indonesia (~1.230 kata) |
| 5 | Cover Image (ID) | [image-cover/lab-16-dependency-injection-cover-id.png](image-cover/lab-16-dependency-injection-cover-id.png) | Cover PNG 1920×1080 tema retro CRT (Bahasa Indonesia) |
| 6 | Cover Image (EN) | [image-cover/lab-16-dependency-injection-cover-en.png](image-cover/lab-16-dependency-injection-cover-en.png) | Cover PNG 1920×1080 tema retro CRT (Bahasa Inggris) |
| 7 | Medium Metadata | [medium/lab-16-dependency-injection-metadata.md](medium/lab-16-dependency-injection-metadata.md) | Metadata SEO, slug, tags, social preview, & hitungan karakter |
| 8 | Caption LinkedIn | [linkedin/lab-16-dependency-injection-caption.md](linkedin/lab-16-dependency-injection-caption.md) | Caption teks 189 kata siap salin (tanpa markdown bold) |

## Catatan Verifikasi & Keterbatasan
- **Commit SHA**: `8afe2a824667528d16d6d9b2ac434a18a6ac2db1`
- **Status Test**: 6/6 test cases PASS (`go test -v -race ./...` durasi 1.44s).
- **Batas Implementasi**: Lab membuktikan manual constructor wiring di Go standar dan isolasi mock tanpa jaringan. Tidak mencakup framework DI otomatis (reflection/code-gen) atau container lifecycle rumit.

## Caption LinkedIn (Siap Salin)

```text
Lab 16 kali ini membahas pemisahan konstruksi objek dari logika bisnis melalui Dependency Injection dan Inversion of Control di Go.

Masalah teknis yang sering terjadi adalah dependensi infrastruktur diinstansiasi langsung di dalam method domain, membuat kode terkunci pada provider tertentu dan pengujian unit terpaksa bergantung pada jaringan eksternal.

Lab ini membuktikan tiga hal konkret:
• Constructor Injection: Processor menerima interface PaymentGateway secara eksplisit sehingga state objek selalu lengkap dan dependensi tampak jelas.
• Service Locator Anti-Pattern: BadProcessor menerima Container generik, menyembunyikan dependensi nyata di balik query method runtime.
• Fast Isolated Testing: MockGateway memungkinkan 6 test assertions berjalan instan tanpa koneksi jaringan, membuktikan error propagation serta validasi input yang menolak transaksi sebelum memanggil gateway.

Mental model inti: Externalize object creation. Biarkan assembler luar merakit object graph, sementara logika domain hanya fokus menjalankan kontrak bisnis. Nilai data murni seperti struct Money diinstansiasi langsung tanpa melalui container DI.

Satu batasan penting: Lab ini menggunakan manual wiring di Go standar dan tidak mendemonstrasikan automated DI framework (seperti reflection atau code-gen), serta tidak mengelola container lifecycle (singleton, scoped, transient).

Pelajari implementasi lengkap dan test assertion di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/16-dependency-injection

#softwareengineering #golang #systemdesign #backend #architecture #designpatterns
```
