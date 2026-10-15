#!/usr/bin/env python3
"""
Publication Builder for Software Engineering Lab 16: Dependency Injection.
Generates: PPTX, PDF, PNG covers, Dev.to MD, Medium HTML, Metadata, Caption, DELIVERABLES.md
Adheres strictly to the SLIDE CLEANUP OVERRIDE standard: concise, high readability, professional.
"""
import os
import sys

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/16-dependency-injection/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/16-dependency-injection"
COMMIT_SHA = "9dd0553"
LAB_NO = "16"
SLUG = "dependency-injection"
TITLE_ID = "Dependency Injection & Inversion of Control"
SUBTITLE_ID = "Memisahkan Konfigurasi dari Penggunaan: Constructor Injection vs Service Locator"
TITLE_EN = "Dependency Injection & Inversion of Control"
SUBTITLE_EN = "Separating Configuration from Use: Constructor Injection vs Service Locator"
AUTHOR = "Lukman / lukman-ss"
TOPIC_TAGS = "Dependency Injection • IoC • Go • Software Architecture • Clean Code"

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

# Slide 1: Cover Slide (Dark Navy)
add_slide(
    categories=["COVER"],
    title=TITLE_ID,
    subtitle=SUBTITLE_ID,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "Wiring Diagram: Composition Root (main) -> Processor -> PaymentGateway Interface",
            "Contrast: Explicit Constructor Injection vs Opaque Service Locator Container",
            f"Penulis: {AUTHOR}",
        ]}
    ],
    mental_model="",
    source_path="",
    slide_num="01 / 24"
)

# Slide 2: Real Problem in Production
add_slide(
    categories=["PROBLEM", "ENGINEERING CONFLICT"],
    title="Coupling Keras Merusak Testability & Evolusi",
    subtitle="Instansiasi dependensi langsung di dalam logika bisnis mengunci sistem ke infrastruktur spesifik.",
    content_items=[
        {"type": "cards", "items": [
            {"title": "Vendor Lock-in Kaku", "desc": "Panggilan new RealGateway() langsung di method bisnis membuat pergantian vendor pembayaran memicu refactor massal."},
            {"title": "Test Lambat & Rapuh", "desc": "Logika pembayaran tidak bisa diuji tanpa koneksi jaringan aktif; unit test berubah jadi integration test yang flaky."},
            {"title": "State Tersembunyi", "desc": "Konsumen tidak tahu dependensi apa yang dibutuhkan sampai runtime meledak karena nil reference atau missing env."}
        ]}
    ],
    narrative="Ketika komponen bisnis menginstansiasi dependensinya sendiri, kode menjadi terkunci pada infrastruktur konkret. Pengujian terisolasi menjadi mustahil dan pergantian vendor memicu refactor besar.",
    mental_model="Hard-coded dependencies fuse business logic directly to infrastructure.",
    source_path="content/02-master-draft.md",
    slide_num="02 / 24",
    layout="cards"
)

# Slide 3: Why Naive / Manual Review Approaches Fail
add_slide(
    categories=["WHY FAILS", "APPROACH"],
    title="Mengapa Global State & Manual Swap Gagal",
    subtitle="Pendekatan naif seperti singleton global atau service locator justru menambah kerapuhan baru.",
    content_items=[
        {"type": "cards", "items": [
            {"title": "Global State / Singleton", "desc": "Membuka race condition saat concurrent test; state bocor antar test suite tanpa isolasi."},
            {"title": "Service Locator Container", "desc": "Menyembunyikan dependensi di balik container.Get(); signature constructor menipu caller."},
            {"title": "Manual Mock Patching", "desc": "Bergantung pada monkey patching runtime yang tidak aman untuk compiled languages seperti Go."}
        ]}
    ],
    narrative="Mencoba menyelesaikan coupling dengan Service Locator atau Singleton global hanya memindahkan masalah: dependensi menjadi tersembunyi dan runtime error menggantikan compile-time error.",
    mental_model="Service Locator trades explicit compile-time coupling for hidden runtime lookup failures.",
    source_path="content/02-master-draft.md",
    slide_num="03 / 24",
    layout="cards"
)

