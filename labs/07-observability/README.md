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

### Unsafe Implementation (`unsafe.go`)

```go
func (p *UnsafeInvoiceProcessor) ProcessInvoice(ctx context.Context, invoiceID string) error {
    start := time.Now()
    // ... executes database, inventory, commission, pdf, notification ...
    duration := time.Since(start)
    p.logf("INFO request completed duration_ms=%d", duration.Milliseconds())
    return nil
}
```

**Kelemahan Unsafe:**
1. **Opaque logs**: Hanya mencatat `duration_ms=4865`. Developer tidak tahu apakah DB, PDF, atau Notification yang lambat.
2. **Tidak ada dependency metrics**: Prometheus tidak dapat mengagregasi p95 per dependency.
3. **Tidak ada span tree / tracing**: Tidak dapat melihat korelasi antar I/O sub-operations.
4. **Tidak ada context correlation**: Request ID tidak dipropagasi ke child operations.

### Safe Implementation (`safe.go`)

```go
func (p *SafeInvoiceProcessor) ProcessInvoice(ctx context.Context, invoiceID string) error {
    ctx, rootSpan := p.Tracer.Start(ctx, "invoice.process")
    defer p.Tracer.End(rootSpan)
    
    // Setiap step diinstrumentasi dengan span, metric, dan structured logs
    stepCtx, span := p.Tracer.Start(ctx, "pdf.generate")
    err := p.PDF.Generate(stepCtx, invoiceID)
    p.Tracer.End(span)
    // ...
}
```

**Keunggulan Safe:**
1. **Child Spans**: Setiap tahapan memiliki span tersendiri (`database.load_invoice`, `pdf.generate`, dll).
2. **Prometheus Metrics**: Mengukur counter `dependency_requests_total` dan latency histogram/duration per status.
3. **Structured Slog Correlation**: Setiap log otomatis mengikutsertakan `request_id`, `trace_id`, `span_id`, dan `invoice_id`.
4. **Context Propagation & Cancellation**: Mendukung timeout dan context cancel secara presisi di setiap boundary I/O.

---

## Dependency Scenarios

Lab ini menyediakan fake dependency simulator yang mendukung:
- `normal`
- `slow-pdf`
- `slow-database`
- `slow-external`
- `pdf-error`
- `context-cancellation`

---

## Cara Menjalankan Test

```bash
# Jalankan unit & scenario tests
go test -v ./labs/07-observability/...

# Jalankan dengan race detector
go test -race -v ./labs/07-observability/...
```

---

## Navigasi

- **Previous**: [Lab 06 — API Versioning](../06-api-versioning/)
- **Next**: [Lab 08 — Database Isolation Level](../08-database-isolation-level/)
