package estimation

import (
	"errors"
	"fmt"
	"math"
)

var (
	ErrInvalidRange        = errors.New("invalid estimate range: min must be <= most likely and most likely <= max")
	ErrNegativeEstimate    = errors.New("estimate values cannot be negative")
	ErrInvalidAvailability = errors.New("availability must be between 0.0 (exclusive) and 1.0 (inclusive)")
	ErrNegativeContingency = errors.New("contingency rate cannot be negative")
	ErrEmptyProject        = errors.New("project must have at least one task or feature")
	ErrUnknownRiskNoSpike  = errors.New("task with unknown risk cannot have zero effort without a spike or explicit range")
)

// RiskLevel represents the technical uncertainty or complexity risk of a task.
type RiskLevel string

const (
	RiskLow     RiskLevel = "Low"
	RiskMedium  RiskLevel = "Medium"
	RiskHigh    RiskLevel = "High"
	RiskUnknown RiskLevel = "Unknown"
)

// ConfidenceLevel represents the overall project estimation confidence.
type ConfidenceLevel string

const (
	ConfidenceHigh   ConfidenceLevel = "High"
	ConfidenceMedium ConfidenceLevel = "Medium"
	ConfidenceLow    ConfidenceLevel = "Low"
)

// EstimateRange models effort range using 3-point estimation (Min, MostLikely, Max).
type EstimateRange struct {
	Min        float64
	MostLikely float64
	Max        float64
}

// Validate checks domain invariants for estimate ranges.
func (r EstimateRange) Validate() error {
	if r.Min < 0 || r.MostLikely < 0 || r.Max < 0 {
		return ErrNegativeEstimate
	}
	if r.Min > r.MostLikely || r.MostLikely > r.Max {
		return fmt.Errorf("%w: min(%.2f) <= mostLikely(%.2f) <= max(%.2f)", ErrInvalidRange, r.Min, r.MostLikely, r.Max)
	}
	return nil
}

// Expected returns the PERT expected effort: (Min + 4*MostLikely + Max) / 6.
func (r EstimateRange) Expected() float64 {
	return (r.Min + 4.0*r.MostLikely + r.Max) / 6.0
}

// Task represents a discrete piece of work.
type Task struct {
	Name        string
	Estimate    EstimateRange
	Risk        RiskLevel
	SpikeDays   float64
	Assumptions []string
}

// Validate ensures task data invariants are met.
func (t Task) Validate() error {
	if err := t.Estimate.Validate(); err != nil {
		return fmt.Errorf("task %q: %w", t.Name, err)
	}
	if t.SpikeDays < 0 {
		return fmt.Errorf("task %q: spike days cannot be negative", t.Name)
	}
	if t.Risk == RiskUnknown && t.SpikeDays == 0 && t.Estimate.Max == 0 {
		return fmt.Errorf("task %q: %w", t.Name, ErrUnknownRiskNoSpike)
	}
	return nil
}

// Feature groups multiple tasks across activities (Backend, Frontend, Testing, etc.).
type Feature struct {
	Name  string
	Tasks []Task
}

// Project represents the full project scope with estimation parameters.
type Project struct {
	Name             string
	Features         []Feature
	Tasks            []Task
	EngineerCount    int
	Availability     float64 // e.g., 0.70 for 70% productive project allocation
	ContingencyRate  float64 // e.g., 0.15 for 15% contingency buffer (0 = auto-calculated from risk)
	AutoContingency  bool    // if true, calculates contingency dynamically based on project risk profile
	Assumptions      []string
}

// EffortBreakdown summarizes calculated effort across ranges.
type EffortBreakdown struct {
	MinDays        float64
	ExpectedDays   float64
	MaxDays        float64
	SpikeDays      float64
	BaseEffort     float64 // expected implementation + spike
	ContingencyDays float64
	TotalEffortDays float64 // BaseEffort + ContingencyDays
}

// DurationBreakdown summarizes calendar time based on availability and engineering capacity.
type DurationBreakdown struct {
	CalendarDays  float64
	CalendarWeeks float64
}

// EstimationResult is the comprehensive result of structured project estimation.
type EstimationResult struct {
	ProjectName       string
	Effort            EffortBreakdown
	Duration          DurationBreakdown
	OverallRisk       RiskLevel
	Confidence        ConfidenceLevel
	RequiredSpikes    []string
	AllAssumptions    []string
	TotalTaskCount    int
}

