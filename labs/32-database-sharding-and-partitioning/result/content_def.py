#!/usr/bin/env python3
"""
Publication Builder Definition for Software Engineering Lab 32: Database Sharding and Partitioning.
Defines metadata, colors, 24 slides, and mental models.
"""
import os
import sys

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/32-database-sharding-and-partitioning/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/32-database-sharding-and-partitioning"
COMMIT_SHA = "cfba3b76824f0ee03b5fa526462d1b72858144e7"
LAB_NO = "32"
SLUG = "database-sharding-and-partitioning"
TITLE_ID = "Database Sharding & Partitioning"
SUBTITLE_ID = "Skalabilitas Horizontal Melalui Routing, Partition Pruning, dan Consistent Hashing"
TITLE_EN = "Database Sharding & Table Partitioning"
SUBTITLE_EN = "Horizontal Scaling via Routing, Partition Pruning, and Consistent Hashing"
AUTHOR = "Lukman / lukman-ss"

TOPIC_TAGS = "Database Sharding • Table Partitioning • Consistent Hashing • Software Engineering Lab"

# Brand Color Constants
DARK_NAVY = "0b1b2d"
TEAL = "087f76"
CREAM = "f4f1e8"
WHITE = "ffffff"
SLATE = "52606b"
ORANGE = "d97706"

# Slide data structure
SLIDES = []

def add_slide(categories, title, subtitle, content_items, mental_model, source_path="", slide_num="", narrative="", title_en="", narrative_en=""):
    """Add a slide definition to the deck."""
    SLIDES.append({
        "categories": categories,
        "title": title,
        "subtitle": subtitle,
        "content_items": content_items,
        "mental_model": mental_model,
        "source": source_path,
        "slide_num": slide_num,
        "narrative": narrative,
        "title_en": title_en or title,
        "narrative_en": narrative_en or narrative
    })

# =====================================================================
# 24 SLIDES DEFINITIONS
# =====================================================================

# Slide 1: Cover (Dark Navy #0b1320 with visual CRT elements)
add_slide(
    categories=["KONTEKS"],
    title=TITLE_ID,
    subtitle="SOFTWARE ENGINEERING LAB | LAB " + LAB_NO,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "Retro CRT monitor visual with scanlines",
            "Hash routing ring & partition boundary terminal display",
            f"Penulis: {AUTHOR}",
        ]}
    ],
    mental_model="",
    source_path="",
    slide_num=""
)

# Slide 2: Batas Skalabilitas Monolitik
add_slide(
    categories=["MOTIVASI", "BATAS FISIK"],
    title="Batas Fisik Database Single-Instance",
    subtitle="Tiga batasan hardware memaksa transisi dari vertical scaling ke partitioning dan sharding",
    content_items=[
        {"title": "BATAS MESIN TUNGGAL", "items": [
            "Batas RAM: Buffer pool cache tidak lagi menampung working dataset aktif",
            "Batas Disk IOPS: Saturasi write throughput dan antrean fsync log WAL",
            "Lock Contention: Transaksi konkuren berebut lock pada index tree dan page yang sama",
        ], "type": "bullet"},
        {"title": "BIAYA VERTICAL SCALING", "items": [
            "Biaya hardware tumbuh eksponensial untuk peningkatan kapasitas marginal",
            "Single point of failure: Kegagalan hardware host menghentikan seluruh layanan",
            "Batas arsitektural: Bus memori dan jumlah core CPU memiliki batas fisik",
        ], "type": "bullet"},
        {"title": "DUA JALUR SKALASI", "items": [
            "Logical Partitioning: Memecah tabel internal pada satu engine database yang sama",
            "Physical Sharding: Membagi dataset ke beberapa independent database nodes",
            "Kombinasi keduanya menjawab skalabilitas lifecycle dan throughput transaksi",
        ], "type": "bullet"},
    ],
    mental_model="Vertical scaling hits a hard hardware ceiling; horizontal partitioning distributes the blast radius.",
    source_path="labs/32-database-sharding-and-partitioning/README.md",
    slide_num="02 / 24",
    narrative="Sistem penyimpanan monolitik menghadapi batas fisik tak terelakkan ketika volume data dan laju transaksi melonjak. Peningkatan kapasitas server tunggal menemui titik jenuh pada batasan I/O storage, memori buffer pool, dan lock contention internal. Memahami perbedaan antara pemisahan logis di dalam satu mesin dan pemisahan fisik ke kluster multi-node adalah langkah pertama menuju arsitektur data terdistribusi yang berkelanjutan.",
    title_en="Physical Limits of Single-Instance Databases",
    narrative_en="Monolithic database engines inevitably hit physical hardware limits as data volume and transaction rates surge. Scaling up a single host encounters severe disk I/O saturation, memory buffer pool misses, and internal lock contention. Distinguishing internal logical table partitioning from multi-node physical sharding establishes the baseline for scalable data architectures."
)

# Slide 3: Logical Partitioning vs Physical Sharding
add_slide(
    categories=["ARSITEKTUR", "TAKSONOMI"],
    title="Logical Partitioning vs Physical Sharding",
    subtitle="Dua paradigma berbeda dengan trade-off operasional, isolasi, dan koordinasi yang kontras",
    content_items=[
        {"title": "PERBANDINGAN ARSITEKTUR", "items": [
            "Dimensi | Logical Partitioning | Physical Sharding",
            "Instance Host | 1 Database Server | N Node Server Terpisah",
            "Batas Skalasi | Batas Hardware Host | Tak Terbatas (Scale-out)",
            "Cross-Join | Native SQL Engine | Memerlukan Scatter-Gather",
            "Distributed ID | Auto-Increment Lokal | Sequence Generator / UUIDv7",
        ], "type": "table"},
        {"title": "LOGICAL PARTITIONING", "items": [
            "Dikelola oleh query engine database internal (PostgreSQL/MySQL)",
            "Tabel induk logis dipecah menjadi sub-tabel fisik terpisah",
            "Optimal untuk manajemen lifecycle data dan query pruning",
        ], "type": "bullet"},
        {"title": "PHYSICAL SHARDING", "items": [
            "Setiap shard adalah database instance independen tanpa shared state",
            "Membutuhkan layer router pintar untuk mengarahkan transaksi",
            "Menghilangkan single point of failure pada level penyimpanan data",
        ], "type": "bullet"},
    ],
    mental_model="Partitioning optimizes single-node query pruning; sharding breaks single-node resource limits.",
    source_path="labs/32-database-sharding-and-partitioning/content/02-master-draft.md",
    slide_num="03 / 24",
    narrative="Logical partitioning membagi satu tabel besar menjadi beberapa sub-tabel di server yang sama, memungkinkan query engine memindai hanya rentang data yang relevan. Sebaliknya, physical sharding memecah dataset ke beberapa server independen yang memiliki CPU, RAM, dan disk tersendiri. Sharding menuntut koordinasi routing dan generator ID terdistribusi, tetapi membuka jalan untuk skalabilitas horizontal linear melampaui kapasitas mesin mana pun.",
    title_en="Logical Partitioning vs Physical Sharding",
    narrative_en="Logical partitioning splits a single table into sub-tables within the same engine, empowering the optimizer to scan only relevant subsets. In contrast, physical sharding distributes records across isolated instances, each with dedicated CPU, RAM, and disk. Sharding requires routing middleware and distributed ID generation, but unlocks horizontal scalability far beyond single-machine limits."
)

