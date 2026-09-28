# Engineering Audit Plan

Target Lab: labs/35-websocket-and-sse
Implementation Files:
- internal/sse/sse.go
- internal/ws/ws.go
- internal/server/server.go
- cmd/demo/main.go

Tests:
- tests/protocol_test.go

Executable/Demo:
- cmd/demo/main.go

Approved Research Inputs:
- research/05-report.md
- engineering/01-design.md
- engineering/02-implementation-notes.md

Main Claims To Verify:
1. WHATWG SSE implementation supporting custom event names, retry intervals, multi-line data framing, Last-Event-ID buffer replay, and client heartbeats.
2. RFC 6455 WebSocket implementation supporting HTTP Upgrade handshake with Sec-WebSocket-Accept generation, text (0x1) and binary (0x2) frame read/write with masking.
3. Unified server running both `/sse` and `/ws` endpoints.
4. Clean concurrency without data races under `go test -race ./...`.

Commands To Run:
- `go test ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- SSE channel buffer overflow / blocking subscribers when slow clients do not consume events.
- Deadlocks or unhandled goroutine leaks in SSE subscriber cleanup or keep-alive loop.
- Partial frame reading / unbounded payload allocation in WebSocket `ReadFrame`.
- Concurrency race conditions on WebSocket frame write/read if accessed across multiple goroutines.
