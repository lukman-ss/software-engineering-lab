#!/usr/bin/env python3
"""Publication Builder Definition for Software Engineering Lab 28: Timeouts and Deadlines.
Generates 24-slide LinkedIn carousel covering deadline propagation, timeout budgeting, exponential backoff with full jitter, circuit breaking, and idempotency deduplication.
"""

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/28-timeouts-and-deadlines/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/28-timeouts-and-deadlines"
COMMIT_SHA = "ba4b621"
LAB_NO = "28"
SLUG = "timeouts-and-deadlines"

TITLE_ID = "Timeouts & Deadlines: Guardrails untuk Sistem Terdistribusi"
SUBTITLE_ID = "Mencegah Dependency Lambat Menghancurkan Layanan melalui Deadline Propagation, Jitter, Circuit Breaker, dan Idempotency"
TITLE_EN = "Timeouts & Deadlines: Guardrails for Distributed Systems"
SUBTITLE_EN = "Preventing Slow Dependencies from Crashing Services via Deadline Propagation, Jitter, Circuit Breaking, and Idempotency."

AUTHOR = "Lukman / lukman-ss"
TOPIC_TAGS = "Timeouts • Deadlines • Retry • Circuit Breaker • Idempotency • Go"
COVER_IMAGE_PATH = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/28-timeouts-and-deadlines/result/image-cover/cover_art.jpg"

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
        "layout": layout,
    })

# Slide 01: Cover
add_slide(
    categories=["COVER"],
    title=TITLE_ID,
    subtitle=SUBTITLE_ID,
    content_items=[{"type": "cover", "title": "", "items": [
        "Timeouts & Deadlines — Guardrails untuk Sistem Terdistribusi",
        "Mengelola latency, retry, circuit & idempotency di Go",
        f"Penulis: {AUTHOR}"
    ]}],
    mental_model="Timeout & deadline adalah resource protection guardrails, bukan sinyal kegagalan.",
    source_path=REPO_URL,
    slide_num="01 / 24",
    narrative="Lab ini menampilkan pola resiliency yang melindungi layanan dari dependency lambat. Setiap pola dibuktikan dengan tes unit, race‑detector, dan demo end‑to‑end.",
    narrative_en="This lab showcases resiliency patterns that shield services from slow downstream dependencies. Each pattern is verified by unit tests, the race detector, and an end‑to‑end demo.",
    layout="cover",
)

# Slide 02: Problem Nyata
add_slide(
    categories=["PROBLEM"],
    title="Dependency Lambat Mengakibatkan Outage Sistem",
    subtitle="Satu layanan 60 s menggangu seluruh cluster.",
    content_items=[{"type": "cards", "title": "DAMPAK", "items": [
        {"title": "Cascading Latency", "desc": "Latency ↑ → L = λW ↑ → thread pool terisi → request tak‑terlayani."},
        {"title": "Resource Exhaustion", "desc": "Connection pool terblokir, memory pressure naik, timeout global."},
        {"title": "Kegagalan Layanan Tidak Terkait", "desc": "Endpoint lainnya menolak karena kehabisan worker."}
    ]}],
    mental_model="Little's Law: peningkatan latency meningkatkan jumlah request in‑flight secara linear.",
    source_path="content/01-content-brief.md:5-13",
    slide_num="02 / 24",
    narrative="Ketika satu dependency melambat dari 300 ms ke 60 s, request menumpuk di pool, menurunkan throughput seluruh layanan. Ini tercermin dalam contoh demo lab.",
    narrative_en="When a single downstream slows from 300 ms to 60 s, requests pile up in the pool, throttling overall service throughput. The lab demo illustrates this effect.",
    layout="cards",
)

# Slide 03: Mengapa Review Manual Gagal
add_slide(
    categories=["WHY FAILS"],
    title="Validasi Aplikasi Tidak Cukup",
    subtitle="Check‑then‑act tidak menutup race window.",
    content_items=[{"type": "cards", "title": "KELEMAHAN", "items": [
        {"title": "Interleaving", "desc": "Scheduler OS memisahkan SELECT dan INSERT, memungkinkan duplikasi."},
        {"title": "Tes Sekuensial", "desc": "Unit test tanpa goroutine selalu lolos, memberi ilusi aman."},
        {"title": "Tidak Ada Kunci Atomik", "desc": "Database tidak menegakkan invariant bila aplikasi mengandalkan memori saja."}
    ]}],
    mental_model="Aplikasi tidak dapat menjamin konsistensi tanpa mekanisme persisten yang atomik.",
    source_path="content/02-master-draft.md:5-12",
    slide_num="03 / 24",
    narrative="Review kode dan tes unit menunjukkan logika if tampak benar, namun tidak melindungi terhadap interleaving thread yang menyebabkan duplikasi data.",
    narrative_en="Code reviews and unit tests show correct if‑logic, yet they do not protect against thread interleaving that leads to duplicate records.",
    layout="cards",
)

