package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/massimo-ua/quill/internal/domain/common"
)

var (
	// ErrInvalidProjectName indicates that the project name is invalid
	ErrInvalidProjectName = errors.New("invalid project name")
	// ErrInvalidProjectGoals indicates that the project goals are invalid
	ErrInvalidProjectGoals = errors.New("invalid project goals")
	// ErrDuplicateMilestone indicates that a milestone with the same name already exists
	ErrDuplicateMilestone = errors.New("milestone with this name already exists")
)

// Project represents a project entity in the system
type Project struct {
	id          common.ID
	name        string
	description string
	goals       []string
	kpis        []string
	milestones  []Milestone
	createdAt   time.Time
	updatedAt   time.Time
}

// Milestone represents a project milestone
type Milestone struct {
	name     string
	deadline time.Time
}

// NewProject creates a new Project instance
func NewProject(name, description string, goals []string) (*Project, error) {
	if err := validateProjectName(name); err != nil {
		return nil, err
	}
	if err := validateProjectGoals(goals); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Project{
		id:          common.GenerateID(),
		name:        strings.TrimSpace(name),
		description: strings.TrimSpace(description),
		goals:       goals,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// MustNewProject creates a new Project instance and panics if validation fails
func MustNewProject(name, description string, goals []string) *Project {
	p, err := NewProject(name, description, goals)
	if err != nil {
		panic(err)
	}
	return p
}

// ID returns the project's identifier
func (p *Project) ID() common.ID {
	return p.id
}

// Name returns the project's name
func (p *Project) Name() string {
	return p.name
}

// Description returns the project's description
func (p *Project) Description() string {
	return p.description
}

// Goals returns the project's goals
func (p *Project) Goals() []string {
	goals := make([]string, len(p.goals))
	copy(goals, p.goals)
	return goals
}

// KPIs returns the project's KPIs
func (p *Project) KPIs() []string {
	kpis := make([]string, len(p.kpis))
	copy(kpis, p.kpis)
	return kpis
}

// Milestones returns the project's milestones
func (p *Project) Milestones() []Milestone {
	milestones := make([]Milestone, len(p.milestones))
	copy(milestones, p.milestones)
	return milestones
}

// CreatedAt returns the project's creation timestamp
func (p *Project) CreatedAt() time.Time {
	return p.createdAt
}

// UpdatedAt returns the project's last update timestamp
func (p *Project) UpdatedAt() time.Time {
	return p.updatedAt
}

// AddKPI adds a new KPI to the project
func (p *Project) AddKPI(kpi string) {
	kpi = strings.TrimSpace(kpi)
	if kpi != "" {
		p.kpis = append(p.kpis, kpi)
		p.updatedAt = time.Now()
	}
}

// AddMilestone adds a new milestone to the project
func (p *Project) AddMilestone(name string, deadline time.Time) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("milestone name cannot be empty")
	}

	// Check for duplicate milestone names
	for _, m := range p.milestones {
		if strings.EqualFold(m.name, name) {
			return ErrDuplicateMilestone
		}
	}

	milestone := Milestone{
		name:     name,
		deadline: deadline,
	}
	p.milestones = append(p.milestones, milestone)
	p.updatedAt = time.Now()
	return nil
}

// UpdateDescription updates the project's description
func (p *Project) UpdateDescription(description string) {
	p.description = strings.TrimSpace(description)
	p.updatedAt = time.Now()
}

// UpdateGoals updates the project's goals
func (p *Project) UpdateGoals(goals []string) error {
	if err := validateProjectGoals(goals); err != nil {
		return err
	}
	p.goals = goals
	p.updatedAt = time.Now()
	return nil
}

// validation helpers
func validateProjectName(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrInvalidProjectName
	}
	return nil
}

func validateProjectGoals(goals []string) error {
	if len(goals) == 0 {
		return ErrInvalidProjectGoals
	}
	for _, goal := range goals {
		if strings.TrimSpace(goal) == "" {
			return ErrInvalidProjectGoals
		}
	}
	return nil
}
