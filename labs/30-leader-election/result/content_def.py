#!/usr/bin/env python3
"""
Publication Builder for Software Engineering Lab 30: Leader Election with Fencing Tokens.
Generates 24-slide LinkedIn carousel covering lease-based leader election, failure detection via TTL,
transient dual-leader states, monotonic fencing tokens, and split-brain defense at the shared storage layer.
"""
import os

OUTPUT_DIR = "labs/30-leader-election/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/30-leader-election"
COMMIT_SHA = "ba4b6216da6a06dee61959231d124cd44085bbf5"
LAB_NO = "30"
SLUG = "leader-election"
TITLE_ID = "Leader Election dengan Fencing Tokens"
SUBTITLE_ID = "Mencegah Split-Brain dan State Corruption Menggunakan Lease TTL dan Monotonic Revision"
TITLE_EN = "Leader Election with Fencing Tokens"
SUBTITLE_EN = "Preventing Split-Brain and State Corruption via Lease TTL and Monotonic Revision"
AUTHOR = "Lukman / lukman-ss"
COVER_IMAGE_PATH = "labs/30-leader-election/result/image-cover/cover_art.jpg"

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
            "Retro CRT monitor dengan scanlines",
            "Diagram leader election: Lease TTL, pause window, and monotonic fencing check",
            f"Penulis: {AUTHOR}",
        ]}
    ],
    mental_model="Leases grant temporal authority; monotonic fencing tokens enforce permanent ordering at storage.",
    source_path=REPO_URL,
    slide_num="01 / 24",
    narrative="Leader election adalah mekanisme dasar dalam sistem terdistribusi untuk memastikan hanya satu node yang mengeksekusi operasi kritis pada satu waktu. Namun, lease berbasis waktu saja tidak cukup mencegah korupsi data ketika leader mengalami jeda eksekusi tak terduga. Lab ini membuktikan implementasi murni Go untuk lease coordinator berbatas TTL dipadukan dengan fencing token monotonik naik. Dengan memvalidasi token secara atomik di lapisan resource bersama, sistem kebal terhadap anomali split-brain dan penulisan usang.",
    narrative_en="Leader election is a fundamental mechanism in distributed systems ensuring only one node coordinates critical operations at any given moment. However, time-bounded leases alone fail to prevent data corruption when a leader experiences unannounced execution pauses. This lab demonstrates a pure Go implementation of a TTL-bounded lease coordinator paired with monotonically increasing fencing tokens. By enforcing strict token checks atomically at the shared resource layer, the architecture remains fully resilient against split-brain writes.",
    layout="cover"
)

# Slide 2: Problem nyata
add_slide(
    categories=["MOTIVASI", "SPLIT-BRAIN RISK"],
    title="Masalah Nyata: Ilusi Otoritas Tunggal",
    subtitle="Ketika node mengalami stop-the-world GC pause, koordinasi terdistribusi rentan terbelah.",
    content_items=[
        {"title": "UNANNOUNCED PAUSE", "items": [
            "Garbage collector pause atau network hiccup menghentikan thread leader.",
            "Leader tidak dapat mengirim heartbeat renewal sebelum lease TTL habis.",
            "Node tidak menyadari bahwa ia sedang membeku dan kehilangan otoritas."
        ], "type": "card"},
        {"title": "PREMATURE FAILOVER", "items": [
            "Coordinator menganggap leader lama mati karena TTL kedaluwarsa.",
            "Kandidat berikutnya memenangkan election dan dipromosikan jadi leader baru.",
            "Dua node kini sama-sama menganggap dirinya pemegang mandat sah."
        ], "type": "card"},
        {"title": "STALE WRITE HAZARD", "items": [
            "Leader lama terbangun tanpa menyadari posisinya telah digantikan.",
            "Ia mengirim mutasi usang ke storage bersama menimpa data leader baru.",
            "Korupsi data senyap terjadi tanpa mekanisme validasi di storage."
        ], "type": "card"}
    ],
    narrative="Dalam arsitektur terdistribusi nyata, proses tidak pernah mati secara instan dan rapi. Fenomena runtime seperti stop-the-world garbage collection pause pada runtime modern atau network partition parsial dapat membekukan proses selama ratusan milidetik. Selama periode pause tersebut, lease TTL milik leader lama akan kedaluwarsa di coordinator, memicu failover otomatis ke kandidat baru. Ketika leader lama bangun dari pause, ia tidak menyadari waktu telah berlalu dan tetap mengeksekusi operasi tulis yang berpotensi merusak integritas data sistem.",
    narrative_en="In production distributed systems, processes rarely fail cleanly or predictably. Runtime anomalies such as stop-the-world garbage collection pauses or transient network partitions can freeze a node for hundreds of milliseconds. During this pause window, the leader's lease TTL expires on the coordinator, triggering an automatic failover to promote a successor. When the old leader resumes execution, it remains completely unaware that time elapsed and continues writing, creating silent data corruption unless strictly intercepted.",
    mental_model="Time passes invisibly to a paused process; local belief of leadership is never sufficient proof.",
    source_path="labs/30-leader-election/content/01-content-brief.md",
    slide_num="02 / 24",
    layout="cards"
)

# Slide 3: Mengapa review manual / pendekatan naif gagal
add_slide(
    categories=["ANALISIS", "PENDEKATAN NAIF"],
    title="Mengapa Lease Berbasis Waktu Saja Gagal?",
    subtitle="Pemeriksaan status lokal di sisi klien tidak memiliki jaminan linieritas saat eksekusi terdistribusi.",
    content_items=[
        {"title": "KODE MUTASI NAIF", "items": [
            "func (n *Node) DoWorkNaive(val string) error {",
            "    if n.state != StateLeader { return ErrNotLeader }",
            "    // GC pause terjadi tepat di sini selama 350ms!",
            "    // Lease TTL 200ms kedaluwarsa di coordinator.",
            "    // Node-B terpilih sebagai leader baru.",
            "    return n.sharedStorage.Save(val) // Menimpa data Node-B!",
            "}"
        ], "type": "code"},
        {"title": "KEGAGALAN FUNDAMENTAL", "items": [
            "Gap waktu antara pengecekan status (check) dan mutasi storage (act).",
            "Asinkronitas jaringan dan runtime membuat pengecekan waktu lokal tidak dapat dipercaya.",
            "Storage pasif menerima semua write tanpa memverifikasi urutan generasi pengirim."
        ], "type": "bullet"}
    ],
    narrative="Pendekatan naif mengasumsikan pengecekan status boolean lokal `if n.state == StateLeader` cukup menjamin mutual exclusion. Asumsi ini salah fatal dalam sistem asinkron karena adanya time-of-check to time-of-use (TOCTOU) gap. Tepat setelah pengecekan lokal berhasil, runtime thread dapat diinterupsi oleh GC pause atau penjadwalan OS selama ratusan milidetik. Saat thread kembali aktif, status kepemimpinan telah berpindah ke node lain, namun instruksi tulis tetap dikirim ke storage. Tanpa validasi di lapisan penyimpanan, write usang akan menimpa data terkini.",
    narrative_en="The naive approach assumes checking a local boolean flag `if n.state == StateLeader` guarantees mutual exclusion. This assumption fails catastrophically due to the time-of-check to time-of-use (TOCTOU) vulnerability in asynchronous distributed environments. Immediately after the local check evaluates to true, the thread can be suspended by a GC pause or OS preemption for hundreds of milliseconds. When the thread resumes, leadership has already transitioned to another candidate, yet the write proceeds, corrupting storage unless validated at the boundary.",
    mental_model="Local validation suffers from TOCTOU gaps; true mutual exclusion must be verified by the shared target.",
    source_path="labs/30-leader-election/research/03-mechanisms-and-leases.md",
    slide_num="03 / 24",
    layout="code-bullet"
)

