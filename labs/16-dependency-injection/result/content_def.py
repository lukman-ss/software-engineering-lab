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
COMMIT_SHA = "8678a90"
LAB_NO = "16"
SLUG = "dependency-injection"
TITLE_ID = "Dependency Injection"
SUBTITLE_ID = "Memisahkan Konfigurasi dari Penggunaan untuk Testability & Evolvability"
TITLE_EN = "Dependency Injection"
SUBTITLE_EN = "Separating Configuration from Use for Testability and Evolvability"
AUTHOR = "Lukman / lukman-ss"
TOPIC_TAGS = "Dependency Injection • Inversion of Control • Go • Test Doubles • Software Architecture"

SLIDES = []

def add_slide(categories, title, subtitle, content_items, mental_model, source_path="", slide_num="", narrative="", narrative_en="", layout=""):
    SLIDES.append({
        "categories": categories,
        "title": title,
        "subtitle": subtitle,
        "content_items": content_items,
        "mental_model": mental_model,
        "source": source_path,
        "source_path": source_path,
        "slide_num": slide_num,
        "narrative": narrative,
        "narrative_en": narrative_en,
        "layout": layout
    })

# =====================================================================
# SLIDE 1: Cover Slide (Dark Navy)
# =====================================================================
add_slide(
    categories=["COVER"],
    title=TITLE_ID,
    subtitle=SUBTITLE_ID,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "Dependency Graph: Composition Root (main) -> Processor -> PaymentGateway Interface",
            "Contrast: Constructor Injection (Explicit) vs Service Locator (Hidden)",
            f"Penulis: {AUTHOR}",
        ]}
    ],
    mental_model="Separation of configuration from use: outside world assembles, inner logic executes.",
    source_path="",
    slide_num="01 / 24",
    layout="cover"
)

# =====================================================================
# SLIDE 2: Problem Nyata (3 Compact Problem Cards)
# =====================================================================
add_slide(
    categories=["PROBLEM", "ENGINEERING CONFLICT"],
    title="Infrastruktur Terkunci di Logika Bisnis",
    subtitle="Ketika modul bisnis membuat dependensinya sendiri, isolasi dan pengujian runtuh.",
    content_items=[
        {"type": "problem_cards", "title": "MASALAH PRODUKSI", "items": [
            {
                "title": "Tightly Coupled Code",
                "desc": "Processor menginstansiasi RealGateway langsung di method; logika bisnis terkunci ke implementasi vendor tunggal."
            },
            {
                "title": "Pengujian Lambat & Flaky",
                "desc": "Unit test terpaksa membuka koneksi jaringan nyata atau memerlukan mock REST API server yang lambat dan rapuh."
            },
            {
                "title": "Ripple Effect Perubahan",
                "desc": "Mengganti vendor payment gateway memaksa refactoring dan rekompilasi menyeluruh di seluruh lapisan consumer."
            }
        ]}
    ],
    narrative="Kopling keras antara logika domain dan implementasi infrastruktur merusak evolvability sistem. Pengujian unit menjadi mustahil dilakukan secara murni tanpa network overhead.",
    mental_model="Hard-coded dependencies fuse domain logic to external plumbing, destroying testability.",
    source_path="content/02-master-draft.md",
    slide_num="02 / 24",
    layout="problem_cards"
)

# =====================================================================
# SLIDE 3: Mengapa Review Manual / Pendekatan Naif Gagal
# =====================================================================
add_slide(
    categories=["WHY FAILS", "HUMAN REVIEW"],
    title="Review Manual Melewatkan Hidden Dependencies",
    subtitle="Service Locator tampak bersih di permukaan namun menyembunyikan kontrak asli.",
    content_items=[
        {"type": "cards", "title": "KELEMAHAN PENDEKATAN NAIF", "items": [
            {
                "title": "Kontrak Opaque",
                "desc": "Signature NewBadProcessor(c Container) tidak menyatakan butuh PaymentGateway; dependensi baru terlihat di dalam baris method."
            },
            {
                "title": "Mocking Menjadi Rumit",
                "desc": "Test tidak bisa menyuntikkan MockGateway langsung, melainkan harus membuat MockContainer yang membungkus MockGateway."
            },
            {
                "title": "Ledakan Runtime",
                "desc": "Jika container lupa meregistrasikan gateway, compiler tetap lolos namun aplikasi meledak panic saat dipanggil di production."
            }
        ]}
    ],
    narrative="Reviewer manusia melihat parameter Container dan menganggap modularitas tercapai. Faktanya, Service Locator hanya memindahkan kopling dari compile-time ke runtime.",
    mental_model="Service Locator hides dependencies inside methods; constructor injection makes them explicit at creation.",
    source_path="internal/di/locator.go",
    slide_num="03 / 24",
    layout="cards"
)

