# Docs vs Code Audit

## Comparison Matrix

| Component / Claim | Documentation (README / Engineering Notes) | Implementation (`internal/`, `cmd/`) | Tests (`tests/`) | Assessment |
| --- | --- | --- | --- | --- |
| SSE Framing | WHATWG SSE `text/event-stream` with `id`, `event`, `retry`, `data` | Implemented in `internal/sse/sse.go:Event.Format()` | Tested in `TestSSE_Formatting` | PASS |
| SSE Resumption | Replay buffer via `Last-Event-ID` header | Implemented in `internal/sse/sse.go:Hub.Subscribe` | Tested in `TestSSE_LastEventID_Resumption` | PASS |
| SSE Keep-alive | Keep-alive comment heartbeats | Implemented in `internal/sse/sse.go:ServeHTTP` via 15s ticker | Verified in code | PASS |
| RFC 6455 Handshake | HTTP 101 Switching Protocols with Sec-WebSocket-Accept calculation | Implemented in `internal/ws/ws.go:Upgrade` & `Dial` | Tested in `TestWebSocket_TextAndBinary` | PASS |
| WS Text & Binary Frames | Opcode 0x1 (Text) and Opcode 0x2 (Binary) framing with masking | Implemented in `internal/ws/ws.go:ReadFrame` & `WriteFrame` | Tested in `TestWebSocket_TextAndBinary` | PASS |
| Demo Execution | `go run ./cmd/demo` outputs both SSE and WS runs | Implemented in `cmd/demo/main.go` | Executed and verified | PASS |

## Discrepancies Found

None. All claimed structures, entry points, and protocol features match between README.md, design documentation, source code, and integration tests.
