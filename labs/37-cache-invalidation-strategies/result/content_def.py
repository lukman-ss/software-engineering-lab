#!/usr/bin/env python3
"""
Publication Builder for Software Engineering Lab 37: Cache Invalidation Strategies.
Generates: PPTX, PDF, PNG covers, Dev.to MD, Medium HTML, Metadata, Caption, DELIVERABLES.md
Adheres strictly to the Lab 13-15 editorial 2-column layout and in-depth technical narratives.
"""
import os
import sys

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/37-cache-invalidation-strategies/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/37-cache-invalidation-strategies"
COMMIT_SHA = "eac3d80ae5344055dc773e5e798bf50f3675a2cc"
LAB_NO = "37"
SLUG = "cache-invalidation-strategies"
TITLE_ID = "Strategi Invalidation Cache"
SUBTITLE_ID = "Pola Write, Mitigasi Stampede, dan Jitter TTL dalam Go"
TITLE_EN = "Cache Invalidation Strategies"
SUBTITLE_EN = "Write Patterns, Stampede Mitigations, and TTL Jitter in Go"
AUTHOR = "Lukman / lukman-ss"
TOPIC_TAGS = "Cache Invalidation • Stampede • Singleflight • XFetch • SWR • Go"

# Slide data structure
SLIDES = []

def add_slide(categories, title, subtitle, content_items, mental_model, source_path="", slide_num="", narrative="", narrative_en=""):
    """Add a slide to the presentation with support for in-depth editorial narratives."""
    SLIDES.append({
        "categories": categories,
        "title": title,
        "subtitle": subtitle,
        "content_items": content_items,
        "mental_model": mental_model,
        "source": source_path,
        "slide_num": slide_num,
        "narrative": narrative,
        "narrative_en": narrative_en
    })

# =====================================================================
# SLIDE CONTENT DEFINITIONS (24 SLIDES - EDITORIAL LAB 13-15 STANDARD)
# =====================================================================

# Slide 1: Cover (Dark Navy retro terminal / art)
add_slide(
    categories=["COVER"],
    title=TITLE_ID,
    subtitle=SUBTITLE_ID,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "Retro CRT monitor beige dengan scanlines",
            "ASCII art terminal yang menampilkan cache miss/hit ratio",
            f"Ditulis oleh: {AUTHOR}",
            "Go 1.22 • golang.org/x/sync/singleflight",
        ]}
    ],
    mental_model="",
    source_path="",
    slide_num="01 / 24"
)

# Slide 2: Dilema Terbesar Sistem Terdistribusi
add_slide(
    categories=["MOTIVASI", "BUKTI LAB"],
    title="Dilema Terbesar Sistem Terdistribusi",
    subtitle="Cache invalidation adalah salah satu dari tiga masalah sulit dalam ilmu komputer.",
    content_items=[
        {"title": "KUTIPAN KLASIK", "items": [
            "Phil Karlton: \"There are only two hard things in Computer Science:",
            "cache invalidation and naming things.\" (dan off-by-one errors).",
            "Problem: cache menyimpan data lama; DB terupdate; reader melihat stale data.",
        ], "type": "bullet"},
        {"title": "MOTIVASI", "items": [
            "Cache mempercepat read dengan menyimpan hasil query DB di memori.",
            "Tanpa invalidasi tepat: inkonsistensi, cache stampede, dan outage.",
            "Trade-off: freshness vs throughput pada write policies.",
        ], "type": "bullet"},
        {"title": "BATASAN LAB", "items": [
            "Lab 37 fokus pada 3 write policies + 4 mitigasi stampede (SF, XFetch, SWR, Jitter).",
            "Implementasi in-memory Go, tanpa Redis/PostgreSQL eksternal.",
            "Verifikasi via unit test & demo CLI (race-clean).",
        ], "type": "bullet"}
    ],
    narrative="Cache invalidation adalah salah satu dilema terbesar dalam arsitektur sistem terdistribusi skala modern. Cache dirancang sebagai lapisan memori kerja yang sangat cepat di depan database, bertujuan agar ribuan query per detik dapat dilayani tanpa membebani storage persisten. Namun, ketika record di database mengalami perubahan, cache harus segera diperbarui atau dihapus. Tanpa mekanisme invalidasi yang disiplin, cache akan terus menyimpan data basi, memicu inkonsistensi data di level aplikasi, serta dapat merusak integritas transaksi bisnis secara meluas.",
    narrative_en="Cache invalidation remains one of the hardest challenges in distributed systems. While caching shields persistent databases from excessive traffic, maintaining freshness during underlying database updates is critical. Without disciplined invalidation policies, stale reads and data corruption will degrade application reliability and consistency.",
    mental_model="Cache invalidation is the hardest problem because consistency is eventual, never instant.",
    source_path="labs/37-cache-invalidation-strategies/content/02-master-draft.md",
    slide_num="02 / 24"
)

# Slide 3: Taksonomi Lab 37
add_slide(
    categories=["ARSITEKTUR", "TAXONOMY"],
    title="Taksonomi Lab 37",
    subtitle="Tiga pola write + empat mitigasi stampede dalam satu arsitektur Go.",
    content_items=[
        {"title": "WRITE PATTERNS", "items": [
            "Pola | Write Flow | Read Flow | Karakteristik",
            "Cache-Aside | DB write → cache Delete | Cache miss → DB query → Set | Lazy, default pragmatic",
            "Write-Through | DB write → cache Set (sync) | Cache hit | Strong read-after-write",
            "Write-Behind | Cache Set → async queue flush | Cache hit | High write throughput, durability risk",
        ], "type": "table"},
        {"title": "STAMPEDE MITIGATIONS", "items": [
            "Mitigasi | Teknik | Scope | Formula / Mekanisme",
            "Singleflight | Request coalescing | Intra-process Go | singleflight.Group.Do(key, fn)",
            "XFetch | Probabilistic early refresh | Per-request | -Δ·β·ln(U) > TTL_remaining",
            "SWR | Stale serve + async revalidate | Per-key | ExpiresAt + staleDelta window",
            "TTL Jitter | Random TTL offset | Cross-key | TTLWithJitter(base, maxJitter)",
        ], "type": "table"}
    ],
    narrative="Lab 37 mengklasifikasikan strategi pengelolaan cache ke dalam dua dimensi arsitektural utama: write patterns dan stampede mitigations. Tiga pola penulisan—Cache-Aside, Write-Through, dan Write-Behind—menentukan keseimbangan antara latensi baca-tulis serta jaminan konsistensi data. Di sisi lain, empat mitigasi stampede—Singleflight, XFetch, Stale-While-Revalidate, dan TTL Jitter—bekerja sama untuk mencegah kehancuran database saat kunci panas kedaluwarsa secara mendadak di bawah lonjakan konkurensi. Seluruh komponen dimodelkan secara murni dalam Go standar tanpa dependensi luar.",
    narrative_en="Lab 37 categorizes cache engineering into write policies and stampede mitigations. Three write strategies dictate latency and consistency trade-offs, while four stampede mitigations protect databases from thundering herds during hot key expiration. All patterns are implemented natively in Go without external infrastructure.",
    mental_model="Map cache strategies into write policies vs stampede mitigations matrix.",
    source_path="labs/37-cache-invalidation-strategies/content/02-master-draft.md",
    slide_num="03 / 24"
)

