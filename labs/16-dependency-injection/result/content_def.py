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
COMMIT_SHA = "d7fb619"
LAB_NO = "16"
SLUG = "dependency-injection"
TITLE_ID = "Dependency Injection di Go"
SUBTITLE_ID = "Constructor Injection vs Service Locator: Trade-off Konkret untuk Testability"
TITLE_EN = "Dependency Injection in Go"
SUBTITLE_EN = "Constructor Injection vs Service Locator: A Concrete Trade-off for Testability"
AUTHOR = "Lukman / lukman-ss"
TOPIC_TAGS = "DI • IoC • Go • Constructor Injection • Service Locator • Testing"

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
            "Three module cards labeled 'Gateway', 'Payment', 'Notification' plugged into a central socket shape",
            "Asymmetric connection lines from modules into the socket, implying injection",
            f"Penulis: {AUTHOR}",
        ]}
    ],
    mental_model="",
    source_path="",
    slide_num="01 / 24"
)

# Slide 2: Engineering Problem — Tight Coupling
# Pre-render spec:
# - Runtime behavior: Processor constructs RealGateway internally via &RealGateway{}, no way to substitute without editing source.
# - Evidence: none (demo code uses &RealGateway{} directly; problem is structural)
# - Correct interpretation: Tight coupling means swap requires touching domain code; test requires network.
# - Do not overclaim: Lab does not measure defect rate reduction or production coupling metrics. Only demonstrates the structural contrast.
# Final headline: 8 kata, insight tajam
# Final body: maks 70 kata (body rule: maks 45 kata untuk non-code/table)
# Final footer: prinsip spesifik
add_slide(
    categories=["PROBLEM", "ENGINEERING CONFLICT"],
    title="Coupling Mengunci Gateway di Domain",
    subtitle="Processor yang membuat RealGateway sendiri tidak bisa diuji tanpa jaringan.",
    content_items=[
        {"title": "MASALAH STRUKTURAL", "items": [
            "Processor meng-instantiate &RealGateway{} secara internal.",
            "Swap provider gateway memaksa edit domain code.",
            "Test menjadi integration test lambat dan flaky."
        ], "type": "cards"}
    ],
    narrative="Dalam kode tradisional, domain langsung membangun dependensi konkret. Lab ini menunjukkan bahwa pola tersebut mengunci pengujian pada infrastruktur nyata — bukan berarti setiap sistem tanpa DI gagal production, melainkan trade-off antara kemudahan写 dan kemampuan uji.",
    mental_model="When a component builds its own dependencies, replacing them requires editing the component itself.",
    source_path="internal/di/processor.go • README.md",
    slide_num="02 / 24",
    layout="cards"
)

# Slide 3: Why Naive Approach Fails — The Lock-In Cost
# - Runtime: NewProcessor needs Concrete type, not interface. Swapping = grep+replace.
# - Evidence: main.go shows manual wiring; test shows only mock substitution possible with interface.
# - Correct: Interface parameter enables mock substitution without touching Processor source.
# - Do not overclaim: Lab does not prove reduced production bug count or faster TTM.
add_slide(
    categories=["WHY FAILS", "TRADE-OFF"],
    title="Konstruktor Tanpa Interface Menyulitkan Swap",
    subtitle="Menambahkan parameter interface hanya mengubah signature, bukan domain logic.",
    content_items=[
        {"title": "BEFORE (Tight Coupling)", "items": [
            "NewProcessor() — tidak ada parameter dependency.",
            "Setiap test harus menggunakan RealGateway asli.",
            "Ganti provider = ubah seluruh call site."
        ], "type": "bullet"},
        {"title": "AFTER (Constructor Injection)", "items": [
            "NewProcessor(g PaymentGateway) — menerima interface.",
            "Test cukup berikan MockGateway.",
            "Ganti provider = ubah satu baris di composition root."
        ], "type": "bullet"}
    ],
    narrative="Perubahan dari konstruktor tanpa parameter ke konstruktor dengan interface PaymentGateway tidak mengubah satu baris pun logika bisnis di dalam ProcessPayment. Ia hanya memindahkan titik pemilihan implementasi dari dalam domain ke composition root. Ini bukan jaminan test lebih cepat, tetapi ia membuat mock substitution mungkin tanpa refaktor lebih jauh.",
    mental_model="Injecting an interface moves the selection point from domain to the composition root.",
    source_path="internal/di/processor.go • tests/processor_test.go",
    slide_num="03 / 24",
    layout="comparison"
)