# Slide 4: Mental model inti / schema card
add_slide(
    categories=["MENTAL MODEL", "SCHEMA"],
    title="Mental Model: Otoritas Temporal vs Validasi Monotonik",
    subtitle="Lease mengatur durasi giliran kerja, fencing token menegakkan urutan generations di storage.",
    content_items=[
        {"title": "STRUKTUR DATA KOORDINASI", "items": [
            "type Lease struct {",
            "    Key          string",
            "    HolderID     string",
            "    FencingToken int64         // revision monotonik naik",
            "    ExpiresAt    time.Time",
            "    TTL          time.Duration",
            "}"
        ], "type": "code"},
        {"title": "PRINSIP FENCING TOKEN", "items": [
            "Strict Monotonicity: Counter revision coordinator selalu bertambah di setiap kepemimpinan baru.",
            "Resource Gatekeeping: Storage melacak lastSeenToken dan menolak token <= lastSeenToken.",
            "Temporal Decoupling: Integritas data tidak bergantung pada sinkronisasi wall-clock mesin klien.",
            "Self-Defending Store: Sistem kebal meskipun terjadi transient dual-leader state di memori."
        ], "type": "bullet"}
    ],
    narrative="Mental model arsitektur ini memisahkan dua konsep: hak memimpin sementara (lease) dan hak menulis permanen (fencing token). Lease coordinator mengelola kepemimpinan dengan TTL berbatas waktu, menerbitkan nomor tiket revisi yang strictly monotonically increasing setiap kali pemimpin baru dilantik. Penyimpanan bersama tidak lagi mempercayai identitas klien secara buta, melainkan bertindak sebagai gerbang independen yang hanya menerima mutasi jika tiket yang dibawa strictly lebih besar dari tiket terakhir yang tercatat.",
    narrative_en="This architectural mental model decouples two distinct concepts: transient leadership authority (lease) and permanent write eligibility (fencing token). The coordinator administers time-bounded leases, issuing a strictly monotonically increasing revision ticket whenever a new leader is recognized. The shared storage layer no longer blindly trusts client identity; instead, it acts as an autonomous gatekeeper that accepts mutations if and only if the presented token strictly exceeds the last recorded token.",
    mental_model="Leases decide who should try to write; fencing tokens decide whose write is allowed to persist.",
    source_path="internal/coordinator/coordinator.go:15-28 • internal/storage/storage.go:14-25",
    slide_num="04 / 24",
    layout="schema-card"
)

# Slide 5: Diagram alur/pipeline
add_slide(
    categories=["ARSITEKTUR", "PIPELINE"],
    title="Alur Koordinasi & Penolakan Mutasi Basi",
    subtitle="Pipeline端到端: Campaign → Lease Grant → Normal Write → Pause/Failover → Stale Write Rejection.",
    content_items=[
        {"title": "SEKUENSI PIPELINE DISTRIBUSI", "items": [
            "1. Node-A.Acquire() → Coordinator terbitkan Lease(Holder=A, Token=1, TTL=200ms)",
            "2. Node-A.PerformFencedWrite(val) → Storage terima: incoming 1 > lastSeen 0 → lastSeen=1",
            "3. Node-A mengalami Simulasi Pause 350ms (eksekusi berhenti, heartbeat terhenti)",
            "4. Waktu berjalan melewati 200ms → Lease Node-A kedaluwarsa di Coordinator",
            "5. Node-B.Acquire() → Coordinator terbitkan Lease(Holder=B, Token=2, TTL=200ms)",
            "6. Node-B.PerformFencedWrite(val) → Storage terima: incoming 2 > lastSeen 1 → lastSeen=2",
            "7. Node-A bangun dari pause → Mencoba Write dengan Token=1 usang",
            "8. Storage Gatekeeper menolak: token 1 <= lastSeen 2 → ErrStaleFencingToken!"
        ], "type": "diagram"}
    ],
    narrative="Pipeline koordinasi menunjukkan bagaimana kegagalan temporal ditangani secara deterministik. Tahap 1 hingga 2 memperlihatkan alur kerja normal di mana Node-A memenangkan election dan menulis sukses dengan token generasi 1. Pada tahap 3 hingga 5, Node-A mengalami pause melebihi batas TTL, memungkinkan coordinator melantik Node-B dengan token generasi 2 yang berhasil menulis ke storage. Saat Node-A tersadar di tahap 7 dan mencoba menulis dengan token 1, storage mengeksekusi fencing check dan menolak mutasi tersebut secara mutlak.",
    narrative_en="The coordination pipeline demonstrates how temporal anomalies are handled deterministically. Stages 1 through 2 capture normal execution where Node-A wins election and writes successfully under generation token 1. During stages 3 to 5, Node-A pauses beyond its TTL window, enabling the coordinator to expire the lease and promote Node-B with token 2. When Node-A wakes up in stage 7 and attempts to commit with stale token 1, the storage layer evaluates the fencing invariant and rejects the write.",
    mental_model="Coordinators manage failover transitions; storage gatekeepers enforce immutable monotonic ordering.",
    source_path="cmd/demo/main.go:15-77",
    slide_num="05 / 24",
    layout="pipeline-diagram"
)

# Slide 6: Cuplikan kode dominan
add_slide(
    categories=["IMPLEMENTASI", "CODE"],
    title="Implementasi Gerbang Fencing di Shared Storage",
    subtitle="Pemeriksaan atomik token <= lastSeenToken di bawah perlindungan mutex pada storage.go.",
    content_items=[
        {"title": "FENCED STORAGE WRITE GATE", "items": [
            "func (s *FencedStorage) Write(author string, token int64, value string) error {",
            "    s.mu.Lock()",
            "    defer s.mu.Unlock()",
            "",
            "    if token <= s.lastSeenToken {",
            "        return fmt.Errorf(\"%w: token %d <= last seen %d (author: %s)\",",
            "            ErrStaleFencingToken, token, s.lastSeenToken, author)",
            "    }",
            "    s.lastSeenToken = token",
            "    s.records = append(s.records, Record{FencingToken: token, Author: author, Value: value, Timestamp: time.Now()})",
            "    return nil",
            "}"
        ], "type": "code"}
    ],
    narrative="Potongan kode ini merupakan inti pertahanan terhadap split-brain di seluruh arsitektur sistem. Method Write pada struct FencedStorage mengunci mutex untuk menjamin operasi check-and-set berlangsung secara linier dan atomik. Aturan validasi mengevaluasi apakah token yang disertakan bernilai lebih kecil atau sama dengan lastSeenToken. Jika ya, storage langsung menolak operasi dengan mengembalikan ErrStaleFencingToken. Hanya token yang strictly lebih besar yang diizinkan memperbarui lastSeenToken dan mencatat data mutasi.",
    narrative_en="This code block constitutes the primary defense against split-brain corruption across the architecture. The Write method on FencedStorage acquires an exclusive mutex lock to ensure linearizable check-and-set semantics. The validation logic verifies whether the incoming token is less than or equal to lastSeenToken. If so, storage immediately aborts the write with ErrStaleFencingToken. Only strictly greater tokens update the internal watermark and append mutation records.",
    mental_model="Strict inequality incoming > lastSeenToken turns disordered distributed writes into a monotonic stream.",
    source_path="internal/storage/storage.go:33-49",
    slide_num="06 / 24",
    layout="code-dominant"
)