# =====================================================================
# SLIDE 4: Mental Model Inti / Schema Card
# =====================================================================
add_slide(
    categories=["MENTAL MODEL", "CORE CONCEPT"],
    title="Separation of Configuration from Use",
    subtitle="Fowler: fase perakitan objek terpisah total dari fase eksekusi logika bisnis.",
    content_items=[
        {"type": "schema", "title": "DUA FASE TERISOLASI", "items": [
            {
                "phase": "Fase Konfigurasi (Composition Root)",
                "role": "main.go / Assembler",
                "desc": "Bertanggung jawab membuat objek konkret (RealGateway), merakit relasi, dan menginjeksikan ke consumer."
            },
            {
                "phase": "Fase Penggunaan (Domain Execution)",
                "role": "Processor (Business Service)",
                "desc": "Hanya menerima PaymentGateway via interface; mengeksekusi ProcessPayment tanpa tahu kelas konkret di baliknya."
            }
        ]}
    ],
    narrative="Prinsip paling fundamental dari DI bukanlah framework atau container, melainkan pemisahan waktu pembuatan objek dari waktu penggunaan objek.",
    mental_model="Assemble the graph once at startup; execute pure business logic during runtime.",
    source_path="content/02-master-draft.md",
    slide_num="04 / 24",
    layout="schema"
)

# =====================================================================
# SLIDE 5: Diagram Alur / Pipeline
# =====================================================================
add_slide(
    categories=["ARSITEKTUR", "DATA FLOW"],
    title="Pipeline Alur Eksekusi Dependency Injection",
    subtitle="Dari entry point aplikasi hingga eksekusi infrastruktur terabstraksi.",
    content_items=[
        {"type": "flow", "items": [
            "cmd/demo/main.go | Composition Root inisialisasi RealGateway",
            "di.NewProcessor(gateway) | Konstruktor menerima kontrak PaymentGateway",
            "processor.ProcessPayment(100) | Validasi input amount <= 0 & buat Money",
            "gateway.Charge(Money) | Delegasi eksekusi ke implementasi konkret"
        ]}
    ],
    narrative="Pipeline memperlihatkan batasan tanggung jawab yang jelas: composition root menangani konstruksi, processor menangani validasi, dan gateway menangani charge.",
    mental_model="Separate configuration from use: simple linear stages create decoupled architecture.",
    source_path="cmd/demo/main.go • internal/di/processor.go",
    slide_num="05 / 24",
    layout="flow"
)

# =====================================================================
# SLIDE 6: Cuplikan Kode Dominan (Maks 8 Baris)
# =====================================================================
add_slide(
    categories=["CODE", "MODEL"],
    title="Constructor Injection: Valid State at Birth",
    subtitle="Field unexported, dependensi eksplisit, dan objek dijamin lengkap.",
    content_items=[
        {"title": "INTERNAL/DI/PROCESSOR.GO", "items": [
            "type Processor struct {",
            "    gateway PaymentGateway",
            "}",
            "",
            "func NewProcessor(g PaymentGateway) *Processor {",
            "    return &Processor{gateway: g}",
            "}"
        ], "type": "code"}
    ],
    narrative="Dengan Constructor Injection, Processor tidak mungkin dibuat dalam keadaan setengah matang. Signature constructor menyatakan secara gamblang apa yang dibutuhkan.",
    mental_model="Explicit constructor parameters enforce immutable, valid instances at creation time.",
    source_path="internal/di/processor.go:6-13",
    slide_num="06 / 24",
    layout="code"
)

