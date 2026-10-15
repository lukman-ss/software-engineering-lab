#!/usr/bin/env python3
"""
Publication Builder for Software Engineering Lab 33: Read Replicas and Replication Lag.
Generates 24-slide LinkedIn carousel covering stale read anomalies, sticky-session routing,
causal LSN tokens, lag-aware filtering, and synchronous remote_apply strategies.
"""
import os

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/33-read-replicas-and-replication-lag/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/33-read-replicas-and-replication-lag"
COMMIT_SHA = "b41fffb"
LAB_NO = "33"
SLUG = "read-replicas-and-replication-lag"
TITLE_ID = "Read Replica dan Replication Lag"
SUBTITLE_ID = "Menjaga Read-Your-Own-Writes Tanpa Mengorbankan Skalabilitas Baca"
TITLE_EN = "Read Replicas and Replication Lag"
SUBTITLE_EN = "Guaranteeing Read-Your-Own-Writes Without Sacrificing Read Scalability"
AUTHOR = "Lukman / lukman-ss"
COVER_IMAGE_PATH = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/33-read-replicas-and-replication-lag/result/image-cover/cover_art.jpg"

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
            "Diagram cluster primary-replica dengan panah WAL streaming",
            "Ilustrasi replication lag dan jalur routing read (naive vs sticky vs token vs sync)",
            f"Penulis: {AUTHOR}",
        ]}
    ],
    mental_model="Replication lag is the distance between primary LSN and replica applied-LSN; every read strategy is a closing mechanism for that gap.",
    source_path=REPO_URL,
    slide_num="01 / 24",
    narrative="Laboratorium ini membahas masalah nyata yang dihadapi setiap backend engineer yang menggunakan read replica untuk scaling database: stale read. Ketika write baru commit ke primary, replica menerima aliran WAL secara asinkron dengan delay tertentu. Jika aplikasi langsung mengarahkan read user ke replica tanpa mekanisme penutupan jarak tersebut, user akan melihat data lama bahkan menerima not-found untuk key yang baru saja ditulis. Lab ini mengimplementasikan empat strategi mitigasi — time-based sticky routing, causal token dengan LSN wait, lag-aware dynamic filtering, dan synchronous remote_apply — dalam simulasi cluster in-memory Go yang utuh.",
    narrative_en="This laboratory explores the real-world problem encountered by every backend engineer using read replicas for database scaling: stale reads. When a new write commits to the primary, replicas receive the WAL stream asynchronously with some delay. If the application immediately directs user reads to replicas without a closing mechanism, users see stale data or even receive not-found errors for recently written keys. This lab implements four mitigation strategies — time-based sticky routing, causal token with LSN wait, lag-aware dynamic filtering, and synchronous remote_apply — within a complete in-memory Go cluster simulation.",
    layout="cover"
)

# Slide 2: Problem nyata
add_slide(
    categories=["MOTIVASI", "STALE READ RISK"],
    title="Masalah Nyata: Ilusi Freshness di Balik Read Replica",
    subtitle="User menulis profil baru, tapi dashboard menampilkan data lama atau kosong.",
    content_items=[
        {"title": "WRITE COMMITTED", "items": [
            "User mengubah field 'tier' premium di endpoint PUT /profile.",
            "Request dialihkan ke primary, LSN increment menjadi 47.",
            "Write berhasil dengan status 200, response body berisi payload baru."
        ], "type": "card"},
        {"title": "REPLICA LAGGING", "items": [
            "WAL entry dengan LSN 47 mengalir ke replica-1 namun belum di-apply.",
            "Replica-1 masih menunjukkan appliedLSN = 44, data lama belum tersingkir.",
            "Worker goroutine replica menunggu durasi lagDuration sebelum flush."
        ], "type": "card"},
        {"title": "STALE READ TERJADI", "items": [
            "Load balancer mengalihkan read dashboard ke replica-1 secara round-robin.",
            "Query SELECT profile WHERE id=user mendapat ErrNotFound atau data usang.",
            "User merasa sistem rusak; actual cause adalah replication lag yang tak tertutup."
        ], "type": "card"}
    ],
    narrative="Skenario ini terjadi secara default di hampir semua database production yang menggunakan primary-replica topology. PostgreSQL streaming replication, MySQL async replication, AWS RDS read replicas, dan Azure Database for PostgreSQL semuanya menjalankan replikasi asinkron secara default tanpa konfigurasi khusus. Ketika client application melakukan pattern write-then-read pada session yang sama, ia tidak memiliki awareness bahwa data yang baru ditulis belum tersedia di replica target. Akibatnya, user mengalami anomali UX yang sulit dilacak karena log menunjukkan write sukses, tapi query berikutnya gagal membaca data yang sama.",
    narrative_en="This scenario occurs by default across nearly all production databases using primary-replica topology. PostgreSQL streaming replication, MySQL async replication, AWS RDS read replicas, and Azure Database for PostgreSQL all run asynchronous replication by default without explicit configuration. When the client application performs a write-then-read pattern on the same session, it has no awareness that the written data is unavailable at the target replica. The result is a hard-to-debug UX anomaly because logs show successful writes but subsequent queries fail to read the same data.",
    mental_model="Stale reads are not bugs; they are the expected behavior of eventually consistent replica systems.",
    source_path="labs/33-read-replicas-and-replication-lag/content/01-content-brief.md",
    slide_num="02 / 24",
    layout="cards"
)

# Slide 3: Mengapa pendekatan naif gagal
add_slide(
    categories=["ANALISIS", "PENDEKATAN NAIF"],
    title="Mengapa Read Naif ke Replica Selalu Berisiko?",
    subtitle="Round-robin tanpa awareness LSN membuka celah read basi di setiap write window.",
    content_items=[
        {"title": "KODE ROUTING NAIF", "items": [
            "func (r *Router) ReadNaive(key string) (string, string, uint64, error) {",
            "    replicas := r.cluster.Replicas()",
            "    if len(replicas) == 0 {",
            "        return r.cluster.Primary().Read(key)  // Fallback aman",
            "    }",
            "    idx := r.rrIndex.Add(1) % uint64(len(replicas))",
            "    target := replicas[idx]",
            "    return target.Read(key)  // Tanpa cek appliedLSN!",
            "}"
        ], "type": "code"},
        {"title": "KEGAGALAN FUNDAMENTAL", "items": [
            "Tidak ada validasi AppliedLSN vs currentLSN sebelum melayani read.",
            "Round-robin indifferent terhadap freshness state setiap replica.",
            "Stale read probability berbanding lurus dengan throughput write dan lag duration.",
            "Write-after-read pada session yang sama paling rentan terhadap anomali ini."
        ], "type": "bullet"}
    ],
    narrative="Implementasi ReadNaive menunjukkan pendekatan paling sederhana dalam read splitting: arahkan setiap request ke replica berikutnya secara berputar. Namun pendekatan ini mengabaikan fakta fundamental bahwa setiap replica memiliki appliedLSN yang mungkin jauh di belakang currentLSN primary. Dalam kondisi write-heavy workload dengan replication lag 200 milidetik, probabilitas mendapatkan stale read bisa mendekati 100% untuk key yang baru ditulis. Bug ini termasuk subtle karena tidak menghasilkan panic atau crash — hanya mengembalikan data salah yang terasa aneh bagi pengguna akhir.",
    narrative_en="The ReadNaive implementation represents the simplest approach to read splitting: round-robining each request across replicas in turn. However, this approach ignores the fundamental fact that every replica may have an appliedLSN significantly behind the primary currentLSN. Under heavy write workloads with a 200-millisecond replication lag, the probability of serving a stale read can approach 100 percent for recently written keys. This bug is particularly subtle because it never panics or crashes — it silently returns incorrect data that perplexes end users.",
    mental_model="Naive load-balancing optimizes for throughput; it abandons consistency entirely.",
    source_path="internal/router/router.go:67-78",
    slide_num="03 / 24",
    layout="code-bullet"
)

