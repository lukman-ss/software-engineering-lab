# Code Audit

## Finding 1

Location: `internal/sse/sse.go:65-70`
Claimed Behavior: Non-blocking broadcast to all active subscribers.
Observed Implementation: `select { case ch <- evt: default: }` ensures slow subscribers do not block broadcast loop. Buffer size is 64.
Assessment: PASS
Severity: LOW
Notes: Correctly avoids head-of-line blocking on Hub broadcaster.

## Finding 2

Location: `internal/sse/sse.go:74-96`
Claimed Behavior: Buffer replay on reconnection with Last-Event-ID.
Observed Implementation: Queries `h.history` under mutex for `evt.ID > lastID` and returns replay slice. Clean `unsubscribe` closure closes client channel and removes from map.
Assessment: PASS
Severity: LOW
Notes: Replay slice sent prior to streaming dynamic events.

## Finding 3

Location: `internal/sse/sse.go:98-151`
Claimed Behavior: Serve SSE stream with headers, keep-alive heartbeat, and cancellation detection.
Observed Implementation: Sets `text/event-stream`, `Cache-Control: no-cache`, `Connection: keep-alive`, `X-Accel-Buffering: no`. Detects `r.Context().Done()` and unregisters subscriber cleanly.
Assessment: PASS
Severity: LOW
Notes: Keep-alive ticker fires every 15s with SSE comment format `: keep-alive\n\n`.

## Finding 4

Location: `internal/ws/ws.go:31-76`
Claimed Behavior: RFC 6455 Handshake upgrading HTTP connection.
Observed Implementation: Validates headers, computes SHA-1 with standard GUID `258EAFA5-E914-47DA-95CA-C5AB0DC85B11`, encodes base64, hijacks connection, writes 101 Switching Protocols response, and flushes buffer.
Assessment: PASS
Severity: LOW
Notes: Compliant with RFC 6455 Section 4.2.2.

## Finding 5

Location: `internal/ws/ws.go:84-127`
Claimed Behavior: Frame parsing for standard opcodes, extended payload lengths, and client masking.
Observed Implementation: Correctly decodes FIN/opcode, parses 7-bit, 16-bit (126), and 64-bit (127) lengths using big-endian ordering, unmasks payload using 4-byte key if masked bit set.
Assessment: PASS
Severity: LOW
Notes: Uses `io.ReadFull` ensuring complete reads on network buffers.

## Finding 6

Location: `internal/ws/ws.go:129-174`
Claimed Behavior: Safe frame writing with optional client masking.
Observed Implementation: Guarded with `c.mu.Lock()`. Encodes lengths, applies mask if requested, writes to buffer and flushes.
Assessment: PASS
Severity: LOW
Notes: Thread-safe write path.

## Finding 7

Location: `internal/server/server.go:26-56`
Claimed Behavior: WebSocket message routing with text echo, binary reversal, and ping/close control frame handling.
Observed Implementation: Handles `OpClose` (echoes close and terminates), `OpPing` (echoes pong), `OpText` (prepends echo prefix), `OpBinary` (reverses byte slice).
Assessment: PASS
Severity: LOW
Notes: Meets protocol demonstration requirements.
