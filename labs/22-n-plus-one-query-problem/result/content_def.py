#!/usr/bin/env python3
"""
Publication Builder for Software Engineering Lab 22: N+1 Query Problem.
Generates: PPTX, PDF, PNG covers, Dev.to MD, Medium HTML, Metadata, Caption, DELIVERABLES.md
"""
import os
import sys
import math
import subprocess
from datetime import datetime
from pptx import Presentation
from pptx.util import Inches, Pt
from pptx.dml.color import RGBColor
from pptx.enum.text import PP_ALIGN
from pptx.enum.shapes import MSO_AUTO_SHAPE_TYPE, MSO_CONNECTOR

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

# Color constants
DARK_NAVY = "0b1b2d"
TEAL = "087f76"
CREAM = "f4f1e8"
WHITE = "ffffff"
SLATE = "52606b"
ORANGE = "d97706"

# Slide data structure
SLIDES = []

def add_slide(categories, title, subtitle, content_items, mental_model, source_path="", slide_num=""):
    """Add a slide to the presentation."""
    SLIDES.append({
        "categories": categories,
        "title": title,
        "subtitle": subtitle,
        "content_items": content_items,  # list of dicts: {"title": str, "items": [str, str, ...], "type": "bullet"|"code"|"table"|"diagram"}
        "mental_model": mental_model,
        "source": source_path,
        "slide_num": slide_num
    })


# =====================================================================
# SLIDE CONTENT DEFINITIONS
# =====================================================================

# Slide 1: Cover (Dark Navy, visual element)
add_slide(
    categories=["KONTEKS"],
    title=TITLE_ID,
    subtitle="SOFTWARE ENGINEERING LAB | LAB " + LAB_NO,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "Retro CRT monitor beige dengan scanlines",
            "ASCII art terminal yang menampilkan query count",
            f"Penulis: {AUTHOR}",
        ]}
    ],
    mental_model="",
    source_path="",
    slide_num=""
)

# Slide 2: Motivation & Problem Definition
add_slide(
    categories=["MOTIVASI", "KONSEP"],
    title="N+1 Query Problem: Anti-Pattern Latensi Tersembunyi",
    subtitle="Satu query untuk mengambil koleksi, lalu N query terpisah untuk setiap item",
    content_items=[
        {"title": "MOTIVASI", "items": [
            "N+1 adalah masalah 'invisible until it hurts' — berfungsi sempurna pada data uji kecil",
            "Setiap query individu terlihat cepat (<1ms), sehingga lolos dari slow query logs",
            "Degradasi performa akumuler: request yang harus 200ms justru butuh 2.4 detik di produksi",
        ], "type": "bullet"},
        {"title": "KONSEP", "items": [
            "Anti-pattern N+1: 1 query induk + N query anak per item",
            "Kompleksitas menjadi O(N) relatif terhadap jumlah record",
            "Sering terjadi di ORM lazy loading default dan GraphQL resolvers",
        ], "type": "bullet"},
        {"title": "STUDI KASUS", "items": [
            "GET /api/work-orders: 100 work order, 5 item tiapnya =",
            ": 101 query (1 + 100) vs 2 query (batching)",
            "Pada data uji 5 work order: hanya 6 query — tidak terlihat masalahnya",
        ], "type": "bullet"},
    ],
    mental_model="N+1 is invisible until it hurts. Batch related queries, don't loop them.",
    source_path="labs/22-n-plus-one-query-problem/README.md",
    slide_num="2 / 24"
)

# Slide 3: Mechanism & Anti-Pattern Anatomy
add_slide(
    categories=["BUKTI LOGIKA"],
    title="Anatomi N+1: 1 + N = N+1",
    subtitle="Loop-driven retrieval memicu eksekusi query berulang di dalam iterasi",
    content_items=[
        {"title": "ALUR EKSEKUSI N+1", "items": [
            "Query #1: GetAllAuthors() → mengembalikan 3 author",
            "Loop iterasi ke-1: GetPostsByAuthorID(1) → Query #2",
            "Loop iterasi ke-2: GetPostsByAuthorID(2) → Query #3",
            "Loop iterasi ke-3: GetPostsByAuthorID(3) → Query #4",
            "Total: 4 query (1 + N di mana N=3)",
        ], "type": "diagram"},
        {"title": "KODE ANTI-PATTERN", "items": [
            "authors := store.GetAllAuthors()",
            "for _, author := range authors {",
            "    posts := store.GetPostsByAuthorID(author.ID)",
            "}",
        ], "type": "code"},
        {"title": "SKALA NYATA", "items": [
            "Jika N=1000 author, query count = 1001",
            "Jika N=10000 author, query count = 10001",
            "Linear scaling dengan jumlah record — tidak berkelakuan baik",
        ], "type": "bullet"},
    ],
    mental_model="Each item in a loop triggers its own query — batch before you iterate.",
    source_path="labs/22-n-plus-one-query-problem/internal/blog/repository.go:13",
    slide_num="3 / 24"
)

