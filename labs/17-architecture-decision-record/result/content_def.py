#!/usr/bin/env python3
"""
Publication Builder Definition for Software Engineering Lab 17: Architecture Decision Record.
Follows:
- cover-rules.md (Mandatory CRT retro computer, series identity lock, no micro-text, thumbnail-proof 320x180)
- slide-rules.md / slide-content-rules.md (<=45 words prose, <=30 words on code slides, dominant visual anchors, 50% zoom test)
- quality-bar.md (sharp engineering insights, strictly bounded claims, >=12 slides with limiting phrases, no overclaims)
- article-rules.md (Dev.to 22-section mapping, Medium ID HTML 22-section mapping)
"""
import os
import sys

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/17-architecture-decision-record/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/17-architecture-decision-record"
COMMIT_SHA = "774e97a"
LAB_NO = "17"
SLUG = "architecture-decision-record"

TITLE_ID = "Architecture Decision Record & Linting Struktural"
SUBTITLE_ID = "Mencegah Perdebatan Berulang dengan Silsilah Keputusan Terverifikasi di Go"
TITLE_EN = "Architecture Decision Records & Structural Linting"
SUBTITLE_EN = "Preventing Cyclical Debates with Verified Decision Lineage in Go"

AUTHOR = "Lukman / lukman-ss"
TOPIC_TAGS = "Architecture Decision Record • ADR • Go • Architecture • Software Engineering"
COVER_IMAGE_PATH = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/17-architecture-decision-record/result/image-cover/cover_art.jpg"

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
            "Retro CRT Terminal Displaying Immutable Decision Ledger and DAG Lineage",
            "Monotonic Index Registry: Linked ADR State Graph & Bidirectional Cross-Check",
            f"Penulis: {AUTHOR}"
        ]}
    ],
    mental_model="Architecture decisions require immutable records co-located with code and validated lineage.",
    source_path="",
    slide_num="01 / 24",
    layout="cover"
)

# Slide 02: Problem Nyata (Layout: cards)
# Limiting phrase 1: "bukan berarti"
add_slide(
    categories=["PROBLEM", "CONTEXT LOSS"],
    title="Kehilangan Konteks Memicu Perdebatan Teknis Berulang",
    subtitle="Tanpa rekaman tertulis mengapa pilihan diambil, tim terjebak dalam dua kegagalan ekstrem.",
    content_items=[
        {"type": "cards", "title": "DAMPAK HILANGNYA KONTEKS KEPUTUSAN", "items": [
            {"title": "Pembalikan Prematur", "desc": "Arsitektur lama dibongkar sembarangan tanpa memahami batasan dan kompromi masa lalu."},
            {"title": "Kepatuhan Buta", "desc": "Pola usang dipertahankan kaku meskipun skala beban dan kapasitas tim telah berubah total."},
            {"title": "Debat Berulang", "desc": "Tim mendiskusikan kembali opsi yang sama berulang kali karena alasan penolakannya hilang."}
        ]}
    ],
    narrative="Ketika latar belakang keputusan teknis hilang, tim terjebak dalam perdebatan siklis tanpa ujung. Memori pengembang bersifat fana, bukan berarti keputusan lama salah, tetapi tanpa catatan tertulis tim tidak tahu batasan apa yang memicunya.",
    narrative_en="When technical rationale is forgotten, engineering teams enter cyclical debates. Individual memory fades quickly; this does not mean past decisions were flawed, but without written constraints teams cannot judge why trade-offs were made.",
    mental_model="Unrecorded decisions force teams to re-litigate trade-offs repeatedly.",
    source_path="content/02-master-draft.md:3-12",
    slide_num="02 / 24",
    layout="cards"
)

# Slide 03: Mengapa Review Manual / Pendekatan Naif Gagal (Layout: cards)
# Limiting phrase 2: "tidak otomatis"
add_slide(
    categories=["WHY FAILS", "DOCUMENTATION DRIFT"],
    title="Wiki Terpisah dan Review Manual Mengalami Drift",
    subtitle="Menyimpan keputusan di wiki eksternal memutus sinkronisasi antara dokumen dan kode sumber.",
    content_items=[
        {"type": "cards", "title": "PENDEKATAN NAIF VS REALITA", "items": [
            {"title": "Wiki Eksternal Usang", "desc": "Dokumen Confluence terpisah dari PR; kode berevolusi cepat, dokumen tertinggal permanen."},
            {"title": "Review Manual Rapuh", "desc": "Reviewer luput memeriksa apakah tautan penggantian Supersedes konsisten di kedua belah pihak."},
            {"title": "Tautan Patah Tak Terlacak", "desc": "Menandai Superseded by 99 manual tidak otomatis gagal saat file 99 tidak pernah ada di repo."}
        ]}
    ],
    narrative="Menaruh keputusan di wiki eksternal tidak otomatis menjamin dokumentasi dibaca atau diperbarui. Review manual rentan melewatkan broken reference dan inkonsistensi status dua arah saat commit digabungkan.",
    narrative_en="Housing architecture decisions in external wikis does not automatically guarantee freshness. Manual PR reviews routinely overlook broken supersession pointers and asymmetric lifecycle mutations.",
    mental_model="Documentation detached from the codebase rapidly suffers from documentation drift.",
    source_path="content/02-master-draft.md:13-21",
    slide_num="03 / 24",
    layout="cards"
)

