#!/usr/bin/env python3
"""
Publication Builder for Software Engineering Lab 17: Architecture Decision Record.
Generates: PPTX, PDF, PNG covers, Dev.to MD, Medium HTML, Metadata, Caption, DELIVERABLES.md
Adheres strictly to the SLIDE CLEANUP OVERRIDE standard: concise, high readability, professional.
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

# Slide 1: Cover Slide (Dark Navy)
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
    title="Konteks Arsitektur Hilang Tanpa Jejak",
    subtitle="Alasan teknis di balik keputusan kritis sering hilang saat tim bertumbuh.",
    content_items=[
        {"title": "MASALAH NYATA", "items": [
            "Keputusan arsitektur diambil via chat atau meeting tanpa catatan formal.",
            "Pertanyaan berulang muncul di review kode: mengapa sistem dibangun begini?",
            "ADR lama tidak diperbarui, sehingga status dokumen berbeda dari kode aktual."
        ], "type": "bullet"}
    ],
    narrative="Keputusan arsitektur sering lahir di chat atau meeting. Saat tim berubah, alasan awalnya hilang. ADR membuat alasan itu bisa ditelusuri lagi.",
    mental_model="Architectural context without structural records becomes tribal knowledge.",
    source_path="content/02-master-draft.md",
    slide_num="02 / 24"
)

# Slide 3: Why Manual Markdown Review Fails
add_slide(
    categories=["WHY FAILS", "HUMAN REVIEW"],
    title="Review Manual Melewatkan Relasi Supersesi",
    subtitle="Reviewer manusia fokus pada isi teks, tetapi melewatkan sinkronisasi tautan dokumen.",
    content_items=[
        {"title": "KEGAGALAN REVIEW MANUAL", "items": [
            "ADR baru menyatakan menggantikan ADR lama, tetapi ADR lama tidak diubah.",
            "ADR lama tetap bertanda aktif padahal keputusannya sudah ditinggalkan.",
            "Tautan referensi menunjuk ke nomor yang salah tanpa disadari."
        ], "type": "bullet"}
    ],
    narrative="Masalah terbesar review manual adalah supersesi yang tidak sinkron. ADR baru mengklaim menggantikan, tetapi ADR lama tetap bertanda aktif karena lupa diperbarui.",
    mental_model="Manual review catches content typos but misses broken relationship links.",
    source_path="content/02-master-draft.md",
    slide_num="03 / 24"
)

# Slide 4: Mental Model — ADR as Structured Data
add_slide(
    categories=["MENTAL MODEL", "CORE CONCEPT"],
    title="ADR Bukan Dokumen, Tapi Record",
    subtitle="Parser memindahkan Markdown ke data terstruktur agar aturan bisa diperiksa mesin.",
    content_items=[
        {"title": "PRINSIP INTI", "items": [
            "ADR diperlakukan sebagai data terstruktur, bukan sekadar teks bebas.",
            "Field kunci diekstrak: ID, Judul, Status, SupersededBy, dan Supersedes.",
            "Mesin memvalidasi konsistensi; manusia tetap menulis esensi keputusannya."
        ], "type": "bullet"}
    ],
    narrative="ADR bukan sekadar dokumen naratif. Saat diubah menjadi struct Record, mesin bisa mengekstrak field kunci dan memeriksa konsistensi yang sering terlewat oleh manusia.",
    mental_model="ADR is structured data first, human narrative second.",
    source_path="internal/adr/models.go",
    slide_num="04 / 24"
)

# Slide 5: System Architecture — Parser to Linter Pipeline
add_slide(
    categories=["ARSITEKTUR", "DATA FLOW"],
    title="Pipeline Validasi Empat Tahap",
    subtitle="Alur konversi dokumen ADR dari teks mentah hingga laporan integritas.",
    content_items=[
        {"title": "VALIDATION PIPELINE", "items": [
            "Markdown Files → adr.Parser (ekstraksi field & status)",
            "→ []*adr.Record (in-memory struct)",
            "→ adr.Linter (cek urutan ID, siklus, & relasi timbal-balik)",
            "→ []error (laporan pelanggaran integritas)"
        ], "type": "diagram"}
    ],
    narrative="Pipeline-nya sederhana: Markdown dibaca, field penting diekstrak, lalu linter memeriksa aturan graf.",
    mental_model="Separate parsing from validation: simple stages create reliable checks.",
    source_path="internal/adr/parser.go • internal/adr/linter.go",
    slide_num="05 / 24"
)

# Slide 6: Key Code Path — Record Struct
add_slide(
    categories=["CODE", "MODEL"],
    title="Struct Record: Fondasi Validasi Data",
    subtitle="Field eksplisit membuat relasi antar-keputusan dapat diperiksa programmatically.",
    content_items=[
        {"title": "MODELS.GO", "items": [
            "type Record struct {",
            "    ID           int",
            "    Title        string",
            "    Status       Status",
            "    SupersededBy int // 0 jika aktif",
            "    Supersedes   int // 0 jika tidak menggantikan",
            "    Content      string",
            "}"
        ], "type": "code"}
    ],
    narrative="Struct Record mendefinisikan field wajib secara eksplisit. Dengan representasi integer dan enum, linter memeriksa aturan tanpa perlu menebak teks bebas.",
    mental_model="Explicit fields over implicit conventions enable automated linting.",
    source_path="internal/adr/models.go",
    slide_num="06 / 24"
)

# Slide 7: Failure Mode — Dangling Supersession Link
add_slide(
    categories=["FAILURE", "REGRESSION TEST"],
    title="Deteksi Tautan Supersesi Patah",
    subtitle="Linter langsung menolak referensi yang menunjuk ke nomor ADR yang tidak ada.",
    content_items=[
        {"title": "TEST CASE: broken_references", "items": [
            "ADR 1: Status Superseded by 99",
            "→ Error: superseded by non-existent ADR 99",
            "",
            "ADR 1: Status Accepted, Supersedes: 99",
            "→ Error: supersedes non-existent ADR 99"
        ], "type": "code"}
    ],
    narrative="Linter memeriksa keberadaan setiap ID referensi. Jika ADR merujuk ke nomor yang belum ada, error spesifik langsung dihasilkan sebelum merger.",
    mental_model="Dangling references are cheaper to catch at lint time than during incident review.",
    source_path="tests/linter_test.go:36-54",
    slide_num="07 / 24"
)

# Slide 8: Hidden Cost — Unverified Lineage Drift
add_slide(
    categories=["COST", "DOCUMENTATION DRIFT"],
    title="Biaya Tersembunyi: Lineage Drift Menumpuk",
    subtitle="Ketiadaan linter tidak merusak build, tetapi merusak keputusan masa depan.",
    content_items=[
        {"title": "DAMPAK DRIFT DOKUMEN", "items": [
            "Tahap | Deteksi | Dampak",
            "CI Gate | Seketika | PR ditolak, feedback loop cepat",
            "Review Manual | Mingguan | Tim bingung dokumen mana yang aktif",
            "Insiden Sistem | Berbulan-bulan | Keputusan baru berpijak pada asumsi keliru"
        ], "type": "table"}
    ],
    narrative="Masalah ADR rusak bukan langsung terlihat. Ia menumpuk sampai keputusan baru bergantung pada catatan lama yang salah.",
    mental_model="Lineage drift accumulates silently until a new decision depends on stale records.",
    source_path="content/02-master-draft.md",
    slide_num="08 / 24"
)

# Slide 9: Experiment Setup — Three ADR Scenarios
add_slide(
    categories=["EXPERIMENT", "DEMO SETUP"],
    title="Skenario Demo: Tiga ADR Evolusi Arsitektur",
    subtitle="Memodelkan evolusi dari monolit ke microservice dengan satu usulan ditolak.",
    content_items=[
        {"title": "URUTAN KEPUTUSAN", "items": [
            "ADR 0001: Modular Monolith | Superseded by ADR 2",
            "ADR 0002: Notification Microservice | Accepted, Supersedes ADR 1",
            "ADR 0003: Event Sourcing | Rejected (alasan biaya operasional)"
        ], "type": "bullet"}
    ],
    narrative="Eksperimen memvalidasi tiga ADR nyata: monolit awal, pemisahan service pengganti, dan usulan ditolak agar argumen penolakan tidak didebatkan ulang.",
    mental_model="Three records demonstrate full lifecycle: active, superseded, and rejected.",
    source_path="cmd/demo/main.go",
    slide_num="09 / 24"
)

# Slide 10: Result Snapshot — Demo Output
add_slide(
    categories=["RESULT", "EVIDENCE"],
    title="Output Demo: Validasi Integritas Lolos",
    subtitle="Ketiga ADR berhasil diparse dan aturan relasi timbal-balik terpenuhi.",
    content_items=[
        {"title": "TERMINAL OUTPUT", "items": [
            "$ go run ./cmd/demo",
            "Successfully parsed 3 ADRs:",
            "  - [0001] Modular Monolith -> Superseded by ADR 2",
            "  - [0002] Notification Microservice -> Supersedes ADR 1",
            "  - [0003] Event Sourcing -> Rejected",
            "Running ADR Integrity Linter... Integrity check passed!"
        ], "type": "code"}
    ],
    narrative="CLI memverifikasi rantai lineage: ADR 1 dan ADR 2 saling terhubung konsisten, nomor urut rapi, dan tidak ada siklus.",
    mental_model="Clean output verifies structural consistency across the entire record set.",
    source_path="cmd/demo/main.go",
    slide_num="10 / 24"
)

# Slide 11: Result Interpretation — What Passing Means
add_slide(
    categories=["INTERPRETASI", "EVIDENCE"],
    title="Lolos Linter Bukan Berarti Keputusan Benar",
    subtitle="Linter memverifikasi format data, bukan kebenaran arsitektur di dalamnya.",
    content_items=[
        {"title": "YANG DIBUKTIKAN LINTER", "items": [
            "Nomor ID berurutan tanpa celah.",
            "Relasi timbal-balik utuh tanpa tautan gantung.",
            "Tidak ada siklus ketergantungan antar-dokumen."
        ], "type": "bullet"},
        {"title": "TETAP PERLU REVIEW MANUSIA", "items": [
            "Kebenaran argumen teknis di bagian Decision.",
            "Apakah microservice memang solusi tepat bagi tim."
        ], "type": "bullet"}
    ],
    narrative="Linter hanya memastikan dokumen rapi dan tautannya konsisten. Ia tidak bisa menilai apakah keputusan microservice benar secara arsitektur. Kualitas keputusan tetap tanggung jawab tim.",
    mental_model="Linter guarantees structural validity; human judgment ensures architectural soundness.",
    source_path="engineering/03-execution-result.md",
    slide_num="11 / 24"
)

# Slide 12: Trade-off — Structure vs Flexibility
add_slide(
    categories=["TRADE-OFF", "DECISION"],
    title="Trade-off: Struktur Kaku vs Fleksibilitas",
    subtitle="Penegakan format memudahkan otomasi, namun membatasi gaya penulisan bebas.",
    content_items=[
        {"title": "STRUKTUR BAKU (LAB INI)", "items": [
            "Wajib tiga section: Context, Decision, Consequences.",
            "Status dibatasi lima enum resmi.",
            "Memungkinkan CI gate otomatis dan deterministik."
        ], "type": "bullet"},
        {"title": "FLEKSIBILITAS YANG DIKORBANKAN", "items": [
            "Tidak mendukung format paragraf naratif bebas.",
            "Memerlukan disiplin tim mengisi field relasi."
        ], "type": "bullet"}
    ],
    narrative="Format baku membuat dokumen bisa divalidasi mesin di CI. Biayanya: engineer tidak bisa lagi menulis catatan bebas tanpa mengikuti struktur.",
    mental_model="Enforced structure trades narrative flexibility for automated verification.",
    source_path="content/02-master-draft.md",
    slide_num="12 / 24"
)

# Slide 13: Improved Approach — Concurrent Validation with Mutex
add_slide(
    categories=["IMPROVED", "CONCURRENCY"],
    title="Validasi Paralel dengan Mutex",
    subtitle="Setiap record divalidasi di goroutine terpisah; error dikumpulkan aman.",
    content_items=[
        {"title": "CONCURRENT LINTER PATTERN", "items": [
            "for _, r := range records {",
            "    wg.Add(1)",
            "    go func(rec *Record) {",
            "        defer wg.Done()",
            "        localErrs := validate(rec)",
            "        mu.Lock()",
            "        errs = append(errs, localErrs...)",
            "        mu.Unlock()",
            "    }(r)",
            "}"
        ], "type": "code"}
    ],
    narrative="Validasi konkuren mengumpulkan error memakai mutex. Pola ini race-safe pada test, tetapi nilai utamanya ada pada kejelasan aturan, bukan kecepatan komputasi.",
    mental_model="Mutex-protected error aggregation keeps concurrent validation race-safe.",
    source_path="internal/adr/linter.go:64-106",
    slide_num="13 / 24"
)

# Slide 14: Before vs After — Manual Review vs Linter Gate
add_slide(
    categories=["KOMPARASI", "BEFORE AFTER"],
    title="Sebelum vs Sesudah Validasi Otomatis",
    subtitle="Celah relasi mekanis yang sering luput ditutup oleh pengecekan linter.",
    content_items=[
        {"title": "REVIEW MANUAL (SEBELUM)", "items": [
            "Tautan supersesi sepihak sering terlewat.",
            "Nomor meloncat sulit disadari saat banyak file.",
            "Status non-standar masuk tanpa sengaja."
        ], "type": "bullet"},
        {"title": "LINTER GATE (SESUDAH)", "items": [
            "Semua inkonsistensi ditolak otomatis di CI.",
            "Nomor ID wajib monotonik dan unik.",
            "Relasi timbal-balik divalidasi deterministik."
        ], "type": "bullet"}
    ],
    narrative="Linter menutup human error pada aspek mekanis: tautan sepihak, nomor loncat, dan status informal. Arsitek bisa fokus meninjau substansi keputusan.",
    mental_model="Automate mechanical relation checks to let reviewers focus on architectural depth.",
    source_path="tests/linter_test.go",
    slide_num="14 / 24"
)

# Slide 15: Boundary of the Lab — Format Validation Only
add_slide(
    categories=["BATASAN", "BOUNDARY"],
    title="Batasan Lab: Hanya Validasi Format",
    subtitle="Lab membuktikan parsing dan relasi; bukan kualitas keputusan arsitektur.",
    content_items=[
        {"title": "YANG DIUJI DI LAB", "items": [
            "Ekstraksi regex dari string Markdown.",
            "Pemeriksaan nomor urut dan relasi timbal-balik.",
            "Deteksi siklus pada graf supersesi."
        ], "type": "bullet"},
        {"title": "BATASAN BUKTI", "items": [
            "Tidak menilai kebenaran argumen keputusan.",
            "Simulasi in-memory, bukan pembacaan ribuan file disk."
        ], "type": "bullet"}
    ],
    narrative="Lab ini memvalidasi aturan dokumen, bukan kebenaran arsitektur. Parsing regex bekerja pada memori; pengujian tidak membuktikan seluruh perilaku repositori skala raksasa.",
    mental_model="Format validation ensures structural compliance, not business correctness.",
    source_path="engineering-audit/04-docs-vs-code.md",
    slide_num="15 / 24"
)

# Slide 16: Edge Case — Feature Cancellation Without ADR
add_slide(
    categories=["EDGE CASE", "DOCUMENTATION"],
    title="Kasus Batas: Menghentikan Fitur Tanpa ADR",
    subtitle="Keputusan yang ditinggalkan tanpa status Deprecated memicu ambiguitas sistem.",
    content_items=[
        {"title": "SKENARIO RISIKO", "items": [
            "Fitur dimatikan tetapi ADR tetap bertanda Accepted.",
            "Engineer baru mengira fitur masih menjadi standar aktif.",
            "Menghapus file justru menghilangkan konteks historis.",
            "Solusi: wajibkan ADR baru berstatus Deprecated."
        ], "type": "bullet"}
    ],
    narrative="Menghentikan fitur tanpa ADR Deprecated menciptakan gap antara catatan dan kode aktual. Dokumen lama tetap terlihat aktif dan membingungkan tim berikutnya.",
    mental_model="Feature deprecation requires an explicit ADR, not silent file deletion.",
    source_path="content/02-master-draft.md",
    slide_num="16 / 24"
)

# Slide 17: What the Tests Prove
add_slide(
    categories=["BUKTI", "TEST COVERAGE"],
    title="Bukti yang Divalidasi oleh Test Suite",
    subtitle="Suite pengujian membuktikan aturan struktural pada skenario yang dirancang.",
    content_items=[
        {"title": "FAKTA TEST SUITE", "items": [
            "Parser menolak status ilegal dan format tanpa judul.",
            "Linter menangkap 7 skenario referensi patah dan tautan gantung.",
            "Race detector bersih pada eksekusi pengujian konkuren."
        ], "type": "bullet"}
    ],
    narrative="Test membuktikan parser dan linter bekerja sesuai aturan pada skenario uji. Bukti ini menegaskan keandalan logika struktural, bukan jaminan semua kebutuhan tim terpenuhi.",
    mental_model="Tests verify structural logic against defined cases, not all team workflows.",
    source_path="tests/parser_test.go • tests/linter_test.go",
    slide_num="17 / 24"
)

# Slide 18: What the Tests Do Not Prove
add_slide(
    categories=["BATASAN BUKTI", "DO NOT OVERCLAIM"],
    title="Batas Integritas: Apa yang Tidak Diuji",
    subtitle="Test tidak menggantikan peninjauan arsitektur dan kebijakan kerja tim.",
    content_items=[
        {"title": "DI LUAR CAKUPAN PENGUJIAN", "items": [
            "Kebenaran konten teknis di dalam section Decision.",
            "Disiplin tim dalam memperbarui ADR saat arsitektur bergeser.",
            "Kesiapan sistem monolit atau microservice di production.",
            "Performa I/O disk pada ribuan berkas dokumen."
        ], "type": "bullet"}
    ],
    narrative="Test unit membuktikan aturan pada input tertentu. Keberhasilan test bukan bukti bahwa keputusan yang didokumentasikan sudah siap diterapkan di production.",
    mental_model="Passing tests confirm rule execution; production requires holistic architectural governance.",
    source_path="engineering-audit/05-gaps.md",
    slide_num="18 / 24"
)

# Slide 19: Production Risk — False Confidence from Green Linter
add_slide(
    categories=["RISIKO", "PRODUCTION"],
    title="Risiko Kepalsuan: Linter Hijau Bukan Aman",
    subtitle="Lolos validasi struktur dapat memberi rasa aman palsu jika argumen tidak dicek.",
    content_items=[
        {"title": "POTENSI KESALAHAN", "items": [
            "Dokumen lolos linter meski bagian Decision hanya berisi teks formalitas.",
            "Tim mengira arsitektur sudah matang hanya karena status CI hijau.",
            "Mitigasi: pertahankan review arsitek untuk menilai isi substansi."
        ], "type": "bullet"}
    ],
    narrative="Linter hijau hanya menjamin format valid, bukan kualitas desain. Jangan jadikan status lolos linter sebagai pengganti telaah arsitektur mendalam oleh engineer senior.",
    mental_model="Green linter verifies syntax; architectural soundness demands human scrutiny.",
    source_path="engineering-audit/06-verdict.md",
    slide_num="19 / 24"
)

# Slide 20: Implementation Checklist
add_slide(
    categories=["CHECKLIST", "IMPLEMENTASI"],
    title="Panduan Implementasi Linter ADR di Tim",
    subtitle="Langkah praktis mengintegrasikan validasi dokumen ke dalam workflow repositori.",
    content_items=[
        {"title": "LANGKAH KERJA", "items": [
            "[ ] Tetapkan template Markdown dengan Context, Decision, Consequences.",
            "[ ] Gunakan enum status resmi (Accepted, Superseded, Rejected, dll.).",
            "[ ] Pasang parser dan linter pada pull request CI gate.",
            "[ ] Larang edit retroaktif: keputusan baru wajib nomor ADR baru."
        ], "type": "checklist"}
    ],
    narrative="Implementasi sukses menggabungkan template terstandar dan CI gate otomatis. Aturan terpenting: jangan ubah ADR lama secara retroaktif; buat catatan baru jika arah berubah.",
    mental_model="Pair clear Markdown standards with automated CI gates to prevent drift.",
    source_path="engineering/02-implementation-notes.md",
    slide_num="20 / 24"
)

# Slide 21: Anti-Pattern — Common Mistakes When Adopting ADR
add_slide(
    categories=["ANTI-PATTERN", "LESSON"],
    title="Anti-Pattern: Kesalahan Adopsi ADR",
    subtitle="Kebiasaan buruk yang merusak manfaat catatan keputusan arsitektur.",
    content_items=[
        {"title": "HINDARI POLA INI", "items": [
            "Menulis ADR untuk hal sepele seperti penamaan variabel.",
            "Mengedit ADR Accepted diam-diam alih-alih membuat pengganti.",
            "Membiarkan status Draft berbulan-bulan tanpa kejelasan.",
            "Lupa mencatat nomor ADR pengganti saat membuat keputusan baru."
        ], "type": "bullet"}
    ],
    narrative="Anti-pattern paling berbahaya adalah suntingan diam-diam pada keputusan lama. Hal ini menghapus riwayat pembelajaran tim dan memicu asumsi keliru bagi engineer baru.",
    mental_model="Never overwrite historical decisions; publish a new record to maintain lineage.",
    source_path="content/02-master-draft.md",
    slide_num="21 / 24"
)

# Slide 22: Reproduction Guide
add_slide(
    categories=["REPRODUKSI", "COMMAND"],
    title="Reproduksi Mandiri di Mesin Lokal",
    subtitle="Eksperimen dijalankan dengan Go standard library tanpa dependency luar.",
    content_items=[
        {"title": "COMMAND RUNNER", "items": [
            "cd labs/17-architecture-decision-record",
            "go test -v ./...    # validasi parser & linter",
            "go test -race ./... # cek race condition",
            "go run ./cmd/demo   # jalankan simulasi 3 ADR"
        ], "type": "code"}
    ],
    narrative="Lab ini zero-dependency dan dapat direproduksi langsung. Menjalankan demo dan test membuktikan parser dan linter bekerja secara mandiri di lokal.",
    mental_model="Zero-dependency tests ensure fast, repeatable verification on any developer machine.",
    source_path="README.md • go.mod • tests/",
    slide_num="22 / 24"
)

# Slide 23: Technical Takeaway
add_slide(
    categories=["TAKEAWAY", "ENGINEERING PRINCIPLE"],
    title="ADR Diubah Jadi Data yang Bisa Dicek",
    subtitle="Memindahkan keputusan arsitektur dari teks bebas ke struktur yang terverifikasi.",
    content_items=[
        {"title": "PRINSIP REKAYASA", "items": [
            "ADR terstruktur memungkinkan verifikasi otomatis relasi dokumen.",
            "Linter menutup kelalaian mekanis yang sering luput dari review manusia.",
            "Format standar membantu review; kualitas keputusan tetap ditentukan isi ADR."
        ], "type": "bullet"}
    ],
    narrative="Nilai utama lab ini adalah menjadikan keputusan arsitektur sebagai data yang bisa diuji. Linter menjaga integritas silsilah dokumen, sementara manusia fokus pada kualitas keputusan.",
    mental_model="Parser makes ADR structured; linter enforces rules. Decisions remain human.",
    source_path="content/05-key-takeaways.md",
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
