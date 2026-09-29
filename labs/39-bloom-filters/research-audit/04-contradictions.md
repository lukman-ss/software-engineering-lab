# 04 - Contradictions Audit: Bloom Filters Research

## Summary
Tidak ditemukan kontradiksi material internal maupun eksternal yang merusak keabsahan riset Bloom Filters.

---

## Analysis of Investigated Trade-offs

### 1. Hash Independence Assumption vs Double Hashing Approximation
- **Statement A**: Teori orisinal Burton H. Bloom (1970) mengasumsikan $k$ fungsi hash acak independen dan berdistribusi seragam.
- **Statement B**: Praktik rekayasa (Kirsch & Mitzenmacher 2006) menggunakan dua fungsi hash $h_1(x) + i \cdot h_2(x) \pmod m$ untuk menghemat CPU.
- **Type**: SOURCE_REFINEMENT (Bukan kontradiksi fatal, melainkan optimasi terbukti).
- **Assessment**: Teorema Kirsch-Mitzenmacher secara formal membuktikan bahwa laju error asimtotik konvergen ke bound yang sama dengan $k$ hash independen.

### 2. Cache-Line Misses vs Bit Uniformity (Standard vs Blocked Filter)
- **Statement A**: Standard Bloom filter menyebarkan $k$ bit secara acak di seluruh array $m$, memicu hingga $k$ cache miss per lookup.
- **Statement B**: Block-based Bloom Filter (RocksDB) melokalisasi $k$ bit dalam 1 cache line 64-byte untuk lookup berkecepatan tinggi, dengan konsekuensi variasi load factor lokal sedikit menaikkan false positive rate.
- **Type**: ARCHITECTURAL_TRADEOFF.
- **Assessment**: Riset mendokumentasikan trade-off ini secara transparan di `04-contradictions.md` dan `06-open-questions.md`.

### 3. Static Sizing vs Dynamic Resizing
- **Statement A**: Rumus Bloom filter mensyaratkan estimasi $n$ elemen di awal.
- **Statement B**: Jika $n$ terlampaui tanpa batas, false positive rate mendekati 100%.
- **Type**: LIMITATION_CONSISTENCY.
- **Assessment**: Riset secara eksplisit mencantumkan keterbatasan kapasitas statis pada bagian limitasi di `05-report.md:59` dan menyarankan Scalable Bloom Filter (SBF) pada `06-open-questions.md:9`.
