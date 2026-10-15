#!/usr/bin/env python3
"""
Publication Builder Definition for Software Engineering Lab 25: Rate Limiting & Backpressure.
Follows:
- cover-rules.md (Mandatory raster hero image at COVER_IMAGE_PATH, series identity lock, no micro-text)
- slide-rules.md (<=45 words prose, <=30 words on code slides, dominant visual anchors, 50% zoom test)
- quality-bar.md (sharp engineering insights, strictly bounded claims, >=12 slides with limiting phrases, no overclaims)
- article-rules.md (Dev.to 22-section mapping, Medium ID HTML 22-section mapping)
"""
import os
import sys

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/25-rate-limiting-and-backpressure/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/25-rate-limiting-and-backpressure"
COMMIT_SHA = "45cc429"
LAB_NO = "25"
SLUG = "rate-limiting-and-backpressure"

TITLE_ID = "Rate Limiting & Backpressure"
SUBTITLE_ID = "Mengontrol Laju Masuk dan Mempropagasi Tekanan Beban di Sistem Terdistribusi"
TITLE_EN = "Rate Limiting & Backpressure"
SUBTITLE_EN = "Controlling Inflow Rate and Propagating Overload Pressure in Distributed Systems"

AUTHOR = "Lukman / lukman-ss"
TOPIC_TAGS = "Rate Limiting • Backpressure • Go • Distributed Systems • System Design"
COVER_IMAGE_PATH = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/25-rate-limiting-and-backpressure/result/image-cover/cover_art.jpg"

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
            "Retro Terminal Displaying Rate Limiters & Bounded Queue Pressure Dynamics",
            "Token Bucket, Leaky Bucket, Fast Rejection Backpressure, and AWS Full Jitter in Pure Go",
            f"Penulis: {AUTHOR}"
        ]}
    ],
    mental_model="Rate limiting guards system boundaries; backpressure propagates overload backward through the internal pipeline.",
    source_path="",
    slide_num="01 / 24",
    layout="cover"
)

# Slide 02: Problem Nyata — Unbounded Queue & Cascading Failure (Layout: cards)
# Limiting phrase 1: "bukan berarti"
add_slide(
    categories=["PROBLEM", "QUEUE OVERLOAD"],
    title="Antrean Tak Berbatas Memicu Kolaps Memori dan Latensi",
    subtitle="Ketika laju kedatangan melebihi kapasitas proses hilir, antrean bertumpuk tanpa batas.",
    content_items=[
        {"type": "cards", "title": "DAMPAK KETIADAAN LIMITER DAN BACKPRESSURE", "items": [
            {"title": "Penumpukan Antrean", "desc": "5.000.000 request dengan drain 2.000/detik butuh 41 menit 40 detik untuk menguras buffer."},
            {"title": "Kehabisan Memori", "desc": "Buffer tak terbatas menyerap jutaan objek hingga proses terbunuh oleh OOM killer OS."},
            {"title": "Bukan Berarti Rusak", "desc": "Sistem melayani request normal saat sepi, tetapi bukan berarti aman dari lonjakan mendadak."}
        ]}
    ],
    narrative="Ketiadaan batas masuk membiarkan jutaan request menumpuk saat downstream melambat. Sistem yang berjalan lancar di kondisi normal bukan berarti mampu bertahan tanpa rate limiter ketika lonjakan beban traffic tak terduga datang.",
    narrative_en="Without boundary controls, arrival bursts overwhelm downstream consumers and trigger unbounded memory accumulation. A system running smoothly under light load does not automatically withstand unexpected traffic spikes.",
    mental_model="Unbounded buffers trade bounded latency for catastrophic out-of-memory crashes.",
    source_path="content/02-master-draft.md:3-12",
    slide_num="02 / 24",
    layout="cards"
)

# Slide 03: Mengapa Pendekatan Naif Gagal — IP Throttling & Jitterless Retry (Layout: cards)
# Limiting phrase 2: "tidak otomatis"
add_slide(
    categories=["WHY FAILS", "NAIVE APPROACH"],
    title="Throttling Berbasis IP dan Retry Deterministik Melumpuhkan Sistem",
    subtitle="CGNAT mencampur ribuan tenant, sementara retry tanpa jitter memicu thundering herd.",
    content_items=[
        {"type": "cards", "title": "KELEMAHAN PENDEKATAN NAIF", "items": [
            {"title": "Kolateral CGNAT", "desc": "RFC 6598 membuat ribuan pelanggan berbagi satu IP; satu pengguna bising melumpuhkan seluruh tenant."},
            {"title": "Thundering Herd", "desc": "100 client retry di interval deterministik yang sama menggandakan beban saat backend baru pulih."},
            {"title": "Tidak Otomatis Aman", "desc": "Menambah kapasitas server tidak otomatis menyelesaikan bottleneck jika retry storm terus menggempur."}
        ]}
    ],
    narrative="Rate limiting berbasis IP publik merugikan ribuan pengguna di balik CGNAT karena kuota terpakai habis oleh satu pengguna. Di sisi lain, menambah kapasitas CPU tidak otomatis menyelesaikan masalah jika retry klien terus datang serempak.",
    narrative_en="IP-based limiting collaterally penalizes legitimate tenants behind CGNAT gateways. Scaling compute capacity does not automatically resolve outages when clients retry simultaneously in lockstep synchronized waves.",
    mental_model="IP-based rate limits cause collateral damage under CGNAT; jitterless retries synchronize thundering herds.",
    source_path="content/02-master-draft.md:33-36",
    slide_num="03 / 24",
    layout="cards"
)

