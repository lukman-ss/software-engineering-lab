# Research Report

## Research Question

> When should a real-time system use WebSocket (full-duplex) versus Server-Sent Events (SSE, unidirectional server→client), considering protocol correctness, operational complexity, HTTP/2 compatibility, and scaling to ~100,000 concurrent connections?

## Executive Summary

The core decision between WebSocket and SSE hinges on communication direction. WebSocket is a distinct, non-HTTP application protocol (RFC 6455) that provides a full-duplex, binary-and-text frame channel over a single TCP connection; it requires a specialized HTTP/2 bootstrap (RFC 8441) to avoid connection-level Upgrade headers. SSE is a standard, persistent HTTP response with a strict `text/event-stream` format that is natively unidirectional (server→client) and is automatically managed by the browser's `EventSource` API (auto-reconnect, `Last-Event-ID`). For one-way server push scenarios like notifications, live dashboards, and LLM token streaming, SSE offers strictly lower operational complexity and better proxy/cache interoperability than WebSocket. The lab-specified 100,000 concurrent connection scenario is constrained primarily by file descriptor limits and per-connection memory in the backend runtime, not by the protocol choice itself; however, SSE's standard HTTP framing simplifies load balancing and avoids sticky-session requirements for read-only traffic.

## Findings

### Finding 1: WebSocket is a full-duplex, non-HTTP protocol

Claim: WebSocket is a two-way, bidirectional communication channel over a single TCP connection that operates outside of the standard HTTP/1.1 request-response lifecycle.

Evidence:
- RFC 6455 states: "The WebSocket Protocol enables two-way communication between a client and a server."
- The opening handshake is an HTTP/1.1 Upgrade request (`Connection: Upgrade`, `Upgrade: websocket`) resulting in a `101 Switching Protocols`.
- "Once the client and server have both sent their handshakes... the data transfer part starts. This is a two-way communication channel where each side can, independently from the other, send data at will."

Sources:
- RFC 6455 – The WebSocket Protocol (December 2011)
- URL: https://datatracker.ietf.org/doc/html/rfc6455

Confidence: HIGH
Corroborated By: MDN WebSocket documentation, WHATWG HTML Standard.

### Finding 2: WebSocket supports binary and UTF-8 text payloads

Claim: The base WebSocket protocol defines distinct frame types for binary data and text (UTF-8) data.

Evidence:
- RFC 6455 states: "there are types for textual data (which is interpreted as UTF-8 [RFC3629] text), binary data (whose interpretation is left up to the application)."

Sources:
- RFC 6455 – The WebSocket Protocol (December 2011)
- URL: https://datatracker.ietf.org/doc/html/rfc6455

Confidence: HIGH
Corroborated By: MDN WebSocket documentation.

### Finding 3: WebSocket does not provide automatic reconnection

Claim: The WebSocket protocol and browser API do not mandate automatic reconnection; the client must implement reconnect logic.

Evidence:
- RFC 6455 Section 7.2.3 describes recovery from abnormal closure but frames it as an implementation concern, not a mandated algorithmic behavior.
- MDN WebSocket documentation does not list auto-reconnect as part of the standard WebSocket interface.

Sources:
- RFC 6455 – The WebSocket Protocol (December 2011)
- URL: https://datatracker.ietf.org/doc/html/rfc6455
- MDN: WebSocket – Web APIs
- URL: https://developer.mozilla.org/en-US/docs/Web/API/WebSocket

Confidence: HIGH
Corroborated By: Common knowledge among front-end developers and third-party WebSocket libraries (e.g., `ReconnectingWebSocket`) that exist specifically to fill this gap.

### Finding 4: SSE is strictly unidirectional server-to-client

Claim: Server-Sent Events deliver data only from the server to the client; there is no native mechanism for client-to-server messaging over the SSE stream.

Evidence:
- MDN EventSource states: "Unlike WebSockets, server-sent events are unidirectional; that is, data messages are delivered in one direction, from the server to the client."
- The WHATWG HTML Standard defines the protocol as a server-push technology using the MIME type `text/event-stream`.

Sources:
- MDN: EventSource – Web APIs
- URL: https://developer.mozilla.org/en-US/docs/Web/API/EventSource
- WHATWG HTML Living Standard Section 9.2
- URL: https://html.spec.whatwg.org/multipage/server-sent-events.html

Confidence: HIGH
Corroborated By: RFC 8441 (which describes WebSocket as the solution for bidirectional HTTP/2 communication).

### Finding 5: SSE supports native automatic reconnection and event identification

Claim: The browser-managed `EventSource` object automatically reconnects on dropped connections, supports a configurable reconnection time, and transmits the last event ID to the server on reconnection.

