#!/usr/bin/env python3
"""
Publication Builder for Software Engineering Lab 26: Contract Testing (Consumer-Driven Contracts).
Defines metadata, colors, and exactly 24 slides for PPTX, PDF, and cover generation.
"""
import os
import sys

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/26-contract-testing/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/26-contract-testing"
COMMIT_SHA = "e963675c67e7d518e19a07ad930dbdb5be73fff1"
LAB_NO = "26"
SLUG = "contract-testing"
TITLE_ID = "Contract Testing"
SUBTITLE_ID = "Consumer-Driven Contracts & CI/CD Deployment Gates"
TITLE_EN = "Contract Testing"
SUBTITLE_EN = "Consumer-Driven Contracts & CI/CD Deployment Gates"
AUTHOR = "Lukman / lukman-ss"

# Brand Colors
DARK_NAVY = "0b1b2d"
TEAL = "087f76"
CREAM = "f4f1e8"
WHITE = "ffffff"
SLATE = "52606b"
ORANGE = "d97706"

# Slide collection
SLIDES = []

def add_slide(categories, title, subtitle, content_items, mental_model, source_path="", slide_num=""):
    """Append a structured slide definition."""
    SLIDES.append({
        "categories": categories,
        "title": title,
        "subtitle": subtitle,
        "content_items": content_items,
        "mental_model": mental_model,
        "source": source_path,
        "slide_num": slide_num
    })

# =====================================================================
# EXACTLY 24 SLIDES DEFINITION
# =====================================================================

# Slide 1: Cover (Dark Navy, Retro Terminal)
add_slide(
    categories=["KONTEKS"],
    title=TITLE_ID,
    subtitle="SOFTWARE ENGINEERING LAB | LAB " + LAB_NO,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "Retro CRT monitor beige dengan scanlines",
            "Terminal monitor CDC verifier engine output",
            f"Penulis: {AUTHOR}",
        ]}
    ],
    mental_model="",
    source_path="",
    slide_num=""
)

# Slide 2: Motivation & Context
add_slide(
    categories=["MOTIVASI", "KONTEKS"],
    title="Masalah Integrasi Microservices: Rapuhnya E2E Test",
    subtitle="End-to-end testing lambat, flaky, dan mahal untuk mendeteksi breaking schema changes",
    content_items=[
        {"title": "MASALAH E2E TESTING", "items": [
            "Membutuhkan seluruh ekosistem service dan database aktif bersamaan di staging",
            "Flaky network, timeout, dan latency test lingkungan memperlambat pipeline CI/CD",
            "Feedback loop lambat: pengembang menunggu berjam-jam hanya untuk cek kompatibilitas DTO",
        ], "type": "bullet"},
        {"title": "PARADOKS MOCK SERVICE", "items": [
            "Unit test consumer dengan mock lokal selalu lolos karena asumsi schema statis",
            "Saat provider mengubah casing enum atau nama field di produksi, consumer langsung crash",
            "Mock yang dirawat manual cepat usang dan menciptakan ilusi keamanan palsu",
        ], "type": "bullet"},
        {"title": "KONDISI DI LAB 26", "items": [
            "Consumer Mobile App membutuhkan kontrak schema minimal dari OrderService",
            "Provider melakukan perubahan tanpa koordinasi yang melanggar kontrak DTO",
            "Dibutuhkan deployment gate otomatis berbasis CDC tanpa dependensi external runtime",
        ], "type": "bullet"},
    ],
    mental_model="Mocks drift into false confidence; contract tests turn consumer assumptions into executable provider guardrails.",
    source_path="labs/26-contract-testing/README.md",
    slide_num="2 / 24"
)

# Slide 3: Consumer-Driven Contracts Anatomy
add_slide(
    categories=["BUKTI LOGIKA", "ANATOMI CDC"],
    title="Anatomi Consumer-Driven Contracts: Consumer Menentukan Ekspektasi",
    subtitle="Bukan provider yang mendikte payload, melainkan kebutuhan nyata consumer yang didaftarkan",
    content_items=[
        {"title": "PRINSIP KEBUTUHAN MINIMAL", "items": [
            "Consumer hanya mencatat subset field yang benar-benar dibaca oleh kodenya",
            "Field provider yang tidak dikonsumsi (misal internal notes) diabaikan dalam kontrak",
            "Provider bebas berevolusi pada field yang tidak mengikat consumer mana pun",
        ], "type": "bullet"},
        {"title": "STRUKTUR INTERACTION", "items": [
            "Request: HTTP Method GET /v1/orders/ORD-123",
            "Provider State: Order ORD-123 exists and is IN_PROGRESS",
            "Response: Status 200 OK, Content-Type JSON, body dengan id, status, customer.name, total",
        ], "type": "code"},
        {"title": "PERGESERAN PARADIGMA", "items": [
            "Provider-driven: Provider merilis OpenAPI spec raksasa, consumer dipaksa mengikuti",
            "Consumer-driven: Consumer merilis kontrak kebutuhan, provider membuktikan kepatuhan",
            "Hasil: Hubungan dependensi API menjadi eksplisit, terukur, dan dapat diaudit di CI",
        ], "type": "bullet"},
    ],
    mental_model="The consumer defines the contract, but the provider proves compliance before shipping.",
    source_path="labs/26-contract-testing/internal/consumer/client.go:84",
    slide_num="3 / 24"
)

