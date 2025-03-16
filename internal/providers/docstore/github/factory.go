package github

import (
	"fmt"
	"github.com/massimo-ua/quill/internal/domain/ports"
)

// NewGitHubDocumentStoreProvider creates a new DocumentStoreProvider backed by GitHub
func NewGitHubDocumentStoreProvider(config *Config) (ports.DocumentStoreProvider, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	client, err := NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub client: %w", err)
	}

	return NewDocumentStoreProvider(client), nil
}