# Slide 4: Cache-Aside (Lazy Loading) - Mekanisme & Flow
add_slide(
    categories=["WRITE POLICY", "CACHE-ASIDE"],
    title="Cache-Aside (Lazy Loading) - Mekanisme & Flow",
    subtitle="Mekanisme Get: cache miss → DB query → cache Set. Mekanisme Put: DB write → Invalidate cache.",
    content_items=[
        {"title": "GET FLOW (MISS)", "items": [
            "1. cache.Get(key) → ErrCacheMiss (cache kosong)",
            "2. db.Query(ctx, key) → fetch value dengan latency Δ = time.Since(start)",
            "3. cache.Set(key, val, ttl, delta) → populasi cache dengan TTL",
            "4. Return value ke reader",
        ], "type": "bullet"},
        {"title": "UPDATE FLOW (INVALIDATE)", "items": [
            "1. db.Write(ctx, key, newVal) → commit ke store persisten",
            "2. cache.Delete(key) → invalidate secara eksplisit",
            "3. Reader berikutnya: miss → fetch DB → cache.Set",
        ], "type": "bullet"},
        {"title": "ALUR KODE", "items": [
            "func (s *CacheAsideService) Update(ctx, key, val) error {",
            "    s.db.Write(ctx, key, val) // 1. DB dulu",
            "    s.cache.Delete(key)      // 2. baru invalidate",
            "    return nil",
            "}",
        ], "type": "code"}
    ],
    narrative="Pola Cache-Aside atau lazy loading merupakan pendekatan paling populer pada sistem dengan beban baca tinggi. Pada alur pembacaan, aplikasi memeriksa cache terlebih dahulu; jika terjadi cache miss, aplikasi mengambil data langsung dari database lalu menyimpannya ke cache dengan TTL tertentu. Saat data diperbarui, aplikasi wajib menulis perubahan ke database terlebih dahulu sebelum menghapus kunci dari cache. Pola ini memastikan hanya data yang benar-benar dibaca yang menempati ruang memori cache.",
    narrative_en="Cache-Aside operates lazily: reads populate the cache upon a miss, and updates commit to the database before invalidating the cached key. This lazy population saves memory by caching only requested records, while making write-then-invalidate ordering essential for correctness.",
    mental_model="Cache-aside: lazy populate on miss, write-through-then-invalidate on update.",
    source_path="labs/37-cache-invalidation-strategies/internal/cache/patterns.go:22-47",
    slide_num="04 / 24"
)

# Slide 5: Cache-Aside - Anomali & Race Condition
add_slide(
    categories=["ANTI-PATTERN", "RACE CONDITION"],
    title="Cache-Aside - Anomali & Race Condition",
    subtitle="Stale read window dan concurrent write-read sequence.",
    content_items=[
        {"title": "STALE READ WINDOW", "items": [
            "Setelah cache miss → DB query dimulai (bisa memakan waktu Δ).",
            "Sepanjang proses tersebut, cache masih kosong.",
            "Reader lain yang tiba di window ini juga melakukan miss → multiple DB queries.",
        ], "type": "bullet"},
        {"title": "CONCURRENT WRITE-READ RACE", "items": [
            "T1: reader mulai DB query untuk key K",
            "T2: writer melakukan db.Write(K, newVal) + cache.Delete(K)",
            "T1: menyelesaikan query → cache.Set(K, OLD_VALUE) → menimpa nilai yang benar!",
            "Reader lanjutan: cache hit pada OLD_VALUE yang sudah lama seharusnya invalid.",
        ], "type": "bullet"},
        {"title": "CONTOH SEKWENSIAL", "items": [
            "Goroutine A: cache.Get(K) miss → mulai db.Query()",
            "Goroutine B: Update(K, newVal) → db.Write() + cache.Delete()",
            "Goroutine A: selesai db.Query() → cache.Set(K, OLD_VALUE)",
            "Goroutine C: cache.Get(K) → hit OLD_VALUE (stale)",
        ], "type": "bullet"}
    ],
    narrative="Meskipun fleksibel, pola Cache-Aside memiliki celah race condition antara thread pembaca dan thread penulis yang berjalan bersamaan. Jika reader mengalami cache miss lalu memulai query database yang lambat, sementara writer secara bersamaan menulis data baru ke database dan menghapus cache, reader yang terlambat dapat menuliskan kembali data usang ke dalam cache. Akibatnya, nilai usang tersebut akan bertahan di cache hingga masa TTL berakhir, menghasilkan anomali data inkonsisten yang berbahaya.",
    narrative_en="Cache-Aside exhibits a subtle race window under concurrent access. If a reader suffers a cache miss while a writer commits an update, a slow-completing reader can overwrite fresh invalidation with stale data, keeping invalid records in cache until TTL expiration.",
    mental_model="Naive cache-aside lets late miss-writes clobber validated cache entries.",
    source_path="labs/37-cache-invalidation-strategies/content/02-master-draft.md",
    slide_num="05 / 24"
)

# Slide 6: Write-Through - Konsistensi Sinkron
add_slide(
    categories=["WRITE POLICY", "WRITE-THROUGH"],
    title="Write-Through - Konsistensi Sinkron",
    subtitle="Setiap Update: DB write → sync cache Set. Read berikutnya: cache hit, nilai terbaru.",
    content_items=[
        {"title": "WRITE FLOW", "items": [
            "1. db.Write(ctx, key, val) → commit ke store persisten (blocking)",
            "2. delta = time.Since(start)",
            "3. cache.Set(key, val, ttl, delta) → update cache secara sinkron",
            "4. Reader berikutnya: cache.Get(key) → hit nilai terbaru (tidak perlu DB query)",
        ], "type": "bullet"},
        {"title": "KODE IMPLEMENTASI", "items": [
            "func (s *WriteThroughService) Update(ctx, key, val) error {",
            "    start := time.Now()",
            "    if err := s.db.Write(ctx, key, val); err != nil {",
            "        return fmt.Errorf(\"db write failed: %w\", err)",
            "    }",
            "    delta := time.Since(start)",
            "    s.cache.Set(key, val, s.ttl, delta)",
            "    return nil",
            "}",
        ], "type": "code"},
        {"title": "KONTRAK KUAT", "items": [
            "Read-after-write consistency terjamin untuk proses berikutnya.",
            "Tidak ada stale read window pada update path.",
            "Cache selalu berisi nilai terbaru yang pernah ditulis.",
        ], "type": "bullet"}
    ],
    narrative="Pola Write-Through menawarkan jaminan read-after-write consistency yang sangat kuat dengan memperlakukan cache sebagai penyimpanan utama yang sinkron. Setiap operasi pembaruan menulis data ke database persisten dan segera memperbarui kunci cache terkait dalam satu alur eksekusi aplikasi secara berurutan. Hasilnya, setiap request pembacaan berikutnya dijamin langsung memperoleh data terbaru dari cache tanpa perlu query database tambahan. Pola ini mengeliminasi celah data basi yang sering terjadi pada Cache-Aside.",
    narrative_en="Write-Through enforces strong read-after-write consistency by synchronizing database commits directly with cache updates. Subsequent readers immediately hit the cache and receive the freshest data without querying the database, eliminating the stale window inherent in lazy caching.",
    mental_model="Write-through guarantees read-after-write via synchronous DB+cache update.",
    source_path="labs/37-cache-invalidation-strategies/internal/cache/patterns.go:78-89",
    slide_num="06 / 24"
)

