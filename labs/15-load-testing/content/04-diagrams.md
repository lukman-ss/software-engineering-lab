## Load Test Diagram
```text
+-------------------------------------------------------------+
|                     Load Test Runner                        |
|                                                             |
|  [VU 1] -----> HTTP POST /booking                           |
|  [VU 2] -----> HTTP POST /booking                           |
|  ...                                                        |
|  [VU N] -----> HTTP POST /booking                           |
|                                                             |
|  (Tiap VU mencatat latensi ke memory slice independen)     |
+--------------------------------+---------------------------+
                               |
                               v
+-------------------------------------------------------------+
|                  HTTP Server (/booking)                     |
|                                                             |
|        +-------------------------------------------+        |
|        | Semafor Penampung Koneksi (Max = 5 Slot)  |        |
|        +-------------------------------------------+        |
|            | Slot 1 | Slot 2 | Slot 3 | Slot 4 | Slot 5     |
|                                                             |
|   (Request ke-6 dan seterusnya tertahan mengantre)          |
|   (Pemrosesan simulasi query: 20ms, atau 500ms              |
|    dengan 10% probabilitas saat di atas kapasitas pool)        |
+-------------------------------------------------------------+
```

## Latency Distribution Under Different Load Conditions

```text
Smoke Test (2 VUs, 5 DB connections)
P95 ≈ P50 ≈ Average ≈ 21ms
No queue buildup — all requests served immediately.

Stress Test (50 VUs, 5 DB connections)
P95 ≈ 982ms, P99 ≈ 1.2s
Queue buildup — 45 VUs wait while 5 process.
Tail latency spikes due to:
  1. Cumulative wait time in queue (primary mechanism)
  2. 10% probability of 25x query duration spike (additional amplifier on top of queuing)
```
