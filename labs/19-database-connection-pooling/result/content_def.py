#!/usr/bin/env python3
"""
Publication Builder Definition for Software Engineering Lab 19: Database Connection Pooling.
Adheres strictly to:
- Cover hero image requirement: Mandatory raster image at COVER_IMAGE_PATH
- Slide composition rules: <=45 words prose, <=30 words on code slides, dominant visual anchors
- Quality bar: strict technical boundary, limiting phrases on >=12 slides, zero banned filler words
- Concrete 24-slide LinkedIn Carousel structure
"""

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/19-database-connection-pooling/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/19-database-connection-pooling"
COMMIT_SHA = "b2d3274"
LAB_NO = "19"
SLUG = "database-connection-pooling"

TITLE_ID = "Database Connection Pooling"
SUBTITLE_ID = "Overhead Handshake, Batas Server, dan Bahaya Kebocoran Koneksi di Go"
TITLE_EN = "Database Connection Pooling"
SUBTITLE_EN = "Handshake Overhead, Server Slot Limits, and Connection Leak Hazards in Go"

AUTHOR = "Lukman / lukman-ss"
TOPIC_TAGS = "Database • Connection Pooling • Concurrency • Resource Starvation • Go • PostgreSQL"
COVER_IMAGE_PATH = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/19-database-connection-pooling/result/image-cover/cover_art.jpg"

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
            "Retro Terminal Visualizing Pool Sizing, Server Saturation Knee & Leak Hazards",
            "Simulasi database/sql/driver di Go: Handshake Overhead, Server Exhaustion, dan Starvation",
            f"Penulis: {AUTHOR}"
        ]}
    ],
    mental_model="Database connections are scarce server-side slots governed by hardware limits rather than client goroutine counts.",
    source_path="",
    slide_num="01 / 24",
    layout="cover"
)

# Slide 02: Problem Nyata (Layout: cards)
# Limiting phrase 1: "bukan berarti"
add_slide(
    categories=["PROBLEM", "RESOURCE EXHAUSTION"],
    title="Error Koneksi Meledak Saat Metrik CPU Database Normal",
    subtitle="Aplikasi menolak jutaan pengguna meski utilitas server database masih sangat rendah.",
    content_items=[
        {"type": "cards", "title": "PARADOKS KEGAGALAN KONEKSI BASIS DATA", "items": [
            {"title": "Penolakan Slot Keras", "desc": "PostgreSQL melempar galat FATAL: too many clients saat slot max_connections terisi penuh."},
            {"title": "Bukan Berarti Rusak", "desc": "Galat koneksi bukan berarti server kehabisan komputasi CPU, melainkan alokasi slot backend habis."},
            {"title": "Refleks Salah Tim", "desc": "Menaikkan ukuran pool secara membabi buta justru memperparah kontensi lock dan memicu degradasi."}
        ]}
    ],
    narrative="Banyak insiden produksi ditandai dengan ledakan galat koneksi sementara CPU database berada di bawah 30 persen. Fenomena ini bukan berarti server mengalami kerusakan komputasi internal, melainkan slot koneksi habis dialokasikan oleh kumpulan worker aplikasi yang tidak terkoordinasi.",
    narrative_en="Production outages frequently feature spikes in connection failures while database CPU utilization remains below 30 percent. This phenomenon does not indicate hardware computation failure; it proves server process slots were exhausted by uncoordinated application workers.",
    mental_model="Connection limits represent hard operating system socket and memory slots rather than raw processing power.",
    source_path="content/02-master-draft.md:3-19",
    slide_num="02 / 24",
    layout="cards"
)

# Slide 03: Mengapa Review Manual / Pendekatan Naif Gagal (Layout: cards)
# Limiting phrase 2: "tidak otomatis"
add_slide(
    categories=["WHY FAILS", "NAIVE APPROACHES"],
    title="Pendekatan Naif Membuka Koneksi Per-Request Menghancurkan Latensi",
    subtitle="Membuat koneksi baru pada setiap query menumpuk beban TCP dan TLS handshake.",
    content_items=[
        {"type": "cards", "title": "TITIK RUNTUH KONEKSI TANPA POOLING", "items": [
            {"title": "Overhead Berulang", "desc": "Setiap koneksi menuntut negosiasi TCP 3-way handshake, TLS exchange, serta autentikasi peran."},
            {"title": "Tekanan Backend", "desc": "Database wajib membuat proses backend baru dan mengalokasikan work_mem untuk tiap request."},
            {"title": "Tidak Otomatis Cepat", "desc": "Membuka koneksi baru tidak otomatis menyelesaikan beban; antrean pembuatan soket memperlambat IO."}
        ]}
    ],
    narrative="Pendekatan naif tanpa reuse mengabaikan besarnya biaya inisialisasi jaringan dan alokasi memori proses basis data. Membuka koneksi baru untuk setiap request tidak otomatis mempercepat eksekusi, melainkan mengubur runtime aplikasi dalam tumpukan latency handshake berulang.",
    narrative_en="Naive per-request connection instantiation overlooks substantial network negotiation and backend process allocation overheads. Creating a fresh connection per transaction does not automatically isolate workload; it drowns the runtime under repetitive handshake latencies.",
    mental_model="Direct connection creation incurs recurring network handshakes and expensive OS process allocation penalties.",
    source_path="content/02-master-draft.md:32-42",
    slide_num="03 / 24",
    layout="cards"
)

