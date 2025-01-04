package services

import (
	"context"
	"fmt"
	"github.com/massimo-ua/quill/internal/domain"
	"github.com/massimo-ua/quill/internal/domain/ports"
	"path/filepath"
	"time"
)

type DocumentationService struct {
	docStore ports.DocumentStoreProvider
	aiAgent  ports.AiAgentProvider
}

func NewDocumentationService(docs ports.DocumentStoreProvider, ai ports.AiAgentProvider) *DocumentationService {
	if docs == nil {
		panic("docStore cannot be nil")
	}
	if ai == nil {
		panic("aiAgent cannot be nil")
	}
	return &DocumentationService{
		docStore: docs,
		aiAgent:  ai,
	}
}

// CreateDocumentation generates and stores documentation from a message
func (s *DocumentationService) CreateDocumentation(
	ctx context.Context,
	msgType domain.MessageType,
	category domain.Category,
	content string,
	references []*domain.Reference,
) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}

	// Generate documentation using AI
	metadata := map[string]interface{}{
		"type":       msgType.String(),
		"category":   category.String(),
		"created_at": time.Now().UTC(),
		"references": references,
	}

	doc, err := s.aiAgent.GenerateDocumentation(ctx, content, metadata)
	if err != nil {
		return fmt.Errorf("failed to generate documentation: %w", err)
	}

	// Store the documentation
	path := s.generatePath(msgType, category)
	if err := s.docStore.StoreDocument(ctx, path, []byte(doc), metadata); err != nil {
		return fmt.Errorf("failed to store documentation: %w", err)
	}

	return nil
}

// UpdateDocumentation updates existing documentation
func (s *DocumentationService) UpdateDocumentation(
	ctx context.Context,
	path string,
	content string,
	metadata map[string]interface{},
) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}

	// Update metadata with modification time
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	metadata["updated_at"] = time.Now().UTC()

	if err := s.docStore.UpdateDocument(ctx, path, []byte(content), metadata); err != nil {
		return fmt.Errorf("failed to update documentation: %w", err)
	}

	return nil
}

// GetDocumentation retrieves documentation by path
func (s *DocumentationService) GetDocumentation(
	ctx context.Context,
	path string,
) ([]byte, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}

	doc, err := s.docStore.GetDocument(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve documentation: %w", err)
	}

	return doc, nil
}

// ListDocumentation lists all documentation in a category
func (s *DocumentationService) ListDocumentation(
	ctx context.Context,
	category domain.Category,
) ([]string, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}

	basePath := filepath.Join("docs", category.String())
	docs, err := s.docStore.ListDocuments(ctx, basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to list documentation: %w", err)
	}

	return docs, nil
}

// generatePath creates the storage path for documentation
func (s *DocumentationService) generatePath(msgType domain.MessageType, category domain.Category) string {
	timestamp := time.Now().UTC().Format("20060102-150405")
	filename := fmt.Sprintf("%s-%s.md", msgType.String(), timestamp)
	return filepath.Join("docs", category.String(), filename)
}
