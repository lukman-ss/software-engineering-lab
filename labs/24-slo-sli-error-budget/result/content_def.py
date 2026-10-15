#!/usr/bin/env python3
"""
Publication Builder for Software Engineering Lab 24: SLO, SLI & Error Budget.
"""
import os
import sys

COVER_IMAGE_PATH = "labs/24-slo-sli-error-budget/result/image-cover/cover_art.jpg"
REPO_URL = "https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/24-slo-sli-error-budget"
COMMIT_SHA = "13224b48a35ba08f4a350b7d160ae4826f9e1f55"
LAB_NO = "24"
SLUG = "slo-sli-error-budget"
TITLE_ID = "SLO, SLI & Error Budget"
SUBTITLE_ID = "Mengubah 'Harus Stabil' Menjadi Angka yang Bisa Dideploy"
TITLE_EN = "SLO, SLI & Error Budget"
SUBTITLE_EN = "Turning 'Must Be Stable' Into Deployable Numbers"
AUTHOR = "Lukman / lukman-ss"

DARK_NAVY = "0b1b2d"
TEAL = "087f76"
CREAM = "f4f1e8"
WHITE = "ffffff"
SLATE = "52606b"
ORANGE = "d97706"
RED = "c0392b"
GREEN = "27ae60"

SLIDES = []

def add_slide(categories, title, subtitle, content_items, mental_model, source_path="", slide_num="", narrative="", narrative_en="", title_en=""):
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
        "narrative_en": narrative_en,
    })

# Slide 1: Cover
add_slide(
    categories=["COVER"],
    title=TITLE_ID,
    subtitle=SUBTITLE_ID,
    content_items=[
        {"type": "cover", "title": "", "image": COVER_IMAGE_PATH, "items": [
            "Penulis: " + AUTHOR,
            "Go • Internal Metrics • Alerting Engine",
            "LAB " + LAB_NO + " | SOFTWARE ENGINEERING LAB",
        ]}
    ],
    mental_model="",
    source_path="",
    slide_num="01 / 24",
    narrative="Lab ini mengimplementasikan kerangka kerja SLO/SLI/Error Budget berbasis prinsip Google SRE menggunakan Go. Sistem melacak event request secara real-time melalui sliding window, menghitung rasio keandalan pengguna, dan memicu alert burn rate multi-jendela untuk mengunci deployment ketika error budget habis.",
    narrative_en="This lab implements a Google SRE–based SLO/SLI/Error Budget framework in Go. The system tracks request events in real time via sliding windows, calculates user-facing reliability ratios, and triggers multi-window burn rate alerts to halt deployments when the error budget is exhausted.",
)

# Slide 2: Problem Cards
add_slide(
    categories=["MOTIVASI", "MASALAH"],
    title="Infrastruktur Sehat, Pengguna Error",
    title_en="Healthy Infrastructure, Broken User Experience",
    subtitle="Metrik CPU dan memory tidak membaca pengalaman pengguna — false sense of safety.",
    content_items=[
        {"title": "FALSE SAFETY", "items": [
            "CPU idle 30%, RAM usage normal → dashboard hijau semua",
            "Tapi error rate 5% (5xx) pada endpoint utama",
            "Tim infrastruktur tidak melihat masalah karena metrik internal stabil",
        ], "type": "bullet"},
        {"title": "DEPLOYMENT CHAOS", "items": [
            "Tanpa angka keandalan, setiap release = tebak-tebakan risiko",
            "Tidak ada kriteria objektif untuk stop/deploy",
            "Reliability jadi opini, bukan keputusan berbasis data",
        ], "type": "bullet"},
        {"title": "DWELL TIME INCIDENT", "items": [
            "Tanpa error budget, tim baru sadar saat pengguna melaporkan bug massal",
            "Dwell time insiden berhari-hari sebelum terdeteksi",
            "Satu incident bisa menghabiskan uptime setahun sekaligus",
        ], "type": "bullet"},
    ],
    narrative="Masalah fundamental di industri: tim keandalan mengukur server health via metrik infrastructure (CPU, RAM, disk) yang tidak berkorelasi dengan pengalaman pengguna akhir. Hasilnya, dashboard selalu hijau sementara error rate di lapisan aplikasi meningkat tanpa peringatan. Tanpa framework SLO, tidak ada tolok ukur objektif untuk menentukan kapan deployment harus dihentikan demi fokus pada stabilitas.",
    narrative_en="The core industry problem: reliability teams measure server health via infrastructure metrics (CPU, RAM, disk) that do not correlate with end-user experience. The result is always-green dashboards while application-layer error rates climb silently. Without an SLO framework, there is no objective benchmark for when deployments must pause to focus on stability.",
    mental_model="Infrastructure metrics track machines; SLOs track users. Always measure what the customer experiences.",
    source_path="labs/24-slo-sli-error-budget/content/02-master-draft.md",
    slide_num="02 / 24",
)

# Slide 3: Why Manual Review / Naive Approach Fails
add_slide(
    categories=["BATAS", "PENDekatan"],
    title="Mengapa Pendekatan Naif Gagal di Produksi",
    title_en="Why a Naive Approach Fails in Production",
    subtitle="Threshold statis dan alert tunggal tidak bisa membedakan spike transien dari insiden nyata.",
    content_items=[
        {"title": "THRESHOLD STATIS", "items": [
            "Alert jika error rate > 1% → terlalu sensitif terhadap lalu lintas rendah",
            "Spike 10 error dari 100 request = 10%, tapi dari 1.000.000 = 0.001%",
            "Satu threshold tidak bisa menangkap kedua skenario ini",
        ], "type": "bullet"},
        {"title": "METRIK INFRASTRUKTUR", "items": [
            "CPU spike ≠ error rate; memory leak bisa terjadi tanpa CPU pressure",
            "Error rate bisa stabil meski latency P99 memburuk drastis",
            "Monitor infrastructure blind terhadap user-perceived failures",
        ], "type": "bullet"},
        {"title": "NO DECISION FRAMEWORK", "items": [
            "Tanpa error budget, pertanyaan 'apakah boleh deploy?' tidak terjawab",
            "Decision berbasis feeling vs berdasarkan data consumption rate",
            "Tim engineering kehilangan signal untuk auto-halt policy",
        ], "type": "bullet"},
    ],
    narrative="Pendekatan konvensional mengandalkan threshold statis pada metrik infrastruktur. Masalah utamanya adalah ambang batas tunggal tidak bisa membedakan antara spike transient dengan insiden nyata. Contoh: 10 error dari 100 request terlihat kritis (10%), tetapi 10 error dari 1.000.000 request hampir tidak bermakna (0.001%). Alert engine dengan single threshold akan memberikan false positive tinggi atau false negative tergantung konfigurasi.",
    narrative_en="Conventional approaches rely on static thresholds on infrastructure metrics. The flaw is a single threshold cannot distinguish transient spikes from genuine incidents. For example, 10 errors from 100 requests looks critical (10%), but 10 errors from 1,000,000 requests is negligible (0.001%). Single-threshold alerts produce either high false positives or high false negatives.",
    mental_model="Static thresholds confuse noise with signal; multi-window burn rate separates them.",
    source_path="labs/24-slo-sli-error-budget/content/02-master-draft.md",
    slide_num="03 / 24",
)