# Slide 4: Operational Blind Spot
add_slide(
    categories=["MASALAH OPERASIONAL"],
    title="Titik Buta: Slow Query Log Gagal Mendeteksi",
    subtitle="Query individu terlihat sempurna, tapi akumulasi menyebabkan bottleneck",
    content_items=[
        {"title": "MENGAPA LOLOS DETEKSI", "items": [
            "Slow query log hanya menandai query dengan durasi > threshold (mis. 1s)",
            "Query N+1 masing-masing <1ms, kecil di atas kertas",
            "Akumulasi N query kecil menghasilkan total latency besar — invisible dalam log",
        ], "type": "bullet"},
        {"title": "CONTOH PUKUL KETIKA PRODUKSI", "items": [
            "Development: 3 author, 4 query → <5ms",
            "Production: 500 author, 501 query → 2.4s",
            "Slow query log melaporkan 0 query lambat — semuanya cepat",
        ], "type": "bullet"},
        {"title": "DETEKSI YANG BENAR", "items": [
            "Query count profiling, bukan hanya durasi individual",
            "ORM-level query logger (mis. Django nplusone)",
            "Instrumentation APM yang mengobservasi pola query",
        ], "type": "bullet"},
    ],
    mental_model="Fast individual queries can mask slow aggregate performance.",
    source_path="labs/22-n-plus-one-query-problem/content/01-content-brief.md",
    slide_num="4 / 24"
)

# Slide 5: Why It Reaches Production
add_slide(
    categories=["AKAR MASALAH"],
    title="Mengapa N+1 Lolos ke Produksi",
    subtitle="Tiga faktor struktural: data uji kecil, optimisasi prematur, profiling kurang",
    content_items=[
        {"title": "FAKTOR 1: DATA UJI KECIL", "items": [
            "Lab & staging memakai data sample 3-5 record",
            "4 query terlihat cepat — tidak ada bottleneck yang terlihat",
            "Prinsip 'works on my machine' menyebabkan kegagalan skalasi",
        ], "type": "bullet"},
        {"title": "FAKTOR 2: ORM DEFAULT LAZY", "items": [
            "Django: lazy loading default untuk relasi",
            "Rails: lazy loading default untuk association",
            "EF Core: lazy loading proxy (opsional tapi sering diaktifkan)",
        ], "type": "bullet"},
        {"title": "FAKTOR 3: PROFILING TIDAK STANDAR", "items": [
            "Query count tidak diukur sebagai metrik CI/CD",
            "Fokus pada duration, bukan volume per request",
            "Tidak ada test assertion untuk jumlah query",
        ], "type": "bullet"},
    ],
    mental_model="Small test data hides linear degradation — profile volume, not just speed.",
    source_path="labs/22-n-plus-one-query-problem/content/02-master-draft.md",
    slide_num="5 / 24"
)

# Slide 6: Architecture
add_slide(
    categories=["DESAIN SISTEM"],
    title="Arsitektur Lab 22: Tiga Layer",
    subtitle="Pisahan Store (mock DB), Repository (data access), Demo (CLI runner)",
    content_items=[
        {"title": "LAYER STRUKTUR", "items": [
            "Store: Internal/blog/store.go — mock database + query counter",
            "Repository: internal/blog/repository.go — dua strategi (naive vs eager)",
            "Demo: cmd/demo/main.go — CLI runner yang mencetak statistik query",
        ], "type": "bullet"},
        {"title": "THREAD SAFETY", "items": [
            "Store menggunakan sync.Mutex untuk query counter",
            "Query count increment thread-safe — valid untuk concurrent test",
            "Race detector lolos tanpa data race (go test -race)",
        ], "type": "bullet"},
        {"title": "VERIFIKASI KWANTITATIF", "items": [
            "Query counter bukan estimasi — angka pasti",
            "Test assertion: N+1=4 query, Eager=2 query (hard-coded)",
            "Deep equality check: reflect.DeepEqual(eager, nplus1)",
        ], "type": "bullet"},
    ],
    mental_model="Make query counts measurable so assertions are exact, not estimated.",
    source_path="labs/22-n-plus-one-query-problem/internal/blog/store.go",
    slide_num="6 / 24"
)

