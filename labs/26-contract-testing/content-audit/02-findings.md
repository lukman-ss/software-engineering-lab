# Audit Summary

Audit Date: 2026-09-28
Target Lab: labs/26-contract-testing
Scope: Content accuracy audit only (no code/research changes)
Auditor: Kiro (AI Technical Content Auditor)

---

# Findings

## Verified Accurate Claims

1. **Consumer contract generation** — `GenerateMobileContract()` in `internal/consumer/client.go:82-111` produces correct structure with MobileApp/OrderService, one interaction for `/v1/orders/ORD-123`, `json.Number("150000")` for total.
2. **Provider V1 compliance** — `ProviderV1.ServeHTTP` in `internal/provider/server.go:18-44` returns exact schema matching contract (IN_PROGRESS, customer.name, int64 total). Verification passes.
3. **Provider Breaking changes** — `ProviderBreaking.ServeHTTP` in `internal/provider/server.go:53-80` emits 3 breaking mutations: status casing (`in_progress`), field rename (`full_name`), type change (string total). Verification fails with ≥3 errors.
4. **Provider Dual** — `ProviderDual.ServeHTTP` in `internal/provider/server.go:89-131` supports `/v1` (compliant) and `/v2` (evolved schema). `/v1` path passes verification.
5. **Verifier engine** — `Verifier.Verify()` in `internal/contract/verifier.go:57-115` performs subset validation with `json.Number` preservation via `decoder.UseNumber()`. `diffValues()` in `internal/contract/verifier.go:117-176` detects missing fields, type mismatches, value mismatches.
6. **Test coverage** — All 5 tests in `tests/contract_test.go` exist and pass:
   - `TestConsumerContractGeneration` (line 13)
   - `TestProviderV1_ContractVerification_Success` (line 27)
   - `TestProviderBreaking_ContractVerification_Fails` (line 50)
   - `TestProviderDual_ContractVerification_Success` (line 74)
   - `TestConcurrentContractVerification` (line 96)
7. **Demo orchestrator** — `cmd/demo/main.go:14-71` executes 4 stages matching lab claims.
8. **Model DTOs** — `OrderResponseV1`, `OrderResponseBreaking`, `OrderResponseV2` in `internal/model/order.go` match code behavior.
9. **Mobile client field validation** — `FetchOrder` in `internal/consumer/client.go:34-79` enforces contract assumptions (customer.name presence, status enum check).

## Content Discrepancies — Gaps Documented in Engineering Audit

1. **GAP-01 (HIGH)**: Response header validation is declared in docs and contract (e.g., `02-master-draft.md` line 21, `04-diagrams.md` line 33) but **not implemented** in `Verifier.Verify`. The verifier never compares `ResponseDefinition.Headers` against actual `resp.Header`. A provider returning wrong Content-Type would pass verification.
   - Status: Documented correctly in `engineering-audit-opensource/06-verdict.md`, lab content accurately discloses this limitation.
2. **GAP-02 (HIGH)**: `/v2` endpoint exists in `ProviderDual` (`internal/provider/server.go:112-128`) but has **no consumer interaction**, **no test**, and **no demo coverage**. Lab content correctly notes this (e.g., `02-master-draft.md` line 153, `04-diagrams.md` line 130).
   - Status: Documented correctly in `engineering-audit-opensource/06-verdict.md`, lab content accurately discloses this limitation.
3. **GAP-06 (LOW)**: Error ordering is nondeterministic due to map iteration in `diffValues()`. Recorded demo transcript differs from fresh run. Lab content correctly notes this (e.g., `01-content-brief.md` line 30).
   - Status: Documented correctly in `engineering-audit-opensource/06-verdict.md`, lab content accurately discloses this limitation.

## Content Quality

- Language: Indonesian/English mix as specified in content brief (`01-content-brief.md` line 4).
- Formatting: Code blocks preserve exact file paths, line numbers, comments.
- Diagrams: Accurate reflection of implementation (consumer/provider contract flow, verifier subset logic, 3 breaking diffs).
- Sources: All citations trace to lab files.

---

# Verdict

**APPROVED_WITH_WARNINGS**

Rationale: All factual claims about implementation behavior are accurate. Three documented gaps (GAP-01, GAP-02, GAP-06) are correctly disclosed in lab content per engineering audit findings. Content does not hallucinate behavior or overstate capabilities beyond disclosed limitations.

No content revision required beyond existing disclosures.

---

# Audit Output Files

- `content-audit/01-audit-summary.md` — This file.
- `content-audit/02-findings.md` — Detailed findings (this output).
- `content-audit/09-verdict.md` — Final verdict.

Audit performed per pipeline override: audit content only, no code/research modification, no file deletion.