# Slide 4: Arsitektur Single-Node Range Partitioning
add_slide(
    categories=["DESAIN ENGINE", "PARTITIONING"],
    title="Mekanisme Range Partitioning pada Single Node",
    subtitle="Sub-tabel independen berdasarkan rentang temporal waktu transaksi",
    content_items=[
        {"title": "STRUKTUR PARTISI TEMPORAL", "items": [
            "sales_logs (Tabel Induk Logis)",
            "├── p2026_q1: [2026-01-01 s.d. 2026-04-01)",
            "├── p2026_q2: [2026-04-01 s.d. 2026-07-01)",
            "├── p2026_q3: [2026-07-01 s.d. 2026-10-01)",
            "└── p2026_q4: [2026-10-01 s.d. 2027-01-01)",
        ], "type": "diagram"},
        {"title": "DESAIN KODE GO", "items": [
            "type PartitionRange struct {",
            "    Name  string",
            "    Start time.Time",
            "    End   time.Time",
            "}",
            "type Table struct { partitions []*Partition }",
        ], "type": "code"},
        {"title": "MANFAAT OPERASIONAL", "items": [
            "Tiap partisi memiliki index B-tree terpisah yang muat di RAM",
            "Query berbasis tanggal secara deterministik diarahkan ke sub-tabel target",
            "Menghindari degradasi index tree raksasa yang membengkak",
        ], "type": "bullet"},
    ],
    mental_model="Time-bounded partition ranges keep working B-tree indices compact enough to fit in RAM.",
    source_path="labs/32-database-sharding-and-partitioning/internal/partitioning/table.go:14",
    slide_num="04 / 24",
    narrative="Dalam implementasi range partitioning, tabel logis sales_logs dipetakan ke partisi kuartalan yang terisolasi. Setiap partisi hanya menerima record yang nilai timestamp-nya berada dalam interval waktu yang ditentukan. Pemisahan ini menjaga agar struktur index B-tree masing-masing partisi tetap kompak dan tetap berada di dalam memori cache, mencegah degradasi performa yang lazim terjadi pada index tabel tunggal berukuran miliaran baris.",
    title_en="Single-Node Range Partitioning Architecture",
    narrative_en="In range partitioning implementations, the logical sales_logs table is mapped across quarterly isolated partitions. Each partition accepts only records whose timestamp falls strictly within its bounds. This containment ensures each partition's B-tree index remains compact and cache-resident, eliminating the cache thrashing common to multi-billion-row monolithic indexes."
)

# Slide 5: Partition Pruning: Eliminasi I/O Sia-Sia
add_slide(
    categories=["OPTIMASI QUERY", "BUKTI EKSEKUSI"],
    title="Partition Pruning: Eliminasi Pindaian Partisi",
    subtitle="Query optimizer memotong partisi di luar rentang filter filter WHERE",
    content_items=[
        {"title": "ALGORITMA PRUNING", "items": [
            "Query filter: WHERE created_at BETWEEN '2026-04-01' AND '2026-06-30'",
            "Evaluasi rentang interval: Hanya p2026_q2 yang beririsan",
            "Engine mengabaikan p2026_q1, p2026_q3, dan p2026_q4 secara mutlak",
        ], "type": "bullet"},
        {"title": "HASIL VERIFIKASI DEMO", "items": [
            "Query Range: 2026-04-01 to 2026-06-30",
            "Partitions Scanned : 1 / 4",
            "Partitions Pruned  : 3 / 4 (75% I/O dihemat)",
            "Records Found      : 1 record tepat",
        ], "type": "code"},
        {"title": "KONDISI PRUNING GAGAL", "items": [
            "Penggunaan fungsi pada kolom partisi: WHERE DATE(created_at) = ...",
            "Klausa WHERE tanpa mencantumkan kunci partisi memicu full table scan",
            "Tipe data yang tidak presisi memaksa implicit type cast pada runtime",
        ], "type": "bullet"},
    ],
    mental_model="Partition pruning turns sequential scans into targeted sub-tree lookups by checking range boundaries first.",
    source_path="labs/32-database-sharding-and-partitioning/internal/partitioning/table.go:61",
    slide_num="05 / 24",
    narrative="Partition pruning adalah kemampuan query engine untuk mengeliminasi pemindaian sub-tabel yang tidak beririsan dengan predikat WHERE. Pada uji coba lab, rentang tanggal Q2 membuat database hanya memindai 1 dari 4 partisi yang ada, langsung memotong 75 persen pemindaian disk. Namun, pruning akan gagal total jika query membungkus kolom partisi dengan fungsi ekspresi atau jika predikat tidak menyertakan kunci partisi sama sekali.",
    title_en="Partition Pruning: Eliminating Wasteful Disk Scans",
    narrative_en="Partition pruning allows the query optimizer to completely bypass sub-tables that do not intersect with WHERE clauses. In our verified test, a Q2 date-range query scans exactly 1 of 4 partitions, pruning 75 percent of the dataset's I/O overhead. However, pruning fails if queries wrap partition keys inside expressions or omit the partition key entirely from the predicate."
)

# Slide 6: Bulk Lifecycle: Fast Drop Partition vs DELETE
add_slide(
    categories=["OPERASIONAL DB", "LIFECYCLE"],
    title="Drop Partition Instan vs Expensive DELETE",
    subtitle="Operasi metadata O(1) menggantikan jutaan tombstone row dan overhead VACUUM",
    content_items=[
        {"title": "PERBANDINGAN METODE", "items": [
            "Operasi | DELETE FROM ... | DROP PARTITION",
            "Waktu Eksekusi | O(N) Berjam-jam | O(1) Milidetik",
            "Write Amplification | Sangat Tinggi (WAL Row) | Minimum (Metadata File)",
            "Disk Space Reclaim | Butuh VACUUM FULL | Langsung ke File System",
            "Lock Impact | Row Lock Berkepanjangan | Exclusive Lock Singkat",
        ], "type": "table"},
        {"title": "KODE DROP PARTISI", "items": [
            "func (t *Table) DropPartition(name string) bool {",
            "    t.mu.Lock(); defer t.mu.Unlock()",
            "    // Hapus referensi metadata dan bebaskan sub-tabel",
            "    // O(1) unlinking tanpa memindai isi baris",
            "}",
        ], "type": "code"},
        {"title": "IMPLIKASI PRODUKSI", "items": [
            "Data retensi audit/transaksi 90 hari dapat dihapus secara instan",
            "Mencegah tabel membengkak (bloat) yang memperlambat sequential scan",
            "Tidak membebani worker background autovacuum pada database aktif",
        ], "type": "bullet"},
    ],
    mental_model="Dropping a partition unlinks files at the filesystem layer; deleting rows burns write I/O per tombstone.",
    source_path="labs/32-database-sharding-and-partitioning/internal/partitioning/table.go:88",
    slide_num="06 / 24",
    narrative="Menghapus data kadaluwarsa pada tabel monolitik raksasa menggunakan klausa DELETE menghasilkan jutaan catatan transaksi WAL dan tombstone, yang berujung pada table bloat dan kebutuhan VACUUM berat. Dengan table partitioning, data lama dibuang menggunakan perintah DROP PARTITION yang mengeksekusi operasi unlinking file O(1). Ruang penyimpanan langsung kembali ke sistem operasi dalam hitungan milidetik tanpa mengganggu transaksi produksi yang sedang berlangsung.",
    title_en="Instant Partition Drop vs Heavy Row Deletion",
    narrative_en="Purging historical data with bulk DELETE statements generates millions of WAL records and dead row tombstones, requiring aggressive VACUUMing and inducing severe table bloat. With table partitioning, obsolete data is discarded via DROP PARTITION, an O(1) metadata unlinking operation that reclaims disk blocks instantly without saturating transaction write logs."
)

