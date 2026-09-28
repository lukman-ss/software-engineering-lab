# Docs vs Code Audit

## Documentation Verification

| Documented Item | Source Code / Test / Executable | Match Status | Notes |
|---|---|---|---|
| Architecture layout | `README.md:5-12` vs directory structure | MATCH | `internal/cors`, `internal/csrf`, `internal/bank`, `cmd/demo`, `tests` directories exist and match descriptions. |
| Test Commands | `README.md:17-25` (`go test -v ./...`, `go test -race ./...`) | MATCH | Both commands run cleanly and pass 100%. |
| Demo Command | `README.md:31-33` (`go run ./cmd/demo`) | MATCH | Executes deterministically, demonstrating vulnerable vs protected flows. |
| Spec Compliance | `internal/cors/middleware.go` vs CORS spec | MATCH | Disallows wildcard credentials, returns proper Vary and preflight response headers. |
| Token Mechanism | `internal/csrf/token.go` vs research notes | MATCH | HMAC-SHA256 session-bound token implementation directly aligns with research requirements. |

No `DOC_CODE_MISMATCH`, `TEST_CLAIM_MISMATCH`, or `RESEARCH_IMPLEMENTATION_MISMATCH` identified.