# =====================================================================
# SLIDE 7: Failure Path / Dangling Reference Konkret
# =====================================================================
add_slide(
    categories=["FAILURE PATH", "EARLY VALIDATION"],
    title="Failure Path: Pencegahan Eksekusi I/O Sia-sia",
    subtitle="Validasi domain memotong eksekusi sebelum memanggil payment gateway.",
    content_items=[
        {"type": "flow", "items": [
            "Input Transaksi | amount = -10 (invalid input)",
            "Processor Pre-check | if amount <= 0 { return error }",
            "Short-Circuit Exit | Return errors.New(\"invalid amount\")",
            "Infrastruktur Aman | gateway.Charge() tidak pernah dipanggil (Amount = 0)"
        ]}
    ],
    narrative="DI tidak menghilangkan kebutuhan validasi. Pre-check di Processor memastikan gateway tidak pernah dipanggil ketika parameter bisnis tidak valid.",
    mental_model="Validate at domain boundary; never leak invalid data into infrastructure calls.",
    source_path="internal/di/processor.go:15-21 • tests/processor_test.go:47-59",
    slide_num="07 / 24",
    layout="flow"
)

# =====================================================================
# SLIDE 8: Tabel Perbandingan Mekanisme vs Gejala
# =====================================================================
add_slide(
    categories=["TABEL", "PERBANDINGAN"],
    title="Tabel Evaluasi Mekanisme Dependensi",
    subtitle="Perbandingan struktural antara Constructor Injection, Service Locator, dan Hard-Coded.",
    content_items=[
        {"title": "MATRIKS ARSITEKTUR DEPENDENSI", "items": [
            "Kriteria | Hard-Coded New | Service Locator | Constructor Injection",
            "Transparansi Kontrak | Tertutup di method | Tertutup di container | Terbuka di constructor",
            "Validasi Compiler | Kuat (tipe konkret) | Lemah (runtime get) | Sangat Kuat (interface)",
            "Kompleksitas Test | Sangat tinggi (mock I/O) | Sedang (mock container) | Minimal (mock interface)",
            "Standar Industri | Anti-pattern monolit | PSR-11: SHOULD NOT | Rekomendasi Fowler / Default"
        ], "type": "table"}
    ],
    narrative="Constructor injection memberikan kombinasi terbaik: dependensi terlihat transparan, compiler memvalidasi pemenuhan interface, dan test double dapat diinjeksi tanpa overhead.",
    mental_model="Constructor injection provides compiler-enforced transparency; service locator delays failures to runtime.",
    source_path="content/02-master-draft.md",
    slide_num="08 / 24",
    layout="table"
)

# =====================================================================
# SLIDE 9: Decision Sequence (3 Focused Decision Cards)
# =====================================================================
add_slide(
    categories=["DECISION", "SEQUENCE"],
    title="Tiga Keputusan Desain Dependensi",
    subtitle="Urutan keputusan arsitektur untuk menjaga kemurnian boundary sistem.",
    content_items=[
        {"type": "decision_cards", "title": "DECISION SEQUENCE", "items": [
            {
                "step": "01. Kontrak Interface",
                "title": "Program to Interface",
                "desc": "Definisikan PaymentGateway Charge(Money) error untuk mengabstraksi seluruh varian vendor pembayaran."
            },
            {
                "step": "02. Wajibkan Ctor Injection",
                "title": "Explicit Dependencies",
                "desc": "Tolak Service Locator; oper PaymentGateway langsung ke konstruktor NewProcessor."
            },
            {
                "step": "03. Value Object Direct",
                "title": "Exempt Pure Data",
                "desc": "Instansiasi Money{Amount, Currency} langsung tanpa DI container karena tidak memiliki efek I/O."
            }
        ]}
    ],
    narrative="Keputusan ini memastikan bahwa abstraksi hanya diterapkan pada batas sistem eksternal, sementara struktur data murni tetap ringan tanpa beban overhead DI.",
    mental_model="Interface external boundaries, inject at constructors, instantiate value objects directly.",
    source_path="internal/di/gateway.go • internal/di/processor.go",
    slide_num="09 / 24",
    layout="decision_cards"
)

