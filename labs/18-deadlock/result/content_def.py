#!/usr/bin/env python3
"""
Publication Builder Definition for Software Engineering Lab 18: Deadlock in Concurrent Systems.
Adheres strictly to:
- Cover hero image requirement: Mandatory raster image at COVER_IMAGE_PATH
- Slide composition rules: <=45 words prose, <=30 words on code slides, dominant visual anchors
- Quality bar: strict technical boundary, limiting phrases on >=12 slides, zero banned filler words
- Concrete 24-slide LinkedIn Carousel structure
"""

OUTPUT_DIR = "labs/18-deadlock/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/18-deadlock"
COMMIT_SHA = "e76d13f"
LAB_NO = "18"
SLUG = "deadlock"

TITLE_ID = "Deadlock di Sistem Konkuren"
SUBTITLE_ID = "Mekanisme Siklus Tunggu, Penanganan Database, dan Strategi Aplikasi di Go"
TITLE_EN = "Deadlock in Concurrent Systems"
SUBTITLE_EN = "Circular Wait Mechanics, Database Engine Aborts, and Application Strategies in Go"

AUTHOR = "Lukman / lukman-ss"
TOPIC_TAGS = "Deadlock • Concurrency • Database Locks • Lock Ordering • Retry Pattern • Go"
COVER_IMAGE_PATH = "labs/18-deadlock/result/image-cover/cover_art.jpg"

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

# Slide 01: Cover (Layout: cover)
add_slide(
    categories=["COVER"],
    title=TITLE_ID,
    subtitle=SUBTITLE_ID,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "Retro Terminal Visualizing Wait-For Graph, Lock Ordering & Abort Mechanics",
            "Simulasi Konkurensi Mutasi Saldo Bank dengan Goroutines dan Context Timeout di Go",
            f"Penulis: {AUTHOR}"
        ]}
    ],
    mental_model="Deadlock is an inescapable circular wait resolved only by external engine aborts or strict lock ordering.",
    source_path="",
    slide_num="01 / 24",
    layout="cover"
)

# Slide 02: Problem Nyata (Layout: cards)
# Limiting phrase 1: "bukan berarti"
add_slide(
    categories=["PROBLEM", "CONCURRENCY FAILURE"],
    title="Lonjakan Transaksi Konkuren Memicu Kemacetan Sistemik",
    subtitle="Ketika dua proses mengunci data yang sama dengan urutan berlawanan, sistem terhenti total.",
    content_items=[
        {"type": "cards", "title": "DAMPAK DEADLOCK PADA SISTEM TRANSAKSIONAL", "items": [
            {"title": "Peningkatan Latensi", "desc": "Transaksi tertahan tanpa batas waktu hingga mekanisme deteksi siklus basis data terpicu."},
            {"title": "Gugurnya Transaksi", "desc": "Database membatalkan satu transaksi sepihak; klien menerima galat pembatalan mendadak."},
            {"title": "Bukan Berarti Rusak", "desc": "Deadlock bukan berarti basis data rusak, melainkan akibat pola urutan akuisisi lock aplikasi yang acak."}
        ]}
    ],
    narrative="Ketika volume transaksi tinggi mengakses baris data secara bersamaan, urutan akses yang tidak teratur memicu kebuntuan fatal. Munculnya galat pembatalan transaksi bukan berarti database mengalami kerusakan internal, melainkan sinyal bahwa lapisan aplikasi membiarkan pola penguncian saling silang.",
    narrative_en="When high-volume transactions access shared data concurrently, uncoordinated access sequences cause systemic freezes. Encountering aborted transaction errors does not imply database malfunction; it indicates that the application layer permitted conflicting lock acquisition patterns.",
    mental_model="Uncoordinated concurrent resource locks inevitably produce circular dependencies under sustained traffic.",
    source_path="content/02-master-draft.md:3-16",
    slide_num="02 / 24",
    layout="cards"
)

# Slide 03: Mengapa Review Manual / Pendekatan Naif Gagal (Layout: cards)
# Limiting phrase 2: "tidak otomatis"
add_slide(
    categories=["WHY FAILS", "CONCURRENCY ANOMALY"],
    title="Review Kode Manual Gagal Menangkap Siklus Tunggu Konkuren",
    subtitle="Mengandalkan insting tanpa standarisasi urutan lock membiarkan celah deadlock lolos ke produksi.",
    content_items=[
        {"type": "cards", "title": "KELEMAHAN PENDEKATAN NAIF PENGUJIAN", "items": [
            {"title": "Interleaving Gaib", "desc": "Jadwal eksekusi thread bergantung pada thread scheduler runtime dan OS, sulit direproduksi lokal."},
            {"title": "Unit Test Tunggal Lolos", "desc": "Pengujian fungsional sekuensial selalu berhasil karena kompetisi sumber daya tidak pernah terpicu."},
            {"title": "Tidak Otomatis Aman", "desc": "Membungkus transaksi dengan blok ACID tidak otomatis aman dari kebuntuan bila urutan query berbeda."}
        ]}
    ],
    narrative="Review manual dan unit test sekuensial tidak mampu mereproduksi deadlock karena eksekusi konkuren bergantung pada timing scheduler CPU. Menjalankan query di dalam transaksi ACID tidak otomatis melindungi sistem dari deadlock jika dua alur bisnis mengunci entitas dalam urutan berlawanan.",
    narrative_en="Manual inspection and single-threaded unit tests fail to expose deadlocks because concurrency depends on operating system thread interleaving. Encapsulating queries within ACID transactions does not automatically guarantee deadlock safety if separate endpoints query tables in opposing orders.",
    mental_model="Sequential testing cannot reveal non-deterministic timing windows in concurrent resource acquisition.",
    source_path="content/02-master-draft.md:17-33",
    slide_num="03 / 24",
    layout="cards"
)