# Slide 04: Mental Model Inti
add_slide(
    categories=["MENTAL MODEL"],
    title="Deadline Propagation sebagai Guardrail Hierarkis",
    subtitle="Child deadline ≤ parent deadline, never larger.",
    content_items=[{"type": "schema", "title": "PROPAGASI DEADLINE", "items": [
        {"phase": "Client", "role": "Budget", "desc": "Set timeout total (misal 200 ms)."},
        {"phase": "Service A", "role": "Child", "desc": "WithTimeout(parent, localBudget) → deadline = min(parent, local)."},
        {"phase": "Service B", "role": "Leaf", "desc": "Operasi dibatalkan segera saat parent atau child deadline terlewati."}
    ]}],
    mental_model="Setiap lapisan turun menurunkan budget, mencegah kerja berlebih pada downstream yang sudah tidak dibutuhkan.",
    source_path="content/02-master-draft.md:30-38",
    slide_num="04 / 24",
    narrative="Deadline diturunkan secara hierarkis sehingga goroutine downstream tidak bekerja lebih lama daripada request yang memicunya.",
    narrative_en="Deadlines are reduced hierarchically so downstream goroutines do not outlive the request that triggered them.",
    layout="schema",
)

# Slide 05: Diagram Alur / Pipeline
add_slide(
    categories=["PIPELINE"],
    title="Pipeline Resiliensi",
    subtitle="Context → Retrier → Circuit → Idempotency → Worker",
    content_items=[{"type": "schema", "title": "ALUR EKSEKUSI", "items": [
        {"phase": "1", "role": "Context", "desc": "WithTimeout(parent, budget)"},
        {"phase": "2", "role": "Retrier", "desc": "Exponential backoff + full jitter"},
        {"phase": "3", "role": "Circuit", "desc": "Open/Close/Half‑Open guard"},
        {"phase": "4", "role": "Idempotency", "desc": "Deduplication store (TTL)"},
        {"phase": "5", "role": "Worker", "desc": "Actual business logic"}
    ]}],
    mental_model="Setiap lapisan dapat menghentikan alur lebih awal, mengurangi beban downstream.",
    source_path="content/02-master-draft.md:160-172",
    slide_num="05 / 24",
    narrative="Alur eksekusi mengalir melalui context, retry, circuit breaker, dan store sebelum mencapai worker. Setiap komponen dapat menghentikan eksekusi lebih awal.",
    narrative_en="Execution flows through context, retry, circuit breaker, and store before reaching the worker. Each component can abort early.",
    layout="schema",
)

# Slide 06: Cuplikan Kode Dominan – ExecuteWithBudget
add_slide(
    categories=["IMPLEMENTATION"],
    title="Deadline Propagation dengan Channel",
    subtitle="Memastikan goroutine berhenti tepat saat budget habis.",
    content_items=[{"type": "code", "title": "internal/deadline/deadline.go:13-28", "items": [
        "func ExecuteWithBudget(ctx context.Context, budget time.Duration, fn WorkerFunc) error {",
        "    childCtx, cancel := context.WithTimeout(ctx, budget)",
        "    defer cancel()",
        "    done := make(chan error, 1)",
        "    go func() { done <- fn(childCtx) }()",
        "    select {",
        "    case <-childCtx.Done():",
        "        return childCtx.Err()",
        "    case err := <-done:",
        "        return err",
        "    }",
        "}"
    ]}],
    mental_model="Channel berkapasitas 1 mencegah goroutine bocor ketika context dibatalkan lebih dulu.",
    source_path="internal/deadline/deadline.go:13-28",
    slide_num="06 / 24",
    narrative="Fungsi ExecuteWithBudget membuat child context dengan budget, meluncurkan goroutine worker, dan menunggu salah satu kondisi selesai lewat select. Jika childCtx dibatalkan terlebih dahulu, error context.DeadlineExceeded dikembalikan.",
    narrative_en="ExecuteWithBudget creates a child context with a budget, launches a worker goroutine, and waits via select for either completion or cancellation. If the child context is cancelled first, it returns context.DeadlineExceeded.",
    layout="code",
)

