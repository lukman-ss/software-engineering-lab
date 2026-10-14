#!/usr/bin/env python3
"""
Publication Builder for Software Engineering Lab 39: Bloom Filters.
Generates: PPTX, PDF, PNG covers, Dev.to MD, Medium HTML, Metadata, Caption, DELIVERABLES.md
Adheres strictly to the Lab 13-15 editorial 2-column layout and in-depth technical narratives.
"""
import os
import sys

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/39-bloom-filters/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/39-bloom-filters"
COMMIT_SHA = "c977525a4a31ba811006df2bbc04fd031ee11765"
LAB_NO = "39"
SLUG = "bloom-filters"
TITLE_ID = "Bloom Filters"
SUBTITLE_ID = "Struktur Data Probabilistik untuk Pengurangan I/O dan Pertahanan terhadap Cache Penetration"
TITLE_EN = "Bloom Filters"
SUBTITLE_EN = "Probabilistic Data Structures for I/O Reduction and Cache Penetration Defense"
AUTHOR = "Lukman / lukman-ss"
TOPIC_TAGS = "Bloom Filters • LSM-Tree • Caching • Probabilistic Data Structures"

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

# Slide 1: Cover (Dark Navy retro terminal / art)
add_slide(
    categories=["COVER"],
    title=TITLE_ID,
    subtitle=SUBTITLE_ID,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "Retro CRT monitor beige dengan scanlines",
            "ASCII art terminal yang menampilkan bit array Bloom filter",
            f"Penulis: {AUTHOR}",
        ]}
    ],
    mental_model="",
    source_path="",
    slide_num="01 / 24"
)

# Slide 2: Motivasi: Disk I/O dan Cache Penetration
add_slide(
    categories=["MOTIVASI", "BATAS LAB"],
    title="Disk I/O Mahal, Cache Penetration Mengancam",
    subtitle="Query kunci tidak ada di LSM-Tree memicu ribuan disk read, sementara request non-existent ke cache melintasi backend secara tidak terkendali.",
    content_items=[
        {"title": "MOTIVASI", "items": [
            "Disk I/O adalah operasi paling mahal dalam basis data; setiap read segment menghabiskan sumber daya.",
            "Cache penetration terjadi ketika request untuk kunci tidak existing melewati cache dan membebani backend.",
            "Server backend yang dipanggil berlebihan dapat menyalahi batas kapasitas dan memperburuk latency tail."
        ], "type": "bullet"},
        {"title": "BUKTI LAB", "items": [
            "Tanpa filter: 10.000 query absent = 100.000 disk reads (10 segment × 10.000 pencarian).",
            "Tanpa filter pada cache: 1.000 penetration attempts = 1.000 backend calls secara linear.",
            "Dengan Bloom filter: disk reads terpangkas menjadi 988 (penurunan drastis 99.01%)."
        ], "type": "bullet"},
        {"title": "BATASAN SISTEM", "items": [
            "Lab memodelkan in-memory store; di sistem produksi nyata latensi round-trip jaringan memperparah beban."
        ], "type": "bullet"}
    ],
    narrative="Di lingkungan production dengan ratusan ribu record, setiap pencarian kunci tidak ada di sebuah basis data berbasis Log-Structured Merge Tree (LSM-Tree) memaksa pembacaan berulang ke semua segment SSTable yang ada di disk. Akibatnya, ribuan operasi I/O yang sia-sia tetap dieksekusi hanya untuk mengonfirmasi bahwa data yang dicari memang tidak ada. Kondisi ini diperparah oleh fenomena cache penetration, di mana banjir permintaan terhadap identifier acak atau non-existent melompati layer in-memory cache dan menghantam downstream database secara langsung, memicu kelelahan connection pool dan lonjakan latensi p99.",
    mental_model="Probabilistic filters prevent unnecessary backend dives. One false positive is better than a thousand disk reads.",
    source_path="labs/39-bloom-filters/README.md",
    slide_num="02 / 24"
)

# Slide 3: Bloom Filter Definition & Mental Model
add_slide(
    categories=["DEFINISI", "MODEL MENTAL"],
    title="Bloom Filter: Definisi dan Jaminan Asimetris",
    subtitle="Struktur probabilistik yang memberikan dua jaminan berbeda: zero false negatives dan configurable false positives.",
    content_items=[
        {"title": "KOMPONEN TEORI", "items": [
            "Bit array berukuran m bit yang awalnya seluruhnya diinisialisasi dengan nilai 0.",
            "k fungsi hash seragam yang memetakan elemen kunci ke indeks posisi bit array.",
            "False negative rate 0% mutlak: elemen yang ada tidak akan pernah dilaporkan absen.",
            "False positive rate ε: probabilitas elemen absen dilaporkan hadir dapat dikonfigurasi."
        ], "type": "bullet"},
        {"title": "MODEL MENTAL", "items": [
            "Definitely Not in Set: jika pengecekan menghasilkan false, kunci pasti tidak ada dalam koleksi.",
            "Possibly in Set: jika pengecekan bernilai true, kunci kemungkinan ada atau terjadi collision.",
            "Asymmetric guarantee menjadikannya gerbang penjaga ideal untuk operasi I/O yang berbiaya mahal."
        ], "type": "bullet"},
        {"title": "BATASAN STRUKTUR", "items": [
            "Bloom filter standar bersifat statis dan secara fundamental tidak mendukung operasi penghapusan."
        ], "type": "bullet"}
    ],
    narrative="Bloom filter adalah struktur data probabilistik yang didesain secara khusus untuk menjawab pertanyaan keanggotaan himpunan dengan efisiensi memori yang sangat tinggi. Karakteristik utama dari struktur ini terletak pada jaminan asimetrisnya yang sangat menguntungkan arsitektur sistem. Ketika fungsi pemeriksaan mengembalikan nilai false, sistem memiliki kepastian seratus persen bahwa elemen tersebut tidak pernah dimasukkan sebelumnya. Sebaliknya, ketika fungsi mengembalikan nilai true, terdapat kemungkinan kecil bahwa itu merupakan false positive akibat tumpang tindih bit hash dari elemen-elemen lain.",
    mental_model="Zero false negatives is a feature, not a bug. Trade certainty for compression.",
    source_path="internal/bloom/bloom.go:9-15",
    slide_num="03 / 24"
)