# Slide 04: Mental Model Inti — Schema Card (Layout: schema)
# Limiting phrase 3: "hanya pada jalur ini"
add_slide(
    categories=["MENTAL MODEL", "CIRCULAR WAIT"],
    title="Anatomi Siklus Tunggu: Deadly Embrace Dua Transaksi",
    subtitle="Empat syarat Coffman terpenuhi secara simultan hingga membentuk dependensi melingkar.",
    content_items=[
        {"type": "schema", "title": "EMPAT TAHAP TERBENTUKNYA CIRCULAR WAIT", "items": [
            {"phase": "1. Mutual Exclusion", "role": "Row Exclusive", "desc": "Satu entitas data akun hanya dapat dikunci dan dimutasi oleh satu transaksi pada satu waktu."},
            {"phase": "2. Hold and Wait", "role": "Lock Retained", "desc": "Transaksi 1 memegang Lock Akun A sambil menunggu pelepasan Lock Akun B dari pihak lain."},
            {"phase": "3. No Preemption", "role": "Voluntary Release", "desc": "Kunci data tidak dapat direbut paksa; hanya dilepas ketika commit atau rollback terjadi."},
            {"phase": "4. Circular Wait", "role": "Deadly Embrace", "desc": "Hanya pada jalur ini transaksi saling mengunci dan terhenti permanen tanpa intervensi luar."}
        ]}
    ],
    narrative="Deadlock terjadi ketika empat syarat Coffman hadir bersamaan: eksklusivitas, menahan lock lama, tanpa perebutan paksa, dan siklus tunggu melingkar. Hanya pada jalur ini kondisi saling sandera terkunci permanen, menuntut adanya intervensi eksternal seperti pembatalan paksa oleh database engine.",
    narrative_en="Deadlocks manifest when all four Coffman conditions align: mutual exclusion, hold-and-wait, no preemption, and circular wait. Only along this path does an unbreakable mutual lock freeze execution permanently, necessitating external abort intervention by the engine.",
    mental_model="A closed circular dependency graph locks all participants until an external entity terminates one node.",
    source_path="content/02-master-draft.md:17-42",
    slide_num="04 / 24",
    layout="schema"
)

# Slide 05: Diagram Alur / Pipeline (Layout: schema)
# Limiting phrase 4: "belum membuktikan"
add_slide(
    categories=["PIPELINE", "DETECTION FLOW"],
    title="Alur Deteksi Mesin Basis Data: Dari Wait-For Graph ke Abort",
    subtitle="Bagaimana database engine mendeteksi siklus dan memilih satu transaksi sebagai korban.",
    content_items=[
        {"type": "schema", "title": "TAHAPAN PENANGANAN DEADLOCK INTERNAL DATABASE", "items": [
            {"phase": "1. Lock Contention", "role": "Wait Queue", "desc": "Transaksi meminta lock pada baris data sibuk dan dimasukkan ke dalam antrean tunggu."},
            {"phase": "2. Timeout Trigger", "role": "deadlock_timeout", "desc": "Setelah ambang waktu terlewati (default PG 1s), thread monitor pengecekan siklus dijalankan."},
            {"phase": "3. Cycle Detection", "role": "Wait-For Graph", "desc": "Algoritma DFS memeriksa graf dependensi transaksi untuk menemukan siklus tertutup."},
            {"phase": "4. Victim Selection", "role": "Engine Abort", "desc": "Mesin membatalkan transaksi korban; belum membuktikan algoritma cost-based di luar lab."}
        ]}
    ],
    narrative="Basis data tidak memeriksa deadlock secara terus-menerus demi menghemat CPU, melainkan menunggu jeda timeout tertentu sebelum menelusuri Wait-For Graph. Simulasi lab membatalkan transaksi via context timeout, tetapi belum membuktikan algoritma pemilihan korban berbasis biaya modifikasi terkecil seperti MySQL InnoDB.",
    narrative_en="Database engines postpone cycle checks to conserve CPU overhead, inspecting the Wait-For Graph only after a timeout elapses. While our lab models this via context timeouts, it has not yet proven production cost-based victim selection heuristics like MySQL InnoDB.",
    mental_model="Periodic background graph analysis identifies cyclic dependencies and evicts minimal-cost victim transactions.",
    source_path="content/02-master-draft.md:53-70",
    slide_num="05 / 24",
    layout="schema"
)

# Slide 06: Cuplikan Kode Dominan (Layout: code)
# Limiting phrase 5: "bukan berarti"
add_slide(
    categories=["IMPLEMENTATION", "DEADLOCK SIMULATION"],
    title="Pola Penguncian Naif yang Menimbulkan Deadlock",
    subtitle="Mengakuisisi lock sesuai urutan input membiarkan arah transaksi saling bertubrukan.",
    content_items=[
        {"type": "code", "title": "internal/transfer/transfer.go", "items": [
            "func TransferNaive(ctx context.Context, from, to *bank.Account, amount int, delay time.Duration) error {",
            "    if err := from.Lock(ctx); err != nil { return err }",
            "    defer from.Unlock()",
            "    time.Sleep(delay) // simulasi durasi transaksi internal",
            "    if err := to.Lock(ctx); err != nil { return err }",
            "    defer to.Unlock()",
            "    from.Balance -= amount",
            "    to.Balance += amount",
            "    return nil",
            "}"
        ]}
    ],
    narrative="Fungsi TransferNaive mengunci akun pengirim lalu akun penerima sesuai urutan parameter pemanggilan. Menambahkan jeda waktu eksekusi bukan berarti memperlambat kode secara sengaja, melainkan memodelkan latensi query nyata yang membuka celah tabrakan bagi transaksi berlawanan arah.",
    narrative_en="TransferNaive acquires locks sequentially based purely on function argument order. Injecting artificial processing delay does not mean introducing arbitrary bloat; it accurately simulates query execution latency that opens the concurrency collision window.",
    mental_model="Binding lock acquisition order directly to call arguments creates bidirectional lock acquisition hazards.",
    source_path="internal/transfer/transfer.go:9-26",
    slide_num="06 / 24",
    layout="code"
)

