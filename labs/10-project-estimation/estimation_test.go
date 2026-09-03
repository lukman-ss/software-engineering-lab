package estimation_test

import (
	"testing"

	"github.com/lukman-ss/software-engineering-lab/labs/10-project-estimation"
)

func TestNaiveEstimator(t *testing.T) {
	days := estimation.EstimateByPageCount(10)
	if days != 10 {
		t.Errorf("naive estimator: expected 10 days for 10 pages, got %d", days)
	}
}

func TestNaiveEstimatorZeroPages(t *testing.T) {
	days := estimation.EstimateByPageCount(0)
	if days != 0 {
		t.Errorf("naive estimator: expected 0 days for 0 pages, got %d", days)
	}
}

func TestNaiveEstimatorNegativePages(t *testing.T) {
	days := estimation.EstimateByPageCount(-5)
	if days != 0 {
		t.Errorf("naive estimator: expected 0 days for negative pages, got %d", days)
	}
}

func TestSimpleLowRiskProject(t *testing.T) {
	project := estimation.Project{
		Name: "Simple Feature",
		Tasks: []estimation.Task{
			{
				Name: "Login Page",
				Estimate: estimation.EstimateRange{
					Min:        0.5,
					MostLikely: 1.0,
					Max:        2.0,
				},
				Risk: estimation.RiskLow,
			},
		},
		Availability:    0.8,
		EngineerCount:   1,
		AutoContingency: true,
	}

	result, err := project.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Effort.FinalEffort.MinDays < 0.5 {
		t.Errorf("expected min >= 0.5, got %f", result.Effort.FinalEffort.MinDays)
	}
	if result.OverallRisk != estimation.RiskLow {
		t.Errorf("expected Low risk, got %s", result.OverallRisk)
	}
}

func TestMediumRiskProject(t *testing.T) {
	project := estimation.Project{
		Name: "Medium Complexity",
		Tasks: []estimation.Task{
			{Name: "Task 1", Estimate: estimation.EstimateRange{Min: 2, MostLikely: 3, Max: 5}, Risk: estimation.RiskMedium},
			{Name: "Task 2", Estimate: estimation.EstimateRange{Min: 1, MostLikely: 2, Max: 4}, Risk: estimation.RiskMedium},
			{Name: "Task 3", Estimate: estimation.EstimateRange{Min: 3, MostLikely: 4, Max: 6}, Risk: estimation.RiskMedium},
			{Name: "Task 4", Estimate: estimation.EstimateRange{Min: 2, MostLikely: 3, Max: 5}, Risk: estimation.RiskMedium},
		},
		Availability:    0.7,
		EngineerCount:   1,
		AutoContingency: true,
	}

	result, err := project.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.OverallRisk != estimation.RiskMedium {
		t.Errorf("expected Medium risk for 4 medium-risk tasks, got %s", result.OverallRisk)
	}

	expectedMin := 8.0
	if result.Effort.FinalEffort.MinDays < expectedMin {
		t.Errorf("expected final min effort >= %f, got %f", expectedMin, result.Effort.FinalEffort.MinDays)
	}
}

func TestHighRiskExternalAPI(t *testing.T) {
	project := estimation.Project{
		Name: "External API Integration",
		Tasks: []estimation.Task{
			{
				Name: "Payment Gateway",
				Estimate: estimation.EstimateRange{
					Min:        3,
					MostLikely: 5,
					Max:        10,
				},
				Risk: estimation.RiskHigh,
			},
		},
		Availability:    1.0,
		EngineerCount:   1,
		AutoContingency: true,
	}

	result, err := project.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.OverallRisk != estimation.RiskHigh {
		t.Errorf("expected High risk for high-risk task, got %s", result.OverallRisk)
	}
}

func TestUnknownDependencyRequiringSpike(t *testing.T) {
	project := estimation.Project{
		Name: "Vendor API Spike",
		Tasks: []estimation.Task{
			{
				Name: "WhatsApp API Investigation",
				Estimate: estimation.EstimateRange{
					Min:        0,
					MostLikely: 0,
					Max:        0,
				},
				Risk:      estimation.RiskUnknown,
				SpikeDays: 1.0,
			},
			{
				Name: "WhatsApp Integration",
				Estimate: estimation.EstimateRange{
					Min:        2,
					MostLikely: 3,
					Max:        4,
				},
				Risk: estimation.RiskMedium,
			},
		},
		Availability:    1.0,
		EngineerCount:   1,
		AutoContingency: true,
	}

	result, err := project.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.RequiredSpikes) != 1 {
		t.Errorf("expected 1 required spike, got %d", len(result.RequiredSpikes))
	}

	if result.Effort.SpikeEffort.MinDays != 1.0 {
		t.Errorf("expected spike days 1.0, got %f", result.Effort.SpikeEffort.MinDays)
	}
}

