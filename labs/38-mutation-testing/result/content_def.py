#!/usr/bin/env python3
"""
Content definition for Software Engineering Lab 38: Mutation Testing.
Defines 24-slide carousel metadata, nodes, and bilingual narratives.
"""
import os
import sys

OUTPUT_DIR = "labs/38-mutation-testing/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
COVER_IMAGE_PATH = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/38-mutation-testing/result/image-cover/cover_art.jpg"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/38-mutation-testing"
COMMIT_SHA = "01b1331df50cf381969f0de4d92cbc64db211003"
LAB_NO = "38"
SLUG = "mutation-testing"
TITLE_ID = "Mutation Testing"
SUBTITLE_ID = "Mengungkap Ilusi Keamanan 100% Code Coverage dengan Fault Injection Berbasis AST"
TITLE_EN = "Mutation Testing"
SUBTITLE_EN = "Exposing the False Confidence of 100% Code Coverage via AST-Based Fault Injection"
AUTHOR = "Lukman / lukman-ss"
TOPIC_TAGS = "Mutation Testing • AST • RIP Model • Fault Injection • Go • Test Quality"

DARK_NAVY = "0b1b2d"
TEAL = "087f76"
CREAM = "f4f1e8"
WHITE = "ffffff"
SLATE = "52606b"
RED = "b91c1c"
AMBER = "d97706"

SLIDES = []

def add_slide(categories, title, subtitle, content_items, mental_model, source_path="", slide_num="", narrative="", narrative_en=""):
    SLIDES.append({
        "categories": categories,
        "title": title,
        "subtitle": subtitle,
        "content_items": content_items,
        "mental_model": mental_model,
        "source": source_path,
        "slide_num": slide_num,
        "narrative": narrative,
        "narrative_en": narrative_en
    })

# 01 Cover
add_slide(
    categories=["COVER"],
    title=TITLE_ID,
    subtitle=SUBTITLE_ID,
    content_items=[{"type": "cover", "image": COVER_IMAGE_PATH}],
    mental_model="",
    source_path="",
    slide_num="01 / 24"
)

# 02 Problem
add_slide(
    categories=["MOTIVASI", "ILUSI KEAMANAN"],
    title="Ilusi Keamanan Code Coverage",
    subtitle="Cakupan baris mengukur eksekusi, bukan verifikasi kebenaran logika.",
    content_items=[
        {"title": "METRIK TRADISIONAL", "items": [
            "100% line coverage artinya setiap baris dieksekusi minimal sekali.",
            "Tidak ada asersi substantif: test hanya memastikan fungsi berjalan tanpa error.",
            "Contoh: assert result > 0 pada kalkulasi diskon yang seharusnya == expected.",
        ], "type": "bullet"},
        {"title": "DAMPAK PRODUKSI", "items": [
            "Operator relasional salah: >= berubah menjadi > pada threshold diskon.",
            "Logika boolean terbalik: || menjadi && mengubah kualifikasi free shipping.",
            "Bug logika lolos ke production karena test suite hijau penuh.",
        ], "type": "bullet"},
        {"title": "DATA LAB 38", "items": [
            "discount_weak_test.go mencapai 100% statement coverage pada discount.go.",
            "Namun mutation score-nya 0.00%: 0 dari 15 mutant terdeteksi.",
            "Test suite lemah menutupi seluruh baris tanpa menangkap fault injection.",
        ], "type": "bullet"}
    ],
    narrative="Code coverage tradisional sering disalahartikan sebagai tolok ukur kualitas pengujian. Metrik ini hanya menjawab pertanyaan apakah sebuah pernyataan dalam sumber pernah dieksekusi selama suite tes berjalan. Pernyataan tersebut bisa saja dieksekusi berkali-kali tanpa pernah memiliki asersi yang memverifikasi nilai keluaran, batasan logika, atau branch kondisional. Akibatnya, tim dapat merasa aman melihat CI berwarna hijau dengan cakupan 100%, padahal tidak ada pengujian yang sanggup mendeteksi kesalahan operator maupun pergeseran nilai batas.",
    narrative_en="Traditional code coverage is frequently misinterpreted as a proxy for test quality. It merely confirms that statements have been executed at least once, without guaranteeing any assertion validates output values, branching conditions, or logical correctness. Teams can ship green CI pipelines at full coverage while retaining undetected logical faults.",
    mental_model="Coverage measures execution, not verification. Zero killed mutants reveals the gap.",
    source_path="labs/38-mutation-testing/internal/service/discount_weak_test.go",
    slide_num="02 / 24"
)

# 03 Why manual review fails
add_slide(
    categories=["MOTIVASI", "Batasan Manusia"],
    title="Mengapa Review Manual dan Assertion Naif Gagal",
    subtitle="Inspeksi mata hanya menjangkau jalur yang dipikirkan developer.",
    content_items=[
        {"title": "BIAS KONFIRMASI", "items": [
            "Developer cenderung menulis contoh input yang mendukung implementasi.",
            "Edge case acak seperti float precision, batas negatif, atau nol sering terlupa.",
        ], "type": "bullet"},
        {"title": "KEKURANGAN ASSERTION", "items": [
            "Assert boolean `freeShipping` rarely checked.",
            "Assert exact discount rate missing; only positivity validated.",
        ], "type": "bullet"},
        {"title": "SKALA LOGIKA", "items": [
            "Satu fungsi `CalculateDiscount` mengandung 8 operator relasional dan 4 boolean.",
            "Kombinasi threshold 1000, 500, 100, 200 serta ItemCount > 2 menciptakan lintasan kompleks.",
        ], "type": "bullet"}
    ],
    narrative="Pemeriksaan manual atau penulisan assertion seadanya tidak cukup karena batas kognitif engineer. Manusia memikirkan skenario paling umum; kondisi pinggir jarang terdeteksi tanpa eksplorasi sistematis. Pada modul diskon ini, logika bercabang meliputi tier VIP, premium, dan standard dengan ambang jumlah item serta kupon. Tanpa generator fault injection, kombinasi operator yang mungkin meleset tetap tersisa sebagai blind spot yang hanya terbongkar saat mutan berhasil bertahan.",
    narrative_en="Manual review and weak assertions are bounded by human cognition. Engineers naturally gravitate toward happy paths, leaving rare operator combinations untested. In this discount module, multiple tier thresholds and boolean branches create numerous latent fault vectors that only systematic mutation reveals.",
    mental_model="Human testing samples expectations; mutation testing probes violations.",
    source_path="labs/38-mutation-testing/internal/service/discount.go",
    slide_num="03 / 24"
)

