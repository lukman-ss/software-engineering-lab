# Code Audit

## Finding 1

Location: internal/contract/verifier.go:95-111
Claimed Behavior: Provider V1 verification passes; Breaking provider fails with >=3 diffs.
Observed Implementation: `decoder.UseNumber()` decodes actual body; `diffValues` compares expected `json.Number` to actual `json.Number` for `total`, strings for `status`/`customer.name`. For V1: all match → PASS. For Breaking: status string mismatch (1), missing `customer.name` (1), total type mismatch json.Number vs string (1) = 3 errors → FAIL.
Assessment: PASS
Severity: LOW
Notes: Map iteration is non-deterministic; error ordering varies run-to-run but count is stable at 3. Logic correct.

## Finding 2

Location: internal/consumer/client.go:61-78
Claimed Behavior: Consumer parses provider JSON, enforces status enum, requires customer.name.
Observed Implementation: Unmarshals into anonymous struct, rejects empty `customer.name`, rejects status not in {IN_PROGRESS, COMPLETED}.
Assessment: PASS
Severity: LOW
Notes: Aligns with consumer contract expectations; correctly fails against Breaking provider (`in_progress` status + missing `customer.name`).

## Finding 3

Location: internal/provider/server.go:18-44, 53-80, 89-131
Claimed Behavior: ProviderV1 returns IN_PROGRESS/Customer.Name(int64 total); ProviderBreaking returns in_progress/full_name/string total; ProviderDual returns V1 at /v1 and V2 at /v2.
Observed Implementation: Three ServeHTTP implementations with explicit JSON responses via json.NewEncoder.
Assessment: PASS
Severity: LOW
Notes: No DB/external I/O; deterministic. ProviderDual /v1 path satisfies contract, /v2 path adds currency.

## Finding 4

Location: internal/contract/verifier.go:47-55, Verify method
Claimed Behavior: Concurrent contract verification safe under -race.
Observed Implementation: Verifier holds `*http.Client` (concurrency-safe). Each Verify call builds local request, local response decode, no shared mutable state across interactions.
Assessment: PASS
Severity: LOW
Notes: Race detector passed; http.Client is documented safe for concurrent use.

## Finding 5

Location: internal/contract/verifier.go:87-92
Claimed Behavior: Response body read failure handled gracefully.
Observed Implementation: `io.ReadAll` result checked; `resp.Body.Close()` called before decode. Error appended to result, `continue` to next interaction.
Assessment: PASS
Severity: LOW
Notes: Body always closed; errors propagated into VerificationResult.Errors without crashing.

## Finding 6

Location: cmd/demo/main.go:20
Claimed Behavior: Consumer generates contract; demo prints it.
Observed Implementation: `json.MarshalIndent(c, ...)` with ignored error (`_`). Contract struct is in-memory and always serializable; no error path reachable.
Assessment: WARNING
Severity: LOW
Notes: Ignored error acceptable for non-failing marshal of known-good struct; not a real risk but deviates from strict error-handling conventions.

## Finding 7

Location: internal/model/order.go
Claimed Behavior: V1/Breaking/V2 DTOs carry breaking-change semantics.
Observed Implementation: Field tags and comments mark breaking semantics: V1 uses `name`/`IN_PROGRESS`/`int64`; Breaking uses `full_name`/`in_progress`/`string`; V2 uses `full_name`/`in_progress`/`string`+`currency`.
Assessment: PASS
Severity: LOW
Notes: V2 shares Breaking's field naming (full_name) but only exposed via /v2 endpoint; V1 preserved separately. Safe evolution.