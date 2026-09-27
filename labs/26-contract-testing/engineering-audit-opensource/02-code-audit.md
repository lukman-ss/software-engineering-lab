# Code Audit

Target Lab: `labs/26-contract-testing`

## Finding 1

Location: `internal/provider/server.go:24`, `49`, `95` — all three Provider handlers
Claimed Behavior: Provider handlers route `/v1/orders/{id}`; Dual provider also routes `/v2/orders/{id}`.
Observed Implementation: All three handlers use `strings.HasPrefix(r.URL.Path, "/v1/orders/")`. Dual adds `/v2/orders/`. Method guard rejects non-GET with 405. Path prefix check correctly trims prefix to extract ID.
Assessment: PASS
Severity: LOW
Notes: Path extraction via `strings.TrimPrefix` could over-match e.g. `/v1/ordersfoo` is guarded by exact prefix `/v1/orders/` (trailing slash), so it is actually safe. No false matches observed.

## Finding 2

Location: `internal/model/order.go:14-48` — DTOs
Claimed Behavior: V1 uses `IN_PROGRESS` (uppercase enum), `total` int64, `customer.name`; Breaking uses lowercase `in_progress`, `total` string, `customer.full_name`; V2 adds `currency`, keeps V1 enum lowercase but exposed on `/v2`.
Observed Implementation: Struct field tags and types match the claimed breaking mutations exactly. V2 omits `notes` (V1 extra field) — consistent with minimal subset consumer not requiring it.
Assessment: PASS
Severity: LOW
Notes: The design notes say V2 uses lowercase status `in_progress`; code uses `in_progress` on V2 — consistent. This is fine because V2 is a new endpoint with its own (new) consumer; the breaking concern only applies to V1 consumers.

## Finding 3

Location: `internal/consumer/client.go:82-111` — `GenerateMobileContract`
Claimed Behavior: Contract declares minimal subset {id, status=IN_PROGRESS, customer.name=Budi Santoso, total=150000}, excludes `notes`, Content-Type header asserted.
Observed Implementation: Interaction body contains exactly the four fields; `json.Number("150000")` used for total, preserving numeric type for diff engine.
Assessment: PASS
Severity: LOW
Notes: `json.Number("150000")` is a numeric literal form so `diffValues` treats it as a number path rather than a string. This is the root of type-matching correctness.

## Finding 4

Location: `internal/consumer/client.go:61-71` — `FetchOrder` runtime parsing
Claimed Behavior: Mobile consumer parses live response and fails fast on missing `customer.name` or unrecognized `status` enum.
Observed Implementation: `json.Unmarshal` into anonymous struct; explicit checks for empty `customer.name` and for `status ∉ {IN_PROGRESS, COMPLETED}`. Breaking provider (lowercase status, full_name) triggers errors.
Assessment: PASS
Severity: LOW
Notes: Runtime client-side validation independently catches breaking changes; tests assert `FetchOrder` returns error on breaking provider. Good defense-in-depth.

## Finding 5

Location: `internal/contract/verifier.go:58-115` — `Verify`
Claimed Behavior: Verifies each interaction: builds HTTP request, checks status code, reads body, decodes JSON with `UseNumber`, then recursively diffs declared fields.
Observed Implementation: Status mismatch recorded before body diff. `defer resp.Body.Close()` — wait: body is closed inline via `_ = resp.Body.Close()` after `io.ReadAll`, not deferred; this is acceptable since the function does not return between. Errors collected into `result.Errors`, `result.Passed` flipped to false on any failure.
Assessment: PASS
Severity: LOW
Notes: `defer resp.Body.Close()` is not used; instead closed immediately after read. Functionally fine. The `Verify` method is stateless per call (no mutable shared state), so concurrent reuse is safe.

## Finding 6

Location: `internal/contract/verifier.go:117-176` — `diffValues`
Claimed Behavior: Recursive subset comparison; reports: missing field, type mismatch (number vs non-number), type mismatch (other), value mismatch.
Observed Implementation: Handles nil expected (returns empty — "extra provider fields ignored"). Map branch recurses per key. Number branch: if either side is `json.Number` and the other is not → type mismatch; if both `json.Number` → string compare; if neither → falls through to reflect.TypeOf compare then value compare.

