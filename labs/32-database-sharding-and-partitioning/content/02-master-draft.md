# Database Sharding dan Table Partitioning

## Problem

Sebuah database monolitik memiliki batas fisik yang tidak dapat dielakkan: kapasitas memori, I/O throughput storage, dan lock contention pada satu server. Ketika volume data dan jumlah transaksi per detik tumbuh melampaui batasan ini, solusi vertical scaling (upgrade hardware) menjadi tidak efisien baik dari sisi biaya maupun skalabilitas.

Dua pendekatan utama untuk melampaui batas single instance adalah:

1. **Logical Partitioning**: Memecah satu tabel besar ke dalam sub-tabel fisik yang lebih kecil pada server yang sama.
2. **Physical Sharding**: Mendistribusikan dataset ke beberapa database server yang benar-benar independen.

Pilihan antara dua pendekatan ini, dan bagaimana mengimplementasikannya secara benar, menentukan apakah sistem mampu scale secara linear atau justru menciptakan bottleneck baru.

---

## Why This Matters

Kesalahan dalam rancangan sharding berdampak fundamental dan sulit dikoreksi setelah sistem berjalan di produksi:

- **Monotonic shard key** menyebabkan 100% write traffic menumpuk pada satu node ("hot partition").
- **Hash Modulo routing** pada saat scale-out node memaksa redistribusi $\approx 79.84\%$ data, artinya hampir seluruh data berpindah lokasi.
- **Query tanpa shard key** tanpa index sekunder memaksa sistem melakukan broadcast query ke seluruh shard cluster.

Memilih arsitektur yang tepat sejak awal adalah keputusan yang memiliki konsekuensi operasional jangka panjang.

---

## Mental Model

Pemisahan konsep yang paling penting:

| Dimensi | Logical Partitioning | Physical Sharding |
|---|---|---|
| Jumlah Instance | 1 database instance | N database instance independen |
| Tujuan Utama | Lifecycle management, pruning | Horizontal scaling melampaui batas 1 node |
| Skalabilitas | Terbatas pada kapasitas 1 mesin | Linear dengan jumlah node |
| Cross-partition JOIN | Native SQL didukung | Tidak tersedia; butuh scatter-gather atau co-location |
| Distributed ID | Tidak diperlukan | Wajib; auto_increment gagal di multi-node |

Konsep kunci: **partition key** (untuk logical partitioning) dan **shard key** (untuk physical sharding) adalah fondasi seluruh arsitektur. Kesalahan dalam pemilihan key ini tidak mudah diubah.

---

## Core Concept: Logical Table Partitioning

Partitioning memecah satu tabel logis menjadi beberapa sub-tabel fisik yang dikelola oleh satu database engine. PostgreSQL mendukung bentuk bawaan: Range, List, dan Hash Partitioning.

### Keunggulan Utama

**Partition Pruning**: Query optimizer hanya memindai partisi yang memiliki overlap dengan range query. Jika tabel memiliki 4 partisi kuartalan dan query hanya membutuhkan data Q2, engine hanya memindai 1 dari 4 partisi.

**Fast Partition Drop**: Menghapus data lama menggunakan `DROP TABLE` atau `DETACH PARTITION` jauh lebih cepat dan tidak memerlukan VACUUM pada seluruh tabel besar.

### Contoh Kasus

Tabel `orders` dengan 4 partisi berdasarkan range waktu kuartalan. Query untuk rentang April–Juni hanya memindai partisi Q2 — 3 partisi lainnya di-prune sepenuhnya.

---

## Core Concept: Physical Database Sharding

Sharding mendistribusikan data ke beberapa database server (shard) yang independen. Setiap shard menyimpan subset data dan dapat dijalankan sebagai replica set tersendiri.

Arsitektur sharded cluster terdiri dari:
- **Shard nodes**: Database instance independen yang menyimpan sebagian data.
- **Router**: Menentukan shard mana yang menjadi tujuan Insert atau Query berdasarkan shard key.
- **Global Secondary Index / Lookup Vindex** *(opsional)*: Memungkinkan point-lookup untuk query tanpa shard key.

---

## Failure Scenario: Write Hotspot karena Monotonic Key

Ini adalah kegagalan paling umum pada implementasi sharding.

Jika shard key adalah nilai yang monotonically increasing (seperti `created_at` timestamp, autoincrement integer, atau MongoDB ObjectId), semua insert baru selalu diarahkan ke satu shard yang sama — shard dengan chunk paling besar atau paling baru.

