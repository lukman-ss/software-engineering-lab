# Diagrams

All diagrams derived from the actual implementation of `labs/26-contract-testing`. No components invented.

---

## Diagram 1 — Contract Testing Lifecycle & CI Gate

```text
                       LAB 26 CONTRACT TESTING LIFECYCLE

[1] Consumer (Mobile App)
    └─ defines minimal expectations
       (only id, status, customer.name, total)

[2]  consumer.GenerateMobileContract()
    └─ produces Contract JSON (1 interaction)
       └─ contract.json artifact

[3]           ┌──────────────────────────────┐
              │   Contract Verifier (CI gate)   │
              │   verifier.Verify(baseURL, c)   │
              │   validates: status + body only │
              └──────────────────────────────┘
                        │
       ┌────────────────┼──────────────┐
       │                │              │
       ▼                ▼              ▼
 [Provider V1]   [Provider Breaking]  [Provider Dual]
  (/v1 route)     (/v1 route)        (/v1 + /v2 routes)
       │                │              │
       │                │              │
    ┌──▼──┐          ┌──▼──┐        ┌──▼──┐
    │ PASS│        │ FAIL │        │ PASS│
    │ (V1 │        │ (V1 │        │ (V1 │
    │ route)│        │ route)│      │ route)│
    └──┬──┘          └──┼──┘        └──┬──┘
       │                │              │
       ▼                ▼              ▼
  ALLOWED    BLOCKED + error list   ALLOWED
  deploy      (3 diffs):            canary /
                  - missing 'customer.name'
                  - type mismatch 'total'
                    int64 vs string
                  - value mismatch 'status'
                    "IN_PROGRESS" vs "in_progress"
```

**What the tests prove:**
- `TestProviderV1_ContractVerification_Success` → V1 path PASS → ALLOWED.
- `TestProviderBreaking_ContractVerification_Fails` → Breaking path FAIL with ≥3 diffs → BLOCKED.
- `TestProviderDual_ContractVerification_Success` → Dual V1 path PASS → ALLOWED.

---

## Diagram 2 — Subset Verification (Minimal Contract vs Full Response)

```text
CONSUMER CONTRACT BODY (expected — minimal):        PROVIDER V1 RESPONSE (actual — full):

  {                                                    {
    "id": "ORD-123",     ✓ check                      "id": "ORD-123",
    "status": "IN_PROGRESS", ✓ check                   "status": "IN_PROGRESS",
    "customer": {          ✓ recurse                  "customer": {
      "name": "Budi Santoso" ✓ check                    "id": "cust-100",
    },                              ...                     "name": "Budi Santoso"
    "total": 150000  ✓  json.Number                    },
                       vs int64                          "total": 150000,
  }                                                    "notes": "...",      ← ignored
                                                   }                          (not in contract)

VERIFICATION: PASS  (every field in contract exists
                     in actual with correct path/type/value;
                     extra provider fields are not checked)
```

**Key point:** `decoder.UseNumber()` in `verifier.go:97` makes `json.Number` the expected type; the verifier compares types explicitly (see `diffValues` Snippet 7), so string `"150000"` ≠ number `150000`.

---

## Diagram 3 — Breaking Provider vs Contract (3 Diffs)

```text
CONTRACT EXPECTS:                BREAKING PROVIDER EMITS:

  status  : "IN_PROGRESS"          status     : "in_progress"
  customer.name : "Budi Santoso"   customer   : { "full_name" : "Budi Santoso" }
  total   : 150000  (number)       total      : "150000"  (string)

VERIFIER diffValues OUTPUT:
  ┌─ childPath = "total"
  │    expected: 150000 [json.Number]
  │    actual:   150000 [string]
  │    → path 'total': type mismatch
  ├─ childPath = "status"
  │    expected: "IN_PROGRESS"
  │    actual:   "in_progress"
  │    → path 'status': value mismatch
  └─ childPath = "customer.name"
       expected exists, actual missing (customer only has "full_name")
       → missing expected field 'customer.name'

Result: VerificationResult{Passed: false, Errors: [3 entries]}
```

**Verified by:** `TestProviderBreaking_ContractVerification_Fails` asserts `len(result.Errors) >= 3` and that `FetchOrder` returns an error.

---

## Diagram 4 — Safe Evolution: Expand/Contract via Dual Provider

```text
SAFE API EVOLUTION (expand → migrate → contract)

PHASE 1 (Expand)                    PHASE 2 (Migrate)            PHASE 3 (Contract)
────────────────────────────────    ─────────────────────        ──────────────────────────
ProviderDual (current lab state)    Consumer migrates            Old V1 removed
────────────────────────────────    ─────────────────────        ──────────────────────────

  /v1/orders/{id}  ← V1 schema       Consumer A: still on /v1    (future)
  /v2/orders/{id}  ← V2 schema       Consumer B: switched to /v2  /v1 routes retired
  (both live)       (consumer       (independent                 (only /v2)
                     chooses route)  per-consumer)

CONTRACT STATUS at each phase:
  V1 contract vs /v1: PASS           V1 contract vs /v1: PASS      V2 contract vs /v2: PASS
  (V2 has own /v2 contract TBD)      (V2 contract vs /v2: PASS)      (old /v1 contract removed)
```

**Verified by:** `TestProviderDual_ContractVerification_Success` — V1 route of Dual provider passes the existing V1 contract; `/v2` route available for new consumers. Lab implements only Phase 1 (Expand); Migrate and Contract are documented migration phases (expand/contract pattern from research Finding 7).

---

## Diagram 5 — Test Coverage Matrix

```text
                          ┌─ Verified ──┐
                          │             │
TEST                     │ PASS  FAIL   │ LAB OUTPUT REFERENCE
──────────────────────────────────────────────────────────────
TestConsumerContractGen ─┤  ✓           │ asserts Consumer=MobileApp, 1 interaction
                         │              │ path=/v1/orders/ORD-123

TestProviderV1_* ────────┤  ✓           │ V1 server → Passed=true; FetchOrder OK

TestProviderBreaking_* ──┤      ✓       │ Breaking server → Passed=false; ≥3 errors;
                         │              │   FetchOrder returns error

TestProviderDual_* ──────┤  ✓           │ Dual server (V1 route) → Passed=true;
                         │              │   FetchOrder OK on /v1

TestConcurrent_* ────────┤  ✓           │ 20 goroutine verify calls → no race;
                         │              │   go test -race PASS
```

**Note:** The 6 non-blocking engineering gaps (GAP-01 through GAP-06, engineering-audit-opensource/06-verdict.md) are explicitly out of the verified coverage set.

**Verification scope disclaimer:** The current implementation validates HTTP status code and response body fields only. Response header validation is a planned enhancement (see GAP-01). Verifier error ordering depends on map iteration; order may vary between runs (see GAP-06). The `/v2` endpoint exists in ProviderDual but has no associated consumer contract or verification test (see GAP-02). Production Pact implementations include header validation and deterministic ordering.
