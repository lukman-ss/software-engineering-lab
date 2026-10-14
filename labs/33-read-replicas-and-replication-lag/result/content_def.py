#!/usr/bin/env python3
"""
Publication Builder Definition for Software Engineering Lab 33: Read Replicas and Replication Lag.
Defines metadata, colors, 24 slides, and helper structures.
DO NOT DELETE THIS FILE.
"""

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/33-read-replicas-and-replication-lag/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/33-read-replicas-and-replication-lag"
COMMIT_SHA = "cfba3b76824f0ee03b5fa526462d1b72858144e7"
LAB_NO = "33"
SLUG = "read-replicas-and-replication-lag"
TITLE_ID = "Read Replicas & Replication Lag"
SUBTITLE_ID = "Mitigasi Stale Read & Menjaga Konsistensi Read-Your-Own-Writes"
TITLE_EN = "Read Replicas & Replication Lag"
SUBTITLE_EN = "Mitigating Stale Reads and Enforcing Read-Your-Own-Writes Guarantees"
AUTHOR = "Lukman / lukman-ss"

# Brand color constants
DARK_NAVY = "0b1b2d"
TEAL = "087f76"
CREAM = "f4f1e8"
WHITE = "ffffff"
SLATE = "52606b"
ORANGE = "d97706"

DEVTO_DESC = "Explore read replicas, replication lag anomalies, and session guarantees (read-your-own-writes) in Go with causal tokens and sticky routing."
DEVTO_TAGS = "database, go, systemdesign, backend"
MEDIUM_DESC = "Panduan lengkap mitigasi replication lag pada read replica database: stale read anomaly, sticky session, causal token LSN, dan synchronous replication."

SLIDES = []

def add_slide(categories, title, subtitle, content_items, mental_model, source_path="", slide_num="", title_en="", narrative_en="", narrative=""):
    SLIDES.append({
        "categories": categories,
        "title": title,
        "subtitle": subtitle,
        "content_items": content_items,
        "mental_model": mental_model,
        "source": source_path,
        "slide_num": slide_num,
        "title_en": title_en or title,
        "narrative_en": narrative_en,
        "narrative": narrative,
    })

# Slide 1: Cover (Dark Navy)
add_slide(
    categories=["COVER"],
    title=TITLE_ID,
    subtitle=SUBTITLE_ID,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "Retro CRT monitor with scanlines and glowing database WAL stream",
            "Primary node broadcasting monotonic LSN to lagging replicas",
            f"Author: {AUTHOR}",
        ]}
    ],
    mental_model="",
    source_path="",
    slide_num=""
)

# Slide 2: Problem Definition (Stale Read Anomaly)
add_slide(
    categories=["MASALAH SISTEM"],
    title="Anomali Stale Read pada Read Replica",
    subtitle="User menulis data baru, namun query berikutnya membaca data lama atau kosong",
    content_items=[
        {"title": "GEJALA ANOMALI", "items": [
            "User melakukan write profil atau checkout pesanan ke primary node",
            "Request read berikutnya dialihkan round-robin ke replica",
            "Replica belum menerapkan WAL entry terbaru, menghasilkan data basi atau 404",
        ], "type": "bullet"},
        {"title": "AKAR TEKNIS", "items": [
            "Replikasi asinkron default pada PostgreSQL, MySQL, Aurora, dan MongoDB",
            "Ada jeda propagasi jaringan dan waktu commit WAL antara master dan replica",
            "Replica menjamin eventual consistency, bukan linearizability real-time",
        ], "type": "bullet"},
        {"title": "DAMPAK PENGGUNA", "items": [
            "User mengira transaksi gagal lalu mengklik tombol submit berulang kali",
            "Saldo rekening atau item keranjang belanja terlihat hilang sesaat",
            "Inkonsistensi session yang merusak kepercayaan pengguna terhadap sistem",
        ], "type": "bullet"},
    ],
    mental_model="Replication lag is the distance between primary commit and replica apply.",
    source_path="labs/33-read-replicas-and-replication-lag/README.md",
    slide_num="02 / 24",
    title_en="Stale Read Anomalies on Read Replicas",
    narrative_en="When scaling database reads across secondary replicas, asynchronous replication creates a propagation window where replicas lag behind the primary. An immediate subsequent read by the writer can hit a replica that has not yet applied the latest write-ahead log entry, returning stale records or false not-found errors.",
    narrative="Ketika arsitektur basis data memisahkan beban tulis ke primary dan beban baca ke beberapa replika sekunder, replikasi asinkron menimbulkan celah propagasi. Pengguna yang baru saja memperbarui data dapat langsung diarahkan ke replika yang belum menerapkan perubahan tersebut. Akibatnya, sistem menyajikan data lama atau bahkan melaporkan data tidak ditemukan, memicu kebingungan operasional dan duplikasi transaksi."
)

# Slide 3: Asynchronous Replication Anatomy
add_slide(
    categories=["MEKANISME"],
    title="Anatomi Aliran Asynchronous Replication",
    subtitle="Write commit langsung di primary sementara replica mengejar secara background",
    content_items=[
        {"title": "ALIRAN WAL ENTRY", "items": [
            "1. Primary menerima write, catat data, naikkan LSN monotonik",
            "2. Primary kembalikan status commit sukses ke client tanpa menunggu replica",
            "3. WAL entry disalurkan via stream jaringan ke antrean worker replica",
            "4. Worker replica memproses dan menerapkan entry ke memory/disk lokal",
        ], "type": "diagram"},
        {"title": "PERSAMAAN LAG", "items": [
            "Replication Lag = LSN_primary - AppliedLSN_replica",
            "Lag Waktu (Δt) = Waktu Terima Primary - Waktu Apply Replica",
            "Kapasitas antrean WAL dan beban I/O lokal menentukan besaran lag",
        ], "type": "bullet"},
        {"title": "TRADE-OFF UTAMA", "items": [
            "Keuntungan: Latensi write sangat rendah dan throughput primary maksimal",
            "Konsekuensi: Jendela waktu staleness tidak terhindarkan bagi replica read",
        ], "type": "bullet"},
    ],
    mental_model="Async replication optimizes write latency by accepting eventual consistency on reads.",
    source_path="labs/33-read-replicas-and-replication-lag/internal/cluster/cluster.go:193",
    slide_num="03 / 24",
    title_en="Anatomy of Asynchronous Replication Streams",
    narrative_en="In asynchronous replication mode, the primary commits transactions locally, assigns a monotonically increasing Log Sequence Number, and immediately acknowledges the client. Write-ahead log records are queued and streamed to replicas in the background, making replication lag structurally inevitable under heavy traffic.",
    narrative="Dalam mode replikasi asinkron, primary memproses transaksi, mencatat entri Write-Ahead Log dengan nomor urut monotonik LSN, lalu seketika mengembalikan respons sukses ke aplikasi. Salinan log dikirimkan ke worker replika di latar belakang. Karena proses pengiriman dan penerapan log bergantung pada latensi jaringan serta antrean lokal, perbedaan antara LSN primary dan applied LSN replika selalu muncul."
)