# =====================================================================
# SLIDE 10: Bukti Eksekusi / Terminal Output (1 Kalimat Penjelasan)
# =====================================================================
add_slide(
    categories=["RESULT", "EVIDENCE"],
    title="Output Demo: Eksekusi Dua Pendekatan",
    subtitle="Hasil run membuktikan kesamaan fungsional dengan perbedaan struktural total.",
    content_items=[
        {"title": "TERMINAL DEMO OUTPUT", "items": [
            "$ go run ./cmd/demo",
            "--- Running Constructor Injection ---",
            "RealGateway charging 100 USD",
            "--- Running Service Locator ---",
            "RealGateway charging 200 USD"
        ], "type": "code"}
    ],
    narrative="Demo membuktikan bahwa kedua pola berhasil mendistribusikan pembayaran ke RealGateway, namun Constructor Injection mencapai hasil ini tanpa keterikatan pada API container.",
    mental_model="Functional equivalence in execution; architectural superiority in isolation.",
    source_path="cmd/demo/main.go",
    slide_num="10 / 24",
    layout="code"
)

# =====================================================================
# SLIDE 11: Analisis Batasan Bukti (Heading Kanan: YANG TIDAK DIBUKTIKAN)
# =====================================================================
add_slide(
    categories=["BATASAN", "BOUNDARY"],
    title="Batasan Bukti: Cakupan Verifikasi Lab",
    subtitle="Memisahkan fakta terverifikasi laboratorium dari asumsi ekosistem produksi.",
    content_items=[
        {"title": "YANG TERBUKTI DI LAB", "items": [
            "MockGateway mengisolasi Processor 100% tanpa jaringan.",
            "Validasi pre-check mencegah pemanggilan gateway.",
            "Error gateway dipropagasikan utuh ke caller.",
            "Race detector bersih tanpa data race."
        ], "type": "bullet"},
        {"title": "YANG TIDAK DIBUKTIKAN", "items": [
            "Overhead performa container vs direct instantiation.",
            "Korelasi penurunan angka bug/defect di tim engineering.",
            "Lifecycle manajemen (Singleton, Scoped, Transient).",
            "Efektivitas DI framework otomatis (Wire, Dig) skala besar."
        ], "type": "bullet"}
    ],
    narrative="Lab membuktikan keunggulan struktural dan kemudahan pengujian unit. Lab tidak mengukur klaim performa runtime microbenchmark atau lifecycle container tingkat lanjut.",
    mental_model="Tests verify isolation and contract adherence, not ecosystem-scale lifecycle management.",
    source_path="content/01-content-brief.md • engineering-audit/05-gaps.md",
    slide_num="11 / 24",
    layout="boundary"
)

# =====================================================================
# SLIDE 12: Trade-off Arsitektur
# =====================================================================
add_slide(
    categories=["TRADE-OFF", "DECISION"],
    title="Trade-off: Kontrol Bersih vs Biaya Boilerplate",
    subtitle="Keuntungan modularitas harus dibayar dengan disiplin perakitan dependensi.",
    content_items=[
        {"title": "KEUNTUNGAN DIPEROLEH", "items": [
            "Isolasi total: unit test berjalan cepat dalam orde mikrodetik.",
            "Vendor agnostic: ganti gateway tanpa sentuh logika bisnis.",
            "Immutabilitas: dependensi tidak dapat diganti sembarangan di runtime."
        ], "type": "bullet"},
        {"title": "KONSEKUENSI DIHADAPI", "items": [
            "Manual wiring menambah baris inisialisasi di composition root.",
            "Tracing alur kode memerlukan navigasi antarmuka interface.",
            "Risiko constructor bloat jika kelas melanggar Single Responsibility."
        ], "type": "bullet"}
    ],
    narrative="DI bukan solusi tanpa biaya. Keuntungan testability dan fleksibilitas harus diimbangi dengan kehati-hatian agar constructor tidak membengkak karena service terlalu besar.",
    mental_model="Every decoupled boundary trades runtime simplicity for construction ceremony.",
    source_path="content/02-master-draft.md",
    slide_num="12 / 24",
    layout="trade_off"
)