# Slide 4: Real-World Outage Scenario
add_slide(
    categories=["MASALAH OPERASIONAL", "SKENARIO"],
    title="Skenario Nyata: Tiga Mutasi Payload Pemicu Downtime",
    subtitle="Perubahan sederhana di backend merusak parsing JSON aplikasi mobile",
    content_items=[
        {"title": "MUTASI 1: CASING ENUM STATUS", "items": [
            "Backend mengubah 'IN_PROGRESS' menjadi huruf kecil 'in_progress' untuk standardisasi",
            "Client mobile memiliki strict enum parser: status tidak dikenali memicu fatal error",
            "User tertahan di layar loading tanpa pesan kesalahan yang informatif",
        ], "type": "bullet"},
        {"title": "MUTASI 2: FIELD RENAME CUSTOMER", "items": [
            "Backend merapikan schema: field customer.name diubah menjadi customer.full_name",
            "Parser consumer mengharapkan customer.name tidak boleh string kosong",
            "Hasil: unmarshal error 'contract violation: customer.name is missing'",
        ], "type": "bullet"},
        {"title": "MUTASI 3: TIPE PRIMITIF MUTATION", "items": [
            "Backend mengirim total dalam bentuk string '150000' daripada integer 150000",
            "Go json.Unmarshal gagal melakukan deserialisasi string ke int64",
            "Kerusakan terjadi seketika di level deserializer sebelum domain logic dijalankan",
        ], "type": "bullet"},
    ],
    mental_model="Small unannounced payload changes cause catastrophic client-side deserialization outages.",
    source_path="labs/26-contract-testing/internal/model/order.go:27",
    slide_num="4 / 24"
)

# Slide 5: The CI/CD Deployment Gate Concept
add_slide(
    categories=["ARSITEKTUR CI/CD", "GATE"],
    title="CI/CD Deployment Gate: Can I Deploy?",
    subtitle="Memindahkan verifikasi kompatibilitas dari staging runtime ke build pipeline",
    content_items=[
        {"title": "PIPELINE CONSUMER", "items": [
            "Menjalankan unit test client dengan mock provider",
            "Membangkitkan file kontrak JSON (MobileApp -> OrderService)",
            "Mempublikasikan kontrak ke artifact store atau git repo",
        ], "type": "bullet"},
        {"title": "PIPELINE PROVIDER", "items": [
            "Menjalankan server provider di lingkungan in-memory ephemeral",
            "Menjalankan CDC Verifier membaca kontrak consumer aktif",
            "Gate check: Jika verification gagal -> CI Provider FAIL -> Deployment BLOCKED",
        ], "type": "bullet"},
        {"title": "KEUNTUNGAN DETERMINISTIK", "items": [
            "Tidak butuh staging environment yang selalu online 24/7",
            "Eksekusi verifikasi selesai dalam milidetik (0.00s pada Lab 26)",
            "Breaking change tertahan sebelum binary masuk ke container registry",
        ], "type": "bullet"},
    ],
    mental_model="Verify consumer contracts at build time so broken payloads never reach the cluster.",
    source_path="labs/26-contract-testing/engineering/01-design.md",
    slide_num="5 / 24"
)

# Slide 6: Lab Architecture & Module Layout
add_slide(
    categories=["DESAIN SISTEM", "MODUL"],
    title="Arsitektur Lab 26: Desain Sistem Tanpa Dependency Daemon",
    subtitle="Implementasi CDC murni Go standard library tanpa binary external pact-ruby",
    content_items=[
        {"title": "STRUKTUR MODUL", "items": [
            "internal/model: Definisi DTO OrderResponseV1, Breaking, V2",
            "internal/consumer: MobileOrderClient & GenerateMobileContract()",
            "internal/provider: ProviderV1, ProviderBreaking, ProviderDual handlers",
            "internal/contract: Verifier engine dengan recursive schema diffing",
        ], "type": "bullet"},
        {"title": "ALASAN DESAIN PURE GO", "items": [
            "Menghilangkan dependensi cgo, daemon broker eksternal, atau CLI binary berat",
            "Dapat dijalankan instan via go test -v ./... dan go run ./cmd/demo",
            "Membuktikan semantik CDC secara murni, transparan, dan mudah dipahami",
        ], "type": "bullet"},
        {"title": "KEAMANAN KONKURENSI", "items": [
            "Verifier stateless menggunakan http.Client Go thread-safe",
            "Uji konkurensi 20 goroutine paralel lolos uji race detector tanpa data race",
            "Dapat diverifikasi ulang dengan perintah go test -race ./...",
        ], "type": "bullet"},
    ],
    mental_model="Pure standard-library architecture isolates contract mechanics from tooling friction.",
    source_path="labs/26-contract-testing/internal/contract/verifier.go",
    slide_num="6 / 24"
)

