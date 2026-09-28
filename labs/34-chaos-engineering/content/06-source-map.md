# Source Map

## Content Brief

Research:
- research/01-plan.md (Research Topic, Objective)
- research/02-sources.md (Source URLs: Principles of Chaos Engineering, AWS, Netflix, Google SRE)
- research/03-evidence.md (Evidence definitions iya claim utama)
- research/05-report.md (Research Report, Findings 1-3)
- research/06-open-questions.md (Open Questions)

Implementation:
- internal/fault/injector.go (Fault Injector)
- internal/circuitbreaker/circuitbreaker.go (Circuit Breaker state machine)
- internal/monitor/monitor.go (Steady-State Monitor)
- internal/experiment/runner.go (Experiment Runner & Auto-Abort)
- cmd/demo/main.go (Interactive Demonstration)

Tests:
- tests/chaos_test.go (TestFaultInjector, TestCircuitBreakerStateTransitions, TestCircuitBreakerGracefulDegradation, TestExperimentAutoAbortOnSteadyStateViolation, TestConcurrencyAndRace)

---

## Master Draft Sections

### Problem
- research/01-plan.md:4-12 (Research Questions 1-4)
- research/05-report.md:7-8 (Executive Summary)

### Why This Matters
- research/03-evidence.md:5-9 (Evidence 1: definition of Chaos Engineering)

### Mental Model
- research/03-evidence.md:19-22 (Evidence 2: steady state as output metrics)
- research/03-evidence.md:35-38 (Evidence 3: hypothesis formulation)
- research/03-evidence.md:51-55 (Evidence 4: blast radius control)

### Core Concept
- internal/fault/injector.go (Fault injection logic)
- internal/circuitbreaker/circuitbreaker.go (State machine)
- internal/experiment/runner.go (Auto-abort)

### How It Works
- engineering/01-design.md:26-36 (Architecture & Components)
- engineering/02-implementation-notes.md:12-18 (Design Decisions)

### Code Walkthrough
- internal/fault/injector.go (full file)
- internal/circuitbreaker/circuitbreaker.go (full file)
- internal/monitor/monitor.go (full file)
- internal/experiment/runner.go (full file)

### What the Tests Prove
- tests/chaos_test.go (all test functions with assertions)
- engineering/03-execution-result.md:19-32 (Test execution results)
- engineering/03-execution-result.md:39-42 (Race detector results)

### Demo Behavior
- cmd/demo/main.go (demo logic)
- engineering/03-execution-result.md:49-90 (Demo output)

### Recovery / Rollback
- internal/experiment/runner.go:91-97 (terminate function with Clear)
- cmd/demo/main.go:127-132 (traffic after experiment)

### Common Mistakes
- engineering/01-design.md:37-38 (What Is Not Demonstrated)
- research/06-open-questions.md:1-3 (Open Questions)

---

## Code Snippets

### Snippet 1 — Fault Injector
Source File: internal/fault/injector.go

### Snippet 2 — Circuit Breaker
Source File: internal/circuitbreaker/circuitbreaker.go

### Snippet 3 — Steady-State Monitor
Source File: internal/monitor/monitor.go

### Snippet 4 — Experiment Runner
Source File: internal/experiment/runner.go

---

## Audit Verdicts

Research Approval:
- research-audit/07-verdict.md (APPROVED)

Engineering Approval:
- engineering-audit/06-verdict.md (APPROVED)

Research Warnings:
- research-audit/06-gaps.md (Source URLs not deep-linked, language-specific differences)

Engineering Warnings:
- engineering-audit/05-gaps.md (Unused field `errorRate` in injector, directory naming discrepancy)

---

## Full Source List

```
internal/
├── fault/
│   └── injector.go
├── circuitbreaker/
│   └── circuitbreaker.go
├── monitor/
│   └── monitor.go
└── experiment/
    └── runner.go

cmd/
└── demo/
    └── main.go

tests/
└── chaos_test.go

research/
├── 01-plan.md
├── 02-sources.md
├── 03-evidence.md
├── 04-contradictions.md
├── 05-report.md
├── 06-open-questions.md
└── runs/

research-audit/
├── 01-audit-plan.md
├── 02-source-audit.md
├── 03-claim-audit.md
├── 04-contradictions.md
├── 05-code-audit.md
├── 06-gaps.md
└── 07-verdict.md

engineering/
├── 01-design.md
├── 02-implementation-notes.md
├── 03-execution-result.md

engineering-audit/
├── 01-audit-plan.md
├── 02-code-audit.md
├── 03-test-audit.md
├── 04-docs-vs-code.md
├── 05-gaps.md
└── 06-verdict.md

README.md
go.mod
```