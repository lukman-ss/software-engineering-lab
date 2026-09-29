## Diagram 1 — Struktur Bit Array dengan Double Hashing

```
+--------------------------------------------------------------------+
| Bit Array (m bits)                                                 |
|  [0] [1] [2] ... [m-1]                                             |
|                                                                     |
|  Key X di-hash ke h1(X) dan h2(X)                                  |
|  Probe i:  pos_i = (h1 + i * h2) mod m                             |
|                                                                     |
|  i=0:  pos_0 = h1 mod m    [xxxxx|xxxxxx]                          |
|  i=1:  pos_1 = (h1+h2) mod m   [xxx|xxxxxx]                        |
|  i=2:  pos_2 = (h1+2*h2) mod m  [xxxx|xxxxxx]                      |
|  ...                                                                |
+--------------------------------------------------------------------+
```

## Diagram 2 — Alur Lookup pada LSM‑Tree dengan Filter

```
Query key → Segment N (newest) → Check filter
                                ├─ Filter says "definitely not" → Skip (no disk read)
                                └─ Filter says "possibly" → Load data block (disk read)
                          Segment N-1 → Check filter
                                ├─ Skip / Read ...
                          ...
                          Segment 1 (oldest) → Check filter
```

## Diagram 3 — Cache Penetration Prevention

```
Request untuk key absent
        │
        ▼
┌─────────────────┐
│ Cache + Bloom   │── filter.Check(key) == false ──► Reject (0 backend call)
│    Filter Gate  │
│                 │── filter.Check(key) == true ──► Check in‑memory map
│  in‑memory map  │── cache hit ──► Return value
│                 │── cache miss ──► Backend call (Fetch)
└─────────────────┘
```

## Diagram 4 — Hubungan Komponen Lab

```
┌──────────────┐      ┌──────────────────┐      ┌──────────────┐
│ cmd/demo     │─────▶│ internal/bloom   │◀─────│ internal/    │
│ (main.go)    │      │ Filter / Sync    │      │ store        │
└──────────────┘      │ (core logic)     │      │ LSMStore     │
                      └──────────────────┘      │ (segments)   │
                                                └──────────────┘
                                                      │
                                                      ▼
                                                ┌──────────────┐
                                                │ store.Cache  │
                                                │ (with filter │
                                                │  gate)       │
                                                └──────────────┘
```

## Diagram 5 — Parameter Optimal untuk ε = 0.01

```
  n = 10,000 elements
  ε = 0.01 (1% FP rate)

  m = ceil(-10000 * ln(0.01) / (ln 2)^2) ≈ 95,851 bits ≈ 11.7 KB
  k = round((95851/10000) * ln 2)        ≈ 7 hash functions
  Bits/element ≈ 9.59
```