# Slide 4: Mental Model — Schema Card
add_slide(
    categories=["MENTAL MODEL", "SKEMA"],
    title="Mental Model: Reliability ≠ 100%",
    title_en="Mental Model: Reliability Is Not 100%",
    subtitle="Error budget adalah 'uang' yang bisa dihabiskan untuk inovasi — kalau habis, berhenti.",
    content_items=[
        {"title": "FORMULA INTI", "items": [
            "SLO: Target keandalan realistis (99.9% = tiga sembilan)",
            "SLI: Ukuran aktual yang dihitung dari good/total events",
            "Error Budget: 1 − SLO = toleransi kesalahan yang diizinkan",
        ], "type": "code"},
        {"title": "ANALOGI DOMPET", "items": [
            "SLO 99.9% → 0.1% budget = 1 dollar per 1000 request",
            "Setiap bad event = pengeluaran 1 dollar dari dompet",
            "Budget habis → stop deploy, fokus perbaikan",
        ], "type": "code"},
        {"title": "BURN RATE", "items": [
            "Burn Rate = actual_error_rate / allowed_error_rate",
            "14.4x = budget habis dalam 2% waktu (fast burn)",
            "6.0x = budget habis dalam 5% waktu (slow burn)",
        ], "type": "code"},
    ],
    narrative="Model mental SLO/Error Budget mengubah kepercayaan subjektif 'layanan harus stabil' menjadi kerangka numerik yang dapat ditindaklanjuti. SLO menetapkan target yang wajar — bukan 100% yang mustahil — dan selisihnya menjadi error budget. Setiap event buruk mengonsumi budget. Burn rate mempercepat perhitungan: 14.4x artinya budget akan habis 14.4 kali lebih cepat dari yang diizinkan, memberi peringatan dini sebelum krisis.",
    narrative_en="The SLO/Error Budget mental model transforms subjective 'service must be stable' into a numerically actionable framework. SLO sets a realistic target — not 100% which is impossible — and the gap becomes the error budget. Each bad event spends from that budget. Burn rate accelerates the timeline: 14.4x means the budget will be exhausted 14.4 times faster, providing early warning before crisis.",
    mental_model="Error budget is spendable risk capital. When it runs out, innovation pauses until reliability recovers.",
    source_path="labs/24-slo-sli-error-budget/internal/slo/evaluator.go",
    slide_num="04 / 24",
)

# Slide 5: Pipeline Diagram
add_slide(
    categories=["ARSITEKTUR", "PIPELINE"],
    title="Pipeline: Event → Tracker → Evaluator → Alert",
    title_en="Pipeline: Event to Alert",
    subtitle="Arsitektur tiga komponen: WindowTracker untuk aggregasi, Evaluator untuk SLO, AlertEngine untuk burn rate.",
    content_items=[
        {"title": "DATA FLOW", "items": [
            "Event{Timestamp, Duration, StatusCode, Endpoint}",
            "↓ WindowTracker.Record() (RWMutex + bucket eviction)",
            "↓ Summary(now) → (total, good, bad) for window",
            "↓ Evaluator.Evaluate() → SLI, Budget Remaining, CanDeploy",
            "↓ AlertEngine.Check() → Multi-window Burn Rate Alerts",
        ], "type": "code"},
        {"title": "KOMPONEN", "items": [
            "metrics.WindowTracker — sliding window, O(bucket_count) query",
            "slo.Evaluator — SLI calculation, budget math, deploy gate",
            "alerting.AlertEngine — dual-window burn rate threshold check",
        ], "type": "code"},
        {"title": "TIME WINDOWS", "items": [
            "SLO tracker: window 30d (compressed to 30min in demo)",
            "Alert short tracker: 5 minutes",
            "Alert long tracker: 60 minutes",
        ], "type": "code"},
    ],
    narrative="Arsitektur pipeline terdiri dari tiga komponen utama yang terkomposisikan secara modular. WindowTracker menerima event HTTP (status code, duration, endpoint) dan mengelompokkannya ke dalam bucket time-series berdasarkan bucketSize. Evaluator memanggil Summary() untuk mendapatkan agregat good/bad/total, lalu menghitung SLI ratio dan sisa error budget. AlertEngine menggunakan dua tracker terpisah (short dan long window) untuk menghitung burn rate dan memutuskan apakah alert harus di-trigger.",
    narrative_en="The pipeline architecture comprises three modularly composed components. WindowTracker accepts HTTP events and buckets them by time. Evaluator calls Summary() to get good/bad/total aggregates, then computes the SLI ratio and remaining error budget. AlertEngine uses two separate trackers to calculate burn rate across short and long windows, deciding whether to trigger an alert.",
    mental_model="Events flow through a read-optimized pipeline: record lock-free, summarize locked, evaluate stateless.",
    source_path="labs/24-slo-sli-error-budget/cmd/demo/main.go",
    slide_num="05 / 24",
)

# Slide 6: Dominant Code — WindowTracker Record
add_slide(
    categories=["KODE", "METRICS"],
    title="WindowTracker: Sliding Window dengan Out-of-Order Support",
    title_en="WindowTracker: Sliding Window with Out-of-Order Support",
    subtitle="RWMutex menjamin thread-safety; bucket insertion menjaga urutan kronologis.",
    content_items=[
        {"title": "RECORD METHOD (tracker.go:46)", "items": [
            "func (w *WindowTracker) Record(e Event) {",
            "    w.mu.Lock()",
            "    defer w.mu.Unlock()",
            "    w.evictStaleLocked(e.Timestamp)",
            "    bucketStart := e.Timestamp.Truncate(w.bucketSize)",
            "    // insert/update bucket in-order ...",
            "}",
        ], "type": "code"},
        {"title": "EVICT LOGIC", "items": [
            "cutoff := now.Add(-w.windowSize)",
            "idx := 0",
            "for idx < len(w.buckets) && buckets[idx].StartTime.Before(cutoff) {",
            "    idx++",
            "}",
            "w.buckets = w.buckets[idx:]",
        ], "type": "code"},
    ],
    narrative="WindowTracker menggunakan sync.RWMutex untuk menjamin thread-safety saat ratusan goroutine merekam event secara konkuren. Metode Record() membuka lock, mengevict bucket yang sudah melewati windowSize, lalu mencari atau membuat bucket baru untuk timestamp yang ditruncate ke bucketSize. Jika timestamp datang out-of-order, bucket disisipkan di posisi yang benar menggunakan linear search — memastikan slice selalu terurut untuk eviction yang efisien.",
    narrative_en="WindowTracker uses sync.RWMutex for thread-safety when hundreds of goroutines record events concurrently. Record() opens a lock, evicts buckets past windowSize, then searches for or creates a new bucket for the truncated timestamp. If timestamps arrive out-of-order, the bucket is inserted at the correct position via linear search, ensuring the slice stays sorted for efficient eviction.",
    mental_model="Sorted bucket insertion + stale eviction = O(1) amortized summary on sliding windows.",
    source_path="labs/24-slo-sli-error-budget/internal/metrics/tracker.go",
    slide_num="06 / 24",
)