# Slide 7: Naive Implementation
add_slide(
    categories=["BUKTI KODE"],
    title="Implementasi Naif: Lazy Loading per Item",
    subtitle="GetPostsByAuthorID dipanggil di dalam loop — masing-masing menghasilkan satu query",
    content_items=[
        {"title": "REPOSITORY N+1", "items": [
            "authors := r.store.GetAllAuthors()  // Query #1",
            "for _, author := range authors {",
            "    posts := r.store.GetPostsByAuthorID(author.ID) // Query #2,#3,#4",
            "    result = append(result, AuthorWithPosts{...})",
            "}",
        ], "type": "code"},
        {"title": "STORE QUERY COUNT", "items": [
            "GetAllAuthors(): queryCount++ → 1",
            "GetPostsByAuthorID(id): queryCount++ per panggilan",
            "3 author → 3 pemanggilan → 3 query tambahan",
            "Total akumulasi: 4",
        ], "type": "code"},
        {"title": "KARAKTERISTIK", "items": [
            "Intuisi kode natural / readable — ini yang membuatnya sulit diketahui",
            "Mirip dengan ORM lazy loading default",
            "Sering ditulis oleh developer junior yang belum sadar",
        ], "type": "bullet"},
    ],
    mental_model="Lazy loading in a loop multiplies query count — batch before iterating.",
    source_path="labs/22-n-plus-one-query-problem/internal/blog/repository.go:13",
    slide_num="7 / 24"
)

# Slide 8: Eager Loading Solution
add_slide(
    categories=["BUKTI KODE"],
    title="Solusi Eager Loading: Batch & Memory Map",
    subtitle="Kumpulkan semua ID, satu query IN(...), kelompokkan di memori",
    content_items=[
        {"title": "REPOSITORY EAGER", "items": [
            "authors := r.store.GetAllAuthors()  // Query #1",
            "ids := extractAuthorIDs(authors)    // slice [1,2,3]",
            "allPosts := r.store.GetPostsByAuthorIDs(ids)  // Query #2 batch",
            "groupByAuthor(allPosts) → map[int][]Post",
            "Build final result from map",
        ], "type": "code"},
        {"title": "STORE BATCH METHOD", "items": [
            "GetPostsByAuthorIDs(authorIDs []int):",
            "    idMap := make(map[int]bool)",
            "    for _, id := range authorIDs { idMap[id] = true }",
            "    return posts where idMap[post.AuthorID]",
            "    queryCount++ once — satu query untuk semua",
        ], "type": "code"},
        {"title": "HASIL", "items": [
            "Total query: 2 (1 untuk author + 1 untuk batch post)",
            "Terlepas dari berapa banyak author — selalu 2 query",
            "Data hasil identik dengan N+1 (verified via DeepEqual)",
        ], "type": "bullet"},
    ],
    mental_model="Collect all IDs, fetch in one IN(...) batch, join in memory — O(1) queries.",
    source_path="labs/22-n-plus-one-query-problem/internal/blog/repository.go:32",
    slide_num="8 / 24"
)

# Slide 9: Quantitative Results
add_slide(
    categories=["HASIL PENGUJIAN"],
    title="Eksekusi Demo: Bukti Kuantitatif",
    subtitle="Output CLI: query count dan jumlah record yang dimuat",
    content_items=[
        {"title": "DEMO OUTPUT — N+1", "items": [
            "--- Simulating N+1 Query Problem ---",
            "Loaded 3 authors with their posts.",
            "Total queries executed: 4 (1 for authors + 3 for posts)",
        ], "type": "code"},
        {"title": "DEMO OUTPUT — EAGER", "items": [
            "--- Simulating Eager Loading (Batching) ---",
            "Loaded 3 authors with their posts.",
            "Total queries executed: 2 (1 for authors + 1 batched query)",
        ], "type": "code"},
        {"title": "INTERPRETASI", "items": [
            "N+1: query count = N+1 (bergantung pada jumlah record)",
            "Eager: query count = 2 (konstan, tidak bergantung pada N)",
            "Pada N=1000: N+1=1001 query vs Eager=2 query — perbedaan 999 query",
        ], "type": "bullet"},
    ],
    mental_model="Constant 2 queries beats N+1, and the gap widens linearly with N.",
    source_path="labs/22-n-plus-one-query-problem/cmd/demo/main.go",
    slide_num="9 / 24"
)

