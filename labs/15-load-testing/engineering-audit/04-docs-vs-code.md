# Docs vs Code Analysis

## Comparisons

### 1. Structure and File Names
- **Documented in README.md:**
  - `cmd/demo`
  - `internal/server`
  - `internal/loadtest`
  - `tests`
  - `engineering/`
- **Actual Files:** Exactly matching. No phantom directories or missing paths.
- **Status:** PASS

### 2. Execution Commands
- **Documented in README.md:**
  - `go run ./cmd/demo`
  - `go test -v ./...`
  - `go test -race ./...`
- **Actual Verification:** All commands run as documented without requiring configuration or external dependencies.
- **Status:** PASS

### 3. Engineering Claims vs Implementation
- **Claim:** Tail latency spikes under constrained resource (connection pool) saturation.
- **Code:** Implemented via bounded buffered channel (`chan struct{}`) and measured via exact-sorted percentiles.
- **Status:** PASS

### 4. Demo Claims vs Code Execution
- **Documented in `engineering/03-execution-result.md`:**
  - Smoke: ~90 RPS, ~21.8ms Average, ~22.6ms P95.
  - Stress: ~234 RPS, ~200.4ms Average, ~212.3ms P95.
- **Observed in Audit Demo Run:**
  - Smoke: 92.96 RPS, 21.45ms Average, 21.59ms P95.
  - Stress: 236.35 RPS, 200.50ms Average, 210.76ms P95.
- **Status:** PASS (Results are authentic, reproducible, and reflect actual workload dynamics).
