#!/usr/bin/env python3
"""
Publication Builder for Software Engineering Lab 32: Database Sharding & Partitioning.
Generates 24-slide LinkedIn carousel covering Single-Node Logical Table Partitioning vs Physical Sharding,
Write Hotspots with Monotonic vs High-Cardinality Keys, Consistent Hashing vs Modulo Resharding,
Scatter-Gather vs Global Secondary Index (Lookup Vindex), and Distributed Unique ID Generation (UUIDv7 & Sequences).
"""
import os

OUTPUT_DIR = "labs/32-database-sharding-and-partitioning/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/32-database-sharding-and-partitioning"
COMMIT_SHA = "04c8c8a4a7923bcdf081f946284fca58c5a38bea"
LAB_NO = "32"
SLUG = "database-sharding-and-partitioning"
TITLE_ID = "Database Sharding & Table Partitioning"
SUBTITLE_ID = "Logical Pruning vs Physical Sharding, Consistent Hashing, dan Eliminasi Scatter-Gather"
TITLE_EN = "Database Sharding & Table Partitioning"
SUBTITLE_EN = "Logical Pruning vs Physical Sharding, Consistent Hashing, and Scatter-Gather Elimination"
AUTHOR = "Lukman / lukman-ss"
COVER_IMAGE_PATH = "labs/32-database-sharding-and-partitioning/result/image-cover/cover_art.jpg"

DARK_NAVY = "0b1b2d"
TEAL = "087f76"
CREAM = "f4f1e8"
WHITE = "ffffff"
SLATE = "52606b"
ORANGE = "d97706"

SLIDES = []

def add_slide(categories, title, subtitle, content_items, mental_model, source_path="", slide_num="", narrative="", narrative_en="", layout=""):
    SLIDES.append({
        "categories": categories,
        "title": title,
        "subtitle": subtitle,
        "content_items": content_items,
        "mental_model": mental_model,
        "source": source_path,
        "slide_num": slide_num,
        "narrative": narrative,
        "narrative_en": narrative_en,
        "layout": layout
    })

# =====================================================================
# 24 SLIDE DEFINITIONS
# =====================================================================

# Slide 1: Cover
add_slide(
    categories=["COVER"],
    title=TITLE_ID,
    subtitle=SUBTITLE_ID,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "Visual Hero: Skalabilitas Data Terdistribusi & Partisi Database",
            "Mekanisme Inti: Partition Pruning, Consistent Hash Ring, Global Secondary Index, UUIDv7",
            f"Penulis: {AUTHOR}",
        ]}
    ],
    mental_model="Logical partitioning divides tables within one instance; physical sharding distributes data across independent network nodes.",
    source_path=REPO_URL,
    slide_num="01 / 24",
    narrative="Skalabilitas database relasional merupakan tantangan paling kritis ketika volume data melampaui kapasitas I/O dan memori satu server. Kesalahan arsitektural yang paling sering terjadi adalah mencampuradukkan partisi logis satu mesin dengan sharding fisik terdistribusi, memilih sharding key monotonik yang memicu write hotspot, atau membiarkan query non-shard key membanjiri seluruh node via scatter-gather broadcast. Lab ini membuktikan implementasi murni Go untuk logical table partitioning dengan query pruning, router consistent hashing ring dengan virtual nodes untuk migrasi data minimal saat scale-out, eliminasi broadcast via Global Secondary Index, serta pembangkitan ID unik terdistribusi berbasis RFC 9562 UUIDv7 dan sequence block allocation.",
    narrative_en="Relational database scalability represents the ultimate bottleneck when data volume exceeds the single-node ceiling for storage I/O and RAM. The most prevalent architectural pitfalls include confusing single-engine table partitioning with multi-node physical sharding, selecting monotonic keys that induce severe write hotspots, or allowing un-sharded queries to exhaust cluster resources through naive scatter-gather broadcasts. This lab demonstrates a pure Go implementation of range-based logical table partitioning with query pruning, consistent hashing ring routers with virtual nodes for minimal resharding shuffle, point-lookup routing via Global Secondary Indexes, and distributed unique ID generation via RFC 9562 UUIDv7 and sequence block allocation.",
    layout="cover"
)

# Slide 2: Problem nyata
add_slide(
    categories=["MOTIVASI", "SCALABILITY BOTTLENECK"],
    title="Masalah Nyata: Batasan Skala Database Monolitik",
    subtitle="Ketika beban write dan volume data menembus batas komputasi satu mesin database.",
    content_items=[
        {"title": "SINGLE-INSTANCE I/O CEILING", "items": [
            "Ukuran working set melampaui buffer pool RAM, memicu lonjakan disk I/O seek.",
            "Write throughput tercekik oleh single WAL lock contention dan disk fsync queue.",
            "Operasi maintenance tabel raksasa (ALTER TABLE, VACUUM) mengunci database berjam-jam."
        ], "type": "card"},
        {"title": "MONOTONIC WRITE HOTSPOTS", "items": [
            "Sharding key berbasis waktu (`created_at`) memusatkan 100% traffic tulis pada 1 shard aktif.",
            "Shard lain menganggur sementara shard target kehabisan CPU dan disk IOPS.",
            "Partisi fisik kehilangan fungsi load balancing akibat skew distribusi ekstrem."
        ], "type": "card"},
        {"title": "SCATTER-GATHER EXHAUSTION", "items": [
            "Query tanpa sharding key (`email`) memaksa broadcast paralel ke 100% shard cluster.",
            "Satu query lambat di satu shard menahan total latency (straggler problem).",
            "Throughput cluster ambruk dari $O(N)$ menjadi $O(1)$ di bawah konkurensi tinggi."
        ], "type": "card"}
    ],
    narrative="Dalam arsitektur backend skala besar, database monolitik pada akhirnya akan membentur batasan fisik perangkat keras. Ketika kapasitas buffer pool RAM terlewati, operasi baca-tulis terpaksa bergantung pada I/O disk yang lambat, sementara proses pemeliharaan tabel seperti VACUUM atau DDL migrations melumpuhkan throughput produksi. Lebih parah lagi, ketika tim mencoba melakukan sharding horizontal namun keliru memilih kunci partisi berbasis timestamp, seluruh beban tulis terkonsentrasi pada satu node aktif saja. Sementara itu, query pencarian sekunder tanpa shard key memicu banjir broadcast scatter-gather yang menghabiskan CPU di seluruh penjuru cluster.",
    narrative_en="In hyper-scale backend architectures, monolithic database engines inevitably strike the physical boundaries of vertical hardware scaling. Once working set memory exceeds buffer pool limits, operations degrade into disk-bound latency spirals while administrative maintenance like VACUUM or DDL migrations severely locks production throughput. Worse, when engineering teams transition to horizontal sharding using sequential timestamp keys, write traffic bottlenecks entirely onto a single active shard. Meanwhile, secondary lookups without shard keys trigger scatter-gather broadcast storms that saturate CPU across all nodes.",
    mental_model="Vertical scaling hits hardware walls; naive horizontal sharding without key topology introduces devastating write hotspots and broadcast storms.",
    source_path="labs/32-database-sharding-and-partitioning/content/01-content-brief.md:3-6",
    slide_num="02 / 24",
    layout="cards"
)

# Slide 3: Mengapa review manual / pendekatan naif gagal
add_slide(
    categories=["ANALISIS", "PENDEKATAN NAIF"],
    title="Mengapa Hash Modulo Sederhana & Shard Monotonik Gagal?",
    subtitle="Formula hash(key) % N memicu bencana migrasi 80% data saat menambah satu server.",
    content_items=[
        {"title": "KODE ROUTING HASH MODULO NAIF", "items": [
            "func (m *ModuloRouter) GetShard(key string) (string, error) {",
            "    // BAHAYA FATAL: Indeks berubah total saat N bertambah!",
            "    h := fnv64a(key)",
            "    idx := int(h % uint64(len(m.shards))) // Modulo N",
            "    return m.shards[idx], nil",
            "}",
            "// Scale-out 4 -> 5 nodes: 79.84% data berpindah shard!",
            "// Akibat: Network saturasi, cache invalidation masif, outage produksi."
        ], "type": "code"},
        {"title": "KEGAGALAN FUNDAMENTAL TEKNIK NAIF", "items": [
            "Data Shuffling Masif: Menambah 1 node dari $N$ ke $N+1$ memindahkan $\\approx \\frac{N}{N+1}$ ($75\\% - 80\\%$) seluruh dataset.",
            "Write Hotspot Monotonik: Shard key urut waktu (`2026-09-28`) mengarahkan 1.000 dari 1.000 insert ke satu node saja.",
            "Scatter-Gather Stragglers: Mencari tanpa shard key membroadcast $N$ node; latensi query sama dengan node terlambat ($P_{99}$).",
            "Distributed ID Collision: Auto-increment native lokal menghasilkan ID duplikat lintas shard tanpa sequence coordinator."
        ], "type": "bullet"}
    ],
    narrative="Pendekatan naif kerap mengimplementasikan sharding menggunakan operator modulo sederhana hash(key) % N. Ketika sistem perlu menambah node ke-5 untuk mengatasi lonjakan beban, rumus modulo berubah seketika dan memaksa hampir 80 persen data berpindah antar server di jaringan. Proses resharding ini menyebabkan lonjakan transfer data terabyte, hilangnya cache memory, dan downtime sistemik. Selain itu, penggunaan auto-increment lokal di masing-masing shard mengakibatkan bentrokan primary key, sedangkan query tanpa shard key membebani seluruh server secara bersamaan.",
    narrative_en="Naive distributed implementations commonly route records via elementary modulo hashing hash(key) % N. When scaling from 4 to 5 nodes to handle traffic surges, the modulo denominator shifts across the entire keyspace, forcing nearly 80 percent of all stored records to migrate across the network simultaneously. This resharding storm exhausts network bandwidth, flushes caches, and often triggers catastrophic cascading outages. Concurrently, localized auto-increment sequences generate colliding IDs, and un-routed lookups stall waiting on the slowest broadcast response.",
    mental_model="Hash modulo forces global data relocation on resize; monotonic keys destroy horizontal write concurrency.",
    source_path="internal/sharding/sharding.go:31-53 • content/04-diagrams.md:58-75",
    slide_num="03 / 24",
    layout="code-bullet"
)

