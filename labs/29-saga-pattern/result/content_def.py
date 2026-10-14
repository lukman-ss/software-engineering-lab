#!/usr/bin/env python3
"""
Publication Definition for Software Engineering Lab 29: Saga Pattern.
Defines 24 slides, metadata, color constants, and content structure.
DO NOT DELETE THIS FILE.
"""
import os
import sys

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/29-saga-pattern/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/29-saga-pattern"
COMMIT_SHA = "9c525944dc2de00e16513830a25ffa1c4834d742"
LAB_NO = "29"
SLUG = "saga-pattern"
TITLE_ID = "Saga Pattern: Transaksi Terdistribusi Tanpa 2PC"
SUBTITLE_ID = "Eventual Consistency melalui LIFO Compensation & Semantic Locks"
TITLE_EN = "Saga Pattern: Distributed Transactions Without 2PC"
SUBTITLE_EN = "Eventual Consistency via LIFO Compensation & Semantic Locks"
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
            "Retro CRT monitor beige dengan scanlines",
            "ASCII art terminal yang menampilkan alur saga execution",
            f"Penulis: {AUTHOR}",
        ]}
    ],
    mental_model="",
    source_path="",
    slide_num=""
)

# Slide 2: Problem Definition (Database-per-Service & Failure of 2PC)
add_slide(
    categories=["MOTIVASI", "KONSEP"],
    title="Keterbatasan ACID Monolitik di Dunia Microservices",
    subtitle="Database-per-Service memecah transaksi lintas boundary layanan mandiri",
    content_items=[
        {"title": "MASALAH DASAR", "items": [
            "Arsitektur microservices menuntut setiap service memiliki database sendiri",
            "Transaksi bisnis seperti Checkout membutuhkan koordinasi Order, Payment, dan Inventory",
            "Pola BEGIN TRANSACTION...COMMIT tunggal tidak dapat melintasi batas jaringan database",
        ], "type": "bullet"},
        {"title": "KEGAGALAN TWO-PHASE COMMIT (2PC)", "items": [
            "2PC mengunci baris data di seluruh node database selama fase persiapan (prepare)",
            "Koordinator tunggal menjadi single point of failure (SPOF) dan latency multiplier",
            "Partisi jaringan membuat resource terkunci tanpa batas (blocking protocol)",
        ], "type": "bullet"},
        {"title": "PERGESERAN PARADIGMA", "items": [
            "Sistem beralih dari ACID ketat menuju BASE (Basically Available, Soft state, Eventual consistency)",
            "Saga pattern memecah transaksi global menjadi rangkaian transaksi lokal independen",
            "Setiap transaksi lokal langsung melakukan commit di databasenya masing-masing",
        ], "type": "bullet"},
    ],
    mental_model="2PC locks distributed databases synchronously; Sagas commit locally and coordinate asynchronously.",
    source_path="labs/29-saga-pattern/README.md",
    slide_num="2 / 24"
)

# Slide 3: Anatomy of a Saga
add_slide(
    categories=["ANATOMI", "ALUR EKSEKUSI"],
    title="Anatomi Saga: Rangkaian Transaksi Lokal",
    subtitle="Eksekusi maju teratur dan pembatalan mundur terarah",
    content_items=[
        {"title": "FORWARD EXECUTION", "items": [
            "T1 (Order Created) -> T2 (Payment Processed) -> T3 (Inventory Reserved) -> T4 (Order Approved)",
            "Setiap step mengeksekusi mutasi data lokal dan mencatat status EXECUTED",
            "Jika semua langkah berhasil, transaksi terdistribusi selesai secara konsisten",
        ], "type": "diagram"},
        {"title": "COMPENSATION ROLLBACK (LIFO)", "items": [
            "Jika T3 gagal (misal out-of-stock), orchestrator memicu kompensasi",
            "Urutan mundur: C2 (Refund Payment) -> C1 (Cancel Order)",
            "Prinsip tumpukan LIFO menjamin efek dibatalkan dalam urutan terbalik",
        ], "type": "diagram"},
        {"title": "KARAKTERISTIK KUNCI", "items": [
            "Kompensasi adalah transaksi bisnis baru, bukan pembatalan level storage engine",
            "Setiap langkah menghasilkan log audit permanen di sistem",
            "Rollback tidak menghapus record masa lalu melainkan menuliskan saldo koreksi",
        ], "type": "bullet"},
    ],
    mental_model="A saga is a chain of local commits; failure triggers a reverse chain of business compensations.",
    source_path="labs/29-saga-pattern/internal/saga/orchestrator.go:49",
    slide_num="3 / 24"
)

