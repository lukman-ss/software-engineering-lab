# Key Takeaways

1. **SLI = good events ÷ total events** (ratio 0–100%). Ukur pengalaman pengguna langsung, bukan CPU/memory.

2. **SLO < 100% adalah target nyata**. 100% tidak realistis karena device dan jaringan pengguna tidak terkontrol.

3. **Error Budget = 1 − SLO**. Ini "uang" untuk berinovasi. Budget habis = deployment berhenti, fokus perbaiken.

4. **Budget formula**: `totalBudget = (1 − targetSlo) × totalEvents`; `remaining = totalBudget − badEvents`; `CanDeploy = false` ketika `remaining ≤ 0` dan ada traffic.

5. **Burn Rate = actualErrorRate ÷ allowedErrorRate**. >1 berarti consumption melebihi target. Semakin tinggi, semaakin cepat intervention diperlukan.

6. **Multi-window alerting (short + long) mencegah false positive**. Spike transit di short window saja tidak memicu alert bila long window masih bersih.

7. **Criticality berbeda → SLO berbeda**. Payment 99.9% (0.1% error tolerance) jauh lebih ketat daripada Reports 95% (5% error tolerance) — hal yang sama (10% error) habiskan budget keduanya, tapi Payment jauh lebih rentan.

8. **Thread-safety terjamin**. `sync.RWMutex` di `WindowTracker` memastikan invariant `good + bad == total` pada kondisi 20 goroutine konkuren. Race detector bersih.

9. **Zero traffic = SLI 1.0, CanDeploy = true**. Evaluator menghindari false freeze ketika tidak ada request.

10. **Implementasi in-memory saja untuk demo**. Produksi memerlukan persistence ke TSDB (Prometheus/Datadog) dan status corrections untuk maintenance windows.

11. **Burn rate thresholds (14.4× PAGE, 6.0× TICKET) adalah rekomendasi Google SRE**, bukan standar universal. Datadog menggunakan skema berbeda (6× 2h window) tetapi prinsipnya sama.

12. **Semua sumber asli Google SRE, Datadog, Prometheus, GCP** — konsisten pada definisi inti, namun statistik seperti "70% outages from changes" adalah observasi internal Google tanpa verifikasi independen.