func TestContingencyIncreasesEstimate(t *testing.T) {
	project := estimation.Project{
		Name: "High Risk Project",
		Tasks: []estimation.Task{
			{
				Name: "Critical Task",
				Estimate: estimation.EstimateRange{
					Min:        2,
					MostLikely: 4,
					Max:        8,
				},
				Risk: estimation.RiskHigh,
			},
		},
		Availability:    1.0,
		EngineerCount:   1,
		ContingencyRate: 0.15,
		AutoContingency: false,
	}

	result, err := project.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Effort.ContingencyEffort <= 0 {
		t.Error("expected positive contingency days")
	}

	impl := result.Effort.ImplementationEffort
	baseMin := impl.MinDays + result.Effort.SpikeEffort.MinDays
	if result.Effort.FinalEffort.MinDays < baseMin {
		t.Error("final effort should include base efforts and contingency")
	}
}

func TestEffortNotCalendarDuration(t *testing.T) {
	project := estimation.Project{
		Name: "Booking Service",
		Tasks: []estimation.Task{
			{Name: "Feature 1", Estimate: estimation.EstimateRange{Min: 5, MostLikely: 10, Max: 20}, Risk: estimation.RiskMedium},
		},
		Availability:    0.5,
		EngineerCount:   1,
		AutoContingency: true,
	}

	result, err := project.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	effort := result.Effort.FinalEffort.ExpectedDays
	calendar := result.Duration.CalendarDuration.ExpectedDays

	if calendar <= effort {
		t.Errorf("calendar duration (%.1f days) should be greater than effort (%.1f days) with 50%% availability", calendar, effort)
	}
}

func TestMultipleEngineersReducesCalendar(t *testing.T) {
	project := estimation.Project{
		Name: "Team Project",
		Tasks: []estimation.Task{
			{Name: "Task", Estimate: estimation.EstimateRange{Min: 10, MostLikely: 15, Max: 20}, Risk: estimation.RiskLow},
		},
		Availability:    1.0,
		EngineerCount:   2,
		AutoContingency: true,
	}

	result, err := project.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	project2 := project
	project2.EngineerCount = 1

	result2, err := project2.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Duration.CalendarDuration.ExpectedDays >= result2.Duration.CalendarDuration.ExpectedDays {
		t.Errorf("2 engineers should reduce calendar duration vs 1 engineer")
	}
}

func TestFeatureBreakdownAffectsTotal(t *testing.T) {
	project := estimation.Project{
		Name: "Booking Service",
		Features: []estimation.Feature{
			{
				Name: "Booking Flow",
				Tasks: []estimation.Task{
					{Name: "Booking Online", Estimate: estimation.EstimateRange{Min: 2, MostLikely: 3, Max: 5}, Risk: estimation.RiskLow},
					{Name: "Pilih Mekanik", Estimate: estimation.EstimateRange{Min: 1, MostLikely: 2, Max: 3}, Risk: estimation.RiskLow},
				},
			},
			{
				Name: "Payment Flow",
				Tasks: []estimation.Task{
					{Name: "Payment Gateway", Estimate: estimation.EstimateRange{Min: 3, MostLikely: 5, Max: 7}, Risk: estimation.RiskHigh},
				},
			},
		},
		Availability:    0.7,
		EngineerCount:   1,
		AutoContingency: true,
	}

	result, err := project.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedImplMin := 6.0 // 2 + 1 + 3 (Payment Gateway is in Features)
	if result.Effort.ImplementationEffort.MinDays != expectedImplMin {
		t.Errorf("expected implementation min %f, got %f", expectedImplMin, result.Effort.ImplementationEffort.MinDays)
	}

	if result.TotalTaskCount != 3 {
		t.Errorf("expected 3 total tasks, got %d", result.TotalTaskCount)
	}
}