# Slide 4: Transaction Taxonomy
add_slide(
    categories=["TAKSONOMI", "KLASIFIKASI"],
    title="Tiga Jenis Transaksi dalam Saga",
    subtitle="Klasifikasi formal langkah transaksi berdasarkan kemampuan kompensasi",
    content_items=[
        {"title": "1. COMPENSABLE TRANSACTIONS", "items": [
            "Langkah yang dapat dibatalkan atau dikoreksi dengan aksi berlawanan",
            "Contoh: CreateOrder (dibatalkan oleh CancelOrder), ChargeCard (RefundPayment)",
            "Wajib memiliki pasangan kompensasi eksplisit sebelum titik balik tercapai",
        ], "type": "bullet"},
        {"title": "2. PIVOT TRANSACTION", "items": [
            "Titik balik tanpa kembali (point of no return) dalam proses bisnis saga",
            "Jika transaksi pivot berhasil, saga dijamin akan terus berjalan hingga tuntas",
            "Jika transaksi pivot gagal, seluruh langkah compensable sebelumnya harus di-rollback",
        ], "type": "bullet"},
        {"title": "3. RETRYABLE TRANSACTIONS", "items": [
            "Langkah yang dieksekusi setelah pivot dan dijamin pada akhirnya akan berhasil",
            "Contoh: Kirim notifikasi email, approval final order, pembaruan analitik",
            "Wajib bersifat idempotent karena dapat di-retry berulang kali oleh worker",
        ], "type": "bullet"},
    ],
    mental_model="Compensable before pivot, pivot decides the fate, retryable guarantees completion after pivot.",
    source_path="labs/29-saga-pattern/content/02-master-draft.md:23",
    slide_num="4 / 24"
)

# Slide 5: Orchestration vs Choreography
add_slide(
    categories=["ARSITEKTUR", "PERBANDINGAN"],
    title="Dua Pola Koordinasi: Orchestration vs Choreography",
    subtitle="Sentralisasi kendali alur versus desentralisasi reaktif berbasis event",
    content_items=[
        {"title": "ORCHESTRATION (CENTRALIZED)", "items": [
            "Koordinator sentral memegang state machine dan memanggil service satu per satu",
            "Alur bisnis eksplisit, mudah di-audit, dan penanganan LIFO rollback terpusat",
            "Kelemahan: Risiko kopling logika bisnis ke satu koordinator jika tidak dirawat",
        ], "type": "bullet"},
        {"title": "CHOREOGRAPHY (EVENT-DRIVEN)", "items": [
            "Setiap service mempublikasikan domain event dan bereaksi terhadap event service lain",
            "Tingkat decoupling tinggi tanpa orchestrator tunggal di tengah",
            "Kelemahan: Alur implisit, sulit melacak state global, risiko siklus event tak berujung",
        ], "type": "bullet"},
        {"title": "PANDUAN PEMILIHAN", "items": [
            "Gunakan Orchestration untuk workflow kompleks dengan >3 langkah dan kompensasi rumit",
            "Gunakan Choreography untuk integrasi sederhana 2-3 service yang minim ketergantungan balik",
            "Lab 29 mengimplementasikan kedua model untuk perbandingan langsung yang objektif",
        ], "type": "bullet"},
    ],
    mental_model="Orchestration gives explicit central control; Choreography gives loose event-driven autonomy.",
    source_path="labs/29-saga-pattern/internal/saga/choreography.go:27",
    slide_num="5 / 24"
)

# Slide 6: System Architecture of Lab 29
add_slide(
    categories=["DESAIN SISTEM"],
    title="Arsitektur Lab 29: Komponen & Boundary",
    subtitle="Pemisahan domain services, state machine orchestrator, dan decoupling event bus",
    content_items=[
        {"title": "SAGA CORE (`internal/saga`)", "items": [
            "Orchestrator: Mengelola slice Step, eksekusi sekuensial, dan kompensasi LIFO",
            "StepLog: Mencatat status PENDING, EXECUTED, FAILED, COMPENSATED, COMPENSATE_FAILED",
            "EventBus: Memfasilitasi publish-subscribe berbasis thread-safe sync.RWMutex",
        ], "type": "bullet"},
        {"title": "DOMAIN SERVICES (`internal/services`)", "items": [
            "OrderService: Mengelola status pesanan dan semantic lock untuk pencegahan race condition",
            "PaymentService: Memproses pembayaran dengan idempotency key per payment ID",
            "InventoryService: Reservasi dan pelepasan stok inventaris barang",
        ], "type": "bullet"},
        {"title": "VERIFIKASI & RUNNER", "items": [
            "cmd/demo/main.go: Mendemonstrasikan Happy Path dan Rollback Out-of-Stock di CLI",
            "tests/saga_test.go: 9 skenario pengujian komprehensif termasuk concurrency race test",
            "Zero external dependencies — murni menggunakan Go standard library untuk akurasi tertinggi",
        ], "type": "bullet"},
    ],
    mental_model="Clean component separation isolates domain logic from distributed transaction coordination.",
    source_path="labs/29-saga-pattern/README.md:5",
    slide_num="6 / 24"
)