# Slide 4: Mental model inti / schema card
add_slide(
    categories=["MENTAL MODEL", "SCHEMA"],
    title="Mental Model: Logical Partitioning vs Physical Sharding",
    subtitle="Membedakan partisi internal satu mesin dari kluster terdistribusi lintas node jaringan.",
    content_items=[
        {"title": "STRUKTUR DATA GO: TABLE PARTITION VS SHARD CLUSTER", "items": [
            "type PartitionRange struct {",
            "    Name  string; Start time.Time; End time.Time // Range [Start, End)",
            "}",
            "type Table struct { // Single Engine Logical Partitioning",
            "    Name string; partitions []*Partition // Pruning lokal via overlap check",
            "}",
            "type Cluster struct { // Multi-Node Physical Sharding",
            "    shards map[string]*Shard // Server fisik independen",
            "    router Router           // Consistent Hash Ring dengan Virtual Nodes",
            "    gsi    *GlobalSecondaryIndex // Lookup Vindex mapping non-shard keys",
            "}"
        ], "type": "code"},
        {"title": "DIKOTOMI ARSITEKTUR FUNDAMENTAL", "items": [
            "Logical Partitioning (1 Mesin): Membagi 1 tabel menjadi sub-tabel terpisah pada disk yang sama. CPU & RAM tetap shared.",
            "Keunggulan Partisi: Query pruning mengeliminasi scan sub-tabel lain; `DROP PARTITION` instan menghapus data tanpa VACUUM.",
            "Physical Sharding (N Mesin): Membagi dataset ke beberapa server mandiri (Shared-Nothing). Melampaui limit RAM/IOPS.",
            "Router Ring & GSI: Router konsisten memetakan shard key ke node; GSI memetakan atribut sekunder (`email` -> `shardKey`)."
        ], "type": "bullet"}
    ],
    narrative="Mental model sistem data terdistribusi bertumpu pada pemisahan konsep yang tegas antara partisi logis dan sharding fisik. Logical Partitioning adalah fitur mesin database tunggal (seperti declarative partitioning di PostgreSQL) yang membagi file fisik tabel berdasarkan rentang waktu untuk mempercepat query pruning dan pembuangan data historis tanpa beban garbage collection. Sebaliknya, Physical Sharding adalah arsitektur shared-nothing di mana setiap shard merupakan instance database mandiri dengan memori dan disk terisolasi. Sharding router dan Global Secondary Index (GSI) bertindak sebagai koordinator untuk mengarahkan trafik aplikasi secara akurat.",
    narrative_en="The architectural mental model requires maintaining an uncompromising conceptual divide between logical partitioning and physical sharding. Logical Partitioning resides within a single database engine (such as PostgreSQL declarative partitioning), segmenting table storage into discrete sub-tables bounded by time ranges to accelerate query pruning and enable instantaneous data drops without table bloat. In stark contrast, Physical Sharding is a shared-nothing distributed architecture where each shard is an isolated database server with dedicated CPU, RAM, and storage, governed by a cluster router and Global Secondary Indexes.",
    mental_model="Partitioning optimizes data lifecycle on one server; sharding scales compute and storage across independent physical nodes.",
    source_path="internal/partitioning/table.go:20-57 • internal/sharding/sharding.go:253-272",
    slide_num="04 / 24",
    layout="schema-card"
)

# Slide 5: Diagram alur/pipeline
add_slide(
    categories=["ARSITEKTUR", "PIPELINE"],
    title="Pipeline End-to-End: Routing, Pruning, & GSI Resolution",
    subtitle="Alur keputusan query: Evaluasi Shard Key → Consistent Hash Ring → Direct vs Scatter-Gather.",
    content_items=[
        {"title": "SEKUENSI EKSEKUSI DATABASE TERDISTRIBUSI", "items": [
            "1. Client Query Ingestion: Request membawa payload record atau query kriteria pencarian.",
            "2. Single-Node Path: Jika query temporal pada Table lokal → QueryRange() mengecek overlap [start, end) → Scan hanya 1 partisi (Pruned 3).",
            "3. Sharded Write Path: Insert(rec) → Hash(rec.ShardKey) dievaluasi pada Consistent Hash Ring → Arahkan langsung ke Shard target.",
            "4. Sharded Read with ShardKey: GetByShardKey(sk, id) → Router lookup $O(\\log V)$ via binary search vnodes → 1 Direct Point Lookup.",
            "5. Secondary Key Read (with GSI): GetByEmailUsingGSI(email, id) → GSI Lookup Vindex memetakan email → ShardKey → Direct Fetch (1 node).",
            "6. Fallback Read (No ShardKey / No GSI): ScatterGatherBroadcast() → Goroutines paralel menembak seluruh N Shard → Gabungkan hasil."
        ], "type": "diagram"}
    ],
    narrative="Pipeline eksekusi ini menggambarkan bagaimana database menangani berbagai pola akses data secara efisien. Pada operasi temporal di single node, query planner langsung memangkas (pruning) partisi yang tidak relevan sehingga I/O disk berkurang drastis. Pada cluster sharded, operasi write dan read yang menyertakan sharding key diarahkan secara presisi ke satu node spesifik via consistent hash ring. Ketika client mencari data berdasarkan atribut non-shard key seperti email, Global Secondary Index memetakan email ke shard key terkait sehingga kueri tetap menjadi point-lookup 1-node dan menghindari bencana broadcast paralel.",
    narrative_en="This comprehensive pipeline illustrates the decision topology for distributed data operations. For temporal single-node queries, the query planner executes partition pruning, reading only overlapping partitions and slashing disk I/O. Within the sharded cluster, read and write operations supplying the shard key route directly through the consistent hashing ring via binary search into the target node. When querying by secondary attributes like email, the Global Secondary Index resolves the shard key upfront, converting an expensive cluster-wide broadcast into a targeted single-node fetch.",
    mental_model="Shard keys route directly; GSI resolves secondary lookups; scatter-gather remains the expensive last resort.",
    source_path="internal/sharding/sharding.go:280-340 • internal/partitioning/table.go:87-119",
    slide_num="05 / 24",
    layout="pipeline-diagram"
)

# Slide 6: Cuplikan kode dominan
add_slide(
    categories=["IMPLEMENTASI", "CODE"],
    title="Implementasi Inti: Consistent Hashing Ring dengan Virtual Nodes",
    subtitle="Algoritma routing ring internal/sharding/sharding.go: binary search vnode & FNV-1a hash.",
    content_items=[
        {"title": "CONSISTENT HASH ROUTER IMPLEMENTATION", "items": [
            "func (ch *ConsistentHashRouter) GetShard(key string) (string, error) {",
            "    ch.mu.RLock(); defer ch.mu.RUnlock()",
            "    if len(ch.ring) == 0 { return \"\", ErrShardNotFound }",
            "    h := hashKey(key)",
            "    // Binary search vnode pertama dengan hash >= h",
            "    idx := sort.Search(len(ch.ring), func(i int) bool {",
            "        return ch.ring[i].hash >= h",
            "    })",
            "    if idx == len(ch.ring) { idx = 0 } // Wrap-around ke awal ring",
            "    return ch.ring[idx].shardID, nil",
            "}"
        ], "type": "code"}
    ],
    narrative="Potongan kode ini merupakan mesin pengarah utama dalam modul internal/sharding/sharding.go. Struktur ConsistentHashRouter merepresentasikan ruang cincin 64-bit yang dipetakan oleh ribuan virtual nodes (misal 100 vnodes per physical shard). Saat kunci data masuk, fungsi hashKey menghitung nilai FNV-1a 64-bit, lalu fungsi sort.Search mengeksekusi binary search dalam waktu logaritmik O(log V) untuk menemukan virtual node pertama yang memiliki nilai hash lebih besar atau sama. Jika nilai hash melewati batas maksimum ring, pencarian wrap-around kembali ke indeks 0, menjamin distribusi seragam.",
    narrative_en="This snippet forms the operational core of the routing engine inside internal/sharding/sharding.go. The ConsistentHashRouter manages a continuous 64-bit circular keyspace populated by virtual nodes (100 vnodes per physical shard). When a key is evaluated, hashKey computes its 64-bit FNV-1a hash, and sort.Search executes a logarithmic binary search O(log V) to identify the first virtual node whose token meets or exceeds the key hash. If the hash overshoots the maximum vnode on the ring, it wraps around to index 0, guaranteeing uniform deterministic routing.",
    mental_model="Binary search on a sorted virtual node ring delivers deterministic, uniform O(log V) routing.",
    source_path="labs/32-database-sharding-and-partitioning/internal/sharding/sharding.go:143-160",
    slide_num="06 / 24",
    layout="code-dominant"
)