# Slide 7: Failure path / dangling reference konkret
add_slide(
    categories=["FAILURE PATH", "BUKTI"],
    title="Failure Path: Penolakan Stale Write Pascajeda GC",
    subtitle="Log terminal membuktikan Node-A ditolak setelah Node-B mengambil alih kepemimpinan.",
    content_items=[
        {"title": "TERMINAL FAILURE LOG", "items": [
            "[4] Simulating Stop-the-World GC Pause / Network Partition on Node-A (350ms > TTL 200ms)...",
            "[5] Status after TTL expiration: Node-A=LEADER, Node-B=LEADER",
            "[6] New Leader promoted: Node-B (Fencing Token: 2)",
            "[7] New Leader Node-B performing fenced write...",
            "Write SUCCESS by Node-B (Token 2)",
            "[8] Old Leader Node-A wakes up from GC pause (still holding stale token 1)...",
            "[9] Old Leader Node-A attempts write to shared storage using stale token 1...",
            "Write REJECTED by FencedStorage: stale fencing token: rejected to prevent split-brain write: token 1 <= last seen 2 (author: Node-A)"
        ], "type": "terminal"},
        {"title": "ANALISIS KEGAGALAN", "items": [
            "Pada baris [5], terjadi transient dual-leader: Node-A dan Node-B sama-sama berstatus LEADER.",
            "Meskipun Node-A merasa sah, storage berhasil mendeteksi anomali generasi (token 1 <= 2).",
            "Shared storage history tetap bersih: hanya mencatat record dari generasi kepemimpinan yang valid."
        ], "type": "bullet"}
    ],
    narrative="Failure path konkret ini diambil langsung dari hasil eksekusi program demo lab. Pada langkah 4 hingga 5, Node-A dipause selama 350 milidetik, melampaui batas TTL 200 milidetik sehingga lease kedaluwarsa. Menariknya, sistem sempat berada pada kondisi transient dual-leader di mana kedua node menganggap dirinya sebagai LEADER aktif. Namun, saat Node-A terbangun dan mencoba menulis dengan token 1, lapisan penyimpanan menolak operasi tersebut secara seketika karena token 1 lebih kecil dari token 2 milik Node-B.",
    narrative_en="This failure path is extracted directly from the executable lab demo. In steps 4 through 5, Node-A is paused for 350ms, comfortably exceeding the 200ms TTL window and triggering lease expiration. Crucially, the cluster momentarily enters a transient dual-leader state where both nodes view themselves as active leaders. However, when Node-A wakes up and attempts a write with token 1, storage intercepts the call and rejects it because token 1 is less than Node-B's token 2.",
    mental_model="Transient dual-leader states are inevitable in asynchronous networks; fencing tokens make them safe.",
    source_path="cmd/demo/main.go:61-75 • engineering/03-execution-result.md:61-70",
    slide_num="07 / 24",
    layout="failure-path"
)

# Slide 8: Tabel perbandingan mekanisme vs gejala
add_slide(
    categories=["PERBANDINGAN", "MEKANISME"],
    title="Mekanisme vs Gejala: Evaluasi Solusi Koordinasi",
    subtitle="Mengapa solusi parsial tetap menyisakan celah korupsi data pada kegagalan terdistribusi.",
    content_items=[
        {"title": "TABEL PERBANDINGAN MEKANISME", "items": [
            "Pendekatan | Mekanisme Koordinasi | Gejala / Failure Mode | Garansi Split-Brain",
            "Heartbeat Naif | Ping status via interval | Network partition memicu dual-leader | Tidak Aman (korupsi data)",
            "Lease TTL Murni | Kunci kedaluwarsa otomatis | GC pause memicu stale write pasca-wake | Tidak Aman (TOCTOU gap)",
            "Distributed Lock | Mutex terpusat (Redis SETNX) | Lock lepas sebelum job tuntas | Sebagian (hanya efisiensi)",
            "Lease + Fencing | Monotonic revision + check-and-set | Stale write dihentikan di resource | Aman Mutlak (linierizable)"
        ], "type": "table"}
    ],
    narrative="Tabel ini membandingkan empat pendekatan koordinasi terdistribusi dalam menghadapi kegagalan jaringan dan runtime pause. Heartbeat naif dan lease TTL murni gagal melindungi integritas data karena keduanya rentan terhadap jeda eksekusi dan kondisi split-brain. Penggunaan distributed lock tanpa token urutan hanya cocok untuk optimalisasi efisiensi seperti deduplikasi job. Hanya penggabungan antara lease berbatas waktu dengan fencing token monotonik yang memberikan garansi keamanan mutlak pada level resource.",
    narrative_en="This table evaluates four distributed coordination patterns under network partitions and runtime pauses. Naive heartbeats and bare TTL leases fail to preserve data integrity because both remain susceptible to execution pauses and split-brain phenomena. Distributed locks devoid of generation ordering only suffice for efficiency optimizations like job deduplication. Combining time-bounded leases with monotonic fencing tokens guarantees absolute write safety at the storage boundary.",
    mental_model="Locks optimize efficiency; monotonic fencing tokens protect data correctness.",
    source_path="research/03-mechanisms-and-leases.md • research/05-system-comparison-and-best-practices.md",
    slide_num="08 / 24",
    layout="comparison-table"
)

# Slide 9: Decision sequence (3 focused decision cards)
add_slide(
    categories=["DESIGN DECISION", "ARSITEKTUR"],
    title="3 Keputusan Desain Kritis Leader Election",
    subtitle="Urutan keputusan arsitektur untuk membangun koordinasi terdistribusi yang tangguh.",
    content_items=[
        {"title": "1. MONOTONIC CLOCK FOR LEASES", "items": [
            "Wajib gunakan time.Now() (monotonic reading) untuk evaluasi TTL dan renewal.",
            "Mencegah anomali perpindahan waktu akibat sinkronisasi NTP step clock drift.",
            "Terapkan interval renewal jauh lebih cepat dari TTL (contoh: renewRate <= TTL / 3)."
        ], "type": "card"},
        {"title": "2. STRICT MONOTONIC REVISION", "items": [
            "Coordinator wajib mengincrement counter revision secara atomik di setiap lease baru.",
            "Hindari penggunaan UUID atau hash acak karena tidak memiliki relasi urutan (ordering).",
            "Renewal pada holder yang sama mempertahankan nilai token yang sudah ada."
        ], "type": "card"},
        {"title": "3. ATOMIC CHECK-AND-SET GATE", "items": [
            "Storage wajib menegakkan aturan incomingToken > lastSeenToken di bawah mutex/CAS.",
            "Tolak operasi dengan status kode kesalahan eksplisit (ErrStaleFencingToken).",
            "Simpan histori watermark generasi tertinggi yang pernah berhasil diterima."
        ], "type": "card"}
    ],
    narrative="Implementasi koordinasi yang aman bergantung pada tiga keputusan desain fundamental. Pertama, perhitungan masa berlaku lease harus menggunakan monotonic clock untuk menghindari disrupsi penyesuaian jam sistem via NTP. Kedua, setiap pergantian kepemimpinan wajib menerbitkan revision token yang monotonik naik; UUID acak dilarang karena tidak memiliki relasi urutan kausal. Ketiga, gerbang validasi di lapisan penyimpanan bersama harus mengimplementasikan check-and-set atomik yang menolak token yang tidak lebih besar.",
    narrative_en="Building secure coordination hinges upon three foundational architectural decisions. First, lease lifetime computations must rely strictly on monotonic clock readings to eliminate clock jumps from NTP synchronization. Second, every leadership handover must allocate a strictly monotonic revision token; random UUIDs are unacceptable because they lack causal ordering. Third, the validation gate at the shared resource layer must enforce atomic check-and-set semantics that reject non-increasing tokens.",
    mental_model="Monotonic clocks measure duration; monotonic counters establish immutable authority sequence.",
    source_path="labs/30-leader-election/content/05-key-takeaways.md",
    slide_num="09 / 24",
    layout="decision-cards"
)