# Slide 7: Orchestrator Engine Implementation
add_slide(
    categories=["BUKTI KODE", "ORCHESTRATOR"],
    title="Implementasi Orchestrator & Struktur Step",
    subtitle="Fungsi Execute dan Compensate sebagai kontrak utama setiap tahapan",
    content_items=[
        {"title": "STRUKTUR STEP & STATUS", "items": [
            "type Step struct {",
            "    Name       string",
            "    Execute    func(ctx context.Context) error",
            "    Compensate func(ctx context.Context) error",
            "}",
            "Tracking status: PENDING, EXECUTED, FAILED, COMPENSATED, COMPENSATE_FAILED",
        ], "type": "code"},
        {"title": "THREAD-SAFE STATE LOGGING", "items": [
            "type Orchestrator struct {",
            "    mu    sync.Mutex",
            "    steps []Step",
            "    logs  []StepLog",
            "}",
            "Setiap mutasi log dilindungi mutex untuk mendukung akses paralel",
        ], "type": "code"},
        {"title": "KONTRAK KOMPENSASI", "items": [
            "Compensate bersifat opsional untuk langkah read-only atau pivot murni",
            "Jika Compensate bernilai nil, langkah tersebut dilewati tanpa error saat rollback",
            "Menjamin fleksibilitas dalam memodelkan ragam jenis transaksi bisnis",
        ], "type": "bullet"},
    ],
    mental_model="Every orchestrator step pairs a forward action with an optional reverse compensation.",
    source_path="labs/29-saga-pattern/internal/saga/orchestrator.go:19",
    slide_num="7 / 24"
)

# Slide 8: LIFO Compensation Algorithm
add_slide(
    categories=["BUKTI KODE", "ALGORITMA"],
    title="Algoritma Kompensasi LIFO: Membalik Riwayat",
    subtitle="Eksekusi mundur dari langkah terakhir yang sukses menuju langkah pertama",
    content_items=[
        {"title": "KODE ALGORITMA LIFO", "items": [
            "func (o *Orchestrator) compensate(ctx context.Context, executed []Step) error {",
            "    var compErrors []error",
            "    for i := len(executed) - 1; i >= 0; i-- {",
            "        step := executed[i]",
            "        if step.Compensate != nil {",
            "            err := step.Compensate(ctx)",
            "            // append error or mark COMPENSATED...",
            "        }",
            "    }",
            "    return aggregateErrors(compErrors)",
            "}",
        ], "type": "code"},
        {"title": "MENGAPA HARUS LIFO?", "items": [
            "Langkah belakangan seringkali bergantung pada status yang dibuat langkah sebelumnya",
            "Contoh: Dana harus dikembalikan (refund) sebelum pesanan dibatalkan (cancel)",
            "Eksekusi acak atau paralel berisiko melanggar batasan konsistensi domain",
        ], "type": "bullet"},
        {"title": "PENANGANAN CONTEXT SAAT ROLLBACK", "items": [
            "Kompensasi dipanggil menggunakan context.Background()",
            "Tujuannya mencegah pembatalan context request induk mematikan proses pembersihan data",
            "Menjamin proses rollback tetap berjalan meski client melakukan disconnect",
        ], "type": "bullet"},
    ],
    mental_model="Unwind executed steps in reverse order (LIFO) so dependent actions undo gracefully.",
    source_path="labs/29-saga-pattern/internal/saga/orchestrator.go:92",
    slide_num="8 / 24"
)

# Slide 9: Choreography Event Bus
add_slide(
    categories=["BUKTI KODE", "CHOREOGRAPHY"],
    title="Event-Driven Choreography via EventBus",
    subtitle="Koordinasi desentralisasi melalui publish/subscribe antar domain service",
    content_items=[
        {"title": "STRUKTUR EVENT BUS", "items": [
            "type EventBus struct {",
            "    mu       sync.RWMutex",
            "    handlers map[EventType][]EventHandler",
            "}",
            "Mendukung pendaftaran banyak subscriber untuk satu event type",
        ], "type": "code"},
        {"title": "ALUR EVENT CHOREOGRAPHY", "items": [
            "OrderCreated -> diproses PaymentService -> PaymentCompleted dipublish",
            "PaymentCompleted -> diproses InventoryService -> InventoryReserved dipublish",
            "InventoryReserved -> diproses OrderService -> OrderApproved dipublish",
        ], "type": "diagram"},
        {"title": "PENANGANAN EVENT ROLLBACK", "items": [
            "Jika reservasi inventaris gagal -> InventoryFailed dipublish",
            "PaymentService mendengar InventoryFailed -> melakukan refund pembayaran",
            "OrderService mendengar PaymentFailed/InventoryFailed -> membatalkan pesanan",
        ], "type": "bullet"},
    ],
    mental_model="Choreography lets services react to domain events without knowing who triggered them.",
    source_path="labs/29-saga-pattern/internal/saga/choreography.go:27",
    slide_num="9 / 24"
)

# Slide 10: Semantic Locking Countermeasure
add_slide(
    categories=["KEAMANAN DATA", "ISOLATION"],
    title="Semantic Locking: Melawan Anomali Tanpa Isolasi ACID",
    subtitle="Menandai record sebagai PENDING untuk mencegah modifikasi konkuren yang merusak",
    content_items=[
        {"title": "MASALAH DIRTY READS DI SAGA", "items": [
            "Karena transaksi lokal langsung di-commit, transaksi lain dapat membaca data sementara",
            "Jika saga kemudian di-rollback, data yang dibaca transaksi lain menjadi invalid",
            "Saga tidak memiliki isolasi level database bawaan antar service",
        ], "type": "bullet"},
        {"title": "IMPLEMENTASI SEMANTIC LOCK", "items": [
            "if s.locks[orderID] {",
            "    return errors.New(\"order is semantically locked\")",
            "}",
            "s.orders[orderID] = OrderPending",
            "s.locks[orderID] = true",
        ], "type": "code"},
        {"title": "PELEPASAN KUNCI SEMANTIK", "items": [
            "Kunci dilepas saat saga mencapai status terminal (ApproveOrder atau CancelOrder)",
            "Request lain yang mencoba memodifikasi pesanan selama pending akan ditolak",
            "Mencegah overwriting data yang sedang dalam proses verifikasi pembayaran",
        ], "type": "bullet"},
    ],
    mental_model="Semantic locks flag uncommitted business state to prevent concurrent dirty modifications.",
    source_path="labs/29-saga-pattern/internal/services/services.go:33",
    slide_num="10 / 24"
)