# Slide 7: Failure Path — Dangling Budget Exhaustion
add_slide(
    categories=["FAILURE PATH", "INSIDEN"],
    title="Failure Path: Regresi Latency Menghabiskan Budget",
    title_en="Failure Path: Latency Regression Exhausts Budget",
    subtitle="Alur kegagalan dari deploy regresi hingga CanDeploy=false dalam hitungan menit.",
    content_items=[
        {"title": "STEP 1 — DEPLOY REGRESI", "items": [
            "Deployment baru mengubah latency P95 dari 50ms → 500ms",
            "isGood predicate: StatusCode < 500 && Duration <= 200ms",
            "Semua request baru sekarang dihitung 'bad' karena latency threshold dilampaui",
        ], "type": "bullet"},
        {"title": "STEP 2 — BUDGET DIMAKAN", "items": [
            "SLO 99.9% → allowed error rate = 0.1%",
            "Dengan 10% error rate dari latency regression → burn rate = 100x",
            "Total error budget untuk 1000 request = 1 event; 10 bad = -9 remaining",
        ], "type": "bullet"},
        {"title": "STEP 3 — ALERT & HALT", "items": [
            "AlertEngine.Check() membandingkan short + long burn rate",
            "Kedua window >= 6.0x → TICKET alert triggered",
            "Evaluator.CanDeploy = false → deployment pipeline berhenti otomatis",
        ], "type": "bullet"},
    ],
    narrative="Failure path ini menggambarkan skenario nyata: tim deploy fitur baru yang Introduces latency regression. Request yang sebelumnya latency-nya 50ms kini menjadi 500ms, melewati threshold 200ms yang didefinisikan closure isGood. Setiap request buruk mengonsumi error budget 99.9%. Dengan 10% error rate, burn rate mencapai 100x — budget 30 hari akan habis dalam 43 menit. Alert sistem memaksa tim menghentikan semua deployment non-P0.",
    narrative_en="This failure path depicts a real scenario: a team deploys a feature introducing latency regression. Requests previously at 50ms latency now take 500ms, exceeding the 200ms threshold defined by the isGood closure. Each bad event consumes the 99.9% error budget. At 10% error rate, burn rate hits 100x — the 30-day budget exhausts in 43 minutes. The alert system forces all non-P0 deployments to halt.",
    mental_model="Budget exhaustion is a leading indicator: catch it before users file tickets.",
    source_path="labs/24-slo-sli-error-budget/content/02-master-draft.md",
    slide_num="07 / 24",
)

# Slide 8: Comparison Table — Mechanisms vs Symptoms
add_slide(
    categories=["TABEL", "KOMPARASI"],
    title="Mekanisme vs Gejala: Apa yang Terukur",
    title_en="Mechanism vs Symptom: What Gets Measured",
    subtitle="Perbandingan metrik infrastruktur konvensional vs SLO-based observability.",
    content_items=[
        {"title": "TIPE METRIK", "items": [
            "Metrik | Sifat | Hubungan ke Pengguna | Alert Respons",
            "CPU Usage | Internal system | Tidak langsung (korosi lemah) | Sangat lamban",
            "Memory Leak | Internal system | Tidak langsung (korosi lambat) | Lamban ( hours)",
            "Error Rate (SLI) | User-perceived | Langsung dan real-time | Menit",
            "Latency P99 (SLI) | User-perceived | Langsung dan real-time | Menit",
            "Burn Rate |复合 metric | Prediktif (forecasting) | Jam",
        ], "type": "table"},
        {"title": "RESPONS TIM", "items": [
            "Infrastruktur alert: tim ops merespons, tapi tidak tahu 'berapa parah'",
            "SLO alert: tim engineering tahu persis berapa persen budget tersisa",
            "Error budget policy: data-driven decision untuk stop/deploy",
        ], "type": "bullet"},
    ],
    narrative="Tabel ini membandingkan dua paradigm observability: metrik infrastruktur tradisional yang mengukur kesehatan mesin versus SLO-based metrics yang mengukur pengalaman pengguna. CPU dan memory adalah diagnostic signal — berguna untuk debugging setelah incident, tapi tidak prediktif. SLI seperti error rate dan latency P99 memberikan sinyal real-time langsung dari perspektif pengguna, memungkinkan tim merespons sebelum masalah menjadi sistemik.",
    narrative_en="This table contrasts two observability paradigms: traditional infrastructure metrics measuring machine health versus SLO-based metrics measuring user experience. CPU and memory are diagnostic signals — useful for post-incident debugging but not predictive. SLIs like error rate and P99 latency provide real-time signals from the user perspective, enabling response before problems become systemic.",
    mental_model="Infrastructure metrics diagnose machines; SLO metrics protect users.",
    source_path="labs/24-slo-sli-error-budget/content/02-master-draft.md",
    slide_num="08 / 24",
)

# Slide 9: Decision Sequence — Three Focus Cards
add_slide(
    categories=["KEPUTUSAN", "FLOW"],
    title="Tiga Keputusan Kritis dalam Error Budgeting",
    title_en="Three Critical Decisions in Error Budgeting",
    subtitle="Dari definisi SLI hingga respons alert — setiap langkah menentukan arah kebijakan.",
    content_items=[
        {"title": "KEPUTUSAN 1: DEFINISI SLI", "items": [
            "Pilih metrik yang benar-benar diukur pengguna (bukan CPU/memory)",
            "Tentukan closure isGood: status code + latency threshold",
            "SLI harus merefleksikan dampak bisnis, bukan kondisi internal",
        ], "type": "card"},
        {"title": "KEPUTUSAN 2: TARGET SLO", "items": [
            "Pilih target realistis: 99.9% (3九) untuk critical service",
            "Jangan gunakan 100% — network/device errors tidak terkontrol",
            "SLO berbeda per endpoint berdasarkan criticality bucketing",
        ], "type": "card"},
        {"title": "KEPUTUSAN 3: KEBIJAKAN CANDEPLOY", "items": [
            "CanDeploy = false saat budgetRemaining <= 0 (dan total > 0)",
            "Halt semua non-P0 deployment hingga budget recovery",
            "Postmortem wajib jika consumed > 20% budget dalam satu insiden",
        ], "type": "card"},
    ],
    narrative="Tiga keputusan ini membentuk fondasi error budgeting yang efektif. Pertama, memilih SLI yang tepat — metrik yang benar-benar mencerminkan pengalaman pengguna, bukan kondisi internal server. Kedua, menetapkan target SLO realistis yang memperhitungkan failure domains di luar kontrol tim (network, device). Ketiga, defining kebijakan CanDeploy yang mengunci pipeline otomatis ketika budget habis, menghilangkan subjeftivitas dalam pengambilan keputusan release.",
    narrative_en="These three decisions form the foundation of effective error budgeting. First, choosing the right SLI — a metric that truly reflects user experience, not internal server state. Second, setting a realistic SLO target that accounts for out-of-scope failure domains. Third, defining a CanDeploy policy that auto-locks the pipeline when budget is exhausted, removing subjectivity from release decisions.",
    mental_model="SLI choice → SLO target → deploy policy: each decision compounds the next.",
    source_path="labs/24-slo-sli-error-budget/content/02-master-draft.md",
    slide_num="09 / 24",
)

