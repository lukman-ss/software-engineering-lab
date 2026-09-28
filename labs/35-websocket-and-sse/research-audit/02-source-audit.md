# Source Audit

Target Lab: `labs/35-websocket-and-sse`

---

## Source 1

Claimed Title: RFC 6455: The WebSocket Protocol
Claimed Publisher: Internet Engineering Task Force (IETF)
URL: https://datatracker.ietf.org/doc/html/rfc6455

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Normative specification for WebSocket protocol framing, handshake, masking, and closing procedures.

Assessment:
PASS

---

## Source 2

Claimed Title: RFC 8441: Bootstrapping WebSockets with HTTP/2
Claimed Publisher: Internet Engineering Task Force (IETF)
URL: https://datatracker.ietf.org/doc/html/rfc8441

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Normative specification defining Extended CONNECT method for running WebSocket over HTTP/2.

Assessment:
PASS

---

## Source 3

Claimed Title: HTML Living Standard: Section 9.2 Server-sent events
Claimed Publisher: Web Hypertext Application Technology Working Group (WHATWG)
URL: https://html.spec.whatwg.org/multipage/server-sent-events.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Definitive web standard specification for EventSource, text/event-stream syntax, reconnection algorithm, and Last-Event-ID processing.

Assessment:
PASS

---

## Source 4

Claimed Title: EventSource - Web APIs | MDN
Claimed Publisher: Mozilla Developer Network (MDN)
URL: https://developer.mozilla.org/en-US/docs/Web/API/EventSource

Reachable:
YES

Source Type:
SECONDARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Mentions 6-connection limit per browser+domain. While true for browser implementations (Chrome/Firefox), it is browser implementation policy rather than an IETF/WHATWG normative requirement.

Assessment:
PASS

---

## Source 5

Claimed Title: WebSocket - Web APIs | MDN
Claimed Publisher: Mozilla Developer Network (MDN)
URL: https://developer.mozilla.org/en-US/docs/Web/API/WebSocket

Reachable:
YES

Source Type:
SECONDARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Authoritative developer documentation describing the browser WebSocket interface and missing built-in auto-reconnect.

Assessment:
PASS

---

## Source 6

Claimed Title: RFC 7540: Hypertext Transfer Protocol Version 2 (HTTP/2)
Claimed Publisher: Internet Engineering Task Force (IETF)
URL: https://datatracker.ietf.org/doc/html/rfc7540

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Note: RFC 7540 was obsoleted by RFC 9113 in June 2022. While core stream multiplexing and connection-header prohibitions remain identical, citing RFC 9113 as the current standard is preferred.
- Typo in `05-report.md` line 138 references `https://datatracker.ietf.org/doc/html/rfc7541` (HPACK) under Source 6 instead of `rfc7540`.

Assessment:
WARNING

---

## Source 7

Claimed Title: Module ngx_http_proxy_module
Claimed Publisher: NGINX Docs
URL: https://nginx.org/en/docs/http/ngx_http_proxy_module.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Covers HTTP proxy buffering (`proxy_buffering`) and timeouts (`proxy_read_timeout`). Does not cover WebSocket-specific connection upgrade proxy directives (`proxy_set_header Upgrade $http_upgrade`), which are documented in NGINX websocket guide.

Assessment:
PASS

---

## Source 8 (Inline in Evidence 12/13 & Finding 11)

Claimed Title: Lab specification (Senior Software Engineer Lab #35)
Claimed Publisher: Internal Project
URL: N/A

Reachable:
NO (Self-referential / Unverifiable external artifact)

Source Type:
COMMUNITY / UNKNOWN

Relevant:
PARTIAL

Supports Claimed Topic:
PARTIAL

Problems:
- Circular evidence. Using lab prompt / internal spec as evidence for factual scaling claims (file descriptor consumption, memory overhead, Redis pub/sub necessity) without citing external systems literature (e.g. Linux socket man pages, POSIX file descriptor specs, or published benchmarks).

Assessment:
WARNING
