# Diagrams

## Diagram 1 — Load Test Type Taxonomy

```
Performance Testing Types
-----------------------------------------------------------
| Type      | Pattern              | Purpose                   |
|-----------|----------------------|---------------------------|
| Smoke     | 2-20 VUs, seconds    | Script validation         |
| Load      | Ramp-up → Plateau    | Baseline/Typical traffic  |
| Stress    | Above normal         | Find breaking point       |
| Spike     | Sudden high load     | Flash sale, launch        |
| Soak      | Hours/Days load      | Memory/resource leaks     |
| Breakpoint| Gradual increase     | Capacity boundary         |
-----------------------------------------------------------
```

## Diagram 2 — Server Semaphore Workflow

```
Request Arrives
      │
      ▼
[Acquire Slot] ──┐
                  │  Full?
      │          └── Yes → Wait in queue
      ▼
[Process Request] (DBQueryDuration ~ 20ms)
      │
      ▼
[Release Slot]
      │
      ▼
Response Sent
```

`MaxDBConnections = 5` → Satu sampai 5 request proses paralel; request ke-6 menunggu.

## Diagram 3 — Latency Distribution: Smoke vs Stress

```
Response Time Distribution
--------------------------------------------------
Smoke Test (2 VUs vs 5 DB slots)
  Count: 188 requests
  Latency:
    Min ────┐
    P50 ────┼───┐
    P95 ────┴───┼───────┐     (all clustered)
    Max ────────┴───┐   │
                    │   ▼
                    └── 21-22ms

Stress Test (50 VUs vs 5 DB slots)
  Count: 120 requests
  Latency:
    Min ────┐
    P50 ────┼───────────┐
    Avg ────┴────┐      │
    P95 ─────────┼──────┼──────────┐
    P99 ─────────┴──────┼──────┐   │
                       │      ▼   │
                       │    Queue │
                       │  wait   │
                       └─────────┴──► 700ms-1.5s
--------------------------------------------------
```

## Diagram 4 — Latency Breakdown (HTTP Request Lifecycle)

```
http_req_duration = client-side metrics
  ├─ http_req_blocked      (waiting for idle connection)
  ├─ http_req_connecting   (TCP handshake)
  ├─ http_req_tls_handshaking (SSL negotiation)
  ├─ http_req_sending      (request body)
  ├─ http_req_waiting      (server processing - TTFB)
  └─ http_req_receiving    (response body)

Diagnosis:
- ↑ http_req_waiting  →  App/DB bottleneck
- ↑ http_req_connecting →  Network/conn issue
- ↑ http_req_tls_handshaking →  SSL cert problem

Catatan server lab: ketika `activeReq > MaxDBConnections`, 10% request mendapat query 25x lebih lambat — amplifier di atas antrean, bukan penyebab utama; stress test tetap menunjukkan degradasi P95 dari penumpukan antrean saja.
```

## Diagram 5 — Load Generator Architecture

```
                    ┌─────────────────────┐
                    │   Load Tester       │
                    │   (Runner)          │
                    │                     │
                    └─┬─┬─┬─┬─┬───────────┘
                      │ │ │ │ │
          ┌───────────┘ │ │ │ └──────────┐
          ▼             ▼ ▼ ▼             ▼
    ┌───────┐     ┌───────┐       ┌───────┐
    │  VU 1 │ ... │  VU N │       │  VU M │
    └───────┘     └───────┘       └───────┘
      │               │               │
      ▼               ▼               ▼
  POST /booking   POST /booking   POST /booking
                      │
                      ▼
              ┌─────────────────┐
              │   Server        │
              │   (semaphore)   │
              ├─────────────────┤
              │ MaxDBConnections│
              │      5 slots    │
              └─────────────────┘
```

Setiap VU memiliki koneksi HTTP terisolasi; server menyebar request ke dalam semaphore 5-slot.

## Diagram 6 — Bottleneck Identification Flow

```
Load Test Running
       │
       ▼
┌─────────────────┐
│ Collect Metrics │
│  - P50, P95, P99 │
│  - Error Rate    │
│  - CPU, Memory   │
└─────────────────┘
       │
       ▼
┌─────────────────┐
│ Is P95 > Target?│────Yes───► Check Resource Utilization
└─────────────────┘
       │ No
       ▼
       OK
       
Check Resource Utilization:
  CPU > 75%?  ──Yes──► CPU bottleneck
  Mem > 80%?  ──Yes──► Memory bottleneck
  DB conn = 0 free? ──Yes──► Connection pool exhausted
  
If http_req_waiting high:
  → App/DB processing slow
If http_req_connecting high:
  → Network latency
```

## Diagram 7 — Metrics Calculation Flow

```
Load Test Runner
       │
       ▼
┌─────────────────┐
│ Per-VU Request  │
│ loop            │
│  ├─ record start │
│  ├─ Do(req)      │
│  └─ record latency=now-start
└─────────────────┘
       │
       ▼ (all VUs done → WaitGroup)
       │
       ▼
┌─────────────────┐
│ Aggregate all   │
│ latencies       │
└─────────────────┘
       │
       ▼
┌─────────────────┐
│ CalculateMetrics│
│  1. sort()      │
│  2. min/max     │
│  3. avg         │
│  4. percentile  │
│     idx = (n-1)*p/100
└─────────────────┘
       │
       ▼
       Result{Total, Success, Errors, P50, P95, P99, ...}
```