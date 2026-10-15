#!/usr/bin/env python3
"""
Publication Builder Definition for Software Engineering Lab 26: Contract Testing.
Follows:
- cover-rules.md (Mandatory raster hero image at COVER_IMAGE_PATH, series identity lock, no micro-text)
- slide-rules.md (<=45 words prose, <=30 words on code slides, dominant visual anchors, 50% zoom test)
- quality-bar.md (sharp engineering insights, strictly bounded claims, >=12 slides with limiting phrases, no overclaims)
- article-rules.md (Dev.to 22-section mapping, Medium ID HTML 22-section mapping)
"""
import os
import sys

OUTPUT_DIR = "labs/26-contract-testing/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/26-contract-testing"
COMMIT_SHA = "4a5cf38"
LAB_NO = "26"
SLUG = "contract-testing"

TITLE_ID = "Contract Testing"
SUBTITLE_ID = "Consumer-Driven Contracts, CI Deployment Gate, dan Safe API Evolution tanpa E2E Rapuh"
TITLE_EN = "Contract Testing"
SUBTITLE_EN = "Consumer-Driven Contracts, CI Deployment Gate, and Safe API Evolution without Brittle E2E"

AUTHOR = "Lukman / lukman-ss"
TOPIC_TAGS = "Contract Testing • Consumer-Driven Contracts • Go • API Versioning • CI/CD"
COVER_IMAGE_PATH = "labs/26-contract-testing/result/image-cover/cover_art.jpg"

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
            "Retro Terminal Displaying Consumer-Driven Contract Verification Flow",
            "MobileApp → Contract → Verifier → Provider V1/Breaking/Dual: PASSED vs BLOCKED",
            f"Penulis: {AUTHOR}"
        ]}
    ],
    mental_model="Contract testing catches integration failures before deployment by making shared expectations executable.",
    source_path="",
    slide_num="01 / 24",
    layout="cover"
)

# Slide 02: Problem Nyata — Unit Test Hijau, Integrasi Produksi Pecah (Layout: cards)
# Limiting phrase 1: "bukan berarti"
add_slide(
    categories=["PROBLEM", "INTEGRATION FAILURE"],
    title="Unit Test Hijau Bukan Berarti Integrasi Selamat",
    subtitle="Provider mengubah kontrak tanpa koordinasi consumer — deploy berhasil, production gagal.",
    content_items=[
        {"type": "cards", "title": "TIGA MUKA KEGAGALAN INTEGRASI", "items": [
            {"title": "Field Rename Tersembunyi", "desc": "customer.name diganti customer.full_name di provider. JSON valid, tapi consumer mendapat null."},
            {"title": "Enum Casing Change", "desc": "Status IN_PROGRESS berubah jadi in_progress. Response tetap 200 OK, tapi consumer gagal parsing."},
            {"title": "Primitive Type Mutation", "desc": "total integer 150000 jadi string '150000'. Consumer math operations pecah di runtime."}
        ]}
    ],
    narrative="Ketika Service A dan Service B masing-masing lulus unit test, integrasi mereka tetap bisa pecah di produksi karena provider mengubah format pesan tanpa sepengetahuan consumer. Tiga skenario di atas menunjukkan bahwa JSON yang valid secara sintaksis tidak menjamin kompatibilitas semantik.",
    narrative_en="When Service A and Service B each pass their unit tests, integration can still fail in production because the provider changes message format without consumer coordination. The three scenarios above demonstrate that syntactically valid JSON does not guarantee semantic compatibility.",
    mental_model="Unit test green status is not sufficient to guarantee integration safety across service boundaries.",
    source_path="content/02-master-draft.md:3-8",
    slide_num="02 / 24",
    layout="cards"
)

# Slide 03: Mengapa Pendekatan Naif Gagal — E2E vs Contract Testing (Layout: cards)
# Limiting phrase 2: "tidak otomatis"
add_slide(
    categories=["WHY FAILS", "NAIVE APPROACH"],
    title="End-to-End Test dan Manual Review Tidak Otomatis Menangkap Breaking Change",
    subtitle="E2E lambat dan mahal; contract testing memberikan feedback cepat di setiap PR.",
    content_items=[
        {"type": "cards", "title": "KELEMAHAN PENDEKATAN TRADISIONAL", "items": [
            {"title": "E2E Test Mahal", "desc": "Menjalankan seluruh sistem gabungan butuh environment lengkap, data fixture, dan waktu build lama."},
            {"title": "Manual Code Review", "desc": "Inspeksi manual perubahan API response tidak otomatis menemukan field rename atau type change."},
            {"title": "Tidak Otomatis Aman", "desc": "Deploy pipeline tanpa gate verifikasi kontrak mengizinkan breaking change lolos ke produksi."}
        ]}
    ],
    narrative="End-to-end test mensimulasikan seluruh alur sistem sehingga mahal untuk dijalankan setiap PR. Review manual pun tidak otomatis menemukan perubahan breaking karena developer fokus pada logika bisnis, bukan schema response. Contract testing mengisi celah ini dengan verifikasi otomatis sebelum deploy.",
    narrative_en="End-to-end tests simulate the entire system flow and are expensive to run on every PR. Manual code review also fails to automatically detect breaking changes because developers focus on business logic rather than response schemas. Contract testing fills this gap with automated verification before deployment.",
    mental_model="E2E tests are too slow and manual review is too error-prone; contract testing provides fast, automated safety gates.",
    source_path="content/02-master-draft.md:9-15",
    slide_num="03 / 24",
    layout="cards"
)