# Slide 7: Anatomi Arsitektur Physical Sharding
add_slide(
    categories=["ARSITEKTUR SHARDING", "KOMPONEN"],
    title="Anatomi Kluster Sharding Fisik",
    subtitle="Arsitektur shared-nothing: Routing layer, Shard nodes, dan Metadata catalog",
    content_items=[
        {"title": "KOMPONEN KLUSTER", "items": [
            "Aplikasi Klien: Mengirim query SQL tanpa perlu tahu topologi node",
            "Routing Coordinator: Menerjemahkan Shard Key ke Shard Node target",
            "Shard Nodes: Server PostgreSQL/MySQL independen (Shared-Nothing)",
            "Catalog / GSI: Menyimpan pemetaan metadata dan indeks sekunder",
        ], "type": "bullet"},
        {"title": "INTERFACES KUNCI DI GO", "items": [
            "type Router interface {",
            "    GetShard(key string) (string, error)",
            "    AddShard(shardID string)",
            "    RemoveShard(shardID string)",
            "}",
            "type Cluster struct { shards map[string]*Shard }",
        ], "type": "code"},
        {"title": "SIFAT SHARED-NOTHING", "items": [
            "Tidak ada memori atau disk yang dibagi antar server shard",
            "Kegagalan satu node hanya berdampak pada subset data terkait",
            "Kapasitas komputasi dan write I/O bertambah secara linear per node",
        ], "type": "bullet"},
    ],
    mental_model="Shared-nothing sharding separates state completely; routers translate logical keys into physical node targets.",
    source_path="labs/32-database-sharding-and-partitioning/internal/sharding/sharding.go:23",
    slide_num="07 / 24",
    narrative="Arsitektur physical sharding menganut prinsip shared-nothing di mana setiap server database beroperasi mandiri tanpa berbagi memori maupun disk storage. Lapisan koordinator routing bertindak sebagai perantara yang mengevaluasi shard key pada query untuk menemukan shard tujuan secara deterministik. Desain ini mengisolasi kegagalan node dan memastikan penambahan node baru benar-benar menambah kapasitas total sistem secara proporsional.",
    title_en="Anatomy of a Physical Sharding Cluster",
    narrative_en="Physical sharding follows a shared-nothing paradigm where database instances operate without shared memory or disk pools. A routing coordination layer intercepts application queries, inspects the shard key, and routes requests to the authoritative shard node. This isolation contains infrastructure failures and provides horizontal scalability with every additional machine."
)

# Slide 8: Bahaya Fatal Monotonic Shard Key
add_slide(
    categories=["FAILURE SCENARIO", "ANTI-PATTERN"],
    title="Write Hotspot Akibat Monotonic Shard Key",
    subtitle="Kunci sekuensial mengalirkan 100% traffic penulisan ke satu node aktif",
    content_items=[
        {"title": "PENYEBAB HOTSPOT", "items": [
            "Penggunaan timestamp (created_at) atau autoincrement ID sebagai shard key",
            "Semua transaksi baru pada jendela waktu yang sama memiliki nilai kunci identik",
            "Router mengarahkan seluruh traffic write ke satu server fisik yang sama",
        ], "type": "bullet"},
        {"title": "EVIDENSI PENGUJIAN LAB", "items": [
            "Ingest 1.000 transaksi berurutan dengan kunci tanggal:",
            "  Node shard-0 :    0 records (0%)",
            "  Node shard-1 :    0 records (0%)",
            "  Node shard-2 : 1000 records (100% [##########])",
            "  Node shard-3 :    0 records (0%)",
            "Hasil: Write Hotspot Parah! Shard lain menganggur.",
        ], "type": "code"},
        {"title": "DAMPAK SISTEMIK", "items": [
            "Satu node kehabisan disk I/O dan connection pool",
            "Node lain di kluster mengalami idle resource pemborosan biaya",
            "Kluster multi-node berkinerja lebih buruk daripada single-instance",
        ], "type": "bullet"},
    ],
    mental_model="Sequential keys serialize writes into a single shard, turning multi-node clusters into single-node bottlenecks.",
    source_path="labs/32-database-sharding-and-partitioning/cmd/demo/main.go:52",
    slide_num="08 / 24",
    narrative="Memilih kunci yang bertambah secara monotonik seperti tanggal atau auto-increment integer adalah jebakan paling berbahaya dalam arsitektur sharding. Pengujian lab menunjukkan bahwa 1.000 transaksi berturut-turut pada hari yang sama seluruhnya bermuara pada shard-2, sementara node lainnya menganggur. Alih-alih mendapatkan distribusi beban merata, sistem justru menciptakan bottleneck parah di satu server sementara sumber daya node lain terbuang sia-sia.",
    title_en="Write Hotspot Catastrophe from Monotonic Shard Keys",
    narrative_en="Selecting monotonically increasing keys such as created_at timestamps or sequential integers is a classic sharding anti-pattern. Lab tests confirm that ingesting 1,000 records on the same date dumps 100 percent of write operations onto shard-2, leaving the remaining nodes completely idle. The cluster degenerates into a single bottlenecked host with massive wasted infrastructure."
)

# Slide 9: Solusi High-Cardinality Sharding Key
add_slide(
    categories=["SOLUSI DESAIN", "DISTRIBUSI DATA"],
    title="Kunci Berkardinalitas Tinggi untuk Distribusi Merata",
    subtitle="Distribusi seimbang di seluruh node fisik dengan tenant_id atau user_id",
    content_items=[
        {"title": "KRITERIA KUNCI IDEAL", "items": [
            "Kardinalitas Tinggi: Memiliki jutaan nilai distinct (user_id, account_id)",
            "Frekuensi Rendah: Tidak ada satu entitas yang menguasai >5% transaksi",
            "Non-Monotonik: Pola penulisan acak di sepanjang ruang domain kunci",
        ], "type": "bullet"},
        {"title": "HASIL VERIFIKASI LAB", "items": [
            "Ingest 1.000 data menggunakan kunci tenant_id:",
            "  Node shard-0 : 400 records (40%)",
            "  Node shard-2 : 200 records (20%)",
            "  Node shard-3 : 400 records (40%)",
            "Distribusi menyebar ke seluruh node fisik secara proporsional.",
        ], "type": "code"},
        {"title": "HASHED SHARDING STRATEGY", "items": [
            "Jika entitas bisnis monotonik, bungkus dengan fungsi hash (MD5, Murmur, FNV)",
            "Hash mengacak urutan temporal menjadi distribusi pseudorandom seragam",
            "Kelemahan: Query berbasis range pada kunci asli tidak lagi terlokalisasi",
        ], "type": "bullet"},
    ],
    mental_model="High-cardinality keys spread transactions uniformly; hashing scrambles temporal patterns into random scatter.",
    source_path="labs/32-database-sharding-and-partitioning/internal/sharding/sharding.go:75",
    slide_num="09 / 24",
    narrative="Kunci sharding yang baik wajib memiliki kardinalitas tinggi dan distribusi nilai acak, seperti user_id atau tenant_id. Dengan mendistribusikan data berdasarkan identifier ini, beban penulisan tersebar merata di seluruh node kluster. Jika aplikasi terpaksa menggunakan kunci berurutan, pendekatan hashed sharding diterapkan untuk mengacak nilai kunci melalui fungsi kriptografis atau hash non-kriptografis sebelum menentukan nomor shard tujuan.",
    title_en="High-Cardinality Shard Keys for Uniform Distribution",
    narrative_en="Effective sharding keys demand high cardinality and non-monotonic distribution, such as tenant_id or user_id. Hashing or picking high-cardinality values scatters incoming writes evenly across all available nodes. When natural business keys are sequential, hashed sharding computes a pseudo-random digest before node mapping, trading away range scans for uniform load distribution."
)

