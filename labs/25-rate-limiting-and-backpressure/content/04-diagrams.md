# Diagrams

## System Flow

```
[Incoming Request]
        │
        ▼
[HTTP Middleware / Tenant Extractor]
        │ (X-API-Key / "anonymous")
        ▼
[Token Bucket / Rate Limiter]
        │
        ├─> Gagal? ──YES──> [HTTP 429 + Retry-After]
        │
        ▼ (Lulus)
[Bounded Queue / Backpressure]
        │
        ├─> Penuh? ──YES──> [ErrQueueFull / 503]
        │
        ▼ (Dienqueue)
[Worker Pool / Consumer]
        │ (L = λW steady state)
        ▼
[Job Selesai]
```

## Token Bucket State Machine

```
┌─────────────────────────────────────────────────────────────┐
│                        Token Bucket                         │
│                                                             │
│  Capacity: B         Refill Rate: R tokens/detik            │
│                                                             │
│  ┌──────────────┐    Allow()       ┌──────────────┐         │
│  │  Tokens = B  │ ────────────────>│  Tokens >= 1 │─ YES ──>│
│  └──────────────┘                  └──────────────┘         │
│       │                                   │                  │
│       │ refill                            │                  │
│       ▼                                   ▼                  │
│  ┌──────────────┐                     ┌──────────────┐      │
│  │  Tokens < 1  │◄───────────────────│  Tokens -= 1 │      │
│  └──────────────┘                     └──────────────┘      │
│       │                                   │                  │
│       └───────────────── NO ──────────────┘                  │
│                                                             │
│                └─> Reject & RetryAfterSeconds()             │
└─────────────────────────────────────────────────────────────┘
```

## Leaky Bucket State Machine

```
┌─────────────────────────────────────────────────────────────┐
│                        Leaky Bucket                         │
│                                                             │
│  Capacity: C         Leak Rate: R units/detik               │
│                                                             │
│  ┌──────────────┐    Allow()       ┌──────────────┐         │
│  │ Water = C    │ ────────────────>│ Water + 1 ≤ C│─ YES ──>│
│  └──────────────┘                  └──────────────┘         │
│       │                                   │                  │
│       │ leak                              │                  │
│       ▼                                   ▼                  │
│  ┌──────────────┐                     ┌──────────────┐      │
│  │ Water < C    │◄───────────────────│ Water += 1   │      │
│  └──────────────┘                     └──────────────┘      │
│       │                                   │                  │
│       └───────────────── NO ──────────────┘                  │
│                                                             │
│                └─> Reject                                   │
└─────────────────────────────────────────────────────────────┘
```

## Bounded Queue Backpressure

```
┌─────────────────────────────────────────────────────────────┐
│                      Bounded Queue                          │
│                                                             │
│  Capacity: 3         Workers: 1                             │
│                                                             │
│  ┌───────────────────────────────────────────────────┐      │
│  │  ┌─────┐ ┌─────┐ ┌─────┐                          │      │
│  │  │Job 1│ │Job 2│ │Job 3│  [Channel Buffer]        │      │
│  │  └─────┘ └─────┘ └─────┘                          │      │
│  └─────────────────────���─────────────────────────────┘      │
│          │                   │                              │
│          ▼                   ▼                              │
│    ┌─────────────┐     ┌─────────────┐                     │
│    │  Worker #1  │     │  Worker #2  │   ...               │
│    └─────────────┘     └─────────────┘                     │
│                                                             │
│  TrySubmit() → Full? ──YES──> [ErrQueueFull] (Fast Drop)   │
│  TrySubmit() → Empty? ──NO──> [Enqueue]                    │
└─────────────────────────────────────────────────────────────┘
```

## Retry Backoff Distribution (FullJitter)

