#!/usr/bin/env python3
"""
Publication Builder Definition for Software Engineering Lab 16: Dependency Injection.
Follows:
- cover-rules.md (Mandatory raster hero image at COVER_IMAGE_PATH, series identity lock, no micro-text)
- slide-rules.md / slide-content-rules.md (<=45 words prose, <=30 words on code slides, dominant visual anchors, 50% zoom test)
- quality-bar.md (sharp engineering insights, strictly bounded claims, >=12 slides with limiting phrases, no overclaims)
- article-rules.md (Dev.to 22-section mapping, Medium ID HTML 22-section mapping)
"""
import os
import sys

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/16-dependency-injection/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/16-dependency-injection"
COMMIT_SHA = "5563838"
LAB_NO = "16"
SLUG = "dependency-injection"

TITLE_ID = "Dependency Injection & Inversion of Control"
SUBTITLE_ID = "Memisahkan Konfigurasi dari Penggunaan untuk Pengujian Terisolasi di Go"
TITLE_EN = "Dependency Injection & Inversion of Control"
SUBTITLE_EN = "Separating Configuration from Use for Isolated Testing in Go"

AUTHOR = "Lukman / lukman-ss"
TOPIC_TAGS = "Dependency Injection • Inversion of Control • Go • Testing • Software Architecture"
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

# Slide 01: Cover (Layout: cover)
add_slide(
    categories=["COVER"],
    title=TITLE_ID,
    subtitle=SUBTITLE_ID,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "Retro Terminal Displaying Decoupled Object Graph & Interface Boundaries",
            "Constructor Injection vs Service Locator Anti-Pattern in Pure Go",
            f"Penulis: {AUTHOR}"
        ]}
    ],
    mental_model="Separate object configuration from business logic to enable fast, isolated verification.",
    source_path="",
    slide_num="01 / 24",
    layout="cover"
)

# Slide 02: Problem Nyata (Layout: cards)
# Limiting phrase 1: "bukan berarti"
add_slide(
    categories=["PROBLEM", "TIGHT COUPLING"],
    title="Instansiasi Langsung di Logika Bisnis Mengunci Dependensi",
    subtitle="Ketika class memanggil new concrete service di dalam method, sistem terkunci pada infrastruktur.",
    content_items=[
        {"type": "cards", "title": "DAMPAK INSTANSIASI LANGSUNG DI LOGIKA BISNIS", "items": [
            {"title": "Test Lambat & Rapuh", "desc": "Pengujian terpaksa menyentuh jaringan atau database eksternal karena dependensi tidak dapat diganti."},
            {"title": "Vendor Lock-in Kaku", "desc": "Mengganti gateway pembayaran memaksa refaktorisasi menyeluruh di seluruh lapisan logika domain."},
            {"title": "Bukan Berarti Rusak", "desc": "Kode berjalan normal di production, tetapi bukan berarti arsitektur mudah berevolusi tanpa rewrite."}
        ]}
    ],
    narrative="Ketika kelas menginstansiasi dependensi secara langsung di dalam method, logika bisnis terkunci ke infrastruktur konkret. Kode mungkin bekerja tanpa error di tahap awal, bukan berarti arsitekturnya sehat ketika vendor pihak ketiga harus diganti.",
    narrative_en="Hardcoding concrete service instantiations inside business methods couples domain logic directly to external infrastructure. The code may run successfully initially; this does not mean the architecture can evolve cleanly when vendor integrations change.",
    mental_model="Direct object instantiation inside business logic locks domain code to external infrastructure.",
    source_path="content/02-master-draft.md:3-12",
    slide_num="02 / 24",
    layout="cards"
)

# Slide 03: Mengapa Review Manual / Pendekatan Naif Gagal (Layout: cards)
# Limiting phrase 2: "tidak otomatis"
add_slide(
    categories=["WHY FAILS", "SERVICE LOCATOR"],
    title="Service Locator Terlihat Bersih Namun Menyembunyikan Keterikatan",
    subtitle="Menyuntikkan generic container memindahkan instansiasi keluar tetapi mengaburkan kontrak dependensi.",
    content_items=[
        {"type": "cards", "title": "PENDEKATAN NAIF VS REALITA TEKNIS", "items": [
            {"title": "Kontrak Tersembunyi", "desc": "Signature NewBadProcessor(Container) menyembunyikan service apa saja yang sebenarnya dibutuhkan."},
            {"title": "Mocking Menjadi Rumit", "desc": "Unit test wajib mem-mock container dan method lookup alih-alih langsung menyuplai service mock."},
            {"title": "Tidak Otomatis Aman", "desc": "Menghilangkan kata kunci 'new' tidak otomatis membuat kode loosely coupled jika container bocor ke domain."}
        ]}
    ],
    narrative="Banyak tim mengira Service Locator menyelesaikan kopling karena objek tidak lagi memanggil new secara langsung. Faktanya, menyuntikkan container tidak otomatis menghasilkan kode decoupled jika dependensi sebenarnya tersembunyi di balik API lookup.",
    narrative_en="Teams often adopt Service Locators believing they eliminate coupling since direct instantiations disappear. In practice, injecting a container does not automatically decouple code when underlying dependencies remain hidden behind runtime lookups.",
    mental_model="Passing a container into business logic hides true dependencies behind runtime lookups.",
    source_path="content/02-master-draft.md:35-38",
    slide_num="03 / 24",
    layout="cards"
)

