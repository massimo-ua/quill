package ports

import (
	"context"
	"github.com/massimo-ua/quill/internal/domain"
	"github.com/massimo-ua/quill/internal/domain/common"
)

// ChatAccessProvider defines interface for chat platform interactions
type ChatAccessProvider interface {
	// SendMessage sends a message to a channel
	SendMessage(ctx context.Context, channelID, content string) error

	// ReplyToMessage replies to a specific message
	ReplyToMessage(ctx context.Context, messageID, content string) error

	// ListenForMessages returns a channel for receiving messages
	ListenForMessages(ctx context.Context) (<-chan *domain.Message, error)

	// HandleInteraction processes interactive components
	HandleInteraction(ctx context.Context, interaction interface{}) error
}

// DocumentStoreProvider defines interface for document storage operations
type DocumentStoreProvider interface {
	// StoreDocument stores a new document
	StoreDocument(ctx context.Context, path string, content []byte, metadata map[string]interface{}) error

	// GetDocument retrieves a document
	GetDocument(ctx context.Context, path string) ([]byte, error)

	// UpdateDocument updates an existing document
	UpdateDocument(ctx context.Context, path string, content []byte, metadata map[string]interface{}) error

	// ListDocuments lists documents in a path
	ListDocuments(ctx context.Context, path string) ([]string, error)

	// DeleteDocument deletes a document
	DeleteDocument(ctx context.Context, path string) error
}

// AiAgentProvider defines interface for AI operations
type AiAgentProvider interface {
	// AnalyzeMessage analyzes message content
	AnalyzeMessage(ctx context.Context, content string) (*domain.MessageAnalysisResult, error)

	// GenerateDocumentation generates documentation from message
	GenerateDocumentation(ctx context.Context, message string, metadata map[string]interface{}) (string, error)

	// CategorizeContent categorizes content
	CategorizeContent(ctx context.Context, content string) (*domain.Category, error)

	// DetectReferences finds references in content
	DetectReferences(ctx context.Context, content string) ([]*domain.Reference, error)
}

// ProjectRepository defines interface for project persistence
type ProjectRepository interface {
	// Save persists a project
	Save(ctx context.Context, project *domain.Project) error

	// FindByID retrieves a project by ID
	FindByID(ctx context.Context, id common.ID) (*domain.Project, error)

	// Update updates project information
	Update(ctx context.Context, project *domain.Project) error

	// Delete removes a project
	Delete(ctx context.Context, id common.ID) error
}

// MessageRepository defines interface for message persistence
type MessageRepository interface {
	// Save persists a message
	Save(ctx context.Context, message *domain.Message) error

	// FindByID retrieves a message by ID
	FindByID(ctx context.Context, id string) (*domain.Message, error)

	// FindByThread retrieves messages in a thread
	FindByThread(ctx context.Context, threadID string) ([]*domain.Message, error)
}