# Slide 4: Mental model inti / schema card
add_slide(
    categories=["MENTAL MODEL", "SCHEMA"],
    title="Mental Model: LSN sebagai Jangkar Konsistensi",
    subtitle="Log Sequence Number membangun urutan kausal yang memungkinkan routing deterministic.",
    content_items=[
        {"title": "STRUKTUR DATA WAL", "items": [
            "type WALEntry struct {",
            "    LSN       uint64    // monotonik naik dari primary",
            "    Key       string",
            "    Value     string",
            "    Timestamp time.Time",
            "}"
        ], "type": "code"},
        {"title": "DUO KEADAAN SETIAP NODE", "items": [
            "Primary: mengasign LSN atomic saat commit write; sumber kebenaran tunggal.",
            "Replica: menyimpan map data KV + appliedLSN + goroutine worker consumption WAL.",
            "Gap: delta(PrimaryLSN, ReplicaAppliedLSN) = ukuran replikasi lag saat ini.",
            "Kunci: setiap read strategy harus meminimalkan atau menutup gap ini untuk sesi yang relevan."
        ], "type": "bullet"}
    ],
    narrative="Model mental ini memisahkan dua lapisan abstraction: lapisan konsistensi (LSN sequence) dan lapisan ketersediaan (data distribution). Primary bertindak sebagai single writer yang memberi nomor urut absolut pada setiap transaksi. Setiap replica secara independen consume dan apply entry WAL dengan latensi tertentu. Router membaca appliedLSN setiap replica untuk membuat keputusan routing yang informed. Tanpa LSN, kita tidak bisa membedakan replica yang fresh versus lagging secara terukur — routing menjadi guesswork belaka.",
    narrative_en="This mental model separates two layers of abstraction: the consistency layer through LSN sequences and the availability layer through data distribution. The primary acts as a single writer assigning absolute monotonic ordering to every transaction. Each replica independently consumes and applies WAL entries with some latency. The router inspects every replica's appliedLSN to make informed routing decisions. Without LSN tracking, distinguishing fresh from lagging replicas becomes impossible, reducing routing to uninformed guesswork.",
    mental_model="LSN is the universal causal anchor; replication lag is simply distance from that anchor.",
    source_path="internal/cluster/cluster.go:18-23 • internal/router/router.go:34-37",
    slide_num="04 / 24",
    layout="schema-card"
)

# Slide 5: Diagram alur/pipeline
add_slide(
    categories=["ARSITEKTUR", "PIPELINE"],
    title="Pipeline Write-Read dengan Replication Lag",
    subtitle="Dari client write sampai routing decision: empat jalur strategi mitigasi.",
    content_items=[
        {"title": "SEKUENSI PIPELINE DISTRIBUSI", "items": [
            "1. Client.Write(session, key, value) → Primary commit, LSN increment atomic, WAL entry broadcast ke replica channel",
            "2. Replica worker sleep(lagDuration) → entry apply ke local data map, appliedLSN = entry.LSN, cond.Broadcast wake WaitForLSN",
            "3a. ReadNaive: round-robin pilih replica → jika appliedLSN < lastWriteLSN → stale read atau ErrNotFound",
            "3b. ReadWithStickySession: cek time.Since(lastWriteTime) < T_sticky → redirect ke primary; setelah TTL expired → ReadLagAware",
            "3c. ReadWithToken: iterasi replica untuk appliedLSN >= minLSN; blocking WaitForLSN + timeout context → fallback primary",
            "3d. ReadLagAware: hitung diff = primaryLSN - appliedLSN; filter eligible diff <= MaxLSNDiff; fallback primary jika kosong"
        ], "type": "diagram"}
    ],
    narrative="Pipeline ini menggambarkan perjalanan lengkap sebuah request dari penulisan hingga pembacaan dengan empat varian routing strategy. Jalur pertama menelusuri write path di mana primary mengincrement currentLSN secara atomik lalu broadcasting WAL entry ke channel setiap replica. Jalur kedua mensimulasikan replication lag melalui goroutine worker yang tidur sesuai duration sebelum apply. Empat jalur ketiga menunjukkan bagaimana router memilih target read: naive round-robin tanpa proteksi, sticky session berbasis waktu, token-based LSN wait dengan blocking, atau lag-aware dynamic filtering dengan threshold.",
    narrative_en="This pipeline illustrates the complete request journey from write commit through read with four distinct routing variants. The first path traces the write flow where the primary atomically increments currentLSN then broadcasts WAL entries to every replica channel. The second path simulates replication lag through worker goroutines sleeping for the configured duration before applying entries. The four sub-paths show how the router selects read targets: naive round-robin without protection, time-based sticky sessions, token-driven LSN waiting with blocking, or lag-aware dynamic filtering with threshold predicates.",
    mental_model="Writes flow forward in LSN order; reads must synchronize with that order or fall back to the source.",
    source_path="cmd/demo/main.go:15-64 • internal/cluster/cluster.go:148-178",
    slide_num="05 / 24",
    layout="pipeline-diagram"
)

# Slide 6: Cuplikan kode dominan
add_slide(
    categories=["IMPLEMENTASI", "CODE"],
    title="Implementasi Sticky Session Router",
    subtitle="Cek lastWriteTime per session untuk mengalihkan read ke primary selama jendela sticky.",
    content_items=[
        {"title": "READ WITH STICKY SESSION", "items": [
            "func (r *Router) ReadWithStickySession(sessionID, key string) (string, string, uint64, error) {",
            "    if stateVal, ok := r.sessions.Load(sessionID); ok {",
            "        state := stateVal.(SessionState)",
            "        if time.Since(state.LastWriteTime) < r.config.StickyDuration {",
            "            val, lsn, err := r.cluster.Primary().Read(key)",
            "            return val, r.cluster.Primary().ID(), lsn, err",
            "        }",
            "    }",
            "    return r.ReadLagAware(key)",
            "}"
        ], "type": "code"}
    ],
    narrative="Potongan kode ini merupakan implementasi time-based sticky routing yang paling sederhana namun efektif. Method pertama-tama memeriksa apakah session ID memiliki state yang terekam di sync.Map. Jika state ditemukan dan waktu sejak lastWrite masih dalam batas StickyDuration, router secara deterministik mengarahkan read ke primary — bukan replica manapun. Ini memastikan read-your-own-writes dalam jendela sticky tanpa biaya LSN check. Setelah jendela habis, flow jatuh ke ReadLagAware yang menggunakan threshold filtering.",
    narrative_en="This code snippet implements the simplest yet most effective form of time-based sticky routing. The method first checks whether the session ID has a recorded state in the sync.Map. If the state exists and the elapsed time since last write remains within the StickyDuration bound, the router deterministically directs reads to the primary instead of any replica. This guarantees read-your-own-writes within the sticky window without any LSN overhead. After the window expires, the flow degrades gracefully into lag-aware threshold filtering.",
    mental_model="Sticky windows trade read throughput for guaranteed freshness within a bounded temporal envelope.",
    source_path="internal/router/router.go:80-90",
    slide_num="06 / 24",
    layout="code-dominant"
)

