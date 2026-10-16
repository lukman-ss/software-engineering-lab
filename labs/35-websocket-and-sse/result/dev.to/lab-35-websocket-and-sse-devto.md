---
title: "WebSocket vs Server-Sent Events: Selecting Real-Time Protocols via Directionality"
published: false
description: "Explore WebSocket vs SSE in Go: RFC 6455 full-duplex binary framing vs WHATWG text/event-stream push with Last-Event-ID resumption."
tags: go, webdev, performance, architecture
---

# WebSocket vs Server-Sent Events: Selecting Real-Time Protocols via Directionality

Selecting between WebSocket and Server-Sent Events dictates connection complexity, proxy resilience, and operational overhead. Lab 35 implements both protocols from scratch in pure Go to expose their architectural trade-offs.

> Directionality dictates the protocol: server-push belongs to SSE, bidirectional belongs to WebSocket.

## Dilema Protokol Real-Time Modern

Engineering teams frequently default to WebSocket for every real-time feature, from live notifications to market tickers. Adopting WebSocket for strictly server-to-client updates introduces stateful TCP connection overhead, manual heartbeat maintenance, sticky reverse proxy routing, and reconnection storm vulnerability. Teams pay the operational penalty of a full-duplex protocol while only utilizing passive unidirectional streaming.

## Mental Model: Directionality Rule

Protocol selection reduces to communication directionality rather than raw speed. When the interaction pattern is strictly server-to-client (notifications, metrics, LLM token streaming), Server-Sent Events (SSE) operates directly over standard HTTP. WebSocket (RFC 6455) is justified only when the client must transmit independent, high-frequency upstream messages over the identical low-latency connection.

## Topologi Arsitektur Lab 35

Lab 35 hosts both protocols side by side under a single Go 1.22 HTTP router with zero external third-party dependencies. The server routes `GET /sse` to an event-stream streaming hub, `GET /ws` to an RFC 6455 connection hijacker, and `POST /publish` to broadcast triggers across active subscribers.

```text
[Client] ─── GET /sse ─────────► [SSE Hub Handler] (text/event-stream + Last-Event-ID)
[Client] ─── GET /ws ──────────► [WS Upgrader]     (101 Switching Protocols + XOR Framing)
[Client] ─── POST /publish ────► [Broadcast Bus]  (In-memory Hub Distribution)
```

## Anatomi Server-Sent Events (SSE)

WHATWG Server-Sent Events standardizes persistent unidirectional text streams over vanilla HTTP. The server returns MIME type `text/event-stream` with `Cache-Control: no-cache` and `Connection: keep-alive`. Each event consists of plaintext fields (`id`, `event`, `retry`, `data`) terminated by a double newline (`\n\n`) delimiter.

## Format Event SSE: Multi-Line Framing

The SSE specification requires payloads with newline characters to split across consecutive `data:` lines. The client EventSource parser reassembles these lines sequentially. The trailing double newline signals the completion of the event frame.

```go
func (e Event) Format() []byte {
    var sb strings.Builder
    if e.ID > 0 { fmt.Fprintf(&sb, "id: %d\n", e.ID) }
    if e.Event != "" { fmt.Fprintf(&sb, "event: %s\n", e.Event) }
    if e.Retry > 0 { fmt.Fprintf(&sb, "retry: %d\n", e.Retry) }
    lines := strings.Split(e.Data, "\n")
    for _, l := range lines {
        fmt.Fprintf(&sb, "data: %s\n", l)
    }
    sb.WriteString("\n")
    return []byte(sb.String())
}
```

## SSE Hub: Resumption via Last-Event-ID

SSE provides native connection recovery. When a network hiccup severs the socket, the browser client reconnects automatically, passing its most recent offset in the `Last-Event-ID` header. The server hub inspects buffer history, replays dropped events (`evt.ID > lastID`), and transitions smoothly to live event streaming.

