#!/usr/bin/env python3
"""
Publication Definition for Software Engineering Lab 34: Chaos Engineering & Fault Injection.
Defines 24 slides, metadata, colors, mental models, and content items.
DO NOT DELETE THIS FILE - reproducible builder source.
"""
import os
import sys

OUTPUT_DIR = "/Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/34-chaos-engineering/result"
CHROME_BIN = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/34-chaos-engineering"
COMMIT_SHA = "b41fffba81521102d1ddfa66c51ba9b4cdae5758"
LAB_NO = "34"
SLUG = "chaos-engineering"
TITLE_ID = "Chaos Engineering & Fault Injection"
SUBTITLE_ID = "Membuktikan Ketahanan Sistem Terdistribusi Secara Empiris"
TITLE_EN = "Chaos Engineering & Fault Injection"
SUBTITLE_EN = "Empirically Proving Distributed System Resilience"
AUTHOR = "Lukman / lukman-ss"

# Standard Brand Colors
DARK_NAVY = "0b1b2d"
TEAL = "087f76"
CREAM = "f4f1e8"
WHITE = "ffffff"
SLATE = "52606b"
ORANGE = "d97706"

SLIDES = []

def add_slide(categories, title, subtitle, content_items, mental_model, source_path="", slide_num=""):
    """Add a slide to the presentation definition."""
    SLIDES.append({
        "categories": categories,
        "title": title,
        "subtitle": subtitle,
        "content_items": content_items,
        "mental_model": mental_model,
        "source": source_path,
        "slide_num": slide_num
    })

# Slide 1: Cover (Dark Navy, Visual CRT Element)
add_slide(
    categories=["KONTEKS"],
    title=TITLE_ID,
    subtitle="SOFTWARE ENGINEERING LAB | LAB " + LAB_NO,
    content_items=[
        {"type": "cover", "title": "", "items": [
            "Retro CRT monitor dengan scanlines dan glow terminal",
            "Diagram alur kegagalan downstream dan circuit breaker",
            f"Penulis: {AUTHOR}",
        ]}
    ],
    mental_model="",
    source_path="",
    slide_num=""
)

# Slide 2: Problem & Motivation
add_slide(
    categories=["MOTIVASI", "PROBLEMATIKA"],
    title="Ketahanan Sistem Tanpa Bukti Empiris Adalah Asumsi",
    subtitle="Layanan terdistribusi rentan ambruk saat downstream mengalami latensi atau outage mendadak",
    content_items=[
        {"title": "MASALAH UTAMA", "items": [
            "Klaim ketahanan arsitektur microservice sering kali hanya diuji pada happy path dan mock unit test sederhana.",
            "Saat ketergantungan downstream melambat atau gagal, ketiadaan perlindungan memicu antrean thread dan resource exhaustion.",
            "Tanpa pengujian empiris terencana, insiden production pertama menjadi tempat uji coba yang mahal bagi reputasi.",
        ], "type": "bullet"},
        {"title": "BAHAYA CASCADING", "items": [
            "Latensi downstream 500ms dapat mengunci connection pool HTTP client dalam hitungan detik.",
            "Kegagalan satu payment gateway merembet ke checkout, cart, hingga mematikan API gateway hulu.",
            "Pengujian integrasi statis tidak mampu mensimulasikan kegagalan jaringan yang fluktuatif.",
        ], "type": "bullet"},
        {"title": "SOLUSI LAB", "items": [
            "Membangun framework eksperimen chaos deterministik di Go 1.22+ murni tanpa external library.",
            "Menyuntikkan latensi dan error downstream secara terprogram di memori.",
            "Membuktikan isolasi kegagalan lewat Circuit Breaker dan graceful degradation fallback.",
        ], "type": "bullet"},
    ],
    mental_model="Resilience without empirical verification is merely wishful thinking.",
    source_path="labs/34-chaos-engineering/README.md",
    slide_num="2 / 24"
)

# Slide 3: Principles of Chaos Engineering
add_slide(
    categories=["PRINSIP", "METODOLOGI"],
    title="Empat Langkah Baku Eksperimen Ilmiah Perangkat Lunak",
    subtitle="Chaos engineering merupakan disiplin verifikasi terstruktur, bukan aksi merusak server secara acak",
    content_items=[
        {"title": "1. STEADY STATE", "items": [
            "Definisikan metrik kesehatan output yang dirasakan pengguna, seperti rasio error dan throughput transaksi.",
            "Fokus pada metrik eksternal bisnis, bukan metrik internal sistem seperti beban CPU atau penggunaan RAM node.",
            "Tentukan baseline angka normal sebelum eksperimen dijalankan sebagai acuan toleransi.",
        ], "type": "bullet"},
        {"title": "2. HIPOTESIS & INJEKSI", "items": [
            "Rumuskan hipotesis: 'Steady-state error rate tetap di bawah ambang batas meskipun downstream mati'.",
            "Suntikkan gangguan terkontrol: latensi jaringan tiruan atau pemaksaan error status downstream.",
            "Jalankan injeksi pada durasi terukur dengan interval observasi metrik periodik.",
        ], "type": "bullet"},
        {"title": "3. BLAST RADIUS & ABORT", "items": [
            "Tetapkan safety guardrails ketat: auto-abort langsung aktif saat metrik melanggar threshold toleransi.",
            "Sinkronisasi pembersihan fault injector seketika untuk mengembalikan trafik normal.",
            "Buktikan bahwa sistem memulihkan diri secara mandiri tanpa intervensi manual tim engineer.",
        ], "type": "bullet"},
    ],
    mental_model="Define steady state first, inject controlled turbulence, and abort automatically on breach.",
    source_path="labs/34-chaos-engineering/content/02-master-draft.md",
    slide_num="3 / 24"
)

