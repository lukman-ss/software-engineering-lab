#!/usr/bin/env python3
"""
Publication Builder for Software Engineering Lab 17: Architecture Decision Record.
Generates: PPTX, PDF, PNG covers, Dev.to MD, Medium HTML, Metadata, Caption, DELIVERABLES.md
Adheres strictly to the Lab 13-15 editorial 2-column layout and in-depth technical narratives.
"""
import os
import sys

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/17-architecture-decision-record/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/17-architecture-decision-record"
COMMIT_SHA = "def0f5a"
LAB_NO = "17"
SLUG = "architecture-decision-record"
TITLE_ID = "Architecture Decision Record"
SUBTITLE_ID = "Parser & Linter untuk Validasi Struktural Keputusan Arsitektur"
TITLE_EN = "Architecture Decision Record"
SUBTITLE_EN = "Parser and Linter for Structural Validation of Architectural Decisions"
AUTHOR = "Lukman / lukman-ss"
TOPIC_TAGS = "ADR • Parser • Linter • Go • Documentation Integrity"

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

# Slide 1: Cover (Dark Navy)
add_slide(
    categories=["COVER"],
    title=TITLE_ID,
    subtitle=SUBTITLE_ID,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "DAG node tree showing ADR decision lineage with Superseded/Supersedes links",
            "Three connected boxes: ADR 1 (Superseded) → ADR 2 (Accepted) ← ADR 3 (Rejected)",
            f"Penulis: {AUTHOR}",
        ]}
    ],
    mental_model="",
    source_path="",
    slide_num="01 / 24"
)

# Slide 2: Context Loss in Engineering Teams
add_slide(
    categories=["PROBLEM", "ENGINEERING CONFLICT"],
    title="Konteks Keputusan Arsitektur Hilang Tanpa Jejak",
    subtitle="Engineer baru bertanya mengapa sistem dibangun begini; tidak ada jawaban yang terekam.",
    content_items=[
        {"title": "MASALAH NYATA", "items": [
            "Putaran architect merangkap code review: pertanyaan berulang 'kenapa kita pilih X?' di every PR.",
            "Keputusan diambil via chat atau meeting, lalu tenggelam dalam commit history.",
            "ADR lama tidak diperbarui saat sistem berubah; status dokumentasi drift dari kode aktual."
        ], "type": "bullet"},
        {"title": "RISIKO BISNIS", "items": [
            "Re-debate arsitektur lama memakan waktu engineering sprint tanpa solusi baru.",
            "Keputusan yang salah diterapkan karena engineer tidak mengetahui ADR penggantian."
        ], "type": "bullet"}
    ],
    narrative="Dalam praktik engineering sehari-hari, keputusan arsitektur sering kali diambil secara informal—lewat chat, voice note, atau diskusi meeting. Ketika engineer baru bergabung atau sistem mengalami pivot strategis, konteks historis mengapa pilihan tertentu dibuat hilang. Tim kemudian terjebak debat berulang tentang keputusan yang sebenarnya sudah pernah diverifikasi dan ditolak. Ini bukan masalah dokumentasi biasa; ini adalah kehilangan kapabilitas keputusan yang berharga.",
    mental_model="Architectural context without structural records becomes tribal knowledge.",
    source_path="labs/17-architecture-decision-record/content/02-master-draft.md",
    slide_num="02 / 24"
)

# Slide 3: Why Manual Markdown Review Fails
add_slide(
    categories=["WHY FAILS", "HUMAN REVIEW"],
    title="Review Markdown Manual Lupa Relasi Penggantian",
    subtitle="Engineer bisa saja update isi ADR lama tapi lupa tambahkan tautan supersesi ke ADR baru.",
    content_items=[
        {"title": "KEGAGALAN REVIEW MANUAL", "items": [
            "ADR 1 tetap berstatus 'Accepted' padahal sudah diganti ADR 2 — lupa di-update.",
            "ADR 2 menulis 'Supersedes: 1', tapi ADR 1 tidak pernah dicatat 'Superseded by 2'.",
            "Penomoran loncat: ADR 1, ADR 3 langsung — engineer lupa insert ADR 2.",
            "Status informal seperti 'Draft', 'In Review' diterima tanpa penolakan."
        ], "type": "bullet"},
        {"title": "CONTOH KASUS NYATA", "items": [
            "ADR 1 ber-status Superseded by 99, padahal ADR 99 belum pernah dibuat."
        ], "type": "bullet"}
    ],
    narrative="Ketika tim mengandalkan review manual terhadap dokumen Markdown ADR, pola kegagalan yang paling sering muncul adalah tautan supersesi yang tidak sinkron. Engineer bisa jadi memperbarui isi keputusan pada ADR baru namun lupa mengembalikan marker pada ADR lama. Atau sebaliknya: ADR baru mengklanimsupersedes_ADR lama, tapi ADR lama tidak pernah mengubah statusnya menjadi Superseded. Review manual juga sulit mendeteksi loncatan nomor urut karena tidak ada mekanisme enforcement struktural.",
    mental_model="Manual review catches content errors but misses structural relation gaps.",
    source_path="labs/17-architecture-decision-record/content/02-master-draft.md",
    slide_num="03 / 24"
)

