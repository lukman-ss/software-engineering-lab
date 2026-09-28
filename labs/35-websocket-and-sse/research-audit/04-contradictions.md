# Contradiction Audit: Lab 35 (WebSocket vs SSE Research)

## Contradiction 1: Load Balancer Sticky Session Requirement
Statement A:
"AWS Elastic Load Balancing documentation confirms WebSocket (but not SSE) requires sticky session support..."
Location: `05-report.md: Finding 10`

Statement B:
"...SSE's standard HTTP framing simplifies load balancing and avoids sticky-session requirements for read-only traffic."
Location: `05-report.md: Executive Summary`

Statement C (Reality / Architecture Nuance):
Sticky sessions on AWS ALB route initial requests based on a cookie. Once a WebSocket connection finishes the 101 Switching Protocols handshake, the TCP connection remains pinned to that target instance until closed; it does not require sticky sessions merely to function as a WebSocket. Furthermore, if an application using SSE tracks client state locally in server memory without a central broker, SSE client reconnections will also fail or produce inconsistent stream state unless sticky sessions are configured or the server architecture is stateless/broker-backed.

Type:
SOURCE_CONFLICT / OVERGENERALIZATION

Impact:
Misleading architectural guidance regarding load balancer configuration for WebSocket vs SSE. Both protocols require sticky sessions if state is tied to node memory upon reconnection, and neither requires cookie stickiness during the lifetime of an established open stream/connection.

Assessment:
MEDIUM severity. The claim should specify that sticky sessions are an application state architecture concern (handling reconnections to stateful nodes) rather than an inherent protocol requirement of WebSocket.

---

## Overall Assessment
No material contradictions were found between RFC standards (RFC 6455, RFC 8441, RFC 7540) and the WHATWG HTML standard. Protocol-level mechanics, framing, encoding, and reconnection behaviors are described consistently across all research documents.
