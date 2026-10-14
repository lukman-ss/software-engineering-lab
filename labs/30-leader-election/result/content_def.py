#!/usr/bin/env python3
"""
Publication Definition for Software Engineering Lab 30: Leader Election with Fencing Tokens.
Defines 24 slides, metadata, color constants, and content structure.
DO NOT DELETE THIS FILE.
"""
import os
import sys

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/30-leader-election/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/30-leader-election"
COMMIT_SHA = "e963675c67e7d518e19a07ad930dbdb5be73fff1"
LAB_NO = "30"
SLUG = "leader-election"
TITLE_ID = "Leader Election dengan Fencing Tokens"
SUBTITLE_ID = "Mencegah Split-Brain & Stale Writes Menggunakan Monotonic Revision"
TITLE_EN = "Leader Election with Fencing Tokens"
SUBTITLE_EN = "Defending Against Split-Brain & Stale Writes via Monotonic Revisions"
AUTHOR = "Lukman / lukman-ss"

# Color constants
DARK_NAVY = "0b1b2d"
TEAL = "087f76"
CREAM = "f4f1e8"
WHITE = "ffffff"
SLATE = "52606b"
ORANGE = "d97706"

# Slide data structure
SLIDES = []

def add_slide(categories, title, subtitle, content_items, mental_model, source_path="", slide_num=""):
    """Add a slide to the presentation definition."""
    SLIDES.append({
        "categories": categories,
        "title": title,
        "subtitle": subtitle,
        "content_items": content_items,
        "mental_model": mental_model,
        "source": source_path,
        "slide_num": slide_num
    })

# =====================================================================
# 24 SLIDE DEFINITIONS
# =====================================================================

# Slide 1: Cover (Dark Navy)
add_slide(
    categories=["KONTEKS"],
    title=TITLE_ID,
    subtitle="SOFTWARE ENGINEERING LAB | LAB " + LAB_NO,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "Retro CRT monitor beige dengan scanlines bercahaya hijau toska",
            "Diagram terminal interaktif: Node-A vs Node-B lease contention",
            f"Penulis: {AUTHOR}",
        ]}
    ],
    mental_model="",
    source_path="",
    slide_num=""
)

# Slide 2: Problem Definition (The Transient Split-Brain Dilemma)
add_slide(
    categories=["MOTIVASI", "PROBLEM"],
    title="Masalah Mendasar: Transient Split-Brain",
    subtitle="Ketika pause eksekusi melampaui batas waktu lease koordinator",
    content_items=[
        {"title": "BATAS KENDALI TUNGGAL", "items": [
            "Sistem terdistribusi menuntut tepat satu primary worker untuk job eksekusi kritis",
            "Lease-based leader election membatasi masa jabatan pemimpin via Time-To-Live (TTL)",
            "Leader wajib memperbarui sewa secara periodik sebelum TTL habis (heartbeat renewal)",
        ], "type": "bullet"},
        {"title": "ANOMALI ASYNCHRONOUS PAUSE", "items": [
            "Stop-the-world GC pause atau partisi jaringan dapat menghentikan node lebih lama dari TTL",
            "Koordinator menganggap leader mati dan melimpahkan kepemimpinan ke node lain",
            "Node lama terbangun tanpa sadar bahwa masa sewanya telah hangus di koordinator",
        ], "type": "bullet"},
        {"title": "BAHAYA DATA OVERWRITE", "items": [
            "Dua node meyakini dirinya sebagai leader yang sah secara bersamaan (dual-leader window)",
            "Tulisan lama dari node yang tertunda berisiko menimpa mutasi baru di shared storage",
            "Lease temporal semata tidak cukup menjamin mutual exclusion mutlak tanpa proteksi token",
        ], "type": "bullet"},
    ],
    mental_model="A lease grants time-bounded authority; uncoordinated pauses cause silent split-brain writes.",
    source_path="labs/30-leader-election/content/02-master-draft.md:3",
    slide_num="2 / 24"
)

# Slide 3: The Access Card Mental Model
add_slide(
    categories=["MENTAL MODEL", "ANALOGI"],
    title="Mental Model: Pintu Akses & Kartu Bernomor",
    subtitle="Mengubah kewenangan berbasis waktu menjadi urutan angka monotonik",
    content_items=[
        {"title": "ANALOGI KARTU AKSES ELEKTRONIK", "items": [
            "Bayangkan pintu kamar hotel dengan kunci kartu digital sewaan berbatas jam",
            "Penyewa A lupa memperpanjang sewa tepat waktu karena tertidur di luar ruangan",
            "Resepsionis (koordinator) menerbitkan kartu baru untuk Penyewa B dengan nomor urut 2",
        ], "type": "bullet"},
        {"title": "PENOLAKAN DI GAGANG PINTU", "items": [
            "Penyewa A terbangun membawa kartu nomor 1 dan mencoba membuka pintu kamar",
            "Sensor pintu (shared resource) mencatat angka terakhir yang masuk adalah nomor 2",
            "Pintu menolak akses kartu 1 secara instan karena nomornya kalah urutan",
        ], "type": "bullet"},
        {"title": "PRINSIP FENCING MONOTONIK", "items": [
            "Setiap akuisisi sewa baru wajib menaikkan angka token revision secara strictly increasing",
            "Resource bersama menolak setiap request dengan incomingToken <= lastSeenToken",
            "Keamanan ditegakkan pada gerbang data, bukan pada persepsi internal klien",
        ], "type": "bullet"},
    ],
    mental_model="A fencing token converts temporal leases into monotonic ordering enforced by the storage gate.",
    source_path="labs/30-leader-election/content/01-content-brief.md:6",
    slide_num="3 / 24"
)