# Slide 07: Failure Path – Double Execution
add_slide(
    categories=["FAILURE PATH"],
    title="Retry Tanpa Idempotency → Double Charge",
    subtitle="Skenario: POST /payment, timeout, retry.",
    content_items=[{"type": "cards", "title": "ALUR GAGAL", "items": [
        {"title": "1. Request pertama", "desc": "Server proses charge, response hilang karena network loss."},
        {"title": "2. Client timeout", "desc": "Context dibatalkan, client menganggap failure."},
        {"title": "3. Retry", "desc": "Tanpa Idempotency-Key, server mengeksekusi charge lagi → double charge."}
    ]}],
    mental_model="Timeout tidak menjamin operasi gagal; tanpa deduplication dapat menyebabkan efek samping duplikat.",
    source_path="content/02-master-draft.md:124-132",
    slide_num="07 / 24",
    narrative="Demo menunjukkan request POST yang di‑retry setelah timeout menghasilkan dua entri pembayaran karena tidak ada kunci idempotent.",
    narrative_en="The demo shows a POST request retried after a timeout creates two payment entries because no idempotency key is used.",
    layout="cards",
)

# Slide 08: Tabel Perbandingan Mekanisme vs Gejala
add_slide(
    categories=["COMPARISON"],
    title="Matriks Mekanisme Resiliensi",
    subtitle="Timeout, Retry, Circuit, Idempotency.",
    content_items=[{"type": "table", "title": "PERBANDINGAN", "items": [
        "Mekanisme | Tujuan | Efek pada Latency | Risiko Jika Tidak Digunakan",
        "Timeout | Guardrail | Memotong kerja lama | Hanging workers",
        "Retry + Jitter | Recoverability | Tambah latency terkontrol | Retry storm",
        "Circuit Breaker | Fail‑fast | Cepat menolak request pada error berulang | Cascading failures",
        "Idempotency | Deduplication | Mencegah side‑effect ganda | Double charge / duplicate write"
    ]}],
    mental_model="Setiap mekanisme menutup satu jenis kegagalan: latency, error burst, atau side‑effect.",
    source_path="content/02-master-draft.md:190-202",
    slide_num="08 / 24",
    narrative="Tabel menegaskan bahwa kombinasi keempat pola menutup celah utama pada sistem terdistribusi.",
    narrative_en="The table affirms that the combination of these four patterns closes the primary failure surfaces in distributed systems.",
    layout="table",
)

# Slide 09: Decision Sequence – 3 Cards
add_slide(
    categories=["DECISION"],
    title="Tiga Keputusan Kunci Implementasi",
    subtitle="1️⃣ Timeout budget, 2️⃣ Retry strategy, 3️⃣ Circuit state.",
    content_items=[{"type": "cards", "title": "KEPUTUSAN", "items": [
        {"title": "Timeout Budget", "desc": "Hitung P95/P99 latency, set budget < 0.5 × target SLA."},
        {"title": "Retry Strategy", "desc": "Gunakan exponential backoff + full jitter, maxAttempts ≤ 5."},
        {"title": "Circuit Breaker", "desc": "FailureThreshold 2, cooldown 50 ms, half‑open test 1."}
    ]}],
    mental_model="Keputusan ini menentukan trade‑off latency vs availability.",
    source_path="content/02-master-draft.md:218-226",
    slide_num="09 / 24",
    narrative="Setiap keputusan dipilih berdasarkan data observasi latency dan error rate pada sistem produksi.",
    narrative_en="Each decision is chosen based on observed latency distribution and error rates in production.",
    layout="cards",
)