# Slide 10: Test Verification
add_slide(
    categories=["VERIFIKASI TEST"],
    title="Unit Test: Assert Kuantitatif Query Count",
    subtitle="Deep equality & query assertion mencegah regresi secara otomatis",
    content_items=[
        {"title": "TEST RESULTS", "items": [
            "=== RUN   TestGetAuthorsWithPostsNPlusOne",
            "--- PASS: TestGetAuthorsWithPostsNPlusOne (0.00s)",
            "=== RUN   TestGetAuthorsWithPostsEager",
            "--- PASS: TestGetAuthorsWithPostsEager (0.00s)",
            "=== RUN   TestEmptyStore",
            "--- PASS: TestEmptyStore (0.00s)",
            "PASS — 4/4 test cases pass",
        ], "type": "code"},
        {"title": "ASSERTI QUERY COUNT", "items": [
            "N+1: queryCount == 4 (1 + N, di mana N=3)",
            "Eager: queryCount == 2 (1 + 1 batch)",
            "Assertion hard-coded: tidak estimasi — angka pasti",
        ], "type": "code"},
        {"title": "DEEP EQUALITY", "items": [
            "reflect.DeepEqual(eagerResult, nPlusOneResult) → true",
            "Batching tidak mengubah semantik data — hanya mengurangi I/O",
            "Race detector: go test -race — lulus tanpa data race",
        ], "type": "bullet"},
    ],
    mental_model="Assert exact query counts in unit tests to prevent N+1 regressions.",
    source_path="labs/22-n-plus-one-query-problem/internal/blog/repository_test.go",
    slide_num="10 / 24"
)

# Slide 11: Complexity & Scalability
add_slide(
    categories=["ANALISIS TEKNIK"],
    title="Kompleksitas: O(N) → O(1) Query Count",
    subtitle="Grafik query count vs jumlah record: linear vs konstan",
    content_items=[
        {"title": "GRAFIK KOMPLEKSITAS", "items": [
            "N=3:   N+1=4 query  | Eager=2 query",
            "N=100: N+1=101 query | Eager=2 query",
            "N=1K:  N+1=1001 query| Eager=2 query",
            "N=10K: N+1=10001 query| Eager=2 query",
        ], "type": "table"},
        {"title": "ANALISIS O-NOTASI", "items": [
            "N+1: Jumlah query = O(N) relatif terhadap jumlah record",
            "Eager: Jumlah query = O(1) — konstan, tidak bergantung pada N",
            "Memory: O(N) keduanya — sama, hanya query saja yang berbeda",
        ], "type": "bullet"},
        {"title": "FACTOR KONSTAN", "items": [
            "Eager memiliki overhead memory untuk map building",
            "Tetapi menghemat ribuan round-trip jaringan DB",
            "Nilai tukar: memory vs latency — biasanya latency lebih mahal",
        ], "type": "bullet"},
    ],
    mental_model="Query count complexity drops from O(N) to O(1) — the gap widens linearly.",
    source_path="labs/22-n-plus-one-query-problem/content/05-key-takeaways.md",
    slide_num="11 / 24"
)

# Slide 12: Beyond Database: API & GraphQL N+1
add_slide(
    categories=["EKSPANSI MASALAH"],
    title="N+1 di Luar Database: API & GraphQL Resolvers",
    subtitle="Pola yang sama muncul pada HTTP calls dan GraphQL field resolvers",
    content_items=[
        {"title": "GRAPHQL N+1", "items": [
            "Resolver untuk posts(author) dipanggil per author",
            "Setiap resolver = 1 HTTP request ke backend",
            "Solusi standar: DataLoader (Facebook) — batch & cache",
        ], "type": "bullet"},
        {"title": "HTTP MICROSERVICE N+1", "items": [
            "Endpoint mengembalikan 100 order, lalu fetch item per order",
            "100 HTTP request seril, latensi menumpuk",
            "Solusi: batch endpoint /items?order_id=1,2,3,... atau aggregation service",
        ], "type": "bullet"},
        {"title": "POLA YANG SAMA", "items": [
            "Database: ORM lazy loading",
            "HTTP: per-object request di dalam loop",
            "GraphQL: per-field resolver di dalam list",
            "Akar masalah: eksekusi I/O berulang di dalam iterasi",
        ], "type": "bullet"},
    ],
    mental_model="N+1 is a pattern, not a layer — it happens wherever I/O runs inside a loop.",
    source_path="labs/22-n-plus-one-query-problem/content/05-key-takeaways.md",
    slide_num="12 / 24"
)