# Slide 4: Core Triad Architecture
add_slide(
    categories=["ARSITEKTUR", "KOMPONEN"],
    title="Tiga Pilar Pola Fencing Tokens",
    subtitle="Pemisahan tanggung jawab antara koordinasi, kandidat, dan proteksi resource",
    content_items=[
        {"title": "1. LINEARIZABLE COORDINATOR", "items": [
            "Menyimpan mapping key-lease dan counter revision monotonik global",
            "Memberikan hak sewa berdurasi TTL dan menerbitkan fencing token baru",
            "Menolak klaim tumpang tindih selama lease lama masih aktif (ErrLeaseHeld)",
        ], "type": "bullet"},
        {"title": "2. ELECTION CANDIDATE (NODE)", "items": [
            "Menjalankan election loop reaktif: Acquire saat follower, Renew saat leader",
            "Menyimpan status internal (FOLLOWER vs LEADER) dan salinan lease lokal",
            "Meneruskan fencing token pada setiap payload penulisan ke shared resource",
        ], "type": "bullet"},
        {"title": "3. FENCED STORAGE ENGINE", "items": [
            "Menyimpan state data bersama dan melacak variabel lastSeenToken",
            "Menjalankan validasi check-and-set atomik sebelum mutasi diizinkan",
            "Menolak paket transaksi basi berstatus stale dengan ErrStaleFencingToken",
        ], "type": "bullet"},
    ],
    mental_model="Coordination issues tokens, candidates campaign and execute, storage enforces strict ordering.",
    source_path="labs/30-leader-election/engineering/01-design.md:10",
    slide_num="4 / 24"
)

# Slide 5: System Architecture Diagram
add_slide(
    categories=["DESAIN SISTEM", "ALUR TOPOLOGI"],
    title="Topologi Sistem & Alur Interaksi Komponen",
    subtitle="Hubungan relasional antara Node A, Node B, Coordinator, dan Storage",
    content_items=[
        {"title": "DIAGRAM TOPOLOGI KOMPONEN", "items": [
            "+-------------------+       +-------------------+",
            "| Node A (Leader)   |       | Node B (Follower) |",
            "+---------+---------+       +---------+---------+",
            "          |                           |          ",
            "          | Acquire / Renew / Release |          ",
            "          v                           v          ",
            "+-----------------------------------------------+",
            "|            Coordinator (In-Memory)            |",
            "| - leases: map[string]*Lease                   |",
            "| - revision: int64 (monotonic fencing counter) |",
            "| - sync.Mutex (linearizable operations)        |",
            "+-----------------------+-----------------------+",
            "                        |                        ",
            "                        | Write(author, token, v)",
            "                        v                        ",
            "+-----------------------------------------------+",
            "|          FencedStorage (Shared Resource)      |",
            "| - lastSeenToken: int64                        |",
            "| - Check-and-set: incomingToken > lastSeenToken|",
            "+-----------------------------------------------+",
        ], "type": "diagram"},
        {"title": "POIN INTERAKSI UTAMA", "items": [
            "Komunikasi antara Node dan Coordinator terisolasi dari jalur penulisan data",
            "Node tidak pernah mengkoordinasikan konsistensi data secara langsung antar rekan (peer)",
            "Shared resource bertindak sebagai benteng pertahanan terakhir konsistensi",
        ], "type": "bullet"},
        {"title": "SIFAT IN-MEMORY LAB", "items": [
            "Seluruh komponen berjalan thread-safe dalam runtime Go murni tanpa dependensi daemon eksternal",
        ], "type": "bullet"},
    ],
    mental_model="Storage is the ultimate arbiter; nodes never negotiate consistency directly with peers.",
    source_path="labs/30-leader-election/content/04-diagrams.md:5",
    slide_num="5 / 24"
)

# Slide 6: Lifecycle State Machine
add_slide(
    categories=["STATE MACHINE", "LIFECYCLE"],
    title="Siklus Hidup Kandidat: Follower & Leader",
    subtitle="Transisi state deterministik berbasis heartbeat renewal dan akuisisi",
    content_items=[
        {"title": "KONDISI FOLLOWER", "items": [
            "State awal setiap node kandidat yang bergabung dalam kluster",
            "Pada setiap tick interval renewal, memanggil coord.Acquire(key, id, ttl)",
            "Jika sukses memperoleh lease, beralih ke state LEADER dan menyimpan token",
            "Jika gagal (ErrLeaseHeld), tetap berstatus FOLLOWER dan menunggu tick berikutnya",
        ], "type": "bullet"},
        {"title": "KONDISI LEADER", "items": [
            "Bertanggung jawab mengeksekusi operasi bisnis dan penulisan berfencing",
            "Secara periodik memanggil coord.Renew(key, id, token, ttl) sebelum batas waktu habis",
            "Jika Renew berhasil, batas waktu kadaluwarsa diperpanjang dengan token tetap",
            "Jika Renew gagal (ErrLeaseExpired / ErrNotLeaseOwner), langsung revert ke FOLLOWER",
        ], "type": "bullet"},
        {"title": "SIMULASI PAUSE & PARTISI", "items": [
            "Panggilan SimulatePause() menahan loop election tanpa mengirim heartbeat",
            "Node tidak sempat menurunkan statusnya sendiri karena thread dalam keadaan terhenti",
        ], "type": "bullet"},
    ],
    mental_model="Heartbeat renews keep leaders active; missed renewals force an immediate step-down to follower.",
    source_path="labs/30-leader-election/internal/candidate/candidate.go:68",
    slide_num="6 / 24"
)