# Slide 11: Idempotency Keys in Distributed Steps
add_slide(
    categories=["RELIABILITAS", "IDEMPOTENCY"],
    title="Idempotency Keys: Fondasi Retry yang Aman",
    subtitle="Menjamin pemanggilan ulang akibat network glitch tidak menduplikasi mutasi dana",
    content_items=[
        {"title": "TANTANGAN NETWORK TIMEOUT", "items": [
            "Saat memanggil PaymentService, timeout jaringan dapat terjadi setelah dana terpotong",
            "Orchestrator tidak tahu apakah request berhasil atau gagal sehingga wajib melakukan retry",
            "Tanpa proteksi, retry akan memotong saldo pelanggan dua kali",
        ], "type": "bullet"},
        {"title": "KODE PENGECEKAN IDEMPOTENSI", "items": [
            "func (s *PaymentService) ProcessPayment(id string, amt float64, ...) error {",
            "    s.mu.Lock(); defer s.mu.Unlock()",
            "    if s.processedID[id] {",
            "        return nil // Idempotent success: lewati eksekusi",
            "    }",
            "    // proses pembayaran baru...",
            "}",
        ], "type": "code"},
        {"title": "APLIKASI DI SEMUA TRANSAKSI", "items": [
            "Kunci idempotensi wajib disertakan pada semua Compensable dan Retryable transactions",
            "Memungkinkan worker mencoba kembali berkali-kali hingga mencapai kepastian status",
            "Menjaga konsistensi data finansial di tengah ketidakpastian jaringan terdistribusi",
        ], "type": "bullet"},
    ],
    mental_model="Idempotency turns repeated network retries into safe no-ops without double execution.",
    source_path="labs/29-saga-pattern/internal/services/services.go:86",
    slide_num="11 / 24"
)

# Slide 12: Concurrency & Context Cancellation
add_slide(
    categories=["KONKURENSI", "RESILIENSI"],
    title="Manajemen Konkurensi & Pembatalan Konteks",
    subtitle="Deteksi ctx.Done() dan perlindungan terhadap race condition",
    content_items=[
        {"title": "DETEKSI CANCEL DI SETIAP LANGKAH", "items": [
            "select {",
            "case <-ctx.Done():",
            "    o.logs = append(o.logs, StepLog{Name: step.Name, Status: StatusFailed})",
            "    compErr := o.compensate(context.Background(), executed)",
            "    return fmt.Errorf(\"saga cancelled: %w\", ctx.Err())",
            "default:",
            "}",
        ], "type": "code"},
        {"title": "ISOLASI MEMORI VIA MUTEX", "items": [
            "Setiap service (Order, Payment, Inventory) memiliki sync.Mutex mandiri",
            "Saga Orchestrator memproteksi slice log dan steps dari konkurensi",
            "Menjamin integritas internal state saat diakses oleh banyak goroutine",
        ], "type": "bullet"},
        {"title": "HASIL VERIFIKASI -RACE", "items": [
            "10 transaksi saga dijalankan bersamaan dalam goroutine terpisah",
            "Setiap saga mengelola pesanan unik dengan lock dan state lokalnya sendiri",
            "Lulus pengujian go test -race tanpa terdeteksi satu pun data race condition",
        ], "type": "bullet"},
    ],
    mental_model="Check context before every step and execute compensations with an uncancelled context.",
    source_path="labs/29-saga-pattern/internal/saga/orchestrator.go:58",
    slide_num="12 / 24"
)

# Slide 13: Demo Walkthrough: Happy Path
add_slide(
    categories=["DEMO RUNNER", "HASIL NYATA"],
    title="Eksekusi Demo: Skenario Happy Path",
    subtitle="Seluruh langkah bisnis tuntas dan disetujui tanpa hambatan",
    content_items=[
        {"title": "OUTPUT TERMINAL HAPPY PATH", "items": [
            "=== Scenario 1: Happy Path (Orchestrator) ===",
            " -> [OrderService] Creating order: ord-success",
            " -> [PaymentService] Processing payment for: ord-success",
            " -> [InventoryService] Reserving 1 laptop",
            " -> [OrderService] Finalizing order approval",
            "Scenario 1 Result: error=<nil>, OrderState=APPROVED, Remaining Stock=0",
        ], "type": "code"},
        {"title": "VERIFIKASI STATE AKHIR", "items": [
            "OrderService mencatat orderState = APPROVED",
            "PaymentService menyimpan record transaksi senilai 1500.0",
            "InventoryService mendegradasi stok laptop dari 1 menjadi 0",
        ], "type": "bullet"},
        {"title": "STATUS LOG ORCHESTRATOR", "items": [
            "Seluruh 4 langkah mencatat status EXECUTED",
            "Tidak ada fungsi kompensasi yang dipanggil selama proses berjalan",
            "Sistem mencapai eventual consistency dalam hitungan milidetik",
        ], "type": "bullet"},
    ],
    mental_model="When all steps succeed sequentially, the saga completes and leaves domain states committed.",
    source_path="labs/29-saga-pattern/cmd/demo/main.go:20",
    slide_num="13 / 24"
)