# Slide 13: DataLoader Pattern
add_slide(
    categories=["SOLUSI JARINGAN"],
    title="DataLoader: Request De-duplication & Batching",
    subtitle="Pola batching di level aplikasi untuk GraphQL dan HTTP",
    content_items=[
        {"title": "DATA LOADER (GRAPHQL)", "items": [
            "Collect requests selama satu tick event loop",
            "Batch semua request dengan key yang sama ke satu backend call",
            "Cache hasil per key — request berulang tidak menghitung lagi",
        ], "type": "bullet"},
        {"title": "HTTP BATCHING EQUIVALENT", "items": [
            "Aggregate IDs, kirim satu request batch ke backend",
            "Gunakan cache layer (Redis) untuk hasil yang relatif statis",
            "Pattern ini sama dengan database batch query — konsep universal",
        ], "type": "bullet"},
        {"title": "KONEKSI KE LAB", "items": [
            "Lab ini: database-level batch via GetPostsByAuthorIDs",
            "DataLoader: application-level batch via IN(...) equivalent",
            "Kedua pola: kumpulkan, batch, satukan di memori",
        ], "type": "bullet"},
    ],
    mental_model="Batch and cache: collect requests across a tick, fire one call, cache per key.",
    source_path="labs/22-n-plus-one-query-problem/content/05-key-takeaways.md",
    slide_num="13 / 24"
)

# Slide 14: Mapping-Level Eager Loading Danger
add_slide(
    categories=["PERINGATAN ARSITEKTUR"],
    title="Bahaya Mapping-Level Eager Loading",
    subtitle="FetchType.EAGER dan semacamnya memuat semua relasi — memory bloat",
    content_items=[
        {"title": "MASALAH", "items": [
            "JPA FetchType.EAGER: memuat semua relasi otomatis di setiap query",
            "Tidak ada kontrol selektif — over-fetching terjadi",
            "Memory bloat: entitas tidak relevan ikut ter-load",
        ], "type": "bullet"},
        {"title": "PERBEDAAN KRITIS", "items": [
            "Mapping-level: deklaratif di class — selalu eager",
            "Query-level: eksplisit per query — kontrol penuh",
            "Lab ini mendemonstrasikan query-level — solusi yang benar",
        ], "type": "bullet"},
        {"title": "CONTOH JPA", "items": [
            "Entity Author dengan @OneToMany EAGER",
            "Setiap kali query Author, semua Post ikut ter-load",
            "Bayangkan: 100 author, masing-masing 1000 post = 100K post loaded",
        ], "type": "bullet"},
    ],
    mental_model="Eager at the mapping level loads everything everywhere — control loading per query.",
    source_path="labs/22-n-plus-one-query-problem/content/02-master-draft.md",
    slide_num="14 / 24"
)

# Slide 15: Cartesian Explosion Trade-off
add_slide(
    categories=["TRADE-OFF TEKNIK"],
    title="Trade-off: Cartesian Product Explosion",
    subtitle="Eager loading via JOIN menyebabkan hasil dikali jumlah child records",
    content_items=[
        {"title": "PHENOMENA", "items": [
            "Parent (1) JOIN Child (N) → hasil = 1 x N rows",
            "Untuk 1 author dengan 50 post → 50 rows yang sama untuk author",
            "Dengan multiple one-to-many → cartesian explosion paralel",
        ], "type": "bullet"},
        {"title": "IMPACT", "items": [
            "Network bandwidth: data yang dikirim jauh lebih besar",
            "Memory: result set di aplikasi jauh lebih besar",
            "CPU: parsing dan mapping hasil yang lebih lambat",
        ], "type": "bullet"},
        {"title": "SOLUSI ORM", "items": [
            "Django: prefetch_related() — gunakan query terpisah, join di Python",
            "EF Core: split query — eksekusi query terpisah per relasi",
            "SQLAlchemy: selectinload() — batch select, join di memori",
        ], "type": "bullet"},
    ],
    mental_model="JOIN-based eager loading causes cartesian explosion — prefer split/batch queries.",
    source_path="labs/22-n-plus-one-query-problem/content/05-key-takeaways.md",
    slide_num="15 / 24"
)

