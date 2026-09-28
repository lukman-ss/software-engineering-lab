# Gap Analysis — labs/31-oauth2-and-oidc

Gap types considered:
MISSING_TEST, BROKEN_IMPLEMENTATION, DOC_CODE_MISMATCH, RACE_CONDITION,
UNHANDLED_ERROR, MISSING_EDGE_CASE, IMPLEMENTATION_OVERCLAIM,
RESEARCH_MISMATCH, FAKE_DEMO, FAKE_BENCHMARK, UNVERIFIED_RESULT.

Findings:

- MISSING_TEST: No test for `ErrIssuedInFuture` (iat > now+5min) in oidc_test.go. (LOW)
- MISSING_TEST: No test for expired authorization code, access token, refresh token via time manipulation (could mock time). (LOW)
- MISSING_EDGE_CASE: ID Token with missing `azp` (optional) claim not exercised; code supports optional via struct omitempty. (LOW)
- DOC_CODE_MISMATCH: README omits mention of symmetric (HS256) justification vs asymmetric JWKS; noted in engineering notes. (LOW)
- No BROKEN_IMPLEMENTATION, RACE_CONDITION, UNHANDLED_ERROR, IMPLEMENTATION_OVERCLAIM, RESEARCH_MISMATCH (research not audited), FAKE_DEMO, FAKE_BENCHMARK, UNVERIFIED_RESULT.

All gaps are LOW severity; none block approval.