func TestInvalidRangeMinGreaterThanMostLikely(t *testing.T) {
	project := estimation.Project{
		Name: "Invalid Range",
		Tasks: []estimation.Task{
			{
				Name: "Bad Estimate",
				Estimate: estimation.EstimateRange{
					Min:        10,
					MostLikely: 5,
					Max:        15,
				},
				Risk: estimation.RiskLow,
			},
		},
		Availability:  1.0,
		EngineerCount: 1,
	}

	_, err := project.Estimate()
	if err == nil {
		t.Error("expected error for min > mostLikely")
	}
}

func TestInvalidRangeMostLikelyGreaterThanMax(t *testing.T) {
	project := estimation.Project{
		Name: "Invalid Range",
		Tasks: []estimation.Task{
			{
				Name: "Bad Estimate",
				Estimate: estimation.EstimateRange{
					Min:        2,
					MostLikely: 10,
					Max:        5,
				},
				Risk: estimation.RiskLow,
			},
		},
		Availability:  1.0,
		EngineerCount: 1,
	}

	_, err := project.Estimate()
	if err == nil {
		t.Error("expected error for mostLikely > max")
	}
}

func TestNegativeEstimateRejected(t *testing.T) {
	project := estimation.Project{
		Name: "Negative Estimate",
		Tasks: []estimation.Task{
			{
				Name: "Bad Task",
				Estimate: estimation.EstimateRange{
					Min:        -1,
					MostLikely: 1,
					Max:        2,
				},
				Risk: estimation.RiskLow,
			},
		},
		Availability:  1.0,
		EngineerCount: 1,
	}

	_, err := project.Estimate()
	if err == nil {
		t.Error("expected error for negative estimate")
	}
}

func TestNegativeSpikeRejected(t *testing.T) {
	task := estimation.Task{
		Name: "Task",
		Estimate: estimation.EstimateRange{
			Min:        1,
			MostLikely: 2,
			Max:        3,
		},
		Risk:      estimation.RiskLow,
		SpikeDays: -1,
	}

	if err := task.Validate(); err == nil {
		t.Error("expected error for negative spike days")
	}
}

func TestInvalidAvailabilityZero(t *testing.T) {
	project := estimation.Project{
		Name: "Zero Availability",
		Tasks: []estimation.Task{
			{Name: "Task", Estimate: estimation.EstimateRange{Min: 1, MostLikely: 2, Max: 3}, Risk: estimation.RiskLow},
		},
		Availability:  0,
		EngineerCount: 1,
	}

	_, err := project.Estimate()
	if err == nil {
		t.Error("expected error for zero availability")
	}
}

func TestInvalidAvailabilityOverOne(t *testing.T) {
	project := estimation.Project{
		Name: "Over 100% Availability",
		Tasks: []estimation.Task{
			{Name: "Task", Estimate: estimation.EstimateRange{Min: 1, MostLikely: 2, Max: 3}, Risk: estimation.RiskLow},
		},
		Availability:  1.5,
		EngineerCount: 1,
	}

	_, err := project.Estimate()
	if err == nil {
		t.Error("expected error for availability > 1")
	}
}

func TestNegativeContingencyRejected(t *testing.T) {
	project := estimation.Project{
		Name: "Negative Contingency",
		Tasks: []estimation.Task{
			{Name: "Task", Estimate: estimation.EstimateRange{Min: 1, MostLikely: 2, Max: 3}, Risk: estimation.RiskLow},
		},
		Availability:    1.0,
		EngineerCount:   1,
		ContingencyRate: -0.1,
	}

	_, err := project.Estimate()
	if err == nil {
		t.Error("expected error for negative contingency")
	}
}

func TestEmptyProjectRejected(t *testing.T) {
	project := estimation.Project{
		Name:          "Empty Project",
		Tasks:         []estimation.Task{},
		Availability:  1.0,
		EngineerCount: 1,
	}

	_, err := project.Estimate()
	if err == nil {
		t.Error("expected error for empty project")
	}
}

func TestInvalidRiskLevel(t *testing.T) {
	task := estimation.Task{
		Name: "Task",
		Estimate: estimation.EstimateRange{
			Min:        1,
			MostLikely: 2,
			Max:        3,
		},
		Risk: estimation.RiskLevel("Critical"),
	}

	if err := task.Validate(); err == nil {
		t.Error("expected error for invalid risk level 'Critical'")
	}
}