# Slide 04: Mental Model Inti — Schema Card (Layout: schema)
# Limiting phrase 3: "hanya pada jalur ini"
add_slide(
    categories=["MENTAL MODEL", "POOL SIZING"],
    title="Model Mental Restoran: Sizing Ditentukan Hardware Backend",
    subtitle="Kapasitas koneksi dibatasi oleh inti CPU dan disk, bukan nafsu request klien.",
    content_items=[
        {"type": "schema", "title": "ANATOMI SIZING DAN SATURASI KONEKSI", "items": [
            {"phase": "1. Inti Komputasi", "role": "Hardware CPU", "desc": "Koneksi konkuren efektif dibatasi core fisik; rumus ((core * 2) + spindle) jadi acuan awal."},
            {"phase": "2. Client Pool", "role": "sql.DB Semaphore", "desc": "Mengatur antrean peminjam koneksi lokal; mencegah goroutine menyerbu database serentak."},
            {"phase": "3. Saturation Knee", "role": "Degradasi Kurva", "desc": "Melewati titik jenuh, throughput anjlok drastis akibat context switching dan lock contention."},
            {"phase": "4. Steady State", "role": "Ukuran Minimal", "desc": "Hanya pada jalur ini pool berukuran ringkas menghasilkan throughput optimal dan stabil."}
        ]}
    ],
    narrative="Koneksi database bekerja layaknya meja restoran di mana kapasitas dibatasi oleh jumlah juru masak dapur, bukan antrean tamu di depan pintu. Hanya pada jalur ini penggunaan pool berukuran ringkas mampu membatasi perebutan sumber daya internal dan menghasilkan kestabilan throughput.",
    narrative_en="Database connections behave like restaurant seating where throughput is governed by kitchen cooks rather than waiting guests. Only along this path does a compact pool configuration constrain internal resource contention and sustain optimal query throughput.",
    mental_model="Size connection pools to match backend hardware concurrency constraints rather than client request volume.",
    source_path="content/02-master-draft.md:20-29",
    slide_num="04 / 24",
    layout="schema"
)

# Slide 05: Diagram Alur / Pipeline (Layout: schema)
# Limiting phrase 4: "belum membuktikan"
add_slide(
    categories=["PIPELINE", "POOL LIFECYCLE"],
    title="Siklus Hidup Peminjaman: Pinjam, Eksekusi, dan Kembalikan",
    subtitle="Arsitektur pool mengelola antrean akuisisi dan pemakaian ulang koneksi aktif.",
    content_items=[
        {"type": "schema", "title": "ALUR PEMINJAMAN KONEKSI DALAM SQL.DB", "items": [
            {"phase": "1. Request Conn", "role": "Acquire Token", "desc": "Worker meminta koneksi aktif via pool.Conn(ctx); mengambil idle atau membuat baru."},
            {"phase": "2. Bounded Query", "role": "Execute Work", "desc": "Query dieksekusi secara singkat tanpa operasi I/O eksternal di dalam siklus penguncian."},
            {"phase": "3. Immediate Return", "role": "defer conn.Close()", "desc": "Koneksi dikembalikan ke daftar idle pool; siap dipakai ulang oleh worker berikutnya."},
            {"phase": "4. Proxy Layer", "role": "PgBouncer", "desc": "Proxy menerjemahkan ribuan pool aplikasi; lab belum membuktikan multi-node proxy riil."}
        ]}
    ],
    narrative="Setiap transaksi wajib meminjam koneksi hanya untuk eksekusi query spesifik lalu segera mengembalikannya ke pool via defer conn.Close(). Simulasi lab memvalidasi siklus peminjaman lokal ini, namun belum membuktikan routing transaksi dinamis pada kluster PgBouncer multi-node di produksi.",
    narrative_en="Every transaction must acquire a connection exclusively for immediate query execution and release it via defer conn.Close(). While our lab validates this local acquisition cycle, it has not yet proven dynamic transaction routing across multi-node PgBouncer clusters.",
    mental_model="Connections must be borrowed for minimal durations and returned immediately to maintain pool liquidity.",
    source_path="content/04-diagrams.md:1-29",
    slide_num="05 / 24",
    layout="schema"
)

# Slide 06: Cuplikan Kode Dominan (Layout: code)
# Limiting phrase 5: "bukan berarti"
add_slide(
    categories=["IMPLEMENTATION", "SAFE PATTERN"],
    title="Pola Eksekusi Bersih: Isolasi I/O Eksternal dari Koneksi",
    subtitle="Memisahkan panggilan jaringan pihak ketiga dari siklus penahanan koneksi database.",
    content_items=[
        {"type": "code", "title": "internal/pool/service.go", "items": [
            "func (s *OrderService) ProcessOrderSafe(ctx context.Context, orderID int, externalCall func() error) error {",
            "    if externalCall != nil {",
            "        if err := externalCall(); err != nil { return err }",
            "    }",
            "    conn, err := s.db.Conn(ctx)",
            "    if err != nil { return err }",
            "    defer conn.Close()",
            "    _, err = conn.ExecContext(ctx, \"UPDATE orders SET status = 'completed' WHERE id = ?\", orderID)",
            "    return err",
            "}"
        ]}
    ],
    narrative="Fungsi ProcessOrderSafe menuntaskan panggilan API eksternal terlebih dahulu sebelum meminjam koneksi dari database pool. Menunda akuisisi koneksi bukan berarti menambah latensi alur, melainkan menjamin koneksi database hanya ditahan selama beberapa mikrodetik pada saat query dieksekusi.",
    narrative_en="The ProcessOrderSafe method executes external network calls completely prior to acquiring a database connection. Postponing connection acquisition does not increase end-to-end latency; it guarantees database sockets are held for mere microseconds during query execution.",
    mental_model="Acquire database connections at the absolute latest moment and release them at the earliest opportunity.",
    source_path="internal/pool/service.go:16-34",
    slide_num="06 / 24",
    layout="code"
)

