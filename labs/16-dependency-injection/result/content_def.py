#!/usr/bin/env python3
"""
Publication Builder Definition for Software Engineering Lab 16: Dependency Injection.
Follows:
- cover-rules.md (Mandatory CRT retro computer, series identity lock, no micro-text, thumbnail-proof 320x180)
- slide-rules.md / slide-content-rules.md (<=45 words prose, <=30 words on code slides, dominant visual anchors, 50% zoom test)
- quality-bar.md (sharp engineering insights, strictly bounded claims, >=12 slides with limiting phrases, no overclaims)
- article-rules.md (Dev.to 22-section mapping, Medium ID HTML 22-section mapping)
"""
import os
import sys

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/16-dependency-injection/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/16-dependency-injection"
COMMIT_SHA = "9c271ae"
LAB_NO = "16"
SLUG = "dependency-injection"

TITLE_ID = "Dependency Injection & Inversion of Control"
SUBTITLE_ID = "Memisahkan Konstruksi Objek dari Logika Domain di Go"
TITLE_EN = "Dependency Injection & Inversion of Control"
SUBTITLE_EN = "Separating Object Construction from Domain Logic in Go"

AUTHOR = "Lukman / lukman-ss"
TOPIC_TAGS = "Dependency Injection • Inversion of Control • Go • Architecture • Testing"
COVER_IMAGE_PATH = os.path.join(OUTPUT_DIR, "image-cover", "cover_art.jpg")

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
            "Retro CRT Computer with Modular Interface Cartridges",
            "Hardware Socket Boundary: Swappable Network Module vs Internal Logic",
            f"Penulis: {AUTHOR}"
        ]}
    ],
    mental_model="Separation of configuration from use is the foundation of modular software.",
    source_path="",
    slide_num="01 / 24",
    layout="cover"
)

# Slide 02: Hardcoded Dependencies Create Tight Coupling (Layout: cards)
# Limiting phrase 1: "bukan berarti"
add_slide(
    categories=["PROBLEM", "TIGHT COUPLING"],
    title="Instansiasi Langsung Mengunci Logika Domain",
    subtitle="Saat kelas membuat dependensinya sendiri, boundary infrastruktur runtuh.",
    content_items=[
        {"type": "cards", "title": "DAMPAK HARDCODED DEPENDENCY", "items": [
            {"title": "Logika Terkunci", "desc": "Domain memanggil new RealGateway() langsung; perubahan gateway memaksa edit domain."},
            {"title": "Test Terikat Jaringan", "desc": "Unit test terpaksa memanggil network eksternal nyata, menjadi lambat dan rapuh."},
            {"title": "Batas Kabur", "desc": "Domain mengurus lifecycle socket jaringan, bukan murni memproses kalkulasi bisnis."}
        ]}
    ],
    narrative="Ketika Processor membuat RealGateway di dalamnya, dependensi konkret terkunci. Ini bukan berarti sistem gagal seketika, tetapi setiap test terpaksa menjadi integration test yang lambat dan rapuh.",
    narrative_en="Hardcoding RealGateway inside Processor binds domain logic to infrastructure. This does not mean instant failure, but forces every test into a brittle, network-dependent integration run.",
    mental_model="Hardcoded instantiation locks infrastructure into domain rules.",
    source_path="content/02-master-draft.md",
    slide_num="02 / 24",
    layout="cards"
)

# Slide 03: Why Naive Approaches Fail (Layout: comparison)
# Limiting phrase 2: "tidak otomatis"
add_slide(
    categories=["WHY FAILS", "NAIVE FIX"],
    title="Global Singleton dan Flag Mengaburkan Kontrak",
    subtitle="Memakai variabel global atau flag env hanya memindahkan coupling ke runtime.",
    content_items=[
        {"type": "cards", "title": "PENDEKATAN NAIF VS RISIKO", "items": [
            {"title": "Global Gateway Singleton", "desc": "State bersama sulit diisolasi; race condition antar test saat mutating instance."},
            {"title": "Runtime If-Else Flag", "desc": "Logika if isTest penuh di kode produksi; menambah branch tanpa kontrak bersih."},
            {"title": "Sembunyikan Dependency", "desc": "Caller tidak tahu bahwa Processor membutuhkan jaringan aktif sebelum dipanggil."}
        ]}
    ],
    narrative="Mengandalkan singleton global tidak otomatis menyelesaikan coupling. Pendekatan ini menyembunyikan dependensi dari fungsi caller dan memicu kontaminasi state antar pengujian paralel.",
    narrative_en="Global singletons do not automatically solve coupling. They obscure dependencies from the caller signature and introduce state pollution across concurrent test executions.",
    mental_model="Global variables disguise coupling rather than resolving it.",
    source_path="content/02-master-draft.md",
    slide_num="03 / 24",
    layout="cards"
)