# Slide 4: Steady-State Metric Definition
add_slide(
    categories=["METRIK", "OBSERVABILITAS"],
    title="Steady State: Mengukur Dampak Pengguna, Bukan Node",
    subtitle="Kesehatan arsitektur diukur dari error rate transaksi akhir, bukan persentase utilitas hardware",
    content_items=[
        {"title": "METRIK OUTPUT VS INTERNAL", "items": [
            "Metrik internal (CPU 80%, RAM 4GB) sering menipu: node sibuk belum tentu melayani pengguna dengan benar.",
            "Metrik output (Error Rate, Order Completion Rate) langsung mencerminkan kepuasan konsumen.",
            "Lab menetapkan ambang batas error rate maksimum (contoh: 20% - 25%) sebagai batas steady state.",
        ], "type": "bullet"},
        {"title": "SAMPLE WARMUP MINIMUM", "items": [
            "Evaluasi kesehatan membutuhkan minimal 5 sampel request awal sebelum perhitungan error rate dieksekusi.",
            "Mencegah false alarm saat request pertama kebetulan gagal pada saat inisialisasi koneksi.",
            "Memastikan rasio error dihitung berdasarkan data agregat yang cukup representatif.",
        ], "type": "bullet"},
        {"title": "STRUKTUR DATA METRIK", "items": [
            "TotalRequests: akumulasi seluruh request masuk yang dilayani client.",
            "FailedRequests: jumlah request yang berujung kegagalan atau error.",
            "SuccessRequests: request yang berhasil diselesaikan via jalur utama maupun fallback.",
        ], "type": "bullet"},
    ],
    mental_model="Measure customer-facing output metrics to define health, never component internals.",
    source_path="labs/34-chaos-engineering/internal/monitor/monitor.go:8",
    slide_num="4 / 24"
)

# Slide 5: System Architecture & Closed-Loop Control
add_slide(
    categories=["ARSITEKTUR", "DESAIN SISTEM"],
    title="Topologi Komponen dan Siklus Closed-Loop Kontrol",
    subtitle="Pemisahan tegas antara injeksi kegagalan, mitigasi client, monitor metrik, dan abort runner",
    content_items=[
        {"title": "DIAGRAM ALUR SISTEM", "items": [
            "[Client Request] --> [Resilient Client (Circuit Breaker)]",
            "                  |--> Downstream Service (with Injector)",
            "                  |--> Graceful Fallback (on Failure / Open)",
            "[Steady State Monitor] <-- (Records Success / Failure)",
            "[Experiment Runner] checks Monitor Ticker --> Triggers Auto-Abort",
        ], "type": "diagram"},
        {"title": "PEMBAGIAN KOMPONEN", "items": [
            "internal/fault: Mengatur latensi buatan dan pemaksaan error downstream secara thread-safe.",
            "internal/circuitbreaker: State machine 3 fasa (Closed, Open, Half-Open) pengendali akses.",
            "internal/monitor: Pelacak metrik rasio kegagalan berbasis operasi atomik berkecepatan tinggi.",
            "internal/experiment: Orkestrator siklus hidup pengujian dengan kemampuan auto-abort sinkron.",
        ], "type": "bullet"},
        {"title": "CLOSED-LOOP FEEDBACK", "items": [
            "Runner tidak berjalan buta; runner membaca kondisi monitor secara periodik via background ticker.",
            "Jika monitor melaporkan status tidak sehat, runner seketika menetralkan injector downstream.",
            "Mengeliminasi risiko eksperimen tertinggal yang merusak trafik production secara permanen.",
        ], "type": "bullet"},
    ],
    mental_model="A closed-loop monitor guarantees turbulence stops before customer harm becomes permanent.",
    source_path="labs/34-chaos-engineering/content/04-diagrams.md",
    slide_num="5 / 24"
)

# Slide 6: Fault Injector Mechanism
add_slide(
    categories=["BUKTI KODE", "FAULT INJECTION"],
    title="Injeksi Fault: Simulasi Latensi dan Error Terprogram",
    subtitle="Implementasi thread-safe menggunakan sync.RWMutex dan time.After dengan pembatalan konteks",
    content_items=[
        {"title": "KODE INJECTOR", "items": [
            "func (fi *Injector) Execute(ctx context.Context) error {",
            "    fi.mu.RLock()",
            "    enabled, latency, forceErr := fi.enabled, fi.latency, fi.forceError",
            "    fi.mu.RUnlock()",
            "    if !enabled { return nil }",
            "    if latency > 0 {",
            "        select {",
            "        case <-time.After(latency):",
            "        case <-ctx.Done(): return ctx.Err()",
            "        }",
            "    }",
            "    if forceErr { return ErrInjectedFault }",
            "    return nil",
            "}",
        ], "type": "code"},
        {"title": "DUKUNGAN CONTEXT", "items": [
            "Latensi injeksi mendengarkan ctx.Done() agar request timeout client langsung dibatalkan.",
            "Mencegah goroutine leak saat client membatalkan transaksi di tengah penundaan waktu.",
            "Mendukung timeout standar library context Go secara alami tanpa blocking abadi.",
        ], "type": "bullet"},
        {"title": "MUTEX SYNCHRONIZATION", "items": [
            "Operasi SetFault() dan Clear() menggunakan Lock eksklusif untuk mutasi state instan.",
            "Operasi Execute() dan IsEnabled() menggunakan RLock untuk throughput pembacaan tinggi.",
            "Terbukti bebas race condition di bawah pengujian go test -race.",
        ], "type": "bullet"},
    ],
    mental_model="Inject faults at system boundaries with context awareness to prevent goroutine leaks.",
    source_path="labs/34-chaos-engineering/internal/fault/injector.go:46",
    slide_num="6 / 24"
)