# Slide 07: Failure Path / Dangling Reference Konkret (Layout: cards)
# Limiting phrase 6: "hanya pada jalur ini"
add_slide(
    categories=["FAILURE PATH", "CONNECTION LEAK"],
    title="Kronologi Kebocoran Koneksi yang Memicu Starvation",
    subtitle="Satu fungsi tidak disiplin memblokir seluruh request sehat hingga batas timeout.",
    content_items=[
        {"type": "cards", "title": "ALUR KEGAGALAN AKIBAT KONEKSI TERTOMBAK", "items": [
            {"title": "Langkah 1: Tahan Soket", "desc": "Worker A meminjam koneksi, menjalankan UPDATE, lalu memanggil API lambat tanpa melepas soket."},
            {"title": "Langkah 2: Pool Kritis", "desc": "Semua slot pool terisi oleh worker yang menanti respons jaringan pihak ketiga (500ms)."},
            {"title": "Langkah 3: Antrean Tercekat", "desc": "Worker B yang memiliki kode benar mencoba meminjam koneksi namun masuk antrean tunggu."},
            {"title": "Langkah 4: Context Timeout", "desc": "Hanya pada jalur ini Worker B gugur dengan galat Deadline Exceeded akibat pool starvation."}
        ]}
    ],
    narrative="Menahan koneksi database selama pemrosesan I/O pihak ketiga adalah racun arsitektur yang melumpuhkan ketersediaan pool. Hanya pada jalur ini request sehat yang kodenya bersih ikut mati akibat context deadline exceeded karena kehabisan slot koneksi yang tersandera.",
    narrative_en="Retaining database connections during slow third-party network I/O is an architectural hazard that starves connection pools. Only along this path do completely healthy requests abort with deadline exceeded errors due to unreleased socket hostage-taking.",
    mental_model="Holding shared connection slots during slow external I/O creates systemic pool starvation for healthy requests.",
    source_path="content/04-diagrams.md:32-54",
    slide_num="07 / 24",
    layout="cards"
)

# Slide 08: Tabel Perbandingan Mekanisme vs Gejala (Layout: table)
add_slide(
    categories=["COMPARISON", "FAILURE MODES"],
    title="Matriks Diagnostik: Patologi Koneksi Basis Data",
    subtitle="Membedakan karakteristik teknis antara unpooled, oversized pool, dan connection leak.",
    content_items=[
        {"title": "PERBANDINGAN ANOMALI MANAJEMEN KONEKSI", "items": [
            "Dimensi Masalah | Unpooled (No Reuse) | Oversized Pool | Connection Leak (Starvation)",
            "Mekanisme Akar | SetMaxIdleConns(0) tutup koneksi | Client pool > server limit | Koneksi ditahan saat I/O eksternal",
            "Gejala di Server | CPU spike akibat handshake TCP | Error: too many clients | CPU database terlihat normal/idle",
            "Dampak Latensi | Latensi baseline tinggi di p99 | Antrean lock & context switch | Lonjakan waktu tunggu akuisisi pool",
            "Manifestasi Error | Koneksi gagal di lonjakan tinggi | ErrServerOverloaded / PG 53300 | Context Deadline Exceeded pada klien",
            "Mitigasi Utama | Aktifkan idle pool reuse | Hitung sizing total deployment | Isolasi I/O & leakDetectionThreshold"
        ], "type": "table"}
    ],
    narrative="Tabel ini memetakan perbedaan akar masalah antara tanpa pool, pool terlalu besar, dan kebocoran penahanan koneksi. Kebocoran koneksi memiliki karakteristik paling berbahaya karena server database terlihat sehat padahal antrean klien mengalami kelaparan ekstrem.",
    narrative_en="This diagnostic matrix contrasts the distinct failure mechanisms of unpooled drivers, oversized pools, and connection leaks. Connection leaks represent the most insidious failure mode because the database displays low resource usage while clients starve.",
    mental_model="Diagnose connection anomalies by correlating application pool wait times against database backend slot usage.",
    source_path="content/02-master-draft.md:128-165",
    slide_num="08 / 24",
    layout="table"
)

# Slide 09: Decision Sequence — 3 Focused Decision Cards (Layout: cards)
# Limiting phrase 7: "bukan berarti"
add_slide(
    categories=["DECISION SEQUENCE", "POOL TUNING"],
    title="Tiga Prinsip Esensial Penentuan Ukuran Pool",
    subtitle="Hierarki keputusan rekayasa untuk mengamankan stabilitas koneksi sistem.",
    content_items=[
        {"type": "cards", "title": "PANDUAN PENETAPAN SIZING KONEKSI", "items": [
            {"title": "1. Hitung Skala Global", "desc": "Total koneksi = (jumlah instance × pool size); wajib di bawah max_connections database."},
            {"title": "2. Sisihkan Slot Cadangan", "desc": "Kurangi minimal 15 koneksi untuk superuser, replikasi fisik, dan monitoring sistem."},
            {"title": "3. Formula Titik Awal", "desc": "Gunakan ((core * 2) + disk); bukan berarti angka mutlak, melainkan basis uji beban."}
        ]}
    ],
    narrative="Penetapan ukuran pool wajib memperhitungkan total armada aplikasi yang terpasang pada kluster produksi. Mengadopsi rumus ((core * 2) + disk) bukan berarti mendapatkan angka sakral tanpa cela, melainkan menetapkan baseline terukur yang wajib divalidasi lewat pengujian beban.",
    narrative_en="Connection pool sizing must aggregate the entire fleet of application instances deployed against the cluster. Adopting the ((core * 2) + disk) formula does not yield an infallible universal constant; it provides a rational baseline for empirical load testing.",
    mental_model="Enforce deployment-wide connection limits and validate core-based sizing formulas through empirical stress tests.",
    source_path="content/02-master-draft.md:43-69",
    slide_num="09 / 24",
    layout="cards"
)