# Slide 4: Read-Your-Own-Writes Guarantee
add_slide(
    categories=["KONSISTENSI"],
    title="Garansi Read-Your-Own-Writes",
    subtitle="User harus selalu melihat update yang baru saja dilakukannya sendiri",
    content_items=[
        {"title": "DEFINISI FORMAL", "items": [
            "Diformalkan oleh Terry et al. (1994) sebagai jaminan sesi terdistribusi",
            "Jika operasi W1 dilakukan oleh sesi S pada waktu t, maka semua read R",
            "berikutnya oleh sesi S pada t' > t wajib mengamati efek dari W1",
        ], "type": "bullet"},
        {"title": "SKOP GARANSI", "items": [
            "Berlaku khusus untuk sesi atau user yang melakukan penulisan",
            "User lain diperbolehkan melihat data yang sedikit tertinggal (stale)",
            "Menghindari kebutuhan linearizability penuh pada seluruh cluster",
        ], "type": "bullet"},
        {"title": "STRATEGI PENYELESAIAN", "items": [
            "Time-Based Sticky: Paksa read ke primary selama jendela waktu tertentu",
            "Causal Token / Min LSN: Tunggu replika mengejar LSN commit transaksi",
            "Synchronous Mode: Terapkan perubahan ke semua replika sebelum commit",
        ], "type": "bullet"},
    ],
    mental_model="Read-your-own-writes is a session guarantee, not global linearizability.",
    source_path="labs/33-read-replicas-and-replication-lag/content/02-master-draft.md:20",
    slide_num="04 / 24",
    title_en="The Read-Your-Own-Writes Consistency Guarantee",
    narrative_en="Read-your-own-writes ensures that an individual user session immediately observes its own updates, even if global readers observe eventual consistency. Formally proven by Terry et al., this session-scoped constraint prevents user-facing anomalies without the crippling overhead of global consensus.",
    narrative="Garansi read-your-own-writes memastikan pengguna selalu dapat membaca hasil perubahannya sendiri sesaat setelah commit berhasil. Berbeda dari konsistensi ketat global, aturan ini dibatasi pada level sesi pengguna. Pengguna lain di belahan sistem lain masih dapat menerima data eventual consistency, sehingga sistem tetap memiliki skalabilitas baca yang tinggi tanpa mengorbankan pengalaman pengguna penulis."
)

# Slide 5: Why Naive Read-Split Fails
add_slide(
    categories=["ANALISIS KEGAGALAN"],
    title="Mengapa Naive Read-Split Gagal",
    subtitle="Memisahkan write ke primary dan read round-robin ke replica memicu bug",
    content_items=[
        {"title": "KODE READ NAIF", "items": [
            "func (r *Router) ReadNaive(key string) {",
            "    target := replicas[rrIndex.Add(1) % len(replicas)]",
            "    return target.Read(key)",
            "}",
        ], "type": "code"},
        {"title": "URUTAN KEGAGALAN", "items": [
            "1. Client Write('user:profile:123', 'Alice') -> Primary LSN 1",
            "2. Client ReadNaive('user:profile:123') -> Replika-2 (Applied LSN 0)",
            "3. Replika-2 mengembalikan ErrNotFound karena WAL belum di-apply",
        ], "type": "bullet"},
        {"title": "TITIK BUTA PENGEMBANG", "items": [
            "Di lingkungan lokal, lag bernilai ~0ms sehingga bug tidak pernah muncul",
            "Di staging, volume data rendah menyembunyikan antrean replikasi",
            "Di produksi dengan network latency dan disk I/O, bug meledak sporadis",
        ], "type": "bullet"},
    ],
    mental_model="Naive round-robin blindly routes reads to replicas without freshness verification.",
    source_path="labs/33-read-replicas-and-replication-lag/internal/router/router.go:67",
    slide_num="05 / 24",
    title_en="Why Naive Read-Splitting Strategies Fail",
    narrative_en="A naive read-write splitter routes all writes to the primary and distributes all reads round-robin across replicas without inspecting replication status. In local development environments with zero network delay, this appears flawless. In production under real I/O pressure, it causes intermittent not-found and stale read regressions.",
    narrative="Pola pemisahan read-write konvensional mengarahkan seluruh operasi tulis ke primary dan membagi query baca secara round-robin ke seluruh replika tanpa memeriksa status replikasi. Pada mesin pengembangan lokal dengan latensi nol, pola ini tampak berjalan sempurna. Namun di lingkungan produksi dengan beban I/O nyata, replika tertinggal beberapa ratus milidetik dan memicu anomali baca kosong."
)

# Slide 6: Architecture Overview
add_slide(
    categories=["ARSITEKTUR"],
    title="Arsitektur Lab 33: Tiga Layer Utama",
    subtitle="Cluster simulator in-memory, Dynamic Router, dan Test Verifier",
    content_items=[
        {"title": "KOMPONEN CLUSTER", "items": [
            "Node: Menyimpan map KV, appliedLSN, dan antrean channel WAL",
            "Primary: Atomic increment currentLSN, eksekusi write real-time",
            "Replica: Worker goroutine dengan delay lag terkonfigurasi dan sync.Cond",
        ], "type": "bullet"},
        {"title": "DYNAMIC ROUTER", "items": [
            "Session State: Tracking LastWriteTime dan LastWriteLSN via sync.Map",
            "Strategi Routing: Naive, Sticky Session, Causal Token, Lag-Aware",
            "Primary Fallback: Pengalihan otomatis saat replika melebihi batas toleransi",
        ], "type": "bullet"},
        {"title": "VERIFIKASI SISTEM", "items": [
            "Test Suite deterministik: 7 skenario pengujian unit dan integrasi",
            "Race detector: Bebas data race pada ratusan goroutine konkuren",
        ], "type": "bullet"},
    ],
    mental_model="Layered separation decouples cluster replication mechanics from client routing decisions.",
    source_path="labs/33-read-replicas-and-replication-lag/engineering/01-design.md:32",
    slide_num="06 / 24",
    title_en="Architecture Overview: A Three-Layer Simulation System",
    narrative_en="The lab architecture isolates cluster mechanics from client-side routing intelligence. The Cluster layer handles concurrent WAL broadcasting and replica worker routines, while the Router layer maintains session metadata and evaluates freshness policies before executing physical reads.",
    narrative="Arsitektur sistem Lab 33 dirancang memisahkan mekanisme replikasi fisik dari lapisan perutean query klien. Lapisan Cluster mensimulasikan penyimpanan data in-memory dengan nomor urut log LSN dan goroutine worker replika. Lapisan Router memelihara status sesi pengguna dan menjalankan strategi pemilihan simpul baca berdasarkan kebijakan kesegaran data yang diinginkan."
)

