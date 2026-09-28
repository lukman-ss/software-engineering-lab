# SLI, SLO, Error Budget & Burn Rate Alerting

## Problem

Tim reliability sering mengukur keandalan dengan metrik infrastructure seperti CPU, memory, atau queue depth. Metrik ini tidak mencerminkan pengalaman pengguna. Konklusi: server "sehat" tetapi pengguna mengalami error. Solusi: gunakan SLI (Service Level Indicator) yang mengukur langsung pengalaman pengguna, bukan internal state sistem.

## Why This Matters

Mengukur dengan metrik yang salah menyebabkan dua masalah utama:

1. **False Sense of Safety**: Sistem melaparkan "healthy" infrastructure, padahal error rate nyata tinggi.
2. **Deployment Chaos**: Tim tidak tahu kapan harus berhenti merilis fitur baru dan fokus pada perbaikan keandalan.

SLI/SLO/Error Budget mengubah "harus stabil" menjadi angka yang dapat diputuskan. Error budget berfungsi sebagai "uang" untuk mengalokasikan risiko: semakin besar budget yang tersisa, semakin banyak ruang untuk inovasi.

## Mental Model

Keandalan bukanlah 100%. Lebar dari target keandalan adalah **error budget**.

```
SLO:       99.9% uptime
Error Budget: 0.1% = 1000 request × 0.1% = 1 request boleh error
Real SLI:  99.8%
Budget Consumed: 2 errors > budget
```

Setiap request yang gagal atau melanggar latency threshold mengonsumsi budget. Latency threshold ditentukan oleh caller via closure `isGood`, bukan dibaca `Evaluator`. Ketika budget habis, kebijakan mengatakan: **STOP deployment berisiko**, fokus pada perbaikan.

## Core Concept

### 1. SLI (Service Level Indicator)

SLI adalah ukuran kuantitatif aspek kualitas layanan yang diberikan. Contohnya:

- Request latency (P50, P95, P99)
- Error rate (4xx/5xx per total request)
- Availability (request yang sukses / total request)

Lab menggunakan **ratio good events / total events** sebagai SLI utama. Kriteria "good" ditentukan oleh caller melalui closure `isGood`: `StatusCode < 500 && Duration <= latencyThreshold`. Evaluator tidak membaca `LatencyThreshold` secara langsung.

```go
isGood := func(e Event) bool {
    return e.StatusCode < 500 && e.Duration <= latencyThreshold
}
```

### 2. SLO (Service Level Objective)

SLO adalah target nilai untuk SLI. Tidak ada konsekuensi eksplisit (bukan SLA). Contoh:

```go
Config{
    Name:         "Payment Service Availability",
    TargetUptime: 0.999, // 99.9% SLO
}
```

### 3. Error Budget

Error budget = 1 - SLO. Ini adalah banyaknya "kesalahan yang masih diizinkan".

```
Total Error Budget = (1 - 0.999) × 1000 = 1 error
```

Jika terjadi 2 error, budget habis (–1). Kebijakan `CanDeploy = false`.

## Failure Scenario

1. Traffic normal: 1000 request, 0 error → SLI 100%, budget remaining 1.00
2. Regresi baru: latency naik 500ms, error rate 5%
3. Dalam 1 jam: 100 request, 5 error → SLI 95%, budget habis
4. Sistem mendorong alert: "Budget habis, SEMUA perubahan harus dihentikan"

Tanpa error budget, tim tidak sadar sampai **dwell time incident** sudah parah — tidak ada metrik yang memaksa perhatian pada budget consumption rate.

## How It Works

### Architecture

```
Event → WindowTracker → (total, good, bad) → Evaluator → SLI, Budget
                                   ↓
                        AlertEngine ← dual trackers (short, long)
```

Komponen utama:

- **WindowTracker** (internal/metrics/tracker.go): Sliding window dengan bucket time-based. Menyimpan totalCount, goodCount, badCount per bucket (tidak ada histogram bucket). Evicts data yang sudah terlewati jendela.
- **Evaluator** (internal/slo/evaluator.go): Menghitung SLI, error budget, dan kebijakan deploy. Catatan: field `LatencyThreshold` di `Config` tidak dibaca oleh evaluator; penilaian latency dilakukan di closure `isGood` milik caller.
- **AlertEngine** (internal/alerting/engine.go): Multi-window burn rate check pada short (5m) dan long (60m) windows.

### Error Budget Calculation

```go
totalErrorBudget := (1.0 - e.config.TargetUptime) * float64(total)
budgetRemaining := totalErrorBudget - float64(bad)
```

Jika `budgetRemaining <= 0` dan ada traffic → `CanDeploy = false`.

Formula alternatif (Finding 13, MEDIUM confidence, vendor-specific): Datadog menghitung sisa budget sebagai `100 × (current_status − target) / (100 − target)` dalam persen. Prinsipnya sama dengan rumus lab; lab menggunakan bentuk absolut karena SLI dihitung dalam range 0–1.

Zero traffic edge case: `SLI = 1.0` dan `CanDeploy = true` ketika `total = 0` — menghindari false freeze saat service baru start atau tidak ada trafik (evaluator.go:44-47).