# Slide 04: Mental Model Inti — Schema Card (Layout: schema)
# Limiting phrase 3: "hanya pada jalur ini"
add_slide(
    categories=["MENTAL MODEL", "SEPARATION OF CONCERNS"],
    title="Pemisahan Konfigurasi dari Penggunaan: Dua Fase Terisolasi",
    subtitle="Konstruksi graf objek didelegasikan ke composition root; consumer hanya menerima kontrak.",
    content_items=[
        {"type": "schema", "title": "EMPAT PILAR SEPARATION OF CONFIGURATION FROM USE", "items": [
            {"phase": "1. Composition Root", "role": "Wiring Engine", "desc": "Hanya main.go atau bootstrap yang berhak membuat struct konkret RealGateway."},
            {"phase": "2. Interface Contract", "role": "Boundary Guard", "desc": "Processor hanya mengenal PaymentGateway interface tanpa tahu implementasi jaringan."},
            {"phase": "3. Constructor Injection", "role": "State Verifier", "desc": "NewProcessor(g) menjamin dependensi valid saat inisialisasi tanpa mutasi runtime."},
            {"phase": "4. Value Objects", "role": "Direct Instantiation", "desc": "Money dibuat langsung; hanya pada jalur ini objek tanpa IO eksternal tidak di-inject."}
        ]}
    ],
    narrative="Prinsip utama Fowler membedakan dua fase: perakitan graf objek di composition root dan eksekusi logika di domain. Hanya pada jalur ini komponen bisnis tetap bersih dari detail infrastruktur, sementara value object tetap diinstansiasi langsung.",
    narrative_en="Fowler's foundational principle establishes two separate phases: object graph assembly at the composition root and business execution within domain logic. Only along this path do business components remain completely insulated from infrastructure details.",
    mental_model="Separate configuration from use: assemble object graphs externally, pass dependencies via interfaces.",
    source_path="content/02-master-draft.md:13-34",
    slide_num="04 / 24",
    layout="schema"
)

# Slide 05: Diagram Alur / Pipeline Validasi (Layout: schema)
# Limiting phrase 4: "belum membuktikan"
add_slide(
    categories=["PIPELINE", "DEPENDENCY FLOW"],
    title="Alur Eksekusi: Composition Root ke Processor dan Gateway",
    subtitle="Dari wiring di main.go hingga verifikasi input dan delegasi transaksi payment.",
    content_items=[
        {"type": "schema", "title": "ALUR TRANSAKSI PEMBAYARAN TERISOLASI", "items": [
            {"phase": "1. Instantiation", "role": "main.go", "desc": "Menginstansiasi RealGateway{} dan menginjeksi ke NewProcessor(realGateway)."},
            {"phase": "2. Input Validation", "role": "Processor Guard", "desc": "Memeriksa amount <= 0; mengembalikan error instan sebelum menyentuh gateway."},
            {"phase": "3. Value Assembly", "role": "Domain Layer", "desc": "Membentuk struct Money{Amount, Currency} langsung di stack memory."},
            {"phase": "4. Gateway Charge", "role": "PaymentGateway", "desc": "Mendelegasikan Charge(m); belum membuktikan performa di jaringan latency tinggi."}
        ]}
    ],
    narrative="Alur eksekusi dimulai dari instansiasi konkret di composition root, dilanjutkan validasi input lokal sebelum memanggil antarmuka gateway. Arsitektur ini mengisolasi logika, namun belum membuktikan keandalan sistem pada jaringan berlatensi tinggi.",
    narrative_en="Execution initiates with concrete assembly at the composition root, followed by local input guards prior to delegating to the gateway interface. This architecture cleanly isolates logic, though it has not yet proven resiliency under high network latency.",
    mental_model="Validate inputs locally at domain boundaries before delegating to external interface contracts.",
    source_path="content/04-diagrams.md:79-101",
    slide_num="05 / 24",
    layout="schema"
)

# Slide 06: Cuplikan Kode Dominan (Layout: code)
# Limiting phrase 5: "bukan berarti"
add_slide(
    categories=["IMPLEMENTATION", "CONSTRUCTOR INJECTION"],
    title="Constructor Injection: Dependensi Eksplisit dan Valid Saat Lahir",
    subtitle="Menerima PaymentGateway via fungsi konstruktor NewProcessor tanpa dependensi tersembunyi.",
    content_items=[
        {"type": "code", "title": "internal/di/processor.go", "items": [
            "type Processor struct {",
            "    gateway PaymentGateway",
            "}",
            "func NewProcessor(g PaymentGateway) *Processor {",
            "    return &Processor{gateway: g}",
            "}",
            "func (p *Processor) ProcessPayment(amount int) error {",
            "    if amount <= 0 { return errors.New(\"invalid amount\") }",
            "    m := Money{Amount: amount, Currency: \"USD\"}",
            "    return p.gateway.Charge(m)",
            "}"
        ]}
    ],
    narrative="Processor mendeklarasikan gateway sebagai unexported field dan menerimanya saat konstruksi. Penggunaan interface bukan berarti arsitektur bebas alokasi memori; teknik ini menukar overhead pointer interface demi isolasi total dependensi domain.",
    narrative_en="Processor stores its gateway interface in an unexported field, accepting it strictly during initialization. Relying on interfaces does not imply zero-cost abstractions; it deliberately accepts interface pointer overhead to achieve complete domain isolation.",
    mental_model="Constructor injection guarantees immutable references and makes all dependencies explicit in the constructor signature.",
    source_path="internal/di/processor.go:6-21",
    slide_num="06 / 24",
    layout="code"
)

