#!/usr/bin/env python3
"""
Publication Builder for Software Engineering Lab 16: Dependency Injection.
Generates: PPTX, PDF, PNG covers, Dev.to MD, Medium HTML, Metadata, Caption
"""
import os
import sys

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/16-dependency-injection/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/16-dependency-injection"
COMMIT_SHA = "8772235"
LAB_NO = "16"
SLUG = "dependency-injection"
TITLE_ID = "Dependency Injection"
SUBTITLE_ID = "Inversi Kontrol & Modularitas Kode"
TITLE_EN = "Dependency Injection"
SUBTITLE_EN = "Inversion of Control & Modular Code"
AUTHOR = "Lukman / lukman-ss"
TOPIC_TAGS = "DI • IoC • Go • Testing"

SLIDES = []

def add_slide(categories, title, subtitle, content_items, mental_model="", source_path="", slide_num="", narrative="", narrative_en="", layout=""):
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

# Slide 1: Cover Slide
add_slide(
    categories=["COVER"],
    title=TITLE_ID,
    subtitle=SUBTITLE_ID,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "CRT computer showing modular circuit with injectable components plugged into sockets",
            "Hewlett Packard CRT monitor with physical cables connecting hardware modules",
            f"Penulis: {AUTHOR}"
        ]}
    ],
    mental_model="",
    source_path="",
    slide_num="01 / 24"
)

# Slide 2: Taped Architecture Loss
add_slide(
    categories=["PROBLEM", "COUPLING"],
    title="Konteks Arsitektur Hilang Tanpa Jejak",
    subtitle="Hard-coded dependencies menghabiskan testability dan swap implementasi.",
    content_items=[
        {"title": "MASALAH NYATA", "items": [
            "Kelas membuat objek dependensinya sendiri → tight coupling.",
            "Ganti payment gateway memaksa ubah banyak tempat.",
            "Test menjadi integration test yang lambat."
        ], "type": "bullet"}
    ],
    narrative="Ketika kelas menginstansiasi dependensinya sendiri, alasan teknisnya tidak terdokumentasi. Tim lama mengerti, tapi engineer baru bingung.",
    mental_model="Tight coupling erases architectural context from code.",
    source_path="content/02-master-draft.md",
    slide_num="02 / 24"
)

# Slide 3: Why This Matters
add_slide(
    categories=["WHY MATTERS"],
    title="Mengapa Ini Penting",
    subtitle="Coupling keras menghambat evolvability sistem jangka panjang.",
    content_items=[
        {"title": "DAMPAG", "items": [
            "Test lambat karena bergantung pada status eksternal (API, database).",
            "Perubahan vendor memicu modifikasi menyeluruh lapisan bisnis.",
            "Arsitektur kaku karena tiap consumer terikat pada implementasi spesifik."
        ], "type": "bullet"}
    ],
    mental_model="Separate construction from use for maintainable evolution.",
    source_path="content/02-master-draft.md",
    slide_num="03 / 24"
)

# Slide 4: Mental Model
add_slide(
    categories=["CONCEPT"],
    title="Pemisahan Konfigurasi dari Penggunaan",
    subtitle="Fase konfigurasi = buat objek; fase penggunaan = jalankan logika.",
    content_items=[
        {"title": "PRINSIP", "items": [
            "Modul logika tidak boleh mencari atau membangun dependensinya.",
            "Assembler/Composition Root yang bertanggung jawab merakit object graph.",
            "Hollywood Principle: Don't call us, we'll call you."
        ], "type": "bullet"}
    ],
    mental_model="External injection enables loose coupling.",
    source_path="content/02-master-draft.md",
    slide_num="04 / 24"
)

