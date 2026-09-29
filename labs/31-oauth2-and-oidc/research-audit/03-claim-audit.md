# Claim Audit

Target Lab: `labs/31-oauth2-and-oidc`

---

## Claim 1

Claim: OAuth 2.0 is an authorization (delegated access) framework, not an authentication protocol.

Location: `research/03-evidence.md` Evidence 1; `research/05-report.md` Finding 1

Evidence Provided: RFC 6749 Abstract: "enables a third-party application to obtain limited access to an HTTP service..."; OIDC Core Sec 1: "without profiling OAuth 2.0, it is incapable of providing information about the authentication of an End-User."

Source: RFC 6749 Abstract/Sec 1; OIDC Core 1.0 Sec 1

Source Actually Supports Claim:
YES — Both sources verified, quotes are accurate.

Classification:
FACT

Severity:
LOW (fully supported)

Notes:
- OAuth 2.1 draft also explicitly echoes: "OAuth is an authorization protocol, not an authentication protocol."

---

## Claim 2

Claim: Access tokens are typically opaque to the client.

Location: `research/03-evidence.md` Evidence 2

Evidence Provided: RFC 6749 Sec 1.4: "The string is usually opaque to the client."

Source: RFC 6749 Sec 1.4

Source Actually Supports Claim:
YES — Quote is accurate.

Classification:
FACT

Severity:
LOW (fully supported)

Notes:
- Research also correctly references RFC 9068 (JWT Profile for OAuth 2.0 Access Tokens) for structured access tokens. RFC 9068 was not cited as a formal numbered source in `02-sources.md` but is referenced inline in Evidence 2. This is a MEDIUM gap — RFC 9068 is a significant modern standard that should be listed as a formal source if cited.

---

## Claim 3

Claim: OIDC is an identity layer on top of OAuth 2.0; authentication is signaled by `openid` scope and returned as an ID Token JWT; profile via UserInfo Endpoint.

Location: `research/03-evidence.md` Evidence 3; `research/05-report.md` Finding 2

Evidence Provided: OIDC Core Abstract/Sec 1/Sec 1.3/Sec 2/Sec 5.3 verbatim text.

Source: OIDC Core 1.0

Source Actually Supports Claim:
YES — Verified. OIDC Core Sec 1 explicitly states "a simple identity layer on top of the OAuth 2.0 protocol."

Classification:
FACT

Severity:
LOW (fully supported)

---

## Claim 4

Claim: ID Tokens MUST contain `iss`, `sub`, `aud` (MUST contain RP client_id), `exp`, `iat`.

Location: `research/03-evidence.md` Evidence 4; `research/05-report.md` Finding 2

Evidence Provided: OIDC Core Sec 2 verbatim claim table.

Source: OIDC Core 1.0 Sec 2

Source Actually Supports Claim:
YES — Verified. OIDC Core Sec 2 lists `iss`, `sub`, `aud`, `exp`, `iat` as REQUIRED. `aud` MUST contain the OAuth 2.0 `client_id` of the Relying Party.

Classification:
FACT

Severity:
LOW

---

## Claim 5

Claim: RP MUST validate ID Token via 13-step check including issuer match, audience validation, signature, expiry, and nonce.

Location: `research/03-evidence.md` Evidence 5; `research/05-report.md` Finding 3

Evidence Provided: OIDC Core Sec 3.1.3.7 with section-specific reference.

Source: OIDC Core 1.0 Sec 3.1.3.7

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
- Research accurately notes nonce is part of the validation. The OIDC Core Sec 2 confirms nonce must equal the value sent in the authentication request when present.

---

## Claim 6

Claim: Using OAuth access tokens for login is insecure because access tokens carry no standardized authentication semantics or audience binding to the RP.

Location: `research/03-evidence.md` Evidence 6; `research/05-report.md` Finding 1

Evidence Provided: OIDC Core Sec 1 (incapable of providing authentication information); RFC 8725 Sec 2.7/2.8/3.9/3.12 (substitution/cross-JWT attacks, audience validation).

Source: OIDC Core 1.0 Sec 1; RFC 8725

Source Actually Supports Claim:
PARTIAL

Classification:
INTERPRETATION

Severity:
LOW

Notes:
- Claim is valid; supported by the cited sources. However, it is an inferential composite claim, not a direct normative prohibition.
- RFC 9700 Sec 2.3 also corroborates (access tokens SHOULD be audience-restricted to resource servers, not RPs), verified.

