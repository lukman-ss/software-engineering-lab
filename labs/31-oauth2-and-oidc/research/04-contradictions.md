# Contradictions and Areas of Uncertainty

No material contradictions discovered among the primary Tier 1 sources on the core research questions.

All five primary sources (RFC 6749, OpenID Connect Core 1.0, RFC 7636, RFC 8252, RFC 9700) are consistent on:

- OAuth 2.0 is Authorization; OIDC is Authentication layered on OAuth 2.0
- Authorization Code Flow + PKCE is the mandatory modern standard
- Implicit Flow is deprecated/insecure
- ID Token validation must include signature, `iss`, `aud`, `exp`
- PKCE `S256` is mandatory-to-implement; `plain` is legacy/compatibility only
- Refresh token rotation is required for public clients

## Minor Nuances (Not Contradictions)

1. **Scope of PKCE requirement:**
   - RFC 7636 (2015): PKCE originally designed for public clients only.
   - RFC 9700 (2025): extends recommendation --- PKCE MUST for public clients, RECOMMENDED for confidential clients (as injection protection + CSRF protection).
   - ASSESSMENT: Evolution, not contradiction. Current best practice covers all client types.

2. **CSRF protection mechanism:**
   - RFC 6749: relies on `state` parameter.
   - RFC 9700 Sec 2.1: clients MAY rely on PKCE CSRF protection if AS supports PKCE; OIDC flows use `nonce`; otherwise one-time `state` bound to user agent MUST be used.
   - ASSESSMENT: Complementary mechanisms, layered defense. No conflict.

3. **Token storage specifics:**
   - Primary RFCs do not explicitly mandate "HttpOnly Secure SameSite cookie vs BFF" wording; this is derived from threat-model reasoning (XSS exposure of `localStorage`, token in URL fragment exposure documented in RFC 6819 Sec 10.3 / RFC 9700 Sec 2.1.2 / oauth.net implicit page).
   - ASSESSMENT: Interpretation from evidence, not direct normative quote. Confidence remains HIGH due to converging secondary guidance, but exact phrasing attributed as interpretation.
