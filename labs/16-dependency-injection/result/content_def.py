#!/usr/bin/env python3
"""
Publication Builder Definition for Software Engineering Lab 16: Dependency Injection.
Generates: PPTX, PDF, PNG covers, Dev.to MD, Medium HTML, Metadata, Caption, DELIVERABLES.md
Adheres strictly to the SLIDE CLEANUP OVERRIDE & VISUAL ANCHOR POLISH standards.
"""
import os
import sys

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/16-dependency-injection/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/16-dependency-injection"
COMMIT_SHA = "76f0059"
LAB_NO = "16"
SLUG = "dependency-injection"
TITLE_ID = "Dependency Injection & IoC"
SUBTITLE_ID = "Constructor Injection vs Service Locator dalam Arsitektur Go"
TITLE_EN = "Dependency Injection & IoC"
SUBTITLE_EN = "Constructor Injection vs Service Locator in Go Architecture"
AUTHOR = "Lukman / lukman-ss"
TOPIC_TAGS = "Dependency Injection • Inversion of Control • Go • Software Architecture • Clean Code"

# Custom Cover SVG: Component Containers & Inversion of Control Visualization
COVER_SVG = """<svg viewBox="0 0 850 780" fill="none" xmlns="http://www.w3.org/2000/svg">
    <!-- Outer Container -->
    <rect x="30" y="40" width="790" height="310" rx="14" fill="#0f172a" stroke="#087f76" stroke-width="2"/>
    <rect x="30" y="40" width="790" height="48" rx="14" fill="#0d2838"/>
    <circle cx="58" cy="64" r="6" fill="#f43f5e"/>
    <circle cx="78" cy="64" r="6" fill="#f59e0b"/>
    <circle cx="98" cy="64" r="6" fill="#10b981"/>
    <text x="125" y="70" fill="#5eead4" font-family="monospace" font-size="16" font-weight="bold">COMPOSITION ROOT: MAIN() DECOUPLING</text>
    
    <!-- Real Gateway Node -->
    <rect x="60" y="115" width="220" height="190" rx="10" fill="#0b1b2d" stroke="#38bdf8" stroke-width="2"/>
    <rect x="60" y="115" width="220" height="38" rx="10" fill="#0c253b"/>
    <text x="80" y="140" fill="#38bdf8" font-family="monospace" font-size="14" font-weight="bold">CONCRETE GATEWAY</text>
    <text x="80" y="180" fill="#e2e8f0" font-family="monospace" font-size="13">type RealGateway</text>
    <text x="80" y="210" fill="#94a3b8" font-family="monospace" font-size="12">Network Call: HTTP</text>
    <rect x="75" y="240" width="190" height="45" rx="6" fill="#071322"/>
    <text x="88" y="268" fill="#20c997" font-family="monospace" font-size="12">Charge(m Money) error</text>

    <!-- Arrow / Injection Flow -->
    <path d="M295 210 L355 210" stroke="#20c997" stroke-width="3" stroke-dasharray="6 4"/>
    <polygon points="365,210 350,203 350,217" fill="#20c997"/>
    <text x="290" y="195" fill="#5eead4" font-family="monospace" font-size="11" font-weight="bold">INJECTS</text>

    <!-- Interface Boundary -->
    <rect x="375" y="115" width="190" height="190" rx="10" fill="#062327" stroke="#20c997" stroke-width="2"/>
    <rect x="375" y="115" width="190" height="38" rx="10" fill="#0a373b"/>
    <text x="395" y="140" fill="#20c997" font-family="monospace" font-size="14" font-weight="bold">&lt;&lt;INTERFACE&gt;&gt;</text>
    <text x="395" y="185" fill="#ffffff" font-family="monospace" font-size="13">PaymentGateway</text>
    <text x="395" y="220" fill="#5eead4" font-family="monospace" font-size="12">Contract Boundary</text>
    <text x="395" y="260" fill="#94a3b8" font-family="monospace" font-size="11">Decoupled from infra</text>

    <!-- Arrow to Processor -->
    <path d="M575 210 L615 210" stroke="#38bdf8" stroke-width="3"/>
    <polygon points="625,210 610,203 610,217" fill="#38bdf8"/>

    <!-- Consumer Node -->
    <rect x="635" y="115" width="160" height="190" rx="10" fill="#1e293b" stroke="#f59e0b" stroke-width="2"/>
    <rect x="635" y="115" width="160" height="38" rx="10" fill="#2d2010"/>
    <text x="650" y="140" fill="#fcd34d" font-family="monospace" font-size="14" font-weight="bold">PROCESSOR</text>
    <text x="650" y="185" fill="#cbd5e1" font-family="monospace" font-size="12">Domain Logic</text>
    <text x="650" y="220" fill="#f59e0b" font-family="monospace" font-size="12">Requires Gateway</text>
    <text x="650" y="260" fill="#94a3b8" font-family="monospace" font-size="11">Zero direct 'new'</text>

    <!-- Bottom Comparison -->
    <rect x="30" y="380" width="790" height="360" rx="16" fill="#0b1320" stroke="#38bdf8" stroke-width="2"/>
    <rect x="30" y="380" width="790" height="52" rx="16" fill="#0c253b"/>
    <text x="60" y="414" fill="#7dd3fc" font-family="monospace" font-size="18" font-weight="bold">STRUCTURAL CONTRACT: CONSTRUCTOR INJECTION VS LOCATOR</text>

    <!-- Card 1: Constructor Injection -->
    <rect x="60" y="455" width="345" height="245" rx="12" fill="#062327" stroke="#10b981" stroke-width="2"/>
    <text x="85" y="495" fill="#10b981" font-family="monospace" font-size="18" font-weight="bold">CONSTRUCTOR INJECTION</text>
    <text x="85" y="530" fill="#ffffff" font-family="monospace" font-size="14">NewProcessor(g PaymentGateway)</text>
    <path d="M85 548 L375 548" stroke="#087f76" stroke-width="1"/>
    <text x="85" y="580" fill="#5eead4" font-family="monospace" font-size="13">✓ Explicit dependencies in signature</text>
    <text x="85" y="615" fill="#5eead4" font-family="monospace" font-size="13">✓ Fast unit testing with MockGateway</text>
    <text x="85" y="650" fill="#5eead4" font-family="monospace" font-size="13">✓ Immutable references at birth</text>
    <text x="85" y="680" fill="#94a3b8" font-family="monospace" font-size="11">Compile-time signature validation</text>

    <!-- Card 2: Service Locator -->
    <rect x="435" y="455" width="355" height="245" rx="12" fill="#1c101c" stroke="#f43f5e" stroke-width="2"/>
    <text x="460" y="495" fill="#f43f5e" font-family="monospace" font-size="18" font-weight="bold">SERVICE LOCATOR (ANTI-PATTERN)</text>
    <text x="460" y="530" fill="#ffffff" font-family="monospace" font-size="14">NewBadProcessor(c Container)</text>
    <path d="M460 548 L760 548" stroke="#831843" stroke-width="1"/>
    <text x="460" y="580" fill="#fda4af" font-family="monospace" font-size="13">❌ Hidden dependencies behind container</text>
    <text x="460" y="615" fill="#fda4af" font-family="monospace" font-size="13">❌ Couples domain to container API</text>
    <text x="460" y="650" fill="#fda4af" font-family="monospace" font-size="13">❌ Missing dependencies fail at runtime</text>
    <text x="460" y="680" fill="#94a3b8" font-family="monospace" font-size="11">PSR-11: Container injection discouraged</text>
</svg>"""

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

