#!/usr/bin/env python3
"""
Publication Definition for Software Engineering Lab 40: Property-Based Testing.
Defines 24 slides, metadata, colors, and full narratives for builder scripts.
"""
import os
import sys

OUTPUT_DIR = "labs/40-property-based-testing/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/40-property-based-testing"
COMMIT_SHA = "7d244968e91a2222291a27f346512976a26017e2"
LAB_NO = "40"
SLUG = "property-based-testing"
TITLE_ID = "Property-Based Testing"
SUBTITLE_ID = "Menemukan Edge-Case Tersembunyi via Universal Invariant & Shrinking"
TITLE_EN = "Property-Based Testing"
SUBTITLE_EN = "Finding Hidden Edge Cases via Universal Invariants and Shrinking"
AUTHOR = "Lukman / lukman-ss"
TOPIC_TAGS = "Property-Based Testing • Universal Invariants • Shrinking • Go"

DARK_NAVY = "0b1b2d"
TEAL = "087f76"
CREAM = "f4f1e8"
WHITE = "ffffff"
SLATE = "52606b"
ORANGE = "d97706"

SLIDES = []

def add_slide(categories, title, subtitle, content_items, mental_model, source_path="", slide_num="", narrative="", narrative_en="", title_en=""):
    """Add a slide to the presentation with support for in-depth editorial narratives."""
    SLIDES.append({
        "categories": categories,
        "title": title,
        "title_en": title_en or title,
        "subtitle": subtitle,
        "content_items": content_items,
        "mental_model": mental_model,
        "source": source_path,
        "slide_num": slide_num,
        "narrative": narrative,
        "narrative_en": narrative_en
    })

# Slide 1: Cover (Dark Navy retro terminal / art)
add_slide(
    categories=["COVER"],
    title=TITLE_ID,
    subtitle=SUBTITLE_ID,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "Retro CRT monitor beige dengan scanlines",
            "Terminal menampilkan invariant verification and counterexample shrinking",
            f"Penulis: {AUTHOR}",
        ]}
    ],
    mental_model="",
    source_path="",
    slide_num="01 / 24"
)

# Slide 2: Ilusi Keamanan Example-Based Testing
add_slide(
    categories=["MOTIVASI", "BATAS MODEL"],
    title="Ilusi Keamanan Example-Based Testing",
    title_en="The False Confidence of Example-Based Testing",
    subtitle="Contoh statis buatan manusia membatasi cakupan uji pada jalur yang diharapkan.",
    content_items=[
        {"title": "BIAS PENGEMBANG", "items": [
            "Unit test konvensional menguji nilai yang dipikirkan developer (happy path).",
            "Contoh hand-picked mengabaikan kombinasi liar, angka negatif, dan float precision.",
            "Test suite hijau 100% code coverage tetap menyisakan blindspot fatal di produksi."
        ], "type": "bullet"},
        {"title": "BUKTI LAB", "items": [
            "Format naif float lolos 5/5 contoh statis ($1.00, $5.25, $99.99, $10.50, $0.01).",
            "Saat diuji invariant properti dengan input acak, 1.000/1.000 input gagal roundtrip."
        ], "type": "bullet"},
        {"title": "PARADIGMA PBT", "items": [
            "Uji kebenaran universal (invariants) terhadap ratusan input acak tanpa asumsi sempit."
        ], "type": "bullet"}
    ],
    narrative="Unit test berbasis contoh (example-based testing) bekerja dengan menegaskan bahwa f(x) == y untuk segelintir sampel statis. Masalah mendasarnya adalah bias konfirmasi: engineer cenderung hanya menguji kombinasi nilai yang telah mereka antisipasi dalam implementasi. Pada dataset kecil atau input tipikal, kode tampak sempurna dan seluruh baris assertions berstatus PASS. Namun ketika berhadapan dengan data acak tak terduga di lingkungan produksi, asumsi tersembunyi runtuh seketika dan memicu kegagalan sistemik yang fatal.",
    narrative_en="Conventional example-based unit tests check whether f(x) == y for a few hand-picked samples. The core vulnerability is developer confirmation bias: engineers only test scenarios they already anticipated. While test suites show 100% green coverage, edge-case assumptions break immediately when exposed to real-world, uncurated input distributions.",
    mental_model="Example-based tests check developer assumptions; property-based tests explore domain boundaries.",
    source_path="labs/40-property-based-testing/README.md",
    slide_num="02 / 24"
)

# Slide 3: Tiga Pilar Siklus Property-Based Testing
add_slide(
    categories=["ARSITEKTUR", "LIFECYCLE"],
    title="Tiga Pilar Siklus Property-Based Testing",
    title_en="Three Pillars of Property-Based Testing",
    subtitle="Generasi domain, evaluasi invariant universal, dan reduksi counterexample.",
    content_items=[
        {"title": "PILAR SIKLUS PBT", "items": [
            "Tahap | Tanggung Jawab | Output Sistem",
            "1. Generator | Buat ratusan sampel input acak terdistribusi | Slices, numbers, structs",
            "2. Invariant | Evaluasi properti universal f(x) | Validasi true / false",
            "3. Shrinker | Pangkas kegagalan ke bentuk terkecil | Minimal counterexample (e.g. [-1])"
        ], "type": "table"}
    ],
    narrative="Siklus kerja Property-Based Testing membagi tugas antara komputasi mesin dan perancangan manusia secara elegan. Manusia bertanggung jawab mendefinisikan batasan domain serta relasi invariant matematika yang wajib selalu bernilai benar. Mesin PBT kemudian bertindak sebagai musuh (adversary) yang mengeksekusi ratusan kombinasi data acak melalui generator, memverifikasi invariant pada setiap iterasi, dan segera mengisolasi kegagalan ke bentuk data paling sederhana menggunakan algoritma shrinking deterministik.",
    narrative_en="Property-based testing splits responsibilities cleanly between human and machine. Engineers define domain constraints and mathematical invariants that must unconditionally hold. The PBT engine acts as an automated adversary, generating random candidates, checking invariants, and shrinking failures to their simplest form.",
    mental_model="Generate broadly, assert universally, and shrink failures ruthlessly to their core.",
    source_path="labs/40-property-based-testing/content/02-master-draft.md",
    slide_num="03 / 24"
)