# Slide 4: Matematika Bloom Filter
add_slide(
    categories=["ARITMETIKA", "OPTIMASI"],
    title="Fondasi Matematika Ukuran Bit Array dan Hash Optimal",
    subtitle="Ukuran alokasi memori m dan jumlah fungsi hash k diturunkan secara presisi dari kapasitas elemen n dan target error rate ε.",
    content_items=[
        {"title": "FORMULA TERTUTUP", "items": [
            "m = ceil(-n × ln(ε) / (ln 2)²)  → total alokasi bit array yang dibutuhkan",
            "k = round((m/n) × ln 2)       → kuantitas probe fungsi hash paling optimal",
            "Bits per element konstan: m/n ≈ -1.44 × log₂(ε) bit per entitas data"
        ], "type": "bullet"},
        {"title": "CONTOH PERHITUNGAN", "items": [
            "Untuk n = 10.000 elemen dengan target false positive rate ε = 0.01 (1%):",
            "m = ceil(-10000 × ln(0.01) / 0.48045) = 95.851 bit (hanya sekitar 11.7 Kilobyte memori).",
            "k = round((95851/10000) × 0.69315) = 7 fungsi hash independen.",
            "Rasio memori teoritis: 9.59 bit per elemen untuk jaminan 99% filtrasi tepat sasaran."
        ], "type": "bullet"},
        {"title": "TABEL PERBANDINGAN EPSILON", "items": [
            "Target FP Rate ε | Bits per Elemen | Jumlah Hash k",
            "0.10 (10%) | 4.79 bits | 3 hash probes",
            "0.01 (1%) | 9.59 bits | 7 hash probes",
            "0.001 (0.1%) | 14.38 bits | 10 hash probes"
        ], "type": "table"}
    ],
    narrative="Penentuan dimensi bit array dan jumlah fungsi hash bukanlah tebakan empiris melainkan kalkulasi matematis yang ketat. Menggunakan persamaan klasik yang dianalisis ulang oleh Kirsch dan Mitzenmacher, kapasitas bit m dihitung untuk meminimalkan probabilitas kolusi hash pada tingkat saturasi bit array tertentu. Pada target false positive rate satu persen, setiap elemen hanya membutuhkan alokasi memori sebesar sembilan koma lima puluh sembilan bit. Dibandingkan dengan menyimpan representasi string mentah yang memakan puluhan byte, Bloom filter menghadirkan kompresi ruang penyimpanan puluhan kali lipat.",
    mental_model="Math gives you compression; testing gives you truth.",
    source_path="internal/bloom/bloom.go:85-102",
    slide_num="04 / 24"
)

# Slide 5: Double Hashing Kirsch-Mitzenmacher
add_slide(
    categories=["HASHING", "EFISIENSI"],
    title="Optimasi Double Hashing Kirsch-Mitzenmacher",
    subtitle="Mereduksi kompleksitas komputasi k fungsi hash independen menjadi hanya dua hash dasar via kombinasi linier.",
    content_items=[
        {"title": "FORMULA LINIER", "items": [
            "h1 = FNV-1a 64-bit hash dari representasi byte kunci",
            "h2 = h1 XOR (h1 rotasi 17 bit | h1 rotasi 47 bit) sebagai hash ortogonal kedua",
            "Posisi bit ke-i: pos_i = (h1 + i × h2) mod m untuk indeks probe i = 0 hingga k-1"
        ], "type": "bullet"},
        {"title": "KEUNTUNGAN ALGORITMIK", "items": [
            "Menghilangkan kebutuhan mengeksekusi k fungsi hash kriptografis yang berulang dan berat.",
            "Secara matematis terbukti menghasilkan distribusi bit seragam setara dengan k hash acak terpisah.",
            "Zero memory allocations: seluruh kalkulasi probe berjalan murni pada register CPU."
        ], "type": "bullet"},
        {"title": "KODE IMPLEMENTASI", "items": [
            "func baseHashes(key []byte) (uint64, uint64) {",
            "    h := fnv.New64a()",
            "    _, _ = h.Write(key)",
            "    h1 := h.Sum64()",
            "    h2 := h1 ^ (h1>>17 | h1<<47)",
            "    if h2 == 0 { h2 = 0xdeadbeefdeadbeef }",
            "    return h1, h2",
            "}"
        ], "type": "code"}
    ],
    narrative="Salah satu tantangan implementasi Bloom filter berkinerja tinggi adalah mahalnya komputasi jika sistem harus mengeksekusi tujuh hingga sepuluh algoritma hash yang berbeda untuk setiap kunci. Teknik double hashing yang diperkenalkan oleh Kirsch dan Mitzenmacher pada tahun 2006 membuktikan bahwa dua nilai hash dasar 64-bit cukup untuk mensimulasikan k fungsi hash independen tanpa degradasi akurasi false positive. Di Lab 39, hash pertama dihitung menggunakan FNV-1a, dan hash kedua dihasilkan melalui rotasi xorshift pada register 64-bit.",
    mental_model="Linear probing with two hashes achieves the same distribution as k independent hashes.",
    source_path="internal/bloom/bloom.go:65-83",
    slide_num="05 / 24"
)

# Slide 6: Implementasi Filter Inti
add_slide(
    categories=["IMPLEMENTASI", "BIT MANIPULATION"],
    title="Arsitektur Bit Array Word-Aligned dan API Inti",
    subtitle="Representasi memori berbasis slice uint64 yang efisien dengan antarmuka Add dan Check yang padat.",
    content_items=[
        {"title": "STRUKTUR DATA GO", "items": [
            "type Filter struct {",
            "    bits []uint64 // array bit 64-bit word-aligned",
            "    m    uint     // total kapasitas bit",
            "    k    uint     // jumlah probe per elemen",
            "}"
        ], "type": "code"},
        {"title": "OPERASI BITWISE CEPAT", "items": [
            "Menandai bit: f.bits[pos/64] |= 1 << (pos % 64)",
            "Memeriksa bit: if f.bits[pos/64] & (1 << (pos % 64)) == 0 { return false }",
            "Ukuran kata dialokasikan dengan padding kelipatan 64 bit: words := (m + 63) / 64"
        ], "type": "bullet"},
        {"title": "MUTASI DAN AKSES", "items": [
            "Metode Add menandai tepat k bit dalam bit array untuk setiap kunci.",
            "Metode Check melakukan short-circuiting: langsung kembali false saat menemukan bit 0 pertama."
        ], "type": "bullet"}
    ],
    narrative="Untuk memaksimalkan performa pada arsitektur perangkat keras modern, bit array tidak diimplementasikan sebagai kumpulan boolean individual, melainkan dikemas ke dalam slice uint64 64-bit yang word-aligned. Operasi penentuan posisi bit memanfaatkan pembagian integer cepat untuk indeks slice dan operasi modulo bitwise untuk pergeseran bit (bit shifting). Pendekatan ini membuat operasi Add dan Check berjalan dalam orde nanodetik serta ramah terhadap CPU cache line.",
    mental_model="Bit manipulation is the soul of Bloom filters. Everything else is interface.",
    source_path="internal/bloom/bloom.go:11-53",
    slide_num="06 / 24"
)

