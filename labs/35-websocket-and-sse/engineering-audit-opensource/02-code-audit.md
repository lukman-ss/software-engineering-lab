## Finding 1
Location: internal/ws/ws.go:31-76
Claimed Behavior: Upgrade validates handshake, returns Conn.
Observed Implementation: Checks headers, hijacks, writes response, returns Conn with mutex.
Assessment: PASS
Severity: LOW
Notes: Proper error handling, mutex protects Close/Write.

## Finding 2
Location: internal/ws/ws.go:84-174
Claimed Behavior: ReadFrame parses frames, WriteFrame encodes, supports masking.
Observed Implementation: Implements RFC 6455 framing, mask handling, payload length handling.
Assessment: PASS
Severity: LOW
Notes: No unchecked errors, mask key fixed but acceptable for demo.

## Finding 3
Location: internal/sse/sse.go:52-71
Claimed Behavior: Broadcast creates event, increments ID, notifies clients.
Observed Implementation: Locks, creates Event, appends to history, non-blocking send.
Assessment: PASS
Severity: LOW

## Finding 4
Location: internal/sse/sse.go:74-96
Claimed Behavior: Subscribe replays missed events, registers client.
Observed Implementation: Returns replay slice, channel, unsubscribe.
Assessment: PASS
Severity: LOW

## Finding 5
Location: internal/server/server.go:24-56
Claimed Behavior: /ws endpoint upgrades and echoes messages.
Observed Implementation: Uses ws.Upgrade, loops reading frames, handles OpClose, OpPing, OpText, OpBinary.
Assessment: PASS
Severity: LOW

## Finding 6
Location: internal/server/server.go:58-66
Claimed Behavior: /publish endpoint broadcasts notification.
Observed Implementation: Reads query param, defaults, broadcasts, writes ID.
Assessment: PASS
Severity: LOW

## Finding 7
Location: internal/sse/sse.go:98-150
Claimed Behavior: SSE HTTP streaming with replay, keep-alive.
Observed Implementation: Writes headers, handles Last-Event-ID, replay, ticker keep-alive.
Assessment: PASS
Severity: LOW

## Finding 8
Location: cmd/demo/main.go:24-95
Claimed Behavior: Demonstrates SSE resumption and WebSocket echo.
Observed Implementation: Broadcasts, client with Last-Event-ID=1 receives events 2+, WebSocket text/binary frames echoed correctly.
Assessment: PASS
Severity: LOW