# Slide 1: Cover (Dark Navy)
add_slide(
    categories=["COVER"],
    title=TITLE_ID,
    subtitle=SUBTITLE_ID,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "Component containers and inversion of control visualization",
            "Decoupling construction from execution in Go services",
            f"Penulis: {AUTHOR}"
        ]}
    ],
    mental_model="",
    source_path="",
    slide_num="01 / 24"
)

# Slide 2: Engineering Problem
add_slide(
    categories=["PROBLEM", "ENGINEERING CONFLICT"],
    title="Hardcoded Dependency Merusak Uji Sistem",
    subtitle="Instansiasi langsung provider eksternal di dalam domain logic membuat unit test mustahil.",
    content_items=[
        {"title": "MASALAH ARSITEKTUR NYATA", "items": [
            "Processor memanggil RealGateway secara langsung di dalam method.",
            "Unit test terpaksa membuat panggilan jaringan ke gateway eksternal.",
            "Kegagalan koneksi pihak ketiga membuat build CI merah seketika."
        ], "type": "bullet"}
    ],
    narrative="Instansiasi langsung RealGateway di dalam Processor mengunci dependensi konkret ke domain. Test tidak bisa mengisolasi kegagalan gateway dari bug logika bisnis.",
    mental_model="Hardcoded dependencies force integration tests where unit tests belong.",
    source_path="internal/di/processor.go",
    slide_num="02 / 24"
)