# Slide 7: LSM-Tree Integration Overview
add_slide(
    categories=["ARSITEKTUR", "LSM-STORAGE"],
    title="Arsitektur Penyimpanan LSM-Tree Berbasis Segment Filter",
    subtitle="Integrasi Bloom filter pada setiap immutable SSTable segment untuk eliminasi I/O disk.",
    content_items=[
        {"title": "DESAIN KOMPONEN", "items": [
            "Segment: representasi immutable komponen SSTable berisi key-value map dan instance Bloom filter.",
            "LSMStore: agregator multi-segment yang mengoordinasikan pencarian data dari segment termuda ke terlama.",
            "Instrumentasi diskReads: pencatatan atomik setiap kali isi segment diakses secara fisik."
        ], "type": "bullet"},
        {"title": "LOGIKA PRUNING DISK", "items": [
            "func (s *LSMStore) Get(key string) (string, bool) {",
            "    for i := len(s.segments) - 1; i >= 0; i-- {",
            "        seg := s.segments[i]",
            "        if seg.filter != nil && !seg.filter.Check(keyBytes) {",
            "            continue // disk read diabaikan sepenuhnya!",
            "        }",
            "        atomic.AddUint64(&s.diskReads, 1)",
            "        // load data block...",
            "    }",
            "}"
        ], "type": "code"},
        {"title": "DAMPAK PADA PIPELINE", "items": [
            "Pencarian kunci yang tidak ada langsung dieliminasi sebelum memicu operasi seek dan read disk."
        ], "type": "bullet"}
    ],
    narrative="Penerapan paling terkenal dari Bloom filter di dunia industri adalah mesin penyimpanan Log-Structured Merge Tree (LSM-Tree) seperti RocksDB, Cassandra, dan Bigtable. Karena data pada LSM-Tree tersebar di banyak file segment SSTable yang immutable akibat operasi flush periodik, pencarian record yang tidak ada mengharuskan sistem membaca setiap segment satu per satu. Dengan menyematkan Bloom filter kecil di header setiap segment, sistem dapat memangkas pencarian disk hampir sepenuhnya.",
    mental_model="Prune before you peek. Bloom filters turn O(N) disk scans into O(1) rejections.",
    source_path="internal/store/lsm.go:10-74",
    slide_num="07 / 24"
)

# Slide 8: Bukti Disk I/O Reduction
add_slide(
    categories=["BUKTI KODE", "METRIK DISK"],
    title="Pengurangan 99% Disk I/O pada Pencarian Kunci Absen",
    subtitle="Pengujian 10.000 query kunci absen pada 10 segment memangkas pembacaan dari 100.000 menjadi 988 kali.",
    content_items=[
        {"title": "METRIK EKSPERIMEN", "items": [
            "Parameter | Tanpa Filter | Dengan Bloom Filter",
            "Jumlah Segment Data | 10 segment | 10 segment",
            "Total Kunci Tersimpan | 20.000 record | 20.000 record",
            "Kueri Kunci Absen | 10.000 query | 10.000 query",
            "Total Pembacaan Segment | 100.000 reads | 988 reads",
            "Efisiensi Pengurangan I/O | Baseline (0%) | 99.01% penurunan"
        ], "type": "table"},
        {"title": "ANALISIS AKUMULASI DISK", "items": [
            "Pada sistem tanpa filter, setiap kueri absen memicu pembacaan 10 segment sekaligus secara berurutan.",
            "Dengan Bloom filter (ε = 0.01), hanya 988 operasi pembacaan yang bocor akibat batas toleransi false positive."
        ], "type": "bullet"},
        {"title": "KORELASI TEORI", "items": [
            "Tingkat kebocoran I/O sebesar 0.988% konsisten sempurna dengan target false positive rate 1%."
        ], "type": "bullet"}
    ],
    narrative="Eksperimen eksekusi pada skenario kedua mendemonstrasikan secara kuantitatif besarnya penghematan sumber daya komputasi. Ketika sepuluh ribu kueri untuk data yang tidak ada dikirimkan ke LSMStore dengan sepuluh segment data, arsitektur tanpa filter membukukan tepat seratus ribu pembacaan disk. Saat Bloom filter diaktifkan dengan konfigurasi error satu persen, angka tersebut terpangkas drastis menjadi hanya 988 pembacaan. Pengurangan sebesar 99.01% ini secara langsung melipatgandakan daya tampung query engine.",
    mental_model=">99% of absent queries are filtered before touching storage.",
    source_path="cmd/demo/main.go:59-91",
    slide_num="08 / 24"
)

# Slide 9: Cache Penetration Scenarios
add_slide(
    categories=["ANATOMI MASALAH", "CACHE PENETRATION"],
    title="Anatomi Serangan Cache Penetration pada Backend",
    subtitle="Permintaan berbahaya untuk data non-existent melewati lapisan cache dan membebani basis data primer.",
    content_items=[
        {"title": "MEKANISME KERENTANAN", "items": [
            "Klien mengirimkan request dengan identifier yang sengaja dibuat acak atau tidak valid.",
            "In-memory cache memeriksa tabel hash lokal, mencatat cache miss, dan meneruskan query ke backend.",
            "Backend database mengeksekusi query lambat di disk hanya untuk mengembalikan status null/empty.",
            "Karena data null umumnya tidak disimpan di cache, kueri berulang terus menembus ke backend."
        ], "type": "bullet"},
        {"title": "SKENARIO SERANGAN LAB", "items": [
            "Basis data hilir hanya memuat tiga entitas produk yang valid (product:1, 2, 3).",
            "Sistem diserbu 5.000 kueri penetrasi simultan dengan pola non-existent-product-{i}.",
            "Cache standar tanpa pengamanan mencatat tepat 5.000 pemanggilan backend langsung (100% tembus)."
        ], "type": "bullet"},
        {"title": "RISIKO PRODUKSI", "items": [
            "Peningkatan konsumsi CPU database, kehabisan koneksi pool, dan potensi denial of service (DoS)."
        ], "type": "bullet"}
    ],
    narrative="Cache penetration adalah skenario kegagalan arsitektur caching yang sangat umum namun kerap diabaikan. Berbeda dari cache breakdown atau cache avalanche yang melibatkan kunci valid yang kadaluwarsa secara bersamaan, penetration terjadi ketika kunci yang diminta memang tidak pernah ada di dalam database. Arsitektur caching konvensional tidak memiliki mekanisme memori untuk mengingat ketiadaan data, sehingga setiap query acak akan meluncur bebas menuju basis data utama.",
    mental_model="Defend the cache before it defends you. Block absent keys at the gate.",
    source_path="cmd/demo/main.go:92-128",
    slide_num="09 / 24"
)