# Slide 7: Primary Node & Atomic LSN
add_slide(
    categories=["BUKTI KODE"],
    title="Primary Node & Monotonic LSN Tracking",
    subtitle="Setiap operasi tulis menaikkan sequence number global secara atomic",
    content_items=[
        {"title": "STRUKTUR WAL & NODE", "items": [
            "type WALEntry struct {",
            "    LSN       uint64",
            "    Key       string",
            "    Value     string",
            "    Timestamp time.Time",
            "}",
        ], "type": "code"},
        {"title": "EKSEKUSI TULIS ATOMIK", "items": [
            "lsn := c.currentLSN.Add(1)",
            "c.primary.mu.Lock()",
            "c.primary.data[key] = value",
            "c.primary.appliedLSN = lsn",
            "c.primary.mu.Unlock()",
        ], "type": "code"},
        {"title": "KARAKTERISTIK LSN", "items": [
            "Monotonik strictly increasing tanpa gap",
            "Menjadi representasi posisi kemajuan commit log basis data",
            "Dasar penentuan komparasi kesegaran data antar simpul",
        ], "type": "bullet"},
    ],
    mental_model="Monotonic LSN provides a global logical clock for ordering database state.",
    source_path="labs/33-read-replicas-and-replication-lag/internal/cluster/cluster.go:193",
    slide_num="07 / 24",
    title_en="Primary Node Execution and Atomic LSN Generation",
    narrative_en="To establish a rigorous order of mutations, the primary node utilizes an atomic 64-bit unsigned integer counter representing the Log Sequence Number. Every write acquires the lock, mutates local state, assigns the new LSN, and packages a WALEntry for asynchronous replica distribution.",
    narrative="Untuk menetapkan urutan mutasi data yang tidak terbantahkan, simpul primary memanfaatkan pencacah atomik 64-bit yang merepresentasikan Log Sequence Number. Setiap kali transaksi tulis terjadi, primary memperbarui data lokal, mencatat nomor urut LSN terbaru, dan membungkus data tersebut ke dalam paket entri WAL untuk didistribusikan ke simpul replika."
)

# Slide 8: Replica Worker & Lag Simulation
add_slide(
    categories=["BUKTI KODE"],
    title="Replica Worker & Simulasi Lag Terkontrol",
    subtitle="Worker goroutine menerapkan jeda waktu sebelum memperbarui memory dan LSN",
    content_items=[
        {"title": "GOROUTINE WORKER", "items": [
            "for entry := range replica.walChannel {",
            "    if lag > 0 { time.Sleep(lag) }",
            "    replica.mu.Lock()",
            "    replica.data[entry.Key] = entry.Value",
            "    replica.appliedLSN = entry.LSN",
            "    replica.cond.Broadcast()",
            "    replica.mu.Unlock()",
            "}",
        ], "type": "code"},
        {"title": "MEKANISME NOTIFIKASI", "items": [
            "sync.Cond digunakan untuk mem-broadcast event ke reader yang menunggu",
            "Channel buffer 1024 mencegah blocking write pada primary",
            "Lag dapat diubah secara dinamis melalui SetLag(d) saat pengujian",
        ], "type": "bullet"},
        {"title": "KONDISI REPLIKASI", "items": [
            "Selama jeda lag berlangsung, AppliedLSN simpul tetap bernilai lama",
            "Query yang masuk pada rentang ini akan membaca snapshot usang",
        ], "type": "bullet"},
    ],
    mental_model="Replicas lag when ingestion outpaces application speed or network delays intervene.",
    source_path="labs/33-read-replicas-and-replication-lag/internal/cluster/cluster.go:148",
    slide_num="08 / 24",
    title_en="Replica Worker Goroutine and Lag Ingestion Pipeline",
    narrative_en="Each replica runs a dedicated worker goroutine consuming from a buffered WAL channel. Artificial delay simulates network transit and disk synchronization latencies. Once applied, sync.Cond broadcasts trigger waiting causal read queries to wake up deterministically.",
    narrative="Setiap replika menjalankan sebuah goroutine worker yang membaca antrean channel WAL berkapasitas 1024 entri. Durasi lag disimulasikan melalui jeda waktu terkontrol sebelum data disinkronkan ke memori lokal. Begitu data diterapkan, variabel kondisi sync.Cond memicu sinyal siaran guna membangunkan operasi baca yang sedang menunggu pencapaian LSN tertentu."
)

# Slide 9: Strategy 1 - Time-Based Sticky Routing
add_slide(
    categories=["STRATEGI ROUTING"],
    title="Strategi 1: Time-Based Sticky Routing",
    subtitle="Arahkan read ke primary selama jendela waktu T_sticky setelah write",
    content_items=[
        {"title": "LOGIKA SESI STICKY", "items": [
            "if stateVal, ok := r.sessions.Load(sessionID); ok {",
            "    state := stateVal.(SessionState)",
            "    if time.Since(state.LastWriteTime) < r.config.StickyDuration {",
            "        return r.cluster.Primary().Read(key)",
            "    }",
            "}",
            "return r.ReadLagAware(key)",
        ], "type": "code"},
        {"title": "KEUNGGULAN OPERASIONAL", "items": [
            "Implementasi sangat ringan: hanya membandingkan timestamp lokal",
            "Tidak memerlukan query status LSN atau koordinasi antar simpul",
            "Secara drastis mereduksi anomali post-write pada navigasi umum user",
        ], "type": "bullet"},
        {"title": "KELEMAHAN HEURISTIK", "items": [
            "Jendela 5 detik adalah estimasi, bukan jaminan pasti replika sudah catch-up",
            "Beban primary tetap tinggi jika user aktif menulis dan membaca bersamaan",
        ], "type": "bullet"},
    ],
    mental_model="Time-based sticky routing trades primary capacity for simplicity without cluster round-trips.",
    source_path="labs/33-read-replicas-and-replication-lag/internal/router/router.go:80",
    slide_num="09 / 24",
    title_en="Strategy 1: Time-Based Sticky Session Routing",
    narrative_en="Time-based sticky routing tracks the timestamp of each session's last write. Any subsequent read within a configurable TTL (e.g., 5 seconds) routes unconditionally to the primary node. It eliminates stale reads for the writer with zero coordination overhead, though it relies on empirical latency assumptions.",
    narrative="Strategi sticky routing berbasis waktu mencatat stempel waktu penulisan terakhir dari masing-masing sesi pengguna. Setiap operasi baca yang dilakukan dalam rentang jendela waktu TTL akan dialihkan langsung ke primary. Pendekatan ini sangat efisien karena tidak membutuhkan pengecekan status simpul, meskipun bersifat heuristik dan bergantung pada estimasi latensi."
)

