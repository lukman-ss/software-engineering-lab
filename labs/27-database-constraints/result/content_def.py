#!/usr/bin/env python3
"""
Publication Builder Definition for Software Engineering Lab 27: Database Constraints.
Generates 24-slide LinkedIn carousel covering declarative database constraints,
race condition prevention under concurrency, SQLSTATE error taxonomy, partial unique indexes,
and production data integrity trade-offs.
"""

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/27-database-constraints/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/27-database-constraints"
COMMIT_SHA = "c317c1d"
LAB_NO = "27"
SLUG = "database-constraints"

TITLE_ID = "Database Constraints: Benteng Terakhir Integritas Data"
SUBTITLE_ID = "Mencegah Race Condition, Menjaga Invarian Skema, dan Menolak Duplikasi Atomik di Go"
TITLE_EN = "Database Constraints: The Last Line of Data Integrity"
SUBTITLE_EN = "Preventing Race Conditions, Enforcing Schema Invariants, and Rejecting Duplicates Atomically in Go"

AUTHOR = "Lukman / lukman-ss"
TOPIC_TAGS = "Database Constraints • Data Integrity • Concurrency • SQLSTATE • Partial Index • Go"
COVER_IMAGE_PATH = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/27-database-constraints/result/image-cover/cover_art.jpg"

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
            "Retro Terminal Visualizing Schema Invariants, SQLSTATE Taxonomy & Atomic Locks",
            "Simulasi Relasional In-Memory dengan 50 Goroutines Konkuren dan Error Mapping di Go",
            f"Penulis: {AUTHOR}"
        ]}
    ],
    mental_model="Database constraints provide declarative, ACID-guaranteed enforcement that eliminates check-then-act vulnerabilities.",
    source_path=REPO_URL,
    slide_num="01 / 24",
    narrative="Banyak sistem backend mempercayakan integritas data sepenuhnya pada kode aplikasi. Ketika beban konkuren melonjak, pola validasi naif membiarkan data duplikat dan relasi yatim piatu lolos ke penyimpanan. Lab ini membuktikan secara empiris bagaimana database constraints bertindak sebagai benteng pertahanan terakhir yang atomik dan tak tertembus.",
    narrative_en="Many backend architectures delegate data integrity entirely to application code. Under concurrent load, naive validation lets duplicate records and orphaned relations slip into persistent storage. This lab empirically demonstrates how database constraints act as the unbreachable, atomic last line of defense.",
    layout="cover"
)

# Slide 02: Problem Nyata (Layout: cards)
add_slide(
    categories=["PROBLEM", "DATA INTEGRITY HAZARD"],
    title="Validasi Aplikasi Runtuh di Bawah Beban Konkuren",
    subtitle="Pola check-then-act di memori aplikasi membuka jendela tabrakan data duplikat.",
    content_items=[
        {"type": "cards", "title": "DAMPAK KEGAGALAN VALIDASI HANYA DI APLIKASI", "items": [
            {"title": "Race Condition Registrasi", "desc": "Dua request membaca 'email belum ada' secara bersamaan, lalu keduanya melakukan INSERT email identik."},
            {"title": "Deduplikasi Pasca-Fakta Mahal", "desc": "Membersihkan data duplikat di produksi membutuhkan rekonsiliasi manual, migration rumit, dan downtime."},
            {"title": "Bukan Berarti Bug Kode", "desc": "Kegagalan ini bukan berarti logika if/else salah, melainkan batasan fundamental arsitektur tanpa penguncian atomik."}
        ]}
    ],
    narrative="Pola check-then-act yang memeriksa ketersediaan data via SELECT sebelum melakukan INSERT selalu tampak berhasil pada pengujian sekuensial. Namun di bawah beban konkuren, dua thread membaca kondisi yang sama sebelum salah satunya menulis. Kegagalan ini bukan berarti kode if-check rusak, melainkan akibat ketiadaan isolasi atomik pada titik penyimpanan.",
    narrative_en="The check-then-act pattern inspecting records via SELECT before issuing an INSERT invariably succeeds in sequential testing. Under concurrent traffic, however, two threads observe the same pre-write state before either commits. This failure does not stem from faulty conditional logic, but from lacking atomic isolation at the persistence tier.",
    mental_model="Application-level checks cannot bridge the race window between read and write operations under concurrency.",
    source_path="content/02-master-draft.md:3-24",
    slide_num="02 / 24",
    layout="cards"
)

# Slide 03: Mengapa Review Manual / Pendekatan Naif Gagal (Layout: cards)
add_slide(
    categories=["WHY FAILS", "CONCURRENCY ANOMALY"],
    title="Mengapa Review Manual dan Pengecekan Memori Gagal?",
    subtitle="Pemeriksaan memori tanpa proteksi engine database tidak mampu menahan beban paralel.",
    content_items=[
        {"type": "cards", "title": "KELEMAHAN PENDEKATAN VALIDASI NAIF", "items": [
            {"title": "Interleaving CPU Scheduling", "desc": "Jeda milidetik antara query SELECT dan INSERT dieksploitasi oleh OS scheduler untuk menyisipkan operasi lain."},
            {"title": "Unit Test Sekuensial Menipu", "desc": "Test standar tanpa konkurensi goroutine selalu pass 100%, menciptakan ilusi keamanan palsu bagi tim engineer."},
            {"title": "Tidak Otomatis Terlindungi", "desc": "Menjalankan validasi di microservice terpisah tidak otomatis menjamin keunikan tanpa koordinasi locking atomik."}
        ]}
    ],
    narrative="Review manual sering melewatkan celah integritas karena kode validasi tampak sangat rapi dan logis di editor. Unit test sekuensial selalu memberikan status hijau karena tidak ada kompetisi thread. Menjalankan validasi di lapisan aplikasi tidak otomatis melindungi sistem jika engine penyimpanan tidak menegakkan kunci eksklusif pada saat penulisan.",
    narrative_en="Code reviews routinely miss integrity vulnerabilities because application-level checks appear impeccably structured on screen. Sequential unit tests pass reliably, fostering false confidence across engineering teams. Running validations purely in application services does not automatically ensure integrity without atomic database constraints.",
    mental_model="Sequential testing creates an illusion of correctness that shatters under concurrent thread interleaving.",
    source_path="content/02-master-draft.md:15-37",
    slide_num="03 / 24",
    layout="cards"
)