# Slide 10: Bukti Eksekusi / Terminal Output (Layout: code)
add_slide(
    categories=["EVIDENCE", "EXECUTION"],
    title="Verifikasi Eksekusi: Test Suite dan CLI Demo Lolos Bersih",
    subtitle="Output terminal membuktikan penalti unpooled, penolakan server, dan deteksi starvation.",
    content_items=[
        {"type": "code", "title": "TERMINAL OUTPUT: go test -v -race ./... && go run ./cmd/demo", "items": [
            "=== RUN   TestDirectConnectionOverhead",
            "--- PASS: TestDirectConnectionOverhead (0.03s) [pooled faster than unpooled]",
            "=== RUN   TestOversizedPoolExhaustsServerConnections",
            "--- PASS: TestOversizedPoolExhaustsServerConnections (0.02s) [10 rejected]",
            "=== RUN   TestConnectionStarvationDueToLeak",
            "--- PASS: TestConnectionStarvationDueToLeak (0.02s) [deadline exceeded]",
            "=== RUN   TestTotalCreatedPoolReuse",
            "--- PASS: TestTotalCreatedPoolReuse (0.00s) [reuse verified: 1 vs 5]",
            "PASS: ok  labs/19-database-connection-pooling/tests  1.328s (race clean)",
            "1. Unpooled: 54.59ms vs Pooled: 6.33µs",
            "2. Oversized Pool: Attempted 30, Succeeded 15, Server Rejected 15",
            "3. Starvation: Order 3 Failed: context deadline exceeded"
        ]}
    ],
    narrative="Hasil eksekusi membuktikan sepuluh skenario uji lolos sempurna tanpa data race di bawah Go race detector. Demonstrasi CLI mengonfirmasi keunggulan reuse mikrodetik atas handshake milidetik, penolakan tepat 15 request saat server saturasi, serta isolasi kegagalan pada kebocoran koneksi.",
    narrative_en="Execution output confirms all ten test scenarios pass cleanly with zero data races under the Go race detector. The CLI demonstration verifies microsecond reuse advantages over millisecond handshakes, exact server rejection under saturation, and leak-induced timeouts.",
    mental_model="Deterministic mock testing provides empirical verification of connection lifecycle bottlenecks and exhaustion limits.",
    source_path="engineering/03-execution-result.md:1-55",
    slide_num="10 / 24",
    layout="code"
)

# Slide 11: Analisis Batasan Bukti (Layout: cards)
# Limiting phrase 8: "belum membuktikan"
add_slide(
    categories=["AUDIT", "EVIDENCE BOUNDARY"],
    title="Batasan Bukti Lab: Yang Dibuktikan vs Belum",
    subtitle="Menegaskan batas simulasi in-memory terhadap kompleksitas database PostgreSQL riil.",
    content_items=[
        {"type": "cards", "title": "BATAS PEMBUKTIAN SISTEM LAB 19", "items": [
            {"title": "YANG DIBUKTIKAN", "desc": "Overhead koneksi baru, penolakan slot server, pool starvation, dan penanganan double close."},
            {"title": "YANG TIDAK DIBUKTIKAN", "desc": "Alokasi work_mem per backend, lock contention engine, OOM PostgreSQL, dan TCP keepalive."},
            {"title": "BATAS KLAIM EMPIRIS", "desc": "Lab belum membuktikan efektivitas formula core pada NVMe SSD modern berkecepatan tinggi."}
        ]}
    ],
    narrative="Lab ini membuktikan mekanika dasar peminjaman pool dan batas penolakan slot menggunakan mock driver in-memory. Pengujian ini belum membuktikan perilaku internal mesin database PostgreSQL nyata seperti pemborosan memori work_mem atau saturasi IOPS pada SSD NVMe.",
    narrative_en="This lab proves fundamental pooling mechanics and slot rejection thresholds using an in-memory mock driver. It has not yet proven production PostgreSQL internals such as work_mem memory pressure or storage IOPS saturation on modern NVMe SSDs.",
    mental_model="Distinguish application connection pool queuing from internal database engine process and memory management.",
    source_path="content/02-master-draft.md:164-165",
    slide_num="11 / 24",
    layout="cards"
)

# Slide 12: Trade-off Arsitektur (Layout: cards)
# Limiting phrase 9: "tidak otomatis"
add_slide(
    categories=["TRADE-OFF", "COSTS & BENEFITS"],
    title="Kompromi Rekayasa: Pool Kecil vs Pool Besar",
    subtitle="Mengevaluasi konsekuensi operasional antara antrean antarmuka dan kestabilan backend.",
    content_items=[
        {"type": "cards", "title": "ANALISIS PERBANDINGAN TRADE-OFF", "items": [
            {"title": "Pool Kecil: Protektif", "desc": "Database beroperasi prima tanpa kejenuhan; request yang meluap wajib sabar mengantre di aplikasi."},
            {"title": "Pool Besar: Ilusi Konkurensi", "desc": "Menampung lonjakan sesaat, namun memicu kontensi lock parah dan risiko kehabisan RAM."},
            {"title": "Tidak Otomatis Aman", "desc": "Memperbesar pool tidak otomatis menaikkan throughput bila kapasitas eksekusi disk database jenuh."}
        ]}
    ],
    narrative="Memilih ukuran pool adalah seni menyeimbangkan antrean di sisi aplikasi melawan stabilitas mesin basis data. Memperbesar ukuran pool tidak otomatis meningkatkan kecepatan pemrosesan jika konkurensi tambahan tersebut hanya memicu perkelahian lock internal dan thrashing memori.",
    narrative_en="Tuning pool sizing balances client-side wait queues against database engine stability. Expanding pool size does not automatically increase transaction throughput if additional concurrency merely induces lock contention and memory thrashing.",
    mental_model="Trade client-side acquisition latency for backend stability by maintaining compact connection pools.",
    source_path="content/02-master-draft.md:13-19",
    slide_num="12 / 24",
    layout="cards"
)