# Slide 07: Failure Path / Dangling Reference Konkret (Layout: cards)
# Limiting phrase 6: "hanya pada jalur ini"
add_slide(
    categories=["FAILURE PATH", "DEADLY EMBRACE"],
    title="Kronologi Tabrakan: Goroutine A vs Goroutine B",
    subtitle="Langkah demi langkah runtuhnya sistem akibat dua transaksi saling menyandera kunci.",
    content_items=[
        {"type": "cards", "title": "FASE KEGAGALAN TRANSAKSI SILANG", "items": [
            {"title": "Langkah 1: Akuisisi Parsial", "desc": "Goroutine 1 berhasil mengunci Akun A; Goroutine 2 bersamaan mengunci Akun B."},
            {"title": "Langkah 2: Saling Menunggu", "desc": "Goroutine 1 mencoba mengunci Akun B (sibuk); Goroutine 2 mencoba mengunci Akun A (sibuk)."},
            {"title": "Langkah 3: Kebuntuan Mutlak", "desc": "Kedua goroutine terblokir; tidak ada yang dapat melanjutkan atau melepas lock secara sukarela."},
            {"title": "Langkah 4: Terminasi Paksa", "desc": "Hanya pada jalur ini context timeout memutus rantai dan mengembalikan galat ErrDeadlock."}
        ]}
    ],
    narrative="Ketika dua transfer berjalan simultan dengan arah berlawanan, masing-masing berhasil mendapatkan kunci pertama namun gagal memperoleh kunci kedua. Hanya pada jalur ini kebuntuan absolut terjadi, memaksa timer pembatalan mengorbankan salah satu transaksi agar yang lain dapat selesai.",
    narrative_en="When bidirectional transfers execute concurrently, each successfully secures its first lock but blocks on the second. Only along this path does an absolute freeze emerge, forcing timeout cancellation to sacrifice one transaction so the survivor can complete.",
    mental_model="Simultaneous cross-resource acquisition without preemption leads to an unresolvable terminal deadlock.",
    source_path="content/04-diagrams.md:8-27",
    slide_num="07 / 24",
    layout="cards"
)

# Slide 08: Tabel Perbandingan Mekanisme vs Gejala (Layout: table)
add_slide(
    categories=["COMPARISON", "DIAGNOSTICS"],
    title="Matriks Diagnostik: Fenomena Penguncian di Database",
    subtitle="Membedakan karakteristik teknis deadlock, lock timeout, dan lock contention biasa.",
    content_items=[
        {"title": "KARAKTERISTIK GANGGUAN PENGUNCIAN SISTEM", "items": [
            "Dimensi Masalah | Normal Lock Contention | Lock Timeout | Deadlock (Circular Wait)",
            "Kondisi Dependensi | Antrean linear satu arah | Tunggu melebihi ambang batas | Siklus tertutup saling sandera (A->B->A)",
            "Penyelesaian Mesin | Menunggu lock dilepas | Transaksi pemohon dibatalkan | Satu transaksi dikorbankan via rollback",
            "Kode Galat Umum | Tidak ada galat (hanya antre) | MySQL 1205 (wait timeout) | PG 40P01 / MySQL 1213 / SQLServer 1205",
            "Akar Masalah | Beban transaksi tinggi | Transaksi lambat / long-running | Pola urutan akuisisi lock tidak konsisten",
            "Solusi Arsitektural | Indexing & isolasi lebih rendah | Perpendek durasi transaksi | Lock Ordering konsisten & Retry Pattern"
        ], "type": "table"}
    ],
    narrative="Tabel ini memetakan perbedaan fundamental antara antrean lock normal, lock timeout, dan kondisi deadlock murni. Deadlock memiliki ciri khas dependensi siklik yang tidak dapat selesai sendiri, sehingga selalu memerlukan intervensi rollback aktif dari engine basis data.",
    narrative_en="This matrix delineates the technical boundaries separating standard lock contention, lock timeouts, and genuine deadlocks. Deadlocks uniquely exhibit cyclic dependencies that cannot self-resolve, demanding proactive engine rollback intervention.",
    mental_model="Distinguish linear queuing delays from cyclic deadlock dependencies requiring automatic abort intervention.",
    source_path="content/02-master-draft.md:34-42",
    slide_num="08 / 24",
    layout="table"
)