# Slide 16: Column Selection Optimization
add_slide(
    categories=["OPTIMASI LANJUTAN"],
    title="Strategi Alternatif: Column Selection",
    subtitle="select/pluck/values/withCount — ambil skalar, bukan entitas penuh",
    content_items=[
        {"title": "PERBANDINGAN", "items": [
            "Eager loading penuh: SELECT * FROM posts WHERE ... + mapping",
            "Column selection: SELECT title FROM posts WHERE ... + raw values",
            "Jika hanya butuh judul atau counter — column selection jauh lebih efisien",
        ], "type": "bullet"},
        {"title": "API DI BEREBERAPA ORM", "items": [
            "Django: values('title'), withCount('posts')",
            "Rails: pluck(:title), posts.count",
            "Laravel: select('title')->withCount('posts')",
            "EF Core: Select(p => p.Title), ToQueryString",
        ], "type": "bullet"},
        {"title": "MANFAAT", "items": [
            "Bandwidth: jauh lebih sedikit data dikirim",
            "Memory: tidak perlu materialisasi entitas penuh",
            "CPU: parsing lebih cepat — tidak perlu mapper ORM",
        ], "type": "bullet"},
    ],
    mental_model="For scalars and counts, column selection beats eager loading of full entities.",
    source_path="labs/22-n-plus-one-query-problem/content/05-key-takeaways.md",
    slide_num="16 / 24"
)

# Slide 17: Profiling-First Workflow
add_slide(
    categories=["ALUR KERJA"],
    title="Workflow Profiling-First",
    subtitle="Langkah demi langkah: ukur → identifikasi → optimasi",
    content_items=[
        {"title": "TIGA LANGKAH", "items": [
            "1. Ukur query count per request (bukan hanya durasi)",
            "2. Identifikasi pola N+1 (lazy loading di loop, per-object API call)",
            "3. Optimasi: batch eager loading, column selection, atau caching",
        ], "type": "bullet"},
        {"title": "TOOLING", "items": [
            "Django: django-nplusone, django-debug-toolbar",
            "Rails: bullet gem",
            "Laravel: beyondcode/laravel-query-detector",
            "EF Core: Microsoft.Extensions.Logging + query tag",
        ], "type": "bullet"},
        {"title": "CI/CD INTEGRASI", "items": [
            "Query count test assertion di CI pipeline",
            "Threshold maksimal query per endpoint",
            "Strict loading mode untuk mencegah lazy loading tak disadari",
        ], "type": "bullet"},
    ],
    mental_model="Measure first, optimize second — query count profiling before duration.",
    source_path="labs/22-n-plus-one-query-problem/content/05-key-takeaways.md",
    slide_num="17 / 24"
)

# Slide 18: Cross-ORM Comparison
add_slide(
    categories=["STANDAR INDUSTRI"],
    title="Eager Loading di Berbagai ORM",
    subtitle="API naming berbeda, konsep dasar sama",
    content_items=[
        {"title": "TABEL PERBANDINGAN", "items": [
            "ORM        | Method                  | Type        | SQL",
            "-----------|--------------------------|-------------|---------",
            "Django     | select_related()         | JOIN        | 1 query",
            "Django     | prefetch_related()       | Split       | 2 queries",
            "Rails      | includes()               | Auto-switch | 1/2",
            "Laravel    | with()                   | Split       | 2 queries",
            "EF Core    | Include()                | JOIN        | 1 query",
            "EF Core    | ThenInclude()            | JOIN        | 1 query",
        ], "type": "table"},
        {"title": "KUNCI PERBEDAAN", "items": [
            "select_related: JOIN — risiko cartesian explosion",
            "prefetch_related: split query — aman dari explosion",
            "includes: Rails memilih otomatis (smart)",
        ], "type": "bullet"},
        {"title": "KONEKSI KE LAB", "items": [
            "Lab ini: split query via GetPostsByAuthorIDs",
            "Mapping ini setara dengan prefetch_related / split query",
            "Bukan JOIN-based eager loading — menghindari cartesian explosion",
        ], "type": "bullet"},
    ],
    mental_model="ORMs agree on the concept — split/batch to avoid cartesian explosion, JOIN to minimize queries.",
    source_path="labs/22-n-plus-one-query-problem/content/02-master-draft.md",
    slide_num="18 / 24"
)