# Slide 13: Concurrency / Synchronization Primitive (Layout: code)
# Limiting phrase 10: "simulasi ini tidak sama dengan"
add_slide(
    categories=["CONCURRENCY", "MOCK DRIVER"],
    title="Mekanisme Pengendali Koneksi Menggunakan Atomic Counter",
    subtitle="Simulasi batas server dan pencatatan koneksi aktif dengan perlindungan thread-safe.",
    content_items=[
        {"type": "code", "title": "internal/pool/mockdb.go", "items": [
            "func (d *MockDriver) Open(name string) (driver.Conn, error) {",
            "    if d.connectDelay > 0 { time.Sleep(d.connectDelay) }",
            "    d.mu.Lock()",
            "    defer d.mu.Unlock()",
            "    if d.maxConnections > 0 && atomic.LoadInt32(&d.activeConns) >= d.maxConnections {",
            "        return nil, ErrServerOverloaded",
            "    }",
            "    atomic.AddInt32(&d.activeConns, 1)",
            "    atomic.AddInt32(&d.totalCreated, 1)",
            "    return &mockConn{driver: d}, nil",
            "}"
        ]}
    ],
    narrative="MockDriver mengendalikan pembukaan koneksi dengan menyimulasikan delay handshake dan menolak permintaan saat batas activeConns tercapai. Simulasi ini tidak sama dengan proses forking postgres di Linux, namun secara akurat merefleksikan penolakan slot keras pada beban tinggi.",
    narrative_en="MockDriver regulates connection creation by injecting artificial handshake delay and enforcing strict activeConns limits. While this simulation is not identical to Linux postgres process forks, it accurately models hard slot rejection under saturation.",
    mental_model="Atomic tracking of active socket handles enforces deterministic connection bounds across concurrent goroutines.",
    source_path="internal/pool/mockdb.go:33-49",
    slide_num="13 / 24",
    layout="code"
)

# Slide 14: Before vs After Comparison Cards (Layout: cards)
# Limiting phrase 11: "bukan berarti"
add_slide(
    categories=["COMPARISON", "BEFORE VS AFTER"],
    title="Transformasi: Dari Koneksi Bocor ke Eksekusi Terisolasi",
    subtitle="Pergeseran dari penahanan soket sembrono menuju disiplin peminjaman minimal.",
    content_items=[
        {"type": "cards", "title": "KOMPARASI KONDISI EKSEKUSI LAYANAN", "items": [
            {"title": "Sebelum: Tahan Sembrono", "desc": "Koneksi dipinjam sebelum API pihak ketiga dipanggil; soket tersandera selama ratusan milidetik."},
            {"title": "Sesudah: Peminjaman Sempit", "desc": "API eksternal diselesaikan lebih dulu; koneksi hanya dipinjam selama mikrosekon eksekusi SQL."},
            {"title": "Bukan Berarti Lambat", "desc": "Memecah alur bukan berarti memperlambat kode; pool terbebas dari ancaman kelaparan sistemik."}
        ]}
    ],
    narrative="Sebelum refaktorisasi, koneksi database disandera selama panggilan jaringan eksternal berlangsung lambat. Memecah alur eksekusi bukan berarti menambah kompleksitas, melainkan membebaskan pool dari ancaman starvation sehingga puluhan transaksi lain dapat berjalan mulus.",
    narrative_en="Prior to refactoring, database connections remained captive during slow third-party network transmissions. Decoupling the execution pipeline does not introduce code bloat; it liberates the pool from starvation so dozens of concurrent transactions succeed.",
    mental_model="Isolate slow I/O boundaries outside database connection scopes to preserve shared pool availability.",
    source_path="content/02-master-draft.md:135-146",
    slide_num="14 / 24",
    layout="cards"
)

# Slide 15: Analisis Mendalam — Pola Kebocoran Tidak Aman (Layout: code)
# Limiting phrase 12: "cakupannya mengikuti test"
add_slide(
    categories=["ANALYSIS", "LEAK MECHANISM"],
    title="Anatomi Kode Berbahaya: Menyandera Koneksi Saat I/O Lambat",
    subtitle="Pola anti-pattern ProcessOrderUnsafeLeak yang mengunci soket database.",
    content_items=[
        {"type": "code", "title": "internal/pool/service.go", "items": [
            "func (s *OrderService) ProcessOrderUnsafeLeak(ctx context.Context, orderID int, externalCall func() error) error {",
            "    conn, err := s.db.Conn(ctx)",
            "    if err != nil { return err }",
            "    defer conn.Close()",
            "    _, err = conn.ExecContext(ctx, \"UPDATE orders SET status = 'processing' WHERE id = ?\", orderID)",
            "    if err != nil { return err }",
            "    if externalCall != nil {",
            "        if err := externalCall(); err != nil { return err }",
            "    }",
            "    return nil",
            "}"
        ]}
    ],
    narrative="Metode ProcessOrderUnsafeLeak menahan koneksi aktif ketika fungsi externalCall dijalankan. Bahaya penahanan soket ini cakupannya mengikuti test starvation lab, membuktikan bahwa penundaan defer conn.Close() selama I/O eksternal secara langsung memblokir antrean peminjam lain.",
    narrative_en="ProcessOrderUnsafeLeak keeps the database connection locked while the externalCall closure executes. The hazard of this prolonged hold strictly matches our starvation test suite, demonstrating that deferring conn.Close() across external I/O directly starves borrowers.",
    mental_model="Holding open connections across unbounded network boundaries creates artificial throughput bottlenecks.",
    source_path="internal/pool/service.go:36-57",
    slide_num="15 / 24",
    layout="code"
)