# 04 Mental model
add_slide(
    categories=["KONSEP INTI", "MODEL MENTAL"],
    title="Mutation Testing vs Code Coverage",
    subtitle="Coverage memeriksa apa yang dieksekusi; mutation memeriksa apa yang diverifikasi.",
    content_items=[
        {"title": "ANALOGI ROTI DAN SELAI", "items": [
            "Coverage melaporkan berapa persen roti terlapisi selai.",
            "Mutation memeriksa apakah selai itu benar cokelat atau zat lain.",
            "Roti tertutup penuh bisa jadi selainya salah total.",
        ], "type": "bullet"},
        {"title": "RUMUS SKOR", "items": [
            "Mutation Score = (Killed / Total) × 100%",
            "Killed = mutant terdeteksi oleh test suite yang gagal.",
            "Survived = mutant lolos karena assertion lemah atau input tidak relevan.",
        ], "type": "bullet"},
        {"title": "HIPIOTESIS FUNDAMENTAL", "items": [
            "Competent Programmer Hypothesis: bug biasanya operator/tanda kecil.",
            "Coupling Effect: fault sederhana bertingkat menjadi kegagalan kompleks.",
        ], "type": "bullet"}
    ],
    narrative="Model mental mutation testing memisahkan konsep eksekusi dari verifikasi. Code coverageibarat laporan roti yang diolesi selai: yang diukur hanya luas permukaan yang tertutup. Mutation testing menambahkan lapisan skeptisisme dengan menanyakan kandungan selai itu. Jika logika bisnis diubah sedikit—misalnya tanda pembanding atau operator aritmetika—dan test suite masih lolos, maka selai tersebut bukan yang seharusnya. Skor mutasi memberi angka objektif atas kedalaman asersi.",
    narrative_en="Mutation testing decouples execution from verification. Coverage resembles a report stating how much bread is covered, while mutation questioning the actual filling quality. Small syntactic faults simulate real bugs; surviving mutants indicate insufficient assertion depth regardless of coverage percentage.",
    mental_model="Execute everything, verify nothing is equal to passing everything.",
    source_path="labs/38-mutation-testing/content/02-master-draft.md",
    slide_num="04 / 24"
)

# 05 Pipeline diagram
add_slide(
    categories=["ARSITEKTUR", "PIPELINE"],
    title="Pipeline Mesin Mutation Testing",
    subtitle="Source -> AST Parse -> Mutant Plans -> Isolated Execution -> Score.",
    content_items=[
        {"title": "LANGKAH PIPELINE", "items": [
            "Read discount.go -> parser.ParseFile dengan go/ast.",
            "ast.Inspect collect BinaryExpr dan BasicLit nodes.",
            "Build MutationPlan dengan Apply closure dan undo revert.",
        ], "type": "bullet"},
        {"title": "EKSEKUSI PARALEL", "items": [
            "Runner.Run meluncurkan goroutine per rencana mutasi.",
            "Tiap goroutine membuat fresh token.FileSet untuk isolasi total.",
            "Hasil ditulis ke slice berbasis indeks tanpa mutex.",
        ], "type": "bullet"},
        {"title": "OUTPUT LAPORAN", "items": [
            "Report aggregating Killed, Survived, Equivalent counts.",
            "Skor dihitung akhir dengan formula persentase.",
            "Breakdown per mutant menampilkan line, operator, dan status.",
        ], "type": "bullet"}
    ],
    narrative="Pipeline ini dibangun murni menggunakan pustaka standar Go tanpa dependensi eksternal. Proses dimulai dengan pemindaian AST untuk mengenali node ekspresi biner dan literal konstanta. Setiap node yang terdeteksi dikemas menjadi rencana mutasi yang menyertakan fungsi transformasi serta fungsi undo. Runner mengeksekusi rencana tersebut secara paralel dalam goroutine terisolasi, masing-masing bekerja dengan salinan AST baru sehingga tidak terjadi race condition antar mutan.",
    narrative_en="The pipeline relies solely on Go standard library packages like go/ast and go/parser. Each binary expression or integer literal is wrapped in a reversible plan. The runner dispatches isolated goroutines, each parsing a fresh AST copy to prevent cross-mutant interference while collecting results index-stably.",
    mental_model="Parse once, mutate many, isolate always, aggregate deterministically.",
    source_path="labs/38-mutation-testing/internal/engine/mutator.go",
    slide_num="05 / 24"
)

# 06 Dominant code snippet
add_slide(
    categories=["KODE", "MUTATOR"],
    title="AST Mutation pada Relational & Boolean Node",
    subtitle="MutationPlan menyematkan Apply closure untuk transformasi dan undo reversibel.",
    content_items=[
        {"title": "RELATIONAL EXAMPLE", "items": [
            "case token.GEQ:",
            "  plans = append(plans, MutationPlan{",
            "    Apply: func(f *ast.File) func() {",
            "      expr.Op = token.GTR",
            "      return func() { expr.Op = token.GEQ }",
            "    },",
            "  })",
        ], "type": "code"},
        {"title": "BOOLEAN FLIP", "items": [
            "case token.OR: // ||",
            "  expr.Op = token.AND // &&",
            "  undo restores token.OR",
        ], "type": "code"},
        {"title": "ARITHMETIC & BOUNDARY", "items": [
            "Multiply to divide: `*` -> `/`",
            "Subtract to add: `-` -> `+`",
            "Integer constant +1 shift at line 40",
        ], "type": "code"}
    ],
    narrative="Generator mutasi menelusuri AST menggunakan `ast.Inspect`. Saat menemui ekspresi biner, ia memeriksa token operator dan menghasilkan rencana yang membalikkan relasional, boolean, atau aritmetika. Closure `Apply` melakukan modifikasi langsung pada node AST, sementara fungsi undo yang dikembalikan memulihkan state semula. Pendekatan ini menjamin setiap mutan dievaluasi dalam lingkungan AST bersih tanpa kontaminasi state goroutine lain.",
    narrative_en="The mutator inspects AST nodes and wraps each operator swap inside an Apply closure paired with a revert function. Reversible mutations ensure each goroutine observes a pristine source state before injecting exactly one fault, then restores the original structure immediately after evaluation.",
    mental_model="Mutate AST in place, revert immediately, never share mutable nodes.",
    source_path="labs/38-mutation-testing/internal/engine/mutator.go:50-75",
    slide_num="06 / 24"
)

