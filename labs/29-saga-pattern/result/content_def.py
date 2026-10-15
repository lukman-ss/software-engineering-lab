#!/usr/bin/env python3
"""
Publication Builder for Software Engineering Lab 29: Saga Pattern.
Generates 24-slide LinkedIn carousel covering Orchestration vs Choreography,
compensating transactions (LIFO rollback), idempotency keys, semantic locks,
and concurrency safety in distributed transaction management.
"""
import os

OUTPUT_DIR = "labs/29-saga-pattern/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/29-saga-pattern"
COMMIT_SHA = "ba4b6216da6a06dee61959231d124cd44085bbf5"
LAB_NO = "29"
SLUG = "saga-pattern"
TITLE_ID = "Saga Pattern"
SUBTITLE_ID = "Mengelola Transaksi Terdistribusi Tanpa Two-Phase Commit di Microservices"
TITLE_EN = "Saga Pattern"
SUBTITLE_EN = "Distributed Transaction Management Without Two-Phase Commit in Microservices"
AUTHOR = "Lukman / lukman-ss"
COVER_IMAGE_PATH = "labs/29-saga-pattern/result/image-cover/cover_art.jpg"

DARK_NAVY = "0b1b2d"
TEAL = "087f76"
CREAM = "f4f1e8"
WHITE = "ffffff"
SLATE = "52606b"
ORANGE = "d97706"
RED = "dc2626"
GREEN = "16a34a"

SLIDES = []

def add_slide(categories, title, subtitle, content_items, mental_model, source_path="", slide_num="", narrative="", narrative_en="", layout="", title_en=""):
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
        "layout": layout,
        "title_en": title_en or title
    })

# Slide 1: Cover
add_slide(
    categories=["COVER"],
    title=TITLE_ID,
    subtitle=SUBTITLE_ID,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "Retro CRT monitor menampilkan sequence diagram Orkestrasi Saga",
            "Diagram LIFO compensation rollback dengan panah merah dan hijau",
            f"Penulis: {AUTHOR}"
        ]}
    ],
    mental_model="Saga pattern replaces 2PC with independent local transactions and compensating rollbacks.",
    source_path=COVER_IMAGE_PATH,
    slide_num="01 / 24",
    layout="cover"
)

# Slide 2: Problem nyata
add_slide(
    categories=["PROBLEM", "DISTRIBUTED TRANSACTION"],
    title="Ketika Database Terpisah Tidak Bisa Berkomitmen Bersama",
    subtitle="Setiap layanan memiliki database sendiri — transaksi ACID klasik runtuh di boundary microservice.",
    content_items=[
        {"type": "cards", "title": "MASALAH DISTRIBUSI TRANSAKSI", "items": [
            {"title": "Database Per-Service", "desc": "Order, Payment, dan Inventory terpisah; tidak ada shared transaction boundary antar layanan."},
            {"title": "Two-Phase Commit Rapuh", "desc": "2PC mengunci semua database secara sinkron hingga seluruh peserta setuju — blocking dan lambat."},
            {"title": "Partisi Jaringan Fatal", "desc": "Satu node down mengunci seluruh transaksi global; availability menurun drastis di bawah network partition."}
        ]}
    ],
    narrative="Di arsitektur microservices modern, setiap service mengelola database sendiri. Saat flow bisnis melibatkan order-creation, payment-processing, dan inventory-reservation secara berurutan, transaksi ACID klasik yang mengandalkan single database transaction langsung runtuh. Dua-phase commit (2PC) muncul sebagai solusi teoretis, namun protokol ini mengunci semua resource secara sinkron hingga seluruh peserta memberikan suara commit — menyebabkan latency tinggi dan vulnerabilitas terhadap network partition.",
    narrative_en="In modern microservice architectures, each service manages its own database. When a business flow involves order creation, payment processing, and inventory reservation sequentially, classical ACID transactions relying on a single database boundary break down. Two-phase commit emerges as a theoretical solution, but it locks all resources synchronously until every participant votes commit — causing high latency and vulnerability to network partitions.",
    mental_model="When services own their data, distributed transactions need local independence plus compensation.",
    source_path="content/02-master-draft.md:3-15",
    slide_num="02 / 24",
    layout="cards"
)

# Slide 3: Mengapa review manual / pendekatan naif gagal
add_slide(
    categories=["WHY FAILS", "NAIVE APPROACH"],
    title="Mengapa Pendekatan Naif Gagal Menangani Transaksi Multi-Layanan",
    subtitle="Orkestrasi manual tanpa kompensasi menciptakan inkonsistensi yang tersembunyi dan sulit dideteksi.",
    content_items=[
        {"type": "code", "title": "internal/services/services.go — Manual Cross-Service Call", "items": [
            "func (o *OrderService) CreateFullOrder(orderID string) error {",
            "    o.CreateOrder(orderID)     // DB-commit #1",
            "    p.ProcessPayment(orderID)  // DB-commit #2 — jika gagal?",
            "    i.ReserveInventory(...)    // DB-commit #3 — jika gagal?",
            "    o.ApproveOrder(orderID)    // DB-commit #4",
            "}"
        ]},
        {"type": "bullet", "title": "KEGAGALAN TANPA SAGA", "items": [
            "Payment gagal setelah order dibuat → order tercatat tanpa pembayaran (status PENDING selamanya).",
            "Inventory reserve gagal setelah payment sukses → uang terpotong tapi stok tidak dipesan.",
            "Tidak ada jejak kompensasi; rekonsiliasi manual diperlukan untuk memperbaiki state."
        ]}
    ],
    narrative="Implementasi naif mencoba memanggil setiap service berurutan tanpa mekanisme rollback lintas domain. Jika payment service gagal setelah order berhasil dibuat, database order mencatat status PENDING namun payment mencatat gagal. State sistem menjadi inkonsisten dan tidak ada jalan mundur otomatis. Pendekatan manual ini juga menyulitkan observabilitas karena tidak ada log eksekusi terpusat yang merekam status setiap langkah.",
    narrative_en="The naive implementation attempts sequential service calls without cross-domain rollback mechanisms. If payment fails after the order is created, the order database records a PENDING status while payment records failure. The system enters an inconsistent state with no automatic recovery path. This manual approach also hinders observability since there is no centralized execution log recording each step's status.",
    mental_model="Manual sequential calls without compensation create silent distributed inconsistency.",
    source_path="internal/services/services.go • cmd/demo/main.go:11-113",
    slide_num="03 / 24",
    layout="code-bullet"
)

