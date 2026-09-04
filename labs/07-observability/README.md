# Senior Software Engineer Daily #7

## Lab 07 — Observability: Kenapa Logging Saja Tidak Cukup?

## Mental Model Utama

Sistem harus mampu menjelaskan apa yang terjadi, di mana waktu dihabiskan, dan kenapa request lambat atau gagal.

Tiga pilar observability memiliki fungsi spesifik:

- **Logs**
  → *Apa yang terjadi pada event tertentu?* (Discrete event log, konteks detail, error message)
- **Metrics**
  → *Bagaimana kondisi sistem secara agregat dari waktu ke waktu?* (Counters, gauge, p95/p99 latency)
- **Traces**
  → *Ke mana waktu sebuah request dihabiskan?* (Causal relationships, dependency breakdown, span hierarchy)

Ketiganya saling melengkapi, tetapi tidak *interchangeable*.

---

## Failure Scenario: The Opaque Invoice Flow

Alur pemrosesan invoice melibatkan beberapa dependensi I/O:

```
HTTP Request
    ↓
Load Invoice / Database
    ↓
Reserve Inventory
    ↓
Calculate Commission
    ↓
Generate PDF
    ↓
Send External Notification / WhatsApp
```

### Normal Execution Latencies

| Step | Operation | Latency |
|------|-----------|---------|
| 1 | Database Load | 20 ms |
| 2 | Reserve Inventory | 15 ms |
| 3 | Calculate Commission | 10 ms |
| 4 | Generate PDF | 30 ms |
| 5 | Notification | 20 ms |
| **Total** | | **~95 ms** |

### Slow PDF Outage Scenario

| Step | Operation | Latency |
|------|-----------|---------|
| 1 | Database Load | 20 ms |
| 2 | Reserve Inventory | 15 ms |
| 3 | Calculate Commission | 10 ms |
| 4 | **Generate PDF** | **4800 ms (Bottleneck)** |
| 5 | Notification | 20 ms |
| **Total** | | **~4865 ms** |

---

## Unsafe vs Safe Comparison

### Unsafe Implementation (`unsafe_service.go`)

```go
func (s *UnsafeInvoiceService) ProcessWithDeps(ctx context.Context, invoiceID string, deps Dependencies) error {
    start := time.Now()
    // ... executes database, inventory, commission, pdf, notification ...
    duration := time.Since(start)
    s.logf("INFO request completed duration_ms=%d", duration.Milliseconds())
    return nil
}
```

**Karakteristik Unsafe:**
1. **Opaque logs**: Hanya mencatat total durasi tanpa breakdown sub-operasi I/O.
2. **Tidak ada correlation**: Meneruskan `context.Context` tanpa mengaitkan `trace_id` atau `span_id` ke logs.
3. **Tidak ada dependency metrics & tracing**: Tidak ada p95 latency atau tracing span per dependency.

### Safe Implementation (`safe_service.go`)

```go
func (s *SafeInvoiceService) ProcessWithDeps(ctx context.Context, invoiceID string, deps Dependencies) error {
    ctx, rootSpan := s.Tracer.Start(ctx, "invoice.process",
        trace.WithAttributes(attribute.String("invoice_id", invoiceID)),
    )
    defer rootSpan.End()

    for _, step := range steps {
        stepCtx, stepSpan := s.Tracer.Start(ctx, step.component+"."+step.operation)
        err := step.execute(stepCtx)
        // ... record OTel spans, prometheus metrics, and structured logs ...
    }
}
```

**Keunggulan Safe:**
1. **Span Hierarchy**:
   ```
   POST /invoices/{id}/process
   └── invoice.process
       ├── database.load_invoice
       ├── inventory.reserve
       ├── commission.calculate
       ├── pdf.generate
       └── notification.send
   ```
2. **Prometheus Golden Signals**:
   - Traffic: `lab07_http_requests_total{method="POST",route="/invoices/{id}/process",status_class="2xx"}`
   - Latency: `lab07_http_request_duration_seconds` & `lab07_dependency_duration_seconds`
   - Errors: `lab07_http_request_errors_total` & `lab07_dependency_errors_total`
   - Saturation proxy: `lab07_http_in_flight_requests`