# Slide 09: Decision Sequence — 3 Focused Decision Cards (Layout: cards)
# Limiting phrase 7: "bukan berarti"
add_slide(
    categories=["DECISION SEQUENCE", "MITIGATION STRATEGY"],
    title="Tiga Lapisan Pertahanan Menghadapi Deadlock",
    subtitle="Hierarki strategi rekayasa perangkat lunak dari pencegahan dini hingga pemulihan aktif.",
    content_items=[
        {"type": "cards", "title": "STRATEGI PENANGANAN DEADLOCK BERTAHAP", "items": [
            {"title": "1. Lock Ordering", "desc": "Urutkan seluruh resource berdasarkan ID sebelum akuisisi kunci; memutus circular wait secara matematis."},
            {"title": "2. Perpendek Transaksi", "desc": "Keluarkan panggilan jaringan pihak ketiga dan kalkulasi berat dari dalam blok transaksi SQL."},
            {"title": "3. Application Retry", "desc": "Tangkap kode galat deadlock; bukan berarti pasrah, melainkan pulihkan transaksi via backoff."}
        ]}
    ],
    narrative="Penanganan deadlock yang kokoh memadukan pencegahan dan ketahanan aplikasi. Mengimplementasikan retry otomatis bukan berarti mengabaikan pencegahan, melainkan jaring pengaman wajib untuk skenario konkurensi ekstrem yang tidak terduga oleh skema ordering.",
    narrative_en="Robust deadlock management pairs proactive architectural avoidance with operational resilience. Implementing automatic retry logic does not mean surrendering to lazy design; it serves as an indispensable safety net for complex concurrency patterns.",
    mental_model="Prevent cycles deterministically via lock ordering, minimize lock duration, and recover residual aborts via retries.",
    source_path="content/02-master-draft.md:64-70",
    slide_num="09 / 24",
    layout="cards"
)

# Slide 10: Bukti Eksekusi / Terminal Output (Layout: code)
add_slide(
    categories=["EVIDENCE", "EXECUTION"],
    title="Eksekusi Test Suite dan Demo cmd/demo Lolos Bersih",
    subtitle="Output terminal membuktikan reproduksi deadlock, prevensi ordering, dan pemulihan retry.",
    content_items=[
        {"type": "code", "title": "TERMINAL OUTPUT: go test -v -race ./... && go run ./cmd/demo", "items": [
            "=== RUN   TestDeadlockOccurrence",
            "--- PASS: TestDeadlockOccurrence (0.02s) [victim identified]",
            "=== RUN   TestLockOrderingPreventsDeadlock",
            "--- PASS: TestLockOrderingPreventsDeadlock (0.01s) [balances consistent]",
            "=== RUN   TestRetryRecoversDeadlock",
            "--- PASS: TestRetryRecoversDeadlock (0.02s) [both transfers complete]",
            "=== RUN   TestTransactionDurationImpact",
            "--- PASS: TestTransactionDurationImpact (0.05s) [long >= short]",
            "PASS: ok  labs/18-deadlock/tests  1.715s (race detector: clean)",
            "1. Simulating Naive Concurrent Transfers -> deadlock victim detected",
            "2. Simulating Lock Ordered Transfers    -> ACC-1=1100, ACC-2=900 (success)",
            "3. Simulating Retry Mechanism           -> Final Balances: ACC-1=1100, ACC-2=900"
        ]}
    ],
    narrative="Verifikasi terminal membuktikan seluruh skenario uji lolos tanpa race condition di bawah pengawasan Go race detector. Demonstrasi CLI mengonfirmasi bahwa pendekatan naif menghasilkan deadlock victim, sedangkan pengurutan kunci dan mekanisme retry menyelesaikan mutasi dengan konsisten.",
    narrative_en="Terminal verification demonstrates that all test suites pass cleanly under the Go race detector with zero data races. The CLI demo confirms that naive execution yields deadlock victims, whereas lock ordering and application retries ensure consistent final state.",
    mental_model="Executable test proof confirms deterministic reproduction, mathematical prevention, and transparent recovery.",
    source_path="engineering/03-execution-result.md:14-61",
    slide_num="10 / 24",
    layout="code"
)

# Slide 11: Analisis Batasan Bukti (Layout: cards)
# Limiting phrase 8: "belum membuktikan"
add_slide(
    categories=["AUDIT", "EVIDENCE BOUNDARY"],
    title="Batasan Bukti Lab: Yang Dibuktikan vs Belum",
    subtitle="Menegaskan batasan simulasi in-memory terhadap kompleksitas database engine riil.",
    content_items=[
        {"type": "cards", "title": "BATAS PEMBUKTIAN SISTEM LAB 18", "items": [
            {"title": "YANG DIBUKTIKAN", "desc": "Reproduksi circular wait, keampuhan pengurutan ID kunci, dan pemulihan saldo via retry loop."},
            {"title": "YANG TIDAK DIBUKTIKAN", "desc": "Mesin graph parser WFG riil, multi-table foreign key cascades, dan gap locking di database."},
            {"title": "BATAS KLAIM EMPIRIS", "desc": "Lab belum membuktikan ketahanan sistem pada skala 50.000 TPS dengan distribusi kunci acak."}
        ]}
    ],
    narrative="Lab ini membuktikan mekanika dasar deadlock dan validitas teknik pengurutan kunci menggunakan channel in-memory. Pengujian ini belum membuktikan fenomena kompleks seperti gap locking pada indeks database riil atau perilaku sistem pada beban puluhan ribu TPS.",
    narrative_en="This lab proves fundamental deadlock mechanics and lock ordering efficacy using in-memory channel primitives. It has not yet proven production phenomena such as database index gap locking or systemic resilience under 50,000 transactions per second.",
    mental_model="Strictly separate application-level synchronization models from relational storage engine locking internals.",
    source_path="content/01-content-brief.md:12-16",
    slide_num="11 / 24",
    layout="cards"
)

