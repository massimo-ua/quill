package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewProject(t *testing.T) {
	t.Run("creates valid project", func(t *testing.T) {
		name := "Test Project"
		desc := "Test Description"
		goals := []string{"Goal 1", "Goal 2"}

		project, err := NewProject(name, desc, goals)

		assert.NoError(t, err)
		assert.NotNil(t, project)
		assert.NotEmpty(t, project.ID())
		assert.Equal(t, name, project.Name())
		assert.Equal(t, desc, project.Description())
		assert.Equal(t, goals, project.Goals())
		assert.NotZero(t, project.CreatedAt())
		assert.Equal(t, project.CreatedAt(), project.UpdatedAt())
	})

	t.Run("fails with empty name", func(t *testing.T) {
		project, err := NewProject("", "desc", []string{"goal"})

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidProjectName)
		assert.Nil(t, project)
	})

	t.Run("fails with empty goals", func(t *testing.T) {
		project, err := NewProject("name", "desc", []string{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidProjectGoals)
		assert.Nil(t, project)
	})

	t.Run("fails with empty goal in goals slice", func(t *testing.T) {
		project, err := NewProject("name", "desc", []string{"goal1", ""})

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidProjectGoals)
		assert.Nil(t, project)
	})
}

func TestMustNewProject(t *testing.T) {
	t.Run("creates valid project", func(t *testing.T) {
		assert.NotPanics(t, func() {
			project := MustNewProject("name", "desc", []string{"goal"})
			assert.NotNil(t, project)
		})
	})

	t.Run("panics with invalid input", func(t *testing.T) {
		assert.Panics(t, func() {
			MustNewProject("", "desc", []string{"goal"})
		})
	})
}

func TestProject_AddKPI(t *testing.T) {
	project := MustNewProject("name", "desc", []string{"goal"})
	originalTime := project.UpdatedAt()

	time.Sleep(time.Millisecond) // Ensure time difference

	t.Run("adds valid KPI", func(t *testing.T) {
		kpi := "New KPI"
		project.AddKPI(kpi)

		assert.Contains(t, project.KPIs(), kpi)
		assert.True(t, project.UpdatedAt().After(originalTime))
	})

	t.Run("ignores empty KPI", func(t *testing.T) {
		originalKPIs := project.KPIs()
		project.AddKPI("")

		assert.Equal(t, originalKPIs, project.KPIs())
	})
}

func TestProject_AddMilestone(t *testing.T) {
	project := MustNewProject("name", "desc", []string{"goal"})
	deadline := time.Now().Add(24 * time.Hour)

	t.Run("adds valid milestone", func(t *testing.T) {
		err := project.AddMilestone("Milestone 1", deadline)

		assert.NoError(t, err)
		milestones := project.Milestones()
		assert.Len(t, milestones, 1)
		assert.Equal(t, "Milestone 1", milestones[0].name)
		assert.Equal(t, deadline, milestones[0].deadline)
	})

	t.Run("prevents duplicate milestone names", func(t *testing.T) {
		err := project.AddMilestone("Milestone 1", deadline)

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrDuplicateMilestone)
		assert.Len(t, project.Milestones(), 1)
	})

	t.Run("fails with empty name", func(t *testing.T) {
		err := project.AddMilestone("", deadline)

		assert.Error(t, err)
		assert.Len(t, project.Milestones(), 1)
	})
}

func TestProject_UpdateDescription(t *testing.T) {
	project := MustNewProject("name", "desc", []string{"goal"})
	originalTime := project.UpdatedAt()

	time.Sleep(time.Millisecond) // Ensure time difference

	t.Run("updates description", func(t *testing.T) {
		newDesc := "New Description"
		project.UpdateDescription(newDesc)

		assert.Equal(t, newDesc, project.Description())
		assert.True(t, project.UpdatedAt().After(originalTime))
	})

	t.Run("trims whitespace", func(t *testing.T) {
		project.UpdateDescription("  trimmed  ")
		assert.Equal(t, "trimmed", project.Description())
	})
}

func TestProject_UpdateGoals(t *testing.T) {
	project := MustNewProject("name", "desc", []string{"goal"})
	originalTime := project.UpdatedAt()

	time.Sleep(time.Millisecond) // Ensure time difference

	t.Run("updates valid goals", func(t *testing.T) {
		newGoals := []string{"New Goal 1", "New Goal 2"}
		err := project.UpdateGoals(newGoals)

		assert.NoError(t, err)
		assert.Equal(t, newGoals, project.Goals())
		assert.True(t, project.UpdatedAt().After(originalTime))
	})

	t.Run("fails with empty goals", func(t *testing.T) {
		err := project.UpdateGoals([]string{})

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidProjectGoals)
	})

	t.Run("fails with empty goal in slice", func(t *testing.T) {
		err := project.UpdateGoals([]string{"goal", ""})

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidProjectGoals)
	})
}

func TestProject_DataImmutability(t *testing.T) {
	t.Run("goals slice is immutable", func(t *testing.T) {
		originalGoals := []string{"Goal 1", "Goal 2"}
		project := MustNewProject("name", "desc", originalGoals)

		goals := project.Goals()
		goals[0] = "Modified Goal"

		assert.Equal(t, originalGoals, project.Goals())
	})

	t.Run("KPIs slice is immutable", func(t *testing.T) {
		project := MustNewProject("name", "desc", []string{"goal"})
		project.AddKPI("KPI 1")

		kpis := project.KPIs()
		kpis[0] = "Modified KPI"

		assert.Equal(t, "KPI 1", project.KPIs()[0])
	})

	t.Run("milestones slice is immutable", func(t *testing.T) {
		project := MustNewProject("name", "desc", []string{"goal"})
		deadline := time.Now().Add(24 * time.Hour)
		_ = project.AddMilestone("Milestone 1", deadline)

		milestones := project.Milestones()
		milestones[0].name = "Modified Milestone"

		assert.Equal(t, "Milestone 1", project.Milestones()[0].name)
	})
}
