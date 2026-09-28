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
- research-revision/03-revision-result.md

Main Claims To Verify:
1. WHATWG SSE unidirectional streaming with custom events, multi-line data, retry headers, keep-alive comments, and Last-Event-ID replay.
2. RFC 6455 WebSocket handshake (HTTP 101, Sec-WebSocket-Accept calculation), text frame (0x1), binary frame (0x2), Ping (0x9), Pong (0xA), and Close (0x8) handling with payload masking/unmasking.
3. Clean execution of tests (`go test ./...` and `go test -race ./...`) and demo (`go run ./cmd/demo`).
4. Correctness of concurrency controls in SSE Hub and WS Conn.

Commands To Run:
- go test -v -count=1 ./...
- go test -race -v -count=1 ./...
- go run ./cmd/demo

Primary Risks:
- Race conditions during SSE broadcast / subscribe / unsubscribe under high concurrency.
- Incomplete WebSocket RFC 6455 framing (e.g. unhandled fragmentation bit, mask key fixed values, partial reads).
- Unbounded memory accumulation in SSE history buffer.
