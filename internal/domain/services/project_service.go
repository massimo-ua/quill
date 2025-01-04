package services

import (
	"context"
	"fmt"
	"github.com/massimo-ua/quill/internal/domain"
	"github.com/massimo-ua/quill/internal/domain/common"
	"github.com/massimo-ua/quill/internal/domain/ports"
)

type ProjectService struct {
	docStore    ports.DocumentStoreProvider
	projectRepo ports.ProjectRepository
}

func NewProjectService(docs ports.DocumentStoreProvider, repo ports.ProjectRepository) *ProjectService {
	return &ProjectService{
		docStore:    docs,
		projectRepo: repo,
	}
}

func (s *ProjectService) CreateProject(ctx context.Context, metadata *domain.ProjectMetadata) error {
	project, err := domain.NewProject(metadata.Name, metadata.Description, metadata.BusinessGoals)
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	// Add KPIs and dates from metadata
	for _, kpi := range metadata.KPIs {
		project.AddKPI(kpi)
	}

	if err := s.projectRepo.Save(ctx, project); err != nil {
		return fmt.Errorf("failed to save project: %w", err)
	}

	// Generate project documentation content
	docContent := fmt.Sprintf("# %s\n\n## Description\n%s\n\n## Business Goals\n",
		project.Name(), project.Description())

	for _, goal := range project.Goals() {
		docContent += fmt.Sprintf("* %s\n", goal)
	}

	docContent += "\n## KPIs\n"
	for _, kpi := range project.KPIs() {
		docContent += fmt.Sprintf("* %s\n", kpi)
	}

	docPath := fmt.Sprintf("projects/%s/README.md", project.ID())
	if err := s.docStore.StoreDocument(ctx, docPath, []byte(docContent), nil); err != nil {
		return fmt.Errorf("failed to store project documentation: %w", err)
	}

	return nil
}

func (s *ProjectService) GetProject(ctx context.Context, id common.ID) (*domain.Project, error) {
	return s.projectRepo.FindByID(ctx, id)
}

func (s *ProjectService) UpdateProject(ctx context.Context, project *domain.Project) error {
	if err := s.projectRepo.Update(ctx, project); err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	// Update documentation with current project state
	docContent := fmt.Sprintf("# %s\n\n## Description\n%s\n\n## Business Goals\n",
		project.Name(), project.Description())

	for _, goal := range project.Goals() {
		docContent += fmt.Sprintf("* %s\n", goal)
	}

	docContent += "\n## KPIs\n"
	for _, kpi := range project.KPIs() {
		docContent += fmt.Sprintf("* %s\n", kpi)
	}

	docPath := fmt.Sprintf("projects/%s/README.md", project.ID())
	if err := s.docStore.UpdateDocument(ctx, docPath, []byte(docContent), nil); err != nil {
		return fmt.Errorf("failed to update project documentation: %w", err)
	}

	return nil
}