# Slide 7: Failure path / dangling reference konkret
add_slide(
    categories=["FAILURE PATH", "BUKTI"],
    title="Failure Path: Bencana Shard Key Monotonik & Resharding Modulo",
    subtitle="Log terminal membuktikan 100% write hotspot pada 1 shard dan lonjakan 79.84% migrasi data.",
    content_items=[
        {"title": "TERMINAL HOTSPOT & RESHARDING TRACE", "items": [
            "[Scenario A] Monotonic Key (date-string '2026-09-28') Sharding:",
            "  Node shard-0: 0 records    [                              ]",
            "  Node shard-1: 0 records    [                              ]",
            "  Node shard-2: 1000 records [##############################] (100% HOTSPOT!)",
            "  Node shard-3: 0 records    [                              ]",
            "Result: Severe Write Hotspot! 100% writes hit single shard.",
            "--------------------------------------------------------------------------------",
            "[Scenario C] Resharding Scale-Out Comparison (4 Shards -> 5 Shards, 10,000 Keys):",
            "  Hash Modulo (N % M) Keys Remapped : 7,984 / 10,000 (79.84% moved! BENCANA NETWORK)",
            "  Consistent Hashing Keys Remapped  : 1,200 / 10,000 (12.00% moved! MINIMAL SHUFFLE)"
        ], "type": "terminal"},
        {"title": "ANALISIS KEGAGALAN SISTEMATIS", "items": [
            "Hotspot Monotonik: Shard key berbasis tanggal memusatkan 1.000 transaksi ke `shard-2`; 3 shard lainnya menganggur total.",
            "Bencana Modulo Resharding: Penambahan 1 node merombak 79.84% data; jaringan tercekik sinkronisasi dan database terkunci.",
            "Mitigasi Lab: Consistent hashing ring membatasi relokasi data ke angka 12.00%, mendekati batas teoritis minimal $1/(N+1) = 20\\%$."
        ], "type": "bullet"}
    ],
    narrative="Bukti failure path ini diambil langsung dari hasil eksekusi program demo cmd/demo/main.go. Pada Skenario A, penggunaan kunci sharding monotonik berbasis tanggal kalender menyebabkan 1.000 dari 1.000 data masuk ke satu node shard-2, membuktikan bagaimana shard key yang salah memusatkan beban 100 persen dan melumpuhkan kapasitas I/O satu server. Pada Skenario C, saat cluster melakukan scale-out dari 4 menjadi 5 node, algoritma hash modulo memindahkan 7.984 dari 10.000 key (79.84 persen). Sebaliknya, consistent hashing hanya memindahkan 1.200 key (12.00 persen), membuktikan keunggulan arsitekturalnya.",
    narrative_en="This concrete failure trace is captured directly from the execution of cmd/demo/main.go. In Scenario A, employing a monotonic date-based sharding key channels 1,000 out of 1,000 inserts into a single physical instance (shard-2), visually proving how improper key selection creates a 100% write bottleneck that nullifies horizontal capacity. In Scenario C, expanding the cluster from 4 to 5 shards under modulo routing relocates 7,984 of 10,000 keys (79.84%). In contrast, consistent hashing moves merely 1,200 keys (12.00%), preserving stability.",
    mental_model="Monotonic keys crush single instances with 100% write skew; modulo resizing triggers network-saturating data avalanches.",
    source_path="cmd/demo/main.go:52-139 • engineering/03-execution-result.md:74-95",
    slide_num="07 / 24",
    layout="failure-path"
)

# Slide 8: Tabel perbandingan mekanisme vs gejala
add_slide(
    categories=["PERBANDINGAN", "MEKANISME"],
    title="Mekanisme vs Gejala: Evaluasi Strategi Partisi & Sharding",
    subtitle="Perbandingan mendalam: Karakteristik, skala beban, trade-off, dan ketahanan produksi.",
    content_items=[
        {"title": "TABEL EVALUASI MEKANISME SKALABILITAS DATABASE", "items": [
            "Mekanisme | Arsitektur Target | Keunggulan Utama | Limitasi / Kerentanan Kritis | Ketahanan Produksi",
            "Logical Partitioning | Single Engine (1 Node) | Query pruning, fast `DROP PARTITION` | Shared CPU/RAM; tidak menyelesaikan limit I/O fisik | Sangat Baik (Lifecycle Data)",
            "Hash Modulo ($K \\% N$) | Multi-Node Cluster | Algoritma $O(1)$ sangat sederhana | Relokasi data masif $\\approx \\frac{N}{N+1}$ saat resharding | Buruk (Dilarang saat Resizing Dinamis)",
            "Consistent Hashing | Multi-Node Cluster | Relokasi minimal $\\approx \\frac{1}{N+1}$ saat scale-out | Membutuhkan virtual nodes untuk mencegah data skew | Unggul (Standar Industri)",
            "Scatter-Gather Broadcast | Multi-Node Cluster | Mendukung query atribut non-shard key | Latensi P99 dibatasi node terlambat (straggler) | Rawan (Bottleneck pada Skala Besar)",
            "Global Secondary Index | Multi-Node Cluster | Point-lookup terarah 1-node via lookup key | Double-write overhead saat Insert/Update | Unggul (Untuk Query Frekuensi Tinggi)"
        ], "type": "table"}
    ],
    narrative="Tabel perbandingan ini mengevaluasi spektrum lengkap solusi skalabilitas database. Logical partitioning sangat ampuh untuk mengelola siklus hidup data besar pada satu mesin (seperti membuang data lama secara instan tanpa operasi DELETE/VACUUM yang lambat), namun tidak mampu mengatasi keterbatasan memori dan CPU server tunggal. Untuk arsitektur multi-node, modulo hashing terbukti berbahaya saat cluster bertumbuh karena merombak mayoritas lokasi data. Consistent hashing dengan virtual nodes menjadi standar emas karena meminimalkan migrasi, sementara Global Secondary Index menghilangkan ketergantungan pada scatter-gather broadcast yang boros sumber daya.",
    narrative_en="This comparative matrix evaluates the entire spectrum of database partitioning and sharding architectures. Logical partitioning excels at managing single-node data lifecycles—enabling instant partition drops without locking tables—yet cannot bypass single-server CPU and memory ceilings. Across distributed multi-node clusters, simple modulo hashing is brittle due to catastrophic resharding churn. Consistent hashing with virtual nodes provides the industry-standard scale-out foundation, paired with Global Secondary Indexes to eliminate latency-destroying scatter-gather broadcasts.",
    mental_model="Logical partitioning manages local table lifecycles; consistent hashing and GSI power scalable distributed clusters.",
    source_path="labs/32-database-sharding-and-partitioning/content/02-master-draft.md:9-40",
    slide_num="08 / 24",
    layout="comparison-table"
)

# Slide 9: Decision sequence (3 focused decision cards)
add_slide(
    categories=["DESIGN DECISION", "ARSITEKTUR"],
    title="3 Keputusan Desain Kritis Arsitektur Sharding",
    subtitle="Urutan keputusan wajib sebelum memecah database ke kluster multi-node terdistribusi.",
    content_items=[
        {"title": "1. SHARD KEY: HIGH CARDINALITY & UNIFORMITY", "items": [
            "Pilih atribut dengan kardinalitas tinggi dan frekuensi rendah (contoh: `account_id`, `tenant_id`, `uuid`).",
            "Haramkan nilai monotonik (`created_at`, auto-increment) sebagai range shard key untuk mencegah 100% write hotspot.",
            "Terapkan Hashed Sharding jika beban kerja didominasi oleh operasi tulis sekuensial intensif."
        ], "type": "card"},
        {"title": "2. ROUTER: CONSISTENT HASH RING & VNODES", "items": [
            "Gunakan cincin hash konsisten (FNV-1a / Murmur3) dengan 50–150 virtual nodes per server fisik.",
            "Virtual nodes mendistribusikan token secara acak-merata untuk mencegah ketimpangan beban (hot shard).",
            "Saat scale-out dari $N$ ke $N+1$, batasi perpindahan data hanya pada porsi $\\approx 1/(N+1)$ dari total dataset."
        ], "type": "card"},
        {"title": "3. SECONDARY ACCESS: GSI INDEXING POLICY", "items": [
            "Identifikasi query non-shard key yang berfrekuensi tinggi (misal login via `email` pada cluster bertarget `user_id`).",
            "Bangun Global Secondary Index (Lookup Vindex) untuk memetakan atribut sekunder ke ShardKey target.",
            "Terima trade-off double-write saat insert demi menjamin latensi baca tetap 1-node point-lookup."
        ], "type": "card"}
    ],
    narrative="Tiga keputusan arsitektural ini menentukan keberhasilan transisi dari database monolitik ke sistem terdistribusi. Pertama, pemilihan sharding key harus memprioritaskan atribut dengan sebaran nilai yang sangat beragam (high cardinality) seperti tenant ID atau user ID guna mendistribusikan data secara seimbang ke seluruh node. Kedua, router wajib menerapkan consistent hashing dengan virtual nodes yang cukup (50 hingga 150 vnodes) agar penambahan node baru tidak memicu perombakan data global. Ketiga, query non-shard key harus dilayani oleh Global Secondary Index untuk mencegah broadcast paralel yang menghancurkan throughput cluster.",
    narrative_en="These three architectural decisions determine the operational viability of any distributed database transition. First, shard key selection must strictly enforce high-cardinality, uniformly distributed identifiers like tenant IDs to ensure write distribution across nodes without skew. Second, the cluster routing tier must adopt consistent hashing with 50 to 150 virtual nodes per physical shard, dampening data reshuffling during cluster expansion to minimal increments. Third, secondary query patterns must be indexed via Global Secondary Indexes, protecting cluster CPU from broadcast exhaustion.",
    mental_model="Select high-cardinality keys, govern routing with virtual node rings, and resolve secondary lookups via GSI tables.",
    source_path="labs/32-database-sharding-and-partitioning/content/05-key-takeaways.md",
    slide_num="09 / 24",
    layout="decision-cards"
)