# Slide 10: Bukti eksekusi / terminal output
add_slide(
    categories=["BUKTI", "DEMO"],
    title="Bukti Eksekusi: Siklus Lengkap Election & Fencing",
    subtitle="Output langsung dari go run ./cmd/demo memperlihatkan transisi kepemimpinan dan integritas storage.",
    content_items=[
        {"title": "DEMO EXECUTION TRACE", "items": [
            "=== Distributed Leader Election with Fencing Tokens Demo ===",
            "[1] Starting Node-A and Node-B election campaigns...",
            "Status: Node-A=LEADER, Node-B=FOLLOWER",
            "[2] Leader elected: Node-A (Fencing Token: 1)",
            "[3] Node-A performing legitimate fenced write... Write SUCCESS by Node-A (Token 1)",
            "[4] Simulating Stop-the-World GC Pause / Network Partition on Node-A (350ms > TTL 200ms)...",
            "[5] Status after TTL expiration: Node-A=LEADER, Node-B=LEADER",
            "[6] New Leader promoted: Node-B (Fencing Token: 2)",
            "[7] New Leader Node-B performing fenced write... Write SUCCESS by Node-B (Token 2)",
            "[8] Old Leader Node-A wakes up from GC pause (still holding stale token 1)...",
            "[9] Old Leader Node-A attempts write to shared storage using stale token 1...",
            "Write REJECTED by FencedStorage: stale fencing token: rejected to prevent split-brain write: token 1 <= last seen 2 (author: Node-A)",
            "Result: Split-brain prevented! Old leader's stale write was safely rejected.",
            "[10] Shared Storage History:",
            "  1. Token: 1 | Author: Node-A | Value: record-from-initial-leader",
            "  2. Token: 2 | Author: Node-B | Value: record-from-failover-leader",
            "=== Demo completed successfully ==="
        ], "type": "terminal"}
    ],
    narrative="Output terminal ini membuktikan bahwa seluruh siklus kepemimpinan berjalan sesuai spesifikasi desain. Node-A pertama kali dilantik dengan Token 1 dan berhasil menulis record pertama. Saat Node-A dibekukan selama 350 milidetik, Node-B mengambil alih kepemimpinan dengan Token 2 dan menulis record kedua. Ketika Node-A terbangun dan berusaha menulis dengan Token 1, sistem menolak operasi tersebut secara eksplisit. Riwayat penyimpanan bersama membuktikan hanya ada dua entri sah dan tidak ada data yang tertimpa oleh mutasi usang.",
    narrative_en="This terminal execution trace verifies that the full leadership lifecycle operates precisely as designed. Node-A initially attains leadership with Token 1 and commits the first record. When Node-A is paused for 350ms, Node-B succeeds it with Token 2 and writes the second record. When Node-A resumes and attempts to write using Token 1, the storage layer rejects the request. The final audit history confirms exactly two legitimate records exist, completely clean of any overwritten state.",
    mental_model="Terminal logs provide indisputable evidence: stale writes fail loudly, preserving storage state.",
    source_path="cmd/demo/main.go:1-78 • engineering/03-execution-result.md:52-77",
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
            "Coordinator bersifat single-node in-memory; tidak membuktikan konsensus multi-node Raft riil.",
            "Tidak menguji fault injeksi disk fisik, network loss paket parsial, atau Byzantine failures.",
            "Clock offset diuji melalui manipulasi software; tidak menguji hardware clock skew ekstrem.",
            "Storage diimplementasikan via in-memory sync.RWMutex; bukan PostgreSQL CAS atau etcd revision check.",
            "Nilai waktu (TTL 200ms, renew 50ms) bersifat ilustratif lab; bukan panduan tuning produksi.",
            "Tidak mencakup mekanisme automated fencing token compaction atau database sequence recycling."
        ], "type": "bullet"}
    ],
    narrative="Klaim teknis dalam laboratorium ini dibatasi oleh parameter lingkungan pengujian yang disederhanakan. Coordinator yang digunakan berjalan in-memory di dalam satu proses Go, bukan kluster terdistribusi multi-node berbasis konsensus Raft seperti etcd atau Paxos. Uji coba tidak mensimulasikan kegagalan jaringan fisik nyata seperti packet loss berkala, korupsi hardware disk, atau perilaku Byzantine. Nilai parameter waktu seperti TTL 200 milidetik hanya dirancang untuk kecepatan eksekusi demo dan tidak boleh langsung diaplikasikan ke lingkungan produksi tanpa benchmarking p99 latency.",
    narrative_en="The technical claims established in this lab are explicitly bounded by the in-process test harness. The coordinator runs in-memory within a single Go runtime rather than across a multi-node Raft consensus cluster such as etcd or ZooKeeper. Real-world physical network degradation like packet loss, storage corruption, or Byzantine behaviors were not exercised. Furthermore, illustrative timing constants (200ms TTL, 50ms renewal) were chosen for rapid test turnaround and must be re-calibrated against production p99 network latency and GC statistics.",
    mental_model="Acknowledge lab constraints to avoid confusing conceptual proofs with production-grade topology.",
    source_path="labs/30-leader-election/content/01-content-brief.md • engineering-audit/06-verdict.md",
    slide_num="11 / 24",
    layout="constraints-list"
)

# Slide 12: Trade-off arsitektur
add_slide(
    categories=["TRADE-OFF", "ARSITEKTUR"],
    title="Trade-off Arsitektur: Fencing vs Kompleksitas",
    subtitle="Menilai biaya operasional penambahan gerbang fencing dibandingkan manfaat konsistensi data.",
    content_items=[
        {"title": "MANFAAT & KEUNGGULAN", "items": [
            "Jaminan Konsistensi Mutlak: Menghilangkan risiko split-brain write pada shared storage.",
            "Decoupled Failure Recovery: Node baru dapat langsung menulis tanpa menunggu node lama mengonfirmasi mati.",
            "Auditability: Setiap record memiliki fencing token yang merefleksikan generasi pemimpin penulis."
        ], "type": "bullet"},
        {"title": "KONSEKUENSI & BIAYA", "items": [
            "Modifikasi Storage Layer: Resource bersama wajib mendukung skema token dan check-and-set.",
            "Failover Latency: Masa transisi pemimpin baru terikat pada durasi TTL lease yang dipilih.",
            "Client Error Handling: Aplikasi klien wajib memiliki alur penanganan saat write ditolak ErrStaleFencingToken."
        ], "type": "bullet"}
    ],
    narrative="Penerapan fencing tokens menuntut kompromi arsitektural antara keamanan data dan kompleksitas sistem. Keunggulan utamanya adalah eliminasi total terhadap korupsi data akibat split-brain serta kemampuan failover tanpa perlu menunggu konfirmasi kematian node lama. Namun, biaya penerapannya tidak kecil: lapisan penyimpanan harus dimodifikasi agar mampu mengevaluasi token, waktu failover dibatasi oleh durasi TTL lease, dan klien aplikasi harus memiliki logika fallback ketika transaksi ditolak akibat token kedaluwarsa.",
    narrative_en="Adopting fencing tokens introduces trade-offs between guaranteed data safety and overall system complexity. The primary benefit is absolute immunity from split-brain data corruption and the freedom to promote new leaders without waiting for explicit old-leader termination. Conversely, the operational overhead requires modifying backend storage engines to support check-and-set semantics, accepting failover dead-time proportional to the lease TTL, and engineering client resilience when writes are rejected.",
    mental_model="Safety requires storage collaboration; you trade write simplicity for bulletproof consistency.",
    source_path="labs/30-leader-election/content/02-master-draft.md",
    slide_num="12 / 24",
    layout="tradeoff"
)