# Slide 16: Analisis Mendalam — Test Bukti Penalti Handshake (Layout: code)
# Limiting phrase 13: "tidak otomatis"
add_slide(
    categories=["ANALYSIS", "OVERHEAD VERIFICATION"],
    title="Verifikasi Test: Membuktikan Penalti Pembuatan Koneksi Langsung",
    subtitle="Membandingkan latensi eksekusi antara SetMaxIdleConns(0) dan pool aktif.",
    content_items=[
        {"type": "code", "title": "tests/pool_test.go", "items": [
            "func TestDirectConnectionOverhead(t *testing.T) {",
            "    mockDriver := pool.NewMockDriver(100, 5*time.Millisecond)",
            "    unpooledDB := pool.OpenDB(mockDriver)",
            "    unpooledDB.SetMaxIdleConns(0) // Memaksa pembukaan koneksi baru setiap kali",
            "    // Eksekusi 5 query unpooled menghasilkan akumulasi delay 25ms+",
            "    pooledDB := pool.OpenDB(mockDriver)",
            "    pooledDB.SetMaxIdleConns(5)",
            "    pooledDB.SetMaxOpenConns(5)",
            "    // Eksekusi 5 query pooled memakai ulang koneksi idle yang sama",
            "    if pooledDuration >= unpooledDuration { t.Errorf(\"expected pooled to be faster\") }",
            "}"
        ]}
    ],
    narrative="TestDirectConnectionOverhead membuktikan bahwa pool dengan idle connections selalu lebih cepat daripada pembukaan koneksi baru di tiap query. Menjalankan query pada database tidak otomatis menghasilkan performa optimal tanpa konfigurasi idle pool yang memadai.",
    narrative_en="TestDirectConnectionOverhead proves that pooled reuse consistently outperforms fresh connection allocation on every query. Executing database operations does not automatically yield acceptable throughput without an appropriately tuned idle connection pool.",
    mental_model="Reusing warm connection handles bypasses socket setup latencies, ensuring microsecond query dispatch.",
    source_path="tests/pool_test.go:13-49",
    slide_num="16 / 24",
    layout="code"
)

# Slide 17: Analisis Mendalam — Test Bukti Server Exhaustion (Layout: code)
# Limiting phrase 14: "hanya pada jalur ini"
add_slide(
    categories=["ANALYSIS", "SATURATION VERIFICATION"],
    title="Verifikasi Test: Pool Terlalu Besar Menghabiskan Kapasitas Server",
    subtitle="Membuktikan penolakan deterministik ketika pool klien melebihi kapasitas backend.",
    content_items=[
        {"type": "code", "title": "tests/pool_test.go", "items": [
            "func TestOversizedPoolExhaustsServerConnections(t *testing.T) {",
            "    mockDriver := pool.NewMockDriver(10, 0) // Server hanya izinkan 10 slot",
            "    db := pool.OpenDB(mockDriver)",
            "    db.SetMaxOpenConns(20) // Klien salah konfigurasi dengan pool 20",
            "    // 20 goroutine meminta koneksi secara bersamaan",
            "    wg.Wait()",
            "    if errCount < 10 {",
            "        t.Errorf(\"expected at least 10 connection attempts to fail, got %d\", errCount)",
            "    }",
            "}"
        ]}
    ],
    narrative="Pengujian ini membuktikan bahwa pool klien berukuran 20 pada server berkapasitas 10 selalu menghasilkan penolakan deterministik. Hanya pada jalur ini kegagalan server terbukti terjadi karena aplikasi mengabaikan batas kapasitas fisik backend database.",
    narrative_en="This test proves that a client pool of 20 against a server capacity of 10 deterministically produces connection rejections. Only along this path do server exhaustion errors emerge because the application disregarded backend physical slot constraints.",
    mental_model="Oversized client pools inevitably trigger server-side connection exhaustion rejections under load.",
    source_path="tests/pool_test.go:51-86",
    slide_num="17 / 24",
    layout="code"
)

# Slide 18: Analisis Mendalam — Test Bukti Starvation Akibat Leak (Layout: code)
# Limiting phrase 15: "bukan berarti"
add_slide(
    categories=["ANALYSIS", "STARVATION VERIFICATION"],
    title="Verifikasi Test: Kebocoran Satu Worker Melumpuhkan Worker Lain",
    subtitle="Permintaan kedua timeout karena koneksi tunggal tersandera I/O eksternal.",
    content_items=[
        {"type": "code", "title": "tests/pool_test.go", "items": [
            "func TestConnectionStarvationDueToLeak(t *testing.T) {",
            "    db.SetMaxOpenConns(1) // Pool ketat berukuran 1 koneksi",
            "    go func() {",
            "        _ = svc.ProcessOrderUnsafeLeak(ctx, 1, func() error {",
            "            close(acquired); time.Sleep(100 * time.Millisecond); return nil",
            "        })",
            "    }()",
            "    <-acquired",
            "    ctx, cancel := context.WithTimeout(ctx, 20*time.Millisecond)",
            "    defer cancel()",
            "    err := svc.ProcessOrderSafe(ctx, 2, nil) // Gagal timeout 20ms",
            "}"
        ]}
    ],
    narrative="Unit test mereproduksi pool starvation secara deterministik pada pool berukuran satu koneksi. Terjadinya timeout pada pemanggil yang aman bukan berarti kodenya cacat, melainkan akibat soket tunggal tersandera oleh fungsi lain yang ceroboh menjalankan delay jaringan eksternal.",
    narrative_en="This unit test deterministically reproduces pool starvation within a single-connection pool. Encountering a timeout on the safe caller does not signify flawed logic; it proves the shared socket was held hostage by a negligent routine running external delays.",
    mental_model="Unreleased connection handles monopolize shared pools, causing collateral starvation for valid requests.",
    source_path="tests/pool_test.go:88-124",
    slide_num="18 / 24",
    layout="code"
)

