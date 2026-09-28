# Code Audit: Contract Testing (lab-26)

## Test 1: verifier.Verify flow
- **Expected:** iterates contract interactions; performs HTTP request; captures status/JSON mismatches; returns `Passed: false` when differences. Fatal bug: verification passes when provider behavior diverges.
- **Observed result:** PASS
- **Implementation:**
  - `internal/contract/verifier.go:58-115` : builds explicit `http.NewRequest` per interaction, propagates custom headers, compares `resp.StatusCode` against expected, reads/decode body with `json.Number` preservation, records errors per mismatched interaction
  - `internal/contract/verifier.go:117-176` : recursive subset matcher `diffValues`; flags missing nested fields, primitive type differences, and exact value divergences
  - Runs demonstrated expected behavior: `srvV1 ✅ Passed`, `srvBreaking ❌ Blocked`, `srvDual ✅ Passed`
- **Notes:** body leak protected (`_ = resp.Body.Close()` post-read); response header declared in contract but not enforced in `Verify`; `ProviderState` stored but no state-setup hook; no request timeout.

## Test 2: provider V1 compliance
- **Expected:** `ProviderV1` serves `status=IN_PROGRESS`, `customer.name`, `total: int64`, `200 OK`, `application/json`.
- **Observed result:** PASS
- **Implementation:** `internal/provider/server.go:18-44` — fixed payload via `OrderResponseV1`, method gate returns `405` for non-GET, unknown route `404`.
- **Notes:** deterministic responses enable stable contract; extra `notes` field tolerated by design (subset rule).

## Test 3: breaking provider mutation fidelity
- **Expected:** introduces exactly 3 advertised mutations: `status` casing, `customer.name → full_name`, `total int → string`.
- **Observed result:** PASS
- **Implementation:**
  - `internal/provider/server.go:62-71` — `Status: "in_progress"`, `FullName`, `Total: "150000"`
  - Verified error capture: 3 diffs (`status value`, `customer.name missing`, `total type`) in demo Stage 3 and `engineering/03-execution-result.md:103-108`
- **Notes:** Verifier uses `json.Number` so `150000` (number) vs `"150000"` (string) correctly flagged as type mismatch — no false equivalence.

## Test 4: dual provider (V1+V2 routing)
- **Expected:** `/v1/*` remains contract-compliant while `/v2/*` exposes evolved schema.
- **Observed result:** PASS
- **Implementation:** `internal/provider/server.go:89-131` — branch on prefix `/v1/orders/` vs `/v2/orders/`; V1 path returns `OrderResponseV1`, V2 path returns `OrderResponseV2` with `currency`.
- **Notes:** `internal/model/order.go:42-48` — V2 `total` as `string` + `full_name` intentional evolution, never served on V1 route. No cross-contamination observed in test/demo.

## Test 5: consumer client validation
- **Expected:** `FetchOrder` maps `customer.name`, rejects unknown enums and missing fields.
- **Observed result:** PASS
- **Implementation:** `internal/consumer/client.go:34-79` — non-200 rejected, malformed JSON rejected, empty `customer.name` → `contract violation`, unknown status → `contract violation`.
- **Notes:** end-to-end test asserts client succeeds vs V1, fails vs Breaking.

## Test 6: concurrency behavior
- **Expected:** concurrent `Verify` / `FetchOrder` must not race.
- **Observed result:** PASS
- **Implementation:**
  - `Verifier` holds single shared `*http.Client` (safe for concurrent use per stdlib), no mutable state written in `Verify` besides locals; providers stateless
  - `go test -race ./...` passed
- **Notes:** `TestConcurrentContractVerification` spawns 20 goroutines reusing one `Verifier` — correct contention probe.