# Slide 13: Concurrency / validation pipeline
add_slide(
    categories=["CONCURRENCY", "VALIDATION"],
    title="Validasi Konkurensi: 5 Node Berlomba Rebut Lease",
    subtitle="TestConcurrentElectionRace membuktikan linearitas coordinator: tepat 1 leader terpilih.",
    content_items=[
        {"title": "KODE TEST KONKURENSI", "items": [
            "func TestConcurrentElectionRace(t *testing.T) {",
            "    coord := coordinator.New()",
            "    store := storage.New()",
            "    nodes := make([]*candidate.Node, 5)",
            "    for i := 0; i < 5; i++ {",
            "        nodes[i] = candidate.NewNode(string(rune('A'+i)), \"leader-lock\", 100*time.Millisecond, 20*time.Millisecond, coord, store)",
            "        nodes[i].Start(ctx)",
            "    }",
            "    time.Sleep(100 * time.Millisecond)",
            "    // Assert: leaderCount == 1 di antara 5 node yang berlomba",
            "}"
        ], "type": "code"},
        {"title": "HASIL VERIFIKASI", "items": [
            "Race Detector: PASS (bebas dari data race condition).",
            "Linearizable Mutex: Semua permintaan Acquire dan Renew dieksekusi secara serial.",
            "Tepat 1 Leader: 4 kandidat lainnya tetap berstatus FOLLOWER secara stabil."
        ], "type": "bullet"}
    ],
    narrative="Uji konkurensi membuktikan ketahanan implementasi saat lima node kandidat mencoba mengakuisisi lease secara simultan. Coordinator memanfaatkan sync.Mutex internal untuk memastikan setiap pengecekan dan penetapan lease berjalan secara linearizable. Dalam pengujian TestConcurrentElectionRace, setelah 100 milidetik kompetisi intensif, tepat satu node berhasil dipromosikan menjadi LEADER sementara empat node lainnya tertahan sebagai FOLLOWER. Verifikasi dengan flag -race membuktikan tidak ada data race di seluruh memori proses.",
    narrative_en="Concurrency stress testing validates implementation correctness when five candidate nodes compete for lease ownership simultaneously. The coordinator leverages an internal sync.Mutex to guarantee linearizable evaluation of lease acquisition and renewal requests. In TestConcurrentElectionRace, after 100ms of concurrent campaigning, exactly one node transitions to LEADER while the remaining four stably remain as FOLLOWERS. Running under Go's race detector flags zero data races across all components.",
    mental_model="Serialized coordinator primitives resolve concurrent contention into a single unambiguous leader.",
    source_path="tests/election_test.go:156-191",
    slide_num="13 / 24",
    layout="concurrency-flow"
)

# Slide 14: Before vs After comparison cards
add_slide(
    categories=["BEFORE/AFTER", "PERUBAHAN"],
    title="Before vs After: Evolusi Sistem Koordinasi",
    subtitle="Transformasi dari koordinasi rapuh berbasis waktu lokal menuju arsitektur fenced mutlak.",
    content_items=[
        {"title": "BEFORE (Lease TTL Murni)", "items": [
            "Koordinasi: Node mengandalkan local state check sebelum eksekusi.",
            "Failure Mode: Stop-the-world pause memicu penulisan ganda pasca-wake.",
            "Penyimpanan: Storage pasif menerima semua data tanpa validasi urutan.",
            "Dampak: Korupsi data senyap dan state database tidak konsisten."
        ], "type": "card"},
        {"title": "AFTER (Lease + Fencing Tokens)", "items": [
            "Koordinasi: Lease berbatas waktu dipadukan dengan revision counter monotonik.",
            "Failure Mode: Leader lama terbangun, namun mutasinya otomatis ditolak.",
            "Penyimpanan: Storage aktif memeriksa aturan incomingToken > lastSeenToken.",
            "Dampak: Integritas data terjamin 100% aman dari ancaman split-brain."
        ], "type": "card"}
    ],
    narrative="Perbandingan before-after ini menunjukkan pergeseran paradigma dalam menangani koordinasi terdistribusi. Sebelum menerapkan fencing tokens, sistem sangat rapuh terhadap jeda eksekusi karena penyimpanan mempercayai klaim kepemimpinan klien secara pasif. Akibatnya, jeda eksekusi leader lama dapat merusak data yang baru saja ditulis oleh leader pengganti. Setelah menerapkan arsitektur fencing, penyimpanan bersama bertindak sebagai filter aktif yang menjamin data tidak akan pernah tertimpa oleh mutasi yang berasal dari generasi kepemimpinan masa lalu.",
    narrative_en="This before-and-after comparison highlights the architectural paradigm shift in distributed coordination. Prior to adopting fencing tokens, the topology was exceptionally fragile under execution pauses because storage passively trusted client assertions. As a result, paused legacy leaders overwrote fresh state written by active successors. With the fencing token pattern in place, shared storage functions as an active verification barrier, ensuring mutations from obsolete leadership epochs are systematically rejected.",
    mental_model="Before: passive storage vulnerable to split-brain. After: active fencing guarantees immutable epochs.",
    source_path="labs/30-leader-election/content/02-master-draft.md",
    slide_num="14 / 24",
    layout="before-after"
)

# Slide 15: Analisis mendalam - Test 1 Lease Acquisition & Renewal
add_slide(
    categories=["BUKTI TEST", "LEASE & TTL"],
    title="Test 1: Akuisisi, Renewal, dan Kedaluwarsa Lease",
    subtitle="TestCoordinatorLeaseAcquisitionAndRenewal membuktikan siklus hidup lease di coordinator.",
    content_items=[
        {"title": "TEST KONTRAK LEASE", "items": [
            "func TestCoordinatorLeaseAcquisitionAndRenewal(t *testing.T) {",
            "    coord := coordinator.New()",
            "    ttl := 200 * time.Millisecond",
            "    lease1, err := coord.Acquire(\"leader-lock\", \"node-1\", ttl)",
            "    // Assert: lease1.FencingToken == 1",
            "    _, err = coord.Acquire(\"leader-lock\", \"node-2\", ttl)",
            "    // Assert: errors.Is(err, coordinator.ErrLeaseHeld)",
            "    time.Sleep(250 * time.Millisecond) // Lampaui TTL 200ms",
            "    lease2, err := coord.Acquire(\"leader-lock\", \"node-2\", ttl)",
            "    // Assert: lease2.FencingToken == 2",
            "}"
        ], "type": "code"},
        {"title": "VERIFIKASI INVARIAN", "items": [
            "Mutual Exclusion: Node-2 ditolak ErrLeaseHeld selama lease Node-1 masih dalam masa aktif.",
            "Token Stability: Memanggil Renew() memperpanjang masa aktif tanpa mengubah nilai token (tetap 1).",
            "Monotonic Increment: Setelah TTL kedaluwarsa, akuisisi oleh Node-2 menghasilkan token 2."
        ], "type": "bullet"}
    ],
    narrative="Unit test pertama membuktikan kebenaran siklus hidup lease pada modul coordinator. Node-1 berhasil memperoleh lease pertama dengan fencing token bernilai 1. Selama lease masih aktif, permintaan akuisisi dari Node-2 langsung ditolak dengan error ErrLeaseHeld. Operasi pembaruan berkala (Renew) terbukti memperpanjang masa berlaku lease tanpa menaikkan nilai token. Setelah dilakukan sleep selama 250 milidetik yang melampaui batas TTL 200 milidetik, lease kedaluwarsa secara otomatis dan Node-2 berhasil mengakuisisi lease dengan token 2.",
    narrative_en="The first unit test validates the operational lifecycle of leases managed by the coordinator module. Node-1 successfully acquires the initial lease, receiving fencing token 1. While this lease remains valid, Node-2's concurrent acquisition attempt is blocked with ErrLeaseHeld. Heartbeat renewal operations demonstrably extend the lease expiration without perturbing the token value. After pausing for 250ms to exceed the 200ms TTL window, the lease expires cleanly, enabling Node-2 to secure ownership with token 2.",
    mental_model="Lease renewal preserves the generation token; lease expiration monotonically advances it.",
    source_path="tests/election_test.go:14-56",
    slide_num="15 / 24",
    layout="test-analysis"
)

