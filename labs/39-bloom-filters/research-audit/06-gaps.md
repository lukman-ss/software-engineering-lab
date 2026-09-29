# Research Gap Analysis: Bloom Filters Research

## Gap 1
Type: WEAK_SOURCE  
Severity: LOW  
Location: `03-evidence.md` (Evidence 4)  
Problem: Klaim pencegahan cache penetration disokong oleh RocksDB Wiki dan DDIA (Kleppmann). RocksDB Wiki adalah sumber relevan untuk disk I/O pruning di LSM-Tree, namun bukan sumber primer untuk pola cache penetration dalam konteks cache layer (Redis/Memcached). Referensi ke "Redis Bloom module docs" dalam bagian "Corroborated By" tidak masuk daftar sumber formal (`02-sources.md`).  
Required Revision: Tambahkan sumber resmi untuk pola cache penetration (misalnya Redis Stack atau AWS ElastiCache documentation, atau literatur system design yang lebih eksplisit membahas cache penetration protection pattern).  
Can Be Approved Without Fix: YES (Klaim secara konsep akurat dan dapat dipahami dari sumber-sumber yang tersedia)

---

## Gap 2
Type: UNVERIFIED_CLAIM  
Severity: MEDIUM  
Location: `05-report.md` (Executive Summary)  
Problem: Klaim bahwa "overhead pointer, bucket table, dan string key object pada Go map atau Java HashSet umumnya mengonsumsi rata-rata 48-96 byte per entri" digunakan sebagai acuan perbandingan efisiensi memori Bloom Filter. Klaim ini tidak didukung oleh sumber formal manapun dalam `02-sources.md`. Angka 48-96 byte per entri belum dikaitkan ke sumber primer atau benchmark yang dapat direproduksi.  
Required Revision: Tambahkan sumber yang mengukur/mendokumentasikan overhead memori Go map atau Java HashSet secara empiris.  
Can Be Approved Without Fix: YES (Klaim ini adalah ilustrasi komparatif, bukan hasil utama riset)

---

## Gap 3
Type: MISSING_SOURCE  
Severity: LOW  
Location: `05-report.md` (Finding 2) & `03-evidence.md` (Evidence 3)  
Problem: Corroboration menyebut "Implementasi Google Guava BloomFilter.java" sebagai bukti adopsi Double Hashing, namun tidak ada URL atau referensi formal ke kode sumber Guava yang dapat diakses dan diverifikasi dalam daftar sumber.  
Required Revision: Tambahkan URL ke source code Google Guava `BloomFilter.java` (misalnya GitHub Guava repository) sebagai referensi kode yang dapat diverifikasi.  
Can Be Approved Without Fix: YES

---

## Gap 4
Type: UNVERIFIED_CLAIM  
Severity: MEDIUM  
Location: `06-open-questions.md` (Section 2.1)  
Problem: Riset secara terbuka mengakui bahwa klaim "MurmurHash3 dan xxHash lebih baik untuk Bloom Filter dibanding MD5/SHA" belum dikonfirmasi dari benchmark primer terpercaya. Ini adalah kelemahan yang sudah diidentifikasi oleh riset itu sendiri.  
Required Revision: Tidak memerlukan perbaikan segera mengingat sudah dicatat sebagai open question, namun direkomendasikan validasi dari SMHasher suite atau paper benchmark yang terpublikasi.  
Can Be Approved Without Fix: YES (Sebagai labelling open question, ini adalah integritas ilmiah yang baik)

---

## Gap 5
Type: MISSING_CASE  
Severity: LOW  
Location: `05-report.md` (Limitations)  
Problem: Diskusi limitasi Bloom Filter tidak membahas implikasi hash seed atau non-determinism antar restart proses, yang relevan untuk penggunaan pada distributed system (misalnya saat SSTable dibuat oleh proses berbeda dan filter perlu dikonsultasi lintas node).  
Required Revision: Pertimbangkan menambahkan catatan tentang pentingnya seed deterministik pada hash function dalam konteks persisted Bloom Filter.  
Can Be Approved Without Fix: YES

---

## Summary
- 5 gap diidentifikasi, semua dengan severity LOW–MEDIUM.
- Tidak ada gap yang mengancam keakuratan inti riset.
- Seluruh gap dapat disetujui tanpa perbaikan wajib sebelum penggunaan sebagai dasar lab/artikel.
