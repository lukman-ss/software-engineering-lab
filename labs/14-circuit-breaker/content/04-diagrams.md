# Diagram Arsitektur dan Alur Circuit Breaker

Diagram-diagram di bawah ini merefleksikan model arsitektur dan mesin status yang secara riil diimplementasikan dan diverifikasi dalam pengujian Lab 14.

---

## 1. Diagram Mesin Status (State Machine Transitions)

```text
               ┌────────────────────────────────────────────────────────┐
               │                                                        │
               │                   Sukses /                             │
               │                   Reset Gagal                          │
               │                                                        │
               ▼                                                        │
       ┌───────────────┐     failures >= FailureThreshold       ┌───────────────┐
       │               │ ─────────────────────────────────────► │               │
       │    CLOSED     │                                        │     OPEN      │
       │  (Terkunci)   │ ◄───────────────────────────────────── │  (Terbuka)    │
       │               │           Probe Berhasil               │               │
       └───────────────┘       (consecutiveSuccesses >=         └───────────────┘
                                   HalfOpenMaxCalls)                    │
                                                                        │
                                                                        │ Waktu Istirahat
                                                                        │ Berlalu
                                                                        │ (now - lastStateChange
                                                                        │  >= OpenTimeout)
                                                                        ▼
                                                                ┌───────────────┐
                                                                │               │
                                                                │   HALF-OPEN   │
                                                                │  (Uji Coba)   │
                                                                │               │
                                                                └───────────────┘
                                                                        │
                                                                        │ Probe Gagal
                                                                        │ (err != nil)
                                                                        ▼
                                                                [Kembali ke OPEN]
```

---

## 2. Diagram Komponen Sistem

```text
+-------------------------+
|     Layanan Checkout    |
| (internal/checkout)     |
+-------------------------+
             │
             │ Panggilan Checkout(ctx)
             ▼
+---------------------------------------------------+
|               Circuit Breaker Proxy               |
|            (internal/circuitbreaker)              |
|                                                   |
|   - CLOSED: Loloskan ke Client                    |
|   - OPEN: Gagalkan Langsung (ErrCircuitOpen)      |
|   - HALF_OPEN: Izinkan Canary Probe               |
+---------------------------------------------------+
             │
             │ Lolos jika CLOSED / Uji Coba HALF_OPEN
             ▼
+-------------------------+
|      Payment Client     |
|   (internal/payment)    |
+-------------------------+
             │
             │ HTTP POST /pay (Network I/O)
             ▼
+---------------------------------------------------+
|             Fake Payment Server HTTP              |
|   (ModeHealthy: 200 OK | ModeDown: 500 Error |   |
|                 ModeSlow: Sleep)                  |
+---------------------------------------------------+
```

---

## 3. Perbandingan Alur: Tanpa vs Dengan Circuit Breaker

### Alur Tanpa Circuit Breaker (Menyebabkan Penumpukan Thread)
```text
Klien            Checkout Service             Payment Server (Lambat/Down)
  │                     │                                  │
  │── Request 1 ───────►│── HTTP POST (Menunggu Timeout) ─►│ (Hang 100ms)
  │                     │◄── Error Timeout (100ms) ────────│
  │◄── Response Gagal ──│                                  │
  │                     │                                  │
  │── Request 2 ───────►│── HTTP POST (Menunggu Timeout) ─►│ (Hang 100ms)
  │                     │◄── Error Timeout (100ms) ────────│
  │◄── Response Gagal ──│                                  │
  │                     │                                  │
  │── Request N ───────►│── HTTP POST (Thread Habis!) ────►│ (Hang 100ms)
  ▼                     ▼                                  ▼
(Dampak: Latensi terakumulasi sebesar N * timeout, koneksi habis, server macet)
```

### Alur Dengan Circuit Breaker (Fail-Fast dan Proteksi Beban)
```text
Klien            Checkout Service       Circuit Breaker         Payment Server
  │                     │                      │                       │
  │── Request 1 ───────►│── Execute() ────────►│── HTTP POST ─────────►│ (Gagal 500)
  │                     │                      │◄── Error ─────────────│
  │◄── Response Gagal ──│◄── Kembalikan Error ─│ (Gagal: 1/3)          │
  │                     │                      │                       │
  │── Request 2 ───────►│── Execute() ────────►│── HTTP POST ─────────►│ (Gagal 500)
  │                     │                      │◄── Error ─────────────│
  │◄── Response Gagal ──│◄── Kembalikan Error ─│ (Gagal: 2/3)          │
  │                     │                      │                       │
  │── Request 3 ───────►│── Execute() ────────►│── HTTP POST ─────────►│ (Gagal 500)
  │                     │                      │◄── Error ─────────────│
  │◄── Response Gagal ──│◄── Kembalikan Error ─│ [TRIP -> OPEN]        │
  │                     │                      │                       │
  │── Request 4 ───────►│── Execute() ────────►│ (Status OPEN)         │
  │                     │                      │ [TOLAK INSTAN: <1µs]  │ (Tidak ada I/O)
  │◄── ErrCircuitOpen ──│◄── ErrCircuitOpen ───│                       │
  │                     │                      │                       │
  ▼                     ▼                      ▼                       ▼
(Dampak: Menghentikan trafik ke downstream, menghemat thread pemanggil)
```

---

## 4. Alur Pemulihan Bertahap (HALF-OPEN Recovery)

```text
Checkout Service           Circuit Breaker                     Payment Server
       │                          │                                   │
       │                          │── Masa Cooldown Selesai (300ms)   │
       │                          │   (Status: HALF-OPEN)             │
       │                          │                                   │
       │── Request Baru ─────────►│── Izinkan 1 Canary Probe ────────►│
       │                          │                                   │
       │                          │◄── 200 OK (Respons Sukses) ───────│
       │                          │                                   │
       │                          │── Evaluasi: Sukses >= MaxCalls    │
       │                          │   (Status: CLOSED)                │
       │◄── Berhasil (Nil Error) ─│                                   │
       │                          │                                   │
       │── Request Berikutnya ───►│── Trafik Normal Diteruskan ──────►│
```