# Slide 4: Mental model inti / schema card
add_slide(
    categories=["MENTAL MODEL", "SAGA CORE"],
    title="Mental Model: Estafet Transaksi dengan Kompensasi Mundur",
    subtitle="Setiap langkah adalah transaksi lokal independen; jika satu gagal, langkah sebelumnya dikompensasi.",
    content_items=[
        {"type": "schema", "title": "SAGA COMPOSITION MODEL", "items": [
            {"phase": "1. Step Execute", "role": "Local Transaction", "desc": "Setiap langkah mengeksekusi transaksi lokal ke database-nya sendiri dan berhasil segera."},
            {"phase": "2. Step Fail", "role": "Trigger Compensation", "desc": "Jika langkah ke-n gagal, orchestrator memicu compensating transactions pada langkah 1 sampai n-1."},
            {"phase": "3. LIFO Rollback", "role": "Last-In First-Out", "desc": "Kompensasi berjalan dari belakang ke depan: refund payment sebelum cancel order."},
            {"phase": "4. Eventual Consistency", "role": "All or Nothing", "desc": "Setelah kompensasi selesai, seluruh sistem kembali ke state sebelum saga dimulai."}
        ]}
    ],
    narrative="Bayangkan proses estafet di mana setiap pelari menyelesaikan bagiannya dan langsung mencatat hasil ke database lokal. Jika pelari ketiga terjatuh, pelari satu dan dua berlari mundur untuk mengembalikan posisi awal mereka. Setiap langkah berdiri sendiri sebagai transaksi lokal yang commit-segera; tidak ada locking lintas layanan. Ketika kegagalan terjadi, orchestrator memicu chain of compensations secara Last-In-First-Out untuk memulihkan konsistensi.",
    narrative_en="Imagine a relay race where each runner completes their leg and immediately records results to a local database. If the third runner stumbles, the first and second runners run backward to restore their original positions. Each step stands alone as a local transaction with immediate commit — no cross-service locking exists. When failure occurs, the orchestrator triggers a Last-In-First-Out chain of compensations to restore consistency.",
    mental_model="Local transactions commit eagerly; compensations unwind them in reverse order for eventual consistency.",
    source_path="content/02-master-draft.md:13-21",
    slide_num="04 / 24",
    layout="schema"
)

# Slide 5: Diagram alur/pipeline
add_slide(
    categories=["ARCHITECTURE", "ORCHESTRATION PIPELINE"],
    title="Pipeline Orkestrasi: Dari Langkah Eksekusi Hingga Kompensasi",
    subtitle="Flow terpusat di Orchestrator: AddStep → Execute → LIFO Compensate pada kegagalan.",
    content_items=[
        {"type": "diagram", "title": "ORCHESTRATOR EXECUTION PIPELINE", "items": [
            "Client Request",
            "↓",
            "Orchestrator.Execute(ctx)",
            "↓",
            "for each Step in steps:",
            "  step.Execute(ctx) → Local DB Write → Commit",
            "  on success: append to executed list",
            "  on failure: break loop → trigger compensate()",
            "↓",
            "compensate(executed[]):",
            "  for i = len-1 down to 0:",
            "    executed[i].Compensate(ctx) → Undo DB Write",
            "  on comp fail: aggregate errors, log COMPENSATE_FAILED",
            "↓",
            "Return: nil (success) | error (step failed + compensation)"
        ]}
    ],
    narrative="Pipeline orkestrasi dimulai ketika client mengirimkan request ke orchestrator. Orchestrator mengiterasi daftar langkah secara sekuensial, mengeksekusi setiap langkah ke database lokal masing-masing service. Pada setiap keberhasilan, langkah dicatat ke dalam stack executed. Jika langkah manapun gagal, loop berhenti dan orchestrator memanggil fungsi compensate() yang berjalan mundur dari stack — memastikan refund terjadi sebelum cancel order. Log Status setiap langkah direkam secara atomik di mutex untuk konsistensi konkuren.",
    narrative_en="The orchestration pipeline begins when a client submits a request to the orchestrator. The orchestrator iterates through steps sequentially, executing each into the respective local service database. On each success, the step is recorded to the executed stack. If any step fails, the loop breaks and the orchestrator calls compensate(), which unwinds backward from the stack — ensuring refund happens before cancel order. Each step's status log is recorded atomically under mutex for concurrency consistency.",
    mental_model="Forward steps consume; backward compensations undo — orchestrated by a single stateful coordinator.",
    source_path="internal/saga/orchestrator.go:49-120",
    slide_num="05 / 24",
    layout="pipeline-diagram"
)

# Slide 6: Cuplikan kode dominan
add_slide(
    categories=["IMPLEMENTATION", "CORE CODE"],
    title="Implementasi Inti: Orchestrator Step & LIFO Compensate",
    subtitle="Struktur Step dengan Execute/Compensate pair dan fungsi compensate() berurutan mundur.",
    content_items=[
        {"type": "code", "title": "internal/saga/orchestrator.go", "items": [
            "type Step struct {",
            "    Name       string",
            "    Execute    func(ctx context.Context) error",
            "    Compensate func(ctx context.Context) error",
            "}",
            "",
            "func (o *Orchestrator) compensate(ctx context.Context, executed []Step) error {",
            "    for i := len(executed) - 1; i >= 0; i-- {",
            "        step := executed[i]",
            "        if step.Compensate != nil {",
            "            err := step.Compensate(ctx)  // LIFO unwind",
            "            ...",
            "        }",
            "    }",
            "}"
        ]}
    ],
    narrative="Struktur data Step mendefinisikan sepasang fungsi Execute dan Compensate untuk setiap transaksi lokal. Orchestrator menyimpan slice Step ini dan mengeksekusinya secara berurutan. Fungsi compensate() melakukan iterasi mundur dari index terakhir ke nol — inilah inti LIFO rollback. Jika compensate() sendiri gagal pada suatu langkah, error di agregasi ke dalam slice compErrors dan dilaporkan bersama error langkah utama.",
    narrative_en="The Step data structure defines an Execute-Compensate function pair for each local transaction. The orchestrator stores this Step slice and executes them sequentially. The compensate() function iterates backward from the last index to zero — this is the core of LIFO rollback. If compensate() itself fails on a step, the error is aggregated into a compErrors slice and reported alongside the primary step error.",
    mental_model="Each step owns both its forward action and its undo action — encapsulated as a single Step struct.",
    source_path="internal/saga/orchestrator.go:19-23,92-112",
    slide_num="06 / 24",
    layout="code"
)