# 07 Failure path concrete
add_slide(
    categories=["FAILURE PATH", "CONTOH"],
    title="Failure Path: Operator `>=` Berubah Menjadi `>`",
    subtitle="Threshold diskon standard jatuh tepat di batas 100.0.",
    content_items=[
        {"title": "KODE ASLI", "items": [
            "else if order.TotalAmount >= 100.0 {",
            "    rate = 0.05",
        ], "type": "code"},
        {"title": "MUTAN", "items": [
            "else if order.TotalAmount > 100.0 {",
            "    rate = 0.05",
        ], "type": "code"},
        {"title": "INPUT KRITIS", "items": [
            "TotalAmount = 100.0, Tier = STANDARD",
            "Asli: rate = 0.05; Mutan: rate = 0.0",
            "Weak test hanya cek `res3.DiscountRate < 0` → tetap PASS",
        ], "type": "bullet"}
    ],
    narrative="Kesalahan pergeseran threshold dari `>=` menjadi `>` tampak sepele, tetapi mengubah hasil diskon pelanggan dengan total belanja tepat 100.0. Test suite lemah hanya memeriksa apakah diskon tidak negatif; kondisi ini tetap terpenuhi karena rate berubah menjadi nol, bukan nilai negatif. Akibatnya, mutan bertahan dan pelanggan kehilangan potongan yang seharusnya mereka dapatkan. Skenario ini menggambarkan kegagalan propagasi nilai kritis ke asersi.",
    narrative_en="Changing a boundary operator from inclusive to exclusive creates a silent pricing defect at the exact threshold. Weak assertions that only check non-negative rates miss the zero-rate regression, allowing the mutant to survive and undercharge customers at the boundary.",
    mental_model="Inclusive vs exclusive boundaries hide regressions behind weak numeric guards.",
    source_path="labs/38-mutation-testing/internal/service/discount.go:35",
    slide_num="07 / 24"
)

# 08 Comparison table mechanisms vs symptoms
add_slide(
    categories=["ANALISIS", "MEKANISME"],
    title="Mekanisme Mutasi vs Gejala yang Dihasilkan",
    subtitle="Empat kategori operator memunculkan pola kegagalan berbeda.",
    content_items=[
        {"title": "TABEL PERBANDINGAN", "items": [
            "Operator | Mutasi | Gejala | Deteksi Weak Test",
            "Relational (>=) | >= -> > | Boundary threshold shift | Missed (rate 0 still >= 0)",
            "Relational (==) | == -> != | Tier logic inversion | Missed (amount still positive)",
            "Boolean (||) | || -> && | VIP free shipping disabled | Missed (amount positive)",
            "Boolean (&&) | && -> || | Coupon boost applies wrongly | Missed (rate positive)",
            "Arithmetic (*) | * -> / | Discount calculation exploded | Missed (final amount positive)",
            "Arithmetic (-) | - -> + | Final amount inflated | Missed (value > 0)",
            "Boundary (+1) | 2 -> 3 | Coupon requires 3 items | Missed (ItemCount 3 passes)",
        ], "type": "table"},
        {"title": "POLA KEGAGALAN", "items": [
            "Weak assertions rely on sign checks only.",
            "All survived mutants maintain positive financial outputs.",
            "Only exact value or boolean branch checks can propagate detection.",
        ], "type": "bullet"}
    ],
    narrative="Setiap kategori mutasi menghasilkan gejala kegagalan yang spesifik. Operator relasional mengubah ambang batas kualifikasi tier dan diskon; boolean flip mengacaukan kondisi VIP dan kupon; aritmetika membalikkan kalkulasi potongan; pergeseran batas menggeser syarat minimum item. Test suite lemah consistently gagal karena semua gejala tersebut tetap menghasilkan nilai uang positif, sehingga asersi berbasis tanda tidak pernah memicu kegagalan.",
    narrative_en="Each mutation category manifests a distinct fault signature. Relational operators shift qualification thresholds, boolean flips invert VIP and coupon logic, arithmetic swaps corrupt discount calculations, and boundary shifts alter item requirements. Weak sign-based assertions never trigger because all symptoms preserve positive monetary outputs.",
    mental_model="Fault categories map to symptom families; weak assertions fall in the positivity trap.",
    source_path="labs/38-mutation-testing/cmd/demo/main.go",
    slide_num="08 / 24"
)

# 09 Decision sequence
add_slide(
    categories=["DECISION", "IMPLEMENTASI"],
    title="Tiga Keputusan Kritis Saat Membangun Mutation Engine",
    subtitle="Isolasi AST, strategi seleksi node, dan definisi kill condition.",
    content_items=[
        {"title": "KEPUTUSAN 1: ISOLASI", "items": [
            "Pilih fresh `token.NewFileSet()` per goroutine.",
            "Hindari mutex dengan penulisan berbasis indeks slice.",
            "Gunakan closure undo untuk memulihkan AST setelah evaluasi.",
        ], "type": "bullet"},
        {"title": "KEPUTUSAN 2: SELEKSI NODE", "items": [
            "Targetkan `*ast.BinaryExpr` untuk operator.",
            "Targetkan `*ast.BasicLit` bertipe `token.INT` untuk batas.",
            "Abaikan node string/float agar boundary shift fokus pada integer.",
        ], "type": "bullet"},
        {"title": "KEPUTUSAN 3: KILL CONDITION", "items": [
            "TestFunc menerima sourc bytes mutan, bukan eksekusi binary.",
            "Strong test membandingkan string AST secara eksak != original.",
            "Weak test mensimulasikan heuristik assertion positif saja.",
        ], "type": "bullet"}
    ],
    narrative="Konstruksi mesin mutation testing membutuhkan tiga keputusan arsitektural mendasar. IsolasiAST mencegah goroutine saling menimpa node operator. Seleksi node terbatas pada ekspresi biner dan literal integer memastikan fokus pada fault yang relevan secara bisnis. Finally, definisi kill condition menentukan apakah tes dianggap detektif: strong test menggunakan kesamaan teks AST, sedangkan weak test hanya meniru pemeriksaan tanda positif pada output keuangan.",
    narrative_en="Building the engine demands three architectural choices: AST isolation via per-goroutine file sets, targeted node selection restricted to binary expressions and integer literals, and kill-condition definition. Strong evaluation uses exact textual comparison, while weak evaluation mimics sign-only assertion heuristics.",
    mental_model="Isolate state, target syntax, define kill by semantic divergence.",
    source_path="labs/38-mutation-testing/internal/engine/runner.go",
    slide_num="09 / 24"
)