# Slide 16: Analisis mendalam - Test 2 Fenced Storage Rejection
add_slide(
    categories=["BUKTI TEST", "STORAGE REJECTION"],
    title="Test 2: Storage Menolak Token Usang & Duplikat",
    subtitle="TestFencedStorageRejectsStaleTokens memverifikasi syarat strict inequality incoming > lastSeen.",
    content_items=[
        {"title": "TEST PENOLAKAN STORAGE", "items": [
            "func TestFencedStorageRejectsStaleTokens(t *testing.T) {",
            "    store := storage.New()",
            "    _ = store.Write(\"node-1\", 10, \"state-1\") // Sukses: 10 > 0",
            "    _ = store.Write(\"node-2\", 11, \"state-2\") // Sukses: 11 > 10",
            "    err1 := store.Write(\"node-1\", 10, \"stale-state\")",
            "    // Assert: errors.Is(err1, storage.ErrStaleFencingToken)",
            "    err2 := store.Write(\"node-2\", 11, \"duplicate-token\")",
            "    // Assert: errors.Is(err2, storage.ErrStaleFencingToken)",
            "    // Assert: len(store.History()) == 2 (mutasi ilegal tidak dicatat)",
            "}"
        ], "type": "code"},
        {"title": "ATURAN STRICT INEQUALITY", "items": [
            "Token Mundur Ditolak: Token 10 ditolak setelah storage melihat token 11 (10 < 11).",
            "Token Sama Ditolak: Token 11 duplikat ditolak karena tidak memenuhi syarat strict greater (11 == 11).",
            "Integritas Riwayat: Penyimpanan bersama hanya berisi tepat 2 record legal."
        ], "type": "bullet"}
    ],
    narrative="Unit test kedua menguji ketatnya aturan validasi pada lapisan FencedStorage. Penyimpanan pertama kali menerima token 10, kemudian menerima token 11. Saat Node-1 mencoba menulis kembali dengan token 10, operasi ditolak karena 10 lebih kecil dari 11. Begitu pula saat Node-2 mencoba menulis ulang dengan token 11 yang identik, storage menolaknya karena aturan mensyaratkan token harus strictly lebih besar. Hasil pemeriksaan History() membuktikan riwayat penyimpanan hanya berisi tepat dua data yang sah dan menolak semua mutasi ilegal.",
    narrative_en="The second unit test rigorously examines the validation invariants implemented inside FencedStorage. The storage engine sequentially admits tokens 10 and 11. When Node-1 subsequently re-submits a write under token 10, the call is rejected since 10 is strictly less than 11. Similarly, when Node-2 attempts a duplicate write under token 11, storage rejects it because the contract requires strict inequality. Inspecting History() confirms exactly two legitimate records exist, guaranteeing invalid writes never pollute state.",
    mental_model="Strict inequality incomingToken > lastSeenToken closes both stale replay and duplicate mutation vulnerabilities.",
    source_path="tests/election_test.go:58-85",
    slide_num="16 / 24",
    layout="test-analysis"
)

# Slide 17: Analisis mendalam - Test 3 Failover & Split-Brain Defense
add_slide(
    categories=["BUKTI TEST", "SPLIT-BRAIN DEFENSE"],
    title="Test 3: Failover Otomatis & Pertahanan Split-Brain",
    subtitle="TestLeaderElectionFailoverAndSplitBrainDefense membuktikan skenario GC pause end-to-end.",
    content_items=[
        {"title": "SIMULASI GC PAUSE PADA NODE AKTIF", "items": [
            "// 1. Leader-1 terpilih dan menulis sukses dengan token1",
            "leader1.SimulatePause(250 * time.Millisecond) // Lebih lama dari TTL 150ms",
            "time.Sleep(200 * time.Millisecond)",
            "// 2. Leader-2 mengambil alih kepemimpinan dengan token2 > token1",
            "err = leader2.PerformFencedWrite(\"leader-2-write\") // SUKSES",
            "// 3. Leader-1 bangun dan mencoba menulis dengan token1 usang",
            "err = store.Write(leader1.ID(), token1, \"stale-write-after-gc-pause\")",
            "// Assert: errors.Is(err, storage.ErrStaleFencingToken)"
        ], "type": "code"},
        {"title": "FAKTA PERILAKU SISTEM", "items": [
            "Leader 1 dibekukan selama 250ms melalui SimulatePause yang menghentikan loop renewal.",
            "Kandidat kedua mendeteksi lease habis dan dipromosikan dengan token lebih tinggi.",
            "Saat Leader 1 bangun, percobaan tulisnya digagalkan oleh pagar validasi storage."
        ], "type": "bullet"}
    ],
    narrative="Test ketiga mendemonstrasikan skenario inti split-brain yang dipecahkan oleh lab ini. Node leader pertama disimulasikan mengalami stop-the-world pause selama 250 milidetik, melampaui batas lease TTL 150 milidetik. Kandidat kedua mendeteksi kekosongan lease, mengambil alih kepemimpinan, dan sukses menulis dengan token 2. Ketika leader pertama bangun dan mencoba menulis ke storage menggunakan token 1 yang dipegangnya, operasi ditolak mentah-mentah dengan ErrStaleFencingToken. Tes ini mengonfirmasi pertahanan split-brain bekerja sempurna.",
    narrative_en="The third test orchestrates the pivotal split-brain scenario solved by this lab. The primary leader is subjected to a simulated stop-the-world pause of 250ms, deliberately exceeding the 150ms TTL threshold. The standby candidate detects the expired lease, claims leadership, and commits state under token 2. When the original leader resumes and attempts to flush data using its stale token 1, the write is stopped cold by ErrStaleFencingToken. This verifies end-to-end split-brain defense under real execution pauses.",
    mental_model="Simulated pause triggers clean failover; resumed stale leader is intercepted at storage boundary.",
    source_path="tests/election_test.go:87-154",
    slide_num="17 / 24",
    layout="test-analysis"
)

# Slide 18: Studi Kasus Industri: etcd vs ZooKeeper vs Redis Redlock
add_slide(
    categories=["STUDI KASUS", "INDUSTRI"],
    title="Studi Kasus Industri: Mekanisme Fencing Nyata",
    subtitle="Bagaimana sistem terdistribusi kelas dunia mengimplementasikan ordering tokens.",
    content_items=[
        {"title": "SISTEM KOORDINASI MODERN", "items": [
            "etcd v3: Menggunakan Lease TTL dipadukan dengan Revision Raft 64-bit sebagai fencing token bawaan.",
            "Apache ZooKeeper: Menghasilkan Ephemeral Sequential Znodes; nomor urut zxid bertindak sebagai token.",
            "Consul: Mengandalkan internal Raft sessions dengan Session TTL dan periodic health checks."
        ], "type": "bullet"},
        {"title": "DEBAT REDIS REDLOCK (KLEPPMANN VS ANTIREZ)", "items": [
            "Martin Kleppmann (2016): Redlock tidak menghasilkan fencing token monotonik; rentan GC pause & NTP jump.",
            "Salvatore Sanfilippo (antirez): Redlock dirancang untuk timed lock praktis, bukan linearizable consensus.",
            "Konsensus Industri: Jika resource tidak mendukung fencing check, Redlock tidak menjamin mutual exclusion mutlak."
        ], "type": "bullet"}
    ],
    narrative="Riset sistem terdistribusi menunjukkan kesepakatan luas tentang perlunya fencing tokens. etcd memanfaatkan nomor Raft Revision 64-bit yang selalu bertambah sebagai fencing token alami untuk setiap transaksi mutasi. Apache ZooKeeper menggunakan zxid dan urutan pembuatan ephemeral znode. Debat bersejarah antara Martin Kleppmann dan antirez mengenai algoritma Redis Redlock menegaskan prinsip utama: kunci terdistribusi tanpa fencing token monotonik hanya cukup untuk efisiensi sistem, bukan untuk keamanan data yang membutuhkan korektivitas mutlak.",
    narrative_en="Industry research reveals widespread consensus regarding the necessity of fencing tokens. etcd leverages its 64-bit Raft Revision counter as a native, strictly monotonic fencing token across transactional mutations. Apache ZooKeeper provides zxid timestamps and sequential ephemeral node numbers. The famous debate between Martin Kleppmann and Salvatore Sanfilippo (antirez) over Redis Redlock established a lasting maxim: distributed locks lacking monotonic tokens serve efficiency purposes, never safety-critical correctness.",
    mental_model="Without monotonic fencing tokens at the storage boundary, distributed locks provide efficiency, not safety.",
    source_path="research/04-fencing-and-redlock.md • research/05-system-comparison-and-best-practices.md",
    slide_num="18 / 24",
    layout="cards"
)