# Slide 7: Failure path / dangling reference konkret
add_slide(
    categories=["FAILURE PATH", "ROLLBACK ANATOMY"],
    title="Failure Path: ReserveInventory Gagal — Urutan Kompensasi LIFO",
    subtitle="Kegagalan stok kosong memicu refund payment sebelum cancel order, sesuai urutan balik stack.",
    content_items=[
        {"type": "cards", "title": "ALUR KEGAGALAN DARI LATAR DEPAN HINGGA BELAKANG", "items": [
            {"title": "Langkah 1: CreateOrder", "desc": "Order dibuat dengan status PENDING. Semantic lock diterapkan pada orderID. Berhasil ✓"},
            {"title": "Langkah 2: ProcessPayment", "desc": "Payment diproses dengan idempotency key. Dana terpotong. Berhasil ✓"},
            {"title": "Langkah 3: ReserveInventory GAGAL", "desc": "Stok item-1 hanya 1 unit, diminta 5. Error: out of stock. Loop berhenti."},
            {"title": "Kompensasi mundur (LIFO)", "desc": "Compensate Payment (Refund) → Compensate Order (Cancel). State pulih ke awal. ✓"}
        ]}
    ],
    narrative="Skenario ini menunjukkan failure path realistis di mana reservasi inventaris gagal karena stok tidak mencukupi. Order sudah dibuat dan payment sudah diproses sebelum kegagalan terdeteksi. Orchestrator kemudian memicu kompensasi LIFO: pertama-tama memproses refund payment, kemudian membatalkan order. Dengan urutan ini, sistem kembali ke keadaan sebelum saga dimulai — payment refund terjadi sebelum order cancel, menjaga konsistensi urutan finansial.",
    narrative_en="This scenario demonstrates a realistic failure path where inventory reservation fails due to insufficient stock. The order has been created and payment processed before the failure is detected. The orchestrator then triggers LIFO compensation: first processing payment refund, then cancelling the order. Through this order, the system returns to its pre-saga state — payment refund occurs before order cancellation, maintaining financial sequence consistency.",
    mental_model="Compensations must unwind in reverse execution order to prevent orphaned side effects.",
    source_path="tests/saga_test.go:74-147",
    slide_num="07 / 24",
    layout="cards"
)

# Slide 8: Tabel perbandingan mekanisme vs gejala
add_slide(
    categories=["COMPARISON", "MECHANISM VS MANIFESTATION"],
    title="Tabel Perbandingan: Mekanisme Kompensasi vs Gejala Sistem",
    subtitle="Memetakan setiap tipe langkah saga ke manifestasi perilaku saat sukses dan gagal.",
    content_items=[
        {"title": "MATRIKS PERILAKU SAGA STEP", "items": [
            "Tipe Langkah | Eksekusi Sukses | Eksekusi Gagal | Kompensasi",
            "Compensable Transaction | Local DB commit persisten | Rollback via Compensate() | Refund, Cancel, Release",
            "Pivot Transaction | Titik balik tidak bisa dibatalkan | Saga dijamin mencapai final state | Tidak ada kompensasi",
            "Retryable Transaction | Post-pivot, dijamin sukses | Idempotent retry | Tidak diperlukan",
            "CreateOrder Step | Status PENDING + semantic lock | Rollback: delete lock + CANCELLED | ✓",
            "ProcessPayment Step | Payment recorded + idempotency key | Rollback: delete payment entry | ✓ Refund",
            "ReserveInventory Step | Stock decremented | Rollback: stock restored via Release() | ✓",
            "ApproveOrder Step | Status APPROVED + lock released | Tanpa kompensasi (final step) | —"
        ], "type": "table"}
    ],
    narrative="Tabel ini memetakan taksonomi tiga tipe transaksi saga menurut Azure Architecture Center ke implementasi spesifik lab. Compensation step adalah yang paling umum — setiap langkah kecuali ApproveOrder memiliki fungsi Compensate. Pivot transaction merepresentasikan titik balik yang tidak dapat digagalkan, sedangkan retryable transaction menjamin kesuksesan setelah pivot. Dalam lab ini, ApproveOrder berperan sebagai final step tanpa kompensasi karena ia sudah merupakan state terminal.",
    narrative_en="This table maps the three saga transaction types from Azure Architecture Center to specific lab implementations. The compensable transaction is most common — every step except ApproveOrder has a Compensate function. Pivot transactions represent irreversible turning points, while retryable transactions guarantee success post-pivot. In this lab, ApproveOrder serves as the final terminal state step without compensation.",
    mental_model="Not all saga steps need compensation — final approval marks the pivot of no return.",
    source_path="content/02-master-draft.md:23-27 • internal/saga/orchestrator.go:9-17",
    slide_num="08 / 24",
    layout="table"
)

# Slide 9: Decision sequence — 3 focused decision cards
add_slide(
    categories=["DECISION SEQUENCE", "ARCHITECTURE GUIDE"],
    title="Keputusan Arsitektur: Orkestrasi vs Koreografi",
    subtitle="Tiga kriteria teknis menentukan pilihan pola koordinasi saga yang tepat.",
    content_items=[
        {"type": "cards", "title": "PILIHAN POLA KOORDINASI SAGA", "items": [
            {"title": "1. Orkestrasi (Orchestrator)", "desc": "Koordinator terpusat mengatur urutan, menyimpan log, dan memicu kompensasi. Cocok untuk flow linear dengan dependencies jelas."},
            {"title": "2. Koreografi (Choreography)", "desc": "Layanan saling mendengarkan event tanpa koordinator pusat. Cocok untuk flow event-driven dengan loose coupling antar service."},
            {"title": "3. Hybrid Pattern", "desc": "Orkestrasi untuk orchestration-critical path; koreografi untuk notification-only paths. Menggabungkan kekuatan kedua pendekatan."}
        ]}
    ],
    narrative="Pemilihan antara orkestrasi dan koreografi bukanlah keputusan hitam-putih. Orkestrasi unggul ketika urutan langkah memiliki dependensi ketat dan perlu logging terpusat, seperti pada flow checkout e-commerce. Koreografi lebih cocok untuk sistem event-driven di mana layanan perlu bereaksi mandiri terhadap event tanpa mengetahui orchestrator. Lab ini mengimplementasikan kedua pola untuk menunjukkan trade-off masing-masing.",
    narrative_en="The choice between orchestration and choreography is not binary. Orchestration excels when step ordering has tight dependencies and centralized logging is needed, such as in e-commerce checkout flows. Choreography suits event-driven systems where services react independently to events without knowing about a central orchestrator. This lab implements both patterns to demonstrate each trade-off.",
    mental_model="Orchestration gives control; choreography gives autonomy. Choose by dependency density.",
    source_path="internal/saga/orchestrator.go • internal/saga/choreography.go",
    slide_num="09 / 24",
    layout="cards"
)

