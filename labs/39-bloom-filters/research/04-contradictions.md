# Contradictions & Trade-Offs

## Material Contradictions
Tidak ditemukan kontradiksi material terkait kebenaran matematis dasar Bloom Filter. Semua sumber otoritatif (Bloom 1970, Kirsch-Mitzenmacher 2006, Google Bigtable, RocksDB, Cassandra) sepakat bahwa Bloom Filter tidak menghasilkan false negative dan formula false positive rate terbukti konsisten.

## Divergence & Architectural Trade-offs

### 1. Hash Independence vs Kirsch-Mitzenmacher Optimization
- **Source A (Teori Klasik Bloom 1970)**: Mensyaratkan $k$ fungsi hash independen yang seragam (uniform independent hash functions).
- **Source B (Kirsch & Mitzenmacher 2006)**: Membuktikan bahwa menghasilkan $k$ nilai hash dari dua fungsi hash via $g_i(x) = h_1(x) + i \cdot h_2(x) \pmod m$ memiliki false positive rate asimtotik yang sama secara praktis dengan $k$ hash independen.
- **Assessment**: Dalam rekayasa perangkat lunak modern (RocksDB, Guava, Redis), pendekatan Kirsch-Mitzenmacher adalah standar de facto karena menghemat siklus CPU tanpa menambah false positive yang terukur.

### 2. Standard Bloom Filter vs Block-based / Split Bloom Filter
- **Source A (Standard Bloom Filter)**: Menyebarkan $k$ bit di seluruh bit array ukuran $m$. Akibatnya, setiap query dapat memicu $k$ kali cache miss pada CPU L1/L2/L3 cache karena bit tersebar acak di memori.
- **Source B (RocksDB Block-based / FastFilter)**: Membatasi bit-bit untuk satu kunci hanya pada 1 blok cache-line (misal 64 byte / 512 bit).
- **Assessment**: Block-based filter meningkatkan kecepatan lookup CPU (1 cache miss per query vs $k$ cache misses), namun mengalami sedikit peningkatan false positive rate empiris akibat variansi load factor antar blok. Untuk sistem berkapasitas besar di mana throughput query menjadi bottleneck, varian blocked lebih disukai.

### 3. Bloom Filter vs Cuckoo Filter
- **Source A (Bloom Filter Advocates)**: Bloom Filter lebih sederhana diimplementasikan, alokasi memori linear bit-level, performa insert sangat cepat dan deterministik tanpa evictions.
- **Source B (Fan et al. 2014, Cuckoo Filter)**: Cuckoo Filter mendukung operasi `delete`, memiliki lookup locality lebih baik, dan menghemat memori pada target error rate rendah ($p < 3\%$). Namun, penambahan elemen (insert) dapat gagal jika tabel penuh (*hash table max capacity / cascade evictions*).
- **Assessment**: Jika sistem membutuhkan operasi `delete` dinamis di RAM, Cuckoo Filter atau Counting Bloom Filter lebih cocok. Untuk immutable/append-only storage (seperti SSTable pada LSM-Tree), standard Bloom filter tetap menjadi pilihan utama industri.