# Slide 3: Why Naive Approach Fails
add_slide(
    categories=["WHY FAILS", "SERVICE LOCATOR"],
    title="Service Locator Menyembunyikan Kebutuhan Komponen",
    subtitle="Menyuntikkan container ke seluruh struct hanya memindahkan kopling ke runtime.",
    content_items=[
        {"title": "KELEMAHAN SERVICE LOCATOR", "items": [
            "Signature fungsi menyembunyikan dependensi asli di balik interface Container.",
            "Compiler tidak tahu apakah service yang dicari benar-benar ada di registry.",
            "Missing dependency baru memicu error saat transaksi dieksekusi di runtime."
        ], "type": "bullet"}
    ],
    narrative="Service Locator tampak fleksibel tetapi menyembunyikan kontrak dependensi. Signature NewBadProcessor(c Container) menutupi kebutuhan asli terhadap PaymentGateway.",
    mental_model="Service Locator trades compile-time visibility for hidden runtime lookups.",
    source_path="internal/di/locator.go",
    slide_num="03 / 24"
)

# Slide 4: Mental Model
add_slide(
    categories=["MENTAL MODEL", "CORE CONCEPT"],
    title="IoC Bukan Fitur, Tapi Arah Kendali",
    subtitle="Hollywood Principle: modul bisnis menerima dependensi, bukan mencarinya sendiri.",
    content_items=[
        {"title": "PRINSIP INTI", "items": [
            "Pisahkan fase konfigurasi (wiring) dari fase penggunaan (execution).",
            "Modul domain tidak boleh mengimpor paket infrastruktur konkret.",
            "Composition root bertindak sebagai perakit tunggal object graph."
        ], "type": "bullet"}
    ],
    narrative="Prinsip utama IoC adalah membalik inisiatif pemanggilan. Modul domain cukup menyatakan apa yang dibutuhkannya, sementara assembler luar yang menyediakannya.",
    mental_model="Separate configuration from use: assemble components before running logic.",
    source_path="content/02-master-draft.md",
    slide_num="04 / 24"
)

# Slide 5: System / Request Flow
add_slide(
    categories=["FLOW", "ARCHITECTURE"],
    title="Alur Request Melalui Boundary Abstraksi",
    subtitle="Panggilan proses pembayaran didelegasikan melalui interface PaymentGateway.",
    content_items=[
        {"type": "flow", "items": [
            "Client Request | Memanggil ProcessPayment(amount)",
            "Processor Domain | Validasi amount > 0 & buat Money",
            "PaymentGateway | Abstraksi kontrak Charge(Money)",
            "Real/Mock Gateway | Eksekusi jaringan atau verifikasi test"
        ]}
    ],
    narrative="Request masuk ke domain Processor, melewati guard validasi, lalu diteruskan ke PaymentGateway. Implementasi konkret tidak diketahui oleh domain.",
    mental_model="Interfaces establish clear boundaries between domain rules and external side effects.",
    source_path="internal/di/gateway.go • internal/di/processor.go",
    slide_num="05 / 24",
    layout="flow"
)

