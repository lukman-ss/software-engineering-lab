#!/usr/bin/env python3
"""
Publication Builder for Software Engineering Lab 22: N+1 Query Problem.
Generates: PPTX, PDF, PNG covers, Dev.to MD, Medium HTML, Metadata, Caption, DELIVERABLES.md
Adheres strictly to the Lab 13-15 editorial 2-column layout and in-depth technical narratives.
"""
import os
import sys

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/22-n-plus-one-query-problem/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/22-n-plus-one-query-problem"
COMMIT_SHA = "6649acb946725e993e78f99ec47833a5498cc870"
LAB_NO = "22"
SLUG = "n-plus-one-query-problem"
TITLE_ID = "N+1 Query Problem"
SUBTITLE_ID = "Batching Relasi Data untuk Eliminasi Latensi Tersembunyi"
TITLE_EN = "N+1 Query Problem"
SUBTITLE_EN = "Batching Database Relations to Eliminate Hidden Latency"
AUTHOR = "Lukman / lukman-ss"
TOPIC_TAGS = "N+1 Query Problem • Batching • Relational Performance • Go"

# Slide data structure
SLIDES = []

def add_slide(categories, title, subtitle, content_items, mental_model, source_path="", slide_num="", narrative="", narrative_en=""):
    """Add a slide to the presentation with support for in-depth editorial narratives."""
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

# =====================================================================
# SLIDE CONTENT DEFINITIONS (24 SLIDES - EDITORIAL LAB 13-15 STANDARD)
# =====================================================================

# Slide 1: Cover (Dark Navy retro terminal / art)
add_slide(
    categories=["COVER"],
    title=TITLE_ID,
    subtitle=SUBTITLE_ID,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "Retro CRT monitor beige dengan scanlines",
            "ASCII art terminal yang menampilkan query count",
            f"Penulis: {AUTHOR}",
        ]}
    ],
    mental_model="",
    source_path="",
    slide_num="01 / 24"
)

# Slide 2: Satu Query Awal, Puluhan Query Mengikuti
add_slide(
    categories=["MOTIVASI", "BUKTI LAB"],
    title="Satu Query Awal, Puluhan Query Mengikuti",
    subtitle="Pengambilan data berulang di dalam iterasi melipatgandakan latensi jaringan.",
    content_items=[
        {"title": "MOTIVASI", "items": [
            "Perilaku cepat di local development menipu engineer saat masuk ke beban production.",
            "Setiap query individu terlihat sub-milidetik, sehingga lolos dari slow query logs.",
            "Degradasi akumuler: request yang seharusnya 20ms membengkak menjadi beberapa detik."
        ], "type": "bullet"},
        {"title": "BUKTI LAB", "items": [
            "Store query counter mencatat 1 query author + 3 query post = 4 query total pada 3 data.",
            "Eager loading memangkas eksekusi menjadi 1 query author + 1 query post = 2 query total."
        ], "type": "bullet"},
        {"title": "BATASAN SISTEM", "items": [
            "Lab memodelkan in-memory Store; pada RDBMS nyata overhead socket TCP memperparah bottleneck."
        ], "type": "bullet"}
    ],
    narrative="N+1 query problem merupakan anti-pattern performa paling umum pada data access layer karena sifatnya yang tidak kasat mata pada volume data kecil. Ketika backend mengambil daftar entitas induk lalu menerbitkan query anak di dalam loop perulangan, jumlah round-trip jaringan meningkat secara linear O(N). Pada tahap development dengan segelintir data, eksekusi selesai sub-milidetik. Namun begitu memasuki lingkungan production dengan ratusan record, akumulasi latensi I/O dan perebutan connection pool akan langsung melumpuhkan throughput layanan.",
    mental_model="N+1 is invisible until it hurts. Batch related queries, don't loop them.",
    source_path="labs/22-n-plus-one-query-problem/README.md",
    slide_num="02 / 24"
)

# Slide 3: Arsitektur Tiga Layer Lab 22
add_slide(
    categories=["ARSITEKTUR", "DATA ACCESS"],
    title="Arsitektur Tiga Layer Lab 22",
    subtitle="Store terisolasi, Repository data access, dan executable demo runner.",
    content_items=[
        {"title": "KOMPONEN ARSITEKTUR", "items": [
            "Komponen | Peran | File Sumber",
            "Store | Mock database thread-safe & query counter | internal/blog/store.go",
            "Repository | Logika query Naive loop vs Eager batching | internal/blog/repository.go",
            "Demo | Eksekusi CLI & verifikasi statistik I/O | cmd/demo/main.go"
        ], "type": "table"}
    ],
    narrative="Arsitektur Lab 22 dibangun dengan pemisahan tanggung jawab yang ketat untuk mengisolasi pengukuran interaksi database. Paket Store bertindak sebagai mock database in-memory yang mencatat setiap operasi baca secara atomik melalui query counter thread-safe. Lapisan Repository mengonsumsi Store untuk menyediakan dua variasi endpoint data access: implementasi naive yang mengeksekusi query di dalam loop, serta implementasi eager yang melakukan batching via filter multi-ID. Pemisahan ini memungkinkan pengujian matematis yang deterministik tanpa dependensi eksternal.",
    mental_model="Clean layer isolation enables precise query instrumentation and verification.",
    source_path="internal/blog/store.go • internal/blog/repository.go",
    slide_num="03 / 24"
)

