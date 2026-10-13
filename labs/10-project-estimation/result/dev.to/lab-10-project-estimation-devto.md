---
title: "Project Estimation: Eight Screens Are Not Eight Days of Work"
published: false
description: "Estimate a service-booking project with task ranges, spikes, risk, contingency, and capacity—and inspect what the Go model does not know."
tags: go, softwareengineering, projectmanagement, backend
---

An application has eight features. Assigning one day to each produces an eight-day estimate, but says nothing about webhook behavior, sandbox access, booking validation, or the engineer's actual availability.

Lab 10 replaces that shortcut with a small Go estimation model. Its case study is a vehicle-service booking application: login, online booking, branch selection, mechanic selection, payment, WhatsApp notifications, an admin dashboard, and Excel reports. The repository implements estimation arithmetic and validation—not those application features or vendor integrations.

The important result is not a more convincing single number. It is an explanation of scope, effort, uncertainty, dependencies, assumptions, and the capacity needed to turn effort into elapsed working time.

## Page count hides the work

The naive estimator deliberately returns one day per page:

```go
// estimation.go; complete function
func EstimateByPageCount(pageCount int) int {
    if pageCount <= 0 {
        return 0
    }
    return pageCount
}
```

`TestNaiveEstimator` expects ten days for ten pages. That assertion confirms the function's behavior; it does not establish that a ten-page project can be delivered in ten days.

The README's payment breakdown explains what page count misses: documentation discovery, authentication, transaction creation, callback verification, error paths, retry/idempotency, and sandbox testing. A screen is an interface boundary, not a unit of implementation effort.

Similarly, branch selection involves master data and a selector, mechanic selection needs relationships and availability checks, and Excel reporting includes filtered data and file generation. These are planning examples in the README, not implemented booking endpoints.

## Make the estimation unit explicit

Each executable planning task carries a name, a three-point effort range, a risk label, spike effort, and assumptions:

```go
// estimation.go; complete type definitions
type EstimateRange struct {
    Min        float64
    MostLikely float64
    Max        float64
}

type Task struct {
    Name        string
    Estimate    EstimateRange
    Risk        RiskLevel
    SpikeDays   float64
    Assumptions []string
}
```

`Project.Estimate` combines the tasks directly attached to the project with tasks inside each feature. Feature names organize the input; the calculation flattens their tasks into one collection. `TestFeatureBreakdownAffectsTotal` checks that three feature tasks contribute six minimum implementation days and a task count of three.

There is no task identity or deduplication. Listing the same work in both collections counts it twice. Nor does the model store task-to-task dependencies or create a schedule from feature order.

## Unknown work needs a bounded investigation

The README distinguishes familiar work from work whose authentication, webhook behavior, or vendor environment is not understood. Unknown does not mean zero effort. It also should not become an unsupported implementation guess.

A spike is timeboxed discovery with a specific output: enough evidence to identify and estimate the implementation tasks. The payment example investigates authentication, transaction creation, signature verification, sandbox behavior, and failure/retry handling. It is not a production integration delivered inside the discovery timebox.

The actual validation rule is strict:

```go
// estimation.go; excerpt from Task.Validate
if t.Risk == RiskUnknown && t.SpikeDays == 0 {
    return fmt.Errorf("task %q: %w", t.Name, ErrUnknownRiskNoSpike)
}
```

Every `RiskUnknown` task must have positive `SpikeDays`, even when its implementation range is nonzero. `TestUnknownRiskWithoutSpikeRejected` supplies a 1/2/3 estimate with no spike and expects rejection. `TestUnknownRiskWithSpikeAndExplicitEstimateAccepted` supplies a positive spike and expects acceptance.

This contradicts the README's failure-scenario row claiming that an unknown task with zero spike and a nonzero estimate is allowed after discovery. The executable model has no completed-spike flag. Once discovery is complete, keeping `RiskUnknown` while setting `SpikeDays` to zero still fails validation; a re-estimated task needs an appropriate supported risk classification.

