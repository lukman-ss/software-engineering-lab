# Lab 35: WebSocket and Server-Sent Events (SSE)

Demonstrates the core differences between WebSocket (RFC 6455 full-duplex text & binary protocol) and Server-Sent Events (WHATWG SSE `text/event-stream` unidirectional push with automatic replay via `Last-Event-ID`).

## Structure

- `internal/sse`: SSE Hub with support for custom event names, retry intervals, multi-line data framing, and buffer replay using `Last-Event-ID`.
- `internal/ws`: Pure Go RFC 6455 framing and connection upgrader supporting text (0x1) and binary (0x2) frames.
- `internal/server`: Unified HTTP router hosting `/sse` and `/ws` endpoints.
- `cmd/demo`: Executable demonstration of SSE event resumption and WebSocket bidirectional message passing.
- `tests/`: End-to-end integration tests and concurrency tests.

## Running Tests and Demo

```bash
# Run tests
go test ./...

# Run tests with race detection
go test -race ./...

# Run the demo
go run ./cmd/demo
```