# Slide 4: Mental Model — Separation of Configuration from Use
# - Runtime: main.go builds RealGateway, passes it. Processor never sees construction.
# - Evidence: cmd/demo/main.go lines 17-21
# - Correct: Configuration phase (main) is separate from use phase (ProcessPayment).
# - Do not overclaim: Lab does not address lifecycle management or scope.
add_slide(
    categories=["MENTAL MODEL", "CORE CONCEPT"],
    title="IoC Bukan Fitur, Tapi Arah Kendali",
    subtitle="Domain menerima dependensi, bukan memilih atau membuatnya sendiri.",
    content_items=[
        {"title": "PRINSIP INTI", "items": [
            "Constructor injection memindahkan instantiation ke caller.",
            "Processor bergantung pada PaymentGateway interface, bukan implementasi.",
            "Komposisi terjadi sekali di composition root (main.go)."
        ], "type": "bullet"}
    ],
    narrative="Inversion of Control adalah prinsip umum; DI adalah bentuk spesifiknya untuk dependensi. Dalam lab ini, arahnya sederhana: main.go memegang kendali atas concrete type, Processor hanya menerima interface yang sudah jadi. Domain validation dan gateway error tetap berjalan di processor, tetapi pemilihan gateway bukan urusan domain lagi.",
    mental_model="IoC is a principle of control direction; DI is one specific mechanism to realize it for dependencies.",
    source_path="content/02-master-draft.md • cmd/demo/main.go",
    slide_num="04 / 24",
    layout="bullet"
)

# Slide 5: System / Data Flow — Request Lifecycle
# - Visual: flow diagram (boxes & arrows)
# - No code, no table; pure flow
# - Body maks 45 kata
# - Do not overclaim: Flow is simplified; no concurrency, no retries shown.
add_slide(
    categories=["FLOW DIAGRAM", "ARCHITECTURE"],
    title="Alur Request dari Composition Root ke Gateway",
    subtitle="Empat langkah: konfigurasi → injeksi → validasi → charge.",
    content_items=[
        {"title": "REQUEST FLOW", "items": [
            "main.go membuat RealGateway & menyuntikkannya ke NewProcessor.",
            "Call site memanggil ProcessPayment(100).",
            "Processor validasi amount > 0, lalu buat Money{100, USD}.",
            "Gateway.Charge(m) dipanggil; error dikembalikan ke caller."
        ], "type": "flow"}
    ],
    narrative="Flow ini menggambarkan eksekusi konkrit dari composition root hingga gateway dipanggil. Validasi invalid amount terjadi sebelum gateway terlibat, sehingga MockGateway.ChargedMoney tetap 0 — terbukti oleh TestProcessor_InvalidAmount. Simulasi ini tidak sama dengan flow production yang memiliki retry, timeout, dan circuit breaker.",
    mental_model="A clean request flow separates configuration, validation, and external call concerns.",
    source_path="internal/di/processor.go • cmd/demo/main.go",
    slide_num="05 / 24",
    layout="flow"
)

# Slide 6: Key Code Path — Constructor Injection
# - Code block: 8 lines max, visually dominant
# - Body maks 30 kata
# - Do not overclaim: nil check absent; nil would panic at runtime if used.
add_slide(
    categories=["CODE", "IMPLEMENTATION"],
    title="Constructor Membuat Dependency Terlihat",
    subtitle="Signature NewProcessor(g PaymentGateway) memindahkan keputusan ke caller.",
    content_items=[
        {"title": "KODE KRITIS", "items": [
            "func NewProcessor(g PaymentGateway) *Processor {",
            "    return &Processor{gateway: g}",
            "}",
            "",
            "func (p *Processor) ProcessPayment(amount int) error {",
            "    if amount <= 0 { return errors.New(\"invalid amount\") }",
            "    m := Money{Amount: amount, Currency: \"USD\"}",
            "    return p.gateway.Charge(m)",
            "}"
        ], "type": "code"}
    ],
    narrative="Parameter constructor g PaymentGateway membuat dependency terlihat di signature publik. Domain tidak lagi perlu import RealGateway. Ini membuat test bisa memberi mock, tetapi constructor masih perlu guard jika nil dependency ingin ditolak eksplisit — saat ini tidak ada nil check.",
    mental_model="Constructor parameters are the public contract for dependencies; missing nil guard means runtime risk remains.",
    source_path="internal/di/processor.go",
    slide_num="06 / 24",
    layout="code"
)

