# Docs vs Code Audit

Target Lab: labs/35-websocket-and-sse

## Comparison Matrix

| Component / Claim | README Description | Code Implementation | Status |
|---|---|---|---|
| SSE Hub & Framing | `internal/sse`: custom event names, retry intervals, multi-line data framing | Implemented in `internal/sse/sse.go` with `Event.Format()` | MATCH |
| SSE Resumption | Replay buffer via `Last-Event-ID` | Implemented in `Hub.Subscribe()` using `Last-Event-ID` header | MATCH |
| RFC 6455 WebSocket | `internal/ws`: pure Go RFC 6455 framing & upgrader (text 0x1, binary 0x2) | Implemented in `internal/ws/ws.go` (Upgrade, ReadFrame, WriteFrame) | MATCH |
| Unified Server | `internal/server`: unified HTTP router hosting `/sse` and `/ws` | Implemented in `internal/server/server.go` | MATCH |
| Executable Demo | `cmd/demo`: demonstrates event resumption and WS bidirectional messaging | Implemented in `cmd/demo/main.go` | MATCH |
| Test Commands | `go test ./...`, `go test -race ./...`, `go run ./cmd/demo` | All commands execute cleanly with zero errors or race warnings | MATCH |

## Identified Mismatches
None. Documentation accurately describes the package structure, flags, commands, and protocol features without exaggeration.