# Slide 6: Key Code Path
add_slide(
    categories=["CODE", "CONSTRUCTOR"],
    title="Constructor Membuat Dependency Terlihat",
    subtitle="Signature NewProcessor menyatakan kontrak dependensi wajib secara eksplisit.",
    content_items=[
        {"title": "PROCESSOR.GO", "items": [
            "type Processor struct {",
            "    gateway PaymentGateway",
            "}",
            "func NewProcessor(g PaymentGateway) *Processor {",
            "    return &Processor{gateway: g}",
            "}"
        ], "type": "code"}
    ],
    narrative="NewProcessor memindahkan pemilihan gateway ke caller. Domain tidak lagi membuat RealGateway sendiri, sehingga test dapat menyuntikkan mock.",
    mental_model="Explicit constructor parameters make missing dependencies impossible at compile time.",
    source_path="internal/di/processor.go:6-13",
    slide_num="06 / 24"
)

# Slide 7: Failure Mode
add_slide(
    categories=["FAILURE", "FAILURE PATH"],
    title="Jalur Kegagalan: Nil Interface Panics",
    subtitle="Go mengizinkan passing nil ke interface; runtime panics jika constructor tidak dijaga.",
    content_items=[
        {"type": "failure_path", "items": [
            "Caller Passes Nil | NewProcessor(nil) lolos kompilasi",
            "Constructor Unchecked | Struct menyimpan pointer gateway nil",
            "ProcessPayment Invoked | Logika bisnis lolos validasi amount",
            "Runtime Nil Panic | p.gateway.Charge(m) memicu panic nil pointer"
        ]}
    ],
    narrative="Constructor injection di Go tidak otomatis menolak nil tanpa pengecekan manual. Jika caller mengoper nil, panic baru meledak saat method dipanggil.",
    mental_model="Interface types permit nil values: constructors still require explicit nil guards.",
    source_path="internal/di/processor.go:11-13",
    slide_num="07 / 24",
    layout="failure_path"
)

# Slide 8: Hidden Cost
add_slide(
    categories=["COST", "HIDDEN COST"],
    title="Biaya Tersembunyi: Wiring Boilerplate",
    subtitle="Manual injection menjaga kejernihan arsitektur namun menambah kode perakitan.",
    content_items=[
        {"title": "PERBANDINGAN BIAYA", "items": [
            "Aspek | Manual Wiring (Lab Ini) | DI Framework (Wire/Dig)",
            "Dependency Graph | Terbaca di main.go | Tersembunyi di reflection / gen",
            "Compile Overhead | Zero overhead | Code-gen build step tambahan",
            "Refactoring Safety | Compiler menangkap perubahan | Bergantung pada framework container"
        ], "type": "table"}
    ],
    narrative="Manual wiring menambah boilerplate pada composition root. Ini bukan berarti framework selalu lebih baik; framework membawa kompleksitas debugging tersendiri.",
    mental_model="Boilerplate in the composition root is the explicit cost of architectural transparency.",
    source_path="cmd/demo/main.go",
    slide_num="08 / 24"
)

# Slide 9: Experiment Setup
add_slide(
    categories=["EXPERIMENT", "TEST SETUP"],
    title="Setup Uji: Isolasi Dua Jalur Arsitektur",
    subtitle="Menguji perilaku validasi dan propagasi error pada Processor dan BadProcessor.",
    content_items=[
        {"type": "metric_cards", "items": [
            {"val": "6 Test Cases", "label": "Unit Suite", "sub": "3 Processor • 3 BadProcessor", "highlight": True},
            {"val": "MockGateway", "label": "Test Double", "sub": "In-memory state inspection"},
            {"val": "Go 1.22", "label": "Race Detector", "sub": "go test -race ./..."}
        ]}
    ],
    narrative="Eksperimen membandingkan 6 skenario uji: transaksi valid, gateway error, dan input negatif pada kedua varian processor.",
    mental_model="Symmetric test coverage isolates architectural trade-offs under identical conditions.",
    source_path="tests/processor_test.go",
    slide_num="09 / 24",
    layout="metric_cards"
)

