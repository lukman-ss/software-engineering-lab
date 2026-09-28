# Code Audit

## Finding 1

Location: `internal/sse/sse.go:19-36` — `Event.Format()`
Claimed Behavior: Formats SSE events per WHATWG spec with `id`, `event`, `retry`, multi-line `data` prefixing, and double-newline terminator.
Observed Implementation: Correctly prefixes multi-line data by splitting on `\n`. ID field guarded by `> 0` (IDs start from 1). Event/retry fields emitted when non-zero/non-empty. Double-newline terminator appended via `sb.WriteString("\n")` after the last `data:` line (which also has a trailing `\n`). This produces the standard `\n\n` block delimiter.
Assessment: PASS
Severity: LOW
Notes: ID guard `> 0` is reasonable since Hub starts `nextID` at 1, but means ID 0 is silently omitted. Not a behavioral problem for this lab.

---

## Finding 2

Location: `internal/sse/sse.go:52-72` — `Hub.Broadcast()`
Claimed Behavior: Thread-safe broadcast to all registered clients with history buffering.
Observed Implementation: Acquires full write lock, assigns monotonic ID, appends to history, then sends to all client channels via non-blocking `select`/`default`. If client channel is full (capacity 64), event is dropped silently for that client.
Assessment: PASS
Severity: LOW
Notes: Silent channel-full drop is intentional for slow clients. Channel buffer size of 64 is adequate for lab purposes. History grows unbounded — acceptable for a lab but worth noting.

---

## Finding 3

Location: `internal/sse/sse.go:74-96` — `Hub.Subscribe()`
Claimed Behavior: Subscribe returns channel for live events plus replay slice of missed events (since `lastID`). Returns unsubscribe closure.
Observed Implementation: Acquires write lock (needed since it writes to `clients` map). Computes replay slice correctly: iterates history, includes events where `evt.ID > lastID`. Returns `ch`, replay slice, and unsubscribe closure. Unsubscribe deletes from map under write lock and closes channel. This is safe: close happens under the same lock that Broadcast uses to range over clients.
Assessment: PASS
Severity: LOW
Notes: Replay slice is taken before the live channel is registered, creating a potential small gap: events broadcast between end of replay-iteration and channel registration would be in history but also queued in channel (double delivery). However, channel is registered immediately after replay is computed within the same lock, so no gap actually exists.

---

## Finding 4

Location: `internal/sse/sse.go:98-151` — `Hub.ServeHTTP()`
Claimed Behavior: Streams SSE events over HTTP, replays missed events from `Last-Event-ID` header, emits keep-alive comments every 15 seconds.
Observed Implementation: Sets required headers (`Content-Type: text/event-stream`, `Cache-Control: no-cache`, `Connection: keep-alive`, `X-Accel-Buffering: no`). Reads `Last-Event-ID` header with safe `strconv.Atoi`. Calls `Subscribe`, defers `unsub`. Writes replay events then flushes. Main select loop handles context cancellation, live events (with individual flushes), and 15s heartbeat comments.
Assessment: PASS
Severity: LOW
Notes: Write errors immediately return — this is correct. Keep-alive comment `": keep-alive\n\n"` uses SSE comment syntax (`:` prefix) which is spec-compliant.

---

## Finding 5

Location: `internal/ws/ws.go:31-76` — `Upgrade()`
Claimed Behavior: RFC 6455 compliant WebSocket handshake: verify `Upgrade: websocket`, read `Sec-WebSocket-Key`, compute `Sec-WebSocket-Accept` as base64(SHA1(key + wsGUID)), respond with 101.
Observed Implementation: Validates `Upgrade` header (exact string match). Validates `Sec-WebSocket-Key` presence. SHA1 computed correctly with the RFC-mandated GUID `258EAFA5-E914-47DA-95CA-C5AB0DC85B11`. Writes 101 response with all required headers. Hijacks connection for raw TCP access.
Assessment: PASS
Severity: LOW
Notes: Does not validate `Sec-WebSocket-Version: 13` from client. In production this matters but for lab purposes it's acceptable.

---

## Finding 6

Location: `internal/ws/ws.go:84-127` — `ReadFrame()`
Claimed Behavior: RFC 6455 frame reading with masking, 7/16/64-bit payload length support.
Observed Implementation: Reads 2-byte header. Extracts opcode, mask bit, payload length. Extends length for 126 (2-byte) and 127 (8-byte) values with `binary.BigEndian`. Reads 4-byte mask key if masked. Reads payload. Applies XOR unmask correctly with `data[i] ^= mask[i%4]`. No locking — read is not protected by `c.mu`.
Assessment: PASS
Severity: LOW
Notes: `ReadFrame` intentionally has no mutex — callers should not read concurrently. This is an acceptable design for this lab (single-reader pattern). The WriteFrame uses `c.mu` because the server writes from a single goroutine per connection in this lab's usage pattern.