# Slide 4: Empat Pola Invariant Kanonikal
add_slide(
    categories=["TAKSONOMI", "INVARIANTS"],
    title="Empat Pola Invariant Kanonikal",
    title_en="Four Canonical Invariant Patterns",
    subtitle="Pola relasional universal untuk membuktikan kebenaran sistem tanpa oracle eksternal.",
    content_items=[
        {"title": "POLA INVARIANT", "items": [
            "Pola Invariant | Formula Formal | Contoh Kasus Nyata",
            "Roundtrip | Decode(Encode(x)) == x | Serialisasi JSON, mata uang, kompresi",
            "Idempotence | f(f(x)) == f(x) | Normalisasi teks, merge interval, cleanup",
            "Oracle Equiv | f_robust(x) == f_naive(x) | Verifikasi algoritma teroptimasi vs referensi",
            "Easy-to-Verify | HardToSolve(x) -> EasyVerify(res) | Validasi shortest path, partisi graf"
        ], "type": "table"}
    ],
    narrative="Kunci keberhasilan PBT terletak pada pemilihan pola invariant yang tepat tanpa menulis ulang logika implementasi di dalam assertion test. Pola Roundtrip menjamin integritas transformasi bolak-balik tanpa data corruption. Pola Idempotence memastikan stabilitas state di mana pemanggilan berulang tidak mengubah hasil. Pola Test Oracle membandingkan algoritma berperforma tinggi terhadap implementasi referensi sederhana. Pola Easy-to-Verify memeriksa keabsahan properti solusi akhir dengan biaya komputasi murah.",
    narrative_en="The core of PBT lies in choosing relational invariants rather than duplicating implementation logic in test assertions. Roundtrips ensure lossless data transformations, idempotence guarantees state stability, test oracles verify complex optimizations against naive references, and easy-to-verify checks validate structural correctness at low computational cost.",
    mental_model="Invariant patterns are executable mathematical specifications for software behavior.",
    source_path="labs/40-property-based-testing/content/02-master-draft.md",
    slide_num="04 / 24"
)

# Slide 5: Kegagalan Presisi Float pada Mata Uang
add_slide(
    categories=["ANTI-PATTERN", "CURRENCY"],
    title="Kegagalan Presisi Float pada Mata Uang",
    title_en="Floating-Point Precision Loss in Currency",
    subtitle="Representasi pecahan biner IEEE 754 merusak presisi desimal secara laten.",
    content_items=[
        {"title": "NAIVE FLOAT IMPLEMENTATION", "items": [
            "type NaiveCurrency struct{}",
            "",
            "func (n NaiveCurrency) Format(dollars float64) string {",
            "    return fmt.Sprintf(\"$%.2f\", dollars)",
            "}",
            "",
            "func (n NaiveCurrency) Parse(s string) (float64, error) {",
            "    s = strings.TrimPrefix(s, \"$\")",
            "    return strconv.ParseFloat(s, 64)",
            "}"
        ], "type": "code"}
    ],
    narrative="Menyimpan nilai moneter dalam tipe floating-point float64 adalah anti-pattern klasik dalam rekayasa perangkat lunak. Pemformatan dengan sintaks pembulatan desimal string memotong digit pecahan sub-cent tanpa kendali. Lebih buruk lagi, konversi string kembali ke float64 tunduk pada aproksimasi biner IEEE 754. Pada pengujian konvensional dengan nilai bulat rapi, bug ini tidak pernah muncul, menciptakan ilusi bahwa modul akuntansi bekerja normal padahal pembulatan sen terus bocor.",
    narrative_en="Storing currency as float64 is a notorious software engineering anti-pattern. Format strings truncate sub-cent precision, while string-to-float conversions suffer from IEEE 754 binary approximation limits. Hand-picked unit test values hide this defect entirely, allowing financial rounding leaks into production.",
    mental_model="Floating point binary approximations inevitably corrupt fractional decimal values.",
    source_path="internal/currency/currency.go:9-21",
    slide_num="05 / 24"
)

# Slide 6: PBT Mendeteksi Kebocoran Float 100%
add_slide(
    categories=["BUKTI KODE", "QUICK CHECK"],
    title="PBT Mendeteksi Kebocoran Float 100%",
    title_en="PBT Catches Float Leaks with 100% Failure Rate",
    subtitle="Pemeriksaan Roundtrip Parse(Format(x)) == x menggagalkan seluruh sampel acak.",
    content_items=[
        {"title": "PROPERTY TEST CHECK", "items": [
            "func TestPropertyNaiveCurrencyFails(t *testing.T) {",
            "    naive := NaiveCurrency{}",
            "    prop := func(val NaiveFloat) bool {",
            "        dollars := float64(val)",
            "        formatted := naive.Format(dollars)",
            "        parsed, err := naive.Parse(formatted)",
            "        if err != nil { return false }",
            "        return parsed == dollars // Gagal pada sub-cent!",
            "    }",
            "    err := quick.Check(prop, &quick.Config{MaxCount: 100})",
            "    // quick.Check mengembalikan error counterexample",
            "}"
        ], "type": "code"}
    ],
    narrative="Ketika properti roundtrip diuji menggunakan paket standar Go testing/quick, kelemahan implementasi float64 langsung terekspos tanpa ampun. Dari 1.000 iterasi angka acak, seluruh 1.000 input mengalami kegagalan roundtrip karena nilai desimal asli terpotong oleh format dua angka di belakang koma. PBT membuktikan secara empiris bahwa fungsi Parse(Format(x)) bukanlah fungsi identitas bagi domain float64, mematahkan asumsi keamanan yang sebelumnya diberikan oleh unit test statis.",
    narrative_en="When the roundtrip property is tested against Go's testing/quick package, the float64 implementation fails unconditionally. Across 1,000 randomized iterations, all 1,000 inputs violate the identity invariant Parse(Format(x)) == x due to sub-cent truncation, proving that example tests produced false confidence.",
    mental_model="A roundtrip property check exposes latent serialization and truncation flaws immediately.",
    source_path="internal/currency/currency_test.go:37-61",
    slide_num="06 / 24"
)