# Slide 4: Core Mental Model — Separation of Configuration from Use
add_slide(
    categories=["MENTAL MODEL", "CORE CONCEPT"],
    title="Pemisahan Konfigurasi dari Penggunaan",
    subtitle="Martin Fowler: Memisahkan perakitan objek dari eksekusi bisnis lebih mendasar dibanding pilihan DI.",
    content_items=[
        {"type": "cards", "items": [
            {"title": "Fase 1: Konfigurasi (Composition Root)", "desc": "Objek konkret diinisialisasi dan dirakit di satu tempat (main.go). Tidak ada eksekusi domain."},
            {"title": "Fase 2: Penggunaan (Business Logic)", "desc": "Layanan bisnis murni menjalankan algoritma domain melalui interface. Bebas dari pembuatan objek eksternal."},
            {"title": "Hollywood Principle", "desc": "'Don't call us, we'll call you.' Framework atau assembler yang menyuntikkan dependensi, bukan consumer yang mencari."}
        ]}
    ],
    narrative="Inti DI adalah memisahkan dua siklus hidup: pembuatan grafik objek di satu composition root, dan eksekusi logika domain yang hanya berinteraksi melalui kontrak interface.",
    mental_model="Separate construction from consumption: objects receive dependencies, never build them.",
    source_path="content/02-master-draft.md • internal/di/processor.go",
    slide_num="04 / 24",
    layout="cards"
)

# Slide 5: System Architecture & Pipeline
add_slide(
    categories=["ARSITEKTUR", "DATA FLOW"],
    title="Alur Dependensi: Composition Root ke Gateway",
    subtitle="Kontrak interface membalikkan kendali implementasi dari consumer ke composition root.",
    content_items=[
        {"type": "flow", "items": [
            "cmd/demo/main.go | Composition root: inisialisasi RealGateway & Processor",
            "PaymentGateway Interface | Kontrak abstrak: Charge(Money) error",
            "Processor.ProcessPayment | Validasi domain amount > 0 & bungkus Money",
            "MockGateway / RealGateway | Eksekusi jaringan nyata atau mock test terisolasi"
        ]}
    ],
    narrative="Alur kerja DI menghubungkan composition root ke gateway melalui kontrak interface. Processor hanya tahu interface PaymentGateway, sehingga implementasi konkret dapat diganti sesuka hati.",
    mental_model="High-level modules depend on abstractions, decoupling business rules from network adapters.",
    source_path="cmd/demo/main.go • internal/di/processor.go",
    slide_num="05 / 24",
    layout="flow"
)

# Slide 6: Dominant Code Snippet — Constructor Injection
add_slide(
    categories=["CODE", "BEST PRACTICE"],
    title="Constructor Injection: Valid at Birth",
    subtitle="Signature konstruktor secara eksplisit menuntut dependensi kontrak yang wajib diisi.",
    content_items=[
        {"title": "INTERNAL/DI/PROCESSOR.GO", "items": [
            "type Processor struct {",
            "    gateway PaymentGateway // unexported field",
            "}",
            "",
            "func NewProcessor(g PaymentGateway) *Processor {",
            "    return &Processor{gateway: g}",
            "}",
            "",
            "func (p *Processor) ProcessPayment(amount int) error {",
            "    if amount <= 0 { return errors.New(\"invalid amount\") }",
            "    return p.gateway.Charge(Money{Amount: amount, Currency: \"USD\"})",
            "}"
        ], "type": "code"}
    ],
    narrative="Constructor Injection memastikan objek selalu berada dalam state valid sejak dibuat. Tidak ada hidden dependency, tidak ada container import, dan compile-time type safe.",
    mental_model="Explicit constructor parameters guarantee valid object state at instantiation.",
    source_path="internal/di/processor.go:6-21",
    slide_num="06 / 24"
)