func TestUnknownRiskWithoutSpikeRejected(t *testing.T) {
	task := estimation.Task{
		Name: "Task",
		Estimate: estimation.EstimateRange{
			Min:        1,
			MostLikely: 2,
			Max:        3,
		},
		Risk: estimation.RiskUnknown,
	}

	if err := task.Validate(); err == nil {
		t.Error("expected error for unknown risk without spike")
	}
}

func TestEmptyTaskNameRejected(t *testing.T) {
	task := estimation.Task{
		Name:    "",
		Estimate: estimation.EstimateRange{Min: 1, MostLikely: 2, Max: 3},
		Risk:    estimation.RiskLow,
	}

	if err := task.Validate(); err == nil {
		t.Error("expected error for empty task name")
	}
}

func TestUnknownRiskWithSpikeAndExplicitEstimateAccepted(t *testing.T) {
	task := estimation.Task{
		Name: "Task",
		Estimate: estimation.EstimateRange{
			Min:        1,
			MostLikely: 2,
			Max:        3,
		},
		Risk:      estimation.RiskUnknown,
		SpikeDays: 0.5,
	}

	if err := task.Validate(); err != nil {
		t.Errorf("unexpected error for unknown risk with spike: %v", err)
	}
}

func TestAssumptionsIncreaseConfidence(t *testing.T) {
	project := estimation.Project{
		Name: "With Assumptions",
		Tasks: []estimation.Task{
			{
				Name: "Task",
				Estimate: estimation.EstimateRange{
					Min:        2,
					MostLikely: 3,
					Max:        5,
				},
				Risk:        estimation.RiskLow,
				Assumptions: []string{"UI final", "Vendor API available"},
			},
		},
		Availability:    1.0,
		EngineerCount:   1,
		AutoContingency: true,
	}

	result, err := project.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Confidence != estimation.ConfidenceHigh {
		t.Errorf("expected High confidence with low risk + assumptions, got %s", result.Confidence)
	}
}

func TestMissingAssumptionsReduceConfidence(t *testing.T) {
	projectWithoutAssumptions := estimation.Project{
		Name: "No Assumptions",
		Tasks: []estimation.Task{
			{
				Name: "Task",
				Estimate: estimation.EstimateRange{
					Min:        2,
					MostLikely: 3,
					Max:        5,
				},
				Risk:        estimation.RiskLow,
				Assumptions: []string{},
			},
		},
		Availability:    1.0,
		EngineerCount:   1,
		AutoContingency: true,
	}

	result, err := projectWithoutAssumptions.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Confidence != estimation.ConfidenceLow {
		t.Errorf("expected Low confidence without assumptions, got %s", result.Confidence)
	}
}

func TestFinalEstimateRangeOrder(t *testing.T) {
	project := estimation.Project{
		Name: "Range Validation",
		Tasks: []estimation.Task{
			{
				Name: "Task",
				Estimate: estimation.EstimateRange{
					Min:        5,
					MostLikely: 10,
					Max:        20,
				},
				Risk: estimation.RiskLow,
			},
		},
		Availability:    1.0,
		EngineerCount:   1,
		AutoContingency: true,
	}

	result, err := project.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	effort := result.Effort.FinalEffort
	if effort.MinDays > effort.ExpectedDays {
		t.Errorf("min (%.1f) should be <= expected (%.1f)", effort.MinDays, effort.ExpectedDays)
	}

	if effort.ExpectedDays > effort.MaxDays {
		t.Errorf("expected (%.1f) should be <= max (%.1f)", effort.ExpectedDays, effort.MaxDays)
	}
}

func TestDurationRange(t *testing.T) {
	project := estimation.Project{
		Name: "Duration Test",
		Tasks: []estimation.Task{
			{Name: "Task", Estimate: estimation.EstimateRange{Min: 10, MostLikely: 15, Max: 20}, Risk: estimation.RiskLow},
		},
		Availability:    0.7,
		EngineerCount:   1,
		AutoContingency: false,
	}

	result, err := project.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dur := result.Duration.CalendarDuration
	if dur.MinDays <= 0 || dur.ExpectedDays <= 0 || dur.MaxDays <= 0 {
		t.Error("duration range should be positive")
	}

	if dur.MinDays > dur.ExpectedDays || dur.ExpectedDays > dur.MaxDays {
		t.Errorf("duration range invalid: min(%.1f) <= expected(%.1f) <= max(%.1f)", dur.MinDays, dur.ExpectedDays, dur.MaxDays)
	}
}