# Slide 7: Failure Mode — Service Locator Anti-Pattern
# - Visual: failure path showing hidden dependency through container
# - Body maks 45 kata
# - Do not overclaim: BadProcessor works functionally; the cost is testability/maintainability.
add_slide(
    categories=["FAILURE MODE", "ANTI-PATTERN"],
    title="Service Locator Menyembunyikan Dependency",
    subtitle="BadProcessor menerima Container, bukan PaymentGateway langsung.",
    content_items=[
        {"title": "ANTI-PATTERN", "items": [
            "NewBadProcessor(c Container) menerima abstraksi yang lebih tinggi.",
            "ProcessPayment memanggil container.GetPaymentGateway() secara internal.",
            "Dependensi nyata PaymentGateway tersembunyi dari signature constructor.",
            "Test harus mock Container, bukan gateway — tambahan indirection."
        ], "type": "bullet"}
    ],
    narrative="BadProcessor tidak salah secara fungsional: ia tetap memanggil Charge dan menghasilkan output yang sama. Masalahnya terletak pada keterbacaan dan testability — setiap developer baru harus menelusuri container untuk memahami dependensi sebenarnya. Ini bukan berarti pattern ini selalu rusak di semua konteks, tetapi ia menambah biaya kognitif yang tidak perlu.",
    mental_model="Hiding dependencies behind a locator increases cognitive load and makes mocking harder.",
    source_path="internal/di/locator.go",
    slide_num="07 / 24",
    layout="bullet"
)

# Slide 8: Hidden Cost — Indirection Tax
# - Table: cost comparison between Constructor Injection and Service Locator
# - Body maks 35 kata (table rule)
# - Do not overclaim: Costs are qualitative estimates from the lab's perspective, not benchmarked metrics.
add_slide(
    categories=["HIDDEN COST", "TRADE-OFF"],
    title="Biaya Tersembunyi: Indirection vs Eksplisitas",
    subtitle="Konstruktor injection menambah boilerplate kecil; Service Locator menambah ketidakjelasan.",
    content_items=[
        {"title": "PERBANDINGAN BIAYA", "items": [
            "Aspek | Constructor | Service Locator",
            "Testability | Mock langsung | Mock container + gateway",
            "Readability | Signature jelas | Butuh telusuri container",
            "Swap provider | Satu baris di root | Satu baris di root",
            "Risk nil | Panic jika tanpa guard | Panic jika Get returns nil"
        ], "type": "table"}
    ],
    narrative="Tabel ini membandingkan dua pola berdasarkan empat aspek yang relevan dengan perawatan kode. Constructor injection memiliki boilerplate lebih tinggi di composition root, tetapi membuat dependensi terlihat di setiap constructor call site. Service locator mengurangi jumlah parameter, tetapi menambatkan domain pada API container. Lab ini tidak mengukur overhead runtime kedua pendekatan.",
    mental_model="The hidden cost of a pattern is often measured in maintenance time, not execution time.",
    source_path="internal/di/processor.go • internal/di/locator.go",
    slide_num="08 / 24",
    layout="table"
)

# Slide 9: Experiment / Test Setup
# - Visual: test harness description with mock
# - Body maks 45 kata
# - Do not overclaim: Race detector only runs sequential tests; no concurrent test exists in lab.
add_slide(
    categories=["EXPERIMENT", "TEST SETUP"],
    title="Setup Pengujian: Mock Isolasi Gateway",
    subtitle="MockGateway merekam panggilan dan bisa diforced fail tanpa jaringan.",
    content_items=[
        {"title": "TEST HARNESS", "items": [
            "MockGateway menyimpan ChargedMoney dan toggle ShouldFail.",
            "6 test cases: 3 untuk Processor (success, error, invalid), 3 untuk BadProcessor.",
            "go test -race ./... dijalankan pada semua test sequential.",
            "Tidak ada test concurrent atau benchmark throughput."
        ], "type": "bullet"}
    ],
    narrative="Test suite dirancang untuk memverifikasi tiga hal: happy path, error propagation, dan validasi pre-check. MockGateway menggantikan RealGateway sepenuhnya — tidak ada panggilan HTTP, tidak ada dependency eksternal. Race detector berjalan pada test sequential; ia bukan bukti bahwa produksi tanpa concurrent access aman, hanya menunjukkan tidak ada data race pada skenario yang diuji.",
    mental_model="A focused test harness isolates one behavior at a time; it proves the path, not the whole system.",
    source_path="tests/processor_test.go",
    slide_num="09 / 24",
    layout="bullet"
)