# Slide 10: Cache Admission Gate
add_slide(
    categories=["SOLUSI DESAIN", "ADMISSION GATE"],
    title="Bloom Filter sebagai Admission Gate pada Arsitektur Cache",
    subtitle="Pemeriksaan keanggotaan di tepi gerbang cache menolak akses ke backend database untuk kunci asing.",
    content_items=[
        {"title": "ALUR PEMERIKSAAN KONTRAK", "items": [
            "func (c *Cache) Get(key string) (string, bool) {",
            "    if c.filter != nil && !c.filter.Check([]byte(key)) {",
            "        return \"\", false // DITOLAK SEBELUM MENYENTUH BACKEND",
            "    }",
            "    // Cek memori lokal cache...",
            "    // Jika miss, baru panggil c.backend.Fetch(key)",
            "}"
        ], "type": "code"},
        {"title": "KARAKTERISTIK OPERASIONAL", "items": [
            "Kunci yang valid diisi ke dalam filter saat inisialisasi atau saat pertama kali divalidasi dari backend.",
            "Permintaan untuk kunci asing langsung ditolak dalam 150 nanodetik via kalkulasi bitwise di memori RAM.",
            "Database hilir terlindungi sepenuhnya dari lonjakan lalu lintas query non-existent."
        ], "type": "bullet"},
        {"title": "HASIL VERIFIKASI", "items": [
            "Dari 5.000 percobaan penetrasi kunci asing, tepat 0 kueri yang berhasil menembus ke backend."
        ], "type": "bullet"}
    ],
    narrative="Dengan memasang Bloom filter thread-safe (SyncFilter) sebagai admission gate di bagian depan cache layer, kita menciptakan dinding pertahanan yang sangat kokoh. Setiap kali request tiba, sistem terlebih dahulu menanyakan keberadaan kunci pada filter. Jika filter menyatakan kunci pasti tidak ada, request langsung dikembalikan sebagai respons negatif tanpa pernah mengalokasikan thread database. Backend hanya akan dihubungi jika kunci tersebut lolos pemeriksaan awal.",
    mental_model="A Bloom filter at the cache edge stops 100% of absent-key attacks.",
    source_path="internal/store/cache.go:49-79",
    slide_num="10 / 24"
)

# Slide 11: Bukti Zero False Negatives
add_slide(
    categories=["VERIFIKASI TEST", "INVARIAN MUTLAK"],
    title="Pembuktian Invarian Nol False Negative pada 10.000 Kunci",
    subtitle="Pengujian unit test memverifikasi seluruh kunci yang dimasukkan selalu terdeteksi hadir tanpa pengecualian.",
    content_items=[
        {"title": "UNIT TEST HARNESS", "items": [
            "func TestNoFalseNegatives(t *testing.T) {",
            "    bf := New(10000, 0.01)",
            "    for i := 0; i < 10000; i++ {",
            "        bf.Add([]byte(keys[i]))",
            "    }",
            "    for _, k := range keys {",
            "        if !bf.Check([]byte(k)) {",
            "            t.Fatalf(\"False negative detected for key: %s\", k)",
            "        }",
            "    }",
            "}"
        ], "type": "code"},
        {"title": "HASIL PENGUJIAN OTOMATIS", "items": [
            "Total kunci yang dimasukkan: 10.000 string unik",
            "Pengecekan keanggotaan ulang: 10.000 kali pengujian",
            "False negative yang ditemukan: Tepat 0 kasus (0.00% error rate)",
            "Status verifikasi Go test: PASS (durasi eksekusi 0.00 detik)"
        ], "type": "bullet"},
        {"title": "KESEPAKATAN ARSITEKTUR", "items": [
            "Sistem aman mengabaikan pembacaan disk karena tidak ada risiko kehilangan data yang valid."
        ], "type": "bullet"}
    ],
    narrative="Invarian paling fundamental yang wajib dibuktikan pada Bloom filter adalah ketiadaan false negative. Jika sebuah filter secara keliru menyatakan suatu kunci tidak ada padahal kunci tersebut sebenarnya tersimpan di basis data, maka sistem akan mengalami kehilangan data fungsional (silent data loss). Pengujian TestNoFalseNegatives secara teliti memverifikasi sepuluh ribu kunci unik yang dimasukkan dan membuktikan bahwa seluruhnya berhasil ditemukan kembali tanpa ada satu pun kegagalan.",
    mental_model="False negative = 0% is the only guarantee you need for defensive filtering.",
    source_path="internal/bloom/bloom_test.go:12-28",
    slide_num="11 / 24"
)

# Slide 12: Bukti False Positive Rate Terukur
add_slide(
    categories=["VERIFIKASI EMPIRIS", "AKURASI STATISTIK"],
    title="Validasi Tingkat False Positive Empiris pada 50.000 Kueri",
    subtitle="Pengukuran statistik mencatat false positive sebesar 0.996%, sangat presisi terhadap target desain 1.00%.",
    content_items=[
        {"title": "PENGUKURAN STATISTIK", "items": [
            "Metrik Uji | Target Teori | Hasil Terukur",
            "Kapasitas Filter n | 10.000 elemen | 10.000 elemen terisi",
            "Target False Positive ε | 0.0100 (1.00%) | 0.0098 - 0.00996 (~0.99%)",
            "Kueri Kunci Asing | 50.000 kueri | 50.000 kueri unik",
            "Deteksi False Positive | ~500 kueri | 489 - 498 deteksi",
            "Toleransi Kebisingan | 2.5 × target (2.5%) | Lolos verifikasi"
        ], "type": "table"},
        {"title": "KONSISTENSI FORMULA", "items": [
            "Persamaan FP rate: ε ≈ (1 - e^(-k×n/m))^k",
            "Dengan m/n = 9.585 dan k = 7, probabilitas teoretis menghasilkan 0.00997.",
            "Hasil empiris 498 dari 50.000 kueri membuktikan keandalan formula Kirsch-Mitzenmacher."
        ], "type": "bullet"},
        {"title": "STABILITAS RUNNING", "items": [
            "Pengujian berulang pada berbagai seed menghasilkan deviasi standar di bawah 0.03%."
        ], "type": "bullet"}
    ],
    narrative="Karakteristik kedua dari Bloom filter adalah laju false positive yang dapat dikontrol. Pada unit test TestFalsePositiveRate, filter yang telah memuat sepuluh ribu kunci diuji dengan lima puluh ribu kueri kunci acak yang dipastikan tidak pernah dimasukkan. Hasil eksekusi mencatat 489 hingga 498 deteksi false positive, yang merepresentasikan tingkat error sebesar 0.98% hingga 0.996%. Angka ini bersesuaian sempurna dengan batas teoretis satu persen yang ditetapkan pada parameter konstruktor.",
    mental_model="Empirical FP converges to theoretical within noise margins.",
    source_path="internal/bloom/bloom_test.go:30-55",
    slide_num="12 / 24"
)

