# Source Map

## Problem Statement

Research: `research/05-report.md` — Research Question
Implementation: `README.md` — Lab description
Engineering: `engineering/01-design.md`

---

## WebSocket vs SSE: Directionality

Research: `research/05-report.md` — Finding 1 (WebSocket full-duplex), Finding 4 (SSE unidirectional)
Sources: RFC 6455, MDN WebSocket, WHATWG HTML Standard
Implementation: `internal/ws/ws.go` (bidirectional framing), `internal/sse/sse.go` (server→client push)
Tests: `tests/protocol_test.go:83` — `TestWebSocket_TextAndBinary` (bidirectional), `tests/protocol_test.go:31` — `TestSSE_LastEventID_Resumption` (unidirectional)

---

## SSE Event Formatting

Research: `research/05-report.md` — Finding 5 (event identification), Finding 6 (UTF-8 text only)
Sources: WHATWG HTML Standard Section 9.2.1, MDN EventSource
Implementation: `internal/sse/sse.go:19-36` — `Event.Format()`
Tests: `tests/protocol_test.go:18` — `TestSSE_Formatting` (exact wire format)

---

## SSE Last-Event-ID Resumption

Research: `research/05-report.md` — Finding 5 (auto-reconnect, Last-Event-ID)
Sources: WHATWG HTML Standard Section 9.2.3, 9.2.4
Implementation: `internal/sse/sse.go:74-96` — `Subscribe()`, `internal/sse/sse.go:110-115` (header parsing)
Tests: `tests/protocol_test.go:31` — `TestSSE_LastEventID_Resumption`
Demo: `cmd/demo/main.go:28-51` — SSE resumption demo

---

## SSE Keep-alive Heartbeat

Research: `research/05-report.md` — Finding 10 (proxy timeout)
Sources: NGINX proxy module docs
Implementation: `internal/sse/sse.go:128-150` — ticker loop with `: keep-alive` comment
Tests: Verified indirectly via `TestSSE_Concurrency`

---

## WebSocket Handshake (RFC 6455)

Research: `research/05-report.md` — Finding 1 (HTTP/1.1 Upgrade), Finding 7 (HTTP/2 RFC 8441)
Sources: RFC 6455 Section 4, RFC 8441
Implementation: `internal/ws/ws.go:31-76` — `Upgrade()`
Tests: `tests/protocol_test.go:83` — `TestWebSocket_TextAndBinary` (implicit handshake via `Dial`)

---

## WebSocket Frame Encoding (RFC 6455 Framing)

Research: `research/05-report.md` — Finding 2 (binary + UTF-8 text payloads)
Sources: RFC 6455 Section 5
Implementation: `internal/ws/ws.go:84-127` — `ReadFrame()`, `internal/ws/ws.go:129-174` — `WriteFrame()`
Tests: `tests/protocol_test.go:83` — `TestWebSocket_TextAndBinary`
Demo: `cmd/demo/main.go:69-92` — text + binary frame demo

---

## WebSocket Text Echo

Research: `research/05-report.md` — Finding 2
Implementation: `internal/server/server.go:45-46` — `OpText` case
Tests: `tests/protocol_test.go:101-116` — text echo assertion

---

## WebSocket Binary Reversal

Research: `research/05-report.md` — Finding 2
Implementation: `internal/server/server.go:47-53` — `OpBinary` case
Tests: `tests/protocol_test.go:118-133` — reversed binary assertion
Demo: `cmd/demo/main.go:79-92` — binary reversal demo

---

## SSE Concurrent Broadcast

Research: `research/05-report.md` — Finding 11 (resource management)
Implementation: `internal/sse/sse.go:52-72` — `Broadcast()`, `internal/sse/sse.go:74-96` — `Subscribe()`
Tests: `tests/protocol_test.go:136` — `TestSSE_Concurrency`
Audit: `engineering-audit/02-code-audit.md` — Finding 11 (concurrency analysis)

---

## SSE History Buffer (Unbounded)

Research: `research/05-report.md` — Finding 11 (memory considerations)
Implementation: `internal/sse/sse.go:40,63` — `h.history = append(...)`
Audit: `engineering-audit/05-gaps.md` — Gap 4 (unbounded memory)
Warning: Production systems require ring buffer or TTL eviction

---

## WebSocket Static Mask Key (Warning)

Research: `research/05-report.md` — Finding 1 (RFC 6455 compliance)
Implementation: `internal/ws/ws.go:158` — static `{0x12, 0x34, 0x56, 0x78}`
Audit: `engineering-audit/02-code-audit.md` — Finding 7 (WARNING)
Audit: `engineering-audit/05-gaps.md` — Gap 1 (static masking key)

---

## Extended Payload Framing (Untested)

Implementation: `internal/ws/ws.go:94-106` (read), `internal/ws/ws.go:145-155` (write)
Audit: `engineering-audit/05-gaps.md` — Gap 2 (extended payload not exercised by tests)

---

## HTTP/2 Compatibility

Research: `research/05-report.md` — Finding 7 (WebSocket RFC 8441), Finding 8 (SSE native HTTP/2), Finding 9 (browser connection limits)
Sources: RFC 8441, RFC 9113 (obsoletes RFC 7540)
Note: HTTP/2 multiplexing not demonstrated in this lab

---

## 100k Concurrent Connections

Research: `research/05-report.md` — Finding 11
Sources: Linux epoll(7), Redis Pub/Sub, C10k/C10M principles
Note: Architectural principle, not empirically benchmarked in this lab

---

## Source Files

| Path | Purpose |
|------|---------|
| `internal/sse/sse.go` | SSE Hub, event formatting, Subscribe, Broadcast, ServeHTTP |
| `internal/ws/ws.go` | WebSocket Upgrade, ReadFrame, WriteFrame, Dial |
| `internal/server/server.go` | HTTP router with /sse, /ws, /publish endpoints |
| `cmd/demo/main.go` | Executable demo of SSE resumption and WebSocket bidirectional |
| `tests/protocol_test.go` | Integration tests: formatting, resumption, text/binary, concurrency |
| `go.mod` | Module definition, Go 1.22, zero external dependencies |