# Slide 04: Mental Model — Separation of Configuration from Use (Layout: schema)
# Limiting phrase 3: "hanya pada jalur ini"
add_slide(
    categories=["MENTAL MODEL", "SEPARATION OF CONCERNS"],
    title="Pisahkan Perakitan Objek dari Eksekusi Bisnis",
    subtitle="Fase konstruksi di luar domain; fase eksekusi hanya menerima kontrak.",
    content_items=[
        {"type": "schema", "title": "DUA FASE ARSITEKTUR", "items": [
            {"phase": "1. Configuration Phase", "role": "Composition Root", "desc": "Menginstansiasi RealGateway dan menyuntikkannya ke NewProcessor(g)."},
            {"phase": "2. Consumption Phase", "role": "Processor Runtime", "desc": "Hanya menjalankan ProcessPayment(amount) tanpa tahu vendor gateway."}
        ]}
    ],
    narrative="Fowler menegaskan pemisahan konfigurasi dari penggunaan adalah prinsip terpenting. Objek dirakit di composition root; domain hanya pada jalur ini memproses logika bisnis murni.",
    narrative_en="Separation of configuration from use isolates assembly from runtime logic. Components are wired at the composition root and execute strictly against interfaces.",
    mental_model="Separate object wiring from execution flow to enable runtime flexibility.",
    source_path="research/05-report.md",
    slide_num="04 / 24",
    layout="schema"
)

# Slide 05: Core Concept — IoC vs DI (Layout: comparison)
# Limiting phrase 4: "bukan berarti"
add_slide(
    categories=["CORE CONCEPT", "IOC VS DI"],
    title="IoC adalah Arah Kendali, DI adalah Bentuk Spesifik",
    subtitle="Inversion of Control mengontrol alur; DI menginversi perakitan dependensi.",
    content_items=[
        {"type": "cards", "title": "PERBEDAAN FUNDAMENTAL", "items": [
            {"title": "Inversion of Control (IoC)", "desc": "Prinsip luas: framework memanggil kode Anda (Hollywood Principle: Don't call us, we'll call you)."},
            {"title": "Dependency Injection (DI)", "desc": "Penerapan spesifik IoC: assembler menyuntikkan dependensi objek dari luar, bukan objek yang mencarinya."}
        ]}
    ],
    narrative="IoC adalah konsep payung pembalikan kontrol alur. Memakai DI bukan berarti mengadopsi seluruh kontrol framework, melainkan hanya menyerahkan tanggung jawab pembuatan dependensi ke assembler eksternal.",
    narrative_en="IoC inverts control flow broadly, while DI specifically externalizes dependency injection. Adopting DI does not imply handing total application control to a framework.",
    mental_model="DI is a specialized subset of IoC focused strictly on dependency binding.",
    source_path="research/05-report.md",
    slide_num="05 / 24",
    layout="cards"
)

# Slide 06: Interface Abstraction and Value Object (Layout: code)
# Limiting phrase 5: "belum membuktikan"
add_slide(
    categories=["IMPLEMENTATION", "INTERFACE"],
    title="Abstraksi Gateway dan Value Object Tanpa DI",
    subtitle="Interface membatasi infrastruktur; value object dibuat langsung.",
    content_items=[
        {"type": "code", "title": "internal/di/gateway.go", "items": [
            "type Money struct { Amount int; Currency string }",
            "",
            "type PaymentGateway interface {",
            "    Charge(m Money) error",
            "}",
            "",
            "type RealGateway struct{}",
            "func (g *RealGateway) Charge(m Money) error {",
            "    fmt.Printf(\"RealGateway charging %d %s\\n\", m.Amount, m.Currency)",
            "    return nil",
            "}"
        ]}
    ],
    narrative="PaymentGateway mendefinisikan kontrak Charge tunggal. Money adalah value object murni tanpa dependensi sehingga dibuat langsung. Adanya interface ini belum membuktikan gateway bebas kegagalan jaringan.",
    narrative_en="PaymentGateway defines a minimal contract while Money remains a simple value object. Defining an interface isolates the consumer, but does not prove network resiliency on its own.",
    mental_model="Abstract external infrastructure; instantiate pure value objects directly.",
    source_path="internal/di/gateway.go",
    slide_num="06 / 24",
    layout="code"
)