# Slide 10: Execution Evidence — Terminal Output
add_slide(
    categories=["BUKTI", "EKSEKUSI"],
    title="Bukti Eksekusi: Demo Phase 1–4",
    title_en="Execution Evidence: Demo Phases 1–4",
    subtitle="Output terminal asli dari cmd/demo/main.go — empat fase simulasi.",
    content_items=[
        {"title": "PHASE 1: BASELINE", "items": [
            "[PHASE 1] Simulating Baseline Traffic (1,000 requests, 100% success)...",
            "Total: 1000 | Good: 1000 | Bad: 0",
            "Target SLO: 99.900% | Current SLI: 100.0000% | Budget Remaining: 1.00",
            "Deployment Allowed: true",
        ], "type": "code"},
        {"title": "PHASE 2: INCIDENT", "items": [
            "[PHASE 2] Simulating Severe Incident (100 total, 10 errors)...",
            "Total: 1100 | Good: 1090 | Bad: 10",
            "Target SLO: 99.900% | Current SLI: 99.0909% | Budget Remaining: -8.90",
            "Deployment Allowed: false (Budget exhausted)",
        ], "type": "code"},
        {"title": "PHASE 3: ALERT", "items": [
            "[PHASE 3] Checking Multi-Window Burn Rate Alerts...",
            ">>> ALERT TRIGGERED: [TICKET] Slow Burn Alert (6.0x - 5% in 6h)",
            "    | ShortBurn: 9.09x | LongBurn: 9.09x (Threshold: 6.00x)",
        ], "type": "code"},
    ],
    narrative="Output terminal di atas adalah bukti eksekusi langsung dari demo program. Fase 1 menampilkan traffic baseline 1000 request tanpa error — budget remaining 1.00 (satu error masih diizinkan untuk SLO 99.9%). Fase 2 mensimulasikan insiden dengan 10 error dari 100 request tambahan: budget langsung negatif -8.90 dan CanDeploy berubah false. Fase 3 membuktikan alert engine memicu TICKET severity karena kedua window (short dan long) melebihi threshold 6.0x burn rate.",
    narrative_en="The terminal output above is direct execution evidence from the demo program. Phase 1 shows baseline traffic of 1,000 requests with zero errors — budget remaining at 1.00 (one error still allowed for 99.9% SLO). Phase 2 simulates an incident with 10 errors from 100 additional requests: budget immediately goes to -8.90 and CanDeploy flips to false. Phase 3 proves the alert engine triggers TICKET severity because both windows exceed the 6.0x burn rate threshold.",
    mental_model="Raw terminal output is the strongest evidence — numbers don't lie.",
    source_path="labs/24-slo-sli-error-budget/cmd/demo/main.go",
    slide_num="10 / 24",
)

# Slide 11: What Is NOT Proven
add_slide(
    categories=["BATAS", "TIDAK DIBUKTIKAN"],
    title="YANG TIDAK DIBUKTIKAN",
    title_en="WHAT IS NOT PROVEN",
    subtitle="Batasan tegas implementasi lab — apa yang tidak dijamin oleh kode saat ini.",
    content_items=[
        {"title": "YANG TIDAK DIBUKTIKAN", "items": [
            "❌ Persistensi metrics: semua data in-memory, hilang saat restart process",
            "❌ Histogram latency: hanya boolean good/bad, tidak ada P50/P95/P99 distribusi",
            "❌ Window 30-day asli: demo menggunakan kompresi waktu 30 menit (bukan 30 hari)",
            "❌ Threshold burn rate universal: 14.4x/6.0x adalah rekomendasi Google, bukan standar industri mutlak",
            "❌ LatencyThreshold field: ada di Config tapi tidak dibaca oleh Evaluator (dead field)",
            "❌ Recovery simulation: tidak ada simulasi recovery/pengisian ulang budget setelah incident",
            "❌ Multi-endpoint scaling: demo hanya 2 endpoint, tidak menguji ratusan endpoint konkuren",
        ], "type": "bullet"},
    ],
    narrative="Transparansi batasan adalah bagian penting dari engineering integrity. Lab ini secara eksplisit tidak membuktikan persistensi metrics — semua data hilang saat process restart karena menggunakan slice in-memory. Histogram latency juga tidak diimplementasikan; hanya predicate boolean good/bad yang digunakan. Window size 30 hari direpresentasikan sebagai kompresi waktu 30 menit. Field LatencyThreshold di Config Evaluator merupakan dead code yang tidak mempengaruhi perhitungan SLI.",
    narrative_en="Boundary transparency is a critical component of engineering integrity. This lab does not prove metrics persistence — all data is lost on process restart because it uses an in-memory slice. Latency histogram is also not implemented; only a boolean good/bad predicate is used. The 30-day window is represented as a 30-minute time compression. The LatencyThreshold field in the Evaluator Config is dead code that does not affect SLI calculation.",
    mental_model="Know your proof boundaries — unproven claims are liabilities, not features.",
    source_path="labs/24-slo-sli-error-budget/engineering-audit/06-verdict.md",
    slide_num="11 / 24",
)

# Slide 12: Architecture Trade-offs
add_slide(
    categories=["TRADE-OFF", "DESAIN"],
    title="Trade-off Arsitektur: In-Memory vs TSDB",
    title_en="Architecture Trade-off: In-Memory vs TSDB",
    subtitle="Solusi demo optimised untuk correctness dan zero-dependency; produksi membutuhkan persistensi.",
    content_items=[
        {"title": "SIM (DEMO)", "items": [
            "+ Zero-dependency: hanya Go standard library",
            "+ Thread-safe: sync.RWMutex proven dengan race detector",
            "+ Out-of-order safe: bucket insertion maintains chronological order",
            "- No persistence: data loss on process restart",
            "- No horizontal scaling: single-process bottleneck",
        ], "type": "bullet"},
        {"title": "PRODUKSI (REKOMENDASI)", "items": [
            "+ Prometheus/Grafana untuk time-series persistence",
            "+ Horizontal scaling via distributed TSDB cluster",
            "+ Native SLO query language (Prometheus recording rules)",
            "- Additional infrastructure complexity and cost",
            "- Integration testing overhead for data pipeline",
        ], "type": "bullet"},
        {"title": "DESAIN DECISION", "items": [
            "Lab memilih in-memory untuk pedagogi dan zero-dependency",
            "AlertEngine dual-tracker pattern bisa di-replikasi di Prometheus",
            "Burn rate math portable ke any time-series database",
        ], "type": "bullet"},
    ],
    narrative="Trade-off utama antara demo in-memory dan implementasi produksi adalah persistensi versus kesederhanaan. Demo menggunakan slice Go dengan mutex untuk bukti konsep correctness — semua test berlaju green tanpa dependensi eksternal. Untuk produksi, metrics perlu di-dump ke TSDB seperti Prometheus setiap interval singkat agar tetap tersedia setelah restart. Pola dual-window burn rate yang sama dapat direplikasi menggunakan Prometheus recording rules dan alerting expressions.",
    narrative_en="The primary trade-off between the in-memory demo and production implementation is persistence versus simplicity. The demo uses Go slices with mutex for correctness proof — all tests pass green without external dependencies. For production, metrics need to be dumped to a TSDB like Prometheus at short intervals to survive restarts. The same dual-window burn rate pattern can be replicated using Prometheus recording rules and alerting expressions.",
    mental_model="In-memory proves correctness; TSDB provides durability. Choose based on lifecycle stage.",
    source_path="labs/24-slo-sli-error-budget/content/02-master-draft.md",
    slide_num="12 / 24",
)