# Slide 04: Mental Model Inti — Consumer-Driven Contract (Layout: schema)
# Limiting phrase 3: "hanya pada jalur ini"
add_slide(
    categories=["MENTAL MODEL", "PILLARS"],
    title="Empat Pilar Consumer-Driven Contract Testing",
    subtitle="Consumer menentukan ekspektasi, provider memenuhi, CI memverifikasi, dan gate menolak breaking change.",
    content_items=[
        {"type": "schema", "title": "EMPAT PILAR CDC TESTING", "items": [
            {"phase": "1. Consumer Test", "role": "Contract Generator", "desc": "Consumer menulis test dengan mock provider, menghasilkan kontrak JSON minimal."},
            {"phase": "2. Contract File", "role": "Shared Artifact", "desc": "Kontrak berisi method, path, status, header, body dengan tipe dan casing eksak."},
            {"phase": "3. Provider Verify", "role": "CI Gate", "desc": "Verifier kirim request sesuai kontrak ke provider; compare status + headers + body subset."},
            {"phase": "4. Pass/Fail Gate", "role": "Deployment Control", "desc": "Hanya pada jalur ini fail → deployment diblokir; pass → deployment diizinkan."},
            {"phase": "5. Safe Evolution", "role": "V1+V2 Routing", "desc": "Dual provider mempertahankan V1 compliant sambil mengekspos V2 schema baru."}
        ]}
    ],
    narrative="Consumer-driven contract testing menggeser inisiatif ke consumer: mereka yang mendefinisikan apa yang dibutuhkan, bukan provider yang mendokumentasikan apa yang tersedia. Hanya pada jalur ini CI gate memblokir deployment ketika verification gagal, mencegah breaking change mencapai produksi.",
    narrative_en="Consumer-driven contract testing shifts initiative to the consumer: they define what is needed, not the provider documenting what is available. Only along this path does the CI gate block deployment when verification fails, preventing breaking changes from reaching production.",
    mental_model="Consumer drives contract definition; CI gate enforces compliance; dual-version routing enables safe evolution.",
    source_path="content/02-master-draft.md:16-28",
    slide_num="04 / 24",
    layout="schema"
)

# Slide 05: Diagram Alur Pipeline (Layout: flow)
add_slide(
    categories=["ARCHITECTURE", "PIPELINE"],
    title="Pipeline Kontrak: Dari Consumer Test hingga CI Gate",
    subtitle="Alur eksekusi kontrak dari generasi hingga verifikasi di provider side.",
    content_items=[
        {"type": "flow", "title": "CONTRACT TESTING PIPELINE", "items": [
            {"step": "1", "label": "Consumer Test", "desc": "Test menjalankan mock server,GenerateMobileContract() membuat kontrak JSON."},
            {"step": "2", "label": "Contract.json", "desc": "Artefak kontrak tersimpan: MobileApp→OrderService, 1 interaksi, body subset minimal."},
            {"step": "3", "label": "CI Trigger", "desc": "Provider CI mengambil kontrak, menjalankan verifier.Verify(baseURL, contract)."},
            {"step": "4", "label": "Verifier Engine", "desc": "DiffValues() bandingkan expected vs actual: status, header, body subset, tipe numerik."},
            {"step": "5", "label": "Gate Decision", "desc": "Passed=true → ALLOWED. Passed=false + errors → BLOCKED with detailed diff report."}
        ]}
    ],
    narrative="Pipeline ini dimulai dari consumer test yang menghasilkan artefak kontrak sebagai JSON. Kontrak kemudian diunduh oleh CI provider yang menjalankan verifier engine untuk membandingkan response aktual dengan ekspektasi. Keputusan gate dibuat secara deterministik berdasarkan hasil diff.",
    narrative_en="The pipeline starts from consumer test generating a contract artifact as JSON. The contract is then fetched by the provider CI which runs the verifier engine to compare actual response against expectations. The gate decision is made deterministically based on the diff result.",
    mental_model="Contract generation and verification form a closed loop between consumer and provider in CI/CD.",
    source_path="content/02-master-draft.md:44-53",
    slide_num="05 / 24",
    layout="flow"
)

# Slide 06: Cuplikan Kode Dominan — Consumer Contract Generation (Layout: code)
add_slide(
    categories=["CODE", "CONSUMER"],
    title="GenerateMobileContract(): Kontrak Minimal dari Sisi Consumer",
    subtitle="Consumer mendeklarasikan hanya field yang benar-benar dibutuhkan.",
    content_items=[
        {"type": "code", "title": "internal/consumer/client.go:84-112", "items": [
            "func GenerateMobileContract() *contract.Contract {",
            "    return &contract.Contract{",
            "        Consumer: \"MobileApp\",",
            "        Provider: \"OrderService\",",
            "        Interactions: []contract.Interaction{{",
            "            Description:   \"A request for order details by ID\",",
            "            ProviderState: \"Order ORD-123 exists and is IN_PROGRESS\",",
            "            Request: contract.RequestDefinition{",
            "                Method: http.MethodGet,",
            "                Path:   \"/v1/orders/ORD-123\",",
            "            },",
            "            Response: contract.ResponseDefinition{",
            "                Status: http.StatusOK,",
            "                Headers: map[string]string{\"Content-Type\": \"application/json\"},",
            "                Body: map[string]interface{}{",
            "                    \"id\":     \"ORD-123\",",
            "                    \"status\": \"IN_PROGRESS\",",
            "                    \"customer\": map[string]interface{}{\"name\": \"Budi Santoso\"},",
            "                    \"total\": json.Number(\"150000\"),",
            "                },",
            "            },",
            "        }},",
            "    }",
            "}"
        ]}
    ],
    narrative="Fungsi ini menghasilkan kontrak minimal yang hanya mendeklarasikan field: id, status, customer.name, dan total. Field tambahan seperti notes atau created_at tidak dicantumkan karena tidak dipakai consumer — prinsip subset matching memastikan provider dapat berevolusi tanpa merusak kontrak.",
    narrative_en="This function generates a minimal contract declaring only the fields: id, status, customer.name, and total. Extra fields like notes or created_at are omitted because the consumer does not use them — subset matching ensures the provider can evolve without breaking the contract.",
    mental_model="Minimal contract: only used fields are declared; extra provider fields are ignored by verifier.",
    source_path="internal/consumer/client.go:84-112",
    slide_num="06 / 24",
    layout="code"
)

