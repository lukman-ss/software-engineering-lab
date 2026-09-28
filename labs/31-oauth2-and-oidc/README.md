# Lab 31: OAuth 2.0 & OpenID Connect (OIDC)

A reference implementation demonstrating authentication vs authorization, PKCE code interception protection, ID Token validation, and refresh token rotation with family revocation.

## Overview

- **OAuth 2.0 (Authorization)**: Access Tokens represent delegated authorization grants to protected resource servers.
- **OIDC (Authentication)**: ID Tokens (signed JWTs) provide cryptographically verifiable claims about user authentication (`iss`, `sub`, `aud`, `exp`, `nonce`).
- **PKCE (RFC 7636 / RFC 9700)**: Proof Key for Code Exchange (`S256`) prevents authorization code interception attacks.
- **Refresh Token Rotation (RFC 9700 Section 4.14)**: Single-use refresh tokens with lineage tracking that revoke the entire token family upon reuse detection.

## Structure

- `pkg/pkce`: PKCE code verifier and S256 challenge generation and verification.
- `pkg/oidc`: ID Token creation, HMAC-SHA256 signature verification, and standard claims validation.
- `pkg/server`: In-memory Authorization Server supporting Auth Code grant, PKCE, OIDC ID Tokens, and Refresh Token Rotation.
- `pkg/client`: Client helper orchestrating requests and token parsing.
- `cmd/demo`: Executable walkthrough of legitimate flows and attack defenses.
- `tests`: Unit tests and race detector tests.

## Running Tests and Demo

Run all tests:
```bash
go test -v ./...
```

Run race detector:
```bash
go test -race ./...
```

Run demo:
```bash
go run ./cmd/demo
```