Also, `RequiredSpikes` includes any task with positive spike days, not just unknown-risk tasks. Spike time is counted as fixed effort in all three range positions.

## A weighted estimate is not a deadline probability

```go
// estimation.go; complete function
func (r EstimateRange) Expected() float64 {
    return (r.Min + 4.0*r.MostLikely + r.Max) / 6.0
}
```

The middle input gets four times the weight of each endpoint. The model validates nonnegative values and `Min <= MostLikely <= Max`, then sums minimum, weighted expected, and maximum values across tasks.

The word `Expected` is an output name for this weighted calculation. There is no sampled distribution, Monte Carlo simulation, percentile, or historical calibration. A maximum estimate is not a promised latest delivery date, and the output range is not a statistical confidence interval.

For the arithmetic below, the source snapshot is commit `82af988b89439fa92b1200406c5afbd59ce72a56`. Go is unavailable in this environment, so the tests were inspected but not executed. The numbers are independently checked arithmetic derived from the fixture and formulas, not captured Go test output.

## Work through the booking fixture

`TestBookingServiceCaseStudy` has six direct tasks and two tasks inside an admin feature:

| Task | Min | Most likely | Max | Risk | Spike days |
|---|---:|---:|---:|---|---:|
| Login | 1 | 2 | 3 | Low | 0 |
| Booking Online | 3 | 5 | 8 | Medium | 0 |
| Pilih Cabang | 1 | 2 | 4 | Low | 0 |
| Pilih Mekanik | 2 | 3 | 5 | Medium | 0 |
| Payment Gateway | 4 | 6 | 12 | Unknown | 2 |
| WhatsApp Notification | 1 | 2 | 4 | Unknown | 0.5 |
| Dashboard Admin | 3 | 5 | 7 | Medium | 0 |
| Laporan Excel | 2 | 4 | 6 | High | 0 |

Task estimates use engineer-days. These are illustrative fixture inputs, not measured productivity or a quotation for a real booking project.

The minimum values sum to 17, the most-likely inputs sum to 29, and the maximum values sum to 49. The sum of weighted expectations is `(17 + 4 × 29 + 49) / 6 = 30.333…`, rounded by the implementation to 30.33.

Adding the two fixed spikes gives 2.5 engineer-days and a base range of 19.50 / 32.83 / 51.50. This is one complete fixture calculation; it should not be combined with the separate reference example in the README, which uses different task estimates and spike lengths.

The fixture does not contain separate testing/UAT or deployment tasks. Some README breakdowns include testing inside a task, but the fixture does not specify how much of those activities its ranges cover. The model cannot detect omitted work.

## Risk controls a heuristic, not external reality

Project risk is derived from task counts:

- Any unknown task makes overall risk High.
- Otherwise, at least 30% High tasks makes overall risk High.
- Otherwise, any High or Medium task makes overall risk Medium.
- Otherwise, it is Low.

This is not effort-weighted: a short unknown investigation can elevate the whole project. The booking fixture is High because it contains unknown tasks.

Automatic contingency applies only when `AutoContingency` is true and the supplied rate is zero. It selects 25% for High, 15% for Medium, and 10% for Low. A positive explicit rate takes precedence, even if automatic selection is enabled. A zero rate with automatic selection disabled means no contingency.

The defaults are teaching choices, not universal planning rules. The README's broader percentage ranges should not be presented as additional executable settings.

### Contingency is applied to every endpoint

The model computes `ContingencyEffort` as base expected effort times the rate. It scales minimum and maximum by the same percentage, rather than adding that single expected contingency amount to both endpoints:

```go
// estimation.go; excerpt from Project.Estimate
contingencyDays := baseRange.ExpectedDays * contingencyRate
```

For this fixture, 25% produces expected contingency of 8.2075 engineer-days and a rounded final range of **24.38 / 41.04 / 64.38 engineer-days**. `ContingencyEffort` itself is not rounded in the stored result; formatted output may display fewer decimals.

