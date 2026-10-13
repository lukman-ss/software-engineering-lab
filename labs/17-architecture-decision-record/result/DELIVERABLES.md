# Publication Deliverables — Software Engineering Lab 17

## Tabel File Deliverable

| No | Tipe Deliverable | Lokasi File | Deskripsi |
|---|---|---|---|
| 1 | PowerPoint Carousel | [linkedin/lab-17-architecture-decision-record-carousel.pptx](linkedin/lab-17-architecture-decision-record-carousel.pptx) | File .pptx 24 slide 16:9 editable dengan Speaker Notes |
| 2 | LinkedIn Carousel PDF | [linkedin/lab-17-architecture-decision-record-linkedin.pdf](linkedin/lab-17-architecture-decision-record-linkedin.pdf) | File PDF 24 slide (Cover CRT + 22 slide krem + Penutup) |
| 3 | Dev.to Article | [dev.to/lab-17-architecture-decision-record-devto.md](dev.to/lab-17-architecture-decision-record-devto.md) | Artikel Markdown Bahasa Inggris dengan frontmatter |
| 4 | Medium Article | [medium/lab-17-architecture-decision-record-medium.html](medium/lab-17-architecture-decision-record-medium.html) | Artikel HTML lengkap Bahasa Indonesia padat (~1.550 kata, 22 seksi, 5 tabel) |
| 5 | Cover Image (ID) | [image-cover/lab-17-architecture-decision-record-cover-id.png](image-cover/lab-17-architecture-decision-record-cover-id.png) | Cover PNG 1920×1080 tema retro CRT (Bahasa Indonesia) |
| 6 | Cover Image (EN) | [image-cover/lab-17-architecture-decision-record-cover-en.png](image-cover/lab-17-architecture-decision-record-cover-en.png) | Cover PNG 1920×1080 tema retro CRT (Bahasa Inggris) |
| 7 | Medium Metadata | [medium/lab-17-architecture-decision-record-metadata.md](medium/lab-17-architecture-decision-record-metadata.md) | Metadata SEO, slug, tags, social preview, & hitungan karakter |
| 8 | Caption LinkedIn | [linkedin/lab-17-architecture-decision-record-caption.md](linkedin/lab-17-architecture-decision-record-caption.md) | Caption teks siap salin (tanpa markdown bold) |

## Catatan Verifikasi & Keterbatasan
- **Commit SHA**: `eb612371e53183397649b7091a33d4068f70a7cb`
- **Status Test**: Seluruh parser & linter assertions PASS (`go test -v -race ./...`).
- **Batas Implementasi**: Lab membuktikan algoritma regex parsing MADR, model state machine 5 status formal, penomoran monotonik, dan validasi timbal-balik graf silsilah konkuren. Tidak mencakup webhook GitHub Actions otomatis atau parser format AsciiDoc.

## Caption LinkedIn (Siap Salin)

```text
Lab 17 kali ini membahas pengelolaan keputusan arsitektur menggunakan Architecture Decision Record (ADR) serta validasi struktural otomatis di Go.

Masalah teknis yang sering terjadi adalah hilangnya konteks historis di balik keputusan desain penting saat anggota tim berganti. Dokumentasi yang terpisah di wiki luar cepat usang, memicu perdebatan teknis berulang atau pembatalan arsitektur secara sembarangan tanpa memahami batasan aslinya.

Lab ini membuktikan tiga mekanisme konkret:
1. Regex-Based Parsing: Mengekstrak ID, Title, Status, serta klausa silsilah SupersededBy dan Supersedes dari berkas Markdown, sekaligus memastikan bagian Context, Decision, dan Consequences tidak kosong.
2. Monotonic & Lineage Integrity: Linter mendeteksi penomoran non-monotonik, ID duplikat, tautan sepihak, target fiktif, self-supersession, serta siklus penggantian tertutup melalui DFS cycle detection.
3. Concurrent Verification: Validasi graf silsilah timbal-balik dieksekusi paralel dengan goroutine dan lolos seluruh test assertions di bawah race detector tanpa peringatan data race.

Mental model inti: Decisions as versioned invariants. Simpan keputusan arsitektural langsung bersama kode sumber sebagai catatan permanen yang tidak dapat disunting retroaktif, melainkan digantikan melalui silsilah penggantian timbal-balik yang divalidasi mesin. Usulan yang ditolak dicatat sebagai status Rejected agar argumen penolakan terdokumentasi.

Satu batasan penting: Lab ini menggunakan parser ekspresi reguler dengan format heading Markdown yang ketat dan bekerja pada struktur in-memory, belum menyediakan AST CommonMark umum atau pengait otomatis git pre-commit hook.

Pelajari implementasi lengkap dan rangkaian pengujian di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/17-architecture-decision-record

#softwareengineering #golang #systemdesign #backend #architecture #cleancode
```