# Slide 04: Mental Model Inti — Schema Card (Layout: schema)
add_slide(
    categories=["MENTAL MODEL", "DEFENSE IN DEPTH"],
    title="Model Mental Dua Lapis: Kolaborasi Aplikasi dan Database",
    subtitle="Memisahkan peran kenyamanan pengguna (UX) dari jaminan kebenaran data mutlak.",
    content_items=[
        {"type": "schema", "title": "ARSITEKTUR PERTAHANAN DUA LAPIS", "items": [
            {"phase": "Lapis 1: Application Validation", "role": "User Experience (UX)", "desc": "Memberikan pesan error ramah, validasi format cepat, dan sanitasi input sebelum request menyentuh database."},
            {"phase": "Lapis 2: Database Constraints", "role": "Last Line of Truth", "desc": "Menjamin invarian skema secara deklaratif dan atomik dengan lock internal engine saat baris data ditulis."},
            {"phase": "Jembatan: Error Mapping", "role": "SQLSTATE Taxonomy", "desc": "Hanya pada jalur ini error database 23xxx diterjemahkan menjadi domain error terstruktur seperti HTTP 409 Conflict."}
        ]}
    ],
    narrative="Arsitektur modern membagi tanggung jawab menjadi dua lapis: aplikasi melayani keramahan pesan pengguna, sedangkan database menjaga kebenaran data absolut. Validasi aplikasi menyaring input yang salah format dengan cepat. Hanya pada jalur ini, ketika terjadi tabrakan konkuren, database constraint menolak pelanggaran secara atomik dan mengembalikan kode SQLSTATE terstandar.",
    narrative_en="Modern software architectures divide integrity responsibilities into two layers: the application handles user experience, while the database guarantees absolute correctness. Application validation filters invalid formats swiftly. Only along this path, when concurrent collisions occur, does the database reject violations atomically with standard SQLSTATE codes.",
    mental_model="Application validation provides user-friendly guidance; database constraints provide inviolable invariant enforcement.",
    source_path="content/02-master-draft.md:25-37",
    slide_num="04 / 24",
    layout="schema"
)

# Slide 05: Diagram Alur / Pipeline (Layout: schema)
add_slide(
    categories=["PIPELINE", "EVALUATION ORDER"],
    title="Urutan Evaluasi Constraint: Alur Validasi Internal Mesin",
    subtitle="Bagaimana storage engine memproses aturan skema secara berurutan sebelum commit baris.",
    content_items=[
        {"type": "schema", "title": "TAHAPAN EVALUASI CONSTRAINT PADA OPERASI INSERT", "items": [
            {"phase": "1. NOT NULL (23502)", "role": "Column Scope", "desc": "Memeriksa kekosongan nilai field wajib; gagal seketika jika kolom mandatori kosong."},
            {"phase": "2. CHECK (23514)", "role": "Predicate Scope", "desc": "Mengevaluasi ekspresi boolean pada baris (misalnya age >= 18 atau status valid)."},
            {"phase": "3. UNIQUE / PARTIAL (23505)", "role": "Index Scope", "desc": "Memeriksa tabrakan nilai pada B-tree index aktif di bawah proteksi lock eksklusif."},
            {"phase": "4. FOREIGN KEY (23503)", "role": "Relation Scope", "desc": "Memverifikasi keberadaan ID induk pada tabel referensi sebelum commit data."},
            {"phase": "5. ATOMIC COMMIT", "role": "Storage Commit", "desc": "Baris data dan entri index ditulis bersamaan; belum membuktikan engine multi-node di lab ini."}
        ]}
    ],
    narrative="Storage engine mengevaluasi batasan data dalam urutan deterministik: NOT NULL diperiksa lebih dulu, disusul ekspresi boolean CHECK, lalu penelusuran UNIQUE index, dan terakhir FOREIGN KEY. Jika salah satu tahap gagal, operasi dibatalkan seketika tanpa ada baris data atau indeks yang tertulis. Lab ini memodelkan alur sekuensial tersebut secara atomik di dalam satu critical section.",
    narrative_en="The storage engine evaluates schema constraints in a deterministic sequence: NOT NULL runs first, followed by CHECK boolean expressions, UNIQUE index scans, and FOREIGN KEY verification. If any check fails, execution aborts immediately without mutating tables or indexes. This lab models this sequential evaluation atomically within a unified critical section.",
    mental_model="Deterministic constraint pipelines reject violations early, guaranteeing zero partial writes upon schema conflicts.",
    source_path="content/04-diagrams.md:69-95",
    slide_num="05 / 24",
    layout="schema"
)

# Slide 06: Cuplikan Kode Dominan (Layout: code)
add_slide(
    categories=["IMPLEMENTATION", "STORAGE ENGINE"],
    title="Evaluasi Atomik Constraint di Bawah Mutex Lock",
    subtitle="Memastikan tidak ada jeda antara pemeriksaan aturan dan penulisan baris data.",
    content_items=[
        {"type": "code", "title": "internal/engine/engine.go", "items": [
            "func (e *Engine) InsertUser(u model.User, usePartial bool) (model.User, error) {",
            "    e.mu.Lock()",
            "    defer e.mu.Unlock()",
            "    if u.Email == \"\" { return model.User{}, dberr.NewNotNullViolation(\"users\", \"email\", \"users_email_not_null\") }",
            "    if u.Age < 18 { return model.User{}, dberr.NewCheckViolation(\"users\", \"users_age_check\", \"age >= 18\") }",
            "    if exists := e.emailIndex[u.Email]; exists {",
            "        return model.User{}, dberr.NewUniqueViolation(\"users\", \"users_email_key\", \"duplicate email\")",
            "    }",
            "    u.ID = e.userSeq.Add(1); e.users[u.ID] = u; e.emailIndex[u.Email] = u.ID",
            "    return u, nil",
            "}"
        ]}
    ],
    narrative="Metode InsertUser menjalankan seluruh pengecekan di dalam satu blok mutex lock yang merepresentasikan ROW EXCLUSIVE lock pada PostgreSQL. Mengunci seluruh proses bukan berarti implementasi ini memakai latch halaman B-tree yang kompleks, melainkan memodelkan jaminan atomisitas agar tidak ada goroutine lain yang dapat menyisipkan email serupa.",
    narrative_en="The InsertUser method executes all constraint checks inside an exclusive mutex lock representing PostgreSQL's ROW EXCLUSIVE lock. Locking the critical section does not mean this lab implements low-level B-tree page latches; it cleanly models atomic guarantees so competing goroutines cannot insert conflicting emails.",
    mental_model="Encapsulating validation and mutation within an atomic critical section eliminates check-then-act vulnerabilities.",
    source_path="internal/engine/engine.go:45-99",
    slide_num="06 / 24",
    layout="code"
)