# Slide 14: Demo Walkthrough: Failure & Rollback
add_slide(
    categories=["DEMO RUNNER", "ROLLBACK"],
    title="Eksekusi Demo: Skenario Kegagalan Stok & Rollback",
    subtitle="Out-of-stock memicu pembatalan pembayaran dan pesanan via LIFO",
    content_items=[
        {"title": "OUTPUT TERMINAL ROLLBACK", "items": [
            "=== Scenario 2: Rollback on Failure (Orchestrator) ===",
            " -> [OrderService] Creating order: ord-failed",
            " -> [PaymentService] Processing payment for: ord-failed",
            " -> [InventoryService] Reserving 1 laptop (current stock: 0)",
            " <- [PaymentService] Compensating: Refunding payment for: ord-failed",
            " <- [OrderService] Compensating: Cancelling order: ord-failed",
            "Scenario 2 Result: error=step ReserveInventory failed: out of stock",
            "Final State: OrderState=CANCELLED, HasPayment=false, Stock=0",
        ], "type": "code"},
        {"title": "ANALISIS URUTAN KOMPENSASI", "items": [
            "Langkah ke-3 (ReserveInventory) gagal karena stok bernilai 0",
            "Orchestrator memutar balik: RefundPayment dipanggil terlebih dahulu",
            "Setelah refund selesai, CancelOrder dipanggil untuk menutup pesanan",
        ], "type": "bullet"},
        {"title": "INTEGRITAS BISNIS TERJAGA", "items": [
            "Pelanggan tidak dirugikan karena uang berhasil di-refund otomatis",
            "Status pesanan tercatat CANCELLED bukan hilang dari database",
            "Konsistensi bisnis pulih tanpa memerlukan kunci database global",
        ], "type": "bullet"},
    ],
    mental_model="A failed step immediately halts forward execution and triggers exact LIFO compensation.",
    source_path="labs/29-saga-pattern/cmd/demo/main.go:70",
    slide_num="14 / 24"
)

# Slide 15: Quantitative Test Verification
add_slide(
    categories=["VERIFIKASI TEST", "EVIDENCE"],
    title="Bukti Pengujian Otomatis: 9 Test Suites",
    subtitle="Verifikasi formal happy path, kompensasi LIFO, idempotensi, dan error propagation",
    content_items=[
        {"title": "HASIL EKSEKUSI TEST", "items": [
            "=== RUN   TestOrchestrator_HappyPath             --- PASS (0.00s)",
            "=== RUN   TestOrchestrator_FailureCompensatesLIFO--- PASS (0.00s)",
            "=== RUN   TestPayment_Idempotency                --- PASS (0.00s)",
            "=== RUN   TestSemanticLock                       --- PASS (0.00s)",
            "=== RUN   TestOrchestrator_Concurrency           --- PASS (0.00s)",
            "=== RUN   TestChoreography_Flow                  --- PASS (0.00s)",
            "=== RUN   TestChoreography_FailureCompensates    --- PASS (0.00s)",
            "=== RUN   TestOrchestrator_CompensationErrorPropagated --- PASS (0.00s)",
            "=== RUN   TestOrchestrator_ContextCancellation   --- PASS (0.00s)",
            "PASS: 9/9 test suites passed in 0.428s",
        ], "type": "code"},
        {"title": "PENGUJIAN SPESIFIK", "items": [
            "TestSemanticLock: Memastikan pencegahan order creation duplikat saat status pending",
            "TestCompensationErrorPropagated: Memverifikasi penanganan jika refund payment gagal",
            "TestChoreography_Flow: Memvalidasi rantai event pub-sub tanpa koordinator sentral",
        ], "type": "bullet"},
        {"title": "JAMINAN KUALITAS KODE", "items": [
            "Seluruh pengujian lolos uji deteksi data race (`go test -race ./...`)",
            "Test case mencakup skenario happy path, sad path, dan edge cases kegagalan jaringan",
            "Bukti deterministik bahwa koordinasi terdistribusi bekerja sesuai spesifikasi",
        ], "type": "bullet"},
    ],
    mental_model="Unit and concurrency tests strictly validate LIFO rollback, idempotency, and semantic locks.",
    source_path="labs/29-saga-pattern/tests/saga_test.go:13",
    slide_num="15 / 24"
)