# Slide 7: Write-Through - Biaya Latensi & Partial Failure
add_slide(
    categories=["FAILURE", "PERFORMA"],
    title="Write-Through - Biaya Latensi & Partial Failure",
    subtitle="Database write penalty dan risiko kegagalan sebagian (partial failure).",
    content_items=[
        {"title": "BIAYA LATENSI", "items": [
            "Write latency = DB_write_time + Cache_set_time (dua operasi blocking berurutan)",
            "Jika DB latency naik, write latency ikut naik secara linier.",
            "Tidak cocok untuk write-heavy workload dengan throughput tinggi.",
        ], "type": "bullet"},
        {"title": "PARTIAL FAILURE SCENARIO", "items": [
            "1. db.Write berhasil commit ke database",
            "2. cache.Set gagal (mis. memory pressure, node eviction)",
            "3. State: DB terupdate, cache kosong/stale",
            "4. Reader berikutnya: miss → fetch DB lama yang sudah benar (recover otomatis via TTL)",
        ], "type": "bullet"},
        {"title": "TABEL RISIKO", "items": [
            "Failure Point | Deteksi | Recovery | Impact",
            "DB write fail | Error return | Retry caller | Write lost",
            "Cache set fail | Silent (swallowed) | TTL safety net | Temporary miss",
            "Both succeed | - | - | 0 impact, consistent",
        ], "type": "table"}
    ],
    narrative="Kelemahan mendasar Write-Through terletak pada penalti latensi tulis dan risiko kegagalan parsial di lingkungan terdistribusi. Latensi setiap operasi update membengkak karena harus menunggu dua operasi I/O sinkron secara berturut-turut: penulisan database dan pembaruan cache. Selain itu, jika database commit berhasil namun penulisan cache gagal akibat kendala memori, status cache menjadi inkonsisten. Dalam skenario ini, TTL bertindak sebagai jaring pengaman agar data tidak terkunci dalam status usang selamanya.",
    narrative_en="Write-Through incurs a latency penalty by executing two consecutive synchronous writes for every update. Furthermore, partial failures—such as a committed database write followed by a failed cache set—demand TTL safety nets to prevent prolonged cache inconsistency across application instances.",
    mental_model="Write-through doubles sync write cost; cache-set failures degrade gracefully via TTL.",
    source_path="labs/37-cache-invalidation-strategies/content/02-master-draft.md",
    slide_num="07 / 24"
)

# Slide 8: Write-Behind (Write-Back) - Throughput Ekstrem
add_slide(
    categories=["WRITE POLICY", "WRITE-BEHIND"],
    title="Write-Behind (Write-Back) - Throughput Ekstrem",
    subtitle="Cache Set langsung, enqueue async queue, flushWorker background ke DB.",
    content_items=[
        {"title": "WRITE FLOW", "items": [
            "1. cache.Set(key, val, ttl, delta) → update cache secara langsung (synchronous)",
            "2. select { case writeQueue <- req: | default: } → enqueue ASYNC",
            "3. flushWorker goroutine: writeQueue <- range → db.Write() di background",
            "4. Reader berikutnya: cache.Get() → hit nilai terbaru (zero DB round-trip)",
        ], "type": "bullet"},
        {"title": "KODE IMPLEMENTASI", "items": [
            "func (s *WriteBehindService) Update(key, val string) {",
            "    s.cache.Set(key, val, s.ttl, 1*time.Millisecond)",
            "    select {",
            "    case s.writeQueue <- WriteRequest{Key: key, Value: val}:",
            "    default: // queue full — drop (demo limitation)",
            "    }",
            "}",
        ], "type": "code"},
        {"title": "THROUGHPUT GAIN", "items": [
            "Write latency ≈ cache.Set time (milidetik), bukan DB round-trip.",
            "flushWorker bisa memproses batch flush untuk efisiensi lebih lanjut.",
            "Reader selalu mendapatkan nilai terbaru dari cache, bukan dari DB.",
        ], "type": "bullet"}
    ],
    narrative="Write-Behind atau Write-Back memprioritaskan throughput penulisan ekstrem dengan memutus ketergantungan langsung terhadap latensi database persisten. Saat pembaruan terjadi, aplikasi langsung memperbarui cache di memori dan memasukkan tugas penulisan ke dalam antrean asinkron, lalu segera merespons pemanggil. Goroutine background secara berkala menguras antrean dan mengeksekusi penulisan ke database di latar belakang. Dengan mekanisme ini, latensi tulis ditekan hingga mendekati nol milidetik, sangat cocok untuk beban tulis masif.",
    narrative_en="Write-Behind maximizes write throughput by updating memory cache immediately and delegating database persistence to an asynchronous queue drained by background workers. Write operations return instantly with near-zero latency, making it ideal for high-velocity telemetry and write-heavy workloads.",
    mental_model="Write-behind decouples write latency by buffering to async queue; cache stays fresh.",
    source_path="labs/37-cache-invalidation-strategies/internal/cache/patterns.go:91-161",
    slide_num="08 / 24"
)

# Slide 9: Write-Behind - Risiko Durabilitas & Data Loss
add_slide(
    categories=["FAILURE", "DATA LOSS"],
    title="Write-Behind - Risiko Durabilitas & Data Loss",
    subtitle="Queue crash, queue overflow, dan data loss saat process termination.",
    content_items=[
        {"title": "CRASH SCENARIO", "items": [
            "Process mengalami crash (panic, OOM, kill) sebelum flushWorker mengeksekusi writeQueue.",
            "Item dalam queue yang belum di-flush = data loss.",
            "Cache (in-memory) juga hilang → key tidak dapat dipulihkan.",
        ], "type": "bullet"},
        {"title": "QUEUE OVERFLOW", "items": [
            "select-default: bila queue penuh, write request di-drop secara silent.",
            "Tidak ada sinyal kirim ke caller bahwa write gagal.",
            "Pada demo: 10 writes dengan buffer size 2 → hanya 3 yang diproses.",
            "Sisa 7 write hilang tanpa jejak.",
        ], "type": "bullet"},
        {"title": "MITIGASI PRODUKSI", "items": [
            "Persistent queue (Kafka, WAL, file log) sebelum flush ke DB.",
            "Acknowledge semantics: caller tunggu hingga flushWorker meng-ack.",
            "Overflow handling: retry, backpressure, atau dead-letter queue.",
        ], "type": "bullet"}
    ],
    narrative="Kecepatan luar biasa Write-Behind harus dibayar dengan kompromi serius pada aspek durabilitas data. Jika proses aplikasi mengalami crash, panic, atau kehabisan memori sebelum worker latar belakang sempat mengosongkan antrean, seluruh data pembaruan yang masih tertahan di antrean akan hilang permanen. Selain itu, jika antrean penuh di bawah beban tinggi, implementasi yang menggunakan select-default akan membuang request penulisan secara diam-diam. Sistem produksi wajib menggunakan antrean persisten seperti Kafka atau WAL.",
    narrative_en="Write-Behind trades durability for raw execution speed. Unflushed queue items are permanently lost if the process crashes unexpectedly. In high-concurrency bursts, overflowing queues can silently discard writes unless backed by persistent message brokers like Kafka or Write-Ahead Logs.",
    mental_model="Write-behind trades write durability for latency; use persistent queues in production.",
    source_path="labs/37-cache-invalidation-strategies/content/02-master-draft.md",
    slide_num="09 / 24"
)