# Slide 07: Failure Path / Dangling Reference Konkret (Layout: cards)
add_slide(
    categories=["FAILURE PATH", "RACE CONDITION"],
    title="Kronologi Tabrakan: Dua Goroutine pada Validasi Naif",
    subtitle="Langkah demi langkah runtuhnya integritas data pada UnsafeStore tanpa constraint.",
    content_items=[
        {"type": "cards", "title": "FASE TERJADINYA DUPLIKASI DATA SECARA KONKUREN", "items": [
            {"title": "1. Pengecekan Bersamaan", "desc": "Goroutine A dan B sama-sama memindai tabel dan membaca bahwa email target belum terdaftar."},
            {"title": "2. Jeda Penjadwalan CPU", "desc": "Keduanya memasuki jeda scheduling (time.Sleep) yang membuka race window di memori aplikasi."},
            {"title": "3. Penulisan Ganda Sukses", "desc": "Goroutine A menulis row; Goroutine B langsung menyusul menulis row dengan email yang sama persis."},
            {"title": "4. Integritas Runtuh", "desc": "Hanya pada jalur ini tabel memuat 2 baris dengan email sama, merusak invarian keunikan sistem."}
        ]}
    ],
    narrative="Ketika dua goroutine mengeksekusi UnsafeStore.RegisterUser secara paralel, keduanya sama-sama menemukan bahwa email belum terdaftar di database. Jeda scheduling CPU mengekspos race window sehingga kedua proses sama-sama melanjutkan ke langkah INSERT. Hanya pada jalur ini duplikasi data berhasil lolos ke penyimpanan, membuktikan bahwa pengecekan memori tanpa constraint tidak berdaya.",
    narrative_en="When two goroutines execute UnsafeStore.RegisterUser concurrently, both observe that the requested email does not yet exist. CPU scheduling delays expose a vulnerability window, prompting both threads to issue INSERT calls. Only along this path do duplicate records infiltrate storage, proving memory checks without constraints are powerless.",
    mental_model="Uncoordinated concurrent reads before unconstrained writes guarantee permanent data corruption.",
    source_path="content/04-diagrams.md:19-40",
    slide_num="07 / 24",
    layout="cards"
)

# Slide 08: Tabel Perbandingan Mekanisme vs Gejala (Layout: table)
add_slide(
    categories=["COMPARISON", "ERROR TAXONOMY"],
    title="Matriks Taksonomi: Jenis Constraint vs Karakteristik",
    subtitle="Membedakan cakupan validasi, kode SQLSTATE standar, dan mitigasi arsitektur.",
    content_items=[
        {"title": "PERBANDINGAN MEKANISME CONSTRAINT DATABASE", "items": [
            "Jenis Constraint | SQLSTATE | Cakupan Evaluasi | Pencegahan Kegagalan | HTTP Domain Mapping",
            "NOT NULL | 23502 | Kolom Tunggal | Mencegah record kosong tak beridentitas | 400 Bad Request",
            "CHECK | 23514 | Kolom / Baris | Menolak nilai di luar batas bisnis logis | 422 Unprocessable Entity",
            "UNIQUE | 23505 | Antar Baris (Index) | Mencegah race condition registrasi ganda | 409 Conflict",
            "FOREIGN KEY | 23503 | Lintas Tabel | Mencegah orphan record (induk tidak ada) | 400 Reference Error",
            "PARTIAL UNIQUE | 23505 | Baris Bersyarat | Mengizinkan reuse email pasca soft-delete | 409 Conflict"
        ], "type": "table"}
    ],
    narrative="Tabel ini memetakan lima mekanisme constraint yang diuji di dalam lab beserta kode SQLSTATE kelas 23 standar PostgreSQL. Masing-masing mekanisme menangani dimensi integritas yang berbeda, mulai dari kelengkapan kolom hingga relasi lintas tabel. Menginspeksi kode SQLSTATE memungkinkan backend memetakan kegagalan teknis menjadi status HTTP yang semantik bagi klien.",
    narrative_en="This matrix charts the five constraint mechanisms verified in our lab alongside standard PostgreSQL Class 23 SQLSTATE codes. Each mechanism addresses a distinct integrity dimension, ranging from single-column presence to cross-table referential integrity. Inspecting SQLSTATE codes enables backends to map low-level database errors into semantic HTTP responses.",
    mental_model="Class 23 SQLSTATE error codes provide unambiguous machine-readable contracts for client error mapping.",
    source_path="internal/dberr/errors.go:8-99",
    slide_num="08 / 24",
    layout="table"
)

# Slide 09: Decision Sequence — 3 Focused Decision Cards (Layout: cards)
add_slide(
    categories=["DECISION SEQUENCE", "SCHEMA DESIGN"],
    title="Tiga Keputusan Kunci Merancang Integritas Skema",
    subtitle="Panduan memilih mekanisme penegakan data yang tepat pada arsitektur produksi.",
    content_items=[
        {"type": "cards", "title": "HIERARKI KEPUTUSAN PENEGAKAN INTEGRITAS DATA", "items": [
            {"title": "1. Mandatori & Predikat Baris", "desc": "Terapkan NOT NULL dan CHECK pada setiap kolom yang memiliki batasan domain eksplisit."},
            {"title": "2. Keunikan & Pola Penghapusan", "desc": "Pilih UNIQUE biasa untuk entitas permanen, atau Partial Index jika sistem menerapkan soft-delete."},
            {"title": "3. Integritas Relasional & Index", "desc": "Pasang FOREIGN KEY untuk mencegah orphan row; bukan berarti cukup, kolom relasi wajib di-index manual."}
        ]}
    ],
    narrative="Desain skema yang tangguh dimulai dari mendefinisikan kolom wajib dengan NOT NULL dan batasan nilai dengan CHECK. Ketika mengelola entitas unik yang mendukung soft-delete, gunakan partial unique index agar email lama dapat didaftarkan kembali. Memasang foreign key bukan berarti pekerjaan selesai; kolom referensi tetap wajib dibuatkan indeks terpisah agar performa query tetap optimal.",
    narrative_en="Robust schema design begins by enforcing mandatory fields via NOT NULL and permissible ranges via CHECK. When handling unique entities that support soft deletion, implement partial unique indexes so recycled identifiers remain valid. Adding foreign keys does not mean database optimization is complete; referenced foreign keys require manual indexing to sustain query performance.",
    mental_model="Design schema invariants strictly at the column, table, and relational boundaries with appropriate index backing.",
    source_path="content/02-master-draft.md:175-195",
    slide_num="09 / 24",
    layout="cards"
)

