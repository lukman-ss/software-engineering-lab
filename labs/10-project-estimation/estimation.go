package estimation

import (
	"errors"
	"fmt"
	"math"
)

var (
	ErrInvalidRange         = errors.New("invalid estimate range: min must be <= most likely and most likely <= max")
	ErrNegativeEstimate     = errors.New("estimate values cannot be negative")
	ErrInvalidAvailability  = errors.New("availability must be between 0.0 (exclusive) and 1.0 (inclusive)")
	ErrNegativeContingency  = errors.New("contingency rate cannot be negative")
	ErrEmptyProject         = errors.New("project must have at least one task or feature")
	ErrUnknownRiskNoSpike   = errors.New("task with unknown risk must have SpikeDays > 0")
	ErrInvalidRiskLevel     = errors.New("invalid risk level")
	ErrInvalidEngineerCount = errors.New("engineer count must be at least 1")
	ErrEmptyTaskName        = errors.New("task name cannot be empty")
)

type RiskLevel string

const (
	RiskLow     RiskLevel = "Low"
	RiskMedium  RiskLevel = "Medium"
	RiskHigh    RiskLevel = "High"
	RiskUnknown RiskLevel = "Unknown"
)

type ConfidenceLevel string

const (
	ConfidenceHigh   ConfidenceLevel = "High"
	ConfidenceMedium ConfidenceLevel = "Medium"
	ConfidenceLow    ConfidenceLevel = "Low"
)

type EstimateRange struct {
	Min        float64
	MostLikely float64
	Max        float64
}

func (r EstimateRange) Validate() error {
	if r.Min < 0 || r.MostLikely < 0 || r.Max < 0 {
		return ErrNegativeEstimate
	}
	if r.Min > r.MostLikely || r.MostLikely > r.Max {
		return fmt.Errorf("%w: min(%.2f) <= mostLikely(%.2f) <= max(%.2f)", ErrInvalidRange, r.Min, r.MostLikely, r.Max)
	}
	return nil
}

func (r EstimateRange) Expected() float64 {
	return (r.Min + 4.0*r.MostLikely + r.Max) / 6.0
}

type Task struct {
	Name        string
	Estimate    EstimateRange
	Risk        RiskLevel
	SpikeDays   float64
	Assumptions []string
}

func (t Task) Validate() error {
	if t.Name == "" {
		return ErrEmptyTaskName
	}
	if err := t.Estimate.Validate(); err != nil {
		return fmt.Errorf("task %q: %w", t.Name, err)
	}
	if t.SpikeDays < 0 {
		return fmt.Errorf("task %q: spike days cannot be negative", t.Name)
	}
	if t.Risk != RiskLow && t.Risk != RiskMedium && t.Risk != RiskHigh && t.Risk != RiskUnknown {
		return fmt.Errorf("task %q: %w", t.Name, ErrInvalidRiskLevel)
	}
	if t.Risk == RiskUnknown && t.SpikeDays == 0 {
		return fmt.Errorf("task %q: %w", t.Name, ErrUnknownRiskNoSpike)
	}
	return nil
}

type Feature struct {
	Name  string
	Tasks []Task
}

type Project struct {
	Name            string
	Features        []Feature
	Tasks           []Task
	EngineerCount   int
	Availability    float64
	ContingencyRate float64
	AutoContingency bool
	Assumptions     []string
}

type EffortRange struct {
	MinDays      float64
	ExpectedDays float64
	MaxDays      float64
}

func (r EffortRange) Validate() error {
	if r.MinDays > r.ExpectedDays || r.ExpectedDays > r.MaxDays {
		return errors.New("effort range validation error: min <= expected <= max")
	}
	return nil
}

func (r EffortRange) String() string {
	return fmt.Sprintf("%.1f–%.1f engineer-days", r.MinDays, r.MaxDays)
}

type DurationRange struct {
	MinDays      float64
	ExpectedDays float64
	MaxDays      float64
}

func (r DurationRange) Weeks() (min, expected, max float64) {
	return r.MinDays / 5.0, r.ExpectedDays / 5.0, r.MaxDays / 5.0
}

func (r DurationRange) String() string {
	minW, _, maxW := r.Weeks()
	return fmt.Sprintf("%.0f–%.0f working days (%.1f–%.1f weeks)", r.MinDays, r.MaxDays, minW, maxW)
}

type EffortBreakdown struct {
	ImplementationEffort EffortRange
	SpikeEffort          EffortRange
	BaseEffort           EffortRange
	ContingencyEffort    float64
	FinalEffort          EffortRange
}

type DurationBreakdown struct {
	CalendarDuration DurationRange
}

type EstimationResult struct {
	ProjectName    string
	Effort         EffortBreakdown
	Duration       DurationBreakdown
	OverallRisk    RiskLevel
	Confidence     ConfidenceLevel
	RequiredSpikes []SpikeInfo
	AllAssumptions []string
	TotalTaskCount int
	EngineerCount  int
	Availability   float64
}

type SpikeInfo struct {
	TaskName string
	Days     float64
}

func (s SpikeInfo) String() string {
	return fmt.Sprintf("%s (%.1f days)", s.TaskName, s.Days)
}

