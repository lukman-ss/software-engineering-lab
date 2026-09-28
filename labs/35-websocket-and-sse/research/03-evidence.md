# Research Evidence

## Evidence 1

Claim: WebSocket enables two-way (full-duplex) communication between a browser client and a server over a single TCP connection.
Evidence: "The WebSocket Protocol enables two-way communication between a client running untrusted code in a controlled environment to a remote host." RFC 6455 Section 1.1.
Source: RFC 6455 (2011)
URL: https://datatracker.ietf.org/doc/html/rfc6455
Confidence: HIGH
Corroborated By: RFC 7540, MDN WebSocket
Notes: WebSocket is distinct from HTTP/1.1's request-response model; it allows independent sending at both ends.

## Evidence 2

Claim: The WebSocket opening handshake is an HTTP/1.1 `Upgrade` request that transitions to the WebSocket protocol on the same TCP connection.
Evidence: "The opening handshake from the client looks as follows: GET /chat HTTP/1.1 ... Upgrade: websocket ... Connection: Upgrade."
Source: RFC 6455 Section 1.2
URL: https://datatracker.ietf.org/doc/html/rfc6455
Confidence: HIGH
Corroborated By: RFC 8441 (describes how HTTP/2 changes this requirement)
Notes: Because it uses connection-level headers, it cannot natively work over pure HTTP/2 without RFC 8441.

## Evidence 3

Claim: WebSocket data frames may be either UTF-8 text or binary payloads.
Evidence: "Broadly speaking, there are types for textual data (which is interpreted as UTF-8 [RFC3629] text), binary data (whose interpretation is left up to the application), and control frames." RFC 6455 Section 1.2.
Source: RFC 6455 (2011)
URL: https://datatracker.ietf.org/doc/html/rfc6455
Confidence: HIGH
Corroborated By: MDN WebSocket
Notes: WebSocket is a generic binary/text framing layer; it does not impose an application-level format.

## Evidence 4

Claim: WebSocket client browsers must implement manual reconnection and keep-alive (ping/pong) if they are desired.
Evidence: The spec defines Ping/Pong control frames (RFC 6455 Section 5.5.2) but requires application-level logic to track liveness and reconnect.
Source: RFC 6455 (2011)
URL: https://datatracker.ietf.org/doc/html/rfc6455
Confidence: HIGH
Corroborated By: MDN WebSocket ("WebSockets API has no way to apply backpressure", no mention of built-in reconnection)
Notes: The spec's Section 7.2.3 describes recovering from abnormal closure but leaves implementation choices to the developer.

## Evidence 5

Claim: Server-Sent Events (SSE) via `EventSource` are strictly unidirectional (server → client).
Evidence: "Unlike WebSockets, server-sent events are unidirectional; that is, data messages are delivered in one direction, from the server to the client."
Source: MDN EventSource
URL: https://developer.mozilla.org/en-US/docs/Web/API/EventSource
Confidence: HIGH
Corroborated By: WHATWG HTML Standard Section 9.2.1
Notes: If the client must send data, an alternate HTTP method (POST/fetch) must be used outside the SSE stream.

## Evidence 6

Claim: The `EventSource` constructor handles connection drops automatically by re-establishing the HTTP connection.
Evidence: "Clients will reconnect if the connection is closed." AND the fetch algorithm invokes "reestablish the connection" when encountering network errors.
Source: WHATWG HTML Standard Section 9.2.1 & Section 9.2.3
URL: https://html.spec.whatwg.org/multipage/server-sent-events.html
Confidence: HIGH
Corroborated By: MDN EventSource
Notes: The reconnection time is implementation-defined, but behavior is standardized.

## Evidence 7

Claim: SSE requires UTF-8 text encoding only, with no provision for binary frames.
Evidence: "Event streams are always decoded as UTF-8. There is no way to specify another character encoding."
Source: WHATWG HTML Standard Section 9.2.1
URL: https://html.spec.whatwg.org/multipage/server-sent-events.html
Confidence: HIGH
Corroborated By: MDN EventSource, RFC 6455 contrast
Notes: Binary data must be base64-encoded if sent over SSE.

## Evidence 8