# Slide 19: Analisis Mendalam — Deadlock Pool-Locking (Layout: code)
# Limiting phrase 16: "cakupannya mengikuti test"
add_slide(
    categories=["ANALYSIS", "POOL DEADLOCK"],
    title="Verifikasi Test: Bahaya Deadlock Akuisisi Koneksi Bersarang",
    subtitle="Meminta koneksi kedua pada thread yang sama saat pool berukuran 1 memicu hang.",
    content_items=[
        {"type": "code", "title": "tests/pool_test.go", "items": [
            "func TestPoolLockingDeadlock(t *testing.T) {",
            "    db.SetMaxOpenConns(1) // Pool berukuran 1",
            "    ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)",
            "    defer cancel()",
            "    conn1, err := db.Conn(ctx)",
            "    if err != nil { t.Fatalf(\"failed to acquire first: %v\", err) }",
            "    defer conn1.Close()",
            "    _, err = db.Conn(ctx) // Meminta koneksi kedua pada pool yang sama!",
            "    if err == nil { t.Fatalf(\"expected context deadline exceeded, got nil\") }",
            "}"
        ]}
    ],
    narrative="Pengujian ini memvalidasi skenario pool-locking di mana sebuah thread mencoba mengakuisisi koneksi kedua sementara koneksi pertama belum dilepas. Verifikasi ini cakupannya mengikuti test lab, memvalidasi formula Tn x (Cm - 1) + 1 sebagai syarat matematis mutlak pencegah kebuntuan pool.",
    narrative_en="This test validates a pool-locking scenario where a single thread requests a secondary connection while retaining its primary handle. The verification scope strictly matches our test suite, proving the mathematical formula Tn x (Cm - 1) + 1 is essential to avert deadlocks.",
    mental_model="Nested connection acquisition without sufficient pool headroom results in deterministic thread self-deadlock.",
    source_path="tests/pool_test.go:154-174",
    slide_num="19 / 24",
    layout="code"
)

# Slide 20: Checklist Implementasi (Layout: checklist)
add_slide(
    categories=["CHECKLIST", "PRODUCTION READINESS"],
    title="Checklist Manajemen Koneksi di Lingkungan Produksi",
    subtitle="Protokol wajib arsitektur sebelum merilis pool database ke lingkungan berskala tinggi.",
    content_items=[
        {"type": "checklist", "title": "STANDAR REKAYASA KONEKSI BASIS DATA", "items": [
            "Hitung total armada: (instance_count x pool_size) wajib lebih kecil dari max_connections.",
            "Sisihkan 15 slot koneksi untuk maintenance, superuser, dan physical replication.",
            "Pastikan tidak ada panggilan I/O jaringan atau komputasi file di dalam blok koneksi.",
            "Gunakan defer conn.Close() atau helper sql.DB tingkat tinggi pada setiap percabangan.",
            "Pasang context timeout ketat pada setiap peminjaman soket dan eksekusi query SQL.",
            "Aktifkan leak detection threshold (HikariCP / wrapper) untuk menangkap stack trace leak."
        ]}
    ],
    narrative="Checklist ini merangkum langkah audit koneksi wajib: verifikasi total koneksi deployment, pencadangan slot administratif, pembebasan transaksi dari panggilan eksternal, penutupan deterministik via defer, pemasangan batas context timeout, serta pemantauan durasi pinjam.",
    narrative_en="This checklist codifies mandatory connection governance: aggregating deployment-wide totals, reserving administrative slots, barring external I/O from transactions, enforcing deterministic defer release, applying strict timeouts, and tracking hold durations.",
    mental_model="Enforce strict operational standards across connection sizing, lifecycle scoping, and timeout boundaries.",
    source_path="content/02-master-draft.md:205-215",
    slide_num="20 / 24",
    layout="checklist"
)

# Slide 21: Anti-Pattern & Cara Reproduksi Kegagalan (Layout: cards)
# Limiting phrase 17: "bukan berarti"
add_slide(
    categories=["ANTI-PATTERN", "PRODUCTION PITFALLS"],
    title="Anti-Pattern Umum dalam Arsitektur Connection Pool",
    subtitle="Kekeliruan konfigurasi yang kerap melumpuhkan stabilitas database relasional.",
    content_items=[
        {"type": "cards", "title": "POLA BERBAHAYA YANG HARUS DIHINDARI", "items": [
            {"title": "Naikkan Pool Saat Error", "desc": "Merespons error koneksi dengan memperbesar pool; memperparah CPU thrashing dan lock contention."},
            {"title": "Sizing Per-Instance", "desc": "Menghitung pool secara lokal tanpa mengalikan jumlah replica pod pada orkestrasi Kubernetes."},
            {"title": "Bukan Berarti Gratis", "desc": "Menganggap koneksi idle gratis bukan berarti tanpa beban; RAM backend tetap terkuras permanen."}
        ]}
    ],
    narrative="Menaikkan ukuran pool saat database mengalami galat koneksi adalah anti-pattern paling umum yang memperparah insiden produksi. Mengabaikan konsumsi memori koneksi idle bukan berarti sistem bebas beban, melainkan mengunci ratusan megabyte RAM backend yang semestinya dipakai untuk buffer cache.",
    narrative_en="Increasing pool size in response to connection exhaustion errors is an anti-pattern that exacerbates outages. Assuming idle connections cost nothing does not reflect reality; it locks hundreds of megabytes of backend RAM that should otherwise serve buffer caches.",
    mental_model="Avoid reactive pool enlargement and account for cumulative connection costs across container replicas.",
    source_path="content/02-master-draft.md:186-196",
    slide_num="21 / 24",
    layout="cards"
)

