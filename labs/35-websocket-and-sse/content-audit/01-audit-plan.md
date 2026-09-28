# Content Audit Plan

Target: labs/35-websocket-and-sse/content/
Mode: Audit content only, do not audit research/code, do not modify files

Scope:
- 01-content-brief.md
- 02-master-draft.md
- 03-code-snippets.md
- 04-diagrams.md
- 05-key-takeaways.md
- 06-source-map.md

References:
- research/05-report.md (11 findings, 8 sources)
- research-audit/07-verdict.md (APPROVED)
- engineering/01-design.md, 02-implementation-notes.md, 03-execution-result.md
- engineering-audit/02-code-audit.md, 04-docs-vs-code.md, 05-gaps.md, 06-verdict.md (APPROVED, 3 warnings)
- engineering-audit-opensource/* (same)
- Code: internal/sse/sse.go, internal/ws/ws.go, internal/server/server.go, tests/protocol_test.go, cmd/demo/main.go, go.mod, README.md

Checks:
1. Accuracy vs research (RFC 6455, WHATWG SSE, RFC 8441/9113, proxy, scaling)
2. Accuracy vs engineering implementation (handshake, framing, hub, endpoints, tests)
3. Clarity, formatting, completeness
4. Hallucinated facts / platform bias
5. Code snippet line-number and content fidelity
6. Diagram fidelity
7. Source-map traceability
