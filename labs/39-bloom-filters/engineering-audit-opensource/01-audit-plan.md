# Engineering Audit Plan

Target Lab: labs/39-bloom-filters
Implementation Files: src/
Tests: tests/
Executable/Demo: cmd/demo (if present)
Approved Research Inputs: (none for this stage)
Main Claims To Verify: Bloom filter correctness, false positive rate, concurrency safety
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks: Incorrect hash functions, concurrency races, missing edge‑case tests