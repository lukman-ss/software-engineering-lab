# Code Audit

Target Lab: labs/35-websocket-and-sse

## Finding 1

Location: `internal/sse/sse.go:19-36`
Claimed Behavior: SSE formatting conforms to WHATWG event stream syntax supporting id, event, retry, and multi-line data terminated by double newlines.
Observed Implementation: `Event.Format()` formats `id`, `event`, `retry`, splits `Data` by newline for `data: ` prefixes, and terminates with `\n\n`.
Assessment: PASS
Severity: LOW
Notes: Correct standard WHATWG format.

## Finding 2

Location: `internal/sse/sse.go:38-96`
Claimed Behavior: Thread-safe SSE Hub broadcasting and subscription management with `Last-Event-ID` message replay.
Observed Implementation: Hub uses `sync.RWMutex`. `Broadcast` locks and appends to `history`, non-blocking send to client channels via `select default`. `Subscribe` filters `history` where `evt.ID > lastID` and returns buffered channel and unsubscribe callback. Unsubscribe removes client and closes channel under lock.
Assessment: PASS
Severity: LOW
Notes: Non-blocking client drop prevents slow consumer stalls. Concurrency is race-free.

## Finding 3

Location: `internal/sse/sse.go:110-150`
Claimed Behavior: HTTP handler stream delivery, context cancellation teardown, and periodic keep-alive heartbeat.
Observed Implementation: Verifies `http.Flusher`. Parses `Last-Event-ID` header. Flushes historical replay before entering select loop. Ticker emits SSE comment `: keep-alive\n\n` every 15s. Cleans up subscription via deferred unsubscribe on client disconnect or handler exit.
Assessment: PASS
Severity: LOW
Notes: Resilient HTTP streaming pattern.

## Finding 4

Location: `internal/ws/ws.go:31-76`
Claimed Behavior: RFC 6455 WebSocket Upgrade handshake using SHA-1 + standard GUID base64 hash.
Observed Implementation: Verifies Upgrade header and Sec-WebSocket-Key. Computes SHA1(key + wsGUID) and base64 encodes it for `Sec-WebSocket-Accept`. Hijacks HTTP connection and writes HTTP/1.1 101 Switching Protocols.
Assessment: PASS
Severity: LOW
Notes: Standard compliance verified.

## Finding 5

Location: `internal/ws/ws.go:84-174`
Claimed Behavior: RFC 6455 frame reading and writing with text (0x1), binary (0x2), close (0x8), ping (0x9), pong (0xA), and masking/unmasking.
Observed Implementation: Decodes frame headers, 7-bit, 16-bit (126), and 64-bit (127) payload lengths, XOR unmasking with 4-byte key. `WriteFrame` serializes frames with mutex synchronization.
Assessment: PASS
Severity: LOW
Notes: Fragmented frames (FIN=0) are not handled, which is documented as an intentional minimal lab scope limitation in `engineering/02-implementation-notes.md`.

## Finding 6

Location: `internal/server/server.go:26-56`
Claimed Behavior: Unified server handling `/sse`, `/ws`, and `/publish` routes.
Observed Implementation: Handles text echo, binary reverse payload echo, ping/pong frames, and clean close opcode handling.
Assessment: PASS
Severity: LOW
Notes: Correct handling across protocol multiplexing.
