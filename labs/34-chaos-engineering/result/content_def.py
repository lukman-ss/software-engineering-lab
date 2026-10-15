#!/usr/bin/env python3
"""
Publication Builder for Software Engineering Lab 34: Chaos Engineering & Fault Injection.
Generates 24-slide LinkedIn carousel covering steady-state verification, controlled fault injection,
circuit breaker state transitions, graceful degradation via fallback, and automated blast-radius containment.
"""
import os

OUTPUT_DIR = "labs/34-chaos-engineering/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/34-chaos-engineering"
COMMIT_SHA = "200fe73bb5e43da154526620201573c7af03cf8e"
LAB_NO = "34"
SLUG = "chaos-engineering"
TITLE_ID = "Chaos Engineering & Fault Injection"
SUBTITLE_ID = "Memverifikasi Ketahanan Layanan Terdistribusi Menggunakan Steady-State, Circuit Breaker, dan Auto-Abort"
TITLE_EN = "Chaos Engineering & Fault Injection"
SUBTITLE_EN = "Verifying Distributed System Resilience via Steady-State Metrics, Circuit Breakers, and Auto-Abort"
AUTHOR = "Lukman / lukman-ss"
COVER_IMAGE_PATH = "labs/34-chaos-engineering/result/image-cover/cover_art.jpg"

DARK_NAVY = "0b1b2d"
TEAL = "087f76"
CREAM = "f4f1e8"
WHITE = "ffffff"
SLATE = "52606b"
ORANGE = "d97706"

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

# =====================================================================
# 24 SLIDE DEFINITIONS
# =====================================================================

# Slide 1: Cover
add_slide(
    categories=["COVER"],
    title=TITLE_ID,
    subtitle=SUBTITLE_ID,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "Controlled fault injection: Latency & forced error pipeline",
            "Resilience verification: Circuit Breaker state transitions, fallback protection, and auto-abort monitor",
            f"Penulis: {AUTHOR}",
        ]}
    ],
    mental_model="Resilience is not an assumption; it is an empirical property proven under controlled turbulence.",
    source_path=REPO_URL,
    slide_num="01 / 24",
    narrative="Dalam arsitektur microservices modern, dependensi eksternal seperti payment gateway dan downstream API pasti akan mengalami gangguan latensi atau kegagalan fatal. Tanpa pengujian proaktif, asumsi ketahanan sistem hanyalah ilusi teoretis yang runtuh saat insiden production terjadi. Laboratorium ini mengimplementasikan framework Chaos Engineering murni dalam Go untuk membuktikan secara empiris bagaimana kombinasi steady-state monitor, Circuit Breaker, fallback, dan auto-abort blast radius mengisolasi kegagalan tanpa merusak metrik pengguna.",
    narrative_en="In modern distributed architectures, downstream dependencies such as payment gateways inevitably experience latency spikes and service outages. Without proactive failure testing, architectural resilience remains an untested theoretical hypothesis until production incidents strike. This laboratory implements a pure Go Chaos Engineering framework proving empirically how steady-state monitoring, Circuit Breakers, graceful fallbacks, and automated blast-radius containment isolate faults without degrading end-user availability.",
    layout="cover"
)

# Slide 2: Problem nyata
add_slide(
    categories=["MOTIVASI", "SYSTEM COLLAPSE"],
    title="Masalah Nyata: Ilusi Ketahanan Pasif",
    subtitle="Downstream failure yang tidak ditangani memicu cascading failure dan resource exhaustion.",
    content_items=[
        {"title": "UNTESTED HYPOTHESIS", "items": [
            "Arsitek mengasumsikan timeout dan failover bekerja tanpa pernah mengujinya secara sengaja.",
            "Integrasi downstream dianggap selalu beroperasi pada SLA latensi sub-100 milidetik normal.",
            "Ketahanan hanya berupa diagram arsitektur pasif di atas kertas tanpa bukti empiris."
        ], "type": "card"},
        {"title": "CASCADING THREAD STARVATION", "items": [
            "Layanan downstream melambat atau gagal merespons, menahan thread koneksi klien pemanggil.",
            "Pool koneksi HTTP/gRPC dan worker thread hulu cepat habis (resource exhaustion).",
            "Satu dependensi non-kritis yang macet melumpuhkan seluruh layanan upstream secara berantai."
        ], "type": "card"},
        {"title": "RUNAWAY EXPERIMENT RISK", "items": [
            "Pengujian chaos yang dilakukan secara serampangan berpotensi merusak seluruh basis pengguna.",
            "Tidak adanya tombol pembatalan darurat otomatis saat error rate melonjak tajam.",
            "Blast radius tidak terkendali mengubah pengujian keandalan menjadi pemadaman sistem total."
        ], "type": "card"}
    ],
    narrative="Kelemahan paling umum dalam sistem terdistribusi adalah ilusi bahwa konfigurasi timeout atau failover pasif otomatis menjamin ketersediaan. Saat layanan pihak ketiga seperti gateway pembayaran mengalami freeze jaringan atau lonjakan latensi, request pemanggil menumpuk dan menahan thread sistem hulu. Dalam hitungan detik, seluruh worker pool habis dan memicu cascading failure ke seluruh ekosistem backend. Selain itu, menyuntikkan kegagalan tanpa mekanisme batas dampak (blast radius) justru membahayakan ketersediaan sistem.",
    narrative_en="The most critical flaw in distributed systems is the illusion that passive timeout configurations guarantee high availability. When a third-party dependency suffers a network freeze or latency surge, inbound requests stall and hold upstream worker threads hostage. Within seconds, socket pools and execution threads are exhausted, triggering catastrophic cascading failures across the entire cluster. Furthermore, executing chaos experiments without strict blast-radius containment risks transforming test routines into full-blown customer outages.",
    mental_model="Passive resilience is untested faith; unmitigated downstream latency starves upstream execution pools.",
    source_path="labs/34-chaos-engineering/content/01-content-brief.md:9-14",
    slide_num="02 / 24",
    layout="cards"
)

# Slide 3: Mengapa review manual / pendekatan naif gagal
add_slide(
    categories=["ANALISIS", "PENDEKATAN NAIF"],
    title="Mengapa Penanganan Error Konvensional Gagal?",
    subtitle="Pola try-catch sederhana dan mock unit test statis buta terhadap dinamika kegagalan terdistribusi.",
    content_items=[
        {"title": "KODE PEMANGGILAN NAIF", "items": [
            "func (c *Client) CallPaymentNaive(ctx context.Context) error {",
            "    // Tidak ada Circuit Breaker atau Fallback",
            "    resp, err := http.Post(\"https://pay-gateway/charge\", ...)",
            "    if err != nil {",
            "        // Error diteruskan mentah ke hulu; user melihat 500",
            "        return err",
            "    }",
            "    return process(resp)",
            "}"
        ], "type": "code"},
        {"title": "KEGAGALAN FUNDAMENTAL", "items": [
            "Blind Re-execution: Klien terus membombardir downstream yang sekarat, memperparah downtime.",
            "Mock Deception: Unit test statis menguji mock sukses; tidak pernah menguji latensi dan partial failure.",
            "Absence of Degradation: Tidak ada rute alternatif (fallback queue) saat jalur utama terputus.",
            "Metrik Terisolasi: Tim memantau CPU/RAM internal alih-alih metrik kepuasan pengguna (error rate)."
        ], "type": "bullet"}
    ],
    narrative="Pendekatan naif mengandalkan penanganan error lokal konvensional yang langsung meneruskan kegagalan downstream ke pengguna akhir. Saat gateway eksternal down, setiap request masuk tetap mencoba membuka koneksi HTTP baru, membombardir server yang sedang sekarat dan mempercepat kehabisan socket pemanggil. Unit test tradisional memperparah ilusi ini karena hanya memverifikasi mock objek dengan waktu respons instan. Tanpa pengujian kegagalan dinamis dan jalur degradasi anggun (fallback), pengguna langsung dihadapkan pada layar error 500.",
    narrative_en="The naive approach relies on conventional error handling that blindly propagates downstream failures directly to end users. When an external gateway falters, every incoming client request continues attempting new HTTP handshakes, battering the failing dependency and exhausting caller sockets. Traditional unit tests exacerbate this false confidence by asserting against instant in-memory mocks. Without dynamic turbulence testing and graceful degradation pathways, any remote outage directly yields HTTP 500 errors to customers.",
    mental_model="Direct passthrough error handling couples your availability directly to your weakest external dependency.",
    source_path="labs/34-chaos-engineering/content/02-master-draft.md:3-8",
    slide_num="03 / 24",
    layout="code-bullet"
)