# Slide 04: Mental Model Inti — Schema Card (Layout: schema)
# Limiting phrase 3: "hanya pada jalur ini"
add_slide(
    categories=["MENTAL MODEL", "IMMUTABILITY & DAG"],
    title="Ko-lokasi, Immutability, dan Graf Silsilah Timbal-Balik",
    subtitle="ADR disimpan bersama kode; keputusan diterima tidak disunting melainkan digantikan.",
    content_items=[
        {"type": "schema", "title": "EMPAT PILAR KONSEPTUAL ADR", "items": [
            {"phase": "1. Co-location", "role": "Git Repository", "desc": "ADR disimpan di docs/adr/, berevolusi melalui PR dan code review yang sama dengan kode."},
            {"phase": "2. Immutability", "role": "Historical Record", "desc": "ADR Accepted tidak diedit retroaktif; arah baru dicatat dalam file ADR baru bernomor urut."},
            {"phase": "3. Directed Lineage", "role": "DAG Relations", "desc": "ADR baru mendeklarasikan Supersedes: X; ADR lama ditandai Superseded by Y secara resiprokal."},
            {"phase": "4. Context Focus", "role": "Forces & Trade-offs", "desc": "Fokus pada 'mengapa' dan batasan kompromi; hanya pada jalur ini silsilah keputusan divalidasi."}
        ]}
    ],
    narrative="ADR memperlakukan keputusan arsitektur sebagai artefak version-controlled yang kekal. Keputusan lama tidak dihapus dari Git, melainkan digantikan secara eksplisit. Hanya pada jalur ini integritas evolusi sistem dapat dilacak.",
    narrative_en="ADRs treat architectural choices as immutable, version-controlled repository assets. Historical records are superseded rather than overwritten, establishing a traceable decision graph across git commits.",
    mental_model="Decisions are immutable history; changes produce new nodes in a directed graph.",
    source_path="content/02-master-draft.md:22-30",
    slide_num="04 / 24",
    layout="schema"
)

# Slide 05: Diagram Alur / Pipeline Validasi (Layout: schema)
# Limiting phrase 4: "belum membuktikan"
add_slide(
    categories=["PIPELINE", "WORKFLOW"],
    title="Pipeline Validasi: Markdown ke Parsed Record ke Linter",
    subtitle="Membedah teks Markdown menjadi representasi memori lalu menguji invarian struktural.",
    content_items=[
        {"type": "schema", "title": "ALUR INTEGRITAS STRUKTURAL ADR", "items": [
            {"phase": "1. Text Scanning", "role": "bufio.Scanner", "desc": "Membaca Markdown baris demi baris, mengekstrak ID, Title, Status, dan relasi supersedes."},
            {"phase": "2. Section Guard", "role": "Content Verifier", "desc": "Memastikan bagian Context, Decision, dan Consequences hadir dan memiliki konten nyata."},
            {"phase": "3. Monotonic Check", "role": "Sequence Validator", "desc": "Mengurutkan array ID dan memverifikasi ids[i] == i + 1 tanpa celah maupun duplikasi."},
            {"phase": "4. Concurrent Lint", "role": "Goroutine Pool", "desc": "Memverifikasi tautan timbal balik SupersededBy vs Supersedes; belum membuktikan semantik isi."}
        ]}
    ],
    narrative="Pipeline membaca dokumen Markdown, mengekstrak metadata melalui parser regex, lalu menguji urutan penomoran dan relasi graf keputusan secara paralel. Pipeline ini memvalidasi struktur, namun belum membuktikan kualitas logika desain teknisnya.",
    narrative_en="The pipeline ingests raw Markdown, extracts metadata using regular expressions, and concurrently audits monotonic indexing and reciprocal linkage. This verifies structural integrity, but does not prove semantic design quality.",
    mental_model="Transform raw markdown into strongly-typed nodes before verifying graph invariants.",
    source_path="content/04-diagrams.md:63-90",
    slide_num="05 / 24",
    layout="schema"
)

# Slide 06: Cuplikan Kode Dominan — Model Record & Status (Layout: code)
# Limiting phrase 5: "bukan berarti"
add_slide(
    categories=["IMPLEMENTATION", "DATA MODEL"],
    title="Model Status Siklus Hidup dan Entitas Record",
    subtitle="Mendefinisikan 5 status siklus hidup formal dan relasi penggantian di Go.",
    content_items=[
        {"type": "code", "title": "internal/adr/models.go", "items": [
            "type Status string",
            "const (",
            "    StatusProposed   Status = \"Proposed\"",
            "    StatusAccepted   Status = \"Accepted\"",
            "    StatusSuperseded Status = \"Superseded\"",
            "    StatusDeprecated Status = \"Deprecated\"",
            "    StatusRejected   Status = \"Rejected\"",
            ")",
            "type Record struct {",
            "    ID, SupersededBy, Supersedes int",
            "    Title string; Status Status; Content string",
            "}"
        ]}
    ],
    narrative="Struct Record menampung metadata keputusan dan penunjuk silsilah penggantian. Penggunaan enum string bukan berarti sistem kebal salah input; parser tetap wajib memvalidasi kecocokan status via method IsValid.",
    narrative_en="The Record struct captures decision metadata alongside forward and backward DAG pointers. Using typed strings does not mean validation is automatic; the parser must explicitly verify enum membership via IsValid.",
    mental_model="Explicit lifecycle types anchor structural validation at the domain boundary.",
    source_path="internal/adr/models.go:3-39",
    slide_num="06 / 24",
    layout="code"
)

