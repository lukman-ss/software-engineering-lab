# Docs vs Code

Target Lab: labs/35-websocket-and-sse

## Compared
README.md vs server.go/sse.go/ws.go/tests/demo; engineering/01-design.md, 02-implementation-notes.md, 03-execution-result.md vs code.

## Findings
1. DOC_CODE_MISMATCH (LOW): `/publish` endpoint exists in server.go:58-66, undocumented in README.
2. RESEARCH_IMPLEMENTATION_MISMATCH (MEDIUM): design claims "WS detects drop and reconnects with app-level frame recovery" — no reconnect logic in ws.go/server.go.
3. RESEARCH_IMPLEMENTATION_MISMATCH (MEDIUM): design architecture claims WS "echo & broadcast" — code echo-only, no broadcast.
4. MATCH: SSE id/event/data/retry + Last-Event-ID replay, heartbeat, text 0x1/binary 0x2, zero-dep stdlib, race clean — all verified in code/tests/demo.
5. MATCH: execution-result demo output reproduces exactly on re-run (port differs only).
6. MATCH: known limitations (no fragmentation/extensions, unbounded SSE history) accurately disclosed in 02-implementation-notes.

## Result
Documentation Accuracy: WARNING (undoc endpoint + two design overclaims, no fake output)