# Slide 7: Failure Path & Dangling Reference
add_slide(
    categories=["FAILURE", "RUNTIME BEHAVIOR"],
    title="Failure Path: Validasi Dini Mencegah I/O Sia-sia",
    subtitle="Pre-check input memotong eksekusi sebelum memicu panggilan jaringan infrastruktur.",
    content_items=[
        {"title": "ALUR KEGAGALAN SISTEM", "items": [
            "Client Input: ProcessPayment(-10)",
            "  -> Guard Clause: if amount <= 0 (VALIDATION FAILURE)",
            "  -> Return: errors.New(\"invalid amount\")",
            "  -> Short-Circuit: Gateway.Charge() TIDAK PERNAH DIPANGGIL",
            "",
            "Gateway Down: MockGateway.ShouldFail = true",
            "  -> Input Valid: ProcessPayment(50) -> Gateway.Charge(50 USD)",
            "  -> Infrastructure Error: \"gateway unavailable\" dipropagasi utuh"
        ], "type": "code"}
    ],
    narrative="DI tidak menghilangkan validasi. Guard clause domain memastikan panggilan jaringan tidak pernah dieksekusi jika input tidak valid, sementara error infrastruktur dipropagasikan secara transparan.",
    mental_model="Domain pre-conditions protect infrastructure boundaries from invalid execution.",
    source_path="internal/di/processor.go:15-21 • tests/processor_test.go:47-59",
    slide_num="07 / 24"
)

# Slide 8: Mechanism Comparison Table
add_slide(
    categories=["KOMPARASI", "TABLE"],
    title="Tabel Komparasi: Tiga Pola Pengelolaan Dependensi",
    subtitle="Perbandingan mendasar antara Hard-coded, Service Locator, dan Constructor Injection.",
    content_items=[
        {"title": "PERBANDINGAN MEKANISME", "items": [
            "Dimensi | Hard-Coded (new) | Service Locator | Constructor Injection",
            "Visibilitas Dependensi | Terkunci internal | Opaque / Tersembunyi | Eksplisit di constructor",
            "Validasi Compile-Time | Ya (tapi kaku) | Tidak (runtime lookup) | Ya (interface-checked)",
            "Isolasi Unit Test | Sangat Sulit | Perlu Mock Container | Mudah (Mock Interface)",
            "Kopling ke Framework | Nol | Tinggi (Container API) | Nol (POPO / Standard Go)"
        ], "type": "table"}
    ],
    narrative="Tabel ini memetakan trade-off ketiga pola. Constructor Injection memberikan keseimbangan terbaik: dependensi terlihat jelas di signature tanpa mengikat kode ke framework container.",
    mental_model="Explicit interfaces beat both rigid hardcoding and opaque container lookups.",
    source_path="content/02-master-draft.md",
    slide_num="08 / 24"
)

# Slide 9: Decision Sequence
add_slide(
    categories=["DECISION", "SELECTION FLOW"],
    title="Decision Sequence: Kapan Menggunakan Pola Tertentu",
    subtitle="Tiga keputusan arsitektural untuk menentukan strategi instansiasi objek.",
    content_items=[
        {"type": "cards", "items": [
            {"title": "Layanan Ber-I/O & State", "desc": "Gunakan Constructor Injection dengan Interface. Wajib jika berinteraksi dengan DB, network API, atau filesystem."},
            {"title": "Value Object & Entity Murni", "desc": "Instansiasi langsung tanpa DI. Struct seperti Money, Coordinates, dan Address tidak punya efek samping."},
            {"title": "Dynamic Plugin Resolution", "desc": "Gunakan Factory atau Registry terkontrol. Hindari Service Locator generik di lapisan logika domain."}
        ]}
    ],
    narrative="Jangan perlakukan semua objek sama. Service dengan I/O eksternal membutuhkan Constructor Injection, sedangkan value object murni harus diinstansiasi langsung tanpa perantara.",
    mental_model="Inject infrastructure boundaries; instantiate domain value objects directly.",
    source_path="content/02-master-draft.md • internal/di/gateway.go",
    slide_num="09 / 24",
    layout="cards"
)