```
Attempt 0: base=100ms, cap=1000ms
├─────────────────────────────────────────────────────────────┤
│[random(0, 100ms)]                                           │
└─────────────────────────────────────────────────────────────┘

Attempt 1: base=100ms, cap=1000ms, 2^1=200ms
├─────────────────────────────────────────────────────────────┤
│[random(0, 200ms)]                                           │
└─────────────────────────────────────────────────────────────┘

Attempt 2: base=100ms, cap=1000ms, 2^2=400ms
├─────────────────────────────────────────────────────────────┤
│[random(0, 400ms)]                                           │
└─────────────────────────────────────────────────────────────┘

Attempt 3: base=100ms, cap=1000ms, 2^3=800ms
├─────────────────────────────────────────────────────────────┤
│[random(0, 800ms)]                                           │
└─────────────────────────────────────────────────────────────┘

Attempt 4+: capped at cap=1000ms
├─────────────────────────────────────────────────────────────┤
│[random(0, 1000ms)]                                          │
└─────────────────────────────────────────────────────────────┘
```

## Little's Law Capacity Planning

```
Little's Law: L = λ × W

L = Average number in system (queue depth)
λ = Arrival rate (requests per second)
W = Average time in system (service time)

Contoh:
λ = 2000 requests/detik
W = 1.25 detik (pemrosesan + antrian)
L = 2000 × 1.25 = 2500 jobs dalam sistem

Jika λ meningkat menjadi 2500/detik:
- Queue depth meningkat
- W meningkat (antrian lebih panjang)
- Latency meningkat

Jika λ melebihi service rate R:
- Queue tumbuh tak terbatas
- Sistem runtuh

Mitigasi:
- Rate limiting → turunkan λ masuk
- Scaling → naikkan service rate R
```

## Tenant Isolation

```
Tenant A (API Key: "tenant-a")          Tenant B (API Key: "tenant-b")
          │                                       │
          ▼                                       ▼
  ┌─────────────────┐                     ┌─────────────────┐
  │   Token Bucket  │                     │   Token Bucket  │
  │  Capacity: 10   │                     │  Capacity: 10   │
  │  Refill: 5/s    │                     │  Refill: 5/s    │
  └─────────────────┘                     └─────────────────┘
          │                                       │
          ▼                                       ▼
   [Registry] ──────────────────────────────────> [Tenant Isolation]
          │
          ▼
   IP-only = ❌ (CGNAT affects many users)
   API key = ✓ (Per-tenant bucket, fair sharing)
```

## Complete System Flow with Timing

```
Time Sequence (t0 → t1000ms):

t0:   [R1] ──> Token Bucket (3 tokens) ──> [Allow, tokens=2]
t1:   [R2] ──> Token Bucket (2 tokens) ──> [Allow, tokens=1]
t2:   [R3] ──> Token Bucket (1 tokens) ──> [Allow, tokens=0]
t3:   [R4] ──> Token Bucket (0 tokens) ──> [429, Retry-After=1]
t4:   [R5] ──> Token Bucket (0 tokens) ──> [429, Retry-After=1]

t500: [R6] ──> Token Bucket (refill ~2.5) ──> [Allow, tokens=1.5]

t0:   [J1] ──> Bounded Queue (cap=3) ──> [Enqueue] ──> [Worker processes]
t1:   [J2] ──> Bounded Queue (cap=3) ──> [Enqueue]
t2:   [J3] ──> Bounded Queue (cap=3) ──> [Enqueue]
t3:   [J4] ──> Bounded Queue (full) ───> [ErrQueueFull]
t4:   [J5] ──> Bounded Queue (full) ───> [ErrQueueFull]
t5:   [J6] ──> Bounded Queue (full) ───> [ErrQueueFull]

t0:   Client retry ──> FullJitter(0ms)     ──> [Wait 0ms, retry]
t1:   Client retry ──> FullJitter(100ms)   ──> [Wait ~50ms, retry]
t2:   Client retry ──> FullJitter(200ms)   ──> [Wait ~150ms, retry]
t3:   Client retry ──> FullJitter(400ms)   ──> [Wait ~300ms, retry]
```