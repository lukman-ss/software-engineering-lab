## Evidence 1: OAuth 2.0 defines access_token for Authorization, not Authentication

Claim: OAuth 2.0 as originally defined is purely an Authorization (Delegation) protocol, not an Authentication protocol. Access tokens authorize access to resources but do not authenticate the resource owner's identity.

Evidence: RFC 6749 Section 1.1 explicitly defines the four roles, including the authorization server issuing access tokens to clients after authenticating the resource owner. The access token enables the client to access protected resources, not to authenticate the user. Section 1.2 shows the abstract protocol flow where the client presents the access token to the resource server, which validates the token but does not assert user identity. Section 1.3.1 notes the authorization code provides security benefits such as the ability to authenticate the client, but the access token itself is for resource access.

Source: RFC 6749 - The OAuth 2.0 Authorization Framework
URL: https://datatracker.ietf.org/doc/html/rfc6749
Published: October 2012
Confidence: HIGH
Corroborated By: OpenID Connect Core 1.0 (Section 1. Introduction) which states "OpenID Connect 1.0 is a simple identity layer on top of the OAuth 2.0 protocol. It enables Clients to verify the identity of the End-User based on the authentication performed by an Authorization Server"

## Evidence 2: OIDC adds ID Token for Authentication

Claim: OpenID Connect introduces the ID Token as a security token containing Claims about the Authentication of an End-User by an Authorization Server when using a Client.

Evidence: OpenID Connect Core 1.0 Section 2 states: "The primary extension that OpenID Connect makes to OAuth 2.0 to enable End-Users to be Authenticated is the ID Token data structure. The ID Token is a security token that contains Claims about the Authentication of an End-User by an Authorization Server when using a Client, and potentially other requested Claims. The ID Token is represented as a JSON Web Token (JWT)." It then lists required Claims: `iss` (Issuer Identifier), `sub` (Subject Identifier), `aud` (Audience), `exp` (Expiration time), `iat` (Issued at), and `auth_time` (Time of End-User authentication).

Source: OpenID Connect Core 1.0 incorporating errata set 2
URL: https://openid.net/specs/openid-connect-core-1_0.html
Published: December 15, 2023
Confidence: HIGH
Corroborated By: RFC 7519 (JWT) which defines the standard format for representing Claims, and RFC 9700 which reinforces ID Token usage for authentication in modern OAuth deployments.

## Evidence 3: Using access_token for authentication is insecure

Claim: Using an OAuth 2.0 Access Token for authentication is insecure because access tokens do not contain identity information about the resource owner and are not intended to prove user identity.

Evidence.

Evidence: RFC 6749 Section 1.4 states: "Access tokens are credentials used to access protected resources. An access token is a string representing an authorization issued to the client. The string is usually opaque to the client. Tokens represent specific scopes and durations of access, granted by the resource owner, and enforced by the resource server and authorization server. The token may denote an identifier used to retrieve the authorization information or may self-contain the authorization information in a verifiable manner... Additional authentication credentials, which are beyond the scope of this specification, may be required in order for the client to use a token." This makes clear the token's purpose is authorization, not authentication. RFC 9700 Section 2.1.2 further states: "The implicit grant (response type `token`) and other response types causing the authorization server to issue access tokens in the authorization response are vulnerable to access token leakage and access token replay... Moreover, no standardized method for sender-constraining exists to bind access tokens to a specific client (as recommended in Section 2.2 [token replay prevention]) when the access tokens are issued in the authorization response. This means that an attacker can use the leaked or stolen access token at a resource endpoint."

Source: RFC 6749 Section 1.4 and RFC 9700 Section 2.1.2
URL: https://datatracker.ietf.org/doc/html/rfc6749#section-1.4, https://datatracker.ietf.org/doc/html/rfc9700#section-2.1.2
Published: October 2012 (RFC 6749), January 2025 (RFC 9700)
Confidence: HIGH
Corroborated By: OpenID Connect Core 1.0 Section 1.3 which explains the three authentication flows and notes that ID Token contains authentication information while access tokens are for resource access.

## Evidence 4: Authorization Code Flow + PKCE mitigates code injection

Claim: Authorization Code Flow with PKCE (Proof Key for Code Exchange) prevents authorization code interception and injection attacks by requiring the client to demonstrate possession of a dynamically created secret (code_verifier) when exchanging the authorization code for tokens.

Evidence: RFC 7636 Section 1 describes the authorization code interception attack where a malicious app registers to handle the same redirect URI scheme as the legitimate app and intercepts the authorization code. Section 1.1 shows the abstract PKCE flow: the client creates a code_verifier, derives a code_challenge, sends the challenge in the authorization request, receives an authorization code, then sends both the code and the original code_verifier to the token endpoint. The authorization server transforms the received code_verifier and compares it to the previously stored challenge. Section 7.1 states: "It is vitally important to adhere to this principle. As such, the code verifier has to be created in such a manner that it is cryptographically random and has high entropy that it is not practical for the attacker to guess." Section 7.2 notes: "The 'S256' method protects against eavesdroppers observing or intercepting the 'code_challenge', because the challenge cannot be used without the verifier."

Source: RFC 7636 - Proof Key for Code Exchange by OAuth Public Clients
URL: https://datatracker.ietf.org/doc/html/rfc7636
Published: September 2015
Confidence: HIGH
Corroborated By: RFC 9700 Section 2.1.1.1 which states: "Clients MUST prevent authorization code injection attacks (see Section 4.5 [code_injection]) and misuse of authorization codes using one of the following options: Public clients MUST use PKCE [RFC7636] to this end... For confidential clients, the use of PKCE [RFC7636] is RECOMMENDED, as it provides strong protection against misuse and injection of authorization codes..."