# Slide 10: Result Snapshot — Test Output
# - Visual: large terminal output block
# - Body maks 30 kata
# - Do not overclaim: Output is from local run; numbers are not benchmarks.
add_slide(
    categories=["EVIDENCE", "RESULTS"],
    title="Semua 6 Test Pass, Race Detector Clean",
    subtitle="Output lokal: PASS 6/6, 0.00s per test case.",
    content_items=[
        {"title": "TERMINAL OUTPUT", "items": [
            "$ go test -race -v ./tests/",
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
            "PASS: ok lab16/tests (cached)"
        ], "type": "code"}
    ],
    narrative="Enam test case berlari dalam 0.00s per kasus. Race detector tidak melaporkan warning. Ini membuktikan jalur yang diuji aman dari data race pada goroutine tunggal — bukan berarti seluruh lifecycle production bebas race.",
    mental_model="Zero race warnings on a small sequential suite prove the tested paths; they do not prove the whole runtime safe.",
    source_path="tests/processor_test.go • README.md",
    slide_num="10 / 24",
    layout="code"
)

# Slide 11: Result Interpretation — What the Numbers Mean
# - Visual: metric cards or interpretation
# - Body maks 45 kata
# - Do not overclaim: Metrics are from local CI run, not production SLOs.
add_slide(
    categories=["INTERPRETATION", "EVIDENCE"],
    title="Angka 0.00s Tidak Otomatis Jadi SLO Production",
    subtitle="Waktu eksekusi test ≠ latency gateway nyata; mereka diukur di kondisi berbeda.",
    content_items=[
        {"title": "METRIK HASIL", "items": [
            "Metric | Nilai | Makna | Batas",
            "Test duration | 0.00s/case | Unit test in-memory | Bukan network latency",
            "Race warnings | 0 | Sequential test clean | Belum concurrent test",
            "Mock coverage | Charge captured | Path: valid+error+invalid | Belum edge: timeout",
            "Demo output | RealGateway charging | Integration manual verified | Bukan automated claim"
        ], "type": "table"}
    ],
    narrative="Tabel ini memetakan setiap angka ke interpretasi yang benar dan batasannya. Test duration 0.00s berarti operasi in-memory cepat, bukan gateway berlatensi nol. Race warnings nol berarti jalur sequential aman, bukan berarti konkurensi production aman. Cakupannya mengikuti test.",
    mental_model="Benchmark numbers are context-bound; they measure what the harness measures, nothing more.",
    source_path="tests/processor_test.go • cmd/demo/main.go",
    slide_num="11 / 24",
    layout="table"
)

# Slide 12: Trade-off — Boilerplate vs Testability
# - Visual: comparison two-column
# - Body maks 45 kata
# - Do not overclaim: Trade-off is qualitative; no quantitative measurement provided.
add_slide(
    categories=["TRADE-OFF", "ARCHITECTURE"],
    title="Trade-off: Boilerplate Konstruksi vs Kemampuan Uji",
    subtitle="Menambah satu parameter constructor membuka jalan untuk mock isolation.",
    content_items=[
        {"title": "BIAYA & MANFAAT", "items": [
            "Biaya: Tambah 1 parameter di constructor",
            "Manfaat: Mocking tanpa network call di unit test",
            "Biaya: Composition root perlu tahu concrete type",
            "Manfaat: Error propagation tetap terverifikasi rapi"
        ], "type": "bullet"}
    ],
    narrative="Trade-off ini sebanding: satu baris boilerplate di constructor mendapatkannya isolated testing. Namun mock isolation bukan berarti production safety — ia hanya menghilangkan dependency eksternal dari test. Validasi amount dan error wrapping masih membutuhkan test terpisah.",
    mental_model="DI trades a small construction cost for isolation gains; it does not buy production safety.",
    source_path="internal/di/processor.go • tests/processor_test.go",
    slide_num="12 / 24",
    layout="comparison"
)

