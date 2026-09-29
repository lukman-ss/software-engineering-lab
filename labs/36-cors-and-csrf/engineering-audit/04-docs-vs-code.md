# Docs vs Code Audit

## README Alignment

File: `README.md`
- Claim: Architecture lists `internal/cors`, `internal/csrf`, `internal/bank`, `cmd/demo`, and `tests`.
  - Observed: All 5 directories and modules exist and match specified roles.
- Claim: Commands `go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo`.
  - Observed: All 3 commands execute without failures or warnings.

## Engineering Notes Alignment

File: `engineering/01-design.md`, `engineering/02-implementation-notes.md`, `engineering/03-execution-result.md`
- Claim: Execution outputs match observed execution of `cmd/demo/main.go`.
  - Observed:
    - Initial: Victim $1000, Attacker $50
    - After vulnerable attack: Victim $600, Attacker $450
    - After protected attack: Blocked 403, Victim $600, Attacker $450
    - After legit transfer: Victim $500, Attacker $550
  - Output matches exact recorded output in `engineering/03-execution-result.md`.

## Research Claims Alignment

File: `research/05-report.md`
- Claim: CORS is an SOP opt-in relaxation mechanism for reads/preflights, NOT a backend authorization barrier.
  - Code & Demo Proof: Proved in `cmd/demo/main.go` and `TestIntegration_CORS_Does_Not_Prevent_CSRF_Execution`.
- Claim: Anti-CSRF requires token-based or origin/header-based validation.
  - Code & Demo Proof: Implemented via HMAC-SHA256 session-bound tokens, `Sec-Fetch-Site`, and custom headers.

## Discrepancies Found

None. Documentation accurately depicts codebase structure, commands, behavior, and output.