# Slide 7: Consumer Contract Generation
add_slide(
    categories=["BUKTI KODE", "CONSUMER"],
    title="Consumer Code: Mendefinisikan Kontrak Minimal",
    subtitle="Mobile consumer mendefinisikan ekspektasi HTTP, headers, dan struktur response",
    content_items=[
        {"title": "KODE PEMBUAT KONTRAK", "items": [
            "func GenerateMobileContract() *contract.Contract {",
            "    return &contract.Contract{",
            "        Consumer: 'MobileApp', Provider: 'OrderService',",
            "        Interactions: []contract.Interaction{{",
            "            Request: contract.RequestDefinition{Method: 'GET', Path: '/v1/orders/ORD-123'},",
            "            Response: contract.ResponseDefinition{Status: 200, ...},",
            "        }},",
            "    }",
            "}",
        ], "type": "code"},
        {"title": "PAYLOAD MINIMAL", "items": [
            "id: 'ORD-123'",
            "status: 'IN_PROGRESS'",
            "customer: { name: 'Budi Santoso' }",
            "total: json.Number('150000')",
        ], "type": "code"},
        {"title": "SUBSET FIELD", "items": [
            "Consumer mengabaikan field provider seperti 'notes' atau 'customer.id'",
            "Menggunakan json.Number untuk membedakan integer vs string secara eksak",
            "Mendokumentasikan provider state agar provider dapat menyiapkan data uji",
        ], "type": "bullet"},
    ],
    mental_model="Express only what the client consumes; ignore all ancillary provider attributes.",
    source_path="labs/26-contract-testing/internal/consumer/client.go:85",
    slide_num="7 / 24"
)

# Slide 8: Consumer Client Validation Logic
add_slide(
    categories=["BUKTI KODE", "VALIDASI"],
    title="Mobile Client: Parsing Ketat Terhadap Pelanggaran Kontrak",
    subtitle="Client mobile memvalidasi keberadaan field wajib dan konsistensi enum",
    content_items=[
        {"title": "UNMARSHAL & ASSERTION", "items": [
            "if err := json.Unmarshal(body, &raw); err != nil {",
            "    return nil, fmt.Errorf('failed to parse contract schema: %w', err)",
            "}",
            "if raw.Customer.Name == '' {",
            "    return nil, fmt.Errorf('contract violation: customer.name is missing')",
            "}",
            "if raw.Status != 'IN_PROGRESS' && raw.Status != 'COMPLETED' {",
            "    return nil, fmt.Errorf('contract violation: unknown status %q', raw.Status)",
            "}",
        ], "type": "code"},
        {"title": "PERILAKU SAAT BREACH", "items": [
            "Jika customer.name hilang, client langsung mengembalikan error spesifik",
            "Jika status berhuruf kecil 'in_progress', client menolak memproses pesanan",
            "Mencegah nilai null merembet ke lapisan UI/domain mobile app",
        ], "type": "bullet"},
        {"title": "KEPERLUAN KONTRAK", "items": [
            "Aturan validasi inilah yang wajib tercermin pada kontrak CDC",
            "CDC menjamin provider tidak merilis response yang melanggar aturan ini",
            "Saling melengkapi: client strict, CDC memvalidasi provider sebelum deploy",
        ], "type": "bullet"},
    ],
    mental_model="Strict client validation protects app runtime; CDC guarantees the provider respects those rules.",
    source_path="labs/26-contract-testing/internal/consumer/client.go:68",
    slide_num="8 / 24"
)

# Slide 9: CDC Verifier Engine Implementation
add_slide(
    categories=["BUKTI KODE", "VERIFIER"],
    title="CDC Verifier: Eksekusi Request & Validasi Status Code",
    subtitle="Mengeksekusi HTTP interaction terhadap URL provider yang sedang berjalan",
    content_items=[
        {"title": "ALUR VERIFIKASI", "items": [
            "targetURL := strings.TrimRight(baseURL, '/') + interaction.Request.Path",
            "req, err := http.NewRequest(interaction.Request.Method, targetURL, nil)",
            "resp, err := v.Client.Do(req)",
            "if resp.StatusCode != interaction.Response.Status {",
            "    result.Passed = false",
            "    result.Errors = append(result.Errors, fmt.Sprintf('status code mismatch...'))",
            "}",
        ], "type": "code"},
        {"title": "VALIDASI HEADER", "items": [
            "Verifier memeriksa header Content-Type secara case-insensitive",
            "Mencegah response text/html atau raw string saat consumer meminta JSON",
            "Memastikan response HTTP valid sebelum masuk ke deserialisasi JSON",
        ], "type": "bullet"},
        {"title": "DETERMINISTIK & CEPAT", "items": [
            "Menggunakan ephemeral httptest.Server di dalam proses testing",
            "Tanpa roundtrip jaringan luar — selesai dalam 0.00 detik",
            "Ideal untuk dimasukkan ke dalam setiap pull request dan pre-commit check",
        ], "type": "bullet"},
    ],
    mental_model="Run the interaction against a live ephemeral provider instance and assert HTTP semantics.",
    source_path="labs/26-contract-testing/internal/contract/verifier.go:61",
    slide_num="9 / 24"
)