# Slide 7: Coordinator & Monotonic Revision Implementation
add_slide(
    categories=["BUKTI KODE", "COORDINATOR"],
    title="Penerbitan Token Monotonik di Koordinator",
    subtitle="Struktur Lease dan increment atomik revision di bawah sync.Mutex",
    content_items=[
        {"title": "STRUKTUR DATA LEASE", "items": [
            "type Lease struct {",
            "    Key          string",
            "    HolderID     string",
            "    FencingToken int64",
            "    ExpiresAt    time.Time",
            "    TTL          time.Duration",
            "}",
        ], "type": "code"},
        {"title": "LOGIKA PENERBITAN TOKEN", "items": [
            "c.revision++",
            "lease := &Lease{",
            "    Key:          key,",
            "    HolderID:     holderID,",
            "    FencingToken: c.revision,",
            "    ExpiresAt:    now.Add(ttl),",
            "    TTL:          ttl,",
            "}",
            "c.leases[key] = lease",
        ], "type": "code"},
        {"title": "JAMINAN KONSISTENSI", "items": [
            "Operasi revision++ dilindungi sync.Mutex koordinator secara ketat",
            "Setiap lease baru dijamin memiliki angka FencingToken yang lebih tinggi",
            "Salinan nilainya didistribusikan ke pemegang baru sebagai bukti kewenangan",
        ], "type": "bullet"},
    ],
    mental_model="Every newly granted lease atomically increments a global revision counter.",
    source_path="labs/30-leader-election/internal/coordinator/coordinator.go:21",
    slide_num="7 / 24"
)

# Slide 8: Lease Acquisition & Contention Guard
add_slide(
    categories=["BUKTI KODE", "KONTENSI"],
    title="Pemeriksaan Masa Berlaku & Kepemilikan Sewa",
    subtitle="Mekanisme Acquire mencegah klaim ganda selama lease masih aktif",
    content_items=[
        {"title": "KODE KOORDINATOR ACQUIRE", "items": [
            "now := time.Now().Add(c.clockOffset)",
            "existing, found := c.leases[key]",
            "if found && now.Before(existing.ExpiresAt) {",
            "    if existing.HolderID == holderID {",
            "        existing.ExpiresAt = now.Add(ttl)",
            "        return copyLease(existing), nil",
            "    }",
            "    return nil, ErrLeaseHeld",
            "}",
        ], "type": "code"},
        {"title": "IDEMPOTENSI HOLDER SAMA", "items": [
            "Jika pemanggil adalah pemilik yang sama, expiry diperbarui tanpa menaikkan revision",
            "Hal ini menjaga stabilitas nomor token selama pemegang kepemimpinan tidak berganti",
        ], "type": "bullet"},
        {"title": "PROTEKSI ERR_LEASE_HELD", "items": [
            "Node kompetitor yang mencoba mengambil alih saat lease aktif akan ditolak mentah-mentah",
            "Menjamin kepemimpinan tunggal pada kondisi normal tanpa partisi",
        ], "type": "bullet"},
    ],
    mental_model="Active leases block foreign acquisitions but permit idempotent extensions by the rightful holder.",
    source_path="labs/30-leader-election/internal/coordinator/coordinator.go:37",
    slide_num="8 / 24"
)

# Slide 9: Lease Renewal & Expiry Mechanics
add_slide(
    categories=["BUKTI KODE", "HEARTBEAT"],
    title="Mekanisme Perpanjangan Sewa (Renew)",
    subtitle="Validasi ganda: masa berlaku waktu dan pencocokan kepemilikan token",
    content_items=[
        {"title": "IMPLEMENTASI KOORDINATOR RENEW", "items": [
            "if now.After(existing.ExpiresAt) || now.Equal(existing.ExpiresAt) {",
            "    delete(c.leases, key)",
            "    return nil, ErrLeaseExpired",
            "}",
            "if existing.HolderID != holderID || existing.FencingToken != token {",
            "    return nil, ErrNotLeaseOwner",
            "}",
            "existing.ExpiresAt = now.Add(ttl)",
        ], "type": "code"},
        {"title": "CLEANUP OTOMATIS SAAT KADALUWARSA", "items": [
            "Jika now >= ExpiresAt, lease langsung dihapus dari map dan dikembalikan ErrLeaseExpired",
            "Kandidat yang telat memperpanjang kehilangan status sewa seketika",
        ], "type": "bullet"},
        {"title": "PENCEGAHAN ZOMBIE RENEWAL", "items": [
            "Validasi existing.FencingToken != token menolak usaha renewal dari leader purba",
            "Mencegah node yang terlambat bangun merebut kembali lease yang sudah dialihkan",
        ], "type": "bullet"},
    ],
    mental_model="Renewals require both an unexpired deadline and an exact match of the holder and token.",
    source_path="labs/30-leader-election/internal/coordinator/coordinator.go:76",
    slide_num="9 / 24"
)