# Slide 4: Entitas Domain: Author dan Post
add_slide(
    categories=["DOMAIN", "DATA MODEL"],
    title="Entitas Domain: Author dan Post",
    subtitle="Relasi satu-ke-banyak (one-to-many) antara penulis blog dan artikelnya.",
    content_items=[
        {"title": "DOMAIN MODELS", "items": [
            "type Author struct {",
            "    ID   int",
            "    Name string",
            "}",
            "",
            "type Post struct {",
            "    ID       int",
            "    AuthorID int",
            "    Title    string",
            "}",
            "",
            "type AuthorWithPosts struct {",
            "    Author Author",
            "    Posts  []Post",
            "}"
        ], "type": "code"}
    ],
    narrative="Model domain merepresentasikan kasus klasik relational data: satu Author memiliki banyak entitas Post melalui referensi AuthorID. Struktur komposit AuthorWithPosts menjadi agregat yang harus dikembalikan kepada consumer API. Pada relasi one-to-many seperti ini, keputusan arsitektural apakah relasi child diambil secara lazily per item atau secara batched akan menentukan total beban jaringan antara service backend dan database engine.",
    mental_model="Aggregate roots must balance relation completeness with network round-trips.",
    source_path="internal/blog/models.go",
    slide_num="04 / 24"
)

# Slide 5: Loop-Driven Data Retrieval (Anti-Pattern)
add_slide(
    categories=["IMPLEMENTASI", "ANTI-PATTERN"],
    title="Loop-Driven Data Retrieval",
    subtitle="Mengambil data relasi child satu per satu di dalam perulangan entitas parent.",
    content_items=[
        {"title": "NAIVE REPOSITORY", "items": [
            "func (r *Repository) GetAuthorsWithPostsNPlusOne() []AuthorWithPosts {",
            "    authors := r.store.GetAllAuthors()",
            "    if len(authors) == 0 { return []AuthorWithPosts{} }",
            "",
            "    var result []AuthorWithPosts",
            "    for _, author := range authors {",
            "        posts := r.store.GetPostsByAuthorID(author.ID)",
            "        result = append(result, AuthorWithPosts{",
            "            Author: author,",
            "            Posts:  posts,",
            "        })",
            "    }",
            "    return result",
            "}"
        ], "type": "code"}
    ],
    narrative="Implementasi naive pada GetAuthorsWithPostsNPlusOne memperlihatkan bagaimana kode yang tampak bersih dan intuitif dapat menyembunyikan bencana skalabilitas. Pertama, metode memanggil GetAllAuthors yang menghasilkan 1 query. Selanjutnya, perulangan for mengiterasi setiap author dan memanggil GetPostsByAuthorID secara individual. Jika terdapat N author, sistem memicu N query tambahan sehingga total query yang dieksekusi menjadi 1 + N. Pola loop-query ini sangat sering dihasilkan secara otomatis oleh mekanisme lazy loading bawaan ORM.",
    mental_model="Querying inside loops turns iteration into cumulative network latency.",
    source_path="internal/blog/repository.go:13-28",
    slide_num="05 / 24"
)

# Slide 6: Mengapa Slow Query Log Gagal Mendeteksi
add_slide(
    categories=["OBSERVABILITAS", "DATABASE"],
    title="Mengapa Slow Query Log Gagal Mendeteksi",
    subtitle="Setiap query individual sangat cepat, tetapi total durasi request meledak.",
    content_items=[
        {"title": "TITIK BUTA DETEKSI", "items": [
            "Slow query log hanya menandai eksekusi dengan durasi di atas ambang batas (mis. >1000ms).",
            "Query N+1 masing-masing berjalan sub-milidetik (<1ms) karena memakai primary/foreign key index.",
            "Database engine menganggap sistem bekerja normal tanpa anomali query berat."
        ], "type": "bullet"},
        {"title": "AKUMULASI WAKTU", "items": [
            "Development: 3 author memicu 4 query, selesai dalam tempo <5ms.",
            "Production: 500 author memicu 501 query, menghabiskan waktu 2.4 detik.",
            "Beban berpindah ke latency transfer jaringan dan serialisasi paket di aplikasi."
        ], "type": "bullet"}
    ],
    narrative="Mekanisme pemantauan database standar seperti slow query log dirancang untuk menangkap eksekusi query tunggal yang tidak efisien—misalnya akibat full table scan atau hilangnya index. Namun, query pada N+1 biasanya berupa pencarian sederhana berbasis index (WHERE author_id = ?) yang selesai dalam waktu 0.5 milidetik. Akibatnya, database engine mencatat performa normal tanpa ada alert. Beban justru terakumulasi pada latensi perjalanan paket bolak-balik di jaringan dan serialisasi data di application server.",
    mental_model="Sub-millisecond individual queries can aggregate into multi-second latency.",
    source_path="labs/22-n-plus-one-query-problem/content/02-master-draft.md",
    slide_num="06 / 24"
)