# Slide 10: Recursive Schema Diffing
add_slide(
    categories=["BUKTI KODE", "DIFFING"],
    title="Algoritma Diffing: Validasi Recursive Subset Schema",
    subtitle="Membandingkan struktur nested JSON hanya pada field yang diminta consumer",
    content_items=[
        {"title": "LOGIKA DIFFING OBJEK", "items": [
            "func diffValues(path string, expected, actual interface{}) []string {",
            "    for key, expVal := range expMap {",
            "        actVal, exists := actMap[key]",
            "        if !exists { diffs = append(diffs, 'missing expected field') }",
            "        diffs = append(diffs, diffValues(childPath, expVal, actVal)...)",
            "    }",
            "}",
        ], "type": "code"},
        {"title": "PENGECEKAN TIPE PRIMITIF", "items": [
            "Membedakan json.Number (angka) dari string JSON secara presisi",
            "Memverifikasi kecocokan tipe via reflect.TypeOf()",
            "Mendeteksi perubahan nilai enum seperti 'IN_PROGRESS' vs 'in_progress'",
        ], "type": "bullet"},
        {"title": "TOLERANSI OVER-PROVISIONING", "items": [
            "Provider memiliki field 'notes' dan 'customer.id' yang tidak ada di kontrak",
            "diffValues mengabaikan field ekstra tersebut tanpa menghasilkan error",
            "Memungkinkan provider menambahkan fitur baru tanpa merusak consumer lama",
        ], "type": "bullet"},
    ],
    mental_model="Validate subset presence and exact types; tolerate unrequested provider fields.",
    source_path="labs/26-contract-testing/internal/contract/verifier.go:130",
    slide_num="10 / 24"
)

# Slide 11: Provider V1 Verification
add_slide(
    categories=["HASIL PENGUJIAN", "PROVIDER V1"],
    title="Verifikasi Provider V1: Memenuhi Seluruh Kontrak",
    subtitle="Server V1 mengirimkan DTO asli yang sesuai dengan harapan MobileApp",
    content_items=[
        {"title": "RESPON PROVIDER V1", "items": [
            "ID: 'ORD-123'",
            "Status: 'IN_PROGRESS' (uppercase string)",
            "Customer: { ID: 'cust-100', Name: 'Budi Santoso' }",
            "Total: 150000 (int64)",
            "Notes: 'Internal note not used by mobile app'",
        ], "type": "code"},
        {"title": "HASIL VERIFIKASI", "items": [
            "Status Code: 200 OK — MATCH",
            "Header Content-Type: application/json — MATCH",
            "Payload schema: 4/4 field consumer ditemukan dengan tipe data yang tepat",
            "Field 'notes' dan 'customer.id' diabaikan secara aman",
        ], "type": "bullet"},
        {"title": "STATUS CI GATE", "items": [
            "Verification result: Passed = true, Errors = 0",
            "CI Deployment Gate: ALLOWED — Siap dideploy ke cluster produksi",
            "Eksekusi test TestProviderV1_ContractVerification_Success: PASS (0.00s)",
        ], "type": "bullet"},
    ],
    mental_model="When expected subset matches completely, deployment gate turns green.",
    source_path="labs/26-contract-testing/internal/provider/server.go:12",
    slide_num="11 / 24"
)

# Slide 12: Breaking Provider Triggering Failures
add_slide(
    categories=["HASIL PENGUJIAN", "BREAKING CHANGES"],
    title="Breaking Provider: Terdeteksi & Diblokir di CI Gate",
    subtitle="Tiga perubahan schema langsung memicu kegagalan verifikasi kontrak",
    content_items=[
        {"title": "OUTPUT VERIFIKASI GAGAL", "items": [
            "[Stage 3] Running Breaking Provider Contract Verification:",
            "Result: BLOCKED! Breaking changes detected before deployment:",
            "  1. path 'total': type mismatch (expected 150000 [json.Number], got 150000 [string])",
            "  2. path 'status': value mismatch (expected 'IN_PROGRESS', got 'in_progress')",
            "  3. missing expected field 'customer.name'",
            "CI Deployment Gate: PREVENTED PRODUCTION OUTAGE.",
        ], "type": "code"},
        {"title": "PENYEBAB KERUSAKAN", "items": [
            "Total: tipe bermutasi dari integer ke string",
            "Status: casing bermutasi dari UPPERCASE ke lowercase",
            "Customer.name: nama field direname menjadi customer.full_name",
        ], "type": "bullet"},
        {"title": "DAMPAK PADA PIPELINE", "items": [
            "Exit code non-nol menghentikan build runner CI/CD",
            "Engineer backend langsung melihat daftar field yang melanggar kontrak",
            "Produksi tetap aman: tidak ada deployment versi cacat ke live server",
        ], "type": "bullet"},
    ],
    mental_model="Contract test failures pinpoint exact field diffs before bad code reaches production.",
    source_path="labs/26-contract-testing/cmd/demo/main.go:40",
    slide_num="12 / 24"
)