# Slide 07: Constructor Injection Implementation (Layout: code)
# Limiting phrase 6: "tidak otomatis"
add_slide(
    categories=["IMPLEMENTATION", "CONSTRUCTOR INJECTION"],
    title="Constructor Injection Membuat Dependensi Terlihat",
    subtitle="Signature NewProcessor mewajibkan gateway tersedia sejak instansiasi.",
    content_items=[
        {"type": "code", "title": "internal/di/processor.go", "items": [
            "type Processor struct { gateway PaymentGateway }",
            "",
            "func NewProcessor(g PaymentGateway) *Processor {",
            "    return &Processor{gateway: g}",
            "}",
            "",
            "func (p *Processor) ProcessPayment(amount int) error {",
            "    if amount <= 0 { return errors.New(\"invalid amount\") }",
            "    m := Money{Amount: amount, Currency: \"USD\"}",
            "    return p.gateway.Charge(m)",
            "}"
        ]}
    ],
    narrative="NewProcessor menerima PaymentGateway secara eksplisit. Dependensi terlihat di signature fungsi. Namun, parameter interface tidak otomatis menolak nil pointer tanpa validasi guard tambahan.",
    narrative_en="NewProcessor explicitly demands its PaymentGateway dependency at birth. While the dependency is transparent, an interface parameter does not automatically guard against nil values.",
    mental_model="Constructor injection exposes dependencies explicitly in the type signature.",
    source_path="internal/di/processor.go",
    slide_num="07 / 24",
    layout="code"
)

# Slide 08: Domain Pre-validation Boundary (Layout: cards)
# Limiting phrase 7: "hanya pada jalur ini"
add_slide(
    categories=["VALIDATION", "FAILURE SCENARIO"],
    title="Validasi Domain Berjalan Sebelum Panggilan Gateway",
    subtitle="Input tidak valid ditolak sebelum menyentuh lapisan jaringan eksternal.",
    content_items=[
        {"type": "cards", "title": "ALUR PENGECEKAN INPUT", "items": [
            {"title": "Input Check (amount <= 0)", "desc": "Mengembalikan error 'invalid amount' seketika tanpa menyentuh gateway."},
            {"title": "Zero Side Effect", "desc": "Mock membuktikan ChargedMoney.Amount tetap bernilai 0 saat input invalid."},
            {"title": "Batas Tanggung Jawab", "desc": "Hanya pada jalur ini input divalidasi; DI tidak menggantikan perlindungan domain."}
        ]}
    ],
    narrative="DI tidak menghilangkan kewajiban validasi. Validasi amount <= 0 memastikan gateway tidak pernah dipanggil sia-sia. Hanya pada jalur ini domain bertanggung jawab atas integritas data angka.",
    narrative_en="DI does not replace input validation. Pre-checking amount <= 0 ensures infrastructure is never invoked on invalid data, guarding against unnecessary network trips.",
    mental_model="Validate input invariants before delegating to external dependencies.",
    source_path="internal/di/processor.go:15-18",
    slide_num="08 / 24",
    layout="cards"
)

# Slide 09: Service Locator Anti-Pattern (Layout: code)
# Limiting phrase 8: "bukan berarti"
add_slide(
    categories=["ANTI-PATTERN", "SERVICE LOCATOR"],
    title="Service Locator Menyembunyikan Dependensi Nyata",
    subtitle="Menyuntikkan container mengikat domain pada API lookup dinamis.",
    content_items=[
        {"type": "code", "title": "internal/di/locator.go", "items": [
            "type Container interface { GetPaymentGateway() PaymentGateway }",
            "",
            "type BadProcessor struct { container Container }",
            "",
            "func NewBadProcessor(c Container) *BadProcessor {",
            "    return &BadProcessor{container: c}",
            "}",
            "",
            "func (p *BadProcessor) ProcessPayment(amount int) error {",
            "    if amount <= 0 { return errors.New(\"invalid amount\") }",
            "    m := Money{Amount: amount, Currency: \"USD\"}",
            "    return p.container.GetPaymentGateway().Charge(m)",
            "}"
        ]}
    ],
    narrative="BadProcessor menerima Container alih-alih PaymentGateway. Kode tetap berjalan, tetapi bukan berarti arsitekturnya bersih: dependensi PaymentGateway tersembunyi di balik container lookup.",
    narrative_en="BadProcessor takes a generic Container instead of PaymentGateway. The code compiles, but this does not mean it is decoupled: actual dependencies are hidden behind service lookups.",
    mental_model="Injecting the container turns dependency management into opaque runtime queries.",
    source_path="internal/di/locator.go",
    slide_num="09 / 24",
    layout="code"
)