# Slide 10: Bukti Eksekusi / Terminal Output (Layout: code)
add_slide(
    categories=["EVIDENCE", "EXECUTION"],
    title="Verifikasi Eksekusi: Test Suite dan Stress Test 50 Goroutine",
    subtitle="Output terminal membuktikan penolakan atomik SQLSTATE 23505 dan kelulusan race detector.",
    content_items=[
        {"type": "code", "title": "TERMINAL OUTPUT: go test -v -race ./... && go run ./cmd/demo", "items": [
            "=== RUN   TestNotNullConstraints",
            "--- PASS: TestNotNullConstraints (0.00s)",
            "=== RUN   TestCheckConstraints",
            "--- PASS: TestCheckConstraints (0.00s)",
            "=== RUN   TestUniqueConstraint",
            "--- PASS: TestUniqueConstraint (0.00s)",
            "=== RUN   TestForeignKeyConstraint",
            "--- PASS: TestForeignKeyConstraint (0.00s)",
            "=== RUN   TestPartialUniqueIndex",
            "--- PASS: TestPartialUniqueIndex (0.00s)",
            "=== RUN   TestConcurrentRegistration_Safe_EnforcesUniqueness",
            "--- PASS: TestConcurrentRegistration_Safe_EnforcesUniqueness (0.00s)",
            "=== RUN   TestErrorClassification",
            "--- PASS: TestErrorClassification (0.00s)",
            "PASS: ok  labs/27-database-constraints/internal/store  1.330s (race: clean)",
            "[DEMO] CONCURRENCY STRESS TEST: 50 GOROUTINES FOR SAME EMAIL",
            "Results: Total=50 | Successful=1 | Rejected 23505=49 | Integrity Intact=true"
        ]}
    ],
    narrative="Bukti eksekusi menunjukkan seluruh skenario pengujian lolos bersih tanpa peringatan dari Go race detector. Pada stress test demo dengan 50 goroutine yang menyerang email yang sama persis, storage engine berhasil meloloskan tepat 1 registrasi dan menolak 49 request lainnya dengan error 23505. Hal ini mengonfirmasi integritas data terjaga 100% di bawah beban paralel tinggi.",
    narrative_en="Execution output confirms that all test suites pass with zero warnings from the Go race detector. In the demo stress test featuring 50 concurrent goroutines submitting an identical email, the storage engine admitted exactly 1 registration while rejecting 49 attempts with error 23505. This demonstrates flawless 100% data integrity under sustained parallel load.",
    mental_model="Stress testing confirms exactly one writer succeeds while all concurrent competitors receive atomic constraint rejections.",
    source_path="engineering/03-execution-result.md:14-83",
    slide_num="10 / 24",
    layout="code"
)

# Slide 11: Analisis Batasan Bukti (Layout: cards)
add_slide(
    categories=["AUDIT", "EVIDENCE BOUNDARY"],
    title="Batasan Bukti Lab: Yang Dibuktikan vs Belum",
    subtitle="Menegaskan ruang lingkup simulasi in-memory terhadap fitur engine database produksi.",
    content_items=[
        {"type": "cards", "title": "BATAS PEMBUKTIAN SISTEM LAB 27", "items": [
            {"title": "YANG DIBUKTIKAN", "desc": "Semantik SQLSTATE 23xxx, pencegahan race condition via atomic lock, dan siklus hidup partial index."},
            {"title": "YANG TIDAK DIBUKTIKAN", "desc": "Constraint terdistribusi multi-node, isolasi SERIALIZABLE, operator GiST, dan B-tree latch riil."},
            {"title": "BATAS KLAIM EMPIRIS", "desc": "Lab belum membuktikan teknik migrasi ALTER TABLE ADD CONSTRAINT NOT VALID pada tabel produksi raksasa."}
        ]}
    ],
    narrative="Lab ini membuktikan semantik SQLSTATE, keefektifan atomic locking, dan mekanisme partial index menggunakan simulator in-memory Go murni. Pengujian ini belum membuktikan constraint pada database terdistribusi multi-node atau penggunaan operator eksklusi GiST. Batasan klaim ini penting agar engineer tidak menyamakan simulator lab dengan benchmark engine produksi penuh.",
    narrative_en="This lab proves SQLSTATE semantics, atomic locking efficacy, and partial unique index mechanics using a pure Go in-memory simulator. The test suite has not yet proven distributed constraints across multi-node databases or GiST exclusion operators. Maintaining this empirical boundary ensures engineers do not conflate simulator testing with production database benchmarking.",
    mental_model="Delineate between educational memory simulations and production storage engine operational complexities.",
    source_path="content/02-master-draft.md:167-175",
    slide_num="11 / 24",
    layout="cards"
)

# Slide 12: Trade-off Arsitektur (Layout: cards)
add_slide(
    categories=["TRADE-OFF", "COSTS & BENEFITS"],
    title="Kompromi Rekayasa: Constraints vs Fleksibilitas Skema",
    subtitle="Menimbang biaya operasional dan performa di balik perlindungan integritas absolut.",
    content_items=[
        {"type": "cards", "title": "ANALISIS TRADE-OFF PENERAPAN DATABASE CONSTRAINTS", "items": [
            {"title": "Overhead Tulis (Write Cost)", "desc": "Setiap operasi INSERT/UPDATE wajib mengevaluasi index B-tree dan relasi FK, menambah sedikit latensi penulisan."},
            {"title": "Friksi Migrasi Produksi", "desc": "Menambah constraint pada tabel bervolume besar membutuhkan table lock eksklusif kecuali memakai pola validasi bertahap."},
            {"title": "Tidak Otomatis Lambat", "desc": "Penambahan index constraint tidak otomatis memperlambat sistem; overhead tulis terbayar lunas oleh integritas yang terjamin."}
        ]}
    ],
    narrative="Menerapkan database constraints menambahkan sedikit biaya komputasi pada saat INSERT karena engine harus memeriksa B-tree index dan tabel induk. Selain itu, menambahkan constraint pada tabel miliaran baris di produksi memerlukan strategi khusus agar tidak mengunci tabel terlalu lama. Namun, overhead tersebut tidak otomatis memperlambat sistem jika dibandingkan biaya besar rekonsiliasi data rusak.",
    narrative_en="Enforcing database constraints introduces minor write overhead during INSERTs as the engine inspects B-tree indexes and parent records. Furthermore, attaching constraints to massive production tables demands phased migration workflows to prevent prolonged table locks. However, this overhead does not degrade overall throughput when weighed against the catastrophic cost of data repair.",
    mental_model="Accept minor write latency at the storage tier to permanently eradicate catastrophic data reconciliation costs.",
    source_path="content/02-master-draft.md:175-186",
    slide_num="12 / 24",
    layout="cards"
)

