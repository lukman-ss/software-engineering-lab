# Content Revision Record

Target Lab: labs/35-websocket-and-sse
Revision Date: 2026-09-28
Reviser: Technical Writer Reviser
Scope: Content only — 02-master-draft.md

---

## Changes Applied

### Fix 1 — Exponential Backoff (§3.1, Line 94)
- **Before:** "secara otomatis reconnect dengan exponential backoff"
- **After:** "secara otomatis reconnect dengan jeda sesuai field `retry` (implementation-defined, biasanya beberapa detik)"
- **Reason:** WHATWG does not mandate exponential backoff; reconnection time is implementation-defined per the `retry` field.

### Fix 2 — Demo Event-Count (§3.2, Line 124)
- **Before:** "event ID 2 dan 3 harus di-replay"
- **After:** "event ID 2 harus di-replay"
- **Reason:** Demo broadcasts only 2 events (IDs 1,2). "2 dan 3" conflated with the 3-event test scenario.

### Fix 3 — SSE Sticky-Session Oversimplification (§3.3, Line 32)
- **Before:** Ends at "Proses recovery terjadi tanpa intervensi developer."
- **After:** Added clause: "Untuk replay lintas node tetap perlu shared history/broker seperti Redis Pub/Sub — in-memory history bersifat per-node."
- **Reason:** In-memory history is per-node; cross-node `Last-Event-ID` replay fails without shared state.

### Fix 4 — WS Sticky-Session Absolute Claim (§3.4, Line 30)
- **Before:** "harus dikonfigurasi untuk mendukung WebSocket sticky sessions — koneksi harus tetap di node yang sama karena state koneksi ada di memori lokal."
- **After:** "umumnya membutuhkan konfigurasi WebSocket sticky sessions — koneksi harus tetap di node yang sama karena state koneksi ada di memori lokal — kecuali ada broker pesan terpusat."
- **Reason:** Sticky session avoidable with centralized message broker (Redis Pub/Sub per Finding 11).

### Fix 5 — WS Exempt from Browser Limits (§3.5, Line 104)
- **Before:** "WebSocket tidak terkena batasan ini karena menggunakan koneksi TCP terpisah yang tidak dihitung dalam batas HTTP."
- **After:** "WebSocket tidak dihitung dalam batas 6 koneksi HTTP/1.1 tersebut (tetapi browser tetap membatasi jumlah koneksi WS secara terpisah)."
- **Reason:** No research source asserts complete WS exemption; browsers impose separate WS connection limits.

---

## Files Modified

- `content/02-master-draft.md` — 5 edits

## Files Not Modified

- `content/03-code-snippets.md` — "ID 2 dan 3" in test description (line 380) is accurate (3-event test scenario, not demo).
- `content/01-content-brief.md` — No audit warnings target this file.
- `content/04-diagrams.md` — No audit warnings target this file.
- `content/05-key-takeaways.md` — No audit warnings target this file.
- `content/06-source-map.md` — No audit warnings target this file.

---

## Verification

All 5 audit warnings from `content-audit/01-content-audit.md` §3.1–§3.5 addressed. No regressions introduced.
