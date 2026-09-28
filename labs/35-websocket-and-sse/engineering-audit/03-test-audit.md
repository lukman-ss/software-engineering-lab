# Test Audit

## Test Suite Overview

| Test Name | File | Purpose | Assertion Quality | Execution Result |
|---|---|---|---|---|
| `TestSSE_Formatting` | `tests/protocol_test.go:18` | Verifies multi-line data, id, event, retry format strings | Exact string comparison against WHATWG wire format | PASS |
| `TestSSE_LastEventID_Resumption` | `tests/protocol_test.go:31` | Verifies events prior to Last-Event-ID are skipped, subsequent events are replayed | Checks omission of event 1, inclusion of events 2 and 3 via HTTP stream | PASS |
| `TestWebSocket_TextAndBinary` | `tests/protocol_test.go:83` | Tests text echo and reversed binary payload delivery over RFC 6455 framing | Verifies text echo prefix + exact binary byte reversal | PASS |
| `TestSSE_Concurrency` | `tests/protocol_test.go:136` | Concurrently subscribes 10 clients and broadcasts 5 events | Race detector clean, all goroutines complete | PASS |

---

## Execution Logs

### Standard Test Execution

```text
$ go test -v -count=1 ./...
?       labs/35-websocket-and-sse/cmd/demo      [no test files]
?       labs/35-websocket-and-sse/internal/server       [no test files]
?       labs/35-websocket-and-sse/internal/sse  [no test files]
?       labs/35-websocket-and-sse/internal/ws   [no test files]
=== RUN   TestSSE_Formatting
--- PASS: TestSSE_Formatting (0.00s)
=== RUN   TestSSE_LastEventID_Resumption
--- PASS: TestSSE_LastEventID_Resumption (0.00s)
=== RUN   TestWebSocket_TextAndBinary
--- PASS: TestWebSocket_TextAndBinary (0.00s)
=== RUN   TestSSE_Concurrency
--- PASS: TestSSE_Concurrency (0.05s)
PASS
ok      labs/35-websocket-and-sse/tests 0.211s
```

### Race Detector Execution

```text
$ go test -race -v -count=1 ./...
=== RUN   TestSSE_Formatting
--- PASS: TestSSE_Formatting (0.00s)
=== RUN   TestSSE_LastEventID_Resumption
--- PASS: TestSSE_LastEventID_Resumption (0.00s)
=== RUN   TestWebSocket_TextAndBinary
--- PASS: TestWebSocket_TextAndBinary (0.00s)
=== RUN   TestSSE_Concurrency
--- PASS: TestSSE_Concurrency (0.05s)
PASS
ok      labs/35-websocket-and-sse/tests 1.209s
```

### Demo Execution

```text
$ go run ./cmd/demo
=== DEMO: Server-Sent Events (SSE) vs WebSocket ===
Server listening on http://127.0.0.1:52805

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

---

## Test Coverage Evaluation

- Happy Path: Fully covered for both SSE and WS.
- Failure / Edge Paths:
  - SSE empty ID handling: Tested indirectly in `TestSSE_Formatting` (ID 42).
  - SSE multi-line data: Tested in `TestSSE_Formatting`.
  - WS text and binary frames: Tested in `TestWebSocket_TextAndBinary`.
  - WS Ping/Pong/Close control frames: Uncovered in tests (implemented in `internal/server/server.go:39-44`, but no dedicated unit test asserts Ping/Pong or Close handshake).
  - WS payload length > 125 (16-bit and 64-bit lengths): Uncovered in tests (implementation exists in `ws.go:94-106` and `ws.go:145-155`, but tests only send small payloads < 125 bytes).
  - Missing coverage: No negative test for invalid `Upgrade` header or missing `Sec-WebSocket-Key`.
