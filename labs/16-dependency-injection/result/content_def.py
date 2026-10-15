#!/usr/bin/env python3
"""
Publication Builder Definition for Software Engineering Lab 16: Dependency Injection.
"""

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/16-dependency-injection/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/16-dependency-injection"
COMMIT_SHA = "3d936ba"
LAB_NO = "16"
SLUG = "dependency-injection"

TITLE_ID = "Dependency Injection & Inversion of Control"
SUBTITLE_ID = "Memisahkan Konfigurasi dari Penggunaan di Go"
TITLE_EN = "Dependency Injection & Inversion of Control"
SUBTITLE_EN = "Separating Configuration from Use in Go"

AUTHOR = "Lukman / lukman-ss"
TOPIC_TAGS = "Dependency Injection • IoC • Constructor Injection • Service Locator Anti-Pattern • Go • Mocking"
COVER_IMAGE_PATH = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/16-dependency-injection/result/image-cover/cover_art.jpg"

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

# Slide 01: Cover
add_slide(
    categories=["COVER"],
    title=TITLE_ID,
    subtitle=SUBTITLE_ID,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "Constructor Injection vs Service Locator Anti-Pattern in Go",
            "Payment Gateway Dependency Separation with Manual Wiring",
            f"Penulis: {AUTHOR}"
        ]}
    ],
    mental_model="Externalize object creation to achieve separation of configuration from use.",
    source_path="",
    slide_num="01 / 24",
    layout="cover"
)

# Slide 02: Problem Nyata
add_slide(
    categories=["PROBLEM", "TIGHT COUPLING"],
    title="Hard-Coded Dependencies Mengunci Logika Bisnis pada Infrastruktur",
    subtitle="Ketika kelas menginstansiasi dependensinya sendiri, pengujian menjadi lambat dan penggantian vendor memaksa modifikasi menyeluruh.",
    content_items=[
        {"type": "cards", "title": "DAMPAK TIGHT COUPLING", "items": [
            {"title": "Pengujian Lambat", "desc": "Setiap test memerlukan setup jaringan, database, atau payment provider riil yang flaky."},
            {"title": "Vendor Lock-in", "desc": "Pergantian payment gateway memaksa edit banyak kelas yang bergantung pada implementasi konkret."},
            {"title": "Bukan Berarti Malas", "desc": "Tim aktif menulis kode, tetapi bukan berarti produktif jika setiap perubahan memicu domino refactoring."}
        ]}
    ],
    narrative="Ketika sebuah kelas membangun dependensinya sendiri, logika bisnis menjadi terkunci pada infrastruktur eksternal. Pengujian unit jadi sulit karena harus mengandalkan status jaringan nyata yang tidak deterministik.",
    narrative_en="When a class builds its own dependencies, business logic becomes locked to external infrastructure. Unit testing becomes difficult because it must rely on non-deterministic live network states.",
    mental_model="Hard-coded dependencies couple business logic to infrastructure, breaking testability and evolvability.",
    source_path="content/02-master-draft.md:3-9",
    slide_num="02 / 24",
    layout="cards"
)

# Slide 03: Mengapa Review Manual Gagal
add_slide(
    categories=["WHY FAILS", "MANUAL REVIEW LIMITATIONS"],
    title="Code Review Manual Tidak Menangkap Dependensi Tersembunyi",
    subtitle="Dependensi yang diinstansiasi secara internal sulit dilacak hanya dengan membaca kode secara sekilas.",
    content_items=[
        {"type": "cards", "title": "KEGAGALAN APPROACH NAIF", "items": [
            {"title": "Tidak Lengkap", "desc": "Reviewer mungkin melewatkan instantiate baru di dalam method ketika fokus pada alur utama."},
            {"title": "Tidak Otomatis", "desc": "Manual review tidak dapat memverifikasi bahwa semua dependensi telah dipisahkan ke constructor."},
            {"title": "Skala Buruk", "desc": "Pada codebase besar, jejak dependensi tersebar di banyak file; bukan berarti review manual cukup."}
        ]}
    ],
    narrative="Code review manual mengandalkan perhatian manusia untuk menangkap pola coupling. Pendekatan ini tidak otomatis menjamin bahwa setiap dependensi eksternal telah diinjeksi melalui constructor.",
    narrative_en="Manual code review relies on human attention to catch coupling patterns. This approach does not automatically guarantee that every external dependency has been injected through constructors.",
    mental_model="Manual review cannot systematically enforce dependency separation across a large codebase.",
    source_path="content/02-master-draft.md:3-9",
    slide_num="03 / 24",
    layout="cards"
)