# Slide 13: Concurrency / Validation Pipeline
add_slide(
    categories=["KONKUREN", "THREAD-SAFE"],
    title="Concurrency: Goroutine-Safe Recording Pipeline",
    title_en="Concurrency: Goroutine-Safe Recording Pipeline",
    subtitle="20 goroutine × 100 request tanpa race condition — terbukti dengan race detector.",
    content_items=[
        {"title": "CONCURRENCY TEST", "items": [
            "func TestConcurrencyMetrics(t *testing.T) {",
            "    var wg sync.WaitGroup",
            "    numGoroutines := 20",
            "    requestsPerGoroutine := 100",
            "    for g := 0; g < numGoroutines; g++ {",
            "        wg.Add(1)",
            "        go func(id int) {",
            "            defer wg.Done()",
            "            for i := 0; i < requestsPerGoroutine; i++ {",
            "                tracker.Record(metrics.Event{...}) // RWMutex protects",
            "            }",
            "        }(g)",
            "    }",
            "    wg.Wait()",
            "}",
        ], "type": "code"},
        {"title": "HASIL", "items": [
            "Expected total: 20 × 100 = 2000",
            "Actual total: 2000 (matched)",
            "good + bad == total: verified",
            "go test -race: clean, zero data races detected",
        ], "type": "bullet"},
    ],
    narrative="Uji konkurensi mengeksekusi 20 goroutine yang masing-masing merekam 100 request secara bersamaan ke WindowTracker yang sama. sync.RWMutex di dalam Record() dan Summary() menjamin bahwa tidak ada interleaving yang merusak counter. Hasilnya: total 2000 event tercatat dengan benar, good + bad == total, dan race detector tidak menemukan data race sama sekali. Ini membuktikan bahwa desain lock granularity cukup untuk workload tinggi.",
    narrative_en="The concurrency test executes 20 goroutines each recording 100 requests simultaneously to the same WindowTracker. The sync.RWMutex inside Record() and Summary() guarantees no interleaving corrupts counters. Result: 2,000 events recorded correctly, good + bad equals total, and the race detector finds zero data races. This proves the lock granularity design is sufficient for high-throughput workloads.",
    mental_model="Mutex-protected shared state with deterministic aggregation is concurrency-safe by construction.",
    source_path="labs/24-slo-sli-error-budget/tests/slo_test.go",
    slide_num="13 / 24",
)

# Slide 14: Before vs After Comparison
add_slide(
    categories=["PERBANDINGAN", "BEFORE/AFTER"],
    title="Before vs After: Observability Transformation",
    title_en="Before vs After: Observability Transformation",
    subtitle="Transformasi dari monitoring infrastructure ke SLO-based reliability engineering.",
    content_items=[
        {"title": "BEFORE (Tanpa SLO)", "items": [
            "Dashboard: CPU 30%, Memory 60%, Disk I/O normal",
            "Alert: 'CPU high' saat load test, tapi user experience tidak terukur",
            "Decision: deploy berdasarkan feeling, tidak ada data budget",
            "Incident response: reactive, dwell time berjam-jam",
        ], "type": "bullet"},
        {"title": "AFTER (Dengan SLO)", "items": [
            "Dashboard: SLI 99.09%, Budget Remaining -8.90, CanDeploy=false",
            "Alert: 'Slow burn 9.09x' — kedua window terdeteksi",
            "Decision: auto-halt deployment, fokus postmortem",
            "Incident response: proactive, budget consumption rate terukur",
        ], "type": "bullet"},
    ],
    narrative="Transformasi observability dari pendekatan infrastructure-centric ke SLO-based mengubah cara tim merespons insiden. Sebelum SLO, dashboard hanya menampilkan metrik mesin — CPU, memory, disk — yang tidak berkorelasi dengan error rate pengguna. Setelah implementasi, dashboard menampilkan SLI aktual, sisa error budget, dan status CanDeploy yang menentukan otomatis apakah deployment baru boleh berjalan. Alert berbasis burn rate memberikan early warning sebelum user mulai melaporkan masalah.",
    narrative_en="The observability transformation from infrastructure-centric to SLO-based changes how teams respond to incidents. Before SLO, dashboards only show machine metrics — CPU, memory, disk — that do not correlate with user error rates. After implementation, dashboards display actual SLI, remaining error budget, and a CanDeploy status that automatically determines whether new deployments can proceed. Burn rate-based alerts provide early warning before users start reporting issues.",
    mental_model="SLO dashboards show user reality; infrastructure dashboards show machine state.",
    source_path="labs/24-slo-sli-error-budget/content/02-master-draft.md",
    slide_num="14 / 24",
)

# Slide 15: SLO Evaluator — Core Math
add_slide(
    categories=["SLO", "MATEMATIKA"],
    title="Evaluator: Matematika SLI dan Error Budget",
    title_en="Evaluator: SLI and Error Budget Mathematics",
    subtitle="Formula matematis yang sederhana namun powerful untuk keputusan deploy.",
    content_items=[
        {"title": "EVALUATE METHOD", "items": [
            "func (e *Evaluator) Evaluate(now time.Time) Status {",
            "    total, good, bad := e.tracker.Summary(now)",
            "",
            "    var sli float64 = 1.0",
            "    if total > 0 {",
            "        sli = float64(good) / float64(total)",
            "    }",
            "    // zero traffic → SLI 1.0, CanDeploy true",
            "",
            "    totalErrorBudget := (1.0 - e.config.TargetUptime) * float64(total)",
            "    budgetRemaining := totalErrorBudget - float64(bad)",
            "    canDeploy := total == 0 || budgetRemaining > 0",
            "}",
        ], "type": "code"},
        {"title": "ZERO TRAFFIC EDGE CASE", "items": [
            "Total = 0: SLI = 1.0 (assumed perfect), CanDeploy = true",
            "Mencegah false freeze saat service baru start atau tidak ada trafik",
            "Validasi: TestEvaluatorZeroTraffic di slo_test.go line 179",
        ], "type": "bullet"},
    ],
    narrative="Evaluator adalah komponen inti yang menghitung tiga nilai: SLI ratio, total error budget, dan remaining budget. Rumus SLI sederhana: good dibagi total, dengan edge case total=0 yang mengembalikan SLI=1.0 untuk mencegah false freeze pada service baru yang belum memiliki trafik. Error budget dihitung sebagai (1 - TargetUptime) × total, yang berarti untuk SLO 99.9% dan 1000 request, budget total adalah 1 error yang diizinkan. CanDeploy=false hanya jika budgetRemaining <= 0 DAN total > 0.",
    narrative_en="The Evaluator is the core component that computes three values: the SLI ratio, total error budget, and remaining budget. The SLI formula is simple: good divided by total, with a total=0 edge case returning SLI=1.0 to prevent false freezes on new services with no traffic. Error budget is (1 - TargetUptime) × total — for 99.9% SLO and 1,000 requests, the total budget is exactly 1 allowed error. CanDeploy is false only when budgetRemaining is ≤ 0 AND total > 0.",
    mental_model="SLI = good/total; budget = (1-SLO)×N; deploy gate = budgetRemaining > 0.",
    source_path="labs/24-slo-sli-error-budget/internal/slo/evaluator.go",
    slide_num="15 / 24",
)

