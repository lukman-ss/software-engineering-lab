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
func (s *UnsafeInvoiceService) Process(ctx context.Context, invoiceID string) error {
    start := time.Now()
    // ... executes database, inventory, commission, pdf, notification ...
    duration := time.Since(start)
    s.logf("INFO request completed duration_ms=%d", duration.Milliseconds())
    return nil
}
```

**Kelemahan Unsafe:**
1. **Opaque logs**: Hanya mencatat `INFO request completed duration_ms=4865`. Developer tidak tahu apakah DB, PDF, atau Notification yang lambat.
2. **Tidak ada dependency metrics**: Prometheus tidak dapat mengagregasi p95 per dependency.
3. **Tidak ada span tree / tracing**: Tidak dapat melihat korelasi antar I/O sub-operations.
4. **Tidak ada context correlation**: Request ID tidak dipropagasi ke child operations.

### Safe Implementation (`safe_service.go`)

```go
func (s *SafeInvoiceService) Process(ctx context.Context, invoiceID string) error {
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
1. **Child Spans**: Setiap tahapan memiliki span tersendiri (`database.load_invoice`, `pdf.generate`, dll).
2. **Prometheus Metrics (Golden Signals)**:
   - Traffic: `lab07_http_requests_total{method="POST",route="/invoices/{id}/process",status_class="2xx"}`
   - Latency: `lab07_http_request_duration_seconds` & `lab07_dependency_duration_seconds`
   - Errors: `lab07_http_request_errors_total` & `lab07_dependency_errors_total`
   - Saturation: `lab07_http_in_flight_requests`
3. **Structured Slog Correlation**:
```json
{
  "time": "2026-09-03T10:56:10Z",
  "level": "INFO",
  "msg": "dependency completed",
  "request_id": "demo-slow-pdf-001",
  "trace_id": "ece356994a7f660b01a9a2d75a9e4171",
  "span_id": "06111ee7e4a6daae",
  "component": "pdf",
  "operation": "generate",
  "duration_ms": 4800,
  "outcome": "success"
}
```
4. **Context Propagation & Cancellation**: Request ID (`X-Request-ID`) dipropagasi dan dipertahankan.

---

## Logs dan Trace Correlation: Alur Diagnosis

```
cari request_id pada log
    ↓
temukan trace_id
    ↓
buka trace
    ↓
lihat child span yang lambat (pdf.generate)
```

---

## Menjalankan Demo Server

Jalankan server demo:

```bash
go run ./labs/07-observability/cmd/demo
```

Server berjalan pada port `:8087`.

### Scenario Endpoints (Safe)

- `POST http://localhost:8087/invoices/INV-101/process?scenario=normal`
- `POST http://localhost:8087/invoices/INV-101/process?scenario=slow-pdf`
- `POST http://localhost:8087/invoices/INV-101/process?scenario=slow-database`
- `POST http://localhost:8087/invoices/INV-101/process?scenario=slow-external`
- `POST http://localhost:8087/invoices/INV-101/process?scenario=pdf-error`

### Unsafe Endpoint

- `POST http://localhost:8087/unsafe/invoices/INV-101/process?scenario=slow-pdf`

### Prometheus Metrics

- `GET http://localhost:8087/metrics`

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

## Navigasi

- **Previous**: [Lab 06 — API Versioning](../06-api-versioning/)
- **Next**: [Lab 08 — Retry](../08-retry/)