# Slide 7: Circuit Breaker State Machine
add_slide(
    categories=["RESILIENCE", "PATTERN"],
    title="State Machine Circuit Breaker: Tiga Fasa Proteksi",
    subtitle="Mendeteksi kegagalan beruntun, memutus arus panggilan, dan menguji pemulihan otomatis",
    content_items=[
        {"title": "CLOSED STATE (NORMAL)", "items": [
            "Permintaan dialirkan langsung ke downstream service tanpa hambatan.",
            "Setiap kegagalan menaikkan counter failures secara kumulatif.",
            "Jika counter failures mencapai threshold, state seketika melompat ke OPEN.",
        ], "type": "bullet"},
        {"title": "OPEN STATE (ISOLASI)", "items": [
            "Permintaan langsung dialihkan ke fallback tanpa pernah memanggil downstream (fail-fast).",
            "Melindungi downstream yang sedang sekarat dari gempuran beban trafik baru.",
            "Setelah masa cooldown terlewati, pemanggilan berikutnya menggeser state ke HALF-OPEN.",
        ], "type": "bullet"},
        {"title": "HALF-OPEN (PROBING)", "items": [
            "Mengizinkan satu panggilan uji coba menerobos ke downstream untuk mengetes kesehatan.",
            "Jika panggilan uji coba berhasil: state kembali ke CLOSED dan failure counter di-reset nol.",
            "Jika panggilan uji coba gagal: state langsung terhempas kembali ke OPEN dengan cooldown baru.",
        ], "type": "bullet"},
    ],
    mental_model="A circuit breaker isolates downstream failures and probes recovery without manual restarts.",
    source_path="labs/34-chaos-engineering/internal/circuitbreaker/circuitbreaker.go:12",
    slide_num="7 / 24"
)

# Slide 8: Circuit Breaker Implementation Details
add_slide(
    categories=["BUKTI KODE", "IMPLEMENTASI"],
    title="Eksekusi Terproteksi dan Pengalihan Fallback",
    subtitle="Struktur eksekusi fungsi membungkus logika utama dengan pemulihan otomatis dan isolasi error",
    content_items=[
        {"title": "KODE CIRCUIT BREAKER", "items": [
            "func (cb *CircuitBreaker) Execute(fn, fallback func() error) error {",
            "    cb.mu.Lock()",
            "    cb.checkStateLocked()",
            "    if cb.state == StateOpen {",
            "        cb.mu.Unlock()",
            "        if fallback != nil { return fallback() }",
            "        return ErrCircuitOpen",
            "    }",
            "    cb.mu.Unlock()",
            "    err := fn()",
            "    // Update state & trigger fallback on error",
            "    return cb.handleResult(err, fallback)",
            "}",
        ], "type": "code"},
        {"title": "MANAJEMEN KUNCI MUTEX", "items": [
            "Mutex dilepas sebelum eksekusi fn() untuk mencegah contention saat downstream lambat.",
            "Pengecekan cooldown dilakukan di checkStateLocked() sebelum status OPEN dievaluasi.",
            "State transitions terlindungi penuh tanpa mengunci thread eksekusi I/O utama.",
        ], "type": "bullet"},
        {"title": "GRACEFUL ERROR HANDLING", "items": [
            "Jika fallback didefinisikan, error downstream ditelan dan fallback dieksekusi transparan.",
            "Jika fallback bernilai nil, error asli atau ErrCircuitOpen diteruskan ke pemanggil.",
            "Memungkinkan fleksibilitas tinggi antara fail-silent atau fail-explicit.",
        ], "type": "bullet"},
    ],
    mental_model="Release locks before downstream calls to avoid thread lockups during degraded latency.",
    source_path="labs/34-chaos-engineering/internal/circuitbreaker/circuitbreaker.go:64",
    slide_num="8 / 24"
)

# Slide 9: Graceful Degradation & Fallback Strategy
add_slide(
    categories=["STRATEGI", "USER EXPERIENCE"],
    title="Graceful Degradation: Menyelamatkan Metrik Bisnis",
    subtitle="Ketika primary gateway gagal, antrean fallback memastikan pesanan tidak hilang sia-sia",
    content_items=[
        {"title": "JALUR PRIMARY VS FALLBACK", "items": [
            "Jalur Utama: Pemanggilan sinkron langsung ke Payment Gateway eksternal.",
            "Jalur Fallback: Memasukkan payload transaksi ke antrean asinkron (Payment Queued).",
            "Pengguna menerima konfirmasi antrean alih-alih layar error HTTP 500 yang merusak konversi.",
        ], "type": "bullet"},
        {"title": "DAMPAK PADA METRIK STEADY-STATE", "items": [
            "Monitor mencatat request fallback sebagai SUCCESS di sisi pengalaman konsumen.",
            "Tingkat kegagalan sistem tetap terjaga pada 0.00% meskipun downstream mengalami pemadaman 100%.",
            "Circuit Breaker mencegah lonjakan latensi dengan mengalihkan request secara instan (<1ms).",
        ], "type": "bullet"},
        {"title": "PRINSIP FALLBACK ARSITEKTUR", "items": [
            "Fallback harus beroperasi lokal atau menggunakan infrastruktur penyimpanan yang terisolasi.",
            "Jangan menggunakan fallback yang bergantung pada dependensi jaringan yang sama dengan primary.",
            "Pastikan ada worker rekonsiliasi asinkron untuk memproses ulang data antrean saat downstream pulih.",
        ], "type": "bullet"},
    ],
    mental_model="Graceful fallbacks convert hard outages into delayed asynchronous completions.",
    source_path="labs/34-chaos-engineering/cmd/demo/main.go:38",
    slide_num="9 / 24"
)