func TestInvalidEngineerCount(t *testing.T) {
	project := estimation.Project{
		Name:          "Zero Engineers",
		Tasks:         []estimation.Task{{Name: "Task", Estimate: estimation.EstimateRange{Min: 1, MostLikely: 2, Max: 3}, Risk: estimation.RiskLow}},
		Availability:  1.0,
		EngineerCount: 0,
	}

	_, err := project.Estimate()
	if err == nil {
		t.Error("expected error for zero engineer count")
	}
}

func TestNegativeEngineerCount(t *testing.T) {
	project := estimation.Project{
		Name: "Negative Engineers",
		Tasks: []estimation.Task{
			{Name: "Task", Estimate: estimation.EstimateRange{Min: 1, MostLikely: 2, Max: 3}, Risk: estimation.RiskLow},
		},
		Availability:  1.0,
		EngineerCount: -1,
	}

	_, err := project.Estimate()
	if err == nil {
		t.Error("expected error for negative engineer count")
	}
}

func TestBookingServiceCaseStudy(t *testing.T) {
	project := estimation.Project{
		Name: "Aplikasi Booking Servis",
		Tasks: []estimation.Task{
			{
				Name: "Login",
				Estimate: estimation.EstimateRange{
					Min:        1,
					MostLikely: 2,
					Max:        3,
				},
				Risk: estimation.RiskLow,
			},
			{
				Name: "Booking Online",
				Estimate: estimation.EstimateRange{
					Min:        3,
					MostLikely: 5,
					Max:        8,
				},
				Risk: estimation.RiskMedium,
			},
			{
				Name: "Pilih Cabang",
				Estimate: estimation.EstimateRange{
					Min:        1,
					MostLikely: 2,
					Max:        4,
				},
				Risk: estimation.RiskLow,
			},
			{
				Name: "Pilih Mekanik",
				Estimate: estimation.EstimateRange{
					Min:        2,
					MostLikely: 3,
					Max:        5,
				},
				Risk: estimation.RiskMedium,
			},
			{
				Name: "Payment Gateway",
				Estimate: estimation.EstimateRange{
					Min:        4,
					MostLikely: 6,
					Max:        12,
				},
				Risk:        estimation.RiskUnknown,
				SpikeDays:   2,
				Assumptions: []string{"Vendor sandbox available", "Integration docs complete"},
			},
			{
				Name: "WhatsApp Notification",
				Estimate: estimation.EstimateRange{
					Min:        1,
					MostLikely: 2,
					Max:        4,
				},
				Risk:        estimation.RiskUnknown,
				SpikeDays:   0.5,
				Assumptions: []string{"Provider API accessible"},
			},
		},
		Features: []estimation.Feature{
			{Name: "Admin Dashboard", Tasks: []estimation.Task{
				{Name: "Dashboard Admin", Estimate: estimation.EstimateRange{Min: 3, MostLikely: 5, Max: 7}, Risk: estimation.RiskMedium},
				{Name: "Laporan Excel", Estimate: estimation.EstimateRange{Min: 2, MostLikely: 4, Max: 6}, Risk: estimation.RiskHigh},
			}},
		},
		Availability:    0.7,
		EngineerCount:   1,
		AutoContingency: true,
		Assumptions:     []string{"UI design final", "1 engineer available", "No major requirement changes"},
	}

	result, err := project.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.OverallRisk != estimation.RiskHigh {
		t.Errorf("expected High risk due to unknown dependencies, got %s", result.OverallRisk)
	}

	if len(result.RequiredSpikes) < 2 {
		t.Errorf("expected at least 2 required spikes, got %d", len(result.RequiredSpikes))
	}

	t.Logf("Project: %s", result.ProjectName)
	t.Logf("Implementation Effort: %s", result.Effort.ImplementationEffort)
	t.Logf("Spike Effort: %.1f days", result.Effort.SpikeEffort.ExpectedDays)
	t.Logf("Base Effort: %s", result.Effort.BaseEffort)
	t.Logf("Contingency: %.1f days", result.Effort.ContingencyEffort)
	t.Logf("Final Effort: %s", result.Effort.FinalEffort)
	t.Logf("Calendar Duration: %s", result.Duration.CalendarDuration)
	t.Logf("Risk: %s, Confidence: %s", result.OverallRisk, result.Confidence)
}