# Slide 10: Fenced Storage Check-and-Set
add_slide(
    categories=["BUKTI KODE", "STORAGE GATE"],
    title="Shared Storage Check-and-Set Gate",
    subtitle="Benteng pertahanan penolakan mutasi stale berbasis perbandingan token",
    content_items=[
        {"title": "VALIDASI STRICT CHECK-AND-SET", "items": [
            "func (s *FencedStorage) Write(author string, token int64, value string) error {",
            "    s.mu.Lock(); defer s.mu.Unlock()",
            "    if token <= s.lastSeenToken {",
            "        return fmt.Errorf(\"%w: token %d <= last seen %d (author: %s)\",",
            "            ErrStaleFencingToken, token, s.lastSeenToken, author)",
            "    }",
            "    s.lastSeenToken = token",
            "    s.records = append(s.records, Record{...})",
            "    return nil",
            "}",
        ], "type": "code"},
        {"title": "SYARAT KETAT: STRICTLY GREATER", "items": [
            "Perhatikan tanda <= : Token yang sama persis (duplikat) wajib ditolak",
            "Setiap mutasi yang sah harus memiliki nomor token yang lebih tinggi secara absolut",
        ], "type": "bullet"},
        {"title": "ATOMISITAS EKSEKUSI", "items": [
            "Pemeriksaan token dan append record berada di dalam satu blok s.mu.Lock()",
            "Menghilangkan window of vulnerability akibat race condition antar worker",
        ], "type": "bullet"},
    ],
    mental_model="Reject writes if incomingToken <= lastSeenToken; storage is the final monotonic gatekeeper.",
    source_path="labs/30-leader-election/internal/storage/storage.go:27",
    slide_num="10 / 24"
)

# Slide 11: Candidate Election Loop & GC Simulation
add_slide(
    categories=["BUKTI KODE", "ELECTION LOOP"],
    title="Loop Pemilihan Kandidat & Simulasi Pause",
    subtitle="Bagaimana kandidat bereaksi terhadap interval waktu dan penundaan eksekusi",
    content_items=[
        {"title": "KODE LOOP ELECTION NODE", "items": [
            "for {",
            "    select {",
            "    case <-ctx.Done(): return",
            "    case now := <-ticker.C:",
            "        n.mu.Lock()",
            "        if now.Before(n.pauseUntil) { n.mu.Unlock(); continue }",
            "        if n.state == StateLeader {",
            "            renewed, err := n.coord.Renew(...)",
            "            if err != nil { n.state = StateFollower; n.currentLease = nil }",
            "        } else {",
            "            acquired, err := n.coord.Acquire(...)",
            "            if err == nil { n.state = StateLeader; n.currentLease = acquired }",
            "        }",
            "        n.mu.Unlock()",
            "    }",
            "}",
        ], "type": "code"},
        {"title": "SIMULATE PAUSE MONOTONIK", "items": [
            "n.pauseUntil = time.Now().Add(duration) menahan loop tanpa intervensi OS sinyal",
            "Mensimulasikan full stop-the-world GC pause atau delay partisi kabel jaringan",
        ], "type": "bullet"},
    ],
    mental_model="Periodic election loops react to lease deadlines while pause markers freeze execution safely.",
    source_path="labs/30-leader-election/internal/candidate/candidate.go:106",
    slide_num="11 / 24"
)

# Slide 12: Split-Brain Timeline Sequence
add_slide(
    categories=["ALUR WAKTU", "SKENARIO"],
    title="Kronologi Kegagalan & Penolakan Stale Write",
    subtitle="Rekonstruksi detik demi detik pencegahan korupsi data di shared resource",
    content_items=[
        {"title": "URUTAN WAKTU EKSEKUSI", "items": [
            "t0: Node A acquire lease -> Token 1 diterbitkan -> Node A jadi Leader",
            "t1: Node A menulis 'record-1' dengan Token 1 -> Storage menerima (lastSeen=1)",
            "t2: Node A mengalami Pause 350ms (TTL=200ms) -> Heartbeat berhenti terkirim",
            "t3: TTL habis di t0+200ms -> Koordinator menandai lease Node A kedaluwarsa",
            "t4: Node B acquire lease -> Token 2 diterbitkan -> Node B jadi Leader baru",
            "t5: Node B menulis 'record-2' dengan Token 2 -> Storage menerima (lastSeen=2)",
            "t6: Node A terbangun dari pause di t2+350ms -> Mengira masih memegang token 1",
            "t7: Node A menulis dengan Token 1 -> Storage REJECT: Token 1 <= lastSeen 2!",
        ], "type": "diagram"},
        {"title": "KONDISI DUAL-LEADER TRANSIENT", "items": [
            "Pada t6, Node A dan Node B sempat sama-sama meyakini dirinya adalah LEADER",
            "Sistem terdistribusi asinkron tidak dapat mencegah persepsi ganda sementara ini",
            "Namun integritas data tetap 100% terjaga berkat rejection di FencedStorage",
        ], "type": "bullet"},
    ],
    mental_model="Transient dual-leader states are inevitable; fencing tokens render stale leaders powerless at the gate.",
    source_path="labs/30-leader-election/content/04-diagrams.md:46",
    slide_num="12 / 24"
)