# Slide 4: Mental Model — ADR as Structured Data
add_slide(
    categories=["MENTAL MODEL", "CORE CONCEPT"],
    title="ADR Bukan Dokumen, Tapi Record yang Bisa Dicek",
    subtitle="Parser mengubah baris Markdown menjadi struct Record; linter memeriksa aturan struktural.",
    content_items=[
        {"title": "PRINSIP PENTING", "items": [
            "ADR diolah sebagai data, bukan sekadar teks human-readable.",
            "Parser mengekstrak ID, Title, Status, SupersededBy, Supersedes, Content.",
            "Linter memvalidasi aturan: monotonik ID, valid status, relasi timbal-balik.",
            "Keputusan arsitektur tetap ditulis manusiawi; mesin hanya memeriksa format."
        ], "type": "bullet"}
    ],
    narrative="Inti dari pendekatan lab ini adalah memandang Architecture Decision Record bukan sebagai dokumen naratif yang perlu dibaca manusia secara manual, melainkan sebagai record terstruktur yang bisa diparse dan divalidasi oleh mesin. Parser membaca Markdown, mengekstrak field-field kunci menjadi struct Record, dan linter kemudian menjalankan serangkaian aturan struktural. Nilai utama bukan pada otomasi pengambilan keputusan, melainkan pada kemampuan mendeteksi anomali struktural yang sering luput dari review manual.",
    mental_model="ADR is structured data first, human narrative second; validation catches what humans miss.",
    source_path="labs/17-architecture-decision-record/internal/adr/models.go • internal/adr/parser.go",
    slide_num="04 / 24"
)

# Slide 5: System Architecture — Parser to Linter Pipeline
add_slide(
    categories=["ARSITEKTUR", "DATA FLOW"],
    title="Pipeline: Markdown → Parser → Records → Linter",
    subtitle="Alur konversi ADR dari teks bebas menjadi record terstruktur dan validasi paralel.",
    content_items=[
        {"title": "ARCHITECTURE FLOW", "items": [
            "[Markdown ADR Text] → adr.Parser (Regex extraction: ID, Title, Status, refs)",
            "→ []*adr.Record (In-Memory Structs)",
            "→ adr.Linter.Validate(records)",
            "  ├─ Monotonic numbering check: sort + sequential scan",
            "  ├─ Cycle detection: DFS with visiting/visited states",
            "  └─ Bidirectional lineage: goroutines + sync.Mutex per-record validation",
            "→ []error (empty = valid, non-empty = integrity breach)"
        ], "type": "diagram"}
    ],
    narrative="Arsitektur sistem validasi ADR pada lab ini dibangun sebagai pipeline tiga tahap. Tahap pertama, parser membaca Markdown baris demi baris menggunakan regex dan mengekstrak field-field kunci. Hasilnya adalah slice struct Record yang berisi metadata keputusan. Tahap kedua, linter menjalankan validasi urutan ID monotonik dan deteksi siklus melalui depth-first search. Tahap ketiga, validasi graf silsilah dilakukan secara paralel menggunakan goroutine dengan mutex untuk agregasi error, memastikan setiap tautan SupersededBy dan Supersedes saling konsisten.",
    mental_model="Pipeline separates parsing from validation: each stage has a single responsibility.",
    source_path="internal/adr/parser.go • internal/adr/linter.go",
    slide_num="05 / 24"
)

# Slide 6: Key Code Path — Record Struct
add_slide(
    categories=["CODE", "MODEL"],
    title="Struct Record: Definisi Inti ADR",
    subtitle="Lima field wajib membuat setiap keputusan arsitektur dapat direferensikan dan diverifikasi.",
    content_items=[
        {"title": "MODELS.GO", "items": [
            "type Status string",
            "",
            "const (",
            "    StatusProposed   Status = \"Proposed\"",
            "    StatusAccepted   Status = \"Accepted\"",
            "    StatusSuperseded Status = \"Superseded\"",
            "    StatusDeprecated Status = \"Deprecated\"",
            "    StatusRejected   Status = \"Rejected\"",
            ")",
            "",
            "type Record struct {",
            "    ID           int",
            "    Title        string",
            "    Status       Status",
            "    SupersededBy int // 0 if not superseded",
            "    Supersedes   int // 0 if does not supersede",
            "    Content      string",
            "}"
        ], "type": "code"}
    ],
    narrative="Struktur data Record pada models.go merupakan fondasi seluruh sistem validasi. Field Status menggunakan enum string yang hanya mengakui lima nilai resmi: Proposed, Accepted, Superseded, Deprecated, dan Rejected. Field SupersededBy mereferensikan ADR pengganti, sedangkan Supersedes mereferensikan ADR yang digantikan. Keduanya bernilai 0 (nol) ketika tidak ada relasi. Dengan representasi terstruktur ini, linter tidak perlu menebak-tebak — ia hanya memeriksa konsistensi integer dan string.",
    mental_model="Explicit fields over implicit conventions make structural validation possible.",
    source_path="internal/adr/models.go",
    slide_num="06 / 24"
)

