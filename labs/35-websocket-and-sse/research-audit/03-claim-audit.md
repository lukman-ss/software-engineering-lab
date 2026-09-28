# Claim Audit

## Claim 1

Claim: WebSocket is a full-duplex, bidirectional communication channel over a single TCP connection that operates outside of the standard HTTP/1.1 request-response lifecycle.
Location: `research/05-report.md:Finding 1`
Evidence Provided: Direct quotes from RFC 6455 Section 1.1 and 1.2 on Upgrade handshake and bidirectional frames.
Source: RFC 6455
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Core protocol behavior accurately described.

---

## Claim 2

Claim: The base WebSocket protocol defines distinct frame types for binary data and text (UTF-8) data.
Location: `research/05-report.md:Finding 2`
Evidence Provided: RFC 6455 Section 1.2 quote regarding textual and binary data types.
Source: RFC 6455
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Well-supported.

---

## Claim 3

Claim: WebSocket protocol and browser API do not mandate automatic reconnection; client must implement reconnect logic.
Location: `research/05-report.md:Finding 3`
Evidence Provided: RFC 6455 Section 7.2.3 and MDN WebSocket API documentation.
Source: RFC 6455, MDN WebSocket
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Accurate negative claim substantiated by specification absence and MDN documentation.

---

## Claim 4

Claim: Server-Sent Events deliver data only from the server to the client; there is no native mechanism for client-to-server messaging over the SSE stream.
Location: `research/05-report.md:Finding 4`
Evidence Provided: MDN EventSource quote and WHATWG HTML Living Standard Section 9.2 MIME type definition.
Source: MDN EventSource, WHATWG HTML Living Standard
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Accurate protocol directionality constraint.

---

## Claim 5

Claim: Browser-managed `EventSource` object automatically reconnects on dropped connections, supports configurable reconnection time, and transmits the last event ID to the server on reconnection.
Location: `research/05-report.md:Finding 5`
Evidence Provided: WHATWG HTML Living Standard Section 9.2.3 and Section 9.2.4 normative steps.
Source: WHATWG HTML Living Standard
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Fully supported by authoritative specification.

---

## Claim 6

Claim: Server-Sent Events are restricted to UTF-8 text; the spec provides no mechanism for binary payloads.
Location: `research/05-report.md:Finding 6`
Evidence Provided: WHATWG HTML Standard Section 9.2.1 quotation ("Event streams are always decoded as UTF-8").
Source: WHATWG HTML Living Standard
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Direct normative quote.

---

## Claim 7

Claim: WebSocket over HTTP/2 requires RFC 8441 (Extended CONNECT) using `:protocol = websocket` pseudo-header because HTTP/2 forbids connection-wide headers like `Upgrade`.
Location: `research/05-report.md:Finding 7`
Evidence Provided: RFC 8441 Section 1 and Section 5 citations.
Source: RFC 8441, RFC 9113
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Authoritative protocol requirement accurately cited.

---

## Claim 8

Claim: SSE requires no protocol-level extension for HTTP/2 because it operates as standard persistent HTTP responses with `text/event-stream`.
Location: `research/05-report.md:Finding 8`
Evidence Provided: WHATWG HTML Standard Section 9.2 and RFC 9113 multiplexing semantics.
Source: WHATWG HTML Standard, RFC 9113
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Fully valid.

---

## Claim 9

Claim: Browsers historically enforce an HTTP/1.1 connection limit of ~6 per domain across all tabs, limiting SSE under HTTP/1.1, whereas HTTP/2 defaults to 100 concurrent streams.
Location: `research/05-report.md:Finding 9`
Evidence Provided: MDN EventSource documentation and RFC 7540 / RFC 9113 stream concurrency defaults.
Source: MDN EventSource, RFC 7540 / RFC 9113
Source Actually Supports Claim: YES
Classification: IMPLEMENTATION-SPECIFIC
Severity: MEDIUM
Notes: The 6-connection cap is a browser client policy rather than an RFC standard rule. The report properly notes this in Limitations.

---

## Claim 10

Claim: Reverse proxies buffer responses by default; administrators must tune `proxy_buffering` and `proxy_read_timeout`.
Location: `research/05-report.md:Finding 10`
Evidence Provided: NGINX documentation and AWS ELB documentation references.
Source: NGINX Docs, AWS ELB Docs
Source Actually Supports Claim: PARTIAL
Classification: IMPLEMENTATION-SPECIFIC
Severity: MEDIUM
Notes: Proxy buffering defaults depend on specific reverse proxy software and configuration. Correctly marked as MEDIUM confidence in report.

---

## Claim 11

Claim: Operating at ~100,000 concurrent connections requires addressing OS file descriptors (`ulimit -n`, `fs.file-max`), per-connection runtime memory overhead, and message broadcasting (Redis Pub/Sub).
Location: `research/05-report.md:Finding 11`
Evidence Provided: Linux epoll(7) man page and Redis Pub/Sub docs.
Source: Linux kernel docs, Redis docs
Source Actually Supports Claim: PARTIAL
Classification: INTERPRETATION
Severity: MEDIUM
Notes: General architectural consensus, but lacks empirical benchmarks or measured memory footprints for specific runtimes. The report transparently flags this in its Limitations section.