3. **Structured Slog Correlation**:
   ```json
   {
     "time": "2026-09-03T10:56:10Z",
     "level": "INFO",
     "msg": "dependency completed",
     "request_id": "demo-slow-pdf-001",
     "trace_id": "ece356994a7f660b01a9a2d75a9e4171",
     "span_id": "5d8fa7cc0d77c85e",
     "invoice_id": "INV-1001",
     "component": "pdf",
     "operation": "generate",
     "duration_ms": 4800,
     "outcome": "success"
   }
   ```
4. **W3C Trace Context Propagation**: Header `traceparent` diinjeksi ke panggilan HTTP downstream.

---

## Hubungan Golden Signals (Di Lab vs Production)

Meskipun di lab ini kita menggunakan simulasi (sleep/timer) dan local observability stack, metrik-metrik ini memetakan kondisi sebenarnya di production:

1. **Latency**
   - **Lab:** Waktu total request invoice & waktu eksekusi dependency (misal `pdf.generate` 4800ms).
   - **Production:** Request duration user-facing API & latency real-world (slow query, API third-party lelet).
2. **Traffic**
   - **Lab:** Jumlah hit HTTP Request ke endpoint.
   - **Production:** Request per Second (RPS) pada API Gateway / Service.
3. **Errors**
   - **Lab:** HTTP 5xx dan metric Dependency Error.
   - **Production:** Kegagalan payment, database connection reset, timeout ke external API.
4. **Saturation (Proxy)**
   - **Lab:** Seberapa banyak in-flight request pada aplikasi kita.
   - **Production:** CPU load penuh, kehabisan memori, antrean worker kepenuhan, max DB connection tercapai.

---

## Investigation Story: Troubleshooting Seperti Engineer Beneran

**Kasus:** Customer mengeluh, *"Kadang-kadang generate invoice lambat banget, sampai 20 detik."*

Sebagai engineer, kamu nggak perlu nebak. Ikuti jejak telemetri ini (kamu bisa coba lakukan dengan skenario `slow-pdf` di bawah):

1. **Metrics memberitahu ADA masalah:**
   Buka Grafana, lihat panel p95 Latency API naik drastis (dari 95ms jadi ~4.8s).
2. **Metrics mempersempit area:**
   Lihat panel dependency metric, *ternyata komponen PDF yang naik grafiknya*. Database & inventory normal.
3. **Tracing memberitahu DI MANA persisnya:**
   Buka Jaeger, filter request yang butuh > 4s. Buka Trace ID-nya. Terlihat jelas dari "waterfall" view bahwa span `pdf.generate` memakan 99% durasi.
4. **Correlation ID menghubungkan kepingan:**
   Dari Jaeger, copy `Trace ID` atau `Request ID` (`demo-slow-pdf-001`).
5. **Logs memberi detail APA yang terjadi:**
   Cari ID tersebut di Logs. Kamu akan melihat jejak yang ditinggalkan setiap fungsi, durasi spesifiknya per baris, dan apakah ada pesan ERROR khusus (misal "PDF template invalid").

**Expected Mental Model:**
- **Metrics** → "Ada masalah, latency sistem naik."
- **Tracing** → "Masalahnya di service PDF."
- **Logs** → "Detailnya, worker kehabisan memori saat render template X."
- **Correlation ID** → Tali yang mengikat semuanya jadi satu cerita utuh.

---

## Observability Bukan Sekadar Install Tools

Ingat: Observability **bukanlah** sekadar meng-install Prometheus, Grafana, atau Jaeger. Tools bisa ganti (Datadog, New Relic, ELK), tapi konsep dasarnya sama:

> **System → menghasilkan Telemetry (Logs/Metrics/Traces) → Engineer melakukan diagnosa internal state.**

---

## Common Observability Mistakes

