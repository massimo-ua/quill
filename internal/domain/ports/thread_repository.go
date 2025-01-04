package ports

import (
	"context"
	"github.com/massimo-ua/quill/internal/domain"
	"github.com/massimo-ua/quill/internal/domain/common"
)

// ThreadRepository defines interface for thread persistence
type ThreadRepository interface {
	// Save persists a thread
	Save(ctx context.Context, thread *domain.Thread) error

	// FindByID retrieves a thread by ID
	FindByID(ctx context.Context, id common.ID) (*domain.Thread, error)

	// Update updates thread information
	Update(ctx context.Context, thread *domain.Thread) error

	// Delete removes a thread
	Delete(ctx context.Context, id common.ID) error

	// FindByProject retrieves all threads for a project
	FindByProject(ctx context.Context, projectID common.ID) ([]*domain.Thread, error)
}
