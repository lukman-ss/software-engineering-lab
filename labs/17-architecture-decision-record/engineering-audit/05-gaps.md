# Gap Analysis

## Gap 1

Type: MISSING_EDGE_CASE
Severity: LOW
Description: The parser verifies that markdown section headers (`## Context`, `## Decision`, `## Consequences`) exist, but does not verify whether there is any non-empty body text under each section.
Impact: An ADR with empty sections would pass parsing validation.

## Gap 2

Type: MISSING_EDGE_CASE
Severity: LOW
Description: The linter verifies 1:1 bidirectional supersession links and prevents direct self-supersession, but does not perform full cycle detection across longer chains (e.g., A supersedes B, B supersedes C, C supersedes A).
Impact: Circular supersession references spanning more than one hop could bypass validation.