# Slide 04: Mental Model Inti
add_slide(
    categories=["MENTAL MODEL", "SEPARATION OF CONCERNS"],
    title="Separation of Configuration from Use",
    subtitle="Fase konfigurasi (composition root) terpisah dari fase penggunaan (business logic).",
    content_items=[
        {"type": "schema", "title": "DUA FASE DALAM SISTEM", "items": [
            {"phase": "1. Konfigurasi", "role": "Composition Root", "desc": "Instansiasi implementasi konkret dan penyuntikan ke dalam constructor (main.go)."},
            {"phase": "2. Penggunaan", "role": "Business Logic", "desc": "Kelas menerima dependensi melalui parameter constructor; tidak mengetahui implementasi konkret."}
        ]}
    ],
    narrative="Prinsip inti Dependency Injection adalah memisahkan pembentukan objek dari penggunaannya. Composition root bertugas merakit object graph, sementara modul bisnis hanya menggunakan dependensi yang telah tersedia.",
    narrative_en="The core principle of Dependency Injection is separating object construction from its use. The composition root assembles the object graph, while business modules only consume dependencies that are provided.",
    mental_model="Separate configuration (composition root) from use (business logic) to achieve loose coupling.",
    source_path="content/02-master-draft.md:13-17",
    slide_num="04 / 24",
    layout="schema"
)

# Slide 05: Diagram Alur / Pipeline
add_slide(
    categories=["PIPELINE", "DI FLOW"],
    title="Alur Kerja Dependency Injection Manual",
    subtitle="Dari instansiasi gateway hingga delegasi ke method Charge.",
    content_items=[
        {"type": "schema", "title": "TAHAPAN EKSEKUSI PAYMENT PROCESSOR", "items": [
            {"phase": "1. Composition Root", "role": "main.go", "desc": "Membuat realGateway := &di.RealGateway{}."},
            {"phase": "2. Constructor Injection", "role": "NewProcessor", "desc": "processor := di.NewProcessor(realGateway); gateway tersimpan sebagai field unexported."},
            {"phase": "3. Validasi Input", "role": "ProcessPayment", "desc": "Periksa amount <= 0; kembalikan error jika invalid."},
            {"phase": "4. Delegasi", "role": "gateway.Charge", "desc": "Buat Money value object dan panggil gateway.Charge(m)."},
            {"phase": "5. Propagasi Error", "role": "Error Return", "desc": "Error dari gateway dipropagasikan ke caller tanpa ditelan."}
        ]}
    ],
    narrative="Pipeline DI manual dimulai dari composition root yang menginialisasi implementasi konkret. Konstruktor menyimpan dependensi sebagai field, method bisnis memvalidasi input lalu mendelegasikan ke gateway.",
    narrative_en="Manual DI pipeline starts at the composition root instantiating concrete implementations. The constructor stores dependencies as fields, business methods validate input then delegate to the gateway.",
    mental_model="Manual wiring sequence: instantiate → inject → validate → delegate → propagate error.",
    source_path="cmd/demo/main.go:16-36",
    slide_num="05 / 24",
    layout="schema"
)

# Slide 06: Cuplikan Kode Dominan — Constructor Injection
add_slide(
    categories=["IMPLEMENTATION", "CONSTRUCTOR INJECTION"],
    title="Constructor Injection: NewProcessor Menerima Interface",
    subtitle="Processor menyimpan gateway PaymentGateway tanpa mengetahui implementasinya.",
    content_items=[
        {"type": "code", "title": "internal/di/processor.go", "items": [
            "type Processor struct {",
            "    gateway PaymentGateway",
            "}",
            "",
            "func NewProcessor(g PaymentGateway) *Processor {",
            "    return &Processor{gateway: g}",
            "}",
            "",
            "func (p *Processor) ProcessPayment(amount int) error {",
            "    if amount <= 0 {",
            "        return errors.New(\"invalid amount\")",
            "    }",
            "    m := Money{Amount: amount, Currency: \"USD\"}",
            "    return p.gateway.Charge(m)",
            "}"
        ]}
    ],
    narrative="Constructor injection menjamin Processor hanya mengenal interface PaymentGateway, bukan implementasi konkrét RealGateway. Dependensi dibuat eksplisit di signature constructor.",
    narrative_en="Constructor injection ensures Processor only knows the PaymentGateway interface, not the concrete RealGateway implementation. Dependencies are made explicit in the constructor signature.",
    mental_model="Inject dependencies via constructor parameters to enforce immutability and explicit contracts.",
    source_path="internal/di/processor.go",
    slide_num="06 / 24",
    layout="code"
)