# Slide 13: Concurrency / Validation Pipeline (Layout: schema)
add_slide(
    categories=["CONCURRENCY", "PIPELINE FLOW"],
    title="Pipeline Validasi Konkuren: Goroutine vs Mutex Lock",
    subtitle="Bagaimana 50 goroutine bersaing memperebutkan satu slot pendaftaran email unik.",
    content_items=[
        {"type": "schema", "title": "ALUR PENANGANAN 50 REQUEST PENDAFTARAN PARALEL", "items": [
            {"phase": "1. 50 Goroutines Spawn", "role": "Traffic Surge", "desc": "50 worker dijalankan serentak mencoba mendaftarkan email race-target@domain.com."},
            {"phase": "2. Mutex Lock Serialization", "role": "Atomic Critical Section", "desc": "Goroutine pertama memperoleh lock eksklusif; 49 worker lainnya mengantre tertib."},
            {"phase": "3. First Insert Success", "role": "Winner Commit", "desc": "Worker 1 lolos pengecekan, mencatatkan email ke tabel dan map index, lalu melepas lock."},
            {"phase": "4. 49 Abort SQLSTATE 23505", "role": "Atomic Rejection", "desc": "Worker 2-50 bergantian masuk lock, mendeteksi email sudah ada, dan langsung menerima error 23505."}
        ]}
    ],
    narrative="Diagram ini memvisualisasikan penanganan 50 goroutine konkuren yang mencoba menulis email yang sama persis. Antrean mutex memaksa eksekusi berjalan serial pada saat penulisan, sehingga goroutine pemenang pertama berhasil mendaftarkan data. Seluruh goroutine berikutnya secara deterministik ditolak oleh pemeriksaan indeks unik, membuktikan bahwa isolasi atomik memotong habis celah race condition.",
    narrative_en="This diagram visualizes how the storage engine handles 50 concurrent goroutines racing to register an identical email. The mutex serializes execution at the write boundary, allowing the initial winner to commit its record. All subsequent goroutines are deterministically rejected by the unique index lookup, proving that atomic isolation eliminates race windows entirely.",
    mental_model="Atomic serialization at the storage layer transforms concurrent race conditions into deterministic rejections.",
    source_path="cmd/demo/main.go:59-83",
    slide_num="13 / 24",
    layout="schema"
)

# Slide 14: Before vs After Comparison Cards (Layout: cards)
add_slide(
    categories=["COMPARISON", "BEFORE VS AFTER"],
    title="Transformasi: Dari UnsafeStore ke SafeStore",
    subtitle="Perbedaan nyata antara validasi manual di aplikasi versus deklarasi constraint database.",
    content_items=[
        {"type": "cards", "title": "PERBANDINGAN KETAHANAN DUA PENDEKATAN", "items": [
            {"title": "Sebelum: UnsafeStore (App Check)", "desc": "Mengandalkan for-loop di memori aplikasi; 20 goroutine menghasilkan >1 user duplikat di penyimpanan."},
            {"title": "Sesudah: SafeStore (DB Constraint)", "desc": "Mendelegasikan ke storage engine; 50 goroutine menghasilkan tepat 1 user dan 49 error SQLSTATE 23505."},
            {"title": "Bukan Berarti Redundan", "desc": "Mempertahankan validasi aplikasi bukan berarti pemborosan kode, melainkan filter awal bagi kenyamanan pengguna."}
        ]}
    ],
    narrative="Sebelum menerapkan database constraints, UnsafeStore mengandalkan pemindaian memori yang mudah ditembus oleh 20 goroutine konkuren hingga menghasilkan baris duplikat. Setelah beralih ke SafeStore, delegasi ke engine database menjamin tepat satu pemenang dari 50 request paralel. Mempertahankan validasi aplikasi di depan bukan berarti redundan, melainkan jembatan penting untuk menyajikan pesan error yang informatif.",
    narrative_en="Prior to implementing database constraints, UnsafeStore relied on in-memory loops that collapsed under 20 concurrent goroutines, allowing duplicate rows. Migrating to SafeStore delegates integrity to the storage engine, guaranteeing exactly one winner among 50 parallel requests. Retaining application validation upfront does not represent redundancy; it provides an essential bridge for friendly client feedback.",
    mental_model="Shift data integrity enforcement from vulnerable application loops to unyielding storage engine constraints.",
    source_path="content/02-master-draft.md:141-151",
    slide_num="14 / 24",
    layout="cards"
)

# Slide 15: Analisis Mendalam — Test Bukti Race Condition Unsafe (Layout: code)
add_slide(
    categories=["ANALYSIS", "RACE VERIFICATION"],
    title="Verifikasi Test: Membuktikan Kebocoran Data pada UnsafeStore",
    subtitle="Pengujian membuktikan bahwa pengecekan di memori gagal menahan 20 goroutine.",
    content_items=[
        {"type": "code", "title": "internal/store/store_test.go", "items": [
            "func TestConcurrentRegistration_Unsafe_SuffersRaceCondition(t *testing.T) {",
            "    eng := engine.NewEngine(); s := store.NewUnsafeStore(eng)",
            "    goroutines := 20; var wg sync.WaitGroup; wg.Add(goroutines)",
            "    for i := 0; i < goroutines; i++ {",
            "        go func(idx int) {",
            "            defer wg.Done()",
            "            _, _ = s.RegisterUser(ctx, model.User{Email: \"unsafe@example.com\", Username: fmt.Sprintf(\"u_%d\", idx)})",
            "        }(i)",
            "    }",
            "    wg.Wait()",
            "    if eng.GetUsersCount() <= 1 { t.Errorf(\"expected race condition causing >1 user entries, got %d\", eng.GetUsersCount()) }",
            "}"
        ]}
    ],
    narrative="Unit test ini membuktikan secara empiris bahwa pendekatan validasi di aplikasi gagal total di bawah konkurensi. Dua puluh goroutine yang mendaftarkan email yang sama menghasilkan lebih dari satu entri di tabel data. Kegagalan ini cakupannya mengikuti test lab, mengonfirmasi secara faktual bahwa sistem tanpa constraint penyimpanan pasti mengalami korupsi data saat menerima lonjakan request.",
    narrative_en="This unit test empirically demonstrates that application-level validation collapses under concurrent workloads. Twenty goroutines submitting an identical email produce multiple duplicate records in the underlying table. The scope of this failure aligns strictly with lab tests, proving that systems lacking database constraints suffer inevitable data corruption under traffic spikes.",
    mental_model="Empirical test assertions prove that application-only checks reliably permit duplicate writes under thread competition.",
    source_path="internal/store/store_test.go:210-239",
    slide_num="15 / 24",
    layout="code"
)