# Slide 04: Mental Model Inti — Dua Sisi Pertahanan (Layout: schema)
# Limiting phrase 3: "hanya pada jalur ini"
add_slide(
    categories=["MENTAL MODEL", "PILLARS"],
    title="Tiga Pilar Resiliensi: Limiter, Bounded Queue, dan Jitter",
    subtitle="Pengendalian pintu masuk, penolakan cepat internal, dan desinkronisasi percobaan ulang.",
    content_items=[
        {"type": "schema", "title": "TIGA PILAR PENGENDALIAN OVERLOAD", "items": [
            {"phase": "1. Rate Limiting", "role": "Boundary Guard", "desc": "Token bucket mengontrol laju masuk di API gateway dan memotong burst liar."},
            {"phase": "2. Backpressure", "role": "Internal Defense", "desc": "Bounded queue menolak kelebihan job seketika dengan sinyal ErrQueueFull."},
            {"phase": "3. Jittered Retry", "role": "Client Backoff", "desc": "Hanya pada jalur ini acakan [0, temp] mencegah sinkronisasi thundering herd."},
            {"phase": "4. RFC 6585 429", "role": "Protocol Signal", "desc": "Memberi tahu klien kapan waktu aman untuk mencoba kembali via Retry-After."}
        ]}
    ],
    narrative="Resiliensi sistem terdistribusi bertumpu pada isolasi pintu masuk via token bucket dan propagasi tekanan ke hulu via bounded queue. Hanya pada jalur ini klien yang menerima sinyal 429 mendistribusikan waktu tunggu secara acak merata.",
    narrative_en="Distributed resilience relies on entrance policing via token buckets combined with backward pressure signals from bounded buffers. Only along this path do throttled clients disperse retry intervals uniformly across time.",
    mental_model="Rate limiting polices external boundaries; backpressure rejects internal overload; jitter disperses client retries.",
    source_path="content/01-content-brief.md:6-10",
    slide_num="04 / 24",
    layout="schema"
)

# Slide 05: Alur Arsitektur Berlapis (Layout: schema)
# Limiting phrase 4: "belum membuktikan"
add_slide(
    categories=["ARCHITECTURE", "PIPELINE"],
    title="Alur Arsitektur: Dari Ekstraksi Tenant ke Konsumer",
    subtitle="Memisahkan layer gateway, rate limiter, bounded buffer, hingga worker pool.",
    content_items=[
        {"type": "schema", "title": "PIPELINE PEMROSESAN PER TENANT", "items": [
            {"phase": "1. Tenant Extractor", "role": "HTTP Middleware", "desc": "Membaca header X-API-Key untuk mendapatkan identitas tenant terisolasi."},
            {"phase": "2. Token Bucket", "role": "Registry Lookup", "desc": "Memeriksa kuota token; melebihi batas menghasilkan 429 Too Many Requests."},
            {"phase": "3. Bounded Queue", "role": "Channel Buffer", "desc": "Menerima job ke antrean; buffer penuh langsung menolak dengan ErrQueueFull."},
            {"phase": "4. Worker Pool", "role": "Job Consumer", "desc": "Konsumer memproses job; belum membuktikan isolasi resource CPU per worker."}
        ]}
    ],
    narrative="Request masuk diidentifikasi per tenant, divalidasi oleh token bucket, lalu diserahkan ke bounded queue. Arsitektur ini mencegah kelebihan beban di memori, namun belum membuktikan jaminan alokasi CPU per tenant di level container OS.",
    narrative_en="Incoming requests are identified per tenant, gated by token buckets, and buffered into bounded channels. This protects memory from unbounded growth, though it does not prove OS-level CPU core isolation across workers.",
    mental_model="Layered defense isolates tenants at the gateway and sheds excess load immediately at queue boundaries.",
    source_path="content/04-diagrams.md:7-24",
    slide_num="05 / 24",
    layout="schema"
)

# Slide 06: Algoritma Token Bucket vs Leaky Bucket (Layout: table)
add_slide(
    categories=["ALGORITHM", "COMPARISON"],
    title="Perbandingan Algoritma: Token Bucket vs Leaky Bucket",
    subtitle="Keduanya setara secara matematis sebagai meter, namun berbeda sudut pandang kendali.",
    content_items=[
        {"title": "MATRIKS KOMPARASI TOKEN BUCKET VS LEAKY BUCKET", "items": [
            "Karakteristik | Token Bucket (Traffic Burst) | Leaky Bucket (Traffic Smoothing)",
            "Toleransi Burst | Mengizinkan burst hingga kapasitas B | Menahan burst; output rata pada laju R",
            "State Utama | Jumlah token yang siap dikonsumsi | Level air yang menunggu dibocorkan",
            "Mekanisme Tolak | Token < n yang diminta | Penambahan air melebihi kapasitas B",
            "Kasus Penggunaan | API Gateway, bursty interactive user | Egress traffic shaping, NGINX limit_req",
            "Respons Standar | HTTP 429 Too Many Requests | HTTP 429 atau 503 Service Unavailable"
        ], "type": "table"}
    ],
    narrative="Token bucket mengizinkan akumulasi token untuk mengakomodasi burst jangka pendek, sedangkan leaky bucket memaksa laju keluaran konstan. Pilihan bergantung pada kebutuhan: fleksibilitas burst atau penghalusan trafik ke hilir.",
    narrative_en="Token buckets permit transient bursts up to capacity while enforcing average rates; leaky buckets strictly smooth outflow. Selection hinges on whether downstream services prefer accommodating bursts or demanding paced traffic.",
    mental_model="Token bucket permits bursts up to capacity; leaky bucket smooths output to a strict constant rate.",
    source_path="content/02-master-draft.md:21-28",
    slide_num="06 / 24",
    layout="table"
)