```go
func (h *Hub) Subscribe(lastID int) (chan Event, []Event, func()) {
    h.mu.Lock(); defer h.mu.Unlock()
    var replay []Event
    for _, evt := range h.history {
        if evt.ID > lastID { replay = append(replay, evt) }
    }
    ch := make(chan Event, 64)
    h.clients[ch] = struct{}{}
    return ch, replay, func() {
        h.mu.Lock()
        delete(h.clients, ch)
        close(ch)
        h.mu.Unlock()
    }
}
```

## Loop Streaming HTTP & Heartbeat

Persistent HTTP streaming requires calling `http.Flusher.Flush()` immediately after writing frames to push bytes through Go runtime buffers. An internal ticker emits `: keep-alive\n\n` comments every 15 seconds. Browser parsers ignore leading colon lines while intermediate reverse proxies keep idle connections alive.

```go
flusher, ok := w.(http.Flusher)
for {
    select {
    case <-ctx.Done(): return
    case evt, ok := <-ch:
        if !ok { return }
        w.Write(evt.Format())
        flusher.Flush()
    case <-ticker.C:
        w.Write([]byte(": keep-alive\n\n"))
        flusher.Flush()
    }
}
```

## Anatomi Protokol WebSocket (RFC 6455)

RFC 6455 provides a transport-independent, full-duplex protocol over a single TCP socket. The connection initiates via an HTTP/1.1 handshake requesting protocol switching with a 16-byte random `Sec-WebSocket-Key`. After returning HTTP `101 Switching Protocols`, HTTP parsing terminates and the socket switches to binary framing.

## Handshake Upgrade & Connection Hijacking

In Go, taking over a raw TCP socket requires invoking `http.Hijacker`. The server computes the challenge response using SHA-1 and the standardized RFC 6455 GUID (`258EAFA5-E914-47DA-95CA-C5AB0DC85B11`), writes the 101 header directly, and assumes complete control over the TCP buffer stream.

```go
h := sha1.New()
h.Write([]byte(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
accept := base64.StdEncoding.EncodeToString(h.Sum(nil))

hj, ok := w.(http.Hijacker)
netConn, bufrw, err := hj.Hijack()
resp := "HTTP/1.1 101 Switching Protocols\r\n" +
    "Upgrade: websocket\r\nConnection: Upgrade\r\n" +
    "Sec-WebSocket-Accept: " + accept + "\r\n\r\n"
bufrw.WriteString(resp); bufrw.Flush()
```

## Struktur Frame Biner RFC 6455

RFC 6455 minimizes wire framing overhead with a base 2-byte header. The header encodes FIN bit, opcode (Text, Binary, Close, Ping, Pong), mask indicator, and dynamic payload length tiers (7-bit, 16-bit uint, or 64-bit uint). Client frames must include a 4-byte masking key to prevent cache poisoning on intermediary gateways.

```text
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-------+-+-------------+-------------------------------+
|F|R|R|R| opcode|M| Payload len |    Extended payload length    |
|I|S|S|S|  (4)  |A|     (7)     |             (16/64)           |
|N|V|V|V|       |S|             |   (if payload len==126/127)   |
| |1|2|3|       |K|             |                               |
+-+-+-+-+-------+-+-------------+ - - - - - - - - - - - - - - - +
|     Extended payload length continued, if payload len == 127  |
+ - - - - - - - - - - - - - - - +-------------------------------+
|                               |Masking-key, if MASK set to 1  |
+-------------------------------+-------------------------------+
| Masking-key (continued)       |          Payload Data         |
+-------------------------------- - - - - - - - - - - - - - - - +
:                     Payload Data continued ...                :
+---------------------------------------------------------------+
```

## Parsing Frame & Operasi XOR Masking

`ReadFrame` reads the wire header, resolves dynamic payload lengths, and reads the 4-byte client mask. Payload restoration is executed in-place via bitwise XOR (`data[i] ^= mask[i%4]`). XOR symmetry allows fast decoding in place without allocating secondary heap buffers.