# Slide 7: Failure Mode — Dangling Supersession Link
add_slide(
    categories=["FAILURE", "REGRESSION TEST"],
    title="Supersesi Patah: ADR 1 Merujuk ke ADR 99",
    subtitle="Tautan Superseded by menunjuk ke nomor yang tidak ada — linter langsung menandai error.",
    content_items=[
        {"title": "TEST CASE: broken_references", "items": [
            "superseded points to non-existent:",
            "  ADR 1: Status Superseded by 99",
            "  → Error: superseded by non-existent ADR 99",
            "",
            "supersedes points to non-existent:",
            "  ADR 1: Status Accepted, Supersedes: 99",
            "  → Error: supersedes non-existent ADR 99",
        ], "type": "code"},
        {"title": "BUKTI DARI TEST", "items": [
            "TestLinter_BrokenReferences / superseded points to non-existent → PASS",
            "TestLinter_BrokenReferences / supersedes points to non-existent → PASS"
        ], "type": "bullet"}
    ],
    narrative="Skenario kegagalan paling umum pada pemeliharaan ADR manual adalah tautan supersesi yang merujuk ke nomor ADR yang belum ada atau sudah dihapus. Linter menangkap kasus ini dengan memverifikasi keberadaan setiap ID referensi di dalam recordMap. Jika ADR 1 menyatakan Superseded by 99 namun tidak ada Record dengan ID 99, pesan error spesifik langsung dihasilkan. Bukti ini berasal dari test case TestLinter_BrokenReferences yang dijalankan secara deterministik.",
    mental_model="A dangling reference is cheaper to catch at lint time than during an architectural review.",
    source_path="tests/linter_test.go:36-54",
    slide_num="07 / 24"
)

# Slide 8: Hidden Cost — Unverified Lineage Drift
add_slide(
    categories=["COST", "DOCUMENTATION DRIFT"],
    title="Biaya Tersembunyi: Lineage Drift Tanpa Deteksi",
    subtitle="Setiap ADR baru yang tidak mengikuti aturan relasi meningkatkan kerumunan manual review.",
    content_items=[
        {"title": "DRIFT PATTERN", "items": [
            "Time to detect | Cause | Impact",
            "Instant (CI) | Linter reject di PR | Build gagal, engineer feedback loop cepat",
            "Mingguan | Manual audit dokumentasi | Engineer kebingungan siapa yang aktif",
            "Bulanan | Production incident review | Tim kehilangan konteks keputusan penting",
            "Selamanya | Keputusan salah diterapkan | Sistem berkembang dengan asumsi keliru"
        ], "type": "table"}
    ],
    narrative="Biaya tersembunyi dari ketiadaan validasi otomatis bukan hanya waktu deteksi, tetapi akumulasi kepercayaan palsu bahwa dokumentasi ADR masih akurat. Setiap kali engineer menambahkan ADR baru tanpa memutakhirkan ADR lama, tercipta discrepancy yang hanya akan terdeteksi saat inspeksi manual. Pada tim kecil, hal ini mungkin tidak terasa. Namun begitu skala tim bertambah, waktu yang terbuang untuk mencari konsistensi manual berlipat ganda tanpa bound yang jelas.",
    mental_model="Unverified lineage drift accumulates silently until a production decision depends on it.",
    source_path="labs/17-architecture-decision-record/content/02-master-draft.md",
    slide_num="08 / 24"
)

# Slide 9: Experiment Setup — Three ADR Scenarios
add_slide(
    categories=["EXPERIMENT", "DEMO SETUP"],
    title="Setup Eksperimen: 3 ADR Evolusi Arsitektur",
    subtitle="Demo memodelkan transisi dari Modular Monolith ke Microservice dengan satu ADR ditolak.",
    content_items=[
        {"title": "ADR SEQUENCE", "items": [
            "ADR 0001: Use Modular Monolith for Core SaaS ERP | Status: Superseded by 2",
            "ADR 0002: Extract Notification Service to Microservice | Status: Accepted, Supersedes: 1",
            "ADR 0003: Reject Event Sourcing for Order Management | Status: Rejected",
            "",
            "Test suite: parser_test.go (valid + superseded + supersedes + invalid cases)",
            "Test suite: linter_test.go (valid sequence + 7 broken reference scenarios + concurrency)"
        ], "type": "bullet"}
    ],
    narrative="Eksperimen lab ini memodelkan tiga skenario ADR yang merepresentasikan evolusi arsitektur nyata. ADR 1 merekam keputusan awal menggunakan Modular Monolith. ADR 2 menggantikan ADR 1 setelah tim tumbuh melebihi 20 engineer dan membutuhkan scaling independen untuk notification service. ADR 3 menolak usulan Event Sourcing karena kebutuhan audit trail masih terpenuhi oleh relational log table standar. Ketiganya divalidasi bersama melalui linter untuk memastikan konsistensi lineagenya.",
    mental_model="Three ADRs capture one realistic architectural evolution path with valid lineage.",
    source_path="cmd/demo/main.go • tests/linter_test.go",
    slide_num="09 / 24"
)

