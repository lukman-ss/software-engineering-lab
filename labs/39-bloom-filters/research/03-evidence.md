# Evidence

## Evidence 1
Claim: Bloom filter menjamin tidak ada false negative (0% false negative rate), namun dapat menghasilkan false positive ($p > 0$).
Evidence: Dalam skema Bloom filter standar, ketika sebuah elemen dimasukkan, seluruh $k$ bit pada posisi hash diatur menjadi 1. Jika suatu elemen sebelumnya telah dimasukkan, pengecekan ke-$k$ posisi bit tersebut selalu menghasilkan 1. Oleh karena itu, jika setidaknya ada satu bit bernilai 0 saat query, elemen tersebut dijamin 100% tidak ada dalam himpunan.
Source: Burton H. Bloom (1970), "Space/Time Trade-offs in Hash Coding with Allowable Errors", Communications of the ACM.
URL: https://dl.acm.org/doi/10.1145/362686.362692
Confidence: HIGH
Corroborated By: Kirsch & Mitzenmacher (2006); RocksDB Wiki; Apache Cassandra Docs.
Notes: Fundamental axiom of Bloom Filters.

## Evidence 2
Claim: Hubungan matematis ukuran bit array optimal ($m$), jumlah elemen ($n$), target false positive ($p$), dan jumlah hash ($k$) adalah:
$m = - \frac{n \ln p}{(\ln 2)^2} \approx -1.4427 \cdot n \log_2 p$
$k = \frac{m}{n} \ln 2 \approx 0.6931 \cdot \frac{m}{n}$
Untuk target $p = 0.01$ (1%), rasio $m/n \approx 9.585$ bit per elemen, dan $k \approx 7$.
Evidence: Derivasi probabilitas false positive $p \approx (1 - e^{-kn/m})^k$. Nilai minimum $p$ tercapai ketika $k = (m/n) \ln 2$. Substitusi $k$ optimal ke dalam persamaan $p$ menghasilkan $p = 2^{-k} = (1/2)^{(m/n)\ln 2}$, sehingga $m = - \frac{n \ln p}{(\ln 2)^2}$.
Source: Burton H. Bloom (1970); Kirsch & Mitzenmacher (2006).
URL: https://dl.acm.org/doi/10.1145/362686.362692
Confidence: HIGH
Corroborated By: RocksDB Wiki; Broder & Mitzenmacher (2004) "Network Applications of Bloom Filters".
Notes: Rumus baku industri untuk kalkulasi alokasi memori filter.

## Evidence 3
Claim: Dua fungsi hash independen ($h_1(x)$ dan $h_2(x)$) dapat menghasilkan $k$ fungsi hash efektif tanpa meningkatkan laju false positive asimtotik melalui formula $g_i(x) = h_1(x) + i \cdot h_2(x) \pmod m$ untuk $i = 0, \dots, k-1$.
Evidence: Teorema Kirsch-Mitzenmacher membuktikan bahwa kombinasi linier dua fungsi hash seragam independen menghasilkan distribusi keanggotaan bit yang ekivalen secara asimtotik dengan penggunaan $k$ fungsi hash acak independen penuh, memangkas beban komputasi CPU secara signifikan.
Source: Adam Kirsch and Michael Mitzenmacher (2006), "Less Hashing, Same Performance: Building a Better Bloom Filter", ESA 2006.
URL: https://www.eecs.harvard.edu/~michaelm/postscripts/esa2006.pdf
Confidence: HIGH
Corroborated By: Implementasi Google Guava `BloomFilter.java`, RocksDB `DynamicBloom`.
Notes: Sering diimplementasikan dengan mengekstrak dua 64-bit integer dari single 128-bit hash seperti MurmurHash3 atau xxHash.

## Evidence 4
Claim: Bloom filter mencegah masalah "Cache Penetration" dengan memfilter query untuk kunci nonexistent sebelum query menyentuh cache atau database storage.
Evidence: Pada cache penetration, bot/penyerang meminta ID acak nonexistent yang menyebabkan cache miss permanen dan memaksa disk query berulang. Dengan menempatkan Bloom Filter di depan lookup layer, request untuk ID yang tidak ada langsung diidentifikasi dengan kepastian 100% (kecuali probabilitas kecil false positive $p$), memotong hingga $(1 - p)$ atau ~99% beban query disk/database.
Source: RocksDB Documentation / System Design Literature (Designing Data-Intensive Applications, Martin Kleppmann, Bab 3).
URL: https://github.com/facebook/rocksdb/wiki/RocksDB-Bloom-Filter
Confidence: HIGH
Corroborated By: Google Bigtable (OSDI 2006); Redis Bloom module docs.
Notes: Pola standar pertahanan arsitektural database throughput tinggi.

## Evidence 5
Claim: Mesin penyimpanan berbasis Log-Structured Merge-Tree (LSM-Tree) seperti Bigtable, RocksDB, dan Cassandra menggunakan Bloom Filter pada setiap SSTable untuk menghindari pembacaan disk yang tidak perlu.
Evidence: Chang et al. (2006) menyatakan: "Bigtable allows clients to specify that Bloom filters should be created for SSTables in a particular locality group. A Bloom filter allows us to ask whether an SSTable might contain any data for a specified row/column pair... Drastically reduces the number of disk seeks required for read operations."
Source: Fay Chang et al. (2006), "Bigtable: A Distributed Storage System for Structured Data", Google Inc.
URL: https://static.googleusercontent.com/media/research.google.com/en//archive/bigtable-osdi06.pdf
Confidence: HIGH
Corroborated By: Apache Cassandra Architecture Docs; RocksDB Wiki.
Notes: Tanpa Bloom filter, point query pada LSM tree harus memeriksa setiap level SSTable secara sekuensial.

## Evidence 6
Claim: Standard Bloom Filter tidak mendukung operasi penghapusan (`delete`). Menghapus bit 1 menjadi 0 dapat menyebabkan false negative pada elemen lain.
Evidence: Karena multiple keys berbagi bit array yang sama melalui hash collisions, mengatur bit dari 1 ke 0 saat menghapus kunci A dapat merusak representasi kunci B yang juga memetakan bit yang sama. Untuk mendukung deletion, diperlukan struktur varian seperti Counting Bloom Filter (CBF) atau Cuckoo Filter.
Source: Fan et al. (2014), "Cuckoo Filter: Practically Better Than Bloom", ACM CoNEXT 2014.
URL: https://www.cs.cmu.edu/~dga/papers/cuckoo-conext2014.pdf
Confidence: HIGH
Corroborated By: Broder & Mitzenmacher (2004).
Notes: Penghapusan pada standard filter hanya bisa dilakukan dengan me-rebuild seluruh filter dari awal.