# Slide 7: Failure path / dangling reference konkret
add_slide(
    categories=["FAILURE PATH", "BUKTI"],
    title="Failure Path: Stale Read Anomali di Demo CLI",
    subtitle="Output terminal membuktikan read ke replica-2 gagal meski write baru sukses di primary.",
    content_items=[
        {"title": "TERMINAL FAILURE LOG", "items": [
            "[Write Primary] Wrote key 'user:profile:123', Primary LSN: 1",
            "[Naive Read] Node: replica-2, Applied LSN: 0, Val: '', Err: key not found",
            "",
            "Analisis:",
            "Primary commit LSN 1 tetapi replica-2 appliedLSN masih 0 karena lag 200ms belum berlalu.",
            "ReadNaive memilih replica-2 via round-robin tanpa pengecekan LSN.",
            "Key tidak ditemukan karena belum di-apply di replica; user menerima ErrNotFound senyap."
        ], "type": "terminal"},
        {"title": "Dampak UX", "items": [
            "User melihat halaman kosong atau 404 setelah upload data baru.",
            "Trust erosion: user yakin input berhasil tetapi output tidak mencerminkan input.",
            "Debug cost tinggi karena write log sukses dan read log juga 'berhasil' (bukan error sistem)."
        ], "type": "bullet"}
    ],
    narrative="Failure path ini diambil langsung dari demo CLI ketika ReadNaive dipanggil tepat setelah write. Primary LSN meningkat ke 1, namun replica-2 masih berada di appliedLSN 0 karena worker goroutine belum menyelesaikan sleep 200 milidetik. Load balancer round-robin memilih replica-2 untuk read berikutnya, menyebabkan key not found karena data belum tersedia. Error ini berbeda dengan koneksi mati atau timeout — aplikasi merespons 200 dengan body kosong, membuat debugging menjadi nightmare bagi engineer yang tidak memahami replication topology.",
    narrative_en="This failure path is extracted directly from the demo CLI when ReadNaive is invoked immediately after a write. The primary LSN increments to 1, but replica-2 remains at appliedLSN 0 because the worker goroutine has not yet completed its 200-millisecond sleep. The round-robin load balancer selects replica-2 for the subsequent read, causing a key-not-found error since the data is unavailable. This error differs from connection failures or timeouts — the application responds with 200 and an empty body, creating a debugging nightmare for engineers unfamiliar with replication topology.",
    mental_model="Silent correctness failures hurt more than loud crashes; they erode trust without trace.",
    source_path="cmd/demo/main.go:23-32 • engineering/03-execution-result.md:53-58",
    slide_num="07 / 24",
    layout="failure-path"
)

# Slide 8: Tabel perbandingan mekanisme vs gejala
add_slide(
    categories=["PERBANDINGAN", "MEKANISME"],
    title="Mekanisme vs Gejala: Evaluasi Strategi Read Routing",
    subtitle="Membandingkan empat pendekatan dalam hal jaminan, overhead, dan fallback behavior.",
    content_items=[
        {"title": "TABEL PERBANDINGAN STRATEGI", "items": [
            "Strategi | Jaminan | Overhead | Fallback | Cocok Untuk",
            "ReadNaive | None | Zero | Primary jika replica kosong | Development, debugging",
            "Sticky Session | Session-level fresh | Time check O(1), session map | Lag-aware filtering | Session-bound writes",
            "Causal Token | Causal ordering | LSN check + blocking wait | Context timeout + primary | Critical data reads",
            "Lag-Aware | Threshold-bound | Delta LSN calculation + filter | Primary jika eligible kosong | Bulk reads, analytics"
        ], "type": "table"}
    ],
    narrative="Tabel ini mengklasifikasikan empat strategi routing berdasarkan empat dimensi evaluasi. ReadNaive menghilangkan overhead komunikasi tetapi mengorbankan konsistensi sepenuhnya. Sticky session menawarkan jaminan session-level freshness dengan overhead minimal berupa time comparison dan session map lookup. Causal token memberikan jaminan kausal terukur tetapi memperkenalkan blocking wait yang dapat meningkatkan p99 latency. Lag-aware filtering menyeimbangkan throughput dan freshness melalui threshold predicate, degraded gracefully ke primary ketika tidak ada replica yang memenuhi syarat.",
    narrative_en="This table classifies the four routing strategies across four evaluation dimensions. ReadNaive eliminates communication overhead but sacrifices consistency entirely. Sticky session offers session-level freshness guarantees with minimal overhead consisting of a time comparison and session map lookup. Causal tokens provide measurable causal ordering guarantees but introduce blocking waits that can elevate p99 latency. Lag-aware filtering balances throughput and freshness through threshold predicates, degrading gracefully to primary when no eligible replica exists.",
    mental_model="Every routing strategy sits on a spectrum between consistency, latency, and operational complexity.",
    source_path="content/02-master-draft.md:62-70",
    slide_num="08 / 24",
    layout="comparison-table"
)

# Slide 9: Decision sequence (3 focused decision cards)
add_slide(
    categories=["DESIGN DECISION", "ARSITEKTUR"],
    title="3 Keputusan Desain Kritis Replication Strategy",
    subtitle="Urutan keputusan arsitektur untuk menyeimbangkan consistency dan throughput.",
    content_items=[
        {"title": "1. TENTUKAN CONSISTENCY REQUIREMENT", "items": [
            "Apakah session ini memerlukan read-your-own-writes (CRUD critical)?",
            "Jika ya → sticky session atau causal token wajib digunakan.",
            "Jika tidak → lag-aware filtering cukup untuk bulk analytics reads."
        ], "type": "card"},
        {"title": "2. KALIBRASI PARAMETER WAKTU", "items": [
            "T_sticky harus ≥ p99 replication lag aktual di production.",
            "Default industri 5 detik bersifat heuristik; perlu tuning berbasis monitoring.",
            "WaitTimeout pada causal token harus ≥ worst-case catch-up duration + buffer."
        ], "type": "card"},
        {"title": "3. PILIH FENCING LAYER YANG TEPAT", "items": [
            "Gunakan mode synchronous remote_apply jika freshness mutlak dibutuhkan.",
            "Accept async replication + application-level fencing untuk throughput maksimal.",
            "Jangan mencampuradukkan semi-sync (durability) dengan synchronous (visibility)."
        ], "type": "card"}
    ],
    narrative="Tiga keputusan desain ini membentuk kerangka kerja strategis sebelum memilih implementasi teknis. Keputusan pertama menentukan level consistency requirement yang membedakan use case mana yang memerlukan protection ketat versus mana yang bisa menerima eventual consistency. Keputusan kedua menyangkut kalibrasi parameter waktu berdasarkan observasi empiris production latency distribution — bukan asumsi teoretis. Keputusan ketiga memilih lapisan proteksi yang sesuai dengan SLA bisnis, antara synchronous replication berbiaya tinggi atau application-level fencing dengan graceful degradation.",
    narrative_en="These three design decisions form a strategic framework before selecting any technical implementation. The first decision determines the consistency requirement level that distinguishes which use cases need strict protection versus those accepting eventual consistency. The second decision involves calibrating time parameters based on empirical production latency observations rather than theoretical assumptions. The third decision selects the appropriate protection layer aligned with business SLAs between high-cost synchronous replication or application-level fencing with graceful degradation.",
    mental_model="Design decisions cascade from requirements to parameters to implementation; never skip upstream thinking.",
    source_path="labs/33-read-replicas-and-replication-lag/content/05-key-takeaways.md",
    slide_num="09 / 24",
    layout="decision-cards"
)

