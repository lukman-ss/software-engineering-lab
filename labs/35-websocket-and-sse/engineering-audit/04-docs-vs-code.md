# Docs vs Code Comparison

Target Lab: labs/35-websocket-and-sse

## Comparison Matrix

| Component | Documentation Claim | Actual Code Implementation | Status |
|---|---|---|---|
| `internal/sse` | SSE Hub with custom event names, retry intervals, multi-line data framing, and buffer replay using `Last-Event-ID`. | Implemented in `internal/sse/sse.go`. Formats `id`, `event`, `retry`, `data`, and replays missing events when `Last-Event-ID` header is present. | MATCH |
| `internal/ws` | Pure Go RFC 6455 framing and connection upgrader supporting text (0x1) and binary (0x2) frames. | Implemented in `internal/ws/ws.go`. Performs HTTP 101 upgrade, SHA-1 key accept header validation, and masked text/binary frame codec. | MATCH |
| `internal/server` | Unified HTTP router hosting `/sse` and `/ws` endpoints. | Implemented in `internal/server/server.go`. Mounts `/sse`, `/ws`, and `/publish`. | MATCH |
| `cmd/demo` | Executable demonstration of SSE event resumption and WebSocket bidirectional message passing. | Implemented in `cmd/demo/main.go`. Executes live server and runs SSE stream resumption + WS text/binary echo. | MATCH |
| Test Commands | `go test ./...`, `go test -race ./...`, `go run ./cmd/demo` | All commands compile and execute cleanly without errors or race warnings. | MATCH |

## Discrepancy Findings

- DOC_CODE_MISMATCH: None detected.
- TEST_CLAIM_MISMATCH: None detected.
- RESEARCH_IMPLEMENTATION_MISMATCH: None detected.
