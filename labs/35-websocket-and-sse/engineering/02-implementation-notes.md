# Implementation Notes

Target Lab: labs/35-websocket-and-sse

## Files Added
- `go.mod`
- `README.md`
- `internal/sse/sse.go`
- `internal/ws/ws.go`
- `internal/server/server.go`
- `cmd/demo/main.go`
- `tests/protocol_test.go`
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`
- `engineering/03-execution-result.md`

## Core Design Decisions
- Pure Go standard library: Built RFC 6455 handshake, frame masking/unmasking, and framing directly via `net/http` connection hijacking to eliminate external dependencies.
- SSE Hub manages in-memory event buffers to demonstrate `Last-Event-ID` resumption seamlessly.
- Thread-safe channels and mutex protections across SSE subscriber registries and WebSocket connection writers.

## Implementation-Specific Choices
- Replay buffer size is unconstrained in-memory slice for simplicity of demonstration.
- WebSocket binary test handles reversing bytes to explicitly prove byte-level payload integrity over the wire.
- Fixed 4-byte client mask key for WS frame writing demo.

## Known Limitations
- Not a full RFC 6455 spec implementation (does not handle fragmented frames or complex extensions).
- SSE buffer is bounded only by process memory; production systems require ring buffers or durable log backing.

## Trade-offs
- Standard library minimal implementation vs full third-party library (`nhooyr.io/websocket` or `gorilla/websocket`): Avoided third-party packages to guarantee zero-dependency reproducibility and direct insight into raw byte framing.

## What Is Demonstrated
- Unidirectional SSE push with `id`, `event`, `data`, `retry`, and stream resume from `Last-Event-ID`.
- Full-duplex WebSocket frames supporting both UTF-8 text and binary payloads.
- Keep-alive heartbeat comments in SSE (`: keep-alive`).

## What Is Not Demonstrated
- HTTP/2 multiplexing bootstrap (RFC 8441).
- 100,000 concurrent connection load testing (deferred to dedicated infra benchmark).