# Slide 10: Steady-State Monitor Implementation
add_slide(
    categories=["BUKTI KODE", "OBSERVABILITAS"],
    title="Pelacakan Rasio Kegagalan Menggunakan Operasi Atomik",
    subtitle="Mengeliminasi overhead penguncian pada jalur transaksi dengan memanfaatkan paket sync/atomic",
    content_items=[
        {"title": "KODE MONITOR ATOMIK", "items": [
            "func (m *Monitor) RecordSuccess() {",
            "    atomic.AddUint64(&m.totalRequests, 1)",
            "}",
            "func (m *Monitor) RecordFailure() {",
            "    atomic.AddUint64(&m.totalRequests, 1)",
            "    atomic.AddUint64(&m.failedRequests, 1)",
            "}",
            "func (m *Monitor) IsHealthy() bool {",
            "    total := atomic.LoadUint64(&m.totalRequests)",
            "    if total < 5 { return true } // Warmup minimum",
            "    return m.ErrorRate() <= m.maxErrorRate",
            "}",
        ], "type": "code"},
        {"title": "EFISIENSI ATOMIK", "items": [
            "atomic.AddUint64 berjalan pada level instruksi CPU (lock-free) tanpa blokade mutex antrean.",
            "Cocok untuk throughput jutaan request per detik di lingkungan microservice berskala besar.",
            "Menjamin akurasi penghitungan total dan kegagalan pada pengujian multi-goroutine.",
        ], "type": "bullet"},
        {"title": "EVALUASI THRESHOLD", "items": [
            "ErrorRate() menghitung rasio gagal terhadap total secara real-time.",
            "Threshold kesehatan maxErrorRate (misal 0.20 atau 20%) menjadi dasar evaluasi IsHealthy().",
            "Kombinasi ambang batas dan minimum sample mencegah false positive selama fase rintisan.",
        ], "type": "bullet"},
    ],
    mental_model="Lock-free atomic counters track health metrics without adding contention to hot paths.",
    source_path="labs/34-chaos-engineering/internal/monitor/monitor.go:27",
    slide_num="10 / 24"
)

# Slide 11: Experiment Runner & Auto-Abort
add_slide(
    categories=["BLAST RADIUS", "ORCHESTRATION"],
    title="Orkestrator Eksperimen dan Mekanisme Auto-Abort",
    subtitle="Loop monitoring periodik menjamin penghentian darurat seketika saat metrik terdegradasi",
    content_items=[
        {"title": "KODE SIKLUS RUNNER", "items": [
            "func (e *Experiment) Run(ctx context.Context) error {",
            "    e.injector.SetFault(e.config.Latency, e.config.ForceError)",
            "    ticker := time.NewTicker(e.config.MonitorInterval)",
            "    defer ticker.Stop()",
            "    for {",
            "        select {",
            "        case <-ctx.Done(): e.terminate(StateAborted, ...)",
            "        case <-timeout:   e.terminate(StateCompleted, ...)",
            "        case <-ticker.C:",
            "            if !e.monitor.IsHealthy() {",
            "                e.terminate(StateAborted, reason)",
            "                return fmt.Errorf(\"aborted: %s\", reason)",
            "            }",
            "        }",
            "    }",
            "}",
        ], "type": "code"},
        {"title": "TERMINASI SINKRON", "items": [
            "Fungsi terminate() memanggil e.injector.Clear() secara langsung di dalam blok mutasi state.",
            "Menjamin injeksi latensi dan error mati dalam hitungan milidetik saat ambang batas terlanggar.",
            "Mencegah efek samping tertinggal yang membahayakan siklus pengujian berikutnya.",
        ], "type": "bullet"},
        {"title": "STATUS REASON AUDIT", "items": [
            "Runner menyimpan string alasan pembatalan (abortReason) secara transparan untuk audit log.",
            "Menyajikan persentase pasti pelanggaran error rate pada saat eksperimen dihentikan paksa.",
            "Mempermudah identifikasi kelemahan arsitektur sebelum melanjutkan ke perbaikan berikutnya.",
        ], "type": "bullet"},
    ],
    mental_model="An automated kill switch is mandatory before injecting turbulence into any system.",
    source_path="labs/34-chaos-engineering/internal/experiment/runner.go:60",
    slide_num="11 / 24"
)

# Slide 12: Blast Radius Containment
add_slide(
    categories=["SAFETY", "TATA KELOLA"],
    title="Pengendalian Blast Radius: Mencegah Malapetaka",
    subtitle="Tiga pilar proteksi untuk memastikan pengujian chaos tidak menimbulkan kerugian nyata",
    content_items=[
        {"title": "BATASAN TEMPORAL (WAKTU)", "items": [
            "Setiap eksperimen wajib memiliki durasi waktu maksimum yang kaku (misal 300ms - 500ms).",
            "Jika waktu habis sebelum terjadi anomali, eksperimen selesai dengan status COMPLETED.",
            "Tidak boleh ada eksperimen chaos yang dibiarkan aktif tanpa batas waktu akhir.",
        ], "type": "bullet"},
        {"title": "BATASAN METRIK (HEALTH)", "items": [
            "Pemeriksaan kesehatan dilakukan berulang pada interval pendek (misal 20ms - 40ms).",
            "Pelanggaran ambang batas langsung memicu auto-abort tanpa menunggu durasi eksperimen selesai.",
            "Menahan dampak kegagalan agar tidak pernah menyebar ke lebih dari persentase toleransi.",
        ], "type": "bullet"},
        {"title": "PEMBERSIHAN TOTAL (CLEANUP)", "items": [
            "Metode Clear() mengembalikan flag enabled ke false dan mereset nilai latensi menjadi nol.",
            "Pembersihan dieksekusi sinkron di terminate(), bukan mengandalkan garbage collector atau cron.",
            "Sistem kembali ke status baseline sejati dalam durasi kurang dari 1 milidetik.",
        ], "type": "bullet"},
    ],
    mental_model="Timeouts, metric thresholds, and immediate cleanups define the blast radius sandbox.",
    source_path="labs/34-chaos-engineering/internal/experiment/runner.go:91",
    slide_num="12 / 24"
)