# Slide 13: Breaking Changes Taxonomy
add_slide(
    categories=["ANALISIS TEKNIK", "TAKSONOMI"],
    title="Taksonomi Breaking Changes pada REST API",
    subtitle="Klasifikasi perubahan payload yang memicu kegagalan consumer",
    content_items=[
        {"title": "TABEL TAKSONOMI PERUBAHAN", "items": [
            "Tipe Mutasi        | Contoh Perubahan              | Dampak Consumer",
            "-------------------|-------------------------------|------------------------",
            "Field Rename       | name -> full_name             | Field missing / null",
            "Type Mutation      | int 150000 -> string '150000' | Deserialization crash",
            "Enum Value Change  | IN_PROGRESS -> in_progress    | Validation/Logic error",
            "Field Removal      | Menghapus field total         | Missing required field",
            "HTTP Status Change | 200 OK -> 202 Accepted        | Client error branch",
        ], "type": "table"},
        {"title": "PERUBAHAN YANG AMAN (NON-BREAKING)", "items": [
            "Menambahkan field baru di level root atau child (additive changes)",
            "Menambahkan endpoint HTTP baru tanpa mengubah endpoint lama",
            "Memperluas toleransi input request pada provider",
        ], "type": "bullet"},
        {"title": "PRINSIP HYRUM'S LAW", "items": [
            "Dengan jumlah consumer yang cukup, semua perilaku API yang dapat diamati akan menjadi dependensi",
            "CDC mendokumentasikan observasi tersebut secara eksplisit dalam kode",
        ], "type": "bullet"},
    ],
    mental_model="Know what breaks: renames, type mutations, and enum alterations break clients; additions do not.",
    source_path="labs/26-contract-testing/internal/model/order.go",
    slide_num="13 / 24"
)

# Slide 14: Safe API Evolution Strategy
add_slide(
    categories=["SOLUSI ARSITEKTUR", "EVOLUSI API"],
    title="Strategi Evolusi API: Dual DTO & Parallel Routing",
    subtitle="Memperkenalkan fitur baru tanpa merusak kontrak consumer lama",
    content_items=[
        {"title": "POLA DUAL ROUTING", "items": [
            "GET /v1/orders/{id} -> Mengembalikan schema V1 (kontrak MobileApp tetap terpenuhi)",
            "GET /v2/orders/{id} -> Mengembalikan schema V2 (fitur baru: full_name, string total, currency)",
            "Backend melayani kedua endpoint secara simultan dari domain model yang sama",
        ], "type": "bullet"},
        {"title": "KODE PROVIDER DUAL", "items": [
            "if strings.HasPrefix(r.URL.Path, '/v1/orders/') {",
            "    respV1 := model.OrderResponseV1{ ... }",
            "    json.NewEncoder(w).Encode(respV1)",
            "}",
            "if strings.HasPrefix(r.URL.Path, '/v2/orders/') {",
            "    respV2 := model.OrderResponseV2{ ... }",
            "    json.NewEncoder(w).Encode(respV2)",
            "}",
        ], "type": "code"},
        {"title": "HASIL VERIFIKASI", "items": [
            "ProviderDual lolos verifikasi kontrak MobileApp V1 dengan 0 error",
            "Mobile App dapat terus beroperasi tanpa migrasi tergesa-gesa",
            "Web / client baru dapat langsung mengadopsi schema V2",
        ], "type": "bullet"},
    ],
    mental_model="Evolve via parallel endpoints: satisfy existing contracts on V1 while innovating on V2.",
    source_path="labs/26-contract-testing/internal/provider/server.go:82",
    slide_num="14 / 24"
)

# Slide 15: Expand-Contract Migration Pattern
add_slide(
    categories=["PANDUAN OPERASIONAL", "MIGRASI"],
    title="Pola Expand-Contract: Migrasi Tanpa Downtime",
    subtitle="Tiga fase evolusi API untuk mematikan schema lama dengan aman",
    content_items=[
        {"title": "FASE 1: EXPAND", "items": [
            "Provider menambahkan endpoint /v2 atau field baru",
            "Menjaga endpoint /v1 tetap beroperasi melayani kontrak aktif",
            "CI gate memverifikasi kontrak lama tetap 100% compliant",
        ], "type": "bullet"},
        {"title": "FASE 2: MIGRATE CONSUMERS", "items": [
            "Consumer secara bertahap memperbarui kodenya ke schema V2",
            "Setiap consumer merilis kontrak baru (V2 interaction) ke broker/repo",
            "Saat semua consumer telah migrasi, kontrak V1 menjadi tidak memiliki consumer",
        ], "type": "bullet"},
        {"title": "FASE 3: CONTRACT (DECOMMISSION)", "items": [
            "Provider mengecek broker: tidak ada consumer yang terikat pada kontrak V1",
            "Endpoint /v1 dihapus secara aman tanpa risiko breaking change",
            "Evolusi selesai tanpa pernah sekalipun memicu downtime",
        ], "type": "bullet"},
    ],
    mental_model="Expand new capabilities, migrate consumers independently, contract old paths safely.",
    source_path="labs/26-contract-testing/engineering/01-design.md",
    slide_num="15 / 24"
)