# Slide 5: DI vs Service Locator
add_slide(
    categories=["COMPARISON", "DI"],
    title="DI vs Service Locator",
    subtitle="Keduanya decorrelate, tapi mekanisme berbeda.",
    content_items=[
        {"title": "CONSTRUCTOR INJECTION", "items": [
            "Dependensi muncul lewat parameter constructor.",
            "Tidak ada request eksplisit ke container.",
            "Jelas: Processor(PaymentGateway)"
        ], "type": "before_after"},
        {"title": "SERVICE LOCATOR", "items": [
            "Consumer punya dependensi pada locator.",
            "Memanggil container.Get(...) secara aktif.",
            "Tersembunyi: processor hanya lihat Container"
        ], "type": "before_after"}
    ],
    mental_model="DI hides look-up; Service Locator exposes it explicitly.",
    source_path="content/02-master-draft.md",
    slide_num="05 / 24",
    layout="before_after"
)

# Slide 6: Interface & Value Object
add_slide(
    categories=["CODE", "MODEL"],
    title="Interface & Value Object",
    subtitle="Money is VO; PaymentGateway is contract.",
    content_items=[
        {"type": "code", "items": [
            "type Money struct {",
            "    Amount   int",
            "    Currency string",
            "}",
            "type PaymentGateway interface {",
            "    Charge(m Money) error",
            "}"
        ]}
    ],
    mental_model="Interface enables swap; vo created directly without DI.",
    source_path="content/03-code-snippets.md",
    slide_num="06 / 24"
)

# Slide 7: Constructor Injection Code
add_slide(
    categories=["CODE"],
    title="Constructor Injection",
    subtitle="Processor menerima gateway lewat constructor.",
    content_items=[
        {"type": "code", "items": [
            "type Processor struct {",
            "    gateway PaymentGateway",
            "}",
            "func NewProcessor(g PaymentGateway) *Processor {",
            "    return &Processor{gateway: g}",
            "}",
            "func (p *Processor) ProcessPayment(amount int) error {",
            "    if amount <= 0 { return errors.New(\"invalid amount\") }",
            "    m := Money{Amount: amount, Currency: \"USD\"}",
            "    return p.gateway.Charge(m) }"
        ]}
    ],
    mental_model="Constructor injection ensures valid state and explicit deps.",
    source_path="content/03-code-snippets.md",
    slide_num="07 / 24"
)

# Slide 8: Failure Scenario
add_slide(
    categories=["FAILURE"],
    title="Skenario Kegagalan: Invalid Input",
    subtitle="Validasi amount <= 0 mencegah panggilan sia-sia.",
    content_items=[
        {"type": "metric_cards", "items": [
            {"val": "ProcessPayment(-10)", "label": "Result", "sub": "Error: invalid amount, gateway NOT called"},
            {"val": "ProcessPayment(100)", "label": "Result", "sub": "Money created, gateway.Charge() invoked"}
        ]}
    ],
    mental_model="Input validation blocks infrastructure calls cleanly.",
    source_path="tests/processor_test.go",
    slide_num="08 / 24"
)

# Slide 9: Test Evidence
add_slide(
    categories=["TEST", "EVIDENCE"],
    title="Bukti: Isolated Testing Tanpa Infrastruktur",
    subtitle="MockGateway menggantikan RealGateway secara penuh.",
    content_items=[
        {"type": "code", "items": [
            "mock := &MockGateway{}",
            "proc := di.NewProcessor(mock)",
            "err := proc.ProcessPayment(50)",
            "if err != nil { t.Fatalf(\"expected nil, got %v\", err) }",
            "if mock.ChargedMoney.Amount != 50 { t.Errorf(\"wrong\")"
        ]}
    ],
    mental_model="Interface injection enables fully isolated unit tests.",
    source_path="content/03-code-snippets.md",
    slide_num="09 / 24"
)

