# Revision Plan

Target Lab: labs/14-circuit-breaker

Previous Audit Status: APPROVED_WITH_WARNINGS

## Blocking Issues

None identified.

## Non-Blocking Issues

1. **AWS URL Redirection**: The link `https://aws.amazon.com/builders-library/timeouts-retries-and-backoff-with-jitter/` permanently redirects (301) to AWS Builder Center (`https://builder.aws.com/content/3EumjoZascWd1oZiEgL8ORlv3qE/timeouts-retries-and-backoff-with-jitter`).

2. **Error Filter Simplification**: While `research/09-failure-modes.md` identifies counting 4xx errors as a failure mode, the implementation in `circuitbreaker.go` increments failure counters on any non-nil error without providing an error predicate.

3. **Unimplemented Metrics**: Observability metrics described in `research/08-observability.md` and `README.md` are not exposed via any metrics export in the Go code.

4. **HALF_OPEN Concurrency Test**: Test suite verifies concurrency in general, but lacks a targeted test verifying probe throttling when multiple concurrent callers hit `HALF_OPEN`.

## Files To Modify

- `research/02-sources.md` - Update AWS URL
- `README.md` - Add clarification about error filtering and metrics being architectural guidance
- `research/09-failure-modes.md` - May need clarification note
- `research/08-observability.md` - May need clarification note

## Verification Plan

- source verification: Check that updated URL resolves correctly
- documentation consistency: Ensure README matches research files
- build: Run `go build ./...` to ensure no syntax errors
- tests: Run `go test ./...` to ensure tests still pass
- demo: Run `go run ./cmd/demo` to verify demo still works