# Slide 19: Evaluation Table
add_slide(
    categories=["TABEL EVALUASI"],
    title="Evaluasi: N+1 vs Eager Loading",
    subtitle="Perbandingan kuantitatif dan kualitatif",
    content_items=[
        {"title": "TABEL PERBANDINGAN", "items": [
            "Dimensi            | N+1 (Lazy)        | Eager (Batch)",
            "-------------------|------------------|---------------",
            "Jumlah query       | N+1              | 2",
            "Skalabilitas       | Linear (O(N))    | Konstan (O(1))",
            "Memory overhead    | Minimal          | Map building",
            "Risk cartesian     | Tidak ada        | JOIN = ya (hindar)",
            "Readability        | Intuitive        | Perlu pemahaman",
            "Over-fetching      | Hanya per item   | Semua per satu",
        ], "type": "table"},
        {"title": "QUERY COUNT SUMMARY", "items": [
            "N=3:     N+1=4     | Eager=2",
            "N=100:   N+1=101   | Eager=2",
            "N=1000:  N+1=1001  | Eager=2",
            "N=10000: N+1=10001 | Eager=2",
        ], "type": "table"},
        {"title": "REKOMENDASI", "items": [
            "Gunakan eager loading untuk data yang pasti dibutuhkan",
            "Gunakan column selection untuk skalar/count",
            "Gunakan lazy loading hanya untuk data opsional/relasional",
        ], "type": "bullet"},
    ],
    mental_model="Trade intuition for correctness — eager batch wins on query count and stays O(1).",
    source_path="labs/22-n-plus-one-query-problem/internal/blog/repository.go",
    slide_num="19 / 24"
)

# Slide 20: Failure Scenario
add_slide(
    categories=["STUDI KASUS"],
    title="Studi Kasus: GET /api/work-orders",
    subtitle="Dari 6 query di dev menjadi 101 di prod — kisah nyata N+1",
    content_items=[
        {"title": "SCENARIO AWAL", "items": [
            "Endpoint: GET /api/work-orders",
            "Data uji dev: 5 work order, 1 item each = 6 query",
            "Response time: <10ms — terlihat cepat",
            "Lolos deployment ke production",
        ], "type": "bullet"},
        {"title": "SCENARIO PRODUKSI", "items": [
            "Production: 100 work order, 5 item each",
            "N+1: 1 + 100 = 101 query seril",
            "Response time: 2.4 detik (ilustratif)",
            "Connection pool exhausted — request pending",
        ], "type": "bullet"},
        {"title": "MITIGASI", "items":[
            "Batched query: WHERE work_order_id IN (...) → 2 query",
            "Column selection: hanya ambil kolom yang perlu",
            "Pagination: tidak load 100 work order sekaligus",
            "Query count assertion di test — mencegah regresi",
        ], "type": "bullet"},
    ],
    mental_model="Dev data hides prod pain — assert query counts for every realistic data size.",
    source_path="labs/22-n-plus-one-query-problem/content/02-master-draft.md",
    slide_num="20 / 24"
)

# Slide 21: Checklist
add_slide(
    categories=["PANDUAN OPERASIONAL"],
    title="Checklist Audit Performa Data Access",
    subtitle="9 poin verifikasi sebelum deploy ke production",
    content_items=[
        {"title": "CHECKLIST", "items": [
            "[ ] Profiling query count sebelum deploy",
            "[ ] Audit lazy loading dalam loop — ganti ke batch/eager",
            "[ ] Eager loading per-query, bukan per-mapping (hindari EAGER global)",
            "[ ] Cek cartesian explosion pada relasi one-to-many/many-to-many",
            "[ ] Gunakan column selection (pluck/values) untuk skalar",
            "[ ] Pagination untuk hasil besar",
            "[ ] Query count assertion di unit/integration test",
            "[ ] Pertimbangkan strict loading mode",
            "[ ] Query count test assertion di CI pipeline",
        ], "type": "bullet"},
        {"title": "ALAT BANTU", "items": [
            "django-nplusone — deteksi N+1 di Django",
            "bullet gem — deteksi N+1 di Rails",
            "laravel-query-detector — deteksi di Laravel",
            "EF Core logging — tracking query patterns",
        ], "type": "bullet"},
        {"title": "PRINSIP", "items": [
            "Jika query count meningkat linear dengan N → ada N+1",
            "Jika query count konstan → batching sudah tepat",
            "Query count assertion = automated guardrail",
        ], "type": "bullet"},
    ],
    mental_model="Nine checkpoints cover detection, batching, and prevention — assert query counts in CI.",
    source_path="labs/22-n-plus-one-query-problem/content/02-master-draft.md",
    slide_num="21 / 24"
)

