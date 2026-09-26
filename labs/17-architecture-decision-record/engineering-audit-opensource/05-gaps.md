MISSING_TEST: parse-then-lint end-to-end pipeline (demo path) is not covered by unit tests. Verified manually via go run ./cmd/demo.
MISSING_TEST: empty/nil records input to Linter.Validate not tested.
MISSING_EDGE_CASE: parser does not test malformed header (e.g. "# Title" without number), though spec requires numbered headers.
MISSING_EDGE_CASE: parser does not test lowercase/mixed-case status with leading/trailing whitespace beyond the Title() normalization.
LOW: monotonic numbering check returns first violation only; multi-gap behavior untested (acceptable design).

No BROKEN_IMPLEMENTATION RACE_CONDITION UNHANDLED_ERROR FAKE_DEMO FAKE_BENCHMARK IMPLEMENTATION_OVERCLAIM found.