# Slide 10: Bukti eksekusi / terminal output
add_slide(
    categories=["EVIDENCE", "EXECUTION OUTPUT"],
    title="Bukti Eksekusi: Semua Test Suite Lolos Bersih",
    subtitle="Output Go test membuktikan happy path, LIFO rollback, idempotency, semantic lock, dan concurrency safety.",
    content_items=[
        {"type": "code", "title": "TERMINAL OUTPUT: go test -v ./...", "items": [
            "=== RUN   TestOrchestrator_HappyPath",
            "    --- PASS: TestOrchestrator_HappyPath (0.00s)",
            "=== RUN   TestOrchestrator_FailureCompensatesLIFO",
            "    --- PASS: TestOrchestrator_FailureCompensatesLIFO (0.00s)",
            "=== RUN   TestPayment_Idempotency",
            "    --- PASS: TestPayment_Idempotency (0.00s)",
            "=== RUN   TestSemanticLock",
            "    --- PASS: TestSemanticLock (0.00s)",
            "=== RUN   TestOrchestrator_Concurrency",
            "    --- PASS: TestOrchestrator_Concurrency (0.00s)",
            "=== RUN   TestChoreography_Flow",
            "    --- PASS: TestChoreography_Flow (0.00s)",
            "=== RUN   TestChoreography_FailureCompensates",
            "    --- PASS: TestChoreography_FailureCompensates (0.00s)",
            "=== RUN   TestOrchestrator_CompensationErrorPropagated",
            "    --- PASS: TestOrchestrator_CompensationErrorPropagated (0.00s)",
            "=== RUN   TestOrchestrator_ContextCancellation",
            "    --- PASS: TestOrchestrator_ContextCancellation (0.00s)",
            "PASS: ok  labs/29-saga-pattern/tests  0.003s"
        ]}
    ],
    narrative="Kesembilan test case berjalan mulus tanpa kegagalan satupun. Happy path membuktikan empat langkah sukses mengubah order status menjadi APPROVED dengan stock berkurang 2 unit. Failure test membuktikan LIFO rollback: order kembali CANCELLED, payment di-refund, stock kembali 1. Concurrency test dengan 10 goroutine simultan memastikan atomicity mutex tetap terjaga.",
    narrative_en="All nine test cases run cleanly without a single failure. The happy path proves four steps succeed, changing order status to APPROVED with stock reduced by 2 units. The failure test proves LIFO rollback: order returns to CANCELLED, payment refunded, stock restored to 1. The concurrency test with 10 simultaneous goroutines ensures mutex atomicity remains intact.",
    mental_model="Clean test suite with zero failures validates both forward execution and backward compensation paths.",
    source_path="tests/saga_test.go",
    slide_num="10 / 24",
    layout="code"
)

# Slide 11: Analisis batasan bukti
add_slide(
    categories=["AUDIT", "EVIDENCE BOUNDARY"],
    title="Batasan Bukti Lab: Yang Dibuktikan vs Belum",
    subtitle="Simulasi in-memory tidak membuktikan behavior di bawah network partition atau database crash riil.",
    content_items=[
        {"type": "cards", "title": "YANG TIDAK DIBUKTIKAN", "items": [
            {"title": "Network Partition", "desc": "Tidak ada simulasi network partition antara orchestrator dan service; semua service berjalan dalam proses Go yang sama."},
            {"title": "Database Crash", "desc": "Tidak diuji apa yang terjadi jika database crashed di tengah eksekusi step — tidak ada persistence layer riil."},
            {"title": "Persistensi Log", "desc": "Step log hanya disimpan di memori; tidak ada audit trail yang persisten ke disk atau database."},
            {"title": "Distributed EventBus", "desc": "EventBus di choreography berjalan in-memory; tidak diuji dengan message broker riil seperti Kafka atau RabbitMQ."}
        ]}
    ],
    narrative="Lab ini membuktikan kebenaran logika komputasional dan concurrency safety dalam lingkungan in-memory. Namun, batasan tegas harus diakui: tidak ada simulasi network partition, database crash, atau persistence layer riil. Dalam produksi, kegagalan kompensasi bisa terjadi karena service unavailable — scenario yang tidak tercakup oleh test suite saat ini. Dead-letter queue dan manual reconciliation adalah countermeasures yang diperlukan di luar cakupan lab.",
    narrative_en="This lab proves computational correctness and concurrency safety in an in-memory environment. However, strict limitations must be acknowledged: there is no network partition simulation, no real database crash, and no persistence layer. In production, compensation failure can occur due to service unavailability — scenarios not covered by the current test suite. Dead-letter queues and manual reconciliation are countermeasures needed beyond this lab's scope.",
    mental_model="In-memory proofs validate logic; production requires persistence, observability, and DLQ handling.",
    source_path="engineering/03-execution-result.md • engineering-audit/06-verdict.md",
    slide_num="11 / 24",
    layout="cards"
)

# Slide 12: Trade-off arsitektur
add_slide(
    categories=["TRADE-OFF", "ARCHITECTURAL COMPROMISE"],
    title="Trade-off Arsitektur: Availability vs Isolation",
    subtitle="Saga mengorbankan isolasi ACID demi ketersediaan tinggi — eventual consistency adalah harga yang harus dibayar.",
    content_items=[
        {"title": "MATRIKS TRADE-OFF SAGA", "items": [
            "Dimensi | ACID (Monolitik) | Saga Pattern | Dampak Produksi",
            "Isolasi (I dalam ACID) | serializable penuh | tidak ada isolasi lintas service | dirty reads, lost updates mungkin terjadi",
            "Ketersediaan | blocking 2PC | non-blocking, each step commit segera | saga jauh lebih available di bawah load tinggi",
            "Latensi | satu round-trip ke semua peserta | N round-trips sekuensial + kompensasi | total latency lebih tinggi per transaksi",
            "Konsistensi | strong consistency | eventual consistency | window inconsistency antar step execute sampai kompensasi",
            "Complexity | database engine handles it | application-level orchestration | developer wrote compensation logic manually",
            "Observability | transaction logs | step-level logs + compensation trace | lebih transparan, lebih mudah debug"
        ], "type": "table"}
    ],
    narrative="Saga menukar isolation level tinggi menjadi availability tinggi. Setiap langkah committed segera tanpa lock lintas service, sehingga tidak ada blocking antar layanan. Namun, window inconsistency muncul antara step execute dan kompensasi — order bisa tercatat APPROVED sementara payment belum diproses. Developer harus menulis compensation logic secara eksplisit, menambah complexity aplikasi dibanding solusi monolitik.",
    narrative_en="Saga trades high isolation for high availability. Each step commits immediately without cross-service locking, so there is no blocking between services. However, a consistency window exists between step execution and compensation — an order may appear APPROVED while payment remains unprocessed. Developers must write compensation logic explicitly, adding application complexity compared to monolithic solutions.",
    mental_model="Availability over isolation: commit early, compensate later, accept eventual consistency.",
    source_path="content/02-master-draft.md:120-135",
    slide_num="12 / 24",
    layout="table"
)