# Slide 12: Trade-off Arsitektur (Layout: cards)
# Limiting phrase 9: "tidak otomatis"
add_slide(
    categories=["TRADE-OFF", "COSTS & BENEFITS"],
    title="Kompromi Rekayasa: Lock Ordering vs Application Retry",
    subtitle="Mengevaluasi konsekuensi operasional antara pengurutan ketat dan toleransi pembatalan.",
    content_items=[
        {"type": "cards", "title": "ANALISIS PERBANDINGAN TRADE-OFF", "items": [
            {"title": "Lock Ordering: Disiplin Kode", "desc": "Bebas deadlock sepenuhnya, namun mewajibkan kesepakatan urutan global di seluruh tim pengembang."},
            {"title": "Retry: Beban Komputasi", "desc": "Mudah dipasang di wrapper, namun memboroskan CPU dan I/O jika terjadi abort berulang-ulang."},
            {"title": "Tidak Otomatis Cepat", "desc": "Menambah retry tidak otomatis meningkatkan throughput bila contention pada satu baris data terlalu pekat."}
        ]}
    ],
    narrative="Penerapan lock ordering memberikan kepastian absolut namun menuntut disiplin penulisan query yang seragam di seluruh repository. Menambahkan retry logic tidak otomatis membuat sistem berkinerja tinggi jika sumber daya yang diperebutkan mengalami saturasi ekstrem.",
    narrative_en="Enforcing global lock ordering guarantees safety but requires rigorous query conventions across all services. Introducing retry wrappers does not automatically enhance throughput if underlying shared resources suffer from acute lock contention.",
    mental_model="Choose between upfront developmental discipline via ordering or operational compute overhead via retries.",
    source_path="content/02-master-draft.md:64-70",
    slide_num="12 / 24",
    layout="cards"
)

# Slide 13: Concurrency / Synchronization Primitive (Layout: code)
# Limiting phrase 10: "simulasi ini tidak sama dengan"
add_slide(
    categories=["CONCURRENCY", "PRIMITIVES"],
    title="Model Sinkronisasi Akun Berbasis Channel dan Timeout",
    subtitle="Menggunakan channel kapasitas 1 sebagai mutex yang peka terhadap pembatalan context.",
    content_items=[
        {"type": "code", "title": "internal/bank/account.go", "items": [
            "type Account struct {",
            "    ID      string",
            "    Balance int",
            "    ch      chan struct{} // channel buffer 1 sebagai token lock",
            "}",
            "func (a *Account) Lock(ctx context.Context) error {",
            "    select {",
            "    case <-a.ch: return nil",
            "    case <-ctx.Done(): return ErrDeadlock // representasi deadlock abort",
            "    }",
            "}"
        ]}
    ],
    narrative="Model Account memanfaatkan buffered channel 1 slot untuk mengendalikan hak akses mutasi saldo. Simulasi ini tidak sama dengan locking row-level di PostgreSQL yang memiliki status lock modes bertingkat, namun secara elegan merefleksikan pembatalan transaksi akibat timeout eksternal.",
    narrative_en="The Account struct utilizes a 1-slot buffered channel to govern balance mutation access. While this in-memory simulation is not identical to PostgreSQL row-level lock trees with shared/exclusive modes, it cleanly models external deadlock abortion.",
    mental_model="Channel token acquisition coupled with context cancellation models interruptible resource acquisition.",
    source_path="internal/bank/account.go:10-33",
    slide_num="13 / 24",
    layout="code"
)

# Slide 14: Before vs After Comparison Cards (Layout: cards)
# Limiting phrase 11: "bukan berarti"
add_slide(
    categories=["COMPARISON", "BEFORE VS AFTER"],
    title="Transformasi: Dari Tabrakan Acak ke Eksekusi Teratur",
    subtitle="Pergeseran dari pemanggilan naif yang rentan menuju orkestrasi penguncian deterministik.",
    content_items=[
        {"type": "cards", "title": "KOMPARASI KONDISI SISTEM TRANSAKSI", "items": [
            {"title": "Sebelum: Urutan Dinamis", "desc": "Kunci diambil mengikuti urutan parameter fungsi; transfer dua arah saling menjebak dalam circular wait."},
            {"title": "Sesudah: Leksikografis", "desc": "ID akun diurutkan secara konsisten; kedua proses antre secara teratur pada lock pertama yang sama."},
            {"title": "Bukan Berarti Lambat", "desc": "Mengantre pada satu kunci bukan berarti menurunkan throughput; sistem terhindar dari biaya rollback."}
        ]}
    ],
    narrative="Sebelum standarisasi urutan, dua transaksi yang bersilangan akan saling menunggu dan berujung pada pembatalan sepihak. Mengantrekan transaksi pada kunci pertama bukan berarti memperlambat sistem, melainkan menyelamatkan CPU dari pemborosan rollback dan retry.",
    narrative_en="Prior to lock ordering, intersecting transactions mutually block and trigger destructive engine rollbacks. Queuing transactions deterministically on an initial shared lock does not degrade system throughput; it prevents wasted rollback cycles.",
    mental_model="Convert circular blocking hazards into deterministic linear queues via strict lexical key sorting.",
    source_path="content/02-master-draft.md:148-166",
    slide_num="14 / 24",
    layout="cards"
)

# Slide 15: Analisis Mendalam — Logika Lock Ordering (Layout: code)
# Limiting phrase 12: "cakupannya mengikuti test"
add_slide(
    categories=["ANALYSIS", "PREVENTION LOGIC"],
    title="Pencegahan Deterministik: Pengurutan Kunci Berdasarkan ID",
    subtitle="Menjamin bahwa semua proses selalu meminta lock dengan urutan seragam tanpa siklus.",
    content_items=[
        {"type": "code", "title": "internal/transfer/transfer.go", "items": [
            "func TransferOrdered(ctx context.Context, acc1, acc2 *bank.Account, amount int, delay time.Duration) error {",
            "    first, second := acc1, acc2",
            "    if acc1.ID > acc2.ID { first, second = acc2, acc1 }",
            "    if err := first.Lock(ctx); err != nil { return err }",
            "    defer first.Unlock()",
            "    time.Sleep(delay)",
            "    if err := second.Lock(ctx); err != nil { return err }",
            "    defer second.Unlock()",
            "    acc1.Balance -= amount; acc2.Balance += amount",
            "    return nil",
            "}"
        ]}
    ],
    narrative="Dengan membandingkan identifier string secara leksikografis, entitas dengan ID lebih kecil selalu dikunci terlebih dahulu. Logika pencegahan ini cakupannya mengikuti test konkurensi lab, menjamin bahwa circular wait mustahil terbentuk karena seluruh goroutine mengantre pada kunci pertama yang sama.",
    narrative_en="By sorting resource IDs lexicographically, the entity with the lower ID is invariably locked first. The scope of this prevention logic strictly matches our concurrency test suite, proving circular waits cannot form because competing routines queue on the identical first lock.",
    mental_model="Global monotonic resource acquisition eliminates cycles across all concurrent transaction paths.",
    source_path="internal/transfer/transfer.go:28-50",
    slide_num="15 / 24",
    layout="code"
)

