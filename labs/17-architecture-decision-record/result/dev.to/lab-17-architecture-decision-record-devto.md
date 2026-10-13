---
title: "Architecture Decision Records: Enforcing Lineage and Integrity with Automated Linting"
published: false
description: "Maintain architectural history in Git. Parse Markdown ADRs, enforce monotonic IDs, and validate bidirectional supersession links with automated linting."
tags: go, architecture, testing, devops
---

# Architecture Decision Records: Enforcing Lineage and Integrity with Automated Linting

Software systems suffer when engineering teams lose the historical context behind technical choices. When engineers rotate or scale increases, questions inevitably emerge: why was a modular monolith chosen over microservices? Why was event sourcing rejected for order management?

Without a durable record of why a choice was made and what constraints dictated it at that moment, teams encounter two predictable failure modes:
1. Reversing past decisions arbitrarily without understanding original constraints.
2. Clinging to outdated patterns blindly even after system load and organization size outgrow them.

Traditional architecture documentation stored in third-party wikis drifts out of sync with code branches. Architecture Decision Records (ADRs) solve this by co-locating lightweight Markdown records directly inside the codebase. 

> Architecture decisions must form an immutable, verifiable directed acyclic graph. Once accepted, records are never edited retroactively; they are superseded through explicit bidirectional lineage.

Lab 17 demonstrates an automated parser and structural integrity linter in Go that treats ADRs as executable architectural invariants.

```
ADR LINEAGE LIFECYCLE:

  [ Proposed ] ──> [ Accepted ] ──> [ Superseded by N ] ──> [ Deprecated ]
        │                                  ▲
        └──> [ Rejected ]                  │ (Supersedes: M)
                                    [ ADR N (Accepted) ]
```

Repository implementation reference:
- Git Commit: `4636359a4376246c0af2abdea360ac0b014c7bc4`
- Target Path: `labs/17-architecture-decision-record`

---

## What Qualifies as an Architectural Decision?

Not every code change or pull request warrants an ADR. Following Michael Nygard (2011) and Olaf Zimmermann (2020), an ADR is reserved strictly for Architecturally Significant Requirements (ASRs):

- **Structure:** Subsystem boundaries, modular monolith vs microservices, packaging constraints.
- **Non-Functional Requirements (NFRs):** Latency SLAs, scalability limits, failure isolation boundaries.
- **Dependencies:** Selection of primary databases, message brokers, core frameworks.
- **Interfaces:** Inter-service protocols (REST vs gRPC vs asynchronous messaging).
- **Construction Techniques:** Testing standards, transaction handling strategies.

Routine bug fixes, method refactoring, and local variable naming stay in commit messages and pull requests.

---

## The ADR State Machine

The parser models decision lifecycle states via strongly-typed status constants in `internal/adr/models.go`:

```go
// internal/adr/models.go:3-11
type Status string

const (
	StatusProposed   Status = "Proposed"
	StatusAccepted   Status = "Accepted"
	StatusSuperseded Status = "Superseded"
	StatusDeprecated Status = "Deprecated"
	StatusRejected   Status = "Rejected"
)
```

Each record struct stores identifier metadata and bidirectional replacement links:

```go
// internal/adr/models.go:22-29
type Record struct {
	ID           int
	Title        string
	Status       Status
	SupersededBy int // 0 if not superseded
	Supersedes   int // 0 if does not supersede
	Content      string
}
```

The lifecycle follows strict semantic rules:
- `Proposed`: Under active team evaluation.
- `Accepted`: Approved and actively binding for the architecture.
- `Superseded`: No longer active; replaced by a subsequent record.
- `Deprecated`: Decommissioned without a direct replacement.
- `Rejected`: Evaluated and formally denied. Documenting rejection grounds prevents cyclical debates from resurfacing without fresh technical evidence.

---

## Failure Scenarios in Unchecked ADR Directories

When ADRs are maintained as free-form Markdown without automated validation, documentation integrity deteriorates rapidly:

1. **Dangling Supersession Link:** An ADR marks `Status: Superseded by 99`, but ADR 99 does not exist in the repository.
2. **One-Sided Supersession:** A new ADR declares `Supersedes: 1`, but ADR 1 remains marked as `Status: Accepted`, creating contradictory active guidance.
3. **Non-Monotonic Numbering:** ADR IDs skip numbers (e.g., jump from 0001 to 0003), leaving ambiguity over whether a record was deleted or lost in a merge conflict.
4. **Self-Supersession & Loops:** ADR 1 states `Superseded by 1` or cyclic dependencies emerge (`1 -> 2 -> 3 -> 1`).
5. **Missing Mandatory Sections:** Authors omit `Context`, `Decision`, or `Consequences`, degrading the record to an unsubstantiated edict.

---

## Parsing Markdown into Structured Invariants

The parser in `internal/adr/parser.go` uses deterministic regular expressions to extract structured metadata and verify section integrity:

```go
// internal/adr/parser.go:12-18
var (
	titleRegex        = regexp.MustCompile(`^#\s+(\d+)\.\s+(.+)$`)
	statusRegex       = regexp.MustCompile(`(?i)^Status:\s*([A-Za-z]+)(?:\s+by\s+(\d+))?`)
	supersedesRegex   = regexp.MustCompile(`(?i)^Supersedes:\s*(\d+)`)
	contextRegex      = regexp.MustCompile(`(?i)^##\s+Context`)
	decisionRegex     = regexp.MustCompile(`(?i)^##\s+Decision`)
	consequencesRegex = regexp.MustCompile(`(?i)^##\s+Consequences`)
)
```

During scanning, `Parse` confirms that `Context`, `Decision`, and `Consequences` are present and contain substantive non-empty text:

```go
// internal/adr/parser.go:118-128
if !foundContext || !contextHasContent {
	return nil, fmt.Errorf("context section missing or empty")
}
if !foundDecision || !decisionHasContent {
	return nil, fmt.Errorf("decision section missing or empty")
}
if !foundConsequences || !consequencesHasContent {
	return nil, fmt.Errorf("consequences section missing or empty")
}
```

---

## Validating Lineage and Monotonicity Concurrently

The linter (`internal/adr/linter.go`) performs integrity checks in three sequential phases:

### 1. Monotonic ID Sequence Check

```go
// internal/adr/linter.go:34-39
for i := 0; i < len(ids); i++ {
	if ids[i] != i+1 {
		errs = append(errs, fmt.Errorf("non-monotonic numbering, expected %d but got %d", i+1, ids[i]))
		break
	}
}
```

### 2. Cycle Detection in Supersession Chains

A depth-first search graph traversal detects directed cycles in `SupersededBy` chains:

```go
// internal/adr/linter.go:42-56
state := make(map[int]int) // 0: unvisited, 1: visiting, 2: visited
var detectCycle func(id int)
detectCycle = func(id int) {
	state[id] = 1
	rec, exists := recordMap[id]
	if exists && rec.SupersededBy != 0 && rec.SupersededBy != rec.ID {
		nextID := rec.SupersededBy
		if state[nextID] == 1 {
			errs = append(errs, fmt.Errorf("cyclical supersession detected involving ADR %d", nextID))
		} else if state[nextID] == 0 {
			detectCycle(nextID)
		}
	}
	state[id] = 2
}
```

### 3. Concurrent Bidirectional Verification

Each record is validated across goroutines using a `sync.WaitGroup` and `sync.Mutex` for error aggregation:

```go
// internal/adr/linter.go:72-85
if rec.Status == StatusSuperseded {
	if rec.SupersededBy == rec.ID {
		localErrs = append(localErrs, fmt.Errorf("ADR %d cannot supersede itself", rec.ID))
	} else if rec.SupersededBy == 0 {
		localErrs = append(localErrs, fmt.Errorf("ADR %d is superseded but missing superseded_by reference", rec.ID))
	} else {
		replacement, exists := recordMap[rec.SupersededBy]
		if !exists {
			localErrs = append(localErrs, fmt.Errorf("ADR %d superseded by non-existent ADR %d", rec.ID, rec.SupersededBy))
		} else if replacement.Supersedes != rec.ID {
			localErrs = append(localErrs, fmt.Errorf("ADR %d superseded by ADR %d, but ADR %d does not declare it supersedes ADR %d", rec.ID, rec.SupersededBy, rec.SupersededBy, rec.ID))
		}
	}
}
```

If record $A$ is marked `Superseded by B`, record $B$ must exist and explicitly declare `Supersedes: A`. Conversely, if record $B$ declares `Supersedes: A`, record $A$ must exist and have status `Superseded` with `SupersededBy: B`.

---

## Case Study: SaaS ERP Architecture Evolution

The demo executable (`cmd/demo/main.go`) evaluates a realistic architectural progression across three records:

```
[ ADR 0001: Modular Monolith ]  (Superseded by 2)
              │
              ▼