# Slide 10: Contrast Matrix — DI vs Service Locator (Layout: table)
add_slide(
    categories=["COMPARISON", "TRADE-OFF"],
    title="Perbandingan Struktural DI vs Service Locator",
    subtitle="Evaluasi keterbacaan kontrak, kemudahan testing, dan coupling API.",
    content_items=[
        {"title": "DI VS SERVICE LOCATOR", "items": [
            "Dimensi | Constructor Injection | Service Locator",
            "Kontrak Signature | Eksplisit (PaymentGateway) | Opaque (Container)",
            "Ketergantungan API | Interface Domain Murni | Terikat Container Interface",
            "Mocking di Test | Inject Mock Langsung | Wajib Mock Container & Service",
            "Deteksi Compile | Error jika interface hilang | Bisa lolos tapi lookup runtime"
        ], "type": "table"}
    ],
    narrative="Tabel ini membandingkan dua pendekatan: DI membuat kebutuhan dependensi transparan pada konstruktor, sedangkan Service Locator memaksa pengujian membuat mock berlapis untuk container.",
    narrative_en="Constructor DI makes component dependencies transparent at compile-time, whereas Service Locator forces tests to construct artificial mock containers to satisfy queries.",
    mental_model="Prefer explicit parameters over general-purpose container lookups.",
    source_path="research/05-report.md",
    slide_num="10 / 24",
    layout="table"
)

# Slide 11: Composition Root in main.go (Layout: code)
# Limiting phrase 9: "simulasi ini tidak sama dengan"
add_slide(
    categories=["COMPOSITION ROOT", "MANUAL WIRING"],
    title="Composition Root Menjadi Titik Tunggal Perakitan",
    subtitle="cmd/demo/main.go merakit gateway konkret dan menyuntikkannya ke prosesor.",
    content_items=[
        {"type": "code", "title": "cmd/demo/main.go", "items": [
            "func main() {",
            "    realGateway := &di.RealGateway{}",
            "",
            "    // Constructor Injection (Bersih & Eksplisit)",
            "    processor := di.NewProcessor(realGateway)",
            "    _ = processor.ProcessPayment(100)",
            "",
            "    // Service Locator (Anti-Pattern Demo)",
            "    container := &SimpleContainer{gateway: realGateway}",
            "    badProcessor := di.NewBadProcessor(container)",
            "    _ = badProcessor.ProcessPayment(200)",
            "}"
        ]}
    ],
    narrative="main.go bertindak sebagai composition root. RealGateway dibuat satu kali lalu disuntikkan. Simulasi ini tidak sama dengan container enterprise kompleks; ini manual wiring murni Go.",
    narrative_en="main.go serves as the application composition root. This simulation is not equivalent to an enterprise DI container; it demonstrates pure, explicit manual wiring.",
    mental_model="Keep instantiation isolated to the application entrypoint.",
    source_path="cmd/demo/main.go",
    slide_num="11 / 24",
    layout="code"
)

# Slide 12: Test Double Isolation (Layout: code)
# Limiting phrase 10: "cakupannya mengikuti test"
add_slide(
    categories=["TESTING", "MOCK DOUBLE"],
    title="Test Double Menggantikan Infrastruktur Jaringan",
    subtitle="MockGateway mencatat transaksi di memori tanpa panggilan HTTP.",
    content_items=[
        {"type": "code", "title": "tests/processor_test.go", "items": [
            "type MockGateway struct {",
            "    ChargedMoney di.Money",
            "    ShouldFail   bool",
            "}",
            "",
            "func (m *MockGateway) Charge(money di.Money) error {",
            "    if m.ShouldFail {",
            "        return errors.New(\"gateway unavailable\")",
            "    }",
            "    m.ChargedMoney = money",
            "    return nil",
            "}"
        ]}
    ],
    narrative="MockGateway menggantikan RealGateway sepenuhnya saat pengujian. Variabel ShouldFail memicu simulasi error jaringan secara deterministik, namun cakupannya mengikuti test skenario sederhana ini.",
    narrative_en="MockGateway stands in for the concrete gateway during unit tests. The ShouldFail flag simulates transport failures deterministically within the scoped test conditions.",
    mental_model="Test doubles isolate business logic verification from network dependencies.",
    source_path="tests/processor_test.go:10-23",
    slide_num="12 / 24",
    layout="code"
)

