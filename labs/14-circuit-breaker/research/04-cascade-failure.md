# 04 Cascade Failure

## Failure Propagation Mechanics
1. Dependency degrades or slows down.
2. Incoming callers block waiting for responses or timeouts.
3. Thread pools, HTTP connection pools, and memory saturate.
4. Upstream callers (and health checks) fail to obtain threads/sockets.
5. The caller crashes or becomes unresponsive to its own clients.
6. The failure cascades upstream transitively to the edge/gateway.