[ ADR 0002: Notification Microservice ]  (Supersedes: 1)

[ ADR 0003: Reject Event Sourcing ]  (Rejected)
```

1. **ADR 0001 (Modular Monolith):** Adopted when the team has 5 engineers and domain boundaries are evolving rapidly.
   - *Forces:* Fast iteration, ACID transactional integrity, single deployable artifact.
   - *Review Trigger:* Team size exceeds 20 engineers or specific subdomains experience burst loads.
2. **ADR 0002 (Extract Notification Microservice):** Adopted when team size reaches 25 engineers and notification bursts impair core transaction throughput.
   - Updates ADR 0001 to `Status: Superseded by 2`.
   - ADR 0002 declares `Supersedes: 1`.
   - Isolates failure domains at the expense of distributed consistency overhead.
3. **ADR 0003 (Reject Event Sourcing):** Evaluated for order management audit history.
   - Rejected because standard relational audit tables satisfy legal compliance without introducing event schema evolution and eventual consistency complexity.
   - Preserves historical consensus against premature architectural rewrites.

Executing `cmd/demo/main.go` parses all three records and passes all linter invariants:

```text
=== Architecture Decision Record (ADR) Lab ===
Successfully parsed 3 ADRs:
  - [0001] Use Modular Monolith for Core SaaS ERP        | Status: Superseded -> Superseded by ADR 2
  - [0002] Extract Notification Service to Microservice  | Status: Accepted   -> Supersedes ADR 1
  - [0003] Reject Event Sourcing for Order Management    | Status: Rejected  

Running ADR Integrity Linter...
Integrity check passed! Decision lineage and lifecycle invariants are intact.
```

---

## What the Test Suite Actually Proves

The test suite consists of unit tests in `tests/parser_test.go` and `tests/linter_test.go`:

| Test Function | Target Behavior | Verified Invariant |
| :--- | :--- | :--- |
| `TestParse_Valid` | Valid ADR markdown | Correct extraction of ID, Title, Status, and Section contents |
| `TestParse_Superseded` | `Status: Superseded by 2` | Parses `SupersededBy` field value accurately |
| `TestParse_Supersedes` | `Supersedes: 1` clause | Parses `Supersedes` dependency pointer accurately |
| `TestParse_Invalid` | Missing header, status, sections | Rejects malformed records and empty Context/Decision/Consequences |
| `TestLinter_ValidSequence` | Coherent 3-record lineage | Zero errors reported on valid lineage |
| `TestLinter_BrokenReferences` | Structural anomalies | Catches dangling pointer, one-sided link, non-monotonic ID, self-supersession, duplicate ID, and cyclic links |
| `TestLinter_ConcurrencyStress` | 100 alternating records | Passes concurrent validation without data race (`-race`) |

All 23 tests pass cleanly under `go test -race ./...`.

---

## Production Boundaries and Limitations

This implementation isolates structural ADR verification. Note the explicit boundaries:

1. **Strict Regex Parsing:** The parser expects exact `# <ID>. <Title>` and markdown section headings. It is not an arbitrary CommonMark AST parser.
2. **In-Memory Strings:** The demo and tests execute against in-memory markdown strings rather than directly traversing the Git working tree.
3. **No Automated Git Hook Integration:** The lab demonstrates the validation logic; wiring it into `.git/hooks/pre-commit` or GitHub Actions workflows is left to the infrastructure pipeline.
4. **Contextual Example:** The SaaS ERP case study illustrates lifecycle mechanics, not an absolute architectural doctrine.
5. **Structural vs. Semantic Validation:** Automated linters guarantee structural integrity, ID continuity, and referential completeness. They cannot evaluate whether an architectural decision is sound. Peer review remains essential.