```go
var head [2]byte
io.ReadFull(c.rw, head[:])
opcode = head[0] & 0x0F
masked := (head[1] & 0x80) != 0
length := uint64(head[1] & 0x7F)
// ... handle extended length ...
if masked {
    io.ReadFull(c.rw, mask[:])
    io.ReadFull(c.rw, data)
    for i := range data {
        data[i] ^= mask[i%4] // Bitwise XOR unmasking
    }
}
```

## WriteFrame: Konstruksi Frame Outbound

`WriteFrame` constructs server outbound frames. The server sets FIN bit `0x80` combined with the target opcode and encodes payload size boundaries. Outbound server frames do not apply client masking. Mutex synchronization guarantees atomic socket writes across concurrent goroutines.

```go
func (c *Conn) WriteFrame(opcode byte, payload []byte, mask bool) error {
    c.mu.Lock(); defer c.mu.Unlock()
    var buf []byte
    b0 := byte(0x80) | (opcode & 0x0F) // FIN=1
    buf = append(buf, b0)
    // ... encode length (7, 16, or 64 bit) ...
    if mask {
        // append maskKey and maskedPayload
    } else {
        buf = append(buf, payload...)
    }
    _, err := c.rw.Write(buf)
    return err
}
```

## Server Loop: Echo Teks & Reversal Biner

The server event loop dispatches frames based on incoming opcodes. Text frames return with an `echo: ` prefix. Binary frames invert payload byte order before transmission. Reversing `0xDEADBEEF` into `0xEFBEADDE` proves transparent raw binary transfer without string transcoding distortions.

```go
switch op {
case ws.OpClose:
    conn.WriteFrame(ws.OpClose, nil, false); return
case ws.OpPing:
    conn.WriteFrame(ws.OpPong, data, false)
case ws.OpText:
    conn.WriteFrame(ws.OpText, append([]byte("echo: "), data...), false)
case ws.OpBinary:
    rev := make([]byte, len(data))
    for i := range data { rev[i] = data[len(data)-1-i] }
    conn.WriteFrame(ws.OpBinary, rev, false)
}
```

## Eksekusi Demo: Bukti Terminal Langsung

Running `cmd/demo` validates both protocol mechanisms deterministically. The SSE client transmits `Last-Event-ID: 1`, skipping event 1 and receiving replayed event 2 seamlessly. The WebSocket client connects, echoes UTF-8 text, and verifies byte-for-byte binary reversal.

```sh
> 1. SSE: Unidirectional Server Push with Last-Event-ID
Client connected with Last-Event-ID: 1. Received resumed events:
  [SSE Stream] id: 2
  [SSE Stream] event: news
  [SSE Stream] retry: 2000
  [SSE Stream] data: Update: Faster routing in net/http
> 2. WebSocket: Full-Duplex Text and Binary Protocol
Client sending WS text frame: "Hello RFC 6455"
Server responded with Opcode 0x1 (Text): "echo: Hello RFC 6455"
Client sending WS binary frame: DEADBEEF
Server responded with Opcode 0x2 (Binary): EFBEADDE
Binary frame payload verified successfully.
```

## Hasil Pengujian Test Suite & Concurrency

The protocol test suite in `tests/protocol_test.go` executes cleanly with the Go 1.22 race detector. Concurrency tests confirm that non-blocking channel broadcasting safely delivers 5 asynchronous events across 10 concurrent subscribers without deadlocks or race conditions.

```sh
=== RUN   TestSSE_Formatting            -> PASS (0.00s)
=== RUN   TestSSE_LastEventID_Resumption -> PASS (0.00s)
=== RUN   TestWebSocket_TextAndBinary    -> PASS (0.00s)
=== RUN   TestSSE_Concurrency           -> PASS (0.05s)
PASS: go test -race ./... (1.559s - Race Clean)
```

## Matriks Evaluasi: WebSocket vs SSE

Evaluating WebSocket against Server-Sent Events reveals clear architectural trade-offs across 8 critical dimensions.