# Slide 13: Improved Approach — Constructor Injection Pattern
# - Visual: code block (improved code) or flow
# - Body maks 30 kata
# - Do not overclaim: Still lacks nil guard; improved but not complete.
add_slide(
    categories=["IMPROVED APPROACH", "PATTERN"],
    title="Nil Guard Menutupi Runtime Risk",
    subtitle="Menambahkan pengecekan nil pada constructor membuat dependency wajib terlihat.",
    content_items=[
        {"title": "IMPROVED CODE", "items": [
            "func NewProcessor(g PaymentGateway) *Processor {",
            "    if g == nil {",
            "        return nil // atau panic, tergantung kebijakan tim",
            "    }",
            "    return &Processor{gateway: g}",
            "}"
        ], "type": "code"}
    ],
    narrative="Menambahkan nil check pada constructor merupakan langkah selanjutnya yang disarankan. Saat ini, constructor menerima nil tanpa peringatan — caller akan menghadapi nil pointer panic saat gateway.Charge dipanggil. Lab ini belum menerapkan guard tersebut; ia sengaja mempertahankan pola minimal untuk fokus pada kontras DI vs Service Locator.",
    mental_model="Explicit nil checks convert runtime surprises into compile-time or call-site failures.",
    source_path="internal/di/processor.go",
    slide_num="13 / 24",
    layout="code"
)

# Slide 14: Before vs After — Visual Comparison
# - Visual: before/after code or architecture diagram
# - Body maks 45 kata
# - Do not overclaim: Visual comparison is structural; performance not measured.
add_slide(
    categories=["COMPARISON", "VISUAL"],
    title="Sebelum vs Sesudah: Konstruktor Injection",
    subtitle="Perubahan hanya di signature; domain logic tetap identik.",
    content_items=[
        {"title": "BEFORE (Tight Coupling)", "items": [
            "Processor { gateway RealGateway }",
            "NewProcessor() → &Processor{gateway: &RealGateway{}}",
            "Test: harus pakai RealGateway atau setup mocking kompleks",
            "Swap: edit seluruh call site"
        ], "type": "bullet"},
        {"title": "AFTER (Constructor Injection)", "items": [
            "Processor { gateway PaymentGateway }",
            "NewProcessor(g PaymentGateway) → &Processor{gateway: g}",
            "Test: cukup passed MockGateway{}",
            "Swap: satu baris di main.go"
        ], "type": "bullet"}
    ],
    narrative="Perbandingan struktural ini menunjukkan bahwa perubahan terbesar bukan di dalam Processor, melainkan di titik komposisinya. Domain logic tetap sama; yang berubah adalah siapa yang bertanggung jawab atas concrete type. Ini bukan berarti semua sistem harus melalui transformasi serupa — ia hanya menunjukkan trade-off yang diambil lab ini.",
    mental_model="The biggest change is at the composition boundary, not inside the domain logic.",
    source_path="internal/di/processor.go • cmd/demo/main.go",
    slide_num="14 / 24",
    layout="comparison"
)

# Slide 15: Boundary of the Lab
# - Visual: boundary statement with clear demarcation
# - Body maks 45 kata
# - Do not overclaim: Explicitly state what lab does NOT cover.
add_slide(
    categories=["BOUNDARY", "LIMITATION"],
    title="Apa yang Tidak Dicakup oleh Lab Ini",
    subtitle="Laboratorium ini hanya membuktikan pola struktural, bukan skalabilitas production.",
    content_items=[
        {"title": "YANG DILAB", "items": [
            "Kontras struktural antara Constructor Injection dan Service Locator.",
            "Verifikasi mock substitution bekerja pada test sequential.",
            "Race detector bersih pada test suite saat ini."
        ], "type": "bullet"},
        {"title": "YANG TIDAK DILAB", "items": [
            "Lifecycle management (singleton, scoped, transient).",
            "Performance overhead reflection-based DI frameworks.",
            "Concurrent access patterns pada shared dependencies.",
            "Production defect rate reduction claim."
        ], "type": "bullet"}
    ],
    narrative="Lab ini sengaja membatasi diri pada pola struktural dasar. Ia tidak mengklaim bahwa constructor injection menghasilkan production-ready code tanpa pertimbangan lebih lanjut. Lifecycle management, concurrency safety, dan framework integration adalah topik yang masih perlu penyelidikan terpisah.",
    mental_model="A lab proves what it designs to prove; it does not absolve you from production considerations.",
    source_path="content/02-master-draft.md • README.md",
    slide_num="15 / 24",
    layout="comparison"
)

