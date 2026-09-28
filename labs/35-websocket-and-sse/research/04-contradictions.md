# Contradictions

No material contradictions discovered between the authoritative sources consulted.

The following areas were cross-checked and found consistent:

1. **Protocol directionality**: RFC 6455 (bidirectional), RFC 8441 (tunneling), and WHATWG HTML (unidirectional server→client) all describe the fundamental directionality of each protocol consistently.

2. **WebSocket handshake**: RFC 6455's HTTP/1.1 upgrade mechanism and RFC 8441's HTTP/2 extension are consistent. RFC 8441 explicitly notes that the traditional upgrade headers cannot be used on HTTP/2.

3. **SSE reconnection**: WHATWG HTML standard's "reestablish the connection" algorithm matches MDN EventSource's description of automatic reconnection.

4. **Payload types**: RFC 6455 (binary and UTF-8 text) and WHATWG HTML (UTF-8 only) are consistent.

5. **HTTP/2 limits**: MDN EventSource's mention of the 6-connection HTTP/1.1 limit and 100-stream HTTP/2 default aligns with RFC 7540 / RFC 9113's stream multiplexing model.

6. **RFC 7540 currency**: The original research cited RFC 7540 (May 2015) without noting that it was obsoleted by RFC 9113 (June 2022). This is an editorial gap rather than a contradiction; the technical content on stream multiplexing and connection-header prohibitions remains unchanged in RFC 9113. Source 6 in `02-sources.md` has been updated to note the obsolescence, and Finding 7 / Finding 8 now cite RFC 9113 as the current standard.

No conflicting claims were identified across the primary sources.