---

## Claim 7

Claim: PKCE prevents authorization-code interception: client creates `code_verifier` (43-128 unreserved chars, >=256-bit entropy), sends `code_challenge` (S256 = `BASE64URL(SHA256(verifier))`); server verifies verifier at token endpoint.

Location: `research/03-evidence.md` Evidence 7

Evidence Provided: RFC 7636 Sec 1.1/4.1/4.2/4.6 formulas; worked example in Appendix B.

Source: RFC 7636

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
- RFC 7636 Appendix B example values (verifier `dBjftJeZ...`, challenge `E9Melhoa...`) are correctly cited.

---

## Claim 8

Claim: S256 is Mandatory-To-Implement (MTI); `plain` MUST NOT be used in new implementations.

Location: `research/03-evidence.md` Evidence 8

Evidence Provided: RFC 7636 Sec 4.2: "If the client is capable of using S256, it MUST use S256, as S256 is Mandatory To Implement on the server." Sec 7.2: "plain SHOULD NOT be used."

Source: RFC 7636 Sec 4.2, 7.2

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
- Minor precision issue: The evidence states `plain SHOULD NOT be used`; the research claim says `must not be used in new implementations`. RFC 7636 Sec 7.2 says "plain SHOULD NOT be used in new implementations, unless they cannot support S256 for some technical reason." This is an accurate characterization, not an overstatement.

---

## Claim 9

Claim: Implicit flow is deprecated; clients SHOULD NOT use `response_type=token`; should use `response_type=code` + PKCE instead.

Location: `research/03-evidence.md` Evidence 9; `research/05-report.md` Finding 5

Evidence Provided: RFC 9700 Sec 2.1.2: "clients SHOULD NOT use the implicit grant (response type token)..."; OAuth 2.1 draft: "The Implicit grant is omitted."

Source: RFC 9700 Sec 2.1.2; draft-ietf-oauth-v2-1

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
- Research also correctly acknowledges OIDC Core 1.0 still documents Implicit/Hybrid flows historically (Contradiction 1).

---

## Claim 10

Claim: Refresh tokens for public clients MUST be sender-constrained or use refresh token rotation; rotation means server issues new refresh token with every response and invalidates prior token; reuse signals breach.

Location: `research/03-evidence.md` Evidence 10; `research/05-report.md` Finding 7

Evidence Provided: RFC 9700 Sec 2.2.2 (MUST), Sec 4.14.2 (rotation semantics and breach detection).

Source: RFC 9700 Sec 2.2.2, 4.14.2

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
- RFC 9700 Sec 2.2.2 verified: "Refresh tokens for public clients MUST be sender-constrained or use refresh token rotation." Evidence accurately quoted.

---

## Claim 11

Claim: JWT validation must pin the algorithm (never trust `alg` header blindly), reject `none` unless explicitly configured, validate all nested crypto operations, and use mutually-exclusive validation rules per JWT kind.

Location: `research/03-evidence.md` Evidence 11; `research/05-report.md` Finding 3

Evidence Provided: RFC 8725 Sec 3.1 (algorithm pinning), Sec 3.3 (validate all crypto), Sec 3.11/3.12 (explicit typing, mutually-exclusive rules).

Source: RFC 8725

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
- RFC 8725 verified. Sec 3.1 confirmed: "Libraries MUST enable the caller to specify a supported set of algorithms and MUST NOT use any other algorithms." Sec 3.2 confirmed: "JWT libraries SHOULD NOT generate JWTs using none unless explicitly requested."

---

## Claim 12

Claim: Tokens must not live in browser JS-reachable storage; recommended pattern is BFF: confidential client holds tokens server-side in cookie session and proxies resource requests.

Location: `research/03-evidence.md` Evidence 12; `research/05-report.md` Finding 6

Evidence Provided: draft-ietf-oauth-browser-based-apps-27 Sec 6.1 BFF architecture description; Sec 5 malicious JS threat model.

Source: draft-ietf-oauth-browser-based-apps-27 Sec 5, 6.1

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
- BFF core responsibilities verified in draft-27 Sec 6.1: "(1) interacts with AS as a confidential OAuth client; (2) manages OAuth access and refresh tokens avoiding direct exposure to the browser-based application; (3) forwards all requests to a resource server, augmenting them with the correct access token." Fully supports the research claim.