# Slide 10: Thundering Herd / Cache Stampede
add_slide(
    categories=["STAMPEDE", "FAILURE MODE"],
    title="Thundering Herd / Cache Stampede",
    subtitle="N concurrent requests pada hot key expiry → N DB queries → DB overload.",
    content_items=[
        {"title": "MEKANISME STAMPEDE", "items": [
            "1. Hot key K kadaluarsa (TTL habis serentak)",
            "2. N concurrent goroutine melakukan cache.Get(K) → semua miss",
            "3. Setiap goroutine menjalankan db.Query(K) secara paralel",
            "4. DB menerima N query identik → CPU & connection pool terbebani",
            "5. Latency naik berhampa → hit rate jatuh, beban makin buruk (cascading failure)",
        ], "type": "bullet"},
        {"title": "UKURAN BEBAN", "items": [
            "Workload | Concurrency N | DB Query Count (Naive)",
            "Local dev | 3 | 3",
            "Staging | 50 | 50",
            "Production hot key | 10,000 | 10,000 (DB overload)",
        ], "type": "table"},
        {"title": "DEMO VERIFIKASI", "items": [
            "20 concurrent goroutines → cache miss semua",
            "NaiveStampedeService: 20 DB queries (λ·T amplification)",
            "SingleFlightService: 1 DB query (coalescing)",
        ], "type": "bullet"}
    ],
    narrative="Cache stampede atau thundering herd terjadi ketika kunci cache populer yang sering diakses kedaluwarsa secara mendadak. Puluhan hingga ribuan goroutine pembaca yang datang bersamaan akan mendapati cache miss, lalu serentak mengeksekusi query database yang identik. Fenomena ini melipatgandakan beban database sebesar perkalian laju kedatangan request dan durasi query. Akibatnya, connection pool terkuras habis, waktu tunggu melonjak tajam, dan sistem dapat mengalami kegagalan berantai hingga ketersediaan layanan runtuh sepenuhnya.",
    narrative_en="A cache stampede occurs when a high-traffic hot key expires suddenly. Hundreds of concurrent requests experience simultaneous cache misses and flood the database with identical expensive queries, exhausting connection pools and causing systemic congestion collapse across the infrastructure.",
    mental_model="Stampede amplifies concurrent misses into redundant DB query storms.",
    source_path="labs/37-cache-invalidation-strategies/content/02-master-draft.md",
    slide_num="10 / 24"
)

# Slide 11: Singleflight - Request Coalescing
add_slide(
    categories=["STAMPEDE", "SINGLEFLIGHT"],
    title="Singleflight - Request Coalescing",
    subtitle="golang.org/x/sync/singleflight, N requests ditumpuk menjadi 1 DB call.",
    content_items=[
        {"title": "PRINSIP", "items": [
            "singleflight.Group.Do(key, fn) menjamin fn dieksekusi tepat 1x per key.",
            "Goroutine pendatang yang request key yang sama akan menunggu hasil dari goroutine pertama.",
            "Semua waiter menerima return value yang sama (deduplication barrier).",
        ], "type": "bullet"},
        {"title": "IMPORT", "items": [
            "import \"golang.org/x/sync/singleflight\"",
            "",
            "type SingleFlightService struct {",
            "    cache  *MemoryCache",
            "    db     *MockDB",
            "    ttl    time.Duration",
            "    flight singleflight.Group",
            "}",
        ], "type": "code"},
        {"title": "DEMO VERIFIKASI", "items": [
            "20 concurrent goroutines → cache.Get miss semua",
            "flight.Do(key, fn) → fn dieksekusi 1x → db.Query 1x",
            "20 goroutines menerima hasil yang sama: db.QueryCount() == 1",
        ], "type": "bullet"}
    ],
    narrative="Mitigasi pertama terhadap cache stampede adalah request coalescing menggunakan paket resmi Go golang.org/x/sync/singleflight. Pola ini memastikan bahwa untuk kunci yang sama, fungsi pengambilan data ke database hanya dijalankan tepat satu kali oleh satu goroutine pemenang. Goroutine-goroutine lain yang datang bersamaan akan ditahan pada barrier internal dan menunggu selesainya eksekusi goroutine pertama. Setelah selesai, seluruh pemanggil menerima salinan hasil yang identik, memangkas puluhan query database menjadi satu panggilan tunggal.",
    narrative_en="Singleflight provides intra-process request coalescing via Go official sync library. When multiple goroutines request an expired hot key, Singleflight executes the underlying database fetch exactly once while suspending other callers until the shared result is broadcast to all waiters.",
    mental_model="Singleflight dedupes concurrent calls: N requests collapse to 1 DB query.",
    source_path="labs/37-cache-invalidation-strategies/internal/cache/stampede.go:45-84",
    slide_num="11 / 24"
)

# Slide 12: Singleflight - Kode & Mekanika Do()
add_slide(
    categories=["SINGLEFLIGHT", "KODE"],
    title="Singleflight - Kode & Mekanika Do()",
    subtitle="Double-check cache hit di dalam closure, deduplication barrier, shared result.",
    content_items=[
        {"title": "IMPLEMENTASI FULL", "items": [
            "func (s *SingleFlightService) Get(ctx, key) (string, error) {",
            "    item, err := s.cache.Get(key) // outer cache check",
            "    if err == nil { return item.Value, nil }",
            "    res, err, _ := s.flight.Do(key, func() (interface{}, error) {",
            "        // Double-check cache inside flight execution",
            "        item, err := s.cache.Get(key)",
            "        if err == nil { return item.Value, nil }",
            "        start := time.Now()",
            "        val, err := s.db.Query(ctx, key)",
            "        if err != nil { return \"\", err }",
            "        delta := time.Since(start)",
            "        s.cache.Set(key, val, s.ttl, delta)",
            "        return val, nil",
            "    })",
            "    if err != nil { return \"\", err }",
            "    return res.(string), nil",
            "}",
        ], "type": "code"},
        {"title": "MEKANIKA INTERNAL", "items": [
            "1. singleflight.Group memelihara map[key]callEntry",
            "2. Goroutine pertama: membuat entry, mengeksekusi fn, menyimpan hasil",
            "3. Goroutine follow: mendaftar sebagai waiter pada channel yang sama",
            "4. Setelah fn selesai: semua waiter dibangunkan, menerima hasil yang sama",
            "Result channel: shared = true",
        ], "type": "bullet"}
    ],
    narrative="Implementasi SingleFlightService menerapkan teknik double-check locking di dalam closure fungsi Do untuk efisiensi maksimal. Sebelum memulai query database yang mahal, goroutine memeriksa kembali isi cache untuk mengantisipasi kemungkinan kunci tersebut telah berhasil diisi oleh eksekusi paralel lain. Singleflight mengelola tabel pemetaan pemanggilan aktif secara internal dan membagikan return value melalui channel bersama. Namun, perlindungan ini terbatas pada level satu proses; sistem multi-pod tetap membutuhkan distributed lock seperti Redis SETNX.",
    narrative_en="Inside the Singleflight closure, a double-check pattern inspects the cache once more before triggering the database fetch, ensuring optimal efficiency. While Singleflight eliminates intra-process stampedes, multi-node deployments still require distributed coordination such as Redis distributed locks.",
    mental_model="In-flight closure double-checks cache to avoid redundant re-fetch on coalescing.",
    source_path="labs/37-cache-invalidation-strategies/internal/cache/stampede.go:56-84",
    slide_num="12 / 24"
)