# Slide 07: Cuplikan Kode Token Bucket: Replenish & Debit (Layout: code)
# Limiting phrase 5: "bukan berarti"
add_slide(
    categories=["IMPLEMENTATION", "TOKEN BUCKET"],
    title="Token Bucket: Pengisian Waktu Berjalan dan Konsumsi",
    subtitle="Menghitung elapsed time secara monotonik tanpa timer background yang boros.",
    content_items=[
        {"type": "code", "title": "internal/ratelimit/bucket.go", "items": [
            "func (tb *TokenBucket) AllowN(n float64) bool {",
            "    tb.mu.Lock()",
            "    defer tb.mu.Unlock()",
            "    now := time.Now()",
            "    elapsed := now.Sub(tb.lastRefill).Seconds()",
            "    tb.tokens = math.Min(tb.capacity, tb.tokens+elapsed*tb.refillRate)",
            "    tb.lastRefill = now",
            "    if tb.tokens >= n {",
            "        tb.tokens -= n",
            "        return true",
            "    }",
            "    return false",
            "}"
        ]}
    ],
    narrative="TokenBucket menghitung penambahan token secara pasif berdasarkan selisih waktu antar-panggilan. Penggunaan mutex bukan berarti implementasi tidak optimal; lock lokal menjamin atomisitas pembaruan multi-field secara bersih dan aman.",
    narrative_en="TokenBucket lazily replenishes tokens on each invocation via monotonic timestamps, avoiding background timer routines. Relying on mutex locks does not imply inefficient code; it ensures robust multi-field consistency.",
    mental_model="Lazy replenishment on call converts expensive timer ticker goroutines into cheap timestamp arithmetic.",
    source_path="internal/ratelimit/bucket.go:29-46",
    slide_num="07 / 24",
    layout="code"
)

# Slide 08: Cuplikan Kode Retry-After Calculation (Layout: code)
# Limiting phrase 6: "hanya pada jalur ini"
add_slide(
    categories=["IMPLEMENTATION", "RFC 6585"],
    title="Kalkulasi Retry-After: Panduan Presisi Waktu Tunggu",
    subtitle="Menghitung waktu tunggu pembulatan ke atas agar klien tidak retry prematur.",
    content_items=[
        {"type": "code", "title": "internal/ratelimit/bucket.go", "items": [
            "func (tb *TokenBucket) RetryAfterSeconds(n float64) int {",
            "    tb.mu.Lock()",
            "    defer tb.mu.Unlock()",
            "    // ... perbarui token dari elapsed time ...",
            "    if tb.tokens >= n { return 0 }",
            "    needed := n - tb.tokens",
            "    secs := needed / tb.refillRate",
            "    rounded := int(secs)",
            "    if float64(rounded) < secs { rounded++ }",
            "    return rounded",
            "}"
        ]}
    ],
    narrative="Kekurangan token dibagi dengan laju pengisian lalu dibulatkan ke atas menjadi detik bulat. Hanya pada jalur ini klien yang ditolak menerima informasi presisi kapan kuota mereka akan terisi kembali tanpa menebak-nebak.",
    narrative_en="Token shortfall is divided by refill rate and rounded up to integer seconds. Only along this path do rejected callers receive accurate resumption timestamps rather than resorting to arbitrary retry delays.",
    mental_model="Round up retry delay to integer seconds to prevent clients from retrying before token replenishment.",
    source_path="internal/ratelimit/bucket.go:55-79",
    slide_num="08 / 24",
    layout="code"
)

# Slide 09: Cuplikan Kode Bounded Queue Fast Rejection (Layout: code)
# Limiting phrase 7: "bukan berarti"
add_slide(
    categories=["IMPLEMENTATION", "BACKPRESSURE"],
    title="Bounded Queue: Penolakan Cepat dengan Select-Default",
    subtitle="Mencegah tumpukan memori tak terbatas dengan membuang beban berlebih seketika.",
    content_items=[
        {"type": "code", "title": "internal/backpressure/queue.go", "items": [
            "func (bq *BoundedQueue) TrySubmit(job Job) error {",
            "    bq.stopMu.RLock()",
            "    defer bq.stopMu.RUnlock()",
            "    if bq.stopped { return ErrQueueStopped }",
            "    select {",
            "    case bq.queue <- job:",
            "        bq.accepted.Add(1)",
            "        return nil",
            "    default:",
            "        bq.rejected.Add(1)",
            "        return ErrQueueFull",
            "    }",
            "}"
        ]}
    ],
    narrative="Pola select-default mengantrekan job jika kanal longgar atau langsung mengembalikan ErrQueueFull saat penuh. Menolak job seketika bukan berarti sistem tidak andal; ini tindakan proteksi mandiri agar proses tidak kehabisan memori.",
    narrative_en="The select-default pattern pushes jobs if buffer capacity exists or drops instantly with ErrQueueFull. Fast rejection does not signify poor service quality; it is an essential survival mechanism preventing memory exhaustion.",
    mental_model="Select-default on bounded channels gives zero-allocation, instant backpressure shedding.",
    source_path="internal/backpressure/queue.go:67-85",
    slide_num="09 / 24",
    layout="code"
)