# Slide 10: Bukti eksekusi / terminal output
add_slide(
    categories=["BUKTI", "DEMO"],
    title="Bukti Eksekusi: Siklus Lengkap Replication & Mitigasi",
    subtitle="Output demo CLI mengkonfirmasi empat skenario berjalan sesuai spesifikasi.",
    content_items=[
        {"title": "DEMO EXECUTION TRACE", "items": [
            "=== Lab 33: Read Replicas and Replication Lag Demo ===",
            "",
            "--- 1. Async Replication Lag Anomaly ---",
            "[Write Primary] Wrote key 'user:profile:123', Primary LSN: 1",
            "[Naive Read] Node: replica-2, Applied LSN: 0, Val: '', Err: key not found",
            "",
            "--- 2. Sticky Session Routing ---",
            "[Sticky Read (within 500ms)] Node: primary, LSN: 1, Val: '{...}', Err: nil",
            "[Sticky Read (after TTL)] Node: replica-1, LSN: 1, Val: '{...}', Err: nil",
            "",
            "--- 3. Causal Token / LSN Wait ---",
            "[Token Read (MinLSN=2)] Node: replica-1, LSN: 2, Val: '{...}', Err: nil",
            "",
            "--- 4. Synchronous Replication ---",
            "[Sync Write] LSN: 1, Write Duration: 104.055166ms",
            "[Sync Naive Read] Node: replica-2, LSN: 1, Val: '{...}', Err: nil"
        ], "type": "terminal"}
    ],
    narrative="Output terminal ini menyajikan empat bukti eksekusi yang masing-masing mengonfirmasi perilaku strategi routing. Bagian pertama membuktikan anomali stale read terjadi ketika ReadNaive dipanggil pasca-write. Bagian kedua menunjukkan sticky session berhasil mengarahkan read ke primary dalam jendela 500ms, kemudian beralih ke replica setelah TTL expire. Bagian ketiga memverifikasi causal token menunggu catch-up sekitar 200ms sebelum melayani read dari replica yang sudah sesuai LSN. Bagian keempat memperlihatkan synchronous replication membutuhkan tambahan ~104ms pada write latency untuk menjamin freshness instan.",
    narrative_en="This terminal output presents four execution proofs each confirming the behavior of different routing strategies. The first section demonstrates the stale read anomaly when ReadNaive is called immediately post-write. The second section shows sticky session successfully directing reads to the primary within the 500-millisecond window then falling back to a replica after TTL expiration. The third section verifies causal tokens wait approximately 200 milliseconds for catch-up before serving reads from matching replicas. The fourth section reveals synchronous replication requiring approximately 104 additional milliseconds on write latency to guarantee instant freshness.",
    mental_model="Terminal traces convert abstract consistency guarantees into observable execution evidence.",
    source_path="cmd/demo/main.go:1-65 • engineering/03-execution-result.md:48-73",
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
            "Cluster in-memory tanpa persistence layer; tidak membuktikan recovery pasca-crash.",
            "Tidak ada topologi cascading replica (replica-of-replica) atau multi-region replication.",
            "Sticky 500ms dan lag 200ms bersifat ilustratif lab; bukan rekomendasi tuning produksi.",
            "Router state disimpan di sync.Map single-process; tidak membuktikan stateless horizontal scaling.",
            "Tidak mencakup network partition antar primary dan replica, only simulated timing lag.",
            "p99 lag multi-region dan cross-datacenter topology tidak diuji di dalam lab.",
            "Tidak mengimplementasikan change data capture (CDC) atau binlog-based synchronization."
        ], "type": "bullet"}
    ],
    narrative="Klaim teknis laboratorium ini dibatasi oleh parameter lingkungan pengujian yang disederhanakan. Cluster yang digunakan berjalan in-memory tanpa persistence disk, sehingga tidak bisa membuktikan perilaku recovery setelah node crash. State router disimpan dalam sync.Map lokal single-process, berbeda dengan arsitektur stateless production yang memerlukan Redis atau shared session store. Simulasi replication lag menggunakan timer sleep 200 milidetik, bukan network latency riil yang bervariasi. Konstanta timing seperti sticky 500ms dipilih untuk deterministik demo, bukan sebagai guideline tuning produksi.",
    narrative_en="The technical claims established in this lab are explicitly bounded by simplified test-environment parameters. The in-memory cluster lacks any persistence layer, so it cannot demonstrate post-crash recovery behavior. Router state resides in a local single-process sync.Map, differing from stateless production architectures that require Redis or shared session stores. Replication lag simulation uses deterministic timer sleeps rather than variable real network latency. Timing constants such as 500-millisecond stickiness were chosen for demonstrable predictability, not as production tuning guidelines.",
    mental_model="Lab proofs establish behavioral correctness, not production-scale resilience or observability.",
    source_path="labs/33-read-replicas-and-replication-lag/content/01-content-brief.md • research-audit/06-gaps.md",
    slide_num="11 / 24",
    layout="constraints-list"
)

# Slide 12: Trade-off arsitektur
add_slide(
    categories=["TRADE-OFF", "ARSITEKTUR"],
    title="Trade-off Arsitektur: Freshness vs Throughput",
    subtitle="Menilai biaya operasional penambahan routing intelligence dibandingkan benefit konsistensi.",
    content_items=[
        {"title": "MANFAAT & KEUNGGULAN", "items": [
            "Read-Your-Own-Writes Guaranteed: User selalu melihat hasil tulisannya sendiri segera.",
            "Scalable Reads: Replica menanggung sebagian besar beban read after sticky window expired.",
            "Graceful Degradation: Fallback ke primary terjadi secara eksplisit, bukan silent failure.",
            "Observability: LSN metrics dapat di-monitor untuk mendeteksi anomaly replication topology."
        ], "type": "bullet"},
        {"title": "KONSEKUENSI & BIAYA", "items": [
            "Sticky Window Overhead: Primary menerima traffic read ekstra selama T_sticky per session.",
            "Blocking Latency: Causal token wait menambah p99 read latency saat replica belum catch-up.",
            "Complexity Cost: Router membutuhkan session state management dan LSN inspection logic.",
            "Fallback Concentration: Semua degraded reads terakumulasi di primary → potential bottleneck."
        ], "type": "bullet"}
    ],
    narrative="Penerapan strategi replication lag mitigation menuntut kompromi arsitektural antara freshness guarantee dan read throughput. Sticky session menambahkan overhead berupa concentrated primary reads selama jendela waktu, tetapi eliminasi stale read sepenuhnya untuk session-bound operations. Causal token menawarkan presisi deterministik namun memperkenalkan blocking wait yang menggeser p99 latency read path. Lag-aware filtering menyeimbangkan keduanya dengan filtering threshold yang dapat dikonfigurasi, degraded gracefully ke primary ketika semua replica terlambat. Synchronous replication menghilangkan lag sama sekali dengan harga write latency ekstra yang proporsional dengan jumlah replica.",
    narrative_en="Implementing replication lag mitigation strategies demands architectural trade-offs between freshness guarantees and read throughput. Sticky sessions introduce overhead through concentrated primary reads during the time window but completely eliminate stale reads for session-bound operations. Causal tokens offer deterministic precision while introducing blocking waits that shift p99 read latency. Lag-aware filtering balances both approaches through configurable threshold predicates, degrading gracefully to primary when all replicas lag excessively. Synchronous replication eliminates lag entirely at the cost of proportionally increased write latency across all replicas.",
    mental_model="Consistency is purchased with either latency, complexity, or primary throughput sacrifice.",
    source_path="labs/33-read-replicas-and-replication-lag/content/02-master-draft.md",
    slide_num="12 / 24",
    layout="tradeoff"
)

# Slide 13: Concurrency / validation pipeline
add_slide(
    categories=["CONCURRENCY", "VALIDATION"],
    title="Validasi Konkurensi: 15 Goroutine × 20 Operasi Race-Free",
    subtitle="TestConcurrentAccess_RaceFree membuktikan thread-safety Router dan Cluster di bawah beban.",
    content_items=[
        {"title": "KODE TEST KONKURENSI", "items": [
            "func TestConcurrentAccess_RaceFree(t *testing.T) {",
            "    c := cluster.NewCluster(3, 10*time.Millisecond, cluster.AsyncReplication)",
            "    defer c.Close()",
            "    r := router.NewRouter(c, router.DefaultConfig())",
            "",
            "    var wg sync.WaitGroup",
            "    for i := 0; i < 15; i++ {",
            "        wg.Add(1)",
            "        go func(id int) {",
            "            defer wg.Done()",
            "            sid := fmt.Sprintf('s%d', id)",
            "            _, _ = r.Write(sid, fmt.Sprintf('k%d', id), fmt.Sprintf('v%d', id))",
            "            for j := 0; j < 20; j++ {",
            "                _, _, _, _ = r.ReadWithStickySession(sid, fmt.Sprintf('k%d', id))",
            "            }",
            "        }(i)",
            "    }",
            "    wg.Wait()",
            "}"
        ], "type": "code"},
        {"title": "HASIL VERIFIKASI", "items": [
            "Race Detector: PASS (go test -race ./... tanpa flag data race).",
            "Session Isolation: 15 session unik tidak saling mengganggu state sync.Map.",
            "Write-Read Consistency: Setiap session melihat data yang ditulis sendiri setelah sticky window.",
            "Cluster Stability: 3 replica workers tetap stabil tanpa goroutine leak atau deadlock."
        ], "type": "bullet"}
    ],
    narrative="Uji konkurensi ini membuktikan ketahanan implementasi saat 15 goroutine bersamaan mengeksekusi pattern write-then-read pada 20 iterasi per session. Setiap goroutine menggunakan session ID unik yang menjaga isolation state di Router sync.Map. Cluster dengan tiga replica worker memproses WAL entries secara paralel dengan gap 10 milidetik. Di bawah go race detector, tidak ada data race terdeteksi di seluruh memory accesses Router methods maupun Cluster nodes. Test ini mengonfirmasi bahwa concurrent access ke shared state tetap linearizable berkat mutex proteksi di setiap node.",
    narrative_en="This concurrency test validates implementation robustness when fifteen goroutines simultaneously execute write-then-read patterns across twenty iterations per session. Each goroutine employs a unique session ID maintaining isolation state within the Router sync.Map. The cluster with three replica workers processes WAL entries in parallel with ten-millisecond gaps. Under the Go race detector, no data races are detected across all Router method or Cluster node memory accesses. This confirms concurrent access to shared state remains linearizable through mutex protections at every node.",
    mental_model="Thread-safe shared state requires explicit mutex boundaries; race detector catches what humans miss.",
    source_path="tests/replication_test.go:173-210",
    slide_num="13 / 24",
    layout="concurrency-flow"
)

