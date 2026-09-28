# Engineering Audit Plan

Target Lab: labs/35-websocket-and-sse
Implementation Files:
- `internal/sse/sse.go`
- `internal/ws/ws.go`
- `internal/server/server.go`
Tests:
- `tests/protocol_test.go`
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research/05-report.md`
- `research-audit/07-verdict.md`
Main Claims To Verify:
- SSE unidirectional server push with `text/event-stream` framing (`id`, `event`, `retry`, multi-line `data`).
- SSE buffer replay and automatic resumption via `Last-Event-ID` header.
- Pure Go RFC 6455 WebSocket handshake (`Sec-WebSocket-Accept` computation using GUID `258EAFA5-E914-47DA-95CA-C5AB0DC85B11`).
- WebSocket binary (0x2) and text (0x1) framing and masking/unmasking.
- Thread-safe SSE Hub broadcasting and connection lifecycle.
Commands To Run:
- `go test ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
Primary Risks:
- Data race conditions in `sse.Hub` client map or buffer history operations.
- Incorrect WebSocket frame masking/unmasking or binary length parsing (16-bit / 64-bit extended length fields).
- Deadlocks or unbuffered write blocks on SSE broadcast channels under heavy client load.