# Slide 19: Analisis Teknis: Monotonic Clock vs Wall-Clock
add_slide(
    categories=["ANALISIS TEKNIS", "RUNTIME"],
    title="Mengapa Lease Wajib Menggunakan Monotonic Clock?",
    subtitle="Bahaya lompatan wall-clock akibat penyesuaian NTP terhadap perhitungan kedaluwarsa lease.",
    content_items=[
        {"title": "ANOMALI WALL-CLOCK (REAL-TIME)", "items": [
            "Sinkronisasi waktu jaringan (NTP step) dapat memajukan atau memundurkan wall-clock sistem.",
            "Jika jam melompat mundur: Lease dianggap aktif selamanya, mematikan kemampuan failover otomatis.",
            "Jika jam melompat maju: Lease dianggap kedaluwarsa prematur saat leader masih sehat."
        ], "type": "bullet"},
        {"title": "GARANSI MONOTONIC CLOCK DI GO", "items": [
            "time.Now() di Go membaca monotonic clock hardware secara otomatis untuk pengukuran durasi.",
            "time.Since(start) dan now.Before(expiresAt) kebal terhadap manipulasi penyesuaian jam sistem.",
            "Coordinator dan Candidate lab ini murni menggunakan durasi monotonik untuk semua logika TTL."
        ], "type": "bullet"}
    ],
    narrative="Penggunaan tipe jam yang tepat merupakan detail implementasi yang membedakan sistem amatir dari sistem kelas produksi. Wall-clock (jam dinding) dirancang untuk mencatat waktu kalender manusia dan dapat melompat mundur atau maju secara drastis saat terjadi penyesuaian NTP. Jika perhitungan lease bergantung pada wall-clock, lompatan waktu dapat memicu kegagalan fatal: failover palsu atau lease yang terkunci abadi. Go standard library mengatasi ini dengan menyematkan pembacaan monotonic hardware clock pada fungsi time.Now(), menjamin perhitungan durasi TTL selalu linier.",
    narrative_en="Selecting the correct clock abstraction separates resilient distributed architectures from fragile prototypes. System wall-clocks track human calendar time and can step backward or forward abruptly during NTP synchronizations. If lease calculations depend on wall-clocks, time jumps cause severe anomalies: phantom failovers or indefinite lease starvation. The Go standard library eliminates this vulnerability by transparently embedding hardware monotonic clock readings inside time.Now(), guaranteeing linear duration calculations.",
    mental_model="Wall-clocks tell human time and jump; monotonic clocks measure duration and never step backwards.",
    source_path="labs/30-leader-election/research/03-mechanisms-and-leases.md • internal/coordinator/coordinator.go:40",
    slide_num="19 / 24",
    layout="cards"
)

# Slide 20: Checklist implementasi
add_slide(
    categories=["CHECKLIST", "PRODUKSI"],
    title="Checklist Implementasi Leader Election & Fencing",
    subtitle="9 poin verifikasi wajib sebelum menggelar sistem koordinasi ke lingkungan produksi.",
    content_items=[
        {"title": "CHECKLIST ARSITEKTUR KOORDINASI", "items": [
            "[ ] Parameter TTL lease dikalibrasi berdasarkan observasi empiris p99 GC pause dan RTT jaringan.",
            "[ ] Loop renewal diatur dengan frekuensi aman (renewRate <= TTL / 3) dengan backoff terukur.",
            "[ ] Semua perhitungan durasi lease dan renewal murni menggunakan monotonic clock.",
            "[ ] Coordinator menghasilkan revision counter yang strictly monotonically increasing di setiap lease baru.",
            "[ ] Pembaharuan lease (renew) oleh pemegang yang sama mempertahankan fencing token yang ada.",
            "[ ] Shared resource (database/storage) mengimplementasikan check-and-set atomik (incoming > lastSeen).",
            "[ ] Permintaan tulis dengan token <= lastSeenToken ditolak dengan error eksplisit yang dapat dipantau.",
            "[ ] Metric alerting dipasang untuk mendeteksi lonjakan kejadian penolakan token basi (stale write spike).",
            "[ ] Sistem aplikasi memahami batas antara timed lock untuk efisiensi vs linearizable fencing untuk keamanan."
        ], "type": "bullet"}
    ],
    narrative="Daftar periksa ini merangkum sembilan prinsip krusial yang wajib dipenuhi sebelum menerapkan sistem leader election di produksi. Penentuan batas waktu TTL dan laju renewal tidak boleh ditebak, melainkan harus diukur berdasarkan data empiris p99 latency jaringan dan durasi GC pause runtime. Coordinator wajib menyediakan nomor generasi monotonik, sementara shared storage wajib mengisolasi validasi token di balik operasi atomik. Setiap penolakan stale token harus dimonitor sebagai indikasi awal degradasi performa atau partisi jaringan.",
    narrative_en="This checklist codifies nine critical rules required before deploying distributed leader election systems to production. Lease TTLs and heartbeat renewal intervals must never be guessed; they require empirical calibration against observed p99 network RTT and runtime GC pause distributions. Coordinators must issue monotonically increasing generation numbers, and shared resources must isolate token validations within atomic check-and-set transactions. Any rejection spike must trigger alerts indicating node stalls or network partitions.",
    mental_model="Calibrate timeouts empirically, increment tokens monotonically, and enforce boundaries atomically.",
    source_path="labs/30-leader-election/content/02-master-draft.md:225-236",
    slide_num="20 / 24",
    layout="checklist"
)