Claim: Over HTTP/2, WebSockets require RFC 8441 (Extended CONNECT method) because HTTP/2 does not allow connection-wide headers like `Upgrade`.
Evidence: "Due to its multiplexing nature, HTTP/2 does not allow connection-wide header fields or status codes, such as the Upgrade and Connection request-header fields or the 101 (Switching Protocols) response code. These are all required by the [RFC6455] opening handshake."
Source: RFC 8441 Section 1
URL: https://datatracker.ietf.org/doc/html/rfc8441
Confidence: HIGH
Corroborated By: RFC 7540 Section 8.3, MDN
Notes: RFC 8441 tunnels the WebSocket handshake inside a single HTTP/2 stream.

## Evidence 9

Claim: Server-Sent Events use standard HTTP semantics, making them naturally compatible with HTTP/2 multiplexing without special extensions.
Evidence: SSE uses standard HTTP requests and the MIME type `text/event-stream`, so HTTP/2's multiplexed streams apply directly.
Source: WHATWG HTML Standard, MDN EventSource
URL: https://developer.mozilla.org/en-US/docs/Web/API/EventSource
Confidence: HIGH
Corroborated By: RFC 7540
Notes: No `Upgrade` is required; each SSE stream is just a persistent HTTP response body.

## Evidence 10

Claim: Browsers limit the number of simultaneous persistent HTTP connections to any given origin.
Evidence: "When not used over HTTP/2, SSE suffers from a limitation to the maximum number of open connections, which can be specially painful when opening various tabs as the limit is per browser and set to a very low number (6)."
Source: MDN EventSource
URL: https://developer.mozilla.org/en-US/docs/Web/API/EventSource
Confidence: HIGH
Corroborated By: RFC 7540 (HTTP/2 defaults to 100 concurrent streams per connection), common browser behavior documentation
Notes: Over HTTP/2, the default concurrent stream limit is negotiated (usually 100), mitigating the 6-connection cap.

## Evidence 11

Claim: Reverse proxies (such as nginx) often buffer HTTP responses by default, which can break or delay real-time event streaming if not configured.
Evidence: NGINX documentation explains `proxy_buffering` and `proxy_read_timeout` directives that affect streaming behavior.
Source: NGINX Docs: ngx_http_proxy_module
URL: https://nginx.org/en/docs/http/ngx_http_proxy_module.html
Confidence: HIGH
Corroborated By: Common architectural practices for SSE/WebSocket proxying
Notes: For SSE, buffering must typically be disabled or timeouts adjusted to prevent dropouts.

## Evidence 12

Claim: Scaling a WebSocket backend to tens of thousands of concurrent clients requires handling file descriptor limits (`ulimit -n`) and allocating memory per TCP socket.
Evidence: Each open TCP socket consumes a file descriptor at the OS level; Linux exposes the per-process limit via `ulimit -n` and the system-wide limit via the `fs.file-max` sysctl, while `epoll(7)` provides the scalable event-notification primitive used by high-connection-count servers (see also C10K/C10M engineering principles).
Source: Linux kernel documentation (epoll(7), sysctl)
URL: https://man7.org/linux/man-pages/man7/epoll.7.html
Confidence: MEDIUM (OS-level constraint is well-established; exact per-connection memory numbers vary by runtime and are not benchmarked here)
Corroborated By: C10K/C10M engineering literature
Notes: Specific memory-per-connection metrics were not found in the authoritative sources consulted; they are runtime-dependent.

## Evidence 13

Claim: For multi-node WebSocket backends, a central message broker (e.g., Redis Pub/Sub) is typically required to route events between nodes.
Evidence: Multi-node WebSocket architectures require a centralized message broker for inter-node event routing; Redis Pub/Sub is a common choice, with alternatives including Kafka, NATS, or gRPC-based pubsub (as documented in cloud-native reference architectures).
Source: Redis Pub/Sub documentation; cloud-native pattern references
URL: https://redis.io/docs/manual/pubsub/
Confidence: MEDIUM (common architectural pattern; specifics were not independently verified by external specs in this research session)
Corroborated By: Common cloud-native patterns
Notes: A production system may also use Kafka, NATS, or gRPC-based service discovery instead of Redis.