# Slide 10: Execution Evidence / Terminal Output
add_slide(
    categories=["EVIDENCE", "EXECUTION"],
    title="Bukti Eksekusi: Demo Wiring & Test Suite",
    subtitle="Verifikasi runtime membuktikan kedua pola bekerja fungsional, namun dengan ergonomi berbeda.",
    content_items=[
        {"title": "TERMINAL OUTPUT DEMO & TEST", "items": [
            "$ go run ./cmd/demo",
            "--- Running Constructor Injection ---",
            "RealGateway charging 100 USD",
            "--- Running Service Locator ---",
            "RealGateway charging 200 USD",
            "",
            "$ go test -v -count=1 ./...",
            "PASS: TestProcessor_Success (0.00s) | TestProcessor_GatewayError (0.00s)",
            "PASS: TestProcessor_InvalidAmount (0.00s) | TestBadProcessor_Success (0.00s)",
            "ok   lab16/tests   0.387s (all 6 tests PASS, race-clean)"
        ], "type": "code"}
    ],
    narrative="Output demo dan test suite membuktikan kedua pola berfungsi, namun Constructor Injection mencapai hasil yang sama dengan nol dependensi terhadap container.",
    mental_model="Both patterns execute identically, but constructor injection eliminates runtime container friction.",
    source_path="cmd/demo/main.go • engineering/03-execution-result.md",
    slide_num="10 / 24"
)

# Slide 11: Evidence Boundary Analysis
add_slide(
    categories=["BATASAN BUKTI", "BOUNDARY"],
    title="Analisis Batasan Bukti: Fakta vs Ekstrapolasi",
    subtitle="Memisahkan temuan empiris laboratorium dari asumsi industri yang belum teruji di lab.",
    content_items=[
        {"title": "YANG DIBUKTIKAN LAB", "items": [
            "Mock injection mengisolasi test tanpa network call (0.00s).",
            "Constructor injection menjamin valid-at-birth compile-time.",
            "Guard clause mencegah pemanggilan gateway saat input invalid.",
            "Value object (Money) dapat dipakai langsung tanpa DI."
        ], "type": "bullet"},
        {"title": "YANG TIDAK DIBUKTIKAN", "items": [
            "Lifecycle container (singleton vs scoped vs transient).",
            "Performa overhead auto-wiring framework (wire / dig).",
            "Klaim penurunan defect rate sistem production secara empiris.",
            "Batas 12 parameter constructor (hanya heuristik lab)."
        ], "type": "bullet"}
    ],
    narrative="Lab membuktikan isolasi test dan pemisahan perakitan, tetapi tidak membuktikan siklus hidup scoped/transient container maupun metrik penurunan bug di production.",
    mental_model="Scope claims strictly to verified test doubles; do not extrapolate lifecycle benchmarks.",
    source_path="engineering/02-implementation-notes.md • content/01-content-brief.md",
    slide_num="11 / 24"
)

# Slide 12: Architectural Trade-offs
add_slide(
    categories=["TRADE-OFF", "ANALYSIS"],
    title="Trade-off Arsitektur: Eksplisit vs Fleksibilitas",
    subtitle="Constructor injection mengorbankan keringkasan kode awal demi transparansi dan keandalan.",
    content_items=[
        {"title": "KEUNTUNGAN (CONSTRUCTOR DI)", "items": [
            "Transparansi total: semua kebutuhan terlihat di signature.",
            "Immutabilitas: dependensi disetel saat inisialisasi dan unexported.",
            "Compile-time safety: compiler Go menolak jika dependensi kurang."
        ], "type": "bullet"},
        {"title": "KONSEKUENSI / BIAYA", "items": [
            "Boilerplate konstruktor bertambah di composition root (main.go).",
            "Perubahan dependensi memerlukan modifikasi signature konstruktor.",
            "Tracing alur kode memerlukan navigasi manual tanpa magic wiring."
        ], "type": "bullet"}
    ],
    narrative="Constructor Injection tidak gratis: ia menambah boilerplate di entry point. Namun biaya ini sebanding dengan jaminan type-safety dan kemudahan pembacaan dependensi.",
    mental_model="Pay upfront boilerplate in main.go to buy absolute clarity and testability across domain services.",
    source_path="content/02-master-draft.md",
    slide_num="12 / 24"
)