# Slide 16: Verification Execution Results
add_slide(
    categories=["BUKTI PENGUJIAN", "EKSEKUSI"],
    title="Hasil Eksekusi Test: Lolos Unit, Integration & Race Detector",
    subtitle="Seluruh rangkaian test lolos 100% tanpa memory leak atau race condition",
    content_items=[
        {"title": "RINGKASAN GO TEST", "items": [
            "=== RUN   TestConsumerContractGeneration -> PASS (0.00s)",
            "=== RUN   TestProviderV1_ContractVerification_Success -> PASS (0.00s)",
            "=== RUN   TestProviderBreaking_ContractVerification_Fails -> PASS (0.00s)",
            "=== RUN   TestProviderDual_ContractVerification_Success -> PASS (0.00s)",
            "=== RUN   TestConcurrentContractVerification -> PASS (0.01s)",
            "=== RUN   TestVerifier_HeaderValidation_And_ErrorBranches -> PASS (0.00s)",
            "=== RUN   TestProviderDual_V2Endpoint_DirectAssertion -> PASS (0.00s)",
            "PASS: ok labs/26-contract-testing/tests 0.390s",
        ], "type": "code"},
        {"title": "VERIFIKASI RACE DETECTOR", "items": [
            "Perintah: go test -race ./...",
            "Hasil: ok labs/26-contract-testing/tests 1.425s",
            "20 goroutine paralel memverifikasi kontrak tanpa konflik data",
        ], "type": "bullet"},
        {"title": "BUKTI KUANTITATIF", "items": [
            "7 test scenario mencakup branch sukses, gagal, header, dan concurrent",
            "TestProviderBreaking membuktikan minimal 3 error ditangkap secara eksak",
            "Determinisme penuh: tidak ada sleep timer atau flakiness",
        ], "type": "bullet"},
    ],
    mental_model="Deterministic verification completes in milliseconds, clean under Go race detector.",
    source_path="labs/26-contract-testing/engineering/03-execution-result.md",
    slide_num="16 / 24"
)

# Slide 17: Testing Strategy Comparison
add_slide(
    categories=["TABEL EVALUASI", "PERBANDINGAN"],
    title="Perbandingan Metode Pengujian Integrasi",
    subtitle="Mengevaluasi Unit Test, End-to-End Test, dan Contract Test",
    content_items=[
        {"title": "TABEL PERBANDINGAN METODE", "items": [
            "Karakteristik      | Unit Test (Mock) | End-to-End Test | Contract Test (CDC)",
            "-------------------|------------------|-----------------|---------------------",
            "Kecepatan Eksekusi | Sangat Cepat (<1s)| Sangat Lambat(mnt)| Sangat Cepat (ms)",
            "Kestabilan (Flaky) | Nol (Stabil)     | Tinggi (Flaky)  | Nol (Stabil)",
            "Cakupan Integrasi  | Tidak Ada (Palsu)| Penuh           | Format & Semantik",
            "Kebutuhan Lingk.   | In-Process       | Multi-Server/DB | In-Process Ephemeral",
            "Deteksi Breaking   | Gagal            | Terlambat       | Tepat Waktu (CI Gate)",
            "Biaya Pemeliharaan | Rendah           | Sangat Tinggi   | Rendah",
        ], "type": "table"},
        {"title": "ANALISIS NILAI", "items": [
            "CDC mengisi gap kritis antara unit test yang terisolasi dan E2E test yang lambat",
            "Memberikan keyakinan integrasi setara E2E dengan biaya eksekusi setara unit test",
            "Paling efektif diterapkan pada arsitektur microservices terdistribusi",
        ], "type": "bullet"},
        {"title": "POSISI PADA PIRAMIDA TEST", "items": [
            "Contract tests berada di tengah piramida: di atas unit tests, di bawah E2E",
            "Mengurangi kebutuhan suite E2E hingga lebih dari 70%",
        ], "type": "bullet"},
    ],
    mental_model="Contract tests deliver integration-level confidence at unit-test speed and stability.",
    source_path="labs/26-contract-testing/content/01-content-brief.md",
    slide_num="17 / 24"
)

# Slide 18: Pact & Open Source Ecosystem
add_slide(
    categories=["STANDAR INDUSTRI", "EKOSISTEM"],
    title="Standar Industri: Ekosistem Pact & Pact Broker",
    subtitle="Bagaimana prinsip lab ini diaplikasikan pada skala enterprise",
    content_items=[
        {"title": "KOMPONEN PACT FRAMEWORK", "items": [
            "Pact Specification: Standar formal JSON untuk mendefinisikan interaksi HTTP dan message",
            "Pact Language SDKs: pact-go, pact-js, pact-jvm, pact-python",
            "Pact Broker: Server sentral yang menyimpan dan memvalidasi matriks versi kontrak",
        ], "type": "bullet"},
        {"title": "FITUR CAN-I-DEPLOY", "items": [
            "CLI command: pact-broker can-i-deploy --pacticipant OrderService --version git-sha",
            "Memverifikasi apakah versi provider kompatibel dengan semua consumer di environment target",
            "Menjadi penentu exit-code pada pipeline CD sebelum rilis ke Kubernetes/Cloud",
        ], "type": "bullet"},
        {"title": "HUBUNGAN KE LAB 26", "items": [
            "Lab 26 mengimplementasikan esensi verifikasi Pact murni dalam Go standard library",
            "Pola diffing, provider states, dan CI gate yang dipelajari di sini 100% kompatibel secara konseptual",
            "Memahami prinsip inti terlebih dahulu sebelum mengadopsi tooling enterprise",
        ], "type": "bullet"},
    ],
    mental_model="Master the CDC mechanics in pure code before scaling to enterprise Pact Broker workflows.",
    source_path="labs/26-contract-testing/engineering/02-implementation-notes.md",
    slide_num="18 / 24"
)