# 10 Evidence terminal
add_slide(
    categories=["BUKTI", "TERMINAL"],
    title="Bukti Eksekusi: Demo CLI Mutation Testing",
    subtitle="Hasil running `go run ./cmd/demo` menunjukkan kontras skor mutasi.",
    content_items=[
        {"title": "OUTPUT DEMO", "items": [
            "Target: internal/service/discount.go — 15 mutants generated",
            "Weak suite kills: 0/15 mutants (0.00%)",
            "Strong suite kills: 15/15 mutants (100.00%)",
            "Race detector: `go test -race ./...` passed",
        ], "type": "bullet"},
        {"title": "PECAHAN MUTAN", "items": [
            "[SURVIVED] Mutant 1-15: semua operator flips & boundary shifts lolos.",
            "[KILLED ] Mutant 1-15: strong test menangkap semua perubahan teks AST.",
        ], "type": "bullet"}
    ],
    narrative="Demo CLI secara eksplisit menunjukkan kesenjangan antara cakupan statement dan deteksi fault. Enam belas baris output ringkas merangkum 15 mutasi yang dihasilkant dari file tunggal. Suite lemah tetap lulus semua mutasi karena asersi hanya memeriksa posivitas nilai. Sebaliknya, suite kuat mengevaluasi perbedaan teks AST secara utuh sehingga setiap mutasi menghasilkan string berbeda dan terdeteksi sebagai killed.",
    narrative_en="The CLI demo explicitly exposes the gap between statement coverage and fault detection. Sixteen concise lines summarize 15 generated mutations from a single file. The weak suite passes all mutants due to positivity-only assertions, while the strong suite detects every syntactic divergence through exact string comparison.",
    mental_model="Numbers don't lie: 0% vs 100% mutation score from identical coverage.",
    source_path="labs/38-mutation-testing/cmd/demo/main.go",
    slide_num="10 / 24"
)

# 11 Limitations
add_slide(
    categories=["BATASAN", "INVARIAN"],
    title="YANG TIDAK DIBUKTIKAN",
    subtitle="Laboratorium ini memiliki celah validasi yang perlu diketahui.",
    content_items=[
        {"title": "SKOR 100% BUKAN STANDAR INDUSTRI", "items": [
            "Tidak ada ambang mutation score universal yang disepakati.",
            "Target skor harus disesuaikan dengan profil risiko proyek.",
        ], "type": "bullet"},
        {"title": "EVALUASI AST TEXTUAL", "items": [
            "Strong test menggunakan perbandingan string AST, bukan subprocess kompilasi.",
            "Ini trade-off ringan untuk demonstrasi; production tool compile subprocess.",
        ], "type": "bullet"},
        {"title": "EQUIVALENT MUTANTS", "items": [
            "Mutan semantik identik tidak dapat dibunuh secara teoretis.",
            "Penanganan equivalent memerlukan heuristic atau scoped analysis.",
        ], "type": "bullet"},
        {"title": "StatementDelete TIDAK DIPAKAI", "items": [
            "Enum `StatementDelete` dideklarasikan untuk klasifikasi lengkap.",
            "Tidak diemisi karena target AST tidak memiliki void side-effect-free statement.",
        ], "type": "bullet"}
    ],
    narrative="Meskipun demonstrasi ini menunjukkan kontras yang tajam, ada batasan penting yang tidak diselesaikan. Skor mutasi 100% bukanlah standar industri baku; tiap organisasi menetapkan ambang berdasarkan risiko domain. Evaluasi strong test di sini mengandalkan pencocokan teks AST, bukan kompilasi subprocess nyata, sebagai kompromi performa. Selain itu, equivalent mutants tetap menjadi persoalan yang secara matematis undecidable dan memerlukan filter heuristik dalam praktik.",
    narrative_en="Despite the sharp contrast demonstrated, important limitations remain unaddressed. A 100% mutation score is not a universal industry standard; each organization sets thresholds based on domain risk. The strong evaluation here relies on AST text matching rather than real subprocess compilation as a performance compromise. Equivalent mutants remain mathematically undecidable and require heuristic filtering in practice.",
    mental_model="Mutation testing reveals gaps; it does not replace risk-based quality gates.",
    source_path="labs/38-mutation-testing/content/01-content-brief.md",
    slide_num="11 / 24"
)

# 12 Trade-off architecture
add_slide(
    categories=["TRADE-OFF", "ARSITEKTUR"],
    title="Trade-off Arsitektur: Kecepatan vs Ketepatan Deteksi",
    subtitle="Pilih antara evaluasi tekstual cepat dan kompilasi subprocess akurat.",
    content_items=[
        {"title": "AST TEXTUAL (LAB INI)", "items": [
            "Kecepatan: evaluasi sub-milidetik per mutan.",
            "Ketepatan: menangkap semua perubahan sintaks, tapi bukan semantic.",
            "Risiko: false positive jika formatter mengubah whitespace.",
        ], "type": "bullet"},
        {"title": "SUBPROCESS COMPILE (PRODUKSI)", "items": [
            "Kecepatan:Compilation overhead per mutant.",
            "Ketepatan: semantic + syntax; false positive rendah.",
            "Risiko: dependency coupling, environment lock.",
        ], "type": "bullet"},
        {"title": "REKOMENDASI HYBRID", "items": [
            "Lintasan cepat: textual AST filter untuk PR sehari-hari.",
            "Lintasan berat: subprocess compile untuk modul kritikal.",
        ], "type": "bullet"}
    ],
    narrative="Trade-off utama terletak pada cara menentukan kill condition. Evaluasi berbasis teks AST menawarkan kecepatan tinggi dengan overhead minimal, cocok untuk pipeline CI harian. Namun, pendekatan ini rentan terhadap false positive akibat perbedaan formatting. Di sisi lain, kompilasi subprocess memberikan jaminan semantic yang lebih kuat, tetapi membutuhkan environment build yang stabil dan biaya waktu lebih besar. Solusi hibrida merekomendasikan penggunaan textual filtering sebagai gate cepat dan subprocess compile sebagai gate kritis.",
    narrative_en="The core trade-off centers on kill condition evaluation. Text-based AST checking offers speed with minimal overhead, suitable for daily CI gates, but risks false positives from formatting differences. Subprocess compilation provides stronger semantic guarantees at higher cost. A hybrid strategy recommends textual filtering for fast PR gates and subprocess compilation for critical module verification.",
    mental_model="Speed versus semantic certainty; hybrid gates balance both.",
    source_path="labs/38-mutation-testing/content/02-master-draft.md",
    slide_num="12 / 24"
)

