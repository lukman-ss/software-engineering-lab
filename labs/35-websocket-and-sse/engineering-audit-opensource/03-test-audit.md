# Test Audit — labs/35-websocket-and-sse

## Commands Executed (actual)
- `go build ./...` → PASS (exit 0)
- `go test ./...` → PASS (`ok labs/35-websocket-and-sse/tests (cached)`)
- `go test -count=1 -v ./tests/...` → PASS, 4/4:
  - TestSSE_Formatting PASS
  - TestSSE_LastEventID_Resumption PASS
  - TestWebSocket_TextAndBinary PASS
  - TestSSE_Concurrency PASS (0.05s)
- `go test -race ./tests/...` → PASS (`ok`, cached)
- `go vet ./...` → PASS (exit 0)
- `go run ./cmd/demo` → PASS, output matches engineering/03-execution-result.md (SSE id:2 replay + WS echo text/binary reversal verified)

## Coverage vs Required Matrix
- Happy path: COVERED (SSE format, Last-Event-ID replay of IDs 2,3 skipping 1; WS text echo + binary reversal)
- Failure path: NOT COVERED (no upgrade rejection, no invalid Last-Event-ID, no closed-conn read, no /publish test)
- Edge cases: PARTIAL (multi-line data formatting covered; extended-length WS frames, fragmented frames, ping/pong, close handshake not covered)
- Transitions: PARTIAL (SSE subscribe→replay→live covered; WS handshake→frame→echo covered; close/ping transitions untested)
- Recovery: WEAK (SSE resumption is the recovery path and is covered end-to-end; WS reconnect/recovery not implemented or tested — consistent with design, which promises only client-side reconnect conceptually)
- Rollback: NOT_APPLICABLE (no stateful rollback semantics in lab)
- Concurrency: COVERED but thin (10 SSE subscribers + 5 broadcasts under race detector; asserts no crash/hang only, not delivery completeness; WS concurrency untested)
- Negative cases: NOT COVERED (no malformed frame, missing upgrade header, bad key assertions)

## Assessment
Passing suite is genuine, not faked: replay test asserts absence of id:1 and presence of id:2/3; WS test asserts exact echo payloads and opcodes. But suite proves only the claimed core (SSE resume + WS text/binary echo) and would not catch handshake strictness, extended-length OOM, close/ping, or slow-subscriber drop gaps.
Strength: core claim coverage is real and race-clean. Weakness: failure/negative/edge coverage missing. Severity MEDIUM, non-blocking given lab scope.