# Slide 10: Routing Sederhana: Hash Modulo
add_slide(
    categories=["ROUTING ALGORITHM", "MODULO"],
    title="Routing Hash Modulo: Rumus Naif & Sederhana",
    subtitle="Pemetaan langsung hash(key) % N: Cepat dihitung, rapuh saat cluster resize",
    content_items=[
        {"title": "RUMUS MATEMATIS", "items": [
            "ShardID = hash(key) % N",
            "Di mana N adalah jumlah node fisik aktif di kluster",
            "Contoh: hash('tenant_12') = 982312, 982312 % 4 = Shard 0",
        ], "type": "bullet"},
        {"title": "IMPLEMENTASI KODE GO", "items": [
            "func (m *ModuloRouter) GetShard(key string) (string, error) {",
            "    m.mu.RLock(); defer m.mu.RUnlock()",
            "    h := hashKey(key)",
            "    idx := int(h % uint64(len(m.shards)))",
            "    return m.shards[idx], nil",
            "}",
        ], "type": "code"},
        {"title": "KELEBIHAN & KETERBATASAN", "items": [
            "Kelebihan: Waktu pencarian O(1) konstan tanpa struktur data tambahan",
            "Kelebihan: Distribusi kunci sangat merata selama hash function baik",
            "Kelemahan: Mengubah nilai N merusak pemetaan hampir seluruh data di kluster",
        ], "type": "bullet"},
    ],
    mental_model="Hash modulo delivers instant O(1) routing, but ties node count directly into the key-to-node mapping.",
    source_path="labs/32-database-sharding-and-partitioning/internal/sharding/sharding.go:30",
    slide_num="10 / 24",
    narrative="Metode Hash Modulo membagi nilai hash kunci dengan jumlah total node N menggunakan operator sisa bagi. Pendekatan ini sangat sederhana dan menghasilkan lookup berkecepatan O(1) tanpa memori tambahan. Namun, kelemahannya terletak pada keterikatan nilai N dalam rumus; saat satu node ditambahkan atau dikurangi, hampir semua kunci lama akan menghasilkan indeks baru yang berbeda dari lokasi fisiknya saat ini.",
    title_en="Hash Modulo Routing: The Naive Division Approach",
    narrative_en="Hash Modulo routes records by taking the modulus of a hashed key against the total node count N. While trivial to implement and executing in O(1) time without extra memory overhead, it binds physical node topology directly to data addressing. Modifying N alters the computed destination for the vast majority of previously assigned keys."
)

# Slide 11: Bencana Skalasi Hash Modulo
add_slide(
    categories=["SKALASI & RESIZE", "EVIDENSI DATA"],
    title="Bencana Resharding Hash Modulo",
    subtitle="Menambah satu node memicu relokasi masif ~80% data di seluruh sistem",
    content_items=[
        {"title": "PERSENTASE TEORITIS RELOKASI", "items": [
            "Rasio data yang berpindah saat scale-out dari N ke N+1:",
            "  Moved Ratio = N / (N + 1)",
            "  Dari 4 node ke 5 node = 4/5 = 80.00% kunci wajib dipindahkan!",
        ], "type": "bullet"},
        {"title": "BUKTI KUANTITATIF LAB (10.000 KEYS)", "items": [
            "Cluster Resize: 4 Shards -> 5 Shards",
            "Keys Remapped   : 7.984 / 10.000 keys",
            "Persentase Pindah: 79.84% (Sangat akurat dengan teori 80%)",
            "Kunci yang bertahan di tempat: Hanya 20.16%",
        ], "type": "code"},
        {"title": "KONSEKUENSI PRODUKSI", "items": [
            "Saturasi bandwidth jaringan antar database untuk transfer data",
            "Periode read-only atau locking berkepanjangan pada aplikasi",
            "Tingkat risiko kegagalan migrasi dan inkonsistensi data sangat tinggi",
        ], "type": "bullet"},
    ],
    mental_model="Scaling hash modulo clusters reshuffles N/(N+1) of all existing data, saturating networks during rebalance.",
    source_path="labs/32-database-sharding-and-partitioning/cmd/demo/main.go:89",
    slide_num="11 / 24",
    narrative="Kelemahan fatal Hash Modulo terungkap saat kluster melakukan scale-out. Secara teoritis, menambah node dari 4 menjadi 5 akan memindahkan N/(N+1) atau 80% data. Hasil uji coba lab dengan 10.000 record membuktikan bahwa 7.984 kunci (79,84%) berubah shard. Di lingkungan produksi, pemindahan empat per lima data ini menimbulkan beban jaringan masif, lonjakan latensi, dan risiko downtime yang tidak dapat diterima.",
    title_en="The Hash Modulo Resharding Catastrophe",
    narrative_en="The fatal flaw of Hash Modulo surfaces during cluster expansion. Adding a single node from N to N+1 displaces N/(N+1) of the entire dataset. In our lab simulation with 10,000 keys, 7,984 records (79.84%) shifted nodes upon adding the 5th shard. In production, reshuffling 80 percent of records floods storage networks and risks severe system downtime."
)

# Slide 12: Solusi: Consistent Hashing Ring
add_slide(
    categories=["ALGORITMA", "CONSISTENT HASHING"],
    title="Consistent Hashing: Memutus Ketergantungan N",
    subtitle="Ruang hash sirkular melingkar 64-bit membatasi relokasi data hanya pada 1/N",
    content_items=[
        {"title": "KONSEP RING SIRKULAR", "items": [
            "Ruang kunci dipetakan ke lingkaran angka 0 hingga 2^64 - 1",
            "Setiap node fisik ditempatkan pada titik tertentu di sepanjang cincin",
            "Kunci data berjalan searah jarum jam hingga bertemu node pertama",
        ], "type": "bullet"},
        {"title": "PENCARIAN BINARY SEARCH", "items": [
            "func (ch *ConsistentHashRouter) GetShard(key string) (string, error) {",
            "    h := hashKey(key)",
            "    idx := sort.Search(len(ch.ring), func(i int) bool {",
            "        return ch.ring[i].hash >= h",
            "    })",
            "    if idx == len(ch.ring) { idx = 0 }",
            "    return ch.ring[idx].shardID, nil",
            "}",
        ], "type": "code"},
        {"title": "MINIMAL RELOKASI DATA", "items": [
            "Saat node baru disisipkan, ia hanya mengambil data dari tetangga terdekat",
            "Hanya ~1/N data yang berpindah lokasi secara keseluruhan",
            "Node lain yang tidak bertetangga sama sekali tidak terpengaruh",
        ], "type": "bullet"},
    ],
    mental_model="Consistent hashing places nodes and keys on a circular ring; new nodes only capture keys from immediate neighbors.",
    source_path="labs/32-database-sharding-and-partitioning/internal/sharding/sharding.go:143",
    slide_num="12 / 24",
    narrative="Consistent Hashing memecahkan masalah redistribusi data dengan memetakan node dan kunci data ke dalam cincin sirkular berukuran 64-bit. Kunci diarahkan ke node terdekat searah jarum jam menggunakan algoritma pencarian biner sort.Search. Saat node baru ditambahkan, node tersebut hanya mengambil sebagian beban dari node tetangga terdekatnya di cincin, menjaga agar sisa data di node lain tetap utuh tanpa perpindahan.",
    title_en="Consistent Hashing: Breaking the N-Dependence",
    narrative_en="Consistent hashing maps both server nodes and dataset keys onto a 64-bit circular ring. Lookups locate the nearest clockwise node using binary search across the ring array. Adding a new node only intercepts keys from its immediate predecessor, preserving existing shard assignments across the rest of the cluster and eliminating widespread rebalancing."
)