# Slide 07: Failure Path / Dangling Reference Konkret (Layout: cards)
# Limiting phrase 6: "hanya pada jalur ini"
add_slide(
    categories=["FAILURE PATH", "ANTI-PATTERN"],
    title="Kegagalan Service Locator: Keterikatan Tersembunyi Container",
    subtitle="Bagaimana NewBadProcessor memecahkan modularitas dan menyulitkan pelacakan dependensi.",
    content_items=[
        {"type": "cards", "title": "SKENARIO KEGAGALAN SERVICE LOCATOR", "items": [
            {"title": "Dependensi Gelap", "desc": "Membaca signature konstruktor tidak memberi tahu reviewer bahwa kelas membutuhkan PaymentGateway."},
            {"title": "Runtime Resolution Bug", "desc": "Jika container mengembalikan nil saat GetPaymentGateway(), terjadi nil-pointer panic saat runtime."},
            {"title": "Pelanggaran PSR-11", "desc": "Passing container ke dalam domain object melanggar panduan PSR-11: 'SHOULD NOT pass container'."},
            {"title": "Rigid Coupling", "desc": "Kelas terikat mati pada kontrak Container; hanya pada jalur ini pengujian terpaksa membuat mock container."}
        ]}
    ],
    narrative="Service Locator menyembunyikan dependensi nyata di dalam tubuh method dan berisiko memicu panic saat lookup gagal. Hanya pada jalur ini unit test dipaksa membuat MockContainer tambahan yang tidak relevan dengan kebutuhan bisnis.",
    narrative_en="The Service Locator masks actual dependencies inside method bodies, risking runtime panics when resolution fails. Only along this path are unit tests forced to construct irrelevant container mocks instead of directly supplying mock services.",
    mental_model="Hiding dependencies inside a container transforms compile-time missing dependency checks into runtime failures.",
    source_path="internal/di/locator.go:6-25",
    slide_num="07 / 24",
    layout="cards"
)

# Slide 08: Tabel Perbandingan Mekanisme vs Gejala (Layout: table)
add_slide(
    categories=["COMPARISON", "DIAGNOSTICS"],
    title="Matriks Komparasi: Constructor Injection vs Service Locator",
    subtitle="Membedah perbedaan mendasar transparansi kontrak, pengujian, dan siklus hidup komponen.",
    content_items=[
        {"title": "FITUR ARSITEKTUR VS POLA DESAIN", "items": [
            "Dimensi Evaluasi | Constructor Injection (Processor) | Service Locator (BadProcessor)",
            "Deklarasi Dependensi | Eksplisit di parameter NewProcessor | Tersembunyi di balik Container",
            "Type Safety & Validitas | Compile-time check, valid at birth | Runtime resolution, rawan nil panic",
            "Beban Mocking Pengujian | Minimal: langsung injeksi MockGateway | Ganda: wajib MockContainer + MockGateway",
            "Kopling Domain Logic | Nol ke framework / container | Terikat langsung ke interface Container",
            "Rekomendasi Standar | Best practice industri modern | Discouraged (PSR-11: SHOULD NOT)"
        ], "type": "table"}
    ],
    narrative="Matriks ini memetakan keunggulan mutlak Constructor Injection atas Service Locator. Injeksi via konstruktor menjamin verifikasi tipe pada saat kompilasi, sedangkan Service Locator menunda deteksi dependensi yang hilang ke runtime.",
    narrative_en="This matrix contrasts Constructor Injection against the Service Locator pattern. Constructor injection guarantees compile-time type safety and explicit dependencies, whereas service locators defer missing-dependency detection to runtime.",
    mental_model="Prefer compile-time safety and explicit dependencies over opaque runtime container lookups.",
    source_path="content/02-master-draft.md:35-38",
    slide_num="08 / 24",
    layout="table"
)