# Slide 7: Desain Robust: Integer Cents Representation
add_slide(
    categories=["SOLUSI", "REPRESENTASI"],
    title="Desain Robust: Integer Cents Representation",
    title_en="Robust Representation: Exact Integer Cents",
    subtitle="Menghilangkan ambiguitas desimal dengan menyimpan nilai moneter sebagai int64.",
    content_items=[
        {"title": "INTEGER CENTS STRUCT", "items": [
            "type RobustAmount struct {",
            "    Cents int64",
            "}",
            "",
            "func (r RobustAmount) Format() string {",
            "    absCents := r.Cents; sign := \"\"",
            "    if absCents < 0 { sign = \"-\"; absCents = -absCents }",
            "    dollars := absCents / 100; cents := absCents % 100",
            "    return fmt.Sprintf(\"%s$%d.%02d\", sign, dollars, cents)",
            "}"
        ], "type": "code"}
    ],
    narrative="Solusi rekayasa yang terbukti secara industri untuk domain mata uang adalah merepresentasikan nominal terkecil sebagai bilangan bulat utuh (integer cents bertipe int64). Format visual desimal hanyalah proyeksi tampilan pada lapisan presentasi string, sedangkan operasi aritmatika internal dan penyimpanan data tetap berada di ranah bilangan diskrit. Pendekatan ini secara definitif menutup peluang terjadinya distorsi presisi IEEE 754 dan menjamin determinisme mutlak.",
    narrative_en="The industry-proven solution for monetary values is storing the smallest fractional unit as an integer (int64 cents). Decimal string formatting becomes a presentation-layer projection, while arithmetic remains purely discrete, eliminating floating-point imprecision permanently.",
    mental_model="Represent fractional quantities in discrete atomic units to eliminate rounding artifacts.",
    source_path="internal/currency/currency.go:23-42",
    slide_num="07 / 24"
)

# Slide 8: Custom Generator & Boundary Value Bias
add_slide(
    categories=["PBT PATTERN", "GENERATOR"],
    title="Custom Generator & Boundary Value Bias",
    title_en="Custom Generator with Boundary Value Bias",
    subtitle="Mengarahkan distribusi acak ke titik-titik rawan: batas nol, negatif, dan ekstrem.",
    content_items=[
        {"title": "BIASED GENERATOR LOGIC", "items": [
            "func (RobustAmount) Generate(r *rand.Rand, size int) reflect.Value {",
            "    switch r.Intn(5) {",
            "    case 0: cents = 0                         // Boundary zero",
            "    case 1: cents = int64(r.Intn(200) - 100)  // Near-zero negatives",
            "    case 2: cents = int64(r.Intn(1000000))    // Normal positive",
            "    case 3: cents = -int64(r.Intn(1000000))   // Normal negative",
            "    default: cents = r.Int63n(1e12) - 5e11    // Multi-billion bounds",
            "    }",
            "    return reflect.ValueOf(RobustAmount{Cents: cents})",
            "}"
        ], "type": "code"}
    ],
    narrative="Generator acak yang murni seragam (uniform distribution) sering kali melewatkan titik-titik krusial tempat bersarangnya bug batas logika. Implementasi generator RobustAmount pada Lab 40 sengaja menyuntikkan bias probabilitas ke arah nilai rawan: batas nol, transisi bilangan negatif kecil di sekitar nol (-100 hingga +99 sen), serta nilai nominal raksasa multi-miliar sen. Generator yang cerdas secara agresif membombardir tepi spektrum domain untuk menguji ketahanan boundary handling.",
    narrative_en="Uniform random generators often miss critical edge cases where boundary bugs reside. The RobustAmount custom generator injects deliberate probability biases toward zero, near-zero negative transitions, and multi-billion cent bounds, actively probing boundary handling resilience.",
    mental_model="Biased generators force tests toward dangerous boundary conditions instead of safe averages.",
    source_path="internal/currency/currency_test.go:63-81",
    slide_num="08 / 24"
)

# Slide 9: 1.000 Iterasi Lolos Roundtrip Integer Cents
add_slide(
    categories=["VERIFIKASI", "ROUNDTRIP"],
    title="1.000 Iterasi Lolos Roundtrip Integer Cents",
    title_en="1,000 Iterations Pass Integer Roundtrip",
    subtitle="Membuktikan secara deterministik bahwa ParseRobust(Format(x)) == x.",
    content_items=[
        {"title": "HASIL VERIFIKASI ROUNDTRIP", "items": [
            "Metrik Uji | Naive Float64 | Robust Integer Cents",
            "Total Iterasi Acak | 1.000 iterasi | 1.000 iterasi",
            "Tingkat Kelulusan | 0% (1.000 gagal) | 100% (1.000 lolos)",
            "Cakupan Domain | Terbatas hand-picked | Nol, negatif, triliun sen",
            "Invarian Terbukti | Gagal total | Parse(Format(x)) == x VALID"
        ], "type": "table"}
    ],
    narrative="Dengan beralih ke integer cents dan parser yang presisi, pengujian properti TestPropertyRobustAmountRoundtrip berhasil mengeksekusi 1.000 iterasi berturut-turut tanpa satu pun pelanggaran invariant. Terlepas dari apakah input berupa angka negatif dengan tanda minus, nilai nol bersih, atau nominal transaksi skala enterprise, relasi Parse(Format(x)) == x selalu terpenuhi. Ini membuktikan bahwa arsitektur data baru telah memenuhi kontrak integritas serialisasi secara matematis.",
    narrative_en="By adopting integer cents and exact string parsing, TestPropertyRobustAmountRoundtrip executes 1,000 randomized iterations without a single invariant violation. Across negative amounts, zero, and enterprise-scale values, Parse(Format(x)) == x holds true unconditionally.",
    mental_model="A verified roundtrip invariant proves mathematical isomorphism between data and representation.",
    source_path="internal/currency/currency_test.go:83-102",
    slide_num="09 / 24"
)

