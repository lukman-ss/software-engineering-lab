# Diagram 1 — Perbandingan Protocol Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        CLIENT                                  │
│  ┌──────────────┐              ┌──────────────┐                │
│  │ EventSource  │              │   WebSocket  │                │
│  │   (SSE)      │              │   (RFC 6455) │                │
│  └──────┬───────┘              └──────┬───────┘                │
│         │                             │                        │
└─────────┼─────────────────────────────┼────────────────────────┘
          │                             │
     HTTP GET                     HTTP GET (Upgrade)
  /sse (persistent)             /ws (101 Switching Protocols)
          │                             │
          ▼                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                        SERVER                                  │
│  ┌──────────────┐              ┌──────────────┐                │
│  │ SSE Hub      │              │  WS Upgrade  │                │
│  │ Broadcast    │              │  + Handler   │                │
│  │ Subscribe    │              │  Loop        │                │
│  └──────────────┘              └──────────────┘                │
└─────────────────────────────────────────────────────────────────┘

SSE:   Unidirectional (server → client), standard HTTP
WS:    Bidirectional (server ↔ client), custom framing over TCP
```

**Source:** Derived from `internal/sse/sse.go` (SSE Hub) dan `internal/ws/ws.go` (WS Upgrade/Handler) + `internal/server/server.go` (router)

---

# Diagram 2 — SSE Event Flow

```
Server                          Client
  │                               │
  │  Broadcast("news", "Go 1.22") │
  │  ─────── ID:1 ──────────────> │  (client belum connected)
  │  Broadcast("news", "Update")  │
  │  ─────── ID:2 ──────────────> │  (client belum connected)
  │                               │
  │                               │ GET /sse
  │                               │ Last-Event-ID: 1
  │ <──────────────────────────── │
  │                               │
  │  Subscribe(lastID=1)          │
  │  replay: [Event{ID:2}]       │
  │  ──────── id: 2 ────────────> │
  │  ──────── event: news ──────> │
  │  ──────── retry: 2000 ──────> │
  │  ──────── data: Update ─────> │
  │                               │
  │  ──────── : keep-alive ─────> │  (heartbeat setiap 15s)
  │                               │
  │  Broadcast("alert", "Bug!")   │
  │  ──────── id: 3 ────────────> │  (live event, langsung dikirim)
```

**Source:** Derived from `internal/sse/sse.go:98-151` (ServeHTTP) dan `internal/sse/sse.go:52-72` (Broadcast) + `internal/sse/sse.go:74-96` (Subscribe)

---

# Diagram 3 — WebSocket Frame Structure

```
RFC 6455 Frame:
┌─────────┬─────────┬───────────┬──────────┬──────────┐
│ FIN(1b) │Opcode(4)│  Mask(1)  │ Length   │Mask Key  │  Payload
│         │         │           │ (7b)     │ (0/4b)   │
├─────────┼─────────┼───────────┼──────────┼──────────┤
│  1      │ 0x1-0x2 │   0/1     │  7/16/64 │   0/4    │  N bytes
│  FIN=1  │Text/Bin │           │  bits    │          │
└─────────┴─────────┴───────────┴──────────┴──────────┘

Length encoding:
  ≤ 125    → 7-bit in header byte
  126-65535 → 126 + 2-byte extended (BigEndian)
  > 65535  → 127 + 8-byte extended (BigEndian)

XOR Unmasking:
  data[i] ^= mask[i % 4]    (RFC 6455 §5.3)
```

**Source:** `internal/ws/ws.go:84-127` (ReadFrame) dan `internal/ws/ws.go:129-174` (WriteFrame)

---

# Diagram 4 — SSE Hub State Model

```
┌─────────────────────────────────────────────────────┐
│                    SSE Hub                          │
│                                                     │
│  ┌──────────┐     ┌────────────────────────┐       │
│  │ history  │     │    clients map          │       │
│  │ []Event  │     │ map[chan Event]struct{} │       │
│  │          │     │                        │       │
│  │ [1]news  │     │  ch1 ───► client A     │       │
│  │ [2]update│     │  ch2 ───► client B     │       │
│  │ [3]alert │     │  ch3 ───► client C     │       │
│  │ ...      │     │  (capacity: 64 each)   │       │
│  └──────────┘     └────────────────────────┘       │
│                                                     │
│  Broadcast() → append history, non-blocking send    │
│  Subscribe(lastID) → replay missed + register ch   │
│  ServeHTTP() → flush replay, stream live           │
└─────────────────────────────────────────────────────┘
```

**Source:** `internal/sse/sse.go:38-96` (Hub struct, NewHub, Broadcast, Subscribe)

---

# Diagram 5 — Lab Endpoint Architecture

```
                    ┌──────────────┐
                    │  HTTP Server │
                    │ (ServeMux)   │
                    └──────┬───────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
              ▼            ▼            ▼
         ┌────────┐  ┌────────┐  ┌────────┐
         │  /sse  │  │  /ws   │  │/publish│
         │        │  │        │  │        │
         │ Hub    │  │ Upgrade│  │Broadcast│
         │.Serve  │  │ +Loop  │  │ Event  │
         │  HTTP()│  │        │  │ to Hub │
         └────────┘  └────────┘  └────────┘
              │            │
              ▼            ▼
         EventSource    WebSocket
         (browser)     (client)
```

**Source:** `internal/server/server.go:15-68` (NewServer)
