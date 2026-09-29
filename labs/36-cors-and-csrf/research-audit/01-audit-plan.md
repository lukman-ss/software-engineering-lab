# 01 — Audit Plan: Research Audit

## Target Lab
`labs/36-cors-and-csrf`

## Pipeline Scope & Override
- **Scope**: Research audit only (`research/` folder).
- **Excluded**: Code, tests, demo, engineering docs, implementation audit (handled in engineering audit stages).
- **Target Output Directory**: `labs/36-cors-and-csrf/research-audit/`

## Files Reviewed
1. `labs/36-cors-and-csrf/research/01-plan.md`
2. `labs/36-cors-and-csrf/research/02-sources.md`
3. `labs/36-cors-and-csrf/research/03-evidence.md`
4. `labs/36-cors-and-csrf/research/04-contradictions.md`
5. `labs/36-cors-and-csrf/research/05-report.md`
6. `labs/36-cors-and-csrf/research/06-open-questions.md`

## Major Claims To Verify
1. **SOP Semantics**: SOP permits cross-origin writes (e.g. form submissions) while blocking cross-origin reads.
2. **CORS Mechanism**: CORS is an HTTP-header mechanism enforced by browsers to control response readability; server executes the incoming request regardless of CORS validation.
3. **CORS Wildcard vs Credentials**: `Access-Control-Allow-Origin: *` cannot be combined with credentials (`credentials: include`); browsers block response access if wildcard is returned.
4. **Simple Requests**: GET, HEAD, POST with safelisted Content-Type bypass preflight and execute directly on target servers.
5. **Preflight Behavior**: Non-simple requests (PUT, DELETE, custom headers like `Authorization` or `X-CSRF-Token`) mandate an `OPTIONS` preflight fetch.
6. **SameSite Cookie Values**: `SameSite=Lax` blocks cross-site POST but permits top-level GET; `SameSite=Strict` blocks all cross-site requests; `SameSite=None` requires `Secure`.
7. **SameSite Level**: Operates at site level (eTLD+1), not origin level (scheme+host+port).
8. **CSRF Mitigation**: Naive Double-Submit is vulnerable to subdomain/plaintext injection; Signed Double-Submit (HMAC) or Synchronizer Token Pattern is required.
9. **Fetch Metadata**: `Sec-Fetch-Site` header provides robust cross-site detection directly at the HTTP layer.

## Primary Risks in Research
- Overgeneralization of browser defaults (e.g. Chrome Lax-by-default rollout nuances and 2-minute POST grace period).
- Treating CORS preflight side-effects as "CORS protecting against CSRF".
- Assumption that SameSite mitigates all CSRF without token-based defenses.
- Citation validity and accessibility for Tier 1 web standards (WHATWG, MDN, OWASP, PortSwigger, Google web.dev).

## Audit Strategy
1. **Source Audit**: Cross-verify every cited source in `02-sources.md` against real URLs, authoritative ownership, and scope.
2. **Claim Audit**: Map each factual statement in `03-evidence.md` and `05-report.md` to primary sources.
3. **Contradiction Verification**: Evaluate flagged contradictions in `04-contradictions.md` to ensure no subtle inaccuracies were overlooked.
4. **Gap Analysis**: Assess whether critical threat models or browser constraints were omitted.
5. **Verdict Generation**: Issue final verdict based on empirical evidence quality.