# Slide 10: Kasus Interval Merge & Asumsi Tersembunyi
add_slide(
    categories=["STUDI KASUS", "INTERVAL"],
    title="Kasus Interval Merge & Asumsi Tersembunyi",
    title_en="Interval Merging & Hidden Preconditions",
    subtitle="Asumsi implisit bahwa urutan input selalu terurut memicu kecacatan algoritma.",
    content_items=[
        {"title": "NAIVE MERGE CODE", "items": [
            "func NaiveMerge(intervals []Interval) []Interval {",
            "    if len(intervals) == 0 { return nil }",
            "    out := []Interval{intervals[0]}",
            "    for _, iv := range intervals[1:] {",
            "        last := &out[len(out)-1]",
            "        if iv.Start <= last.End {",
            "            if iv.End > last.End { last.End = iv.End }",
            "        } else { out = append(out, iv) }",
            "    }",
            "    return out // BUG: Mengasumsikan input terurut!",
            "}"
        ], "type": "code"}
    ],
    narrative="Kelemahan algoritma sering kali berakar pada prakondisi implisit yang tidak dinyatakan secara tegas dalam kontrak kode. Fungsi NaiveMerge dirancang untuk menggabungkan interval waktu yang saling tumpang tindih. Namun, loop perulangan langsung membandingkan elemen berjalan terhadap elemen terakhir dari slice output, mengasumsikan input sudah terurut berdasarkan titik awal (Start). Ketika dipanggil dengan contoh terurut, fungsi bekerja mulus; namun begitu urutan berubah, output menjadi rusak.",
    narrative_en="Algorithmic bugs frequently stem from unstated implicit preconditions. NaiveMerge merges overlapping intervals by checking incoming elements against the last emitted interval, implicitly assuming presorted inputs. Hand-picked sorted examples succeed, but arbitrary ordering produces silent corruption.",
    mental_model="Unenforced implicit preconditions are bugs waiting for real-world input to trigger them.",
    source_path="internal/interval/interval.go:10-28",
    slide_num="10 / 24"
)

# Slide 11: Test Oracle: Membongkar 85% Discrepancy
add_slide(
    categories=["ORACLE", "EQUIVALENCE"],
    title="Test Oracle: Membongkar 85% Discrepancy",
    title_en="Test Oracle Uncovers 85% Discrepancy",
    subtitle="Membandingkan implementasi naif terhadap implementasi robust secara komparatif.",
    content_items=[
        {"title": "ORACLE INVARIANT CHECK", "items": [
            "func TestPropertyNaiveMergeFails(t *testing.T) {",
            "    prop := func(ivs IntervalSlice) bool {",
            "        naive := NaiveMerge([]Interval(ivs))",
            "        robust := RobustMerge([]Interval(ivs))",
            "        return intervalsEqual(naive, robust)",
            "    }",
            "    err := quick.Check(prop, &quick.Config{MaxCount: 100})",
            "    // 85 dari 100 input acak menghasilkan ketidaksesuaian!",
            "}"
        ], "type": "code"}
    ],
    narrative="Teknik Test Oracle memanfaatkan dua implementasi berbeda untuk input yang sama: f_naive(x) versus f_robust(x). RobustMerge secara eksplisit melakukan pengurutan sort.Slice sebelum proses penggabungan dilakukan. Ketika dievaluasi dengan generator irisan acak, pengujian membongkar bahwa 85 dari 100 input menghasilkan inkonsistensi keluaran antara NaiveMerge dan RobustMerge. Generator acak secara instan mengekspos betapa rapuhnya kode naif di hadapan array yang tidak terurut.",
    narrative_en="The Test Oracle pattern compares two distinct implementations against identical inputs: f_naive(x) versus f_robust(x). RobustMerge enforces pre-sorting via sort.Slice. Under property evaluation, 85 out of 100 randomized slices expose discrepancies, demonstrating the fragility of unsorted inputs.",
    mental_model="A reference implementation acts as an oracle to evaluate alternative algorithms at scale.",
    source_path="internal/interval/interval_test.go:43-59",
    slide_num="11 / 24"
)

# Slide 12: Invariant Idempotence pada Penggabungan Data
add_slide(
    categories=["IDEMPOTENCE", "STABILITAS"],
    title="Invariant Idempotence pada Penggabungan Data",
    title_en="Idempotence Invariant on Data Merging",
    subtitle="Menjalankan operasi berulang pada hasil penggabungan tidak boleh mengubah state.",
    content_items=[
        {"title": "IDEMPOTENCE EQUATION", "items": [
            "Formula Teoretis: f(f(x)) == f(x)",
            "",
            "func TestPropertyRobustMergeIdempotence(t *testing.T) {",
            "    prop := func(ivs IntervalSlice) bool {",
            "        once := RobustMerge([]Interval(ivs))",
            "        twice := RobustMerge(once)",
            "        return intervalsEqual(once, twice)",
            "    }",
            "    quick.Check(prop, &quick.Config{MaxCount: 1000})",
            "}"
        ], "type": "code"}
    ],
    narrative="Idempotensi adalah sifat matematika fundamental dalam arsitektur komputasi: f(f(x)) == f(x). Dalam konteks pemrosesan data, menjalankan fungsi penggabungan pada kumpulan interval yang sudah digabungkan wajib menghasilkan output yang sama persis tanpa mutasi tambahan. Pengujian TestPropertyRobustMergeIdempotence memvalidasi 1.000 iterasi irisan acak, membuktikan bahwa operasi RobustMerge bersifat stabil dan aman dari efek samping mutasi berulang.",
    narrative_en="Idempotence is a foundational invariant: f(f(x)) == f(x). In data processing, re-applying a merge operation to an already merged slice must yield an identical result without side effects. Running 1,000 iterations confirms that RobustMerge is mathematically stable across all permutations.",
    mental_model="Idempotence guarantees operational convergence: repetition causes zero state drift.",
    source_path="internal/interval/interval_test.go:61-77",
    slide_num="12 / 24"
)