### Burn Rate

Burn rate mengukur laju konsumsi budget relatif terhadap target.

```go
// Contoh: 20 error dari 1000 total request → error rate 2%
actualErrorRate := float64(bad) / float64(total)
allowedErrorRate := 1.0 - targetSLO                 // 0.1% untuk SLO 99.9%
burnRate := actualErrorRate / allowedErrorRate      // 2% / 0.1% = 20x
```

Burn rate 20x artinya Anda mengonsumsi 20× lebih cepat dari yang diizinkan. Pada demo Fase 2, akumulasi 10 error / 1100 total = 0.91% → 9.09x (melewati threshold 6.0x).

### Multi-Window Alerting

Agar tidak memicu alert pada spike singkat:

- **Short window** (5 menit) = sensitif
- **Long window** (60 menit) = normalisasi noise
- Alert **hanya** terpicu jika **kedua window** melewati threshold

Threshold yang direkomendasikan:

| Rule | Burn Rate | Window | Budget |
|------|-----------|--------|--------|
| Page Fast | 14.4× | 1h + 5m | 2% |
| Ticket Slow | 6.0× | 6h + 30m | 5% |

Hanya dua rule di atas yang diimplementasikan di demo (`cmd/demo/main.go:38-49`). Rule Ticket 1×/3d+6h/10% dari rekomendasi Google SRE bersifat research-only, tidak diimplementasikan.

Catatan divergensi: rekomendasi Google SRE (Finding 9) mengklasifikasikan 6.0×/6h+30m sebagai PAGE, sedangkan demo menetapkan 6.0× sebagai TICKET dan 14.4× sebagai PAGE. Implementasi sengaja menurunkan severity 6.0× agar sesuai kebijakan demo dua-tier.

Catatan: field `LongWindow`, `ShortWindow`, `BudgetConsumedPct` di `BurnRateRule` diabaikan oleh `engine.Check()` — semua rule berbagi `shortTracker` dan `longTracker` konstruksi-time.

## Implementation

### WindowTracker

Penyimpanan in-memory dengan slice bucket. Setiap bucket berisi:

```go
type Bucket struct {
    StartTime  time.Time
    TotalCount int64
    GoodCount  int64
    BadCount   int64
}
```

Method `Record()`:

1. Buka mutex lock
2. Evict bucket yang sudah kadaluarsa (lebih lama dari windowSize)
3. Cari atau buat bucket untuk timestamp truncated ke bucketSize
4. Increment totalCount, goodCount, atau badCount

Method `Summary()`:

1. Lock seluruh window
2. Jumlahkan semua bucket
3. Kembalikan total, good, bad

### Evaluator

```go
func (e *Evaluator) Evaluate(now time.Time) Status {
    total, good, bad := e.tracker.Summary(now)
    sli := float64(good) / float64(total) if total > 0 else 1.0
    totalErrorBudget := (1.0 - target) * float64(total)
    budgetRemaining := totalErrorBudget - float64(bad)
    canDeploy := budgetRemaining > 0 || total == 0
    return Status{...}
}
```

Catatan: SLIs selalu dalam range 0–1 (0%–100%).

## Code Walkthrough

### Event Recording

```go
ev := metrics.Event{
    Timestamp:  now.Add(time.Duration(i) * 50 * time.Millisecond),
    Duration:   50 * time.Millisecond,
    StatusCode: 200,
    Endpoint:   "/api/v1/pay",
}
sloTracker.Record(ev)
```

### Evaluasi SLO

```go
status := evaluator.Evaluate(simTime)
fmt.Printf("Target SLO: %.3f%% | SLI: %.4f%% | Budget: %.2f\n",
    status.TargetSLO*100, status.CurrentSLI*100, status.BudgetRemaining)
```

### Burn Rate Alert

```go
alerts := alertEngine.Check(evalTime)
for _, a := range alerts {
    fmt.Printf("ALERT: [%s] %s | Burn: %.2fx\n",
        a.Severity, a.RuleName, a.ShortBurnRate)
}
```

## What the Tests Prove

### Unit Tests

| Test | Verifikasi |
|------|------------|
| TestMetricsWindowTracker | Aggregation benar (12 total, 10 good, 2 bad) |
| TestOutOfOrderTimestamps | Bucket ordering dan eviction benar |
| TestEvaluatorZeroTraffic | SLI = 1.0 pada traffic nol, CanDeploy = true |
| TestSLOEvaluator | CanDeploy = false ketika budget habis |

### Concurrency Test

`TestConcurrencyMetrics` mengeksekusi 20 goroutine × 100 request bersamaan. Semua request tercatat dengan benar:

```go
expectedTotal := int64(numGoroutines * requestsPerGoroutine) // 2000
```

Race detector tidak mengembalikan error.

### Alert Engine Test

`TestAlertEngineBurnRate` memverifikasi:

1. 98 request OK + 2 error (2% error rate) → 20x burn rate → alert Page
2. Short window spike (100x) tetapi long window bersih (0.1x) → **tidak ada alert**