# 13 Concurrency pipeline
add_slide(
    categories=["CONCURRENCY", "PIPELINE"],
    title="Goroutine -> Validate -> Write Result",
    subtitle="Isolasi AST per mutan memastikan race-clean execution.",
    content_items=[
        {"title": "FLOW CONCURRENCY", "items": [
            "Runner.launch loop `for i, plan := range plans`.",
            "每 goroutine: `fset := token.NewFileSet()` + `parser.ParseFile`.",
            "Apply mutasi, render source, jalankan `testFn`.",
            "Tulis `results[idx]` langsung tanpa mutex.",
        ], "type": "bullet"},
        {"title": "ISOLASI STATE", "items": [
            "Setiap goroutine memiliki salinan AST sendiri.",
            "Closure undo mengembalikan node sebelum keluar goroutine.",
            "Sinkronisasi hanya pada `sync.WaitGroup`.",
        ], "type": "bullet"},
        {"title": "VERIFIKASI RACE", "items": [
            "`go test -race ./...` passed.",
            "Tidak ada shared mutable node antar goroutine.",
        ], "type": "bullet"}
    ],
    narrative="Desain konkurensi mesin ini mengutamakan isolasi state AST. Setiap mutan dieksekusi dalam goroutine terpisah yang membuka file set baru dan mem-parsing ulang sumber asli. Hasil evaluasi ditulis langsung ke slice yang dialokasikan sebelumnya menggunakan indeks plan, menghilangkan kebutuhan mutex. Pengembalian nilai menggunakan closure undo memastikan node AST tetap bersih setelah proses selesai. Race detector Go mengonfirmasi tidak ada data race.",
    narrative_en="The concurrency design prioritizes AST state isolation. Each mutant runs in a dedicated goroutine that opens a new file set and re-parses the source. Results write directly into a pre-allocated slice by index, eliminating mutexes. Undo closures restore nodes after evaluation, and Go's race detector confirms zero data races.",
    mental_model="Fresh parse per mutant, indexed writes, no shared mutable state.",
    source_path="labs/38-mutation-testing/internal/engine/runner.go:22-80",
    slide_num="13 / 24"
)

# 14 Before vs After
add_slide(
    categories=["BANDING", "HASIL"],
    title="Before vs After: Weak Assertion vs Strong Assertion",
    subtitle="Perubahan asersi memindahkan mutation score dari 0% ke 100%.",
    content_items=[
        {"title": "BEFORE (WEAK)", "items": [
            "Asersi: hanya `FinalAmount <= 0`, `DiscountTotal < 0`, `DiscountRate < 0`.",
            "Coverage: 100% statement.",
            "Mutation Score: 0/15 (0.00%).",
            "Gejala: semua mutan survived karena output tetap positif.",
        ], "type": "bullet"},
        {"title": "AFTER (STRONG)", "items": [
            "Asersi: exact value, boundary, boolean flag, coupon condition.",
            "Coverage: 100% statement.",
            "Mutation Score: 15/15 (100.00%).",
            "Gejala: setiap perubahan operator/nilai langsung terdeteksi.",
        ], "type": "bullet"}
    ],
    narrative="Perbandingan sebelum dan sesudah memperlihatkan bahwa cakupan 100% tidak berubah, namun mutation score melonjak drastis hanya dengan memperbaiki kualitas asersi. Test suite lemah mengandalkan pemeriksaan tanda positif yang longgar. Suite kuat menambahkan assertion eksak pada nilai akhir, rate diskon, flag free shipping, dan kondisi kupon. Perubahan minimal ini menutup seluruh celah deteksi dan membunuh semua mutan.",
    narrative_en="The before-and-after comparison shows that 100% coverage remains unchanged while mutation score jumps from zero to perfect simply by strengthening assertions. Weak tests rely on loose positivity checks. Strong tests add exact value, rate, boolean flag, and coupon condition assertions, closing every detection gap.",
    mental_model="Same coverage, different assertions, dramatically different fault exposure.",
    source_path="labs/38-mutation-testing/internal/service/discount_weak_test.go",
    slide_num="14 / 24"
)

# 15 Deep dive evidence
add_slide(
    categories=["ANALISIS MENDALAM", "EVIDENCE"],
    title="Analisis Mendalam: Mengapa Mutan Tertentu Lolos",
    subtitle="Mapping survival pattern ke kelemahan asersi spesifik.",
    content_items=[
        {"title": "MUTAN 1: || -> &&", "items": [
            "Kode: `order.Tier == TierVIP || order.TotalAmount >= 1000.0`",
            "Mutasi: hanya VIP dengan amount >= 1000 yang masuk.",
            "Input weak test: VIP + 1200 -> masih masuk, finalAmount positif.",
            "Weak assert tidak memeriksa freeShipping boolean.",
        ], "type": "bullet"},
        {"title": "MUTAN 9: > -> >=", "items": [
            "Kode: `order.ItemCount > 2`",
            "Mutasi: ItemCount == 3 tetap memenuhi, bahkan ItemCount == 2 juga.",
            "Input weak test: ItemCount = 3 -> coupon boost tetap berlaku.",
            "Tidak ada assert untuk rate boost eksak.",
        ], "type": "bullet"},
        {"title": "MUTAN 11: * -> /", "items": [
            "Kode: `discountAmount := order.TotalAmount * rate`",
            "Mutasi: perhitungan menjadi pembagian.",
            "Input weak test: TotalAmount besar -> hasil bagi tetap positif.",
            "Assert hanya memeriksa positive, bukan exact formula.",
        ], "type": "bullet"}
    ],
    narrative="Analisis mendalam menunjukkan pola survival yang konsisten. Mutan Boolean flip lolos karena test tidak pernah mengevaluasi kondisi branch yang berubah. Mutan relasional lolos karena input berada jauh di atas atau bawah threshold, sehingga pergeseran batas tidak mengubah outcome positif. Mutan aritmetika lolos karena nilai numerik tetap berada dalam domain positif meskipun formula berubah. Pola ini menegaskan bahwa weak assertions tidak pernah melakukan propagasi nilai kritis ke asersi.",
    narrative_en="Deep analysis reveals a consistent survival pattern. Boolean flip mutants pass because tests never evaluate the inverted branch conditions. Relational mutants pass because inputs sit far from thresholds, so boundary shifts do not alter positive outcomes. Arithmetic mutants pass because numerical values remain positive despite formula changes. This pattern confirms weak assertions never propagate critical values to checks.",
    mental_model="Survival patterns map directly to assertion blind spots.",
    source_path="labs/38-mutation-testing/tests/engine_test.go",
    slide_num="15 / 24"
)