# Slide 4: Mental model inti / schema card
add_slide(
    categories=["MENTAL MODEL", "METRIKA"],
    title="Mental Model: Eksperimen Ilmiah & Steady State",
    subtitle="Chaos engineering bukan merusak sistem, melainkan memverifikasi hipotesis stabilitas output.",
    content_items=[
        {"title": "STRUKTUR METRIK STEADY-STATE", "items": [
            "type SteadyStateMetrics struct {",
            "    TotalRequests   uint64 // Total beban request masuk",
            "    FailedRequests  uint64 // Jumlah kegagalan di sisi klien",
            "    SuccessRequests uint64 // Total transaksi yang berhasil/degraded",
            "}",
            "",
            "// Evaluasi kesehatan sistem berdasarkan output batas pengguna",
            "func (m *Monitor) IsHealthy() bool {",
            "    if atomic.LoadUint64(&m.totalRequests) < 5 { return true }",
            "    return m.ErrorRate() <= m.maxErrorRate // Ambang batas lab: 25%",
            "}"
        ], "type": "code"},
        {"title": "4 PILAR KETAHANAN TERUKUR", "items": [
            "1. Steady-State Metric: Ukur metrik output pengguna (error rate), bukan kesehatan internal node.",
            "2. Empirical Hypothesis: Hipotesiskan bahwa error rate tetap di bawah batas aman meski fault aktif.",
            "3. Controlled Injection: Sisipkan fault deterministik (latensi & error) pada titik batas integrasi.",
            "4. Automated Abort: Hentikan injeksi seketika saat metrik melanggar threshold untuk melindungi sistem."
        ], "type": "bullet"}
    ],
    narrative="Mental model Chaos Engineering dibangun di atas metode ilmiah yang terstruktur. Sistem didefinisikan pertama kali melalui steady-state metric, yaitu indikator kesehatan output sistem dari sudut pandang konsumen akhir (seperti error rate atau transaksi sukses), bukan metrik internal seperti penggunaan memori. Hipotesis empiris dirumuskan: 'Meskipun downstream disuntikkan gangguan, sistem mampu menjaga error rate di bawah ambang batas toleransi.' Jika hipotesis gagal dan batas dilanggar, monitor otomatis membatalkan eksperimen.",
    narrative_en="The core mental model of Chaos Engineering is anchored in rigorous scientific experimentation. A system is defined primarily through steady-state output metrics—observing client-facing error rates and transaction completions rather than internal infrastructure telemetry like CPU usage. An empirical hypothesis is constructed: 'Even when downstream faults are injected, the architecture maintains error rates below defined safety margins.' If this hypothesis fails and boundaries breach, monitoring triggers an immediate automated abort.",
    mental_model="Define steady state on user output; formulate an empirical hypothesis; auto-abort on boundary breach.",
    source_path="internal/monitor/monitor.go:10-48 • labs/34-chaos-engineering/content/01-content-brief.md:12-14",
    slide_num="04 / 24",
    layout="schema-card"
)

# Slide 5: Diagram alur/pipeline
add_slide(
    categories=["ARSITEKTUR", "PIPELINE"],
    title="Pipeline Arsitektur: Client Berketahanan & Auto-Abort",
    subtitle="Alur request端到端: Client Request → Circuit Breaker → Fault Injector → Steady-State Monitor.",
    content_items=[
        {"title": "SEKUENSI PIPELINE EKSPERIMEN CHAOS", "items": [
            "1. Client Request masuk ke antarmuka Resilient Client.",
            "2. CircuitBreaker memeriksa status internal: CLOSED, OPEN, atau HALF-OPEN.",
            "3. Jika state CLOSED: Request dieksekusi menuju Downstream Service.",
            "4. Downstream melewati Fault Injector (mengevaluasi injeksi latency dan forced error).",
            "5. Jika Downstream Error: Circuit Breaker menambah failure counter.",
            "6. Jika Failure >= Threshold: Circuit Breaker trip ke OPEN dan langsung memicu Fallback.",
            "7. Graceful Fallback mengembalikan respons alternatif dan mencatat SUKSES pada Monitor.",
            "8. Monitor Ticker secara periodik memeriksa IsHealthy() (ErrorRate <= MaxThreshold).",
            "9. Jika unmitigated error menembus threshold: Experiment runner memicu terminate(ABORTED) & Clear()."
        ], "type": "diagram"}
    ],
    narrative="Diagram pipeline ini menunjukkan interaksi dinamis antara pemanggil, pelindung sirkuit, dan instrumen eksperimen. Ketika request masuk, Circuit Breaker menentukan apakah panggilan aman diteruskan. Jalur pemanggilan melewati Fault Injector yang dapat menyisipkan jeda waktu atau memaksakan error secara terprogram. Jika terjadi kegagalan berturut-turut, Circuit Breaker beralih ke state OPEN dan mengalihkan trafik ke fungsi fallback. Pada saat yang sama, steady-state monitor mengevaluasi rasio error; jika terjadi lonjakan tidak termitigasi, runner membatalkan eksperimen seketika.",
    narrative_en="This pipeline diagram captures the dynamic interplay between the caller, protective circuit breakers, and chaos orchestration components. Inbound requests pass through the Circuit Breaker to assess execution safety. The invocation route crosses a Fault Injector programmed to simulate network latency or hard errors. If consecutive failures exceed threshold, the breaker trips to OPEN and delegates execution to the graceful fallback routine. Concurrently, the monitor polices error rates, triggering an instant abort if unmitigated faults compromise system stability.",
    mental_model="Circuit breakers protect upstream resources; steady-state monitors guard the safety of the chaos experiment.",
    source_path="labs/34-chaos-engineering/content/04-diagrams.md:38-73",
    slide_num="05 / 24",
    layout="pipeline-diagram"
)

# Slide 6: Cuplikan kode dominan
add_slide(
    categories=["IMPLEMENTASI", "CODE"],
    title="State Machine Circuit Breaker dengan Fallback",
    subtitle="Eksekusi aman, deteksi threshold, transisi cooldown, dan delegasi fallback di circuitbreaker.go.",
    content_items=[
        {"title": "CIRCUIT BREAKER EXECUTE & FALLBACK ENGINE", "items": [
            "func (cb *CircuitBreaker) Execute(fn func() error, fallback func() error) error {",
            "    cb.mu.Lock()",
            "    cb.checkStateLocked() // Cooldown elapsed? Transisi OPEN -> HALF-OPEN",
            "    if cb.state == StateOpen {",
            "        cb.mu.Unlock()",
            "        if fallback != nil { return fallback() }",
            "        return ErrCircuitOpen",
            "    }",
            "    cb.mu.Unlock()",
            "    err := fn()",
            "    cb.mu.Lock()",
            "    defer cb.mu.Unlock()",
            "    if err != nil {",
            "        cb.failures++",
            "        if cb.failures >= cb.threshold || cb.state == StateHalfOpen {",
            "            cb.state = StateOpen",
            "            cb.lastStateChg = time.Now()",
            "        }",
            "        if fallback != nil { return fallback() }",
            "        return err",
            "    }",
            "    // Eksekusi sukses memulihkan state dari HALF-OPEN ke CLOSED",
            "    if cb.state == StateHalfOpen || cb.state == StateClosed { cb.state = StateClosed; cb.failures = 0 }",
            "    return nil",
            "}"
        ], "type": "code"}
    ],
    narrative="Implementasi Circuit Breaker ini adalah inti proteksi terhadap kegagalan downstream. Sebelum mengeksekusi fungsi utama `fn`, method memeriksa apakah waktu jeda `cooldown` telah terlewati untuk mencoba transisi ke `StateHalfOpen`. Jika sirkuit dalam keadaan terbuka (`StateOpen`), panggilan eksternal langsung dilewati (fast-fail) dan dialihkan ke fungsi `fallback`. Jika fungsi utama menghasilkan error, counter kegagalan dinaikkan hingga memicu pemutusan sirkuit. Sebaliknya, eksekusi sukses pada mode probation mereset sirkuit kembali ke `StateClosed`.",
    narrative_en="This Circuit Breaker implementation serves as the primary barrier against downstream service degradation. Prior to executing the primary function `fn`, it evaluates whether the cooldown window has elapsed to allow speculative probes under `StateHalfOpen`. If the breaker is OPEN, downstream execution is bypassed immediately (fast-fail) and delegated directly to the fallback closure. If the primary call fails, the failure counter increments toward the trip threshold. Conversely, a successful attempt restores the breaker to `StateClosed`.",
    mental_model="Fast-fail saves upstream threads; graceful fallback preserves downstream business continuity.",
    source_path="internal/circuitbreaker/circuitbreaker.go:155-195",
    slide_num="06 / 24",
    layout="code-dominant"
)

