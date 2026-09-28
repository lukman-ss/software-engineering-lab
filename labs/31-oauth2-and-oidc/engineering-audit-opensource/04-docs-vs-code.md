# Docs vs Code Audit

## README Claims
- Describes PKCE S256, OIDC ID token, refresh token rotation, demo steps – matches implementation.
- Lists commands to run tests and demo – accurate.

## Engineering Notes
- Design doc states token family revocation on replay – implementation matches.
- Claims validation list matches oidc.go implementation.
- Concurrency safety claimed – mutexes and race test confirm.

## Observed Mismatches
- README mentions `refresh token rotation with family revocation` – demo and code confirm.
- No mismatch detected.

## Gap Types Detected
- DOC_CODE_MISMATCH: none
- TEST_CLAIM_MISMATCH: none
- RESEARCH_IMPLEMENTATION_MISMATCH: none (research not audited per pipeline)