# Slide 13: Invariant Struktural: Non-Overlapping Slices
add_slide(
    categories=["STRUKTUR", "INVARIANTS"],
    title="Invariant Struktural: Non-Overlapping Slices",
    title_en="Structural Invariant: Non-Overlapping Slices",
    subtitle="Setiap pasangan interval berurutan dijamin tidak saling bertumpukan.",
    content_items=[
        {"title": "STRUKTURAL ASSERTION", "items": [
            "func TestPropertyRobustMergeNonOverlapping(t *testing.T) {",
            "    prop := func(ivs IntervalSlice) bool {",
            "        merged := RobustMerge([]Interval(ivs))",
            "        for i := 0; i < len(merged)-1; i++ {",
            "            if merged[i].End >= merged[i+1].Start {",
            "                return false // Ada tumpang tindih!",
            "            }",
            "        }",
            "        return true",
            "    }",
            "    quick.Check(prop, &quick.Config{MaxCount: 1000})",
            "}"
        ], "type": "code"}
    ],
    narrative="Selain relasi f(f(x)), kebenaran algoritma penggabungan interval harus memenuhi invariant struktural domain: untuk setiap indeks i dan i+1 pada irisan keluaran, titik akhir interval pertama harus lebih kecil dari titik awal interval berikutnya (out[i].End < out[i+1].Start). Pengujian 1.000 iterasi acak memastikan bahwa tidak ada irisan ganda atau fragmen bertumpuk yang lolos, menjamin partisi interval bebas anomali.",
    narrative_en="Beyond relational invariants, merge algorithms must fulfill structural domain properties: for all consecutive intervals in the output, out[i].End < out[i+1].Start. Across 1,000 iterations, the test proves that no overlapping fragments escape the merging logic.",
    mental_model="Structural invariants verify domain-specific physical constraints on output shapes.",
    source_path="internal/interval/interval_test.go:79-99",
    slide_num="13 / 24"
)

# Slide 14: Anatomi Masalah: Kompleksitas Counterexample
add_slide(
    categories=["SHRINKING", "DIAGNOSTIK"],
    title="Anatomi Masalah: Kompleksitas Counterexample",
    title_en="Problem Anatomy: Counterexample Noise",
    subtitle="Kegagalan awal dari generator acak mengandung derau data masif yang mengaburkan akar masalah.",
    content_items=[
        {"title": "SAMPLE INITIAL FAILING INPUT", "items": [
            "Array Awal (10 elemen acak bernilai besar & campur aduk):",
            "[137, 18, 100, -27, 95, 107, 126, 78, -7, 79]",
            "",
            "Pertanyaan Developer saat Melihat Log:",
            "• Apakah bug dipicu oleh angka 137? Atau ukuran array bernilai 10?",
            "• Apakah urutan elemen saling mempengaruhi?",
            "• Elemen mana yang menjadi pemicu minimal pelanggaran invariant?"
        ], "type": "bullet"}
    ],
    narrative="Ketika mesin PBT mendeteksi pelanggaran invariant pada input acak berukuran besar, output mentahnya sering kali berupa struktur data raksasa yang membingungkan engineer. Pada contoh demo, input yang memicu kegagalan predicate adalah array 10 elemen acak yang bercampur antara ratusan angka positif dan dua bilangan negatif (-27 dan -7). Tanpa mekanisme reduksi, developer akan menghabiskan waktu berjam-jam memilah data mana yang sesungguhnya menjadi biang kerok kegagalan sistem.",
    narrative_en="When a PBT engine catches an invariant violation on randomized inputs, the raw failure payload is often large and noisy. In our demo, an arbitrary 10-element slice contains mixed positive values and two negatives. Without automated reduction, developers waste hours isolating the exact trigger.",
    mental_model="Raw random failures contain heavy noise; automated shrinking isolates the minimal signal.",
    source_path="labs/40-property-based-testing/content/02-master-draft.md",
    slide_num="14 / 24"
)

# Slide 15: Konsep dan Algoritma Counterexample Shrinking
add_slide(
    categories=["ALGORITMA", "SHRINKER"],
    title="Konsep dan Algoritma Shrinking",
    title_en="Counterexample Shrinking Algorithm",
    subtitle="Reduksi deterministik dari kegagalan kompleks menuju kasus reproduksi terkecil.",
    content_items=[
        {"title": "TIGA STRATEGI REDUKSI", "items": [
            "Fase Algoritma | Aksi Reduksi | Target Optimasi",
            "1. Binary Sectioning | Potong separuh kiri / kanan array | Pangkas ukuran O(log N)",
            "2. Element Removal | Hapus elemen individual satu per satu | Eliminasi noise non-kritis",
            "3. Value Reduction | Bagi dua nilai elemen menuju 0 | Minimalisasi magnitudo nilai"
        ], "type": "table"}
    ],
    narrative="Algoritma counterexample shrinking dirancang untuk mentransformasikan input kegagalan mentah menjadi *minimal reproducible counterexample*. Algoritma bekerja secara deterministik dalam tiga fase bertingkat: pertama, melakukan pemotongan biner (halving) terhadap panjang array; kedua, mencoba menghapus elemen satu per satu; dan ketiga, memperkecil magnitudo nilai individual mendekati nol. Setiap langkah yang tetap mempertahankan status gagal (FAIL) diadopsi sebagai basis input baru.",
    narrative_en="Counterexample shrinking deterministically transforms complex failing inputs into minimal reproducible test cases. It operates in three sequential phases: binary list halving, individual element removal, and numerical reduction toward zero, retaining failing states at each stage.",
    mental_model="Shrinking is binary search applied to bug reproducers: drop halves, drop items, drop magnitude.",
    source_path="internal/shrinker/shrinker.go:24-60",
    slide_num="15 / 24"
)