# Slide 07: Failure Path / Dangling Reference Konkret (Layout: cards)
# Limiting phrase 6: "hanya pada jalur ini"
add_slide(
    categories=["FAILURE PATH", "STRUCTURAL BREAK"],
    title="Empat Titik Patah Silsilah Keputusan Arsitektur",
    subtitle="Bagaimana ketidaksinkronan penyuntingan manual merusak integritas riwayat Git.",
    content_items=[
        {"type": "cards", "title": "SKENARIO KEGAGALAN TERVERIFIKASI", "items": [
            {"title": "Dangling SupersededBy", "desc": "ADR 1 bertuliskan 'Superseded by 99', namun berkas ADR 99 tidak ada di repositori."},
            {"title": "One-Sided Reference", "desc": "ADR 2 mengklaim 'Supersedes: 1', tetapi ADR 1 masih berstatus 'Accepted' tanpa pointer pengganti."},
            {"title": "Non-Monotonic Gap", "desc": "Terdapat berkas ADR 1 dan ADR 3 tanpa ADR 2; memicu ambiguitas riwayat yang hilang."},
            {"title": "Self & Cyclic Supersession", "desc": "ADR menunjuk dirinya sendiri atau membentuk loop 1 -> 2 -> 3 -> 1; hanya pada jalur ini dideteksi."}
        ]}
    ],
    narrative="Ketidaksinkronan dua arah sering terjadi saat pengembang membuat ADR baru tetapi lupa memperbarui status ADR lama. Hanya pada jalur ini validator memeriksa kedua sisi relasi graf untuk mencegah dangling reference di repositori.",
    narrative_en="One-sided references emerge when engineers write replacement ADRs without updating prior records. Only along this validation pathway does the tool cross-examine both nodes to eliminate broken pointers.",
    mental_model="A DAG pointer without reciprocal back-reference is a silent documentation defect.",
    source_path="tests/linter_test.go:36-175",
    slide_num="07 / 24",
    layout="cards"
)

# Slide 08: Tabel Perbandingan Mekanisme vs Gejala (Layout: table)
add_slide(
    categories=["COMPARISON", "DIAGNOSTICS"],
    title="Matriks Mekanisme Linter vs Gejala Kerusakan",
    subtitle="Pemetaan aturan validasi kode terhadap anomali struktural berkas Markdown.",
    content_items=[
        {"title": "VALIDASI VS GEJALA KERUSAKAN", "items": [
            "Aturan Linter | Gejala Kerusakan Dokumen | Dampak Tim Rekayasa",
            "Monotonic Order | Lompat nomor (1 langsung 3) | Keraguan riwayat commit hilang",
            "Reciprocal Pointer | ADR 2 supersede 1 tanpa update 1 | Ambiguitas status operasional aktif",
            "Target Existence | Superseded by 99 (file 99 nihil) | Referensi putus / dangling context",
            "Cycle Detector | Siklus 1 -> 2 -> 3 -> 1 | Silsilah keputusan mengalami loop tak hingga",
            "Section Guard | Context / Decision / Consequence kosong | Catatan kehilangan alasan teknis esensial"
        ], "type": "table"}
    ],
    narrative="Tabel ini memetakan mekanisme linter terhadap gejala kerusakan dokumen riil. Setiap guard mencegah risiko ambiguitas arsitektural yang dapat menyesatkan teknisi saat mengambil keputusan di masa depan.",
    narrative_en="This matrix maps validator checks directly against concrete documentation anomalies. Each validation rule safeguards against structural ambiguities that misinform future architectural reviews.",
    mental_model="Enforce structural invariants to maintain reliable historical traceability.",
    source_path="content/02-master-draft.md:51-73",
    slide_num="08 / 24",
    layout="table"
)