# Slide 16: Burn Rate Alert Engine
add_slide(
    categories=["ALERT", "BURN RATE"],
    title="AlertEngine: Multi-Window Burn Rate Calculation",
    title_en="AlertEngine: Multi-Window Burn Rate Calculation",
    subtitle="Dua window (short + long) menghindari false positive dari spike transient.",
    content_items=[
        {"title": "BURN RATE FORMULA", "items": [
            "func (a *AlertEngine) CalculateBurnRate(total, bad int64) float64 {",
            "    if total == 0 { return 0.0 }",
            "    actualErrorRate := float64(bad) / float64(total)",
            "    allowedErrorRate := 1.0 - a.targetSLO",
            "    return actualErrorRate / allowedErrorRate",
            "}",
        ], "type": "code"},
        {"title": "CHECK METHOD (dual window)", "items": [
            "shortBurn := CalculateBurnRate(shortTotal, shortBad)",
            "longBurn := CalculateBurnRate(longTotal, longBad)",
            "triggered := shortBurn >= rule.BurnRateFactor && longBurn >= rule.BurnRateFactor",
            "# AND condition — both must exceed threshold",
        ], "type": "code"},
        {"title": "TRANSIENT SPIKE PROOF", "items": [
            "Short window: 100x burn rate (10% errors in 5 min)",
            "Long window: 0.1x burn rate (0.01% errors in 60 min)",
            "Result: NO alert triggered — spike is transient, not systemic",
        ], "type": "bullet"},
    ],
    narrative="AlertEngine menghitung burn rate untuk dua window secara paralel: short window (5 menit) yang sensitif dan long window (60 menit) yang menormalkan noise. Alert hanya di-trigger jika KEDUA window melewati threshold yang sama — ini mekanisme anti-false-positive. Test kasus Transient Spike membuktikan: short window menunjukkan 100x burn rate (spike akut), tapi long window hanya 0.1x (latar belakang normal bersih). Hasilnya: tidak ada alert, karena spike bersifat sementara dan tidak mengancam budget jangka panjang.",
    narrative_en="AlertEngine calculates burn rate for two windows in parallel: a sensitive short window (5 minutes) and a noise-normalizing long window (60 minutes). An alert triggers only when BOTH windows exceed the same threshold — this is the anti-false-positive mechanism. The Transient Spike test case proves the point: short window shows 100x burn rate (acute spike), but long window shows only 0.1x (clean background). Result: no alert, because the spike is transient and does not threaten long-term budget.",
    mental_model="Dual-window AND logic filters transient spikes while catching sustained budget burns.",
    source_path="labs/24-slo-sli-error-budget/internal/alerting/engine.go",
    slide_num="16 / 24",
)

# Slide 17: Test Evidence — SLO Evaluator
add_slide(
    categories=["TEST", "BUKTI"],
    title="Test: SLO Evaluator — Budget Exhaustion & Zero Traffic",
    title_en="Test: SLO Evaluator — Budget Exhaustion and Zero Traffic",
    subtitle="Dua test case membuktikan boundary condition CanDeploy dan zero-traffic fallback.",
    content_items=[
        {"title": "TEST: BUDGET EXHAUSTION", "items": [
            "99 good + 1 bad → SLI = 0.99 → budgetRemaining = 0 → CanDeploy = true ✓",
            "Add 1 more bad → 99 good + 2 bad → SLI = 0.9804",
            "budgetRemaining = -1 → CanDeploy = false ✓",
            "Assertion: !status.CanDeploy saat budget < 0",
        ], "type": "code"},
        {"title": "TEST: ZERO TRAFFIC", "items": [
            "No events recorded → total = 0",
            "SLI forced to 1.0 (no division by zero)",
            "CanDeploy = true (avoid false freeze on fresh service)",
            "Assertion: status.CurrentSLI == 1.0 && status.CanDeploy == true",
        ], "type": "code"},
    ],
    narrative="Dua test case ini menutupi boundary condition paling kritis dari Evaluator. Test pertama memverifikasi bahwa CanDeploy berubah false tepat saat budget.remaining turun ke nol atau negatif — 99 good + 2 bad pada SLO 99% memicu CanDeploy=false. Test kedua menangani edge case zero traffic: ketika tidak ada event sama sekali, SLI diset default ke 1.0 dan CanDeploy tetap true untuk menghindari false freeze pada service yang baru dimulai atau belum menerima traffik.",
    narrative_en="These two test cases cover the most critical boundary conditions of the Evaluator. The first verifies that CanDeploy flips to false exactly when budgetRemaining drops to zero or negative — 99 good plus 2 bad on 99% SLO triggers CanDeploy=false. The second handles the zero-traffic edge case: when no events exist at all, SLI defaults to 1.0 and CanDeploy stays true to prevent false freezes on newly started services.",
    mental_model="Boundary conditions are where reliability bugs hide — test them explicitly.",
    source_path="labs/24-slo-sli-error-budget/tests/slo_test.go",
    slide_num="17 / 24",
)

# Slide 18: Test Evidence — Alert Engine
add_slide(
    categories=["TEST", "BUKTI"],
    title="Test: AlertEngine — Sustained vs Transient Burn",
    title_en="Test: AlertEngine — Sustained versus Transient Burn",
    subtitle="Bukti bahwa dual-window check membedakan insiden nyata dari spike sementara.",
    content_items=[
        {"title": "TEST: SUSTAINED BURN → ALERT", "items": [
            "98 OK + 2 error dalam short (5m) dan long (60m) window",
            "Error rate = 2% / 0.1% allowed = 20x burn rate",
            "Both windows >= 14.4x threshold → Page alert triggered ✓",
            "Severity = PAGE, RuleName = 'Critical Fast Burn'",
        ], "type": "code"},
        {"title": "TEST: TRANSIENT SPIKE → NO ALERT", "items": [
            "Short window: 90 OK + 10 error = 10% → 100x burn rate",
            "Long window: 9999 OK + 1 error = 0.01% → 0.1x burn rate",
            "Short >= 14.4x ✓, Long >= 14.4x ✗",
            "Result: 0 alerts — spike is transient, long trend is healthy ✓",
        ], "type": "code"},
    ],
    narrative="AlertEngine diuji dengan dua skenario berlawanan untuk memverifikasi dual-window logic. Skenario pertama: 2% error rate persisten di kedua window menghasilkan burn rate 20x — kedua window melewati threshold 14.4x, memicu PAGE alert. Skenario kedua: spike transien 10% error di short window (100x burn rate) tapi long window tetap bersih di 0.1x. Karena long window tidak melewati threshold, alert tidak di-trigger — membuktikan bahwa mekanisme dual-window berhasil menyaring false positive dari spike singkat.",
    narrative_en="AlertEngine is tested with two opposite scenarios to verify dual-window logic. Scenario one: a persistent 2% error rate across both windows produces a 20x burn rate — both windows exceed the 14.4x threshold, triggering a PAGE alert. Scenario two: a transient 10% error spike in the short window (100x burn rate) while the long window remains clean at 0.1x. Since the long window does not exceed the threshold, no alert fires — proving the dual-window mechanism successfully filters false positives from short spikes.",
    mental_model="If only one window screams, it might be noise. Both must agree for action.",
    source_path="labs/24-slo-sli-error-budget/tests/slo_test.go",
    slide_num="18 / 24",
)