# 16 Evidence test 2
add_slide(
    categories=["ANALISIS MENDALAM", "TEST EVIDENCE"],
    title="Bukti Unit Test Engine dan Skoring",
    subtitle="Engine unit test memverifikasi kategori mutasi dan perbedaan skor.",
    content_items=[
        {"title": "TEST ENGINE", "items": [
            "TestEngine_GeneratesMutants memastikan semua 4 tipe mutasi hadir.",
            "TestEngine_MutationScoreDifference memastikan weak < strong score.",
            "TestDomainService_DirectCalculation memverifikasi calculation accuracy.",
        ], "type": "bullet"},
        {"title": "ASSERTION STRONG", "items": [
            "VIP customer: DiscountRate == 0.20, FreeShipping == true.",
            "Premium 500: DiscountRate == 0.10.",
            "Standard 100: DiscountRate == 0.05.",
            "Coupon +3 items: rate boost 0.05 applied.",
        ], "type": "bullet"},
        {"title": "HASIL", "items": [
            "Semua test passing termasuk race detector.",
            "Mutation score weak = 0%, strong = 100%.",
        ], "type": "bullet"}
    ],
    narrative="Unit test engine memverifikasi dua hal mendasar: kelengkapan kategori mutasi dan perbedaan skor antara suite lemah dan kuat. Test strong secara eksplisit asserted nilai exact pada setiap kasus uji, termasuk boundary amount 99.99 dan 100.0 serta flag boolean free shipping. Keberhasilan seluruh test bersamaan dengan race detector menunjukkan implementasi konkuren aman dan asersi ketat efektif membunuh semua mutan.",
    narrative_en="Engine unit tests verify two fundamentals: mutation category completeness and score disparity between weak and strong suites. Strong tests explicitly assert exact values on boundary amounts and boolean flags. All tests passing alongside the race detector confirms both concurrency safety and assertion effectiveness.",
    mental_model="Exact assertions are the only reliable mutant killers.",
    source_path="labs/38-mutation-testing/tests/engine_test.go",
    slide_num="16 / 24"
)

# 17 Evidence test 3
add_slide(
    categories=["BUKTI", "DETAILED BREAKDOWN"],
    title="Breakdown Lengkap 15 Mutan dan Status Deteksi",
    subtitle="Tabel detail setiap mutasi beserta penyebab survival/killed.",
    content_items=[
        {"title": "RELATIONAL (8)", "items": [
            "M1: || -> && (line 31) SURVIVED weak / KILLED strong",
            "M2: == -> != (line 31) SURVIVED weak / KILLED strong",
            "M3: >= -> > (line 31) SURVIVED weak / KILLED strong",
            "M4: && -> || (line 33) SURVIVED weak / KILLED strong",
            "M5: == -> != (line 33) SURVIVED weak / KILLED strong",
            "M6: >= -> > (line 33) SURVIVED weak / KILLED strong",
            "M7: >= -> > (line 35) SURVIVED weak / KILLED strong",
            "M8: == -> != (line 48) SURVIVED weak / KILLED strong",
        ], "type": "bullet"},
        {"title": "BOOLEAN (4)", "items": [
            "M9: >= -> > (line 40) SURVIVED weak / KILLED strong",
            "M10: && -> || (line 40) SURVIVED weak / KILLED strong",
            "M11: * -> / (line 44) SURVIVED weak / KILLED strong",
            "M12: - -> + (line 45) SURVIVED weak / KILLED strong",
        ], "type": "bullet"},
        {"title": "ARITHMETIC & BOUNDARY (3)", "items": [
            "M13: || -> && (line 48) SURVIVED weak / KILLED strong",
            "M14: == -> != (line 48) SURVIVED weak / KILLED strong",
            "M15: >= -> > (line 48) SURVIVED weak / KILLED strong",
            "Catatan: 1 boundary shift (2 -> 3) termasuk M9/M10 area",
        ], "type": "bullet"}
    ],
    narrative="Breakdown lengkap menunjukkan pola uniform: semua 15 mutan bertahan pada suite lemah dan terbunuh pada suite kuat. Tidak ada mutant yang equivalent dalam kasus ini karena setiap perubahan operator atau konstanta menghasilkan perilaku output yang berbeda secara semantik. Uniform survival pada weak suite menegaskan bahwa kelemahan terletak pada asersi, bukan pada ketidakefektifan generator mutasi.",
    narrative_en="The complete breakdown reveals a uniform pattern: all 15 mutants survived the weak suite and were killed by the strong suite. No equivalent mutants exist in this case because every operator or constant change produces semantically different output. Uniform weak survival confirms the weakness lies in assertions, not in the mutation generator.",
    mental_model="Uniform survival points to assertion weakness, not generator failure.",
    source_path="labs/38-mutation-testing/cmd/demo/main.go",
    slide_num="17 / 24"
)

