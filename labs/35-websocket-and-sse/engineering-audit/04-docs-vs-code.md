# Docs vs Code Audit

## Verification Matrix

| Claim Source | Claim | Code Implementation | Status |
|---|---|---|---|
| README.md:3 | WebSocket (RFC 6455 full-duplex text & binary protocol) | `internal/ws/ws.go`: Supports text (0x1) & binary (0x2) frames | PASS |
| README.md:3 | WHATWG SSE `text/event-stream` unidirectional push with replay via `Last-Event-ID` | `internal/sse/sse.go`: Header `text/event-stream`, replay via `Last-Event-ID` | PASS |
| README.md:7 | SSE Hub with custom event names, retry intervals, multi-line data, buffer replay | `internal/sse/sse.go`: Fields `Event`, `Retry`, `ID`, multi-line `Data` splitting | PASS |
| README.md:8 | Pure Go RFC 6455 framing and connection upgrader | `internal/ws/ws.go`: Implements handshake, masking, frame decode/encode | PASS |
| README.md:9 | Unified HTTP router hosting `/sse` and `/ws` endpoints | `internal/server/server.go`: Wires `/sse` and `/ws` on `http.ServeMux` | PASS |
| README.md:10 | Demo of SSE event resumption and WS message passing | `cmd/demo/main.go`: Demonstrates both SSE resumption and WS text/binary echo | PASS |
| README.md:11 | End-to-end integration tests and concurrency tests | `tests/protocol_test.go`: 4 tests covering formatting, resumption, text/binary, concurrency | PASS |

---

## Detailed Findings

1. `DOC_CODE_MISMATCH`: None identified. All claims in README accurately describe the codebase.
2. `TEST_CLAIM_MISMATCH`: None identified.
3. `RESEARCH_IMPLEMENTATION_MISMATCH`: None identified. Research report claims standard SSE text stream and RFC 6455 binary frame decoding; implementation delivers both pure Go without external dependencies.
