package adr_test

import (
	"strings"
	"testing"

	"labs/17-architecture-decision-record/internal/adr"
)

func TestParse_Valid(t *testing.T) {
	content := `
# 1. Modular Monolith over Microservices

Date: 2026-09-25
Status: Accepted

## Context
Small team of 5 engineers.

## Decision
Modular Monolith.
`
	record, err := adr.Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if record.ID != 1 {
		t.Errorf("expected id 1, got %d", record.ID)
	}
	if record.Status != adr.StatusAccepted {
		t.Errorf("expected status Accepted, got %s", record.Status)
	}
}

func TestParse_Superseded(t *testing.T) {
	content := `
# 1. Modular Monolith over Microservices

Date: 2026-09-25
Status: Superseded by 2

## Context
Small team of 5 engineers.
`
	record, err := adr.Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if record.ID != 1 {
		t.Errorf("expected id 1, got %d", record.ID)
	}
	if record.Status != adr.StatusSuperseded {
		t.Errorf("expected status Superseded, got %s", record.Status)
	}
	if record.SupersededBy != 2 {
		t.Errorf("expected superseded by 2, got %d", record.SupersededBy)
	}
}

func TestParse_Supersedes(t *testing.T) {
	content := `
# 2. Extract Notification Microservice

Date: 2026-10-01
Status: Accepted
Supersedes: 1

## Context
Team scale exceeded 20 engineers.
`
	record, err := adr.Parse(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if record.ID != 2 {
		t.Errorf("expected id 2, got %d", record.ID)
	}
	if record.Status != adr.StatusAccepted {
		t.Errorf("expected status Accepted, got %s", record.Status)
	}
	if record.Supersedes != 1 {
		t.Errorf("expected supersedes 1, got %d", record.Supersedes)
	}
}

func TestParse_Invalid(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr string
	}{
		{
			name: "missing title",
			content: `
Status: Accepted
`,
			wantErr: "title not found",
		},
		{
			name: "missing status",
			content: `
# 1. Title
`,
			wantErr: "status not found",
		},
		{
			name: "invalid status",
			content: `
# 1. Title
Status: Draft
`,
			wantErr: "invalid status: Draft",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := adr.Parse(tt.content)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("expected error containing %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}
