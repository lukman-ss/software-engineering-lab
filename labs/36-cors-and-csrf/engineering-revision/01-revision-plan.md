# Engineering Revision Plan

Target Lab: `labs/36-cors-and-csrf`
Previous Verdict: APPROVED — zero blocking issues; 4 LOW non-blocking test-coverage gaps.

## Blocking Issues
None.

## Non-Blocking Issues

| ID     | Severity | Description |
|--------|----------|-------------|
| GAP-01 | LOW | `/api/transfer/custom-header` has no integration test verifying attacker POST (missing `X-Requested-With`) receives 403. |
| GAP-02 | LOW | CSRF token submitted via `X-CSRF-Token` header not exercised in integration tier (only form-field path covered). |
| GAP-03 | LOW | Cross-session token binding (token-A rejected with session-B cookie) not tested at integration tier. |
| GAP-04 | LOW | `Sec-Fetch-Site: same-origin` positive flow not exercised in integration tests. |

## Files To Change

- `tests/integration_test.go` — add 4 tests, one per gap.

## Tests To Add/Modify

1. `TestIntegration_CustomHeader_BlocksAttacker` (GAP-01)
2. `TestIntegration_CSRF_Header_Submission` (GAP-02)
3. `TestIntegration_CrossSession_TokenRejected` (GAP-03)
4. `TestIntegration_SecFetchSite_SameOrigin_Allowed` (GAP-04)

## Validation Commands

```bash
cd labs/36-cors-and-csrf
go test ./...
go test -race ./...
go run ./cmd/demo
```
