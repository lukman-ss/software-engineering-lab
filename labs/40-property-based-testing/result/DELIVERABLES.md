# Publication Deliverables — Lab 40: Property-Based Testing

## Tabel Rekapitulasi Deliverables

| No | Tipe Deliverable | Lokasi File | Deskripsi & Verifikasi |
|---|---|---|---|
| 1 | PowerPoint Carousel (editable) | [linkedin/lab-40-property-based-testing-carousel.pptx](linkedin/lab-40-property-based-testing-carousel.pptx) | Tepat 24 slide 16:9, speaker notes lengkap di setiap slide, dark navy cover & closing, interior editorial light cream, ukuran 113 KB (>100 KB) |
| 2 | LinkedIn Carousel PDF (selectable) | [linkedin/lab-40-property-based-testing-linkedin.pdf](linkedin/lab-40-property-based-testing-linkedin.pdf) | Tepat 24 halaman 16:9 landscape, teks dapat diseleksi, render via Chrome headless print-to-pdf, ukuran 1.04 MB (>200 KB) |
| 3 | Dev.to Markdown Article (English) | [dev.to/lab-40-property-based-testing-devto.md](dev.to/lab-40-property-based-testing-devto.md) | 2.073 kata (>1.500 kata), frontmatter lengkap, 22 subbagian memetakan slide carousel, fenced code blocks Go/sh, link commit & repo |
| 4 | Medium HTML Article (Indonesian) | [medium/lab-40-property-based-testing-medium.html](medium/lab-40-property-based-testing-medium.html) | 2.383 kata (>2.000 kata), standalone HTML Georgia font, dark code blocks, 22 subbagian teknis, tabel komparasi, ukuran 21.8 KB |
| 5 | Cover Image (Indonesia) | [image-cover/lab-40-property-based-testing-cover-id.png](image-cover/lab-40-property-based-testing-cover-id.png) | 1920×1080 PNG monitor CRT retro beige dengan scanlines, tipografi rapi subtitle Bahasa Indonesia, ukuran 753 KB (>50 KB) |
| 6 | Cover Image (English) | [image-cover/lab-40-property-based-testing-cover-en.png](image-cover/lab-40-property-based-testing-cover-en.png) | 1920×1080 PNG monitor CRT retro beige dengan scanlines, tipografi rapi subtitle Bahasa Inggris, ukuran 747 KB (>50 KB) |
| 7 | Medium Metadata (Markdown) | [medium/lab-40-property-based-testing-metadata.md](medium/lab-40-property-based-testing-metadata.md) | Title, subtitle, SEO title (49 kar), SEO description (154 kar), slug, tepat 5 tags, alt-text cover, social metadata |
| 8 | LinkedIn Caption (Plain Text) | [linkedin/lab-40-property-based-testing-caption.md](linkedin/lab-40-property-based-testing-caption.md) | 204 kata (target 180–250 kata), awalan persis *Lab 40 kali ini membahas ...*, bullet points, mental model Bahasa Inggris, tanpa markdown bold |
| 9 | Python Content Definition | [content_def.py](content_def.py) | Definisi 24 slide, konstanta warna, metadata, dan narasi lengkap untuk builder deterministik (disimpan permanen) |

---

## Catatan Verifikasi & Batas Implementasi

- **Commit SHA**: `7d244968e91a2222291a27f346512976a26017e2`
- **Status Test**: 7/7 test cases PASS (`go test -v -race ./...` selesai tanpa data race)
  * `internal/currency`: `TestExampleBasedNaiveCurrency` (PASS), `TestPropertyNaiveCurrencyFails` (PASS), `TestPropertyRobustAmountRoundtrip` (PASS)
  * `internal/interval`: `TestExampleBasedIntervalMerge` (PASS), `TestPropertyNaiveMergeFails` (PASS), `TestPropertyRobustMergeIdempotence` (PASS), `TestPropertyRobustMergeNonOverlapping` (PASS)
  * `internal/shrinker`: `TestFindAndShrink` (PASS)
- **Batas Implementasi**:
  * Pengujian menggunakan paket standar Go `testing/quick` dengan generator kustom; tidak menggunakan library pihak ketiga eksternal seperti `gopter`.
  * Algoritma counterexample shrinking diimplementasikan khusus untuk slice bilangan bulat (`[]int`) melalui tiga strategi terurut: binary sectioning, element removal, dan numerical reduction toward zero.

---

## Teks Caption LinkedIn (Siap Salin)

```text
Lab 40 kali ini membahas Property-Based Testing (PBT), kelemahan mendasar unit test berbasis contoh, serta pembuktian kebenaran sistem via universal invariant dan counterexample shrinking di Go.

Masalah teknis yang sering terjadi di produksi adalah unit test konvensional memberikan rasa aman palsu. Pengembang hanya menguji skenario yang mereka antisipasi, sehingga kode berstatus lulus seratus persen meski menyimpan cacat batas domain. Pada contoh mata uang berbasis float64, format naif lulus seluruh contoh statis, namun gagal total saat menerima angka pecahan acak akibat keterbatasan aproksimasi biner IEEE 754. Begitu pula pada penggabungan interval, algoritma naif gagal pada 85 dari 100 irisan acak karena berasumsi data selalu terurut.

Lab ini membuktikan empat pilar utama:
• Roundtrip Invariant: Representasi integer cents int64 lulus 1.000 iterasi acak tanpa kehilangan presisi pecahan sedikit pun.
• Test Oracle & Idempotence: RobustMerge membuktikan relasi f(f(x)) == f(x) dan hasil non-overlapping pada 1.000 kombinasi irisan acak.
• Automated Shrinking: Algoritma pemotongan biner dan reduksi nilai memangkas array acak 10 elemen menjadi satu counterexample minimal [-1] hanya dalam 21 langkah evaluasi.
• Zero Dependency: Seluruh pengujian properti dibangun murni menggunakan pustaka standar Go testing/quick.

Mental model inti: Example-based tests check developer assumptions; property-based tests explore domain boundaries.

Pelajari implementasi lengkap dan hasil pengujian di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/40-property-based-testing

#softwareengineering #golang #testing #qualityassurance #systemdesign
```