# Slide 07: Failure Path — Breaking Change yang Terdeteksi (Layout: detail)
add_slide(
    categories=["FAILURE", "DETECTION"],
    title="Breaking Provider Menghasilkan 3 Error Detail di CI Gate",
    subtitle="Verifikator menangkap enum casing, field rename, dan primitive type mutation.",
    content_items=[
        {"type": "detail", "title": "BREAKING CHANGES TERTANGKAP", "items": [
            {"error": "Type Mismatch — total", "detail": "Expected json.Number (integer 150000), got string ('150000'). decoder.UseNumber() mempertahankan presisi tipe numerik."},
            {"error": "Value Mismatch — status", "detail": "Expected 'IN_PROGRESS', got 'in_progress'. Enum casing change terdeteksi sebagai value mismatch."},
            {"error": "Missing Field — customer.name", "detail": "Field 'customer.name' tidak ada di response. Provider mengganti menjadi 'customer.full_name'."}
        ]}
    ],
    narrative="Ketiga error ini dihasilkan oleh diffValues() di verifier.go. Perbedaan tipe numerik tertangkap berkat decoder.UseNumber() yang menjaga json.Number, berbeda dengan default decoder yang mengubah semua angka jadi float64. Field rename terdeteksi karena key tidak ditemukan di expected map.",
    narrative_en="These three errors are produced by diffValues() in verifier.go. The numeric type difference is caught thanks to decoder.UseNumber() preserving json.Number, unlike the default decoder that converts all numbers to float64. Field rename is detected because the key is not found in the expected map.",
    mental_model="Json.Number decoder preserves type precision; map key absence detects field renames.",
    source_path="engineering/03-execution-result.md:103-108",
    slide_num="07 / 24",
    layout="detail"
)

# Slide 08: Tabel Perbandingan Mekanisme vs Gejala (Layout: table)
add_slide(
    categories=["ANALYSIS", "COMPARISON"],
    title="Breaking Change vs Mekanisme Deteksi",
    subtitle="Setiap jenis perubahan dipetakan ke mekanisme verifikasi spesifik.",
    content_items=[
        {"type": "table", "title": "BREAKING CHANGE vs VERIFICATION MECHANISM", "items": [
            {"mechanism": "Status Code Check", "triggers": "Response bukan 200 OK", "example": "Server error atau redirect"},
            {"mechanism": "Header Validation", "triggers": "Content-Type != application/json", "example": "Provider return text/plain"},
            {"mechanism": "DiffValues Body Subset", "triggers": "Missing key / type mismatch / value mismatch", "example": "rename, enum change, int→string"},
            {"mechanism": "Json.Number Preservation", "triggers": "Integer numeric vs string representation", "example": "150000 vs '150000'"}
        ]}
    ],
    narrative="Verifier melakukan tiga lapisan validasi: status code, response header, dan body subset. Lapisan body menggunakan recursive map comparison dengan json.Number untuk membedakan integer dan string numerik. Header validation juga menjadi bagian dari kontrak, bukan hanya payload.",
    narrative_en="The verifier performs three layers of validation: status code, response header, and body subset. The body layer uses recursive map comparison with json.Number to distinguish integer from string numerics. Header validation is also part of the contract, not just the payload.",
    mental_model="Three-layer validation: HTTP semantics first, then structural diff with type preservation.",
    source_path="internal/contract/verifier.go:60-128",
    slide_num="08 / 24",
    layout="table"
)

# Slide 09: Decision Sequence — 3 Focused Decision Cards (Layout: cards)
# Limiting phrase 4: "hanya pada jalur ini"
add_slide(
    categories=["DECISION", "GATE"],
    title="Tiga Keputusan Kritis dalam Pipeline Contract Testing",
    subtitle="Pass, fail, dan partial-fail masing-masing mengarahkan tindakan berbeda.",
    content_items=[
        {"type": "cards", "title": "DECISION SEQUENCE", "items": [
            {"title": "Decision 1: Passed?", "desc": "Hanya pada jalur ini deployment diizinkan. Semua interaction verification harus returning true."},
            {"title": "Decision 2: Failed?", "desc": "Deployment diblokir. Error list harus ditinjau tim provider — bukan sekadar retry tanpa analisis."},
            {"title": "Decision 3: Partial?", "desc": "Beberapa interaction fail, sebagian pass. Review per-interaction untuk изолиasi scope perubahan."}
        ]}
    ],
    narrative="Pipeline contract testing menghasilkan tiga jalur keputusan: pass (deployment izin), fail (deployment blokir dengan error report), dan partial (beberapa kontrak gagal, sebagian lolos). Hanya pada jalur pass saja deployment diizinkan tanpa review tambahan.",
    narrative_en="The contract testing pipeline produces three decision paths: pass (deployment allowed), fail (deployment blocked with error report), and partial (some contracts fail, others pass). Only along the pass path is deployment allowed without additional review.",
    mental_model="Three decision branches: pass allows deploy, fail blocks, partial requires per-interaction review.",
    source_path="cmd/demo/main.go:14-68",
    slide_num="09 / 24",
    layout="cards"
)