# Slide 10: Boundary of Claims
add_slide(
    categories=["QUALITY", "BOUNDARY"],
    title="Batasan Klaim",
    subtitle="Tidak ada jaminan produksi otomatis.",
    content_items=[
        {"title": "PERNYATAAN YANG DAPAT DIBUKTIKAN", "items": [
            "DI memudahkan testing tanpa infrastructure nyata.",
            "MockGateway menggantikan network call secara penuh.",
            "Tes berjalan dalam 0.00s per kasus."
        ], "type": "bullet"},
        {"title": "BATASAN BUKTI", "items": [
            "\"bukan berarti\" behaving di production sama dengan test.",
            "\"tidak otomatis\" performa atau defect rate menurun.",
            "\"perlu bukti tambahan\" untuk skala enterprise."
        ], "type": "bullet"}
    ],
    mental_model="Lab membuktikan testability; production butuh validasi terpisah.",
    source_path="modules/quality-bar.md",
    slide_num="10 / 24"
)

# Slide 11: What Tests Do Not Prove
add_slide(
    categories=["BATASAN BUKTI", "DO NOT OVERCLAIM"],
    title="Batas Integritas",
    subtitle="Tes tidak membuktikan keamanan produksi.",
    content_items=[
        {"title": "DI LUAR CAKUPAN", "items": [
            "Lifecycle management singleton/scoped/transient.",
            "Performance overhead container vs direct new.",
            "Pengurangan defect rate akibat DI."
        ], "type": "bullet"}
    ],
    mental_model="Passing unit tests validate structure; production needs holistic validation.",
    source_path="content/02-master-draft.md",
    slide_num="11 / 24"
)

# Slide 12: Implementation Approach
add_slide(
    categories=["IMPLEMENTATION"],
    title="Manual Wiring Tanpa Framework",
    subtitle="YAGNI: Go stand-alone cukup membuktikan konsep.",
    content_items=[
        {"title": "ALUR KERJA", "items": [
            "main.go menginisialisasi RealGateway.",
            "NewProcessor(realGateway) – Constructor Injection.",
            "Output: RealGateway charging 100 USD"
        ], "type": "bullet"}
    ],
    mental_model="Manual wiring bukan berarti tidak profesional; type-safe dan compile-time valid.",
    source_path="content/02-master-draft.md",
    slide_num="12 / 24"
)

# Slide 13: Case Study
add_slide(
    categories=["CASE STUDY"],
    title="Payment Processor: Contrast",
    subtitle="Processor vs BadProcessor pada kasus nyata.",
    content_items=[
        {"title": "CONTOH", "items": [
            "Processor: NewProcessor(realGateway)",
            "BadProcessor: NewBadProcessor(container)",
            "Output keduanya identik: RealGateway charging ..."
        ], "type": "bullet"}
    ],
    mental_model="DI pattern yields better testability; Service Locator hides dependencies.",
    source_path="content/02-master-draft.md",
    slide_num="13 / 24"
)

# Slide 14: Code Review
add_slide(
    categories=["CODE"],
    title="Review Kode: Anti-Pattern",
    subtitle="BadProcessor menyembunyikan dependensi di container.",
    content_items=[
        {"type": "code", "items": [
            "type BadProcessor struct {",
            "    container Container",
            "}",
            "func (p *BadProcessor) ProcessPayment(amount int) error {",
            "    return p.container.GetPaymentGateway().Charge(...)",
            "}"
        ]}
    ],
    mental_model="Service Locator membuat dependensi tak terlihat pada signature.",
    source_path="content/03-code-snippets.md",
    slide_num="14 / 24"
)

# Slide 15: Production Checklist
add_slide(
    categories=["CHECKLIST"],
    title="Checklist DI dalam Tim",
    subtitle="Langkah implementasi yang direkomendasikan.",
    content_items=[
        {"title": "WAJIB", "items": [
            "Deklarasikan dependensi sebagai interface.",
            "Gunakan constructor injection sebagai default.",
            "Composition root SATU-satunya tempat buat implementasi konkret."
        ], "type": "checklist"},
        {"title": "HINDARI", "items": [
            "Container.Get(...) di logic domain.",
            "Inject value object ke container.",
            "Constructor dengan banyak parameter (>12 sebagai sinyal)."
        ], "type": "checklist"}
    ],
    mental_model="Checklist memastikan DI implementation konsisten dan testable.",
    source_path="content/02-master-draft.md",
    slide_num="15 / 24",
    layout="checklist"
)