### Sub-point 6a
Location: `diffValues` number branch at `expected.(json.Number)` and `actual.(json.Number)`
Claimed Behavior: Breaking `total` (int 150000 as json.Number vs string "150000") → type mismatch.
Observed Implementation: Contract declares `json.Number("150000")`. Breaking provider emits JSON number `150000` for... wait — no. Breaking provider emits `"total":"150000"` (string). `json.NewDecoder` with `UseNumber` parses `"150000"` (quoted) as `string`, and unquoted `150000` as `json.Number`. Contract expected is `json.Number("150000")`. So expected is number, actual is string → `isExpNum||isActNum` true, and `!(isExpNum&&isActNum)` true → type mismatch emitted. Correct.
Assessment: PASS
Severity: LOW

### Sub-point 6b
Location: `diffValues` enum `status` comparison
Claimed Behavior: `status` expected `IN_PROGRESS` vs actual `in_progress` → value mismatch.
Observed Implementation: Both expected and actual are Go strings (`IN_PROGRESS`/`in_progress`). Neither is json.Number → reflect.TypeOf equal → falls to `expected != actual` → value mismatch. Correct.
Assessment: PASS
Severity: LOW

### Sub-point 6c
Location: `diffValues` missing `customer.name`
Claimed Behavior: Breaking provider emits `customer.full_name`, contract expects `customer.name` → "missing expected field".
Observed Implementation: Contract body nested `customer: {name: ...}`. Actual has only `full_name` key. `diffValues` recurses into `customer` map, key `name` not found in actual → diff "missing expected field 'customer.name'". Correct.
Assessment: PASS
Severity: LOW

## Finding 7

Location: `internal/contract/verifier.go` — `Verify` concurrent reuse
Claimed Behavior: Verifier instances shared across goroutines (test shares one `verifier` across 20 goroutines).
Observed Implementation: `Verifier` has only a `*http.Client` field (stateless methods). `http.Client` is documented as safe for concurrent use. `Verify` uses only local variables. No shared mutable state mutated.
Assessment: PASS
Severity: LOW
Notes: `httptest.Server` is also concurrency-safe; the same server handles concurrent requests fine.

## Finding 8

Location: `cmd/demo/main.go:14`
Claimed Behavior (implementation-notes line 13): "Interactive CLI demonstrating full lifecycle."
Observed Implementation: Program is fully non-interactive; runs all stages deterministically and exits 0 on success / 1 on failure. No stdin reads, no prompts.
Assessment: WARNING
Severity: LOW
Notes: Mislabeling "Interactive CLI". Not a functional defect. Documentation accuracy issue.

## Finding 9

Location: `cmd/demo/main.go` Stage 3 ordering
Claimed Behavior (engineering/03-execution-result.md line 105-107): Breaking diffs emitted as: (1) total type mismatch, (2) status value mismatch, (3) missing customer.name.
Observed Implementation: Actual execution emits: (1) status value mismatch, (2) missing customer.name, (3) total type mismatch.
Assessment: WARNING
Severity: LOW
Notes: `diffValues` iterates Go map `for key, expVal := range expMap` → nondeterministic ordering. Documented execution-result output does not match actual output ordering. The 3 breaking changes ARE all detected (count ≥ 3 asserted in test), so core behavior holds, but the documented example output is stale/incorrect. Real demo output verified via `go run ./cmd/demo`.

## Finding 10

Location: README.md:36-49
Claimed Behavior: README instructs `cd labs/26-contract-testing` then `go test -v ./...` / `go test -race ./...` / `go run ./cmd/demo`.
Observed Implementation: All commands verified runnable from repo root; README paths assume cwd = repo root. Commands match.
Assessment: PASS
Severity: LOW

## Summary

Code is functionally correct and matches core claims. Two documentation/labeling discrepancies: (F8) "Interactive CLI" mislabel; (F9) stale ordering in execution-result doc vs actual nondeterministic map iteration order. Neither is a functional defect. Concurrency path safe (F7). All three breaking-change detection paths verified (F6a–6c).