# Slide 07: Failure Path — Service Locator
add_slide(
    categories=["FAILURE PATH", "SERVICE LOCATOR ANTI-PATTERN"],
    title="Service Locator Menyembunyikan Dependensi Sesungguhnya",
    subtitle="BadProcessor menerima Container alih-alih PaymentGateway, membuat dependensi tersembunyi.",
    content_items=[
        {"type": "cards", "title": "SKENARIO KEGAGALAN LOKATOR", "items": [
            {"title": "Dependensi Tersembunyi", "desc": "Constructor hanya menerima Container; developer tidak tahu method butuh PaymentGateway."},
            {"title": "Pencarian Dinamis", "desc": "ProcessPayment memanggil p.container.GetPaymentGateway() saat runtime, menambah indirection."},
            {"title": "Pengujian Ribet", "desc": "Test harus membuat MockContainer yang mengembalikan MockGateway, bukan MockGateway langsung."},
            {"title": "Bukan Berarti Lebih Baik", "desc": "Service Locator bukan berarti anti-pattern mutlak; PSR-11 menggunakan SHOULD NOT, bukan MUST NOT."}
        ]}
    ],
    narrative="Service Locator menyembunyikan dependensi sebenarnya di balik interface container. Hal ini mengurangi keterbacaan dan mempersulit mocking karena test harus menyiapkan container, bukan dependensi langsung.",
    narrative_en="Service Locator hides actual dependencies behind a container interface. This reduces readability and makes mocking harder because tests must prepare a container instead of the dependency directly.",
    mental_model="Service Locator obscures dependencies, complicating testing and reducing code transparency.",
    source_path="internal/di/locator.go",
    slide_num="07 / 24",
    layout="cards"
)

# Slide 08: Tabel Perbandingan Mekanisme vs Gejala
add_slide(
    categories=["COMPARISON", "DI vs SERVICE LOCATOR"],
    title="Kontras Constructor Injection vs Service Locator",
    subtitle="Perbandingan mekanisme, pengujian, dan konsekuensi arsitektur.",
    content_items=[
        {"title": "ASPEK VS KONSTRUKTOR INJECTION VS SERVICE LOCATOR", "items": [
            "Aspek | Constructor Injection | Service Locator",
            "Dependensi | Eksplisit di constructor | Tersembunyi di container",
            "Mocking | Langsung inject mock | Harus mock container + get()",
            "Kekekalan | Field immutable setelah konstruksi | Bisa berubah jika container berubah",
            "Kompleksitas | Mudah dipahami | Menambah indirection layer",
            "Standar Industri | Diproleh sebagai best practice | Disebabkan sebagai anti-pattern"
        ], "type": "table"}
    ],
    narrative="Tabel ini memetakan perbedaan mendasar antara dua pendekatan. Constructor injection menghasilkan dependensi yang eksplisit dan mudah diuji; service locator memperkenalkan lapisan tambahan yang menyulitkan isolasi.",
    narrative_en="This table maps the fundamental differences between two approaches. Constructor injection yields explicit, testable dependencies; service locator introduces an extra indirection layer that complicates isolation.",
    mental_model="Compare injection mechanisms along axes of visibility, testability, and industry acceptance.",
    source_path="internal/di/processor.go internal/di/locator.go",
    slide_num="08 / 24",
    layout="table"
)