# Slide 13: Concurrency / validation pipeline
add_slide(
    categories=["CONCURRENCY", "THREAD SAFETY"],
    title="Pipeline Konkurensi: Goroutine → Validate → Mutex → Append Error",
    subtitle="Orchestrator menggunakan sync.Mutex untuk melindungi step list dan log agar aman dari concurrent Execute.",
    content_items=[
        {"type": "code", "title": "internal/saga/orchestrator.go — Concurrency Guards", "items": [
            "func (o *Orchestrator) AddStep(step Step) {",
            "    o.mu.Lock()",
            "    defer o.mu.Unlock()",
            "    o.steps = append(o.steps, step)",
            "}",
            "",
            "func (o *Orchestrator) Execute(ctx context.Context) error {",
            "    o.mu.Lock()",
            "    steps := make([]Step, len(o.steps))",
            "    copy(steps, o.steps)  // snapshot tanpa lock holder",
            "    o.mu.Unlock()",
            "    // ...execute steps without holding mu during step calls",
            "    o.mu.Lock()",
            "    o.logs = append(o.logs, StepLog{Name: step.Name, Status: ...})",
            "    o.mu.Unlock()",
            "}"
        ]},
        {"type": "bullet", "title": "KONTRAK KEAMANAN", "items": [
            "Mu memegang slice steps saat AddStep danExecute — mencegah race pada read.",
            "Snapshot steps saat Execute: salinan sementara dilepaskan sebelum step call.",
            "Log append dilindungi mutex — konsistensi log terjamin di bawah 10 concurrent goroutine.",
            "Hasil: go test -race laporan clean pada TestOrchestrator_Concurrency (10 workers)."
        ]}
    ],
    narrative="Keamanan konkuren dicapai melalui strategi snapshot-and-release. Saat Execute dipanggil, orchestrator menyalin slice steps ke variabel lokal sementara mutex dilepaskan, menghindari blocking selama step execution berlangsung. Log append dan read dilindungi mutex individual untuk memastikan konsistensi pembacaan. Test concurrency menjalankan 10 goroutine simultan dengan mutex dan race detector Go melaporkan status bersih tanpa data race.",
    narrative_en="Concurrency safety is achieved through a snapshot-and-release strategy. When Execute is called, the orchestrator copies the steps slice to a local variable while releasing the mutex, avoiding blocking during step execution. Log append and read operations are protected by individual mutexes to ensure read consistency. The concurrency test runs 10 simultaneous goroutines with the Go race detector reporting a clean result with no data races.",
    mental_model="Snapshot-then-release mutex pattern prevents deadlock while ensuring log consistency under concurrency.",
    source_path="internal/saga/orchestrator.go:43-89",
    slide_num="13 / 24",
    layout="code-bullet"
)

# Slide 14: Before vs After comparison cards
add_slide(
    categories=["COMPARISON", "BEFORE-AFTER"],
    title="Before vs After: Dari Manual Calls ke Saga Orkestrasi",
    subtitle="Perbandingan visual antara implementasi naif tanpa kompensasi dan implementasi saga dengan LIFO rollback.",
    content_items=[
        {"type": "cards", "title": "PERBANDINGAN SEBELUM DAN SESUDAH SAGA", "items": [
            {"title": "BEFORE: Manual Calls", "desc": "Orde → Payment → Inventory → Approve. Gagal di Inventory: order PENDING, payment tercatat, stok untouched. Tidak ada recovery otomatis. State inconsistency permanen."},
            {"title": "AFTER: Saga Orchestration", "desc": "Orde → Payment → Inventory(GAGAL) → Compensate Payment(refund) → Compensate Order(cancel). State pulih sempurna. Log lengkap记录了每一步状态。"},
            {"title": "BEFORE: No Observability", "desc": "Tidak ada log terpusat; harus trace masing-masing service log untuk reconstruct failure."},
            {"title": "AFTER: Centralized Step Log", "desc": "Orchestrator mencatat STATUS_EXECUTED, STATUS_FAILED, STATUS_COMPENSATED untuk setiap langkah. Debugging satu command: orch.Logs()."}
        ]}
    ],
    narrative="Perbandingan sebelum dan sesudah menunjukkan transformasi fundamental dari state management manual ke orchestration terstruktur. Sebelum saga, kegagalan di satu service meninggalkan state inconsistency yang memerlukan intervensi manual. Setelah saga, orchestrator secara otomatis mencatut log status dan memicu kompensasi berurutan. Setiap langkah memiliki jejak audit trail yang memungkinkan debugging cepat tanpa perlu menggabungkan log dari multiple service.",
    narrative_en="The before-after comparison shows a fundamental transformation from manual state management to structured orchestration. Before saga, a failure in one service leaves persistent state inconsistency requiring manual intervention. After saga, the orchestrator automatically logs status and triggers sequential compensation. Each step has an audit trail enabling fast debugging without combining logs from multiple services.",
    mental_model="Orchestration replaces chaos with traceable state transitions and automatic recovery.",
    source_path="cmd/demo/main.go • tests/saga_test.go:74-147",
    slide_num="14 / 24",
    layout="cards"
)

# Slide 15: Analisis mendalam — Service Implementation
add_slide(
    categories=["DEEP DIVE", "SERVICE LAYER"],
    title="Service Layer: Semantic Lock, Idempotency, dan Atomic Stock Mutation",
    subtitle="Setiap service memiliki mekanisme proteksi mandiri sebelum saga menyentuhnya.",
    content_items=[
        {"title": "SEMANTIC LOCK (OrderService)", "items": [
            "func CreateOrder checks s.locks[orderID] — menolak duplikasi saat saga masih berjalan.",
            "ApproveOrder menghapus lock via delete(s.locks, orderID) setelah approved.",
            "CancelOrder juga menghapus lock agar order yang dikompensasi bisa dibuat ulang."
        ], "type": "bullet"},
        {"title": "IDEMPOTENCY KEY (PaymentService)", "items": [
            "processedID map tracks已processed payment IDs.",
            "func ProcessPayment: if processedID[id] == true → return nil (idempotent success).",
            "RefundPayment: delete(s.payments, id) — aman dipanggil berkali-kali."
        ], "type": "bullet"},
        {"title": "ATOMIC STOCK MUTATION (InventoryService)", "items": [
            "Reserve: atomic stock decrement under mutex — stock[item] -= qty, reserved[item] += qty.",
            "Release: atomic stock restore under mutex — reserved[item] -= qty, stock[item] += qty.",
            "Kedua operasi dilindungi sync.Mutex tunggal — nol race condition pada concurrent access."
        ], "type": "bullet"}
    ],
    narrative="Service layer bukan sekadar storage abstraction — setiap service mengandung mekanisme proteksi mandiri. OrderService menggunakan semantic lock untuk mencegah double-creation order dalam satu saga instance. PaymentService mengandalkan idempotency key untuk menangani retry berulang pada payment step. InventoryService menjamin atomicity stock mutation melalui mutex tunggal. Tanpa proteksi ini, saga kompensasi sendiri bisa menimbulkan race condition.",
    narrative_en="The service layer is not merely a storage abstraction — each service contains its own protection mechanisms. OrderService uses semantic locking to prevent double order creation within a single saga instance. PaymentService relies on idempotency keys to handle repeated retry on the payment step. InventoryService guarantees atomic stock mutation through a single mutex. Without these protections, saga compensation itself could introduce race conditions.",
    mental_model="Service-level guards complement saga-level compensation — defense in depth across layers.",
    source_path="internal/services/services.go:29-157",
    slide_num="15 / 24",
    layout="code-bullet"
)