# Slide 10: Result Snapshot
add_slide(
    categories=["RESULT", "EVIDENCE"],
    title="Output Demo: Verifikasi Eksekusi Terpisah",
    subtitle="Kedua pola memproses pembayaran dengan implementasi RealGateway yang sama.",
    content_items=[
        {"title": "TERMINAL DEMO OUTPUT", "items": [
            "$ go run ./cmd/demo",
            "--- Running Constructor Injection ---",
            "RealGateway charging 100 USD",
            "--- Running Service Locator ---",
            "RealGateway charging 200 USD"
        ], "type": "code"}
    ],
    narrative="Demo membuktikan kedua pola menghasilkan output fungsional yang sama. Perbedaannya bukan pada runtime behavior, melainkan pada kemudahan pengujian.",
    mental_model="Identical runtime output does not imply identical architectural maintainability.",
    source_path="cmd/demo/main.go",
    slide_num="10 / 24"
)

# Slide 11: Result Interpretation
add_slide(
    categories=["INTERPRETATION", "EVIDENCE"],
    title="Test Memvalidasi Alur, Bukan Kesiapan Skala",
    subtitle="Keberhasilan mock membuktikan isolasi logika, bukan keandalan jaringan riil.",
    content_items=[
        {"title": "YANG DIBUKTIKAN TEST", "items": [
            "Processor mendelegasikan pembayaran ke MockGateway tanpa jaringan.",
            "Validasi amount <= 0 menolak eksekusi sebelum memanggil gateway.",
            "Race detector bersih pada seluruh skenario pengujian unit."
        ], "type": "bullet"},
        {"title": "YANG TIDAK DIBUKTIKAN", "items": [
            "Ketahanan koneksi HTTP saat gateway riil mengalami timeout.",
            "Manajemen connection pool dan lifecycle pada beban production.",
            "Keamanan race condition pada shared state antar-goroutine."
        ], "type": "bullet"}
    ],
    narrative="",
    mental_model="Passing mock tests prove logic isolation; production still requires resilience tests.",
    source_path="tests/processor_test.go",
    slide_num="11 / 24"
)

# Slide 12: Trade-off
add_slide(
    categories=["TRADE-OFF", "DECISION"],
    title="Trade-off: Fleksibilitas vs Kejelasan Kontrak",
    subtitle="Service locator memudahkan penambahan service baru tetapi menutupi kontrak.",
    content_items=[
        {"title": "CONSTRUCTOR INJECTION", "items": [
            "Kontrak terlihat jelas pada signature konstruktor.",
            "Compiler menolak caller jika argumen kurang.",
            "Constructor berisiko panjang jika kelas terlalu sibuk."
        ], "type": "bullet"},
        {"title": "SERVICE LOCATOR", "items": [
            "Signature konstruktor tetap ringkas menerima Container tunggal.",
            "Kebutuhan dependensi tersembunyi di dalam implementasi method.",
            "Kompilasi lolos meski service belum terdaftar di container."
        ], "type": "bullet"}
    ],
    narrative="",
    mental_model="Explicit constructor parameters protect architectural clarity at the cost of signature length.",
    source_path="content/02-master-draft.md",
    slide_num="12 / 24"
)