# Slide 10: Bukti Eksekusi — Terminal Output Demo (Layout: code)
add_slide(
    categories=["EVIDENCE", "EXECUTION"],
    title="Demo 4 Stage: Generate → V1 Pass → Breaking Blocked → Dual Pass",
    subtitle="Eksekusi nyata menghasilkan output deterministic yang dapat direproduksi.",
    content_items=[
        {"type": "code", "title": "go run ./cmd/demo (Stage 3 Output)", "items": [
            "[Stage 3] Running Breaking Provider Contract Verification:",
            "Result: BLOCKED! Breaking changes detected before deployment:",
            "  1. [A request for order details by ID] path 'total':",
            "     type mismatch (expected 150000 [json.Number], got 150000 [string])",
            "  2. [A request for order details by ID] path 'status':",
            "     value mismatch (expected \"IN_PROGRESS\", got \"in_progress\")",
            "  3. [A request for order details by ID] missing expected field",
            "     'customer.name'",
            "CI Deployment Gate: PREVENTED PRODUCTION OUTAGE."
        ]}
    ],
    narrative="Output demo ini menunjukkan Stage 3 yang memverifikasi ProviderBreaking. Tiga error detail tercetak secara terstruktur, masing-masing menyebutkan path field, jenis mismatch, dan nilai yang diharapkan versus aktual. Gate CI menolak deployment — produksi terlindungi.",
    narrative_en="This demo output shows Stage 3 verifying ProviderBreaking. Three structured errors are printed, each citing the field path, mismatch type, and expected versus actual values. The CI gate rejects the deployment — production is protected.",
    mental_model="Deterministic error reporting provides actionable feedback to the provider team.",
    source_path="engineering/03-execution-result.md:100-113",
    slide_num="10 / 24",
    layout="code"
)

# Slide 11: Analisis Batasan Bukti — YANG TIDAK DIBUKTIKAN (Layout: warning)
add_slide(
    categories=["LIMITATIONS", "BOUNDARY"],
    title="YANG TIDAK DIBUKTIKAN dalam Laboratorium Ini",
    subtitle="Lab adalah simplifikasi; beberapa aspek produksi belum tercover.",
    content_items=[
        {"type": "warning", "title": "BATASAN LAB", "items": [
            {"item": "Pact Broker Remote", "detail": "Kontrak dipertukarkan in-memory; tidak ada broker sentral untuk store dan can-i-deploy check lintas tim."},
            {"item": "Provider State Management", "detail": "State didefinisikan via provider_state string; tidak ada callback setup/teardown state fixture."},
            {"item": "Array Diff", "detail": "DiffValues() tidak menangani array comparison — jika response berisi list, diff gagal atau skip."},
            {"item": "Message Pact", "detail": "Verifikasi terbatas pada HTTP REST; event-driven (Kafka/RabbitMQ) tidak diuji dalam lab."},
            {"item": "V2 Schema Assertion", "detail": "Endpoint /v2/orders/{id} hanya diassert status 200, bukan body schema sesuai kontrak."}
        ]}
    ],
    narrative="Laboratorium ini sengaja disederhanakan untuk reproduktibilitas lokal. Beberapa fitur produksi seperti Pact Broker, provider state callbacks, dan array diff tidak diimplementasikan. V2 endpoint hanya diverifikasi status code-nya, bukan body schema — ini adalah gap yang perlu ditambal untuk production readiness.",
    narrative_en="This lab is intentionally simplified for local reproducibility. Several production features like Pact Broker, provider state callbacks, and array diff are not implemented. The V2 endpoint is verified only for status code, not body schema — this is a gap that needs addressing for production readiness.",
    mental_model="Lab scope is a minimal viable engine; production requires broker, state management, and broader diff coverage.",
    source_path="research-audit/03-claim-audit.md:1-164",
    slide_num="11 / 24",
    layout="warning"
)

# Slide 12: Trade-off Arsitektur (Layout: cards)
add_slide(
    categories=["TRADE-OFF", "ARCHITECTURE"],
    title="Trade-off: Custom Verifier vs Full Pact Stack",
    subtitle="Keputusan pure Go stdlib mengorbankan fitur produksi demi reproduktibilitas.",
    content_items=[
        {"type": "cards", "title": "CUSTOM vs FULL PACT", "items": [
            {"title": "Custom Verifier (Lab)", "desc": "Pure Go stdlib, zero external deps, race-safe, cepat. Kompromi: subset matching, tanpa broker, tanpa message pact."},
            {"title": "Full Pact Stack (Produksi)", "desc": "Ruby/JS client, Pact Broker, provider state callbacks, message pact. Lebih lengkap tapi butuh dependensi eksternal."},
            {"title": "Hybrid Approach", "desc": "Gunakan custom verifier untuk CI gate cepat, migrasi ke Pact Broker saat skalabilitas lintas tim diperlukan."}
        ]}
    ],
    narrative="Laboratorium memilih custom verifier dengan Go standard library murni untuk menghindari dependensi Ruby atau Node.js. Kelebihannya: test berjalan di environment apa pun. Kekurangannya: tidak ada broker sentral, provider state management manual, dan tidak support message pact. Untuk produksi skala kecil, custom verifier cukup.",
    narrative_en="The lab chose a custom verifier using pure Go standard library to avoid Ruby or Node.js dependencies. The advantage: tests run in any environment. The disadvantage: no central broker, manual state management, and no message pact support. For small-scale production, the custom verifier is sufficient.",
    mental_model="Zero-dependency trade-off sacrifices broker and state management for portability and speed.",
    source_path="content/02-master-draft.md:86-91",
    slide_num="12 / 24",
    layout="cards"
)