# Slide 13: Concurrency & Validation Pipeline
add_slide(
    categories=["CONCURRENCY", "ISOLATION"],
    title="Pipeline Validasi & Keamanan Konkurensi",
    subtitle="Struct immutable dan stateless service memastikan eksekusi aman di multi-goroutine.",
    content_items=[
        {"title": "CONCURRENCY & VALIDATION PIPELINE", "items": [
            "Goroutine Client Request -> NewProcessor(mock)",
            "  -> ProcessPayment(amount) [Stateless Validation]",
            "  -> if amount <= 0 -> Instant Error (No I/O)",
            "  -> Value Object: Money{Amount, USD} [Immutable, Stack-allocated]",
            "  -> Gateway.Charge(Money) [Thread-Safe Interface Call]",
            "Race Detector Verification: go test -race ./... -> 0 warnings"
        ], "type": "code"}
    ],
    narrative="Dengan menginjeksi interface stateless dan memakai value object immutable, Processor bebas dari race condition saat diakses oleh banyak goroutine secara bersamaan.",
    mental_model="Immutable dependencies and stack-allocated value objects yield natural thread safety.",
    source_path="tests/processor_test.go • internal/di/processor.go",
    slide_num="13 / 24"
)

# Slide 14: Before vs After Comparison Cards
add_slide(
    categories=["KOMPARASI", "BEFORE AFTER"],
    title="Sebelum vs Sesudah Constructor Injection",
    subtitle="Transformasi dari tight coupling ke decoupled architecture berbasis interface.",
    content_items=[
        {"title": "SEBELUM: SERVICE LOCATOR / TIGHT", "items": [
            "Processor memanggil new RealGateway() atau container.Get().",
            "Test unit wajib menjalankan server mock HTTP atau mock container.",
            "Signature NewBadProcessor(Container) menyembunyikan kebutuhan gateway.",
            "Missing dependency meledak di runtime saat method dipanggil."
        ], "type": "before_after"},
        {"title": "SESUDAH: CONSTRUCTOR INJECTION", "items": [
            "Processor hanya bergantung pada interface PaymentGateway.",
            "Test unit menyuntikkan MockGateway in-memory tanpa network (0.00s).",
            "Signature NewProcessor(PaymentGateway) eksplisit dan terverifikasi.",
            "Missing dependency gagal di compile-time jika argumen kosong."
        ], "type": "before_after"}
    ],
    narrative="Sebelumnya, dependensi tersembunyi dan runtime error mengintai. Sesudahnya, interface membuat kontrak jelas dan compiler menjamin kelengkapan parameter sejak awal.",
    mental_model="Shift dependency discovery from hidden runtime lookups to strict compile-time contracts.",
    source_path="internal/di/locator.go • internal/di/processor.go",
    slide_num="14 / 24",
    layout="before_after"
)