# Slide 13: Improved Approach
add_slide(
    categories=["IMPROVED", "VALUE OBJECT"],
    title="Value Object Tidak Memerlukan Injeksi",
    subtitle="Struktur data tanpa dependensi infrastruktur cukup diinstansiasi langsung.",
    content_items=[
        {"title": "DIRECT INSTANTIATION", "items": [
            "func (p *Processor) ProcessPayment(amount int) error {",
            "    if amount <= 0 {",
            "        return errors.New(\"invalid amount\")",
            "    }",
            "    m := Money{Amount: amount, Currency: \"USD\"}",
            "    return p.gateway.Charge(m)",
            "}"
        ], "type": "code"}
    ],
    narrative="Money adalah value object murni tanpa efek samping. Menyuntikkan factory untuk struct sederhana seperti Money adalah over-engineering yang tidak perlu.",
    mental_model="Reserve dependency injection for behavioral services, not pure value objects.",
    source_path="internal/di/processor.go:15-21",
    slide_num="13 / 24"
)

# Slide 14: Before vs After
add_slide(
    categories=["COMPARISON", "BEFORE AFTER"],
    title="Sebelum vs Sesudah Constructor Injection",
    subtitle="Mengubah instansiasi internal menjadi penyuntikan kontrak eksternal.",
    content_items=[
        {"title": "HARDCODED REAL GATEWAY", "items": [
            "Processor membuat new RealGateway() di method.",
            "Test terpaksa memanggil network eksternal.",
            "Gagal ganti provider tanpa mengubah domain code."
        ], "type": "before_after"},
        {"title": "CONSTRUCTOR INJECTION", "items": [
            "Processor menerima interface PaymentGateway.",
            "Test menyuntikkan MockGateway in-memory seketika.",
            "Caller bebas menukar provider di composition root."
        ], "type": "before_after"}
    ],
    narrative="",
    mental_model="Externalizing construction decouples business logic from concrete infrastructure.",
    source_path="internal/di/processor.go • cmd/demo/main.go",
    slide_num="14 / 24",
    layout="before_after"
)

# Slide 15: Boundary of the Lab
add_slide(
    categories=["BOUNDARY", "LIMITATION"],
    title="Batasan Lab: Manual Wiring In-Memory",
    subtitle="Lab ini membuktikan struktur injeksi murni tanpa framework otomatis.",
    content_items=[
        {"title": "CAKUPAN IMPLEMENTASI", "items": [
            "Manual composition root di cmd/demo/main.go.",
            "Single interface PaymentGateway dan implementasi mock.",
            "Validasi domain sebelum delegasi ke infrastruktur."
        ], "type": "bullet"},
        {"title": "TIDAK TERCAKUP DALAM LAB", "items": [
            "DI framework otomatis seperti Wire, Dig, atau Fx.",
            "Lifecycle management: Singleton, Scoped, dan Transient.",
            "Scope validation untuk mencegah memory leak antar-request."
        ], "type": "bullet"}
    ],
    narrative="",
    mental_model="Manual wiring demonstrates the core principle without the complexity of automated frameworks.",
    source_path="engineering-audit/05-gaps.md",
    slide_num="15 / 24"
)

# Slide 16: Edge Case
add_slide(
    categories=["EDGE CASE", "DEFENSE"],
    title="Kasus Batas: Penanganan Parameter Tidak Valid",
    subtitle="Validasi domain wajib memblokir eksekusi sebelum menyentuh gateway.",
    content_items=[
        {"title": "TEST INVALID AMOUNT", "items": [
            "proc := di.NewProcessor(mock)",
            "err := proc.ProcessPayment(-10)",
            "if err == nil {",
            "    t.Fatalf(\"expected error for negative amount\")",
            "}",
            "if mock.ChargedMoney.Amount != 0 {",
            "    t.Errorf(\"gateway should not have been called\")",
            "}"
        ], "type": "code"}
    ],
    narrative="Test memastikan ProcessPayment(-10) gagal seketika dan ChargedMoney.Amount tetap 0. DI tidak meniadakan validasi input di trust boundary.",
    mental_model="Validate input at domain boundaries before delegating to injected infrastructure.",
    source_path="tests/processor_test.go:47-59",
    slide_num="16 / 24"
)

