# Test Audit

Target Lab: labs/35-websocket-and-sse

## Test Suite Overview

| Test Name | File | Targets | Type |
|---|---|---|---|
| `TestSSE_Formatting` | `tests/protocol_test.go:18` | `sse.Event.Format` | Unit / Spec Compliance |
| `TestSSE_LastEventID_Resumption` | `tests/protocol_test.go:31` | `sse.Hub.ServeHTTP` & Header Resumption | Integration |
| `TestWebSocket_TextAndBinary` | `tests/protocol_test.go:83` | `ws.Upgrade`, `ws.Dial`, Text & Binary Echo | Integration |
| `TestSSE_Concurrency` | `tests/protocol_test.go:136` | `sse.Hub` with 10 parallel clients & broadcasts | Concurrency & Race Check |

## Execution Results

### 1. Standard Test Run
Command: `go test -v -count=1 ./...`
Output:
```text
=== RUN   TestSSE_Formatting
--- PASS: TestSSE_Formatting (0.00s)
=== RUN   TestSSE_LastEventID_Resumption
--- PASS: TestSSE_LastEventID_Resumption (0.00s)
=== RUN   TestWebSocket_TextAndBinary
--- PASS: TestWebSocket_TextAndBinary (0.00s)
=== RUN   TestSSE_Concurrency
--- PASS: TestSSE_Concurrency (0.05s)
PASS
ok  	labs/35-websocket-and-sse/tests	0.413s
```

### 2. Race Detection Run
Command: `go test -race -v -count=1 ./...`
Output:
```text
=== RUN   TestSSE_Formatting
--- PASS: TestSSE_Formatting (0.00s)
=== RUN   TestSSE_LastEventID_Resumption
--- PASS: TestSSE_LastEventID_Resumption (0.00s)
=== RUN   TestWebSocket_TextAndBinary
--- PASS: TestWebSocket_TextAndBinary (0.00s)
=== RUN   TestSSE_Concurrency
--- PASS: TestSSE_Concurrency (0.05s)
PASS
ok  	labs/35-websocket-and-sse/tests	1.451s
```

### 3. Demo Execution
Command: `go run ./cmd/demo`
Output:
```text
=== DEMO: Server-Sent Events (SSE) vs WebSocket ===
Server listening on http://127.0.0.1:51021

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

## Coverage Evaluation

- Happy path: Covered (SSE stream reception, WS text/binary echo).
- Failure/Reconnection: Covered via `Last-Event-ID` resumption test.
- Concurrency: Covered via 10 concurrent subscribers receiving broadcasts under race detector.
- Edge Cases: Multi-line SSE formatting tested.
