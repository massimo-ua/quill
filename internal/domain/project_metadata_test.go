package domain

import (
	"strings"
	"testing"
	"time"
)

func TestNewProjectMetadata(t *testing.T) {
	now := time.Now()
	future := now.AddDate(0, 1, 0)

	tests := []struct {
		name        string
		projName    string
		desc        string
		goals       []string
		kpis        []string
		start       time.Time
		end         time.Time
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid project metadata",
			projName:    "Test Project",
			desc:        "A test project",
			goals:       []string{"Goal 1", "Goal 2"},
			kpis:        []string{"KPI 1", "KPI 2"},
			start:       now,
			end:         future,
			expectError: false,
		},
		{
			name:        "missing project name",
			projName:    "",
			desc:        "A test project",
			goals:       []string{"Goal 1"},
			kpis:        []string{"KPI 1"},
			start:       now,
			end:         future,
			expectError: true,
			errorMsg:    "project name is required",
		},
		{
			name:        "missing description",
			projName:    "Test Project",
			desc:        "",
			goals:       []string{"Goal 1"},
			kpis:        []string{"KPI 1"},
			start:       now,
			end:         future,
			expectError: true,
			errorMsg:    "project description is required",
		},
		{
			name:        "empty business goals",
			projName:    "Test Project",
			desc:        "A test project",
			goals:       []string{},
			kpis:        []string{"KPI 1"},
			start:       now,
			end:         future,
			expectError: true,
			errorMsg:    "at least one business goal is required",
		},
		{
			name:        "empty KPIs",
			projName:    "Test Project",
			desc:        "A test project",
			goals:       []string{"Goal 1"},
			kpis:        []string{},
			start:       now,
			end:         future,
			expectError: true,
			errorMsg:    "at least one KPI is required",
		},
		{
			name:        "end date before start date",
			projName:    "Test Project",
			desc:        "A test project",
			goals:       []string{"Goal 1"},
			kpis:        []string{"KPI 1"},
			start:       future,
			end:         now,
			expectError: true,
			errorMsg:    "start date cannot be after end date",
		},
		{
			name:        "zero start date",
			projName:    "Test Project",
			desc:        "A test project",
			goals:       []string{"Goal 1"},
			kpis:        []string{"KPI 1"},
			start:       time.Time{},
			end:         future,
			expectError: true,
			errorMsg:    "start and end dates must be set",
		},
		{
			name:        "zero end date",
			projName:    "Test Project",
			desc:        "A test project",
			goals:       []string{"Goal 1"},
			kpis:        []string{"KPI 1"},
			start:       now,
			end:         time.Time{},
			expectError: true,
			errorMsg:    "start and end dates must be set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm, err := NewProjectMetadata(tt.projName, tt.desc, tt.goals, tt.kpis, tt.start, tt.end)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if pm.Name != tt.projName {
				t.Errorf("expected name %q, got %q", tt.projName, pm.Name)
			}
			if pm.Description != tt.desc {
				t.Errorf("expected description %q, got %q", tt.desc, pm.Description)
			}
			if len(pm.BusinessGoals) != len(tt.goals) {
				t.Errorf("expected %d goals, got %d", len(tt.goals), len(pm.BusinessGoals))
			}
			if len(pm.KPIs) != len(tt.kpis) {
				t.Errorf("expected %d KPIs, got %d", len(tt.kpis), len(pm.KPIs))
			}
		})
	}
}

func TestMustNewProjectMetadata(t *testing.T) {
	now := time.Now()
	future := now.AddDate(0, 1, 0)

	tests := []struct {
		name        string
		projName    string
		desc        string
		goals       []string
		kpis        []string
		start       time.Time
		end         time.Time
		shouldPanic bool
	}{
		{
			name:        "valid project metadata",
			projName:    "Test Project",
			desc:        "A test project",
			goals:       []string{"Goal 1"},
			kpis:        []string{"KPI 1"},
			start:       now,
			end:         future,
			shouldPanic: false,
		},
		{
			name:        "invalid project metadata",
			projName:    "",
			desc:        "A test project",
			goals:       []string{"Goal 1"},
			kpis:        []string{"KPI 1"},
			start:       now,
			end:         future,
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tt.shouldPanic && r == nil {
					t.Error("expected panic but got none")
				}
				if !tt.shouldPanic && r != nil {
					t.Errorf("unexpected panic: %v", r)
				}
			}()

			pm := MustNewProjectMetadata(tt.projName, tt.desc, tt.goals, tt.kpis, tt.start, tt.end)
			if tt.shouldPanic {
				return
			}

			if pm.Name != tt.projName {
				t.Errorf("expected name %q, got %q", tt.projName, pm.Name)
			}
		})
	}
}

func TestFormatContext(t *testing.T) {
	now := time.Now()
	future := now.AddDate(0, 1, 0)

	pm := MustNewProjectMetadata(
		"Test Project",
		"A test project description",
		[]string{"Goal 1", "Goal 2"},
		[]string{"KPI 1", "KPI 2"},
		now,
		future,
	)

	formatted := pm.FormatContext()

	expectedParts := []string{
		"Project: Test Project",
		"Description: A test project description",
		"Business Goals:",
		"- Goal 1",
		"- Goal 2",
		"Key Performance Indicators:",
		"- KPI 1",
		"- KPI 2",
		"Project Timeline:",
		now.Format("2006-01-02"),
		future.Format("2006-01-02"),
	}

	for _, part := range expectedParts {
		if !strings.Contains(formatted, part) {
			t.Errorf("expected formatted context to contain %q, but it doesn't\nGot: %s", part, formatted)
		}
	}
}

func TestProjectMetadata_Validate(t *testing.T) {
	now := time.Now()
	future := now.AddDate(0, 1, 0)

	validMetadata := &ProjectMetadata{
		Name:          "Test Project",
		Description:   "A test project",
		BusinessGoals: []string{"Goal 1"},
		KPIs:          []string{"KPI 1"},
		StartDate:     now,
		EndDate:       future,
	}

	t.Run("valid metadata", func(t *testing.T) {
		if err := validMetadata.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("validate copy doesn't affect original", func(t *testing.T) {
		metadata := *validMetadata
		metadata.Name = ""

		// This should fail
		err := metadata.Validate()
		if err == nil {
			t.Error("expected error but got none")
		}

		// Original should still be valid
		if err := validMetadata.Validate(); err != nil {
			t.Errorf("original metadata should still be valid, got error: %v", err)
		}
	})
}