# Slide 13: Demo Walkthrough - Phase 1 & 2
add_slide(
    categories=["DEMO KASUS", "VERIFIKASI"],
    title="Demo Kasus Pembayaran: Baseline dan Mitigasi Sukses",
    subtitle="Simulasi interaktif membuktikan kestabilan baseline dan keampuhan graceful fallback",
    content_items=[
        {"title": "FASE 1: BASELINE (REQ #1 - #5)", "items": [
            "5 request dikirim dalam kondisi downstream normal tanpa adanya gangguan buatan.",
            "Seluruh request berhasil dilayani lewat jalur utama: 'Payment Success (Primary)'.",
            "Circuit Breaker bertahan di state CLOSED; error rate monitor tercatat 0.00%.",
        ], "type": "bullet"},
        {"title": "FASE 2: RESILIENT CLIENT (REQ #6 - #15)", "items": [
            "Eksperimen menyuntikkan kegagalan error 100% pada downstream service selama 300ms.",
            "Req #6 & #7: Gagal di primary, dialihkan mulus ke 'Payment Queued (Fallback)' (CB: CLOSED).",
            "Req #8: Threshold 3 kegagalan terlampaui; Circuit Breaker seketika meloncat ke OPEN.",
            "Req #9 - #15: Fast-fail tanpa memanggil downstream, langsung mengeksekusi fallback (CB: OPEN).",
        ], "type": "bullet"},
        {"title": "HASIL STEADY-STATE", "items": [
            "Total Requests: 15 | Success: 15 | Failed: 0 | Error Rate: 0.00%.",
            "Eksperimen selesai dengan status COMPLETED karena steady-state kesehatan tidak pernah terlanggar.",
            "Pelanggan tetap terlayani tanpa menyadari adanya pemadaman total di downstream gateway.",
        ], "type": "bullet"},
    ],
    mental_model="Graceful degradation keeps error rate at zero even while the circuit breaker trips open.",
    source_path="labs/34-chaos-engineering/cmd/demo/main.go:52",
    slide_num="13 / 24"
)

# Slide 14: Demo Walkthrough - Phase 3 & 4
add_slide(
    categories=["DEMO KASUS", "PEMULIHAN"],
    title="Demo Fase 3 & 4: Auto-Abort dan Pemulihan Mandiri",
    subtitle="Eksperimen tanpa mitigasi memicu pembatalan otomatis, diikuti pemulihan penuh trafik normal",
    content_items=[
        {"title": "FASE 3: AUTO-ABORT UNMITIGATED", "items": [
            "Menjalankan eksperimen kedua pada raw service call tanpa perlindungan fallback.",
            "Raw Call #1 & #2 gagal berturut-turut, mendongkrak error rate melompat ke 33.33%.",
            "Ambang batas 20% dilanggar: Runner langsung membatalkan pengujian dengan status ABORTED.",
            "Alasan pembatalan: 'Steady state breached: error rate 33.33% exceeded threshold'.",
            "Fault Injector langsung dinetralkan otomatis (Active: false).",
        ], "type": "bullet"},
        {"title": "FASE 4: TRAFIK RECOVERY (REQ #16 - #20)", "items": [
            "Setelah jeda cooldown 200ms, trafik produksi normal dialirkan kembali ke sistem.",
            "Req #16 melintas pada mode HALF-OPEN, downstream sukses merespons, state reset ke CLOSED.",
            "Req #17 - #20 meluncur mulus via jalur utama 'Payment Success (Primary)' (CB: CLOSED).",
        ], "type": "bullet"},
        {"title": "BUKTI PENUTUP", "items": [
            "Sistem berhasil melindungi diri sendiri dari gangguan yang melampaui batas toleransi.",
            "Pemulihan berlangsung sepenuhnya otomatis tanpa restart service atau intervensi operator.",
        ], "type": "bullet"},
    ],
    mental_model="Unmitigated turbulence triggers instant auto-abort, enabling self-healing recovery.",
    source_path="labs/34-chaos-engineering/cmd/demo/main.go:90",
    slide_num="14 / 24"
)

# Slide 15: Unit & Concurrency Test Proof
add_slide(
    categories=["TEST HARNESS", "BUKTI EMPIRIS"],
    title="Verifikasi Pengujian: Lima Test Case Lolos 100%",
    subtitle="Rangkaian unit test dan race detector membuktikan akurasi state dan thread-safety",
    content_items=[
        {"title": "TABEL UJI LABORATORIUM", "items": [
            "Test Suite | Target Verifikasi | Hasil",
            "TestFaultInjector | Latensi, Force Error, & Pembersihan Clear() | PASS (0.01s)",
            "TestCircuitBreakerStateTransitions | Closed -> Open -> Half-Open -> Closed | PASS (0.06s)",
            "TestCircuitBreakerGracefulDegradation | Eksekusi fallback menelan error downstream | PASS (0.00s)",
            "TestExperimentAutoAbortOnSteadyStateViolation | Auto-abort saat steady-state terlanggar | PASS (0.01s)",
            "TestConcurrencyAndRace | 50 Goroutine konkuren membaca & memutasi state | PASS (0.00s)",
        ], "type": "table"},
        {"title": "HASIL RACE DETECTOR", "items": [
            "Perintah eksekusi: go test -race ./...",
            "Output: ok labs/34-chaos-engineering/tests 1.481s.",
            "Nol data race terdeteksi pada seluruh interaksi mutex dan atomic pointer.",
        ], "type": "bullet"},
        {"title": "DETERMINISME PENGUJIAN", "items": [
            "Seluruh test dirancang deterministik tanpa ketergantungan pada port jaringan nyata.",
            "Menjamin build pipeline CI/CD dapat mengeksekusi suite chaos dalam waktu di bawah 2 detik.",
        ], "type": "bullet"},
    ],
    mental_model="Automated test assertions with race detection turn theoretical resilience into verified fact.",
    source_path="labs/34-chaos-engineering/tests/chaos_test.go:1",
    slide_num="15 / 24"
)