# Slide 7: Failure path konkret
add_slide(
    categories=["FAILURE PATH", "BLAST RADIUS"],
    title="Failure Path: Pemadaman Akibat Unmitigated Fault",
    subtitle="Ketika fault downstream tidak dimitigasi fallback, steady-state monitor memicu auto-abort seketika.",
    content_items=[
        {"title": "TERMINAL TRACE: AUTO-ABORT PADA VIOLATION", "items": [
            "[3] Running Second Chaos Experiment: Unmitigated Downstream Fault (Triggers Auto-Abort)...",
            "Raw Call #1 Failed: chaos: injected fault failure",
            "Raw Call #2 Failed: chaos: injected fault failure",
            "Experiment State: ABORTED",
            "Abort Reason: Steady state breached: error rate 33.33% exceeded threshold",
            "Fault Injector Active: false (Neutralized on Abort)",
            "",
            "[4] Traffic After Experiment Ends / Recovery...",
            "Req #16 -> Payment Success (Primary) | CB State: CLOSED"
        ], "type": "terminal"},
        {"title": "ANALISIS KEGAGALAN TERKENDALI", "items": [
            "Panggilan raw service tanpa perlindungan Circuit Breaker/Fallback langsung mencatat kegagalan.",
            "Setelah 2 kali kegagalan berturut-turut, error rate kumulatif melonjak mencapai 33.33%.",
            "Threshold keamanan lab (20%) terlampaui; background ticker mendeteksi pelanggaran steady-state.",
            "Runner memicu abort darurat dan memanggil injector.Clear() secara sinkron untuk memulihkan sistem."
        ], "type": "bullet"}
    ],
    narrative="Failure path ini merekam perilaku sistem ketika gangguan disuntikkan tanpa perlindungan Circuit Breaker dan Fallback. Tanpa mitigasi, dua request downstream langsung gagal akibat `ErrInjectedFault`, mendorong rasio error kumulatif melonjak ke 33.33%, melampaui ambang batas toleransi 20%. Background monitor yang berjalan via ticker mendeteksi pelanggaran steady-state ini dalam hitungan milidetik. Eksperimen dihentikan dengan status ABORTED, dan fungsi injector dinetralkan seketika, membuktikan penahanan blast radius bekerja sempurna.",
    narrative_en="This failure trace documents system behavior when turbulence is injected without the buffer of Circuit Breakers or Fallbacks. Lacking mitigation, two unshielded calls fail instantly from `ErrInjectedFault`, driving cumulative error rates to 33.33% and breaching the 20% safety threshold. The asynchronous background monitor detects this steady-state breach within milliseconds. The experiment runner immediately aborts the run and synchronously invokes `injector.Clear()`, proving that blast-radius containment shields system availability.",
    mental_model="Unmitigated faults trigger automated aborts, neutralizing injectors before blast radius escapes.",
    source_path="cmd/demo/main.go:90-125 • engineering/03-execution-result.md:75-81",
    slide_num="07 / 24",
    layout="failure-path"
)

# Slide 8: Tabel perbandingan mekanisme vs gejala
add_slide(
    categories=["PERBANDINGAN", "MEKANISME"],
    title="Mekanisme Ketahanan vs Gejala Kegagalan",
    subtitle="Membandingkan efektivitas proteksi arsitektur dalam meredam anomali sistem terdistribusi.",
    content_items=[
        {"title": "TABEL PERBANDINGAN MEKANISME RESILIENSI", "items": [
            "Pendekatan | Mekanisme Proteksi | Gejala Saat Downstream Down | Garansi Ketersediaan",
            "Panggilan Naif | Direct HTTP call tanpa timeout | Resource exhaustion, thread hanging | Buruk (0% ketersediaan)",
            "Timeout Saja | Batas durasi context.WithTimeout | Thread tertahan hingga batas waktu | Sedang (kegagalan tetap lolos)",
            "Circuit Breaker | Fast-fail setelah failure threshold | Eksekusi diputus, error cepat | Baik (upstream terlindungi)",
            "CB + Fallback | Fast-fail dialihkan ke rute alternatif | Pengguna menerima hasil terdegradasi | Unggul (error rate 0.00%)",
            "Auto-Abort Monitor | Pemantauan steady-state berkelanjutan | Eksperimen dibatalkan saat bahaya | Mutlak (blast radius aman)"
        ], "type": "table"}
    ],
    narrative="Tabel ini mengevaluasi berbagai lapisan perlindungan sistem dalam menghadapi gangguan downstream. Panggilan naif adalah yang paling rentan karena dapat menahan thread eksekusi hingga memicu pemadaman total. Penambahan timeout mencegah kebocoran thread tak terbatas namun tetap meneruskan error mentah ke pengguna. Circuit Breaker memberikan perlindungan proaktif dengan memutus panggilan ke dependensi bermasalah. Ketika dipadukan dengan Fallback dan Auto-Abort, sistem mencapai ketahanan paripurna: pengguna tetap dilayani dan eksperimen tidak membahayakan sistem.",
    narrative_en="This comparison matrix evaluates layered architectural defenses against downstream degradation. Naive client invocations represent the highest risk, risking unbounded thread exhaustion and total cluster failure. Timeouts bound thread exposure but continue forwarding raw failures upstream. Circuit Breakers isolate troubled dependencies by failing fast. When paired with graceful fallback routines and automated abort monitors, the topology achieves optimal resilience: end users experience uninterrupted service while chaos experiments remain strictly bounded.",
    mental_model="Timeouts protect threads; circuit breakers protect downstream dependencies; fallbacks protect users.",
    source_path="labs/34-chaos-engineering/content/02-master-draft.md:16-26",
    slide_num="08 / 24",
    layout="comparison-table"
)

# Slide 9: Decision sequence (3 focused decision cards)
add_slide(
    categories=["DESIGN DECISION", "ARSITEKTUR"],
    title="3 Keputusan Desain Kritis Chaos Engineering",
    subtitle="Urutan keputusan teknis untuk membangun framework eksperimen kegagalan yang aman dan deterministik.",
    content_items=[
        {"title": "1. STEADY-STATE AT THE BOUNDARY", "items": [
            "Pilih metrik output yang dirasakan langsung oleh pengguna akhir (error rate, transaction throughput).",
            "Hindari menjadikan metrik internal (CPU, memory, garbage collection time) sebagai steady state utama.",
            "Terapkan syarat minimum sampel (contoh: minimal 5 request) sebelum mengevaluasi pelanggaran batas."
        ], "type": "card"},
        {"title": "2. TWO-PHASE CIRCUIT TRANSITIONS", "items": [
            "Konfigurasikan failure threshold terukur untuk memutus sirkuit ke OPEN sebelum thread pool jenuh.",
            "Wajib sediakan durasi cooldown sebelum probation (HALF-OPEN) guna menguji kesembuhan dependensi.",
            "Gunakan fungsi fallback terpisah untuk menelan error dan menjaga degradasi layanan tetap anggun."
        ], "type": "card"},
        {"title": "3. SYNCHRONOUS INJECTOR NEUTRALIZATION", "items": [
            "Experiment runner wajib mengikat ticker pemantau dengan mekanisme pembatalan konteks.",
            "Panggil injector.Clear() secara sinkron di dalam fungsi terminate() pada status ABORTED atau COMPLETED.",
            "Pastikan pembersihan fault tidak terblokir oleh goroutine lain untuk mencegah kebocoran gangguan."
        ], "type": "card"}
    ],
    narrative="Membangun sistem pengujian chaos yang andal menuntut tiga keputusan arsitektural penting. Pertama, metrik steady-state wajib didefinisikan pada batas pengalaman pengguna akhir, bukan pada telemetri internal server, dengan batas sampel minimum agar tidak memicu alarm palsu. Kedua, Circuit Breaker harus mengimplementasikan transisi dua fase yang jelas (termasuk fase probasi Half-Open) serta mendukung fungsi fallback. Ketiga, penetralan alat injeksi harus dieksekusi secara sinkron dan instan saat terjadi pembatalan darurat.",
    narrative_en="Engineering a reliable chaos testing harness demands three foundational architectural decisions. First, steady-state metrics must be measured at the consumer boundary rather than through internal node telemetry, backed by minimum sample thresholds to prevent spurious alerts. Second, the Circuit Breaker must follow a disciplined three-state machine featuring an intermediate Half-Open probe and first-class fallback integration. Third, injector neutralization must occur synchronously during emergency abort routines to eliminate fault leakage.",
    mental_model="Measure user outcomes, decouple failing paths via breakers, and guarantee synchronous fault teardown.",
    source_path="labs/34-chaos-engineering/content/05-key-takeaways.md:1-10",
    slide_num="09 / 24",
    layout="decision-cards"
)