# Slide 09: Decision Sequence — 3 Focused Decision Cards
add_slide(
    categories=["DECISION SEQUENCE", "ARCHITECTURAL CHOICES"],
    title="Keputusan Desain DI pada Lab Ini",
    subtitle="Tiga keputusan arsitektural yang menentukan struktur kode.",
    content_items=[
        {"type": "cards", "title": "KEPUTUSAN PENTING", "items": [
            {"title": "Pilih Constructor Injection", "desc": "Default modern; menjamin valid state sejak objek lahir dan dependensi eksplisit."},
            {"title": "Hindari Service Locator di Domain", "desc": "Container hanya digunakan di composition root; tidak disuntikkan ke service."},
            {"title": "Gunakan Value Object Langsung", "desc": "Money diinstansiasi langsung karena tidak memiliki dependensi infrastruktur."}
        ]}
    ],
    narrative="Lab ini memilih constructor injection sebagai pola default karena memberikan kekekalan dan kejelasan dependensi. Service Locator hanya digunakan untuk mendemonstrasikan anti-pattern, bukan sebagai solusi.",
    narrative_en="This lab chooses constructor injection as the default pattern because it provides immutability and dependency clarity. Service Locator is used only to demonstrate the anti-pattern, not as a solution.",
    mental_model="Choose constructor injection for explicit dependencies, avoid service locator in domain code, instantiate value objects directly.",
    source_path="content/02-master-draft.md:91-101",
    slide_num="09 / 24",
    layout="cards"
)

# Slide 10: Bukti Eksekusi / Terminal Output
add_slide(
    categories=["EVIDENCE", "TEST RESULTS"],
    title="6 Test Case Lolos dengan Race Detector",
    subtitle="Semua test berjalan dalam 0.00s per kasus tanpa dependency eksternal.",
    content_items=[
        {"type": "code", "title": "TERMINAL OUTPUT: go test -v -race ./...", "items": [
            "=== RUN   TestProcessor_Success",
            "--- PASS: TestProcessor_Success (0.00s)",
            "=== RUN   TestProcessor_GatewayError",
            "--- PASS: TestProcessor_GatewayError (0.00s)",
            "=== RUN   TestProcessor_InvalidAmount",
            "--- PASS: TestProcessor_InvalidAmount (0.00s)",
            "=== RUN   TestBadProcessor_Success",
            "--- PASS: TestBadProcessor_Success (0.00s)",
            "=== RUN   TestBadProcessor_GatewayError",
            "--- PASS: TestBadProcessor_GatewayError (0.00s)",
            "=== RUN   TestBadProcessor_InvalidAmount",
            "--- PASS: TestBadProcessor_InvalidAmount (0.00s)",
            "PASS",
            "ok      lab16/tests   (cached)"
        ]}
    ],
    narrative="Output terminal membuktikan seluruh enam test case lulus termasuk race detector. Test berhasil mengisolasi Processor dari gateway nyata, memverifikasi happy path, error propagation, dan validasi input.",
    narrative_en="Terminal output proves all six test cases pass including race detector. Tests successfully isolate Processor from real gateway, verifying happy path, error propagation, and input validation.",
    mental_model="Executable test suite proves isolation, error handling, and input validation without external infrastructure.",
    source_path="tests/processor_test.go",
    slide_num="10 / 24",
    layout="code"
)

# Slide 11: Analisis Batasan Bukti
add_slide(
    categories=["AUDIT", "EVIDENCE BOUNDARY"],
    title="Batasan Bukti Lab: Yang Dibuktikan vs Belum",
    subtitle="Menegaskan cakupan verifikasi terhadap skenario produksi.",
    content_items=[
        {"type": "cards", "title": "BATAS PEMBUKTIAN SISTEM", "items": [
            {"title": "YANG DIBUKTIKAN", "desc": "Isolasi mock, propagation error, validasi input, race safety pada 6 test case."},
            {"title": "YANG TIDAK DIBUKTIKAN", "desc": "Lifecycle management (singleton/scoped/transient), performance overhead container, pengurangan defect rate."},
            {"title": "BATAS KLAIM EMPIRIS", "desc": "Tidak ada data kuantitatif tentang kecepatan pengujian atau stabilitas sistem setelah penerapan DI."}
        ]}
    ],
    narrative="Lab ini membuktikan sifat isolasi dan propagasi error pada skenario sederhana. Namun belum membuktikan aspek lifecycle, kinerja, atau dampak empiris terhadap metrik rekayasa perangkat lunak.",
    narrative_en="This lab proves isolation and error propagation in simple scenarios. It does not prove lifecycle management, performance overhead, or empirical impact on software engineering metrics.",
    mental_model="Strictly separate proven unit-test behavior from unproven production-scale claims.",
    source_path="content/01-content-brief.md:13-17",
    slide_num="11 / 24",
    layout="cards"
)

