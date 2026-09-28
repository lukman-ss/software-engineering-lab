# Test Audit

## Test Suite Overview

File: `tests/protocol_test.go`

Test cases reviewed:
1. `TestSSE_Formatting`: Verifies multiline data, custom event name, retry field, and id field rendering according to WHATWG event stream specs.
2. `TestSSE_LastEventID_Resumption`: Starts test HTTP server, emits 3 events, connects with `Last-Event-ID: 1`, confirms receipt of events 2 and 3 while verifying event 1 is skipped.
3. `TestWebSocket_TextAndBinary`: Verifies full RFC 6455 handshake over raw TCP, writes masked text frame and asserts echo response, writes masked binary frame and asserts byte-reversed response.
4. `TestSSE_Concurrency`: Spawns 10 concurrent HTTP client connections streaming SSE while broadcaster emits events, confirming thread safety and absence of deadlocks.

## Execution Verification

1. Standard Test Execution:
```bash
$ go test -v ./...
=== RUN   TestSSE_Formatting
--- PASS: TestSSE_Formatting (0.00s)
=== RUN   TestSSE_LastEventID_Resumption
--- PASS: TestSSE_LastEventID_Resumption (0.00s)
=== RUN   TestWebSocket_TextAndBinary
--- PASS: TestWebSocket_TextAndBinary (0.00s)
=== RUN   TestSSE_Concurrency
--- PASS: TestSSE_Concurrency (0.05s)
PASS
ok  	labs/35-websocket-and-sse/tests	0.276s
```

2. Race Detector Execution:
```bash
$ go test -race ./...
ok  	labs/35-websocket-and-sse/tests	1.258s
```
Result: PASS, zero data races detected.

3. Demo Execution:
```bash
$ go run ./cmd/demo
=== DEMO: Server-Sent Events (SSE) vs WebSocket ===
Server listening on http://127.0.0.1:51908

--- 1. SSE: Unidirectional Server Push with Last-Event-ID Resumption ---
Client connected with Last-Event-ID: 1. Received resumed events:
  [SSE Stream] id: 2
  [SSE Stream] event: news
  [SSE Stream] retry: 2000
  [SSE Stream] data: Update: Faster routing in net/http

--- 2. WebSocket: Full-Duplex Text and Binary Protocol ---
Client sending WS text frame: "Hello RFC 6455"
Server responded with Opcode 0x1 (Text): "echo: Hello RFC 6455"
Client sending WS binary frame: DEADBEEF
Server responded with Opcode 0x2 (Binary - Reversed): EFBEADDE
Binary frame payload verified successfully.

=== Demo Complete ===
```
Result: PASS, clean output matching protocol specifications.

## Coverage Assessment

- Happy Path: Fully covered (SSE formatting, resumption, WS text/binary echo).
- Failure Path: Connection hijacking check and basic handshake header validation covered.
- Concurrency: Covered via `TestSSE_Concurrency` and validated under Go race detector.
- Edge Cases: Multi-line SSE formatting covered; binary byte reversal verified.
