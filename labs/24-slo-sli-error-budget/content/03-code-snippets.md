# Code Snippets

Semua snippet diambil verbatim dari implementasi yang telah di-audit dan disetujui. Tidak ada kode sintetis.

---

## Snippet 1 — Definisi Event dan Bucket

Source File: `internal/metrics/tracker.go` (line 8-20)
Purpose: Struktur data inti untuk merepresentasikan request dan agregasinya per bucket waktu.

```go
type Event struct {
	Timestamp  time.Time
	Duration   time.Duration
	StatusCode int
	Endpoint   string
}

type Bucket struct {
	StartTime  time.Time
	TotalCount int64
	GoodCount  int64
	BadCount   int64
}
```

Explanation: `Event` merepresentasikan satu HTTP request dengan timestamp, durasi, status code, dan endpoint. `Bucket` mengagregasi event per interval waktu (bucketSize) untuk efisiensi memori pada sliding window.

---

## Snippet 2 — WindowTracker: Inisialisasi dan Struktur

Source File: `internal/metrics/tracker.go` (line 22-44)
Purpose: Sliding-window time-bucketed tracker dengan thread-safety.

```go
type WindowTracker struct {
	mu           sync.RWMutex
	windowSize   time.Duration
	bucketSize   time.Duration
	buckets      []Bucket
	isGoodEvent  func(e Event) bool
}

func NewWindowTracker(windowSize time.Duration, bucketSize time.Duration, isGood func(e Event) bool) *WindowTracker {
	if bucketSize <= 0 {
		bucketSize = time.Second
	}
	numBuckets := int(windowSize / bucketSize)
	if numBuckets < 1 {
		numBuckets = 1
	}
	return &WindowTracker{
		windowSize:  windowSize,
		bucketSize:  bucketSize,
		buckets:     make([]Bucket, 0, numBuckets),
		isGoodEvent: isGood,
	}
}
```

Explanation: `WindowTracker` menyimpan slice bucket dengan `sync.RWMutex` untuk thread-safety. `windowSize` menentukan jendela rolling (mis. 30 menit), `bucketSize` menentukan granularitas agregasi. `isGoodEvent` adalah predicate yang menentukan apakah event dianggap "good" — di lab: `StatusCode < 500 && Duration <= threshold`.

---

## Snippet 3 — WindowTracker: Record dengan Out-of-Order Handling

Source File: `internal/metrics/tracker.go` (line 46-104)
Purpose: Mencatat event ke bucket yang benar, termasuk handling timestamp out-of-order.

```go
func (w *WindowTracker) Record(e Event) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.evictStaleLocked(e.Timestamp)

	bucketStart := e.Timestamp.Truncate(w.bucketSize)
	good := w.isGoodEvent(e)

	n := len(w.buckets)
	if n > 0 && w.buckets[n-1].StartTime.Equal(bucketStart) {
		w.buckets[n-1].TotalCount++
		if good {
			w.buckets[n-1].GoodCount++
		} else {
			w.buckets[n-1].BadCount++
		}
		return
	}

	if n > 0 && bucketStart.Before(w.buckets[n-1].StartTime) {
		for i := 0; i < n; i++ {
			if w.buckets[i].StartTime.Equal(bucketStart) {
				w.buckets[i].TotalCount++
				if good {
					w.buckets[i].GoodCount++
				} else {
					w.buckets[i].BadCount++
				}
				return
			}
			if w.buckets[i].StartTime.After(bucketStart) {
				b := Bucket{
					StartTime:  bucketStart,
					TotalCount: 1,
				}
				if good {
					b.GoodCount = 1
				} else {
					b.BadCount = 1
				}
				w.buckets = append(w.buckets[:i], append([]Bucket{b}, w.buckets[i:]...)...)
				return
			}
		}
	}

	b := Bucket{
		StartTime:  bucketStart,
		TotalCount: 1,
	}
	if good {
		b.GoodCount = 1
	} else {
		b.BadCount = 1
	}
	w.buckets = append(w.buckets, b)
}
```

Explanation: Tiga jalur: (1) bucket terakhir cocok → increment langsung (fast path), (2) timestamp lebih awal → cari posisi terurut dan insert atau increment, (3) timestamp baru → append bucket baru. `evictStaleLocked` dipanggil terlebih dahulu untuk membuang data kadaluarsa.

---

## Snippet 4 — WindowTracker: Eviction dan Summary

Source File: `internal/metrics/tracker.go` (line 106-128)
Purpose: Membuang bucket kadaluarsa dan mengagregasi total/good/bad.