# Slide 7: Biaya Tersembunyi: Network Round-Trip Time
add_slide(
    categories=["NETWORK", "PERFORMA"],
    title="Biaya Tersembunyi: Network Round-Trip Time",
    subtitle="Overhead protokol, packet traversal, dan context switching aplikasi.",
    content_items=[
        {"title": "OVERHEAD JARINGAN", "items": [
            "Komponen Biaya | 1 Query Batch | 100 Query N+1",
            "Network RTT | 1 x 2ms = 2ms | 100 x 2ms = 200ms",
            "TCP Wire Overhead | 1 frame exchange | 100 frame exchanges",
            "DB Engine Parsing | 1x AST parse | 100x AST parse",
            "App Context Switch | Minimal | Sangat tinggi (100x I/O)"
        ], "type": "table"}
    ],
    narrative="Ketika aplikasi berkomunikasi dengan database melalui jaringan TCP, setiap query harus menanggung biaya Round-Trip Time (RTT), handshake protocol packet, context switching, serta parsing SQL abstract syntax tree di sisi database. Pada topologi cloud modern di mana database berada di server terpisah, RTT tipikal berkisar antara 1 hingga 5 milidetik. Menjalankan 100 query terpisah berarti membuang minimal 200 hingga 500 milidetik murni hanya untuk waktu transit paket jaringan, terlepas dari seberapa cepat mesin database mengeksekusinya.",
    mental_model="Network round-trips dominate relational latency over raw query execution.",
    source_path="labs/22-n-plus-one-query-problem/research/03-evidence.md",
    slide_num="07 / 24"
)

# Slide 8: Bahaya Connection Pool Exhaustion
add_slide(
    categories=["RESOURCES", "CONCURRENCY"],
    title="Bahaya Connection Pool Exhaustion",
    subtitle="Rentetan query berurutan menahan koneksi dan memicu antrean thread backend.",
    content_items=[
        {"title": "SEQUENTIAL HOLD", "items": [
            "Satu HTTP worker menahan koneksi database selama seluruh siklus iterasi loop N.",
            "Koneksi tidak dapat dipakai oleh worker lain hingga iterasi terakhir selesai."
        ], "type": "bullet"},
        {"title": "POOL STARVATION", "items": [
            "Kenaikan concurrency membuat kapasitas connection pool langsung habis (exhausted).",
            "Request baru mengalami blocking queue di pool dan berakhir dengan gateway timeout (504)."
        ], "type": "bullet"},
        {"title": "CASCADING FAILURE", "items": [
            "Satu endpoint lambat melumpuhkan seluruh service yang berbagi connection pool database."
        ], "type": "bullet"}
    ],
    narrative="Dampak sistemik paling destruktif dari N+1 query problem adalah pengurasan connection pool (connection pool starvation). Dalam arsitektur pooling standar, satu thread web request meminjam koneksi database dan menahannya selama iterasi berlangsung. Jika 50 request konkuren masing-masing membutuhkan 100 query berurutan, pool database yang berkapasitas 20 koneksi akan langsung tersaturasi. Request lain yang membutuhkan query kilat 2 milidetik terpaksa mengantre di queue pool hingga timeout, menyebabkan kegagalan berantai di seluruh layanan.",
    mental_model="Sequential queries monopolize connection pools and trigger cascading timeouts.",
    source_path="labs/22-n-plus-one-query-problem/content/02-master-draft.md",
    slide_num="08 / 24"
)

# Slide 9: Solusi Eager Loading: Pola Batch Data Retrieval
add_slide(
    categories=["SOLUSI", "BATCHING"],
    title="Solusi Eager Loading: Pola Batch Data Retrieval",
    subtitle="Ekstraksi seluruh identifier relasi sebelum mengirim satu query gabungan.",
    content_items=[
        {"title": "ALUR 4 TAHAP EAGER BATCHING", "items": [
            "TAHAP 1: Ambil Parent ── GetAllAuthors() → [ID: 1, 2, 3] (Query #1)",
            "TAHAP 2: Ekstraksi Kunci ── authorIDs = [1, 2, 3] di memori",
            "TAHAP 3: Ambil Child Sekaligus ── GetPostsByAuthorIDs([1, 2, 3]) (Query #2)",
            "TAHAP 4: Asosiasi Memori ── Kelompokkan posts ke author via map[int][]Post"
        ], "type": "diagram"}
    ],
    narrative="Prinsip utama penyelesaian N+1 adalah memisahkan fase pengumpulan kebutuhan data dari fase eksekusi query database. Pendekatan eager batching membaca seluruh parent record terlebih dahulu, mengekstrak seluruh primary key yang relevan ke dalam slice memori, lalu menerbitkan satu query tunggal untuk mengambil seluruh child record yang cocok. Dengan pola ini, terlepas dari apakah terdapat 3, 50, atau 1000 author, jumlah interaksi database dipangkas secara drastis menjadi tepat dua query I/O.",
    mental_model="Batch before you iterate: decouple query planning from data association.",
    source_path="internal/blog/repository.go:30-58",
    slide_num="09 / 24"
)