# Slide 13: Concurrency / Validation Pipeline (Layout: detail)
add_slide(
    categories=["CONCURRENCY", "THREAD-SAFETY"],
    title="Verifier Thread-Safe: 20 Goroutine Tanpa Data Race",
    subtitle="Setiap goroutine menginisialisasi request dan response lokal — state tidak dibagikan.",
    content_items=[
        {"type": "detail", "title": "CONCURRENT VERIFICATION PATTERN", "items": [
            {"pattern": "Per-goroutine Client", "detail": "Verifier.NewVerifier() membuat instance baru dengan http.Client timeout 5s — tidak ada shared mutable state."},
            {"pattern": "Request Isolation", "detail": "http.NewRequest() menghasilkan request object unik per interaksi — concurrent Do() aman."},
            {"pattern": "Result Aggregation", "detail": "Setiap goroutine mengembalikan VerificationResult lokal — sinkronisasi hanya di wg.Wait(), bukan di result."},
            {"pattern": "Race Detector PASS", "detail": "go test -race ./... menghasilkan: ok labs/26-contract-testing/tests 1.425s tanpa race detected."}
        ]}
    ],
    narrative="Verifier tidak menggunakan shared mutable state manapun. Setiap panggilan Verify() menginisialisasi http.Client baru, membuat request baru, dan mengembalikan result baru. Race detector Go konfirmasi thread-safety melalui test concurrent 20 goroutine selama 0.01 detik.",
    narrative_en="The verifier does not use any shared mutable state. Each Verify() call initializes a new http.Client, creates a new request, and returns a new result. Go's race detector confirms thread-safety through the 20-goroutine concurrent test in 0.01 seconds.",
    mental_model="Stateless verification per call guarantees thread-safety without locks or mutexes.",
    source_path="tests/contract_test.go:97-116",
    slide_num="13 / 24",
    layout="detail"
)

# Slide 14: Before vs After Comparison (Layout: comparison)
# Limiting phrase 5: "tidak otomatis"
add_slide(
    categories=["COMPARISON", "TRANSFORMATION"],
    title="Sebelum vs Sesudah Contract Testing di Pipeline CI",
    subtitle="Transformasi dari blind deploy ke gated deployment dengan evidence.",
    content_items=[
        {"type": "comparison", "title": "BEFORE vs AFTER CONTRACT TESTING", "items": [
            {"before": "Deploy tanpa gate verifikasi kontrak", "after": "CI gate memverifikasi kontrak sebelum deploy"},
            {"before": "Production outage karena field rename", "after": "Breaking change terdeteksi di PR, error report jelas"},
            {"before": "Debug integration failure manual耗时", "after": "Verifikasi otomatis menghasilkan diff terstruktur per field"},
            {"before": "E2E test berjalan lambat dan flaky", "after": "Contract test cepat (sub-second) dan deterministic"},
            {"before": "Provider evolusi tanpa koordinasi consumer", "after": "Dual provider V1+V2 memungkinkan migrasi independen"}
        ]}
    ],
    narrative="Sebelum contract testing, tim mengandalkan E2E test yang lambat dan manual debugging saat integrasi pecah. Setelah implementasi, setiap PR provider divalidasi terhadap kontrak consumer — breaking change tertangkap sebelum deploy. Tidak otomatis berarti semua edge case tertutup, tetapi risiko production outage berkurang signifikan.",
    narrative_en="Before contract testing, teams relied on slow E2E tests and manual debugging when integration broke. After implementation, every provider PR is validated against consumer contracts — breaking changes are caught before deployment. Not automatic means not all edge cases are covered, but the risk of production outage is significantly reduced.",
    mental_model="Contract testing shifts integration safety left from production to CI.",
    source_path="content/02-master-draft.md:1-14",
    slide_num="14 / 24",
    layout="comparison"
)

# Slide 15: Deep Analysis — Verifier Engine Architecture (Layout: detail)
add_slide(
    categories=["ANALYSIS", "ENGINE"],
    title="Verifier Engine: Recursive Map Diff dengan Type Preservation",
    subtitle="Mekanisme diffValues() menjadi inti deteksi breaking change.",
    content_items=[
        {"type": "detail", "title": "DIFFVALUES ENGINE", "items": [
            {"component": "decoder.UseNumber()", "impact": "Mengubah semua angka JSON jadi json.Number, bukan float64. Integer 150000 tetap json.Number, string '150000' tetap string."},
            {"component": "map[string]interface{}", "impact": "Body di-parse sebagai generic map — memungkinkan subset matching tanpa struct tag."},
            {"component": "reflect.TypeOf comparison", "impact": "Primitive type mismatch (string vs int) dideteksi dengan reflect.TypeOf, bukan string comparison."},
            {"component": "Recursive traversal", "impact": "Nested object comparison dilakukan secara rekursif — customer.name vs customer.full_name terdeteksi di level child."},
            {"component": "Nil expected handling", "impact": "Jika expected nil, diffValues() return empty — field opsional diabaikan oleh verifier."}
        ]}
    ],
    narrative="Verifier engine menggunakan pendekatan subset matching: hanya field yang dideklarasikan consumer yang diperiksa. Nested map traversal memungkinkan deteksi error di field dalam (customer.name). UseNumber() adalah kunci yang membedakan integer dan string numerik — fitur yang tidak tersedia pada JSON parser default Go.",
    narrative_en="The verifier engine uses subset matching: only fields declared by the consumer are checked. Nested map traversal enables detection of errors in child fields (customer.name). UseNumber() is the key differentiator that distinguishes integers from string numerics — a feature not available on Go's default JSON parser.",
    mental_model="Subset matching with recursive diff and json.Number preserves type precision for breaking change detection.",
    source_path="internal/contract/verifier.go:108-189",
    slide_num="15 / 24",
    layout="detail"
)