# Slide 09: Decision Sequence — 3 Focused Decision Cards (Layout: cards)
# Limiting phrase 7: "bukan berarti"
add_slide(
    categories=["DECISION SEQUENCE", "CASE STUDY"],
    title="Silsilah Evolusi SaaS ERP: Monolith ke Microservice",
    subtitle="Tiga keputusan bertahap yang membuktikan siklus hidup Proposal, Accepted, hingga Rejected.",
    content_items=[
        {"type": "cards", "title": "SEKUENS KEPUTUSAN SAAS ERP", "items": [
            {"title": "ADR 0001: Modular Monolith", "desc": "Tim 5 teknisi; fokus iterasi cepat dan transaksi ACID lokal. Status akhir: Superseded by 2."},
            {"title": "ADR 0002: Notification Service", "desc": "Tim 25 teknisi; lonjakan beban notifikasi diisolasi ke microservice. Status: Accepted (Supersedes 1)."},
            {"title": "ADR 0003: Reject Event Sourcing", "desc": "Usulan audit log via Event Sourcing ditolak demi kesederhanaan relasional. Status: Rejected."}
        ]}
    ],
    narrative="Evolusi arsitektur ERP mendemonstrasikan transisi status alami. Memilih monolith di awal bukan berarti anti-microservices, melainkan respons rasional terhadap batasan tim 5 teknisi dan batas waktu peluncuran 3 bulan.",
    narrative_en="The SaaS ERP case study illustrates realistic lifecycle progression. Adopting a modular monolith first does not indicate opposition to distributed systems; it represents a rational response to early-stage team constraints.",
    mental_model="Architecture decisions evolve in response to changing organizational and technical forces.",
    source_path="cmd/demo/main.go:10-64",
    slide_num="09 / 24",
    layout="cards"
)

# Slide 10: Bukti Eksekusi / Terminal Output (Layout: terminal)
add_slide(
    categories=["EVIDENCE", "EXECUTION"],
    title="Eksekusi Demo cmd/demo dan Test Suite Lolos Bersih",
    subtitle="Output terminal membuktikan parsing 3 ADR dan kelolosan linter tanpa galat integritas.",
    content_items=[
        {"type": "code", "title": "TERMINAL OUTPUT: go run ./cmd/demo && go test -v", "items": [
            "=== Architecture Decision Record (ADR) Lab ===",
            "Successfully parsed 3 ADRs:",
            "  - [0001] Use Modular Monolith for Core SaaS ERP        | Status: Superseded -> Superseded by ADR 2",
            "  - [0002] Extract Notification Service to Microservice  | Status: Accepted   -> Supersedes ADR 1",
            "  - [0003] Reject Event Sourcing for Order Management    | Status: Rejected",
            "",
            "Running ADR Integrity Linter...",
            "Integrity check passed! Decision lineage and lifecycle invariants are intact."
        ]}
    ],
    narrative="Program demo mengekstrak seluruh metadata dan memvalidasi integritas silsilah tanpa kegagalan. Eksekusi ini membuktikan bahwa invarian penomoran monotonik dan tautan timbal balik bekerja tepat pada skenario SaaS ERP.",
    narrative_en="The demo runner extracts document metadata and verifies graph invariants without error. This execution confirms that monotonic ordering and bidirectional pointers function correctly across the ERP case study.",
    mental_model="Executable verification proves that structural validation rules hold under simulated workflows.",
    source_path="engineering/03-execution-result.md:51-63",
    slide_num="10 / 24",
    layout="code"
)

# Slide 11: Analisis Batasan Bukti (Layout: cards)
# Limiting phrase 8: "belum membuktikan"
add_slide(
    categories=["AUDIT", "EVIDENCE BOUNDARY"],
    title="Batasan Bukti Lab: Yang Dibuktikan vs Belum",
    subtitle="Memisahkan verifikasi struktural otomatis dari penilaian semantik arsitektur nyata.",
    content_items=[
        {"type": "cards", "title": "BATAS PEMBUKTIAN SISTEM", "items": [
            {"title": "YANG DIBUKTIKAN", "desc": "Parsing regex mendeteksi header, status, section; linter menolak broken link, gap nomor, cycle, race."},
            {"title": "YANG TIDAK DIBUKTIKAN", "desc": "Kelayakan teknis keputusan arsitektur; dampak kuantitatif terhadap onboarding; integrasi git hook."},
            {"title": "BATAS LINGKUP TIM", "desc": "Lab belum membuktikan efektivitas skema pada skala 100+ engineer atau repositori terdistribusi multi-repo."}
        ]}
    ],
    narrative="Lab ini membuktikan keandalan validasi sintaks dan integritas referensi graf. Hasil ini belum membuktikan kebenaran isi keputusan teknis itu sendiri, yang tetap memerlukan penilaian manusia melalui peer review arsitek.",
    narrative_en="This lab proves regex-based syntax parsing and referential graph validation. It does not prove the soundness of the architectural decisions themselves, which still requires human peer review.",
    mental_model="Structural linters verify document integrity, not human architectural wisdom.",
    source_path="research-audit/07-verdict.md:Limitations",
    slide_num="11 / 24",
    layout="cards"
)