# Slide 10: Result Snapshot — Demo Output
add_slide(
    categories=["RESULT", "EVIDENCE"],
    title="Output Demo: Parsing Berhasil, Linter Lolos",
    subtitle="Terminal menampilkan 3 ADR terparsing dan integrity check passed tanpa error.",
    content_items=[
        {"title": "CLI OUTPUT", "items": [
            "$ go run ./cmd/demo",
            "=== Architecture Decision Record (ADR) Lab ===",
            "Successfully parsed 3 ADRs:",
            "  - [0001] Use Modular Monolith for Core SaaS ERP | Status: Superseded -> Superseded by ADR 2",
            "  - [0002] Extract Notification Service to Microservice  | Status: Accepted   -> Supersedes ADR 1",
            "  - [0003] Reject Event Sourcing for Order Management    | Status: Rejected  ",
            "",
            "Running ADR Integrity Linter...",
            "Integrity check passed! Decision lineage and lifecycle invariants are intact."
        ], "type": "code"}
    ],
    narrative="Output demo CLI membuktikan bahwa ketiga ADR berhasil diparse dari string Markdown dan diterjemahkan menjadi struct Record. Lineage chain ADR 1 → ADR 2 terdeteksi sebagai valid: ADR 1 berstatus Superseded dengan reference ke 2, dan ADR 2 berstatus Accepted dengan klaim Supersedes: 1. Linter tidak menemukan error apa pun pada validasi monotonicity, cycle detection, maupun bidirectional link integrity.",
    mental_model="Clean linter output confirms structural consistency, not architectural correctness.",
    source_path="cmd/demo/main.go • tests/linter_test.go",
    slide_num="10 / 24"
)

# Slide 11: Result Interpretation — What Passing Means
add_slide(
    categories=["INTERPRETASI", "EVIDENCE"],
    title="Apa Artinya Linter Passed di Lab Ini?",
    subtitle="Linter passing membuktikan struktur record konsisten, bukan kebenaran keputusan arsitektur.",
    content_items=[
        {"title": "YANG DIBUKTIKAN", "items": [
            "Format parsing: judul, status, dan tautan supersesi terstruktur dengan benar.",
            "Monotonik: ID 1, 2, 3 berurutan tanpa celah.",
            "Bidirectional: ADR 1 ↔ ADR 2 memiliki tautan timbal-balik yang lengkap.",
            "Cycle-free: tidak ada sirkuit dalam graf supersesi."
        ], "type": "bullet"},
        {"title": "YANG TIDAK DIBUKTIKAN", "items": [
            "Kebenaran keputusan arsitektur di dalam field Content.",
            "Apakah ADR 1 seharusnya benar-benar digantikan ADR 2 di dunia nyata.",
            "Kesiapan sistem untuk production deployment."
        ], "type": "bullet"}
    ],
    narrative="Hasil linter yang passing adalah bukti kondisi necessary, bukan sufficient. Ia membuktikan bahwa ketiga ADR memenuhi invariant struktural: format parsable, nomor urut monotonik, dan tautan supersesi timbal-balik utuh. Namun, linter tidak menilai apakah keputusan Modular Monolith memang tepat, atau apakah penolakan Event Sourcing didasarkan pada analisis yang komprehensif. Kualitas substansi keputusan tetap tanggung jawab review manusia.",
    mental_model="Linter passing is necessary but not sufficient for architectural soundness.",
    source_path="labs/17-architecture-decision-record/engineering/03-execution-result.md",
    slide_num="11 / 24"
)

# Slide 12: Trade-off — Structure vs Flexibility
add_slide(
    categories=["TRADE-OFF", "DECISION"],
    title="Trade-off: Enforced Structure vs Dokumentasi Fleksibel",
    subtitle="Standarisasi format ADR meningkatkan deteksi error, tetapi membatasi ekspresi nuansa konteks.",
    content_items=[
        {"title": "ENFORCED STRUCTURE (Lab Ini)", "items": [
            "Regex-based parsing memaksa format fixed (## Context, ## Decision, ## Consequences).",
            "Status wajib dari 5 enum resmi; status informal ditolak parser.",
            "Supersedes/SupersededBy wajib konsisten secara timbal-balik."
        ], "type": "bullet"},
        {"title": "FLEKSIBILITAS YANG DILEWATI", "items": [
            "Tidak mendukung ADR yang memerlukan struktur sub-section kustom di luar 3 section wajib.",
            "Tidak memvalidasi kualitas argumen dalam section Decision.",
            "Tidak menangani kasus ADR multi-author atau decision yang melibatkan beberapa sub-system."
        ], "type": "bullet"}
    ],
    narrative="Validasi struktural ADR membawa trade-off mendasar: semakin ketat format yang dipaksakan, semakin rendah fleksibilitas penulisan konteks arsitektur. Lab ini memilih path enforceable machine-checkable, artinya parser akan menolak ADR yang tidak mengikuti skema tetap. Di sisi lain, pendekatan ini memungkinkan CI/CD gate yang andal. Tim yang membutuhkan nuansa lebih kaya bisa menambahkan section tambahan di luar yang diparsing, selama tiga section wajib tetap ada.",
    mental_model="Structure enforcement gains verification reliability at the cost of narrative flexibility.",
    source_path="labs/17-architecture-decision-record/content/02-master-draft.md",
    slide_num="12 / 24"
)