# Slide 16: Analisis mendalam — Event Bus untuk Choreography
add_slide(
    categories=["DEEP DIVE", "EVENT BUS"],
    title="Choreography Pattern: Event Bus Decoupled via Pub/Sub",
    subtitle="Setiap layanan subscribe ke event tertentu dan publish event setelah proses selesai.",
    content_items=[
        {"type": "code", "title": "internal/saga/choreography.go — EventBus pub/sub", "items": [
            "type EventBus struct {",
            "    mu       sync.RWMutex",
            "    handlers map[EventType][]EventHandler",
            "}",
            "",
            "func (b *EventBus) Subscribe(t EventType, h EventHandler) {",
            "    b.mu.Lock()",
            "    b.handlers[t] = append(b.handlers[t], h)",
            "    b.mu.Unlock()",
            "}",
            "",
            "func (b *EventBus) Publish(ctx context.Context, e Event) {",
            "    b.mu.RLock()",
            "    handlers := append([]EventHandler(nil), b.handlers[e.Type]...)",
            "    b.mu.RUnlock()",
            "    for _, h := range handlers { h(ctx, e) }",
            "}"
        ]},
        {"type": "bullet", "title": "FLOW CHOREOGRAPHY", "items": [
            "OrderCreated event → PaymentService.ProcessPayment → PaymentCompleted/PaymentFailed",
            "PaymentCompleted → InventoryService.Reserve → InventoryReserved/InventoryFailed",
            "InventoryReserved → OrderService.ApproveOrder",
            "InventoryFailed/PaymentFailed → PaymentService.RefundPayment + OrderService.CancelOrder"
        ]}
    ],
    narrative="Choreography pattern menghilangkan koordinator terpusat sepenuhnya. Setiap layanan mendengarkan event dari layanan lain dan bereaksi sesuai handler yang terdaftar. Event bus menggunakan RWMutex untuk keamanan konkuren: Subscribe memegang write lock, Publish memegang read lock sementara menyalin handler slice. Pola ini menghasilkan loose coupling tetapi mengurangi observabilitas karena tidak ada central log eksekusi.",
    narrative_en="The choreography pattern eliminates the central coordinator entirely. Each service listens to events from other services and reacts according to registered handlers. The event bus uses RWMutex for concurrency safety: Subscribe holds write locks, Publish holds read locks while copying the handler slice. This pattern yields loose coupling but reduces observability since there is no centralized execution log.",
    mental_model="Decentralized event-driven coordination trades observability for service autonomy.",
    source_path="internal/saga/choreography.go • tests/saga_test.go:225-316",
    slide_num="16 / 24",
    layout="code-bullet"
)

# Slide 17: Bukti test mendalam — Idempotency & Semantic Lock
add_slide(
    categories=["TEST ANALYSIS", "INVARIANT PROOFS"],
    title="Bukti Invarian: Idempotency Payment dan Semantic Lock Order",
    subtitle="Dua test case mengonfirmasi bahwa state mutation tetap konsisten di bawah kondisi edge case.",
    content_items=[
        {"type": "code", "title": "tests/saga_test.go — Idempotency & Semantic Lock Proofs", "items": [
            "func TestPayment_Idempotency(t *testing.T) {",
            "    paymentSvc := services.NewPaymentService()",
            "    id := \"tx-1\"",
            "    // First call: records payment",
            "    paymentSvc.ProcessPayment(id, 100, false)",
            "    // Second call with same ID: idempotent no-op",
            "    err := paymentSvc.ProcessPayment(id, 100, false)",
            "    // Must succeed (nil error) — double charge prevented",
            "    if err != nil { t.Fatalf(...) }",
            "}",
            "",
            "func TestSemanticLock(t *testing.T) {",
            "    orderSvc := services.NewOrderService()",
            "    id := \"ord-lock\"",
            "    orderSvc.CreateOrder(id)  // sets locks[id] = true",
            "    err := orderSvc.CreateOrder(id)  // should FAIL",
            "    if err == nil { t.Fatal(\"second create should fail due to semantic lock\") }",
            "}"
        ]},
        {"type": "bullet", "title": "INARIAN YANG DIBUKTIKAN", "items": [
            "Idempotency: ProcessPayment dengan ID sama tidak menduplikasi charge.",
            "Semantic lock: CreateOrder dengan orderID yang sama ditolak selama lock masih aktif."
        ]}
    ],
    narrative="Test idempotency membuktikan bahwa payment service menerima panggilan berulang dengan ID yang sama tanpa melakukan double charge. Test semantic lock membuktikan bahwa order service menolak kreasi order duplikat selama saga instance masih berjalan — mencegah race condition antara concurrent saga instances yang menggunakan orderID sama.",
    narrative_en="The idempotency test proves that the payment service accepts repeated calls with the same ID without double-charging. The semantic lock test proves that the order service rejects duplicate order creation while the saga instance is still active — preventing race conditions between concurrent saga instances using the same orderID.",
    mental_model="Idempotency and semantic locks are the cornerstones of reliable saga step execution.",
    source_path="tests/saga_test.go:149-171",
    slide_num="17 / 24",
    layout="code"
)

# Slide 18: Bukti test — LIFO Log Verification
add_slide(
    categories=["TEST ANALYSIS", "LOG PROOFS"],
    title="Bukti LIFO Log: Verifikasi Urutan Status Step pada Kompensasi",
    subtitle="Test menegaskan log records tepat 5 entri: EXECUTED, EXECUTED, FAILED, COMPENSATED, COMPENSATED.",
    content_items=[
        {"type": "code", "title": "tests/saga_test.go — LIFO Log Verification", "items": [
            "func TestOrchestrator_FailureCompensatesLIFO(t *testing.T) {",
            "    // ...setup: CreateOrder, ProcessPayment, ReserveInventory(fail)...",
            "    err := orch.Execute(ctx)",
            "    // Verify state recovery",
            "    if orderSvc.GetOrder(orderID) != services.OrderCancelled {",
            "        t.Fatalf(\"expected order cancelled, got %v\", orderSvc.GetOrder(orderID))",
            "    }",
            "    if paymentSvc.HasPayment(orderID) {",
            "        t.Fatal(\"expected payment refunded\")",
            "    }",
            "    // Verify log order (5 entries exactly)",
            "    logs := orch.Logs()",
            "    expectedOrder := []saga.StepStatus{",
            "        saga.StatusExecuted, saga.StatusExecuted, saga.StatusFailed,",
            "        saga.StatusCompensated, saga.StatusCompensated,",
            "    }",
            "    for i, l := range logs {",
            "        if l.Status != expectedOrder[i] { ... }",
            "    }",
            "}"
        ]},
        {"type": "bullet", "title": "LOG EKSPEKTASI", "items": [
            "Idx 0: CreateOrder → EXECUTED (order PENDING)",
            "Idx 1: ProcessPayment → EXECUTED (payment recorded)",
            "Idx 2: ReserveInventory → FAILED (out of stock)",
            "Idx 3: ProcessPayment → COMPENSATED (refund executed)",
            "Idx 4: CreateOrder → COMPENSATED (order CANCELLED)"
        ]}
    ],
    narrative="Test ini memverifikasi tiga hal sekaligus: state recovery (order CANCELLED, payment refunded), stock integrity (stock tetap 1), dan urutan log yang tepat. Enam entri log — dua EXECUTED, satu FAILED, dua COMPENSATED — membentuk sequence yang membuktikan LIFO rollback bekerja sesuai desain. Urutan ini penting untuk debugging: engineer bisa merekonstruksi kronologi exact failure.",
    narrative_en="This test verifies three things simultaneously: state recovery (order CANCELLED, payment refunded), stock integrity (stock remains 1), and correct log sequence. Six log entries — two EXECUTED, one FAILED, two COMPENSATED — form a sequence proving LIFO rollback works as designed. This order is critical for debugging: engineers can reconstruct the exact failure chronology.",
    mental_model="Log sequence integrity is as important as state recovery — every compensation must be traced.",
    source_path="tests/saga_test.go:74-147",
    slide_num="18 / 24",
    layout="code"
)