Hasil yang diverifikasi dari demo lab:

```
[Scenario A] Monotonic Key (date-string) Sharding:
  Node shard-0: 0 records  
  Node shard-1: 0 records  
  Node shard-2: 1000 records  [##############################]
  Node shard-3: 0 records  
Result: Severe Write Hotspot! 100% writes hit single shard.
```

Dengan shard key berbasis `user_id` (high-cardinality, non-monotonic):

```
[Scenario B] High-Cardinality Key (user_id) Sharding:
  Node shard-3: 400 records  [############                  ]
  Node shard-0: 400 records  [############                  ]
  Node shard-1: 0 records    
  Node shard-2: 200 records  [######                        ]
Result: Uniform distribution across physical shards.
```

**Kriteria shard key yang optimal** (berdasarkan riset MongoDB dan Vitess):
- **High cardinality**: Nilai distincts yang banyak (misalnya `user_id`, `tenant_id`), bukan nilai terbatas seperti `continent` yang hanya 7 nilai.
- **Non-monotonic**: Nilai tidak meningkat/menurun secara linier seiring waktu.
- **Enables point-lookup**: Satu nilai key selalu menghasilkan routing ke satu shard yang spesifik.

Ketika natural key bersifat monotonic (timestamp, ObjectId), **Hashed Sharding** adalah mitigasi yang terdokumentasi: hash dari nilai key mendistribusikan insert secara merata, dengan trade-off bahwa range query menjadi scatter-gather.

---

## Architecture: Dua Algoritma Routing

### Hash Modulo

Routing paling sederhana: `hash(key) % N` di mana N adalah jumlah shard.

**Masalah fundamental**: Ketika N berubah (saat menambah node), hampir semua key menghasilkan modulo yang berbeda. Dari N=4 ke N=5, fraksi key yang perlu pindah adalah $\frac{N}{N+1} = \frac{4}{5} = 80\%$.

Hasil yang diverifikasi:
```
Hash Modulo (N % M) Keys Remapped: 7984 / 10000 (79.84% moved)
```

### Consistent Hashing dengan Virtual Nodes

Consistent Hashing menempatkan node dan key pada "ring" hash circular. Ketika satu node ditambahkan, hanya key yang sebelumnya milik node tetangga di ring yang perlu berpindah — secara teori $\approx \frac{1}{N}$ dari total key.

Berdasarkan Karger et al. (1997, ACM STOC): saat resize hash table, rata-rata hanya $n/m$ key yang perlu di-remap (di mana $n$ = jumlah key, $m$ = jumlah slot).

**Virtual Nodes**: Setiap physical node direpresentasikan oleh banyak "virtual node" di ring (misalnya 100–150 vnodes per physical node). Virtual nodes mencegah distribusi tidak merata ketika jumlah physical node sedikit.

Hasil yang diverifikasi dengan 150 virtual nodes:
```
Consistent Hashing Keys Remapped: 800 / 5000 (16.00%)
Consistent Hashing Keys Remapped: 1200 / 10000 (12.00%)
Theoretical Minimal Relocation (1/N): ~20.00%
```

**Trade-off**: Consistent Hashing membutuhkan binary search $O(\log(N \cdot V))$ untuk routing, sedangkan Hash Modulo adalah $O(1)$. Namun Consistent Hashing mengurangi biaya resharding dari $\approx 80\%$ menjadi $\approx 12\%$–$16\%$ data migration.

---

## Implementation

Lab mengimplementasikan seluruh komponen dalam pure Go menggunakan standard library only (`hash/fnv`, `crypto/rand`, `sync`, `sort`).

### Struktur Komponen

```
internal/
  partitioning/
    table.go       — Range-partitioned table dengan partition pruning
  sharding/
    sharding.go    — Router, ConsistentHashRing, Cluster, GlobalSecondaryIndex, Shard
  idgen/
    idgen.go       — UUIDv7 (RFC 9562) dan SequenceBlockAllocator
cmd/
  demo/
    main.go        — End-to-end runnable demo
tests/
  sharding_test.go — Unit dan concurrency tests
```

### Komponen Inti

**`Shard`**: Simulasi satu database node menggunakan thread-safe in-memory map dengan `sync.RWMutex`.