# Slide 13: XFetch - Probabilistic Early Expiration
add_slide(
    categories=["X-FETCH", "PROBABILISTIC"],
    title="XFetch - Probabilistic Early Expiration",
    subtitle="Formula: -Δ · β · ln(U) > TTL_remaining, compute cost Δ terukur per request.",
    content_items=[
        {"title": "FORMULA XFETCH", "items": [
            "offset = -Δ · β · ln(U)",
            "Recompute sebelum TTL resmi ketika:",
            "  -Δ · β · ln(U) > TTL_remaining",
            "",
            "Parameters:",
            "- Δ = time.Since(start) pada DB query terakhir (compute cost)",
            "- β = tuning parameter (default = 1.0, >1 = lebih agresif, <1 = lebih konservatif)",
            "- U ~ Uniform(0,1), random draw per request",
        ], "type": "bullet"},
        {"title": "KODE ShouldRecompute", "items": [
            "func ShouldRecompute(delta time.Duration, beta float64,",
            "    ttlRemaining time.Duration, u float64) bool {",
            "    if u <= 0 || u >= 1 { return false }",
            "    deltaSec := delta.Seconds()",
            "    ttlRemainingSec := ttlRemaining.Seconds()",
            "    expiryCompute := -deltaSec * beta * math.Log(u)",
            "    return expiryCompute > ttlRemainingSec",
            "}",
        ], "type": "code"},
        {"title": "DISTRIBUSI", "items": [
            "Traffic tinggi (Δ besar) → ekspresi -ln(U) cenderung trigger refresh lebih awal",
            "Traffic jarang → sering melewati kondisi, gunakan cache sampai TTL resmi",
            "Menggeser puncak rebuild dari saturasi satu titik ke distribusi exponential",
        ], "type": "bullet"}
    ],
    narrative="XFetch adalah algoritma mitigasi stampede probabilistik berbasis riset akademis Vattani et al. Setiap kali request membaca cache, algoritma menghitung ambang batas kedaluwarsa dini berdasarkan biaya komputasi query database terakhir, parameter tuning beta, dan bilangan acak seragam. Jika nilai komputasi probabilistik melebihi sisa TTL, request tersebut akan berinisiatif merefresh cache secara proaktif sebelum waktu kedaluwarsa resmi tiba. Pola ini menyebarkan beban pembaruan cache ke sepanjang kurva eksponensial.",
    narrative_en="XFetch employs probabilistic early expiration derived from academic research by Vattani et al. By evaluating background fetch duration and uniform random draws against remaining TTL, high-traffic keys trigger proactive refreshes before official expiration, smoothing out concurrent rebuild spikes.",
    mental_model="XFetch probabilistically shifts cache rebuild from peak expiry to exponential spread.",
    source_path="labs/37-cache-invalidation-strategies/internal/cache/stampede.go:125-136",
    slide_num="13 / 24"
)

# Slide 14: XFetch - Tanda Negatif Matematika
add_slide(
    categories=["MATH", "BUG"],
    title="XFetch - Tanda Negatif Matematika",
    subtitle="Koreksi -ln(U) > 0 vs ln(U) < 0 bug, beta parameter menggerakkan distribusi.",
    content_items=[
        {"title": "FORMULA BENAR", "items": [
            "expiryCompute = -deltaSec × beta × Log(u)",
            "karena U ∈ (0,1) → ln(U) < 0 → -ln(U) > 0",
            "Sehingga LHS selalu POSITIF, bisa lebih besar dari TTL_remaining.",
        ], "type": "bullet"},
        {"title": "BUG FORMULA PEDAGOGIS", "items": [
            " expiryCompute = deltaSec × beta × Log(u)  (tanpa minus!)",
            " ln(U) ∈ (-∞, 0) untuk U ∈ (0,1), jadi LHS selalu NEGATIF",
            " RHS = ttlRemainingSec > 0 (selalu positif)",
            " Hasil: NEGATIF > POSITIF = selalu FALSE",
            " Efek: early refresh TIDAK PERNAH trigger (stampede tetap terjadi).",
        ], "type": "bullet"},
        {"title": "VERIFIKASI TEST", "items": [
            "Test 1: u=0.5 → -ln(0.5)=0.6931, Δ×0.6931=69.31ms < 200ms → False",
            "Test 2: u=0.3 → -ln(0.3)=1.2039, Δ×1.2039=120.39ms > 50ms → True",
            "Test 3: tanpa minus sign → hasil selalu False (bug teguh)",
        ], "type": "bullet"}
    ],
    narrative="Penerapan matematis formula XFetch menuntut pemahaman mendalam terhadap tanda minus logaritma natural. Karena variabel U berada pada rentang nol hingga satu, nilai logaritma natural ln(U) selalu menghasilkan bilangan negatif. Penambahan tanda minus di depan rumus membuat hasil kalkulasi bernilai positif, sehingga memungkinkan perbandingan yang valid terhadap sisa TTL. Kesalahan umum dalam mengabaikan tanda minus akan menghasilkan nilai negatif permanen, menyebabkan kondisi penyegaran awal tidak pernah terpenuhi dan stampede tetap terjadi.",
    narrative_en="Mathematical precision in XFetch requires retaining the negative sign in front of the natural logarithm. Because ln(U) is strictly negative for random values between 0 and 1, omitting the minus sign forces the expression negative, silently disabling proactive early recomputation.",
    mental_model="-ln(U) maps uniform draws into positive offsets; missing minus sign silently disables early refresh.",
    source_path="labs/37-cache-invalidation-strategies/tests/cache_test.go:161-184",
    slide_num="14 / 24"
)

# Slide 15: Stale-While-Revalidate (SWR) - Toleransi Data Basi
add_slide(
    categories=["SWR", "STALE-WHILE-REVALIDATE"],
    title="Stale-While-Revalidate (SWR) - Toleransi Data Basi",
    subtitle="Fresh vs Stale window vs Hard miss, RFC 5861 Cache-Control extension.",
    content_items=[
        {"title": "WINDOW STATE MACHINE", "items": [
            "State | Kondisi | Action | Return",
            "Fresh | now < ExpiresAt | Return cached | Cached value",
            "Stale-1 | ExpiresAt ≤ now < staleUntil | Trigger async revalidate | Stale value",
            "Stale-2 | now ≥ staleUntil (past grace) | Synchronous fetch | New DB value (blocking)",
            "Miss | Key not in cache | Synchronous fetch | New DB value",
        ], "type": "table"},
        {"title": "ALUR KODE GET", "items": [
            "item, ok := s.cache.GetRaw(key)",
            "now := time.Now()",
            "if ok {",
            "    if now.Before(item.ExpiresAt) { return item.Value, nil }       // Fresh",
            "    staleUntil := item.ExpiresAt.Add(s.staleDelta)",
            "    if now.Before(staleUntil) {",
            "        s.triggerRevalidate(key)                                   // Stale-1",
            "        return item.Value, nil",
            "    }",
            "}",
            "// Hard miss / Stale-2 → synchronous fetch",
            "val, err := s.db.Query(ctx, key)",
            "s.cache.Set(key, val, ttl, delta)",
            "return val, nil",
        ], "type": "code"}
    ],
    narrative="Stale-While-Revalidate yang diadaptasi dari RFC 5861 menawarkan solusi performa tinggi dengan memperluas toleransi terhadap data kedaluwarsa. Sistem membagi siklus hidup cache menjadi tiga jendela: fresh di mana data disajikan langsung, stale window di mana data usang segera disajikan ke klien sementara revalidasi dipicu di latar belakang, serta hard miss di mana query database harus dieksekusi secara sinkron. Pendekatan ini menjaga latensi pembacaan tetap berada di level sub-milidetik.",
    narrative_en="Stale-While-Revalidate, inspired by RFC 5861, balances latency and freshness across three lifecycle phases. Within the stale grace window, clients receive immediate cached responses while background revalidation fetches updated records asynchronously, guaranteeing low end-user latency.",
    mental_model="SWR extends cache lifetime into a grace window while async-revalidating in background.",
    source_path="labs/37-cache-invalidation-strategies/internal/cache/stampede.go:197-224",
    slide_num="15 / 24"
)