# Slide 16: Analisis Mendalam — Test Bukti Enforce Uniqueness (Layout: code)
add_slide(
    categories=["ANALYSIS", "CONSTRAINT VERIFICATION"],
    title="Verifikasi Test: SafeStore Menjamin Tepat Satu Pemenang",
    subtitle="20 goroutine paralel menghasilkan tepat 1 keberhasilan dan 19 penolakan atomik.",
    content_items=[
        {"type": "code", "title": "internal/store/store_test.go", "items": [
            "func TestConcurrentRegistration_Safe_EnforcesUniqueness(t *testing.T) {",
            "    eng := engine.NewEngine(); s := store.NewSafeStore(eng)",
            "    goroutines := 20; var wg sync.WaitGroup; wg.Add(goroutines)",
            "    successCount, errorCount := 0, 0; var mu sync.Mutex",
            "    for i := 0; i < goroutines; i++ {",
            "        go func(idx int) {",
            "            defer wg.Done()",
            "            _, err := s.RegisterUser(ctx, model.User{Email: \"safe@example.com\", Username: fmt.Sprintf(\"u_%d\", idx), Age: 25, Status: \"active\"})",
            "            mu.Lock(); defer mu.Unlock()",
            "            if err == nil { successCount++ } else { errorCount++ }",
            "        }(i)",
            "    }",
            "    wg.Wait()",
            "    if successCount != 1 || errorCount != 19 || eng.GetUsersCount() != 1 { t.Fatalf(\"integrity violated!\") }",
            "}"
        ]}
    ],
    narrative="Pengujian ini membuktikan ketangguhan SafeStore yang mengandalkan constraint engine. Dari 20 goroutine yang dijalankan serentak, tepat satu pendaftaran yang berhasil ditulis ke penyimpanan, sementara 19 goroutine lainnya menerima penolakan. Hanya pada jalur ini integritas tabel tetap terjaga pada angka tepat 1 baris, membuktikan keampuhan proteksi level database.",
    narrative_en="This test verifies the resilience of SafeStore backed by engine constraints. Among 20 concurrently spawned goroutines, exactly one registration succeeds in persisting to storage, while the remaining 19 receive immediate rejections. Only along this path is table integrity strictly preserved at exactly 1 row, confirming the power of database-level protection.",
    mental_model="Database constraints deterministically ensure a single commit winner and N-1 clean rejections under high concurrency.",
    source_path="internal/store/store_test.go:162-208",
    slide_num="16 / 24",
    layout="code"
)

# Slide 17: Analisis Mendalam — Partial Unique Index & Soft Delete (Layout: code)
add_slide(
    categories=["ANALYSIS", "PARTIAL INDEX"],
    title="Solusi Soft-Delete: Partial Unique Index WHERE deleted_at IS NULL",
    subtitle="Mengizinkan daur ulang email lama tanpa mengorbankan keunikan data pengguna aktif.",
    content_items=[
        {"type": "code", "title": "internal/engine/engine.go", "items": [
            "// PARTIAL UNIQUE: CREATE UNIQUE INDEX users_active_email_idx ON users (email) WHERE deleted_at IS NULL;",
            "if usePartialUniqueIndex {",
            "    if u.DeletedAt == nil {",
            "        if existingID, exists := e.activeEmails[u.Email]; exists {",
            "            return model.User{}, dberr.NewUniqueViolation(\"users\", \"users_active_email_idx\", \"duplicate active email\")",
            "        }",
            "    }",
            "}",
            "// Saat soft delete: delete(e.activeEmails, u.Email) -> email langsung dapat didaftarkan kembali",
            "// Baris dengan deleted_at != nil tidak diindeks, sehingga multi-duplikat soft-delete tetap valid."
        ]}
    ],
    narrative="Penerapan UNIQUE biasa pada tabel dengan soft-delete menimbulkan masalah: email pengguna yang sudah dihapus tidak akan pernah bisa didaftarkan kembali. Dengan Partial Unique Index bersyarat WHERE deleted_at IS NULL, hanya akun aktif yang diperiksa keunikannya. Pengguna yang sudah soft-delete dikeluarkan dari indeks, sehingga email yang sama dapat dipakai ulang oleh pendaftar baru.",
    narrative_en="Applying a standard UNIQUE constraint to tables with soft deletion creates severe friction: deactivated emails can never be re-registered. By leveraging a conditional Partial Unique Index with WHERE deleted_at IS NULL, the database enforces uniqueness exclusively on active rows. Deactivated records are evicted from the index, allowing identifier recycling without conflict.",
    mental_model="Partial indexes scope uniqueness predicates exclusively to active records, resolving soft-delete reuse dilemmas.",
    source_path="internal/engine/engine.go:71-77",
    slide_num="17 / 24",
    layout="code"
)