# Slide 10: Bukti eksekusi / terminal output
add_slide(
    categories=["BUKTI", "DEMO"],
    title="Bukti Eksekusi: Verifikasi Pruning, Hashing, GSI, & UUIDv7",
    subtitle="Output terminal go run ./cmd/demo membuktikan keunggulan pruning, distribusi vnodes, dan kecepatan GSI.",
    content_items=[
        {"title": "DEMO EXECUTION TRACE", "items": [
            "=== LAB 32: DATABASE SHARDING AND PARTITIONING DEMONSTRATION ===",
            "[1. Single-Node Partitioning] Query Range: 2026-04-01 to 2026-06-30",
            "  Records Found: 1 | Partitions Scanned: 1 / 4 (Pruned 3 partitions)",
            "[2. Sharding Key Selection: Monotonic vs High-Cardinality]",
            "  Scenario A (date-string): Node shard-2 got 1000 records (100% Hotspot)",
            "  Scenario B (user_id)    : Node shard-0: 400 | shard-2: 200 | shard-3: 400 (Distributed)",
            "[3. Resharding Scale-out: 4 Shards -> 5 Shards (10,000 Keys)]",
            "  Hash Modulo (N % M) Keys Remapped : 7,984 / 10,000 (79.84% moved)",
            "  Consistent Hashing Keys Remapped  : 1,200 / 10,000 (12.00% moved) | Theoretical: ~20.00%",
            "[4. Non-Sharded Attributes: Scatter-Gather vs GSI]",
            "  Scatter-Gather Query : Nodes Broadcasted: 4 / 4 | Execution Time: 92µs",
            "  Global Secondary Index: Nodes Broadcasted: 1 / 4 | Execution Time: 1µs (92x Faster!)",
            "[5. Distributed Unique ID Generation]",
            "  UUIDv7 (Time-Ordered 128-bit) : 01a0eb21-fc8c-7aba-abc5-bcb4361a5e7b",
            "  Sequence Block Allocator IDs  : 1 2 3 4 5 6 7 8 (Block Size=5)",
            "=== [DEMO COMPLETE] All database sharding & partitioning concepts executed ==="
        ], "type": "terminal"}
    ],
    narrative="Bukti eksekusi terminal membuktikan secara empiris seluruh klaim arsitektur dalam laboratorium ini. Pada modul partisi logis, query range pada kuartal kedua hanya memindai 1 partisi dan memangkas 3 partisi lainnya secara instan. Pada perbandingan resharding, consistent hashing ring membatasi relokasi data pada angka 12.00 persen dibandingkan modulo hashing yang merombak 79.84 persen data. Pada pengujian akses non-shard key, Global Secondary Index berhasil mengarahkan query langsung ke 1 shard dalam waktu 1 mikrodetik, mengungguli metode scatter-gather yang membroadcast 4 shard dalam 92 mikrodetik.",
    narrative_en="This terminal trace provides empirical verification for all architectural claims established in the lab. In the logical partitioning benchmark, a quarterly range query scans exactly 1 partition while pruning the remaining 3 instantly. In the scale-out resharding trial, consistent hashing restricts key migration to 12.00% compared to 79.84% churn under modulo hashing. Finally, secondary key evaluation proves Global Secondary Indexes execute targeted 1-shard point lookups in 1 microsecond, outperforming 4-shard scatter-gather broadcasts running at 92 microseconds.",
    mental_model="Pruning eliminates 75% partition I/O; consistent hashing cuts migration 6.6x; GSI delivers 92x faster point lookups.",
    source_path="cmd/demo/main.go:1-210 • engineering/03-execution-result.md:65-111",
    slide_num="10 / 24",
    layout="terminal-output"
)

# Slide 11: Analisis batasan bukti
add_slide(
    categories=["BATASAN", "ANALISIS"],
    title="YANG TIDAK DIBUKTIKAN",
    subtitle="Batasan ruang lingkup laboratorium yang wajib dipahami sebelum implementasi produksi.",
    content_items=[
        {"title": "YANG TIDAK DIBUKTIKAN", "items": [
            "Penyimpanan tabel dan shard berbasis in-memory Go struct; tidak menguji disk persistence engine (InnoDB, RocksDB, LSM-Tree).",
            "Tidak menyertakan Two-Phase Commit (2PC) coordinator untuk transaksi cross-shard atomik (XA transaction).",
            "Global Secondary Index (GSI) dimodelkan secara in-memory sinkron; tidak menguji asynchronous eventual consistency replikasi GSI.",
            "Tidak menyertakan dynamic shard split / auto-rebalancing engine otomatis saat kapasitas fisik storage mencapai batas.",
            "Tidak menyimulasikan network partitions, node failures, atau consensus replication (Raft/Paxos) antar replika shard.",
            "Tidak menguji distributed deadlock detection saat dua transaksi cross-shard saling mengunci record di shard berlawanan."
        ], "type": "bullet"}
    ],
    narrative="Laboratorium ini dirancang mandiri untuk membuktikan mekanisme algoritma routing, partisi logis, dan mitigasi query tanpa dependensi eksternal. Namun demikian, implementasi ini memiliki batasan ruang lingkup tegas yang harus diperhatikan oleh arsitek sistem. Komponen shard dan GSI berjalan di dalam memori proses tunggal tanpa persistensi disk fisik. Lab ini tidak menguji protokol Two-Phase Commit (2PC) terdistribusi yang diperlukan untuk transaksi atomik lintas shard, dan tidak memodelkan replikasi konsensus Raft atau deteksi deadlock terdistribusi yang umum terjadi pada database skala enterprise.",
    narrative_en="This laboratory is purposefully self-contained to isolate and verify routing mechanics, logical pruning, and index resolution without heavyweight external engines. Consequently, its technical guarantees are strictly bounded. Shard data structures and GSI tables reside in-process in Go memory without persistent block device drivers. The lab does not implement distributed Two-Phase Commit (2PC) engines for cross-shard ACID transactions, nor does it simulate consensus-based replication protocols like Raft or distributed deadlock detection algorithms required in production clusters.",
    mental_model="Acknowledge educational bounds: in-memory routing models routing algorithms; production requires 2PC engines and consensus replication.",
    source_path="labs/32-database-sharding-and-partitioning/content/01-content-brief.md:25-29",
    slide_num="11 / 24",
    layout="constraints-list"
)

# Slide 12: Trade-off arsitektur
add_slide(
    categories=["TRADE-OFF", "ARSITEKTUR"],
    title="Trade-off Arsitektur: Sharding vs Kompleksitas Terdistribusi",
    subtitle="Menilai keuntungan kapasitas horisontal terhadap beban biaya operasional dan konsistensi data.",
    content_items=[
        {"title": "KEUNTUNGAN & KAPASITAS SKALABILITAS", "items": [
            "Skalabilitas Horisontal Linear: Menembus batasan memori RAM, disk IOPS, dan write bandwidth mesin database tunggal.",
            "Blast Radius Terisolasi: Kerusakan hardware pada Shard 2 tidak mematikan ketersediaan data pada Shard 0, 1, dan 3.",
            "Query Pruning Efisien: Rentang partisi logis memangkas scan disk secara drastis untuk tabel berbasis waktu."
        ], "type": "bullet"},
        {"title": "BIAYA OPERASIONAL & KOMPLEKSITAS", "items": [
            "Kehilangan JOIN Relasional: Operasi SQL JOIN lintas shard key mustahil dilakukan tanpa scatter-gather atau agregasi di aplikasi.",
            "Overhead Transaksi Terdistribusi: Two-Phase Commit (2PC) melipatgandakan latensi jaringan dan berisiko fractured reads.",
            "Double-Write Overhead pada GSI: Menjaga indeks sekunder global membutuhkan write ganda yang memperlambat operasi insert."
        ], "type": "bullet"}
    ],
    narrative="Keputusan mengadopsi physical sharding selalu membawa kompromi teknik yang signifikan. Di satu sisi, sistem memperoleh skalabilitas linear dan isolasi kegagalan di mana matinya satu server tidak melumpuhkan seluruh aplikasi. Namun di sisi lain, pengembang harus membayar harga arsitektural yang mahal: hilangnya kapabilitas SQL JOIN antar tabel yang berbeda shard, tingginya latensi koordinasi transaksi Two-Phase Commit, serta keharusan mengelola overhead double-write untuk menjaga konsistensi Global Secondary Index.",
    narrative_en="Architecting a distributed sharded database demands a deliberate reckoning with architectural trade-offs. While sharding unlocks unbounded horizontal write throughput and insulates blast radiuses across physical nodes, it dismantles native relational conveniences. Developers forfeit efficient cross-shard relational JOIN operations, incur punishing network latency penalties under Two-Phase Commit protocols, and must maintain dual-write overheads to keep Global Secondary Indexes synchronized.",
    mental_model="Sharding buys unbounded scale and fault isolation at the cost of cross-shard joins, 2PC latency, and index dual-writes.",
    source_path="labs/32-database-sharding-and-partitioning/content/02-master-draft.md:270-295",
    slide_num="12 / 24",
    layout="tradeoff"
)

