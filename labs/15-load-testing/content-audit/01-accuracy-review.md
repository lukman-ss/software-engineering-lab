# Content Accuracy Review

Target Lab: labs/15-load-testing
Audit Date: 2026-09-26

## Methodology
Compare all content files against approved research and verified engineering implementation (code, tests, execution results).

---

## Finding 1 — INCORRECT STRESS TEST METRICS (Critical)

**Location**: `content/02-master-draft.md` lines 192–206

**Claim**: Stress Test (50 VUs, capacity 5):
- Average: ~504.8ms
- P50: ~609.0ms
- P95: ~981.6ms
- P99: ~1.175s

**Actual** (from `engineering/03-execution-result.md`):
- Average: 739.091366ms
- P50: 669.568542ms
- P95: 1.35666975s
- P99: 1.587917041s

**Impact**: Content understates average latency by ~31% (504.8ms vs 739ms) and P95 by ~28% (981ms vs 1.357s). This misrepresents the severity of tail latency degradation.

---

## Finding 2 — INCORRECT LATENCY MULTIPLIER CLAIM (Major)

**Location**: `content/02-master-draft.md` line 206

**Claim**: "naik ~46x lipat dari ~21ms ke ~982ms pada P95" (increase ~46x from ~21ms to ~982ms on P95)

**Actual**: Stress P95 (1.357s) / Smoke P95 (21.37ms) = ~63.5x increase

**Impact**: Understates the magnitude of queue-induced latency amplification by ~38%.

---

## Finding 3 — INCOMPLETE CODE SNIPPET (Major)

**Location**: `content/03-code-snippets.md` Snippet 2, line 87

**Claimed code**: `_ = resp.Body.Close()`

**Actual implementation** (`internal/loadtest/runner.go` lines 91–92):
```go
_, _ = io.Copy(io.Discard, resp.Body)
_ = resp.Body.Close()
```

**Impact**: The snippet omits the critical `io.Copy(io.Discard, resp.Body)` call, which is necessary to drain the response body before closing. Without this, HTTP connection reuse may fail, potentially causing resource leaks in the load runner.

---

## Finding 4 — Smoke Test Metrics: VERIFIED (Minor, No Issue)

**Location**: `content/02-master-draft.md` lines 193–198

**Claim**: Smoke Test Average ~21.2ms, P50 ~21.2ms, P95 ~21.4ms, P99 ~22.2ms

**Actual**: Average 21.25ms, P50 21.22ms, P95 21.37ms, P99 22.29ms

**Verdict**: Accurate within rounding tolerance.

---

## Findings NOT Flagged (Verified Correct)

- Code snippets 1 and 3 accurately reflect `internal/server/server.go` and `internal/loadtest/metrics.go`
- Architecture diagram in `04-diagrams.md` matches implementation
- All conceptual content (Mental Model, Failure Scenario, Core Concepts) aligns with approved research
- `source-map.md` correctly maps to research and engineering files
- `05-key-takeaways.md` accurately summarizes key lessons
- Warning notes in content brief (sorting slice limitation, time.Timer usage, 10% delay probability) are correctly reflected
- Test claims (`go test -race ./...` passes) verified in engineering-audit
- Thread-safe per-VU result collection correctly described
- MaxIdleConns: 1000 correctly documented as preventing client-side bottleneck