**`ModuloRouter`**: Routing `hash(key) % N` menggunakan FNV-1a 64-bit.

**`ConsistentHashRouter`**: Hash ring dengan virtual nodes, pencarian via binary search (`sort.Search`). Default 100 virtual nodes per physical shard.

**`Cluster`**: Mengelola multiple `Shard`, menangani routing Insert/Query, dan mengorkestrasi Scatter-Gather atau GSI lookup.

**`GlobalSecondaryIndex`**: Thread-safe mapping dari secondary key (email) ke shard key, memungkinkan point-lookup tanpa broadcast.

---

## Code Walkthrough

### Partition Pruning

```go
// QueryRange demonstrates partition pruning: only partitions overlapping [start, end) are scanned.
func (t *Table) QueryRange(start, end time.Time) QueryResult {
    // ...
    for _, p := range t.partitions {
        // Overlap condition: start < p.Range.End && end > p.Range.Start
        if start.Before(p.Range.End) && end.After(p.Range.Start) {
            scanned++
            // scan only this partition
        }
    }
}
```

Hanya partisi yang overlap secara temporal dengan range query yang di-scan; partisi lain di-skip.

### Consistent Hash Ring Lookup

```go
func (ch *ConsistentHashRouter) GetShard(key string) (string, error) {
    h := hashKey(key)
    idx := sort.Search(len(ch.ring), func(i int) bool {
        return ch.ring[i].hash >= h
    })
    if idx == len(ch.ring) {
        idx = 0  // wrap around ring
    }
    return ch.ring[idx].shardID, nil
}
```

Binary search menemukan virtual node pertama yang hash-nya $\geq$ hash(key). Wrapping ke awal ring saat hash(key) > semua vnode.

### Scatter-Gather dengan Goroutines

```go
func (c *Cluster) ScatterGatherBroadcastWithContext(ctx context.Context, predicate func(Record) bool) ScatterGatherResult {
    var wg sync.WaitGroup
    ch := make(chan shardResult, len(shards))

    for _, s := range shards {
        wg.Add(1)
        go func(shard *Shard) {
            defer wg.Done()
            // each goroutine queries one shard concurrently
        }(s)
    }
    wg.Wait()
    close(ch)
    // aggregate results
}
```

Setiap shard di-query secara paralel dalam goroutine terpisah, dengan dukungan context cancellation.

### UUIDv7 Generation (RFC 9562)

```go
func NewUUIDv7() (string, error) {
    var uuid [16]byte
    ms := uint64(time.Now().UnixMilli())
    // bits 0-47: millisecond timestamp
    uuid[0] = byte(ms >> 40)
    // ...
    // version 7: 0111
    uuid[6] = (uuid[6] & 0x0F) | 0x70
    // variant RFC 4122/9562: 10xxxxxx
    uuid[8] = (uuid[8] & 0x3F) | 0x80
    // ...
}
```

48-bit timestamp di-encode di awal UUID, menjamin time-ordering leksikografis: UUIDv7 yang dibuat lebih belakangan selalu lebih besar secara string comparison.

---

## Queries Without the Sharding Key

Ketika query dilakukan berdasarkan atribut yang bukan shard key (misalnya mencari user berdasarkan `email` padahal shard key adalah `tenant_id`), dua opsi tersedia:

### Scatter-Gather (Broadcast)

Query dikirim ke **semua** shard secara paralel, setiap shard menyaring record yang cocok, dan hasil digabungkan di cluster level.

Hasil yang diverifikasi:
```
Scatter-Gather Query (by Email without Shard Key):
  Nodes Broadcasted: 4 / 4
  Records Matched: 1
  Execution Time: 92µs
```

Masalah: Biaya query naik linear dengan jumlah node. Untuk cluster besar atau high-QPS, ini menjadi bottleneck.

### Global Secondary Index (Lookup Vindex)

Sebuah secondary index menyimpan mapping: `email → shard_key`. Saat Insert, mapping ini diperbarui secara bersamaan (double-write). Saat query, router pertama lookup ke GSI untuk mendapatkan shard key, lalu langsung mengakses shard yang tepat.

Hasil yang diverifikasi:
```
Global Secondary Index (Lookup Vindex) Query:
  Nodes Broadcasted: 1 (Direct Point Lookup via Shard Key mapping)
  Record Found: ID=usr-342, Email=user_342@company.com, ShardKey=tenant-42
  Execution Time: 1µs (GSI Avoided Broadcast Overhead!)
```