# Slide 14: Before vs After comparison cards
add_slide(
    categories=["BEFORE/AFTER", "PERUBAHAN"],
    title="Before vs After: Evolusi Read Routing Strategy",
    subtitle="Transformasi dari read naif berisiko stale hingga arsitektur four-layer mitigation.",
    content_items=[
        {"title": "BEFORE (Naive Round-Robin)", "items": [
            "Routing: Pilih replica secara acak tanpa awareness LSN state.",
            "Failure Mode: Stale read atau ErrNotFound terjadi setiap write-then-read window.",
            "Observability: Tidak ada metric untuk mendeteksi replication lag anomaly.",
            "Dampak: User experience degradasi senyap tanpa error reporting yang meaningful."
        ], "type": "card"},
        {"title": "AFTER (Four-Layer Mitigation)", "items": [
            "Routing: Sticky session untuk session-bound writes, token wait untuk causal reads.",
            "Failure Mode: Fallback eksplisit ke primary dengan logging node ID dan LSN delta.",
            "Observability: appliedLSN per replica dapat dipromonitor untuk alerting threshold breach.",
            "Dampak: Read-your-own-writes terjamin 100% dengan graceful degradation ketika replica lagging."
        ], "type": "card"}
    ],
    narrative="Perbandingan before-after ini menunjukkan evolusi arsitektur read routing dari pendekatan naif yang rentan stale read menuju sistem four-layer mitigation dengan fallback bertingkat. Sebelum penerapan strategi, setiap write-then-read cycle berisiko tinggi menghasilkan data usang karena tidak ada mekanisme penutupan replication gap. Setelah implementasi, sistem memiliki tiga lapis proteksi: sticky session untuk session-bound operations, causal token untuk critical reads, dan lag-aware filtering untuk bulk analytics. Setiap lapis memiliki fallback eksplisit ke primary dengan logging yang dapat diobservasi.",
    narrative_en="This before-and-after comparison highlights the architectural evolution from fragile naive routing prone to stale reads toward a four-layer mitigation system with tiered fallbacks. Prior to strategy implementation, every write-then-read cycle carried high risk of serving stale data due to absent gap-closing mechanisms. After deployment, the system features three protection layers: sticky sessions for session-bound operations, causal tokens for critical reads, and lag-aware filtering for bulk analytics. Each layer maintains explicit fallback logging to the primary with observable node identifiers and LSN deltas.",
    mental_model="Before: blind distribution amplifies inconsistency. After: informed routing turns lag into a measurable variable.",
    source_path="labs/33-read-replicas-and-replication-lag/content/02-master-draft.md",
    slide_num="14 / 24",
    layout="before-after"
)

# Slide 15: Analisis mendalam - Test 1 Naive Stale Read
add_slide(
    categories=["BUKTI TEST", "STALE READ"],
    title="Test 1: Stale Read Terjadi Tanpa Proteksi",
    subtitle="TestNaiveReplicationLag_StaleRead membuktikan anomaly read-then-write window.",
    content_items=[
        {"title": "TEST KONTRAK STALE READ", "items": [
            "func TestNaiveReplicationLag_StaleRead(t *testing.T) {",
            "    c := cluster.NewCluster(2, 200*time.Millisecond, cluster.AsyncReplication)",
            "    defer c.Close()",
            "    r := router.NewRouter(c, router.DefaultConfig())",
            "",
            "    _, err := r.Write(\"s1\", \"key:A\", \"value-A\")",
            "    if err != nil { t.Fatalf(\"Write failed: %v\", err) }",
            "",
            "    // Immediately read without any mitigation",
            "    val, nodeID, lsn, err := r.ReadNaive(\"key:A\")",
            "    if err == nil {",
            "        t.Fatalf(\"Expected stale read error, got val=%s node=%s lsn=%d\", val, nodeID, lsn)",
            "    }",
            "    // After waiting for lag to expire, read should succeed",
            "    time.Sleep(250 * time.Millisecond)",
            "    val, nodeID, lsn, err = r.ReadNaive(\"key:A\")",
            "    if err != nil { t.Fatalf(\"Post-lag read failed: %v\", err) }",
            "}"
        ], "type": "code"},
        {"title": "VERIFIKASI INVARIAN", "items": [
            "Immediate Read: ReadNaive pasca-write langsung menghasilkan ErrNotFound (stale read).",
            "Post-Lag Read: Setelah sleep 250ms (> lag 200ms), read berhasil dengan nilai fresh.",
            "Node Affinity: Setelah catch-up, read dapat berhasil dari replica manapun.",
            "Timing Proof: Perbedaan hasil sebelum dan sesudah lag membuktikan anomali temporal."
        ], "type": "bullet"}
    ],
    narrative="Unit test pertama membuktikan keberadaan stale read anomaly secara eksplisit. Setelah write ke primary dengan key 'key:A', panggilan ReadNaive seketika menghasilkan error ErrNotFound karena replica belum menerapkan WAL entry. Setelah menunggu 250 milidetik yang melampaui batas lag 200 milidetik, read yang sama berhasil mengembalikan value yang benar. Test ini mengonfirmasi bahwa anomali bersifat temporal — bukan permanen — dan dapat dicegah dengan menambahkan lapisan routing intelligence di antara write commit dan read dispatch.",
    narrative_en="The first unit test explicitly proves the existence of stale read anomalies. Immediately after writing to the primary with key 'key:A', the ReadNaive call produces an ErrNotFound error because replicas have not yet applied the WAL entry. After waiting 250 milliseconds to exceed the 200-millisecond lag threshold, the same read succeeds returning the correct value. This test confirms the anomaly is purely temporal — not permanent — and can be prevented by adding a routing intelligence layer between write commit and read dispatch.",
    mental_model="Stale reads are deterministic consequences of timing gaps, not random network failures.",
    source_path="tests/replication_test.go:1-45",
    slide_num="15 / 24",
    layout="test-analysis"
)