# Slide 16: When Compensation Fails (Operational Realities)
add_slide(
    categories=["KEGAGALAN OPERASIONAL", "RECOVERY"],
    title="Ketika Kompensasi Gagal: Realitas Lapangan",
    subtitle="Penanganan status COMPENSATE_FAILED dan strategi pemulihan sistem",
    content_items=[
        {"title": "KONDISI COMPENSATE_FAILED", "items": [
            "Bagaimana jika refund gateway menghasilkan timeout atau status 500 error?",
            "Saga tidak dapat melakukan rollback otomatis tak terbatas",
            "Orchestrator mencatat status langkah menjadi COMPENSATE_FAILED dan agregasi error",
        ], "type": "bullet"},
        {"title": "STRATEGI MITIGASI PRODUKSI", "items": [
            "Dead-Letter Queue (DLQ): Menyimpan pesan kompensasi yang gagal untuk proses retry terpisah",
            "Automated Alerting: Memberi sinyal instan kepada tim engineer on-call",
            "Idempotent Retry Worker: Mencoba kembali kompensasi dengan pola exponential backoff",
        ], "type": "bullet"},
        {"title": "MANUAL RECONCILIATION", "items": [
            "Dashboard operasional untuk meninjau status saga yang mengalami kebuntuan",
            "Customer service atau finance dapat melakukan penyesuaian manual (out-of-band)",
            "Audit trail log langkah menjadi sumber kebenaran tunggal investigasi",
        ], "type": "bullet"},
    ],
    mental_model="When business compensation fails, isolate to DLQs and rely on idempotent retry workers or manual fixes.",
    source_path="labs/29-saga-pattern/internal/saga/orchestrator.go:100",
    slide_num="16 / 24"
)

# Slide 17: ACID vs 2PC vs Saga Comparison
add_slide(
    categories=["TABEL EVALUASI", "PERBANDINGAN"],
    title="Evaluasi Arsitektur: ACID vs 2PC vs Saga",
    subtitle="Perbandingan mendalam dimensi konsistensi, ketersediaan, dan skalabilitas",
    content_items=[
        {"title": "TABEL PERBANDINGAN TIGA PENDEKATAN", "items": [
            "Dimensi          | ACID Monolitik | Two-Phase Commit (2PC) | Saga Pattern",
            "-----------------|----------------|------------------------|---------------------",
            "Konsistensi      | Strict Immed.  | Strict Immed.          | Eventual Consistency",
            "Isolasi (I)      | Didukung DBMS  | Didukung DBMS          | Perlu Semantic Lock",
            "Ketersediaan     | Sangat Tinggi  | Rendah (Blocking)      | Sangat Tinggi",
            "Skalabilitas     | Terbatas 1 DB  | Buruk lintas jaringan  | Linear horizontal",
            "Kompleksitas     | Sangat Rendah  | Menengah               | Tinggi (Kompensasi)",
        ], "type": "table"},
        {"title": "TRADE-OFF KONSISTENSI", "items": [
            "ACID menjamin tidak ada dirty read tanpa coding tambahan",
            "Saga menukar isolasi instan demi mempertahankan ketersediaan sistem 24/7",
            "Cocok untuk sistem high-throughput di mana lock contention 2PC tidak dapat diterima",
        ], "type": "bullet"},
        {"title": "KESIMPULAN PEMILIHAN", "items": [
            "Tetap gunakan ACID lokal di dalam boundary tiap microservice",
            "Gunakan Saga secara eksklusif untuk koordinasi lintas service terpisah",
        ], "type": "bullet"},
    ],
    mental_model="Choose ACID for local boundaries; use Sagas to preserve availability across distributed systems.",
    source_path="labs/29-saga-pattern/content/02-master-draft.md:5",
    slide_num="17 / 24"
)

# Slide 18: Industry Standards & Production Tooling
add_slide(
    categories=["STANDAR INDUSTRI", "TOOLING"],
    title="Ekosistem & Framework Saga di Industri",
    subtitle="Solusi open-source dan cloud-native untuk orchestrator produksi",
    content_items=[
        {"title": "TEMPORAL.IO & CADENCE", "items": [
            "Stateful workflow engine terkemuka untuk orchestration saga skala enterprise",
            "Menjamin eksekusi fungsi Go/Java/TypeScript tahan banting terhadap crash server",
            "Secara otomatis mencatat event history dan menangani retry serta timeout kompensasi",
        ], "type": "bullet"},
        {"title": "AWS STEP FUNCTIONS", "items": [
            "Serverless workflow orchestrator berbasis declarative JSON/YAML state machine",
            "Terintegrasi langsung dengan AWS Lambda, SQS, SNS, dan DynamoDB",
            "Mendukung visualisasi alur eksekusi grafis dan penanganan Catch/Retry terstruktur",
        ], "type": "bullet"},
        {"title": "EVENTUM & MASS TRANSIT", "items": [
            "Eventuate Tram / MassTransit untuk implementasi choreography berbasis transactional messaging",
            "Sering dipadukan dengan Outbox Pattern (Lab 21) untuk menjamin atomilitas event publishing",
            "Solusi tepat saat membangun arsitektur event-driven terdistribusi murni",
        ], "type": "bullet"},
    ],
    mental_model="Production sagas use workflow engines like Temporal or Step Functions to survive infrastructure crashes.",
    source_path="labs/29-saga-pattern/content/05-key-takeaways.md:12",
    slide_num="18 / 24"
)

