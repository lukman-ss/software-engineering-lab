# Content Revision Record

## Target Lab
labs/15-load-testing

## Audit Source
engineering-audit-opensource/05-gaps.md — G5 (MEDIUM): Latency not recorded for failed requests.

## Issue Summary
- **Severity**: MEDIUM
- **Finding**: Engineering code-audit Finding 3 (G5) identified that the load runner only records latency for successful (HTTP 2xx) requests. Transport errors and HTTP >= 400 responses are counted as errors but their durations are not captured. The content previously implied all request latencies were collected for percentile computation, without explicitly clarifying this scoping.

## Changes Made

### 01-content-brief.md
- **Line 24 (Warnings)**: Added new warning bullet clarifying that percentile metrics (P50/P90/P95/P99) only cover successful requests (HTTP 2xx); failed requests are counted as errors but their latency is not measured.

### 02-master-draft.md
- **Line 34 (How It Works)**: Clarified that the load runner collects latency only for successful requests (HTTP 2xx). Added explicit statement that failed requests (transport errors or HTTP >= 400) are not recorded in latency metrics but are counted as errors.
- **Line 170 (Code Walkthrough Snippet 2)**: Added note after the explanation that latency is only recorded for HTTP 2xx responses, and that error responses are counted separately as `errs`.
- **Line 218 (Production Considerations)**: Added bullet noting that in production, error response latencies should also be captured—not just error count—to reveal timeout or dependency failure patterns.

### 03-code-snippets.md
- **Line 100 (Snippet 2 Explanation)**: Added explicit note that latency is only recorded for successful (HTTP 2xx) requests; transport errors and HTTP >= 400 responses do not produce entries in the `lats` slice but are counted as `errs`. Clarified that P50/P95/P99 metrics reflect only successful-request latency.

### 05-key-takeaways.md
- **Line 7 (new takeaway)**: Added point #7 clarifying that percentile metrics only cover HTTP 2xx responses; error latencies are not measured (only error count is tracked). Warned this can mislead analysis if failures carry high latency.

## Files Not Modified (Out of Scope per Pipeline Override)
- `engineering/01-design.md` — DOC_CODE_MISMATCH (iterations vs bounded duration) — engineering doc, not content file.
- `engineering/02-implementation-notes.md` — Tail latency "strictly queuing" claim — engineering doc, not content file. Content master-draft correctly describes the 10% random slowdown.
- `engineering/03-execution-result.md` — Stale test list — engineering doc, not content file.
- `cmd/demo/main.go` — P90 not printed — code file, not content.

## Verification
- `go test -v ./...` — all 12 tests pass
- `go test -race ./...` — no races detected
- Content reviewed against code audit findings from engineering-audit-opensource/
- Content brief updated to include the new limitation warning

## New Changes (Addressing Content Audit Gaps)

### Audit Source
content-audit/09-verdict.md — Non-blocking warnings from engineering-audit-opensource/05-gaps.md G2 and G3.

### Issues Addressed

1. **G3: P90 Metric Omission in Demo Output**
   - **Severity**: LOW
   - **Finding**: `Result` struct computes P90 but demo's `printResults` only prints P50/P95/P99, creating inconsistency with content documenting P50/P90/P95/P99 as standard output.
   - **Action**: Clarify in content that demo prints subset while P90 remains computed.

2. **G2: Missing Clarification on 10% Slowdown vs Queuing**
   - **Severity**: LOW
   - **Finding**: Content describes 10% slowdown as additive but doesn't explicitly contrast with primary queuing mechanism.
   - **Action**: Add sentence clarifying it's an amplifier on top of queuing delays.

### Files Modified

#### 02-master-draft.md
- **Line 70 (Implementation section)**: Added note that `printResults` in demo prints subset (P50/P95/P99) while `Result` computes all seven metrics including P90.
- **Line 120 (Code Walkthrough)**: Added clarifying sentence that 10% slowdown is an amplifier on top of queuing delays; stress test would still show P95 degradation from queue buildup alone.

#### 03-code-snippets.md
- **Line 100 (Snippet 2 Explanation)**: Added note about demo printing subset of metrics.
- **Line 48 (Snippet 1 Explanation)**: Added clarifying sentence about 10% slowdown being additive amplifier to queuing delays.

#### 04-diagrams.md
- **Line 41 (Tail latency causes)**: Updated wording to clarify 10% probability is an amplifier on top of queue wait.

#### 05-key-takeaways.md
- **Line 8 (existing takeaway)**: Added note about P90 being computed but not printed in demo output.

## Files Not Modified (Confirmed)
- `engineering/` files — out of scope per pipeline override
- `cmd/demo/main.go` — code file, not content
- `internal/` — code files, not content
- `research/` — research files, not content

## Verification of New Changes
- Content now accurately reflects that demo prints P50/P95/P99 subset while P90 is computed
- Content clarifies 10% slowdown as amplifier on top of queuing mechanism
- All changes maintain consistency with existing content style and technical accuracy
- No code or research files modified

## Status
READY_FOR_CONTENT_ADAPTER