# Slide 10: Cuplikan Kode AWS Full Jitter Retry (Layout: code)
# Limiting phrase 8: "tidak otomatis"
add_slide(
    categories=["IMPLEMENTATION", "FULL JITTER"],
    title="AWS Full Jitter: Mencegah Sinkronisasi Thundering Herd",
    subtitle="sleep = random(0, min(cap, base × 2^attempt)) menyebarkan lonjakan beban.",
    content_items=[
        {"type": "code", "title": "internal/retry/backoff.go", "items": [
            "func ComputeBackoff(strategy BackoffStrategy, attempt int,",
            "    cfg Config, prevSleep time.Duration) time.Duration {",
            "    expBackoff := float64(cfg.Base) * math.Pow(2, float64(attempt))",
            "    temp := math.Min(float64(cfg.Cap), expBackoff)",
            "    switch strategy {",
            "    case FullJitter:",
            "        if temp <= 0 { return 0 }",
            "        return time.Duration(rand.Float64() * temp)",
            "    // ... case EqualJitter, DecorrelatedJitter, NoJitter ...",
            "    }",
            "}"
        ]}
    ],
    narrative="Full Jitter mengambil sampel acak seragam antara nol hingga batas plafon eksponensial. Menambahkan jitter tidak otomatis mempercepat waktu pemulihan jika kapasitas downstream memang telah runtuh secara permanen.",
    narrative_en="Full Jitter samples uniformly across zero to the exponential ceiling, desynchronizing concurrent clients. Introducing jitter does not automatically restore services if downstream capacity has experienced permanent collapse.",
    mental_model="Full jitter samples uniformly from zero to the exponential ceiling to disperse client retry waves.",
    source_path="internal/retry/backoff.go:24-43",
    slide_num="10 / 24",
    layout="code"
)

# Slide 11: Cuplikan Kode HTTP Middleware 429 & Retry-After (Layout: code)
# Limiting phrase 9: "bukan berarti"
add_slide(
    categories=["IMPLEMENTATION", "MIDDLEWARE"],
    title="Middleware HTTP: Menegakkan Status 429 RFC 6585",
    subtitle="Mengekstraksi tenant API key dan mengembalikan header Retry-After standar.",
    content_items=[
        {"type": "code", "title": "internal/httputil/middleware.go", "items": [
            "func RateLimitMiddleware(reg *ratelimit.Registry, next http.Handler) http.Handler {",
            "    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {",
            "        key := r.Header.Get(\"X-API-Key\")",
            "        if key == \"\" { key = \"anonymous\" }",
            "        bucket := reg.Get(key)",
            "        if !bucket.Allow() {",
            "            retryAfter := bucket.RetryAfterSeconds(1.0)",
            "            w.Header().Set(\"Retry-After\", strconv.Itoa(retryAfter))",
            "            w.WriteHeader(http.StatusTooManyRequests) // 429",
            "            _ = json.NewEncoder(w).Encode(RateLimitResponse{...})",
            "            return",
            "        }",
            "        next.ServeHTTP(w, r)",
            "    })",
            "}"
        ]}
    ],
    narrative="Middleware mengisolasi token bucket per tenant key sebelum meneruskan request ke handler berikutnya. Mengirimkan status 429 bukan berarti client dilarang selamanya; header Retry-After memandu kapan pemanggilan boleh diulang.",
    narrative_en="Middleware validates per-tenant token buckets prior to passing execution to application handlers. Returning HTTP 429 does not imply permanent rejection; the Retry-After header guides intelligent rescheduling.",
    mental_model="Extract tenant keys at edge middleware, reject excess with 429, and provide explicit Retry-After headers.",
    source_path="internal/httputil/middleware.go:17-39",
    slide_num="11 / 24",
    layout="code"
)

# Slide 12: Eksekusi Demo Mandiri (Layout: code)
add_slide(
    categories=["EVIDENCE", "TERMINAL DEMO"],
    title="Eksekusi Demo cmd/demo: Empat Verifikasi Perilaku",
    subtitle="Output terminal membuktikan burst, smoothing, penolakan beban, dan sebaran jitter.",
    content_items=[
        {"type": "code", "title": "TERMINAL OUTPUT: go run ./cmd/demo", "items": [
            "=== 1. Token Bucket Burst & Rate Limiting ===",
            "Request #1..3: Allowed=true | Request #4..5: Allowed=false",
            "After 300ms pause: Allowed=true (Remaining Tokens: 0.5)",
            "=== 2. Leaky Bucket Traffic Smoothing ===",
            "Request #1..3: Allowed=true | Request #4..5: Allowed=false (Water: 3.0)",
            "=== 3. Bounded Queue Backpressure (Load Shedding) ===",
            "Job #1..3: ACCEPTED | Job #4, #6: REJECTED (queue capacity exceeded)",
            "Stats: Accepted=4, Rejected=2, Processed=1",
            "=== 4. AWS Retry Backoff Strategies (Attempts 0..3) ===",
            "Attempt 0 -> NoJitter: 100ms  | FullJitter: 61ms   | EqualJitter: 80ms",
            "Attempt 3 -> NoJitter: 800ms  | FullJitter: 342ms  | EqualJitter: 647ms"
        ]}
    ],
    narrative="Demo membuktikan seluruh skenario berjalan sesuai spesifikasi: 3 request awal lolos burst lalu ditolak, buffer antrean menolak job saat melampaui kapasitas 3, dan Full Jitter mendispersi waktu jeda secara acak.",
    narrative_en="The demo run verifies all four components: 3 initial burst requests succeed before immediate rejection, the bounded buffer drops excess jobs beyond capacity 3, and Full Jitter desynchronizes retry sleeps.",
    mental_model="Executable verification proves burst accommodation, queue shedding, and jitter dispersion in live runs.",
    source_path="cmd/demo/main.go:13-65",
    slide_num="12 / 24",
    layout="code"
)