# Slide 13: Peran Virtual Nodes (vnodes)
add_slide(
    categories=["REKAYASA SISTEM", "VIRTUAL NODES"],
    title="Virtual Nodes: Mengeliminasi Hash Ring Skew",
    subtitle="Memetakan satu node fisik ke ratusan titik token acak di cincin hash",
    content_items=[
        {"title": "MASALAH HASH SKEW", "items": [
            "Pada jumlah node sedikit (misal 3-4 node), titik acak tidak menyebar rata",
            "Satu node bisa menguasai 70% bentang cincin sementara yang lain hanya 10%",
            "Menimbulkan ketimpangan beban penyimpanan dan throughput yang tajam",
        ], "type": "bullet"},
        {"title": "MEKANISME VIRTUAL NODES", "items": [
            "Setiap node fisik dialokasikan R virtual token (misal 100-150 vnodes)",
            "Format token vkey: shardID + '#' + strconv.Itoa(i)",
            "Titik-titik virtual terdistribusi acak merata di sekeliling cincin",
        ], "type": "bullet"},
        {"title": "KODE INISIALISASI VNODE", "items": [
            "for i := 0; i < ch.vnodeCount; i++ {",
            "    vkey := shardID + \"#\" + strconv.Itoa(i)",
            "    h := hashKey(vkey)",
            "    ch.ring = append(ch.ring, vnode{hash: h, shardID: shardID})",
            "}",
            "sort.Slice(ch.ring, func(i, j int) bool { ... })",
        ], "type": "code"},
    ],
    mental_model="Virtual nodes interleave physical servers across the ring, statistically guaranteeing balanced distribution.",
    source_path="labs/32-database-sharding-and-partitioning/internal/sharding/sharding.go:115",
    slide_num="13 / 24",
    narrative="Jika setiap server hanya memiliki satu titik token di cincin hash, posisi acak sering kali memicu ketimpangan distribusi (hash skew). Menggunakan virtual nodes (vnodes) mengatasi hal ini dengan memberikan 100 hingga 150 token virtual bagi setiap server fisik. Titik-titik ini terselip secara berselang-seling di sepanjang cincin, memastikan probabilitas penempatan kunci terbagi seimbang dan seragam ke seluruh server fisik.",
    title_en="Virtual Nodes: Eradicating Hash Ring Skew",
    narrative_en="Assigning a single ring token per physical host frequently introduces severe hash skew, where one server disproportionately covers large arc segments. Virtual nodes (vnodes) solve this by assigning 100 to 150 distinct token coordinates per physical node. Interleaving virtual positions across the ring guarantees statistically balanced key distribution and uniform workload allocation."
)

# Slide 14: Perbandingan Empiris Resharding
add_slide(
    categories=["EVIDENSI TEST", "ANALISIS BENCHMARK"],
    title="Pembuktian Resharding: Modulo vs Consistent Hash",
    subtitle="Verifikasi matematis pergeseran kunci saat kluster scale-out dari 4 ke 5 node",
    content_items=[
        {"title": "TABEL PERBANDINGAN LAB", "items": [
            "Metrik | Hash Modulo | Consistent Hashing (100 vnodes)",
            "Dataset Diuji | 10.000 Kunci Akun | 10.000 Kunci Akun",
            "Kunci Berpindah | 7.984 Kunci | 1.200 Kunci",
            "Rasio Relokasi | 79.84% Data Pindah | 12.00% Data Pindah",
            "Target Minimal | ~80.00% | ~20.00% (1/N teoritis)",
        ], "type": "table"},
        {"title": "TEST RUNNER ASSERTIONS", "items": [
            "// tests/sharding_test.go:111",
            "if modMoveRatio < 0.65 { t.Errorf(...) }",
            "// Modulo terbukti memindahkan 75.12% - 79.84% kunci",
            "if chMoveRatio <= 0.05 || chMoveRatio > 0.40 { t.Errorf(...) }",
            "// Consistent Hash terbukti memindahkan hanya 12.00% - 16.00% kunci",
        ], "type": "code"},
        {"title": "EFISIENSI MIGRASI DATA", "items": [
            "Consistent Hashing menghemat hingga 85% bandwidth transfer resharding",
            "Mereduksi durasi migrasi kluster dari jam menjadi hitungan menit",
            "Menjaga ketersediaan aplikasi selama proses rebalancing berlangsung",
        ], "type": "bullet"},
    ],
    mental_model="Consistent hashing reduces data migration during scale-out from 80% down to 12%, saving 85% network I/O.",
    source_path="labs/32-database-sharding-and-partitioning/tests/sharding_test.go:63",
    slide_num="14 / 24",
    narrative="Uji coba otomatis di tests/sharding_test.go menguji langsung perbedaan kedua pendekatan saat menambah kapasitas dari 4 menjadi 5 node. Modulo memindahkan 7.984 dari 10.000 kunci (79,84%), sedangkan Consistent Hashing hanya memindahkan 1.200 kunci (12,00%). Hasil ini membuktikan Consistent Hashing mengeliminasi hingga 85 persen kebutuhan migrasi data, menjadikan proses penambahan node di sistem produksi jauh lebih efisien.",
    title_en="Empirical Resharding Proof: Modulo vs Consistent Hash",
    narrative_en="Automated tests in tests/sharding_test.go benchmark both strategies during cluster expansion from 4 to 5 shards. Modulo relocates 7,984 out of 10,000 keys (79.84%), whereas consistent hashing relocates only 1,200 keys (12.00%). This confirms an 85 percent reduction in data movement, enabling live, low-latency scale-out operations in production environments."
)

# Slide 15: Dilema Non-Sharded Query
add_slide(
    categories=["PROBLEM", "QUERY ENGINE"],
    title="Dilema Query Tanpa Menyertakan Shard Key",
    subtitle="Ketika pencarian dilakukan berdasarkan kolom sekunder seperti email atau status",
    content_items=[
        {"title": "SKENARIO KASUS NYATA", "items": [
            "Tabel data dishard menggunakan tenant_id sebagai shard key utama",
            "Alur Login User: Pengguna login hanya memasukkan email dan password",
            "Aplikasi tidak mengetahui tenant_id pengguna sebelum membaca datanya",
        ], "type": "bullet"},
        {"title": "DUA JALUR PENYELESAIAN", "items": [
            "Scatter-Gather: Broadcast query paralel ke seluruh shard di kluster",
            "Global Secondary Index (GSI): Menjaga indeks global pemetaan email -> tenant_id",
            "Masing-masing memiliki trade-off performa, biaya, dan konsistensi",
        ], "type": "bullet"},
        {"title": "RISIKO OPERASIONAL", "items": [
            "Broadcast ke puluhan shard menghabiskan connection pool kluster",
            "Tail latency query ditentukan oleh shard terlambat (straggler problem)",
            "GSI memerlukan penulisan ganda (double-write) dan mekanisme sinkronisasi",
        ], "type": "bullet"},
    ],
    mental_model="Queries lacking the sharding key face a hard trade-off: scatter-gather broadcast storm vs GSI lookup overhead.",
    source_path="labs/32-database-sharding-and-partitioning/content/02-master-draft.md",
    slide_num="15 / 24",
    narrative="Dalam database ter-shard, query yang menyertakan shard key dapat langsung diarahkan ke satu node target. Masalah muncul ketika aplikasi harus mencari data berdasarkan kolom sekunder, seperti mencari pengguna berdasarkan alamat email pada tabel yang di-shard berdasarkan tenant_id. Tanpa mengetahui shard key, query engine dihadapkan pada dua opsi: mem-broadcast query ke setiap shard (scatter-gather) atau memanfaatkan Global Secondary Index.",
    title_en="The Non-Sharded Query Dilemma",
    narrative_en="In a sharded architecture, queries including the shard key route directly to the designated host. A major architectural challenge emerges when searching against secondary attributes like email in a cluster sharded by tenant_id. Without knowing the shard key upfront, the application must choose between broadcasting across every shard (scatter-gather) or maintaining a Global Secondary Index."
)