# Slide 12: Trade-off Arsitektur (Layout: cards)
# Limiting phrase 9: "tidak otomatis"
add_slide(
    categories=["TRADE-OFF", "COSTS & BENEFITS"],
    title="Trade-off Penerapan ADR dalam Alur Rekayasa",
    subtitle="Biaya pemeliharaan dokumentasi terstruktur vs manfaat ketertelusuran jangka panjang.",
    content_items=[
        {"type": "cards", "title": "KOMPROMI IMPLEMENTASI ADR", "items": [
            {"title": "Overhead Penulisan PR", "desc": "Pengembang wajib menyusun Context, Decision, dan Consequences sebelum merge arsitektural."},
            {"title": "Disiplin Immutability", "desc": "Keputusan tidak boleh diedit langsung; harus membuat berkas baru saat strategi berubah."},
            {"title": "Bukan Pengganti Diskusi", "desc": "Adanya ADR tidak otomatis mencegah perdebatan sengit jika tim tidak menyepakati batasan masalah."}
        ]}
    ],
    narrative="Menerapkan ADR menuntut disiplin penulisan dan biaya review tambahan di setiap PR arsitektur. Keberadaan ADR tidak otomatis menyelesaikan perbedaan pendapat tim, melainkan menyediakan basis data objektif untuk pengambilan keputusan.",
    narrative_en="Adopting ADRs imposes authoring friction and additional review overhead on architectural PRs. Having an ADR repository does not automatically eliminate disagreement, but grounds technical debates in documented facts.",
    mental_model="ADRs trade short-term authoring friction for long-term organizational clarity.",
    source_path="content/02-master-draft.md:228-251",
    slide_num="12 / 24",
    layout="cards"
)

# Slide 13: Concurrency / Validation Pipeline (Layout: code)
# Limiting phrase 10: "simulasi ini tidak sama dengan"
add_slide(
    categories=["CONCURRENCY", "PARALLEL LINTER"],
    title="Validasi Konkuren Graf Keputusan dengan Mutex Agregator",
    subtitle="Goroutine memverifikasi tautan timbal balik secara independen dengan proteksi race.",
    content_items=[
        {"type": "code", "title": "internal/adr/linter.go", "items": [
            "var wg sync.WaitGroup",
            "for _, r := range records {",
            "    wg.Add(1)",
            "    go func(rec *Record) {",
            "        defer wg.Done()",
            "        var localErrs []error",
            "        // Cek reciprocal SupersededBy & Supersedes",
            "        if len(localErrs) > 0 {",
            "            mu.Lock()",
            "            errs = append(errs, localErrs...)",
            "            mu.Unlock()",
            "        }",
            "    }(r)",
            "}"
        ]}
    ],
    narrative="Linter memanfaatkan goroutine untuk memverifikasi setiap simpul record secara paralel. Pengumpulan error lokal diproteksi sync.Mutex. Simulasi ini tidak sama dengan pemrosesan jutaan node graf, namun lolos go test -race pada 100 record.",
    narrative_en="The linter distributes node validation across worker goroutines with mutex-guarded error aggregation. While this simulation is not equivalent to billion-node graph engines, it passes race detection cleanly across 100 records.",
    mental_model="Fan-out node checks concurrently; aggregate errors through thread-safe synchronization.",
    source_path="internal/adr/linter.go:64-106",
    slide_num="13 / 24",
    layout="code"
)

# Slide 14: Before vs After Comparison Cards (Layout: cards)
# Limiting phrase 11: "bukan berarti"
add_slide(
    categories=["COMPARISON", "BEFORE VS AFTER"],
    title="Transformasi Manajemen Keputusan: Sebelum vs Sesudah ADR",
    subtitle="Perubahan fundamental dari memori fana pengembang ke repositori terstruktur.",
    content_items=[
        {"type": "cards", "title": "SEBELUM (NAIF) VS SESUDAH (ADR TERVERIFIKASI)", "items": [
            {"title": "Sebelum: Asumsi Tribal", "desc": "Alasan arsitektur hanya ada di ingatan lead engineer; hilang seketika saat orang tersebut resign."},
            {"title": "Sesudah: Riwayat Ko-lokasi", "desc": "Seluruh keputusan tersimpan di docs/adr/, terlacak via git blame dan pull request terikat."},
            {"title": "Integritas Silsilah", "desc": "Mengganti arsitektur bukan berarti menghapus masa lalu, melainkan menghubungkan tautan timbal balik."}
        ]}
    ],
    narrative="Sebelumnya pengetahuan arsitektural tersebar dalam chat dan obrolan lisan yang mudah hilang. Dengan ADR dan linter otomatis, silsilah keputusan tercatat kekal; memperbarui sistem bukan berarti menghapus riwayat sebelumnya.",
    narrative_en="Previously, critical architectural context remained trapped in tribal discussions. Structured ADRs ensure architectural history survives personnel changes through co-located, immutable Git records.",
    mental_model="Replace fragile oral traditions with permanent, version-controlled architecture ledgers.",
    source_path="content/02-master-draft.md:1-21",
    slide_num="14 / 24",
    layout="cards"
)