# Slide 10: Bukti eksekusi / terminal output
add_slide(
    categories=["BUKTI", "DEMO"],
    title="Bukti Eksekusi: Siklus Lengkap Eksperimen Chaos",
    subtitle="Output terminal dari go run ./cmd/demo memperlihatkan baseline, fallback, auto-abort, dan recovery.",
    content_items=[
        {"title": "HASIL EKSEKUSI DEMO LENGKAP", "items": [
            "=== Chaos Engineering & Fault Injection Demo ===",
            "[1] Baseline Steady-State Traffic (Normal Operation)...",
            "Req #1 -> Payment Success (Primary) | CB State: CLOSED",
            "Req #5 -> Payment Success (Primary) | CB State: CLOSED",
            "[2] Running Chaos Experiment: Injecting Downstream Failures with Resilient Client...",
            "Req #6 -> Payment Queued (Fallback) | CB State: CLOSED",
            "Req #8 -> Payment Queued (Fallback) | CB State: OPEN",
            "Req #15 -> Payment Queued (Fallback) | CB State: OPEN",
            "Experiment State: COMPLETED",
            "Steady-State Metrics with Fallback: Total=15, Success=15, Failed=0, ErrorRate=0.00%",
            "[3] Running Second Chaos Experiment: Unmitigated Downstream Fault (Triggers Auto-Abort)...",
            "Raw Call #1 Failed: chaos: injected fault failure",
            "Experiment State: ABORTED",
            "Abort Reason: Steady state breached: error rate 33.33% exceeded threshold",
            "Fault Injector Active: false (Neutralized on Abort)",
            "[4] Traffic After Experiment Ends / Recovery...",
            "Req #16 -> Payment Success (Primary) | CB State: CLOSED",
            "=== Demo Complete ==="
        ], "type": "terminal"}
    ],
    narrative="Bukti eksekusi terminal ini mengonfirmasi keempat fase eksperimen berjalan sesuai rancangan. Fase 1 membuktikan baseline operasi normal berjalan sukses pada state CLOSED. Pada Fase 2, injeksi fault downstream ditahan sempurna oleh Resilient Client; setelah 3 kegagalan awal, Circuit Breaker membuka (OPEN) dan mengalihkan 10 request ke fallback ('Payment Queued'), menjaga error rate pengguna tetap 0.00%. Fase 3 membuktikan bahwa pemanggilan tanpa mitigasi memicu auto-abort seketika saat error rate mencapai 33.33%, diikuti oleh pemulihan penuh pada Fase 4.",
    narrative_en="This terminal trace verifies that all four experiment phases execute precisely as designed. Phase 1 demonstrates normal baseline execution operating stably in the CLOSED state. In Phase 2, downstream turbulence is absorbed seamlessly by the Resilient Client; after three initial downstream errors, the breaker trips to OPEN, routing all requests to fallback ('Payment Queued') and preserving a 0.00% user error rate. Phase 3 proves that unmitigated traffic triggers an immediate auto-abort when error rates breach 33.33%, followed by pristine recovery in Phase 4.",
    mental_model="Empirical execution proves resilience: zero user-facing errors with fallback, immediate abort without it.",
    source_path="cmd/demo/main.go:50-135 • engineering/03-execution-result.md:50-90",
    slide_num="10 / 24",
    layout="terminal-output"
)

# Slide 11: Analisis batasan bukti
add_slide(
    categories=["BATASAN", "ANALISIS"],
    title="YANG TIDAK DIBUKTIKAN",
    subtitle="Batasan ruang lingkup laboratorium yang wajib dipahami sebelum implementasi produksi berskala besar.",
    content_items=[
        {"title": "YANG TIDAK DIBUKTIKAN", "items": [
            "Injeksi fault dilakukan in-memory Go murni; tidak membuktikan simulasi network packet loss atau kernel drop riil.",
            "Steady-state monitor menggunakan akumulator total kumulatif; bukan sliding time window (misal 1 menit terakhir).",
            "Ambang batas error rate (20%, 25%) dan durasi cooldown (50ms, 200ms) adalah parameter lab, bukan tuning produksi.",
            "Eksperimen tidak menguji latensi persentil p99 dinamis atau degradasi parsial berbasis tracing terdistribusi.",
            "Tidak mencakup canary request routing bertahap (1% -> 5% -> 100%) atau automated rollback berbasis eBPF/service mesh.",
            "Field errorRate di injector.go tidak aktif; pengujian murni menggunakan flag deterministik forceError."
        ], "type": "bullet"}
    ],
    narrative="Meskipun laboratorium ini berhasil membuktikan prinsip-prinsip fundamental Chaos Engineering, terdapat batasan ketat yang harus digarisbawahi. Injeksi kegagalan disimulasikan secara in-memory menggunakan channel dan error buatan, bukan manipulasi soket jaringan tingkat rendah via iptables atau Chaos Mesh. Metrik error rate dihitung secara kumulatif sederhana tanpa sliding window, sehingga histori lama mempengaruhi rasio saat ini. Selain itu, nilai ambang batas waktu dan error dirancang khusus untuk kenyamanan eksekusi test dan bukan merupakan rekomendasi default di sistem live.",
    narrative_en="While this lab successfully validates core Chaos Engineering mechanics, its architectural boundaries must be explicitly noted. Fault injection operates purely in-memory via timer channels and synthetic Go errors, rather than low-level network packet drop simulations via iptables or Envoy proxies. The steady-state monitor aggregates cumulative request counters rather than computing a sliding time window. Furthermore, test thresholds and cooldown intervals were tailored for rapid execution and should not be treated as production baselines.",
    mental_model="Acknowledge lab boundaries: in-memory simulation proves logic; production requires sliding windows and mesh routing.",
    source_path="labs/34-chaos-engineering/content/01-content-brief.md:40-46 • engineering-audit/06-verdict.md:40-42",
    slide_num="11 / 24",
    layout="constraints-list"
)

# Slide 12: Trade-off arsitektur
add_slide(
    categories=["TRADE-OFF", "ARSITEKTUR"],
    title="Trade-off Arsitektur: Fallback vs Beban Sistem",
    subtitle="Menimbang implikasi operasional dari graceful degradation dan pengujian chaos otomatis.",
    content_items=[
        {"title": "KEUNTUNGAN RESILIENSI", "items": [
            "Availability Terjaga: Pengguna tidak menerima error mentah saat dependensi utama terganggu.",
            "Blast Radius Terisolasi: Circuit Breaker menghentikan pemborosan thread dan soket jaringan secara instan.",
            "Kepercayaan Empiris: Perilaku failover terverifikasi secara berkala sebelum menjadi insiden darurat."
        ], "type": "bullet"},
        {"title": "KONSEKUENSI & BIAYA SISTEM", "items": [
            "Eventual Consistency Burden: Fallback berbasis antrean (queue) menuntut rekonsiliasi state asinkron lanjutan.",
            "Stale Data Risk: Fallback berbasis cache lokal berpotensi menyajikan data yang tidak mutakhir.",
            "Blast Radius Overhead: Pemantauan steady-state berkala memerlukan alokasi CPU untuk polling metrik.",
            "Fallback Complexity: Jalur fallback itu sendiri harus diuji agar tidak menjadi titik kegagalan kedua (secondary crash)."
        ], "type": "bullet"}
    ],
    narrative="Penerapan arsitektur berketahanan menuntut pertimbangan trade-off yang matang antara ketersediaan antarmuka dan konsistensi data. Graceful degradation melalui antrean fallback memastikan transaksi pengguna tetap diterima, namun memindahkan beban ke sistem pemrosesan latar belakang yang harus menangani rekonsiliasi data asinkron. Jika fallback menggunakan data cache, pengguna berisiko melihat data usang. Selain itu, jalur fallback itu sendiri harus diperlakukan sebagai kode produksi kritis yang wajib diuji agar tidak memicu kegagalan sekunder.",
    narrative_en="Architecting resilient systems introduces nuanced trade-offs between immediate availability and eventual consistency. Implementing graceful degradation via queuing mechanisms preserves user interactions, yet shifts operational burden onto background reconciliation workers. When fallback relies on local caches, callers risk consuming stale reads. Moreover, fallback handlers represent critical production execution paths that must be thoroughly validated to ensure they do not become single points of secondary failure.",
    mental_model="Fallbacks preserve immediate availability at the cost of eventual data reconciliation complexity.",
    source_path="labs/34-chaos-engineering/content/02-master-draft.md:19-21 • research/05-report.md",
    slide_num="12 / 24",
    layout="tradeoff"
)

