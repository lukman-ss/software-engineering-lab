# Test Audit

Target Lab: labs/35-websocket-and-sse

## Commands Executed
- `go build ./...` → PASS (no output)
- `go test ./... -v` → PASS (4/4: TestSSE_Formatting, TestSSE_LastEventID_Resumption, TestWebSocket_TextAndBinary, TestSSE_Concurrency)
- `go test -race ./... -v` → PASS (4/4, 1.245s, no race)
- `go run ./cmd/demo` → PASS (SSE resumption + WS text/binary echo, matches recorded output)

## Coverage Matrix
- happy path: PASS — SSE format, SSE resumption IDs 2,3, WS text echo `echo: ping test`, WS binary reverse
- failure path: MISSING — no test for bad Upgrade header, missing Sec-WebSocket-Key, hijack failure, malformed frame
- edge cases: PARTIAL — multi-line data covered; missing: invalid Last-Event-ID, large payload >125/>64k, ping/pong, close frame, continuation/fragmentation
- transitions: PASS — Subscribe→replay→live, WS text→binary sequence
- recovery: PARTIAL — SSE Last-Event-ID replay proven; WS reconnect not tested (design mentions app-level recovery, not implemented)
- rollback: NOT_APPLICABLE — no transactional state
- concurrency: WEAK PASS — 10 clients + 5 broadcasts, no panic, race clean; no delivery-count assertion
- negative cases: MISSING — no assertion on dropped broadcast, full channel, /publish validation

## Assessment
Suite proves core claims. Passing suite is real but thin on failure/negative paths. No fake tests observed.