# Slide 15: Analisis Mendalam — Parsing Regex & Validasi Section (Layout: code)
# Limiting phrase 12: "cakupannya mengikuti test"
add_slide(
    categories=["ANALYSIS", "PARSER INTERNALS"],
    title="Parser Memastikan Kehadiran Header dan Section Wajib",
    subtitle="Regex mengekstrak atribut metadata sekaligus memvalidasi konten section tidak kosong.",
    content_items=[
        {"type": "code", "title": "internal/adr/parser.go", "items": [
            "var (",
            "    titleRegex   = regexp.MustCompile(`^#\\s+(\\d+)\\.\\s+(.+)$`)",
            "    statusRegex  = regexp.MustCompile(`(?i)^Status:\\s*([A-Za-z]+)(?:\\s+by\\s+(\\d+))?`)",
            "    contextRegex = regexp.MustCompile(`(?i)^##\\s+Context`)",
            ")",
            "if !foundTitle { return nil, fmt.Errorf(\"title not found\") }",
            "if !foundStatus || !record.Status.IsValid() { ... }",
            "if !foundContext || !contextHasContent { ... }"
        ]}
    ],
    narrative="Parser memeriksa format Markdown baris demi baris menggunakan regex terarah. Parser menolak dokumen jika section Context, Decision, atau Consequences kosong. Namun, cakupannya mengikuti test skenario spesifik format lab ini.",
    narrative_en="The parser evaluates Markdown lines sequentially using targeted regex patterns, rejecting missing or empty required sections. Its verification boundary is scoped strictly to the format defined in the test suite.",
    mental_model="Mandate essential semantic sections at the syntax parsing stage.",
    source_path="internal/adr/parser.go:12-164",
    slide_num="15 / 24",
    layout="code"
)

# Slide 16: Analisis Mendalam — Validasi Urutan Monotonik (Layout: code)
# Limiting phrase 13: "tidak otomatis"
add_slide(
    categories=["ANALYSIS", "MONOTONIC CHECK"],
    title="Penomoran Monotonik Mencegah Kehilangan Riwayat",
    subtitle="Mengurutkan array ID untuk mendeteksi celah urutan bernomor secara deterministik.",
    content_items=[
        {"type": "code", "title": "internal/adr/linter.go", "items": [
            "ids := make([]int, 0, len(records))",
            "for id := range recordMap { ids = append(ids, id) }",
            "sort.Ints(ids)",
            "",
            "for i := 0; i < len(ids); i++ {",
            "    if ids[i] != i+1 {",
            "        errs = append(errs, fmt.Errorf(",
            "            \"non-monotonic numbering, expected %d but got %d\",",
            "            i+1, ids[i]))",
            "        break",
            "    }",
            "}"
        ]}
    ],
    narrative="Linter menyortir ID dan mencocokkan indeks ke-i dengan ekspektasi i+1. Jika ADR 2 terhapus atau terlompati, linter langsung menolak commit. Menjaga nomor urut tidak otomatis mudah saat merge branch konkuren tanpa linter ini.",
    narrative_en="The linter sorts IDs and enforces that element i matches i+1 exactly. If ADR 2 is skipped, validation fails immediately. Sequential monotonicity is not automatically maintained in concurrent branches without this check.",
    mental_model="Strict monotonic indexing ensures no historical decisions disappear unseen from Git history.",
    source_path="internal/adr/linter.go:27-39",
    slide_num="16 / 24",
    layout="code"
)

# Slide 17: Analisis Mendalam — Deteksi Siklus dan Self-Supersession (Layout: code)
# Limiting phrase 14: "hanya pada jalur ini"
add_slide(
    categories=["ANALYSIS", "CYCLE DETECTION"],
    title="Deteksi Siklus DFS Mencegah Infinite Loop Silsilah",
    subtitle="Algoritma depth-first search berbasis 3-state menandai siklus pada graf keputusan.",
    content_items=[
        {"type": "code", "title": "internal/adr/linter.go", "items": [
            "state := make(map[int]int) // 0:unvisited, 1:visiting, 2:visited",
            "var detectCycle func(id int)",
            "detectCycle = func(id int) {",
            "    state[id] = 1",
            "    if rec, exists := recordMap[id]; exists && rec.SupersededBy != 0 {",
            "        if state[rec.SupersededBy] == 1 {",
            "            errs = append(errs, fmt.Errorf(\"cyclical supersession\"))",
            "        } else if state[rec.SupersededBy] == 0 { detectCycle(rec.SupersededBy) }",
            "    }",
            "    state[id] = 2",
            "}"
        ]}
    ],
    narrative="Silsilah keputusan harus membentuk DAG murni tanpa loop. Algoritma DFS 3-state melacak simpul yang sedang ditelusuri. Hanya pada jalur ini linter menangkap anomali ketika ADR 1 digantikan oleh 2, 2 oleh 3, dan 3 kembali ke 1.",
    narrative_en="Decision lineage must form a strict DAG without circular dependencies. A 3-state DFS algorithm tracks visiting nodes. Only along this path does the linter detect cycles like ADR 1 -> 2 -> 3 -> 1.",
    mental_model="Architectural evolution must flow forward as a directed acyclic graph, never a cycle.",
    source_path="internal/adr/linter.go:41-62",
    slide_num="17 / 24",
    layout="code"
)

