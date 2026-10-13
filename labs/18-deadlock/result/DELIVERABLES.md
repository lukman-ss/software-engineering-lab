# Publication Deliverables — Software Engineering Lab 18

## Tabel File Deliverable

| No | Tipe Deliverable | Lokasi File | Deskripsi |
|---|---|---|---|
| 1 | PowerPoint Carousel | [linkedin/lab-18-deadlock-carousel.pptx](linkedin/lab-18-deadlock-carousel.pptx) | File .pptx 24 slide 16:9 editable dengan Speaker Notes (via python-pptx) |
| 2 | LinkedIn Carousel PDF | [linkedin/lab-18-deadlock-linkedin.pdf](linkedin/lab-18-deadlock-linkedin.pdf) | File PDF 24 slide rasio 16:9 (Cover CRT + 22 slide krem + Penutup) |
| 3 | Dev.to Article | [dev.to/lab-18-deadlock-devto.md](dev.to/lab-18-deadlock-devto.md) | Artikel Markdown Bahasa Inggris lengkap dengan frontmatter |
| 4 | Medium Article | [medium/lab-18-deadlock-medium.html](medium/lab-18-deadlock-medium.html) | Artikel HTML lengkap Bahasa Indonesia (~1.250 kata) |
| 5 | Cover Image (ID) | [image-cover/lab-18-deadlock-cover-id.png](image-cover/lab-18-deadlock-cover-id.png) | Cover PNG 1920×1080 retro CRT terminal (Bahasa Indonesia) |
| 6 | Cover Image (EN) | [image-cover/lab-18-deadlock-cover-en.png](image-cover/lab-18-deadlock-cover-en.png) | Cover PNG 1920×1080 retro CRT terminal (Bahasa Inggris) |
| 7 | Medium Metadata | [medium/lab-18-deadlock-metadata.md](medium/lab-18-deadlock-metadata.md) | Metadata SEO, slug, tags, social preview, & hitungan karakter |
| 8 | Caption LinkedIn | [linkedin/lab-18-deadlock-caption.md](linkedin/lab-18-deadlock-caption.md) | Caption teks 201 kata siap salin (tanpa markdown bold) |

## Catatan Verifikasi & Keterbatasan
- **Commit SHA**: `6649acb946725e993e78f99ec47833a5498cc870`
- **Status Test**: 4/4 test cases PASS (`go test -v -race ./...` durasi 1.73s).
- **Batas Implementasi**: Lab ini memodelkan simulasi tingkat aplikasi menggunakan channel Go dan batas waktu context untuk mensimulasikan pemilihan transaksi korban oleh database engine. Tidak mengimplementasikan Wait-For Graph (WFG) internal storage engine dan menggunakan backoff tetap 2ms untuk kebutuhan test yang cepat dan deterministik.

## Caption LinkedIn (Siap Salin)

```text
Lab 18 kali ini membahas mekanisme terbentuknya deadlock pada sistem konkuren, penanganan oleh database engine, serta strategi pencegahan di tingkat aplikasi.

Masalah teknis yang sering terjadi adalah dua transaksi atau goroutine mengakses dan mengunci sumber daya bersama secara bersamaan dengan urutan berlawanan, menciptakan siklus circular wait yang membekukan pemrosesan transaksi.

Lab ini membuktikan empat hal konkret:
• Deadlock Occurrence: Transfer naif dua arah yang dieksekusi simultan terbukti memicu kegagalan deadlock victim saat context timeout terlampaui.
• Lock Ordering: Pengurutan akuisisi lock berdasarkan identitas akun memutus terbentuknya siklus secara deterministik dan menjaga konsistensi saldo.
• Application Retry: Mekanisme retry pada tingkat aplikasi berhasil memulihkan transaksi yang sempat dibatalkan oleh database engine.
• Transaction Duration: Pengujian durasi transaksi membuktikan bahwa menahan lock lebih lama secara nyata melipatgandakan peluang terjadinya deadlock.

Mental model inti: Deadlock is a circular dependency. Database engine membatalkan salah satu transaksi sebagai victim. Lock ordering mencegah siklus terbentuk, sementara application retry memulihkan transaksi korban.

Satu batasan penting: Lab ini adalah simulasi tingkat aplikasi berbasis channel Go dan context timeout, bukan implementasi Wait-For Graph (WFG) internal database engine, serta menggunakan jeda retry tetap 2ms untuk kebutuhan pengujian.

Pelajari implementasi lengkap dan test assertion di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/18-deadlock

#softwareengineering #golang #concurrency #database #systemdesign #backend
```
