package openai

import (
	"fmt"
	"github.com/massimo-ua/quill/internal/domain/ports"
)

// NewOpenAIProvider creates a new AiAgentProvider that uses OpenAI
func NewOpenAIProvider(cfg *Config) (ports.AiAgentProvider, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	client, err := NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI client: %w", err)
	}

	return NewProvider(client), nil
}