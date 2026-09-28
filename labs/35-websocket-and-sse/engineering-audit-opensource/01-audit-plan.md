# Engineering Audit Plan

Target Lab: labs/35-websocket-and-sse
Implementation Files:
- internal/sse/sse.go (151 lines: Event.Format, Hub Broadcast/Subscribe/ServeHTTP)
- internal/ws/ws.go (202 lines: Upgrade, ReadFrame, WriteFrame, Dial, Conn)
- internal/server/server.go (69 lines: /sse, /ws, /publish router)
- cmd/demo/main.go (95 lines: SSE resumption + WS text/binary demo)
Tests:
- tests/protocol_test.go (166 lines, 4 tests: SSE_Formatting, SSE_LastEventID_Resumption, WebSocket_TextAndBinary, SSE_Concurrency)
Executable/Demo:
- cmd/demo (httptest server + SSE Last-Event-ID:1 replay + WS masked text/binary echo)
Approved Research Inputs:
- NOT AUDITED in this stage per PIPELINE OVERRIDE (research/ + research-audit/ skipped; research-audit/07-verdict.md = APPROVED noted for context only)
Main Claims To Verify:
1. SSE Hub emits spec-shaped frames (id/event/retry/data + multi-line data) and replays history via Last-Event-ID
2. WS implements RFC 6455 handshake + text(0x1)/binary(0x2) framing with masking, echo text with prefix, reverse binary
3. Server hosts /sse + /ws on one router; demo output is real
4. Tests prove claimed behavior including concurrency safety
Commands To Run:
- go build ./...
- go test ./... (+ -count=1 -v ./tests/...)
- go test -race ./tests/...
- go vet ./...
- go run ./cmd/demo
Primary Risks:
- Test suite is happy-path only (no handshake-failure, ping/pong, close-frame, extended-length, invalid-Last-Event-ID tests)
- ws.Dial ignores its urlStr argument (request path hardcoded); works only because server path is always /ws
- SSE history unbounded; broadcast drops to slow subscribers silently