# Slide 19: Trade-offs & Engineering Costs
add_slide(
    categories=["TRADE-OFF TEKNIK", "BATASAN"],
    title="Trade-off & Biaya Penerapan Contract Testing",
    subtitle="CDC bukan pengganti seluruh jenis pengujian: kenali batasannya",
    content_items=[
        {"title": "BIAYA & TANTANGAN", "items": [
            "Disiplin Organisasi: Membutuhkan kolaborasi erat antara tim consumer dan tim provider",
            "Manajemen Provider State: Menyiapkan data fixture agar provider merespons sesuai kontrak",
            "State Setup Overhead: Endpoint mock/helper untuk seeding database sebelum verifikasi",
        ], "type": "bullet"},
        {"title": "APA YANG TIDAK DIUJI CDC", "items": [
            "Performance & Latency: Tidak mengukur throughput, p99 latency, atau beban DB",
            "Deep Business Logic: Tidak memverifikasi kalkulasi internal akuntansi yang kompleks",
            "Security & Authorization: Tidak menggantikan vulnerability scanner atau penetration test",
        ], "type": "bullet"},
        {"title": "STRATEGI MITIGASI", "items": [
            "Fokuskan CDC pada batas kontrak DTO dan semantik interaksi API",
            "Gabungkan dengan load testing (Lab 15) dan unit testing domain logic",
            "Otomatiskan regenerasi kontrak saat consumer test dijalankan di local",
        ], "type": "bullet"},
    ],
    mental_model="Contract tests guard boundary semantics; business logic and performance still need their own tests.",
    source_path="labs/26-contract-testing/engineering/02-implementation-notes.md:34",
    slide_num="19 / 24"
)

# Slide 20: Implementation Boundaries of Lab 26
add_slide(
    categories=["BATAS LAB", "TRANSPARANSI"],
    title="Batasan Implementasi Lab 26",
    subtitle="Asumsi penyederhanaan yang diambil untuk menjaga fokus pembelajaran",
    content_items=[
        {"title": "BATASAN 1: PURE GO RUNNER", "items": [
            "Menggunakan verifier kustom berbasis Go net/http, bukan binary resmi Pact",
            "Menghilangkan ketergantungan runtime ruby daemon untuk kemudahan eksekusi lokal",
            "Algoritma diffing mencakup subset JSON, primitive types, dan exact string/number",
        ], "type": "bullet"},
        {"title": "BATASAN 2: IN-MEMORY BROKER", "items": [
            "Kontrak dioper secara in-memory/in-process tanpa server Pact Broker eksternal",
            "Pada produksi enterprise, kontrak disimpan di Pactflow atau dedicated server",
            "Pola data flow tetap identik: consumer produce -> provider verify",
        ], "type": "bullet"},
        {"title": "BATASAN 3: PROVIDER STATE SEDERHANA", "items": [
            "Provider state disimulasikan melalui deterministic endpoint routing (/v1/orders/ORD-123)",
            "Tidak mengimplementasikan state setup callback endpoint dinamis",
            "Cukup untuk membuktikan deteksi kegagalan pada 3 jenis breaking changes utama",
        ], "type": "bullet"},
    ],
    mental_model="Understand the deliberate simplifications: pure Go eliminates setup toil while proving the core model.",
    source_path="labs/26-contract-testing/engineering/02-implementation-notes.md:29",
    slide_num="20 / 24"
)

# Slide 21: When to Use Contract Testing
add_slide(
    categories=["PANDUAN KEPUTUSAN", "KAPAN PAKAI"],
    title="Panduan Keputusan: Kapan Menggunakan CDC?",
    subtitle="Kriteria arsitektur dan organisasi yang paling diuntungkan oleh contract testing",
    content_items=[
        {"title": "SANGAT DISARANKAN KETIKA", "items": [
            "Banyak tim microservices independen yang merilis software dengan siklus berbeda",
            "Aplikasi Mobile / Frontend SPA yang terikat pada backend API pihak ketiga",
            "Lingkungan staging lambat, tidak stabil, atau memakan biaya cloud sangat mahal",
            "Regresi API sering terjadi karena perubahan schema yang tidak sengaja",
        ], "type": "bullet"},
        {"title": "KURANG COCOK KETIKA", "items": [
            "Monolith tunggal: seluruh pemanggilan adalah function call internal in-process",
            "API publik eksternal dengan ribuan konsumen tidak dikenal (gunakan schema validation/OpenAPI)",
            "Tim sangat kecil (1-2 orang) di mana komunikasi verbal jauh lebih murah daripada tooling",
        ], "type": "bullet"},
        {"title": "MATRIKS EFISIENSI", "items": [
            "Semakin tinggi jumlah dependensi antar-layanan, semakin besar ROI dari contract testing",
            "Memangkas waktu koordinasi rilis antar-tim dari harian menjadi otomatis di CI",
        ], "type": "bullet"},
    ],
    mental_model="Adopt CDC when cross-team microservice deployments slow down due to integration fear.",
    source_path="labs/26-contract-testing/content/01-content-brief.md",
    slide_num="21 / 24"
)