# Slide 13: Verifikasi Test Suite & Race Detector (Layout: code)
# Limiting phrase 10: "simulasi ini tidak sama dengan"
add_slide(
    categories=["VERIFICATION", "TEST SUITE"],
    title="Verifikasi Test Suite: 15/15 PASS dan Bersih dari Race",
    subtitle="Seluruh modul rate limiting, queue, retry, dan middleware teruji konkurensi.",
    content_items=[
        {"type": "code", "title": "TERMINAL OUTPUT: go test -v -race ./...", "items": [
            "=== RUN   TestTokenBucket_BurstAndRefill / ConcurrencyRace",
            "--- PASS: TestTokenBucket_* (0.05s)",
            "=== RUN   TestLeakyBucket_LeakRate / ConcurrencyRace",
            "--- PASS: TestLeakyBucket_* (0.02s)",
            "=== RUN   TestRegistry_TenantIsolation / ConcurrentSameKeyGet",
            "--- PASS: TestRegistry_* (0.01s)",
            "=== RUN   TestBoundedQueue_RejectionUnderLoad / ConcurrencySafety",
            "--- PASS: TestBoundedQueue_* (0.08s)",
            "=== RUN   TestComputeBackoff_Bounds / DecorrelatedJitter_Bounds",
            "--- PASS: TestComputeBackoff_* (0.00s)",
            "=== RUN   TestRateLimitMiddleware_RFC6585 / AnonymousFallback",
            "--- PASS: TestRateLimitMiddleware_* (0.00s)",
            "PASS: ok labs/25-rate-limiting-and-backpressure (race detector: clean)"
        ]}
    ],
    narrative="Seluruh 15 unit test lolos tanpa kegagalan di bawah race detector Go dengan goroutine agresif. Simulasi in-memory ini tidak sama dengan cluster Redis terdistribusi yang harus menangani latency jaringan antar node.",
    narrative_en="All 15 automated test cases pass cleanly under the Go race detector under concurrent stress. While this verifies in-memory algorithms, this simulation is not equivalent to distributed Redis clusters spanning networks.",
    mental_model="Race-clean unit tests prove concurrent synchronization correctness in local shared-memory runtimes.",
    source_path="engineering/03-execution-result.md:8-60",
    slide_num="13 / 24",
    layout="code"
)

# Slide 14: Batasan Bukti Lab (Layout: cards)
# Limiting phrase 11: "belum membuktikan"
add_slide(
    categories=["AUDIT", "EVIDENCE BOUNDARY"],
    title="Batasan Bukti Lab: Yang Dibuktikan vs Belum",
    subtitle="Transparansi batasan rekayasa antara simulator lab lokal dan sistem produksi.",
    content_items=[
        {"type": "cards", "title": "BATAS PEMBUKTIAN REKAYASA", "items": [
            {"title": "YANG DIBUKTIKAN", "desc": "Penolakan burst token bucket, fast drop bounded queue, dan sebaran rentang Full Jitter."},
            {"title": "YANG TIDAK DIBUKTIKAN", "desc": "Sinkronisasi token lintas multi-instance pod tanpa shared state eksternal."},
            {"title": "BATAS KLAIM EMPIRIS", "desc": "Lab belum membuktikan ketahanan cluster terhadap badai DDoS ratusan ribu node eksternal."}
        ]}
    ],
    narrative="Lab ini membuktikan kebenaran matematis algoritma bucket, isolasi tenant lokal, dan pemotongan beban antrean. Namun, pengujian lokal ini belum membuktikan ketahanan cluster multi-region terhadap serangan DDoS skala masif.",
    narrative_en="This lab proves mathematical correctness for bucket algorithms, tenant isolation, and bounded queue shedding. However, it has not yet proven distributed resilience against massive multi-region distributed denial-of-service attacks.",
    mental_model="Clearly delimit verified algorithmic invariants from unproven distributed scale heuristics.",
    source_path="content/01-content-brief.md:11-12",
    slide_num="14 / 24",
    layout="cards"
)

# Slide 15: Trade-off Arsitektural (Layout: cards)
# Limiting phrase 12: "tidak otomatis"
add_slide(
    categories=["TRADE-OFF", "ENGINEERING COSTS"],
    title="Trade-off Arsitektur Rate Limiting & Backpressure",
    subtitle="Kompromi antara kepuasan pengguna, akurasi sinkronisasi, dan kompleksitas sistem.",
    content_items=[
        {"type": "cards", "title": "KOMPROMI REKAYASA SISTEM", "items": [
            {"title": "False Rejections", "desc": "Limiter agresif berisiko menolak pengguna legitimate saat aktivitas bisnis melonjak wajar."},
            {"title": "Shared State Latency", "desc": "Sinkronisasi Redis menambah latensi 1-2ms di hot-path setiap request HTTP yang masuk."},
            {"title": "Tidak Otomatis Kebal", "desc": "Memasang middleware limiter tidak otomatis membuat database kebal terhadap query N+1 berat."}
        ]}
    ],
    narrative="Rate limiting mencegah kehabisan kapasitas tetapi berisiko menolak transaksi sah pengguna di saat peak season. Menambahkan limiter di pintu gerbang tidak otomatis melindungi storage jika kueri internal tidak terindeks dengan baik.",
    narrative_en="Rate limiters prevent capacity exhaustion but introduce false rejection risks for legitimate users during peak traffic. Deploying gateway limiters does not automatically safeguard storage engines against unindexed slow queries.",
    mental_model="Rate limiting trades edge false-rejection risks for core service survival and latency stability.",
    source_path="content/02-master-draft.md:99-108",
    slide_num="15 / 24",
    layout="cards"
)