# Slide 13: Improved Approach — Concurrent Validation with Mutex
add_slide(
    categories=["IMPROVED", "CONCURRENCY"],
    title="Validasi Paralel: Goroutine + Mutex untuk Agregasi Error",
    subtitle="Setiap ADR divalidasi secara independent goroutine, error dikumpulkan dengan sync.Mutex.",
    content_items=[
        {"title": "CONCURRENT VALIDATION PATTERN", "items": [
            "var wg sync.WaitGroup",
            "for _, r := range records {",
            "    wg.Add(1)",
            "    go func(rec *Record) {",
            "        defer wg.Done()",
            "        var localErrs []error",
            "        // ... validate rec, append localErrs ...",
            "        if len(localErrs) > 0 {",
            "            mu.Lock(); errs = append(errs, localErrs...); mu.Unlock()",
            "        }",
            "    }(r)",
            "}",
            "wg.Wait()"
        ], "type": "code"}
    ],
    narrative="Pola validasi konkuren pada linter menunjukkan praktik Go yang tepat untuk parallel error aggregation. Setiap Record divalidasi dalam goroutine terpisah, menghasilkan slice error lokal yang hanyadigabungkan ke global errs slice setelah acquire mutex. Pola ini menghindari race condition pada write ke shared error slice. Meskipun validasi pada 3 record tidak memberikan gain performa signifikan, pola ini dapat diskalakan untuk hundred-an ADR tanpa perubahan struktur kode.",
    mental_model="Per-record goroutine with mutex aggregation scales validation without shared-state races.",
    source_path="internal/adr/linter.go:64-106",
    slide_num="13 / 24"
)

# Slide 14: Before vs After — Manual Review vs Linter Gate
add_slide(
    categories=["KOMPARASI", "BEFORE AFTER"],
    title="Before vs After: Otomasi Linter Menghilangkan Human Error",
    subtitle="Membandingkan gap deteksi antara review manual dan automated structural linting.",
    content_items=[
        {"title": "SEBELUM LINTER (Manual)", "items": [
            "Lineage break: ADR 1 Superseded by 2, ADR 2 tidak memiliki Supersedes: 1 → TIDAK TERDETEKSI",
            "Non-monotonic ID: ADR 1, ADR 3 meloncat → TIDAK TERDETEKSI selagi tidak inspect semua file",
            "Invalid status 'Draft' masuk ke repo → TIDAK TERDETEKSI oleh regex sederhana",
            "Self-supersession: ADR 1 Superseded by 1 → TIDAK TERDETEKSI"
        ], "type": "bullet"},
        {"title": "SESUDAH LINTER (Automated)", "items": [
            "Semua kasus di atas ditolak oleh TestLinter_BrokenReferences dalam <1ms",
            "CI gate menolak merge jika linter menemukan error",
            "Concurrent stress test membuktikan tidak ada race pada agregasi error"
        ], "type": "bullet"}
    ],
    narrative="Perbandingan before-after menunjukkan bahwa automasi linter menutup celah deteksi yang secara konsisten terlewatkan pada review manual. Empat skenario failure yang tercatat dalam test suite—dangling reference, mismatched link, non-monotonic ID, dan self-supersession—semuanya lolos dari radar reviewer manusia tetapi langsung tertangkap oleh automated check. Ini bukan berarti linter menggantikan peran arsitek, melainkan mengurangi beban kognitif pada hal-hal yang bersifat mekanis.",
    mental_model="Automation closes mechanical review gaps that human reviewers consistently miss.",
    source_path="tests/linter_test.go",
    slide_num="14 / 24"
)

# Slide 15: Boundary of the Lab — Format Validation Only
add_slide(
    categories=["BATASAN", "BOUNDARY"],
    title="Batasan Lab: Validasi Format, Bukan Kualitas Keputusan",
    subtitle="Linter memeriksa kepatuhan struktur; ia tidak menilai substansi atau kebenaran argumen.",
    content_items=[
        {"title": "YANG DILAKUKAN LAB", "items": [
            "Parsing teks Markdown ke struct Record dengan regex.",
            "Memvalidasi urutan ID monotonik dan uniqueness.",
            "Memeriksa konsistensi timbal-balik SupersededBy ↔ Supersedes.",
            "Mendeteksi siklus dalam graf supersesi."
        ], "type": "bullet"},
        {"title": "YANG TIDAK DILAKUKAN LAB", "items": [
            "Tidak menilai kualitas argumen di dalam section Decision.",
            "Tidak memverifikasi apakah ADR cocok dengan implementasi kode aktual.",
            "Tidak memeriksa apakah engineer sudah mendapatkan approval stakeholder sebelum Accepted.",
            "Simulasi ini tidak sama dengan production ADR management system seperti ADRtool atau Architect.js."
        ], "type": "bullet"}
    ],
    narrative="Batasan paling penting dari lab ini adalah ia hanya memeriksa kepatuhan format, bukan kebenaran konten. Parser tidak tahu apakah keputusan 'Extract Notification Service to Microservice' adalah keputusan yang tepat—ia hanya memverifikasi bahwa keputusan tersebut ditulis dalam format yang benar dengan relasi supersesi yang konsisten. Klaim validasi otomatis ini berlaku pada jalur yang diuji, bukan jaminan bahwa seluruh ekosistem dokumentasi arsitektur tim sudah bersih.",
    mental_model="Format validation ensures structural integrity; it does not evaluate architectural soundness.",
    source_path="labs/17-architecture-decision-record/engineering-audit/04-docs-vs-code.md",
    slide_num="15 / 24"
)

