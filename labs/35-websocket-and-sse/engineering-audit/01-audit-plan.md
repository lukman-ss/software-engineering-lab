# Engineering Audit Plan

Target Lab: labs/35-websocket-and-sse
Implementation Files:
- `internal/sse/sse.go` (SSE event model, formatting, hub broadcast/subscription, HTTP stream handler)
- `internal/ws/ws.go` (RFC 6455 frame encoding/decoding, handshake upgrade, dialer)
- `internal/server/server.go` (Unified HTTP router wiring `/sse`, `/ws`, and `/publish`)
Tests:
- `tests/protocol_test.go` (`TestSSE_Formatting`, `TestSSE_LastEventID_Resumption`, `TestWebSocket_TextAndBinary`, `TestSSE_Concurrency`)
Executable/Demo:
- `cmd/demo/main.go`
Approved Research Inputs:
- `research/05-report.md`
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`
- `engineering/03-execution-result.md`

Main Claims To Verify:
1. WHATWG SSE spec conformance: framing format (`id`, `event`, `retry`, `data\n\n`), multi-line data framing, header handling (`text/event-stream`, `Cache-Control: no-cache`, `Connection: keep-alive`), and replay from `Last-Event-ID`.
2. RFC 6455 WebSocket conformance: handshake upgrade (`Upgrade: websocket`, `Connection: Upgrade`, `Sec-WebSocket-Accept` SHA1+GUID hash computation), frame binary parsing/serialization (FIN bit, opcode, masking bit, payload length 7-bit/16-bit/64-bit, unmasking/masking with 4-byte key).
3. Text (0x1) and Binary (0x2) payload handling across client and server.
4. Concurrency safety: Hub client map synchronization, channel lifecycle, WebSocket connection mutex safety under simultaneous read/write operations.
5. No artificial delays or fabricated demo results.

Commands To Run:
- `go test -v -count=1 ./...`
- `go test -race -v -count=1 ./...`
- `go run ./cmd/demo`

Primary Risks:
- WebSocket connection half-duplex locking: `c.mu` in `ws.go` protects both `Close` and `WriteFrame`, but `ReadFrame` does not lock `c.mu`. Need to verify if concurrent read and write on `net.Conn` / `bufio.ReadWriter` is thread-safe or prone to data races.
- Channel closing and broadcasting race in SSE Hub (`unsubscribe` vs `Broadcast`).
- Unbounded history growth in SSE Hub memory.
- WebSocket frame fragmentation (FIN=0) and control frame handling (Ping/Pong/Close) completeness.
- Lack of error return / client ping-pong heartbeats in real-world scenarios.