# Slide 09: Decision Sequence — 3 Focused Decision Cards (Layout: cards)
# Limiting phrase 7: "bukan berarti"
add_slide(
    categories=["DECISION SEQUENCE", "ARCHITECTURE"],
    title="Keputusan Arsitektur: Memilih Pola Perakitan Dependensi",
    subtitle="Tiga pertimbangan terstruktur dalam menentukan strategi injeksi dependensi di Go.",
    content_items=[
        {"type": "cards", "title": "SEKUENS KEPUTUSAN WIRING DEPENDENSI", "items": [
            {"title": "1. Constructor Default", "desc": "Gunakan Constructor Injection untuk semua service inti demi validitas status dan immutability."},
            {"title": "2. Bypass Value Objects", "desc": "Instansiasi langsung struct data murni seperti Money; jangan kelola DTO di dalam injector."},
            {"title": "3. Manual Composition", "desc": "Gunakan manual wiring di main.go; bukan berarti anti-framework, melainkan prioritas kesederhanaan YAGNI."}
        ]}
    ],
    narrative="Keputusan arsitektural menetapkan Constructor Injection sebagai standar baku dan menginstansiasi value object secara langsung. Memilih manual wiring di main.go bukan berarti anti-framework, melainkan keputusan pragmatis menjaga kesederhanaan aplikasi Go.",
    narrative_en="The architectural sequence establishes constructor injection as default while directly instantiating plain value objects. Choosing manual wiring in main.go does not reflect anti-framework dogma; it reflects a pragmatic commitment to Go simplicity.",
    mental_model="Default to constructor injection for services, instantiate value objects directly, and assemble graphs in main.",
    source_path="content/05-key-takeaways.md",
    slide_num="09 / 24",
    layout="cards"
)

# Slide 10: Bukti Eksekusi / Terminal Output (Layout: code)
add_slide(
    categories=["EVIDENCE", "EXECUTION"],
    title="Eksekusi Demo cmd/demo dan Test Suite Lolos Bersih",
    subtitle="Output terminal membuktikan kedua pola berjalan dan seluruh 6 unit test lulus tanpa warning.",
    content_items=[
        {"type": "code", "title": "TERMINAL OUTPUT: go test -v -race ./... && go run ./cmd/demo", "items": [
            "=== RUN   TestProcessor_Success",
            "--- PASS: TestProcessor_Success (0.00s)",
            "=== RUN   TestProcessor_GatewayError",
            "--- PASS: TestProcessor_GatewayError (0.00s)",
            "=== RUN   TestProcessor_InvalidAmount",
            "--- PASS: TestProcessor_InvalidAmount (0.00s)",
            "=== RUN   TestBadProcessor_Success / GatewayError / InvalidAmount",
            "--- PASS: TestBadProcessor_* (0.00s)",
            "PASS: ok lab16/tests 1.356s (race detector: clean)",
            "--- Running Constructor Injection: RealGateway charging 100 USD",
            "--- Running Service Locator: RealGateway charging 200 USD"
        ]}
    ],
    narrative="Verifikasi eksekusi membuktikan 6 skenario unit test lolos secara deterministik dalam hitungan detik. Eksekusi demo memperlihatkan kedua mekanisme berhasil memproses transaksi pembayaran tanpa race condition.",
    narrative_en="Verification demonstrates all 6 unit test scenarios passing deterministically in sub-second execution. The demo run confirms both injection mechanisms process payments successfully without triggering data races.",
    mental_model="Executable verification proves functional parity between patterns while exposing architectural divergence.",
    source_path="engineering/03-execution-result.md:8-46",
    slide_num="10 / 24",
    layout="code"
)

# Slide 11: Analisis Batasan Bukti (Layout: cards)
# Limiting phrase 8: "belum membuktikan"
add_slide(
    categories=["AUDIT", "EVIDENCE BOUNDARY"],
    title="Batasan Bukti Lab: Yang Dibuktikan vs Belum",
    subtitle="Menegaskan batasan verifikasi lab terhadap skenario enterprise berskala masif.",
    content_items=[
        {"type": "cards", "title": "BATAS PEMBUKTIAN SISTEM", "items": [
            {"title": "YANG DIBUKTIKAN", "desc": "Isolasi unit test via mock, validasi guard amount <= 0, dan propagasi gateway error."},
            {"title": "YANG TIDAK DIBUKTIKAN", "desc": "Lifecycle management singleton/scoped/transient dan performa container otomatis."},
            {"title": "BATAS KLAIM EMPIRIS", "desc": "Lab belum membuktikan klaim reduksi defect rate production atau batas baku 12 parameter."}
        ]}
    ],
    narrative="Lab ini membuktikan isolasi pengujian unit dan validasi input sebelum pemanggilan infrastruktur eksternal. Hasil ini belum membuktikan efektivitas lifecycle management kompleks atau klaim penurunan tingkat cacat di production.",
    narrative_en="This lab proves isolated unit testing and input guard execution prior to external infrastructure calls. It does not prove automated lifecycle resolution or empirical claims regarding defect rate reduction in production.",
    mental_model="Clearly delimit verified architectural properties from unproven enterprise heuristics.",
    source_path="content/01-content-brief.md:12-17",
    slide_num="11 / 24",
    layout="cards"
)