# Slide 13: Demo Walkthrough: Election & Normal Write
add_slide(
    categories=["DEMO RUNNER", "FASE 1"],
    title="Eksekusi Demo: Akuisisi Kepemimpinan Awal",
    subtitle="Node-A memenangkan kampanye awal dan melakukan penulisan berfencing sah",
    content_items=[
        {"title": "LOG TERMINAL DEMO FASE 1", "items": [
            "[DEMO] Memulai 2 kandidat bersaing untuk lease 'app-primary' (TTL: 200ms, Renew: 50ms)...",
            "[STATUS] Node-A: LEADER | Node-B: FOLLOWER",
            "[LEADER] Node-A memenangkan pemilihan, token: 1",
            "[WRITE] Node-A menulis dengan token 1...",
            "[SUCCESS] Storage menerima tulisan dari Node-A dengan token 1",
            "[PAUSE] Mensimulasikan stop-the-world GC pause 350ms pada Node-A...",
        ], "type": "code"},
        {"title": "OBSERVASI PERILAKU", "items": [
            "Node-A berhasil memperoleh lease pertama dan menerima FencingToken bernilai 1",
            "Storage mengizinkan penulisan record pertama karena token 1 > lastSeenToken 0",
            "Nilai lastSeenToken pada FencedStorage diperbarui menjadi 1",
        ], "type": "bullet"},
        {"title": "KONDISI STANDBY NODE-B", "items": [
            "Node-B tetap berstatus FOLLOWER karena Acquire ditolak dengan ErrLeaseHeld",
        ], "type": "bullet"},
    ],
    mental_model="The winning candidate obtains token 1 and successfully writes to storage initialized at 0.",
    source_path="labs/30-leader-election/cmd/demo/main.go:37",
    slide_num="13 / 24"
)

# Slide 14: Demo Walkthrough: Failover & Stale Rejection
add_slide(
    categories=["DEMO RUNNER", "FASE 2"],
    title="Eksekusi Demo: Failover & Penolakan Mutasi Stale",
    subtitle="Node-B mengambil alih dan request zombie Node-A digagalkan oleh storage",
    content_items=[
        {"title": "LOG TERMINAL DEMO FASE 2", "items": [
            "[FAILOVER] Node-B mengambil alih kepemimpinan, token baru: 2",
            "[STATUS] Node-A: LEADER (pause) | Node-B: LEADER",
            "[WRITE] Node-B menulis dengan token 2...",
            "[SUCCESS] Storage menerima tulisan dari Node-B dengan token 2",
            "[RESUME] Node-A terbangun dari pause dan mencoba menulis dengan token lama (1)...",
            "[FENCING DEFENSE] Storage MENOLAK tulisan dari Node-A!",
            "  Error: stale fencing token: token 1 <= last seen 2 (author: Node-A)",
            "[VERIFIKASI] Storage records berisi 2 entri yang sah, split-brain berhasil dicegah!",
        ], "type": "code"},
        {"title": "BUKTI PENOLAKAN DETERMINISTIK", "items": [
            "Node-B sukses memvalidasi token 2 di koordinator setelah lease Node-A expired",
            "Tulisan Node-A pasca-bangun ditolak secara presisi oleh ErrStaleFencingToken",
            "Tidak ada data corrupt atau overwriting yang terjadi di dalam penyimpanan bersama",
        ], "type": "bullet"},
    ],
    mental_model="Failover elevates candidate B to token 2, causing the awakened leader A's write to be denied.",
    source_path="labs/30-leader-election/cmd/demo/main.go:68",
    slide_num="14 / 24"
)

# Slide 15: Quantitative Test Verification
add_slide(
    categories=["VERIFIKASI TEST", "EVIDENCE"],
    title="Bukti Pengujian Otomatis: 4/4 Test Suites Lulus",
    subtitle="Verifikasi formal akuisisi sewa, penolakan token, failover, dan balapan paralel",
    content_items=[
        {"title": "LOG EKSEKUSI TEST & RACE DETECTOR", "items": [
            "=== RUN   TestCoordinatorLeaseAcquisitionAndRenewal",
            "--- PASS: TestCoordinatorLeaseAcquisitionAndRenewal (0.25s)",
            "=== RUN   TestFencedStorageRejectsStaleTokens",
            "--- PASS: TestFencedStorageRejectsStaleTokens (0.00s)",
            "=== RUN   TestLeaderElectionFailoverAndSplitBrainDefense",
            "--- PASS: TestLeaderElectionFailoverAndSplitBrainDefense (0.28s)",
            "=== RUN   TestConcurrentElectionRace",
            "--- PASS: TestConcurrentElectionRace (0.10s)",
            "PASS: 4 passed in 5 packages (total time: 2.000s, zero race warnings)",
        ], "type": "code"},
        {"title": "CAKUPAN PENGUJIAN UTAMA", "items": [
            "TestCoordinatorLeaseAcquisitionAndRenewal: Memvalidasi renewal dan kenaikan revision",
            "TestFencedStorageRejectsStaleTokens: Menolak pengulangan token dan token mundur",
            "TestLeaderElectionFailoverAndSplitBrainDefense: Skenario full lifecycle pause & failover",
            "TestConcurrentElectionRace: 5 kandidat bersaing simultan, tepat 1 leader terpilih",
        ], "type": "bullet"},
        {"title": "HASIL ZERO RACE WARNINGS", "items": [
            "Pengujian dijalankan dengan flag go test -race ./... tanpa terdeteksi data race apa pun",
        ], "type": "bullet"},
    ],
    mental_model="Unit, failover, and race tests verify the single-leader invariant and monotonic fencing.",
    source_path="labs/30-leader-election/tests/election_test.go:13",
    slide_num="15 / 24"
)