# Slide 16: SWR - Async Revalidation & Deduplikasi
add_slide(
    categories=["SWR", "CONCURRENCY"],
    title="SWR - Async Revalidation & Deduplikasi",
    subtitle="Serving stale immediately, background goroutine update, revalidating map guard.",
    content_items=[
        {"title": "REVALIDASI ASYNC", "items": [
            "triggerRevalidate(key): spawn single goroutine per key",
            "go func() { db.Query → cache.Set }",
            "Reader tidak blocking; langsung menerima stale value.",
            "Revalidasi berlangsung di background, tidak menunggu selesai.",
        ], "type": "bullet"},
        {"title": "DEDUPLIKASI (revalidating map)", "items": [
            "s.mu.Lock(); if s.revalidating[key] { s.mu.Unlock(); return }",
            "s.revalidating[key] = true; s.mu.Unlock()",
            "// goroutine: defer delete(s.revalidating, key)",
            "Hanya 1 goroutine yang melakukan revalidasi per key.",
        ], "type": "bullet"},
        {"title": "KODE TRIGGER", "items": [
            "func (s *SWRService) triggerRevalidate(key string) {",
            "    s.mu.Lock()",
            "    if s.revalidating[key] { s.mu.Unlock(); return }",
            "    s.revalidating[key] = true",
            "    s.mu.Unlock()",
            "    s.revalCount.Add(1)",
            "    go func() { // async fetch & set",
            "        defer func() {",
            "            s.mu.Lock()",
            "            delete(s.revalidating, key)",
            "            s.mu.Unlock()",
            "        }()",
            "        ...db.Query → s.cache.Set...",
            "    }()",
            "}",
        ], "type": "code"}
    ],
    narrative="Keberhasilan SWR sangat bergantung pada mekanisme deduplikasi revalidasi latar belakang agar tidak memicu badai permintaan baru ke database. Tanpa pembatas, ratusan request pada data basi akan masing-masing meluncurkan goroutine revalidasi baru. SWRService mengatasi ancaman ini dengan memelihara tabel penanda revalidasi berbasis mutex; hanya satu goroutine yang diizinkan memperbarui database untuk kunci yang sama pada satu waktu. Jika revalidasi gagal, klien tetap aman menerima nilai usang sebagai fallback.",
    narrative_en="SWR requires robust deduplication to prevent background revalidation storms. Using a mutex-protected tracking map, SWRService guarantees that only a single background goroutine refreshes the database per key at any given time, gracefully falling back to stale data upon upstream errors.",
    mental_model="Serve stale immediately, dedupe background revalidation to one goroutine per key.",
    source_path="labs/37-cache-invalidation-strategies/internal/cache/stampede.go:226-254",
    slide_num="16 / 24"
)

# Slide 17: TTL Jitter - Anti-Sinkronisasi Expiry
add_slide(
    categories=["JITTER", "DESINKRONISASI"],
    title="TTL Jitter - Anti-Sinkronisasi Expiry",
    subtitle="Base TTL + random offset [0, maxJitter), desynchronization.",
    content_items=[
        {"title": "FORMULA JITTER", "items": [
            "TTL_effective = base + random(0, maxJitter)",
            "Implementasi: time.Duration(rand.Int63n(int64(maxJitter)))",
            "Hasil berada di range [base, base + maxJitter)",
            "Tiap key kadaluarsa pada waktu yang berbeda.",
        ], "type": "bullet"},
        {"title": "KODE IMPLEMENTASI", "items": [
            "func TTLWithJitter(base, maxJitter time.Duration) time.Duration {",
            "    if maxJitter <= 0 { return base }",
            "    jitter := time.Duration(rand.Int63n(int64(maxJitter)))",
            "    return base + jitter",
            "}",
        ], "type": "code"},
        {"title": "MANFAAT DESINKRONISASI", "items": [
            "Cross-key: 1000 key dengan base=10s, jitter=[0,5s) → expire selama 10-15s",
            "Tidak ada lonjakan simultan expiry pada t = 10s semua sekaligus",
            "Mengurangi probabilitas banyak key miss secara bersamaan",
            "Bukan solusi untuk single hot key stampede (lihat slide 10)",
        ], "type": "bullet"}
    ],
    narrative="TTL Jitter adalah teknik pertahanan sederhana namun ampuh untuk mencegah sinkronisasi waktu kedaluwarsa lintas kunci cache. Ketika ribuan kunci diisi secara bersamaan saat startup sistem atau deploy aplikasi, penggunaan TTL statis akan menyebabkan seluruh kunci tersebut kedaluwarsa pada detik yang sama persis. Dengan menambahkan offset acak pada rentang nol hingga batas maksimum jitter, masa kedaluwarsa tersebar merata, mengeliminasi lonjakan serentak pada database di masa mendatang.",
    narrative_en="TTL Jitter adds a random time offset to base TTL values, preventing cross-key expiration synchronization. By desynchronizing bulk-primed keys across a bounded time window, Jitter eliminates simultaneous mass expirations that often overwhelm backend storage following application deployments.",
    mental_model="Jitter desynchronizes cross-key expiries, spreading cache misses over time.",
    source_path="labs/37-cache-invalidation-strategies/internal/cache/store.go:78-86",
    slide_num="17 / 24"
)

# Slide 18: Pengujian: Cache-Aside vs Write-Through vs Write-Behind
add_slide(
    categories=["TESTING", "VERIFIKASI"],
    title="Pengujian: Cache-Aside vs Write-Through vs Write-Behind",
    subtitle="TestCachePatterns, test output & query counters.",
    content_items=[
        {"title": "TEST CASE OVERVIEW", "items": [
            "Test | Pattern | Assertion | Expected DB Queries",
            "Cache-Aside Read | Cache-Aside | Miss → 1 query | 1",
            "Cache-Aside Hit | Cache-Aside | Hit → 0 query | 1 (stays)",
            "Cache-Aside Update | Cache-Aside | Write+Delete → next miss | 2",
            "Write-Through Write | Write-Through | Update sync set cache | 1 (no extra query)",
            "Write-Behind Flush | Write-Behind | Async flush after sleep | 1 write async",
        ], "type": "table"},
        {"title": "DEMO OUTPUT", "items": [
            "[Cache-Aside]",
            "  1st GET (cache miss)  -> DB queries: 1, val: Laptop",
            "  2nd GET (cache hit)   -> DB queries: 1, val: Laptop",
            "  After update & GET    -> DB queries: 2, val: Laptop Pro",
            "",
            "[Write-Through]",
            "  1st GET (cache miss)  -> DB queries: 1",
            "  After update & GET    -> DB queries: 1, val: Phone Pro",
            "",
            "[Write-Behind]",
            "  After Update, GET from cache -> DB writes: 0, val: Tablet",
            "  After 50ms flush delay       -> DB writes: 1",
        ], "type": "code"}
    ],
    narrative="Pengujian TestCachePatterns memvalidasi perilaku ketiga write pattern menggunakan MockDB dengan counter query atomik yang presisi. Pada Cache-Aside, query database hanya bertambah pada cache miss pertama dan tetap stabil pada hit kedua, lalu bertambah lagi setelah operasi update menghapus cache. Pada Write-Through, pembacaan pasca-update langsung mendapati cache hit tanpa query database. Pada Write-Behind, update cache terjadi seketika, dan penulisan database terbukti dieksekusi secara asinkron beberapa milidetik kemudian.",
    narrative_en="Unit testing via TestCachePatterns verifies all three write policies deterministically using MockDB query counters. Invalidation cycles, synchronous write-through hits, and asynchronous write-behind flushes are validated against exact query count assertions, confirming predictable IO behavior.",
    mental_model="MockDB query counters make write-policy behavior deterministic and testable.",
    source_path="labs/37-cache-invalidation-strategies/tests/cache_test.go:13-104",
    slide_num="18 / 24"
)