# Slide 18: Analisis Mendalam — Foreign Key Integrity (Layout: code)
add_slide(
    categories=["ANALYSIS", "REFERENTIAL INTEGRITY"],
    title="Mencegah Data Yatim Piatu dengan Foreign Key Constraint",
    subtitle="Memastikan pesanan baru tidak dapat dibuat untuk pengguna yang tidak terdaftar.",
    content_items=[
        {"type": "code", "title": "internal/engine/engine.go", "items": [
            "func (e *Engine) InsertOrder(o model.Order) (model.Order, error) {",
            "    e.mu.Lock(); defer e.mu.Unlock()",
            "    if o.UserID == 0 { return model.Order{}, dberr.NewNotNullViolation(\"orders\", \"user_id\", \"orders_user_id_not_null\") }",
            "    if o.TotalCents <= 0 { return model.Order{}, dberr.NewCheckViolation(\"orders\", \"orders_total_cents_check\", \"total > 0\") }",
            "    // FOREIGN KEY: CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES users(id)",
            "    if _, exists := e.users[o.UserID]; !exists {",
            "        return model.Order{}, dberr.NewForeignKeyViolation(\"orders\", \"fk_orders_user\", fmt.Sprintf(\"key (user_id)=(%d) not present\", o.UserID))",
            "    }",
            "    o.ID = e.orderSeq.Add(1); e.orders[o.ID] = o; return o, nil",
            "}"
        ]}
    ],
    narrative="Fungsi InsertOrder mengevaluasi tiga batasan secara berurutan: NOT NULL pada UserID, CHECK pada TotalCents, dan verifikasi relasional FOREIGN KEY. Jika UserID bernilai 99999 yang tidak tercatat di tabel users, sistem mengembalikan SQLSTATE 23503. Mekanisme ini mencegah terciptanya order yatim piatu yang dapat merusak rekonsiliasi keuangan sistem.",
    narrative_en="The InsertOrder function sequentially evaluates three constraints: NOT NULL on UserID, CHECK on TotalCents, and referential FOREIGN KEY validation. If a UserID of 99999 is submitted without existing in the users table, the engine returns SQLSTATE 23503. This mechanism prevents orphaned order records that could compromise financial reconciliation.",
    mental_model="Foreign keys enforce relational parentage at the storage boundary, eradicating orphaned dependent records.",
    source_path="internal/engine/engine.go:122-148",
    slide_num="18 / 24",
    layout="code"
)

# Slide 19: Analisis Mendalam — SQLSTATE Taxonomy & Domain Mapping (Layout: code)
add_slide(
    categories=["ANALYSIS", "ERROR HANDLING"],
    title="Inspeksi SQLSTATE: Mengubah Error Database Menjadi Domain Error",
    subtitle="Hindari parsing teks error yang rapuh; periksa kode SQLSTATE Class 23 yang stabil.",
    content_items=[
        {"type": "code", "title": "internal/dberr/errors.go", "items": [
            "func MapToDomainError(err error) error {",
            "    var cErr *ConstraintError",
            "    if errors.As(err, &cErr) {",
            "        switch cErr.Code {",
            "        case SQLStateUniqueViolation: // 23505",
            "            return fmt.Errorf(\"conflict: resource with this unique attribute already exists (rule: %s)\", cErr.ConstraintName)",
            "        case SQLStateNotNullViolation: // 23502",
            "            return fmt.Errorf(\"invalid input: mandatory field is missing (rule: %s)\", cErr.ConstraintName)",
            "        case SQLStateCheckViolation: // 23514",
            "            return fmt.Errorf(\"validation failed: value outside boundary (rule: %s)\", cErr.ConstraintName)",
            "        case SQLStateForeignKeyViolation: // 23503",
            "            return fmt.Errorf(\"reference error: referenced entity does not exist (rule: %s)\", cErr.ConstraintName)",
            "        }",
            "    }",
            "    return err",
            "}"
        ]}
    ],
    narrative="Banyak pengembang melakukan parsing string error seperti duplicate key yang sangat rapuh terhadap pergantian bahasa basis data. Kode SQLSTATE Class 23 bersifat universal dan terstandarisasi. Menggunakan MapToDomainError memungkinkan aplikasi mengubah kode mesin menjadi pesan domain yang bersih tanpa mengekspos struktur internal database ke publik.",
    narrative_en="Many engineers parse raw error strings like duplicate key, which break whenever database locales change. Class 23 SQLSTATE codes provide a standardized, permanent machine contract. Employing MapToDomainError allows applications to translate engine codes into clean domain errors without leaking internal database schemas to clients.",
    mental_model="Inspect stable SQLSTATE codes rather than brittle localized text strings for robust application error mapping.",
    source_path="internal/dberr/errors.go:83-100",
    slide_num="19 / 24",
    layout="code"
)

# Slide 20: Checklist Implementasi (Layout: checklist)
add_slide(
    categories=["CHECKLIST", "PRODUCTION READINESS"],
    title="Checklist Penerapan Database Constraints di Produksi",
    subtitle="Langkah-langkah rekayasa wajib untuk menjamin integritas data sebelum deployment.",
    content_items=[
        {"type": "checklist", "title": "STANDAR REKAYASA INTEGRITAS SKEMA DATABASE", "items": [
            "Tandai seluruh kolom mandatori dengan NOT NULL, jangan hanya mengandalkan validasi form.",
            "Terapkan UNIQUE constraint pada atribut unik alih-alih melakukan SELECT sebelum INSERT.",
            "Gunakan Partial Unique Index (WHERE deleted_at IS NULL) untuk entitas yang mendukung soft-delete.",
            "Definisikan CHECK constraint untuk membatasi rentang nilai angka dan enum status kolom.",
            "Pasang indeks eksplisit pada kolom referencing FOREIGN KEY demi mencegah full table scan saat DELETE parent.",
            "Periksa kode SQLSTATE Class 23 (23502, 23505, 23514, 23503) pada error handler aplikasi."
        ]}
    ],
    narrative="Checklist ini merangkum langkah-langkah praktis sebelum merilis skema ke produksi: penetapan NOT NULL pada kolom wajib, penerapan UNIQUE untuk proteksi race condition, pemanfaatan partial index untuk soft-delete, pemasangan CHECK untuk validasi baris, penambahan indeks pada foreign key, serta inspeksi SQLSTATE pada layer error handling.",
    narrative_en="This checklist summarizes mandatory engineering steps before releasing schemas to production: declaring NOT NULL on mandatory fields, applying UNIQUE constraints to block race conditions, leveraging partial indexes for soft deletes, defining CHECK constraints for row limits, indexing foreign key columns, and inspecting SQLSTATE codes in application error handlers.",
    mental_model="Systematically verify declarative constraints, index backing, and structured error mapping across the persistence tier.",
    source_path="content/02-master-draft.md:204-215",
    slide_num="20 / 24",
    layout="checklist"
)