# Slide 12: Trade-off Arsitektur (Layout: cards)
# Limiting phrase 9: "tidak otomatis"
add_slide(
    categories=["TRADE-OFF", "COSTS & BENEFITS"],
    title="Trade-off Penerapan Dependency Injection",
    subtitle="Biaya indireksi dan konfigurasi tambahan berbanding lurus dengan fleksibilitas sistem.",
    content_items=[
        {"type": "cards", "title": "KOMPROMI REKAYASA DEPENDENCY INJECTION", "items": [
            {"title": "Indirection & Tracing", "desc": "Pemisahan konstruksi menyulitkan penelusuran implementasi konkret saat membaca kode tanpa IDE canggih."},
            {"title": "Boilerplate Inisialisasi", "desc": "Manual wiring menambah baris perakitan di composition root seiring membesarnya graph objek."},
            {"title": "Tidak Otomatis Sederhana", "desc": "Menerapkan DI tidak otomatis menyederhanakan kode jika kelas belum memiliki abstraksi yang stabil."}
        ]}
    ],
    narrative="Penerapan DI membawa konsekuensi penelusuran alur yang lebih rumit karena implementasi konkret terpisah dari titik penggunaan. Menambahkan interface dan injector tidak otomatis membuat arsitektur lebih baik jika sistemnya masih monolitik sederhana.",
    narrative_en="Adopting DI introduces cognitive indirection because concrete implementations are detached from invocation sites. Introducing interfaces and assemblers does not automatically improve architecture if the domain does not require implementation swapping.",
    mental_model="DI trades code colocation and simple tracing for loose coupling and isolated testability.",
    source_path="content/02-master-draft.md:126-134",
    slide_num="12 / 24",
    layout="cards"
)

# Slide 13: Concurrency / Validation Pipeline (Layout: code)
# Limiting phrase 10: "simulasi ini tidak sama dengan"
add_slide(
    categories=["CONCURRENCY", "RACE SAFETY"],
    title="Keamanan Konkurensi pada Komponen Stateless",
    subtitle="Processor yang hanya membaca interface gateway aman dieksekusi secara konkuren.",
    content_items=[
        {"type": "code", "title": "internal/di/processor.go", "items": [
            "// Processor stateless: gateway bersifat read-only setelah konstruksi",
            "type Processor struct {",
            "    gateway PaymentGateway",
            "}",
            "// Aman dipanggil lintas goroutine selama gateway thread-safe",
            "func (p *Processor) ProcessPayment(amount int) error {",
            "    if amount <= 0 { return errors.New(\"invalid amount\") }",
            "    m := Money{Amount: amount, Currency: \"USD\"}",
            "    return p.gateway.Charge(m)",
            "}"
        ]}
    ],
    narrative="Processor bersifat stateless setelah dikonstruksi sehingga aman dipanggil konkuren lintas goroutine tanpa mutasi internal. Simulasi ini tidak sama dengan gateway jaringan riil yang mengelola koneksi socket dan state pool secara terdistribusi.",
    narrative_en="The Processor struct remains stateless post-construction, rendering it inherently safe for concurrent execution across goroutines. While this simulation is not equivalent to distributed connection pools under network load, it passes race checks cleanly.",
    mental_model="Stateless service components receiving read-only interface references are inherently thread-safe.",
    source_path="tests/processor_test.go:23-59",
    slide_num="13 / 24",
    layout="code"
)

# Slide 14: Before vs After Comparison Cards (Layout: cards)
# Limiting phrase 11: "bukan berarti"
add_slide(
    categories=["COMPARISON", "BEFORE VS AFTER"],
    title="Transformasi Desain: Hardcoded vs Constructor Injection",
    subtitle="Pergeseran dari kopling ketat infrastruktur ke pemisahan tanggung jawab yang modular.",
    content_items=[
        {"type": "cards", "title": "SEBELUM (HARDCODED) VS SESUDAH (DI VIA CONSTRUCTOR)", "items": [
            {"title": "Sebelum: new RealGateway()", "desc": "Logika pembayaran menginstansiasi network client di dalam method; unit test mustahil tanpa internet."},
            {"title": "Sesudah: Constructor Injection", "desc": "Processor menerima PaymentGateway; unit test berjalan 0.00s menggunakan MockGateway."},
            {"title": "Bukan Berarti Rumit", "desc": "Refaktorisasi ke constructor injection hanya menambah 3 baris kode namun membuka kemampuan testability."}
        ]}
    ],
    narrative="Sebelum menerapkan DI, logika pembayaran terkunci pada implementasi RealGateway konkret yang menyulitkan pengujian offline. Memisahkan konstruksi ke NewProcessor bukan berarti menambah kompleksitas, melainkan membebaskan domain dari detail eksternal.",
    narrative_en="Prior to DI, business workflows directly instantiated concrete network gateways, making offline unit testing impossible. Refactoring to constructor injection does not mean adding unnecessary complexity; it frees domain models from external details.",
    mental_model="Extract concrete instantiation to constructor parameters to transform flaky integration tests into instant unit tests.",
    source_path="content/02-master-draft.md:3-12",
    slide_num="14 / 24",
    layout="cards"
)