# Slide 16: Edge Case — Cancellation Without ADR
add_slide(
    categories=["EDGE CASE", "RACE CONDITION"],
    title="Edge Case: Keputusan Dibatalkan Tanpa ADR Baru",
    subtitle="Engineer menghentikan fitur tanpa membuat ADR Deprecated; state sistem menjadi ambigu.",
    content_items=[
        {"title": "SKENARIO EDGE CASE", "items": [
            "ADR X berstatus Accepted, tapi tim memutuskan berhenti menggunakan fitur tersebut.",
            "Tanpa ADR baru, status ADR X tetap Accepted —误导 engineer lain以为 masih relevan.",
            "Solusi: buat ADR baru berstatus Deprecated merujuk ke alasan penghentian.",
            "Kegagalan: engineer menghapus file ADR langsung → lineage history hilang."
        ], "type": "bullet"},
        {"title": "REKOMENDASI PRODUKSI", "items": [
            "Wajibkan policy tim: setiap penghentian fitur memerlukan ADR Deprecated baru.",
            "Linter bisa diperluas untuk mendeteksi ADR Accepted yang tidak memiliki successor dalam 12 bulan."
        ], "type": "bullet"}
    ],
    narrative="Kasus tepi yang tidak ditangani oleh linter saat ini adalah penghentian keputusan tanpa dokumentasi formal. Ketika engineer secara diametral menghentikan fitur yang sebelumnya di-accept via ADR, status dokumen tetap Accepted meskipun praktik aktual sudah berubah. Ini menciptakan gap antara dokumentasi dan realitas sistem. Solusi yang direkomendasikan adalah policy tim yang mewajibkan ADR Deprecated sebagai pengganti, bukan penghapusan file.",
    mental_model="Silent feature deletion without Deprecated ADR creates documentation-reality gap.",
    source_path="labs/17-architecture-decision-record/content/02-master-draft.md",
    slide_num="16 / 24"
)

# Slide 17: What the Tests Prove
add_slide(
    categories=["BUKTI", "TEST COVERAGE"],
    title="Apa yang Dibuktikan Oleh Test Suite",
    subtitle="Parser dan linter menutupi semua skenario struktural dasar yang didefinisikan dalam lab.",
    content_items=[
        {"title": "PARSER TESTS (parser_test.go)", "items": [
            "TestParse_Valid: ID=1, Status=Accepted → PASS",
            "TestParse_Superseded: SupersededBy=2 berhasil diekstrak → PASS",
            "TestParse_Supersedes: Supersedes=1 berhasil diekstrak → PASS",
            "TestParse_Invalid: missing_title, missing_status, invalid_status, empty sections → semua PASS"
        ], "type": "bullet"},
        {"title": "LINTER TESTS (linter_test.go)", "items": [
            "TestLinter_ValidSequence: 3 ADR valid, zero errors → PASS",
            "TestLinter_BrokenReferences: 7 sub-case (non-existent ref, mismatch, non-monotonic, self-super, duplicate ID, cyclical) → semua PASS",
            "TestLinter_ConcurrencyStress: concurrent validation, no data race → PASS"
        ], "type": "bullet"}
    ],
    narrative="Suite pengujian lab ini mencakup 10 test function dengan 18 sub-case yang membuktikan seluruh invariant struktural yang didefinisikan. Parser menangani 6 kasus validasi format (termasuk 9 sub-case invalid). Linter menangani 8 skenario integritas graf (valid sequence, 7 broken reference types). Race detector pada go test -race juga lulus tanpa warning. Bukti ini cakupannya mengikuti skenario test yang didefinisikan — bukan klaim bahwa seluruh kemungkinan edge case production tertutup.",
    mental_model="Test coverage proves structural invariants on defined scenarios, not all possible production states.",
    source_path="tests/parser_test.go • tests/linter_test.go",
    slide_num="17 / 24"
)

# Slide 18: What the Tests Do Not Prove
add_slide(
    categories=["BATASAN BUKTI", "DO NOT OVERCLAIM"],
    title="Apa yang Tidak Dibuktikan Oleh Test Suite",
    subtitle="Test hanya menjangkau jalur yang di-skripsikan; production masih perlu validasi tambahan.",
    content_items=[
        {"title": "BUKAN BUKNIT TEST", "items": [
            "Kebenaran konten arsitektur di dalam section Decision dan Consequences.",
            "Adopsi ADR oleh seluruh anggota tim sesuai format yang ditetapkan.",
            "Integrasi linter dengan CI/CD pipeline (hanya test Go murni, bukan GitHub Actions).",
            "Skalabilitas linter pada ribuan ADR (hanya 3 record dalam demo, 3-7 dalam test)."
        ], "type": "bullet"},
        {"title": "SIMULASI IN-MEMORY", "items": [
            "Parser bekerja pada string literal, bukan file system read dari ribuan dokumen Markdown."
        ], "type": "bullet"}
    ],
    narrative="Penting untuk mengenali batas bukti dari test suite ini. Test tidak menilai kebenaran substansi keputusan arsitektur yang terdapat di dalam field Content. Test juga tidak membuktikan bahwa linter akan bekerja dengan baik ketika diterapkan pada ratusan file ADR di repositori production. Simulasi berjalan pada string literal di dalam kode, bukan pada file system I/O aktual. Klaim validasi struktural yang dihasilkan test berlaku pada jalur yang diuji, bukan generalisasi ke seluruh spektrum kemungkinan error di lapangan.",
    mental_model="Unit tests prove behavior on defined inputs; they do not guarantee production behavior.",
    source_path="labs/17-architecture-decision-record/engineering-audit/05-gaps.md",
    slide_num="18 / 24"
)