# =====================================================================
# SLIDE 13: Concurrency / Validation Pipeline
# =====================================================================
add_slide(
    categories=["CONCURRENCY", "PIPELINE"],
    title="Pipeline Validasi & Thread-Safety",
    subtitle="Processor berstatus stateless dan aman dieksekusi lintas goroutine paralel.",
    content_items=[
        {"type": "flow", "items": [
            "Goroutine Request | Panggilan concurrent ke ProcessPayment(amount)",
            "Pure Pre-Validation | amount <= 0 dicek tanpa mutasi shared state",
            "Value Object Build | Money{Amount, Currency} dialokasikan per request",
            "Thread-Safe Charge | Gateway dieksekusi secara aman (Race Clean)"
        ]}
    ],
    narrative="Karena Processor tidak menyimpan state yang dapat dimutasi (hanya menyimpan reference interface gateway yang stateless), pemrosesan pembayaran bebas dari race condition.",
    mental_model="Immutable dependency references make services inherently safe for concurrent execution.",
    source_path="internal/di/processor.go • tests/processor_test.go",
    slide_num="13 / 24",
    layout="flow"
)

# =====================================================================
# SLIDE 14: Before vs After Comparison Cards
# =====================================================================
add_slide(
    categories=["KOMPARASI", "BEFORE AFTER"],
    title="Sebelum vs Sesudah Penerapan DI",
    subtitle="Transformasi kode rapuh menjadi komponen yang teruji dan adaptif.",
    content_items=[
        {"title": "SEBELUM: TIGHTLY COUPLED", "items": [
            "Processor memanggil new RealGateway() secara internal.",
            "Unit test terpaksa menghubungi server gateway nyata.",
            "Pergantian vendor menuntut perubahan di logika bisnis.",
            "Sulit menguji skenario kegagalan jaringan secara deterministik."
        ], "type": "before_after"},
        {"title": "SESUDAH: CONSTRUCTOR INJECTION", "items": [
            "Processor menerima PaymentGateway via NewProcessor.",
            "Unit test menyuntikkan MockGateway tanpa sentuh jaringan.",
            "Pergantian vendor hanya mengubah satu baris di main.go.",
            "Skenario kegagalan gateway dapat diuji deterministik 100%."
        ], "type": "before_after"}
    ],
    narrative="Sebelumnya, setiap pengujian bergantung pada kestabilan koneksi internet ke vendor. Sesudah DI, seluruh skenario ekstrim dapat disimulasikan seketika di memori lokal.",
    mental_model="Decoupled interfaces turn brittle integration tests into ultra-fast deterministic unit tests.",
    source_path="content/02-master-draft.md • tests/processor_test.go",
    slide_num="14 / 24",
    layout="before_after"
)

# =====================================================================
# SLIDE 15: Analisis Mendalam: Happy Path Isolation (Cards/Rows)
# =====================================================================
add_slide(
    categories=["ANALISIS TEST", "HAPPY PATH"],
    title="Test Suite 01: Verifikasi Happy Path Terisolasi",
    subtitle="Pengujian mengeksekusi transfer 50 USD tanpa membuka koneksi jaringan.",
    content_items=[
        {"type": "cards", "title": "VERIFIKASI SUKSES", "items": [
            {
                "title": "Substitusi Bersih",
                "desc": "mock := &MockGateway{} disuntikkan ke NewProcessor(mock); tidak ada panggilan TCP atau setup database."
            },
            {
                "title": "Validasi Status Mutasi",
                "desc": "mock.ChargedMoney mencatat persis Amount: 50 dan Currency: 'USD' setelah ProcessPayment(50) berhasil."
            },
            {
                "title": "Zero-Latency Run",
                "desc": "Eksekusi test selesai dalam 0.00s deterministik, bebas dari network flakiness atau rate-limit vendor."
            }
        ]}
    ],
    narrative="Test happy path membuktikan bahwa Processor mendelegasikan data secara akurat ke gateway tanpa perlu mengetahui implementasi di belakang interface.",
    mental_model="Mock substitution proves delegation accuracy without touching concrete infrastructure.",
    source_path="tests/processor_test.go:23-35",
    slide_num="15 / 24",
    layout="cards"
)