# Slide 16: Common Mistakes
add_slide(
    categories=["MISTAKES"],
    title="Kesalahan Umum",
    subtitle="Yang sering dilakukan developer baru.",
    content_items=[
        {"title": "1. Service Locator sebagai DI", "items": [
            "Inject container ke kelas; panggil container.Get() di method.",
            "Solusi: injeksikan interface spesifik yang dibutuhkan."
        ], "type": "bullet"},
        {"title": "2. Over-abstraksi value object", "items": [
            "Masukkan Money, DateTime, Address ke container.",
            "Solusi: buat langsung, tidak butuh interface."
        ], "type": "bullet"}
    ],
    mental_model="Abstraction tanpa manfaat = overhead yang tidak perlu.",
    source_path="content/02-master-draft.md",
    slide_num="16 / 24"
)

# Slide 17: Trade-offs
add_slide(
    categories=["TRADE-OFF"],
    title="Biaya DI",
    subtitle="Tidak ada Solution gratis; ada trade-off.",
    content_items=[
        {"title": "KEUNTUNGAN", "items": [
            "Testing terisolasi tanpa infrastructure.",
            "Swap implementasi lewat composition root.",
            "Clear dependency contracts."
        ], "type": "bullet"},
        {"title": "KOS-KOSAN", "items": [
            "Konfigurasi ekstra di main.go.",
            "Tracing kesulitan (behavior dipisah dari construction).",
            "Refactoring penting untuk memilih DI framework."
        ], "type": "bullet"}
    ],
    mental_model="DI improves testability; adds configuration overhead that must be managed.",
    source_path="content/02-master-draft.md",
    slide_num="17 / 24"
)

# Slide 18: Conclusion
add_slide(
    categories=["CONCLUSION"],
    title="Kesimpulan",
    subtitle="DI meningkatkan modularitas, bukan penyelesaian semua masalah.",
    content_items=[
        {"title": "INSIGHT UTAMA", "items": [
            "Constructor injection default untuk testability.",
            "Avoid Service Locator di domain logic.",
            "Mock injection = testing isolation tanpa network.",
            "Value objects langsung dibuat, tidak lewat DI."
        ], "type": "bullet"}
    ],
    mental_model="DI memisahkan konfigurasi dari penggunaan; bukan silver bullet.",
    source_path="content/05-key-takeaways.md",
    slide_num="18 / 24"
)

# Slide 19: Key Takeaways
add_slide(
    categories=["TAKEAWAYS"],
    title="8 Poin Kunci",
    subtitle="Ringkasan prinsip Dependency Injection.",
    content_items=[
        {"type": "metric_cards", "items": [
            {"val": "1", "label": "Sepetakan konfigurasi dari penggunaan"},
            {"val": "2", "label": "IoC umum; DI spesifik untuk dependensi"},
            {"val": "3", "label": "Constructor injection default"},
            {"val": "4", "label": "Service Locator dilarang di domain"},
            {"val": "5", "label": "Program to interfaces"},
            {"val": "6", "label": "Mock injection untuk pengujian"},
            {"val": "7", "label": "Value objects langsung dibuat"},
            {"val": "8", "label": "Over-injection = SRP violation"}
        ]}
    ],
    mental_model="Eight core principles guide DI adoption in Go applications.",
    source_path="content/05-key-takeaways.md",
    slide_num="19 / 24"
)