Ini membuktikan bahwa **multi-window check** mengurangi false positive pada spike transien.

## Recovery / Rollback (Hypothetical Procedure — NOT demonstrated in demo)

Ketika `CanDeploy = false` (budget habis):

1. **STOP** semua deployments kecuali P0/security (postmortem requirement)
2. Dianalisis root cause di `/api/v1/pay`
3. Perbaikan dideploy sebagai hotfix (setelah postmortem jika incident >20% budget)
4. Setelah error rate turun, budget kembali positif → `CanDeploy = true`

Catatan: Demo hanya menampilkan fase 1-4 (baseline → incident → alert → comparison). Tidak ada simulasi recovery yang benar-benar dijalankan.

## Production Considerations

- **Window Size**: Lab gunakan 30 menit *kompresi waktu* (mewakili 30 hari dalam simulasi); Google rekomendasi 28 hari (4 minggu) untuk konsistensi weekend.
- **Persistency**: Semua metrics in-memory. Untuk produksi, dump ke TSDB (Prometheus, Datadog) setiap interval singkat.
- **Status Corrections**: Google SRE Workbook membahas "status corrections" untuk mengecualikan jendela pemeliharaan.
- **Multi-endpoint**: Setiap endpoint dapat punya SLO berbeda. Di demo, Payment 99.9% vs Reports 95.0%.

## Common Mistakes

1. **Menggunakan 100% sebagai SLO**: Tidak realistis karena error device/jaringan tidak terkontrol.
2. **Mengukur latency rata-rata**: Mean menyembunyikan tail latency. Gunakan percentile.
3. **Alert pada satu window**: Spike singkat akan memicu false positive. Gunakan multi-window.
4. **SLO untuk infrastructure metrics**: CPU/RAM adalah diagnostic signal, bukan user experience.
5. **Tidak ada kebijakan setelah budget habis**: Error budget tanpa kebijakan = sekadar dashboard.

## Case Study

### Demo Output (Ringkasan)

**PHASE 1 — Baseline Traffic**
```
Total: 1000 | Good: 1000 | Bad: 0
SLI: 100.00% | Budget Remaining: 1.00
CanDeploy: true
```

**PHASE 2 — Incident 10% Error Rate**
```
Total: 1100 | Good: 1090 | Bad: 10
SLI: 99.09% | Budget Remaining: -8.90
CanDeploy: false (Budget exhausted)
```

Perhitungan: SLO 99.9% → allowed error rate 0.1%. 10/1100 = 0.91% error rate → membutuhkan 9.1× burn rate.

**PHASE 3 — Burn Rate Alert**
```
>>> ALERT TRIGGERED: [TICKET] Slow Burn Alert (6.0x) 
    | ShortBurn: 9.09x | LongBurn: 9.09x (Threshold: 6.00x)
```

**PHASE 4 — Criticality Comparison**
```
Payment Target SLO: 99.9% | SLI: 99.09% | Budget: -8.90
Reports Target SLO: 95.0% | SLI: 90.00% | Budget: -5.00
CanDeploy: Both false
```

Meskipun kedua endpoint melewati SLO, **Reports memiliki buffer 5%** yang jauh lebih lebar daripada Payment (0.1%).

## Checklist

- [ ] Tentukan SLI yang mengukur user experience (good/total ratio)
- [ ] Pilih SLO realistis (≤ 99.99%, bukan 100%)
- [ ] Definisikan latency threshold yang relevan
- [ ] Implementasikan error budget policy (halo deployment, postmortem)
- [ ] Setup multi-window burn rate alerting
- [ ] Buat kebijakan per-endpoint based on criticality

## Key Takeaways

1. **SLI = good/total events** (ratio 0–100%). Ukur pengalaman, bukan infrastructure.
2. **SLO < 100%** adalah target yang realistis. Error budget berfungsi sebagai toleransi risiko.
3. **Error budget = 1 - SLO**. Consumtion happened on bad events.
4. **Burn rate = actual_error_rate / allowed_error_rate**. >1 artinya Anda "berlebihan".
5. **Multi-window alert** (short + long) mengurangi false positive spike transit.
6. **CanDeploy policy** memberi keputusan otomatis berdasar budget.
7. **Konfigurasi SLO per endpoint** based on business criticality (Payment 99.9%, Reports 95%).
8. **Thread-safety** terjamin dengan mutex sync.RWMutex di WindowTracker.
9. **Implementasi in-memory** untuk demo; produksi butuh TSDB persistence.
10. **Source**: Google SRE Book, SRE Workbook, Datadog, Prometheus — semua konsisten pada definisi inti.

## Sources

- Google SRE Book Chapter 4: Service Level Objectives
- Google SRE Book Chapter 3: Embracing Risk  
- Google SRE Workbook Chapter 2: Implementing SLOs
- Google SRE Workbook Chapter 5: Alerting on SLOs
- Google SRE Workbook Appendix A: Availability Table
- Google SRE Workbook Appendix B: Error Budget Policy
- Datadog SLO Documentation
- Prometheus Alerting Best Practices