```go
func (w *WindowTracker) evictStaleLocked(now time.Time) {
	cutoff := now.Add(-w.windowSize)
	idx := 0
	for idx < len(w.buckets) && w.buckets[idx].StartTime.Before(cutoff) {
		idx++
	}
	if idx > 0 {
		w.buckets = w.buckets[idx:]
	}
}

func (w *WindowTracker) Summary(now time.Time) (total int64, good int64, bad int64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.evictStaleLocked(now)
	for _, b := range w.buckets {
		total += b.TotalCount
		good += b.GoodCount
		bad += b.BadCount
	}
	return total, good, bad
}
```

Explanation: `evictStaleLocked` menghapus semua bucket dengan `StartTime < now - windowSize`. `Summary` melakukan eviction lalu menjumlahkan semua bucket. Menggunakan `Lock` (bukan `RLock`) karena eviction memodifikasi slice.

---

## Snippet 5 — SLO Evaluator: Config dan Status

Source File: `internal/slo/evaluator.go` (line 10-32)
Purpose: Definisi konfigurasi SLO dan hasil evaluasi.

```go
type Config struct {
	Name             string
	TargetUptime     float64
	LatencyThreshold time.Duration
}

type Status struct {
	SLOName          string
	TargetSLO        float64
	CurrentSLI       float64
	TotalEvents      int64
	GoodEvents       int64
	BadEvents        int64
	TotalErrorBudget float64
	BudgetConsumed   float64
	BudgetRemaining  float64
	CanDeploy        bool
}
```

Explanation: `Config` menyimpan nama layanan, target SLO (mis. 0.999), dan threshold latensi. `Status` adalah hasil evaluasi yang berisi SLI aktual, budget total/terpakai/sisa, dan flag `CanDeploy` untuk kebijakan release.

---

## Snippet 6 — SLO Evaluator: Evaluate

Source File: `internal/slo/evaluator.go` (line 41-71)
Purpose: Menghitung SLI, error budget, dan kebijakan deployment.

```go
func (e *Evaluator) Evaluate(now time.Time) Status {
	total, good, bad := e.tracker.Summary(now)

	var sli float64 = 1.0
	if total > 0 {
		sli = float64(good) / float64(total)
	}

	allowedFailureRate := 1.0 - e.config.TargetUptime
	totalErrorBudget := allowedFailureRate * float64(total)
	budgetConsumed := float64(bad)
	budgetRemaining := totalErrorBudget - budgetConsumed

	canDeploy := true
	if total > 0 && budgetRemaining <= 0 {
		canDeploy = false
	}

	return Status{
		SLOName:          e.config.Name,
		TargetSLO:        e.config.TargetUptime,
		CurrentSLI:       math.Round(sli*10000) / 10000,
		TotalEvents:      total,
		GoodEvents:       good,
		BadEvents:        bad,
		TotalErrorBudget: math.Round(totalErrorBudget*100) / 100,
		BudgetConsumed:   budgetConsumed,
		BudgetRemaining:  math.Round(budgetRemaining*100) / 100,
		CanDeploy:        canDeploy,
	}
}
```

Explanation: SLI = good/total (default 1.0 saat tidak ada trafik). `totalErrorBudget` = allowed failure rate × total events. `budgetRemaining` negatif berarti budget habis. `CanDeploy` menjadi false hanya jika ada trafik dan budget habis — pada zero traffic selalu true untuk menghindari false freeze.

---

## Snippet 7 — Alerting: Burn Rate Rule dan Engine

Source File: `internal/alerting/engine.go` (line 9-49)
Purpose: Definisi rule burn rate dan konstruktor AlertEngine.

```go
type Severity string

const (
	SeverityPage   Severity = "PAGE"
	SeverityTicket Severity = "TICKET"
	SeverityNone   Severity = "NONE"
)

type BurnRateRule struct {
	Name              string
	Severity          Severity
	LongWindow        time.Duration
	ShortWindow       time.Duration
	BurnRateFactor    float64
	BudgetConsumedPct float64
}

type AlertEngine struct {
	targetSLO    float64
	shortTracker *metrics.WindowTracker
	longTracker  *metrics.WindowTracker
	rules        []BurnRateRule
}

type AlertResult struct {
	Triggered     bool
	RuleName      string
	Severity      Severity
	LongBurnRate  float64
	ShortBurnRate float64
	ThresholdRate float64
}
```

Explanation: `BurnRateRule` mendefinisikan threshold burn rate untuk alert. `AlertEngine` memegang dua tracker terpisah (short dan long window) dan daftar rule. `PAGE` untuk insiden kritis, `TICKET` untuk slow burn.

---

## Snippet 8 — Alerting: CalculateBurnRate

Source File: `internal/alerting/engine.go` (line 51-61)
Purpose: Menghitung burn rate dari total dan bad events.

