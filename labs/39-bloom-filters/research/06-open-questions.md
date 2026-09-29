# Open Questions

## 1. Pertanyaan Belum Terjawab Sepenuhnya

### 1.1 Ribbon Filter (Successor RocksDB)
RocksDB v6.15+ memperkenalkan Ribbon Filter sebagai pengganti Block-based Bloom Filter pada LSM-Tree. Claim: lebih hemat memori pada error rate rendah. Perlu investigasi lebih lanjut tentang tradeoff CPU waktu build vs lookup dan apakah ini layak diimplementasikan dari scratch untuk laboratorium.
- **Sumber potensial**: Facebook Engineering Blog; Peter Dillinger & Martijn Moritz "Ribbon Filters" (2021).

### 1.2 Scalable Bloom Filter (SBF)
Ketika jumlah elemen aktual melebihi estimasi $n$ awal, perlu strategi ekspansi filter tanpa menyebabkan false positive rate melonjak. Scalable Bloom Filter (Almeida et al. 2007) menggunakan rangkaian filter dengan error rate geometrik. Implementasi praktisnya untuk sistem backend belum diinvestigasi mendalam.

### 1.3 Thread Safety & Concurrent Access
Penggunaan Bloom Filter di lingkungan multi-threaded (misalnya server Go / Java dengan goroutine/thread pool) memerlukan mekanisme sinkronisasi atau penggunaan atomic bit operations. Apakah operasi `insert` dan `lookup` pada bit array aman secara konkurensi bergantung pada implementasi bahasa pemrograman. Perlu investigasi khusus untuk bahasa target lab.

## 2. Bukti yang Masih Lemah

### 2.1 Perbandingan Empiris MurmurHash3 vs xxHash vs FNV untuk Bloom Filter
Klaim "MurmurHash3 dan xxHash lebih baik untuk Bloom Filter dibanding MD5/SHA karena kecepatan" banyak beredar di komunitas tetapi belum di-cross-check dengan benchmark primer terpercaya. Benchmark dari sumber terpercaya (smhasher suite) perlu diverifikasi.

### 2.2 Ukuran Filter Praktis di Cassandra Produksi
Dokumentasi Cassandra menyebutkan nilai `bloom_filter_fp_chance` default 0.1 (10% false positive) untuk LocalStrategy dan 0.01 (1% false positive) untuk NetworkTopologyStrategy, tetapi data konsumsi RAM off-heap aktual dalam deployment skala besar belum terverifikasi dari data observability produksi nyata.

## 3. Penelitian Lanjutan yang Disarankan

1. **Empirical Validation**: Implementasikan dan ukur false positive rate empiris vs teoritis pada $n = 100.000$ elemen dengan berbagai nilai $k$ untuk validasi formula optimal.
2. **Benchmark Hash Functions**: Benchmarking MurmurHash3 vs xxHash3 vs FNV-1a menggunakan `smhasher` dan mengukur distribusi bit untuk 100K sampel.
3. **Block-based vs Standard BF**: Ukur latency query antara Standard Bloom Filter dan Block-based Filter pada memori besar (>1 GB bit array) untuk membuktikan perbedaan cache-line miss secara empiris.
4. **Ribbon Filter Deep Dive**: Baca "Ribbon Filter: Practically Smaller Than Bloom and Xor" (Dillinger, Walzer, 2021) untuk evaluasi apakah varian ini layak dipresentasikan pada lab lanjutan.