# Slide 16: Binary Sectioning: Memangkas Panjang Array
add_slide(
    categories=["SHRINKER", "CODE"],
    title="Binary Sectioning: Memangkas Panjang Array",
    title_en="Binary Sectioning: Halving Array Length",
    subtitle="Menguji apakah separuh bagian data masih mempertahankan kegagalan invarian.",
    content_items=[
        {"title": "CHUNK REMOVAL STRATEGY", "items": [
            "if len(current) > 1 {",
            "    mid := len(current) / 2",
            "    left := current[:mid]",
            "    if !predicate(left) {",
            "        current = append([]int(nil), left...)",
            "        changed = true; continue",
            "    }",
            "    right := current[mid:]",
            "    if !predicate(right) {",
            "        current = append([]int(nil), right...)",
            "        changed = true; continue",
            "    }",
            "}"
        ], "type": "code"}
    ],
    narrative="Fase pertama dari shrinking adalah binary chunking. Algoritma membagi slice menjadi dua bagian (kiri dan kanan), lalu mengevaluasi predicate pada masing-masing belahan. Jika belahan kiri masih memicu kegagalan, separuh bagian kanan langsung dibuang secara permanen. Dalam demo, array 10 elemen terpotong menjadi 5 elemen [137, 18, 100, -27, 95] pada langkah 1, lalu terpotong lagi menjadi 3 elemen [100, -27, 95] pada langkah 3, mempercepat pencarian secara eksponensial.",
    narrative_en="The first phase of shrinking is binary list sectioning. The slice is halved and each section evaluated against the predicate. If the left half still fails, the entire right half is discarded immediately, reducing input size from 10 elements to 5, then to 3 within a few steps.",
    mental_model="Halving input payloads quickly eliminates large clusters of irrelevant test data.",
    source_path="internal/shrinker/shrinker.go:61-95",
    slide_num="16 / 24"
)

# Slide 17: Reduksi Elemen & Magnitudo Menuju Nol
add_slide(
    categories=["SHRINKER", "VALUE REDUCTION"],
    title="Reduksi Elemen & Magnitudo Menuju Nol",
    title_en="Element Removal and Value Reduction",
    subtitle="Menghapus elemen per satuan dan memperkecil nilai numerik hingga batas minimal.",
    content_items=[
        {"title": "VALUE REDUCTION LOOP", "items": [
            "// Coba set elemen menjadi 0 atau bagi dua nilainya",
            "candidates := []int{0, current[i] / 2}",
            "for _, cand := range candidates {",
            "    testSlice := append([]int(nil), current...)",
            "    testSlice[i] = cand",
            "    if !predicate(testSlice) {",
            "        current = testSlice",
            "        changed = true; break",
            "    }",
            "}"
        ], "type": "code"}
    ],
    narrative="Setelah panjang array dipangkas hingga tersisa 1 elemen tunggal berisi [-27], algoritma melanjutkan ke strategi ketiga: reduksi nilai numerik. Nilai -27 secara bertahap dibagi dua mendekati nol: -27 -> -13 -> -6 -> -3 -> -1. Setiap langkah diverifikasi apakah invarian masih dilanggar. Ketika mencapai -1, langkah berikutnya yaitu 0 menghasilkan PASS (karena 0 bukan negatif). Algoritma berhenti dan menetapkan [-1] sebagai nilai batas minimal mutlak.",
    narrative_en="Once array length is reduced to a single-element slice [-27], the shrinker applies numerical reduction toward zero: -27 -> -13 -> -6 -> -3 -> -1. When candidate 0 passes the non-negative check, the shrinker terminates, locking in [-1] as the exact minimal boundary counterexample.",
    mental_model="Gradual numerical reduction isolates the exact mathematical threshold of failure.",
    source_path="internal/shrinker/shrinker.go:119-151",
    slide_num="17 / 24"
)

# Slide 18: Jejak 21 Langkah Menuju Counterexample Minimal
add_slide(
    categories=["DEMO", "TRACE"],
    title="Jejak 21 Langkah Menuju Minimal",
    title_en="21-Step Shrinking Trace to Minimal Bug",
    subtitle="Bukti log eksekusi nyata dari FindAndShrink pada internal/shrinker.",
    content_items=[
        {"title": "JEJAK LANGKAH SHRINKING NYATA", "items": [
            "Input Awal (10 elemen): [137 18 100 -27 95 107 126 78 -7 79]",
            "Step  1: remove-right-half         [5 elem] [137 18 100 -27 95] -> FAIL",
            "Step  3: remove-left-half          [3 elem] [100 -27 95]        -> FAIL",
            "Step  5: remove-left-half          [2 elem] [-27 95]            -> FAIL",
            "Step  6: remove-right-half         [1 elem] [-27]               -> FAIL",
            "Step  9: reduce-val to -13         [1 elem] [-13]               -> FAIL",
            "Step 12: reduce-val to -6          [1 elem] [-6]                -> FAIL",
            "Step 15: reduce-val to -3          [1 elem] [-3]                -> FAIL",
            "Step 18: reduce-val to -1          [1 elem] [-1]                -> FAIL",
            "Hasil Akhir (Minimal Counterexample): [-1] (Tepat 21 langkah evaluasi!)"
        ], "type": "code"}
    ],
    narrative="Log terminal dari eksekusi demo memperlihatkan transparansi penuh dari proses reduksi data. Dari input awal acak sebanyak 10 elemen yang berantakan, mesin hanya membutuhkan 21 evaluasi predikat untuk mengeliminasi 9 elemen yang tidak relevan serta memperkecil nilai -27 menjadi -1. Hasil akhir [-1] memberikan kejelasan diagnostik instan bagi developer: masalah bukan pada panjang array atau interaksi elemen, melainkan keberadaan bilangan negatif sekecil apa pun.",
    narrative_en="The demo terminal trace illustrates complete diagnostic clarity. From a chaotic 10-element array, the engine executes 21 predicate evaluations to eliminate 9 irrelevant elements and reduce -27 to -1. The resulting [-1] reveals the root cause instantly: any negative value violates the invariant.",
    mental_model="Shrinking eliminates debugging guesswork by delivering the minimal reproducible unit.",
    source_path="cmd/demo/main.go:160-210",
    slide_num="18 / 24"
)