# Slide 15: Deep Dive — Service Locator Anti-Pattern
add_slide(
    categories=["ANALISIS", "ANTI-PATTERN"],
    title="Bedah Anatomi: Service Locator Anti-Pattern",
    subtitle="Mengapa menyuntikkan Container ke kelas domain dianggap buruk oleh standar PSR-11 & Fowler.",
    content_items=[
        {"title": "INTERNAL/DI/LOCATOR.GO", "items": [
            "type BadProcessor struct {",
            "    container Container // Anti-Pattern: Mengikat ke container API",
            "}",
            "",
            "func (p *BadProcessor) ProcessPayment(amount int) error {",
            "    // Opaque lookup: pembaca kode tidak tahu dependensi asli",
            "    return p.container.GetPaymentGateway().Charge(m)",
            "}"
        ], "type": "code"}
    ],
    narrative="Service Locator tampak fleksibel, tetapi ia mengikat kelas bisnis ke API kontainer dan menyembunyikan dependensi sebenarnya dari siapapun yang membaca signature kelas.",
    mental_model="Container injection pollutes domain logic with infrastructure orchestration concerns.",
    source_path="internal/di/locator.go:10-25",
    slide_num="15 / 24"
)

# Slide 16: Deep Dive — Value Objects Bypass DI
add_slide(
    categories=["ANALISIS", "VALUE OBJECT"],
    title="Value Objects Tidak Boleh Menggunakan DI",
    subtitle="Membedakan service berinfrastruktur dari struktur data murni tanpa side-effect.",
    content_items=[
        {"type": "cards", "items": [
            {"title": "Karakteristik Money Struct", "desc": "Hanya membawa data (Amount, Currency). Tidak memiliki koneksi jaringan, database, atau state eksternal."},
            {"title": "Instansiasi Langsung", "desc": "Money{Amount: 100, Currency: 'USD'} dibuat langsung di dalam method bisnis tanpa melalui container atau injector."},
            {"title": "Anti-Pattern Over-Engineering", "desc": "Membuat MoneyFactory atau mendaftarkan Money ke container DI adalah abstraksi sia-sia yang membebani alokasi memori."}
        ]}
    ],
    narrative="DI dicadangkan untuk service yang berinteraksi dengan infrastruktur eksternal. Struktur data murni seperti Money cukup dibuat langsung tanpa melibatkan container.",
    mental_model="Inject behavior; instantiate state.",
    source_path="internal/di/gateway.go:5-9 • internal/di/processor.go:19",
    slide_num="16 / 24",
    layout="cards"
)

# Slide 17: Deep Dive — Isolated Mocking Strategy
add_slide(
    categories=["BUKTI", "TEST STRATEGY"],
    title="Strategi Mocking: Pengujian Cepat & Terisolasi",
    subtitle="Test double menguji interaksi kontrak tanpa memerlukan infrastruktur pihak ketiga.",
    content_items=[
        {"title": "TESTS/PROCESSOR_TEST.GO", "items": [
            "type MockGateway struct {",
            "    ChargedMoney di.Money",
            "    ShouldFail   bool",
            "}",
            "",
            "func (m *MockGateway) Charge(money di.Money) error {",
            "    if m.ShouldFail { return errors.New(\"gateway unavailable\") }",
            "    m.ChargedMoney = money",
            "    return nil",
            "}"
        ], "type": "code"}
    ],
    narrative="Dengan menginjeksi mock sederhana di unit test, kita dapat memverifikasi payload transaksi dan simulasi error tanpa menyentuh jaringan sama sekali.",
    mental_model="Lightweight in-memory test doubles verify contract adherence in microseconds.",
    source_path="tests/processor_test.go:10-22",
    slide_num="17 / 24"
)

# Slide 18: What the Tests Do Not Prove
add_slide(
    categories=["BATASAN BUKTI", "DO NOT OVERCLAIM"],
    title="Batas Integritas: Apa yang Tidak Diuji",
    subtitle="Keberhasilan test suite unit tidak menjamin kesiapan performa di skala produksi.",
    content_items=[
        {"title": "DI LUAR CAKUPAN PENGUJIAN LAB", "items": [
            "Ketahanan koneksi TCP/HTTP saat RealGateway menghadapi latency spike.",
            "Penyelesaian cyclic dependency kompleks antar puluhan service.",
            "Manajemen memori lifecycle singleton vs transient under high load.",
            "Kecepatan eksekusi compile-time vs reflection-based DI containers."
        ], "type": "bullet"}
    ],
    narrative="Suite test memverifikasi isolasi logika unit dan propagasi error. Pengujian ini tidak membuktikan ketahanan jaringan produksi atau manajemen memori container skala besar.",
    mental_model="Unit test passing confirms logical contracts; production resiliency requires integration stress tests.",
    source_path="engineering-audit/05-gaps.md",
    slide_num="18 / 24"
)