# Slide 19: Production Risk — False Confidence from Green Linter
add_slide(
    categories=["RISIKO", "PRODUCTION"],
    title="Risiko Produksi: False Confidence dari Linter Hijau",
    subtitle="Linter passing bisa memberikan ilusi keamanan yang menyesatkan jika substansi tidak dicek manusia.",
    content_items=[
        {"title": "FALSE POSITIVE RISKS", "items": [
            "Tim percaya ADR collection-nya 'bersih' hanya karena linter passing, padahal argumen Decision lemah.",
            "Engineer bisa mengisi Content dengan placeholder kosong selama section ada → linter tetap passing.",
            "Policy tim yang tidak enforcing bisa di-skip tanpa konsekuensi jika gate CI tidak dipasang."
        ], "type": "bullet"},
        {"title": "MITIGASI", "items": [
            "Tambah human review checkpoint: ADR Accepted wajib ditinjau architecture board.",
            "Validasi minimum content length pada section Context dan Decision."
        ], "type": "bullet"}
    ],
    narrative="Risiko terbesar dari automasi validasi adalah false confidence. Ketika linter menghasilkan output hijau tanpa error, tim cenderung menganggap seluruh aspek dokumentasi arsitektur sudah OK. Padahal, Content field bisa saja berisi penjelasan yang sangat minim atau bahkan placeholder. Linter passing bukan berarti keputusan arsitektur tersebut berkualitas tinggi — ia hanya berarti formatnya konsisten secara struktural. Mitigasi utamanya adalah mempertahankan human review sebagai gate tambahan.",
    mental_model="A green linter pass validates structure, not wisdom. Human review remains irreplaceable.",
    source_path="labs/17-architecture-decision-record/engineering-audit/06-verdict.md",
    slide_num="19 / 24"
)

# Slide 20: Implementation Checklist
add_slide(
    categories=["CHECKLIST", "IMPLEMENTASI"],
    title="Implementation Checklist: Menerapkan ADR Linter di Tim",
    subtitle="Langkah-langkah praktis untuk mengintegrasikan parser-linter ke workflow tim.",
    content_items=[
        {"title": "CHECKLIST LANGKAH", "items": [
            "[ ] Define ADR template Markdown dengan 3 section wajib: Context, Decision, Consequences.",
            "[ ] Implement parser regex yang mengekstrak ID, Title, Status, SupersededBy, Supersedes.",
            "[ ] Tambahkan status enum valid (Proposed, Accepted, Superseded, Deprecated, Rejected).",
            "[ ] Implement linter dengan validasi monotonik + bidirectional supersession + cycle detection.",
            "[ ] Integrasikan ke CI/CD: go test ./... atau custom binary di PR gate.",
            "[ ] Dokumentasikan policy tim: setiap perubahan keputusan wajib ADR baru, bukan edit retroaktif."
        ], "type": "checklist"}
    ],
    narrative="Implementasi ADR linter di lingkungan production memerlukan enam langkah praktis. Dimulai dari definisi template Markdown standar yang disepakati tim, dilanjutkan implementasi parser dan linter sesuai pola yang divalidasi di lab. Integrasi CI/CD memastikan violation langsung ditolak pada PR. Langkah terakhir — yang paling sering terlewat — adalah policy enforcement: mengubah budaya tim agar setiap perubahan keputusan arsitektur dicatat sebagai ADR baru, bukan suntingan retroaktif pada ADR lama.",
    mental_model="Checklist ensures consistent adoption; culture change ensures sustained compliance.",
    source_path="labs/17-architecture-decision-record/engineering/02-implementation-notes.md",
    slide_num="20 / 24"
)