# Slide 21: Anti-pattern & cara reproduksi kegagalan
add_slide(
    categories=["ANTI-PATTERN", "REPRODUKSI"],
    title="Anti-Pattern: Kesalahan Desain yang Harus Dihindari",
    subtitle="Kekeliruan arsitektur yang melumpuhkan jaminan keselamatan koordinasi terdistribusi.",
    content_items=[
        {"title": "ANTI-PATTERN UMUM", "items": [
            "Fencing di Klien: Melakukan validasi token di memori klien sebelum memanggil storage (rentan TOCTOU).",
            "Menerima Token Sama: Menggunakan syarat token >= lastSeen membuka celah eksekusi mutasi ganda.",
            "Token Non-Monotonik: Memakai UUID acak atau timestamp lokal yang tidak memiliki garansi urutan kausal.",
            "Wall-Clock TTL: Menggunakan jam dinding sistem rentan manipulasi penyesuaian NTP clock skew."
        ], "type": "card"},
        {"title": "CARA REPRODUKSI GAGAL DI LAB", "items": [
            "1. Ubah internal/storage/storage.go baris 37 menjadi token < s.lastSeenToken (izinkan token sama).",
            "2. Jalankan go test -v ./tests: TestFencedStorageRejectsStaleTokens langsung gagal pada duplicate check.",
            "3. Hapus validasi storage sama sekali: stale write dari Node-A akan menimpa histori sah milik Node-B!"
        ], "type": "card"}
    ],
    narrative="Mengabaikan prinsip dasar koordinasi dapat berakibat fatal bagi integritas data perusahaan. Anti-pattern yang paling sering dijumpai adalah memvalidasi status kepemimpinan di sisi klien aplikasi alih-alih di gerbang penyimpanan, membuka celah TOCTOU yang lebar. Mengizinkan token yang bernilai sama (menggunakan operator lebih besar atau sama dengan) juga merusak keunikan generasi pemimpin. Untuk membuktikannya di lab, kita cukup melonggarkan validasi di storage.go dan melihat unit test langsung mendeteksi pelanggaran keamanan.",
    narrative_en="Violating distributed coordination fundamentals compromises enterprise data consistency. The most prevalent anti-pattern is executing leadership validation exclusively in client memory instead of at the storage gate, inviting catastrophic TOCTOU race conditions. Accepting equal tokens by substituting strict inequality with greater-than-or-equal operators likewise shatters epoch exclusivity. This can be reproduced in the lab by relaxing storage.go line 37, causing the test suite to immediately catch the security violation.",
    mental_model="Client-side checks are advisory; only server-side resource fencing guarantees mutual exclusion.",
    source_path="labs/30-leader-election/content/02-master-draft.md:206-213",
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
            "cd labs/30-leader-election",
            "go test -v ./...",
            "go test -race ./...",
            "go run ./cmd/demo"
        ], "type": "code"},
        {"title": "RINGKASAN METRIK LAB", "items": [
            "Test Suite: 4 test cases lulus 100% tanpa kegagalan (1.17s total runtime).",
            "Race Detector: Bersih dari thread contention dan data races (2.01s pass).",
            "Dependencies: Murni Go Standard Library (sync, time, context, errors) tanpa dependensi pihak ketiga."
        ], "type": "bullet"}
    ],
    narrative="Lab ini dirancang agar dapat diverifikasi secara mandiri dalam hitungan detik. Seluruh implementasi dibangun murni menggunakan standard library Go tanpa memerlukan server eksternal seperti Redis atau etcd. Perintah go test -v mengeksekusi empat pengujian komprehensif mulai dari siklus hidup lease hingga pertahanan split-brain. Perintah go test -race memvalidasi keamanan konkurensi di bawah beban goroutine intensif, dan go run demo menampilkan simulasi visual interaktif penolakan mutasi usang.",
    narrative_en="This laboratory is engineered for immediate, frictionless local verification. The entire codebase is implemented with the Go standard library, eliminating external dependencies like Redis clusters or ZooKeeper daemons. Running go test -v executes four exhaustive test scenarios validating lease lifecycles and split-brain defenses. Executing go test -race proves thread safety under multi-goroutine contention, while go run demo delivers an interactive terminal demonstration of stale write rejection.",
    mental_model="Zero-dependency standard library code enables instant, deterministic peer verification.",
    source_path="README.md:37-50 • engineering/03-execution-result.md",
    slide_num="22 / 24",
    layout="command-block"
)

# Slide 23: Prinsip arsitektural jangka panjang
add_slide(
    categories=["PRINSIP", "ARSITEKTUR"],
    title="Prinsip Arsitektur: Desain untuk Asinkronitas",
    subtitle="Pelajaran abadi dalam merancang sistem terdistribusi yang tahan terhadap kegagalan parsial.",
    content_items=[
        {"title": "1. LEASE MEMBERIKAN EFISIENSI, BUKAN SAFETY", "items": [
            "Lease berbatas waktu sangat baik untuk mencegah duplikasi pekerjaan komputasi.",
            "Namun, lease tidak pernah cukup untuk menjamin integritas penulisan pada shared storage non-idempoten."
        ], "type": "card"},
        {"title": "2. RESOURCE ADALAH HAKIM TERTINGGI", "items": [
            "Klien yang mengalami pause tidak pernah tahu bahwa dirinya sudah kedaluwarsa.",
            "Hanya shared resource yang memiliki perspektif global untuk menolak instruksi usang melalui fencing."
        ], "type": "card"},
        {"title": "3. URUTAN GENERASI MENGALAHKAN ASUMSI WAKTU", "items": [
            "Jangan pernah bergantung pada asumsi bahwa proses Anda akan selesai sebelum waktu T habis.",
            "Gunakan nomor urut generasi (monotonic epoch) yang dapat divalidasi secara deterministik."
        ], "type": "card"}
    ],
    narrative="Tiga prinsip arsitektur ini berlaku universal di luar implementasi lab Go ini. Pertama, pahami bahwa lease berbasis waktu hanya memberikan timed mutual exclusion yang cocok untuk optimalisasi efisiensi, bukan keselamatan mutlak. Kedua, tempatkan gerbang keamanan di lapisan resource bersama karena proses klien yang membeku tidak memiliki kesadaran akan berlalunya waktu. Ketiga, gantikan asumsi waktu dengan urutan generasi eksplisit; nomor token yang monotonik naik memberikan jaminan kausalitas yang kebal terhadap jeda sistem terdistribusi.",
    narrative_en="These three architectural principles transcend this specific Go implementation. First, recognize that time-bounded leases merely provide timed mutual exclusion ideal for efficiency optimization, never absolute safety. Second, position the enforcement gate at the shared resource layer because paused clients remain fundamentally blind to elapsed time. Third, replace fragile temporal assumptions with explicit generation numbers; strictly monotonic tokens provide deterministic causal ordering immune to distributed execution pauses.",
    mental_model="Time is subjective in distributed systems; monotonic order at the resource is absolute truth.",
    source_path="labs/30-leader-election/content/05-key-takeaways.md",
    slide_num="23 / 24",
    layout="principles-list"
)

# Slide 24: Closing slide
add_slide(
    categories=["CLOSING"],
    title="Otoritas Terdistribusi Terjaga, Split-Brain Dicegah.",
    subtitle="",
    content_items=[
        {"title": "KOMPONEN ARSITEKTUR LAB", "items": [
            "internal/coordinator/coordinator.go — Linearizable lease & monotonic revision issuer",
            "internal/candidate/candidate.go — Campaign runner, heartbeat renewal, & pause simulation",
            "internal/storage/storage.go — Fenced storage gate with atomic check-and-set",
            "tests/election_test.go — 4 comprehensive suites (concurrency, failover, rejection)",
            "cmd/demo/main.go — Complete interactive lifecycle demonstration"
        ], "type": "bullet"},
        {"title": "VERIFIKASI & REPOSITORI", "items": [
            f"Commit SHA: {COMMIT_SHA[:8]}",
            "Unit Tests: 4/4 PASS (100% clean)",
            "Race Detector: PASS (zero data race)",
            f"Repositori: github.com/lukman-ss/software-engineering-lab/tree/main/labs/{LAB_NO}-{SLUG}"
        ], "type": "bullet"}
    ],
    narrative="Laboratorium Leader Election dengan Fencing Tokens ini membuktikan secara konkret bahwa permasalahan klasik split-brain dan penulisan data basi dapat diselesaikan secara elegan tanpa infrastruktur yang berbelit-belit. Melalui koordinasi lease berbasis waktu yang dipadukan dengan validasi token revisi monotonik naik di lapisan resource bersama, sistem menjamin konsistensi data mutlak bahkan saat node mengalami stop-the-world GC pause yang ekstrem. Seluruh kode sumber, dokumentasi arsitektur, dan instruksi verifikasi tersedia lengkap di repositori.",
    narrative_en="This Leader Election with Fencing Tokens laboratory provides definitive proof that classical split-brain hazards and stale data overwrites can be eliminated cleanly without excessive infrastructure complexity. By harmonizing time-bounded coordinator leases with monotonically increasing revision tokens enforced at the shared storage layer, the architecture guarantees absolute data consistency even under extreme stop-the-world pauses. Complete source code, design specifications, and test harnesses are published in the repository.",
    mental_model="Lease TTL grants temporal leadership; monotonic fencing tokens protect permanent storage.",
    source_path=REPO_URL,
    slide_num="24 / 24",
    layout="closing-dark-navy"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