# Slide 10: Strategy 2 - Causal Token & LSN Wait
add_slide(
    categories=["STRATEGI ROUTING"],
    title="Strategi 2: Causal Token / Minimum LSN Routing",
    subtitle="Client membawa commit token LSN dan menunggu replica mengejar posisi tersebut",
    content_items=[
        {"title": "PENCARIAN SIMPUL FRESH", "items": [
            "for _, replica := range replicas {",
            "    if replica.AppliedLSN() >= minLSN {",
            "        return replica.Read(key)",
            "    }",
            "}",
        ], "type": "code"},
        {"title": "BLOCKING WAIT & TIMEOUT", "items": [
            "waitCtx, cancel := context.WithTimeout(ctx, r.config.WaitTimeout)",
            "defer cancel()",
            "if err := targetReplica.WaitForLSN(waitCtx, minLSN); err == nil {",
            "    return targetReplica.Read(key)",
            "}",
            "return r.cluster.Primary().Read(key) // Fallback aman",
        ], "type": "code"},
        {"title": "GARANSI MATEMATIS", "items": [
            "Pasti mengamati mutasi transaksi karena AppliedLSN >= minLSN terjamin",
            "Jika replika lambat, fallback ke primary melindungi ketersediaan sistem",
        ], "type": "bullet"},
    ],
    mental_model="Causal token routing guarantees causal consistency by conditioning reads on commit LSN.",
    source_path="labs/33-read-replicas-and-replication-lag/internal/router/router.go:92",
    slide_num="10 / 24",
    title_en="Strategy 2: Causal Token and Minimum LSN Routing",
    narrative_en="Causal token routing binds each write's commit LSN to subsequent read requests. The router evaluates replicas for applied LSN progress. If no replica has reached the target LSN, it executes an event-driven wait with a deadline, falling back to the primary if timeout thresholds expire.",
    narrative="Pendekatan causal token menyertakan token LSN hasil transaksi tulis pada setiap query baca berikutnya. Router memeriksa apakah terdapat replika yang telah mencapai nomor urut tersebut. Jika belum ada, sistem menunggu secara efisien menggunakan variabel kondisi hingga batas waktu tercapai, lalu melakukan fallback ke primary bila replika mengalami keterlambatan parah."
)

# Slide 11: Event-Driven WaitForLSN
add_slide(
    categories=["BUKTI KODE"],
    title="Event-Driven WaitForLSN dengan sync.Cond",
    subtitle="Menghilangkan polling CPU yang boros dengan sinkronisasi variabel kondisi",
    content_items=[
        {"title": "KODE SINKRONISASI", "items": [
            "func (n *Node) WaitForLSN(ctx context.Context, target uint64) error {",
            "    go func() {",
            "        n.mu.Lock(); defer n.mu.Unlock()",
            "        for n.appliedLSN < target {",
            "            n.cond.Wait() // Dibangunkan oleh Broadcast worker",
            "        }",
            "        close(ch)",
            "    }()",
            "    select {",
            "    case <-ctx.Done(): return ctx.Err()",
            "    case <-ch: return nil",
            "    }",
            "}",
        ], "type": "code"},
        {"title": "EFISIENSI KOMPUTASI", "items": [
            "Goroutine reader tertidur tanpa memakan siklus komputasi CPU",
            "Langsung aktif seketika worker replika selesai menerapkan WAL",
            "Mendukung pembatalan timeout context secara aman dan bersih",
        ], "type": "bullet"},
        {"title": "SAFETY CONCURRENCY", "items": [
            "Mutex melindungi pembacaan dan penulisan variabel appliedLSN",
            "Tidak terjadi kebocoran goroutine saat context kadaluarsa",
        ], "type": "bullet"},
    ],
    mental_model="Event-driven condition variables eliminate busy-waiting during replication catch-up.",
    source_path="labs/33-read-replicas-and-replication-lag/internal/cluster/cluster.go:79",
    slide_num="11 / 24",
    title_en="Event-Driven Catch-Up with Condition Variables",
    narrative_en="Rather than wasting CPU cycles on repetitive polling, WaitForLSN suspends the reader goroutine on a sync.Cond condition variable. When the background replica worker applies the transaction, it broadcasts a wake-up signal, immediately resuming the reader with sub-millisecond precision.",
    narrative="Alih-alih membebani CPU dengan mekanisme polling berulang, method WaitForLSN menidurkan goroutine pembaca menggunakan sync.Cond. Ketika worker replika selesai menerapkan transaksi WAL, sinyal siaran dikirimkan untuk membangunkan goroutine secara instan, menghasilkan waktu tunggu yang presisi dan hemat sumber daya."
)

# Slide 12: Strategy 3 - Lag-Aware Dynamic Routing
add_slide(
    categories=["STRATEGI ROUTING"],
    title="Strategi 3: Lag-Aware Dynamic Routing",
    subtitle="Mengecualikan replica yang tertinggal jauh dan fallback saat replica habis",
    content_items=[
        {"title": "FILTERING BERBASIS DELTA LSN", "items": [
            "primaryLSN := r.cluster.CurrentLSN()",
            "for _, replica := range replicas {",
            "    diff := primaryLSN - replica.AppliedLSN()",
            "    if diff <= r.config.MaxLSNDiff {",
            "        eligible = append(eligible, replica)",
            "    }",
            "}",
        ], "type": "code"},
        {"title": "FALLBACK OTOMATIS", "items": [
            "if len(eligible) == 0 {",
            "    return r.cluster.Primary().Read(key) // Fallback-lag",
            "}",
            "idx := r.rrIndex.Add(1) % len(eligible)",
            "return eligible[idx].Read(key)",
        ], "type": "code"},
        {"title": "PROTEKSI SISTEM", "items": [
            "Mencegah query umum membaca data yang terlalu usang",
            "Mengisolasi simpul replika yang mengalami gangguan I/O atau jaringan",
            "Degradasi anggun (graceful degradation) menjaga ketersediaan layanan",
        ], "type": "bullet"},
    ],
    mental_model="Dynamic lag-aware routing preserves consistency SLAs by shedding severely lagging nodes.",
    source_path="labs/33-read-replicas-and-replication-lag/internal/router/router.go:117",
    slide_num="12 / 24",
    title_en="Strategy 3: Lag-Aware Dynamic Routing and Primary Fallback",
    narrative_en="For general reads without explicit session constraints, the lag-aware routing strategy evaluates the delta between primary LSN and replica applied LSN. Nodes exceeding the maximum tolerable difference are temporarily removed from rotation, triggering automated fallback to primary if all replicas degrade.",
    narrative="Untuk operasi baca umum yang tidak terikat sesi penulis, strategi lag-aware menghitung selisih delta antara LSN primary dan LSN terapan replika. Simpul yang melewati batas toleransi MaxLSNDiff dikeluarkan dari antrean beban. Jika seluruh replika mengalami gangguan, router secara otomatis mengalihkan beban baca ke primary demi menjaga integritas data."
)