// NaiveEstimator implements a simplistic, single-point calculation for demonstration.
// Junior mental model: 10 pages * 1 day = 10 days (ignores risk, uncertainty, activities, availability).
func EstimateByPageCount(pageCount int) int {
	return pageCount * 1
}

// Estimate calculates a comprehensive, risk-aware project estimate.
func (p Project) Estimate() (*EstimationResult, error) {
	allTasks := make([]Task, 0, len(p.Tasks))
	allTasks = append(allTasks, p.Tasks...)
	for _, f := range p.Features {
		allTasks = append(allTasks, f.Tasks...)
	}

	if len(allTasks) == 0 {
		return nil, ErrEmptyProject
	}

	if p.Availability <= 0 || p.Availability > 1.0 {
		return nil, ErrInvalidAvailability
	}

	if p.ContingencyRate < 0 {
		return nil, ErrNegativeContingency
	}

	engineerCount := p.EngineerCount
	if engineerCount <= 0 {
		engineerCount = 1
	}

	var (
		minSum        float64
		expectedSum   float64
		maxSum        float64
		spikeSum      float64
		highRiskCount int
		unknownCount  int
		requiredSpikes []string
		assumptions   []string
	)

	assumptions = append(assumptions, p.Assumptions...)

	for _, task := range allTasks {
		if err := task.Validate(); err != nil {
			return nil, err
		}

		minSum += task.Estimate.Min
		expectedSum += task.Estimate.Expected()
		maxSum += task.Estimate.Max
		spikeSum += task.SpikeDays

		if task.SpikeDays > 0 {
			requiredSpikes = append(requiredSpikes, fmt.Sprintf("%s (%.1f days spike)", task.Name, task.SpikeDays))
		}

		assumptions = append(assumptions, task.Assumptions...)

		switch task.Risk {
		case RiskHigh:
			highRiskCount++
		case RiskUnknown:
			unknownCount++
		}
	}

	// Calculate overall project risk
	overallRisk := RiskLow
	totalTasks := len(allTasks)
	if unknownCount > 0 || float64(highRiskCount)/float64(totalTasks) >= 0.3 {
		overallRisk = RiskHigh
	} else if highRiskCount > 0 || float64(totalTasks) > 10 {
		overallRisk = RiskMedium
	}

	// Calculate Contingency
	contingencyRate := p.ContingencyRate
	if p.AutoContingency || contingencyRate == 0 {
		switch overallRisk {
		case RiskHigh:
			contingencyRate = 0.25 // 25% for high risk / unknown projects
		case RiskMedium:
			contingencyRate = 0.15 // 15% standard buffer
		case RiskLow:
			contingencyRate = 0.10 // 10% minimal buffer
		}
	}

	baseEffort := expectedSum + spikeSum
	contingencyDays := baseEffort * contingencyRate
	totalEffortDays := baseEffort + contingencyDays

	// Calendar duration = Total Effort / (Engineers * Availability)
	effectiveDailyCapacity := float64(engineerCount) * p.Availability
	calendarDays := totalEffortDays / effectiveDailyCapacity
	calendarWeeks := calendarDays / 5.0 // 5 working days per week

	// Determine confidence level
	confidence := ConfidenceHigh
	if len(assumptions) == 0 || unknownCount > 0 || overallRisk == RiskHigh {
		confidence = ConfidenceLow
	} else if overallRisk == RiskMedium || contingencyRate > 0.15 {
		confidence = ConfidenceMedium
	}

	return &EstimationResult{
		ProjectName: p.Name,
		Effort: EffortBreakdown{
			MinDays:         math.Round(minSum*100) / 100,
			ExpectedDays:    math.Round(expectedSum*100) / 100,
			MaxDays:         math.Round(maxSum*100) / 100,
			SpikeDays:       math.Round(spikeSum*100) / 100,
			BaseEffort:      math.Round(baseEffort*100) / 100,
			ContingencyDays: math.Round(contingencyDays*100) / 100,
			TotalEffortDays: math.Round(totalEffortDays*100) / 100,
		},
		Duration: DurationBreakdown{
			CalendarDays:  math.Round(calendarDays*10) / 10,
			CalendarWeeks: math.Round(calendarWeeks*10) / 10,
		},
		OverallRisk:    overallRisk,
		Confidence:     confidence,
		RequiredSpikes: requiredSpikes,
		AllAssumptions: assumptions,
		TotalTaskCount: totalTasks,
	}, nil
}