# Slide 16: Analisis Matematika: Little's Law & Backlog Growth (Layout: cards)
# Limiting phrase 13: "hanya pada jalur ini"
add_slide(
    categories=["ANALYSIS", "LITTLE'S LAW"],
    title="Analisis Beban: Little's Law dan Pertumbuhan Backlog",
    subtitle="L = λ × W membuktikan bagaimana penumpukan antrean merusak latensi sistem.",
    content_items=[
        {"type": "cards", "title": "FORMULA MATEMATIKA SISTEM ANTREAN", "items": [
            {"title": "Little's Law (L = λW)", "desc": "Kapasitas in-flight (L) adalah perkalian antara laju kedatangan (λ) dan durasi proses (W)."},
            {"title": "Akumulasi Backlog", "desc": "Pertumbuhan antrean = (Arrival Rate − Processing Rate) × Waktu durasi kelebihan beban."},
            {"title": "Kalkulasi Lab 25", "desc": "5.000.000 request ÷ 2.000 drain/detik = 2.500 detik (41 menit 40 detik); hanya pada jalur bounded buffer batas latensi terjaga."}
        ]}
    ],
    narrative="Jika laju kedatangan melampaui laju pemrosesan, waktu tunggu membengkak hingga puluhan menit. Hanya pada jalur dengan bounded queue penumpukan dipotong seketika sehingga latensi job yang diterima tetap dalam batas SLA.",
    narrative_en="When arrival rates outpace processing throughput, wait times explode exponentially into tens of minutes. Only along bounded queue pathways is backlog accumulation arrested, keeping accepted job latency within SLA targets.",
    mental_model="L = lambda * W: when processing capacity is saturated, queue wait times degrade proportionally to backlog depth.",
    source_path="content/02-master-draft.md:29-32",
    slide_num="16 / 24",
    layout="cards"
)

# Slide 17: Studi Kasus Produksi: Stripe's 4-Tier Rate Limiter (Layout: schema)
# Limiting phrase 14: "bukan berarti"
add_slide(
    categories=["CASE STUDY", "STRIPE ARCHITECTURE"],
    title="Studi Kasus Stripe: Arsitektur Limiter 4 Lapisan",
    subtitle="Pola pembagian tanggung jawab proteksi API gateway pada skala pembayaran global.",
    content_items=[
        {"type": "schema", "title": "EMPAT LAPISAN LIMITER STRIPE", "items": [
            {"phase": "1. Request Rate Limiter", "role": "Per-User Quota", "desc": "Membatasi laju request per pengguna menggunakan token bucket terdistribusi."},
            {"phase": "2. Concurrent Limiter", "role": "In-Flight Bound", "desc": "Membatasi jumlah request simultan yang sedang aktif berjalan per pengguna."},
            {"phase": "3. Fleet Load Shedder", "role": "Criticality Reserve", "desc": "Menolak trafik berprioritas rendah demi menjamin kapasitas transaksi kritis."},
            {"phase": "4. Worker Load Shedder", "role": "Execution Guard", "desc": "Membuang uji coba test-mode, GET lalu POST; bukan berarti request hilang tanpa status."}
        ]}
    ],
    narrative="Stripe menerapkan 4 tier limiter independen untuk mendegradasi sistem secara anggun saat kelebihan beban. Membuang trafik berprioritas rendah saat CPU tinggi bukan berarti sistem gagal, melainkan langkah penyelamatan transaksi keuangan kritis.",
    narrative_en="Stripe deploys four distinct limiter tiers to achieve graceful degradation under extreme load. Shedding low-criticality requests during CPU saturation does not mean system failure; it deliberately safeguards mission-critical payment workflows.",
    mental_model="Layer multiple specialized limiters to drop low-criticality traffic first during severe overload.",
    source_path="content/02-master-draft.md:117-124",
    slide_num="17 / 24",
    layout="schema"
)

# Slide 18: Studi Kasus Google SRE: Criticality & Retry Budget (Layout: cards)
# Limiting phrase 15: "cakupannya mengikuti test"
add_slide(
    categories=["CASE STUDY", "GOOGLE SRE"],
    title="Google SRE: Criticality Level dan Kuota Retry",
    subtitle="Mencegah retry storm cascading failure melalui pembatasan rasio percobaan ulang.",
    content_items=[
        {"type": "cards", "title": "DISIPLIN KETAHANAN GOOGLE SRE", "items": [
            {"title": "4 Level Criticality", "desc": "CRITICAL_PLUS, CRITICAL, SHEDDABLE_PLUS, SHEDDABLE; tolak level rendah lebih dulu."},
            {"title": "Retry Budget Klien", "desc": "Maksimal 3 attempt per-request dan rasio retry maksimal 10% dari total trafik klien."},
            {"title": "Token Bucket Quota", "desc": "AWS SDK mengadopsi 500 token retry pool; cakupannya mengikuti test kalibrasi error transien."}
        ]}
    ],
    narrative="Google SRE membagi request ke 4 tingkat kepentingan dan membatasi retry maksimal 10% dari volume panggilan. Strategi kuota retry token bucket ini diadopsi luas di industri, dan cakupannya mengikuti test evaluasi latensi backend.",
    narrative_en="Google SRE categorizes RPC traffic into four criticality tiers and enforces a strict 10% retry budget per client. Managing retries via token quotas prevents cascade amplification, with parameters calibrated to downstream recovery profiles.",
    mental_model="Enforce strict retry budgets (max 3 attempts, <10% retry ratio) to eliminate cascading retry storms.",
    source_path="content/02-master-draft.md:121-124",
    slide_num="18 / 24",
    layout="cards"
)

