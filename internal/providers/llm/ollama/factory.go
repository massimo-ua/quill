package ollama

import (
	"fmt"
	"github.com/massimo-ua/quill/internal/domain/ports"
)

// NewOllamaProvider creates a new AiAgentProvider that uses Ollama
func NewOllamaProvider(cfg *Config) (ports.AiAgentProvider, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	client, err := NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Ollama client: %w", err)
	}

	return NewProvider(client), nil
}