```go
func (a *AlertEngine) CalculateBurnRate(total, bad int64) float64 {
	if total == 0 {
		return 0.0
	}
	actualErrorRate := float64(bad) / float64(total)
	allowedErrorRate := 1.0 - a.targetSLO
	if allowedErrorRate <= 0 {
		return 0.0
	}
	return actualErrorRate / allowedErrorRate
}
```

Explanation: `burn rate = actual_error_rate / allowed_error_rate`. Contoh: SLO 99.9% → allowed 0.1%. Jika actual 2% → burn rate = 20x. Mengembalikan 0.0 untuk zero traffic atau SLO 100% (allowed = 0).

---

## Snippet 9 — Alerting: Check (Multi-Window)

Source File: `internal/alerting/engine.go` (line 63-89)
Purpose: Evaluasi multi-window burn rate — alert hanya jika kedua window melewati threshold.

```go
func (a *AlertEngine) Check(now time.Time) []AlertResult {
	shortTotal, _, shortBad := a.shortTracker.Summary(now)
	longTotal, _, longBad := a.longTracker.Summary(now)

	shortBurn := a.CalculateBurnRate(shortTotal, shortBad)
	longBurn := a.CalculateBurnRate(longTotal, longBad)

	var results []AlertResult
	for _, rule := range a.rules {
		triggered := false
		if shortBurn >= rule.BurnRateFactor && longBurn >= rule.BurnRateFactor {
			triggered = true
		}

		if triggered {
			results = append(results, AlertResult{
				Triggered:     true,
				RuleName:      rule.Name,
				Severity:      rule.Severity,
				LongBurnRate:  longBurn,
				ShortBurnRate: shortBurn,
				ThresholdRate: rule.BurnRateFactor,
			})
		}
	}
	return results
}
```

Explanation: Mengambil summary dari kedua tracker, menghitung burn rate masing-masing, lalu memeriksa setiap rule. Kondisi `shortBurn >= threshold && longBurn >= threshold` mencegah false positive dari spike transit yang hanya muncul di short window.

---

## Snippet 10 — Demo: Konfigurasi dan Simulasi Baseline

Source File: `cmd/demo/main.go` (line 17-52)
Purpose: Setup SLO, tracker, dan alert rules untuk demonstrasi.

```go
targetSLO := 0.999 // 99.9% availability
latencyThreshold := 200 * time.Millisecond

isGood := func(e metrics.Event) bool {
	return e.StatusCode < 500 && e.Duration <= latencyThreshold
}

window30d := 30 * time.Minute
shortWin := 5 * time.Minute
longWin := 60 * time.Minute

sloTracker := metrics.NewWindowTracker(window30d, 10*time.Second, isGood)
shortTracker := metrics.NewWindowTracker(shortWin, 5*time.Second, isGood)
longTracker := metrics.NewWindowTracker(longWin, 10*time.Second, isGood)

evaluator := slo.NewEvaluator(slo.Config{
	Name:             "Payment Service Availability",
	TargetUptime:     targetSLO,
	LatencyThreshold: latencyThreshold,
}, sloTracker)

alertRules := []alerting.BurnRateRule{
	{
		Name:           "Critical Fast Burn (14.4x - 2% in 1h)",
		Severity:       alerting.SeverityPage,
		BurnRateFactor: 14.4,
	},
	{
		Name:           "Slow Burn Alert (6.0x - 5% in 6h)",
		Severity:       alerting.SeverityTicket,
		BurnRateFactor: 6.0,
	},
}

alertEngine := alerting.NewAlertEngine(targetSLO, shortTracker, longTracker, alertRules)
```

Explanation: Demo menggunakan skala waktu terkompresi (30 menit mewakili 30 hari). Tiga tracker terpisah untuk SLO evaluation, short-window alerting, dan long-window alerting. Predicate `isGood` menggabungkan status code dan latency threshold.

---

## Snippet 11 — Demo: Simulasi Insiden dan Evaluasi

Source File: `cmd/demo/main.go` (line 74-104)
Purpose: Simulasi insiden 10% error rate dan pengecekan burn rate.