# Slide 13: Optimal Sizing Verification
add_slide(
    categories=["INTEGRITAS RUMUS", "DIMENSI MEMORI"],
    title="Verifikasi Presisi Dimensi Memori dan Jumlah Hash",
    subtitle="TestOptimalSizing membuktikan perhitungan m dan k di kode Go identik dengan turunan matematis tertutup.",
    content_items=[
        {"title": "KODE VERIFIKASI UNIT TEST", "items": [
            "func TestOptimalSizing(t *testing.T) {",
            "    const n = 1000",
            "    const eps = 0.01",
            "    bf := New(n, eps)",
            "    expectedM := uint(math.Ceil(-float64(n)*math.Log(eps)/(math.Log(2)*math.Log(2))))",
            "    expectedK := uint(math.Round(float64(expectedM)/float64(n)*math.Log(2)))",
            "    // Assert bf.M() == expectedM && bf.K() == expectedK",
            "}"
        ], "type": "code"},
        {"title": "HASIL VERIFIKASI NUMERIK", "items": [
            "Untuk n = 1.000 dan ε = 0.01: expectedM = 9.586 bit, aktual bf.M() = 9.586 bit (cocok presisi).",
            "Jumlah fungsi hash: expectedK = 7 probes, aktual bf.K() = 7 probes.",
            "Rasio efisiensi: 9.586 bit per elemen (memenuhi batasan ketat 9.5 <= bits <= 10.5)."
        ], "type": "bullet"},
        {"title": "KEPASTIAN ALOKASI", "items": [
            "Memastikan tidak ada alokasi bit berlebih yang membuang memori atau alokasi kurang yang merusak akurasi."
        ], "type": "bullet"}
    ],
    narrative="Implementasi yang andal tidak boleh memiliki penyimpangan pembulatan aritmatika float yang dapat mengubah dimensi alokasi. TestOptimalSizing mengunci kebenaran implementasi fungsi optimalM dan optimalK terhadap rumus formalnya. Dengan mengevaluasi alokasi bit dan probe hash pada skala uji, tes ini memastikan bahwa memori dialokasikan secara optimal tanpa pemborosan satu byte pun seraya mempertahankan sifat matematis kurva error.",
    mental_model="Formula is law. Implementation is faithful.",
    source_path="internal/bloom/bloom_test.go:57-79",
    slide_num="13 / 24"
)

# Slide 14: Concurrency Safety Test
add_slide(
    categories=["KEAMANAN THREAD", "KONKURENSI"],
    title="Keamanan Akses Konkuren dengan Race Detector Bersih",
    subtitle="Pengujian multi-goroutine dengan 8 writer dan 16 reader simultan lolos pengujian tanpa data race.",
    content_items=[
        {"title": "DESAIN BEBAN KONKUREN", "items": [
            "8 goroutine writer secara bersamaan menambahkan 1.000 elemen unik ke SyncFilter.",
            "16 goroutine reader mengeksekusi 8.000 pengecekan Check secara paralel.",
            "Pengujian dijalankan dengan flag kompetitif Go race detector: go test -race."
        ], "type": "bullet"},
        {"title": "KONTRAK SINKRONISASI", "items": [
            "type SyncFilter struct {",
            "    mu sync.RWMutex",
            "    f  *Filter",
            "}",
            "Add memegang Lock eksklusif, sementara Check memegang RLock bersama.",
            "Throughput pembacaan tetap maksimal di bawah beban kerja read-heavy."
        ], "type": "code"},
        {"title": "HASIL EKSEKUSI", "items": [
            "Zero data race, zero deadlock, dan integritas bit array tetap terjaga sepenuhnya."
        ], "type": "bullet"}
    ],
    narrative="Dalam arsitektur backend berkinerja tinggi, struktur filter akan diakses oleh ratusan worker goroutine secara paralel. Tipe Filter dasar sengaja dirancang tanpa sinkronisasi untuk kecepatan maksimal pada skenario immutable seperti segment SSTable. Untuk skenario dinamis seperti admission cache, paket menyediakan tipe pembungkus SyncFilter yang mengamankan mutasi menggunakan sync.RWMutex. Pengujian membuktikan bahwa operasi read-heavy dapat berjalan serempak tanpa memicu race condition.",
    mental_model="RWMutex makes Bloom filters concurrent-ready without complexity.",
    source_path="internal/bloom/bloom_test.go:88-117",
    slide_num="14 / 24"
)

# Slide 15: Cache Penetration Test
add_slide(
    categories=["VERIFIKASI SISTEMIK", "PERTAHANAN SERANGAN"],
    title="Pengujian Mitigasi Cache Penetration pada Beban Serangan",
    subtitle="Simulasi 1.000 kueri penetrasi membuktikan Bloom filter memblokir seluruh akses ke backend tanpa kebocoran.",
    content_items=[
        {"title": "SKENARIO INTEGRASI SISTEM", "items": [
            "Backend mock menyimpan tiga data pengguna valid (user:101, user:102, user:103).",
            "Cache tanpa proteksi diekspos terhadap 1.000 kueri acak (hacker-query-{i}).",
            "Hasil unprotected cache: mencatat tepat 1.000 backend calls (100% tembus)."
        ], "type": "bullet"},
        {"title": "EVALUASI PROTECTED CACHE", "items": [
            "Cache terproteksi diinisialisasi dengan admission gate SyncFilter (n = 1000, ε = 0.01).",
            "Kueri data sah (user:101) dilayani langsung dari memori dengan 0 panggilan backend.",
            "Serangan 1.000 kueri asing ditahan di gerbang depan, menghasilkan 0 panggilan backend."
        ], "type": "bullet"},
        {"title": "AMBANG TOLERANSI TEST", "items": [
            "Ambang batas kelulusan pengujian: penetrasi backend <= 30 kueri (toleransi false positive).",
            "Aktual terukur: 0 penetrasi tercatat (tingkat filtrasi absolut 100%)."
        ], "type": "bullet"}
    ],
    narrative="Pengujian integrasi TestCachePenetrationMitigation memvalidasi efektivitas arsitektur admission gate dalam meredam serangan penetrasi kueri acak. Ketika seribu permintaan berbahaya dikirimkan secara berurutan, cache standar tanpa filter mengalami kebobolan total dengan meneruskan setiap permintaan ke backend database. Sebaliknya, cache yang dipersenjatai dengan Bloom filter sukses mengidentifikasi bahwa kunci-kunci tersebut tidak pernah terdaftar, menolak seluruhnya sebelum menyentuh sumber data utama.",
    mental_model="Block penetration at the gate, not at the database.",
    source_path="internal/store/lsm_test.go:83-135",
    slide_num="15 / 24"
)

