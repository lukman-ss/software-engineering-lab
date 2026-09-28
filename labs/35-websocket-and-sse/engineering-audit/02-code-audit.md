# Code Audit

Target Lab: labs/35-websocket-and-sse

## Finding 1

Location: `internal/sse/sse.go:19-36`
Claimed Behavior: Format WHATWG SSE text/event-stream events with `id`, `event`, `retry`, and multi-line `data`.
Observed Implementation: `Event.Format()` emits optional `id`, `event`, and `retry`, splits `Data` on newline and prepends `data: ` to each line, and terminates with `\n\n`.
Assessment: PASS
Severity: LOW
Notes: Correctly complies with W3C / WHATWG SSE specification format.

## Finding 2

Location: `internal/sse/sse.go:52-96`
Claimed Behavior: Safe concurrent SSE event broadcasting, subscription registration, and history replay via `Last-Event-ID`.
Observed Implementation: `Hub` protects `history` and `clients` map with `mu sync.RWMutex`. `Broadcast` acquires lock, appends to `history`, and sends to buffered client channels (`select default` prevents slow consumers from blocking hub). `Subscribe` replays events with `evt.ID > lastID` while lock is held.
Assessment: PASS
Severity: LOW
Notes: Concurrency safety verified with race detector. Slow consumer non-blocking drop logic handles slow readers gracefully.

## Finding 3

Location: `internal/sse/sse.go:98-151`
Claimed Behavior: HTTP streaming handler with flush and keep-alive heartbeats.
Observed Implementation: Sets headers `Content-Type: text/event-stream`, `Cache-Control: no-cache`, `Connection: keep-alive`. Replays history, streams live events from channel, and sends periodic comment `: keep-alive\n\n` on 15s ticker. Unsubscribes on disconnect (`ctx.Done()`).
Assessment: PASS
Severity: LOW
Notes: Standard flusher and context lifecycle cleanly implemented.

## Finding 4

Location: `internal/ws/ws.go:31-76`
Claimed Behavior: RFC 6455 HTTP upgrade handshake.
Observed Implementation: Validates `Upgrade: websocket`, computes SHA-1 of `Sec-WebSocket-Key` + RFC GUID `258EAFA5-E914-47DA-95CA-C5AB0DC85B11`, base64 encodes it to `Sec-WebSocket-Accept`, hijacks TCP connection and returns `HTTP/1.1 101 Switching Protocols`.
Assessment: PASS
Severity: LOW
Notes: Implements RFC 6455 section 4 handshake without third-party dependencies.

## Finding 5

Location: `internal/ws/ws.go:84-174`
Claimed Behavior: RFC 6455 framing with support for text (0x1), binary (0x2), close (0x8), ping/pong, and mask unmasking.
Observed Implementation: Correctly handles 7-bit, 16-bit (`length == 126`), and 64-bit (`length == 127`) extended payload lengths. Correctly reads and applies 4-byte XOR masking key on payload read. `WriteFrame` sets FIN bit (0x80) with opcode and formats payload lengths appropriately.
Assessment: PASS
Severity: LOW
Notes: Masking applied properly when writing client frames and verified during server reads.
