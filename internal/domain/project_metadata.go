package domain

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidProjectMetadata = errors.New("invalid project metadata")
)

// ProjectMetadata represents essential project information and goals
type ProjectMetadata struct {
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	BusinessGoals []string  `json:"businessGoals"`
	KPIs          []string  `json:"kpis"`
	StartDate     time.Time `json:"startDate"`
	EndDate       time.Time `json:"endDate"`
}

// NewProjectMetadata creates a new ProjectMetadata instance with validation
func NewProjectMetadata(name, description string, goals, kpis []string, start, end time.Time) (*ProjectMetadata, error) {
	pm := &ProjectMetadata{
		Name:          name,
		Description:   description,
		BusinessGoals: goals,
		KPIs:          kpis,
		StartDate:     start,
		EndDate:       end,
	}

	if err := pm.Validate(); err != nil {
		return nil, fmt.Errorf("failed to create project metadata: %w", err)
	}

	return pm, nil
}

// MustNewProjectMetadata creates a new ProjectMetadata instance and panics if validation fails
func MustNewProjectMetadata(name, description string, goals, kpis []string, start, end time.Time) *ProjectMetadata {
	pm, err := NewProjectMetadata(name, description, goals, kpis, start, end)
	if err != nil {
		panic(err)
	}
	return pm
}

// Validate ensures all project metadata fields are valid
func (pm *ProjectMetadata) Validate() error {
	if pm.Name == "" {
		return fmt.Errorf("%w: project name is required", ErrInvalidProjectMetadata)
	}

	if pm.Description == "" {
		return fmt.Errorf("%w: project description is required", ErrInvalidProjectMetadata)
	}

	if len(pm.BusinessGoals) == 0 {
		return fmt.Errorf("%w: at least one business goal is required", ErrInvalidProjectMetadata)
	}

	if len(pm.KPIs) == 0 {
		return fmt.Errorf("%w: at least one KPI is required", ErrInvalidProjectMetadata)
	}

	if err := pm.validateDates(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidProjectMetadata, err)
	}

	return nil
}

// validateDates ensures project dates are logical
func (pm *ProjectMetadata) validateDates() error {
	if pm.StartDate.IsZero() || pm.EndDate.IsZero() {
		return errors.New("start and end dates must be set")
	}

	if pm.StartDate.After(pm.EndDate) {
		return errors.New("start date cannot be after end date")
	}

	return nil
}

// FormatContext returns a formatted string suitable for use in prompts
func (pm *ProjectMetadata) FormatContext() string {
	context := fmt.Sprintf("Project: %s\nDescription: %s\n\nBusiness Goals:\n", pm.Name, pm.Description)

	for _, goal := range pm.BusinessGoals {
		context += fmt.Sprintf("- %s\n", goal)
	}

	context += "\nKey Performance Indicators:\n"
	for _, kpi := range pm.KPIs {
		context += fmt.Sprintf("- %s\n", kpi)
	}

	context += fmt.Sprintf("\nProject Timeline: %s to %s",
		pm.StartDate.Format("2006-01-02"),
		pm.EndDate.Format("2006-01-02"))

	return context
}