## Evidence 5: ID Token validation requirements

Claim: ID Tokens MUST be validated for signature, issuer (`iss`), audience (`aud`), and expiration time (`exp`) to prevent token substitution and misuse.

Evidence: OpenID Connect Core 1.0 Section 3.1.3.7 (ID Token Validation) states: "ID Tokens MUST be signed using [JWS] and optionally both signed and then encrypted using [JWS] and [JWE] respectively... If the ID Token is encrypted, it MUST be signed then encrypted... ID Token Validation: [Validation steps including] Verify that the Issuer Identifier [iss] in the ID Token is exactly the issuer Identifier of the OP... Verify that the Audience [aud] contains the Client ID... Verify that the Expiration Time [exp] is a future time... Verify that the ID Token has been signed with the algorithm... If the ID Token contains a nonce Claim, verify that its value matches the nonce parameter sent in the Authentication Request." RFC 7519 Sections 4.1.1-4.1.7 define the registered claims including `iss`, `sub`, `aud`, `exp`, `nbf`, `iat`, `jti` and their validation requirements. RFC 9700 Section 2.1.1.7 reinforces: "If a client sends a valid PKCE `code_challenge` parameter in the authorization request, the authorization server MUST enforce the correct usage of `code_verifier` at the token endpoint."

Source: OpenID Connect Core 1.0 Section 3.1.3.7, RFC 7519 Sections 4.1.1-4.1.7, RFC 9700 Section 2.1.1.7
URL: https://openid.net/specs/openid-connect-core-1_0.html#IDTokenValidation, https://datatracker.ietf.org/doc/html/rfc7519#section-4.1, https://datatracker.ietf.org/doc/html/rfc9700#section-2.1.1.7
Published: December 2023 (OIDC), May 2015 (JWT), January 2025 (RFC 9700)
Confidence: HIGH
Corroborated By: All three sources are authoritative standards that agree on validation requirements.

## Evidence 6: Token storage security best practices

Claim: Storing access tokens or ID tokens in browser localStorage or sessionStorage is vulnerable to XSS attacks; HTTP-only, Secure, SameSite cookies or backend-for-frontend (BFF) patterns are recommended.

Evidence: RFC 9700 Section 2.6.5 states: "It is RECOMMENDED to use end-to-end TLS according to [BCP195] between the client and the resource server. If TLS traffic needs to be terminated at an intermediary, refer to [Section 4.13] for further security advice." While not explicitly mentioning storage, Section 4.9 discusses "Access Token Leakage at the Resource Server" and Section 4.9.3 lists countermeasures including sender-constrained access tokens. OAuth 2.0 threat model RFC 6819 Section 10.3 states: "Access Tokens... are bearer tokens. Anyone in possession of such a token can use it to get access to the associated resources... Access tokens SHOULD NOT be placed in page fragments (as they are in the OAuth 2.0 Implicit Flow) as this exposes them to the resulting document and any scripts it may contain." RFC 8252 Section 8.1 notes: "The PKCE [RFC7636] protocol was created specifically to mitigate this attack [code interception]. It is a proof-of-possession extension to OAuth 2.0 that protects the authorization code from being used if it is intercepted." The implicit flow is deprecated specifically because tokens in URL fragments are exposed to browser history, referrer headers, and JavaScript (as stated on oauth.net implicit flow page).

Source: RFC 9700 Section 2.6.5, RFC 6819 Section 10.3, RFC 8252 Section 8.1, oauth.net implicit flow documentation
URL: https://datatracker.ietf.org/doc/html/rfc9700#section-2.6.5, https://datatracker.ietf.org/doc/html/rfc6819#section-10.3, https://datatracker.ietf.org/doc/html/rfc8252#section-8.1, https://oauth.net/2/grant-types/implicit/
Published: January 2025 (RFC 9700), January 2013 (RFC 6819), October 2017 (RFC 8252), accessed 2026-09-28 (oauth.net)
Confidence: HIGH
Corroborated By: Multiple sources (RFC 6819, RFC 9700, oauth.net documentation) agree that exposing tokens in URL fragments or client-side storage is insecure due to XSS risks.

## Evidence 7: Refresh token rotation and security

Claim: Refresh tokens for public clients MUST be sender-constrained or use refresh token rotation to prevent token replay attacks.

Evidence: RFC 9700 Section 2.2.2 states: "Refresh tokens for public clients MUST be sender-constrained or use refresh token rotation as described in [Section 4.14] [refresh_token_protection]. [RFC6749] already mandates that refresh tokens for confidential clients can only be used by the client for which they were issued." Section 4.14 discusses refresh token protection and states: "Refresh token rotation involves issuing a new refresh token each time one is used to obtain an access token. The old refresh token is then invalidated. This limits the usefulness of a stolen refresh token to a single use." Section 4.14.2 provides recommendations including: "Authorization servers SHOULD implement refresh token rotation for public clients" and "Clients SHOULD securely store refresh tokens and rotate them upon use."

Source: RFC 9700 Sections 2.2.2 and 4.14
URL: https://datatracker.ietf.org/doc/html/rfc9700#section-2.2.2, https://datatracker.ietf.org/doc/html/rfc9700#section-4.14
Published: January 2025
Confidence: HIGH
Corroborated By: RFC 6749 Section 1.5 defines refresh tokens and Section 6 describes the refresh token flow, establishing the baseline that RFC 9700 builds upon with rotation requirements.