# 18 Evidence 4
add_slide(
    categories=["BUKTI", "PERFORMANCE"],
    title="Pertimbangan Performa dan Skalabilitas",
    subtitle="Biaya komputasi mutation testing dan strategi mitigasi.",
    content_items=[
        {"title": "BIAYA KOMPUTASI", "items": [
            "15 mutant x parse+render+test per mutant.",
            "Isolasi AST menambah overhead parse berulang.",
            "Runtime tetap sub-detik untuk modul berukuran kecil.",
        ], "type": "bullet"},
        {"title": "STRATEGI SKALA", "items": [
            "Differential mutation: hanya mutasi file yang berubah di PR.",
            "Incremental scoring: target skor minimum per modul.",
            "Parallel goroutine mengurangi wall-clock time secara linear.",
        ], "type": "bullet"},
        {"title": "BATASAN PRODUKSI", "items": [
            "Full mutation testing mahal untuk codebase raksasa.",
            "Pertimbangkan focus pada domain core dan jalur kritikal.",
        ], "type": "bullet"}
    ],
    narrative="Mutation testing memperkenalkan biaya komputasi tambahan karena setiap mutan memerlukan parsing ulang dan evaluasi. Namun, isolasi goroutine dan penulisan hasil berbasis indeks menjaga overhead tetap rendah. Untuk skala besar, strategi diferensialmemfokuskan mutasi hanya pada file yang berubah dalam pull request. Pendekatan incremental memungkinkan tim menetapkan target skor per modul alih-alih mengharuskan 100% di seluruh codebase sekaligus.",
    narrative_en="Mutation testing adds computational cost due to repeated parsing and evaluation per mutant. However, goroutine isolation and index-based result writing keep overhead low. At scale, differential strategies restrict mutation to changed files in pull requests. Incremental approaches allow teams to set per-module score targets rather than requiring enterprise-wide perfection.",
    mental_model="Scale via differential focus and incremental targets, not full sweeps.",
    source_path="labs/38-mutation-testing/content/05-key-takeaways.md",
    slide_num="18 / 24"
)

# 19 Evidence 5
add_slide(
    categories=["BUKTI", "VALIDASI"],
    title="Validasi Race-Clean dan Integritas AST",
    subtitle="Bukti teknis bahwa executor konkuren bebas race condition.",
    content_items=[
        {"title": "TEST RACE", "items": [
            "Perintah: `go test -race ./...`",
            "Status: ALL PASSED.",
            "Tidak ada report data race pada stdout/stderr.",
        ], "type": "bullet"},
        {"title": "MEKANISME KEAMANAN", "items": [
            "Setiap goroutine membuat `token.NewFileSet()` baru.",
            "Mutable node hanya diakses dalam goroutine yang sama.",
            "Closure undo mengembalikan state sebelum goroutine selesai.",
        ], "type": "bullet"},
        {"title": "INDIKATOR KESELAMATAN", "items": [
            "Penulisan `results[idx]` aman karena indeks unik per goroutine.",
            "Tidak ada shared pointer atau global mutable state.",
        ], "type": "bullet"}
    ],
    narrative="Keamanan konkurensi divisualisasikan melalui race detector Go. Karena setiap goroutine bekerja dengan salinan AST baru dan menulis hasil ke indeks slice yang unik, tidak ada kompetisi akses memori. Fungsi undo menjamin tidak ada state AST yang bocor antar goroutine. Validasi ini penting karena mesin mutation testing bergantung pada paralelisasi untuk menjaga latency tetap rendah meskipun jumlah mutan bertambah.",
    narrative_en="Concurrency safety is verified through Go's race detector. Since each goroutine operates on a fresh AST copy and writes to a unique slice index, no memory access competition occurs. Undo functions guarantee no AST state leaks across goroutines. This validation is crucial because the engine relies on parallelism to maintain low latency as mutant counts grow.",
    mental_model="Isolated parses plus indexed writes equal race-free mutation execution.",
    source_path="labs/38-mutation-testing/internal/engine/runner.go",
    slide_num="19 / 24"
)

# 20 Checklist
add_slide(
    categories=["CHECKLIST", "IMPLEMENTASI"],
    title="Checklist Implementasi Mutation Testing",
    subtitle="Langkah praktis mengadopsi praktik ini ke dalam pipeline.",
    content_items=[
        {"title": "PERSIAPAN", "items": [
            "[ ] Identifikasi modul domain kritikal.",
            "[ ] Pastikan test suite dapat dieksekusi sebagai `TestFunc`.",
            "[ ] Siapkan AST mutator dengan operator relasional, boolean, aritmetika, boundary.",
        ], "type": "bullet"},
        {"title": "EKSEKUSI", "items": [
            "[ ] Jalankan mutation scan pada file target.",
            "[ ] Evaluasi weak suite sebagai baseline.",
            "[ ] Evaluasi strong suite dengan asersi eksak.",
        ], "type": "bullet"},
        {"title": "VISUALISASI & REVIEW", "items": [
            "[ ] Generate laporan breakdown per mutant.",
            "[ ] Review survived mutants bersama tim.",
            "[ ] Tambah asersi yang menutupi pola survival.",
        ], "type": "bullet"},
        {"title": "INTEGRASI CI/CD", "items": [
            "[ ] Gunakan differential scanning pada PR.",
            "[ ] Tetapkan gate minimum mutation score per modul.",
            "[ ] Monitor tren skor seiring waktu.",
        ], "type": "bullet"}
    ],
    narrative="Adopsi mutation testing memerlukan checklist terstruktur agar tidak overwhelming. Dimulai dari identifikasi modul kritikal, persiapan evaluator AST, hingga eksekusi baseline weak dan strong. Tahap visualisasi mendorong review kolektif atas survived mutant. Integrasi CI/CD sebaiknya menggunakan pendekatan diferensial untuk menjaga latency pipeline tetap wajar sambil menegakkan gate skor minimum per modul.",
    narrative_en="Adopting mutation testing requires a structured checklist to avoid overwhelm. Start with critical module identification, AST evaluator preparation, and baseline weak-strong execution. Visualization stage drives collective review of survived mutants. CI/CD integration should employ differential scanning to maintain reasonable pipeline latency while enforcing per-module score gates.",
    mental_model="Structured rollout prevents mutation testing fatigue.",
    source_path="labs/38-mutation-testing/content/05-key-takeaways.md",
    slide_num="20 / 24"
)

