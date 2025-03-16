package github

import (
	"context"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// DocumentStoreProvider implements the domain.DocumentStoreProvider interface
// for storing documentation in a GitHub repository
type DocumentStoreProvider struct {
	client *Client
}

// NewDocumentStoreProvider creates a new GitHubDocumentStoreProvider
func NewDocumentStoreProvider(client *Client) *DocumentStoreProvider {
	if client == nil {
		panic("client cannot be nil")
	}
	return &DocumentStoreProvider{
		client: client,
	}
}

// StoreDocument implements the ports.DocumentStoreProvider.StoreDocument method
// It stores a document in a GitHub repository
func (p *DocumentStoreProvider) StoreDocument(ctx context.Context, path string, content []byte, metadata map[string]interface{}) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}

	// Ensure parent directory exists
	dirPath := filepath.Dir(path)
	if dirPath != "." && dirPath != "/" {
		if err := p.client.ensureDirectoryExists(ctx, dirPath); err != nil {
			return fmt.Errorf("failed to ensure directory exists: %w", err)
		}
	}

	// Build commit message
	message := "Add documentation"
	if metadata != nil {
		if msgType, ok := metadata["type"].(string); ok {
			message = fmt.Sprintf("Add %s documentation", msgType)
		}
		if category, ok := metadata["category"].(string); ok {
			message = fmt.Sprintf("%s (%s)", message, category)
		}
	}

	// Create content
	_, err := p.client.CreateContent(ctx, path, content, message)
	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}

	return nil
}

// GetDocument implements the ports.DocumentStoreProvider.GetDocument method
// It retrieves a document from a GitHub repository
func (p *DocumentStoreProvider) GetDocument(ctx context.Context, path string) ([]byte, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}

	content, err := p.client.GetContent(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	// Decode base64 content
	decoded, err := base64.StdEncoding.DecodeString(content.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to decode content: %w", err)
	}

	return decoded, nil
}

// UpdateDocument implements the ports.DocumentStoreProvider.UpdateDocument method
// It updates an existing document in a GitHub repository
func (p *DocumentStoreProvider) UpdateDocument(ctx context.Context, path string, content []byte, metadata map[string]interface{}) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}

	// Build commit message
	message := "Update documentation"
	if metadata != nil {
		if msgType, ok := metadata["type"].(string); ok {
			message = fmt.Sprintf("Update %s documentation", msgType)
		}
		if category, ok := metadata["category"].(string); ok {
			message = fmt.Sprintf("%s (%s)", message, category)
		}
		if timestamp, ok := metadata["updated_at"].(time.Time); ok {
			message = fmt.Sprintf("%s at %s", message, timestamp.Format(time.RFC3339))
		}
	}

	// Update content
	_, err := p.client.UpdateContent(ctx, path, content, message)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	return nil
}

// ListDocuments implements the ports.DocumentStoreProvider.ListDocuments method
// It lists documents in a path in a GitHub repository
func (p *DocumentStoreProvider) ListDocuments(ctx context.Context, path string) ([]string, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}

	items, err := p.client.ListContents(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}

	// Extract file paths, ignoring .gitkeep files and directories
	var paths []string
	for _, item := range items {
		if item.Type == "file" && !strings.HasSuffix(item.Name, ".gitkeep") {
			paths = append(paths, item.Path)
		}
	}

	return paths, nil
}

// DeleteDocument implements the ports.DocumentStoreProvider.DeleteDocument method
// It deletes a document from a GitHub repository
func (p *DocumentStoreProvider) DeleteDocument(ctx context.Context, path string) error {
	if ctx == nil {
		return fmt.Errorf("context cannot be nil")
	}

	_, err := p.client.DeleteContent(ctx, path, "Delete documentation "+path)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	return nil
}