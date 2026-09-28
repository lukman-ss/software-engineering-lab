# Docs vs Code — labs/35-websocket-and-sse

## README Claims
1. "Demonstrates core differences between WebSocket (RFC 6455 full-duplex text & binary protocol) and Server-Sent Events (WHATWG SSE `text/event-stream` unidirectional push with automatic replay via `Last-Event-ID`)."
2. "Running demo shows SSE event resumption and WebSocket bidirectional message passing."
3. "`go test ./...` and `go test -race ./...` pass clean."
4. "`cmd/demo` runs cleanly and demonstrates both protocols."

## Code Reality
- SSE Hub (`internal/sse`) implements event formatting, broadcast, subscribe with replay via Last-Event-ID, and keep-alive heartbeat. ✅ matches claim 1.
- WS (`internal/ws`) implements RFC 6455 handshake and framing (text, binary, ping/pong, close). Echo server demonstrates bidirectional messaging. ✅ matches claim 1.
- Demo (`cmd/demo`) reproduces the described scenario: broadcasts two SSE events, reconnects with Last-Event-ID:1, reads resumed event 2, then opens raw TCP, upgrades to WS, sends masked text and binary frames, receives echoed text and reversed binary. Output matches `engineering/03-execution-result.md`. ✅ matches claim 2 & 4.
- Tests (`tests/protocol_test.go`) cover SSE formatting, Last-Event-ID replay, WS echo for text & binary, and SSE concurrency safety. They all PASS. ✅ matches claim 3.

## Mismatches
- README does not list `/publish` endpoint, yet code implements it (internal/server/server.go). No documentation, could mislead users. → DOC_CODE_MISMATCH.
- README claims "full-duplex" for WS; implementation echoes only; does not expose a generic broadcast channel. Still demonstrates full-duplex echo, acceptable.
- README mentions "automatic replay via `Last-Event-ID`" – implementation works, test confirms.

## Summary
Only discrepancy is undocumented `/publish` endpoint. No other DOC_CODE_MISMATCH, TEST_CLAIM_MISMATCH, or RESEARCH_IMPLEMENTATION_MISMATCH (research audit skipped).