# Slide 16: Scatter-Gather (Broadcast Query)
add_slide(
    categories=["QUERY EXECUTION", "SCATTER-GATHER"],
    title="Scatter-Gather: Eksekusi Broadcast Paralel",
    subtitle="Membagi query ke semua shard secara paralel dan menggabungkan hasilnya",
    content_items=[
        {"title": "ALUR KERJA SCATTER-GATHER", "items": [
            "1. Router menduplikasi query dan mengirimkannya ke N shard serentak",
            "2. Shard memproses query lokal dan mengembalikan hasilnya via channel",
            "3. Router menggabungkan (gather), mengurutkan, dan membatasi (LIMIT) hasil",
        ], "type": "bullet"},
        {"title": "IMPLEMENTASI CONTEXT DI GO", "items": [
            "for _, s := range shards {",
            "    go func(shard *Shard) {",
            "        select {",
            "        case <-ctx.Done(): return",
            "        default: ch <- shardResult{records: searchLocal(shard)}",
            "        }",
            "    }(s)",
            "}",
        ], "type": "code"},
        {"title": "THE STRAGGLER PROBLEM", "items": [
            "Respon total dibatasi oleh node yang paling lambat merespon (P99 tail)",
            "Jika 1 dari 50 shard mengalami disk stall, seluruh query ikut tertahan",
            "Membengkakkan CPU dan memory overhead pada router agregator",
        ], "type": "bullet"},
    ],
    mental_model="Scatter-gather latency is bound by the slowest responding shard in the broadcast fan-out.",
    source_path="labs/32-database-sharding-and-partitioning/internal/sharding/sharding.go:343",
    slide_num="16 / 24",
    narrative="Scatter-gather mengatasi query tanpa shard key dengan menembakkan query paralel ke seluruh node shard menggunakan goroutine, lalu mengumpulkan hasilnya di koordinator. Meskipun fleksibel tanpa membutuhkan tabel indeks tambahan, pendekatan ini sangat rentan terhadap straggler problem. Latensi total kueri ditentukan oleh node paling lambat, dan overhead broadcast meningkat drastis seiring bertambahnya jumlah shard.",
    title_en="Scatter-Gather: Parallel Broadcast Execution",
    narrative_en="Scatter-gather addresses non-shard-key lookups by dispatching parallel queries to every shard node via concurrent workers and aggregating responses. While requiring no auxiliary indexing, it suffers severely from the straggler problem. Total query latency is dictated by the slowest responding node, and broadcast overhead scales linearly with cluster size."
)

# Slide 17: Global Secondary Index (Lookup Vindex)
add_slide(
    categories=["INDEXING", "GSI"],
    title="Global Secondary Index: Point Lookup Terarah",
    subtitle="Menghindari broadcast badai dengan pemetaan kunci sekunder ke shard key",
    content_items=[
        {"title": "STRUKTUR PEMETAAN GSI", "items": [
            "Tabel Lookup Global (Lookup Vindex):",
            "  email: 'user_342@company.com' -> shardKey: 'tenant-42'",
            "Point lookup tahap 1: Query GSI untuk mendapatkan shardKey",
            "Point lookup tahap 2: Query langsung ke node target (1 node saja)",
        ], "type": "bullet"},
        {"title": "KODE IMPLEMENTASI GSI", "items": [
            "func (c *Cluster) GetByEmailUsingGSI(email, id string) (Record, error) {",
            "    shardKey, ok := c.gsi.Lookup(email)",
            "    if !ok { return Record{}, errors.New(\"not found\") }",
            "    return c.GetByShardKey(shardKey, id)",
            "}",
        ], "type": "code"},
        {"title": "TRADE-OFF DOUBLE WRITE", "items": [
            "Write Amplification: Setiap insert/update email wajib memperbarui tabel GSI",
            "Konsistensi Data: Membutuhkan 2PC atau eventual consistency asinkron",
            "Kapasitas GSI: Tabel GSI itu sendiri pada akhirnya harus di-shard",
        ], "type": "bullet"},
    ],
    mental_model="GSI trades write amplification for point-lookup speed, replacing cluster broadcasts with directed two-step reads.",
    source_path="labs/32-database-sharding-and-partitioning/internal/sharding/sharding.go:228",
    slide_num="17 / 24",
    narrative="Global Secondary Index (GSI), yang dalam teknologi seperti Vitess dikenal sebagai Lookup Vindex, memetakan atribut sekunder langsung ke shard key utama. Query login cukup membaca indeks GSI untuk mengetahui nomor shard pemilik data, lalu melakukan point-lookup ke node bersangkutan. Keuntungan latensi baca yang drastis ini diimbangi biaya penulisan ganda (double-write) dan kebutuhan koordinasi konsistensi saat data dimutasi.",
    title_en="Global Secondary Index: Directed Point Lookups",
    narrative_en="A Global Secondary Index (GSI), styled as a Lookup Vindex in Vitess, maps secondary attributes directly back to primary shard keys. Authentication workflows query the GSI to retrieve the hosting shard identifier, then execute a single targeted point-lookup. This latency reduction trades against write amplification and eventual consistency coordination during insert operations."
)

# Slide 18: Pembuktian Benchmark: Broadcast vs GSI
add_slide(
    categories=["BENCHMARK", "BUKTI KUANTITATIF"],
    title="Hasil Uji Coba: Broadcast vs GSI",
    subtitle="Pengukuran aktual eksekusi query non-shard key pada kluster 4 node",
    content_items=[
        {"title": "LOG HASIL EKSEKUSI DEMO", "items": [
            "1. Scatter-Gather Broadcast (by Email):",
            "   Nodes Broadcasted : 4 / 4 Shards",
            "   Records Matched   : 1",
            "   Execution Time    : 72.834 µs",
            "2. Global Secondary Index Lookup:",
            "   Nodes Broadcasted : 1 (Point Lookup Terarah)",
            "   Execution Time    : 1.208 µs (60x Lebih Cepat!)",
        ], "type": "code"},
        {"title": "ANALISIS LATENSI PRODUKSI", "items": [
            "Pada in-memory test: Selisih latensi mencapai 60 kali lipat",
            "Pada database fisik dengan jaringan: Perbedaan menjadi 1 RTT vs N RTT",
            "Scatter-gather menghabiskan thread pool dan memory buffer di seluruh node",
        ], "type": "bullet"},
        {"title": "KESIMPULAN DESAIN", "items": [
            "Gunakan GSI untuk akses data frekuensi tinggi (misal: otentikasi login)",
            "Gunakan Scatter-Gather hanya untuk laporan analitik batch frekuensi rendah",
            "GSI mutlak diperlukan untuk mencegah connection starvation di skala besar",
        ], "type": "bullet"},
    ],
    mental_model="GSI point-lookups bypass cluster fan-out entirely, cutting network round trips and query latency by orders of magnitude.",
    source_path="labs/32-database-sharding-and-partitioning/cmd/demo/main.go:141",
    slide_num="18 / 24",
    narrative="Eksekusi demo lab memperlihatkan perbedaan performa nyata antara kedua strategi. Scatter-gather mem-broadcast query ke seluruh 4 shard dan selesai dalam waktu 72,83 µs, sementara GSI menyelesaikan point-lookup hanya dalam 1,20 µs—60 kali lebih cepat. Di dunia nyata dengan latensi jaringan dan koneksi I/O fisik, eliminasi broadcast ini krusial untuk mencegah kehabisan connection pool pada jam sibuk.",
    title_en="Benchmark Verification: Broadcast vs GSI",
    narrative_en="Lab execution logs demonstrate the stark performance gap between both strategies. Scatter-gather broadcasts across all 4 shards, completing in 72.83 µs, while GSI lookup finishes in 1.20 µs—a 60x speed advantage. Across physical networks where latency and TCP overhead dominate, eliminating broadcast storms prevents connection pool starvation under heavy concurrency."
)

