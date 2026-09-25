# Source Map

## Problem & Why This Matters

Research:
research/runs/2026-09-25-architecture-decision-record/05-report.md (Finding 1)
research/runs/2026-09-25-architecture-decision-record/03-evidence.md (Evidence 1, 3)

Sources:
Documenting Architecture Decisions (Nygard)
AWS Prescriptive Guidance

## Mental Model & Core Concept

Research:
research/runs/2026-09-25-architecture-decision-record/05-report.md (Finding 2, Finding 3)
research/runs/2026-09-25-architecture-decision-record/03-evidence.md (Evidence 2, 4, 5)

Implementation:
internal/adr/models.go (Status Lifecycle)

## Failure Scenario & How It Works

Engineering:
engineering/01-design.md
engineering-audit/05-gaps.md

Implementation:
internal/adr/linter.go

## Architecture & Implementation

Implementation:
internal/adr/models.go
internal/adr/parser.go
internal/adr/linter.go

## Code Walkthrough

Implementation:
internal/adr/linter.go

Tests:
tests/linter_test.go

## What the Tests Prove

Tests:
tests/parser_test.go
tests/linter_test.go

Engineering:
engineering/03-execution-result.md

## Case Study

Research:
research/runs/2026-09-25-architecture-decision-record/05-report.md (Finding 4)

Research Audit:
research-audit/03-claim-audit.md (Claim 6 Warning)
research-audit/04-contradictions.md (Contradiction 1)

Implementation:
cmd/demo/main.go

Engineering Audit:
engineering-audit/04-docs-vs-code.md