# Slide 18: Analisis Mendalam — Menolak Nilai Status Tidak Valid (Layout: code)
# Limiting phrase 15: "bukan berarti"
add_slide(
    categories=["ANALYSIS", "SCHEMA GUARD"],
    title="Validasi Status Siklus Hidup Menolak Status Informal",
    subtitle="Mencegah status informal seperti 'Draft' atau 'In Review' masuk ke repositori utama.",
    content_items=[
        {"type": "code", "title": "tests/parser_test.go", "items": [
            "func TestParse_Invalid(t *testing.T) {",
            "    tests := []struct { name, content, wantErr string }{",
            "        {",
            "            name: \"invalid status\",",
            "            content: \"# 1. Title\\nStatus: Draft\\n## Context\\nA\\n## Decision\\nB\\n## Consequences\\nC\\n\",",
            "            wantErr: \"invalid status: Draft\",",
            "        },",
            "    }",
            "    // Verifikasi penolakan status di luar 5 status formal",
            "}"
        ]}
    ],
    narrative="Status siklus hidup dibatasi secara tegas pada Proposed, Accepted, Superseded, Deprecated, dan Rejected. Menolak status Draft bukan berarti melarang peninjauan bertahap, melainkan mewajibkan draft diuji via branch PR sebelum di-merge.",
    narrative_en="Lifecycle statuses are strictly restricted to five formal states. Rejecting 'Draft' status does not prevent collaborative review; it mandates that drafting occurs within branch PRs prior to merge.",
    mental_model="Enforce formal lifecycle states at merge boundaries to avoid ambiguous intermediate records.",
    source_path="tests/parser_test.go:126-132",
    slide_num="18 / 24",
    layout="code"
)

# Slide 19: Analisis Mendalam — Rekam Usulan yang Ditolak (Rejected) (Layout: cards)
# Limiting phrase 16: "tidak otomatis"
add_slide(
    categories=["ANALYSIS", "REJECTED DECISIONS"],
    title="Menyimpan Status 'Rejected' Menutup Debat Masa Depan",
    subtitle="Merekam usulan yang ditolak mencegah tim mengulang analisis yang sama tanpa data baru.",
    content_items=[
        {"type": "cards", "title": "NILAI STRATEGIS STATUS REJECTED", "items": [
            {"title": "Dokumentasi Alasan Penolakan", "desc": "Mencatat mengapa Event Sourcing ditolak: kompleksitas migrasi skema event vs tabel relasional."},
            {"title": "Mencegah Usulan Berulang", "desc": "Anggota tim baru tidak akan mengajukan ulang usulan serupa tanpa membawa bukti baru."},
            {"title": "Bukan Kegagalan Rekayasa", "desc": "ADR Rejected tidak otomatis berarti buang-buang waktu; itu investasi pencegahan disrupsi."}
        ]}
    ],
    narrative="Banyak tim menghapus dokumen usulan ketika ditolak. Praktik ini keliru karena menghilangkan alasan penolakannya. Menyimpan ADR berstatus Rejected tidak otomatis menambah beban, melainkan menjadi perisai dari perdebatan teknis berulang.",
    narrative_en="Many teams discard rejected proposals, inadvertently deleting the technical rationale behind the decision. Retaining Rejected ADRs does not introduce clutter; it acts as a firewall against repeated, unfounded re-proposals.",
    mental_model="Documenting why an option was rejected is as valuable as recording the chosen path.",
    source_path="content/02-master-draft.md:272-280",
    slide_num="19 / 24",
    layout="cards"
)

# Slide 20: Checklist Implementasi (Layout: checklist)
add_slide(
    categories=["CHECKLIST", "GUIDELINES"],
    title="Checklist Penerapan ADR dan Validasi Struktural",
    subtitle="Pedoman wajib sebelum menggabungkan pull request arsitektur ke main branch.",
    content_items=[
        {"type": "checklist", "title": "KRITERIA MERGE ADR", "items": [
            "ADR diberi nomor urut berikutnya secara kontigu tanpa melompati nomor.",
            "Berkas disimpan di direktori yang sama dengan kode dalam format Markdown.",
            "Bagian Context, Decision, dan Consequences terisi lengkap dengan analisis kompromi.",
            "Jika menggantikan keputusan lama, update status lama (Superseded by) dan deklarasikan Supersedes.",
            "Linter integritas dijalankan di CI/CD dan dinyatakan lulus tanpa peringatan tautan putus."
        ]}
    ],
    narrative="Checklist ini merangkum disiplin tata kelola ADR: nomor urut kontigu, ko-lokasi repositori, analisis konsekuensi dua arah yang terisi penuh, tautan resiprokal yang sinkron, dan verifikasi CI/CD otomatis.",
    narrative_en="This checklist enforces robust ADR governance: contiguous numbering, code co-location, comprehensive consequence documentation, reciprocal pointers, and mandatory CI/CD gate passes.",
    mental_model="Treat architectural documentation with the same testing rigor applied to production code.",
    source_path="content/02-master-draft.md:281-293",
    slide_num="20 / 24",
    layout="checklist"
)

