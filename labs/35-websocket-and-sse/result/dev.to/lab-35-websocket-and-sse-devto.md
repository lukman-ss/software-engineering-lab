---
title: "WebSocket vs Server-Sent Events: Selecting Real-Time Protocols via Directionality and Complexity"
published: false
description: "Explore WebSocket vs SSE in Go: RFC 6455 full-duplex binary framing vs WHATWG text/event-stream push with Last-Event-ID resumption."
tags: go, webdev, performance, architecture
---

# WebSocket vs Server-Sent Events: Selecting Real-Time Protocols via Directionality and Complexity

Selecting between WebSocket and Server-Sent Events dictates connection complexity, proxy resilience, and operational overhead. Lab 35 implements both protocols from scratch in pure Go to expose their architectural trade-offs.

> Do not pay bidirectional connection tax for unidirectional server push.

## The Modern Real-Time Protocol Dilemma

Engineering teams often treat WebSocket as the default choice for all real-time requirements. Adopting WebSocket for strictly unidirectional use cases introduces unnecessary operational overhead. Server-Sent Events offers standard HTTP push that is lightweight and proxy-friendly.

## Core Mental Model: Directionality Rule

Architectural selection reduces to one question: does the client send messages upstream independently over the same channel? If communication flows strictly server-to-client, SSE is optimal. WebSocket is reserved for genuine full-duplex requirements.

## Lab 35 Comparative Architecture

Lab 35 hosts both protocols under a single Go HTTP ServeMux. The zero-dependency implementation reveals how standard library interfaces handle persistent HTTP streaming via http.Flusher and raw protocol transitions via http.Hijacker.

## Anatomy of Server-Sent Events (SSE)

The WHATWG specification defines SSE as a long-lived HTTP response using text/event-stream. Each event comprises id, event, retry, and data fields terminated by double newlines. Operating over pure HTTP ensures immediate compatibility with web infrastructure.

## SSE Event Formatting in Go

The Format() method enforces WHATWG rules. When payloads contain newlines, each line must carry a 'data: ' prefix. Client parsers assemble these lines sequentially. The final double newline marks the completion of the event frame.

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

SSE's defining advantage is native browser reconnection coupled with the Last-Event-ID header. Upon disconnection, the browser retransmits its latest ID. The hub replays missed events from buffer memory before piping real-time events.

```go
func (h *Hub) Subscribe(lastID int) (chan Event, []Event, func()) {
    h.mu.Lock(); defer h.mu.Unlock()
    var replay []Event
    for _, evt := range h.history {
        if evt.ID > lastID { replay = append(replay, evt) }
    }
    ch := make(chan Event, 64)
    h.clients[ch] = struct{}{}
    return ch, replay, func() { ... }
}
```

## HTTP Streaming Loop & Keep-Alive

SSE streaming requires http.Flusher to push frames through Go runtime buffers immediately. The select loop handles client termination contexts, live channel broadcasts, and a 15-second ticker sending ': keep-alive' comments to prevent intermediate proxy timeouts.

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

## Anatomy of WebSocket Protocol (RFC 6455)

WebSocket was standardized by IETF RFC 6455 to provide low-latency full-duplex communication over single TCP connections. Following an HTTP 101 Switching Protocols handshake, the underlying socket shifts entirely to raw binary framing.

## WebSocket Upgrade & Connection Hijacking

The Upgrade() function computes the Sec-WebSocket-Accept challenge. By invoking http.Hijacker, Go transfers socket ownership from net/http. The server writes 101 Switching Protocols headers directly to TCP buffers and wraps the connection into a custom struct.

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

## RFC 6455 Binary Frame Structure

RFC 6455 framing minimizes wire overhead with a 2-byte base header. Payload lengths encode dynamically across 7-bit, 16-bit, and 64-bit boundaries. Client-to-server masking protects intermediaries against cache poisoning vulnerabilities.

## Frame Parsing & XOR Unmasking

The ReadFrame() method decodes wire bytes from the stream. When the mask bit is set, the server reads the 4-byte key and performs in-place XOR restoration. This guarantees protocol fidelity with minimal memory allocation.

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

## WriteFrame: Outbound Frame Construction

The WriteFrame() method serializes outbound binary frames. Mutex synchronization guarantees that frame writes remain atomic across concurrent goroutines, avoiding frame byte interleaving on TCP streams.

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

## Server Loop: Text Echo & Binary Reversal

The WebSocket event loop routes incoming frames by opcode. Text frames receive an echo prefix, while binary payloads undergo byte-order reversal. This reversal proves transparent raw binary transfer without encoding loss.

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