# Slide 19: Kegagalan Auto-Increment pada Multi-Shard
add_slide(
    categories=["DISTRIBUTED ID", "PROBLEM IDENTITAS"],
    title="Masalah Kunci Primer pada Multi-Shard",
    subtitle="Auto-increment lokal menghasilkan ID duplikat di berbagai node independen",
    content_items=[
        {"title": "ANATOMI TABRAKAN ID", "items": [
            "Shard A menghasilkan ID baru: 101, 102, 103",
            "Shard B menghasilkan ID baru: 101, 102, 103 (Duplikasi!)",
            "Primary key global tidak lagi unik di level kluster",
        ], "type": "diagram"},
        {"title": "KELEMAHAN UUIDv4 ACAK", "items": [
            "UUIDv4 (128-bit acak murni) menyelesaikan masalah collision unik",
            "Namun keacakan bit merusak index B-tree database fisik",
            "Memicu random page split, cache thrashing, dan fragmentasi disk parah",
        ], "type": "bullet"},
        {"title": "PERSYARATAN ID TERDISTRIBUSI", "items": [
            "Unik global tanpa perlu koordinasi jaringan berulang",
            "Monotonik terurut waktu (time-ordered) untuk menjaga efisiensi B-tree",
            "Mendukung throughput tinggi hingga jutaan ID per detik",
        ], "type": "bullet"},
    ],
    mental_model="Database auto-increment breaks across independent shards; random UUIDs degrade B-tree write performance through cache thrashing.",
    source_path="labs/32-database-sharding-and-partitioning/content/02-master-draft.md",
    slide_num="19 / 24",
    narrative="Mekanisme auto-increment bawaan database hanya menjamin keunikan di dalam satu mesin. Ketika dataset dipecah ke beberapa shard independen, masing-masing node akan menghasilkan ID lokal yang bertabrakan. Menggantinya dengan UUIDv4 memang menjamin keunikan global, tetapi sifat acaknya menghancurkan struktur indeks B-tree database fisik karena memicu page split terus-menerus dan menghancurkan cache locality.",
    title_en="Primary Key Collisions Across Independent Shards",
    narrative_en="Native database auto-increment guarantees uniqueness only within a single engine boundary. When sharding across multi-node topologies, separate shards generate colliding sequence integers. Replacing sequences with random UUIDv4 guarantees global uniqueness, but its non-sequential entropy induces severe B-tree page splits and cache fragmentation in disk-backed storage engines."
)

# Slide 20: RFC 9562 UUIDv7 & Sequence Block Allocator
add_slide(
    categories=["SOLUSI ID", "STANDAR INDUSTRI"],
    title="Strategi Distributed ID: UUIDv7 vs Block Allocator",
    subtitle="Dua pola standar industri untuk primary key terdistribusi berkinerja tinggi",
    content_items=[
        {"title": "RFC 9562 UUIDv7 (TIME-ORDERED)", "items": [
            "48-bit Unix timestamp millisecond + version 7 + 74-bit random entropy",
            "Terurut waktu secara leksikografis (u1 < u2) — ramah append pada B-tree",
            "Dihasilkan sepenuhnya lokal di memori aplikasi tanpa koordinasi jaringan",
        ], "type": "bullet"},
        {"title": "SEQUENCE BLOCK ALLOCATOR (VITESS)", "items": [
            "Layanan pusat mengalokasikan blok ID numerik (misal: 1.000 ID sekaligus)",
            "Node aplikasi mengonsumsi ID secara lokal hingga blok habis sebelum alokasi baru",
            "Menghasilkan integer 64-bit padat dan hemat ukuran memori",
        ], "type": "bullet"},
        {"title": "HASIL VERIFIKASI DEMO", "items": [
            "UUIDv7 Tergenerate : 01a13905-fefc-7db1-adc7-ed5774231cd0",
            "Sequence Block IDs : 1 2 3 4 5 6 7 8 (Block Size=5)",
            "Keduanya terbukti bebas tabrakan pada pengujian konkuren.",
        ], "type": "code"},
    ],
    mental_model="UUIDv7 provides zero-coordination time-ordered 128-bit keys; Sequence block allocation provides compact 64-bit integers with batched coordination.",
    source_path="labs/32-database-sharding-and-partitioning/internal/idgen/idgen.go:18",
    slide_num="20 / 24",
    narrative="Untuk mengatasi collision ID tanpa merusak B-tree, dua solusi diterapkan. Standar RFC 9562 UUIDv7 menggabungkan timestamp 48-bit dengan entropi acak, menghasilkan ID 128-bit yang terurut leksikografis tanpa memerlukan koordinasi jaringan. Pendekatan kedua, Sequence Block Allocation (gaya Vitess), mengambil blok nomor integer berurutan dari server pusat dalam batch besar, menghasilkan ID integer 64-bit yang padat dan efisien.",
    title_en="Distributed ID Strategies: UUIDv7 vs Block Allocators",
    narrative_en="Two industry patterns solve distributed identity without corrupting B-tree locality. The RFC 9562 UUIDv7 standard embeds a 48-bit millisecond timestamp ahead of random entropy, producing lexicographically sortable 128-bit identifiers without network round-trips. Alternatively, Vitess-style Sequence Block Allocation leases integer blocks in batches, providing compact 64-bit monotonic sequences with minimal centralized lock contention."
)

# Slide 21: Kompleksitas Transaksi Cross-Shard & 2PC
add_slide(
    categories=["KONSISTENSI", "TRANSAKSI TERDISTRIBUSI"],
    title="Kompleksitas Transaksi Lintas Shard & TwoPC",
    subtitle="Ketika atomisitas menuntut koordinasi antar server: Two-Phase Commit dan Fractured Reads",
    content_items=[
        {"title": "TWO-PHASE COMMIT (2PC)", "items": [
            "Phase 1 (Prepare): Koordinator meminta vote commit ke semua shard terkait",
            "Phase 2 (Commit): Jika semua setuju, koordinator menerbitkan perintah commit",
            "Penalti: Menambah latensi round-trip dan menahan lock baris lebih lama",
        ], "type": "bullet"},
        {"title": "BAHAYA FRACTURED READS", "items": [
            "2PC hanya menjamin atomisitas eksekusi (All or Nothing)",
            "2PC TIDAK menyediakan isolasi transaksi lintas shard secara otomatis",
            "Query baca konkuren dapat melihat shard A sudah commit tetapi shard B belum",
        ], "type": "bullet"},
        {"title": "ATURAN ARSITEKTUR UTAMA", "items": [
            "Desain skema agar 95%+ transaksi terjadi di dalam 1 shard yang sama",
            "Gunakan Co-location (Entity Grouping) berdasarkan Shard Key yang identik",
            "Gunakan Saga Pattern atau Eventual Consistency untuk transaksi lintas shard",
        ], "type": "bullet"},
    ],
    mental_model="Two-Phase Commit guarantees cross-shard atomicity, not global isolation; co-locate related records to avoid distributed transactions entirely.",
    source_path="labs/32-database-sharding-and-partitioning/content/05-key-takeaways.md",
    slide_num="21 / 24",
    narrative="Ketika satu transaksi harus memperbarui data di dua shard berbeda, arsitektur terdistribusi membutuhkan protokol Two-Phase Commit (2PC). Meskipun 2PC menjamin atomisitas penulisan, protokol ini melipatgandakan latensi dan tidak memberikan isolasi baca lintas shard (rawan fractured reads). Strategi terbaik adalah mengelompokkan data terkait dalam satu shard key yang sama (co-location) sehingga transaksi terdistribusi dapat dihindari sepenuhnya.",
    title_en="Cross-Shard Transaction Complexity & Two-Phase Commit",
    narrative_en="When transactions span across distinct shard nodes, systems rely on Two-Phase Commit (2PC) coordination. While 2PC enforces write atomicity, it compounds network round-trips and fails to provide global snapshot isolation, exposing concurrent readers to fractured reads. The gold standard in distributed schema design is co-locating relational entities under the same shard key, eliminating distributed transactions entirely."
)