# Slide 10: Query Batching di Sisi Store
add_slide(
    categories=["STORE", "IN-MEMORY"],
    title="Query Batching di Sisi Store",
    subtitle="Memanfaatkan hash lookup dan filter multi-ID untuk mengambil child records.",
    content_items=[
        {"title": "BATCH QUERY METHOD", "items": [
            "func (s *Store) GetPostsByAuthorIDs(authorIDs []int) []Post {",
            "    s.mu.Lock()",
            "    defer s.mu.Unlock()",
            "    s.queryCount++",
            "",
            "    idMap := make(map[int]bool)",
            "    for _, id := range authorIDs { idMap[id] = true }",
            "",
            "    var result []Post",
            "    for _, p := range s.posts {",
            "        if idMap[p.AuthorID] {",
            "            result = append(result, p)",
            "        }",
            "    }",
            "    return result",
            "}"
        ], "type": "code"}
    ],
    narrative="Di level Store, metode GetPostsByAuthorIDs merepresentasikan eksekusi klausa SQL WHERE author_id IN (...). Implementasi in-memory membangun idMap untuk lookup cepat bertipe O(1) sebelum memfilter koleksi posts. Yang paling krusial: fungsi ini hanya menaikkan counter queryCount sebanyak tepat 1 kali untuk seluruh daftar authorIDs. Di sistem produksi, operasi ini diterjemahkan menjadi satu perintah SQL dengan transmisi payload tunggal, menghemat 99% overhead protokol database.",
    mental_model="A single set-based filter replaces hundreds of point lookups.",
    source_path="internal/blog/store.go:63-80",
    slide_num="10 / 24"
)

# Slide 11: In-Memory Map Grouping pada Layer Aplikasi
add_slide(
    categories=["ALGORITMA", "GO LOGIC"],
    title="In-Memory Map Grouping pada Layer Aplikasi",
    subtitle="Mengelompokkan hasil batch ke dalam hash map untuk asosiasi instan O(1).",
    content_items=[
        {"title": "REKONSTRUKSI HASH MAP", "items": [
            "// 1. Kelompokkan posts berdasarkan author ID",
            "postsByAuthor := make(map[int][]Post)",
            "for _, post := range allPosts {",
            "    postsByAuthor[post.AuthorID] = append(",
            "        postsByAuthor[post.AuthorID], post)",
            "}",
            "",
            "// 2. Susun agregat akhir dengan pencarian O(1)",
            "var result []AuthorWithPosts",
            "for _, author := range authors {",
            "    result = append(result, AuthorWithPosts{",
            "        Author: author,",
            "        Posts:  postsByAuthor[author.ID],",
            "    })",
            "}"
        ], "type": "code"}
    ],
    narrative="Tantangan setelah mengambil sekumpulan child records adalah merekonstruksi relasi hierarkis ke entitas parent aslinya. Solusi optimal di level aplikasi adalah mengelompokkan hasil batch ke dalam tabel hash lokal map[int][]Post dengan key berupa foreign key AuthorID. Iterasi rekonstruksi kemudian hanya membutuhkan pencarian hash O(1) untuk menyematkan slice posts ke masing-masing author. Menghindari nested loop di memori menjamin kompleksitas waktu pemrosesan aplikasi tetap linear O(N + M).",
    mental_model="Group child entities in memory to preserve O(1) relational association.",
    source_path="internal/blog/repository.go:45-56",
    slide_num="11 / 24"
)

# Slide 12: Komparasi Kuantitatif Kompleksitas Sistem
add_slide(
    categories=["KOMPARASI", "BENCHMARK"],
    title="Komparasi Kuantitatif Kompleksitas Sistem",
    subtitle="Transformasi kompleksitas I/O dari ketergantungan data O(N) ke deterministik O(1).",
    content_items=[
        {"title": "MATRIKS KOMPARASI", "items": [
            "Parameter | Naive (Loop) | Eager (Batch)",
            "Query Count | 1 + N (Linear) | 1 + 1 (Konstan 2)",
            "Kompleksitas I/O | O(N) relatif record | O(1) relatif record",
            "Data Transfer | Header repetitif | Single packet payload",
            "Memory Allocation | Fragmented allocations | Single pre-sized buffer",
            "Network RTT | N + 1 round-trips | 2 round-trips"
        ], "type": "table"}
    ],
    narrative="Tabel komparasi menggarisbawahi perubahan fundamental pada profil operasional sistem. Dalam pendekatan loop-driven, kompleksitas I/O bersifat O(N), yang berarti setiap penambahan data author di database secara langsung meningkatkan beban I/O jaringan. Sebaliknya, pendekatan eager loading menurunkan kompleksitas interaksi database menjadi O(1) relatif terhadap jumlah baris data—hanya bertambah jika terdapat entitas relasi tambahan (O(k) di mana k adalah jumlah entitas relasional). Ini menghasilkan stabilitas throughput yang terprediksi.",
    mental_model="Shift complexity from database network I/O into deterministic in-memory grouping.",
    source_path="labs/22-n-plus-one-query-problem/content/02-master-draft.md",
    slide_num="12 / 24"
)