# Slide 19: Outbox Pattern Integration
add_slide(
    categories=["POLA GABUNGAN", "ARSITEKTUR"],
    title="Sinergi: Saga Pattern + Transactional Outbox",
    subtitle="Menjamin penerbitan event tidak pernah hilang saat database lokal di-commit",
    content_items=[
        {"title": "DUAL-WRITE PROBLEM DI SAGA", "items": [
            "Sebuah service harus melakukan dua hal: update database lokal dan publish event ke broker",
            "Jika database commit berhasil tapi broker pesan mati, saga akan terhenti permanen",
            "Jika broker publish berhasil tapi database commit gagal, data menjadi tidak konsisten",
        ], "type": "bullet"},
        {"title": "SOLUSI TRANSACTIONAL OUTBOX", "items": [
            "Simpan domain event ke dalam tabel outbox di database yang sama dalam satu transaksi ACID",
            "Proses worker independen membaca tabel outbox dan mempublikasikannya ke message broker",
            "Menghilangkan risiko dual-write inconsistency pada arsitektur saga choreography",
        ], "type": "bullet"},
        {"title": "KORELASI LAB", "items": [
            "Lab 21 (Transactional Outbox) memberikan fondasi pengiriman pesan andal",
            "Lab 29 (Saga Pattern) mengorkestrasi alur transaksi bisnis tingkat tinggi",
            "Kombinasi keduanya merupakan gold standard arsitektur microservices modern",
        ], "type": "bullet"},
    ],
    mental_model="Pair Sagas with Transactional Outbox to eliminate the dual-write problem across service boundaries.",
    source_path="labs/29-saga-pattern/content/02-master-draft.md:37",
    slide_num="19 / 24"
)

# Slide 20: Real-World Case Study: E-Commerce Checkout
add_slide(
    categories=["STUDI KASUS", "PRODUKSI"],
    title="Studi Kasus Nyata: E-Commerce Checkout",
    subtitle="Menyusuri siklus hidup 4 tahap pesanan dari keranjang hingga pengiriman",
    content_items=[
        {"title": "TAHAP 1: ORDER CREATION (COMPENSABLE)", "items": [
            "Client menekan tombol Bayar Sekarang di aplikasi",
            "OrderService mencatat pesanan PENDING dan memasang semantic lock pada order ID",
            "Kompensasi jika gagal di langkah lanjutan: CancelOrder",
        ], "type": "bullet"},
        {"title": "TAHAP 2 & 3: PAYMENT & INVENTORY", "items": [
            "PaymentService memotong saldo kartu dengan idempotency key (Kompensasi: RefundPayment)",
            "InventoryService mengalokasikan stok gudang fisik (Kompensasi: ReleaseInventory)",
            "Langkah inventaris bertindak sebagai Pivot Transaction",
        ], "type": "bullet"},
        {"title": "TAHAP 4: APPROVAL (RETRYABLE)", "items": [
            "Setelah stok aman, order status diubah menjadi APPROVED dan semantic lock dilepas",
            "Email invoice dikirimkan ke pelanggan via retryable worker",
            "Pesanan kini siap diproses oleh divisi logistik pengiriman",
        ], "type": "bullet"},
    ],
    mental_model="E-commerce checkouts balance fast client confirmation with background eventual consistency.",
    source_path="labs/29-saga-pattern/cmd/demo/main.go:21",
    slide_num="20 / 24"
)

# Slide 21: Production Readiness Checklist
add_slide(
    categories=["PANDUAN OPERASIONAL", "CHECKLIST"],
    title="Checklist Kesiapan Produksi Saga",
    subtitle="9 poin verifikasi arsitektur sebelum merilis ke lingkungan live",
    content_items=[
        {"title": "CHECKLIST PENERAPAN SAGA", "items": [
            "[ ] Setiap compensable step memiliki fungsi kompensasi yang teruji",
            "[ ] Seluruh endpoint kompensasi dan retry bersifat 100% idempotent",
            "[ ] Semantic lock atau versioning diterapkan untuk mencegah dirty reads",
            "[ ] Dead-Letter Queue (DLQ) disiapkan untuk kompensasi yang gagal",
            "[ ] Log audit status langkah disimpan di storage persisten",
            "[ ] Context timeout dan cancellation dikelola tanpa membatalkan rollback",
            "[ ] Monitoring dan alerting aktif untuk status COMPENSATE_FAILED",
            "[ ] Prosedur rekonsiliasi manual didokumentasikan untuk tim finance",
            "[ ] Concurrency race detector lulus tanpa anomali data",
        ], "type": "bullet"},
        {"title": "ATURAN UTAMA", "items": [
            "Jangan pernah mengasumsikan kompensasi akan selalu berhasil tanpa cela",
            "Rancang sistem dengan asumsi jaringan akan gagal pada momen terburuk",
            "Saga yang baik selalu memiliki jalur pemulihan darurat",
        ], "type": "bullet"},
        {"title": "PERAN OBSERVABILITY", "items": [
            "Distributed tracing (OpenTelemetry) sangat krusial untuk melacak Saga ID lintas service",
        ], "type": "bullet"},
    ],
    mental_model="A complete production checklist covers idempotency, semantic locks, DLQs, and observability.",
    source_path="labs/29-saga-pattern/content/02-master-draft.md:143",
    slide_num="21 / 24"
)