This is an explicit assumption about risk buffering. It does not discover vendor delays, validate assumptions, or guarantee that the chosen buffer is sufficient.

## Effort becomes duration through capacity

```go
// estimation.go; excerpt from Project.Estimate
effectiveDailyCapacity := float64(p.EngineerCount) * p.Availability
```

Each final effort endpoint is divided by this capacity and rounded to one decimal. One engineer at 70% availability gives 0.7 engineer-days of modeled daily capacity. For the fixture, that produces **34.8 / 58.6 / 92.0 working days**.

`DurationRange.Weeks` divides each duration by five. The endpoint range is therefore approximately **7.0–18.4 five-working-day weeks**, not literal calendar days and not a date-aware schedule. Holidays, vendor waiting, dependency ordering, and the critical path are not represented.

`TestMultipleEngineersReducesCalendar` checks that the formula gives a smaller duration with two engineers. It does not prove that every real task can be split evenly or that a project finishes twice as fast. The linear capacity formula assumes that the modeled effort can use that capacity.

The README's separate reference example says 15–23 engineer-days with one engineer at 70%, yet mentions three to four weeks and converts only the lower endpoint. Dividing both endpoints yields approximately 21.4–32.9 working days, or 4.3–6.6 five-working-day weeks. The executable fixture above uses different inputs; neither example should be quoted as a verified delivery commitment.

## Confidence is a communication label

`calculateConfidence` returns Low for any unknown task or High overall risk. Medium risk returns Medium. Low risk with at least one assumption returns High; without assumptions it returns Low.

The function does not inspect whether assumption text is true, specific, or meaningful. `contingencyRate` is passed into it but is unused. Adding buffer therefore does not increase confidence under this implementation. The booking fixture remains Low confidence despite its assumptions and contingency.

`TestAssumptionsIncreaseConfidence` and `TestMissingAssumptionsReduceConfidence` exercise the low-risk label behavior. They do not establish a statistical probability of meeting a deadline.

## Validation is useful, but bounded

Tests cover invalid range order, negative effort/spike/contingency, unsupported risk labels, empty projects, empty task names, invalid availability, and nonpositive engineer counts. Other tests check feature aggregation, range ordering, risk labels, spike accounting, and the distinction between effort and duration.

Several validations only assert that an error exists. `TestBookingServiceCaseStudy` asserts High risk and at least two spikes, then logs calculated values; it does not pin every numeric endpoint to an expected constant.

The implementation checks an empty task name with `Name == ""`, so whitespace-only names pass. Floating-point validation does not explicitly reject `NaN` or infinity. It has no dependency graph, waiting-time field, task deduplication, missing-scope detection, or automatic re-estimation after a spike. These are limits of the teaching model, not features silently provided by its result type.

## Communicate the conditions with the number

For the fixture, an honest summary is: eight planning tasks, 2.5 engineer-days of spikes, 24.38–64.38 final engineer-days under the model's 25% contingency, and 34.8–92.0 modeled working days with one engineer at 70% availability. Overall risk is High and confidence is Low. Assumptions include final UI design, available vendor sandbox/documentation, accessible provider API, and no major scope changes.

Vendor credentials are dependencies. Unexpected webhook behavior and changing booking allocation rules are risks. A capacity formula does not automatically add time spent waiting for either.

**An estimate is a range with assumptions, not a promise without conditions.** The model is valuable when its arithmetic and limitations make those conditions visible.

## Repository sources

- [Lab 10](https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/10-project-estimation)
- [README.md](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/10-project-estimation/README.md): scope, breakdown, risks, and communication examples.
- [estimation.go](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/10-project-estimation/estimation.go): types, validation, arithmetic, and heuristics.
- [estimation_test.go](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/10-project-estimation/estimation_test.go): fixtures and assertions.
- [go.mod](https://github.com/lukman-ss/software-engineering-lab/blob/main/labs/10-project-estimation/go.mod): module and Go 1.25 requirement.