# Slide 13: SQL Translation: IN Clause vs Single JOIN
add_slide(
    categories=["SQL", "RDBMS"],
    title="SQL Translation: IN Clause vs Single JOIN",
    subtitle="Membandingkan strategi batch query terpisah versus query tunggal dengan JOIN.",
    content_items=[
        {"title": "DUA PENDEKATAN SQL", "items": [
            "-- Strategi 1: Batching (Dua Query Terpisah)",
            "SELECT id, name FROM authors;",
            "SELECT id, author_id, title FROM posts",
            "WHERE author_id IN (1, 2, 3);",
            "",
            "-- Strategi 2: Single JOIN (Satu Query Tunggal)",
            "SELECT a.id, a.name, p.id, p.title",
            "FROM authors a",
            "LEFT JOIN posts p ON a.id = p.author_id;"
        ], "type": "code"}
    ],
    narrative="Dalam dunia RDBMS, terdapat dua pendekatan utama untuk eager loading: dua query terpisah menggunakan klausa IN (split queries), atau satu query tunggal menggunakan LEFT JOIN. Meskipun single JOIN terdengar menggoda karena hanya membutuhkan 1 query, teknik ini menduplikasi data kolom parent pada setiap baris hasil child yang dikembalikan. Pada relasi one-to-many dengan kolom teks besar, duplikasi ini memboroskan bandwidth jaringan. Split queries via IN clause umumnya menjadi pilihan default ORM modern karena lebih efisien terhadap memori.",
    mental_model="Split batch queries avoid parent column duplication inherent in wide JOINs.",
    source_path="labs/22-n-plus-one-query-problem/research/03-evidence.md",
    slide_num="13 / 24"
)

# Slide 14: Batasan Parameter IN Clause pada RDBMS
add_slide(
    categories=["BATASAN", "PRODUCTION"],
    title="Batasan Parameter IN Clause pada RDBMS",
    subtitle="Batas maksimum parameter dan overhead parsing abstract syntax tree database.",
    content_items=[
        {"title": "BATASAN PARAMETER", "items": [
            "PostgreSQL & SQLite memiliki limit bind parameter (mis. 65,535 atau 999 parameter).",
            "Melebihi batas ini menyebabkan error query execution langsung dari driver database."
        ], "type": "bullet"},
        {"title": "OVERHEAD PARSER", "items": [
            "Statement SQL dengan ribuan literal ID membebani query planner dan parser CPU.",
            "Waktu penyusunan AST di engine database dapat melebihi waktu eksekusi indeks."
        ], "type": "bullet"},
        {"title": "CHUNKED BATCHING", "items": [
            "Pecah daftar ID besar menjadi chunk terkendali (misal 200–500 ID per batch query)."
        ], "type": "bullet"}
    ],
    narrative="Batching menggunakan klausa IN bukan tanpa batas. Database engine memiliki batasan terhadap jumlah parameter bind dalam satu statement (misalnya batas 65.535 parameter pada PostgreSQL atau 999 pada SQLite versi terdahulu). Lebih jauh lagi, statement SQL dengan ribuan literal parameter memaksa query planner menghabiskan waktu CPU berharga hanya untuk mem-parsing sintaks string query. Oleh karena itu, pada skala data masif, aplikasi wajib menerapkan teknik chunked batching—memecah daftar ID menjadi batch terkendali (misal 500 ID per batch).",
    mental_model="Unbounded IN clauses choke SQL parsers: chunk bulk IDs into manageable sizes.",
    source_path="labs/22-n-plus-one-query-problem/content/02-master-draft.md",
    slide_num="14 / 24"
)