```go
fmt.Println("\n[PHASE 2] Simulating Severe Incident (100 total requests, 10 errors = 10% error rate)...")
incStart := simTime.Add(time.Second)
for i := 0; i < 90; i++ {
	ev := metrics.Event{
		Timestamp:  incStart.Add(time.Duration(i) * 10 * time.Millisecond),
		Duration:   50 * time.Millisecond,
		StatusCode: 200,
		Endpoint:   "/api/v1/pay",
	}
	sloTracker.Record(ev)
	shortTracker.Record(ev)
	longTracker.Record(ev)
}
for i := 0; i < 10; i++ {
	ev := metrics.Event{
		Timestamp:  incStart.Add(time.Duration(90+i) * 10 * time.Millisecond),
		Duration:   500 * time.Millisecond,
		StatusCode: 500,
		Endpoint:   "/api/v1/pay",
	}
	sloTracker.Record(ev)
	shortTracker.Record(ev)
	longTracker.Record(ev)
}

evalTime := incStart.Add(2 * time.Second)
status = evaluator.Evaluate(evalTime)
fmt.Printf("Total: %d | Good: %d | Bad: %d\n", status.TotalEvents, status.GoodEvents, status.BadEvents)
fmt.Printf("Deployment Allowed: %v (Budget exhausted)\n", status.CanDeploy)

fmt.Println("\n[PHASE 3] Checking Multi-Window Burn Rate Alerts...")
alerts := alertEngine.Check(evalTime)
for _, a := range alerts {
	fmt.Printf(">>> ALERT TRIGGERED: [%s] %s | ShortBurn: %.2fx | LongBurn: %.2fx (Threshold: %.2fx)\n",
		a.Severity, a.RuleName, a.ShortBurnRate, a.LongBurnRate, a.ThresholdRate)
}
```

Explanation: 90 request sukses + 10 request gagal (status 500, latency 500ms) menghasilkan 10% error rate. Pada SLO 99.9% (allowed 0.1%), burn rate = 100x. Evaluator mendeteksi budget habis (CanDeploy = false). AlertEngine memicu slow burn alert (6.0x) karena kedua window melewati threshold.

---

## Snippet 12 — Test: Verifikasi Alert Tidak Terpicu pada Spike Transit

Source File: `tests/slo_test.go` (line 130-151)
Purpose: Membuktikan multi-window mencegah false positive.

```go
shortOnlyTracker := metrics.NewWindowTracker(5*time.Minute, time.Second, isGood)
longCleanTracker := metrics.NewWindowTracker(60*time.Minute, time.Second, isGood)
engineTransient := alerting.NewAlertEngine(0.999, shortOnlyTracker, longCleanTracker, rules)

for i := 0; i < 90; i++ {
	shortOnlyTracker.Record(metrics.Event{Timestamp: now, StatusCode: 200})
}
for i := 0; i < 10; i++ {
	shortOnlyTracker.Record(metrics.Event{Timestamp: now, StatusCode: 500})
}
for i := 0; i < 9999; i++ {
	longCleanTracker.Record(metrics.Event{Timestamp: now, StatusCode: 200})
}
longCleanTracker.Record(metrics.Event{Timestamp: now, StatusCode: 500})

transientAlerts := engineTransient.Check(now)
if len(transientAlerts) != 0 {
	t.Fatalf("expected 0 alerts for transient spike (long window below threshold), got %d", len(transientAlerts))
}
```

Explanation: Short window: 10% error (100x burn rate). Long window: 10.000 request, 1 error = 0.01% error rate (0.1x burn rate < 14.4x). Alert TIDAK terpicu karena long window di bawah threshold — membuktikan nilai multi-window untuk mengurangi alert fatigue.

---

## Snippet 13 — Test: Concurrency dengan 20 Goroutine

Source File: `tests/slo_test.go` (line 200-236)
Purpose: Memverifikasi thread-safety WindowTracker di bawah beban konkuren.

```go
func TestConcurrencyMetrics(t *testing.T) {
	tracker := metrics.NewWindowTracker(10*time.Second, 100*time.Millisecond, func(e metrics.Event) bool {
		return e.StatusCode == 200
	})

	now := time.Now()
	var wg sync.WaitGroup
	numGoroutines := 20
	requestsPerGoroutine := 100

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < requestsPerGoroutine; i++ {
				status := 200
				if i%10 == 0 {
					status = 500
				}
				tracker.Record(metrics.Event{
					Timestamp:  now.Add(time.Duration(i) * time.Millisecond),
					StatusCode: status,
				})
			}
		}(g)
	}
	wg.Wait()

	total, good, bad := tracker.Summary(now.Add(time.Second))
	expectedTotal := int64(numGoroutines * requestsPerGoroutine)
	if total != expectedTotal {
		t.Fatalf("expected total %d, got %d", expectedTotal, total)
	}
	if good+bad != total {
		t.Fatalf("expected good (%d) + bad (%d) == total (%d)", good, bad, total)
	}
}
```

Explanation: 20 goroutine masing-masing mencatat 100 request (10% error rate) secara konkuren ke tracker yang sama. Race detector (`go test -race`) mengonfirmasi tidak ada data race. Invariant `good + bad == total` selalu terpenuhi.