# Slide 13: Concurrency / validation pipeline
add_slide(
    categories=["CONCURRENCY", "VALIDATION"],
    title="Validasi Konkurensi: 200 Worker Simultan pada Sharded Cluster",
    subtitle="TestConcurrentClusterAccess membuktikan integritas state di bawah konkurensi Insert, Read, dan GSI.",
    content_items=[
        {"title": "KODE PENGUJIAN STRESS KONKURENSI", "items": [
            "func TestConcurrentClusterAccess(t *testing.T) {",
            "    router := sharding.NewConsistentHashRouter(50, []string{\"s1\", \"s2\", \"s3\", \"s4\"})",
            "    cluster := sharding.NewCluster(router)",
            "    var wg sync.WaitGroup; numOps := 200",
            "    for i := 0; i < numOps; i++ {",
            "        wg.Add(1)",
            "        go func(idx int) {",
            "            defer wg.Done()",
            "            sk := fmt.Sprintf(\"tenant-%d\", idx%10)",
            "            recID := fmt.Sprintf(\"id-%d\", idx)",
            "            email := fmt.Sprintf(\"user-%d@test.com\", idx)",
            "            _ = cluster.Insert(sharding.Record{ID: recID, ShardKey: sk, Email: email, Data: \"concurrent-test\"})",
            "            _, _ = cluster.GetByShardKey(sk, recID)",
            "            _, _ = cluster.GetByEmailUsingGSI(email, recID)",
            "        }(i)",
            "    }",
            "    wg.Wait()",
            "    // Assert: total seluruh record di seluruh shard == 200 (Zero Data Loss)",
            "}"
        ], "type": "code"},
        {"title": "VERIFIKASI THREAD SAFETY & RACE DETECTOR", "items": [
            "RWMutex Granular: Setiap Shard dan GSI dilindungi oleh `sync.RWMutex` independen, mencegah lock contention antar shard.",
            "Zero Data Race: Eksekusi `go test -race` lulus 100% bersih tanpa mendeteksi data race condition pada 200 goroutine simultan.",
            "Total Invariant: 200 transaksi insert terdistribusi secara atomik; verifikasi cluster count membuktikan tidak ada data hilang."
        ], "type": "bullet"}
    ],
    narrative="Uji konkurensi TestConcurrentClusterAccess memverifikasi ketahanan cluster saat diserbu oleh 200 goroutine simultan yang melakukan insert, pembacaan berbasis shard key, dan point-lookup via GSI. Karena setiap instance Shard dan struktur GlobalSecondaryIndex memiliki sync.RWMutex masing-masing, penguncian data bersifat granular sehingga operasi pada shard-1 tidak menghalangi operasi pada shard-2. Pengujian dengan flag -race membuktikan ketiadaan memory race condition dan menegaskan integritas data tersimpan secara utuh tepat 200 record.",
    narrative_en="The concurrency validation suite TestConcurrentClusterAccess proves system resilience under a barrage of 200 concurrent goroutines executing parallel inserts, shard-key lookups, and GSI queries. Because each Shard instance and the GlobalSecondaryIndex manage independent sync.RWMutex primitives, thread contention remains isolated—writes on shard-1 proceed without blocking reads on shard-2. Execution under the Go race detector confirms zero data races while maintaining an exact 200-record state invariant.",
    mental_model="Granular per-shard mutexes isolate concurrency contention; race detectors verify thread-safe distributed routing.",
    source_path="labs/32-database-sharding-and-partitioning/tests/sharding_test.go:212-249",
    slide_num="13 / 24",
    layout="concurrency-flow"
)

# Slide 14: Before vs After comparison cards
add_slide(
    categories=["BEFORE/AFTER", "PERUBAHAN"],
    title="Before vs After: Evolusi Arsitektur Manajemen Data",
    subtitle="Transformasi dari database monolitik rentan bottleneck menuju cluster sharded terkoordinasi.",
    content_items=[
        {"title": "BEFORE (Monolithic DB & Naive Modulo)", "items": [
            "Kapasitas: Terkunci pada limit hardware 1 server fisik; disk I/O mentok saat data membesar.",
            "Query Temporal: Full table scan membebani jutaan row; maintenance `DELETE` memicu table bloat.",
            "Resharding: Menggunakan rumus hash modulo; menambah node baru memindahkan 80% data.",
            "Query Sekunder: Menembak seluruh database tanpa filter; lonjakan CPU membakar server.",
            "Primary Keys: Auto-increment lokal bertabrakan jika dipisah ke beberapa node."
        ], "type": "card"},
        {"title": "AFTER (Partitioning, Consistent Ring, & GSI)", "items": [
            "Kapasitas: Skalabilitas horisontal multi-node terdistribusi; beban write terbagi merata.",
            "Query Temporal: Range logical partitioning memangkas 75% partisi (query pruning); drop instan.",
            "Resharding: Consistent hashing ring dengan vnodes membatasi migrasi hanya $\\approx 12\\% - 16\\%$.",
            "Query Sekunder: Global Secondary Index memetakan atribut sekunder langsung ke 1 shard target.",
            "Primary Keys: RFC 9562 UUIDv7 dan Sequence Block Allocator menjamin ID unik terurut waktu."
        ], "type": "card"}
    ],
    narrative="Perbandingan before-after ini merangkum evolusi menyeluruh dari arsitektur penyimpanan data. Pada paradigma lama, tim bergantung pada satu server monolitik besar yang rentan mengalami disk bottleneck, menggunakan rumus modulo rapuh yang memicu badai migrasi data saat scaling, serta menderita bentrokan primary key. Dengan arsitektur modern, tabel lokal dipangkas secara logis untuk efisiensi siklus hidup data, cluster didistribusikan via consistent hashing ring dengan virtual nodes, pencarian sekunder diarahkan via GSI, dan ID unik digenerasi secara aman via standar UUIDv7.",
    narrative_en="This before-and-after comparison highlights the architectural paradigm shift in data layer engineering. Legacy architectures relied on single monolithic instances constrained by disk throughput ceilings, fragile modulo routing that triggered resharding avalanches, and local auto-increments prone to collisions. The modern distributed standard leverages logical table partitioning for temporal pruning, consistent hashing rings with virtual nodes for minimal data reshuffling, GSI tables for targeted point lookups, and RFC 9562 UUIDv7 for collision-free time-ordered primary keys.",
    mental_model="Before: hardware ceilings, 80% resharding churn, and broadcast storms. After: 75% pruning, 12% scale-out shuffle, and 1-node GSI routing.",
    source_path="labs/32-database-sharding-and-partitioning/content/02-master-draft.md:410-435",
    slide_num="14 / 24",
    layout="before-after"
)

# Slide 15: Analisis mendalam - Test 1 Logical Partition Pruning
add_slide(
    categories=["BUKTI TEST", "LOGICAL PARTITIONING"],
    title="Test 1: Validasi Query Pruning pada Range Partitioning",
    subtitle="TestPartitionPruning membuktikan hanya 1 dari 3 partisi yang discan dan drop tabel instan.",
    content_items=[
        {"title": "KODE PENGUJIAN QUERY PRUNING & FAST DROP", "items": [
            "func TestPartitionPruning(t *testing.T) {",
            "    ranges := []partitioning.PartitionRange{",
            "        {Name: \"p2026_01\", Start: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)},",
            "        {Name: \"p2026_02\", Start: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)},",
            "        {Name: \"p2026_03\", Start: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)},",
            "    }",
            "    tbl := partitioning.NewTable(\"orders\", ranges)",
            "    // Insert data Jan dan Feb...",
            "    res := tbl.QueryRange(febStart, febEnd)",
            "    // Assert: res.PartitionsScanned == 1 && res.TotalPartitions == 3 (Pruning Sukses!)",
            "    // Drop partition Jan p2026_01...",
            "    // Assert: tbl.PartitionCount() == 2 (Hapus instan tanpa scan tabel)"
            "}"
        ], "type": "code"},
        {"title": "VERIFIKASI INVARIAN QUERY PRUNING", "items": [
            "Overlap Condition Pruning: Logika `start.Before(p.Range.End) && end.After(p.Range.Start)` mengabaikan partisi Jan dan Mar.",
            "Eliminasi Disk I/O: Query hanya membaca memori partisi Feb, menghemat 66.7% operasi I/O pemindaian.",
            "Fast Drop Partition: Menghapus partisi `p2026_01` membuang seluruh row di dalamnya secara instan tanpa VACUUM bloat."
        ], "type": "bullet"}
    ],
    narrative="Unit test pertama memvalidasi mekanisme pemangkasan partisi pada internal/partitioning/table.go. Tabel orders dikonfigurasikan dengan tiga rentang partisi bulanan. Saat kueri range dieksekusi khusus untuk bulan Februari 2026, sistem secara cerdas mengevaluasi batasan waktu partisi dan hanya memindai partisi p2026_02 (PartitionsScanned=1 dari 3). Selain itu, pengujian membuktikan bahwa method DropPartition menghapus seluruh data bulan Januari seketika secara $O(1)$ tanpa perlu memindai dan menghapus baris satu per satu.",
    narrative_en="The first unit test suite validates partition pruning mechanics inside internal/partitioning/table.go. The orders table is instantiated with three consecutive monthly ranges. When a range query targeting February 2026 executes, the pruning evaluator inspects interval boundaries, scanning strictly partition p2026_02 (PartitionsScanned=1 out of 3 total). Furthermore, DropPartition demonstrates instant O(1) partition eviction for January data without table-scanning row-by-row DELETE operations.",
    mental_model="Pruning logic start < end && end > start bypasses irrelevant partitions; dropping partitions avoids VACUUM overhead.",
    source_path="labs/32-database-sharding-and-partitioning/tests/sharding_test.go:15-61",
    slide_num="15 / 24",
    layout="test-analysis"
)