## Demo Execution & Verification Results

Running cmd/demo validates theoretical claims. SSE resumption honors Last-Event-ID: 1 by skipping initial events and replaying the subsequent update. WebSocket mirrors text and reverses binary DEADBEEF into EFBEADDE flawlessly.

```go
> 1. SSE: Unidirectional Server Push with Last-Event-ID
Client connected with Last-Event-ID: 1. Received resumed events:
  [SSE Stream] id: 2
  [SSE Stream] event: news
  [SSE Stream] retry: 2000
  [SSE Stream] data: Update: Faster routing in net/http
```

```go
> 2. WebSocket: Full-Duplex Text and Binary Protocol
Client sending WS text frame: "Hello RFC 6455"
Server responded with Opcode 0x1 (Text): "echo: Hello RFC 6455"
Client sending WS binary frame: DEADBEEF
Server responded with Opcode 0x2 (Binary): EFBEADDE
Binary frame payload verified successfully.
```

## Test Suite Verification & Concurrency

All test suites pass under the Go race detector. The concurrency test confirms that the SSE Hub non-blocking select pattern safely serves 10 parallel subscribers across 5 concurrent broadcasts without data races or lockups.

```sh
=== RUN   TestSSE_Formatting            -> PASS (0.00s)
=== RUN   TestSSE_LastEventID_Resumption -> PASS (0.00s)
=== RUN   TestWebSocket_TextAndBinary    -> PASS (0.00s)
=== RUN   TestSSE_Concurrency           -> PASS (0.05s)
PASS: go test -race ./... (1.559s - Race Clean)
```

## Evaluation Matrix: WebSocket vs SSE

The evaluation matrix highlights the architectural trade-offs. WebSocket offers full-duplex streaming and binary payloads at the cost of operational proxy complexity. SSE delivers native network resilience and automatic HTTP/2 multiplexing.

Dimensi           | WebSocket (RFC 6455)      | Server-Sent Events (SSE)
|---|---|---|
Arah Komunikasi   | Full-Duplex (2 Arah)      | Unidirectional (Server->Client)
Transport / Pipa  | Raw TCP (setelah Upgrade) | Standard HTTP Streaming
Tipe Payload      | UTF-8 Teks & Raw Biner    | UTF-8 Teks Murni Saja
Auto-Reconnect    | Manual / Custom Client    | Bawaan Browser EventSource
Resumption State  | Tidak ada di spec level   | Bawaan via Last-Event-ID
Proxy / LB Compat | Butuh Sticky / WS Bypass  | Standar HTTP (No buffering)
HTTP/2 Multiplex  | Butuh RFC 8441 CONNECT    | Native HTTP/2 Multiplexing

## Proxy, Firewall & Gateway Interoperability

Long-lived connections are vulnerable to gateway buffering and idle terminations. WebSocket often encounters strict enterprise firewalls suspicious of protocol switching. SSE requires disabling proxy buffering via X-Accel-Buffering: no headers.

## Impact of HTTP/2 & Stream Multiplexing

HTTP/2 amplifies SSE's viability by multiplexing many event streams across a single TCP connection. Conversely, WebSocket over HTTP/2 mandates RFC 8441 Extended CONNECT mechanics because traditional connection upgrades are prohibited.

## Case Study: Notifications & LLM Streaming

Production adoption reflects directionality rules. Generative AI providers stream tokens via SSE because traffic is strictly server-to-client following the initial prompt. Real-time collaborative tools like Figma depend on WebSocket for low-latency bidirectional interaction.

## Lab 35 Model Boundaries & Constraints

This lab focuses on educational clarity. Static masking keys simplify byte-level inspection, and the SSE hub retains history without buffer eviction. Production deployments require bounded ring buffers, cryptographically random masks, and message brokers.

## Protocol Decision Checklist

Utilize this checklist during architectural design phases. If upstream requests can be satisfied via standard HTTP POST while clients listen passively, SSE is the superior path. Verify reverse proxy settings and OS file descriptor limits.

## Three Core Real-Time Principles

Resilient systems choose transport protocols matching traffic patterns. Avoid over-engineering by deploying WebSocket indiscriminately. Reserving SSE for unidirectional streams and WebSocket for true full-duplex interactions yields scalable architectures.

## Source

[Software Engineering Lab 35 - WebSocket vs Server-Sent Events](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/35-websocket-and-sse), snapshot commit `b48fe6c3`. Source: README, implementation, tests, and execution records. Code excerpts belong in their original package context.
