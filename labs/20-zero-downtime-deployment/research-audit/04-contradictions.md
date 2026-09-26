# 04 — Contradictions

Target Lab: `labs/20-zero-downtime-deployment`
Research Run: `research/runs/2026-09-26-zero-downtime-deployment/`
Audit Date: 2026-09-26

## Evaluation

Comparison across research files (`01-plan.md`, `02-sources.md`, `03-evidence.md`, `04-contradictions.md`, `05-report.md`, `06-open-questions.md`):

1. **Internal Consistency**:
   - `01-plan.md` identifies NGINX OSS drain limitations and Octane vs PHP-FPM differences as risks; `05-report.md` consistently qualifies these throughout the findings and limitations section.
   - `03-evidence.md` marks Redis upgrade as NOT VERIFIED; `04-contradictions.md`, `05-report.md`, and `06-open-questions.md` consistently maintain that this is unverified.

2. **Source Conflict vs Implementation Specifics**:
   - The distinction between Blue-Green (instant router rollback, 2x cost) and Rolling (capacity efficient, slower rollout undo) is accurately presented as an architectural trade-off rather than conflicting doctrine.
   - NGINX OSS passive checks vs NGINX Plus active checks are explicitly documented without claiming OSS has commercial features.

## Findings

No material contradictions found.