# Slide 16: Analisis Mendalam — Test Bukti Terjadinya Deadlock (Layout: code)
# Limiting phrase 13: "tidak otomatis"
add_slide(
    categories=["ANALYSIS", "FAILURE VERIFICATION"],
    title="Verifikasi Test: Membuktikan Terjadinya Deadlock Victim",
    subtitle="Dua goroutine bersaing memicu tabrakan deterministik pada penguncian naif.",
    content_items=[
        {"type": "code", "title": "tests/transfer_test.go", "items": [
            "func TestDeadlockOccurrence(t *testing.T) {",
            "    // Goroutine 1 mentransfer A -> B (delay 10ms, timeout 20ms)",
            "    go func() { errs[0] = transfer.TransferNaive(ctx1, accA, accB, 10, 10*time.Millisecond) }()",
            "    // Goroutine 2 bersamaan mentransfer B -> A (delay 10ms, timeout 20ms)",
            "    go func() { errs[1] = transfer.TransferNaive(ctx2, accB, accA, 10, 10*time.Millisecond) }()",
            "    wg.Wait()",
            "    hasDeadlock := (errs[0] == bank.ErrDeadlock) || (errs[1] == bank.ErrDeadlock)",
            "    if !hasDeadlock { t.Fatalf(\"expected at least one deadlock victim\") }",
            "}"
        ]}
    ],
    narrative="Unit test membuktikan secara deterministik bahwa transfer silang tanpa pengurutan selalu mengorbankan minimal satu transaksi. Menjalankan dua goroutine paralel tidak otomatis memicu deadlock tanpa delay internal yang memadai untuk membuka jendela kontensi.",
    narrative_en="This unit test deterministically proves that cross transfers without ordering reliably produce at least one aborted victim. Spawning two goroutines concurrently does not automatically reproduce deadlocks unless processing delay sufficiently widens the contention window.",
    mental_model="Reproducing deadlocks requires synchronized contention windows where resources are held during secondary requests.",
    source_path="tests/transfer_test.go:13-41",
    slide_num="16 / 24",
    layout="code"
)

# Slide 17: Analisis Mendalam — Test Bukti Lock Ordering (Layout: code)
# Limiting phrase 14: "hanya pada jalur ini"
add_slide(
    categories=["ANALYSIS", "ORDERING VERIFICATION"],
    title="Verifikasi Test: Lock Ordering Menuntaskan Semua Transaksi",
    subtitle="Kedua transfer konkuren berhasil diselesaikan dengan saldo akhir yang konsisten sempurna.",
    content_items=[
        {"type": "code", "title": "tests/transfer_test.go", "items": [
            "func TestLockOrderingPreventsDeadlock(t *testing.T) {",
            "    go func() { errs[0] = transfer.TransferOrdered(ctx1, accA, accB, 10, 5*time.Millisecond) }()",
            "    go func() { errs[1] = transfer.TransferOrdered(ctx2, accB, accA, 10, 5*time.Millisecond) }()",
            "    wg.Wait()",
            "    if errs[0] != nil || errs[1] != nil { t.Fatalf(\"expected both to succeed\") }",
            "    if accA.Balance != 100 || accB.Balance != 100 {",
            "        t.Fatalf(\"unexpected balance: A=%d, B=%d\", accA.Balance, accB.Balance)",
            "    }",
            "}"
        ]}
    ],
    narrative="Pengujian membuktikan bahwa kedua transaksi bolak-balik berhasil dieksekusi tanpa error sama sekali saat menggunakan TransferOrdered. Hanya pada jalur ini konsistensi saldo terjaga utuh (100 dan 100) karena transaksi saling mengantre secara tertib alih-alih saling mematikan.",
    narrative_en="Testing confirms both bidirectional transfers complete successfully with zero errors when routed through TransferOrdered. Only along this path do balances remain perfectly consistent (100 and 100) as transactions serialize gracefully rather than aborting.",
    mental_model="Enforcing a single acquisition order linearizes resource contention, preserving full transaction integrity.",
    source_path="tests/transfer_test.go:43-74",
    slide_num="17 / 24",
    layout="code"
)