# Slide 13: Strategy 4 - Synchronous Replication
add_slide(
    categories=["MODE REPLIKASI"],
    title="Strategi 4: Synchronous Replication (remote_apply)",
    subtitle="Write diblokir hingga semua replica menerapkan perubahan secara tuntas",
    content_items=[
        {"title": "EKSEKUSI TULIS SINKRON", "items": [
            "if c.replMode == SyncReplication {",
            "    for _, replica := range c.replicas {",
            "        if lag > 0 { time.Sleep(lag) }",
            "        replica.data[key] = value",
            "        replica.appliedLSN = lsn",
            "    }",
            "}",
        ], "type": "code"},
        {"title": "GARANSI KESEGARAN", "items": [
            "Semua replika seketika memiliki AppliedLSN == PrimaryLSN",
            "Naive read sekalipun dijamin 100% bebas dari anomali stale read",
            "Setara dengan parameter synchronous_commit = remote_apply di PostgreSQL",
        ], "type": "bullet"},
        {"title": "BIAYA PERFORMA TULIS", "items": [
            "Latensi write membengkak sebanding dengan round-trip terlama replika",
            "Satu replika lambat akan memperlambat seluruh proses penulisan sistem",
            "Hanya direkomendasikan untuk sistem perbankan atau data finansial kritis",
        ], "type": "bullet"},
    ],
    mental_model="Synchronous replication achieves immediate read freshness by penalizing write latency.",
    source_path="labs/33-read-replicas-and-replication-lag/internal/cluster/cluster.go:214",
    slide_num="13 / 24",
    title_en="Strategy 4: Synchronous Replication and the Freshness Tax",
    narrative_en="Synchronous replication halts the primary write process until all designated replica nodes acknowledge full local execution. While this completely eliminates stale reads across the entire cluster, it imposes a severe latency penalty on every write operation proportional to network round-trip delays.",
    narrative="Replikasi sinkron menahan konfirmasi transaksi tulis pada primary hingga seluruh simpul replika selesai menerapkan perubahan ke penyimpanan lokal. Pola ini sepenuhnya melenyapkan anomali stale read pada semua replika, namun mengorbankan performa tulis secara drastis karena setiap operasi penulisan terhambat oleh latensi jaringan terjauh."
)

# Slide 14: Quantitative Comparison Table
add_slide(
    categories=["TABEL EVALUASI"],
    title="Perbandingan 4 Strategi Mitigasi Lag",
    subtitle="Evaluasi matriks trade-off antara garansi, latensi, dan beban primary",
    content_items=[
        {"title": "MATRIKS STRATEGI", "items": [
            "Strategi | Garansi | Latensi Read | Overhead Primary | Kompleksitas",
            "Naive Read | Eventual | Terendah | Nol | Sangat Rendah",
            "Sticky Session | Read-Own-Writes | Terendah | Sedang (T_sticky) | Rendah",
            "Causal Token | Causal Session | Sedang (Wait) | Rendah (Fallback) | Sedang",
            "Lag-Aware | Bounded Lag | Rendah | Kondisional | Sedang",
            "Sync Mode | Linearizable | Terendah | Nol (Write Mahal) | Tinggi",
        ], "type": "table"},
        {"title": "ANALISIS KEBUTUHAN", "items": [
            "E-commerce umum: Kombinasi Sticky Session (checkout) + Lag-Aware (katalog)",
            "Fintech / Ledger: Causal Token untuk transaksi atau Synchronous Replication",
            "Analytics Dashboard: Lag-Aware dengan toleransi batas LSN yang longgar",
        ], "type": "bullet"},
        {"title": "ATURAN UTAMA", "items": [
            "Tidak ada solusi perak: pilih strategi berdasarkan nilai toleransi bisnis",
        ], "type": "bullet"},
    ],
    mental_model="Match your replication mitigation strategy to the cost of reading stale state.",
    source_path="labs/33-read-replicas-and-replication-lag/content/02-master-draft.md:62",
    slide_num="14 / 24",
    title_en="Quantitative Strategy Matrix: Trade-Off Evaluation",
    narrative_en="Choosing a mitigation strategy requires balancing consistency guarantees against latency overhead and primary load. While sticky routing offers the lowest operational friction, causal token routing guarantees mathematical correctness without permanently overloading the primary cluster.",
    narrative="Pemilihan strategi penanganan replication lag menuntut keseimbangan antara garansi konsistensi, latensi query, dan beban kerja simpul primary. Sticky routing menawarkan implementasi paling praktis untuk alur pengguna umum, sedangkan causal token memberikan kepastian matematis tanpa membebani primary secara berlebihan."
)

# Slide 15: Demo Execution & Quantitative Proof
add_slide(
    categories=["BUKTI DEMO"],
    title="Hasil Eksekusi Demo CLI Lab 33",
    subtitle="Bukti empiris terjadinya anomali dan keberhasilan masing-masing strategi",
    content_items=[
        {"title": "LOG SKENARIO 1 & 2", "items": [
            "[Write Primary] key='user:profile:123', LSN: 1",
            "[Naive Read] Node: replica-2, Applied LSN: 0, Err: key not found",
            "[Sticky Read (within 500ms)] Node: primary, LSN: 1, Val: Alice",
            "[Sticky Read (after TTL)] Node: replica-1, LSN: 1, Val: Alice",
        ], "type": "code"},
        {"title": "LOG SKENARIO 3 & 4", "items": [
            "[Write Primary] key='order:101', LSN: 2",
            "[Token Read (MinLSN=2)] Node: replica-1, LSN: 2, Val: completed",
            "[Sync Write] LSN: 1, Write Duration: 102.19ms",
            "[Sync Naive Read] Node: replica-2, LSN: 1, Val: Bob, Err: <nil>",
        ], "type": "code"},
        {"title": "INTERPRETASI HASIL", "items": [
            "Skenario 1 membuktikan kegagalan fatal metode naif saat lag aktif",
            "Skenario 2 membuktikan transisi mulus dari primary ke replika setelah TTL",
            "Skenario 3 membuktikan penungguan deterministik hingga replika siap",
            "Skenario 4 membuktikan kesegaran instan dengan penalti durasi tulis 102ms",
        ], "type": "bullet"},
    ],
    mental_model="Empirical demo execution validates theoretical invariants under controlled conditions.",
    source_path="labs/33-read-replicas-and-replication-lag/cmd/demo/main.go",
    slide_num="15 / 24",
    title_en="Interactive CLI Demonstration Execution Logs",
    narrative_en="Running the lab demo CLI provides reproducible evidence of every failure mode and mitigation strategy. Naive reads fail immediately under asynchronous delay, sticky routing shields the author during post-write navigation, causal tokens block cleanly until catch-up, and synchronous writes record a measurable 102ms latency overhead.",
    narrative="Eksekusi demo CLI membuktikan validitas seluruh skenario secara nyata. Pendekatan naif langsung memicu error data tidak ditemukan saat lag aktif, sticky routing berhasil melindungi pembacaan pengguna ke primary, causal token menunggu hingga replika siap secara presisi, dan mode sinkron mencatatkan penalti waktu tulis terukur sebesar 102ms."
)

