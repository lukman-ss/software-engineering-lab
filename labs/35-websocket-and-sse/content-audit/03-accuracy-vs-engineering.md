# Accuracy vs Engineering Implementation

Code checked: `internal/sse/sse.go`, `internal/ws/ws.go`, `internal/server/server.go`, `tests/protocol_test.go`, `cmd/demo/main.go`, `go.mod`.

## PASS — matches code

| Claim | Evidence |
|---|---|
| SSE Format id/event/retry/data, multi-line `data: `, `\n\n` | `sse.go:19-36`, `TestSSE_Formatting` exact string |
| Hub `Broadcast` monotonic ID, history append, non-blocking send, cap 64 | `sse.go:52-72`, `Subscribe` `make(chan Event, 64)` |
| `Subscribe` replay `evt.ID > lastID`, register under same lock | `sse.go:74-96` |
| Headers: `text/event-stream`, `no-cache`, `keep-alive`, `X-Accel-Buffering: no` | `sse.go:105-108` |
| Heartbeat `: keep-alive\n\n` every 15s | `sse.go:128-148` |
| Retry hardcoded 2000 | `sse.go:60`, demo output `retry: 2000` |
| WS Upgrade: Upgrade==websocket, SHA1(key+GUID), 101, Hijack | `ws.go:31-76` |
| GUID `258EAFA5-E914-47DA-95CA-C5AB0DC85B11` | `ws.go:22` |
| ReadFrame 7/16/64-bit length, XOR `mask[i%4]` | `ws.go:84-127` |
| WriteFrame FIN=1, static mask `{0x12,0x34,0x56,0x78}` | `ws.go:129-174` |
| `/sse`, `/ws`, `/publish` on ServeMux | `server.go:24-66` |
| OpClose / OpPing→Pong / OpText `"echo: "` / OpBinary reverse | `server.go:38-54` |
| 4 tests: Formatting, LastEventID, TextAndBinary, Concurrency | `protocol_test.go` |
| Demo Last-Event-ID:1 skips id:1, replays id:2 news Update… | `cmd/demo/main.go`, `engineering/03-execution-result.md` |
| Demo binary `DEADBEEF` → `EFBEADDE` | demo L80-91 |
| Go 1.22, zero deps | `go.mod` |
| Not full RFC 6455; no fragmented frames | `ws.go` switch ignores OpContinuation |
| Unbounded history | `history = append(...)` never trimmed |
| Extended length untested | all test payloads ≤125 bytes |

## WARN — incomplete / slightly off

1. **Keep-alive “verified indirectly via TestSSE_Concurrency”** (`06-source-map.md` L44). Test timeout 500ms; ticker 15s. Heartbeat path never runs. Race-detector only.
2. **Snippet 8 source `sse.go:128-150`**. Snippet starts at `ctx := r.Context()` = L127. Off-by-one.
3. **`Dial` lab shortcuts omitted from Warnings.** Engineering audit Finding 8: static RFC-example `Sec-WebSocket-Key`, single `Read` for 101. Draft mentions Dial exists; does not list these next to static mask / unbounded history.
4. **`Upgrade` does not check `Sec-WebSocket-Version: 13`.** Audit Finding 5. Content lists “not full spec” but not this concrete gap.
5. **Frame text omits RSV (3 bit).** First byte = FIN+RSV+opcode. Diagram 3 same omission. Code never sets/checks RSV. Pedagogical compression, not a false claim about this lab.

## Snippet fidelity

Snippets 1–7, 10: byte-accurate vs source (snippet 9/10 tests truncated with comments — acceptable).
Snippet 7 `/ws` handler: exact `server.go:26-56`.