# Slide 19: Integrasi Go Standard Library testing/quick
add_slide(
    categories=["GO STDLIB", "TOOLING"],
    title="Integrasi Go Standard Library testing/quick",
    title_en="Go Standard Library: testing/quick",
    subtitle="Membangun PBT murni tanpa dependensi eksternal pihak ketiga.",
    content_items=[
        {"title": "QUICK CHECK CONFIGURATION", "items": [
            "cfg := &quick.Config{",
            "    MaxCount: 1000, // Jumlah iterasi acak yang dieksekusi",
            "    Rand:     rand.New(rand.NewSource(seed)), // Determinisme",
            "}",
            "",
            "err := quick.Check(func(a RobustAmount) bool {",
            "    return ParseRobust(a.Format()) == a",
            "}, cfg)",
            "if err != nil {",
            "    t.Fatalf(\"property violation: %v\", err)",
            "}"
        ], "type": "code"}
    ],
    narrative="Salah satu keunggulan ekosistem Go adalah ketersediaan paket bawaan testing/quick di dalam standard library tanpa memerlukan dependensi pihak ketiga. Paket ini menyediakan fungsi quick.Check yang mengeksekusi properti hingga ribuan kali dengan konfigurasi MaxCount dan seed kustom. Generator bawaan secara otomatis merefleksikan tipe data dasar, sementara interface quick.Generator memungkinkan engineer menyuntikkan generator buatan sendiri untuk tipe data kompleks.",
    narrative_en="Go provides built-in property testing via testing/quick in the standard library, requiring zero third-party dependencies. quick.Check evaluates properties across thousands of iterations with configurable counts and seeds. Custom types implement quick.Generator to supply domain-specific generation logic.",
    mental_model="Zero-dependency property testing is already built into the Go standard library.",
    source_path="internal/currency/currency_test.go:85-102",
    slide_num="19 / 24"
)

# Slide 20: Determinisme, Seeding, & CI/CD Pipeline
add_slide(
    categories=["CI/CD", "DETERMINISME"],
    title="Determinisme, Seeding, & CI/CD Pipeline",
    title_en="Determinism, Seeding, and CI/CD Pipelines",
    subtitle="Mengelola trade-off antara keacakan eksploratif dan reproduktibilitas pengujian.",
    content_items=[
        {"title": "STRATEGI SEED & ITERASI", "items": [
            "Lingkungan | Konfigurasi Seed | Iterasi | Tujuan",
            "Local PR Check | Fixed Seed (deterministik) | 100 - 500 | Cepat & lolos tanpa flakiness",
            "Nightly CI | Dynamic Time Seed | 10.000 - 50.000 | Eksplorasi bug liar & boundary",
            "Bug Reproduction | Injeksi Failing Seed | 1 run minimal | Reproduksi 100% konsisten"
        ], "type": "table"}
    ],
    narrative="Kritik yang kerap diarahkan pada PBT adalah risiko munculnya flaky tests jika seed acak berubah di setiap run CI. Solusi profesional adalah menerapkan dual-strategy: pada local testing dan PR verification, gunakan seed tetap dengan 100–500 iterasi agar eksekusi berlangsung sub-detik dan stabil. Pada nightly build atau background fuzzer, aktifkan seed dinamis berbasis timestamp dengan 50.000 iterasi. Setiap kali kegagalan ditemukan, log seed wajib dicatat agar developer dapat mereproduksinya secara instan.",
    narrative_en="A common critique of PBT is test flakiness from changing random seeds. Professional teams adopt a dual strategy: fixed seeds with 100-500 iterations for sub-second PR checks, and dynamic timestamp seeds with 50,000 iterations in nightly CI. Logging the failure seed ensures 100% reproducible debugging locally.",
    mental_model="Exploratory randomness belongs in nightly CI; deterministic seeds belong in pull requests.",
    source_path="labs/40-property-based-testing/content/02-master-draft.md",
    slide_num="20 / 24"
)

# Slide 21: Matriks Perbandingan Example vs Property-Based
add_slide(
    categories=["KOMPARASI", "MATRIKS"],
    title="Matriks Komparasi Example vs PBT",
    title_en="Comparison Matrix: Example vs Property-Based",
    subtitle="Perbedaan mendasar dalam cakupan, bias pengujian, dan kemampuan isolasi bug.",
    content_items=[
        {"title": "TABEL PERBANDINGAN KOMPREHENSIF", "items": [
            "Dimensi Evaluasi | Example-Based Testing | Property-Based Testing",
            "Cakupan Input | Titik diskrit statis (hand-picked) | Spektrum kontinu acak terdistribusi",
            "Bias Pengujian | Tinggi (mengikuti intuisi engineer) | Rendah (adversarial fuzzing)",
            "Bug Boundary | Kerap terlewat (sub-cent, negative) | Terdeteksi via biased generator",
            "Isolasi Bug | Manual via trace & debugger | Otomatis via counterexample shrinking",
            "Tingkat Keyakinan | False confidence (jalur aman) | Calibrated confidence (invarian kokoh)"
        ], "type": "table"}
    ],
    narrative="Perbandingan matriks memperlihatkan pergeseran paradigma kualitas pengujian. Unit test konvensional sangat baik untuk memvalidasi kontrak fungsional spesifik dan regresi skenario bisnis tertentu. Namun, untuk sistem yang menangani kalkulasi finansial, protokol serialisasi, atau manipulasi struktur data kompleks, PBT memberikan jaminan kebenaran yang jauh lebih kokoh melalui eksplorasi batas domain secara otomatis dan reduksi kegagalan yang deterministik.",
    narrative_en="The comparison matrix underscores a paradigm shift in test quality. While example tests excel at regression validation for explicit business user flows, PBT delivers mathematical confidence for financial calculations, serialization protocols, and complex data structures through automated exploration.",
    mental_model="Use example tests for concrete user flows; use property tests for structural correctness.",
    source_path="cmd/demo/main.go:215-245",
    slide_num="21 / 24"
)

