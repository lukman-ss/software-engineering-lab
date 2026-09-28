# Content Brief

Topic:
Chaos Engineering dan Fault Injection — memverifikasi ketahanan sistem terdistribusi melalui eksperimen terkontrol: steady-state metric, injection latency/error, Circuit Breaker, fallback, dan auto-abort blast radius.

Target Reader:
Software engineer dan SRE yang merancang layanan terdistribusi, ingin membuktikan klaim ketahanan (timeout, circuit breaker, graceful degradation) sebelum insiden production, dan perlu model mental yang terikat pada implementasi lab yang sudah diverifikasi.

Problem:
Klaim ketahanan (failover, circuit breaker, fallback) sering tidak diuji sampai kegagalan nyata terjadi. Tanpa eksperimen terstruktur, fault downstream bisa menumpuk menjadi cascading failure. Tanpa batas blast radius dan tombol abort, eksperimen itu sendiri bisa merusak seluruh pengguna.

Core Mental Model:
Chaos Engineering bukan merusak sistem secara acak. Alurnya: definisikan steady state (metrik output, bukan kondisi internal) → hipotesiskan bahwa metrik itu tetap terjaga meski fault X disuntikkan → suntikkan fault terkontrol → ukur → batalkan otomatis jika threshold dilanggar. Ketahanan terbukti hanya jika fallback/circuit breaker menjaga error rate pengguna; tanpa mitigasi, monitor wajib mematikan injector.

Approved Research Status:
APPROVED (research-audit/07-verdict.md, 2026-09-28). Non-blocking: URL Netflix/Google SRE masih tag/TOC; perbedaan pustaka bahasa ditunda ke engineering.

Approved Engineering Status:
APPROVED (engineering-audit/06-verdict.md, 2026-09-28). Non-blocking: field `errorRate` di `internal/fault/injector.go` tidak dipakai; injeksi error memakai `forceError bool`.

Main Concepts:
1. Steady state sebagai metrik output (error rate, bukan health internal komponen).
2. Hipotesis eksperimen: "steady state tetap terjaga meski fault X disuntikkan."
3. Fault injection in-memory: latency + forced error.
4. Circuit Breaker: Closed → Open → Half-Open, dengan fallback.
5. Blast radius: auto-abort + `injector.Clear()` sinkron saat threshold dilanggar.
6. Graceful degradation: fallback mencatat success di sisi pengguna meski primary gagal.

Verified Behaviors:
- Fault injector: disabled → SetFault (latency + error) → Clear; diuji `TestFaultInjector`.
- Circuit breaker: 2 failure trip ke Open, fast-fail `ErrCircuitOpen`, cooldown ke Half-Open, sukses reset ke Closed; diuji `TestCircuitBreakerStateTransitions`.
- Fallback menelan error primary; diuji `TestCircuitBreakerGracefulDegradation`.
- Auto-abort saat error rate > threshold; injector dinetralkan; diuji `TestExperimentAutoAbortOnSteadyStateViolation`.
- Concurrency aman di race detector; diuji `TestConcurrencyAndRace` + `go test -race ./...`.
- Demo: baseline 5 sukses CLOSED; 10 request dengan fallback (CB OPEN, ErrorRate=0.00%); eksperimen tanpa mitigasi ABORTED (error rate 33.33%); recovery 5 sukses CLOSED.

Available Case Studies:
Lab demo `cmd/demo` — simulasi payment gateway. Bukan insiden production. Bukan studi Netflix/AWS empiris di dalam lab ini.

Warnings:
- Lab disederhanakan: metrik kumulatif, bukan sliding window; injeksi in-memory, bukan network/proxy; threshold error rate, bukan p99 latency.
- Field `errorRate` tidak aktif; jangan klaim injeksi probabilistik.
- Nilai threshold (20%, 25%), failure threshold (2, 3), cooldown (50ms, 200ms) adalah contoh lab, bukan rekomendasi production.
- Sumber Netflix/Google SRE adalah indeks/tag, bukan deep-link bab.
- Tidak didemonstrasikan: canary routing, OpenTelemetry, compound multi-service failure, otomatisasi CI/CD ke production.
- `pkg/*` di design awal; kode aktual di `internal/*`.