# Slide 16: Edge Case — Nil Dependency
# - Visual: edge case scenario
# - Body maks 45 kata
# - Do not overclaim: Nil panic is a runtime issue, not a DI flaw per se.
add_slide(
    categories=["EDGE CASE", "RISK"],
    title="Edge Case: Nil Dependency Menghasilkan Panic",
    subtitle="NewProcessor(nil) lulus konstrutor; panic terjadi saat Charge dipanggil.",
    content_items=[
        {"title": "Skenario Edge", "items": [
            "Caller mengoper nil sebagai PaymentGateway.",
            "Constructor tidak menolak; Processor tetap terbuat.",
            "ProcessPayment(100) → gateway.Charge(m) → nil pointer panic.",
            "Error tidak ter-catch; crash terjadi di runtime."
        ], "type": "bullet"}
    ],
    narrative="Nil dependency adalah edge case yang belum ditangani oleh constructor saat ini. Ia bukan kelemahan DI secara prinsip, melainkan konsekuensi dari ketiadaan explicit guard. Dalam production, caller yang bertanggung jawab memastikan non-nil; namun test suite saat ini tidak memiliki case untuk nil input.",
    mental_model="Absence of a guard is not a DI flaw; it is a defensive programming gap that any pattern must address.",
    source_path="internal/di/processor.go",
    slide_num="16 / 24",
    layout="bullet"
)

# Slide 17: What the Tests Prove
# - Visual: list of proven facts
# - Body maks 45 kata
# - Do not overclaim: Only the tested paths are proven; untested paths remain unproven.
add_slide(
    categories=["EVIDENCE", "PROVEN FACTS"],
    title="Test Suite Membuktikan Tiga Jalur Saja",
    subtitle="Happy path, gateway error, dan invalid amount — tidak lebih.",
    content_items=[
        {"title": "YANG DIBUKTIKAN", "items": [
            "MockGateway berhasil menangkap panggilan Charge dengan amount dan currency benar.",
            "Error dari gateway dipropagasikan ke caller tanpa di-suppress.",
            "Pre-check amount <= 0 mencegah panggilan gateway (ChargedMoney tetap 0)."
        ], "type": "bullet"}
    ],
    narrative="Ketiga jalur ini membentuk inti verifikasi behavior Processor. Namun test ini membuktikan jalur A, bukan klaim production B. Misalnya, concurrency safety, memory leak, atau integration dengan payment provider sungguhan tidak tercover oleh test suite yang ada.",
    mental_model="Tests prove the paths they cover; they do not prove the paths they do not traverse.",
    source_path="tests/processor_test.go",
    slide_num="17 / 24",
    layout="bullet"
)

# Slide 18: What the Tests Do Not Prove
# - Visual: list of unproven aspects
# - Body maks 45 kata
# - Do not overclaim: Explicitly name the gaps.
add_slide(
    categories=["EVIDENCE", "UNPROVEN"],
    title="Yang Belum Dibuktikan oleh Test Suite",
    subtitle="Test sequential tidak menjamin behavior pada konkurensi atau skala besar.",
    content_items=[
        {"title": "YANG TIDAK DIBUKTIKAN", "items": [
            "Concurrency safety: tidak ada concurrent test yang memicu race condition.",
            "Lifecycle management: tidak adauji singleton vs transient vs scoped.",
            "Framework DI overhead: tidak diukur latency injection vs manual new.",
            "Production defect reduction: tidak ada studi empiris dalam lab ini."
        ], "type": "bullet"}
    ],
    narrative="Ruang kosong ini bukan kegagalan lab — ia adalah batasan yang disengaja. Lab bertujuan menunjukkan kontras struktural, bukan menyelidiki aspek non-fungsional. Klaim production safety memerlukan evidence tambahan dari load test, chaos engineering, atau audit kode manual.",
    mental_model="The absence of evidence is not evidence of absence; it is a scope declaration.",
    source_path="tests/processor_test.go",
    slide_num="18 / 24",
    layout="bullet"
)

