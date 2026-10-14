# Publication Deliverables — Lab 30: Leader Election with Fencing Tokens

## Daftar File Deliverable (8 Deliverables + Python Definition)

| No | Tipe Deliverable | Lokasi File | Deskripsi | Ukuran |
|---|---|---|---|---|
| 0 | Python Slide Definition (Master) | [content_def.py](content_def.py) | Master script Python berisi definisi 24 slide, konstanta warna, metadata, dan mental models (DILARANG DIHAPUS) | 39 KB |
| 1 | PowerPoint Carousel (editable) | [linkedin/lab-30-leader-election-carousel.pptx](linkedin/lab-30-leader-election-carousel.pptx) | Tepat 24 slide 16:9 landscape, editable vector shapes/text, speaker notes lengkap di setiap slide | 116 KB |
| 2 | LinkedIn Carousel PDF (selectable) | [linkedin/lab-30-leader-election-linkedin.pdf](linkedin/lab-30-leader-election-linkedin.pdf) | Tepat 24 halaman 16:9, selectable text, dark navy cover & closing, editorial light interior, rendered via Chrome headless | 746 KB |
| 3 | Dev.to Markdown Article (English) | [dev.to/lab-30-leader-election-devto.md](dev.to/lab-30-leader-election-devto.md) | Artikel teknis bahasa Inggris, valid YAML frontmatter, 2.760 kata (>1.500 kata), syntax-highlighted Go code, link repo & commit hash | 29 KB |
| 4 | Medium HTML Article (Indonesian) | [medium/lab-30-leader-election-medium.html](medium/lab-30-leader-election-medium.html) | Artikel teknis bahasa Indonesia, standalone HTML Georgia/Arial styling, 3.376 kata (>2.000 kata), lengkap 9 section wajib, tanpa JS | 32 KB |
| 5 | Cover Image (Indonesia) | [image-cover/lab-30-leader-election-cover-id.png](image-cover/lab-30-leader-election-cover-id.png) | 1920×1080 PNG retro CRT terminal monitor, subtitle Bahasa Indonesia, rendered via Chrome screenshot | 752 KB |
| 6 | Cover Image (English) | [image-cover/lab-30-leader-election-cover-en.png](image-cover/lab-30-leader-election-cover-en.png) | 1920×1080 PNG retro CRT terminal monitor, subtitle Bahasa Inggris, rendered via Chrome screenshot | 752 KB |
| 7 | Medium Metadata (Markdown) | [medium/lab-30-leader-election-metadata.md](medium/lab-30-leader-election-metadata.md) | Title, subtitle, SEO title (38 kar), SEO description (147 kar), slug, 5 tags, alt-text, social metadata | 1.5 KB |
| 8 | LinkedIn Caption (plain text) | [linkedin/lab-30-leader-election-caption.md](linkedin/lab-30-leader-election-caption.md) | Tepat 243 kata (180–250 kata), awalan persis *Lab 30 kali ini membahas ...*, English mental model, tanpa markdown bold | 1.9 KB |

---

## Catatan Verifikasi & Batas Implementasi

- **Commit SHA**: `e963675c67e7d518e19a07ad930dbdb5be73fff1`
- **Status Test**: 4/4 test suites PASS (`go test -v -race ./...` in 2.000s, 0 data races)
  * `TestCoordinatorLeaseAcquisitionAndRenewal` (PASS, 0.25s)
  * `TestFencedStorageRejectsStaleTokens` (PASS, 0.00s)
  * `TestLeaderElectionFailoverAndSplitBrainDefense` (PASS, 0.28s)
  * `TestConcurrentElectionRace` (PASS, 0.10s)
- **Batas Implementasi Utama**:
  1. Koordinator dijalankan secara in-memory dalam satu proses sistem operasi menggunakan `sync.Mutex`, bukan sebagai kluster konsensus multi-node terdistribusi (seperti Raft di etcd).
  2. Jeda eksekusi dimodelkan via penanda waktu monotonik `pauseUntil` pada candidate loop, bukan sinyal OS level `SIGSTOP` atau partisi antarmuka jaringan fisik.
  3. Nilai durasi waktu lab (TTL 200 ms, Renew 50 ms) dioptimalkan demi kecepatan eksekusi demo dan unit test, bukan angka rekomendasi untuk lingkungan produksi (di mana TTL umumnya 5–15 detik).
  4. Fencing token mewajibkan lapisan penyimpanan hilir mendukung conditional check-and-set atomik (`incomingToken > lastSeenToken`); proteksi tidak dapat ditegakkan jika backend hanya menyediakan API pembaruan buta tanpa versioning.

---

## Caption LinkedIn (Siap Salin)

Lab 30 kali ini membahas implementasi Leader Election dengan Fencing Tokens untuk mencegah split-brain dan mutasi zombie pada sistem terdistribusi Go.

Tantangan utama pada koordinasi terdistribusi adalah ilusi waktu sinkron. Ketika leader primer mengalami stop-the-world GC pause atau partisi jaringan melebihi batas sewa Time-To-Live, koordinator menganggap node tersebut mati dan melimpahkan mandat ke kandidat baru. Saat node lama terbangun tanpa sadar masa sewanya telah hangus, ia mengeksekusi operasi penulisan basi yang menimpa mutasi baru dan merusak integritas shared storage.

Lab ini membuktikan lima prinsip rekayasa sistem:
• Linearizable Coordinator: Koordinator thread-safe menerbitkan sewa berdurasi waktu dan menaikkan nomor revisi monotonik secara atomik.
• Bounded Failure Detection: Heartbeat renewal berkala mendeteksi kegagalan leader dan mengeksekusi failover otomatis ke node siaga.
• Strict Check-and-Set: Gerbang penyimpanan bersama menolak penulisan jika nomor token yang datang bernilai kurang dari atau sama dengan nomor token terakhir yang terlihat.
• Concurrency Safety: Lima kandidat bersaing simultan memperebutkan sewa dengan hasil deterministik tepat satu pemenang tanpa race condition.
• Monotonic Clocks: Perhitungan masa kadaluwarsa sewa memanfaatkan jam monotonik Go untuk menghindari lompatan jam dinding NTP.

Mental model inti: A lease grants time-bounded authority; a fencing token converts that temporal authority into a monotonic ordering that shared storage enforces, rejecting stale zombie writes.

Batasan penting: Lab ini menguji logika koordinasi via memori internal dan mutex Go, bukan kluster etcd fisik multi-node. Parameter waktu disederhanakan untuk kecepatan eksekusi demo.

Pelajari kode sumber dan pengujian di:
https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/30-leader-election

#softwareengineering #golang #distributedsystems #concurrency #systemdesign #backend
