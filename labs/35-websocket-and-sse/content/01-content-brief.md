# Content Brief

**Topic:** Membandingkan WebSocket (full-duplex) dan Server-Sent Events (SSE, unidirectional) untuk real-time server push

**Target Reader:** Backend developer dan system architect yang membangun sistem real-time dan perlu memilih protokol transport

**Problem:** Saat membangun fitur real-time (notification, live dashboard, LLM token streaming), developer sering tidak yakin kapan harus pakai WebSocket atau SSE. Masing-masing punya trade-off berbeda dalam hal kompleksitas, kompatibilitas proxy, dan skabilitas.

**Core Mental Model:** Keputusan utama bergantung pada **directionality** — apakah komunikasi hanya server→client (SSE) atau bidirectional (WebSocket)

**Approved Research Status:** APPROVED (11 claims verified, 0 contradictions, 0 unsupported claims)

**Approved Engineering Status:** APPROVED (4 tests PASS, race detector PASS, demo PASS)

**Main Concepts:**
- WebSocket RFC 6455: full-duplex, non-HTTP protocol over TCP
- SSE WHATWG: unidirectional server→client, standard HTTP response
- `text/event-stream` MIME type
- RFC 8441 WebSocket over HTTP/2
- HTTP/2 stream multiplexing
- `Last-Event-ID` resumption mechanism

**Verified Behaviors:**
- SSE event formatting: `id`, `event`, `data`, `retry` fields
- SSE `Last-Event-ID` resumption: skip event di bawah ID, replay event setelahnya
- WebSocket text frame echo
- WebSocket binary frame reversal
- Concurrent broadcast ke multiple SSE clients
- WebSocket handshake: `101 Switching Protocols`

**Available Case Studies:**
- Lab demo: SSE resumption dari `Last-Event-ID: 1` melewati event pertama
- Lab demo: WebSocket binary payload reversal terverifikasi byte-level
- 10 concurrent SSE clients tanpa race condition

**Warnings:**
- Lab ini bukan RFC 6455 full spec implementation
- Static mask key untuk WebSocket masking (bukan crypto-random)
- Unbounded history buffer di SSE Hub (hanya untuk lab)
- Extended payload framing (length >125) ada di code tapi belum tested