# Slide 19: Production Risk — Applying Pattern Raw
# - Visual: risk items
# - Body maks 45 kata
# - Do not overclaim: Risks are inferred from pattern characteristics, not measured in production.
add_slide(
    categories=["RISK", "PRODUCTION"],
    title="Risiko Menerapkan Pola Tanpa Pertimbangan Lifecycle",
    subtitle="Constructor injection tanpa lifecycle management berpotensi memory leak.",
    content_items=[
        {"title": "RISIKO POTENSIAL", "items": [
            "Setiap call site membuat Processor baru → object graph tumbuh tanpa batas.",
            "Shared state di dalam PaymentGateway (misal token cache) perlu sinkronisasi.",
            "Nil guard tidak ada → caller error menghasilkan panic di runtime.",
            "Komposisi root tunggal (`main.go`) bisa menjadi bottleneck pada sistem besar."
        ], "type": "bullet"}
    ],
    narrative="Risiko-rikiko ini muncul ketika pola diterapkan secara naif tanpa mempertimbangkan scope dan lifecycle. Lab ini tidak menyelesaikan masalah-masalah tersebut; ia hanya menunjukkan bahwa pola konstruktif injection lebih mudah diuji daripada tight coupling. Solusi lengkap memerlukan dependency container atau manual wiring yang lebih matang.",
    mental_model="Pattern application without lifecycle awareness turns a testing benefit into an operational risk.",
    source_path="content/02-master-draft.md • internal/di/processor.go",
    slide_num="19 / 24",
    layout="bullet"
)

# Slide 20: Implementation Checklist
# - Visual: checklist format
# - Body maks 45 kata
# - Do not overclaim: Checklist is heuristic, not enforced by tooling in lab.
add_slide(
    categories=["CHECKLIST", "IMPLEMENTATION"],
    title="Checklist Implementasi Dependency Injection yang Aman",
    subtitle="Langkah-langkah praktis untuk menerapkan DI tanpa overengineering.",
    content_items=[
        {"title": "LANGKAH KERJA", "items": [
            "[ ] Dependeansi dideklarasikan sebagai interface, bukan concrete type.",
            "[ ] Constructor injection digunakan; semua dependency eksplisit di signature.",
            "[ ] Composition root (`main.go`) adalah satu-satunya tempat konstruksi concrete.",
            "[ ] Value objects (seperti `Money`) diinstansiasi langsung, bukan melalui DI.",
            "[ ] Validasi input dilakukan sebelum delegasi ke dependency eksternal.",
            "[ ] `go test -race ./...` dijalankan dan PASS sebelum merge."
        ], "type": "checklist"}
    ],
    narrative="Checklist ini merupakan panduan minimum yang diturunkan dari eksperimen lab. Ia bukan standar universal — masing-masing tim mungkin menambahkan aturan tambahan seperti dependency audit, interface segregation, atau DI container integration. Poin kelima hingga keenam belum sepenuhnya diuji pada skenario concurrent production.",
    mental_model="Checklists guide practice; they do not replace architectural judgment or runtime verification.",
    source_path="content/02-master-draft.md • tests/processor_test.go",
    slide_num="20 / 24",
    layout="checklist"
)

# Slide 21: Anti-Pattern — Common Mistakes
# - Visual: anti-pattern list
# - Body maks 45 kata
# - Do not overclaim: Anti-patterns are observed from the lab's code structure.
add_slide(
    categories=["ANTI-PATTERN", "LESSON"],
    title="Anti-Pattern: Service Locator di Domain Logic",
    subtitle="Menyuntikkan container ke kelas bisnis menyembunyikan dependensi aktual.",
    content_items=[
        {"title": "HINDARI POLA INI", "items": [
            "NewBadProcessor(c Container) — signature menyembunyikan PaymentGateway.",
            "container.GetPaymentGateway() dipanggil di dalam method domain.",
            "Test harus mock Container alih-alih gateway — tambahan kompleksitas.",
            "Swap provider tetap satu baris, tetapi readability dan traceability menurun."
        ], "type": "bullet"}
    ],
    narrative="Service Locator bukanlah pola yang salah secara absolut — ia useful dalam konteks tertentu seperti plugin systems. Namun untuk domain logic yang memerlukan testability dan traceability, pola ini menambah biaya kognitif tanpa memberikan manfaat yang setara. Lab ini menunjukkannya sebagai kontras, bukan sebagai solusi yang dilarang mutlak.",
    mental_model="Service Locator is a tool, not a sin; misuse depends on where and why it is applied.",
    source_path="internal/di/locator.go",
    slide_num="21 / 24",
    layout="bullet"
)