# Slide 19: Bukti test — Context Cancellation
add_slide(
    categories=["TEST ANALYSIS", "CONTEXT SAFETY"],
    title="Bukti Context Cancellation: Kompensasi Saat Saga Dibatalkan",
    subtitle="Test memastikan step yang sudah dieksekusi dikompensasi saat context dibatalkan di tengah jalan.",
    content_items=[
        {"type": "code", "title": "tests/saga_test.go — Context Cancellation Test", "items": [
            "func TestOrchestrator_ContextCancellation(t *testing.T) {",
            "    ctx, cancel := context.WithCancel(context.Background())",
            "    orch := saga.NewOrchestrator()",
            "    step1Compensated := false",
            "    orch.AddStep(saga.Step{",
            "        Name: \"Step1\",",
            "        Execute: func(ctx context.Context) error {",
            "            cancel() // cancel context before step 2",
            "            return nil",
            "        },",
            "        Compensate: func(ctx context.Context) error {",
            "            step1Compensated = true",
            "            return nil",
            "        },",
            "    })",
            "    orch.AddStep(saga.Step{",
            "        Name: \"Step2\",",
            "        Execute: func(ctx context.Context) error {",
            "            t.Fatal(\"step 2 should not execute when context cancelled\")",
            "            return nil",
            "        },",
            "    })",
            "    err := orch.Execute(ctx)",
            "    if !step1Compensated { t.Fatal(\"expected step 1 compensated\") }",
            "}"
        ]},
        {"type": "bullet", "title": "INARIAN CONTEXT CANCELLATION", "items": [
            "Context cancellation saat step 1 execute memicu cancel() sebelum step 2 dimulai.",
            "Step 2 tidak pernah dieksekusi — select case ctx.Done() interrupts loop.",
            "Step 1 tetap dikompensasi — compensate() dipanggil dengan context.Background().",
            "Error yang dikembalikan: 'saga cancelled: context canceled' dengan compensation details."
        ]}
    ],
    narrative="Test ini memvalidasi perilaku saga ketika context dibatalkan di tengah eksekusi. Context cancellation dideteksi pada select case di loop Execute, menghentikan iterasi sebelum step berikutnya dijalankan. Meskipun context dibatalkan, compensate() tetap dipanggil dengan context.Background() — memastikan kompensasi tidak terhenti oleh cancellation propagation. Step yang sudah dieksekusi pasti dikompensasi, mempertahankan invarian atomicity.",
    narrative_en="This test validates saga behavior when context is cancelled mid-execution. Context cancellation is detected via the select case in the Execute loop, halting iteration before the next step runs. Despite context cancellation, compensate() is still called with context.Background() — ensuring compensation is not aborted by cancellation propagation. Every executed step is guaranteed compensation, maintaining the atomicity invariant.",
    mental_model="Context cancellation must not abandon in-flight compensations — use detached context for rollback.",
    source_path="tests/saga_test.go:353-385 • internal/saga/orchestrator.go:57-69",
    slide_num="19 / 24",
    layout="code"
)

# Slide 20: Checklist implementasi
add_slide(
    categories=["CHECKLIST", "IMPLEMENTATION GUIDE"],
    title="Checklist Implementasi Saga Pattern yang Benar",
    subtitle="Delapan poin wajib sebelum saga pattern siap deployment ke production.",
    content_items=[
        {"title": "CHECKLIST IMPLEMENTASI", "items": [
            "[x] Setiap step memiliki fungsi Compensate yang valid (tidak nil untuk non-final steps).",
            "[x] Setiap step Execute dan Compensate adalah idempotent — aman dipanggil berulang.",
            "[x] Semantic lock diterapkan pada entity yang bisa diakses concurrent saga instances.",
            "[x] Orchestrator menggunakan snapshot pattern untuk concurrent-safe step execution.",
            "[x] Context cancellation handled: compensate() dipanggil dengan context independent.",
            "[x] Step log mencatat semua status: PENDING, EXECUTED, FAILED, COMPENSATED, COMPENSATE_FAILED.",
            "[x] go test -race dijalankan dan melaporkan clean — zero data race.",
            "[x] Dead-letter queue planned untuk compensation failures yang tidak bisa otomatis recover."
        ], "type": "bullet"}
    ],
    narrative="Checklist ini merangkum delapan persyaratan kritis untuk implementasi saga pattern yang production-ready. Poin paling sering terlewat adalah idempotensi kompensasi — jika compensate() dipanggil dua kali, sistem harus tetap stabil. Poin keenam tentang dead-letter queue juga krusial: compensation failure permanen memerlukan mekanisme manual intervention yang harus direncanakan sebelum deployment.",
    narrative_en="This checklist summarizes eight critical requirements for production-ready saga pattern implementation. The most commonly overlooked point is compensation idempotency — if compensate() is called twice, the system must remain stable. The sixth point about dead-letter queues is also crucial: permanent compensation failures require manual intervention mechanisms that must be planned before deployment.",
    mental_model="Idempotent compensations and observable logs are non-negotiable for production sagas.",
    source_path="internal/saga/orchestrator.go • tests/saga_test.go • content/02-master-draft.md:143-148",
    slide_num="20 / 24",
    layout="code"
)

# Slide 21: Anti-pattern & cara reproduksi kegagalan
add_slide(
    categories=["ANTI-PATTERN", "FAILURE REPRODUCTION"],
    title="Anti-Pattern Saga: Yang Harus Dihindari di Production",
    subtitle="Empat kesalahan desain umum yang menyebabkan data corruption atau infinite rollback loops.",
    content_items=[
        {"type": "cards", "title": "ANTI-PATTERNS YANG HARUS DIHINDARI", "items": [
            {"title": "Non-Idempotent Compensate", "desc": "Compensate() memanggil external API tanpa idempotency key → refund ganda saat retry kompensasi."},
            {"title": "Circular Compensation Dependency", "desc": "Step A compensate memanggil Step B execute → infinite loop kompensasi timbal balik."},
            {"title": "Missing Compensate for Final Step", "desc": "Final step (ApproveOrder) tidak punya Compensate, tapi jika di tengah jalan di-rollback, approval permanen."},
            {"title": "Blocking Long-Running Compensate", "desc": "Compensate() memanggil synchronous HTTP ke external service → orch stall, lock held, deadlock."}
        ]}
    ],
    narrative="Empat anti-pattern ini paling umum ditemukan di implementasi saga production. Non-idempotent compensate adalah penyebab utama double-refund yang merugikan financial reconciliation. Circular dependency antara kompensasi step menciptakan infinite loop yang melumpuhkan orchestrator. Missing compensate pada final step berarti approval state permanen meskipun step lain dikompensasi. Blocking long-running compensate mengunci goroutine orchestrator dan mencegah step selanjutnya dieksekusi.",
    narrative_en="These four anti-patterns are most commonly found in production saga implementations. Non-idempotent compensate is the primary cause of double-refund that harms financial reconciliation. Circular compensation dependency creates infinite loops that paralyze the orchestrator. Missing compensate on the final step means approval state persists permanently even when other steps are compensated. Blocking long-running compensate locks the orchestrator goroutine and prevents subsequent steps from executing.",
    mental_model="Every compensation must be idempotent, non-blocking, and free from circular dependencies.",
    source_path="content/02-master-draft.md:126-132",
    slide_num="21 / 24",
    layout="cards"
)