---

## Finding 7

Location: `internal/ws/ws.go:129-174` — `WriteFrame()`
Claimed Behavior: Writes RFC 6455 frame with FIN=1, correct length encoding, optional masking.
Observed Implementation: FIN bit set via `byte(0x80) | opcode`. Payload length encoded correctly for all three ranges. When `mask=true`, uses a static key `{0x12, 0x34, 0x56, 0x78}`. Static mask key is not random — RFC 6455 requires mask keys be unpredictable. However, for a lab demonstrating the protocol mechanism (not security), this is a simplification.
Assessment: WARNING
Severity: MEDIUM
Notes: Static mask key violates RFC 6455 Section 5.3 requirement for client-to-server masking randomness. In this lab context the server does not validate mask key randomness. Does not affect protocol correctness in the test harness. Mark as non-blocking warning.

---

## Finding 8

Location: `internal/ws/ws.go:176-201` — `Dial()`
Claimed Behavior: Client-side WebSocket handshake over raw TCP.
Observed Implementation: Uses a static, hardcoded `Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==` (the RFC 6455 example key). Sends handshake, reads response, checks for `HTTP/1.1 101`. Uses a single `rawConn.Read` to read the 101 response — this could return partial data if the TCP buffer is slow.
Assessment: WARNING
Severity: MEDIUM
Notes: Single `Read` for HTTP response is not robust. The static key is fine for the lab (the GUID-based `Sec-WebSocket-Accept` doesn't matter for loopback/test usage), but a production-quality implementation would need a buffered HTTP response parser. In this lab, all connections are loopback so in practice this works. No test failures observed.

---

## Finding 9

Location: `internal/server/server.go:33-56` — WebSocket handler loop
Claimed Behavior: Handles text echo and binary reverse for WebSocket connections.
Observed Implementation: Continuous `for` loop calling `ReadFrame`. Handles `OpClose` (sends Close, returns), `OpPing` (sends Pong), `OpText` (echoes with `"echo: "` prefix), `OpBinary` (reverses payload). All write errors are silently ignored with `_ =`.
Assessment: WARNING
Severity: MEDIUM
Notes: Write errors on echo are silently discarded. If the connection is broken, the next `ReadFrame` will catch it. Acceptable for lab. However, `OpContinuation` (fragmented frames) is not handled — falls through `switch` silently. For this lab, no fragmented frames are sent, so this does not affect correctness.

---

## Finding 10

Location: `internal/sse/sse.go:39,85` — SSE Hub history unbounded growth
Claimed Behavior: History buffer enables Last-Event-ID replay.
Observed Implementation: `history []Event` is appended to on every `Broadcast` call, never trimmed.
Assessment: WARNING
Severity: LOW
Notes: Memory leak for long-running servers. Not a correctness issue for lab.

---

## Finding 11

Location: Concurrency analysis — SSE Hub
Claimed Behavior: Concurrent safe.
Observed Implementation: `Broadcast` takes write lock; `Subscribe` takes write lock; `ServeHTTP` reads from returned channel (no lock needed since channel is per-client). The `unsubscribe` closure acquires write lock to delete and close. `Broadcast` uses range-over-map under lock — safe since `unsubscribe` also takes lock to delete. No deadlocks observed in race-detector run.
Assessment: PASS
Severity: LOW
Notes: Race detector test passed cleanly.

---

## Finding 12

Location: Concurrency analysis — WebSocket `Conn`
Claimed Behavior: `mu sync.Mutex` in `Conn` guards writes.
Observed Implementation: `WriteFrame` locks `c.mu`. `ReadFrame` does not lock `c.mu`. `Close` locks `c.mu`. In the server loop, reads and writes are sequential (read then write) in a single goroutine — no concurrent read/write from multiple goroutines exists in the server implementation. In tests, same pattern. Race detector passes.
Assessment: PASS
Severity: LOW
Notes: The server pattern is single-goroutine per connection (read loop with writes inline), which is safe. If `WriteFrame` were called from a separate goroutine, the partial protection (write mutex, no read mutex) could theoretically be incomplete, but that usage does not occur in this lab.