# Slide 13: Concurrency / validation pipeline
add_slide(
    categories=["CONCURRENCY", "RACE CONDITIONS"],
    title="Validasi Konkurensi: 20 Worker di Bawah Beban Chaos",
    subtitle="TestConcurrencyAndRace membuktikan sinkronisasi thread-safe bebas race condition dengan go test -race.",
    content_items=[
        {"title": "KODE TEST STRES KONKURENSI", "items": [
            "func TestConcurrencyAndRace(t *testing.T) {",
            "    inj := fault.NewInjector(); mon := monitor.NewMonitor(0.50)",
            "    cb := circuitbreaker.NewCircuitBreaker(5, 20*time.Millisecond)",
            "    var wg sync.WaitGroup; workers := 20; iterations := 50",
            "    for i := 0; i < workers; i++ {",
            "        wg.Add(1)",
            "        go func(id int) {",
            "            defer wg.Done()",
            "            for j := 0; j < iterations; j++ {",
            "                if id%2 == 0 { inj.SetFault(time.Microsecond, true) } else { inj.Clear() }",
            "                _ = cb.Execute(func() error { return inj.Execute(context.Background()) },",
            "                               func() error { return nil })",
            "                mon.RecordSuccess(); _ = mon.IsHealthy()",
            "            }",
            "        }(i)",
            "    }",
            "    wg.Wait()",
            "}"
        ], "type": "code"},
        {"title": "HASIL VERIFIKASI THREAD-SAFETY", "items": [
            "Race Detector: PASS (1.481 detik pengujian bersih tanpa data race).",
            "Primitive Mutex: RWMutex pada Injector dan Mutex pada Circuit Breaker mengisolasi state mutasi.",
            "Atomic Counting: Monitor menggunakan sync/atomic uint64 untuk pembaruan metrik tanpa contention lock."
        ], "type": "bullet"}
    ],
    narrative="Uji konkurensi intensif dijalankan untuk membuktikan ketahanan komponen di bawah kondisi perebutan thread yang ekstrem. Sebanyak 20 goroutine secara paralel melakukan 50 iterasi operasi mutasi fault (`SetFault` dan `Clear`), eksekusi Circuit Breaker, dan pencatatan metrik steady-state. Penggunaan `sync/atomic` pada struct Monitor dan `sync.Mutex` pada Circuit Breaker menjamin konsistensi data memori. Pengujian dengan flag `-race` pada Go compiler membuktikan seluruh sistem 100% bebas dari data race.",
    narrative_en="Rigorous concurrency stress testing verifies component resilience under extreme thread contention. Twenty concurrent goroutines execute 50 iterations each, simultaneously toggling faults via `SetFault` and `Clear`, dispatching requests through the Circuit Breaker, and recording steady-state metrics. Combining `sync/atomic` primitives inside Monitor with exclusive mutex locks inside CircuitBreaker eliminates shared memory corruption. Running Go's race detector flags zero race conditions across 1,000 concurrent cycles.",
    mental_model="Atomic metric updates paired with explicit mutex state barriers guarantee zero race conditions.",
    source_path="tests/chaos_test.go:156-189 • engineering/03-execution-result.md:35-42",
    slide_num="13 / 24",
    layout="concurrency-flow"
)

# Slide 14: Before vs After comparison cards
add_slide(
    categories=["BEFORE/AFTER", "ARSITEKTUR"],
    title="Before vs After: Evolusi Ketahanan Sistem",
    subtitle="Transformasi dari ketergantungan downstream yang rapuh menuju arsitektur berketahanan teruji.",
    content_items=[
        {"title": "BEFORE (Arsitektur Pasif)", "items": [
            "Pola Integrasi: Panggilan langsung tanpa Circuit Breaker atau pembatasan kegagalan.",
            "Reaksi Downstream: Thread pemanggil membeku saat downstream mengalami lonjakan latensi.",
            "Dampak Kegagalan: Resource exhaustion, lonjakan HTTP 500 mentah pada sisi pengguna.",
            "Validasi Ketahanan: Mengandalkan asumsi teoritis dan mock unit test statis yang dangkal."
        ], "type": "card"},
        {"title": "AFTER (Chaos-Engineered Architecture)", "items": [
            "Pola Integrasi: Dibungkus Circuit Breaker dengan rute fallback graceful degradation.",
            "Reaksi Downstream: Fast-fail otomatis memutus sirkuit dan menyajikan respons alternatif terdegradasi.",
            "Dampak Kegagalan: Error rate pengguna bertahan di 0.00% meski downstream mengalami outage.",
            "Validasi Ketahanan: Dibuktikan empiris via injeksi fault terkontrol dan monitor auto-abort."
        ], "type": "card"}
    ],
    narrative="Perbandingan before-after ini mengilustrasikan pergeseran paradigma dari asumsi ketahanan pasif menuju ketahanan aktif yang terbukti empiris. Sebelum perbaikan, sistem sangat rentan terhadap gangguan eksternal karena ketiadaan isolasi kegagalan langsung melumpuhkan thread pemanggil dan memunculkan error 500 ke pengguna. Setelah menerapkan pola Chaos Engineering dengan Circuit Breaker dan Fallback, kegagalan dependensi diisolasi secara instan. Hasilnya, metrik ketersediaan pengguna tetap sempurna bahkan di tengah disrupsi downstream total.",
    narrative_en="This before-and-after comparison highlights the architectural paradigm shift from passive resilience assumptions to empirical, self-healing topologies. Prior to adopting these patterns, systems remained fragile; unshielded downstream outages exhausted caller threads and propagated raw 500 errors to customers. Adopting chaos-engineered defenses with Circuit Breakers and Fallbacks isolates failure domains instantly. As a result, consumer availability metrics remain stable at 100% even during total remote dependency outages.",
    mental_model="Before: fragile direct coupling. After: decoupled execution with circuit protection and graceful fallbacks.",
    source_path="labs/34-chaos-engineering/content/02-master-draft.md:3-26",
    slide_num="14 / 24",
    layout="before-after"
)

# Slide 15: Analisis mendalam - Test 1 Fault Injector Primitives
add_slide(
    categories=["BUKTI TEST", "FAULT INJECTOR"],
    title="Test 1: Presisi Injeksi Latensi & Forced Error",
    subtitle="TestFaultInjector memverifikasi eksekusi latensi, propagasi error, dan netralisasi bersih via Clear().",
    content_items=[
        {"title": "KODE TEST INJEKTOR GANGGUAN", "items": [
            "func TestFaultInjector(t *testing.T) {",
            "    fi := fault.NewInjector(); ctx := context.Background()",
            "    // Verifikasi awal: default disabled",
            "    if err := fi.Execute(ctx); err != nil { t.Fatalf(\"...\") }",
            "    // Injeksi latensi 10ms dan forceError=true",
            "    fi.SetFault(10*time.Millisecond, true)",
            "    start := time.Now(); err := fi.Execute(ctx); elapsed := time.Since(start)",
            "    if !errors.Is(err, fault.ErrInjectedFault) || elapsed < 10*time.Millisecond {",
            "        t.Fatalf(\"injeksi gagal memenuhi syarat latensi atau error\")",
            "    }",
            "    // Netralisasi gangguan",
            "    fi.Clear()",
            "    if fi.IsEnabled() || fi.Execute(ctx) != nil { t.Fatalf(\"...\") }",
            "}"
        ], "type": "code"},
        {"title": "POIN VERIFIKASI INVARIAN", "items": [
            "Presisi Waktu: Latensi terukur elapsed >= 10 milidetik sesuai durasi yang disuntikkan.",
            "Error Determinism: Mengembalikan sentry error ErrInjectedFault saat forceError diaktifkan.",
            "Clean Teardown: Pemanggilan Clear() mengembalikan state isEnabled ke false dan error kembali nil."
        ], "type": "bullet"}
    ],
    narrative="Unit test pertama memverifikasi ketepatan primitif Fault Injector dalam menyisipkan gangguan terprogram. Saat berada pada state non-aktif, injektor tidak memberikan overhead dan mengembalikan nilai nil. Ketika dikonfigurasi dengan latensi 10 milidetik dan forced error, fungsi `Execute` menahan thread pemanggil secara akurat selama minimal 10 milidetik sebelum mengembalikan `ErrInjectedFault`. Begitu method `Clear()` dipanggil, seluruh konfigurasi gangguan terhapus secara atomik dan pemanggilan berikutnya kembali bersih seketika.",
    narrative_en="The first unit test verifies the surgical precision of the Fault Injector primitive when introducing simulated turbulence. In its default disabled state, the injector incurs zero operational overhead and returns nil. When programmed with 10ms latency and forced error flags, `Execute` holds execution for at least 10ms before emitting `ErrInjectedFault`. Calling `Clear()` atomically wipes the fault configuration, ensuring subsequent executions pass cleanly without residual effects.",
    mental_model="Fault injection primitives must be deterministic, timed precisely, and neutralized instantly on demand.",
    source_path="tests/chaos_test.go:16-50",
    slide_num="15 / 24",
    layout="test-analysis"
)