Hindari anti-pattern berikut saat bekerja:
- **Logging semua hal (Verbose):** Noise terlalu tinggi, tagihan log storage bengkak, susah cari error sungguhan.
- **Log tanpa context (Opaque):** `"Error saving to DB"`. Request ID berapa? User ID siapa?
- **Semua log pakai ERROR:** Warning/Info dicatat sebagai ERROR, membuat alert fatigue (engineer jadi kebal/cuek dengan notifikasi error).
- **Nggak pakai Request ID:** Gagal melacak perjalanan satu request lintas service/goroutine.
- **Baru memikirkan monitoring setelah Production down:** Panik pas incident karena blind-spot.

---

## Metrics Cardinality

Metric labels dibatasi hanya pada nilai *bounded*: `method`, `route`, `status_class`, `component`, `operation`, `outcome`.

Dilarang menjadi label:
- `request_id`, `trace_id`, `span_id`
- `invoice_id`, `customer_id`, `user_id`
- `raw URL` atau `raw error` string
- `timestamp`

---

## W3C Trace Context & Request Correlation

Dua identifier memiliki fungsi berbeda:
- `request_id`: Correlation identifier yang diteruskan lewat HTTP header `X-Request-ID` untuk pencarian log end-to-end.
- `trace_id`: OpenTelemetry distributed tracing identifier yang diteruskan via W3C header `traceparent` untuk melacak span dan latency.

Keduanya bukan hal yang sama.

### Alur Propagasi:

```
Client request
X-Request-ID: demo-slow-ext-001
traceparent: ...

        ↓

Invoice API log
request_id=demo-slow-ext-001
trace_id=<same-trace>

        ↓

Notification HTTP request
X-Request-ID: demo-slow-ext-001
traceparent=<same-trace>

        ↓

Notification server log
request_id=demo-slow-ext-001
trace_id=<same-trace>
span_id=<downstream-server-span>
```

`HTTPNotificationClient` menggunakan OpenTelemetry TextMapPropagator resmi untuk menginjeksi header standar W3C:
```
traceparent: 00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01
```
serta header `X-Request-ID`. Memastikan trace context dan request correlation tidak terputus saat request keluar menuju layanan eksternal (service boundary disimulasikan dalam satu proses untuk keperluan lab).

---

## Error Correlation

Ketika sub-operasi gagal:
1. Record error pada dependency span & set status `Error`.
2. Record error pada `invoice.process` span.
3. Record error pada root HTTP span & set status `Error`.
4. Naikkan `lab07_dependency_errors_total` dan `lab07_http_request_errors_total`.
5. Log terstruktur dengan level `ERROR` yang memuat `request_id`, `trace_id`, dan detail error.
6. Kembalikan respons HTTP `502 Bad Gateway` (atau `504 Gateway Timeout`) dengan sanitized JSON body yang tetap menyertakan `request_id`.

---

## Context Cancellation

Operasi simulasi menggunakan timer yang dapat dihentikan (`time.NewTimer` + `defer timer.Stop()`). Ketika context dibatalkan atau timeout tercapai, timer segera dibatalkan dan I/O langsung dihentikan.

---

## Cara Menjalankan dengan Docker

Jalankan perintah berikut dari root repository:

```bash
docker compose up -d
```

*(Atau jika Anda sedang berada di dalam direktori `labs/07-observability`, gunakan: `docker compose -f ../../docker-compose.yml up -d`)*

### Endpoints Layanan
- **Lab 07 API**: `http://localhost:8087`
- **Prometheus Metrics**: `http://localhost:8087/metrics`
- **Prometheus Server**: `http://localhost:9090`
- **Grafana**: `http://localhost:3000` (admin/admin)
- **Jaeger UI**: `http://localhost:16686`

---

## Cara Membaca Observability Tools

### 1. Membaca Prometheus (`:9090`)
- Query p95 Latency: `histogram_quantile(0.95, sum(rate(lab07_dependency_duration_seconds_bucket[1m])) by (le, component, operation))`
- Query Error Rate: `sum(rate(lab07_dependency_errors_total[1m])) by (component, operation)`