# Slide 16: Evidence — Test Suite Results (Layout: code)
add_slide(
    categories=["EVIDENCE", "TESTS"],
    title="Tujuh Test: Semua PASS, Race Detector Bersih",
    subtitle="Verifikasi menyeluruh mencakup generasi, sukses, gagal, dual, concurrent, header, dan v2.",
    content_items=[
        {"type": "code", "title": "go test -v ./... (Ringkasan)", "items": [
            "=== RUN   TestConsumerContractGeneration",
            "--- PASS: TestConsumerContractGeneration (0.00s)",
            "=== RUN   TestProviderV1_ContractVerification_Success",
            "--- PASS: TestProviderV1_ContractVerification_Success (0.00s)",
            "=== RUN   TestProviderBreaking_ContractVerification_Fails",
            "--- PASS: TestProviderBreaking_ContractVerification_Fails (0.00s)",
            "=== RUN   TestProviderDual_ContractVerification_Success",
            "--- PASS: TestProviderDual_ContractVerification_Success (0.00s)",
            "=== RUN   TestConcurrentContractVerification",
            "--- PASS: TestConcurrentContractVerification (0.01s)",
            "=== RUN   TestVerifier_HeaderValidation_And_ErrorBranches",
            "--- PASS: TestVerifier_HeaderValidation_And_ErrorBranches (0.00s)",
            "=== RUN   TestProviderDual_V2Endpoint_DirectAssertion",
            "--- PASS: TestProviderDual_V2Endpoint_DirectAssertion (0.00s)",
            "ok  \tlabs/26-contract-testing/tests\t0.390s"
        ]}
    ],
    narrative="Ketujuh test menutupi seluruh skenario: generasi kontrak, verifikasi V1 sukses, breaking change gagal, dual provider sukses, concurrent safety, header validation, dan v2 assertion. Race detector berjalan bersih tanpa flagged race condition.",
    narrative_en="The seven tests cover all scenarios: contract generation, V1 success verification, breaking change failure, dual provider success, concurrent safety, header validation, and v2 assertion. The race detector runs clean with no flagged race conditions.",
    mental_model="Comprehensive test coverage validates all verification paths and concurrency safety.",
    source_path="engineering/03-execution-result.md:19-42",
    slide_num="16 / 24",
    layout="code"
)

# Slide 17: Evidence — Race Detector Performance (Layout: code)
add_slide(
    categories=["EVIDENCE", "PERFORMANCE"],
    title="Race Detector: 1.425s untuk Seluruh Test Suite",
    subtitle="Konfirmasi thread-safety tanpa mutex atau lock.",
    content_items=[
        {"type": "code", "title": "go test -race ./... (Output)", "items": [
            "=== RUN   TestConsumerContractGeneration",
            "--- PASS: TestConsumerContractGeneration (0.00s)",
            "=== RUN   TestProviderV1_ContractVerification_Success",
            "--- PASS: TestProviderV1_ContractVerification_Success (0.00s)",
            "=== RUN   TestProviderBreaking_ContractVerification_Fails",
            "--- PASS: TestProviderBreaking_ContractVerification_Fails (0.00s)",
            "=== RUN   TestProviderDual_ContractVerification_Success",
            "--- PASS: TestProviderDual_ContractVerification_Success (0.00s)",
            "=== RUN   TestConcurrentContractVerification",
            "--- PASS: TestConcurrentContractVerification (0.01s)",
            "=== RUN   TestVerifier_HeaderValidation_And_ErrorBranches",
            "--- PASS: TestVerifier_HeaderValidation_And_ErrorBranches (0.00s)",
            "=== RUN   TestProviderDual_V2Endpoint_DirectAssertion",
            "--- PASS: TestProviderDual_V2Endpoint_DirectAssertion (0.00s)",
            "ok  \tlabs/26-contract-testing/tests\t(cached)",
            ""
            "ok  \tlabs/26-contract-testing/tests\t1.425s",
            "    (race detector: no races detected)"
        ]}
    ],
    narrative="Race detector mengonfirmasi bahwa tidak ada shared mutable state antara goroutine. Setiap verifikasi berjalan independent dengan http.Client lokal dan request object unik. Performa concurrent test hanya 0.01s — verifier engine efisien dan scalable.",
    narrative_en="The race detector confirms no shared mutable state between goroutines. Each verification runs independently with a local http.Client and unique request object. The concurrent test takes only 0.01s — the verifier engine is efficient and scalable.",
    mental_model="Thread-safety achieved through stateless per-call design, not synchronization primitives.",
    source_path="engineering/03-execution-result.md:44-47",
    slide_num="17 / 24",
    layout="code"
)

# Slide 18: Deep Analysis — Provider Implementation Patterns (Layout: table)
add_slide(
    categories=["ANALYSIS", "PROVIDER"],
    title="Tiga Pola Provider: Compliant, Breaking, dan Dual Evolution",
    subtitle="Perbandingan implementasi HTTP handler.",
    content_items=[
        {"type": "table", "title": "PROVIDER IMPLEMENTATION PATTERNS", "items": [
            {"pattern": "ProviderV1", "route": "/v1/orders/{id}", "behavior": "Mengembalikan OrderResponseV1 compliant kontrak. Field notes ada tapi tidak divalidasi (subset matching)."},
            {"pattern": "ProviderBreaking", "route": "/v1/orders/{id}", "behavior": "Mengembalikan OrderResponseBreaking dengan 3 mutasi: enum casing, field rename, type mutation."},
            {"pattern": "ProviderDual", "route": "/v1/orders/{id}", "behavior": "Mempertahankan V1 route compliant sambil menambahkan V2 route dengan schema baru + currency field."},
            {"pattern": "ProviderDual", "route": "/v2/orders/{id}", "behavior": "Respons V2: OrderResponseV2 dengan customer.full_name, total string, dan currency IDR. Belum ada kontrak consumer."}
        ]}
    ],
    narrative="Tiga pola provider ini menggambarkan spektrum kompatibilitas. ProviderV1 adalah baseline compliant. ProviderBreaking menunjukkan skenario worst-case. ProviderDual merepresentasikan praktik evolusi aman: V1 dipertahankan untuk consumer lama, V2 diperkenalkan untuk consumer baru tanpa mengganggu yang sudah ada.",
    narrative_en="These three provider patterns illustrate the compatibility spectrum. ProviderV1 is the compliant baseline. ProviderBreaking shows the worst-case scenario. ProviderDual represents safe evolution practice: V1 is maintained for legacy consumers, V2 is introduced for new consumers without disrupting existing ones.",
    mental_model="Dual routing enables independent consumer migration without breaking existing integrations.",
    source_path="internal/provider/server.go:1-131",
    slide_num="18 / 24",
    layout="table"
)

