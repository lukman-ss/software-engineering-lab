# Code Audit — labs/35-websocket-and-sse

Scope: implementation + tests only. Research not audited (pipeline override).

## Finding 1

Location: internal/sse/sse.go:19-36 (Event.Format)
Claimed Behavior: SSE frames with id/event/retry/multi-line data.
Observed Implementation: Writes `id:` only if ID>0, `event:` if non-empty, `retry:` if >0, one `data:` line per split line, trailing blank line. Matches WHATWG framing shape.
Assessment: PASS
Severity: LOW
Notes: None.

## Finding 2

Location: internal/sse/sse.go:52-96 (Hub.Broadcast/Subscribe)
Claimed Behavior: Buffer replay using Last-Event-ID; thread-safe subscribers.
Observed Implementation: Broadcast assigns incrementing ID, Retry 2000, appends history, non-blocking send under mu. Subscribe replays evt.ID > lastID with buffered ch(64); unsubscribe deletes + closes under mu. Broadcast and unsubscribe serialized on same mutex; no send-on-closed reachable.
Assessment: PASS
Severity: LOW
Notes: Unbounded history and silent drop to slow subscribers are disclosed in engineering/02-implementation-notes.md. Fine for lab.

## Finding 3

Location: internal/sse/sse.go:98-151 (ServeHTTP)
Claimed Behavior: text/event-stream push with Last-Event-ID resumption + keep-alive.
Observed Implementation: Flusher guard, correct Content-Type/Cache-Control/Connection/X-Accel-Buffering headers, invalid Last-Event-ID tolerated as 0 (full replay), replay written before live stream, ctx.Done exits, 15s `: keep-alive` heartbeat, write errors return.
Assessment: PASS
Severity: LOW
Notes: None.

## Finding 4

Location: internal/ws/ws.go:31-76 (Upgrade)
Claimed Behavior: RFC 6455 upgrader.
Observed Implementation: Accept key = base64(sha1(key+GUID)) correct; hijack + 101 correct. Upgrade header compared case-sensitively; Sec-WebSocket-Version not validated.
Assessment: WARNING
Severity: LOW
Notes: Works with demo/tests; strictness gap only matters against diverse clients.

## Finding 5

Location: internal/ws/ws.go:84-127 (ReadFrame)
Claimed Behavior: RFC 6455 framing incl. 126/127 extended lengths + unmasking.
Observed Implementation: Correct 7-bit/16-bit/64-bit length decode, mask application, io.ReadFull. No max-payload cap: `make([]byte, length)` trusts remote 64-bit length (OOM/panic on malicious peer). No FIN/continuation reassembly; unknown opcodes surfaced raw to caller.
Assessment: WARNING
Severity: MEDIUM
Notes: Fragmentation gap is disclosed in engineering notes; missing length cap is not. Lab-only risk, no remote exposure in demo/tests.

## Finding 6

Location: internal/ws/ws.go:129-174 (WriteFrame)
Claimed Behavior: Server frame writer, concurrent-safe.
Observed Implementation: FIN=1, correct length encoding, mu-serialized writes, flush when available. Client-mask path uses fixed key 0x12,0x34,0x56,0x78 instead of random.
Assessment: WARNING
Severity: LOW
Notes: Fixed mask is disclosed in engineering notes; interoperable, just not RFC-conformant randomness. Server-to-client unmasked path (mask=false) correct.

## Finding 7

Location: internal/ws/ws.go:176-202 (Dial)
Claimed Behavior: Test/demo dial helper.
Observed Implementation: `urlStr` parameter ignored; request line hardcodes `GET /ws`. Handshake check reads once (≤1024B) and compares first 12 bytes to `HTTP/1.1 101`; Sec-WebSocket-Accept not verified.
Assessment: WARNING
Severity: LOW
Notes: Callers always pass "/ws" so behavior is correct; API signature is misleading. No retry/version logic claimed, none present.

## Finding 8

Location: internal/server/server.go:26-56 (/ws handler)
Claimed Behavior: Text echo + binary echo (reversed), ping/pong, close handshake.
Observed Implementation: Text prefixed `echo: `, binary byte-reversed, ping→pong, close→close+return, read-error→return with deferred conn.Close(). Write errors discarded (`_ =`); continuation/other opcodes silently ignored.
Assessment: PASS
Severity: LOW
Notes: Matches claimed echo behavior; swallowed write errors acceptable for demo echo server.

## Finding 9

Location: internal/server/server.go:58-66 (/publish)
Claimed Behavior: (Undocumented helper.)
Observed Implementation: Broadcasts `notification` event, returns `Published event ID %d`. Functional, untested, absent from README endpoint list.
Assessment: WARNING
Severity: LOW
Notes: Tracked as DOC_CODE_MISMATCH in 04-docs-vs-code.md / 05-gaps.md.

Overall: core SSE replay and WS text/binary echo verified by reading; no FAIL-level defect. Warnings are edge-case/robustness gaps, all lab-scoped.