# Slide 19: Perbandingan Before vs After (Layout: cards)
# Limiting phrase 16: "bukan berarti"
add_slide(
    categories=["COMPARISON", "BEFORE VS AFTER"],
    title="Transformasi Desain: Unprotected vs Resilient Pipeline",
    subtitle="Pergeseran dari sistem rentan crash ke arsitektur dengan degradasi terukur.",
    content_items=[
        {"type": "cards", "title": "SEBELUM (UNPROTECTED) VS SESUDAH (RESILIENT)", "items": [
            {"title": "Sebelum: Unbounded Channel", "desc": "Antrean bertambah hingga jutaan item; latensi membengkak 40+ menit sebelum proses OOM crash."},
            {"title": "Sesudah: Bounded + Jitter", "desc": "Kelebihan beban ditolak cepat (429/ErrQueueFull); klien desinkronisasi via Full Jitter."},
            {"title": "Bukan Berarti Nol Error", "desc": "Resiliensi bukan berarti bebas error 100%, melainkan memastikan sistem tetap hidup melayani trafik sehat."}
        ]}
    ],
    narrative="Sebelum proteksi, lonjakan trafik menelan seluruh memori dan mematikan server. Menerapkan limiter dan backpressure bukan berarti meniadakan error, melainkan mengendalikan error secara sadar agar layanan inti tetap beroperasi.",
    narrative_en="Prior to rate limiting, traffic surges exhausted heap memory and triggered catastrophic process crashes. Building resilient boundaries does not eliminate errors; it trades uncontrolled crashes for controlled, bounded degradation.",
    mental_model="Transform uncontrolled systemic collapse into bounded, graceful rejection under peak overload.",
    source_path="content/02-master-draft.md:29-36",
    slide_num="19 / 24",
    layout="cards"
)

# Slide 20: Checklist Kesiapan Produksi (Layout: checklist)
add_slide(
    categories=["CHECKLIST", "PRODUCTION READINESS"],
    title="Checklist Rekayasa Rate Limiting & Backpressure",
    subtitle="Sepuluh kriteria wajib verifikasi sebelum merilis layanan ke traffic produksi.",
    content_items=[
        {"type": "checklist", "title": "KRITERIA IMPLEMENTASI RESILIENSI", "items": [
            "Identifikasi batas sistem: tentukan rate limiter di edge API gateway dan service mesh.",
            "Pilih algoritma yang tepat: Token Bucket untuk burst wajar, Leaky Bucket untuk output konstan.",
            "Gunakan kunci identitas per-tenant (X-API-Key atau User ID) untuk mencegah kolateral CGNAT.",
            "Terapkan Bounded Queue dengan kapasitas terukur sesuai laju pemrosesan konsumer downstream.",
            "Kembalikan respons standar HTTP 429 dengan header Retry-After presisi (RFC 6585).",
            "Gunakan Exponential Backoff dengan Full Jitter pada semua klien pemanggil untuk mencegah thundering herd.",
            "Tetapkan Retry Budget: maksimal 3 percobaan per request dan batas rasio retry di bawah 10%.",
            "Pantau Queue Age, bukan hanya Queue Length, untuk deteksi dini pembengkakan latensi pemrosesan."
        ]}
    ],
    narrative="Checklist ini merangkum disiplin perancangan sistem terdistribusi: isolasi batas, pilihan algoritma, identifikasi tenant, buffer terbatas, respons standar 429, desinkronisasi retry, kuota percobaan, dan observabilitas usia antrean.",
    narrative_en="This checklist enforces robust distributed system discipline: boundary limits, algorithm selection, tenant keying, bounded channels, standard 429 headers, jittered retries, retry budgets, and queue age monitoring.",
    mental_model="Apply strict architectural gates across limits, buffers, retries, and metrics before production rollout.",
    source_path="content/02-master-draft.md:126-138",
    slide_num="20 / 24",
    layout="checklist"
)

