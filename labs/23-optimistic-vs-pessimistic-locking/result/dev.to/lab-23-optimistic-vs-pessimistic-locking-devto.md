---
title: Optimistic vs Pessimistic Locking in Go
published: false
description: Explore optimistic and pessimistic locking patterns in Go with practical code examples.
tags:
  - go
  - concurrency
  - locking
  - software-engineering
---

## Introduction
Locking ensures data consistency in concurrent systems. This lab demonstrates two common strategies—optimistic and pessimistic locking—implemented in Go.

## Motivation
High‑traffic services need to balance throughput and safety. Pessimistic locking protects resources preemptively, while optimistic locking validates after the fact, reducing contention.

## Lab Limits
- No external packages beyond the standard library.
- Focus on in‑memory data structures.
- Tests run with `go test ./...`.

## Code Evidence
```go
// Pessimistic lock using sync.Mutex
var mu sync.Mutex
func update(value int) {
    mu.Lock()
    defer mu.Unlock()
    // critical section
}

// Optimistic lock using version check
type Record struct { v int; data string }
func (r *Record) Update(newData string, expectedV int) bool {
    if r.v != expectedV { return false }
    r.data = newData
    r.v++
    return true
}
```

## Lab Path
1. Clone repo.
2. Run `go test ./...`.
3. Observe race conditions with `-race` flag.

## Conclusion
Choose locking strategy based on contention level and data modification patterns.

---
Repository: https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/23-optimistic-vs-pessimistic-locking
Commit: {{COMMIT_SHA}}
