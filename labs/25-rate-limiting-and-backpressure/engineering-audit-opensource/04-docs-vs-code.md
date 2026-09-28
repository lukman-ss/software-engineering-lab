# Docs vs Code Audit

Comparing README.md (source of truth for claims) with implementation + tests + demo.

## DOC_CODE_MISMATCH: None detected.

README claims matched code line-for-line:

- Token Bucket & Leaky Bucket algorithms: present.
- Per-tenant registry key isolation: present.
- Non-blocking submission returning ErrQueueFull: present.
- AWS jitter strategies with formulas: present.
- HTTP 429 middleware with Retry-After header: present.
- Demo sections: matches four printed blocks.

## TEST_CLAIM_MISMATCH: None detected.

Tests exercise exactly the claims in README:
- Burst, refill, leak-rate.
- Tenant isolation.
- Queue rejection under load + concurrency safety + stop safety.
- Jitter bounds.
- HTTP 429 and anonymous.

Gaps in test coverage flagged in gaps, not mismatches.

## RESEARCH_IMPLEMENTATION_MISMATCH: Not audited (pipeline override).

## DOC_TEST_MISMATCH: None.

Test names and assertions align with README behaviors; no test contradicts README.

## DOC_DEMO_MISMATCH: None.

README demo commands (`go test ./...`, `go test -race ./...`, `go run ./cmd/demo`) produce outputs matching engineering/execution-result.md (and live demo above). Executed commands succeed.

## README_ACCURACY: PASS.

README reflects actual code behavior; no overclaiming; limitations listed (in-memory only, no distributed); test instructions work. Severity LOW: no issue.

## ENGINEERING NOTES vs CODE

Engineering/01-design.md:  
- Architecture diagram matches code flow.  
- Components list matches packages.  
- Test strategy matches files.  
- Execution plan matches audit-run commands.  
- Implementation Decisions: pure stdlib (true), monotonic time (true), deterministic test helpers (claimed but unit tests use wall-clock sleeps — mild mismatch). Severity LOW note only.  

Engineering/02-implementation-notes.md:  
- Files added list exact.  
- Core Design Decisions: mutex sync (true), monotonic time (true), select-default backpressure (true), AWS Full Jitter citation (true).  
- Implementation-Specific Choices: floating point tokens (true), X-API-Key header (true).  
- Known Limitations: in-memory only (true).  
- Trade-offs: mutex vs CAS noted.  
- What Is Demonstrated: matches four demo sections.  
- What Is Not Demonstrated: distributed & autoscaling (true).  
Alignment strong; Severity LOW none.

## Summary

No DOC_CODE_MISMATCH, TEST_CLAIM_MISMATCH, RESEARCH_IMPLEMENTATION_MISMATCH, or DOC_TEST_MISMATCH detected. README accurately describes code. Minor doc vs note mismatch: design claims deterministic test helpers while tests use real sleeps. Not material; flagged as WARNING in gaps.