Evidence:
- WHATWG HTML Standard Section 9.2.3: "When a user agent is to reestablish the connection, the user agent must run the following steps... Wait a delay equal to the reconnection time of the event source."
- Section 9.2.4: The `Last-Event-ID` HTTP request header "reports an EventSource object's last event ID string to the server when the user agent is to reestablish the connection."
- MDN confirms the `onmessage`, `onerror`, `onopen` lifecycle.

Sources:
- WHATWG HTML Living Standard (Updated September 2026)
- URL: https://html.spec.whatwg.org/multipage/server-sent-events.html
- MDN: EventSource – Web APIs
- URL: https://developer.mozilla.org/en-US/docs/Web/API/EventSource

Confidence: HIGH
Corroborated By: Consistent behavior across all modern browser engines (Firefox, Chrome, Safari, Edge).

### Finding 6: SSE uses UTF-8 text only; no binary support

Claim: Server-Sent Events are restricted to UTF-8 text; the spec provides no mechanism for binary payloads.

Evidence:
- WHATWG HTML Standard Section 9.2.1: "Event streams are always decoded as UTF-8. There is no way to specify another character encoding."

Sources:
- WHATWG HTML Living Standard
- URL: https://html.spec.whatwg.org/multipage/server-sent-events.html

Confidence: HIGH
Corroborated By: MDN EventSource documentation.

### Finding 7: WebSocket over HTTP/2 requires RFC 8441 (Extended CONNECT)

Claim: Because HTTP/2 forbids connection-wide headers such as `Upgrade` and `Connection`, the WebSocket handshake cannot operate over plain HTTP/2. RFC 8441 defines a new Extended CONNECT method that uses the `:protocol = websocket` pseudo-header to establish a WebSocket tunnel on a single HTTP/2 stream.

Evidence:
- RFC 8441 Section 1: "HTTP/2 does not allow connection-wide header fields or status codes, such as the Upgrade and Connection request-header fields or the 101 (Switching Protocols) response code. These are all required by the [RFC6455] opening handshake."
- Section 5: "The :protocol pseudo-header field MUST be included in the CONNECT request, and it MUST have a value of 'websocket'."

Sources:
- RFC 8441 – Bootstrapping WebSockets with HTTP/2 (September 2018)
- URL: https://datatracker.ietf.org/doc/html/rfc8441
- RFC 7540 (May 2015, obsoleted by RFC 9113 June 2022) / RFC 9113 – HTTP/2
- URL: https://datatracker.ietf.org/doc/html/rfc9113

Confidence: HIGH

### Finding 8: SSE requires no protocol-level extension for HTTP/2

Claim: Because SSE is standard HTTP semantics, it runs natively over HTTP/2 with the connection's built-in stream multiplexing. No specialized extension like RFC 8441 is required.

Evidence:
- SSE uses persistent, chunked HTTP responses with the `text/event-stream` MIME type.
- HTTP/2's multiplexing applies to all HTTP semantics, including long-lived responses, without modification.
- MDN notes that the HTTP/2 default concurrent stream limit (defaulting to 100) mitigates the traditional HTTP/1.1 6-connection-cap.

Sources:
- WHATWG HTML Living Standard Section 9.2
- URL: https://html.spec.whatwg.org/multipage/server-sent-events.html
- RFC 9113 Sections 5 and 8.1 (Streams and Multiplexing, HTTP message framing) — obsoletes RFC 7540 with identical multiplexing semantics
- URL: https://datatracker.ietf.org/doc/html/rfc9113
- MDN: EventSource – Web APIs
- URL: https://developer.mozilla.org/en-US/docs/Web/API/EventSource

Confidence: HIGH
Corroborated By: MDN.

### Finding 9: Browser connection limits constrain SSE (not WebSocket) over HTTP/1.1

Claim: Browsers enforce a maximum number of simultaneous HTTP/1.1 connections per origin, historically limiting SSE to roughly six concurrent streams per browser across all tabs to that origin. HTTP/2 removes this bottleneck with negotiated stream limits (default 100).

Evidence:
- MDN EventSource explicitly states: "When not used over HTTP/2, SSE suffers from a limitation to the maximum number of open connections... set to a very low number (6)... the limit is per browser + domain."

Sources:
- MDN: EventSource – Web APIs
- URL: https://developer.mozilla.org/en-US/docs/Web/API/EventSource
- RFC 7540 Section 5.1.2 (Stream Concurrency)
- URL: https://datatracker.ietf.org/doc/html/rfc7540

Confidence: HIGH (MDN cites Chrome/Firefox behavior).

### Finding 10: Reverse proxy configuration must account for buffering and timeouts

Claim: Reverse proxies and load balancers often buffer responses by default, which can delay or break real-time event delivery. Administrators must tune `proxy_buffering` (disable for SSE) and `proxy_read_timeout` appropriately.