# Slide 16: Analisis mendalam - Test 2 Circuit Breaker State Transitions
add_slide(
    categories=["BUKTI TEST", "CIRCUIT BREAKER"],
    title="Test 2: Transisi State Closed → Open → Half-Open",
    subtitle="TestCircuitBreakerStateTransitions memvalidasi threshold kegagalan, fast-fail, dan pemulihan.",
    content_items=[
        {"title": "VERIFIKASI TRANSISI TIGA STATE", "items": [
            "// Threshold = 2 failures, Cooldown = 50ms",
            "cb := circuitbreaker.NewCircuitBreaker(2, 50*time.Millisecond)",
            "_ = cb.Execute(failFn, nil) // Failure #1 -> Masih CLOSED",
            "_ = cb.Execute(failFn, nil) // Failure #2 -> Trip ke OPEN!",
            "",
            "// Pemanggilan instan menghasilkan fast-fail ErrCircuitOpen tanpa memanggil fn",
            "err := cb.Execute(okFn, nil)",
            "// Assert: errors.Is(err, circuitbreaker.ErrCircuitOpen)",
            "",
            "time.Sleep(60 * time.Millisecond) // Lewati cooldown 50ms",
            "// State otomatis bergeser ke HALF-OPEN untuk uji coba",
            "err = cb.Execute(okFn, nil) // Eksekusi probasi SUKSES",
            "// Assert: cb.State() kembali ke CLOSED dan failure counter direset ke 0"
        ], "type": "code"},
        {"title": "GARANSI PERILAKU SIRKUIT", "items": [
            "Threshold Boundary: Sirkuit tetap tertutup pada kegagalan pertama, dan baru terbuka tepat di kegagalan kedua.",
            "Fast Failure: State OPEN mengembalikan ErrCircuitOpen tanpa menyentuh layanan downstream.",
            "Self-Healing: Eksekusi sukses pada fase HALF-OPEN secara otomatis mereset sistem ke CLOSED."
        ], "type": "bullet"}
    ],
    narrative="Unit test kedua membuktikan kepatuhan transisi state machine pada modul Circuit Breaker. Sirkuit dikonfigurasi dengan ambang batas 2 kegagalan dan cooldown 50 milidetik. Kegagalan pertama mempertahankan status CLOSED. Kegagalan kedua langsung memicu transisi sirkuit ke OPEN. Pada status OPEN, pemanggilan fungsi langsung diputus (fast-fail) dengan error `ErrCircuitOpen`. Setelah jeda 60 milidetik melewati masa cooldown, sirkuit masuk ke mode HALF-OPEN; pemanggilan sukses pada masa percobaan ini langsung mereset sirkuit ke CLOSED.",
    narrative_en="The second unit test validates the state machine invariants of the Circuit Breaker module. Configured with a failure threshold of 2 and a 50ms cooldown window, the initial error preserves the CLOSED state. The second consecutive error instantly trips the breaker to OPEN. Under OPEN, execution calls fast-fail immediately with `ErrCircuitOpen`. After waiting 60ms to exceed cooldown, the breaker transitions to HALF-OPEN; a subsequent successful invocation cleanly resets the state to CLOSED.",
    mental_model="Trip on threshold breach, bypass downstream during open, and heal automatically on probation success.",
    source_path="tests/chaos_test.go:52-90",
    slide_num="16 / 24",
    layout="test-analysis"
)

# Slide 17: Analisis mendalam - Test 3 Graceful Degradation
add_slide(
    categories=["BUKTI TEST", "FALLBACK DEGRADATION"],
    title="Test 3: Graceful Degradation Menelan Error Primary",
    subtitle="TestCircuitBreakerGracefulDegradation memverifikasi proteksi pengalaman pengguna via fallback.",
    content_items=[
        {"title": "UJI KONTRAK FALLBACK", "items": [
            "func TestCircuitBreakerGracefulDegradation(t *testing.T) {",
            "    cb := circuitbreaker.NewCircuitBreaker(1, 50*time.Millisecond)",
            "    failFn := func() error { return errors.New(\"upstream failed\") }",
            "",
            "    fallbackInvoked := false",
            "    fallback := func() error {",
            "        fallbackInvoked = true",
            "        return nil // Menelan error dan menyajikan degraded response",
            "    }",
            "",
            "    err := cb.Execute(failFn, fallback)",
            "    // Assert: err == nil (error ditelan)",
            "    // Assert: fallbackInvoked == true"
            "}"
        ], "type": "code"},
        {"title": "ANALISIS MEKANISME FALLBACK", "items": [
            "Error Masking: Error kegagalan remote ditelan oleh fallback handler, mencegah status 500.",
            "Alternative Payload: Klien menerima respons terdegradasi yang valid (contoh: antrean pembayaran asinkron).",
            "Upstream Shield: Konsumen API tidak mengalami interupsi transaksi meskipun dependensi downstream padam."
        ], "type": "bullet"}
    ],
    narrative="Unit test ketiga menguji efektivitas fungsi fallback dalam menerapkan pola graceful degradation. Ketika fungsi utama `failFn` mengembalikan error, Circuit Breaker mengeksekusi penanganan alternatif `fallback` yang telah disediakan. Handler fallback berhasil menelan error remote dan mengembalikan nilai nil ke klien pemanggil. Mekanisme ini menjamin bahwa kegagalan downstream tidak pernah bocor menjadi error mentah di sisi pengguna, melainkan dialihkan menjadi respons terdegradasi yang sah dan terkontrol.",
    narrative_en="The third unit test exercises graceful degradation via fallback handlers. When the primary function `failFn` returns an error, the Circuit Breaker executes the alternative `fallback` closure. The fallback handler successfully intercepts the remote failure and returns nil to the caller. This pattern ensures that downstream service outages never bubble up as raw customer errors, substituting transparent failures with structured, controlled degradation responses.",
    mental_model="Fallback functions absorb raw dependency failures, returning gracefully degraded responses to users.",
    source_path="tests/chaos_test.go:92-109",
    slide_num="17 / 24",
    layout="test-analysis"
)

# Slide 18: Analisis mendalam - Test 4 Auto-Abort Blast Radius
add_slide(
    categories=["BUKTI TEST", "BLAST RADIUS AUTO-ABORT"],
    title="Test 4: Auto-Abort & Netralisasi Saat Pelanggaran Metrik",
    subtitle="TestExperimentAutoAbortOnSteadyStateViolation membuktikan penahanan dampak eksperimen secara otomatis.",
    content_items=[
        {"title": "SIMULASI AUTO-ABORT PADA STEADY-STATE BREACH", "items": [
            "inj := fault.NewInjector(); mon := monitor.NewMonitor(0.20) // Batas 20%",
            "// Baseline 4 request sukses",
            "for i := 0; i < 4; i++ { mon.RecordSuccess() }",
            "exp := experiment.NewExperiment(cfg, inj, mon)",
            "go func() { errCh <- exp.Run(ctx) }()",
            "",
            "// Suntikkan kegagalan langsung ke monitor (error rate = 2/6 = 33.3%)",
            "mon.RecordFailure(); mon.RecordFailure()",
            "",
            "// Assert: err != nil (\"Steady state breached: error rate 33.33%...\")",
            "// Assert: exp.State() == experiment.StateAborted",
            "// Assert: inj.IsEnabled() == false (dinutralkan seketika!)"
        ], "type": "code"},
        {"title": "FAKTA KONTROL BLAST RADIUS", "items": [
            "Proactive Termination: Background ticker mendeteksi pelanggaran dalam interval 10 milidetik.",
            "State Transition: Eksperimen beralih ke StateAborted dengan pesan alasan yang jelas dan terstruktur.",
            "Synchronous Teardown: Injektor dinonaktifkan seketika, mencegah dampak meluas ke request berikutnya."
        ], "type": "bullet"}
    ],
    narrative="Unit test keempat membuktikan penahanan blast radius yang menjadi syarat mutlak dalam Chaos Engineering profesional. Eksperimen dikonfigurasi dengan batas toleransi error rate 20%. Ketika dua kegagalan disuntikkan sehingga rasio error melonjak menjadi 33.33%, background ticker pemantau mendeteksi pelanggaran tersebut dan langsung membatalkan eksperimen dengan status StateAborted. Yang terpenting, method `inj.IsEnabled()` terbukti bernilai false, mengonfirmasi bahwa seluruh gangguan dinetralkan seketika tanpa kebocoran efek samping.",
    narrative_en="The fourth test verifies blast-radius containment, a mandatory prerequisite for production-grade Chaos Engineering. Configured with a 20% error rate tolerance, the experiment is subjected to failures pushing the error ratio to 33.33%. The background polling ticker catches the boundary breach and immediately aborts the run under StateAborted. Crucially, asserting `inj.IsEnabled() == false` confirms that all fault injection primitives are synchronously dismantled without leaking into subsequent requests.",
    mental_model="Automated abort mechanisms act as emergency circuit breakers for the chaos experiment itself.",
    source_path="tests/chaos_test.go:111-154",
    slide_num="18 / 24",
    layout="test-analysis"
)