# Slide 16: Analisis mendalam - Test 2 Consistent Hash vs Modulo Resharding
add_slide(
    categories=["BUKTI TEST", "ROUTING & RESHARDING"],
    title="Test 2: Evaluasi Relokasi Data Scale-Out (5.000 Keys)",
    subtitle="TestRoutingAndConsistentHashRelocation membuktikan Modulo merelokasi 75.12% vs Consistent Hash 16.00%.",
    content_items=[
        {"title": "KODE PERBANDINGAN RELOKASI PADA SCALE-OUT", "items": [
            "func TestRoutingAndConsistentHashRelocation(t *testing.T) {",
            "    numKeys := 5000",
            "    // 1. Modulo router test: 3 nodes -> 4 nodes",
            "    modRouter := sharding.NewModuloRouter([]string{\"shard-1\", \"shard-2\", \"shard-3\"})",
            "    // Catat penempatan awal...",
            "    modRouter.AddShard(\"shard-4\") // Scale-out 3 -> 4",
            "    // Hash Modulo memindahkan 3756 / 5000 keys (75.12%)",
            "    // 2. Consistent Hash router test (150 vnodes): 3 nodes -> 4 nodes",
            "    chRouter := sharding.NewConsistentHashRouter(150, []string{\"shard-1\", \"shard-2\", \"shard-3\"})",
            "    // Catat penempatan awal...",
            "    chRouter.AddShard(\"shard-4\") // Scale-out 3 -> 4",
            "    // Consistent Hash memindahkan 800 / 5000 keys (16.00%)",
            "}"
        ], "type": "code"},
        {"title": "ANALISIS RASIO PERPINDAHAN KUNCI", "items": [
            "Modulo Churn: Formula $K \\pmod N$ merelokasi $75.12\\%$ kunci (mendekati prediksi teoritis $3/4 = 75.00\\%$).",
            "Consistent Hash Stability: Cincin hash dengan 150 vnodes hanya merelokasi $16.00\\%$ kunci (dalam batas aman $5\\% - 40\\%$).",
            "Efisiensi Jaringan: Consistent hashing menghemat $78.7\\%$ beban transfer data antar server saat cluster melakukan ekspansi kapasitas."
        ], "type": "bullet"}
    ],
    narrative="Pengujian kedua membandingkan dampak scale-out cluster dari 3 node menjadi 4 node terhadap 5.000 kunci data unik. Algoritma modulo hashing terbukti merelokasi 3.756 kunci (75.12 persen), memvalidasi model matematika di mana penambahan node ke-4 mengacak ulang tiga perempat penempatan data. Sebaliknya, router consistent hashing dengan 150 virtual nodes hanya memindahkan 800 kunci (16.00 persen). Data yang tersisa (84 persen) tetap berada di shard aslinya tanpa perlu dipindahkan sama sekali melalui jaringan.",
    narrative_en="The second test directly contrasts cluster scale-out behavior when expanding from 3 to 4 nodes across 5,000 distinct keys. Modulo hashing forcefully remaps 3,756 keys (75.12%), corroborating the theoretical prediction where adding a 4th node shuffles three-quarters of the entire keyspace. In contrast, the consistent hashing router configured with 150 virtual nodes moves merely 800 keys (16.00%). The remaining 84% of keys retain their physical location untouched, conserving cluster network capacity.",
    mental_model="Modulo reshuffling displaces 75% of keys; consistent hashing confines key migration to the newly added sector.",
    source_path="labs/32-database-sharding-and-partitioning/tests/sharding_test.go:63-120",
    slide_num="16 / 24",
    layout="test-analysis"
)

# Slide 17: Analisis mendalam - Test 3 Scatter-Gather vs GSI Resolution
add_slide(
    categories=["BUKTI TEST", "QUERY OPTIMIZATION"],
    title="Test 3: Eliminasi Broadcast via Global Secondary Index",
    subtitle="TestClusterScatterGatherAndGSI membuktikan point-lookup 1-node dan timeout cancellation aman.",
    content_items=[
        {"title": "PENGUJIAN GSI POINT-LOOKUP & TIMEOUT", "items": [
            "// 1. Direct Point-Lookup menggunakan ShardKey",
            "rec, err := cluster.GetByShardKey(\"tenant-A\", \"rec-1\") // Langsung ke Shard 0",
            "// 2. Lookup via Global Secondary Index (GSI)",
            "recGSI, err := cluster.GetByEmailUsingGSI(\"bob@example.com\", \"rec-2\")",
            "// Assert: recGSI.Data == \"Bob Data\" (Tepat 1 shard dihubungi, 0 broadcast overhead)",
            "// 3. Fallback Scatter-Gather Broadcast",
            "sgRes := cluster.ScatterGatherBroadcast(func(r sharding.Record) bool { return r.Email == \"carol@example.com\" })",
            "// Assert: sgRes.ShardsQueried == 3 (Broadcast menembus seluruh 3 shard)",
            "// 4. Scatter-Gather dengan Context Timeout / Cancellation",
            "ctx, cancel := context.WithCancel(context.Background()); cancel()",
            "canceledRes := cluster.ScatterGatherBroadcastWithContext(ctx, predicate)",
            "// Assert: canceledRes.ShardResponded == 0 (Goroutine dibatalkan seketika)"
        ], "type": "code"},
        {"title": "CHECKPOINT OPTIMASI NON-SHARDED QUERY", "items": [
            "GSI Direct Routing: Memetakan `email -> shardKey` memungkinkan pemanggilan langsung tanpa menyentuh shard lain.",
            "Scatter-Gather Broadcast: Menanyakan 3/3 shard secara paralel via Goroutines; mahal namun menjamin kelengkapan data.",
            "Context Safety: Channel cancellation `ctx.Done()` menghentikan pemrosesan goroutine di semua shard jika client timeout."
        ], "type": "bullet"}
    ],
    narrative="Unit test ketiga menguji arsitektur pencarian atribut sekunder yang tidak dijadikan sharding key. Fungsi GetByEmailUsingGSI membuktikan bahwa pemetaan Global Secondary Index memungkinkan sistem menemukan data hanya dengan mengakses 1 shard spesifik tanpa membroadcast ke shard lain. Ketika pencarian dilakukan tanpa bantuan GSI melalui ScatterGatherBroadcast, seluruh 3 shard terpaksa dipindai secara paralel. Pengujian juga memvalidasi bahwa jika klien membatalkan request melalui context.Context, seluruh goroutine broadcast langsung berhenti bekerja dengan aman.",
    narrative_en="The third unit test evaluates querying mechanics for non-shard key attributes. The GetByEmailUsingGSI implementation confirms that maintaining a Global Secondary Index table allows the router to execute direct single-shard point lookups without broadcasting across the cluster. When querying without GSI via ScatterGatherBroadcast, all 3 shards must be scanned in parallel. Furthermore, the test verifies that supplying a pre-canceled context cleanly halts goroutine execution across all shards.",
    mental_model="GSI transforms cluster-wide broadcast queries into targeted single-node reads; context cancels rogue scatter goroutines.",
    source_path="labs/32-database-sharding-and-partitioning/tests/sharding_test.go:121-170",
    slide_num="17 / 24",
    layout="test-analysis"
)

