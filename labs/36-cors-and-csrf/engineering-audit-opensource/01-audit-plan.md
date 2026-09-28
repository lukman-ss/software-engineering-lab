# Engineering Audit Plan

Target Lab: labs/36-cors-and-csrf
Implementation Files: None found
Tests: None found
Executable/Demo: None found
Approved Research Inputs: labs/36-cors-and-csrf/research/
Main Claims To Verify:
- CORS is not a backend firewall or CSRF protection mechanism
- Robust defense-in-depth measures against CSRF operate via HMAC-SHA256 signed session-bound tokens, Fetch Metadata, custom headers
Commands To Run:
- go test ./... (expect failure due to no Go files)
- go run ./cmd/demo (expect failure due to missing cmd/)
Primary Risks:
- Missing implementation entirely
- Missing test suite
- README describes non-existent code
- No verifiable artifact to audit