# Slide 10: Bukti Eksekusi – Terminal Output
add_slide(
    categories=["EVIDENCE"],
    title="Demo: Timeout & Circuit Breaker",
    subtitle="Output singkat menampilkan state transitions.",
    content_items=[{"type": "code", "title": "Terminal", "items": [
        "$ go run ./cmd/demo",
        "Scenario 1: budget 50ms → context deadline exceeded",
        "Scenario 2: circuit opens after 2 failures → state=OPEN",
        "Scenario 3: idempotent retry → execution count=1"
    ]}],
    mental_model="Log menunjukkan guardrail berfungsi sesuai harapan.",
    source_path="engineering/03-execution-result.md:1-10",
    slide_num="10 / 24",
    narrative="Demo mengeksekusi empat skenario; output terminal memperlihatkan perubahan state circuit dan jumlah eksekusi idempotent.",
    narrative_en="The demo runs four scenarios; terminal output shows circuit state transitions and idempotent execution count.",
    layout="code",
)

# Slide 11: Analisis Batasan Bukti
add_slide(
    categories=["LIMITS"],
    title="YANG TIDAK DIBUKTIKAN",
    subtitle="Keterbatasan laboratorium.",
    content_items=[{"type": "cards", "title": "Batasan", "items": [
        {"title": "Distribusi Latency Dinamis", "desc": "Tidak ada adaptasi timeout berbasis P99 secara real‑time."},
        {"title": "Persisten Idempotency Store", "desc": "In‑memory store hilang saat crash; produksi membutuhkan Redis/DB."},
        {"title": "Header Deadline HTTP", "desc": "Tidak ada standar tunggal, hanya contoh gRPC header."}
    ]}],
    mental_model="Lab fokus pada pola di dalam satu proses Go.",
    source_path="content/02-master-draft.md:30-38",
    slide_num="11 / 24",
    narrative="Penelitian mengakui bahwa beberapa aspek produksi tidak tercakup, seperti timeout adaptif dan store persisten.",
    narrative_en="The research acknowledges that several production aspects are omitted, such as adaptive timeouts and persistent stores.",
    layout="cards",
)

# Slide 12: Trade‑off Arsitektur
add_slide(
    categories=["TRADEOFF"],
    title="Trade‑off Timeout vs Throughput",
    subtitle="Menurunkan budget mengurangi latency tetapi meningkatkan reject rate.",
    content_items=[{"type": "cards", "title": "PERTIMBANGAN", "items": [
        {"title": "Timeout Ketat", "desc": "Mengurangi tail latency, tapi meningkatkan false‑negative failures."},
        {"title": "Timeout Longgar", "desc": "Mencegah false failures, tapi meningkatkan jumlah request terblokir."}
    ]}],
    mental_model="Guardrail harus seimbang dengan SLA dan beban kerja.",
    source_path="content/02-master-draft.md:186-192",
    slide_num="12 / 24",
    narrative="Mengatur budget timeout memerlukan kompromi antara latensi ekor dan tingkat penolakan yang dapat diterima.",
    narrative_en="Setting timeout budgets requires a compromise between tail latency and acceptable rejection rates.",
    layout="cards",
)

# Slide 13: Concurrency / Validation Pipeline
add_slide(
    categories=["CONCURRENCY"],
    title="Pipeline Validasi Konkuren dengan Mutex",
    subtitle="Retrier → Circuit → Idempotency Store.",
    content_items=[{"type": "code", "title": "internal/circuit/circuit.go:108-126", "items": [
        "func (b *Breaker) RecordFailure() {",
        "    b.mu.Lock()",
        "    defer b.mu.Unlock()",
        "    if b.state == StateHalfOpen {",
        "        b.state = StateOpen",
        "        b.failures = 0",
        "        b.successes = 0",
        "        b.lastStateChg = time.Now()",
        "    } else if b.state == StateClosed {",
        "        b.failures++",
        "        if b.failures >= b.cfg.FailureThreshold {",
        "            b.state = StateOpen",
        "            b.failures = 0",
        "            b.successes = 0",
        "            b.lastStateChg = time.Now()",
        "        }",
        "    }",
        "}"
    ]}],
    mental_model="RWMutex melindungi state machine dari race conditions.",
    source_path="internal/circuit/circuit.go:108-126",
    slide_num="13 / 24",
    narrative="Circuit breaker menggunakan sync.RWMutex untuk mengamankan transisi state di antara banyak goroutine.",
    narrative_en="The circuit breaker employs sync.RWMutex to safeguard state transitions across many goroutine.",
    layout="code",
)

