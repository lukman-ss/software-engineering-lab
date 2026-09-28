# Docs vs Code

Target Lab: labs/26-contract-testing
Date: 2026-09-28
Scope: README + engineering notes vs code/tests/demo. Research/content excluded per pipeline override.

## README.md vs Code — PASS

- Structure diagram matches actual files (`cmd/demo`, `internal/{consumer,contract,model,provider}`, `tests/contract_test.go`, `engineering/`, `go.mod`). Verified by directory listing.
- Commands (`go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo`) all executed successfully as documented.
- Claims (CDC contracts, CI gate verification, 3 breaking-change classes, V1-preserving V2 evolution) all demonstrated by code + tests + live demo.
- No benchmarks, metrics, or performance claims → nothing to flag as FAKE_BENCHMARK.

## engineering/01-design.md vs Code

1. DOC_CODE_MISMATCH (LOW): §Components claims verifier compares "headers"; implementation ignores response headers (02-code-audit Finding 2).
2. DOC_CODE_MISMATCH (LOW): §Expected Behavior / architecture diagram imply V2 contract verification against `/v2/orders/{id}`; no V2 contract generator or V2-path verification exists — only V1-path verification against the dual provider. The implemented (and tested) claim is "dual provider maintains V1 compatibility", which holds.
3. DOC_CODE_MISMATCH (LOW): Architecture shows `contracts/mobile_order_v1.json` file; contract is in-memory only, no JSON file is written. Cosmetic.
4. "Verification fails, exit code non-zero" — accurate for a CI gate consuming `Verify` result; demo itself exits 0 after correctly blocking (by design). Not a mismatch, noted for precision.

## engineering/02-implementation-notes.md vs Code

- MATCH: file list, pure-Go stdlib choice (`net/http`, `httptest`, no external deps — `go.mod` has zero requires), subset rule, strict declared-field validation, V2-via-dual-routing, stated non-goals (no broker, no async).
- DOC_CODE_MISMATCH (LOW): "verification engine comparing … headers" — same header overclaim as above (single source, counted once in gaps).
- `ponytail:` simplification markers present and honest about ceilings. No overclaim detected.

## engineering/03-execution-result.md vs Independent Run — PASS

- Recorded `go test -v` output reproduced exactly (5/5 PASS, same `[no test files]` lines).
- Recorded race result (`ok … 1.425s`) reproduced: `ok … 1.171s`, no warnings.
- Recorded demo transcript reproduced; only difference is the order of the 3 breaking diffs (map iteration nondeterminism: recorded total/status/name vs observed status/name/total). Tests assert count, not order — correct handling. No invented output: the recorded log is genuine.

## Verdict

Documentation Accuracy: WARNING (3 LOW mismatches, all header/V2-scope/file-diagram cosmetic; core behavioral claims fully match).
No TEST_CLAIM_MISMATCH, no RESEARCH_IMPLEMENTATION_MISMATCH in scope, no FAKE_DEMO, no FAKE_BENCHMARK, no UNVERIFIED_RESULT.