# Slide 19: Studi kasus industri & best practices
add_slide(
    categories=["STUDI KASUS", "BEST PRACTICES"],
    title="Penerapan Industri: Dari Simian Army ke Production",
    subtitle="Bagaimana raksasa teknologi mengelola blast radius dan mengotomatiskan pembuktian ketahanan.",
    content_items=[
        {"title": "EVOLUSI CHAOS DI INDUSTRI", "items": [
            "Netflix Chaos Monkey & FIT: Menyuntikkan kegagalan langsung pada tingkat RPC protokol dengan kontrol context.",
            "AWS Fault Injection Simulator (FIS): Framework terkelola untuk pengujian disrupsi instance dan partisi jaringan.",
            "Google Disaster Recovery Testing (DiRT): Latihan berskala besar untuk memvalidasi kesiapan failover multi-region."
        ], "type": "bullet"},
        {"title": "PRINSIP EMAS BLAST RADIUS", "items": [
            "Start in Staging / Canary: Mulai dari lingkungan pengujian atau fraksi trafik terkecil (1% pengguna).",
            "Automated Kill Switch: Jangan pernah menjalankan eksperimen chaos tanpa pemantau metrik real-time dan auto-abort.",
            "Business Metric Alignment: Evaluasi steady-state berdasarkan throughput checkout dan checkout conversion."
        ], "type": "bullet"}
    ],
    narrative="Dalam industri software terdistribusi, Chaos Engineering telah berevolusi dari sekadar skrip pemati server acak menjadi disiplin rekayasa yang sangat presisi. Netflix memelopori pendekatan ini dengan Failure Injection Testing (FIT) yang menyisipkan gangguan pada level RPC menggunakan konteks request spesifik. AWS dan Google mengintegrasikan pengujian kegagalan terkelola dengan batas blast radius yang ketat. Kunci sukses implementasi industri adalah mengikat hipotesis eksperimen langsung ke metrik bisnis utama dan selalu menyertakan tombol pembatalan darurat otomatis.",
    narrative_en="Across modern enterprise systems, Chaos Engineering has evolved from random instance termination into a highly disciplined scientific practice. Netflix pioneered fine-grained failure injection via FIT, injecting latency at the RPC layer using synthetic request tracing headers. AWS and Google mandate automated blast-radius containment for cross-region recovery exercises. The universal golden rule across industry leaders is anchoring steady-state hypotheses directly to core business metrics backed by automated kill switches.",
    mental_model="Industrial chaos engineering mandates strict canary blast radius and automated kill switches.",
    source_path="labs/34-chaos-engineering/content/01-content-brief.md:44-46 • research/02-sources.md",
    slide_num="19 / 24",
    layout="cards"
)

# Slide 20: Checklist implementasi
add_slide(
    categories=["CHECKLIST", "PRODUKSI"],
    title="Checklist Implementasi Chaos Engineering",
    subtitle="8 langkah verifikasi wajib sebelum menggelar pengujian kegagalan terkontrol ke sistem produksi.",
    content_items=[
        {"title": "CHECKLIST KESIAPAN EKSPERIMEN", "items": [
            "[ ] Tentukan steady-state metric berbasis output pengguna (error rate, p99 latency) dengan batas jelas.",
            "[ ] Tetapkan ukuran sampel minimum sebelum evaluasi breach agar terhindar dari alarm palsu startup.",
            "[ ] Bungkus semua pemanggilan remote dengan Circuit Breaker yang memiliki parameter threshold dan cooldown.",
            "[ ] Sediakan rute fallback graceful degradation (antrean asinkron atau cache lokal) yang teruji.",
            "[ ] Pasang background monitor yang mengevaluasi steady-state metric secara periodik via interval ticker.",
            "[ ] Hubungkan monitor dengan mekanisme auto-abort dan pembersihan injector sinkron (injector.Clear()).",
            "[ ] Batasi durasi maksimum eksperimen (timeout window) untuk mencegah eksekusi tak terhingga.",
            "[ ] Uji seluruh komponen di bawah race detector Go (go test -race ./...) untuk menjamin thread-safety."
        ], "type": "bullet"}
    ],
    narrative="Daftar periksa ini merangkum delapan langkah penting yang wajib dipenuhi sebelum menggelar eksperimen chaos ke lingkungan produksi. Eksperimen tidak boleh dimulai tanpa definisi metrik steady-state yang terikat pada pengalaman pengguna serta batas toleransi yang terukur. Pemanggilan dependensi eksternal harus diproteksi oleh Circuit Breaker dan rute fallback alternatif. Terakhir, mekanisme pembatalan otomatis (auto-abort) dan netralisasi injektor harus diverifikasi secara ketat agar tidak ada gangguan yang tertinggal di sistem.",
    narrative_en="This checklist codifies eight mandatory architectural requirements before scheduling chaos experiments in production. No experiment may begin without defining user-facing steady-state metrics paired with concrete threshold boundaries. External dependencies must be guarded behind Circuit Breakers backed by verified fallback degradation routes. Finally, automated abort triggers and synchronous injector teardown routines must be validated under concurrency race detectors to guarantee zero residual turbulence.",
    mental_model="Define metrics, protect with breakers, enforce fallbacks, and guarantee automated abort teardown.",
    source_path="labs/34-chaos-engineering/content/02-master-draft.md:90-96",
    slide_num="20 / 24",
    layout="checklist"
)

# Slide 21: Anti-pattern & cara reproduksi kegagalan
add_slide(
    categories=["ANTI-PATTERN", "REPRODUKSI"],
    title="Anti-Pattern: Kesalahan Fatal dalam Chaos Testing",
    subtitle="Kekeliruan desain yang mengubah pengujian ketahanan menjadi malapetaka production.",
    content_items=[
        {"title": "ANTI-PATTERN UMUM", "items": [
            "Chaos Tanpa Auto-Abort: Menyuntikkan kegagalan tanpa pemantauan kontinu dan tanpa tombol darurat pemutus.",
            "Memantau Metrik Internal Saja: Berfokus pada utilisasi CPU/RAM alih-alih tingkat kegagalan transaksi pengguna.",
            "Circuit Breaker Tanpa Fallback: Sirkuit terbuka memutus panggilan tetapi pengguna tetap menerima error 500 mentah.",
            "Untested Fallback Logic: Fallback handler memicu panic atau deadlock baru saat dieksekusi di bawah tekanan."
        ], "type": "card"},
        {"title": "CARA REPRODUKSI GAGAL DI LAB", "items": [
            "1. Ubah cmd/demo/main.go dengan menghapus blok fallback pada pemanggilan clientRequest.",
            "2. Jalankan go run ./cmd/demo: saat downstream disuntikkan fault, error rate langsung meledak ke 100%.",
            "3. Matikan auto-abort di internal/experiment/runner.go: eksperimen terus menyiksa sistem tanpa henti!"
        ], "type": "card"}
    ],
    narrative="Mengabaikan prinsip dasar Chaos Engineering dapat berakibat fatal bagi reputasi dan ketersediaan layanan perusahaan. Anti-pattern paling berbahaya adalah menjalankan injeksi gangguan tanpa monitor steady-state yang dilengkapi auto-abort otomatis, membiarkan eksperimen merusak trafik pengguna secara tak terbatas. Kesalahan umum lainnya adalah memasang Circuit Breaker tanpa menyiapkan fungsi fallback, sehingga pengguna tetap mengalami kegagalan transaksi. Untuk membuktikannya di lab, kita cukup mematikan fungsi fallback dan melihat error rate pengguna melonjak drastis.",
    narrative_en="Violating chaos engineering principles introduces grave operational risks. The most dangerous anti-pattern is executing failure injection without continuous steady-state monitoring and automated abort mechanisms, leaving experiments to ravage user traffic indefinitely. Another common flaw is implementing Circuit Breakers without fallback handlers, resulting in fast-fail errors that still crash client workflows. Disabling fallbacks in the lab immediately reproduces this vulnerability, causing user error rates to spike instantly to 100%.",
    mental_model="Never inject faults without auto-abort; never trip breakers without graceful fallbacks.",
    source_path="labs/34-chaos-engineering/content/02-master-draft.md:78-82",
    slide_num="21 / 24",
    layout="anti-pattern-cards"
)

