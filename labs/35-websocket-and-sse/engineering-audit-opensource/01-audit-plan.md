# Engineering Audit Plan

Target Lab: labs/35-websocket-and-sse
Implementation Files:
- internal/sse/sse.go
- internal/ws/ws.go
- internal/server/server.go
- cmd/demo/main.go
Tests: tests/protocol_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: engineering/01-design.md, engineering/02-implementation-notes.md, engineering/03-execution-result.md
Main Claims To Verify:
- SSE unidirectional server push with id, event, data, retry lines and Last-Event-ID resumption
- WebSocket full-duplex binary and text frame handling per RFC 6455
- Thread-safe concurrency for both protocols
- Demo shows both protocols in action
Commands To Run:
- go build ./...
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Memory leak in SSE history buffer (unbounded slice)
- WS framing may not handle all RFC 6455 edge cases (fragmented frames, extensions)
- No authentication or authorization in demo