# Slide 13: Verified Test Suite Execution (Layout: terminal)
add_slide(
    categories=["EVIDENCE", "TEST VERIFICATION"],
    title="Enam Kasus Uji PASS Termasuk Race Detector",
    subtitle="Pengujian membuktikan isolasi sukses, propagasi error, dan penolakan input.",
    content_items=[
        {"type": "code", "title": "TERMINAL OUTPUT: go test -v -race ./...", "items": [
            "=== RUN   TestProcessor_Success           --- PASS (0.00s)",
            "=== RUN   TestProcessor_GatewayError      --- PASS (0.00s)",
            "=== RUN   TestProcessor_InvalidAmount     --- PASS (0.00s)",
            "=== RUN   TestBadProcessor_Success        --- PASS (0.00s)",
            "=== RUN   TestBadProcessor_GatewayError   --- PASS (0.00s)",
            "=== RUN   TestBadProcessor_InvalidAmount  --- PASS (0.00s)",
            "PASS",
            "ok  github.com/lukman-ss/.../tests  1.034s"
        ]}
    ],
    narrative="Semua 6 test suite berhasil dieksekusi dalam 0.00s per kasus tanpa dependensi eksternal. Mode -race memastikan tidak ada data race terdeteksi pada eksekusi test yang dijalankan.",
    narrative_en="All 6 test cases executed in 0.00s without external infrastructure. Race detector confirmed zero data races across the simulated test scenarios.",
    mental_model="Passing unit tests verify contract compliance and error handling pathways.",
    source_path="engineering/03-execution-result.md",
    slide_num="13 / 24",
    layout="code"
)

# Slide 14: Evidence Boundary: What Was and Wasn't Proven (Layout: comparison)
# Limiting phrase 11: "belum membuktikan"
add_slide(
    categories=["AUDIT", "EVIDENCE BOUNDARY"],
    title="Batas Bukti Lab: Yang Dibuktikan vs Belum",
    subtitle="Memisahkan verifikasi interface dari klaim performa container produksi.",
    content_items=[
        {"type": "cards", "title": "BATAS BUKTI RISET & KODE", "items": [
            {"title": "YANG DIBUKTIKAN", "desc": "Substitusi mock tanpa sentuh domain; input invalid cegah panggilan gateway; race clean pada test."},
            {"title": "YANG TIDAK DIBUKTIKAN", "desc": "Lifecycle scope (singleton/transient); overhead performa reflection container; penurunan bug production."}
        ]}
    ],
    narrative="Lab ini membuktikan decouling interface dan isolasi mock. Hasil ini belum membuktikan perbandingan performa DI container otomatis berskala besar maupun pengelolaan lifecycle kompleks.",
    narrative_en="This lab proves interface decoupling and test double substitution. It has not proven automated DI container overhead or production lifecycle management strategies.",
    mental_model="Distinguish verified unit test mechanics from broad architectural claims.",
    source_path="research-audit/07-verdict.md",
    slide_num="14 / 24",
    layout="cards"
)

# Slide 15: Value Objects Should Bypass DI (Layout: cards)
# Limiting phrase 12: "bukan berarti"
add_slide(
    categories=["BEST PRACTICE", "VALUE OBJECTS"],
    title="Value Object Diinstansiasi Langsung Tanpa DI",
    subtitle="Entitas data murni tanpa dependensi eksternal tidak memerlukan container.",
    content_items=[
        {"type": "cards", "title": "HEURISTIK VALUE OBJECT LAB", "items": [
            {"title": "Money, DateTime, Address", "desc": "Struktur data murni tanpa stateful network socket diinstansiasi dengan struct literal."},
            {"title": "Bukan Standar Universal", "desc": "Daftar tipe ini merupakan heuristik internal lab, bukan standar industri mutlak."},
            {"title": "Hindari Over-Abstraction", "desc": "Memasukkan value object ke dalam container hanya menambah boilerplate kosong."}
        ]}
    ],
    narrative="Mengabstraksi semua objek ke dalam DI adalah anti-pattern over-engineering. Value object seperti Money dibuat langsung; ini bukan berarti mengorbankan prinsip decoupling arsitektur.",
    narrative_en="Injecting pure data structures creates unnecessary abstraction layers. Value objects should be instantiated directly; this does not violate core decoupling principles.",
    mental_model="Instantiate pure immutable data directly; reserve DI for external infrastructure.",
    source_path="internal/di/gateway.go:8-11",
    slide_num="15 / 24",
    layout="cards"
)