# Slide 20: Architecture Components
add_slide(
    categories=["ARCHITECTURE"],
    title="Komponen DI",
    subtitle="Tabel komponen dan lokasinya.",
    content_items=[
        {"type": "table", "items": [
            "Komponen | Lokasi | Peran",
            "PaymentGateway | internal/di/gateway.go | Interface abstraksi payment provider",
            "RealGateway | internal/di/gateway.go | Implementasi produksi",
            "Money | internal/di/gateway.go | Value object",
            "Processor | internal/di/processor.go | Service dengan Constructor Injection",
            "Container | internal/di/locator.go | Interface Service Locator",
            "BadProcessor | internal/di/locator.go | Service dengan Anti-pattern",
            "MockGateway | tests/processor_test.go | Test double"
        ]}
    ],
    mental_model="Clear component mapping enables modular testing.",
    source_path="content/02-master-draft.md",
    slide_num="20 / 24"
)

# Slide 21: Reproduction Guide
add_slide(
    categories=["REPRODUKSI"],
    title="Cara Menduplikasi",
    subtitle="Jalankan langkah ini di mesin lokal.",
    content_items=[
        {"type": "code", "items": [
            "cd labs/16-dependency-injection",
            "go test -v ./...    # 6 test cases, semua PASS",
            "go test -race ./... # race-detector clean",
            "go run ./cmd/demo   # jalankan demo DI vs Service Locator"
        ]}
    ],
    mental_model="Reproduction requires only Go standard library, no external dependencies.",
    source_path="README.md",
    slide_num="21 / 24"
)

# Slide 22: Verification Results
add_slide(
    categories=["VERIFICATION"],
    title="Hasil Verifikasi",
    subtitle="Semua test lulus dengan race detector bersih.",
    content_items=[
        {"type": "metric_cards", "items": [
            {"val": "6", "label": "Test Cases", "sub": "3 untuk Processor • 3 untuk BadProcessor"},
            {"val": "0", "label": "Failures", "sub": "Semua lolos"},
            {"val": "PASS", "label": "Race Detector", "sub": "go test -race bersih"}
        ]}
    ],
    mental_model="Comprehensive test coverage validates DI implementation behavior.",
    source_path="tests/processor_test.go",
    slide_num="22 / 24"
)

# Slide 23: Before/After
add_slide(
    categories=["BEFORE AFTER"],
    title="Sebelum vs Sesudah DI",
    subtitle="Perbandingan pola yang direkomendasikan.",
    content_items=[
        {"title": "TANPA DI", "items": [
            "Processor membuat gateway sendiri.",
            "Hard-coded pada RealGateway.",
            "Test lambat, butuh network."
        ], "type": "before_after"},
        {"title": "DENGAN DI", "items": [
            "Processor menerima interface.",
            "Swappable lewat constructor.",
            "Test cepat, mock-only."
        ], "type": "before_after"}
    ],
    mental_model="DI mengekspos dependency lewat kontrak, bukan hard-code.",
    source_path="content/02-master-draft.md",
    slide_num="23 / 24",
    layout="before_after"
)

# Slide 24: Closing
add_slide(
    categories=["CLOSING"],
    title="DI untuk Kode Lebih Bersih",
    subtitle="Memisahkan konfigurasi dari penggunaan adalah investasi pada evolvability.",
    content_items=[
        {"title": "MODUL LAB", "items": [
            "internal/di/gateway.go • internal/di/processor.go • internal/di/locator.go",
            "cmd/demo/main.go • tests/processor_test.go",
            "README.md • content/*"
        ], "type": "bullet"},
        {"title": "STATUS", "items": [
            f"Sumber: commit {COMMIT_SHA}",
            "Test passing: 6 test cases, 0 failures.",
            "Race detector: clean (go test -race ./...)",
            "Manual wiring demo; production may need framework."
        ], "type": "bullet"}
    ],
    mental_model="DI enables modular, testable code but adds configuration overhead.",
    source_path=REPO_URL,
    slide_num="24 / 24"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")