# Slide 19: Pengujian: Naive Stampede vs Singleflight
add_slide(
    categories=["TESTING", "STAMPEDE MITIGATION"],
    title="Pengujian: Naive Stampede vs Singleflight",
    subtitle="20 DB queries vs 1 DB query under 20 concurrent goroutines.",
    content_items=[
        {"title": "NAIVE STAMPEDE (TEST)", "items": [
            "concurrency := 20",
            "svc := NewNaiveStampedeService(mem, db, 1*time.Minute)",
            "20 goroutines → svc.Get(ctx, 'hot-key') secara paralel",
            "wg.Wait()",
            "ASSERT: db.QueryCount() > 1  (stampede terjadi!)",
            "RESULT: 20 concurrent → 20 DB queries",
        ], "type": "bullet"},
        {"title": "SINGLEFLIGHT (TEST)", "items": [
            "concurrency := 20",
            "svc := NewSingleFlightService(mem, db, 1*time.Minute)",
            "20 goroutines → svc.Get(ctx, 'hot-key') secara paralel",
            "wg.Wait()",
            "ASSERT: db.QueryCount() == 1  (coalescing)",
            "RESULT: 20 concurrent → 1 DB query, 20 correct results",
        ], "type": "bullet"},
        {"title": "DEMO OUTPUT", "items": [
            "[Naive] 20 concurrent goroutines -> 20 DB queries (stampede!)",
            "[SF]    20 concurrent goroutines -> 1 DB queries (singleflight), 20 correct results",
        ], "type": "code"}
    ],
    narrative="Pengujian mitigasi stampede membuktikan efektivitas Singleflight secara empiris di bawah beban konkurensi tinggi. Ketika dua puluh goroutine diluncurkan secara simultan untuk mengakses kunci panas yang belum ada di cache, implementasi Naive menghasilkan dua puluh query database paralel—memperlihatkan stampede nyata. Sebaliknya, SingleFlightService berhasil menggabungkan seluruh dua puluh panggilan tersebut menjadi tepat satu query database tunggal, dan seluruh goroutine menerima nilai yang benar tanpa terjadi data race.",
    narrative_en="Concurrency tests empirically contrast Naive Stampede against Singleflight. Under twenty concurrent goroutines requesting a cold hot-key, the naive service triggers twenty redundant database queries, whereas Singleflight successfully collapses all twenty requests into exactly one database invocation.",
    mental_model="Singleflight collapses concurrent misses: 20 requests → 1 DB query.",
    source_path="labs/37-cache-invalidation-strategies/tests/cache_test.go:106-159",
    slide_num="19 / 24"
)

# Slide 20: Pengujian: XFetch Deterministic Simulation
add_slide(
    categories=["TESTING", "X-FETCH"],
    title="Pengujian: XFetch Deterministic Simulation",
    subtitle="Low rand draw triggers refresh, high rand draw serves cache.",
    content_items=[
        {"title": "SETUP TEST", "items": [
            "delta = 100ms, beta = 1.0, ttl = 1s",
            "Initial: cache.Get miss → db.Query → cache.Set, db.QueryCount = 1",
            "Two deterministic randFunc injected via SetRandFunc()",
        ], "type": "bullet"},
        {"title": "HIGH U (0.5) → NO EARLY REFRESH", "items": [
            "remaining = 200ms, u = 0.5",
            "-ln(0.5) = 0.6931",
            "expiryCompute = 100ms × 0.6931 = 69.31ms < 200ms",
            "ShouldRecompute = FALSE → serve from cache",
            "RESULT: db.QueryCount stays 1 (cache hit)",
        ], "type": "bullet"},
        {"title": "LOW U (0.3) → TRIGGER EARLY REFRESH", "items": [
            "remaining = 50ms, u = 0.3",
            "-ln(0.3) = 1.2039",
            "expiryCompute = 100ms × 1.2039 = 120.39ms > 50ms",
            "ShouldRecompute = TRUE → recompute from DB",
            "RESULT: db.QueryCount grows to 2 (proactive refresh)",
        ], "type": "bullet"},
        {"title": "DEMO OUTPUT", "items": [
            "Initial GET       -> DB queries: 1, val: initial-val",
            "Low rand draw     -> DB queries: 2, val: refreshed-val (early proactive refresh)",
            "High rand draw    -> DB queries: 1, val: fresh-data (served from cache, no early refresh)",
        ], "type": "code"}
    ],
    narrative="Pengujian logika XFetch dirancang secara deterministik melalui injeksi fungsi pengacak randFunc tiruan untuk menguji kedua percabangan algoritma. Ketika generator acak menghasilkan nilai tinggi mendekati satu, formula matematis memutuskan tidak perlu recompute, sehingga data disajikan langsung dari cache. Sebaliknya, saat diinjeksi nilai acak yang sangat kecil, rumus menghasilkan nilai komputasi yang melampaui sisa TTL, memicu eksekusi query pembaruan dini ke database secara terbukti dan terukur.",
    narrative_en="Deterministic XFetch testing leverages injected random generators to validate decision paths reliably. High uniform draws preserve cache hits, while simulated low draws reliably satisfy the probabilistic formula, triggering early proactive recomputations before hard TTL expiration.",
    mental_model="Deterministic randFunc injection tests both XFetch refresh and cache-hit branches.",
    source_path="labs/37-cache-invalidation-strategies/tests/cache_test.go:161-305",
    slide_num="20 / 24"
)

# Slide 21: Pengujian: SWR & Race Detection
add_slide(
    categories=["TESTING", "SWR"],
    title="Pengujian: SWR & Race Detection",
    subtitle="Immediate stale serve, async reload verified, go test -race clean.",
    content_items=[
        {"title": "TEST SCENARIO", "items": [
            "TTL=20ms, staleDelta=300ms",
            "1. Initial Get → cache.Set, return 'initial-data'",
            "2. time.Sleep(30ms) → TTL expired, masuk stale window",
            "3. db.SetData('key', 'updated-data')",
            "4. Concurrent Get (10 goroutines) → semua dapat 'initial-data' (stale)",
            "5. triggerRevalidate dispatch 1 goroutine (deduplicated!)",
            "6. time.Sleep(50ms) → async revalidasi selesai",
            "7. Next Get → return 'updated-data' (fresh)",
        ], "type": "bullet"},
        {"title": "RACE DETECTION", "items": [
            "go test -race ./... → PASS (1.553s, 0 races detected)",
            "revalidateCount atomic: hanya 1 revalidasi untuk 10 concurrent stale requests",
            "revalidating map + mutex: benar-benar mencegah storm",
        ], "type": "bullet"},
        {"title": "DEMO OUTPUT", "items": [
            "Initial GET       -> val: v1",
            "GET (stale)       -> val: v1  (stale served immediately)",
            "GET (fresh)       -> val: v2  (async revalidation complete)",
        ], "type": "code"}
    ],
    narrative="Pengujian Stale-While-Revalidate memverifikasi bahwa data usang langsung dikembalikan kepada pemanggil tanpa jeda latensi begitu TTL kedaluwarsa namun masih dalam rentang toleransi stale. Sepuluh goroutine konkuren yang membaca kunci basi secara simultan terbukti hanya memicu tepat satu proses revalidasi latar belakang berkat mekanisme deduplikasi. Pengujian race detector Go dengan flag -race lulus sempurna tanpa kebocoran data, membuktikan keandalan thread-safety implementasi SWR di lingkungan produksi.",
    narrative_en="SWR test suites confirm immediate delivery of stale data alongside single-goroutine background revalidation. Under ten concurrent readers, deduplication prevents revalidation storms, while Go race detector passes cleanly, validating strict thread-safety across concurrent cache accesses.",
    mental_model="SWR serves stale immediately while async revalidation dedupes concurrent callers.",
    source_path="labs/37-cache-invalidation-strategies/tests/cache_test.go:186-221",
    slide_num="21 / 24"
)