# Slide 16: Constructor Bloat & Single Responsibility (Layout: cards)
# Limiting phrase 13: "tidak otomatis"
add_slide(
    categories=["SMELL", "CONSTRUCTOR BLOAT"],
    title="Parameter Terlalu Banyak Menandakan SRP Violation",
    subtitle="Konstruktor dengan belasan argumen mengindikasikan kelas terlalu gemuk.",
    content_items=[
        {"type": "cards", "title": "DIAGNOSIS CONSTRUCTOR BLOAT", "items": [
            {"title": "Gejala Constructor Bloat", "desc": "NewService(...) menerima banyak interface; memicu konstruksi yang rapuh."},
            {"title": "Heuristik Lab vs Realita", "desc": "Batas 12 parameter adalah heuristik lab; Fowler hanya menyebut 'a lot of parameters'."},
            {"title": "Solusi Arsitektural", "desc": "Pecah kelas menjadi service lebih kecil; jangan tutup masalah dengan Service Locator."}
        ]}
    ],
    narrative="Banyak parameter konstruktor adalah alarm bahwa kelas melanggar SRP. Mengganti konstruktor gemuk dengan container tidak otomatis menyehatkan desain; solusinya adalah refactoring service.",
    narrative_en="A bloated constructor signature signals Single Responsibility violations. Hiding parameters inside a container does not automatically fix design flaws; decomposing the class is required.",
    mental_model="Decompose bloated services rather than hiding dependencies inside a locator.",
    source_path="content/02-master-draft.md:131-135",
    slide_num="16 / 24",
    layout="cards"
)

# Slide 17: PSR-11 and Industry Guidance (Layout: cards)
# Limiting phrase 14: "bukan berarti"
add_slide(
    categories=["STANDARD", "PSR-11"],
    title="Rekomendasi Industri: Hindari Menginjeksi Container",
    subtitle="Standar PSR-11 menegaskan container tidak boleh diteruskan ke dalam objek.",
    content_items=[
        {"type": "cards", "title": "STANDAR RFC 2119 & ARSITEKTUR", "items": [
            {"title": "PSR-11: SHOULD NOT", "desc": "Spesifikasi menegaskan: users SHOULD NOT pass a container into an object."},
            {"title": "Rekomendasi Kuat", "desc": "SHOULD NOT berarti rekomendasi kuat arsitektur, bukan larangan teknis (MUST NOT)."},
            {"title": "Dampak Ekosistem", "desc": "Konsisten di Spring, .NET, dan Laravel: framework menginjeksi dependency spesifik."}
        ]}
    ],
    narrative="PSR-11 melarang pengoperan container ke objek dengan klausul SHOULD NOT. Ini bukan berarti container ilegal, tetapi peringatan keras bahwa domain tidak boleh bergantung pada mekanisme lookup.",
    narrative_en="PSR-11 explicitly discourages passing containers into objects via SHOULD NOT. This is a strong architectural guideline against coupling domain entities to runtime lookup tools.",
    mental_model="Expose specific interfaces rather than leaking the container into domain boundaries.",
    source_path="research/03-evidence.md:Evidence 17",
    slide_num="17 / 24",
    layout="cards"
)