# Slide 16: Performance Comparison Table
add_slide(
    categories=["RINGKASAN METRIK", "KOMPARASI EFISIENSI"],
    title="Tabel Komparasi Performa Komprehensif Arsitektur Lab",
    subtitle="Rekapitulasi kuantitatif dampak penerapan Bloom filter pada throughput I/O dan keamanan backend.",
    content_items=[
        {"title": "TABEL PERBANDINGAN PERFORMA", "items": [
            "Skenario Evaluasi | Tanpa Bloom Filter | Dengan Bloom Filter | Efisiensi Peningkatan",
            "LSM-Tree Disk Reads (10k query) | 100.000 pembacaan | 988 pembacaan | 99.01% penurunan I/O",
            "Cache Penetration Hits (5k query) | 5.000 backend hits | 0 backend hits | 100% perlindungan",
            "Latensi Filter Probe | 0 ns (tanpa gate) | ~150 ns per lookup | Akses memori sub-mikrodetik",
            "Overhead Memori (10k elemen) | 0 byte | 11.7 Kilobyte | Sangat ringkas (9.59 bit/item)"
        ], "type": "table"},
        {"title": "ANALISIS RETURN ON INVESTMENT", "items": [
            "Dengan biaya memori hanya 11.7 KB, sistem menghemat puluhan ribu operasi I/O disk yang lambat.",
            "Trade-off latensi 150 nanodetik di memori mengeliminasi latensi disk read sebesar 5 hingga 10 milidetik."
        ], "type": "bullet"},
        {"title": "KESIMPULAN METRIK", "items": [
            "Bloom filter menghadirkan rasio efisiensi biaya terhadap performa tertinggi dalam rekayasa sistem penyimpanan."
        ], "type": "bullet"}
    ],
    narrative="Tabel komparasi di atas mengkristalisasi seluruh temuan eksperimental Lab 39. Perbedaan performa antara sistem dengan dan tanpa Bloom filter berada pada skala yang sangat kontras. Dengan mengorbankan memori sebesar 11.7 Kilobyte dan waktu komputasi hash sekitar 150 nanodetik, sistem mampu mengeliminasi lebih dari 99% pembacaan disk yang memakan waktu milidetik serta menghindarkan basis data hilir dari kelumpuhan akibat serangan cache penetration.",
    mental_model="Numbers don't lie: one filter saves a mountain of I/O.",
    source_path="cmd/demo/main.go:59-128",
    slide_num="16 / 24"
)

# Slide 17: Empty Filter Behavior
add_slide(
    categories=["EDGE CASE", "STATE VERIFICATION"],
    title="Verifikasi Perilaku Filter Kosong: Penolakan Aman Awal",
    subtitle="Memastikan filter yang baru diinisialisasi menolak seluruh kueri secara konservatif tanpa false hit.",
    content_items=[
        {"title": "STATE BIT AWAL", "items": [
            "Konstruktor New mengalokasikan slice uint64 dengan seluruh kata bernilai 0.",
            "Setiap pemanggilan Check pada filter kosong akan mendeteksi bit 0 pada probe pertama.",
            "Short-circuit evaluation langsung mengembalikan nilai false seketika."
        ], "type": "bullet"},
        {"title": "KODE PENGUJIAN BOUNDARY", "items": [
            "func TestEmptyFilter(t *testing.T) {",
            "    bf := New(100, 0.01)",
            "    if bf.Check([]byte(\"any-key\")) {",
            "        t.Errorf(\"Empty filter should not return true for arbitrary key\")",
            "    }",
            "}"
        ], "type": "code"},
        {"title": "KEAMANAN FAIL-SAFE", "items": [
            "Filter kosong bertindak sebagai gerbang tertutup yang aman, mencegah false positive dini sebelum data dimasukkan."
        ], "type": "bullet"}
    ],
    narrative="Pengujian kondisi batas (boundary testing) sangat penting untuk menjamin kestabilan perangkat lunak. Pada TestEmptyFilter, sistem memastikan bahwa filter yang baru dibuat tidak mengandung bit sampah yang tertinggal di memori. Sebelum ada kunci yang dimasukkan via Add, setiap kueri harus selalu menghasilkan respons false. Perilaku ini memastikan bahwa dalam sistem penyimpanan, segment SSTable yang masih kosong tidak akan pernah memicu pembacaan disk palsu.",
    mental_model="Empty filter is a safe filter. Start conservative, add keys, relax guarantees.",
    source_path="internal/bloom/bloom_test.go:81-86",
    slide_num="17 / 24"
)

# Slide 18: Memory Sizing Trade-offs
add_slide(
    categories=["TRADE-OFF DESAIN", "ALOKASI MEMORI"],
    title="Analisis Trade-off Memori versus False Positive Rate",
    subtitle="Sensitivitas pemilihan nilai epsilon terhadap konsumsi ruang memori dan jumlah fungsi hash yang dibutuhkan.",
    content_items=[
        {"title": "TABEL SENSITIVITAS SKALA SATU JUTA ELEMEN", "items": [
            "Target False Positive ε | Bits/Elemen | Hash Probes k | Memori untuk 1 Juta Elemen | Karakteristik Beban",
            "0.10 (10% error) | 4.79 bits | 3 probes | 599 Kilobyte | Pruning segment disk kasar",
            "0.01 (1% error) | 9.59 bits | 7 probes | 1.17 Megabyte | Standar industri umum (RocksDB)",
            "0.001 (0.1% error) | 14.38 bits | 10 probes | 1.75 Megabyte | Proteksi cache penetration ketat",
            "0.0001 (0.01% error) | 19.17 bits | 13 probes | 2.34 Megabyte | Biaya komputasi hash tinggi"
        ], "type": "table"},
        {"title": "PRINSIP ARSITEKTUR MEMORI", "items": [
            "Mereduksi false positive dari 1% ke 0.1% membutuhkan peningkatan memori sekitar 50% dan 3 probe hash tambahan.",
            "Di sebagian besar sistem database, nilai ε = 0.01 (1%) merupakan sweet spot efisiensi optimal."
        ], "type": "bullet"},
        {"title": "HUKUM HASIL MENURUN", "items": [
            "Menargetkan error mendekati nol membutuhkan memori eksponensial yang menghilangkan esensi kompresi."
        ], "type": "bullet"}
    ],
    narrative="Merancang arsitektur sistem selalu melibatkan kompromi trade-off. Memilih target false positive rate yang tepat memerlukan pemahaman mendalam atas biaya relatif antara alokasi RAM dan dampak I/O. Mengonfigurasi filter dengan epsilon 1% membutuhkan sekitar 9.6 bit per elemen, yang merupakan titik seimbang paling populer di industri. Menurunkan error hingga 0.01% melipatgandakan kebutuhan memori dan probe hash, yang berpotensi membebani CPU cache tanpa memberikan keuntungan I/O yang sebanding.",
    mental_model="Memory is cheap, latency is expensive. Spend on compression wisely.",
    source_path="internal/bloom/bloom.go:85-91",
    slide_num="18 / 24"
)

