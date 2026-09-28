# Contradiction Audit: OAuth 2.0 & OIDC Research

Target Lab: `labs/31-oauth2-and-oidc`
Audit Scope: Research Files and Sources Consistency

---

## Contradiction Audit Result

No material contradictions found across the research artifacts (`01-plan.md`, `02-sources.md`, `03-evidence.md`, `04-contradictions.md`, `05-report.md`, `06-open-questions.md`) or between cited primary sources (RFC 6749, OpenID Connect Core 1.0, RFC 7636, RFC 8252, RFC 7519, RFC 9700).

---

## Evaluated Nuances

1. **Evolution of PKCE Scope:**
   - Statement A: RFC 7636 (2015) specifies PKCE for public clients.
   - Statement B: RFC 9700 (2025) mandates PKCE for public clients and recommends it for confidential clients.
   - Assessment: Consistent. Standard evolution over time.

2. **CSRF vs PKCE Binding:**
   - Statement A: RFC 6749 specifies `state` parameter for CSRF protection.
   - Statement B: RFC 9700 allows PKCE to fulfill code injection / CSRF protection if supported, while OIDC uses `nonce`.
   - Assessment: Complementary layered defense mechanisms.

3. **Normative vs Derived Guidance on Token Storage:**
   - Statement A: RFC 6819 / RFC 9700 focus on bearer token leakage risks in front-channel / JavaScript contexts.
   - Statement B: Research report recommends HttpOnly cookies / BFF pattern for browser apps.
   - Assessment: Appropriately qualified in `04-contradictions.md` and `05-report.md` as architectural interpretation of standard threat models.

---

## Conclusion
No internal or source-level contradictions identified.
