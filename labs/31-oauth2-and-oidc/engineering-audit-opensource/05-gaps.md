# Gaps & Open Issues

## Gap list

### G1 — `plain` PKCE downgrade surface (DOC_CODE_MISMATCH)
- README emphasizes S256; `pkg/pkce` + `pkg/server` still accept `plain`.
- RFC 9700 / RFC 8252 §7.6 deprecate `plain` for native apps.
- Impact: a client could downgrade PKCE to a no-op; server permits it.
- Severity: LOW

### G2 — HS256 signing (implementation note, not a bug)
- OIDC production uses RS256/JWKS. HS256 with shared secret used here.
- Acceptable for a reference/demo impl if communicated. Currently undocumented as a limitation.
- Severity: LOW

### G3 — `aud` modeled as single string
- `IDTokenClaims.Audience string` — OIDC allows `[]string`.
- Scope: single-audience reference impl. Fine for this lab.
- Severity: LOW

### G4 — Error wrapping non-unwrappable for inner errors
- `fmt.Errorf("%w: %v", sentinel, inner)` — inner cause not in `errors.Is`/`As` chain.
- Severity: LOW

### G5 — Missing boundary tests
- PKCE 43/128 length bounds; auth-code expiry; refresh-token expiry; iat +300s skew edge.
- Severity: LOW

### G6 — No JWKS / no asymmetric key path
- README doesn't claim JWKS. `oidc` package is self-contained HMAC. Aligned.
- Severity: LOW (informational)

## Open questions (carried from research/06-open-questions.md)
- Should `plain` be removed or gated behind a strict mode? → Recommend removing or deprecating.
- Is HS256 signing explicitly dev-only? → Should add a note.

## No fake benchmarks / results
- `engineering/03-execution-result.md` and demo output both regenerated live — verified real.
- Severity: none.

## Risk rollup
| Severity | Count |
|---|---|
| CRITICAL | 0 |
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 6 |