# Slide 12: Trade-off Arsitektur
add_slide(
    categories=["TRADE-OFF", "COSTS & BENEFITS"],
    title="Trade-off Penerapan DI Manual vs Framework",
    subtitle="Manual wiring menambah boilerplate tetapi menjamin type safety tanpa dependensi eksternal.",
    content_items=[
        {"type": "cards", "title": "KOMPROMI REKAYASA DI", "items": [
            {"title": "Boilerplate vs Convenience", "desc": "Manual wiring memerlukan kode extra di main.go; framework mengurangi boilerplate tetapi menambah kompleksitas."},
            {"title": "Type Safety vs Flexibilitas", "desc": "Constructor injection enforced pada compile time; reflection-based framework lebih fleksibel tetapi kurang aman."},
            {"title": "Tidak Otomatis_scalable", "desc": "Pada object graph sangat besar, manual wiring menjadi sulit dikelola; bukan berarti DI framework selalu diperlukan."}
        ]}
    ],
    narrative="Penerapan DI manual menuntut penulis composition root yang jelas. Pada skala besar, overhead boilerplate dapat signifikan, namun trade-off ini sering diterima demi kejelasan dan kemudahan pengujian.",
    narrative_en="Manual DI implementation requires a clear composition root. At large scale, boilerplate overhead can become significant, but this trade-off is often accepted for clarity and testability.",
    mental_model="DI manual wiring trades boilerplate for type safety and zero external dependencies.",
    source_path="content/02-master-draft.md:76-81",
    slide_num="12 / 24",
    layout="cards"
)

# Slide 13: Concurrency / Validation Pipeline
add_slide(
    categories=["CONCURRENCY", "RACE SAFETY"],
    title="Race Safety pada Implementasi DI",
    subtitle="Test dengan race detector memastikan tidak ada shared mutable state.",
    content_items=[
        {"type": "code", "title": "tests/processor_test.go", "items": [
            "func TestProcessor_Success(t *testing.T) {",
            "    mock := &MockGateway{}",
            "    proc := di.NewProcessor(mock)",
            "    err := proc.ProcessPayment(50)",
            "    if err != nil {",
            "        t.Fatalf(\"expected nil error, got %v\", err)",
            "    }",
            "    // Race detector passes; no shared mutable state",
            "}"
        ]}
    ],
    narrative="Race detector Go memeriksa akses konkuren ke memori. Implementasi DI ini tidak menggunakan goroutine atau shared state, sehingga test berlari bersih tanpa warning race condition.",
    narrative_en="Go's race detector checks concurrent memory accesses. This DI implementation uses no goroutines or shared state, so tests run cleanly without race condition warnings.",
    mental_model="Dependency injection promotes stateless, testable components that naturally avoid race conditions.",
    source_path="tests/processor_test.go:23-35",
    slide_num="13 / 24",
    layout="code"
)

# Slide 14: Before vs After Comparison Cards
add_slide(
    categories=["COMPARISON", "BEFORE VS AFTER"],
    title="Transformasi: Dari Hard-Coded ke Constructor Injection",
    subtitle="Perbandingan kode sebelum dan sesudah penerapan DI.",
    content_items=[
        {"type": "cards", "title": "SEBELUM (TIGHT COUPLING) VS SESUDAH (DI)", "items": [
            {"title": "Sebelum: Hard-Coded", "desc": "Processor baru RealGateway di dalam method; tidak bisa di-mock tanpa mengubah kode sumber."},
            {"title": "Sesudah: DI Manual", "desc": "NewProcessor menerima PaymentGateway; mock disuntikkan langsung pada test."},
            {"title": "Bukan Berarti Sempurna", "desc": "DI bukan solusi ajaib; tetap perlu disiplin composition root dan避免 over-injection."}
        ]}
    ],
    narrative="Sebelum DI, kelas bisnis langsung menginstal dependensi kongkrit, mengunci pada satu implementasi. Setelah DI, dependensi datang dari luar, memungkinkan penggantian dan pengujian terisolasi.",
    narrative_en="Before DI, business classes directly instantiate concrete dependencies, locking to one implementation. After DI, dependencies come from outside, enabling substitution and isolated testing.",
    mental_model="Replace internal instantiation with external injection to achieve loose coupling.",
    source_path="internal/di/processor.go internal/di/locator.go",
    slide_num="14 / 24",
    layout="cards"
)

