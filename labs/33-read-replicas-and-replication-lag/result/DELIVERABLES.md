# Publication Deliverables — Lab 33: Read Replicas and Replication Lag

## Daftar File Deliverable (24 slide, 2 cover, 2 dokumen teks, 1 PDF, 1 PPTX, 1 metadata)

| No | Tipe Deliverable | Lokasi File | Deskripsi |
|---|---|---|---|
| 1 | PowerPoint Carousel (editable) | [linkedin/lab-33-read-replicas-and-replication-lag-carousel.pptx](linkedin/lab-33-read-replicas-and-replication-lag-carousel.pptx) | 24 slide 16:9, speaker notes lengkap, dark navy cover, editorial light slides 2‑23, dark navy penutup, ukuran 115 KB (>100 KB) |
| 2 | LinkedIn Carousel PDF (selectable text) | [linkedin/lab-33-read-replicas-and-replication-lag-linkedin.pdf](linkedin/lab-33-read-replicas-and-replication-lag-linkedin.pdf) | 24 halaman 16:9, dark navy cover & penutup, editorial light interior, ukuran 879 KB (>200 KB) |
| 3 | Dev.to Markdown Article (English) | [dev.to/lab-33-read-replicas-and-replication-lag-devto.md](dev.to/lab-33-read-replicas-and-replication-lag-devto.md) | 1.567 kata, frontmatter YAML lengkap, syntax‑highlighted Go/sh code, link repo & commit, batas implementasi transparan |
| 4 | Medium HTML Article (Indonesian) | [medium/lab-33-read-replicas-and-replication-lag-medium.html](medium/lab-33-read-replicas-and-replication-lag-medium.html) | 2.661 kata, 9 section wajib lengkap, cover image, tabel komparasi, callout, tanpa JavaScript, styling Georgia/Arial |
| 5 | Cover Image (Indonesia) | [image-cover/lab-33-read-replicas-and-replication-lag-cover-id.png](image-cover/lab-33-read-replicas-and-replication-lag-cover-id.png) | 1920×1080 PNG retro CRT monitor, scanlines, subtitle Bahasa Indonesia, ukuran 764 KB (>50 KB) |
| 6 | Cover Image (English) | [image-cover/lab-33-read-replicas-and-replication-lag-cover-en.png](image-cover/lab-33-read-replicas-and-replication-lag-cover-en.png) | 1920×1080 PNG retro CRT monitor, scanlines, subtitle Bahasa Inggris, ukuran 759 KB (>50 KB) |
| 7 | Medium Metadata (Markdown) | [medium/lab-33-read-replicas-and-replication-lag-metadata.md](medium/lab-33-read-replicas-and-replication-lag-metadata.md) | Title, subtitle, SEO title (47 kar), SEO description (150 kar), slug, 5 tags, alt‑text cover, social metadata |
| 8 | LinkedIn Caption (plain text) | [linkedin/lab-33-read-replicas-and-replication-lag-caption.md](linkedin/lab-33-read-replicas-and-replication-lag-caption.md) | 235 kata (target 180‑250), awalan persis *Lab 33 kali ini membahas …*, mental model Inggris, tanpa markdown bold |

## Catatan Verifikasi & Batas Implementasi
- **Commit SHA**: `cfba3b76824f0ee03b5fa526462d1b72858144e7`
- **Status Test**: 7/7 test cases PASS (`go test -v -race ./...`)
- **Batas Implementasi**: Simulasi kluster in-memory key-value dengan pencatatan Log Sequence Number (LSN) atomik; tidak menggunakan database fisik multi-node atau disk persistence; latensi numerik adalah nilai demonstratif terkontrol.

## Caption LinkedIn (Siap Salin)
Lab 33 kali ini membahas anomali replication lag pada database read replicas, kegagalan pola pemisahan baca-tulis konvensional, serta empat strategi perutean cerdas untuk mempertahankan konsistensi read-your-own-writes.

Masalah teknis terjadi saat aplikasi mengarahkan kueri baca ke replika sekunder untuk skala beban. Karena replikasi basis data berjalan asinkron di latar belakang, terdapat selisih waktu antara commit di primary dengan penerapan data di replika. Pengguna yang baru saja memperbarui profil atau checkout pesanan dapat disajikan data usang atau menerima pesan data tidak ditemukan saat kueri masuk ke replika yang tertinggal.

Lab ini membuktikan empat hal konkret:
• Stale Read Anomaly: Kueri naif ke replika ber-lag 500ms langsung memicu ErrNotFound sesaat setelah transaksi tulis committed di primary.
• Time-Based Sticky Routing: Penahanan kueri baca ke primary selama jendela 300ms berhasil melindungi sesi sebelum beralih ke replika.
• Causal Token Routing: Perutean token LSN menunggu secara efisien via sync.Cond hingga replika mencapai posisi log yang dipersyaratkan.
• Synchronous Freshness: Replikasi sinkron menjamin ketersediaan data seketika pada replika dengan penalti waktu penulisan sebesar 102ms.

Mental model inti: Replication lag is the distance between primary commit and replica apply. Match your routing strategy to the cost of reading stale state.

Satu batasan penting: Lab ini menggunakan simulasi kluster in-memory dengan urutan LSN atomik, bukan pengujian basis data fisik multi-region, serta berfokus murni pada logika perutean dan sinkronisasi sesi.

Pelajari implementasi lengkap dan unit test di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/33-read-replicas-and-replication-lag

#softwareengineering #golang #database #distributeddatabase #systemdesign #backend