# Slide 22: Scope & Boundaries of Lab 29
add_slide(
    categories=["BATAS LAB", "TRANSPARANSI"],
    title="Batasan & Lingkup Implementasi Lab 29",
    subtitle="Membedakan simplifikasi edukasi dengan kebutuhan enterprise sesungguhnya",
    content_items=[
        {"title": "IN-MEMORY STATE & MOCK STORAGE", "items": [
            "Lab ini menggunakan in-memory map dan sync.Mutex bukan SQL/NoSQL database nyata",
            "State machine orchestrator tidak disimpan di persistent storage disk",
            "Jika proses Go di-kill di tengah eksekusi, in-flight saga state akan hilang",
        ], "type": "bullet"},
        {"title": "SIMPLIFIKASI JARINGAN", "items": [
            "Pemanggilan antar service dilakukan secara in-process function call, bukan gRPC/HTTP",
            "EventBus menggunakan Go channel/slice in-memory, bukan Apache Kafka atau RabbitMQ",
            "Tidak menyimulasikan network partition nyata, tetapi memverifikasi logika koordinasi",
        ], "type": "bullet"},
        {"title": "NILAI INTI YANG DIBUKTIKAN", "items": [
            "Algoritma kompensasi LIFO terbukti secara matematis dan logis",
            "Mekanisme idempotency key dan semantic lock diverifikasi 100% fungsional",
            "Fondasi konseptual yang kokoh sebelum mengadopsi engine seperti Temporal",
        ], "type": "bullet"},
    ],
    mental_model="This lab proves algorithmic correctness in-memory; production requires persistent state engines.",
    source_path="labs/29-saga-pattern/README.md:1",
    slide_num="22 / 24"
)

# Slide 23: Core Engineering Principles
add_slide(
    categories=["PRINSIP INTI", "RINGKASAN"],
    title="Tiga Prinsip Inti Transaksi Terdistribusi",
    subtitle="Pedoman mendasar bagi software engineer dalam mendesain sistem microservices",
    content_items=[
        {"title": "PRINSIP 1: COMMIT LOKAL, KOMPENSASI GLOBAL", "items": [
            "Jangan pernah mencoba mengunci database jarak jauh secara sinkron",
            "Lakukan commit cepat di batas lokal dan andalkan kompensasi jika gagal",
            "Maksimumkan ketersediaan sistem dan ketahanan terhadap partisi jaringan",
        ], "type": "bullet"},
        {"title": "PRINSIP 2: IDEMPOTENSI ADALAH HARGA MATI", "items": [
            "Dalam sistem terdistribusi, kegagalan jaringan tidak terhindarkan",
            "Retry yang aman hanya dapat terwujud jika setiap langkah memiliki idempotency key",
            "Desain fungsi agar dapat dipanggil berkali-kali dengan hasil yang sama persis",
        ], "type": "bullet"},
        {"title": "PRINSIP 3: LINDUNGI ISOLASI SECARA SEMANTIK", "items": [
            "Ketiadaan isolasi ACID harus diimbangi oleh semantic locking di level aplikasi",
            "Tandai status data PENDING agar modifikasi konkuren tidak merusak integritas",
            "Lindungi sistem dari anomali data dirty reads dan lost updates",
        ], "type": "bullet"},
    ],
    mental_model="Three rules of sagas: commit locally, make every step idempotent, and lock semantically.",
    source_path="labs/29-saga-pattern/content/05-key-takeaways.md:1",
    slide_num="23 / 24"
)

# Slide 24: Closing & Resources (Dark Navy)
add_slide(
    categories=["RINGKASAN"],
    title="Ringkasan Lab 29: Saga Pattern",
    subtitle=f"Commit: {COMMIT_SHA[:8]} | 9/9 Test Suites PASS | Repository: lukman-ss/software-engineering-lab",
    content_items=[
        {"title": "HASIL VERIFIKASI", "items": [
            "Happy path orkestrasi: 4 langkah tuntas dengan status APPROVED",
            "Rollback kompensasi LIFO: Berhasil refund payment dan cancel order saat inventaris gagal",
            "9 unit, concurrency, dan cancellation test suites lulus sempurna (go test -v -race)",
        ], "type": "bullet"},
        {"title": "MENTAL MODEL UTAMA", "items": [
            "A saga is a sequence of local transactions where each step commits immediately.",
            "If any step fails, compensating transactions execute in reverse order (LIFO).",
            "Eventual consistency replaces strict 2PC locking across service boundaries.",
        ], "type": "bullet"},
        {"title": "RESOURCES & AUDIT", "items": [
            f"Repository: {REPO_URL}",
            f"Commit SHA: {COMMIT_SHA}",
            "Jalankan Demo: go run ./cmd/demo (di folder lab)",
            "Jalankan Test: go test -v -race ./... (di folder lab)",
        ], "type": "bullet"},
    ],
    mental_model="Sagas trade strict ACID isolation for resilient, scalable eventual consistency across microservices.",
    source_path="labs/29-saga-pattern/README.md",
    slide_num=""
)

if __name__ == "__main__":
    print(f"Total slides defined: {len(SLIDES)}")
    for i, s in enumerate(SLIDES):
        print(f"Slide {i+1}: {s['title']} ({len(s['content_items'])} content items)")