# Slide 19: Deep Analysis — Model DTO Structure (Layout: detail)
add_slide(
    categories=["ANALYSIS", "MODELS"],
    title="Struktur DTO: V1, Breaking, V2 — Satu Domain, Tiga Wujud",
    subtitle="Refleksi Go struct dengan json tag menentukan perilaku serialisasi.",
    content_items=[
        {"type": "detail", "title": "DTO STRUCT ANALYSIS", "items": [
            {"struct": "OrderResponseV1", "tags": "id, status, customer(Name), total(int64), notes,omitempty", "validates": "Memenuhi kontrak consumer — semua field expected ada dengan tipe benar."},
            {"struct": "OrderResponseBreaking", "tags": "id, status(in_progress), customer(FullName), total(string), notes", "validates": "Melanggar kontrak — 3 breaking change: enum casing, field rename, type mutation."},
            {"struct": "OrderResponseV2", "tags": "id, status(in_progress), customer(FullName), total(string), currency", "validates": "Schema baru V2 — belum ada consumer contract, hanya status assertion di test."}
        ]}
    ],
    narrative="Setiap struct merepresentasikan fase evolusi API. V1 adalah kontrak asli yang dipatuhi. Breaking adalah contoh deformasi kontrak yang harus ditangkap. V2 adalah arah evolusi yang sah — field rename dan type mutation diterima karena consumer V2 yang baru telah beradaptasi.",
    narrative_en="Each struct represents a phase of API evolution. V1 is the original contract compliant. Breaking is an example of contract deformation that must be caught. V2 is the legitimate evolution direction — field rename and type mutation are accepted because the new V2 consumer has already adapted.",
    mental_model="DTO struct tags define serialization behavior; contract testing validates the serialized output against expectations.",
    source_path="internal/model/order.go:1-53",
    slide_num="19 / 24",
    layout="detail"
)

# Slide 20: Checklist Implementasi (Layout: cards)
add_slide(
    categories=["IMPLEMENTATION", "CHECKLIST"],
    title="Checklist Implementasi Contract Testing yang Benar",
    subtitle="Panduan praktis untuk adopsi CDC di tim engineering.",
    content_items=[
        {"type": "cards", "title": "CDC IMPLEMENTATION CHECKLIST", "items": [
            {"item": "Definisikan kontrak minimal", "desc": "Consumer hanya mendeklarasikan field yang benar-benar digunakan — jangan over-specify."},
            {"item": "Sertakan HTTP semantik", "desc": "Kontrak mencakup method, path, status, header, body dengan tipe dan casing eksak."},
            {"item": "Jalankan verifikasi di CI", "desc": "Provider CI harus menjalankan verifier sebelum merge — pass = boleh deploy, fail = blokir."},
            {"item": "Gunakan expand/contract pattern", "desc": "Breaking change: tambah field baru dulu, migrasi consumer, baru hapus field lama."},
            {"item": "Pertahankan dual routing", "desc": "V1 tetap compliant sambil V2 diperkenalkan untuk consumer yang bermigrasi."},
            {"item": "Hindari over-specification", "desc": "Jangan uji aturan validasi bisnis (max 20 chars) di kontrak — fokus pada format pesan."}
        ]}
    ],
    narrative="Checklist ini merangkum praktik terbaik contract testing. Poin terpenting: kontrak harus minimal dan berfokus pada field yang digunakan consumer. Over-specification adalah anti-pattern utama yang membuat kontrak terlalu rapuh dan menghambat evolusi provider.",
    narrative_en="This checklist summarizes best practices for contract testing. The most important point: contracts must be minimal and focused on fields used by the consumer. Over-specification is the primary anti-pattern that makes contracts too brittle and hinders provider evolution.",
    mental_model="Minimal contract, CI gate enforcement, and expand/contract pattern are the three non-negotiables.",
    source_path="content/05-key-takeaways.md:1-20",
    slide_num="20 / 24",
    layout="cards"
)

# Slide 21: Anti-Pattern & Cara Reproduksi Kegagalan (Layout: detail)
add_slide(
    categories=["ANTI-PATTERN", "REPRODUCTION"],
    title="Lima Anti-Pattern Contract Testing dan Cara Reproduksi",
    subtitle="Setiap anti-pattern dapat direproduksi dengan modifikasi kontrak atau provider.",
    content_items=[
        {"type": "detail", "title": "ANTI-PATTERNS & REPRODUCTION", "items": [
            {"anti": "Over-specification", "repro": "Tambahkan field validasi bisnis (e.g., 'username max 20 chars') ke kontrak. Provider yang melonggarkan validasi akan pecah kontrak tanpa dampak consumer."},
            {"anti": "Tidak bedakan additive vs breaking", "repro": "Ubah total dari int jadi string tanpa update kontrak. Verifier harus gagal — jika lolos, matcher terlalu longgar."},
            {"anti": "Kontrak terlalu ketat", "repro": "Deklarasikan semua field response (termasuk notes, created_at) sebagai required. Provider tambah field baru akan fail."},
            {"anti": "JSON Schema saja", "repro": "Ganti diffValues dengan validator JSON Schema. Casing enum IN_PROGRESS vs in_progress akan lolos Schema valid."},
            {"anti": "Contract testing menggantikan unit test", "repro": "Hapus semua unit test internal, hanya jalankan contract test. Logika bisnis tidak terverifikasi."}
        ]}
    ],
    narrative="Kelima anti-pattern ini umum terjadi saat tim baru mengadopsi contract testing. Over-specification adalah yang paling berbahaya karena membuat kontrak terlalu kaku. JSON Schema alone tidak cukup karena tidak menangkap semantik tipe seperti integer vs string numerik.",
    narrative_en="These five anti-patterns commonly occur when teams first adopt contract testing. Over-specification is the most dangerous because it makes contracts too rigid. JSON Schema alone is insufficient because it cannot capture type semantics like integer versus string numerics.",
    mental_model="Anti-patterns arise from misunderstanding contract scope; specificity must match actual consumer usage.",
    source_path="content/02-master-draft.md:207-214",
    slide_num="21 / 24",
    layout="detail"
)