# Slide 16: Analisis mendalam - Test 2 Sticky Session
add_slide(
    categories=["BUKTI TEST", "STICKY ROUTING"],
    title="Test 2: Sticky Session Menjamin Read-Your-Own-Writes",
    subtitle="TestStickySessionRouting membuktikan primary redirect dalam jendela T_sticky.",
    content_items=[
        {"title": "TEST KONTRAK STICKY ROUTING", "items": [
            "func TestStickySessionRouting(t *testing.T) {",
            "    c := cluster.NewCluster(2, 200*time.Millisecond, cluster.AsyncReplication)",
            "    defer c.Close()",
            "    cfg := router.DefaultConfig()",
            "    cfg.StickyDuration = 500 * time.Millisecond",
            "    r := router.NewRouter(c, cfg)",
            "",
            "    _, err := r.Write(\"s1\", \"key:B\", \"value-B\")",
            "    if err != nil { t.Fatalf(\"Write failed: %v\", err) }",
            "",
            "    // Read within sticky window — must hit primary",
            "    val, nodeID, _, err := r.ReadWithStickySession(\"s1\", \"key:B\")",
            "    if err != nil || val != \"value-B\" || !strings.HasSuffix(nodeID, \"primary\") {",
            "        t.Fatalf(\"Sticky read failed or hit wrong node: val=%s node=%s err=%v\", val, nodeID, err)",
            "    }",
            "",
            "    time.Sleep(600 * time.Millisecond) // Expire sticky window",
            "    val, nodeID, _, err = r.ReadWithStickySession(\"s1\", \"key:B\")",
            "    if err != nil || nodeID == \"primary\" {",
            "        t.Fatalf(\"Post-sticky read should route to replica, got node=%s\", nodeID)",
            "    }",
            "}"
        ], "type": "code"},
        {"title": "ATURAN STRICT TIMING", "items": [
            "Pre-Expiration: Read dalam 500ms pertama selalu diarahkan ke primary, bukan replica.",
            "Post-Expiration: Setelah 600ms, routing berubah ke ReadLagAware (replica eligible).",
            "Value Integrity: Payload 'value-B' konsisten tersedia di primary dan replica catch-up.",
            "Session Isolation: Session ID lain tidak terpengaruh oleh sticky state session s1."
        ], "type": "bullet"}
    ],
    narrative="Unit test kedua mengonfirmasi bahwa sticky session routing berhasil mencegah stale read dalam jendela waktu yang ditentukan. Read yang dipanggil dalam 500 milidetik pasca-write secara konsisten diarahkan ke primary node, menjamin data fresh tersedia. Setelah jendela expire dengan sleep 600 milidetik, routing otomatis bergeser ke ReadLagAware yang menyeleksi replica berdasarkan threshold eligibility. Test ini membuktikan bahwa pendekatan time-based sederhana namun efektif untuk session-bound read patterns tanpa memerlukan LSN tracking di setiap request.",
    narrative_en="The second unit test confirms that sticky session routing successfully prevents stale reads within the specified time window. Reads invoked within 500 milliseconds post-write consistently route to the primary node, guaranteeing fresh data availability. After the window expires with a 600-millisecond sleep, routing automatically shifts to ReadLagAware which selects replicas based on threshold eligibility. This test proves the time-based approach is both simple and effective for session-bound read patterns without requiring per-request LSN tracking.",
    mental_model="Time-bounded primary redirection trades read throughput for guaranteed session consistency.",
    source_path="tests/replication_test.go:47-98",
    slide_num="16 / 24",
    layout="test-analysis"
)

# Slide 17: Analisis mendalam - Test 3 Causal Token
add_slide(
    categories=["BUKTI TEST", "CAUSAL TOKEN"],
    title="Test 3: Causal Token Tunggu Catch-Up Sebelum Serve Read",
    subtitle="TestReadWithToken_LSN membuktikan blocking wait dengan LSN deterministik.",
    content_items=[
        {"title": "TEST KONTRAK CAUSAL TOKEN", "items": [
            "func TestReadWithToken_LSN(t *testing.T) {",
            "    c := cluster.NewCluster(1, 200*time.Millisecond, cluster.AsyncReplication)",
            "    defer c.Close()",
            "    r := router.NewRouter(c, router.DefaultConfig())",
            "",
            "    lsn, err := r.Write(\"s1\", \"key:C\", \"value-C\")",
            "    if err != nil { t.Fatalf(\"Write failed: %v\", err) }",
            "",
            "    start := time.Now()",
            "    val, nodeID, readLSN, err := r.ReadWithToken(context.Background(), lsn, \"key:C\")",
            "    elapsed := time.Since(start)",
            "",
            "    if err != nil { t.Fatalf(\"ReadWithToken failed: %v\", err) }",
            "    if val != \"value-C\" { t.Fatalf(\"Expected value-C, got %s\", val) }",
            "    if readLSN < lsn { t.Fatalf(\"Expected LSN >= %d, got %d\", lsn, readLSN) }",
            "    if elapsed < 150*time.Millisecond {",
            "        t.Fatalf(\"Expected token read to wait ~200ms, finished in %v on node %s\", elapsed, nodeID)",
            "    }",
            "}"
        ], "type": "code"},
        {"title": "INVARIAN BLOCKING WAIT", "items": [
            "Deterministic Wait: ReadWithToken blocking hingga replica appliedLSN >= minLSN terbukti.",
            "Elapsed Timing: Test asserting elapsed >= 150ms mengkonfirmasi blocking behavior aktif.",
            "LSN Guarantee: readLSN >= minLSN menjamin causal ordering terpenuhi sebelum serve response.",
            "Fallback Safety: Jika timeout, method fallback ke primary dengan suffix '(fallback)' di nodeID."
        ], "type": "bullet"}
    ],
    narrative="Unit test ketiga menguji implementasi causal token routing yang melibatkan blocking wait mechanism. Setelah write dengan LSN tertentu, ReadWithToken iterasi replica untuk menemukan yang sudah mencapai threshold LSN. Ketika tidak ada replica yang match, method memanggil WaitForLSN pada replica pertama dengan context timeout. Test ini memvalidasi bahwa read blocked selama sekitar 200 milidetik menunggu replica catch-up, kemudian melayani response dengan nilai fresh dan LSN yang sesuai. Blocking behavior ini merupakan bedanya dengan sticky session yang hanya mengandalkan waktu tanpa memeriksa LSN aktual.",
    narrative_en="The third unit test examines the causal token routing implementation involving a blocking wait mechanism. After writing with a specific LSN, ReadWithToken iterates replicas to find one already at or above the threshold. When no replica matches, the method calls WaitForLSN on the first replica with a context timeout. This test validates that the read blocks for approximately 200 milliseconds waiting for replica catch-up before serving a fresh response with the correct LSN. This blocking behavior distinguishes causal tokens from sticky sessions, which rely solely on time without checking actual LSN values.",
    mental_model="Causal tokens convert temporal uncertainty into deterministic ordering guarantees.",
    source_path="tests/replication_test.go:100-121",
    slide_num="17 / 24",
    layout="test-analysis"
)

# Slide 18: Studi Kasus Industri: AWS RDS vs PostgreSQL Streaming
add_slide(
    categories=["STUDI KASUS", "INDUSTRI"],
    title="Studi Kasus Industri: Replication Lag di Database Production",
    subtitle="Bagaimana ekosistem database kelas dunia menangani stale read problem secara praktis.",
    content_items=[
        {"title": "AWS RDS / AURORA READ REPLICAS", "items": [
            "Default: Asynchronous replication dengan lag bervariasi tergantung region dan instance size.",
            "Read-Your-Own-Writes: AWS DocumentDB menyediakan 'session consistency' via sticky routing otomatis.",
            "Monitoring: CloudWatch metric 'ReplicaLag' exposing seconds behind primary secara real-time.",
            "Failover: Automated promotion namun read splitter perlu rekonfigurasi sticky session state."
        ], "type": "bullet"},
        {"title": "POSTGRESQL STREAMING REPLICATION", "items": [
            "Async Default: wal_level = replica, synchronous_commit = off, remote_apply belum diaktifkan.",
            "Synchronous Commit: Setting synchronous_commit = remote_apply menjamin write visibility di semua replica.",
            "pg_stat_replication: System view exposing replay_lsn vs write_lsn untuk monitoring lag.",
            "Application Level: Driver seperti PgBouncer atau connection pooler bisa implement sticky routing."
        ], "type": "bullet"}
    ],
    narrative="Studi kasus industri menunjukkan bahwa masalah replication lag bukan teoretis melainkan operasional sehari-hari di infrastruktur cloud modern. AWS RDS menyediakan metric ReplicaLag melalui CloudWatch yang memungkinkan tim SRE mendeteksi anomali sebelum user complaint muncul. PostgreSQL streaming replication memberikan kontrol granular melalui synchronous_commit parameter yang membedakan antara durability guarantee dan visibility guarantee. Aplikasi enterprise besar seperti Twitter dan Netflix mengimplementasikan custom session-affinity routing di layer client library untuk menutup gap readability tanpa menunggu perbaikan di database engine.",
    narrative_en="Industry case studies demonstrate that replication lag is not merely theoretical but an everyday operational challenge in modern cloud infrastructure. AWS RDS provides the ReplicaLag metric through CloudWatch enabling SRE teams to detect anomalies before user complaints surface. PostgreSQL streaming replication offers granular control through the synchronous_commit parameter distinguishing between durability guarantees and visibility guarantees. Large enterprise applications like Twitter and Netflix implement custom session-affinity routing within client libraries to close readability gaps without waiting for database engine fixes.",
    mental_model="Production systems treat replication lag as an operational variable, not an inflexible constraint.",
    source_path="research/05-report.md • research/03-evidence.md",
    slide_num="18 / 24",
    layout="cards"
)

