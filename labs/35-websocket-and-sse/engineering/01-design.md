# Engineering Design

Target Lab: labs/35-websocket-and-sse
Research Status: APPROVED

## Concept To Prove
Demonstrate the fundamental protocol differences between WebSocket (RFC 6455 full-duplex text/binary) and Server-Sent Events (WHATWG SSE text/event-stream unidirectional server push), highlighting native SSE auto-reconnect with `Last-Event-ID` versus custom client-side reconnect handling in WebSockets.

## Expected Behavior
- **WebSocket**: Bidirectional messaging (echo text + binary frame pong).
- **SSE**: Unidirectional server streaming with `id`, `event`, `data`, `retry` lines. Responds to `Last-Event-ID` header on re-connection to resume message stream from last received ID.

## Failure Scenario
- Stream interruption handled via client reconnect for both mechanisms.
- SSE automatically resumes stream from `Last-Event-ID`.
- WS detects connection drop and reconnects with application-level frame recovery logic.

## Success Criteria
- Native standard library HTTP implementation (using `nhooyr.io/websocket` or stdlib net/http for SSE + RFC 6455 framing / low-dependency WS handler).
- `go test ./...` and `go test -race ./...` pass clean.
- `cmd/demo` runs cleanly and demonstrates both protocols in action.

## Architecture
- `internal/sse`: HTTP handler for `text/event-stream` with buffer resume via `Last-Event-ID`.
- `internal/ws`: HTTP handler for WebSocket RFC 6455 upgrade/framing (text/binary echo & broadcast).
- `internal/server`: Combined HTTP server serving both endpoints.
- `tests/`: Integration tests verifying WebSocket duplex and SSE stream resumption.

## Components
- SSE Handler & Store
- WebSocket Handler & Protocol Framer
- Test Runner & CLI Demo

## Test Strategy
- Unit tests for SSE line parser/formatter.
- Integration tests for SSE resumption using `Last-Event-ID`.
- Integration tests for WS binary & text frame handling.
- Race detector verification.

## Execution Plan
1. Implement `internal/sse`.
2. Implement `internal/ws`.
3. Implement `internal/server` and `cmd/demo`.
4. Implement `tests/`.
5. Execute `go test -race ./...` and `go run ./cmd/demo`.
6. Record output in `engineering/03-execution-result.md`.

## Implementation Decisions
- Standard library Go `net/http` for SSE.
- Minimal zero-dependency WebSocket implementation for RFC 6455 handshake and framing (or standard library websocket helper).