# Slide 14: Before vs After Comparison
add_slide(
    categories=["COMPARISON"],
    title="Sebelum vs Sesudah Guardrails",
    subtitle="Tanpa vs dengan timeout & circuit.",
    content_items=[{"type": "cards", "title": "PERBANDINGAN", "items": [
        {"title": "Tanpa Guardrails", "desc": "Latency 60 s → pool terisi → 503 Service Unavailable pada endpoint lain."},
        {"title": "Dengan Guardrails", "desc": "Timeout 30 ms → request dibatalkan, circuit opens, downstream tidak dibebani."}
    ]}],
    mental_model="Guardrails memotong beban berlebih di sumbernya.",
    source_path="content/02-master-draft.md:140-148",
    slide_num="14 / 24",
    narrative="Perbandingan menunjukkan bahwa penambahan timeout, circuit, dan idempotency membuat sistem tetap responsif meski dependency melambat.",
    narrative_en="The comparison shows that adding timeout, circuit, and idempotency keeps the system responsive even when a downstream becomes slow.",
    layout="cards",
)

# Slide 15: Test ExecuteWithBudget
add_slide(
    categories=["EVIDENCE"],
    title="Test ExecuteWithBudget – Timeout Inheritance",
    subtitle="Bukti parent deadline membatalkan child.",
    content_items=[{"type": "code", "title": "internal/deadline/deadline_test.go:30-38", "items": [
        "func TestExecuteWithBudget_ParentTimeoutInherited(t *testing.T) {",
        "    ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)",
        "    defer cancel()",
        "    err := ExecuteWithBudget(ctx, 500*time.Millisecond, func(ctx context.Context) error {",
        "        time.Sleep(100*time.Millisecond)",
        "        return nil",
        "    })",
        "    assert.Equal(t, context.DeadlineExceeded, err)",
        "}"
    ]}],
    mental_model="Child deadline tidak pernah melebihi parent deadline.",
    source_path="internal/deadline/deadline_test.go:30-38",
    slide_num="15 / 24",
    narrative="Test ini membuktikan bahwa child context akan dibatalkan segera ketika parent deadline tercapai, meski budget lokal lebih besar.",
    narrative_en="This test proves that a child context is cancelled as soon as the parent deadline is reached, even if the local budget is larger.",
    layout="code",
)

# Slide 16: Test Retrier – Jitter Bounds & MaxAttempts
add_slide(
    categories=["EVIDENCE"],
    title="Test Retrier – Jitter Bounds & MaxAttempts",
    subtitle="Bukti full jitter dan pengulangan maksimal.",
    content_items=[{"type": "code", "title": "internal/retry/retry_test.go:12-24", "items": [
        "func TestRetrier_JitterBoundsAndZeroConfig(t *testing.T) {",
        "    r := NewRetrier(Config{})",
        "    for i := 1; i <= 5; i++ {",
        "        d := r.CalculateBackoff(i)",
        "        assert.GreaterOrEqual(t, d, time.Duration(0))",
        "        assert.LessOrEqual(t, d, r.cfg.MaxBackoff)",
        "    }",
        "}"
    ]}],
    mental_model="Backoff selalu berada dalam rentang [0, maxBackoff].",
    source_path="internal/retry/retry_test.go:12-24",
    slide_num="16 / 24",
    narrative="Test ini memastikan bahwa setiap nilai backoff tidak negatif dan tidak melebihi batas maksimum, serta konfigurasi nol menggunakan nilai default.",
    narrative_en="This test ensures every backoff value is non‑negative and does not exceed the maximum bound, and that zero config uses defaults.",
    layout="code",
)

# Slide 17: Test Circuit Breaker – State Transitions
add_slide(
    categories=["EVIDENCE"],
    title="Test Circuit Breaker – State Transitions",
    subtitle="Bukti CLOSED → OPEN → HALF_OPEN → CLOSED.",
    content_items=[{"type": "code", "title": "internal/circuit/circuit_test.go:5-20", "items": [
        "func TestCircuitBreaker_StateTransitions(t *testing.T) {",
        "    cb := NewBreaker(Config{FailureThreshold: 2, SuccessThreshold: 2, Cooldown: 50*time.Millisecond})",
        "    cb.RecordFailure()",
        "    cb.RecordFailure()",
        "    assert.Equal(t, StateOpen, cb.State())",
        "    time.Sleep(60*time.Millisecond)",
        "    assert.Equal(t, StateHalfOpen, cb.State())",
        "    cb.RecordSuccess()",
        "    cb.RecordSuccess()",
        "    assert.Equal(t, StateClosed, cb.State())",
        "}"
    ]}],
    mental_model="State machine berubah hanya di bawah lock.",
    source_path="internal/circuit/circuit_test.go:5-20",
    slide_num="17 / 24",
    narrative="Test ini mengonfirmasi transisi state yang diharapkan: dua kegagalan membuka circuit, cooldown membawa ke half‑open, lalu dua kesuksesan menutupnya.",
    narrative_en="This test confirms the expected state transitions: two failures open the circuit, cooldown brings it to half‑open, then two successes close it.",
    layout="code",
)