# Slide 22: Praktik Terbaik & Anti-Pattern PBT
add_slide(
    categories=["BEST PRACTICE", "GUIDELINES"],
    title="Praktik Terbaik & Anti-Pattern PBT",
    title_en="Best Practices and Anti-Patterns in PBT",
    subtitle="Panduan menghindari jebakan umum dalam mendesain invariant dan generator.",
    content_items=[
        {"title": "YANG HARUS DIHINDARI", "items": [
            "1. Menulis ulang implementasi di dalam assertion test (duplikasi logika kode).",
            "2. Generator terlalu sempit sehingga mengabaikan batas domain ekstrem.",
            "3. Menguji tanpa shrinker sehingga menyulitkan proses debugging hasil gagal."
        ], "type": "bullet"},
        {"title": "POLA YANG DIREKOMENDASIKAN", "items": [
            "1. Fokus pada sifat relasional: roundtrip, idempotence, oracle, dan simetri.",
            "2. Suntikkan boundary bias: selalu sertakan 0, -1, MaxInt, dan slice kosong.",
            "3. Simpan seed yang gagal ke dalam unit test regresi berbasis contoh permanen."
        ], "type": "bullet"}
    ],
    narrative="Kesalahan paling sering terjadi saat mengadopsi PBT adalah kecenderungan developer menduplikasi logika fungsi yang diuji ke dalam assertion test. Praktik ini salah karena jika implementasi memiliki bug logika, assertion duplikat akan ikut menyetujuinya. Fokuslah murni pada sifat relasional independen. Selain itu, pastikan generator selalu mencakup nilai-nilai ekstrem (nol, batas maksimum integer, irisan kosong), dan simpan counterexample yang ditemukan sebagai unit test regresi permanen.",
    narrative_en="The most common PBT anti-pattern is re-implementing system logic inside test assertions, leading to shared blind spots. Focus strictly on relational invariants. Always bias generators toward boundary conditions (zeros, empty slices, integer limits), and save discovered counterexamples as permanent regression tests.",
    mental_model="Assert properties of the relationship, never re-implement the algorithm in the test.",
    source_path="labs/40-property-based-testing/content/02-master-draft.md",
    slide_num="22 / 24"
)

# Slide 23: Panduan Reproduksi Mandiri & Verifikasi
add_slide(
    categories=["PANDUAN", "REPRODUKSI"],
    title="Panduan Reproduksi Mandiri & Verifikasi",
    title_en="Hands-On Reproduction and Verification Guide",
    subtitle="Seluruh kode Go standar siap dijalankan dan diuji langsung pada terminal Anda.",
    content_items=[
        {"title": "LANGKAH EKSEKUSI CLI", "items": [
            "# 1. Masuk ke direktori lab",
            "cd labs/40-property-based-testing",
            "",
            "# 2. Jalankan seluruh property & unit test",
            "go test -v ./...",
            "",
            "# 3. Verifikasi ketiadaan race condition",
            "go test -race ./...",
            "",
            "# 4. Eksekusi CLI interactive demo",
            "go run ./cmd/demo"
        ], "type": "code"}
    ],
    narrative="Seluruh implementasi pada Lab 40 dibangun di atas Go 1.22 standar tanpa dependensi luar, menjamin reproduksibilitas penuh di lingkungan macOS, Linux, maupun Windows. Perintah go test -v ./... mengeksekusi suite unit test, uji properti testing/quick, serta algoritma shrinking. Perintah go run ./cmd/demo menampilkan perbandingan interaktif antara kegagalan float, perbandingan oracle interval, dan jejak reduksi 21 langkah secara real-time.",
    narrative_en="Lab 40 runs on standard Go 1.22 without third-party dependencies, guaranteeing portability across macOS, Linux, and Windows. Running go test -v ./... executes all unit tests, quick properties, and shrinking verification. Running go run ./cmd/demo displays real-time CLI comparisons and shrinking traces.",
    mental_model="Standard library tests provide reproducible, zero-dependency validation for engineering teams.",
    source_path="labs/40-property-based-testing/README.md • go.mod",
    slide_num="23 / 24"
)

# Slide 24: Closing Slide (Dark Navy `#0b1b2d`)
add_slide(
    categories=["CLOSING"],
    title="Eksperimen kecil, invariant universal terbukti.",
    title_en="Small experiment, universal invariants proven.",
    subtitle="",
    content_items=[
        {"title": "MODUL LAB", "items": [
            "internal/currency • internal/interval • internal/shrinker",
            "cmd/demo/main.go • README.md • go.mod",
            "Zero third-party dependencies (pure Go 1.22 stdlib)"
        ], "type": "bullet"},
        {"title": "VERIFIKASI & STATUS", "items": [
            f"Sumber: commit {COMMIT_SHA[:8]}",
            "Test passing & Race-Clean: 7/7 unit & property tests PASS",
            "Invarian terbukti: Roundtrip Currency, Idempotence Interval, dan 21-Step Shrinking"
        ], "type": "bullet"}
    ],
    mental_model="Stop guessing edge cases with static examples. Let invariants and shrinking find them for you.",
    source_path=REPO_URL,
    slide_num="24 / 24"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
