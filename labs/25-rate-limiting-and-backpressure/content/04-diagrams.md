# Diagrams

Semua diagram diturunkan dari implementasi dan arsitektur yang terverifikasi. Tidak ada komponen yang ditemukan selain yang ada di kode.

## Diagram 1 — Arsitektur Lapisan (dari engineering/01-design.md)

```text
[Incoming Request]
        │
        ▼
[HTTP Middleware / Tenant Extractor] ── (RFC 6598 Tenant Key via X-API-Key)
        │
        ▼
[Token Bucket / Rate Limiter] ── (Exceeded? → HTTP 429 + Retry-After)
        │ (Passed)
        ▼
[Bounded Queue / Backpressure Channel] ── (Queue Full? → 503 Overloaded)
        │ (Enqueued)
        ▼
[Worker Pool / Consumer] ── (Little's Law L = λW steady state)
        │
        ▼
[Client with Full Jitter Retry] ── (AWS Jitter Backoff)
```

Komponen nyata dalam kode:
- `internal/httputil/middleware.go` — HTTP Middleware
- `internal/ratelimit/bucket.go` + `internal/ratelimit/registry.go` — Token Bucket
- `internal/backpressure/queue.go` — Bounded Queue
- `internal/retry/backoff.go` — Jitter Backoff
- `cmd/demo/main.go` — CLI demo

## Diagram 2 — Token Bucket State Transitions

```text
       Capacity (B)
   ┌──────────────┐
   │ █████████████│ ← Bucket Full (burst allowed)
   │ ████████░░░░░│ ← Partial
   │ ███░░░░░░░░░░│ ← Low
   │ ░░░░░░░░░░░░░│ ← Empty (reject until refill)
   └──────────────┘
         ▲
         │  refill tokens at rate R per second
         │  tokens += elapsed × R
         │  tokens = min(tokens, B)
```

Berdasarkan `internal/ratelimit/bucket.go:34-38` — `AllowN()` replenishes berdasarkan elapsed time dan meng-cap ke capacity.

## Diagram 3 — Leaky Bucket State Transitions

```text
   Capacity (B)
   ┌──────────────┐
   │ █████████████│ ← Reject (water + 1 > B)
   │ ████████░░░░░│ ← Partial
   │ ███░░░░░░░░░░│ ← Draining
   │ ░░░░░░░░░░░░░│ ← Empty
   └──────────────┘
         ▲
         │  leak water at rate R per second
         │  water -= elapsed × R
         │  water = max(water, 0)
```

Berdasarkan `internal/ratelimit/bucket.go:104-114` — `Allow()` mengurangi water dengan elapsed × leakRate, menolak jika `water + 1 > capacity`.

## Diagram 4 — Bounded Queue Backpressure

```text
    Producer          BoundedQueue          Worker Pool
   ┌──────┐        ┌──────────┐        ┌──────────┐
   │      │──Try──▶│  buffer  │──deque──▶│ worker 1 │
   │      │  Submit│  [cap=N] │        │ worker 2 │
   │      │        │          │        │ worker 3 │
   └──────┘        └──────────┘        └──────────┘
     │                    │
     │  if full:          │  context cancel /
     │  return ErrQueueFull│  Stop()
     ▼                    ▼
  Fast rejection     Graceful shutdown
```

Berdasarkan `internal/backpressure/queue.go`:
- `TrySubmit()` (line 67-85): pola `select-default` → fast rejection
- `workerLoop()` (line 49-63): `select` antara `ctx.Done()` dan channel receive
- `Stop()` (line 91-103): cancel context, close channel, tunggu worker selesai

## Diagram 5 — Full Jitter Distribution

```text
Sleep Duration
    │
    │  ██
    │  ██  ██
    │  ██  ██  ██
    │  ██  ██  ██  ██
    └─────────────────────▶ Attempt
    Attempt 0  Attempt 1  Attempt 2  Attempt 3
    [0, base]  [0, 2×base]  [0, 4×base]  [0, 8×base]
    (capped at Cap)
```

Berdasarkan `internal/retry/backoff.go:36-41` — `FullJitter` menghasilkan `rand.Float64() × min(cap, base × 2^attempt)`. Interval tersebar merata dari 0 ke cap eksponensial, berbeda dengan `NoJitter` yang selalu menghasilkan nilai maksimum.

## Diagram 6 — Retry Budget (Google SRE Model)

```text
   ┌──────────────────┐
   │  Retry Budget    │
   ├──────────────────┤
   │ Per-request:     │
   │   max 3 attempts │
   ├──────────────────┤
   │ Per-client:      │
   │   max 10% retries│
   └──────────────────┘
         ▲
         │  Google SRE: Handling Overload (Source 10)
         │  Mencegah retry storm saat cascading overload
```

Berdasarkan `research/05-report.md` Finding 9 — per-request budget max 3 attempts, per-client retry ratio below 10%.