func EstimateByPageCount(pageCount int) int {
	if pageCount <= 0 {
		return 0
	}
	return pageCount
}

func (p Project) Estimate() (*EstimationResult, error) {
	allTasks := make([]Task, 0, len(p.Tasks)+len(p.Features)*5)
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

	if p.EngineerCount < 1 {
		return nil, ErrInvalidEngineerCount
	}

	var (
		implMinSum      float64
		implExpectedSum float64
		implMaxSum      float64
		spikeSum        float64
		highRiskCount   int
		mediumRiskCount int
		unknownCount    int
		requiredSpikes  []SpikeInfo
		assumptions     []string
	)

	assumptions = append(assumptions, p.Assumptions...)

	for _, task := range allTasks {
		if err := task.Validate(); err != nil {
			return nil, err
		}

		implMinSum += task.Estimate.Min
		implExpectedSum += task.Estimate.Expected()
		implMaxSum += task.Estimate.Max
		spikeSum += task.SpikeDays

		if task.SpikeDays > 0 {
			requiredSpikes = append(requiredSpikes, SpikeInfo{
				TaskName: task.Name,
				Days:     task.SpikeDays,
			})
		}

		assumptions = append(assumptions, task.Assumptions...)
		switch task.Risk {
		case RiskHigh:
			highRiskCount++
		case RiskMedium:
			mediumRiskCount++
		case RiskUnknown:
			unknownCount++
		}
	}

	overallRisk := RiskLow
	totalTasks := len(allTasks)
	if unknownCount > 0 {
		overallRisk = RiskHigh
	} else if totalTasks > 0 && float64(highRiskCount)/float64(totalTasks) >= 0.3 {
		overallRisk = RiskHigh
	} else if highRiskCount > 0 || mediumRiskCount > 0 {
		overallRisk = RiskMedium
	}

	contingencyRate := p.ContingencyRate
	if p.AutoContingency && contingencyRate == 0 {
		switch overallRisk {
		case RiskHigh:
			contingencyRate = 0.25
		case RiskMedium:
			contingencyRate = 0.15
		case RiskLow:
			contingencyRate = 0.10
		}
	}

	implRange := EffortRange{
		MinDays:      math.Round(implMinSum*100) / 100,
		ExpectedDays: math.Round(implExpectedSum*100) / 100,
		MaxDays:      math.Round(implMaxSum*100) / 100,
	}

	spikeRange := EffortRange{
		MinDays:      spikeSum,
		ExpectedDays: spikeSum,
		MaxDays:      spikeSum,
	}

	baseRange := EffortRange{
		MinDays:      math.Round((implRange.MinDays + spikeRange.MinDays) * 100) / 100,
		ExpectedDays: math.Round((implRange.ExpectedDays + spikeRange.ExpectedDays) * 100) / 100,
		MaxDays:      math.Round((implRange.MaxDays + spikeRange.MaxDays) * 100) / 100,
	}

	contingencyDays := baseRange.ExpectedDays * contingencyRate

	finalRange := EffortRange{
		MinDays:      math.Round((baseRange.MinDays * (1 + contingencyRate)) * 100) / 100,
		ExpectedDays: math.Round((baseRange.ExpectedDays + contingencyDays) * 100) / 100,
		MaxDays:      math.Round((baseRange.MaxDays * (1 + contingencyRate)) * 100) / 100,
	}

	effectiveDailyCapacity := float64(p.EngineerCount) * p.Availability
	calendarDuration := DurationRange{
		MinDays:      math.Round(finalRange.MinDays/effectiveDailyCapacity*10) / 10,
		ExpectedDays: math.Round(finalRange.ExpectedDays/effectiveDailyCapacity*10) / 10,
		MaxDays:      math.Round(finalRange.MaxDays/effectiveDailyCapacity*10) / 10,
	}

	confidence := calculateConfidence(overallRisk, unknownCount, len(assumptions), contingencyRate)

	return &EstimationResult{
		ProjectName:    p.Name,
		Effort: EffortBreakdown{
			ImplementationEffort: implRange,
			SpikeEffort:          spikeRange,
			BaseEffort:           baseRange,
			ContingencyEffort:    contingencyDays,
			FinalEffort:          finalRange,
		},
		Duration: DurationBreakdown{
			CalendarDuration: calendarDuration,
		},
		OverallRisk:    overallRisk,
		Confidence:     confidence,
		RequiredSpikes: requiredSpikes,
		AllAssumptions: assumptions,
		TotalTaskCount: totalTasks,
		EngineerCount:  p.EngineerCount,
		Availability:   p.Availability,
	}, nil
}

func calculateConfidence(overallRisk RiskLevel, unknownCount, assumptionCount int, contingencyRate float64) ConfidenceLevel {
	if unknownCount > 0 || overallRisk == RiskHigh {
		return ConfidenceLow
	}
	if assumptionCount == 0 {
		if overallRisk == RiskMedium {
			return ConfidenceMedium
		}
		return ConfidenceLow
	}
	if overallRisk == RiskMedium {
		return ConfidenceMedium
	}
	return ConfidenceHigh
}