# =====================================================================
# SLIDE 16: Analisis Mendalam: Error Propagation (Cards/Rows)
# =====================================================================
add_slide(
    categories=["ANALISIS TEST", "ERROR PROPAGATION"],
    title="Test Suite 02: Propagasi Kegagalan Gateway",
    subtitle="Memverifikasi error dari infrastruktur diteruskan utuh tanpa ditelan.",
    content_items=[
        {"type": "cards", "title": "VERIFIKASI PROPAGASI ERROR", "items": [
            {
                "title": "Simulasi Down Time",
                "desc": "MockGateway dengan ShouldFail: true mensimulasikan kondisi 'gateway unavailable' secara instan."
            },
            {
                "title": "Error Forwarding",
                "desc": "Processor mengembalikan error gateway langsung ke caller tanpa menelan (swallow) pesan kesalahan."
            },
            {
                "title": "Audit Integritas",
                "desc": "Caller menerima sinyal kegagalan deterministik untuk memicu mekanisme retry atau fallback otomatis."
            }
        ]}
    ],
    narrative="Menguji kegagalan gateway pada arsitektur tightly coupled sangat sulit. Dengan DI, kegagalan vendor dapat disimulasikan hanya dengan mengaktifkan flag pada test double.",
    mental_model="Deterministic error simulation verifies resilient error propagation paths.",
    source_path="tests/processor_test.go:37-45",
    slide_num="16 / 24",
    layout="cards"
)

# =====================================================================
# SLIDE 17: Analisis Mendalam: Invalid Input Short-Circuit (Cards/Rows)
# =====================================================================
add_slide(
    categories=["ANALISIS TEST", "DEFENSIVE CHECK"],
    title="Test Suite 03: Short-Circuit Invalid Amount",
    subtitle="Memastikan gateway tidak pernah dipanggil jika parameter nominal salah.",
    content_items=[
        {"type": "cards", "title": "VERIFIKASI SHORT-CIRCUIT", "items": [
            {
                "title": "Input Negatif / Nol",
                "desc": "ProcessPayment(-10) menghasilkan error 'invalid amount' seketika pada boundary validasi."
            },
            {
                "title": "Assertion Gateway Call",
                "desc": "Assertion mock.ChargedMoney.Amount == 0 membuktikan Charge() tidak pernah disentuh sama sekali."
            },
            {
                "title": "Cost & Traffic Saving",
                "desc": "Mencegah overhead biaya API call vendor dan transmisi data yang ditakdirkan gagal sejak awal."
            }
        ]}
    ],
    narrative="Kombinasi assertion error dan pemeriksaan nilai 0 pada mock membuktikan bahwa validasi pre-check berhasil menghentikan alur sebelum mencapai lapisan infrastruktur.",
    mental_model="Asserting zero interaction with downstream services proves boundary validation effectiveness.",
    source_path="tests/processor_test.go:47-59",
    slide_num="17 / 24",
    layout="cards"
)

# =====================================================================
# SLIDE 18: Analisis Mendalam: Value Object Exemption (Cards/Rows)
# =====================================================================
add_slide(
    categories=["ANALISIS DESAIN", "VALUE OBJECT"],
    title="Desain Rasional: Value Object Bebas DI",
    subtitle="Money dan DTO murni dibuat langsung tanpa melewati container injector.",
    content_items=[
        {"type": "cards", "title": "PRINSIP VALUE OBJECT", "items": [
            {
                "title": "Data Murni Tanpa Efek",
                "desc": "Struct Money{Amount, Currency} tidak berinteraksi dengan database, network, atau state eksternal."
            },
            {
                "title": "Bahaya Over-Injection",
                "desc": "Memasukkan struct data ke dalam container DI menambah layer indirection tanpa keuntungan testability."
            },
            {
                "title": "Heuristik Desain",
                "desc": "Gunakan DI untuk Service ber-efek samping (I/O); gunakan direct instantiation untuk data value murni."
            }
        ]}
    ],
    narrative="Kesalahan umum pemula adalah menginjeksikan segalanya. Value object seperti Money, Date, atau Alamat adalah representasi data murni yang harus diinstansiasi langsung.",
    mental_model="Inject services that cross boundaries; instantiate pure value objects directly.",
    source_path="internal/di/gateway.go:5-9 • internal/di/processor.go:19",
    slide_num="18 / 24",
    layout="cards"
)