# Slide 16: Concurrency & Thread-Safety Design
add_slide(
    categories=["ENGINEERING", "THREAD SAFETY"],
    title="Desain Konkurensi: Harmonisasi RWMutex dan Mutex",
    subtitle="Pemisahan strategi sinkronisasi berdasarkan pola akses read-heavy dan write-heavy",
    content_items=[
        {"title": "RWMUTEX PADA FAULT INJECTOR", "items": [
            "Pola pembacaan (Execute, IsEnabled) sangat dominan selama ribuan request per detik.",
            "sync.RWMutex memungkinkan banyak goroutine membaca konfigurasi fault secara simultan.",
            "Write lock hanya diambil sesaat saat injeksi diaktifkan (SetFault) atau dinetralkan (Clear).",
        ], "type": "bullet"},
        {"title": "MUTEX PADA CIRCUIT BREAKER", "items": [
            "Setiap transisi state melibatkan pembaruan variabel kegagalan dan timestamp lastStateChg.",
            "sync.Mutex standar menjamin serialisasi transisi state yang presisi tanpa race condition.",
            "Lock dilepas sebelum pemanggilan delegat fn() untuk mencegah thread bottleneck.",
        ], "type": "bullet"},
        {"title": "ATOMIC PADA MONITOR", "items": [
            "Penghitungan metrik tidak menggunakan lock sama sekali pada level agregasi.",
            "atomic.AddUint64 dan LoadUint64 memastikan pembacaan throughput zero-latency.",
            "Menghilangkan risiko deadlock antara monitor dan komponen pemanggil transaksi.",
        ], "type": "bullet"},
    ],
    mental_model="Match the concurrency primitive to the access pattern: atomic for counters, RWMutex for flags.",
    source_path="labs/34-chaos-engineering/engineering/02-implementation-notes.md",
    slide_num="16 / 24"
)

# Slide 17: Comparison: Naive vs Resilient System
add_slide(
    categories=["KOMPARASI", "BEFORE AFTER"],
    title="Perbandingan: Sistem Naif vs Sistem Berketahanan",
    subtitle="Perbedaan perilaku sistem saat menghadapi kegagalan downstream pada skala produksi",
    content_items=[
        {"title": "TABEL PERBANDINGAN", "items": [
            "Dimensi Evaluasi | Sistem Naif (Tanpa Chaos/CB) | Resilient Client (Lab 34)",
            "Downstream Latency | Thread terkunci, connection pool habis | Fast-fail lewat Circuit Breaker (<1ms)",
            "Downstream Outage | HTTP 500 error langsung ke konsumen | Graceful fallback (antrean transaksi)",
            "Error Rate Pengguna | Melonjak drastis hingga 100% | Tetap 0.00% berkat fallback terisolasi",
            "Verifikasi Ketahanan | Tebakan pasif / manual staging | Eksperimen terprogram & auto-abort",
            "Pemulihan Layanan | Membutuhkan restart service manual | Otomatis pulih via probe Half-Open",
        ], "type": "table"},
        {"title": "ANALISIS PERBEDAAN", "items": [
            "Sistem naif membiarkan kegagalan downstream merembet ke seluruh ekosistem layanan.",
            "Sistem berketahanan mengisolasi anomali pada batas perimeter dan meredam dampaknya.",
        ], "type": "bullet"},
        {"title": "KEUNTUNGAN BISNIS", "items": [
            "Mencegah kerugian finansial akibat transaksi belanja yang gagal terproses.",
            "Menjaga service level objective (SLO) ketersediaan tetap prima di hadapan pelanggan.",
        ], "type": "bullet"},
    ],
    mental_model="Resilient systems bend under turbulence through isolation; naive systems shatter entirely.",
    source_path="labs/34-chaos-engineering/content/02-master-draft.md",
    slide_num="17 / 24"
)

# Slide 18: Trade-offs & Implementation Caveats
add_slide(
    categories=["TRADE OFF", "BATASAN LAB"],
    title="Trade-off dan Batasan Implementasi Laboratorium",
    subtitle="Penyederhanaan model lab untuk kejelasan konsep versus kebutuhan sistem production nyata",
    content_items=[
        {"title": "PENYEDERHANAAN LAB", "items": [
            "Metrik Kumulatif: Lab menggunakan counter kumulatif sederhana; production membutuhkan sliding window.",
            "Injeksi In-Memory: Simulasi dilakukan dalam runtime Go lokal, bukan manipulasi paket jaringan (iptables/eBPF).",
            "Error Rate Sederhana: Belum memperhitungkan latensi persentil p95/p99 atau anomali jitter.",
        ], "type": "bullet"},
        {"title": "FIELD UNUSED DI INJECTOR", "items": [
            "Field errorRate float64 pada internal/fault/injector.go tidak diaktifkan pada logika eksekusi.",
            "Lab sengaja menggunakan determinisme penuh via forceError bool untuk menjamin hasil test yang stabil.",
            "Hindari klaim pengujian probabilistik parsial pada implementasi lab saat ini.",
        ], "type": "bullet"},
        {"title": "OVERHEAD DAN BIAYA", "items": [
            "Penerapan circuit breaker menambahkan beberapa nanodetik latensi CPU per transaksi.",
            "Penyimpanan fallback antrean membutuhkan alokasi storage dan worker pemrosesan asinkron.",
            "Investasi arsitektur sebanding dengan proteksi downtime yang diperoleh di production.",
        ], "type": "bullet"},
    ],
    mental_model="Acknowledge lab simplifications honestly: cumulative counters teach the logic, sliding windows run production.",
    source_path="labs/34-chaos-engineering/content/01-content-brief.md",
    slide_num="18 / 24"
)