# Slide 15: Bahaya Cartesian Product Explosion
add_slide(
    categories=["ARSITEKTUR", "RISIKO"],
    title="Bahaya Cartesian Product Explosion",
    subtitle="Risiko ekspansi baris berlipat ganda saat melakukan eager JOIN multi-relasi.",
    content_items=[
        {"title": "MULTIPLIKASI BARIS", "items": [
            "Tingkat Relasi | Entitas Record | Total Baris Hasil JOIN",
            "1 Author | 1 author | 1 baris",
            "10 Posts | 10 post per author | 1 x 10 = 10 baris",
            "5 Tags | 5 tag per post | 1 x 10 x 5 = 50 baris",
            "10 Comments | 10 komentar per post | 1 x 10 x 5 x 10 = 500 baris!"
        ], "type": "table"}
    ],
    narrative="Ketika pengembang mencoba mengoptimalkan beberapa relasi one-to-many sekaligus (misalnya mengambil Author beserta Posts dan Tags-nya) menggunakan JOIN tunggal, sistem akan terkena fenomena Cartesian Product Explosion. Database mengalikan setiap entitas parent dengan seluruh permutasi child records-nya. Data yang dikembalikan melalui jaringan membengkak ribuan kali lipat dan berpotensi memicu Out-Of-Memory (OOM) pada backend. Solusi arsitektural yang aman adalah membagi relasi menjadi batch queries independen terpisah.",
    mental_model="Multiple relational JOINs cause combinatorial row explosions; use separate batches.",
    source_path="labs/22-n-plus-one-query-problem/content/02-master-draft.md",
    slide_num="15 / 24"
)

# Slide 16: Mapping-Level Eager vs Query-Level Eager
add_slide(
    categories=["ORM", "ANTI-PATTERN"],
    title="Mapping-Level Eager vs Query-Level Eager",
    subtitle="Bahaya anotasi global FetchType.EAGER dibandingkan selective fetching.",
    content_items=[
        {"title": "MAPPING-LEVEL (ANTI-PATTERN)", "items": [
            "FetchType.EAGER pada schema entity menarik seluruh relasi setiap kali entitas dibaca.",
            "Operasi sederhana seperti validasi email author terpaksa memuat seluruh riwayat postingan.",
            "Menimbulkan memory bloat parah dan membebani garbage collector aplikasi."
        ], "type": "bullet"},
        {"title": "QUERY-LEVEL (BEST PRACTICE)", "items": [
            "Relasi dibiarkan lazy secara default pada level definisi domain model.",
            "Eager loading hanya diaktifkan secara eksplisit pada query yang benar-benar membutuhkannya."
        ], "type": "bullet"}
    ],
    narrative="Kesalahan umum yang sering terjadi di level framework (khususnya JPA/Hibernate) adalah menerapkan mapping-level eager loading (FetchType.EAGER) secara global pada entity schema. Hal ini menyebabkan setiap pembacaan entitas Author otomatis menarik seluruh Posts, Comments, dan Metadata terkait di mana pun dipanggil. Pendekatan profesional menuntut query-level eager loading: biarkan relasi default terisolasi, dan aktifkan batching hanya secara eksplisit pada lapisan repository atau service yang benar-benar membutuhkan agregat data tersebut.",
    mental_model="Enforce query-level eager fetching; avoid global schema-level eager mappings.",
    source_path="labs/22-n-plus-one-query-problem/content/02-master-draft.md",
    slide_num="16 / 24"
)

# Slide 17: Solusi Ekosistem: GraphQL DataLoader Pattern
add_slide(
    categories=["EKOSISTEM", "GRAPHQL"],
    title="Solusi Ekosistem: GraphQL DataLoader Pattern",
    subtitle="Mengotomatisasi penundaan dan batching query pada resolver independen.",
    content_items=[
        {"title": "ALUR RESOLVER COALESCING", "items": [
            "Resolver #1 (Author 1) ──┐",
            "Resolver #2 (Author 2) ──┼──> [DataLoader Queue] ──> Single Batch Query",
            "Resolver #3 (Author 3) ──┘        (Tick Event)        WHERE author_id IN (1, 2, 3)"
        ], "type": "diagram"}
    ],
    narrative="Pada arsitektur GraphQL, masalah N+1 muncul secara alami karena setiap field resolver dieksekusi secara independen. Filosofi batching diwujudkan melalui pola DataLoader. DataLoader bekerja dengan menampung (coalescing) permintaan fetch ID individual yang terjadi dalam satu event loop tick, menunda eksekusi ke database, lalu mengirimkan satu query batch terkonsolidasi. DataLoader membuktikan bahwa prinsip pemisahan pengumpulan kunci dari eksekusi I/O dapat diterapkan secara elegan pada lapisan middleware.",
    mental_model="DataLoaders coalesce concurrent point queries across async resolver ticks.",
    source_path="labs/22-n-plus-one-query-problem/research/03-evidence.md",
    slide_num="17 / 24"
)