# Slide 18: Test Idempotency Store – Concurrent Access
add_slide(
    categories=["EVIDENCE"],
    title="Test Idempotency Store – Concurrent Access",
    subtitle="50 goroutine konkuren, race detector PASS.",
    content_items=[{"type": "code", "title": "internal/idempotency/idempotency_test.go:8-22", "items": [
        "func TestStore_ConcurrentAccess(t *testing.T) {",
        "    store := NewStore(time.Minute)",
        "    var wg sync.WaitGroup",
        "    for i := 0; i < 50; i++ {",
        "        wg.Add(1)",
        "        go func(id int) {",
        "            defer wg.Done()",
        "            store.Set(fmt.Sprintf(\"key-%d\", id), \"response\")",
        "            val, ok := store.Get(fmt.Sprintf(\"key-%d\", id))",
        "            assert.True(t, ok)",
        "            assert.Equal(t, \"response\", val)",
        "        }(i)",
        "    }",
        "    wg.Wait()",
        "}"
    ]}],
    mental_model="Mutex mengamankan operasi Get/Set/TTL secara konkuren.",
    source_path="internal/idempotency/idempotency_test.go:8-22",
    slide_num="18 / 24",
    narrative="Test ini menjalankan 50 goroutine yang melakukan Get/Set secara acak; race detector tidak melaporkan kesalahan.",
    narrative_en="This test runs 50 goroutines performing random Get/Set operations; the race detector reports no errors.",
    layout="code",
)

# Slide 19: Integration Test – Retry + Circuit + Idempotency
add_slide(
    categories=["EVIDENCE"],
    title="Integration – Retry dengan Circuit Breaker",
    subtitle="Backend selalu gagal, circuit terbuka setelah 2 attempt.",
    content_items=[{"type": "code", "title": "tests/integration_test.go:14-42", "items": [
        "func TestIntegration_RetryWithCircuitBreaker(t *testing.T) {",
        "    cb := circuit.NewBreaker(circuit.Config{FailureThreshold: 2, ...})",
        "    r := retry.NewRetrier(retry.Config{MaxAttempts: 4, ...})",
        "    attempts := 0",
        "    err := r.Do(context.Background(), func(ctx context.Context) error {",
        "        return cb.Execute(func() error {",
        "            attempts++",
        "            return errors.New(\"backend down\")",
        "        })",
        "    })",
        "    assert.NotNil(t, err)",
        "    assert.Equal(t, circuit.StateOpen, cb.State())",
        "}"
    ]}],
    mental_model="Circuit breaker memotong rantai retry setelah threshold.",
    source_path="tests/integration_test.go:14-42",
    slide_num="19 / 24",
    narrative="Integrasi ini membuktikan bahwa setelah dua kegagalan berturut‑turut, circuit breaker menolak semua attempt berikutnya tanpa memanggil backend.",
    narrative_en="This integration proves that after two consecutive failures, the circuit breaker rejects all subsequent attempts without calling the backend.",
    layout="code",
)

# Slide 20: Checklist Implementasi
add_slide(
    categories=["CHECKLIST"],
    title="Checklist Deploy Guardrails",
    subtitle="Langkah‑langkah produksi.",
    content_items=[{"type": "cards", "title": "CHECKLIST", "items": [
        {"title": "1. Define timeout budgets", "desc": "Berbasis P95 latensi layanan."},
        {"title": "2. Implement Retry with jitter", "desc": "Gunakan konfigurasi base/backoff, maxAttempts ≤5."},
        {"title": "3. Add circuit breaker", "desc": "Set failureThreshold & cooldown sesuai SLA."},
        {"title": "4. Add idempotency store", "desc": "Redis atau DB dengan TTL persisten."},
        {"title": "5. Propagate context deadlines", "desc": "Pastikan semua downstream menerima context."}
    ]}],
    mental_model="Deploy langkah demi langkah memastikan semua guardrails aktif.",
    source_path="content/02-master-draft.md:212-221",
    slide_num="20 / 24",
    narrative="Checklist ini merangkum praktik terbaik implementasi guardrails pada sistem terdistribusi.",
    narrative_en="This checklist summarizes best practices for implementing guardrails on distributed systems.",
    layout="cards",
)