# Slide 22: Takeaway praktis & command block dominan
add_slide(
    categories=["PRACTICAL TAKEAWAY", "COMMAND REFERENCE"],
    title="Takeaway Praktis: Command & Runbook untuk Saga Pattern",
    subtitle="Kumpulan command reproduksi dan prinsip operasional untuk saga di lingkungan production.",
    content_items=[
        {"type": "code", "title": "REPRODUKSI & OPERASIONAL COMMAND", "items": [
            "# Jalankan seluruh test suite",
            "go test -v ./...",
            "",
            "# Verifikasi race condition",
            "go test -race -v ./...",
            "",
            "# Demo interaktif (happy path + rollback)",
            "go run ./cmd/demo",
            "",
            "# Production runbook:",
            "# 1. Monitor orch.Logs() untuk STATUS_FAILED / STATUS_COMPENSATE_FAILED",
            "# 2. Alert on COMPENSATE_FAILED → trigger DLQ processing",
            "# 3. Manual reconciliation via admin console for stuck compensations",
            "# 4. Regular healthcheck: verify all in-flight sagas reached terminal state"
        ]},
        {"type": "bullet", "title": "PRINSIP OPERASIONAL", "items": [
            "Kompensasi adalah transaksi biasa — perlakukan dengan sama seriusnya seperti eksekusi forward.",
            "Setiap saga instance memiliki SagaID unik untuk tracing end-to-end.",
            "Dead letter queue diperlukan untuk compensation failure yang tidak bisa auto-retry."
        ]}
    ],
    narrative="Takeaway praktis untuk engineer yang mengadopsi saga pattern di production: command di atas adalah titik awal reproduksi mandiri. Prinsip operasional ketiga poin terakhir merupakan panduan monitoring yang harus diimplementasikan sebelum go-live. Compensation failure yang masuk ke DLQ harus diproses secara manual oleh on-call engineer dengan konsol admin — tidak ada fully-automatic resolution untuk semua kasus.",
    narrative_en="The practical takeaway for engineers adopting saga pattern in production: the commands above are the starting point for independent reproduction. The three operational principles at the bottom are monitoring guidelines that must be implemented before go-live. Compensation failures entering the DLQ must be processed manually by on-call engineers through the admin console — there is no fully automatic resolution for all cases.",
    mental_model="Runbooks matter as much as code — compensation failure handling requires operational discipline.",
    source_path="cmd/demo/main.go • tests/saga_test.go • README.md",
    slide_num="22 / 24",
    layout="code"
)

# Slide 23: Prinsip arsitektural jangka panjang
add_slide(
    categories=["ARCHITECTURAL PRINCIPLE", "LONG-TERM DESIGN"],
    title="Prinsip Arsitektural Jangka Panjang: Eventual Consistency sebagai Kontrak",
    subtitle="Saga pattern mengajarkan bahwa konsistensi eventual adalah harga yang sah untuk availability tinggi.",
    content_items=[
        {"type": "schema", "title": "PRINSIP ARSITEKTUR SAGA", "items": [
            {"phase": "1. Local Autonomy", "role": "Each Service Owns Data", "desc": "Setiap service mengelola database sendiri; tidak ada shared transaction boundary."},
            {"phase": "2. Eager Local Commit", "role": "Commit Before Next Step", "desc": "Setiap langkah committed segera ke database lokal — tidak ada pending write antar step."},
            {"phase": "3. Explicit Compensation", "role": "Undo as First-Class Operation", "desc": "Kompensasi adalah operation explicit, bukan implicit rollback — didesain, di-test, di-monitor."},
            {"phase": "4. Observable State", "role": "Logs adalah Source of Truth", "desc": "Step log adalah satu-satunya sumber kebenaran untuk debugging distributed failure."},
            {"phase": "5. Resilient Design", "role": "Accept Partial Failure", "desc": "Sistem diterima sebagai entitas yang bisa gagal parsial — compensations dirancang untuk recover."}
        ]}
    ],
    narrative="Prinsip arsitektur ini berlaku jauh melampaui implementasi lab. Sagas bukan sekadar pola teknis — mereka mewakili filosofi desain sistem terdistribusi yang menerima eventual consistency sebagai kontrak fundamental. Setiap prinsip di atas mengarahkan keputusan arsitektural: dari pemilihan database per service hingga desain monitoring dan alerting system. Engineer yang memahami ini akan menghindari godaan untuk kembali ke ACID monolitik.",
    narrative_en="These architectural principles extend far beyond the lab implementation. Sagas are not merely a technical pattern — they represent a distributed systems design philosophy that accepts eventual consistency as a fundamental contract. Each principle above guides architectural decisions: from per-service database selection to monitoring and alerting system design. Engineers who understand this will avoid the temptation to revert to monolithic ACID.",
    mental_model="Eventual consistency is not a compromise — it is the architecture of scale.",
    source_path="content/02-master-draft.md:10-12,126-135",
    slide_num="23 / 24",
    layout="schema"
)

# Slide 24: Closing slide
add_slide(
    categories=["CLOSING"],
    title="Eksperimen kecil, interpretasi yang teliti.",
    subtitle="",
    content_items=[
        {"type": "bullet", "title": "SUMBER DAYA LAB", "items": [
            f"Repo: {REPO_URL}",
            f"Commit SHA: {COMMIT_SHA[:8]}",
            "internal/saga/orchestrator.go — Orchestrator dengan LIFO compensation",
            "internal/saga/choreography.go — Event bus untuk pattern choreography",
            "internal/services/services.go — Order, Payment, Inventory dengan semantic lock & idempotency",
            "tests/saga_test.go — 9 test cases: happy path, LIFO rollback, idempotency, concurrency",
            "cmd/demo/main.go — Interactive CLI demonstration"
        ]},
        {"type": "bullet", "title": "VERIFIKASI", "items": [
            "Semua test suite: PASS (Race-Clean, Go 1.22+)",
            "Invarian terbukti: LIFO rollback, idempotency, semantic lock, concurrency safety.",
            "Batasan: in-memory simulation — belum mencakup network partition atau persistence crash."
        ]}
    ],
    mental_model="Distributed consistency is earned through deliberate design, not inherited from databases.",
    source_path=REPO_URL,
    slide_num="24 / 24",
    layout="closing"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