# Slide 21: Anti-Pattern & Cara Reproduksi Kegagalan (Layout: cards)
add_slide(
    categories=["ANTI-PATTERN", "PITFALLS"],
    title="Anti-Pattern Umum dalam Menjaga Integritas Data",
    subtitle="Kesalahan desain arsitektur yang sering merusak konsistensi data di sistem backend.",
    content_items=[
        {"type": "cards", "title": "POLA BERBAHAYA YANG WAJIB DIHINDARI", "items": [
            {"title": "Check-Then-Act di Aplikasi", "desc": "Mengecek keberadaan record via SELECT di aplikasi lalu melakukan INSERT; pasti bocor di bawah beban konkuren."},
            {"title": "Menebak String Error", "desc": "Melakukan regex matching pada teks error database yang dapat berubah sewaktu-waktu akibat pembaruan versi DB."},
            {"title": "Bukan Berarti Ter-Index", "desc": "Membuat FOREIGN KEY bukan berarti kolom otomatis terindeks; ketiadaan indeks manual merusak performa query DELETE."}
        ]}
    ],
    narrative="Anti-pattern paling fatal adalah merasa aman dengan pola SELECT sebelum INSERT pada aplikasi. Kesalahan lainnya mencakup parsing teks error alih-alih membaca SQLSTATE, serta berasumsi bahwa foreign key secara otomatis membuatkan indeks pada kolom anak. Mengabaikan aspek ini bukan berarti sistem rusak hari ini, melainkan menanam bom waktu performa di masa depan.",
    narrative_en="The most perilous anti-pattern is assuming safety from application-level SELECT-then-INSERT flows. Other prevalent traps include parsing localized error strings rather than SQLSTATE codes and assuming foreign keys automatically index referencing columns. Overlooking these realities creates ticking performance and data corruption time-bombs.",
    mental_model="Avoid application-only uniqueness checks, error string parsing, and unindexed foreign key dependencies.",
    source_path="content/02-master-draft.md:186-195",
    slide_num="21 / 24",
    layout="cards"
)

# Slide 22: Takeaway Praktis & Command Block Dominan (Layout: code)
add_slide(
    categories=["REPRODUCTION", "COMMANDS"],
    title="Reproduksi Mandiri di Mesin Lokal Menggunakan Go",
    subtitle="Jalankan seluruh pengujian constraint, race condition, dan demo 50 worker tanpa dependensi eksternal.",
    content_items=[
        {"type": "code", "title": "COMMAND RUNNER (ZERO DEPENDENCIES)", "items": [
            "cd labs/27-database-constraints",
            "go test -v ./...    # menjalankan 7 unit test untuk NOT NULL, CHECK, UNIQUE, FK, & Partial Index",
            "go test -race ./... # memvalidasi eksekusi konkuren bebas dari data race (PASS in 1.33s)",
            "go run ./cmd/demo   # eksekusi demonstrasi CLI interaktif: SQLSTATE demo & 50-worker stress test"
        ]}
    ],
    narrative="Seluruh pengujian lab dibangun murni menggunakan Go standard library tanpa perlu mengunduh package eksternal atau menjalankan Docker. Eksperimen lokal ini membuktikan seluruh semantik penegakan constraint secara deterministik, namun sistem produksi masih perlu mengonfigurasi strategi migrasi bertahap seperti ADD CONSTRAINT NOT VALID.",
    narrative_en="The entire lab simulation is engineered strictly with the Go standard library, requiring no external packages or Docker containers. While local tests prove constraint semantics deterministically, production environments still require phased migration strategies such as ADD CONSTRAINT NOT VALID.",
    mental_model="Deterministic local simulations validate distributed consistency principles without heavy infrastructure overhead.",
    source_path="README.md:1-25",
    slide_num="22 / 24",
    layout="code"
)

# Slide 23: Prinsip Arsitektural Jangka Panjang (Layout: cards)
add_slide(
    categories=["TAKEAWAY", "ARCHITECTURAL PRINCIPLES"],
    title="Prinsip Inti: Database Adalah Benteng Terakhir Kebenaran",
    subtitle="Aplikasi boleh berganti dan berkembang, namun integritas data harus abadi.",
    content_items=[
        {"type": "cards", "title": "TIGA PILAR INTEGRITAS DATA ARSITEKTUR MODERN", "items": [
            {"title": "Deklaratif di Titik Tulis", "desc": "Letakkan aturan keabsahan data sedekat mungkin dengan media penyimpanan agar terlindungi secara atomik."},
            {"title": "Pemisahan Peran yang Jelas", "desc": "Aplikasi melayani kenyamanan interaksi pengguna; database constraints menjamin kebenaran absolut."},
            {"title": "Tidak Otomatis Kebal", "desc": "Memasang constraint tidak otomatis menyelesaikan aturan lintas tabel; pahami batas row-scoped versus triggers."}
        ]}
    ],
    narrative="Aplikasi backend akan terus ditulis ulang, microservices akan terus bertambah, namun data di database akan bertahan melampaui siklus hidup kode tersebut. Menyerahkan integritas data hanya pada kode aplikasi adalah perjudian berisiko tinggi. Database constraints memastikan bahwa apa pun bug yang terjadi di lapisan aplikasi, integritas data penyimpanan tidak akan pernah ternoda.",
    narrative_en="Application layers will inevitably be refactored and microservices will proliferate, but persisted data outlives software lifecycles. Entrusting data integrity solely to application code is a high-risk architectural gamble. Database constraints guarantee that regardless of application-layer regressions, persistent data remains uncorrupted.",
    mental_model="Applications are ephemeral; persisted data is durable. Enforce invariant truth at the storage tier.",
    source_path="content/05-key-takeaways.md:1-20",
    slide_num="23 / 24",
    layout="cards"
)

# Slide 24: Closing Slide (Layout: closing)
add_slide(
    categories=["CLOSING"],
    title="Jangan pernah percayakan integritas data hanya pada application code; database constraint adalah benteng terakhir kebenaran absolut.",
    subtitle="",
    content_items=[
        {"type": "bullet", "title": "MODUL IMPLEMENTASI", "items": [
            "internal/engine/engine.go • internal/store/store.go",
            "internal/dberr/errors.go • internal/model/model.go",
            "cmd/demo/main.go • internal/store/store_test.go"
        ]},
        {"type": "bullet", "title": "VERIFIKASI & STATUS", "items": [
            f"Commit: {COMMIT_SHA[:7]} • Status: APPROVED",
            "Test passing: 8 test scenarios (NotNull, Check, Unique, FK, PartialIndex, ConcurrencySafe, ConcurrencyUnsafe, ErrorClassification)",
            "Race detector: Clean (go test -race ./... PASS in 1.330s)",
            "Zero external dependencies: 100% Go standard library."
        ]}
    ],
    mental_model="Database constraints serve as the final, immutable bastion of data correctness across evolving system architectures.",
    source_path=REPO_URL,
    slide_num="24 / 24",
    layout="closing"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
    for s in SLIDES:
        print(f"[{s['slide_num']}] {s['layout']}: {s['title']}")