# Slide 19: Production Risk — Over-Injection & SRP Violation
add_slide(
    categories=["RISIKO", "ARCHITECTURE"],
    title="Risiko Produksi: Constructor Bloat & Over-Injection",
    subtitle="Banyaknya dependensi di konstruktor adalah 'code smell' pelanggaran Single Responsibility.",
    content_items=[
        {"type": "cards", "items": [
            {"title": "Gejala Constructor Bloat", "desc": "Konstruktor menerima 8-12 dependensi sekaligus (gateway, repo, cache, logger, metrics, mailer, queue)."},
            {"title": "Akar Masalah: SRP Violation", "desc": "Layanan tersebut melakukan terlalu banyak hal. Menambahkan parameter konstruktor hanya menutupi masalah desain."},
            {"title": "Solusi Arsitektural", "desc": "Pecah layanan monolitik menjadi komponen terfokus: pisahkan notification handler, billing ledger, dan audit logging."}
        ]}
    ],
    narrative="Jika konstruktor Anda membutuhkan belasan parameter, jangan cari framework auto-wiring untuk menyembunyikannya. Pecah kelas tersebut karena ia melanggar Single Responsibility Principle.",
    mental_model="Constructor parameter explosion is a distress signal that the class has too many responsibilities.",
    source_path="content/02-master-draft.md",
    slide_num="19 / 24",
    layout="cards"
)

# Slide 20: Implementation Checklist
add_slide(
    categories=["CHECKLIST", "IMPLEMENTASI"],
    title="Checklist Adopsi Dependency Injection",
    subtitle="Panduan langkah demi langkah menerapkan DI yang bersih pada repositori tim.",
    content_items=[
        {"title": "PANDUAN LANGKAH IMPLEMENTASI", "items": [
            "Deklarasikan dependensi sebagai interface minimal di dekat consumer.",
            "Gunakan Constructor Injection sebagai mekanisme default (valid-at-birth).",
            "Pusatkan instansiasi konkret di Composition Root (main.go).",
            "Larang passing Container atau Service Locator ke domain logic.",
            "Instansiasi Value Object (Money, DTO) secara langsung tanpa DI.",
            "Pastikan semua unit test berjalan terisolasi menggunakan Test Double."
        ], "type": "checklist"}
    ],
    narrative="Checklist ini menjadi panduan praktis untuk memastikan arsitektur DI tetap bersih, mudah diuji, dan bebas dari jebakan container lock-in.",
    mental_model="Enforce explicit constructors, isolate composition roots, and keep domain logic container-free.",
    source_path="content/02-master-draft.md • content/05-key-takeaways.md",
    slide_num="20 / 24",
    layout="checklist"
)

# Slide 21: Anti-Pattern & Reproduction
add_slide(
    categories=["ANTI-PATTERN", "REPRODUKSI"],
    title="Anti-Pattern: Jebakan Umum dalam Praktik DI",
    subtitle="Pola keliru yang sering muncul saat tim mencoba menerapkan Inversion of Control.",
    content_items=[
        {"title": "JEBAKAN ARSITEKTUR UMUM", "items": [
            "Framework Lock-in: Memasukkan decorator / tag framework ke domain core.",
            "Service Locator Terselubung: Mengoper map[string]any sebagai dependency bag.",
            "Interface Bloat: Membuat interface 1-to-1 untuk struct yang tidak pernah di-swap.",
            "Ambient Context: Menggunakan global singleton yang dimutasi saat testing."
        ], "type": "bullet"},
        {"title": "CARA REPRODUKSI MASALAH", "items": [
            "Inject mock yang mengembalikan error -> verifikasi penanganan error.",
            "Kirim amount <= 0 -> verifikasi gateway mock tidak menerima hit sama sekali."
        ], "type": "bullet"}
    ],
    narrative="Hindari mengotori logika domain dengan API framework. Gunakan interface Go yang idiomatis di sisi konsumen dan jangan membuat abstraksi prematur untuk objek yang tidak perlu di-mock.",
    mental_model="Do not invent abstractions without test or substitution pressure.",
    source_path="content/02-master-draft.md",
    slide_num="21 / 24"
)