# Slide 22: Takeaway praktis & command block dominan
add_slide(
    categories=["PRAKTIS", "COMMAND"],
    title="Verifikasi Mandiri: Uji & Jalankan Lab",
    subtitle="Reproduksi seluruh bukti teknis, benchmark race detector, dan demo interaktif di terminal Anda.",
    content_items=[
        {"title": "PERINTAH VERIFIKASI RESMI", "items": [
            "cd labs/34-chaos-engineering",
            "go test -v ./...",
            "go test -race ./...",
            "go run ./cmd/demo"
        ], "type": "code"},
        {"title": "RINGKASAN METRIK LAB", "items": [
            "Test Suite: 5/5 test cases lulus 100% (FaultInjector, CircuitBreaker, Degradation, AutoAbort, Race).",
            "Race Detector: Bersih dari data race di bawah 20 concurrent goroutines (1.481 detik pengujian).",
            "Pustaka Mandiri: Murni Go Standard Library (sync, sync/atomic, time, context, errors) tanpa third-party."
        ], "type": "bullet"}
    ],
    narrative="Laboratorium ini dirancang untuk dapat diverifikasi secara mandiri dalam hitungan detik di mesin lokal mana pun. Seluruh kode dibangun murni menggunakan standard library Go tanpa memerlukan instalasi database eksternal, container Docker, atau framework pengujian pihak ketiga. Perintah `go test -v` menjalankan lima rangkaian pengujian komprehensif mulai dari injeksi primitif hingga auto-abort. Perintah `go test -race` memvalidasi thread-safety konkurensi, dan `go run ./cmd/demo` menyajikan simulasi interaktif dari seluruh siklus ketahanan.",
    narrative_en="This laboratory is engineered for immediate, deterministic local verification on any workstation. The entire codebase is implemented with the Go standard library, completely free of external database daemons, Docker dependencies, or third-party test runners. Running `go test -v` executes five comprehensive suites validating injection primitives, state transitions, fallbacks, and auto-abort triggers. Running `go test -race` guarantees thread-safety under concurrency, while `go run ./cmd/demo` delivers an interactive terminal simulation.",
    mental_model="Zero-dependency standard library code enables immediate, deterministic local resilience verification.",
    source_path="labs/34-chaos-engineering/README.md:30-45 • engineering/03-execution-result.md",
    slide_num="22 / 24",
    layout="command-block"
)

# Slide 23: Prinsip arsitektural jangka panjang
add_slide(
    categories=["PRINSIP", "ARSITEKTUR"],
    title="Prinsip Arsitektur: Desain untuk Turbulensi",
    subtitle="Pelajaran abadi dalam merancang sistem terdistribusi yang tangguh terhadap kegagalan parsial.",
    content_items=[
        {"title": "1. FAILURE IS INEVITABLE; COLLAPSE IS A CHOICE", "items": [
            "Layanan downstream pasti akan mengalami degradasi performa atau pemadaman total.",
            "Arsitektur yang baik tidak berharap sistem kebal, melainkan merancang isolasi agar kegagalan tidak merembet."
        ], "type": "card"},
        {"title": "2. USER STEADY STATE IS THE ONLY TRUTH", "items": [
            "Kesehatan internal komponen (CPU, thread pool, memori) tidak berarti apa-apa jika transaksi pengguna gagal.",
            "Jadikan metrik output batas konsumen sebagai kompas utama dalam mengevaluasi keandalan sistem."
        ], "type": "card"},
        {"title": "3. BLAST RADIUS IS SACROSANCT", "items": [
            "Pengujian ketahanan tidak boleh menjadi sumber insiden baru bagi pengguna produksi.",
            "Otomatisasi pembatalan eksperimen (auto-abort) dan degradasi anggun (fallback) adalah pengaman mutlak."
        ], "type": "card"}
    ],
    narrative="Tiga prinsip arsitektur ini berlaku universal di seluruh perancangan sistem terdistribusi modern. Pertama, kegagalan komponen eksternal adalah kepastian hukum alam dalam sistem terdistribusi; yang membedakan sistem rapuh dari sistem tangguh adalah kemampuannya mengisolasi kegagalan agar tidak menjadi pemadaman total. Kedua, selalu jadikan metrik output pengguna sebagai tolok ukur utama ketersediaan sistem. Ketiga, kesucian batasan dampak (blast radius) harus selalu dijaga melalui otomasi pemantauan dan tombol pembatalan darurat.",
    narrative_en="These three architectural principles transcend this specific Go implementation and apply universally to modern software engineering. First, remote dependency degradation is an absolute mathematical certainty; resilient architectures isolate failures rather than pretending they will never occur. Second, customer-facing steady-state metrics represent the ultimate ground truth of system health. Third, the sanctity of blast-radius containment must be fiercely preserved through automated kill switches and graceful fallbacks.",
    mental_model="Embrace failure through isolation, measure availability at the user boundary, and strictly bound blast radius.",
    source_path="labs/34-chaos-engineering/content/05-key-takeaways.md",
    slide_num="23 / 24",
    layout="principles-list"
)

# Slide 24: Closing slide
add_slide(
    categories=["CLOSING"],
    title="Ketahanan Terdistribusi Terbukti, Blast Radius Terkendali.",
    subtitle="",
    content_items=[
        {"title": "KOMPONEN ARSITEKTUR LAB", "items": [
            "internal/fault/injector.go — Thread-safe latency and forced error fault injection primitive",
            "internal/circuitbreaker/circuitbreaker.go — Closed, Open, Half-Open state machine with fallback",
            "internal/monitor/monitor.go — Atomic steady-state metrics tracker & breach evaluator",
            "internal/experiment/runner.go — Experiment lifecycle coordinator with auto-abort & synchronous clear",
            "tests/chaos_test.go — 5 exhaustive suites covering injection, states, degradation, abort, and race",
            "cmd/demo/main.go — Complete interactive lifecycle simulation with payment gateway example"
        ], "type": "bullet"},
        {"title": "VERIFIKASI & REPOSITORI", "items": [
            f"Commit SHA: {COMMIT_SHA[:8]}",
            "Unit Tests: 5/5 PASS (100% clean across all packages)",
            "Race Detector: PASS (zero data race detected)",
            f"Repositori: github.com/lukman-ss/software-engineering-lab/tree/main/labs/{LAB_NO}-{SLUG}"
        ], "type": "bullet"}
    ],
    narrative="Laboratorium Chaos Engineering & Fault Injection ini membuktikan secara konkret bahwa ketahanan sistem terdistribusi bukanlah teori abstrak, melainkan properti empiris yang dapat dibuktikan dan diukur. Melalui kombinasi injeksi kegagalan terkontrol, Circuit Breaker adaptif, rute fallback graceful degradation, dan steady-state monitor yang dilengkapi auto-abort blast radius, arsitektur microservices terbukti kebal terhadap disrupsi downstream. Seluruh kode sumber, dokumentasi audit, dan prosedur pengujian tersedia lengkap di repositori.",
    narrative_en="This Chaos Engineering & Fault Injection laboratory provides definitive proof that distributed resilience is not an abstract theoretical ideal, but an empirically verifiable discipline. Harmonizing controlled fault injection primitives, adaptive Circuit Breakers, graceful fallback degradation, and steady-state monitors backed by automated blast-radius containment ensures systems survive external disruptions unscathed. Complete source code, architectural audit trails, and verification scripts are available in the repository.",
    mental_model="Steady state defines stability; circuit breakers isolate faults; auto-abort guarantees safe experimentation.",
    source_path=REPO_URL,
    slide_num="24 / 24",
    layout="closing-dark-navy"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