# Slide 18: Analisis Mendalam — Mekanisme Application-Level Retry (Layout: code)
# Limiting phrase 15: "bukan berarti"
add_slide(
    categories=["ANALYSIS", "RETRY PATTERN"],
    title="Pemulihan Transparan: Penanganan Abort dengan Retry",
    subtitle="Membungkus eksekusi dengan batas timeout percobaan dan jeda backoff terencana.",
    content_items=[
        {"type": "code", "title": "internal/transfer/transfer.go", "items": [
            "func TransferWithRetry(ctx context.Context, from, to *bank.Account, amount int, delay time.Duration, maxRetries int) error {",
            "    for i := 0; i < maxRetries; i++ {",
            "        attemptCtx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)",
            "        err := TransferNaive(attemptCtx, from, to, amount, delay)",
            "        cancel()",
            "        if err == nil { return nil }",
            "        if err != bank.ErrDeadlock && err != context.DeadlineExceeded { return err }",
            "        time.Sleep(2 * time.Millisecond) // backoff sebelum mencoba kembali",
            "    }",
            "    return bank.ErrDeadlock",
            "}"
        ]}
    ],
    narrative="Fungsi TransferWithRetry menangkap galat pembatalan deadlock dan mencoba kembali transaksi secara transparan hingga batas maksimal. Menggunakan fixed backoff 2ms bukan berarti pola ini ideal untuk produksi; implementasi ini adalah penyederhanaan terencana yang membutuhkan jitter pada skala besar.",
    narrative_en="TransferWithRetry intercepts deadlock aborts and transparently reissues the transaction up to maxRetries. Adopting a fixed 2ms backoff does not represent production best practice; it is a documented simplification that requires randomized jitter in production systems.",
    mental_model="Catch transient concurrency aborts and replay idempotent transactions following brief backoff intervals.",
    source_path="internal/transfer/transfer.go:52-71",
    slide_num="18 / 24",
    layout="code"
)

# Slide 19: Analisis Mendalam — Korelasi Durasi Transaksi (Layout: code)
# Limiting phrase 16: "cakupannya mengikuti test"
add_slide(
    categories=["ANALYSIS", "EMPIRICAL METRIC"],
    title="Dampak Durasi: Transaksi Panjang Memperbesar Risiko",
    subtitle="Membandingkan frekuensi deadlock antara transaksi dengan delay versus instan.",
    content_items=[
        {"type": "code", "title": "tests/transfer_test.go", "items": [
            "func TestTransactionDurationImpact(t *testing.T) {",
            "    longDeadlocks := countDeadlocks(5*time.Millisecond, 10)",
            "    shortDeadlocks := countDeadlocks(0, 10)",
            "    if longDeadlocks < shortDeadlocks {",
            "        t.Fatalf(\"expected long duration to have >= deadlocks: long=%d, short=%d\",",
            "                 longDeadlocks, shortDeadlocks)",
            "    }",
            "}"
        ]}
    ],
    narrative="Pengujian empiris membuktikan bahwa transaksi yang menahan lock selama 5 milidetik menghasilkan frekuensi deadlock yang lebih tinggi atau sama dibandingkan transaksi tanpa delay. Pengamatan ini cakupannya mengikuti test lab, memvalidasi prinsip bahwa memperpanjang waktu penahanan lock memperbesar kemungkinan tabrakan.",
    narrative_en="Empirical testing demonstrates that transactions holding locks for 5ms produce higher or equal deadlock counts compared to zero-delay transactions. This finding directly reflects test suite metrics, confirming that prolonged lock durations expand the concurrency collision surface.",
    mental_model="Deadlock probability scales proportionally with the duration locks are held open by active transactions.",
    source_path="tests/transfer_test.go:105-145",
    slide_num="19 / 24",
    layout="code"
)

# Slide 20: Checklist Implementasi (Layout: checklist)
add_slide(
    categories=["CHECKLIST", "PRODUCTION READINESS"],
    title="Checklist Proteksi Deadlock di Lingkungan Produksi",
    subtitle="Panduan komprehensif sebelum merilis kode transaksi konkuren ke basis data.",
    content_items=[
        {"type": "checklist", "title": "STANDAR REKAYASA KONEKSI BASIS DATA", "items": [
            "Terapkan Lock Ordering seragam pada seluruh mutasi multi-tabel di seluruh repositori.",
            "Lakukan ORDER BY primary_key ASC sebelum query UPDATE batch atau SELECT FOR UPDATE.",
            "Keluarkan panggilan HTTP eksternal dan komputasi berat dari dalam transaksi database.",
            "Tangkap kode galat deadlock vendor (Postgres 40P01, MySQL 1213, SQL Server 1205).",
            "Terapkan exponential backoff dengan randomized jitter pada lapisan retry aplikasi.",
            "Pastikan seluruh operasi yang di-retry bersifat idempotent untuk mencegah mutasi ganda."
        ]}
    ],
    narrative="Checklist ini merangkum langkah mitigasi wajib: urutan penguncian global yang seragam, pengurutan kunci pada update massal, pembebasan transaksi dari panggilan eksternal, penangkapan kode galat vendor spesifik, retry berbasis jitter, serta jaminan idempotensi.",
    narrative_en="This checklist codifies production defenses: uniform global lock ordering, primary key sorting on bulk updates, stripping external I/O from transactions, vendor error code interception, jittered exponential backoff, and strict idempotency.",
    mental_model="Enforce strict operational criteria to design, isolate, and recover concurrent database transactions.",
    source_path="content/02-master-draft.md:273-280",
    slide_num="20 / 24",
    layout="checklist"
)

# Slide 21: Anti-Pattern & Cara Reproduksi Kegagalan (Layout: cards)
# Limiting phrase 17: "bukan berarti"
add_slide(
    categories=["ANTI-PATTERN", "PITFALLS"],
    title="Anti-Pattern Umum dalam Penanganan Transaksi",
    subtitle="Kesalahan desain yang sering memicu deadlock persisten pada sistem produksi.",
    content_items=[
        {"type": "cards", "title": "POLA BERBAHAYA YANG HARUS DIHINDARI", "items": [
            {"title": "I/O Eksternal di Transaksi", "desc": "Memanggil API pihak ketiga atau menunggu I/O disk saat koneksi database masih memegang row lock."},
            {"title": "Batch Update Acak", "desc": "Mengupdate ratusan baris data dari array memori tanpa mengurutkan ID terlebih dahulu."},
            {"title": "Bukan Berarti Aman", "desc": "Memasang retry tanpa idempotensi bukan berarti aman, melainkan memicu mutasi saldo ganda."}
        ]}
    ],
    narrative="Menahan koneksi transaksi database sambil menunggu respons API pihak ketiga adalah penyebab paling umum lonjakan deadlock di produksi. Memasang wrapper retry tanpa memastikan idempotensi mutasi bukan berarti sistem menjadi andal, melainkan berisiko mencatatkan perubahan saldo ganda.",
    narrative_en="Holding database locks while waiting on third-party HTTP responses represents the most frequent cause of production deadlocks. Wrapping endpoints in retry logic without verifying mutation idempotency does not yield safety; it risks double financial debits.",
    mental_model="Avoid external I/O within transaction boundaries and never retry non-idempotent business operations.",
    source_path="content/02-master-draft.md:257-262",
    slide_num="21 / 24",
    layout="cards"
)