# Slide 16: The Redlock Controversy Explained
add_slide(
    categories=["DEBAT INDUSTRI", "ANALISIS"],
    title="Debat Redlock: Martin Kleppmann vs antirez",
    subtitle="Mengapa distributed lock tanpa fencing token tidak memadai untuk korektivitas",
    content_items=[
        {"title": "KRITIK MARTIN KLEPPMANN (2016)", "items": [
            "Lock berbasis waktu di Redis (Redlock) mengasumsikan model waktu sinkron semu",
            "GC pause, swap delay, atau lompatan jam (NTP shift) dapat membuat lock kadaluwarsa",
            "Klien yang tidak sadar lock-nya habis akan menulis ke storage dan merusak data",
            "Kesimpulan: Lock service saja tidak cukup; storage wajib melakukan fencing token",
        ], "type": "bullet"},
        {"title": "TANGGAPAN SALVATORE SANFILIPPO (ANTIREZ)", "items": [
            "Redlock dirancang praktis untuk menjamin mutual exclusion probabilistik tinggi",
            "Jika resource backend mendukung fencing token, algoritma koordinasi apa pun jadi aman",
            "Redlock tetap sangat efisien untuk tugas yang tidak membutuhkan strict correctness",
        ], "type": "bullet"},
        {"title": "KONSENSUS ARSITEKTUR MODERN", "items": [
            "Gunakan lock sederhana (Redis) untuk optimasi efisiensi (pencegahan duplikasi job)",
            "Gunakan konsensus terbukti (etcd/ZooKeeper) + fencing token untuk data kritis",
        ], "type": "bullet"},
    ],
    mental_model="Distributed locks only optimize efficiency; safety requires monotonic fencing at the resource layer.",
    source_path="labs/30-leader-election/content/02-master-draft.md:200",
    slide_num="16 / 24"
)

# Slide 17: Industry Implementations Comparison
add_slide(
    categories=["TABEL EVALUASI", "PERBANDINGAN"],
    title="Perbandingan Implementasi Sistem Terdistribusi",
    subtitle="Bagaimana etcd, ZooKeeper, Consul, dan Redis mengelola pemilihan leader",
    content_items=[
        {"title": "TABEL ANALISIS 4 TEKNOLOGI KOORDINASI", "items": [
            "Sistem    | Mekanisme Pemilihan    | Fencing Token Asli   | Sifat Konsensus",
            "----------|------------------------|----------------------|-----------------",
            "etcd      | Lease + Campaign Keys  | Revision (int64)     | Raft (Strict)",
            "ZooKeeper | Ephemeral Sequentials  | zxid / Node Version  | ZAB (Strict)",
            "Consul    | Session TTL + Locks    | ModifyIndex (uint64) | Raft (Strict)",
            "Redis     | SET NX EX (Redlock)    | Tidak ada bawaan     | Asynchronous",
        ], "type": "table"},
        {"title": "ETCD REVISION MODEL", "items": [
            "Revision pada etcd adalah counter 64-bit global yang monotonik bertambah di setiap mutasi",
            "Kandidat dapat menggunakan revision key saat lease diberikan sebagai fencing token",
        ], "type": "bullet"},
        {"title": "ZOOKEEPER ZXID PATTERN", "items": [
            "ZooKeeper transaction ID (zxid) mencerminkan tatanan kausal total dari seluruh proposal",
        ], "type": "bullet"},
    ],
    mental_model="etcd provides revisions, ZooKeeper provides zxids; both offer native monotonic fencing tokens.",
    source_path="labs/30-leader-election/content/02-master-draft.md:215",
    slide_num="17 / 24"
)

# Slide 18: Monotonic Clocks vs Wall Clocks
add_slide(
    categories=["RELIABILITAS", "CLOCK SAFETY"],
    title="Pentingnya Jam Monotonik: Hindari NTP Steps",
    subtitle="Mengapa wall-clock dilarang keras untuk pengukuran durasi lease",
    content_items=[
        {"title": "BAHAYA WALL-CLOCK (REAL TIME)", "items": [
            "Wall-clock merefleksikan waktu kalender manusia (UTC / local time)",
            "Sinkronisasi Network Time Protocol (NTP) dapat menyebabkan waktu melompat maju atau mundur",
            "Lompatan mundur membuat lease tampak aktif selamanya; lompatan maju mematikan lease prematur",
        ], "type": "bullet"},
        {"title": "MONOTONIC CLOCK DI RUNTIME GO", "items": [
            "Go runtime secara otomatis menyertakan monotonic reading pada panggilan time.Now()",
            "Operasi durasi seperti time.Since(t) atau now.Before(expiry) membaca os monotonic timer",
            "Tidak terpengaruh oleh penyesuaian waktu sistem operasi maupun leap second",
        ], "type": "bullet"},
        {"title": "PEDOMAN IMPLEMENTASI", "items": [
            "Gunakan Monotonic Clock untuk internal expiry, renew timeouts, dan heartbeat loop",
            "Gunakan Wall-Clock HANYA untuk pencatatan log audit manusia dan timestamp tampilan",
        ], "type": "bullet"},
    ],
    mental_model="Always calculate lease deadlines with monotonic timers; wall-clock shifts corrupt time checks.",
    source_path="labs/30-leader-election/content/02-master-draft.md:73",
    slide_num="18 / 24"
)

