# Sources

## Source 1
Title: Space/Time Trade-offs in Hash Coding with Allowable Errors
Publisher: Communications of the ACM (Vol. 13, Issue 7, pp. 422–426)
URL: https://dl.acm.org/doi/10.1145/362686.362692
Published: 1970-07-01
Accessed: 2026-09-29
Source Tier: Tier 1 (Foundational Academic Paper)
Relevance: Makalah seminal penemu Bloom Filter (Burton H. Bloom), mendefinisikan trade-off ruang-waktu dan properti "allowable errors" (false positive tanpa false negative).

## Source 2
Title: Less Hashing, Same Performance: Building a Better Bloom Filter
Publisher: Harvard University / European Symposium on Algorithms (ESA 2006, LNCS 4168, pp. 456–467)
URL: https://www.eecs.harvard.edu/~michaelm/postscripts/esa2006.pdf
Published: 2006-09-11
Accessed: 2026-09-29
Source Tier: Tier 1 (Academic Paper)
Relevance: Membuktikan teknik $g_i(x) = h_1(x) + i \cdot h_2(x) \pmod m$ cukup untuk mensimulasikan $k$ fungsi hash independen tanpa degradasi false positive rate asimtotik.

## Source 3
Title: Bigtable: A Distributed Storage System for Structured Data
Publisher: USENIX OSDI 2006 (Google, Inc.)
URL: https://static.googleusercontent.com/media/research.google.com/en//archive/bigtable-osdi06.pdf
Published: 2006-11-06
Accessed: 2026-09-29
Source Tier: Tier 1 (Industry Research Paper)
Relevance: Dokumentasi arsitektur orisinal penggunaan Bloom Filter pada SSTable untuk memangkas disk lookup bagi baris/kolom nonexistent.

## Source 4
Title: RocksDB Bloom Filter
Publisher: Meta Open Source / RocksDB GitHub Wiki
URL: https://github.com/facebook/rocksdb/wiki/RocksDB-Bloom-Filter
Published: 2023-01-15 (Updated)
Accessed: 2026-09-29
Source Tier: Tier 1 (Official Engine Documentation)
Relevance: Menjelaskan implementasi Bloom Filter modern (Block-based & Full Bloom Filter, Ribbon filter) untuk meminimalkan I/O pada LSM-Tree.

## Source 5
Title: Apache Cassandra Architecture: Bloom Filters
Publisher: Apache Software Foundation
URL: https://cassandra.apache.org/doc/latest/cassandra/operating/bloom_filters.html
Published: 2024-01-01 (Latest docs)
Accessed: 2026-09-29
Source Tier: Tier 1 (Official Database Documentation)
Relevance: Dokumentasi tuning `bloom_filter_fp_chance`, trade-off konsumsi RAM off-heap SSTable vs disk reads.

## Source 6
Title: Cuckoo Filter: Practically Better Than Bloom
Publisher: Carnegie Mellon University / ACM CoNEXT 2014
URL: https://www.cs.cmu.edu/~dga/papers/cuckoo-conext2014.pdf
Published: 2014-12-02
Accessed: 2026-09-29
Source Tier: Tier 1 (Academic Paper)
Relevance: Analisis komparatif varian modern pendukung operasi `delete` dan perbandingan efisiensi memori terhadap Bloom Filter.
