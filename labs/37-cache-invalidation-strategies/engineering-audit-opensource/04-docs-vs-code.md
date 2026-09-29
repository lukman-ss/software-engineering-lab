# Docs vs Code

Scope: README.md + engineering/01-design.md, 02-implementation-notes.md, 03-execution-result.md vs code/tests/demo. Research content out of scope per pipeline override.

## README.md vs Code — PASS
- Structure listing matches: cmd/demo/main.go, internal/cache/{store,repo,patterns,stampede}.go, tests/cache_test.go, engineering/{01,02,03}.md, go.mod. All exist.
- 7 feature bullets match code: Cache-Aside (patterns.go:22-47), Write-Through (51-89), Write-Behind (98-166), SingleFlight via x/sync (stampede.go:45-84), XFetch formula (125-136), SWR (173-254), TTL jitter (store.go:79-85). Formula string in README matches ShouldRecompute exactly.
- Run commands (`go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo`) all verified working.

## Engineering 01-design.md vs Code — WARNING (minor)
- DOC_CODE_MISMATCH (LOW): §Architecture lists `jitter.go: TTL jitter calculation` — no such file; jitter lives in store.go:79-85. Only stale-filename reference found.
- All other claims match: MemoryCache/MockDB/CacheAside/WriteThrough/WriteBehind/SingleFlight/XFetch/SWR descriptions match implementations; beta=1.0, 20-goroutine test setup, TTL/stale windows match tests.

## Engineering 02-implementation-notes.md vs Code — PASS
- Files-added list accurate. GetRaw/ReadDelta/formula-guard/SWR-guard/overflow-drop descriptions verified against code. Known-limitations section (in-process singleflight, queue drop, no persistence, SWR timing) honest and matches observed code. Without it these would be overclaims; with it they are scoped correctly.

## Engineering 03-execution-result.md vs Observed — PASS
- Test output block matches fresh run exactly (same 5 tests, same subtests, same PASS).
- Race block matches fresh `go test -race -count=1` (ok, no warnings).
- Demo block matches live `go run ./cmd/demo` line-for-line except jitter sample values (random by design — expected divergence, not fabrication). Query counts reproduced: Cache-Aside 1/1/2, Write-Through 1/1, WriteBehind 0→1, Naive 20, SF 1/20, XFetch 1/2/1, SWR v1/v1/v2.

## Verdict
- No TEST_CLAIM_MISMATCH, no RESEARCH_IMPLEMENTATION_MISMATCH (out of scope), no FAKE_DEMO, no FAKE_BENCHMARK (none claimed).
- One DOC_CODE_MISMATCH (LOW): phantom `jitter.go` in 01-design.md.