# Slide 22: Batasan Implementasi Lab 32
add_slide(
    categories=["BATAS IMPLEMENTASI", "DISCLAIMER"],
    title="Batas Model dan Ruang Lingkup Lab 32",
    subtitle="Simulasi in-memory terfokus pada kebenaran algoritma dan invariant matematis",
    content_items=[
        {"title": "TERSEDIA DI LAB INI", "items": [
            "Implementasi Range Partitioning dengan verifikasi pruning dan fast drop",
            "Modulo Router vs Consistent Hash Router dengan 100-150 virtual nodes",
            "Scatter-Gather Broadcast paralel dengan konteks pembatalan di Go",
            "Global Secondary Index (GSI) dan dua generator ID terdistribusi",
        ], "type": "bullet"},
        {"title": "BELUM DICAKUP / BATASAN", "items": [
            "Penyimpanan berbasis disk nyata (menggunakan Go map in-memory terlindungi mutex)",
            "Mesin Two-Phase Commit (2PC) transaksi lintas shard fisik",
            "Mekanisme data replication (Raft/Paxos) antar shard replicas",
            "Proses online data streaming resharding tanpa downtime produksi",
        ], "type": "bullet"},
        {"title": "TUJUAN EDUKASI", "items": [
            "Membuktikan secara deterministik angka pergeseran kunci dan invariant routing",
            "Menghubungkan teori buku sistem terdistribusi dengan baris kode Go yang dapat di-debug",
        ], "type": "bullet"},
    ],
    mental_model="This lab models core distributed algorithms in-memory to prove routing invariants without disk storage noise.",
    source_path="labs/32-database-sharding-and-partitioning/README.md",
    slide_num="22 / 24",
    narrative="Lab ini dirancang sebagai model simulasi in-memory yang ketat untuk menguji invariant algoritma sharding dan partitioning di level kode Go. Implementasi ini tidak menyertakan penyimpanan file fisik, replikasi konsensus Raft, atau engine transaksi 2PC terdistribusi. Fokus utama lab adalah membuktikan secara kuantitatif efisiensi pruning, stabilitas relokasi consistent hashing, perbandingan kecepatan GSI, serta karakteristik ID terdistribusi.",
    title_en="Implementation Boundaries of Lab 32",
    narrative_en="This laboratory models sharding and partitioning mechanics in-memory to rigorously verify algorithmic invariants within pure Go code. The codebase deliberately omits physical disk engines, Raft consensus replication, and distributed 2PC transaction managers. Its mission is providing concrete, verifiable metrics on partition pruning, consistent hashing relocation, GSI performance, and distributed identifier behavior."
)

# Slide 23: Checklist Keputusan Arsitektur Data
add_slide(
    categories=["PANDUAN OPERASIONAL", "CHECKLIST"],
    title="Checklist Keputusan Sharding vs Partitioning",
    subtitle="Kerangka kerja pengambilan keputusan sebelum memecah arsitektur database Anda",
    content_items=[
        {"title": "CHECKLIST EVALUASI SISTEM", "items": [
            "[ ] Apakah volume data melampaui kapasitas 1 server? Jika belum, gunakan Partitioning.",
            "[ ] Apakah bottleneck pada lifecycle data historis? Gunakan Range Partitioning & Drop.",
            "[ ] Shard Key memiliki kardinalitas tinggi dan terdistribusi acak (hindari tanggal murni).",
            "[ ] Gunakan Consistent Hashing dengan 100+ vnodes untuk meminimalkan data shuffle.",
            "[ ] Pasang GSI untuk pencarian atribut sekunder frekuensi tinggi.",
            "[ ] Gunakan RFC 9562 UUIDv7 untuk primary key terdistribusi ramah B-tree.",
            "[ ] Rancang skema co-location agar 95%+ query tidak membutuhkan cross-shard join.",
        ], "type": "bullet"},
        {"title": "KAPAN PARTITIONING CUKUP", "items": [
            "Dataset besar tetapi I/O server masih mencukupi",
            "Query membutuhkan partition pruning berbasis tanggal/kategori",
            "Retensi data membutuhkan penghapusan tabel lama secara rutin",
        ], "type": "bullet"},
        {"title": "KAPAN HARUS SHARDING", "items": [
            "Kapasitas write I/O server tunggal telah mencapai 100% saturasi",
            "Ukuran dataset melebihi kapasitas storage fisik server terbesar",
            "Kebutuhan ketersediaan tinggi di berbagai zona ketersediaan geografis",
        ], "type": "bullet"},
    ],
    mental_model="Exhaust single-node table partitioning before adopting physical sharding; sharding is a last-resort scalability lever.",
    source_path="labs/32-database-sharding-and-partitioning/content/02-master-draft.md",
    slide_num="23 / 24",
    narrative="Sebelum melangkah ke kompleksitas sharding fisik, tim rekayasa wajib mengevaluasi apakah logical table partitioning pada server tunggal sudah cukup untuk mengatasi masalah performa dan retensi data. Sharding fisik adalah senjata pamungkas ketika batasan I/O dan memori mesin tunggal benar-benar telah terlampaui. Ikuti checklist ini untuk memastikan pemilihan shard key, algoritma routing, dan strategi indeks sekunder Anda dirancang secara tepat sejak hari pertama.",
    title_en="Architectural Decision Checklist: Sharding vs Partitioning",
    narrative_en="Engineering teams must evaluate whether single-node logical partitioning suffices before introducing physical sharding complexity. Sharding represents a high-overhead architectural shift justified only when single-node write throughput and storage limits are entirely exhausted. Follow this checklist to ensure proper shard key selection, routing algorithms, and distributed secondary indexing from day one."
)

# Slide 24: Closing (Dark Navy #0b1b2d)
add_slide(
    categories=["RINGKASAN"],
    title="Eksperimen kecil, interpretasi yang teliti.",
    subtitle=f"Commit: {COMMIT_SHA[:8]} | Test: 5/5 PASS | Repository: lukman-ss/software-engineering-lab",
    content_items=[
        {"title": "MODUL LAB DIUJI", "items": [
            "internal/partitioning • internal/sharding",
            "internal/idgen • tests/sharding_test.go",
            "cmd/demo/main.go • README.md",
        ], "type": "bullet"},
        {"title": "VERIFIKASI & METRIK", "items": [
            f"Sumber commit: {COMMIT_SHA[:8]}",
            "Test passing 5/5 dengan race detector lolos tanpa data race",
            "Consistent hash: 12% relokasi vs Modulo: 79.84% relokasi",
            "GSI lookup: 1.20 µs vs Scatter-gather: 72.83 µs (60x lebih cepat)",
        ], "type": "bullet"},
        {"title": "REPOSITORI", "items": [
            f"GitHub: {REPO_URL.replace('https://', '')}",
            "Jalankan demo : go run ./cmd/demo",
            "Jalankan test : go test -v -race ./...",
        ], "type": "bullet"},
    ],
    mental_model="Logical partitioning prunes scans locally; consistent sharding scales writes globally without resharding shock.",
    source_path="labs/32-database-sharding-and-partitioning/README.md",
    slide_num="24 / 24",
    narrative="",
    title_en="Small experiments, rigorous interpretations.",
    narrative_en=""
)

if __name__ == "__main__":
    print(f"Total slides defined in content_def.py: {len(SLIDES)}")