# Slide 18: Studi Kasus Industri: Vitess Vindexes & Citus Distributed Tables
add_slide(
    categories=["STUDI KASUS", "INDUSTRI"],
    title="Studi Kasus Industri: Vitess (YouTube) & Citus (Postgres)",
    subtitle="Bagaimana sistem skala internet memecahkan masalah sharding, routing, dan cross-shard queries.",
    content_items=[
        {"title": "VITESS: YOUTUBE & SLACK HORIZONTAL SCALING", "items": [
            "Vindex Routing Engine: Vitess menggunakan abstraksi Vindex (Functional & Lookup Vindexes) untuk memetakan SQL query ke tablet MySQL.",
            "Lookup Vindex as GSI: Indeks sekunder global disimpan di keyspace terpisah untuk melayani query non-shard key (`user_email`).",
            "Sequence Keyspace: Menghindari bentrokan ID numerik dengan mengalokasikan blok auto-increment secara terpusat."
        ], "type": "bullet"},
        {"title": "CITUS: DISTRIBUTED POSTGRESQL EXTENSION", "items": [
            "Co-Location Strategy: Tabel yang memiliki relasi (misal `orders` dan `order_items`) di-shard menggunakan `tenant_id` yang sama.",
            "Local Foreign Key Joins: Data terkait berada pada physical worker yang sama, memungkinkan query SQL JOIN dieksekusi lokal.",
            "Distributed Query Planner: Mentransformasikan SQL kompleks menjadi paralel query snippets ke masing-masing node worker."
        ], "type": "bullet"}
    ],
    narrative="Arsitektur industri membuktikan bahwa prinsip-prinsip dalam lab ini adalah fondasi dari platform data terdistribusi terkemuka. Vitess yang digunakan oleh YouTube dan Slack memelopori konsep Vindexes di mana router cerdas VTGate memetakan query ke ratusan shard MySQL tanpa mengubah kode aplikasi. Sementara itu, Citus pada ekosistem PostgreSQL menekankan strategi Co-Location: dengan menyelaraskan sharding key antara tabel induk dan anak, transaksi dan relasi JOIN dapat diselesaikan sepenuhnya di level server lokal tanpa melintasi jaringan.",
    narrative_en="Industry production architectures validate that the core tenets implemented in this lab power modern enterprise data platforms. Vitess, scaling transactional backends at YouTube and Slack, pioneered Vindexes where the VTGate router steers queries to hundreds of MySQL shards transparently. Concurrently, Citus within the PostgreSQL ecosystem champions table Co-Location: by sharing common distribution keys across related tables, foreign key constraints and relational JOINs execute locally without cross-shard network chatter.",
    mental_model="Co-locate related tables by tenant ID to preserve local SQL joins; rely on Vindex engines for transparent routing.",
    source_path="labs/32-database-sharding-and-partitioning/content/02-master-draft.md:300-330",
    slide_num="18 / 24",
    layout="cards"
)

# Slide 19: Analisis Teknis: Distributed Unique ID Generation (UUIDv7 vs Sequences)
add_slide(
    categories=["ANALISIS TEKNIS", "DISTRIBUTED ID"],
    title="Analisis Teknis: Solusi Primary Key Terdistribusi",
    subtitle="Mencegah tabrakan primary key: RFC 9562 UUIDv7 vs Sequence Block Allocation (Vitess).",
    content_items=[
        {"title": "RFC 9562 UUIDv7 (TIME-ORDERED 128-BIT)", "items": [
            "Struktur Bit: 48-bit Unix timestamp ms | 4-bit Ver (0111) | 12-bit rand_a | 2-bit Var (10) | 62-bit rand_b.",
            "Keunggulan B-Tree: Terurut secara temporal ($u_1 < u_2$); mencegah page splitting masif pada indeks database.",
            "Zero Coordination: Tidak memerlukan koordinasi jaringan, server sentral, atau lock terpusat antar worker."
        ], "type": "bullet"},
        {"title": "SEQUENCE BLOCK ALLOCATION (VITESS-STYLE)", "items": [
            "Mekanisme Blok: Worker meminta rentang blok ID (misal ukuran blok = 1.000 ID) dari central sequence counter.",
            "Eksekusi Lokal: 1.000 operasi insert berikutnya mengalokasikan ID sekuensial di memori lokal tanpa panggilan jaringan.",
            "Efisiensi Jaringan: Mengurangi 99.9% round-trip database dibandingkan query auto-increment sekuensial tradisional."
        ], "type": "bullet"}
    ],
    narrative="Dalam database terdistribusi, fitur auto-increment bawaan satu mesin tidak dapat digunakan karena akan menimbulkan tabrakan primary key antar shard. Lab ini menguji dua solusi standar industri. Pertama, RFC 9562 UUIDv7 menghasilkan ID 128-bit independen yang menyematkan 48-bit timestamp pada bit awal, menjamin sifat time-ordered yang ramah indeks B-tree tanpa koordinasi jaringan. Kedua, Sequence Block Allocation meminta jatah alokasi blok angka (misal blok 1.000 ID) ke server sentral, memungkinkan worker menerbitkan ID integer berurutan secara cepat di memori lokal.",
    narrative_en="Distributed databases cannot rely on local auto-increment sequences without triggering catastrophic primary key collisions across shards. This lab evaluates two industry standard solutions. First, RFC 9562 UUIDv7 generates 128-bit globally unique IDs embedding a 48-bit millisecond timestamp upfront, ensuring monotonic lexicographical sorting that prevents B-tree index fragmentation with zero network coordination. Second, Sequence Block Allocation reserves ID ranges in bulk from a central coordinator, enabling workers to dispense sequential integer IDs locally.",
    mental_model="UUIDv7 provides zero-coordination B-tree friendly IDs; sequence block allocators preserve sequential integers efficiently.",
    source_path="labs/32-database-sharding-and-partitioning/internal/idgen/idgen.go:1-120",
    slide_num="19 / 24",
    layout="cards"
)

# Slide 20: Checklist implementasi
add_slide(
    categories=["CHECKLIST", "PRODUKSI"],
    title="Checklist Implementasi Produksi Sharding & Partitioning",
    subtitle="9 poin verifikasi arsitektur sebelum menggelar cluster database terdistribusi ke produksi.",
    content_items=[
        {"title": "CHECKLIST KESIAPAN DISTRIBUTED DATABASE", "items": [
            "[ ] Mulai dari Logical Table Partitioning pada 1 mesin sebelum memutuskan sharding multi-node.",
            "[ ] Evaluasi Shard Key: Wajib berkardinalitas tinggi (tenant_id / user_id) dan tidak menggunakan nilai monotonik.",
            "[ ] Terapkan Hashed Sharding jika beban kerja didominasi oleh operasi tulis temporal intensif.",
            "[ ] Gunakan Consistent Hashing dengan 50–150 virtual nodes per server fisik untuk stabilitas scale-out.",
            "[ ] Siapkan Global Secondary Index (Lookup Vindex) untuk query non-shard key yang berfrekuensi tinggi.",
            "[ ] Terapkan timeout dan context cancellation ketat pada seluruh query fallback Scatter-Gather broadcast.",
            "[ ] Gunakan RFC 9562 UUIDv7 atau Sequence Block Allocation untuk primary key terdistribusi tanpa tabrakan.",
            "[ ] Kelompokkan tabel berelasi dalam satu shard key yang sama (Co-Location) demi menjaga integritas JOIN lokal.",
            "[ ] Rancang arsitektur tanpa ketergantungan pada transaksi Two-Phase Commit (2PC) cross-shard berkepanjangan."
        ], "type": "bullet"}
    ],
    narrative="Daftar periksa ini merangkum sembilan prinsip krusial yang wajib dievaluasi oleh arsitek database sebelum beralih ke arsitektur sharded. Langkah pertama selalu memaksimalkan logical partitioning dan optimasi query pada server tunggal. Ketika sharding fisik mutlak diperlukan, pemilihan sharding key berkardinalitas tinggi dan penggunaan consistent hash ring dengan virtual nodes adalah syarat mutlak. Kueri pencarian sekunder harus dilindungi oleh Global Secondary Index, dan struktur data relasional harus di-co-locate dalam shard yang sama demi menghindari latensi Two-Phase Commit.",
    narrative_en="This checklist codifies nine critical prerequisites before embarking on database sharding in production. Engineering teams must first exhaust single-node logical partitioning and hardware vertical headroom. When physical sharding becomes unavoidable, selecting high-cardinality shard keys and provisioning consistent hashing rings with virtual nodes are non-negotiable requirements. Secondary lookups must be secured via Global Secondary Indexes, and relational schemas must be co-located to avoid the crushing latency penalties of distributed 2PC transactions.",
    mental_model="Maximize single-node partitioning first; when sharding, mandate high-cardinality keys, vnode rings, and table co-location.",
    source_path="labs/32-database-sharding-and-partitioning/content/02-master-draft.md:335-350",
    slide_num="20 / 24",
    layout="checklist"
)

# Slide 21: Anti-pattern & cara reproduksi kegagalan
add_slide(
    categories=["ANTI-PATTERN", "REPRODUKSI"],
    title="Anti-Pattern: Kesalahan Desain Sharding & Reproduksi Lab",
    subtitle="Kekeliruan arsitektur yang melumpuhkan cluster dan cara membuktikannya langsung di lab.",
    content_items=[
        {"title": "ANTI-PATTERN DISTRIBUTED DATA UMUM", "items": [
            "Monotonic Shard Key: Menggunakan tanggal atau auto-increment sebagai shard key, memicu 100% write hotspot.",
            "Modulo Hash Resharding: Menggunakan rumus `key % N` yang memaksa 80% data berpindah saat penambahan node.",
            "Unindexed Scatter-Gather: Membiarkan query pencarian sekunder melakukan broadcast 100% node di setiap request.",
            "Distributed Foreign Keys: Mencoba menegakkan relasi JOIN antar shard berbeda tanpa co-location data."
        ], "type": "card"},
        {"title": "CARA REPRODUKSI KEGAGALAN DI LAB", "items": [
            "1. Eksekusi `demoShardingHotspots()` di `cmd/demo/main.go`: Saksikan 1.000 record menumpuk 100% di `shard-2`.",
            "2. Jalankan `TestRoutingAndConsistentHashRelocation`: Saksikan Hash Modulo memindahkan 75.12% key pada scale-out.",
            "3. Bandingkan latensi Scatter-Gather (92µs) vs GSI (1µs) di `demoScatterGatherVsGSI()` untuk membuktikan broadcast penalty."
        ], "type": "card"}
    ],
    narrative="Mengabaikan prinsip dasar partisi database dapat melumpuhkan seluruh infrastruktur backend. Anti-pattern paling umum adalah memilih sharding key monotonik seperti tanggal transaksi yang memusatkan seluruh write traffic pada satu mesin database yang malang. Kesalahan fatal lainnya adalah menggunakan rumus hash modulo sederhana yang memicu banjir data saat menambah kapasitas server. Pengembang dapat mereproduksi skenario kegagalan ini secara langsung di lab melalui fungsi demo dan unit test yang telah disediakan.",
    narrative_en="Violating fundamental distributed data principles introduces catastrophic performance failures. The most devastating anti-pattern remains adopting monotonic timestamp shard keys, concentrating 100% of cluster write throughput onto a single suffocating instance. Equally hazardous is relying on primitive modulo arithmetic that triggers massive data rebalancing storms during cluster expansions. Engineers can reproduce these exact degradation scenarios locally through the automated demo and benchmark suites provided in this lab.",
    mental_model="Never shard on monotonic keys; never scale out using raw modulo hashing; verify failure thresholds through lab tests.",
    source_path="labs/32-database-sharding-and-partitioning/content/02-master-draft.md:275-290",
    slide_num="21 / 24",
    layout="anti-pattern-cards"
)