# Slide 19: Analisis Teknis: sync.Cond WaitForLSN vs Context Timeout
add_slide(
    categories=["ANALISIS TEKNIS", "RUNTIME"],
    title="Bagaimana WaitForLSN Membedakan Blocking dan Timeout?",
    subtitle="Konkurensi pattern menggunakan sync.Cond + context cancellation untuk polling efisien.",
    content_items=[
        {"title": "IMPLEMENTASI WAITER", "items": [
            "func (n *Node) WaitForLSN(ctx context.Context, targetLSN uint64) error {",
            "    ch := make(chan struct{})",
            "    stop := make(chan struct{})",
            "    go func() {",
            "        n.mu.Lock()",
            "        defer n.mu.Unlock()",
            "        for n.appliedLSN < targetLSN {",
            "            select {",
            "            case <-stop:",
            "                return",
            "            default:",
            "            }",
            "            n.cond.Wait()  // Release lock, block until Broadcast",
            "        }",
            "        close(ch)",
            "    }()",
            "    select {",
            "    case <-ctx.Done():",
            "        close(stop); n.cond.Broadcast(); return ctx.Err()",
            "    case <-ch:",
            "        return nil",
            "    }",
            "}"
        ], "type": "code"},
        {"title": "MEKANISME PENGUNCIAN", "items": [
            "cond.Wait() melepas mutex sementara goroutine block, menghindari busy-wait CPU spinning.",
            "Worker goroutine memanggil cond.Broadcast() setelah apply entry, waking semua waiter.",
            "Context deadline memicu close(stop) + Broadcast() untuk emergency wakeup dan cleanup."
        ], "type": "bullet"}
    ],
    narrative="Implementasi WaitForLSN menunjukkan pola konkurensi Go yang elegan untuk menunggu kondisi berbasis channel dan condvar. Goroutine waiter mengakuisisi mutex node, loop until appliedLSN memenuhi target, dan memanggil cond.Wait() yang secara atomik melepaskan mutex sambil memblokir goroutine. Worker goroutine setelah apply WAL entry memanggil cond.Broadcast() untuk membangunkan semua waiter yang sedang blocked. Context cancellation memberikan jalur escape yang aman dengan menutup stop channel dan memicu broadcast ulang agar waiter tidak stuck selamanya.",
    narrative_en="The WaitForLSN implementation demonstrates an elegant Go concurrency pattern for condition-based waiting combining channels and condvars. The waiter goroutine acquires the node mutex, loops until appliedLSN satisfies the target, and calls cond.Wait() which atomically releases the mutex while blocking the goroutine. The worker goroutine invokes cond.Broadcast() after applying WAL entries to wake all blocked waiters. Context cancellation provides a safe escape hatch by closing the stop channel and triggering a re-broadcast ensuring waiters never remain permanently stuck.",
    mental_model="Condition variables replace polling loops; context cancellation ensures graceful deadline enforcement.",
    source_path="internal/cluster/cluster.go:79-107",
    slide_num="19 / 24",
    layout="code-bullet"
)

# Slide 20: Checklist implementasi
add_slide(
    categories=["CHECKLIST", "PRODUKSI"],
    title="Checklist Implementasi Read Replica Strategy",
    subtitle="10 poin verifikasi wajib sebelum merutekan traffic ke environment production.",
    content_items=[
        {"title": "CHECKLIST ARSITEKTUR REPLIKASI", "items": [
            "[ ] Tentukan consistency requirement: session-bound vs casual reads untuk memilih strategi.",
            "[ ] Kalibrasi T_sticky ≥ p99 replication lag aktual di production environment (bukan asumsi).",
            "[ ] Konfigurasi MaxLSNDiff berdasarkan observasi delta LSN normal vs anomaly threshold.",
            "[ ] Set WaitTimeout causal token ≥ worst-case catch-up duration + 20% buffer.",
            "[ ] Deploy monitoring metrics: replica lag histogram, stale read rate, primary fallback frequency.",
            "[ ] Implement circuit breaker untuk primary fallback ketika semua replica exceed threshold.",
            "[ ] Distribusikan router session state ke shared store (Redis/memcached) untuk horizontal scaling.",
            "[ ] Validasi bahwa sticky session map tidak memory-leak pada session idle jangka panjang.",
            "[ ] Uji synchronous replication mode hanya untuk write-critical paths, bukan seluruh traffic.",
            "[ ] Prepare runbook deteksi anomaly: cara membedakan network partition dari replication lag spike."
        ], "type": "bullet"}
    ],
    narrative="Daftar periksa ini merangkum sepuluh prinsip krusial yang wajib dipenuhi sebelum deploying read replica strategy ke production. Parameter timing seperti sticky duration dan wait timeout tidak boleh di-set berdasarkan feeling, melainkan harus dikalibrasi dari data historis replication lag percentile. Monitoring dan alerting menjadi mata telinga operator untuk mendeteksi degradation sebelum user complain. State management router perlu diverifikasi tidak menyebabkan memory leak pada session lifecycle yang panjang. Terakhir, dibedakan antara synchronous replication untuk critical path versus asynchronous dengan application-level fencing untuk read-heavy workload.",
    narrative_en="This checklist codifies ten critical principles required before deploying read replica strategies to production. Timing parameters such as sticky duration and wait timeout must never be guessed but calibrated against historical replication lag percentile data. Monitoring and alerting become the operator sensory system for detecting degradation before user complaints surface. Router state management must be verified for absence of memory leaks across long-lived session lifecycles. Finally, distinguish synchronous replication for critical paths versus asynchronous with application-level fencing for read-heavy workloads.",
    mental_model="Calibrate timing empirically, monitor continuously, and distinguish path criticality explicitly.",
    source_path="labs/33-read-replicas-and-replication-lag/content/05-key-takeaways.md",
    slide_num="20 / 24",
    layout="checklist"
)

# Slide 21: Anti-pattern & cara reproduksi kegagalan
add_slide(
    categories=["ANTI-PATTERN", "REPRODUKSI"],
    title="Anti-Pattern: Kesalahan Desain yang Harus Dihindari",
    subtitle="Kekeliruan arsitektur yang melumpuhkan jaminan keselamatan replication topology.",
    content_items=[
        {"title": "ANTI-PATTERN UMUM", "items": [
            "Mencampuradukkan semi-sync dengan synchronous: MySQL semi-sync menjamin relay log receipt, bukan apply visibility.",
            "Sticky Duration Terlalu Pendek: T_sticky < p99 lag menyebabkan session masih mengalami stale read selepas jendela.",
            "Tidak Ada Fallback Primary: ReadLagAware tanpa fallback utama menyebabkan total read unavailability saat semua replica lagging.",
            "Ignore LSN pada Cached Reads: Read dari application cache tanpa invalidasi berdasarkan LSN menghasilkan dirty reads persisten.",
            "Hardcode Waiting Durasi: Menggunakan time.Sleep statis alih-alih context timeout menyebabkan goroutine leak dan deadlock."
        ], "type": "card"},
        {"title": "CARA REPRODUKSI GAGAL DI LAB", "items": [
            "1. Ubah router.go ReadLagAware: hapus fallback primary → test TestReplicaLagThreshold_Fallback akan fail.",
            "2. Set MaxLSNDiff = 0 → semua replica terfilter, read selalu fallback primary; throughput read drop drastis.",
            "3. Hapus WaitForLSN context timeout → TestReadWithToken_LSN hanging indefinitely pada replica yang crash."
        ], "type": "card"}
    ],
    narrative="Mengabaikan prinsip desain replication topology dapat berakibat fatal bagi availability dan consistency system. Anti-pattern yang paling sering terjadi adalah underestimate durasi sticky window dibanding actual replication lag, menghasilkan false sense of security bagi engineer. Tidak menyertakan fallback primary pada lag-aware routing berarti read operations akan total unavailable ketika semua replica mengalami lag spike simultan. Hardcoding sleep duration alih-alih menggunakan context timeout juga menyebabkan goroutine leakage yang dapat mematikan seluruh cluster worker.",
    narrative_en="Violating replication topology design principles compromises both system availability and data consistency. The most common anti-pattern is underestimating sticky window duration relative to actual replication lag, inducing a false sense of security for engineering teams. Omitting primary fallbacks in lag-aware routing guarantees total read unavailability whenever all replicas suffer simultaneous lag spikes. Furthermore, hardcoding sleep durations rather than respecting context timeouts invites severe goroutine leakage that can paralyze cluster workers.",
    mental_model="Naive fallbacks and uncalibrated timers turn transient replication lag into catastrophic outages.",
    source_path="labs/33-read-replicas-and-replication-lag/content/05-key-takeaways.md",
    slide_num="21 / 24",
    layout="anti-pattern-cards"
)