# Slide 19: Production Evolution: Sliding Windows & eBPF
add_slide(
    categories=["PRODUKSI", "ROADMAP"],
    title="Jalur Evolusi Menuju Skala Produksi Enterprise",
    subtitle="Mentransformasikan implementasi lab menjadi platform chaos engineering tingkat lanjut",
    content_items=[
        {"title": "SLIDING TIME WINDOWS", "items": [
            "Gantikan counter seumur hidup dengan ring-buffer 60 detik (misal Hystrix / Resilience4j style).",
            "Mencegah riwayat kegagalan masa lalu membebani keputusan isolasi pada trafik masa kini.",
            "Menghasilkan evaluasi kesehatan real-time yang jauh lebih responsif terhadap perubahan kondisi.",
        ], "type": "bullet"},
        {"title": "EBPF & SERVICE MESH INJECTION", "items": [
            "Pindahkan injeksi dari kode aplikasi internal ke layer infrastruktur (Envoy proxy / Chaos Mesh).",
            "Menyuntikkan packet loss, TCP drop, dan DNS failure tanpa mengubah satu baris pun kode bisnis.",
            "Menguji ketahanan sistem secara transparan pada level kernel sistem operasi.",
        ], "type": "bullet"},
        {"title": "CANARY EXPERIMENT ROUTING", "items": [
            "Batasi injeksi hanya pada 1% trafik pengguna nyata melalui header routing khusus (canary).",
            "Meminimalkan dampak bisnis jika terjadi unexpected cascading failure saat pengujian berlangsung.",
            "Integrasi dengan OpenTelemetry distributed tracing untuk memvisualisasikan jalur fault.",
        ], "type": "bullet"},
    ],
    mental_model="Progress from in-memory mocks to eBPF mesh injection as chaos maturity grows.",
    source_path="labs/34-chaos-engineering/content/02-master-draft.md",
    slide_num="19 / 24"
)

# Slide 20: Anti-Patterns in Chaos Engineering
add_slide(
    categories=["ANTI PATTERN", "JEBAKAN"],
    title="Jebakan Umum dalam Mengadopsi Chaos Engineering",
    subtitle="Kesalahan fatal yang dapat mengubah latihan ketahanan menjadi bencana downtime nyata",
    content_items=[
        {"title": "1. MERUSAK TANPA AUTO-ABORT", "items": [
            "Menyuntikkan kegagalan tanpa mekanisme monitoring otomatis dan tombol pemutus darurat.",
            "Ketika sistem benar-benar ambruk, engineer panik mencari cara mematikan skrip pengujian.",
            "Jangan pernah mengeksekusi eksperimen tanpa loop kill-switch yang terverifikasi.",
        ], "type": "bullet"},
        {"title": "2. MENGUKUR METRIK YANG SALAH", "items": [
            "Mengukur utilitas CPU, memori, atau IOPS server alih-alih keberhasilan transaksi bisnis.",
            "CPU rendah saat server down sering disalahartikan sebagai kondisi sistem yang 'santai'.",
            "Selalu jadikan kepuasan pengguna akhir dan error rate transaksi sebagai kompas utama.",
        ], "type": "bullet"},
        {"title": "3. CIRCUIT BREAKER TANPA FALLBACK", "items": [
            "Menerapkan Circuit Breaker tetapi membiarkan error ErrCircuitOpen meledak ke pengguna akhir.",
            "Circuit Breaker hanya memotong beban downstream; tanpa fallback, pengguna tetap menerima error.",
            "Kombinasikan selalu circuit breaker dengan strategi graceful degradation yang matang.",
        ], "type": "bullet"},
    ],
    mental_model="Never inject chaos without an automated kill switch and a customer-preserving fallback.",
    source_path="labs/34-chaos-engineering/content/05-key-takeaways.md",
    slide_num="20 / 24"
)

# Slide 21: When to Use & When to Avoid
add_slide(
    categories=["PANDUAN", "KAPAN PAKAI"],
    title="Panduan Keputusan: Kapan Menggunakan Chaos Engineering",
    subtitle="Kriteria kematangan arsitektur sebelum memulai eksperimen gangguan terprogram",
    content_items=[
        {"title": "SAAT TEPAT MENGGUNAKAN", "items": [
            "Arsitektur Microservices dengan ketergantungan downstream kritis (Payment, Logistik, Auth).",
            "Sistem yang memiliki Service Level Agreement (SLA) 99.9% ke atas dengan beban transaksi masif.",
            "Ketika tim telah mengimplementasikan timeout, retry, dan circuit breaker dan ingin bukti ketahanan.",
            "Sebelum menghadapi lonjakan trafik musiman (Flash Sale, Black Friday, promo puncak).",
        ], "type": "bullet"},
        {"title": "SAAT HARUS MENGHINDARI", "items": [
            "Aplikasi monolitik sederhana tanpa ketergantungan jaringan eksternal.",
            "Sistem yang belum memiliki observabilitas dasar (metrik error rate dan logging belum tersedia).",
            "Ketika tim belum menerapkan penanganan error dan timeout standar pada client HTTP/gRPC.",
            "Menginjeksi chaos pada sistem yang sudah diketahui rapuh hanya membuang waktu tanpa nilai tambah.",
        ], "type": "bullet"},
        {"title": "PRASYARAT WAJIB", "items": [
            "Observabilitas matang, arsitektur fallback terpasang, dan safety guardrails teruji 100%.",
        ], "type": "bullet"},
    ],
    mental_model="Fix fundamental error handling before adopting chaos; test resilience only after building it.",
    source_path="labs/34-chaos-engineering/content/01-content-brief.md",
    slide_num="21 / 24"
)