# Slide 18: Column Selection & Aggregation Optimization
add_slide(
    categories=["OPTIMASI", "SQL"],
    title="Column Selection & Aggregation Optimization",
    subtitle="Memilih kolom spesifik atau kalkulasi database saat entitas penuh tidak diperlukan.",
    content_items=[
        {"title": "AGGREGATION VS FULL ENTITY", "items": [
            "-- 1. Ambil ringkasan agregasi alih-alih seluruh baris post",
            "SELECT author_id, COUNT(*) AS post_count",
            "FROM posts WHERE author_id IN (1, 2, 3)",
            "GROUP BY author_id;",
            "",
            "-- 2. Proyeksikan kolom esensial saja",
            "SELECT id, author_id, title",
            "FROM posts WHERE author_id IN (1, 2, 3);"
        ], "type": "code"}
    ],
    narrative="Eager loading penuh tidak selalu menjadi solusi terbaik jika kebutuhan bisnis hanyalah informasi ringkasan atau segelintir kolom. Jika use-case hanya ingin menampilkan jumlah artikel tiap penulis, memuat seluruh objek Post ke dalam memori aplikasi adalah pemborosan besar. Eksekusi agregasi langsung di database melalui GROUP BY atau subquery scalar menghasilkan efisiensi transfer data yang jauh lebih superior, menjaga footprint RAM aplikasi tetap ramping.",
    mental_model="Push aggregation to the database engine before pulling bulk entities into memory.",
    source_path="labs/22-n-plus-one-query-problem/content/02-master-draft.md",
    slide_num="18 / 24"
)

# Slide 19: Demo CLI Execution Output
add_slide(
    categories=["DEMO", "VERIFIKASI"],
    title="Demo CLI Execution Output",
    subtitle="Bukti eksekusi nyata dari cmd/demo/main.go pada dataset 3 author.",
    content_items=[
        {"title": "OUTPUT TERMINAL", "items": [
            "$ go run ./cmd/demo",
            "--- 1. Simulating N+1 Query Problem ---",
            "Loaded 3 authors with their posts.",
            "Total queries executed: 4 (1 for authors + 3 queries for posts)",
            "",
            "--- 2. Simulating Eager Loading (Batching) ---",
            "Loaded 3 authors with their posts.",
            "Total queries executed: 2 (1 for authors + 1 batched query for posts)"
        ], "type": "code"}
    ],
    narrative="Eksekusi demo pada cmd/demo/main.go mendemonstrasikan kontras performa secara langsung di terminal console. Pada dataset uji dengan 3 author (Alice, Bob, Charlie) dan 5 total post, implementasi naive mencatat tepat 4 query database. Sebaliknya, implementasi eager loading menyelesaikan seluruh pemuatan agregat dengan tepat 2 query. Verifikasi langsung ini menegaskan validitas teori tanpa bergantung pada konfigurasi database pihak ketiga yang rumit.",
    mental_model="Empirical CLI verification validates algorithmic query reduction.",
    source_path="labs/22-n-plus-one-query-problem/cmd/demo/main.go",
    slide_num="19 / 24"
)

# Slide 20: Pengujian Otomatis Query Count Invariant
add_slide(
    categories=["TESTING", "CI/CD"],
    title="Pengujian Otomatis Query Count Invariant",
    subtitle="Memverifikasi jumlah query secara deterministik dalam unit test CI/CD.",
    content_items=[
        {"title": "TEST QUERY COUNT", "items": [
            "func TestGetAuthorsWithPostsNPlusOne(t *testing.T) {",
            "    store := NewStore()",
            "    repo := NewRepository(store)",
            "    store.ResetQueryCount()",
            "",
            "    result := repo.GetAuthorsWithPostsNPlusOne()",
            "    if len(result) != 3 { t.Fatalf(\"expected 3, got %d\", len(result)) }",
            "",
            "    queryCount := store.GetQueryCount()",
            "    if queryCount != 4 {",
            "        t.Errorf(\"expected 4 queries (N+1), got %d\", queryCount)",
            "    }",
            "}"
        ], "type": "code"}
    ],
    narrative="Salah satu tantangan terbesar N+1 adalah kerentanannya terhadap regresi kode di masa depan: seorang engineer baru dapat secara tidak sengaja menambahkan pemanggilan accessor di dalam template atau loop. Lab 22 menunjukkan bagaimana menulis unit test yang memverifikasi query count secara deterministik. Dengan menginstruksikan Store mock untuk menghitung pemanggilan query, pipeline CI/CD dapat secara otomatis menggagalkan build jika jumlah interaksi database melebihi ambang batas invarian.",
    mental_model="Lock query count invariants into CI tests to permanently prevent regressions.",
    source_path="internal/blog/repository_test.go:8-25",
    slide_num="20 / 24"
)

