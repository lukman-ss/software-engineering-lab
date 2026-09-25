package main

import (
	"fmt"
	"os"

	"labs/17-architecture-decision-record/internal/adr"
)

const adr1 = `# 1. Use Modular Monolith for Core SaaS ERP

Date: 2026-09-25
Status: Superseded by 2

## Context
We are an early-stage startup building a B2B SaaS ERP with a team of 5 engineers and a 3-month launch deadline. 
Domain boundaries are still evolving.

## Decision
We will build the system as a Modular Monolith with strictly separated domain packages and clear in-process interfaces.

## Consequences
- Positive: Fast iteration, unified deployment, no network latency between domains, transactional integrity.
- Negative: Monolithic deployment artifact, single shared database.
- Neutral: Requires discipline to maintain package isolation.

## Review Triggers
- Team size exceeds 20 engineers.
- Need for independent deployment cadences for specific domains (e.g., high-throughput notifications).
`

const adr2 = `# 2. Extract Notification Service to Microservice

Date: 2027-01-15
Status: Accepted
Supersedes: 1

## Context
The engineering team has grown to 25 engineers. The notification engine is experiencing high burst traffic that impacts core ERP performance.

## Decision
Extract the notification module into an independent microservice communicating asynchronously via message queue.

## Consequences
- Positive: Independent scaling and deployment; isolates failure domain from core ERP.
- Negative: Distributed system overhead, eventual consistency, increased observability burden.
`

const adr3 = `# 3. Reject Event Sourcing for Order Management

Date: 2027-02-01
Status: Rejected

## Context
A proposal was made to rewrite the Order Management system to use full Event Sourcing to retain audit history.

## Decision
We reject Event Sourcing at this time. Standard relational audit log tables fulfill legal requirements without event schema migration complexity.

## Consequences
- Positive: Keeps simple relational queries and familiar transaction semantics.
- Negative: Manual audit trail logic required in application layer.
`

func main() {
	fmt.Println("=== Architecture Decision Record (ADR) Lab ===")
	rawDocs := []string{adr1, adr2, adr3}
	records := make([]*adr.Record, 0, len(rawDocs))

	for i, raw := range rawDocs {
		rec, err := adr.Parse(raw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing ADR #%d: %v\n", i+1, err)
			os.Exit(1)
		}
		records = append(records, rec)
	}

	fmt.Printf("Successfully parsed %d ADRs:\n", len(records))
	for _, r := range records {
		lineage := ""
		if r.SupersededBy != 0 {
			lineage += fmt.Sprintf(" -> Superseded by ADR %d", r.SupersededBy)
		}
		if r.Supersedes != 0 {
			lineage += fmt.Sprintf(" -> Supersedes ADR %d", r.Supersedes)
		}
		fmt.Printf("  - [%04d] %-45s | Status: %-10s%s\n", r.ID, r.Title, r.Status, lineage)
	}

	fmt.Println("\nRunning ADR Integrity Linter...")
	linter := adr.NewLinter()
	errs := linter.Validate(records)

	if len(errs) > 0 {
		fmt.Println("Linter detected integrity errors:")
		for _, err := range errs {
			fmt.Printf("  [FAIL] %v\n", err)
		}
		os.Exit(1)
	}

	fmt.Println("Integrity check passed! Decision lineage and lifecycle invariants are intact.")
}