# Slide 16: Test Suite Verification & Race Detector
add_slide(
    categories=["VERIFIKASI TEST"],
    title="Verifikasi 7 Test Case & Bebas Data Race",
    subtitle="Pengujian unit dan konkurensi membuktikan ketahanan sistem di bawah beban",
    content_items=[
        {"title": "DAFTAR TEST SUITE", "items": [
            "TestNaiveReplicationLag_StaleRead: PASS (0.60s)",
            "TestStickySessionRouting: PASS (0.55s)",
            "TestReadWithToken_LSN: PASS (0.20s)",
            "TestReplicaLagThreshold_Fallback: PASS (0.00s)",
            "TestSynchronousReplication_Freshness: PASS (0.10s)",
            "TestConcurrentAccess_RaceFree: PASS (1.10s)",
            "TestWaitForLSN_ContextTimeout: PASS (0.05s)",
        ], "type": "code"},
        {"title": "UJI KONKURENSI EKSTRIM", "items": [
            "5 Goroutine Writer + 10 Goroutine Reader",
            "Masing-masing 20 operasi paralel (Total 300 operasi simultan)",
            "Dijalankan dengan go test -race tanpa memicu race condition",
        ], "type": "bullet"},
        {"title": "UJI BATAS WAKTU", "items": [
            "TestWaitForLSN membuktikan pengembalian context.DeadlineExceeded",
            "Sistem tidak mengalami kebocoran memori atau dead lock goroutine",
        ], "type": "bullet"},
    ],
    mental_model="Automated test assertions with race detection guarantee thread-safe correctness.",
    source_path="labs/33-read-replicas-and-replication-lag/tests/replication_test.go",
    slide_num="16 / 24",
    title_en="Comprehensive Unit and Concurrency Test Verification",
    narrative_en="All seven test cases pass under Go's race detector. The test suite verifies exact stale read occurrences, sticky session TTL hand-offs, causal LSN assertions, threshold fallback to primary, immediate synchronous visibility, and rigorous thread safety under concurrent writer-reader pressure.",
    narrative="Seluruh 7 skenario pengujian unit lulus sempurna dengan pengaktifan Go race detector. Suite pengujian memvalidasi munculnya anomali stale read, transisi TTL sticky session, penungguan token LSN, fallback akibat lag berlebih, kesegaran replikasi sinkron, serta keamanan akses memori konkuren tanpa data race."
)

# Slide 17: Production Mapping: PostgreSQL & MySQL
add_slide(
    categories=["STANDAR INDUSTRI"],
    title="Pemetaan ke PostgreSQL dan MySQL",
    subtitle="Bagaimana konsep lab berkorespondensi dengan konfigurasi database populer",
    content_items=[
        {"title": "POSTGRESQL PRIMITIVES", "items": [
            "LSN: pg_current_wal_lsn() & pg_last_wal_replay_lsn()",
            "Sync Commit: synchronous_commit = on | remote_apply | off",
            "Lag Metric: pg_stat_replication.replay_lag",
            "Catch-Up Query: pg_wal_replay_pause() / pg_last_wal_receive_lsn()",
        ], "type": "bullet"},
        {"title": "MYSQL / VITESS PRIMITIVES", "items": [
            "GTID: Global Transaction Identifier setara dengan LSN monotonik",
            "Semi-Sync: rpl_semi_sync_master_enabled (durability, bukan visibility)",
            "Lag Metric: Seconds_Behind_Master (kurang akurat) vs Heartbeat GTID",
            "WAIT_FOR_EXECUTED_GTID_SET: Ekivalen native dari WaitForLSN",
        ], "type": "bullet"},
        {"title": "PERINGATAN SEMI-SYNC", "items": [
            "MySQL semi-sync hanya menjamin relay log sampai di replika",
            "Data belum tentu sudah di-apply ke storage engine (bisa tetap stale)",
        ], "type": "bullet"},
    ],
    mental_model="Database engines provide native LSN/GTID primitives; applications orchestrate routing.",
    source_path="labs/33-read-replicas-and-replication-lag/content/02-master-draft.md:38",
    slide_num="17 / 24",
    title_en="Mapping In-Memory Concepts to Production PostgreSQL and MySQL",
    narrative_en="The concepts implemented in this lab directly map to production database mechanisms. PostgreSQL exposes Write-Ahead Log positions via LSN functions, while MySQL leverages Global Transaction Identifiers (GTID) and WAIT_FOR_EXECUTED_GTID_SET to implement deterministic causal synchronization.",
    narrative="Konsep yang disimulasikan pada Lab 33 memiliki padanan langsung pada basis data produksi. PostgreSQL menyediakan fungsi pemantauan LSN fisik dan konfigurasi synchronous_commit, sedangkan MySQL mengandalkan Global Transaction Identifiers (GTID) serta fungsi WAIT_FOR_EXECUTED_GTID_SET untuk sinkronisasi kausal deterministik."
)

# Slide 18: Cloud Managed Services (RDS, Aurora, Azure)
add_slide(
    categories=["CLOUD DATABASE"],
    title="Replication Lag di Cloud Managed Services",
    subtitle="Perbedaan karakteristik arsitektur replikasi shared-storage vs shared-nothing",
    content_items=[
        {"title": "AWS AURORA (SHARED STORAGE)", "items": [
            "Arsitektur storage terdistribusi: Primary dan replica membaca volume sama",
            "Replication lag sangat rendah (umumnya < 20ms)",
            "Lag hanya berupa invalidasi buffer cache memori lokal di replica",
        ], "type": "bullet"},
        {"title": "AWS RDS / AZURE PG (SHARED NOTHING)", "items": [
            "Replikasi streaming berbasis pengiriman byte WAL fisik antar instans",
            "Lag berfluktuasi drastis saat terjadi transaksi batch besar atau spike I/O",
            "Metrik: ReplicaLag (CloudWatch) & physical_replication_delay (Azure)",
        ], "type": "bullet"},
        {"title": "IMPLIKASI DESAIN", "items": [
            "Aurora tetap membutuhkan mitigasi read-your-own-writes pada microsecond",
            "Shared-nothing membutuhkan threshold MaxLSNDiff dan fallback yang ketat",
        ], "type": "bullet"},
    ],
    mental_model="Shared storage reduces replication lag magnitude but does not eliminate cache staleness.",
    source_path="labs/33-read-replicas-and-replication-lag/research/02-sources.md",
    slide_num="18 / 24",
    title_en="Cloud Managed Database Replicas: Aurora vs Traditional RDS",
    narrative_en="Managed cloud databases approach replication differently. AWS Aurora utilizes a distributed log-structured storage volume where replicas observe minimal cache invalidation lag (<20ms). Standard AWS RDS and Azure PostgreSQL rely on traditional network streaming where I/O bottlenecks can induce significant replication lag.",
    narrative="Layanan cloud terkelola memiliki arsitektur replikasi yang berbeda. AWS Aurora memisahkan komputasi dari penyimpanan terdistribusi bersama, sehingga replika hanya mengalami lag pembaruan cache memori dalam orde milidetik rendah. Sebaliknya, RDS standar dan Azure PG menggunakan aliran WAL konvensional yang rentan terhadap lonjakan lag saat lalu lintas data padat."
)

