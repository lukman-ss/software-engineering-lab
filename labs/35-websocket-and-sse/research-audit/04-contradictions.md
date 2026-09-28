# Contradiction Audit

Target Lab: `labs/35-websocket-and-sse`

---

## Contradiction 1: RFC 7540 vs RFC 7541 Citation Typo

Statement A:
"RFC 7540 Sections 5 and 8.1 (Streams and Multiplexing, HTTP request/response exchange)"

Location:
`05-report.md` (Finding 8, line 137)

Statement B:
"URL: https://datatracker.ietf.org/doc/html/rfc7541"

Location:
`05-report.md` (Finding 8, line 138)

Type:
INTERNAL

Impact:
RFC 7541 is HPACK (header compression), whereas RFC 7540 is the HTTP/2 core protocol specification. The URL typo points readers to HPACK instead of the stream specification.

Assessment:
LOW. Editorial typo; does not alter substantive findings.

---

## Contradiction 2: RFC 7540 Currency Status

Statement A:
"RFC 7540: Hypertext Transfer Protocol Version 2 (HTTP/2)... Published: May 2015"

Location:
`02-sources.md` (Source 6)

Statement B:
RFC 7540 was officially obsoleted by RFC 9113 (HTTP/2) in June 2022.

Location:
IETF Standards Track Registry

Type:
SOURCE_CONFLICT

Impact:
Minor standard tracking omission. The technical multiplexing and header rules relevant to WebSocket and SSE remain unchanged in RFC 9113.

Assessment:
LOW.

---

## Summary
No material contradictions found across the core technical claims regarding protocol framing, directionality, reconnection mechanics, or HTTP/2 transport semantics.
