# Detailed Audit Findings

## HIGH — Factual Inaccuracies

### 1. Fabricated Term: "DINP (Dwell Time Incident Performance)"
**Location:** content/02-master-draft.md:75
**Issue:** "Tanpa error budget, tim tidak sadar sampai DINP (Dwell Time Incident Performance) sudah parah."
**Reality:** "DINP" is not a standard SRE term. No such acronym exists in Google SRE Book, Workbook, Datadog, Prometheus, or industry literature. This appears hallucinated.
**Impact:** Misleads readers with fake terminology.

### 2. Failure Scenario Budget Calculation Error
**Location:** content/02-master-draft.md:70
**Text:** "Traffic normal: 1000 request, 0 error → SLI 100%, budget +0.1"
**Reality:** Demo output shows "Budget Remaining: 1.00" (engineering/03-execution-result.md:55). The budget remaining = (1-0.999)*1000 = 1.00 error-count. "+0.1" incorrectly states the allowed error RATE (0.1%) as the budget VALUE.
**Impact:** Contradicts verified demo output; confuses error rate vs budget count.

### 3. LatencyThreshold Dead Field Presented as Functional
**Location:** content/02-master-draft.md:39, 49-56, content/03-code-snippets.md:5
**Text:** "Sebuah request dikatakan 'good' jika statusnya < 500 **dan** durasinya ≤ latency threshold" with code showing `LatencyThreshold` in Config.
**Reality:** `evaluator.go` never reads `LatencyThreshold`. Latency judgment is delegated to caller's `isGood` closure (demo line 20-22). The evaluator has no knowledge of latency threshold.
**Engineering Audit (opensource):** D3 — "evaluator.go never reads LatencyThreshold; latency judged only by caller's isGood closure" (DOC_CODE_MISMATCH, LOW).
**Impact:** Readers believe evaluator enforces latency; it does not.

### 4. Per-Rule Window Fields Described as Functional
**Location:** content/02-master-draft.md:120-128, content/04-diagrams.md (D3 table)
**Text:** Table showing "Rule | Burn Rate | Window | Budget" with separate windows per rule. Snippet 7 shows `LongWindow`, `ShortWindow`, `BudgetConsumedPct` in `BurnRateRule`.
**Reality:** `engine.Check()` ignores all three fields. All rules share construction-time `shortTracker` and `longTracker` (engine.go:64-65, 73). Per-rule window sizing is unimplemented.
**Engineering Audit (opensource):** D4 — "Per-rule window fields unimplemented; IMPLEMENTATION_OVERCLAIM" (MEDIUM).
**Impact:** Content describes functionality that doesn't exist.

---

## HIGH — Implementation Gaps Not Disclosed

### 5. Demo Has No Recovery Phase
**Location:** content/02-master-draft.md:236-243 ("Recovery / Rollback" section with 4 steps)
**Reality:** Demo ends after Phase 4 (engineering/03-execution-result.md:48-73, cmd/demo/main.go). No recovery simulation exists. The section presents operational procedure as if demonstrated.
**Impact:** Readers believe recovery was demonstrated; it was not.

### 6. "Histogram Latency Buckets" Claim False
**Location:** engineering/01-design.md:26 (design doc), implied in content architecture descriptions
**Reality:** `WindowTracker.Bucket` has only `TotalCount`, `GoodCount`, `BadCount` — no histogram buckets (tracker.go:15-20). Latency is boolean via `isGood` predicate.
**Engineering Audit (opensource):** D5 — "plain bucket counts; no histogram" (DOC_CODE_MISMATCH).
**Impact:** Architecture misrepresented as having latency histograms.

### 7. Zero-Traffic SLI Behavior Undocumented in Key Sections
**Location:** content/02-master-draft.md (no mention in core concept), content/05-key-takeaways.md:19
**Reality:** Evaluator returns SLI=1.0, CanDeploy=true on zero traffic (evaluator.go:44-47, test line 192-196). This is a critical edge-case behavior not highlighted in main concepts.
**Impact:** Operational surprise for users implementing this pattern.

---

## MEDIUM — Overclaims / Misrepresentations

### 8. "100% Test Coverage" Claim Unverified
**Location:** engineering/01-design.md:21 (success criteria), implied in content
**Reality:** No coverage report generated. 6 unit tests cover happy paths and some edge cases, but "100% coverage on core math" is unproven (branches in Record/eviction not fully covered).
**Engineering Audit (opensource):** D5 — "No coverage report (edge branches uncovered)" (TEST_CLAIM_MISMATCH).
**Impact:** Quality claim unsupported by evidence.

### 9. Component Name Mismatch: SLOEvaluator vs Evaluator
**Location:** engineering/01-design.md:34, content/02-master-draft.md:90, content/03-code-snippets.md section header "SLO Evaluator"
**Reality:** Actual type is `Evaluator` (evaluator.go:16, 34). No type named `SLOEvaluator` exists.
**Impact:** API documentation mismatch.

### 10. Burn Rate Threshold Table Includes Unimplemented Rules
**Location:** content/02-master-draft.md:124-129
**Table includes:** "Page Fast 14.4× 1h+5m 2%", "Page Slow 6.0× 6h+30m 5%", "Ticket 1× 3d+6h 10%"
**Reality:** Only two rules implemented: 14.4x PAGE and 6.0x TICKET (demo lines 38-49). The 1×/3d rule is research-only.
**Impact:** Content presents research recommendations as implemented features.

### 11. Source Map References Wrong Test Path
**Location:** content/06-source-map.md:22, 37, 56, etc.
**References:** "tests/slo_test.go" but no `tests/` dir in content or engineering dirs (exists at lab root only).
**Reality:** Test file at `labs/24-slo-sli-error-budget/tests/slo_test.go`. Paths are relative to lab root but source map doesn't clarify.
**Impact:** Minor navigation confusion.

---

## LOW — Clarity Issues

### 12. Compressed Time Ambiguity
**Location:** content/01-content-brief.md:53, content/02-master-draft.md:248
**Text:** "Window demo menggunakan 30 hari" / "Lab gunakan 30 hari rolling window"
**Reality:** Demo uses `30 * time.Minute` compressed to represent 30 days (demo line 24 comment). The warning mentions this but phrasing suggests actual 30-day window.
**Impact:** Ambiguous whether demo runs real 30 days or compressed.

### 13. Engineering Status Claim vs Open-Source Audit
**Location:** content/01-content-brief.md:18-19
**Text:** "Approved Engineering Status: APPROVED"
**Reality:** engineering-audit-opensource found 2 non-blocking issues (DOC_CODE_MISMATCH, IMPLEMENTATION_OVERCLAIM). Internal audit approved; open-source audit flagged gaps.
**Impact:** Status claim depends on which audit is considered authoritative.

### 14. Content Brief Case Study Math Inconsistency
**Location:** content/01-content-brief.md:42
**Text:** "Insiden 10% error rate (10/100) pada SLO 99.9% → burn rate ~100x"
**Reality:** Demo Phase 2 uses 1100 total (1000 baseline + 100 incident), 10 errors = 0.91% error rate, burn rate = 9.09x (not 100x). The 100x would be true if only incident traffic counted (10/100).
**Impact:** Confuses window scope (incident-only vs rolling window).