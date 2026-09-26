## 1. Connection Pooling Flow (Basic)

```text
┌─────────────────────────────────────────────────────────────────────┐
│  Aplikasi Konkuren (banyak goroutines / threads)                     │
│  ├─ Request A → ──────────────────────────────────┐                 │
│  ├─ Request B → ─────────────────────────────┐    │                 │
│  ├─ Request C → ───────────────────────┐    │    │                 │
│  └─ Request D → ──────────────────┐    │    │    │                 │
└────────────────────────────────────┴────┴────┴────┴─────────────────┘
                                     ▼    ▼    ▼    ▼
                            [ Client Connection Pool ]
                            (maxOpenConns = 5)
                                     │    │    │
                             (reuse) │    │    │
                                     ▼    ▼    ▼
                          [ Database Backend ]
                          (max_connections = 15)
                          ┌─ Connection 1: active
                          ├─ Connection 2: active
                          └─ Connection 3: active
```

**Key Points:**
- Pool membatasi jumlah koneksi maksimal yang bisa dibuka ke database
- Koneksi dipakai ulang, bukan dibuat tutup setiap request
- Database memiliki batas hard `max_connections` (misal: 15)
- Jika client pool + instance aplikasi melebihi limit, server menolak

---

## 2. Connection Leak & Starvation Flow

```text
┌─────────────────────────────────────────────────────────────────────┐
│  Goroutine 1 (Unsafe Pattern)                                        │
│  ├─ db.Conn(ctx)                          [Acquired]                 │
│  ├─ Exec("UPDATE orders...")              [Active]                   │
│  ├─ externalCall() ← SLOW NETWORK CALL              [HOLDING]       │
│  │     Sleep(500ms)                                                  │
│  └─ (connection tidak di-release sampai selesai)                    │
└─────────────────────────────────────────────────────────────────────┘
                                             ▼
                                    [Pool Habis: 2/2]
                                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│  Goroutine 2 (Safe Pattern)                                          │
│  ├─ db.Conn(ctx)                          [Waiting...]               │
│  └─ context timeout (100ms) ← Deadline Exceeded                    │
└─────────────────────────────────────────────────────────────────────┘
```

**Root Cause:** Eksekusi I/O eksternal (panggil API pihak ketiga) dilakukan dalam blok yang masih memegang koneksi database.

---

## 3. Oversized Pool Exhaustion Flow

```text
Database Server: max_connections = 15
Client Pool Size: 50 (terlalu besar)

Concurrent Requests (30)
├─ Request 1..15  →  Acquired  (success)
├─ Request 16..30 →  Rejected (ErrServerOverloaded)
                   └─ "FATAL: sorry, too many clients already"
```

**Reality Check:** Kesalahan umum adalah memperbesar pool sisi client tanpa mempertimbangkan:
- Jumlah instance aplikasi (autoscaling)
- Total potensi koneksi = `instance_count × pool_size`
- Reserved connections untuk superuser/monitoring (biasanya 15)

---

## 4. Pool Sizing Formula

```text
Optimal Pool Size ≈ (core_count × 2) + effective_spindle_count

Contoh:
- 4-core CPU, 1 HDD → (4 × 2) + 1 = 9 connections
- 8-core CPU, NVMe SSD → (8 × 2) + 0 = 16 connections (SSD = no seeks)
```

**Catatan dari PostgreSQL Wiki & HikariCP:**
- Core count TIDAK termasuk hyperthreading threads
- effective_spindle_count = 0 jika data fully cached (RAM)
- Formula ini adalah titik awal, bukan solusi mutlak
- SSD: less blocking → fewer threads lebih baik (effect_spindle ≈ 0)