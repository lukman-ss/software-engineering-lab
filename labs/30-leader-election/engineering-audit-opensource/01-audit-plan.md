# Engineering Audit Plan

Target Lab: labs/30-leader-election
Implementation Files: internal/**, cmd/**, tests/**
Tests: tests/election_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: (none for this stage)
Main Claims To Verify: lease acquisition, TTL expiration, fencing token safety, single leader invariant, split‑brain prevention
Commands To Run: go test ./..., go test -race ./..., go run ./cmd/demo
Primary Risks: temporary dual‑leader state during pause, in‑memory coordinator not distributed