# Slide 19: Preventing Herd Effects & Failure Thundering
add_slide(
    categories=["OPTIMISASI", "SKALABILITAS"],
    title="Mencegah Herd Effect Saat Failover",
    subtitle="Strategi antrian teratur agar ratusan follower tidak membanjiri koordinator",
    content_items=[
        {"title": "PROBLEMATIKA THUNDERING HERD", "items": [
            "Jika 100 node serentak melakukan polling Acquire saat leader mati, koordinator overload",
            "Lonjakan network traffic dan lock contention menurunkan responsivitas kluster",
        ], "type": "bullet"},
        {"title": "STRATEGI 1: DIRECT PREDECESSOR WATCH", "items": [
            "Pola klasik ZooKeeper: Setiap kandidat membuat ephemeral sequential node",
            "Kandidat nomor N HANYA memantau (watch) event penghapusan pada node nomor N-1",
            "Saat leader gugur, hanya node berikutnya dalam antrean yang terbangun dan berpromosi",
        ], "type": "bullet"},
        {"title": "STRATEGI 2: RANDOMIZED EXPONENTIAL BACKOFF", "items": [
            "Jika menggunakan polling, tambahkan jitter acak pada durasi tunggu kandidat",
            "Mendistribusikan waktu kedatangan request akuisisi secara merata di timeline",
        ], "type": "bullet"},
    ],
    mental_model="Watch immediate predecessors or add randomized jitter to eliminate thundering herd storms.",
    source_path="labs/30-leader-election/content/02-master-draft.md:183",
    slide_num="19 / 24"
)

# Slide 20: Observability & Rejection Monitoring
add_slide(
    categories=["OBSERVABILITY", "MONITORING"],
    title="Observabilitas & Metrik Penolakan Fencing",
    subtitle="Menjadikan ErrStaleFencingToken sebagai indikator kesehatan sistem",
    content_items=[
        {"title": "METRIK UTAMA YANG HARUS DIUKUR", "items": [
            "fencing_token_rejections_total: Counter tulisan yang ditolak oleh resource gate",
            "leader_lease_renew_duration_seconds: Latensi round-trip pembaruan heartbeat",
            "leader_active_duration_seconds: Berapa lama seorang leader mampu mempertahankan jabatan",
            "candidate_state: Gauge status node (0=Follower, 1=Candidate, 2=Leader)",
        ], "type": "bullet"},
        {"title": "SIGNIFIKANSI REJECTION ALERT", "items": [
            "Satu kali penolakan stale write adalah bukti keberhasilan mekanisme pertahanan split-brain",
            "Namun lonjakan penolakan menandakan adanya degradasi performa: GC pause parah atau jaringan flap",
            "Gunakan rejection rate sebagai Early Warning System sebelum kegagalan hardware meluas",
        ], "type": "bullet"},
        {"title": "AUDIT TRAIL LOGGING", "items": [
            "Catat setiap kejadian penolakan dengan menyertakan author ID, token lama, dan token aktif",
        ], "type": "bullet"},
    ],
    mental_model="A fencing rejection proves safety worked, but frequent rejections signal severe latency spikes.",
    source_path="labs/30-leader-election/content/02-master-draft.md:185",
    slide_num="20 / 24"
)

# Slide 21: Production Readiness Checklist
add_slide(
    categories=["PANDUAN OPERASIONAL", "CHECKLIST"],
    title="Checklist Kesiapan Produksi Leader Election",
    subtitle="9 poin verifikasi arsitektur sebelum menggelar sistem koordinasi ke live cluster",
    content_items=[
        {"title": "CHECKLIST PENERAPAN KOORDINASI", "items": [
            "[ ] TTL lease dikalibrasi berdasarkan p99 durasi GC pause dan network RTT",
            "[ ] Margin interval renewal diatur minimal 3x lebih cepat dari durasi TTL",
            "[ ] Koordinator menerbitkan token monotonik strictly increasing di setiap lease baru",
            "[ ] Resource penyimpanan menerapkan check-and-set atomik (token <= lastSeen -> REJECT)",
            "[ ] Semua logika durasi dan timeout menggunakan Monotonic Clock",
            "[ ] Leader melepaskan sewa secara sukarela (Release) saat proses shutdown anggun",
            "[ ] Metrik dan alerting aktif untuk pemantauan fencing_token_rejections",
            "[ ] Penanganan herd effect (watch predecessor / jitter backoff) telah diuji",
            "[ ] Lulus pengujian konkurensi race detector (go test -race)",
        ], "type": "bullet"},
        {"title": "ATURAN UTAMA", "items": [
            "Jangan pernah mempercayai klien untuk menentukan apakah dirinya masih sah sebagai leader",
            "Pastikan shared resource memiliki otoritas mutlak menolak tulisan basi",
        ], "type": "bullet"},
    ],
    mental_model="Calibrate TTL margins, enforce monotonic fencing at the resource, and verify under race tests.",
    source_path="labs/30-leader-election/content/02-master-draft.md:225",
    slide_num="21 / 24"
)