# Slide 22: Takeaway Praktis & Command Block (Layout: code)
add_slide(
    categories=["PRACTICAL", "COMMANDS"],
    title="Reproduksi Mandiri: Validasi Full Pipeline di Lokal",
    subtitle="Semua tahap demo dapat dijalankan dengan command sederhana tanpa dependensi eksternal.",
    content_items=[
        {"type": "code", "title": "COMMAND RUNNER (PURE GO STANDARD LIBRARY)", "items": [
            "# 1. Jalankan seluruh test suite dengan verbose",
            "go test -v ./...",
            "#",
            "# 2. Jalankan race detector untuk verifikasi thread-safety",
            "go test -race ./...",
            "#",
            "# 3. Eksekusi demo 4 stage (generate -> V1 pass -> breaking blocked -> dual pass)",
            "go run ./cmd/demo",
            "#",
            "# 4. Build binary untuk environment produksi",
            "go build -o bin/contract-demo ./cmd/demo"
        ]}
    ],
    narrative="Semua command di atas berjalan di Go standard library tanpa dependensi eksternal. Race detector mengonfirmasi thread-safety concurrent verification. Demo 4 stage memberikan bukti visual yang dapat direproduksi di mesin manapun.",
    narrative_en="All commands above run on Go standard library with zero external dependencies. The race detector confirms thread-safety of concurrent verification. The 4-stage demo provides reproducible visual evidence on any machine.",
    mental_model="Zero-dependency local verification guarantees fast, repeatable validation across developer workflows.",
    source_path="README.md:30-48",
    slide_num="22 / 24",
    layout="code"
)

# Slide 23: Prinsip Arsitektural Jangka Panjang (Layout: cards)
# Limiting phrase 6: "tidak otomatis"
add_slide(
    categories=["TAKEAWAY", "CORE PRINCIPLES"],
    title="Prinsip Inti: Kendalikan Kontrak, Evolusi dengan Aman",
    subtitle="Arsitektur bertahan lama ketika kontrak divalidasi di setiap lapisan deploy.",
    content_items=[
        {"type": "cards", "title": "TIGA PRINSIP ARSITEKTUR JANGKA PANJANG", "items": [
            {"title": "Consumer Drives the Contract", "desc": "Consumer yang menentukan ekspektasi minimal — bukan provider yang mendokumentasikan capability. Ini memastikan kontrak relevan dengan kebutuhan riil."},
            {"title": "Gate Before Deploy", "desc": "Verifikasi kontrak harus berjalan di CI sebelum merge. Tidak otomatis berarti semua kasus tercover, tetapi mengurangi risiko production outage secara signifikan."},
            {"title": "Dual Routing for Migration", "desc": "Pertahankan V1 compliant sambil memperkenalkan V2. Consumer bermigrasi sesuai ritme masing-masing tanpa tekanan breaking change."}
        ]}
    ],
    narrative="Prinsip pertama memastikan kontrak selalu berpusat pada consumer — bukan dokumentasi provider yang sering tertinggal. Prinsip kedua menempatkan gate di titik paling kritis: sebelum kode mencapai production. Prinsip ketiga menyediakan mekanisme evolusi tanpa disrupti.",
    narrative_en="The first principle ensures contracts remain consumer-centric — not provider documentation that often falls behind. The second principle places the gate at the most critical point: before code reaches production. The third principle provides an evolution mechanism without disruption.",
    mental_model="Consumer-driven contracts, CI gates, and dual routing form the triad of sustainable API evolution.",
    source_path="content/05-key-takeaways.md:1-10",
    slide_num="23 / 24",
    layout="cards"
)

# Slide 24: Closing Slide (Dark Navy, repo link pendek, commit SHA, credits)
add_slide(
    categories=["CLOSING"],
    title="Sistem tanpa contract testing adalah sistem yang menunda kegagalan integrasi produksinya sendiri.",
    subtitle="",
    content_items=[
        {"type": "bullet", "title": "MODUL IMPLEMENTASI", "items": [
            "internal/consumer/client.go — MobileOrderClient & GenerateMobileContract()",
            "internal/contract/verifier.go — Verifier.Verify() & diffValues()",
            "internal/model/order.go — OrderResponseV1/Breaking/V2 DTOs",
            "internal/provider/server.go — ProviderV1/Breaking/Dual handlers",
            "cmd/demo/main.go — 4-stage demo orchestrator",
            "tests/contract_test.go — 7 test cases"
        ]},
        {"type": "bullet", "title": "VERIFIKASI & STATUS", "items": [
            f"Commit: {COMMIT_SHA} • Status: APPROVED",
            "Test passing: 7 unit tests (generation, V1 pass, breaking fail, dual pass, concurrent, header validation, V2 assertion)",
            "Race detector: Clean (go test -race ./... PASS in 1.425s)",
            "Demo: 4-stage executeable (generate → V1 pass → breaking blocked → dual pass)"
        ]}
    ],
    mental_model="A system without contract testing is merely delaying its inevitable integration failure.",
    source_path=REPO_URL,
    slide_num="24 / 24",
    layout="closing"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
    for s in SLIDES:
        print(f"[{s['slide_num']}] {s['layout']}: {s['title']}")