# Slide 15: Analisis Mendalam — Interface PaymentGateway
add_slide(
    categories=["ANALYSIS", "INTERFACE DESIGN"],
    title="Interface PaymentGateway: Kontrak Tunggal",
    subtitle="Definisi metode Charge yang menerima value object Money.",
    content_items=[
        {"type": "code", "title": "internal/di/gateway.go", "items": [
            "type Money struct {",
            "    Amount   int",
            "    Currency string",
            "}",
            "",
            "type PaymentGateway interface {",
            "    Charge(m Money) error",
            "}"
        ]}
    ],
    narrative="Interface PaymentGateway mendefinisikan kontrak tunggal untuk semua implementasi gateway. Money adalah value object yang mewakili uang dengan amount dan currency.",
    narrative_en="The PaymentGateway interface defines a single contract for all gateway implementations. Money is a value object representing currency with amount and currency fields.",
    mental_model="Define minimal interfaces that capture the required behavior without exposing implementation details.",
    source_path="internal/di/gateway.go",
    slide_num="15 / 24",
    layout="code"
)

# Slide 16: Analisis Mendalam — Value Object Money
add_slide(
    categories=["ANALYSIS", "VALUE OBJECT"],
    title="Money: Value Object yang Diinstansiasi Langsung",
    subtitle="Value object tidak memerlukan DI karena tidak memiliki dependensi eksternal.",
    content_items=[
        {"type": "code", "title": "internal/di/processor.go", "items": [
            "func (p *Processor) ProcessPayment(amount int) error {",
            "    if amount <= 0 {",
            "        return errors.New(\"invalid amount\")",
            "    }",
            "    m := Money{Amount: amount, Currency: \"USD\"}",
            "    return p.gateway.Charge(m)",
            "}"
        ]}
    ],
    narrative="Money dibuat langsung di dalam method karena merupakan struktur data murni tanpa dependensi infrastruktur. Menggunakan DI untuk value object akan menambah abstraksi tanpa manfaat.",
    narrative_en="Money is created directly inside the method because it is a pure data structure without infrastructure dependencies. Using DI for value objects would add abstraction without benefit.",
    mental_model="Instantiate value objects directly; reserve DI for services that depend on external infrastructure.",
    source_path="internal/di/processor.go:15-21",
    slide_num="16 / 24",
    layout="code"
)

# Slide 17: Analisis Mendalam — MockGateway
add_slide(
    categories=["ANALYSIS", "TEST DOUBLE"],
    title="MockGateway: Test Double untuk Isolasi Pengujian",
    subtitle="Mock meniru PaymentGateway tanpa panggilan jaringan nyata.",
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
    narrative="MockGateway merekam panggilan Charge dan dapat dikonfigurasi untuk gagal. Hal ini memungkinkan pengujian error propagation tanpa koneksi jaringan riil.",
    narrative_en="MockGateway records Charge calls and can be configured to fail. This enables testing error propagation without real network connections.",
    mental_model="Use mock doubles to isolate unit tests from external side effects.",
    source_path="tests/processor_test.go:10-21",
    slide_num="17 / 24",
    layout="code"
)

# Slide 18: Analisis Mendalam — Composition Root
add_slide(
    categories=["ANALYSIS", "COMPOSITION ROOT"],
    title="Composition Root: main.go sebagai Titik Penyambungan",
    subtitle="Satu-satunya tempat yang tahu implementasi konkret dan menyuntikkannya.",
    content_items=[
        {"type": "code", "title": "cmd/demo/main.go", "items": [
            "func main() {",
            "    realGateway := &di.RealGateway{}",
            "    processor := di.NewProcessor(realGateway)",
            "    fmt.Println(\"--- Running Constructor Injection ---\")",
            "    err := processor.ProcessPayment(100)",
            "    if err != nil {",
            "        fmt.Printf(\"Error: %v\\n\", err)",
            "    }",
            "}"
        ]}
    ],
    narrative="Composition root adalah fungsi main yang bertanggung jawab merakit object graph. Ia satu-satunya place yang mengimpor implementasi konkret RealGateway.",
    narrative_en="The composition root is the main function responsible for assembling the object graph. It is the only place that imports the concrete RealGateway implementation.",
    mental_model="Centralize concrete instantiation in a single composition root to maintain separation of concerns.",
    source_path="cmd/demo/main.go:16-26",
    slide_num="18 / 24",
    layout="code"
)

