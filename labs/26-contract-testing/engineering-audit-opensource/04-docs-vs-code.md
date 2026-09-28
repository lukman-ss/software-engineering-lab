# Docs vs Code

Target Lab: labs/26-contract-testing

Sources compared:
  - README.md (implementation & usage claims)
  - engineering/01-design.md (claims + success criteria)
  - engineering/02-implementation-notes.md (limitations + trade-offs)
  - internal/*, cmd/demo (code)
  - tests/* (tests)

## Finding 1 — DOC_CODE_MISMATCH (low)

README "## Running the Lab":
```bash
go test -v ./...
go test -race ./...
go run ./cmd/demo
```
Implementation: commands valid for go 1.22 (go.mod). Verified runnable. README structure listing matches actual tree (cmd/demo/main.go, internal/*, tests/contract_test.go). 
Assessment: PASS — README accurate.

## Finding 2 — TEST_CLAIM_MISMATCH (medium)

engineering/02-design.md §7 Test Strategy lists `TestConcurrentVerification`.
Actual test file defines `TestConcurrentContractVerification`.
Assessment: WARNING — naming mismatch between design spec and actual test. Tests themselves present and passing.

## Finding 3 — DOC_CODE_MISMATCH (medium)

engineering/01-design.md §2 / Components §4 claim the verification engine compares "response schema/types/enums" including "headers".
Code (internal/contract/verifier.go): Response `Headers` field is never read or asserted. Request headers are set (but there are none in the contract).
Tests: No header assertions.
Assessment: WARNING — design over-claims header comparison not implemented.

## Finding 4 — DOC_CODE_MISMATCH (low)

engineering/02-implementation-notes.md §3 What Is Demonstrated lists four bullets, including "Verification failure with 3 exact breaking change diffs".
Demo (cmd/demo/main.go) reproduces this live → matches.
engineering/03-execution-result.md log shows breaking-provider errors; count and semantic content match actual execution (order varies due to map nondeterminism; both are unordered listings).
Assessment: PASS — demo not forged.

## Finding 5 — TEST_CLAIM_MISMATCH (low)

Success criteria #7 claims demo executes "all 3 lifecycle stages (Stage 1/2/3)".
Code's demo has four printed stages: Stage 1 generate, Stage 2 V1 pass, Stage 3 breaking block, Stage 4 dual pass.
Assessment: WARNING — count mismatch in wording (design says 3, demo prints 4). Semantic coverage complete.