# Slide 19: Distributed Session State in Production
add_slide(
    categories=["ARSITEKTUR SKALA"],
    title="Tantangan Scaling Session State di Produksi",
    subtitle="Transisi dari sync.Map single-process menuju arsitektur stateless",
    content_items=[
        {"title": "BATASAN IN-MEMORY MAP", "items": [
            "Router lab menyimpan sesi pada sync.Map lokal di satu instance aplikasi",
            "Pada arsitektur horizontal (banyak pod API), map lokal tidak tersinkronisasi",
            "Request berikutnya dari user bisa masuk ke pod yang tidak memiliki state sesi",
        ], "type": "bullet"},
        {"title": "SOLUSI 1: DISTRIBUTED CACHE", "items": [
            "Simpan LastWriteTime dan LastWriteLSN di Redis dengan TTL terkonfigurasi",
            "Semua pod API membaca status sesi dari cluster Redis terpusat",
            "Biaya: Tambahan satu round-trip network lookup sebelum eksekusi query",
        ], "type": "bullet"},
        {"title": "SOLUSI 2: CLIENT TOKEN (JWT/COOKIE)", "items": [
            "Kirim commit LSN kembali ke client via header HTTP response",
            "Client mengirimkan kembali LSN via cookie/header pada request berikutnya",
            "Stateless penuh pada backend, router langsung membaca nilai token dari header",
        ], "type": "bullet"},
    ],
    mental_model="Stateless routing pushes causal tokens to client cookies or shared Redis clusters.",
    source_path="labs/33-read-replicas-and-replication-lag/content/02-master-draft.md:190",
    slide_num="19 / 24",
    title_en="Scaling Session State: From In-Memory Maps to Distributed Systems",
    narrative_en="While the lab models session tracking via an in-memory sync.Map, horizontal production environments require distributed state coordination. Production systems either store session write tokens in Redis or encode minimum LSN requirements directly into signed client cookies and HTTP headers.",
    narrative="Implementasi lab menggunakan sync.Map in-memory untuk melacak status sesi pada satu proses aplikasi. Pada sistem produksi berskala horizontal dengan banyak container, status sesi harus didistribusikan. Pendekatan praktis meliputi penyimpanan stempel sesi di Redis atau menyematkan token LSN secara stateless ke dalam cookie terenkripsi milik klien."
)

# Slide 20: Polling Overhead & Heartbeat Optimization
add_slide(
    categories=["OPTIMASI PERFORMA"],
    title="Overhead Polling LSN & Optimasi Heartbeat",
    subtitle="Menghindari query metadata basis data berulang pada setiap request pembacaan",
    content_items=[
        {"title": "MASALAH PER-QUERY LSN CHECK", "items": [
            "Menjalankan query LSN ke database pada setiap read melipatgandakan round-trip",
            "Beban CPU basis data tersita untuk melayani query metadata, bukan query data",
            "Koneksi pool cepat jenuh akibat query administratif yang berulang",
        ], "type": "bullet"},
        {"title": "SOLUSI BACKGROUND HEARTBEAT", "items": [
            "Worker aplikasi melakukan polling lag replica secara berkala (misal tiap 100ms)",
            "Status LSN dan kelaikan simpul dicatat di memory cache lokal aplikasi",
            "Router membuat keputusan perutean murni berdasarkan data cache tanpa I/O",
        ], "type": "bullet"},
        {"title": "HEARTBEAT TABLE TECHNIQUE", "items": [
            "Primary menulis timestamp periodik ke tabel khusus __heartbeat",
            "Replika mereplikasi baris tersebut, router membaca selisih waktu secara pasif",
        ], "type": "bullet"},
    ],
    mental_model="Cache replica lag status asynchronously via heartbeats to avoid per-query metadata overhead.",
    source_path="labs/33-read-replicas-and-replication-lag/research-audit/06-gaps.md",
    slide_num="20 / 24",
    title_en="Mitigating LSN Polling Overhead with Asynchronous Heartbeats",
    narrative_en="Querying replica positions on every client request introduces severe database round-trip overhead. High-throughput architectures deploy background health workers that poll replica status periodically or read dedicated heartbeat tables, allowing routers to make routing decisions entirely in memory.",
    narrative="Mengeksekusi query pemeriksaan LSN ke replika pada setiap request baca dapat memicu overhead round-trip yang masif. Arsitektur berkinerja tinggi mengatasi kendala ini dengan menjalankan worker background asinkron yang memperbarui metrik lag secara periodik ke dalam memori aplikasi, sehingga keputusan routing dapat diambil tanpa latensi tambahan."
)