# Slide 22: Scope & Boundaries of Lab 30
add_slide(
    categories=["BATAS LAB", "TRANSPARANSI"],
    title="Batasan & Lingkup Implementasi Lab 30",
    subtitle="Membedakan simplifikasi edukasi dengan kebutuhan enterprise sesungguhnya",
    content_items=[
        {"title": "IN-MEMORY LINEARIZABILITY", "items": [
            "Koordinator lab menggunakan map in-memory yang dilindungi Go sync.Mutex",
            "Membuktikan kebenaran algoritma fencing tanpa perlu menjalankan cluster etcd multi-node",
            "Tidak menyimulasikan crash koordinator atau kehilangan state akibat proses restart",
        ], "type": "bullet"},
        {"title": "SIMULASI PAUSE LOKAL", "items": [
            "Penundaan eksekusi dilakukan via field pauseUntil bukan stop-the-world signal OS nyata",
            "Tidak memodelkan Byzantine faults atau manipulasi jam sistem yang disengaja",
            "Parameter lab (TTL 200ms, Renew 50ms) dioptimalkan untuk kecepatan eksekusi demo",
        ], "type": "bullet"},
        {"title": "NILAI INTI YANG TERBUKTI", "items": [
            "Invarian single-leader terbukti deterministik pada pengujian konkuren 5 node",
            "Pertahanan split-brain melalui penolakan token terbukti 100% menangkal korupsi data",
        ], "type": "bullet"},
    ],
    mental_model="This lab proves algorithmic correctness in-memory; production relies on consensus engines like etcd.",
    source_path="labs/30-leader-election/README.md:1",
    slide_num="22 / 24"
)

# Slide 23: Core Distributed Systems Principles
add_slide(
    categories=["PRINSIP INTI", "RINGKASAN"],
    title="Tiga Prinsip Inti Koordinasi Terdistribusi",
    subtitle="Pedoman fundamental bagi software engineer dalam mendesain arsitektur cluster",
    content_items=[
        {"title": "PRINSIP 1: WAKTU ADALAH ILUSI DALAM DISTRIBUTED SYSTEMS", "items": [
            "Jangan pernah mengandalkan jam fisik atau durasi absolut untuk keamanan data mutlak",
            "Lease temporal hanya memberikan mutual exclusion probabilistik",
            "Ubah batas waktu menjadi tatanan logis monotonik yang dapat diverifikasi",
        ], "type": "bullet"},
        {"title": "PRINSIP 2: VALIDASI HARUS BERADA DI RESOURCE LAYER", "items": [
            "Klien yang mengalami pause tidak pernah tahu bahwa dirinya sudah kedaluwarsa",
            "Klien zombie akan selalu merasa dirinya adalah leader yang sah",
            "Hanya lapisan penyimpanan bersama yang dapat menjadi gerbang validasi akhir",
        ], "type": "bullet"},
        {"title": "PRINSIP 3: STRICT CHECK-AND-SET ATOMIK", "items": [
            "Terapkan aturan perbandingan tegas: incomingToken > lastSeenToken",
            "Tolak token yang sama (duplikat) maupun token masa lalu yang lebih rendah",
            "Kunci atomisitas check-and-set menjamin integritas data dari race condition",
        ], "type": "bullet"},
    ],
    mental_model="Time is unreliable; safety demands monotonic ordering validated atomically at the shared resource.",
    source_path="labs/30-leader-election/content/05-key-takeaways.md:1",
    slide_num="23 / 24"
)

# Slide 24: Closing & Resources (Dark Navy)
add_slide(
    categories=["RINGKASAN"],
    title="Ringkasan Lab 30: Leader Election with Fencing Tokens",
    subtitle=f"Commit: {COMMIT_SHA[:8]} | 4/4 Test Suites PASS | Repository: lukman-ss/software-engineering-lab",
    content_items=[
        {"title": "HASIL VERIFIKASI UTAMA", "items": [
            "Single-leader invariant terbukti deterministik di bawah balapan 5 kandidat simultan",
            "Failover otomatis sukses mengambil alih kepemimpinan saat heartbeat leader terhenti",
            "Shared storage menolak mutasi stale dari leader lama dengan ErrStaleFencingToken",
            "4 unit dan race detector test suites lulus sempurna (go test -v -race ./...)",
        ], "type": "bullet"},
        {"title": "MENTAL MODEL UTAMA", "items": [
            "A lease grants time-bounded authority; a fencing token converts that temporal authority",
            "into a monotonic ordering that shared storage enforces — rejecting stale zombie writes.",
        ], "type": "bullet"},
        {"title": "RESOURCES & AUDIT", "items": [
            f"Repository: {REPO_URL}",
            f"Commit SHA: {COMMIT_SHA}",
            "Jalankan Demo: go run ./cmd/demo (di folder lab 30)",
            "Jalankan Test: go test -v -race ./... (di folder lab 30)",
        ], "type": "bullet"},
    ],
    mental_model="Leases grant temporary authority; fencing tokens enforce immutable safety at the storage gate.",
    source_path="labs/30-leader-election/README.md",
    slide_num=""
)

if __name__ == "__main__":
    print(f"Total slides defined: {len(SLIDES)}")
    for i, s in enumerate(SLIDES):
        print(f"Slide {i+1}: {s['title']} ({len(s['content_items'])} content items)")