# Slide 18: Concurrency & Lifecycle Limitations (Layout: cards)
# Limiting phrase 15: "production masih perlu"
add_slide(
    categories=["PRODUCTION", "CONCURRENCY RISK"],
    title="Race Clean di Test Bukan Bukti Lifecycle Production Aman",
    subtitle="Test suite menguji stateless injection; stateful singleton butuh sinkronisasi.",
    content_items=[
        {"type": "cards", "title": "BATAS KESIMPULAN KONKURENSI", "items": [
            {"title": "Stateless Mock Lolos", "desc": "Test race bersih karena RealGateway dan MockGateway tidak berbagi mutable state."},
            {"title": "Stateful Singleton", "desc": "Production masih perlu mutex atau atomic jika dependensi menyimpan counter atau pool koneksi."},
            {"title": "Lifecycle Management", "desc": "Lab ini tidak menguji manajemen lifecycle transient vs singleton berumur panjang."}
        ]}
    ],
    narrative="Status race-clean di test memvalidasi alur fungsi yang diuji. Namun, production masih perlu validasi konkurensi terpisah jika implementasi gateway berbagi koneksi socket atau mutable cache.",
    narrative_en="A race-free unit test verifies current execution paths. However, production systems still require rigorous concurrency validation whenever dependencies share pooled connections.",
    mental_model="Unit test race-safety does not guarantee concurrency safety for stateful singletons.",
    source_path="engineering-audit/03-test-audit.md",
    slide_num="18 / 24",
    layout="cards"
)

# Slide 19: The Hidden Costs of DI (Layout: cards)
# Limiting phrase 16: "tidak otomatis"
add_slide(
    categories=["TRADE-OFF", "OVERHEAD"],
    title="Biaya Tersembunyi: Kapan DI Tidak Diperlukan?",
    subtitle="Pemisahan abstraksi menambah overhead mental, tracing file, dan konfigurasi.",
    content_items=[
        {"type": "cards", "title": "BIAYA DAN OVERHEAD DI", "items": [
            {"title": "Tracing Lebih Sulit", "desc": "Melihat pemanggil butuh navigasi interface; konstruksi terpisah jauh dari titik guna."},
            {"title": "Boilerplate Tambahan", "desc": "Membuat interface untuk kelas tunggal tanpa varian implementasi menambah file kosong."},
            {"title": "Prinsip YAGNI", "desc": "Jika sebuah modul tidak butuh swapping provider atau isolated test, DI tidak otomatis untung."}
        ]}
    ],
    narrative="DI bukan obat mujarab gratis. Memisahkan konstruksi dari penggunaan menyulitkan tracing alur kode. Menambahkan abstraksi tidak otomatis menguntungkan jika modul tersebut tidak pernah diganti.",
    narrative_en="DI introduces architectural overhead and complicates code tracing. Introducing abstractions does not automatically provide value when components have only one stable implementation.",
    mental_model="Apply DI where swapping or test isolation justifies the abstraction cost.",
    source_path="research/05-report.md:Finding 10",
    slide_num="19 / 24",
    layout="cards"
)

# Slide 20: Implementation Checklist (Layout: checklist)
add_slide(
    categories=["CHECKLIST", "GUIDELINES"],
    title="Checklist Penerapan Dependency Injection di Go",
    subtitle="Panduan langkah praktis merancang dependency injection yang aman dan bersih.",
    content_items=[
        {"type": "checklist", "title": "PRINSIP CHECKLIST", "items": [
            "Deklarasikan dependensi sebagai interface, bukan struct konkret.",
            "Gunakan constructor injection sebagai pola bawaan (default).",
            "Pusatkan instansiasi konkret hanya di Composition Root (main.go).",
            "Instansiasi value object murni secara langsung tanpa kontainer.",
            "Hindari menyuntikkan container ke dalam fungsi domain logic."
        ]}
    ],
    narrative="Checklist ini merangkum aturan main DI: interface pada consumer, constructor injection eksplisit, composition root tunggal, dan instansiasi langsung untuk tipe data murni.",
    narrative_en="This checklist enforces DI best practices: consume interfaces, require explicit constructors, localize wiring in the composition root, and bypass DI for value types.",
    mental_model="Explicit interfaces and centralized wiring maintain long-term architectural sanity.",
    source_path="content/02-master-draft.md:155-166",
    slide_num="20 / 24",
    layout="checklist"
)

# Slide 21: Anti-Patterns to Avoid (Layout: cards)
# Limiting phrase 17: "bukan berarti"
add_slide(
    categories=["ANTI-PATTERN", "PITFALLS"],
    title="Anti-Pattern Umum dalam Penerapan DI",
    subtitle="Kesalahan umum yang sering merusak manfaat pemisahan dependensi.",
    content_items=[
        {"type": "cards", "title": "HINDARI POLA INI", "items": [
            {"title": "Container as God Object", "desc": "Menjadikan container parameter umum di seluruh konstruktor kelas aplikasi."},
            {"title": "Interface Pollution", "desc": "Membuat interface 1:1 untuk setiap struct internal yang tidak memiliki varian lain."},
            {"title": "Over-Injection Value", "desc": "Menyuntikkan konfigurasi primitif atau entitas data murni ke container."}
        ]}
    ],
    narrative="Menggunakan container sebagai god object merusak encapsulation. Menerapkan interface untuk semua struct bukan berarti sistem menjadi bersih, melainkan menciptakan interface pollution yang melelahkan.",
    narrative_en="Treating the container as a universal parameter destroys encapsulation. Creating 1:1 interfaces for every struct does not mean clean code; it produces severe interface pollution.",
    mental_model="Avoid interface pollution; abstract only across genuine architectural boundaries.",
    source_path="content/02-master-draft.md:135-146",
    slide_num="21 / 24",
    layout="cards"
)