# Slide 22: Reproduction Guide
# - Visual: command blocks
# - Body maks 45 kata
# - Do not overclaim: Commands reproduce the lab; they do not prove production readiness.
add_slide(
    categories=["REPRODUKSI", "COMMAND"],
    title="Reproduksi Mandiri di Mesin Lokal",
    subtitle="Go standard library, zero external dependencies.",
    content_items=[
        {"title": "COMMAND RUNNER", "items": [
            "cd labs/16-dependency-injection",
            "go test -race -v ./tests/     # validasi 6 test cases",
            "go run ./cmd/demo            # jalankan demo kedua pola",
            "go vet ./...                 # statik analisis tambahan"
        ], "type": "code"}
    ],
    narrative="Lab ini zero-dependency dan dapat direproduksi langsung di mesin lokal manapun. Menjalankan test dan demo membuktikan pola struktural bekerja sesuai deskripsi. Ia tidak membuktikan bahwa pola tersebut cocok untuk sistem dengan ratusan call site atau multi-threaded production workload.",
    mental_model="Reproducibility on a local machine is necessary but not sufficient for production confidence.",
    source_path="README.md • go.mod • tests/",
    slide_num="22 / 24",
    layout="code"
)

# Slide 23: Technical Takeaway
# - Visual: key insight
# - Body maks 45 kata
# - Do not overclaim: Takeaway is specific to this lab's scope.
add_slide(
    categories=["TAKEAWAY", "ENGINEERING PRINCIPLE"],
    title="Constructor Membuat Dependency Terlihat",
    subtitle="Signature yang eksplisit memudahkan mock, swap, dan code review.",
    content_items=[
        {"title": "PRINSIP REKAYASA", "items": [
            "DI bukan fitur framework — ia arah kendali atas instantiation.",
            "Constructor injection memindahkan keputusan ke composition root.",
            "Service Locator mengurangi boilerplate parameter, tetapi menambah indirection.",
            "Nil guard dan lifecycle management adalah tanggung jawab implementasi, bukan pola."
        ], "type": "bullet"}
    ],
    narrative="Pelajaran utama lab ini: dependency injection mengubah bagaimana kita berpikir tentang ownership. Domain tidak owning dependensinya; caller yang owning dan menyuntikkannya. Ini tidak otomatis membuat codebase production-ready, tetapi ia membuat trade-off antara testability dan complexity lebih terlihat.",
    mental_model="DI shifts ownership from the component to its assembler; visibility improves testability.",
    source_path="content/05-key-takeaways.md",
    slide_num="23 / 24",
    layout="bullet"
)

# Slide 24: Closing Slide (Dark Navy)
add_slide(
    categories=["CLOSING"],
    title="Constructor memperlihatkan dependency wajib; validasi nil tetap tanggung jawab constructor.",
    subtitle="",
    content_items=[
        {"title": "MODUL LAB", "items": [
            "internal/di/processor.go • internal/di/locator.go • internal/di/gateway.go",
            "cmd/demo/main.go • tests/processor_test.go",
            "README.md • content/02-master-draft.md"
        ], "type": "bullet"},
        {"title": "VERIFIKASI & STATUS", "items": [
            f"Sumber: commit {COMMIT_SHA[:8]}",
            "Test passing: 6 test cases, 0 failures.",
            "Race detector: clean (go test -race ./...)",
            "Simulasi in-memory mock; bukan proof production DI governance completeness."
        ], "type": "bullet"}
    ],
    mental_model="Constructor injection makes dependencies explicit; it does not replace lifecycle management or nil safety.",
    source_path=REPO_URL,
    slide_num="24 / 24"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