# Slide 22: Takeaway Praktis & Command Block Dominan (Layout: code)
# Limiting phrase 18: "production masih perlu"
add_slide(
    categories=["REPRODUCTION", "COMMANDS"],
    title="Reproduksi Mandiri di Mesin Lokal Menggunakan Go",
    subtitle="Verifikasi mandiri simulasi deadlock tanpa ketergantungan modul eksternal.",
    content_items=[
        {"type": "code", "title": "COMMAND RUNNER (ZERO DEPENDENCIES)", "items": [
            "cd labs/18-deadlock",
            "go test -v ./...    # verifikasi 4 pengujian deadlock, ordering, retry & durasi",
            "go test -race ./... # pembuktian sinkronisasi konkuren bebas data race",
            "go run ./cmd/demo   # eksekusi demonstrasi CLI interaktif skenario 1, 2, dan 3"
        ]}
    ],
    narrative="Seluruh simulasi dibangun murni menggunakan pustaka standar Go tanpa modul pihak ketiga. Eksperimen lokal ini memvalidasi mekanika deadlock di memori, namun production masih perlu mengonfigurasi parameter database nyata seperti deadlock_timeout dan log_lock_waits.",
    narrative_en="The entire simulation runs on the Go standard library with zero external dependencies. While local testing confirms in-memory algorithmic principles, production systems still require configuring database parameters like deadlock_timeout and log_lock_waits.",
    mental_model="Zero-dependency concurrency simulations provide fast, deterministic validation of distributed failure modes.",
    source_path="README.md:14-34",
    slide_num="22 / 24",
    layout="code"
)

# Slide 23: Prinsip Arsitektural Jangka Panjang (Layout: cards)
# Limiting phrase 19: "tidak otomatis"
add_slide(
    categories=["TAKEAWAY", "ARCHITECTURAL PRINCIPLES"],
    title="Prinsip Inti: Deadlock adalah Keniscayaan Konkurensi",
    subtitle="Pondasi arsitektur tangguh bertumpu pada pencegahan terstruktur dan toleransi kegagalan.",
    content_items=[
        {"type": "cards", "title": "TIGA PILAR KETAHANAN KONKURENSI SISTEM", "items": [
            {"title": "Disiplin Urutan Global", "desc": "Jadikan pengurutan kunci sebagai standar baku query di seluruh lapisan arsitektur microservices."},
            {"title": "Keringkasan Transaksi", "desc": "Pangkas waktu hidup transaksi seminimal mungkin; lock hanya saat mutasi siap dieksekusi."},
            {"title": "Toleransi Transparan", "desc": "Rancang aplikasi agar tidak otomatis panik saat menerima abort; pulihkan via retry transparan."}
        ]}
    ],
    narrative="Deadlock tidak boleh dipandang sebagai aib atau bug fatal, melainkan mekanisme perlindungan integritas data dari basis data. Merancang sistem transaksional tidak otomatis membebaskan kita dari kontensi; sistem yang tangguh mencegah siklus sedini mungkin dan memulihkan korban dengan anggun.",
    narrative_en="Deadlocks should not be viewed as catastrophic bugs, but as the database engine's ultimate defense of data consistency. Designing transactional architectures does not automatically eliminate lock contention; resilient systems prevent cycles upfront and recover victims gracefully.",
    mental_model="Design for deterministic lock ordering upfront, minimize critical sections, and treat abort recovery as standard operating procedure.",
    source_path="content/05-key-takeaways.md:1-8",
    slide_num="23 / 24",
    layout="cards"
)

# Slide 24: Closing Slide (Layout: closing)
add_slide(
    categories=["CLOSING"],
    title="Deadlock bukan kegagalan basis data, melainkan alarm bahwa aplikasi Anda mengunci sumber daya secara acak.",
    subtitle="",
    content_items=[
        {"type": "bullet", "title": "MODUL IMPLEMENTASI", "items": [
            "internal/bank/account.go • internal/transfer/transfer.go",
            "cmd/demo/main.go • tests/transfer_test.go",
            "engineering/01-design.md • content/02-master-draft.md"
        ]},
        {"type": "bullet", "title": "VERIFIKASI & STATUS", "items": [
            f"Commit: {COMMIT_SHA[:7]} • Status: APPROVED",
            "Test passing: 4 test suites (DeadlockOccurrence, LockOrdering, RetryRecovery, DurationImpact)",
            "Race detector: Clean (go test -race ./... PASS in 1.715s)",
            "Zero external dependencies: 100% Go standard library."
        ]}
    ],
    mental_model="Deadlocks are structural alarms indicating uncoordinated concurrent resource access patterns.",
    source_path=REPO_URL,
    slide_num="24 / 24",
    layout="closing"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
    for s in SLIDES:
        print(f"[{s['slide_num']}] {s['layout']}: {s['title']}")