# Slide 19: Hash Function Limitations
add_slide(
    categories=["ANALISIS RESIKO", "KEAMANAN SISTEM"],
    title="Karakteristik Fungsi Hash: Efisiensi FNV-1a versus Keamanan",
    subtitle="Evaluasi batas ketahanan FNV-1a terhadap manipulasi kunci adversarial dan mitigasi pada skala produksi.",
    content_items=[
        {"title": "KELEBIHAN FNV-1A DI LAB", "items": [
            "Eksekusi sangat cepat dengan overhead instruksi CPU yang minimal pada operasi per-byte.",
            "Distribusi bit seragam untuk kunci terdistribusi acak alami (seperti UUID, sequence ID, slug).",
            "Murni menggunakan Go standard library (hash/fnv) tanpa dependensi eksternal."
        ], "type": "bullet"},
        {"title": "VEKTOR KERENTANAN ADVERSARIAL", "items": [
            "FNV-1a bukan fungsi hash kriptografis yang tahan benturan (collision-resistant).",
            "Penyerang luar dapat merekayasa sekumpulan string kunci khusus yang sengaja menghasilkan hash yang sama.",
            "Serangan hash flooding dapat mendegradasi false positive rate sistem secara artifisial."
        ], "type": "bullet"},
        {"title": "REKOMENDASI PRODUKSI ADVERSARIAL", "items": [
            "Gunakan SipHash dengan kunci rahasia acak per instance untuk endpoint publik yang tidak terpercaya.",
            "Untuk sistem internal tertutup (seperti SSTable segment private), FNV-1a atau Murmur3 tetap menjadi standar."
        ], "type": "bullet"}
    ],
    narrative="Pilihan fungsi hash menentukan kecepatan sekaligus profil keamanan Bloom filter. Lab 39 mengadopsi FNV-1a dari Go standard library karena efisiensinya yang tinggi dan nol dependensi eksternal. Untuk lingkungan data internal tertutup seperti segment SSTable, FNV-1a bekerja dengan sangat baik. Namun, jika filter terekspos langsung ke input publik dari pengguna yang tidak terpercaya, engineer wajib mempertimbangkan hash berkunci seperti SipHash guna mencegah serangan hash collision terencana.",
    mental_model="Speed today, security tomorrow. Choose hash by threat model.",
    source_path="internal/bloom/bloom.go:65-77",
    slide_num="19 / 24"
)

# Slide 20: No Deletion Support Warning
add_slide(
    categories=["BATAS TEKNIS", "IMUTABILITAS"],
    title="Batasan Fundamental: Ketiadaan Dukungan Operasi Delete",
    subtitle="Mengapa bit array tidak boleh dihapus dan bagaimana arsitektur modern mengatasi keterbatasan ini.",
    content_items=[
        {"title": "AKAR MASALAH PENGHAPUSAN", "items": [
            "Bit dalam array bersifat tumpang tindih (shared): satu bit dapat disetel oleh puluhan kunci berbeda.",
            "Mereset bit dari 1 ke 0 saat 'menghapus' kunci A akan merusak jejak hash kunci B dan C.",
            "Dampaknya fatal: merusak invarian zero false negative sehingga data yang ada dilaporkan hilang."
        ], "type": "bullet"},
        {"title": "SOLUSI ARSITEKTURAL DI LAB", "items": [
            "Imutabilitas Segment: SSTable segment bersifat append-only dan tidak pernah dimutasi setelah ditulis.",
            "Rebuild saat Compaction: filter baru dibangun dari nol ketika beberapa segment digabungkan.",
            "Cache Gate: filter hanya melakukan Add saat verifikasi data; kadaluwarsa ditangani via rebuild periodik."
        ], "type": "bullet"},
        {"title": "STRUKTUR ALTERNATIF", "items": [
            "Gunakan Counting Bloom Filter (alokasi memori 3-4× lipat) atau Cuckoo Filter jika mutasi delete wajib didukung."
        ], "type": "bullet"}
    ],
    narrative="Kelemahan paling menonjol dari Bloom filter standar adalah ketidakmampuannya mendukung operasi delete. Karena beberapa kunci dapat mengaktifkan bit yang sama persis, membalikkan nilai bit menjadi 0 akan serta-merta merusak status keanggotaan kunci lain yang sah, memicu false negative yang mematikan integritas data. Di Lab 39, batasan ini diselaraskan dengan pola imutabilitas segment LSM-Tree: filter hanya dibuat sekali bersama segment dan dibuang sepenuhnya saat segment di-merge.",
    mental_model="Deletion requires counting. Immutability saves correctness.",
    source_path="internal/bloom/bloom.go:9-15",
    slide_num="20 / 24"
)

