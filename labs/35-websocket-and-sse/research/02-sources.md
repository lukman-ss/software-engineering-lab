# Research Sources

## Source 1

Title: RFC 6455: The WebSocket Protocol
Publisher: Internet Engineering Task Force (IETF)
URL: https://datatracker.ietf.org/doc/html/rfc6455
Published: December 2011
Accessed: 2026-09-28
Source Tier: Tier 1 (Official Standards Track Specification)
Relevance: Defines the base WebSocket protocol, frame types, HTTP/1.1 upgrade handshake, masking, connection closing, and security model.

## Source 2

Title: RFC 8441: Bootstrapping WebSockets with HTTP/2
Publisher: Internet Engineering Task Force (IETF)
URL: https://datatracker.ietf.org/doc/html/rfc8441
Published: September 2018
Accessed: 2026-09-28
Source Tier: Tier 1 (Official Standards Track Specification)
Relevance: Defines Extended CONNECT method for running RFC 6455 over an HTTP/2 stream, removing the requirement for hop-by-hop upgrade headers and enabling multiplexing.

## Source 3

Title: HTML Living Standard: Section 9.2 Server-sent events
Publisher: Web Hypertext Application Technology Working Group (WHATWG)
URL: https://html.spec.whatwg.org/multipage/server-sent-events.html
Published: Living Standard (Updated September 2026)
Accessed: 2026-09-28
Source Tier: Tier 1 (Official Web Standard)
Relevance: Normative specification for the `EventSource` API, `text/event-stream` format, automatic reconnection algorithm, and `Last-Event-ID` header.

## Source 4

Title: EventSource - Web APIs | MDN
Publisher: Mozilla Developer Network (MDN)
URL: https://developer.mozilla.org/en-US/docs/Web/API/EventSource
Published: Last modified March 13, 2025
Accessed: 2026-09-28
Source Tier: Tier 1 / Tier 2 (Authoritative Technical Documentation)
Relevance: Practical documentation on `EventSource`, noting the 6-connection per-domain limit over HTTP/1.1 and default multiplexing up to 100 over HTTP/2.

## Source 5

Title: WebSocket - Web APIs | MDN
Publisher: Mozilla Developer Network (MDN)
URL: https://developer.mozilla.org/en-US/docs/Web/API/WebSocket
Published: Last modified September 25, 2024
Accessed: 2026-09-28
Source Tier: Tier 1 / Tier 2 (Authoritative Technical Documentation)
Relevance: API details, absence of built-in backpressure in the base `WebSocket` interface, lack of native auto-reconnect, and events.

## Source 6

Title: RFC 7540: Hypertext Transfer Protocol Version 2 (HTTP/2)
Publisher: Internet Engineering Task Force (IETF)
URL: https://datatracker.ietf.org/doc/html/rfc7540
Published: May 2015 (Obsoleted by RFC 9113 in June 2022)
Accessed: 2026-09-28
Source Tier: Tier 1 (Official Standards Track Specification)
Relevance: Details HTTP/2 framing, streams, multiplexing, and prohibition of connection-specific headers like `Upgrade`. Note: Technical multiplexing rules are preserved in RFC 9113.

## Source 7

Title: epoll(7) - Linux Manual Page & sysctl kernel documentation
Publisher: Linux Kernel Organization / Michael Kerrisk (man7.org)
URL: https://man7.org/linux/man-pages/man7/epoll.7.html
Published: Ongoing
Accessed: 2026-09-28
Source Tier: Tier 1 (Official OS / Kernel Documentation)
Relevance: Describes Linux event notification facility (`epoll`) and kernel tuning parameters (`fs.file-max`, `net.ipv4.tcp_rmem`, `net.ipv4.tcp_wmem`) required for handling large numbers (~100,000) of concurrent open file descriptors / TCP sockets.

## Source 8

Title: Module ngx_http_proxy_module
Publisher: NGINX Docs
URL: https://nginx.org/en/docs/http/ngx_http_proxy_module.html
Published: Ongoing
Accessed: 2026-09-28
Source Tier: Tier 1 (Official Software Documentation)
Relevance: Details on `proxy_buffering`, `proxy_read_timeout`, and handling long-lived streaming connections behind reverse proxies.