# Slide 19: Out-of-Order & Eviction Test
add_slide(
    categories=["TEST", "EDGE CASE"],
    title="Test: Out-of-Order Timestamps & Bucket Eviction",
    title_en="Test: Out-of-Order Timestamps and Bucket Eviction",
    subtitle="WindowTracker menangani event datang acak dan evict bucket kadaluarsa dengan benar.",
    content_items=[
        {"title": "OUT-OF-ORDER RECORDING", "items": [
            "Record event at t+5s FIRST, then t+2s SECOND",
            "Bucket at t+2s created after t+5s bucket exists",
            "Linear search inserts new bucket at correct chronological position",
            "Summary(t+6s): total=3, good=2, bad=1 ✓",
        ], "type": "code"},
        {"title": "PARTIAL EVICTION", "items": [
            "Advance time to t+13s → cutoff = t+3s",
            "Bucket at t+2s is BEFORE cutoff → evicted",
            "Bucket at t+5s is AFTER cutoff → retained",
            "Summary(t+13s): total=1, good=1, bad=0 ✓ (only +5s bucket remains)",
        ], "type": "code"},
    ],
    narrative="WindowTracker dirancang untuk menangani event yang datang dengan timestamp tidak berurutan — pola umum di distributed systems di mana clock skew atau network delay menyebabkan reorder. Test membuktikkan bahwa event di t+5 detik direkam pertama, kemudian event di t+2 detik, dan bucket insertion secara otomatis Menyisipkan bucket baru di posisi kronologis yang benar. Partial eviction juga terverifikasi: saat waktu maju ke t+13 detik, cutoff di t+3 detik menghapus bucket t+2s tapi mempertahankan bucket t+5s.",
    narrative_en="WindowTracker is designed to handle events arriving out-of-order — a common pattern in distributed systems where clock skew or network delay causes reordering. The test proves that a t+5s event recorded first, followed by a t+2s event, results in automatic chronological bucket insertion. Partial eviction is also verified: when time advances to t+13s, the cutoff at t+3s evicts the t+2s bucket while retaining the t+5s bucket.",
    mental_model="Sorted bucket insertion with time-based eviction handles distributed clock skew gracefully.",
    source_path="labs/24-slo-sli-error-budget/tests/slo_test.go",
    slide_num="19 / 24",
)

# Slide 20: Implementation Checklist
add_slide(
    categories=["CHECKLIST", "IMPLEMENTASI"],
    title="Checklist Implementasi SLO/Error Budget",
    title_en="SLO/Error Budget Implementation Checklist",
    subtitle="Daftar periksa untuk mengimplementasikan kerangka kerja ini di sistem produksi.",
    content_items=[
        {"title": "FUNDAMENTAL", "items": [
            "□ Tentukan SLI yang mengukur user experience (good/total ratio)",
            "□ Pilih SLO realistis ≤ 99.99% — hindari 100% perfectionism",
            "□ Definisi closure isGood: statusCode threshold + latency threshold",
            "□ Implementasikan WindowTracker dengan sliding window + bucket eviction",
        ], "type": "bullet"},
        {"title": "ALERTING", "items": [
            "□ Setup dual tracker: short window (5-15m) dan long window (1-6h)",
            "□ Pilih burn rate threshold: 14.4x (PAGE) dan 6.0x (TICKET)",
            "□ Implementasikan AND logic: alert hanya jika kedua window exceed threshold",
            "□ Integrasi dengan notification system (PagerDuty, Slack, email)",
        ], "type": "bullet"},
        {"title": "POLICY", "items": [
            "□ CanDeploy gate: auto-halt deployment saat budgetRemaining <= 0",
            "□ Policy exception: P0/security hotfix tetap bisa deploy dengan approval",
            "□ Postmortem requirement: wajib jika single incident消耗 > 20% budget",
            "□ Dashboard SLO: real-time SLI, budget gauge, burn rate trend",
        ], "type": "bullet"},
    ],
    narrative="Checklist ini mencakup lima lapisan implementasi: fundamental measurement, alerting infrastructure, dan policy enforcement. Lapisan fundamental memastikan SLI terdefinisi dengan benar dari perspektif pengguna. Lapisan alerting menetapkan dual-window burn rate dengan threshold yang sesuai severity. Lapisan policy adalah yang paling krusial — tanpa kebijakan CanDeploy yang enforcement, error budget hanya menjadi dashboard cantik tanpa tindakan.",
    narrative_en="This checklist covers five implementation layers: fundamental measurement, alerting infrastructure, and policy enforcement. The fundamental layer ensures SLI is correctly defined from the user perspective. The alerting layer establishes dual-window burn rate with severity-appropriate thresholds. The policy layer is the most critical — without an enforced CanDeploy policy, error budget is merely a pretty dashboard with no action.",
    mental_model="Measurement without policy is vanity. Policy without measurement is opinion.",
    source_path="labs/24-slo-sli-error-budget/content/02-master-draft.md",
    slide_num="20 / 24",
)

# Slide 21: Anti-Pattern & Reproduction
add_slide(
    categories=["ANTI-PATTERN", "REPRODUKSI"],
    title="Anti-Pattern & Cara Reproduksinya",
    title_en="Anti-Patterns and How to Reproduce Them",
    subtitle="Pola kesalahan umum yang menghambat implementasi SLO yang efektif.",
    content_items=[
        {"title": "ANTI-PATTERN 1: 100% SLO", "items": [
            "Pola: Menetapkan SLO 99.9999% (six nines)",
            "Dampak: budget hampir nihil, setiap error = freeze instan",
            "Reproduksi: Set TargetUptime=0.999999, record 1 error → CanDeploy=false",
        ], "type": "bullet"},
        {"title": "ANTI-PATTERN 2: Single-Window Alert", "items": [
            "Pola: Alert hanya berdasarkan short window (5 menit)",
            "Dampak: false positive tinggi dari spike transient (deploy, cron job)",
            "Reproduksi: Record 10 error di short window, 0 error di long window → alert false trigger",
        ], "type": "bullet"},
        {"title": "ANTI-PATTERN 3: Infrastructure SLI", "items": [
            "Pola: Menggunakan CPU/memory sebagai SLI primary metric",
            "Dampak: dashboard hijau tapi pengguna mengalami error rate tinggi",
            "Reproduksi: Simulasikan CPU 30% + error rate 10% → CPU alert tidak nyala, user complaint naik",
        ], "type": "bullet"},
        {"title": "ANTI-PATTERN 4: No Deploy Policy", "items": [
            "Pola: Budget habis tapi tidak ada mekanisme auto-halt",
            "Dampak: tim tetap deploy fitur baru sambil budget kosong",
            "Reproduksi: Set budgetRemaining=-5 tapi skip CanDeploy check → deployment jalan terus",
        ], "type": "bullet"},
    ],
    narrative="Empat anti-pattern ini adalah jebakan umum yang dihadapi tim saat mengadopsi SLO. Anti-pattern pertama adalah perfeksionisme 100% — menetapkan SLO enam sembilan membuat budget hampir mustahil dipenuhi, memicu freeze deployment konstan. Anti-pattern kedua, single-window alert, menyebabkan noise alert berlebihan dari spike transient. Anti-pattern ketiga, menggunakan metrik infrastruktur sebagai SLI, menciptakan false confidence. Anti-pattern keempat, tidak ada kebijakan deploy otomatis, membuat error budget hanya menjadi dashboard tanpa actionability.",
    narrative_en="These four anti-patterns are common traps teams face when adopting SLOs. The first is perfectionism — setting six-nines SLO makes the budget nearly impossible to maintain, triggering constant deployment freezes. The second, single-window alerting, causes alert noise from transient spikes. The third, using infrastructure metrics as SLI, creates false confidence. The fourth, having no auto-deploy policy, renders error budget a dashboard without actionability.",
    mental_model="The best SLO is the highest one your users will accept — not the highest number you can write.",
    source_path="labs/24-slo-sli-error-budget/content/02-master-draft.md",
    slide_num="21 / 24",
)