# Slide 15: Analisis Mendalam — Kontrak Interface & Value Object (Layout: code)
# Limiting phrase 12: "cakupannya mengikuti test"
add_slide(
    categories=["ANALYSIS", "DOMAIN MODEL"],
    title="Interface Minimal dan Pengecualian untuk Value Object",
    subtitle="Mendefinisikan kontrak abstraksi ramping serta aturan kapan DI tidak boleh digunakan.",
    content_items=[
        {"type": "code", "title": "internal/di/gateway.go", "items": [
            "// Money adalah value object data murni: instansiasi langsung",
            "type Money struct {",
            "    Amount   int",
            "    Currency string",
            "}",
            "// PaymentGateway mendefinisikan kontrak minimal abstraction",
            "type PaymentGateway interface {",
            "    Charge(m Money) error",
            "}"
        ]}
    ],
    narrative="Interface PaymentGateway dirancang ramping hanya dengan satu method Charge, sedangkan Money diinstansiasi langsung tanpa perantara container. Keputusan membedakan service dan value object ini cakupannya mengikuti test skenario spesifik lab ini.",
    narrative_en="The PaymentGateway interface declares a single Charge method while Money is instantiated directly as a value object. Scoping DI strictly to services rather than value objects reflects design boundaries validated within this lab.",
    mental_model="Inject external infrastructure services; instantiate pure immutable value objects directly.",
    source_path="internal/di/gateway.go:5-15",
    slide_num="15 / 24",
    layout="code"
)

# Slide 16: Analisis Mendalam — Pengujian Terisolasi via Mock (Layout: code)
# Limiting phrase 13: "tidak otomatis"
add_slide(
    categories=["ANALYSIS", "ISOLATED TESTING"],
    title="Pengujian Terisolasi Menggunakan Test Double MockGateway",
    subtitle="Menguji logika Processor tanpa menyentuh koneksi jaringan riil dengan determinisme mutlak.",
    content_items=[
        {"type": "code", "title": "tests/processor_test.go", "items": [
            "type MockGateway struct {",
            "    ChargedMoney di.Money",
            "    ShouldFail   bool",
            "}",
            "func (m *MockGateway) Charge(money di.Money) error {",
            "    if m.ShouldFail { return errors.New(\"gateway unavailable\") }",
            "    m.ChargedMoney = money",
            "    return nil",
            "}"
        ]}
    ],
    narrative="MockGateway mensimulasikan kegagalan jaringan dan mencatat payload transaksi tanpa koneksi HTTP riil. Menulis unit test dengan mock tidak otomatis menjamin integritas integrasi vendor jika payload skema production mengalami perubahan.",
    narrative_en="MockGateway simulates network outages and records transaction payloads without touching HTTP sockets. Passing unit tests with mocks does not automatically guarantee third-party integration integrity if production API schemas change.",
    mental_model="Test doubles verify business logic deterministically without incurring network latency or side effects.",
    source_path="tests/processor_test.go:10-21",
    slide_num="16 / 24",
    layout="code"
)

# Slide 17: Analisis Mendalam — Validasi Guard Mencegah Panggilan Sia-sia (Layout: code)
# Limiting phrase 14: "hanya pada jalur ini"
add_slide(
    categories=["ANALYSIS", "INPUT VALIDATION"],
    title="Validasi Pre-Check Mencegah Panggilan Jaringan Sia-sia",
    subtitle="Verifikasi amount <= 0 memastikan gateway tidak pernah dipanggil pada payload invalid.",
    content_items=[
        {"type": "code", "title": "tests/processor_test.go", "items": [
            "func TestProcessor_InvalidAmount(t *testing.T) {",
            "    mock := &MockGateway{}",
            "    proc := di.NewProcessor(mock)",
            "    err := proc.ProcessPayment(-10)",
            "    if err == nil { t.Fatalf(\"expected error for negative amount\") }",
            "    if mock.ChargedMoney.Amount != 0 {",
            "        t.Errorf(\"gateway should not have been called, got %d\",",
            "            mock.ChargedMoney.Amount)",
            "    }",
            "}"
        ]}
    ],
    narrative="Test memastikan bahwa input negatif langsung ditolak di lapisan domain sebelum menyentuh gateway. Hanya pada jalur ini mock memverifikasi bahwa amount tetap 0, membuktikan panggilan jaringan berhasil dicegah sejak awal.",
    narrative_en="Tests confirm negative amounts are rejected at the domain boundary prior to contacting the gateway. Only along this pathway does the mock assert that Amount remains 0, confirming wasteful network invocations were prevented.",
    mental_model="Enforce business validation guards before delegating to downstream infrastructure dependencies.",
    source_path="tests/processor_test.go:47-59",
    slide_num="17 / 24",
    layout="code"
)

# Slide 18: Analisis Mendalam — Propagasi Error Tanpa Modifikasi (Layout: code)
# Limiting phrase 15: "bukan berarti"
add_slide(
    categories=["ANALYSIS", "ERROR PROPAGATION"],
    title="Propagasi Error Transparan dari Lapisan Infrastruktur",
    subtitle="Ketika gateway eksternal gagal, error diteruskan ke caller tanpa ditelan secara diam-diam.",
    content_items=[
        {"type": "code", "title": "tests/processor_test.go", "items": [
            "func TestProcessor_GatewayError(t *testing.T) {",
            "    mock := &MockGateway{ShouldFail: true}",
            "    proc := di.NewProcessor(mock)",
            "    err := proc.ProcessPayment(50)",
            "    if err == nil {",
            "        t.Fatalf(\"expected error from gateway, got nil\")",
            "    }",
            "    // Verifikasi error gateway sampai ke pemanggil",
            "}"
        ]}
    ],
    narrative="Processor meneruskan error dari gateway secara transparan ke caller. Meneruskan error mentah bukan berarti sistem tidak butuh retry policy; ini membuktikan kontrak interface mampu mengalirkan status kegagalan infrastruktur secara akurat.",
    narrative_en="The Processor bubbles gateway errors transparently back to the caller. Propagating raw errors does not imply the system lacks retry policies; it confirms interface boundaries accurately convey downstream failure states.",
    mental_model="Transparent error propagation preserves downstream failure context for upstream callers.",
    source_path="tests/processor_test.go:37-45",
    slide_num="18 / 24",
    layout="code"
)