# Slide 22: Lab Boundaries
add_slide(
    categories=["BATAS LAB"],
    title="Batasan Implementasi Lab 22",
    subtitle="Lab ini adalah simulasi konseptual — bukan produksi database asli",
    content_items=[
        {"title": "BATAS 1: IN-MEMORY MOCK STORE", "items": [
            "Store menggunakan slice in-memory, bukan database nyata",
            "Query counter hanya counter — tidak ada latency jaringan nyata",
            "Angka 4 query vs 2 query adalah verifikasi logik, bukan benchmark",
        ], "type": "bullet"},
        {"title": "BATAS 2: KONEKSI KE ORM", "items":[
            "Lab implementasi manual, bukan integrasi dengan Django/Rails/EF Core",
            "API batch (GetPostsByAuthorIDs) merepresentasikan konsep IN(...)",
            "ORM-level lazy vs eager loading tidak langsung ditunjukkan di sini",
        ], "type": "bullet"},
        {"title": "BATAS 3: ANGKA ILUSTRATIF", "items":[
            "Angka 2.4s pada studi kasus adalah angka ilustratif",
            "Dataset hanya 3 author — produksi bisa ribuan",
            "Lab membuktikan konsep query count, bukan performa absolut",
        ], "type": "bullet"},
    ],
    mental_model="This lab proves query-count reduction conceptually — real-world latency varies.",
    source_path="labs/22-n-plus-one-query-problem/content/01-content-brief.md",
    slide_num="22 / 24"
)

# Slide 23: Core Design Principles
add_slide(
    categories=["PRINSIP INTI"],
    title="Prinsip Desain Akses Data yang Benar",
    subtitle="Tiga prinsip dasar mencegah N+1 secara sistematis",
    content_items=[
        {"title": "PRINSIP 1: PROFILING- FIRST", "items":[
            "Hitung query count sebelum optimasi",
            "Tooling: ORM query logger, APM, slow query analysis",
            "CI/CD: query count threshold assertion",
        ], "type": "bullet"},
        {"title": "PRINSIP 2: BATCH, JANGAN LOOP", "items":[
            "Kumpulkan ID → satu query IN(...) → join di memori",
            "Pola ini universal: database, HTTP, GraphQL",
            "Query-level eager loading — kontrol per query",
        ], "type": "bullet"},
        {"title": "PRINSIP 3: PILIH KOLOM YANG DIPERLUKAN", "items":[
            "select/pluck/values/withCount — ambil skalar",
            "Hindari materialisasi entitas penuh jika tidak dipakai",
            "Reduksi bandwidth, memory, dan CPU parsing",
        ], "type": "bullet"},
    ],
    mental_model="Three principles: measure first, batch not loop, select only what you need.",
    source_path="labs/22-n-plus-one-query-problem/content/02-master-draft.md",
    slide_num="23 / 24"
)

# Slide 24: Closing (Dark Navy, rich visual)
add_slide(
    categories=["RINGKASAN"],
    title="Ringkasan Lab 22",
    subtitle=f"Commit: {COMMIT_SHA[:8]} | Test: 4/4 PASS | Repository: lukman-ss/software-engineering-lab",
    content_items=[
        {"title": "KEY RESULTS", "items":[
            "N+1: 1 + N query (linear dengan record count)",
            "Eager: 2 query (konstan, terlepas dari N)",
            "Test assertion: query count verified, deep equality confirmed",
        ], "type": "bullet"},
        {"title": "MENTAL MODEL", "items":[
            "N+1 is a pattern, not a layer — it strikes wherever I/O runs in a loop.",
            "Batch all related identifiers into a single IN(...) query.",
            "Assert query counts in unit tests — guardrails don't regress silently.",
        ], "type": "bullet"},
        {"title": "RESOURCES", "items":[
            f"Repository: {REPO_URL}",
            f"Commit SHA: {COMMIT_SHA}",
            "Run demo: go run ./cmd/demo",
            "Run tests: go test -v -race ./...",
        ], "type": "bullet"},
    ],
    mental_model="N+1 is a pattern, not a layer — batch all IDs into ONE query, assert count in tests.",
    source_path="labs/22-n-plus-one-query-problem/README.md",
    slide_num=""
)

if __name__ == "__main__":
    print(f"Total slides defined: {len(SLIDES)}")