| Dimension | WebSocket (RFC 6455) | Server-Sent Events (SSE) |
| --- | --- | --- |
| Directionality | Full-Duplex (Bidirectional) | Unidirectional (Server -> Client) |
| Transport Pipeline | Raw TCP (Post-HTTP Upgrade) | Standard HTTP Streaming |
| Payload Support | UTF-8 Text and Raw Binary | Plain UTF-8 Text Only |
| Auto-Reconnection | Manual Client Implementation | Native Browser EventSource Support |
| State Resumption | Custom Application Layer | Built-in via Last-Event-ID |
| Proxy & Load Balancer | Requires Sticky / Protocol Bypass | Standard HTTP (Disable Buffering) |
| HTTP/2 Multiplexing | Requires RFC 8441 CONNECT | Native Out-of-the-Box Multiplexing |
| Firewall Traversal | Can be blocked by Deep Inspection | Standard Web Port 80/443 Native |

## Tantangan Proxy, Firewall & Gateway

Long-lived connections face hostile intermediate network middleboxes. WebSocket connections can be severed by corporate firewalls inspecting non-HTTP TCP traffic or dropped by load balancers due to idle timeouts. SSE requires disabling reverse proxy buffering by passing `X-Accel-Buffering: no` to Nginx so frames flush downstream immediately.

## Dampak HTTP/2 & Multiplexing Aliran

HTTP/2 natively enhances SSE by multiplexing dozens of concurrent event streams across a single physical TCP connection, eliminating the browser HTTP/1.1 limit of 6 connections per host. Conversely, WebSocket over HTTP/2 requires implementing RFC 8441 Extended CONNECT, which lacks universal edge proxy support.

## Studi Kasus Nyata: Notifikasi & AI Tokens

Production systems align with directionality rules. Generative AI services (OpenAI, Anthropic) stream tokens via SSE because the request prompt completes upfront and token delivery flows strictly server-to-client. Collaborative editors (Figma) and interactive messaging platforms rely on WebSocket for continuous bidirectional event exchange.

## Batasan Model & Implementasi Lab 35

This lab deliberately constrains scope to highlight core protocol mechanics. Client masking keys use a constant byte sequence rather than cryptographically random bytes. SSE history uses an unbounded in-memory slice without eviction policies. Production environments require bounded circular ring buffers, dynamic crypto-random masking, and multi-node Redis Pub/Sub backplanes.

## Checklist Keputusan Arsitektur Real-Time

Before committing to a real-time protocol, review this architectural evaluation checklist:

- [x] Does the client need to send data upstream on the same socket? (Yes -> WS, No -> SSE)
- [x] Does the payload require arbitrary binary encoding? (Yes -> WS, No -> SSE)
- [x] Is automatic connection recovery and event resumption required out of the box? (SSE)
- [x] Are requests routed through restrictive enterprise firewalls or standard CDNs? (SSE)
- [x] Has `X-Accel-Buffering: no` been configured on edge reverse proxies for SSE?
- [x] Are heartbeat comments or ping/pong frames scheduled to defeat idle connection drops?

## Tiga Prinsip Inti Komunikasi Real-Time

Resilient real-time architectures respect three axioms:

1. **Directionality is King**: Do not pay bidirectional operational taxes for unidirectional data flows. SSE is the ideal tool for notifications, feeds, and streaming generation.
2. **Leverage Platform Primitives**: Native browser `EventSource` handles reconnection and offset tracking automatically; HTTP/2 eliminates connection pooling overhead.
3. **Protocol Respects Boundaries**: WebSocket remains the industry standard when low-latency upstream client input and binary payloads are genuinely required.

## Source

[Software Engineering Lab 35 - WebSocket and Server-Sent Events](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/35-websocket-and-sse), snapshot commit `b48fe6c30d0b3635f36b8549c5621bbf5ce7c5f8`. Source: README, `internal/sse`, `internal/ws`, `internal/server`, `cmd/demo`, and `tests/protocol_test.go`. Code excerpts belong in their original package context.
