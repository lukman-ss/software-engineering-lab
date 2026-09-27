# SLI, SLO, Error Budget & Burn Rate Alerting

## Problem

Tim reliability sering mengukur keandalan dengan metrik infrastructure seperti CPU, memory, atau queue depth. Metrik ini tidak mencerminkan pengalaman pengguna. Konklusi: server "sehat" tetapi pengguna mengalami error. Solusi: gunakan SLI (Service Level Indicator) yang mengukur langsung pengalaman pengguna, bukan internal state sistem.

## Why This Matters

Mengukur dengan metrik yang salah menyebabkan dua masalah utama:

1. **False Sense of Safety**: Sistem melaparkan "healthy" infrastructure, padahal error rate nyata tinggi.
2. **Deployment Chaos**: Tim tidak tahu kapan harus berhenti merilis fitur baru dan fokus pada perbaikan keandalan.

SLI/SLO/Error Budget mengubah "harus stabil" menjadi angka yang dapat diputuskan. Error budget berfungsi sebagai "uang" untuk mengalokasikan risiko: semakin besar budget yang tersisa, semakin banyak ruang untuk inovasi.

## Mental Model

Keandalan bukanlah 100%. Lebar dari target keandaran adalah **error budget**.

```
SLO:       99.9% uptime
Error Budget: 0.1% = 1000 request × 0.1% = 1 request boleh error
Real SLI:  99.8%
Budget Consumed: 2 errors > budget
```

Setiap request yang gagal atau melanggar latency threshold mengonsumsi budget. Ketika budget habis, kebijakan mengatakan: **STOP deployment berisiko**, fokus pada perbaikan.

## Core Concept

### 1. SLI (Service Level Indicator)

SLI adalah ukuran kuantitatif aspek kualitas layanan yang diberikan. Contohnya:

- Request latency (P50, P95, P99)
- Error rate (4xx/5xx per total request)
- Availability (request yang sukses / total request)

Lab menggunakan **ratio good events / total events** sebagai SLI utama. Sebuah request dikatakan "good" jika statusnya < 500 **dan** durasinya ≤ latency threshold.

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

1. Traffic normal: 1000 request, 0 error → SLI 100%, budget +0.1
2. Regresi baru: latency naik 500ms, error rate 5%
3. Dalam 1 jam: 100 request, 5 error → SLI 95%, budget habis
4. Sistem mendorong alert: "Budget habis, SEMUA perubahan harus dihentikan"

Tanpa error budget, tim tidak sadar sampai DINP (Dwell Time Incident Performance) sudah parah.

## How It Works

### Architecture

```
Event → WindowTracker → (total, good, bad) → Evaluator → SLI, Budget
                                   ↓
                        AlertEngine ← dual trackers (short, long)
```

Komponen utama:

- **WindowTracker** (internal/metrics/tracker.go): Sliding window dengan bucket time-based. Menyimpan totalCount, goodCount, badCount per bucket. Evicts data yang sudah terlewati jendela.
- **Evaluator** (internal/slo/evaluator.go): Menghitung SLI, error budget, dan kebijakan deploy.
- **AlertEngine** (internal/alerting/engine.go): Multi-window burn rate check pada short (5m) dan long (60m) windows.

### Error Budget Calculation

```go
totalErrorBudget := (1.0 - e.config.TargetUptime) * float64(total)
budgetRemaining := totalErrorBudget - float64(bad)
```

Jika `budgetRemaining <= 0` dan ada traffic → `CanDeploy = false`.

### Burn Rate

Burn rate mengukur laju konsumsi budget relatif terhadap target.

```go
actualErrorRate := float64(bad) / float64(total)   // mis. 10% / 1000
allowedErrorRate := 1.0 - targetSLO                 // 0.1%
burnRate := actualErrorRate / allowedErrorRate      // 10% / 0.1% = 100x
```

Burn rate 100x artinya Anda mengonsumsi 100× lebih cepat dari boleh.

### Multi-Window Alerting

Agar tidak memicu alert pada spike singkat:

- **Short window** (5 menit) = sensitif
- **Long window** (60 menit) = normalisasi noise
- Alert **hanya** terpicu jika **kedua window** melewati threshold

Threshold yang direkomendasikan:

| Rule | Burn Rate | Window | Budget |
|------|-----------|--------|--------|
| Page Fast | 14.4× | 1h + 5m | 2% |
| Page Slow | 6.0× | 6h + 30m | 5% |
| Ticket | 1× | 3d + 6h | 10% |

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
2. Short window sengekan (100x) tetapi long window bersih (0.1x) → **tidak ada alert**

Ini membuktikan bahwa **multi-window check** mengurangi false positive pada spike transit.

## Recovery / Rollback

Ketika `CanDeploy = false`:

1. **STOP** semua non-essential deployments
2. Dianalisis root cause di `/api/v1/pay`
3. Perbaikan dideploy sebagai hotfix **dengan approval P0**
4. Setelah error rate turun, budget kembali positif → `CanDeploy = true`

No magic automation: keputusan tetap manual, hanya informasi yang lebih akurat.

## Production Considerations

- **Window Size**: Lab gunakan 30 hari rolling window; Google rekomendasi 28 hari (4 minggu) untuk konsistensi weekend.
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
7. **Congure SLO per endpoint** based on business criticality (Payment 99.9%, Reports 95%).
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