# Slide 19: Analisis Mendalam — Wiring di Composition Root (Layout: code)
# Limiting phrase 16: "cakupannya mengikuti test"
add_slide(
    categories=["ANALYSIS", "COMPOSITION ROOT"],
    title="Composition Root di cmd/demo: Satu Titik Perakitan Objek",
    subtitle="Inisialisasi komponen konkret dan perakitan graf dilakukan di titik masuk aplikasi.",
    content_items=[
        {"type": "code", "title": "cmd/demo/main.go", "items": [
            "func main() {",
            "    // Instansiasi dependensi konkret",
            "    realGateway := &di.RealGateway{}",
            "    // Constructor Injection eksplisit",
            "    processor := di.NewProcessor(realGateway)",
            "    _ = processor.ProcessPayment(100)",
            "    // Service Locator anti-pattern untuk demonstrasi",
            "    container := &SimpleContainer{gateway: realGateway}",
            "    badProcessor := di.NewBadProcessor(container)",
            "    _ = badProcessor.ProcessPayment(200)",
            "}"
        ]}
    ],
    narrative="Fungsi main bertindak sebagai composition root tunggal tempat semua struct konkret dibuat dan dirakit. Struktur sederhana ini sengaja dipilih tanpa framework DI pihak ketiga, dan cakupannya mengikuti test verifikasi perilaku lab.",
    narrative_en="The main entrypoint functions as the sole composition root where concrete structs are instantiated and wired. This minimal structure avoids third-party DI reflection frameworks, with its scope strictly matching lab verification needs.",
    mental_model="Centralize concrete component instantiations strictly within the application composition root.",
    source_path="cmd/demo/main.go:16-36",
    slide_num="19 / 24",
    layout="code"
)

# Slide 20: Checklist Implementasi (Layout: checklist)
add_slide(
    categories=["CHECKLIST", "GUIDELINES"],
    title="Checklist Penerapan Dependency Injection di Go",
    subtitle="Pedoman wajib sebelum menggabungkan pull request yang mengelola dependensi infrastruktur.",
    content_items=[
        {"type": "checklist", "title": "KRITERIA IMPLEMENTASI DI", "items": [
            "Dependensi infrastruktur dideklarasikan sebagai interface, bukan struct konkret.",
            "Constructor injection digunakan sebagai default; semua dependensi eksplisit di parameter.",
            "Composition root (main.go atau wire bootstrap) adalah satu-satunya perakit implementasi konkret.",
            "Tidak ada pemanggilan generic container.Get() di dalam lapisan logika domain bisnis.",
            "Value object murni (seperti Money) diinstansiasi langsung tanpa melalui injection container.",
            "Seluruh service memiliki unit test terisolasi menggunakan mock tanpa menyentuh jaringan."
        ]}
    ],
    narrative="Checklist ini merangkum disiplin perancangan DI: kontrak berbasis interface, injeksi eksplisit via constructor, isolasi composition root, larangan container di domain, instansiasi langsung value object, dan pengujian unit terisolasi.",
    narrative_en="This checklist enforces robust DI practices: interface-driven contracts, explicit constructor injection, centralized composition root, container prohibition in domain logic, and isolated mock testing.",
    mental_model="Apply strict architectural gates to ensure dependencies remain explicit and testable.",
    source_path="content/02-master-draft.md:155-167",
    slide_num="20 / 24",
    layout="checklist"
)

# Slide 21: Anti-Pattern & Cara Reproduksi Kegagalan (Layout: cards)
# Limiting phrase 17: "bukan berarti"
add_slide(
    categories=["ANTI-PATTERN", "PITFALLS"],
    title="Anti-Pattern Umum dalam Penerapan Dependency Injection",
    subtitle="Kekeliruan perancangan yang sering merusak modularitas dan keterbacaan kode.",
    content_items=[
        {"type": "cards", "title": "HINDARI KESALAHAN UMUM INI", "items": [
            {"title": "Service Locator Leak", "desc": "Mengoper Container ke dalam domain; menyembunyikan dependensi nyata di balik lookup method runtime."},
            {"title": "Over-Injection & Bloat", "desc": "Konstruktor menerima belasan dependensi; sinyal keras pelanggaran Single Responsibility Principle."},
            {"title": "Value Object Injection", "desc": "Memasukkan struct data murni seperti Money atau DateTime ke container DI; bukan berarti lebih bersih."}
        ]}
    ],
    narrative="Menyuntikkan container ke domain menyembunyikan dependensi, sedangkan konstruktor dengan belasan parameter menandakan pelanggaran SRP. Memasukkan value object ke dalam container bukan berarti arsitektur lebih rapi, melainkan menambah boilerplate sia-sia.",
    narrative_en="Passing containers into domain logic conceals actual dependencies, while parameter-bloated constructors signal SRP violations. Managing simple value objects inside DI containers does not make code cleaner; it adds pointless ceremony.",
    mental_model="Avoid container leakage, prevent constructor parameter bloat, and never inject pure value objects.",
    source_path="content/02-master-draft.md:135-146",
    slide_num="21 / 24",
    layout="cards"
)