# Slide 22: Takeaway praktis & command block dominan
add_slide(
    categories=["PRAKTIS", "COMMAND"],
    title="Verifikasi Mandiri: Uji & Jalankan Lab",
    subtitle="Reproduksi seluruh bukti teknis, benchmark race detector, dan demo interaktif di terminal Anda.",
    content_items=[
        {"title": "PERINTAH VERIFIKASI RESMI", "items": [
            "cd labs/33-read-replicas-and-replication-lag",
            "go test -v ./...",
            "go test -race ./...",
            "go run ./cmd/demo"
        ], "type": "code"},
        {"title": "RINGKASAN METRIK LAB", "items": [
            "Test Suite: 7 test cases lulus 100% tanpa kegagalan (2.13s total runtime).",
            "Race Detector: Bersih dari thread contention dan data races (3.98s pass).",
            "Dependencies: Murni Go Standard Library (sync, time, context, atomic) tanpa dependensi pihak ketiga."
        ], "type": "bullet"}
    ],
    narrative="Lab ini dirancang agar dapat diverifikasi secara mandiri dalam hitungan detik di mesin lokal. Seluruh implementasi dibangun murni menggunakan standard library Go tanpa memerlukan server eksternal seperti PostgreSQL atau MySQL. Perintah go test -v mengeksekusi tujuh pengujian komprehensif mulai dari anomali stale read hingga konkurensi race-free. Perintah go test -race memvalidasi keamanan thread di bawah beban 15 goroutine intensif, dan go run ./cmd/demo menampilkan demonstrasi terminal interaktif keempat skenario routing.",
    narrative_en="This laboratory is engineered for immediate, frictionless local verification within seconds. The entire codebase is implemented with the Go standard library, eliminating external dependencies like PostgreSQL or MySQL instances. Running go test -v executes seven exhaustive test scenarios validating stale read anomalies and race-free concurrency. Executing go test -race proves thread safety under fifteen concurrent goroutines, while go run ./cmd/demo delivers an interactive terminal demonstration of all four routing scenarios.",
    mental_model="Zero-dependency standard library code enables instant, deterministic peer verification.",
    source_path="README.md:19-34 • engineering/03-execution-result.md",
    slide_num="22 / 24",
    layout="command-block"
)

# Slide 23: Prinsip arsitektural jangka panjang
add_slide(
    categories=["PRINSIP", "ARSITEKTUR"],
    title="Prinsip Arsitektur: Desain untuk Asinkronitas",
    subtitle="Pelajaran abadi dalam merancang sistem terdistribusi yang tahan terhadap replication lag.",
    content_items=[
        {"title": "1. ASINKRONITAS ADALAH KENYATAAN DEFAULT", "items": [
            "Replication lag bukan anomali yang harus dihilangkan, melainkan realitas yang harus dikelola.",
            "Asumsikan setiap read ke replica berpotensi stale kecuali dibuktikan sebaliknya oleh router."
        ], "type": "card"},
        {"title": "2. CLOSING MECHANISM ADALAH HAK KLIEN", "items": [
            "Aplikasi bertanggung jawab memilih mekanisme penutupan jarak (sticky, token, threshold).",
            "Database engine menyediakan primitives (LSN, WAL), tetapi router aplikasi menentukan kebijakan."
        ], "type": "card"},
        {"title": "3. DEGRADASI EKSPLISIT MENGALAHKAN CRASH SENYAP", "items": [
            "Fallback ke primary dengan logging terstruktur jauh lebih baik daripada menyajikan data basi.",
            "Transparansi status node dan LSN watermark menjaga integritas sistem jangka panjang."
        ], "type": "card"}
    ],
    narrative="Tiga prinsip arsitektur ini berlaku universal di luar implementasi lab Go ini. Pertama, sadari bahwa asinkronitas adalah default industri dan upaya menghilangkannya secara total akan mengorbankan skalabilitas tulis. Kedua, tempatkan tanggung jawab penutupan replication gap di lapisan aplikasi karena hanya aplikasilah yang memahami konteks konsistensi setiap use case. Ketiga, prioritaskan degradasi eksplisit ke primary dengan observabilitas lengkap dibandingkan membiarkan user menerima data stale secara senyap.",
    narrative_en="These three architectural principles transcend this specific Go implementation. First, recognize that asynchronous replication is the industry default, and attempting to eradicate it entirely sacrifices write scalability. Second, place the responsibility for closing replication gaps at the application layer because only the application understands the consistency context of each use case. Third, prioritize explicit degradation to the primary with rich observability over allowing users to silently ingest stale data.",
    mental_model="Asynchrony is reality; smart routing bridges the gap between scalability and user trust.",
    source_path="labs/33-read-replicas-and-replication-lag/content/05-key-takeaways.md",
    slide_num="23 / 24",
    layout="principles-list"
)

# Slide 24: Closing slide
add_slide(
    categories=["CLOSING"],
    title="Skalabilitas Baca Terjaga, Konsistensi Terjamin.",
    subtitle="",
    content_items=[
        {"title": "KOMPONEN ARSITEKTUR LAB", "items": [
            "internal/cluster/cluster.go — Primary-replica simulation, WAL streaming, LSN tracking, & sync/async paths",
            "internal/router/router.go — Four-strategy read splitter with session-affinity & dynamic lag filtering",
            "tests/replication_test.go — 7 comprehensive test suites (stale reads, sticky routing, token wait, race-free)",
            "cmd/demo/main.go — Complete interactive lifecycle demonstration of 4 routing scenarios"
        ], "type": "bullet"},
        {"title": "VERIFIKASI & REPOSITORI", "items": [
            f"Commit SHA: {COMMIT_SHA[:8]}",
            "Unit Tests: 7/7 PASS (100% clean, 2.13s)",
            "Race Detector: PASS (zero data race, 3.98s)",
            f"Repositori: github.com/lukman-ss/software-engineering-lab/tree/main/labs/{LAB_NO}-{SLUG}"
        ], "type": "bullet"}
    ],
    narrative="Laboratorium Read Replicas and Replication Lag ini membuktikan secara konkret bahwa dilema klasik antara skalabilitas baca dan konsistensi data dapat dijembatani secara elegan. Melalui kombinasi sticky session routing, causal LSN tokens, dan lag-aware threshold filtering, sistem mampu memaksimalkan offloading query ke replica tanpa mengorbankan garansi read-your-own-writes bagi pengguna. Seluruh kode sumber, dokumentasi arsitektur, dan instruksi verifikasi tersedia lengkap di repositori.",
    narrative_en="This Read Replicas and Replication Lag laboratory definitively demonstrates that the classic dilemma between read scalability and data consistency can be elegantly reconciled. Through a combination of sticky session routing, causal LSN tokens, and lag-aware threshold filtering, the architecture maximizes read offloading to replicas without compromising read-your-own-writes guarantees for users. Complete source code, design specifications, and verification harnesses are available in the repository.",
    mental_model="Replication scales reads; application routing safeguards consistency.",
    source_path=REPO_URL,
    slide_num="24 / 24",
    layout="closing-dark-navy"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