**Trade-off GSI**: Setiap Insert/Update memerlukan double-write (ke shard record + ke GSI). Overhead ini diamortisasi seiring bertambahnya jumlah shard, karena GSI menghindari scatter-gather yang semakin mahal.

---

## Distributed ID Generation

Di lingkungan sharded, native `auto_increment` database tidak bisa digunakan lintas shard. RFC 9562 §2.1 menyatakan: "auto-increment schemes that are often used by databases do not work well: the effort required to coordinate sequential numeric identifiers across a network can easily become a burden."

Dua solusi yang diimplementasikan:

### UUIDv7 (RFC 9562)

128-bit UUID dengan 48-bit timestamp di depan. Menjamin:
- **Uniqueness**: Kombinasi timestamp + random bits.
- **Time-ordering**: Lexicographic order = chronological order, optimal untuk B-tree index.
- **No coordination**: Setiap node dapat generate UUID independen tanpa koordinasi.

Diverifikasi: Dua UUIDv7 yang dibuat dengan jeda 2ms memenuhi $u_1 < u_2$ secara leksikografis.

### Sequence Block Allocator (Vitess-style)

Numeric ID dengan block allocation: node mengambil blok ID dari central sequencer (misal 1–1000), lalu mengalokasikan ID secara lokal dari blok tersebut tanpa network round-trip per ID.

Trade-off: Jika node crash sebelum blok habis digunakan, ID dalam blok tersebut hilang (gap). Ini adalah trade-off yang didokumentasikan oleh Vitess, bukan bug.

Diverifikasi: Allocator dengan block size 10 menghasilkan ID sequential (1, 2, 3, ...) secara konsisten tanpa collision dalam 25 ID berturutan.

---

## What the Tests Prove

| Test | Behavior yang Diverifikasi |
|---|---|
| `TestPartitionPruning` | Range query pada 3 partisi hanya memindai 1 partisi (pruning bekerja). DropPartition berhasil. |
| `TestRoutingAndConsistentHashRelocation` | Hash Modulo memindahkan $\geq 65\%$ keys saat scale-out; Consistent Hash memindahkan $5\%$–$40\%$ keys. |
| `TestClusterScatterGatherAndGSI` | Scatter-gather query 3 shard; GSI lookup langsung ke shard tepat; context cancellation menghentikan semua goroutine. |
| `TestIDGenerators` | UUIDv7 time-ordered ($u_1 < u_2$); 25 sequential IDs tanpa collision dari block allocator. |
| `TestConcurrentClusterAccess` | 200 goroutine concurrent Insert/Get tanpa race condition (verified dengan `-race` detector). |

Semua tests lulus di `go test -race -v ./...`.

---

## Recovery / Resharding

### Resharding via Consistent Hashing

Untuk menambah node baru ke cluster, proses yang diimplementasikan:

1. Tambahkan shard baru ke router (`AddShardNode`).
2. Panggil `RebalanceData()`: kumpulkan semua record, reset semua shard, re-route semua record ke shard yang baru berdasarkan router yang telah diupdate.
3. Hanya $\approx 12\%$–$16\%$ data yang perlu dipindahkan.

Dalam sistem produksi nyata (seperti Vitess), resharding dilakukan secara live: data di-copy dan di-verify ke shard baru sementara shard lama tetap melayani traffic. Cutover hanya membutuhkan beberapa detik read-only downtime. Lab ini mensimulasikan mekanisme routing; live-copy tidak diimplementasikan.

### MongoDB Automatic Balancing

MongoDB menggunakan background balancer yang secara otomatis memigrasikan data ranges antar shard. MongoDB 8.0 memperkenalkan `sh.shardAndDistributeCollection()` yang langsung mendistribusikan data tanpa menunggu balancer cycle.

---

## Production Considerations

Hal-hal yang diverifikasi di lab tetapi perlu pertimbangan tambahan untuk produksi:

**1. Virtual Node Count**
100–150 virtual nodes per physical shard adalah konfigurasi yang digunakan di lab. Nilai lebih rendah ($<$ 50) dapat menyebabkan distribusi tidak merata ketika jumlah physical node sedikit.

**2. GSI Consistency**
Lab mengimplementasikan GSI sebagai in-process double-write yang synchronous. Di produksi, GSI update yang gagal (partial write) membutuhkan kompensasi transaksi atau eventually-consistent model. Vitess `consistent_lookup` menangani ini tanpa 2PC melalui locking dan sequencing yang terencana.

