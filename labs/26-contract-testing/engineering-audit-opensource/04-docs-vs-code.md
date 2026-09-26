# Docs vs Code

Target Lab: labs/26-contract-testing
Sources: README.md, engineering/0[1-3].md, tests/contract_test.go, cmd/demo/main.go, internal/*.go

## MISMATCHES FOUND

DOC_CODE_MISMATCH (LOW): README.md:15-16 claims project structure shows `cmd/demo/main.go` as "Executable multi-stage demo". Correct. No mismatch.

RESEARCH_IMPLEMENTATION_MISMATCH: Not applicable (PIPELINE OVERRIDE: research audit skipped).

TEST_CLAIM_MISMATCH: NONE observed.

## Docs Claims Verification

README.md:34-42 "Run Tests and Race Detector": claims `go test -v ./...` and `go test -race ./...`. Match: both commands executed and passed.

README.md:44-48 "Run Demo": claims `go run ./cmd/demo`. Match: demo ran and produced expected 4-stage output.

README.md:8-12 Overview claims:
1. Consumer-Driven Contracts: PASS (GenerateMobileContract yields JSON interaction)
2. Provider CI Gate Verification: PASS (Verifier.Verify against httptest servers)
3. Breaking Change Detection: PASS (3 exact diffs on breaking provider)
4. Safe API Evolution: PARTIAL (V1 preserved in dual; V2 handler exists but no contract/test)

Engineering Design Claims (01-design.md:7-12):
1. CDC capture expectations without testing internal provider: PASS (verifier only sees HTTP in/out)
2. Breaking changes break consumer contract: PASS (3 diffs observed)
3. CI/CD Gate blocks on breaking: PASS (exit code 1 triggered in demo on unexpected pass)
4. Safe API Evolution via V2 DTO preserving V1: PARTIAL (see Finding 9)

Engineering Implementation Notes (02-implementation-notes.md:19-48):
1. Pure Go Implementation: PASS (stdlib only; no external deps)
2. Minimal Subset Verification: PASS (notes/customer.id ignored)
3. Strict Validation on Declared Fields: PASS (enum/int/string diffs triggered)
4. Custom verifier vs Pact: PASS (homemade but correct semantics demonstrated)
5. Known Limitations (no broker, provider state via routing): PASS (accurate)
6. Trade-offs justification: PASS (valid)
7. What Is Demonstrated (items 39-48): All PASS or PARTIAL with justification
8. What Is Not Demonstrated: accurately describes omissions

Engineering Execution Result (03-execution-result.md:5-116):
- Build: SUCCESS (match: `go build ./...` clean)
- Tests: PASS (match: all tests pass)
- Race Detector: PASS (match: `-race` clean)
- Demo: PASS (match: live output identical to engineering/03-execution-result.md:66-114 line-for-line)

## Conclusion

No DOC_CODE_MISMATCH, TEST_CLAIM_MISMATCH, or RESEARCH_IMPLEMENTATION_MISMATCH found that affect core claims. Minor PARTIAL alignments noted where V2 contract/test missing but V1 preservation proven.