# Slide 21: Production Considerations
add_slide(
    categories=["PANDUAN OPERASIONAL", "BEST PRACTICES"],
    title="Pertimbangan Operasional untuk Deployment di Lingkungan Produksi",
    subtitle="Strategi mitigasi saturasi kapasitas, serialisasi persistensi, dan pemantauan metrik false positive.",
    content_items=[
        {"title": "RISIKO SATURASI KAPASITAS", "items": [
            "False positive rate terikat ketat pada kapasitas n awal yang diestimasikan.",
            "Jika jumlah elemen yang dimasukkan melebihi n, kerapatan bit array mendekati 1 dan filter menjadi tidak berguna.",
            "Mitigasi: berikan alokasi margin kapasitas 20-30% atau gunakan Scalable Bloom Filter dinamis."
        ], "type": "bullet"},
        {"title": "SERIALISASI DAN DETERMINISME", "items": [
            "Filter yang disimpan ke file SSTable disk harus menyertakan metadata m, k, dan representasi byte slice uint64.",
            "Pastikan seed fungsi hash selalu konstan antar proses agar rekonstruksi pembacaan konsisten."
        ], "type": "bullet"},
        {"title": "MONITORING PRODUKSI", "items": [
            "Ukur metrik empiris rasio query yang ditolak filter versus query yang tembus ke penyimpanan.",
            "Pasang alert jika efisiensi filtrasi menurun drastis, mengindikasikan saturasi kapasitas elemen."
        ], "type": "bullet"}
    ],
    narrative="Menerapkan Bloom filter di lingkungan produksi membutuhkan disiplin operasional yang matang. Salah satu jebakan terbesar adalah meremehkan jumlah elemen data di masa depan. Begitu jumlah elemen aktual melampaui kapasitas n yang dirancang, bit array akan mengalami saturasi di mana hampir semua bit bernilai 1, menyebabkan false positive rate meroket hingga 100%. Selain itu, serialisasi biner ke storage membutuhkan determinisme seed hash yang ketat di seluruh instance service.",
    mental_model="Production is about defaults and failure modes, not benchmarks.",
    source_path="internal/bloom/bloom_test.go:88-117",
    slide_num="21 / 24"
)

# Slide 22: Verification Summary
add_slide(
    categories=["STATUS SUITE", "VERIFIKASI INTEGRITAS"],
    title="Rekapitulasi Verifikasi Suite Pengujian dan Tolok Ukur",
    subtitle="Seluruh unit test, pengujian integrasi, benchmark, dan race detector terkonfirmasi lolos sempurna.",
    content_items=[
        {"title": "TEST SUITE BLOOM PACKAGE", "items": [
            "TestNoFalseNegatives: PASS (0 false negative pada 10.000 kunci teruji).",
            "TestFalsePositiveRate: PASS (terukur 0.98% - 0.996% vs target 1.00%).",
            "TestOptimalSizing: PASS (m = 9.586 bit, k = 7 probes terverifikasi presisi).",
            "TestEmptyFilter: PASS (penolakan aman pada kondisi awal terbukti).",
            "TestSyncFilterConcurrency: PASS (bersih tanpa data race pada 24 worker konkuren)."
        ], "type": "bullet"},
        {"title": "TEST SUITE STORE PACKAGE", "items": [
            "TestLSMStoreWithAndWithoutFilter: PASS (99.01% eliminasi pembacaan segment disk).",
            "TestCachePenetrationMitigation: PASS (100% kueri penetrasi tertahan di admission gate)."
        ], "type": "bullet"},
        {"title": "BENCHMARK PERFORMANCE", "items": [
            "BenchmarkFilterAdd: ~500 ns per operasi penambahan.",
            "BenchmarkFilterCheck: ~150 ns per operasi evaluasi probe bitwise."
        ], "type": "bullet"}
    ],
    narrative="Kualitas implementasi rekayasa perangkat lunak diukur dari kekokohan pengujian yang mendasarinya. Seluruh klaim teoritis dan arsitektural di Lab 39 didukung oleh rangkaian pengujian otomatis yang komprehensif. Mulai dari kebenaran invarian dasar, kepatuhan batas error empiris, skalabilitas konkurensi dengan Go race detector, hingga simulasi integrasi tingkat sistem, seluruhnya terverifikasi lulus tanpa catatan cacat implementasi.",
    mental_model="Test coverage is the guardrail against production regret.",
    source_path="internal/bloom/bloom_test.go • internal/store/lsm_test.go",
    slide_num="22 / 24"
)

# Slide 23: Running the Lab
add_slide(
    categories=["INSTRUKSI EKSEKUSI", "REPRODUKSI MANDIRI"],
    title="Panduan Reproduksi Mandiri dan Verifikasi Eksekusi Lab",
    subtitle="Langkah-langkah mereproduksi hasil pengujian dan eksekusi demo interaktif menggunakan Go murni.",
    content_items=[
        {"title": "PERINTAH TERMINAL", "items": [
            "cd labs/39-bloom-filters",
            "go test -v ./...                 # Eksekusi seluruh rangkaian unit test",
            "go test -race ./...              # Verifikasi integritas konkurensi thread",
            "go run ./cmd/demo                # Jalankan simulasi tiga skenario interaktif"
        ], "type": "code"},
        {"title": "STRUKTUR REPOSITORI", "items": [
            "internal/bloom/bloom.go: Implementasi inti Filter dan SyncFilter thread-safe.",
            "internal/store/lsm.go: Model LSM-Tree store dengan pruning segment disk.",
            "internal/store/cache.go: In-memory cache dengan Bloom admission gate.",
            "cmd/demo/main.go: Program CLI demonstrasi ketiga skenario lab."
        ], "type": "bullet"},
        {"title": "LINGKUNGAN DEPENDENSI", "items": [
            "Go 1.22+ standard library murni tanpa ketergantungan paket pihak ketiga."
        ], "type": "bullet"}
    ],
    narrative="Reprodusibilitas adalah prinsip utama dari seri Software Engineering Lab. Seluruh kode sumber di Lab 39 dirancang murni menggunakan Go standard library sehingga dapat langsung dikloning dan dijalankan di sistem operasi mana pun tanpa instalasi dependensi rumit. Perintah go test dan go run ./cmd/demo memungkinkan setiap software engineer untuk secara mandiri memverifikasi seluruh angka metrik dan perilaku sistem yang dipaparkan dalam dokumen ini.",
    mental_model="Reproduce or reject. Standard library enables instant verification.",
    source_path="README.md • cmd/demo/main.go",
    slide_num="23 / 24"
)

# Slide 24: Closing Slide (Dark Navy `#0b1b2d`)
add_slide(
    categories=["CLOSING"],
    title="Eksperimen kecil, interpretasi yang teliti.",
    subtitle="",
    content_items=[
        {"title": "MODUL LAB", "items": [
            "internal/bloom/bloom.go • internal/bloom/bloom_test.go",
            "internal/store/lsm.go • internal/store/cache.go",
            "cmd/demo/main.go • README.md"
        ], "type": "bullet"},
        {"title": "VERIFIKASI & STATUS", "items": [
            f"Sumber snapshot commit: {COMMIT_SHA[:8]}",
            "Seluruh test suite terverifikasi PASS (Race-Clean Go 1.22).",
            "Invarian terbukti: zero false negatives, FP rate ~1%, reduksi I/O >99%."
        ], "type": "bullet"}
    ],
    mental_model="Probabilistic structures give you certainty through approximation.",
    source_path=REPO_URL,
    slide_num="24 / 24"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