# Slide 22: Step-by-Step CI Adoption Checklist
add_slide(
    categories=["PANDUAN OPERASIONAL", "CHECKLIST"],
    title="Checklist Adopsi CDC di Pipeline Produksi",
    subtitle="8 langkah taktis mengimplementasikan deployment gate di organisasi",
    content_items=[
        {"title": "CHECKLIST PERSIAPAN & CONSUMER", "items": [
            "[ ] Identifikasi batasan service dan consumer yang paling kritis (misal Mobile App)",
            "[ ] Buat consumer test yang menghasilkan file kontrak interaksi minimal",
            "[ ] Setup publish kontrak ke Pact Broker atau repository terpusat pada pipeline consumer",
        ], "type": "bullet"},
        {"title": "CHECKLIST PROVIDER & GATE", "items": [
            "[ ] Implementasikan provider verification step di CI provider sebelum build container",
            "[ ] Tambahkan automated provider state handler untuk menyiapkan mock DB",
            "[ ] Pasang can-i-deploy gate untuk memblokir merge PR yang melanggar kontrak",
        ], "type": "bullet"},
        {"title": "CHECKLIST EVOLUSI", "items": [
            "[ ] Terapkan pola Expand-Contract saat memutasi schema DTO",
            "[ ] Pantau utilisasi kontrak lama sebelum mematikan endpoint deprecated",
        ], "type": "bullet"},
    ],
    mental_model="Eight concrete steps turn integration anxiety into automated CI/CD confidence.",
    source_path="labs/26-contract-testing/README.md",
    slide_num="22 / 24"
)

# Slide 23: Core Design Principles
add_slide(
    categories=["PRINSIP INTI", "FILOSOFI"],
    title="Tiga Prinsip Inti Contract Testing",
    subtitle="Fondasi filosofis untuk membangun integrasi layanan yang tangguh",
    content_items=[
        {"title": "PRINSIP 1: MINIMAL SUBSET", "items": [
            "Konsumen hanya meminta apa yang benar-benar digunakan",
            "Mencegah pengikatan yang terlalu ketat (tight coupling) pada atribut yang tidak relevan",
            "Memberi ruang bagi provider untuk berevolusi secara independen",
        ], "type": "bullet"},
        {"title": "PRINSIP 2: VERIFIKASI PADA BUILD TIME", "items": [
            "Validasi terjadi sebelum artifact masuk ke container image",
            "Menghilangkan ketergantungan pada runtime staging environment yang rapuh",
            "Feedback loop instan bagi engineer saat masih berada di konteks kode",
        ], "type": "bullet"},
        {"title": "PRINSIP 3: EVOLUSI ADITIF", "items": [
            "Perubahan selalu bersifat aditif melalui dual routing atau new optional fields",
            "Jangan pernah mengubah tipe atau merename field yang terikat kontrak aktif",
            "Migrasi konsumen secara asinkron sebelum mematikan jalur lama",
        ], "type": "bullet"},
    ],
    mental_model="Minimal subsets preserve provider agility; build-time gates eliminate runtime integration dread.",
    source_path="labs/26-contract-testing/engineering/01-design.md",
    slide_num="23 / 24"
)

# Slide 24: Closing (Dark Navy, Summary & Credits)
add_slide(
    categories=["RINGKASAN"],
    title="Ringkasan Lab 26: Contract Testing",
    subtitle=f"Commit: {COMMIT_SHA[:8]} | Test: 7/7 PASS | Repository: lukman-ss/software-engineering-lab",
    content_items=[
        {"title": "HASIL EKSEKUSI UTAMA", "items": [
            "Provider V1 lolos verifikasi kontrak consumer MobileApp secara instan",
            "Breaking changes (enum, rename, primitive mutation) terdeteksi & diblokir di CI gate",
            "Dual Provider (V1 + V2) mendemonstrasikan evolusi schema yang aman dan kompatibel",
        ], "type": "bullet"},
        {"title": "MENTAL MODEL INTI", "items": [
            "Contract tests turn consumer assumptions into executable provider guardrails.",
            "Verify minimal subsets at build time so broken payloads never reach production.",
            "Evolve APIs through parallel endpoints to maintain contract compliance.",
        ], "type": "bullet"},
        {"title": "SUMBER DAYA REPOSITORI", "items": [
            f"Repository: {REPO_URL}",
            f"Commit SHA: {COMMIT_SHA}",
            "Jalankan Demo: go run ./cmd/demo",
            "Jalankan Test: go test -v -race ./...",
        ], "type": "bullet"},
    ],
    mental_model="Contract tests turn consumer assumptions into executable provider guardrails at build time.",
    source_path="labs/26-contract-testing/README.md",
    slide_num=""
)

if __name__ == "__main__":
    print(f"[SUCCESS] Lab {LAB_NO} ({SLUG}): {len(SLIDES)} slides defined successfully.")
