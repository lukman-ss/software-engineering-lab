# Code Audit

Target Lab: labs/26-contract-testing

## Finding 1

Location: internal/contract/verifier.go:58-115 (`Verify`)
Claimed Behavior: Provider verification engine verifies HTTP interactions (method, path, headers, status, response schema/types/enums) against generated contracts.
Observed Implementation: Verifies method/path via live request, status code, body subset via `diffValues`. Response `Headers` in contract (`Content-Type: application/json`) never read or compared. Request headers are set, response headers ignored.
Assessment: WARNING
Severity: MEDIUM
Notes: Core body/status verification correct. Header comparison claimed in engineering/01-design.md Components §4 but absent in code. No test covers headers.

## Finding 2

Location: internal/contract/verifier.go:117-176 (`diffValues`)
Claimed Behavior: Strict validation on declared fields: type safety, exact enum case, nested field existence.
Observed Implementation: Map-subset recursion correct. `json.Number` vs `string` distinguished via `UseNumber` decoder. Exact string compare catches `IN_PROGRESS` vs `in_progress`. Missing key reported. Extra provider fields (`notes`) ignored. Verified live: breaking provider yields 3 diffs (status value, missing `customer.name`, total type).
Assessment: PASS
Severity: LOW
Notes: Minimal-subset rule correctly implemented. Error order nondeterministic (Go map iteration) — tests correctly assert count, not order.

## Finding 3

Location: internal/contract/verifier.go:164-175 (primitive fallback)
Claimed Behavior: N/A (generic comparison)
Observed Implementation: After map/number handling, falls to `reflect.TypeOf` check then `expected != actual`. If a contract body ever contains an array (`[]interface{}`), `expected != actual` panics (uncomparable type). Current contract has no arrays so untriggered.
Assessment: WARNING
Severity: LOW
Notes: Out of lab scope (ponytail documents minimal runner, no arrays used). No test exercises array values. Would need recursive slice handling for general use.

## Finding 4

Location: internal/contract/verifier.go:47-55, internal/consumer/client.go:26-31 (`http.Client{}` with no timeout)
Claimed Behavior: N/A
Observed Implementation: Both verifier and consumer use `&http.Client{}` with no `Timeout`. Hung provider blocks `Verify`/`FetchOrder` indefinitely.
Assessment: WARNING
Severity: LOW
Notes: Irrelevant under `httptest` scope; no hang observed. Production/CI gate would want timeout. Not claimed anywhere.

## Finding 5

Location: internal/consumer/client.go:34-79 (`FetchOrder`)
Claimed Behavior: Mobile client maps provider payload to consumer view; fails on contract violation.
Observed Implementation: Unmarshals into strict struct (`Total int64`, `Customer.Name`, status enum check, empty-name check). Breaking provider (total as string) fails at `json.Unmarshal` with `cannot unmarshal string into ... int64` — before reaching name/status checks. V1/Dual parse correctly.
Assessment: PASS
Severity: LOW
Notes: Client failure on breaking provider proven by test. Specific branches `customer.name is missing` / `unknown status` unreachable when total-type breaks first — covered only generically (see test audit).

## Finding 6

Location: internal/provider/server.go (ProviderV1, ProviderBreaking, ProviderDual)
Claimed Behavior: V1 fulfills contract; Breaking introduces 3 breaks; Dual preserves V1 alongside V2.
Observed Implementation: All three stateless `ServeHTTP`, GET-only (405 otherwise), 404 otherwise. V1 returns `IN_PROGRESS`/`name`/`150000 int`. Breaking returns `in_progress`/`full_name`/`"150000" string`. Dual serves identical V1 on `/v1/...` plus V2 on `/v2/...`. No shared state; concurrency-safe.
Assessment: PASS
Severity: LOW
Notes: Empty-ID edge (`/v1/orders/` → 200 with empty id) unvalidated but out of scope. `/v2` handler exists but no contract/test exercises it (see gaps).

## Finding 7

Location: internal/model/order.go, cmd/demo/main.go
Claimed Behavior: V1/Breaking/V2 DTOs; 4-stage demo (generate, V1 pass, breaking block, dual pass).
Observed Implementation: DTO tags correct. Demo runs all stages against `httptest` servers, exits 1 on unexpected outcome. Live run reproduced claimed output semantically (error order varies — map nondeterminism, expected).
Assessment: PASS
Severity: LOW
Notes: Demo is real, not scripted. No fake output.