# Slide 22: Reproduction Guide (Layout: code)
add_slide(
    categories=["REPRODUCTION", "COMMANDS"],
    title="Reproduksi Mandiri di Mesin Lokal",
    subtitle="Eksperimen dijalankan dengan Go standard library tanpa library eksternal.",
    content_items=[
        {"type": "code", "title": "COMMAND RUNNER", "items": [
            "cd labs/16-dependency-injection",
            "go test -v ./...    # validasi 6 skenario pengujian",
            "go test -race ./... # verifikasi race detector",
            "go run ./cmd/demo   # simulasi manual wiring & locator"
        ]}
    ],
    narrative="Lab ini sepenuhnya mandiri tanpa dependensi pihak ketiga. Perintah ini menjalankan verifikasi happy path, skenario error jaringan, validasi input, serta simulasi manual wiring.",
    narrative_en="This lab is self-contained with zero third-party dependencies. Running these commands validates happy paths, error propagation, input pre-checks, and manual wiring mechanics.",
    mental_model="Zero-dependency implementations ensure fully reproducible technical verification.",
    source_path="README.md",
    slide_num="22 / 24",
    layout="code"
)

# Slide 23: Key Engineering Takeaways (Layout: cards)
# Limiting phrase 18: "tidak otomatis"
add_slide(
    categories=["TAKEAWAY", "CORE PRINCIPLES"],
    title="Kesimpulan Inti: Konstruksi Terpisah dari Penggunaan",
    subtitle="Dependency Injection membebaskan domain dari perakitan infrastruktur.",
    content_items=[
        {"type": "cards", "title": "TIGA PILAR REKAYASA", "items": [
            {"title": "Separation of Concerns", "desc": "Composition root merakit; domain mengeksekusi kontrak interface."},
            {"title": "Testability Alami", "desc": "Substitusi mock mengisolasi pengujian tanpa perlu infrastruktur jaringan nyata."},
            {"title": "Batas Realistis", "desc": "DI memperjelas kontrak dependensi, tetapi tidak otomatis menjamin kebenaran runtime."}
        ]}
    ],
    narrative="DI memisahkan perakitan objek dari logika domain. Interface memungkinkan isolasi pengujian secara deterministik. Namun, DI tidak otomatis menjamin keamanan runtime tanpa validasi menyeluruh.",
    narrative_en="DI cleanly decouples object construction from domain execution. While it enables deterministic mock isolation, DI does not automatically guarantee overall runtime correctness.",
    mental_model="Separate configuration from use; verify domain invariants at explicit boundaries.",
    source_path="content/05-key-takeaways.md",
    slide_num="23 / 24",
    layout="cards"
)

# Slide 24: Closing Slide (Layout: closing)
add_slide(
    categories=["CLOSING"],
    title="Pemisahan konstruksi dari penggunaan adalah investasi testability dan evolusi sistem.",
    subtitle="",
    content_items=[
        {"type": "bullet", "title": "MODUL IMPLEMENTASI", "items": [
            "internal/di/gateway.go • internal/di/processor.go • internal/di/locator.go",
            "cmd/demo/main.go • tests/processor_test.go",
            "README.md • content/02-master-draft.md"
        ]},
        {"type": "bullet", "title": "VERIFIKASI & STATUS", "items": [
            f"Commit: {COMMIT_SHA[:7]} • Status: APPROVED",
            "Test passing: 6 test cases (0.00s per case)",
            "Race detector: Clean (go test -race ./...)",
            "Manual wiring Go murni; bukan container DI berbasis refleksi enterprise."
        ]}
    ],
    mental_model="Constructor injection exposes dependencies; runtime safety remains the engineer's duty.",
    source_path=REPO_URL,
    slide_num="24 / 24",
    layout="closing"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
