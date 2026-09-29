# Diagrams & Flowcharts

Dokumen ini memuat diagram arsitektur, diagram alur eksekusi request, dan perbandingan alur serangan pada `labs/36-cors-and-csrf`.

---

## 1. SOP vs CORS vs CSRF Boundary

Diagram berikut mengilustrasikan batasan Same-Origin Policy, peran CORS, dan celah yang dieksploitasi oleh CSRF di level browser:

```text
+-------------------------------------------------------------------------------+
|                                 BROWSER                                       |
|                                                                               |
|  Origin: https://attacker.com                                                 |
|  +-------------------------------------------------------------------------+  |
|  |  Cross-Origin WRITE (Form POST / Simple Request)                        |  |
|  |  ---------------------------------------------> Diizinkan oleh SOP      |  |
|  |                                                 (Memicu CSRF)           |  |
|  |                                                                         |  |
|  |  Cross-Origin READ (Fetch response payload)                             |  |
|  |  ---------------------------------------------> Dilarang oleh SOP       |  |
|  |                                                 (Kecuali diizinkan CORS)|  |
|  +-------------------------------------------------------------------------+  |
|                                     │                                         |
|                                     │ Request + Cookie Sesi Otomatis           |
|                                     ▼                                         |
+-------------------------------------------------------------------------------+
                                      │
                                      ▼
+-------------------------------------------------------------------------------+
|                                 BACKEND                                       |
|                                                                               |
|  [ CORS Middleware ]                                                          |
|  - Memeriksa header Origin                                                    |
|  - Menolak Preflight OPTIONS jika origin terlarang (HTTP 403)                  |
|  - MELEWATKAN simple POST ke handler aplikasi (tidak mencegah eksekusi!)       |
|                                     │                                         |
|                                     ▼                                         |
|  [ Banking Service Business Logic ]                                           |
|  - Mengeksekusi mutasi saldo di database ($1000 -> $700)                      |
|  - Mengembalikan HTTP 200 OK                                                  |
+-------------------------------------------------------------------------------+
```

---

## 2. Alur Serangan CSRF pada Endpoint Rentan (Vulnerable Flow)

Miskonsepsi: Pengembang mengira CORS menolak penyerang, padahal transaksi tetap tereksekusi di database.

```text
Victim Browser                 Attacker Site (evil.com)              Bank Backend
      │                                    │                              │
      │── (1) Buka evil.com ──────────────>│                              │
      │<── (2) Return Form Auto-Submit ────│                              │
      │                                                                   │
      │── (3) POST /api/transfer/vulnerable ─────────────────────────────>│
      │       Origin: https://evil.com                                    │
      │       Cookie: session_id=victim-secret                            │
      │       Body: to=acc-attacker&amount=400                            │
      │                                                                   │
      │                                                                   │── (4) CORS: Origin tidak ada di whitelist
      │                                                                   │       Tapi ini simple POST, teruskan ke next.
      │                                                                   │── (5) Cek Cookie valid: acc-victim
      │                                                                   │── (6) Mutasi: Kurangi saldo korban $400,
      │                                                                   │       Tambah saldo penyerang $400.
      │                                                                   │
      │<── (7) Response HTTP 200 OK (Tanpa header ACAO) ──────────────────│
      │                                                                   │
      │── [Browser memblokir JS attacker membaca response]                │
      │   (TAPI UANG SUDAH BERHASIL DITRANSFER!)                          │
```

---

## 3. Alur Pertahanan Anti-CSRF Token Berlapis (Protected Flow)

Pada endpoint terproteksi (`/api/transfer/protected`), pipeline middleware memvalidasi token kriptografis sebelum handler bisnis dipanggil:

```text
Legitimate Client                 Bank Server                  Attacker (evil.com)
       │                               │                                │
       │── (1) GET /api/csrf-token ───>│                                │
       │       Cookie: session_id      │                                │
       │<── (2) 200 OK (signed token) ─│                                │
       │                               │                                │
       │── (3) POST /api/transfer ────>│                                │
       │       Cookie: session_id      │                                │
       │       X-CSRF-Token: <token>   │                                │
       │                               │── (4) Validasi HMAC Signature  │
       │                               │── (5) Validasi Session ID Match│
       │                               │── (6) Validasi Non-Expired TTL │
       │                               │── (7) Eksekusi Transfer Aman   │
       │<── (8) 200 OK (Success) ──────│                                │
       │                               │                                │
       │                               │<── (9) Attacker POST (No token)│
       │                               │        Cookie: session_id      │
       │                               │                                │
       │                               │── (10) Middleware Tolak Token  │
       │                               │<── (11) 403 Forbidden ─────────│
       │                               │    (Saldo Tidak Berubah)       │
```

---

## 4. Pipeline Middleware Backend (`labs/36-cors-and-csrf`)

Urutan eksekusi middleware pada bank server:

```text
                       [ HTTP Request ]
                              │
                              ▼
               +─────────────────────────────+
               |       CORS Middleware       |
               |  (internal/cors/middleware) |
               +─────────────────────────────+
                 /                         \
    [ Preflight OPTIONS ]             [ Actual Request ]
      /               \                        │
[ Origin OK ]   [ Origin Bad ]                 │
     │                 │                       │
 204 No Content   403 Forbidden                ▼
                               +─────────────────────────────+
                               |   Fetch Metadata / Custom   |
                               |      Header Middleware      |
                               +─────────────────────────────+
                                 /                         \
                   [ Sec-Fetch-Site: cross-site ]    [ Header Valid / Same-Origin ]
                                 │                                 │
                           403 Forbidden                           ▼
                                                   +─────────────────────────────+
                                                   |       CSRF Middleware       |
                                                   |  (internal/csrf/middleware) |
                                                   +─────────────────────────────+
                                                     /                         \
                                      [ Invalid/Missing Token ]         [ Valid Token ]
                                                     │                         │
                                               403 Forbidden                   ▼
                                                               +─────────────────────────────+
                                                               |   Bank Business Handler     |
                                                               |     (internal/bank/app)     |
                                                               +─────────────────────────────+
                                                                               │
                                                                         200 OK Response
```