# Slide 22: Practical Takeaway & Command Runner
add_slide(
    categories=["REPRODUKSI", "COMMAND"],
    title="Reproduksi Mandiri & Command Runner",
    subtitle="Verifikasi seluruh konsep di mesin lokal menggunakan Go standard library.",
    content_items=[
        {"title": "COMMAND RUNNER MANDIRI", "items": [
            "cd labs/16-dependency-injection",
            "go test -v ./...       # Menjalankan 6 test cases (mock isolation)",
            "go test -race ./...    # Verifikasi data race pada goroutine",
            "go run ./cmd/demo      # Menjalankan demo perbandingan kedua pola"
        ], "type": "code"}
    ],
    narrative="Eksperimen ini zero-external-dependency dan dapat direproduksi dalam waktu kurang dari satu detik pada mesin lokal manapun.",
    mental_model="Zero-dependency architecture enables sub-second reproducible verification.",
    source_path="README.md • go.mod • tests/",
    slide_num="22 / 24"
)

# Slide 23: Long-term Architectural Principle
add_slide(
    categories=["TAKEAWAY", "PRINSIP ARSITEKTUR"],
    title="Prinsip Inti: Program to Interfaces, Not Implementations",
    subtitle="Arsitektur yang tangguh dibangun di atas batas kontrak yang jelas dan modular.",
    content_items=[
        {"type": "cards", "items": [
            {"title": "Pemisahan Peran Tegas", "desc": "Logika bisnis tidak boleh tahu cara membuat dependensinya; ia hanya tahu cara memakainya."},
            {"title": "Isolasi Total Pengujian", "desc": "Unit test cepat dan deterministik hanya tercapai jika batas I/O eksternal dapat disubstitusi tanpa trik kotor."},
            {"title": "Inversi Kontrol Sejati", "desc": "Kendali perakitan diserahkan ke composition root, membebaskan domain service untuk fokus murni pada aturan bisnis."}
        ]}
    ],
    narrative="Prinsip sejati DI bukan tentang framework atau anotasi canggih, melainkan tentang memisahkan concern: biarkan composition root merakit, dan biarkan modul domain mengeksekusi.",
    mental_model="True inversion of control liberates domain services from infrastructure assembly.",
    source_path="content/05-key-takeaways.md",
    slide_num="23 / 24",
    layout="cards"
)

# Slide 24: Closing Slide (Dark Navy)
add_slide(
    categories=["CLOSING"],
    title="Pisahkan konfigurasi dari penggunaan; suntikkan kontrak, bukan container.",
    subtitle="",
    content_items=[
        {"title": "MODUL REPOSITORI", "items": [
            "internal/di/gateway.go • internal/di/processor.go • internal/di/locator.go",
            "cmd/demo/main.go • tests/processor_test.go",
            "README.md • content/02-master-draft.md"
        ], "type": "bullet"},
        {"title": "VERIFIKASI & STATUS", "items": [
            f"Commit: {COMMIT_SHA[:7]} • Status: APPROVED",
            "Test Suite: 6/6 PASS (0.00s execution per case)",
            "Race Safety: Clean (go test -race ./...)",
            "Batasan: Manual wiring demo; bukan evaluasi container lifecycle production."
        ], "type": "bullet"}
    ],
    mental_model="Program to interfaces, wire at the composition root, and let domain logic remain pure.",
    source_path=REPO_URL,
    slide_num="24 / 24"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