# Slide 21: Anti-Pattern Umum & Jebakan Desain (Layout: cards)
# Limiting phrase 17: "bukan berarti"
add_slide(
    categories=["ANTI-PATTERN", "PITFALLS"],
    title="Anti-Pattern Umum: Kesalahan Fatal Pengendalian Overload",
    subtitle="Kekeliruan perancangan yang sering memicu pemadaman massal di industri.",
    content_items=[
        {"type": "cards", "title": "HINDARI KESALAHAN ARSITEKTURAL INI", "items": [
            {"title": "IP Limiting di CGNAT", "desc": "Membatasi IP publik memblokir ribuan pengguna sah di balik satu gerbang ISP seluler."},
            {"title": "Retry Tanpa Jitter", "desc": "Interval eksponensial kaku menyinkronkan ribuan klien menjadi gelombang thundering herd penghancur."},
            {"title": "Unbounded Buffer", "desc": "Antrean tanpa batas bukan berarti aman; itu bom waktu memori yang menjamin kegagalan total."}
        ]}
    ],
    narrative="Throttling berbasis IP merugikan tenant di balik CGNAT, sedangkan retry deterministik melipatgandakan beban downstream saat recovery. Memberi buffer antrean tak terbatas bukan berarti sistem toleran, melainkan menimbun kebangkrutan memori.",
    narrative_en="IP-based limiting penalizes innocent users sharing carrier gateways, while deterministic retries generate destructive thundering herds. Unbounded buffers do not reflect high availability; they guarantee out-of-memory crashes.",
    mental_model="Avoid IP-only throttling, eliminate deterministic retries, and never allow unbounded queues.",
    source_path="content/02-master-draft.md:109-116",
    slide_num="21 / 24",
    layout="cards"
)

# Slide 22: Reproduksi Mandiri & Command Runner (Layout: code)
# Limiting phrase 18: "production masih perlu"
add_slide(
    categories=["REPRODUCTION", "COMMANDS"],
    title="Reproduksi Mandiri di Mesin Lokal Menggunakan Go",
    subtitle="Validasi penuh tanpa dependensi eksternal pihak ketiga (pure Go standard library).",
    content_items=[
        {"type": "code", "title": "COMMAND RUNNER (PURE GO STANDARD LIBRARY)", "items": [
            "cd labs/25-rate-limiting-and-backpressure",
            "go test -v ./...    # validasi unit test bucket, queue, retry, dan middleware",
            "go test -race ./... # verifikasi konkurensi stateless bebas data race",
            "go run ./cmd/demo   # eksekusi demonstrasi langsung 4 pilar arsitektur"
        ]}
    ],
    narrative="Seluruh mekanisme rate limiting dan backpressure ditulis murni menggunakan Go standard library tanpa framework pihak ketiga. Pengujian lokal ini memverifikasi algoritma, namun production masih perlu Redis untuk koordinasi multi-pod.",
    narrative_en="The entire test harness relies strictly on pure Go standard library packages with zero external dependencies. While local testing confirms algorithmic correctness, production deployments still require Redis for multi-instance state sync.",
    mental_model="Zero-dependency local verification guarantees fast, repeatable validation across developer workflows.",
    source_path="README.md:55-70",
    slide_num="22 / 24",
    layout="code"
)

# Slide 23: Prinsip Arsitektural Jangka Panjang (Layout: cards)
# Limiting phrase 19: "tidak otomatis"
add_slide(
    categories=["TAKEAWAY", "CORE PRINCIPLES"],
    title="Prinsip Inti: Kendalikan Pintu Masuk, Propagasikan Tekanan",
    subtitle="Arsitektur bertahan lama ketika batas beban ditegakkan secara disiplin di setiap lapisan.",
    content_items=[
        {"type": "cards", "title": "TIGA PILAR KELANGSUNGAN ARSITEKTUR", "items": [
            {"title": "Polisikan Pintu Gerbang", "desc": "Gunakan token bucket per tenant untuk membatasi kuota sebelum request membebani resource inti."},
            {"title": "Tolak Cepat di Hilir", "desc": "Terapkan bounded queue agar kelebihan beban dibuang seketika tanpa menumpuk latensi."},
            {"title": "Desinkronisasi Klien", "desc": "Terapkan Full Jitter; tidak otomatis menjamin backend pulih jika downstream rusak permanen."}
        ]}
    ],
    narrative="Rate limiting menjaga gerbang masuk, sementara backpressure mempropagasikan sinyal tekanan balik ke hulu. Menerapkan Full Jitter pada klien tidak otomatis memulihkan sistem jika kapasitas downstream memang tidak memadai untuk beban dasar.",
    narrative_en="Rate limiting protects ingress boundaries while backpressure propagates overload signals backward. Adding full jitter to client retries does not automatically compensate for fundamentally undersized downstream infrastructure.",
    mental_model="Police external ingress, reject internal overload fast, and desynchronize client retries with jitter.",
    source_path="content/05-key-takeaways.md",
    slide_num="23 / 24",
    layout="cards"
)

# Slide 24: Closing Slide (Layout: closing)
add_slide(
    categories=["CLOSING"],
    title="Sistem tanpa backpressure adalah sistem yang menunda kebangkrutan memorinya sendiri.",
    subtitle="",
    content_items=[
        {"type": "bullet", "title": "MODUL IMPLEMENTASI", "items": [
            "internal/ratelimit/bucket.go • internal/ratelimit/registry.go",
            "internal/backpressure/queue.go • internal/retry/backoff.go",
            "internal/httputil/middleware.go • cmd/demo/main.go"
        ]},
        {"type": "bullet", "title": "VERIFIKASI & STATUS", "items": [
            f"Commit: {COMMIT_SHA[:7]} • Status: APPROVED",
            "Test passing: 15 unit tests (burst, refill, tenant isolation, queue shedding, jitter)",
            "Race detector: Clean (go test -race ./... PASS in 0.28s)",
            "Implementasi Go standard library murni; in-memory single process architecture."
        ]}
    ],
    mental_model="A system without backpressure is merely borrowing time until an unavoidable memory collapse.",
    source_path=REPO_URL,
    slide_num="24 / 24",
    layout="closing"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
    for s in SLIDES:
        print(f"[{s['slide_num']}] {s['layout']}: {s['title']}")
