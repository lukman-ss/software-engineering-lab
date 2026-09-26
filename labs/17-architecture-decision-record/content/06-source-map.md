# Source Map

## Problem & Why This Matters

Research:
research/05-report.md (Finding 1, Finding 2)
research/03-evidence.md (Evidence 1, Evidence 2)

Sources:
Source 1: Documenting Architecture Decisions (Nygard)
Source 5: AWS Prescriptive Guidance

## Mental Model & Core Concept

Research:
research/05-report.md (Finding 1, Finding 5)
research/03-evidence.md (Evidence 1, Evidence 4)

Implementation:
internal/adr/models.go (Status & Record struct)

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
research/05-report.md (Finding 5)
research-audit/03-claim-audit.md (Claim 6 - MEDIUM warning on anecdotal evidence)

Implementation:
cmd/demo/main.go

Engineering Audit:
engineering-audit/04-docs-vs-code.md
