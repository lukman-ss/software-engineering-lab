# Research Plan: Bloom Filters

## Research Topic
Bloom Filters — Struktur Data Probabilistik untuk Memangkas Disk I/O dan Mencegah Cache Penetration pada Backend Berkinerja Tinggi.

## Objective
Menginvestigasi dasar matematis, karakteristik operasional, formula optimalisasi ukuran ($m$) dan jumlah fungsi hash ($k$), strategi hashing efisien (Kirsch-Mitzenmacher), pencegahan cache penetration, penggunaan industri pada LSM-Tree (RocksDB, Cassandra), serta trade-off dan limitasi Bloom Filter.

## Research Questions
1. **Mathematical Foundations**: Bagaimana perumusan probabilitas false positive ($p$), penentuan ukuran bit array optimal ($m$), dan jumlah fungsi hash optimal ($k$) untuk $n$ elemen?
2. **Double Hashing Optimization**: Apakah implementasi $k$ fungsi hash independen memerlukan $k$ algoritma hash berbeda, atau cukup 2 fungsi hash (Kirsch-Mitzenmacher technique)?
3. **Cache Penetration Mitigation**: Bagaimana pola arsitektur Bloom Filter melindungi database dari serangan pencarian kunci nonexistent?
4. **LSM-Tree Disk I/O Pruning**: Bagaimana mesin penyimpanan LSM-Tree (RocksDB, Google Bigtable, Apache Cassandra) memanfaatkan Bloom Filter pada SSTable?
5. **Limitations & Variants**: Mengapa Bloom Filter standar tidak mendukung operasi `delete`, dan varian apa yang mengatasi limitasi tersebut (Counting Bloom Filter, Cuckoo Filter)?

## Search Strategy
1. Peninjauan paper primer: Burton H. Bloom (1970).
2. Peninjauan paper optimasi hash: Adam Kirsch & Michael Mitzenmacher (2006).
3. Peninjauan dokumentasi arsitektur industri: RocksDB (Meta), Apache Cassandra, Google Bigtable paper (Chang et al., 2006).
4. Peninjauan literatur algoritma hash non-kriptografis: MurmurHash3 (Austin Appleby), xxHash (Yann Collet), FNV.

## Expected Primary Sources
- Burton H. Bloom (1970), "Space/Time Trade-offs in Hash Coding with Allowable Errors", Communications of the ACM.
- Adam Kirsch and Michael Mitzenmacher (2006), "Less Hashing, Same Performance: Building a Better Bloom Filter", ESA 2006 / Harvard University.
- Fay Chang et al. (2006), "Bigtable: A Distributed Storage System for Structured Data", Google Inc., OSDI 2006.
- RocksDB Documentation / Wiki on RocksDB Bloom Filter implementation (Meta Open Source).
- Apache Cassandra Architecture Documentation on Bloom Filters.

## Risks / Unknowns
- Kinerja CPU overhead dari hashing berulang jika fungsi hash tidak optimal.
- Efek cache-line locality pada bit array berukuran besar (Block-based / Split Bloom Filter).
- Kebutuhan alokasi ulang / re-hashing jika estimasi $n$ terlampaui (Scalable Bloom Filter).