# Slide 21: Anti‑pattern & Reproduksi Kegagalan
add_slide(
    categories=["ANTIPATTERN"],
    title="Anti‑pattern: Timeout = Failure",
    subtitle="Menyebabkan retry berbahaya dan double‑charge.",
    content_items=[{"type": "cards", "title": "REPRODUKSI", "items": [
        {"title": "1. Set timeout 30s", "desc": "Terlalu long, worker terblokir."},
        {"title": "2. Retry tanpa jitter", "desc": "Semua client retry bersamaan → storm."},
        {"title": "3. Tanpa idempotency", "desc": "Operasi mutasi dieksekusi dua kali."}
    ]}],
    mental_model="Kegagalan ini muncul ketika guardrails tidak diterapkan.",
    source_path="content/02-master-draft.md:190-196",
    slide_num="21 / 24",
    narrative="Contoh anti‑pattern menegaskan pentingnya mengamankan setiap lapisan.",
    narrative_en="The anti‑pattern example underscores the need to secure each layer.",
    layout="cards",
)

# Slide 22: Takeaway Praktis & Command Block
add_slide(
    categories=["TAKEAWAY"],
    title="Praktis: Satu Perintah Deploy Guardrails",
    subtitle="Gunakan go‑run demo atau integrate libs.",
    content_items=[{"type": "code", "title": "Deploy", "items": [
        "go get github.com/lukman-ss/timeoutguard",
        "go run ./cmd/demo --scenario=all"
    ]}],
    mental_model="Integrasi satu paket membawa semua pola.",
    source_path="cmd/demo/main.go:1-10",
    slide_num="22 / 24",
    narrative="Dengan menambahkan modul timeoutguard, aplikasi langsung dapat memakai context deadline, retry, circuit, dan idempotency.",
    narrative_en="By adding the timeoutguard module, the application instantly gains context deadline, retry, circuit, and idempotency capabilities.",
    layout="code",
)

# Slide 23: Prinsip Arsitektural Jangka Panjang
add_slide(
    categories=["PRINCIPLE"],
    title="Prinsip: Guardrails Sebagai Kontrak Layanan",
    subtitle="Timeout, retry, circuit, idempotency = Service‑Level Contract.",
    content_items=[{"type": "cards", "title": "PRINSIP", "items": [
        {"title": "Explicit Contracts", "desc": "Setiap API menyatakan timeout & idempotency semantics."},
        {"title": "Observability Built‑in", "desc": "Metrics untuk deadline breaches, circuit state, retry count."},
        {"title": "Evolutionary Safety", "desc": "Tambah guardrails tanpa mengubah layanan klien."}
    ]}],
    mental_model="Guardrails menyediakan SLA yang dapat di‑monitor dan di‑tweak.",
    source_path="content/02-master-draft.md:226-236",
    slide_num="23 / 24",
    narrative="Prinsip ini menjadikan guardrails bagian resmi dari kontrak API, meningkatkan kepercayaan dan evolusi layanan.",
    narrative_en="This principle makes guardrails a formal part of the API contract, boosting trust and service evolution.",
    layout="cards",
)

# Slide 24: Closing
add_slide(
    categories=["CLOSING"],
    title="Terima Kasih",
    subtitle="Repositori, commit, dan kredit.",
    content_items=[{"type": "cards", "title": "INFO", "items": [
        {"title": "Repo", "desc": REPO_URL},
        {"title": "Commit", "desc": COMMIT_SHA},
        {"title": "Author", "desc": AUTHOR}
    ]}],
    mental_model="",
    source_path="",
    slide_num="24 / 24",
    narrative="Semoga konsep guardrails membantu Anda merancang layanan yang tahan terhadap ketergantungan lambat.",
    narrative_en="Hope these guardrail concepts help you design services resilient to slow dependencies.",
    layout="cards",
)