# Slide 17: What the Tests Prove
add_slide(
    categories=["BUKTI", "TEST COVERAGE"],
    title="Fakta yang Dibuktikan Suite Pengujian",
    subtitle="Enam test case memverifikasi perilaku terisolasi pada kedua pola arsitektur.",
    content_items=[
        {"title": "FAKTA TERVERIFIKASI", "items": [
            "Processor dan BadProcessor sukses mendelegasikan transaksi valid.",
            "Error dari PaymentGateway dipropagasikan utuh ke caller.",
            "Input negatif ditolak tanpa memanggil gateway eksternal.",
            "go test -race tidak mendeteksi data race pada skenario uji."
        ], "type": "bullet"}
    ],
    narrative="Suite pengujian membuktikan bahwa dependency mocking berhasil memisahkan logika domain dari efek samping jaringan pada semua kasus yang diuji.",
    mental_model="Clean unit tests confirm structural wiring under specified test fixtures.",
    source_path="tests/processor_test.go",
    slide_num="17 / 24"
)

# Slide 18: What the Tests Do Not Prove
add_slide(
    categories=["BATASAN BUKTI", "DO NOT OVERCLAIM"],
    title="Batas Integritas: Apa yang Tidak Diuji",
    subtitle="Hasil unit test bersih bukan jaminan arsitektur concurrency production aman.",
    content_items=[
        {"title": "DI LUAR BUKTI LAB", "items": [
            "Bukan jaminan zero-overhead pada graf dependensi kompleks.",
            "Belum membuktikan keandalan retries saat gateway riil timeout.",
            "Race detector hanya membaca jalur pengujian yang dijalankan.",
            "Tidak membuktikan efisiensi memory lifecycle skala production."
        ], "type": "bullet"}
    ],
    narrative="Test suite yang lolos hanya membuktikan skenario mock yang dirancang. Jangan menyamakan clean unit test dengan kesiapan sistem menghadapi traffic production.",
    mental_model="Unit test assertions verify local contracts, not distributed system stability.",
    source_path="engineering-audit/05-gaps.md",
    slide_num="18 / 24"
)

# Slide 19: Production Risk
add_slide(
    categories=["RISIKO", "PRODUCTION"],
    title="Risiko: Constructor Over-Injection",
    subtitle="Constructor dengan terlalu banyak dependensi menandakan pelanggaran SRP.",
    content_items=[
        {"title": "GEJALA KERUSAKAN DESAIN", "items": [
            "Constructor menerima belasan parameter interface berbeda.",
            "Kelas mengurusi terlalu banyak domain logic (God Object).",
            "Setiap unit test memerlukan belasan mock setup yang rapuh.",
            "Solusi: pecah service menjadi komponen yang lebih granular."
        ], "type": "bullet"}
    ],
    narrative="Menumpuk banyak dependensi ke dalam satu constructor bukan salah DI, melainkan sinyal bahwa class tersebut memikul terlalu banyak tanggung jawab.",
    mental_model="A constructor with too many dependencies is a code smell for single-responsibility violation.",
    source_path="content/02-master-draft.md",
    slide_num="19 / 24"
)

# Slide 20: Implementation Checklist
add_slide(
    categories=["CHECKLIST", "IMPLEMENTASI"],
    title="Panduan Implementasi DI yang Aman",
    subtitle="Langkah praktis menerapkan dependency injection pada codebase Go.",
    content_items=[
        {"title": "CHECKLIST ARSITEKTUR", "items": [
            "Definisikan interface di sisi consumer, bukan di sisi provider.",
            "Gunakan Constructor Injection sebagai default untuk seluruh service.",
            "Tambahkan guard validasi nil pada setiap konstruktor.",
            "Pusatkan instansiasi konkret di composition root (main.go)."
        ], "type": "checklist"}
    ],
    narrative="",
    mental_model="Define interfaces at consumption boundaries and centralize wiring in the composition root.",
    source_path="internal/di/processor.go • cmd/demo/main.go",
    slide_num="20 / 24",
    layout="checklist"
)