# Slide 21: Anti-Pattern — Common Mistakes When Adopting ADR
add_slide(
    categories=["ANTI-PATTERN", "LESSON"],
    title="Anti-Pattern: 5 Kebiasaan Buruk Saat Mengadopsi ADR",
    subtitle="Kesalahan umum yang mengurangi efektivitas ADR dan cara menghindarinya.",
    content_items=[
        {"title": "ANTI-PATTERN SERING TERJADI", "items": [
            "Mengklik ADR untuk segala hal: refactor kecil, penamaan variabel → noise berlebihan.",
            "Mengedit ADR yang sudah Accepted secara silent, tanpa membuat ADR baru (retroactive edit).",
            "Menulis ADR hanya sebagai catatan meeting tanpa section Decision yang eksplisit.",
            "Status 'Draft' bertahan terlalu lama tanpa transisi ke Accepted atau Rejected.",
            "Melupakan update SupersededBy pada ADR lama saat ADR baru dibuat."
        ], "type": "bullet"}
    ],
    narrative=" Lima anti-pattern ini sering muncul saat tim pertama kali mengadopsi ADR. Pola yang paling merusak adalah retroactive edit — suntingan diametral pada ADR yang sudah Accepted. Ini menghancurkan jejak keputusan dan membuat engineer berikutnya membaca konteks yang sudah tidak relevan. ADR linter lab ini dirancang untuk mencegah sebagian besar anti-pattern tersebut secara struktural, meskipun tidak menutupi masalah adopsi budaya tim.",
    mental_model="ADR anti-patterns are cultural; linter catches structural symptoms, not behavioral roots.",
    source_path="labs/17-architecture-decision-record/content/02-master-draft.md",
    slide_num="21 / 24"
)

# Slide 22: Reproduction Guide
add_slide(
    categories=["REPRODUKSI", "COMMAND"],
    title="Reproduksi Mandiri: Jalankan Demo dan Test di Lokal",
    subtitle="Semua dependency menggunakan Go standard library, tanpa framework eksternal.",
    content_items=[
        {"title": "COMMAND RUNNER", "items": [
            "cd labs/17-architecture-decision-record",
            "go test -v ./...",
            "go test -race ./...",
            "go run ./cmd/demo"
        ], "type": "code"},
        {"title": "EXPECTED OUTPUT", "items": [
            "TestParse_*: semua PASS",
            "TestLinter_*: semua PASS (14 sub-case total)",
            "go test -race: race detector clean",
            "Demo: 3 ADR parsed, Integrity check passed"
        ], "type": "bullet"}
    ],
    narrative="Laboratorium ini dirancang untuk reproduksi instan di mesin manapun yang menjalankan Go 1.22+. Tidak ada dependensi eksternal, tidak ada database, tidak ada network call. Command go run ./cmd/demo menjalankan demo CLI yang memuat tiga ADR dari string literal, mem-parsenya, dan menjalankan linter. Command go test -race memastikan tidak ada race condition pada konkurensi linter. Seluruh artefak testing dan catatan audit tersedia di direktori lab untuk studi lebih lanjut.",
    mental_model="Zero-dependency Go standard library ensures instant reproduction on any machine.",
    source_path="README.md • go.mod • tests/",
    slide_num="22 / 24"
)

# Slide 23: Technical Takeaway
add_slide(
    categories=["TAKEAWAY", "ENGINEERING PRINCIPLE"],
    title="ADR Diubah Jadi Data yang Bisa Dicek, Bukan Hanya Dokumen",
    subtitle="Parser+linter memindahkan keputusan arsitektur dari domain narasi ke domain struktur terverifikasi.",
    content_items=[
        {"title": "PELAJARAN REKAYASA", "items": [
            "ADR yang diparse menjadi struct Record memungkinkan validasi otomatis yang tidak bisa dilakukan humans dengan andal.",
            "Bidirectional supersession link adalah invarian yang paling sering rusak; linter menangkapnya di detik pertama.",
            "Konkurensi goroutine + mutex pattern untuk error aggregation adalah reusable pattern untuk validation framework lain.",
            "Format standar membantu review; kualitas keputusan tetap ditentukan isi ADR."
        ], "type": "bullet"}
    ],
    narrative="Pelajaran rekayasa perangkat lunak utama dari lab ini adalah pergeseran perspektif: dari ADR sebagai dokumen naratif statis menjadi ADR sebagai record terstruktur yang dapat diverifikasi mesin. Parser dan linter tidak menggantikan judgment arsitek, tetapi mereka menutupi gap mekanisme yang konsisten terlewat oleh review manual. Pola validasi konkuren yang diterapkan di linter juga bersifat generik dan dapat diadaptasi untuk validasi dokumen struktural lain di luar konteks ADR.",
    mental_model="Structured validation catches what human review consistently misses: mechanical relation gaps.",
    source_path="labs/17-architecture-decision-record/content/05-key-takeaways.md",
    slide_num="23 / 24"
)

# Slide 24: Closing Slide (Dark Navy)
add_slide(
    categories=["CLOSING"],
    title="Parser membuat ADR terstruktur; linter baru memeriksa aturannya.",
    subtitle="",
    content_items=[
        {"title": "MODUL LAB", "items": [
            "internal/adr/models.go • internal/adr/parser.go • internal/adr/linter.go",
            "cmd/demo/main.go • tests/parser_test.go • tests/linter_test.go",
            "README.md • content/02-master-draft.md"
        ], "type": "bullet"},
        {"title": "VERIFIKASI & STATUS", "items": [
            f"Sumber: commit {COMMIT_SHA[:8]}",
            "Test passing: 14 test cases, 0 failures.",
            "Race detector: clean (go test -race ./...)",
            "Simulasi in-memory string parsing; bukan proof production ADR governance completeness."
        ], "type": "bullet"}
    ],
    mental_model="Parser makes ADR machine-readable; linter enforces structural invariants. Human wisdom remains irreplaceable.",
    source_path=REPO_URL,
    slide_num="24 / 24"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
