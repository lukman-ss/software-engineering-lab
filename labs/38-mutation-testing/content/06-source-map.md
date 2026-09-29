# Source Map

Mapping of Master Draft & Publication Sections to Lab Verification Artifacts:

## 1. Problem & Mental Model (Code Coverage Illusion vs Mutation Testing)
- **Research**: `labs/38-mutation-testing/research/05-report.md` (Finding 1, Finding 2), `03-evidence.md`
- **Research Audit**: `labs/38-mutation-testing/research-audit/03-claim-audit.md`
- **Implementation**: `labs/38-mutation-testing/internal/service/discount_weak_test.go`
- **Tests & Verification**: `weak_cov.out` (100.0% statement coverage verification)

## 2. Core Concepts & Theoretical Foundations (RIP Model, Hypotheses, Score Formula)
- **Research**: `labs/38-mutation-testing/research/05-report.md` (Finding 2, Finding 4, Finding 5, Finding 6, Finding 10)
- **Implementation Data Structures**: `labs/38-mutation-testing/internal/engine/types.go` (MutationType, MutantStatus, Report)
- **Tests**: `labs/38-mutation-testing/tests/engine_test.go` (`TestEngine_MutationScoreDifference`)

## 3. Mutation Operators & AST Mutation Engine
- **Research**: `labs/38-mutation-testing/research/05-report.md` (Finding 3)
- **Implementation**: `labs/38-mutation-testing/internal/engine/mutator.go` (`NewASTMutator`, `GenerateMutations`, `RenderSource`)
- **Tests**: `labs/38-mutation-testing/tests/engine_test.go` (`TestEngine_GeneratesMutants`)

## 4. Concurrent Isolated Runner Architecture
- **Engineering Design & Notes**: `labs/38-mutation-testing/engineering/01-design.md`, `02-implementation-notes.md`
- **Engineering Audit**: `labs/38-mutation-testing/engineering-audit/02-code-audit.md`, `06-verdict.md`
- **Implementation**: `labs/38-mutation-testing/internal/engine/runner.go` (`Runner.Run`)
- **Tests**: `labs/38-mutation-testing/tests/engine_test.go`, race detector validation (`go test -race ./...`)

## 5. Domain Case Study & Assertions Comparison (Weak vs Strong)
- **Implementation**:
  - Domain Target: `labs/38-mutation-testing/internal/service/discount.go`
  - Weak Suite: `labs/38-mutation-testing/internal/service/discount_weak_test.go`
  - Strong Suite: `labs/38-mutation-testing/internal/service/discount_strong_test.go`
- **Execution Output & Demo**: `labs/38-mutation-testing/engineering/03-execution-result.md`, `cmd/demo/main.go`
- **Tests**: `labs/38-mutation-testing/tests/engine_test.go` (`TestDomainService_DirectCalculation`)

## 6. Production Considerations & Tooling Landscape
- **Research**: `labs/38-mutation-testing/research/05-report.md` (Finding 7, Finding 8, Finding 9, Finding 11)
- **Engineering Trade-offs & Limitations**: `labs/38-mutation-testing/engineering/02-implementation-notes.md` (Lines 48-76)