# =====================================================================
# SLIDE 19: Analisis Mendalam: Service Locator Coupling Risk (Cards/Rows)
# =====================================================================
add_slide(
    categories=["ANALISIS DESAIN", "ANTI-PATTERN"],
    title="Bahaya Service Locator: PSR-11 Warning",
    subtitle="Menyuntikkan container ke domain menyamarkan dependensi dan merusak pengujian.",
    content_items=[
        {"type": "cards", "title": "DAMPAK NEGATIF SERVICE LOCATOR", "items": [
            {
                "title": "Pelanggaran PSR-11",
                "desc": "Standar PHP-FIG PSR-11 menegaskan: 'Users SHOULD NOT pass a container into an object'."
            },
            {
                "title": "Double-Mocking Overhead",
                "desc": "Pengujian BadProcessor mewajibkan pembuatan MockContainer hanya untuk mengembalikan MockGateway."
            },
            {
                "title": "Kopling ke API Container",
                "desc": "Kelas bisnis menjadi terikat pada method GetPaymentGateway() milik container, bukan antarmuka generik."
            }
        ]}
    ],
    narrative="Service Locator adalah ilusi decoupling. Alih-alih terikat pada RealGateway, kelas justru terikat erat pada API Container yang lebih rumit untuk di-mock.",
    mental_model="Passing a container to domain logic replaces one tight coupling with an even broader one.",
    source_path="internal/di/locator.go • tests/processor_test.go:61-108",
    slide_num="19 / 24",
    layout="cards"
)

# =====================================================================
# SLIDE 20: Checklist Implementasi
# =====================================================================
add_slide(
    categories=["CHECKLIST", "BEST PRACTICE"],
    title="Checklist Adopsi Dependency Injection",
    subtitle="Panduan langkah praktis menerapkan DI yang bersih pada proyek produksi.",
    content_items=[
        {"title": "CHECKLIST ARSITEKTUR", "items": [
            "Deklarasikan dependensi sebagai interface, bukan tipe struct konkret.",
            "Gunakan Constructor Injection sebagai standar default; tolak setter opsional.",
            "Pusatkan instansiasi objek konkret hanya pada Composition Root (main.go).",
            "Jangan pernah mengoper Container atau Service Locator ke dalam domain logic.",
            "Instansiasi struct data murni (Value Objects) langsung tanpa melalui DI.",
            "Pastikan setiap method bisnis memiliki unit test dengan test double terisolasi."
        ], "type": "checklist"}
    ],
    narrative="Checklist ini menjadi panduan praktis dalam code review untuk mencegah bocornya dependensi konkret ke dalam domain logic tim.",
    mental_model="Follow a disciplined checklist to keep architectural boundaries clean and tests fast.",
    source_path="content/02-master-draft.md:155-167",
    slide_num="20 / 24",
    layout="checklist"
)

# =====================================================================
# SLIDE 21: Anti-Pattern & Cara Reproduksi Kegagalan
# =====================================================================
add_slide(
    categories=["ANTI-PATTERN", "PITFALLS"],
    title="Tiga Anti-Pattern Umum Seputar DI",
    subtitle="Kesalahan arsitektural yang sering ditemui saat mengadopsi pola injeksi.",
    content_items=[
        {"type": "cards", "title": "ANTI-PATTERN UTAMA", "items": [
            {
                "title": "Constructor Bloat",
                "desc": "Konstruktor dengan belasan parameter adalah sinyal kuat pelanggaran Single Responsibility (SRP)."
            },
            {
                "title": "Service Locator Disguise",
                "desc": "Membungkus dependensi ke dalam satu Map atau Container serbaguna demi menyederhanakan signature ctor."
            },
            {
                "title": "Over-Abstracting Everything",
                "desc": "Membuat interface untuk setiap kelas tanpa pernah memiliki lebih dari satu implementasi atau kebutuhan mock."
            }
        ]}
    ],
    narrative="DI bukan alasan untuk membiarkan sebuah kelas menangani terlalu banyak hal. Jika parameter konstruktor membengkak, solusinya adalah memecah kelas, bukan memakai Service Locator.",
    mental_model="Constructor bloat is a design warning: split the class, don't hide the bloat in a container.",
    source_path="content/02-master-draft.md:135-146",
    slide_num="21 / 24",
    layout="cards"
)