### 2. Membaca Grafana (`:3000`)
- Buka dashboard: **Lab 07 - Observability**.
- Panel `Dependency p95 Latency` akan menunjukkan lonjakan pada komponen yang lambat (misal `pdf.generate` pada scenario `slow-pdf`).

### 3. Membaca Jaeger (`:16686`)
- Pilih Service: `lab07-observability`.
- Klik **Find Traces**.
- Buka trace untuk melihat hierarki `POST /invoices/{id}/process` -> `invoice.process` -> child spans beserta durasi masing-masing.

---

## Cara Menjalankan Test

```bash
# Jalankan unit & integration tests
make lab-07-test

# Jalankan dengan race detector
make lab-07-test-race

# Linter & Vet
make lab-07-vet
```

---

## Production Considerations & Trade-offs

- **Sampling Rate**: Tracing 100% traffic pada traffic tinggi membebani storage/network; gunakan probabilistic sampling (misal 5-10%) pada production.
- **Log Volume**: Batasi log INFO pada level dependency untuk hot path, atau turunkan ke DEBUG.
- **Metric Cardinality**: Hindari penggunaan user-generated identifiers sebagai labels untuk mencegah ledakan memori TSDB.

---

### Korelasi Signal Saat Diagnosa Bottleneck

| Signal | Yang Terlihat |
|---|---|
| **Log** | Request selesai sekitar 4,8 detik |
| **HTTP metric** | Request latency meningkat |
| **Dependency metric** | Latency komponen `pdf` meningkat |
| **Trace** | `pdf.generate` memakai hampir seluruh waktu |
| **Request ID** | Log request yang sama dapat dicari |
| **Trace ID** | Log dapat dibuka sebagai distributed trace |
| **Span ID** | Log dapat diarahkan ke dependency span tertentu |

---

## Command Reproduksi & Manual Testing

Jalankan perintah berikut untuk menguji berbagai skenario I/O:

```bash
# 1. Normal scenario (~95ms)
# Pembelajaran: Baseline system yang sehat
curl -i -X POST \
  -H "X-Request-ID: demo-normal-001" \
  "http://localhost:8087/invoices/INV-1001/process?scenario=normal"

# 2. Slow PDF scenario (~4800ms bottleneck)
# Pembelajaran: Mendeteksi bottleneck pada application/internal processing (CPU/Mem heavy task)
curl -i -X POST \
  -H "X-Request-ID: demo-slow-pdf-001" \
  "http://localhost:8087/invoices/INV-1001/process?scenario=slow-pdf"

# 3. Slow Database scenario (~3000ms bottleneck)
# Pembelajaran: Mendeteksi bottleneck pada internal infrastructure dependency (slow query/DB contention)
curl -i -X POST \
  -H "X-Request-ID: demo-slow-db-001" \
  "http://localhost:8087/invoices/INV-1001/process?scenario=slow-database"

# 4. Slow External Notification scenario (~3500ms downstream HTTP delay)
# Pembelajaran: Mendeteksi bottleneck pada third-party API latency (melewati HTTP network boundary via HTTPNotificationClient)
curl -i -X POST \
  -H "X-Request-ID: demo-slow-ext-001" \
  "http://localhost:8087/invoices/INV-1001/process?scenario=slow-external"

# 5. PDF Error scenario (502 Bad Gateway)
# Pembelajaran: Memeriksa error telemetry (error span ter-record, http 5xx, structured error log)
curl -i -X POST \
  -H "X-Request-ID: demo-pdf-err-001" \
  "http://localhost:8087/invoices/INV-1001/process?scenario=pdf-error"
```

---

## Out of Scope

- Alerting & PagerDuty integration
- SLO/SLI engineering
- Log aggregation cluster (Elasticsearch / Loki)
- Tail-based sampling & massive-scale sampling strategy
- Multi-region telemetry
- Incident response runbooks
- Capacity planning
- Full APM platform setup

---

## Navigasi

- **Previous**: [Lab 06 — API Versioning](../06-api-versioning/)
- **Next**: [Lab 08 — Database Isolation Level](../08-database-isolation-level/)