# Slide 21: Validasi Kesetaraan Semantik (Deep Equivalence)
add_slide(
    categories=["INTEGRITAS", "ASSERTION"],
    title="Validasi Kesetaraan Semantik (Deep Equivalence)",
    subtitle="Memastikan optimasi query tidak merusak struktur atau kebenaran data bisnis.",
    content_items=[
        {"title": "DEEP EQUIVALENCE TEST", "items": [
            "// Memastikan hasil batching 100% identik dengan naive",
            "nPlusOneResult := repo.GetAuthorsWithPostsNPlusOne()",
            "if !reflect.DeepEqual(result, nPlusOneResult) {",
            "    t.Errorf(\"eager result mismatch with N+1 result:\\n\"+",
            "        \"Eager: %+v\\nN+1: %+v\", result, nPlusOneResult)",
            "}"
        ], "type": "code"}
    ],
    narrative="Optimasi performa tidak ada artinya jika mengubah integritas atau susunan data yang diterima oleh pemanggil. Dalam TestGetAuthorsWithPostsEager, paket menguji bahwa slice AuthorWithPosts yang dihasilkan oleh algoritma batching identik secara mendalam (deeply equal) dengan hasil algoritma naive menggunakan reflect.DeepEqual. Pengujian ini membuktikan bahwa pengelompokan memori berbasis hash map mempertahankan struktur data bisnis secara sempurna tanpa efek samping.",
    mental_model="Performance refactoring must guarantee semantic data equivalence.",
    source_path="internal/blog/repository_test.go:45-50",
    slide_num="21 / 24"
)

# Slide 22: Observability & Deteksi N+1 di Production
add_slide(
    categories=["OBSERVABILITAS", "TRACING"],
    title="Observability & Deteksi N+1 di Production",
    subtitle="Instrumentasi OpenTelemetry spans, query profiling, dan static linters.",
    content_items=[
        {"title": "DISTRIBUTED TRACING", "items": [
            "APM trace menampilkan 'waterfall graph' berisi puluhan span query identik berurutan.",
            "Waktu tunggu jaringan antar span memperjelas latensi akumulatif yang dialami client."
        ], "type": "bullet"},
        {"title": "MIDDLEWARE AUDITING", "items": [
            "Pasang query logger di database driver untuk mendeteksi repetisi query template per HTTP request.",
            "Beri threshold alarm jika satu request mengeksekusi >20 query serupa dalam durasi <50ms."
        ], "type": "bullet"},
        {"title": "STATIC ANALYSIS", "items": [
            "Gunakan linter kode untuk menandai eksekusi database calls di dalam perulangan loop for."
        ], "type": "bullet"}
    ],
    narrative="Di lingkungan produksi skala besar, deteksi N+1 dilakukan melalui instrumentasi Application Performance Monitoring (APM) dan OpenTelemetry distributed tracing. Trace visual memperlihatkan serangkaian span database berurutan dengan query template yang sama persis (waterfall trace). Melengkapi tracing dengan query counter middleware dan static code analysis linter memastikan bahwa setiap anomali rasio query terhadap HTTP request langsung memicu notifikasi peringatan sebelum mengganggu ketersediaan sistem.",
    mental_model="Waterfall APM traces and query-to-request ratios surface hidden runtime loops.",
    source_path="labs/22-n-plus-one-query-problem/content/02-master-draft.md",
    slide_num="22 / 24"
)

# Slide 23: Reproduksi Mandiri & Panduan Verifikasi
add_slide(
    categories=["PANDUAN", "REPRODUKSI"],
    title="Reproduksi Mandiri & Panduan Verifikasi",
    subtitle="Seluruh implementasi berbasis Go standard library siap dijalankan dan diuji.",
    content_items=[
        {"title": "COMMAND RUNNER", "items": [
            "cd labs/22-n-plus-one-query-problem",
            "go test -v ./...",
            "go test -race ./...",
            "go run ./cmd/demo"
        ], "type": "code"}
    ],
    narrative="Modul Lab 22 ditulis menggunakan Go 1.22 standar tanpa framework eksternal, memastikan portabilitas tinggi dan reproduksi instan di mesin apa pun. Perintah go test menjalankan seluruh suite pengujian query count dan edge cases tanpa data race. Perintah demo menjalankan perbandingan langsung antara kedua strategi. Semua artefak pengujian dan catatan audit terdokumentasi lengkap di direktori lab untuk studi arsitektur lebih lanjut.",
    mental_model="Standard library benchmarks provide reliable, zero-dependency validation.",
    source_path="labs/22-n-plus-one-query-problem/README.md • go.mod",
    slide_num="23 / 24"
)

# Slide 24: Closing Slide (Dark Navy `#0b1b2d`)
add_slide(
    categories=["CLOSING"],
    title="Eksperimen kecil, eliminasi latensi tersembunyi.",
    subtitle="",
    content_items=[
        {"title": "MODUL LAB", "items": [
            "internal/blog/models.go • internal/blog/store.go",
            "internal/blog/repository.go • repository_test.go",
            "cmd/demo/main.go • README.md"
        ], "type": "bullet"},
        {"title": "VERIFIKASI & STATUS", "items": [
            f"Sumber: commit {COMMIT_SHA[:8]}",
            "Test passing dan diverifikasi (Race-Clean Go 1.22).",
            "Invarian terbukti: 4 query (N+1) dipangkas menjadi 2 query (Batching)."
        ], "type": "bullet"}
    ],
    mental_model="N+1 is invisible until it hurts. Batch related queries, don't loop them.",
    source_path=REPO_URL,
    slide_num="24 / 24"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