# Slide 22: Takeaway praktis & command block dominan
add_slide(
    categories=["PRAKTIS", "COMMAND"],
    title="Verifikasi Mandiri: Uji & Jalankan Lab di Mesin Anda",
    subtitle="Reproduksi seluruh unit test, deteksi data race, dan simulasi cluster sharded di terminal Anda.",
    content_items=[
        {"title": "PERINTAH VERIFIKASI RESMI", "items": [
            "cd labs/32-database-sharding-and-partitioning",
            "go test -v -count=1 ./...",
            "go test -race -v -count=1 ./...",
            "go run ./cmd/demo"
        ], "type": "code"},
        {"title": "RINGKASAN METRIK LABORATORIUM", "items": [
            "Test Suite: 5/5 test suites lulus 100% (Pruning, Resharding Relocation, Scatter-Gather & GSI, UUIDv7, Concurrency).",
            "Race Detector: Bersih mutlak dari data race condition pada 200 goroutine simultan (1.405s pass).",
            "Zero Dependencies: Murni Go Standard Library (`sync`, `hash/fnv`, `crypto/rand`, `sort`, `time`) tanpa dependensi eksternal."
        ], "type": "bullet"}
    ],
    narrative="Laboratorium ini dirancang untuk dapat diverifikasi secara mandiri dalam hitungan detik tanpa memerlukan instalasi database eksternal. Seluruh logika sharding router, virtual nodes, partition pruning, dan generator UUIDv7 dibangun murni menggunakan pustaka standar Go. Perintah go test -v mengeksekusi pengujian komprehensif mulai dari verifikasi query pruning hingga rasio relokasi consistent hashing. Perintah go test -race menjamin thread-safety di bawah konkurensi tinggi, dan go run demo memperlihatkan perbandingan performa secara visual.",
    narrative_en="This laboratory is engineered for immediate, deterministic peer verification without running external database containers or daemons. All consistent hashing algorithms, virtual node rings, partition pruning evaluators, and UUIDv7 generators are built strictly with Go standard library primitives. Running go test -v executes rigorous test suites verifying pruning bounds and scale-out relocation ratios. Executing go test -race confirms thread safety under high concurrency, while go run demo renders the full interactive benchmark.",
    mental_model="Zero-dependency standard library code enables deterministic, transparent verification of distributed data architectures.",
    source_path="labs/32-database-sharding-and-partitioning/README.md • engineering/03-execution-result.md",
    slide_num="22 / 24",
    layout="command-block"
)

# Slide 23: Prinsip arsitektural jangka panjang
add_slide(
    categories=["PRINSIP", "ARSITEKTUR"],
    title="Prinsip Arsitektur: Desain Database Skala Besar",
    subtitle="Tiga hukum fundamental dalam merancang penyimpanan data yang tahan uji skala waktu.",
    content_items=[
        {"title": "1. PARTISI LOGIS SEBELUM SHARDING FISIK", "items": [
            "Gunakan declarative table partitioning pada satu mesin untuk query pruning dan manajemen siklus hidup data.",
            "Tunda sharding fisik terdistribusi sampai batasan hardware vertikal (I/O, RAM, Write lock) benar-benar terlampaui."
        ], "type": "card"},
        {"title": "2. SHARD KEY MENENTUKAN STRUKTUR SISTEM", "items": [
            "Pilih sharding key berkardinalitas tinggi yang menyelaraskan pola akses baca dan tulis aplikasi.",
            "Co-locate entitas data terkait di bawah satu kunci partisi guna menghindari transaksi terdistribusi Two-Phase Commit."
        ], "type": "card"},
        {"title": "3. MINIMALKAN PERPINDAHAN & BROADCAST", "items": [
            "Gunakan Consistent Hashing dengan virtual nodes agar ekspansi cluster hanya memindahkan porsi data minimal $O(1/N)$.",
            "Lindungi sistem dari broadcast scatter-gather menggunakan Global Secondary Indexes untuk akses atribut sekunder."
        ], "type": "card"}
    ],
    narrative="Tiga prinsip arsitektural ini merupakan pedoman abadi bagi arsitek perangkat lunak dalam merancang sistem penyimpanan data. Pertama, jangan terburu-buru melakukan sharding terdistribusi; maksimalkan partisi logis dan optimasi query pada mesin tunggal sebelum menambah kompleksitas jaringan. Kedua, pahami bahwa sharding key adalah keputusan topologi yang menentukan struktur sistem: co-locate data terkait untuk mempertahankan transaksi dan relasi lokal. Ketiga, rancang arsitektur routing yang meminimalkan pergerakan data saat scaling dan mengeliminasi ketergantungan pada query broadcast.",
    narrative_en="These three timeless architectural principles govern high-scale data layer design. First, resist premature distributed sharding; exhaust single-node logical partitioning, indexing, and vertical headroom before introducing network boundaries. Second, recognize that shard key selection defines the topology of your system: co-locate related records to preserve local transactional boundaries and avoid distributed commits. Third, build routing infrastructures that minimize data relocation during scale-outs and systematically eliminate scatter-gather broadcast dependencies.",
    mental_model="Exhaust single-node partitioning first; co-locate data by high-cardinality keys; govern scale-out with consistent hash rings.",
    source_path="labs/32-database-sharding-and-partitioning/content/05-key-takeaways.md",
    slide_num="23 / 24",
    layout="principles-list"
)

# Slide 24: Closing slide
add_slide(
    categories=["CLOSING"],
    title="Partisi Cerdas, Sharding Terukur, Skalabilitas Tanpa Batas.",
    subtitle="",
    content_items=[
        {"title": "KOMPONEN ARSITEKTUR LABORATORIUM", "items": [
            "internal/partitioning/table.go — Single-node range-based logical partitioning with query pruning & fast drop",
            "internal/sharding/sharding.go — ModuloRouter vs ConsistentHashRouter (with vnodes), Cluster, & Scatter-Gather",
            "internal/idgen/idgen.go — RFC 9562 UUIDv7 (128-bit time-ordered) & Vitess-style Sequence Block Allocator",
            "tests/sharding_test.go — 5 comprehensive test suites covering pruning, resharding, GSI, IDs, & concurrency",
            "cmd/demo/main.go — Interactive benchmarks demonstrating hotspots, resharding churn, and 92x GSI speedup"
        ], "type": "bullet"},
        {"title": "VERIFIKASI & REPOSITORI", "items": [
            f"Commit SHA: {COMMIT_SHA[:8]}",
            "Unit Tests: 5/5 PASS (100% clean, 0.446s)",
            "Race Detector: PASS (zero data race, 1.405s)",
            f"Repositori: github.com/lukman-ss/software-engineering-lab/tree/main/labs/{LAB_NO}-{SLUG}"
        ], "type": "bullet"}
    ],
    narrative="Laboratorium Database Sharding and Partitioning ini membuktikan secara konkret bahwa skalabilitas penyimpanan terdistribusi dapat dicapai dengan elegan melalui pemisahan tanggung jawab yang jelas. Dengan memadukan logical table partitioning untuk efisiensi siklus hidup lokal, router consistent hashing ring dengan virtual nodes untuk ekspansi minim migrasi, serta Global Secondary Index untuk mengeliminasi broadcast scatter-gather, sistem mampu melampaui batasan mesin tunggal dengan performa terprediksi. Seluruh kode sumber, dokumentasi arsitektur, dan hasil uji tersedia lengkap di repositori.",
    narrative_en="This Database Sharding and Partitioning laboratory provides empirical proof that distributed data scalability is cleanly attainable through disciplined architectural boundaries. By combining logical table partitioning for local query pruning, consistent hashing rings with virtual nodes for minimal resharding shuffle, and Global Secondary Indexes to eliminate scatter-gather broadcast storms, backend systems scale beyond single-node ceilings with deterministic performance. Complete source code and verification harnesses are available in the repository.",
    mental_model="Logical partitioning prunes locally; consistent hashing scales horizontally; GSI eliminates broadcast overhead.",
    source_path=REPO_URL,
    slide_num="24 / 24",
    layout="closing-dark-navy"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