# =====================================================================
# SLIDE 22: Takeaway Praktis & Command Block Dominan
# =====================================================================
add_slide(
    categories=["REPRODUKSI", "COMMAND"],
    title="Reproduksi Mandiri Suite Pengujian",
    subtitle="Verifikasi seluruh invarian dependency injection secara mandiri di mesin lokal.",
    content_items=[
        {"title": "COMMAND RUNNER (GO 1.22+)", "items": [
            "cd labs/16-dependency-injection",
            "go test -v ./...       # Verifikasi 6 skenario test suite",
            "go test -race ./...    # Verifikasi thread-safety zero data race",
            "go run ./cmd/demo      # Eksekusi demo komparasi Ctor vs Locator"
        ], "type": "code"}
    ],
    narrative="Lab ini sepenuhnya berbasis Go standard library tanpa external dependencies. Seluruh verifikasi dapat direproduksi dalam hitungan detik di lingkungan apa pun.",
    mental_model="Zero-dependency standard tests guarantee instant, reproducible verification anywhere.",
    source_path="README.md • go.mod • tests/processor_test.go",
    slide_num="22 / 24",
    layout="code"
)

# =====================================================================
# SLIDE 23: Prinsip Arsitektural Jangka Panjang
# =====================================================================
add_slide(
    categories=["PRINSIP", "LONG-TERM"],
    title="Prinsip Arsitektur: Hollywood Principle",
    subtitle="Don't call us, we'll call you — membalikkan arah kendali dependensi sistem.",
    content_items=[
        {"title": "PRINSIP REKAYASA PERANGKAT LUNAK", "items": [
            "Inversion of Control: Logika bisnis fokus mengeksekusi use case; framework/assembler yang menyediakan instrumen.",
            "Program to Interfaces: Modul tingkat tinggi tidak boleh bergantung pada modul tingkat rendah.",
            "Evolvability: Kemampuan menukar provider atau infrastruktur tanpa memicu modifikasi pada logika inti domain."
        ], "type": "bullet"}
    ],
    narrative="Dependency Injection adalah perwujudan konkret dari Dependency Inversion Principle. Dengan membalik kontrol pembuatan objek, sistem menjadi adaptif terhadap perubahan zaman.",
    mental_model="Invert the control of construction to protect domain logic from infrastructure churn.",
    source_path="content/05-key-takeaways.md",
    slide_num="23 / 24",
    layout="bullet"
)

# =====================================================================
# SLIDE 24: Closing Slide (Dark Navy)
# =====================================================================
add_slide(
    categories=["CLOSING"],
    title="Pisahkan perakitan dari pemakaian; isolasi logika dari infrastruktur.",
    subtitle="",
    content_items=[
        {"title": "MODUL LAB 16", "items": [
            "internal/di/processor.go • internal/di/gateway.go • internal/di/locator.go",
            "cmd/demo/main.go • tests/processor_test.go • content/02-master-draft.md"
        ], "type": "bullet"},
        {"title": "VERIFIKASI & STATUS", "items": [
            f"Commit: {COMMIT_SHA[:8]} | Go 1.22 Standard Library",
            "Test Suite: 6 cases passing, 0 failures, 100% deterministic.",
            "Race Detector: Clean pass (go test -race ./...)",
            f"Repository: {REPO_URL}"
        ], "type": "bullet"}
    ],
    mental_model="Separate configuration from use: outside world assembles, inner logic executes.",
    source_path=REPO_URL,
    slide_num="24 / 24",
    layout="closing"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
