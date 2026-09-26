package adr_test

import (
	"strings"
	"testing"

	"labs/17-architecture-decision-record/internal/adr"
)

func TestLinter_ValidSequence(t *testing.T) {
	adr1 := &adr.Record{
		ID:           1,
		Title:        "Modular Monolith",
		Status:       adr.StatusSuperseded,
		SupersededBy: 2,
	}
	adr2 := &adr.Record{
		ID:         2,
		Title:      "Microservices",
		Status:     adr.StatusAccepted,
		Supersedes: 1,
	}
	adr3 := &adr.Record{
		ID:     3,
		Title:  "Event Sourcing",
		Status: adr.StatusRejected,
	}

	linter := adr.NewLinter()
	errs := linter.Validate([]*adr.Record{adr1, adr2, adr3})
	if len(errs) > 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestLinter_BrokenReferences(t *testing.T) {
	tests := []struct {
		name    string
		records []*adr.Record
		wantErr string
	}{
		{
			name: "superseded points to non-existent",
			records: []*adr.Record{
				{
					ID:           1,
					Title:        "Modular Monolith",
					Status:       adr.StatusSuperseded,
					SupersededBy: 99,
				},
			},
			wantErr: "superseded by non-existent ADR 99",
		},
		{
			name: "supersedes points to non-existent",
			records: []*adr.Record{
				{
					ID:         1,
					Title:      "Microservices",
					Status:     adr.StatusAccepted,
					Supersedes: 99,
				},
			},
			wantErr: "supersedes non-existent ADR 99",
		},
		{
			name: "mismatched supersession link",
			records: []*adr.Record{
				{
					ID:           1,
					Title:        "Modular Monolith",
					Status:       adr.StatusSuperseded,
					SupersededBy: 2,
				},
				{
					ID:     2,
					Title:  "Microservices",
					Status: adr.StatusAccepted,
					// Forgot Supersedes: 1
				},
			},
			wantErr: "ADR 1 superseded by ADR 2, but ADR 2 does not declare it supersedes ADR 1",
		},
		{
			name: "non-monotonic numbering",
			records: []*adr.Record{
				{
					ID:     1,
					Title:  "One",
					Status: adr.StatusAccepted,
				},
				{
					ID:     3,
					Title:  "Three",
					Status: adr.StatusAccepted,
				},
			},
			wantErr: "non-monotonic numbering, expected 2 but got 3",
		},
		{
			name: "self supersession",
			records: []*adr.Record{
				{
					ID:           1,
					Title:        "Self",
					Status:       adr.StatusSuperseded,
					SupersededBy: 1,
				},
			},
			wantErr: "ADR 1 cannot supersede itself",
		},
		{
			name: "duplicate ADR ID",
			records: []*adr.Record{
				{
					ID:     1,
					Title:  "First",
					Status: adr.StatusAccepted,
				},
				{
					ID:     1,
					Title:  "Duplicate First",
					Status: adr.StatusAccepted,
				},
			},
			wantErr: "duplicate ADR ID: 1",
		},
		{
			name: "cyclical supersession",
			records: []*adr.Record{
				{
					ID:           1,
					Title:        "One",
					Status:       adr.StatusSuperseded,
					SupersededBy: 2,
				},
				{
					ID:           2,
					Title:        "Two",
					Status:       adr.StatusSuperseded,
					Supersedes:   1,
					SupersededBy: 3,
				},
				{
					ID:           3,
					Title:        "Three",
					Status:       adr.StatusSuperseded,
					Supersedes:   2,
					SupersededBy: 1,
				},
			},
			wantErr: "cyclical supersession detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			linter := adr.NewLinter()
			errs := linter.Validate(tt.records)
			if len(errs) == 0 {
				t.Fatalf("expected errors, got none")
			}

			found := false
			for _, err := range errs {
				if strings.Contains(err.Error(), tt.wantErr) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected error containing %q, got %v", tt.wantErr, errs)
			}
		})
	}
}

func TestLinter_ConcurrencyStress(t *testing.T) {
	const count = 100
	records := make([]*adr.Record, count)

	for i := 1; i <= count; i += 2 {
		records[i-1] = &adr.Record{
			ID:           i,
			Title:        "Decision Old",
			Status:       adr.StatusSuperseded,
			SupersededBy: i + 1,
		}
		records[i] = &adr.Record{
			ID:         i + 1,
			Title:      "Decision New",
			Status:     adr.StatusAccepted,
			Supersedes: i,
		}
	}

	linter := adr.NewLinter()
	errs := linter.Validate(records)
	if len(errs) > 0 {
		t.Fatalf("expected 0 errors for stress batch, got %d: %v", len(errs), errs)
	}
}

