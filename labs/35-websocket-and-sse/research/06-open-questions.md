# Open Questions

1. **Exact memory overhead per WebSocket connection in Go (goroutine stack) vs Node.js (event loop)**: The lab specification mentions memory overhead but does not provide quantified figures.

2. **Empirical 100k-connection benchmark data**: No independent benchmark confirming the resource footprint of 100,000 concurrent connections (either protocol) was located.

3. **SSE behavior under HTTP/2 with reverse proxies**: Whether specific reverse proxy implementations (nginx, Envoy, AWS ALB) buffer SSE streams differently from standard HTTP responses was not definitively tested.

4. **SSE over HTTP/3 (QUIC)**: HTTP/3 adoption may change connection-limit dynamics. Research needed. Noted as a future alternative to the HTTP/1.1 / HTTP/2 regime discussed here; WebTransport and QUIC-based protocols may further shift the trade-off landscape.

5. **SSE as a replacement for LLM streaming at scale**: While OpenAI and Anthropic use SSE for token streaming, the operational limits under extreme load (100k+ streaming connections) are not fully documented.

6. **Browser behavior on SSE reconnection**: The exact backoff algorithm and maximum retry limits are implementation-defined ("probably a few seconds"). Specific vendor behavior across Chrome, Firefox, and Safari is not standardized.

7. **WebSocket vs SSE for mobile data usage**: While WHATWG notes SSE can reduce battery via push proxies, a formal comparison of data overhead between the two protocols on mobile networks was not researched.

8. **Redis Pub/Sub vs alternatives for multi-node broadcasting**: The lab suggests Redis Pub/Sub, but no comparison with Kafka, NATS, or gRPC-based pub/sub was performed.
