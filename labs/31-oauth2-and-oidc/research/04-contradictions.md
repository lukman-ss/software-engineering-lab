# Contradictions

## Contradiction 1
**Topic**: OIDC Core still documents Implicit Flow and Hybrid Flow, while RFC 9700 and OAuth 2.1 deprecate/omit Implicit.

**Source A (OIDC Core 1.0 Sec 3.2/3.3)**:
Defines "Implicit Flow" (response_type=id_token token / id_token) and "Hybrid Flow" (response_type=code id_token / code token / code id_token token) as valid authentication flows. Flow comparison table lists them alongside Authorization Code Flow.

**Source B (RFC 9700 Sec 2.1.2 / OAuth 2.1)**:
"Clients SHOULD NOT use the implicit grant (response type token) ... Clients SHOULD instead use the response type code". OAuth 2.1: "The Implicit grant (response_type=token) is omitted from this specification."

**Assessment**:
OIDC Core 1.0 (final 2014, errata 2023) predates RFC 9700 (Jan 2025) and OAuth 2.1 consolidation. The Implicit Flow is historically part of OIDC but is now considered insecure for browser-based apps because it returns access tokens in the URL fragment, exposing them to leakage via referer headers, browser history, and XSS. RFC 9700 is the authoritative current BCP; OIDC Core flows should be read in historical context. Implementers should use Authorization Code Flow + PKCE (response_type=code) with openid scope. The Hybrid Flow similarly exposes tokens in the front channel and is discouraged by the same logic, though not explicitly named in RFC 9700 Sec 2.1.2.

## Contradiction 2
**Topic**: Refresh token rotation mandate.

**Source A (OIDC Core 1.0 Sec 12.2)**:
Example successful refresh response contains a new `refresh_token`: `{"access_token": "...", "token_type": "Bearer", "refresh_token": "9yNOxJtZa5", "expires_in": 3600}`. The text does not use the term "rotation" and does not require revocation of the old token.

**Source B (RFC 9700 Sec 4.14.2)**:
Explicitly mandates rotation: "the authorization server issues a new refresh token with every access token refresh response. The previous refresh token is invalidated... If a refresh token is compromised and subsequently used by both the attacker and the legitimate client, one of them will present an invalidated refresh token, which will inform the authorization server of the breach... it will revoke the active refresh token."

**Assessment**:
OIDC Core example is compatible with rotation (new refresh_token issued) but does not specify the security requirement of invalidating the old one. RFC 9700 is the authoritative security BCP and imposes the rotation requirement for public clients. The OIDC example should be read as a structural example, not a security prescription. No material contradiction — RFC 9700 supersedes on security hardening.

## Contradiction 3
**Topic**: OIDC Core allows `none` algorithm for ID Tokens under limited conditions; RFC 8725 forbids `none` except when JWT is protected by other means.

**Source A (OIDC Core Sec 2)**:
"ID Tokens MUST NOT use none as the alg value unless the Response Type used returns no ID Token from the Authorization Endpoint (such as when using the Authorization Code Flow) and the Client explicitly requested the use of none at Registration time."

**Source B (RFC 8725 Sec 3.2)**:
"Applications MUST only allow the use of cryptographically current algorithms... if a JWT is cryptographically protected end-to-end by a transport layer... there may be no need to apply another layer... the use of the 'none' algorithm can be perfectly acceptable. The 'none' algorithm should only be used when the JWT is cryptographically protected by other means."

**Assessment**:
Both align: `none` is only acceptable when integrity/authenticity is guaranteed by another layer (TLS + registration-time consent in OIDC; transport-layer crypto in RFC 8725). No contradiction, but RFC 8725 is stricter about algorithm pinning in general (Sec 3.1: libraries MUST reject `none` unless explicitly requested). Implementers should follow RFC 8725 defaults and only allow `none` when explicitly configured.

## Contradiction 4
**Topic**: Browser-based apps token storage — localStorage vs httpOnly cookies vs BFF.

**Source A (RFC 9700 Sec 4.2/4.3)**:
Credentials leak via Referer headers and browser history; tokens in URL fragments can be logged by browsers/proxies.

**Source B (Browser-based-apps draft Sec 8.5/6.1)**:
Recommends BFF as most secure; for browser-only clients, in-memory storage is preferred over localStorage/IndexedDB; persistent storage discouraged. Does NOT forbid localStorage but says "applications MUST NOT use persistent token storage (e.g., localStorage) unless the tokens are sender-constrained or encrypted."

**Assessment**:
RFC 9700 identifies leakage vectors; draft provides concrete architecture patterns. The draft's allowance of in-memory for short-lived tokens (with sender-constrained fallback) is the pragmatic current guidance. The "trap" in the lab ("Menyimpan Access Token di localStorage") is validated by both sources as a real risk.