# Slide 19: Analisis Mendalam — Validasi Input
add_slide(
    categories=["ANALYSIS", "INPUT VALIDATION"],
    title="Validasi Input Sebelum Delegasi ke Gateway",
    subtitle="Processor menolak amount <= 0 tanpa memanggil gateway.",
    content_items=[
        {"type": "code", "title": "internal/di/processor.go", "items": [
            "func (p *Processor) ProcessPayment(amount int) error {",
            "    if amount <= 0 {",
            "        return errors.New(\"invalid amount\")",
            "    }",
            "    // ... create Money and call gateway",
            "}"
        ]}
    ],
    narrative="Validasi input dilakukan sebelum delegasi ke dependensi eksternal. Ini mencegah panggilan jaringan sia-sia dan memastikan hanya amount yang valid yang diproses.",
    narrative_en="Input validation occurs before delegating to external dependencies. This prevents unnecessary network calls and ensures only valid amounts are processed.",
    mental_model="Validate inputs early to fail fast and avoid side effects in downstream dependencies.",
    source_path="internal/di/processor.go:15-18",
    slide_num="19 / 24",
    layout="code"
)

# Slide 20: Checklist Implementasi
add_slide(
    categories=["CHECKLIST", "BEST PRACTICES"],
    title="Checklist Penerapan Dependency Injection",
    subtitle="Standar wajib sebelum menggabungkan modul bisnis dengan dependensi eksternal.",
    content_items=[
        {"type": "checklist", "title": "KRITERIA IMPLEMENTASI DI YANG BAIK", "items": [
            "Dependensi dideklarasikan sebagai interface, bukan tipe konkret.",
            "Constructor injection digunakan sebagai default untuk semua service.",
            "Composition root (main.go atau setara) adalah satu-satunya tempat yang membuat implementasi konkret.",
            "Tidak ada container.Get(...) di dalam logika domain.",
            "Value objects (Money, DateTime, Address) diinstansiasi langsung.",
            "Setiap service dapat diuji dengan mock tanpa setup infrastruktur.",
            "Validasi input dilakukan sebelum delegasi ke dependensi eksternal.",
            "Test mencakup happy path, error propagation, dan invalid input.",
            "go test -race ./... PASS tanpa warning.",
            "Heuristik lab (seperti 12 parameter) tidak dikutip sebagai standar universal."
        ]}
    ],
    narrative="Checklist ini merangkum prinsip-prinsip DI yang harus dipenuhi agar kode tetap teruji, terelaborasi, dan mudah berkembang.",
    narrative_en="This checklist summarizes the DI principles that must be met to keep code testable, decoupled, and easily evolving.",
    mental_model="Enforce dependency injection guidelines through a systematic checklist.",
    source_path="content/02-master-draft.md:155-167",
    slide_num="20 / 24",
    layout="checklist"
)

# Slide 21: Anti-Pattern & Cara Reproduksi Kegagalan
add_slide(
    categories=["ANTI-PATTERN", "PITFALLS"],
    title="Anti-Pattern Umum dalam Dependency Injection",
    subtitle="Kesalahan desain yang sering merusak manfaat DI.",
    content_items=[
        {"type": "cards", "title": "HINDARI KESALAHAN INI", "items": [
            {"title": "Service Locator di Domain", "desc": "Menyuntikkan container ke kelas bisnis; dependensi menjadi tersembunyi."},
            {"title": "Mengabstraksi Value Object", "desc": "Memasukkan Money ke dalam container DI padahal hanya struktur data."},
            {"title": "Constructor Bloat", "desc": "Konstruktor dengan banyak parameter menandakan kelas melakukan terlalu banyak hal."},
            {"title": "DI Overuse", "desc": "Menggunakan DI untuk segala hal, termasuk fungsi murni tanpa state."}
        ]}
    ],
    narrative="Service Locator adalah anti-pattern karena menyembunyikan dependensi. Value object seharusnya tidak diinject. Constructor bloat menunjukkan pelanggaran Single Responsibility Principle.",
    narrative_en="Service Locator is an anti-pattern because it hides dependencies. Value objects should not be injected. Constructor bloat indicates a violation of the Single Responsibility Principle.",
    mental_model="Identify and avoid common DI anti-patterns to preserve testability and clarity.",
    source_path="content/02-master-draft.md:135-146",
    slide_num="21 / 24",
    layout="cards"
)