# Slide 22: Matriks Komparasi 7 Strategi
add_slide(
    categories=["KOMPARASI", "MATRIKS"],
    title="Matriks Komparasi 7 Strategi",
    subtitle="Tabel lengkap: Pola, Latensi Tulis, Latensi Baca, Konsistensi, Risiko.",
    content_items=[
        {"title": "MATRIKS KOMPARASI LENGKAP", "items": [
            "Strategi | Latensi Tulis | Latensi Baca | Konsistensi | Risiko Utama",
            "Cache-Aside | Sync DB Write | Cache hit (0 RTT) | Eventual read-after-write | Stampede pada miss massal",
            "Write-Through | DB+Cache (2 sync ops) | Cache hit (0 RTT) | Strong read-after-write | Latency penulisan tinggi",
            "Write-Behind | Cache Set + async queue | Cache hit (0 RTT) | Eventual consistency | Data loss pada crash/overflow",
            "Singleflight | Sama | Cache hit setelah 1st miss | Eventual | Tidak proteksi cross-process",
            "XFetch | Sama | Cache hit/proactive refresh | Eventual | Sign error bug (formulasi salah)",
            "SWR | Sama | Stale serve (0 blocking) | Eventual | Staleness window, revalidation storm",
            "TTL Jitter | Sama | Cache hit (tersebar waktu) | Eventual | Tidak proteksi single hot key",
        ], "type": "table"},
        {"title": "SPESIFIKASI", "items": [
            "Kolom tambahan | Cache-Aside | Write-Through | Write-Behind",
            "Memory overhead | Lazy, hanya data yang dibaca | Semua key yang ditulis | Semua key yang ditulis",
            "Write throughput | Medium | Low | High (async)",
            "Read-after-write guarantee | No | Yes | No (eventual)",
            "Durability | DB committed | DB committed | Risk of loss before flush",
        ], "type": "table"}
    ],
    narrative="Matriks komparasi tujuh strategi menyajikan panduan evaluasi menyeluruh berdasarkan profil latensi tulis, latensi baca, tingkat konsistensi, dan risiko arsitektural utama. Tidak ada satu strategi tunggal yang sempurna untuk semua skenario; setiap pilihan merupakan kompromi teknik yang sadar. Cache-Aside unggul pada efisiensi memori, Write-Through memberikan konsistensi tertinggi dengan mengorbankan kecepatan tulis, sementara Write-Behind menawarkan performa penulisan maksimum namun menuntut mitigasi khusus terhadap risiko hilangnya data saat kegagalan sistem.",
    narrative_en="The comprehensive comparison matrix evaluates all seven strategies across write latency, read latency, consistency guarantees, and operational failure risks. Engineering teams must weigh strong read-after-write consistency against write latency and durability risks to identify optimal caching architectures.",
    mental_model="No single strategy dominates; trade write latency, read latency, consistency, risk.",
    source_path="labs/37-cache-invalidation-strategies/content/02-master-draft.md",
    slide_num="22 / 24"
)

# Slide 23: Panduan Adopsi & Rekomendasi Lapangan
add_slide(
    categories=["ADOPSI", "REKOMENDASI"],
    title="Panduan Adopsi & Rekomendasi Lapangan",
    subtitle="Decision tree, kapan pakai Cache-Aside, Write-Behind, SF, SWR.",
    content_items=[
        {"title": "DECISION TREE", "items": [
            "Read-heavy + cold writes?",
            "  → Cache-Aside (default) + TTL Jitter",
            "",
            "Hot key + stampede observed?",
            "  → Singleflight (single proc) OR Redis lock (multi-proc)",
            "  → XFetch (probabilistic proactive refresh)",
            "",
            "Tolerir staleness untuk latency ultra rendah?",
            "  → SWR (serve stale, async revalidate)",
            "",
            "Write-heavy + throughput kritis?",
            "  → Write-Behind + persistent queue (Kafka/WAL)",
            "  → Hindari jika durability kritis (financial)",
            "",
            "Read-after-write consistency diperlukan?",
            "  → Write-Through (accept write latency)",
        ], "type": "bullet"},
        {"title": "REKOMENDASI LAPANGAN", "items": [
            "Cache-Aside + TTL Jitter adalah default pragmatis untuk read-heavy API.",
            "Singleflight wajib dipasang pada hot key Go service (intra-proc).",
            "XFetch beta=1.0+ untuk sistem dengan compute cost variabel (DB query lambat).",
            "SWR ideal untuk API yang tahan terhadap data 30-60s lama.",
            "Write-Behind hanya untuk non-critical writes (telemetry, analytics).",
            "Gunakan Write-Through untuk financial/audit-sensitive data.",
        ], "type": "bullet"},
        {"title": "PRODUCTION UPGRADE PATH", "items": [
            "In-memory lab → Redis/Memcached distributed cache",
            "MockDB queryCount → Prometheus counter/histogram",
            "singleflight intra-proc → Redis SET NX PX distributed lock",
            "select-default queue → Kafka persistent queue + ack",
        ], "type": "bullet"}
    ],
    narrative="Panduan adopsi menetapkan pohon keputusan praktis untuk mengimplementasikan strategi cache di lingkungan produksi. Untuk sebagian besar API read-heavy, Cache-Aside yang dipadukan dengan TTL Jitter adalah pilihan default terbaik. Ketika beban lalu lintas tinggi memicu stampede pada kunci panas, pasang Singleflight untuk layanan monolitik atau kunci terdistribusi untuk arsitektur multi-pod. Terapkan SWR jika sistem menoleransi data usang demi latensi sub-milidetik, dan gunakan Write-Behind hanya jika kehilangan data non-kritis dapat diterima.",
    narrative_en="Adoption guidelines provide a decision tree mapping workload characteristics to caching strategies. Cache-Aside with Jitter serves as the pragmatic baseline, Singleflight guards hot keys, SWR satisfies ultra-low latency requirements, and Write-Behind suits non-critical high-throughput ingestion pipelines.",
    mental_model="Choose invalidation strategy by matching workload profile to tolerance for staleness and durability.",
    source_path="labs/37-cache-invalidation-strategies/content/02-master-draft.md",
    slide_num="23 / 24"
)

# Slide 24: Closing Slide (Dark Navy `#0b1b2d`)
add_slide(
    categories=["CLOSING"],
    title="Cache invalidation: di mana performance bertemu kebenaran data.",
    subtitle="",
    content_items=[
        {"title": "MODUL LAB", "items": [
            "internal/cache/store.go • internal/cache/repo.go",
            "internal/cache/patterns.go • stampede.go",
            "tests/cache_test.go • cmd/demo/main.go • README.md",
        ]},
        {"title": "VERIFIKASI & STATUS", "items": [
            f"Sumber: commit {COMMIT_SHA[:8]}",
            "Test passing dan diverifikasi (Race-Clean Go 1.22).",
            "Invarian terbukti: 20 goroutine naive → 20 DB queries; dengan SingleFlight → 1 DB query.",
        ], "type": "bullet"},
        {"title": "TAKEAWAYS", "items": [
            "Cache-aside adalah default pragmatis; urutan write→invalidate kritis.",
            "Write-through memberi read-after-write sinkron; write-behind mengorbankan durability.",
            "Singleflight + XFetch + SWR + Jitter membentukkan defense-in-depth terhadap stampede.",
        ], "type": "bullet"},
    ],
    mental_model="No cache is perfect—combine lazy, coalesced, and probabilistic refresh for resilience.",
    source_path=REPO_URL,
    slide_num="24 / 24"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