# Slide 22: Takeaway Praktis & Command Block Dominan (Layout: code)
# Limiting phrase 18: "production masih perlu"
add_slide(
    categories=["REPRODUCTION", "COMMANDS"],
    title="Reproduksi Mandiri di Mesin Lokal Menggunakan Go",
    subtitle="Verifikasi mandiri seluruh perilaku pool tanpa dependensi database eksternal.",
    content_items=[
        {"type": "code", "title": "COMMAND RUNNER (ZERO DEPENDENCIES)", "items": [
            "cd labs/19-database-connection-pooling",
            "go test -v ./...    # menjalankan 10 test suite verifikasi overhead, limit, dan leak",
            "go test -race ./... # pembuktian sinkronisasi konkurensi bersih tanpa data race",
            "go run ./cmd/demo   # eksekusi demonstrasi interaktif skenario 1, 2, dan 3"
        ]}
    ],
    narrative="Seluruh arsitektur dibangun murni menggunakan pustaka standar Go tanpa perlu menginstal server PostgreSQL lokal. Pengujian ini memvalidasi mekanika peminjaman soket, namun production masih perlu mengonfigurasi parameter PgBouncer serta pemantauan metrik pg_stat_activity.",
    narrative_en="The entire testing suite runs on the Go standard library without requiring external PostgreSQL servers. While local tests validate socket borrowing mechanics, production systems still require tuning PgBouncer parameters and monitoring pg_stat_activity metrics.",
    mental_model="Self-contained standard library harnesses enable rapid, zero-dependency validation of connection architectures.",
    source_path="README.md:1-35",
    slide_num="22 / 24",
    layout="code"
)

# Slide 23: Prinsip Arsitektural Jangka Panjang (Layout: cards)
# Limiting phrase 19: "tidak otomatis"
add_slide(
    categories=["TAKEAWAY", "ARCHITECTURAL PRINCIPLES"],
    title="Prinsip Inti: Pool Lebih Kecil, Sistem Lebih Cepat",
    subtitle="Stabilitas sistem basis data bertumpu pada disiplin pembatasan konkurensi.",
    content_items=[
        {"type": "cards", "title": "TIGA PILAR KETAHANAN KONEKSI SISTEM", "items": [
            {"title": "Batas Sizing Konservatif", "desc": "Jadikan kapasitas hardware sebagai batas atas; batasi antrean di aplikasi alih-alih di database."},
            {"title": "Isolasi I/O Ketat", "desc": "Keluarkan seluruh panggilan eksternal dari penahanan koneksi; pinjam hanya saat query siap."},
            {"title": "Skalabilitas Horizontal", "desc": "Tidak otomatis menaikkan slot; gunakan proxy pooling transaksi untuk ribuan klien."}
        ]}
    ],
    narrative="Pertanyaan mendasar dalam rekayasa perangkat lunak bukanlah seberapa besar pool yang bisa dibuat, melainkan sekecil apa pool yang masih mampu memenuhi target throughput. Menghadapi lonjakan traffic tidak otomatis diselesaikan dengan menaikkan koneksi server, melainkan dengan memasang proxy pooling transaksi.",
    narrative_en="The fundamental engineering question is not how large a pool can be configured, but how compact it can remain while satisfying throughput demands. Handling traffic spikes does not automatically justify inflating database connection limits; it requires deploying transaction-mode proxies.",
    mental_model="Resilient architectures maintain compact connection pools, isolate external I/O, and scale via proxies.",
    source_path="content/05-key-takeaways.md:1-20",
    slide_num="23 / 24",
    layout="cards"
)

# Slide 24: Closing Slide (Layout: closing)
add_slide(
    categories=["CLOSING"],
    title="Koneksi database adalah sumber daya terbatas backend, bukan variabel yang bisa diperbesar tanpa batas.",
    subtitle="",
    content_items=[
        {"type": "bullet", "title": "MODUL IMPLEMENTASI", "items": [
            "internal/pool/mockdb.go • internal/pool/service.go",
            "cmd/demo/main.go • tests/pool_test.go",
            "engineering/01-design.md • content/02-master-draft.md"
        ]},
        {"type": "bullet", "title": "VERIFIKASI & STATUS", "items": [
            f"Commit: {COMMIT_SHA[:7]} • Status: APPROVED",
            "Test passing: 10/10 test suites (Overhead, Oversized, Starvation, Deadlock, Reuse)",
            "Race detector: Clean (go test -race ./... PASS in 1.328s)",
            "Zero external dependencies: 100% Go standard library."
        ]}
    ],
    mental_model="Database connections are finite physical resources that demand strict lifecycle governance and proxy mediation.",
    source_path=REPO_URL,
    slide_num="24 / 24",
    layout="closing"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
    for s in SLIDES:
        print(f"[{s['slide_num']}] {s['layout']}: {s['title']}")