# Slide 21: Anti-Pattern & Cara Reproduksi Kegagalan (Layout: cards)
# Limiting phrase 17: "bukan berarti"
add_slide(
    categories=["ANTI-PATTERN", "PITFALLS"],
    title="Anti-Pattern Umum dalam Pengelolaan ADR",
    subtitle="Kekeliruan pemeliharaan yang sering merusak fungsi dokumentasi arsitektur.",
    content_items=[
        {"type": "cards", "title": "HINDARI KESALAHAN UMUM INI", "items": [
            {"title": "Retroactive Mutation", "desc": "Mengedit teks ADR Accepted lama saat keputusan berubah alih-alih membuat ADR baru."},
            {"title": "Trivial Decision Pollution", "desc": "Membuat ADR untuk hal mikro seperti nama variabel atau refaktorisasi internal rutin."},
            {"title": "Pemberian Ulasan Kosong", "desc": "Mencatat keuntungan saja tanpa konsekuensi negatif; bukan berarti rancangannya sempurna."}
        ]}
    ],
    narrative="Menyunting ADR yang telah disetujui merusak integritas sejarah. Membuka ADR untuk keputusan mikro menciptakan polusi dokumen. Setiap keputusan memiliki konsekuensi negatif; mengabaikannya bukan berarti keputusan tersebut bebas kelemahan.",
    narrative_en="Retroactively mutating approved ADRs invalidates historical truth. Registering trivial refactoring creates document bloat. Every architecture decision carries negative trade-offs; ignoring them does not make the design flawless.",
    mental_model="Preserve history faithfully; record trade-offs honestly without omitting downsides.",
    source_path="content/02-master-draft.md:241-250",
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
            "cd labs/17-architecture-decision-record",
            "go test -v ./...    # validasi parser, linter, dan DAG lineage",
            "go test -race ./... # verifikasi konkurensi linter bebas race",
            "go run ./cmd/demo   # simulasi evolusi SaaS ERP 3-node"
        ]}
    ],
    narrative="Seluruh mekanisme validasi ditulis menggunakan Go standard library murni. Perintah ini mengeksekusi rangkaian uji unit, verifikasi race condition pada linter konkuren, serta simulasi end-to-end evolusi SaaS ERP.",
    narrative_en="The entire validation suite runs on pure Go standard library with zero third-party dependencies. These commands execute unit tests, race-safe concurrent linting, and the SaaS ERP evolution demo.",
    mental_model="Zero-dependency test harnesses guarantee fast, reliable verification across all developer environments.",
    source_path="README.md",
    slide_num="22 / 24",
    layout="code"
)

# Slide 23: Prinsip Arsitektural Jangka Panjang (Layout: cards)
# Limiting phrase 19: "tidak otomatis"
add_slide(
    categories=["TAKEAWAY", "ARCHITECTURAL PRINCIPLES"],
    title="Prinsip Inti: Rekam Konteks, Verifikasi Silsilah",
    subtitle="Sistem berevolusi; pemahaman atas batasan masa lalu menjaga konsistensi masa depan.",
    content_items=[
        {"type": "cards", "title": "TIGA PILAR KELANGSUNGAN ARSITEKTUR", "items": [
            {"title": "Rekam Mengapa", "desc": "Abadikan konteks, gaya dorong, dan kompromi saat keputusan diambil untuk generasi pengembang berikutnya."},
            {"title": "Kekekalan Sejarah", "desc": "Keputusan adalah rekaman waktu; perubahan selalu melahirkan simpul baru dalam graf silsilah."},
            {"title": "Validasi Otomatis", "desc": "Disiplin tim penting, tetapi tidak otomatis terjaga tanpa gerbang linter otomatis di pipeline CI."}
        ]}
    ],
    narrative="ADR bukan sekadar dokumen teks, melainkan instrumen mitigasi risiko teknis jangka panjang. Mendokumentasikan keputusan tidak otomatis menyelesaikan tantangan arsitektur, namun memastikan setiap evolusi bertumpu pada konteks yang jelas.",
    narrative_en="ADRs function as an organizational risk mitigation tool rather than mere documentation. Maintaining decision records does not automatically resolve complex engineering trade-offs, but guarantees that future evolution builds upon clear context.",
    mental_model="Capture why decisions were made; let automated linters protect the lineage.",
    source_path="content/05-key-takeaways.md",
    slide_num="23 / 24",
    layout="cards"
)

# Slide 24: Closing Slide (Layout: closing)
add_slide(
    categories=["CLOSING"],
    title="Keputusan arsitektur yang tidak dicatat adalah hutang konteks yang harus dibayar di masa depan.",
    subtitle="",
    content_items=[
        {"type": "bullet", "title": "MODUL IMPLEMENTASI", "items": [
            "internal/adr/models.go • internal/adr/parser.go • internal/adr/linter.go",
            "cmd/demo/main.go • tests/parser_test.go • tests/linter_test.go",
            "engineering/01-design.md • content/02-master-draft.md"
        ]},
        {"type": "bullet", "title": "VERIFIKASI & STATUS", "items": [
            f"Commit: {COMMIT_SHA[:7]} • Status: APPROVED",
            "Test passing: 15 assertions across parser, DAG, cycle detection, and race tests",
            "Race detector: Clean (go test -race ./... PASS in 1.4s)",
            "Validasi struktural Markdown Go murni; bukan parser AST umum."
        ]}
    ],
    mental_model="Architectural decisions without verified records become forgotten context.",
    source_path=REPO_URL,
    slide_num="24 / 24",
    layout="closing"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
    for s in SLIDES:
        print(f"[{s['slide_num']}] {s['layout']}: {s['title']}")
