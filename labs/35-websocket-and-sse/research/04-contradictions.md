# Contradictions

No material contradictions discovered between the authoritative sources consulted.

The following areas were cross-checked and found consistent:

1. **Protocol directionality**: RFC 6455 (bidirectional), RFC 8441 (tunneling), and WHATWG HTML (unidirectional server→client) all describe the fundamental directionality of each protocol consistently.

2. **WebSocket handshake**: RFC 6455's HTTP/1.1 upgrade mechanism and RFC 8441's HTTP/2 extension are consistent. RFC 8441 explicitly notes that the traditional upgrade headers cannot be used on HTTP/2.

3. **SSE reconnection**: WHATWG HTML standard's "reestablish the connection" algorithm matches MDN EventSource's description of automatic reconnection.

4. **Payload types**: RFC 6455 (binary and UTF-8 text) and WHATWG HTML (UTF-8 only) are consistent.

5. **HTTP/2 limits**: MDN EventSource's mention of the 6-connection HTTP/1.1 limit and 100-stream HTTP/2 default aligns with RFC 7540's stream multiplexing model.

No conflicting claims were identified across the primary sources.