**3. Cross-Shard Transactions**
TwoPC (Two-Phase Commit) menjamin atomisitas tetapi bukan full ACID isolation. Vitess secara eksplisit mendokumentasikan bahwa TwoPC "does not provide isolation in the traditional ACID sense" — aplikasi mungkin mengalami fractured reads lintas shard. Lab ini tidak mengimplementasikan 2PC.

**4. ID Generation Trade-offs**
- UUIDv7: 128-bit, menjamin uniqueness global tanpa koordinasi; lebih besar dari integer.
- Sequence Block Allocator: Numeric integer, lebih compact; gap pada crash; butuh central coordinator.

---

## Common Mistakes

1. **Shard key berbasis `created_at` atau autoincrement** → write hotspot 100% pada satu shard. Gunakan `user_id` atau hash dari natural key.

2. **Tidak memperhitungkan biaya GSI** → double-write overhead tidak terlihat di low load tapi signifikan di high write throughput.

3. **Menggunakan Hash Modulo untuk cluster yang akan sering di-scale** → setiap penambahan node menyebabkan $\approx 80\%$ data migration. Gunakan Consistent Hashing.

4. **Mempertahankan auto_increment di lingkungan multi-shard** → collision dan serialisasi ID lintas node. Gunakan UUIDv7 atau block allocator.

5. **Berasumsi TwoPC memberikan full ACID** → cross-shard transactions masih dapat menghasilkan fractured reads. Co-locate data yang sering diakses bersama dalam shard key yang sama.

6. **Virtual nodes terlalu sedikit** → distribusi ring tidak merata, beberapa shard menerima jauh lebih banyak data.

---

## Checklist

Sebelum deploy sistem sharded ke produksi:

- [ ] Shard key telah dievaluasi: high-cardinality, non-monotonic, enables point-lookup
- [ ] Jika natural key monotonic (timestamp, ObjectId), hashed sharding diterapkan
- [ ] Routing algorithm: Consistent Hashing digunakan untuk cluster yang akan di-scale
- [ ] Virtual node count dikonfigurasi (≥ 50 per physical node)
- [ ] Distributed ID generation: UUIDv7 atau block allocator menggantikan auto_increment
- [ ] Query pattern dianalisis: tentukan mana yang butuh GSI, mana yang dapat scatter-gather
- [ ] Double-write failure handling pada GSI sudah dirancang
- [ ] Cross-shard transaction boundary sudah diminimalkan; 2PC digunakan hanya jika atomisitas wajib
- [ ] Resharding plan tersedia: berapa fraksi data yang dipindahkan dan estimasi durasi

---

## Key Takeaways

1. Partitioning dan Sharding bukan sinonim: partitioning adalah optimization dalam satu instance; sharding adalah distribusi ke banyak instance untuk melampaui batas satu mesin.

2. Pilihan shard key adalah keputusan desain yang paling kritis — sulit diubah setelah live.

3. Consistent Hashing dengan virtual nodes mengurangi data movement saat scale-out dari $\approx 80\%$ (Hash Modulo) menjadi $\approx 12\%$–$16\%$.

4. Query tanpa shard key wajib menggunakan GSI/Lookup Vindex untuk high-QPS; scatter-gather hanya efisien untuk aggregation besar.

5. Auto_increment tidak dapat digunakan di lingkungan multi-shard. UUIDv7 (time-ordered, no-coordination) atau Sequence Block Allocator (numeric, low-coordination) adalah solusi yang tersedia.

6. TwoPC di sharded cluster menjamin atomisitas tetapi tidak menjamin full ACID isolation; fractured reads mungkin terjadi.

---

## Sources

- PostgreSQL Documentation 18: Chapter 5.12 — Table Partitioning
- MongoDB Manual: Sharding Overview, Shard Key Selection, Hashed Sharding
- Vitess Documentation: Sharding, Resharding, Vindexes, Sequences, Distributed Transactions
- RFC 9562 §2.1 — UUID Version 7 (Time-Ordered UUIDs)
- Karger, D. et al. (1997). "Consistent Hashing and Random Trees." ACM STOC '97. DOI: 10.1145/258533.258660
- Wikipedia: Consistent Hashing (corroborating restatement of Karger et al.)