# Slide 22: Takeaway Praktis & Command Block
add_slide(
    categories=["REPRODUCTION", "COMMANDS"],
    title="Reproduksi Mandiri di Mesin Lokal Menggunakan Go",
    subtitle="Verifikasi mandiri DI dengan menjalankan test suite.",
    content_items=[
        {"type": "code", "title": "COMMAND RUNNER", "items": [
            "cd labs/16-dependency-injection",
            "go test -v ./...    # jalankan semua test case",
            "go test -race ./... # verifikasi race safety",
            "go run ./cmd/demo   # eksekusi demo constructor injection & service locator"
        ]}
    ],
    narrative="Semua kode ditulis dengan Go standard library tanpa dependensi eksternal. Test suite mencakup isolasi, error propagation, dan validasi input.",
    narrative_en="All code uses the Go standard library with zero external dependencies. The test suite covers isolation, error propagation, and input validation.",
    mental_model="Use minimal reproduction commands to verify DI implementation locally.",
    source_path="README.md:14-23",
    slide_num="22 / 24",
    layout="code"
)

# Slide 23: Prinsip Arsitektural Jangka Panjang
add_slide(
    categories=["TAKEAWAY", "ARCHITECTURAL PRINCIPLES"],
    title="Prinsip Inti: Dependency Harus Eksplisit",
    subtitle="Konfigurasi terpisah dari penggunaan; interface sebagai kontrak.",
    content_items=[
        {"type": "cards", "title": "TIGA PILAR ARSITEKTUR DI", "items": [
            {"title": "Separation of Concerns", "desc": "Configuration (composition root) terpisah dari use (business logic)."},
            {"title": "Program to Interfaces", "desc": "Consumer bergantung pada kontrak interface, bukan implementasi konkret."},
            {"title": "Explicit over Implicit", "desc": "Semua dependensi dinyatakan eksplisit di constructor; tidak ada magic lookup."}
        ]}
    ],
    narrative="Prinsip DI menekankan kejelasan dependensi. Dengan memisahkan konfigurasi dan penggunaan, sistem menjadi lebih mudah diuji, dipahami, dan dikembangkan.",
    narrative_en="DI principles emphasize dependency clarity. By separating configuration from use, systems become easier to test, understand, and evolve.",
    mental_model="Architect systems with explicit dependencies and separated concerns for long-term maintainability.",
    source_path="content/05-key-takeaways.md",
    slide_num="23 / 24",
    layout="cards"
)

# Slide 24: Closing Slide
add_slide(
    categories=["CLOSING"],
    title="Dependency Injection adalah tentang kejelasan, bukan kerumitan.",
    subtitle="",
    content_items=[
        {"type": "bullet", "title": "MODUL IMPLEMENTASI", "items": [
            "internal/di/gateway.go • internal/di/processor.go • internal/di/locator.go",
            "cmd/demo/main.go • tests/processor_test.go",
            "content/02-master-draft.md • content/01-content-brief.md"
        ]},
        {"type": "bullet", "title": "VERIFIKASI & STATUS", "items": [
            f"Commit: {COMMIT_SHA[:7]} • Status: APPROVED",
            "Test passing: 6 test scenarios (success, error propagation, invalid amount)",
            "Race detector: Clean (go test -race ./... PASS)",
            "Zero external dependencies: 100% Go standard library."
        ]}
    ],
    mental_model="Dependency Injection separates concerns, making systems testable and evolvable.",
    source_path=REPO_URL,
    slide_num="24 / 24",
    layout="closing"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
    for s in SLIDES:
        print(f"[{s['slide_num']}] {s['layout']}: {s['title']}")