# Slide 22: Practical Checklist for Production Chaos
add_slide(
    categories=["CHECKLIST", "OPERASIONAL"],
    title="Checklist Kesiapan Pengujian Chaos di Lingkungan Nyata",
    subtitle="Daftar periksa esensial sebelum menyuntikkan gangguan terkontrol ke sistem produksi",
    content_items=[
        {"title": "FASE PERSIAPAN (PRE-FLIGHT)", "items": [
            "[x] Definisikan steady-state berbasis metrik bisnis pengguna (maksimal error rate, latensi p95).",
            "[x] Tentukan hipotesis eksperimen yang spesifik dan terukur.",
            "[x] Konfigurasikan durasi eksperimen dan ambang batas toleransi blast radius.",
        ], "type": "bullet"},
        {"title": "FASE EKSEKUSI (IN-FLIGHT)", "items": [
            "[x] Pastikan background ticker monitor aktif membaca metrik kesehatan secara kontinu.",
            "[x] Verifikasi bahwa fallback menangani kegagalan downstream tanpa mencemari pengalaman pengguna.",
            "[x] Uji fungsi auto-abort darurat untuk memastikan pembatalan otomatis berfungsi saat threshold breached.",
        ], "type": "bullet"},
        {"title": "FASE PASCA-UJI (POST-FLIGHT)", "items": [
            "[x] Verifikasi metode Clear() menetralkan seluruh gangguan seketika.",
            "[x] Pantau pemulihan otomatis Circuit Breaker dari Half-Open kembali ke Closed.",
            "[x] Dokumentasikan temuan eksperimen ke dalam engineering audit record untuk perbaikan sistem.",
        ], "type": "bullet"},
    ],
    mental_model="Run through the pre-flight checklist: measure steady state, automate abort, verify clean reset.",
    source_path="labs/34-chaos-engineering/content/02-master-draft.md",
    slide_num="22 / 24"
)

# Slide 23: Key Architectural Takeaways
add_slide(
    categories=["KESIMPULAN", "INTI TEKNIS"],
    title="Rangkuman Esensial: Ketahanan Terverifikasi Secara Empiris",
    subtitle="Tiga pilar fundamental dalam membangun sistem yang tangguh menghadapi turbulensi",
    content_items=[
        {"title": "1. EKSPERIMEN BUKAN KERUSAKAN", "items": [
            "Chaos engineering adalah sains verifikasi terencana dengan hipotesis dan batas toleransi.",
            "Tujuannya adalah menemukan kelemahan arsitektur sebelum insiden nyata merugikan pengguna.",
            "Steady state yang sehat berpusat pada output akhir, bukan kondisi internal hardware.",
        ], "type": "bullet"},
        {"title": "2. ISOLASI & DEGRADASI", "items": [
            "Circuit Breaker bertindak sebagai pemutus sekring untuk mencegah kegagalan downstream meluas.",
            "Graceful fallback mengubah error fatal menjadi proses tunda yang menyelamatkan transaksi bisnis.",
            "Kombinasi keduanya menjaga metrik kepuasan pengguna tetap sempurna di tengah bencana sistem.",
        ], "type": "bullet"},
        {"title": "3. SAFETY FIRST VIA AUTO-ABORT", "items": [
            "Jangan pernah menjalankan pengujian chaos tanpa loop monitor otomatis dan safety kill-switch.",
            "Netralisasi fault secara sinkron menjamin kembalinya stabilitas sistem tanpa efek samping.",
            "Verifikasi konkurensi dengan race detector untuk menjamin keandalan pada trafik tinggi.",
        ], "type": "bullet"},
    ],
    mental_model="Resilience is proven only when graceful degradation absorbs faults and safety switches halt chaos.",
    source_path="labs/34-chaos-engineering/content/05-key-takeaways.md",
    slide_num="23 / 24"
)

# Slide 24: Closing & References (Dark Navy)
add_slide(
    categories=["PENUTUP", "REFERENSI"],
    title="Eksperimen terencana, interpretasi yang teliti.",
    subtitle="Membuktikan ketahanan sistem terdistribusi melalui kode Go yang deterministik dan terverifikasi",
    content_items=[
        {"title": "INFORMASI LABORATORIUM", "items": [
            f"Software Engineering Lab: {LAB_NO} - {TITLE_ID}",
            "Implementasi: internal/fault, internal/circuitbreaker, internal/monitor, internal/experiment",
            f"Repository: {REPO_URL}",
            f"Commit Snapshot: {COMMIT_SHA[:8]}",
        ], "type": "bullet"},
        {"title": "VERIFIKASI & AUDIT", "items": [
            f"Penulis: {AUTHOR}",
            "Status Pengujian: PASS (5/5 unit & race tests lolos)",
            "Race Detection: PASS (go test -race ./... clean)",
            "Demo Output: Baseline OK, Resilient Fallback OK, Auto-Abort OK, Recovery OK",
        ], "type": "bullet"},
        {"title": "CATATAN RUNTIME", "items": [
            "Lab menggunakan counter kumulatif in-memory untuk determinisme pengujian.",
            "Di lingkungan produksi, terapkan sliding time window dan mesh network injection.",
        ], "type": "bullet"},
    ],
    mental_model="Continuous resilience verification separates robust engineering from naive assumption.",
    source_path="labs/34-chaos-engineering/cmd/demo/main.go",
    slide_num="24 / 24"
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides definition for Lab {LAB_NO}: {TITLE_ID}")