# Slide 21: Anti-Pattern
add_slide(
    categories=["ANTI-PATTERN", "SERVICE LOCATOR"],
    title="Anti-Pattern: Passing Container ke Domain",
    subtitle="Menyuntikkan service locator merusak enkapsulasi dan transparansi kontrak.",
    content_items=[
        {"title": "ALASAN MENGHINDARI LOCATOR", "items": [
            "PSR-11 dan standar industri menegaskan: 'Users SHOULD NOT pass container'.",
            "Domain terikat pada interface container milik framework.",
            "Mocking menjadi sulit karena seluruh container harus ditiru.",
            "Kompilasi berhasil meski service target belum terdaftar."
        ], "type": "bullet"}
    ],
    narrative="Menyuntikkan container ke domain logic membuat komponen terikat pada framework dan mengaburkan kebutuhan sistem sesungguhnya.",
    mental_model="Never inject the container: pass only the specific service required by the domain.",
    source_path="internal/di/locator.go",
    slide_num="21 / 24"
)

# Slide 22: Reproduction Guide
add_slide(
    categories=["REPRODUKSI", "COMMAND"],
    title="Reproduksi Mandiri di Mesin Lokal",
    subtitle="Eksperimen dijalankan dengan Go standard library tanpa dependency eksternal.",
    content_items=[
        {"title": "COMMAND RUNNER", "items": [
            "cd labs/16-dependency-injection",
            "go test -v ./...    # validasi 6 test cases",
            "go test -race ./... # verifikasi race detector",
            "go run ./cmd/demo   # simulasi DI vs Service Locator"
        ], "type": "code"}
    ],
    narrative="Lab ini sepenuhnya zero-dependency dan dapat direproduksi langsung. Eksekusi test dan demo memverifikasi pemisahan konstruksi dan penggunaan secara deterministik.",
    mental_model="Standard library tests ensure deterministic verification without external tooling overhead.",
    source_path="README.md • go.mod • tests/",
    slide_num="22 / 24"
)

# Slide 23: Technical Takeaway
add_slide(
    categories=["TAKEAWAY", "ENGINEERING PRINCIPLE"],
    title="DI Memperlihatkan Dependency Komponen",
    subtitle="Constructor injection memindahkan keputusan perakitan ke pemanggil.",
    content_items=[
        {"title": "PRINSIP REKAYASA", "items": [
            "Constructor injection membuat kontrak ketergantungan eksplisit.",
            "Validasi nil tetap menjadi tanggung jawab konstruktor.",
            "Pemisahan konstruksi dari penggunaan memungkinkan pengujian terisolasi."
        ], "type": "bullet"}
    ],
    narrative="Nilai sejati DI bukan pada framework otomatis, melainkan pada kejelasan kontrak arsitektur. Memisahkan konstruksi dari penggunaan membuat sistem dapat diuji dan berevolusi.",
    mental_model="Dependency injection makes dependencies visible; it does not replace runtime validation.",
    source_path="content/05-key-takeaways.md",
    slide_num="23 / 24"
)

# Slide 24: Closing (Dark Navy)
add_slide(
    categories=["CLOSING"],
    title="Constructor membuat dependency terlihat; domain tetap independen.",
    subtitle="",
    content_items=[
        {"title": "MODUL LAB", "items": [
            "internal/di/gateway.go • internal/di/processor.go",
            "internal/di/locator.go • cmd/demo/main.go",
            "tests/processor_test.go • README.md"
        ], "type": "bullet"},
        {"title": "VERIFIKASI & STATUS", "items": [
            f"Sumber: commit {COMMIT_SHA[:8]}",
            "Test passing: 6 test cases, 0 failures.",
            "Race detector: clean (go test -race ./...)",
            "Manual in-memory wiring; bukan proof framework DI production completeness."
        ], "type": "bullet"}
    ],
    mental_model="Constructor injection exposes dependencies explicitly; configuration remains external to use.",
    source_path=REPO_URL,
    slide_num="24 / 24"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