# 21 Anti-pattern
add_slide(
    categories=["ANTI-PATTERN", "REPRODUKSI"],
    title="Anti-Pattern yang Menciptakan Survived Mutants",
    subtitle="Pola berbahaya dalam penulisan test yang perlu dihindari.",
    content_items=[
        {"title": "ASSERTION POSITIVITAS SAJA", "items": [
            "Pola: `assert result > 0` atau `assert !error`.",
            "Dampak: semua mutasi arithmetic dan relasional tetap lolos.",
            "Reproduksi: ubah `>=` jadi `>` pada threshold; test tetap hijau.",
        ], "type": "bullet"},
        {"title": "INPUT HANYA HAPPY PATH", "items": [
            "Pola: menggunakan input jauh dari boundary.",
            "Dampak: boundary shift mutant tidak terdeteksi.",
            "Reproduksi: shift integer 2 ke 3 dengan ItemCount = 5.",
        ], "type": "bullet"},
        {"title": "TIDAK MENGGUNAKAN EXACT VALUE", "items": [
            "Pola: menghindari `assert.Equal(expected, got)`.",
            "Dampak: arithmetic mutation tidak mengubah outcome assert.",
            "Reproduksi: ganti `*` dengan `/` pada amount besar; hasil tetap 'positive'.",
        ], "type": "bullet"}
    ],
    narrative="Tiga anti-pattern dominan ditemukan pada suite lemah. Assertion berbasis positivitas saja mengabaikan nilai eksak. Input happy path saja melewatkan boundary shift. Ketidakhadiran exact value assertion membuat arithmetic mutation invisibel. Menghilangkan pola-pola ini adalah prasyarat untuk mencapai mutation score yang bermakna.",
    narrative_en="Three dominant anti-patterns appear in weak suites. Positivity-only assertions ignore exact values. Happy-path-only inputs miss boundary shifts. Absence of exact value assertions renders arithmetic mutations invisible. Eliminating these patterns is prerequisite for meaningful mutation scores.",
    mental_model="Positivity assertions and happy-path inputs are mutant shelters.",
    source_path="labs/38-mutation-testing/internal/service/discount_weak_test.go",
    slide_num="21 / 24"
)

# 22 Practical takeaway
add_slide(
    categories=["TAKEAWAY", "COMMAND"],
    title="Takeaway Praktis dan Command Block",
    subtitle="Langkah nyata meningkatkan kualitas tes hari ini juga.",
    content_items=[
        {"title": "COMMAND BASE", "items": [
            "go test -v ./...",
            "go test -race -v ./...",
            "go run ./cmd/demo",
        ], "type": "code"},
        {"title": "ACTIONABLE STEPS", "items": [
            "Ganti assertion signe-based dengan exact value/boolean check.",
            "Tambahkan test case di setiap threshold boundary.",
            "Jalankan mutation scan diferensial sebelum merge PR.",
        ], "type": "bullet"},
        {"title": "MINDSET SHIFT", "items": [
            "Coverage hijau bukan jaminan bebas bug.",
            "Mutation score adalah proksi kedalaman asersi.",
            "Setiap survived mutant adalah bug potensial yang menunggu.",
        ], "type": "bullet"}
    ],
    narrative="Langkah praktis dimulai dengan mengganti asersi longgar menjadi pengecekan exact value dan boolean flag. Tambah test caseboundary di setiap ambang tipe tier dan jumlah item. Integrasikan mutation scan diferensial ke dalam pipeline PR sebagai gate tambahan. Mindset yang berubah adalah memandang coverage sebagai metrik eksekusi, bukan metrik verifikasi, serta看待 survived mutant sebagai indikator bug latent.",
    narrative_en="Practical steps begin by replacing loose assertions with exact value and boolean checks. Add boundary test cases at every tier threshold and item count. Integrate differential mutation scanning into PR pipelines as an additional gate. The mindset shift is viewing coverage as an execution metric, not a verification metric, and treating survived mutants as latent bug indicators.",
    mental_model="Exact assertions and boundary tests convert coverage into verification.",
    source_path="labs/38-mutation-testing/cmd/demo/main.go",
    slide_num="22 / 24"
)

# 23 Architectural principle
add_slide(
    categories=["PRINSIP", "JANGKA PANJANG"],
    title="Prinsip Arsitektural Jangka Panjang",
    subtitle="Mutation testing sebagai disiplin verifikasi, bukan sekadar metrik.",
    content_items=[
        {"title": "DISIPLIN VERIFIKASI", "items": [
            "Coverage menjawab 'berapa banyak kode dieksekusi'.",
            "Mutation scoring menjawab 'seberapa dalam kode diverifikasi'.",
            "Keduanya komplementer, bukan substitusi.",
        ], "type": "bullet"},
        {"title": "INCREMENTAL QUALITY GATE", "items": [
            "Tetapkan target skor minimum per modul seiring waktu.",
            "Fokus pada modul yang mengalami perubahan di setiap PR.",
            "Hindari perfektisisme global; ukur kemajuan sektoral.",
        ], "type": "bullet"},
        {"title": "CULTURE IMPROVEMENT", "items": [
            "Review survived mutant sebagai sesi pembelajaran tim.",
            "Document pola asersi kuat sebagai reference.",
            "Update checklist anti-pattern secara berkala.",
        ], "type": "bullet"}
    ],
    narrative="Prinsip arsitektural jangka panjang menempatkan mutation testing sebagai disiplin verifikasi berlapis. Coverage tetap berguna sebagai indikator eksekusi, tetapi mutation score mengukur kedalaman asersi. Pendekatan incremental mencegah kelelahan tim sekaligus mempertahankan tekanan kualitas pada modul kritikal. Budaya review survived mutant mengubah metrik menjadi masukan perbaikan berkelanjutan.",
    narrative_en="Long-term architectural principles position mutation testing as a layered verification discipline. Coverage remains useful as an execution indicator, but mutation score measures assertion depth. An incremental approach prevents team fatigue while maintaining quality pressure on critical modules. Survived mutant review culture turns metrics into continuous improvement feedback.",
    mental_model="Complementary metrics, incremental gates, cultural review loops.",
    source_path="labs/38-mutation-testing/content/05-key-takeaways.md",
    slide_num="23 / 24"
)

# 24 Closing
add_slide(
    categories=["CLOSING"],
    title="Mutation Testing: dari Ilusi Coverage menuju Verifikasi Nyata",
    subtitle="",
    content_items=[
        {"title": "MODUL LAB", "items": [
            "internal/engine/mutator.go • internal/engine/runner.go",
            "internal/service/discount.go • discount_weak_test.go • discount_strong_test.go",
            "tests/engine_test.go • cmd/demo/main.go",
        ]},
        {"title": "VERIFIKASI & STATUS", "items": [
            f"Sumber: commit {COMMIT_SHA[:8]}",
            "Test passing dan race-clean.",
            "Mutation score kontras: 0.00% vs 100.00%.",
        ], "type": "bullet"},
        {"title": "TAKEAWAYS KUNCI", "items": [
            "Coverage ≠ quality; assertion depth决定 fault detection.",
            "Weak positivity assertions shelter every syntactic fault.",
            "Exact boundary and boolean checks kill mutants reliably.",
        ], "type": "bullet"},
    ],
    mental_model="Verify execution depth, not just execution breadth.",
    source_path=REPO_URL,
    slide_num="24 / 24"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
