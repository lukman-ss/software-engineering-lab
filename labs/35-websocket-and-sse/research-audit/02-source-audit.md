# Source Audit: Lab 35 (WebSocket vs SSE Research)

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
- None. Normative specification defining frame format, HTTP/1.1 Upgrade handshake, masking, close handshake, and error handling.

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
- None. Authoritative standard for running WebSocket over HTTP/2 via Extended CONNECT method with `:protocol = websocket`.

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
- None. Living standard specifying the `EventSource` web interface, `text/event-stream` parser, reconnection algorithm, and `Last-Event-ID` mechanism.

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
- None. Accurately documents browser implementation details, including the 6-connection per origin limit under HTTP/1.1 and stream multiplexing over HTTP/2.

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
- None. Accurately documents browser `WebSocket` constructor, lack of native backpressure controls, and absence of standardized auto-reconnection in the browser interface.

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
- Minor: In `05-report.md` Finding 8, the citation URL is given as `https://datatracker.ietf.org/doc/html/rfc7541` (which is HPACK - Header Compression for HTTP/2) instead of `rfc7540`. In `02-sources.md`, the URL is correct (`rfc7540`). Also note RFC 7540 has been obsoleted by RFC 9113 (HTTP/2 core spec, June 2022), though the technical stream semantics cited remain identical.

Assessment:
PASS

---

## Source 7

Claimed Title: Module ngx_http_proxy_module
Claimed Publisher: NGINX Docs
URL: https://nginx.org/en/docs/http/ngx_http_proxy_module.html

Reachable:
YES

Source Type:
PRIMARY (Software Vendor Specification)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Contains standard definitions for `proxy_buffering` and `proxy_read_timeout`. The Research Agent properly noted in limitations that specific standalone WebSocket proxy guide page returned 404, so this general module doc was used.

Assessment:
PASS

---

## Additional Sources Cited in Report

### Source 8 (AWS ELB Documentation - Cited in 05-report.md Finding 10)

Claimed Title: AWS ELB User Guide: WebSocket support
Claimed Publisher: AWS Docs
URL: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/websockets-support.html

Reachable:
YES

Source Type:
PRIMARY (Vendor Documentation)

Relevant:
YES

Supports Claimed Topic:
PARTIAL

Problems:
- Missing from `02-sources.md`. Cited only in `05-report.md`.
- Claim mentions "AWS Elastic Load Balancing documentation confirms WebSocket (but not SSE) requires sticky session support". In Application Load Balancers, sticky sessions are typically needed if handshakes/reconnections must map to a specific target server maintaining state, but established WebSocket connections are pinned to the target instance for the life of the TCP connection once upgraded. SSE also requires sticky sessions if stateful server-side event tracking is pinned to a single server instance without a shared broker.

Assessment:
WARNING