# Slide 22: Practical Takeaway & Command Block
add_slide(
    categories=["TAKEAWAY", "COMMAND"],
    title="Takeaway Praktis: Jalankan Lab Ini",
    title_en="Practical Takeaway: Run This Lab",
    subtitle="Semua kode Go standar — tanpa dependensi eksternal — siap dijalankan di terminal Anda.",
    content_items=[
        {"title": "COMMAND BLOCK", "items": [
            "# 1. Navigate to lab directory",
            "cd labs/24-slo-sli-error-budget",
            "",
            "# 2. Run all tests with race detector",
            "go test -race -v ./...",
            "",
            "# 3. Run the demo to see live SLO evaluation",
            "go run ./cmd/demo",
        ], "type": "code"},
        {"title": "EXPECTED OUTPUT", "items": [
            "Phase 1: BudgetRemaining=1.00, CanDeploy=true (baseline)",
            "Phase 2: BudgetRemaining=-8.90, CanDeploy=false (incident)",
            "Phase 3: TICKET alert triggered at 9.09x burn rate",
            "Phase 4: Payment vs Reports comparison (99.9% vs 95% SLO)",
        ], "type": "bullet"},
    ],
    narrative="Untuk memahami cara kerja SLO/Error Budget secara mendalam, jalankan seluruh test suite dengan race detector aktif. Perintah go test -race -v ./... akan menjalankan lima test case: WindowTracker aggregation, SLO evaluator budget math, AlertEngine dual-window logic, Out-of-order timestamp handling, dan Concurrency stress test. Perintah go run ./cmd/demo menampilkan simulasi empat fase secara real-time — dari baseline traffic hingga incident response dengan alert triggering.",
    narrative_en="To deeply understand how SLO and Error Budget work, run the full test suite with the race detector enabled. The go test -race -v ./... command runs five test cases: WindowTracker aggregation, SLO evaluator budget math, AlertEngine dual-window logic, out-of-order timestamp handling, and concurrency stress testing. The go run ./cmd/demo command displays a four-phase simulation in real time — from baseline traffic through incident response with alert triggering.",
    mental_model="Run the tests. See the numbers. Trust the evidence.",
    source_path="labs/24-slo-sli-error-budget/README.md",
    slide_num="22 / 24",
)

# Slide 23: Long-term Architectural Principles
add_slide(
    categories=["PRINSIP", "JANGKA PANJANG"],
    title="Prinsip Arsitektural Jangka Panjang",
    title_en="Long-term Architectural Principles",
    subtitle="Tiga prinsip yang mengatur hubungan antara reliability engineering dan product velocity.",
    content_items=[
        {"title": "PRINSIP 1: RELIABILITY IS A FEATURE", "items": [
            "SLO bukan meta-metric — ia adalah user-facing quality contract",
            "Setiap percentage point SLO bernilai uang dan trust pengguna",
            "Budget allocation = risk allocation — decide consciously",
        ], "type": "card"},
        {"title": "PRINSIP 2: AUTOMATE THE HALT", "items": [
            "CanDeploy=false harus memblokir pipeline CI/CD secara otomatis",
            "Manual override hanya untuk P0/security dengan audit trail",
            "Policy enforcement > human judgment untuk mencegah drift",
        ], "type": "card"},
        {"title": "PRINSIP 3: MEASURE USER EXPERIENCE", "items": [
            "SLI harus mengukur apa yang dirasakan pengguna, bukan mesin",
            "CPU 30% meaningless jika error rate 5%; CPU 90% fine jika SLI 99.99%",
            "Dashboard harus menjawab: 'Apakah pengguna kami bahagia?'",
        ], "type": "card"},
    ],
    narrative="Tiga prinsip arsitektural ini membentuk filosofi long-term reliability engineering. Prinsip pertama menegaskan bahwa reliability bukan side-feature tapi core product value — setiap percentage point SLO bernilai uang dan trust pengguna. Prinsip kedua menekankan automasi: kebijakan CanDeploy harus enforcement secara teknis di pipeline CI/CD, bukan bergantung pada disiplin manual. Prinsip ketiga mengingatkan bahwa metrik harus user-centric — dashboard yang indah tapi tidak menjawab 'apakah pengguna bahagia' adalah wasted instrumentation.",
    narrative_en="These three architectural principles form the philosophy of long-term reliability engineering. The first principle asserts that reliability is not a side-feature but core product value — every SLO percentage point is real money and user trust. The second principle emphasizes automation: the CanDeploy policy must be technically enforced in the CI/CD pipeline, not reliant on manual discipline. The third principle reminds us that metrics must be user-centric — a beautiful dashboard that does not answer 'are our users happy' is wasted instrumentation.",
    mental_model="Automate the guardrails. Measure the user. Spend the budget wisely.",
    source_path="labs/24-slo-sli-error-budget/content/05-key-takeaways.md",
    slide_num="23 / 24",
)

# Slide 24: Closing
add_slide(
    categories=["CLOSING"],
    title="Error Budget: Uang Risiko yang Harus Dihabiskan",
    title_en="Error Budget: Risk Capital That Must Be Spent",
    subtitle="",
    content_items=[
        {"title": "RINGKASAN LAB", "items": [
            "3 komponen: WindowTracker (metrics) + Evaluator (SLO) + AlertEngine (burn rate)",
            "5 test cases: aggregation, budget math, dual-window alert, out-of-order, concurrency",
            "Demo 4 phase: baseline → incident → alert → cross-SLO comparison",
            "Zero third-party dependencies — pure Go 1.22 standard library",
        ], "type": "bullet"},
        {"title": "BUKTI VERIFIKASI", "items": [
            f"Source: {REPO_URL}",
            f"Commit: {COMMIT_SHA[:8]}",
            "Tests: ALL PASS (go test -race ./... clean)",
            "Demo: 4-phase execution verified",
        ], "type": "bullet"},
    ],
    narrative="Lab 24 membuktikan bahwa kerangka kerja SLO/Error Budget dapat diimplementasikan dalam 128 baris Go code murni tanpa dependensi eksternal. Tiga komponen modular — WindowTracker untuk aggregasi event, Evaluator untuk perhitungan SLI dan budget, AlertEngine untuk multi-window burn rate — bekerja bersama membentuk sistem observability yang mesure用户体验 langsung. Lima test case mencakup semua boundary condition dan konkurensi. Pesan terakhir: error budget harus dibakar, bukan ditabung. Budget yang tidak terpakai artinya Anda terlalu konservatif dan melewatkan kesempatan innovasi.",
    narrative_en="Lab 24 proves that an SLO/Error Budget framework can be implemented in 128 lines of pure Go code with zero external dependencies. Three modular components — WindowTracker for event aggregation, Evaluator for SLI and budget calculation, AlertEngine for multi-window burn rate — work together to form a user-experience-direct observability system. Five test cases cover all boundary conditions and concurrency. The final message: error budget must be spent, not saved. Unused budget means you are being too conservative and missing innovation opportunities.",
    mental_model="An unspent error budget is wasted opportunity. Spend it on controlled risk, not hoarded perfection.",
    source_path=REPO_URL,
    slide_num="24 / 24",
)

if __name__ == "__main__":
    print(f"Loaded {len(SLIDES)} slides for Lab {LAB_NO}: {TITLE_ID}")