# Slide 22: Takeaway Praktis & Command Block Dominan (Layout: code)
# Limiting phrase 18: "production masih perlu"
add_slide(
    categories=["REPRODUCTION", "COMMANDS"],
    title="Reproduksi Mandiri di Mesin Lokal Menggunakan Go",
    subtitle="Eksperimen validasi mandiri tanpa dependensi pihak ketiga (zero-dependency).",
    content_items=[
        {"type": "code", "title": "COMMAND RUNNER (ZERO DEPENDENCIES)", "items": [
            "cd labs/16-dependency-injection",
            "go test -v ./...    # validasi unit test terisolasi Processor dan BadProcessor",
            "go test -race ./... # verifikasi konkurensi stateless bebas race condition",
            "go run ./cmd/demo   # eksekusi komparasi Constructor Injection vs Service Locator"
        ]}
    ],
    narrative="Seluruh mekanisme validasi ditulis menggunakan Go standard library murni tanpa framework eksternal. Skenario lokal ini memverifikasi isolasi logika, namun production masih perlu mempertimbangkan auto-wiring tool seperti wire pada skala ribuan komponen.",
    narrative_en="The entire test harness relies solely on the pure Go standard library with zero third-party dependencies. While this validates architectural isolation locally, production systems may still require code-generation tools like wire at enterprise scale.",
    mental_model="Zero-dependency local verification guarantees fast, repeatable validation across developer workflows.",
    source_path="README.md:13-23",
    slide_num="22 / 24",
    layout="code"
)

# Slide 23: Prinsip Arsitektural Jangka Panjang (Layout: cards)
# Limiting phrase 19: "tidak otomatis"
add_slide(
    categories=["TAKEAWAY", "ARCHITECTURAL PRINCIPLES"],
    title="Prinsip Inti: Pisahkan Konfigurasi, Eksplisitkan Kontrak",
    subtitle="Arsitektur bertahan lama ketika logika bisnis terbebas dari mekanisme perakitan infrastruktur.",
    content_items=[
        {"type": "cards", "title": "TIGA PILAR KELANGSUNGAN ARSITEKTUR", "items": [
            {"title": "Pisahkan Konfigurasi", "desc": "Biarkan composition root merakit graf objek; modul bisnis fokus murni pada aturan domain."},
            {"title": "Kontrak Eksplisit", "desc": "Gunakan Constructor Injection agar dependensi terlihat transparan sejak kompilasi."},
            {"title": "Disiplin Batas Domain", "desc": "Menerapkan DI tidak otomatis menyelesaikan semua masalah jika batasan bounded context tidak dirancang benar."}
        ]}
    ],
    narrative="Dependency Injection adalah teknik memisahkan fase konfigurasi dari eksekusi logika. Menyuntikkan dependensi via interface tidak otomatis menyelesaikan kekacauan arsitektur jika tanggung jawab domain sejak awal tidak dirancang dengan batas yang tegas.",
    narrative_en="Dependency Injection fundamentally separates configuration concerns from operational business logic. Leveraging interfaces and constructors does not automatically fix architectural debt if domain boundaries are poorly defined initially.",
    mental_model="Decouple configuration from execution to preserve long-term evolvability and testability.",
    source_path="content/05-key-takeaways.md",
    slide_num="23 / 24",
    layout="cards"
)

# Slide 24: Closing Slide (Layout: closing)
add_slide(
    categories=["CLOSING"],
    title="Logika bisnis yang menginstansiasi dependensinya sendiri adalah sistem yang menolak untuk diuji.",
    subtitle="",
    content_items=[
        {"type": "bullet", "title": "MODUL IMPLEMENTASI", "items": [
            "internal/di/gateway.go • internal/di/processor.go • internal/di/locator.go",
            "cmd/demo/main.go • tests/processor_test.go",
            "engineering/01-design.md • content/02-master-draft.md"
        ]},
        {"type": "bullet", "title": "VERIFIKASI & STATUS", "items": [
            f"Commit: {COMMIT_SHA[:7]} • Status: APPROVED",
            "Test passing: 6 unit tests (happy path, gateway error, input validation pre-check)",
            "Race detector: Clean (go test -race ./... PASS in 1.35s)",
            "Manual wiring Go murni; bukan framework DI berbasis refleksi runtime."
        ]}
    ],
    mental_model="Business logic that instantiates its own dependencies creates code that resists testing.",
    source_path=REPO_URL,
    slide_num="24 / 24",
    layout="closing"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
    for s in SLIDES:
        print(f"[{s['slide_num']}] {s['layout']}: {s['title']}")