# Slide 21: Operational Checklist for Database Replicas
add_slide(
    categories=["PANDUAN OPERASIONAL"],
    title="Checklist Penerapan Read Replica di Produksi",
    subtitle="Panduan verifikasi arsitektur sebelum mengaktifkan perutean replika",
    content_items=[
        {"title": "CHECKLIST ARSITEKTUR", "items": [
            "[ ] Identifikasi query kritis pasca-write dan pasang garansi sesi",
            "[ ] Konfigurasi threshold MaxLSNDiff dan aktifkan fallback primary otomatis",
            "[ ] Pasang monitoring metrik lag replika pada persentil p95 dan p99",
            "[ ] Pastikan ukuran connection pool replika memadai untuk lonjakan read",
            "[ ] Tetapkan batas waktu timeout yang ketat saat menunggu replika catch-up",
            "[ ] Uji skenario failover dan isolasi replika yang mengalami disk stall",
        ], "type": "bullet"},
        {"title": "ATURAN THRESHOLD", "items": [
            "Gunakan metrik lag berbasis byte/LSN, bukan hanya unit detik",
            "Alarm harus berbunyi sebelum antrean replika menyebabkan OOM pada worker",
        ], "type": "bullet"},
        {"title": "DISIPLIN REKAYASA", "items": [
            "Selalu siapkan jalur fallback ke primary untuk melindungi ketersediaan",
        ], "type": "bullet"},
    ],
    mental_model="Enforce strict lag SLAs and automated primary fallbacks before scaling reads.",
    source_path="labs/33-read-replicas-and-replication-lag/content/02-master-draft.md:240",
    slide_num="21 / 24",
    title_en="Production Checklist for Database Read Replication",
    narrative_en="Before distributing read traffic across secondary replicas, engineering teams must establish clear operational guardrails. Enforce automated primary fallbacks, monitor p99 replication lag metrics, and ensure post-write query pathways explicitly request session guarantees.",
    narrative="Sebelum mendistribusikan lalu lintas baca ke replika sekunder di lingkungan produksi, tim pengembang wajib menerapkan parameter operasional yang ketat. Siapkan mekanisme fallback otomatis ke primary, pantau metrik lag replikasi pada persentil p99, serta pastikan alur query pasca-tulis dilindungi oleh garansi sesi read-your-own-writes."
)

# Slide 22: Lab Boundaries & Simulation Scope
add_slide(
    categories=["BATAS LAB"],
    title="Batasan Model & Ruang Lingkup Lab 33",
    subtitle="Transparansi batasan simulasi in-memory dibanding database fisik",
    content_items=[
        {"title": "YANG DIMODELKAN (TERSEDIA)", "items": [
            "Aliran data WAL terurut dengan LSN monotonik atomik",
            "Simulasi lag deterministik via channel buffer dan time delay",
            "Empat strategi routing lengkap dengan fallback dan event synchronization",
            "Pengujian konkurensi multithreading bebas race condition",
        ], "type": "bullet"},
        {"title": "BATASAN LAB (BELUM DIMODELKAN)", "items": [
            "Penyimpanan murni in-memory map tanpa persistensi disk / fsync",
            "Tidak memodelkan engine SQL, indexing B-Tree, atau parsing transaksi",
            "Tidak memodelkan topologi cascading replication atau partisi multi-region",
            "Tidak memodelkan algoritma konsensus pemilihan primary (Raft / Paxos)",
        ], "type": "bullet"},
        {"title": "TUJUAN EDUKASI", "items": [
            "Membuktikan logika mitigasi lag tanpa kompleksitas instalasi DB eksternal",
        ], "type": "bullet"},
    ],
    mental_model="Simulations prove architectural invariants cleanly by stripping away disk and network noise.",
    source_path="labs/33-read-replicas-and-replication-lag/content/01-content-brief.md:33",
    slide_num="22 / 24",
    title_en="Implementation Scope and Architectural Boundaries",
    narrative_en="This lab focuses strictly on validating routing logic and session guarantees. By implementing an in-memory key-value cluster with atomic LSN sequencing, it demonstrates replication lag phenomena deterministically without requiring heavyweight database infrastructure or multi-node clustering.",
    narrative="Lab 33 berfokus murni pada pembuktian logika perutean query dan garansi konsistensi sesi. Dengan membangun kluster key-value in-memory yang dilengkapi pencatatan LSN atomik, sistem mampu mendemonstrasikan anomali lag dan solusinya secara deterministik tanpa hambatan instalasi basis data eksternal yang rumit."
)

# Slide 23: Core Takeaways & Architectural Principles
add_slide(
    categories=["KESIMPULAN"],
    title="Prinsip Inti Mengelola Replication Lag",
    subtitle="Ringkasan pemikiran arsitektural untuk sistem basis data terdistribusi",
    content_items=[
        {"title": "1. ASINKRON ADALAH DEFAULT", "items": [
            "Semua database populer berjalan asinkron demi mempertahankan throughput",
            "Mengharapkan konsistensi instan dari replika tanpa strategi adalah ilusi",
        ], "type": "bullet"},
        {"title": "2. GARANSI BERBASIS KEBUTUHAN", "items": [
            "Gunakan Sticky Routing untuk alur pengguna umum yang berbiaya murah",
            "Gunakan Causal Token saat kebenaran matematis mutlak dibutuhkan",
        ], "type": "bullet"},
        {"title": "3. DEGRADASI ANGGUN", "items": [
            "Batasi replika dengan MaxLSNDiff dan selalu sediakan fallback ke primary",
            "Sistem yang baik tetap melayani pembacaan meski seluruh replika tertinggal",
        ], "type": "bullet"},
    ],
    mental_model="Asynchronous replication trades consistency for throughput; route by session when freshness matters.",
    source_path="labs/33-read-replicas-and-replication-lag/content/05-key-takeaways.md",
    slide_num="23 / 24",
    title_en="Core Architectural Principles for Replication Management",
    narrative_en="Asynchronous replication remains the dominant industry standard for high-throughput database systems. Enforcing read-your-own-writes requires thoughtful client routing rather than hoping for zero-lag replicas. Combining lightweight sticky routing with causal token fallbacks ensures both performance and consistency.",
    narrative="Replikasi asinkron tetap menjadi pilihan standar industri untuk mencapai throughput penulisan yang tinggi. Menjaga integritas data pengguna pasca-transaksi menuntut kecerdasan pada lapisan perutean query. Kombinasi antara sticky routing yang ringan dan causal token yang presisi menghasilkan sistem yang cepat sekaligus konsisten."
)

# Slide 24: Closing Slide (Dark Navy)
add_slide(
    categories=["RINGKASAN"],
    title="Eksperimen kecil, interpretasi yang teliti.",
    subtitle=f"Commit: {COMMIT_SHA[:8]} | Tests: 7/7 PASS | Repository: lukman-ss/software-engineering-lab",
    content_items=[
        {"title": "MODUL LAB 33", "items": [
            "internal/cluster/cluster.go (Node, WAL, Replication)",
            "internal/router/router.go (Sticky, Token, Lag-Aware)",
            "tests/replication_test.go (7 Test Cases, Race Free)",
            "cmd/demo/main.go (CLI Interactive Demonstration)",
        ], "type": "bullet"},
        {"title": "VERIFIKASI & METADATA", "items": [
            f"Penulis: {AUTHOR}",
            f"Commit Hash: {COMMIT_SHA}",
            "Status: PASS (go test -v -race ./...)",
            "Model: In-Memory Primary-Replica Simulation",
        ], "type": "bullet"},
    ],
    mental_model="Replication lag is inevitable in async systems; smart routing bridges the gap to consistency.",
    source_path="labs/33-read-replicas-and-replication-lag/README.md",
    slide_num=""
)

if __name__ == "__main__":
    print(f"Total slides defined: {len(SLIDES)}")