Evidence:
- NGINX's proxy module documentation defines `proxy_buffering` and `proxy_read_timeout` as relevant to forwarding long-lived responses.
- AWS Elastic Load Balancing documentation confirms WebSocket (but not SSE) requires sticky session support, while standard HTTP responses (including SSE) can leverage HTTP/2 multiplexing.

Sources:
- NGINX Docs: ngx_http_proxy_module
- URL: https://nginx.org/en/docs/http/ngx_http_proxy_module.html
- AWS ELB User Guide: WebSocket support
- URL: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/websockets-support.html

Confidence: MEDIUM (Authoritative vendor documentation; NGINX's specific WebSocket buffering guidance was not directly retrieved.)

### Finding 11: Scaling to 100k connections requires careful resource management

Claim: Operating at ~100,000 concurrent connections requires addressing three resource domains: (1) OS file descriptor limits (`ulimit -n`), (2) per-connection memory overhead in the runtime, and (3) cross-node message broadcasting.

Evidence:
- Linux kernel documentation confirms each open TCP socket consumes a file descriptor; per-process limit is managed via `ulimit -n`, system-wide via `fs.file-max` sysctl; `epoll(7)` is the scalable event notification facility used by high-concurrency servers.
- Per-connection memory overhead is runtime-specific (Go goroutine stacks, Node.js event loop buffers, JVM thread stacks); exact figures require empirical benchmarking.
- Multi-node broadcast requires a centralized message broker; Redis Pub/Sub is a common choice, with Kafka, NATS, or gRPC-based pubsub as alternatives (Redis docs, cloud-native reference architectures).

Sources:
- Linux epoll(7) & sysctl documentation
- URL: https://man7.org/linux/man-pages/man7/epoll.7.html
- Redis Pub/Sub documentation
- URL: https://redis.io/docs/manual/pubsub/

Confidence: MEDIUM (Architectural principle, but not independently benchmarked or measured in this research).
Corroborated By: RFC 6455 security considerations (Section 10.4) referencing "implementation-specific limits", and standard C10k/C10M engineering principles.

## Areas of Agreement

- **Directionality is the primary distinguishing factor**: WebSocket is full-duplex; SSE is unidirectional (server→client).
- **WebSocket requires a distinct protocol upgrade; SSE is standard HTTP**: WebSocket uses an `Upgrade` request, SSE uses a persistent `text/event-stream` HTTP response.
- **SSE offers built-in browser reconnection; WebSocket does not**: This is a spec-defined, browser-implemented feature of `EventSource`.
- **Binary vs text**: WebSocket supports both; SSE is UTF-8 text only.
- **HTTP/2**: WebSocket requires RFC 8441 to operate over HTTP/2; SSE is native HTTP/2.
- **Proxy complexity**: WebSocket is more prone to proxy/firewall issues; SSE benefits from standard HTTP tooling.

## Areas of Disagreement

No material disagreements were found among the authoritative sources consulted.

## Limitations

- **No empirical benchmark data**: The claim that 100,000 concurrent connections are achievable was accepted from the lab specification without independent benchmarking.
- **Proxy documentation was partially unavailable**: Specific NGINX WebSocket proxy guidance page (`websocket_proxy.html`) returned a 404; the general proxy module documentation was used instead.
- **Per-connection memory overhead is runtime-specific**: Exact byte-per-connection figures for Go goroutines, Node.js event loops, or JVM-based servers were not researched.
- **Browser-specific connection limits**: While MDN documents the 6-connection limit, this is a de facto browser policy rather than a web standard; actual values may vary by browser version.
- **Source priority**: The lab specification was used as a primary source for the SSE streaming example and the 100k scenario, but the underlying claims should not be verified using the lab text alone.

## Conclusion

For real-time use cases where the server only pushes data to the client—such as notifications, LLM token streaming, live dashboards, and stock tickers—Server-Sent Events provide the simplest, most interoperable solution. SSE avoids the non-HTTP upgrade mechanism, proxy complexity, and manual reconnection logic that WebSocket requires, while seamlessly leveraging HTTP/2 multiplexing.

WebSocket should remain the choice whenever bidirectional communication is necessary, such as chat applications, multiplayer games, collaborative editing, or any scenario where the client must send asynchronous messages of its own. The additional complexity of WebSocket (handshake management, reconnection, and proxying) is justified only when full-duplex communication is a hard requirement.

The 100,000 concurrent connection scenario is feasible with both protocols but is primarily constrained by OS-level file descriptor limits and backend runtime efficiency, not by the protocol choice. For pure server push at this scale, SSE's standard HTTP properties make it simpler to load-balance and scale horizontally.
